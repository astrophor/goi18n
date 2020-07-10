package i18n

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type XliffDoc struct {
	XMLName xml.Name    `xml:"xliff"`
	Version string      `xml:"version,attr"`
	Xmlns   string      `xml:"xmlns,attr"`
	Files   []XliffFile `xml:"file"`
}

type XliffFile struct {
	XMLName    xml.Name        `xml:"file"`
	Original   string          `xml:"original,attr"`
	SourceLang string          `xml:"source-language,attr,omitempty"`
	DataType   string          `xml:"datatype,attr,omitempty"`
	TargetLang string          `xml:"target-language,attr,omitempty"`
	Header     XliffFileHeader `xml:"header"`
	Body       XliffFileBody   `xml:"body"`
}

type XliffFileHeader struct {
	XMLName xml.Name            `xml:"header"`
	Tool    XliffFileHeaderInfo `xml:"tool"`
}

type XliffFileHeaderInfo struct {
	XMLName     xml.Name `xml:"tool"`
	ToolId      string   `xml:"tool-id,attr,omitempty"`
	ToolName    string   `xml:"tool-name,attr,omitempty"`
	ToolVersion string   `xml:"tool-version,attr,omitempty"`
	BuildNum    string   `xml:"build-num,attr,omitempty"`
}

type XliffFileBody struct {
	XMLName   xml.Name                 `xml:"body"`
	TransUnit []XliffFileBodyTransUnit `xml:"trans-unit"`
}

type XliffFileBodyTransUnit struct {
	XMLName  xml.Name                    `xml:"trans-unit"`
	Id       string                      `xml:"id,attr"`
	Approved string                      `xml:"approved,attr,omitempty"`
	Source   XliffFileBodyTransUnitInner `xml:"source"`
	Target   XliffFileBodyTransUnitInner `xml:"target"`
}

type XliffFileBodyTransUnitInner struct {
	Inner string `xml:",chardata"`
	Lang  string `xml:"lang,attr,omitempty"`
}

func LoadXliffFromFile(file_name string) (*XliffDoc, error) {
	fd, err := os.Open(file_name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	doc := new(XliffDoc)
	dec := xml.NewDecoder(fd)
	if err := dec.Decode(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

type XliffTranslator struct {
	data map[string]map[string]string
}

func (x *XliffTranslator) Load(file_path string) error {
	file_path, err := CheckFilePath(file_path)
	if err != nil {
		return err
	}

	file_list, err := GetFileList(file_path)
	if err != nil {
		return err
	}

	x.data = make(map[string]map[string]string)

	for _, file_name := range file_list {
		lang, err := GetLanguage(file_name)
		if err != nil {
			return err
		}

		dict, err := GetTranslation(file_name)
		if err != nil {
			return err
		}

		x.data[lang] = dict
	}

	return nil
}

func CheckFilePath(file_path string) (string, error) {
	fi, err := os.Stat(file_path)
	if err != nil {
		return "", err
	}

	if !fi.IsDir() {
		return "", fmt.Errorf("path %v is not a directory", file_path)
	}

	return file_path, nil
}

func GetFileList(file_path string) ([]string, error) {
	files, err := ioutil.ReadDir(file_path)
	if err != nil {
		return nil, err
	}

	var file_list []string
	for _, file := range files {
		file_name := file_path + "/" + file.Name()
		file_list = append(file_list, file_name)
	}

	return file_list, nil
}

func GetLanguage(file_name string) (string, error) {
	items := strings.Split(file_name, "/")
	if len(items) == 0 {
		return "", fmt.Errorf("invalid file name: %v\n", file_name)
	}

	parts := strings.Split(items[len(items)-1], ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid file name: %v\n", file_name)
	} else if parts[1] != "xliff" {
		return "", fmt.Errorf("invalid file name: %v\n", file_name)
	}

	return parts[0], nil
}

func GetTranslation(file_name string) (map[string]string, error) {
	lang, err := GetLanguage(file_name)
	if err != nil {
		return nil, err
	}

	doc, err := LoadXliffFromFile(file_name)
	if err != nil {
		return nil, err
	}

	dict := make(map[string]string)

	for file_idx, _ := range doc.Files {
		for unit_idx, _ := range doc.Files[file_idx].Body.TransUnit {
			if lang == "en" {
				dict[doc.Files[file_idx].Body.TransUnit[unit_idx].Id] = doc.Files[file_idx].Body.TransUnit[unit_idx].Source.Inner
			} else {
				dict[doc.Files[file_idx].Body.TransUnit[unit_idx].Id] = doc.Files[file_idx].Body.TransUnit[unit_idx].Target.Inner
			}
		}
	}

	return dict, nil
}

func (x *XliffTranslator) T(id, lang string) string {
	mapping, ok := x.data[lang]
	if !ok {
		log.Printf("language %v is not supported\n", lang)
		mapping = x.data["en"]
	}

	return mapping[id]
}
