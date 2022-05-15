package proxy

import (
	"crypto/tls"
	"go.uber.org/atomic"
	"log"
	"net"
)

type Listener interface {
	Accept() (net.Conn, error)
	Close() error
}

type Connection struct {
	local    ListenConfig
	remote   ListenConfig
	listener Listener
	stopFlag atomic.Bool
}

func NewConnection(px ProxyConfig) (*Connection, error) {
	c := &Connection{
		local:  px.Listen,
		remote: px.Remote,
	}

	if px.Listen.TLS {
		cert, err := tls.LoadX509KeyPair(px.Listen.PrivFile, px.Listen.PubFile)
		if err != nil {
			return nil, err
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		listener, err := tls.Listen("tcp", px.Listen.Addr, config)
		if err != nil {
			return nil, err
		}
		c.listener = listener
	} else {
		listener, err := net.Listen("tcp", px.Listen.Addr)
		if err != nil {
			return nil, err
		}
		c.listener = listener
	}

	return c, nil
}

func (c *Connection) Start() {
	for {
		if c.stopFlag.Load() {
			return
		}
		conn, err := c.listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection, err :%v\n", err)
			continue
		}
		var p *Proxy
		p = New(conn, c.local, c.remote)
		//p.Nagles = *nagles

		go p.Start()
	}
}

func (c *Connection) Stop() {
	c.stopFlag.Store(true)
	c.listener.Close()
}
