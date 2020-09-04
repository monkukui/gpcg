package gpt_test

import (
	"bytes"
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

	generatedFile, err := os.Open("./gen/gen.go")
	defer generatedFile.Close()
	if err != nil {
		t.Fatal(err)
	}
	expectedFile, err := os.Open("./testdata/src/a/expected/expected.go")
	defer expectedFile.Close()
	if err != nil {
		t.Fatal(err)
	}

	// 生成されたコードと期待値が一致しているかをチェック
	generatedCode, err := ioutil.ReadAll(generatedFile)
	if err != nil {
		t.Fatal(err)
	}
	expectedCode, err := ioutil.ReadAll(expectedFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(generatedCode, expectedCode) {
		t.Fatal("generated code is different from expected code")
	}
}
