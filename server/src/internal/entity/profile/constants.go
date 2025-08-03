package profile

import "fmt"

type AccessRight string

const (
	AccessRightFull    AccessRight = "full_access"
	AccessRightManager AccessRight = "manager_access"
	AccessRightTeacher AccessRight = "teacher_access"
)

var AccessRightCheckMap = map[AccessRight]struct{}{
	AccessRightFull:    {},
	AccessRightManager: {},
	AccessRightTeacher: {},
}

func (s *AccessRight) Scan(value any) error {
	if b, ok := value.([]byte); ok {
		*s = AccessRight(string(b))
		return nil
	}
	if sVal, ok := value.(string); ok {
		*s = AccessRight(sVal)
		return nil
	}

	if s.HasInAccessRightMap() {
		return fmt.Errorf("право доступа %v не найдено в списке доступных", value)
	}

	return fmt.Errorf("не удалось записать %v в тип AccessRight", value)
}

func (s AccessRight) HasInAccessRightMap() bool {
	_, exists := AccessRightCheckMap[s]
	return exists
}

const (
	AuthSuccessfulyMessage = "Здравствуйте %s, Авторизация прошла успешно, можете вернутся обратно!"
)
