import {downloadExcelWorkouts} from "../api/excel.ts";
import {Download, Loader} from "lucide-react";
import Button from "../components/Button.tsx";
import React, {useEffect, useState} from "react";
import Toast from "../components/Toast.tsx";
import {getExerciseGroups} from "../api/exercises.ts";
import {useNavigate} from "react-router-dom";

const StatsPage: React.FC = () => {
    const [toast, setToast] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [groups, setGroups] = useState<Group[]>([]);
    const navigate = useNavigate();

    useEffect(() => {
        const loadGroups = async () => {
            try {
                const data = await getExerciseGroups();
                setGroups(data);
            } catch (e) {
                setToast("Ошибка загрузки групп ❌");
            } finally {
                setLoading(false);
            }
        };

        loadGroups();
    }, []);

    return (
        <div>

            <div className={"page stack"}>

                <h1>Динамика</h1>

                {loading && <Loader/>}

                {!loading && <b>Выберите группу:</b>}
                {!loading && groups.map(g =>
                    <Button
                        variant="ghost"
                        onClick={() => navigate(`/statistics/${g.code}`)}
                    >{g.name}</Button>)
                }

                {!loading && <b>Или экспортируйте все данные разом:</b>}

                {!loading && <Button
                    variant="primary"
                    onClick={async () => {
                        try {
                            await downloadExcelWorkouts();
                            setToast("Файл Excel успешно скачан ✅");
                        } catch (err) {
                            console.error(err);
                            setToast("Ошибка при скачивании Excel ❌");
                        }
                    }}
                >
                    <Download size={16}/> Экспорт в Excel
                </Button>}

            </div>

            {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
        </div>
    );
}

export default StatsPage;
