package model

import (
	"sync"
	"testing"
)

func TestNewTCPConnect(t *testing.T) {
	tests := []struct {
		name               string
		times              int
		totalConnectNumber int
		nowConnectNumber   int
	}{
		{
			name:               "TestNewTCPConnect tcp connect +10 user totalConnectNumber should be 10",
			times:              10,
			totalConnectNumber: 10,
			nowConnectNumber:   10,
		},
		{
			name:               "TestNewTCPConnect tcp connect +10 user totalConnectNumber should be 20",
			times:              10,
			totalConnectNumber: 20,
			nowConnectNumber:   20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			for i := 0; i < tt.times; i++ {
				wg.Add(1)
				go NewTCPConnect(&wg)
			}
			wg.Wait()
			status := <-TCPConnectStatusChannel
			if status.TotalConnectNumber != tt.totalConnectNumber {
				t.Errorf("status.TotalConnectNumber = %v, want %v", status.TotalConnectNumber, tt.totalConnectNumber)
			}
			if status.NowConnectNumber != tt.nowConnectNumber {
				t.Errorf("status.NowConnectNumber = %v, want %v", status.NowConnectNumber, tt.nowConnectNumber)
			}
			TCPConnectStatusChannel <- status
		})
	}

	// reset status
	status := <-TCPConnectStatusChannel
	status.TotalConnectNumber = 0
	status.NowConnectNumber = 0
	TCPConnectStatusChannel <- status
}

func TestTCPDisconnected(t *testing.T) {
	// mock connect status
	status := <-TCPConnectStatusChannel
	status.NowConnectNumber = 20
	TCPConnectStatusChannel <- status

	tests := []struct {
		name             string
		times            int
		nowConnectNumber int
	}{
		{
			name:             "TestTCPDisconnected -10 user nowConnectNumber should be 10",
			times:            10,
			nowConnectNumber: 10,
		},
		{
			name:             "TestTCPDisconnected -10 user nowConnectNumber should be 0",
			times:            10,
			nowConnectNumber: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			for i := 0; i < tt.times; i++ {
				wg.Add(1)
				go TCPDisconnected(&wg)
			}
			wg.Wait()
			status := <-TCPConnectStatusChannel
			if status.NowConnectNumber != tt.nowConnectNumber {
				t.Errorf("status.NowConnectNumber = %v, want %v", status.NowConnectNumber, tt.nowConnectNumber)
			}
			TCPConnectStatusChannel <- status
		})
	}

	// reset status
	status = <-TCPConnectStatusChannel
	status.TotalConnectNumber = 0
	status.NowConnectNumber = 0
	TCPConnectStatusChannel <- status
}

func TestExternalAPIRequest(t *testing.T) {
	tests := []struct {
		name                   string
		times                  int
		allExternalAPIRequest  int
		waitExternalAPIRequest int
	}{
		{
			name:  "TestExternalAPIRequest +40 user allExternalAPIRequest should be 40, waitExternalAPIRequest should be 40",
			times: 40,
			allExternalAPIRequest:  40,
			waitExternalAPIRequest: 40,
		},
		{
			name:  "TestExternalAPIRequest +60 user allExternalAPIRequest should be 100, waitExternalAPIRequest should be 100",
			times: 60,
			allExternalAPIRequest:  100,
			waitExternalAPIRequest: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			for i := 0; i < tt.times; i++ {
				wg.Add(1)
				go ExternalAPIRequest(&wg)
			}
			wg.Wait()
			status := <-TCPConnectStatusChannel
			if status.AllExternalAPIRequest != tt.allExternalAPIRequest {
				t.Errorf("status.AllExternalAPIRequest = %v, want %v", status.AllExternalAPIRequest, tt.allExternalAPIRequest)
			}
			if status.WaitExternalAPIRequest != tt.waitExternalAPIRequest {
				t.Errorf("status.WaitExternalAPIRequest = %v, want %v", status.WaitExternalAPIRequest, tt.waitExternalAPIRequest)
			}
			TCPConnectStatusChannel <- status
		})
	}

	// reset status
	status := <-TCPConnectStatusChannel
	status.AllExternalAPIRequest = 0
	status.WaitExternalAPIRequest = 0
	TCPConnectStatusChannel <- status
}

func TestExternalAPIRequestFinish(t *testing.T) {
	// mock connect status
	status := <-TCPConnectStatusChannel
	status.FinishExternalAPIRequest = 40
	status.WaitExternalAPIRequest = 60
	TCPConnectStatusChannel <- status

	tests := []struct {
		name                     string
		times                    int
		finishExternalAPIRequest int
		waitExternalAPIRequest   int
	}{
		{
			name:  "TestExternalAPIRequestFinish +30 user finishExternalAPIRequest should be 70, waitExternalAPIRequest should be 30",
			times: 30,
			finishExternalAPIRequest: 70,
			waitExternalAPIRequest:   30,
		},
		{
			name:  "TestExternalAPIRequestFinish +30 user finishExternalAPIRequest should be 100, waitExternalAPIRequest should be 0",
			times: 30,
			finishExternalAPIRequest: 100,
			waitExternalAPIRequest:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			for i := 0; i < tt.times; i++ {
				wg.Add(1)
				go ExternalAPIRequestFinish(&wg)
			}
			wg.Wait()
			status := <-TCPConnectStatusChannel
			if status.FinishExternalAPIRequest != tt.finishExternalAPIRequest {
				t.Errorf("status.FinishExternalAPIRequest = %v, want %v", status.FinishExternalAPIRequest, tt.finishExternalAPIRequest)
			}
			if status.WaitExternalAPIRequest != tt.waitExternalAPIRequest {
				t.Errorf("status.WaitExternalAPIRequest = %v, want %v", status.WaitExternalAPIRequest, tt.waitExternalAPIRequest)
			}
			TCPConnectStatusChannel <- status
		})
	}

	// reset status
	status = <-TCPConnectStatusChannel
	status.FinishExternalAPIRequest = 0
	status.WaitExternalAPIRequest = 0
	TCPConnectStatusChannel <- status
}
