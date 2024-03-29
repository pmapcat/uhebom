// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:35<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package uhebom

import ()

type trees_utils_ struct{}

var trees_utils = trees_utils_{}

func (h *trees_utils_) pairwise(data []*DataTree, K, start int) [][][]*DataTree {
	// TODO: check index sizes
	result := [][][]*DataTree{}
	// _ = "breakpoint"
	for k := 1; k < K+1; k++ {
		for i := 0; i < K; i++ {
			for j := start + i; j < len(data); j += k {
				slice_ax, slice_ay := utils.bindSlice(j, j+k, len(data))
				slice_bx, slice_by := utils.bindSlice(j+k, j+2*k, len(data))

				slice_a := data[slice_ax:slice_ay]
				slice_b := data[slice_bx:slice_by]
				if len(slice_a) >= k && len(slice_b) >= k {
					result = append(result, [][]*DataTree{slice_a, slice_b})
				}
			}
		}
	}
	return result
}

func (t *trees_utils_) create2dMatrix(x, y int) [][]float64 {
	result := make([][]float64, x)
	for i := 0; i < x; i++ {
		result[i] = make([]float64, y)
	}
	return result
}

func (u *trees_utils_) treeMatch(t1, t2 *DataTree) float64 {
	t1_root, t1_exist := t1.getRoot()
	t2_root, t2_exist := t2.getRoot()
	if !t1_exist || !t2_exist {
		return 0
	}
	if t1_root != t2_root {
		return 0
	}
	rows := t1.getChildrenCount() + 1
	cols := t2.getChildrenCount() + 1
	m := u.create2dMatrix(rows, cols)
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			child1, _ := t1.getChild(i - 1)
			child2, _ := t2.getChild(j - 1)
			m[i][j] = utils.maxf([]float64{m[i][j], m[i-1][j-1] + u.treeMatch(child1, child2)})
		}
	}
	return 1 + m[rows-1][cols-1]
}
