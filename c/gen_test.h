#include "Colfer.h"

#include <math.h>
#include <stdint.h>

const struct {
	const char* serial_hex;
	const struct gen_base_types values;
} golden_base_types[] = {

	// all zero
	{"1d0001"
		"00" // bool
		"00" // int8
		"00" // uint8
		"01" // int16
		"01" // uint16
		"01" // int32
		"01" // uint32
		"01" // int64
		"01" // uint64
		"00000000" // float32
		"0000000000000000" // float64
		"0000000000000000" // timestamp
		"01" // text
		, {
			.s = { .utf8 = NULL, .len = 0 },
		}
	},

	// small values
	{"1d0003" // fixed size 78, variable size 2
		"01" // bool
		"02" // int8
		"03" // uint8
		"11" // int16
		"0b" // uint16
		"19" // int32
		"0f" // uint32
		"21" // int64
		"13" // uint64
		"00002041" // float32
		"0000000000002640" // float64
		"0d00000003000000" // timestamp
		"03" // text
		// variable section (reversed order)
		"63" // text
		, {
			.bools = 1, .i8 = 2, .u8 = 3, .i16= 4, .u16 = 5,
			.i32 = 6, .u32 = 7, .i64 = 8, .u64 = 9,
			.f32 = 10, .f64 = 11, .t = { .tv_sec = 12, .tv_nsec = 13 },
			.s = { .utf8 = "c", .len = 1 },
		}
	},
};
