package config

import (
	"os"
	"reflect"
	"text/template"

	"github.com/sirupsen/logrus"
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
	//OpenIDRealm       string   `envconfig:"SERVER_OPENID_REALM" default:"0.0.0.0:8081" doc:"The OpenID Realm (See: [Section 9.2 of the OpenID Spec](https://openid.net/specs/openid-authentication-2_0-12.html#realms))"`
	AllowedOrigins    []string `envconfig:"ALLOWED_ORIGINS" default:"*"`

	AdminStudentIDs	string	`envconfig:"ADMIN_STUDENT_IDS" default:"Andrew_Eggleston" doc:"Separated by commas and no spaces e.x. \"Andrew_Eggleston,Nick_Fay\""`

	LoginRedirectPath string   `envconfig:"SERVER_REDIRECT_PATH" default:"http://10.0.0.5.nip.io:8080/" doc:"URL to redirect user to after a successful login"`
	//TODO: not :8081, make :8080 after frontend works

	RabbitMQURL     string  `envconfig:"RABBITMQ_URL" default:"amqp://guest:guest@localhost:5672/" doc:"URL for AMQP server"`
	RabbitMQQueue     string   `envconfig:"RABBITMQ_QUEUE" default:"events" doc:"Name of queue over which events are sent"`


	DbAddr		string	`envconfig:"DATABASE_ADDR" default:"seniorassassin.caeprqfnq40y.us-east-2.rds.amazonaws.com:5432" doc:"Database Address"`
	DbDatabase	string 	`envconfig:"DATABASE_NAME" default:"seniorassassin" doc:"Database Name"`
	DbUsername	string 	`envconfig:"DATABASE_USERNAME" default:"SeniorAssassin" doc:"Database Username"`
	DbPassword 	string 	`envconfig:"DATABASE_PASSWORD" default:"assassinpass" doc:"Database password"`

	CookieStoreSecret string   `envconfig:"COOKIE_STORE_SECRET" default:"secret" doc:"base64 encoded key to use for encrypting cookies"`
	CookieDomain      string   `envconfig:"SERVER_COOKIE_DOMAIN" default:"" doc:"Cookie URL domain"`
	SecureCookies      bool     `envconfig:"SECURE_COOKIE" doc:"Enable 'secure' flag on cookies" default:"false"`

	SlackbotURL       string   `envconfig:"SLACK_URL" doc:"Slack webhook URL"`

	IDWhitelist  	string   `envconfig:"ID_WHITELIST" doc:"ID Group XML page to use to filter logins"`
	MockupAuth        bool     `envconfig:"MOCKUP_AUTH" default:"false" doc:"Enable Mockup Authentication"`
	ServeStatic       bool     `envconfig:"SERVE_STATIC" default:"true" doc:"Serve /static/"`

}

var Constants = constants{}

func init() {
	err:= envconfig.Process("PARM", &Constants)
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