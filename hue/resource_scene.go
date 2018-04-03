package hue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lawsontyler/ghue/sdk/lights"
	"fmt"
	"github.com/lawsontyler/ghue/sdk/common"
	"github.com/lawsontyler/ghue/sdk/scenes"
	"strconv"
)


func resourceScene() *schema.Resource {
	return &schema.Resource{
		Create: resourceSceneCreate,
		Read:   resourceSceneRead,
		Update: resourceSceneUpdate,
		Delete: resourceSceneDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"recycle": {
				Type: schema.TypeBool,
				Optional: false,
				Default: true,
			},
			"light_state": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"light_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"bri": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{"xy", "ct"},
						},
						"hue": {
							Type: schema.TypeString,
							Optional: true,
							ConflictsWith: []string{"xy", "ct"},
						},
						"sat": {
							Type: schema.TypeString,
							Optional: true,
							ConflictsWith: []string{"xy", "ct"},
						},
						"xy": {
							Type: schema.TypeSet,
							Optional: true,
							ConflictsWith: []string{"bri", "hue", "sat", "ct"},
							Elem: &schema.Schema{Type: schema.TypeFloat},
							MaxItems: 2,
						},
						"ct": {
							Type: schema.TypeString,
							Optional: true,
							ConflictsWith: []string{"bri", "hue", "sat", "xy"},
						},
						"transitiontime": {
							Type: schema.TypeInt,
							Optional: true,
							Default: 4,
						},
					},
				},
			},
		},
	}
}

func resourceSceneCreate(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)
	d.Partial(true)

	// Step 1: Set the light states
	lightStates := d.Get("light_state").(*schema.Set).List()

	var lightsInScene []string

	for _, lightState := range lightStates {
		lightState := lightState.(map[string]interface{})

		stateValue := lights.SetStateValues{}

		if brightness := lightState["bri"]; brightness != nil {
			stateValue.Bri = brightness.(string)
		}

		if hue := lightState["hue"]; hue != nil {
			stateValue.Hue = hue.(string)
		}

		if saturation := lightState["sat"]; saturation != nil {
			stateValue.Sat = saturation.(string)
		}

		if xy := lightState["xy"]; xy != nil {
			xy := xy.([]float64)
			// The library expects two float64s as a string.  Neat.
			stateValue.XY = fmt.Sprintf("%0.4fx%0.4f", xy[0], xy[1])
		}

		stateValue.TransitionTime = "0"

		lightId := lightState["light_id"].(string)
		lightsInScene = append(lightsInScene, lightId)

		_, _, err := lights.SetState(connection, lightId , &stateValue)

		if err != nil {
			return err
		}

		d.SetPartial("light_state")
	}

	// Now that we know the lights in the scene and have set their state, we can create the scene and capture it.

	scene := scenes.Create{
		Name: d.Get("name").(string),
		Lights: lightsInScene,
		Recycle: d.Get("recycle").(bool),
	}

	sceneResult, _, err := scenes.CreateApi(connection, &scene)

	if err != nil {
		return err
	}

	d.Partial(false)

	d.SetId(strconv.Itoa(sceneResult.Success.Id))

	return nil
}

func resourceSceneRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSceneUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSceneDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}