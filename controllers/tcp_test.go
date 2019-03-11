package controllers

import (
	"bufio"
	"bytes"
	"fmt"
	"golang-proxy-server/config"
	"net"
	"strings"
	"sync"
	"testing"
)

func TestTCPConfig_Create(t *testing.T) {
	type fields struct {
		Type string
		Host string
		Port string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "create tcp service should be success.",
			fields: fields{
				Type: "tcp",
				Host: "127.0.0.1",
				Port: "9999",
			},
			wantErr: false,
		},
		{
			name: "use (TCPConfig).Create create udp service should be false.",
			fields: fields{
				Type: "udp",
				Host: "127.0.0.1",
				Port: "9999",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcp := &TCPConfig{
				Type: tt.fields.Type,
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			if err := tcp.Create(); (err != nil) != tt.wantErr {
				t.Errorf("TCPConfig.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			// close
			tcp.Close()

		})
	}
}

func TestTCPService(t *testing.T) {
	// create a mock TCPService()
	configTmp := config.DeploySet
	config.DeploySet.Env.EnvHost = "127.0.0.1"
	config.DeploySet.Env.EnvPort = "9998"
	config.DeploySet.External.ExternalURL = "mock"

	var wg sync.WaitGroup
	wg.Add(1)

	go TCPService()

	for {
		conn, _ := net.Dial("tcp", "127.0.0.1:9998")

		if conn == nil {
			continue
		}

		t.Run("TestTCPService() mock api response should be success", func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewBuffer([]byte("test")))
			text, _ := reader.ReadString('\n')
			// send to msg
			fmt.Fprintf(conn, text+"\n")
			// listen for reply
			message1, _ := bufio.NewReader(conn).ReadString('\n')
			message2, _ := bufio.NewReader(conn).ReadString('\n')

			if strings.Index(message1, "tcp server is connect") == -1 {
				t.Errorf("TestTCPService() request mock api response '%v', want '%v'", message1, "tcp server is connect")
			}
			if strings.Index(message2, "MockExternal request data: (test) ing...") == -1 {
				t.Errorf("TestTCPService() request mock api response '%v', want '%v'", message2, "MockExternal request data: (test) ing...")
			}
		})

		// send 'quit' to 'disconnected'
		t.Run("TestTCPService() mock api send 'quit' to 'tcp server is connect' should be success", func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewBuffer([]byte("quit")))
			text, _ := reader.ReadString('\n')
			// send to msg
			fmt.Fprintf(conn, text+"\n")
			// listen for reply
			message, _ := bufio.NewReader(conn).ReadString('\n')

			t.Run("TestTCPService() mock api response should be success", func(t *testing.T) {
				if strings.Index(message, "disconnected") == -1 {
					t.Errorf("TestTCPService() request mock api response '%v', want contain '%v'", message, "disconnected")
				}
			})
		})
		wg.Done()
		break
	}

	wg.Wait()
	// restore config
	config.DeploySet.Env.EnvHost = configTmp.Env.EnvHost
	config.DeploySet.Env.EnvPort = configTmp.Env.EnvPort
	config.DeploySet.External.ExternalURL = configTmp.External.ExternalURL
}
