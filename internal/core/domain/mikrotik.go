package domain

type MikrotikInterface string

const (
	Ether1  MikrotikInterface = "ether1"
	Ether2  MikrotikInterface = "ether2"
	Ether3  MikrotikInterface = "ether3"
	Ether4  MikrotikInterface = "ether4"
	Ether5  MikrotikInterface = "ether5"
	Ether6  MikrotikInterface = "ether6"
	Ether7  MikrotikInterface = "ether7"
	Ether8  MikrotikInterface = "ether8"
	Ether9  MikrotikInterface = "ether9"
	Ether10 MikrotikInterface = "ether10"
	Ether11 MikrotikInterface = "ether11"
	Ether12 MikrotikInterface = "ether12"
	Ether13 MikrotikInterface = "ether13"
	SFP1    MikrotikInterface = "sfp-sfpplus1"
	SFP2    MikrotikInterface = "sfp-sfpplus2"
	SFP3    MikrotikInterface = "sfp-sfpplus3"
	SFP4    MikrotikInterface = "sfp-sfpplus4"
	SFP5    MikrotikInterface = "sfp-sfpplus5"
	SFP6    MikrotikInterface = "sfp-sfpplus6"
	SFP7    MikrotikInterface = "sfp-sfpplus7"
	SFP8    MikrotikInterface = "sfp-sfpplus8"
	SFP9    MikrotikInterface = "sfp-sfpplus9"
	SFP10   MikrotikInterface = "sfp-sfpplus10"
	SFP11   MikrotikInterface = "sfp-sfpplus11"
	SFP12   MikrotikInterface = "sfp-sfpplus12"
)

type Resource struct {
	Source *string
	Cpu    string
	Uptime string
}

type Traffic struct {
	Source *string
	Rx     string
	Tx     string
}

type Routers struct {
	Name    string  `yaml:"name"`
	Address string  `yaml:"address"`
	Routes  []Route `yaml:"routes"`
}
type Route struct {
	Event       string `yaml:"event"`
	Destination string `yaml:"destination"`
	Gateway     string `yaml:"gateway"`
	Distance    int    `yaml:"distance"`
	RoutingMark string `yaml:"routing_mark"`
	Comment     string `yaml:"comment"`
	Disabled    bool   `yaml:"disabled"`
}

type RouterRoutes struct {
	Routers []Routers `yaml:"routers"`
}
