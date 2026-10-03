package services_test

// Aturan Marketing Officer TANPA Oracle (gudang tiruan). Sumber setiap aturan: dokumentasi paket services.

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/marketingofficer/backend/models"
	"nusantarare/modul/marketingofficer/backend/services"
	"nusantarare/modul/marketingofficer/backend/tiruan"
)

var (
	ctx   = context.Background()
	admin = inti.Pelaku{AkunID: "UJI-ADMIN"}
)

func layanan(g *tiruan.Gudang) *services.Layanan { return services.BaruLayanan(g, tiruan.Transaksi) }

func anggota(login string) models.Isian {
	return models.Isian{AksesLogin: login, LeaderID: "10000101", BranchParent: "UJI-P00", BranchDetailID: "UJI-B01", Aktif: true}
}

// Anggota baru: Marketing Code = CONTACT_ID akun, nama dari M_LOGIN_GO, leader dan cabang disalin seperti
// SaveMarketingOfficer_Act, ID = 1 + 7 digit sequence, pelaku di USERUPDATE.
func TestTambahAnggotaBaru(t *testing.T) {
	g := tiruan.Contoh()
	m, err := layanan(g).Tambah(ctx, admin, anggota(" UJI-MKT01 "))
	if err != nil {
		t.Fatal(err)
	}
	mau := models.MarketingOfficer{ID: "10000201", ClientID: "CON-1002", ClientName: "UJI Marketing Satu",
		ClientID2: "10000101", MOLeader: "UJI Leader Satu", MOStatus: "1", BranchParent: "UJI-P00",
		BranchDetailID: "UJI-B01", BranchDetailName: "UJI Cabang Satu", TeamGroup: "1", AksesLogin: "UJI-MKT01",
		UserUpdate: "UJI-ADMIN", Tanggal: "2026-10-03 09:00"}
	if m != mau {
		t.Errorf("tambah anggota:\n dapat %+v\n mau   %+v", m, mau)
	}
}

// Leader baru: CLIENTID2 = LEADER, MOLEADER = nama sendiri ("Jika leader").
func TestTambahLeaderBaru(t *testing.T) {
	g := tiruan.Contoh()
	isi := anggota("UJI-MKT01")
	isi.Leader, isi.LeaderID = true, ""
	m, err := layanan(g).Tambah(ctx, admin, isi)
	if err != nil {
		t.Fatal(err)
	}
	if m.ClientID2 != models.NilaiLeader || m.MOLeader != "UJI Marketing Satu" {
		t.Errorf("leader: CLIENTID2 %q MOLEADER %q", m.ClientID2, m.MOLeader)
	}
}

// Akun yang sudah dipakai baris MO lama: Marketing Code orang itu dipakai lagi, bukan CONTACT_ID.
func TestTambahMemakaiMarketingCodeLama(t *testing.T) {
	g := tiruan.Contoh()
	g.MO["10000150"] = models.MarketingOfficer{ID: "10000150", ClientID: "UJI-KONTAK-50", ClientName: "UJI Marketing Dua",
		ClientID2: "10000101", MOStatus: models.StatusNonaktif, AksesLogin: "uji-mkt02"}
	m, err := layanan(g).Tambah(ctx, admin, anggota("UJI-MKT02"))
	if err != nil {
		t.Fatal(err)
	}
	if m.ClientID != "UJI-KONTAK-50" {
		t.Errorf("Marketing Code %q, mau kode lama UJI-KONTAK-50", m.ClientID)
	}
}

func TestTambahDitolakTanpaTulisan(t *testing.T) {
	for _, k := range []struct {
		nama  string
		ubah  func(*tiruan.Gudang, *models.Isian)
		mau   error
		pesan string
	}{
		{"tanpa akun", func(_ *tiruan.Gudang, i *models.Isian) { i.AksesLogin = " " }, services.ErrMasukanTidakSah, "Login Account is required"},
		{"akun tidak ada", func(_ *tiruan.Gudang, i *models.Isian) { i.AksesLogin = "UJI-TIDAKADA" }, services.ErrMasukanTidakSah, "not in the user list"},
		{"akun nonaktif", func(_ *tiruan.Gudang, i *models.Isian) { i.AksesLogin = "UJI-OFF" }, services.ErrMasukanTidakSah, "is not active"},
		{"akun lebih dari 50", func(g *tiruan.Gudang, i *models.Isian) {
			id := "UJI-" + strings.Repeat("P", 47)
			g.Akun[id] = models.Akun{LoginID: id, Nama: "UJI Panjang", ContactID: "CON-1009", Aktif: true}
			i.AksesLogin = id
		}, services.ErrMasukanTidakSah, "longer than 50"},
		{"nama akun lebih dari 100 byte", func(g *tiruan.Gudang, i *models.Isian) {
			g.Akun["UJI-NAMA"] = models.Akun{LoginID: "UJI-NAMA", Nama: strings.Repeat("N", 101), ContactID: "CON-1010", Aktif: true}
			i.AksesLogin = "UJI-NAMA"
		}, services.ErrMasukanTidakSah, "fix it in Kelola User"},
		{"leader kosong", func(_ *tiruan.Gudang, i *models.Isian) { i.LeaderID = "" }, services.ErrMasukanTidakSah, "Leader is required"},
		{"leader nonaktif", func(_ *tiruan.Gudang, i *models.Isian) { i.LeaderID = "10000102" }, services.ErrMasukanTidakSah, "not an active leader"},
		{"leader bukan leader", func(_ *tiruan.Gudang, i *models.Isian) { i.LeaderID = "10000103" }, services.ErrMasukanTidakSah, "not an active leader"},
		{"sub branch nonaktif", func(_ *tiruan.Gudang, i *models.Isian) { i.BranchDetailID = "UJI-B09" }, services.ErrMasukanTidakSah, "is not active"},
		{"sub branch bukan milik branch", func(_ *tiruan.Gudang, i *models.Isian) { i.BranchParent = "UJI-B01" }, services.ErrMasukanTidakSah, "does not belong"},
		{"orang sudah aktif", func(g *tiruan.Gudang, i *models.Isian) {
			g.MO["10000160"] = models.MarketingOfficer{ID: "10000160", ClientID: "CON-1002", MOStatus: models.StatusAktif}
		}, services.ErrSudahAktif, "already has an active marketing officer row (ID 10000160)"},
		// Akun yang sudah dipakai baris AKTIF lain = orang yang sama (Marketing Code-nya dipakai lagi), jadi yang
		// menolak adalah satu-aktif per Marketing Code.
		{"akun dipakai baris aktif lain", func(g *tiruan.Gudang, i *models.Isian) {
			g.MO["10000161"] = models.MarketingOfficer{ID: "10000161", ClientID: "UJI-KONTAK-61", MOStatus: models.StatusAktif, AksesLogin: "uji-mkt01"}
		}, services.ErrSudahAktif, "already has an active marketing officer row (ID 10000161)"},
	} {
		g := tiruan.Contoh()
		isi := anggota("UJI-MKT01")
		k.ubah(g, &isi)
		_, err := layanan(g).Tambah(ctx, admin, isi)
		if !errors.Is(err, k.mau) || err == nil || !strings.Contains(err.Error(), k.pesan) {
			t.Errorf("%s: %v, mau %v berisi %q", k.nama, err, k.mau, k.pesan)
		}
		if len(g.Disisip) != 0 {
			t.Errorf("%s: tetap menyisipkan %v", k.nama, g.Disisip)
		}
	}
	if _, err := layanan(tiruan.Contoh()).Tambah(ctx, inti.Pelaku{}, anggota("UJI-MKT01")); !errors.Is(err, services.ErrTanpaPelaku) {
		t.Errorf("tanpa pelaku: %v", err)
	}
}

// Nonaktif tidak tunduk aturan satu-aktif: riwayat orang yang sama boleh banyak baris nonaktif.
func TestTambahNonaktifBolehWalauOrangSudahAktif(t *testing.T) {
	g := tiruan.Contoh()
	isi := anggota("UJI-LEAD01")
	isi.Aktif = false
	m, err := layanan(g).Tambah(ctx, admin, isi)
	if err != nil {
		t.Fatal(err)
	}
	if m.MOStatus != models.StatusNonaktif || m.ClientID != "CON-1001" {
		t.Errorf("nonaktif: %+v", m)
	}
}

// Ubah: Marketing Code dan nama TIDAK berubah walau akun diganti; AKSES_LOGIN lama Pega yang tidak diganti tidak
// diperiksa ulang; leader dan cabang yang tidak diganti tidak diperiksa ulang.
func TestUbahMenjagaIdentitasDanNilaiLama(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	lama := g.MO["10000103"]
	isi := models.Isian{AksesLogin: "UJIOPERATORLAMA", LeaderID: "10000101", BranchParent: "UJI-P00", BranchDetailID: "UJI-B02", Aktif: true}
	// Leader dinonaktifkan sesudah baris ini dibuat - ubah lain tetap boleh.
	ld := g.MO["10000101"]
	ld.MOStatus = models.StatusNonaktif
	g.MO["10000101"] = ld
	m, err := l.Ubah(ctx, admin, "10000103", isi)
	if err != nil {
		t.Fatal(err)
	}
	if m.AksesLogin != "UJIOPERATORLAMA" || m.BranchDetailName != "UJI Cabang Dua Lama" || m.ClientID2 != "10000101" {
		t.Errorf("nilai lama tidak dipertahankan: %+v", m)
	}
	// Ganti akun: Marketing Code dan nama tetap.
	isi.AksesLogin = "UJI-MKT02"
	m, err = l.Ubah(ctx, admin, "10000103", isi)
	if err != nil {
		t.Fatal(err)
	}
	if m.AksesLogin != "UJI-MKT02" || m.ClientID != lama.ClientID || m.ClientName != lama.ClientName || m.UserUpdate != "UJI-ADMIN" {
		t.Errorf("ganti akun: %+v", m)
	}
	// Ganti cabang: salinan cabang baru.
	isi.BranchDetailID = "UJI-B01"
	if m, err = l.Ubah(ctx, admin, "10000103", isi); err != nil || m.BranchDetailName != "UJI Cabang Satu" || m.TeamGroup != "1" {
		t.Errorf("ganti cabang: %+v %v", m, err)
	}
}

func TestUbahDitolak(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Ubah(ctx, admin, "19999999", anggota("UJI-MKT01")); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("tidak ada: %v", err)
	}
	isi := models.Isian{LeaderID: "10000101", Aktif: true}
	if _, err := l.Ubah(ctx, admin, "10000101", isi); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("leader diri sendiri: %v", err)
	}
	// Mengaktifkan baris lama orang yang sudah punya baris aktif.
	g.MO["10000170"] = models.MarketingOfficer{ID: "10000170", ClientID: "CON-1001", ClientName: "UJI Leader Satu",
		ClientID2: models.NilaiLeader, MOLeader: "UJI Leader Satu", MOStatus: models.StatusNonaktif}
	if _, err := l.Ubah(ctx, admin, "10000170", models.Isian{Leader: true, Aktif: true}); !errors.Is(err, services.ErrSudahAktif) {
		t.Errorf("aktif ganda: %v", err)
	}
	// Mengganti akun ke akun yang dipakai baris aktif orang lain.
	isiLama := models.Isian{AksesLogin: "UJI-LEAD01", LeaderID: "10000101", BranchParent: "UJI-P00", BranchDetailID: "UJI-B02", Aktif: true}
	if _, err := l.Ubah(ctx, admin, "10000103", isiLama); !errors.Is(err, services.ErrSudahAktif) ||
		!strings.Contains(err.Error(), "Login Account UJI-LEAD01 is already used by active marketing officer ID 10000101") {
		t.Errorf("akun ganda: %v", err)
	}
	if len(g.Diperbarui) != 0 {
		t.Errorf("tetap memperbarui: %v", g.Diperbarui)
	}
	// Menonaktifkan selalu boleh.
	if m, err := l.Ubah(ctx, admin, "10000101", models.Isian{Leader: true, AksesLogin: "UJI-LEAD01", BranchParent: "UJI-P00",
		BranchDetailID: "UJI-B01"}); err != nil || m.MOStatus != models.StatusNonaktif {
		t.Errorf("nonaktifkan: %+v %v", m, err)
	}
}

func TestDaftarMenandaiAkun(t *testing.T) {
	g := tiruan.Contoh()
	g.MO["10000180"] = models.MarketingOfficer{ID: "10000180", ClientName: "UJI Akun Nonaktif", MOStatus: models.StatusAktif, AksesLogin: "uji-off"}
	g.MO["10000181"] = models.MarketingOfficer{ID: "10000181", ClientName: "UJI Tanpa Akun", MOStatus: models.StatusAktif}
	d, err := layanan(g).Daftar(ctx)
	if err != nil {
		t.Fatal(err)
	}
	status := map[string]string{}
	for _, b := range d {
		status[b.ID] = b.StatusAkun + "|" + b.EmailAkun
	}
	for id, mau := range map[string]string{
		"10000101": "aktif|uji.lead01@nusantara.example",
		"10000103": "tidak-ada|",
		"10000180": "nonaktif|",
		"10000181": "|",
	} {
		if status[id] != mau {
			t.Errorf("%s: %q, mau %q", id, status[id], mau)
		}
	}
}

func TestPilihanHanyaYangAktif(t *testing.T) {
	p, err := layanan(tiruan.Contoh()).Pilihan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var akun []string
	for _, a := range p.Akun {
		akun = append(akun, a.LoginID)
	}
	if strings.Join(akun, ",") != "UJI-LEAD01,UJI-MKT02,UJI-MKT01" {
		t.Errorf("akun aktif urut nama: %v", akun)
	}
	if len(p.Leader) != 1 || p.Leader[0].ID != "10000101" {
		t.Errorf("leader aktif: %+v", p.Leader)
	}
	if len(p.Branch) != 1 || p.Branch[0].ID != "UJI-P00" || len(p.SubBranch) != 3 {
		t.Errorf("branch %+v sub branch %d", p.Branch, len(p.SubBranch))
	}
}
