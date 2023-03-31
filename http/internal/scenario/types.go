package scenario

// session is a struct that contains a slice of pointers to paragraph
type session struct {
	paragraphArray []*paragraph
}

// paragraph is one pair of question and answer, and it records the tree structure
type paragraph struct {
	// prev is the previous paragraph, root is the root paragraph, hidden(should be renamed as shadow)
	// is the paragraph that is not shown in the tree structure
	prev, root, hidden *paragraph
	// children is the slice of pointers to the children paragraphs
	children []*paragraph
	origin   string // the original question
	param    string // the parameter of the question
	actual   string // the actual question, which is the origin question with the parameter replaced + script executed
	answer   string // the answer of the question
	rule     *rule  // the rule that is used to generate the answer
}
