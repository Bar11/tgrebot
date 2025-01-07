// Package compare
package compare

import (
	"fmt"
	"strings"

	"github.com/pierrre/compare"
)

func CompareStruct[T any](s1, s2 T) (compare.Result, map[string]interface{}) {
	diff := compare.Compare(s1, s2)
	diffCount := make(map[string]interface{})
	for _, difference := range diff {
		path := fmt.Sprintf("%v", difference.Path) // 新内容
		path = strings.Replace(path, ".", "", 1)
		if _, ok := diffCount[path]; ok {
			continue
		}
		diffCount[path] = path
		pathFull := difference.V2 // 原内容
		pathFull = strings.Replace(pathFull, ".", "", 1)
		splits := strings.Split(pathFull, ".")
		var pathFull1 = ""
		for i, split := range splits {
			if i == 0 {
				pathFull1 = split
			} else {
				pathFull1 = pathFull1 + "." + split
			}
			if _, ok := diffCount[pathFull1]; ok {
				continue
			}
			diffCount[pathFull1] = pathFull1
		}
	}
	return diff, diffCount
}
