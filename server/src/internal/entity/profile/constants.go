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

var AccessRightIDMap = map[AccessRight]int{
	AccessRightFull:    1,
	AccessRightManager: 2,
	AccessRightTeacher: 3,
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

func (s AccessRight) HasAccessRightInList(arList []AccessRight) bool {
	arMap := make(map[AccessRight]struct{}, len(arList))
	for _, ar := range arList {
		arMap[ar] = struct{}{}
	}

	_, exists := arMap[s]
	return exists
}
