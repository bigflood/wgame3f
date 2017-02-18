package db

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func (db *datastoreDB) GetGame(ctx context.Context, id string) (*Game, error) {
	k := datastore.NewKey(ctx, "Game", id, 0, nil)
	game := &Game{}
	if err := datastore.Get(ctx, k, game); err != nil {
		return nil, fmt.Errorf("datastoredb: could not get Game: %v", err)
	}
	game.ID = id
	return game, nil
}

//
// func (db *datastoreDB) AddGame(ctx context.Context, b *Game) (id int64, err error) {
// 	k := datastore.NewIncompleteKey(ctx, "Game", nil)
// 	k, err = datastore.Put(ctx, k, b)
// 	if err != nil {
// 		return 0, fmt.Errorf("datastoredb: could not put Game: %v", err)
// 	}
// 	return k.IntID(), nil
// }

func (db *datastoreDB) DeleteGame(ctx context.Context, id string) error {
	k := datastore.NewKey(ctx, "Game", id, 0, nil)
	if err := datastore.Delete(ctx, k); err != nil {
		return fmt.Errorf("datastoredb: could not delete Game: %v", err)
	}
	return nil
}

func (db *datastoreDB) UpdateGame(ctx context.Context, b *Game) error {
	k := datastore.NewKey(ctx, "Game", b.ID, 0, nil)
	if _, err := datastore.Put(ctx, k, b); err != nil {
		return fmt.Errorf("datastoredb: could not update Game: %v", err)
	}
	return nil
}

func (db *datastoreDB) ListGames(ctx context.Context) ([]*Game, error) {
	var games []*Game
	q := datastore.NewQuery("Game").
		Order("Title")

	keys, err := q.GetAll(ctx, &games)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list games: %v", err)
	}

	for i, k := range keys {
		games[i].ID = k.StringID()
	}

	return games, nil
}
