import { CommonModule } from '@angular/common';
import {
  AfterContentInit,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  ContentChildren,
  EventEmitter,
  inject,
  Input,
  Output,
  QueryList,
  TemplateRef,
} from '@angular/core';
import { PrimeTemplate } from 'primeng/api';
import { TableModule } from 'primeng/table';
import { DEFAULT_ROWS_PER_PAGE_OPTIONS, LazyPageEvent } from '@arda/core';

@Component({
  selector: 'arda-data-table',
  standalone: true,
  imports: [CommonModule, TableModule],
  template: `
    <div class="bg-surface-0 dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 overflow-hidden">
      <p-table
        [value]="value"
        [loading]="loading"
        [lazy]="lazy"
        [lazyLoadOnInit]="lazyLoadOnInit"
        [first]="first"
        [rows]="rows"
        [totalRecords]="totalRecords"
        [paginator]="paginator"
        [rowsPerPageOptions]="rowsPerPageOptions"
        [showFirstLastIcon]="showFirstLastIcon"
        [showPageLinks]="showPageLinks"
        [responsiveLayout]="responsiveLayout"
        [styleClass]="styleClass"
        (onLazyLoad)="lazyLoad.emit($event)"
      >
        @if (captionTemplate) {
          <ng-template pTemplate="caption">
            <ng-container *ngTemplateOutlet="captionTemplate" />
          </ng-template>
        }

        @if (headerTemplate) {
          <ng-template pTemplate="header" let-columns>
            <ng-container *ngTemplateOutlet="headerTemplate; context: { $implicit: columns, columns: columns }" />
          </ng-template>
        }

        @if (bodyTemplate) {
          <ng-template pTemplate="body" let-rowData let-columns="columns" let-rowIndex="rowIndex">
            <ng-container
              *ngTemplateOutlet="
                bodyTemplate;
                context: { $implicit: rowData, rowData: rowData, columns: columns, rowIndex: rowIndex }
              "
            />
          </ng-template>
        }

        @if (emptyMessageTemplate) {
          <ng-template pTemplate="emptymessage" let-columns>
            <ng-container *ngTemplateOutlet="emptyMessageTemplate; context: { $implicit: columns, columns: columns }" />
          </ng-template>
        }

        @if (footerTemplate) {
          <ng-template pTemplate="footer" let-columns>
            <ng-container *ngTemplateOutlet="footerTemplate; context: { $implicit: columns, columns: columns }" />
          </ng-template>
        }
      </p-table>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ArdaDataTable<T = unknown> implements AfterContentInit {
  private cd = inject(ChangeDetectorRef);

  @ContentChildren(PrimeTemplate) templates?: QueryList<PrimeTemplate>;

  @Input() value: T[] = [];
  @Input() loading = false;
  @Input() lazy = true;
  @Input() lazyLoadOnInit = false;
  @Input() first = 0;
  @Input() rows = 20;
  @Input() totalRecords = 0;
  @Input() paginator = true;
  @Input() rowsPerPageOptions: number[] = [...DEFAULT_ROWS_PER_PAGE_OPTIONS];
  @Input() showFirstLastIcon = false;
  @Input() showPageLinks = false;
  @Input() responsiveLayout: 'stack' | 'scroll' = 'scroll';
  @Input() styleClass = 'p-datatable-sm';

  @Output() lazyLoad = new EventEmitter<LazyPageEvent>();

  captionTemplate?: TemplateRef<unknown>;
  headerTemplate?: TemplateRef<unknown>;
  bodyTemplate?: TemplateRef<unknown>;
  emptyMessageTemplate?: TemplateRef<unknown>;
  footerTemplate?: TemplateRef<unknown>;

  ngAfterContentInit(): void {
    this.captureTemplates();
    this.templates?.changes.subscribe(() => this.captureTemplates());
  }

  private captureTemplates(): void {
    this.captionTemplate = undefined;
    this.headerTemplate = undefined;
    this.bodyTemplate = undefined;
    this.emptyMessageTemplate = undefined;
    this.footerTemplate = undefined;

    this.templates?.forEach((item) => {
      switch (item.getType()) {
        case 'caption':
          this.captionTemplate = item.template;
          break;
        case 'header':
          this.headerTemplate = item.template;
          break;
        case 'body':
          this.bodyTemplate = item.template;
          break;
        case 'emptymessage':
          this.emptyMessageTemplate = item.template;
          break;
        case 'footer':
          this.footerTemplate = item.template;
          break;
      }
    });

    this.cd.markForCheck();
  }
}
