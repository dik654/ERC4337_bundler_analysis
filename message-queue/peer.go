package main

import (
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

type Peer interface {
	Send([]byte) error
}

type WSPeer struct {
	conn   *websocket.Conn
	server *Server
}

func NewWSPeer(conn *websocket.Conn, s *Server) *WSPeer {
	p := &WSPeer{
		conn:   conn,
		server: s,
	}
	// web socket peer에서 반복해서 읽기
	go p.readLoop()
	return p
}

func (p *WSPeer) readLoop() {
	var msg WSMessage
	for {
		// web socket peer에서 JSON 타입의 데이터를 읽은 뒤(Action, topics)
		if err := p.conn.ReadJSON(&msg); err != nil {
			slog.Error("ws peer read error", "err", err)
			return
		}
		// subscribe action일 경우 해당 topic들 안의 data들을 모두 peer에게 전송
		if err := p.handleMessage(msg); err != nil {
			slog.Error("ws peer handle msg error", "err", err)
			return
		}
	}
}

func (p *WSPeer) handleMessage(msg WSMessage) error {
	// subscribe action일 경우 해당 topic들 안의 data들을 모두 peer에게 전송
	if msg.Action == "subscribe" {
		p.server.AddPeerToTopics(p, msg.Topics...)
	}
	return nil
}

func (p *WSPeer) Send(b []byte) error {
	// binary를 web socket으로 전송
	return p.conn.WriteMessage(websocket.TextMessage, b)
}
