import {api} from "./client.ts";

export const createMeasurement = (newMeasurement: Partial<ToCreateMeasurement>) =>
    api<Measurement>(`/api/measurements`, {
        method: "POST",
        body: JSON.stringify(newMeasurement),
    });

export const deleteMeasurement = (id: number) =>
    api(`/api/measurements/${id}`, {method: "DELETE"});

export const getMeasurements = (offset, limit: number) =>
    api<FindWithOffsetLimitMeasurement>(`/api/measurements?offset=${offset}&limit=${limit}`);

export const getMeasurementTypes = () =>
    api<MeasurementTypeDTO[]>(`/api/measurements/types`);
