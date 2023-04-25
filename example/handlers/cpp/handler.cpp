#include <iostream>
#include <string>

int main() {
    std::cout << "[C++ Handler] Started"<< std::endl;
    std::string input;
    std::getline(std::cin, input);
    std::cout << input << std::endl;
    return 0;
}