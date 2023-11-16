// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file test.colf for package gen.

#include <stdint.h>
#include <stdlib.h>
#include <time.h>

// Upper boundary for octets in a Colfer serial.
#define COLFER_MAX (16 * 1024 * 1024)

#ifdef __cplusplus
extern "C" {
#endif


// BaseTypes contains common data-types.
struct gen_base_types {
 	// B tests binary values.
	// GEN_BASE_TYPES_B_FLAG in bools bit-field; GEN_BASE_TYPES_B_FLAG;
 
 	// I8 tests signed integer values.
	int8_t i8;
 
 	// U8 tests unsigned integer values.
	uint8_t u8;
 
 	// I16 tests signed integer values.
	int16_t i16;
 
 	// U16 tests unsigned integer values.
	uint16_t u16;
 
 	// I32 tests signed integer values.
	int32_t i32;
 
 	// U32 tests unsigned integer values.
	uint32_t u32;
 
 	// I64 tests signed integer values.
	int64_t i64;
 
 	// U64 tests unsigned integer values.
	uint64_t u64;
 
 	// F32 tests floating-point values.
	float f32;
 
 	// F64 tests floating-point values.
	double f64;
 
 	// T tests timestamps (with nanosecond precision).
	struct timespec t;
 
 	// C tests Unicode strings of variable size.
	struct {
		const char* utf8;
		size_t len; // octet count
	} s;

	// Bit field for booleans.
	unsigned int bools;
};
// Boolean fields of struct gen_base_types.
#define GEN_BASE_TYPES_B_FLAG (1 << 0)

// Marshal encodes o as Colfer at start, up to COLFER_MAX octets in size. A zero
// return signals that the data in o exceeds COLFER_MAX.
size_t
gen_base_types_marshal(const struct gen_base_types* o, void* start);

// Unmarshal decodes o as Colfer from start. The number of octets consumed is at
// least 3, and at most COLFER_MAX. A zero return signals malformed data. String
// fields are allocated including null terminator. Caller owns the memory.
size_t
gen_base_types_unmarshal(struct gen_base_types* o, const void* start);


// ListTypes contains each BaseType supported in list form.
struct gen_list_types {
 	// F32s tests a variable-size listing.
	struct {
		float* list;
		size_t len; // element count
	} f32s;
 
 	// F64s tests a variable-size listing.
	struct {
		double* list;
		size_t len; // element count
	} f64s;
 
 	// Ss tests Unicode strings of variable size.
	struct {
		struct {
			const char* utf8;
			size_t len; // octet count
		}* list;
		size_t len; // element count
	} ss;
};

// Marshal encodes o as Colfer at start, up to COLFER_MAX octets in size. A zero
// return signals that the data in o exceeds COLFER_MAX.
size_t
gen_list_types_marshal(const struct gen_list_types* o, void* start);

// Unmarshal decodes o as Colfer from start. The number of octets consumed is at
// least 3, and at most COLFER_MAX. A zero return signals malformed data. String
// fields are allocated including null terminator. Caller owns the memory.
size_t
gen_list_types_unmarshal(struct gen_list_types* o, const void* start);


// ArrayTypes contains each BaseType supported in array form.
// The odd order is to breach some word boundaries in the fixed section.
struct gen_array_types {
 
	float f32a2[2];
 
 
	double f64a3[3];
 
 
	uint64_t u64a2[2];
 
 
	int32_t i32a2[2];
 
 
	uint32_t u32a2[2];
 
 
	int16_t i16a2[2];
 
 
	uint16_t u16a2[2];
 
 
	int8_t i8a2[2];
 
 
	uint8_t u8a2[2];
 
 
	struct timespec ta2[2];
 
 
	struct {
		const char* utf8;
		size_t len; // octet count
	} sa2[2];
};

// Marshal encodes o as Colfer at start, up to COLFER_MAX octets in size. A zero
// return signals that the data in o exceeds COLFER_MAX.
size_t
gen_array_types_marshal(const struct gen_array_types* o, void* start);

// Unmarshal decodes o as Colfer from start. The number of octets consumed is at
// least 3, and at most COLFER_MAX. A zero return signals malformed data. String
// fields are allocated including null terminator. Caller owns the memory.
size_t
gen_array_types_unmarshal(struct gen_array_types* o, const void* start);


// OpaqueTypes mixes fixed and variable-byte values.
struct gen_opaque_types {
 	// A8 tests 8-bit values.
	uint8_t a8;
 
 	// A16 tests 16-bit values.
	uint16_t a16;
 
 	// A32 tests 32-bit values.
	uint32_t a32;
 
 	// A64 tests 64-bit values.
	uint64_t a64;
};

// Marshal encodes o as Colfer at start, up to COLFER_MAX octets in size. A zero
// return signals that the data in o exceeds COLFER_MAX.
size_t
gen_opaque_types_marshal(const struct gen_opaque_types* o, void* start);

// Unmarshal decodes o as Colfer from start. The number of octets consumed is at
// least 3, and at most COLFER_MAX. A zero return signals malformed data. String
// fields are allocated including null terminator. Caller owns the memory.
size_t
gen_opaque_types_unmarshal(struct gen_opaque_types* o, const void* start);


// DromedaryCase mixes name conventions.
struct gen_dromedary_case {
 
	struct {
		const char* utf8;
		size_t len; // octet count
	} pascal_case;
 
 
	uint8_t with_snake;
};

// Marshal encodes o as Colfer at start, up to COLFER_MAX octets in size. A zero
// return signals that the data in o exceeds COLFER_MAX.
size_t
gen_dromedary_case_marshal(const struct gen_dromedary_case* o, void* start);

// Unmarshal decodes o as Colfer from start. The number of octets consumed is at
// least 3, and at most COLFER_MAX. A zero return signals malformed data. String
// fields are allocated including null terminator. Caller owns the memory.
size_t
gen_dromedary_case_unmarshal(struct gen_dromedary_case* o, const void* start);


#ifdef __cplusplus
} // extern "C"
#endif
