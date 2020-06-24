package googlecloudfunctions

import (
	"fmt"
	"net/http"
	"os"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"strconv"
)

var client *restclient.TookanClient
var secret string
var _baseURL string
var _smsSecret string
var _smsSid string
var _smsNumber string

func init() {

	secret = os.Getenv("SECRET")
	if secret == "" {
		panic("Please set SECRET")
	}

	_baseURL = os.Getenv("BASE_URL")
	_smsSecret = os.Getenv("SMS_SECRET")
	_smsSid = os.Getenv("SMS_SID")
	_smsNumber = os.Getenv("SMS_NUMBER")
	client = restclient.FromEnv()
}

func baseURL() string {
	if _baseURL == "" {
		panic("Please set BASE_URL")
	}
	return _baseURL
}

func smsSecret() string {
	if _smsSecret == "" {
		panic("Please set SMS_SECRET")
	}
	return _smsSecret
}

func smsSid() string {
	if _smsSid == "" {
		panic("Please set SMS_SID")
	}
	return _smsSid
}

func smsNumber() string {
	if _smsNumber == "" {
		panic("Please set SMS_NUMBER")
	}
	return _smsNumber
}

func checkSecret(w http.ResponseWriter, r *http.Request, expectedSecret string) bool {
	givenSecret, ok := r.URL.Query()["secret"]
	if !ok || len(givenSecret) != 1 || givenSecret[0] != expectedSecret {
		fmt.Fprintf(w, "Bad secret: ")
		return false
	}
	return true
}

func getTeamId(w http.ResponseWriter, r *http.Request) (int, error) {
	teamIdStr, ok := r.URL.Query()["teamId"]
	if !ok || len(teamIdStr) != 1 ||teamIdStr[0] == "" {
		fmt.Fprintf(w, "Please set teamId")
		return 0, fmt.Errorf("Please set teamId")
	}
	teamId, err := strconv.Atoi(teamIdStr[0])
	if err != nil {
		fmt.Fprintf(w, "Please set a teamId that is a integer")
		return 0, err
	}

	return teamId, nil
}
