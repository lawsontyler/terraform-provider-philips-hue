package hue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lawsontyler/ghue/sdk/groups"
	"github.com/lawsontyler/ghue/sdk/common"
	"fmt"
	"github.com/lawsontyler/ghue/sdk/rules"
	"github.com/lawsontyler/terraform-provider-philips-hue/hue/lib/constants"
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
							Type: schema.TypeMap,
							Required: true,
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
			conditionArray = append(conditionArray, rules.Condition{
				Address: v["address"].(string),
				Operator: v["operator"].(string),
				Value: v["value"].(string),
			})
		}
	}

	return conditionArray
}

func dataToActionArray(actions *schema.Set) []rules.Action {
	var actionArray []rules.Action

	if v := actions; v.Len() > 0 {
		for _, v := range v.List() {
			v := v.(map[string]interface{})
			actionArray = append(actionArray, rules.Action{
				Address: v["address"].(string),
				Method: v["method"].(string),
				Body: v["body"].(map[string]interface{}),
			})
		}
	}

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
	}

	d.Set("name", rule.Name)

	var conditions []map[string]interface{}

	for _, ruleCondition := range rule.Conditions {
		condition := make(map[string]interface{})

		condition["address"] = ruleCondition.Address
		condition["operator"] = ruleCondition.Operator
		if ruleCondition.Value != "" {
			condition["value"] = ruleCondition.Value
		}

		conditions = append(conditions, condition)
	}

	var actions []map[string]interface{}

	for _, ruleAction := range rule.Actions {
		action := make(map[string]interface{})

		action["address"] = ruleAction.Address
		action["method"] = ruleAction.Method
		action["body"] = ruleAction.Body

		actions = append(actions, action)
	}

	d.Set("condition", conditions)
	d.Set("action", actions)

	return nil
}

func resourceRuleUpdate(d *schema.ResourceData, m interface{}) error {
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

func resourceRuleDelete(d *schema.ResourceData, m interface{}) error {
	connection := m.(*common.Connection)

	_, _, err := groups.DeleteAPI(connection, d.Id())

	if err != nil {
		return err
	}

	return nil
}