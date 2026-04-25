import {
  Component,
  signal,
  HostListener,
  ViewEncapsulation,
  inject,
  effect,
} from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { ButtonModule } from 'primeng/button';
import { AvatarModule } from 'primeng/avatar';
import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';
import { MenuService } from '../../services/menu.service';
import { SidebarMenu } from './sidebar-menu/sidebar-menu';
import { AppHeader } from './header/header';

@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [RouterModule, CommonModule, ButtonModule, AvatarModule, SidebarMenu, AppHeader],
  templateUrl: './layout.html',
  styleUrl: './layout.css',
  encapsulation: ViewEncapsulation.None,
})
export class Layout {
  private authService = inject(AuthService);
  private menuService = inject(MenuService);
  private tenantService = inject(TenantService);

  userInitials = this.authService.userInitials;
  currentUser = this.authService.currentUser;
  menuItems = this.menuService.menuItems;

  sidebarVisible = signal(
    typeof window !== 'undefined' ? window.innerWidth >= 1024 : true,
  );

  constructor() {
    effect(() => {
      const tenantId = this.tenantService.selectedTenantId();
      if (tenantId) {
        this.menuService.loadMenu();
      }
    });
  }

  toggleSidebar(): void {
    this.sidebarVisible.update((v) => !v);
  }

  @HostListener('window:resize')
  onResize(): void {
    if (typeof window !== 'undefined') {
      if (window.innerWidth < 1024 && this.sidebarVisible()) {
        this.sidebarVisible.set(false);
      } else if (window.innerWidth >= 1024 && !this.sidebarVisible()) {
        this.sidebarVisible.set(true);
      }
    }
  }
}
