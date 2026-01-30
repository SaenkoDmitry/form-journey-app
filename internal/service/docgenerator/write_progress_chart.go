package docgenerator

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"

	"github.com/xuri/excelize/v2"

	"github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
)

func (s *serviceImpl) writeAllProgressCharts(
	f *excelize.File,
	progresses map[string]map[string]*summary.Progress,
	redHeaderStyle int,
	greenHeaderStyle int,
	blueHeaderStyle int,
) {

	sheet := ByExerciseSummarySheet
	_, _ = f.NewSheet(sheet)

	exercises := make([]string, 0)
	for exercise, progress := range progresses {
		if len(progress) == 0 {
			continue
		}
		exercises = append(exercises, exercise)
	}
	sort.Strings(exercises)

	row := 1
	for i, e := range exercises {
		progress, ok := progresses[e]
		if !ok {
			continue
		}
		style := blueHeaderStyle
		switch {
		case i%3 == 1:
			style = redHeaderStyle
		case i%3 == 2:
			style = greenHeaderStyle
		}
		row = s.writeProgressChart(f, sheet, e, progress, row, style)
	}
}

func (s *serviceImpl) writeProgressChart(
	f *excelize.File,
	sheet string,
	exercise string,
	progress map[string]*summary.Progress,
	firstRow int,
	headerStyle int,
) int {
	dates := make([]string, 0, len(progress))
	for d := range progress {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	exerciseUnitType := constants.RepsUnit
	for _, p := range progress {
		if p.SumMinutes > 0 {
			exerciseUnitType = constants.MinutesUnit
			break
		} else if p.SumMeters > 0 {
			exerciseUnitType = constants.MetersUnit
			break
		} else {
			break
		}
	}

	switch exerciseUnitType {
	case constants.RepsUnit:
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", firstRow), messages.WorkoutDate)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", firstRow), "Макс вес (кг)")
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", firstRow), "Макс кол-во повторов")
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", firstRow), "Средний вес (кг)")
	case constants.MinutesUnit:
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", firstRow), messages.WorkoutDate)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", firstRow), "Макс время (минут)")
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", firstRow), "Мин время (минут)")
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", firstRow), "Всего (минут)")
	case constants.MetersUnit:
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", firstRow), messages.WorkoutDate)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", firstRow), "Макс дистанция (метры)")
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", firstRow), "Мин дистанция (метры)")
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", firstRow), "Всего (метров)")
	}
	_ = f.SetRowStyle(ByExerciseSummarySheet, firstRow, firstRow, headerStyle)

	lastRow := firstRow + len(dates)

	fmt.Println("len(dates):", len(dates))
	row := firstRow + 1
	for _, d := range dates {
		switch exerciseUnitType {
		case constants.RepsUnit:
			_ = f.SetCellValue(sheet, "A"+strconv.Itoa(row), d)
			_ = f.SetCellValue(sheet, "B"+strconv.Itoa(row), float64(progress[d].MaxWeight))
			_ = f.SetCellValue(sheet, "C"+strconv.Itoa(row), float64(progress[d].MaxReps))
			_ = f.SetCellValue(sheet, "D"+strconv.Itoa(row), math.Round(float64(progress[d].AvgWeight)))
		case constants.MinutesUnit:
			_ = f.SetCellValue(sheet, "A"+strconv.Itoa(row), d)
			_ = f.SetCellValue(sheet, "B"+strconv.Itoa(row), float64(progress[d].MaxMinutes))
			_ = f.SetCellValue(sheet, "C"+strconv.Itoa(row), float64(progress[d].MinMinutes))
			_ = f.SetCellValue(sheet, "D"+strconv.Itoa(row), math.Round(float64(progress[d].SumMinutes)))
		case constants.MetersUnit:
			fmt.Println("progress[d]", progress[d])
			_ = f.SetCellValue(sheet, "A"+strconv.Itoa(row), d)
			_ = f.SetCellValue(sheet, "B"+strconv.Itoa(row), float64(progress[d].MaxMeters))
			_ = f.SetCellValue(sheet, "C"+strconv.Itoa(row), float64(progress[d].MinMeters))
			_ = f.SetCellValue(sheet, "D"+strconv.Itoa(row), math.Round(float64(progress[d].SumMeters)))
		}
		row++
	}

	chart := &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				Name:       fmt.Sprintf("%s!$B$%d", sheet, firstRow),
				Categories: fmt.Sprintf("%s!$A$%d:$A$%d", sheet, firstRow+1, lastRow),
				Values:     fmt.Sprintf("%s!$B$%d:$B$%d", sheet, firstRow+1, lastRow),
				Marker: excelize.ChartMarker{
					Symbol: "circle",
					Size:   6,
				},
			},
			{
				Name:       fmt.Sprintf("%s!$C$%d", sheet, firstRow),
				Categories: fmt.Sprintf("%s!$A$%d:$A$%d", sheet, firstRow+1, lastRow),
				Values:     fmt.Sprintf("%s!$C$%d:$C$%d", sheet, firstRow+1, lastRow),
				Marker: excelize.ChartMarker{
					Symbol: "circle",
					Size:   6,
				},
			},
			{
				Name:       fmt.Sprintf("%s!$D$%d", sheet, firstRow),
				Categories: fmt.Sprintf("%s!$A$%d:$A$%d", sheet, firstRow+1, lastRow),
				Values:     fmt.Sprintf("%s!$D$%d:$D$%d", sheet, firstRow+1, lastRow),
				Marker: excelize.ChartMarker{
					Symbol: "circle",
					Size:   6,
				},
			},
		},
		Title: []excelize.RichTextRun{
			{Text: exercise},
		},
		Legend: excelize.ChartLegend{
			Position: "bottom",
		},
	}

	err := f.AddChart(sheet, fmt.Sprintf("F%d", firstRow+1), chart)
	if err != nil {
		fmt.Println("error while build chart:", err.Error())
	}

	return firstRow + betweenExerciseRows
}

const (
	betweenExerciseRows = 17
)
