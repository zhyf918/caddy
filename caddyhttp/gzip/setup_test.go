package gzip

import (
	"testing"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("http", `gzip`)
	err := setup(c)
	if err != nil {
		t.Errorf("Expected no errors, but got: %v", err)
	}
	mids := httpserver.GetConfig(c).Middleware()
	if mids == nil {
		t.Fatal("Expected middleware, was nil instead")
	}

	handler := mids[0](httpserver.EmptyNext)
	myHandler, ok := handler.(Gzip)
	if !ok {
		t.Fatalf("Expected handler to be type Gzip, got: %#v", handler)
	}

	if !httpserver.SameNext(myHandler.Next, httpserver.EmptyNext) {
		t.Error("'Next' field of handler was not set properly")
	}

	tests := []struct {
		input     string
		shouldErr bool
	}{
		{`gzip {`, true},
		{`gzip {}`, true},
		{`gzip a b`, true},
		{`gzip a {`, true},
		{`gzip { not f } `, true},
		{`gzip { not } `, true},
		{`gzip { not /file
		 ext .html
		 level 1
		} `, false},
		{`gzip { level 9 } `, false},
		{`gzip { ext } `, true},
		{`gzip { ext /f
		} `, true},
		{`gzip { not /file
		 ext .html
		 level 1
		}
		gzip`, false},
		{`gzip {
		 ext ""
		}`, false},
		{`gzip { not /file
		 ext .html
		 level 1
		}
		gzip { not /file1
		 ext .htm
		 level 3
		}
		`, false},
		{`gzip { not /file
		 ext .html
		 level 1
		}
		gzip { not /file1
		 ext .htm
		 level 3
		}
		`, false},
		{`gzip { not /file
		 ext *
		 level 1
		}
		`, false},
		{`gzip { not /file
		 ext *
		 level 1
		 min_length ab
		}
		`, true},
		{`gzip { not /file
		 ext *
		 level 1
		 min_length 1000
		}
		`, false},
	}
	for i, test := range tests {
		_, err := gzipParse(caddy.NewTestController("http", test.input))
		if test.shouldErr && err == nil {
			t.Errorf("Test %v: Expected error but found nil", i)
		} else if !test.shouldErr && err != nil {
			t.Errorf("Test %v: Expected no error but found error: %v", i, err)
		}
	}
}
