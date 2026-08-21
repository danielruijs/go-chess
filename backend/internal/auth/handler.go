package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	userStore    *UserStore
	sessionStore *SessionStore
	cookieDomain string
}

func NewHandler(userStore *UserStore, sessionStore *SessionStore, cookieDomain string) *Handler {
	return &Handler{
		userStore:    userStore,
		sessionStore: sessionStore,
		cookieDomain: cookieDomain,
	}
}

type credentials struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

type PlayerInfoData struct {
	Username        string `json:"username"`
	DisplayName     string `json:"displayName"`
	IsAuthenticated bool   `json:"isAuthenticated"`
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, sessionID SessionID) {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    string(sessionID),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionCookieMaxAge.Seconds()),
		Domain:   h.cookieDomain,
	}
	http.SetCookie(w, cookie)
}

func (h *Handler) clearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Domain:   h.cookieDomain,
	}
	http.SetCookie(w, cookie)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := h.userStore.Register(r.Context(), creds.Username, creds.Password, creds.DisplayName)
	if err != nil {
		if valErr, ok := errors.AsType[*UserRegistrationError](err); ok {
			http.Error(w, valErr.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create session for authenticated user
	session := h.sessionStore.CreateSession(user.Username, user.DisplayName)
	h.setSessionCookie(w, session.ID)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(PlayerInfoData{
		Username:        user.Username,
		DisplayName:     session.DisplayName,
		IsAuthenticated: true,
	})
	if err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := h.userStore.Login(r.Context(), creds.Username, creds.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session := h.sessionStore.CreateSession(user.Username, user.DisplayName)
	h.setSessionCookie(w, session.ID)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(PlayerInfoData{
		Username:        user.Username,
		DisplayName:     session.DisplayName,
		IsAuthenticated: true,
	})
	if err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		h.sessionStore.DeleteSession(SessionID(cookie.Value))
	}

	h.clearSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		if session, exists := h.sessionStore.GetSession(SessionID(cookie.Value)); exists {
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(PlayerInfoData{
				Username:        session.Username,
				DisplayName:     session.DisplayName,
				IsAuthenticated: session.Username != "",
			})
			if err != nil {
				log.Printf("encode response: %v", err)
			}
			return
		}
	}

	// No session or invalid session: create a new anonymous session
	session := h.sessionStore.CreateAnonSession()
	h.setSessionCookie(w, session.ID)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(PlayerInfoData{
		Username:        "",
		DisplayName:     session.DisplayName,
		IsAuthenticated: false,
	})
	if err != nil {
		log.Printf("encode response: %v", err)
	}
}
