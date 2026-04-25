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
    50: '{blue.50}',
    100: '{blue.100}',
    200: '{blue.200}',
    300: '{blue.300}',
    400: '{blue.400}',
    500: '{blue.500}',
    600: '{blue.600}',
    700: '{blue.700}',
    800: '{blue.800}',
    900: '{blue.900}',
    950: '{blue.950}',
  },
  indigo: {
    50: '{indigo.50}',
    100: '{indigo.100}',
    200: '{indigo.200}',
    300: '{indigo.300}',
    400: '{indigo.400}',
    500: '{indigo.500}',
    600: '{indigo.600}',
    700: '{indigo.700}',
    800: '{indigo.800}',
    900: '{indigo.900}',
    950: '{indigo.950}',
  },
  teal: {
    50: '{teal.50}',
    100: '{teal.100}',
    200: '{teal.200}',
    300: '{teal.300}',
    400: '{teal.400}',
    500: '{teal.500}',
    600: '{teal.600}',
    700: '{teal.700}',
    800: '{teal.800}',
    900: '{teal.900}',
    950: '{teal.950}',
  },
  rose: {
    50: '{rose.50}',
    100: '{rose.100}',
    200: '{rose.200}',
    300: '{rose.300}',
    400: '{rose.400}',
    500: '{rose.500}',
    600: '{rose.600}',
    700: '{rose.700}',
    800: '{rose.800}',
    900: '{rose.900}',
    950: '{rose.950}',
  },
  emerald: {
    50: '{emerald.50}',
    100: '{emerald.100}',
    200: '{emerald.200}',
    300: '{emerald.300}',
    400: '{emerald.400}',
    500: '{emerald.500}',
    600: '{emerald.600}',
    700: '{emerald.700}',
    800: '{emerald.800}',
    900: '{emerald.900}',
    950: '{emerald.950}',
  },
  purple: {
    50: '{purple.50}',
    100: '{purple.100}',
    200: '{purple.200}',
    300: '{purple.300}',
    400: '{purple.400}',
    500: '{purple.500}',
    600: '{purple.600}',
    700: '{purple.700}',
    800: '{purple.800}',
    900: '{purple.900}',
    950: '{purple.950}',
  },
  orange: {
    50: '{orange.50}',
    100: '{orange.100}',
    200: '{orange.200}',
    300: '{orange.300}',
    400: '{orange.400}',
    500: '{orange.500}',
    600: '{orange.600}',
    700: '{orange.700}',
    800: '{orange.800}',
    900: '{orange.900}',
    950: '{orange.950}',
  },
  cyan: {
    50: '{cyan.50}',
    100: '{cyan.100}',
    200: '{cyan.200}',
    300: '{cyan.300}',
    400: '{cyan.400}',
    500: '{cyan.500}',
    600: '{cyan.600}',
    700: '{cyan.700}',
    800: '{cyan.800}',
    900: '{cyan.900}',
    950: '{cyan.950}',
  },
  pink: {
    50: '{pink.50}',
    100: '{pink.100}',
    200: '{pink.200}',
    300: '{pink.300}',
    400: '{pink.400}',
    500: '{pink.500}',
    600: '{pink.600}',
    700: '{pink.700}',
    800: '{pink.800}',
    900: '{pink.900}',
    950: '{pink.950}',
  },
  amber: {
    50: '{amber.50}',
    100: '{amber.100}',
    200: '{amber.200}',
    300: '{amber.300}',
    400: '{amber.400}',
    500: '{amber.500}',
    600: '{amber.600}',
    700: '{amber.700}',
    800: '{amber.800}',
    900: '{amber.900}',
    950: '{amber.950}',
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
    xs: '{spacing.2}',
    sm: '{spacing.3}',
    md: '{spacing.4}',
    lg: '{spacing.5}',
    xl: '{spacing.6}',
  },
  default: {
    xs: '{spacing.4}',
    sm: '{spacing.5}',
    md: '{spacing.6}',
    lg: '{spacing.8}',
    xl: '{spacing.10}',
  },
  large: {
    xs: '{spacing.6}',
    sm: '{spacing.8}',
    md: '{spacing.10}',
    lg: '{spacing.12}',
    xl: '{spacing.14}',
  },
};

export const ACCENT_COLORS = {
  primary: '{primary.500}',
  secondary: '{secondary.500}',
  success: '{success.500}',
  info: '{info.500}',
  warning: '{warning.500}',
  danger: '{danger.500}',
  surface: '{surface.500}',
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
