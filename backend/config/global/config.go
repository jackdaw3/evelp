package global

import "time"

var (
	Conf *Config
)

type AppConfig struct {
	ServerPort string `yaml:"serverPort"`
	LogLevel   string `yaml:"logLevel"`
}

type CryptoConfig struct {
	KeyPath string `yaml:"keyPath"`
}

type LocalData struct {
	Path    string `yaml:"path"`
	Refresh bool   `yaml:"refresh"`
}

type RemoteData struct {
	Address    string `yaml:"address"`
	DataSource string `yaml:"datasource"`
	Refresh    bool   `yaml:"refresh"`
}

type DataConfig struct {
	Local  LocalData  `yaml:"local"`
	Remote RemoteData `yaml:"remote"`
}

type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	Database        string `yaml:"database"`
	UserName        string `yaml:"username"`
	Password        string `yaml:"password"`
	Loc             string `yaml:"loc"`
	Charset         string `yaml:"charset"`
	MaxIdleConn     int    `yaml:"maxIdleConn"`
	MaxOpenConn     int    `yaml:"maxOpenConn"`
	ConnMaxLifeTime int    `yaml:"connMaxLifetime"`
	AutoMigrate     bool   `yaml:"autoMigrate"`
}

type RedisConfig struct {
	Address           string        `yaml:"address"`
	Password          string        `yaml:"password"`
	Database          int           `yaml:"database"`
	HistoryExpireTime time.Duration `yaml:"history_expire_time"`
	OrderExpireTime   time.Duration `yaml:"order_expire_time"`
	Refresh           bool          `yaml:"refresh"`
}

type Config struct {
	App    AppConfig    `yaml:"app"`
	Crypto CryptoConfig `yaml:"crypto"`
	Data   DataConfig   `yaml:"data"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
}
