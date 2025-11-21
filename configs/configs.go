package configs

type Config struct {
	System SystemConfig `mapstructure:"system"`
	Mysql  MysqlConfig  `mapstructure:"mysql"`
	Log    LogConfig    `mapstructure:"log"`
	Jwt    JWTConfig    `mapstructure:"jwt"`
}

type SystemConfig struct {
	Port          string `mapstructure:"port"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime uint   `mapstructure:"expiretime"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Config   string `mapstructure:"config"`
	DBName   string `mapstructure:"db-name"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	OutputFile string `mapstructure:"output-file"`
}

// root:psw@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
func (m *MysqlConfig) Dsn() string {
	return m.UserName + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DBName + "?" + m.Config
}
