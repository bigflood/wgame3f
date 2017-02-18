package db

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func (db *datastoreDB) GetIapSetting(ctx context.Context, id string) (*IapSettingInfo, error) {
	k := datastore.NewKey(ctx, "IapSettingInfo", id, 0, nil)
	info := &IapSettingInfo{}
	if err := datastore.Get(ctx, k, info); err != nil {
		return nil, fmt.Errorf("datastoredb: could not get iap setting: %v", err)
	}
	info.ID = id
	return info, nil
}

func (db *datastoreDB) DeleteIapSetting(ctx context.Context, id string) error {
	k := datastore.NewKey(ctx, "IapSettingInfo", id, 0, nil)
	if err := datastore.Delete(ctx, k); err != nil {
		return fmt.Errorf("datastoredb: could not delete iap setting: %v", err)
	}
	return nil
}

func (db *datastoreDB) UpdateIapSetting(ctx context.Context, b *IapSettingInfo) error {
	k := datastore.NewKey(ctx, "IapSettingInfo", b.ID, 0, nil)
	if _, err := datastore.Put(ctx, k, b); err != nil {
		return fmt.Errorf("datastoredb: could not update iap setting: %v", err)
	}
	return nil
}

func (db *datastoreDB) ListIapSettings(ctx context.Context) ([]*IapSettingInfo, error) {
	var infos []*IapSettingInfo
	q := datastore.NewQuery("IapSettingInfo").
		Order("ID")

	keys, err := q.GetAll(ctx, &infos)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list iap settings: %v", err)
	}

	for i, k := range keys {
		infos[i].ID = k.StringID()
	}

	return infos, nil
}
