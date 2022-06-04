package fieldmask

import (
	"strings"
)

// FieldMask represents a subset of fields that should be returned by
// a read operation or modified by an update operation.
type FieldMask map[string]interface{}

// From creates a new field mask, which only contains the values for the given
// dot-separated paths from a source map src. For each path that does not exist
// in src, a nil value will be inserted into the final field mask.
//
// Specially, empty paths indicates all values from src should be included, which
// is now implemented by just returning src internally. Therefore, in this case,
// note that any further change of the created field mask will affect src as well.
func From(src map[string]interface{}, paths ...string) FieldMask {
	if len(paths) == 0 {
		return src
	}

	fm := map[string]interface{}{}
	for _, path := range paths {
		copy(fm, src, path)
	}
	return fm
}

// Get looks up the given dot-separated path from fm. It will return the value
// for the specified path, or ok=false if the path does not exist.
func (fm FieldMask) Get(path string) (value interface{}, ok bool) {
	var obj interface{} = map[string]interface{}(fm)
	var m map[string]interface{}

	parts := strings.Split(path, ".")
	for _, p := range parts {
		m, ok = obj.(map[string]interface{})
		if !ok {
			return nil, false
		}

		obj, ok = m[p]
		if !ok {
			return nil, false
		}

		value = obj
	}

	return value, ok
}

// Has reports whether the given dot-separated path exists or not.
func (fm FieldMask) Has(path string) bool {
	_, ok := fm.Get(path)
	return ok
}

// FieldMask returns the nested field mask at the given dot-separated path.
// It will return ok=false if the path does not exist, or the corresponding
// non-nil value is not a map.
func (fm FieldMask) FieldMask(path string) (mask FieldMask, ok bool) {
	value, ok := fm.Get(path)
	if ok && value == nil {
		// Nil value is treated as an empty field mask.
		return nil, true
	}

	// Non-nil value must be a map.
	m, ok := value.(map[string]interface{})
	return m, ok
}

func copy(dst, src map[string]interface{}, path string) {
	var srcObj interface{} = src
	var dstM = dst

	parts := strings.Split(path, ".")
	for i, p := range parts {
		srcM, _ := srcObj.(map[string]interface{})
		srcV := srcM[p]

		if i < len(parts)-1 {
			// Since p points to a non-tail part of the whole path, we must
			// ensure that the corresponding value in dstM is always a non-nil map.
			dstV, ok := dstM[p].(map[string]interface{})
			if !ok {
				dstV = make(map[string]interface{})
				dstM[p] = dstV
			}

			dstM = dstV   // Enter into the next deeper level of dst.
			srcObj = srcV // Enter into the next deeper level of src.

			continue
		}

		// Now p points to the last part, copy the final value (might be nil) into dstM.
		dstM[p] = srcV
	}
}
