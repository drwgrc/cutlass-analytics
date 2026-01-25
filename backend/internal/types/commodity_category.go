package types

type CommodityCategory string

const (
	CommodityCategoryBasic    CommodityCategory = "basic"    // Wood, Iron, Hemp, etc.
	CommodityCategoryHerb     CommodityCategory = "herb"     // Cowslip, Indigo, etc.
	CommodityCategoryMineral  CommodityCategory = "mineral"  // Chalcocite, Cubanite, etc.
	CommodityCategoryForaged  CommodityCategory = "foraged"  // Forageable items
	CommodityCategoryRefined  CommodityCategory = "refined"  // Processed goods (cloth, dye, etc.)
	CommodityCategoryShipSupply CommodityCategory = "ship_supply" // Rum, Cannonballs
)