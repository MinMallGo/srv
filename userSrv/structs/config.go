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

type NacosCnf struct {
	Host      string `json:"host"`
	Port      uint64 `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Namespace string `json:"namespace"`
	DataID    string `json:"data_id"`
	Group     string `json:"group"`
}
