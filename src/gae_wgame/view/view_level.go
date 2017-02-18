package view

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	//	"gae_wgame/auth"
	"gae_wgame/db"
)

func levelListHandler(w http.ResponseWriter, r *http.Request) *AppError {
	ctx := appengine.NewContext(r)
	levels, err := db.DB.ListLevels(ctx)
	if err != nil {
		return AppErrorf(err, "could not list levels: %v", err)
	}

	db.SortLevels(levels)
	return levelListTmpl.Execute(w, r, levels)
}

func levelListOfGameHandler(w http.ResponseWriter, r *http.Request) *AppError {
	game := mux.Vars(r)["game"]

	ctx := appengine.NewContext(r)
	levels, err := db.DB.ListLevelsWithGame(ctx, game)
	if err != nil {
		return AppErrorf(err, "could not list levels: %v", err)
	}

	if r.URL.Query().Get("format") == "json" {
		bytes, err := json.Marshal(levels)
		if err != nil {
			return AppErrorf(err, "marshal failed: %v", err)
		}

		w.Write(bytes)
		return nil
	}

	db.SortLevels(levels)
	return levelListTmpl.Execute(w, r, levels)
}

func levelListOfGameStageHandler(w http.ResponseWriter, r *http.Request) *AppError {
	gamestage := mux.Vars(r)["gamestage"]

	ctx := appengine.NewContext(r)
	levels, err := db.DB.ListLevelsWithGameStage(ctx, gamestage)
	if err != nil {
		return AppErrorf(err, "could not list levels: %v", err)
	}

	if r.URL.Query().Get("format") == "json" {
		if levels == nil {
			w.Write(([]byte)("[]"))
			return nil
		}

		bytes, err := json.Marshal(levels)
		if err != nil {
			return AppErrorf(err, "marshal failed: %v", err)
		}

		w.Write(bytes)
		return nil
	}

	db.SortLevels(levels)
	return levelListTmpl.Execute(w, r, levels)
}

func levelFromRequest(ctx context.Context, r *http.Request) (*db.LevelInfo, error) {
	id := mux.Vars(r)["id"]
	level, err := db.DB.GetLevel(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not find level: %v", err)
	}
	return level, nil
}

func levelDetail(r *http.Request) (*db.LevelInfoAndData, error) {
	ctx := appengine.NewContext(r)

	level, err := levelFromRequest(ctx, r)
	if err != nil {
		return nil, err
	}

	data, err := db.DB.GetLevelData(ctx, level.ID)
	if err != nil {
		return nil, err
	}

	responseStruct := &db.LevelInfoAndData{LevelInfo: *level, Data: data.Data}
	return responseStruct, nil
}

func levelDetailHandler(w http.ResponseWriter, r *http.Request) *AppError {
	if r.URL.Query().Get("format") == "json" {
		responseStruct, err := levelDetail(r)

		if err != nil {
			w.Write(([]byte)(fmt.Sprintf("failed: %v", err)))
			return nil
		}

		bytes, err := json.Marshal(responseStruct)
		if err != nil {
			w.Write(([]byte)(fmt.Sprintf("failed: %v", err)))
			return nil
		}

		w.Write(bytes)
		return nil
	} else {
		responseStruct, err := levelDetail(r)

		if err != nil {
			return AppErrorf(err, "%v", err)
		}

		return levelDetailTmpl.Execute(w, r, responseStruct)
	}

}

func levelAddFormHandler(w http.ResponseWriter, r *http.Request) *AppError {
	return levelEditTmpl.Execute(w, r, nil)
}

type LevelForm struct {
	ID                 string
	Game, Stage, Level string
	Title, Description string
	Data               string
}

func levelEditFormHandler(w http.ResponseWriter, r *http.Request) *AppError {
	ctx := appengine.NewContext(r)

	level, err := levelFromRequest(ctx, r)
	if err != nil {
		return AppErrorf(err, "%v", err)
	}

	data, err := db.DB.GetLevelData(ctx, level.ID)
	if err != nil {
		return AppErrorf(err, "%v", err)
	}

	stageName := level.GameStage[len(level.Game)+1:]
	levelName := level.ID[len(level.GameStage)+1:]

	return levelEditTmpl.Execute(w, r,
		LevelForm{
			ID:   level.ID,
			Game: level.Game, Stage: stageName, Level: levelName,
			Title: level.Title, Description: level.Description,
			Data: data.Data})
}

// levelFromForm populates the fields of a Level from form values
// (see templates/edit.html).
func levelFromForm(r *http.Request) (*db.LevelInfo, *db.LevelData, error) {
	gameName := r.FormValue("game")
	stageName := r.FormValue("stage")
	levelName := r.FormValue("level")

	if gameName == "" || stageName == "" || levelName == "" {
		return nil, nil, fmt.Errorf("empty parameter (game or stage or level)")
	}

	gs := fmt.Sprintf("%v-%v", gameName, stageName)
	gsl := fmt.Sprintf("%v-%v", gs, levelName)

	level := &db.LevelInfo{
		ID:          gsl,
		Game:        gameName,
		GameStage:   gs,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	data := &db.LevelData{
		ID:   gsl,
		Data: r.FormValue("data"),
	}

	return level, data, nil
}

func levelUpdate(ctx context.Context, r *http.Request) (*db.LevelInfo, error) {
	level, data, err := levelFromForm(r)
	if err != nil {
		return nil, err
	}

	level.UpdateTime = time.Now()

	err = db.DB.UpdateLevel(ctx, level)
	if err != nil {
		return nil, err
	}

	err = db.DB.UpdateLevelData(ctx, data)
	if err != nil {
		return nil, err
	}

	return level, nil
}

func levelUpdateHandler(w http.ResponseWriter, r *http.Request) *AppError {
	ctx := appengine.NewContext(r)
	level, err := levelUpdate(ctx, r)
	if err != nil {
		return AppErrorf(err, "failed: %v", err)
	}

	//go publishUpdate(level.ID)
	http.Redirect(w, r, fmt.Sprintf("/levels/%v", level.ID), http.StatusFound)
	return nil
}

func appLevelUpdate(r *http.Request) error {
	ctx := appengine.NewContext(r)

	gameName := r.FormValue("game")
	if gameName == "" {
		return fmt.Errorf("empty game field")
	}

	editkey := r.FormValue("editkey")
	if editkey == "" {
		return fmt.Errorf("empty editkey field")
	}

	game, err := db.DB.GetGame(ctx, gameName)
	if err != nil {
		return err
	}

	if game.EditKey != editkey {
		return fmt.Errorf("invalid editkey")
	}

	_, err = levelUpdate(ctx, r)
	return err
}

func appLevelUpdateHandler(w http.ResponseWriter, r *http.Request) *AppError {
	err := appLevelUpdate(r)

	if err != nil {
		w.Write(([]byte)(fmt.Sprintf("failed: %v", err)))
		return nil
	}

	//go publishUpdate(level.ID)
	w.Write(([]byte)("ok"))
	return nil
}

func levelDelete(r *http.Request) error {
	id := mux.Vars(r)["id"]
	ctx := appengine.NewContext(r)
	return db.DB.DeleteLevel(ctx, id)
}

func levelDeleteHandler(w http.ResponseWriter, r *http.Request) *AppError {
	err := levelDelete(r)

	if err != nil {
		return AppErrorf(err, "could not delete level: %v", err)
	}

	http.Redirect(w, r, "/levels", http.StatusFound)
	return nil
}

func appLevelDeleteHandler(w http.ResponseWriter, r *http.Request) *AppError {
	err := levelDelete(r)

	if err != nil {
		w.Write(([]byte)(fmt.Sprintf("failed: %v", err)))
		return nil
	}

	w.Write(([]byte)("ok"))
	return nil
}
