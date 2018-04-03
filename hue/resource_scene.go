package hue

import "github.com/hashicorp/terraform/helper/schema"


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
		},
	}
}

func resourceSceneCreate(d *schema.ResourceData, m interface{}) error {
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