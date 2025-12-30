// Allows for the idiomatic modification of NZB metadata fields, particularly useful when dealing with automated systems reliant on it.

package metaeditor

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/jgr0sz/nzbgo/parser"
)

/*
	Central set of functions for NZB metadata editing. Note that operation
*/

//Takes and instantiates an NzbMetaEditor object from a file path.
func From_file(path string) (*NzbMetaEditor, error) {
	//Extracting an Nzb's metadata 
	Nzb, err := parser.FromFile(path) 

	//parser.FromFile() returns a nil pointer if it couldn't find a file, therefore we must check if it hasn't before dereferencing it (segfault).
	if err != nil {
		log.Printf("Error parsing the NZB for metadata. %v", err)
		return nil, err
	}

	//Assigning and returning our editor object
	return &NzbMetaEditor{
		Metadata: Nzb.Head.Meta,
		Nzb: Nzb,
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
		//Useful unction to efficiently remove located duplicates using a defined predicate func.
		editor.Metadata = slices.DeleteFunc(editor.Metadata, func(m parser.Meta) bool {
			return m.Type == attribute
		})
	}
}

//Set metadata fields, replacing one or more of the same type of fields.
func Set(editor *NzbMetaEditor, attribute string, content string) {
	if IsAttribute(attribute) {
		for _, m := range editor.Metadata {
			if m.Type == attribute {
				m.Value = content
			}
		}
	}
}

//Default sorting pattern.
var defaultPattern = map[string]int {
	"title": 0,
	"category": 1, 
	"password": 2, 
	"tag": 3,
}

//Sorting function for metadata fields. Default is "title": 0,"category": 1, "password": 2, "tag": 3, (where 0 is highest-priority). A custom pattern may be provided, given it followed map[string]int; note that patterns missing fields or files with non-NZB fields will have their ordering preserved.
func Sort(editor *NzbMetaEditor, pattern map[string]int) {
	var sortPattern = defaultPattern

	if pattern != nil {
		sortPattern = pattern
	}

	//SliceStable() is used to preserve original ordering in the case that some metadata field either:
	//Doesn't exist in defaultPattern, or is missing in pattern.
	sort.SliceStable(editor.Metadata, func(a, b int) bool {
		typeA, okA := sortPattern[editor.Metadata[a].Type]
		typeB, okB := sortPattern[editor.Metadata[b].Type]
		
		//If either key doesn't exist in the map
		if !okA || !okB {
			return false
		}
	
		return typeA < typeB
	})
}

//Using NzbMetaEditor's associated Nzb class, we write the modified metadata to it and export it. Optional overwriting of same-named files.
func ToFile(editor *NzbMetaEditor, path string, overwrite bool) {
	//Sanitizes the path 
	path = filepath.Clean(path)
	//Using the stored pointer to the original Nzb object in NzbMetaEditor, we can apply the changes and marshal the updated NZB.
	editor.Nzb.Head.Meta = editor.Metadata
	nzbFile, err := xml.MarshalIndent(editor.Nzb, "", "  ")

	if err != nil {
		log.Printf("Unable to marshal Nzb instance to JSON: %v", err)
	}

	//In the case that the file exists and overwrite is disabled.
	if _, err := os.Stat(path); err == nil && !overwrite {
		log.Printf("Cannot write to file: file exists, overwriting disabled.")
		return
	}
	err = os.WriteFile(path, nzbFile, 0644)

	if err != nil {
		log.Printf("Unable to write file to %s: %v", path, err)
	}
}

//Using NzbMetaEditor's associated Nzb class, we write the modified metadata to it and return a string of the NZB in XML.
func ToStr(editor *NzbMetaEditor) string {
	editor.Nzb.Head.Meta = editor.Metadata
	nzbData, err := xml.MarshalIndent(editor.Nzb, "", "  ")

	if err != nil {
		log.Printf("Unable to marshal Nzb instance to JSON: %v", err)
	}
	return string(nzbData)
}
