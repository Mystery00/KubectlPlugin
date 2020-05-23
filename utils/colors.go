package utils

func Black(s string) string {
	return Parse("{{ . | black }}", s)
}

func Red(s string) string {
	return Parse("{{ . | red }}", s)
}

func Green(s string) string {
	return Parse("{{ . | green }}", s)
}

func Yellow(s string) string {
	return Parse("{{ . | yellow }}", s)
}

func Blue(s string) string {
	return Parse("{{ . | blue }}", s)
}

func Magenta(s string) string {
	return Parse("{{ . | magenta }}", s)
}

func Cyan(s string) string {
	return Parse("{{ . | cyan }}", s)
}

func White(s string) string {
	return Parse("{{ . | white }}", s)
}
