package login

// Identitas akun - CONTACT_ID, username dan email sudah terdaftar (migrasi 905, keputusan work owner 03-10-2026 V1:
// "M_LOGIN_GO ID nya pake CON-xxx"; "tambahkan proteksi format email, atau email sudah terdaftar"; "tambahkan juga
// proteksi username sudah ada"). Login lewat email DIBATALKAN work owner ("LOGIN LEWAT EMAIL TIDAK JADI!!"): login
// tetap dengan username saja. TANPA Oracle.

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// go-ora mengikat menurut URUTAN placeholder: CONTACT_ID dibaca tepat sesudah LOGIN_ID (urutan pindaiRingkas).
func TestSQLDaftarMembacaContactIDSesudahLoginID(t *testing.T) {
	for _, satu := range []bool{false, true} {
		q := strings.Join(strings.Fields(sqlDaftarAkun("T", satu)), " ")
		if !strings.HasPrefix(q, "SELECT LOGIN_ID, CONTACT_ID, NAME, ORGANIZATION_CODE,") {
			t.Errorf("daftar akun (satu=%v): %s", satu, q)
		}
	}
}

func TestKelolaBuatMemberiContactIDYangTidakBerubah(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	for _, id := range []string{"UJI-CON-1", "UJI-CON-2"} {
		if err := k.Buat(ctxUji, "UJI-ADMIN", AkunBaru{ID: id, Nama: "Uji", Organisasi: "RNM", Menu: []string{"claimlife"}},
			"Sandi-Admin-01", true); err != nil {
			t.Fatal(err)
		}
	}
	r1, err := k.Rinci(ctxUji, "UJI-CON-1")
	if err != nil {
		t.Fatal(err)
	}
	r2, err := k.Rinci(ctxUji, "UJI-CON-2")
	if err != nil {
		t.Fatal(err)
	}
	if r1.IDKontak != "CON-1001" || r2.IDKontak != "CON-1002" {
		t.Fatalf("CONTACT_ID akun baru: %q dan %q, mau CON-1001 dan CON-1002", r1.IDKontak, r2.IDKontak)
	}
	ubah, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-CON-1", IsianAkun{Nama: "Uji Diubah", Organisasi: "RNM", Menu: []string{"claimlife"}})
	if err != nil {
		t.Fatal(err)
	}
	if ubah.IDKontak != "CON-1001" {
		t.Errorf("CONTACT_ID berubah sesudah ubah: %q", ubah.IDKontak)
	}
	keluar, err := json.Marshal(ubah.RingkasAkun)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(keluar), `"contactId":"CON-1001"`) {
		t.Errorf("jawaban tanpa contactId: %s", keluar)
	}
}

// Username sudah ada (permintaan work owner 03-10-2026 "tambahkan juga proteksi username sudah ada"): sama persis
// atau beda huruf saja - keduanya ditolak dan tidak ada akun yang dibuat.
func TestBuatUsernameSudahAdaDitolak(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	for _, id := range []string{"UJI-ADMIN", "uji-admin", "Uji-Baru"} {
		err := k.Buat(ctxUji, "UJI-ADMIN", AkunBaru{ID: id, Nama: "Uji", Organisasi: "RNM", Menu: []string{"claimlife"}},
			"Sandi-Admin-01", true)
		if !errors.Is(err, ErrAkunSudahAda) {
			t.Errorf("username %q: %v, mau ErrAkunSudahAda", id, err)
		}
	}
	if len(g.dibuat) != 0 {
		t.Errorf("akun tetap dibuat: %v", g.dibuat)
	}
}

// Email sudah terdaftar (permintaan work owner 03-10-2026 "tambahkan proteksi format email, atau email sudah
// terdaftar"): email akun lain, tanpa beda huruf. Format salah tetap ErrEmailTidakSah.
func TestEmailSudahTerdaftarDitolak(t *testing.T) {
	g := gudangUji(t)
	g.tulisKontak("UJI-BARU", Kontak{Email: "uji.baru@nusantara.example"})
	k := kelolaUji(g)
	buat := func(id, email string) error {
		return k.Buat(ctxUji, "UJI-ADMIN", AkunBaru{ID: id, Nama: "Uji", Organisasi: "RNM", Menu: []string{"claimlife"},
			Kontak: Kontak{Email: email}}, "Sandi-Admin-01", true)
	}
	for _, email := range []string{"uji.baru@nusantara.example", " UJI.Baru@Nusantara.EXAMPLE "} {
		if err := buat("UJI-EMAIL-1", email); !errors.Is(err, ErrEmailSudahTerdaftar) {
			t.Errorf("buat dengan email %q: %v, mau ErrEmailSudahTerdaftar", email, err)
		}
	}
	if err := buat("UJI-EMAIL-1", "bukan-email"); !errors.Is(err, ErrEmailTidakSah) {
		t.Errorf("format email salah: %v, mau ErrEmailTidakSah", err)
	}
	if len(g.dibuat) != 0 {
		t.Fatalf("akun tetap dibuat: %v", g.dibuat)
	}
	// Email disimpan huruf kecil.
	if err := buat("UJI-EMAIL-2", " Uji.Dua@Nusantara.EXAMPLE "); err != nil {
		t.Fatal(err)
	}
	if got := g.kontak["UJI-EMAIL-2"].Email; got != "uji.dua@nusantara.example" {
		t.Errorf("email tersimpan %q, mau huruf kecil", got)
	}
	// Ubah: email akun lain ditolak dan tidak tertulis; email milik sendiri tetap boleh.
	isi := IsianAkun{Nama: "Uji Dua", Organisasi: "RNM", Menu: []string{"claimlife"}, Kontak: Kontak{Email: "UJI.BARU@nusantara.example"}}
	if _, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-EMAIL-2", isi); !errors.Is(err, ErrEmailSudahTerdaftar) {
		t.Errorf("ubah ke email akun lain: %v", err)
	}
	if got := g.kontak["UJI-EMAIL-2"].Email; got != "uji.dua@nusantara.example" {
		t.Errorf("email tertimpa walau ditolak: %q", got)
	}
	isi.Kontak.Email = "Uji.Dua@nusantara.example"
	if _, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-EMAIL-2", isi); err != nil {
		t.Errorf("ubah dengan email sendiri: %v", err)
	}
}

// Pengaman terakhir di Oracle (dua simpan serentak lolos pemeriksaan aplikasi): ORA-00001 index email = email
// terdaftar; PK atau index username = username sudah ada.
func TestGalatGandaOracle(t *testing.T) {
	for _, k := range []struct {
		pesan string
		mau   error
	}{
		{"ORA-00001: unique constraint (POOLDATA.UX_M_LOGIN_GO_EMAIL) violated", ErrEmailSudahTerdaftar},
		{"ORA-00001: unique constraint (POOLDATA.UX_M_LOGIN_GO_LOGIN_ID) violated", ErrAkunSudahAda},
		{"ORA-00001: unique constraint (POOLDATA.PK_M_LOGIN_GO) violated", ErrAkunSudahAda},
	} {
		if err := galatGanda(errors.New(k.pesan)); !errors.Is(err, k.mau) {
			t.Errorf("%s: %v, mau %v", k.pesan, err, k.mau)
		}
	}
	if err := galatGanda(errors.New("ORA-01400: cannot insert NULL")); err != nil {
		t.Errorf("bukan ORA-00001 harus nil: %v", err)
	}
}

// Username dan email diperiksa TERPISAH, masing-masing tanpa beda huruf (index unik LOWER(...) migrasi 905).
func TestSQLPemakaiTanpaBedaHuruf(t *testing.T) {
	if q := strings.Join(strings.Fields(sqlPemakaiUsername("T")), " "); q != "SELECT LOGIN_ID FROM T WHERE LOWER(LOGIN_ID) = :1 ORDER BY LOGIN_ID" {
		t.Errorf("pemakai username: %s", q)
	}
	if q := strings.Join(strings.Fields(sqlPemakaiEmail("T")), " "); q != "SELECT LOGIN_ID FROM T WHERE LOWER(EMAIL) = :1 ORDER BY LOGIN_ID" {
		t.Errorf("pemakai email: %s", q)
	}
}

// Pesan 409 berbahasa Inggris, sama dengan label layar.
func TestRuteGandaMenjawab409(t *testing.T) {
	for _, k := range []struct {
		err   error
		pesan string
	}{
		{ErrAkunSudahAda, "Username is already registered"},
		{ErrEmailSudahTerdaftar, "Email is already registered to another account"},
	} {
		w := httptest.NewRecorder()
		tulisGalatKelola(w, k.err, "uji")
		if w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), k.pesan) {
			t.Errorf("%v: %d %s", k.err, w.Code, w.Body.String())
		}
	}
}

// Login lewat email DIBATALKAN work owner 03-10-2026 ("LOGIN LEWAT EMAIL TIDAK JADI!!"): email akun yang sah dan
// sandi yang benar tetap ditolak sebagai kredensial salah; username tetap jalan.
func TestMasukTidakLewatEmail(t *testing.T) {
	g := gudangUji(t)
	g.tulisKontak("UJI-ADMIN", Kontak{Email: "uji.admin@nusantara.example"})
	l := layananUji(g, saatUji)
	if _, _, err := l.Masuk(ctxUji, "uji.admin@nusantara.example", "Sandi-Benar-01"); !errors.Is(err, ErrKredensial) {
		t.Errorf("login dengan email: %v, mau ErrKredensial", err)
	}
	if _, _, err := l.Masuk(ctxUji, "UJI-ADMIN", "Sandi-Benar-01"); err != nil {
		t.Errorf("login username: %v", err)
	}
}
