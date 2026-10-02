package services_test

// Reinsurer (paket 4, tiket 05/11): `SaveSecurityLife_Act` tanpa procedure -
// wajib-isi langkah 3 b595 VERBATIM b299; 0..100 untuk PCTSHARE/COMMISION
// (maksud `SetErrorMessageReinsurer`, R3); total share `CountingPercentShare_Act`
// dijumlah dan TIDAK memblokir; laporan kontrak total ≠ 100.

import (
	"context"
	"errors"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/services"
	"nusantarare/modul/mastercontractretrolife/backend/tiruan"
)

func angkaUji(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	d, _, err := apd.NewFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func gudangReinsurer() *tiruan.Gudang {
	g := gudangKontrak()
	g.Kontrak["UJI-K1"] = models.Kontrak{ID: "UJI-K1", IDTreatyYear: "UJI-T1", ReinsTypeID: "10196", ReinsTypeName: "QS"}
	g.MasterRe = []models.MasterReinsurer{{ID: "UJI-L01", ClientName: "UJI REASURANSI SATU", StatusActive: models.StatusMasterReinsurerAktif},
		{ID: "UJI-L02", ClientName: "UJI REASURANSI DUA", StatusActive: models.StatusMasterReinsurerAktif}}
	return g
}

func reinsurerLengkap() services.ReinsurerMasuk {
	return services.ReinsurerMasuk{ReinsurerName: "UJI REASURANSI SATU", ReinsurerID: "UJI-L01", PctShare: "40",
		Komisi: "27,5", OvrComm: "2.5"}
}

func TestReinsurerWajibLimaMedanVerbatim(t *testing.T) {
	kosong := map[string]func(*services.ReinsurerMasuk){
		"REINSURER NAME": func(m *services.ReinsurerMasuk) { m.ReinsurerName = "" },
		"REINS ID":       func(m *services.ReinsurerMasuk) { m.ReinsurerID = "" },
		"(%) SHARE":      func(m *services.ReinsurerMasuk) { m.PctShare = "" },
		"(%) DISCOUNT":   func(m *services.ReinsurerMasuk) { m.Komisi = "" },
		"(%) OVR COMM":   func(m *services.ReinsurerMasuk) { m.OvrComm = "" },
	}
	for nama, ubah := range kosong {
		g := gudangReinsurer()
		m := reinsurerLengkap()
		ubah(&m)
		_, err := layananUji(g).SimpanReinsurer(context.Background(), pelaku, "UJI-K1", m)
		if !errors.Is(err, services.ErrWajibIsi) || services.Pesan(err) != "All value cannot be empty." || len(g.Reinsurer) != 0 {
			t.Errorf("%s kosong: %v (%d baris)", nama, err, len(g.Reinsurer))
		}
	}
}

func TestReinsurerShareDanKomisi0Sampai100OvrCommTidak(t *testing.T) {
	for _, kasus := range []struct {
		ubah  func(*services.ReinsurerMasuk)
		tolak bool
		medan string
	}{
		{func(m *services.ReinsurerMasuk) { m.PctShare = "100.0000001" }, true, "(%) SHARE"},
		{func(m *services.ReinsurerMasuk) { m.PctShare = "-0.1" }, true, "(%) SHARE"},
		{func(m *services.ReinsurerMasuk) { m.Komisi = "101" }, true, "(%) DISCOUNT"},
		{func(m *services.ReinsurerMasuk) { m.PctShare, m.Komisi = "100", "0" }, false, ""},
		// `SetErrorMessageReinsurer` tidak memeriksa OVR_COMM - tidak ditegakkan.
		{func(m *services.ReinsurerMasuk) { m.OvrComm = "150" }, false, ""},
	} {
		g := gudangReinsurer()
		m := reinsurerLengkap()
		kasus.ubah(&m)
		_, err := layananUji(g).SimpanReinsurer(context.Background(), pelaku, "UJI-K1", m)
		if kasus.tolak != (err != nil) {
			t.Errorf("%+v: galat %v, mau tolak=%v", m, err, kasus.tolak)
			continue
		}
		if kasus.tolak && (!errors.Is(err, services.ErrMasukanTidakSah) || !containsAll(services.Pesan(err), kasus.medan, "0", "100")) {
			t.Errorf("%+v: pesan %q harus menyebut %s dan rentangnya", m, services.Pesan(err), kasus.medan)
		}
	}
}

func TestReinsurerDariMasterDanSalinanKontrak(t *testing.T) {
	g := gudangReinsurer()
	m := reinsurerLengkap()
	m.ReinsurerName = "UJI NAMA KETIKAN"
	r, err := layananUji(g).SimpanReinsurer(context.Background(), pelaku, "UJI-K1", m)
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != "1000044" || r.ReinsurerName != "UJI REASURANSI SATU" || r.UserID != "UJI-PELAKU" {
		t.Errorf("identitas / nama master / pelaku: %+v", r)
	}
	if r.TreatyYearID != "UJI-T1" || r.TreatyContractID != "UJI-K1" || r.ReinsTypeID != "10196" || r.ReinsTypeName != "QS" {
		t.Errorf("salinan kontrak (K4): %+v", r)
	}
	if r.Komisi.Text('f') != "27.5" || r.OvrComm.Text('f') != "2.5" {
		t.Errorf("komisi koma -> titik: %s %s", r.Komisi.Text('f'), r.OvrComm.Text('f'))
	}
	m.ReinsurerID = "UJI-TAK-ADA"
	if _, err := layananUji(g).SimpanReinsurer(context.Background(), pelaku, "UJI-K1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("reinsurer di luar master: %v", err)
	}
}

func TestReinsurerUbahTetapDiKontraknya(t *testing.T) {
	g := gudangReinsurer()
	g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1"}
	m := reinsurerLengkap()
	m.ID = "UJI-R1"
	r, err := layananUji(g).SimpanReinsurer(context.Background(), pelaku, "", m)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Reinsurer) != 1 || r.PctShare.Text('f') != "40" {
		t.Errorf("upsert reinsurer: %+v (%d)", r, len(g.Reinsurer))
	}
	if _, err := layananUji(g).SimpanReinsurer(context.Background(), pelaku, "UJI-K-LAIN", m); !errors.Is(err, services.ErrReinsurerTidakAda) {
		t.Errorf("reinsurer lewat kontrak lain: %v", err)
	}
}

// ⛔ Total share TIDAK memblokir (tiket 05 AC, `[keputusan work owner]`).
func TestTotalShareDijumlahDanTidakMemblokir(t *testing.T) {
	g := gudangReinsurer()
	l := layananUji(g)
	for _, share := range []string{"60", "50"} {
		m := reinsurerLengkap()
		m.PctShare = share
		if _, err := l.SimpanReinsurer(context.Background(), pelaku, "UJI-K1", m); err != nil {
			t.Fatalf("share %s: %v (total > 100 tidak boleh memblokir)", share, err)
		}
	}
	j, err := l.DaftarReinsurer(context.Background(), pelaku, "UJI-K1")
	if err != nil {
		t.Fatal(err)
	}
	if j.TotalShare != "110" || !j.TotalBukan100 {
		t.Errorf("total %q bukan100 %v, mau 110 / true", j.TotalShare, j.TotalBukan100)
	}
}

func TestLaporanTotalShareBukan100(t *testing.T) {
	g := gudangReinsurer()
	g.Kontrak["UJI-K2"] = models.Kontrak{ID: "UJI-K2", IDTreatyYear: "UJI-T1", ReinsTypeName: "2ND QS"}
	g.Kontrak["UJI-K3"] = models.Kontrak{ID: "UJI-K3", IDTreatyYear: "UJI-T1", ReinsTypeName: "SURPLUS"}
	g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1",
		PctShare: angkaUji(t, "33.3333333333333333333")}
	g.Reinsurer["UJI-R2"] = models.Reinsurer{ID: "UJI-R2", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1",
		PctShare: angkaUji(t, "66.6666666666666666667")}
	g.Reinsurer["UJI-R3"] = models.Reinsurer{ID: "UJI-R3", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K2",
		PctShare: angkaUji(t, "120")}
	d, err := layananUji(g).LaporanTotalShareBukan100(context.Background(), pelaku, "")
	if err != nil {
		t.Fatal(err)
	}
	ada := map[string]string{}
	for _, b := range d {
		ada[b.KontrakID] = b.TotalShare + "|" + b.Selisih
	}
	// Tepat 100 (desimal presisi arbitrer) TIDAK muncul; >100 dan tanpa reinsurer (0) muncul.
	if _, muncul := ada["UJI-K1"]; muncul || ada["UJI-K2"] != "120|-20" || ada["UJI-K3"] != "0|100" || len(ada) != 2 {
		t.Errorf("laporan: %v", ada)
	}
	d, _ = layananUji(g).LaporanTotalShareBukan100(context.Background(), pelaku, "UJI-TAHUN-LAIN")
	if len(d) != 0 {
		t.Errorf("saringan tahun: %v", d)
	}
}
