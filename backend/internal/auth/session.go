package auth

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	SessionCookieName = "session_id"
	SessionDuration   = 24 * time.Hour
)

type SessionID string

type Session struct {
	ID          SessionID
	Username    string // Empty if anonymous
	DisplayName string
	ExpiresAt   time.Time
}

func (s Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

type SessionStore struct {
	mu          sync.RWMutex
	sessions    map[SessionID]Session
	anonCounter atomic.Uint64
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[SessionID]Session),
	}
}

func (s *SessionStore) nextAnonName() string {
	val := s.anonCounter.Add(1)
	return fmt.Sprintf("Anonymous %d", val)
}

func GenerateSessionID() SessionID {
	return SessionID(uuid.New().String())
}

func IsValidSessionID(sessionID SessionID) bool {
	_, err := uuid.Parse(string(sessionID))
	return err == nil
}

func (s *SessionStore) CreateAnonSession() Session {
	return s.CreateAnonSessionWithID(GenerateSessionID())
}

func (s *SessionStore) CreateAnonSessionWithID(sessionID SessionID) Session {
	return s.CreateSessionWithID(sessionID, "", s.nextAnonName())
}

func (s *SessionStore) CreateSession(username, displayName string) Session {
	return s.CreateSessionWithID(GenerateSessionID(), username, displayName)
}

func (s *SessionStore) CreateSessionWithID(sessionID SessionID, username, displayName string) Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := Session{
		ID:          sessionID,
		Username:    username,
		DisplayName: displayName,
		ExpiresAt:   time.Now().Add(SessionDuration),
	}
	s.sessions[sessionID] = session
	return session
}

// returns the session and a boolean indicating if it was found and valid
func (s *SessionStore) GetSession(sessionID SessionID) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists || session.IsExpired() {
		return Session{}, false
	}
	return session, true
}

func (s *SessionStore) DeleteSession(sessionID SessionID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}

// PlayerKey represents a unique identifier for a player's cached session.
// It is derived from a user's username (for authenticated users) or their session ID (for anonymous users).
type PlayerKey string

func (s Session) PlayerKey() PlayerKey {
	if s.Username != "" {
		return PlayerKey("user:" + s.Username)
	}
	return PlayerKey("anon:" + string(s.ID))
}
