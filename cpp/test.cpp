#include "gen/Colfer.h"

#include <iostream>
#include <chrono>


int main() {
	gen::O x = {};
	x.t = std::chrono::seconds(7);

	std::cout << "Hello ";
	std::cout << x.marshal_len();
	std::cout << "!\n";
}
