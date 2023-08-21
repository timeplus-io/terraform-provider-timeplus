// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStreamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "timeplus_stream" "test" {
  name = "test_resource"
  description = "testing stream resource"
  column {
    name = "col_1"
    type = "string"
  }
  column {
    name = "col_2"
    type = "int64"
  }
  retention_size = 1024
  retention_period = 3600
  historical_data_ttl = "to_datetime(_tp_time) + INTERVAL 1 DAY"
}
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("timeplus_stream.test", "name", "test_resource"),
					resource.TestCheckResourceAttr("timeplus_stream.test", "description", "testing stream resource"),
					resource.TestCheckResourceAttr("timeplus_stream.test", "retention_size", "1024"),
					resource.TestCheckResourceAttr("timeplus_stream.test", "retention_period", "3600"),
					resource.TestCheckResourceAttr("timeplus_stream.test", "historical_data_ttl", "to_datetime(_tp_time) + INTERVAL 1 DAY"),
				),
			},
			// Update and Read testing
			{
				Config: `
resource "timeplus_stream" "test" {
  name = "test_resource"
  description = "testing stream resource update"
  column {
    name = "col_1"
    type = "string"
  }
  column {
    name = "col_2"
    type = "int64"
  }
  retention_size = 2048
  retention_period = 7200
  historical_data_ttl = "to_datetime(_tp_time) + INTERVAL 1 DAY"
}
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("timeplus_stream.test", "name", "test_resource"),
					resource.TestCheckResourceAttr("timeplus_stream.test", "description", "testing stream resource update"),
					resource.TestCheckResourceAttr("timeplus_stream.test", "retention_size", "2048"),
					resource.TestCheckResourceAttr("timeplus_stream.test", "retention_period", "7200"),
					resource.TestCheckResourceAttr("timeplus_stream.test", "historical_data_ttl", "to_datetime(_tp_time) + INTERVAL 1 DAY"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
