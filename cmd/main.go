package main

import (
	"os"

	"github.com/com-gft-tsbo-source/go-fe-measure/femeasure"
)

// ###########################################################################
// ###########################################################################
// MAIN
// ###########################################################################
// ###########################################################################

var usage []byte = []byte("fe-measure: [OPTIONS] ")

func main() {

	var ms fe-measure.FeMeasure

	fe-measure.InitFromArgs(&ms, os.Args, nil)

	ms.Run()
}
