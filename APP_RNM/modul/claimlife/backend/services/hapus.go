package services

// Menghapus klaim beserta seluruh isinya - tiket 15.
//
// Untuk apa berkas ini: penghapusan yang DIDAHULUI angka. Pengguna berhak tahu
// berapa baris tiap jenis yang akan ikut terhapus sebelum ia menekan Ya, dan
// angka itu wajib sama persis dengan yang benar-benar terhapus.
//
// Dibaca sesudah: tolak.go.
//
// Istilah:
//   - dampak : cacah baris per jenis yang akan ikut terhapus.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// PeranHapusKlaim adalah peran yang boleh menghapus klaim.
//
// ⚠️ Korpus TIDAK memuat satu pun rule yang menghapus kasus klaim -
// pencarian `Obj-Delete`, `pxDelete`, dan nama berkas `*Delete*`/`*Hapus*` di
// 136 berkas `Claim Life` hanya menemukan penghapusan DOKUMEN dan penomoran
// ulang peserta. Jadi peran ini BUKAN bacaan XML melainkan turunan dari tiket
// 15 sendiri ("Sebagai ReasLifeAdmin, saya dapat menghapus"). Tiket 07 sudah
// lewat dan tidak menemukan gerbang XML untuk penghapusan; nilainya tetap,
// sumbernya kini dinyatakan terus terang.
const PeranHapusKlaim = inti.PeranAdmin

var (
	// ErrKlaimSudahDiKomite - klaim yang sudah diserahkan tidak dihapus.
	ErrKlaimSudahDiKomite = errors.New(
		"services: klaim sudah diserahkan ke Komite dan tidak dapat dihapus")
	// ErrHapusFisikDilarang - jalur pengguna tidak boleh menghapus fisik.
	ErrHapusFisikDilarang = errors.New(
		"services: penghapusan klaim di jalur pengguna adalah PENANDA, bukan hapus fisik")
)

// Penghapusan adalah layanan penghapusan klaim.
//
// ⚠️ TANPA medan jejak. Ronde pertama memasangnya beserta `DenganJejak`, dan
// keduanya tidak pernah dibaca sekali pun - perkabelan yang menjanjikan
// perekaman yang tidak dilakukannya. Jejak penghapusan lahir bersama bentuk
// logisnya (ADR-U-0031) dan tabelnya (butir am), bukan sebelum itu.
type Penghapusan struct {
	svc *Service
}

// Penghapusan menyusun layanan itu.
func (s *Service) Penghapusan() *Penghapusan { return &Penghapusan{svc: s} }

// Dampak menghitung baris yang akan ikut terhapus - TANPA menghapus apa pun.
//
// ⛔ Fungsi ini hanya MEMBACA. Itulah yang membuat "Batal" benar-benar
// membatalkan: jalur pratinjau tidak punya satu pun tulisan untuk dibatalkan.
func (h *Penghapusan) Dampak(ctx context.Context, pelaku inti.Pelaku, klaimID string) (
	models.DampakHapus, error) {
	var d models.DampakHapus
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return d, err
	}
	if err := inti.WajibPeran(pelaku, PeranHapusKlaim); err != nil {
		return d, err
	}
	if strings.TrimSpace(klaimID) == "" {
		return d, fmt.Errorf("%w: pengenal klaim wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if !h.svc.PunyaDatabase() {
		return d, db.ErrTanpaOracle
	}
	if _, err := h.pastikanAda(ctx, klaimID); err != nil {
		return d, err
	}
	// ⛔ CASE_ID DIBACA, tidak diandaikan sama dengan pengenal klaim. Butir ae1
	// memang mengisinya begitu untuk klaim yang sistem ini buat sendiri, tetapi
	// itu keputusan pengisian - bukan jaminan bentuk - dan klaim yang kelak
	// dimigrasikan membawa CASE_ID warisannya sendiri. Mengandaikannya berarti
	// mencacah baris milik klaim lain.
	caseID, err := repository.NewKlaimLife(h.svc.DB()).CaseIDKlaim(ctx, klaimID)
	if err != nil {
		return d, err
	}
	return repository.NewPohonKlaim(h.svc.DB()).Dampak(ctx, klaimID, caseID)
}

// Hapus TIDAK menghapus - dan itu keputusan yang sudah diambil, bukan celah.
//
// ⛔ ADR-U-0031 `[keputusan work owner, 23-09-2026]`: *"Penghapusan adalah
// penanda dan nilai balik, bukan hapus fisik … Nol perintah hapus fisik pada
// jalur pengguna di lapisan mana pun."* Teks tiket 15 ditulis 18-09-2026 dan
// berbicara tentang kaskade fisik lima tingkat; ADR itu LIMA HARI lebih muda
// dan memutuskan sebaliknya.
//
// Yang menahan pelaksanaannya: **kolom penandanya belum ada**. Menambahkannya
// adalah langkah migrasi baru, dan langkah migrasi baru hanya lahir dari paket
// keputusan yang masih `[USULAN]`. Menebak nama kolom dan artinya berarti
// mengarang bentuk penyimpanan.
//
// Karena itu fungsi ini GAGAL TERANG, sementara `Dampak` tetap bekerja penuh:
// peringatan berisi angka - bagian yang benar-benar dilihat pengguna - dapat
// dibangun dan diuji sekarang, dan yang menunggu keputusan hanya tulisannya.
//
// ⚠️ `repository.PohonKlaim.Hapus` yang menghapus fisik TETAP ADA dan tetap
// hanya dipanggil test: ia membuktikan kaskade `ON DELETE CASCADE` tiket 14,
// dan itu pembuktian bentuk basis data - bukan jalur pengguna.
//
// ⚠️ TANPA argumen waktu. Ronde pertama menerimanya dan tidak pernah
// membacanya - tanda tangan yang menjanjikan penstempelan yang tidak terjadi.
func (h *Penghapusan) Hapus(ctx context.Context, pelaku inti.Pelaku,
	klaimID string) (models.DampakHapus, error) {

	// Gerbang dan pembacaan dampak tetap dijalankan lebih dulu: pemanggil yang
	// tidak berwenang berhak mendapat jawaban yang benar tentang wewenangnya,
	// bukan jawaban tentang keputusan yang belum diambil.
	d, err := h.Dampak(ctx, pelaku, klaimID)
	if err != nil {
		return d, err
	}
	if err := h.tolakBilaSudahDiKomite(ctx, klaimID); err != nil {
		return d, err
	}

	// ⛔ BUTIR bb: kasus yang sudah ditutup tidak dapat diubah lagi.
	// Satu pintu untuk seluruh rute pengubah - lihat
	// services.PastikanKasusTerbuka, yang pemanggilannya ditagih penjaga
	// statik. Diperiksa SESUDAH wewenang: pemanggil yang tidak berhak tidak
	// berhak pula tahu keadaan kasusnya.
	if err := h.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return d, err
	}
	return d, fmt.Errorf("%w: kolom penandanya belum ada, dan menambahkannya "+
		"menuntut langkah migrasi yang keputusannya belum disahkan work owner",
		ErrHapusFisikDilarang)
}

// tolakBilaSudahDiKomite menolak klaim yang keputusannya sedang berjalan di
// luar modul ini.
//
// ⛔ Menghapus - fisik maupun logis - induk dari kasus komite yang masih
// berjalan meninggalkan kasus yang menunjuk klaim yang tidak ada lagi.
func (h *Penghapusan) tolakBilaSudahDiKomite(ctx context.Context, klaimID string) error {
	perBaris, err := repository.NewKlaimLife(h.svc.DB()).AmbilBaris(ctx, klaimID)
	if err != nil {
		return err
	}
	for _, daftar := range perBaris {
		for _, b := range daftar {
			if strings.TrimSpace(b.KomiteID) != "" {
				return fmt.Errorf("%w: baris %s menunjuk komite %s",
					ErrKlaimSudahDiKomite, b.ID, b.KomiteID)
			}
		}
	}
	return nil
}

// pastikanAda menolak pengenal klaim yang tidak menunjuk apa pun.
//
// ⛔ Dipisah supaya "klaim tidak ada" tidak terbaca sebagai "nol baris ikut
// terhapus": keduanya menghasilkan angka nol, dan artinya berlawanan.
func (h *Penghapusan) pastikanAda(ctx context.Context, klaimID string) (
	*models.Klaim, error) {
	hdr, err := repository.NewKlaimLife(h.svc.DB()).AmbilHeader(ctx, klaimID)
	if err != nil {
		return nil, err
	}
	if hdr == nil {
		return nil, fmt.Errorf("%w: %q", ErrKlaimTidakAda, klaimID)
	}
	return hdr, nil
}
