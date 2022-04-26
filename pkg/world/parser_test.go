package world

import (
	"strings"
	"testing"
)

type testLinkMap map[string]map[Direction]string

func TestParse(t *testing.T) {
	var tests = []struct {
		input        string
		parseError   bool
		nAliens      int
		wantedCities []string
		wantedLinks  testLinkMap
	}{
		// Successful parsing
		{
			input: `A north=B
							B south=C`,
			parseError:   false,
			nAliens:      0,
			wantedCities: []string{"A", "B", "C"},
			wantedLinks: testLinkMap{
				"A": map[Direction]string{North: "B"},
				"B": map[Direction]string{South: "C"},
			},
		},
		// Failed direction parsing
		{
			input:        `A north=B asd=C`,
			parseError:   true,
			nAliens:      0,
			wantedCities: []string{"A"},
			wantedLinks:  testLinkMap{},
		},
		// Failed destination parsing
		{
			input:        "this is not a valid destination string",
			parseError:   true,
			nAliens:      0,
			wantedCities: []string{},
			wantedLinks:  testLinkMap{},
		},
		// Empty world (no cities)
		{
			input:        "",
			parseError:   false,
			nAliens:      0,
			wantedCities: []string{},
			wantedLinks:  testLinkMap{},
		},
		// Directions are overridden on the same line
		{
			input:        "A north=B north=C",
			parseError:   false,
			nAliens:      0,
			wantedCities: []string{"A", "B", "C"},
			wantedLinks: testLinkMap{
				"A": map[Direction]string{North: "C"},
			},
		},
		// Aliens are correctly deployed
		{
			input: `A north=B
							B south=C`,
			parseError:   false,
			nAliens:      4,
			wantedCities: []string{"A", "B", "C"},
			wantedLinks:  testLinkMap{},
		},
	}

	for _, test := range tests {
		w, err := Parse(
			strings.NewReader(test.input),
			test.nAliens,
		)

		if test.parseError {
			if err == nil {
				t.Errorf("Expected error, got nil")
				continue
			} else {
				continue
			}
		}

		if test.nAliens != len(w.Aliens) {
			t.Errorf("Expected %d aliens, got %d", test.nAliens, len(w.Aliens))
			continue
		}

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
			continue
		}

		for _, city := range test.wantedCities {
			if _, exists := w.Cities[city]; !exists {
				t.Errorf("Expected city %s to exist", city)
				continue
			}
		}

		for _, city := range test.wantedCities {
			if links, exists := test.wantedLinks[city]; exists {
				for direction, arrival := range links {
					strDirection, _ := directionString(direction)

					if _, exists := w.Links[city][direction]; !exists {
						t.Errorf(
							"City %s should have links to %s, but it doesn't",
							city,
							strDirection,
						)
						continue
					}

					if _, exists := w.Links[city][direction]; !exists {
						t.Errorf(
							"Expected link %s->%s to exist",
							city,
							arrival,
						)
						continue
					}
					if w.Links[city][direction].Name != arrival {
						t.Errorf(
							"Expected link %s=%s, got %s",
							strDirection,
							arrival,
							w.Links[city][direction].Name,
						)
					}
				}
			}
		}
	}
}
