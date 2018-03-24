package caddyutil

type Config struct {
	MysqlServer   string
	Username string
	Password string
	Dataname string
	Rdserver string
	RdPW     string
	AvatarSrc string
	FileSrc string
	RdWord string
	UpFileSrc string
}

// DefaultConfig for the loginsrv handler
func DefaultConfig() *Config {
	return &Config{
		MysqlServer:"192.168.0.170:3306"
		Username:"root"
		Password:"root"
		Dataname:"college"
		Rdserver:"192.168.0.114:6379"
		RdPW:"root"
	}
}