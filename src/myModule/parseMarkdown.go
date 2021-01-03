package myModule

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type calenderMdString struct {
	year        int
	monthDay    string
	events      []string
	description string
}

type CalenderMdInt struct {
	Year, Month, Day int
	Events           []event
	Description      string
}

type event struct {
	start time.Duration
	end   time.Duration
}

func CalenderConverter(cStr calenderMdString) (CalenderMdInt, error) {
	//yyyy/mm/dd
	timeStr := fmt.Sprintf("%v/%v", cStr.year, cStr.monthDay)
	timeStr = strings.ReplaceAll(timeStr, "/", "-")
	yearMonthDay, parseErr := ParseYearMonthDay(timeStr)

	if parseErr != nil {
		return CalenderMdInt{}, parseErr
	}
	ev := parseEvents(cStr.events)

	result := CalenderMdInt{
		Year: yearMonthDay.Year(), Month: int(yearMonthDay.Month()), Day: yearMonthDay.Day(), Events: ev, Description: cStr.description,
	}
	return result, nil
}

func parseEvents(ev []string) []event {
	eventSeparator := "~"
	timeSeparator := ":"
	space := regexp.MustCompile(`\s`)

	var result []event
	for _, i := range ev {
		i = space.ReplaceAllString(i, "")
		split := strings.Split(i, eventSeparator)
		startSplit := strings.Split(split[0], timeSeparator)
		endSplit := strings.Split(split[1], timeSeparator)

		startHour, _ := time.ParseDuration(startSplit[0] + "h")
		startMinutes, _ := time.ParseDuration(startSplit[1] + "m")
		endHour, _ := time.ParseDuration(endSplit[0] + "h")
		endMinutes, _ := time.ParseDuration(endSplit[1] + "m")

		result = append(result, event{start: startHour + startMinutes, end: endHour + endMinutes})
	}
	return result
}

func SplitFromMarkdown(strPath string) []calenderMdString {
	splitDir := strings.Split(strPath, "/")
	dirYearStr := splitDir[len(splitDir)-2]

	dirYear, errYear := strconv.Atoi(dirYearStr)
	file, errRead := ioutil.ReadFile(strPath)

	if errYear != nil || errRead != nil {
		log.Print("fail at 'SplitFromMarkdown'")
		return []calenderMdString{}
	}

	newLine := regexp.MustCompile(`\n\s+`)
	thirdHeader := regexp.MustCompile(`###\s+`)
	lines := newLine.ReplaceAllString(string(file), "")

	space := regexp.MustCompile(`\s`)
	splitBySharp := thirdHeader.Split(lines, -1)

	monthDayRegx := regexp.MustCompile("[0-1]?[\\d][\\s]?/[\\s]?[0-3]?[\\d]")
	eventTimeRegx := regexp.MustCompile("[0-2]?[\\d][\\s]?:[\\d]{2}[\\s]*~[\\s]*[0-2]?[\\d][\\s]?:[\\d]{2}")

	var mdStrings []calenderMdString
	for _, item := range splitBySharp[1:] {
		monthDay := monthDayRegx.FindAllString(item, -1)[0]
		dayEvents := eventTimeRegx.FindAllString(item, -1)

		var events []string
		for _, ev := range dayEvents {
			ev = space.ReplaceAllString(ev, "")
			events = append(events, ev)
		}
		//mdStrings = append(mdStrings, calenderMdString{monthDay: monthDay, events: events})
		mdStrings = append(mdStrings, calenderMdString{year: dirYear, monthDay: monthDay, events: events})
	}
	return mdStrings
}
