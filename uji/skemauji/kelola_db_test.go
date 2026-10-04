//go:build db

package skemauji_test

// Kelola User di Oracle (keputusan work owner 01-10-2026): isi awal 903,
// buat-ubah-hapus akun beserta workbasket dan menunya, dan penjaga admin
// terakhir di dalam transaksi gudang - SQL-nya dijalankan sungguhan, bukan
// dibaca polanya.

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti/backend/login"
	"nusantarare/inti/backend/menu"
	"nusantarare/inti/backend/migrasi"
	"nusantarare/uji/skemauji"
)

func TestKelolaUserDariOracle(t *testing.T) {
	sqlDB, skema, ctx := pasangSkemaInti(t)
	jalan := pasangMasterTiruan(t, sqlDB, skema, ctx)
	hitung := func(q string) int {
		t.Helper()
		var n int
		if err := sqlDB.QueryRowContext(ctx, q).Scan(&n); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		return n
	}

	// Isi awal 903 atas akun yang SUDAH ADA: 20 modul + Kelola User.
	// CONTACT_ID NOT NULL sejak migrasi 905 - baris tiruan mengisinya sendiri.
	jalan(`INSERT INTO ` + skema + `.M_LOGIN_GO (LOGIN_ID, NAME, PASSWORD_HASH, CONTACT_ID) VALUES ('UJI-LAMA', 'Uji Lama', 'x', 'CON-UJI-LAMA')`)
	langkah, err := migrasi.PernyataanLangkah(os.DirFS("../../inti/backend"), "903_m_login_go_menu.sql")
	if err != nil {
		t.Fatal(err)
	}
	jalan(strings.ReplaceAll(langkah[1], "{skema}", skema))
	if n := hitung(`SELECT COUNT(*) FROM ` + skema + `.M_LOGIN_GO_MENU WHERE LOGIN_ID = 'UJI-LAMA'`); n != 21 {
		t.Fatalf("isi awal 903: %d menu, mau 21", n)
	}

	repo, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = repo.Close() }()
	gudang := login.NewGudangOracle(repo)
	k := login.NewKelola(gudang, menu.NewPembaca(repo))
	l := login.NewLayanan(gudang, []byte("rahasia-uji-yang-panjangnya-32-byte!"))

	if err := k.Buat(ctx, "UJI-LAMA", login.AkunBaru{ID: "UJI-ADM", Nama: "Uji Adm", Organisasi: "UJI-ORG",
		Menu: []string{menu.KodeKelolaUser, "claimlife"}}, "Sandi-Admin-01", true); err != nil {
		t.Fatalf("buat admin: %v", err)
	}
	if err := k.Buat(ctx, "UJI-ADM", login.AkunBaru{ID: "UJI-B", Nama: "Uji B", Organisasi: "UJI-ORG", Divisi: "UJI-DIV",
		Unit: "UJI-UNIT", Workbasket: []string{"ReasLifeSPV"}, Menu: []string{"claimlife"}}, "Sandi-Admin-01", true); err != nil {
		t.Fatalf("buat B: %v", err)
	}
	if err := k.Buat(ctx, "UJI-ADM", login.AkunBaru{ID: "UJI-C", Nama: "Uji C", Menu: []string{"tidakada"}},
		"Sandi-Admin-01", true); !errors.Is(err, login.ErrMenuTidakDikenal) {
		t.Errorf("menu tak dikenal: %v", err)
	}
	p, _, err := l.Masuk(ctx, "UJI-B", "Sandi-Admin-01")
	if err != nil || !p.WajibGantiSandi || !reflect.DeepEqual(p.Menu, []string{"claimlife"}) {
		t.Fatalf("masuk B: %+v %v", p, err)
	}

	// Hapus permanen: akun lama beserta 21 menunya.
	if err := k.Hapus(ctx, "UJI-ADM", "UJI-LAMA"); err != nil {
		t.Fatalf("hapus akun lama: %v", err)
	}
	if n := hitung(`SELECT COUNT(*) FROM `+skema+`.M_LOGIN_GO_MENU WHERE LOGIN_ID = 'UJI-LAMA'`) +
		hitung(`SELECT COUNT(*) FROM `+skema+`.M_LOGIN_GO WHERE LOGIN_ID = 'UJI-LAMA'`); n != 0 {
		t.Errorf("hapus permanen menyisakan %d baris", n)
	}

	// Ubah mengganti workbasket dan menu.
	r, err := k.Ubah(ctx, "UJI-ADM", "UJI-B", login.IsianAkun{Nama: "Uji B Baru", Organisasi: "UJI-ORG",
		Workbasket: []string{"ReasLifeAdmin"}, Menu: []string{"premiumlistlife", menu.KodeKelolaUser}})
	if err != nil || r.Nama != "Uji B Baru" || r.Divisi != "" || !reflect.DeepEqual(r.Workbasket, []string{"ReasLifeAdmin"}) ||
		!reflect.DeepEqual(r.Menu, []string{menu.KodeKelolaUser, "premiumlistlife"}) || r.LoginTerakhir == "" {
		t.Fatalf("ubah B: %+v %v", r, err)
	}

	// Admin terakhir: B dinonaktifkan oleh ADM; lalu menghapus ADM - satu-satunya
	// admin aktif - ditolak, dan transaksinya tidak menyisakan apa pun.
	if _, err := k.SetelAktif(ctx, "UJI-ADM", "UJI-B", false); err != nil {
		t.Fatalf("nonaktifkan B: %v", err)
	}
	if err := k.Hapus(ctx, "UJI-B", "UJI-ADM"); !errors.Is(err, login.ErrAdminTerakhir) {
		t.Errorf("hapus admin terakhir: %v", err)
	}
	if n := hitung(`SELECT COUNT(*) FROM ` + skema + `.M_LOGIN_GO_MENU WHERE LOGIN_ID = 'UJI-ADM'`); n != 2 {
		t.Errorf("ROLLBACK admin terakhir: menu ADM tinggal %d, mau 2", n)
	}
	if _, err := k.Ubah(ctx, "UJI-B", "UJI-ADM", login.IsianAkun{Nama: "Uji Adm", Menu: []string{"claimlife"}}); !errors.Is(err, login.ErrAdminTerakhir) {
		t.Errorf("cabut Kelola User admin terakhir: %v", err)
	}
	if _, err := k.SetelAktif(ctx, "UJI-ADM", "UJI-ADM", false); !errors.Is(err, login.ErrDiriSendiri) {
		t.Errorf("nonaktifkan diri sendiri: %v", err)
	}

	// Buka kunci.
	jalan(`UPDATE ` + skema + `.M_LOGIN_GO SET FAILED_COUNT = 5, LOCKED_UNTIL = SYSDATE + 1 WHERE LOGIN_ID = 'UJI-B'`)
	if r, err := k.Rinci(ctx, "UJI-B"); err != nil || !r.Terkunci || r.Aktif {
		t.Fatalf("B terkunci dan nonaktif: %+v %v", r, err)
	}
	if r, err := k.BukaKunci(ctx, "UJI-ADM", "UJI-B"); err != nil || r.Terkunci {
		t.Errorf("buka kunci: %+v %v", r, err)
	}

	// Tab Security: password baru tanpa wajib ganti - login dengan password itu,
	// tidak wajib ganti; lalu centangnya saja dinyalakan.
	if err := k.Buat(ctx, "UJI-ADM", login.AkunBaru{ID: "UJI-SEC", Nama: "Uji Sec"}, "Sandi-Admin-01", false); err != nil {
		t.Fatalf("buat tanpa wajib ganti: %v", err)
	}
	if r, err := k.AturSandi(ctx, "UJI-ADM", "UJI-SEC", "Sandi-Baru-Oracle-1", false); err != nil || r.WajibGantiSandi {
		t.Fatalf("atur sandi: %+v %v", r, err)
	}
	if p, _, err := l.Masuk(ctx, "UJI-SEC", "Sandi-Baru-Oracle-1"); err != nil || p.WajibGantiSandi {
		t.Errorf("masuk dengan password dari admin: %+v %v", p, err)
	}
	if r, err := k.AturSandi(ctx, "UJI-ADM", "UJI-SEC", "", true); err != nil || !r.WajibGantiSandi {
		t.Errorf("hanya centang: %+v %v", r, err)
	}
	if err := k.Hapus(ctx, "UJI-ADM", "UJI-SEC"); err != nil {
		t.Fatalf("hapus UJI-SEC: %v", err)
	}

	// Daftar dan pilihan form.
	daftar, err := k.Daftar(ctx)
	if err != nil || len(daftar) != 2 || daftar[0].AkunID != "UJI-ADM" || daftar[1].AkunID != "UJI-B" {
		t.Errorf("daftar %+v %v", daftar, err)
	}
	pil, err := k.Pilihan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(pil.Organisasi) != 1 || len(pil.Divisi) != 1 || pil.Divisi[0].Induk != "UJI-ORG" ||
		len(pil.Unit) != 1 || pil.Unit[0].Induk != "UJI-DIV" || len(pil.Workbasket) != 2 || len(pil.Menu) != 21 ||
		pil.Menu[20].Kode != menu.KodeKelolaUser {
		t.Errorf("pilihan %+v", pil)
	}

	// Identitas (migrasi 905, 03-10-2026): CONTACT_ID dari sequence, username dan email unik tanpa beda huruf, index
	// unik sebagai pengaman terakhir, dan login TIDAK lewat email (dibatalkan work owner).
	if len(daftar) == 2 && (!strings.HasPrefix(daftar[0].IDKontak, login.AwalanIDKontak) ||
		!strings.HasPrefix(daftar[1].IDKontak, login.AwalanIDKontak) || daftar[0].IDKontak == daftar[1].IDKontak) {
		t.Errorf("CONTACT_ID: %q dan %q", daftar[0].IDKontak, daftar[1].IDKontak)
	}
	if _, err := k.Ubah(ctx, "UJI-ADM", "UJI-ADM", login.IsianAkun{Nama: "Uji Adm", Organisasi: "UJI-ORG",
		Menu: []string{menu.KodeKelolaUser, "claimlife"}, Kontak: login.Kontak{Email: " UJI.Adm@Nusantara.EXAMPLE "}}); err != nil {
		t.Fatalf("email ADM: %v", err)
	}
	if _, _, err := l.Masuk(ctx, "uji.adm@nusantara.example", "Sandi-Admin-01"); !errors.Is(err, login.ErrKredensial) {
		t.Errorf("masuk dengan email harus ditolak: %v", err)
	}
	if _, err := k.Ubah(ctx, "UJI-ADM", "UJI-B", login.IsianAkun{Nama: "Uji B Baru", Organisasi: "UJI-ORG",
		Menu: []string{"claimlife"}, Kontak: login.Kontak{Email: "uji.adm@nusantara.example"}}); !errors.Is(err, login.ErrEmailSudahTerdaftar) {
		t.Errorf("email ganda: %v", err)
	}
	if err := k.Buat(ctx, "UJI-ADM", login.AkunBaru{ID: "uji-b", Nama: "Uji b"}, "Sandi-Admin-01", true); !errors.Is(err, login.ErrAkunSudahAda) {
		t.Errorf("username beda huruf: %v", err)
	}
	if _, err := sqlDB.ExecContext(ctx, `INSERT INTO `+skema+`.M_LOGIN_GO (LOGIN_ID, NAME, PASSWORD_HASH, CONTACT_ID, EMAIL)
		VALUES ('UJI-GANDA', 'Uji Ganda', 'x', 'CON-UJI-GANDA', 'UJI.ADM@nusantara.example')`); err == nil ||
		!strings.Contains(err.Error(), "UX_M_LOGIN_GO_EMAIL") {
		t.Errorf("index unik email: %v", err)
	}
}
