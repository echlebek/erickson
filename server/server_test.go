package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/echlebek/erickson/db"
	"github.com/echlebek/erickson/diff"
	"github.com/echlebek/erickson/review"
)

const diffTxt = `diff --git a/data_cleanup.go b/data_cleanup.go
index c45cb07..e9b8ff6 100644
--- a/data_cleanup.go
+++ b/data_cleanup.go
@@ -1,60 +1,58 @@
 package main

 import (
        "bufio"
        "flag"
        "fmt"
        "io"
+       "log"
        "os"
 )

 func main() {
        flag.Parse()
        args := flag.Args()
        if len(args) != 2 {
                fmt.Println("usage: ./data_cleanup input.txt output.csv")
                os.Exit(1)
        }
        inFile, err := os.Open(args[0])
        defer inFile.Close()
        if err != nil {
-               fmt.Println(err)
-               os.Exit(1)
+               log.Fatal(err)
        }
        rd := bufio.NewReader(inFile)

        outFile, err := os.Create(args[1])
        defer outFile.Close()
        if err != nil {
-               fmt.Println(err)
-               os.Exit(1)
+               log.Fatal(err)
        }
        wr := bufio.NewWriter(outFile)

        var r, lastR rune

        for err == nil {
                r, _, err = rd.ReadRune()
                switch r {
                case '"':
                        switch lastR {
                        case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
                                wr.WriteString(" inch")
                        }
                case '~':
                case '^':
                        wr.WriteByte('\t')
                default:
                        wr.WriteRune(r)
                }

                lastR = r
        }

        if err != io.EOF {
-               fmt.Println(err)
-               os.Exit(1)
+               log.Fatal(err)
        }

        wr.Flush()
 }
diff --git a/import.sql b/import.sql
index 90e913e..c31da41 100644
--- a/import.sql
+++ b/import.sql
@@ -1,18 +1,18 @@
 BEGIN;

     SET CONSTRAINTS ALL DEFERRED;

-    COPY food_groups FROM '/Users/eric/code/nutr/sr26/fd_group.csv' DELIMITER '        ' CSV;
-    COPY foods FROM '/Users/eric/code/nutr/sr26/food_des.csv' DELIMITER '      ' CSV;
-    COPY langua_l_desc FROM '/Users/eric/code/nutr/sr26/langdesc.csv' DELIMITER '      ' CSV;
-    COPY langua_l_factors FROM '/Users/eric/code/nutr/sr26/langual.csv' DELIMITER '    ' CSV;
-    COPY nutrients FROM '/Users/eric/code/nutr/sr26/nutr_def.csv' DELIMITER '  ' CSV;
-    COPY source_codes FROM '/Users/eric/code/nutr/sr26/src_cd.csv' DELIMITER ' ' CSV;
-    COPY data_derivation_codes FROM '/Users/eric/code/nutr/sr26/deriv_cd.csv' DELIMITER '      ' CSV;
-    COPY nutrient_data FROM '/Users/eric/code/nutr/sr26/nut_data.csv' DELIMITER '      ' CSV;
-    COPY weights FROM '/Users/eric/code/nutr/sr26/weight.csv' DELIMITER '      ' CSV;
-    COPY footnotes FROM '/Users/eric/code/nutr/sr26/footnote.csv' DELIMITER '  ' CSV;
-    COPY sources_of_data FROM '/Users/eric/code/nutr/sr26/data_src.csv' DELIMITER '    ' CSV;
-    COPY sources_of_data_assoc FROM '/Users/eric/code/nutr/sr26/datsrcln.csv' DELIMITER '      ' CSV;
+    \copy food_groups FROM 'fd_group.csv' DELIMITER '  ' CSV;
+    \copy foods FROM 'food_des.csv' DELIMITER '        ' CSV;
+    \copy langua_l_desc FROM 'langdesc.csv' DELIMITER '        ' CSV;
+    \copy langua_l_factors FROM 'langual.csv' DELIMITER '      ' CSV;
+    \copy nutrients FROM 'nutr_def.csv' DELIMITER '    ' CSV;
+    \copy source_codes FROM 'src_cd.csv' DELIMITER '   ' CSV;
+    \copy data_derivation_codes FROM 'deriv_cd.csv' DELIMITER '        ' CSV;
+    \copy nutrient_data FROM 'nut_data.csv' DELIMITER '        ' CSV;
+    \copy weights FROM 'weight.csv' DELIMITER '        ' CSV;
+    \copy footnotes FROM 'footnote.csv' DELIMITER '    ' CSV;
+    \copy sources_of_data FROM 'data_src.csv' DELIMITER '      ' CSV;
+    \copy sources_of_data_assoc FROM 'datsrcln.csv' DELIMITER '        ' CSV;

 COMMIT;
`

var Do = http.DefaultTransport.RoundTrip

func TestServer(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "erickson")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	db, err := db.NewBoltDB(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	handler := NewRootHandler(db, wd+"./..")
	fmt.Println(os.Getwd())

	server := httptest.NewServer(handler)
	defer server.Close()

	*handler.URL = server.URL
	url := server.URL

	url = create(t, url)
	read(t, url)
	annotate(t, url+"/rev/0")
	update(t, url)
	destroy(t, url)
}

func create(t *testing.T, url string) string {
	files, err := diff.ParseFiles(diffTxt)
	if err != nil {
		t.Fatal(err)
	}
	r := review.R{
		Summary:   review.Summary{},
		Revisions: []review.Revision{{Files: files}},
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", url+"/reviews", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 303 {
		msg, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("bad response: code %d, %s", resp.StatusCode, msg)
	}
	expURL := url + "/reviews/1"
	if resp.Header.Get("Location") != expURL {
		t.Errorf("wrong url. got %s, want %s", resp.Header.Get("Location"), expURL)
	}
	return expURL
}

func read(t *testing.T, url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) == 0 {
		t.Error("empty response")
	}
	return body
}

func update(t *testing.T, url string) {

	//req, err := http.NewRequest("PATCH", url, body)
}

func annotate(t *testing.T, url string) {
	anno := review.Annotation{
		File:    0,
		Hunk:    0,
		Line:    123,
		Comment: "I don't like this line",
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&anno); err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PATCH", url, &buf)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("couldn't annotate review: %d", resp.StatusCode)
	}
}

func destroy(t *testing.T, url string) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("bad response code %d", resp.StatusCode)
	}
}
