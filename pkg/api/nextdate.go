package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const F = "20060102"

func afterNow(date, now time.Time) bool {
	if date.After(now) {
		return true
	}
	return false
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	errMsg := "Data format error"

	date, err := time.Parse(F, dstart)
	if err != nil {
		log.Println("Incorrect date", err)
		return "", err
	}

	splitRepeat := strings.Split(repeat, " ")

	switch splitRepeat[0] {
	case "":
		return "", nil
	case "d":
		if len(splitRepeat) != 2 {
			return "", fmt.Errorf(errMsg)
		}
		day, err := strconv.Atoi(splitRepeat[1])
		if err != nil {
			return "", fmt.Errorf(errMsg)
		}
		if day > 0 && day <= 400 {
			for {
				date = date.AddDate(0, 0, day)
				if afterNow(date, now) {
					break
				}
			}
		} else {
			return "", fmt.Errorf(errMsg)
		}
	case "y":
		if len(splitRepeat) > 1 {
			return "", fmt.Errorf(errMsg)
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	case "w":
		if len(splitRepeat) != 2 {
			return "", fmt.Errorf(errMsg)
		}
		dayString := strings.Split(splitRepeat[1], ",")
		dayReapeat := []int{}
		for _, dayStr := range dayString {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf(errMsg)
			}
			dayReapeat = append(dayReapeat, day)
		}
		if len(dayReapeat) > 7 {
			return "", fmt.Errorf(errMsg)
		}
		for i := 1; i < len(dayReapeat); i++ {
			if dayReapeat[i] <= dayReapeat[i-1] {
				return "", fmt.Errorf(errMsg)
			}
		}
		for {
			numDayOfWeek := int(date.Weekday())
			checkSunday := func() {
				if numDayOfWeek == 0 {
					numDayOfWeek = 7
				}
			}
			check := false
			for _, dayWeek := range dayReapeat {
				if numDayOfWeek >= dayWeek {
					if numDayOfWeek >= dayReapeat[len(dayReapeat)-1] {
						date = date.AddDate(0, 0, 7-numDayOfWeek+dayReapeat[0])
						numDayOfWeek = int(date.Weekday())
						checkSunday()
						if afterNow(date, now) {
							check = true
							break
						}
					} else {
						continue
					}
				} else {
					date = date.AddDate(0, 0, dayWeek-numDayOfWeek)
					numDayOfWeek = int(date.Weekday())
					checkSunday()
					if afterNow(date, now) {
						check = true
						break
					}
				}
			}
			if check {
				break
			}
		}
	case "m":
		if len(splitRepeat) < 2 && len(splitRepeat) > 3 {
			return "", fmt.Errorf(errMsg)
		}
		date = date.AddDate(0, 0, 1)
		day := [32]bool{}
		month := [13]bool{}
		monthDayMap := map[int]int{
			1:  1,
			2:  3,
			3:  1,
			4:  2,
			5:  1,
			6:  2,
			7:  1,
			8:  1,
			9:  2,
			10: 1,
			11: 2,
			12: 1,
		}
		dayString := strings.Split(splitRepeat[1], ",")
		checkDate := func() (string, error) {
			monthDate := int(date.Month())
			yearDate := date.Year()
			for i := range day {
				for _, repeatDays := range dayString {
					numDays, err := strconv.Atoi(repeatDays)
					if err != nil || numDays < -2 || numDays > 31 || numDays == 0 {
						return "", fmt.Errorf(errMsg)
					}
					switch numDays {
					case i:
						day[i] = true
					case -1:
						switch monthDayMap[monthDate] {
						case 1:
							day[31] = true
						case 2:
							day[30] = true
						case 3:
							if yearDate%4 == 0 && (yearDate%100 != 0 || yearDate%400 == 0) {
								day[29] = true
							} else {
								day[28] = true
							}
						}
					case -2:
						switch monthDayMap[monthDate] {
						case 1:
							day[30] = true
						case 2:
							day[29] = true
						case 3:
							if yearDate%4 == 0 {
								day[28] = true
							} else {
								day[27] = true
							}
						}
					}
				}
			}
			return "", nil
		}
		if len(splitRepeat) == 2 {
			for i := 1; i < len(month); i++ {
				month[i] = true
			}
		}
		if len(splitRepeat) == 3 {
			monthString := strings.Split(splitRepeat[2], ",")
			for _, repeatMonth := range monthString {
				numMonth, err := strconv.Atoi(repeatMonth)
				if err != nil || numMonth < 1 || numMonth > 12 {
					return "", fmt.Errorf(errMsg)
				}
				for i := range month {
					if i == numMonth {
						month[i] = true
					}
				}
			}
		}
		if _, err := checkDate(); err != nil {
			return "", fmt.Errorf(errMsg)
		}
		for {
			checkDate()
			if day[date.Day()] == true && month[date.Month()] == true {
				if afterNow(date, now) {
					break
				}
			}
			date = date.AddDate(0, 0, 1)
		}
	default:
		return "", fmt.Errorf(errMsg)
	}
	return date.Format(F), nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	if r.FormValue("now") != "" {
		now, _ = time.Parse(F, r.FormValue("now"))
	}
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")
	result, _ := NextDate(now, dstart, repeat)
	w.Write([]byte(result))
}