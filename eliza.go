// Package eliza provides an implementation of the Eliza chatbot
package eliza

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/necrophonic/log"
)

func init() {
	log.InitFromString(os.Getenv("ELIZA_LOG_LEVEL"))
}

// Analyse performs psychoanalysis on the given sentance
func Analyse(this []byte) ([]byte, error) {
	response, err := AnalyseString(string(this))
	if err != nil {
		return nil, err
	}
	return []byte(response), nil
}

// AnalyseString performs psychoanalysis on the given sentance string
func AnalyseString(this string) (string, error) {
	// nb. These steps aren't necessarily the most effecient as some things
	// could be combined - but they're laid out like this to more clearly
	// document the alogrithm.

	// Firstly split sentance into words separated by spaces
	words := split(strings.Trim(this, "\n"))

	// Second, perform pre-substitution in the word list
	words = preSubstitute(words)

	// Third, make a list of all keywords in the input words sorted into
	// descending weight
	keywords := identifyKeywords(words)

	// Fourth, run through each keyword and match against decomposition
	// sequences until a match is found. If a match is found then process
	// the reassembly for that word and move to post processing, otherwise
	// move to the next keyword.
	// This will also post process any post-sub words.
	words, err := processKeywords(keywords, words)
	if err != nil {
		return "", err
	}

	return strings.Join(words, " "), nil
}

func split(said string) []string {
	words := strings.Split(said, " ")
	for i, w := range words {
		words[i] = strings.ToLower(strings.Trim(w, ".!?"))
	}
	return words
}

func preSubstitute(words []string) []string {
	log.Trace("Running pre substitutions")
	for i, w := range words {
		if sub, ok := pre[w]; ok {
			words[i] = sub
		}
	}
	return words
}

func postSubstitute(words []string) []string {
	log.Trace("Running post substitutions")
	for i, w := range words {
		if sub, ok := post[w]; ok {
			words[i] = sub
		}
	}
	return words
}

func chooseAssembly(d *decomp) string {
	// Grab the next asseumbly and then
	// increment (and loop if needed) the counter to
	// call the next asseumbly next time around.
	chosen := d.Assemblies[d.AssemblyNext]
	d.AssemblyNext = (d.AssemblyNext + 1) % uint8(len(d.Assemblies))

	if strings.HasPrefix(chosen, "goto ") {
		// It's a jump command rather than an actual reassembly
		// Find where to jump to and then retrieve the proper key
		g := strings.TrimPrefix(chosen, "goto ")
		chosen = chooseAssembly(keywordMap[g].Decompositions[0])
	}
	return chosen
}

// Sort keys by weight - implements sort.Interface
type byWeight []keyword

func (a byWeight) Len() int           { return len(a) }
func (a byWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byWeight) Less(i, j int) bool { return a[i].Weight > a[j].Weight }

func identifyKeywords(words []string) (keys []keyword) {
	log.Debug("Attempting to identify keywords")
	for _, w := range words {
		log.Tracef("Checking if '%s' is a keyword", w)
		if k, ok := keywordMap[w]; ok {
			log.Tracef("Identified keyword -> '%s'", w)
			keys = append(keys, k)
		}
	}
	// Sort in descending order and then append the default case to the end
	sort.Sort(byWeight(keys))
	keys = append(keys, keywordMap["xnone"])
	return
}

func processKeywords(keywords []keyword, words []string) ([]string, error) {

	for _, kw := range keywords {
		// Get the pattern for the keyword and attempt to match it to the words we have
		for _, d := range kw.Decompositions {

			pattern := d.Pattern

			// Deal with synonyms
			// If we have a word in the pattern prefixed with a @ then it needs to be
			// substituted with all possible synonyms.
			// nb. May be more effecient to bake these directly into the pattern definitions
			// but that negates ease of adding new synonyms in future.
			for k := range synonyms {
				synonymKey := "@" + k
				if strings.Contains(pattern, synonymKey) {
					pattern = strings.Replace(pattern, synonymKey, fmt.Sprintf("(?:%s)", strings.Join(synonyms[k], "|")), -1)
				}
			}

			sentance := strings.Join(words, " ")
			log.Tracef("Process keywords: Attempt to match pattern '%s' to '%s'\n", pattern, sentance)

			re := regexp.MustCompile(pattern)
			results := re.FindStringSubmatch(sentance)
			if len(results) > 0 {
				resassmbly := chooseAssembly(d)
				log.Debugf("Process keywords: Matched regex [%s] -> now using assembly [%s]\n", pattern, resassmbly)

				for i, match := range results {
					// Before the matched text is subbed back in, it needs to post substituted
					// replace "I" with "You" etc using the 'post' substitution list
					match := strings.Join(postSubstitute(strings.Split(strings.Trim(match, " "), " ")), " ")
					resassmbly = strings.Replace(resassmbly, fmt.Sprintf("(%d)", i), match, -1)
				}

				return strings.Split(resassmbly, " "), nil
			}
		}
	}

	return nil, errors.New("Failed to process keywords - no clauses matched")
}
