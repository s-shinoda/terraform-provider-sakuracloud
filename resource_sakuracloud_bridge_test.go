package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudBridge_Basic(t *testing.T) {
	var bridge sacloud.Bridge
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudBridgeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudBridgeConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeExists("sakuracloud_bridge.foobar", &bridge),
					resource.TestCheckResourceAttr(
						"sakuracloud_bridge.foobar", "name", "mybridge"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bridge.foobar", "switch_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccSakuraCloudBridge_Update(t *testing.T) {
	var bridge sacloud.Bridge
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudBridgeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudBridgeConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeExists("sakuracloud_bridge.foobar", &bridge),
					resource.TestCheckResourceAttr(
						"sakuracloud_bridge.foobar", "name", "mybridge"),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudBridgeConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBridgeExists("sakuracloud_bridge.foobar", &bridge),
					resource.TestCheckResourceAttr(
						"sakuracloud_bridge.foobar", "name", "mybridge_upd"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudBridgeExists(n string, bridge *sacloud.Bridge) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Bridge ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)
		originalZone := client.Zone
		client.Zone = "is1a"
		defer func() { client.Zone = originalZone }()

		foundBridge, err := client.Bridge.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundBridge.ID != rs.Primary.ID {
			return fmt.Errorf("Bridge not found")
		}

		*bridge = *foundBridge

		return nil
	}
}

func testAccCheckSakuraCloudBridgeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "is1a"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_bridge" {
			continue
		}

		_, err := client.Bridge.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Bridge still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudBridgeConfig_basic = `
resource "sakuracloud_bridge" "foobar" {
    name = "mybridge"
    description = "Bridge from TerraForm for SAKURA CLOUD"
    zone = "is1a"
}`

var testAccCheckSakuraCloudBridgeConfig_update = `
resource "sakuracloud_bridge" "foobar" {
    name = "mybridge_upd"
    description = "Bridge from TerraForm for SAKURA CLOUD"
    zone = "is1a"
}`