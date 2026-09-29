package inti

// Pelaku permintaan HTTP - penunda sampai tiket 07 / IAM.
//
// ⛔ BERKAS INI BUKAN AUTENTIKASI, dan tidak boleh disalahartikan begitu.
// Header HTTP dapat ditulis siapa saja yang dapat mengirim permintaan, jadi
// `X-Pelaku` tidak membuktikan apa pun tentang siapa pengirimnya.
//
// `[keputusan work owner 26-09-2026, butir ab]`. Kenapa ia tetap ada: tiket 05
// menuntut peran `ReasLifeAdmin` sekarang, sedangkan sumber peran yang
// sebenarnya adalah satu tabel (ADR-U-0030) yang belum ada dan tiket 07 yang
// menegakkannya. Tanpa penunda ini, jalur wewenang tidak dapat dibangun
// maupun diuji sampai tiket 07 selesai.
//
// Tiga pagar, dan ketiganya perlu:
//  1. mati secara bawaan - AUTH_STUB harus disetel sadar
//  2. ditolak saat memuat konfigurasi bila IS_PEGA_PROD=true
//  3. seluruh uji wewenang berjalan di seam services dengan Pelaku langsung,
//     sehingga aturannya tidak bergantung pada jalur ini sama sekali
//
// Dibaca sesudah: handlers.go.

import (
	"net/http"
	"strings"
)

// pelakuDari membaca pelaku permintaan.
//
// Tanpa stub, ia mengembalikan pelaku KOSONG - bukan pelaku istimewa. Jalur
// yang memerlukan identitas karena itu DITOLAK, dan itu memang yang benar
// sampai sumber peran yang sebenarnya ada.
//
// ⛔ Keadaan stub dibawa sebagai ARGUMEN, bukan variabel paket. Variabel paket
// membuat dua Router dalam satu proses - yang biasa terjadi di test - berbagi
// satu saklar: Router(svc, true) akan menyalakan stub bagi server yang dibuat
// Router(svc, false) dan masih hidup. Test menjadi bergantung urutan, dan
// `go test -race` menandainya sebagai tulis-baca serentak.
func PelakuDari(r *http.Request, stubAktif bool) Pelaku {
	if !stubAktif {
		return Pelaku{}
	}
	var peran []string
	for _, p := range strings.Split(r.Header.Get("X-Peran"), ",") {
		if p = strings.TrimSpace(p); p != "" {
			peran = append(peran, p)
		}
	}
	return Pelaku{
		AkunID: strings.TrimSpace(r.Header.Get("X-Pelaku")),
		Peran:  peran,
	}
}
