import { initFederation } from '@angular-architects/native-federation';

// Lấy tham số cấu hình từ env.js (được nạp vào runtime trước Angular)
const env = (window as any).__env || {};

// Thiết lập Map động trỏ đến các module con
const federationManifest = {
  iam: `${env.mfeIamUrl || 'http://localhost:4201'}/remoteEntry.json`,
  mdm: `${env.mfeMdmUrl || 'http://localhost:4202'}/remoteEntry.json`,
  ntf: `${env.mfeNtfUrl || 'http://localhost:4204'}/remoteEntry.json`,
  crm: `${env.mfeCrmUrl || 'http://localhost:4210'}/remoteEntry.json`,
  hrm: `${env.mfeHrmUrl || 'http://localhost:4211'}/remoteEntry.json`,
  loan: `${env.mfeLoanUrl || 'http://localhost:4212'}/remoteEntry.json`,
  bpm: `${env.mfeBpmUrl || 'http://localhost:4213'}/remoteEntry.json`
};

// Truyền thẳng Object vào initFederation (không dùng file json cứng nữa)
initFederation(federationManifest)
  .catch(err => console.error(err))
  .then(_ => import('./bootstrap'))
  .catch(err => console.error(err));
