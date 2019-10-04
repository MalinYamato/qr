//
// Copyright 2017 Malin Yamato --  All rights reserved.
// https://github.com/MalinYamato
//
// MIT License
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Rakuen. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	//"github.com/dghubble/gologin/v2/google"

	"github.com/dghubble/gologin/google"
	"github.com/rs/xid"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	sessionName  = "secure.raku.cloud"
	sessionToken = "SessionToken"
)

// Config configures the main ServeMux.

func checkSet(a string, b string) string {
	if a == "" {
		return b
	}
	return a
}

func timestamp() string {
	date := time.Now()
	return strconv.Itoa(date.Day()) + ":" + strconv.Itoa(date.Hour()) + ":" + strconv.Itoa(date.Minute()) + ":" + strconv.Itoa(date.Second())
}

func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		var sessionUserKey = "GoogleUser"
		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		secret := xid.New() // used as a secret to verify identity of users who sends websocket messages from the brwoser to the server
		// remove possible old cookies
		if isAuthenticated(req) {
			log.Println("Login: There was an old cookie. Removing it")
			_sessionStore.Destroy(w, sessionName)
		}
		session := _sessionStore.New(sessionName)
		session.Values[sessionUserKey] = googleUser.Id
		session.Values[sessionToken] = secret.String()
		err = session.Save(w)
		if err != nil {
			log.Println("Login: could not set session ", err)
		}
		_admin.Email = googleUser.Email
		_admin.UserID, _ = strconv.Atoi(googleUser.Id)
		_admin.LastName = googleUser.FamilyName
		_admin.FirstName = googleUser.GivenName

		log.Printf("Login: Successful Login of %s %s Email: %s ID: %$ ",
			googleUser.GivenName, googleUser.FamilyName, googleUser.Email, googleUser.Id)
		http.Redirect(w, req, "/session", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("logoiutHandler called")
	if req.Method == "POST" {
		req.ParseForm()
		_sessionStore.Destroy(w, sessionName)
		log.Println("Session destroyed!")
	}
	// redirect does not work for AJAX calls. Redirects have to be implemtend by client
	w.Write([]byte("login " + SUCCESS))
	http.Redirect(w, req, "/", http.StatusFound)
}

func requireLogin(next http.Handler) http.Handler {
	log.Println("RequireLogin called")
	fn := func(w http.ResponseWriter, req *http.Request) {
		if !isAuthenticated(req) {
			log.Println("requireLogin: login is required as cookie not found")
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func isAuthenticated(req *http.Request) bool {
	_, err := _sessionStore.Get(req, sessionName)
	if err == nil {
		return true
	}
	return false
}
