package constant

var SportCss = map[rune]string{
	'T': "#00FF00", // Type
	'V': "#FF0000", // Var
	'C': "#FF00FF", // Const
}

const (
	Line_Normal string = "--"
	Line_Splite string = ".."
)

const (
	Arrow_ExpandLeft     string = "<|"
	Arrow_ExpandRight    string = "|>"
	Arrow_Composition    string = "*"
	Arrow_Polymerization string = "o"
)

const (
	Visibility_Private        rune = '-'
	Visibility_Protected      rune = '#'
	Visibility_PackagePrivate rune = '~'
	Visibility_public         rune = '+'
)

const (
	Color_Red   string = "#FF0000"
	Color_Green string = "green"
	Color_blue  string = "blue"
)

const (
	LineStyle_Bold   string = "bold"
	LineStyle_Dashed string = "dashed"
	LineStyle_Dotted string = "dotted"
	LineStyle_Hidden string = "hidden"
)
