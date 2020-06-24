package planning

import (
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"strings"
	"time"
	"regexp"
	"strconv"
	"fmt"
)

type APlan struct {
	label string
	cachedEnd *int
}

var allPlans = []APlan{
	APlan{
		label: "12:00 à 15:00",
	},
	APlan{
		label: "15:00 à 18:00",
	},
	APlan{
		label: "18:00 à 21:00",
	},
}

var labelRegexp = regexp.MustCompile("\\d{1,2}:\\d{1,2} à (\\d{1,2}):(\\d{1,2})")

func GetNextPlans() []APlan {
	result := make([]APlan, 0)
	now := time.Now()
	for _, plan := range allPlans {
		if !plan.EndBefore(now) {
			result = append(result, plan)
		}
	}
	return result
}

func GetAllPlans() []APlan {
	return allPlans
}

func NewPlan(s string) (APlan, error) {
	matchs := labelRegexp.FindStringSubmatch(s)
	if matchs == nil ||  len(matchs) == 0 || matchs[0] != strings.TrimSpace(s) {
		return APlan{}, fmt.Errorf("%s is not a valid plan", s)
	}
	return APlan{
		label: s,
	}, nil
}

func UnsafeNewPlan(s string) APlan {
	plan, err := NewPlan(s)
	if err != nil {
		panic(err)
	}
	return plan
}

func (plan *APlan) extractEndInTotalMinutes() int {
	if (plan.cachedEnd != nil) {
		return *plan.cachedEnd
	}

	matchs := labelRegexp.FindStringSubmatch(plan.label)
	if len(matchs) != 3 {
		panic("Should have find two subgroups")
	}
	hour, err := strconv.Atoi(matchs[1])
	if err != nil {
		panic(err)
	}
	minute, err := strconv.Atoi(matchs[2])
	if err != nil {
		panic(err)
	}
	plan.cachedEnd = new(int)
	*plan.cachedEnd = hour*60 + minute
	return *plan.cachedEnd
}

func (plan *APlan) EndBefore(t time.Time) bool {
	location, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		panic(err)
	}
	tInLocal := t.In(location)
	currentTotalMinute := tInLocal.Hour()*60 + tInLocal.Minute()
	return currentTotalMinute >= plan.extractEndInTotalMinutes()
}

func (plan *APlan) ToString() string {
	return plan.label
}

func (plan *APlan) AbsoluteString(day time.Time) string {
	dayString := day.Format("02/01/2006")
	return plan.label + " " + dayString
}

func SeparatePlannedAndUnplannedTask(tasks []restclient.Task) (inTime []restclient.Task, outdated []restclient.Task, notPlannedToday []restclient.Task){
	inTime = make([]restclient.Task, 0)
	outdated = make([]restclient.Task, 0)
	notPlannedToday = make([]restclient.Task, 0)
	for _, task := range tasks {
		plans := FindTodayPlanTask(task)
		if (len(plans) == 0) {
			notPlannedToday = append(notPlannedToday, task)
		} else if (AllEndBefore(plans, time.Now())) {
			outdated = append(outdated, task)
		} else {
			inTime = append(inTime, task)
		}
	}

	return inTime, outdated, notPlannedToday
}

func AllEndBefore(plans []APlan, aTime time.Time) bool {
	if len(plans) == 0 {
		return false
	}
	for _, plan := range plans {
		if !plan.EndBefore(aTime) {
			return false
		}
	}

	return true
}

func FindTodayPlanTask(task restclient.Task) []APlan {
	if !(TaskForToday(task.Name)) {
		return nil
	}
	labels := labelRegexp.FindAllString(task.Name, -1)
	if (labels != nil && len(labels) != 0) {
		result := make([]APlan, len(labels))
		for i, plan := range labels {
			result[i]  = APlan{label: plan}
		}
		return result
	}
	return []APlan{}
}

var dateRegexp = regexp.MustCompile("\\d{1,2}/\\d{1,2}/\\d{4}")
func TaskForToday(taskName string) bool {
	dateStr := dateRegexp.FindString(taskName)
	if dateStr == "" {
		//if no date information we assume the task is for today
		return true
	}

	date, err := time.Parse("2/1/2006", dateStr)
	if err != nil {
		//invalid date information we assume the task is for today
		println("Could not parse date: " + dateStr)
		println(err)
		return true
	}
	now := time.Now()
	nowYear, nowMonth, nowDay := now.Date()
	realYear, realMonth, realDay := date.Date()
	return nowYear == realYear && nowMonth == realMonth && nowDay == realDay
}

var doubleSpace = regexp.MustCompile("\\s+")
var multipleLabelRegexp = regexp.MustCompile("(\\d{1,2}:\\d{1,2} à \\d{1,2}:\\d{1,2}[\n, ]*)+")
func ReplacePlanLabelOnTaskName(oldTaskName string, newPlan APlan, day time.Time) string {
	withoutPreviousDate := dateRegexp.ReplaceAllString(oldTaskName, "")
	mayContainDoubleSpace := multipleLabelRegexp.ReplaceAllString(withoutPreviousDate, newPlan.AbsoluteString(day) + " ")
	return strings.TrimSpace(doubleSpace.ReplaceAllString(mayContainDoubleSpace, " "))
}
