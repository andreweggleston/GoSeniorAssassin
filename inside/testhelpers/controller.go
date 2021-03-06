package testhelpers

import (
	"net/http/httptest"
	"net"
	"testing"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
	"net/http/cookiejar"
	"net/url"
	"errors"
	"strconv"
	"net/http"
	"github.com/gorilla/websocket"
	"math/rand"
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"github.com/andreweggleston/GoSeniorAssassin/routes"
)

const InitMessages int = 4

type SuffixList struct{}

var (
	options = &cookiejar.Options{PublicSuffixList: SuffixList{}}
)

func (SuffixList) PublicSuffix(_ string) string {
	return ""
}

func (SuffixList) String() string {
	return ""
}

var DefaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	}).Dial,
}

func NewClient() (client *http.Client) {
	client = new(http.Client)
	client.Transport = DefaultTransport
	DefaultTransport.CloseIdleConnections()
	client.Jar, _ = cookiejar.New(options)
	return
}

func Login(steamid string, client *http.Client) (*http.Response, error) {
	addr, _ := url.Parse("http://localhost:8080/startMockLogin/" + steamid)
	return client.Do(&http.Request{Method: "GET", URL: addr})
}

func ConnectWS(client *http.Client) (*websocket.Conn, error) {
	ws := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/websocket/"}
	domain := &url.URL{Scheme: "http", Host: "localhost:8080"}

	if len(client.Jar.Cookies(domain)) == 0 {
		return nil, errors.New("Client cookiejar has no cookies D:")
	}

	header := http.Header{"Cookie": []string{client.Jar.Cookies(domain)[0].String()}}

	conn, _, err := websocket.DefaultDialer.Dial(ws.String(), header)
	return conn, err
}

func LoginAndConnectWS() (string, *websocket.Conn, *http.Client, error) {
	steamid := strconv.Itoa(rand.Int())
	client := NewClient()

	_, err := Login(steamid, client)
	if err != nil {
		return "", nil, nil, err
	}

	conn, err := ConnectWS(client)
	if err != nil {
		return "", nil, nil, err
	}

	_, err = ReadMessages(conn, InitMessages, nil)

	return steamid, conn, client, err
}

func EmitJSONWithReply(conn *websocket.Conn, req map[string]interface{}) (map[string]interface{}, error) {
	if err := conn.WriteJSON(req); err != nil {
		return nil, errors.New("Error while marshing request: " + err.Error())
	}

	resp := make(map[string]interface{})

	if err := conn.ReadJSON(&resp); err != nil {
		return nil, errors.New("Error while marshing response: " + err.Error())
	}

	return resp["data"].(map[string]interface{}), nil
}

func StartServer() *httptest.Server {
	var mux = http.NewServeMux()
	config.Constants.MockupAuth = true
	routes.SetupHTTP(mux)

	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		for err != nil {
			l, err = net.Listen("tcp", "localhost:8080")
		}
	}

	server := &httptest.Server{Listener: l, Config: &http.Server{Handler: mux}}
	go server.Start()
	return server
}

func ReadMessages(conn *websocket.Conn, n int, t *testing.T) ([]map[string]interface{}, error) {
	var messages []map[string]interface{}
	for i := 0; i < n; i++ {
		data := ReadJSON(conn)
		messages = append(messages, data)

		if t != nil {
			bytes, _ := json.MarshalIndent(data, "", "  ")
			t.Logf("%s", string(bytes))
		}
	}

	return messages, nil
}

func ReadJSON(conn *websocket.Conn) map[string]interface{} {
	reply := make(map[string]interface{})

	err := conn.ReadJSON(&reply)
	if err != nil {
		logrus.Error(err.Error())
	}

	return reply["data"].(map[string]interface{})
}