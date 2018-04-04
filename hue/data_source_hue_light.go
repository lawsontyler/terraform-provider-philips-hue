package hue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lawsontyler/ghue/sdk/common"
	"github.com/lawsontyler/ghue/sdk/lights"
)

func dataSourceHueLight() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHueLightRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				ConflictsWith: []string{"light_id"},
				Optional: true,
			},
			"light_id": {
				Type: schema.TypeString,
				ConflictsWith: []string{"name"},
				Optional: true,
			},
		},
	}
}

func dataSourceHueLightRead(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*common.Connection)

	lightName := d.Get("name")
	lightId := d.Get("light_id")

	if lightId != nil {
		lightId := lightId.(string)
		_, _, err := lights.GetLight(connection, lightId)

		if err != nil {
			return err
		}

		d.SetId(lightId)

	} else if lightName != nil {
		lightName := lightName.(string)
		lightId, _, err := lights.GetLightIdByName(connection, lightName)

		if err != nil {
			return err
		}

		d.SetId(lightId)
	}

	return nil
}