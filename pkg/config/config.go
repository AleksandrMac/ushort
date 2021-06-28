package config

type Server struct {
	Port    string `yaml:"Port" env:"PORT" envDefault:"8000"`
	TimeOut int64  `yaml:"TimeOut" env:"SERVER_TIMEOUT" envDefault:"30"`
}

type DB struct {
	URL string `yaml:"db_url" env:"DATABASE_URL" envDefault:"postgres://postgres:password@localhost:5432/ushort?sslmode=disable"`

	// Driver   string `yaml:"driver" env:"DB_DRIVER" envDefault:"postgres"`
	// Name     string `yaml:"name" env:"DB_NAME" envDefault:"ushort"`
	// Host     string `yaml:"host" env:"DB_HOST" envDefault:"localhost"`
	// Port     string `yaml:"port" env:"DB_PORT" envDefault:"5432"`
	// User     string `yaml:"user" env:"DB_USER" envDefault:"postgres"`
	// Password string `yaml:"password" env:"DB_PASSWORD,unset"`
	// SslMode  string `yaml:"sslMode" env:"DB_SSLMODE" envDefault:"disable"`
	// TimeZone string `yaml:"timeZone" env:"DB_TIMEZONE" envDefault:"Europe/Moscow"`
}

type Config struct {
	Server Server `yaml:"Server"`
	DB     DB     `yaml:"DB"`
}

func New() *Config {
	return new(Config)
}
