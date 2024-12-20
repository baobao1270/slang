// Slang is a package for parsing language codes in different formats or standards.
// It ships with a embedded database from ISO and Microsoft.
//
// Copyright (c) 2024 Joseph Chris, under the MIT license.
//
// Use of this package is subject to [ISO 963-3 Terms of Use] and [Microsoft Open Specifications Documentation Copyrihgt Notes].
//
// [Microsoft Open Specifications Documentation Copyrihgt Notes]: https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-lcid#intellectual-property-rights-notice-for-open-specifications-documentation
// [ISO 963-3 Terms of Use]: https://iso639-3.sil.org/code_tables/download_tables#termsofuse
package slang

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"errors"
	"io"
	"sort"
	"strconv"
	"strings"
)

//go:embed langdb.csv
var db []byte

var (
	ErrParse        = errors.New("error parsing csv database")  // ErrParse is an error when parsing the database.
	ErrInvalidWinID = errors.New("invalid Windows language ID") // ErrInvalidWindowsID is an error when encountering an invalid Microsoft Windows language ID.
)

// LangParser is a parser for language database.
type LangParser struct {
	data []Lang
}

// Lang is an entry from the language database.
type Lang struct {
	// Displaying name of the language, in ASCII (may contain spaces and special characters)
	Name string

	// Location of the language, in ASCII (may contain spaces and special characters)
	Location string

	// Microsoft's LCID of the language.
	//
	// See: https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-lcid
	MSLCID uint32

	// BCP47 tag of the language (example: en-US).
	// See: https://tools.ietf.org/html/bcp47
	BCP47 string

	// Microsoft's Windows language ID of the language (example: CHS).
	//
	// See: https://learn.microsoft.com/en-us/dotnet/api/system.globalization.cultureinfo.threeletterwindowslanguagename
	WinID string

	// ISO 639-1 code of the language (example: zh).
	//
	// If the language does not have an ISO 639-1 code, this field will be same as ISO 639-2.
	ISO639Set1 string

	// ISO 639-2 code of the language (example: zuo).
	ISO639Set2 string

	// ISO 639-3 code of the language (example: cmn).
	//
	// For most languages, this field will be same as ISO 639-2.
	//
	// If the language is a sub-language of macrolanguage, this field will be different from ISO 639-2.
	ISO639Set3 string
}

// IsValidWinID checks if the Windows language ID is valid.
func IsValidWinID(id string) bool {
	if len(id) != 3 {
		return false
	}
	for _, c := range id {
		// Check if the character is an ASCII letter.
		if c < 'A' || (c > 'Z' && c < 'a') || c > 'z' {
			return false
		}
	}
	return strings.ToUpper(id) != "ZZZ"
}

// NewParser creates a default language parser.
func NewParser() (*LangParser, error) {
	lp := make([]Lang, 0)
	r := csv.NewReader(bytes.NewReader(db))

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, ErrParse
		}
		if len(line) != 9 {
			return nil, ErrParse
		}

		fID, fMSLCID := line[0], line[3]
		if fID == "id" {
			continue
		}

		mslcid, err := strconv.ParseUint(fMSLCID[2:], 16, 32)
		if err != nil {
			return nil, ErrParse
		}

		lp = append(lp, Lang{
			Name:       line[1],
			Location:   line[2],
			MSLCID:     uint32(mslcid),
			BCP47:      line[4],
			WinID:      line[5],
			ISO639Set1: line[6],
			ISO639Set2: line[7],
			ISO639Set3: line[8],
		})
	}
	return &LangParser{data: lp}, nil
}

// AddCustom adds custom language to the parser.
func (p *LangParser) AddCustom(lang Lang) *LangParser {
	p.data = append(p.data, lang)
	return p
}

// FindAllByBCP47 returns all possible values matching the BCP47 tag with best matching order.
//
// Case insensitive, support both dash (-) and underscore (_) as separator.
//
// # Examples
//  1. "en-US" will return [en-US en], but no "en-GB".
//  2. "bho-Deva" will return [bho-Deva bho bho-Deva-IN].
//  3. "bho-Deva-IN" will return [bho-Deva-IN bho-Deva bho] (best match goes first).
//  4. "be" will return [be be-BY] but no "bem" or "bem-ZM".
//  5. "en-Invalid" will return [en] but no "en-Invalid".
func (p *LangParser) FindAllByBCP47(bcp47 string) []Lang {
	results := []Lang{}
	tagSlices := strings.Split(stdBCP47Tag(bcp47), "-")

	// Find up
	for pos := range tagSlices {
		tag := strings.Join(tagSlices[:len(tagSlices)-pos], "-")
		for _, lang := range p.data {
			if strings.EqualFold(lang.BCP47, tag) {
				results = append(results, lang)
			}
		}
	}

	// Find down
	for _, lang := range p.data {
		if strings.HasPrefix(strings.ToLower(lang.BCP47), stdBCP47Tag(bcp47)+"-") {
			results = append(results, lang)
		}
	}

	return results
}

// FindAllByWinID returns all possible values matching the Windows language ID.
//
// Case insensitive. Result is sorted by BCP47 tag length.
//
// If provided Windows language ID is invalid, it will return an empty slice.
func (p *LangParser) FindAllByWinID(winID string) []Lang {
	if !IsValidWinID(winID) {
		return []Lang{}
	}

	return p.selectEqualFold(winID, func(lang Lang) string {
		return lang.WinID
	})
}

// FindAllByISO639Set1 returns all possible values matching the ISO 639-1 code.
//
// Case insensitive. Result is sorted by BCP47 tag length.
func (p *LangParser) FindAllByISO639Set1(iso639 string) []Lang {
	return p.selectEqualFold(iso639, func(lang Lang) string {
		return lang.ISO639Set1
	})
}

// FindAllByISO639Set2 returns all possible values matching the ISO 639-2 code.
//
// Case insensitive. Result is sorted by BCP47 tag length.
func (p *LangParser) FindAllByISO639Set2(iso639 string) []Lang {
	return p.selectEqualFold(iso639, func(lang Lang) string {
		return lang.ISO639Set2
	})
}

// FindAllByISO639Set3 returns all possible values matching the ISO 639-3 code.
//
// Case insensitive. Result is sorted by BCP47 tag length.
func (p *LangParser) FindAllByISO639Set3(iso639 string) []Lang {
	return p.selectEqualFold(iso639, func(lang Lang) string {
		return lang.ISO639Set3
	})
}

// FindAllByISO639Alpah3 returns all possible values matching the given ISO 639 code.
//
// Case insensitive. Result is sorted by BCP47 tag length.
//
// This function will try to find the language by ISO 639-3, then ISO 639-2, and finally ISO 639-1.
// If any found in the previous step, it will skip the next step.
func (p *LangParser) FindAllByISOCode(iso639 string) []Lang {
	results := p.FindAllByISO639Set3(iso639)
	if len(results) == 0 {
		results = p.FindAllByISO639Set2(iso639)
	}
	if len(results) == 0 {
		results = p.FindAllByISO639Set1(iso639)
	}
	return results
}

// FindByBCP47 returns the first possible best value matching the BCP47 tag.
//
// Case insensitive, support both dash (-) and underscore (_) as separator.
//
// If no value is found, it will return nil.
func (p *LangParser) FindByBCP47(bcp47 string) *Lang {
	return firstOrNil(p.FindAllByBCP47(bcp47))
}

// FindByWinID returns the first possible best value matching the Windows language ID.
//
// Case insensitive.
//
// If there is multiple possible languages found, it will return the language with the shortest BCP47 tag.
//
// If no value is found or the provided Windows language ID is invalid, it will return nil.
func (p *LangParser) FindByWinID(winID string) *Lang {
	return firstOrNil(p.FindAllByWinID(winID))
}

func (p *LangParser) selectEqualFold(value string, fieldGetter func(lang Lang) string) []Lang {
	results := []Lang{}
	for _, lang := range p.data {
		if strings.EqualFold(fieldGetter(lang), value) {
			results = append(results, lang)
		}
	}
	sortByBCP47Tag(results)
	return results
}

// IsValidWinID checks if the Windows language ID is valid.
func (lang *Lang) IsValidWinID() bool {
	return IsValidWinID(lang.WinID)
}

// FindByISO639Set1 returns the first possible best value matching the ISO 639-1 code.
//
// Case insensitive.
//
// If there is multiple possible languages found, it will return the language with the shortest BCP47 tag.
//
// If no value is found, it will return nil.
func (p *LangParser) FindByISO639Set1(iso639 string) *Lang {
	return firstOrNil(p.FindAllByISO639Set1(iso639))
}

// FindByISO639Set2 returns the first possible best value matching the ISO 639-2 code.
//
// Case insensitive.
//
// If there is multiple possible languages found, it will return the language with the shortest BCP47 tag.
//
// If no value is found, it will return nil.
func (p *LangParser) FindByISO639Set2(iso639 string) *Lang {
	return firstOrNil(p.FindAllByISO639Set2(iso639))
}

// FindByISO639Set3 returns the first possible best value matching the ISO 639-3 code.
//
// Case insensitive.
//
// If there is multiple possible languages found, it will return the language with the shortest BCP47 tag.
//
// If no value is found, it will return nil.
func (p *LangParser) FindByISO639Set3(iso639 string) *Lang {
	return firstOrNil(p.FindAllByISO639Set3(iso639))
}

// FindByISOCode returns the first possible best value matching the given ISO 639 code.
//
// Case insensitive.
//
// If there is multiple possible languages found, it will return the language with the shortest BCP47 tag.
//
// If no value is found, it will return nil.
//
// This function will try to find the language by order of ISO 639-3, then ISO 639-2, and finally ISO 639-1.
func (p *LangParser) FindByISOCode(iso639 string) *Lang {
	return firstOrNil(p.FindAllByISOCode(iso639))
}

// Parse tries to parse the language code and return the best possible language.
//
// This function will try to match in following order: BCP47, Windows language ID, ISO 639-3, ISO 639-2, ISO 639-1.
//
// If the language code is not found, it will return nil.
func (p *LangParser) Parse(value string) *Lang {
	if lang := p.FindByBCP47(value); lang != nil {
		return lang
	}
	if lang := p.FindByISOCode(value); lang != nil {
		return lang
	}
	if lang := p.FindByWinID(value); lang != nil {
		return lang
	}
	return nil
}

func stdBCP47Tag(tag string) string {
	return strings.ToLower(strings.ReplaceAll(tag, "_", "-"))
}

func sortByBCP47Tag(langs []Lang) {
	sort.Slice(langs, func(i, j int) bool {
		if len(langs[i].BCP47) == len(langs[j].BCP47) {
			return langs[i].BCP47 < langs[j].BCP47
		}
		return len(langs[i].BCP47) < len(langs[j].BCP47)
	})
}

func firstOrNil(langs []Lang) *Lang {
	if len(langs) == 0 {
		return nil
	}
	return &langs[0]
}
