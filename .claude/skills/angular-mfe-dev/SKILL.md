---
name: angular-mfe-dev
description: Hỗ trợ phát triển Angular 21+, Nx Monorepo, PrimeNG và Module Federation cho dự án Arda
disable-model-invocation: false
---

# Angular MFE Development Skill

Mục đích: Hỗ trợ phát triển Frontend Applications trong thư mục `arda-mfe/`.

## 🎯 Phạm vi

- **Framework**: Angular 21+ (Signals, Standalone Components)
- **Monorepo**: Nx 22+ (Libraries, Affected builds)
- **Architecture**: Micro-Frontend (Module Federation)
- **UI Library**: PrimeNG 21+ (Tailwind CSS)
- **Reactivity**: RxJS & Angular Signals

## 📦 Project Structure (Arda MFE)

```
arda-mfe/
├── apps/
│   ├── shell/             # Host application
│   └── [feature]-mfe/     # Remote micro-frontends
├── libs/
│   ├── ui/                # Shared PrimeNG components
│   ├── auth/              # Zitadel OIDC integration
│   └── shared/            # Utilities & Base models
└── nx.json
```

## 🛠️ Key Patterns

### 1. Signals & Performance
Ưu tiên sử dụng Signals cho state local và UI reactivity.
```typescript
readonly items = signal<Item[]>([]);
readonly total = computed(() => this.items().length);
```

### 2. Module Federation
Cấu hình `webpack.config.js` để chia sẻ dependencies và lazy load remotes.

### 3. PrimeNG & Tailwind
Sử dụng PrimeNG components kết hợp với Tailwind CSS cho layout.

## 🎯 Usage
- `/angular-mfe-dev "Tạo MFE mới cho domain [name]"`
- `/angular-mfe-dev "Tạo shared component sử dụng PrimeNG DataTable"`
- `/angular-mfe-dev "Cấu hình Module Federation cho remote app"`
