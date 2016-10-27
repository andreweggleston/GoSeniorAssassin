package sessions

import (
	"sync"
	"github.com/TF2Stadium/wsevent"
	"time"
)

var (
	socketsMu        = new(sync.RWMutex)
	IDSockets   = make(map[string][]*wsevent.Client) //id -> client array, since players can have multiple tabs open
	connectedMu      = new(sync.Mutex)
	connectedTimer   = make(map[string](*time.Timer))
)

func AddSocket(studentid string, so *wsevent.Client) {
	socketsMu.Lock()
	defer socketsMu.Unlock()

	IDSockets[studentid] = append(IDSockets[studentid], so)
	if len(IDSockets[studentid]) == 1 {
		connectedMu.Lock()
		timer, ok := connectedTimer[studentid]
		if ok {
			timer.Stop()
			delete(connectedTimer, studentid)
		}
		connectedMu.Unlock()
	}
}

func RemoveSocket(sessionID, studentid string) {
	socketsMu.Lock()
	defer socketsMu.Unlock()

	clients := IDSockets[studentid]
	for i, socket := range clients {
		if socket.ID == sessionID {
			clients[i] = clients[len(clients)-1]
			clients[len(clients)-1] = nil
			clients = clients[:len(clients)-1]
			break
		}
	}

	IDSockets[studentid] = clients

	if len(clients) == 0 {
		delete(IDSockets, studentid)
	}
}

func GetSockets(studentid string) (sockets []*wsevent.Client, success bool) {
	socketsMu.RLock()
	defer socketsMu.RUnlock()

	sockets, success = IDSockets[studentid]
	return
}

func IsConnected(studentid string) bool {
	_, ok := GetSockets(studentid)
	return ok
}

func ConnectedSockets(studentid string) int {
	socketsMu.RLock()
	l := len(IDSockets[studentid])
	socketsMu.RUnlock()

	return l
}

func AfterDisconnectedFunc(studentid string, d time.Duration, f func()) {
	connectedMu.Lock()
	connectedTimer[studentid] = time.AfterFunc(d, func() {
		if !IsConnected(studentid) {
			f()
		}

		connectedMu.Lock()
		delete(connectedTimer, studentid)
		connectedMu.Unlock()
	})
	connectedMu.Unlock()
}