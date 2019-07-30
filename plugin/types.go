package plugin

type Matcher struct {
	Path string
}

type Plugin struct {
	Matchers *Matchers `description: "Matchers plugin"`
}

type Matchers map[string]*Matcher
