package db

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/echlebek/erickson/diff"
	"github.com/echlebek/erickson/review"
)

const mockDiff1 = `--- a/setup.py	Mon Aug 05 22:46:08 2013 -0700
+++ b/setup.py	Thu Mar 26 15:32:48 2015 -0700
@@ -30,7 +30,5 @@
     scripts=["scripts/fetta"],
     test_suite="tests.unit",
     cmdclass={"cram": CramTest},
-    install_requires=[
-        "numpy", "h5py"
-    ]
+    install_requires=["h5py"]
 )`

const mockDiff2 = `--- a/setup.py	Mon Aug 05 22:46:08 2013 -0700
+++ b/setup.py	Thu Mar 26 15:36:47 2015 -0700
@@ -9,6 +9,7 @@


 class CramTest(TestCommand):
+    """Runs all the cram tests."""

     def run(self):
         import cram
@@ -30,7 +31,5 @@
     scripts=["scripts/fetta"],
     test_suite="tests.unit",
     cmdclass={"cram": CramTest},
-    install_requires=[
-        "numpy", "h5py"
-    ]
+    install_requires=["h5py"]
 )`

var mockReview review.R

func init() {
	files, err := diff.ParseFiles(mockDiff1)
	if err != nil {
		panic(err)
	}
	mockReview = review.R{
		Summary: review.Summary{
			Submitter:   "eric",
			SubmittedAt: time.Now(),
			UpdatedAt:   time.Now(),
		},
		Revisions: []review.Revision{
			{
				Files: files,
				Annotations: []review.Annotation{
					{File: 0, Hunk: 0, Line: 0, Comment: "foo"},
				},
			},
		},
	}
}

func TestMissingDB(t *testing.T) {
	// Deliberately create a DB that is missing the necessary info.
	tmpdir, err := ioutil.TempDir("/tmp", "erickson")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)
	var boltDB BoltDB
	db, err := bolt.Open(tmpdir+"/erickson.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	boltDB.DB = db
	_, err = boltDB.CreateReview(mockReview)
	if err != ErrNoDB {
		t.Errorf("expected ErrNoDB")
	}
}

// Test an entire review lifecycle.
func TestCRUD(t *testing.T) {
	tmpdir, err := ioutil.TempDir("/tmp", "erickson")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)
	db, err := NewBoltDB(tmpdir + "/erickson.db")
	if err != nil {
		t.Fatal(err)
	}
	id, err := db.CreateReview(mockReview)
	if err != nil {
		t.Error(db.String())
		t.Error(db.GoString())
		t.Fatal(err)
	}
	mockReview.Summary.ID = id
	if exp := 1; id != exp {
		t.Errorf("wrong review id. got %d, want %d", id, exp)
	}

	gotReview, err := db.GetReview(1)
	if err != nil {
		t.Fatal(err)
	}

	if s1, s2 := gotReview.Summary, mockReview.Summary; !reflect.DeepEqual(s1, s1) {
		t.Fatalf("bad summary data. got %+v, want %+v", s1, s2)
	}

	if r1, r2 := gotReview.Revisions, mockReview.Revisions; !reflect.DeepEqual(r1, r2) {
		t.Fatalf("bad revision data. got %+v, want %+v", r1, r2)
	}

	mockReview2 := mockReview

	id, err = db.CreateReview(mockReview2)
	if err != nil {
		t.Error(db.String())
		t.Error(db.GoString())
		t.Fatal(err)
	}
	mockReview2.Summary.ID = id
	if exp := 2; id != exp {
		t.Errorf("wrong review id. got %d, want %d", id, exp)
	}
	files, err := diff.ParseFiles(mockDiff2)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AddRevision(2, review.Revision{Files: files}); err != nil {
		t.Fatal(err)
	}
	gotReview, err = db.GetReview(2)
	if err != nil {
		t.Error(db.String())
		t.Error(db.GoString())
		t.Fatal(err)
	}

	mockReview2.Revisions = append(mockReview2.Revisions, review.Revision{Files: files})

	if s1, s2 := gotReview.Summary, mockReview2.Summary; !reflect.DeepEqual(s1, s1) {
		t.Fatalf("bad summary data. got %+v, want %+v", s1, s2)
	}

	if r1, r2 := gotReview.Revisions, mockReview2.Revisions; !reflect.DeepEqual(r1, r2) {
		t.Fatalf("bad revision data. got %+v, want %+v", r1, r2)
	}

	summaries, err := db.GetSummaries()
	if err != nil {
		t.Fatal(err)
	}
	if sumLen := len(summaries); sumLen != 2 {
		t.Fatalf("wrong number of review summaries: got %d, want 2", sumLen)
	}
	if summaries[0] != mockReview.Summary {
		t.Errorf("wrong review summary: got %+v, want %+v", summaries[0], mockReview.Summary)
	}
	if summaries[1] != mockReview2.Summary {
		t.Errorf("wrong review summary: got %+v, want %+v", summaries[1], mockReview.Summary)
	}

	newSum := mockReview2.Summary

	newSum.Submitter = "boris"

	db.SetSummary(2, newSum)

	summaries, err = db.GetSummaries()
	if err != nil {
		t.Fatal(err)
	}

	if summaries[1] != newSum {
		t.Errorf("wrong review summary: got %+v, want %+v", summaries[1], newSum)
	}

	anno := review.Annotation{
		File:    0,
		Hunk:    0,
		Line:    0,
		Comment: "Fitter, happier, more productive",
		User:    "Fred",
	}
	gotReview.Revisions[0].Annotate(anno)

	if err := db.UpdateRevision(2, 0, gotReview.Revisions[0]); err != nil {
		t.Fatal(err)
	}

	mr2, err := db.GetReview(2)
	if err != nil {
		t.Fatal(err)
	}
	if annotations := mr2.Revisions[0].Annotations; len(annotations) != 2 {
		t.Errorf("wrong number of annotations: got %d, want %d", len(annotations), 2)
	} else if annotations[1] != anno {
		t.Errorf("wrong annotation: got %+v, want %+v", annotations[0], anno)
	}

	if err := db.DeleteReview(3); err == nil {
		t.Error("expected error")
	} else if _, ok := err.(ErrNoReview); !ok {
		t.Error("expected ErrNoReview")
	}

	if err := db.DeleteReview(1); err != nil {
		t.Fatal(err)
	}

	if summaries, err := db.GetSummaries(); err != nil {
		t.Fatal(err)
	} else if len(summaries) != 1 {
		t.Error("delete failed")
	}
}
