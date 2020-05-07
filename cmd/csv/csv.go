package main

import (
	"strconv"
	"os"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"encoding/csv"
	"strings"
	"regexp"
)

func removeComma(s string) string {
	return strings.ReplaceAll(s, ",", " ")
}

func extractRepasNum(task restclient.Task) string {
	repasRegexp := regexp.MustCompile(`(\d+) Repas`)
	result := repasRegexp.FindSubmatch([]byte(task.Name))
	if len(result) == 2 {
		return string(result[1])
	} else {
		return "?"
	}
}

func main() {
	client := restclient.FromEnv()
	teamId := restclient.TeamIdFromEnv()
	tasks := client.GetCompletedTasks(teamId)

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()
	writer.Write([]string{
		"job id",
		"nom",
		"description",
		"date de creation",
		"date de fin",
		"nombre de repas",
	})
	for _, task := range tasks {
		err := writer.Write([]string{
			strconv.Itoa(int(*task.JobId)),
			removeComma(task.Name),
			removeComma(task.Description),
			task.CreationDate,
			task.CompletedDate,
			extractRepasNum(task),
		})

		if err != nil {
			panic(err)
		}
	}
}
