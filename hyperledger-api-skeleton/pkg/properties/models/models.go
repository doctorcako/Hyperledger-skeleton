package models

import "fmt"

type HttpHandler struct {
	Type        string `yaml:"type" default:"restful"`
	Url         string `yaml:"url" default:"127.0.0.1"`
	Port        string `yaml:"port" default:"3000"`
	NetworkType string `yaml:"network_type" default:"tcp"`
	Swagger     string `yaml:"swagger" default:"../docs/swagger.yaml"`
	LogLevel    string `yaml:"log_level" default:"INFO"`
}

type Micro struct {
	Name string `yaml:"name" default:"micro-template"`
}

type GatewayParams struct {
	MspID        string `yaml:"msp_id"`
	PeerEndpoint string `yaml:"peer_endpoint"`
	GatewayPeer  string `yaml:"gateway_peer"`
	ChannelName  string `yaml:"channel_name"`
}

type Kafka struct {
	Brokers      []string     `yaml:"brokers"`
	InputTopics  InputTopics  `yaml:"input_topics"`
	OutputTopics OutputTopics `yaml:"output_topics"`
	Retries      Retries      `yaml:"retries"`
	Consumer     Consumer     `yaml:"consumer"`
}

type Retries struct {
	MaxRetries int `yaml:"max_retries"`
	Delay      int `yaml:"delay"`
}

type InputTopics struct {
	Prediction string `yaml:"prediction"`
}

type OutputTopics struct {
	DltTopic string `yaml:"dlt_topic"`
}

type Consumer struct {
	NumWorkers int `yaml:"num_workers"`
}

type Db struct {
	Postgres Postgres `yaml:"postgres"`
	Redis    Redis    `yaml:"redis"`
}

type Postgres struct {
	Host         string `yaml:"host" default:"localhost"`
	Port         string `yaml:"port" default:"5432"`
	User         string `yaml:"user" default:"postgres"`
	Password     string `yaml:"password" default:"12345678"`
	Database     string `yaml:"database" default:"postgres"`
	Schema       string `yaml:"schema" default:""`
	MaxOpenConns int    `yaml:"max_open_conns" default:"8"`
	MaxIdleConns int    `yaml:"max_idle_conns" default:"4"`
}

func (p Postgres) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", p.Host, p.Port, p.User, p.Password, p.Database)
}

type Redis struct {
	Host         string `yaml:"host" default:"localhost"`
	Port         string `yaml:"port" default:"5432"`
	User         string `yaml:"user" default:""`
	Password     string `yaml:"password" default:""`
	Database     int    `yaml:"database" default:"0"`
	MaxOpenConns int    `yaml:"max_open_conns" default:"8"`
	MaxIdleConns int    `yaml:"max_idle_conns" default:"4"`
}

func (p Redis) GetDSN() string {
	if p.User != "" && p.Password != "" {
		return fmt.Sprintf("redis://%s:%s@%s:%s/%d", p.User, p.Password, p.Host, p.Port, p.Database)
	} else {
		return fmt.Sprintf("redis://%s:%s/%d", p.Host, p.Port, p.Database)
	}
}
