package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *MdmRepo) ListBusinessCalendars(ctx context.Context, filter biz.PageFilter) ([]*biz.BusinessCalendar, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, timezone, calendar_type, description, status, created_at, updated_at
		FROM business_calendars
		WHERE deleted_at IS NULL
		  AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%' OR calendar_type ILIKE '%' || $2 || '%')
		ORDER BY code ASC, id ASC
		LIMIT $3 OFFSET $4`,
		filter.Status, filter.Keyword, page.Limit+1, page.Offset,
	)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	list, err := scanBusinessCalendars(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) GetBusinessCalendar(ctx context.Context, id string) (*biz.BusinessCalendar, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, timezone, calendar_type, description, status, created_at, updated_at
		FROM business_calendars
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanBusinessCalendar(row)
}

func (r *MdmRepo) GetBusinessCalendarByCode(ctx context.Context, code string) (*biz.BusinessCalendar, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, timezone, calendar_type, description, status, created_at, updated_at
		FROM business_calendars
		WHERE code = $1 AND deleted_at IS NULL`, code)
	return scanBusinessCalendar(row)
}

func (r *MdmRepo) CreateBusinessCalendar(ctx context.Context, item *biz.BusinessCalendar) (*biz.BusinessCalendar, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO business_calendars (code, name, timezone, calendar_type, description, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text`,
		item.Code, item.Name, item.Timezone, item.CalendarType, item.Description, item.Status,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.GetBusinessCalendar(ctx, item.ID)
}

func (r *MdmRepo) UpdateBusinessCalendar(ctx context.Context, item *biz.BusinessCalendar) (*biz.BusinessCalendar, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE business_calendars
		SET code = $2, name = $3, timezone = $4, calendar_type = $5,
		    description = $6, status = $7, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.Timezone, item.CalendarType, item.Description, item.Status,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetBusinessCalendar(ctx, item.ID)
}

func (r *MdmRepo) DeleteBusinessCalendar(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "business_calendars", id)
}

func (r *MdmRepo) ListWorkingHours(ctx context.Context, calendarID string) ([]*biz.WorkingHour, error) {
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, calendar_id::text, day_of_week, is_working_day, start_time, end_time,
		       cutoff_time, session_name, sort_order, created_at, updated_at
		FROM working_hours
		WHERE deleted_at IS NULL AND calendar_id::text = $1
		ORDER BY day_of_week ASC, sort_order ASC, id ASC`, calendarID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWorkingHours(rows)
}

func (r *MdmRepo) CreateWorkingHour(ctx context.Context, item *biz.WorkingHour) (*biz.WorkingHour, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO working_hours (
			calendar_id, day_of_week, is_working_day, start_time, end_time,
			cutoff_time, session_name, sort_order
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id::text`,
		item.CalendarID, item.DayOfWeek, item.IsWorkingDay, item.StartTime, item.EndTime,
		item.CutoffTime, item.SessionName, item.SortOrder,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.getWorkingHour(ctx, item.ID)
}

func (r *MdmRepo) UpdateWorkingHour(ctx context.Context, item *biz.WorkingHour) (*biz.WorkingHour, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE working_hours
		SET calendar_id = $2, day_of_week = $3, is_working_day = $4,
		    start_time = $5, end_time = $6, cutoff_time = $7, session_name = $8,
		    sort_order = $9, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		item.ID, item.CalendarID, item.DayOfWeek, item.IsWorkingDay, item.StartTime,
		item.EndTime, item.CutoffTime, item.SessionName, item.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getWorkingHour(ctx, item.ID)
}

func (r *MdmRepo) DeleteWorkingHour(ctx context.Context, id string) error {
	tag, err := r.data.db.Pool.Exec(ctx,
		`UPDATE working_hours SET deleted_at = now(), updated_at = now() WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return biz.ErrNotFound
	}
	return nil
}

func (r *MdmRepo) ListCalendarExceptions(ctx context.Context, filter biz.CalendarExceptionFilter) ([]*biz.CalendarException, error) {
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, calendar_id::text, to_char(date, 'YYYY-MM-DD'), exception_type, name,
		       is_working_day, start_time, end_time, cutoff_time, source, note, created_at, updated_at
		FROM calendar_exceptions
		WHERE deleted_at IS NULL
		  AND calendar_id::text = $1
		  AND ($2 = '' OR date >= $2::date)
		  AND ($3 = '' OR date <= $3::date)
		  AND ($4 = '' OR exception_type = $4)
		ORDER BY date ASC, id ASC`,
		filter.CalendarID, filter.FromDate, filter.ToDate, filter.ExceptionType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCalendarExceptions(rows)
}

func (r *MdmRepo) CreateCalendarException(ctx context.Context, item *biz.CalendarException) (*biz.CalendarException, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO calendar_exceptions (
			calendar_id, date, exception_type, name, is_working_day,
			start_time, end_time, cutoff_time, source, note
		)
		VALUES ($1, $2::date, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id::text`,
		item.CalendarID, item.Date, item.ExceptionType, item.Name, item.IsWorkingDay,
		item.StartTime, item.EndTime, item.CutoffTime, item.Source, item.Note,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.getCalendarException(ctx, item.ID)
}

func (r *MdmRepo) UpdateCalendarException(ctx context.Context, item *biz.CalendarException) (*biz.CalendarException, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE calendar_exceptions
		SET calendar_id = $2, date = $3::date, exception_type = $4, name = $5,
		    is_working_day = $6, start_time = $7, end_time = $8, cutoff_time = $9,
		    source = $10, note = $11, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		item.ID, item.CalendarID, item.Date, item.ExceptionType, item.Name, item.IsWorkingDay,
		item.StartTime, item.EndTime, item.CutoffTime, item.Source, item.Note,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getCalendarException(ctx, item.ID)
}

func (r *MdmRepo) DeleteCalendarException(ctx context.Context, id string) error {
	tag, err := r.data.db.Pool.Exec(ctx,
		`UPDATE calendar_exceptions SET deleted_at = now(), updated_at = now() WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return biz.ErrNotFound
	}
	return nil
}

func (r *MdmRepo) getWorkingHour(ctx context.Context, id string) (*biz.WorkingHour, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, calendar_id::text, day_of_week, is_working_day, start_time, end_time,
		       cutoff_time, session_name, sort_order, created_at, updated_at
		FROM working_hours
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanWorkingHour(row)
}

func (r *MdmRepo) getCalendarException(ctx context.Context, id string) (*biz.CalendarException, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, calendar_id::text, to_char(date, 'YYYY-MM-DD'), exception_type, name,
		       is_working_day, start_time, end_time, cutoff_time, source, note, created_at, updated_at
		FROM calendar_exceptions
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanCalendarException(row)
}

func scanBusinessCalendars(rows pgx.Rows) ([]*biz.BusinessCalendar, error) {
	var list []*biz.BusinessCalendar
	for rows.Next() {
		item, err := scanBusinessCalendar(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanBusinessCalendar(row pgx.Row) (*biz.BusinessCalendar, error) {
	item := &biz.BusinessCalendar{}
	err := row.Scan(
		&item.ID, &item.Code, &item.Name, &item.Timezone, &item.CalendarType,
		&item.Description, &item.Status, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func scanWorkingHours(rows pgx.Rows) ([]*biz.WorkingHour, error) {
	var list []*biz.WorkingHour
	for rows.Next() {
		item, err := scanWorkingHour(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanWorkingHour(row pgx.Row) (*biz.WorkingHour, error) {
	item := &biz.WorkingHour{}
	err := row.Scan(
		&item.ID, &item.CalendarID, &item.DayOfWeek, &item.IsWorkingDay,
		&item.StartTime, &item.EndTime, &item.CutoffTime, &item.SessionName,
		&item.SortOrder, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func scanCalendarExceptions(rows pgx.Rows) ([]*biz.CalendarException, error) {
	var list []*biz.CalendarException
	for rows.Next() {
		item, err := scanCalendarException(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanCalendarException(row pgx.Row) (*biz.CalendarException, error) {
	item := &biz.CalendarException{}
	err := row.Scan(
		&item.ID, &item.CalendarID, &item.Date, &item.ExceptionType, &item.Name,
		&item.IsWorkingDay, &item.StartTime, &item.EndTime, &item.CutoffTime,
		&item.Source, &item.Note, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
