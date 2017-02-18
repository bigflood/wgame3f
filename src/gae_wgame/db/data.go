package db

import (
	"golang.org/x/net/context"
)

type GameDatabase interface {
	ListGames(ctx context.Context) ([]*Game, error)

	GetGame(ctx context.Context, id string) (*Game, error)
	UpdateGame(ctx context.Context, b *Game) error
	DeleteGame(ctx context.Context, id string) error

	ListLevels(ctx context.Context) ([]*LevelInfo, error)
	ListLevelsWithGame(ctx context.Context, game string) ([]*LevelInfo, error)
	ListLevelsWithGameStage(ctx context.Context, gamestage string) ([]*LevelInfo, error)
	GetLevel(ctx context.Context, id string) (*LevelInfo, error)
	UpdateLevel(ctx context.Context, b *LevelInfo) error
	DeleteLevel(ctx context.Context, id string) error

	GetLevelData(ctx context.Context, id string) (*LevelData, error)
	UpdateLevelData(ctx context.Context, b *LevelData) error
	DeleteLevelData(ctx context.Context, id string) error

	ListIapSettings(ctx context.Context) ([]*IapSettingInfo, error)
	GetIapSetting(ctx context.Context, id string) (*IapSettingInfo, error)
	UpdateIapSetting(ctx context.Context, b *IapSettingInfo) error
	DeleteIapSetting(ctx context.Context, id string) error

	GetCount(ctx context.Context, id string) (int, error)
	AddCount(ctx context.Context, id string, v int) (int, error)

	// Close closes the database, freeing up any available resources.
	// TODO(cbro): Close() should return an error.
	Close()
}
