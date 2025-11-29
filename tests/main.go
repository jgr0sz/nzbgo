package main

import (
	"fmt"
	"jgr0sz/nzbgo/src"
)

//Testing file for parser features.

func main() {
	nzb, err := src.FromFile("tests/big_buck_bunny.nzb")
	if err != nil {
		fmt.Println("Error parsing NZB.", err)
		return
	}
	src.ToString(nzb)
	largestFile := src.MainFile(nzb)
	fmt.Printf("\n\nThe largest file in the nzb is %s with %d bytes.\n", largestFile.Subject, src.FileSize(largestFile))
	fmt.Printf("The unique groups are %v.\n", src.Groups(nzb))

	fmt.Printf("This NZB is %d bytes. %f percent of it is par2 files.\n", src.Size(nzb), src.Par2_percentage(nzb))
	
	for _, s := range nzb.Files {
		fmt.Printf("File date: %s\n", src.DatePosted(s).Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("The filename of this NZB's first file is: %s\n", src.ExtractFilename(nzb.Files[0]))
	filenames := src.Filenames(nzb)

	for _, f := range filenames {
		fmt.Printf("The full list of this NZB's filenames are: %s\n", f)
	}
}
