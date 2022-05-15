package proxy

import (
	"go.uber.org/atomic"
	"log"
	"sync"
)

var proxyMG = &ProxyManagement{pools: make(map[string]*Connection)}

type ProxyManagement struct {
	sync.Mutex

	pools    map[string]*Connection
	stopFlag atomic.Bool
}

func Start() {
	pxConfigs := GetConfig().Proxy
	for name, px := range pxConfigs {
		if px.Enabled {
			conn, err := NewConnection(px)
			if err != nil {
				log.Fatalf("new connection failed, name: %s, err: %v\n", name, err)
			}
			proxyMG.pools[name] = conn
			go conn.Start()
			log.Printf("new connection pool, name: %v\n, config: %+v\n", name, px)
		}
	}
}

func Reload() {
	if proxyMG.stopFlag.Load() {
		return
	}
	pxConfigs := GetConfig().Proxy
	for _, px := range pxConfigs {
		if px.Enabled {

		}
	}
}

func Stop() {
	if proxyMG.stopFlag.Load() {
		return
	}
	proxyMG.Lock()
	defer proxyMG.Unlock()

	for name, conn := range proxyMG.pools {
		log.Printf("stop proxy, name: %v\n", name)
		conn.Stop()
	}
}
