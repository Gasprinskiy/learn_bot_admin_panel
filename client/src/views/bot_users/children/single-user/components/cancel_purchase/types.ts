import type { SubscriptionCancelReason } from '@/shared/types/bot_users';

export type CancelPurchaseEmits = {
  onSubmit: [reason: SubscriptionCancelReason];
};
