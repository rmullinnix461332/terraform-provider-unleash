package unleash

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUnleashStrategy_basic(t *testing.T) {
	//resourceName := "unleash_strategy.test"
	rName := "default"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:            testAccUnleashStrategy_basic(rName),
				ResourceName:      "unleash_strategy.test1",
				ImportState:       true,
				ImportStateId:     rName,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccUnleashStrategy_default(t *testing.T) {
	resourceName := "unleash_strategy.test2"
	rName := "structuredRollout"
	descUpdate := "Test Strategy Updated"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUnleashStrategy_default(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStrategyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test Strategy Type"),
				),
			},
			{
				Config: testAccUnleashStrategy_update(rName, descUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStrategyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", descUpdate),
				),
			},
		},
	})
}

func testAccCheckStrategyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Strategy Not found: %s", n)
		}

		return nil
	}
}

func testAccUnleashStrategy_basic(rName string) string {
	return fmt.Sprintf(`
resource "unleash_strategy" "test1" {
	name       = %q
}
`, rName)
}

func testAccUnleashStrategy_default(rName string) string {
	return fmt.Sprintf(`
resource "unleash_strategy" "test2" {
	name        = %q
	description = "Test Strategy Type"
	enabled	    = true
}
`, rName)
}

func testAccUnleashStrategy_update(rName string, description string) string {
	return fmt.Sprintf(`
resource "unleash_strategy" "test2" {
	name        = %q
	description = %q
	enabled	    = true

	parameters {
			name        = "Test"
			type        = "string"
			description = "Test parameter"
			required    = false
	}
}
`, rName, description)
}
