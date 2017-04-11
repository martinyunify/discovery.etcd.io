package main

import (
	"log"

	_ "github.com/martinyunify/discovery.etcd.io/http"
	"github.com/rsc/devweb/slave"
)

func main() {
	log.SetFlags(0)
	slave.Main()
}
