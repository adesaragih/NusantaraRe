package services_test

// Uji panel Attachment modul Adjustment - NOL koneksi Oracle.
//
// ⭐ 8 Oktober 2026: nama kategori dari `M_KATEGORIMASTERTREATY` (katalog
// RD Pega), urut seperti `GetMasterTreatyCategory_SQL` — sama dengan
// modul Treaty In.

import (
	"context"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

// katalogUji - kesebelas pasangan `M_KATEGORIMASTERTREATY` (terukur
// 8 Oktober 2026).
func katalogUji() map[string]string {
	return map[string]string{
		"00000": "Others",
		"00001": "Analysed Email",
		"00002": "Approval Email",
		"00003": "Binding, signed share Email",
		"00004": "Info Pack",
		"00005": "Summary Treaty Leader",
		"00006": "Assessment Inward Treaty Form / Format Analisa Treaty",
		"00007": "Pega Proportional Calculation /Perhitungan Pega Proportional",
		"00008": "Letter of Acknowledgment / LOA",
		"00009": "Claim Data",
		"00010": "Offer Email",
	}
}

func namaBaris(baris []models.BarisKategoriLampiran) []string {
	out := make([]string, 0, len(baris))
	for _, b := range baris {
		out = append(out, b.Nama)
	}
	return out
}

// ⭐ Penamaan dan urutan = tangkapan layar Pega (`order by note`).
func TestKategoriBernamaDanUrutSepertiPega(t *testing.T) {
	got := namaBaris(services.SusunKategoriLampiran(katalogUji(), nil, true))
	mau := []string{
		"Analysed Email",
		"Approval Email",
		"Assessment Inward Treaty Form / Format Analisa Treaty",
		"Binding, signed share Email",
		"Claim Data",
		"Info Pack",
		"Letter of Acknowledgment / LOA",
		"Offer Email",
		"Others",
		"Pega Proportional Calculation /Perhitungan Pega Proportional",
		"Summary Treaty Leader",
	}
	if strings.Join(got, "|") != strings.Join(mau, "|") {
		t.Errorf("urutan\n got %q\nmau %q", got, mau)
	}
}

// Keempat kode yang dulu tampil "belum dipastikan" kini BERNAMA dari katalog.
func TestEmpatKodeDuluTanpaNamaKiniBernama(t *testing.T) {
	nama := map[string]string{}
	for _, b := range services.SusunKategoriLampiran(katalogUji(), nil, true) {
		if !b.Dipastikan {
			t.Errorf("kode %s belum dipastikan", b.Kode)
		}
		nama[b.Kode] = b.Nama
	}
	for kode, mau := range map[string]string{
		"00003": "Binding, signed share Email", "00004": "Info Pack",
		"00008": "Letter of Acknowledgment / LOA", "00009": "Claim Data",
	} {
		if nama[kode] != mau {
			t.Errorf("%s = %q, mau %q", kode, nama[kode], mau)
		}
	}
}

// `GetMasterTreatyCategory_Act` [2.2]: kontrak Non-Prop, kode 00007 bernama
// Non-Prop di tempat yang sama.
func TestKontrakNonPropMemakaiNamaNonProp00007(t *testing.T) {
	for _, b := range services.SusunKategoriLampiran(katalogUji(), nil, false) {
		if b.Kode == "00007" && b.Nama != "Pega Non Proportional Calculation /Perhitungan Pega Non Proportional" {
			t.Errorf("00007 Non-Prop = %q", b.Nama)
		}
	}
	if !services.SifatProporsional("Proportional") || !services.SifatProporsional("") || services.SifatProporsional("NonProportional") {
		t.Error("SifatProporsional salah")
	}
}

// Kategori BERNOL berkas tetap tampil; cacah menurut kode, atau nama bila
// baris lama tidak berkode.
func TestKategoriNolBerkasTetapTampil(t *testing.T) {
	lampiran := []models.BarisLampiranWarisan{
		{KodeKategori: "00001"}, {KodeKategori: "00001"},
		{NamaKategori: "Claim Data"},
	}
	baris := services.SusunKategoriLampiran(katalogUji(), lampiran, true)
	if len(baris) != 11 {
		t.Fatalf("%d baris, mau 11", len(baris))
	}
	cacah := map[string]int{}
	for _, b := range baris {
		cacah[b.Kode] = b.Cacah
	}
	if cacah["00001"] != 2 || cacah["00009"] != 1 || cacah["00000"] != 0 || cacah["00003"] != 0 {
		t.Errorf("cacah salah: %v", cacah)
	}
}

// ⛔ Kode di data yang tidak ada di katalog tetap TAMPIL (di ujung), tanpa
// nama tebakan.
func TestKodeAsingTetapTampilTanpaNama(t *testing.T) {
	baris := services.SusunKategoriLampiran(katalogUji(), []models.BarisLampiranWarisan{{KodeKategori: "00099"}}, true)
	u := baris[len(baris)-1]
	if len(baris) != 12 || u.Kode != "00099" || u.Nama != "" || u.Dipastikan || u.Cacah != 1 {
		t.Errorf("baris asing %+v (%d baris)", u, len(baris))
	}
}

// ⭐ Spanduk biru itu ATURAN. Diuji DUA ARAH.
func TestNamaBerkasAmanDuaArah(t *testing.T) {
	for _, n := range []string{
		"laporan treaty.pdf", "laporan(1).pdf", "a/b.pdf", "a\b.pdf",
		"../rahasia.pdf", "surat&lampiran.pdf", "", "   ", ".", "..",
	} {
		if services.NamaBerkasAman(n) {
			t.Errorf("%q DITERIMA; spanduk menyebut titik atau garis bawah sebagai pengganti aman", n)
		}
	}
	// Pemeriksa yang menolak SEGALANYA lulus uji di atas dengan sempurna,
	// dan menutup pintu yang sistem lama buka.
	for _, n := range []string{
		"laporan.pdf", "laporan_treaty.pdf", "laporan-treaty.pdf",
		"laporan.2026.pdf", "LAPORAN.PDF", "Treaty_In_2026-01.xlsx",
	} {
		if !services.NamaBerkasAman(n) {
			t.Errorf("%q DITOLAK; ia hanya huruf, angka, titik, garis bawah, tanda hubung", n)
		}
	}
}

// Tanpa identitas ditolak SEBELUM gudang disentuh.
func TestLampiranMenolakTanpaIdentitasSebelumGudang(t *testing.T) {
	g := &gudangTiruan{katalogKategori: katalogUji()}
	l := services.LayananDengan(g)

	if _, err := l.LampiranKontrakWarisan(context.Background(), inti.Pelaku{}, "1001851", ""); err == nil {
		t.Fatal("permintaan tanpa identitas DITERIMA")
	}
}

// Pengenal kosong ditolak - "tidak ada" berbeda dari "tidak ditanyakan".
func TestLampiranMenolakPengenalKosong(t *testing.T) {
	g := &gudangTiruan{katalogKategori: katalogUji()}
	l := services.LayananDengan(g)

	for _, id := range []string{"", "   "} {
		if _, err := l.LampiranKontrakWarisan(context.Background(), pelakuAda, id, ""); err == nil {
			t.Errorf("pengenal %q DITERIMA", id)
		}
	}
}

// Uji POSITIF: kategori tersusun, cacahnya benar, berkasnya diteruskan.
func TestLampiranPositif(t *testing.T) {
	g := &gudangTiruan{
		katalogKategori: katalogUji(),
		lampiran: []models.BarisLampiranWarisan{
			{ID: "1", KodeKategori: "00001", NamaBerkas: "a.pdf"},
			{ID: "2", KodeKategori: "00001", NamaBerkas: "b.pdf"},
		},
	}
	l := services.LayananDengan(g)

	hasil, err := l.LampiranKontrakWarisan(context.Background(), pelakuAda, "1001851", "Proportional")
	if err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	if len(hasil.Berkas) != 2 {
		t.Errorf("%d berkas, mau 2", len(hasil.Berkas))
	}
	if len(hasil.Kategori) != 11 {
		t.Errorf("%d kategori, mau 11", len(hasil.Kategori))
	}
	for _, k := range hasil.Kategori {
		if k.Kode == "00001" && k.Cacah != 2 {
			t.Errorf("kategori 00001 cacahnya %d, mau 2", k.Cacah)
		}
	}
}
