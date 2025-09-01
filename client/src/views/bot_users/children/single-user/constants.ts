import { PaymentTypeID, SubscriptionCancelReason } from '@/shared/types/bot_users';
import type { DropdownMixedOption } from 'naive-ui/es/dropdown/src/interface';
import { DropdownKey } from './types';

export const UploadErrorMessagesByCode: Record<number, string> = {
  413: 'Файл размером больше 5 мб',
  400: 'Не валидные параметры',
  500: 'Внутреняя ошибка сервера, попробуйте позже',
};

export const CancelErrorMessagesByCode: Record<number, string> = {
  500: 'Внутреняя ошибка сервера, попробуйте позже',
};

export const PaymentTypeNameByID: Record<PaymentTypeID, string> = {
  [PaymentTypeID.P2P]: 'Перевод на карту',
  [PaymentTypeID.Payme]: 'Payme',
};

export const CacnelReasonLabelMap: Record<SubscriptionCancelReason, string> = {
  [SubscriptionCancelReason.EXPIRED]: 'Подписка истекла',
  [SubscriptionCancelReason.MONEYBACK]: 'Возврат средств',
};

export const DropdownOptions: Array<DropdownMixedOption> = [
  {
    label: 'Подключить подписку',
    key: DropdownKey.CONNECTSUBS,
  },
  {
    label: 'Отменить подписку',
    key: DropdownKey.CACNELSUBS,
  },
];
