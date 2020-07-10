package i18n

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestLoadXliffFromFile(t *testing.T) {
	doc, err := LoadXliffFromFile("./test.xliff")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("version:", doc.Version)
	fmt.Println("xmlns:", doc.Xmlns)

	fmt.Println("original:", doc.Files[0].Original)
	fmt.Println("source lang:", doc.Files[0].SourceLang)
	fmt.Println("data type:", doc.Files[0].DataType)
	fmt.Println("target lang:", doc.Files[0].TargetLang)

	fmt.Println("tool id:", doc.Files[0].Header.Tool.ToolId)
	fmt.Println("tool name:", doc.Files[0].Header.Tool.ToolName)
	fmt.Println("tool version:", doc.Files[0].Header.Tool.ToolVersion)
	fmt.Println("build num:", doc.Files[0].Header.Tool.BuildNum)

	fmt.Println("id:", doc.Files[0].Body.TransUnit[0].Id)
	fmt.Println("approved:", doc.Files[0].Body.TransUnit[0].Approved)
	fmt.Println("source inner:", doc.Files[0].Body.TransUnit[0].Source.Inner)
	fmt.Println("source lang:", doc.Files[0].Body.TransUnit[0].Source.Lang)

	fmt.Println("target inner:", doc.Files[0].Body.TransUnit[0].Target.Inner)
	fmt.Println("target lang:", doc.Files[0].Body.TransUnit[0].Target.Lang)
}

func TestXliffMarshal(t *testing.T) {
	doc := XliffDoc{
		Version: "1.1",
		Xmlns:   "test",
		Files: []XliffFile{
			XliffFile{
				Original:   "1",
				SourceLang: "en",
				DataType:   "text",
				TargetLang: "fr",
				Header: XliffFileHeader{
					Tool: XliffFileHeaderInfo{
						ToolId:      "1",
						ToolName:    "sublime",
						ToolVersion: "3",
						BuildNum:    "2",
					},
				},
				Body: XliffFileBody{
					TransUnit: []XliffFileBodyTransUnit{
						XliffFileBodyTransUnit{
							Id:       "about",
							Approved: "yes",
							Source: XliffFileBodyTransUnitInner{
								Inner: "this is",
							},
							Target: XliffFileBodyTransUnitInner{
								Inner: "aaa",
								Lang:  "fr",
							},
						},
					},
				},
			},
		},
	}

	if buf, err := xml.Marshal(doc); err != nil {
		t.Error(err)
	} else {
		fmt.Println(string(buf))
	}
}

func TestGetLanguage(t *testing.T) {
	name, err := GetLanguage("./abc.xliff")
	if err != nil {
		t.Error(err)
	} else if name != "abc" {
		t.Error("wrong answer")
	}

	name, err = GetLanguage("./data/abc.xliff")
	if err != nil {
		t.Error(err)
	} else if name != "abc" {
		t.Error("wrong answer")
	}

	name, err = GetLanguage("./abc.xli")
	if err == nil {
		t.Error("wrong answer")
	}

	name, err = GetLanguage("./abcxli")
	if err == nil {
		t.Error("wrong answer")
	}

	name, err = GetLanguage("")
	if err == nil {
		t.Error("wrong answer")
	}
}

func TestGetTranslation(t *testing.T) {
	dict, err := GetTranslation("./test.xliff")
	if err != nil {
		t.Error(err)
	}

	if dict["CFBundleName"] != "Air Matters - Environs" {
		t.Error("wrong answer")
	}
}

func TestT(t *testing.T) {
	translator := XliffTranslator{}
	if err := translator.Load("./data"); err != nil {
		t.Error(err)
	}

	if translator.T("CFBundleName", "fr") != "Air Matters - Environs" {
		t.Error("wrong answer")
	}

	if translator.T("CFBundleName", "ch") != "" {
		t.Error("wrong answer")
	}

	if translator.T("C", "en") != "this is c" {
		t.Error("wrong answer")
	}
}
