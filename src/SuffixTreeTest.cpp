#include "SuffixTree.h"
#include <iostream>

int main() {
    std::string inputString = "abccdkdsngjkdgasgaxgaaksdnfkdlslfaaadvldsmkvldsmklfdnsjfdsfndsabaaabab$";
    SuffixTree st(inputString);
    std::cout << "Suffix tree created" << std::endl;
    st.printSuffixTree(st.root);
    std::cout << st.printSuffixTree2(st.root, 0) << std::endl;

    return 0;
}