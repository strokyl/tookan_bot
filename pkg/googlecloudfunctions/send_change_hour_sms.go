package googlecloudfunctions

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/strokyl/tookan_bot/pkg/planning"
	"github.com/strokyl/tookan_bot/pkg/googlecloudfunctions/sms_util"
	"github.com/strokyl/tookan_bot/pkg/generate_hash"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"github.com/sfreiberg/gotwilio"
)

func SendChangeHourSMS(w http.ResponseWriter, r *http.Request) {
	if !checkSecret(w, r, secret) {
		return
	}
	teamId, err := getTeamId(w, r)
	if err != nil {
		return
	}
	tasks := client.GetUnasignedTasks(teamId)
	inTime, outDated, notForToday  := planning.SeparatePlannedAndUnplannedTask(tasks)
	for _, task := range outDated {
		SendSMS(w, task)
	}
	response, err := json.MarshalIndent(map[string]interface{}{
		"inTime" : inTime,
		"outDated" : outDated,
		"notForToday" : notForToday,
	}, "", "  ")

	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(response) + "\n")
}

func SendSMS(w http.ResponseWriter, task restclient.Task) {
	to := task.Phone
	hashSalt := generate_hash.Secret(secret)
	data := sms_util.GetSMSContent(hashSalt, baseURL(), *task.JobId)

	twilio := gotwilio.NewTwilioClient(smsSid(), smsSecret())
	twilio.SendSMS(smsNumber(), to, data, "", "")
	fmt.Fprintln(w, "SMS sent:")
	fmt.Fprintln(w, to)
	fmt.Fprintln(w, data)
	fmt.Fprintf(w, "\n\n")
}
