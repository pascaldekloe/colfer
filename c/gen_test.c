#include "gen/Colfer.h"
#include "gen_test.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>

const unsigned char hex_table[] = "0123456789abcdef";

// hexstr maps data into buf as a null terminated hex string.
void hexstr(char* buf, const void* data, size_t datalen) {
	const uint8_t* p = data;
	for (; datalen-- != 0; p++) {
		uint8_t c = *p;
		*buf++ = hex_table[c >> 4];
		*buf++ = hex_table[c & 15];
	}
	*buf = 0;
}

int gen_o_equal(const gen_o* pa, const gen_o* pb) {
	if (pa == NULL || pb == NULL) return pa == pb;
	const gen_o a = *pa, b = *pb;

	if (! (
		a.b == b.b
		&& a.u8 == b.u8
		&& a.u16 == b.u16
		&& a.u32 == b.u32
		&& a.u64 == b.u64
		&& a.i32 == b.i32
		&& a.i64 == b.i64
		&& (a.f32 == b.f32 || (a.f32 != a.f32 && b.f32 != b.f32))
		&& a.f32s.len == b.f32s.len
		&& (a.f64 == b.f64 || (a.f64 != a.f64 && b.f64 != b.f64))
		&& a.f64s.len == b.f64s.len
		&& !memcmp(&a.t, &b.t, sizeof(struct timespec))
		&& a.s.len == b.s.len && !memcmp(a.s.utf8, b.s.utf8, a.s.len)
		&& a.ss.len == b.ss.len
		&& a.a.len == b.a.len && !memcmp(a.a.octets, b.a.octets, a.a.len)
		&& a.as.len == b.as.len
		&& gen_o_equal(a.o, b.o)
		&& a.os.len == b.os.len
	))
		return 0;

	for (size_t i = 0, n = a.f32s.len; i < n; ++i) {
		float fa = a.f32s.list[i], fb = b.f32s.list[i];
		if (fa != fb && (fa == fa || fb == fb)) return 0;
	}

	for (size_t i = 0, n = a.f64s.len; i < n; ++i) {
		double fa = a.f64s.list[i], fb = b.f64s.list[i];
		if (fa != fb && (fa == fa || fb == fb)) return 0;
	}

	for (size_t i = 0, n = a.ss.len; i < n; ++i) {
		colfer_text sa = a.ss.list[i], sb = b.ss.list[i];
		if (sa.len != sb.len || memcmp(sa.utf8, sb.utf8, sa.len)) return 0;
	}

	for (size_t i = 0, n = a.as.len; i < n; ++i) {
		colfer_binary ba = a.as.list[i], bb = b.as.list[i];
		if (ba.len != bb.len || memcmp(ba.octets, bb.octets, ba.len)) return 0;
	}

	for (size_t i = 0, n = a.os.len; i < n; ++i)
		if (!gen_o_equal(&a.os.list[i], &b.os.list[i])) return 0;

	return 1;
}

void gen_o_dump(const gen_o o) {
	char* buf = malloc(colfer_size_max * 2 + 1);

	printf("{ ");
	if (o.b) printf("b=true ");
	if (o.u8) printf("u8=%" PRIu8 " ", o.u8);
	if (o.u16) printf("u16=%" PRIu16 " ", o.u16);
	if (o.u32) printf("u32=%" PRIu32 " ", o.u32);
	if (o.i64) printf("i64=%" PRId64 " ", o.i64);
	if (o.i32) printf("i32=%" PRId32 " ", o.i32);
	if (o.i64) printf("i64=%" PRId64 " ", o.i64);
	if (o.f32) printf("f32=%f ", o.f32);
	if (o.f32s.len) {
		printf("f32s=[");
		for (size_t i = 0; i < o.f32s.len; ++i)
			printf(" %f", o.f32s.list[i]);
		printf(" ] ");
	}
	if (o.f64) printf("f64=%f ", o.f64);
	if (o.f64s.len) {
		printf("f64s=[");
		for (size_t i = 0; i < o.f64s.len; ++i)
			printf(" %f", o.f64s.list[i]);
		printf(" ] ");
	}
	if (o.t.tv_sec) printf("t.tv_sec=%lld ", (long long) o.t.tv_sec);
	if (o.t.tv_nsec) printf("t.tv_nsec=%ld ", o.t.tv_nsec);
	if (o.s.len) {
		hexstr(buf, o.s.utf8, o.s.len);
		printf("s=0x%s", buf);
	}
	if (o.ss.len) {
		printf("ss=[");
		for (size_t i = 0; i < o.ss.len; ++i) {
			hexstr(buf, o.ss.list[i].utf8, o.ss.list[i].len);
			printf(" 0x%s", buf);
		}
		printf(" ] ");
	}
	if (o.a.len) {
		hexstr(buf, o.a.octets, o.a.len);
		printf("a=0x%s", buf);
	}
	if (o.as.len) {
		printf("as=[");
		for (size_t i = 0; i < o.as.len; ++i) {
			hexstr(buf, o.as.list[i].octets, o.as.list[i].len);
			printf(" 0x%s", buf);
		}
		printf(" ] ");
	}
	if (o.o) {
		printf("o=");
		gen_o_dump(*o.o);
		printf(" ");
	}
	if (o.os.len) {
		printf("os=[");
		for (size_t i = 0; i < o.os.len; ++i) {
			putchar(' ');
			gen_o_dump(o.os.list[i]);
		}
		printf("] ");
	}
	putchar('}');

	free(buf);
}

int main() {
	const int n = sizeof(golden_cases) / sizeof(golden);
	printf("got %d golden cases\n", n);

	printf("TEST equality...\n");
	for (int i = 0; i < n; ++i) {
		const gen_o* a = &golden_cases[i].o;
		for (int j = 0; j < n; ++j) {
			const gen_o* b = &golden_cases[j].o;

			if (i == j) {
				if (!gen_o_equal(a, b))
					printf("0x%s: struct not equal to itself\n", golden_cases[i].hex);
			} else {
				if (gen_o_equal(a, b))
					printf("0x%s: struct equal to 0x%s\n", golden_cases[i].hex, golden_cases[j].hex);
			}
		}
	}

	printf("TEST marshal length...\n");
	for (int i = 0; i < n; ++i) {
		golden g = golden_cases[i];
		size_t got = gen_o_marshal_len(&g.o);
		size_t want = strlen(g.hex) / 2;
		if (got != want)
			printf("0x%s: got marshal length %zu, want %zu\n", g.hex, got, want);

		// size maximum
		for (colfer_size_max = 0; colfer_size_max < want; ++colfer_size_max) {
			got = gen_o_marshal_len(&g.o);
			if (got || errno != EFBIG)
				printf("0x%s: got marshal length %zu and errno %d with Colfer size maximum %zu\n", g.hex, got, errno, colfer_size_max);


			errno = 0;
		}
		colfer_size_max = 16 * 1024 * 1024;
	}

	void* buf = malloc(colfer_size_max);
	void* hex = malloc(colfer_size_max * 2 + 1);

	printf("TEST marshalling...\n");
	for (int i = 0; i < n; ++i) {
		golden g = golden_cases[i];
		size_t wrote = gen_o_marshal(&g.o, buf);
		hexstr(hex, buf, wrote);
		if (strcmp(hex, g.hex)) {
			printf("0x%s: got marshal data 0x%s\n", g.hex, hex);
			continue;
		}

		gen_o got = {0};
		size_t read = gen_o_unmarshal(&got, buf, wrote);
		if (errno != 0) {
			perror("unmarshal error");
			errno = 0;
		}
		if (read != wrote || !gen_o_equal(&got, &g.o)) {
			printf("0x%s: unmarshal read %zu bytes:\n\tgot: ", g.hex, read);
			gen_o_dump(got);
			printf("\n\twant: ");
			gen_o_dump(g.o);
			putchar('\n');
		}
	}

	printf("TEST unmarshal limits...\n");
	for (int i = 0; i < n; ++i) {
		golden g = golden_cases[i];
		size_t len = gen_o_marshal(&g.o, buf);
		if (!len) {
			printf("0x%s: skip due marshal fail\n", g.hex);
			continue;
		}

		// buffer length:
		for (size_t lim = 0; lim < len; lim++) {
			gen_o o = {0};
			size_t read = gen_o_unmarshal(&o, buf, lim);
			if (read || errno != EWOULDBLOCK)
				printf("0x%s[0:%zu]: unmarshal read %zu and errno %d\n", g.hex, lim, read, errno);

			errno = 0;
		}

		// size maximum:
		for (colfer_size_max = 0; colfer_size_max < len; ++colfer_size_max) {
			gen_o o = {0};
			size_t read = gen_o_unmarshal(&o, buf, len);
			if (read || errno != EFBIG)
				printf("0x%s: unmarshal read %zu with errno %d for size maximum %zu\n", g.hex, read, errno, colfer_size_max);

			errno = 0;
		}
		colfer_size_max = 16 * 1024 * 1024;
	}

	free(buf);
	free(hex);
}
