package usecase

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/timer"
	"gorm.io/gorm"

	daytypeusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/daytypes"
	exerciseusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/exercises"
	exportusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/exports"
	groupusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/groups"
	programusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/programs"
	sessionusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/session"
	setusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/sets"
	statsusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/stats"
	timerusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/timers"
	userusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/users"
	workoutusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/docgenerator"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
)

type Container struct {

	// workouts
	ConfirmDeleteWorkoutUC *workoutusecases.ConfirmDeleteUseCase
	DeleteWorkoutUC        *workoutusecases.DeleteUseCase
	CreateWorkoutUC        *workoutusecases.CreateUseCase
	ShowWorkoutProgressUC  *workoutusecases.ShowProgressUseCase
	FindMyWorkoutsUC       *workoutusecases.FindMyUseCase
	StartWorkoutUC         *workoutusecases.StartUseCase
	ConfirmFinishWorkoutUC *workoutusecases.ConfirmFinishUseCase
	FinishWorkoutUC        *workoutusecases.FinishUseCase
	FindWorkoutsByUserUC   *workoutusecases.FindByUserIDUseCase
	StatsWorkoutUC         *workoutusecases.StatsUseCase

	// exercises
	ShowCurrentExerciseSessionUC *sessionusecases.ShowCurrentExerciseSessionUseCase
	ExerciseTypeListUC           *exerciseusecases.ExerciseTypeListUseCase
	FindTypesByGroupUC           *exerciseusecases.FindTypesByGroupUseCase
	ConfirmDeleteExerciseUC      *exerciseusecases.ConfirmDeleteUseCase
	DeleteExerciseUC             *exerciseusecases.DeleteUseCase
	GetExerciseUC                *exerciseusecases.GetUseCase
	CreateExerciseUC             *exerciseusecases.CreateUseCase

	// timers
	StopTimerUC  *timerusecases.StopUseCase
	StartTimerUC *timerusecases.StartUseCase

	// exports
	ExportToExcelUC *exportusecases.ExportToExcelUseCase

	// stats
	PeriodStatsUC *statsusecases.GetPeriodStatsUseCase

	// sets
	CompleteSetUC   *setusecases.CompleteUseCase
	AddOneMoreSetUC *setusecases.AddOneMoreUseCase
	RemoveLastSetUC *setusecases.RemoveLastUseCase
	UpdateNextSetUC *setusecases.UpdateNextUseCase

	// programs
	DeleteProgramUC         *programusecases.DeleteUseCase
	CreateProgramUC         *programusecases.CreateUseCase
	ActivateProgramUC       *programusecases.ActivateUseCase
	GetProgramUC            *programusecases.GetUseCase
	FindAllProgramsByUserUC *programusecases.FindAllByUserUseCase
	RenameProgramUC         *programusecases.RenameUseCase
	GetByUserProgramUC      *programusecases.GetByUserUseCase

	// dayTypes
	DayTypesCreateUC *daytypeusecases.CreateUseCase
	UpdateDateTypeUC *daytypeusecases.UpdateUseCase
	GetDayTypeUC     *daytypeusecases.GetUseCase

	// groups
	GetGroupUC     *groupusecases.GetUseCase
	GetAllGroupsUC *groupusecases.GetAllUseCase

	// sessions
	MoveSessionToExerciseUC *sessionusecases.MoveToUseCase

	// users
	CreateUserUC    *userusecases.CreateUseCase
	GetUserUC       *userusecases.GetUseCase
	FindUserUC      *userusecases.FindUseCase
	DeleteDayTypeUC *daytypeusecases.DeleteUseCase
}

func NewContainer(db *gorm.DB) *Container {
	usersRepo := users.NewRepo(db)
	programsRepo := programs.NewRepo(db)
	dayTypesRepo := daytypes.NewRepo(db)
	workoutsRepo := workouts.NewRepo(db)
	exercisesRepo := exercises.NewRepo(db)
	setsRepo := sets.NewRepo(db)
	sessionsRepo := sessions.NewRepo(db)
	exerciseTypesRepo := exercisetypes.NewRepo(db)
	exerciseGroupTypesRepo := exercisegrouptypes.NewRepo(db)

	timerStore := timer.NewStore()
	summaryService := summary.NewService()
	docGeneratorService := docgenerator.NewService(summaryService)

	return &Container{

		// workouts
		DeleteWorkoutUC:        workoutusecases.NewDeleteUseCase(workoutsRepo, setsRepo, exercisesRepo),
		ConfirmDeleteWorkoutUC: workoutusecases.NewConfirmDeleteUseCase(workoutsRepo, dayTypesRepo),
		CreateWorkoutUC:        workoutusecases.NewCreateUseCase(workoutsRepo, exercisesRepo, usersRepo, dayTypesRepo),
		StartWorkoutUC:         workoutusecases.NewStartUseCase(workoutsRepo, sessionsRepo),
		FindMyWorkoutsUC:       workoutusecases.NewFindMyUseCase(workoutsRepo, usersRepo),
		ShowWorkoutProgressUC:  workoutusecases.NewShowProgressUseCase(workoutsRepo, sessionsRepo),
		ConfirmFinishWorkoutUC: workoutusecases.NewConfirmFinishUseCase(workoutsRepo, dayTypesRepo),
		FinishWorkoutUC:        workoutusecases.NewFinishUseCase(workoutsRepo, sessionsRepo),
		FindWorkoutsByUserUC:   workoutusecases.NewFindByUserUseCase(workoutsRepo, usersRepo),
		StatsWorkoutUC:         workoutusecases.NewStatsUseCase(workoutsRepo, dayTypesRepo, exerciseTypesRepo),

		// exercises
		ExerciseTypeListUC:      exerciseusecases.NewExerciseTypeListUseCase(exerciseTypesRepo),
		FindTypesByGroupUC:      exerciseusecases.NewFindTypesByGroupUseCase(exerciseTypesRepo),
		ConfirmDeleteExerciseUC: exerciseusecases.NewConfirmDeleteUseCase(exerciseTypesRepo, exercisesRepo),
		DeleteExerciseUC:        exerciseusecases.NewDeleteUseCase(exercisesRepo),
		GetExerciseUC:           exerciseusecases.NewGetUseCase(exercisesRepo, exerciseTypesRepo),
		CreateExerciseUC:        exerciseusecases.NewCreateUseCase(exercisesRepo, workoutsRepo, exerciseTypesRepo),

		// timers
		StopTimerUC:  timerusecases.NewStopUseCase(timerStore),
		StartTimerUC: timerusecases.NewStartUseCase(timerStore, exercisesRepo),

		// exports
		ExportToExcelUC: exportusecases.NewExportToExcelUseCase(usersRepo, exerciseGroupTypesRepo, workoutsRepo,
			exercisesRepo, summaryService, docGeneratorService),

		// stats
		PeriodStatsUC: statsusecases.NewGetPeriodStatsUseCase(usersRepo, workoutsRepo),

		// sets
		CompleteSetUC:   setusecases.NewCompleteUseCase(setsRepo, exercisesRepo, exerciseTypesRepo),
		AddOneMoreSetUC: setusecases.NewAddOneMoreUseCase(setsRepo, exercisesRepo),
		RemoveLastSetUC: setusecases.NewRemoveLastUseCase(setsRepo, exercisesRepo),
		UpdateNextSetUC: setusecases.NewUpdateNextUseCase(setsRepo, exercisesRepo),

		// programs
		DeleteProgramUC:         programusecases.NewDeleteUseCase(programsRepo, usersRepo),
		CreateProgramUC:         programusecases.NewCreateUseCase(programsRepo, usersRepo),
		ActivateProgramUC:       programusecases.NewActivateUseCase(usersRepo),
		GetProgramUC:            programusecases.NewGetUseCase(programsRepo, exerciseTypesRepo),
		FindAllProgramsByUserUC: programusecases.NewFindAllByUserUseCase(programsRepo, usersRepo),
		RenameProgramUC:         programusecases.NewRenameUseCase(programsRepo),
		GetByUserProgramUC:      programusecases.NewGetByUserUseCase(programsRepo, usersRepo),

		// groups
		GetGroupUC:     groupusecases.NewGetUseCase(exerciseGroupTypesRepo),
		GetAllGroupsUC: groupusecases.NewGetAllUseCase(exerciseGroupTypesRepo),

		// dayTypes
		DayTypesCreateUC: daytypeusecases.NewCreateUseCase(dayTypesRepo),
		UpdateDateTypeUC: daytypeusecases.NewUpdateUseCase(dayTypesRepo),
		GetDayTypeUC:     daytypeusecases.NewGetUseCase(dayTypesRepo),
		DeleteDayTypeUC:  daytypeusecases.NewDeleteUseCase(dayTypesRepo),

		// sessions
		ShowCurrentExerciseSessionUC: sessionusecases.NewShowCurrentExerciseUseCase(workoutsRepo, sessionsRepo, exerciseTypesRepo, dayTypesRepo),
		MoveSessionToExerciseUC:      sessionusecases.NewMoveToExerciseUseCase(sessionsRepo, exercisesRepo),

		// users
		CreateUserUC: userusecases.NewCreateUseCase(usersRepo, programsRepo),
		FindUserUC:   userusecases.NewFindUseCase(usersRepo, programsRepo),
		GetUserUC:    userusecases.NewGetUseCase(usersRepo),
	}
}
