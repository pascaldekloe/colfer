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

// prevents compiler optimization:
void* serial;
size_t serial_size;

void marshal_colfer() {
	auto start = std::chrono::high_resolution_clock::now();
	for (size_t i = 0; i < rounds; ++i)
		serial_size = bench_colfer_marshal(&test_data[i % test_data_len], serial);
	auto end = std::chrono::high_resolution_clock::now();

	auto took = std::chrono::duration_cast<std::chrono::nanoseconds>(end - start).count();
	std::cout << "BENCH Colfer " << rounds << " marshals avg " << took / rounds << "ns\n";
}

void unmarshal_colfer() {
	void* serials[test_data_len];
	size_t serial_sizes[test_data_len];
	for (size_t i = 0; i < test_data_len; i++) {
		serials[i] = malloc(colfer_size_max);
		serial_sizes[i] = bench_colfer_marshal(&test_data[i], serials[i]);
	}

	auto o = new bench_colfer;

	auto start = std::chrono::high_resolution_clock::now();
	for (size_t i = 0; i < rounds; ++i)
		serial_size = bench_colfer_unmarshal(o, serials[i % test_data_len], colfer_size_max);
	auto end = std::chrono::high_resolution_clock::now();

	auto took = std::chrono::duration_cast<std::chrono::nanoseconds>(end - start).count();
	std::cout << "BENCH Colfer " << rounds << " umarshals avg " << took / rounds << "ns\n";
}

void marshal_fb() {
	flatbuffers::FlatBufferBuilder fbb;

	auto start = std::chrono::high_resolution_clock::now();
	for (size_t i = 0; i < rounds; ++i) {
		auto item = test_data[i % test_data_len];

		fbb.Clear();
		auto host = fbb.CreateString(item.host.utf8, item.host.len);
		auto o = bench::CreateFlatBuffers(fbb, item.key, host, item.port, item.size, item.hash, item.ratio, item.route);
		fbb.Finish(o);

		serial = fbb.GetBufferPointer();
		serial_size = fbb.GetSize();
		fbb.ReleaseBufferPointer();
	}
	auto end = std::chrono::high_resolution_clock::now();

	auto took = std::chrono::duration_cast<std::chrono::nanoseconds>(end - start).count();
	std::cout << "BENCH FlatBuffers " << rounds << " marshals avg " << took / rounds << "ns\n";
}

void unmarshal_fb() {
	flatbuffers::FlatBufferBuilder fbb;

	void* serials[test_data_len];
	size_t serial_sizes[test_data_len];
	for (size_t i = 0; i < test_data_len; ++i) {
		auto item = test_data[i % test_data_len];

		fbb.Clear();
		auto host = fbb.CreateString(item.host.utf8, item.host.len);
		auto o = bench::CreateFlatBuffers(fbb, item.key, host, item.port, item.size, item.hash, item.ratio, item.route);
		fbb.Finish(o);

		serial_sizes[i] = fbb.GetSize();
		serials[i] = malloc(fbb.GetSize());
		memcpy(serials[i], fbb.GetBufferPointer(), fbb.GetSize());
	}

	bench_colfer o = {};

	auto start = std::chrono::high_resolution_clock::now();
	for (size_t i = 0; i < rounds; ++i) {
		auto view = bench::GetFlatBuffers(serials[i % test_data_len]);
		o.key = view->key();
		auto s = view->host()->str();
		o.host.utf8 = &s[0];
		o.host.len = s.size();
		o.port = view->port();
		o.size = view->size();
		o.hash = view->hash();
		o.ratio = view->ratio();
		o.route = view->route();
	}
	auto end = std::chrono::high_resolution_clock::now();

	auto took = std::chrono::duration_cast<std::chrono::nanoseconds>(end - start).count();
	std::cout << "BENCH FlatBuffers " << rounds << " unmarshals avg " << took / rounds << "ns\n";
}

int main(int argc, char **argv) {
	if (argc >= 2) {
		rounds = strtol(argv[1], NULL, 0);
	}

	serial = malloc(colfer_size_max);

	marshal_colfer();
	unmarshal_colfer();
	marshal_fb();
	unmarshal_fb();
}
