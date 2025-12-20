package metaeditor

import (
	"jgr0sz/nzbgo/parser"
	"log"
	"slices"
	"strings"
)

//Preserves the Nzb file for exporting purposes.
var metaNzb *parser.Nzb

/*
	Central set of functions for NZB metadata editing. Note that operation 
*/

//Takes and instantiates an NzbMetaEditor object from a file path.
func From_file(path string) (*NzbMetaEditor, error) {
	//Extracting an Nzb's metadata 
	meta, err := parser.FromFile(path) 
	metaNzb = meta

	//parser.FromFile() returns a nil pointer if it couldn't find a file, therefore we must check if it hasn't before dereferencing it (segfault).
	if err != nil || metaNzb == nil {
		log.Printf("Error parsing the NZB for metadata. %v", err)
		return nil, err
	}

	//Assigning and returning our editor object
	return &NzbMetaEditor{
		Metadata: meta.Head.Meta,
	}, err
}

//Helper func to verify if Nzb metadata contains an acceptable attribute. Case insensitive.
func IsAttribute(attribute string) bool {
	attribute = strings.ToLower(attribute)
	return attribute == "category" ||
	attribute == "password" ||
	attribute == "tag" ||
	attribute == "title"
}

//Appends a metadata field to an NzbMetaEditors's current splice of them.
func Append(editor *NzbMetaEditor, attribute string, content string) {
	//After verifying our attribute, we append it as a new Meta object to the provided editor's splice.
	if IsAttribute(attribute) {
		editor.Metadata = append(editor.Metadata, parser.Meta{
			Type: attribute,
			Value: content,
		})
	}
}

//Clears all metadata fields from an NzbMetaEditor.
func Clear(editor *NzbMetaEditor) {
	editor.Metadata = nil
}

//Removes a specified metadata field from an NzbMetaEditor, as well as any duplicates it may have.
func Remove(editor *NzbMetaEditor, attribute string) {
	if IsAttribute(attribute) {
		//Crafty function to efficiently remove located duplicates using a defined predicate func.
		editor.Metadata = slices.DeleteFunc(editor.Metadata, func(m parser.Meta) bool {
			return m.Type == attribute
		})
	}
}
