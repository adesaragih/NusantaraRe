package services

// Predikat lini dan siklus dari registry predikat EDM (tiket E01) - pengisi
// struct Predikat (before-image, gerbang E04 langkah 10-16).

import "nusantarare/modul/endorsmentfacin/backend/services/predikat"

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
