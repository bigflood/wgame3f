package db

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func (db *datastoreDB) GetLevel(ctx context.Context, id string) (*LevelInfo, error) {
	k := datastore.NewKey(ctx, "LevelInfo", id, 0, nil)
	level := &LevelInfo{}
	if err := datastore.Get(ctx, k, level); err != nil {
		return nil, fmt.Errorf("datastoredb: could not get Level: %v", err)
	}
	level.ID = id
	return level, nil
}

//
// // AddLevel saves a given level, assigning it a new ID.
// func (db *datastoreDB) AddLevel(ctx context.Context, b *LevelInfo) (id int64, err error) {
// 	k := datastore.NewIncompleteKey(ctx, "LevelInfo", nil)
// 	k, err = datastore.Put(ctx, k, b)
// 	if err != nil {
// 		return 0, fmt.Errorf("datastoredb: could not put Level: %v", err)
// 	}
// 	return k.IntID(), nil
// }

func (db *datastoreDB) DeleteLevel(ctx context.Context, id string) error {
	k := datastore.NewKey(ctx, "LevelInfo", id, 0, nil)
	if err := datastore.Delete(ctx, k); err != nil {
		return fmt.Errorf("datastoredb: could not delete Level: %v", err)
	}
	return nil
}

func (db *datastoreDB) UpdateLevel(ctx context.Context, b *LevelInfo) error {
	k := datastore.NewKey(ctx, "LevelInfo", b.ID, 0, nil)
	if _, err := datastore.Put(ctx, k, b); err != nil {
		return fmt.Errorf("datastoredb: could not update Level: %v", err)
	}
	return nil
}

func (db *datastoreDB) ListLevels(ctx context.Context) ([]*LevelInfo, error) {
	var levels []*LevelInfo
	q := datastore.NewQuery("LevelInfo")

	keys, err := q.GetAll(ctx, &levels)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list levels: %v", err)
	}

	for i, k := range keys {
		levels[i].ID = k.StringID()
	}

	return levels, nil
}

func (db *datastoreDB) ListLevelsWithGame(ctx context.Context, game string) ([]*LevelInfo, error) {
	if game == "" {
		return nil, fmt.Errorf("game is none")
	}

	var levels []*LevelInfo
	q := datastore.NewQuery("LevelInfo").
		Filter("Game =", game)

	keys, err := q.GetAll(ctx, &levels)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list levels: %v", err)
	}

	for i, k := range keys {
		levels[i].ID = k.StringID()
	}

	return levels, nil
}

func (db *datastoreDB) ListLevelsWithGameStage(ctx context.Context, gamestage string) ([]*LevelInfo, error) {
	if gamestage == "" {
		return nil, fmt.Errorf("gamestage is none")
	}

	var levels []*LevelInfo
	q := datastore.NewQuery("LevelInfo").
		Filter("GameStage =", gamestage)

	keys, err := q.GetAll(ctx, &levels)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list levels: %v", err)
	}

	for i, k := range keys {
		levels[i].ID = k.StringID()
	}

	return levels, nil
}

func (db *datastoreDB) GetLevelData(ctx context.Context, id string) (*LevelData, error) {
	k := datastore.NewKey(ctx, "LevelData", id, 0, nil)
	data := &LevelData{}
	if err := datastore.Get(ctx, k, data); err != nil {
		return nil, fmt.Errorf("datastoredb: could not get LevelData: %v", err)
	}
	data.ID = id
	return data, nil
}

func (db *datastoreDB) UpdateLevelData(ctx context.Context, b *LevelData) error {
	k := datastore.NewKey(ctx, "LevelData", b.ID, 0, nil)
	if _, err := datastore.Put(ctx, k, b); err != nil {
		return fmt.Errorf("datastoredb: could not update LevelData: %v", err)
	}
	return nil
}

func (db *datastoreDB) DeleteLevelData(ctx context.Context, id string) error {
	k := datastore.NewKey(ctx, "LevelData", id, 0, nil)
	if err := datastore.Delete(ctx, k); err != nil {
		return fmt.Errorf("datastoredb: could not delete LevelData: %v", err)
	}
	return nil
}
