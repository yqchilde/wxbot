package robot

import (
	"github.com/eatmoreapple/openwechat"
)

type Message struct {
	*openwechat.Message
}

type User struct {
	*openwechat.User
}

func (m *Message) Sender() (*User, error) {
	sender, err := m.Message.Sender()
	return &User{sender}, err
}

func (m *Message) SenderInGroup() (*User, error) {
	group, err := m.Message.SenderInGroup()
	return &User{group}, err
}

func (m *Message) Receiver() (*User, error) {
	receiver, err := m.Message.Receiver()
	return &User{receiver}, err
}
