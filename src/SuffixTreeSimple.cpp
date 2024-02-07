#include "SuffixTree.h"
#include <iostream>
#include <unordered_map>
#include <vector>
#include <string>


// Constructor for SuffixTreeNode
SuffixTreeNode::SuffixTreeNode(int label, SuffixTreeNode* parent, std::unordered_map<char, SuffixTreeNode*> children, int startIdx, int endIdx)
    : label(label), parent(parent), children(children), startIdx(startIdx), endIdx(endIdx) {}


// Destructor for SuffixTreeNode
SuffixTreeNode::~SuffixTreeNode() {
    label = 0;
    parent = nullptr;
    children.clear();
    startIdx = 0;
    endIdx = 0;
}

// Constructor for SuffixTree
SuffixTree::SuffixTree(const std::string& inputString) : root(nullptr), length(0) {
    root = createSuffixTree(inputString);
}

// Destructor for SuffixTree
SuffixTree::~SuffixTree() {

    //delete all nodes in the suffix tree
    std::vector<SuffixTreeNode*> nodesToDelete;
    nodesToDelete.push_back(root);
    while (!nodesToDelete.empty()) {
        SuffixTreeNode* currentNode = nodesToDelete.back();
        nodesToDelete.pop_back();
        for (auto it = currentNode->children.begin(); it != currentNode->children.end(); it++) {
            nodesToDelete.push_back(it->second);
        }
        delete currentNode;
    }
    root = nullptr;
    length = 0;

}

// Function to get the length of an edge
int SuffixTree::edgeLength(SuffixTreeNode* node) const {
    return node->endIdx - node->startIdx + 1; // start and end idx are inclusive
}


//split edge function
void SuffixTree::splitEdge(SuffixTreeNode* originalChild, int startIdx, int splitIdx, int endIdx, const std::string& inputString, int suffixOffset) {
    //create a new child
    SuffixTreeNode* newChild = new SuffixTreeNode(suffixOffset, nullptr, std::unordered_map<char, SuffixTreeNode*>(), startIdx + splitIdx, endIdx);

    //create new internal node
    SuffixTreeNode* internalNode = new SuffixTreeNode(originalChild->startIdx, originalChild->parent, std::unordered_map<char, SuffixTreeNode*>(), originalChild->startIdx, originalChild->startIdx+splitIdx-1);

    //add internal node as parent to new child
    newChild->parent = internalNode;
        
    //update parent by removing original child and adding internal node
    //this is done by overwriting the original child with the internal node
    originalChild->parent->children[inputString[internalNode->startIdx]] = internalNode;

    //update original child
    originalChild->parent = internalNode;
    originalChild->startIdx += splitIdx;

    //check if they have the same starting character
    if (inputString[originalChild->startIdx] == inputString[newChild->startIdx]) {
        std::cout << "problemo mister" << std::endl;
    }

    //add original child and new child to internal node
    std::unordered_map<char, SuffixTreeNode*> internalChildren;
    internalChildren[inputString[originalChild->startIdx]] = originalChild;
    internalChildren[inputString[newChild->startIdx]] = newChild;
    internalNode->children = internalChildren;

}




//insert suffix beginning at idx into the suffix tree
void SuffixTree::insertSuffix(const std::string& str, int suffixOffset, SuffixTreeNode* root, const std::string& inputString) {    //get the length of the suffix
    int suffixLength = str.length() - suffixOffset;

    //start in root
    SuffixTreeNode* currentNode = root;
    
    int depth = 0;
    while(true){
        //check if the current node has a child with the first character of the suffix
        char letter = str[suffixOffset + depth];

        if (currentNode->children.find(letter) != currentNode->children.end()) {
            //if it is, slowscan through edge
            //if edge is longer than our string, we are guaranteed to mismatch on $ character anyways.
            int currentEdgeSize = edgeLength(currentNode->children[letter]);
            for (int j = 0; j < currentEdgeSize; j++) {
                if (str[suffixOffset + depth + j] != str[currentNode->children[letter]->startIdx + j]) {
                    
                    //if the characters do not match, split the edge and insert the suffix
                    SuffixTree::splitEdge(currentNode->children[letter], suffixOffset + depth, j, str.length()-1, inputString, suffixOffset);
                    return;
                } 
            }
            currentNode = currentNode->children[letter];
            depth = depth + currentEdgeSize;
            //check if current node exists 
        } else {            
            //if it does not, create a new node and insert it as a child of the current node
            //note that we will always end here if we match completely (as we have $ character)
            SuffixTreeNode* newNode = new SuffixTreeNode(suffixOffset, currentNode, std::unordered_map<char, SuffixTreeNode*>(), suffixOffset + depth, str.length()-1);
            currentNode->children[str[suffixOffset + depth]] = newNode;

            return;
        }
    }
}


//creaate suffix tree. Takes a string and returns the root of the suffix tree
SuffixTreeNode* SuffixTree::createSuffixTree(const std::string& inputString) {
    //create a root node
    SuffixTreeNode* root = new SuffixTreeNode(-1, nullptr, std::unordered_map<char, SuffixTreeNode*>(), 0, 0);


    for (int i = 0; i < inputString.length(); i++) {
        //insert all suffixes of inputString into the suffix tree
       SuffixTree::insertSuffix(inputString, i, root, inputString); 
    }
    
    return root;

}


int SuffixTree::printSuffixTree2(SuffixTreeNode* root, int depth) const {    int size = 1;
    
    //print this node start and end
    for (int i = 0; i < depth; i++) {
        std::cout << "-";
    }
    std::cout << " Start index: " << root->startIdx << " End index: " << root->endIdx << std::endl;
    for (auto it = root->children.begin(); it != root->children.end(); it++) {
        size += SuffixTree::printSuffixTree2(it->second, depth + 1);
    }
    return size;
}


//print suffix tree
void SuffixTree::printSuffixTree(SuffixTreeNode* root) const {
    int level = 0;
    std::cout << "Printing suffix tree" << std::endl;
    std::cout << "Root" << std::endl;
    for (auto it = root->children.begin(); it != root->children.end(); it++) {
        std::cout << "Edge label: " << it->first << " Start index: " << it->second->startIdx << " End index: " << it->second->endIdx << std::endl;
        for (auto it2 = it->second->children.begin(); it2 != it->second->children.end(); it2++) {
            std::cout << "Edge label: " << it2->first << " Start index: " << it2->second->startIdx << " End index: " << it2->second->endIdx << std::endl;
        }
    }
}


void SuffixTree::searchSubstring(const std::string& substring) const {
    //start from the root
    SuffixTreeNode* currentNode = root;

    //traverse tree until substring is found or we reach the end of the substring
    for (int i = 0; i < substring.length(); i++) {
        char letter = substring[i];
        if (currentNode->children.find(letter) != currentNode->children.end()) {
            //if the letter is found, slowscan through the edge
            //by comparing each character of the substring with the edge
            int currentEdgeSize = edgeLength(currentNode->children[letter]);
            for (int j = 0; j < currentEdgeSize; j++) {
                if (substring[i + j] != substring[currentNode->children[letter]->startIdx + j]) {
                    std::cout << "Substring not found" << std::endl;
                    return;
                }
            }

        } else {
            std::cout << "Substring not found" << std::endl;
            return;
        }
    }
    //if we reach the end of the substring, print the indices of the occurrences of the substring in the input string
    std::cout << "Substring found" << std::endl;

}

