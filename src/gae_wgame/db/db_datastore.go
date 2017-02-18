package db

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// https://cloud.google.com/datastore/docs/concepts/overview
type datastoreDB struct {
}

var _ GameDatabase = &datastoreDB{}

// https://godoc.org/google.golang.org/cloud/datastore
func newDatastoreDB() (GameDatabase, error) {

	return &datastoreDB{}, nil
}

// Close closes the database.
func (db *datastoreDB) Close() {
	// No op.
}

func (db *datastoreDB) GetCount(ctx context.Context, id string) (int, error) {
	counter := &Counter{}
	k := datastore.NewKey(ctx, "Counter", id, 0, nil)
	if err := datastore.Get(ctx, k, counter); err != nil {
		return 0, err
	}
	/*
		if err != nil {
			datastore.RunInTransaction(ctx, func(ctx context.Context) error {
				err := datastore.Get(ctx, k, status)

				if err != nil && err != datastore.ErrNoSuchEntity {
					_, err = datastore.Put(ctx, k, status)
				}
				return err
			}, nil)
		}
	*/

	return counter.Count, nil
}

func (db *datastoreDB) AddCount(ctx context.Context, id string, v int) (int, error) {
	counter := &Counter{}
	k := datastore.NewKey(ctx, "Counter", id, 0, nil)

	datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		err := datastore.Get(ctx, k, counter)

		if err == nil {
			counter.Count += v
			_, err = datastore.Put(ctx, k, counter)
		} else if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		} else {
			counter.Count = v
			_, err = datastore.Put(ctx, k, counter)
		}

		return err
	}, nil)

	return counter.Count, nil
}
