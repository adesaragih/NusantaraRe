package services_test

// Perbaikan /code-review sesudah paket 11 (01-10-2026): jejak audit menyebut pelaku, master saat simpan
// sama ketatnya dengan autocomplete, induk dikunci sebelum anak ditulis/dihapus, galat yang terbaca.

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/services"
	"nusantarare/modul/mastercontractretrolife/backend/tiruan"
)

func layananCatat(g *tiruan.Gudang) (*services.Layanan, *[]string) {
	var baris []string
	return services.BaruLayanan(g, g.Transaksi, func(s string) { baris = append(baris, s) }), &baris
}

func memuat(baris []string, potongan ...string) bool {
	for _, b := range baris {
		if containsAll(b, potongan...) {
			return true
		}
	}
	return false
}

// ⛔ Tiket 08 AC 29, tiket 09 AC jejak: siapa, kapan (stempel log), dan berapa.
func TestJejakHapusDanSalinMenyebutPelaku(t *testing.T) {
	ctx := context.Background()
	l, baris := layananCatat(gudangPohon())
	if _, err := l.Hapus(ctx, pelaku, services.HapusKontrak, "UJI-K1", models.Dampak{Security: 3, Reinsurer: 2, Business: 1}); err != nil {
		t.Fatal(err)
	}
	if !memuat(*baris, "deleted kontrak UJI-K1", "by account UJI-PELAKU") {
		t.Errorf("jejak hapus tanpa pelaku: %q", *baris)
	}

	g := gudangBusiness()
	g.Kontrak["UJI-K2"] = models.Kontrak{ID: "UJI-K2", IDTreatyYear: "UJI-T1", ReinsTypeID: "10200", ReinsTypeName: "OR"}
	g.Business["UJI-B1"] = models.Business{ID: "UJI-B1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1", ReinsTypeID: "10196",
		BizCode: "UJI-B01", BizName: "UJI BUSINESS SATU", RIRateID: "UJI-RATE-7", RIRate: "UJI"}
	l, baris = layananCatat(g)
	if _, err := l.SalinSemua(ctx, pelaku, "UJI-B1", []string{"UJI-K2"}); err != nil {
		t.Fatal(err)
	}
	if !memuat(*baris, "business UJI-B1 copied to 1", "by account UJI-PELAKU") {
		t.Errorf("jejak salin tanpa pelaku: %q", *baris)
	}
}

func TestJejakSalinanAnakMenyebutPelaku(t *testing.T) {
	g := gudangPohon()
	g.Business["UJI-B1"] = models.Business{ID: "UJI-B1", TreatyYearID: "UJI-T1", TreatyYear: "LAMA", TreatyContractID: "UJI-K1"}
	l, baris := layananCatat(g)
	th := g.Tahun["UJI-T1"]
	if _, err := l.SimpanTahun(context.Background(), pelaku, services.TahunMasuk{ID: "UJI-T1", TreatyYear: "2031",
		UnderwritingYear: "2031", StartDate: "2031-01-01", EndDate: "2031-12-31"}, false); err != nil {
		t.Fatalf("simpan tahun %+v: %v", th, err)
	}
	if !memuat(*baris, "treaty year UJI-T1 copied to", "by account UJI-PELAKU") {
		t.Errorf("jejak salinan tahun tanpa pelaku: %q", *baris)
	}
}

// ⛔ Master saat simpan = saringan autocomplete (`BrowseCedingCoLife_RD` b565/b601, `BrowseBusinessLife_RD`
// b651) untuk pilihan BARU; baris lama yang reinsurernya kini nonaktif tetap dapat disunting.
func TestReinsurerBaruHarusLifeDanAktif(t *testing.T) {
	ctx := context.Background()
	g := gudangReinsurer()
	g.MasterRe = append(g.MasterRe,
		models.MasterReinsurer{ID: "UJI-N01", ClientName: "UJI NON LIFE", StatusActive: models.StatusMasterReinsurerAktif},
		models.MasterReinsurer{ID: "UJI-L09", ClientName: "UJI NONAKTIF", StatusActive: "0"})
	l := layananUji(g)
	for _, id := range []string{"UJI-N01", "UJI-L09"} {
		m := reinsurerLengkap()
		m.ReinsurerID = id
		if _, err := l.SimpanReinsurer(ctx, pelaku, "UJI-K1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
			t.Errorf("reinsurer baru %s: %v, mau ErrMasukanTidakSah", id, err)
		}
	}
	// Baris lama ber-reinsurer nonaktif: disunting tanpa mengganti reinsurer = boleh.
	g.Reinsurer["UJI-R7"] = models.Reinsurer{ID: "UJI-R7", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1", ReinsurerID: "UJI-L09"}
	m := reinsurerLengkap()
	m.ID, m.ReinsurerID = "UJI-R7", "UJI-L09"
	if _, err := l.SimpanReinsurer(ctx, pelaku, "", m); err != nil {
		t.Errorf("ubah baris lama tanpa ganti reinsurer: %v", err)
	}
	// Mengganti reinsurer baris lama ke yang nonaktif = pilihan baru, ditolak.
	g.Reinsurer["UJI-R8"] = models.Reinsurer{ID: "UJI-R8", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1", ReinsurerID: "UJI-L01"}
	m.ID = "UJI-R8"
	if _, err := l.SimpanReinsurer(ctx, pelaku, "", m); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("ganti ke reinsurer nonaktif: %v", err)
	}
}

func TestSecurityBaruHarusLifeDanAktif(t *testing.T) {
	g := gudangReinsurer()
	g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1"}
	g.MasterRe = append(g.MasterRe, models.MasterReinsurer{ID: "UJI-N01", ClientName: "UJI NON LIFE",
		StatusActive: models.StatusMasterReinsurerAktif})
	_, err := layananUji(g).SimpanSecurity(context.Background(), pelaku, "UJI-R1",
		services.SecurityMasuk{ReinsurerName: "x", ReinsurerID: "UJI-N01", PctShare: "10"})
	if !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("security non-life: %v", err)
	}
}

func TestBusinessBaruHarusLife(t *testing.T) {
	ctx := context.Background()
	g := gudangBusiness()
	g.MasterBiz = append(g.MasterBiz, models.MasterBusiness{ID: "UJI-B02", Note: "UJI NON LIFE", OldID: "N02"})
	l := layananUji(g)
	m := businessLengkap()
	m.BizCode = "UJI-B02"
	if _, err := l.SimpanBusiness(ctx, pelaku, "UJI-K1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("business non-life baru: %v", err)
	}
	g.Business["UJI-B7"] = models.Business{ID: "UJI-B7", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1", BizCode: "UJI-B02"}
	m.ID = "UJI-B7"
	if _, err := l.SimpanBusiness(ctx, pelaku, "", m); err != nil {
		t.Errorf("ubah business lama tanpa ganti kode: %v", err)
	}
}

// ⛔ K2 tanpa FK: induk DIKUNCI (`SELECT … FOR UPDATE`) sebelum anak ditulis dan sebelum kaskade
// menghitung - anak yang lahir bersamaan dengan hapus induknya tidak dapat menjadi yatim.
func TestIndukDikunciSebelumAnakDitulisAtauDihapus(t *testing.T) {
	ctx := context.Background()
	g := gudangBusiness()
	g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1"}
	l := layananUji(g)
	if _, err := l.SimpanReinsurer(ctx, pelaku, "UJI-K1", reinsurerLengkap()); err != nil {
		t.Fatal(err)
	}
	if _, err := l.SimpanSecurity(ctx, pelaku, "UJI-R1", services.SecurityMasuk{ReinsurerName: "x", ReinsurerID: "UJI-L02", PctShare: "10"}); err != nil {
		t.Fatal(err)
	}
	if _, err := l.SimpanBusiness(ctx, pelaku, "UJI-K1", businessLengkap()); err != nil {
		t.Fatal(err)
	}
	mau := []string{"kontrak|UJI-K1", "reinsurer|UJI-R1", "kontrak|UJI-K1"}
	if !reflect.DeepEqual(g.Kunci, mau) {
		t.Errorf("kunci induk: %q, mau %q", g.Kunci, mau)
	}

	g = gudangPohon()
	if _, err := layananUji(g).Hapus(ctx, pelaku, services.HapusReinsurer, "UJI-R1", models.Dampak{Security: 2}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(g.Kunci, []string{"reinsurer|UJI-R1"}) {
		t.Errorf("kunci hapus: %q", g.Kunci)
	}
}

func TestSalinSemuaMengunciSetiapSasaran(t *testing.T) {
	g := gudangBusiness()
	g.Kontrak["UJI-K2"] = models.Kontrak{ID: "UJI-K2", IDTreatyYear: "UJI-T1", ReinsTypeID: "10200"}
	g.Kontrak["UJI-K3"] = models.Kontrak{ID: "UJI-K3", IDTreatyYear: "UJI-T1", ReinsTypeID: "10201"}
	g.Business["UJI-B1"] = models.Business{ID: "UJI-B1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1", ReinsTypeID: "10196"}
	if _, err := layananUji(g).SalinSemua(context.Background(), pelaku, "UJI-B1", []string{"UJI-K2", "UJI-K3"}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(g.Kunci, []string{"kontrak|UJI-K2", "kontrak|UJI-K3"}) {
		t.Errorf("kunci sasaran: %q", g.Kunci)
	}
}

// Tiket 09 AC: baris yang sudah tidak ada saat penghapus berjalan = pesan jelas (404), bukan 500.
func TestBarisHilangSaatHapusAdalah404(t *testing.T) {
	g := gudangPohon()
	g.SelaHapus = func() { g.Kontrak = map[string]models.Kontrak{} }
	_, err := layananUji(g).Hapus(context.Background(), pelaku, services.HapusKontrak, "UJI-K1",
		models.Dampak{Security: 3, Reinsurer: 2, Business: 1})
	if !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Errorf("kontrak hilang saat hapus: %v, mau ErrKontrakTidakAda", err)
	}
}

// ⛔ K8: kalimat 409 dibaca manusia - nol cetakan struct Go.
func TestPesan409TanpaCetakanStruct(t *testing.T) {
	_, err := layananUji(gudangPohon()).Hapus(context.Background(), pelaku, services.HapusKontrak, "UJI-K1",
		models.Dampak{Security: 1})
	if !errors.Is(err, services.ErrDampakBerubah) {
		t.Fatalf("dampak basi: %v", err)
	}
	p := services.Pesan(err)
	if strings.ContainsAny(p, "{}") || !containsAll(p, "3 security", "2 reinsurer", "1 business") {
		t.Errorf("pesan 409: %q", p)
	}
}

// ⛔ NUMBER Oracle: 38 digit dan rentang eksponen - `1E-200` dan `1E+99999` ditolak 422 menyebut medannya.
func TestAngkaDiLuarJangkauanNumberDitolak(t *testing.T) {
	ctx := context.Background()
	g := gudangReinsurer()
	l := layananUji(g)
	for _, v := range []string{"1E-200", "0.00000000000000000000000000000000000000001", "1E+2"} {
		m := reinsurerLengkap()
		m.OvrComm = v
		_, err := l.SimpanReinsurer(ctx, pelaku, "UJI-K1", m)
		if v == "1E+2" {
			if err != nil {
				t.Errorf("OVR COMM %s (seratus, 3 digit): %v", v, err)
			}
			continue
		}
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(services.Pesan(err), "(%) OVR COMM") {
			t.Errorf("OVR COMM %s: %v", v, err)
		}
	}
	m := kontrakLengkap()
	m.IDR = "1E+99999"
	if _, err := l.SimpanKontrak(ctx, pelaku, "UJI-T1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("MAXIMUM LIMIT (IDR) 1E+99999: %v", err)
	}
}

// ⛔ VARCHAR2(100)/(1000) `[data DBA]`: kelebihan panjang = 422 menyebut medannya, bukan ORA-12899 → 500.
func TestTeksMelebihiLebarKolomDitolakMenyebutMedan(t *testing.T) {
	ctx := context.Background()
	panjang := strings.Repeat("9", 101)
	l := layananUji(gudangBusiness())
	_, err := l.SimpanTahun(ctx, pelaku, services.TahunMasuk{TreatyYear: panjang, UnderwritingYear: "2026",
		StartDate: "2026-01-01", EndDate: "2026-12-31"}, true)
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(services.Pesan(err), "TRANSACTION YEAR") {
		t.Errorf("TRANSACTION YEAR 101 byte: %v", err)
	}
	m := businessLengkap()
	m.RIRateID = panjang
	if _, err := l.SimpanBusiness(ctx, pelaku, "UJI-K1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("RIRATEID 101 byte: %v", err)
	}
	m = businessLengkap()
	m.RIRate = strings.Repeat("x", 1001)
	if _, err := l.SimpanBusiness(ctx, pelaku, "UJI-K1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("R/I RATE 1001 byte: %v", err)
	}
	panjangPelaku := inti.Pelaku{AkunID: strings.Repeat("A", 101)}
	_, err = l.SimpanTahun(ctx, panjangPelaku, services.TahunMasuk{TreatyYear: "2026", UnderwritingYear: "2026",
		StartDate: "2026-01-01", EndDate: "2026-12-31"}, true)
	if !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("akun pelaku 101 byte ke USERID VARCHAR2(100): %v", err)
	}
}
