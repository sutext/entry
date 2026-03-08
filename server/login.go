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
	"sutext.github.io/suid"
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
	token, err := s.createUserToken(user.ID)
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
	token, err := s.createUserToken(user.ID)
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
func (s *server) handleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	ctx := r.Context()
	userID, err := s.ensureLoggedIn(r)
	if err != nil {
		s.writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := s.db.GetUser(ctx, userID)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(user.ToView()); err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (s *server) createUserToken(userID suid.SUID) (string, error) {
	return jwt.Signed(s.signer).Claims(jwt.Claims{
		Subject:  userID.String(),
		Expiry:   jwt.NewNumericDate(time.Now().Add(s.accessTokenDuration)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}).Serialize()
}
func (s *server) getToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("missing authorization token")
	}
	token = token[len("Bearer "):]
	return token, nil
}
func (s *server) ensureLoggedIn(r *http.Request) (uid suid.SUID, err error) {
	token, err := s.getToken(r)
	if err != nil {
		return uid, err
	}
	tok, err := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.EdDSA})
	if err != nil {
		return uid, err
	}
	var claims jwt.Claims
	if err = tok.Claims(s.secret, &claims); err != nil {
		return uid, err
	}
	return suid.Parse(claims.Subject)
}
func (s *server) writeError(w http.ResponseWriter, code int, msg string) {
	http.Error(w, msg, code)
}
