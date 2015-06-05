package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Client struct {
	socket       net.Conn
	br           *bufio.Reader
	sentMessages []message
	// Max messages per 15 seconds
	messageRate int
	server      string
}

type message struct {
	message string
	t       time.Time
}

type Event struct {
	From    string
	Kind    string
	Message string
}

func New(server string, messageRate int) (*Client, error) {
	client := &Client{
		messageRate: messageRate,
		server:      server,
	}

	client.Reload()

	return client, nil
}

func (C *Client) Reload() error {
	conn, err := net.Dial("tcp", C.server)
	br := bufio.NewReaderSize(conn, 512)
	C.br = br
	C.socket = conn

	return err
}

func (C *Client) Send(msg string) {
	C.socket.Write([]byte(msg + "\r\n"))
	log.Printf("< %s", msg)
}

func (C *Client) Say(channel, msg string) {
	// Prune message list
	i := 0
	for i = range C.sentMessages {
		if time.Since(C.sentMessages[i].t) < 15*time.Second {
			break
		}
	}
	if i == len(C.sentMessages) {
		C.sentMessages = []message{}
	} else {
		C.sentMessages = C.sentMessages[i:]
	}

	if len(C.sentMessages) < C.messageRate {
		C.Send(fmt.Sprintf("PRIVMSG %s :%s", channel, msg))
		C.sentMessages = append(C.sentMessages, message{
			message: msg,
			t:       time.Now(),
		})
	}
}

func (C *Client) ReadEvent() (Event, error) {
	t := time.Now()
	t = t.Add(1 * time.Minute)

	C.socket.SetReadDeadline(t)

	msg, err := C.br.ReadString('\n')
	if err != nil {
		return Event{}, err
	}

	space := strings.Index(msg, " ")
	from := msg[:space]

	msg = msg[space+1:]
	space = strings.Index(msg, " ")
	if space == -1 {
		return Event{from, msg, ""}, nil
	}
	kind := msg[:space]

	message := strings.TrimSpace(msg[space:])

	return Event{from, kind, message}, nil
}
