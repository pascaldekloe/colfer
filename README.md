# Colfer

WIP: schema-based binary data format optimized for speed, size, and simplicity.


# Encoding

Colfer documents consists of zero or more field value definitions. Each value is applied in the order of appearance.
Data is represented in a big-endian manner. The format relies on `varints` also known as a [variable-length quantity](https://en.wikipedia.org/wiki/Variable-length_quantity).


## Field Definition

The first byte identifies the field number with it's 7 least significant bits. The most significant bit is applied as follows, depending on the data type.

bool: the actual value
int: the sign: `1` for negative
float: `0` for 32-bit and `1` for 64-bit IEEE 754
timestamp: `0` for seconds, `1` for nanoseconds
blob: always `0`; `1` is reserved for future use

Integers and timestamps are followed by a varint value. Blobs are followed with a varint data length and finally the blob itself.
