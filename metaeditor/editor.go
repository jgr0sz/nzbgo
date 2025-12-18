package metaeditor

import (
	"jgr0sz/nzbgo/parser"
	"log"
)

/*
	Central set of functions for NZB metadata editing.
*/

//Takes and instantiates an NzbMetaEditor object from a file path.
func From_file(path string) (*NzbMetaEditor, error) {
	//Extracting an Nzb's metadata 
	var editor NzbMetaEditor
	metaNzb, err := parser.FromFile(path) 

	//parser.FromFile() returns a nil pointer if it couldn't find a file, therefore we must check if it hasn't before dereferencing it (segfault).
	if err != nil || metaNzb == nil {
		log.Printf("Error parsing the NZB for metadata. %v", err)
		return nil, err
	}

	//Assigning and returning our editor object
	editor.Metadata = metaNzb.Head.Meta
	return &editor, err
}