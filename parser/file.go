package parser

import (
	"regexp"
	"strings"
	"time"
	"unicode"
)

//Compiled regexes for filename extraction.
var EXTRACTION_PATTERNS = [3]*regexp.Regexp {
	regexp.MustCompile(`"(.*)"`),
	//Modified the 2nd statement to best emulate Python's fullmatch(), (^...$) https://github.com/Ravencentric/nzb/blob/aa5d11dfed61b49b3b3ed5c00226b88fad7e591b/src/nzb/_subparsers.py#L35C7-L35C93
	regexp.MustCompile(`^(?:\[|\()(?:\d+/\d+)(?:\]|\))\s-\s(.*)\syEnc\s(?:\[|\()(?:\d+/\d+)(?:\]|\))\s\d+$`),
	regexp.MustCompile(`\b([\w\-+()' .,]+(?:\[[\w\-/+()' .,]*][\w\-+()' .,]*)*\.[A-Za-z0-9]{2,4})\b`),
}

//Compiled regex for the splitting pattern between stem/extension.
var SPLITTING_PATTERN = *regexp.MustCompile(`(\.[a-z]\w{2,5})$`)

//Compiled regexes for stem-obfuscated patterns.
var OBFUSCATION_PATTERNS = [5]*regexp.Regexp {
	regexp.MustCompile(`^[a-f0-9]{32}$`), 
	regexp.MustCompile(`^[a-f0-9.]{40,}$`),

	//These two are meant to be used in a singular if-statement
	regexp.MustCompile(`[a-f0-9]{30}`),
	//len() of this >= 2 is a condition
	regexp.MustCompile(`\[\w+\]`),

	regexp.MustCompile(`^abc\.xyz`),
}

//Converts and retrieves the Unix-timestamped File field into UTC format.
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

//Helper func that when using a precompiled regex to search for a filename, returns it and its length, otherwise an empty string and int.
func FilenameSearch(r *regexp.Regexp, file File) (string){
	//We take the pattern and see if there's a match.
	filenameResult := r.FindStringSubmatch(file.Subject)
	if filenameResult == nil {
		return ""
	}
	//Assuming there is one, we return it and its length.
	return filenameResult[1]
}

//Searches for a filename within a subject. https://github.com/Ravencentric/nzb/blob/aa5d11dfed61b49b3b3ed5c00226b88fad7e591b/src/nzb/_subparsers.py#L24-46 for more info.
func ExtractFilename(file File) string {
	//Using the global var of precompiled regex patterns, we iterate through and apply a filename search, seeing if we get a filename back, otherwise an empty string.
	for _, Fname := range EXTRACTION_PATTERNS {
		filename := FilenameSearch(Fname, file)
		if len(filename) != 0 {
			return strings.TrimSpace(filename)
		}
	}
	return ""
}

//Splits a filename into two components: a stem and an extension. If no valid extension was found, return only the filename, it being a stem. 
func SplitFilename(filename string) (string, string) {
	//Extracting the index, splitIndexes is constructed like this:
	//[fullStart, fullEnd, group1Start, group1End, group2Start...], where groups are the leftmost submatches and full is the entire string matched.
	splitIndexes := SPLITTING_PATTERN.FindStringSubmatchIndex(filename)
	if splitIndexes == nil {
		return filename, ""
	}
	stem := filename[:splitIndexes[0]]
	extension := filename[splitIndexes[2]+1: splitIndexes[3]]
	return stem, extension
}

//Determines if a file stem is likely obfuscated or not. Relies on collected SplitFilename() output.
//More info: https://github.com/sabnzbd/sabnzbd/blob/297455cd35c71962d39a36b7f99622f905d2234e/sabnzbd/deobfuscate_filenames.py#L104

func IsObfuscated(stem string) bool {
	//In the case where no stem exists
	if stem == "" {
		return true
	}

	//Commonly-used obfuscation patterns 
	//(Using FindString() for speed because we need >0 matches)
	if OBFUSCATION_PATTERNS[0].FindString(stem) != "" {
		return true
	}

	if OBFUSCATION_PATTERNS[1].FindString(stem) != "" {
		return true
	}

	if OBFUSCATION_PATTERNS[2].FindString(stem) != "" && 
	len(OBFUSCATION_PATTERNS[3].FindAllString(stem, -1)) >= 2 {
		return true
	}

	if OBFUSCATION_PATTERNS[4].FindString(stem) != "" {
		return true
	}

	//Variables to store the presence of unobfuscated indicators
	var (
		decimals,
		upperchars,
		lowerchars,
		spacesdots int
	) 
	
	//Collecting stem information prior to deducing an unobfuscated file stem
	for _, element := range stem {
		//Checks if the element is within the numerical rangetable
		if unicode.Is(unicode.Number, element) {
			decimals++
		}

		if unicode.IsUpper(element) {
			upperchars++
		}

		if unicode.IsLower(element) {
			lowerchars++
		}

		if element == ' ' || 
		element == '.' || 
		element == '_' {
			spacesdots++
		}
	}

	//Common unobfuscated patterns
	//"Great Pretender"
	if upperchars >= 2 && 
	lowerchars >= 2 && 
	spacesdots >= 1 {
		return false
	}

	//"this is a regular name"
	if spacesdots >= 3 {
		return false
	}

	//"Spiderman 2021"
	if (upperchars + lowerchars) >= 4 && 
	decimals >= 4 && 
	spacesdots >= 1 {
		return false
	} 

	//"Gattaca", our stem starts with a capital letter and most of its letters are lowercase
	if unicode.IsUpper(rune(stem[0])) && 
	lowerchars > 2 && 
	float32(upperchars / lowerchars) <= 0.25 {
		return false
	}
	return true
}

//Determines whether a file extension is present within a filename or not. Case and dot insensitive.
func HasExtension(file File, ext string) bool {
	//Extracting a filename first, then acquiring the extension from splitting it
	_, fileExtension := SplitFilename(ExtractFilename(file))
	//Under case-folding, check if there is said extension within our filename.
	return strings.EqualFold(fileExtension, strings.TrimPrefix(ext, "."))
}
