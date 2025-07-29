package structs

type ServerConfig struct {
	MySQL  MySQLConfig  `json:"mysql"`
	Consul ConsulConfig `json:"consul"`
	Name   string       `json:"name"`
}

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ConsulConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	HealthCheck string `json:"health_check_ip"`
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
