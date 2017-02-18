package view

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"google.golang.org/appengine"

	"gae_wgame/auth"
	"gae_wgame/db"
)

// loginHandler initiates an OAuth flow to authenticate the user.
func loginHandler(w http.ResponseWriter, r *http.Request) *AppError {

	sessionID := uuid.NewV4().String()

	oauthFlowSession, err := db.SessionStore.New(r, sessionID)
	if err != nil {
		return AppErrorf(err, "could not create oauth session: %v", err)
	}
	oauthFlowSession.Options.MaxAge = 10 * 60 // 10 minutes

	redirectURL, err := validateRedirectURL(r.FormValue("redirect"))
	if err != nil {
		return AppErrorf(err, "invalid redirect URL: %v", err)
	}
	oauthFlowSession.Values[auth.OauthFlowRedirectKey] = redirectURL

	if err := oauthFlowSession.Save(r, w); err != nil {
		return AppErrorf(err, "could not save session: %v", err)
	}

	// Use the session ID for the "state" parameter.
	// This protects against CSRF (cross-site request forgery).
	// See https://godoc.org/golang.org/x/oauth2#Config.AuthCodeURL for more detail.
	url := db.OAuthConfig.AuthCodeURL(sessionID, oauth2.ApprovalForce,
		oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

// validateRedirectURL checks that the URL provided is valid.
// If the URL is missing, redirect the user to the application's root.
// The URL must not be absolute (i.e., the URL must refer to a path within this
// application).
func validateRedirectURL(path string) (string, error) {
	if path == "" {
		return "/", nil
	}

	// Ensure redirect URL is valid and not pointing to a different server.
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "/", err
	}
	if parsedURL.IsAbs() {
		return "/", errors.New("URL must be absolute")
	}
	return path, nil
}

// oauthCallbackHandler completes the OAuth flow, retreives the user's profile
// information and stores it in a session.
func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) *AppError {
	ctx := appengine.NewContext(r)

	oauthFlowSession, err := db.SessionStore.Get(r, r.FormValue("state"))
	if err != nil {
		return AppErrorf(err, "invalid state parameter. try logging in again.")
	}

	redirectURL, ok := oauthFlowSession.Values[auth.OauthFlowRedirectKey].(string)
	// Validate this callback request came from the app.
	if !ok {
		return AppErrorf(err, "invalid state parameter. try logging in again.")
	}

	code := r.FormValue("code")
	tok, err := db.OAuthConfig.Exchange(ctx, code)
	if err != nil {
		return AppErrorf(err, "could not get auth token: %v", err)
	}

	session, err := db.SessionStore.New(r, auth.DefaultSessionID)
	if err != nil {
		return AppErrorf(err, "could not get default session: %v", err)
	}

	profile, err := auth.FetchProfile(ctx, tok)
	if err != nil {
		return AppErrorf(err, "could not fetch Google profile: %v", err)
	}

	session.Values[auth.OauthTokenSessionKey] = tok
	// Strip the profile to only the fields we need. Otherwise the struct is too big.
	session.Values[auth.GoogleProfileSessionKey] = auth.StripProfile(profile)
	if err := session.Save(r, w); err != nil {
		return AppErrorf(err, "could not save session: %v", err)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
	return nil
}

// logoutHandler clears the default session.
func logoutHandler(w http.ResponseWriter, r *http.Request) *AppError {
	session, err := db.SessionStore.New(r, auth.DefaultSessionID)
	if err != nil {
		return AppErrorf(err, "could not get default session: %v", err)
	}
	session.Options.MaxAge = -1 // Clear session.
	if err := session.Save(r, w); err != nil {
		return AppErrorf(err, "could not save session: %v", err)
	}
	redirectURL := r.FormValue("redirect")
	if redirectURL == "" {
		redirectURL = "/"
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
	return nil
}
