package services_test

// Security reinsurer (paket 5, tiket 06): `SaveSecurityReinsurerLife_Act`
// tanpa procedure - wajib-isi langkah 3 b580 VERBATIM b284; 0..100 (R3);
// salinan kunci induk dari reinsurer (K4); eksposur = share anak × share induk.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/services"
	"nusantarare/modul/mastercontractretrolife/backend/tiruan"
)

func gudangSecurity(t *testing.T) *tiruan.Gudang {
	g := gudangReinsurer()
	g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1",
		PctShare: angkaUji(t, "40")}
	return g
}

func securityLengkap() services.SecurityMasuk {
	return services.SecurityMasuk{ReinsurerName: "UJI REASURANSI DUA", ReinsurerID: "UJI-L02", PctShare: "10"}
}

func TestSecurityWajibTigaMedanVerbatim(t *testing.T) {
	for nama, ubah := range map[string]func(*services.SecurityMasuk){
		"SECURITY REINSURER NAME": func(m *services.SecurityMasuk) { m.ReinsurerName = "" },
		"REINS ID":                func(m *services.SecurityMasuk) { m.ReinsurerID = "" },
		"(%) SHARE":               func(m *services.SecurityMasuk) { m.PctShare = " " },
	} {
		g := gudangSecurity(t)
		m := securityLengkap()
		ubah(&m)
		_, err := layananUji(g).SimpanSecurity(context.Background(), pelaku, "UJI-R1", m)
		if !errors.Is(err, services.ErrWajibIsi) || services.Pesan(err) != "All value cannot be empty." || len(g.Security) != 0 {
			t.Errorf("%s kosong: %v", nama, err)
		}
	}
	m := securityLengkap()
	m.PctShare = "100.5"
	if _, err := layananUji(gudangSecurity(t)).SimpanSecurity(context.Background(), pelaku, "UJI-R1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("share 100.5: %v", err)
	}
}

func TestSecurityLahirDiBawahReinsurerDenganSalinanKunci(t *testing.T) {
	g := gudangSecurity(t)
	s, err := layananUji(g).SimpanSecurity(context.Background(), pelaku, "UJI-R1", securityLengkap())
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != "1000044" || s.TreatyReinsurerID != "UJI-R1" || s.TreatyYearID != "UJI-T1" || s.TreatyContractID != "UJI-K1" ||
		s.ReinsurerName != "UJI REASURANSI DUA" || s.UserID != "UJI-PELAKU" {
		t.Errorf("security: %+v", s)
	}
	if _, err := layananUji(g).SimpanSecurity(context.Background(), pelaku, "UJI-TAK-ADA", securityLengkap()); !errors.Is(err, services.ErrReinsurerTidakAda) {
		t.Errorf("tanpa reinsurer induk: %v", err)
	}
}

// ⛔ Jebakan terbesar modul (tiket 06): 10% di tingkat dua atas induk 40% = 4%.
func TestEksposurBerjenjangSepuluhKaliEmpatPuluh(t *testing.T) {
	g := gudangSecurity(t)
	l := layananUji(g)
	s, err := l.SimpanSecurity(context.Background(), pelaku, "UJI-R1", securityLengkap())
	if err != nil {
		t.Fatal(err)
	}
	j, err := l.DaftarSecurity(context.Background(), pelaku, "UJI-R1")
	if err != nil {
		t.Fatal(err)
	}
	if got := j.Eksposur[s.ID]; got != "4" {
		t.Errorf("eksposur %q, mau 10 × 40 / 100 = 4", got)
	}
	// Mengubah share INDUK mengubah eksposur anak tanpa menyentuh baris anak.
	induk := g.Reinsurer["UJI-R1"]
	induk.PctShare = angkaUji(t, "25")
	g.Reinsurer["UJI-R1"] = induk
	sebelum := g.Security[s.ID]
	j, _ = l.DaftarSecurity(context.Background(), pelaku, "UJI-R1")
	if got := j.Eksposur[s.ID]; got != "2.5" {
		t.Errorf("eksposur sesudah induk 25%%: %q, mau 2.5", got)
	}
	if g.Security[s.ID] != sebelum {
		t.Error("baris anak berubah padahal hanya induk yang diubah")
	}
}

func TestEksposurKosongBilaShareKosong(t *testing.T) {
	g := gudangSecurity(t)
	g.Security["UJI-S1"] = models.SecurityReinsurer{ID: "UJI-S1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1",
		TreatyReinsurerID: "UJI-R1"}
	j, err := layananUji(g).DaftarSecurity(context.Background(), pelaku, "UJI-R1")
	if err != nil {
		t.Fatal(err)
	}
	if got, ada := j.Eksposur["UJI-S1"]; !ada || got != "" {
		t.Errorf("share NULL: eksposur %q (ada=%v), mau \"\"", got, ada)
	}
}
