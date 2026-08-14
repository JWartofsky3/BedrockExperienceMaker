package operations

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/auth"
	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/packs"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request loginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			auth.WriteJSONError(w, "request body is invalid", http.StatusBadRequest)
			return
		}
		request.Username = strings.TrimSpace(request.Username)
		if request.Username == "" || request.Password == "" {
			auth.WriteJSONError(w, "username and password are required", http.StatusBadRequest)
			return
		}

		var user auth.User
		var passwordHash string
		err := db.QueryRowContext(r.Context(), `SELECT id, username, password_hash FROM users WHERE username = ?`, request.Username).Scan(&user.ID, &user.Username, &passwordHash)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(request.Password)) != nil {
			auth.WriteJSONError(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		user.Name = "users/" + user.ID

		token, err := newToken()
		if err != nil {
			auth.WriteJSONError(w, "could not sign in", http.StatusInternalServerError)
			return
		}
		sessionID, err := packs.NewID()
		if err != nil {
			auth.WriteJSONError(w, "could not sign in", http.StatusInternalServerError)
			return
		}
		if _, err := db.ExecContext(r.Context(), `INSERT INTO sessions (id, user_id, token_hash, expires_at) VALUES (?, ?, SHA2(?, 256), DATE_ADD(NOW(), INTERVAL 7 DAY))`, sessionID, user.ID, token); err != nil {
			auth.WriteJSONError(w, "could not sign in", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: auth.CookieName, Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: int((7 * 24 * time.Hour).Seconds())})
		auth.WriteJSON(w, user, http.StatusOK)
	}
}

func newToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
