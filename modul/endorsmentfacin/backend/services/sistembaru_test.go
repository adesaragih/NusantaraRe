package services

// Uji polis buatan NB FacIn sistem baru (A51 gerbang 1, A52 gerbang bisnis) - DATA_JSON minimal `{"IDNewBisnis": …}`.

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/repository"
	"nusantarare/modul/endorsmentfacin/backend/tiruan"
)

func TestBisnisPolisSistemBaruBelumDiperiksa(t *testing.T) {
	d := dataBersih()
	// Gerbang 2 diperiksa: `lengkap` hanya bergantung pada gerbang bisnis di uji ini.
	d.SistemBaru, d.JenisBisnis, d.G2Diperiksa = true, "", true
	h := susun(t, d)
	if h.Bisnis.Hasil != BD || h.Bisnis.Catatan != catatanBisnisSistemBaru || h.Lengkap || h.Ditolak || len(h.Pesan) != 0 {
		t.Fatalf("bisnis = %+v lengkap %v pesan %v", h.Bisnis, h.Lengkap, teksPesan(h))
	}
	// Polis sistem baru yang DATA_JSON-nya membawa BusinessType: dinilai seperti biasa.
	d.JenisBisnis = "Fire"
	if h := susun(t, d); h.Bisnis.Hasil != L || !h.Lengkap {
		t.Fatalf("sistem baru + BusinessType: %+v", h.Bisnis)
	}
	// Polis sistem baru berpola RNML: langkah 8 "Life" tetap berlaku → ditolak lini (Life di luar EDM Fac In).
	d.PolicyNo, d.JenisBisnis = polisUjiLife, ""
	if h := susun(t, d); h.Bisnis.Hasil != T || h.Bisnis.JenisBisnis != "Life" {
		t.Fatalf("sistem baru RNML: %+v", h.Bisnis)
	}
	// Polis era Pega tanpa BusinessType: tetap "Invalid policy no!" (korpus apa adanya).
	d = dataBersih()
	d.JenisBisnis = ""
	if h := susun(t, d); !samaIrisan(teksPesan(h), []string{"PolicyNo|Invalid policy no!"}) || h.Bisnis.Hasil != T {
		t.Fatalf("era Pega: %+v %v", h.Bisnis, teksPesan(h))
	}
}

func TestGerbangSatuPolisSistemBaru(t *testing.T) {
	for _, k := range []struct {
		nama        string
		status      string
		takTerbaca  bool
		klep        map[KlepPembatal]bool
		want        models.HasilGerbangPeriksa
		catatan     string
		pesanGagal1 bool
	}{
		{"tanpa EDM berproduksi", "", false, nil, L, catatanG1SistemBaru, false},
		{"EDM batal sejak semula", "1", false, nil, T, "", true},
		{"EDM batal prorata", "2", false, nil, T, "", true},
		{"EDM perubahan", "4", false, nil, L, "", false},
		{"batal + klep", "1", false, map[KlepPembatal]bool{KlepBatal: true}, models.GerbangDiklep, catatanDiklep, false},
		{"tabel 262 tidak ada", "", true, nil, BD, catatanKasusBelumAda, false},
	} {
		d := dataBersih()
		d.SistemBaru, d.StatusEDM, d.G1TakTerbaca, d.Klep = true, k.status, k.takTerbaca, k.klep
		h := susun(t, d)
		if h.Gerbang[0].Hasil != k.want || h.Gerbang[0].Catatan != k.catatan {
			t.Errorf("%s: gerbang1 = %+v", k.nama, h.Gerbang[0])
		}
		if k.pesanGagal1 != samaIrisan(teksPesan(h), []string{"PolicyNo|Sudah Di endorsement Batal"}) {
			t.Errorf("%s: pesan = %v", k.nama, teksPesan(h))
		}
		if k.takTerbaca && h.Lengkap {
			t.Errorf("%s: lengkap padahal gerbang 1 belum diperiksa", k.nama)
		}
	}
}

func TestLayananSumberGerbangSatu(t *testing.T) {
	// Era Pega: GetEDMStatus_SQL; tabel kasus sistem baru tidak dibaca.
	g := gudangBersih()
	if _, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji}); err != nil {
		t.Fatal(err)
	}
	if !g.Dipanggil("PolisSistemBaru") || !g.Dipanggil("StatusEDM") || g.Dipanggil("EdmTypeSelesaiProduksi") {
		t.Fatalf("era Pega: panggilan %v", g.Panggilan)
	}
	// Sistem baru: kasus EDM berproduksi; DATA_JSON EdmType tidak dibaca.
	g = gudangBersih()
	g.SistemBaru, g.JenisBisnisNilai, g.EdmTypeProduksi = true, "", "2"
	h, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji})
	if err != nil {
		t.Fatal(err)
	}
	if g.Dipanggil("StatusEDM") || !g.Dipanggil("EdmTypeSelesaiProduksi") {
		t.Fatalf("sistem baru: panggilan %v", g.Panggilan)
	}
	if h.Gerbang[0].Hasil != T || h.Bisnis.Hasil != BD || h.Lengkap {
		t.Fatalf("sistem baru: gerbang1 %+v bisnis %+v lengkap %v", h.Gerbang[0], h.Bisnis, h.Lengkap)
	}
	// RI slip: G1 dilompati, tetapi penanda tetap dibaca (gerbang bisnis memakainya).
	g = gudangBersih()
	g.SistemBaru, g.JenisBisnisNilai = true, ""
	h, err = BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji, EdmType: "4", EdmTypeNew: "4"})
	if err != nil {
		t.Fatal(err)
	}
	if !g.Dipanggil("PolisSistemBaru") || g.Dipanggil("EdmTypeSelesaiProduksi") || h.Bisnis.Hasil != BD {
		t.Fatalf("RI slip sistem baru: %v %+v", g.Panggilan, h.Bisnis)
	}
	// Tabel 262 tidak ada: gerbang 1 belum diperiksa; galat lain diteruskan.
	g = gudangBersih()
	g.SistemBaru, g.GalatEdmTypeProduksi = true, fmt.Errorf("%w: ORA-00942", repository.ErrBelumAda)
	if h, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji}); err != nil ||
		h.Gerbang[0].Hasil != BD {
		t.Fatalf("tabel 262 tidak ada: %+v %v", h.Gerbang[0], err)
	}
	g = gudangBersih()
	g.SistemBaru, g.GalatEdmTypeProduksi = true, errors.New("ORA-03113")
	if _, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji}); err == nil {
		t.Fatal("galat EdmTypeSelesaiProduksi ditelan")
	}
}

// A51 opsi 1a: DATA_JSON dulu; polis sistem baru tanpa BusinessType → FACINPRODUCTION → BUSINESS → kontrak.
func TestJenisBisnisDariTabelLewatKontrak(t *testing.T) {
	jb := func() *tiruan.JenisBisnisTiruan {
		return &tiruan.JenisBisnisTiruan{Peta: map[string]string{"006/22": "FireStyle1", "004/": "MarineCargo"}}
	}
	periksa := func(g *tiruan.Gudang, j *tiruan.JenisBisnisTiruan) models.HasilPeriksa {
		t.Helper()
		h, err := BaruLayanan(g).DenganJenisBisnis(j).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji})
		if err != nil {
			t.Fatal(err)
		}
		return h
	}
	// Satu baris BUSINESS → FireStyle1 → IsFire → lolos, dengan catatan sumber.
	g, j := gudangBersih(), jb()
	g.SistemBaru, g.JenisBisnisNilai, g.Bisnis = true, "", []repository.BarisBisnis{{GroupPanel: "006", OldID: "22"}}
	h := periksa(g, j)
	if h.Bisnis != (models.BisnisPeriksa{Hasil: L, JenisBisnis: "FireStyle1", Catatan: catatanBisnisDariTabel}) || h.Ditolak {
		t.Fatalf("satu baris: %+v %v", h.Bisnis, teksPesan(h))
	}
	if !samaIrisan(j.Masukan, []string{"006/22"}) {
		t.Fatalf("masukan kontrak = %v", j.Masukan)
	}
	// OLDID kosong diteruskan apa adanya (DeT: baris tanpa OldID cocok apa pun).
	g, j = gudangBersih(), jb()
	g.SistemBaru, g.JenisBisnisNilai, g.Bisnis = true, "", []repository.BarisBisnis{{GroupPanel: "004"}}
	if h := periksa(g, j); h.Bisnis.JenisBisnis != "MarineCargo" || h.Bisnis.Hasil != L {
		t.Fatalf("OLDID kosong: %+v", h.Bisnis)
	}
	// Nol atau dua baris → ("", "") → UNKNOWN → "Invalid business, please contact IT!" (pola NB), bukan "Invalid policy no!".
	for _, baris := range [][]repository.BarisBisnis{nil, {{GroupPanel: "006", OldID: "22"}, {GroupPanel: "003", OldID: "27"}}} {
		g, j = gudangBersih(), jb()
		g.SistemBaru, g.JenisBisnisNilai, g.Bisnis = true, "", baris
		h := periksa(g, j)
		if h.Bisnis.JenisBisnis != "UNKNOWN" || !samaIrisan(teksPesan(h), []string{"PolicyNo|Invalid business, please contact IT!"}) ||
			!samaIrisan(j.Masukan, []string{"/"}) {
			t.Fatalf("%d baris: %+v %v masukan %v", len(baris), h.Bisnis, teksPesan(h), j.Masukan)
		}
	}
	// DATA_JSON membawa BusinessType: tabel tidak dibaca (urutan sumber).
	g, j = gudangBersih(), jb()
	g.SistemBaru, g.JenisBisnisNilai = true, "Fire"
	if h := periksa(g, j); g.Dipanggil("BisnisProduksi") || len(j.Masukan) != 0 || h.Bisnis.Catatan != "" {
		t.Fatalf("DATA_JSON terisi: %v %+v", g.Panggilan, h.Bisnis)
	}
	// Polis era Pega tanpa BusinessType: tabel TIDAK dibaca, korpus apa adanya.
	g, j = gudangBersih(), jb()
	g.JenisBisnisNilai = ""
	if h := periksa(g, j); g.Dipanggil("BisnisProduksi") || !samaIrisan(teksPesan(h), []string{"PolicyNo|Invalid policy no!"}) {
		t.Fatalf("era Pega: %v %v", g.Panggilan, teksPesan(h))
	}
	// Galat BUSINESS diteruskan.
	g, j = gudangBersih(), jb()
	g.SistemBaru, g.JenisBisnisNilai, g.GalatBisnis = true, "", errors.New("ORA-03113")
	if _, err := BaruLayanan(g).DenganJenisBisnis(j).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji}); err == nil {
		t.Fatal("galat BUSINESS ditelan")
	}
}

// Buat atas polis sistem baru dengan kontrak tersambung: kasus lahir (G1 tanpa EDM berproduksi, lini dari tabel).
func TestBuatKasusPolisSistemBaruDenganKontrak(t *testing.T) {
	g := gudangKasus()
	g.SistemBaru, g.JenisBisnisNilai, g.Bisnis = true, "", []repository.BarisBisnis{{GroupPanel: "006", OldID: "22"}}
	j := &tiruan.JenisBisnisTiruan{Peta: map[string]string{"006/22": "FireStyle1"}}
	h, err := layananKasus(g).DenganJenisBisnis(j).Buat(context.Background(), akunUji, masukanKasus())
	if err != nil || h.CaseID != "EDM-00001" {
		t.Fatalf("hasil = %+v %v", h, err)
	}
}

// Buat atas polis sistem baru: lini belum dapat diperiksa → tidak lengkap → kasus tidak lahir.
func TestBuatKasusPolisSistemBaru(t *testing.T) {
	g := gudangKasus()
	g.SistemBaru, g.JenisBisnisNilai = true, ""
	h, err := layananKasus(g).Buat(context.Background(), akunUji, masukanKasus())
	if err != nil {
		t.Fatal(err)
	}
	if h.CaseID != "" || h.Periksa == nil || h.Periksa.Lengkap || h.Periksa.Ditolak || h.Periksa.Bisnis.Hasil != BD {
		t.Fatalf("hasil = %+v", h.Periksa)
	}
	tanpaTulisan(t, g)
}

// gerbangtolak.go: BisnisBelumDiperiksa menahan langkah 8-18 (aditif; bawaan = korpus).
func TestGerbangPenolakanBisnisBelumDiperiksa(t *testing.T) {
	m := MasukanGerbang{PolicyNo: polisUji, EdmType: models.EdmPerubahan, EdmTypeNew: "1"}
	h, err := PeriksaGerbangPenolakan(m)
	if err != nil || !samaIrisan(teksPesanMedan(h.Pesan), []string{"Invalid policy no!"}) {
		t.Fatalf("bawaan: %v %v", h.Pesan, err)
	}
	m.BisnisBelumDiperiksa = true
	if h, err := PeriksaGerbangPenolakan(m); err != nil || len(h.Pesan) != 0 || h.Ditolak {
		t.Fatalf("belum diperiksa: %v %v", h.Pesan, err)
	}
	// Pesan gerbang tetap terpasang.
	m.Query.CacahRenewal = 1
	if h, _ := PeriksaGerbangPenolakan(m); !samaIrisan(teksPesanMedan(h.Pesan), []string{"There's RNW with this policy no!"}) {
		t.Fatalf("gerbang 5 hilang: %v", h.Pesan)
	}
}

func teksPesanMedan(p []PesanMedan) []string {
	out := []string{}
	for _, x := range p {
		out = append(out, x.Teks)
	}
	return out
}
