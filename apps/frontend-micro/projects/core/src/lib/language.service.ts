import { Injectable, inject, signal } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { firstValueFrom } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class LanguageService {
  private translate = inject(TranslateService);
  private static readonly STORAGE_KEY = 'arda-lang';

  readonly currentLang = signal<string>(this.getInitialLang());

  constructor() { }

  /**
   * Khởi tạo ngôn ngữ đồng bộ với APP_INITIALIZER
   */
  async init(): Promise<void> {
    const lang = this.currentLang();
    this.translate.setDefaultLang('vi');

    try {
      const translations = await firstValueFrom(this.translate.use(lang));
      if (translations) {
        this.translate.setTranslation(lang, translations as any, true);
      }
      document.documentElement.lang = lang;
    } catch (err) {
      console.error('LanguageService: Could not load initial language', err);
      try {
        await firstValueFrom(this.translate.use('vi'));
      } catch (e) {
        /* ignore fallback failure */
      }
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
