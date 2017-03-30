# Colfer [![Build Status](https://travis-ci.org/pascaldekloe/colfer.svg?branch=master)](https://travis-ci.org/pascaldekloe/colfer)

Colfer is a binary serialization [format](https://github.com/pascaldekloe/colfer/wiki/Spec)
optimized for speed and size.

The project's compiler `colf(1)` generates source code from schema definitions
to marshal and unmarshall data structures.

This is free and unencumbered software released into the
[public domain](http://creativecommons.org/publicdomain/zero/1.0).
The format is inspired by Proto**col** Buf**fer**.


#### Language Support

* C, C99 compliant, C++11 compliant, WIP
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

* Rust and Python
* RMI (WIP)
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
    	list. The expression is applied to the target language under
    	the name ColferListMax. (default "64 * 1024")
  -p prefix
    	Adds a package prefix. Use slash as a separator when nesting.
  -s expression
    	Sets the default upper limit for serial byte sizes. The
    	expression is applied to the target language under the name
    	ColferSizeMax. (default "16 * 1024 * 1024")
  -v	Enables verbose reporting to the standard error.
  -x class
    	Makes all generated classes extend a super class. Use slash as
    	a package separator. Java only.

EXIT STATUS
	The command exits 0 on succes, 1 on compilation failure and 2
	when invoked without arguments.

EXAMPLES
	Compile ./io.colf with compact limits as C:

		colf -b src -s 2048 -l 96 C io.colf

	Compile ./api/*.colf in package com.example as Java:

		colf -p com/example -x com/example/Parent Java api

BUGS
	Report bugs at https://github.com/pascaldekloe/colfer/issues

SEE ALSO
	protoc(1)
```


It is recommended to commit the generated source code to the respective version
control. Modern developers may disagree and use the Maven plugin.

```xml
<plugin>
	<groupId>net.quies.colfer</groupId>
	<artifactId>colfer-maven-plugin</artifactId>
	<version>1.8</version>
	<configuration>
		<packagePrefix>com/example</packagePrefix>
	</configuration>
</plugin>
```



## Schema

Data structures are defined in `.colf` files. The format is quite conventional.

```
// Package demo offers a demonstration.
// These comment lines will end up in the generated code.
package demo

// Course is the grounds where the game of golf is played.
type course struct {
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
| bool		| char		| bool		| boolean	| Boolean	|
| uint8		| uint8_t	| uint8		| byte †	| Number	|
| uint16	| uint16_t	| uint16	| short †	| Number	|
| uint32	| uint32_t	| uint32	| int †		| Number	|
| uint64	| uint64_t	| uint64	| long †	| Number ‡	|
| int32		| int32_t	| int32		| int		| Number	|
| int64		| int64_t	| int64		| long		| Number ‡	|
| float32	| float		| float32	| float		| Number	|
| float64	| double	| float64	| double	| Number	|
| timestamp	| 2 × time_t	| Time ††	| Instant	| Date + Number	|
| text		| char* †‡	| string	| String ‡‡	| String ‡‡	|
| binary	| uint8_t* †‡ 	| []byte	| byte[]	| Uint8Array	|
| list		| †‡		| slice		| array		| Array		|

* † signed representation of unsigned data, i.e. may overflow to negative.
* ‡ range limited to (1 - 2⁵³, 2⁵³ - 1)
* †† timezone not preserved
* †‡ struct of pointer + size_t
* ‡‡ characters limited by UTF-16 (`U+0000`, `U+10FFFF`)

Lists may contain floating points, text, binaries or data structures.



## Compatibility

Name changes do not affect the serialization format. Deprecated fields should be
renamed to clearly discourage their use. For backwards compatibility new fields
must be added to the end of colfer structs. Thus the number of fields can be
seen as the schema version.



## Performance

Colfer aims to be the fastest and the smallest format while still resilient to malicious input. See the [benchmark wiki](https://github.com/pascaldekloe/colfer/wiki/Benchmark) for a comparison.
Suboptimal performance is treated like a bug.
