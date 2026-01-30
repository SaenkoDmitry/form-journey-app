package docgenerator

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	summarysvc "github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
	"github.com/xuri/excelize/v2"
)

type Service interface {
	ExportToFile(
		workouts []models.WorkoutDay,
		summary map[string]*summarysvc.ExerciseSummary,
		byDateSummary map[string]*summarysvc.DateSummary,
		progresses map[string]map[string]*summarysvc.Progress,
		groupCodesMap map[string]string,
		typeSummary map[string]map[string]*summarysvc.WeekSummary,
	) (*excelize.File, error)
}

type serviceImpl struct {
	summaryService summarysvc.Service
}

func NewService(summaryService summarysvc.Service) Service {
	return &serviceImpl{
		summaryService: summaryService,
	}
}

const (
	DefaultSheet = "Sheet1"

	WorkoutSheet                = "Все тренировки"
	TotalSummarySheet           = "Упражнения"
	ByDateSummarySheet          = "По датам"
	ByWeekAndExTypeSummarySheet = "По неделям & типу упражнения"
	ByExerciseSummarySheet      = "Динамика"
)

func (s *serviceImpl) ExportToFile(
	workouts []models.WorkoutDay,
	summary map[string]*summarysvc.ExerciseSummary,
	byDateSummary map[string]*summarysvc.DateSummary,
	progresses map[string]map[string]*summarysvc.Progress,
	groupCodesMap map[string]string,
	byWeekAndExerciseTypeSummary map[string]map[string]*summarysvc.WeekSummary,
) (*excelize.File, error) {
	f := excelize.NewFile()

	redColor := "#FF746C"
	greenColor := "#6FC276"
	blueColor := "#6488EA"

	redHeaderStyle := HeaderStyle(f, redColor)
	greedHeaderStyle := HeaderStyle(f, greenColor)
	blueHeaderStyle := HeaderStyle(f, blueColor)

	s.writeWorkoutsSheet(f, workouts, groupCodesMap)
	s.writeTotalSummarySheet(f, summary)
	s.writeByDateSummarySheet(f, byDateSummary)
	s.writeByWeekAndExTypeSummarySheet(f, byWeekAndExerciseTypeSummary)
	s.writeAllProgressCharts(f, progresses, redHeaderStyle, greedHeaderStyle, blueHeaderStyle)

	_ = f.SetRowStyle(WorkoutSheet, 1, 1, blueHeaderStyle)
	_ = f.SetRowStyle(TotalSummarySheet, 1, 1, redHeaderStyle)
	_ = f.SetRowStyle(ByWeekAndExTypeSummarySheet, 1, 1, greedHeaderStyle)
	_ = f.SetRowStyle(ByDateSummarySheet, 1, 1, greedHeaderStyle)

	AutoFitColumns(f, WorkoutSheet, 1, 8)
	AutoFitColumns(f, TotalSummarySheet, 1, 7)
	AutoFitColumns(f, ByWeekAndExTypeSummarySheet, 1, 10)
	AutoFitColumns(f, ByDateSummarySheet, 1, 6)
	AutoFitColumns(f, ByExerciseSummarySheet, 1, 4)

	_ = f.DeleteSheet(DefaultSheet)

	f.SetActiveSheet(0)
	return f, nil
}
