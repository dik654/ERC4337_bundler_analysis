package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

var upgrader = websocket.Upgrader{}

type Consumer interface {
	Start() error
}

type WSConsumer struct {
	ListenAddr string
	server     *Server
}

func NewWSConsumer(listenAddr string, server *Server) *WSConsumer {
	return &WSConsumer{
		ListenAddr: listenAddr,
		server:     server,
	}
}

func (ws *WSConsumer) Start() error {
	// web socket 서버 실행
	slog.Info("websocket consumer started", "port", ws.ListenAddr)
	return http.ListenAndServe(ws.ListenAddr, ws)
}

// web socket 서버를 열기 위해서 필요한 과정
func (ws *WSConsumer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// websocket 연결, :4000
	p := NewWSPeer(conn, ws.server)
	// 해당 ws peer가 연결됐음을 true로 변경
	ws.server.AddConn(p)
}

type WSMessage struct {
	Action string   `json:"action"`
	Topics []string `json:"topics"`
}
