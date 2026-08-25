package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/api"
)

type Handles struct {
	//secrets  *Secrets
	//postgres *Postgres
	httpClient *http.Client
	logger     *slog.Logger

	AddAccountTimeout  time.Duration
	FindAccountTimeout time.Duration
	JWTCookieTimeout   time.Duration
	//JWTCookiePath      string
	//JWTCookieSecure    bool
	//JWTCookieHTTPOnly  bool
	//JWTCookieDomain    string
}

func NewHTTPHandles(
	//secrets *Secrets,
	//postgres *Postgres,
	logger *slog.Logger,
	addAccountTimeout, findAccountTimeout, JWTCookieTimeout time.Duration,
	// JWTCookiePath string,
	// JWTCookieSecure bool,
	// JWTCookieHTTPOnly bool,
	// JWTCookieDomain string,
) *Handles {
	return &Handles{
		//secrets,
		//postgres,
		&http.Client{Timeout: 5 * time.Second},
		logger,

		addAccountTimeout,
		findAccountTimeout,
		JWTCookieTimeout,
		//JWTCookiePath,
		//JWTCookieSecure,
		//JWTCookieHTTPOnly,
		//JWTCookieDomain,
	}
}

// Healthy handles /ok requests
func (h *Handles) Healthy(w http.ResponseWriter, _ *http.Request) {
	const op = "main.Healthy"

	err := api.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
	if err != nil {
		h.logger.Error("could not write response", "op", op, "err", err)
		return
	}
}

// AvailableVocabs handles /ok requests
func (h *Handles) AvailableVocabs(w http.ResponseWriter, r *http.Request) {
	const op = "main.AvailableVocabs"

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet,
		"http://127.0.0.1:8079/api/available-vocabs", nil)
	if err != nil {
		h.logger.Error("could not build request", "op", op, "err", err)
		api.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal"})
		return
	}

	res, err := h.httpClient.Do(req)
	if err != nil {
		h.logger.Error("vocabs service unreachable", "op", op, "err", err)
		api.WriteJSON(w, http.StatusBadGateway, map[string]any{"error": "upstream error"})
		return
	}
	defer res.Body.Close()

	var vocabs struct {
		Vocabs []string `json:"vocabs"`
	}
	if err = json.NewDecoder(res.Body).Decode(&vocabs); err != nil {
		h.logger.Error("could not decode response", "op", op, "err", err)
		api.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal"})
		return
	}

	if err = api.WriteJSON(w, http.StatusOK, map[string]any{"vocabs": vocabs.Vocabs}); err != nil {
		h.logger.Error("could not write response", "op", op, "err", err)
	}
}

func (h *Handles) Register(w http.ResponseWriter, r *http.Request) {
	const op = "main.Register"

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		h.logger.Error("could not decode credentials", "op", op, "err", err)
		writeErr := api.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"err": "invalid body",
		})
		if writeErr != nil {
			h.logger.Error("could not write response", "op", op, "err", writeErr)
		}
		return
	}

	if err := creds.ValidateForRegister(); err != nil {
		writeErr := api.WriteJSON(w, http.StatusBadRequest, map[string]any{"err": err})
		if writeErr != nil {
			h.logger.Error("could not write response", "op", op, "err", writeErr)
		}
		return
	}

	hashed := h.secrets.hashSecret(creds.Password)
	creds.Password = hashed

	ctx, cancel := context.WithTimeout(r.Context(), h.AddAccountTimeout)
	err := h.postgres.AddAccount(ctx, creds)
	cancel()

	if err != nil {
		h.logger.Info("could not create user", "err", err)
		writeErr := api.WriteJSON(w, http.StatusBadRequest, map[string]any{"err": err.Error()})
		if writeErr != nil {
			h.logger.Error("could not write response", "err", writeErr)
		}
		return
	}

	exp := time.Now().Add(h.JWTCookieTimeout)
	token, err := h.secrets.NewJWT(creds.Login, exp)
	if err != nil {
		h.logger.Info("could not create jwt token", "err", err)
		writeErr := api.WriteJSON(w, http.StatusBadRequest, map[string]any{"err": err.Error()})
		if writeErr != nil {
			h.logger.Error("could not write response", "err", writeErr)
		}
		return
	}

	//http.SetCookie(w, &http.Cookie{
	//	Name:     middleware.LoginCookieName,
	//	Value:    token,
	//	Path:     h.JWTCookiePath,
	//	MaxAge:   int(h.JWTCookieTimeout.Seconds()),
	//	Secure:   h.JWTCookieSecure,
	//	HttpOnly: h.JWTCookieHTTPOnly,
	//	SameSite: http.SameSiteLaxMode,
	//	Domain:   h.JWTCookieDomain,
	//})

	err = api.WriteJSON(w, http.StatusOK, map[string]any{"token": token, "exp": exp.Unix()})
	if err != nil {
		h.logger.Error("could not write response", "err", err)
	}
}

//
//func (h *Handles) Login(w http.ResponseWriter, r *http.Request) {
//	const op = "main.Login"
//
//	var creds Credentials
//	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
//		h.logger.Error("could not decode credentials", "op", op, "err", err)
//		writeErr := api.WriteJSON(w, http.StatusBadRequest, map[string]any{"err": "invalid body"})
//		if writeErr != nil {
//			h.logger.Error("could not write response", "op", op, "err", writeErr)
//		}
//		return
//	}
//
//	if err := creds.ValidateForLogin(); err != nil {
//		writeErr := api.WriteJSON(w, http.StatusBadRequest, map[string]any{"err": err})
//		if writeErr != nil {
//			h.logger.Error("could not write response", "op", op, "err", writeErr)
//		}
//		return
//	}
//
//	ctx, cancel := context.WithTimeout(r.Context(), h.FindAccountTimeout)
//	hash, found := h.postgres.FindAccount(ctx, creds.Login)
//	cancel()
//	if !found {
//		err := api.WriteJSON(w, http.StatusBadRequest, map[string]any{"err": "account not found"})
//		if err != nil {
//			h.logger.Error("could not write response", "op", op, "err", err)
//		}
//		return
//	}
//	if !h.secrets.VerifyPassword(creds.Password, hash) {
//		err := api.WriteJSON(w, http.StatusBadRequest, map[string]any{"err": "wrong login or password"})
//		if err != nil {
//			h.logger.Error("could not write response", "op", op, "err", err)
//		}
//		return
//	}
//
//	exp := time.Now().Add(h.JWTCookieTimeout)
//	token, err := h.secrets.NewJWT(creds.Login, exp)
//	if err != nil {
//		h.logger.Info("could not create jwt token", "err", err)
//		writeErr := api.WriteJSON(w, http.StatusBadRequest, map[string]any{"err": err.Error()})
//		if writeErr != nil {
//			h.logger.Error("could not write response", "err", writeErr)
//		}
//		return
//	}
//	//http.SetCookie(w, &http.Cookie{
//	//	Name:   middleware.LoginCookieName,
//	//	Value:  token,
//	//	Path:   "/",
//	//	MaxAge: h.JWTCookieTimeout,
//	//	//Secure:   true,
//	//	HttpOnly: true,
//	//	SameSite: http.SameSiteLaxMode,
//	//	//Domain:   "xolra0d.com",
//	//})
//
//	//http.SetCookie(w, &http.Cookie{
//	//	Name:     middleware.LoginCookieName,
//	//	Value:    token,
//	//	Path:     h.JWTCookiePath,
//	//	MaxAge:   int(h.JWTCookieTimeout.Seconds()),
//	//	Secure:   h.JWTCookieSecure,
//	//	HttpOnly: h.JWTCookieHTTPOnly,
//	//	SameSite: http.SameSiteLaxMode,
//	//	Domain:   h.JWTCookieDomain,
//	//})
//
//	err = api.WriteJSON(w, http.StatusOK, map[string]any{"token": token, "exp": exp.Unix()})
//	if err != nil {
//		h.logger.Error("could not write response", "err", err)
//	}
//}
//
//func (h *Handles) PublicKeys(w http.ResponseWriter, _ *http.Request) {
//	pem, err := h.secrets.EncodeJWTPublicKey()
//	if err != nil {
//		writeErr := api.WriteJSON(w, http.StatusInternalServerError, map[string]any{"err": err.Error()})
//		if writeErr != nil {
//			h.logger.Error("could not write response", "err", writeErr)
//		}
//		return
//	}
//	err = api.WriteJSON(w, http.StatusOK, map[string]any{
//		"jwt": pem,
//	})
//	if err != nil {
//		h.logger.Error("could not write response", "err", err)
//	}
//}
