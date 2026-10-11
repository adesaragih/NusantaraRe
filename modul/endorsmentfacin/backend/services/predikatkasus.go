package services

// Predikat lini dan siklus dari registry predikat EDM (tiket E01) - pengisi
// struct Predikat (before-image, gerbang E04 langkah 10-16).

import (
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services/predikat"
)

// PredikatDari - menilai seluruh medan Predikat lewat registry varian EDM.
//
// ⛔ Urutan mengikat: `k` wajib sudah membawa hasil query jenis bisnis polis
// lama SESUDAH JenisBisnisUntukPredikat (`SetErrorBatalEndorsement_Act`
// langkah 7-8); bila belum, registry panic. `[terverifikasi]` IsLife varian
// EDM TIDAK membaca hasil query itu, melainkan `pyWorkPage.Quotation.BusinessOldId`
// (L1…L16) - jenis bisnis "Life" dari langkah 8 tidak menyalakannya sendiri.
func PredikatDari(k predikat.Kasus) (Predikat, error) {
	var p Predikat
	for nama, tuju := range map[string]*bool{
		"IsFire": &p.IsFire, "IsGolfInsurance": &p.IsGolfInsurance, "IsAneka": &p.IsAneka, "IsPA": &p.IsPA,
		"IsMarineCargo": &p.IsMarineCargo, "IsMBU": &p.IsMBU, "IsTravel": &p.IsTravel, "IsLife": &p.IsLife,
		"IsEDM": &p.IsEDM,
	} {
		v, err := predikat.Eval(nama, k)
		if err != nil {
			return Predikat{}, err
		}
		*tuju = v
	}
	return p, nil
}

// JalurTypeEDM - `pyWorkPage.OfferFacIn.QuotationData.Type`, dibaca gerbang
// `IsEdmAdj*` / `IsEdmExtendPeriod` / `IsEdmAddObject`.
const JalurTypeEDM = "pyWorkPage.OfferFacIn.QuotationData.Type"

// KasusSesudahRISlip - `CountPaymentEdm_Act` langkah 5 ("Set Type RI to 0 IF
// EDM RI SLIP"): `pyWorkPage.OfferFacIn.QuotationData.EdmTypeNew==4` T=2 F=3 →
// `QuotationData.Type = 0`. Langkah ini berjalan SEBELUM gerbang langkah 9, 11,
// 13-15, jadi PredikatPembayaranDari wajib dinilai atas hasilnya - RI slip
// menjadi `IsEdmAdjRefNo` (Type 0). Langkah 1 (`IsEDM`) tidak membaca Type.
// Mengembalikan salinan; `k` tidak diubah. ⚠️ Pega MENULIS `Type = 0` ke
// halaman kerja (permanen); di sini hanya salinan untuk menilai predikat -
// penulisan balik ke kasus belum diport (milik lapisan yang menyimpan kasus).
func KasusSesudahRISlip(k predikat.KasusEDM, edmTypeNew string) (predikat.KasusEDM, error) {
	riSlip, err := samaKode(edmTypeNew, models.EdmTypeNewRISlip)
	if err != nil {
		return predikat.KasusEDM{}, err
	}
	salinan := k
	salinan.Properti = make(map[string]string, len(k.Properti)+1)
	for j, v := range k.Properti {
		salinan.Properti[j] = v
	}
	if riSlip {
		salinan.Properti[JalurTypeEDM] = string(models.TypeRISlip)
	}
	return salinan, nil
}
