package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"strings"
)

var GmailIcon = &fyne.StaticResource{
	StaticName: "gmail_icon.svg",
	StaticContent: []byte(`
		<?xml version="1.0" standalone="no"?>
		<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 20010904//EN"
		 "http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd">
		<svg version="1.0" xmlns="http://www.w3.org/2000/svg"
		 width="512.000000pt" height="512.000000pt" viewBox="0 0 512.000000 512.000000"
		 preserveAspectRatio="xMidYMid meet">
		
		<g transform="translate(0.000000,512.000000) scale(0.100000,-0.100000)"
		fill="#ffffff" stroke="none">
		<path d="M413 4465 c-144 -31 -287 -140 -351 -266 -65 -130 -63 -50 -60 -1747
		l3 -1537 31 -65 c53 -110 140 -178 256 -200 28 -6 235 -10 459 -10 l409 0 2
		991 3 990 695 -521 c382 -286 700 -518 706 -516 7 3 321 236 698 518 l685 513
		1 754 c0 601 -3 751 -12 744 -7 -6 -320 -240 -696 -522 l-682 -511 -878 658
		c-482 363 -901 671 -931 686 -102 52 -221 66 -338 41z"/>
		<path d="M4538 3047 l-577 -432 -1 -987 0 -988 398 0 c229 0 424 5 462 11 122
		19 226 102 272 217 l23 57 3 1110 c1 611 1 1185 0 1277 l-3 167 -577 -432z"/>
		</g>
		</svg>`),
}

var FacebookIcon = &fyne.StaticResource{
	StaticName: "facebook_icon.svg",
	StaticContent: []byte(`
		<?xml version="1.0" standalone="no"?>
		<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 20010904//EN"
		"http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd">
		<svg version="1.0" xmlns="http://www.w3.org/2000/svg"
		width="512.000000pt" height="512.000000pt" viewBox="0 0 512.000000 512.000000"
		preserveAspectRatio="xMidYMid meet">
		
		<g transform="translate(0.000000,512.000000) scale(0.100000,-0.100000)"
		fill="#ffffff" stroke="none">
		<path d="M555 5110 c-211 -29 -391 -160 -485 -355 -70 -146 -65 15 -65 -2195
		0 -2210 -5 -2049 65 -2195 64 -133 156 -226 285 -291 146 -72 69 -68 1193 -71
		l1012 -4 0 881 0 880 -320 0 -320 0 0 400 0 400 319 0 320 0 3 408 c4 388 5
		411 27 482 101 337 344 580 681 681 71 22 94 23 483 27 l407 3 0 -399 0 -400
		-262 -3 c-241 -4 -266 -6 -297 -24 -19 -11 -43 -33 -55 -48 -20 -27 -21 -41
		-24 -378 l-3 -349 401 0 c379 0 402 -1 396 -17 -4 -10 -75 -190 -159 -400
		l-152 -383 -242 0 -243 0 0 -881 0 -880 533 4 c588 4 574 3 712 71 129 65 221
		158 285 291 70 146 65 -15 65 2195 0 2210 5 2049 -65 2195 -64 133 -156 226
		-285 291 -151 75 46 68 -2165 70 -1092 1 -2012 -2 -2045 -6z"/>
		</g>
		</svg>`),
}

var TwitterIcon = &fyne.StaticResource{
	StaticName: "twitter_icon.svg",
	StaticContent: []byte(`
		<?xml version="1.0" standalone="no"?>
		<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 20010904//EN"
		 "http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd">
		<svg version="1.0" xmlns="http://www.w3.org/2000/svg"
		 width="512.000000pt" height="512.000000pt" viewBox="0 0 512.000000 512.000000"
		 preserveAspectRatio="xMidYMid meet">
		
		<g transform="translate(0.000000,512.000000) scale(0.100000,-0.100000)"
		fill="#ffffff" stroke="none">
		<path d="M3367 4624 c-408 -74 -731 -372 -838 -776 -19 -69 -23 -110 -23 -248
		-1 -91 2 -184 7 -208 l7 -44 -67 7 c-441 40 -797 141 -1152 326 -328 171 -669
		440 -878 691 -32 38 -62 67 -67 65 -6 -2 -27 -42 -48 -90 -144 -325 -117 -700
		70 -992 55 -85 145 -186 218 -244 l75 -61 -41 0 c-88 0 -271 50 -370 100 -21
		11 -43 20 -47 20 -14 0 4 -183 27 -271 91 -345 347 -618 688 -732 56 -19 96
		-37 88 -40 -45 -17 -193 -29 -308 -25 l-130 6 7 -26 c11 -44 85 -184 131 -248
		176 -248 465 -418 747 -441 42 -3 77 -9 77 -14 0 -16 -224 -155 -361 -223
		-312 -155 -621 -222 -984 -214 l-185 4 52 -33 c77 -48 309 -165 413 -208 561
		-233 1216 -286 1834 -149 909 202 1658 834 2036 1719 161 378 242 750 251
		1150 l3 180 73 58 c92 75 196 174 275 262 67 76 165 207 159 212 -2 2 -32 -8
		-67 -22 -133 -54 -440 -132 -486 -123 -11 2 -3 12 26 31 76 50 211 189 268
		275 50 74 123 221 123 247 0 6 -61 -20 -136 -57 -136 -66 -313 -130 -451 -163
		l-73 -17 -61 56 c-234 217 -571 316 -882 260z"/>
		</g>
		</svg>`),
}

var mapStringColors = map[string]string{
	"Dark":  "#87CEFA",
	"Light": "#87CEFA",
	"Blue":  "#267AB6",
	"Pink":  "#E86EC1",
}

// Función que devuelve el ícono adecuado en función del tema
func GetIconForPassword(passwordDetails map[string]string, t string) fyne.Resource {
	label := passwordDetails["label"]
	color, exists := mapStringColors[t]
	if !exists {
		color = "#FFFFFF"
	}

	var svgContent []byte
	switch {
	case containsIgnoreCase(label, "facebook"):
		svgContent = setSVGFill(FacebookIcon.StaticContent, color)
	case containsIgnoreCase(label, "twitter"):
		svgContent = setSVGFill(TwitterIcon.StaticContent, color)
	case containsIgnoreCase(label, "gmail"):
		svgContent = setSVGFill(GmailIcon.StaticContent, color)
	default:
		return theme.FileIcon()
	}

	return &fyne.StaticResource{
		StaticName:    label + "_icon.svg",
		StaticContent: svgContent,
	}
}

func setSVGFill(svg []byte, color string) []byte {
	svgString := string(svg)
	updatedSVG := strings.ReplaceAll(svgString, `fill="#ffffff"`, fmt.Sprintf(`fill="%s"`, color))
	return []byte(updatedSVG)
}

func containsIgnoreCase(label string, keyword string) bool {
	return strings.Contains(strings.ToLower(label), strings.ToLower(keyword))
}
