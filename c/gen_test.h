#include "Colfer.h"

#include <math.h>
#include <stdint.h>

const struct {
	const char* serial_hex;
	const struct gen_base_types values;
} golden_base_types[] = {

	// all zero
	{"211002"
		"00" // uint8
		"00" // int8
		"01" // uint16
		"01" // int16
		"01" // uint32
		"01" // int32
		"01" // uint64
		"01" // int64
		"00000000" // float32
		"0000000000000000" // float64
		"0000000000000000" // timestamp
		"01" // text size
		"00" // bool
		, {
			.s = { .utf8 = NULL, .len = 0 },
		}
	},

	// small values
	{"221002"
		"01" // uint8
		"02" // int8
		"07" // uint16
		"11" // int16
		"0b" // uint32
		"19" // int32
		"0f" // uint64
		"21" // int64
		"00002041" // float32
		"0000000000002640" // float64
		"0d00000003000000" // timestamp
		"03" // text size
		"01" // bool
		"63" // text payload
		, {
			.u8 = 1, .i8 = 2, .u16= 3, .i16 = 4,
			.u32 = 5, .i32 = 6, .u64 = 7, .i64 = 8,
			.f32 = 10, .f64 = 11,
			.t = { .tv_sec = 12, .tv_nsec = 13 },
			.s = { .utf8 = "c", .len = 1 },
			.bools = 1,
		}
	},
};
