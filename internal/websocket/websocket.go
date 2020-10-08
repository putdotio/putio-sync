package websocket

import (
	"errors"
	"net/http"
	"time"

	"github.com/cenkalti/log"
	"github.com/gorilla/websocket"
)

var ErrInvalidToken = errors.New("invalid token")

type Websocket struct {
	url  string
	conn *websocket.Conn
}

func New(wsURL string) *Websocket {
	return &Websocket{
		url: wsURL,
	}
}

func (w *Websocket) Connect(handshakeTimeout time.Duration) error {
	log.Debugf("Connecting to websocket: %s", w.url)
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: handshakeTimeout,
	}
	conn, _, err := dialer.Dial(w.url, nil) // nolint:bodyclose
	if err != nil {
		return err
	}
	log.Debugf("Connected to websocket: %s", w.url)
	w.conn = conn
	return nil
}

func (w *Websocket) Close() error {
	if w.conn == nil {
		return nil
	}
	return w.conn.Close()
}

func (w *Websocket) Auth(token string, timeout time.Duration) error {
	_ = w.conn.SetWriteDeadline(deadlineForTimeout(timeout))
	err := w.conn.WriteMessage(websocket.TextMessage, []byte(token))
	var closeError *websocket.CloseError
	if errors.As(err, &closeError) && closeError.Code == 4001 {
		return ErrInvalidToken
	}
	return err
}

type IncomingMessage struct {
	Type  string `json:"type"`
	Value struct {
		ParentID int64 `json:"parent_id"`
	} `json:"value"`
}

func (w *Websocket) Recv(timeout time.Duration) (msg IncomingMessage, err error) {
	_ = w.conn.SetReadDeadline(deadlineForTimeout(timeout))
	err = w.conn.ReadJSON(&msg)
	return
}

func deadlineForTimeout(timeout time.Duration) time.Time {
	if timeout == 0 {
		return time.Time{}
	}
	return time.Now().Add(timeout)
}
