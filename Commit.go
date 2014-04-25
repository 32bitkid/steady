package main

import "bufio"
import "bytes"
import "strings"
import "log"
import "strconv"
import "regexp"
import "time"
import "fmt"

type Hash string
type HashList []Hash
type AssetNumber string
type AssetNumberList []AssetNumber

type Commit struct {
	Hash       Hash
	Parents    HashList
	AuthorName string
	When       time.Time
	Message    string
	References AssetNumberList
}

type CommitParser func(string) *Commit

type ParserOptions struct {
	CommitParser
	bufio.SplitFunc
	GitArgs []string
}

type Commits []*Commit

func (c *Commit) IsMerge() bool {
	return len(c.Parents) > 1
}

func (c Commits) AllReferences() AssetNumberList {
	cache := make(map[AssetNumber]bool, 0)
	for _, commit := range c {
		for _, number := range commit.References {
			cache[number] = true
		}
	}

	list := make(AssetNumberList, 0, len(cache))
	for number := range cache {
		list = append(list, number)
	}
	return list
}

const formatUnitSeparator string = "%x1f"
const actualUnitSeparator string = "\x1f"

var referenceRegex = regexp.MustCompile(`[A-Za-z]-\d+`)

var formatTokens = []string{"%h", "%p", "%an", "%at", "%B"}
var prettyArg = fmt.Sprintf("--pretty=tformat:%s", strings.Join(formatTokens, formatUnitSeparator))
var formatArgs = []string{"log", prettyArg, "-z"}

func parseCommit(data string) *Commit {
	parts := strings.SplitN(data, actualUnitSeparator, len(formatTokens))

	message := parts[4]

	parents := strings.Split(parts[1], " ")
	unix, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	parentHashes := make(HashList, 0, len(parents))
	for _, parent := range parents {
		parentHashes = append(parentHashes, Hash(parent))
	}

	referenceStrings := referenceRegex.FindAllString(message, -1)
	referencedAssets := make(AssetNumberList, 0, len(referenceStrings))
	for _, number := range referenceStrings {
		referencedAssets = append(referencedAssets, AssetNumber(number))
	}

	return &(Commit{
		Hash(parts[0]),
		parentHashes,
		parts[2],
		time.Unix(unix, 0),
		message,
		referencedAssets,
	})
}

func nulSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\x00'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

var DefaultOptions = ParserOptions{parseCommit, nulSplitter, formatArgs}
