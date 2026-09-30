package menu

// `GET /api/menu` - milik aplikasi, bukan satu modul: dipasang `cmd/api`
// walau modul mana pun nonaktif.

import (
	"net/http"
	"strings"

	"nusantarare/inti"
	"nusantarare/inti/galat"
)

// Rute menjawab `GET /api/menu`.
//
// `pembaca` nil = proses tanpa Oracle. `modulAktif` adalah nama modul yang
// dipasang proses ini (MODUL_AKTIF): butir modul lain tidak dikirim.
// `stubPelaku` = AUTH_STUB, dibaca hanya untuk diteruskan ke
// `SaringMenuUntukPelaku`.
//
// ⛔ Kegagalan dijawab galat TERANG (ADR-0015), tidak pernah menu kosong:
// sidebar yang kosong diam-diam terbaca "aplikasi tanpa menu".
func Rute(pembaca PembacaMenu, modulAktif []string, stubPelaku bool) http.HandlerFunc {
	aktif := append([]string(nil), modulAktif...)
	return func(w http.ResponseWriter, r *http.Request) {
		if pembaca == nil {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi; menu dibaca dari M_NAV_MENU")
			return
		}
		baris, err := pembaca.Baca(r.Context())
		if err != nil {
			if strings.Contains(err.Error(), "ORA-00942") {
				galat.Tulis(w, http.StatusServiceUnavailable,
					"tabel M_NAV_MENU belum ada - migrasi 900 belum dijalankan (-migrate, oleh work owner)")
				return
			}
			galat.Tulis(w, http.StatusInternalServerError, "menu: membaca M_NAV_MENU gagal: "+err.Error())
			return
		}
		galat.TulisJSON(w, SaringMenuUntukPelaku(inti.PelakuDari(r, stubPelaku), Susun(baris, aktif)))
	}
}
