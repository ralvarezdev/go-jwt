package interception

type (
	// Interception is the type of interception
	Interception int
)

// Interception values
const (
	RefreshToken Interception = iota
	AccessToken
	None
)
