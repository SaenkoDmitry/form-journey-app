import {api} from "./client";

export const getPrograms = () =>
    api<ProgramDTO[]>("/api/programs");

export const GetActiveProgramForUser = () =>
    api<ProgramDTO[]>("/api/programs/active");

export const createProgram = (name: string) =>
    api("/api/programs", {
        method: "POST",
        body: JSON.stringify({name}),
    });

export const deleteProgram = (id: number) =>
    api(`/api/programs/${id}`, {method: "DELETE"});

export const renameProgram = (id: number, name: string) =>
    api(`/api/programs/${id}/rename`, {
        method: "POST",
        body: JSON.stringify({name}),
    });

export const chooseProgram = (id: number) =>
    api(`/api/programs/${id}/choose`, {method: "POST"});
