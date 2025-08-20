import { PaymentTypeID } from '@/shared/types/bot_users';

export const UploadErrorMessagesByCode: Record<number, string> = {
  413: 'Файл размером больше 5 мб',
  400: 'Не валидные параметры',
};

export const PaymentTypeNameByID: Record<PaymentTypeID, string> = {
  [PaymentTypeID.P2P]: 'Перевод на карту',
};
