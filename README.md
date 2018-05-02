# Philips Hue Terraform Module

## About 

This is a project of passion, mixing a couple of my favourite things - Philips Hue and Infrastructure as Code.

The underlying API calls are done via lawsontyler/ghue.  Not _everything_ is supported yet, but it's just about enough
for my own personal use :D

## Installation

Since there are no releases as of yet, you should be able to:

```
# Set your GOPATH if not already set:
# export GOPATH="~/go"

go get github.com/lawsontyler/terraform-provider-philips-hue
cd $GOPATH/github.com/lawsontyler/terraform-provider-philips-hue
go build -o terraform-provider-philips-hue

# You'll need to know your OS and Archetecture.
# e.g. darwin_amd64; linux_i386
mkdir -p <path-to-my-terraform>/terraform.d/plugins/<os_arch>
ln -s $GOPATH/github.com/lawsontyler/terraform-provider-philips-hue \
      <path-to-my-terraform>/terraform.d/plugins/<os_arch>/terraform-provider-philips-hue

cd <path-to-my-terraform>
terraform init
# Now use terraform like normal.  You might need to run `terraform init` a bunch more times.  A lot.
# Run it all the time.
```

## Contributing

Pull requests are welcome!  I plan on only developing this as far as I need to for myself.  Please, extend it as you see
fit and make some PRs back.

## Usage

Here's what some terrform might look like.  This is real code I'm using for my house.

```
provider "philips-hue" {
    hub_address = "192.168.1.170"
    hub_username = "totally-my-username"
}

data "philips-hue_light" "basement-1" {
    light_id = "2"
}

data "philips-hue_light" "basement-2" {
    light_id = "6"
}

data "philips-hue_sensor" "basement-dimmer" {
    sensor_id = "8"
}

resource "philips-hue_scene" "basement-red" {
    name = "Basement Red"

    light_state {
        light_id = "${data.philips-hue_light.basement-1.id}"
        bri = "1"
        sat = "255"
        hue = "0"
    }

    light_state {
        light_id = "${data.philips-hue_light.basement-2.id}"
        bri = "1"
        sat = "255"
        hue = "0"
    }
}

resource "philips-hue_group" "basement-group" {
    name = "Basement Lamps"
    lights = [ "${data.philips-hue_light.basement-1.id}", "${data.philips-hue_light.basement-2.id}" ]
}

resource "philips-hue_rule" "basement-dimmer-on-short" {
    name = "Basement Dimmer On Short"

    condition {
        address = "/sensors/${data.philips-hue_sensor.basement-dimmer.id}/state/buttonevent"
        operator = "eq"
        value = "1000"
    }

    condition {
        address = "/sensors/${data.philips-hue_sensor.basement-dimmer.id}/state/lastupdated"
        operator = "dx"
    }

    action {
        address = "/groups/${philips-hue_group.basement-group.id}/action"
        method = "PUT"
        body {
            scene = "${philips-hue_scene.basement-red.id}"
        }
    }
}
```

