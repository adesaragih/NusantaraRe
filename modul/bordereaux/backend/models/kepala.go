package models

// Kepala grid detail menurut format bordereaux 2025 (`D:\XML\RNM_BRD\Bordereaux\IPR - Format_Bordereaux_2025.xlsx`,
// satu sheet per LOB; perintah work owner 05-10-2026: "samain tampilan tabel utamanya saja ... ikuti/samain headernya
// saja ... ubah yang premium fire dulu"; lalu "bisnis lain perbaiki, ikuti sheet nya"; kata "Indonesia Re" di sheet diganti
// "Nusantara Re", 05-10-2026). HANYA
// tampilan grid: judul CSV (`KolomCSV.Judul`), pembacaan Upload CSV, dan pesan validasi tidak berubah.

// KepalaGrid - judul tampilan satu kolom grid detail.
type KepalaGrid struct {
	// Grup - judul sel gabungan di baris kepala pertama (mis. `*PERIOD OF INSURANCE` di atas START dan END); kosong =
	// kolom tanpa grup (judulnya menempati dua baris kepala).
	Grup string
	// Judul - teks kepala kolom persis seperti sheet-nya (boleh ber-baris baru).
	Judul string
	// Sembunyi - kolom tidak ada di format Excel; tidak ditampilkan grid (datanya tetap tersimpan).
	Sembunyi bool
}

// kepalaExcel (kepala_gen.go, dibangkitkan dari sheet Excel) - per kode kombinasi (`KombinasiBdx.Kode`), per kolom
// tabel. Kombinasi tanpa entri memakai judul CSV.

// AdaKepalaExcel - kombinasi ini sudah memakai kepala format Excel.
func AdaKepalaExcel(k KombinasiBdx) bool {
	_, ada := kepalaExcel[k.Kode]
	return ada
}

// KepalaKolom - judul tampilan kolom `c` kombinasi `k`: menurut format Excel bila kombinasinya sudah dipetakan,
// selain itu judul CSV tanpa grup.
func KepalaKolom(k KombinasiBdx, c KolomCSV) KepalaGrid {
	if peta, ada := kepalaExcel[k.Kode]; ada {
		if kg, ada := peta[c.Kolom]; ada {
			return kg
		}
	}
	return KepalaGrid{Judul: c.Judul}
}
