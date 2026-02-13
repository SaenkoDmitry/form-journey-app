import {api} from "./client.ts";

export const deleteWorkout = (id: number) =>
    api(`/api/workouts/${id}`, {method: "DELETE"});

export const getWorkouts = (offset, limit: number) =>
    api<ShowMyWorkoutsResult>(`/api/workouts?offset=${offset}&limit=${limit}`);

export const startWorkout = (dayTypeID: number) =>
    api<StartWorkoutDTO>("/api/workouts/start", {
        method: "POST",
        body: JSON.stringify({
            day_type_id: dayTypeID,
        }),
    });

export const finishWorkout = (workoutID: number) =>
    api<StartWorkoutDTO>(`/api/workouts/${workoutID}/finish`, {
        method: "POST",
    });
