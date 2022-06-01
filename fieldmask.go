package fieldmask

import (
	"strings"
)

// FieldMask represents a subset of fields that should be returned by
// a read operation or modified by an update operation.
type FieldMask map[string]interface{}

// Get looks up the given dot-separated path from fm. It will return the value
// for the specified path, or ok=false if the path does not exist.
func (fm FieldMask) Get(path string) (value interface{}, ok bool) {
	var obj interface{} = map[string]interface{}(fm)

	parts := strings.Split(path, ".")
	for _, p := range parts {
		m, ok := obj.(map[string]interface{})
		if !ok {
			return nil, false
		}

		obj, ok = m[p]
		if !ok {
			return nil, false
		}

		value = obj
	}

	return value, true
}

// Has reports whether the given dot-separated path exists or not.
func (fm FieldMask) Has(path string) bool {
	_, ok := fm.Get(path)
	return ok
}

// FieldMask returns the nested field mask at the given path. It will return
// ok=false if the path does not exist, or the corresponding value is not a map.
func (fm FieldMask) FieldMask(path string) (mask FieldMask, ok bool) {
	value, _ := fm.Get(path)

	m, ok := value.(map[string]interface{})
	return m, ok
}

// Copy copies the values for the given paths into fm from a source map src.
// For each path that does not exist in src, no corresponding value will be
// copied into fm.
func (fm FieldMask) Copy(src map[string]interface{}, paths ...string) {
	for _, path := range paths {
		fm.copy(src, path)
	}
}

func (fm FieldMask) copy(src map[string]interface{}, path string) {
	var srcObj interface{} = src
	var dstM = map[string]interface{}(fm)

	parts := strings.Split(path, ".")
	for i, p := range parts {
		srcM, ok := srcObj.(map[string]interface{})
		if !ok {
			// Ignore non-existent path.
			return
		}

		srcV, ok := srcM[p]
		if !ok {
			// Ignore non-existent path.
			return
		}

		if i < len(parts)-1 {
			// Now p points to a non-leaf value in src.

			// Ensure that the corresponding value in dstM is always a non-nil map.
			dstV, ok := dstM[p].(map[string]interface{})
			if !ok {
				dstV = make(map[string]interface{})
				dstM[p] = dstV
			}

			// Enter into the next deeper level of dstM.
			dstM = dstV
		} else {
			// Now p points to a leaf value in src, copy this value into dstM.
			dstM[p] = srcV
		}

		// Enter into the next deeper level of src.
		srcObj = srcV
	}
}
