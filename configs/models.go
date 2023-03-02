package configs

type Config struct {
	Self struct {
		Host  string `yaml:"host"`
		Port  int    `yaml:"port"`
		Debug bool   `yaml:"debug"`
	} `yaml:"self"`
	Afdian struct {
		UserID   string `yaml:"user_id"`
		APIToken string `yaml:"api_token"`
	} `yaml:"afdian"`
	API struct {
		APIToken          string `yaml:"api_token"`
		OrderEndpoint     string `yaml:"order_endpoint"`
		SponsorsEndpoint  string `yaml:"sponsors_endpoint"`
		GetuseridEndpoint string `yaml:"getuserid_endpoint"`
		IsLimitHost       bool   `yaml:"is_limit_host"`
		Only              string `yaml:"only"`
	} `yaml:"api"`
	Webhook struct {
		Enabled  bool   `yaml:"enabled"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"webhook"`
	Db struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DbName   string `yaml:"db_name"`
	} `yaml:"db"`
}
