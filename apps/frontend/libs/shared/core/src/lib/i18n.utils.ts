import { HttpClient } from '@angular/common/http';
import { TranslateLoader } from '@ngx-translate/core';
import { Observable } from 'rxjs';
import { getEnv } from '@arda-mfe/shared-config';

export class CustomTranslateLoader implements TranslateLoader {
  constructor(private http: HttpClient) {}
  getTranslation(lang: string): Observable<any> {
    // Thêm tham số version để tránh cache nhầm khi deploy bản mới (Cache Busting)
    const env = getEnv();
    const v = env.version ?? Date.now().toString();
    return this.http.get<any>(`/i18n/${lang}.json?v=${v}`);
  }
}

export function HttpLoaderFactory(http: HttpClient) {
  return new CustomTranslateLoader(http);
}
