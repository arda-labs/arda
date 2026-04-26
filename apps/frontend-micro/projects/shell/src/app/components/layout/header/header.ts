import { Component, input, output, inject, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { Button } from 'primeng/button';
import { Avatar } from 'primeng/avatar';
import { Tooltip } from 'primeng/tooltip';
import { Menu } from 'primeng/menu';
import { MenuItem } from 'primeng/api';
import { AuthService } from '../../../services/auth.service';
import { LanguageService } from '@arda/core';
import { TenantSwitcher } from './tenant-switcher/tenant-switcher';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterModule, Button, Avatar, Tooltip, Menu, TenantSwitcher],
  templateUrl: './header.html',
  styleUrl: './header.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AppHeader {
  private authService = inject(AuthService);
  private langService = inject(LanguageService);

  userInitials = this.authService.userInitials;
  userName = this.authService.currentUser;
  currentLang = this.langService.currentLang;

  sidebarVisible = input<boolean>(true);
  toggleSidebar = output<void>();

  changeLang(lang: string) {
    this.langService.setLanguage(lang);
  }

  userMenuItems: MenuItem[] = [
    {
      label: 'Profile',
      icon: 'pi pi-user',
      routerLink: '/app/profile'
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
}
