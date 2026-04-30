export interface ListOptions {
  status?: string;
  keyword?: string;
  pageSize?: number;
  pageToken?: string;
  parentId?: string;
  level?: string;
  areaTypeId?: string;
  groupCode?: string;
}

export interface PageResponse<T> {
  items: T[];
  nextPageToken: string;
}

export interface AdministrativeUnit {
  id: string;
  code: string;
  name: string;
  fullName: string;
  shortName: string;
  level: string;
  unitType: string;
  parentId: string;
  path: string;
  sortOrder: number;
  latitude: number;
  longitude: number;
  status: string;
  effectiveFrom: string;
  effectiveTo: string;
  source: string;
  metadataJson: string;
}

export interface AdministrativeUnitNode {
  unit: AdministrativeUnit;
  children: AdministrativeUnitNode[];
}

export interface AdministrativeUnitSyncResult {
  provinceCount: number;
  wardCount: number;
  effectiveDate: string;
  source: string;
}

export interface AreaType {
  id: string;
  code: string;
  name: string;
  description: string;
  allowHierarchy: boolean;
  status: string;
}

export interface Area {
  id: string;
  areaTypeId: string;
  areaTypeCode: string;
  parentId: string;
  code: string;
  name: string;
  description: string;
  managerUserId: string;
  status: string;
  effectiveFrom: string;
  effectiveTo: string;
  metadataJson: string;
}

export interface AreaAdministrativeUnit {
  id: string;
  areaId: string;
  administrativeUnitId: string;
  scopeType: string;
  createdAt?: string;
}

export interface CodeSet {
  id: string;
  code: string;
  name: string;
  description: string;
  isSystem: boolean;
  status: string;
}

export interface CodeItem {
  id: string;
  codeSetId: string;
  codeSetCode: string;
  code: string;
  name: string;
  value: string;
  parentId: string;
  sortOrder: number;
  color: string;
  icon: string;
  metadataJson: string;
  isDefault: boolean;
  isSystem: boolean;
  status: string;
  effectiveFrom: string;
  effectiveTo: string;
}

export interface SystemParameter {
  id: string;
  key: string;
  name: string;
  groupCode: string;
  valueType: string;
  valueText: string;
  valueNumber: number;
  valueBoolean: boolean;
  valueJson: string;
  defaultValue: string;
  isSecret: boolean;
  isEditable: boolean;
  isSystem: boolean;
  validationRuleJson: string;
  description: string;
  status: string;
  updatedBy: string;
}

export interface CreditInstitution {
  id: string;
  code: string;
  name: string;
  shortName: string;
  address: string;
  phone: string;
  email: string;
  licenseNumber: string;
  issuedDate: string;
  taxCode: string;
  website: string;
  note: string;
  status: string;
}

export interface BusinessCalendar {
  id: string;
  code: string;
  name: string;
  timezone: string;
  calendarType: string;
  description: string;
  status: string;
}

export interface WorkingHour {
  id: string;
  calendarId: string;
  dayOfWeek: number;
  isWorkingDay: boolean;
  startTime: string;
  endTime: string;
  cutoffTime: string;
  sessionName: string;
  sortOrder: number;
}

export interface CalendarException {
  id: string;
  calendarId: string;
  date: string;
  exceptionType: string;
  name: string;
  isWorkingDay: boolean;
  startTime: string;
  endTime: string;
  cutoffTime: string;
  source: string;
  note: string;
}

export interface BusinessDayCalculation {
  calendarId?: string;
  calendarCode?: string;
  startDate: string;
  offsetDays: number;
  adjustmentRule: string;
}

export interface BusinessDayCalculationResult {
  startDate: string;
  resultDate: string;
  calendarDays: number;
  skippedDates: string[];
  isBusinessDay: boolean;
}
