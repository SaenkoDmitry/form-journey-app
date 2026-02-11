import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { createDay, deleteDay, getProgram } from "../api/days";
import Button from "../components/Button";
import DayCard from "../components/DayCard";
import "../styles/ProgramBase.css";

export type WorkoutDayTypeDTO = {
    id: number;
    program_id: number;
    name: string;
    preset?: string;
    created_at: string;
};

export type ProgramDTO = {
    id: number;
    user_id: number;
    name: string;
    created_at: string;
    day_types: WorkoutDayTypeDTO[];
};

export default function ProgramDetailsPage() {
    const { id } = useParams();
    const [program, setProgram] = useState<ProgramDTO | null>(null);
    const [toast, setToast] = useState<string | null>(null);

    const load = async () => {
        const data = await getProgram(Number(id));
        setProgram(data);
    };

    useEffect(() => {
        load();
    }, [id]);

    if (!program) return null;

    const addDay = async () => {
        const name = prompt("Название дня");
        if (!name) return;

        try {
            await createDay(program.id, name);
            showToast("✅ День добавлен");
            await load();
        } catch {
            showToast("❌ Ошибка при добавлении дня");
        }
    };

    const removeDay = async (dayId: number) => {
        const confirmed = window.confirm("Вы уверены, что хотите удалить день?");
        if (!confirmed) return;

        try {
            await deleteDay(program.id, dayId);
            showToast("✅ День удалён");
            await load();
        } catch (e) {
            const msg = e instanceof Error ? e.message : "Ошибка при удалении дня";
            showToast(`❌ ${msg}`);
        }
    };

    const showToast = (text: string) => {
        setToast(text);
        setTimeout(() => setToast(null), 3000);
    };

    return (
        <div className="page stack">
            <h2 className="title">{program.name}</h2>

            <Button variant="primary" onClick={addDay}>
                + Добавить день
            </Button>

            {program.day_types && program.day_types.length > 0 ? (
                program.day_types.map((day) => (
                    <DayCard
                        key={day.id}
                        day={day}
                        programId={program.id}
                        onDelete={removeDay}
                    />
                ))
            ) : (
                <div>Дней пока нет</div>
            )}

            {toast && <div className="toast">{toast}</div>}
        </div>
    );
}
