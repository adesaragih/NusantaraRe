package menu

// `GET /api/menu` - milik aplikasi, bukan satu modul: dipasang `cmd/api`
// walau modul mana pun nonaktif.

import (
	"log"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
)

// Rute menjawab `GET /api/menu`.
//
// `pembaca` nil = proses tanpa Oracle. `modulAktif` adalah nama modul yang
// dipasang proses ini (MODUL_AKTIF): butir modul lain tidak dikirim.
//
// Saringan per akun (Kelola User, 01-10-2026): permintaan bersesi login
// menerima HANYA menu akunnya (`SaringMenuUntukAkun`). Tanpa sesi,
// `stubPelaku` (AUTH_STUB) menerima seluruh menu tabel tanpa menu aplikasi -
// Kelola User menuntut login sungguhan - dan tanpa stub jawabannya 401.
//
// ⛔ Kegagalan dijawab galat TERANG, tidak pernah menu kosong: sidebar yang
// kosong diam-diam terbaca "aplikasi tanpa menu". Rincian galat driver
// tinggal di log server, tidak di badan jawaban - pesannya dapat memuat
// nilai kolom (pola `modul/claimlife/backend/handlers`, `inti/backend/layanan`).
func Rute(pembaca PembacaMenu, modulAktif []string, stubPelaku bool) http.HandlerFunc {
	aktif := append([]string(nil), modulAktif...)
	return func(w http.ResponseWriter, r *http.Request) {
		if pembaca == nil {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi; menu dibaca dari M_NAV_MENU")
			return
		}
		kode, sesi := inti.AksesMenuDari(r.Context())
		if !sesi && !stubPelaku {
			galat.Tulis(w, http.StatusUnauthorized, "belum login atau sesi sudah berakhir")
			return
		}
		baris, err := pembaca.Baca(r.Context())
		if err != nil {
			if strings.Contains(err.Error(), "ORA-00942") {
				galat.Tulis(w, http.StatusServiceUnavailable,
					"tabel M_NAV_MENU belum ada - migrasi 900 belum dijalankan (-migrate, oleh work owner)")
				return
			}
			log.Printf("menu: membaca M_NAV_MENU: %v", err)
			galat.Tulis(w, http.StatusInternalServerError, "menu: membaca M_NAV_MENU gagal; rinciannya di log server")
			return
		}
		m := Susun(baris, aktif)
		if sesi {
			m = SaringMenuUntukAkun(m, kode)
		}
		galat.TulisJSON(w, m)
	}
}
