package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"go-chess/internal/cache"

	"github.com/google/uuid"
)

const (
	sessionCookieName = "session_id"
	// sessionCookieMaxAge is intentionally longer than sessionDuration.
	// The backend is the source of truth for session expiry — sessions are evicted
	// after sessionDuration of inactivity. The cookie just needs to persist long
	// enough to outlive any active session, including across browser restarts.
	sessionCookieMaxAge = 30 * 24 * time.Hour

	sessionDuration        = 24 * time.Hour
	sessionCleanupInterval = 10 * time.Minute
)

type SessionID string

type Session struct {
	ID          SessionID
	Username    string // Empty if anonymous
	DisplayName string
}

type SessionStore struct {
	cache       *cache.Cache[SessionID, Session]
	anonCounter atomic.Uint64
}

func NewSessionStore(ctx context.Context) (*SessionStore, error) {
	cache, err := cache.New[SessionID](cache.Options[Session]{
		Cleanup: &cache.CleanupConfig[Session]{
			Interval: sessionCleanupInterval,
			ShouldEvict: func(s Session, lastUsed time.Time) bool {
				return time.Since(lastUsed) > sessionDuration
			},
		},
	})
	if err != nil {
		return nil, err
	}

	store := &SessionStore{
		cache: cache,
	}
	store.cache.StartCleanup(ctx)
	return store, nil
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
	session := Session{
		ID:          sessionID,
		Username:    username,
		DisplayName: displayName,
	}
	s.cache.Set(sessionID, session)
	return session
}

// returns the session and a boolean indicating if it was found and valid
func (s *SessionStore) GetSession(sessionID SessionID) (Session, bool) {
	return s.cache.Get(sessionID)
}

func (s *SessionStore) DeleteSession(sessionID SessionID) {
	s.cache.Delete(sessionID)
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

// SessionIDFromRequest extracts and validates the session ID from the request cookie, if present.
// Returns the session ID and a boolean indicating if it was found and valid.
func SessionIDFromRequest(r *http.Request) (SessionID, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return "", false
	}
	sessionID := SessionID(cookie.Value)
	if !IsValidSessionID(sessionID) {
		return "", false
	}
	return sessionID, true
}
