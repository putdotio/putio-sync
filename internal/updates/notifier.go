package updates

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/cenkalti/log"
	"github.com/putdotio/putio-sync/v2/internal/websocket"
)

type Notifier struct {
	HasUpdates chan string

	url              string
	handshakeTimeout time.Duration
	writeTimeout     time.Duration
	newConnectionC   chan *websocket.Websocket
	closeC           chan struct{}

	m       sync.Mutex
	token   string
	started bool
}

func NewNotifier(wsURL string, handshakeTimeout, writeTimeout time.Duration) *Notifier {
	return &Notifier{
		HasUpdates:       make(chan string, 1),
		url:              wsURL,
		handshakeTimeout: handshakeTimeout,
		writeTimeout:     writeTimeout,
		newConnectionC:   make(chan *websocket.Websocket),
		closeC:           make(chan struct{}),
	}
}

func (s *Notifier) Close() {
	close(s.closeC)
}

func (s *Notifier) Start() {
	s.m.Lock()
	if !s.started {
		go s.run()
		s.started = true
	}
	s.m.Unlock()
}

func (s *Notifier) run() {
	go s.writer()
	for {
		s.reader()
		select {
		case <-s.closeC:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}

func (s *Notifier) SetToken(token string) {
	s.m.Lock()
	s.token = token
	s.m.Unlock()
}

func (s *Notifier) notifyUpdate(name string) {
	select {
	case s.HasUpdates <- name:
	default:
	}
}

func (s *Notifier) writer() {
	var ws *websocket.Websocket
	for {
		select {
		case ws = <-s.newConnectionC:
			s.m.Lock()
			token := s.token
			s.m.Unlock()
			err := ws.Auth(token, s.writeTimeout)
			if err == websocket.ErrInvalidToken {
				s.m.Lock()
				s.token = ""
				s.m.Unlock()
			}
			if err != nil {
				log.Errorln("websocket send error:", err.Error())
				ws.Close()
				ws = nil
				break
			}
			s.notifyUpdate("WEBSOCKET_CONNECTED")
		case <-s.closeC:
			if ws != nil {
				ws.Close()
			}
			return
		}
	}
}

func (s *Notifier) reader() {
	s.m.Lock()
	token := s.token
	s.m.Unlock()
	if token == "" {
		return
	}

	ws := websocket.New(s.url)
	err := ws.Connect(s.handshakeTimeout)
	if err != nil {
		log.Errorln("websocket connect error:", err.Error())
		return
	}

	// Make sure connection is closed on return
	closed := make(chan struct{})
	defer func() {
		ws.Close()
		close(closed)
	}()
	go func() {
		select {
		case <-closed:
		case <-s.closeC:
			ws.Close()
		}
	}()

	// Notify writer for new connection
	select {
	case s.newConnectionC <- ws:
	case <-s.closeC:
		return
	}

	for {
		msg, err := ws.Recv(0)
		if err != nil {
			log.Errorln("websocket receive error:", err.Error())
			return
		}
		switch msg.Type {
		case "file_create", "file_update", "file_delete":
			name, err := readName(msg.Value)
			if err != nil {
				name = strconv.FormatInt(readID(msg.Value), 10)
			}
			s.notifyUpdate(name)
		}
	}
}

func readName(b json.RawMessage) (string, error) {
	var v struct {
		Name string `json:"name"`
	}
	err := json.Unmarshal(b, &v)
	return v.Name, err
}

func readID(b json.RawMessage) int64 {
	var v struct {
		ID int64 `json:"id"`
	}
	_ = json.Unmarshal(b, &v)
	return v.ID
}
