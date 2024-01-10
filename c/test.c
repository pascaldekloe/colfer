#include "test.h"

#include <inttypes.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// encode data into buf as a null terminated hex string
void hexstr(char* buf, const void* data, size_t datalen);

// field comparison
int seal_base_types_dfr(const struct seal_base_types a, const struct seal_base_types b);
int seal_list_types_dfr(const struct seal_list_types a, const struct seal_list_types b);

// print data in human-readable form
void seal_base_types_dump(const struct seal_base_types o);
void seal_list_types_dump(const struct seal_list_types o);

int main(void) {
	const int n = sizeof(golden_base_types) / sizeof(golden_base_types[0]);
	printf("got %d golden cases\n", n);

	int fail = 0;

	printf("TEST equality...\n");
	for (int i = 0; i < n; ++i) {
		for (int j = 0; j < n; ++j) {
			bool dfr = seal_base_types_dfr(golden_base_types[i].values, golden_base_types[j].values);
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

	char buf[4096];
	char hex[4096 * 2 + 1];

	printf("TEST encoding roundtrip...\n");
	for (int i = 0; i < n; ++i) {
		size_t wrote = seal_base_types_marshal(&golden_base_types[i].values, buf);
		hexstr(hex, &buf, wrote);
		if (strcmp(hex, golden_base_types[i].serial_hex)) {
			fail++;
			printf("0x%s: got marshal data 0x%s\n",
				golden_base_types[i].serial_hex, hex);
			continue;
		}

		struct seal_base_types got;
		size_t read = seal_base_types_unmarshal(&got, buf, &malloc);
		if (read == 0) {
			fail++;
			printf("0x%s: got unmarshal error\n",
				golden_base_types[i].serial_hex);
			continue;
		}
		if (read != wrote || seal_base_types_dfr(got, golden_base_types[i].values) != 0) {
			fail++;
			printf("0x%s: unmarshal read %zu bytes:\n\tgot: ",
				golden_base_types[i].serial_hex, read);
			seal_base_types_dump(got);
			printf("\n\twant: ");
			seal_base_types_dump(golden_base_types[i].values);
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

int seal_base_types_dfr(const struct seal_base_types a, const struct seal_base_types b) {
	return a._flags != b._flags
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

int seal_list_types_dfr(const struct seal_list_types a, const struct seal_list_types b) {
	bool dfr = a.f32l.len != b.f32l.len
	        || a.f64l.len != b.f64l.len
	        || a.sl.len != b.sl.len
	;

	for (size_t i = 0, n = a.f32l.len; i < n; ++i) {
		float fa = a.f32l.list[i];
		float fb = b.f32l.list[i];
		dfr |= fa == fa && fa != fb;
		dfr |= fa != fa && fb == fb;
	}
	for (size_t i = 0, n = a.f64l.len; i < n; ++i) {
		double fa = a.f64l.list[i];
		double fb = b.f64l.list[i];
		dfr |= fa == fa && fa != fb;
		dfr |= fa != fa && fb == fb;
	}

	for (size_t i = 0, n = a.sl.len; i < n; ++i) {
		size_t len = a.sl.list[i].len;
		dfr |= len != b.sl.list[i].len;
		dfr |= memcmp(a.sl.list[i].utf8, b.sl.list[i].utf8, len) != 0;
	}

	return dfr;
}

void seal_base_types_dump(const struct seal_base_types o) {
	char buf[1024];

	printf("{ ");

	printf("u8=%" PRIu8 " ", o.u8);
	printf("i8=%" PRId8 " ", o.i8);
	printf("u16=%" PRIu16 " ", o.u16);
	printf("i16=%" PRId16 " ", o.i16);
	printf("u32=%" PRIu32 " ", o.u32);
	printf("i32=%" PRId32 " ", o.i32);
	printf("u64=%" PRIu64 " ", o.u64);
	printf("i64=%" PRId64 " ", o.i64);

	printf("f32=%f ", o.f32);
	printf("f64=%f ", o.f64);

	printf("t=%lld.%09ld ", (long long) o.t.tv_sec, (long)o.t.tv_nsec);

	if (!o.s.len) {
		printf("s=0x ");
	} else if (o.s.len > sizeof(buf) / 2) {
		printf("s=%zu B ", o.s.len);
	} else {
		hexstr(buf, o.s.utf8, o.s.len);
		printf("s=0x%s ", buf);
	}

	printf("b=%d ", (o._flags & SEAL_BASE_TYPES_B_FLAG));

	putchar('}');
}

void seal_list_types_dump(const struct seal_list_types o) {
	char buf[1024];

	printf("{ ");

	printf("a8l=[");
	for (size_t i = 0; i < o.a8l.len; ++i)
		printf(" %x", o.a8l.list[i]);
	printf(" ] ");
	printf("a16l=[");
	for (size_t i = 0; i < o.a16l.len; ++i)
		printf(" %x", o.a16l.list[i]);
	printf(" ] ");
	printf("a32l=[");
	for (size_t i = 0; i < o.a32l.len; ++i)
		printf(" %x", o.a32l.list[i]);
	printf(" ] ");
	printf("a64l=[");
	for (size_t i = 0; i < o.a64l.len; ++i)
		printf(" %llx", o.a64l.list[i]);
	printf(" ] ");

	printf("f32l=[");
	for (size_t i = 0; i < o.f32l.len; ++i)
		printf(" %f", o.f32l.list[i]);
	printf(" ] ");
	printf("f64l=[");
	for (size_t i = 0; i < o.f64l.len; ++i)
		printf(" %f", o.f64l.list[i]);
	printf(" ] ");

	printf("tl=[");
	for (size_t i = 0; i < o.tl.len; ++i)
		printf(" %lld.%09ld s", (long long)o.tl.list[i].tv_sec, (long)o.tl.list[i].tv_nsec);
	printf(" ] ");

	printf("sl=[");
	for (size_t i = 0; i < o.sl.len; ++i) {
		if (!o.sl.list[i].len) {
			printf(" 0x ");
		} else if (o.sl.list[i].len > sizeof(buf) / 2) {
			printf(" %zuB ", o.sl.list[i].len);
		} else {
			hexstr(buf, o.sl.list[i].utf8, o.sl.list[i].len);
			printf(" 0x%s ", buf);
		}
	}
	printf(" ] ");

	putchar('}');
}
