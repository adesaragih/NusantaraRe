package services

// Kalimat galat untuk layar.
//
// ⛔ Galat Go sendiri WAJIB sampai ke layar berkata-kata - tidak ada yang
// ditelan, tidak ada HTML. Kalimat bawaan modul berbahasa Inggris (label
// korpus berbahasa Inggris); pesan yang ada di korpus dibawa VERBATIM
// (`PesanProductNameKosong`, `PesanPlanKosong`, `PesanSalin`, ...).

import (
	"errors"
	"regexp"
	"strings"
)

// pesanLayar - galat yang punya kalimat layar sendiri, berbeda dari teks
// lognya (sebab Oracle hanya di log).
type pesanLayar interface{ PesanLayar() string }

var polaMarkup = regexp.MustCompile(`<[^>]*>`)

// Pesan - kalimat galat untuk layar, di SATU tempat: kalimat layar galat bila
// ada, awalan lapisan dibuang, markup dibuang.
func Pesan(err error) string {
	if err == nil {
		return ""
	}
	var l pesanLayar
	s := err.Error()
	if errors.As(err, &l) {
		s = l.PesanLayar()
	}
	for _, a := range []string{"services: ", "repository: ", "models: "} {
		s = strings.ReplaceAll(s, a, "")
	}
	s = polaMarkup.ReplaceAllString(s, "")
	return strings.TrimSpace(strings.NewReplacer("<", "", ">", "").Replace(s))
}
