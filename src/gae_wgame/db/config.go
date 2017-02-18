// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package db

import (
	//	"errors"
	"log"
	"os"

	"github.com/gorilla/sessions"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	DB GameDatabase

	// OAuthConfig oauth2 config
	OAuthConfig *oauth2.Config

	// SessionStore var
	SessionStore sessions.Store

	// Force import of mgo library.
	_ mgo.Session
)

const sessionSecret = "--session-secret--"

// Setup setup
func Setup() {
	var err error

	// [START cloudsql]
	// To use MySQL, uncomment the following lines, and update the username,
	// password and host.
	//
	// DB, err = newMySQLDB(MySQLConfig{
	// 	Username: "",
	// 	Password: "",
	// 	Host:     "",
	// 	Port:     3306,
	// })
	// [END cloudsql]

	// [START mongo]
	// To use Mongo, uncomment the next lines and update the address string and
	// optionally, the credentials.
	//
	// var cred *mgo.Credential
	// DB, err = newMongoDB("localhost", cred)
	// [END mongo]

	// [START datastore]
	// To use Cloud Datastore, uncomment the following lines and update the
	// project ID.
	// More options can be set, see the google package docs for details:
	// http://godoc.org/golang.org/x/oauth2/google
	//
	DB, err = configureDatastoreDB()
	// [END datastore]

	if err != nil {
		log.Fatal(err)
	}

	// [START storage]
	// To configure Cloud Storage, uncomment the following lines and update the
	// bucket name.
	//
	// StorageBucketName = "<your-storage-bucket>"
	// StorageBucket, err = configureStorage(StorageBucketName)
	// [END storage]

	if err != nil {
		log.Fatal(err)
	}

	// [START auth]
	// To enable user sign-in, uncomment the following lines and update the
	// Client ID and Client Secret.
	// You will also need to update OAUTH2_CALLBACK in app.yaml when pushing to
	// production.
	//
	OAuthConfig = configureOAuthClient("--Client ID--", "--Client Secret--")
	// [END auth]

	// [START sessions]
	// Configure storage method for session-wide information.
	// Update "something-very-secret" with a hard to guess string or byte sequence.
	cookieStore := sessions.NewCookieStore([]byte(sessionSecret))
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
	}
	SessionStore = cookieStore
	// [END sessions]

	// [START pubsub]
	// To configure Pub/Sub, uncomment the following lines and update the project ID.
	//
	// PubsubClient, err = configurePubsub("<your-project-id>")
	// [END pubsub]

	if err != nil {
		log.Fatal(err)
	}
}

func configureDatastoreDB() (GameDatabase, error) {
	return newDatastoreDB()
}

func configureOAuthClient(clientID, clientSecret string) *oauth2.Config {
	redirectURL := os.Getenv("OAUTH2_CALLBACK")
	if os.Getenv("RUN_WITH_DEVAPPSERVER") == "1" || redirectURL == "" {
		redirectURL = "http://localhost:8080/oauth2callback"
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}
