# Colfer Version 2 — Format Specification

## Abstract

Version 2 of the data format trades simplicity for peak performance.

* backward and forward compatibility
* maximized number of fixed data positions
* minimized number of boundary-checks required
* designed for parallel processing options
* sane size boundaries on all attributes


Colfer is an octet oriented data format.

```bnf
octet :≡ 0..255 ; 8 bits
```


## Integer Compression

Compressed integers consist of one header octet each. Up to eight more octets
may follow in the tail, in little-endian order. Compression should select the
longest `integer-tail` with a `non-zero` trailer, which can be none (`ε`).

```bnf
integer-head :≡ octet ;
integer-tail :≡ ε |
                non-zero |
                octet non-zero |
                octet octet non-zero |
                octet octet octet non-zero |
                octet octet octet octet non-zero |
                octet octet octet octet octet non-zero |
                octet octet octet octet octet octet non-zero |
                octet octet octet octet octet octet octet non-zero ;

non-zero     :≡ 1..255 ;
```

The header octet counts the number of octets in the tail with its trailing-zero
count, range 0..8. Any remaning bits in the header (denoted by `x` in the table
below) hold the least-significant bits of the compressed integer. The rest from
the compressed integer overflows to `integer-tail`.

| Tail Size | Header Bits  | Value Range                         |
|:----------|:-------------|:------------------------------------|
| 0 octets  | `xxxx xxx1`  | 7 bit (128)                         |
| 1 octet   | `xxxx xx10`  | 14 bit (16 384)                     |
| 2 octets  | `xxxx x100`  | 21 bit (2 097 152)                  |
| 3 octets  | `xxxx 1000`  | 28 bit (268 435 456)                |
| 4 octets  | `xxx1 0000`  | 35 bit (34 359 738 368)             |
| 5 octets  | `xx10 0000`  | 42 bit (4 398 046 511 104)          |
| 6 octets  | `x100 0000`  | 49 bit (562 949 953 421 312)        |
| 7 octets  | `1000 0000`  | 56 bit (72 057 594 037 927 936)     |
| 8 octets  | `0000 0000`  | 64 bit (18 446 744 073 709 551 616) |

Signed integers reorganise their bits conform the ZigZag[^1] algorithm before
compression.


## Size Profile

Colfer encoding starts with a header. The three least-significant bits from the
first octet select a size profile. Numeric value 0 selects *compact*, 1 selects
*wide*, and 2 selects *royal*.

```bnf
encoding-head :≡ compact-head | wide-head | royal-head ;
compact-head  :≡ compact-flag octet octet ; 24 bits
wide-head     :≡ wide-flag octet octet octet octet ; 40 bits
royal-head    :≡ royal-flag octet octet octet
                 octet octet octet ; 56 bits

compact-flag  :≡ 0..31 × 8 ;
wide-flag     :≡ 0..31 × 8 + 1 ;
royal-flag    :≡ 0..31 × 8 + 2 ;
```

Each profile has its own limits on the overall size, the UTF-8 size per `text`
field, the element count per list field, and the `fixed` size (explained in the
following section).

| Profile | Encoding Limit | Fixed Data Limit | UTF-8 Limit  | List Limit |
|:--------|:---------------|:-----------------|:-------------|:-----------|
| compact | 4 KiB          | 512 B            | 255 B        | 255        |
| wide    | 2 MiB          | 64 KiB           | 64 KiB − 1 B | 65 535     |
| royal   | 512 MiB        | 16 MiB           | 16 MiB − 1 B | 16 777 215 |

Headers with the compact profile have 12 bits for the total size, and 9 bits for
the fixed data size. The little-endian value of a `compact-head` equals the
index of the last octet in the encoding multiplied by 8, plus the index of the
last octet in the fixed section multiplied by 32 768.

     0                   1                   2
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3
    ┌─────┬───────────────────────┬─────────────────┐
    │0 0 0│ total size − 1        │ fixed size − 1  │
    └─────┴───────────────────────┴─────────────────┘

Headers with the wide profile have 21 bits for the total size, and 16 bits for
the fixed data size. The little-endian value of a `wide-head` equals 1, plus the
index of the last octet in the encoding multiplied by 8, plus the index of the
last octet in the fixed section multiplied by 16 777 216.

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

Headers with the royal profile have 29 bits for the total size, and 24 bits for
the fixed data size. The little-endian value of a `royal-head` equals 2, plus
the index of the last octet in the encoding multiplied by 8, plus the index of
the last octet in the fixed section multiplied by 536 870 912.

     0                   1                   2                 3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 9 0 1
    ┌─────┬───────────────────────────────────────────────────────┐
    │0 1 0│ total size − 1                                        │
    └─────┴───────────────────────────────────────────────────────┘

                     4                 5
     2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 9 0 1 2 3 4 5 6
    ┌───────────────────────────────────────────────┐
    │ fixed size − 1                                │
    └───────────────────────────────────────────────┘

ℹ Note that logically the total size is greater than or equal to the fixed size,
and that the fixed size greater than the header size. As such, every encoding is
at least 4 octets in length, and with the first 4 octets read, the total size is
known.


## Data Structure

Data is split in two sections, namely `fixed` and `variable`. Data types with a
fixed-length encoding append to `fixed` only. Data types with a variable-length
encoding append their minimum size to `fixed`, and the remainder overflows to
`variable`. For example, integer compression has its `integer-head` octet (at a
known location) in `fixed`, and its `integer-tail` in `variable`.

```bnf
encoding      :≡ fixed variable ; packs one data structure
fixed         :≡ encoding-head field-fixes ;
field-fixes   :≡ fix field-fixes | fix ;

variable      :≡ overflow payloads ;
overflow      :≡ ε | integer-tail overflow ;
payloads      :≡ ε | payload payloads ;
```

Fields append in sequential order to `fixed`. However, booleans group by 8 in
the form of little-endian *bit fields*. Thus, the first 8 boolean fields reside
at the position of the first field, and the next 8 booleans at the position of
the ninth field, and so forth.

All `fix` values encode in little-endian order. Opaque and floating-point[^2]
values reside as is without compression. Timestamp values encode nanoseconds in
the 30 least-significant bits, and the 34 most-significant bits hold the number
of seconds that have elapsed since 00:00:00 UTC on 1 January 1970, not counting
leap seconds.

`List-size` contains the element count as either one, two or three octets,
depending on the size profile. `Text-size` contains the UTF-8 octet count,
again, as either one, two or three octets, depending on the size profile.

Arrays with a fixed size encode their elements just as separate fields would do.
Nested data structures also encode their fields inline, all as if they were part
of the hosting data structure.

```bnf
fix       :≡ opaque8 | opaque16 | opaque32 | opaque64 |
             integer-head | float32 | float64 | timestamp |
             text-size | list-size | bit-field ;
opaque8   :≡ octet ;
opaque16  :≡ octet octet ;
opaque32  :≡ octet octet octet octet ;
opaque64  :≡ octet octet octet octet ;
float32   :≡ octet octet octet octet ; IEEE floating-point
float64   :≡ octet octet octet octet
             octet octet octet octet ; IEEE floating-point
timestamp :≡ octet octet octet octet
             octet octet octet octet ;
text-size :≡ octet | octet octet | octet octet octet ;
list-size :≡ octet | octet octet | octet octet octet ;
bit-field :≡ octet ;
```

Remainders of `fixed` appear in `overflow` in corresponding order. The `payload`
components append in reverse field order, this to allow for gaps of unknown data
from format revisions with more fields.

Text encodes as UTF-8[^3], within Unicode range U+0000..U+10FFFF.

```
payload   :≡ text | list ;
text      :≡ ε | utf-8 text ;
list      :≡ opaque8-list | opaque16-list | opaque32-list |
             opaque64-list | float32-list | float64-list |
             timestamp-list ;

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

Both `text` and `list` must match their respective `text-size` and `list-size`
declarations in the fixed section. Text lists have a `text-size` in the payload
section for each of its `text` entries, in corresponding order.

```bnf
opaque8-list   :≡ ε | octet opaque8-list ;
opaque16-list  :≡ ε | octet octet opaque16-list ;
opaque32-list  :≡ ε | octet octet octet octet opaque32-list ;
opaque64-list  :≡ ε | octet octet octet octet
                  octet octet octet octet opaque64-list ;
float32-list   :≡ ε | float32 float32-list ;
float64-list   :≡ ε | float64 float64-list ;
text-list      :≡ text-sizes text ;
text-sizes     :≡ ε | text-size text-sizes ;
timestamp-list :≡ ε | timestamp timestamp-list ;
```


## References

[^1]: [(Protocol Buffers) ZigZag encoding](https://developers.google.com/protocol-buffers/docs/encoding#signed-integers)
[^2]: [IEEE Standard for Floating-Point Arithmetic](https://ieeexplore.ieee.org/document/4610935/)
[^3]: [UTF-8, a transformation format of ISO 10646](https://tools.ietf.org/rfc/rfc3629.txt)
