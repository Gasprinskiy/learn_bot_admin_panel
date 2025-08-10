import { SubscriptionStatus } from '@/shared/types/profile';
import type { SelectMixedOption } from 'naive-ui/es/select/src/interface';

type UIStatusType = 'default' | 'primary' | 'success' | 'info' | 'warning' | 'error';

export const SubscriptionStatusTitleMap: Record<SubscriptionStatus, string> = {
  [SubscriptionStatus.ACTIVE]: 'активна',
  [SubscriptionStatus.EXPIRED]: 'исеткла',
  [SubscriptionStatus.NOTEXISTS]: 'отсуствует',
};

export const SubscriptionUIStatus: Record<SubscriptionStatus, UIStatusType> = {
  [SubscriptionStatus.ACTIVE]: 'success',
  [SubscriptionStatus.EXPIRED]: 'warning',
  [SubscriptionStatus.NOTEXISTS]: 'default',
};

export const SubscriptionStatusSelectOptions: SelectMixedOption[] = [
  {
    label: SubscriptionStatusTitleMap[SubscriptionStatus.ACTIVE],
    value: SubscriptionStatus.ACTIVE,
  },
  {
    label: SubscriptionStatusTitleMap[SubscriptionStatus.EXPIRED],
    value: SubscriptionStatus.EXPIRED,
  },
  {
    label: SubscriptionStatusTitleMap[SubscriptionStatus.NOTEXISTS],
    value: SubscriptionStatus.NOTEXISTS,
  },
];
