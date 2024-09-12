package locale

const (
	ShortcodeEnglish = "eng"
	ShortcodeRussian = "rus"
	ShortcodeUzbek   = "uzb"
)

type Locale struct {
	name      string
	shortcode string
}

func (l Locale) Name() string {
	return l.name
}

func (l Locale) Shortcode() string {
	return l.shortcode
}

var (
	English = Locale{"English", ShortcodeEnglish}
	Russian = Locale{"Russian", ShortcodeRussian}
	Uzbek   = Locale{"Uzbek", ShortcodeUzbek}
)

func LangToLocale(lang string) Locale {
	switch lang {
	case ShortcodeEnglish:
		return English
	case ShortcodeRussian:
		return Russian
	case ShortcodeUzbek:
		return Uzbek
	default:
		return English
	}
}
