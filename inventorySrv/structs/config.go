package structs

type ServerConfig struct {
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Consul ConsulConfig `mapstructure:"consul"`
	Name   string       `mapstructure:"name"`
	Redis  RedisCnf     `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type ConsulConfig struct {
	Host string   `mapstructure:"host"`
	Port int      `mapstructure:"port"`
	Tags []string `mapstructure:"tags"`
}

type NacosCnf struct {
	Host      string `json:"host"`
	Port      uint64 `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Namespace string `json:"namespace"`
	DataID    string `json:"data_id"`
	Group     string `json:"group"`
}

type RedisCnf struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
}
