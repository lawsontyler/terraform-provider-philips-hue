package hue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lawsontyler/ghue/sdk/common"
	"github.com/lawsontyler/ghue/sdk/sensors"
)

func dataSourceHueSensor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHueSensorRead,

		Schema: map[string]*schema.Schema {
			"name": {
				Type: schema.TypeString,
				ConflictsWith: []string{"sensor_id"},
				Optional: true,
			},
			"sensor_id": {
				Type: schema.TypeString,
				ConflictsWith: []string{"name"},
				Optional: true,
			},
		},
	}
}

func dataSourceHueSensorRead(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*common.Connection)

	sensorName := d.Get("name")
	sensorId := d.Get("sensor_id")

	if sensorId != nil {
		sensorId := sensorId.(string)
		_, _, err := sensors.GetSensor(connection, sensorId)

		if err != nil {
			return err
		}

		d.SetId(sensorId)
	} else if sensorName != nil {
		sensorName := sensorName.(string)
		sensorId, _, err := sensors.GetSensorIdByName(connection, sensorName)

		if err != nil {
			return err
		}

		d.SetId(sensorId)
	}

	return nil
}