import { definePreset } from '@primeuix/themes';
import Aura from '@primeuix/themes/aura';
import { PALETTES, RADIUS, SCALE, ColorScale } from './palettes';

export interface ArdaPresetOptions {
  palette?: ColorScale;
  radius?: keyof typeof RADIUS;
  scale?: keyof typeof SCALE;
}

export const ArdaPreset = definePreset(Aura, {
  semantic: {
    primary: PALETTES['blue'],
  },
});

export function createArdaPreset(options: ArdaPresetOptions = {}) {
  const { palette = PALETTES['blue'], radius, scale } = options;
  return definePreset(Aura, {
    semantic: {
      primary: palette,
      ...(radius ? { borderRadius: RADIUS[radius] } : {}),
      ...(scale ? { formField: SCALE[scale] } : {}),
    },
  });
}

export { Aura };
