package main

import (
	"regexp"
	"strings"
	"testing"

	"github.com/jgr0sz/nzbgo/parser"
)

var (
	nzb, _ = parser.FromFile("big_buck_bunny.nzb")
	re = regexp.MustCompile(`^(?:\[|\()(?:\d+/\d+)(?:\]|\))\s-\s(.*)\syEnc\s(?:\[|\()(?:\d+/\d+)(?:\]|\))\s\d+`)
	sinkString string
	sinkInt int
)

func OldFnameSearch (r *regexp.Regexp, file parser.File) (string, int) {
	    if r.MatchString(file.Subject) {
        return strings.TrimSpace(r.FindStringSubmatch(file.Subject)[1]), len(strings.TrimSpace(r.FindStringSubmatch(file.Subject)[1]))
    }
    return "", len("")
}

func BenchmarkFilenameSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkString = parser.FilenameSearch(re, nzb.Files[0])
	}
}

func BenchmarkOldFilenameSearch(b * testing.B) {
	for i := 0; i < b.N; i++ {
		sinkString, sinkInt = OldFnameSearch(re, nzb.Files[0])
	}
}