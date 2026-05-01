const fs = require('fs');

const tsconfigPath = './tsconfig.json';
const data = JSON.parse(fs.readFileSync(tsconfigPath, 'utf8'));

const appsToAdd = ['bpm', 'loan', 'hrm'];

appsToAdd.forEach(app => {
  const appRef = { "path": `./projects/${app}/tsconfig.app.json` };
  const specRef = { "path": `./projects/${app}/tsconfig.spec.json` };
  
  if (!data.references.some(r => r.path === appRef.path)) {
    data.references.push(appRef);
  }
  if (!data.references.some(r => r.path === specRef.path)) {
    data.references.push(specRef);
  }
});

fs.writeFileSync(tsconfigPath, JSON.stringify(data, null, 2));
console.log('Updated tsconfig.json successfully.');
