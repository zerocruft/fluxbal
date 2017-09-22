package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zerocruft/capacitor"
	"github.com/zerocruft/fluxbal/state"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {

	go listener()
	go background()
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
	lightestNode := state.GetNodeWithLightestLoad()
	if lightestNode.Address == "" {
		response.WriteHeader(http.StatusNoContent)
		return
	}

	http.Redirect(response, request, "ws://"+lightestNode.Address, 302)
}

func handleControlClusterPing(response http.ResponseWriter, request *http.Request) {

	//TODO do some research on who this is? maybe authentication
	defer request.Body.Close()

	msgBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusNotAcceptable)
		return
	}

	ping := capacitor.FluxPing{}
	err = json.Unmarshal(msgBytes, &ping)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusNotAcceptable)
		return
	}

	state.AddNode(ping.Node, ping.NumberOfConnections)

	pong := capacitor.FluxPong{
		Peers: state.ToNodeSlice(),
	}
	pongBytes, err := json.Marshal(pong)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusNotAcceptable)
		return
	}
	response.Write(pongBytes)
}

func background() {
	for {
		time.Sleep(5 * time.Second)
		for _, node := range state.CopyOfNodes() {
			if time.Now().Add(-20 * time.Second).After(node.LastPing) {
				state.RemoveNode(node.Node.Address)
			}
		}
	}
}
