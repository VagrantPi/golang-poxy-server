package main

import (
	"golang-proxy-server/controllers"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go controllers.HTTPService()
	go controllers.TCPService()
	wg.Wait()
}
