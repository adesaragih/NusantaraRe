package models

// Untuk apa berkas ini: LAMPIRAN KLAIM INDUK di layar komite - perintah work owner 11-10-2026 ("TAMBAHKAN ATTACHMENT
// ... KLAIM PROP KE KOMITE KLAIM PROP. YANG MUNCUL HANYA KATEGORI YANG ADA ISI NYA. JIKA KATEGORI NYA 0 TIDAK USAH DI
// MUNCULIN!"; "TOMBOL ATTACHNYA DI KOMITE DI DISABLE AJA"). Grid kategori seperti tab lampiran Claim Prop
// (`AttachCategory.pxResults`), dibaca lewat `kontrak.KlaimTreatyKomite`; komite hanya melihat dan mengunduh.

import "nusantarare/inti/backend/kontrak"

// KategoriLampiran - satu baris grid kategori: kategori master PROP dan cacah dokumen klaim induk.
type KategoriLampiran struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	CountAttach int    `json:"countAttach"`
}

// Lampiran - satu dokumen klaim induk (jendela View File).
type Lampiran struct {
	ID       string `json:"id"`
	NamaFile string `json:"namaFile"`
	Kategori string `json:"kategori"`
	// MIME - ekstensi huruf kecil.
	MIME string `json:"mime"`
	// Tanggal - Upload Date `DD/MM/YYYY HH24:MI`; Operator - Attached By.
	Tanggal  string `json:"tanggal"`
	Operator string `json:"operator"`
	// AdaObjek - objek penyimpanan tercatat (syarat View / View Office Online).
	AdaObjek bool `json:"adaObjek"`
}

// LampiranKomite - grid kategori dan dokumen klaim induk yang tampil di layar komite.
type LampiranKomite struct {
	Kategori []KategoriLampiran `json:"kategori"`
	Lampiran []Lampiran         `json:"lampiran"`
}

// SusunLampiranKomite - HANYA kategori bercacah lebih dari 0 (urutan master dipertahankan) dan dokumen berkategori
// itu; daftar kosong = [] (bukan null) di JSON.
func SusunLampiranKomite(d kontrak.DaftarLampiranKlaim) LampiranKomite {
	out := LampiranKomite{Kategori: []KategoriLampiran{}, Lampiran: []Lampiran{}}
	berisi := map[string]bool{}
	for _, k := range d.Kategori {
		if k.Jumlah > 0 {
			berisi[k.ID] = true
			out.Kategori = append(out.Kategori, KategoriLampiran{ID: k.ID, Label: k.Label, CountAttach: k.Jumlah})
		}
	}
	for _, a := range d.Lampiran {
		if berisi[a.Kategori] {
			out.Lampiran = append(out.Lampiran, Lampiran{ID: a.ID, NamaFile: a.NamaFile, Kategori: a.Kategori,
				MIME: a.MIME, Tanggal: a.Tanggal, Operator: a.Operator, AdaObjek: a.AdaObjek})
		}
	}
	return out
}
