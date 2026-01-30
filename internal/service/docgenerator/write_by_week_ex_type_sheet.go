package docgenerator

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
	"github.com/xuri/excelize/v2"
	"sort"
	"strconv"
)

func (s *serviceImpl) writeByWeekAndExTypeSummarySheet(f *excelize.File, summary map[string]map[string]*summary.WeekSummary) {
	sheet := ByWeekAndExTypeSummarySheet
	_, _ = f.NewSheet(sheet)

	exercisesMap := make(map[string]struct{})
	exercises := make([]string, 0)
	weeks := make([]string, 0)
	for week, exMap := range summary {
		weeks = append(weeks, week)
		for exName := range exMap {
			exercisesMap[exName] = struct{}{}
		}
	}
	sort.Strings(weeks)
	for exName := range exercisesMap {
		exercises = append(exercises, exName)
	}
	sort.Strings(exercises)

	headers := []string{
		messages.WorkoutDate,
	}
	for _, ex := range exercises {
		headers = append(headers, ex)
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}

	row := 2

	for _, week := range weeks {
		sum := summary[week]

		_ = f.SetCellValue(sheet, string('A')+strconv.Itoa(row), week)
		for i, ex := range exercises {
			if sum[ex] == nil {
				continue
			}
			switch {
			case sum[ex].SumWeight > 0:
				_ = f.SetCellValue(sheet, string(rune('A'+i+1))+strconv.Itoa(row), fmt.Sprintf("%0.f кг", sum[ex].SumWeight))
			case sum[ex].SumMinutes > 0:
				_ = f.SetCellValue(sheet, string(rune('A'+i+1))+strconv.Itoa(row), fmt.Sprintf("%d мин", sum[ex].SumMinutes))
			case sum[ex].SumMeters > 0:
				_ = f.SetCellValue(sheet, string(rune('A'+i+1))+strconv.Itoa(row), fmt.Sprintf("%d метров", sum[ex].SumMeters))
			}
		}
		row++
	}
}
