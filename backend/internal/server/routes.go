package server

import "net/http"

var mux = http.NewServeMux()

type Routes struct {
	PARSE    string
	GENERATE string
}

var RoutesInstance = Routes{
	PARSE:    "/parse",
	GENERATE: "/generate",
}

func GetMuxInstance() *http.ServeMux {
	return mux
}
