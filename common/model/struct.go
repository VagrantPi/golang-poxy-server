package model

import (
	"golang-proxy-server/config"
	"sync"
	"time"
)

// External - external api interface, it can be any type if you impl
type External interface {
	Request() (result string, err error)
}

var (
	// ExternalRequestQueue - a request queue to external api
	ExternalRequestQueue chan interface{}

	// ExternalAPIRate - a request rate channel
	ExternalAPIRate chan struct{}

	// TCPConnectStatusChannel - tcp server status
	TCPConnectStatusChannel chan TCPConnectStatus
)

func init() {
	if config.DeploySet.External.ExternalRequestQueue == 0 {
		config.DeploySet.External.ExternalRequestQueue = 1024
	}
	ExternalRequestQueue = make(chan interface{}, config.DeploySet.External.ExternalRequestQueue)
	TCPConnectStatusChannel = make(chan TCPConnectStatus, 1)
	TCPConnectStatusChannel <- TCPConnectStatus{}

	ExternalAPIRate = make(chan struct{}, config.DeploySet.External.ExternalLimitPer)

	timeTrigger := time.Tick(5 * time.Second)
	go func() {
		for _ = range timeTrigger {
			// every second init ExternalAPIRate
			for i := 0; i < config.DeploySet.External.ExternalLimitPer; i++ {
				ExternalAPIRate <- struct{}{}
			}
		}
	}()
}

// TCPConnectStatus - tcp server status struct
type TCPConnectStatus struct {
	NowConnectNumber         int `json:"now_connect_number"`
	TotalConnectNumber       int `json:"total_connect_number"`
	FinishExternalAPIRequest int `json:"finish_external_api_request"`
	WaitExternalAPIRequest   int `json:"wait_external_api_request"`
	RejectExternalAPIRequest int `json:"reject_external_api_request"`
	AllExternalAPIRequest    int `json:"all_external_api_request"`
}

func NewTCPConnect(wg *sync.WaitGroup) {
	status := <-TCPConnectStatusChannel
	status.TotalConnectNumber++
	status.NowConnectNumber++
	TCPConnectStatusChannel <- status
	// for test
	if wg != nil {
		wg.Done()
	}
}

func TCPDisconnected(wg *sync.WaitGroup) {
	status := <-TCPConnectStatusChannel
	status.NowConnectNumber--
	TCPConnectStatusChannel <- status
	// for test
	if wg != nil {
		wg.Done()
	}
}

func ExternalAPIRequest(wg *sync.WaitGroup) {
	status := <-TCPConnectStatusChannel
	status.AllExternalAPIRequest++
	status.WaitExternalAPIRequest++
	TCPConnectStatusChannel <- status
	// for test
	if wg != nil {
		wg.Done()
	}
}

func ExternalAPIRequestFinish(wg *sync.WaitGroup) {
	status := <-TCPConnectStatusChannel
	status.FinishExternalAPIRequest++
	status.WaitExternalAPIRequest--
	TCPConnectStatusChannel <- status
	// for test
	if wg != nil {
		wg.Done()
	}
}
