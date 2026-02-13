interface Workout {
    id: number;
    name: string;
    started_at?: string;
    duration?: string;
    status?: string;
    completed: boolean;
}

interface ShowMyWorkoutsResult {
    items: Workout[];
    pagination: Pagination;
}

interface Pagination {
    limit: number;
    offset: number;
    total: number;
}

interface User {
    id: number;
    first_name: string;
    username?: string;
    photo_url?: string;
}

interface ExerciseType {
    id: number;
    name: string;
    url: string;
    exercise_group_type_code: string;
    rest_in_seconds: number;
    accent: string;
    units: string;
    description: string;
}

interface Set {
    ID: number;
    Reps: number;
    FactReps: number;
    Weight: number;
    FactWeight: number;
    Minutes: number;
    FactMinutes: number;
    Meters: number;
    FactMeters: number;
    Completed: boolean;
    CompletedAt: string;
    Index: number;
}

interface Exercise {
    ID: number;
    ExerciseType: ExerciseType;
    Sets: Set[];
    Index: number;
}

interface ReadWorkoutDTO {
    progress: WorkoutProgress;
    Stats: WorkoutStatistic;
}

interface WorkoutProgress {
    workout: FormattedWorkout;
    TotalExercises: number;
    CompletedExercises: number;
    TotalSets: number;
    CompletedSets: number;
    ProgressPercent: number;
    RemainingMin: number | null;
    SessionStarted: boolean;
}

interface WorkoutStatistic {
    DayType: any;
    WorkoutDay: any;
    TotalWeight: number;
    CompletedExercises: number;
    CardioTime: number; // в минутах
    ExerciseTypesMap: Record<number, any>;
    ExerciseWeightMap: Record<number, number>;
    ExerciseTimeMap: Record<number, number>;
}

interface FormattedWorkout {
    id: number
    user_id: number
    status: string
    started_at: string
    duration: string
    ended_at: string
    day_type_name: string
    completed: boolean
    exercises: FormattedExercise[]
}

interface FormattedExercise {
    id: number
    name: string
    url: string
    group_name: string
    rest_in_seconds: number
    accent: string
    units: string
    description: string
    index: number
    sets: FormattedSet[]
    sum_weight: number
}

interface FormattedSet {
    id: number
    formatted_string: string
    completed: boolean
    completed_at: string
    index: number
}

interface Measurement {
    id: number;
    user_id: number;
    created_at: string; // ISO string
    shoulders: number;
    chest: number;
    hand_left: number;
    hand_right: number;
    waist: number;
    buttocks: number;
    hip_left: number;
    hip_right: number;
    calf_left: number;
    calf_right: number;
    weight: number;
}

interface FindWithOffsetLimitMeasurement {
    items: Measurement[];
    count: number;
}


interface ToCreateMeasurement {
    user_id: number;
    shoulders: number;
    chest: number;
    hand_left: number;
    hand_right: number;
    waist: number;
    buttocks: number;
    hip_left: number;
    hip_right: number;
    calf_left: number;
    calf_right: number;
    weight: number;
}

type Group = {
    code: string;
    name: string;
};

interface WorkoutDayTypeDTO {
    id: number;
    program_id: number;   // оставил snake_case, чтобы совпадало с API
    name: string;
    preset: string;
    created_at: string;   // ISO string
}

interface StartWorkoutDTO {
    workout_id: number;
}

interface CurrentExerciseSession {
    exercise: FormattedExercise;
    exercise_type: ExerciseType;
    day_type: WorkoutDayTypeDTO;
    workout: FormattedWorkout;
    exercise_index: number;
}

interface ProgramDTO {
    id: number;
    user_id: number;
    name: string;
    created_at: string;
    day_types: WorkoutDayTypeDTO[];
}

type SetDTO = {
    reps: number;
    weight: number;
    minutes: number;
    meters: number;
};

type ExerciseDTO = {
    id: number;
    name: string;
    sets: SetDTO[];
};

type PresetListDTO = {
    exercises: ExerciseDTO[];
};