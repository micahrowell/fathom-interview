package internal

type Message struct {
	UserID string `json:"UserID"`
	Body   string `json:"Body"`
}

type Configuration struct {
	Server serverconfig
}

type serverconfig struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}
