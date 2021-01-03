package myModule

import "time"

func GetJapaneseHolidays(t time.Time) []int {
	//	todo: 祝日が担保できないから表示年月絞る?
	targetMonth := int(t.Month())

	// todo: 設定ファイルからにしたい(春分とかの式はどうする？)
	holidays := make([][]interface{}, 12)
	holidays[0] = []interface{}{1, "second-friday"} //	第2金曜
	//holidays[0] = []interface{}{1, 21, "second-friday"} //	第2金曜 fixme: これはテスト
	holidays[1] = []interface{}{11, emperorBirth(t)}
	holidays[2] = []interface{}{getSpringDay(t)} //	春分日
	holidays[3] = []interface{}{29}
	holidays[4] = []interface{}{3, 4, 5}
	holidays[5] = []interface{}{}
	holidays[6] = []interface{}{"third-friday"} //	第3金曜
	holidays[7] = []interface{}{11}
	holidays[8] = []interface{}{"third-monday", getAutumnDay(t)} //	第3金曜, 秋分日
	holidays[9] = []interface{}{"second-monday"}                 //	第2月曜
	holidays[10] = []interface{}{3, 23}
	holidays[11] = []interface{}{emperorBirth(t)}

	var result []int
	for _, item := range holidays[targetMonth-1] {
		switch item.(type) {
		case string:
			result = append(result, ParseDay(t, item.(string)))
		case int:
			result = append(result, item.(int))
		default:
			return []int{}
		}
	}
	return result
}

func getSpringDay(t time.Time) int {
	return int(20.8431+0.242194*float32(t.Year()-1980)) - (t.Year()-1980)/4
}

func getAutumnDay(t time.Time) int {
	return int(23.2488+0.242194*float32(t.Year()-1980)) - (t.Year()-1980)/4
}

func emperorBirth(t time.Time) int {
	if t.Year() <= 2018 && int(t.Month()) == 12 || t.Year() >= 2020 && int(t.Month()) == 2 {
		return 23
	}
	return -1
}
