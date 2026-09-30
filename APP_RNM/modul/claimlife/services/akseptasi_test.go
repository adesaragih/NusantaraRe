package services_test

// Akseptasi baris adjustment dari Claim Life sendiri - TANPA Oracle.
//
// Pemilik: audit A0 (brief lanjutan 4 bab 7). Dibaca sesudah: akseptasi.go.

import (
	"errors"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/models"
	"nusantarare/modul/claimlife/services"
)

// TestAwalanNomorAkseptasiPerType mengunci kedua RDB.
//
// `[terverifikasi]` `Claim Life/RDBList/Generate_NoAccept_Life.xml` baris 85:
// `SELECT 'RNML-A'||{pyWorkPage.BusinessCode}||'.'||{TempGenerate.CARI1}||'.'||
// {TempGenerate.CARI2}||'.'||LPAD(TO_CHAR(ACCEPTATIONNOLIFE_SEQ.NEXTVAL),5,'0')`;
// `Generate_NoAccept_LifeRetro.xml` baris 85 sama persis dengan awalan
// `'RNML-AR'`.
func TestAwalanNomorAkseptasiPerType(t *testing.T) {
	for tipe, mau := range map[string]string{
		services.TypeQP: "RNML-A",
		services.TypeQR: "RNML-A",
		services.TypeTP: "RNML-AR",
		services.TypeTR: "RNML-AR",
	} {
		got, err := services.AwalanNomorAkseptasi(tipe)
		if err != nil {
			t.Errorf("Type %q: %v", tipe, err)
			continue
		}
		if got != mau {
			t.Errorf("Type %q -> %q, mau %q", tipe, got, mau)
		}
	}
	// ⛔ Type asing tidak ditebak: ia menentukan awalan nomor yang tercetak di
	// dokumen, dan menebaknya berarti menerbitkan nomor yang salah seri.
	for _, asing := range []string{"", "QQ", "qp", "TPX"} {
		if _, err := services.AwalanNomorAkseptasi(asing); !errors.Is(
			err, services.ErrTypeTidakDikenal) {
			t.Errorf("Type %q: galat = %v, mau ErrTypeTidakDikenal", asing, err)
		}
	}
}

// TestPeriodeNomorAkseptasiAmbangDuaPuluhLima - `[terverifikasi]`
// `SaveAdjustment_Act.xml` pecahan baris 521:
// `@if(@toDecimal(Local.currentdate)>25, CurrentMonth+1, CurrentMonth)`;
// baris 542 memadankannya jadi dua digit; baris 563 mengubah `"13"` jadi
// `"01"`.
//
// ⛔ PENYIMPANGAN SADAR - dilaporkan, bukan disembunyikan. `[terverifikasi]`
// baris 605 berbunyi
// `@if(MM=="12" && NextMonth=="01", @toDecimal(@CurrentDate("YY")), @toDecimal(@CurrentDate("YY")))`
// - KEDUA cabangnya IDENTIK, sehingga tahun Pega TIDAK bergeser ketika bulan
// berguling 12 → 01. Akibatnya nomor Januari memakai tahun Desember. Go
// menggeser tahunnya; bila paritas persis yang diinginkan, ubah di satu tempat
// dan test ini yang merah lebih dulu.
func TestPeriodeNomorAkseptasiAmbangDuaPuluhLima(t *testing.T) {
	for _, k := range []struct {
		saat         string
		mauMM, mauYY string
		apa          string
	}{
		{"2026-03-25 10:00:00", "03", "26", "tepat 25: belum bergeser"},
		{"2026-03-26 10:00:00", "04", "26", "lewat 25: bulan berikutnya"},
		{"2026-03-01 10:00:00", "03", "26", "awal bulan"},
		{"2026-12-25 10:00:00", "12", "26", "Desember tepat 25"},
		{"2026-12-26 10:00:00", "01", "27", "Desember lewat 25: 13 -> 01 DAN tahun bergeser"},
		{"2026-09-30 23:59:59", "10", "26", "akhir September"},
	} {
		saat, err := time.Parse("2006-01-02 15:04:05", k.saat)
		if err != nil {
			t.Fatal(err)
		}
		mm, yy := services.PeriodeNomorAkseptasi(saat)
		if mm != k.mauMM || yy != k.mauYY {
			t.Errorf("%s (%s): (%q,%q), mau (%q,%q)",
				k.saat, k.apa, mm, yy, k.mauMM, k.mauYY)
		}
	}
}

// TestRakitNomorAkseptasi - bentuk utuhnya.
func TestRakitNomorAkseptasi(t *testing.T) {
	got := services.RakitNomorAkseptasi("RNML-A", "L1", "04", "26", "3509")
	const mau = "RNML-AL1.04.26.03509"
	if got != mau {
		t.Errorf("nomor = %q, mau %q", got, mau)
	}
	// ⛔ LPAD lima digit. Urut yang lebih panjang TIDAK dipotong - memotongnya
	// menerbitkan nomor yang bertabrakan dengan nomor lama.
	if got := services.RakitNomorAkseptasi("RNML-AR", "L2", "12", "26", "123456"); got != "RNML-ARL2.12.26.123456" {
		t.Errorf("urut enam digit dipotong: %q", got)
	}
}

// TestGerbangAksepMengikutiPrasyaratXML - gerbang dagang jalur akseptasi.
//
// `[terverifikasi]` `SaveAdjustment_Act.xml` pecahan baris 854 dan 1048:
// `.ACCEPTEDNO=="" && .IsCheck=true && .STS_REJECT=="0"` beserta `Type`.
//
// ⛔ PENYIMPANGAN SADAR - dilaporkan. Prasyarat XML-nya ditulis TANPA KURUNG:
// `.IsCheck=true&&Type=="QP"||Type=="QR"` (baris 837, 1031). Karena `&&`
// mengikat lebih erat daripada `||`, bacaan harfiahnya meloloskan `QR` dan
// `TR` TANPA memeriksa `IsCheck` sama sekali. Go memakai bacaan yang
// DIMAKSUD - `IsCheck && (QP || QR)` - sebab meloloskan baris peserta yang
// tidak dipilih ke akseptasi adalah cacat, bukan fitur.
func TestGerbangAksepMengikutiPrasyaratXML(t *testing.T) {
	baik := models.Peserta{ID: "P-1", IsCheck: "true"}
	barisBaik := models.BarisAdjustment{ID: "A-1", KodeStatus: kontrak.KodeOutstanding}

	if err := services.PeriksaBolehAksep(baik, barisBaik); err != nil {
		t.Errorf("peserta dan baris yang sah ditolak: %v", err)
	}
	for _, k := range []struct {
		apa      string
		p        models.Peserta
		b        models.BarisAdjustment
		sentinel error
	}{
		{"peserta tidak dipilih", models.Peserta{ID: "P-1", IsCheck: "false"},
			barisBaik, services.ErrPesertaTidakDipilih},
		{"IsCheck kosong", models.Peserta{ID: "P-1"},
			barisBaik, services.ErrPesertaTidakDipilih},
		{"baris sudah bernomor akseptasi", baik,
			models.BarisAdjustment{ID: "A-1", KodeStatus: kontrak.KodeOutstanding,
				NomorAkseptasi: "RNML-AL1.04.26.00001"}, services.ErrBarisSudahBernomorAkseptasi},
		{"baris sudah Aksep", baik,
			models.BarisAdjustment{ID: "A-1", KodeStatus: kontrak.KodeAksep},
			services.ErrBarisBukanOutstanding},
		{"baris sudah Ditolak", baik,
			models.BarisAdjustment{ID: "A-1", KodeStatus: kontrak.KodeDitolak},
			services.ErrBarisBukanOutstanding},
	} {
		if err := services.PeriksaBolehAksep(k.p, k.b); !errors.Is(err, k.sentinel) {
			t.Errorf("%s: galat = %v, mau %v", k.apa, err, k.sentinel)
		}
	}
}

// TestWajibPemegangTahap - gerbang peran menurut POHON XML, bukan tebakan.
//
// `[terverifikasi — pohon XML, brief lanjutan 4 bab 7]` Tombol "Save
// Adjustment" di `Section/ClaimLifeDetailGCNM.xml` TIDAK dibungkus gerbang
// peran mana pun; seluruh leluhurnya bervisibilitas ALWAYS. Yang menentukan
// siapa boleh menekannya adalah SIAPA MEMEGANG TAHAP, sebab sectionnya
// diluncurkan dari layar tahap: Outstanding → `ReasLifeAdmin`, Medical Check →
// `ReasLifeMedicalAdvisor`, Claim Analis → `ReasLifeSPV` (`Register_Flow`,
// 14 `Property-Set pyWorkPage.pyPosition`).
//
// ⛔ Karena itu `ErrAksepBukanDariModulIni` DIHAPUS: akseptasi punya DUA
// jalur, dan tiket 04 hanya mengenal satu.
func TestWajibPemegangTahap(t *testing.T) {
	for peran, tahap := range map[string]models.Tahap{
		models.PeranAdminLife:   models.TahapOutstanding,
		models.PeranMedicalLife: models.TahapMedicalCheck,
		models.PeranSPVLife:     models.TahapClaimAnalis,
	} {
		if err := services.WajibPemegangTahap(pelakuBerperan(peran), tahap); err != nil {
			t.Errorf("%s di tahapnya sendiri ditolak: %v", peran, err)
		}
	}
	// ⛔ Peran lain di tahap yang bukan miliknya ditolak.
	if err := services.WajibPemegangTahap(
		pelakuBerperan(models.PeranMedicalLife), models.TahapClaimAnalis); !errors.Is(
		err, inti.ErrTanpaWewenang) {
		t.Errorf("Medical di Claim Analis: galat = %v, mau ErrTanpaWewenang", err)
	}
	// Tahap tak dikenal tidak punya pemegang, dan tidak ditebak.
	if err := services.WajibPemegangTahap(
		pelakuBerperan(models.PeranAdminLife), models.TahapTidakDikenal); err == nil {
		t.Error("tahap tak dikenal lolos tanpa galat")
	}
}

// TestAksepBukanLagiTerlarangDiModulIni - ralat tiket 04 dan 07.
func TestAksepBukanLagiTerlarangDiModulIni(t *testing.T) {
	// ⛔ Dulu: peran mana pun yang menuju Aksep ditolak
	// ErrAksepBukanDariModulIni. Kini Aksep adalah tujuan yang sah, dan yang
	// menggerbanginya pemegang tahap - diperiksa `SimpanAdjustment`.
	err := services.WajibPeranPengubahStatus(
		pelakuBerperan(inti.PeranAdmin), kontrak.StatusAksep)
	if err != nil {
		t.Errorf("Aksep ditolak di lapisan peran: %v; ia kini jalur sah modul "+
			"ini (SaveAdjustment_Act)", err)
	}
}
