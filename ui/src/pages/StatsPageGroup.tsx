import {Loader} from "lucide-react";
import Button from "../components/Button.tsx";
import React, {useEffect, useState} from "react";
import Toast from "../components/Toast.tsx";
import {getExerciseGroups, getExerciseTypesByGroup} from "../api/exercises.ts";
import {useNavigate, useParams} from "react-router-dom";

const StatsPageGroup: React.FC = () => {
    const [toast, setToast] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();
    const [groupsMap, setGroupsMap] = useState<Record<string, Group>>({});
    const {groupCode} = useParams();
    const [exerciseTypes, setExerciseTypes] = useState<ExerciseType[]>([]);

    useEffect(() => {
        if (!groupCode) return;

        const fetchGroupData = async () => {
            try {
                setLoading(true);

                const [exerciseTypes, groups]: [ExerciseType[], Group[]] = await Promise.all([
                    getExerciseTypesByGroup(groupCode),
                    getExerciseGroups()
                ]);

                const groupsMap = groups.reduce<Record<string, Group>>((acc, group) => {
                    acc[group.code] = group;
                    return acc;
                }, {});

                setExerciseTypes(exerciseTypes);
                setGroupsMap(groupsMap);

            } catch (err: any) {
                setError(err.message || 'Не удалось загрузить данные');
            } finally {
                setLoading(false);
            }
        }

        fetchGroupData();
    }, [groupCode]);

    if (loading) return <Loader/>;
    if (error) return <p style={{color: 'red'}}>{error}</p>;

    return (
        <div>

            <div className={"page stack"}>

                {groupCode && <h1>Динамика: {groupsMap[groupCode].name}</h1>}

                {exerciseTypes.map(ex => <Button
                        variant="ghost"
                        onClick={() => navigate(`/statistics/${groupCode}/exercise/${ex.id}`)}>
                        {ex.name}
                    </Button>
                )}
            </div>

            {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
        </div>
    );
}

export default StatsPageGroup;
