package models

// Keadaan efek keluar kasus Komite - tiket 08 Komite Claim Life. MURNI.
//
// ⛔ PERILAKU BARU SEPENUHNYA: korpus tidak punya permukaan pemantauan apa
// pun untuk kegagalan efek keluar (hanya log `DIRECTTOKASIR_LOG`,
// `MONITORING_KLAIM_LOG`). Ini syarat pengaman ADR-0015: "wajib berhasil"
// tanpa permukaan manusia hanya memindahkan kegagalan diam ke antrean.
//
// km5: KODE tinggal di sini (status outbox `antre`/`jalan`/`selesai`/
// `gagal-permanen`, milik repository/efekkeluar.go); layar menerima KATA.

import "strings"

// Kode status outbox - salinan nilai `repository.StatusEfek*`, diuji sama.
const (
	KodeEfekAntre         = "antre"
	KodeEfekJalan         = "jalan"
	KodeEfekSelesai       = "selesai"
	KodeEfekGagalPermanen = "gagal-permanen"
)

// Kata keadaan untuk layar.
const (
	KataEfekTertunda        = "tertunda"
	KataEfekTuntas          = "tuntas"
	KataEfekPerluIntervensi = "perlu intervensi"
	// KataKeputusanTersimpan - ada efek yang belum tuntas (ADR-0015).
	KataKeputusanTersimpan = "tersimpan, belum tuntas"
)

// KataStatusEfek menerjemahkan satu kode status outbox.
func KataStatusEfek(kode string) string {
	switch strings.TrimSpace(kode) {
	case KodeEfekSelesai:
		return KataEfekTuntas
	case KodeEfekGagalPermanen:
		return KataEfekPerluIntervensi
	case KodeEfekAntre, KodeEfekJalan:
		return KataEfekTertunda
	}
	// Kode asing TIDAK dianggap tuntas: yang tidak dikenal harus dilihat orang.
	return KataEfekPerluIntervensi
}

// KeadaanEfekKasus meringkas seluruh efek satu kasus.
//
// ⛔ "Tuntas" HANYA bila setiap efek selesai (AC 25 spec). Satu yang gagal
// permanen membuat kasusnya "perlu intervensi" walau yang lain selesai; kasus
// tanpa satu efek pun (belum ada keputusan) tidak diberi kata sama sekali.
func KeadaanEfekKasus(kode []string) string {
	if len(kode) == 0 {
		return ""
	}
	tertunda := false
	for _, k := range kode {
		switch KataStatusEfek(k) {
		case KataEfekPerluIntervensi:
			return KataEfekPerluIntervensi
		case KataEfekTertunda:
			tertunda = true
		}
	}
	if tertunda {
		return KataKeputusanTersimpan
	}
	return KataEfekTuntas
}
