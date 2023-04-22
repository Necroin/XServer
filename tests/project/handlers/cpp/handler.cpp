#include <iostream>
#include <string>

int main() {
    std::cout << "[C++ Handler] Started";
    std::string input;
    std::getline(std::cin, input);
    std::cout << input;
    return 0;
}