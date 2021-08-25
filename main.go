package main

import (
	"github.com/sreesindhusruthiyadavalli/customSA/structsa"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(structsa.Analyzer) }
