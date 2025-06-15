package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

func main2() {
	fp := os.Getenv("FE_PORT")

	if fp == "" {
		log.Fatal("FE_PORT env var is missing.")
	}

	bp := os.Getenv("BE_PORT")

	if bp == "" {
		log.Fatal("BE_PORT env var is missing.")
	}

	url, err := url.Parse("http://localhost:" + bp)

	if err != nil {
		log.Fatalf("Error parsing backend url: %v", err)
	}

	rp := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = url.String()
		},
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				ta, err := net.ResolveTCPAddr(network, addr)
				if err != nil {
					return nil, err
				}
				return net.DialTCP(network, nil, ta)
			},
		},
	}

	log.Fatal(http.ListenAndServe(":"+fp, rp))

}

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", ":8080")

	list, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Error listen on url: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		log.Print("go routine")
		for {
			log.Print("list.accept")
			conn, err := list.Accept()
			log.Print("pos list.accept")
			if err != nil {
				if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
					log.Fatalf("Error: %v", err)
					return
				}
				continue
			}
			log.Print("pre handleConn")
			go handleConn(conn.(*net.TCPConn))
		}
	}()

	wg.Wait()

}

func handleConn(conn *net.TCPConn) {
	log.Print("handleConn")

	defer func() {
		_ = conn.Close()
	}()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Error Buffer: %v", err)
		return
	}

	payload := buf[:n]
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(payload)))
	if err != nil {
		log.Fatalf("Error Request: %v", err)
		return
	}

	log.Print(req.ProtoMajor)

	if req.ProtoMajor >= 2 {
		err = handleHttp2(bytes.NewBuffer(payload), conn)
		log.Fatalf("Error: %v", err)
	} else {
		err = handleHttpReq(req, conn)
		log.Fatalf("Error: %v", err)
	}

}

func handleHttpReq(req *http.Request, w net.Conn) error {
	resp := &http.Response{}
	log.Print(req.Host)
	// Check the host of the request
	if req.Host == "banned.hostfactor.io" {
		resp.StatusCode = http.StatusNotFound
	} else {
		// Required to forward the request
		req.RequestURI = ""
		var err error
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
	}
	return resp.Write(w)
}

func handleHttp2(initial io.Reader, conn net.Conn) error {
	defer func() {
		_ = conn.Close()
	}()
	dataBuffer := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(conn, dataBuffer)
	f := http2.NewFramer(conn, conn)
	err := f.WriteSettingsAck()
	if err != nil {
		return err
	}

	f = http2.NewFramer(io.Discard, reader)
	//decoder := hpack.NewDecoder(1024, nil)
	log.Print("URL")
	/* auth := ""
	for auth == "" {
		frame, err := f.ReadFrame()
		if err != nil {
			return err
		}

		switch t := frame.(type) {
		case *http2.HeadersFrame:
			out, err := decoder.DecodeFull(t.HeaderBlockFragment())
			if err != nil {
				return err
			}

			for _, v := range out {
				if v.Name == ":authority" {
					auth = v.Value
				}
			}
		}
	}

	if auth == "blocked.hostfactor.io" {
		return nil
	} */
	log.Print("DIAL")
	dialer, err := net.Dial("tcp", "localhost:50080")
	if err != nil {
		return err
	}

	_ = dialer.SetReadDeadline(time.Now().Add(5 * time.Second))

	wg := sync.WaitGroup{}
	wg.Add(1)
	dataSent := int64(0)
	go func() {
		// Copy any data we receive from the host into the original connection
		dataSent, err = io.Copy(conn, dialer)
		wg.Done()
	}()

	_, err = io.Copy(dialer, io.MultiReader(initial, dataBuffer, conn))
	wg.Wait()

	if errors.Is(err, os.ErrDeadlineExceeded) && dataSent > 0 {
		return nil
	}
	return err
}
