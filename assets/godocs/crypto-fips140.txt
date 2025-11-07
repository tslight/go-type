package fips140 // import "crypto/fips140"


FUNCTIONS

func Enabled() bool
    Enabled reports whether the cryptography libraries are operating in FIPS
    140-3 mode.

    It can be controlled at runtime using the GODEBUG setting "fips140".
    If set to "on", FIPS 140-3 mode is enabled. If set to "only", non-approved
    cryptography functions will additionally return errors or panic.

    This can't be changed after the program has started.

