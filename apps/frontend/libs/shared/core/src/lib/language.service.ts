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
    // Đảm bảo có ngôn ngữ mặc định để fallback
    this.translate.setDefaultLang('vi');

    console.log(`LanguageService: Initializing with language "${lang}"...`);

    try {
      // Đảm bảo load bản dịch trước khi ứng dụng render
      const translations = await firstValueFrom(this.translate.use(lang));

      if (translations) {
        // Force set translation để đảm bảo store được cập nhật
        this.translate.setTranslation(lang, translations as any, true);
        console.log(`LanguageService: Successfully loaded language "${lang}"`, Object.keys(translations as any));

        // Kiểm tra xem instant có lấy được giá trị không
        const testVal = this.translate.instant('PAGES.LOGIN.TITLE');
        console.log(`LanguageService: Validation check for 'PAGES.LOGIN.TITLE': ${testVal}`);
      }

      document.documentElement.lang = lang;
    } catch (err) {
      console.error('LanguageService: Could not load initial language', err);
      // Fallback sang tiếng Việt nếu load thất bại
      try {
        await firstValueFrom(this.translate.use('vi'));
      } catch (e) {
        console.error('LanguageService: Fallback to "vi" failed as well', e);
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
