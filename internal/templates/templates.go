package templates

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

func GetLegExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: "Разгибание голени сидя (передняя поверхность бедра)",
			Sets: []models.Set{
				{Reps: 16, Weight: 50, Index: 1},
				{Reps: 12, Weight: 60, Index: 2},
				{Reps: 12, Weight: 60, Index: 3},
				{Reps: 12, Weight: 60, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1O5ZtpBUpuromec5ISnmbi1F6JxPCZc7y/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Сгибание голени сидя (задняя поверхность бедра)",
			Sets: []models.Set{
				{Reps: 14, Weight: 40, Index: 1},
				{Reps: 14, Weight: 40, Index: 2},
				{Reps: 14, Weight: 40, Index: 3},
				{Reps: 14, Weight: 40, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1YuwDtXx2ITjCqzjIldwp3NxH_lIvLz6f/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Жим платформы ногами (передняя поверхность бедра)",
			Sets: []models.Set{
				{Reps: 17, Weight: 100, Index: 1},
				{Reps: 15, Weight: 160, Index: 2},
				{Reps: 12, Weight: 200, Index: 3},
				{Reps: 12, Weight: 220, Index: 4},
				{Reps: 12, Weight: 240, Index: 5},
				{Reps: 12, Weight: 260, Index: 6},
			},
			RestInSeconds: 180,
			Hint:          "<a href=\"https://drive.google.com/file/d/1K56NqY-QwpgAMBN1BZk4l1BwsOBqfsFd/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Подъем ног в висе на локтях (прямая мышца живота)",
			Sets: []models.Set{
				{Reps: 25, Weight: 0, Index: 1},
				{Reps: 25, Weight: 0, Index: 2},
				{Reps: 25, Weight: 0, Index: 3},
			},
			RestInSeconds: 90,
			Hint:          "<a href=\"https://drive.google.com/file/d/1zRS_sbKBZr6LDLqtQwpnn7zO00pW1f2M/view?usp=drive_link\">Google drive</a>",
		},
	}
}

func GetShoulderExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: "Обратные разведения в пек-дек (задняя дельтовидная мышца)",
			Sets: []models.Set{
				{Reps: 15, Weight: 15, Index: 1},
				{Reps: 15, Weight: 15, Index: 2},
				{Reps: 15, Weight: 15, Index: 3},
				{Reps: 15, Weight: 15, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1gf78lwsJ8bLjbM8ib05_LwVJYNlu7dR5/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Протяжка штанги (средняя дельтовидная мышца)",
			Sets: []models.Set{
				{Reps: 12, Weight: 40, Index: 1},
				{Reps: 12, Weight: 40, Index: 2},
				{Reps: 12, Weight: 40, Index: 3},
				{Reps: 12, Weight: 40, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1GJ687cZsaQqH4CWB8vZHCXnpwYd7XoAi/view?usp=drive_link\">Google drive</a>",
		},
	}
}

func GetBackExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: "Подтягивание в гравитроне широким хватом (широчайшая мышца спины)",
			Sets: []models.Set{
				{Reps: 12, Weight: 14, Index: 1},
				{Reps: 12, Weight: 14, Index: 2},
				{Reps: 12, Weight: 14, Index: 3},
				{Reps: 12, Weight: 14, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1PD8_FusA1mHskK0NI4m1F3hwRUaGvrgE/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Вертикальная тяга в рычажном тренажере (широчайшая мышца спины)",
			Sets: []models.Set{
				{Reps: 10, Weight: 100, Index: 1},
				{Reps: 10, Weight: 100, Index: 2},
				{Reps: 10, Weight: 100, Index: 3},
				{Reps: 10, Weight: 100, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1bYdfjJMWW0hmLsf3ExpNuQ-0xNMRS1U6/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Горизонтальная тяга в блочном тренажере с упором в грудь (широчайшая мышца спины)",
			Sets: []models.Set{
				{Reps: 12, Weight: 60, Index: 1},
				{Reps: 12, Weight: 60, Index: 2},
				{Reps: 12, Weight: 60, Index: 3},
				{Reps: 12, Weight: 60, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1fF0cWdCwWDvNRXFgdT5tmwtE9kn7KyQF/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Тяга гантели с упором в скамью (широчайшая мышца спины)",
			Sets: []models.Set{
				{Reps: 12, Weight: 20, Index: 1},
				{Reps: 12, Weight: 20, Index: 2},
				{Reps: 12, Weight: 20, Index: 3},
				{Reps: 12, Weight: 20, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/14GX4r7yNO2vyQda9YJoTzQVxwCXuZkz3/view?usp=drive_link\">Google drive</a>",
		},
	}
}

func GetBicepsExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: "Сгибание рук с супинацией гантелями (двуглавая мышца плеча)",
			Sets: []models.Set{
				{Reps: 14, Weight: 15, Index: 1},
				{Reps: 14, Weight: 15, Index: 2},
				{Reps: 14, Weight: 15, Index: 3},
				{Reps: 14, Weight: 15, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1rBaFPefQgB0wcC5t7uPvMMKFq_LlBvnT/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Молотковые сгибания с гантелями (брахиалис + плечевая мышца)",
			Sets: []models.Set{
				{Reps: 12, Weight: 14, Index: 1},
				{Reps: 10, Weight: 16, Index: 2},
				{Reps: 8, Weight: 18, Index: 3},
				{Reps: 6, Weight: 20, Index: 4},
			},
			Hint: `<a href="https://drive.google.com/file/d/1Z_U7XNG_uzgGetLuYlXKaV6DmuBeJ2Q9/view?usp=drive_link">Google drive</a>

<b>Важно для безопасности плеч в супинации:</b>
			- Не размахивайте гантелями в нижней точке
			- Опускайте на 90%, оставляя легкий сгиб в локте
			- При болях в переднем плече - уменьшите амплитуду и вес`,
			RestInSeconds: 120,
		},
	}
}

func GetChestExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: "Жим лежа широким хватом (грудные мышцы)",
			Sets: []models.Set{
				{Reps: 16, Weight: 45, Index: 1},
				{Reps: 15, Weight: 55, Index: 2},
				{Reps: 14, Weight: 65, Index: 3},
				{Reps: 14, Weight: 65, Index: 4},
				{Reps: 14, Weight: 65, Index: 5},
			},
			RestInSeconds: 180,
			Hint:          "<a href=\"https://drive.google.com/file/d/14UrwIH5SsuFi1HHk0jVjrx8QTl89dgWU/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Жим горизонтально в тренажере Technogym (грудные мышцы)",
			Sets: []models.Set{
				{Reps: 12, Weight: 60, Index: 1},
				{Reps: 12, Weight: 60, Index: 2},
				{Reps: 12, Weight: 60, Index: 3},
				{Reps: 12, Weight: 60, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1cW6OCH1d7Q9T7Qkb-9Ipi3o11WWD-WGa/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Сведение рук в тренажере бабочка (грудные мышцы)",
			Sets: []models.Set{
				{Reps: 14, Weight: 17, Index: 1},
				{Reps: 14, Weight: 17, Index: 2},
				{Reps: 14, Weight: 17, Index: 3},
				{Reps: 14, Weight: 17, Index: 4},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1ig_qeLClNbP6RgZLMoHf8egzIyKjGGWy/view?usp=drive_link\">Google drive</a>",
		},
	}
}

func GetTricepsExercises() []models.Exercise {
	return []models.Exercise{
		{
			Name: "Французский жим с гантелями лежа (трехглавая мышца плеча / трицепс)",
			Sets: []models.Set{
				{Reps: 14, Weight: 16, Index: 1},
				{Reps: 14, Weight: 16, Index: 2},
				{Reps: 14, Weight: 16, Index: 3},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/173bvlP-5G1R_xM0f5TCGHNgaaWyDotej/view?usp=drive_link\">Google drive</a>",
		},
		{
			Name: "Разгибание на трицепс с верхнего блока канатной рукоятью (трехглавая мышца плеча / трицепс)",
			Sets: []models.Set{
				{Reps: 12, Weight: 17, Index: 1},
				{Reps: 12, Weight: 17, Index: 2},
				{Reps: 12, Weight: 17, Index: 3},
			},
			RestInSeconds: 120,
			Hint:          "<a href=\"https://drive.google.com/file/d/1WcDsUYztez0jcwoyoaJr600fRv0DdOCO/view?usp=drive_link\">Google drive</a>",
		},
	}
}
