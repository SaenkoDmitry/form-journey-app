-- +goose Up
-- +goose StatementBegin
INSERT INTO exercise_types (name, url, exercise_group_type_code, rest_in_seconds, accent, units)
VALUES ('Сгибания рук с гантелями обратным хватом', 'https://disk.yandex.ru/i/pcNCAYlDR_mclA', 'biceps', 120,
        'плечелучевая мышца предплечья.<br>Эта мышца отвечает за сгибание в локтевом суставе и играет ключевую роль в развитии силы хвата и общего объёма предплечья.<br>В обычных подъёмах на бицепс плечелучевая мышца почти не работает, а обратный хват загружает её по максимуму.<br>Также нагрузка акцентируется на бицепсе плеча — при этом акцент делается на плечевой мышце (брахиалисе), расположенной под бицепсом, и на длинной (внешней) головке самого бицепса.', 'reps,weight');
INSERT INTO exercise_types (name, url, exercise_group_type_code, rest_in_seconds, accent, units)
VALUES ('Сгибания рук с гантелями', 'https://disk.yandex.ru/i/CdzoIQ-K1iOdjg', 'biceps', 120,
        'двуглавая мышца плеча (бицепс)', 'reps,weight');
INSERT INTO exercise_types (name, url, exercise_group_type_code, rest_in_seconds, accent, units)
VALUES ('Поочередные сгибания рук с гантелями', 'https://disk.yandex.ru/i/JIv3ARvCGCBlkg', 'biceps', 120,
        'двуглавая мышца плеча (бицепс)', 'reps,weight');
INSERT INTO exercise_types (name, url, exercise_group_type_code, rest_in_seconds, accent, units)
VALUES ('Сгибания рук со штангой', 'https://disk.yandex.ru/i/NKcBXhlqKo1VBw', 'biceps', 120,
        '<b>Стандартно на ширине плеч</b> — акцент на весь бицепс.<br>
<b>Узким хватом</b> — акцент на внешнюю головку (длинную) бицепса.<br>
<b>Широким хватом</b> — акцент на внутреннюю головку (короткую) бицепса.<br>
Также можно использовать изогнутый EZ-гриф, который акцентирует нагрузку на брахиалис и обе головки бицепса.', 'reps,weight');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM exercise_types where name = 'Сгибания рук с гантелями обратным хватом';
DELETE FROM exercise_types where name = 'Сгибания рук с гантелями';
DELETE FROM exercise_types where name = 'Поочередные сгибания рук с гантелями';
DELETE FROM exercise_types where name = 'Сгибания рук со штангой';
-- +goose StatementEnd
