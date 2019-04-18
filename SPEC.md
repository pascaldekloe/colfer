# Colfer Version 2 Specification


## Abstract

Colfer is an octet oriented data format.

```bnf
octet :≡ 0–255 ; 8-bit byte
```

Version 2 trades simplicity for extreme performance.

* backward and forward compatibility
* optimal buffer sizes due to fixed worst cases
* maximise branch prediction
* parallel/vector processing options
* single field extraction options


## Integers

All integers are 64-bit wide. Signed integers use *ZigZag encoding*.

Leading zeros are ommitted with the following algorithm. TODO(pascaldekloe)

```bnf
integer      :≡ integer-head integer tail ;
integer-head :≡ octet ;
integer-tail :≡ octet |
                octet octet |
                octet octet octet |
                octet octet octet octet |
                octet octet octet octet octet |
                octet octet octet octet octet octet |
                octet octet octet octet octet octet octet |
                octet octet octet octet octet octet octet octet ;
```


## Record Encoding

Data with a fixed length encoding is placed at the start of a serial to maximise
the number of stationary positions. Variable-length data is split in two parts.
The minimum [size] goes into `fixed` and the remainder, if any, overflows into
`variable`.

```bnf
serial        :≡ fixed-size variable-size fixed variable ; record
fixed-size    :≡ 0–255 0–255 ; 16-bit address space, little-endian order
variable-size :≡ integer-head ; 64-bit address space
```

The `variable-size` may have an `integer-tail` in the beginning of `variable`,
which starts at the 4th octet + the numeric value of `fixed-size`.

Fields appear in sequential order in `fixed`. However, booleans group by 8 in
the form of little-endian *bit fields*. Thus the first 8 boolean fields reside
at the position of the first field, the next 8 at the position of the ninth, and
so forth.

```bnf
fixed     :≡ ε | fixed fix ;
fix       :≡ integer-head | bit-field | float32 | float64 |
             opaque8   | opaque16  | opaque24  | opaque32  |
             opaque40  | opaque48  | opaque56  | opaque64  |
             opaque72  | opaque80  | opaque88  | opaque96  |
             opaque104 | opaque112 | opaque120 | opaque128 |
             opaque136 | opaque144 | opaque152 | opaque160 |
             opaque168 | opaque176 | opaque184 | opaque192 |
             opaque200 | opaque208 | opaque216 | opaque224 |
             opaque232 | opaque240 | opaque248 | opaque256 ;
bit-field :≡ octet ; little-endian bit order
float32   :≡ octet octet octet octet ; big-endian order
float64   :≡ octet octet octet octet octet octet octet octet ; big-endian order
opaque8   :≡ octet ;
opaque16  :≡ octet octet ;
opaque24  :≡ octet octet octet ;
…
opaque256 :≡ octet octet octet octet octet octet octet octet
             octet octet octet octet octet octet octet octet
             octet octet octet octet octet octet octet octet
             octet octet octet octet octet octet octet octet ;
```

Remainders of `fixed` appear in `overflow` in sequential order. The `embedded`
parts append in reverse field order, this to support unknow/newer field gaps.

```bnf
variable :≡ overflow embedded
overflow :≡ ε | overflow integer-tail ;
embedded :≡ ε | embedded text | embedded array ;
```

Text fields encode the octet size as an integer (in `fixed` and `overflow`) and
the actual *UTF-8* payload goes into `embedded`.

```
text      :≡ text char | char ;
char      :≡ char8 | char16 | char24 | char32 ;
char8     :≡ 0x00–0x7F ;
char16    :≡ 0xC2–0xDF char-tail ;
char24    :≡ 0xE0 0xA0–0xBF char-tail |
             0xE1–0xEC char-tail char-tail |
             0xED 0x80–0x9F char-tail |
             0xEE–0xEF char-tail char-tail ;
char32    :≡ 0xF0 0x90–0xBF char-tail char-tail |
             0xF1–0xF3 char-tail char-tail char-tail |
             0xF4 0x80–0x8F char-tail char-tail ;
char-tail :≡ 0x80–0xBF ;
```

Array fields encode the element count as an integer (in `fixed` and `overflow`).
Text arrays contain an integer array with the octet sizes of each corresponding
element.

```bnf
array :≡ integer-array   | boolean-array   | float32-array   | float64-array   |
         opaque8-array   | opaque16-array  | opaque24-array  | opaque32-array  |
         opaque40-array  | opaque48-array  | opaque56-array  | opaque64-array  |
         opaque72-array  | opaque80-array  | opaque88-array  | opaque96-array  |
         opaque104-array | opaque112-array | opaque120-array | opaque128-array |
         opaque136-array | opaque144-array | opaque152-array | opaque160-array |
         opaque168-array | opaque176-array | opaque184-array | opaque192-array |
         opaque200-array | opaque208-array | opaque216-array | opaque224-array |
         opaque232-array | opaque240-array | opaque248-array | opaque256-array |
         text-array      | record array    ;

integer-array      :≡ integer-head-array integer-tail-array ;
integer-head-array :≡ integer-head-array integer-head | integer-head ;
integer-tail-array :≡ integer-tail-array integer-tail | integer-tail ;

boolean-array   :≡ boolean-array bit-field | bit-field ;
float32-array   :≡ float32-array float32 | float32 ;
float64-array   :≡ float64-array float64 | float64 ;
opaque8-array   :≡ opaque8-array opaque8 | opaque8 ;
opaque16-array  :≡ opaque16-array opaque16 | opaque16 ;
opaque24-array  :≡ opaque24-array opaque24 | opaque24 ;
…
opaque256-array :≡ opaque256-array opaque256 | opaque256 ;

text-array    :≡ integer-array text ;
record-array  :≡ record-array record | record ;
```


## References

* [UTF-8, a transformation format of ISO 10646](https://tools.ietf.org/rfc/rfc3629.txt)
* [IEEE Standard for Floating-Point Arithmetic](https://ieeexplore.ieee.org/document/4610935/)
* [(Protocol Buffers) ZigZag encoding](https://developers.google.com/protocol-buffers/docs/encoding#signed-integers)
