package customConfig

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Type       string `yaml:"type"`
	Url        string `yaml:"url"`
	Port       string `yaml:"port"`
	Swagger    string `yaml:"swagger"`
	LogLevel   string `yaml:"logLevel"`
	ServerName string `yaml:"serverName"`
	Container  string `yaml:"container"`
}

func ReadFile(serverName string, filePath string) (ServerConfiguration, customError.Error) {
	v, errC := SetViperByEnvironment(filePath)
	if errC != nil {
		return ServerConfiguration{}, errC
	}

	serverPath := serverName
	server := ServerConfiguration{
		Type:       v.GetString(serverPath + ".type"),
		Url:        v.GetString(serverPath + ".url"),
		Port:       v.GetString(serverPath + ".port"),
		Swagger:    v.GetString(serverPath + ".swagger"),
		LogLevel:   v.GetString(serverPath + ".logLevel"),
		ServerName: serverName,
		Container:  v.GetString(serverPath + ".container"),
	}
	return server, nil
}

func SetViperByEnvironment(fileName string) (*viper.Viper, customError.Error) {
	var effectiveFileName string
	v := viper.New()
	if os.Getenv("ENVIRONMENT") != "" {
		substr := strings.Split(fileName, ".yaml")
		effectiveFileName = substr[0] + "_" + strings.ToLower(os.Getenv("ENVIRONMENT")) + ".yaml"
	} else {
		effectiveFileName = fileName
	}

	errC := existsFile(effectiveFileName)
	if errC != nil {
		return nil, errC
	}

	v.SetConfigType("yaml")
	v.SetConfigFile(effectiveFileName)
	v.AutomaticEnv()
	envReplacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(envReplacer)
	err := v.ReadInConfig()
	if err != nil {
		return nil, customError.NewError(customError.FileOperationError, err.Error())
	}

	return v, nil
}

// This function check if the file exists
func existsFile(path string) customError.Error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	return customError.NewError(customError.FileOperationError, err.Error())
}
