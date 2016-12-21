#include <string>
#include <sstream>
#include <iostream>
#include <iomanip>

#include "./golden.h"

std::string hex_str(const unsigned char* data, int len) {
	assert(len < 1024);

	std::ostringstream buf;
	buf << std::hex;

	for (int i = 0; i < len; ++i)
		buf << std::setw(2) << std::setfill('0') << int(data[i]);

	return buf.str();
}

int main() {
	auto cases = new_golden_cases();
	std::cout << "got " << cases.size() << " cases\n";

	std::cout << "TEST equality operators...\n";
	for (auto pair1 : cases) {
		for (auto pair2 : cases) {
			if (pair1.first == pair2.first) {
				if (pair1.second != pair2.second)
					std::cout << "0x" << pair1.first << ": struct not equal to itself\n";
			} else {
				if (pair1.second == pair2.second)
					std::cout << "0x" << pair1.first << ": struct equals 0x" << pair2.first << "\n";
			}
		}
	}

	std::cout << "TEST marshal_len...\n";
	for (auto pair : cases) {
		auto got = pair.second.marshal_len();
		auto want = pair.first.size() / 2;
		if (got != want)
			std::cout << "0x" << pair.first << ": got marshal_len " << got << ", want " << want << "\n";
	}

	unsigned char* buf = new unsigned char[gen::colfer_size_max];

	std::cout << "TEST marshal...\n";
	for (auto pair : cases) {
		auto n = pair.second.marshal(buf);
		auto got = hex_str(&buf[0], n);
		if (got != pair.first) {
			std::cout << "0x" << pair.first << ": marshal wrote 0x" << got << "\n";
			continue;
		}
	}
}
