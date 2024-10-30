package properties

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"hyperledger-api-skeleton/hyperledger-business-logic/pkg/errors"
	"hyperledger-api-skeleton/hyperledger-business-logic/pkg/properties/models"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customConfig"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customError"
)

type Properties struct {
	ApplicationProperties ApplicationProperties
}

type ApplicationProperties struct {
	PropertiesViper *viper.Viper
	Micro           models.Micro         `yaml:"micro"`
	HttpServer      models.HttpHandler   `yaml:"http_server"`
	GatewayParams   models.GatewayParams `yaml:"gateway_params"`
	Kafka           models.Kafka         `yaml:"kafka"`
}

func (p *Properties) InitAllConfig(applicationProperties string) customError.Error {
	return p.ApplicationProperties.initApplicationConfig(applicationProperties)
}

func (p *ApplicationProperties) initApplicationConfig(applicationProperties string) customError.Error {
	//Define Default value
	p.Micro = p.Micro.DefaultTag()
	p.HttpServer = p.HttpServer.DefaultTag()

	//Read Yaml Properties
	v, errC := customConfig.SetViperByEnvironment(applicationProperties)
	if errC != nil {
		return errC
	}
	mapProperties := v.AllSettings()
	yamlString, err := yaml.Marshal(mapProperties)
	if err != nil {
		return errors.NewCustomError(errors.DataManipulationError, err.Error())
	}
	err = yaml.Unmarshal(yamlString, &p)
	if err != nil {
		return errors.NewCustomError(errors.DataManipulationError, err.Error())
	}
	p.PropertiesViper = v

	return nil
}
