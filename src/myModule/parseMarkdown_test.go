package myModule

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseEvents(t *testing.T) {
	testCase := []struct {
		event  []string
		expect []event
	}{
		{
			[]string{"1:00~2:05", "9 : 00	~	10 : 00"},
			[]event{
				{start: time.Hour*1 + time.Minute*0, end: time.Hour*2 + time.Minute*5},
				{start: time.Hour*9 + time.Minute*0, end: time.Hour*10 + time.Minute*0},
			},
		},
		{
			[]string{"	0:00	~	25:05", "  9 : 00	~	1 0 : 0 0  "},
			[]event{
				{start: time.Hour*0 + time.Minute*0, end: time.Hour*25 + time.Minute*5},
				{start: time.Hour*9 + time.Minute*0, end: time.Hour*10 + time.Minute*0},
			},
		},
	}
	for _, test := range testCase {
		actual := parseEvents(test.event)
		expect := test.expect
		assert.Equal(t, actual, expect)
	}
}

func TestCalenderConverter(t *testing.T) {
	testCase := []struct {
		cStr   calenderMdString
		expect CalenderMdInt
	}{
		{
			calenderMdString{
				year: 2020, monthDay: "12/02", description: "this is test.",
				events: []string{"12:00 ~ 13:00", "13:30	~	13:40"},
			},
			CalenderMdInt{
				Year: 2020, Month: 12, Day: 2, Description: "this is test.",
				Events: []event{
					{start: time.Hour*12 + time.Minute*0, end: time.Hour*13 + time.Minute*0},
					{start: time.Hour*13 + time.Minute*30, end: time.Hour*13 + time.Minute*40},
				},
			},
		},
		{
			calenderMdString{
				year: 2020, monthDay: "12/33", description: "this is test.",
				events: []string{"12:00 ~ 13:00", "13:30	~	13:40"},
			},
			CalenderMdInt{},
		},
	}
	for _, test := range testCase {
		actual, _ := CalenderConverter(test.cStr)
		expect := test.expect
		assert.Equal(t, actual, expect)
	}
}
