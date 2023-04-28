package conf

import (
	"LibManSys/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("conf") // set the config file name. Viper will automatically detect the file extension name
	viper.AddConfigPath("./")   // search the config file under the current directory

	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic(err)
	}

	logrus.Info("Configuration file loaded")

	var confItems = map[string][]string{
		"mysql":  {"user", "password", "host", "port", "db_name"},
		"server": {"port"},
	}

	for k, v := range confItems {
		checkConfIsSet(k, v)
	}

	logrus.Info("All required values in configuration file are set")
}

func checkConfIsSet(name string, keys []string) {
	for i := range keys {
		wholeKey := name + "." + keys[i]
		if !viper.IsSet(wholeKey) {
			logrus.WithField(wholeKey, nil).
				Fatal("The following item of your configuration file hasn't been set properly: ")
		}
	}
}

func GetMysqlLoginConfig() *model.ConnectConfig {
	return &model.ConnectConfig{
		User:     viper.GetString("mysql.user"),
		Password: viper.GetString("mysql.password"),
		Host:     viper.GetString("mysql.host"),
		Port:     viper.GetInt("mysql.port"),
		DBName:   viper.GetString("mysql.db_name"),
	}
}
