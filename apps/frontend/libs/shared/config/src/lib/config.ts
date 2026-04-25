export interface Environment {
  apiUrl: string;
  authClientId: string;
  version: string;
  production: boolean;
}

declare global {
  interface Window {
    __env?: Partial<Environment>;
  }
}

/**
 * Các giá trị mặc định dùng cho môi trường local dev
 */
const DEFAULT_ENV: Environment = {
  apiUrl: 'http://localhost:8000/api',
  authClientId: 'ZITADEL_SPA_CLIENT_ID',
  version: 'local',
  production: false,
};

/**
 * Hàm duy nhất để đọc biến môi trường trong toàn bộ hệ thống.
 * Kết hợp giữa window.__env (từ K8s/env.js) và các giá trị mặc định.
 */
export function getEnv(): Environment {
  const windowEnv = window.__env ?? {};

  return {
    ...DEFAULT_ENV,
    ...windowEnv,
    // Đảm bảo các field quan trọng luôn có giá trị
    apiUrl: windowEnv.apiUrl ?? DEFAULT_ENV.apiUrl,
    authClientId: windowEnv.authClientId ?? DEFAULT_ENV.authClientId,
  };
}
