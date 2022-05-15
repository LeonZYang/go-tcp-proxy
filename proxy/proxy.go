package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

type Proxy struct {
	sentBytes     uint64
	receivedBytes uint64
	local, remote ListenConfig
	lconn, rconn  net.Conn
	erred         bool
	errsig        chan bool

	// Settings
	Nagles bool

	exitedChan chan struct{}
}

func New(lconn net.Conn, local, remote ListenConfig) *Proxy {
	return &Proxy{
		lconn:      lconn,
		local:      local,
		remote:     remote,
		erred:      false,
		errsig:     make(chan bool),
		exitedChan: make(chan struct{}),
	}
}

type setNoDelayer interface {
	SetNoDelay(bool) error
}

// Start - open connection to remote and start proxying data.
func (p *Proxy) Start() {
	defer p.lconn.Close()

	var err error
	//connect to remote
	p.rconn, err = Dial(p.remote)
	if err != nil {
		log.Printf("Remote connection failed: %v\n", err)
		return
	}
	defer p.rconn.Close()

	//nagles?
	if p.Nagles {
		if conn, ok := p.lconn.(setNoDelayer); ok {
			conn.SetNoDelay(true)
		}
		if conn, ok := p.rconn.(setNoDelayer); ok {
			conn.SetNoDelay(true)
		}
	}

	//display both ends
	log.Printf("open connection, local: %v, remote: %v\n", p.local.Addr, p.remote.Addr)
	//bidirectional copy
	go p.pipe(p.lconn, p.rconn)
	go p.pipe(p.rconn, p.lconn)

	//wait for close...
	select {
	case <-p.exitedChan:
	case <-p.errsig:
	}

	log.Printf("close connection, local: %v, remote: %v, sent: %v, received: %v\n", p.local.Addr, p.remote.Addr, p.sentBytes, p.receivedBytes)
	close(p.exitedChan)
}

func (p *Proxy) err(s string, err error) {
	if p.erred {
		return
	}
	if err != io.EOF {
		log.Printf(s, err)
	}
	select {
	case p.errsig <- true:
	case <-p.exitedChan:
	}
	p.erred = true
}

// pip
func (p *Proxy) pipe(src, dst net.Conn) {
	islocal := src == p.lconn

	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)
	for {
		err := src.SetReadDeadline(time.Now().Add(2 * time.Minute))
		if err != nil {
			p.err("set read timeout failed '%s'\n", err)
			return
		}
		n, err := src.Read(buff)
		if err != nil {
			p.err("read failed '%s'\n", err)
			return
		}
		b := buff[:n]
		_ = dst.SetWriteDeadline(time.Now().Add(2 * time.Minute))
		//write out result
		n, err = dst.Write(b)
		if err != nil {
			p.err("write failed '%s'\n", err)
			return
		}
		if islocal {
			p.sentBytes += uint64(n)
		} else {
			p.receivedBytes += uint64(n)
		}
	}
}

func Dial(lc ListenConfig) (net.Conn, error) {
	if lc.TLS {
		return dialTcpTLS(lc)
	}

	return dialTcp(lc.Addr)
}

func dialTcp(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	return conn, err
}

func dialTcpTLS(lc ListenConfig) (net.Conn, error) {
	cert, err := tls.LoadX509KeyPair(lc.PrivFile, lc.PubFile)
	if err != nil {
		return nil, err
	}

	certBytes, err := ioutil.ReadFile(lc.Ca)
	if err != nil {
		return nil, err
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		return nil, errors.New("failed to parse root certificate")
	}

	cf := &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", lc.Addr, cf)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
