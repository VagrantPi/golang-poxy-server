package controllers

import (
	"fmt"
	"golang-proxy-server/common/model"
	"golang-proxy-server/common/util"
	"golang-proxy-server/config"
	"io"
	"net"
)

// TCPConfig - TCP server connect config
type TCPConfig struct {
	Type     string
	Host     string
	Port     string
	Listener *net.TCPListener
}

// Create - a TCPConfig receiver to Create a TCPListener
func (t *TCPConfig) Create() error {
	tcpAddr, err := net.ResolveTCPAddr(t.Type, t.Host+":"+t.Port)
	if err != nil {
		config.Error.Println("[Error] TCP end point Err:", err.Error())
		return err
	}
	t.Listener, err = net.ListenTCP(t.Type, tcpAddr)
	if err != nil {
		config.Error.Println("[Error] listening:", err.Error())
		return err
	}
	return nil
}

// Close - a TCPConfig receiver to Close a TCPListener
func (t *TCPConfig) Close() {
	if t.Listener != nil {
		t.Listener.Close()
	}
}

// TCPService - tcp service
func TCPService() {
	service := &TCPConfig{
		Type: "tcp",
		Host: config.DeploySet.Env.EnvHost,
		Port: config.DeploySet.Env.EnvPort,
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
		clientID := util.GenRandString(5)
		clientIP := tcpConn.RemoteAddr().String()
		config.Info.Printf("client %v is connected : %v", clientID, clientIP)
		tcpConn.Write([]byte("tcp server is connect \n"))

		// update tcp status
		model.NewTCPConnect(nil)

		// create handle goroutine
		go tcpHandle(tcpConn, clientID, clientIP)

		// request queue worker
		go func() {
			for data := range model.ExternalRequestQueue {
				// limit external resquest rate
				<-model.ExternalAPIRate
				// config.Info.Printf("api response : %v", elem)
				requestData, ok := data.(*TCPExternalAPI)
				if ok {
					go Worker(tcpConn, clientID, clientIP, requestData)
				}
			}
		}()
	}
}

func tcpHandle(conn *net.TCPConn, clientID, clientIP string) {
	// conn.SetKeepAlivePeriod(2 * time.Second)
	// conn.SetDeadline(time.Now().Add(time.Second * 2))
	defer func() {
		config.Info.Printf("client %v disconnected ip: %v", clientID, clientIP)
		conn.Write([]byte("tcp server is connect"))
		model.TCPDisconnected(nil)
		conn.Close()
	}()

	for {
		buf := make([]byte, 1024)
		readLen, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if string(buf[0:4]) == "quit" {
			break
		}

		// external api
		var external model.External = &TCPExternalAPI{
			ClientID: clientID,
			IP:       config.DeploySet.External.ExternalURL,
			Method:   config.DeploySet.External.ExternalMethod,
			Data:     string(buf[0 : readLen-1]),
		}

		if data, err := external.Request(); err != nil {
			conn.Write([]byte(err.Error() + "\n"))
		} else {
			conn.Write([]byte(data + "\n"))
		}
		model.ExternalAPIRequest(nil)
	}
}
