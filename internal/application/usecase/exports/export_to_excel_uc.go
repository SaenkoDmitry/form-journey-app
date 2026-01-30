package exports

import (
	"bytes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/docgenerator"
	summarysvc "github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
	"sort"
)

type ExportToExcelUseCase struct {
	usersRepo              users.Repo
	exerciseGroupTypesRepo exercisegrouptypes.Repo
	workoutsRepo           workouts.Repo
	exercisesRepo          exercises.Repo
	summaryService         summarysvc.Service
	docGeneratorService    docgenerator.Service
}

func NewExportToExcelUseCase(
	usersRepo users.Repo,
	exerciseGroupTypesRepo exercisegrouptypes.Repo,
	workoutsRepo workouts.Repo,
	exercisesRepo exercises.Repo,
	summaryService summarysvc.Service,
	docGeneratorService docgenerator.Service,
) *ExportToExcelUseCase {
	return &ExportToExcelUseCase{
		usersRepo:              usersRepo,
		exerciseGroupTypesRepo: exerciseGroupTypesRepo,
		workoutsRepo:           workoutsRepo,
		exercisesRepo:          exercisesRepo,
		summaryService:         summaryService,
		docGeneratorService:    docGeneratorService,
	}
}

func (uc *ExportToExcelUseCase) Name() string {
	return "Экспорт в Excel"
}

func (uc *ExportToExcelUseCase) Execute(chatID int64) (*bytes.Buffer, error) {
	groupCodes, err := uc.exerciseGroupTypesRepo.GetAll()
	if err != nil {
		return nil, err
	}
	groupCodesMap := make(map[string]string)
	for _, code := range groupCodes {
		groupCodesMap[code.Code] = code.Name
	}

	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}

	workoutObjs, err := uc.workoutsRepo.FindAll(user.ID)
	if err != nil {
		return nil, err
	}

	totalSummary := uc.summaryService.BuildTotal(workoutObjs, groupCodesMap)
	byDateSummary := uc.summaryService.BuildByDate(workoutObjs)
	weekExerciseTypeSummary := uc.summaryService.BuildByWeekAndExType(workoutObjs, groupCodesMap)

	exerciseObjs, err := uc.exercisesRepo.FindAllByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	sort.Slice(exerciseObjs, func(i, j int) bool {
		return exerciseObjs[i].ExerciseType.ExerciseGroupTypeCode < exerciseObjs[j].ExerciseType.ExerciseGroupTypeCode
	})

	progresses := make(map[string]map[string]*summarysvc.Progress)
	for _, e := range exerciseObjs {
		progresses[e.ExerciseType.Name] = uc.summaryService.BuildExerciseProgress(workoutObjs, e.ExerciseType.Name)
	}

	file, err := uc.docGeneratorService.ExportToFile(workoutObjs, totalSummary, byDateSummary, progresses, groupCodesMap, weekExerciseTypeSummary)
	if err != nil {
		return nil, err
	}

	buf, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf, nil
}
