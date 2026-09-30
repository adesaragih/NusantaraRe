// Package unggah memuat gerbang berkas unggahan yang dipakai bersama:
// batas ukuran, penulisan berbatas ke disk, pengenal berkas (`ImageID`), dan
// mime dokumen.
//
// Refactor bentuk B (30-09-2026): dulu bagian `services/unggahan.go`,
// `models/imageid.go`, dan `models/mimedokumen.go`.
package unggah

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	// ErrUnggahanDirBelumDisetel - folder unggahan belum dipilih.
	//
	// ⛔ Gagal TERANG, bukan bawaan diam-diam ke `./unggahan`. Berkas
	// pelanggan tidak boleh mendarat di folder kerja siapa pun yang
	// kebetulan menjalankan server.
	ErrUnggahanDirBelumDisetel = errors.New(
		"services: UNGGAHAN_DIR belum disetel; unggahan dokumen ditolak")
	// ErrBerkasTerlaluBesar - melebihi BatasUkuranUnggahan.
	ErrBerkasTerlaluBesar = errors.New("services: berkas melebihi batas ukuran")
	// ErrBerkasKosong - nol byte.
	ErrBerkasKosong = errors.New("services: berkas kosong")
)

// Jenis efek outbox untuk berkas - butir be.
//
// ⛔ TEKS, dan tersimpan di kolom `JENIS_EFEK`. Ia dibaca kembali oleh
// proses lain, mungkin berhari-hari kemudian: nilainya tidak boleh berubah
// hanya karena sebuah konstanta Go diganti namanya.
const (
	// JenisEfekStorageUnggah - padanan `InsertGoogleStorage_Act` b1023.
	JenisEfekStorageUnggah = "storage-unggah"
	// JenisEfekStorageHapus - padanan `DeleteGoogleStorage_Act`.
	JenisEfekStorageHapus = "storage-hapus"
)

// BatasUkuranUnggahan adalah batas byte satu berkas.
//
// ⚠️ `[keputusan kami - tidak ada di korpus]` Nol rule di `Claim Life`
// menyebutkan batas ukuran: Pega menyerahkannya ke `dragDropFileUpload`
// bawaan platform, yang batasnya hidup di konfigurasi sistem dan tidak ikut
// diekspor. 25 MiB dipilih sebab ia lebih besar daripada seluruh jenis
// dokumen di tabel MIME 48 baris yang masuk akal dilampirkan pada klaim, dan
// cukup kecil supaya satu permintaan tidak menahan memori proses.
//
// Angkanya berdiri di SINI, bernama, supaya ia dapat dibantah - bukan
// tersebar sebagai literal di handler.
const BatasUkuranUnggahan = 25 << 20

// BerkasMasuk adalah satu berkas yang sedang diunggah.
//
// ⚠️ `Isi` pembaca, bukan `[]byte`. Berkas 25 MiB yang dibaca seluruhnya ke
// memori sebelum diperiksa adalah 25 MiB yang dapat diminta siapa saja,
// berkali-kali, sebelum satu pun gerbang berjalan.
type BerkasMasuk struct {
	NamaFile string
	// Mime dari pemanggil. Kosong berarti diturunkan dari nama berkas.
	Mime string
	// Kategori adalah `KATEGORI_2` - yang `DocumentLife.xml` tampilkan.
	Kategori string
	Isi      io.Reader
}

// TulisBerkas menyalin isi ke jalur itu, berbatas ukuran.
//
// Refactor bentuk B (30-09-2026): isi `Unggahan.tulisBerkas` dipindah apa
// adanya ke sini, sebab Treaty Contract Out memakai gerbang yang sama.
// `batas` <= 0 berarti `BatasUkuranUnggahan`.
func TulisBerkas(jalur string, isi io.Reader, batas int64) error {
	if err := os.MkdirAll(filepath.Dir(jalur), 0o750); err != nil {
		return fmt.Errorf("services: menyiapkan folder unggahan: %w", err)
	}
	f, err := os.OpenFile(jalur, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o640)
	if err != nil {
		return fmt.Errorf("services: membuat berkas unggahan: %w", err)
	}
	// Satu byte LEBIH daripada batas, supaya kelebihan dapat dibedakan dari
	// berkas yang panjangnya tepat sebesar batas.
	if batas <= 0 {
		batas = BatasUkuranUnggahan
	}
	n, salinErr := io.Copy(f, io.LimitReader(isi, batas+1))
	tutupErr := f.Close()
	switch {
	case salinErr != nil:
		_ = os.Remove(jalur)
		return fmt.Errorf("services: menulis berkas unggahan: %w", salinErr)
	case tutupErr != nil:
		_ = os.Remove(jalur)
		return fmt.Errorf("services: menutup berkas unggahan: %w", tutupErr)
	case n == 0:
		_ = os.Remove(jalur)
		return ErrBerkasKosong
	case n > batas:
		_ = os.Remove(jalur)
		return fmt.Errorf("%w: %d byte, batas %d", ErrBerkasTerlaluBesar, n, batas)
	}
	return nil
}
