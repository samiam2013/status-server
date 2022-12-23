package main

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func basicAuthMiddleware(next http.HandlerFunc, conf config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the auth header or return 401
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Basic ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// split the auth header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// decode the base64
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// split the username and password
		parts = strings.Split(string(decoded), ":")
		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// check the username and password with time invariant compare
		if !constantTimeAuthorize(parts[0], parts[1], conf.Username, conf.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func constantTimeAuthorize(userInput, passInput, user, pass string) bool {
	failAnyway := false
	// pad the username and password to correct length
	if len(userInput) < len(user) {
		userInput += strings.Repeat(" ", len(user)-len(userInput))
	}
	if len(userInput) > len(user) {
		failAnyway = true
		userInput = userInput[:len(user)]
	}
	if len(passInput) < len(pass) {
		passInput += strings.Repeat(" ", len(pass)-len(passInput))
	}
	if len(passInput) > len(pass) {
		failAnyway = true
		passInput = passInput[:len(pass)]
	}
	// check the username and password with time invariant compare
	uAuth := subtle.ConstantTimeCompare([]byte(userInput), []byte(user)) == 1
	pAuth := subtle.ConstantTimeCompare([]byte(passInput), []byte(pass)) == 1
	logrus.Warn("!failAnyway ", !failAnyway, " uAuth ", uAuth, " pAuth ", pAuth)
	auth := longCircuitAnd(!failAnyway, longCircuitAnd(uAuth, pAuth))
	return auth
}

const (
	boolTrue = iota + 1
	boolFalse
)

func convert(b bool) int8 {
	if b {
		return boolTrue
	}
	return boolFalse
}

func longCircuitAnd(p, q bool) bool {
	pInt := convert(p)
	qInt := convert(q)
	sum := pInt & qInt
	return sum == boolTrue
}
