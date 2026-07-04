package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"time"
)

type AuthHandler struct {
	userStore      *UserStore
	sessionStore   *SessionStore
	allowedOrigins []string
	cookieDomain   string
}

func NewAuthHandler(userStore *UserStore, sessionStore *SessionStore, allowedOrigins []string, cookieDomain string) *AuthHandler {
	return &AuthHandler{
		userStore:      userStore,
		sessionStore:   sessionStore,
		allowedOrigins: allowedOrigins,
		cookieDomain:   cookieDomain,
	}
}

func (h *AuthHandler) CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			allowed := slices.Contains(h.allowedOrigins, origin)
			if !allowed {
				http.Error(w, "Forbidden: Invalid CORS Origin", http.StatusForbidden)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
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

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, sessionID SessionID) {
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

func (h *AuthHandler) clearSessionCookie(w http.ResponseWriter) {
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

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
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
