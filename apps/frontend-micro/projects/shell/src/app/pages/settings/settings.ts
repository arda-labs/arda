import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ThemeService, RADIUS, SCALE, FONTS, BASE_FONT_SIZE, ACCENT_COLORS, LINE_HEIGHTS, LETTER_SPACINGS } from '@arda/core';
import { SelectButton } from 'primeng/selectbutton';
import { ToggleButton } from 'primeng/togglebutton';
import { Tabs, TabList, Tab, TabPanels, TabPanel } from 'primeng/tabs';

interface ColorOption {
  label: string;
  value: string;
}

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

function toOptions(obj: Record<string, any>): { label: string; value: string }[] {
  return Object.keys(obj).map((k) => ({ label: capitalize(k), value: k }));
}

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    SelectButton,
    ToggleButton,
    Tabs,
    TabList,
    Tab,
    TabPanels,
    TabPanel,
  ],
  templateUrl: './settings.html',
})
export class Settings {
  private themeService = inject(ThemeService);

  activeTab: 'appearance' | 'typography' | 'layout' = 'appearance';

  colorOptions: ColorOption[] = [
    { label: 'Blue', value: 'blue' },
    { label: 'Indigo', value: 'indigo' },
    { label: 'Teal', value: 'teal' },
    { label: 'Rose', value: 'rose' },
    { label: 'Emerald', value: 'emerald' },
    { label: 'Purple', value: 'purple' },
    { label: 'Orange', value: 'orange' },
    { label: 'Cyan', value: 'cyan' },
    { label: 'Pink', value: 'pink' },
    { label: 'Amber', value: 'amber' },
  ];

  colorHexMap: Record<string, string> = {
    blue: '#3b82f6',
    indigo: '#6366f1',
    teal: '#14b8a6',
    rose: '#f43f5e',
    emerald: '#10b981',
    purple: '#a855f7',
    orange: '#f97316',
    cyan: '#06b6d4',
    pink: '#ec4899',
    amber: '#f59e0b',
  };

  radiusOptions = toOptions(RADIUS);
  scaleOptions = toOptions(SCALE);
  fontOptions = toOptions(FONTS);
  fontSizeOptions = toOptions(BASE_FONT_SIZE);
  accentColorOptions = toOptions(ACCENT_COLORS);
  lineHeightOptions = toOptions(LINE_HEIGHTS);
  letterSpacingOptions = toOptions(LETTER_SPACINGS);

  settings = this.themeService.settings;

  getColorPreview(colorKey: string): string {
    return this.colorHexMap[colorKey] || '#3b82f6';
  }

  getFontFamily(fontKey: string): string {
    const font = FONTS[fontKey as keyof typeof FONTS];
    return font?.family || 'sans-serif';
  }

  getFontSizeValue(fontSizeKey: string): string {
    const size = BASE_FONT_SIZE[fontSizeKey as keyof typeof BASE_FONT_SIZE];
    return size || '16px';
  }

  getAccentColorValue(accentKey: string): string {
    const accent = ACCENT_COLORS[accentKey as keyof typeof ACCENT_COLORS];
    return accent || '{primary.500}';
  }

  onPaletteChange(value: string): void {
    this.themeService.updateSetting('palette', value);
  }

  onRadiusChange(value: string): void {
    this.themeService.updateSetting('radius', value);
  }

  onScaleChange(value: string): void {
    this.themeService.updateSetting('scale', value);
  }

  onFontChange(value: string): void {
    this.themeService.updateSetting('font', value);
  }

  onFontSizeChange(value: string): void {
    this.themeService.updateSetting('baseFontSize', value);
  }

  onDarkModeChange(value: boolean): void {
    this.themeService.updateSetting('darkMode', value);
  }

  onAccentColorChange(value: string): void {
    this.themeService.updateSetting('accentColor', value);
  }

  onLineHeightChange(value: string): void {
    this.themeService.updateSetting('lineHeight', value);
  }

  onLetterSpacingChange(value: string): void {
    this.themeService.updateSetting('letterSpacing', value);
  }
}
