export interface ColorScale {
  50: string;
  100: string;
  200: string;
  300: string;
  400: string;
  500: string;
  600: string;
  700: string;
  800: string;
  900: string;
  950: string;
}

export const PALETTES: Record<string, ColorScale> = {
  blue: {
    50: '#eff6ff',
    100: '#dbeade',
    200: '#bfdbfe',
    300: '#93c5fd',
    400: '#60a5fa',
    500: '#3b82f6',
    600: '#2563eb',
    700: '#1d4ed8',
    800: '#1e40af',
    900: '#1e3a8a',
    950: '#172554',
  },
  indigo: {
    50: '#eef2ff',
    100: '#e0e7ff',
    200: '#c7d2fe',
    300: '#a5b4fc',
    400: '#818cf8',
    500: '#6366f1',
    600: '#4f46e5',
    700: '#4338ca',
    800: '#3730a3',
    900: '#312e81',
    950: '#1e1b4b',
  },
  teal: {
    50: '#f0fdfa',
    100: '#ccfbf1',
    200: '#99f6e4',
    300: '#5eead4',
    400: '#2dd4bf',
    500: '#14b8a6',
    600: '#0d9488',
    700: '#0f766e',
    800: '#115e59',
    900: '#134e4a',
    950: '#042f2e',
  },
  rose: {
    50: '#fff1f2',
    100: '#ffe4e6',
    200: '#fecdd3',
    300: '#fda4af',
    400: '#fb7185',
    500: '#f43f5e',
    600: '#e11d48',
    700: '#be123c',
    800: '#9f1239',
    900: '#881337',
    950: '#4c0519',
  },
  emerald: {
    50: '#ecfdf5',
    100: '#d1fae5',
    200: '#a7f3d0',
    300: '#6ee7b7',
    400: '#34d399',
    500: '#10b981',
    600: '#059669',
    700: '#047857',
    800: '#065f46',
    900: '#064e3b',
    950: '#022c22',
  },
  purple: {
    50: '#faf5ff',
    100: '#f3e8ff',
    200: '#e9d5ff',
    300: '#d8b4fe',
    400: '#c084fc',
    500: '#a855f7',
    600: '#9333ea',
    700: '#7e22ce',
    800: '#6b21a8',
    900: '#581c87',
    950: '#3b0764',
  },
  orange: {
    50: '#fff7ed',
    100: '#ffedd5',
    200: '#fed7aa',
    300: '#fdba74',
    400: '#fb923c',
    500: '#f97316',
    600: '#ea580c',
    700: '#c2410c',
    800: '#9a3412',
    900: '#7c2d12',
    950: '#431407',
  },
  cyan: {
    50: '#ecfeff',
    100: '#cffafe',
    200: '#a5f3fc',
    300: '#67e8f9',
    400: '#22d3ee',
    500: '#06b6d4',
    600: '#0891b2',
    700: '#0e7490',
    800: '#155e75',
    900: '#164e63',
    950: '#083344',
  },
  pink: {
    50: '#fdf2f8',
    100: '#fce7f3',
    200: '#fbcfe8',
    300: '#f9a8d4',
    400: '#f472b6',
    500: '#ec4899',
    600: '#db2777',
    700: '#be185d',
    800: '#9d174d',
    900: '#831843',
    950: '#500724',
  },
  amber: {
    50: '#fffbeb',
    100: '#fef3c7',
    200: '#fde68a',
    300: '#fcd34d',
    400: '#fbbf24',
    500: '#f59e0b',
    600: '#d97706',
    700: '#b45309',
    800: '#92400e',
    900: '#78350f',
    950: '#451a03',
  },
};

export const RADIUS = {
  sharp: { none: '0', xs: '0', sm: '0', md: '0', lg: '0', xl: '0' },
  rounded: {
    none: '0',
    xs: '0.125rem',
    sm: '0.25rem',
    md: '0.375rem',
    lg: '0.5rem',
    xl: '0.75rem',
  },
  pill: {
    none: '0',
    xs: '0.5rem',
    sm: '1rem',
    md: '1.5rem',
    lg: '2rem',
    xl: '3rem',
  },
};

export const SCALE = {
  compact: {
    xs: '0.5rem',
    sm: '0.75rem',
    md: '1rem',
    lg: '1.25rem',
    xl: '1.5rem',
  },
  default: {
    xs: '1rem',
    sm: '1.25rem',
    md: '1.5rem',
    lg: '2rem',
    xl: '2.5rem',
  },
  large: {
    xs: '1.5rem',
    sm: '2rem',
    md: '2.5rem',
    lg: '3rem',
    xl: '3.5rem',
  },
};

export const ACCENT_COLORS = {
  primary: 'var(--p-primary-500)',
  secondary: 'var(--p-secondary-500)',
  success: 'var(--p-success-500)',
  info: 'var(--p-info-500)',
  warning: 'var(--p-warning-500)',
  danger: 'var(--p-danger-500)',
  surface: 'var(--p-surface-500)',
};

export const FONTS: Record<
  string,
  { importUrl: string; family: string }
> = {
  inter: {
    importUrl:
      'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap',
    family: 'Inter, sans-serif',
  },
  roboto: {
    importUrl:
      'https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap',
    family: 'Roboto, sans-serif',
  },
  plus: {
    importUrl:
      'https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap',
    family: '"Plus Jakarta Sans", sans-serif',
  },
  poppins: {
    importUrl:
      'https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700&display=swap',
    family: 'Poppins, sans-serif',
  },
  opensans: {
    importUrl:
      'https://fonts.googleapis.com/css2?family=Open+Sans:wght@400;500;600;700&display=swap',
    family: '"Open Sans", sans-serif',
  },
  system: {
    importUrl: '',
    family: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  },
  lato: {
    importUrl:
      'https://fonts.googleapis.com/css2?family=Lato:wght@400;500;600;700&display=swap',
    family: 'Lato, sans-serif',
  },
  montserrat: {
    importUrl:
      'https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500;600;700&display=swap',
    family: 'Montserrat, sans-serif',
  },
};

export const LINE_HEIGHTS = {
  tight: '1.25',
  snug: '1.375',
  normal: '1.5',
  relaxed: '1.625',
  loose: '2',
};

export const LETTER_SPACINGS = {
  tighter: '-0.05em',
  tight: '-0.025em',
  normal: '0',
  wide: '0.025em',
  wider: '0.05em',
  widest: '0.1em',
};

export const BASE_FONT_SIZE = {
  xs: '12px',
  small: '14px',
  default: '16px',
  large: '18px',
  xl: '20px',
  xxl: '22px',
};
