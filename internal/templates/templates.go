package templates

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

func GetLegExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: constants.ExtensionOfLowerLegWhileSitting,
			Sets: []models.Set{
				{Reps: 16, Weight: 50, Index: 1},
				{Reps: 12, Weight: 60, Index: 2},
				{Reps: 12, Weight: 60, Index: 3},
				{Reps: 12, Weight: 60, Index: 4},
			},
			RestInSeconds: 120,
		},
		{
			Name: constants.FlexionOfLowerLegWhileSitting,
			Sets: []models.Set{
				{Reps: 14, Weight: 40, Index: 1},
				{Reps: 14, Weight: 40, Index: 2},
				{Reps: 14, Weight: 40, Index: 3},
				{Reps: 14, Weight: 40, Index: 4},
			},
			RestInSeconds: 120,
		},
		{
			Name: constants.PlatformLegPress,
			Sets: []models.Set{
				{Reps: 17, Weight: 100, Index: 1},
				{Reps: 15, Weight: 160, Index: 2},
				{Reps: 12, Weight: 200, Index: 3},
				{Reps: 12, Weight: 220, Index: 4},
				{Reps: 12, Weight: 240, Index: 5},
				{Reps: 12, Weight: 260, Index: 6},
			},
			RestInSeconds: 180,
		},
		{
			Name: constants.LiftingLegsAtTheElbow,
			Sets: []models.Set{
				{Reps: 25, Weight: 0, Index: 1},
				{Reps: 25, Weight: 0, Index: 2},
				{Reps: 25, Weight: 0, Index: 3},
			},
			RestInSeconds: 90,
		},
	}
}

func GetShoulderExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: constants.ReverseDilutionsInThePectoral,
			Sets: []models.Set{
				{Reps: 15, Weight: 15, Index: 1},
				{Reps: 15, Weight: 15, Index: 2},
				{Reps: 15, Weight: 15, Index: 3},
				{Reps: 15, Weight: 15, Index: 4},
			},
			RestInSeconds: 120,
		},
		{
			Name: constants.ExtensionOfBarbell,
			Sets: []models.Set{
				{Reps: 12, Weight: 40, Index: 1},
				{Reps: 12, Weight: 40, Index: 2},
				{Reps: 12, Weight: 40, Index: 3},
				{Reps: 12, Weight: 40, Index: 4},
			},
			RestInSeconds: 120,
		},
	}
}

func GetBackExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: constants.PullUpInTheGravitronWithAWideGrip,
			Sets: []models.Set{
				{Reps: 12, Weight: 14, Index: 1},
				{Reps: 12, Weight: 14, Index: 2},
				{Reps: 12, Weight: 14, Index: 3},
				{Reps: 12, Weight: 14, Index: 4},
			},
			RestInSeconds: 120,
		},
		{
			Name: constants.VerticalTractionInALeverSimulator,
			Sets: []models.Set{
				{Reps: 10, Weight: 100, Index: 1},
				{Reps: 10, Weight: 100, Index: 2},
				{Reps: 10, Weight: 100, Index: 3},
				{Reps: 10, Weight: 100, Index: 4},
			},
			RestInSeconds: 120,
		},
		{
			Name: constants.HorizontalDeadliftInABlockSimulatorWithAnEmphasisOnTheChest,
			Sets: []models.Set{
				{Reps: 12, Weight: 60, Index: 1},
				{Reps: 12, Weight: 60, Index: 2},
				{Reps: 12, Weight: 60, Index: 3},
				{Reps: 12, Weight: 60, Index: 4},
			},
			RestInSeconds: 120,
		},
		{
			Name: constants.DumbbellDeadliftWithEmphasisOnTheBench,
			Sets: []models.Set{
				{Reps: 12, Weight: 20, Index: 1},
				{Reps: 12, Weight: 20, Index: 2},
				{Reps: 12, Weight: 20, Index: 3},
				{Reps: 12, Weight: 20, Index: 4},
			},
			RestInSeconds: 120,
		},
	}
}

func GetBicepsExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: constants.ArmFlexionWithDumbbellSupination,
			Sets: []models.Set{
				{Reps: 14, Weight: 15, Index: 1},
				{Reps: 14, Weight: 15, Index: 2},
				{Reps: 14, Weight: 15, Index: 3},
				{Reps: 14, Weight: 15, Index: 4},
			},
			RestInSeconds: 120,
		},
		{
			Name: constants.HammerBendsWithDumbbells,
			Sets: []models.Set{
				{Reps: 12, Weight: 14, Index: 1},
				{Reps: 10, Weight: 16, Index: 2},
				{Reps: 8, Weight: 18, Index: 3},
				{Reps: 6, Weight: 20, Index: 4},
			},
			RestInSeconds: 120,
		},
	}
}

func GetChestExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: constants.BenchPressWithAWideGrip,
			Sets: []models.Set{
				{Reps: 16, Weight: 45, Index: 1},
				{Reps: 15, Weight: 55, Index: 2},
				{Reps: 14, Weight: 65, Index: 3},
				{Reps: 14, Weight: 65, Index: 4},
				{Reps: 14, Weight: 65, Index: 5},
			},
			RestInSeconds: 180,
		},
		{
			Name: constants.HorizontalBenchPressInTheTechnoGymSimulator,
			Sets: []models.Set{
				{Reps: 12, Weight: 60, Index: 1},
				{Reps: 12, Weight: 60, Index: 2},
				{Reps: 12, Weight: 60, Index: 3},
				{Reps: 12, Weight: 60, Index: 4},
			},
			RestInSeconds: 120,
		},
		{
			Name: constants.BringingArmsTogetherInTheButterflySimulator,
			Sets: []models.Set{
				{Reps: 14, Weight: 17, Index: 1},
				{Reps: 14, Weight: 17, Index: 2},
				{Reps: 14, Weight: 17, Index: 3},
				{Reps: 14, Weight: 17, Index: 4},
			},
			RestInSeconds: 120,
		},
	}
}

func GetTricepsExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: constants.FrenchBenchPressWithDumbbells,
			Sets: []models.Set{
				{Reps: 14, Weight: 16, Index: 1},
				{Reps: 14, Weight: 16, Index: 2},
				{Reps: 14, Weight: 16, Index: 3},
			},
			RestInSeconds: 120,
		},
		{
			Name: constants.ExtensionOfTricepsFromTheUpperBlockWithARopeHandle,
			Sets: []models.Set{
				{Reps: 12, Weight: 17, Index: 1},
				{Reps: 12, Weight: 17, Index: 2},
				{Reps: 12, Weight: 17, Index: 3},
			},
			RestInSeconds: 120,
		},
	}
}

func GetCardioExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: constants.Walking,
			Sets: []models.Set{
				{Minutes: 15, Index: 1},
			},
		},
		{
			Name: constants.RunningOnTrack,
			Sets: []models.Set{
				{Minutes: 15, Index: 1},
			},
		},
	}
}
