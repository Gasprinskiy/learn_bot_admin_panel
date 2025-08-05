import { VisualMode } from './types';
import type { VisualModeOption } from './types';

export const VisualModeOptions: Array<VisualModeOption> = [
  {
    value: VisualMode.PERFORMANCE,
    label: 'Режим быстродействия',
  },
  {
    value: VisualMode.DEFAULT,
    label: 'Обычный режим',
  },
];
