#include <iostream>
#include <unordered_map>
#include <vector>
#include <string>



// Structure to represent a node in the suffix tree
struct SuffixTreeNode {
    int label;
    SuffixTreeNode* parent;
    std::unordered_map<char,SuffixTreeNode*> children;
    int startIdx;
    int endIdx;

    //constructor
    SuffixTreeNode(int label, SuffixTreeNode* parent, std::unordered_map<char, SuffixTreeNode*> children, int startIdx, int endIdx) {
        this->label = label;
        this->parent = parent;
        this->children = children;
        //add start and end idx
        this->startIdx = startIdx;
        this->endIdx = endIdx;        
    }
};

struct SuffixTree {
    std::string inputString;
    struct SuffixTreeNode *root;
    int length; //????????

    //constructor
    SuffixTree(std::string inputString, SuffixTreeNode *root, int length) {
        this->inputString = inputString;
        this->root = root;
        this->length = length;
    }
};




//get length of edge
int edgeLength(SuffixTreeNode* node) {
    return node->endIdx - node->startIdx; //start and end index are inclusive
}


//split edge function
void splitEdge(SuffixTreeNode* originalChild, int startIdx, int splitIdx, int endIdx) {
    //create a new child
    SuffixTreeNode* newChild = new SuffixTreeNode(splitIdx, nullptr, std::unordered_map<char, SuffixTreeNode*>(), splitIdx, endIdx);

    //create new internal node
    std::unordered_map<char, SuffixTreeNode*> internalChildren;
    internalChildren[originalChild->label] = originalChild;
    internalChildren[newChild->label] = newChild;
    SuffixTreeNode* internalNode = new SuffixTreeNode(splitIdx-1, originalChild->parent, internalChildren, startIdx, splitIdx-1);

    //add internal node as parent to new child
    newChild->parent = internalNode;

    std::unordered_map<char, SuffixTreeNode*> originalParentChildren;
    //remove originalChild from originalParent->children
    originalParentChildren.erase(originalChild->label);
    //add internal node as child instead
    originalParentChildren[internalNode->label] = internalNode;
    originalChild->parent->children = originalParentChildren;

    //update original child
    originalChild->parent = internalNode;
    originalChild->startIdx = splitIdx;
}




//insert suffix beginning at idx into the suffix tree
void insertSuffix(std::string* strPointer, int suffixOffset, SuffixTreeNode* root) {
    //get the length of the suffix
    int suffixLength = strPointer->length() - suffixOffset;

    //start in root
    SuffixTreeNode* currentNode = root;
    
    int depth = 0;
    while(true){
        //check if the current node has a child with the first character of the suffix
        if (currentNode->children.find((*strPointer)[suffixOffset + depth]) != currentNode->children.end()) {
            //if it is, slowscan through edge
            //if edge is longer than our string, we are guaranteed to mismatch on $ character anyways.
            for (int j = 0; j < edgeLength(currentNode->children[(*strPointer)[suffixOffset + depth]]); j++) {
                if ((*strPointer)[suffixOffset + depth + j] != (*strPointer)[(currentNode->children[(*strPointer)[suffixOffset + depth]])->startIdx + j]) {
                    
                    //if the characters do not match, split the edge and insert the suffix
                    splitEdge(currentNode->children[(*strPointer)[suffixOffset + depth]], suffixOffset + depth, suffixOffset + depth + j, suffixLength - 1);
                    return;
                } 
            }
            depth = depth + edgeLength(currentNode->children[(*strPointer)[suffixOffset + depth]]);
            currentNode = currentNode->children[(*strPointer)[suffixOffset + depth]];
        } else {            
            //if it does not, create a new node and insert it as a child of the current node
            //note that we will always end here if we match completely (as we have $ character)
            SuffixTreeNode* newNode = new SuffixTreeNode(suffixOffset + depth, currentNode, std::unordered_map<char, SuffixTreeNode*>(), suffixOffset + depth + 1, suffixLength - 1);
            currentNode->children[(*strPointer)[suffixOffset + depth]] = newNode;

            currentNode = newNode;

            return;
        }
    }
}


//creaate suffix tree. Takes a string and returns the root of the suffix tree
SuffixTreeNode* createSuffixTree(std::string inputString) {

    //create a root node
    SuffixTreeNode* root = new SuffixTreeNode(-1, nullptr, std::unordered_map<char, SuffixTreeNode*>(), 0, 0);


    for (int i = 0; i < inputString.length(); i++) {
        //insert all suffixes of inputString into the suffix tree
        insertSuffix(&inputString, i, root); 
    }
    
    return root;

}


//print suffix tree
void printSuffixTree(SuffixTreeNode* root) {
    std::cout << "Printing suffix tree" << std::endl;
    std::cout << "Root" << std::endl;
    for (auto it = root->children.begin(); it != root->children.end(); it++) {
        std::cout << "Edge label: " << it->first << " Start index: " << it->second->startIdx << " End index: " << it->second->endIdx << std::endl;
        for (auto it2 = it->second->children.begin(); it2 != it->second->children.end(); it2++) {
            std::cout << "Edge label: " << it2->first << " Start index: " << it2->second->startIdx << " End index: " << it2->second->endIdx << std::endl;
        }
    }
}

int main() {

    std::string inputString = "abcd$";
    SuffixTreeNode* root = createSuffixTree(inputString);
    std::cout << "Suffix tree created" << std::endl;
    printSuffixTree(root);


    return 0;
}