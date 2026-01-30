package summary

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	"math"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

type Service interface {
	BuildTotal(workouts []models.WorkoutDay, groupCodesMap map[string]string) map[string]*ExerciseSummary
	BuildByDate(workouts []models.WorkoutDay) map[string]*DateSummary
	BuildExerciseProgress(workouts []models.WorkoutDay, exerciseName string) map[string]*Progress
}

type serviceImpl struct {
}

func NewService() Service {
	return &serviceImpl{}
}

func (s *serviceImpl) BuildTotal(workouts []models.WorkoutDay, groupCodesMap map[string]string) map[string]*ExerciseSummary {
	summary := make(map[string]*ExerciseSummary)

	for _, w := range workouts {
		if !w.Completed {
			continue
		}
		date := w.StartedAt.Add(3 * time.Hour).Format("2006-01-02")

		for _, e := range w.Exercises {
			if e.CompletedSets() == 0 {
				continue
			}

			sum, ok := summary[e.ExerciseType.Name]
			if !ok {
				sum = &ExerciseSummary{
					Workouts: make(map[string]struct{}),
				}
				summary[e.ExerciseType.Name] = sum
			}

			sum.ExerciseType = groupCodesMap[e.ExerciseType.ExerciseGroupTypeCode]
			sum.Workouts[date] = struct{}{}

			for _, set := range e.Sets {
				if !set.Completed {
					continue
				}
				sum.Sets++
				sum.TotalWeight += float64(set.Weight) * float64(set.Reps)
				sum.TotalReps += set.Reps
				sum.TotalMinutes += set.Minutes

				if set.Weight > sum.MaxWeight {
					sum.MaxWeight = set.Weight
				}
			}
		}
	}

	for _, sum := range summary {
		if sum.TotalReps > 0 {
			sum.AvgWeight = math.Round(sum.TotalWeight / float64(sum.TotalReps))
		}
	}

	return summary
}

func (s *serviceImpl) BuildByDate(workouts []models.WorkoutDay) map[string]*DateSummary {
	result := make(map[string]*DateSummary)

	for _, w := range workouts {
		date := w.StartedAt.Add(3 * time.Hour).Format("2006-01-02")

		d, ok := result[date]
		if !ok {
			d = &DateSummary{
				Workouts:  1,
				Exercises: make(map[string]struct{}),
			}
			result[date] = d
		}

		for _, e := range w.Exercises {
			d.Exercises[e.ExerciseType.Name] = struct{}{}

			for _, sum := range e.Sets {
				d.Sets++
				d.TotalVolume += sum.Weight * float32(sum.Reps)

				if sum.Weight > d.MaxWeight {
					d.MaxWeight = sum.Weight
				}
			}
		}
	}

	return result
}

func (s *serviceImpl) BuildExerciseProgress(
	workouts []models.WorkoutDay,
	exerciseName string,
) map[string]*Progress {

	progress := make(map[string]*Progress)

	for _, w := range workouts {
		if !w.Completed {
			continue
		}

		currDate := w.StartedAt.Add(3 * time.Hour).Format("02.01.2006")
		thisWeek := utils.GetThisWeek(w.StartedAt)

		var key string
		for _, e := range w.Exercises {
			if e.ExerciseType.Name != exerciseName {
				continue
			}
			if e.CompletedSets() == 0 {
				continue
			}

			if e.ExerciseType.ExerciseGroupTypeCode == "cardio" {
				key = thisWeek
			} else {
				key = currDate
			}

			sumWeight := float32(0)
			countOfReps := 0

			if _, ok := progress[key]; !ok {
				progress[key] = &Progress{}
			}

			for _, set := range e.Sets {
				if !set.Completed {
					continue
				}

				countOfReps += set.GetRealReps()
				sumWeight += set.GetRealWeight() * float32(set.GetRealReps())
				if progress[key].MaxWeight < set.GetRealWeight() ||
					progress[key].MaxWeight == set.GetRealWeight() && progress[key].MaxReps < set.GetRealReps() {
					progress[key].MaxWeight = set.GetRealWeight()
					progress[key].MaxReps = set.GetRealReps()
				}

				progress[key].SumMinutes += set.GetRealMinutes()
				if progress[key].MinMinutes == 0 {
					progress[key].MinMinutes = set.GetRealMinutes()
					progress[key].MaxMinutes = set.GetRealMinutes()
				} else {
					progress[key].MinMinutes = min(progress[key].MinMinutes, set.GetRealMinutes())
					progress[key].MaxMinutes = max(progress[key].MaxMinutes, set.GetRealMinutes())
				}

				progress[key].SumMeters += set.GetRealMeters()
				if progress[key].MinMeters == 0 {
					progress[key].MinMeters = set.GetRealMeters()
					progress[key].MaxMeters = set.GetRealMeters()
				} else {
					progress[key].MinMeters = min(progress[key].MinMeters, set.GetRealMeters())
					progress[key].MaxMeters = max(progress[key].MaxMeters, set.GetRealMeters())
				}
			}
			progress[key].AvgWeight = sumWeight / float32(countOfReps)
		}
	}

	return progress
}
