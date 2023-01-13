// Copyright 2020. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"

	"github.com/dalzilio/mcc/pnml"
)

// ----------------------------------------------------------------------

// We use the functionality for Marshalling Go structs into XML to implement the
// generation on PNML files.

const (
	// DOCTYPE for the generated PNML file
	DOCTYPE = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

// wpnml is the type of PNML for a P/T net without graphical information
type wpnml struct {
	XMLName xml.Name `xml:"http://www.pnml.org/version-2009/grammar/pnml pnml"`
	WNET    wnet     `xml:"net"`
}

// wnet is the type of PNML net. We ignore the graphical information contained
// in the net and only consider the first page
type wnet struct {
	Thetype string `xml:"type,attr"`
	ID      string `xml:"id,attr"`
	NAME    string `xml:"name>text"`
	PAGE    wpage  `xml:"page"`
}

type wpage struct {
	ID     string        `xml:"id,attr"`
	PLACES []*Place      `xml:"place"`
	TRANS  []*Trans      `xml:"transition"`
	TOOLS  *toolspecific `xml:"toolspecific,omitempty"`
}

// toolspecific is the type of toolspecific XML annotation added to the PNML
// output to declare association lists between place and transitions in the COL,
// input net and  their counterparts in the output, P/T net
type toolspecific struct {
	TheTool    string    `xml:"tool,attr"`
	TheVersion string    `xml:"version,attr"`
	HLPL       []hlassoc `xml:"places>place"`
	HLTR       []hlassoc `xml:"transitions>transition"`
}

type hlassoc struct {
	ID    string   `xml:"id,attr"`
	NAMES []string `xml:"name"`
}

// ----------------------------------------------------------------------

// MarshalXML encodes the receiver as zero or more XML elements. This makes
// Place a xml.Marshaller
func (v Place) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{{Name: xml.Name{Local: "id"}, Value: v.name}}
	e.EncodeToken(start)
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "name"}})
	e.EncodeElement(v.name, xml.StartElement{Name: xml.Name{Local: "text"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "name"}})
	if v.init != 0 {
		e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "initialMarking"}})
		e.EncodeElement(v.init, xml.StartElement{Name: xml.Name{Local: "text"}})
		e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "initialMarking"}})
	}
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

// MarshalXML encodes the receiver as zero or more XML elements. This makes
// Trans a xml.Marshaller
func (v Trans) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	name := fmt.Sprintf("t%d", v.count)
	start.Attr = []xml.Attr{{Name: xml.Name{Local: "id"}, Value: name}}
	e.EncodeToken(start)
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "name"}})
	e.EncodeElement(name, xml.StartElement{Name: xml.Name{Local: "text"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "name"}})
	e.EncodeToken(xml.EndElement{Name: start.Name})

	for _, c := range v.in {
		encodeArc(e, fmt.Sprintf("p2t-%d-%d", v.count, c.count), c.Place.name, name, c.int)
	}
	for _, c := range v.out {
		encodeArc(e, fmt.Sprintf("t2p-%d-%d", c.count, v.count), name, c.Place.name, c.int)
	}

	return nil
}

func encodeArc(e *xml.Encoder, id, src, tgt string, weight int) {
	arc := xml.StartElement{
		Name: xml.Name{Local: "arc"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "id"}, Value: id},
			{Name: xml.Name{Local: "source"}, Value: src},
			{Name: xml.Name{Local: "target"}, Value: tgt},
		},
	}
	e.EncodeToken(arc)
	if weight != 1 {
		e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "inscription"}})
		e.EncodeElement(weight, xml.StartElement{Name: xml.Name{Local: "text"}})
		e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "inscription"}})
	}
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "arc"}})
}

// ----------------------------------------------------------------------

// PnmlWrite prints a P/T net in PNML format on an io.Writer
func (net Net) PnmlWrite(w io.Writer) error {
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")

	// we start by sorting the slice of places. In the case where the result is
	// not "sliced", place names are all of the form p_k, with k an integer, and
	// are already sorted, so we can just do nothing.
	if net.verbose != pnml.SMPT {
		sort.Slice(net.pl, func(i, j int) bool {
			return net.pl[i].name < net.pl[j].name
		})
	}

	var tools *toolspecific

	// we print out properties. We use the fact that places are sorted by names.
	// Hence (core) places corresponding to the same colored place are grouped
	// together. Same for transitions.
	if net.printprops {
		tools = &toolspecific{TheTool: "mcc-unfolder", TheVersion: "2.0"}
		assoc := hlassoc{ID: ""}
		// output list of places for each colored one
		for _, v := range net.pl {
			if assoc.ID != v.label {
				// if this is not the first assoc list we need to push the
				// previous one in tools
				if assoc.ID != "" {
					tools.HLPL = append(tools.HLPL, assoc)
				}
				assoc = hlassoc{ID: v.label}
			}
			assoc.NAMES = append(assoc.NAMES, v.name)
		}
		// we should also push the last one
		if assoc.ID != "" {
			tools.HLPL = append(tools.HLPL, assoc)
		}
		// we do the same with the list of transitions
		assoc = hlassoc{ID: ""}
		for k, v := range net.tr {
			if assoc.ID != v.label {
				if assoc.ID != "" {
					tools.HLTR = append(tools.HLTR, assoc)
				}
				assoc = hlassoc{ID: v.label}
			}
			assoc.NAMES = append(assoc.NAMES, fmt.Sprintf("t%d", k))
		}
		if assoc.ID != "" {
			tools.HLTR = append(tools.HLTR, assoc)
		}
	}

	// Now we output the file on the io.Writer
	wpnml := wpnml{
		WNET: wnet{
			Thetype: "http://www.pnml.org/version-2009/grammar/ptnet",
			ID:      net.name,
			NAME:    "MCC-PT-" + net.name,
			PAGE: wpage{
				ID:     "page",
				PLACES: net.pl,
				TRANS:  net.tr,
				TOOLS:  tools,
			},
		},
	}
	w.Write([]byte(DOCTYPE))
	return encoder.Encode(wpnml)
}
