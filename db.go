package main

import (
	"bytes"
	"strings"

	"github.com/corona10/goimagehash"
	"github.com/rs/xid"

	"github.com/syndtr/goleveldb/leveldb"
)

func storeAppend(id string, hashes ...*goimagehash.ImageHash) error {
	batch := new(leveldb.Batch)
	buf := &bytes.Buffer{}
	for _, h := range hashes {
		if err := h.Dump(buf); err != nil {
			continue
		}

		batch.Put([]byte(id+"/"+xid.New().String()), buf.Bytes())
	}
	return db.Write(batch, nil)
}

func storeFind(maxDistance int, hashes ...*goimagehash.ImageHash) []string {
	if len(hashes) < 1 {
		return []string{}
	}

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	found := map[string]bool{}
	ret := []string{}

	for iter.Next() {
		id := strings.SplitN(string(iter.Key()), "/", 2)[0]
		buf := bytes.NewReader(iter.Value())
		hash1, err := goimagehash.LoadImageHash(buf)
		if err != nil {
			continue
		}

		for _, hash2 := range hashes {
			d, err := hash2.Distance(hash1)
			if err != nil {
				continue
			}

			// fmt.Println(hash1.GetHash(), hash2.GetHash(), d)

			if d <= maxDistance {
				found[id] = true
			}
		}
	}

	for id := range found {
		ret = append(ret, id)
	}

	return ret
}
