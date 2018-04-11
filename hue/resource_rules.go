package hue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lawsontyler/ghue/sdk/common"
	"fmt"
	"github.com/lawsontyler/ghue/sdk/rules"
	"github.com/lawsontyler/terraform-provider-philips-hue/hue/lib/constants"
	"github.com/Sirupsen/logrus"
	"strconv"
)


func resourceRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleCreate,
		Read:   resourceRuleRead,
		Update: resourceRuleUpdate,
		Delete: resourceRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Required: true,
			},
			"condition": {
				Type: schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type: schema.TypeString,
							Required: true,
						},
						"operator": {
							Type: schema.TypeString,
							Required: true,
							ValidateFunc: func(i interface{}, s string) (_ []string, errors []error) {
								value := i.(string)
								switch value {
								case "eq":
								case "gt":
								case "lt":
								case "dx":
								case "ddx":
								case "stable":
								case "not stable":
								case "in":
								case "not in":
									break
								default:
									errors = append(errors, fmt.Errorf("%s is not a valid option for %q", value, s))
								}

								return
							},
						},
						"value": {
							Type: schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"action": {
				Type: schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type: schema.TypeString,
							Required: true,
						},
						"method": {
							Type: schema.TypeString,
							Required: true,
							ValidateFunc: func(i interface{}, s string) (_ []string, errors []error) {
								value := i.(string)
								switch value {
								case "PUT":
								case "POST":
								case "GET":
									break
								default:
									errors = append(errors, fmt.Errorf("%s is not a valid option for %q", value, s))
								}

								return
							},
						},
						"body": {
							Type: schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validateLightOnOffState,
									},

									"bri": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"hue": {
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"action.body.xy", "action.body.ct"},
									},
									"sat": {
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"action.body.xy", "action.body.ct"},
									},
									"xy": {
										Type:          schema.TypeSet,
										Optional:      true,
										ConflictsWith: []string{"action.body.hue", "action.body.sat", "action.body.ct"},
										Elem:          &schema.Schema{Type: schema.TypeFloat},
										MaxItems:      2,
									},
									"ct": {
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"action.body.hue", "action.body.sat", "action.body.xy"},
									},
									"alert": {
										Type: schema.TypeString,
										Optional: true,
									},
									"effect": {
										Type: schema.TypeString,
										Optional: true,
									},
									"transitiontime": {
										Type: schema.TypeString,
										Optional: true,
										Default: 4,
									},
									"bri_inc": {
										Type: schema.TypeString,
										Optional: true,
									},
									"hue_inc": {
										Type: schema.TypeString,
										Optional: true,
									},
									"sat_inc": {
										Type: schema.TypeString,
										Optional: true,
									},
									"ct_inc": {
										Type: schema.TypeString,
										Optional: true,
									},
									"xy_inc": {
										Type: schema.TypeString,
										Optional: true,
									},
									"scene": {
										Type: schema.TypeString,
										Optional: true,
									},

								},
							},
						},
					},
				},
			},
		},
	}
}

func dataToConditionArray(conditions *schema.Set) []rules.Condition {
	var conditionArray []rules.Condition

	if v := conditions; v.Len() > 0 {
		for _, v := range v.List() {
			v := v.(map[string]interface{})

			condition := rules.Condition{
				Address: v["address"].(string),
				Operator: v["operator"].(string),
			}

			if value := v["value"].(string); value != "" {
				condition.Value = &value
			}

			conditionArray = append(conditionArray, condition)
		}
	}

	return conditionArray
}

func dataToActionArray(actions *schema.Set) []rules.Action {
	var actionArray []rules.Action

	if v := actions; v.Len() > 0 {
		for _, v := range v.List() {
			v := v.(map[string]interface{})

			action := rules.Action{}
			actionBody := rules.ActionBody{}

			action.Address = v["address"].(string)
			action.Method = v["method"].(string)

			logrus.Errorf("Body is: %s", v["body"].(*schema.Set).List()[0].(map[string]interface{}))

			for key, bodyValue := range v["body"].(*schema.Set).List()[0].(map[string]interface{}) {
				logrus.Errorf("Checking %s...", key)
				switch key {
				case "state":
					state := bodyValue.(string)
					if state != "" {
						if state == "on" {
							v := true
							actionBody.On = &v
						} else {
							v := false
							actionBody.On = &v
						}
					}

					break
				case "bri":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue > 0 {
							actionBody.Bri = &bodyValue
						}
					}
					break
				case "hue":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
							actionBody.Hue = &bodyValue
						}
					}
					break
				case "sat":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
							actionBody.Sat = &bodyValue
						}
					}
					break
				case "xy":
					bodyValue := bodyValue.(*schema.Set)

					if bodyValue.Len() > 1 {
						actionBody.XY = &[2]float64{bodyValue.List()[0].(float64),bodyValue.List()[1].(float64)}
					}
					break
				case "ct":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
							actionBody.CT = &bodyValue
						}
					}
					break
				case "alert":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						actionBody.Alert = &bodyValue
					}
					break
				case "effect":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						actionBody.Effect = &bodyValue
					}
					break
				case "bri_inc":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
							actionBody.BriInc = &bodyValue
						}
					}
					break
				case "hue_inc":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
							actionBody.HueInc = &bodyValue
						}
					}
					break
				case "sat_inc":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
							actionBody.SatInc = &bodyValue
						}
					}
					break
				case "ct_inc":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						if bodyValue, _ := strconv.Atoi(bodyValue); bodyValue >= 0 {
							actionBody.CTInc = &bodyValue
						}
					}
					break
				case "xy_inc":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						if bodyValue, _ := strconv.ParseFloat(bodyValue, 64); bodyValue > -1 {
							actionBody.XYInc = &bodyValue
						}
					}
					break
				case "scene":
					if bodyValue := bodyValue.(string); bodyValue != "" {
						actionBody.Scene = &bodyValue
					}
					break
				default:
					continue
				}

			}

			action.Body = actionBody

			actionArray = append(actionArray, action)
		}
	}

	logrus.Errorf("Action Array is: %s", actionArray)

	return actionArray
}

func resourceRuleCreate(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	conditions := dataToConditionArray(d.Get("condition").(*schema.Set))
	actions := dataToActionArray(d.Get("action").(*schema.Set))

	rule := rules.Create{
		Name: d.Get("name").(string),
		Conditions: conditions,
		Actions: actions,
	}

	result, _, err := rules.CreateAPI(connection, &rule)

	if err != nil {
		return err
	}

	d.SetId(result.Success.Id)

	return nil
}

func resourceRuleRead(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	rule, hueErr, err := rules.GetRule(connection, d.Id())

	if err != nil && hueErr != nil && hueErr.Error.Type == int(constants.NOT_FOUND) {
		d.SetId("")
		return nil
	}

	d.Set("name", rule.Name)

	conditions := make([]map[string]interface{}, 0, len(rule.Conditions))
	actions    := make([]map[string]interface{}, 0, len(rule.Actions))

	for _, ruleCondition := range rule.Conditions {
		condition := map[string]interface{}{
			"address": ruleCondition.Address,
			"operator": ruleCondition.Operator,
		}

		if ruleCondition.Value != nil {
			condition["value"] = ruleCondition.Value
		}

		conditions = append(conditions, condition)
	}
	// logrus.Errorf("Rule Actions: %s", rule.Actions)


	for _, ruleAction := range rule.Actions {

		logrus.Errorf("Body: %s", ruleAction.Body)

		body := map[string]interface{} {}

		if ruleAction.Body.On != nil {
			if *ruleAction.Body.On == true {
				body["state"] = "on"
			} else {
				body["state"] = "off"
			}
		}

		if ruleAction.Body.Bri != nil {
			body["bri"] = strconv.Itoa(*ruleAction.Body.Bri)
		}
		if ruleAction.Body.Hue != nil {
			body["hue"] = strconv.Itoa(*ruleAction.Body.Hue)
		}
		if ruleAction.Body.Sat != nil {
			body["sat"] = strconv.Itoa(*ruleAction.Body.Sat)
		}
		if ruleAction.Body.CT != nil {
			body["ct"] = strconv.Itoa(*ruleAction.Body.CT)
		}
		if ruleAction.Body.XY != nil {
			body["xy"] = *ruleAction.Body.XY
		}
		if ruleAction.Body.Alert != nil {
			body["alert"] = ruleAction.Body.Alert
		}
		if ruleAction.Body.Effect != nil {
			body["effect"] = ruleAction.Body.Effect
		}
		if ruleAction.Body.BriInc != nil {
			body["bri_inc"] = strconv.Itoa(*ruleAction.Body.BriInc)
		}
		if ruleAction.Body.HueInc != nil {
			body["hue_inc"] = strconv.Itoa(*ruleAction.Body.HueInc)
		}
		if ruleAction.Body.SatInc != nil {
			body["sat_inc"] = strconv.Itoa(*ruleAction.Body.SatInc)
		}
		if ruleAction.Body.CTInc != nil {
			body["ct_inc"] = strconv.Itoa(*ruleAction.Body.CTInc)
		}
		if ruleAction.Body.XYInc != nil {
			body["xy_inc"] = strconv.FormatFloat(*ruleAction.Body.XYInc, 'f', 3, 64)
		}
		if ruleAction.Body.Scene != nil {
			body["scene"] = *ruleAction.Body.Scene
		}

		bodySet := make([]interface{}, 0, 1)

		bodySet = append(bodySet, body)

		action := map[string]interface{}{
			"address": ruleAction.Address,
			"method": ruleAction.Method,
			"body": bodySet,
		}

		actions = append(actions, action)
	}

	d.Set("condition", conditions)
	d.Set("action", actions)

	return nil
}

func resourceRuleUpdate(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	conditions := dataToConditionArray(d.Get("condition").(*schema.Set))
	actions := dataToActionArray(d.Get("action").(*schema.Set))

	rule := rules.Update {
		Name: d.Get("name").(string),
		Conditions: conditions,
		Actions: actions,
	}

	_, _, err := rules.UpdateAPI(connection, d.Id(), &rule)

	if err != nil {
		return err
	}

	return nil
}

func resourceRuleDelete(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	_, _, err := rules.DeleteAPI(connection, d.Id())

	if err != nil {
		return err
	}

	return nil
}