package utils

import (
	"sort"
	"strconv"
	"time"
)

var subList = map[string][]string{
	"SHFE":  {"ni", "sn", "ru"},                   //1,5,9
	"SHFE1": {"cu", "al", "zn", "pb", "fu", "wr"}, //每月
	"SHFE2": {"au", "ag"},                         //6,12
	"SHFE3": {"bu"},                               //6,9,12
	"SHFE4": {"rb", "hc"},                         //1,5,10
	"SHFE5": {"ni"},                         //1,5,10
	//"SHFE4": {"cu", "al", "zn", "pb", "ni", "sn", "au", "ag", "rb", "wr", "hc", "fu", "bu", "ru"},
	"CZCE":  {"FG", "MA", "SR", "TA", "RM", "OI", "CF", "CY", "ZC", "SM", "SF", "WH", "JR", "PM", "RI", "LR", "RS"},
	"CZCE1": {"AP"}, //5,7,10,11,12
	"DCE":   {"a", "b", "bb", "c", "cs", "fb", "i", "j", "jd", "jm", "l", "m", "p", "pp", "v", "y"},
	"INE":   {"SC"}, //每月
}

var timeList = map[string][]string{
	"SHFE":  {"01", "05", "09"},
	"SHFE1": {"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"},
	"SHFE2": {"06", "12"},
	"SHFE3": {"06", "09", "12"},
	"SHFE4": {"01", "05", "10"},
	"SHFE5": {"01", "03", "05", "07", "09"},
	"CZCE":  {"01", "05", "09"},
	"CZCE1": {"01", "03", "05", "07", "10", "11", "12"},
	"DCE":   {"01", "05", "09"},
	"INE":   {"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"},
}

var closeList = map[int][]string{
	0: {"fu", "wr", "AP", "pp", "v", "l", "bb", "fb", "c", "cs", "jd", "SM", "SF",
		"WH", "JR", "LR", "PM", "RI", "RS"},
	1: {"ru", "bu", "rb", "hc"},                                                                        //23:00
	2: {"FG", "MA", "SR", "TA", "RM", "OI", "CF", "CY", "ZC", "i", "j", "jm", "a", "b", "m", "p", "y"}, //23:30
	3: {"cu", "pb", "al", "zn", "sn", "ni"},                                                            //01:00
	4: {"au", "ag", "sc"},                                                                              //02:30
}

func dateArr(array []string, index int) []string {
	t := time.Now()
	var arr []string
	t1 := t.AddDate(1, 0, 0)
	for _, v := range array {
		tt, _ := time.Parse("200601", strconv.Itoa(t.Year())+v)
		tt1, _ := time.Parse("200601", strconv.Itoa(t1.Year())+v)
		if tt.After(t) {
			arr = append(arr, strconv.Itoa(tt.Year())[index:4]+v)
		}
		if tt1.After(t) {
			arr = append(arr, strconv.Itoa(tt1.Year())[index:4]+v)
		}
		sort.Sort(sort.StringSlice(arr))
	}
	return arr
}

func SubList() []string {
	var arr []string
	for k, v := range subList {
		for _, vv := range v {
			index := On(!InArray(k, []string{"CZCE", "CZCE1"}), 2, 3).(int)
			tt := dateArr(timeList[k], index)
			for i := 0; i < len(tt); i++ {
				if i > 3 {
					continue
				}
				arr = append(arr, vv+tt[i])
			}
		}
	}
	return arr
}

type CloseType struct {
	Symbol    string
	CloseType int
}

func CloseList() []CloseType {
	var arr []CloseType
	for k, v := range subList {
		closeTime := 0

		for _, vv := range v {
			index := On(!InArray(k, []string{"CZCE", "CZCE1"}), 2, 3).(int)
			tt := dateArr(timeList[k], index)
			for key, symbol := range closeList {
				if InArray(vv, symbol) {
					closeTime = key
				}
			}
			for i := 0; i < len(tt); i++ {
				if i > 3 {
					continue
				}
				arr = append(arr, CloseType{vv + tt[i], closeTime})
			}
		}
	}
	return arr
}

func CloseListArray() map[int][]string {
	arr := make(map[int][]string)
	for k, v := range subList {
		closeTime := 0

		for _, vv := range v {
			index := On(!InArray(k, []string{"CZCE", "CZCE1"}), 2, 3).(int)
			tt := dateArr(timeList[k], index)
			for key, symbol := range closeList {
				if InArray(vv, symbol) {
					closeTime = key
				}
			}
			for i := 0; i < len(tt); i++ {
				if i > 3 {
					continue
				}
				arr[closeTime] = append(arr[closeTime], vv+tt[i])
			}
		}
	}
	return arr
}
