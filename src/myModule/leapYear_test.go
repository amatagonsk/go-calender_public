package myModule

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println(" ========== run test ========== ")
	code := m.Run()
	fmt.Println(" ========== end test ========== ")
	os.Exit(code)
}

func TestIsLeapYear(t *testing.T) {
	testCase := []struct {
		year   int
		isLeap bool
	}{
		{1900, false},
		{1904, true},
		{2044, true},
		{2048, true},
	}

	for idx, item := range testCase {
		actual := IsLeapYear(item.year)
		if item.isLeap != actual {
			t.Errorf("pattern %d: want %v, actual %v", idx, item.isLeap, actual)
		}
	}
}

func TestParseDay(t *testing.T) {
	//sep := time.Date(2021, 2, 1,0,0,0,0, targetLocal)
	//expected := []int{1,9,17,28}
	//testStr := []string{"first-monday", "second-tuesday", "third-wednesday", "forth-sunday"}

	feb := time.Date(2021, 9, 1, 0, 0, 0, 0, TargetTimeLocal)
	expected := []int{5, 6, 14, 15, 26, -1}
	testStr := []string{"first-sunday", "first-monday", "second-tuesday", "third-wednesday", "forth-sunday", "fifth-saturday"}

	var actual []int
	for _, s := range testStr {
		actual = append(actual, ParseDay(feb, s))
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestGetJapaneseHolidays(t *testing.T) {
	//feb := time.Date(2021, 2, 1, 0, 0, 0, 0, targetLocal)
	//actualFeb := GetJapaneseHolidays(feb)
	//fmt.Println(actualFeb)

	//mar := time.Date(1999, 3, 1, 0, 0, 0, 0, targetLocal)
	//actualMar := GetJapaneseHolidays(mar)
	//fmt.Println(actualMar)

	sep := time.Date(2021, 9, 1, 0, 0, 0, 0, TargetTimeLocal)
	actualSep := GetJapaneseHolidays(sep)
	expected := []int{20, 23}
	//fmt.Println(actualSep)

	if !reflect.DeepEqual(expected, actualSep) {
		t.Errorf("expected: %v, actual: %v", expected, actualSep)
	}
}

//func test(ary []string) {
//	sort.Strings(ary)
//	fmt.Print(ary,"\n")
//
//	sort.Sort(sort.Reverse(sort.StringSlice(ary)))
//}
//
//func TestSort(t *testing.T) {
//	testCase := []string{"hog", "pogCamp", "ho", "h"}
//	fmt.Print(testCase,"\n")
//	test(testCase)
//	fmt.Print(testCase,"\n")
//}
