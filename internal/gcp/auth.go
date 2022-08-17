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
	"context"
	"encoding/json"
	"os"
	"path"

	logging "cloud.google.com/go/logging/apiv2"
	"google.golang.org/api/option"
)

var IsGCloud = false

type Auth struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"type"`
}

func LoggingClient(ctx context.Context) (*logging.Client, error) {
	if !IsGCloud {
		return logging.NewClient(ctx, option.WithCredentialsFile(authFile()))
	} else {
		return logging.NewClient(ctx)
	}
}

func authDir() string {
	hd, _ := os.UserHomeDir()
	dir := path.Join(hd, ".loggo", "auth")
	return dir
}

func authFile() string {
	return path.Join(authDir(), "gcp.json")
}

func Delete() {
	_ = os.Remove(authFile())
}

func (a *Auth) Save() error {
	if err := os.MkdirAll(authDir(), os.ModePerm); err != nil {
		return err
	}
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(authFile())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(b)
	return err
}
