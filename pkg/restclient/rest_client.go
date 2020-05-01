package restclient

import "github.com/go-resty/resty/v2"
import "encoding/json"
import "fmt"
import "os"
import "time"

type TookanClient struct {
	token string
}

func New(token string) *TookanClient {
	if token == "" {
		panic("Tookan client token can't be empty")
	}

	return &TookanClient{
		token: token,
	}
}



func (tookanClient *TookanClient) GetUnasignedTask(teamId int) []Task {
	const unasigned_status = 6
	const pickup_type = 0

	if teamId == 0 {
		panic("Tookan client teamId can't be 0")
	}

	ignoreDateBefore := time.Now().Add(-3*24*time.Hour)

	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"api_key": tookanClient.token,
			"job_type": pickup_type,
			"job_status": unasigned_status,
			"team_id": teamId,
		}).
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



	filtered := make([]Task, 0)
	for _, task := range *result.Data {
		if task.JobId == nil {
			panic("One of task has no JobId..." + resp.String())
		}
		date, err := time.Parse("2006-01-02 15:04:05", task.CreationDate)
		if err != nil {
			panic(err)
		}
		if (date.After(ignoreDateBefore)) {
			filtered = append(filtered, task)
		} else {
			fmt.Println("Task ignored because it's too old: %+v\n", task)
		}
	}


	return filtered
}

func (tookanClient *TookanClient) AutoAssignATask(id JobId) bool {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"api_key": tookanClient.token,
			"job_id": id,
		}).
		Post("https://api.tookanapp.com/v2/re_autoassign_task")

	if err != nil {
		panic(err)
	}

	if resp.IsError() {
		panic(fmt.Sprintf("Http error %s: %s", resp.Status(), resp.String()))
	}

	result := response{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		panic(err)
	}

	if result.Status != 200 {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Could not auto affect task %d: %+v\n", id, result))
		return false
	}
	return true
}

func (tookanClient *TookanClient) AutoAssignTasks(tasks []Task) (success []Task, failure []Task) {
	success = make([]Task, 0)
	failure = make([]Task, 0)
	for _, task := range tasks {
		if tookanClient.AutoAssignATask(*task.JobId) {
			success = append(success, task)
		} else {
			failure = append(failure, task)
		}
	}
	return success, failure
}

type responseTasks struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data *[]Task `json:"data"`
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
}