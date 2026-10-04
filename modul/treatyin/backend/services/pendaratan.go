package services

// Pemuatan tabel pendaratan - urutan, rekonsiliasi, dan keputusan
// mengikat-atau-membatalkan.
//
// ⛔ YANG DIPUTUSKAN DI SINI: muat dulu, COCOKKAN, baru ikat. Pemuat yang
// mengikat lebih dulu lalu mencocokkan akan menemukan kerusakan sesudah
// kerusakan itu terlihat orang lain - dan 26.536 baris yang salah jauh
// lebih mahal dibereskan daripada satu transaksi yang dibatalkan.
//
// ⚠️ Lapisan ini TIDAK membuka maupun menutup transaksi untuk satu kontrak
// dengan caranya sendiri; pemanggilnya yang memegang transaksi, sebab
// pemanggil itulah yang tahu sepuluh kontrak dulu atau seluruhnya.

import (
	"context"
	"fmt"
	"sort"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/repository"
)

// HasilMuat adalah satu kontrak yang sudah didaratkan, beserta buktinya.
type HasilMuat struct {
	MasterID string
	// Cacah elemen di DOKUMEN, per tabel.
	DiDokumen map[string]int
	// Cacah baris di TABEL sesudah pemuatan, per tabel.
	DiTabel map[string]int
	// Kunci JSON yang belum punya kolom, per tabel.
	TakTerpetakan map[string][]string
	// Selisih yang ditemukan; kosong berarti cocok.
	Selisih []string
}

// Cocok menjawab apakah kontrak ini boleh diikat.
func (h HasilMuat) Cocok() bool { return len(h.Selisih) == 0 }

// MuatSatuKontrak mendaratkan satu kontrak di dalam transaksi yang diberikan,
// lalu MENCOCOKKAN hasilnya tanpa keluar dari transaksi itu.
//
// ⛔ Tidak mengikat dan tidak membatalkan. Keputusan itu milik pemanggil,
// yang membacanya dari `HasilMuat.Cocok()`.
func MuatSatuKontrak(ctx context.Context, g *repository.Gudang, tx *db.Tx, masterID string) (HasilMuat, error) {
	h := HasilMuat{MasterID: masterID}

	teks, err := g.BacaDokumenMentah(ctx, masterID)
	if err != nil {
		return h, err
	}
	doc, err := repository.UraiDokumen(teks)
	if err != nil {
		return h, fmt.Errorf("services: kontrak %s: %w", masterID, err)
	}

	h.DiDokumen = repository.CacahLarik(doc)
	h.TakTerpetakan = repository.KunciTakTerpetakan(doc)

	if _, err := g.MuatKontrak(ctx, tx, masterID, doc); err != nil {
		return h, err
	}

	// ⛔ Dibaca DI DALAM transaksi: yang hendak dibuktikan adalah apa yang
	// baru saja ditulis, dan dari luar transaksi itu belum ada.
	h.DiTabel, err = g.CacahBarisKontrak(ctx, tx, masterID)
	if err != nil {
		return h, err
	}

	h.Selisih = selisihCacah(h.DiDokumen, h.DiTabel)
	return h, nil
}

// selisihCacah mengadu dua peta cacah dan menamai setiap ketidakcocokan.
//
// ⚠️ Menyebut TABEL, dokumen, dan tabel sekaligus. Pesan "tidak cocok" tanpa
// angkanya memaksa orang berikutnya menjalankan ulang seluruh pemuatan hanya
// untuk mengetahui apa yang tidak cocok.
func selisihCacah(dokumen, tabel map[string]int) []string {
	var out []string
	nama := make([]string, 0, len(dokumen))
	for t := range dokumen {
		nama = append(nama, t)
	}
	sort.Strings(nama)
	for _, t := range nama {
		if dokumen[t] != tabel[t] {
			out = append(out, fmt.Sprintf("%s: dokumen %d, tabel %d", t, dokumen[t], tabel[t]))
		}
	}
	// Tabel yang ADA isinya tetapi tidak ada di dokumen - sisa muatan lama
	// yang pengosongan lewatkan. Ini yang membuktikan idempotensi.
	for t, n := range tabel {
		if _, ada := dokumen[t]; !ada && n != 0 {
			out = append(out, fmt.Sprintf("%s: %d baris tersisa tanpa padanan di dokumen", t, n))
		}
	}
	return out
}

// RingkasTakTerpetakan menggabungkan laporan kunci tanpa kolom dari banyak
// kontrak menjadi satu daftar - supaya 1.854 laporan tidak menjadi 1.854
// baris yang sama.
func RingkasTakTerpetakan(hasil []HasilMuat) map[string][]string {
	kumpul := map[string]map[string]bool{}
	for _, h := range hasil {
		for tabel, kunci := range h.TakTerpetakan {
			if kumpul[tabel] == nil {
				kumpul[tabel] = map[string]bool{}
			}
			for _, k := range kunci {
				kumpul[tabel][k] = true
			}
		}
	}
	out := map[string][]string{}
	for tabel, set := range kumpul {
		daftar := make([]string, 0, len(set))
		for k := range set {
			daftar = append(daftar, k)
		}
		sort.Strings(daftar)
		out[tabel] = daftar
	}
	return out
}
