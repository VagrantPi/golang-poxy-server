package controllers

import (
	"encoding/json"
	"golang-proxy-server/common/model"
	"net/http"
)

// Monitor - return some statistics
func Monitor(w http.ResponseWriter, req *http.Request) {
	status := <-model.TCPConnectStatusChannel
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
	model.TCPConnectStatusChannel <- status
}

// HTTPService - Http Service for return some statistics
func HTTPService() {
	http.HandleFunc("/monitor", Monitor)
	http.ListenAndServe(":8081", nil)
}
