---
name: e2e-testing
description: Hỗ trợ viết E2E tests
disable-model-invocation: true

---
# End-to-End Testing Skill

Mục đích: Hỗ trợ viết E2E tests cho Angular applications trong dự án Arda.

## 🎯 Phạm vi

- Setup Playwright
- Tạo E2E test scenarios
- Run E2E tests
- Generate E2E test reports
- Record E2E test videos

## 📦 Setup Playwright

### Install Playwright

```bash
# Install Playwright
npm install -D @playwright/test

# Install browsers
npx playwright install

# Install browsers (with system dependencies)
npx playwright install-deps

# Verify installation
npx playwright test --version
```

### Playwright Configuration

```typescript
// playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,

  use: {
    baseURL: 'http://localhost:4200',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
  ],

  reporter: [
    ['html'],
    ['json', { outputFile: 'test-results/results.json' }],
    ['junit', { outputFile: 'test-results/junit.xml' }],
  ],
});
```

## 📦 Page Object Model

### Base Page

```typescript
// e2e/pages/base-page.ts
import { Page, Locator } from '@playwright/test';

export class BasePage {
  readonly page: Page;
  readonly url: string;

  constructor(page: Page, url: string) {
    this.page = page;
    this.url = url;
  }

  async goto() {
    await this.page.goto(this.url);
  }

  async refresh() {
    await this.page.reload();
  }

  async waitForLoad() {
    await this.page.waitForLoadState('networkidle');
  }

  async takeScreenshot(filename: string) {
    await this.page.screenshot({ path: `screenshots/${filename}` });
  }

  async fillInput(locator: Locator, value: string) {
    await locator.fill(value);
  }

  async clickElement(locator: Locator) {
    await locator.click();
  }

  async getText(locator: Locator): Promise<string> {
    return await locator.textContent() || '';
  }

  async isVisible(locator: Locator): Promise<boolean> {
    return await locator.isVisible();
  }
}
```

### Login Page

```typescript
// e2e/pages/login.page.ts
import { Page } from '@playwright/test';
import { BasePage } from './base-page';

export class LoginPage extends BasePage {
  readonly emailInput: Locator;
  readonly passwordInput: Locator;
  readonly loginButton: Locator;
  readonly errorMessage: Locator;

  constructor(page: Page) {
    super(page, '/auth/login');
    this.emailInput = page.getByPlaceholder('Email');
    this.passwordInput = page.getByPlaceholder('Password');
    this.loginButton = page.getByRole('button', { name: 'Login' });
    this.errorMessage = page.getByTestId('error-message');
  }

  async login(email: string, password: string) {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.loginButton.click();
    await this.page.waitForURL('**/dashboard');
  }

  async getErrorMessage(): Promise<string> {
    return await this.errorMessage.textContent() || '';
  }
}
```

### Journal List Page

```typescript
// e2e/pages/journal-list.page.ts
import { Page } from '@playwright/test';
import { BasePage } from './base-page';

export class JournalListPage extends BasePage {
  readonly journalTable: Locator;
  readonly createButton: Locator;
  readonly searchInput: Locator;
  readonly statusFilter: Locator;
  readonly loadingIndicator: Locator;

  constructor(page: Page) {
    super(page, '/accounting/journals');
    this.journalTable = page.getByRole('table');
    this.createButton = page.getByRole('button', { name: 'Create Journal' });
    this.searchInput = page.getByPlaceholder('Search...');
    this.statusFilter = page.getByRole('combobox', { name: 'Status' });
    this.loadingIndicator = page.getByTestId('loading-indicator');
  }

  async waitForJournals() {
    await this.journalTable.waitFor({ state: 'visible' });
  }

  async getJournalCount(): Promise<number> {
    const rows = await this.journalTable.locator('tbody tr').count();
    return rows;
  }

  async searchJournals(query: string) {
    await this.searchInput.fill(query);
    await this.waitForLoad();
  }

  async filterByStatus(status: string) {
    await this.statusFilter.click();
    await this.page.getByRole('option', { name: status }).click();
    await this.waitForLoad();
  }

  async clickCreateJournal() {
    await this.createButton.click();
  }
}
```

### Journal Form Page

```typescript
// e2e/pages/journal-form.page.ts
import { Page } from '@playwright/test';
import { BasePage } from './base-page';

export class JournalFormPage extends BasePage {
  readonly journalNoInput: Locator;
  readonly journalDateInput: Locator;
  readonly descriptionInput: Locator;
  readonly statusSelect: Locator;
  readonly saveButton: Locator;
  readonly cancelButton: Locator;
  readonly validationErrors: Locator;

  constructor(page: Page) {
    super(page, '/accounting/journals/new');
    this.journalNoInput = page.getByLabel('Journal No');
    this.journalDateInput = page.getByLabel('Journal Date');
    this.descriptionInput = page.getByLabel('Description');
    this.statusSelect = page.getByLabel('Status');
    this.saveButton = page.getByRole('button', { name: 'Save' });
    this.cancelButton = page.getByRole('button', { name: 'Cancel' });
    this.validationErrors = page.getByTestId('validation-error');
  }

  async fillJournalForm(data: {
    journalNo: string;
    journalDate: string;
    description: string;
    status?: string;
  }) {
    await this.journalNoInput.fill(data.journalNo);
    await this.journalDateInput.fill(data.journalDate);
    await this.descriptionInput.fill(data.description);

    if (data.status) {
      await this.statusSelect.click();
      await this.page.getByRole('option', { name: data.status }).click();
    }
  }

  async saveJournal() {
    await this.saveButton.click();
  }

  async cancel() {
    await this.cancelButton.click();
  }

  async getValidationErrors(): Promise<string[]> {
    const errors = await this.validationErrors.allTextContents();
    return errors;
  }
}
```

## 📦 E2E Test Scenarios

### Authentication Flow

```typescript
// e2e/auth/auth.spec.ts
import { test, expect } from '@playwright/test';
import { LoginPage } from '../pages/login.page';

test.describe('Authentication', () => {
  let loginPage: LoginPage;

  test.beforeEach(async ({ page }) => {
    loginPage = new LoginPage(page);
    await loginPage.goto();
  });

  test('should login with valid credentials', async () => {
    await loginPage.login('user@example.com', 'password123');

    await expect(page).toHaveURL('**/dashboard');
    await expect(page.getByText('Welcome')).toBeVisible();
  });

  test('should show error with invalid credentials', async () => {
    await loginPage.emailInput.fill('invalid@example.com');
    await loginPage.passwordInput.fill('wrongpassword');
    await loginPage.loginButton.click();

    const errorMessage = await loginPage.getErrorMessage();
    expect(errorMessage).toContain('Invalid credentials');
  });

  test('should validate email format', async () => {
    await loginPage.emailInput.fill('invalid-email');
    await loginPage.passwordInput.fill('password123');
    await loginPage.loginButton.click();

    await expect(loginPage.emailInput).toHaveAttribute('aria-invalid', 'true');
  });
});
```

### Journal Management Flow

```typescript
// e2e/accounting/journals.spec.ts
import { test, expect } from '@playwright/test';
import { JournalListPage } from '../pages/journal-list.page';
import { JournalFormPage } from '../pages/journal-form.page';

test.describe('Journal Management', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login('user@example.com', 'password123');
  });

  test('should list journals', async ({ page }) => {
    const journalListPage = new JournalListPage(page);
    await journalListPage.goto();

    await journalListPage.waitForJournals();

    const count = await journalListPage.getJournalCount();
    expect(count).toBeGreaterThan(0);
  });

  test('should create new journal', async ({ page }) => {
    const journalListPage = new JournalListPage(page);
    await journalListPage.goto();
    await journalListPage.clickCreateJournal();

    const journalFormPage = new JournalFormPage(page);
    await journalFormPage.fillJournalForm({
      journalNo: `JNL-${Date.now()}`,
      journalDate: '2024-04-25',
      description: 'Test Journal',
      status: 'DRAFT',
    });

    await journalFormPage.saveJournal();

    await expect(page).toHaveURL('**/accounting/journals');
    await expect(page.getByText('Journal created successfully')).toBeVisible();
  });

  test('should validate required fields', async ({ page }) => {
    const journalFormPage = new JournalFormPage(page);
    await journalFormPage.goto();

    await journalFormPage.saveJournal();

    const errors = await journalFormPage.getValidationErrors();
    expect(errors).toContain('Journal No is required');
    expect(errors).toContain('Description is required');
  });

  test('should search journals', async ({ page }) => {
    const journalListPage = new JournalListPage(page);
    await journalListPage.goto();

    await journalListPage.searchJournals('JNL-001');
    await journalListPage.waitForJournals();

    const count = await journalListPage.getJournalCount();
    expect(count).toBe(1);
  });

  test('should filter journals by status', async ({ page }) => {
    const journalListPage = new JournalListPage(page);
    await journalListPage.goto();

    await journalListPage.filterByStatus('POSTED');
    await journalListPage.waitForJournals();

    const count = await journalListPage.getJournalCount();
    expect(count).toBeGreaterThan(0);
  });
});
```

### Customer Management Flow

```typescript
// e2e/crm/customers.spec.ts
import { test, expect } from '@playwright/test';
import { CustomerListPage } from '../pages/customer-list.page';
import { CustomerFormPage } from '../pages/customer-form.page';

test.describe('Customer Management', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login('user@example.com', 'password123');
  });

  test('should list customers', async ({ page }) => {
    const customerListPage = new CustomerListPage(page);
    await customerListPage.goto();

    await customerListPage.waitForCustomers();

    const count = await customerListPage.getCustomerCount();
    expect(count).toBeGreaterThan(0);
  });

  test('should create new customer', async ({ page }) => {
    const customerListPage = new CustomerListPage(page);
    await customerListPage.goto();
    await customerListPage.clickCreateCustomer();

    const customerFormPage = new CustomerFormPage(page);
    await customerFormPage.fillCustomerForm({
      code: `CUST-${Date.now()}`,
      name: 'Test Customer',
      email: 'test@example.com',
      phone: '+84123456789',
      address: '123 Test Street',
    });

    await customerFormPage.saveCustomer();

    await expect(page).toHaveURL('**/crm/customers');
    await expect(page.getByText('Customer created successfully')).toBeVisible();
  });
});
```

## 📦 Running E2E Tests

### Commands

```bash
# Run all E2E tests
npx playwright test

# Run tests in headed mode (visible browser)
npx playwright test --headed

# Run tests in debug mode
npx playwright test --debug

# Run tests on specific browser
npx playwright test --project=chromium

# Run tests in UI mode
npx playwright test --ui

# Run specific test file
npx playwright test e2e/accounting/journals.spec.ts

# Run tests with specific grep pattern
npx playwright test --grep "should create"

# Run tests and show trace on failure
npx playwright test --trace on

# Generate HTML report
npx playwright test --reporter=html
```

### CI/CD Configuration

```yaml
# .github/workflows/e2e.yml
name: E2E Tests

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  e2e:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Install dependencies
        run: npm ci

      - name: Install Playwright Browsers
        run: npx playwright install --with-deps

      - name: Build application
        run: npm run build

      - name: Run E2E tests
        run: npx playwright test

      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: playwright-report/
```

## 📦 Test Data Management

### Test Data Setup

```typescript
// e2e/fixtures/test-data.ts
export const testUsers = {
  admin: {
    email: 'admin@example.com',
    password: 'Admin123!',
    role: 'ADMIN',
  },
  user: {
    email: 'user@example.com',
    password: 'User123!',
    role: 'USER',
  },
};

export const testCustomers = {
  valid: {
    code: 'CUST-TEST-001',
    name: 'Test Customer',
    email: 'test@example.com',
    phone: '+84123456789',
    address: '123 Test Street',
    status: 'ACTIVE',
  },
  invalid: {
    code: '',
    name: '',
    email: 'invalid-email',
    phone: 'invalid-phone',
    address: '',
  },
};

export const testJournals = {
  valid: {
    journalNo: `JNL-${Date.now()}`,
    journalDate: new Date().toISOString().split('T')[0],
    description: 'Test Journal Description',
    status: 'DRAFT',
  },
};
```

### Database Seed Script

```typescript
// e2e/utils/seed-db.ts
import { execSync } from 'child_process';

export async function seedTestData() {
  console.log('Seeding test data...');

  execSync('npm run seed:test-db', {
    stdio: 'inherit',
  });

  console.log('Test data seeded successfully');
}

export async function clearTestData() {
  console.log('Clearing test data...');

  execSync('npm run clear:test-db', {
    stdio: 'inherit',
  });

  console.log('Test data cleared successfully');
}
```

## 🎯 Usage Examples

```
/e2e-testing "Tạo E2E test scenario"

Usage:
/e2e-testing "Tạo E2E test cho customer management flow: list, create, edit, delete"

Sẽ:
1. Tạo page objects
2. Tạo test scenarios
3. Implement assertions
```

```
/e2e-testing "Setup Playwright"

Usage:
/e2e-testing "Setup Playwright configuration với Chrome, Firefox, Safari browsers"

Sẽ:
1. Configure playwright.config.ts
2. Setup projects cho browsers
3. Configure reporters
```

```
/e2e-testing "Run E2E tests"

Usage:
/e2e-testing "Run E2E tests cho accounting module"

Sẽ:
1. Chạy tests
2. Generate HTML report
3. Show results
```

## 📦 Best Practices

### Test Organization

- Use Page Object Model
- Keep tests independent
- Use descriptive test names
- Group related tests
- Reuse page components

### Test Quality

- Test user workflows
- Test error scenarios
- Test edge cases
- Use appropriate waits
- Handle test flakiness

### Performance

- Run tests in parallel
- Use efficient selectors
- Avoid unnecessary waits
- Use test fixtures
- Optimize test execution time

---

*Last Updated: 2026-04-25*
