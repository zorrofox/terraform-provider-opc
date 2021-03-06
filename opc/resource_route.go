package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOPCRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCRouteCreate,
		Read:   resourceOPCRouteRead,
		Update: resourceOPCRouteUpdate,
		Delete: resourceOPCRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"admin_distance": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAdminDistance,
			},

			"ip_address_prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIPPrefixCIDR,
			},

			"next_hop_vnic_set": {
				Type:     schema.TypeString,
				Required: true,
			},

			"tags": tagsOptionalSchema(),
		},
	}
}

func resourceOPCRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).Routes()

	// Get Required attributes
	name := d.Get("name").(string)
	ipPrefix := d.Get("ip_address_prefix").(string)
	nextHop := d.Get("next_hop_vnic_set").(string)

	// Start populating input struct
	input := &compute.CreateRouteInput{
		Name:            name,
		IPAddressPrefix: ipPrefix,
		NextHopVnicSet:  nextHop,
	}

	// Get Optional Attributes
	desc, descOk := d.GetOk("description")
	if descOk {
		input.Description = desc.(string)
	}

	dist, distOk := d.GetOk("admin_distance")
	if distOk {
		input.AdminDistance = dist.(int)
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	// Create Route
	info, err := client.CreateRoute(input)
	if err != nil {
		return fmt.Errorf("Error creating route '%s': %v", name, err)
	}

	d.SetId(info.Name)

	return resourceOPCRouteRead(d, meta)
}

func resourceOPCRouteRead(d *schema.ResourceData, meta interface{}) error {
	computeClient := meta.(*compute.ComputeClient).Routes()

	name := d.Id()
	input := &compute.GetRouteInput{
		Name: name,
	}

	result, err := computeClient.GetRoute(input)
	if err != nil {
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading route '%s': %v", name, err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("admin_distance", result.AdminDistance)
	d.Set("ip_address_prefix", result.IPAddressPrefix)
	d.Set("next_hop_vnic_set", result.NextHopVnicSet)
	d.Set("description", result.Description)
	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}
	return nil
}

func resourceOPCRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).Routes()

	// Get Required attributes
	name := d.Get("name").(string)
	ipPrefix := d.Get("ip_address_prefix").(string)
	nextHop := d.Get("next_hop_vnic_set").(string)

	// Start populating input struct
	input := &compute.UpdateRouteInput{
		Name:            name,
		IPAddressPrefix: ipPrefix,
		NextHopVnicSet:  nextHop,
	}

	// Get Optional Attributes
	desc, descOk := d.GetOk("description")
	if descOk {
		input.Description = desc.(string)
	}

	dist, distOk := d.GetOk("admin_distance")
	if distOk {
		input.AdminDistance = dist.(int)
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	// Create Route
	info, err := client.UpdateRoute(input)
	if err != nil {
		return fmt.Errorf("Error creating route '%s': %v", name, err)
	}

	d.SetId(info.Name)

	return resourceOPCRouteRead(d, meta)
}

func resourceOPCRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).Routes()

	name := d.Id()
	input := &compute.DeleteRouteInput{
		Name: name,
	}
	if err := client.DeleteRoute(input); err != nil {
		return fmt.Errorf("Error deleting route '%s': %v", name, err)
	}
	return nil
}
