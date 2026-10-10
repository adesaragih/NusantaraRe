package models

// Uji penulis laporan pemuat dokumen lama endorsemen (tiket EDM 10): arsip CSV medan tanpa kolom (F3 NB: arsip
// audit, bukan penampung), CSV galat ber-PRODKE, ringkasan (selesai = nol galat dan nol BELUM DIPUTUSKAN), cacah
// penanda migrasi (tiket 09) dan salinan SuggestList.

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

func bacaCSVUji(t *testing.T, b *bytes.Buffer) [][]string {
	t.Helper()
	baris, err := csv.NewReader(bytes.NewReader(b.Bytes())).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	return baris
}

func TestLaporanArsipMedanEDM(t *testing.T) {
	var arsip, gal bytes.Buffer
	l, err := LaporanPemuatBaru(&arsip, &gal, true)
	if err != nil {
		t.Fatal(err)
	}
	d := strings.Replace(dokumenUjiEDMProp, `"UJIMedanFiktif": "UJI-nilai-lewat"`, `"UJIMedanFiktif": "UJI-a, \"b\"\nbaris kedua"`, 1)
	h, err := PecahDokumenEDM(barisUjiEDM(d))
	if err != nil {
		t.Fatal(err)
	}
	p := PenandaMigrasi{Spreading: []PenandaBaris{{true, true}, {false, true}}, Angsuran: []PenandaBaris{{false, true}},
		Lapisan: []PenandaBaris{{true, false}}}
	if err := l.Berhasil(h, UsulanDisalin, p); err != nil {
		t.Fatal(err)
	}
	if err := l.Tutup(); err != nil {
		t.Fatal(err)
	}
	baris := bacaCSVUji(t, &arsip)
	if strings.Join(baris[0], ",") != "POLIS_ID,JALUR,NILAI,KEPUTUSAN" {
		t.Fatalf("kepala arsip %q", baris[0])
	}
	ada := map[string]bool{}
	for _, b := range baris[1:] {
		ada[strings.Join(b, "|")] = true
	}
	for _, b := range [][]string{
		{"ASM-FW-GISFW-WORK EDMT-990002", "PolicyTreatyIn.UJIMedanFiktif", "UJI-a, \"b\"\nbaris kedua", KeputusanBelumDiputuskan},
		{"ASM-FW-GISFW-WORK EDMT-990002", "PolicyTreatyIn.OldData.PremiOgp", "150", "dibuang: old_data"},
		{"ASM-FW-GISFW-WORK EDMT-990002", "PolicyTreatyIn.TreatyDifference.TotalPremium", "130.000000276", "dibuang: selisih_turunan"},
		{"ASM-FW-GISFW-WORK EDMT-990002", "PolicyTreatyIn.TreatyXOLList(1).ValueList(1).GrossPremi", "5", "dibuang: prop_hasil_antara"},
	} {
		if !ada[strings.Join(b, "|")] {
			t.Errorf("baris arsip %q tidak ada", b)
		}
	}
	if n := len(baris) - 1; n != len(h.Arsip) {
		t.Errorf("arsip %d baris, harap %d", n, len(h.Arsip))
	}
	r := l.Ringkasan()
	if r.Dimuat != 1 || r.MedanBelumDiputuskan != 1 || r.Selesai() || r.AngkaKasusTerbesar != 990002 || r.GenerasiTerbesar != 2 {
		t.Errorf("ringkasan %+v", r)
	}
	if r.Penanda != (RingkasanPenanda{BarisSpreading: 2, BarisAngsuran: 1, BarisLapisan: 1, BergeserSpreading: 1,
		BergeserLapisan: 1, Berlapis: 3}) {
		t.Errorf("penanda %+v", r.Penanda)
	}
	if r.Usulan != (RingkasanUsulan{Disalin: 2, TanpaPIC: 1}) {
		t.Errorf("usulan %+v", r.Usulan)
	}
	for _, s := range []string{"BELUM DIPUTUSKAN", "BELUM SELESAI", "PASANGAN_BERGESER", "RUMUS_BERLAPIS", "EDMT-990002",
		"HISTORYAKSEPTASIPRODUCTION"} {
		if !strings.Contains(r.Teks(), s) {
			t.Errorf("teks ringkasan tanpa %q:\n%s", s, r.Teks())
		}
	}
	if strings.Contains(r.Teks(), "UJI-a") {
		t.Error("ringkasan memuat NILAI medan - nilai hanya di berkas CSV")
	}
}

func TestLaporanGalatEDMBerPRODKE(t *testing.T) {
	var tak, gal bytes.Buffer
	l, _ := LaporanPemuatBaru(&tak, &gal, false)
	b := barisUjiEDM(strings.Replace(dokumenUjiEDMProp, `"StartDate": "20171001"`, `"StartDate": "05/06/2017"`, 1))
	h, err := PecahDokumenEDM(b)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Gagal(b, h.Galat); err != nil {
		t.Fatal(err)
	}
	if err := l.Gagal(b, []GalatDokumen{{Err: ErrPercabangan}}); err != nil {
		t.Fatal(err)
	}
	l.Lewat(ErrBukanTreatyIn)
	l.Lewat(ErrBukanGenerasiEndorsemen)
	_ = l.Tutup()
	baris := bacaCSVUji(t, &gal)
	if len(baris) != 3 || strings.Join(baris[0], ",") != "IDPEGA,NOPOLIS,PRODKE,JALUR,NILAI,SEBAB" {
		t.Fatalf("CSV galat %q", baris)
	}
	if baris[1][2] != "2" || baris[1][3] != "PolicyTreatyIn.StartDate" || !strings.Contains(baris[1][5], "ambigu") {
		t.Errorf("baris galat %q", baris[1])
	}
	r := l.Ringkasan()
	if r.DokumenGagal != 2 || r.TanggalAmbigu != 1 || r.BukanTreatyIn != 1 || r.GenerasiNB != 1 || r.Selesai() {
		t.Errorf("ringkasan %+v", r)
	}
	if r.GalatPerJenis[JenisGalat(ErrPercabangan)] != 1 || JenisGalat(ErrPercabangan) == JenisGalatSimpan {
		t.Errorf("galat per jenis %v", r.GalatPerJenis)
	}
}

func TestRingkasanEDMSelesaiHanyaBilaNolGalatDanNolBelumDiputuskan(t *testing.T) {
	var arsip, gal bytes.Buffer
	l, _ := LaporanPemuatBaru(&arsip, &gal, true)
	d := strings.Replace(dokumenUjiEDMProp, `"UJIMedanFiktif": "UJI-nilai-lewat",`, ``, 1)
	h, err := PecahDokumenEDM(barisUjiEDM(d))
	if err != nil || len(h.BelumDiputuskan) != 0 {
		t.Fatalf("%v %+v", err, h.BelumDiputuskan)
	}
	_ = l.Berhasil(h, UsulanDilewati, PenandaMigrasi{})
	l.SudahDimuat()
	_ = l.Tutup()
	r := l.Ringkasan()
	if !r.Selesai() || r.SudahDimuat != 1 || r.Usulan.DokumenSudahAda != 1 || !strings.Contains(r.Teks(), "STATUS: SELESAI") {
		t.Errorf("ringkasan %+v\n%s", r, r.Teks())
	}
	if len(bacaCSVUji(t, &arsip))-1 != len(h.Arsip) || len(h.Arsip) == 0 {
		t.Error("medan yang dibuang tetap diarsipkan")
	}
}
