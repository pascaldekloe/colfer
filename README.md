# Colfer [![GoDoc](https://godoc.org/github.com/pascaldekloe/colfer/go?status.svg)](https://godoc.org/github.com/pascaldekloe/colfer/go) [![Build Status](https://travis-ci.org/pascaldekloe/colfer.svg?branch=master)](https://travis-ci.org/pascaldekloe/colfer)

WIP: schema-based binary data format optimized for speed, size, and simplicity.
The format is inspired by Proto**col** Buf**fer**.


# Encoding

A Colfer representation starts with an 8-bit magic number `0x80` followed by
zero or more field value definitions. Only those fields with a value other than
the zero value may be serialized. Fields appear by field number, in ascending
order.

Data is represented in a big-endian manner. The format relies on *varints* also
known as a
[variable-length quantity](https://en.wikipedia.org/wiki/Variable-length_quantity).


## Value Definiton

Each definition starts with an 8-bit *key*. The 7 least significant bits
identify the field by its numeric value.

Field occurences of type `bool` set the value to `true`.

Types `uint64` and `uint32` are encoded as varints. The most significant bit on
the key means negative for signed types `int64` and `int32`.

Types `float64` and `float32` are encoded conform IEEE 754.

The 64 bits for `timestamp` have the two's complement representation of the
number of miliseconds that have elapsed since 00:00:00 UTC, Thursday, 1 January
1970, not counting leap seconds. If the most significant bit on the key is set
then the following 32 bits contain the nanosecond fraction.

The data for types `text` and `binary` is prefixed with a varint size
declaration.
