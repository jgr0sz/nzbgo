package src

import (
	"fmt"
	"strings"
)

//Checks if the NZB file is .gz using its extension.
func IsGzip(path string) bool {
	//Sanitizing extension case
	path = strings.ToLower(path)
	return strings.HasSuffix(path, ".gz") || strings.HasSuffix(path, ".gzip")
}

//Checks if a file is a .par2 file using its extension.
func IsPar2(file *File) bool {
	return strings.Contains(file.Subject, ".par2")
}

//Retrieves the splice of File objects included in the NZB. Synonymous with nzb.Files.
func GetFiles(nzb *Nzb) []File {
	return nzb.Files
}

//Retrieves file's group name/s.
func GroupName(file *File) []string {
	return file.Groups
}

//Outputs an NZB's important mapped fields.
func ToString(nzb *Nzb) {
	//Prints out NZB metadata, if any.
	fmt.Println("=== META ===")
	for _, m := range nzb.Meta {
		fmt.Printf("%s: %s\n", m.Type, m.Value)
	}

	//Prints out files and the contents of their attributes.
	fmt.Println("\n=== FILES ===")
	for i, f := range nzb.Files {
		fmt.Printf("File %d:\n", i+1)
		fmt.Printf("\tPoster:\t%s\n", f.Poster)
		fmt.Printf("\tDate:\t%d\n", f.Date)
		fmt.Printf("\tSubject:\t%s\n", f.Subject)

		//Prints out groups of a file
		if len(f.Groups) > 0 {
			fmt.Println("\tGroups:")
			for _, g := range f.Groups {
				fmt.Printf("\t- %s\n", g)
			}
		}

		//Prints out file segments, with each segment's attribute information and data
		if len(f.Segments) > 0 {
			fmt.Println("\tSegments:")
			for _, s := range f.Segments {
				fmt.Printf("\t- Number: %d, Bytes: %d, ID: %s\n", s.Number, s.Bytes, s.ID)
			}
		}
	}
}
