import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import {
  BusinessCalendar,
  BusinessDayCalculation,
  BusinessDayCalculationResult,
  CalendarException,
  ListOptions,
  PageResponse,
  WorkingHour,
} from '../models/mdm.models';
import { buildParams } from './mdm-http';

@Injectable({ providedIn: 'root' })
export class BusinessCalendarService {
  private http = inject(HttpClient);

  list(options: ListOptions = {}): Observable<PageResponse<BusinessCalendar>> {
    return this.http.get<any>('/api/v1/mdm/business-calendars', { params: buildParams(options) }).pipe(
      map(resp => ({
        items: (resp.business_calendars ?? resp.businessCalendars ?? []).map((item: any) => this.toCalendar(item)),
        nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '',
      })),
    );
  }

  create(item: Partial<BusinessCalendar>): Observable<BusinessCalendar> {
    return this.http.post<any>('/api/v1/mdm/business-calendars', { business_calendar: this.fromCalendar(item) }).pipe(
      map(resp => this.toCalendar(resp)),
    );
  }

  update(id: string, item: Partial<BusinessCalendar>): Observable<BusinessCalendar> {
    return this.http.put<any>(`/api/v1/mdm/business-calendars/${encodeURIComponent(id)}`, { business_calendar: this.fromCalendar(item) }).pipe(
      map(resp => this.toCalendar(resp)),
    );
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/business-calendars/${encodeURIComponent(id)}`);
  }

  listWorkingHours(calendarId: string): Observable<WorkingHour[]> {
    return this.http.get<any>(`/api/v1/mdm/business-calendars/${encodeURIComponent(calendarId)}/working-hours`).pipe(
      map(resp => (resp.working_hours ?? resp.workingHours ?? []).map((item: any) => this.toWorkingHour(item))),
    );
  }

  createWorkingHour(calendarId: string, item: Partial<WorkingHour>): Observable<WorkingHour> {
    return this.http.post<any>(`/api/v1/mdm/business-calendars/${encodeURIComponent(calendarId)}/working-hours`, { working_hour: this.fromWorkingHour(item) }).pipe(
      map(resp => this.toWorkingHour(resp)),
    );
  }

  updateWorkingHour(id: string, item: Partial<WorkingHour>): Observable<WorkingHour> {
    return this.http.put<any>(`/api/v1/mdm/working-hours/${encodeURIComponent(id)}`, { working_hour: this.fromWorkingHour(item) }).pipe(
      map(resp => this.toWorkingHour(resp)),
    );
  }

  deleteWorkingHour(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/working-hours/${encodeURIComponent(id)}`);
  }

  listExceptions(calendarId: string, options: { fromDate?: string; toDate?: string; exceptionType?: string } = {}): Observable<CalendarException[]> {
    let params = new HttpParams();
    if (options.fromDate) params = params.set('from_date', options.fromDate);
    if (options.toDate) params = params.set('to_date', options.toDate);
    if (options.exceptionType) params = params.set('exception_type', options.exceptionType);
    return this.http.get<any>(`/api/v1/mdm/business-calendars/${encodeURIComponent(calendarId)}/exceptions`, { params }).pipe(
      map(resp => (resp.calendar_exceptions ?? resp.calendarExceptions ?? []).map((item: any) => this.toException(item))),
    );
  }

  createException(calendarId: string, item: Partial<CalendarException>): Observable<CalendarException> {
    return this.http.post<any>(`/api/v1/mdm/business-calendars/${encodeURIComponent(calendarId)}/exceptions`, { calendar_exception: this.fromException(item) }).pipe(
      map(resp => this.toException(resp)),
    );
  }

  updateException(id: string, item: Partial<CalendarException>): Observable<CalendarException> {
    return this.http.put<any>(`/api/v1/mdm/calendar-exceptions/${encodeURIComponent(id)}`, { calendar_exception: this.fromException(item) }).pipe(
      map(resp => this.toException(resp)),
    );
  }

  deleteException(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/calendar-exceptions/${encodeURIComponent(id)}`);
  }

  calculate(request: BusinessDayCalculation): Observable<BusinessDayCalculationResult> {
    return this.http.post<any>('/api/v1/mdm/business-calendars/calculate', {
      calendar_id: request.calendarId,
      calendar_code: request.calendarCode,
      start_date: request.startDate,
      offset_days: request.offsetDays,
      adjustment_rule: request.adjustmentRule,
    }).pipe(
      map(resp => ({
        startDate: resp.start_date ?? resp.startDate ?? '',
        resultDate: resp.result_date ?? resp.resultDate ?? '',
        calendarDays: resp.calendar_days ?? resp.calendarDays ?? 0,
        skippedDates: resp.skipped_dates ?? resp.skippedDates ?? [],
        isBusinessDay: resp.is_business_day ?? resp.isBusinessDay ?? false,
      })),
    );
  }

  private toCalendar(item: any): BusinessCalendar {
    return {
      id: item.id ?? '',
      code: item.code ?? '',
      name: item.name ?? '',
      timezone: item.timezone ?? 'Asia/Ho_Chi_Minh',
      calendarType: item.calendar_type ?? item.calendarType ?? 'BANKING',
      description: item.description ?? '',
      status: item.status ?? 'ACTIVE',
    };
  }

  private fromCalendar(item: Partial<BusinessCalendar>): any {
    return {
      id: item.id,
      code: item.code,
      name: item.name,
      timezone: item.timezone,
      calendar_type: item.calendarType,
      description: item.description,
      status: item.status,
    };
  }

  private toWorkingHour(item: any): WorkingHour {
    return {
      id: item.id ?? '',
      calendarId: item.calendar_id ?? item.calendarId ?? '',
      dayOfWeek: item.day_of_week ?? item.dayOfWeek ?? 1,
      isWorkingDay: item.is_working_day ?? item.isWorkingDay ?? false,
      startTime: item.start_time ?? item.startTime ?? '',
      endTime: item.end_time ?? item.endTime ?? '',
      cutoffTime: item.cutoff_time ?? item.cutoffTime ?? '',
      sessionName: item.session_name ?? item.sessionName ?? '',
      sortOrder: item.sort_order ?? item.sortOrder ?? 0,
    };
  }

  private fromWorkingHour(item: Partial<WorkingHour>): any {
    return {
      id: item.id,
      calendar_id: item.calendarId,
      day_of_week: item.dayOfWeek,
      is_working_day: item.isWorkingDay,
      start_time: item.startTime,
      end_time: item.endTime,
      cutoff_time: item.cutoffTime,
      session_name: item.sessionName,
      sort_order: item.sortOrder,
    };
  }

  private toException(item: any): CalendarException {
    return {
      id: item.id ?? '',
      calendarId: item.calendar_id ?? item.calendarId ?? '',
      date: item.date ?? '',
      exceptionType: item.exception_type ?? item.exceptionType ?? 'HOLIDAY',
      name: item.name ?? '',
      isWorkingDay: item.is_working_day ?? item.isWorkingDay ?? false,
      startTime: item.start_time ?? item.startTime ?? '',
      endTime: item.end_time ?? item.endTime ?? '',
      cutoffTime: item.cutoff_time ?? item.cutoffTime ?? '',
      source: item.source ?? '',
      note: item.note ?? '',
    };
  }

  private fromException(item: Partial<CalendarException>): any {
    return {
      id: item.id,
      calendar_id: item.calendarId,
      date: item.date,
      exception_type: item.exceptionType,
      name: item.name,
      is_working_day: item.isWorkingDay,
      start_time: item.startTime,
      end_time: item.endTime,
      cutoff_time: item.cutoffTime,
      source: item.source,
      note: item.note,
    };
  }
}
