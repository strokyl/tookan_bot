package restclient

import (
	"github.com/go-resty/resty/v2"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"strings"
	"strconv"
)

type TookanClient struct {
	token string
	debug bool
}

func New(token string, debug bool) *TookanClient {
	if token == "" {
		panic("Tookan client token can't be empty")
	}

	return &TookanClient{
		token: token,
		debug: debug,
	}
}

func FromEnv() *TookanClient {
	token := os.Getenv("TOOKAN_V2_TOKEN")
	if token == "" {
		panic("Please set TOOKAN_V2_TOKEN")
	}

	debug := os.Getenv("DEBUG") == "true"
	return New(token, debug)
}

func TeamIdFromEnv() int {
	teamIdStr := os.Getenv("TOOKAN_TEAM_ID")
	if teamIdStr == "" {
		panic("Please set TOOKAN_TEAM_ID")
	}
	var err error
	teamId, err := strconv.Atoi(teamIdStr)
	if err != nil {
		panic("Please set a TOOKAN_TEAM_ID that is an integer")
	}
	return teamId
}

func (tookanClient *TookanClient) getTasks(teamId int, taskStatus *int, startDate *time.Time) []Task {
	const pickup_type = 0
	if teamId == 0 {
		panic("Tookan client teamId can't be 0")
	}

	client := resty.New().SetDebug(tookanClient.debug)

	payload := map[string]interface{}{
		"api_key": tookanClient.token,
		"team_id": teamId,
		"job_type": pickup_type,
		"is_pagination": 1,
	}

	if taskStatus != nil {
		payload["job_status"] = taskStatus
	}

	if startDate != nil {
		payload["start_date"] = startDate.Format("2006-01-02")
	}


	allTasks := make([]Task, 0)
	page := 1

	for {
		payload["requested_page"] = page
		resp, err := client.R().
			SetHeader("Accept", "application/json").
			SetBody(payload).
			Post("https://api.tookanapp.com/v2/get_all_tasks")

		if err != nil {
			panic(err)
		}

		if resp.IsError() {
			panic(fmt.Sprintf("Http error %s: %s", resp.Status(), resp.String()))
		}

		result := responseTasks{}
		err = json.Unmarshal(resp.Body(), &result)
		if err != nil {
			panic(err)
		}

		if result.Status != 200 {
			panic(fmt.Sprintf("Could not fetch unasigned tasks: %+v\n", result))
		}

		if result.Data == nil {
			panic("Result has no data field: " + resp.String())
		}

		allTasks = append(allTasks, *result.Data...)

		if result.TotalPageCount <= page {
			break
		}
		page = page + 1
	}

	return allTasks
}

func (tookanClient *TookanClient) GetCompletedTasks(teamId int) []Task {
	successfulStatus := 2
	beforeDate := time.Now().Add(-3*24*time.Hour)
	allCompletedTasks := tookanClient.getTasks(teamId, &successfulStatus, &beforeDate)

	filtered := make([]Task, 0)
	todayStart := time.Now().Format("2006-01-02")
	for _, task := range allCompletedTasks {
		if (strings.HasPrefix(task.CompletedDate, todayStart)) {
			filtered = append(filtered, task)
		}
	}

	return filtered
}

const unasigned_status = 6

func (tookanClient *TookanClient) GetUnasignedTasks(teamId int) []Task {
	unasignedStatus := unasigned_status

	ignoreDateBefore := time.Now().Add(-3*24*time.Hour)

	return tookanClient.getTasks(teamId, &unasignedStatus, &ignoreDateBefore)
}

func (tookanClient *TookanClient) GetCancelledTasks(teamId int) []Task {
	cancelledStatus := 9

	ignoreDateBefore := time.Now().Add(-3*24*time.Hour)

	return tookanClient.getTasks(teamId, &cancelledStatus, &ignoreDateBefore)
}

func (tookanClient *TookanClient) MarkATaskUnasigned(id JobId) error {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"api_key": tookanClient.token,
			"job_id": id,
			"job_status": unasigned_status,
		}).
		Post("https://api.tookanapp.com/v2/update_task_status")

	if err != nil {
		return fmt.Errorf("Could not make rest call for markin task %d unasigned: %v", id, err)
	}

	if resp.IsError() {
		return fmt.Errorf("Http error when marking task %d unasigned %s: %s", id, resp.Status(), resp.String())
	}

	result := response{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return fmt.Errorf("Could not serialize json %+v\n", err)
	}

	if result.Status != 200 {
		return fmt.Errorf(fmt.Sprintf("Could not mark task %d unasigned: %+v\n", id, result))
	}
	return nil
}

func (tookanClient *TookanClient) AutoAssignATask(id JobId) error {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"api_key": tookanClient.token,
			"job_id": id,
		}).
		Post("https://api.tookanapp.com/v2/re_autoassign_task")

	if err != nil {
		return fmt.Errorf("Could not make rest call for autoassign task %d: %v", id, err)
	}

	if resp.IsError() {
		return fmt.Errorf("Http error when autoassign task %d %s: %s", id, resp.Status(), resp.String())
	}

	result := response{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return fmt.Errorf("Could not serialize json %+v\n", err)
	}

	if result.Status != 200 {
		return fmt.Errorf(fmt.Sprintf("Could not auto affect task %d: %+v\n", id, result))
	}
	return nil
}

func (tookanClient *TookanClient) AutoAssignTasks(tasks []Task) (success []Task, failure []Task) {
	return actionOnTasks(tookanClient.AutoAssignATask, tasks)
}

func (tookanClient *TookanClient) MarkTasksUnasigned(tasks []Task) (success []Task, failure []Task) {
	return actionOnTasks(tookanClient.MarkATaskUnasigned, tasks)
}

func actionOnTasks(action func(JobId) error, tasks []Task) (success []Task, failure []Task) {
	success = make([]Task, 0)
	failure = make([]Task, 0)
	var err error
	for _, task := range tasks {
		err = action(*task.JobId)
		if err == nil {
			success = append(success, task)
		} else {
			fmt.Fprintf(os.Stderr, "%v", err)
			failure = append(failure, task)
		}
	}
	return success, failure
}

type responseTasks struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data *[]Task `json:"data"`
	TotalPageCount int `json:"total_page_count"`
}

type response struct {
	Status int `json:"status"`
	Message string `json:"message"`
}

type JobId int

type Task struct {
	JobId *JobId `json:"job_id"`
	Name string `json:"job_pickup_name"`
	Description string `json:"job_description"`
	CreationDate string `json:"creation_datetime"`
	CompletedDate string `json:"completed_datetime"`
}
