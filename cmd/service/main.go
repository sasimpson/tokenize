package main

import (
	"fmt"
	"net/http"
	"tokenize/api"
)

func main() {

	routes := api.Routes()

	server := &http.Server{
		Addr:    ":8000",
		Handler: routes,
	}

	fmt.Println("starting service on", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
