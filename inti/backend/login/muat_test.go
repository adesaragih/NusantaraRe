package login

// Pemuat akun dari ekspor operator Pega, password tanpa batasan, password bawaan, dan login tanpa beda huruf (perintah
// work owner 11-10-2026). Fixture UJI- saja.

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

const csvLoginUji = "\xef\xbb\xbfPYEMAILADDRESS;PYORGDIVISION;PYORGUNIT;PYORGANIZATION;PYPOSITION;PYUSERIDENTIFIER;PYUSERNAME\n" +
	"uji1@contoh.test;Claim;Life;RNM;Claim Staff;UJI-MUAT1;Uji Muat Satu\n" +
	"UJI1@contoh.test;Claim;NonLife;RNM;;uji-muat2;Uji Muat Dua\n" + // email kembar (beda huruf) -> dikosongkan
	"-;DivisionHead;Unit;RNM;Kepala;UJI-MUAT3;Uji Muat Tiga\n" + // email tidak sah, divisi tanpa padanan
	"uji4@contoh.test;UnderWriting;Unit;RNM;;UJI-MUAT4;Uji Muat Empat\n" + // TECH/UW: UW tidak ada di master uji
	"uji5@contoh.test;Claim;Life;RNM;;uji-admin;Uji Admin Lagi\n" + // akun sudah ada (beda huruf) -> dilewati
	"uji-dipakai@contoh.test;Claim;Life;RNM;;UJI-MUAT6;Uji Muat Enam\n" + // email dipakai akun lain -> dikosongkan
	"uji7@contoh.test;Claim;Life;RNM;;UJI-MUAT7;Uji Muat Tujuh\n" + // CONTACT_ID sheet sudah dipakai -> gagal
	"\n"

const csvKontakUji = "CONTACT,NAMA,AKSES LOGIN\n" +
	"ASM-SFAGIS-WORK-CONTACT CON-49,UJI MUAT SATU,uji-muat1\n" +
	"ASM-SFAGIS-WORK-CONTACT CON-77,UJI MUAT TUJUH,UJI-MUAT7\n" +
	"ASM-SFAGIS-WORK-CONTACT CON-114,UJI TANPA LOGIN,\n"

// gudangMuatUji - gudang uji dengan email dan CONTACT_ID yang sudah dipakai akun lain.
func gudangMuatUji(t *testing.T) *gudangTiruan {
	g := gudangUji(t)
	g.tulisKontak("UJI-ADMIN", Kontak{Email: "uji-dipakai@contoh.test"})
	g.idKontak = map[string]string{"UJI-ADMIN": "CON-77"}
	return g
}

func barisMuatUji(t *testing.T) []BarisMuat {
	t.Helper()
	login, err := BacaCSV(strings.NewReader(csvLoginUji))
	if err != nil {
		t.Fatal(err)
	}
	kb, err := BacaCSV(strings.NewReader(csvKontakUji))
	if err != nil {
		t.Fatal(err)
	}
	kontak, err := PetaKontak(kb)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(kontak, map[string]string{"UJI-MUAT1": "CON-49", "UJI-MUAT7": "CON-77"}) {
		t.Fatalf("peta contact %v", kontak)
	}
	return SusunMuat(login, kontak)
}

func TestSusunMuatDariSheetLogin(t *testing.T) {
	b := barisMuatUji(t)
	if len(b) != 7 {
		t.Fatalf("baris %d, mau 7 (BOM dibuang, titik koma dikenali, baris kosong dilewati)", len(b))
	}
	a := b[0].Akun
	if a.ID != "UJI-MUAT1" || a.Nama != "Uji Muat Satu" || a.Organisasi != "RNM" || a.Divisi != "TECH" ||
		a.Unit != "CLM" || a.Email != "uji1@contoh.test" || a.Jabatan != "Claim Staff" || a.IDKontak != "CON-49" ||
		len(b[0].Catatan) != 0 {
		t.Errorf("baris 2: %+v %v", a, b[0].Catatan)
	}
	if b[1].Akun.Email != "" || b[1].Akun.IDKontak != "" || !strings.Contains(strings.Join(b[1].Catatan, ";"), "sudah dipakai baris 2") {
		t.Errorf("baris 3 (email kembar): %+v %v", b[1].Akun, b[1].Catatan)
	}
	if c := strings.Join(b[2].Catatan, ";"); b[2].Akun.Email != "" || b[2].Akun.Divisi != "" ||
		!strings.Contains(c, "tidak sah") || !strings.Contains(c, "tanpa padanan") {
		t.Errorf("baris 4 (email tidak sah, divisi tanpa padanan): %+v %v", b[2].Akun, b[2].Catatan)
	}
	if b[3].Akun.Divisi != "TECH" || b[3].Akun.Unit != "UW" {
		t.Errorf("baris 5 (UnderWriting): %+v", b[3].Akun)
	}
	if _, err := PetaKontak([]map[string]string{{kolomKontak: "BUKAN-KODE", kolomAksesLogin: "UJI"}}); err == nil {
		t.Error("CONTACT bukan CON-angka harus galat")
	}
}

func TestMuatPenggunaUjiKeringLaluJalankan(t *testing.T) {
	ctx := context.Background()
	mau := []string{StatusMuatSiap, StatusMuatSiap, StatusMuatSiap, StatusMuatGagal, StatusMuatSudahAda, StatusMuatSiap,
		StatusMuatGagal}

	g := gudangMuatUji(t)
	h, err := NewLayanan(g, rahasiaUji).MuatPengguna(ctx, barisMuatUji(t), false)
	if err != nil {
		t.Fatal(err)
	}
	status := func(hs []HasilMuat) []string {
		var s []string
		for _, x := range hs {
			s = append(s, x.Status)
		}
		return s
	}
	if !reflect.DeepEqual(status(h), mau) || len(g.dibuat) != 0 {
		t.Fatalf("uji-kering: %v (dibuat %d)", status(h), len(g.dibuat))
	}
	if c := strings.Join(h[5].Catatan, ";"); !strings.Contains(c, "sudah dipakai akun UJI-ADMIN") {
		t.Errorf("email dipakai akun basis data: %v", h[5].Catatan)
	}
	if c := strings.Join(h[6].Catatan, ";"); !strings.Contains(c, "CON-77") {
		t.Errorf("CONTACT_ID dipakai: %v", h[6].Catatan)
	}

	g = gudangMuatUji(t)
	l := layananUji(g, saatUji)
	sukses := []string{StatusMuatDibuat, StatusMuatDibuat, StatusMuatDibuat, StatusMuatGagal, StatusMuatSudahAda,
		StatusMuatDibuat, StatusMuatGagal}
	if h, err = l.MuatPengguna(ctx, barisMuatUji(t), true); err != nil || !reflect.DeepEqual(status(h), sukses) {
		t.Fatalf("jalankan: %v %v", status(h), err)
	}
	a := g.akun["UJI-MUAT1"]
	if a == nil || !CocokSandi(a.HashSandi, SandiBawaan) || a.WajibGantiSandi || g.idKontak["UJI-MUAT1"] != "CON-49" {
		t.Fatalf("akun UJI-MUAT1: %+v contact %q", a, g.idKontak["UJI-MUAT1"])
	}
	if !strings.HasPrefix(g.idKontak["uji-muat2"], AwalanIDKontak+"1") || g.kontak["UJI-MUAT6"].Email != "" {
		t.Errorf("contact sequence %q / email dikosongkan %q", g.idKontak["uji-muat2"], g.kontak["UJI-MUAT6"].Email)
	}
	if _, ada := g.akun["UJI-MUAT4"]; ada {
		t.Error("akun gagal tetap ditulis")
	}
	// login tanpa beda huruf dengan password bawaan
	p, _, err := l.Masuk(ctx, "uji-MUAT1", SandiBawaan)
	if err != nil || p.AkunID != "UJI-MUAT1" || p.WajibGantiSandi {
		t.Errorf("login akun muatan: %+v %v", p, err)
	}
}

// Login: username tanpa beda huruf (perintah work owner 11-10-2026); sesi memakai LOGIN_ID tersimpan.
func TestMasukUsernameTanpaBedaHuruf(t *testing.T) {
	p, tok, err := layananUji(gudangUji(t), saatUji).Masuk(context.Background(), " uji-admin ", "Sandi-Benar-01")
	if err != nil || p.AkunID != "UJI-ADMIN" || tok.Akun != "UJI-ADMIN" {
		t.Fatalf("login huruf kecil: %+v %+v %v", p, tok, err)
	}
	if q := strings.Join(strings.Fields(sqlAmbilAkun("T")), " "); !strings.Contains(q, "WHERE LOWER(LOGIN_ID) = LOWER(:1)") {
		t.Errorf("SQL login: %s", q)
	}
}

// Kelola User: password dikosongkan = SandiBawaan (perintah work owner 11-10-2026).
func TestKelolaBuatSandiKosongPakaiBawaan(t *testing.T) {
	g := gudangUji(t)
	a := AkunBaru{ID: "UJI-BAWAAN", Nama: "Uji Bawaan", Organisasi: "RNM", Menu: []string{"claimlife"}, IDKontak: "CON-5"}
	if err := kelolaUji(g).Buat(ctxUji, "UJI-ADMIN", a, "", false); err != nil {
		t.Fatal(err)
	}
	if !CocokSandi(g.akun["UJI-BAWAAN"].HashSandi, SandiBawaan) {
		t.Error("password kosong tidak menjadi 123456")
	}
	if g.idKontak["UJI-BAWAAN"] == "CON-5" {
		t.Error("Kelola User tidak boleh menetapkan CONTACT_ID")
	}
}
