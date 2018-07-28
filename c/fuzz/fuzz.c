#include "../gen/Colfer.h"

#include <stdlib.h>
#include <unistd.h>
#include <stdio.h>
#include <string.h>

int main() {
	void* in = malloc(colfer_size_max);
	ssize_t inlen = read(STDIN_FILENO, in, colfer_size_max);
	if (inlen < 0) {
		perror(NULL);
		return 1;
	}

	gen_o o = {0};
	size_t read = gen_o_unmarshal(&o, in, inlen);
	if (!read) {
		return 0;
	}

	size_t len = gen_o_marshal_len(&o);
	if (len != read) {
		return 2;
	}

	void* out = malloc(colfer_size_max);
	size_t wrote = gen_o_marshal(&o, out);
	if (wrote != read) {
		return 3;
	}

	if (memcmp(in, out, len)) return 4;

	return 0;
}
