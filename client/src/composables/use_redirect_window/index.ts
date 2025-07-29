import { shallowRef } from 'vue';
import type { UseRedirectWindowParams } from './types';

export function useRedirectWindow(params: UseRedirectWindowParams) {
  const { name, width, height } = params;

  const windowWidth = width || 500;
  const windowHeight = height || 600;

  const windowProxy = shallowRef<WindowProxy | null>();
  const closed = shallowRef<boolean>(true);

  function listenClosed() {
    const interval = setInterval(() => {
      closed.value = windowProxy.value?.closed || false;
      if (closed.value) {
        clearInterval(interval);
      }
    }, 500);
  }

  function open(url: string) {
    const screenWidth = window.screen.availWidth;
    const screenHeight = window.screen.availHeight;

    const centerX = window.screen.availLeft! + (screenWidth - windowWidth) / 2;
    const centerY = window.screen.availTop! + (screenHeight - windowHeight) / 2;

    windowProxy.value = window.open(
      url,
      name,
      `top=${centerY},left=${centerX},width=${windowWidth},height=${windowHeight},toolbar=no,menubar=no,resizable=false`,
    );

    listenClosed();
  }

  function close() {
    closed.value = true;
    windowProxy.value?.close();
  }

  return {
    open,
    close,
  };
}
