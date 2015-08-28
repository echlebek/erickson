package db

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/echlebek/erickson/review"
)

type BoltDB struct {
	*bolt.DB
}

type metaData struct {
	Summaries map[string]review.Summary `json:"summaries"`
}

var (
	_            Database = &BoltDB{}
	rootKey               = []byte("erickson")
	metaKey               = []byte("metadata")
	reviewsKey            = []byte("reviews")
	revisionsKey          = []byte("revisions")
	summaryKey            = []byte("summary")
)

func getReviewsBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	root := tx.Bucket(rootKey)
	if root == nil {
		return nil, ErrNoDB
	}
	reviews := root.Bucket(reviewsKey)
	if reviews == nil {
		return nil, ErrNoDB
	}
	return reviews, nil
}

func getReviewBucket(tx *bolt.Tx, id int) (*bolt.Bucket, error) {
	reviewsBkt, err := getReviewsBucket(tx)
	if err != nil {
		return nil, err
	}
	bucket := reviewsBkt.Bucket(int2key(id))
	if bucket == nil {
		return nil, ErrNoReview(id)
	}
	return bucket, nil
}

func getMetaData(tx *bolt.Tx) (metaData, error) {
	var meta metaData
	root := tx.Bucket(rootKey)
	if root == nil {
		return meta, ErrNoDB
	}
	metaValue := root.Get(metaKey)
	err := json.Unmarshal(metaValue, &meta)
	return meta, err
}

func setMetaData(tx *bolt.Tx, meta metaData) error {
	root := tx.Bucket(rootKey)
	if root == nil {
		return ErrNoDB
	}
	metaValue, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return root.Put(metaKey, metaValue)
}

func NewBoltDB(path string) (*BoltDB, error) {
	var boltDB BoltDB
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	boltDB.DB = db

	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists(rootKey)
		if err != nil {
			return err
		}
		if v := root.Get(metaKey); v == nil {
			meta := metaData{Summaries: make(map[string]review.Summary)}
			if err := setMetaData(tx, meta); err != nil {
				return err
			}
		}
		_, err = root.CreateBucketIfNotExists(reviewsKey)
		return err
	})

	return &boltDB, err
}

func (db *BoltDB) CreateReview(r review.R) (int, error) {
	var (
		err      error
		reviewID int
	)
	err = db.Update(func(tx *bolt.Tx) error {
		meta, err := getMetaData(tx)
		if err != nil {
			return err
		}
		reviewsBkt, err := getReviewsBucket(tx)
		if err != nil {
			return err
		}
		ns, err := reviewsBkt.NextSequence()
		if err != nil {
			return err
		}
		reviewID = int(ns)
		newReview, err := reviewsBkt.CreateBucket(int2key(reviewID))
		if err != nil {
			return err
		}
		r.Summary.ID = reviewID
		meta.Summaries[strconv.Itoa(reviewID)] = r.Summary
		if err := setMetaData(tx, meta); err != nil {
			return err
		}
		revisionsValue, err := json.Marshal(r.Revisions)
		if err != nil {
			return err
		}
		return newReview.Put(revisionsKey, revisionsValue)
	})
	return int(reviewID), err
}

func (db *BoltDB) GetReview(id int) (review.R, error) {
	var (
		review review.R
		err    error
	)
	err = db.Update(func(tx *bolt.Tx) error {
		review, err = getReview(tx, id)
		return err
	})
	return review, err
}

func getReview(tx *bolt.Tx, id int) (review.R, error) {
	var (
		revisions []review.Revision
		review    review.R
	)

	metaData, err := getMetaData(tx)
	if err != nil {
		return review, err
	}

	root := tx.Bucket(rootKey)
	if root == nil {
		return review, ErrNoDB
	}
	reviews := root.Bucket(reviewsKey)
	if reviews == nil {
		return review, ErrNoDB
	}
	reviewBkt, err := getReviewBucket(tx, id)
	if err != nil {
		return review, err
	}
	review.Summary = metaData.Summaries[strconv.Itoa(id)]
	rev := reviewBkt.Get(revisionsKey)
	if err := json.Unmarshal(rev, &revisions); err != nil {
		return review, err
	}
	review.Revisions = revisions
	return review, nil
}

type reviewSlice []review.Summary

func (r reviewSlice) Len() int {
	return len(r)
}

func (r reviewSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r reviewSlice) Less(i, j int) bool {
	return r[i].ID < r[j].ID
}

func (db *BoltDB) GetSummaries() ([]review.Summary, error) {
	var (
		result   []review.Summary
		metaData metaData
	)
	err := db.Update(func(tx *bolt.Tx) (err error) {
		metaData, err = getMetaData(tx)
		return
	})
	for _, s := range metaData.Summaries {
		result = append(result, s)
	}
	sort.Sort(reviewSlice(result))
	return result, err
}

func (db *BoltDB) SetSummary(id int, summary review.Summary) error {
	return db.Update(func(tx *bolt.Tx) error {
		metaData, err := getMetaData(tx)
		if err != nil {
			return err
		}
		metaData.Summaries[strconv.Itoa(id)] = summary
		return setMetaData(tx, metaData)
	})
}

func int2key(i int) []byte {
	return strconv.AppendInt([]byte{}, int64(i), 10)
}

func (db *BoltDB) DeleteReview(id int) error {
	return db.Update(func(tx *bolt.Tx) error {
		reviewsBkt, err := getReviewsBucket(tx)
		if err != nil {
			return err
		}
		if err := reviewsBkt.DeleteBucket(int2key(id)); err != nil {
			return reviewBucketErr(id, err)
		}
		metaData, err := getMetaData(tx)
		if err != nil {
			return err
		}
		delete(metaData.Summaries, strconv.Itoa(id))
		return setMetaData(tx, metaData)
	})
}

func reviewBucketErr(id int, err error) error {
	switch err {
	case bolt.ErrBucketNotFound:
		return ErrNoReview(id)
	}
	return err
}

func (db *BoltDB) AddRevision(id int, revision review.Revision) error {
	// FIXME: Might be expensive. Give revisions its own bucket if this is too slow.
	return db.Update(func(tx *bolt.Tx) error {
		var revisions []review.Revision
		reviewBkt, err := getReviewBucket(tx, id)
		if err != nil {
			return err
		}
		revisionsValue := reviewBkt.Get(revisionsKey)
		if revisionsValue != nil {
			if err := json.Unmarshal(revisionsValue, &revisions); err != nil {
				return err
			}
		}
		revisions = append(revisions, revision)
		revisionsValue, err = json.Marshal(revisions)
		if err != nil {
			return err
		}
		return reviewBkt.Put(revisionsKey, revisionsValue)
	})
}

func (db *BoltDB) AddAnnotation(id, revId int, an review.Annotation) error {
	return db.Update(func(tx *bolt.Tx) error {
		reviewBkt, err := getReviewBucket(tx, id)
		if err != nil {
			return err
		}
		revisionsValue := reviewBkt.Get(revisionsKey)
		var revisions []review.Revision
		if err := json.Unmarshal(revisionsValue, &revisions); err != nil {
			return err
		}
		if revId >= len(revisions) {
			return ErrNoRevision(revId)
		}
		revisions[revId].Annotations = append(revisions[revId].Annotations, an)
		revisionsValue, err = json.Marshal(revisions)
		if err != nil {
			return err
		}
		return reviewBkt.Put(revisionsKey, revisionsValue)
	})
}
