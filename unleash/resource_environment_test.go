package unleash

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUnleashEnvironment_basic(t *testing.T) {
	//resourceName := "unleash_environment.test"
	rName := "qa-dcs"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:            testAccUnleashEnvironment_basic(rName),
				ResourceName:      "unleash_environment.test1",
				ImportState:       true,
				ImportStateId:     rName,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccUnleashEnvironment_default(t *testing.T) {
	resourceName := "unleash_environment.test2"
	rName := "qa-test"
	typeUpdate := "preproduction"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUnleashEnvironment_default(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "test"),
				),
			},
			{
				Config: testAccUnleashEnvironment_update(rName, typeUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", typeUpdate),
				),
			},
		},
	})
}

func testAccCheckEnvironmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Environment Not found: %s", n)
		}

		return nil
	}
}

func testAccUnleashEnvironment_basic(rName string) string {
	return fmt.Sprintf(`
resource "unleash_environment" "test1" {
	name       = %q
}
`, rName)
}

func testAccUnleashEnvironment_default(rName string) string {
	return fmt.Sprintf(`
resource "unleash_environment" "test2" {
	name       = %q
	type       = "test"
	enabled    = true
}
`, rName)
}

func testAccUnleashEnvironment_update(rName string, envType string) string {
	return fmt.Sprintf(`
resource "unleash_environment" "test2" {
	name       = %q
	type       = %q
	enabled	   = true
}
`, rName, envType)
}
