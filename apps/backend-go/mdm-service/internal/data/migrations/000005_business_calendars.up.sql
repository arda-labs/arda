CREATE TABLE IF NOT EXISTS business_calendars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    timezone TEXT NOT NULL DEFAULT 'Asia/Ho_Chi_Minh',
    calendar_type TEXT NOT NULL DEFAULT 'BANKING',
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_business_calendars_code_active
    ON business_calendars (code)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS ix_business_calendars_status
    ON business_calendars (status)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS working_hours (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    calendar_id UUID NOT NULL REFERENCES business_calendars(id),
    day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    is_working_day BOOLEAN NOT NULL DEFAULT true,
    start_time TEXT NOT NULL DEFAULT '',
    end_time TEXT NOT NULL DEFAULT '',
    cutoff_time TEXT NOT NULL DEFAULT '',
    session_name TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS ix_working_hours_calendar
    ON working_hours (calendar_id, day_of_week, sort_order)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS calendar_exceptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    calendar_id UUID NOT NULL REFERENCES business_calendars(id),
    date DATE NOT NULL,
    exception_type TEXT NOT NULL DEFAULT 'HOLIDAY',
    name TEXT NOT NULL,
    is_working_day BOOLEAN NOT NULL DEFAULT false,
    start_time TEXT NOT NULL DEFAULT '',
    end_time TEXT NOT NULL DEFAULT '',
    cutoff_time TEXT NOT NULL DEFAULT '',
    source TEXT NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_calendar_exceptions_calendar_date_active
    ON calendar_exceptions (calendar_id, date)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS ix_calendar_exceptions_calendar_date
    ON calendar_exceptions (calendar_id, date)
    WHERE deleted_at IS NULL;

INSERT INTO business_calendars (code, name, timezone, calendar_type, description, status)
VALUES (
    'VN_BANKING_CALENDAR',
    'Lịch làm việc ngân hàng Việt Nam',
    'Asia/Ho_Chi_Minh',
    'BANKING',
    'Lịch chuẩn cho xử lý nghiệp vụ ngân hàng: làm việc thứ Hai đến thứ Sáu, nghỉ cuối tuần và ngày lễ.',
    'ACTIVE'
)
ON CONFLICT DO NOTHING;

INSERT INTO working_hours (calendar_id, day_of_week, is_working_day, start_time, end_time, cutoff_time, session_name, sort_order)
SELECT c.id, v.day_of_week, v.is_working_day, v.start_time, v.end_time, v.cutoff_time, v.session_name, v.sort_order
FROM business_calendars c
CROSS JOIN (
    VALUES
        (1, true, '08:00', '17:00', '15:30', 'Thứ hai', 10),
        (2, true, '08:00', '17:00', '15:30', 'Thứ ba', 20),
        (3, true, '08:00', '17:00', '15:30', 'Thứ tư', 30),
        (4, true, '08:00', '17:00', '15:30', 'Thứ năm', 40),
        (5, true, '08:00', '17:00', '15:30', 'Thứ sáu', 50),
        (6, false, '', '', '', 'Thứ bảy', 60),
        (7, false, '', '', '', 'Chủ nhật', 70)
) AS v(day_of_week, is_working_day, start_time, end_time, cutoff_time, session_name, sort_order)
WHERE c.code = 'VN_BANKING_CALENDAR'
  AND NOT EXISTS (
      SELECT 1 FROM working_hours wh
      WHERE wh.calendar_id = c.id
        AND wh.day_of_week = v.day_of_week
        AND wh.deleted_at IS NULL
  );

INSERT INTO calendar_exceptions (
    calendar_id, date, exception_type, name, is_working_day,
    start_time, end_time, cutoff_time, source, note
)
SELECT c.id, v.date::date, v.exception_type, v.name, false, '', '', '', v.source, v.note
FROM business_calendars c
CROSS JOIN (
    VALUES
        ('2026-01-01', 'HOLIDAY', 'Tết Dương lịch', 'VN_PUBLIC_HOLIDAY_2026', 'Nghỉ lễ bắt buộc theo lịch nghỉ lễ Việt Nam.'),
        ('2026-02-16', 'HOLIDAY', 'Tết Nguyên đán 2026', 'VN_PUBLIC_HOLIDAY_2026', 'Ngày nghỉ Tết Âm lịch.'),
        ('2026-02-17', 'HOLIDAY', 'Tết Nguyên đán 2026', 'VN_PUBLIC_HOLIDAY_2026', 'Ngày nghỉ Tết Âm lịch.'),
        ('2026-02-18', 'HOLIDAY', 'Tết Nguyên đán 2026', 'VN_PUBLIC_HOLIDAY_2026', 'Ngày nghỉ Tết Âm lịch.'),
        ('2026-02-19', 'HOLIDAY', 'Tết Nguyên đán 2026', 'VN_PUBLIC_HOLIDAY_2026', 'Ngày nghỉ Tết Âm lịch.'),
        ('2026-02-20', 'HOLIDAY', 'Tết Nguyên đán 2026', 'VN_PUBLIC_HOLIDAY_2026', 'Ngày nghỉ Tết Âm lịch.'),
        ('2026-04-26', 'HOLIDAY', 'Giỗ Tổ Hùng Vương', 'VN_PUBLIC_HOLIDAY_2026', 'Mùng 10 tháng 3 âm lịch; ngày này rơi vào Chủ nhật.'),
        ('2026-04-27', 'COMPENSATION_LEAVE', 'Nghỉ bù Giỗ Tổ Hùng Vương', 'VN_PUBLIC_HOLIDAY_2026', 'Nghỉ bù do Giỗ Tổ Hùng Vương rơi vào ngày nghỉ hằng tuần.'),
        ('2026-04-30', 'HOLIDAY', 'Ngày Giải phóng miền Nam, thống nhất đất nước', 'VN_PUBLIC_HOLIDAY_2026', 'Ngày lễ 30/4.'),
        ('2026-05-01', 'HOLIDAY', 'Ngày Quốc tế Lao động', 'VN_PUBLIC_HOLIDAY_2026', 'Ngày lễ 1/5.'),
        ('2026-09-01', 'HOLIDAY', 'Nghỉ lễ Quốc khánh', 'VN_PUBLIC_HOLIDAY_2026', 'Một ngày liền kề Quốc khánh theo lịch nghỉ lễ năm 2026.'),
        ('2026-09-02', 'HOLIDAY', 'Quốc khánh nước Cộng hòa xã hội chủ nghĩa Việt Nam', 'VN_PUBLIC_HOLIDAY_2026', 'Ngày lễ Quốc khánh 2/9.')
) AS v(date, exception_type, name, source, note)
WHERE c.code = 'VN_BANKING_CALENDAR'
ON CONFLICT DO NOTHING;
