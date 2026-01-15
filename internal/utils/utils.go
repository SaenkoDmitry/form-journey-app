package utils

import (
	"fmt"
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
