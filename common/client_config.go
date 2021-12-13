package common

// Client config
type Client struct {
	Name     string `yaml:"name"` // unique name of client
	Url      string `yaml:"url"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}
