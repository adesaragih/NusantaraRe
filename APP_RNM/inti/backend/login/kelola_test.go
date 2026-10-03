package login

// Aturan Kelola User - TANPA Oracle. Penjaga "admin terakhir" ditiru gudang
// tiruan seperti gudang Oracle menegakkannya: sesudah perubahan, di dalam
// transaksinya.

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"

	"nusantarare/inti/backend/menu"
)

var _ GudangKelola = (*gudangTiruan)(nil)

// kodeAdmin - pemegang Kelola User yang AKTIF, urut.
func (g *gudangTiruan) kodeAdmin() []string {
	var out []string
	for id, a := range g.akun {
		if a.Aktif && menuDipegang(g.menu[id], menu.KodeKelolaUser) {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

func menuDipegang(daftar []string, kode string) bool {
	for _, k := range daftar {
		if k == kode {
			return true
		}
	}
	return false
}

// jagaAdmin menjalankan perubahan lalu membatalkannya bila tidak tersisa
// seorang admin aktif pun - padanan transaksi gudang Oracle.
func (g *gudangTiruan) jagaAdmin(ubah func()) error {
	salinAkun := map[string]Akun{}
	for id, a := range g.akun {
		salinAkun[id] = *a
	}
	salinWB, salinMenu := map[string][]string{}, map[string][]string{}
	for id, w := range g.workbasket {
		salinWB[id] = w
	}
	for id, m := range g.menu {
		salinMenu[id] = m
	}
	ubah()
	if len(g.kodeAdmin()) > 0 {
		return nil
	}
	g.akun = map[string]*Akun{}
	for id, a := range salinAkun {
		a := a
		g.akun[id] = &a
	}
	g.workbasket, g.menu = salinWB, salinMenu
	return ErrAdminTerakhir
}

func (g *gudangTiruan) DaftarAkun(context.Context) ([]RingkasAkun, error) {
	var out []RingkasAkun
	for _, a := range g.akun {
		r := ringkasDari(*a)
		r.Kontak = g.kontak[a.ID]
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].AkunID < out[j].AkunID })
	return out, nil
}
func (g *gudangTiruan) RingkasAkun(_ context.Context, id string) (RingkasAkun, error) {
	a, ada := g.akun[id]
	if !ada {
		return RingkasAkun{}, ErrAkunTidakAda
	}
	r := ringkasDari(*a)
	r.Kontak = g.kontak[id]
	return r, nil
}
func ringkasDari(a Akun) RingkasAkun {
	return RingkasAkun{AkunID: a.ID, Nama: a.Nama, Organisasi: a.Organisasi, Divisi: a.Divisi, Unit: a.Unit,
		Aktif: a.Aktif, Terkunci: a.Terkunci, WajibGantiSandi: a.WajibGantiSandi}
}
func (g *gudangTiruan) UbahAkun(_ context.Context, id string, isi IsianAkun) error {
	return g.jagaAdmin(func() {
		a := g.akun[id]
		a.Nama, a.Organisasi, a.Divisi, a.Unit = isi.Nama, isi.Organisasi, isi.Divisi, isi.Unit
		g.tulisKontak(id, isi.Kontak)
		g.workbasket[id], g.menu[id] = isi.Workbasket, isi.Menu
	})
}
func (g *gudangTiruan) SetelAktif(_ context.Context, id string, aktif bool) error {
	return g.jagaAdmin(func() {
		g.akun[id].Aktif = aktif
		g.akun[id].VersiSesi++
	})
}
func (g *gudangTiruan) AturSandi(_ context.Context, id, hash string, wajibGanti bool) error {
	a, ada := g.akun[id]
	if !ada {
		return ErrAkunTidakAda
	}
	a.WajibGantiSandi = wajibGanti
	if hash != "" {
		a.HashSandi, a.Terkunci = hash, false
		a.VersiSesi++
	}
	return nil
}
func (g *gudangTiruan) BukaKunci(_ context.Context, id string) error {
	g.akun[id].Terkunci = false
	return nil
}
func (g *gudangTiruan) HapusAkun(_ context.Context, id string) error {
	return g.jagaAdmin(func() {
		delete(g.akun, id)
		delete(g.workbasket, id)
		delete(g.menu, id)
	})
}
func (g *gudangTiruan) Master(context.Context) (PilihanKelola, error) {
	return PilihanKelola{
		Organisasi: []OpsiMaster{{Kode: "RNM", Nama: "Nusantara Re"}},
		Divisi:     []OpsiMaster{{Kode: "TECH", Nama: "Teknik", Induk: "RNM"}, {Kode: "FIN", Nama: "Keuangan", Induk: "RNM"}},
		Unit:       []OpsiMaster{{Kode: "CLM", Nama: "Klaim", Induk: "TECH"}, {Kode: "TAX", Nama: "Pajak", Induk: "FIN"}},
		Workbasket: []OpsiMaster{{Kode: "ReasLifeAdmin", Nama: "Admin"}, {Kode: "ReasLifeSPV", Nama: "SPV"}},
	}, nil
}

// pembacaMenuUji - baris modul tabel menu tanpa Oracle.
type pembacaMenuUji struct{ baris []menu.Baris }

func (p pembacaMenuUji) Baca(context.Context) ([]menu.Baris, error) { return p.baris, nil }

func kelolaUji(g *gudangTiruan) *Kelola {
	return NewKelola(g, pembacaMenuUji{baris: []menu.Baris{
		{ID: 2, Kode: "claimlife", Label: "Claim Life", Golongan: "KLAIM", Modul: "claimlife", Urutan: 2, Dimigrasi: true},
		{ID: 1, Kode: "claimfacin", Label: "Claim Fac In", Golongan: "KLAIM", Modul: "claimfacin", Urutan: 1},
		{ID: 3, Kode: "premiumlistlife", Label: "PremiumList Life", Golongan: "TREATY", Modul: "premiumlistlife", Urutan: 5, Dimigrasi: true},
	}})
}

var ctxUji = context.Background()

func TestKelolaBuatAkunDenganSandiAdminWajibGanti(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	baru := AkunBaru{ID: " UJI-KELOLA-1 ", Nama: " Uji Kelola ", Organisasi: "RNM", Divisi: "TECH", Unit: "CLM",
		Workbasket: []string{"ReasLifeSPV", "ReasLifeSPV", ""}, Menu: []string{"claimlife", "claimlife"}}
	if err := k.Buat(ctxUji, "UJI-ADMIN", baru, "Sandi-Admin-01", true); err != nil {
		t.Fatal(err)
	}
	a := g.akun["UJI-KELOLA-1"]
	// ⛔ Sandi yang diketik admin bukan sandi pemiliknya: WAJIB diganti saat
	// login pertama (keputusan work owner).
	if a == nil || !a.WajibGantiSandi || a.Nama != "Uji Kelola" || !CocokSandi(a.HashSandi, "Sandi-Admin-01") {
		t.Fatalf("akun %+v", a)
	}
	if !reflect.DeepEqual(g.workbasket["UJI-KELOLA-1"], []string{"ReasLifeSPV"}) ||
		!reflect.DeepEqual(g.menu["UJI-KELOLA-1"], []string{"claimlife"}) {
		t.Errorf("ganda tidak dibuang: wb %v menu %v", g.workbasket["UJI-KELOLA-1"], g.menu["UJI-KELOLA-1"])
	}
}

func TestKelolaBuatDitolak(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	sah := AkunBaru{ID: "UJI-KELOLA-2", Nama: "Uji", Organisasi: "RNM", Menu: []string{"claimlife"}}
	for _, kasus := range []struct {
		nama  string
		ubah  func(*AkunBaru)
		sandi string
		mau   error
	}{
		{"sandi pendek", func(*AkunBaru) {}, "pendek", ErrSandiTerlaluPendek},
		{"menu tak dikenal", func(a *AkunBaru) { a.Menu = []string{"tidakada"} }, "Sandi-Admin-01", ErrMenuTidakDikenal},
		{"akun ganda", func(a *AkunBaru) { a.ID = "UJI-ADMIN" }, "Sandi-Admin-01", ErrAkunSudahAda},
		{"jenjang", func(a *AkunBaru) { a.Divisi, a.Unit = "TECH", "TAX" }, "Sandi-Admin-01", ErrJenjangTidakCocok},
		{"workbasket tak ada", func(a *AkunBaru) { a.Workbasket = []string{"Hantu"} }, "Sandi-Admin-01", ErrMasterTidakAda},
	} {
		a := sah
		kasus.ubah(&a)
		if err := k.Buat(ctxUji, "UJI-ADMIN", a, kasus.sandi, true); !errors.Is(err, kasus.mau) {
			t.Errorf("%s: %v, mau %v", kasus.nama, err, kasus.mau)
		}
	}
	if _, ada := g.akun["UJI-KELOLA-2"]; ada {
		t.Error("akun tersimpan walau ditolak")
	}
}

func TestKelolaUbahMenggantiProfilWorkbasketMenu(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	isi := IsianAkun{Nama: "Uji Kunci Baru", Organisasi: "RNM", Divisi: "FIN", Unit: "TAX",
		Workbasket: []string{"ReasLifeAdmin"}, Menu: []string{"premiumlistlife", menu.KodeKelolaUser}}
	r, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-KUNCI", isi)
	if err != nil {
		t.Fatal(err)
	}
	if r.Nama != "Uji Kunci Baru" || r.Unit != "TAX" || !reflect.DeepEqual(r.Workbasket, []string{"ReasLifeAdmin"}) ||
		!reflect.DeepEqual(r.Menu, []string{menu.KodeKelolaUser, "premiumlistlife"}) {
		t.Errorf("rinci sesudah ubah %+v", r)
	}
	if _, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-TIDAK-ADA", isi); !errors.Is(err, ErrAkunTidakAda) {
		t.Errorf("akun tak ada: %v", err)
	}
	if _, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-KUNCI", IsianAkun{Nama: " "}); !errors.Is(err, ErrAkunTidakSah) {
		t.Errorf("nama kosong: %v", err)
	}
}

// ⛔ Penjaga diri sendiri (keputusan work owner): admin tidak dapat mencabut
// Kelola User dari dirinya, menonaktifkan, atau menghapus dirinya - yang lain
// tetap boleh.
func TestKelolaPenjagaDiriSendiri(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	if _, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-ADMIN", IsianAkun{Nama: "Uji Admin", Menu: []string{"claimlife"}}); !errors.Is(err, ErrDiriSendiri) {
		t.Errorf("cabut Kelola User dari diri sendiri: %v", err)
	}
	if _, err := k.SetelAktif(ctxUji, "UJI-ADMIN", "UJI-ADMIN", false); !errors.Is(err, ErrDiriSendiri) {
		t.Errorf("nonaktifkan diri sendiri: %v", err)
	}
	if err := k.Hapus(ctxUji, "UJI-ADMIN", "UJI-ADMIN"); !errors.Is(err, ErrDiriSendiri) {
		t.Errorf("hapus diri sendiri: %v", err)
	}
	if g.akun["UJI-ADMIN"] == nil || !g.akun["UJI-ADMIN"].Aktif || !menuDipegang(g.menu["UJI-ADMIN"], menu.KodeKelolaUser) {
		t.Fatal("admin berubah walau ditolak")
	}
	// Mengubah menu LAIN milik diri sendiri tetap boleh, selama Kelola User tinggal.
	if _, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-ADMIN", IsianAkun{Nama: "Uji Admin", Menu: []string{menu.KodeKelolaUser}}); err != nil {
		t.Errorf("ubah diri sendiri dengan Kelola User tetap: %v", err)
	}
}

// ⛔ Admin terakhir: perubahan yang menyisakan NOL akun aktif ber-Kelola User
// ditolak - juga bila pelakunya bukan admin itu (mis. admin kedua yang sedang
// dicabut sesinya).
func TestKelolaAdminTerakhir(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	// UJI-KUNCI dijadikan admin kedua, lalu UJI-ADMIN dinonaktifkan olehnya.
	g.menu["UJI-KUNCI"] = []string{menu.KodeKelolaUser}
	if _, err := k.SetelAktif(ctxUji, "UJI-KUNCI", "UJI-ADMIN", false); err != nil {
		t.Fatalf("admin kedua menonaktifkan admin pertama: %v", err)
	}
	if v := g.akun["UJI-ADMIN"].VersiSesi; v != 4 {
		t.Errorf("nonaktif tidak mencabut sesi: versi %d", v)
	}
	// Kini UJI-KUNCI satu-satunya admin aktif: menghapusnya - oleh siapa pun
	// selain dirinya - ditolak.
	if err := k.Hapus(ctxUji, "UJI-ADMIN", "UJI-KUNCI"); !errors.Is(err, ErrAdminTerakhir) {
		t.Errorf("hapus admin terakhir: %v", err)
	}
	if g.akun["UJI-KUNCI"] == nil {
		t.Fatal("admin terakhir terhapus")
	}
	if _, err := k.SetelAktif(ctxUji, "UJI-KUNCI", "UJI-ADMIN", true); err != nil {
		t.Errorf("aktifkan kembali: %v", err)
	}
	if err := k.Hapus(ctxUji, "UJI-ADMIN", "UJI-KUNCI"); err != nil {
		t.Errorf("hapus admin kedua saat admin lain aktif: %v", err)
	}
	if _, ada := g.akun["UJI-KUNCI"]; ada || g.menu["UJI-KUNCI"] != nil {
		t.Error("hapus permanen menyisakan akun atau menunya")
	}
}

func TestKelolaBukaKunciDanDaftar(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	r, err := k.BukaKunci(ctxUji, "UJI-ADMIN", "UJI-KUNCI")
	if err != nil || r.Terkunci {
		t.Errorf("buka kunci: %+v %v", r, err)
	}
	if _, err := k.BukaKunci(ctxUji, "UJI-ADMIN", "UJI-TIDAK-ADA"); !errors.Is(err, ErrAkunTidakAda) {
		t.Errorf("buka kunci akun tak ada: %v", err)
	}
	daftar, err := k.Daftar(ctxUji)
	if err != nil || len(daftar) != 4 || daftar[0].AkunID != "UJI-ADMIN" {
		t.Errorf("daftar %+v %v", daftar, err)
	}
}

// Pilihan form: master dari gudang, menu = baris modul tabel menu per
// golongan (urutan sidebar, lalu URUTAN), ditambah menu aplikasi di ADMIN.
func TestKelolaPilihanMenuUrutSidebar(t *testing.T) {
	p, err := kelolaUji(gudangUji(t)).Pilihan(ctxUji)
	if err != nil {
		t.Fatal(err)
	}
	var dapat []string
	for _, m := range p.Menu {
		dapat = append(dapat, m.Golongan+"/"+m.Kode)
	}
	mau := []string{"TREATY/premiumlistlife", "KLAIM/claimfacin", "KLAIM/claimlife", menu.GolonganAdmin + "/" + menu.KodeKelolaUser}
	if !reflect.DeepEqual(dapat, mau) {
		t.Errorf("menu pilihan %v, mau %v", dapat, mau)
	}
	if len(p.Workbasket) != 2 || len(p.Unit) != 2 || p.Menu[1].Dimigrasi {
		t.Errorf("pilihan %+v", p)
	}
}

// Ubah TIDAK membuang workbasket yang sedang nonaktif di master - layar tidak
// menampilkannya, jadi tidak pernah mengirimnya - dan tidak menolak akun yang
// organisasi/divisi/unit-nya kini nonaktif selama jenjangnya TIDAK diubah.
// Workbasket yang BARU ditambahkan tetap harus aktif (temuan /code-review).
func TestKelolaUbahMenjagaYangNonaktif(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	g.masterWorkbasket["ReasLifeSPV"] = false
	a := g.akun["UJI-KUNCI"]
	a.Organisasi, a.Divisi, a.Unit = "RNM", "TECH", "CLM"
	g.tidakAktif["CLM"] = true
	g.workbasket["UJI-KUNCI"] = []string{"ReasLifeAdmin", "ReasLifeSPV"}

	r, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-KUNCI", IsianAkun{Nama: "Uji Kunci", Organisasi: "RNM", Divisi: "TECH",
		Unit: "CLM", Workbasket: []string{"ReasLifeAdmin"}, Menu: []string{"claimlife"}})
	if err != nil {
		t.Fatalf("ubah akun berunit nonaktif tanpa mengubah jenjang: %v", err)
	}
	if !reflect.DeepEqual(g.workbasket["UJI-KUNCI"], []string{"ReasLifeAdmin", "ReasLifeSPV"}) {
		t.Errorf("workbasket nonaktif terbuang: %v", g.workbasket["UJI-KUNCI"])
	}
	// Layar tetap menampilkan yang aktif saja.
	if !reflect.DeepEqual(r.Workbasket, []string{"ReasLifeAdmin"}) {
		t.Errorf("rinci workbasket %v", r.Workbasket)
	}
	// Jenjang yang DIUBAH ke unit nonaktif ditolak; menambah workbasket nonaktif ditolak.
	g.tidakAktif["TAX"] = true
	if _, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-KUNCI", IsianAkun{Nama: "Uji Kunci", Organisasi: "RNM", Divisi: "FIN",
		Unit: "TAX"}); !errors.Is(err, ErrMasterTidakAda) {
		t.Errorf("pindah ke unit nonaktif: %v", err)
	}
	g.workbasket["UJI-ADMIN"] = []string{"ReasLifeAdmin"}
	if _, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-ADMIN", IsianAkun{Nama: "Uji Admin", Workbasket: []string{"ReasLifeAdmin", "ReasLifeSPV"},
		Menu: []string{menu.KodeKelolaUser}}); !errors.Is(err, ErrMasterTidakAda) {
		t.Errorf("menambah workbasket nonaktif: %v", err)
	}
}

func TestKelolaBuatAkunTitikDitolak(t *testing.T) {
	g := gudangUji(t)
	for _, id := range []string{".", ".."} {
		if err := kelolaUji(g).Buat(ctxUji, "UJI-ADMIN", AkunBaru{ID: id, Nama: "Titik"}, "Sandi-Admin-01", true); !errors.Is(err, ErrAkunTidakSah) {
			t.Errorf("akun %q: %v", id, err)
		}
	}
}

// Tab Security (permintaan work owner 01-10-2026): "Change Password Next
// Login" dicentang = wajib ganti, tidak dicentang = tidak perlu - saat
// membuat akun juga.
func TestKelolaBuatTanpaWajibGanti(t *testing.T) {
	g := gudangUji(t)
	if err := kelolaUji(g).Buat(ctxUji, "UJI-ADMIN", AkunBaru{ID: "UJI-BEBAS", Nama: "Bebas"}, "Sandi-Admin-01", false); err != nil {
		t.Fatal(err)
	}
	if a := g.akun["UJI-BEBAS"]; a == nil || a.WajibGantiSandi {
		t.Errorf("akun tanpa centang wajib ganti: %+v", a)
	}
}

// Atur password akun lain: hash baru, centang wajib ganti mengikuti isian,
// seluruh sesinya dicabut (versi naik), kuncinya dibuka. Tanpa password =
// hanya centangnya yang berubah - hash dan sesi tetap.
func TestKelolaAturSandi(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	hashLama, versiLama := g.akun["UJI-KUNCI"].HashSandi, g.akun["UJI-KUNCI"].VersiSesi
	r, err := k.AturSandi(ctxUji, "UJI-ADMIN", "UJI-KUNCI", "Sandi-Baru-Admin-1", true)
	a := g.akun["UJI-KUNCI"]
	if err != nil || !r.WajibGantiSandi || r.Terkunci || a.VersiSesi != versiLama+1 || !CocokSandi(a.HashSandi, "Sandi-Baru-Admin-1") {
		t.Fatalf("atur sandi: %+v %+v %v", r, a, err)
	}
	hashBaru, versiBaru := a.HashSandi, a.VersiSesi
	if r, err := k.AturSandi(ctxUji, "UJI-ADMIN", "UJI-KUNCI", "", false); err != nil || r.WajibGantiSandi ||
		a.HashSandi != hashBaru || a.VersiSesi != versiBaru {
		t.Errorf("hanya centang: %+v %+v %v", r, a, err)
	}
	if hashBaru == hashLama {
		t.Error("hash tidak berubah")
	}
	for _, kasus := range []struct {
		id, sandi string
		mau       error
	}{
		{"UJI-KUNCI", "pendek", ErrSandiTerlaluPendek},
		{"UJI-TIDAK-ADA", "Sandi-Baru-Admin-1", ErrAkunTidakAda},
	} {
		if _, err := k.AturSandi(ctxUji, "UJI-ADMIN", kasus.id, kasus.sandi, false); !errors.Is(err, kasus.mau) {
			t.Errorf("%s %q: %v, mau %v", kasus.id, kasus.sandi, err, kasus.mau)
		}
	}
}
