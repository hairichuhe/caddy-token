package caddyutil

type Config struct {
	AvatarSrc   string
	FileSrc     string
	AllowOrigin string
}

// DefaultConfig for the loginsrv handler
func DefaultConfig() *Config {
	return &Config{
		AllowOrigin: "http://127.0.0.1:8080",
	}
}
