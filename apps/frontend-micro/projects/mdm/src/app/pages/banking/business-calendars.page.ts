import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, effect, inject, signal } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { forkJoin } from 'rxjs';
import { ConfirmationService, MessageService } from 'primeng/api';
import { Button } from 'primeng/button';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { DatePicker } from 'primeng/datepicker';
import { Dialog } from 'primeng/dialog';
import { InputText } from 'primeng/inputtext';
import { Select } from 'primeng/select';
import { TableModule } from 'primeng/table';
import { Tag } from 'primeng/tag';
import { Textarea } from 'primeng/textarea';
import { Toast } from 'primeng/toast';
import { Tooltip } from 'primeng/tooltip';
import { BusinessCalendar, CalendarException, WorkingHour } from '../../models/mdm.models';
import { BusinessCalendarService } from '../../services/business-calendar.service';
import { statusSeverity } from '../../services/mdm-http';

interface CalendarDay {
  date: string;
  label: number;
  inMonth: boolean;
  isToday: boolean;
  isWorkingDay: boolean;
  exception?: CalendarException;
}

@Component({
  selector: 'app-business-calendars-page',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    Button,
    ConfirmDialog,
    DatePicker,
    Dialog,
    InputText,
    Select,
    TableModule,
    Tag,
    Textarea,
    Toast,
    Tooltip,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './business-calendars.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class BusinessCalendarsPage {
  private service = inject(BusinessCalendarService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly selectedCalendarId = signal('');
  readonly selectedDay = signal<CalendarDay | null>(null);
  readonly workingHours = signal<WorkingHour[]>([]);
  readonly exceptions = signal<CalendarException[]>([]);
  readonly currentMonth = signal(this.monthStart(new Date()));
  readonly calendarDialogVisible = signal(false);
  readonly hourDialogVisible = signal(false);
  readonly exceptionDialogVisible = signal(false);
  readonly selectedCalendar = signal<BusinessCalendar | null>(null);
  readonly selectedHour = signal<WorkingHour | null>(null);
  readonly selectedException = signal<CalendarException | null>(null);
  readonly isSaving = signal(false);
  readonly calculationResult = signal('');

  readonly resource = rxResource({
    stream: () => this.service.list({ pageSize: 200 }),
  });
  readonly calendars = computed(() => this.resource.value()?.items ?? []);
  readonly activeCalendars = computed(() => this.calendars().filter(item => item.status === 'ACTIVE').length);
  readonly currentCalendar = computed(() => this.calendars().find(item => item.id === this.selectedCalendarId()) ?? null);
  readonly monthTitle = computed(() => {
    const month = this.currentMonth();
    return `Tháng ${month.getMonth() + 1}/${month.getFullYear()}`;
  });
  readonly monthDays = computed(() => this.buildMonthDays());
  readonly monthWorkingDays = computed(() => this.monthDays().filter(day => day.inMonth && day.isWorkingDay).length);
  readonly monthHolidayCount = computed(() => this.monthDays().filter(day => day.inMonth && !day.isWorkingDay).length);

  readonly statusOptions = [
    { label: 'Hiệu lực', value: 'ACTIVE' },
    { label: 'Ngưng hiệu lực', value: 'INACTIVE' },
  ];
  readonly calendarTypeOptions = [
    { label: 'Ngân hàng', value: 'BANKING' },
    { label: 'Giao dịch', value: 'TRADING' },
    { label: 'Vận hành', value: 'OPERATIONS' },
  ];
  readonly dayOptions = [
    { label: 'Thứ hai', value: 1 },
    { label: 'Thứ ba', value: 2 },
    { label: 'Thứ tư', value: 3 },
    { label: 'Thứ năm', value: 4 },
    { label: 'Thứ sáu', value: 5 },
    { label: 'Thứ bảy', value: 6 },
    { label: 'Chủ nhật', value: 7 },
  ];
  readonly exceptionTypeOptions = [
    { label: 'Nghỉ lễ', value: 'HOLIDAY' },
    { label: 'Nghỉ bù', value: 'COMPENSATION_LEAVE' },
    { label: 'Làm bù', value: 'MAKEUP_WORKDAY' },
    { label: 'Đóng sớm', value: 'EARLY_CLOSE' },
  ];
  readonly adjustmentRuleOptions = [
    { label: 'Ngày làm việc kế tiếp', value: 'FOLLOWING' },
    { label: 'Ngày làm việc trước đó', value: 'PRECEDING' },
    { label: 'Không điều chỉnh', value: 'NONE' },
  ];

  readonly calendarForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    timezone: new FormControl('Asia/Ho_Chi_Minh', { nonNullable: true, validators: [Validators.required] }),
    calendarType: new FormControl('BANKING', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  readonly hourForm = new FormGroup({
    dayOfWeek: new FormControl(1, { nonNullable: true }),
    isWorkingDay: new FormControl(true, { nonNullable: true }),
    startTime: new FormControl('08:00', { nonNullable: true }),
    endTime: new FormControl('17:00', { nonNullable: true }),
    cutoffTime: new FormControl('15:30', { nonNullable: true }),
    sessionName: new FormControl('', { nonNullable: true }),
    sortOrder: new FormControl(10, { nonNullable: true }),
  });

  readonly exceptionForm = new FormGroup({
    date: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    exceptionType: new FormControl('HOLIDAY', { nonNullable: true }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    isWorkingDay: new FormControl(false, { nonNullable: true }),
    startTime: new FormControl('', { nonNullable: true }),
    endTime: new FormControl('', { nonNullable: true }),
    cutoffTime: new FormControl('', { nonNullable: true }),
    source: new FormControl('', { nonNullable: true }),
    note: new FormControl('', { nonNullable: true }),
  });

  readonly calculationForm = new FormGroup({
    startDate: new FormControl(this.formatDate(new Date()), { nonNullable: true, validators: [Validators.required] }),
    offsetDays: new FormControl(1, { nonNullable: true }),
    adjustmentRule: new FormControl('FOLLOWING', { nonNullable: true }),
  });

  constructor() {
    effect(() => {
      const calendars = this.calendars();
      if (!this.selectedCalendarId() && calendars.length > 0) {
        this.selectedCalendarId.set(calendars[0].id);
        this.reloadDetails();
      }
    });
  }

  reload(): void {
    this.resource.reload();
    this.reloadDetails();
  }

  selectCalendar(id: string): void {
    this.selectedCalendarId.set(id);
    this.selectedDay.set(null);
    this.reloadDetails();
  }

  previousMonth(): void {
    const month = this.currentMonth();
    this.currentMonth.set(new Date(month.getFullYear(), month.getMonth() - 1, 1));
  }

  nextMonth(): void {
    const month = this.currentMonth();
    this.currentMonth.set(new Date(month.getFullYear(), month.getMonth() + 1, 1));
  }

  openCalendar(item?: BusinessCalendar): void {
    this.selectedCalendar.set(item ?? null);
    this.calendarForm.reset(item ? { ...item } : {
      code: '',
      name: '',
      timezone: 'Asia/Ho_Chi_Minh',
      calendarType: 'BANKING',
      description: '',
      status: 'ACTIVE',
    });
    this.calendarDialogVisible.set(true);
  }

  saveCalendar(): void {
    if (this.calendarForm.invalid) {
      this.calendarForm.markAllAsTouched();
      return;
    }
    const selected = this.selectedCalendar();
    this.isSaving.set(true);
    const request = selected
      ? this.service.update(selected.id, this.calendarForm.getRawValue())
      : this.service.create(this.calendarForm.getRawValue());
    request.subscribe({
      next: item => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: selected ? 'Đã cập nhật lịch' : 'Đã tạo lịch' });
        this.calendarDialogVisible.set(false);
        this.selectedCalendarId.set(item.id);
        this.resource.reload();
        this.reloadDetails();
        this.isSaving.set(false);
      },
      error: () => this.handleSaveError(),
    });
  }

  deleteCalendar(item: BusinessCalendar): void {
    this.confirmationService.confirm({
      header: 'Xác nhận xóa',
      message: `Xóa lịch "${item.name}"?`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => this.service.delete(item.id).subscribe({
        next: () => {
          this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa lịch' });
          this.selectedCalendarId.set('');
          this.resource.reload();
          this.reloadDetails();
        },
        error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa lịch' }),
      }),
    });
  }

  openHour(item?: WorkingHour): void {
    this.selectedHour.set(item ?? null);
    this.hourForm.reset(item ? { ...item } : {
      dayOfWeek: 1,
      isWorkingDay: true,
      startTime: '08:00',
      endTime: '17:00',
      cutoffTime: '15:30',
      sessionName: '',
      sortOrder: 10,
    });
    this.hourDialogVisible.set(true);
  }

  saveHour(): void {
    const calendarId = this.selectedCalendarId();
    if (!calendarId) return;
    const selected = this.selectedHour();
    this.isSaving.set(true);
    const request = selected
      ? this.service.updateWorkingHour(selected.id, { ...this.hourForm.getRawValue(), calendarId })
      : this.service.createWorkingHour(calendarId, this.hourForm.getRawValue());
    request.subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã lưu giờ làm việc' });
        this.hourDialogVisible.set(false);
        this.reloadDetails();
        this.isSaving.set(false);
      },
      error: () => this.handleSaveError(),
    });
  }

  deleteHour(item: WorkingHour): void {
    this.service.deleteWorkingHour(item.id).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa giờ làm việc' });
        this.reloadDetails();
      },
      error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa giờ làm việc' }),
    });
  }

  openException(item?: CalendarException, date?: string): void {
    this.selectedException.set(item ?? null);
    this.exceptionForm.reset(item ? { ...item } : {
      date: date ?? this.formatDate(new Date()),
      exceptionType: 'HOLIDAY',
      name: '',
      isWorkingDay: false,
      startTime: '',
      endTime: '',
      cutoffTime: '',
      source: '',
      note: '',
    });
    this.exceptionDialogVisible.set(true);
  }

  saveException(): void {
    const calendarId = this.selectedCalendarId();
    if (!calendarId || this.exceptionForm.invalid) {
      this.exceptionForm.markAllAsTouched();
      return;
    }
    const selected = this.selectedException();
    this.isSaving.set(true);
    const request = selected
      ? this.service.updateException(selected.id, { ...this.exceptionForm.getRawValue(), calendarId })
      : this.service.createException(calendarId, this.exceptionForm.getRawValue());
    request.subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã lưu ngày ngoại lệ' });
        this.exceptionDialogVisible.set(false);
        this.reloadDetails();
        this.isSaving.set(false);
      },
      error: () => this.handleSaveError(),
    });
  }

  deleteException(item: CalendarException): void {
    this.service.deleteException(item.id).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa ngày ngoại lệ' });
        this.reloadDetails();
      },
      error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa ngày ngoại lệ' }),
    });
  }

  calculate(): void {
    const calendar = this.currentCalendar();
    if (!calendar || this.calculationForm.invalid) return;
    this.service.calculate({
      calendarId: calendar.id,
      ...this.calculationForm.getRawValue(),
    }).subscribe({
      next: result => {
        this.calculationResult.set(`${result.resultDate} (${result.calendarDays} ngày lịch, bỏ qua ${result.skippedDates.length} ngày)`);
      },
      error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tính ngày làm việc' }),
    });
  }

  selectDay(day: CalendarDay): void {
    this.selectedDay.set(day);
  }

  dayName(dayOfWeek: number): string {
    return this.dayOptions.find(item => item.value === dayOfWeek)?.label ?? '';
  }

  statusSeverity = statusSeverity;

  private reloadDetails(): void {
    const calendarId = this.selectedCalendarId();
    if (!calendarId) {
      this.workingHours.set([]);
      this.exceptions.set([]);
      return;
    }
    forkJoin({
      hours: this.service.listWorkingHours(calendarId),
      exceptions: this.service.listExceptions(calendarId),
    }).subscribe({
      next: result => {
        this.workingHours.set(result.hours);
        this.exceptions.set(result.exceptions);
      },
      error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tải cấu hình lịch' }),
    });
  }

  private buildMonthDays(): CalendarDay[] {
    const month = this.currentMonth();
    const start = new Date(month.getFullYear(), month.getMonth(), 1);
    const offset = (start.getDay() + 6) % 7;
    const gridStart = new Date(start);
    gridStart.setDate(start.getDate() - offset);
    const today = this.formatDate(new Date());
    return Array.from({ length: 42 }, (_, index) => {
      const date = new Date(gridStart);
      date.setDate(gridStart.getDate() + index);
      const value = this.formatDate(date);
      const exception = this.exceptions().find(item => item.date === value);
      return {
        date: value,
        label: date.getDate(),
        inMonth: date.getMonth() === month.getMonth(),
        isToday: value === today,
        isWorkingDay: exception ? exception.isWorkingDay : this.isWeekdayWorking(date),
        exception,
      };
    });
  }

  private isWeekdayWorking(date: Date): boolean {
    let weekday = date.getDay();
    if (weekday === 0) weekday = 7;
    return this.workingHours().find(item => item.dayOfWeek === weekday)?.isWorkingDay ?? false;
  }

  private monthStart(date: Date): Date {
    return new Date(date.getFullYear(), date.getMonth(), 1);
  }

  private formatDate(date: Date): string {
    const month = `${date.getMonth() + 1}`.padStart(2, '0');
    const day = `${date.getDate()}`.padStart(2, '0');
    return `${date.getFullYear()}-${month}-${day}`;
  }

  private handleSaveError(): void {
    this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu dữ liệu' });
    this.isSaving.set(false);
  }
}
