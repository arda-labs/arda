import { initFederation } from '@angular-architects/native-federation';

// Lấy tham số cấu hình từ env.js (được nạp vào runtime trước Angular)
const env = (window as any).__env || {};

// Thiết lập Map động trỏ đến các module con
const federationManifest = {
  iam: `${env.mfeIamUrl || 'http://localhost:3002'}/remoteEntry.json`,
  common: `${env.mfeCommonUrl || 'http://localhost:3001'}/remoteEntry.json`
};

// Truyền thẳng Object vào initFederation (không dùng file json cứng nữa)
initFederation(federationManifest)
  .catch(err => console.error(err))
  .then(_ => import('./bootstrap'))
  .catch(err => console.error(err));
