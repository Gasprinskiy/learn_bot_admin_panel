import { AccessRight } from '@/shared/types/profile';

export const StaffListErrorsMap: Record<number, string> = {
  404: 'Работники не найдены',
  500: 'Внутреняя ошибка сервера, попробуйте позже',
};

export const AccesRightTitleMap: Record<AccessRight, string> = {
  [AccessRight.AccessRightFull]: 'Полный',
  [AccessRight.AccessRightManager]: 'Менеджер',
  [AccessRight.AccessRightTeacher]: 'Учитель',
};
