package session

import (
	"sync"
)

type SessionData struct {
	Username string
	Name     string
	Role     string
}

var (
	sessions = make(map[string]SessionData) // session_id -> username
	mu       sync.Mutex                     // para manejar concurrencia
)

// Guarda una sesión
func Set(sessionID, username, name, role string) {
	mu.Lock()
	defer mu.Unlock()
	sessions[sessionID] = SessionData{Username: username, Name: name, Role: role}
}

// Obtiene el usuario por session_id
func Get(sessionID string) (SessionData, bool) {
	mu.Lock()
	defer mu.Unlock()
	data, ok := sessions[sessionID]
	return data, ok
}

// Elimina una sesión
func Delete(sessionID string) {
	mu.Lock()
	defer mu.Unlock()
	delete(sessions, sessionID)
}
