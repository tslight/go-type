package adler32 // import "hash/adler32"

Package adler32 implements the Adler-32 checksum.

It is defined in RFC 1950:

    Adler-32 is composed of two sums accumulated per byte: s1 is
    the sum of all bytes, s2 is the sum of all s1 values. Both sums
    are done modulo 65521. s1 is initialized to 1, s2 to zero.  The
    Adler-32 checksum is stored as s2*65536 + s1 in most-
    significant-byte first (network) order.

CONSTANTS

const Size = 4
    The size of an Adler-32 checksum in bytes.


FUNCTIONS

func Checksum(data []byte) uint32
    Checksum returns the Adler-32 checksum of data.

func New() hash.Hash32
    New returns a new hash.Hash32 computing the Adler-32 checksum. Its Sum
    method will lay the value out in big-endian byte order. The returned Hash32
    also implements encoding.BinaryMarshaler and encoding.BinaryUnmarshaler to
    marshal and unmarshal the internal state of the hash.

