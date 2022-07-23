package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"go-tcp-proxy/proxy"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	cf = flag.String("cf", "server.toml", "server config")
)

func main() {
	flag.Parsed()
	_, err := toml.DecodeFile(*cf, proxy.GetConfig())
	if err != nil {
		log.Fatalf("parse config err: %v\n", err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	proxy.Start()

EXIT:
	for {
		sig := <-sc
		fmt.Println("received signal:", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			break EXIT
		case syscall.SIGHUP:
			log.Println("receive reload signal")
			//proxy.Reload()
		default:
			break EXIT
		}
	}

	proxy.Stop()
}
