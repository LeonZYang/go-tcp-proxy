package proxy

import (
	"go.uber.org/atomic"
	"log"
	"sync"
)

var proxyMG = &ProxyManagement{conns: make(map[string]*Connection)}

type ProxyManagement struct {
	mu sync.RWMutex

	conns    map[string]*Connection
	stopFlag atomic.Bool
}

func (p *ProxyManagement) GetConn(name string) (*Connection, bool) {
	p.mu.RLock()
	conn, ok := p.conns[name]
	p.mu.RUnlock()
	return conn, ok
}

func (p *ProxyManagement) SetConn(name string, conn *Connection) {
	p.mu.Lock()
	p.conns[name] = conn
	p.mu.Unlock()
}

func (p *ProxyManagement) GetConns() map[string]*Connection {
	var m = make(map[string]*Connection)
	p.mu.RLock()
	for k, v := range p.conns {
		m[k] = v
	}
	p.mu.RUnlock()
	return m
}

func Start() {
	pxConfigs := GetConfig().Proxy
	for name, px := range pxConfigs {
		if px.Enabled {
			conn, err := NewConnection(px)
			if err != nil {
				log.Fatalf("new connection failed, name: %s, err: %v\n", name, err)
			}
			proxyMG.SetConn(name, conn)
			go conn.Start()
			log.Printf("new connection pool, name: %v\n, config: %+v\n", name, px)
		}
	}
}

// TODO: need implement
//func Reload() {
//	if proxyMG.stopFlag.Load() {
//		return
//	}
//	pxConfigs := GetConfig().Proxy
//	for name, px := range pxConfigs {
//		if px.Enabled {
//			oldConn, ok := proxyMG.GetConn(name)
//			if ok{
//
//			}
//
//			conn, err := NewConnection(px)
//			if err != nil {
//				log.Fatalf("new connection failed, name: %s, err: %v\n", name, err)
//			}
//			proxyMG.SetConn(name, conn)
//			go conn.Start()
//			log.Printf("new connection pool, name: %v\n, config: %+v\n", name, px)
//		}
//	}
//}

func Stop() {
	if proxyMG.stopFlag.Load() {
		return
	}
	proxyMG.mu.Lock()
	defer proxyMG.mu.Unlock()

	for name, conn := range proxyMG.conns {
		log.Printf("stop proxy, name: %v\n", name)
		conn.Stop()
	}
}
