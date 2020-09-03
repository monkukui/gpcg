package gpt_test

import (
	"fmt"
	"os"
	"testing"

	"io/ioutil"

	"golang.org/x/tools/go/analysis/analysistest"
	"gpt"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, gpt.Analyzer, "a")
}

func TestCodeGenerate(t *testing.T) {

	generatedFile, err := os.Open("./gen/gen.go")
	defer generatedFile.Close()
	if err != nil {
		t.Error(err)
	}
	expectedFile, err := os.Open("./gen/expected.go")
	defer expectedFile.Close()
	if err != nil {
		t.Error(err)
	}

	// 生成されたコードと期待値が一致しているかをチェック
	generatedCode, err := ioutil.ReadAll(generatedFile)
	if err != nil {
		t.Error(err)
	}
	expectedCode, err := ioutil.ReadAll(expectedFile)
	if err != nil {
		t.Error(err)
	}

	if string(generatedCode) != string(expectedCode) {
		t.Error("generated code is different from expected code")
	}
}
