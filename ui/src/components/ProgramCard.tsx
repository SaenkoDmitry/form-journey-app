import Button from "./Button";

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
