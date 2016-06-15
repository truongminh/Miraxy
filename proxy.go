package main

import (
	"net"
    "log"
	"sync"
)

// Proxy connections from Listen to Backend.
type Proxy struct {
	Listen   string
	Backend  string
	listener net.Listener
}

func (p *Proxy) Run() error {
	var err error
	if p.listener, err = net.Listen("tcp", p.Listen); err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for {
		if conn, err := p.listener.Accept(); err == nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.handle(conn)
			}()
		} else {
			return nil
		}
	}
	wg.Wait()
	return nil
}

func (p *Proxy) Close() error {
	return p.listener.Close()
}

func (p *Proxy) handle(upConn net.Conn) {
	defer upConn.Close()
	log.Printf("accepted: %s", upConn.RemoteAddr())
	downConn, err := net.Dial("tcp", p.Backend)
	if err != nil {
		log.Fatalf("unable to connect to %s: %s", p.Backend, err)
		return
	}
	defer downConn.Close()
	if err := Pipe(upConn, downConn); err != nil {
		log.Printf("pipe failed: %s", err)
	} else {
		log.Printf("disconnected: %s", upConn.RemoteAddr())
	}
}