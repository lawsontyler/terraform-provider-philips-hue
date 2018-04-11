package hue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lawsontyler/ghue/sdk/common"
	"github.com/lawsontyler/ghue/sdk/scenes"
	"fmt"
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
				Optional: true,
				Default: true,
			},
			"light_state": {
				Type: schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{

					Schema: map[string]*schema.Schema{
						"light_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"state": {
							Type: schema.TypeString,
							Optional: true,
							ValidateFunc: validateLightOnOffState,
						},

						"bri": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hue": {
							Type: schema.TypeString,
							Optional: true,
							ConflictsWith: []string{"light_state.xy", "light_state.ct"},
						},
						"sat": {
							Type: schema.TypeString,
							Optional: true,
							ConflictsWith: []string{"light_state.xy", "light_state.ct"},
						},
						"xy": {
							Type: schema.TypeSet,
							Optional: true,
							ConflictsWith: []string{"light_state.hue", "light_state.sat", "light_state.ct"},
							Elem: &schema.Schema{Type: schema.TypeFloat},
							MaxItems: 2,
						},
						"ct": {
							Type: schema.TypeString,
							Optional: true,
							ConflictsWith: []string{"light_state.hue", "light_state.sat", "light_state.xy"},
						},
						"transitiontime": {
							Type: schema.TypeString,
							Optional: true,
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
	// Seizure warning is especially important here.

	lightStates := d.Get("light_state").(*schema.Set).List()

	lightsInScene := getLightsInScene(lightStates)

	// Now that we know the lights in the scene and have set their state, we can create the scene and capture it.

	scene := scenes.Create{
		Name:    d.Get("name").(string),
		Lights:  lightsInScene,
		Recycle: d.Get("recycle").(bool),
	}

	sceneResult, _, err := scenes.CreateApi(connection, &scene)

	if err != nil {
		return err
	}

	d.SetPartial("name")

	// Now that it's created...and all the light states are wrong, let's run an update to set them.

	err = setLightStates(connection, sceneResult.Success.Id, lightStates)

	if err != nil {
		return err
	}

	d.Partial(false)

	d.SetId(sceneResult.Success.Id)

	return nil
}

func getLightsInScene(lightStates []interface{}) []string {
	var lightsInScene []string
	for _, lightState := range lightStates {
		lightState := lightState.(map[string]interface{})

		lightId := lightState["light_id"].(string)
		lightsInScene = append(lightsInScene, lightId)
	}
	return lightsInScene
}

func setLightStates(connection *common.Connection, sceneId string, lightStates []interface{}) error {
	for _, lightState := range lightStates {
		lightState := lightState.(map[string]interface{})
		var updateLightState scenes.LightState

		for key, bodyValue := range lightState {
			switch key {
			case "state":
				if bodyValue := bodyValue.(string); bodyValue != "" {
					if bodyValue == "on" {
						v := true
						updateLightState.On = &v
					} else {
						v := false
						updateLightState.On = &v
					}
				}
				break
			case "bri":
				if bodyValue := bodyValue.(string); bodyValue != "" {
					if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
						updateLightState.Bri = &bodyValue
					}
				}
				break
			case "hue":
				if bodyValue := bodyValue.(string); bodyValue != "" {
					if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
						updateLightState.Hue = &bodyValue
					}
				}
				break
			case "sat":
				if bodyValue := bodyValue.(string); bodyValue != "" {
					if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
						updateLightState.Sat = &bodyValue
					}
				}
				break
			case "ct":
				if bodyValue := bodyValue.(string); bodyValue != "" {
					if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
						updateLightState.CT = &bodyValue
					}
				}
				break
			case "xy":
				if bodyValue := bodyValue.(*schema.Set); bodyValue.Len() > 1 {
					updateLightState.XY = &[2]float64{ bodyValue.List()[0].(float64), bodyValue.List()[1].(float64) }
				}
				break
			case "transitiontime":
				if bodyValue := bodyValue.(string); bodyValue != "" {
					if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
						updateLightState.TransitionTime = &bodyValue
					}
				}
				break
			}
		}

		lightId := lightState["light_id"].(string)

		_, _, err := scenes.UpdateSceneLightState(connection, sceneId, lightId , &updateLightState)

		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSceneRead(d *schema.ResourceData, m interface{}) error {

	connection := m.(*common.Connection)

	scene, _, err := scenes.GetScene(connection, d.Id())

	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("name", scene.Name)

	d.Set("recycle", scene.Recycle)

	var lightStates []map[string]interface{}

	for lightId, lightState := range scene.Lightstates {
		state := make(map[string]interface{})
		state["light_id"] = lightId

		if lightState.On != nil {
			if true == *lightState.On {
				state["state"] = "on"
			} else {
				state["state"] = "off"
			}
		}

		if lightState.Bri != nil {
			state["bri"] = strconv.Itoa(*lightState.Bri)
		}

		if lightState.Hue != nil {
			state["hue"] = strconv.Itoa(*lightState.Hue)
		}

		if lightState.Sat != nil {
			state["sat"] = strconv.Itoa(*lightState.Sat)
		}

		if lightState.CT != nil {
			state["ct"] = strconv.Itoa(*lightState.CT)
		}

		if lightState.TransitionTime != nil {
			state["transitiontime"] = strconv.Itoa(*lightState.TransitionTime)
		}

		set := schema.Set{}
		if lightState.XY != nil {
			set.Add(lightState.XY[0])
			set.Add(lightState.XY[1])
			state["xy"] = &set
		}

		if lightState.Effect != nil {
			state["effect"] = *lightState.Effect
		}

		lightStates = append(lightStates, state)
	}

	d.Set("light_state", lightStates)

	return nil
}

func resourceSceneUpdate(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	d.Partial(true)

	lightStates := d.Get("light_state").(*schema.Set).List()
	lightsInScene := getLightsInScene(lightStates)

	// Now that we know the lights in the scene and have set their state, we can create the scene and capture it.

	scene := scenes.Update{
		Name:    d.Get("name").(string),
		Lights:  lightsInScene,
		Recycle: d.Get("recycle").(bool),
	}

	scenes.UpdateAPI(connection, d.Id(), &scene)

	d.SetPartial("name")

	setLightStates(connection, d.Id(), lightStates)

	d.Partial(false)

	return nil
}

func resourceSceneDelete(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	_, _, err := scenes.DeleteAPI(connection, d.Id())

	if err != nil {
		return err
	}

	return nil
}

func validateLightOnOffState(i interface{}, s string) (_ []string, errors []error) {
	value := i.(string)

	switch value {
	case "on":
	case "off":
		break
	default:
		errors = append(errors, fmt.Errorf("%q must be either 'on' or 'off'", s))
	}

	return
}