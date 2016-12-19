#include <cassert>
#include <cstdint>
#include <limits>

#include <chrono>
#include <string>

#include <map>
#include <vector>

#include "gen/Colfer.h"

std::map<std::string, gen::O> new_golden_cases() {
	std::map<std::string, gen::O> m = {
		{"7f", {}},
		{"007f", {.b = true}},
		{"01017f", {.u32 = 1}},
		{"01ff017f", {.u32 = UINT8_MAX}},
		{"01ffff037f", {.u32 = UINT16_MAX}},
		{"81ffffffff7f", {.u32 = UINT32_MAX}},
		{"02017f", {.u64 = 1}},
		{"02ff017f", {.u64 = UINT8_MAX}},
		{"02ffff037f", {.u64 = UINT16_MAX}},
		{"02ffffffff0f7f", {.u64 = UINT32_MAX}},
		{"82ffffffffffffffff7f", {.u64 = UINT64_MAX}},
		{"03017f", {.i32 = 1}},
		{"83017f", {.i32 = -1}},
		{"037f7f", {.i32 = INT8_MAX}},
		{"8380017f", {.i32 = INT8_MIN}},
		{"03ffff017f", {.i32 = INT16_MAX}},
		{"838080027f", {.i32 = INT16_MIN}},
		{"03ffffffff077f", {.i32 = INT32_MAX}},
		{"8380808080087f", {.i32 = INT32_MIN}},
		{"04017f", {.i64 = 1}},
		{"84017f", {.i64 = -1}},
		{"047f7f", {.i64 = INT8_MAX}},
		{"8480017f", {.i64 = INT8_MIN}},
		{"04ffff017f", {.i64 = INT16_MAX}},
		{"848080027f", {.i64 = INT16_MIN}},
		{"04ffffffff077f", {.i64 = INT32_MAX}},
		{"8480808080087f", {.i64 = INT32_MIN}},
		{"04ffffffffffffffff7f7f", {.i64 = INT64_MAX}},
		{"848080808080808080807f", {.i64 = INT64_MIN}},
		{"05000000017f", {.f32 = std::numeric_limits<float>::denorm_min()}},
		{"057f7fffff7f", {.f32 = std::numeric_limits<float>::max()}},
		{"057fc000007f", {.f32 = std::numeric_limits<float>::quiet_NaN()}},
		{"0600000000000000017f", {.f64 = std::numeric_limits<double>::denorm_min()}},
		{"067fefffffffffffff7f", {.f64 = std::numeric_limits<double>::max()}},
		{"067ff80000000000017f", {.f64 = std::numeric_limits<double>::quiet_NaN()}},
		{"0755ef312a2e5da4e77f", {.t = std::chrono::nanoseconds(1441739050777888999)}},
		{"870000000100000000000000007f", {.t = std::chrono::seconds(UINT32_MAX) + std::chrono::seconds(1)}},
		{"87ffffffffffffffff2e5da4e77f", {.t = std::chrono::nanoseconds(222111001)}},
		{"87fffffff14f443f00000000007f", {.t = std::chrono::seconds(-63094636800)}},
		{"0801417f", {.s = "A"}},
		{"080261007f", {.s = "a\x00"}},
		{"0809c280e0a080f09080807f", {.s = "\u0080\u0800\U00010000"}},
		{"08800120202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020207f", {.s = std::string(128, ' ')}},
		{"0901ff7f", {.a = {UINT8_MAX}}},
		{"090202007f", {.a = {2, 0}}},
		{"09c0010909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909097f", {.a = std::vector<uint8_t>(192, 9)}},
		{"0a7f7f", {.o = new gen::O}},
		{"0a007f7f", {.o = new gen::O}},
		{"0b01007f7f", {.os = {{.b = true}}}},
		{"0b027f7f7f", {.os = {{}, {}}}},
		{"0c0300016101627f", {.ss = {"", "a", "b"}}},
		{"0d0201000201027f", {.as = {{0}, {1, 2}}}},
		{"0e017f", {.u8 = 1}},
		{"0eff7f", {.u8 = UINT8_MAX}},
		{"8f017f", {.u16 = 1}},
		{"0fffff7f", {.u16 = UINT16_MAX}},
		{"1002000000003f8000007f", {.f32s = {1}}},
		{"11014058c000000000007f", {.f64s = {99}}}
	};

	// FIXME: do following directry as a literal (how?)
	auto p = m.find("0a007f7f");
	assert(p != m.end());
	assert(p->second.o);
	assert(! p->second.o->o);
	gen::O o = {.b = true};
	p->second.o->o = &o;

	return m;
}
