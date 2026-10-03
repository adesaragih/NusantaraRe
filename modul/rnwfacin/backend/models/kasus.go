// Package models memuat tipe data modul RNW Fac In.
package models

import (
	"time"

	"nusantarare/inti/backend/kontrak"
)

// StatusBusinessRenewal - `QuotationData.StatusBusiness` kasus renewal: pembeda
// siklus (NB 1, Renewal 2, Endorsement 3 - nbfacin butir 29). `[terverifikasi]`
// fixture kasus RNW nyata (`rnw-fire-1`) bernilai 2.
const StatusBusinessRenewal = "2"

// Kasus - halaman kerja kasus Fac In sebagai jalur properti Pega → nilai, mis.
// `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness`. Bentuk yang sama dibaca
// mesin NB lewat kontrak.KasusFacIn.
type Kasus map[string]string

var _ kontrak.KasusFacIn = Kasus{}

// Nilai - jalur yang tidak ada dibaca kosong, seperti clipboard Pega.
func (k Kasus) Nilai(jalur string) (string, bool) {
	v, ada := k[jalur]
	return v, ada
}

// Salin - salinan dangkal; mengubahnya tidak mengubah asal.
func (k Kasus) Salin() Kasus {
	b := make(Kasus, len(k))
	for j, v := range k {
		b[j] = v
	}
	return b
}

// KasusRenewal - kasus renewal yang baru dibuat (tiket R01).
type KasusRenewal struct {
	Kasus Kasus
	// TanggalRenewal - masukan layar. Belum ditulis ke `QuotationData.RNWDate`:
	// format clipboard Pega `yyyyMMddTHHmmss.SSS GMT` menuntut konversi zona waktu
	// yang belum terverifikasi; milik R02 (layar periode).
	TanggalRenewal time.Time
	// Catatan - "Note" layar masuk. `[pertanyaan terbuka]` properti Pega-nya:
	// tidak terikat di section PeriodeRenewal.
	Catatan string
}
