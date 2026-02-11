import {api} from "./client";

export const getProgram = (programId: number) =>
    api(`/api/programs/${programId}`);

export const createDay = (programId: number, name: string) =>
    api(`/api/programs/${programId}/days`, {
        method: "POST",
        body: JSON.stringify({name}),
    });

export const deleteDay = (programId: number, dayId: number) =>
    api(`/api/programs/${programId}/days/${dayId}`, {method: "DELETE"});
