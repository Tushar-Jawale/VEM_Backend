package config

type Config struct {
	Port string
	Env  string
}

func Load() *Config {
	return &Config{
		Port: ":8080",
		Env:  "development",
	}
}
