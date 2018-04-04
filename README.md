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

# SEIZURE WARNING

In order to set and validate scenes, I need to actually set light states.  This means that the lights will change and flick from
scene to scene as it reads and writes scenes.  As with any fast-changing lights, some people may feel discomfort when
running `terraform plan` or `terraform apply`.

I'm not a doctor - if you or anyone else in your household are prone to seizures, consult your doctor and make sure you
follow their advice for rapidly changing or flickering lights.  **By using this plugin, you take on this risk and absolve
the author(s) of any accountability.**

