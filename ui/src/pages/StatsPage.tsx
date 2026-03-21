import {downloadExcelWorkouts} from "../api/excel.ts";
import {Download, Loader} from "lucide-react";
import Button from "../components/Button.tsx";
import React, {useState} from "react";
import Toast from "../components/Toast.tsx";
import {useNavigate} from "react-router-dom";

const StatsPage: React.FC = () => {
    const [toast, setToast] = useState<string | null>(null);
    const navigate = useNavigate();

    return (
        <div>

            <div className={"page stack"}>

                <h1>Динамика</h1>

                <Button variant={"ghost"} onClick={() => navigate(`/statistics/exercise-groups`)}>Упражнения</Button>
                <Button variant={"ghost"} onClick={() => navigate(`/statistics/measurements`)}>Замеры</Button>

                {<b>Или экспортируйте все данные разом:</b>}

                {<Button
                    variant="active"
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
