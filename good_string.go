// Code generated by "stringer -type=Good"; DO NOT EDIT.

package econerra

import "strconv"

const _Good_name = "GrainVegetablesCottonMeatBeerClothingLabour"

var _Good_index = [...]uint8{0, 5, 15, 21, 25, 29, 37, 43}

func (i Good) String() string {
	if i >= Good(len(_Good_index)-1) {
		return "Good(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Good_name[_Good_index[i]:_Good_index[i+1]]
}
