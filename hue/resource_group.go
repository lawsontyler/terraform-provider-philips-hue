package hue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lawsontyler/ghue/sdk/groups"
	"github.com/lawsontyler/ghue/sdk/common"
	"github.com/lawsontyler/terraform-provider-philips-hue/hue/lib/constants"
)


func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lights": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"type": {
				Type: schema.TypeString,
				Optional: true,
				Default:  "LightGroup",
			},
		},
	}
}

func dataToLightArray(lights *schema.Set) []string {
	var lightArray []string

	if v := lights; v.Len() > 0 {
		for _, v := range v.List() {
			lightArray = append(lightArray, v.(string))
		}
	}

	return lightArray
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	lights := dataToLightArray(d.Get("lights").(*schema.Set))

	group := groups.Create{
		Name: d.Get("name").(string),
		Lights: lights,
		Type: d.Get("type").(string),
	}

	result, _, err := groups.CreateAPI(connection, &group)

	if err != nil {
		return err
	}

	d.SetId(result.Success.Id)

	return nil
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	group, hueErr, err := groups.GetGroup(connection, d.Id())

	if err != nil && hueErr != nil && hueErr.Error.Type == int(constants.NOT_FOUND) {
		d.SetId("")
	}

	d.Set("name", group.Name)
	d.Set("lights", group.Lights)
	d.Set("type", group.Type)

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	lights := dataToLightArray(d.Get("lights").(*schema.Set))

	group := groups.Update{
		Name: d.Get("name").(string),
		Lights: lights,
	}

	_, _, err := groups.UpdateAPI(connection, d.Id(), &group)

	if err != nil {
		return err
	}

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	_, _, err := groups.DeleteAPI(connection, d.Id())

	if err != nil {
		return err
	}

	return nil
}