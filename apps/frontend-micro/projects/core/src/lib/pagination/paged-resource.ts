import { computed, signal } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { Observable, of } from 'rxjs';

export const DEFAULT_PAGE_SIZE = 20;
export const DEFAULT_ROWS_PER_PAGE_OPTIONS = [10, 20, 50, 100] as const;

export interface PageRequest {
  pageSize?: number;
  pageToken?: string;
}

export interface PageResponse<T> {
  items: T[];
  nextPageToken: string;
  totalCount?: number;
}

export interface LazyPageEvent {
  first?: number | null;
  rows?: number | null;
}

export interface PagedResourceOptions<TItem, TParams> {
  params: () => TParams | null | undefined;
  load: (params: TParams, page: Required<PageRequest>) => Observable<PageResponse<TItem>>;
  defaultPageSize?: number;
  rowsPerPageOptions?: readonly number[];
}

export function emptyPage<T>(): PageResponse<T> {
  return {
    items: [],
    nextPageToken: '',
  };
}

export function createPagedResource<TItem, TParams>(options: PagedResourceOptions<TItem, TParams>) {
  const pageSize = signal(options.defaultPageSize ?? DEFAULT_PAGE_SIZE);
  const pageIndex = signal(0);
  const pageToken = signal('');
  const pageTokens = signal<string[]>(['']);
  const refreshKey = signal(0);
  const rowsPerPageOptions = [...(options.rowsPerPageOptions ?? DEFAULT_ROWS_PER_PAGE_OPTIONS)];

  const resource = rxResource({
    params: () => ({
      source: options.params(),
      pageSize: pageSize(),
      pageToken: pageToken(),
      refreshKey: refreshKey(),
    }),
    stream: ({ params }) => {
      if (params.source == null || params.source === '') {
        return of(emptyPage<TItem>());
      }

      return options.load(params.source, {
        pageSize: params.pageSize,
        pageToken: params.pageToken,
      });
    },
  });

  const items = computed(() => resource.value()?.items ?? []);
  const first = computed(() => pageIndex() * pageSize());
  const totalRecords = computed(() => {
    const page = resource.value();
    if (page?.totalCount != null) {
      return page.totalCount;
    }

    const loaded = first() + (page?.items.length ?? 0);
    return page?.nextPageToken ? loaded + pageSize() : loaded;
  });

  const reset = () => {
    pageTokens.set(['']);
    pageIndex.set(0);
    pageToken.set('');
  };

  const loadPage = (event: LazyPageEvent) => {
    const rows = event.rows ?? pageSize();
    const firstRow = event.first ?? 0;
    const nextPageIndex = Math.floor(firstRow / rows);

    if (rows !== pageSize()) {
      pageSize.set(rows);
      reset();
      return;
    }

    if (nextPageIndex === pageIndex()) {
      return;
    }

    if (nextPageIndex === 0) {
      pageIndex.set(0);
      pageToken.set('');
      return;
    }

    const tokens = [...pageTokens()];
    const currentPage = resource.value();
    if (nextPageIndex === pageIndex() + 1 && currentPage?.nextPageToken) {
      tokens[nextPageIndex] = currentPage.nextPageToken;
      pageTokens.set(tokens);
      pageIndex.set(nextPageIndex);
      pageToken.set(currentPage.nextPageToken);
      return;
    }

    const token = tokens[nextPageIndex];
    if (token !== undefined) {
      pageIndex.set(nextPageIndex);
      pageToken.set(token);
    }
  };

  const refresh = () => {
    reset();
    refreshKey.update((value) => value + 1);
  };

  return {
    resource,
    items,
    first,
    totalRecords,
    pageIndex,
    pageSize,
    pageToken,
    refreshKey,
    rowsPerPageOptions,
    loadPage,
    refresh,
    reset,
  };
}
