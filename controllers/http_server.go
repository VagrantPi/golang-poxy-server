package controllers

import (
	"encoding/json"
	"golang-proxy-server/common/model"
	"golang-proxy-server/config"
	"net/http"

	"github.com/rs/cors"
)

// Monitor - return some statistics
func Monitor(w http.ResponseWriter, req *http.Request) {
	status := <-model.TCPConnectStatusChannel
	status.ExternalAPIRequestIng = config.DeploySet.External.ExternalLimitPer - len(model.ExternalAPIRate)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
	model.TCPConnectStatusChannel <- status
}

// HTTPService - Http Service for return some statistics
func HTTPService() {
	mux := http.NewServeMux()
	mux.HandleFunc("/monitor", Monitor)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8081", handler)
}
