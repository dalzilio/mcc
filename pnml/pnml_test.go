package pnml

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

type TestPNML struct {
	n  string // file name
	pl int    // expected nuber of places
	tr int    // expected nuber of transitions
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
	expected := make(map[string]TestPNML)
	for _, r := range records {
		pl, _ := strconv.Atoi(r[2])
		tr, _ := strconv.Atoi(r[3])
		expected[r[0]] = TestPNML{n: r[1], pl: pl, tr: tr}
	}
	// Iterate through all the PNML file in the benchmarks folder
	directory := "../benchmarks/simple/"
	files, err := os.ReadDir(directory)
	if err != nil {
		t.Errorf("pnml.Build(): error opening the benchmarks folder (%s)", directory)
	}

	for _, file := range files {
		filename := file.Name()
		if !file.IsDir() && filepath.Ext(filename) == ".pnml" {
			xmlFile, err := os.Open(filepath.Join(directory, filename))
			if err != nil {
				t.Errorf("pnml.Build(): error opening file: %s", err)
			}
			defer xmlFile.Close()

			decoder := NewDecoder(xmlFile)
			var p = new(Net)
			err = decoder.Build(p)
			if err != nil {
				t.Errorf("pnml.Build(): error decoding PNML file %s: %s", filename, err)
			}

			filename = filename[0 : len(filename)-len(".pnml")]
			tt, present := expected[filename]
			if !present {
				t.Errorf("pnml.Build(): error file %s.pnml is not listed in the benchmarks.csv files", filename)
			}
			if tt.n != p.Name {
				t.Errorf("pnml.Build(): error model in file %s.pnml has not the right name (has %s instead of %s)", filename, p.Name, tt.n)
			}
			if ppl := len(p.Page.Places); tt.pl != ppl {
				t.Errorf("pnml.Build(): error model in file %s.pnml has not the right number of places (has %d instead of %d)", filename, ppl, tt.pl)
			}
			if ttr := len(p.Page.Trans); tt.tr != ttr {
				t.Errorf("pnml.Build(): error model in file %s.pnml has not the right number of transitions (has %d instead of %d)", filename, ttr, tt.tr)
			}
		}
	}
}

var result string

func BenchmarkBuild(b *testing.B) {
	// Find all the PNML file in the benchmarks folder
	directory := "../benchmarks/large/"
	files, err := os.ReadDir(directory)
	if err != nil {
		os.Exit(1)
	}

	descriptors := make([]*os.File, 0)
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".pnml" {
			xmlFile, err := os.Open(filepath.Join(directory, file.Name()))
			if err != nil {
				os.Exit(1)
			}
			defer xmlFile.Close()
			descriptors = append(descriptors, xmlFile)
		}
	}

	var p = new(Net)

	for n := 0; n < b.N; n++ {
		for _, xmlFile := range descriptors {
			decoder := NewDecoder(xmlFile)
			_ = decoder.Build(p)
			result = p.Name
		}
	}
}
