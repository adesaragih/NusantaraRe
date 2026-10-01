package login

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// gudangTiruan - M_LOGIN_GO di memori, cukup untuk aturan layanan.
type gudangTiruan struct {
	akun       map[string]*Akun
	workbasket map[string][]string
	// master organisasi: unit -> divisi, divisi -> organisasi; aktif.
	unit, divisi     map[string]string
	organisasi       map[string]bool
	tidakAktif       map[string]bool
	masterWorkbasket map[string]bool
	gagal, berhasil  []string
	dibuat           []AkunBaru
}

// hashUji - satu hash bcrypt untuk seluruh uji (cost 12 lambat).
var hashUji = sync.OnceValues(func() (string, error) { return HashSandi("Sandi-Benar-01") })

func gudangUji(t *testing.T) *gudangTiruan {
	t.Helper()
	h, err := hashUji()
	if err != nil {
		t.Fatal(err)
	}
	return &gudangTiruan{
		akun: map[string]*Akun{
			"UJI-ADMIN":    {ID: "UJI-ADMIN", Nama: "Uji Admin", HashSandi: h, Aktif: true, VersiSesi: 3},
			"UJI-NONAKTIF": {ID: "UJI-NONAKTIF", Nama: "Uji Nonaktif", HashSandi: h, Aktif: false, VersiSesi: 1},
			"UJI-KUNCI":    {ID: "UJI-KUNCI", Nama: "Uji Kunci", HashSandi: h, Aktif: true, Terkunci: true, VersiSesi: 1},
			"UJI-BARU":     {ID: "UJI-BARU", Nama: "Uji Baru", HashSandi: h, Aktif: true, WajibGantiSandi: true, VersiSesi: 1},
		},
		workbasket:       map[string][]string{"UJI-ADMIN": {"ReasLifeAdmin", "ReasLifeSPV"}},
		unit:             map[string]string{"CLM": "TECH", "TAX": "FIN"},
		divisi:           map[string]string{"TECH": "RNM", "FIN": "RNM"},
		organisasi:       map[string]bool{"RNM": true},
		tidakAktif:       map[string]bool{"UW": true},
		masterWorkbasket: map[string]bool{"ReasLifeAdmin": true, "ReasLifeSPV": true},
	}
}

func (g *gudangTiruan) AmbilAkun(_ context.Context, id string) (Akun, error) {
	a, ada := g.akun[id]
	if !ada {
		return Akun{}, ErrAkunTidakAda
	}
	return *a, nil
}
func (g *gudangTiruan) Workbasket(_ context.Context, id string) ([]string, error) {
	return g.workbasket[id], nil
}
func (g *gudangTiruan) CatatGagal(_ context.Context, id string) error {
	g.gagal = append(g.gagal, id)
	return nil
}
func (g *gudangTiruan) CatatBerhasil(_ context.Context, id string) error {
	g.berhasil = append(g.berhasil, id)
	return nil
}
func (g *gudangTiruan) GantiSandi(_ context.Context, id, hash string, versiLama int64) error {
	a := g.akun[id]
	if a.VersiSesi != versiLama {
		return ErrSesiTidakSah
	}
	a.HashSandi, a.WajibGantiSandi, a.VersiSesi = hash, false, a.VersiSesi+1
	return nil
}
func (g *gudangTiruan) NaikkanVersi(_ context.Context, id string) error {
	g.akun[id].VersiSesi++
	return nil
}
func (g *gudangTiruan) InfoUnit(_ context.Context, code string) (string, bool, error) {
	d, ada := g.unit[code]
	if !ada {
		return "", false, ErrMasterTidakAda
	}
	return d, !g.tidakAktif[code], nil
}
func (g *gudangTiruan) InfoDivisi(_ context.Context, code string) (string, bool, error) {
	o, ada := g.divisi[code]
	if !ada {
		return "", false, ErrMasterTidakAda
	}
	return o, !g.tidakAktif[code], nil
}
func (g *gudangTiruan) InfoOrganisasi(_ context.Context, code string) (bool, error) {
	aktif, ada := g.organisasi[code]
	if !ada {
		return false, ErrMasterTidakAda
	}
	return aktif, nil
}
func (g *gudangTiruan) WorkbasketAktif(_ context.Context, id string) (bool, error) {
	aktif, ada := g.masterWorkbasket[id]
	if !ada {
		return false, ErrMasterTidakAda
	}
	return aktif, nil
}
func (g *gudangTiruan) BuatAkun(_ context.Context, a AkunBaru, hash string, wajibGanti bool) error {
	g.dibuat = append(g.dibuat, a)
	g.akun[a.ID] = &Akun{ID: a.ID, Nama: a.Nama, HashSandi: hash, Aktif: true, WajibGantiSandi: wajibGanti, VersiSesi: 1}
	return nil
}

func layananUji(g *gudangTiruan, saat time.Time) *Layanan {
	return NewLayanan(g, rahasiaUji).DenganJam(func() time.Time { return saat })
}

var saatUji = time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)

func TestMasukBerhasilMenerbitkanTokenDanPeran(t *testing.T) {
	g := gudangUji(t)
	p, tok, err := layananUji(g, saatUji).Masuk(context.Background(), " UJI-ADMIN ", "Sandi-Benar-01")
	if err != nil {
		t.Fatal(err)
	}
	if p.AkunID != "UJI-ADMIN" || p.Nama != "Uji Admin" || len(p.Peran) != 2 || p.WajibGantiSandi {
		t.Errorf("profil %+v", p)
	}
	if tok.Akun != "UJI-ADMIN" || tok.Versi != 3 || !tok.Mulai.Equal(saatUji) {
		t.Errorf("token %+v", tok)
	}
	if len(g.berhasil) != 1 || len(g.gagal) != 0 {
		t.Errorf("catatan berhasil %v gagal %v", g.berhasil, g.gagal)
	}
}

// ⛔ Satu pesan untuk tiga keadaan berbeda: akun tidak ada, nonaktif, sandi
// salah - jawaban tidak boleh membocorkan akun mana yang ada.
func TestMasukGagalSatuPesan(t *testing.T) {
	for _, k := range []struct{ akun, sandi string }{
		{"UJI-TIDAK-ADA", "Sandi-Benar-01"},
		{"UJI-NONAKTIF", "Sandi-Benar-01"},
		{"UJI-ADMIN", "Sandi-Salah-01"},
		{"", "Sandi-Benar-01"},
		{"ada spasi", "Sandi-Benar-01"},
	} {
		g := gudangUji(t)
		_, _, err := layananUji(g, saatUji).Masuk(context.Background(), k.akun, k.sandi)
		if !errors.Is(err, ErrKredensial) {
			t.Errorf("%q/%q: %v, mau ErrKredensial", k.akun, k.sandi, err)
		}
		if len(g.berhasil) != 0 {
			t.Errorf("%q: tercatat berhasil", k.akun)
		}
	}
}

func TestSandiSalahDicatatGagal(t *testing.T) {
	g := gudangUji(t)
	_, _, _ = layananUji(g, saatUji).Masuk(context.Background(), "UJI-ADMIN", "Sandi-Salah-01")
	if len(g.gagal) != 1 || g.gagal[0] != "UJI-ADMIN" {
		t.Errorf("gagal tercatat %v", g.gagal)
	}
}

// Terkunci: sandi BENAR pun ditolak, dan percobaan tidak memperpanjang kunci.
func TestAkunTerkunciDitolakWalauSandiBenar(t *testing.T) {
	g := gudangUji(t)
	_, _, err := layananUji(g, saatUji).Masuk(context.Background(), "UJI-KUNCI", "Sandi-Benar-01")
	if !errors.Is(err, ErrTerkunci) {
		t.Errorf("akun terkunci: %v", err)
	}
	if len(g.gagal) != 0 || len(g.berhasil) != 0 {
		t.Errorf("akun terkunci menyentuh hitungan: gagal %v berhasil %v", g.gagal, g.berhasil)
	}
}

func TestSesiDariCookie(t *testing.T) {
	g := gudangUji(t)
	l := layananUji(g, saatUji)
	_, tok, err := l.Masuk(context.Background(), "UJI-ADMIN", "Sandi-Benar-01")
	if err != nil {
		t.Fatal(err)
	}
	nilai := Tandatangani(rahasiaUji, tok)
	nanti := layananUji(g, saatUji.Add(10*time.Minute))
	p, baru, err := nanti.Sesi(context.Background(), nilai)
	if err != nil || p.AkunID != "UJI-ADMIN" || len(p.Peran) != 2 {
		t.Fatalf("sesi: %+v %v", p, err)
	}
	if !baru.Habis.Equal(saatUji.Add(10*time.Minute + BatasDiam)) {
		t.Errorf("sesi tidak diperpanjang: %v", baru.Habis)
	}
	// Logout menaikkan versi: cookie yang sama tidak berlaku lagi.
	if err := nanti.Keluar(context.Background(), "UJI-ADMIN"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := nanti.Sesi(context.Background(), nilai); !errors.Is(err, ErrSesiTidakSah) {
		t.Errorf("cookie sesudah logout: %v", err)
	}
}

func TestSesiAkunNonaktifDitolak(t *testing.T) {
	g := gudangUji(t)
	l := layananUji(g, saatUji)
	_, tok, _ := l.Masuk(context.Background(), "UJI-ADMIN", "Sandi-Benar-01")
	g.akun["UJI-ADMIN"].Aktif = false
	if _, _, err := l.Sesi(context.Background(), Tandatangani(rahasiaUji, tok)); !errors.Is(err, ErrSesiTidakSah) {
		t.Errorf("akun dinonaktifkan sesudah login: %v", err)
	}
}

func TestGantiSandi(t *testing.T) {
	g := gudangUji(t)
	l := layananUji(g, saatUji)
	p, tok, err := l.Masuk(context.Background(), "UJI-BARU", "Sandi-Benar-01")
	if err != nil || !p.WajibGantiSandi {
		t.Fatalf("akun baru: %+v %v", p, err)
	}
	if _, _, err := l.GantiSandi(context.Background(), tok, "Sandi-Salah-01", "Sandi-Baru-0001"); !errors.Is(err, ErrKredensial) {
		t.Errorf("sandi lama salah: %v", err)
	}
	if _, _, err := l.GantiSandi(context.Background(), tok, "Sandi-Benar-01", "pendek"); !errors.Is(err, ErrSandiTerlaluPendek) {
		t.Errorf("sandi baru pendek: %v", err)
	}
	if _, _, err := l.GantiSandi(context.Background(), tok, "Sandi-Benar-01", "Sandi-Benar-01"); !errors.Is(err, ErrSandiSama) {
		t.Errorf("sandi baru sama: %v", err)
	}
	p, baru, err := l.GantiSandi(context.Background(), tok, "Sandi-Benar-01", "Sandi-Baru-0001")
	if err != nil || p.WajibGantiSandi {
		t.Fatalf("ganti sandi: %+v %v", p, err)
	}
	if baru.Versi != tok.Versi+1 {
		t.Errorf("versi sesudah ganti sandi %d, mau %d", baru.Versi, tok.Versi+1)
	}
	if _, _, err := l.Sesi(context.Background(), Tandatangani(rahasiaUji, tok)); !errors.Is(err, ErrSesiTidakSah) {
		t.Error("cookie lama masih berlaku sesudah ganti sandi")
	}
	if !CocokSandi(g.akun["UJI-BARU"].HashSandi, "Sandi-Baru-0001") {
		t.Error("hash sandi baru tidak tersimpan")
	}
}

func TestBuatPenggunaMemeriksaMasterDanJenjang(t *testing.T) {
	for _, k := range []struct {
		nama string
		a    AkunBaru
		mau  error
	}{
		{"lengkap", AkunBaru{ID: "UJI-B1", Nama: "Uji B1", Organisasi: "RNM", Divisi: "TECH", Unit: "CLM", Workbasket: []string{"ReasLifeAdmin"}}, nil},
		{"tanpa unit", AkunBaru{ID: "UJI-B2", Nama: "Uji B2", Organisasi: "RNM", Divisi: "FIN"}, nil},
		{"tanpa organisasi", AkunBaru{ID: "UJI-B3", Nama: "Uji B3"}, nil},
		{"unit milik divisi lain", AkunBaru{ID: "UJI-B4", Nama: "Uji B4", Organisasi: "RNM", Divisi: "FIN", Unit: "CLM"}, ErrJenjangTidakCocok},
		{"unit tanpa divisi", AkunBaru{ID: "UJI-B5", Nama: "Uji B5", Organisasi: "RNM", Unit: "CLM"}, ErrJenjangTidakCocok},
		{"divisi tanpa organisasi", AkunBaru{ID: "UJI-B6", Nama: "Uji B6", Divisi: "TECH"}, ErrJenjangTidakCocok},
		{"organisasi tidak ada", AkunBaru{ID: "UJI-B7", Nama: "Uji B7", Organisasi: "XXX"}, ErrMasterTidakAda},
		{"unit nonaktif", AkunBaru{ID: "UJI-B8", Nama: "Uji B8", Organisasi: "RNM", Divisi: "TECH", Unit: "UW"}, ErrMasterTidakAda},
		{"workbasket tidak ada", AkunBaru{ID: "UJI-B9", Nama: "Uji B9", Workbasket: []string{"ReasTidakAda"}}, ErrMasterTidakAda},
		{"akun tidak sah", AkunBaru{ID: "ada spasi", Nama: "Uji"}, ErrAkunTidakSah},
		{"nama kosong", AkunBaru{ID: "UJI-B10"}, ErrAkunTidakSah},
		{"akun sudah ada", AkunBaru{ID: "UJI-ADMIN", Nama: "Uji"}, ErrAkunSudahAda},
	} {
		g := gudangUji(t)
		g.unit["UW"] = "TECH"
		sandi, err := NewLayanan(g, nil).BuatPengguna(context.Background(), k.a)
		if !errors.Is(err, k.mau) {
			t.Errorf("%s: %v, mau %v", k.nama, err, k.mau)
			continue
		}
		if k.mau != nil {
			if len(g.dibuat) != 0 {
				t.Errorf("%s: tetap dibuat", k.nama)
			}
			continue
		}
		if PeriksaSandiBaru(sandi) != nil || !CocokSandi(g.akun[k.a.ID].HashSandi, sandi) || !g.akun[k.a.ID].WajibGantiSandi {
			t.Errorf("%s: sandi sementara %q tidak tersimpan sebagai hash, atau tidak wajib diganti", k.nama, sandi)
		}
	}
}

func TestLayananTanpaRahasiaTidakMenerbitkanSesi(t *testing.T) {
	g := gudangUji(t)
	if _, _, err := NewLayanan(g, nil).Masuk(context.Background(), "UJI-ADMIN", "Sandi-Benar-01"); !errors.Is(err, ErrTanpaRahasia) {
		t.Errorf("tanpa SESI_RAHASIA: %v", err)
	}
}

// Sandi DITETAPKAN operator (`-sandi-dari-stdin`): aturan sandi tetap berlaku,
// hash tersimpan, dan akunnya tidak wajib ganti - sandinya bukan sementara.
func TestBuatPenggunaDenganSandi(t *testing.T) {
	g := gudangUji(t)
	l := NewLayanan(g, rahasiaUji)
	if err := l.BuatPenggunaDenganSandi(context.Background(), AkunBaru{ID: "UJI-C1", Nama: "Uji C1"}, "pendek"); !errors.Is(err, ErrSandiTerlaluPendek) {
		t.Errorf("sandi pendek: %v", err)
	}
	if len(g.dibuat) != 0 {
		t.Fatal("akun dibuat walau sandinya ditolak")
	}
	a := AkunBaru{ID: "UJI-C2", Nama: "Uji C2", Organisasi: "RNM", Divisi: "TECH", Unit: "CLM", Workbasket: []string{"ReasLifeAdmin"}}
	if err := l.BuatPenggunaDenganSandi(context.Background(), a, "Sandi-Uji-Tetap-1"); err != nil {
		t.Fatal(err)
	}
	if !CocokSandi(g.akun["UJI-C2"].HashSandi, "Sandi-Uji-Tetap-1") || g.akun["UJI-C2"].WajibGantiSandi {
		t.Errorf("akun bersandi tetap: %+v", g.akun["UJI-C2"])
	}
	p, _, err := l.Masuk(context.Background(), "UJI-C2", "Sandi-Uji-Tetap-1")
	if err != nil || p.WajibGantiSandi {
		t.Errorf("masuk dengan sandi tetap: %+v %v", p, err)
	}
	if err := l.BuatPenggunaDenganSandi(context.Background(), AkunBaru{ID: "UJI-C3", Nama: "Uji", Unit: "CLM"}, "Sandi-Uji-Tetap-1"); !errors.Is(err, ErrJenjangTidakCocok) {
		t.Errorf("jenjang tetap diperiksa: %v", err)
	}
}
