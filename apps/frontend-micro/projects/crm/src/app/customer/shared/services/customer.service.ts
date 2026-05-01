import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class CustomerService {
  private http = inject(HttpClient);
  private readonly API_URL = '/api/v1/customers';

  getCustomers(): Observable<any[]> {
    return this.http.get<any[]>(`${this.API_URL}`);
  }

  getCustomerDetails(id: string): Observable<any> {
    return this.http.get<any>(`${this.API_URL}/${id}`);
  }

  createRequest(data: any): Observable<any> {
    return this.http.post(`${this.API_URL}/register`, data);
  }
}
