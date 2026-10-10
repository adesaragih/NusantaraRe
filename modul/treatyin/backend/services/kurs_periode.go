package services

import (
	"fmt"
	"strings"
	"time"

	"nusantarare/modul/treatyin/backend/models"
)

// KursDalamPeriode menyaring grid Rate of Exchange satu kontrak: hanya baris
// yang masa berlakunya BERADA DI DALAM periode kontrak — `Valid From` >=
// Commencement dan `Valid Until` <= Expiry — yang dikembalikan.
//
// ⛔ LAHIR DARI CACAT NYATA, 10 Oktober 2026. Kontrak Non-Prop 1002309
// (01/01/2026–31/12/2026) diisi DUA baris kurs (IDR, USD); sesudah Save
// grid menampilkan SEMBILAN. Sebabnya: `TREATYEXCHANGEYEARLY` berkunci
// TAHUN dan dipakai bersama (`BacaKursTahunan`), jadi setiap kontrak 2026
// yang pernah di-Save menambahkan barisnya ke grid kontrak lain — periode
// 01/04, 01/07, dan 08/10 ikut tampil.
//
// ⭐ Keputusan pemilik proses 10 Oktober 2026: "Saring periode". Tabelnya
// TETAP satu-satunya tempat kurs (nol tabel, nol migrasi); yang berubah
// hanya baris mana yang grid kontrak ini tampilkan — dan, sejak hari yang
// sama, baris mana yang Delete grid boleh hapus (`kursBolehDihapus`).
//
// ⚠️ Akibat yang diterima: dua kontrak dengan periode yang PERSIS sama
// tetap saling melihat barisnya — tabelnya tidak punya kolom kontrak.
//
// ⚠️ Tanggal baris dibaca dari `…Asli` (tanggal WIB, yang grid tampilkan),
// jadi `KursDalamPeriode` dipanggil SESUDAH medan itu diisi. Batas kontrak
// yang kosong atau bukan tanggal tidak menyaring sisi itu; baris yang
// tanggalnya kosong atau bukan tanggal tidak lolos sisi yang disaring.
func KursDalamPeriode(baris []models.BarisKursWarisan, mulai, berakhir string) []models.BarisKursWarisan {
	awal, adaAwal := tanggalSah(TanggalWIB(mulai))
	akhir, adaAkhir := tanggalSah(TanggalWIB(berakhir))
	if !adaAwal && !adaAkhir {
		return baris
	}
	out := make([]models.BarisKursWarisan, 0, len(baris))
	for _, b := range baris {
		if adaAwal {
			if d, ok := tanggalSah(b.BerlakuDariAsli); !ok || d < awal {
				continue
			}
		}
		if adaAkhir {
			if s, ok := tanggalSah(b.BerlakuSampaiAsli); !ok || s > akhir {
				continue
			}
		}
		out = append(out, b)
	}
	return out
}

// KursGridKontrak - baris yang grid Rate of Exchange kontrak tampilkan.
//
// ⭐ Kontrak yang MASIH DAPAT DI-SAVE melihat periodenya saja, persis: Add,
// Delete, lalu Save menghasilkan jumlah baris yang sama dengan yang tampil
// sebelum Save (permintaan pemakai 10 Oktober 2026).
//
// ⛔ Kontrak `Resolve Complete` yang NOL barisnya di dalam periode melihat
// SELURUH TAHUN, seperti sebelum saringan ini ada. Dicacah 10 Oktober 2026
// di DEV_NUSARE2: 1.302 kontrak tuntas tidak punya satu pun baris di dalam
// periodenya — kurs warisan berperiode 01/07–30/06, kontraknya 01/01–31/12
// — dan tanpa cadangan ini grid mereka kosong. Kontrak tuntas tidak dapat
// di-Save, jadi cadangan ini tidak pernah bertemu Add/Delete. Yang masih
// dapat disunting dan nol baris di periodenya: 7 kontrak.
func KursGridKontrak(baris []models.BarisKursWarisan, mulai, berakhir, status string) []models.BarisKursWarisan {
	dalam := KursDalamPeriode(baris, mulai, berakhir)
	if len(dalam) == 0 && status == models.StatusTuntas {
		return baris
	}
	return dalam
}

// kursBolehDihapus - dari `hapus`, hanya pengenal yang ADA di `tahunan`
// (seluruh baris tahun kontrak, `…Asli` sudah terisi) dan berada di dalam
// periode `mulai`–`berakhir`.
//
// ⛔ PENGAMAN, bukan hiasan: tabelnya dipakai bersama. Baris di luar
// periode kontrak — kurs kontrak lain, atau cadangan seluruh tahun — tidak
// pernah tampil sebagai milik kontrak ini, jadi tidak boleh terhapus
// darinya. Yang ditolak dilewati tanpa galat: sesudah Save grid disaring
// ulang dengan periode yang sama, dan baris itu memang tidak tampil.
func kursBolehDihapus(hapus []string, tahunan []models.BarisKursWarisan, mulai, berakhir string) []string {
	if len(hapus) == 0 {
		return nil
	}
	if _, ok := tanggalSah(TanggalWIB(mulai)); !ok {
		return nil
	}
	if _, ok := tanggalSah(TanggalWIB(berakhir)); !ok {
		return nil
	}
	dalam := map[string]bool{}
	for _, b := range KursDalamPeriode(tahunan, mulai, berakhir) {
		dalam[b.ID] = true
	}
	out := []string{}
	for _, id := range hapus {
		if id = strings.TrimSpace(id); id != "" && dalam[id] {
			out = append(out, id)
		}
	}
	return out
}

// tanggalSah - `YYYYMMDD` yang benar-benar tanggal; selain itu `false`.
// Bentuk delapan digit dapat dibandingkan sebagai teks.
func tanggalSah(s string) (string, bool) {
	t := strings.TrimSpace(s)
	if len(t) != 8 {
		return "", false
	}
	if _, err := time.Parse("20060102", t); err != nil {
		return "", false
	}
	return t, true
}

// periksaKursDalamPeriode - setiap baris kurs yang disimpan harus punya Valid
// From dan Valid Until yang sah dan berada di dalam periode kontrak, sebab
// grid hanya menampilkan baris seperti itu (`KursDalamPeriode`). Baris tanpa
// mata uang (baris kosong yang baru ditambahkan) dilewati, seperti
// `tulisKurs`. Periode kontrak yang kosong atau tak terbaca tidak diperiksa.
func periksaKursDalamPeriode(baris []models.BarisKursWarisan, mulai, berakhir string) error {
	awal, adaAwal := tanggalSah(TanggalWIB(mulai))
	akhir, adaAkhir := tanggalSah(TanggalWIB(berakhir))
	if !adaAwal || !adaAkhir {
		return nil
	}
	tampil := func(t string) string { return t[6:8] + "/" + t[4:6] + "/" + t[0:4] }
	for _, b := range baris {
		if strings.TrimSpace(b.MataUangID) == "" && strings.TrimSpace(b.MataUang) == "" {
			continue
		}
		nama := strings.TrimSpace(b.MataUang)
		if nama == "" {
			nama = b.MataUangID
		}
		dari, okDari := tanggalSah(b.BerlakuDariAsli)
		sampai, okSampai := tanggalSah(b.BerlakuSampaiAsli)
		if !okDari || !okSampai {
			return ditolak(fmt.Sprintf("Rate of Exchange %s: Valid From dan Valid Until wajib diisi.", nama))
		}
		if dari < awal || sampai > akhir {
			return ditolak(fmt.Sprintf(
				"Rate of Exchange %s: Valid From–Valid Until (%s–%s) harus berada di dalam periode kontrak (%s–%s).",
				nama, tampil(dari), tampil(sampai), tampil(awal), tampil(akhir)))
		}
	}
	return nil
}
