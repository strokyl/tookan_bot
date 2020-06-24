package googlecloudfunctions

import (
	"fmt"
	"net/http"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"github.com/strokyl/tookan_bot/pkg/planning"
	"github.com/strokyl/tookan_bot/pkg/generate_hash"
	"github.com/strokyl/tookan_bot/pkg/googlecloudfunctions/sms_util"
	"github.com/gobuffalo/packr"
	"html/template"
	"time"
	"strconv"
)

func getJobId(w http.ResponseWriter, r *http.Request) (restclient.JobId, error) {
	jobIdStr, ok := r.URL.Query()["jobId"]
	if !ok || len(jobIdStr) != 1 || jobIdStr[0] == "" {
		fmt.Fprintf(w, "Please set jobId")
		return 0, fmt.Errorf("Please set jobUd")
	}
	jobId, err := strconv.Atoi(jobIdStr[0])
	if err != nil {
		fmt.Fprintf(w, "Please set a jobId that is a integer")
		return 0, err
	}

	return restclient.JobId(jobId), nil
}

func getPlan(w http.ResponseWriter, r *http.Request) (planning.APlan, error) {
	planStr, ok := r.URL.Query()["plan"]
	if !ok || len(planStr) != 1 || planStr[0] == "" {
		fmt.Fprintf(w, "Please set plan")
		return planning.APlan{}, fmt.Errorf("Please set plan")
	}
	plan, err := planning.NewPlan(planStr[0])
	if err != nil {
		fmt.Fprintf(w, "Please set a valid plan: XX:XX à XX:XX")
		return planning.APlan{}, err
	}

	return plan, nil
}

func getDate(w http.ResponseWriter, r *http.Request) (time.Time, error) {
	dateStr, ok := r.URL.Query()["date"]
	if !ok || len(dateStr) != 1 || dateStr[0] == "" {
		fmt.Fprintf(w, "Please set date")
		return time.Time{}, fmt.Errorf("Please set date")
	}
	date, err := time.Parse("2/1/2006", dateStr[0])
	if err != nil {
		fmt.Fprintf(w, "Please set a valid date: DD/MM/YYYY")
		return time.Time{}, err
	}

	return date, nil
}

func ChangeHour(w http.ResponseWriter, r *http.Request) {
	jobId, err := getJobId(w, r)
	if err != nil {
		return
	}
	hashSalt := generate_hash.Secret(secret)
	if !checkSecret(w, r, string(generate_hash.GenerateHash(hashSalt, string(jobId)))) {
		return
	}
	plan, err := getPlan(w, r)
	if err != nil {
		return
	}
	date, err := getDate(w, r)
	if err != nil {
		return
	}
	task, err := client.GetTask(jobId)
	if err != nil {
		panic(err)
	}
	if task == nil {
		fmt.Fprintf(w, "Votre panier n'existe plus")
		return
	}
	newTaskName := planning.ReplacePlanLabelOnTaskName(task.Name, plan, date)
	err = client.RenameTask(jobId, newTaskName)
	if err != nil {
		panic(err)
	}
	_, failure := client.AutoAssignTasks([]restclient.Task{*task})
	if len(failure) != 0 {
		panic(fmt.Errorf("Could not autoassign task %v", *task))
	}

	fmt.Fprintln(w, "Ton Panier à été mis à jour. Tu ne reçois pas d’email de confirmation cette fois mais pas de panique c’est bien au chaud chez nous. Merci de ta générosité et flexibilité #PourEux. On fait de notre max pour que nos biclous arrivent à toi !")
}

func ChangeHourForm(w http.ResponseWriter, r *http.Request) {
	jobId, err := getJobId(w, r)
	if err != nil {
		return
	}
	hashSalt := generate_hash.Secret(secret)
	if !checkSecret(w, r, string(generate_hash.GenerateHash(hashSalt, string(jobId)))) {
		return
	}
	box := packr.NewBox("../../templates")
	templateString, err := box.FindString("form.html")
	if err != nil {
		panic(err)
	}
	template, err := template.New("form").Parse(templateString)
	if err != nil {
		panic(err)
	}

	template.Execute(w, struct {
		Plans []sms_util.PlanLinks
		CancelLink string
	} {
		Plans: sms_util.GetAllPlanLinks(hashSalt, baseURL(), jobId),
			CancelLink: sms_util.GetCancelLink(hashSalt, baseURL(), jobId),
	})
}
