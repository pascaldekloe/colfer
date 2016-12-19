#include <iostream>

#include "./golden.h"

int main() {
	auto cases = new_golden_cases();
	std::cout << "got " << cases.size() << " cases\n";

	std::cout << "TEST equality operators...\n";
	for (auto pair1 : cases) {
		for (auto pair2 : cases) {
			if (pair1.first == pair2.first) {
				if (pair1.second != pair2.second)
					std::cout << "0x" << pair1.first << ": not equal to itself!\n";
			} else {
				if (pair1.second == pair2.second)
					std::cout << "0x" << pair1.first << ": equals 0x" << pair2.first << "!\n";
			}
		}
	}

	std::cout << "TEST marshal_len...\n";
	for (auto pair : cases) {
		auto got = pair.second.marshal_len();
		auto want = pair.first.size() / 2;
		if (got != want)
			std::cout << "0x" << pair.first << ": got marshal_len " << got << "!\n";
	}

	unsigned char buf[1024];

	std::cout << "TEST marshal...\n";
	for (auto pair : cases) {
		auto n = pair.second.marshal_to(buf);
		auto want_n = pair.first.size() / 2;
		if (n != want_n)
			std::cout << "0x" << pair.first << ": marshal wrote " << n << " octets!\n";
	}
}
