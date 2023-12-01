# Colfer

Colfer is a binary serialization [format](https://github.com/pascaldekloe/colfer/blob/v2/SPEC.md)
optimized for speed and size.

The project's compiler `colf(1)` generates source code from schema definitions
to marshal and unmarshall data structures.

This is free and unencumbered software released into the
[public domain](http://creativecommons.org/publicdomain/zero/1.0).
The format is inspired by Proto**col** Buf**fer**s.

[![CI](https://github.com/pascaldekloe/colfer/actions/workflows/ci.yml/badge.svg)](https://github.com/pascaldekloe/colfer/actions/workflows/ci.yml)


#### Language Support

* **C**, ISO/IEC 9899:2011 compliant a.k.a. C11, C++ compatible
* **Go**, a.k.a. golang
* **Java**, base module exclusively, Android compatible
* **JavaScript**, a.k.a. ECMAScript, NodeJS compatible
* 🚧 Gergely Bódi realised a functional **Dart** [port](https://github.com/vendelin8/colfer).
* 🚧 Karthik Kumar Viswanathan has a **Python** [alternative](https://github.com/guilt/colfer-python) under construction.


#### Features

* Simple and straightforward in use
* No dependencies other than the core library
* Both faster and smaller than the competition
* [Robust](#security) against malicious input
* Maximum of 127 fields per data structure
* No support for enumerations
* Framed; suitable for concatenation/streaming


## Use

Download a [prebuilt compiler](https://github.com/pascaldekloe/colfer/releases)
or run `go get -u github.com/pascaldekloe/colfer/cmd/colf` to make one yourself.
Homebrew users can also `brew install colfer`.

The command prints its own manual when invoked without arguments.

```
NAME
	colf — compile Colfer schemas

SYNOPSIS
	colf [-h]
	colf [-vf] [-b directory] [-p package] \
		[-s expression] [-l expression] C [file ...]
	colf [-vf] [-b directory] [-p package] [-t files] \
		[-s expression] [-l expression] Go [file ...]
	colf [-vf] [-b directory] [-p package] [-t files] \
		[-x class] [-i interfaces] [-c file] \
		[-s expression] [-l expression] Java [file ...]
	colf [-vf] [-b directory] [-p package] \
		[-s expression] [-l expression] JavaScript [file ...]

DESCRIPTION
	The output is source code for either C, Go, Java or JavaScript.

	For each operand that names a file of a type other than
	directory, colf reads the content as schema input. For each
	named directory, colf reads all files with a .colf extension
	within that directory. If no operands are given, the contents of
	the current directory are used.

	A package definition may be spread over several schema files.
	The directory hierarchy of the input is not relevant to the
	generated code.

OPTIONS
  -b directory
    	Use a base directory for the generated code. (default ".")
  -c file
    	Insert a code snippet from a file.
  -f	Normalize the format of all schema input on the fly.
  -h	Prints the manual to standard error.
  -i interfaces
    	Make all generated classes implement one or more interfaces.
    	Use commas as a list separator.
  -l expression
    	Set the default upper limit for the number of elements in a
    	list. The expression is applied to the target language under
    	the name ColferListMax. (default "64 * 1024")
  -p package
    	Compile to a package prefix.
  -s expression
    	Set the default upper limit for serial byte sizes. The
    	expression is applied to the target language under the name
    	ColferSizeMax. (default "16 * 1024 * 1024")
  -t files
    	Supply custom tags with one or more files. Use commas as a list
    	separator. See the TAGS section for details.
  -v	Enable verbose reporting to standard error.
  -x class
    	Make all generated classes extend a super class.

TAGS
	Tags, a.k.a. annotations, are source code additions for structs
	and/or fields. Input for the compiler can be specified with the
	-t option. The data format is line-oriented.

		<line> :≡ <qual> <space> <code> ;
		<qual> :≡ <package> '.' <dest> ;
		<dest> :≡ <struct> | <struct> '.' <field> ;

	Lines starting with a '#' are ignored (as comments). Java output
	can take multiple tag lines for the same struct or field. Each
	code line is applied in order of appearance.

EXIT STATUS
	The command exits 0 on success, 1 on error and 2 when invoked
	without arguments.

EXAMPLES
	Compile ./io.colf with compact limits as C:

		colf -b src -s 2048 -l 96 C io.colf

	Compile ./*.colf with a common parent as Java:

		colf -p com.example.model -x com.example.io.IOBean Java

BUGS
	Report bugs at <https://github.com/pascaldekloe/colfer/issues>.

	Text validation is not part of the marshalling and unmarshalling
	process. C and Go just pass any malformed UTF-8 characters. Java
	and JavaScript replace unmappable content with the '?' character
	(ASCII 63).

SEE ALSO
	protoc(1), flatc(1)
```

It is recommended to commit the generated source code into the respective
version control to preserve build consistency and minimise the need for compiler
installations. Alternatively, you may use the
[Maven plugin](https://github.com/pascaldekloe/colfer/wiki/Java#maven).

```xml
<plugin>
	<groupId>net.quies.colfer</groupId>
	<artifactId>colfer-maven-plugin</artifactId>
	<version>1.11.2</version>
	<configuration>
		<packagePrefix>com/example</packagePrefix>
	</configuration>
</plugin>
```



## Schema

Data structures are defined in `.colf` files. The format is quite self-explanatory.

```go
// Package demo offers a demonstration.
// These comment lines end up in the generated code.
package demo

// Course is the grounds where the game of golf is played.
type course struct {
	ID    uint64
	name  text
	holes [18]hole
	image []opaque8
	tags  []text
}

type hole struct {
	// Lat is the latitude of the cup.
	lat float64
	// Lon is the longitude of the cup.
	lon float64
	// Par is the stroke indicator.
	par uint8
	// Water marks the presence of water.
	water bool
	// Sand marks the presence of sand.
	sand bool
}
```

See what the generated code looks like in
[C](https://gist.github.com/pascaldekloe/05e903f12a4f02a995f71d0c18872b65),
[Go](https://gist.github.com/pascaldekloe/786fd46e6e4710c14fee7da1f480c2d4),
[Java](https://gist.github.com/pascaldekloe/b54326e6b7c5e9f036911a8cbea6ccbf) or
[JavaScript](https://gist.github.com/pascaldekloe/5653c8bb074ebd29ffcc0deece7495a4).

The following table shows the default data-types per programming language.

| Colfer	| C		| Go		| Java		| Rust		|
|:--------------|:--------------|:--------------|:--------------|:--------------|
| bool		| int + mask	| bool		| boolean	| bool		|
| int8		| int8_t	| int8		| byte		| i16		|
| uint8		| uint8_t	| uint8		| byte [^1]	| u16		|
| opaque8	| uint8_t	| uint8		| byte [^1]	| u16		|
| int16		| int16_t	| int16		| short		| i16		|
| uint16	| uint16_t	| uint16	| short [^1]	| u16		|
| opaque16	| uint16_t	| uint16	| short [^1]	| u16		|
| int32		| int32_t	| int32		| int		| i32		|
| uint32	| uint32_t	| uint32	| int [^1]	| u32		|
| opaque32	| uint32_t	| uint32	| int [^1]	| u32		|
| int64		| int64_t	| int64		| long		| i64		|
| uint64	| uint64_t	| uint64	| long [^1]	| u64		|
| opaque64	| uint64_t	| uint64	| long [^1]	| u64		|
| float32	| float		| float32	| float		| f32		|
| float64	| double	| float64	| double	| f64		|
| timestamp	| timespec	| Time [^2]	| Instant	| DateTime<Utc>	|
| opaque	| void*		| []byte	| byte[] [^1]	| [u8; usize]	|
| text		| const char*	| string	| String	| String	|
| []T		| T*		| []T		| T[]		| [T]		|
| [n]T		| T[n]		| [n]T		| final T[]	| [T; n]	|

[^1]: may overflow to negative values
[^2]: timezone not preserved


## Security

Colfer is suited for untrusted data sources such as network I/O or bulk streams.
Marshalling and unmarshalling comes with built-in size protection to ensure
predictable memory consumption. The format prevents memory bombs by design.

The marshaller may not produce malformed output, regardless of the data input.
In no event may the unmarshaller read outside the boundaries of a serial. Fuzz
testing did not reveal any volnurabilities yet. Computing power is welcome.


## Security

Colfer is suited for untrusted data sources such as network I/O or bulk streams.
Marshalling and unmarshalling comes with built-in size protection to ensure
predictable memory consumption. The format prevents memory bombs by design.

The marshaller may not produce malformed output, regardless of the data input.
In no event may the unmarshaller read outside the boundaries of a serial. Fuzz
testing did not reveal any volnurabilities yet. Computing power is welcome.


## Compatibility

Name changes do not affect the serialization format. Deprecated fields should be
renamed to clearly discourage their use. For backwards compatibility new fields
must be added to the end of colfer structs. Thus the number of fields can be
seen as the schema version.

Serials from Colfer version 1 always start with an octet in range [0x00, 0x7f].
Any (custom) prefix in range [0x80, 0xff] on version 2 serials will allow for a
distriction between the two when mixing both formats.


## Performance

Colfer aims to be the fastest and the smallest format without compromising on
reliability. See the
[benchmark wiki](https://github.com/pascaldekloe/colfer/wiki/Benchmark) for a
comparison. Suboptimal performance is treated like a bug.
