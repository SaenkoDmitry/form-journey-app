import React, {useEffect, useState} from "react";
import {chooseProgram, createProgram, deleteProgram, getPrograms, renameProgram,} from "../api/programs";
import ProgramCard from "../components/ProgramCard";
import Button from "../components/Button";
import {useNavigate} from "react-router-dom";
import "../styles/ProgramBase.css";
import {Loader, Plus} from "lucide-react";

const ProgramsPage: React.FC = () => {
    const navigate = useNavigate();
    const [programs, setPrograms] = useState<any[]>([]);
    const [toast, setToast] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    const load = async () => {
        try {
            setLoading(true);
            const data = await getPrograms();
            setPrograms(data);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        load();
    }, []);

    const showToast = (text: string) => {
        setToast(text);
        setTimeout(() => setToast(null), 3000);
    };

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
    };

    const handleActivate = async (id: number, name: string) => {
        try {
            await chooseProgram(id);
            showToast(`✅ Программа '${name}' выбрана`);
            setPrograms((prev) =>
                prev.map((p) => ({...p, is_active: p.id === id}))
            );
        } catch {
            showToast("❌ Ошибка при выборе основной программы");
        }
    };

    const handleDelete = async (id: number) => {
        if (!window.confirm("Вы уверены, что хотите удалить программу?")) return;

        try {
            await deleteProgram(id);
            showToast("✅ Программа удалена");
            await load();
        } catch (e) {
            const msg = e instanceof Error ? e.message : "Ошибка при удалении программы";
            showToast(`❌ ${msg}`);
        }
    };

    return (
        <div className="page stack">
            <h1>Программы</h1>

            <Button variant="active" onClick={handleCreate}>
                <Plus size={14}/>Новая программа
            </Button>

            {loading && <Loader />}

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
};

export default ProgramsPage;
