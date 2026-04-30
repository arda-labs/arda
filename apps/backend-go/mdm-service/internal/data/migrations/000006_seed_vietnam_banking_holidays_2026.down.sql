DELETE FROM calendar_exceptions
WHERE source = 'VN_PUBLIC_HOLIDAY_2026'
  AND calendar_id IN (
      SELECT id FROM business_calendars WHERE code = 'VN_BANKING_CALENDAR'
  );
