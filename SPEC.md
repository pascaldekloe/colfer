# Colfer 2 Specification

The data format consists of 3 parts for the following reasons.

* trailing zero compression
* backward and forward data compatibility
* optimal buffer sizes due fixed worst cases
* better branch prediction; fewer jumps
* allow parallel/vector processing

Both integers and floating-points are encoded in little-endian byte order.
Sizes are defined as an (unsigned) integer for the total number of octets.

```bnf
FLIT64        :≡ FLIT64-single | FLIT64-head FLIT64-tail
FLIT64-single :≡ odd-octet
FLIT64-head   :≡ even-octet
FLIT64-tail   :≡ octet |
                 octet octet |
                 octet octet octet |
                 octet octet octet octet |
                 octet octet octet octet octet |
                 octet octet octet octet octet octet |
                 octet octet octet octet octet octet octet |
                 octet octet octet octet octet octet octet octet

octet         :≡ 0 – 255
even-octet    :≡ 0 | 2 | 4 | … | 254
odd-octet     :≡ 1 | 3 | 5 | … | 255
```


## Data Structure Encoding

The first octet in a serial contains the size of the following *fixed part*.

```bnf
struct     :≡ zero | fixed-size fixed-part ranged-part variable-part
zero       :≡ 0
fixed-size :≡ 1 – 255
```

The second octet in a serial contains a FLIT64 for the size of the following
*ranged* and *variable part*. When the size declaration needs more than one
octet, then the rest of the FLIT64 is placed at the beginning of the ranged
part.

```bnf
fixed-part :≡ FLIT64-single fixes | FLIT64-head fixes
```

Each field from an encoded data structure has a fixed element. Note that the
position in the serial does not depend on the actual values.

```bnf
fixes   :≡ fixes fix | ε
fix     :≡ flags | int8 | int16 | float32 | float64 | FLIT64-single | FLIT64-head
flags   :≡ octet
int8    :≡ octet
int16   :≡ octet octet
float32 :≡ octet octet octet octet
float64 :≡ octet octet octet octet octet octet octet octet
```

Booleans are bit flags in octets with the first value in a sequence of eight
using the most significant bit; the second the second most significant bit and
so forth.

Both 8 and 16-bit integers go directly into the fixed part and so are the IEEE
754 floating-points.

All other types encode a FLIT64 with the first octet in the fixed part and the
rest (if any) in the ranged part. Signed integers use ZigZag-encoding.

Binary and text types have their respective size in the ranged FLIT64; arrays
the element count. The actual content goes into the variable part, in reverse
field order.

```bnf
ranged-part   :≡ ranged-part FLIT64-tail | ε
variable-part :≡ variable-part variable | ε

variable      :≡ BLOB | UTF-8 | struct-array | FLIT64-array |
                 int16-array | float32-array | float64-array
BLOB          :≡ BLOB octet | ε
UTF-8         :≡ UTF-8 UTF-8-char | ε
UTF-8-char    :≡ UTF-8-seq-1 | UTF-8-seq-2 | UTF-8-seq-3 | UTF-8-seq-4
UTF-8-seq-1   :≡ 0x00 – 0x7F
UTF-8-seq-2   :≡ 0xC2 – 0xDF UTF-8-tail
UTF-8-seq-3   :≡ 0xE0 0xA0 – 0xBF UTF-8-tail |
                 0xE1 – 0xEC UTF-8-tail UTF-8-tail |
                 0xED 0x80 – 0x9F UTF-8-tail |
                 0xEE – 0xEF UTF-8-tail UTF-8-tail
UTF-8-seq-4   :≡ 0xF0 0x90 – 0xBF UTF-8-tail UTF-8-tail |
                 0xF1 – 0xF3 UTF-8-tail UTF-8-tail UTF-8-tail |
                 0xF4 0x80 – 0x8F UTF-8-tail UTF-8-tail
UTF-8-tail    :≡ 0x80 – 0xBF
struct-array  :≡ struct-array struct | ε
FLIT64-array  :≡ FLIT64-array FLIT64 | ε
int16-array   :≡ int16-array int16 | ε
float32-array :≡ float32-array float32 | ε
float64-array :≡ float64-array float64 | ε
```

Fields not encoded default to the *zero value*, which is false for booleans,
zero for numbers and the empty set for text, binaries and arrays. Additional
(newer) fields may be ignored.


## References

* [Fixed-Length Integer Trim algorithm](https://github.com/pascaldekloe/flit)
* [UTF-8, a transformation format of ISO 10646](https://tools.ietf.org/rfc/rfc3629.txt)
* [IEEE Standard for Floating-Point Arithmetic](https://ieeexplore.ieee.org/document/4610935/)
