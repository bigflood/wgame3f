package view

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"gae_wgame/db"
)

func iapSettingListHandler(w http.ResponseWriter, r *http.Request) *AppError {
	ctx := appengine.NewContext(r)
	infos, err := db.DB.ListIapSettings(ctx)
	if err != nil {
		return AppErrorf(err, "could not list iap settings: %v", err)
	}

	return iapSettingListTmpl.Execute(w, r, infos)
}

func iapSettingFromRequest(r *http.Request) (*db.IapSettingInfo, error) {
	id := mux.Vars(r)["id"]
	ctx := appengine.NewContext(r)
	info, err := db.DB.GetIapSetting(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not find iap setting: %v", err)
	}
	return info, nil
}

func iapSettingDetailHandler(w http.ResponseWriter, r *http.Request) *AppError {
	info, err := iapSettingFromRequest(r)
	if err != nil {
		return AppErrorf(err, "%v", err)
	}

	return iapSettingDetailTmpl.Execute(w, r, info)
}

func iapSettingAddFormHandler(w http.ResponseWriter, r *http.Request) *AppError {
	return iapSettingEditTmpl.Execute(w, r, nil)
}

func iapSettingEditFormHandler(w http.ResponseWriter, r *http.Request) *AppError {
	info, err := iapSettingFromRequest(r)
	if err != nil {
		return AppErrorf(err, "%v", err)
	}

	return iapSettingEditTmpl.Execute(w, r, info)
}

func iapSettingRefreshFormHandler(w http.ResponseWriter, r *http.Request) *AppError {
	info, err := iapSettingFromRequest(r)
	if err != nil {
		return AppErrorf(err, "%v", err)
	}

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	err = iapSettingRefreshGoogleToken(ctx, client, info)
	var msg string
	if err == nil {
		msg = fmt.Sprintf("ok.") //, m.AccessToken)
	} else {
		msg = fmt.Sprintf("failed! \n%v", err.Error())
	}

	http.Redirect(w, r,
		fmt.Sprintf("/admin/iap_setting/%s?%s", info.ID, url.Values{"msg": {msg}}.Encode()),
		http.StatusFound)
	return nil
}

func iapSettingRefreshGoogleToken(ctx context.Context, client *http.Client, info *db.IapSettingInfo) error {
	resp, err := client.PostForm("https://accounts.google.com/o/oauth2/token",
		url.Values{
			"grant_type":    {"refresh_token"},
			"client_id":     {info.GoogleClientID},
			"client_secret": {info.GoogleClientSecret},
			"refresh_token": {info.GoogleRefToken},
		})

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	type Result struct {
		Error       string `json:"error"`
		AccessToken string `json:"access_token"`
	}

	var m Result
	err = json.Unmarshal(body, &m)

	if err != nil {
		return err
	}

	if m.AccessToken != "" {
		info.GoogleAccToken = m.AccessToken
		info.GoogleRefreshTime = time.Now()
		db.DB.UpdateIapSetting(ctx, info)
		return nil
	} else {
		return fmt.Errorf("%v", string(body))
	}
}

// iapSettingFromForm populates the fields of a IapSettingInfo from form values
// (see templates/iap_setting_edit.html).
func iapSettingFromForm(r *http.Request) (*db.IapSettingInfo, error) {

	info := &db.IapSettingInfo{
		ID:                 r.FormValue("id"),
		Title:              r.FormValue("title"),
		Description:        r.FormValue("description"),
		ApplePassword:      r.FormValue("ApplePassword"),
		AmazonSecret:       r.FormValue("AmazonSecret"),
		GooglePublicKey:    r.FormValue("GooglePublicKey"),
		GoogleAccToken:     r.FormValue("GoogleAccToken"),
		GoogleRefToken:     r.FormValue("GoogleRefToken"),
		GoogleClientID:     r.FormValue("GoogleClientID"),
		GoogleClientSecret: r.FormValue("GoogleClientSecret"),
	}

	return info, nil
}

// updateHandler updates the details of a given iap setting.
func iapSettingUpdateHandler(w http.ResponseWriter, r *http.Request) *AppError {
	info, err := iapSettingFromForm(r)
	if err != nil {
		return AppErrorf(err, "could not parse iap setting from form: %v", err)
	}

	if info.ID == "" {
		return AppErrorf(fmt.Errorf("empty ID"), "empty ID")
	}

	info.UpdateTime = time.Now()

	ctx := appengine.NewContext(r)

	err = db.DB.UpdateIapSetting(ctx, info)
	if err != nil {
		return AppErrorf(err, "could not save iap setting: %v", err)
	}
	//go publishUpdate(info.ID)
	http.Redirect(w, r, fmt.Sprintf("/admin/iap_setting/%s", info.ID), http.StatusFound)
	return nil
}

func iapSettingDeleteHandler(w http.ResponseWriter, r *http.Request) *AppError {
	id := mux.Vars(r)["id"]

	ctx := appengine.NewContext(r)
	err := db.DB.DeleteIapSetting(ctx, id)
	if err != nil {
		return AppErrorf(err, "could not delete iap setting: %v", err)
	}
	http.Redirect(w, r, "/admin/iap_setting", http.StatusFound)
	return nil
}
