#include "build/gen/Colfer.h"
#include "build/gen/scheme.pb.h"
#include "build/gen/scheme_generated.h"

// https://github.com/google/benchmark
#include <benchmark/benchmark.h>


const bench_colfer test_data[] = {
        {1234567890L, (char*) "db003lz12", 9, 389, 452, 0x488b5c2428488918ULL, 0.99, 1},
        {1234567891L, (char*) "localhost", 9, 22, 4096, 0x243048899c24c824ULL, 0.20, 0},
        {1234567892L, (char*) "kdc.local", 9, 88, 1984, 0x000048891c24485cULL, 0.06, 0},
        {1234567893L, (char*) "vhost8.dmz.example.com", 22, 27017, 59741, 0x5c2408488b9c2489ULL, 0.0, 1}
};

const size_t test_data_len = sizeof(test_data) / sizeof(bench_colfer);


static void BM_marshal_colfer(benchmark::State& state) {
	void* buf = malloc(colfer_size_max);

	for (int i = 0; state.KeepRunning(); i++) {
		auto data = &test_data[i % test_data_len];

                benchmark::DoNotOptimize(bench_colfer_marshal(data, buf));
                benchmark::DoNotOptimize(buf);
                benchmark::ClobberMemory();
        }
}

static void BM_unmarshal_colfer(benchmark::State& state) {
	void* serials[test_data_len];
	for (size_t i = 0; i < test_data_len; i++) {
		serials[i] = malloc(colfer_size_max);
		bench_colfer_marshal(&test_data[i], serials[i]);
	}

	auto o = new bench_colfer;

	for (int i = 0; state.KeepRunning(); i++) {
		auto serial = serials[i % test_data_len];

                benchmark::DoNotOptimize(bench_colfer_unmarshal(o, serial, colfer_size_max));
                benchmark::DoNotOptimize(o);
                benchmark::ClobberMemory();
        }
}

static void BM_marshal_flatbuffers(benchmark::State& state) {
	flatbuffers::FlatBufferBuilder fbb(colfer_size_max);

	for (int i = 0; state.KeepRunning(); i++) {
		auto data = test_data[i % test_data_len];

		auto host = fbb.CreateString(data.host.utf8, data.host.len);
		auto o = bench::CreateFlatBuffers(fbb, data.key, host, data.port, data.size, data.hash, data.ratio, data.route);
		fbb.Finish(o);

                benchmark::DoNotOptimize(fbb.GetBufferPointer());
                benchmark::DoNotOptimize(fbb.GetSize());
                benchmark::ClobberMemory();

		fbb.Clear();
	}
}

static void BM_unmarshal_flatbuffers(benchmark::State& state) {
	flatbuffers::FlatBufferBuilder fbb(colfer_size_max);

	void* serials[test_data_len];
	for (size_t i = 0; i < test_data_len; ++i) {
		auto data = test_data[i % test_data_len];

		auto host = fbb.CreateString(data.host.utf8, data.host.len);
		auto o = bench::CreateFlatBuffers(fbb, data.key, host, data.port, data.size, data.hash, data.ratio, data.route);
		fbb.Finish(o);

		serials[i] = malloc(fbb.GetSize());
		memcpy(serials[i], fbb.GetBufferPointer(), fbb.GetSize());
		fbb.Clear();
	}

	bench_colfer o = {};

	for (int i = 0; state.KeepRunning(); i++) {
		auto serial = serials[i % test_data_len];

		auto view = bench::GetFlatBuffers(serial);
		o.key = view->key();
		auto s = view->host()->str();
		o.host.utf8 = &s[0];
		o.host.len = s.size();
		o.port = view->port();
		o.size = view->size();
		o.hash = view->hash();
		o.ratio = view->ratio();
		o.route = view->route();

                benchmark::DoNotOptimize(o);
                benchmark::ClobberMemory();

		fbb.Clear();
	}
}


BENCHMARK(BM_marshal_colfer);
BENCHMARK(BM_unmarshal_colfer);
BENCHMARK(BM_marshal_flatbuffers);
BENCHMARK(BM_unmarshal_flatbuffers);

BENCHMARK_MAIN();
