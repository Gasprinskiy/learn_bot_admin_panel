export interface UseRedirectWindowParams {
  name: string;
  width?: number;
  height?: number;
}

export interface ScreenWithAvailCoords extends Screen {
  availLeft?: number;
  availTop?: number;
}

export interface UseRedirectWindowReturnType {
  open: (url: string) => void;
  close: () => void;
}
