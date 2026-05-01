import { Component, ChangeDetectionStrategy, signal, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { TagModule } from 'primeng/tag';
import { CustomerService } from '../shared/services/customer.service';

@Component({
  selector: 'app-customer-list',
  standalone: true,
  imports: [CommonModule, TableModule, ButtonModule, InputTextModule, TagModule],
  templateUrl: './customer-list.html',
  styleUrl: './customer-list.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CustomerListComponent implements OnInit {
  private customerService = inject(CustomerService);
  private router = inject(Router);
  
  customers = signal<any[]>([]);
  loading = signal(false);

  ngOnInit() {
    this.loadData();
  }

  loadData() {
    this.loading.set(true);
    this.customerService.getCustomers().subscribe({
      next: (data) => {
        this.customers.set(data);
        this.loading.set(false);
      },
      error: () => {
        // Fallback to mock if API fails for demo
        this.customers.set([
          { id: 'C001', name: 'Nguyễn Văn A (Mock)', status: 'ACTIVE', lastUpdate: '2026-05-01' }
        ]);
        this.loading.set(false);
      }
    });
  }

  viewDetails(customer: any) {
    this.router.navigate(['/customer/info/details', customer.id]);
  }

  createNew() {
    this.router.navigate(['/customer/register/init']);
  }

  getStatusSeverity(status: string) {
    switch (status) {
      case 'ACTIVE': return 'success';
      case 'PENDING': return 'warn';
      default: return 'secondary';
    }
  }
}
