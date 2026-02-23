package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"golang.org/x/crypto/bcrypt"
	"sutext.github.io/entry/model"
)

type loginRequest struct {
	Email    string
	Password string
}
type loginResponse struct {
	User  *model.UserView `json:"user"`
	Token string          `json:"token"`
}

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	ctx := r.Context()
	var req loginRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(req.Password)); err != nil {
		s.writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	token, err := s.generateToken(user)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := loginResponse{
		User:  user.ToView(),
		Token: token,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

type registerRequest struct {
	Email    string
	Password string
}

func (s *server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req registerRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Email == "" || req.Password == "" {
		s.writeError(w, http.StatusBadRequest, "email or password is empty")
		return
	}
	user := model.NewUser()
	user.Email = &req.Email
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user.Hash = string(hash)
	if err = s.db.CreateUser(r.Context(), user); err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := s.generateToken(user)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := loginResponse{
		User:  user.ToView(),
		Token: token,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}
func (s *server) generateToken(user *model.User) (string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.EdDSA,
		Key:       s.secret,
	}, nil)
	if err != nil {
		return "", err
	}
	builder := jwt.Signed(signer)
	builder.Claims(jwt.Claims{
		Subject:  user.ID.String(),
		Expiry:   jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	})
	return builder.Serialize()
}
func (s *server) ensureLoggedIn(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("missing authorization token")
	}
	token = token[len("Bearer "):]
	tok, err := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.EdDSA})
	if err != nil {
		return "", err
	}
	var claims jwt.Claims
	if err = tok.Claims(s.secret.Public(), &claims); err != nil {
		return "", err
	}
	return claims.Subject, nil
}
func (s *server) writeError(w http.ResponseWriter, code int, msg string) {
	http.Error(w, msg, code)
}
