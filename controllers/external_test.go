package controllers

import (
	"bufio"
	"bytes"
	"fmt"
	"golang-proxy-server/common/model"
	"golang-proxy-server/config"
	"net"
	"strings"
	"sync"
	"testing"
)

func TestTCPExternalAPI_Request(t *testing.T) {
	type fields struct {
		ClientID string
		IP       string
		Method   string
		Data     interface{}
	}
	tests := []struct {
		name           string
		fields         fields
		wantContainStr string
		wantErr        bool
	}{
		{
			name: "TCPExternalAPI Request string data should by success",
			fields: fields{
				ClientID: "mock_client_id",
				IP:       "mock",
				Method:   "GET",
				Data:     "key=value",
			},
			wantContainStr: "MockExternal request data: ",
			wantErr:        false,
		},
		{
			name: "TCPExternalAPI Request []byte data should by success",
			fields: fields{
				ClientID: "mock_client_id",
				IP:       "mock",
				Method:   "GET",
				Data:     []byte("test"),
			},
			wantContainStr: "MockExternal requesting...",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			external := &TCPExternalAPI{
				ClientID: tt.fields.ClientID,
				IP:       tt.fields.IP,
				Method:   tt.fields.Method,
				Data:     tt.fields.Data,
			}
			gotStr, err := external.Request()
			if (err != nil) != tt.wantErr {
				t.Errorf("TCPExternalAPI.Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.Index(gotStr, tt.wantContainStr) == -1 {
				t.Errorf("TCPExternalAPI.Request() = %v, want contain %v", gotStr, tt.wantContainStr)
			}
		})
		<-model.ExternalRequestQueue
	}
}

func TestWorker(t *testing.T) {
	// create a mock mockTCPService()
	var conn *net.TCPConn
	var wg sync.WaitGroup
	wg.Add(1)
	// go func(conn **net.TCPConn, wg *sync.WaitGroup) {
	go mockTCPService(&conn, &wg)
	// }(&conn, &wg)

	// connect mock tcp server
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for {
			clientConn, _ := net.Dial("tcp", "127.0.0.1:9996")

			if clientConn == nil {
				continue
			}
			reader := bufio.NewReader(bytes.NewBuffer([]byte("test")))
			text, _ := reader.ReadString('\n')
			// send to msg
			fmt.Fprintf(clientConn, text+"\n")

			if clientConn != nil {
				break
			}
		}
		wg.Done()
	}(&wg)
	wg.Wait()
	t.Run("TestWorker model.ExternalRequestQueue len should be 1", func(t *testing.T) {
		if len(model.ExternalRequestQueue) != 1 {
			t.Errorf("len(model.ExternalRequestQueue) = %v, want 1", len(model.ExternalRequestQueue))
			return
		}
	})

	t.Run("TestWorker after Work() model.ExternalRequestQueue len should be 0", func(t *testing.T) {
		var wg2 sync.WaitGroup
		wg2.Add(1)
		// fmt.Println("conn:", conn)
		go func(wg2 *sync.WaitGroup) {
			for data := range model.ExternalRequestQueue {
				fmt.Println("data:", data)
				// config.Info.Printf("api response : %v", elem)
				requestData, ok := data.(*TCPExternalAPI)
				fmt.Println("requestData:", requestData)
				if ok {
					go Worker(conn, "testClientID", "http://127.0.0.1", requestData)
				}
			}
			wg2.Done()
		}(&wg2)

		if len(model.ExternalRequestQueue) != 0 {
			t.Errorf("len(model.ExternalRequestQueue) = %v, want 0", len(model.ExternalRequestQueue))
			return
		}
		wg2.Wait()
	})

}

// func mockTCPService(conn **net.TCPConn) {
func mockTCPService(conn **net.TCPConn, wg *sync.WaitGroup) {
	service := &TCPConfig{
		Type: "tcp",
		Host: "127.0.0.1",
		Port: "9996",
	}
	if err := service.Create(); err != nil {
		fmt.Println("create tcp listener Err: ", err.Error())
	}
	defer service.Close()
	config.Info.Println("Listening on " + service.Host + ":" + service.Port)

	for {
		tcpConn, err := service.Listener.AcceptTCP()
		if err != nil {
			continue
		}
		conn = &tcpConn
		go tcpHandle(tcpConn, "testClientID", "http://127.0.0.1")
		wg.Done()
		break
	}
}

// func Test_mockTCPService(t *testing.T) {
// 	type args struct {
// 		conn **net.TCPConn
// 		wg   *sync.WaitGroup
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockTCPService(tt.args.conn, tt.args.wg)
// 		})
// 	}
// }
