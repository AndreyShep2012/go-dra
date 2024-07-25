package main

import (
	"flag"
	"go-dra/client"
	"log"
	"time"
)

func main() {
	raddr := flag.String("raddr", "", "remote address of diameter server")
	transport := flag.String("transport", "tcp", "the address of diameter server")
	laddr := flag.String("laddr", "", "local address of diameter client")
	flag.Parse()

	if raddr == nil || *raddr == "" {
		log.Fatal("server address is empty, use 'raddr=some_server:port' flag")
	}

	if laddr == nil || *laddr == "" {
		log.Fatal("local address is empty, use 'laddr=some_server:port' flag")
	}

	log.Printf("remote addr: %s, local addr:%s, transport: %s\n", *raddr, *laddr, *transport)

	_, err := client.MakeConnection(*laddr, *raddr, *transport)
	if err != nil {
		log.Fatal(err.Error())
	}

	time.Sleep(time.Second * 5)
}
