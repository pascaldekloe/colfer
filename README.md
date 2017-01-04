# Colfer [![Build Status](https://travis-ci.org/pascaldekloe/colfer.svg?branch=master)](https://travis-ci.org/pascaldekloe/colfer)

Colfer is a schema-based binary serialization format optimized for speed and
size.

The project's compiler `colf(1)` generates source code from schema definitions
to marshal and unmarshall data structures.

This is free and unencumbered software released into the
[public domain](http://creativecommons.org/publicdomain/zero/1.0).
The format is inspired by Proto**col** Buf**fer**.


#### Language Support

* C, C99 compliant, C++13 compliant, WIP, API might change
* Go, a.k.a. golang
* Java, Android compatible
* JavaScript, a.k.a. ECMAScript, NodeJS compatible

#### Features

* Simple and straightforward in use
* No dependencies other than the core library
* Both faster and smaller than: Protocol Buffers, FlatBuffers and MessagePack
* Robust including size protection
* Maximum of 127 fields per data structure
* No support for enumerations
* Framed; suitable for concatenation/streaming

#### TODO's

* RMI (WIP
[![GoDoc](https://godoc.org/github.com/pascaldekloe/colfer/rpc?status.svg)](https://godoc.org/github.com/pascaldekloe/colfer/rpc)
)
* List type support for integers and timestamps
* Please [share](https://github.com/pascaldekloe/colfer/wiki/Users#production-use) your experiences



## Use

Download a [prebuilt compiler](https://github.com/pascaldekloe/colfer/releases)
or run `go get -u github.com/pascaldekloe/colfer/cmd/colf` to make one yourself.
Without arguments the command prints its manual.

```
NAME
	colf — compile Colfer schemas

SYNOPSIS
	colf [ options ] language [ file ... ]

DESCRIPTION
	Generates source code for a language. The options are: C, Go,
	Java and JavaScript.
	The file operands specify the input. Directories are scanned for
	files with the colf extension. If file is absent, colf includes
	the working directory.
	A package can have multiple schema files.

OPTIONS
  -b directory
    	Use a specific destination base directory. (default ".")
  -f	Normalizes schemas on the fly.
  -l expression
    	Sets the default upper limit for the number of elements in a
    	list. The expression is applied to the target language under the
    	name ColferListMax. (default "64 * 1024")
  -p prefix
    	Adds a package prefix. Use slash as a separator when nesting.
  -s expression
    	Sets the default upper limit for serial byte sizes. The
    	expression is applied to the target language under the name
    	ColferSizeMax. (default "16 * 1024 * 1024")
  -v	Enables verbose reporting to the standard error.

EXIT STATUS
	The command exits 0 on succes, 1 on compilation failure and 2
	when invoked without arguments.

EXAMPLES
	Compile ./api/*.colf into ./src/ as Java:

		colf -p com/example -b src java api

	Compile ./io.colf with compact limits as C:

		colf -s 2048 -l 96 c io.colf

BUGS
	Report bugs at https://github.com/pascaldekloe/colfer/issues

SEE ALSO
	protoc(1)
```


It is recommended to commit the generated source code to the respective version
control.
Maven users may [disagree](https://github.com/pascaldekloe/colfer/wiki/Java#maven).



## Schema

Data structures are defined in `.colf` files. The format is quite conventional.

```
// Package demo offers a demonstration.
// These comment lines will end up in the generated code.
package demo

// Course is the grounds where the game of golf is played.
type Course struct {
	ID    uint64
	name  text
	holes []hole
	image binary
	tags  []text
}

type hole struct {
	// Lat is the latitude of the cup.
	lat float64
	// Lon is the longitude of the cup.
	lon float64
	// Par is the difficulty index.
	par uint8
	// Water marks the presence of water.
	water bool
	// Sand marks the presence of sand.
	sand bool
}
```

[See](https://gist.github.com/pascaldekloe/f5f15729cceefe430c9858d58e0dd1a3)
what the generated code looks like.

The following table shows how Colfer data types are applied per language.

| Colfer	| C		| Go		| Java		| JavaScript	|
|:--------------|:--------------|:--------------|:--------------|:--------------|
| bool		| bool		| bool		| boolean	| Boolean	|
| uint8		| uint8_t	| uint8		| byte †	| Number	|
| uint16	| uint16_t	| uint16	| short †	| Number	|
| uint32	| uint32_t	| uint32	| int †		| Number	|
| uint64	| uint64_t	| uint64	| long †	| Number ‡	|
| int32		| int32_t	| int32		| int		| Number	|
| int64		| int64_t	| int64		| long		| Number ‡	|
| float32	| float		| float32	| float		| Number	|
| float64	| double	| float64	| double	| Number	|
| timestamp	| 2 × time_t	| Time ††	| Instant	| Date + Number	|
| text		| char †‡	| string	| String ‡‡	| String ‡‡	|
| binary	| uint8_t †‡ 	| []byte	| byte[]	| Uint8Array	|
| list		| †‡		| slice		| array		| Array		|

* † signed representation of unsigned data, i.e. may overflow to negative.
* ‡ range limited to (1 - 2⁵³, 2⁵³ - 1)
* †† timezone not preserved
* †‡ struct of pointer + size_t
* ‡‡ characters limited by UTF-16 (`U+0000`, `U+10FFFF`)

Lists may contain floating points, text, binaries or data structures.


## Compatibility

Name changes do not affect the serialization format. Deprecated fields can be
renamed to clearly discourage its use.

The following changes are backward compatible.
* New fields at the end of Colfer structs
* Change datatype int32 into int64
* Change datatype text into binary



## Performance

```
BenchmarkMarshal/colfer-8   	20000000	        65.9 ns/op	      48 B/op	       1 allocs/op
BenchmarkMarshal/protobuf-8 	20000000	        81.4 ns/op	      52 B/op	       1 allocs/op
BenchmarkMarshal/flatbuf-8  	 2000000	       701 ns/op	     472 B/op	      12 allocs/op
BenchmarkUnmarshal/colfer-8 	20000000	        94.3 ns/op	      84 B/op	       2 allocs/op
BenchmarkUnmarshal/protobuf-8         	10000000	       128 ns/op	      84 B/op	       2 allocs/op
BenchmarkUnmarshal/flatbuf-8          	10000000	       152 ns/op	      84 B/op	       2 allocs/op
BenchmarkMarshalReuse/colfer-8        	50000000	        36.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalReuse/protobuf-8      	30000000	        48.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalReuse/flatbuf-8       	 5000000	       294 ns/op	       0 B/op	       0 allocs/op
BenchmarkUnmarshalReuse/colfer-8      	20000000	        62.7 ns/op	      20 B/op	       1 allocs/op
BenchmarkUnmarshalReuse/protobuf-8    	20000000	        92.5 ns/op	      20 B/op	       1 allocs/op
BenchmarkUnmarshalReuse/flatbuf-8     	10000000	       119 ns/op	      20 B/op	       1 allocs/op
```



## Format

Data structures consist of zero or more field *value definitions* followed by a
termination byte `0x7f`. Only those fields with a value other than the *zero
value* may be serialized. Fields appear in order as stated by the schema.

The zero value for booleans is `false`, integers: `0`, floating points: `0.0`,
timestamps: `1970-01-01T00:00:00.000000000Z`, text & binary: the empty
string, nested data structures: `null` and an empty list for data structure
lists.

Data is represented in a big-endian manner. The format relies on *varints* also
known as a
[variable-length quantity](https://en.wikipedia.org/wiki/Variable-length_quantity).
Bits reserved for future use (*RFU*) must be set to 0.


#### Value Definiton

Each definition starts with an 8-bit header. The 7 least significant bits
identify the field by its (0-based position) index in the schema. The most
significant bit is used as a *flag*.

Boolean occurrences set the value to `true`. The flag is RFU.

Unsigned 8-bit integer values just follow the header byte and the flag is RFU.
Unsigned 16-bit integer values greather than 255 also follow the header byte.
Smaller values are encoded in one byte with the header flag set.
Unsigned 32-bit integer values less than 1<<21 are encoded as varints and
larger values set the header flag for fixed length encoding.
Unsigned 64-bit integer values less than 1<<49 are encoded as varints and
larger values set the header flag for fixed length encoding.

Signed 32- and 64-bit integers are encoded as varints. The flag stands for
negative. The tenth byte for 64-bit integers is skipped for encoding since its
value is fixed to `0x01`.

Floating points are encoded conform IEEE 754. The flag is RFU.

Timestamps are encoded as a 32-bit unsigned integer for the number of seconds
that have elapsed since 00:00:00 UTC, Thursday, 1 January 1970, not counting
leap seconds. When the header flag is set then the number of seconds is encoded
as a 64-bit two's complement integer. In both cases the value is followed with
32 bits for the nanosecond fraction. Note that the first two bits are RFU.

The data for text and binaries is prefixed with a varint byte size declaration.
Text is encoded as UTF-8. The flag is RFU.

Lists of floating points, text, binaries and data structures are prefixed with a
varint element size declaration. The flag is RFU.
