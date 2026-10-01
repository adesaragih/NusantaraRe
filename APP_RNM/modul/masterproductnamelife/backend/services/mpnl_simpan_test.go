package services_test

// Simpan sisi umum (paket 3, tiket 02) di atas gudang tiruan (kodek repository sungguhan).

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
	"nusantarare/modul/masterproductnamelife/backend/tiruan"
)

// produkMasuk - isian form lengkap yang sah (medan wajib terisi, pilihan ada di master).
func produkMasuk() models.Produk {
	return models.Produk{
		Umum: models.ProdukUmum{ProductName: "UJI PRODUK", Ceding: "UJI CEDING", CedingID: "L0UJI", SOBName: "UJI SOB",
			SOBID: "L0SOB", RIComm: "12,5", RIRisk: "UJI RISK", RIRiskID: "1000117", InwardName: "UJI PRODUK UJI PEMEGANG",
			TreatyNumber: "UJI/001", Cause: "ANY CAUSE", CauseID: "100004",
			// Medan mati dikirim klien - harus DIABAIKAN server.
			TypeBasicRider: "KLIEN", Grup: "KLIEN", CreateOp: "KLIEN", UpdateOp: "KLIEN"},
		Inward: models.ProdukInward{PolicyHolder: "UJI-ORG-1", PolicyHolderName: "UJI PEMEGANG", Begin: "2026-03-01",
			Currency: "IDR", CurrencyID: "1"},
	}
}

func layananMaster() (*services.Layanan, *tiruan.Gudang) {
	l, g := layananUji()
	g.Master[models.MasterCeding] = []models.NilaiMaster{{ID: "L0UJI", Nama: "UJI CEDING"}, {ID: "L0BARU", Nama: "UJI CEDING BARU"}}
	g.Master[models.MasterSOB] = []models.NilaiMaster{{ID: "L0SOB", Nama: "UJI SOB"}}
	g.Master[models.MasterRIRisk] = []models.NilaiMaster{{ID: "1000117", Nama: "UJI RISK"}}
	g.Master[models.MasterPenyebab] = []models.NilaiMaster{{ID: "100004", Nama: "ANY CAUSE"}}
	g.Master[models.MasterPemegangPolis] = []models.NilaiMaster{{ID: "UJI-ORG-1", Nama: "UJI PEMEGANG"}}
	g.Master[models.MasterMataUang] = []models.NilaiMaster{{ID: "1", Nama: "IDR"}}
	return l, g
}

func TestSimpanBaruIdentitasDariSequenceDanJejakPelaku(t *testing.T) {
	l, g := layananMaster()
	p, err := l.SimpanProduk(context.Background(), pelakuUji, produkMasuk(), true)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "100044" {
		t.Errorf("ID baru dari sequence '1' ‖ LPAD(44, 5): %q", p.ID)
	}
	if p.Umum.CreateOp != "UJI-PELAKU" || p.Umum.UpdateOp != "UJI-PELAKU" {
		t.Errorf("CREATEOP/UPDATEOP = pelaku, bukan isian klien: %+v", p.Umum)
	}
	if p.Umum.TypeBasicRider != "" || p.Umum.Grup != "" {
		t.Errorf("medan mati tidak pernah dari klien: %+v", p.Umum)
	}
	if p.Umum.PolicyHolder != "UJI-ORG-1" || p.Umum.PolicyHolderName != "UJI PEMEGANG" {
		t.Errorf("SaveProductName_Act 1 b359: POLICYHODER disalin dari inward: %+v", p.Umum)
	}
	if p.Umum.RIComm != "12.5" {
		t.Errorf("koma desimal diterima sebagai titik: %q", p.Umum.RIComm)
	}
	if g.Komit != 1 || !strings.Contains(g.Umum["100044"], `"POLICYHODER":"UJI-ORG-1"`) {
		t.Errorf("tersimpan satu transaksi berkunci Pega: komit %d %s", g.Komit, g.Umum["100044"])
	}
	if g.Datar["100044"] != [4]string{"1000117", "UJI RISK", "UJI PRODUK", "01/03/2026"} {
		t.Errorf("kolom datar RIRISKID, RIRISK, PRODUCTNAME, BEGIN_DATE: %v", g.Datar["100044"])
	}
}

func TestSimpanBaruMenolakIDDariKlien(t *testing.T) {
	l, g := layananMaster()
	m := produkMasuk()
	m.ID = "100001"
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, true); !errors.Is(err, services.ErrIDDariKlien) {
		t.Errorf("ADR-0006: %v", err)
	}
	if g.Komit != 0 {
		t.Error("nol tulisan")
	}
}

func TestUbahMempertahankanPembuatKunciLamaDanMedanMati(t *testing.T) {
	l, g := layananMaster()
	g.Umum["100007"] = `{"ID":"100007","PRODUCTNAME":"LAMA","CREATEOP":"UJI-PEMBUAT","TYPE":"1","GRUP":"2",
		"KUNCILAMA":"tetap","CommentList":[{"Date":"20260101T000000.000 GMT","OperatorName":"UJI-LAMA","Suggest":"x"}]}`
	m := produkMasuk()
	m.ID = "100007"
	m.CommentList = nil // klien tidak dapat menghapus riwayat
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, false)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "100007" || p.Umum.ProductName != "UJI PRODUK" || p.Umum.CreateOp != "UJI-PEMBUAT" ||
		p.Umum.UpdateOp != "UJI-PELAKU" || p.Umum.TypeBasicRider != "1" || p.Umum.Grup != "2" {
		t.Errorf("ubah: %+v", p.Umum)
	}
	if !strings.Contains(g.Umum["100007"], `"KUNCILAMA":"tetap"`) || len(p.CommentList) < 1 ||
		p.CommentList[0].OperatorName != "UJI-LAMA" {
		t.Errorf("kunci lama dan riwayat komentar dipertahankan: %s", g.Umum["100007"])
	}
	if len(g.Umum) != 1 {
		t.Errorf("ubah tidak menggandakan: %d produk", len(g.Umum))
	}
}

func TestUbahProdukTidakAda(t *testing.T) {
	l, _ := layananMaster()
	m := produkMasuk()
	m.ID = "100999"
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, false); !errors.Is(err, services.ErrProdukTidakAda) {
		t.Errorf("404: %v", err)
	}
}

func TestPilihanMasterDiverifikasiDanNamaDariMaster(t *testing.T) {
	l, _ := layananMaster()
	m := produkMasuk()
	m.Umum.CedingID, m.Umum.Ceding = "L0BARU", "DIKETIK KLIEN"
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil {
		t.Fatal(err)
	}
	if p.Umum.Ceding != "UJI CEDING BARU" {
		t.Errorf("nama dibangun ulang dari ID master: %q", p.Umum.Ceding)
	}
	m = produkMasuk()
	m.Umum.CedingID = "L0TIDAKADA"
	_, err = l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(services.Pesan(err), "Ceding") {
		t.Errorf("pilihan di luar master ditolak dengan nama medan: %v", err)
	}
	m = produkMasuk()
	m.Umum.CauseID, m.Umum.Cause = "", "DIKETIK BEBAS"
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, true); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("nama tanpa ID master (diketik bebas) ditolak: %v", err)
	}
}

func TestAngkaTidakSahDitolakDanUangTidakBerubah(t *testing.T) {
	l, _ := layananMaster()
	m := produkMasuk()
	m.Umum.RIComm = "12.3456789012345678901234567890"
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil || p.Umum.RIComm != "12.3456789012345678901234567890" {
		t.Errorf("desimal persis: %q %v", p.Umum.RIComm, err)
	}
	for _, v := range []string{"abc", "1,5.0", "1e400", "NaN"} {
		m := produkMasuk()
		m.Umum.RIComm = v
		_, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(services.Pesan(err), "Deduction (%)") {
			t.Errorf("%q harus ditolak dengan label medan: %v", v, err)
		}
	}
}

func TestPanjangKolomDatarDitolakBukanDipotong(t *testing.T) {
	l, _ := layananMaster()
	m := produkMasuk()
	m.Umum.ProductName = strings.Repeat("X", 1001)
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, true); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(services.Pesan(err), "Product Name") {
		t.Errorf("PRODUCTNAME VARCHAR2(1000): %v", err)
	}
}

func TestGagalTulisTidakPernahTampakBerhasil(t *testing.T) {
	l, g := layananMaster()
	g.GagalTulis = tiruan.ErrTiruan
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, produkMasuk(), true); !errors.Is(err, tiruan.ErrTiruan) {
		t.Errorf("galat tulis diteruskan: %v", err)
	}
	if g.Komit != 0 || len(g.Umum) != 0 {
		t.Errorf("nol baris, nol komit: %d %d", g.Komit, len(g.Umum))
	}
}
