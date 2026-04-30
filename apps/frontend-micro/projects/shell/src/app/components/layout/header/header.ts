import { Component, input, output, inject, ChangeDetectionStrategy, computed, signal, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NavigationEnd, Router, RouterModule } from '@angular/router';
import { Button } from 'primeng/button';
import { Avatar } from 'primeng/avatar';
import { Tooltip } from 'primeng/tooltip';
import { Menu } from 'primeng/menu';
import { MenuItem as PrimeMenuItem } from 'primeng/api';
import { Breadcrumb } from 'primeng/breadcrumb';
import { AuthService } from '../../../services/auth.service';
import { MenuItem as ShellMenuItem, MenuService } from '../../../services/menu.service';
import { LanguageService, ThemeService } from '@arda/core';
import { filter } from 'rxjs';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterModule, Button, Avatar, Tooltip, Menu, Breadcrumb],
  templateUrl: './header.html',
  styleUrl: './header.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AppHeader implements OnDestroy {
  private authService = inject(AuthService);
  private langService = inject(LanguageService);
  private themeService = inject(ThemeService);
  private menuService = inject(MenuService);
  private router = inject(Router);

  userInitials = this.authService.userInitials;
  userName = this.authService.currentUser;
  currentLang = this.langService.currentLang;
  themeSettings = this.themeService.settings;
  private currentUrl = signal(this.normalizeUrl(this.router.url));
  private now = signal(new Date());
  private clockTimer: ReturnType<typeof setInterval> | null = null;

  sidebarVisible = input<boolean>(true);
  toggleSidebar = output<void>();

  breadcrumbHome: PrimeMenuItem = {
    icon: 'pi pi-home',
    routerLink: '/home',
  };

  breadcrumbItems = computed<PrimeMenuItem[]>(() => {
    const url = this.currentUrl();
    const trail = this.findMenuTrail(this.menuService.menuItems(), url);

    if (trail.length > 0) {
      return trail.map((item, index) => ({
        label: item.label,
        routerLink: item.routerLink?.length ? item.routerLink[0] : undefined,
      }));
    }

    return this.fallbackBreadcrumb(url);
  });

  constructor() {
    this.router.events.pipe(filter((event): event is NavigationEnd => event instanceof NavigationEnd)).subscribe(event => {
      this.currentUrl.set(this.normalizeUrl(event.urlAfterRedirects));
    });

    this.clockTimer = setInterval(() => {
      this.now.set(new Date());
    }, 1000);
  }

  marketClock = computed(() => {
    const date = this.now();
    const day = new Intl.DateTimeFormat('en-GB', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    }).format(date);
    const time = new Intl.DateTimeFormat('en-GB', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    }).format(date);

    return `${day} ${time} ICT`;
  });

  currentFlag = computed(() => this.currentLang() === 'vi' ? '🇻🇳' : '🇺🇸');

  languageMenuItems: PrimeMenuItem[] = [
    {
      label: '🇻🇳 Tiếng Việt',
      command: () => this.changeLang('vi'),
    },
    {
      label: '🇺🇸 English',
      command: () => this.changeLang('en'),
    },
  ];

  ngOnDestroy(): void {
    if (this.clockTimer) {
      clearInterval(this.clockTimer);
    }
  }

  changeLang(lang: string) {
    this.langService.setLanguage(lang);
  }

  toggleDarkMode(): void {
    this.themeService.updateSetting('darkMode', !this.themeSettings().darkMode);
  }

  userMenuItems: PrimeMenuItem[] = [
    {
      label: 'Profile',
      icon: 'pi pi-user',
      routerLink: '/iam/profile'
    },
    {
      separator: true
    },
    {
      label: 'Logout',
      icon: 'pi pi-sign-out',
      command: () => {
        this.authService.logout();
      }
    }
  ];

  onToggleSidebar() {
    this.toggleSidebar.emit();
  }

  private findMenuTrail(items: ShellMenuItem[], url: string, parents: ShellMenuItem[] = []): ShellMenuItem[] {
    let bestTrail: ShellMenuItem[] = [];

    for (const item of items) {
      const itemUrl = this.normalizeMenuUrl(item.routerLink);
      const currentTrail = [...parents, item];

      if (itemUrl && (url === itemUrl || url.startsWith(`${itemUrl}/`))) {
        bestTrail = currentTrail;
      }

      const childTrail = this.findMenuTrail(item.items ?? [], url, currentTrail);
      if (childTrail.length > bestTrail.length) {
        bestTrail = childTrail;
      }
    }

    return bestTrail;
  }

  private fallbackBreadcrumb(url: string): PrimeMenuItem[] {
    const labels: Record<string, string> = {
      '/home': 'Dashboard',
      '/settings': 'Cài đặt',
      '/workspaces': 'Workspace',
      '/iam/profile': 'Thông tin cá nhân',
    };

    const label = labels[url] ?? this.toTitle(url.split('/').filter(Boolean).at(-1) ?? 'Trang');
    return [{ label, routerLink: url }];
  }

  private normalizeMenuUrl(routerLink?: string[]): string {
    return this.normalizeUrl(routerLink?.[0] ?? '');
  }

  private normalizeUrl(url: string): string {
    const [path] = url.split(/[?#]/);
    const normalized = path.startsWith('/') ? path : `/${path}`;
    return normalized.length > 1 ? normalized.replace(/\/$/, '') : normalized;
  }

  private toTitle(value: string): string {
    return value
      .split(/[-_]/)
      .filter(Boolean)
      .map(part => part.charAt(0).toUpperCase() + part.slice(1))
      .join(' ');
  }
}
