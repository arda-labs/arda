package biz

import (
	"context"
	"strings"
	"time"
)

type BusinessCalendar struct {
	ID           string
	Code         string
	Name         string
	Timezone     string
	CalendarType string
	Description  string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type WorkingHour struct {
	ID           string
	CalendarID   string
	DayOfWeek    int
	IsWorkingDay bool
	StartTime    string
	EndTime      string
	CutoffTime   string
	SessionName  string
	SortOrder    int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CalendarException struct {
	ID            string
	CalendarID    string
	Date          string
	ExceptionType string
	Name          string
	IsWorkingDay  bool
	StartTime     string
	EndTime       string
	CutoffTime    string
	Source        string
	Note          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CalendarExceptionFilter struct {
	CalendarID    string
	FromDate      string
	ToDate        string
	ExceptionType string
}

type BusinessDayCalculation struct {
	CalendarID     string
	CalendarCode   string
	StartDate      string
	OffsetDays     int
	AdjustmentRule string
}

type BusinessDayCalculationResult struct {
	StartDate     string
	ResultDate    string
	CalendarDays  int
	SkippedDates  []string
	IsBusinessDay bool
}

func (uc *MdmUsecase) ListBusinessCalendars(ctx context.Context, filter PageFilter) ([]*BusinessCalendar, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListBusinessCalendars(ctx, filter)
}

func (uc *MdmUsecase) GetBusinessCalendar(ctx context.Context, id string) (*BusinessCalendar, error) {
	return uc.repo.GetBusinessCalendar(ctx, id)
}

func (uc *MdmUsecase) CreateBusinessCalendar(ctx context.Context, item *BusinessCalendar) (*BusinessCalendar, error) {
	normalizeBusinessCalendar(item)
	return uc.repo.CreateBusinessCalendar(ctx, item)
}

func (uc *MdmUsecase) UpdateBusinessCalendar(ctx context.Context, item *BusinessCalendar) (*BusinessCalendar, error) {
	normalizeBusinessCalendar(item)
	return uc.repo.UpdateBusinessCalendar(ctx, item)
}

func (uc *MdmUsecase) DeleteBusinessCalendar(ctx context.Context, id string) error {
	return uc.repo.DeleteBusinessCalendar(ctx, id)
}

func (uc *MdmUsecase) ListWorkingHours(ctx context.Context, calendarID string) ([]*WorkingHour, error) {
	return uc.repo.ListWorkingHours(ctx, strings.TrimSpace(calendarID))
}

func (uc *MdmUsecase) CreateWorkingHour(ctx context.Context, item *WorkingHour) (*WorkingHour, error) {
	normalizeWorkingHour(item)
	return uc.repo.CreateWorkingHour(ctx, item)
}

func (uc *MdmUsecase) UpdateWorkingHour(ctx context.Context, item *WorkingHour) (*WorkingHour, error) {
	normalizeWorkingHour(item)
	return uc.repo.UpdateWorkingHour(ctx, item)
}

func (uc *MdmUsecase) DeleteWorkingHour(ctx context.Context, id string) error {
	return uc.repo.DeleteWorkingHour(ctx, id)
}

func (uc *MdmUsecase) ListCalendarExceptions(ctx context.Context, filter CalendarExceptionFilter) ([]*CalendarException, error) {
	filter.CalendarID = strings.TrimSpace(filter.CalendarID)
	filter.FromDate = strings.TrimSpace(filter.FromDate)
	filter.ToDate = strings.TrimSpace(filter.ToDate)
	filter.ExceptionType = upperDefault(filter.ExceptionType, "")
	return uc.repo.ListCalendarExceptions(ctx, filter)
}

func (uc *MdmUsecase) CreateCalendarException(ctx context.Context, item *CalendarException) (*CalendarException, error) {
	normalizeCalendarException(item)
	return uc.repo.CreateCalendarException(ctx, item)
}

func (uc *MdmUsecase) UpdateCalendarException(ctx context.Context, item *CalendarException) (*CalendarException, error) {
	normalizeCalendarException(item)
	return uc.repo.UpdateCalendarException(ctx, item)
}

func (uc *MdmUsecase) DeleteCalendarException(ctx context.Context, id string) error {
	return uc.repo.DeleteCalendarException(ctx, id)
}

func (uc *MdmUsecase) CalculateBusinessDay(ctx context.Context, in BusinessDayCalculation) (*BusinessDayCalculationResult, error) {
	in.CalendarID = strings.TrimSpace(in.CalendarID)
	in.CalendarCode = upperDefault(in.CalendarCode, "")
	in.StartDate = strings.TrimSpace(in.StartDate)
	in.AdjustmentRule = upperDefault(in.AdjustmentRule, "FOLLOWING")
	start, err := time.Parse("2006-01-02", in.StartDate)
	if err != nil {
		return nil, ErrInvalidArgument
	}

	calendar, err := uc.resolveBusinessCalendar(ctx, in.CalendarID, in.CalendarCode)
	if err != nil {
		return nil, err
	}
	hours, err := uc.repo.ListWorkingHours(ctx, calendar.ID)
	if err != nil {
		return nil, err
	}
	exceptions, err := uc.repo.ListCalendarExceptions(ctx, CalendarExceptionFilter{CalendarID: calendar.ID})
	if err != nil {
		return nil, err
	}

	result, calendarDays, skipped := calculateBusinessDate(start, in.OffsetDays, in.AdjustmentRule, hours, exceptions)
	return &BusinessDayCalculationResult{
		StartDate:     start.Format("2006-01-02"),
		ResultDate:    result.Format("2006-01-02"),
		CalendarDays:  calendarDays,
		SkippedDates:  skipped,
		IsBusinessDay: isBusinessDay(result, hours, exceptions),
	}, nil
}

func normalizeBusinessCalendar(item *BusinessCalendar) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.Timezone = defaultString(item.Timezone, "Asia/Ho_Chi_Minh")
	item.CalendarType = upperDefault(item.CalendarType, "BANKING")
	item.Description = strings.TrimSpace(item.Description)
	item.Status = upperDefault(item.Status, "ACTIVE")
}

func normalizeWorkingHour(item *WorkingHour) {
	item.CalendarID = strings.TrimSpace(item.CalendarID)
	item.StartTime = strings.TrimSpace(item.StartTime)
	item.EndTime = strings.TrimSpace(item.EndTime)
	item.CutoffTime = strings.TrimSpace(item.CutoffTime)
	item.SessionName = strings.TrimSpace(item.SessionName)
	if item.DayOfWeek < 1 || item.DayOfWeek > 7 {
		item.DayOfWeek = 1
	}
	if item.SortOrder == 0 {
		item.SortOrder = item.DayOfWeek * 10
	}
}

func normalizeCalendarException(item *CalendarException) {
	item.CalendarID = strings.TrimSpace(item.CalendarID)
	item.Date = strings.TrimSpace(item.Date)
	item.ExceptionType = upperDefault(item.ExceptionType, "HOLIDAY")
	item.Name = strings.TrimSpace(item.Name)
	item.StartTime = strings.TrimSpace(item.StartTime)
	item.EndTime = strings.TrimSpace(item.EndTime)
	item.CutoffTime = strings.TrimSpace(item.CutoffTime)
	item.Source = strings.TrimSpace(item.Source)
	item.Note = strings.TrimSpace(item.Note)
}

func (uc *MdmUsecase) resolveBusinessCalendar(ctx context.Context, id, code string) (*BusinessCalendar, error) {
	if id != "" {
		return uc.repo.GetBusinessCalendar(ctx, id)
	}
	if code != "" {
		return uc.repo.GetBusinessCalendarByCode(ctx, code)
	}
	return nil, ErrInvalidArgument
}

func calculateBusinessDate(start time.Time, offset int, adjustmentRule string, hours []*WorkingHour, exceptions []*CalendarException) (time.Time, int, []string) {
	if offset == 0 {
		if isBusinessDay(start, hours, exceptions) || adjustmentRule == "NONE" {
			return start, 0, nil
		}
		return adjustBusinessDate(start, adjustmentRule, hours, exceptions)
	}

	step := 1
	if offset < 0 {
		step = -1
		offset = -offset
	}
	current := start
	calendarDays := 0
	skipped := make([]string, 0)
	for counted := 0; counted < offset; {
		current = current.AddDate(0, 0, step)
		calendarDays += step
		if isBusinessDay(current, hours, exceptions) {
			counted++
			continue
		}
		skipped = append(skipped, current.Format("2006-01-02"))
	}
	return current, calendarDays, skipped
}

func adjustBusinessDate(start time.Time, adjustmentRule string, hours []*WorkingHour, exceptions []*CalendarException) (time.Time, int, []string) {
	direction := 1
	if adjustmentRule == "PRECEDING" {
		direction = -1
	}
	current := start
	calendarDays := 0
	skipped := make([]string, 0)
	for !isBusinessDay(current, hours, exceptions) {
		skipped = append(skipped, current.Format("2006-01-02"))
		current = current.AddDate(0, 0, direction)
		calendarDays += direction
	}
	return current, calendarDays, skipped
}

func isBusinessDay(date time.Time, hours []*WorkingHour, exceptions []*CalendarException) bool {
	day := date.Format("2006-01-02")
	for _, item := range exceptions {
		if item.Date == day {
			return item.IsWorkingDay
		}
	}
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	for _, item := range hours {
		if item.DayOfWeek == weekday {
			return item.IsWorkingDay
		}
	}
	return false
}
