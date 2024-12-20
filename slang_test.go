package slang_test

import (
	"testing"

	"github.com/baobao1270/slang"
)

func TestNewParser(t *testing.T) {
	_, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestAddCustomLanguage(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	lp.AddCustom(slang.Lang{
		Name:       "Klingon",
		Location:   "Star Trek Universe",
		MSLCID:     0x0000,
		BCP47:      "kg-SU",
		WinID:      "KLI",
		ISO639Set1: "kg",
		ISO639Set2: "tlh",
		ISO639Set3: "tlh",
	})

	if lp.Parse("kg-SU").Name != "Klingon" {
		t.Errorf("Error: Custom language 'Klingon' not found (parse: kg-SU)")
	}
	if lp.Parse("kLi").Name != "Klingon" {
		t.Errorf("Error: Custom language 'Klingon' not found (parse: KLI)")
	}
	if lp.Parse("kg").Name != "Klingon" {
		t.Errorf("Error: Custom language 'Klingon' not found (parse: kg)")
	}
	if lp.Parse("TlH").Name != "Klingon" {
		t.Errorf("Error: Custom language 'Klingon' not found (parse: tlh)")
	}
}

func TestFindAllByBCP47English(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	lang := lp.FindAllByBCP47("en-US")
	if lang[0].BCP47 != "en-US" {
		t.Errorf("Error: FindAllByBCP47(en-US)[0] should be 'en-US'")
	}
	if lang[1].BCP47 != "en" {
		t.Errorf("Error: FindAllByBCP47(en-US)[0] should have name 'en'")
	}
}

func TestFindAllByBCP47Deva(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	langs := lp.FindAllByBCP47("bho-Deva")
	if langs[0].BCP47 != "bho-Deva" {
		t.Errorf("Error: FindAllByBCP47(bho-Deva)[0] should be 'bho-Deva'")
	}
	if langs[1].BCP47 != "bho" {
		t.Errorf("Error: FindAllByBCP47(bho-Deva)[1] should have name 'bho'")
	}
	if langs[2].BCP47 != "bho-Deva-IN" {
		t.Errorf("Error: FindAllByBCP47(bho-Deva)[2] should have name 'bho-Deva-IN'")
	}
}

func TestFindAllByBCP47DevaBestMatchWithCaseIgnoranceAndUnderscore(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	langs := lp.FindAllByBCP47("BHO_DEVA_IN")
	if langs[0].BCP47 != "bho-Deva-IN" {
		t.Errorf("Error: FindAllByBCP47(BHO_DEVA)[0] should be 'bho-Deva-IN'")
	}
	if langs[1].BCP47 != "bho-Deva" {
		t.Errorf("Error: FindAllByBCP47(BHO_DEVA)[1] should have name 'bho-Deva'")
	}
	if langs[2].BCP47 != "bho" {
		t.Errorf("Error: FindAllByBCP47(BHO_DEVA)[2] should have name 'bho'")
	}
}

func TestFindAllByBCP47InvalidSubLanguage(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	langs := lp.FindAllByBCP47("en-Invalid")
	if langs[0].BCP47 != "en" {
		t.Errorf("Error: FindAllByBCP47(en-Invalid)[0] should be 'en'")
	}
	if len(langs) != 1 {
		t.Errorf("Error: FindAllByBCP47(en-Invalid) should have only 1 language")
	}
}

func TestFindAllByBCP47Chinese(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	langs := lp.FindAllByBCP47("zh-CN")
	if langs[0].BCP47 != "zh-CN" {
		t.Errorf("Error: FindAllByBCP47(zh-CN)[0] should be 'zh-CN'")
	}
	if langs[1].BCP47 != "zh" {
		t.Errorf("Error: FindAllByBCP47(zh-CN)[1] should have name 'zh'")
	}
	if len(langs) <= 2 {
		t.Errorf("Error: FindAllByBCP47(zh-CN) should have > 2 languages")
	}
}

func TestFindAllByBCP47ChineseTraditional(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	langs := lp.FindAllByBCP47("zh-TW")
	if langs[0].BCP47 != "zh-TW" {
		t.Errorf("Error: FindAllByBCP47(zh-TW)[0] should be 'zh-TW'")
	}
	if langs[1].BCP47 != "zh" {
		t.Errorf("Error: FindAllByBCP47(zh-TW)[1] should have name 'zh'")
	}
	if len(langs) <= 2 {
		t.Errorf("Error: FindAllByBCP47(zh-TW) should have > 2 languages")
	}
}

func TestFindAllByBCP47ChineseTraditionalBestMatchWithCaseIgnoranceAndUnderscore(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	langs := lp.FindAllByBCP47("Zh_tW")
	if langs[0].BCP47 != "zh-TW" {
		t.Errorf("Error: FindAllByBCP47(ZH_TW)[0] should be 'zh-TW'")
	}
	if langs[1].BCP47 != "zh" {
		t.Errorf("Error: FindAllByBCP47(ZH_TW)[1] should have name 'zh'")
	}
	if len(langs) <= 2 {
		t.Errorf("Error: FindAllByBCP47(ZH_TW) should have > 2 languages")
	}
}

func TestNotFoundByBCP47(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	langs := lp.FindAllByBCP47("invalid")
	if len(langs) != 0 {
		t.Errorf("Error: FindAllByBCP47(invalid) should have 0 languages")
	}
}

func TestInvalidWinID(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	lang := lp.FindByWinID("invalid")
	if lang != nil {
		t.Errorf("Error: FindByWinID(invalid) should have nil language")
	}

	lang = lp.FindByWinID("zzz")
	if lang != nil {
		t.Errorf("Error: FindByWinID(zzz) should have nil language")
	}

	lang = lp.FindByWinID("_+#")
	if lang != nil {
		t.Errorf("Error: FindByWinID(_+#) should have nil language")
	}

	lang = lp.FindByWinID("000")
	if lang != nil {
		t.Errorf("Error: FindByWinID(000) should have nil language")
	}
}

func TestAfarInvalidWinID(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	lang := lp.Parse("aa")
	if lang.BCP47 != "aa" {
		t.Errorf("Error: FindByWinID(aa) should be 'aa'")
	}

	if lang.IsValidWinID() {
		t.Errorf("Error: FindByWinID(aa).IsValidWinID() should be false")
	}
}

func TestParseInvalid(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	lang := lp.Parse("invalid")
	if lang != nil {
		t.Errorf("Error: Parse(invalid) should be nil")
	}
}

func TestParseEnglish(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	lang := lp.Parse("en-US")
	if lang.BCP47 != "en-US" {
		t.Errorf("Error: Parse(en-US) should be 'en-US'")
	}
	lang = lp.Parse("enu")
	if lang.BCP47 != "en" {
		t.Errorf("Error: %v", lang)
		t.Errorf("Error: Parse(eng) should be 'en'")
	}
	lang = lp.Parse("eNg")
	if lang.BCP47 != "en" {
		t.Errorf("Error: %v", lang)
		t.Errorf("Error: Parse(eng) should be 'en'")
	}
}

func TestParseChinese(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	lang := lp.Parse("zh-CN")
	if lang.BCP47 != "zh-CN" {
		t.Errorf("Error: Parse(zh-CN) should be 'zh-CN'")
	}
	lang = lp.Parse("CHS")
	if lang.BCP47 != "zh" {
		t.Errorf("Error: Parse(CHS) should be 'zh'")
	}
	lang = lp.Parse("zHo")
	if lang.BCP47 != "zh" {
		t.Errorf("Error: Parse(ZHO) should be 'zh'")
	}
	lang = lp.Parse("WuU")
	if lang.BCP47 != "zh" {
		t.Errorf("Error: Parse(WUU) should be 'zh'")
	}
}

func TestFindByISO639Set1(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	lang := lp.FindByISO639Set1("en")
	if lang.BCP47 != "en" || lang.ISO639Set1 != "en" {
		t.Errorf("Error: FindByISO639Set1(en) should be {en, en}")
	}

	lang = lp.FindByISO639Set1("FR")
	if lang.BCP47 != "fr" || lang.ISO639Set1 != "fr" {
		t.Errorf("Error: FindByISO639Set1(zh) should be {zh, zh}")
	}

	lang = lp.FindByISO639Set1("cmn")
	if lang != nil {
		t.Errorf("Error: FindByISO639Set1(cmn) should be nil")
	}

	lang = lp.FindByISO639Set1("zh-Hans")
	if lang != nil {
		t.Errorf("Error: FindByISO639Set1(zh-Hans) should be nil")
	}
}

func TestFindByISO639Set2(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	lang := lp.FindByISO639Set2("aze")
	if lang.BCP47 != "az" || lang.ISO639Set2 != "aze" {
		t.Errorf("Error: FindByISO639Set2(aze) should be {az, aze}")
	}

	lang = lp.FindByISO639Set2("az")
	if lang != nil {
		t.Errorf("Error: FindByISO639Set2(az) should be nil")
	}

	lang = lp.FindByISO639Set2("cmn")
	if lang != nil {
		t.Errorf("Error: FindByISO639Set2(cmn) should be nil")
	}
}

func TestFindByISO639Set3(t *testing.T) {
	lp, err := slang.NewParser()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	lang := lp.FindByISO639Set3("ara")
	if lang.BCP47 != "ar" || lang.ISO639Set3 != "ara" {
		t.Errorf("Error: FindByISO639Set3(ara) should be 'ar'")
	}

	lang = lp.FindByISO639Set3("YuE")
	if lang.BCP47 != "yue" || lang.ISO639Set3 != "yue" {
		t.Errorf("Error: FindByISO639Set3(yue) should be 'yue'")
	}

	lang = lp.FindByISO639Set3("aZj")
	if lang.BCP47 != "az" || lang.ISO639Set3 != "azj" {
		t.Errorf("Error: FindByISO639Set3(azj) should be 'azj'")
	}
}
