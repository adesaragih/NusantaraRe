package services_test

// Business (paket 6, tiket 07/08): `SaveBusinessLife_Act` dan
// `SaveBusinessToAllLife_Act` tanpa procedure.

import (
	"context"
	"errors"
	"sort"
	"strings"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/services"
	"nusantarare/modul/mastercontractretrolife/backend/tiruan"
)

func gudangBusiness() *tiruan.Gudang {
	g := gudangReinsurer()
	g.MasterBiz = []models.MasterBusiness{{ID: "UJI-B01", Note: "UJI BUSINESS SATU", OldID: "L01"}}
	return g
}

func businessLengkap() services.BusinessMasuk {
	return services.BusinessMasuk{BizCode: "UJI-B01", BizName: "ketikan", RIRateID: "UJI-RATE-7",
		RIRate: "  UJI Tabel Rate, 0,5%  "}
}

func TestBusinessWajibEmpatMedanVerbatim(t *testing.T) {
	for nama, ubah := range map[string]func(*services.BusinessMasuk){
		"BUSINESS CODE": func(m *services.BusinessMasuk) { m.BizCode = "" },
		"BUSINESS NAME": func(m *services.BusinessMasuk) { m.BizName = "" },
		"RIRATEID":      func(m *services.BusinessMasuk) { m.RIRateID = "" },
		"R/I RATE":      func(m *services.BusinessMasuk) { m.RIRate = " " },
	} {
		g := gudangBusiness()
		m := businessLengkap()
		ubah(&m)
		_, err := layananUji(g).SimpanBusiness(context.Background(), pelaku, "UJI-K1", m)
		if !errors.Is(err, services.ErrWajibIsi) || services.Pesan(err) != "All value cannot be empty." || len(g.Business) != 0 {
			t.Errorf("%s kosong: %v", nama, err)
		}
	}
}

func TestBusinessRIRateTeksApaAdanyaDanSalinanInduk(t *testing.T) {
	g := gudangBusiness()
	b, err := layananUji(g).SimpanBusiness(context.Background(), pelaku, "UJI-K1", businessLengkap())
	if err != nil {
		t.Fatal(err)
	}
	// ⛔ R7: RIRATE teks - tidak dipangkas, tidak diformat, tidak diurai.
	if b.RIRate != "  UJI Tabel Rate, 0,5%  " || b.RIRateID != "UJI-RATE-7" {
		t.Errorf("RIRATE diubah: %q / %q", b.RIRate, b.RIRateID)
	}
	if b.BizName != "UJI BUSINESS SATU" || b.ID != "1000044" || b.UserID != "UJI-PELAKU" {
		t.Errorf("nama master / ID / pelaku: %+v", b)
	}
	if b.TreatyYearID != "UJI-T1" || b.TreatyYear != "2026" || b.ReinsTypeID != "10196" || b.ReinsTypeName != "QS" {
		t.Errorf("salinan tahun/kontrak (K4): %+v", b)
	}
	m := businessLengkap()
	m.BizCode = "UJI-TAK-ADA"
	if _, err := layananUji(g).SimpanBusiness(context.Background(), pelaku, "UJI-K1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("business di luar master: %v", err)
	}
}

func gudangSalinSemua() *tiruan.Gudang {
	g := gudangBusiness()
	g.Kontrak["UJI-K2"] = models.Kontrak{ID: "UJI-K2", IDTreatyYear: "UJI-T1", ReinsTypeID: "10197", ReinsTypeName: "2ND QS"}
	g.Kontrak["UJI-K3"] = models.Kontrak{ID: "UJI-K3", IDTreatyYear: "UJI-T1", ReinsTypeID: "10198", ReinsTypeName: "SURPLUS"}
	// Kontrak lain berjenis SAMA di tahun yang sama - dilewati (3.1 WhenTrue 3).
	g.Kontrak["UJI-K4"] = models.Kontrak{ID: "UJI-K4", IDTreatyYear: "UJI-T1", ReinsTypeID: "10196", ReinsTypeName: "QS"}
	// Tahun lain - di luar `GetTreatyContract_life` (`idtreatyyear =`).
	g.Tahun["UJI-T2"] = models.TahunTreaty{ID: "UJI-T2", TreatyYear: "2027"}
	g.Kontrak["UJI-K9"] = models.Kontrak{ID: "UJI-K9", IDTreatyYear: "UJI-T2", ReinsTypeID: "10199"}
	g.Business["UJI-BZ"] = models.Business{ID: "UJI-BZ", TreatyYearID: "UJI-T1", TreatyYear: "2026", TreatyContractID: "UJI-K1",
		ReinsTypeID: "10196", ReinsTypeName: "QS", BizCode: "UJI-B01", BizName: "UJI BUSINESS SATU", RIRateID: "UJI-RATE-7", RIRate: "UJI R"}
	return g
}

// ⭐ R2: sasaran = kontrak LAIN setahun yang jenisnya BERBEDA.
func TestPratinjauSalinSemuaJenisBerbedaTahunSama(t *testing.T) {
	p, err := layananUji(gudangSalinSemua()).PratinjauSalinSemua(context.Background(), pelaku, "UJI-BZ")
	if err != nil {
		t.Fatal(err)
	}
	var id []string
	for _, s := range p.Sasaran {
		id = append(id, s.ID)
	}
	sort.Strings(id)
	if strings.Join(id, ",") != "UJI-K2,UJI-K3" || p.ReinsTypeID != "10196" || p.Business.ID != "UJI-BZ" {
		t.Errorf("sasaran %v (dasar jenis %q)", id, p.ReinsTypeID)
	}
}

func TestSalinSemuaSatuTransaksiDanPesanVerbatim(t *testing.T) {
	g := gudangSalinSemua()
	var log []string
	l := services.BaruLayanan(g, g.Transaksi, func(s string) { log = append(log, s) })
	h, err := l.SalinSemua(context.Background(), pelaku, "UJI-BZ", []string{"UJI-K3", "UJI-K2"})
	if err != nil {
		t.Fatal(err)
	}
	if h.Pesan != "Copied to all reins types." || h.Jumlah != 2 {
		t.Errorf("hasil %+v", h)
	}
	baru := 0
	for _, b := range g.Business {
		if b.ID == "UJI-BZ" {
			continue
		}
		baru++
		k := g.Kontrak[b.TreatyContractID]
		if b.ReinsTypeID != k.ReinsTypeID || b.TreatyYear != "2026" || b.RIRate != "UJI R" || b.UserID != "UJI-PELAKU" {
			t.Errorf("baris salinan: %+v", b)
		}
	}
	if baru != 2 || g.Komit != 1 {
		t.Errorf("%d baris baru dalam %d transaksi, mau 2 dalam 1", baru, g.Komit)
	}
	if len(log) != 1 || !containsAll(log[0], "2", "UJI-BZ") || strings.Contains(log[0], "UJI-PELAKU") {
		t.Errorf("log server (K6, tanpa nama orang): %q", log)
	}
}

func TestSalinSemuaSasaranBerubahDitolakNolBaris(t *testing.T) {
	g := gudangSalinSemua()
	_, err := layananUji(g).SalinSemua(context.Background(), pelaku, "UJI-BZ", []string{"UJI-K2"})
	if !errors.Is(err, services.ErrDampakBerubah) || len(g.Business) != 1 {
		t.Errorf("sasaran beda dari pratinjau: %v (%d business)", err, len(g.Business))
	}
}

func TestSalinSemuaGagalDiTengahTidakAdaYangLahir(t *testing.T) {
	g := gudangSalinSemua()
	g.GagalTulis["SisipBusiness#2"] = errors.New("UJI ORA-00001")
	_, err := layananUji(g).SalinSemua(context.Background(), pelaku, "UJI-BZ", []string{"UJI-K2", "UJI-K3"})
	if err == nil || len(g.Business) != 1 {
		t.Errorf("gagal di tengah: %v, %d business (mau 1 - atomik, R4)", err, len(g.Business))
	}
}
