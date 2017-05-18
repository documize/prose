// The MIT License (MIT)
//
// Copyright (c) 2015 Kevin S. Dias
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package tokenize

import (
	"regexp"
	"strings"
)

// PragmaticSegmenter is a multilingual, rule-based sentence boundary detector.
type PragmaticSegmenter struct {
	processor languageProcessor
}

// NewPragmaticSegmenter creates a new PragmaticSegmenter according to the
// specified language.
//
// TODO: Check that lang is supported.
func NewPragmaticSegmenter(lang string) *PragmaticSegmenter {
	return &PragmaticSegmenter{}
}

// Tokenize splits text into sentences.
func (p *PragmaticSegmenter) Tokenize(text string) []string {
	return p.processor.process(text)
}

// A Rule associates a regular expression with a replacement string.
type Rule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Sub replaces all occurances of Pattern with Replacement.
func (r *Rule) Sub(text string) string {
	f := r.Pattern.FindStringSubmatchIndex
	for loc := f(text); len(loc) > 0; loc = f(text) {
		text = text[:loc[2]] + r.Replacement + text[loc[3]:]
	}
	return text
}

type languageProcessor interface {
	process(text string) []string
}

// numbers

var periodBeforeNumberRule = Rule{
	Pattern: regexp.MustCompile(`(\.)\d`), Replacement: "∯"}
var numberAfterPeriodBeforeLetterRule = Rule{
	Pattern: regexp.MustCompile(`\d(\.)\S`), Replacement: "∯"}
var newLineNumberPeriodSpaceLetterRule = Rule{
	Pattern: regexp.MustCompile(`[\n\r]\d(\.)(?:[\s\S]|\))`), Replacement: "∯"}
var startLineNumberPeriodRule = Rule{
	Pattern: regexp.MustCompile(`^\d(\.)(?:[\s\S]|\))`), Replacement: "∯"}
var startLineTwoDigitNumberPeriodRule = Rule{
	Pattern: regexp.MustCompile(`^\d\d(\.)(?:[\s\S]|\))`), Replacement: "∯"}
var allNumberRules = []Rule{
	periodBeforeNumberRule, numberAfterPeriodBeforeLetterRule,
	newLineNumberPeriodSpaceLetterRule, startLineNumberPeriodRule,
	startLineTwoDigitNumberPeriodRule,
}

// common

var sentenceBoundaryRE = regexp.MustCompile(
	`\x{ff08}(?:[^\x{ff09}])*\x{ff09}(\s?[A-Z])|` +
		`\x{300c}(?:[^\x{300d}])*\x{300d}(\s[A-Z])|` +
		`\((?:[^\)]){2,}\)(\s[A-Z])|` +
		`'(?:[^'])*[^,]'(\s[A-Z])|` +
		`"(?:[^"])*[^,]"(\s[A-Z])|` +
		`“(?:[^”])*[^,]”(\s[A-Z])|` +
		`\S.*?[。．.！!?？ȸȹ☉☈☇☄]`)
var quotationAtEndOfSentenceRE = regexp.MustCompile(
	`[!?\.-][\"\'\x{201d}\x{201c}]\s{1}[A-Z]`)
var parensBetweenDoubleQuotesRE = regexp.MustCompile(`["”]\s\(.*\)\s["“]`)
var betweenDoubleQuotesRE2 = regexp.MustCompile(`(?:[^"])*[^,]"|“(?:[^”])*[^,]”`)
var splitSpaceQuotationAtEndOfSentenceRE = regexp.MustCompile(
	`[!?\.-][\"\'\x{201d}\x{201c}](\s{1})[A-Z]`) // lookahead
var continuousPunctuationRE = regexp.MustCompile(`\S(!|\?){3,}(?:\s|\z|$)`)
var possessiveAbbreviationRule = Rule{
	Pattern: regexp.MustCompile(`(\.)'s\s|(\.)'s$|(\.)'s\z`), Replacement: "∯"}
var kommanditgesellschaftRule = Rule{
	Pattern: regexp.MustCompile(`Co(\.)\sKG`), Replacement: "∯"}
var multiPeriodAbbrevRE = regexp.MustCompile(`\b[a-z](?:\.[a-z])+[.]`)

// AM/PM
var upperCasePmRule = Rule{
	Pattern: regexp.MustCompile(`P∯M(∯)\s[A-Z]`), Replacement: "."}
var upperCaseAmRule = Rule{
	Pattern: regexp.MustCompile(`A∯M(∯)\s[A-Z]`), Replacement: "."}
var lowerCasePmRule = Rule{
	Pattern: regexp.MustCompile(`p∯m(∯)\s[A-Z]`), Replacement: "."}
var lowerCaseAmRule = Rule{
	Pattern: regexp.MustCompile(`a∯m(∯)\s[A-Z]`), Replacement: "."}
var allAmPmRules = []Rule{
	upperCasePmRule, upperCaseAmRule, lowerCasePmRule, lowerCaseAmRule}

// Searches for periods within an abbreviation and replaces the periods.
var singleUpperCaseLetterAtStartOfLineRule = Rule{
	Pattern: regexp.MustCompile(`^[A-Z](\.)\s`), Replacement: "∯"}
var singleUpperCaseLetterRule = Rule{
	Pattern: regexp.MustCompile(`\s[A-Z]\.\s`), Replacement: "∯"}
var allSingleUpperCaseLetterRules = []Rule{
	singleUpperCaseLetterAtStartOfLineRule, singleUpperCaseLetterRule}

// Searches for ellipses within a string and replaces the periods.
var threeConsecutiveRule = Rule{
	Pattern: regexp.MustCompile(`(\.\.\.)\s+[A-Z]`), Replacement: "☏."}
var fourConsecutiveRule = Rule{
	Pattern: regexp.MustCompile(`\S(\.{3}\.)\s[A-Z]`), Replacement: "ƪ"}
var threeSpaceRule = Rule{
	Pattern: regexp.MustCompile(`(\s\.){3}\s`), Replacement: "♟"}
var fourSpaceRule = Rule{
	Pattern: regexp.MustCompile(`[a-z](\.\s{3}\.(?:\z|$|\n))`), Replacement: "♝"}
var otherThreePeriodRule = Rule{Pattern: regexp.MustCompile(`\.\.\.`), Replacement: "ƪ"}
var allEllipsesRules = []Rule{
	threeConsecutiveRule, fourConsecutiveRule, threeSpaceRule, fourSpaceRule,
	otherThreePeriodRule}

// between_punctuation

var betweenSingleQuotesRE = regexp.MustCompile(`\s'(?:[^']|'[a-zA-Z])*'`)
var betweenDoubleQuotesRE = regexp.MustCompile(`"([^"\\]+|\\{2}|\\.)*"`)
var betweenArrowQuotesRE = regexp.MustCompile(`«([^»\\]+|\\{2}|\\.)*»`)
var betweenSmartQuotesRE = regexp.MustCompile(`“([^”\\]+|\\{2}|\\.)*”`)
var betweenSquareBracketsRE = regexp.MustCompile(`\[([^\]\\]+|\\{2}|\\.)*\]`)
var betweenParensRE = regexp.MustCompile(`\(([^\(\)\\]+|\\{2}|\\.)*\)`)
var wordWithLeadingApostropheRE = regexp.MustCompile(`\s'(?:[^']|'[a-zA-Z])*'\S`)

func subPat(text, mtype string, pat *regexp.Regexp) string {
	canidates := []string{}
	for _, s := range pat.FindAllString(text, -1) {
		canidates = append(canidates, strings.TrimSpace(s))
	}
	r := punctuationReplacer{
		matches: canidates, text: text, matchType: mtype}
	return r.replace()
}

func replaceBetweenQuotes(text string) string {
	text = subPat(text, "single", betweenSingleQuotesRE)
	text = subPat(text, "double", betweenDoubleQuotesRE)
	text = subPat(text, "double", betweenSquareBracketsRE)
	text = subPat(text, "double", betweenParensRE)
	text = subPat(text, "double", betweenArrowQuotesRE)
	text = subPat(text, "double", betweenSmartQuotesRE)
	return text
}

// punctuation_replacer

var escapeRegexReservedCharacters = strings.NewReplacer(
	`(`, `\(`, `)`, `\)`, `[`, `\[`, `]`, `\]`, `-`, `\-`,
)

var subEscapeRegexReservedCharacters = strings.NewReplacer(
	`\(`, `(`, `\)`, `)`, `\[`, `[`, `\]`, `]`, `\-`, `-`,
)

type punctuationReplacer struct {
	matches   []string
	text      string
	matchType string
}

func (r *punctuationReplacer) replace() string {
	return r.replacePunctuation(r.matches)
}

func (r *punctuationReplacer) replacePunctuation(matches []string) string {
	r.text = escapeRegexReservedCharacters.Replace(r.text)
	for _, m := range matches {
		m = escapeRegexReservedCharacters.Replace(m)

		s := r.sub(m, ".", "∯")
		sub1 := r.sub(s, "。", "&ᓰ&")
		sub2 := r.sub(sub1, "．", "&ᓱ&")
		sub3 := r.sub(sub2, "！", "&ᓳ&")
		sub4 := r.sub(sub3, "!", "&ᓴ&")
		sub5 := r.sub(sub4, "?", "&ᓷ&")
		sub6 := r.sub(sub5, "? ", "&ᓸ&")
		if r.matchType != "single" {
			r.sub(sub6, "'", "&⎋&")
		}
	}
	return subEscapeRegexReservedCharacters.Replace(r.text)
}

func (r *punctuationReplacer) sub(content, a, b string) string {
	repl := strings.Replace(content, a, b, -1)
	r.text = strings.Replace(r.text, content, repl, -1)
	return repl
}
