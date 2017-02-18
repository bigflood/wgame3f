// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package view

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"google.golang.org/api/plus/v1"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"

	"gae_wgame/auth"
	"gae_wgame/db"
)

var (
	// TemplatesDir var
	TemplatesDir = "templates"
)

// parseTemplate applies a given file to the body of the base template.
func parseTemplate(filename string) *appTemplate {
	tmpl := template.Must(template.ParseFiles(filepath.Join(TemplatesDir, "base.html")))

	// Put the named file into a template called "body"
	path := filepath.Join(TemplatesDir, filename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("could not read template: %v", err))
	}
	template.Must(tmpl.New("body").Parse(string(b)))

	return &appTemplate{tmpl.Lookup("base.html")}
}

// appTemplate is a user login-aware wrapper for a html/template.
type appTemplate struct {
	t *template.Template
}

func (tmpl *appTemplate) Execute(w http.ResponseWriter, r *http.Request, data interface{}) *AppError {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	adminLogin, _ := user.LoginURL(ctx, r.URL.RequestURI())
	adminLogout, _ := user.LogoutURL(ctx, r.URL.RequestURI())

	d := struct {
		Data           interface{}
		AuthEnabled    bool
		Profile        *plus.Person
		AdminUser      *user.User
		IsAdmin        bool
		LoginURL       string
		LogoutURL      string
		AdminLoginURL  string
		AdminLogoutURL string
	}{
		Data:           data,
		AuthEnabled:    db.OAuthConfig != nil,
		LoginURL:       "/login?redirect=" + r.URL.RequestURI(),
		LogoutURL:      "/logout?redirect=" + r.URL.RequestURI(),
		AdminLoginURL:  adminLogin,
		AdminLogoutURL: adminLogout,
	}

	d.AdminUser = u
	d.IsAdmin = user.IsAdmin(ctx)

	if d.AuthEnabled {
		// Ignore any errors.
		d.Profile = auth.ProfileFromSession(r)
	}

	if err := tmpl.t.Execute(w, d); err != nil {
		return AppErrorf(err, "could not write template: %v")
	}
	return nil
}
