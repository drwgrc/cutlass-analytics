package types

type Ocean string

const (
	OceanEmerald  Ocean = "emerald"
	OceanMeridian Ocean = "meridian"
	OceanCerulean Ocean = "cerulean"
	OceanObsidian Ocean = "obsidian"
)

func (o Ocean) String() string {
	return string(o)
}

func (o Ocean) IsValid() bool {
	switch o {
	case OceanEmerald, OceanMeridian, OceanCerulean, OceanObsidian:
		return true
	}
	return false
}