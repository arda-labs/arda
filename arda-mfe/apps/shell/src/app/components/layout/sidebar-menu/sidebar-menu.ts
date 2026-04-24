import { Component, input, signal, ViewEncapsulation, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';

export interface MenuItem {
  label: string;
  icon?: string;
  routerLink?: string[];
  items?: MenuItem[];
}

@Component({
  selector: 'app-sidebar-menu',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './sidebar-menu.html',
  styleUrl: './sidebar-menu.css',
  encapsulation: ViewEncapsulation.None,
})
export class SidebarMenu implements OnInit {
  items = input.required<MenuItem[]>();
  expandedItems = signal<Set<string>>(new Set());
  private router = inject(Router);

  ngOnInit(): void {
    // Tự động mở menu cha nếu có mục con đang active
    this.autoExpandActiveItems(this.items());
  }

  private autoExpandActiveItems(items: MenuItem[]): boolean {
    let hasActiveChild = false;

    for (const item of items) {
      const isCurrentRoute = item.routerLink ? this.router.isActive(this.router.createUrlTree(item.routerLink), {
        paths: 'subset',
        queryParams: 'ignored',
        fragment: 'ignored',
        matrixParams: 'ignored'
      }) : false;

      let childHasActive = false;
      if (item.items) {
        childHasActive = this.autoExpandActiveItems(item.items);
      }

      if (isCurrentRoute || childHasActive) {
        if (item.items && item.items.length > 0) {
          this.expandedItems.update(set => new Set(set).add(item.label));
        }
        hasActiveChild = true;
      }
    }
    return hasActiveChild;
  }

  toggleMenuItem(item: MenuItem): void {
    if (!item.items || item.items.length === 0) return;

    this.expandedItems.update(set => {
      const newSet = new Set(set);
      if (newSet.has(item.label)) {
        newSet.delete(item.label);
      } else {
        newSet.add(item.label);
      }
      return newSet;
    });
  }

  isExpanded(label: string): boolean {
    return this.expandedItems().has(label);
  }
}
