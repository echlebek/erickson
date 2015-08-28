package main

import (
	"log"
	"net/http"
	"time"

	"github.com/echlebek/erickson/db"
	"github.com/echlebek/erickson/server"
)

const diff = `diff --git a/data_cleanup.go b/data_cleanup.go
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

func main() {
	db, err := db.NewBoltDB("my2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler := server.NewRootHandler(db, ".")

	s := http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
