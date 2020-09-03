package cloudflare

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	zoneID = os.Getenv("CLOUDFLARE_ZONE_ID")
	domain = os.Getenv("CLOUDFLARE_DOMAIN")
)

func TestAccCloudflareAccessApplicationBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := fmt.Sprintf("cloudflare_access_application.%s", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudflareAccessApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareAccessApplicationConfigBasic(rnd, zoneID, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "zone_id", zoneID),
					resource.TestCheckResourceAttr(name, "name", rnd),
					resource.TestCheckResourceAttr(name, "domain", fmt.Sprintf("%s.%s", rnd, domain)),
					resource.TestCheckResourceAttr(name, "session_duration", "24h"),
					resource.TestCheckResourceAttr(name, "cors_headers.#", "0"),
					resource.TestCheckResourceAttr(name, "auto_redirect_to_identity", "false"),
				),
			},
		},
	})
}

func TestAccCloudflareAccessApplicationWithCORS(t *testing.T) {
	rnd := generateRandomResourceName()
	name := fmt.Sprintf("cloudflare_access_application.%s", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudflareAccessApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareAccessApplicationConfigWithCORS(rnd, zoneID, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "zone_id", zoneID),
					resource.TestCheckResourceAttr(name, "name", rnd),
					resource.TestCheckResourceAttr(name, "domain", fmt.Sprintf("%s.%s", rnd, domain)),
					resource.TestCheckResourceAttr(name, "session_duration", "24h"),
					resource.TestCheckResourceAttr(name, "cors_headers.#", "1"),
					resource.TestCheckResourceAttr(name, "cors_headers.0.allowed_methods.#", "3"),
					resource.TestCheckResourceAttr(name, "cors_headers.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(name, "cors_headers.0.max_age", "10"),
					resource.TestCheckResourceAttr(name, "auto_redirect_to_identity", "false"),
				),
			},
		},
	})
}

func TestAccCloudflareAccessApplicationWithAutoRedirectToIdentity(t *testing.T) {
	rnd := generateRandomResourceName()
	name := fmt.Sprintf("cloudflare_access_application.%s", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudflareAccessApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareAccessApplicationConfigWithAutoRedirectToIdentity(rnd, zoneID, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "zone_id", zoneID),
					resource.TestCheckResourceAttr(name, "name", rnd),
					resource.TestCheckResourceAttr(name, "domain", fmt.Sprintf("%s.%s", rnd, domain)),
					resource.TestCheckResourceAttr(name, "session_duration", "24h"),
					resource.TestCheckResourceAttr(name, "auto_redirect_to_identity", "true"),
				),
			},
		},
	})
}

func TestAccCloudflareAccessApplicationWithADefinedIdps(t *testing.T) {
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := generateRandomResourceName()
	name := fmt.Sprintf("cloudflare_access_application.%s", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudflareAccessApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareAccessApplicationConfigWithADefinedIdp(rnd, zoneID, domain, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "zone_id", zoneID),
					resource.TestCheckResourceAttr(name, "name", rnd),
					resource.TestCheckResourceAttr(name, "domain", fmt.Sprintf("%s.%s", rnd, domain)),
					resource.TestCheckResourceAttr(name, "session_duration", "24h"),
					resource.TestCheckResourceAttr(name, "auto_redirect_to_identity", "true"),
					resource.TestCheckResourceAttr(name, "allowed_idps.#", "1"),
				),
			},
		},
	})
}

func testAccCloudflareAccessApplicationConfigBasic(rnd, zoneID, domain string) string {
	return fmt.Sprintf(`
resource "cloudflare_access_application" "%[1]s" {
  zone_id                   = "%[2]s"
  name                      = "%[1]s"
  domain                    = "%[1]s.%[3]s"
  session_duration          = "24h"
  auto_redirect_to_identity = false
}
`, rnd, zoneID, domain)
}

func testAccCloudflareAccessApplicationConfigWithCORS(rnd, zoneID, domain string) string {
	return fmt.Sprintf(`
resource "cloudflare_access_application" "%[1]s" {
  zone_id          = "%[2]s"
  name             = "%[1]s"
  domain           = "%[1]s.%[3]s"
  session_duration = "24h"
  cors_headers {
    allowed_methods = ["GET", "POST", "OPTIONS"]
    allowed_origins = ["https://example.com"]
    allow_credentials = true
    max_age = 10
  }
  auto_redirect_to_identity = false
}
`, rnd, zoneID, domain)
}

func testAccCloudflareAccessApplicationConfigWithAutoRedirectToIdentity(rnd, zoneID, domain string) string {
	return fmt.Sprintf(`
resource "cloudflare_access_application" "%[1]s" {
  zone_id                   = "%[2]s"
  name                      = "%[1]s"
  domain                    = "%[1]s.%[3]s"
  session_duration          = "24h"
  auto_redirect_to_identity = true
}
`, rnd, zoneID, domain)
}

func testAccCloudflareAccessApplicationConfigWithADefinedIdp(rnd, zoneID, domain string, accountID string) string {
	return fmt.Sprintf(`
resource "cloudflare_access_identity_provider" "%[1]s" {
  account_id = "%[4]s"
  name = "%[1]s"
  type = "onetimepin"
}
resource "cloudflare_access_application" "%[1]s" {
  zone_id                   = "%[2]s"
  name                      = "%[1]s"
  domain                    = "%[1]s.%[3]s"
  session_duration          = "24h"
  auto_redirect_to_identity = true
  allowed_idps              = [cloudflare_access_identity_provider.%[1]s.id]
}
`, rnd, zoneID, domain, accountID)
}

func testAccCheckCloudflareAccessApplicationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudflare.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudflare_access_application" {
			continue
		}

		_, err := client.AccessApplication(rs.Primary.Attributes["zone_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("AccessApplication still exists")
		}
	}

	return nil
}

func TestAccCloudflareAccessApplicationWithZoneID(t *testing.T) {
	// Temporarily unset CLOUDFLARE_API_TOKEN if it is set as the Access
	// service does not yet support the API tokens and it results in
	// misleading state error messages.
	if os.Getenv("CLOUDFLARE_API_TOKEN") != "" {
		defer func(apiToken string) {
			os.Setenv("CLOUDFLARE_API_TOKEN", apiToken)
		}(os.Getenv("CLOUDFLARE_API_TOKEN"))
		os.Setenv("CLOUDFLARE_API_TOKEN", "")
	}

	rnd := generateRandomResourceName()
	name := "cloudflare_access_application." + rnd
	zone := os.Getenv("CLOUDFLARE_DOMAIN")
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	updatedName := fmt.Sprintf("%s-updated", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAccount(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccessApplicationWithZoneID(rnd, zone, zoneID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rnd),
					resource.TestCheckResourceAttr(name, "account_id", accountID),
				),
			},
			{
				Config: testAccessApplicationWithZoneIDUpdated(rnd, zone, zoneID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", updatedName),
					resource.TestCheckResourceAttr(name, "account_id", accountID),
				),
			},
		},
	})
}

func testAccessApplicationWithZoneID(resourceID, zone, zoneID string) string {
	return fmt.Sprintf(`
		resource "cloudflare_access_application" "%[1]s" {
			name    = "%[1]s"
			zone_id = "%[3]s"
			domain  = "%[1]s.%[2]s"
		}
	`, resourceID, zone, zoneID)
}

func testAccessApplicationWithZoneIDUpdated(resourceID, zone, zoneID string) string {
	return fmt.Sprintf(`
		resource "cloudflare_access_application" "%[1]s" {
			name    = "%[1]s-updated"
			zone_id = "%[3]s"
			domain  = "%[1]s.%[2]s"
		}
	`, resourceID, zone, zoneID)
}