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
