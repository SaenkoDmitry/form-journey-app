package templates

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

const (
	ExtensionOfLowerLegWhileSitting                             = "Разгибание голени сидя (передняя поверхность бедра)"
	FlexionOfLowerLegWhileSitting                               = "Сгибание голени сидя (задняя поверхность бедра)"
	PlatformLegPress                                            = "Жим платформы ногами (передняя поверхность бедра)"
	LiftingLegsAtTheElbow                                       = "Подъем ног в висе на локтях (прямая мышца живота)"
	ReverseDilutionsInThePectoral                               = "Обратные разведения в пек-дек (задняя дельтовидная мышца)"
	ExtensionOfBarbell                                          = "Протяжка штанги (средняя дельтовидная мышца)"
	PullUpInTheGravitronWithAWideGrip                           = "Подтягивание в гравитроне широким хватом (широчайшая мышца спины)"
	VerticalTractionInALeverSimulator                           = "Вертикальная тяга в рычажном тренажере (широчайшая мышца спины)"
	HorizontalDeadliftInABlockSimulatorWithAnEmphasisOnTheChest = "Горизонтальная тяга в блочном тренажере с упором в грудь (широчайшая мышца спины)"
	DumbbellDeadliftWithEmphasisOnTheBench                      = "Тяга гантели с упором в скамью (широчайшая мышца спины)"
	ArmFlexionWithDumbbellSupination                            = "Сгибание рук с супинацией гантелями (двуглавая мышца плеча)"
	HammerBendsWithDumbbells                                    = "Молотковые сгибания с гантелями (брахиалис + плечевая мышца)"
	BenchPressWithAWideGrip                                     = "Жим лежа широким хватом (грудные мышцы)"
	HorizontalBenchPressInTheTechnoGymSimulator                 = "Жим горизонтально в тренажере TechnoGym (грудные мышцы)"
	BringingArmsTogetherInTheButterflySimulator                 = "Сведение рук в тренажере бабочка (грудные мышцы)"
	FrenchBenchPressWithDumbbells                               = "Французский жим с гантелями лежа (трехглавая мышца плеча / трицепс)"
	ExtensionOfTricepsFromTheUpperBlockWithARopeHandle          = "Разгибание на трицепс с верхнего блока канатной рукоятью (трехглавая мышца плеча / трицепс)"
)

var (
	Hints = map[string]string{
		ExtensionOfLowerLegWhileSitting:                             "<a href=\"https://drive.google.com/file/d/1O5ZtpBUpuromec5ISnmbi1F6JxPCZc7y/view?usp=drive_link\">Google drive</a>",
		FlexionOfLowerLegWhileSitting:                               "<a href=\"https://drive.google.com/file/d/1YuwDtXx2ITjCqzjIldwp3NxH_lIvLz6f/view?usp=drive_link\">Google drive</a>",
		PlatformLegPress:                                            "<a href=\"https://drive.google.com/file/d/1K56NqY-QwpgAMBN1BZk4l1BwsOBqfsFd/view?usp=drive_link\">Google drive</a>",
		LiftingLegsAtTheElbow:                                       "<a href=\"https://drive.google.com/file/d/1zRS_sbKBZr6LDLqtQwpnn7zO00pW1f2M/view?usp=drive_link\">Google drive</a>",
		ReverseDilutionsInThePectoral:                               "<a href=\"https://drive.google.com/file/d/1gf78lwsJ8bLjbM8ib05_LwVJYNlu7dR5/view?usp=drive_link\">Google drive</a>",
		ExtensionOfBarbell:                                          "<a href=\"https://drive.google.com/file/d/1GJ687cZsaQqH4CWB8vZHCXnpwYd7XoAi/view?usp=drive_link\">Google drive</a>",
		PullUpInTheGravitronWithAWideGrip:                           "<a href=\"https://drive.google.com/file/d/1PD8_FusA1mHskK0NI4m1F3hwRUaGvrgE/view?usp=drive_link\">Google drive</a>",
		VerticalTractionInALeverSimulator:                           "<a href=\"https://drive.google.com/file/d/1bYdfjJMWW0hmLsf3ExpNuQ-0xNMRS1U6/view?usp=drive_link\">Google drive</a>",
		HorizontalDeadliftInABlockSimulatorWithAnEmphasisOnTheChest: "<a href=\"https://drive.google.com/file/d/1fF0cWdCwWDvNRXFgdT5tmwtE9kn7KyQF/view?usp=drive_link\">Google drive</a>",
		DumbbellDeadliftWithEmphasisOnTheBench:                      "<a href=\"https://drive.google.com/file/d/14GX4r7yNO2vyQda9YJoTzQVxwCXuZkz3/view?usp=drive_link\">Google drive</a>",
		ArmFlexionWithDumbbellSupination:                            "<a href=\"https://drive.google.com/file/d/1rBaFPefQgB0wcC5t7uPvMMKFq_LlBvnT/view?usp=drive_link\">Google drive</a>",
		HammerBendsWithDumbbells: `<a href="https://drive.google.com/file/d/1Z_U7XNG_uzgGetLuYlXKaV6DmuBeJ2Q9/view?usp=drive_link">Google drive</a>

<b>Важно для безопасности плеч в супинации:</b>
			- Не размахивайте гантелями в нижней точке
			- Опускайте на 90%, оставляя легкий сгиб в локте
			- При болях в переднем плече - уменьшите амплитуду и вес`,
		BenchPressWithAWideGrip:                            "<a href=\"https://drive.google.com/file/d/14UrwIH5SsuFi1HHk0jVjrx8QTl89dgWU/view?usp=drive_link\">Google drive</a>",
		HorizontalBenchPressInTheTechnoGymSimulator:        "<a href=\"https://drive.google.com/file/d/1cW6OCH1d7Q9T7Qkb-9Ipi3o11WWD-WGa/view?usp=drive_link\">Google drive</a>",
		BringingArmsTogetherInTheButterflySimulator:        "<a href=\"https://drive.google.com/file/d/1ig_qeLClNbP6RgZLMoHf8egzIyKjGGWy/view?usp=drive_link\">Google drive</a>",
		FrenchBenchPressWithDumbbells:                      "<a href=\"https://drive.google.com/file/d/173bvlP-5G1R_xM0f5TCGHNgaaWyDotej/view?usp=drive_link\">Google drive</a>",
		ExtensionOfTricepsFromTheUpperBlockWithARopeHandle: "<a href=\"https://drive.google.com/file/d/1WcDsUYztez0jcwoyoaJr600fRv0DdOCO/view?usp=drive_link\">Google drive</a>",
	}
)

func GetLegExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: ExtensionOfLowerLegWhileSitting,
			Sets: []models.Set{
				{Reps: 16, Weight: 50, Index: 1},
				{Reps: 12, Weight: 60, Index: 2},
				{Reps: 12, Weight: 60, Index: 3},
				{Reps: 12, Weight: 60, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[ExtensionOfLowerLegWhileSitting],
		},
		{
			Name: FlexionOfLowerLegWhileSitting,
			Sets: []models.Set{
				{Reps: 14, Weight: 40, Index: 1},
				{Reps: 14, Weight: 40, Index: 2},
				{Reps: 14, Weight: 40, Index: 3},
				{Reps: 14, Weight: 40, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[FlexionOfLowerLegWhileSitting],
		},
		{
			Name: PlatformLegPress,
			Sets: []models.Set{
				{Reps: 17, Weight: 100, Index: 1},
				{Reps: 15, Weight: 160, Index: 2},
				{Reps: 12, Weight: 200, Index: 3},
				{Reps: 12, Weight: 220, Index: 4},
				{Reps: 12, Weight: 240, Index: 5},
				{Reps: 12, Weight: 260, Index: 6},
			},
			RestInSeconds: 180,
			Hint:          Hints[PlatformLegPress],
		},
		{
			Name: LiftingLegsAtTheElbow,
			Sets: []models.Set{
				{Reps: 25, Weight: 0, Index: 1},
				{Reps: 25, Weight: 0, Index: 2},
				{Reps: 25, Weight: 0, Index: 3},
			},
			RestInSeconds: 90,
			Hint:          Hints[LiftingLegsAtTheElbow],
		},
	}
}

func GetShoulderExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: ReverseDilutionsInThePectoral,
			Sets: []models.Set{
				{Reps: 15, Weight: 15, Index: 1},
				{Reps: 15, Weight: 15, Index: 2},
				{Reps: 15, Weight: 15, Index: 3},
				{Reps: 15, Weight: 15, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[ReverseDilutionsInThePectoral],
		},
		{
			Name: ExtensionOfBarbell,
			Sets: []models.Set{
				{Reps: 12, Weight: 40, Index: 1},
				{Reps: 12, Weight: 40, Index: 2},
				{Reps: 12, Weight: 40, Index: 3},
				{Reps: 12, Weight: 40, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[ExtensionOfBarbell],
		},
	}
}

func GetBackExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: PullUpInTheGravitronWithAWideGrip,
			Sets: []models.Set{
				{Reps: 12, Weight: 14, Index: 1},
				{Reps: 12, Weight: 14, Index: 2},
				{Reps: 12, Weight: 14, Index: 3},
				{Reps: 12, Weight: 14, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[PullUpInTheGravitronWithAWideGrip],
		},
		{
			Name: VerticalTractionInALeverSimulator,
			Sets: []models.Set{
				{Reps: 10, Weight: 100, Index: 1},
				{Reps: 10, Weight: 100, Index: 2},
				{Reps: 10, Weight: 100, Index: 3},
				{Reps: 10, Weight: 100, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[VerticalTractionInALeverSimulator],
		},
		{
			Name: HorizontalDeadliftInABlockSimulatorWithAnEmphasisOnTheChest,
			Sets: []models.Set{
				{Reps: 12, Weight: 60, Index: 1},
				{Reps: 12, Weight: 60, Index: 2},
				{Reps: 12, Weight: 60, Index: 3},
				{Reps: 12, Weight: 60, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[HorizontalDeadliftInABlockSimulatorWithAnEmphasisOnTheChest],
		},
		{
			Name: DumbbellDeadliftWithEmphasisOnTheBench,
			Sets: []models.Set{
				{Reps: 12, Weight: 20, Index: 1},
				{Reps: 12, Weight: 20, Index: 2},
				{Reps: 12, Weight: 20, Index: 3},
				{Reps: 12, Weight: 20, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[DumbbellDeadliftWithEmphasisOnTheBench],
		},
	}
}

func GetBicepsExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: ArmFlexionWithDumbbellSupination,
			Sets: []models.Set{
				{Reps: 14, Weight: 15, Index: 1},
				{Reps: 14, Weight: 15, Index: 2},
				{Reps: 14, Weight: 15, Index: 3},
				{Reps: 14, Weight: 15, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[ArmFlexionWithDumbbellSupination],
		},
		{
			Name: HammerBendsWithDumbbells,
			Sets: []models.Set{
				{Reps: 12, Weight: 14, Index: 1},
				{Reps: 10, Weight: 16, Index: 2},
				{Reps: 8, Weight: 18, Index: 3},
				{Reps: 6, Weight: 20, Index: 4},
			},
			Hint:          Hints[HammerBendsWithDumbbells],
			RestInSeconds: 120,
		},
	}
}

func GetChestExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: BenchPressWithAWideGrip,
			Sets: []models.Set{
				{Reps: 16, Weight: 45, Index: 1},
				{Reps: 15, Weight: 55, Index: 2},
				{Reps: 14, Weight: 65, Index: 3},
				{Reps: 14, Weight: 65, Index: 4},
				{Reps: 14, Weight: 65, Index: 5},
			},
			RestInSeconds: 180,
			Hint:          Hints[BenchPressWithAWideGrip],
		},
		{
			Name: HorizontalBenchPressInTheTechnoGymSimulator,
			Sets: []models.Set{
				{Reps: 12, Weight: 60, Index: 1},
				{Reps: 12, Weight: 60, Index: 2},
				{Reps: 12, Weight: 60, Index: 3},
				{Reps: 12, Weight: 60, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[HorizontalBenchPressInTheTechnoGymSimulator],
		},
		{
			Name: BringingArmsTogetherInTheButterflySimulator,
			Sets: []models.Set{
				{Reps: 14, Weight: 17, Index: 1},
				{Reps: 14, Weight: 17, Index: 2},
				{Reps: 14, Weight: 17, Index: 3},
				{Reps: 14, Weight: 17, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          Hints[BringingArmsTogetherInTheButterflySimulator],
		},
	}
}

func GetTricepsExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: FrenchBenchPressWithDumbbells,
			Sets: []models.Set{
				{Reps: 14, Weight: 16, Index: 1},
				{Reps: 14, Weight: 16, Index: 2},
				{Reps: 14, Weight: 16, Index: 3},
			},
			RestInSeconds: 120,
			Hint:          Hints[FrenchBenchPressWithDumbbells],
		},
		{
			Name: ExtensionOfTricepsFromTheUpperBlockWithARopeHandle,
			Sets: []models.Set{
				{Reps: 12, Weight: 17, Index: 1},
				{Reps: 12, Weight: 17, Index: 2},
				{Reps: 12, Weight: 17, Index: 3},
			},
			RestInSeconds: 120,
			Hint:          Hints[ExtensionOfTricepsFromTheUpperBlockWithARopeHandle],
		},
	}
}
