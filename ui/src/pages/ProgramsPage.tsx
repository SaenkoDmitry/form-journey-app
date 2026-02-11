import { useEffect, useState } from "react";
import {
    getPrograms,
    createProgram,
    deleteProgram,
    chooseProgram,
    renameProgram,
} from "../api/programs";
import ProgramCard from "../components/ProgramCard";
import Button from "../components/Button";
import { useNavigate } from "react-router-dom";
import "../styles/ProgramBase.css";

export default function ProgramsPage() {
    const [programs, setPrograms] = useState<any[]>([]);
    const navigate = useNavigate();
    const [toast, setToast] = useState<string | null>(null);

    const load = async () => {
        setPrograms(await getPrograms());
    };

    useEffect(() => {
        load();
    }, []);

    const handleCreate = async () => {
        const name = prompt("Название программы");
        if (!name) return;

        try {
            await createProgram(name);
            showToast("✅ Программа создана");
            await load();
        } catch {
            showToast("❌ Ошибка при создании программы");
        }
        load();
    };

    const handleRename = async (id: number, oldName: string) => {
        const name = prompt("Новое название", oldName);
        if (!name) return;

        try {
            await renameProgram(id, name);
            showToast("✅ Программа переименована");
            await load();
        } catch {
            showToast("❌ Ошибка при переименовании программы");
        }
        load();
    };

    const handleActivate = async (id: number, name: string) => {
        try {
            await chooseProgram(id);
            showToast("✅ Программа '" + name + "' выбрана");
            await load();
        } catch {
            showToast("❌ Ошибка при выборе основной программы");
        }

        setPrograms((prev) =>
            prev.map((p) => ({
                ...p, is_active: p.id === id,
            }))
        );
    };

    const handleDelete = async (id: number) => {

        try {
            await deleteProgram(id);
            showToast("✅ Программа удалена");
            await load();
        } catch (e) {
            const msg = e instanceof Error ? e.message : "Ошибка при удалении программы";
            showToast(`❌ ${msg}`);
        }
        load();
    }

    const showToast = (text: string) => {
        setToast(text);
        setTimeout(() => setToast(null), 3000);
    };

    return (
        <div className="page stack">
            <h1>Программы</h1>

            <Button variant="primary" onClick={handleCreate}>
                + Новая программа
            </Button>

            {programs.map((p) => (
                <ProgramCard
                    key={p.id}
                    name={p.name}
                    active={p.is_active}
                    onOpen={() => navigate(`/programs/${p.id}`)}
                    onActivate={() => handleActivate(p.id, p.name)}
                    onRename={() => handleRename(p.id, p.name)}
                    onDelete={() => handleDelete(p.id)}
                />
            ))}

            {toast && <div className="toast">{toast}</div>}
        </div>
    );
}
