//
// Copyright 2018 Malin Yamato Lääkkö --  All rights reserved.
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
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"flag"
	"fmt"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/google"
	"github.com/dghubble/sessions"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type Date struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
}

const (
	ERROR   = "ERROR"
	WARNING = "WARNING"
	SUCCESS = "SUCCESS"
	GREEN   = "GREEN" // sender and target are sending pvt messages to each other
	BLUE    = "BLUE"  // sender sends pvt messages to the target but not the other way around
	//BLACK = "BLACK" // The target is blocking, black listening the sender
)

type Status struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/session") {
		log.Println("Main: Set path ", r.URL.Path)
		r.URL.Path = "/"
	} else if strings.Contains(r.URL.Path, "/user") {
		//log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(_documentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	} else if strings.Contains(r.URL.Path, "/test") {
		//log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(_documentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	} else if strings.Contains(r.URL.Path, "/css") {
		//log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(_documentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	} else if strings.Contains(r.URL.Path, "/js") {
		//log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(_documentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	} else if strings.Contains(r.URL.Path, "/images") {
		//log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(_documentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "Main: Illegal path "+r.URL.Path, 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Main: Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template.Must(template.ParseFiles(_home)).Execute(w, struct {
		Protocol string
		Host     string
		Port     string
	}{
		Protocol: _config.Protocol,
		Host:     _config.Host,
		Port:     _config.Port,
	})
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {

	sess, err := _sessionStore.Get(r, sessionName)
	if err != nil {
		log.Println("Main: sessionHandler: Error in getting and verifying coookie ", err)
	}
	token := sess.Values[sessionToken].(string)
	log.Println("session token from cookie ", token)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template.Must(template.ParseFiles(_home)).Execute(w, struct {
		Protocol  string
		Host      string
		Port      string
		LoggedIn  string
		LoggedOut string
		Name      string
		Email     string
		Coupons   []Coupon
	}{
		Protocol:  _config.Protocol,
		Host:      _config.Host,
		Port:      _config.Port,
		LoggedIn:  "flex",
		LoggedOut: "none",
		Name:      _admin.FirstName + " " + _admin.LastName,
		Email:     _admin.Email,
		Coupons:   _coupons.getAll(),
	})
}

func NewMux(config *Config) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", serveHome)
	mux.Handle("/session/", requireLogin(http.HandlerFunc(sessionHandler)))
	mux.Handle("/coupon", requireLogin(http.HandlerFunc(CouponHandler)))
	mux.HandleFunc("/logout", logoutHandler)

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.url() + "/google/callback",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	mux.Handle("/google/login", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
	mux.Handle("/google/callback", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueSession(), nil)))

	return mux
}

func getCookieAndTokenfromRequest(r *http.Request, onlyTooken bool) (token string, cookie string, err error) {
	if !onlyTooken {
		//retrieve encrypted cookie
		cookieInfo, err := r.Cookie(sessionName)
		if err != nil {
			return "", "", fmt.Errorf("No cookie found for give cookie name %s detail %s", sessionName, err)
		}
		cookie = cookieInfo.Value
	}
	session, err := _sessionStore.Get(r, sessionName)
	if err != nil {
		return "", "", fmt.Errorf("Fail to retrieve cookie to create session %s detail %s", sessionName, err)
	}
	atoken, ok := session.Values[sessionToken]
	if !ok {
		return "", "", fmt.Errorf("The sesstion did not contain %s ", sessionToken)
	}
	if atoken != nil {
		token = atoken.(string)
	} else {
		token = ""
	}
	return token, cookie, nil
}

const _home = "main.html"

var _coupons Coupons
var _documentRoot string
var _sessionStore *sessions.CookieStore
var _config Config
var _admin Person

func main() {
	_coupons = Coupons{make(map[string]Coupon)}
	if os.Getenv("RakuRunMode") == "Test" {
		_config.load("qr_test.conf")
	} else {
		_config.load("qr.conf")
	}
	if _config.ClientID == "" {
		log.Fatal("Missing Google Client ID")
	}
	if _config.ClientSecret == "" {
		log.Fatal("Missing Google Client Secret")
	}
	_sessionStore = sessions.NewCookieStore([]byte(_config.ChatPrivateKey), nil)
	dir, _ := os.Getwd()
	_documentRoot = strings.Replace(dir, " ", "\\ ", -1)
	var addr = flag.String("addr", ":"+_config.Port, "http service address")
	flag.Parse()
	log.Println("Create a hub and run it in a different thread")
	log.Println("Load persons database...")
	_coupons.load()
	log.Println("Starting service at ", _config.url())
	if _config.Protocol == "http" {
		err := http.ListenAndServe(*addr, NewMux(&_config))
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else { // https
		err := http.ListenAndServeTLS(*addr, _config.SSLCert, _config.SSLPrivateKey, NewMux(&_config))
		if err != nil {
			log.Fatal("ListenAndServe TLS: ", err)
		}
	}
}
