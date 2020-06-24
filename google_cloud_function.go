package functions

import (
	"github.com/strokyl/tookan_bot/pkg/googlecloudfunctions"
	"net/http"
)


func ChangeHour(w http.ResponseWriter, r *http.Request) {
	googlecloudfunctions.ChangeHour(w, r)
}

func ChangeHourForm(w http.ResponseWriter, r *http.Request) {
	googlecloudfunctions.ChangeHourForm(w, r)
}

func SendChangeHourSMS(w http.ResponseWriter, r *http.Request) {
	googlecloudfunctions.SendChangeHourSMS(w, r)
}

func CancelTask(w http.ResponseWriter, r *http.Request) {
	googlecloudfunctions.CancelTask(w, r)
}
