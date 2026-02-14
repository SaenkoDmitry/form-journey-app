import Button from "./Button";
import {Edit, Pencil, PencilLine, Star, Trash2} from "lucide-react";

type Props = {
    name: string;
    active?: boolean;
    onOpen: () => void;
    onActivate: () => void;
    onRename: () => void;
    onDelete: () => void;
};

export default function ProgramCard({
                                        name,
                                        active,
                                        onOpen,
                                        onActivate,
                                        onRename,
                                        onDelete,
                                    }: Props) {
    return (
        <div className="card row">
            <div
                onClick={onOpen}
                style={{ cursor: "pointer", flex: 1 }}
            >
                <b>{name}</b>
                {active && <div className="badge">🟢 Активна</div>}
            </div>

            <div className="row-actions">
                <Button
                    onClick={onActivate}
                    variant={active ? "active" : "ghost"}
                >
                    <Star size={14}/>
                </Button>
                <Button onClick={onRename}><PencilLine size={14}/></Button>
                <Button variant="danger" onClick={onDelete}>
                    <Trash2 size={14}/>
                </Button>
            </div>
        </div>
    );
}
