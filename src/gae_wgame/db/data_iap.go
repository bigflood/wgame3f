package db

import (
	"time"
)

type IapSettingInfo struct {
	ID string

	Title       string    `datastore:",noindex"`
	Description string    `datastore:",noindex"`
	UpdateTime  time.Time `datastore:",noindex"`

	ApplePassword      string    `datastore:",noindex"`
	AmazonSecret       string    `datastore:",noindex"`
	GooglePublicKey    string    `datastore:",noindex"`
	GoogleAccToken     string    `datastore:",noindex"`
	GoogleRefToken     string    `datastore:",noindex"`
	GoogleClientID     string    `datastore:",noindex"`
	GoogleClientSecret string    `datastore:",noindex"`
	GoogleRefreshTime  time.Time `datastore:",noindex"`
}
