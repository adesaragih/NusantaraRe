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
		t.Errorf("SaveProductName_Act 1 b361: POLICYHODER disalin dari inward: %+v", p.Umum)
	}
	if p.Umum.RIComm != "12.5" {
		t.Errorf("koma desimal diterima sebagai titik: %q", p.Umum.RIComm)
	}
	s := g.Produk["100044"]
	if g.Komit != 1 || s.Inward.PolicyHolder != "UJI-ORG-1" || s.Umum.PolicyHolder != "UJI-ORG-1" {
		t.Errorf("tersimpan satu transaksi, satu kolom POLICYHOLDER untuk kedua sisi: komit %d %+v", g.Komit, s.Inward)
	}
	if s.Umum.RIRiskID != "1000117" || s.Umum.RIRisk != "UJI RISK" {
		t.Errorf("kolom RIRISKID, RIRISK: %q %q", s.Umum.RIRiskID, s.Umum.RIRisk)
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

// Tabel flat (02-10-2026): pembuat dan riwayat komentar dipertahankan dari baris tersimpan; medan layar mati
// (`TYPE`, `GRUP`) tanpa kolom flat - kosong, dan tetap tidak pernah dari klien; kunci JSON tak dikelola tidak ada (D2).
func TestUbahMempertahankanPembuatDanRiwayatKomentar(t *testing.T) {
	l, g := layananMaster()
	g.IsiJSON("100007", `{"ID":"100007","PRODUCTNAME":"LAMA","CREATEOP":"UJI-PEMBUAT",
		"KUNCILAMA":"tetap","CommentList":[{"Date":"20260101T000000.000 GMT","OperatorName":"UJI-LAMA","Suggest":"x"}]}`, "")
	m := produkMasuk()
	m.ID = "100007"
	m.CommentList = nil // klien tidak dapat menghapus riwayat
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, false)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "100007" || p.Umum.ProductName != "UJI PRODUK" || p.Umum.CreateOp != "UJI-PEMBUAT" ||
		p.Umum.UpdateOp != "UJI-PELAKU" || p.Umum.TypeBasicRider != "" || p.Umum.Grup != "" {
		t.Errorf("ubah: %+v", p.Umum)
	}
	if len(p.CommentList) != 2 || p.CommentList[0].OperatorName != "UJI-LAMA" || len(g.Produk["100007"].CommentList) != 2 {
		t.Errorf("riwayat komentar dipertahankan + satu baris simpan: %+v", p.CommentList)
	}
	if len(g.Produk) != 1 {
		t.Errorf("ubah tidak menggandakan: %d produk", len(g.Produk))
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
	// NUMBER(38,8) (K6, 02-10-2026): desimal PERSIS sampai 8 angka di belakang koma dan 30 di depannya.
	for _, v := range []string{"12.34567891", "123456789012345678901234567890.12345678"} {
		m.Umum.RIComm = v
		p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
		if err != nil || p.Umum.RIComm != v {
			t.Errorf("desimal persis %q: %q %v", v, p.Umum.RIComm, err)
		}
	}
	// Lebih dari itu DITOLAK berkalimat - tidak pernah dibulatkan Oracle diam-diam.
	for v, mau := range map[string]string{
		"12.3456789012345678901234567890": `Deduction (%) "12.3456789012345678901234567890" has more than 8 decimal places`,
		"1234567890123456789012345678901": `Deduction (%) "1234567890123456789012345678901" has more than 30 digits before the decimal point`,
	} {
		m := produkMasuk()
		m.Umum.RIComm = v
		if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, true); !strings.Contains(services.Pesan(err), mau) {
			t.Errorf("%q harus ditolak %q: %v", v, mau, err)
		}
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

// Lebar kolom flat (`RIRISKID` VARCHAR2(10), `RIRISK` VARCHAR2(100) - lebar katalog DEV - dan setiap kolom teks lain,
// tiket 01 bab 02-10-2026) diperiksa atas nilai AKHIR yang ditulis - sesudah nama diganti nama master - dan ditolak
// berkalimat sebelum SQL tulis (nol tulisan), bukan dipotong Oracle.
func TestPanjangKolomDatarDitolakBukanDipotong(t *testing.T) {
	ditolak := func(t *testing.T, siapkan func(g *tiruan.Gudang, m *models.Produk), label string) {
		t.Helper()
		l, g := layananMaster()
		m := produkMasuk()
		siapkan(g, &m)
		_, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(services.Pesan(err), label) {
			t.Errorf("%s: %v", label, err)
		}
		if g.Komit != 0 || len(g.Produk) != 0 {
			t.Errorf("%s: nol tulisan (komit %d, %d baris)", label, g.Komit, len(g.Produk))
		}
	}
	// Nama master R/I Risk lebih dari 100 byte: nilai yang AKAN ditulis ke RIRISK.
	ditolak(t, func(g *tiruan.Gudang, m *models.Produk) {
		g.Master[models.MasterRIRisk] = append(g.Master[models.MasterRIRisk], models.NilaiMaster{ID: "1000999", Nama: strings.Repeat("R", 101)})
		m.Umum.RIRiskID, m.Umum.RIRisk = "1000999", "UJI"
	}, "R/I Risk Name is 101 bytes long; the column holds at most 100")
	// ID master 11 karakter: RIRISKID VARCHAR2(10).
	ditolak(t, func(g *tiruan.Gudang, m *models.Produk) {
		g.Master[models.MasterRIRisk] = append(g.Master[models.MasterRIRisk], models.NilaiMaster{ID: "12345678901", Nama: "UJI RISK 11"})
		m.Umum.RIRiskID, m.Umum.RIRisk = "12345678901", "UJI RISK 11"
	}, "R/I Risk Name ID is 11 bytes long; the column holds at most 10")
	// `PRODUCTNAME` VARCHAR2(1000) tabel flat.
	ditolak(t, func(_ *tiruan.Gudang, m *models.Produk) {
		m.Umum.ProductName = strings.Repeat("P", 1001)
	}, "Product Name is 1001 bytes long; the column holds at most 1000")
	// Kolom anak: `USIA` VARCHAR2(200) - pesan menyebut grid dan nomor baris.
	ditolak(t, func(_ *tiruan.Gudang, m *models.Produk) {
		m.LienClause = []models.BarisLien{{Usia: strings.Repeat("U", 201), Manfaat: "50"}}
	}, "LIEN CLAUSE (Potongan Manfaat Klaim) row 1: Usia saat Klaim is 201 bytes long; the column holds at most 200")
	// NUMBER(5): bilangan bulat, paling banyak 5 digit.
	ditolak(t, func(_ *tiruan.Gudang, m *models.Produk) {
		m.Inward.MinAge = "17.5"
	}, `Minimum Age (Years) "17.5" must be a whole number`)
	ditolak(t, func(_ *tiruan.Gudang, m *models.Produk) {
		m.UnderwritingLimit = []models.BarisUWLimit{{MinAge: "18", MaxAge: "123456"}}
	}, `UNDERWRITING LIMIT row 1: Max Age "123456" has more than 5 digits`)
	// Komentar simpan ini ikut diukur: `SUGGEST` VARCHAR2(4000).
	ditolak(t, func(_ *tiruan.Gudang, m *models.Produk) {
		m.Umum.Comment = strings.Repeat("C", 4001)
	}, "Comment row 1: Comment is 4001 bytes long; the column holds at most 4000")

	// Nama panjang dari KLIEN dengan ID master sah: yang ditulis adalah nama master - diterima.
	l, g := layananMaster()
	m := produkMasuk()
	m.Umum.RIRisk = strings.Repeat("R", 101)
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil || p.Umum.RIRisk != "UJI RISK" || g.Produk[p.ID].Umum.RIRisk != "UJI RISK" {
		t.Errorf("nama klien diganti nama master sebelum diperiksa: %v %q %q", err, p.Umum.RIRisk, g.Produk[p.ID].Umum.RIRisk)
	}
	// Tepat selebar kolom: diterima utuh.
	l, g = layananMaster()
	m = produkMasuk()
	m.Umum.ProductName = strings.Repeat("X", 1000)
	if p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true); err != nil || len(g.Produk[p.ID].Umum.ProductName) != 1000 {
		t.Errorf("Product Name 1000 byte tersimpan utuh: %v", err)
	}
}

func TestGagalTulisTidakPernahTampakBerhasil(t *testing.T) {
	l, g := layananMaster()
	g.GagalTulis = tiruan.ErrTiruan
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, produkMasuk(), true); !errors.Is(err, tiruan.ErrTiruan) {
		t.Errorf("galat tulis diteruskan: %v", err)
	}
	if g.Komit != 0 || len(g.Produk) != 0 {
		t.Errorf("nol baris, nol komit: %d %d", g.Komit, len(g.Produk))
	}
}
