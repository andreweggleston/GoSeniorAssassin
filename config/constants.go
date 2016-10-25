package config

import (
	"net/url"
	"os"
	"reflect"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

var (
	mdTableTemplate = template.Must(template.New("doc").Parse(`
| Environment Variable | Description |
|----------------------|-------------|{{range .}}
|    ` + "`{{index . 0}}`" + `     |{{index . 1}}|{{end}}
`))
)

type constants struct {
	ListenAddress	string	`envconfig:"SA_SERVER_ADDR" default:"0.0.0.0:8081" doc:"Address to serve on"`
	RabbitMQURL       string   `envconfig:"RABBITMQ_URL" default:"amqp://guest:guest@localhost:5672/" doc:"URL for AMQP server"`
}