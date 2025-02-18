package config

type Database struct {
	Type          string `json:"type" mapstructure:"type" yaml:"type"`
	SqlitePath    string `json:"sqlite_path" mapstructure:"sqlite_path" yaml:"sqlite_path"`
	MysqlHost     string `json:"mysql_host" mapstructure:"mysql_host" yaml:"mysql_host"`
	MysqlPort     string `json:"mysql_port" mapstructure:"mysql_port" yaml:"mysql_port"`
	MysqlDbname   string `json:"mysql_dbname" mapstructure:"mysql_dbname" yaml:"mysql_dbname"`
	MysqlUsername string `json:"mysql_username" mapstructure:"mysql_username" yaml:"mysql_username"`
	MysqlPassword string `json:"mysql_password" mapstructure:"mysql_password" yaml:"mysql_password"`
}

type Telegram struct {
	Token string `json:"token" mapstructure:"token" yaml:"token"`
	Id    int64  `json:"id" mapstructure:"id" yaml:"id"`
}

type Weixin struct {
	XcxAppid  string `json:"xcx_appid" mapstructure:"xcx_appid" yaml:"xcx_appid"`
	XcxSecret string `json:"xcx_secret" mapstructure:"xcx_secret" yaml:"xcx_secret"`
	GzhAppid  string `json:"gzh_appid" mapstructure:"gzh_appid" yaml:"gzh_appid"`
	GzhSecret string `json:"gzh_secret" mapstructure:"gzh_secret" yaml:"gzh_secret"`
}
