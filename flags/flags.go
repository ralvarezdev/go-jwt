package flags

import (
	goflags "github.com/ralvarezdev/go-flags"
)

type (
	// PrivateKeyFlag is the name of the private key flag.
	PrivateKeyFlag struct {
		goflags.Flag
	}
	
	// PublicKeyFlag is the name of the public key flag.
	PublicKeyFlag struct {
		goflags.Flag
	}
)

// NewPrivateKeyFlag creates a new PrivateKeyFlag.
// 
// Parameters:
// 
// 	defaultValue - the default value for the flag.
// 
// Returns:
// 
// 	- A pointer to the created PrivateKeyFlag.
func NewPrivateKeyFlag(
	defaultValue *string,
) *PrivateKeyFlag {
	return &PrivateKeyFlag{
		Flag: *goflags.NewFlag(defaultValue, nil, PrivateKeyFlagName, PrivateKeyFlagUsage),
	}
}

// NewPublicKeyFlag creates a new PublicKeyFlag.
// 
// Parameters:
// 
// 	defaultValue - the default value for the flag.
// 
// Returns:
// 
// 	- A pointer to the created PublicKeyFlag.
func NewPublicKeyFlag(
	defaultValue *string,
) *PublicKeyFlag {
	return &PublicKeyFlag{
		Flag: *goflags.NewFlag(defaultValue, nil, PublicKeyFlagName, PublicKeyFlagUsage),
	}
}

// Default returns the default value of the flag.
// 
// Returns:
// 
//	The default value.
func (f *PrivateKeyFlag) Default() string {
	if f == nil {
		return ""
	}
	return f.Default()
}

// Default returns the default value of the flag.
// 
// Returns:
// 
//	The default value.
func (f *PublicKeyFlag) Default() string {
	if f == nil {
		return ""
	}
	return f.Default()
}

// Path returns the current path value.
// 
// Returns:
// 
//   - The path as a string.
func (f *PrivateKeyFlag) Path() string {
	if f == nil {
		return ""
	}
	return f.Value()
}

// Path returns the current path value.
// 
// Returns:
// 
//   - The path as a string.
func (f *PublicKeyFlag) Path() string {
	if f == nil {
		return ""
	}
	return f.Value()
}

// SetPrivateKeyFlag initializes the private key flag.
// 
// Parameters:
// 
//  - flag: The PrivateKeyFlag to initialize.
func SetPrivateKeyFlag(flag *PrivateKeyFlag) {
	if flag != nil {
		flag.SetFlag()
	}
}

// SetPublicKeyFlag initializes the public key flag.
// 
// Parameters:
// 
// - flag: The PublicKeyFlag to initialize.
func SetPublicKeyFlag(flag *PublicKeyFlag) {
	if flag != nil {
		flag.SetFlag()
	}
}