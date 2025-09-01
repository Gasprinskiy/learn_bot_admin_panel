import type { SelectMixedOption } from 'naive-ui/es/select/src/interface';

import { SubscriptionCancelReason } from '@/shared/types/bot_users';
import { CacnelReasonLabelMap } from '../../constants';

export const CacnelReasonOptions: Array<SelectMixedOption> = [
  {
    label: CacnelReasonLabelMap[SubscriptionCancelReason.EXPIRED],
    value: SubscriptionCancelReason.EXPIRED,
  },
  {
    label: CacnelReasonLabelMap[SubscriptionCancelReason.MONEYBACK],
    value: SubscriptionCancelReason.MONEYBACK,
  },
];
