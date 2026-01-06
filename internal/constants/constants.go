package constants

import "fmt"

const (
	LegsAndShouldersWorkoutID   = "legs_and_shoulders"
	LegsAndShouldersWorkoutName = "🦵 Ноги & плечи"

	BackAndBicepsWorkoutID   = "back_and_biceps"
	BackAndBicepsWorkoutName = "🏋️‍♂️ Спина & бицепсы"

	ChestAndTricepsID   = "chest_and_triceps"
	ChestAndTricepsName = "🫀 Грудь & трицепсы"

	CardioID   = "cardio"
	CardioName = "🏃 Кардио"
)

const (
	Legs    = "Ноги"
	Press   = "Пресс"
	Deltas  = "Дельты"
	Back    = "Спина"
	Biceps  = "Бицепс"
	Chest   = "Грудь"
	Triceps = "Трицепс"
)
const (
	FrontSurfaceOfTheThigh = "передняя поверхность бедра"
	BackSurfaceOfTheThigh  = "задняя поверхность бедра"

	RectusAbdominisMuscle = "прямая мышца живота"

	BackDeltoidMuscle   = "задняя дельтовидная мышца"
	MiddleDeltoidMuscle = "средняя дельтовидная мышца"

	LatissimusDorsiMuscle = "широчайшая мышца спины"

	PectoralMuscles = "грудные мышцы"

	BicepsBrachiiMuscle         = "двуглавая мышца плеча"
	TricepsShoulderMuscle       = "трехглавая мышца плеча"
	BrachialisAndshoulderMuscle = "брахиалис + плечевая мышца"
)

type ExerciseObj struct {
	Name   string
	Url    string
	Type   string
	Accent string
}

func (e *ExerciseObj) GetName() string {
	if e == nil {
		return ""
	}
	return e.Name
}

func (e *ExerciseObj) GetAccent() string {
	if e == nil {
		return ""
	}
	return e.Accent
}

func (e *ExerciseObj) GetHint() string {
	if e == nil {
		return ""
	}
	if e.Url == "" {
		return ""
	}
	return WrapYandexLink(e.Url)
}

func WrapYandexLink(url string) string {
	return fmt.Sprintf("<a href=\"%s\">Видео</a>", url)
}

var AllExercises = map[string]*ExerciseObj{
	// Ноги
	ExtensionOfLowerLegWhileSitting: {
		Name:   "Разгибание голени сидя",
		Url:    "https://disk.yandex.ru/i/nevoPFhHbc8l8g",
		Type:   Legs,
		Accent: FrontSurfaceOfTheThigh,
	},
	FlexionOfLowerLegWhileSitting: {
		Name:   "Сгибание голени сидя",
		Url:    "https://disk.yandex.ru/i/PqkWBjSSNvH-Vg",
		Type:   Legs,
		Accent: BackSurfaceOfTheThigh,
	},
	PlatformLegPress: {
		Name:   "Жим платформы ногами",
		Url:    "https://disk.yandex.ru/i/UsaW3YjvDRWm3w",
		Type:   Legs,
		Accent: FrontSurfaceOfTheThigh,
	},
	// Пресс
	LiftingLegsAtTheElbow: {
		Name:   "Подъем ног в висе на локтях",
		Url:    "https://disk.yandex.ru/i/pkAxqVWTe4L_Xw",
		Type:   Press,
		Accent: RectusAbdominisMuscle,
	},
	// Дельты
	ReverseDilutionsInThePectoral: {
		Name:   "Обратные разведения в пек-дек",
		Url:    "https://disk.yandex.ru/i/9lYAV1wr3VjWQQ",
		Type:   Deltas,
		Accent: BackDeltoidMuscle,
	},
	ExtensionOfBarbell: {
		Name:   "Протяжка штанги",
		Url:    "https://disk.yandex.ru/i/0aaEdn5IBOI6zQ",
		Type:   Deltas,
		Accent: MiddleDeltoidMuscle,
	},
	//Спина
	PullUpInTheGravitronWithAWideGrip: {
		Name:   "Подтягивание в гравитроне широким хватом",
		Url:    "https://disk.yandex.ru/i/jp52K-HTe86iLA",
		Type:   Back,
		Accent: LatissimusDorsiMuscle,
	},
	VerticalTractionInALeverSimulator: {
		Name:   "Вертикальная тяга в рычажном тренажере",
		Url:    "https://disk.yandex.ru/i/x6qRCfJBGA7tEQ",
		Type:   Back,
		Accent: LatissimusDorsiMuscle,
	},
	HorizontalDeadliftInABlockSimulatorWithAnEmphasisOnTheChest: {
		Name:   "Горизонтальная тяга в блочном тренажере с упором в грудь",
		Url:    "https://disk.yandex.ru/i/DnyJDcPaJLUyCg",
		Type:   Back,
		Accent: LatissimusDorsiMuscle,
	},
	DumbbellDeadliftWithEmphasisOnTheBench: {
		Name:   "Тяга гантели с упором в скамью",
		Url:    "https://disk.yandex.ru/i/mU9TIaxDPV6nXw",
		Type:   Back,
		Accent: LatissimusDorsiMuscle,
	},
	// Бицепс
	ArmFlexionWithDumbbellSupination: {
		Name:   "Сгибание рук с супинацией гантелями",
		Url:    "https://disk.yandex.ru/i/LWBPrSeWvxNiUw",
		Type:   Biceps,
		Accent: BicepsBrachiiMuscle,
	},
	HammerBendsWithDumbbells: {
		Name:   "Молотковые сгибания с гантелями",
		Url:    "https://disk.yandex.ru/i/OvY5i3YGxyi6gw",
		Type:   Biceps,
		Accent: BrachialisAndshoulderMuscle,
	},
	// Грудные
	BenchPressWithAWideGrip: {
		Name:   "Жим лежа широким хватом",
		Url:    "https://disk.yandex.ru/i/w2FIsYgqMQ-RPA",
		Type:   Chest,
		Accent: PectoralMuscles,
	},
	HorizontalBenchPressInTheTechnoGymSimulator: {
		Name:   "Жим горизонтально в тренажере TechnoGym",
		Url:    "https://disk.yandex.ru/i/vyDhCyusHft5VQ",
		Type:   Chest,
		Accent: PectoralMuscles,
	},
	BringingArmsTogetherInTheButterflySimulator: {
		Name:   "Сведение рук в тренажере бабочка",
		Url:    "https://disk.yandex.ru/i/JADkjm4tiUsAdQ",
		Type:   Chest,
		Accent: PectoralMuscles,
	},
	// Трицепс
	FrenchBenchPressWithDumbbells: {
		Name:   "Французский жим с гантелями лежа",
		Url:    "https://disk.yandex.ru/i/9KPxatabvDYy8g",
		Type:   Triceps,
		Accent: TricepsShoulderMuscle,
	},
	ExtensionOfTricepsFromTheUpperBlockWithARopeHandle: {
		Name:   "Разгибание на трицепс с верхнего блока канатной рукоятью",
		Url:    "https://disk.yandex.ru/i/sG8luvJYQWNgyg",
		Type:   Triceps,
		Accent: TricepsShoulderMuscle,
	},
	// Кардио
	Walking: {
		Name: "Ходьба",
		Url:  "",
		Type: Legs,
	},
	RunningOnTrack: {
		Name: "Бег на дорожке",
		Url:  "",
		Type: Legs,
	},
	RunningOnMechanicalTrack: {
		Name: "Бег на механической дорожке",
		Url:  "",
		Type: Legs,
	},
	Bicycle: {
		Name: "Велосипед",
		Url:  "",
		Type: Legs,
	},
}

const (
	// Ноги
	ExtensionOfLowerLegWhileSitting = "extension_of_lower_leg_while_sitting"
	FlexionOfLowerLegWhileSitting   = "flexion_of_lower_leg_while_sitting"
	PlatformLegPress                = "platform_leg_press"

	// Пресс
	LiftingLegsAtTheElbow = "lifting_legs_at_the_elbow"

	// Дельты
	ReverseDilutionsInThePectoral = "reverse_dilutions_in_the_pectoral"
	ExtensionOfBarbell            = "extension_of_barbell"

	//Спина
	PullUpInTheGravitronWithAWideGrip                           = "pull_up_in_the_gravitron_with_a_wide_grip"
	VerticalTractionInALeverSimulator                           = "vertical_traction_in_a_lever_simulator"
	HorizontalDeadliftInABlockSimulatorWithAnEmphasisOnTheChest = "horizontal_deadlift_in_a_block_simulator_with_an_emphasis_on_the_chest"
	DumbbellDeadliftWithEmphasisOnTheBench                      = "dumbbell_deadlift_with_emphasis_on_the_bench"

	// Бицепс
	ArmFlexionWithDumbbellSupination = "arm_flexion_with_dumbbell_supination"
	HammerBendsWithDumbbells         = "hammer_bends_with_dumbbells"

	// Грудные
	BenchPressWithAWideGrip                     = "bench_press_with_a_wide_grip"
	HorizontalBenchPressInTheTechnoGymSimulator = "horizontal_bench_press_in_the_techno_gym_simulator"
	BringingArmsTogetherInTheButterflySimulator = "bringing_arms_together_in_the_butterfly_simulator"

	// Трицепс
	FrenchBenchPressWithDumbbells                      = "french_bench_press_with_dumbbells"
	ExtensionOfTricepsFromTheUpperBlockWithARopeHandle = "extension_of_triceps_from_the_upper_block_with_a_rope_handle"

	// Кардио
	Walking                  = "walking"
	RunningOnTrack           = "running_on_track"
	RunningOnMechanicalTrack = "running_on_mechanical_track"
	Bicycle                  = "bicycle"
)
