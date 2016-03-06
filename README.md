# Colfer [![GoDoc](https://godoc.org/github.com/pascaldekloe/colfer?status.svg)](https://godoc.org/github.com/pascaldekloe/colfer) [![Build Status](https://travis-ci.org/pascaldekloe/colfer.svg?branch=master)](https://travis-ci.org/pascaldekloe/colfer)

WIP: schema-based binary data format optimized for speed, size, and simplicity.
The format is inspired by Proto**col** Buf**fer**.

# Use

```
NAME
	colf — compile Colfer schemas

SYNOPSIS
	colf [-b dir] [-p path] language [file ...]

DESCRIPTION
	Generates source code for the given language. Both go and java are
	supported.
	The file operands are processed in command-line order. If file is
	absent, colf reads all ".colf" files in the working directory.

  -b string
	Use a specific destination base directory. (default ".")
  -p string
	Adds a package prefix. Use slash as a separator when nesting.

BUGS
	Report bugs at https://github.com/pascaldekloe/colfer/issues

SEE ALSO
	protoc(1)
```


# Build

Run `go get github.com/pascaldekloe/colfer/cmd/colf` to install the compiler.

Run `go generate` before the tests.


# Encoding

Data structures start with an 8-bit magic number `0x80` followed by zero or more
field *value definitions*. Only those fields with a value other than the *zero
value* may be serialized. Fields appear in order as stated by the schema.

The zero value for booleans is `false`, integers: `0`, floating points: `0.0`,
timestamps: `1970-01-01T00:00:00.000000000Z` and for text & binary: the empty
string.

Data is represented in a big-endian manner. The format relies on *varints* also
known as a
[variable-length quantity](https://en.wikipedia.org/wiki/Variable-length_quantity).


## Value Definiton

Each definition starts with an 8-bit header. The 7 least significant bits
identify the field by its (0-based position) index in the schema. The most
significant bit is used as a *flag*.

Boolean occurrences set the value to `true`.

Integers are encoded as varints. The header flag indicates negative for signed
types.

Floating points are encoded conform IEEE 754.

Timestamps are encoded as a 64-bit two's complement integer for the number of
seconds that have elapsed since 00:00:00 UTC, Thursday, 1 January 1970, not
counting leap seconds. When the header flag is set then the value is followed
with 32 bits for the nanosecond fraction. Again, a zero value must not be
serialized.

The data for text and binaries is prefixed with a varint size declaration.
