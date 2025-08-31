import type { User } from '@/shared/types/profile';

export type ProfileFormState = Pick<User, 'first_name' | 'last_name' | 'tg_user_name'>

export interface CreateRedactProfileProps {
  form_data?: ProfileFormState;
};

export type CreateRedactProfileEmits = {
  onSubmit: [state: ProfileFormState];
};
