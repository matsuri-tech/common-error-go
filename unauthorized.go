package merrors

const (
	ErrorInvalidAuthority ErrorType = "invalid_authority"
)

func InvalidAuthority() error {
	return ErrorUnauthorized("invalid authority", ErrorInvalidAuthority)
}
