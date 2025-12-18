package main

import (
	"fmt"
	"jgr0sz/nzbgo/metaeditor"
)

func editortest() {
	editor, err := metaeditor.From_file("tests/meta_editor.nzb")

	if err != nil {
		return
	}

	for _, m := range editor.Metadata {
		fmt.Printf("%s : %s\n", m.Type, m.Value)
	}
}