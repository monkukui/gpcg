package gpt_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"gpt"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, gpt.Analyzer, "a")
}
