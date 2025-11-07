package crc64 // import "hash/crc64"

Package crc64 implements the 64-bit cyclic redundancy check, or CRC-64,
checksum. See https://en.wikipedia.org/wiki/Cyclic_redundancy_check for
information.

CONSTANTS

const (
	// The ISO polynomial, defined in ISO 3309 and used in HDLC.
	ISO = 0xD800000000000000

	// The ECMA polynomial, defined in ECMA 182.
	ECMA = 0xC96C5795D7870F42
)
    Predefined polynomials.

const Size = 8
    The size of a CRC-64 checksum in bytes.


FUNCTIONS

func Checksum(data []byte, tab *Table) uint64
    Checksum returns the CRC-64 checksum of data using the polynomial
    represented by the Table.

func New(tab *Table) hash.Hash64
    New creates a new hash.Hash64 computing the CRC-64 checksum using the
    polynomial represented by the Table. Its Sum method will lay the value
    out in big-endian byte order. The returned Hash64 also implements
    encoding.BinaryMarshaler and encoding.BinaryUnmarshaler to marshal and
    unmarshal the internal state of the hash.

func Update(crc uint64, tab *Table, p []byte) uint64
    Update returns the result of adding the bytes in p to the crc.


TYPES

type Table [256]uint64
    Table is a 256-word table representing the polynomial for efficient
    processing.

func MakeTable(poly uint64) *Table
    MakeTable returns a Table constructed from the specified polynomial.
    The contents of this Table must not be modified.

