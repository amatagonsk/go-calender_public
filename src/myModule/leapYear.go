package myModule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var TargetTimeLocal = time.Local

func IsLeapYear(year int) bool {
	return year%400 == 0 || year%4 == 0 && year%100 != 0
}

func DayOfMonth(year int, month time.Month) int {
	day := 0
	return time.Date(year, month+1, day, 0, 0, 0, 0, TargetTimeLocal).Day()
}

func StartOffsetOfMonth(year int, month time.Month) int {
	day := 1
	return int(time.Date(year, month, day, 0, 0, 0, 0, TargetTimeLocal).Weekday())
}

func IsNowTargetYearMonth(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month()
}

func IndexOfElements(arr []interface{}, target interface{}) int {
	switch arr[0].(type) {
	case int:
		for i, elem := range arr {
			if elem.(int) == target.(int) {
				return i
			}
		}
	case string:
		for i, elem := range arr {
			if elem.(string) == target.(string) {
				return i
			}
		}
	}
	return -1
}

func IndexOfStrElements(arr []string, target string) int {
	for i, elem := range arr {
		if elem == target {
			return i
		}
	}
	return -1
}

func IndexOfIntElement(arr []int, target int) int {
	for i, elem := range arr {
		if elem == target {
			return i
		}
	}
	return -1
}

/*
yyyy-mm-dd
*/
func ParseYearMonthDay(yearMonthDay string) (*time.Time, error) {
	layout := "2006/01/02"
	split := strings.Split(yearMonthDay, "-")

	year, yearErr := strconv.Atoi(split[0])
	if yearErr != nil {
		return nil, yearErr
	}

	month, monthErr := strconv.Atoi(split[1])
	if monthErr != nil {
		return nil, monthErr
	}
	day, dayErr := strconv.Atoi(split[2])
	if dayErr != nil {
		return nil, dayErr
	}

	value := fmt.Sprintf("%04d/%02d/%02d", year, month, day)
	t, e := time.Parse(layout, value)

	//fmt.Println(e)
	if e != nil {
		return nil, e
	}
	return &t, nil
}

/*
yyyy-mm
*/
func ParseYearMonth(yearMonth string) (*time.Time, error) {
	arg := yearMonth + "-01"
	return ParseYearMonthDay(arg)
}

func ParseDay(t time.Time, numberedWeekday string) int {
	nums := []interface{}{"first", "second", "third", "forth", "fifth"}
	weeks := []interface{}{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
	split := strings.Split(numberedWeekday, "-")

	num := IndexOfElements(nums, split[0])
	targetWeek := IndexOfElements(weeks, split[1])

	startWeek := int(t.Weekday())

	if targetWeek >= startWeek {
		num--
	}
	result := targetWeek + startWeek + 7*num + 2

	if result > DayOfMonth(t.Year(), t.Month()) {
		return -1
	}
	return result
}
