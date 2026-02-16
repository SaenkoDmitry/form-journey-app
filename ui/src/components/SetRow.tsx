import EditableValue from "./EditableValue";
import Button from "./Button.tsx";
import "../styles/SetRow.css";
import "../styles/workout.css";
import {Check, X} from "lucide-react";

type Props = {
    set: FormattedSet;
    index: number;
    onDelete: () => void;
    onComplete: () => void;
    onChange: (id: number, reps: number, weight: number, minutes: number, meters: number) => void;
};

export default function SetRow({ set, index, onDelete, onComplete, onChange }: Props) {
    // const columnLength = (set.reps > 0 ? 1 : 0) + (set.weight > 0 ? 1 : 0) + (set.minutes > 0 ? 1 : 0) + (set.meters > 0 ? 1 : 0);

    return (
        <div className={`set-row ${set.completed ? "done" : ""}`}>
            <div className="set-index">{index + 1}</div>

            {set.reps > 0 && (
                <EditableValue
                    // columnLength={columnLength}
                    fact={set.fact_reps}
                    planned={set.reps}
                    suffix="повт."
                    typeParam="int"
                    completed={set.completed}
                    onSave={(v) => onChange(set.id, v, set.fact_weight, set.fact_minutes, set.fact_meters)}
                />
            )}

            {set.weight > 0 && (
                <EditableValue
                    // columnLength={columnLength}
                    fact={set.fact_weight}
                    planned={set.weight}
                    suffix="кг"
                    typeParam="float"
                    completed={set.completed}
                    onSave={(v) => onChange(set.id, set.fact_reps, v, set.fact_minutes, set.fact_meters)}
                />
            )}

            {set.minutes > 0 && (
                <EditableValue
                    // columnLength={columnLength}
                    fact={set.fact_minutes}
                    planned={set.minutes}
                    suffix="мин"
                    typeParam="int"
                    completed={set.completed}
                    onSave={(v) => onChange(set.id, set.fact_reps, set.fact_weight, v, set.fact_meters)}
                />
            )}

            {set.meters > 0 && (
                <EditableValue
                    // columnLength={columnLength}
                    fact={set.fact_meters}
                    planned={set.meters}
                    suffix="м"
                    typeParam="int"
                    completed={set.completed}
                    onSave={(v) => onChange(set.id, set.fact_reps, set.fact_weight, set.fact_minutes, v)}
                />
            )}

            <div className="set-actions">
                <Button variant={"active"} onClick={onComplete}><Check size={10}/></Button>
                <Button variant={"danger"} onClick={onDelete}><X size={10}/></Button>
            </div>
        </div>
    );
}
