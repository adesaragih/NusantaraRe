package services

import (
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestPosisiPutaranRiwayat - bug DEV 12-10-2026 NB-151039: POSISI riwayat Submit = IsCedingConfirm menurut FlagOnGoingPolicy
// (SetValidateDate_PostAct 3.3 / 3.4), PUTARAN = rumus CARI8 SaveViewSuggest. Tahap Offer (flag 0) tetap "Offer" / "1"; flag 1 =
// "Policy" / "2" walau tahap tidak dikenali sebagai tahap Policy; flag lain = "Policy" (tahap tiket 77, repository.TahapDariFlag).
func TestPosisiPutaranRiwayat(t *testing.T) {
	for masuk, mau := range map[string]string{"Offer": "1", "Binding": "2", "Accepted": "2", "Policy": "2", "Retrocession": "", "": ""} {
		if got := putaranRiwayat(masuk); got != mau {
			t.Errorf("putaran %q: %q mau %q", masuk, got, mau)
		}
	}
	for flag, mau := range map[string][2]string{"0": {"Offer", "1"}, "1": {"Policy", "2"}, "2": {"Binding", "2"}, "": {"Policy", "2"},
		"9": {"Policy", "2"}} {
		k := models.KonteksSubmit{Position: "Offer", FlagOnGoingPolicy: flag}
		w := Baru(nil).tulisanSubmit(TahapMarketing, "UJI-NB-1", pelakuMkt, models.IsiSubmit{EmailType: models.EmailTypeAccept}, k,
			models.AkunSubmit{Nama: "Uji"})
		if w.Riwayat.Posisi != mau[0] || w.Riwayat.Putaran != mau[1] {
			t.Errorf("flag %q: %q / %q mau %v", flag, w.Riwayat.Posisi, w.Riwayat.Putaran, mau)
		}
	}
}
