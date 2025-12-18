package main

import (
	"fmt"
	"jgr0sz/nzbgo/parser"
)

func parsertest() {
	nzb, err := parser.FromFile("tests/big_buck_bunny.nzb")
	if err != nil {
		fmt.Printf("Error parsing NZB. %v\n", err)
		return
	}
	parser.ToString(nzb)
	largestFile := parser.MainFile(nzb)
	fmt.Printf("\n\nThe largest file in the nzb is %s with %d bytes.\n", largestFile.Subject, parser.FileSize(largestFile))
	fmt.Printf("The unique groups are %v.\n", parser.Groups(nzb))

	fmt.Printf("This NZB is %d bytes. %f percent of it is par2 files.\n", parser.Size(nzb), parser.Par2_percentage(nzb))
	
	for _, s := range nzb.Files {
		fmt.Printf("File date: %s\n", parser.DatePosted(s).Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("The filename of this NZB's first file is: %s\n", parser.ExtractFilename(nzb.Files[0]))
	filenames := parser.Filenames(nzb)

	for _, f := range filenames {
		fmt.Printf("The full list of this NZB's filenames are: %s\n", f)
	}

	fmt.Printf("Here is our nzb in JSON we can use here in the future: \n%s", parser.ToJSON(nzb))
}