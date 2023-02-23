package configs

type Config struct {
	Self struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		Debug     bool   `yaml:"debug"`
		AuthToken string `yaml:"auth_token"`
	} `yaml:"self"`
	HTTP struct {
		AfdUserID   string `yaml:"afd_user_id"`
		AfdAPIToken string `yaml:"afd_api_token"`
	} `yaml:"http"`
	Webhook struct {
		Enabled bool   `yaml:"enabled"`
		Point   string `yaml:"point"`
	} `yaml:"webhook"`
	Db struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DbName   string `yaml:"db_name"`
	} `yaml:"db"`
}
