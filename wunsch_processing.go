// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-29-08 19:28<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package uhebom

import (
	"bytes"
	"strconv"

	wunsch "github.com/MichaelLeachim/wunsch"
	"strings"
)

type TabularForm struct {
	Tag     string
	Path    string
	Content string
	IsGap   bool
}

func newGapTabularForm() *TabularForm {
	tb := newTabularForm("", "", "")
	tb.IsGap = true
	return tb

}
func newTabularForm(tag, path string, content string) *TabularForm {
	return &TabularForm{Path: path, Content: content, Tag: tag}
}

type wunsch_processing_ struct{}

var wunsch_processing = wunsch_processing_{}

func (w *wunsch_processing_) convertToTabularFormRecursive(path string, result *[]*TabularForm, data *DataTree) {
	// we take in only leaf items
	if len(data.Children) == 0 {
		*result = append(*result, newTabularForm(data.Tag, strings.Join([]string{path, "/", data.Tag}, ""), strings.TrimSpace(data.Data)))
		return
	}
	// otherwise, we continue drilling down
	new_path := strings.Join([]string{path, "/", data.Tag}, "")
	for _, child := range data.Children {
		w.convertToTabularFormRecursive(new_path, result, child)
	}
}

func (w *wunsch_processing_) convertToTabularForm(data []*DataTree) []*TabularForm {
	return_data := []*TabularForm{}
	for _, tree := range data {
		result := []*TabularForm{}
		w.convertToTabularFormRecursive(tree.Tag, &result, tree)
		return_data = append(return_data, result...)
	}
	return return_data
}
func (w *wunsch_processing_) dereferenceListOfForms(data []*TabularForm) []TabularForm {
	est := []TabularForm{}
	for _, v := range data {
		est = append(est, *v)
	}
	return est
}

func (w *wunsch_processing_) tabularAlignmentBasedOnWunschAlgorithm(data [][]*TabularForm) [][]*TabularForm {
	converted := [][]wunsch.Item{}
	cached := map[int]*TabularForm{}
	for tab_index, tabular_form := range data {
		converted_item := []wunsch.Item{}
		for index, tabular_item := range tabular_form {
			item_pos := int(utils.hash(strconv.Itoa(tab_index) + "." + strconv.Itoa(index)))
			converted_item = append(converted_item, wunsch.NewItem(item_pos, int(utils.hash(tabular_item.Path))))
			cached[item_pos] = tabular_item
		}
		converted = append(converted, converted_item)
	}
	result := [][]*TabularForm{}
	aligned_item_lists, _ := wunsch.AlignMany(converted)
	for _, aligned_items := range aligned_item_lists {
		subresult := []*TabularForm{}
		for _, item := range aligned_items {
			if item.IsGap() {
				subresult = append(subresult, newGapTabularForm())
			} else {
				item, ok := cached[item.Index]
				if !ok {
					panic("This cannot happen. Check fundamentals!")
				}
				subresult = append(subresult, item)
			}
		}
		result = append(result, subresult)
	}
	return result
}

func (w *wunsch_processing_) asHTMLTables(data [][][]*TabularForm) []byte {
	var template bytes.Buffer
	template.WriteString("<style>table {border-collapse: collapse;}table, th, td {border: 1px solid black;}</style>")

	for i, item := range data {
		template.WriteString("<h1>Table number: ")
		template.WriteString(strconv.Itoa(i))
		template.WriteString(" </h1>")
		template.WriteString(w.asHTMLTable(item))
	}
	return template.Bytes()
}

func (w *wunsch_processing_) extractionWork(root *DataTree, k int, threshold float64) [][][]*TabularForm {
	regions := newMiningDataRegion(root, k, threshold).findRegions(root)
	record_finder := newMiningDataRecord(threshold)

	result := [][][]*TabularForm{}
	for _, region := range regions {
		records := record_finder.findRecords(region)
		tabular_records := [][]*TabularForm{}
		for _, record := range records {
			tabular_records = append(tabular_records, wunsch_processing.convertToTabularForm(record.convertToBase()))
		}
		result = append(result, wunsch_processing.tabularAlignmentBasedOnWunschAlgorithm(tabular_records))
	}
	return result
}

func (w *wunsch_processing_) extract(data []byte) [][][]*TabularForm {
	return w.extractionWork(html_tools.readHTML(html_tools.parseHTML(bytes.NewReader(data))), 5, 0.75)
}

func (w *wunsch_processing_) asHTMLTable(data [][]*TabularForm) string {
	// convert region to a HTML table
	var buffer bytes.Buffer
	buffer.WriteString("<table>")
	for index, item := range data {
		buffer.WriteString("<tr><td>")
		buffer.WriteString(strconv.Itoa(index))
		buffer.WriteString("</td>")
		for _, field := range item {
			buffer.WriteString(`<td>`)
			buffer.WriteString(field.Content)
			buffer.WriteString("</td>")
		}
		buffer.WriteString("</tr>")
	}
	buffer.WriteString("</table>")
	return buffer.String()
}
