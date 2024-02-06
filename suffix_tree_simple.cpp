#include <iostream>
#include <unordered_map>
#include <vector>
#include <string>



//global variable
std::string inputString;



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
    return node->endIdx - node->startIdx  + 1; //start and end idx are inclusive
}


//split edge function
void splitEdge(SuffixTreeNode* originalChild, int startIdx, int splitIdx, int endIdx) {

    std::cout << "Splitting edge" << std::endl;
    std::cout << "Original child start idx: " << originalChild->startIdx << " split idx: " << splitIdx << " end idx: " << originalChild->endIdx << std::endl;
    std::cout << "the values to the function:" << startIdx << " " << splitIdx << " " << endIdx << " " << std::endl;

    //create a new child
    SuffixTreeNode* newChild = new SuffixTreeNode(startIdx + splitIdx, nullptr, std::unordered_map<char, SuffixTreeNode*>(), startIdx + splitIdx, endIdx);


    std::cout << "New child start idx: " << startIdx + splitIdx << " end idx: " << newChild->endIdx << std::endl;

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

    std::unordered_map<char, SuffixTreeNode*> internalChildren;
    std::cout << "Creating internal node " << originalChild->label << " hugo" << std::endl;
    std::cout << "Creating internal node " << newChild->label << " hugo" << std::endl;
    internalChildren[inputString[originalChild->startIdx]] = originalChild;
    internalChildren[inputString[newChild->startIdx]] = newChild;

    internalNode->children = internalChildren;
    std::cout << "hugo" << std::endl;

    std::cout << "Internal node start idx: " << internalNode->startIdx << " end idx: " << internalNode->endIdx << std::endl;
}




//insert suffix beginning at idx into the suffix tree
void insertSuffix(std::string* strPointer, int suffixOffset, SuffixTreeNode* root) {
    std::cout << "Inserting suffix no: " << suffixOffset << std::endl;
    //get the length of the suffix
    int suffixLength = strPointer->length() - suffixOffset;

    //start in root
    SuffixTreeNode* currentNode = root;
    
    int depth = 0;
    while(true){
        //check if the current node has a child with the first character of the suffix
        char letter = (*strPointer)[suffixOffset + depth];
        std::cout << "Checking if current node has a child with the first character of the suffix" << std::endl;
        std::cout << "Current node start idx: " << letter << std::endl;
        if (currentNode->children.find(letter) != currentNode->children.end()) {
            std::cout << "Current node " << currentNode->children[letter]->startIdx << " "<< currentNode->children[letter]->endIdx << std::endl;
            std::cout << "Current node start idx: " << (*strPointer)[currentNode->children[letter]->startIdx] << std::endl;
            
            //if it is, slowscan through edge
            //if edge is longer than our string, we are guaranteed to mismatch on $ character anyways.
            int currentEdgeSize = edgeLength(currentNode->children[letter]);
            for (int j = 0; j < currentEdgeSize; j++) {
                std::cout << j << " " << currentEdgeSize << std::endl;
                std::cout << suffixOffset + depth + j << " " << (*strPointer).length() << std::endl;
                std::cout << currentNode->children[letter]->startIdx << " " << std::endl;
                std::cout << (*strPointer)[suffixOffset + depth + j] << " " << (*strPointer)[currentNode->children[letter]->startIdx + j] << std::endl;
                if ((*strPointer)[suffixOffset + depth + j] != (*strPointer)[currentNode->children[letter]->startIdx + j]) {
                    
                    //if the characters do not match, split the edge and insert the suffix
                    std::cout << "now calling splitedge to insert new node on an edge" << std::endl;
                    std::cout << suffixOffset << " " << depth << " " << j << std::endl;
                    splitEdge(currentNode->children[letter], suffixOffset + depth, j, (*strPointer).length()-1);
                    return;
                } 
            }
            currentNode = currentNode->children[letter];
            depth = depth + currentEdgeSize;
            //check if current node exists 
            std::cout << currentNode->endIdx << std::endl;
        } else {            
            //if it does not, create a new node and insert it as a child of the current node
            //note that we will always end here if we match completely (as we have $ character)
            std::cout << "Inserting new node on a node" << std::endl;
            SuffixTreeNode* newNode = new SuffixTreeNode(suffixOffset, currentNode, std::unordered_map<char, SuffixTreeNode*>(), suffixOffset + depth, (*strPointer).length()-1);
            currentNode->children[(*strPointer)[suffixOffset + depth]] = newNode;
            

            std::cout << "test" << suffixOffset + depth << (*strPointer)[suffixOffset + depth] << std::endl;
            std::cout << "test" << (*strPointer)[suffixOffset + depth] << std::endl;
            std::cout << newNode->startIdx << " " << newNode->endIdx << std::endl;

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


int printSuffixTree2(SuffixTreeNode* root, int depth = 0) {
    int size = 1;
    
    //print this node start and end
    for (int i = 0; i < depth; i++) {
        std::cout << "-";
    }
    std::cout << " Start index: " << root->startIdx << " End index: " << root->endIdx << std::endl;
    for (auto it = root->children.begin(); it != root->children.end(); it++) {
        size += printSuffixTree2(it->second, depth + 1);
    }
    return size;
}


//print suffix tree
void printSuffixTree(SuffixTreeNode* root) {
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

int main() {

    inputString = "abccdkdsngjkdgasgaxgaaksdnfkdlslfaaadvldsmkvldsmklfdnsjfdsfndsabaaabab$";
    SuffixTreeNode* root = createSuffixTree(inputString);
    std::cout << "Suffix tree created" << std::endl;
    printSuffixTree(root);
    std::cout << printSuffixTree2(root) << std::endl;


    return 0;
}