import { Injectable, signal, effect } from '@angular/core';
import { usePreset } from '@primeuix/themes';
import {
  PALETTES,
  RADIUS,
  SCALE,
  FONTS,
  BASE_FONT_SIZE,
  ACCENT_COLORS,
  LINE_HEIGHTS,
  LETTER_SPACINGS,
  createArdaPreset,
} from '@arda-mfe/shared-theme';
import {
  ThemeSettings,
  DEFAULT_THEME_SETTINGS,
} from './theme-settings';

@Injectable({ providedIn: 'root' })
export class ThemeService {
  private static STORAGE_KEY = 'arda-theme-settings';

  settings = signal<ThemeSettings>(this.loadSettings());
  private fontLink: HTMLLinkElement | null = null;

  constructor() {
    effect(() => {
      this.applySettings(this.settings());
    });
  }

  updateSetting<K extends keyof ThemeSettings>(
    key: K,
    value: ThemeSettings[K],
  ): void {
    const updated = { ...this.settings(), [key]: value };
    this.settings.set(updated);
    this.saveSettings(updated);
  }

  private applySettings(s: ThemeSettings): void {
    // Apply preset (palette, radius, scale)
    const preset = createArdaPreset({
      palette: PALETTES[s.palette],
      radius: s.radius as keyof typeof RADIUS,
      scale: s.scale as keyof typeof SCALE,
    });
    usePreset(preset);

    // Apply font
    const font = FONTS[s.font];
    this.applyFont(font);
    document.documentElement.style.fontSize =
      BASE_FONT_SIZE[s.baseFontSize as keyof typeof BASE_FONT_SIZE];

    // Apply dark mode
    if (s.darkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }

    // Apply accent color
    const accentColor = ACCENT_COLORS[s.accentColor as keyof typeof ACCENT_COLORS];
    document.documentElement.style.setProperty('--accent-color', accentColor);

    // Apply line height
    document.documentElement.style.lineHeight = LINE_HEIGHTS[s.lineHeight as keyof typeof LINE_HEIGHTS];

    // Apply letter spacing
    document.documentElement.style.letterSpacing = LETTER_SPACINGS[s.letterSpacing as keyof typeof LETTER_SPACINGS];
  }

  private applyFont(font: { importUrl: string; family: string }): void {
    if (this.fontLink) {
      this.fontLink.remove();
      this.fontLink = null;
    }
    if (font.importUrl) {
      const link = document.createElement('link');
      link.rel = 'stylesheet';
      link.href = font.importUrl;
      document.head.appendChild(link);
      this.fontLink = link;
    }
    document.documentElement.style.fontFamily = font.family;
  }

  private loadSettings(): ThemeSettings {
    try {
      const stored = localStorage.getItem(ThemeService.STORAGE_KEY);
      if (stored) {
        return { ...DEFAULT_THEME_SETTINGS, ...JSON.parse(stored) };
      }
    } catch {
      /* ignore */
    }
    return { ...DEFAULT_THEME_SETTINGS };
  }

  private saveSettings(settings: ThemeSettings): void {
    localStorage.setItem(ThemeService.STORAGE_KEY, JSON.stringify(settings));
  }
}
