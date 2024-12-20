# Slang - Language Code Parser
Slang is a Go package designed to parse language codes in various formats and standards, including ISO 639, BCP47, and Microsoft's Windows Language ID (LCID).

It comes with an embedded database sourced from ISO and Microsoft to provide language-related data.

## Features
Slang supports parsing language codes in various formats and standards, including:
  - BCP 47 tags (e.g., `en-US`, `fr-CA`).
  - ISO 639-1, 639-2, and 639-3 codes (e.g., `en`, `fra`, `zho`, `wuu`).
  - Microsoft Windows Language IDs (e.g., `CHS` for Simplified Chinese).

Slang also supports **finding the best match** with given BCP 47 tags.

## Installation
```bash
go get github.com/baobao1270/slang
```

## Usage

**Language Code Parsing**

Following is a simple example of how to use Slang to parse a language code.
```go
package main

import (
	"fmt"
	"log"

	"github.com/baobao1270/slang"
)

func main() {
	parser, err := slang.NewParser()
	if err != nil {
		log.Fatal(err)
	}
	
	// Use parser methods (example below)
	lang := parser.Parse("en-US")
	if lang != nil {
		fmt.Println("Found language:", lang.Name)
	} else {
		fmt.Println("Language not found.")
	}
}
```

**Adding Custom Language Data**

You can add custom language data to the parser by using the `AddCustom` method.
```go
parser.AddCustom(slang.Lang{
    Name: "Klingon",
    Location: "Star Trek Universe",
    MSLCID: 0x0000,
    BCP47: "kg-SU",
    WinID: "KLI",
    ISO639Set1: "kg",
    ISO639Set2: "tlh",
    ISO639Set3: "tlh",
})
```

## License
This package is open-source and is licensed under the MIT License.

However, it due to standard references, it is important to note that the data provided by this package is sourced from the following standards:
  - ISO 639-3, from [SIL International](https://iso639-3.sil.org/).
  - Microsoft Windows Language IDs, from [Microsoft](https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-lcid).

Use of this package may subject to the [ISO 639-3 Terms of Use] and [Microsoft Open Specifications Documentation Copyrihgt Notes].

[Microsoft Open Specifications Documentation Copyrihgt Notes]: https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-lcid#intellectual-property-rights-notice-for-open-specifications-documentation
[ISO 963-3 Terms of Use]: https://iso639-3.sil.org/code_tables/download_tables#termsofuse
