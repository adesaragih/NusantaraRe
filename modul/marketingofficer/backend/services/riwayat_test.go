package services_test

// Log perubahan dari MARKETINGOFFICER_LOG (baris LAMA per UPDATE) - TANPA Oracle.

import (
	"reflect"
	"testing"

	"nusantarare/modul/marketingofficer/backend/models"
	"nusantarare/modul/marketingofficer/backend/services"
	"nusantarare/modul/marketingofficer/backend/tiruan"
)

func mo(sebagian func(*models.MarketingOfficer)) models.MarketingOfficer {
	m := models.MarketingOfficer{ID: "10000103", ClientID: "UJI-KONTAK-3", ClientName: "UJI Pega Lama", ClientID2: "10000101",
		MOLeader: "UJI Leader Satu", MOStatus: "1", BranchDetailID: "UJI-B02", BranchDetailName: "UJI Cabang Dua", TeamGroup: "2"}
	sebagian(&m)
	return m
}

// Versi: log lama Pega (tanpa LOG_TIME, AKSES_LOGIN tidak diketahui) -> log Go (bertanda) -> sekarang.
func TestSusunRiwayatTerbaruDuluDenganPerkiraanDanAksesTakDiketahui(t *testing.T) {
	pega := models.BarisLog{MarketingOfficer: mo(func(m *models.MarketingOfficer) { m.MOStatus = "2"; m.Tanggal = "" })}
	goLog := models.BarisLog{MarketingOfficer: mo(func(m *models.MarketingOfficer) {
		m.AksesLogin, m.UserUpdate, m.Tanggal = "UJI-MKT01", "UJIOPERATOR", "2026-10-01 08:00"
	}), Aksi: models.AksiGo, LogTime: "2026-10-03 10:00:01"}
	kini := mo(func(m *models.MarketingOfficer) {
		m.AksesLogin, m.BranchDetailID, m.BranchDetailName, m.TeamGroup = "UJI-MKT02", "UJI-B01", "UJI Cabang Satu", "1"
		m.UserUpdate, m.Tanggal = "UJI-ADMIN", "2026-10-03 10:00"
	})
	r := services.SusunRiwayat(kini, []models.BarisLog{pega, goLog})
	if r.ID != "10000103" || r.JumlahLog != 2 || len(r.Perubahan) != 2 {
		t.Fatalf("riwayat %+v", r)
	}
	terbaru := r.Perubahan[0]
	if terbaru.Waktu != "2026-10-03 10:00:01" || terbaru.Oleh != "UJI-ADMIN" || terbaru.Perkiraan {
		t.Errorf("perubahan terbaru: %+v", terbaru)
	}
	mau := []services.RuasBerubah{
		{Kolom: "AKSES_LOGIN", Sebelum: "UJI-MKT01", Sesudah: "UJI-MKT02"},
		{Kolom: "BRANCHDETAILID", Sebelum: "UJI-B02", Sesudah: "UJI-B01"},
		{Kolom: "BRANCHDETAILNAME", Sebelum: "UJI Cabang Dua", Sesudah: "UJI Cabang Satu"},
		{Kolom: "TEAMGROUP", Sebelum: "2", Sesudah: "1"},
	}
	if !reflect.DeepEqual(terbaru.Ruas, mau) {
		t.Errorf("ruas terbaru:\n dapat %+v\n mau   %+v", terbaru.Ruas, mau)
	}
	// Baris Pega: waktu dari TANGGAL versi baru, ditandai perkiraan; AKSES_LOGIN tidak dibandingkan (tidak diketahui).
	tertua := r.Perubahan[1]
	if tertua.Waktu != "2026-10-01 08:00" || !tertua.Perkiraan || tertua.Oleh != "UJIOPERATOR" ||
		!reflect.DeepEqual(tertua.Ruas, []services.RuasBerubah{{Kolom: "MOSTATUS", Sebelum: "2", Sesudah: "1"}}) {
		t.Errorf("perubahan tertua: %+v", tertua)
	}
	if kosong := services.SusunRiwayat(kini, nil); len(kosong.Perubahan) != 0 || kosong.Perubahan == nil {
		t.Errorf("tanpa log: %+v", kosong)
	}
}

// Ubah lewat layanan: trigger (tiruan) menulis baris lama, layanan menandainya dengan AKSES_LOGIN lama; riwayat
// lalu menunjukkan pergantian akun oleh pelaku.
func TestUbahMenandaiLogDanRiwayatMembacanya(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	isi := models.Isian{AksesLogin: "UJI-MKT02", LeaderID: "10000101", BranchParent: "UJI-P00", BranchDetailID: "UJI-B02", Aktif: true}
	if _, err := l.Ubah(ctx, admin, "10000103", isi); err != nil {
		t.Fatal(err)
	}
	if d := g.Log["10000103"]; len(d) != 1 || d[0].Aksi != models.AksiGo || d[0].AksesLogin != "UJIOPERATORLAMA" {
		t.Fatalf("log ditandai: %+v", d)
	}
	r, err := l.Riwayat(ctx, "10000103")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Perubahan) != 1 || r.Perubahan[0].Oleh != "UJI-ADMIN" ||
		!reflect.DeepEqual(r.Perubahan[0].Ruas, []services.RuasBerubah{{Kolom: "AKSES_LOGIN", Sebelum: "UJIOPERATORLAMA", Sesudah: "UJI-MKT02"}}) {
		t.Errorf("riwayat: %+v", r)
	}
	if _, err := l.Riwayat(ctx, "19999999"); err != services.ErrTidakAda {
		t.Errorf("riwayat MO tidak ada: %v", err)
	}
}
