package routes

import "net/http"

type route struct {
	pattern string
	handler http.HandlerFunc
}

var httpRoutes = []route{
	{"/", controllers.MainHandler},

}