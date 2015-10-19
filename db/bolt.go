package db

import (
	"bytes"
	"encoding/gob"
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
	usersKey              = []byte("users")
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
	dec := gob.NewDecoder(bytes.NewReader(metaValue))
	err := dec.Decode(&meta)
	return meta, err
}

func setMetaData(tx *bolt.Tx, meta metaData) error {
	root := tx.Bucket(rootKey)
	if root == nil {
		return ErrNoDB
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(meta); err != nil {
		return err
	}
	return root.Put(metaKey, buf.Bytes())
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
		if _, err := root.CreateBucketIfNotExists(reviewsKey); err != nil {
			return err
		}
		_, err = root.CreateBucketIfNotExists(usersKey)
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
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(r.Revisions); err != nil {
			return err
		}
		return newReview.Put(revisionsKey, buf.Bytes())
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
	dec := gob.NewDecoder(bytes.NewReader(rev))
	if err := dec.Decode(&revisions); err != nil {
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
	return db.Update(func(tx *bolt.Tx) error {
		var revisions []review.Revision
		reviewBkt, err := getReviewBucket(tx, id)
		if err != nil {
			return err
		}
		revisionsValue := reviewBkt.Get(revisionsKey)
		if revisionsValue != nil {
			dec := gob.NewDecoder(bytes.NewReader(revisionsValue))
			if err := dec.Decode(&revisions); err != nil {
				return err
			}
		}
		revisions = append(revisions, revision)
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(revisions); err != nil {
			return err
		}
		return reviewBkt.Put(revisionsKey, buf.Bytes())
	})
}

// UpdateRevision overwrites the specified revision with a new one.
func (db *BoltDB) UpdateRevision(id, revId int, revision review.Revision) error {
	return db.Update(func(tx *bolt.Tx) error {
		reviewBkt, err := getReviewBucket(tx, id)
		if err != nil {
			return err
		}
		var revisions []review.Revision
		revisionsValue := reviewBkt.Get(revisionsKey)
		dec := gob.NewDecoder(bytes.NewReader(revisionsValue))
		if err := dec.Decode(&revisions); err != nil {
			return err
		}
		if revId >= len(revisions) {
			return ErrNoRevision(revId)
		}
		revisions[revId] = revision
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(revisions); err != nil {
			return err
		}
		return reviewBkt.Put(revisionsKey, buf.Bytes())
	})
}

func getUsersBucket(tx *bolt.Tx) (b *bolt.Bucket, err error) {
	root := tx.Bucket(rootKey)
	if root == nil {
		err = ErrNoDB
		return
	}
	b = root.Bucket(usersKey)
	if b == nil {
		err = ErrNoDB
	}
	return
}

func (db *BoltDB) CreateUser(u review.User) error {
	return db.Update(func(tx *bolt.Tx) error {
		userBkt, err := getUsersBucket(tx)
		if err != nil {
			return err
		}
		if userValue := userBkt.Get([]byte(u.Name)); userValue != nil {
			return ErrUserExists
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(u); err != nil {
			return err
		}
		return userBkt.Put([]byte(u.Name), buf.Bytes())
	})
}

func (db *BoltDB) UpdateUser(u review.User) error {
	return db.Update(func(tx *bolt.Tx) error {
		userBkt, err := getUsersBucket(tx)
		if err != nil {
			return err
		}
		if userValue := userBkt.Get([]byte(u.Name)); userValue == nil {
			return ErrNoUser
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(u); err != nil {
			return err
		}
		return userBkt.Put([]byte(u.Name), buf.Bytes())
	})
}

func (db *BoltDB) GetUser(username string) (u review.User, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		userBkt, err := getUsersBucket(tx)
		if err != nil {
			return err
		}
		userValue := userBkt.Get([]byte(username))
		if userValue == nil {
			return ErrNoUser
		}
		dec := gob.NewDecoder(bytes.NewReader(userValue))
		return dec.Decode(&u)
	})
	return
}

func (db *BoltDB) DeleteUser(username string) error {
	return db.Update(func(tx *bolt.Tx) error {
		userBkt, err := getUsersBucket(tx)
		if err != nil {
			return err
		}
		return userBkt.Delete([]byte(username))
	})
}
