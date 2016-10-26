package sessions

import (
	"sync"
	"github.com/TF2Stadium/wsevent"
	"time"
	"github.com/TF2Stadium/PlayerStatsScraper/steamid"
)

var (
	socketsMu        = new(sync.RWMutex)
	IDSockets   = make(map[string][]*wsevent.Client) //id -> client array, since players can have multiple tabs open
	connectedMu      = new(sync.Mutex)
	connectedTimer   = make(map[string](*time.Timer))
)

func AddSocket(id string, so *wsevent.Client) {
	socketsMu.Lock()
	defer socketsMu.Unlock()

	IDSockets[id] = append(IDSockets[id], so)
	if len(IDSockets[id]) == 1 {
		connectedMu.Lock()
		timer, ok := connectedTimer[id]
		if ok {
			timer.Stop()
			delete(connectedTimer, id)
		}
		connectedMu.Unlock()
	}
}

func RemoveSocket(sessionID, id string) {
	socketsMu.Lock()
	defer socketsMu.Unlock()

	clients := IDSockets[id]
	for i, socket := range clients {
		if socket.ID == sessionID {
			clients[i] = clients[len(clients)-1]
			clients[len(clients)-1] = nil
			clients = clients[:len(clients)-1]
			break
		}
	}

	IDSockets[id] = clients

	if len(clients) == 0 {
		delete(IDSockets, id)
	}
}

func GetSockets(id string) (sockets []*wsevent.Client, success bool) {
	socketsMu.RLock()
	defer socketsMu.RUnlock()

	sockets, success = IDSockets[id]
	return
}

func IsConnected(id string) bool {
	_, ok := GetSockets(id)
	return ok
}

func ConnectedSockets(id string) int {
	socketsMu.RLock()
	l := len(IDSockets[id])
	socketsMu.RUnlock()

	return l
}

func AfterDisconnectedFunc(id string, d time.Duration, f func()) {
	connectedMu.Lock()
	connectedTimer[id] = time.AfterFunc(d, func() {
		if !IsConnected(id) {
			f()
		}

		connectedMu.Lock()
		delete(connectedTimer, id)
		connectedMu.Unlock()
	})
	connectedMu.Unlock()
}