package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ScheduleValue struct {
	Days          []int
	StartDate     time.Time
	NumOfSessions int
}

type ScheduleForm struct {
	Days          string
	StartDate     string
	NumOfSessions string
	Errors        map[string]string
	ValidatedData ScheduleValue
}

type Schedule struct {
	Sessions []string `json:"sessions"`
}

const ChaptersCount = 30

func (form *ScheduleForm) Validate() bool {
	form.Errors = make(map[string]string)

	daysReg := regexp.MustCompile("^[1-7](,[1-7])*$")            // days of week from 1 to 7
	numOfSessionsReg := regexp.MustCompile("^[1-9]\\d*$")        // positive integer
	startDateReg := regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2}$") // yyyy-mm-dd format

	daysMatched := daysReg.Match([]byte(form.Days))
	numOfSessionsMatched := numOfSessionsReg.Match([]byte(form.NumOfSessions))
	dateMatched := startDateReg.Match([]byte(form.StartDate))

	if !daysMatched {
		form.Errors["days"] = "Please enter valid integers (range 1 to 7) separated by ',' ex 1,2,4 "
	}
	if numOfSessionsMatched == false {
		form.Errors["numOfSessions"] = "Please enter a valid numOfSessions integer"
	}
	if !dateMatched {
		form.Errors["startDate"] = "Please enter valid date with format 2020-01-30"
	}

	date, dateErr := time.Parse("2006-01-02", form.StartDate)
	if dateErr != nil {
		form.Errors["startDate"] = "Please enter valid date with format 2020-01-31"
	}

	if len(form.Errors) > 0 {
		return false
	}

	var days = make([]int, 7)
	for _, i := range strings.Split(form.Days, ",") {
		day, _ := strconv.Atoi(i)
		days[day-1]++
	}

	sessions, _ := strconv.Atoi(form.NumOfSessions)
	form.ValidatedData = ScheduleValue{Days: days, StartDate: date, NumOfSessions: sessions}

	return true
}

func (sv *ScheduleValue) getSchedule() Schedule {
	curDate := sv.StartDate
	schedule := Schedule{}
	for {
		if len(schedule.Sessions) == sv.NumOfSessions*ChaptersCount {
			break
		}
		weekday := curDate.Weekday()     // start with sunday
		curDay := (int(weekday) + 1) % 7 //my custom ordering for week days start with Saturday

		if sv.Days[curDay] > 0 {
			schedule.Sessions = append(schedule.Sessions, curDate.Format("2006-01-02"))
		}

		curDate = curDate.Add(time.Hour * 24)
	}
	return schedule
}
