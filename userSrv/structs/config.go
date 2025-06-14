package structs

type ServerConfig struct {
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Consul ConsulConfig `mapstructure:"consul"`
	Name   string       `mapstructure:"name"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
