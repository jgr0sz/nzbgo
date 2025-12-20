package main

import (
	"fmt"
	"jgr0sz/nzbgo/metaeditor"
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
	printEditor()
}