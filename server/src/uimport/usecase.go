package uimport

import "learn_bot_admin_panel/internal/usecase"

type Usecase struct {
	Profile *usecase.Profile
	Jwt     *usecase.Jwt
}
