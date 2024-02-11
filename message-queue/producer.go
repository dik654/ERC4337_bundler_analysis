package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/exp/slog"
)

type Producer interface {
	Start() error
}

type HTTPProducer struct {
	listenAddr string
	server     *Server
	producech  chan<- Message
}

func NewHTTPProducer(listenAddr string, producech chan Message) *HTTPProducer {
	return &HTTPProducer{
		listenAddr: listenAddr,
		producech:  producech,
	}
}

func (p *HTTPProducer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		// ex) /publish/topic_1라면 => publish/topic_1
		path = strings.TrimPrefix(r.URL.Path, "/")
		// ["publish", "topic_1"]
		parts = strings.Split(path, "/")
	)
	// commit
	if r.Method == "GET" {
	}

	// POST 요청이라면
	if r.Method == "POST" {
		// ["publish", "topic_1"]같은 형태가 아니라면 종료
		if len(parts) != 2 {
			fmt.Println("invalid action")
			return
		}

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			// 오류 처리
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		// produce 채널에 topic과 data를 담아서 전송
		p.producech <- Message{
			Data:  body,
			Topic: parts[1],
		}
	}

	fmt.Println(parts)
}

func (p *HTTPProducer) Start() error {
	// http 서버 생성
	slog.Info("HTTP transport started", "port", p.listenAddr)
	return http.ListenAndServe(p.listenAddr, p)
}
