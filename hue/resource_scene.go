package hue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lawsontyler/ghue/sdk/common"
	"github.com/lawsontyler/ghue/sdk/scenes"
	"github.com/lawsontyler/terraform-provider-philips-hue/hue/lib/constants"
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

						"on": {
							Type: schema.TypeBool,
							Required: true,
						},

						"bri": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"hue": {
							Type: schema.TypeInt,
							Optional: true,
							ConflictsWith: []string{"light_state.xy", "light_state.ct"},
						},
						"sat": {
							Type: schema.TypeInt,
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
							Type: schema.TypeInt,
							Optional: true,
							ConflictsWith: []string{"light_state.hue", "light_state.sat", "light_state.xy"},
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

		updateLightState.On = lightState["on"].(bool)

		if brightness := lightState["bri"]; brightness != nil {
			updateLightState.Bri = brightness.(int)
		}

		if hue := lightState["hue"]; hue != nil {
			updateLightState.Hue = hue.(int)
		}

		if saturation := lightState["sat"]; saturation != nil {
			updateLightState.Sat = saturation.(int)
		}

		if ct := lightState["ct"]; ct != nil {
			updateLightState.CT = ct.(int)
		}

		if xy := lightState["xy"].(*schema.Set); xy.Len() > 1 {
			updateLightState.XY = &[2]float64{ xy.List()[0].(float64), xy.List()[1].(float64) }
		}

		if transitiontime := lightState["transitiontime"]; transitiontime != nil {
			updateLightState.TransitionTime = transitiontime.(int)
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
		state["on"] = lightState.On
		state["bri"] = lightState.Bri
		state["transitiontime"] = lightState.TransitionTime
		state["hue"] = lightState.Hue
		state["sat"] = lightState.Sat
		state["ct"] = lightState.CT

		set := schema.Set{}
		if lightState.XY != nil {
			set.Add(lightState.XY[0])
			set.Add(lightState.XY[1])
		}

		state["xy"] = &set
		state["effect"] = lightState.Effect

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

func resourceSceneExists(d *schema.ResourceData, m interface{}) (bool, error) {
	connection := m.(*common.Connection)

	_, hueErr, err := scenes.GetScene(connection, d.Id())

	if err != nil {
		if hueErr.Error.Type == int(constants.NOT_FOUND) {
			return false, nil
		}

		return false, err
	}

	return true, err
}
