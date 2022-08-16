/*
Copyright © 2022 Aurelio Calegari, et al.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package gcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aurc/loggo/internal/char"
	"github.com/aurc/loggo/internal/util"
)

//const AuthorizationEndpoint = "https://accounts.google.com/o/oauth2/auth/oauthchooseaccount"

const AuthorizationEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
const codeChallengeMethod = "S256"

// Client ID from project "usable-auth-library", configured for
// general purpose API testing, extracted from gcloud sdk sourcecode.
const (
	DefaultCredentialsDefaultClientId     = "764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com"
	DefaultCredentialsDefaultClientSecret = "d-FL95Q19q7MQmFpd7hHD0Ty"
)

func OAuth() {
	doOAuthAsync(DefaultCredentialsDefaultClientId, DefaultCredentialsDefaultClientSecret)
}

func doOAuthAsync(clientId, clientSecret string) {
	// Generates state and PKCE values.
	state, _ := CreateCodeVerifier()
	codeVerifier, _ := CreateCodeVerifier()

	// Creates a redirect URI using an available port on the loopback address.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		util.Log().Fatal(err)
	}
	tcp := listener.Addr().(*net.TCPAddr)

	redirectUri := fmt.Sprintf("http://%s:%d/", "127.0.0.1", tcp.Port)
	util.Log().WithField("code", redirectUri).Info("Preparing redirect url")

	scopes := []string{
		"openid",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/accounts.reauth",
	}

	// Creates the OAuth 2.0 authorization request.
	data := url.Values{}
	data.Set("response_type", "code")
	data.Set("scope", strings.Join(scopes, " "))
	data.Set("redirect_uri", redirectUri)
	data.Set("access_type", "offline")
	data.Set("client_id", clientId)
	data.Set("state", state.String())
	data.Set("flowName", "GeneralOAuthFlow")
	data.Set("code_challenge", codeVerifier.CodeChallengeS256())
	data.Set("code_challenge_method", codeChallengeMethod)

	authorizationRequest := fmt.Sprintf(`%s?%s`, AuthorizationEndpoint, data.Encode())
	c := &callbackHandler{
		state: state.String(),
		code:  make(chan string, 1),
	}

	util.Log().WithField("code", authorizationRequest).Info("Assembled authorisation request.")

	go func() {
		err = util.OpenBrowser(authorizationRequest)
		if err != nil {
			util.Log().Fatal(err)
		}

	}()

	go func() {
		util.Log().Infof("Start redirect listener at port %d", tcp.Port)
		err = http.Serve(listener, c)
		if err != nil {
			util.Log().Fatal(err)
		}
	}()

	code := <-c.code
	util.Log().WithField("code", code).Info("Received Auth code.")
	a := exchangeCodeForTokensAsync(code, codeVerifier.String(), redirectUri, clientId, clientSecret)
	a.Save()
}

func exchangeCodeForTokensAsync(code, codeVerifier, redirectUri, clientId, clientSecret string) *Auth {
	// builds the request
	tokenRequestUri := "https://www.googleapis.com/oauth2/v4/token"

	data := url.Values{}
	data.Set("code", code)
	data.Set("redirect_uri", redirectUri)
	data.Set("client_id", clientId)
	data.Set("code_verifier", codeVerifier)
	data.Set("client_secret", clientSecret)
	data.Set("scope", "")
	data.Set("grant_type", "authorization_code")
	encodedData := data.Encode()

	util.Log().WithField("code", tokenRequestUri).WithField("data", encodedData).Info("Requesting token exchange.")

	req, err := http.NewRequest(http.MethodPost, tokenRequestUri, strings.NewReader(encodedData))
	if err != nil {
		util.Log().Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		util.Log().Fatal(err)
	}
	body := bytes.NewBufferString("")
	_, err = body.ReadFrom(response.Body)
	if err != nil {
		util.Log().Fatal(err)
	}

	util.Log().WithField("code", body.String()).Info("Received token response.")

	m := make(map[string]string)
	_ = json.Unmarshal(body.Bytes(), &m)
	return &Auth{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RefreshToken: m["refresh_token"],
		Type:         "authorized_user",
	}
}

type callbackHandler struct {
	state string
	code  chan string
}

func (c *callbackHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Second)
	vals := req.URL.Query()

	util.Log().WithField("code", req.URL.RawQuery).Info("Received request from browser")

	if strings.Contains(req.RequestURI, "favicon") {
		resp.WriteHeader(404)
		return
	}
	if vals.Has("error") {
		_, _ = resp.Write([]byte(fmt.Sprintf(`<html><body style="text-align: center;font-family: Courier;"><h1>Authentication Failed!<h1><h2>%s</h2><span style="padding: 20px;background-color: #b4b2b2">%s</span></body></html>`,
			"OAuth authorization error:", vals.Get("error"))))
		util.Log().Fatal("OAuth authorization error:", vals.Get("error"))
	}

	// extracts the code
	var code string
	var incomingState string
	if vals.Has("code") && vals.Has("state") {
		code = vals.Get("code")
		incomingState = vals.Get("state")
	} else {
		_, _ = resp.Write([]byte(fmt.Sprintf(`<html><body style="text-align: center;font-family: Courier;"><h1>Authentication Failed!<h1><h2>%s</h2><span style="padding: 20px;background-color: #b4b2b2">%s</span></body></html>`,
			"Malformed authorization response:", req.URL.RawQuery)))
		util.Log().Fatal("Malformed authorization response:", req.URL.RawQuery)
	}

	// Compares the receieved state to the expected value, to ensure that
	// this app made the request which resulted in authorization.
	if c.state != incomingState {
		_, _ = resp.Write([]byte(fmt.Sprintf(`<html><body style="text-align: center;font-family: Courier;"><h1>Authentication Failed!<h1><h2>%s</h2><span style="padding: 20px;background-color: #b4b2b2">%s</span></body></html>`,
			"Received request with invalid state:", incomingState)))
		util.Log().Fatal("Received request with invalid state:", incomingState)
	}
	c.code <- code

	builder := strings.Builder{}
	builder.WriteString(`<html><head><meta http-equiv='refresh' content='10;url=https://cloud.google.com/sdk/auth_success'></head>`)
	builder.WriteString(`<style>.bgCol {background-color: #1a1a1a;color: #2f2e2e;}.fgCol {color: #4bcece;}</style>`)
	builder.WriteString(`<body style="font-family: Courier; background-color: black;color: white;text-align: center"><br><h1>You can close your browser window now!</h1><h3>GCP Authentication Complete</h3><br><h2>`)
	builder.WriteString(char.NewCanvas().WithWord(char.LoggoLogo...).PrintCanvasAsHtml())
	builder.WriteString(`</h2><br><span style="font-weight: bold;color: white">l'oGGo:</span> <span style="color:yellow"><u>Rich Terminal User Interface for following JSON logs</u></span><br><span style="color: gray">Copyright &copy; 2022 Aurelio Calegari, et al.</span><br><a href="https://github.com/aurc/loggo">https://github.com/aurc/loggo</a></body></html>`)
	_, _ = resp.Write([]byte(builder.String()))
}
