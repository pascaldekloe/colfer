// when integer fields
#include <cstdint>
// when timestamp fields
#include <chrono>
// when list or binary fields
#include <vector>
// when text fields
#include <string>

// O contains all supported data types.
class O {
	public:
		// B tests booleans.
		bool b;

		// U8 tests unsigned 8-bit integers.
		uint8_t u8;

		// U16 tests unsigned 16-bit integers.
		uint16_t u16;

		// U32 tests unsigned 32-bit integers.
		uint32_t u32;

		// U64 tests unsigned 64-bit integers.
		uint64_t u64;

		// I32 tests signed 32-bit integers.
		int32_t i32;

		// I64 tests signed 64-bit integers.
		int64_t i64;

		// F32 tests 32-bit floating points.
		float f32;

		// F32s tests 32-bit floating point lists.
		std::vector<float> f32s;

		// F64 tests 64-bit floating points.
		double f64;

		// F64s tests 64-bit floating point lists.
		std::vector<double> f64s;

		// T tests timestamps.
		std::chrono::nanoseconds t;

		// A tests binaries.
		std::vector<uint8_t> a;

		// As tests binary lists.
		std::vector<std::vector<uint8_t> > as;

		// S tests text.
		std::string s;

		// Ss tests text lists.
		std::vector<std::string> ss;

		// O tests nested data structures.
		O* o;

		// Os tests data structure lists.
		std::vector<O*> os;

		// MarshalLen returns the Colfer serial byte size.
		size_t marshalLen();

		// MarshalTo encodes O as Colfer into buf and returns the number of bytes
		// written.
		size_t marshalTo(void* buf);

		// Unmarshal decodes data as Colfer and returns the number of bytes read.
		size_t umarshal(void* data);
};

int main() {
}
