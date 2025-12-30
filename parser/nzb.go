// Allows for the parsing of information used in an NZB file to reconstruct a remotely-hosted file.
package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"os"
	"slices"
)

/*
	Central set of functions meant to aid in the parsing of NZBs.
*/

//Takes and instantiates an Nzb from a provided string.
func FromStr(nzbString string) (*Nzb, error) {
	var strNzb Nzb
	err := xml.Unmarshal([]byte(nzbString), &strNzb)

	if err != nil {
		log.Printf("Unable to unmarshal string data for parsing: %v", err)
		return nil, err
	}
	return &strNzb, err
}

//Takes and instantiates an Nzb from a JSON file.
func FromJSON(jsonData []byte) (*Nzb, error) {
	var jsonNzb Nzb
	err := json.Unmarshal(jsonData, &jsonNzb)

	if err != nil {
		log.Printf("Unable to unmarshal JSON data for parsing: %v", err)
		return nil, err
	}
	return &jsonNzb, err
}

//Takes and instantiates an Nzb from a file path.
func FromFile(path string) (*Nzb, error) {
	file, err := os.ReadFile(path)

	//Handling gzip
	if (IsGzip(path)) {
		//gzip.Reader decompresses the file
		gzReader, err := gzip.NewReader(bytes.NewReader(file))
		if err != nil {
			return nil, err
		}
		defer gzReader.Close()

		//Bytes are read back into the file byte splice once decompressed
		file, err = io.ReadAll(gzReader) 
		if err != nil {
			return nil, err
		}
	}

	//Invalid filepath to NZB
	if err != nil {
		return nil, err
	}

	//Accesses NZB file content, which using mapped struct tags groups information together.
	var NzbContents Nzb
	err = xml.Unmarshal(file, &NzbContents)
	//Important check, else unable to proceed
	if err != nil  {
		log.Printf("Unable to unmarshal NZB file for parsing: %v", err)
		return nil, err
	}
	return &NzbContents, err
}

//Finds the main content file in the NZB. This is determined by finding the largest file without the .par2 extension.
func MainFile(nzb *Nzb) File {
	fileSizes := []int{}
	for _, f := range nzb.Files {
		//Ignore the file entirely if it's a .par2
		if IsPar2(&f) {
			fileSizes = append(fileSizes, 0)
			continue
		}
		//Calculate and store filesizes
		fileSizes = append(fileSizes, FileSize(f))
	}

	mainFileSize := 0
	mainFileIdx := -1
	//Once calculated and stored, we can search for the largest filesize, conveniently giving us the index
	//of the file in nzb.Files, as fileSizes dynamically fills proportional to it.
	for i, size := range fileSizes {
		if (size > mainFileSize) {
			mainFileSize = size
			mainFileIdx = i
		}
	}
	return nzb.Files[mainFileIdx]
}

//TODO: see if more efficient methods exist

//Retrieves the splice of unique posters of an NZB's Files[] field.
func Posters(nzb *Nzb) []string {
	postersSlice := []string{}
	for _, f := range nzb.Files {
		if (!slices.Contains(postersSlice, f.Poster)) {
			postersSlice = append(postersSlice, f.Subject)
		}
	}
	return postersSlice
}

//Retrieves the splice of unique filenames of an NZB's Files[] field, if there were any found.
func Filenames(nzb *Nzb) []string {
	filenameSlice := []string{}
	for _, f := range nzb.Files {
		if (ExtractFilename(f) != "") {
			filenameSlice = append(filenameSlice, ExtractFilename(f))
		}
	}
	return filenameSlice
}

//Retrieves the splice of unique groups of an NZB's Files[].Groups field.
func Groups(nzb *Nzb) []string {
	groupsSlice := []string{}
	//We traverse Groups[] and check each group; if they are unique (not in the slice yet), we add them.
	for _, f := range nzb.Files {
		for _, g := range f.Groups {
			if (!slices.Contains(groupsSlice, g)) {
				groupsSlice = append(groupsSlice, g)
			}
		}
	}
	return groupsSlice
}

//Computes the overall size in bytes of all the files within the NZB.
func Size(nzb *Nzb) int {
	totalSize := 0
	for _, s := range nzb.Files {
		totalSize += FileSize(s)
	}
	return totalSize
}

//Retrieves the splice of par2 files in the NZB.
func Par2_files(nzb *Nzb) []File {
	par2_slice := []File{}
	for _, f := range nzb.Files {
		if IsPar2(&f) {
			par2_slice = append(par2_slice, f)
		}
	}
	return par2_slice
}

//Computes the overall size in bytes of all the .par2 files in the NZB.
func Par2_Size(nzb *Nzb) int {
	totalSize := 0
	par2_splice := Par2_files(nzb)
	for _, p := range par2_splice {
		totalSize += FileSize(p)
	}
	return totalSize
}

//Computes the percentage of the size of all .par2 files relative to the overall size.
func Par2_percentage(nzb *Nzb) float64 {
	return float64(Par2_Size(nzb))/ float64(Size(nzb))
}

//Serializes an Nzb instance into a JSON string.
func ToJSON(nzb *Nzb) string {
	NzbJSON, err := json.MarshalIndent(nzb, "", "  ")
	if err != nil {
		log.Printf("Unable to marshal Nzb instance to JSON: %v", err)
		return ""
	}
	return string(NzbJSON)
}
