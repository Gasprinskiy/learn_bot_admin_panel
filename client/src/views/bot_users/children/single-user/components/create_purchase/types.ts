import type { BotSubscriptionType } from '@/shared/types/bot_users';

export interface CreatePurchaseProps {
  subscriptionTypes: BotSubscriptionType[];
}

export type CreatePurchaseEmits = {
  onSubmit: [sub_id: number, upload: File];
};
