package templates

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

func GetLegExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: "Разгибание голени сидя (передняя поверхность бедра)",
			Sets: []models.Set{
				{Reps: 16, Weight: 50},
				{Reps: 12, Weight: 60},
				{Reps: 12, Weight: 60},
				{Reps: 12, Weight: 60},
			},
			RestInSeconds: 120,
		},
		{
			Name: "Сгибание голени сидя (задняя поверхность бедра)",
			Sets: []models.Set{
				{Reps: 14, Weight: 40},
				{Reps: 14, Weight: 40},
				{Reps: 14, Weight: 40},
				{Reps: 14, Weight: 40},
			},
			RestInSeconds: 120,
		},
		{
			Name: "Жим платформы ногами (передняя поверхность бедра)",
			Sets: []models.Set{
				{Reps: 17, Weight: 100},
				{Reps: 15, Weight: 160},
				{Reps: 12, Weight: 200},
				{Reps: 12, Weight: 220},
				{Reps: 12, Weight: 240},
				{Reps: 12, Weight: 260},
			},
			RestInSeconds: 180,
		},
		{
			Name: "Подъем ног в вися на локтях (прямая мышца живота)",
			Sets: []models.Set{
				{Reps: 25, Weight: 0},
				{Reps: 25, Weight: 0},
				{Reps: 25, Weight: 0},
			},
			RestInSeconds: 90,
		},
		{
			Name: "Обратные разведения в пек-дек (задняя дельтовидная мышца)",
			Sets: []models.Set{
				{Reps: 15, Weight: 15},
				{Reps: 15, Weight: 15},
				{Reps: 15, Weight: 15},
				{Reps: 15, Weight: 15},
			},
			RestInSeconds: 120,
		},
		{
			Name: "Протяжка штанги (средняя дельтовидная мышца)",
			Sets: []models.Set{
				{Reps: 12, Weight: 40},
				{Reps: 12, Weight: 40},
				{Reps: 12, Weight: 40},
				{Reps: 12, Weight: 40},
			},
			RestInSeconds: 120,
		},
	}
}

func GetBackExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: "Подтягивание в гравитроне широким хватом (широчайшая мышца спины)",
			Sets: []models.Set{
				{Reps: 12, Weight: 14},
				{Reps: 12, Weight: 14},
				{Reps: 12, Weight: 14},
				{Reps: 12, Weight: 14},
			},
			RestInSeconds: 120,
		},
		{
			Name: "Вертикальная тяга в рычажном тренажере (широчайшая мышца спины)",
			Sets: []models.Set{
				{Reps: 10, Weight: 100},
				{Reps: 10, Weight: 100},
				{Reps: 10, Weight: 100},
				{Reps: 10, Weight: 100},
			},
			RestInSeconds: 120,
		},
		{
			Name: "Горизонтальная тяга в блочном тренажере с упором в грудь (широчайшая мышца спины)",
			Sets: []models.Set{
				{Reps: 12, Weight: 60},
				{Reps: 12, Weight: 60},
				{Reps: 12, Weight: 60},
				{Reps: 12, Weight: 60},
			},
			RestInSeconds: 120,
		},
		{
			Name: "Тяга гантели с упором в скамью (широчайшая мышца спины)",
			Sets: []models.Set{
				{Reps: 12, Weight: 20},
				{Reps: 12, Weight: 20},
				{Reps: 12, Weight: 20},
				{Reps: 12, Weight: 20},
			},
			RestInSeconds: 120,
		},
		{
			Name: "Сгибание рук с супинацией гантелями (двуглавая мышца плеча)",
			Sets: []models.Set{
				{Reps: 16, Weight: 14},
				{Reps: 16, Weight: 14},
				{Reps: 16, Weight: 14},
				{Reps: 16, Weight: 14},
			},
			RestInSeconds: 120,
		},
		{
			Name: "Молотковые сгибания с гантелями (брахиалис + плечевая мышца)",
			Sets: []models.Set{
				{Reps: 12, Weight: 14},
				{Reps: 10, Weight: 16},
				{Reps: 8, Weight: 18},
				{Reps: 6, Weight: 20},
			},
			Hint:          "⚠️ Важно для безопасности плеч в супинации:\n- Не размахивайте гантелями в нижней точке \n- Опускайте на 90%, оставляя легкий сгиб в локте \n- При болях в переднем плече - уменьшите амплитуду и вес",
			RestInSeconds: 120,
		},
	}
}
