export const ErrorMessagesByCode: Record<number, string> = {
  401: 'Требуется авторизация',
  404: 'Данные пользователя не найдены',
  410: 'Сессия истекла',
  500: 'Внутреняя ошибка сервера, попробуйте позже',
};

export const PasswordLoginErrorMessagesByCode: Record<number, string> = {
  404: 'Не верный логин или пароль',
  500: 'Внутреняя ошибка сервера, попробуйте позже',
  400: 'Не верный логин или пароль',
};
