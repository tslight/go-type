package rc4 // import "crypto/rc4"

Package rc4 implements RC4 encryption, as defined in Bruce Schneier's Applied
Cryptography.

RC4 is cryptographically broken and should not be used for secure applications.

TYPES

type Cipher struct {
	// Has unexported fields.
}
    A Cipher is an instance of RC4 using a particular key.

func NewCipher(key []byte) (*Cipher, error)
    NewCipher creates and returns a new Cipher. The key argument should be the
    RC4 key, at least 1 byte and at most 256 bytes.

func (c *Cipher) Reset()
    Reset zeros the key data and makes the Cipher unusable.

    Deprecated: Reset can't guarantee that the key will be entirely removed from
    the process's memory.

func (c *Cipher) XORKeyStream(dst, src []byte)
    XORKeyStream sets dst to the result of XORing src with the key stream.
    Dst and src must overlap entirely or not at all.

type KeySizeError int

func (k KeySizeError) Error() string

