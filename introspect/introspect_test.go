package introspect_test

import (
	"os"
	"testing"

	"github.com/stefanprifti/gql/introspect"
)

func TestIntrospect(t *testing.T) {
	cases := []struct {
		name     string
		url      string
		fileName string
	}{
		{
			name:     "test",
			url:      "https://brotforce-bff-staging.mcmakler.com/",
			fileName: "./testdata/brotforce-bff-staging.mcmakler.com.graphql",
		},
		{
			name:     "test",
			url:      "https://google-maps-develop.mcmakler.com/query",
			fileName: "./testdata/google-maps-develop.mcmakler.com.graphql",
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

			// expectedTxt, err := ioutil.ReadAll(f)
			// if err != nil {
			// 	t.Fatal(err)
			// }

			// if txt != string(expectedTxt) {
			// 	t.Fatalf("expected %s, got %s", expectedTxt, txt)
			// }

			_, err = f.Write([]byte(txt))
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
