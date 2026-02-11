import Button from "./Button";
import "../styles/ProgramCard.css";

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
            {/* Левый кликабельный блок */}
            <div
                className="program-left-block"
                onClick={onOpen}
                style={{ cursor: "pointer", flex: 1 }}
            >
                <div className="program-name">{name}</div>
                {active && <div className="badge">🟢 Активна</div>}
            </div>

            <div className="row-actions">
                <Button
                    onClick={onActivate}
                    variant={active ? "primary" : "ghost"}
                >
                    ⭐
                </Button>
                <Button onClick={onRename}>✏️</Button>
                <Button variant="danger" onClick={onDelete}>
                    🗑
                </Button>
            </div>
        </div>
    );
}
