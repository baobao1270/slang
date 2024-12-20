package main

import "C"

import (
	"fmt"
	"os"
	"strings"

	"github.com/baobao1270/slang"
)

const (
	ExportErrSuccess    = 0x00
	ExportErrParse      = 0x01
	ExportErrNoSuchLang = 0x02

	ProgramName = "slang-parse"
)

func tabSprint(lang *slang.Lang) string {
	return strings.Join([]string{
		lang.Name,
		lang.Location,
		fmt.Sprintf("0x%04X", lang.MSLCID),
		lang.BCP47,
		lang.WinID,
		lang.ISO639Set1,
		lang.ISO639Set2,
		lang.ISO639Set3,
	}, "\t")
}

func SlangParse(langCode string) (int8, string) {
	parser, err := slang.NewParser()
	if err != nil {
		return ExportErrParse, ""
	}

	lang := parser.Parse(langCode)
	if lang == nil {
		return ExportErrNoSuchLang, ""
	}

	return ExportErrSuccess, tabSprint(lang)
}

//export SlangParseLang
func SlangParseLang(langCode string) (int8, *C.char) {
	errCode, langInfo := SlangParse(langCode)
	return errCode, C.CString(langInfo)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <lang_code>\n", ProgramName)
		os.Exit(1)
	}

	errCode, langInfo := SlangParse(os.Args[1])
	if errCode != ExportErrSuccess {
		fmt.Printf("Error: 0x%02X\n", errCode)
		os.Exit(1)
	}

	fmt.Println(langInfo)
}
