package main


type Feature struct {
	Label string                      // Human readable label for Id	
	Archived bool                     // True if retired/rolled out. These will always return true but will continue to count accesses

	Expression string                 // An expression to evaluate eg: true, false 
	
	Enabled map[string]map[string]bool   // 1st string is the environment 2nd String is entity id (eg account or user pk)
	Disabled map[string]map[string]bool  // Ditto
}

var Features map[string]Feature  // The 'store' of features.
