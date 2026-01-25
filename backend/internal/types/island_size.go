package types

type IslandSize string

const (
	IslandSizeOutpost IslandSize = "outpost"
	IslandSizeMedium  IslandSize = "medium"
	IslandSizeLarge   IslandSize = "large"
)

func (s IslandSize) String() string {
	return string(s)
}

func (s IslandSize) MaxBuildings() int {
	switch s {
	case IslandSizeOutpost:
		return 2
	case IslandSizeMedium:
		return 6
	case IslandSizeLarge:
		return -1
	}
	return 0
}