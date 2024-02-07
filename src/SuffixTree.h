#ifndef SUFFIX_TREE_H
#define SUFFIX_TREE_H

#include <iostream>
#include <unordered_map>
#include <string>

// Structure to represent a node in the suffix tree
struct SuffixTreeNode {
    // fields
    int label;
    SuffixTreeNode* parent;
    std::unordered_map<char, SuffixTreeNode*> children;
    int startIdx;
    int endIdx;

    // Constructor
    SuffixTreeNode(int label, SuffixTreeNode* parent, std::unordered_map<char, SuffixTreeNode*> children, int startIdx, int endIdx);

    // Destructor
    ~SuffixTreeNode();
};

// Structure to represent the entire suffix tree
struct SuffixTree {
    // fields
    SuffixTreeNode* root;
    int length;

    // Constructor
    SuffixTree(const std::string& inputString);

    // Destructor
    ~SuffixTree();



    // Function to get the length of an edge
    // Parameters: node whose edge length is to be calculated
    // Returns: length of the edge (inclusive of start and end indices)
    int edgeLength(SuffixTreeNode* node) const;

    // Function to split an edge
    void splitEdge(SuffixTreeNode* originalChild, int startIdx, int splitIdx, int endIdx, const std::string& inputString, int suffixOffset);

    // Function to insert a suffix into the suffix tree
    void insertSuffix(const std::string& str, int suffixOffset, SuffixTreeNode* root, const std::string& inputString);

    // Function to create a suffix tree from the given input string.
    // Constructs a suffix tree data structure using Ukkonen's algorithm.
    // Parameters: input string for which the suffix tree is to be constructed.
    // Returns: A pointer to the root node of the constructed suffix tree.
    SuffixTreeNode* createSuffixTree(const std::string& inputString);

    // Function to print the suffix tree
    int printSuffixTree2(SuffixTreeNode* root, int depth) const;

    // Function to search for a substring in the suffix tree
    void printSuffixTree(SuffixTreeNode* root) const;

    // Function to search for a substring in the suffix tree
    // Should write the indices of the occurrences of the substring in the input string to the standard output
    // in standard cigar format
    void searchSubstring(const std::string& substring) const;
};

#endif // SUFFIX_TREE_H
