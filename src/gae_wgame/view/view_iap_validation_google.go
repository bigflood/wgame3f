package view

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"

	"gae_wgame/db"
)

type googleReceipt struct {
	Data      googleReceiptData `json:"data"`
	Signature string            `json:"signature"`
}

type googleReceiptData struct {
	PackageName   string `json:"packageName"`
	ProductID     string `json:"productId"`
	PurchaseToken string `json:"purchaseToken"`
}

type googleAPIResult struct {
	Error googleAPIError `json:"error"`
	body  string
}

type googleAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func iapValidationGoogle(ctx context.Context, id, receipt string, result *iapValidationResult) (string, error) {
	setting, err := db.DB.GetIapSetting(ctx, id)
	if err != nil {
		return "", err
	}

	client := urlfetch.Client(ctx)

	if setting.GoogleAccToken == "" {
		err = iapSettingRefreshGoogleToken(ctx, client, setting)

		if err != nil {
			return "", err
		}

		if setting.GoogleAccToken == "" {
			return "", fmt.Errorf("failed to refresh google access_token")
		}
	}

	data := &googleReceipt{}
	err = json.Unmarshal(([]byte)(receipt), data)
	if err != nil {
		return "", err
	}

	errorCode, msg, err := doGoogleAPI(client, setting, data)

	if errorCode == 401 { // Authorization failed
		err = iapSettingRefreshGoogleToken(ctx, client, setting)

		if err != nil {
			return "", err
		}

		errorCode, msg, err = doGoogleAPI(client, setting, data)

		if err != nil {
			return "", err
		}
	}

	if errorCode != 0 {
		return fmt.Sprintf("not valid:%v, %v", errorCode, msg), nil
	}

	return fmt.Sprintf("valid:%v", msg), nil
}

func doGoogleAPI(client *http.Client, setting *db.IapSettingInfo, data *googleReceipt) (int, string, error) {
	packageName := data.Data.PackageName
	productID := data.Data.ProductID
	purchaseToken := data.Data.PurchaseToken

	var url = "https://www.googleapis.com/androidpublisher/v2/applications/" + packageName + "/purchases/subscriptions/" + productID + "/tokens/" + purchaseToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, "", err
	}

	req.Header.Add("Authorization", "Bearer "+setting.GoogleAccToken)

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	bodyString := string(body)
	if strings.Contains(bodyString, "error") {
		apiResult := &googleAPIResult{}
		err = json.Unmarshal(body, &apiResult)
		if err != nil {
			return -1, "", err
		}
		return apiResult.Error.Code, apiResult.Error.Message, nil
	}

	return 0, bodyString, nil
}
