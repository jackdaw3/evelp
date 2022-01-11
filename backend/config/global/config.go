package global

var (
	Conf *Config
)

type AppConfig struct {
	ServerPort string `yaml:"serverPort"`
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

type DataConfig struct {
	RefreshStaticData        bool `yaml:"refreshStaticData"`
	RefreshRemoteData        bool `yaml:"refreshRemoteData"`
	RefreshUnusualRemoteData bool `yaml:"refreshUnusualRemoteData"`

	StaticDataPath    string `yaml:"staticDataPath"`
	RemoteDataAddress string `yaml:"remoteDataAddress"`
	RemoteDataSource  string `yaml:"remoteDataSource"`
}

type Config struct {
	App   AppConfig   `yaml:"app"`
	MySQL MySQLConfig `yaml:"mysql"`
	Data  DataConfig  `yaml:"data"`
}
