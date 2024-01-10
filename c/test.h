#include "Colfer.h"

#include <float.h>
#include <math.h>
#include <stdint.h>

const struct {
	const char* serial_hex;
	const struct seal_base_types values;
} golden_base_types[] = {

	// all zero
	{"088110"
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
		"00" // text size
		"00" // bool
		, {
			.s = { .utf8 = NULL }
		}
	},

	// small values
	{"508110"
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
		"09" // text size
		"01" // bool
		"c280" // payload U+0080
		"e0a080" // payload U+0800
		"f0908080" // payload U+10000
		, {
			.u8 = 1, .i8 = 2, .u16= 3, .i16 = 4,
			.u32 = 5, .i32 = 6, .u64 = 7, .i64 = 8,
			.f32 = 10, .f64 = 11,
			.t = { .tv_sec = 12, .tv_nsec = 13 },
			.s = { .utf8 = "\xC2\x80" "\xE0\xA0\x80" "\xF0\x90\x80\x80", .len = 9 },
			._flags = 1,
		}
	},

	// large values
	{"388210"
		"ff" // uint8
		"7f" // int8
		"04" // uint16
		"04" // int16
		"10" // uint32
		"10" // int32
		"00" // uint64
		"00" // int64
		"ffff7f7f" // float32
		"ffffffffffffef7f" // float64
		"ffc99afbffffffff" // timestamp
		"0a" // text size
		"00" // bool (has no large value)
		"ffff" // overflow 65535
		"feff" // overflow 32767
		"ffffffff" // overflow 4294967295
		"feffffff" // overflow 2147483647
		"ffffffffffffffff" // overflow 18446744073709551615
		"feffffffffffffff" // overflow 9223372036854775807
		"7f" // payload U+007F
		"dfbf" // payload U+07FF
		"efbfbf" // payload U+FFFF
		"f48fbfbf" // payload U+10FFFF
		, {
			.u8 = -1, .i8 = 127, .u16= -1, .i16 = 32767,
			.u32 = -1, .i32 = 2147483647, .u64 = -1LL, .i64 = 9223372036854775807LL,
			.f32 = FLT_MAX, .f64 = DBL_MAX,
			.t = { .tv_sec = (1L << 34) - 1, .tv_nsec = 999999999L },
			.s = { .utf8 = "\x7F" "\xDF\xBF" "\xEF\xBF\xBF" "\xF4\x8f\xBF\xBF", .len = 10 },
		}
	},

};
