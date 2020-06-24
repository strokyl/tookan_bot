package sms_util

import (
	"fmt"
	"net/url"
	"time"
	"github.com/strokyl/tookan_bot/pkg/planning"
	"github.com/strokyl/tookan_bot/pkg/generate_hash"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"strings"
)

func GenerateChangeTaskPlanURL(secret generate_hash.Secret, baseURL string, jobId restclient.JobId, date time.Time, plan planning.APlan) string {
	return fmt.Sprintf("%s/ChangeHour?jobId=%d&plan=%s&date=%s&secret=%s",
		strings.TrimRight(baseURL, "/"),
		int(jobId),
		url.QueryEscape(plan.ToString()),
		url.QueryEscape(date.Format("02/01/2006")),
		url.QueryEscape(string(generate_hash.GenerateHash(secret, string(jobId)))),
	)
}

func GenerateChangeTaskFormURL(secret generate_hash.Secret, baseURL string, jobId restclient.JobId) string {
	return fmt.Sprintf("%s/ChangeHourForm?jobId=%d&secret=%s",
		strings.TrimRight(baseURL, "/"),
		int(jobId),
		url.QueryEscape(string(generate_hash.GenerateHash(secret, string(jobId)))),
	)
}

func GetCancelLink(secret generate_hash.Secret, baseURL string, jobId restclient.JobId) string {
	return fmt.Sprintf("%s/CancelTask?jobId=%d&secret=%s",
		strings.TrimRight(baseURL, "/"),
		int(jobId),
		url.QueryEscape(string(generate_hash.GenerateHash(secret, string(jobId)))),
	)
}

type PlanLinks struct {
	Link string
	Text string
}

func GetAllPlanLinks(secret generate_hash.Secret, baseURL string, jobID restclient.JobId) []PlanLinks {
	result := make([]PlanLinks, 0)
	today := time.Now()
	tomorrow := today.Add(24*time.Hour)
	for _, plan := range planning.GetNextPlans() {
		result = append(result, PlanLinks {
			Text: fmt.Sprintf("Aujourd'hui (%s) de %s", today.Format("02/01"), plan.ToString()),
			Link: GenerateChangeTaskPlanURL(secret, baseURL, jobID, today, plan),
		})
	}

	for _, plan := range planning.GetAllPlans() {
		result = append(result, PlanLinks {
				Text: fmt.Sprintf("Demain (%s) de %s", tomorrow.Format("02/01"), plan.ToString()),
				Link: GenerateChangeTaskPlanURL(secret, baseURL, jobID, tomorrow, plan),
			})
	}

	return result
}

func GetSMSContent(secret generate_hash.Secret, baseURL string, jobId restclient.JobId) string {
	start := `Salut,
Ici Manu #PourEux.
On a raté ton créneau.
Dispo plus tard ?
RDV ici: `

	return start + GenerateChangeTaskFormURL(secret, baseURL, jobId)
}
