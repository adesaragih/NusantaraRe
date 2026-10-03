package login

// Kontak akun - email, nomor HP, NIK, jabatan (Kelola User, permintaan work owner 03-10-2026).

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestPeriksaKontakFormat(t *testing.T) {
	for _, k := range []struct {
		nama string
		isi  Kontak
		mau  error
	}{
		{"semua kosong sah (opsional)", Kontak{}, nil},
		{"semua terisi sah", Kontak{Email: "uji.user@nusantara.example", Telepon: "+62 812-3456-7890", NIK: "UJI/2026.001-A", Jabatan: "Underwriter Senior"}, nil},
		{"telepon 8 digit sah", Kontak{Telepon: "02112345"}, nil},
		{"telepon 15 digit sah", Kontak{Telepon: "123456789012345"}, nil},
		{"email tanpa @", Kontak{Email: "uji.nusantara.example"}, ErrEmailTidakSah},
		{"email tanpa domain bertitik", Kontak{Email: "uji@nusantara"}, ErrEmailTidakSah},
		{"email berspasi", Kontak{Email: "uji user@nusantara.example"}, ErrEmailTidakSah},
		{"email lebih dari 254 byte", Kontak{Email: strings.Repeat("a", 250) + "@b.cd"}, ErrEmailTidakSah},
		{"telepon 7 digit", Kontak{Telepon: "1234567"}, ErrTeleponTidakSah},
		{"telepon 16 digit", Kontak{Telepon: "1234567890123456"}, ErrTeleponTidakSah},
		{"telepon berhuruf", Kontak{Telepon: "0812-ABC-7890"}, ErrTeleponTidakSah},
		{"telepon + di tengah", Kontak{Telepon: "0812+34567890"}, ErrTeleponTidakSah},
		{"NIK bertanda lain", Kontak{NIK: "UJI#001"}, ErrNIKTidakSah},
		{"NIK berspasi", Kontak{NIK: "UJI 001"}, ErrNIKTidakSah},
		{"NIK lebih dari 30", Kontak{NIK: strings.Repeat("1", 31)}, ErrNIKTidakSah},
		{"jabatan lebih dari 150 byte", Kontak{Jabatan: strings.Repeat("J", 151)}, ErrJabatanTidakSah},
	} {
		if err := PeriksaKontak(k.isi); !errors.Is(err, k.mau) || (k.mau == nil && err != nil) {
			t.Errorf("%s: dapat %v, mau %v", k.nama, err, k.mau)
		}
	}
}

func TestKelolaBuatDanUbahMenyimpanKontakRapi(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	baru := AkunBaru{ID: "UJI-KONTAK-1", Nama: "Uji Kontak", Organisasi: "RNM", Menu: []string{"claimlife"},
		Kontak: Kontak{Email: " uji.kontak@nusantara.example ", Telepon: " 0812 3456 7890 ", NIK: " UJI-001 ", Jabatan: " Analyst "}}
	if err := k.Buat(ctxUji, "UJI-ADMIN", baru, "Sandi-Admin-01", true); err != nil {
		t.Fatal(err)
	}
	r, err := k.Rinci(ctxUji, "UJI-KONTAK-1")
	if err != nil {
		t.Fatal(err)
	}
	if mau := (Kontak{Email: "uji.kontak@nusantara.example", Telepon: "0812 3456 7890", NIK: "UJI-001", Jabatan: "Analyst"}); r.Kontak != mau {
		t.Errorf("kontak tersimpan rapi (spasi tepi dibuang):\n dapat %+v\n mau   %+v", r.Kontak, mau)
	}
	isi := IsianAkun{Nama: "Uji Kontak", Organisasi: "RNM", Menu: []string{"claimlife"},
		Kontak: Kontak{Email: "baru@nusantara.example", Jabatan: "Head of Unit"}}
	r, err = k.Ubah(ctxUji, "UJI-ADMIN", "UJI-KONTAK-1", isi)
	if err != nil {
		t.Fatal(err)
	}
	if mau := (Kontak{Email: "baru@nusantara.example", Jabatan: "Head of Unit"}); r.Kontak != mau {
		t.Errorf("ubah menulis keempat medan; yang dikosongkan menjadi kosong:\n dapat %+v\n mau   %+v", r.Kontak, mau)
	}
}

func TestKelolaKontakTidakSahDitolakTanpaTulisan(t *testing.T) {
	g := gudangUji(t)
	k := kelolaUji(g)
	buruk := AkunBaru{ID: "UJI-KONTAK-2", Nama: "Uji", Organisasi: "RNM", Menu: []string{"claimlife"}, Kontak: Kontak{Email: "bukan-email"}}
	if err := k.Buat(ctxUji, "UJI-ADMIN", buruk, "Sandi-Admin-01", true); !errors.Is(err, ErrEmailTidakSah) {
		t.Fatalf("buat dengan email tidak sah: %v", err)
	}
	if _, ada := g.akun["UJI-KONTAK-2"]; ada {
		t.Error("akun tidak boleh dibuat bila kontaknya ditolak")
	}
	if _, err := k.Ubah(ctxUji, "UJI-ADMIN", "UJI-BARU", IsianAkun{Nama: "Uji Baru", Menu: []string{"claimlife"},
		Kontak: Kontak{Telepon: "123"}}); !errors.Is(err, ErrTeleponTidakSah) {
		t.Fatalf("ubah dengan telepon tidak sah: %v", err)
	}
	if g.kontak["UJI-BARU"] != (Kontak{}) {
		t.Error("kontak tidak boleh tertulis bila ditolak")
	}
}

// Bentuk API: badan buat/ubah dan jawaban daftar membawa keempat kunci RATA (bukan objek bersarang).
func TestKontakJSONRata(t *testing.T) {
	var m isianPengguna
	badan := `{"akunId":"UJI-X","nama":"Uji","email":"u@x.example","telepon":"08123456789","nik":"UJI-9","jabatan":"Staff"}`
	if err := json.Unmarshal([]byte(badan), &m); err != nil {
		t.Fatal(err)
	}
	if m.Kontak != (Kontak{Email: "u@x.example", Telepon: "08123456789", NIK: "UJI-9", Jabatan: "Staff"}) {
		t.Errorf("badan terurai ke Kontak: %+v", m.Kontak)
	}
	keluar, err := json.Marshal(RingkasAkun{AkunID: "UJI-X", Kontak: m.Kontak})
	if err != nil {
		t.Fatal(err)
	}
	for _, kunci := range []string{`"email":"u@x.example"`, `"telepon":"08123456789"`, `"nik":"UJI-9"`, `"jabatan":"Staff"`} {
		if !strings.Contains(string(keluar), kunci) {
			t.Errorf("jawaban tanpa %s: %s", kunci, keluar)
		}
	}
	if strings.Contains(string(keluar), `"Kontak"`) {
		t.Errorf("kontak harus rata, bukan objek bersarang: %s", keluar)
	}
}

// Pesan galat berbahasa Inggris (permintaan work owner) dan sampai ke layar sebagai 400 berkalimat.
func TestPesanKontakBerbahasaInggris(t *testing.T) {
	for _, e := range []error{ErrEmailTidakSah, ErrTeleponTidakSah, ErrNIKTidakSah, ErrJabatanTidakSah} {
		pesan := strings.TrimPrefix(e.Error(), "login: ")
		if !strings.HasPrefix(pesan, "Email") && !strings.HasPrefix(pesan, "Phone Number") &&
			!strings.HasPrefix(pesan, "Employee ID") && !strings.HasPrefix(pesan, "Position") {
			t.Errorf("pesan diawali label Inggris medannya: %q", pesan)
		}
	}
}

// go-ora mengikat menurut URUTAN placeholder: kolom kontak dan nilainya (`nilaiKontak`) harus berurutan sama.
func TestSQLKontakUrutKolomDanBind(t *testing.T) {
	sisip := strings.Join(strings.Fields(sqlSisipAkun("T", "S")), " ")
	if !strings.Contains(sisip, "MUST_CHANGE_PASSWORD, EMAIL, PHONE_NUMBER, EMPLOYEE_ID, JOB_POSITION, CONTACT_ID)") ||
		!strings.Contains(sisip, ":8, :9, :10, :11, :12, :13 || S.NEXTVAL)") {
		t.Errorf("INSERT: kolom kontak sesudah MUST_CHANGE_PASSWORD, bind :9-:12, CONTACT_ID :13: %s", sisip)
	}
	ubah := strings.Join(strings.Fields(sqlUbahProfil("T")), " ")
	if !strings.Contains(ubah, "EMAIL = :5, PHONE_NUMBER = :6, EMPLOYEE_ID = :7, JOB_POSITION = :8") ||
		!strings.Contains(ubah, "WHERE LOGIN_ID = :9") {
		t.Errorf("UPDATE: kontak :5-:8, LOGIN_ID :9: %s", ubah)
	}
	daftar := strings.Join(strings.Fields(sqlDaftarAkun("T", false)), " ")
	if !strings.Contains(daftar, "LAST_LOGIN, EMAIL, PHONE_NUMBER, EMPLOYEE_ID, JOB_POSITION FROM T") {
		t.Errorf("SELECT: kontak sesudah LAST_LOGIN (urutan pindaiRingkas): %s", daftar)
	}
	if n := nilaiKontak(Kontak{Email: "e", Telepon: "t", NIK: "n", Jabatan: "j"}); len(n) != 4 || n[0] != "e" || n[1] != "t" || n[2] != "n" || n[3] != "j" {
		t.Errorf("nilaiKontak berurutan EMAIL, PHONE_NUMBER, EMPLOYEE_ID, JOB_POSITION: %v", n)
	}
	if n := nilaiKontak(Kontak{}); n[0] != nil || n[3] != nil {
		t.Errorf("kosong = NULL: %v", n)
	}
}
