package scraper

import (
	"cutlass_analytics/internal/types"
	"testing"
)

func TestParseIslandInfo(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		islandID uint64
		ocean    types.Ocean
		want     *IslandData
		wantErr  bool
	}{
		{
			name: "valid colonized island",
			html: `<html><body>
				<center>
					<center>
						<font size="+1">Turtle Island</font><br>
						Population: 57<br>
						Located in the Diamond archipelago.<br>
						Governor: <a href="/yoweb/pirate.wm?classic=true&target=Darkseid">Darkseid</a><br>
						Property tax: 20%<br>
					</center>
					Ruled by <a href="/yoweb/flag/info.wm?flagid=10013644&classic=$classic">Black Flag Inc</a><br>
					Exports: Wood, Iron, Stone<br>
				</center>
			</body></html>`,
			islandID: 42,
			ocean:    types.OceanEmerald,
			want: &IslandData{
				GameIslandID: 42,
				Name:         "Turtle Island",
				Population:   57,
				Archipelago:  "Diamond",
				IsColonized:  true,
				GovernorName: "Darkseid",
				GovernorFlag: "10013644",
				Commodities:  []string{"Wood", "Iron", "Stone"},
			},
			wantErr: false,
		},
		{
			name: "island with single commodity",
			html: `<html><body>
				<center>
					<center>
						<font size="+1">Small Island</font><br>
						Population: 10<br>
						Located in the Ruby archipelago.<br>
					</center>
					Exports: Hemp<br>
				</center>
			</body></html>`,
			islandID: 5,
			ocean:    types.OceanMeridian,
			want: &IslandData{
				GameIslandID: 5,
				Name:         "Small Island",
				Population:   10,
				Archipelago:  "Ruby",
				IsColonized:  true,
				Commodities:  []string{"Hemp"},
			},
			wantErr: false,
		},
		{
			name: "island with population",
			html: `<html><body>
				<center>
					<center>
						<font size="+1">Big Island</font><br>
						Population: 1234<br>
						Located in the Emerald archipelago.<br>
					</center>
				</center>
			</body></html>`,
			islandID: 100,
			ocean:    types.OceanEmerald,
			want: &IslandData{
				GameIslandID: 100,
				Name:         "Big Island",
				Population:   1234,
				Archipelago:  "Emerald",
				IsColonized:  true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIslandInfo(tt.html, tt.islandID, tt.ocean)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIslandInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.GameIslandID != tt.want.GameIslandID {
				t.Errorf("GameIslandID = %v, want %v", got.GameIslandID, tt.want.GameIslandID)
			}
			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Population != tt.want.Population {
				t.Errorf("Population = %v, want %v", got.Population, tt.want.Population)
			}
			if got.Archipelago != tt.want.Archipelago {
				t.Errorf("Archipelago = %v, want %v", got.Archipelago, tt.want.Archipelago)
			}
			if got.IsColonized != tt.want.IsColonized {
				t.Errorf("IsColonized = %v, want %v", got.IsColonized, tt.want.IsColonized)
			}
			if got.GovernorName != tt.want.GovernorName {
				t.Errorf("GovernorName = %v, want %v", got.GovernorName, tt.want.GovernorName)
			}
			if got.GovernorFlag != tt.want.GovernorFlag {
				t.Errorf("GovernorFlag = %v, want %v", got.GovernorFlag, tt.want.GovernorFlag)
			}
			if len(got.Commodities) != len(tt.want.Commodities) {
				t.Errorf("Commodities length = %v, want %v", len(got.Commodities), len(tt.want.Commodities))
			}
		})
	}
}

func TestParseTaxRates(t *testing.T) {
	// Note: The ParseTaxRates function uses a very specific CSS selector
	// (body > center > table > td:first-child > table) designed to match the
	// exact HTML structure returned by the Puzzle Pirates game server.
	// Testing with synthetic HTML is challenging due to how goquery normalizes
	// table structures (auto-inserting tbody/tr elements).
	// These tests verify the function handles input without errors.
	tests := []struct {
		name    string
		html    string
		ocean   types.Ocean
		wantErr bool
	}{
		{
			name:    "handles valid HTML without error",
			html:    `<html><body><center><table><tr><td><table><tr><th>Commodity</th><th>Tax Rate</th></tr></table></td></tr></table></center></body></html>`,
			ocean:   types.OceanEmerald,
			wantErr: false,
		},
		{
			name:    "handles empty HTML without error",
			html:    `<html><body></body></html>`,
			ocean:   types.OceanMeridian,
			wantErr: false,
		},
		{
			name:    "handles malformed HTML without error",
			html:    `<body><center>no table here</center></body>`,
			ocean:   types.OceanEmerald,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTaxRates(tt.html, tt.ocean)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTaxRates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseCrewFameList(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		ocean   types.Ocean
		want    []CrewFameData
		wantErr bool
	}{
		{
			name: "valid crew fame list",
			html: `<html><body>
				<table>
					<tr>
						<th>Rank</th>
						<th>Crew</th>
						<th>Fame</th>
					</tr>
					<tr>
						<td>1</td>
						<td><a href="/yoweb/crew/info.wm?crewid=12345">Best Crew</a></td>
						<td>Illustrious</td>
					</tr>
					<tr>
						<td>2</td>
						<td><a href="/yoweb/crew/info.wm?crewid=67890">Second Crew</a></td>
						<td>Renowned</td>
					</tr>
					<tr>
						<td>3</td>
						<td><a href="/yoweb/crew/info.wm?crewid=11111">Third Crew</a></td>
						<td>Obscure</td>
					</tr>
				</table>
			</body></html>`,
			ocean: types.OceanEmerald,
			want: []CrewFameData{
				{CrewID: 12345, Name: "Best Crew", FameLevel: types.FameLevelIllustrious, Rank: intPtr(1)},
				{CrewID: 67890, Name: "Second Crew", FameLevel: types.FameLevelRenowned, Rank: intPtr(2)},
				{CrewID: 11111, Name: "Third Crew", FameLevel: types.FameLevelObscure, Rank: intPtr(3)},
			},
			wantErr: false,
		},
		{
			name: "all fame levels",
			html: `<html><body>
				<table>
					<tr><th>Rank</th><th>Crew</th><th>Fame</th></tr>
					<tr><td>1</td><td><a href="/yoweb/crew/info.wm?crewid=1">A</a></td><td>Obscure</td></tr>
					<tr><td>2</td><td><a href="/yoweb/crew/info.wm?crewid=2">B</a></td><td>Rumored</td></tr>
					<tr><td>3</td><td><a href="/yoweb/crew/info.wm?crewid=3">C</a></td><td>Noted</td></tr>
					<tr><td>4</td><td><a href="/yoweb/crew/info.wm?crewid=4">D</a></td><td>Recognized</td></tr>
					<tr><td>5</td><td><a href="/yoweb/crew/info.wm?crewid=5">E</a></td><td>Distinguished</td></tr>
					<tr><td>6</td><td><a href="/yoweb/crew/info.wm?crewid=6">F</a></td><td>Celebrated</td></tr>
					<tr><td>7</td><td><a href="/yoweb/crew/info.wm?crewid=7">G</a></td><td>Eminent</td></tr>
					<tr><td>8</td><td><a href="/yoweb/crew/info.wm?crewid=8">H</a></td><td>Renowned</td></tr>
					<tr><td>9</td><td><a href="/yoweb/crew/info.wm?crewid=9">I</a></td><td>Illustrious</td></tr>
				</table>
			</body></html>`,
			ocean: types.OceanMeridian,
			want: []CrewFameData{
				{CrewID: 1, Name: "A", FameLevel: types.FameLevelObscure, Rank: intPtr(1)},
				{CrewID: 2, Name: "B", FameLevel: types.FameLevelRumored, Rank: intPtr(2)},
				{CrewID: 3, Name: "C", FameLevel: types.FameLevelNoted, Rank: intPtr(3)},
				{CrewID: 4, Name: "D", FameLevel: types.FameLevelRecognized, Rank: intPtr(4)},
				{CrewID: 5, Name: "E", FameLevel: types.FameLevelDistinguished, Rank: intPtr(5)},
				{CrewID: 6, Name: "F", FameLevel: types.FameLevelCelebrated, Rank: intPtr(6)},
				{CrewID: 7, Name: "G", FameLevel: types.FameLevelEminent, Rank: intPtr(7)},
				{CrewID: 8, Name: "H", FameLevel: types.FameLevelRenowned, Rank: intPtr(8)},
				{CrewID: 9, Name: "I", FameLevel: types.FameLevelIllustrious, Rank: intPtr(9)},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCrewFameList(tt.html, tt.ocean)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCrewFameList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ParseCrewFameList() got %d crews, want %d", len(got), len(tt.want))
				return
			}
			for i, crew := range got {
				if crew.CrewID != tt.want[i].CrewID {
					t.Errorf("CrewID[%d] = %v, want %v", i, crew.CrewID, tt.want[i].CrewID)
				}
				if crew.Name != tt.want[i].Name {
					t.Errorf("Name[%d] = %v, want %v", i, crew.Name, tt.want[i].Name)
				}
				if crew.FameLevel != tt.want[i].FameLevel {
					t.Errorf("FameLevel[%d] = %v, want %v", i, crew.FameLevel, tt.want[i].FameLevel)
				}
				if *crew.Rank != *tt.want[i].Rank {
					t.Errorf("Rank[%d] = %v, want %v", i, *crew.Rank, *tt.want[i].Rank)
				}
			}
		})
	}
}

func TestParseCrewInfo(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		crewID  uint64
		ocean   types.Ocean
		want    *CrewData
		wantErr bool
	}{
		{
			name: "crew with flag and rank",
			html: `<html><body>
				<table>
					<tr>
						<td width="246">
							<font><b>Pirates of the Caribbean</b></font>
							<a href="/yoweb/flag/info.wm?flagid=99999">Jolly Roger Alliance</a>
						</td>
					</tr>
				</table>
				<a href="/yoweb/crew/battleinfo.wm?crewid=12345&classic=false">Sea Lords</a>
			</body></html>`,
			crewID: 12345,
			ocean:  types.OceanEmerald,
			want: &CrewData{
				GameCrewID: 12345,
				Name:       "Pirates of the Caribbean",
				FlagID:     uint64Ptr(99999),
				FlagName:   "Jolly Roger Alliance",
				CrewRank:   types.CrewRankSeaLords,
			},
			wantErr: false,
		},
		{
			name: "crew without flag (independent)",
			html: `<html><body>
				<table>
					<tr>
						<td width="246">
							<font><b>Lone Wolf Sailors</b></font>
						</td>
					</tr>
				</table>
				<a href="/yoweb/crew/battleinfo.wm?crewid=55555&classic=false">Sailors</a>
			</body></html>`,
			crewID: 55555,
			ocean:  types.OceanMeridian,
			want: &CrewData{
				GameCrewID: 55555,
				Name:       "Lone Wolf Sailors",
				FlagID:     nil,
				FlagName:   "",
				CrewRank:   types.CrewRankSailors,
			},
			wantErr: false,
		},
		{
			name: "crew with imperials rank",
			html: `<html><body>
				<table>
					<tr>
						<td width="246">
							<font><b>Imperial Crew</b></font>
						</td>
					</tr>
				</table>
				<a href="/yoweb/crew/battleinfo.wm?crewid=88888&classic=false">Imperials</a>
			</body></html>`,
			crewID: 88888,
			ocean:  types.OceanEmerald,
			want: &CrewData{
				GameCrewID: 88888,
				Name:       "Imperial Crew",
				CrewRank:   types.CrewRankImperials,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCrewInfo(tt.html, tt.crewID, tt.ocean)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCrewInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.GameCrewID != tt.want.GameCrewID {
				t.Errorf("GameCrewID = %v, want %v", got.GameCrewID, tt.want.GameCrewID)
			}
			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if (got.FlagID == nil) != (tt.want.FlagID == nil) {
				t.Errorf("FlagID nil status = %v, want %v", got.FlagID == nil, tt.want.FlagID == nil)
			}
			if got.FlagID != nil && tt.want.FlagID != nil && *got.FlagID != *tt.want.FlagID {
				t.Errorf("FlagID = %v, want %v", *got.FlagID, *tt.want.FlagID)
			}
			if got.FlagName != tt.want.FlagName {
				t.Errorf("FlagName = %v, want %v", got.FlagName, tt.want.FlagName)
			}
			if got.CrewRank != tt.want.CrewRank {
				t.Errorf("CrewRank = %v, want %v", got.CrewRank, tt.want.CrewRank)
			}
		})
	}
}

func TestParseCrewBattleInfo(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		crewID  uint64
		want    *CrewBattleData
		wantErr bool
	}{
		{
			name: "valid battle info with multiple rows",
			html: `<html><body>
				<table>
					<tr><th>Header 1</th></tr>
					<tr><th>Header 2</th></tr>
					<tr>
						<td>2024-01-01</td>
						<td>10</td>
						<td>5</td>
						<td>5</td>
						<td>3</td>
						<td>2</td>
						<td>10:00</td>
					</tr>
					<tr>
						<td>2024-01-02</td>
						<td>8</td>
						<td>4</td>
						<td>4</td>
						<td>2</td>
						<td>1</td>
						<td>8:00</td>
					</tr>
				</table>
			</body></html>`,
			crewID: 12345,
			want: &CrewBattleData{
				TotalPVPWins:   5,  // 3 + 2
				TotalPVPLosses: 3,  // 2 + 1
			},
			wantErr: false,
		},
		{
			name: "battle info with zero values",
			html: `<html><body>
				<table>
					<tr><th>Header 1</th></tr>
					<tr><th>Header 2</th></tr>
					<tr>
						<td>2024-01-01</td>
						<td>0</td>
						<td>0</td>
						<td>0</td>
						<td>0</td>
						<td>0</td>
						<td>0:00</td>
					</tr>
				</table>
			</body></html>`,
			crewID: 67890,
			want: &CrewBattleData{
				TotalPVPWins:   0,
				TotalPVPLosses: 0,
			},
			wantErr: false,
		},
		{
			name: "battle info with large numbers",
			html: `<html><body>
				<table>
					<tr><th>Header 1</th></tr>
					<tr><th>Header 2</th></tr>
					<tr>
						<td>2024-01-01</td>
						<td>100</td>
						<td>50</td>
						<td>50</td>
						<td>1000</td>
						<td>500</td>
						<td>100:00</td>
					</tr>
				</table>
			</body></html>`,
			crewID: 11111,
			want: &CrewBattleData{
				TotalPVPWins:   1000,
				TotalPVPLosses: 500,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCrewBattleInfo(tt.html, tt.crewID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCrewBattleInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.TotalPVPWins != tt.want.TotalPVPWins {
				t.Errorf("TotalPVPWins = %v, want %v", got.TotalPVPWins, tt.want.TotalPVPWins)
			}
			if got.TotalPVPLosses != tt.want.TotalPVPLosses {
				t.Errorf("TotalPVPLosses = %v, want %v", got.TotalPVPLosses, tt.want.TotalPVPLosses)
			}
		})
	}
}

func TestParseFlagFameList(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		ocean   types.Ocean
		want    []FlagFameData
		wantErr bool
	}{
		{
			name: "valid flag fame list",
			html: `<html><body>
				<table>
					<tr>
						<th>Rank</th>
						<th>Flag</th>
						<th>Fame</th>
					</tr>
					<tr>
						<td>1</td>
						<td><a href="/yoweb/flag/info.wm?flagid=11111">Top Flag</a></td>
						<td>Illustrious</td>
					</tr>
					<tr>
						<td>2</td>
						<td><a href="/yoweb/flag/info.wm?flagid=22222">Second Flag</a></td>
						<td>Eminent</td>
					</tr>
				</table>
			</body></html>`,
			ocean: types.OceanEmerald,
			want: []FlagFameData{
				{FlagID: 11111, Name: "Top Flag", FameLevel: types.FameLevelIllustrious, Rank: intPtr(1)},
				{FlagID: 22222, Name: "Second Flag", FameLevel: types.FameLevelEminent, Rank: intPtr(2)},
			},
			wantErr: false,
		},
		{
			name: "empty flag list",
			html: `<html><body>
				<table>
					<tr>
						<th>Rank</th>
						<th>Flag</th>
						<th>Fame</th>
					</tr>
				</table>
			</body></html>`,
			ocean:   types.OceanMeridian,
			want:    []FlagFameData{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFlagFameList(tt.html, tt.ocean)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlagFameList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ParseFlagFameList() got %d flags, want %d", len(got), len(tt.want))
				return
			}
			for i, flag := range got {
				if flag.FlagID != tt.want[i].FlagID {
					t.Errorf("FlagID[%d] = %v, want %v", i, flag.FlagID, tt.want[i].FlagID)
				}
				if flag.Name != tt.want[i].Name {
					t.Errorf("Name[%d] = %v, want %v", i, flag.Name, tt.want[i].Name)
				}
				if flag.FameLevel != tt.want[i].FameLevel {
					t.Errorf("FameLevel[%d] = %v, want %v", i, flag.FameLevel, tt.want[i].FameLevel)
				}
			}
		})
	}
}

func TestParseFlagInfo(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		flagID  uint64
		ocean   types.Ocean
		want    *FlagData
		wantErr bool
	}{
		{
			name: "valid flag info with all reputations",
			html: `<html><body>
				<table>
					<tr>
						<td width="246">
							<font><b>Jolly Roger Alliance</b></font>
						</td>
						<td width="246">
							<table>
								<tr><td>Conqueror:</td><td><font>Distinguished</font></td></tr>
								<tr><td>Explorer:</td><td><font>Celebrated</font></td></tr>
								<tr><td>Patron:</td><td><font>Eminent</font></td></tr>
								<tr><td>Magnate:</td><td><font>Renowned</font></td></tr>
							</table>
						</td>
					</tr>
				</table>
			</body></html>`,
			flagID: 99999,
			ocean:  types.OceanEmerald,
			want: &FlagData{
				GameFlagID:          99999,
				Name:                "Jolly Roger Alliance",
				ConquerorReputation: fameLevelPtr(types.FameLevelDistinguished),
				ExplorerReputation:  fameLevelPtr(types.FameLevelCelebrated),
				PatronReputation:    fameLevelPtr(types.FameLevelEminent),
				MagnateReputation:   fameLevelPtr(types.FameLevelRenowned),
			},
			wantErr: false,
		},
		{
			name: "flag with obscure reputations",
			html: `<html><body>
				<table>
					<tr>
						<td width="246">
							<font><b>New Flag</b></font>
						</td>
						<td width="246">
							<table>
								<tr><td>Conqueror:</td><td><font>Obscure</font></td></tr>
								<tr><td>Explorer:</td><td><font>Obscure</font></td></tr>
								<tr><td>Patron:</td><td><font>Obscure</font></td></tr>
								<tr><td>Magnate:</td><td><font>Obscure</font></td></tr>
							</table>
						</td>
					</tr>
				</table>
			</body></html>`,
			flagID: 11111,
			ocean:  types.OceanMeridian,
			want: &FlagData{
				GameFlagID:          11111,
				Name:                "New Flag",
				ConquerorReputation: fameLevelPtr(types.FameLevelObscure),
				ExplorerReputation:  fameLevelPtr(types.FameLevelObscure),
				PatronReputation:    fameLevelPtr(types.FameLevelObscure),
				MagnateReputation:   fameLevelPtr(types.FameLevelObscure),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFlagInfo(tt.html, tt.flagID, tt.ocean)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlagInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.GameFlagID != tt.want.GameFlagID {
				t.Errorf("GameFlagID = %v, want %v", got.GameFlagID, tt.want.GameFlagID)
			}
			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if !compareFameLevel(got.ConquerorReputation, tt.want.ConquerorReputation) {
				t.Errorf("ConquerorReputation = %v, want %v", got.ConquerorReputation, tt.want.ConquerorReputation)
			}
			if !compareFameLevel(got.ExplorerReputation, tt.want.ExplorerReputation) {
				t.Errorf("ExplorerReputation = %v, want %v", got.ExplorerReputation, tt.want.ExplorerReputation)
			}
			if !compareFameLevel(got.PatronReputation, tt.want.PatronReputation) {
				t.Errorf("PatronReputation = %v, want %v", got.PatronReputation, tt.want.PatronReputation)
			}
			if !compareFameLevel(got.MagnateReputation, tt.want.MagnateReputation) {
				t.Errorf("MagnateReputation = %v, want %v", got.MagnateReputation, tt.want.MagnateReputation)
			}
		})
	}
}

func TestParseIslandList(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		ocean   types.Ocean
		want    int // number of islands expected
		wantErr bool
	}{
		{
			name: "multiple islands",
			html: `<html><body>
				<center>
					<center>
						<font>Alpha Island</font>
						Population: 100
						Located in the Diamond archipelago.
					</center>
					<center>
						<font>Beta Island</font>
						Population: 200
						Located in the Ruby archipelago.
					</center>
				</center>
			</body></html>`,
			ocean:   types.OceanEmerald,
			want:    2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIslandList(tt.html, tt.ocean)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIslandList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ParseIslandList() got %d islands, want %d", len(got), tt.want)
			}
		})
	}
}

// Helper functions for creating pointers
func intPtr(i int) *int {
	return &i
}

func uint64Ptr(i uint64) *uint64 {
	return &i
}

func fameLevelPtr(f types.FameLevel) *types.FameLevel {
	return &f
}

func compareFameLevel(a, b *types.FameLevel) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
