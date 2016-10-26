package config

import (
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
	PublicAddress   string   `envconfig:"PUBLIC_ADDR" doc:"Publicly accessible address for the server, requires schema"`
	RabbitMQURL     string  `envconfig:"RABBITMQ_URL" default:"amqp://guest:guest@localhost:5672/" doc:"URL for AMQP server"`

	DbAddr		string	`envconfig:"DATABASE_ADDR" default:"127.0.0.1:5432" doc:"Database Address"`
	DbDatabase	string 	`envconfig:"DATABASE_NAME" default:"seniorassassin" doc:"Database Name"`
	DbUsername	string 	`envconfig:"DATABASE_USERNAME" default:"seniorassassin" doc:"Database Username"`
	DbPassword 	string 	`envconfig:"DATABASE_PASSWORD" default:"assassinpass" doc:"Database password"`

	CookieStoreSecret string   `envconfig:"COOKIE_STORE_SECRET" default:"secret" doc:"base64 encoded key to use for encrypting cookies"`
}

var Constants = constants{}

func init() {
	err:= envconfig.Process("HELEN", &Constants)
	if err != nil {
		logrus.Fatal(err)
	}
}

func PrintConfigDoc() {
	var data [][]string
	t := reflect.TypeOf(constants{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("envconfig") == "" {
			continue
		}
		data = append(data, []string{field.Tag.Get("envconfig"), field.Tag.Get("doc")})
	}

	mdTableTemplate.Execute(os.Stdout, data)
}