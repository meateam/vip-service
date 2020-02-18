package main

import (
	"github.com/meateam/vip-service/server"
)

func main() {
	server.NewServer(nil).Serve(nil)
}
