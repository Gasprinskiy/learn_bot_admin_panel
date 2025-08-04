export type SetPasswordEmits = {
  onSubmit: [password: string];
};

export interface SetPasswordState {
  passwordSource: string;
  passwordСonfirm: string;
}
