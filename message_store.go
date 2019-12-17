package gott

import (
	"sync"
)

type ClientMessage struct {
	Topic, Payload []byte
	QoS            byte
	Client         *Client
	Status         byte
}

//func (cm *ClientMessage) String() string {
//	return fmt.Sprintf("%s:%d", cm.Client.ClientId, cm.QoS)
//}

type MessageStore struct {
	messages map[uint16]*ClientMessage
	mutex    sync.RWMutex
}

func NewMessageStore() *MessageStore {
	return &MessageStore{messages: map[uint16]*ClientMessage{}}
}

func (ms *MessageStore) delete(packetId uint16) {
	ms.mutex.Lock()
	delete(ms.messages, packetId)
	ms.mutex.Unlock()
}

func (ms *MessageStore) Store(packetId uint16, msg *ClientMessage) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.messages[packetId] = msg
}

func (ms *MessageStore) Acknowledge(packetId uint16, status byte, delete bool) {
	if cm := ms.Get(packetId); cm != nil {
		cm.Status = status

		if delete {
			ms.delete(packetId)
		}
	}
}

func (ms *MessageStore) Get(packetId uint16) *ClientMessage {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	if m, ok := ms.messages[packetId]; ok {
		return m
	} else {
		return nil
	}
}
