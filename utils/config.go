package utils

import "github.com/spf13/viper"

type Config struct {
	//DB
	User                 string `mapstructure:"USER"`
	Passwd               string `mapstructure:"PASSWORD"`
	Net                  string `mapstructure:"NETWORK"`
	Addr                 string `mapstructure:"ADDR"`
	DBName               string `mapstructure:"DBNAME"`
	AllowNativePasswords bool   `mapstructure:"ALLOWNATIVEPASSWORDS"`
	//Logger
	LogPath string `mapstructure:"LOGPATH"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("auth-server")
	viper.SetConfigType("env")

	//viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
