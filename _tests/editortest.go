package main

import (
	"fmt"

	"github.com/jgr0sz/nzbgo/metaeditor"
)

func editorTest() {
	editor, err := metaeditor.From_file("tests/nzbs/meta_editor.nzb")

	//Anon function we reuse. Note the closure editor is being used, so printEditor will always refer to the same object.
	printEditor := func(){
		fmt.Printf("%+v\n\n", editor)
	}

	if err != nil {
		return
	}
	
	metaeditor.Append(editor, "password", "12345")
	printEditor()
	metaeditor.Remove(editor, "password")
	metaeditor.Append(editor, "password", "foobar!")
	printEditor()
	metaeditor.Sort(editor, nil)

	var customPattern = map[string]int {
		"password": 0,
		"title": 1,
		"category": 2, 
		"tag": 3,
	}
	
	metaeditor.Sort(editor, customPattern)
	printEditor()
	metaeditor.ToFile(editor, "tests/nzbs/samplenzb.nzb", true)
	fmt.Printf("%s", metaeditor.ToStr(editor))
}