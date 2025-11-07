package des // import "crypto/des"

Package des implements the Data Encryption Standard (DES) and the Triple Data
Encryption Algorithm (TDEA) as defined in U.S. Federal Information Processing
Standards Publication 46-3.

DES is cryptographically broken and should not be used for secure applications.

CONSTANTS

const BlockSize = 8
    The DES block size in bytes.


FUNCTIONS

func NewCipher(key []byte) (cipher.Block, error)
    NewCipher creates and returns a new cipher.Block.

func NewTripleDESCipher(key []byte) (cipher.Block, error)
    NewTripleDESCipher creates and returns a new cipher.Block.


TYPES

type KeySizeError int

func (k KeySizeError) Error() string

