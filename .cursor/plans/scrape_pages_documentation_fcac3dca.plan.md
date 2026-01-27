---
name: Scrape Pages Documentation
overview: Create a comprehensive documentation file listing all pages being scraped, their URLs, and the expected values that need to be extracted from each page with HTML selectors.
todos: []
isProject: false
---

# Scrape Pages Documentation

This document lists all pages being scraped by the system and the values expected to be extracted from each page. For each page, provide the HTML selector (CSS selector or XPath) that locates the value.

## 1. Island Info Page ✅

**URL Pattern:** `https://{ocean}.puzzlepirates.com/yoweb/island/info.wm?islandid={id}` (looped from 0 to 120)

**Parser Function:** `ParseIslandInfo` in `backend/internal/scraper/parser.go`

**Scraping Approach:**

- Loop through island IDs from 0 to 120
- Visit each individual island page
- Skip islands with text "Shiver me timbers: The island is uncolonized."
- Parse colonized islands only

**HTML Structure:**

```
<center>
  <center>
    <font size="+1">Turtle Island</font><br>
    Population: 57<br>
    Located in the Diamond archipelago.<br>
    Governor: <a href="/yoweb/pirate.wm?classic=true&target=Darkseid">Darkseid</a><br>
    Property tax: 20%<br>
  </center>
  Ruled by <a href="/yoweb/flag/info.wm?flagid=10013644&classic=$classic">Black Flag Inc</a><br>
  Exports:
    Wood
    ,
    Iron
    ,
    Stone
  <br><br>
</center>
```

**Expected Values:**

- `GameIslandID` (uint64) - From URL parameter `islandid={id}`
- `Name` (string) - Island name from `<font size="+1">` tag
- `Size` (IslandSize enum) - **NOT AVAILABLE ON THIS PAGE**
- `Population` (int) - Extracted from text pattern `Population: [number]`
- `Archipelago` (string) - Extracted from text pattern `Located in the [name] archipelago.`
- `IsColonized` (bool) - Always `true` (uncolonized islands are skipped)
- `Commodities` ([]string) - Extracted from "Exports:" section, comma-separated list
- `GovernorName` (string) - From link text in "Governor: <a>...</a>"
- `GovernorFlag` (string) - Flag ID from "Ruled by <a href="...flagid=...">" link

**Implementation Status:**

- [x] GameIslandID from URL parameter
- [x] Island name from `<font size="+1">` tag
- [x] Population from text pattern `Population: [number]`
- [x] Archipelago from text pattern `Located in the [name] archipelago.`
- [x] Commodities from "Exports:" section (comma-separated)
- [x] IsColonized set to `true` (uncolonized islands skipped)
- [x] Governor name from link after "Governor:"
- [x] Governor flag ID from "Ruled by" link
- [ ] Size - Not available on this page

---

## 2. Tax Rates Page

**URL Pattern:** `https://{ocean}.puzzlepirates.com/yoweb/econ/taxrates.wm`

**Parser Function:** `ParseTaxRates` in `backend/internal/scraper/parser.go`

**Expected Values:**

- `CommodityName` (string) - Name of the commodity
- `TaxValue` (int) - Tax rate value (numeric, may have formatting like commas)

**Current Parsing Logic:** Parses table rows, expects first cell = commodity name, second cell = tax value.

**HTML Selector Needed For:**

- [ ] Commodity name (table cell)
- [ ] Tax value (table cell)

---

## 3. Crew Fame List Page

**URL Pattern:** `https://{ocean}.puzzlepirates.com/ratings/top_fame_97.html`

**Parser Function:** `ParseCrewFameList` in `backend/internal/scraper/parser.go`

**Expected Values:**

- `CrewID` (uint64) - Crew ID extracted from URL parameter `crewid=` in links
- `Name` (string) - Crew name (from link text)
- `FameLevel` (FameLevel enum) - One of: "Obscure", "Rumored", "Noted", "Recognized", "Distinguished", "Celebrated", "Eminent", "Renowned", "Illustrious"
- `Rank` (*int) - Optional rank number (usually first table cell)

**Current Parsing Logic:** Parses table rows, extracts crew ID from `a[href*='crewid=']` links, searches row text for fame level keywords.

**HTML Selector Needed For:**

- [ ] Crew ID from link href
- [ ] Crew name
- [ ] Fame level
- [ ] Rank number

---

## 4. Crew Info Page

**URL Pattern:** `https://{ocean}.puzzlepirates.com/yoweb/crew/info.wm?crewid={crewID}`

**Parser Function:** `ParseCrewInfo` in `backend/internal/scraper/parser.go`

**Expected Values:**

- `Name` (string) - Crew name (from h1, h2, or title)
- `FlagID` (*uint64) - Optional flag ID extracted from URL parameter `flagid=` in links (null if crew is independent)
- `FlagName` (string) - Flag name (from link text)

**Current Parsing Logic:** Extracts name from headings/title, extracts flag ID from `a[href*='flagid=']` links.

**HTML Selector Needed For:**

- [ ] Crew name
- [ ] Flag ID from link href (if crew is in a flag)
- [ ] Flag name

---

## 5. Crew Battle Info Page

**URL Pattern:** `https://{ocean}.puzzlepirates.com/yoweb/crew/battleinfo.wm?crewid={crewID}&classic=$classic`

**Parser Function:** `ParseCrewBattleInfo` in `backend/internal/scraper/parser.go`

**Expected Values:**

- `CrewRank` (CrewRank enum) - One of: "Sailors", "Mostly Harmless", "Scurvy Dogs", "Scoundrels", "Blaggards", "Dread Pirates", "Sea Lords", "Imperials"
- `TotalPVPWins` (int) - Total PVP wins count
- `TotalPVPLosses` (int) - Total PVP losses count

**Current Parsing Logic:** Searches all text for rank keywords, uses regex to find "Wins: X" and "Losses: X" patterns.

**HTML Selector Needed For:**

- [ ] Crew rank
- [ ] Total PVP wins
- [ ] Total PVP losses

---

## 6. Flag Fame List Page

**URL Pattern:** `https://{ocean}.puzzlepirates.com/ratings/top_fame_112.html`

**Parser Function:** `ParseFlagFameList` in `backend/internal/scraper/parser.go`

**Expected Values:**

- `FlagID` (uint64) - Flag ID extracted from URL parameter `flagid=` in links
- `Name` (string) - Flag name (from link text)
- `FameLevel` (FameLevel enum) - One of: "Obscure", "Rumored", "Noted", "Recognized", "Distinguished", "Celebrated", "Eminent", "Renowned", "Illustrious"
- `Rank` (*int) - Optional rank number (usually first table cell)

**Current Parsing Logic:** Parses table rows, extracts flag ID from `a[href*='flagid=']` links, searches row text for fame level keywords.

**HTML Selector Needed For:**

- [ ] Flag ID from link href
- [ ] Flag name
- [ ] Fame level
- [ ] Rank number

---

## 7. Flag Info Page

**URL Pattern:** `https://{ocean}.puzzlepirates.com/yoweb/flag/info.wm?flagid={flagID}`

**Parser Function:** `ParseFlagInfo` in `backend/internal/scraper/parser.go`

**Expected Values:**

- `Name` (string) - Flag name (from h1, h2, or title)

**Current Parsing Logic:** Extracts name from headings/title.

**HTML Selector Needed For:**

- [ ] Flag name

---

## Notes

- All URLs use `{ocean}` placeholder which is replaced with the actual ocean name (e.g., "emerald", "meridian")
- Some parsers use regex patterns on text content rather than specific HTML selectors
- The current implementation is somewhat fragile and relies on text matching - providing proper HTML selectors will make it more robust
- Island Info Page (`ParseIslandInfo`) exists in the parser but is not currently used by the scraper - it may be needed for more detailed island information