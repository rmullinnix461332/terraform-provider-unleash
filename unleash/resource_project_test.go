package unleash

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUnleashProject_basic(t *testing.T) {
	//resourceName := "unleash_project.test"
	rId := "slfus-client-onboard"
	rName := "SLFUS Client Onboard"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:            testAccUnleashProject_basic(rId, rName),
				ResourceName:      "unleash_project.test1",
				ImportState:       true,
				ImportStateId:     rId,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccUnleashProject_default(t *testing.T) {
	resourceName := "unleash_project.test2"
	rId := acctest.RandomWithPrefix("tf-acc-test")
	rName := "tf acc test"
	descUpdate := "My Project Updated"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUnleashProject_default(rId, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "My Project"),
				),
			},
			{
				Config: testAccUnleashProject_update(rId, rName, descUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", descUpdate),
				),
			},
		},
	})
}

func testAccCheckProjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Project Not found: %s", n)
		}

		return nil
	}
}

func testAccUnleashProject_basic(rId string, rName string) string {
	return fmt.Sprintf(`
resource "unleash_project" "test1" {
	project_id = %q
	name       = %q

	environments = ["stage-dcs", "qa-dcs", "prod-dcs"]
}
`, rId, rName)
}

func testAccUnleashProject_default(rId string, rName string) string {
	return fmt.Sprintf(`
resource "unleash_project" "test2" {
	project_id    = %q
	name          = %q
	description   = "My Project"

	environments = ["stage-dcs", "qa-dcs"]
}
`, rId, rName)
}

func testAccUnleashProject_update(rId string, rName string, descUpdate string) string {
	return fmt.Sprintf(`
resource "unleash_project" "test2" {
	project_id    = %q
	name          = %q
	description   = %q

	environments = ["stage-dcs", "qa-dcs"]
}
`, rId, rName, descUpdate)
}
