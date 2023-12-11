# Colfer Version 2 — Format Specification

## Abstract

Version 2 trades simplicity for peak performance.

* backward and forward compatibility
* fixed data positions maximised
* fewer bound-checks required
* parallel/vector processing options


Colfer is an octet oriented data format with built-in integer compression.

```bnf
octet :≡ 0–255 ; 8 bits
```


## Integer Encoding

Integer values encode one header octet each. Up to eight more octets may follow
in the tail. Encoders should select the smallest tail which is can still carry
all input data.

```bnf
integer-head :≡ octet ;
integer-tail :≡ ε |
                octet |
                octet octet |
                octet octet octet |
                octet octet octet octet |
                octet octet octet octet octet |
                octet octet octet octet octet octet |
                octet octet octet octet octet octet octet |
                octet octet octet octet octet octet octet octet ;
```

The header octet counts the number of octets in the tail with its trailing-zero
count, range 0–8. The least-significant bits of the integer value reside in the
tail, in little-endian order. Any remaining bits (which are the most significant
ones) follow the size flag in the header, denoted by `x` in the following table.

| Tail Size | Header Bits  | Range                               |
|:----------|:-------------|:------------------------------------|
| 0 octet   | `xxxx xxx1`  | 7-bit (128)                         |
| 1 octets  | `xxxx xx10`  | 14-bit (16'384)                     |
| 2 octets  | `xxxx x100`  | 21-bit (2'097'152)                  |
| 3 octets  | `xxxx 1000`  | 28-bit (268'435'456)                |
| 4 octets  | `xxx1 0000`  | 35-bit (34'359'738'368)             |
| 5 octets  | `xx10 0000`  | 42-bit (4'398'046'511'104)          |
| 6 octets  | `x100 0000`  | 49-bit (562'949'953'421'312)        |
| 7 octets  | `1000 0000`  | 56-bit (72'057'594'037'927'936)     |
| 8 octets  | `0000 0000`  | 64-bit (18'446'744'073'709'551'616) |

Signed integers pack their bits conform the ZigZag[^1] algorithm first.


## Size Profile

Colfer encoding starts with a header. The three least-significant bits from the
first octet select a size profile. Numeric value 0 selects *compact*, 1 selects
*wide*, and 2 selects *large*.

```bnf
encoding-head :≡ compact-header | wide-header | large-header
compact-head  :≡ octet octet octet ; 24 bits
wide-head     :≡ octet octet octet octet octet ; 40 bits
large-head    :≡ octet octet octet octet octet octet octet ; 56 bits
```

Each profile has its own limits on the overall size, the UTF-8 size per `text`
field, the element count per list field, and the `fixed` size (explained in the
following section).

| Profile | Encoding Limit | Fixed Data Limit | UTF-8 Limit  | List Limit |
|:--------|:---------------|:-----------------|:-------------|:-----------|
| compact | 4 KiB          | 512 B            | 255 B        | 255        |
| wide    | 2 MiB          | 64 KiB           | 64 KiB − 1 B | 65535      |
| large   | 521 MiB        | 64 KiB           | 16 MiB − 1 B | 16777215   |

Headers with the compact profile have 12 bits for the total size, and 9 bits for
the fixed data size. The numeric value of a `compact-header` equals the index of
the last octet in the encoding multiplied by 8, plus the index of the last octet
in the fixed section multiplied by 32768.

     0                   1                   2
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3
    ┌─────┬───────────────────────┬─────────────────┐
    │0 0 0│ total size − 1        │ fixed size − 1  │
    └─────┴───────────────────────┴─────────────────┘

Headers with the wide profile have 21 bits for the total size, and 16 bits for
the fixed data size. The numeric value of a `wide-header` equals 1, plus the
index of the last octet in the encoding multiplied by 8, plus the index of the
last octet in the fixed section multiplied by 16777216.

     0                   1                   2
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3
    ┌─────┬─────────────────────────────────────────┐
    │1 0 0│ total size − 1                          │
    └─────┴─────────────────────────────────────────┘

               3
     4 5 6 7 9 0 1 2 3 4 5 6 7 8 9
    ┌─────────────────────────────┐
    │ fixed size − 1              │
    └─────────────────────────────┘

Headers with the large profile have 29 bits for the total size, and 16 bits for
the fixed data size. The numeric value of a `large-header` equals 2, plus the
index of the last octet in the encoding multiplied by 8, plus the index of the
last octet in the fixed section multiplied by 536870912.

     0                   1                   2                 3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 9 0 1
    ┌─────┬───────────────────────────────────────────────────────┐
    │0 1 0│ total size − 1                                        │
    └─────┴───────────────────────────────────────────────────────┘

                     4
     2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7
    ┌───────────────────────────────┐
    │ fixed size − 1                │
    └───────────────────────────────┘

ℹ Note that logically the total size is greater than or equal to the fixed size,
and that the fixed size greater than the header size. As such, every encoding is
at least 4 octets in length, and with 4 octets read the total size is known.


## Data Structure Encoding

Data is split in two sections, namely `fixed` and `variable`. Data types with a
fixed-length encoding append to `fixed` only. Data types with a variable-length
encoding append their minimum size to `fixed`, and any remainder overflows to
`variable`. For example, integer encoding as described in the previous section
has its `integer-head` octet (at a known location) in `fixed`. When the head
value calls for more data, then its `integer-tail` appends to `variable`.

```bnf
encoding      :≡ fixed variable ; packs one data structure
fixed         :≡ header field-fixes ;
field-fixes   :≡ fix field-fixes | fix ;

variable      :≡ overflow payloads ;
overflow      :≡ integer-tail overflow | ε ;
payloads      :≡ payload payloads | ε ;
```

Fields append in sequential order to `fixed`. However, booleans group by 8 in
the form of little-endian *bit fields*. Thus, the first 8 boolean fields reside
at the position of the first field, and the next 8 booleans at the position of
the ninth field, and so forth.

Single and double-precision floating-points[^2] encode without compression, in
big-endian octet-order. Nested data structures encode their fields inline, as if
they were part of the hosting data structure.

Arrays with a fixed size encode their elements just as separate fields would do.
Nested data structures also encode their fields inline, all as if they were part
of the hosting data structure.

```bnf
fix       :≡ integer-head | bit-field | float32 | float64 | nested |
             opaque8  | opaque16 | opaque24 | … | opaque524288 ;
bit-field :≡ octet ; little-endian bit-order, zero padded
float32   :≡ octet octet octet octet ; IEEE floating-point
float64   :≡ octet octet octet octet
             octet octet octet octet ; IEEE floating-point
opaque8   :≡ octet ;
opaque16  :≡ octet octet ;
opaque24  :≡ octet octet octet ;
…
```

Remainders of `fixed` appear in `overflow` in corresponding order. The `payload`
components append in reverse field order to support unknown gaps from encodings
with more fields for which the data type is unkown.

Opaque data with a variable size is copied as is into a `payload` section with
the octet count encoded as an integer (in `fixed` and `overflow`). Text encodes
similar, with `payload` as UTF-8[^3].

```
payload   :≡ opaque payload | text payload | list payload ;
opaque    :≡ octet opaque | ε ;
text      :≡ utf-8 text | ε ;
utf-8     :≡ utf-seq-1 | utf-seq-2 | utf-seq-3 | utf-seq-4 ;
utf-seq-1 :≡ 0x00–0x7F ;
utf-seq-2 :≡ 0xC2–0xDF utf-tail ;
utf-seq-3 :≡ 0xE0 0xA0–0xBF utf-tail |
             0xE1–0xEC utf-tail utf-tail |
             0xED 0x80–0x9F utf-tail |
             0xEE–0xEF utf-tail utf-tail ;
utf-seq-4 :≡ 0xF0 0x90–0xBF utf-tail utf-tail |
             0xF1–0xF3 utf-tail utf-tail utf-tail |
             0xF4 0x80–0x8F utf-tail utf-tail ;
utf-tail  :≡ 0x80–0xBF ;
```

Lists encode their element count as an integer (in `fixed` and `overflow`). The
`payload` for integers is a FLIT64[^4] sequence in ascending order. Booleans do
little-endian bit-order. Strings encode their element's octet count as a FLIT64
sequence folowed by each value as UTF-8.

```bnf
list           :≡ boolean-list | float32-list | float64-list | integer-list |
                  text-list | opaque8-list | … | structure-list ;
boolean-list   :≡ bit-field boolean-list | ε ;
float32-list   :≡ float32 float32-list | ε ;
float64-list   :≡ float64 float64-list | ε ;
integer-list   :≡ flit64 integer-list | ε ;
text-list      :≡ integer-list | text ;
opaque8-list   :≡ opaque8 opaque8-list | ε ;
opaque16-list  :≡ opaque16 opaque8-list | ε ;
opaque24-list  :≡ opaque24 opaque8-list | ε ;
…
structure-list :≡ encoding structure-list | ε ;
```


## References

[^1]: [(Protocol Buffers) ZigZag encoding](https://developers.google.com/protocol-buffers/docs/encoding#signed-integers)
[^2]: [IEEE Standard for Floating-Point Arithmetic](https://ieeexplore.ieee.org/document/4610935/)
[^3]: [UTF-8, a transformation format of ISO 10646](https://tools.ietf.org/rfc/rfc3629.txt)
[^4]: [Fixed-Length Integer Trim (FLIT)](https://github.com/pascaldekloe/flit)
