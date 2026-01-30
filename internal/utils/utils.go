package utils

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func FormatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d ч %d мин", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%d мин", minutes)
	}
	return fmt.Sprintf("%d сек", seconds)
}

func BetweenTimes(startedAt time.Time, endedAt *time.Time) string {
	duration := endedAt.Sub(startedAt)
	return FormatDuration(duration)
}

type Exercise struct {
	ID   int64
	Sets []Set
}

type Set struct {
	Reps    int
	Weight  float32
	Minutes int
}

func SplitPreset(preset string) []Exercise {
	if preset == "" {
		return []Exercise{}
	}
	result := make([]Exercise, 0)
	exercises := strings.Split(preset, ";")
	if len(exercises) == 0 {
		return make([]Exercise, 0)
	}
	for i := range exercises {
		temp := strings.Split(exercises[i], ":")
		if len(temp) == 0 {
			continue
		}

		first, second := temp[0], temp[1]
		exerciseTypeID, _ := strconv.ParseInt(first, 10, 64) // nolint

		result = append(result, Exercise{
			ID: exerciseTypeID,
		})

		approaches := strings.Split(second[1:len(second)-1], ",")
		for _, approach := range approaches {
			if strings.Contains(approach, "*") {
				temp2 := strings.Split(approach, "*")
				reps, _ := strconv.ParseInt(temp2[0], 10, 64)
				weight, _ := strconv.ParseFloat(temp2[1], 32)
				result[len(result)-1].Sets = append(result[len(result)-1].Sets, Set{
					Reps:   int(reps),
					Weight: float32(weight),
				})
			} else {
				minutes, _ := strconv.ParseInt(approach, 10, 64)
				result[len(result)-1].Sets = append(result[len(result)-1].Sets, Set{
					Minutes: int(minutes),
				})
			}
		}
	}
	return result
}

func IsValidPreset(preset string) bool {
	pattern := `^(\d+|\d+\*\d+)(,(\d+|\d+\*\d+))*$`
	matched, err := regexp.MatchString(pattern, preset)
	if err != nil {
		return false
	}
	return matched
}

func WrapYandexLink(url string) string {
	return fmt.Sprintf("\n<a href=\"%s\"><b>👀</b></a>", url)
}

func SplitUnits(units string) ([]string, bool) {
	m := make(map[string]struct{})
	for _, unit := range strings.Split(units, ",") {
		if strings.EqualFold(unit, constants.RepsUnit) ||
			strings.EqualFold(unit, constants.WeightUnit) ||
			strings.EqualFold(unit, constants.MinutesUnit) ||
			strings.EqualFold(unit, constants.MetersUnit) {
			m[unit] = struct{}{}
			continue
		}
		return []string{}, false
	}
	arr := make([]string, 0, len(m))
	for k := range m {
		arr = append(arr, k)
	}
	return arr, true
}

func EqualArrays(arr1, arr2 []string) bool {
	m1 := make(map[string]struct{})
	m2 := make(map[string]struct{})
	for _, e := range arr1 {
		m1[e] = struct{}{}
	}
	for _, e := range arr2 {
		m2[e] = struct{}{}
	}
	return reflect.DeepEqual(m1, m2)
}

func FormatDateTime(dateTime time.Time) string {
	dateTimeInMSK := dateTime.Add(3 * time.Hour)
	d := dateTimeInMSK.Format("02.01.2006")
	weekDay := getRussianWeekDay(dateTimeInMSK.Weekday())
	t := dateTimeInMSK.Format("15:04")
	return fmt.Sprintf("%s (%s) в %s", d, weekDay, t)
}

func getRussianWeekDay(weekday time.Weekday) string {
	switch weekday {
	case time.Monday:
		return "ПН"
	case time.Tuesday:
		return "ВТ"
	case time.Wednesday:
		return "СР"
	case time.Thursday:
		return "ЧТ"
	case time.Friday:
		return "ПТ"
	case time.Saturday:
		return "СБ"
	case time.Sunday:
		return "ВС"
	}
	return ""
}

func GetThisWeek(date time.Time) string {
	align := (date.Add(3*time.Hour).Weekday() + 6) % 7
	fmt.Println(date.Add(3*time.Hour).Weekday().String(), int(align))
	from := date.AddDate(0, 0, -int(align)).Add(3 * time.Hour).Format("02.01.06")
	to := date.AddDate(0, 0, 7-int(align)-1).Add(3 * time.Hour).Format("02.01.06")
	return fmt.Sprintf("%s – %s", from, to)
}
