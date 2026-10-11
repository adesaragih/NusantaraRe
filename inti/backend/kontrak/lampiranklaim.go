package kontrak

// Lampiran kasus klaim induk yang dibaca layar Komite - perintah work owner 11-10-2026 ("TAMBAHKAN ATTACHMENT KLAIM
// FACIN KE KOMITE KLAIM FACIN, KLAIM PROP KE KOMITE KLAIM PROP ... TOMBOL ATTACHNYA DI KOMITE DI DISABLE AJA").
// Hanya-baca: unggah, hapus, dan pindah kategori tetap milik layar klaim (`GCNMSaveAttachments`). Dipakai
// `KlaimTreatyKomite` dan `KlaimFacInKomite`; penyedianya membaca master kategori dokumen klaim lininya
// (`T_KATEGORI_DOC_KLAIM` TYPE_KLAIM PROP / FAC) dan dokumen klaim (`DOCUMENT_CLAIM`) seperti tab lampiran klaim.

import (
	"errors"
	"io"
)

// KategoriLampiranKlaim - satu baris `AttachCategory.pxResults`: kategori master dan cacah dokumen kasus klaim itu
// (termasuk kategori bercacah 0 - penyaringan urusan pemakai).
type KategoriLampiranKlaim struct {
	ID     string
	Label  string
	Jumlah int
}

// LampiranKlaim - satu dokumen kasus klaim induk.
type LampiranKlaim struct {
	ID       string
	NamaFile string
	Kategori string
	// MIME - ekstensi huruf kecil (`InsertDocument_Act` S3).
	MIME string
	// Tanggal - Upload Date `DD/MM/YYYY HH24:MI`; Operator - PXCREATEOPERATOR.
	Tanggal  string
	Operator string
	// AdaObjek - objek penyimpanan tercatat (syarat View / View Office Online).
	AdaObjek bool
}

// DaftarLampiranKlaim - kategori master dan dokumen kasus klaim induk (terbaru dahulu).
type DaftarLampiranKlaim struct {
	Kategori []KategoriLampiranKlaim
	Lampiran []LampiranKlaim
}

// IsiLampiranKlaim - isi satu dokumen klaim; pemanggil wajib menutup `Isi`.
type IsiLampiranKlaim struct {
	Isi  io.ReadCloser
	Nama string
	Mime string
}

var (
	// ErrLampiranKlaimTidakAda - dokumen tidak ada di kasus klaim induk itu.
	ErrLampiranKlaimTidakAda = errors.New("kontrak: lampiran tidak ada di kasus klaim induk")
	// ErrLampiranKlaimBukanOffice - View Office Online untuk berkas selain Excel / Word / PowerPoint.
	ErrLampiranKlaimBukanOffice = errors.New("kontrak: View Office Online is only for Excel, Word, and PowerPoint files")
)
