package blob

import (
	"sort"
	"strings"
)

var (
	sTrim      = strings.Trim
	sTrimRight = strings.TrimRight
	sSplit     = strings.Split
)

func Map2KVs4BL(m map[int][]*blobline, less4key func(i int, j int) bool) (keys []int, values [][]*blobline) {

	if m == nil {
		return nil, nil
	}
	if len(m) == 0 {
		return []int{}, [][]*blobline{}
	}

	type kv struct {
		key   int
		value []*blobline
	}

	kvSlc := []kv{}
	for k, v := range m {
		kvSlc = append(kvSlc, kv{key: k, value: v})
	}

	switch {
	case less4key != nil:
		sort.SliceStable(kvSlc, func(i, j int) bool { return less4key(kvSlc[i].key, kvSlc[j].key) })
	default:
		// do not sort
	}

	for _, kvEle := range kvSlc {
		keys = append(keys, kvEle.key)
		values = append(values, kvEle.value)
	}
	return
}

func Map2KVs4Blob(m map[int][]*blobline, less4key func(i int, j int) bool) (keys []int, values [][]*blobline) {

	if m == nil {
		return nil, nil
	}
	if len(m) == 0 {
		return []int{}, [][]*blobline{}
	}

	type kv struct {
		key   int
		value []*blobline
	}

	kvSlc := []kv{}
	for k, v := range m {
		kvSlc = append(kvSlc, kv{key: k, value: v})
	}

	switch {
	case less4key != nil:
		sort.SliceStable(kvSlc, func(i, j int) bool { return less4key(kvSlc[i].key, kvSlc[j].key) })
	default:
		// do not sort
	}

	for _, kvEle := range kvSlc {
		keys = append(keys, kvEle.key)
		values = append(values, kvEle.value)
	}
	return
}
