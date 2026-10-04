package services

// Syarat berbeda tiap pemulihan limit — tiket 32.
//
// ⛔ Golongan tiket ini BARU, bukan pelestarian. Yang dilestarikan hanya
// "pemulihan disimpan sebagai daftar"; yang BARU adalah "nilainya boleh
// berbeda antar baris". `SetReinstatementPct.xml` menulis kedua persennya
// tetap "100" sebagai tetapan di dalam kode, jadi sistem lama tidak dapat
// menyatakan "pemulihan pertama gratis, kedua berbayar 100%".
//
// ⚠️ `INV-47` (baris berjumlah utuh) TIDAK berlaku di sini - `INV-49`
// mengecualikan pemulihan sebagai besaran BERULANG. Pengecualian itu disebut
// di sini supaya siapa pun yang kelak menambahkan pemeriksaan jumlah tahu ia
// sedang melanggar, bukan melengkapi.

import (
	"context"
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// ErrPemulihanTidakSah — tiket 32. Pesannya menyebut baris mana dan ruas apa.
var ErrPemulihanTidakSah = errors.New("pemulihan limit tidak sah")

var (
	nol     = apd.New(0, 0)
	seratus = apd.New(100, 0)
)

// diRentangPersen menjawab apakah d berada di 0..100 inklusif (`INV-41`).
func diRentangPersen(d *apd.Decimal) bool {
	return d.Cmp(nol) >= 0 && d.Cmp(seratus) <= 0
}

// CatatPemulihanLimit menyimpan SELURUH pemulihan sebuah layer sekaligus.
//
// ⛔ Sekaligus, bukan satu per satu, dan itu disengaja: nomor urut yang kembar
// hanya terlihat bila seluruh barisnya dilihat bersama.
//
// ⚠️ DAN LAPISAN INI SATU-SATUNYA YANG MEMERIKSANYA. `PEMULIHAN_LIMIT` tidak
// punya `UNIQUE (ID_LAYER, NOMOR_URUT_PEMULIHAN)`, dan ia TIDAK dipasang di
// sini: `SPEC-INVARIAN.md` berhenti di `INV-71` dan belum menomori kunci alami
// tabel ini, sementara `TestTabelTanpaKunciAlamiTidakDiberiDiamDiam` melarang
// menyisipkan `UNIQUE` sebelum nomornya turun - larangan yang benar, sebab
// constraint tanpa nomor invarian tidak dapat diadili siapa pun.
//
// Akibatnya dibayar dan dinyatakan: jalur mana pun yang TIDAK melewati fungsi
// ini - pemuatan data tiket 44, perbaikan manual, skrip - dapat menyisipkan
// nomor urut kembar tanpa ditolak. Itu tagihan ke `SPEC-INVARIAN.md`, bukan
// alasan memasang constraint diam-diam.
func (l *Layanan) CatatPemulihanLimit(ctx context.Context, p inti.Pelaku, idLayer int64, baris []models.PemulihanLimit) error {
	if err := inti.WajibIdentitas(p); err != nil {
		return err
	}
	if idLayer <= 0 {
		return fmt.Errorf("%w: pengenal layer %d bukan pengenal yang sah", ErrMasukanTidakSah, idLayer)
	}
	if len(baris) == 0 {
		return fmt.Errorf("%w: nol baris pemulihan dikirim; menghapus seluruh pemulihan bukan jalur ini", ErrPemulihanTidakSah)
	}

	var salah []string
	terlihat := map[int64]int{}
	for i, b := range baris {
		ke := i + 1
		switch {
		case b.NomorUrut < 1:
			salah = append(salah, fmt.Sprintf("baris ke-%d: nomor urut pemulihan %d, harus 1 atau lebih", ke, b.NomorUrut))
		default:
			if sebelumnya, ada := terlihat[b.NomorUrut]; ada {
				salah = append(salah, fmt.Sprintf("baris ke-%d: nomor urut pemulihan %d sudah dipakai baris ke-%d", ke, b.NomorUrut, sebelumnya))
			} else {
				terlihat[b.NomorUrut] = ke
			}
		}

		// Porsi limit yang dipulihkan - WAJIB, dan kosong bukan nol.
		switch {
		case b.PersenPemulihan == nil:
			salah = append(salah, fmt.Sprintf("baris ke-%d: persen pemulihan kosong; kolomnya wajib isi", ke))
		case !diRentangPersen(b.PersenPemulihan):
			salah = append(salah, fmt.Sprintf("baris ke-%d: persen pemulihan %s di luar rentang 0-100 (INV-41)", ke, b.PersenPemulihan.Text('f')))
		}

		// Tarif premi pemulihan - BOLEH kosong, dan kosong berarti belum
		// dinyatakan. Diisi nol hanya bila memang nol.
		if b.PersenTambahan != nil && !diRentangPersen(b.PersenTambahan) {
			salah = append(salah, fmt.Sprintf("baris ke-%d: persen tambahan %s di luar rentang 0-100 (INV-41)", ke, b.PersenTambahan.Text('f')))
		}
	}
	// ⚠️ TIDAK diurutkan. Urutannya urutan baris masukan, dan itu yang
	// berguna bagi yang memperbaikinya; mengurutkan teksnya menaruh
	// "baris ke-10" sebelum "baris ke-2".
	if len(salah) > 0 {
		return fmt.Errorf("%w: %v", ErrPemulihanTidakSah, salah)
	}

	if err := l.gudang.CatatPemulihanLimit(ctx, idLayer, baris); err != nil {
		return err
	}
	return nil
}
