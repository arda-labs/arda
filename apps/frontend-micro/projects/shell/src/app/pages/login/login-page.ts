import { ChangeDetectionStrategy, Component, OnInit, inject, signal } from '@angular/core';
import { FormField, form, required, email } from '@angular/forms/signals';
import { Button } from 'primeng/button';
import { InputText } from 'primeng/inputtext';
import { Message } from 'primeng/message';
import { TranslateService, TranslatePipe } from '@ngx-translate/core';
import { AuthService } from '../../services/auth.service';
import { ZitadelSessionService } from '../../services/zitadel-session.service';
import { LanguageService, getAuthConfig } from '@arda/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-login-page',
  standalone: true,
  imports: [FormField, Button, InputText, Message, TranslatePipe],
  templateUrl: './login-page.html',
  styleUrl: './login-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginPage implements OnInit {
  private authService = inject(AuthService);
  private sessionService = inject(ZitadelSessionService);
  private langService = inject(LanguageService);
  private translate = inject(TranslateService);
  private route = inject(ActivatedRoute);

  authRequestId = signal<string | null>(null);
  isLoading = signal(false);
  showPassword = signal(false);
  errorMessage = signal<string | null>(null);
  currentLang = this.langService.currentLang;

  loginModel = signal({ loginName: '', password: '' });
  loginForm = form(this.loginModel, (field) => {
    required(field.loginName, { message: 'Email is required' });
    email(field.loginName, { message: 'Invalid email' });
    required(field.password, { message: 'Password is required' });
  });

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

  changeLang(lang: string) {
    this.langService.setLanguage(lang);
  }

  togglePassword(): void {
    this.showPassword.update((v) => !v);
  }

  onSubmit(): void {
    if (this.loginForm().invalid() || this.isLoading()) return;
    const authRequestId = this.authRequestId();
    if (!authRequestId) return;

    const { loginName, password } = this.loginModel();
    this.isLoading.set(true);
    this.errorMessage.set(null);

    this.sessionService.login(loginName, password, authRequestId).subscribe({
      next: ({ callbackUrl }) => {
        window.location.href = callbackUrl;
      },
      error: (err: unknown) => {
        this.isLoading.set(false);
        this.errorMessage.set(this.parseError(err));
      }
    });
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
