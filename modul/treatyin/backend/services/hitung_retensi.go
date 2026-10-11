package services

import inti "nusantarare/inti/backend"

// RUMUS TAB MAXIMUM RETENTION (cabang NON-PROPORSIONAL) - dari ekspor Pega.
//
// ---------------------------------------------------------------------
// ⭐ DUA ATURAN, DARI SATU ACTIVITY
// ---------------------------------------------------------------------
//
//	Activity/TreatyInNPSetTotal.xml   param.type == "retention" (langkah 2-3)
//	                                  tombol `Update Total`
//	Activity/TreatyInNonAddItem.xml   param.Type == "retention" (langkah 2)
//	                                  tombol `Add`
//
// ---------------------------------------------------------------------
// ⛔ TAB INI LEBIH SEDERHANA DARIPADA EGNPI, DAN ITU BUKAN KELALAIAN
// ---------------------------------------------------------------------
// EGNPI punya DUA tombol; retensi hanya SATU. Sebabnya: retensi nol punya
// kolom `Amount in IDR`, jadi nol konversi kurs yang perlu dijalankan.
// Ekspor tab Maximum Retention memang hanya memuat `Update Total` -
// menambahkan `Update Value` di sini berarti membuat tombol yang di Pega
// tidak ada dan tidak punya pekerjaan.
//
// ⚠️ Arah ketergantungannya justru TERBALIK: tab EGNPI yang membaca retensi,
// bukan sebaliknya - `TreatyInNonAddItem(egnpi)` mengambil mata uang baris
// baru dari `TreatyIn.Retention(1)`. Lihat `TambahBarisEgnpi`.

// Aksi tab Maximum Retention.
const (
	// AksiRetensiTotal - tombol `Update Total` (`TreatyInNPSetTotal` retention).
	AksiRetensiTotal = "total"
	// AksiRetensiTambah - tombol `Add` (`TreatyInNonAddItem` Type=retention).
	AksiRetensiTambah = "tambah"
	// AksiRetensiHapus - tombol `Delete` pada baris grid.
	AksiRetensiHapus = "hapus"
)

// RetensiNP - satu baris tab Maximum Retention.
//
// ⚠️ `Note` dan `ClassOfBusiness` ADA di data dan nol terlihat di grid:
// `Note` hanya di rinciannya (`Section/MaxRetention.xml`), `ClassOfBusiness`
// nol di keduanya. Keduanya dibawa apa adanya - medan yang dibuang di
// perjalanan tidak dapat dikembalikan saat Save.
type RetensiNP struct {
	ID                string `json:"ID"`
	TreatyGroup       string `json:"TreatyGroup"`
	TreatyGroupID     string `json:"TreatyGroupID"`
	Currency          string `json:"Currency"`
	CurrencyID        string `json:"CurrencyID"`
	Amount            string `json:"Amount"`
	ClassOfBusiness   string `json:"ClassOfBusiness"`
	ClassOfBusinessID string `json:"ClassOfBusinessID"`
	Note              string `json:"Note"`
}

// MasukanRetensi - satu aksi tab Maximum Retention beserta isian layarnya.
type MasukanRetensi struct {
	Aksi    string      `json:"aksi"`
	Retensi []RetensiNP `json:"retensi"`
	// Indeks baris (mulai 0) untuk `hapus`.
	Indeks int `json:"indeks"`
}

// HasilRetensi - baris sesudah aksi, berikut total per mata uang.
type HasilRetensi struct {
	Retensi []RetensiNP `json:"retensi"`
	// TotalRetentionAmountNP - grid `Total Retention Amount` | `Value`.
	TotalRetentionAmountNP []NilaiMataUang `json:"TotalRetentionAmountNP"`
	Pesan                  []string        `json:"pesan"`
}

// HitungRetensi - bentuk ber-pelaku untuk handler.
func (l *Layanan) HitungRetensi(p inti.Pelaku, m MasukanRetensi) (HasilRetensi, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilRetensi{}, err
	}
	return HitungRetensi(m), nil
}

// HitungRetensi menjalankan satu aksi tab Maximum Retention. Murni.
func HitungRetensi(m MasukanRetensi) HasilRetensi {
	h := HasilRetensi{Retensi: append([]RetensiNP{}, m.Retensi...), Pesan: []string{}}

	switch m.Aksi {
	case AksiRetensiTambah:
		h.Retensi = TambahBarisRetensi(h.Retensi)
	case AksiRetensiHapus:
		if m.Indeks >= 0 && m.Indeks < len(h.Retensi) {
			h.Retensi = append(h.Retensi[:m.Indeks], h.Retensi[m.Indeks+1:]...)
		}
	}

	// ⭐ Total SELALU dihitung ulang, apa pun aksinya - alasan yang sama
	// dengan tab EGNPI: total basi di sebelah baris yang baru berubah lebih
	// menyesatkan daripada total yang belum pernah dihitung.
	h.TotalRetentionAmountNP = NPSetTotalRetensi(h.Retensi)
	return h
}

// NPSetTotalRetensi - `Activity/TreatyInNPSetTotal.xml` cabang
// `param.type == "retention"` (langkah 2-3):
//
//	[2] kosongkan TreatyIn.TotalRetentionAmountNP
//	[3] per baris Retention: SIGMA .Amount per MATA UANG
//	    (pola `Appendflag`: gabung ke baris bermata uang sama, atau tambah)
//
// ⛔ Nol pagar bagi-nol di sini, dan itu benar: cabang retensi tidak pernah
// membagi. Menyalin pagar dari cabang EGNPI akan menolak data yang sah.
func NPSetTotalRetensi(rows []RetensiNP) []NilaiMataUang {
	total := []NilaiMataUang{}
	for _, r := range rows {
		total = tambahPerMataUang(total, r.Currency, r.CurrencyID, angka(r.Amount))
	}
	return total
}

// TambahBarisRetensi - `Activity/TreatyInNonAddItem.xml` cabang
// `param.Type == "retention"`: satu baris KOSONG, `ID = ""`.
//
// ⚠️ Nol nilai awal diwarisi dari mana pun - beda dengan cabang `egnpi`
// yang mewarisi mata uang `Retention(1)`. Arah ketergantungannya satu jalan.
func TambahBarisRetensi(rows []RetensiNP) []RetensiNP {
	return append(rows, RetensiNP{})
}
