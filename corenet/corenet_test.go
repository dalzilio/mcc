package corenet

import (
	"bufio"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/dalzilio/mcc/hlnet"
	"github.com/dalzilio/mcc/pnml"
)

type testPNML struct {
	pl int // expected nuber of places
	tr int // expected nuber of transitions
}

func TestBuild(t *testing.T) {
	// Read the csv file with the expected number of places and trans for each
	// model
	csvFile, _ := os.Open("../benchmarks/models.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comment = '#'
	// Populate two maps to test the results
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	expected := make(map[string]testPNML)
	for _, r := range records {
		pl, _ := strconv.Atoi(r[4])
		tr, _ := strconv.Atoi(r[5])
		expected[r[0]] = testPNML{pl: pl, tr: tr}
	}
	// Iterate through all the PNML file in the benchmarks folder
	directory := "../benchmarks/"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		t.Errorf("corenet.Build(): error opening the benchmarks folder (%s)", directory)
	}

	for _, file := range files {
		filename := file.Name()
		if !file.IsDir() && filepath.Ext(filename) == ".pnml" {
			xmlFile, err := os.Open(filepath.Join(directory, filename))
			if err != nil {
				t.Errorf("corenet.Build(): error opening file: %s", err)
			}
			defer xmlFile.Close()

			// Decode the XML file and build a pnml.Net
			decoder := pnml.NewDecoder(xmlFile)
			var p = new(pnml.Net)
			err = decoder.Build(p)
			if err != nil {
				t.Errorf("corenet.Build(): error decoding PNML file: %s", err)
			}

			// Build a hlnet.Net object from it
			p.SetVerbose(pnml.MINIMAL)
			p.SetFES(false)
			hl := hlnet.Build(p)

			// Then finally build a corenet.Net
			cn := Build(p, hl)

			// Check that the metrics are correct
			filename = filename[0 : len(filename)-len(".pnml")]
			tt, present := expected[filename]
			if !present {
				t.Errorf("corenet.Build(): error file %s.pnml is not listed in the benchmarks.csv files", filename)
			}
			if ppl := len(cn.pl); tt.pl != ppl {
				t.Errorf("corenet.Build(): error model in file %s.pnml has not the right number of places (has %d instead of %d)", filename, ppl, tt.pl)
			}
			if ttr := len(cn.tr); tt.tr != ttr {
				t.Errorf("corenet.Build(): error model in file %s.pnml has not the right number of transitions (has %d instead of %d)", filename, ttr, tt.tr)
			}
		}
	}
}

var result string

func BenchmarkBuild(b *testing.B) {
	// Find all the PNML file in the benchmarks folder
	directory := "../benchmarks/"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		os.Exit(1)
	}

	nets := make([]*pnml.Net, 0)
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".pnml" {
			xmlFile, err := os.Open(filepath.Join(directory, file.Name()))
			if err != nil {
				os.Exit(1)
			}
			decoder := pnml.NewDecoder(xmlFile)
			var p = new(pnml.Net)
			_ = decoder.Build(p)
			p.SetVerbose(pnml.MINIMAL)
			p.SetFES(false)
			nets = append(nets, p)
			xmlFile.Close()
		}
	}

	for n := 0; n < b.N; n++ {
		for _, p := range nets {
			hl := hlnet.Build(p)
			cn := Build(p, hl)
			result = cn.name
		}
	}
}
