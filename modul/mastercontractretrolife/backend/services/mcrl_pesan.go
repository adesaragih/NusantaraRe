package services

// Kalimat galat untuk layar (K8).
//
// ⛔ Galat Go sendiri WAJIB sampai ke layar berkata-kata - tidak ada yang
// ditelan, tidak ada HTML. Kalimat bawaan modul berbahasa Inggris, seperti
// Treaty Contract Out [keputusan work owner 30-09-2026 untuk TCO]; pesan yang
// ada di korpus dibawa VERBATIM (`PesanKosongTahun`, `PesanKosongSemua`, dst).

import (
	"errors"
	"regexp"
	"strings"
)

// Pesan VERBATIM korpus.
const (
	// PesanKosongTahun - `SaveTreatyYearLife_Act.xml` b313.
	PesanKosongTahun = "Value cannot be empty."
	// PesanKosongSemua - `SaveTreatyLimit_Act.xml` b313, `SaveSecurityLife_Act.xml` b299,
	// `SaveSecurityReinsurerLife_Act.xml` b284, `SaveBusinessLife_Act.xml` b293.
	PesanKosongSemua = "All value cannot be empty."
	// PesanHapusBerhasil - `DeleteTreatyLimit_Act.xml` b1057, `DeleteSecurityLife_Act.xml` b942,
	// `DeleteSecurityReinsurerLife_Act.xml` b881.
	PesanHapusBerhasil = "Data Berhasil di Hapus"
	// PesanSalinSemua - `SaveBusinessToAllLife_Act.xml` b1258.
	PesanSalinSemua = "Copied to all reins types."
)

// PesanHapusBusiness - `DeleteRowBusiness.xml` b665:
// "Data Dengan ID" + " " + ID + " " + "Berhasil di Hapus".
func PesanHapusBusiness(id string) string {
	return "Data Dengan ID" + " " + id + " " + "Berhasil di Hapus"
}

// pesanLayar - galat yang punya kalimat layar sendiri, berbeda dari teks
// lognya (mis. `repository.GalatMaster`: sebab Oracle hanya di log).
type pesanLayar interface{ PesanLayar() string }

var polaMarkup = regexp.MustCompile(`<[^>]*>`)

// Pesan - kalimat galat untuk layar, di SATU tempat (tiket 10): kalimat layar
// galat bila ada, awalan lapisan dibuang, dan markup dibuang - nilai masukan
// atau teks luar yang memuat tag tidak pernah kembali sebagai tag.
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
