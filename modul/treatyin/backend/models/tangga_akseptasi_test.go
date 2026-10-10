package models_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
)

func TestTanggaNaikSesuaiKeputusanPemilikProses(t *testing.T) {
	// ⭐ Tangga yang dinyatakan pemilik proses, dinaiki dari ujung ke ujung:
	// ReasTreatyInAdmin > SecHead > DeptHead > Director > Resolve Complete.
	urut := []struct {
		dari, ke, status string
	}{
		{models.PosisiAdmin, models.PosisiSecHead, models.StatusAccept},
		{models.PosisiSecHead, models.PosisiDeptHead, models.StatusAccept},
		{models.PosisiDeptHead, models.PosisiDirector, models.StatusAccept},
		{models.PosisiDirector, models.PosisiKosong, models.StatusTuntas},
	}
	posisi := models.PosisiAdmin
	for _, l := range urut {
		if posisi != l.dari {
			t.Fatalf("tangga putus: ada di %q, seharusnya %q", posisi, l.dari)
		}
		got, err := models.LangkahBerikut(posisi, models.PilihAccept, models.StatusAccept, false)
		if err != nil {
			t.Fatalf("%s Accept: %v", l.dari, err)
		}
		if got.Posisi != l.ke || got.Status != l.status {
			t.Fatalf("%s Accept = %q/%q, mau %q/%q", l.dari, got.Posisi, got.Status, l.ke, l.status)
		}
		posisi = got.Posisi
	}
}

func TestPosisiKosongSamaDenganAdmin(t *testing.T) {
	a, err := models.LangkahBerikut(models.PosisiKosong, models.PilihAccept, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if a.Posisi != models.PosisiSecHead {
		t.Fatalf("kosong Accept = %q, mau %q", a.Posisi, models.PosisiSecHead)
	}
}

func TestJalurTurun(t *testing.T) {
	for _, p := range []string{models.PosisiSecHead, models.PosisiDeptHead, models.PosisiDirector} {
		tolak, err := models.LangkahBerikut(p, models.PilihReject, models.StatusAccept, false)
		if err != nil {
			t.Fatalf("%s Reject: %v", p, err)
		}
		if tolak.Posisi != models.PosisiAdmin || tolak.Status != models.StatusReject {
			t.Fatalf("%s Reject = %q/%q", p, tolak.Posisi, tolak.Status)
		}
		// ⭐ Dikembalikan kepada penulis komentar PERTAMA - yang mengajukannya.
		if tolak.AsalNama != models.AsalNamaKomentarPertama {
			t.Fatalf("%s Reject asal nama = %q, mau komentar-pertama", p, tolak.AsalNama)
		}
		tampik, err := models.LangkahBerikut(p, models.PilihDecline, models.StatusAccept, false)
		if err != nil {
			t.Fatalf("%s Decline: %v", p, err)
		}
		if tampik.Posisi != models.PosisiKosong || tampik.Status != models.StatusDecline {
			t.Fatalf("%s Decline = %q/%q", p, tampik.Posisi, tampik.Status)
		}
	}
}

func TestAdminNolCabangTurun(t *testing.T) {
	// ⛔ Pengaju tidak dapat menolak berkasnya sendiri.
	for _, pilih := range []string{models.PilihReject, models.PilihDecline} {
		if _, err := models.LangkahBerikut(models.PosisiAdmin, pilih, "", false); !errors.Is(err, models.ErrPilihanTakBerlaku) {
			t.Fatalf("Admin %s = %v, mau ErrPilihanTakBerlaku", pilih, err)
		}
	}
}

func TestJalurRevisiTuntasDiSecHead(t *testing.T) {
	// ⛔ Revisi TIDAK lewat Dept Head maupun Director.
	naik, err := models.LangkahBerikut(models.PosisiAdmin, models.PilihAccept, "", true)
	if err != nil || naik.Posisi != models.PosisiSecHead {
		t.Fatalf("revisi Admin Accept = %q, %v", naik.Posisi, err)
	}
	tuntas, err := models.LangkahBerikut(models.PosisiSecHead, models.PilihAccept, models.StatusAccept, true)
	if err != nil {
		t.Fatal(err)
	}
	if tuntas.Status != models.StatusTuntas {
		t.Fatalf("revisi SecHead Accept = %q, mau %q", tuntas.Status, models.StatusTuntas)
	}
	// ⚠️ Keduanya DIKOSONGKAN, bukan dibiarkan.
	if tuntas.RevisionState == nil || *tuntas.RevisionState != "" {
		t.Fatalf("RevisionState = %v, mau pointer ke kosong", tuntas.RevisionState)
	}
	if tuntas.ViewState == nil || *tuntas.ViewState != "" {
		t.Fatalf("ViewState = %v, mau pointer ke kosong", tuntas.ViewState)
	}
	// ⛔ Dua anak tangga atas nol dicapai di jalur revisi.
	for _, p := range []string{models.PosisiDeptHead, models.PosisiDirector} {
		if _, err := models.LangkahBerikut(p, models.PilihAccept, models.StatusAccept, true); err == nil {
			t.Fatalf("revisi %s Accept seharusnya galat", p)
		}
	}
}

func TestRevisiRejectPakaiKomentarTerakhir(t *testing.T) {
	// ⚠️ Beda halus: jalur biasa memakai komentar PERTAMA, revisi TERAKHIR.
	r, err := models.LangkahBerikut(models.PosisiSecHead, models.PilihReject, models.StatusAccept, true)
	if err != nil {
		t.Fatal(err)
	}
	if r.AsalNama != models.AsalNamaKomentarTerakhir {
		t.Fatalf("asal nama = %q, mau komentar-terakhir", r.AsalNama)
	}
	if r.RevisionState == nil || *r.RevisionState != "1" {
		t.Fatalf("RevisionState = %v, mau \"1\"", r.RevisionState)
	}
	biasa, err := models.LangkahBerikut(models.PosisiSecHead, models.PilihReject, models.StatusAccept, false)
	if err != nil {
		t.Fatal(err)
	}
	if biasa.AsalNama == r.AsalNama {
		t.Fatal("jalur biasa dan revisi memakai asal nama yang SAMA - salah satunya salah")
	}
	// ⚠️ Dan jalur biasa nol menyentuh RevisionState.
	if biasa.RevisionState != nil || biasa.ViewState != nil {
		t.Fatal("jalur biasa menyentuh RevisionState/ViewState")
	}
}

func TestKontrakTuntasTakDapatDikirimUlang(t *testing.T) {
	if models.BolehKirim(models.StatusTuntas) {
		t.Fatal("BolehKirim(Resolve Complete) = true")
	}
	for _, p := range models.TanggaAkseptasi {
		if _, err := models.LangkahBerikut(p, models.PilihAccept, models.StatusTuntas, false); !errors.Is(err, models.ErrSudahTuntas) {
			t.Fatalf("%s pada status tuntas = %v", p, err)
		}
	}
}

// ⛔ PENJAGA: nol langkah boleh bermuara ke peran di luar tangga, dan nol kode
// boleh menghidupkan kembali cabang kode regu yang pemilik proses buang.
func TestPeranDiLuarTanggaNolMenjadiTujuan(t *testing.T) {
	diLuar := []string{"ReasTreatyInGroupLeader", "ReasTreatyInUnderwriting"}
	for _, p := range models.TanggaAkseptasi {
		for _, pilih := range []string{models.PilihAccept, models.PilihReject, models.PilihDecline} {
			for _, revisi := range []bool{false, true} {
				l, err := models.LangkahBerikut(p, pilih, models.StatusAccept, revisi)
				if err != nil {
					continue
				}
				for _, x := range diLuar {
					if l.Posisi == x {
						t.Fatalf("%s %s (revisi=%v) bermuara ke %q yang di luar tangga", p, pilih, revisi, x)
					}
				}
				if l.Posisi != models.PosisiKosong && !models.AdalahAnakTangga(l.Posisi) {
					t.Fatalf("%s %s bermuara ke %q yang bukan anak tangga", p, pilih, l.Posisi)
				}
			}
		}
	}
}

// tanpaKomentar membuang baris komentar.
//
// ⛔ Ini BUKAN pengurai Go - ia hanya membuang baris yang seluruhnya
// komentar. Itu cukup untuk penjaga di bawah dan sengaja tidak lebih: penguraian
// setengah jadi yang mengaku lengkap lebih berbahaya daripada yang sederhana
// dan dinyatakan sederhana.
func tanpaKomentar(src string) string {
	var b strings.Builder
	for _, baris := range strings.Split(src, "\n") {
		t := strings.TrimSpace(baris)
		if strings.HasPrefix(t, "//") || strings.HasPrefix(t, "*") || strings.HasPrefix(t, "/*") {
			continue
		}
		b.WriteString(baris)
		b.WriteByte('\n')
	}
	return b.String()
}

// ⚠️ Dan penyaring itu harus TETAP MENGGIGIT kode sungguhan. Penjaga yang
// menyaring terlalu banyak lulus tanpa memeriksa apa pun.
func TestPenyaringKomentarMasihMenggigit(t *testing.T) {
	if strings.Contains(tanpaKomentar("// kode regu SPVTREATY1 dibuang\n"), "SPVTREATY1") {
		t.Fatal("komentar tidak tersaring")
	}
	if !strings.Contains(tanpaKomentar("\tif regu == \"SPVTREATY1\" {\n"), "SPVTREATY1") {
		t.Fatal("kode sungguhan ikut tersaring - penjaga jadi buta")
	}
}

// ⛔ Kode regu `SPVTREATY1`/`TREATY1` - pemilik proses: "jangan digunakan
// dulu". Penjaga ini menangkap KODE yang menghidupkannya kembali; dokumentasi
// yang MENJELASKAN mengapa ia dibuang justru harus boleh menyebutnya, jadi
// komentar disaring lebih dulu.
func TestKodeReguNolDipakaiKode(t *testing.T) {
	akar := filepath.Join("..", "..", "..", "treatyin")
	kode := []string{"SPVTREATY1", "SPVTREATY2", "TREATY1", "TREATY2", "pyTelephone"}
	var temuan []string
	err := filepath.Walk(akar, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		ext := filepath.Ext(p)
		if ext != ".go" && ext != ".ts" && ext != ".tsx" {
			return nil
		}
		if strings.HasSuffix(p, "tangga_akseptasi_test.go") {
			return nil // berkas ini sendiri menyebutnya untuk menjaganya
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		isi := tanpaKomentar(string(b))
		for _, k := range kode {
			if strings.Contains(isi, k) {
				temuan = append(temuan, p+" → "+k)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(temuan) > 0 {
		t.Fatalf("kode regu hidup kembali di: %v", temuan)
	}
}
