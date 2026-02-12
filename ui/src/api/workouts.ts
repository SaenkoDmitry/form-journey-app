import {api} from "./client.ts";

export const deleteWorkout = (id: number) =>
    api(`/api/workouts/${id}`, {method: "DELETE"});
