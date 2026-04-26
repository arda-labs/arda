import { Component, OnInit, inject, signal } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { MessageModule } from 'primeng/message';
import { TranslateModule, TranslateService, TranslatePipe } from '@ngx-translate/core';
import { AuthService } from '../../services/auth.service';
import { ZitadelSessionService } from '../../services/zitadel-session.service';
import { LanguageService, getAuthConfig } from '@arda-mfe/shared-core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-login-page',
  standalone: true,
  imports: [ReactiveFormsModule, ButtonModule, InputTextModule, MessageModule, TranslatePipe],
  templateUrl: './login-page.html',
  styleUrl: './login-page.css',
})
export class LoginPage implements OnInit {
  private authService = inject(AuthService);
  private sessionService = inject(ZitadelSessionService);
  private langService = inject(LanguageService);
  private translate = inject(TranslateService);
  private route = inject(ActivatedRoute);
  private fb = inject(FormBuilder);

  authRequestId = signal<string | null>(null);
  isLoading = signal(false);
  showPassword = signal(false);
  errorMessage = signal<string | null>(null);
  currentLang = this.langService.currentLang;

  changeLang(lang: string) {
    this.langService.setLanguage(lang);
  }

  readonly features = [
    { icon: 'pi-users', text: 'CRM, HRM, Finance trên một nền tảng' },
    { icon: 'pi-shield', text: 'Bảo mật chuẩn ngân hàng — RBAC + 2FA' },
    { icon: 'pi-building', text: 'Multi-tenant, phân quyền theo tổ chức' },
    { icon: 'pi-bolt', text: 'Tối ưu hóa hiệu năng và trải nghiệm người dùng' },
  ];

  readonly badges = [
    { icon: 'pi-lock', text: 'End-to-end' },
    { icon: 'pi-shield', text: '2FA' },
    { icon: 'pi-verified', text: 'SSO' },
  ];

  readonly form = this.fb.group({
    loginName: ['', [Validators.required, Validators.email]],
    password: ['', Validators.required],
  });

  get forgotPasswordUrl(): string {
    return `${getAuthConfig().authority}/ui/login/password/reset`;
  }

  ngOnInit(): void {
    const authRequest = this.route.snapshot.queryParams['authRequest'];
    if (authRequest) {
      this.authRequestId.set(authRequest);
    } else {
      this.authService.login();
    }
  }

  togglePassword(): void {
    this.showPassword.update((v) => !v);
  }

  async onSubmit(): Promise<void> {
    if (this.form.invalid || this.isLoading()) return;
    const authRequestId = this.authRequestId();
    if (!authRequestId) return;

    const { loginName, password } = this.form.value;
    this.isLoading.set(true);
    this.errorMessage.set(null);

    try {
      const { callbackUrl } = await this.sessionService.login(
        loginName!,
        password!,
        authRequestId,
      );
      window.location.href = callbackUrl;
    } catch (err: unknown) {
      this.isLoading.set(false);
      this.errorMessage.set(this.parseError(err));
    }
  }

  private parseError(err: unknown): string {
    if (err && typeof err === 'object' && 'status' in err) {
      const s = (err as { status: number }).status;
      if (s === 400 || s === 401) return this.translate.instant('PAGES.LOGIN.LOGIN_FAILED');
      if (s === 404) return this.translate.instant('PAGES.LOGIN.ACCOUNT_NOT_FOUND');
      if (s === 429) return this.translate.instant('PAGES.LOGIN.TOO_MANY_REQUESTS');
    }
    return this.translate.instant('COMMON.ERROR.SYSTEM');
  }
}
