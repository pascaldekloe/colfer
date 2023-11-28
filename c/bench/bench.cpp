#include "Colfer.h"
#include "scheme.pb.h" // ProtoBuf
#include "scheme_generated.h" // FlatBuffers

#include <benchmark/benchmark.h> // https://github.com/google/benchmark

#include <iostream>
#include <cstdlib>


const bench_colfer test_data[] = {
	{1234567890L, (char*) "db003lz12", 9, 389, 452, 0x488b5c2428488918ULL, 0.99, 1},
	{1234567891L, (char*) "localhost", 9, 22, 4096, 0x243048899c24c824ULL, 0.20, 0},
	{1234567892L, (char*) "kdc.local", 9, 88, 1984, 0x000048891c24485cULL, 0.06, 0},
	{1234567893L, (char*) "vhost8.dmz.example.com", 22, 27017, 59741, 0x5c2408488b9c2489ULL, 0.0, 1}
};

const size_t test_data_len = sizeof(test_data) / sizeof(struct bench_colfer);


static void BM_marshal_colfer(benchmark::State& state) {
	void* buf = malloc(bench_colfer_marshal_max);

	for (int i = 0; state.KeepRunning(); i++) {
		auto data = &test_data[i % test_data_len];

		size_t n = bench_colfer_marshal(data, buf);
		if (n == 0) {
			state.SkipWithError("marshall error");
			return;
		}
		benchmark::DoNotOptimize(n);
		benchmark::DoNotOptimize(buf);
		benchmark::ClobberMemory();
	}
}

static void BM_unmarshal_colfer(benchmark::State& state) {
	void* serials[test_data_len];
	for (size_t i = 0; i < test_data_len; i++) {
		serials[i] = malloc(bench_colfer_marshal_max);
		size_t n = bench_colfer_marshal(&test_data[i], serials[i]);
		if (n == 0) {
			state.SkipWithError("marshall error");
			return;
		}
	}

	struct bench_colfer o;

	for (int i = 0; state.KeepRunning(); i++) {
		auto serial = serials[i % test_data_len];
		size_t n = bench_colfer_unmarshal(&o, serial);
		if (n == 0) {
			state.SkipWithError("unmarshal error");
			return;
		}
		benchmark::DoNotOptimize(n);
		benchmark::DoNotOptimize(o);
		benchmark::ClobberMemory();
		std::free((void*)o.host.utf8);
	}
}

static void BM_marshal_protobuf(benchmark::State& state) {
	std::string buf;

	bench::ProtoBuf data[test_data_len];
	for (size_t i = 0; i < test_data_len; i++) {
		data[i].set_key(test_data[i].key);
		data[i].set_host(test_data[i].host.utf8, test_data[i].host.len);
		data[i].set_port(test_data[i].port);
		data[i].set_size(test_data[i].size);
		data[i].set_hash(test_data[i].hash);
		data[i].set_ratio(test_data[i].ratio);
		data[i].set_route(test_data[i].bools & BENCH_COLFER_ROUTE_FLAG);
		if (!data[i].IsInitialized()) {
			state.SkipWithError("not initialized");
			return;
		}
	}

	for (int i = 0; state.KeepRunning(); i++) {
		auto ok = data[i % test_data_len].SerializeToString(&buf);
		if (!ok) {
			state.SkipWithError("not serialized");
			return;
		}

		benchmark::DoNotOptimize(ok);
		benchmark::DoNotOptimize(buf);
		benchmark::ClobberMemory();

		buf.clear();
	}
}

static void BM_unmarshal_protobuf(benchmark::State& state) {
	std::string serials[test_data_len];
	for (size_t i = 0; i < test_data_len; i++) {
		bench::ProtoBuf data;
		data.set_key(test_data[i].key);
		data.set_host(test_data[i].host.utf8, test_data[i].host.len);
		data.set_port(test_data[i].port);
		data.set_size(test_data[i].size);
		data.set_hash(test_data[i].hash);
		data.set_ratio(test_data[i].ratio);
		data.set_route(test_data[i].bools & BENCH_COLFER_ROUTE_FLAG);
		if (!data.IsInitialized()) {
			state.SkipWithError("not initialized");
			return;
		}
		if (!data.SerializeToString(&serials[i])) {
			state.SkipWithError("not serialized");
			return;
		}
	}


	bench::ProtoBuf o;

	for (int i = 0; state.KeepRunning(); i++) {
		auto ok = o.ParseFromString(serials[i % test_data_len]);
		if (!ok) {
			state.SkipWithError("not parsed");
			return;
		}

		benchmark::DoNotOptimize(ok);
		benchmark::DoNotOptimize(o);
		benchmark::ClobberMemory();
	}
}

static void BM_marshal_flatbuffers(benchmark::State& state) {
	flatbuffers::FlatBufferBuilder fbb(bench_colfer_marshal_max);

	for (int i = 0; state.KeepRunning(); i++) {
		auto data = test_data[i % test_data_len];

		auto host = fbb.CreateString(data.host.utf8, data.host.len);
		auto o = bench::CreateFlatBuffers(fbb, data.key, host, data.port, data.size, data.hash, data.ratio, data.bools & BENCH_COLFER_ROUTE_FLAG);
		fbb.Finish(o);

		benchmark::DoNotOptimize(fbb.GetBufferPointer());
		benchmark::DoNotOptimize(fbb.GetSize());
		benchmark::ClobberMemory();

		fbb.Clear();
	}
}

static void BM_unmarshal_flatbuffers(benchmark::State& state) {
	flatbuffers::FlatBufferBuilder fbb(bench_colfer_marshal_max);

	void* serials[test_data_len];
	for (size_t i = 0; i < test_data_len; ++i) {
		auto data = test_data[i % test_data_len];

		auto host = fbb.CreateString(data.host.utf8, data.host.len);
		auto o = bench::CreateFlatBuffers(fbb, data.key, host, data.port, data.size, data.hash, data.ratio, data.bools & BENCH_COLFER_ROUTE_FLAG);
		fbb.Finish(o);

		serials[i] = malloc(fbb.GetSize());
		memcpy(serials[i], fbb.GetBufferPointer(), fbb.GetSize());
		fbb.Clear();
	}

	struct bench_colfer o;

	for (int i = 0; state.KeepRunning(); i++) {
		auto serial = serials[i % test_data_len];

		auto view = bench::GetFlatBuffers(serial);
		// call each lazy-getter
		o.key = view->key();
		// allocate a copy:
		std::string* s = new std::string(view->host()->str());
		o.host.utf8 = s->c_str();
		o.host.len = s->size();
		o.port = view->port();
		o.size = view->size();
		o.hash = view->hash();
		o.ratio = view->ratio();
		if (view->route()) o.bools |= BENCH_COLFER_ROUTE_FLAG;

		benchmark::DoNotOptimize(o);
		benchmark::ClobberMemory();
		fbb.Clear();
		delete s;
	}
}


BENCHMARK(BM_marshal_colfer);
BENCHMARK(BM_unmarshal_colfer);
BENCHMARK(BM_marshal_protobuf);
BENCHMARK(BM_unmarshal_protobuf);
BENCHMARK(BM_marshal_flatbuffers);
BENCHMARK(BM_unmarshal_flatbuffers);

BENCHMARK_MAIN();
