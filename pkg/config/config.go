package config

type Server struct {
	Port    string `yaml:"Port" env:"SERVER_PORT" envDefault:"8000"`
	TimeOut int64  `yaml:"TimeOut" env:"SERVER_TIMEOUT" envDefault:"30"`
}

type DB struct {
	Driver   string `yaml:"driver" env:"DB_DRIVER" envDefault:"postgres"`
	Name     string `yaml:"name" env:"DB_NAME" envDefault:"ushort"`
	Host     string `yaml:"host" env:"DB_HOST" envDefault:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" envDefault:"5432"`
	User     string `yaml:"user" env:"DB_USER" envDefault:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD,unset"`
	SslMode  string `yaml:"sslMode" env:"DB_SSLMODE" envDefault:"disable"`
	TimeZone string `yaml:"timeZone" env:"DB_TIMEZONE" envDefault:"Europe/Moscow"`
}

type Config struct {
	Server Server `yaml:"Server"`
	DB     DB     `yaml:"DB"`
}
