// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/gob"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"google.golang.org/api/plus/v1"

	"gae_wgame/db"
)

const (
	DefaultSessionID = "default"

	// The following keys are used for the default session. For example:
	//  session, _ := db.SessionStore.New(r, defaultSessionID)
	//  session.Values[oauthTokenSessionKey]
	GoogleProfileSessionKey = "google_profile"
	OauthTokenSessionKey    = "oauth_token"

	// This key is used in the OAuth flow session to store the URL to redirect the
	// user to after the OAuth flow is complete.
	OauthFlowRedirectKey = "redirect"
)

func init() {
	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})
	gob.Register(&plus.Person{})
}

// FetchProfile retrieves the Google+ profile of the user associated with the
// provided OAuth token.
func FetchProfile(ctx context.Context, tok *oauth2.Token) (*plus.Person, error) {
	client := oauth2.NewClient(ctx, db.OAuthConfig.TokenSource(ctx, tok))
	plusService, err := plus.New(client)
	if err != nil {
		return nil, err
	}
	return plusService.People.Get("me").Do()
}

// ProfileFromSession retreives the Google+ profile from the default session.
// Returns nil if the profile cannot be retreived (e.g. user is logged out).
func ProfileFromSession(r *http.Request) *plus.Person {
	session, err := db.SessionStore.Get(r, DefaultSessionID)
	if err != nil {
		return nil
	}
	tok, ok := session.Values[OauthTokenSessionKey].(*oauth2.Token)
	if !ok || !tok.Valid() {
		return nil
	}
	profile, ok := session.Values[GoogleProfileSessionKey].(*plus.Person)
	if !ok {
		return nil
	}
	return profile
}

// StripProfile returns a subset of a plus.Person.
func StripProfile(p *plus.Person) *plus.Person {
	return &plus.Person{
		Id:          p.Id,
		DisplayName: p.DisplayName,
		Image:       p.Image,
		Etag:        p.Etag,
		Name:        p.Name,
		Url:         p.Url,
	}
}
