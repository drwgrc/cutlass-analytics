package types

type CrewPVPResult struct {
	CrewID         uint
	GameCrewID     uint64
	Name           string
	Ocean          Ocean
	FlagID         *uint
	FlagName       string
	CrewRank       CrewRank
	TotalPVPWins   int
	TotalPVPLosses int
}

func (r *CrewPVPResult) WinRate() float64 {
	total := r.TotalPVPWins + r.TotalPVPLosses
	if total == 0 {
		return 0
	}
	return float64(r.TotalPVPWins) / float64(total) * 100
}

func (r *CrewPVPResult) TotalBattles() int {
	return r.TotalPVPWins + r.TotalPVPLosses
}