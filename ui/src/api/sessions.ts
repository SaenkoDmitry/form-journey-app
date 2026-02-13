import {api} from "./client.ts";

export const showCurrentExerciseSession = (workoutID: number) =>
    api<CurrentExerciseSession>(`/api/sessions/${workoutID}`, {
        method: "GET",
    });

export const moveToExerciseSession = (workoutID: number, next: boolean) =>
    api<{}>(`/api/sessions/${workoutID}`, {
        method: "POST",
        body: JSON.stringify({
            next: next,
        }),
    });
