package view

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"

	//	"gae_wgame/auth"
	"gae_wgame/db"
)

var (
	// See template.go
	gameListTmpl, gameEditTmpl, gameDetailTmpl                   *appTemplate
	levelListTmpl, levelEditTmpl, levelDetailTmpl                *appTemplate
	iapSettingListTmpl, iapSettingEditTmpl, iapSettingDetailTmpl *appTemplate

	iapValidationFormTmpl, iapValidationResultTmpl *appTemplate

	statusTmpl *appTemplate
)

// LoadTemplates func
func LoadTemplates() {
	//fmt.Printf("templates dir = %v\n", templatesDir)

	gameListTmpl = parseTemplate("game_list.html")
	gameEditTmpl = parseTemplate("game_edit.html")
	gameDetailTmpl = parseTemplate("game_detail.html")

	levelListTmpl = parseTemplate("level_list.html")
	levelEditTmpl = parseTemplate("level_edit.html")
	levelDetailTmpl = parseTemplate("level_detail.html")

	iapSettingListTmpl = parseTemplate("iap_setting_list.html")
	iapSettingEditTmpl = parseTemplate("iap_setting_edit.html")
	iapSettingDetailTmpl = parseTemplate("iap_setting_detail.html")

	iapValidationFormTmpl = parseTemplate("iap_validation_form.html")
	iapValidationResultTmpl = parseTemplate("iap_validation_result.html")

	statusTmpl = parseTemplate("status.html")
}

// RegisterHandlers func
func RegisterHandlers(r *mux.Router) {

	r.Handle("/", http.RedirectHandler("/games", http.StatusFound))

	r.Methods("GET").Path("/status").
		Handler(AppHandler(statusHandler))

	r.Methods("GET").Path("/games").
		Handler(AppHandler(gameListHandler))
	r.Methods("GET").Path("/games/{id}").
		Handler(AppHandler(gameDetailHandler))
	r.Methods("GET").Path("/admin/games/add").
		Handler(AppHandler(gameAddFormHandler))
	r.Methods("GET").Path("/admin/games/{id}/edit").
		Handler(AppHandler(gameEditFormHandler))

	r.Methods("POST").Path("/admin/games").
		Handler(AppHandler(gameUpdateHandler))
	r.Methods("POST").Path("/admin/games/{id}:delete").
		Handler(AppHandler(gameDeleteHandler)).Name("delete")

	r.Methods("GET").Path("/levels").
		Handler(AppHandler(levelListHandler))
	r.Methods("GET").Path("/levels/game/{game}").
		Handler(AppHandler(levelListOfGameHandler))
	r.Methods("GET").Path("/levels/gamestage/{gamestage}").
		Handler(AppHandler(levelListOfGameStageHandler))
	r.Methods("GET").Path("/levels/{id}").
		Handler(AppHandler(levelDetailHandler))
	r.Methods("GET").Path("/admin/levels/add").
		Handler(AppHandler(levelAddFormHandler))
	r.Methods("GET").Path("/admin/levels/{id}/edit").
		Handler(AppHandler(levelEditFormHandler))

	r.Methods("POST").Path("/admin/levels").
		Handler(AppHandler(levelUpdateHandler))
	r.Methods("POST").Path("/admin/levels/{id}:delete").
		Handler(AppHandler(levelDeleteHandler)).Name("delete")

	r.Methods("POST").Path("/app/levels").
		Handler(AppHandler(appLevelUpdateHandler))
	r.Methods("POST").Path("/app/levels/{id}:delete").
		Handler(AppHandler(appLevelDeleteHandler)).Name("delete")

	r.Methods("GET").Path("/admin/iap_setting").
		Handler(AppHandler(iapSettingListHandler))
	r.Methods("GET").Path("/admin/iap_setting/{id}").
		Handler(AppHandler(iapSettingDetailHandler))
	r.Methods("GET").Path("/admin/iap_setting/{id}/add").
		Handler(AppHandler(iapSettingAddFormHandler))
	r.Methods("GET").Path("/admin/iap_setting/{id}/edit").
		Handler(AppHandler(iapSettingEditFormHandler))
	r.Methods("GET").Path("/admin/iap_setting/{id}/refresh").
		Handler(AppHandler(iapSettingRefreshFormHandler))

	r.Methods("POST").Path("/admin/iap_setting").
		Handler(AppHandler(iapSettingUpdateHandler))
	r.Methods("POST").Path("/admin/iap_setting/{id}:delete").
		Handler(AppHandler(iapSettingDeleteHandler)).Name("delete")

	r.Methods("GET").Path("/iap_validation/{id}/{platform}").
		Handler(AppHandler(iapValidationFormHandler))

	r.Methods("POST").Path("/iap_validation/{id}/{platform}").
		Handler(AppHandler(iapValidationHandler))

	// The following handlers are defined in auth.go and used in the
	// "Authenticating Users" part of the Getting Started guide.
	r.Methods("GET").Path("/login").
		Handler(AppHandler(loginHandler))
	r.Methods("POST").Path("/logout").
		Handler(AppHandler(logoutHandler))
	r.Methods("GET").Path("/oauth2callback").
		Handler(AppHandler(oauthCallbackHandler))
}

func statusHandler(w http.ResponseWriter, r *http.Request) *AppError {
	ctx := appengine.NewContext(r)

	type Result struct {
		StatusCount int
		CronCount   int
		Version     string
		Env         template.HTML
	}

	result := Result{}
	result.StatusCount, _ = db.DB.AddCount(ctx, "StatusCount", 1)
	result.CronCount, _ = db.DB.GetCount(ctx, "CronCount")
	result.Version = appengine.VersionID(ctx)
	result.Env = template.HTML(strings.Join(os.Environ(), "</br>"))

	return statusTmpl.Execute(w, r, result)
}

func int64FromMuxVars(r *http.Request, name string) (int64, error) {
	return strconv.ParseInt(mux.Vars(r)[name], 10, 64)
}

// AppHandler type
// http://blog.golang.org/error-handling-and-go
type AppHandler func(http.ResponseWriter, *http.Request) *AppError

// AppError struct
type AppError struct {
	Error   error
	Message string
	Code    int
}

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

// AppErrorf f
func AppErrorf(err error, format string, v ...interface{}) *AppError {
	return &AppError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
