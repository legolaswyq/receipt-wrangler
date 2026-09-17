package repositories

// CategoryColorPalette is the curated, categorical palette used to assign a
// default color to a category that is created (or backfilled) without one. It is
// mirrored on the desktop client's color picker so hand-picked and auto-assigned
// colors come from the same set.
var CategoryColorPalette = []string{
	// Tableau 10 (the original set) — kept first so already-assigned colors don't churn.
	"#4E79A7",
	"#F28E2B",
	"#E15759",
	"#76B7B2",
	"#59A14F",
	"#EDC948",
	"#B07AA1",
	"#FF9DA7",
	"#9C755F",
	"#BAB0AC",
	// Tableau 20 extension — appended so installs with more categories still get
	// distinct colors.
	"#A0CBE8",
	"#FFBE7D",
	"#8CD17D",
	"#B6992D",
	"#499894",
	"#86BCB6",
	"#FF9D9A",
	"#79706E",
	"#D37295",
	"#D4A6C8",
}

// nextCategoryColor returns the first palette color not already used by an
// existing category. Once every palette color is in use it wraps by count, so
// new categories keep getting a deterministic, evenly-distributed color rather
// than always the first one.
func nextCategoryColor(usedColors []string) string {
	used := make(map[string]bool, len(usedColors))
	for _, c := range usedColors {
		used[c] = true
	}

	for _, color := range CategoryColorPalette {
		if !used[color] {
			return color
		}
	}

	return CategoryColorPalette[len(usedColors)%len(CategoryColorPalette)]
}
