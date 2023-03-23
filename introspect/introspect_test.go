package introspect_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stefanprifti/gqlclientgen/introspect"
)

func TestIntrospect(t *testing.T) {
	cases := []struct {
		name     string
		url      string
		fileName string
	}{
		{
			name:     "brotforce",
			url:      "https://brotforce-bff-staging.mcmakler.com/",
			fileName: "./testdata/brotforce-bff-staging.mcmakler.com.graphql",
		},
		{
			name:     "countries",
			url:      "https://countries.trevorblades.com/graphql",
			fileName: "./testdata/countries.trevorblades.com.graphql",
		},
		{
			name:     "swapi",
			url:      "https://swapi-graphql.netlify.app/.netlify/functions/index",
			fileName: "./testdata/swapi-graphql.netlify.app.graphql",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			schema, err := introspect.URL(tc.url)
			if err != nil {
				t.Fatal(err)
			}

			txt, err := introspect.SchemaToText(schema)
			if err != nil {
				t.Fatal(err)
			}

			f, err := os.OpenFile(tc.fileName, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			expectedTxt, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			if string(txt) != string(expectedTxt) {
				t.Fatalf("expected %s, got %s", expectedTxt, txt)
			}

			// _, err = f.Write([]byte(txt))
			// if err != nil {
			// 	t.Fatal(err)
			// }
		})
	}
}
