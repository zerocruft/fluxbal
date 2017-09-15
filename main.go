package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"fmt"
)

func main() {

	go listener()
	waitGroup.Wait()
}

func listener() {
	router := mux.NewRouter()
	router.HandleFunc("/flux", handleFluxRequest).Methods("GET")
	router.HandleFunc("/control/cluster/ping", handleControlClusterPing).Methods("POST")

	err := http.ListenAndServe(":"+strconv.Itoa(flgPort), router)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleFluxRequest(response http.ResponseWriter, request *http.Request) {

	http.Redirect(response, request, "ws://localhost:8283/flux", 302)
}

func handleControlClusterPing(response http.ResponseWriter, request *http.Request) {


}
