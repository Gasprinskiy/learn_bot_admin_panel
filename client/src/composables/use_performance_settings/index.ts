import { VisualMode } from './types';
import { VisualModeOptions } from './constants';
import { useStorage } from '@vueuse/core';

const state = useStorage('visual_mode', VisualMode.PERFORMANCE);

export function usePerformanceSettings() {
  function addClassToBody() {
    const value = state.value;

    const body = document.body;

    if (body.classList.contains(value)) {
      return;
    }
    body.classList.remove(VisualMode.DEFAULT, VisualMode.PERFORMANCE);
    body.classList.add(value);
  }

  function setMode(value: VisualMode) {
    state.value = value;
    addClassToBody();
  }

  return {
    VisualModeOptions,
    currentMode: state,
    addClassToBody,
    setMode,
  };
}
