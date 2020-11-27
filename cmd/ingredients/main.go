package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/apex/gateway"
	"github.com/schollz/ingredients"
)

func main() {
	port := flag.Int("port", -1, "specify a port to use http rather than AWS Lambda")
	flag.Parse()

	listener := gateway.ListenAndServe
	portStr := fmt.Sprintf(":%d", *port)
	if *port != -1 {
		portStr = fmt.Sprintf(":%d", *port)
		listener = http.ListenAndServe
	}

	fmt.Printf("listening on port %d\n", *port)

	getIngredientsHandler := http.HandlerFunc(getIngredients)
	http.Handle("/", getIngredientsHandler)

	log.Fatal(listener(portStr, logRequest(http.DefaultServeMux)))
}

func getIngredients(resp http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	url := query.Get("url")

	var r *ingredients.Recipe
	var err error

	r, err = ingredients.NewFromURL(url)
	if err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte("Missing 'url' in url parameter"))
		return
	}

	ing := r.IngredientList()
	b, _ := json.MarshalIndent(ing, "", "    ")

	resp.WriteHeader(200)
	resp.Write([]byte(b))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
