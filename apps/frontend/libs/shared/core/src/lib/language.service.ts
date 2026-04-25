import { Injectable, inject, signal } from '@angular/core';
import { registerLocaleData } from '@angular/common';
import localeVi from '@angular/common/locales/vi';
import { TranslateService } from '@ngx-translate/core';
import { firstValueFrom } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class LanguageService {
  private translate = inject(TranslateService);
  private static readonly STORAGE_KEY = 'arda-lang';

  readonly currentLang = signal<string>(this.getInitialLang());

  constructor() {
    registerLocaleData(localeVi);
  }

  /**
   * Khởi tạo ngôn ngữ đồng bộ với APP_INITIALIZER
   */
  async init(): Promise<void> {
    const lang = this.currentLang();
    try {
      await firstValueFrom(this.translate.use(lang));
      document.documentElement.lang = lang;
    } catch (err) {
      console.error('LanguageService: Could not load initial language', err);
    }
  }

  setLanguage(lang: string) {
    this.translate.use(lang).subscribe({
      next: () => {
        this.currentLang.set(lang);
        localStorage.setItem(LanguageService.STORAGE_KEY, lang);
        document.documentElement.lang = lang;
        window.dispatchEvent(new CustomEvent('arda-lang-changed', { detail: lang }));
      }
    });
  }

  private getInitialLang(): string {
    const saved = localStorage.getItem(LanguageService.STORAGE_KEY);
    if (saved) return saved;
    return 'vi';
  }
}
