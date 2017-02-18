package view

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"

	"gae_wgame/db"
)

func iapValidationApple(ctx context.Context, id, receipt string, result *iapValidationResult) (string, error) {
	setting, err := db.DB.GetIapSetting(ctx, id)
	if err != nil {
		return "", err
	}

	client := urlfetch.Client(ctx)

	apiResult, err := doAppleAPI(client, setting, receipt)
	if err != nil {
		return "", err
	}

	if apiResult.Status != 0 {
		return fmt.Sprintf("not valid:%v", apiResult.Status), nil
	}

	return "valid", err
}

type appleAPIRequest struct {
	ReceiptData string `json:"receipt-data"`
}

type appleAPIResponse struct {
	Status int `json:"status"`
}

func doAppleAPI(client *http.Client, setting *db.IapSettingInfo, data string) (*appleAPIResponse, error) {

	reqData := appleAPIRequest{
		ReceiptData: data,
	}

	reqJson, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	var url = "https://buy.itunes.apple.com/verifyReceipt"

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqJson))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	apiResult := &appleAPIResponse{}
	err = json.Unmarshal(body, &apiResult)
	if err != nil {
		return nil, err
	}

	return apiResult, nil
}
