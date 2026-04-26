export interface ThemeSettings {
  palette: string;
  radius: string;
  scale: string;
  font: string;
  baseFontSize: string;
  darkMode: boolean;
  accentColor: string;
  lineHeight: string;
  letterSpacing: string;
}

export const DEFAULT_THEME_SETTINGS: ThemeSettings = {
  palette: 'blue',
  radius: 'rounded',
  scale: 'default',
  font: 'inter',
  baseFontSize: 'default',
  darkMode: false,
  accentColor: 'primary',
  lineHeight: 'normal',
  letterSpacing: 'normal',
};
