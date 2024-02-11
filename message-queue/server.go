package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"golang.org/x/exp/slog"
)

type Message struct {
	// 어떤 토픽으로 쏠지
	Topic string `json:"topic"`
	// 어떤 데이터를 담고있는지
	Data []byte `json:"data"`
}

type Config struct {
	HTTPListenAddr    string
	WSListenAddr      string
	StoreProducerFunc StoreProducerFunc
}

type Server struct {
	*Config

	mu    sync.RWMutex
	peers map[Peer]bool

	topics    map[string]Storer
	consumers []Consumer
	producers []Producer
	producech chan Message
	quitch    chan struct{}
}

// initialize
func NewServer(cfg *Config) (*Server, error) {
	producech := make(chan Message)
	s := &Server{
		Config:    cfg,
		topics:    make(map[string]Storer),
		quitch:    make(chan struct{}),
		peers:     make(map[Peer]bool),
		producech: producech,
		producers: []Producer{
			NewHTTPProducer(cfg.HTTPListenAddr, producech),
		},
		consumers: []Consumer{},
	}
	s.consumers = append(s.consumers, NewWSConsumer(cfg.WSListenAddr, s))
	return s, nil
}

func (s *Server) Start() {
	// web socket 리스닝 서버 생성
	for _, consumer := range s.consumers {
		go func(c Consumer) {
			if err := c.Start(); err != nil {
				fmt.Println(err)
			}
		}(consumer)
	}
	// http 리스닝 서버 생성
	for _, producer := range s.producers {
		go func(p Producer) {
			if err := p.Start(); err != nil {
				fmt.Println(err)
			}
		}(producer)
	}
	// 서버들로부터 받은 신호 읽기
	s.loop()
}

func (s *Server) loop() {
	for {
		select {
		// 종료 신호가 넘어오면 main 종료
		case <-s.quitch:
			return
		//
		case msg := <-s.producech:
			// 받은 message를 store에 추가
			index, err := s.publish(msg)
			if err != nil {
				slog.Error("failed to publish", err)
			} else {
				// 메세지 받았음을 로깅
				slog.Info("produced message", "index", index)
			}
		}
	}
}

func (s *Server) publish(msg Message) (int, error) {
	// topic에 해당하는 memory store 저장공간 가져와서
	store := s.getStoreForTopic(msg.Topic)
	// data [][]byte에 메세지의 데이터 저장
	return store.Push(msg.Data)
}

// topic에 해당하는 memory store 저장공간 가져오기
func (s *Server) getStoreForTopic(topic string) Storer {
	// 해당 topic이 없었다면 새롭게 생성
	if _, ok := s.topics[topic]; !ok {
		// topic의 저장소 추가
		s.topics[topic] = s.StoreProducerFunc()
		slog.Info("created new topic", "topic", topic)
	}
	// 저장소 리턴
	return s.topics[topic]
}

func (s *Server) AddConn(p Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Info("added new peer", "peer", p)
	s.peers[p] = true
}

func (s *Server) AddPeerToTopics(p Peer, topics ...string) {
	// 해당 topic 저장소들에 들어있는 binary 데이터들을 peer에 전송
	for _, topic := range topics {
		store := s.getStoreForTopic(topic)
		size := store.Len()
		for i := 0; i < size; i++ {
			b, _ := store.Get(i)
			message := Message{
				Topic: topic,
				Data:  b,
			}
			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Fatal(err)
			}
			p.Send(messageBytes)
		}
	}
	slog.Info("adding peer to topics", "topics", topics, "peer", p)
}
