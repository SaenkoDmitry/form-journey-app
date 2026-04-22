interface WorkoutCardProps {
    w: Workout;
    idx: number;
}

export default function WorkoutCard({ w, idx }: WorkoutCardProps) {
    return (
        <div
            style={{
                padding: 16,
                borderRadius: 16,
                boxShadow: '0 4px 12px rgba(0,0,0,0.08)',
                transition: '0.2s',
            }}
        >

            <h2 style={{ margin: 0 }}>{idx}.{w.name}</h2>

            <div style={{ padding: '4px 0', opacity: 0.6 }}>
                {w.started_at}
            </div>

            <div
                style={{
                    margin: 4,
                    fontWeight: 600,
                    color:
                        w.status === 'finished'
                            ? 'var(--color-primary-hover)'
                            : w.status === 'in_progress'
                                ? '#f9a825'
                                : '#999',
                }}
            >
                {w.status}
            </div>
        </div>
    );
}
