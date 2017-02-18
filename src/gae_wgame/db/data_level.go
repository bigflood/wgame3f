package db

import (
	"sort"
	"time"
	//"google.golang.org/appengine/datastore"
)

type LevelInfo struct {
	ID string // "Game-Stage-Level"

	Game      string // "Game"
	GameStage string // "Game-Stage"

	Title       string    `datastore:",noindex"`
	Description string    `datastore:",noindex"`
	UpdateTime  time.Time `datastore:",noindex"`
}

type LevelData struct {
	ID   string // "GameStage-Level"
	Data string `datastore:",noindex"`
}

type LevelInfoAndData struct {
	LevelInfo
	Data string
}

func SortLevels(list []*LevelInfo) {
	sort.Sort(levelListSorter{List: list})
}

type levelListSorter struct {
	List []*LevelInfo
}

func (this levelListSorter) Len() int {
	return len(this.List)
}

func (this levelListSorter) Less(i, j int) bool {
	return stringLess(this.List[i].ID, this.List[j].ID)
}

func (this levelListSorter) Swap(i, j int) {
	this.List[i], this.List[j] = this.List[j], this.List[i]
}

func stringLess(a, b string) bool {
	i := 0
	n := min(len(a), len(b))
	for ; i < n; i++ {
		if a[i] != b[i] {
			break
		}
	}

	for ; i < n; i++ {
		x, y := a[i], b[i]
		if !isdigit(x) {
			if !isdigit(y) {
				return a < b
			}
			return true
		} else if !isdigit(y) {
			return false
		}
	}

	if len(a) < len(b) && isdigit(b[n]) {
		return true
	}

	if len(a) > len(b) && isdigit(a[n]) {
		return false
	}

	return a < b
}

func isdigit(a byte) bool {
	return a >= '0' && a <= '9'
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
