import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

interface LoginResponse {
  callbackUrl: string;
}

@Injectable({ providedIn: 'root' })
export class ZitadelSessionService {
  private http = inject(HttpClient);
  private apiUrl = (window as any).__env?.apiUrl ?? 'http://localhost:8000';

  /**
   * Gọi backend (iam-service) để thực hiện flow login trọn gói với Zitadel
   */
  login(email: string, password: string, authRequestId: string): Observable<LoginResponse> {
    const url = `${this.apiUrl}/v1/auth/login`;
    return this.http.post<LoginResponse>(url, {
      email,
      password,
      auth_request_id: authRequestId,
    });
  }
}
