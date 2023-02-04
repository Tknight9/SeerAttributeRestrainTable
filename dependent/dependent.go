package dependent

import (
	"sort"
	"sync"
)

const (
	MinAtbID          = 11
	MaxAtbID          = 36
	AttributeSimpleID = 20
)

type AtbNode struct {
	id     uint16
	isBase bool
	first  uint16
	second uint16

	name      string
	imgPath   string
	sibling   []uint16
	restraint []uint16
}

type RelationTable struct {
	PositiveAttack []uint16
	NegativeAttack []uint16
	InvalidAttack  uint16

	PositiveRecipient []uint16
	NegativeRecipient []uint16
	InvalidRecipient  uint16
}

type TimesPair struct {
	Attacker  uint16
	Recipient uint16
	Times     float64
}

type AtbID struct {
	Fir uint16
	Sec uint16
}

var (
	lock      sync.RWMutex
	AtbIDList = []uint16{
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36,

		1114, 1115, 1118, 1119, 1121, 1122, 1123, 1124, 1125, 1126, 1128, 1132,
		1219, 1221, 1222, 1224, 1227,
		1319, 1321, 1322, 1328,
		1415, 1418, 1419, 1422, 1423, 1424, 1425, 1432,
		1518, 1521, 1522, 1524, 1526, 1527,
		1617, 1619, 1621, 1622, 1624, 1625, 1626, 1628, 1632,
		1719, 1721, 1724, 1725, 1727, 1729,
		1819, 1822, 1823, 1824, 1825, 1830, 1832,
		1923, 1925, 1926, 1932,
		2122, 2123, 2124, 2125, 2126, 2127, 2128, 2130, 2132,
		2223, 2224, 2225, 2226, 2229, 2232,
		2325, 2326, 2327,
		2428, 2429, 2430, 2432,
		2526, 2528, 2529, 2530,
		2627, 2628, 2630, 2632, 2634,
		2732, 2829, 2832, 2932, 2936, 3236,
	}
	AtbNameMap = map[uint16]string{
		11: "火", 12: "水", 13: "草", 14: "飞行", 15: "电",
		16: "地面", 17: "机械", 18: "冰", 19: "超能", 20: "普通",
		21: "战斗", 22: "暗影", 23: "光", 24: "龙", 25: "神秘",
		26: "圣灵", 27: "次元", 28: "远古", 29: "邪灵", 30: "自然",
		31: "王", 32: "混沌", 33: "神灵", 34: "轮回", 35: "虫",
		36: "虚空",

		1114: "火飞行", 1115: "电火", 1118: "冰火", 1119: "火超能", 1121: "战斗火", 1122: "暗影火", 1123: "光火", 1124: "火龙", 1125: "火神秘", 1126: "圣灵火", 1128: "远古火", 1132: "混沌火",
		1219: "水超能", 1221: "水战斗", 1222: "水暗影", 1224: "水龙", 1227: "水次元",
		1319: "草超能", 1321: "草战斗", 1322: "草暗影", 1328: "远古草",
		1415: "电飞行", 1418: "冰飞行", 1419: "飞行超能", 1422: "飞行暗影", 1423: "光飞行", 1424: "飞龙", 1425: "飞行神秘", 1432: "飞行混沌",
		1518: "电冰", 1521: "电战斗", 1522: "电暗影", 1524: "电龙", 1526: "圣灵电", 1527: "电次元",
		1617: "机械地面", 1619: "地面超能", 1621: "地面战斗", 1622: "地面暗影", 1624: "地面龙", 1625: "地面神秘", 1626: "圣灵地面", 1628: "远古地面", 1632: "混沌地面",
		1719: "机械超能", 1721: "机械战斗", 1724: "机械龙", 1725: "机械神秘", 1727: "机械次元", 1729: "机械邪灵",
		1819: "冰超能", 1822: "冰暗影", 1823: "冰光", 1824: "冰龙", 1825: "冰神秘", 1830: "自然冰", 1832: "混沌冰",
		1923: "光超能", 1925: "神秘超能", 1926: "圣灵超能", 1932: "混沌超能",
		2122: "战斗暗影", 2123: "光战斗", 2124: "战斗龙", 2125: "神秘战斗", 2126: "圣灵战斗", 2127: "次元战斗", 2128: "远古战斗", 2130: "自然战斗", 2132: "混沌战斗",
		2223: "光暗影", 2224: "暗影龙", 2225: "神秘暗影", 2226: "圣灵暗影", 2229: "邪灵暗影", 2232: "混沌暗影",
		2325: "光神秘", 2326: "圣灵光", 2327: "光次元",
		2428: "远古龙", 2429: "邪灵龙", 2430: "自然龙", 2432: "混沌龙",
		2526: "圣灵神秘", 2528: "远古神秘", 2529: "邪灵神秘", 2530: "自然神秘",
		2627: "圣灵次元", 2628: "圣灵远古", 2630: "圣灵自然", 2632: "混沌圣灵", 2634: "圣灵轮回",
		2732: "混沌次元",
		2829: "远古邪灵", 2832: "混沌远古",
		2932: "混沌邪灵", 2936: "虚空邪灵",
		3236: "虚空混沌",
	}
	BaseAtbRestraintTable = map[uint16]RelationTable{
		//20: {PositiveAttack: []uint16{}, []uint16{}, 0},
		11: {
			[]uint16{13, 17, 18},
			[]uint16{11, 12, 26, 30, 32, 33},
			0,

			[]uint16{12, 16, 26, 30, 33},
			[]uint16{11, 13, 17, 18, 24, 35},
			0,
		},
		12: {
			[]uint16{11, 16},
			[]uint16{12, 13, 26, 30, 32, 33},
			0,

			[]uint16{13, 15, 26, 30, 33},
			[]uint16{11, 12, 17, 18, 24, 35},
			0,
		},
		13: {
			[]uint16{12, 16, 23},
			[]uint16{11, 13, 14, 17, 26, 28, 32, 33},
			0,

			[]uint16{11, 14, 18, 26, 28, 30, 33, 35},
			[]uint16{12, 13, 15, 16, 24},
			0,
		},
		14: {[]uint16{13, 21, 35},
			[]uint16{15, 17, 27, 29, 30, 32},
			0,

			[]uint16{15, 18, 27, 28, 30, 32},
			[]uint16{13, 36},
			16,
		},
		15: {[]uint16{12, 14, 22, 27, 32, 36},
			[]uint16{13, 15, 25, 26, 30, 33},
			16,

			[]uint16{16, 25, 26, 30, 33},
			[]uint16{14, 15, 17, 24, 32},
			0,
		},
		16: {[]uint16{11, 15, 17, 31, 34},
			[]uint16{13, 19, 22, 24, 26, 30, 33, 35},
			14,

			[]uint16{12, 13, 18, 30, 35},
			[]uint16{25},
			15,
		},
		17: {[]uint16{18, 21, 28, 29, 33},
			[]uint16{11, 12, 15, 17, 27},
			0,

			[]uint16{11, 16, 21, 27},
			[]uint16{13, 14, 17, 18, 19, 22, 23, 28, 29, 30, 32, 33},
			0,
		},
		18: {[]uint16{13, 14, 16, 27, 28, 34, 35},
			[]uint16{11, 12, 17, 18, 26, 32, 33},
			0,

			[]uint16{11, 17, 21, 24, 26, 32, 33},
			[]uint16{18, 22, 23, 27, 28, 29, 34, 35},
			0,
		},
		19: {[]uint16{21, 25, 30},
			[]uint16{17, 19, 35},
			23,

			[]uint16{22, 23, 27, 36},
			[]uint16{16, 19, 21, 29, 30, 31, 34},
			0,
		},
		21: {[]uint16{17, 18, 24, 26},
			[]uint16{19, 21, 22, 29, 31},
			0,

			[]uint16{14, 17, 19, 31, 35, 36},
			[]uint16{21, 25, 26, 30, 32, 33},
			0,
		},
		22: {[]uint16{19, 22, 27},
			[]uint16{17, 18, 23, 26, 29, 33},
			0,

			[]uint16{15, 22, 23, 29, 31, 34},
			[]uint16{16, 21, 30, 36},
			27,
		},
		23: {[]uint16{19, 22, 35},
			[]uint16{17, 18, 23, 26, 29, 30, 33, 34, 36},
			13,

			[]uint16{13, 29, 30, 34, 36},
			[]uint16{22, 23, 35},
			19,
		},
		24: {[]uint16{18, 24, 26, 29},
			[]uint16{11, 12, 13, 15, 28, 35},
			0,

			[]uint16{21, 24, 28},
			[]uint16{16, 26, 33},
			0,
		},
		25: {[]uint16{15, 25, 26, 30, 31, 33, 34},
			[]uint16{16, 21, 29, 32, 35},
			0,

			[]uint16{19, 25, 28, 29, 32, 36},
			[]uint16{15, 26, 30},
			0,
		},
		26: {[]uint16{11, 12, 13, 15, 18, 28, 36},
			[]uint16{21, 24, 25, 34},
			0,

			[]uint16{21, 24, 25, 34},
			[]uint16{11, 12, 13, 15, 16, 18, 22, 23, 29, 36},
			0,
		},
		27: {[]uint16{14, 17, 19, 29, 30, 35, 36},
			[]uint16{18, 31, 32, 33, 34},
			22,

			[]uint16{15, 18, 22, 29, 31, 32, 34},
			[]uint16{14, 17, 30, 36},
			0,
		},
		28: {[]uint16{13, 14, 24, 25, 36},
			[]uint16{17, 18, 31, 34},
			0,

			[]uint16{17, 18, 26, 33},
			[]uint16{13, 24},
			0,
		},
		29: {[]uint16{22, 23, 25, 27, 30},
			[]uint16{17, 18, 19, 26, 31, 32, 34},
			33,

			[]uint16{17, 24, 27, 31, 32, 33, 34},
			[]uint16{14, 21, 22, 23, 25, 30},
			0,
		},
		30: {[]uint16{11, 12, 13, 14, 15, 16, 23, 31, 34},
			[]uint16{17, 19, 21, 22, 25, 27, 29, 32, 36},
			0,

			[]uint16{19, 25, 27, 29, 32, 36},
			[]uint16{11, 12, 14, 15, 16, 23, 31, 34},
			0,
		},
		31: {[]uint16{21, 22, 27, 29},
			[]uint16{19, 30, 35},
			0,

			[]uint16{16, 25, 30},
			[]uint16{21, 27, 28, 29},
			0,
		},
		32: {[]uint16{14, 18, 25, 27, 29, 30, 33},
			[]uint16{15, 17, 21, 34},
			36,

			[]uint16{15, 34, 35, 36},
			[]uint16{11, 12, 13, 14, 18, 25, 27, 29, 30},
			0,
		},
		33: {[]uint16{11, 12, 13, 15, 18, 28, 29, 32},
			[]uint16{17, 21, 24},
			0,

			[]uint16{17, 25, 32},
			[]uint16{11, 12, 13, 15, 16, 18, 22, 23, 27},
			29,
		},
		34: {[]uint16{22, 23, 26, 27, 29, 32},
			[]uint16{18, 19, 30, 36},
			0,

			[]uint16{16, 18, 25, 30, 36},
			[]uint16{23, 26, 27, 28, 29, 32},
			0,
		},
		35: {[]uint16{13, 16, 21, 32, 35},
			[]uint16{11, 12, 18, 23},
			0,

			[]uint16{14, 18, 23, 27, 35},
			[]uint16{16, 19, 24, 25, 31},
			0,
		},
		36: {[]uint16{19, 21, 23, 25, 30, 34},
			[]uint16{14, 22, 26, 27},
			0,

			[]uint16{15, 26, 27, 28},
			[]uint16{23, 30, 34},
			32,
		},
	}
	BaseAtbRelate = make([][]uint16, 40)
)

func GetAtbIDType(id uint16) int {
	_, isExist := AtbNameMap[id]
	if isExist {
		if id < 1000 {
			return 1
		} else {
			return 2
		}
	}
	return 0
}

func GetAtbRelates(id uint16) []uint16 {
	_, isExist := AtbNameMap[id]

	if id >= MinAtbID && id <= MaxAtbID {
		return BaseAtbRelate[id]
	} else if isExist {
		temp := getAtbID(id)
		return []uint16{temp.Fir, temp.Sec}
	} else {
		myApp.ErrorLog.Println("error [func GetAtbRelates()]: Input [ID] is Invalid!")
		return nil
	}
}

func GetRestraintTimes(id uint16) ([]TimesPair, []TimesPair) {
	var (
		timesTableAsAtk []TimesPair
		timesTableAsRec []TimesPair
		isExist         bool
	)

	_, isExist = AtbNameMap[id]

	if !isExist || id == AttributeSimpleID {
		return nil, nil
	} else {
		for index := range AtbNameMap {
			if index == AttributeSimpleID {
				continue
			}

			restraintTimes := calculatePairRestraintTimes(id, index)
			if restraintTimes != 1 {
				timesTableAsAtk = append(timesTableAsAtk, TimesPair{
					Attacker:  id,
					Recipient: index,
					Times:     restraintTimes,
				})
			}

			restraintTimes = calculatePairRestraintTimes(index, id)
			if restraintTimes != 1 {
				timesTableAsRec = append(timesTableAsRec, TimesPair{
					Attacker:  index,
					Recipient: id,
					Times:     restraintTimes,
				})
			}
		}
	}

	sort.Slice(timesTableAsAtk, func(i, j int) bool {
		if timesTableAsAtk[i].Times == timesTableAsAtk[j].Times {
			return timesTableAsAtk[i].Recipient > timesTableAsAtk[j].Recipient
		}
		return timesTableAsAtk[i].Times < timesTableAsAtk[j].Times
	})
	sort.Slice(timesTableAsRec, func(i, j int) bool {
		if timesTableAsRec[i].Times == timesTableAsRec[j].Times {
			return timesTableAsRec[i].Attacker > timesTableAsRec[j].Attacker
		}
		return timesTableAsRec[i].Times < timesTableAsRec[j].Times
	})
	return timesTableAsAtk, timesTableAsRec
}

func InitBaseAtbRelates() {
	if BaseAtbRelate == nil {
		myApp.ErrorLog.Panicln("error [func InitBaseAtbRelates()]: Table 'BaseAtbRelate' is not allocated!")
	} else {
		index := 0
		for i := MinAtbID; i <= MaxAtbID; i++ {
			if i == 20 {
				i++
			}
			go processForInitBaseAtbRelates(i, &index)
		}
		return
	}
}

// Calculate the Restraint Times by ATB id
func calculatePairRestraintTimes(atkID uint16, recID uint16) float64 {
	var (
		atk, rec       AtbID
		restraintTimes [5]float64
	)

	atk = getAtbID(atkID)
	rec = getAtbID(recID)

	if atk.Fir != 0 && rec.Fir != 0 {
		setRestraintTimesValueForCalculatePairRestraintTimes(atk.Fir, rec.Fir, &restraintTimes[1])
	}
	if atk.Fir != 0 && rec.Sec != 0 {
		setRestraintTimesValueForCalculatePairRestraintTimes(atk.Fir, rec.Sec, &restraintTimes[2])
	}
	if atk.Sec != 0 && rec.Fir != 0 {
		setRestraintTimesValueForCalculatePairRestraintTimes(atk.Sec, rec.Fir, &restraintTimes[3])
	}
	if atk.Sec != 0 && rec.Sec != 0 {
		setRestraintTimesValueForCalculatePairRestraintTimes(atk.Sec, rec.Sec, &restraintTimes[4])
	}

	restraintTimes[0] = calculateRestraintTimes(atkID, recID, restraintTimes)
	return restraintTimes[0]
}

func calculateRestraintTimes(atkID uint16, recID uint16, restraintTimes [5]float64) float64 {
	if (atkID >= MinAtbID && atkID <= MaxAtbID) && (recID >= MinAtbID && recID <= MaxAtbID) {
		//single to single
		return restraintTimes[1]
	} else if (atkID >= MinAtbID && atkID <= MaxAtbID) && !(recID >= MinAtbID && recID <= MaxAtbID) {
		//single to multiple
		if restraintTimes[1] == 2 && restraintTimes[2] == 2 {
			return 4
		} else if restraintTimes[1] == 0 || restraintTimes[2] == 0 {
			return (restraintTimes[1] + restraintTimes[2]) / 4.0
		} else {
			return (restraintTimes[1] + restraintTimes[2]) / 2.0
		}
	} else if !(atkID >= MinAtbID && atkID <= MaxAtbID) && (recID >= MinAtbID && recID <= MaxAtbID) {
		//multiple to single
		if restraintTimes[1] == 2 && restraintTimes[3] == 2 {
			return 4
		} else if restraintTimes[1] == 0 || restraintTimes[3] == 0 {
			return (restraintTimes[1] + restraintTimes[3]) / 4.0
		} else {
			return (restraintTimes[1] + restraintTimes[3]) / 2.0
		}
	} else {
		//multiple to multiple
		fir := calculatePairRestraintTimes(atkID, recID/100)
		sec := calculatePairRestraintTimes(atkID, recID%100)
		return (fir + sec) / 2.0
	}
}

func getAtbID(id uint16) AtbID {
	ret := AtbID{
		Fir: id / 100,
		Sec: id % 100,
	}

	if ret.Fir == 0 && ret.Sec >= MinAtbID && ret.Sec <= MaxAtbID {
		ret.Fir, ret.Sec = ret.Sec, ret.Fir
		return ret
	} else if ret.Fir < MinAtbID || ret.Fir > MaxAtbID || ret.Sec < MinAtbID || ret.Sec > MaxAtbID || ret.Fir >= ret.Sec {
		//If Fir or Sec out of range, or Fir is no less than Sec, show error message
		myApp.ErrorLog.Println("error [func getAtbID()]:Attribute's ID is invalid")
		return AtbID{}
	}

	return ret
}

func isInPositiveAttackTable(firID uint16, secID uint16) bool {
	//If we can't find the Recipient's first id in Attacker's first attribute's positive table(table's length equals to the return value),
	//then we return a boolean value.
	tableLength := len(BaseAtbRestraintTable[firID].PositiveAttack)
	sortRet := sort.Search(tableLength, func(i int) bool {
		return BaseAtbRestraintTable[firID].PositiveAttack[i] >= secID
	})

	if tableLength == sortRet || BaseAtbRestraintTable[firID].PositiveAttack[sortRet] != secID {
		return false
	} else {
		return true
	}
}

func isInNegativeAttackTable(firID uint16, secID uint16) bool {
	//If we can't find the Recipient's first id in Attacker's first attribute's positive table(table's length equals to the return value),
	//then we return a boolean value.
	tableLength := len(BaseAtbRestraintTable[firID].NegativeAttack)
	sortRet := sort.Search(tableLength, func(i int) bool {
		return BaseAtbRestraintTable[firID].NegativeAttack[i] >= secID
	})

	if tableLength == sortRet || BaseAtbRestraintTable[firID].NegativeAttack[sortRet] != secID {
		return false
	} else {
		return true
	}
}

func processForInitBaseAtbRelates(i int, index *int) {
	lock.Lock()
	BaseAtbRelate[i] = make([]uint16, 16)
	lock.Unlock()

	if BaseAtbRelate[i] == nil {
		myApp.ErrorLog.Panicln("error [func InitBaseAtbRelates()]: Failed to allocate memory!---", i)
	}

	for j := MinAtbID; j <= MaxAtbID; j++ {
		if i == j {
			continue
		}

		var tempID uint16 = uint16(i*100 + j)

		if i > j {
			tempID = uint16(j*100 + i)
		}
		_, isExist := AtbNameMap[tempID]
		if !isExist {
			continue
		}

		lock.Lock()
		BaseAtbRelate[i][*index] = tempID
		lock.Unlock()
		*index++
	}
	*index = 0
}

func setRestraintTimesValueForCalculatePairRestraintTimes(fir uint16, sec uint16, value *float64) {
	if fir == 0 || sec == 0 {
		*value = -1
	} else if isInPositiveAttackTable(fir, sec) {
		*value = 2
	} else if isInNegativeAttackTable(fir, sec) {
		*value = 0.5
	} else if sec != 0 && BaseAtbRestraintTable[fir].InvalidAttack == sec {
		*value = 0
	} else {
		*value = 1
	}
}

/*
{
	11:"火",	12:"水",	13:"草",	14:"飞行",	15:"电",
	16:"地面",	17:"机械",	18:"冰",	19:"超能",	20:"普通",
	21:"战斗",	22:"暗影",	23:"光",	24:"龙",	25:"神秘",
	26:"圣灵",	27:"次元",	28:"远古",	29:"邪灵",	30:"自然",
	31:"王",	32:"混沌",	33:"神灵",	34:"轮回",	35:"虫",
	36:"虚空",

	1114:"火飞行", 1115:"电火", 1118:"冰火", 1119:"火超能", 1121:"战斗火", 1122:"暗影火", 1123:"光火", 1124:"火龙", 1125:"火神秘", 1126:"圣灵火", 1128:"远古火", 1132:"混沌火",
	1219:"水超能", 1221:"水战斗", 1222:"水暗影", 1224:"水龙", 1227:"水次元",
	1319:"草超能", 1321:"草战斗", 1322:"草暗影", 1328:"远古草",
	1415:"电飞行", 1418:"冰飞行", 1419:"飞行超能", 1422:"飞行暗影", 1423:"光飞行", 1424:"飞龙", 1425:"飞行神秘", 1432:"飞行混沌",
	1518:"电冰", 1521:"电战斗", 1522:"电暗影", 1524:"电龙", 1526:"圣灵电",1527:"电次元",
	1617:"机械地面", 1619:"地面超能", 1621:"地面战斗", 1622:"地面暗影", 1624:"地面龙", 1625:"地面神秘", 1626:"圣灵地面", 1628:"远古地面", 1632:"混沌地面",
	1719:"机械超能", 1721:"机械战斗", 1724:"机械龙", 1727:"机械次元", 1729:"机械邪灵",
	1819:"冰超能", 1822:"冰暗影", 1823:"冰光", 1824:"冰龙", 1825:"冰神秘", 1830:"自然冰", 1832:"混沌冰",
	1923:"光超能", 1925:"神秘超能", 1926:"圣灵超能", 1932:"混沌超能",
	2122:"战斗暗影", 2123:"光战斗", 2124:"战斗龙", 2125:"神秘战斗", 2126:"圣灵战斗", 2127:"次元战斗", 2128:"远古战斗", 2130:"自然战斗", 2132:"混沌战斗",
	2223:"光暗影", 2224:"暗影龙", 2225:"神秘暗影", 2226:"圣灵暗影", 2229:"邪灵暗影", 2232:"混沌暗影",
	2325:"光神秘", 2326:"圣灵光", 2327:"光次元",
	2428:"远古龙", 2429:"邪灵龙", 2430:"自然龙", 2432:"混沌龙",
	2526:"圣灵神秘", 2528:"远古神秘", 2529:"邪灵神秘", 2530:"自然神秘",
	2627:"圣灵次元", 2628:"圣灵远古", 2630:"圣灵自然",  2632:"混沌圣灵", 2634:"圣灵轮回",
	2732:"混沌次元",
	2829:"远古邪灵", 2832:"混沌远古",
	2932:"混沌邪灵", 2936:"虚空邪灵",
	3236:"虚空混沌",
}
*/
