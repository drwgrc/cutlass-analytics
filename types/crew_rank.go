package types

type CrewRank string

const (
	CrewRankSailors        CrewRank = "Sailors"
	CrewRankMostlyHarmless CrewRank = "Mostly Harmless"
	CrewRankScurvyDogs     CrewRank = "Scurvy Dogs"
	CrewRankScoundrels     CrewRank = "Scoundrels"
	CrewRankBlaggards      CrewRank = "Blaggards"
	CrewRankDreadPirates   CrewRank = "Dread Pirates"
	CrewRankSeaLords       CrewRank = "Sea Lords"
	CrewRankImperials      CrewRank = "Imperials"
)

func (r CrewRank) String() string {
	return string(r)
}

func (r CrewRank) Order() int {
	switch r {
	case CrewRankSailors:
		return 1
	case CrewRankMostlyHarmless:
		return 2
	case CrewRankScurvyDogs:
		return 3
	case CrewRankScoundrels:
		return 4
	case CrewRankBlaggards:
		return 5
	case CrewRankDreadPirates:
		return 6
	case CrewRankSeaLords:
		return 7
	case CrewRankImperials:
		return 8
	}
	return 0
}