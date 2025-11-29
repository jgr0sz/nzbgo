package src

import (
	"regexp"
	"strings"
	"time"
)

//Compiled regexes for filename extraction.
var PATTERNS = []*regexp.Regexp {
	regexp.MustCompile(`"(.*)"`),
	regexp.MustCompile(`^(?:\[|\()(?:\d+/\d+)(?:\]|\))\s-\s(.*)\syEnc\s(?:\[|\()(?:\d+/\d+)(?:\]|\))\s\d+`),
	regexp.MustCompile(`\b([\w\-+()' .,]+(?:\[[\w\-/+()' .,]*][\w\-+()' .,]*)*\.[A-Za-z0-9]{2,4})\b`),
}

//Function to convert and retrieve the Unix-timestamped File field into UTC format.
func DatePosted(file File) time.Time {
	return time.Unix(file.Date, 0).UTC()
}

//Adds up all of the segment byte sizes of a file, returning its overall size.
func FileSize(file File) int {
	var segmentByteSize int
	for _, s := range file.Segments {
		segmentByteSize += s.Bytes
	}
	return segmentByteSize
}

//Function that uses a precompiled regex to search for a filename, will return it and its length, otherwise an empty string and int.
func FilenameSearch(r *regexp.Regexp, file File) (string, int){
    if r.MatchString(file.Subject) {
        return strings.TrimSpace(r.FindStringSubmatch(file.Subject)[1]), len(strings.TrimSpace(r.FindStringSubmatch(file.Subject)[1]))
    }
    return "", len("")
}

//Function that searches for a filename within a subject. https://github.com/Ravencentric/nzb/blob/aa5d11dfed61b49b3b3ed5c00226b88fad7e591b/src/nzb/_subparsers.py#L24-46 for more info.
func ExtractFilename(file File) string {
	//Using the global var of precompiled regex patterns, we iterate through and apply a filename search, seeing if we get a filename back, otherwise an empty string.
	for _, Fname := range PATTERNS {
		filename, length := FilenameSearch(Fname, file)
		if length > 0 {
			return filename
		}
	}
	return ""
}
