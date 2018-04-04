package hue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lawsontyler/ghue/sdk/common"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema {
			"hub_address": {
				Type: schema.TypeString,
				Required: true,
				Description: "Address of your Philips Hue Hub",
			},
			"hub_username": {
				Type: schema.TypeString,
				Required: true,
				Description: "Username on your Hub.  See Hue API for details on how to create this.",
			},
			"verbose": {
				Type: schema.TypeBool,
				Optional: true,
				Default: false,
			},
		},

		ResourcesMap: map[string]*schema.Resource {
			"philips-hue_scene": resourceScene(),
			"philips-hue_group": resourceGroup(),
			"philips-hue_rule": resourceRule(),
		},

		DataSourcesMap: map[string]*schema.Resource {
			"philips-hue_light": dataSourceHueLight(),
			"philips-hue_sensor": dataSourceHueSensor(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	connection := &common.Connection{
		Host:     d.Get("hub_address").(string),
		Username: d.Get("hub_username").(string),
		Verbose:  d.Get("verbose").(bool),
	}

	return connection, nil
}