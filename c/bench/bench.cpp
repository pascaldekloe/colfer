#include "gen/Colfer.h"
#include "gen/scheme.pb.h"
#include "gen/scheme_generated.h"

#include <chrono>
#include <iostream>

const bench_colfer test_data[] = {
        {1234567890L, (char*) "db003lz12", 9, 389, 452, 0x488b5c2428488918ULL, 0.99, 1},
        {1234567891L, (char*) "localhost", 9, 22, 4096, 0x243048899c24c824ULL, 0.20, 0},
        {1234567892L, (char*) "kdc.local", 9, 88, 1984, 0x000048891c24485cULL, 0.06, 0},
        {1234567893L, (char*) "vhost8.dmz.example.com", 22, 27017, 59741, 0x5c2408488b9c2489ULL, 0.0, 1}
};

const size_t test_data_len = sizeof(test_data) / sizeof(bench_colfer);

// Rounds is the number of operations to run for each benchmark.
size_t rounds = 10000000;

int main(int argc, char **argv) {
	if (argc >= 2) {
		rounds = strtol(argv[1], NULL, 0);
	}

	void* serials[test_data_len];
	for (size_t i = 0; i < test_data_len; i++)
		serials[i] = malloc(colfer_size_max);
	bench_colfer o = {};

	auto start = std::chrono::high_resolution_clock::now();
	for (size_t i = 0; i < rounds; ++i) {
		size_t item = i % test_data_len;
		if (!bench_colfer_marshal(&test_data[item], serials[item])) {
			std::cout << "marshal error\n";
			return 1;
		}
	}
	auto end = std::chrono::high_resolution_clock::now();

	auto took = std::chrono::duration_cast<std::chrono::nanoseconds>(end - start).count();
	std::cout << "BENCH Colfer " << rounds << " marshals avg " << took / rounds << "ns\n";

	start = std::chrono::high_resolution_clock::now();
	for (size_t i = 0; i < rounds; ++i) {
		size_t item = i % test_data_len;
		if (!bench_colfer_unmarshal(&o, serials[item], colfer_size_max)) {
			std::cout << "unmarshal error\n";
			return 1;
		}
	}
	end = std::chrono::high_resolution_clock::now();

	took = std::chrono::duration_cast<std::chrono::nanoseconds>(end - start).count();
	std::cout << "BENCH Colfer " << rounds << " umarshals avg " << took / rounds << "ns\n";
}
