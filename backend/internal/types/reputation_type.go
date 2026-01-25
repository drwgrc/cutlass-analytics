package types

type ReputationType string

const (
	ReputationConqueror ReputationType = "Conqueror"
	ReputationExplorer  ReputationType = "Explorer"
	ReputationPatron    ReputationType = "Patron"
	ReputationMagnate   ReputationType = "Magnate"
)

func (r ReputationType) String() string {
	return string(r)
}

func (r ReputationType) IsValid() bool {
	switch r {
	case ReputationConqueror, ReputationExplorer, ReputationPatron, ReputationMagnate:
		return true
	}
	return false
}