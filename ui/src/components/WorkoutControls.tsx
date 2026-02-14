import Button from "./Button.tsx";
import {useNavigate} from "react-router-dom";
import {ArrowLeft, ArrowRight, BarChart3, ChevronLeft, ChevronRight, LineChart} from "lucide-react";

interface WorkoutControlsProps {
    onPrev: () => void;
    onNext: () => void;
    workoutID: number;
    disablePrev?: boolean;
    disableNext?: boolean;
}


export default function WorkoutControls({onPrev, onNext, workoutID, disablePrev, disableNext}: WorkoutControlsProps) {
    const navigate = useNavigate();

    return (
        <div className="controls">
            <Button
                variant="ghost"
                onClick={onPrev}
                disabled={disablePrev}
            >
                {/*<ArrowLeft/>*/}
                <ChevronLeft/>
            </Button>
            <Button
                variant="ghost"
                onClick={() => navigate(`/workouts/${workoutID}`)}
            >
                <BarChart3 size={18} className="flex items-center gap-2 [&>svg]:translate-y-[1px]" />
                <div className="leading-none">Прогресс</div>
            </Button>
            <Button
                variant="ghost"
                onClick={onNext}
                disabled={disableNext}
            >
                {/*<ArrowRight/>*/}
                <ChevronRight/>
            </Button>
        </div>
    );
}
