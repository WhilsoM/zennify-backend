package ports

type RegisterRequest struct {
	Username string `validate:"required,min=3,max=64"`
	Password string `validate:"required,min=8,max=128"`
}

type LoginRequest struct {
	Username string `validate:"required,min=1,max=64"`
	Password string `validate:"required,min=1,max=128"`
}

type RefreshTokensRequest struct {
	RefreshToken string `validate:"required,min=1,max=550"`
}
