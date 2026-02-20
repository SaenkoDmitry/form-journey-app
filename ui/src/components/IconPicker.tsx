import React from "react";
import {
    BicepsFlexed,
    Bike,
    Bird,
    Brain,
    Bus,
    Cat,
    Cherry,
    ChessKing,
    ChessQueen,
    Coffee,
    ContactRound,
    Dog,
    Dumbbell,
    Flame,
    Headphones,
    Heart,
    Lollipop,
    Panda,
    Rabbit,
    Smile,
    Squirrel,
    Star,
    Sun,
    UserRound
} from "lucide-react";

export const ICONS = {
    Smile,
    UserRound,
    ContactRound,
    Star,
    BicepsFlexed,
    Dumbbell,
    Bike,
    Brain,
    Headphones,
    Coffee,
    Sun,
    Cherry,
    Lollipop,
    Heart,
    Flame,
    Panda,
    Cat,
    Dog,
    Rabbit,
    Bird,
    Squirrel,
    ChessKing,
    ChessQueen,
    Bus,
};

export type IconName = keyof typeof ICONS;

interface Props {
    selected: IconName;
    onSelect: (name: IconName) => void;
    onClose: () => void;
}

const IconPicker: React.FC<Props> = ({selected, onSelect, onClose}) => {
    return (
        <div
            style={{
                position: "fixed",
                inset: 0,
                background: "rgba(0,0,0,0.4)",
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                zIndex: 1000,
            }}
            onClick={onClose}
        >
            <div
                onClick={(e) => e.stopPropagation()}
                style={{
                    background: "var(--color-bg)",
                    padding: "1.5rem",
                    borderRadius: 20,
                    width: 320,
                    display: "grid",
                    gridTemplateColumns: "repeat(4, 1fr)",
                    gap: 16,
                    boxShadow: "0 10px 30px rgba(0,0,0,0.15)",
                }}
            >
                {Object.entries(ICONS).map(([name, Icon]) => (
                    <div
                        key={name}
                        onClick={() => onSelect(name as IconName)}
                        style={{
                            cursor: "pointer",
                            padding: 12,
                            borderRadius: 14,
                            display: "flex",
                            justifyContent: "center",
                            transition: "all 0.2s ease",
                            background:
                                selected === name
                                    ? "rgba(0,0,0,0.06)"
                                    : "transparent",
                            border:
                                selected === name
                                    ? "2px solid #000"
                                    : "2px solid transparent",
                        }}
                    >
                        <Icon size={28}/>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default IconPicker;
