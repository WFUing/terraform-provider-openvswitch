package openvswitch

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/hashicorp/terraform/helper/schema"
)

// Resource Definition
func resourceBridge() *schema.Resource {
	return &schema.Resource{
		Create: resourceBridgeCreate,
		Read:   resourceBridgeRead,
		Update: resourceBridgeUpdate,
		Delete: resourceBridgeDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// "ofversion": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// 	Default:  "OpenFlow13",
			// },
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// 通过 sudo ovs-vsctl show 鉴权
func checkPermissions() error {
	cmd := exec.Command("sudo", "ovs-vsctl", "show")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("insufficient permissions to run ovs-vsctl: %v", err)
	}
	return nil
}

// 通过sudo运行ovs-vsctl
func runCommandWithSudo(command string, args ...string) error {
	cmdArgs := append([]string{command}, args...)
	cmd := exec.Command("sudo", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run command %s: %v, output: %s", command, err, string(output))
	}
	return nil
}

func resourceBridgeCreate(d *schema.ResourceData, m interface{}) error {
	if err := checkPermissions(); err != nil {
		return err
	}

	bridge := d.Get("name").(string)
	ipAddress := d.Get("ip_address").(string)
	log.Printf("[DEBUG] Creating bridge: %s with ip %s", bridge, ipAddress)

	// 如果ovs已经存在，先删除
	if exists, _ := bridgeExists(nil, bridge); exists {
		log.Printf("[DEBUG] Bridge %s already exists, deleting it first", bridge)
		if err := resourceBridgeDelete(d, m); err != nil {
			return fmt.Errorf("failed to delete existing bridge %s: %v", bridge, err)
		}
	}

	// 创建ovs
	if err := runCommandWithSudo("ovs-vsctl", "add-br", bridge); err != nil {
		return fmt.Errorf("failed to create bridge %s: %v", bridge, err)
	}

	// ver := []string{d.Get("ofversion").(string)}
	// client := ovs.New()
	// if err := client.VSwitch.Set.Bridge(bridge, ovs.BridgeOptions{Protocols: ver}); err != nil {
	// 	return fmt.Errorf("failed to set bridge options for %s: %v", bridge, err)
	// }

	// 设置网络接口和IP地址
	if err := runCommandWithSudo("ip", "link", "set", bridge, "up"); err != nil {
		return fmt.Errorf("failed to set bridge %s up: %v", bridge, err)
	}

	if err := runCommandWithSudo("ip", "addr", "add", ipAddress, "dev", bridge); err != nil {
		return fmt.Errorf("failed to add IP address %s to bridge %s: %v", ipAddress, bridge, err)
	}

	d.SetId(bridge)
	return resourceBridgeRead(d, m)
}

func resourceBridgeRead(d *schema.ResourceData, m interface{}) error {
	bridge := d.Id()
	log.Printf("[DEBUG] Reading bridge: %s", bridge)

	client := ovs.New()
	exists, err := bridgeExists(client, bridge)
	if err != nil {
		return fmt.Errorf("error checking if bridge exists: %v", err)
	}

	if !exists {
		d.SetId("")
	}

	return nil
}

func bridgeExists(_ *ovs.Client, bridge string) (bool, error) {
	if err := runCommandWithSudo("ovs-vsctl", "br-exists", bridge); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceBridgeUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceBridgeRead(d, m)
}

func resourceBridgeDelete(d *schema.ResourceData, m interface{}) error {
	bridge := d.Get("name").(string)
	log.Printf("[DEBUG] Deleting bridge: %s", bridge)
	if err := runCommandWithSudo("ovs-vsctl", "del-br", bridge); err != nil {
		return fmt.Errorf("failed to delete bridge %s: %v", bridge, err)
	}
	d.SetId("")
	return nil
}
