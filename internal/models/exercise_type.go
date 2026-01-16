package models

type ExerciseType struct {
	ID                    int64 `gorm:"primaryKey;autoIncrement"`
	Name                  string
	Url                   string
	ExerciseGroupTypeCode string
	RestInSeconds         int
	Accent                string
	Units                 string
}

func (*ExerciseType) TableName() string {
	return "exercise_types"
}
