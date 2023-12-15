package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

var (
	listenAddr string
	proxyURL   string
)

func init() {
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:443", "listening address")
	flag.StringVar(&proxyURL, "proxy", "socks5://127.0.0.1:1080", "proxy")
	flag.Parse()
}

func main() {
	u, err := url.Parse(proxyURL)
	if err != nil {
		log.Fatalf("parse proxy %s error: %v\n", proxyURL, err)
	}

	p := &Proxy{
		proxy: u,
	}
	if err := p.ListenAndServe("tcp", listenAddr); err != nil {
		log.Fatal(err)
	}
}

type Proxy struct {
	l     net.Listener
	proxy *url.URL
}

func (p *Proxy) ListenAndServe(network, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		return fmt.Errorf("create listener error: %w", err)
	}
	p.l = l
	return p.Serve()
}

func (p *Proxy) Serve() error {
	for {
		conn, err := p.l.Accept()
		if err != nil {
			return fmt.Errorf("accept new conn error: %w", err)
		}

		go func() {
			err := p.handleProxy(conn)
			if err != nil {
				log.Println(conn.RemoteAddr(), err)
			}
		}()
	}
}

func (p *Proxy) handleProxy(conn net.Conn) error {
	defer conn.Close()

	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		return fmt.Errorf("set read deadline for ClientHello error: %w", err)
	}

	var handshakeBuf bytes.Buffer
	hostname, _, err := extractSNI(io.TeeReader(conn, &handshakeBuf))
	if err != nil {
		return fmt.Errorf("extract SNI error: %w", err)
	}

	if hostname == "" {
		return fmt.Errorf("empty SNI")
	}

	network := conn.LocalAddr().Network()
	port := conn.LocalAddr().(*net.TCPAddr).Port

	if err = conn.SetReadDeadline(time.Time{}); err != nil {
		return fmt.Errorf("clearing read deadline for ClientHello error: %w", err)
	}

	log.Println("connect to", hostname)
	backend, err := p.getBackend(network, hostname+fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("dial %q error: %w", hostname, err)
	}
	defer backend.Close()

	if _, err = io.Copy(backend, &handshakeBuf); err != nil {
		return fmt.Errorf("replay handshake to %q error: %w", backend, err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go pipe(&wg, conn, backend)
	go pipe(&wg, backend, conn)
	wg.Wait()
	return nil
}

func (p *Proxy) dialWithProxy(network string, host string, timeout time.Duration) (net.Conn, error) {
	dailer, err := proxy.FromURL(p.proxy, &net.Dialer{Timeout: timeout, KeepAlive: 30 * time.Second})
	if err != nil {
		return nil, err
	}
	return dailer.Dial(network, host)
}

func (p *Proxy) getBackend(network string, host string) (net.Conn, error) {
	return p.dialWithProxy(network, host, time.Second*10)
}

func pipe(wg *sync.WaitGroup, src, dst net.Conn) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("%s<>%s -> %s<>%s: %s\n", src.RemoteAddr(), src.LocalAddr(), dst.LocalAddr(), dst.RemoteAddr(), err)
	}
	wg.Done()
}
