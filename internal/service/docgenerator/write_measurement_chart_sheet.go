package docgenerator

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/docgenerator/helpers"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func (s *serviceImpl) writeMeasurementChartSheet(
	f *excelize.File,
	measurements []*dto.Measurement,
	headerStyle int,
) {
	sheet := MeasurementSheet
	_, _ = f.NewSheet(sheet)

	_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", 1), messages.WorkoutDate)
	_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", 1), "Плечи (см)")
	_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", 1), "Грудь (см)")
	_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", 1), "Руки (см)")
	_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", 1), "Талия (см)")
	_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", 1), "Ягодицы (см)")
	_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", 1), "Бедра (см)")
	_ = f.SetCellValue(sheet, fmt.Sprintf("H%d", 1), "Икры (см)")
	_ = f.SetCellValue(sheet, fmt.Sprintf("I%d", 1), "Вес (кг)")
	//_ = f.SetCellValue(sheet, fmt.Sprintf("H%d", row), "Процент жира")
	//_ = f.SetCellValue(sheet, fmt.Sprintf("H%d", row), "Объем мышц (кг)")
	_ = f.SetRowStyle(sheet, 1, 1, headerStyle)

	for i, m := range measurements {
		_ = f.SetCellValue(sheet, "A"+strconv.Itoa(i+2), m.CreatedAt)
		_ = f.SetCellValue(sheet, "B"+strconv.Itoa(i+2), m.Shoulders)
		_ = f.SetCellValue(sheet, "C"+strconv.Itoa(i+2), m.Chest)
		_ = f.SetCellValue(sheet, "D"+strconv.Itoa(i+2), m.Hands)
		_ = f.SetCellValue(sheet, "E"+strconv.Itoa(i+2), m.Waist)
		_ = f.SetCellValue(sheet, "F"+strconv.Itoa(i+2), m.Buttocks)
		_ = f.SetCellValue(sheet, "G"+strconv.Itoa(i+2), m.Hips)
		_ = f.SetCellValue(sheet, "H"+strconv.Itoa(i+2), m.Calves)
		_ = f.SetCellValue(sheet, "I"+strconv.Itoa(i+2), m.Weight)
	}

	firstRow := 1
	lastRow := firstRow + len(measurements)

	charts := make([]ChartSetting, 0)
	charts = append(charts, ChartSetting{RangeSymbol: "B", CategoryName: "Плечи"})
	charts = append(charts, ChartSetting{RangeSymbol: "C", CategoryName: "Грудь"})
	charts = append(charts, ChartSetting{RangeSymbol: "D", CategoryName: "Руки"})
	charts = append(charts, ChartSetting{RangeSymbol: "E", CategoryName: "Талия"})
	charts = append(charts, ChartSetting{RangeSymbol: "F", CategoryName: "Ягодицы"})
	charts = append(charts, ChartSetting{RangeSymbol: "G", CategoryName: "Бедра"})
	charts = append(charts, ChartSetting{RangeSymbol: "H", CategoryName: "Икры"})
	charts = append(charts, ChartSetting{RangeSymbol: "I", CategoryName: "Вес"})

	for i, ch := range charts {
		chart := makeChart(sheet, firstRow, lastRow, ch.RangeSymbol, ch.CategoryName)
		err := f.AddChart(sheet, fmt.Sprintf("J%d", i*betweenChartRowsCount+2), chart)
		if err != nil {
			fmt.Println("error while build chart:", err.Error())
			return
		}
	}
}

const betweenChartRowsCount = 14

type ChartSetting struct {
	RangeSymbol  string
	CategoryName string
}

func makeChart(sheet string, firstRow, lastRow int, rangeSymbol string, category string) *excelize.Chart {
	return &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				Name:       helpers.FormatCell(sheet, rangeSymbol, firstRow),
				Categories: helpers.FormatDataRange(sheet, "A", "A", firstRow, lastRow),
				Values:     helpers.FormatDataRange(sheet, rangeSymbol, rangeSymbol, firstRow, lastRow),
				Marker: excelize.ChartMarker{
					Symbol: "circle",
					Size:   5,
				},
			},
		},
		Title: []excelize.RichTextRun{
			{
				Text: category,
				Font: &excelize.Font{
					Bold:  true,
					Size:  18,
					Color: constants.SkyBlueColor,
				},
			},
		},
		Legend: excelize.ChartLegend{
			Position: "bottom",
		},
	}
}
