#include "gen_test.h"

#include <inttypes.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>

// encode data into buf as a null terminated hex string
void hexstr(char* buf, const void* data, size_t datalen);

// field comparison
int gen_base_types_dfr(const struct gen_base_types a, const struct gen_base_types b);
int gen_list_types_dfr(const struct gen_list_types a, const struct gen_list_types b);

// print data in human-readable form
void gen_base_types_dump(const struct gen_base_types o);
void gen_list_types_dump(const struct gen_list_types o);

int main(void) {
	const int n = sizeof(golden_base_types) / sizeof(golden_base_types[0]);
	printf("got %d golden cases\n", n);

	int fail = 0;

	printf("TEST equality...\n");
	for (int i = 0; i < n; ++i) {
		for (int j = 0; j < n; ++j) {
			bool dfr = gen_base_types_dfr(golden_base_types[i].values, golden_base_types[j].values);
			if ((i == j) && (dfr != 0)) {
				fail++;
				printf("0x%s: compared not equal to self\n",
					golden_base_types[i].serial_hex);
			}
			if ((i != j) && (dfr == 0)) {
				fail++;
				printf("0x%s: compared equal to 0x%s\n",
					golden_base_types[i].serial_hex,
					golden_base_types[j].serial_hex);
			}
		}
	}

	char buf[COLFER_MAX];
	char hex[COLFER_MAX * 2 + 1];

	printf("TEST encoding roundtrip...\n");
	for (int i = 0; i < n; ++i) {
		size_t wrote = gen_base_types_marshal(&golden_base_types[i].values, buf);
		hexstr(hex, &buf, wrote);
		if (strcmp(hex, golden_base_types[i].serial_hex)) {
			fail++;
			printf("0x%s: got marshal data 0x%s\n",
				golden_base_types[i].serial_hex, hex);
			continue;
		}

		struct gen_base_types got;
		size_t read = gen_base_types_unmarshal(&got, buf);
		if (read == 0) {
			fail++;
			printf("0x%s: got unmarshal error\n",
				golden_base_types[i].serial_hex);
			continue;
		}
		if (read != wrote || gen_base_types_dfr(got, golden_base_types[i].values) != 0) {
			fail++;
			printf("0x%s: unmarshal read %zu bytes:\n\tgot: ",
				golden_base_types[i].serial_hex, read);
			gen_base_types_dump(got);
			printf("\n\twant: ");
			gen_base_types_dump(golden_base_types[i].values);
			putchar('\n');
		}
	}

	return fail;
}

const unsigned char hex_table[] = "0123456789abcdef";

void hexstr(char* buf, const void* data, size_t datalen) {
	const uint8_t* p = data;
	for (; datalen != 0; datalen--) {
		uint8_t c = *p++;
		*buf++ = hex_table[c >> 4];
		*buf++ = hex_table[c & 15];
	}
	*buf = 0;
}

int gen_base_types_dfr(const struct gen_base_types a, const struct gen_base_types b) {
	return a.bools != b.bools
	    || a.u8 != b.u8
	    || a.i8 != b.i8
	    || a.u16 != b.u16
	    || a.i16 != b.i16
	    || a.u32 != b.u32
	    || a.i32 != b.i32
	    || a.u64 != b.u64
	    || a.i64 != b.i64
	    || (a.f32 != b.f32 && a.f32 == a.f32 && b.f32 == b.f32)
	    || (a.f32 == b.f32 && (a.f32 != a.f32 || b.f32 != b.f32))
	    || (a.f64 != b.f64 && a.f64 == a.f64 && b.f64 == b.f64)
	    || (a.f64 == b.f64 && (a.f64 != a.f64 || b.f64 != b.f64))
	    || a.t.tv_sec != b.t.tv_sec || a.t.tv_nsec != b.t.tv_nsec
	    || a.s.len != b.s.len || memcmp(a.s.utf8, b.s.utf8, a.s.len)
	;
}

int gen_list_types_dfr(const struct gen_list_types a, const struct gen_list_types b) {
	bool dfr = a.f32s.len != b.f32s.len
	        || a.f64s.len != b.f64s.len
	        || a.ss.len != b.ss.len
	;

	for (size_t i = 0, n = a.f32s.len; i < n; ++i) {
		float fa = a.f32s.list[i];
		float fb = b.f32s.list[i];
		dfr |= fa == fa && fa != fb;
		dfr |= fa != fa && fb == fb;
	}
	for (size_t i = 0, n = a.f64s.len; i < n; ++i) {
		double fa = a.f64s.list[i];
		double fb = b.f64s.list[i];
		dfr |= fa == fa && fa != fb;
		dfr |= fa != fa && fb == fb;
	}

	for (size_t i = 0, n = a.ss.len; i < n; ++i) {
		size_t len = a.ss.list[i].len;
		dfr |= len != b.ss.list[i].len;
		dfr |= memcmp(a.ss.list[i].utf8, b.ss.list[i].utf8, len) != 0;
	}

	return dfr;
}

void gen_base_types_dump(const struct gen_base_types o) {
	char buf[1024];

	printf("{ ");
	printf("b=%d ", (o.bools & GEN_BASE_TYPES_B_FLAG));
	printf("i8=%" PRId8 " ", o.i8);
	printf("u8=%" PRIu8 " ", o.u8);
	printf("i16=%" PRId16 " ", o.i16);
	printf("u16=%" PRIu16 " ", o.u16);
	printf("i32=%" PRId32 " ", o.i32);
	printf("u32=%" PRIu32 " ", o.u32);
	printf("i64=%" PRId64 " ", o.i64);
	printf("u64=%" PRIu64 " ", o.u64);
	printf("f32=%f ", o.f32);
	printf("f64=%f ", o.f64);
	printf("t.tv_sec=%lld ", (long long) o.t.tv_sec);
	printf("t.tv_nsec=%lld ", (long long) o.t.tv_nsec);

	if (!o.s.len) {
		printf("s=0x ");
	} else if (o.s.len > sizeof(buf) / 2) {
		printf("s=%zuB ", o.s.len);
	} else {
		hexstr(buf, o.s.utf8, o.s.len);
		printf("s=0x%s ", buf);
	}

	putchar('}');
}

void gen_list_types_dump(const struct gen_list_types o) {
	char buf[1024];

	printf("{ ");
	printf("f32s=[");
	for (size_t i = 0; i < o.f32s.len; ++i)
		printf(" %f", o.f32s.list[i]);
	printf(" ] ");
	printf("f64s=[");
	for (size_t i = 0; i < o.f64s.len; ++i)
		printf(" %f", o.f64s.list[i]);
	printf(" ] ");

	printf("ss=[");
	for (size_t i = 0; i < o.ss.len; ++i) {
		if (!o.ss.list[i].len) {
			printf(" 0x ");
		} else if (o.ss.list[i].len > sizeof(buf) / 2) {
			printf(" %zuB ", o.ss.list[i].len);
		} else {
			hexstr(buf, o.ss.list[i].utf8, o.ss.list[i].len);
			printf(" 0x%s ", buf);
		}
	}
	printf(" ] ");

	putchar('}');
}
