package ascii

// HorizontalMargins ...
type HorizontalMargins struct {
	Left  int
	Right int
}

// Width ...
func (margins HorizontalMargins) Width(contentWidth int) int {
	return margins.Left + margins.Right + contentWidth
}

// VerticalMargins ...
type VerticalMargins struct {
	Top    int
	Bottom int
}

// StoneMargins ...
type StoneMargins struct {
	HorizontalMargins
	VerticalMargins
}

// LegendMargins ...
type LegendMargins struct {
	Column VerticalMargins
	Row    HorizontalMargins
}

// Margins ...
type Margins struct {
	Stone  StoneMargins
	Legend LegendMargins
	Board  VerticalMargins
}
