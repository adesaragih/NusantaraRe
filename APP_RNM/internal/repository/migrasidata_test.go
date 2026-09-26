package repository

// Test pembongkar data lama - TANPA Oracle.
//
// Untuk apa berkas ini: membuktikan bahwa pemindahan dari tabel datar warisan
// ke pohon klaim yang baru tidak kehilangan baris, tidak mengubah satu digit
// uang pun, dan melaporkan - bukan menebak - apa yang tidak dapat diurai.
//
// Fixture di sini SINTETIS: nol nama orang, nol nomor polis nyata, nol potongan
// data produksi. Seluruh nilai berawalan UJI.

import (
	"strings"
	"testing"

	"nusantarare/pkg/utils"
)

// contohBaris membuat satu baris lama yang wajar, lalu memberi kesempatan
// memutar satu-dua medan.
func contohBaris(id, caseID, sertifikat, jumlah string) BarisLama {
	return BarisLama{
		ID:             id,
		CASEID:         caseID,
		NO_CLAIM:       "UJI-CLM-0001",
		POLICY_NO:      "UJI-POL-0001",
		CERTIFICATE_NO: sertifikat,
		BUSINESSNAME:   "UJI BISNIS",
		CLAIM_RETRO:    "UJI-RETRO",
		PL_NUMBER:      "UJI-PL-1",
		CURRENCY:       "IDR",
		CLAIM_AMOUNT:   jumlah,
		STS_REJECT:     "1",
		NO_ACCEPTATION: "UJI-AKSEP-1",
		TYPE:           "UJI-TYPE",
		CREATEOPNAME:   "UJI-OPERATOR",
	}
}

// ADR-U-0011: SELURUH baris adjustment ikut pindah, bukan hanya yang terakhir.
func TestSeluruhBarisAdjustmentIkutPindah(t *testing.T) {
	masuk := []BarisLama{
		contohBaris("R1", "UJI-CASE-1", "006", "100.00000001"),
		contohBaris("R2", "UJI-CASE-1", "006", "200"),
		contohBaris("R3", "UJI-CASE-1", "010", "300"),
	}
	pohon, lap := BongkarBarisLama(masuk)

	if len(pohon) != 1 {
		t.Fatalf("klaim terbentuk %d, mau 1", len(pohon))
	}
	if n := len(pohon[0].Klaim.Peserta); n != 2 {
		t.Fatalf("peserta %d, mau 2", n)
	}
	if n := len(pohon[0].Klaim.Peserta[0].Baris); n != 2 {
		t.Errorf("peserta pertama punya %d baris, mau 2 - baris pertama tertimpa?", n)
	}
	// Jumlah baris sesudah migrasi sama dengan sebelumnya.
	if lap.BarisMasuk != lap.AdjustmentTerbentuk {
		t.Errorf("masuk %d baris, terbentuk %d adjustment", lap.BarisMasuk, lap.AdjustmentTerbentuk)
	}
	if pohon[0].Klaim.CacahBaris() != 3 {
		t.Errorf("cacah baris pohon = %d, mau 3", pohon[0].Klaim.CacahBaris())
	}
}

// Setiap baris menggantung pada peserta yang benar.
func TestBarisMenunjukPesertaYangBenar(t *testing.T) {
	masuk := []BarisLama{
		contohBaris("R1", "UJI-CASE-1", "006", "100"),
		contohBaris("R2", "UJI-CASE-1", "010", "200"),
		contohBaris("R3", "UJI-CASE-1", "006", "300"),
	}
	pohon, _ := BongkarBarisLama(masuk)
	ps := pohon[0].Klaim.Peserta

	if ps[0].NomorSertifikat != "006" || ps[1].NomorSertifikat != "010" {
		t.Fatalf("urutan peserta salah: %q, %q", ps[0].NomorSertifikat, ps[1].NomorSertifikat)
	}
	idPeserta0 := []string{ps[0].Baris[0].ID, ps[0].Baris[1].ID}
	if idPeserta0[0] != "R1" || idPeserta0[1] != "R3" {
		t.Errorf("baris peserta 006 = %v, mau [R1 R3]", idPeserta0)
	}
	if ps[1].Baris[0].ID != "R2" {
		t.Errorf("baris peserta 010 = %q, mau R2", ps[1].Baris[0].ID)
	}
}

// ADR-U-0022: nomor sertifikat berawalan nol tetap teks.
func TestSertifikatBerawalanNolTetapUtuh(t *testing.T) {
	pohon, _ := BongkarBarisLama([]BarisLama{contohBaris("R1", "UJI-CASE-1", "006", "1")})
	if got := pohon[0].Klaim.Peserta[0].NomorSertifikat; got != "006" {
		t.Errorf("nomor sertifikat = %q, mau %q", got, "006")
	}
}

// ADR-U-0003 / ADR-U-0021: uang pindah tanpa berubah satu digit, dibandingkan
// TEPAT - bukan dengan toleransi.
func TestUangPindahTanpaBerubahSatuDigit(t *testing.T) {
	mau := []string{"1234567890.12345678", "0.00000001", "250000", "0"}
	var masuk []BarisLama
	for i, v := range mau {
		masuk = append(masuk, contohBaris("R"+string(rune('1'+i)), "UJI-CASE-1", "006", v))
	}
	pohon, lap := BongkarBarisLama(masuk)
	if !lap.Bersih() {
		for _, x := range lap.Temuan {
			if x.Jenis != TemuanNilaiHardcode {
				t.Errorf("temuan tak terduga: %s", x)
			}
		}
	}
	baris := pohon[0].Klaim.Peserta[0].Baris
	if len(baris) != len(mau) {
		t.Fatalf("baris %d, mau %d", len(baris), len(mau))
	}
	for i, v := range mau {
		got := utils.FormatDecimal(baris[i].JumlahKlaim.Amount)
		if got != v {
			t.Errorf("baris %d: jumlah = %q, mau %q", i, got, v)
		}
		if baris[i].JumlahKlaim.Currency != "IDR" {
			t.Errorf("baris %d: mata uang hilang", i)
		}
	}
}

// Uang yang tidak dapat diurai DILAPORKAN dan dibiarkan kosong - tidak
// dibulatkan, tidak menjadi nol (ADR-U-0027).
func TestUangTakTeruraiDilaporkanBukanDitebak(t *testing.T) {
	b := contohBaris("R1", "UJI-CASE-1", "006", "seratus ribu")
	pohon, lap := BongkarBarisLama([]BarisLama{b})

	adj := pohon[0].Klaim.Peserta[0].Baris[0]
	if !adj.JumlahKlaim.Kosong() {
		t.Errorf("jumlah tak terurai menjadi %q", adj.JumlahKlaim.String())
	}
	if !adaTemuan(lap, TemuanUangTakTerurai, "CLAIM_AMOUNT") {
		t.Error("nilai uang tak terurai tidak dilaporkan")
	}
}

// Tanggal teks menjadi DATE tanpa pergeseran zona waktu; yang tidak dapat
// diurai dilaporkan.
func TestTanggalDiuraiTanpaGeserZona(t *testing.T) {
	b := contohBaris("R1", "UJI-CASE-1", "006", "1")
	b.ACCEPTATION_DATE = "2026-09-25"
	pohon, _ := BongkarBarisLama([]BarisLama{b})
	tgl := pohon[0].Klaim.Peserta[0].Baris[0].TanggalAkseptasi
	if utils.FormatTanggal(tgl) != "2026-09-25" {
		t.Errorf("tanggal = %q, mau 2026-09-25 - ada pergeseran zona", utils.FormatTanggal(tgl))
	}
	if y, m, d := tgl.Date(); y != 2026 || int(m) != 9 || d != 25 {
		t.Errorf("komponen tanggal bergeser: %v", tgl)
	}

	b2 := contohBaris("R2", "UJI-CASE-2", "006", "1")
	b2.ACCEPTATION_DATE = "kemarin"
	pohon2, lap2 := BongkarBarisLama([]BarisLama{b2})
	if !pohon2[0].Klaim.Peserta[0].Baris[0].TanggalAkseptasi.IsZero() {
		t.Error("tanggal tak terurai tidak dibiarkan kosong")
	}
	if !adaTemuan(lap2, TemuanTanggalTakTerurai, "ACCEPTATION_DATE") {
		t.Error("tanggal tak terurai tidak dilaporkan")
	}
}

// Atribut tingkat klaim yang BERBEDA antar baris peserta dilaporkan, dan tidak
// diam-diam dipilih salah satu.
func TestAtributKlaimBerbedaDilaporkan(t *testing.T) {
	a := contohBaris("R1", "UJI-CASE-1", "006", "1")
	b := contohBaris("R2", "UJI-CASE-1", "010", "1")
	b.BUSINESSNAME = "UJI BISNIS LAIN"
	_, lap := BongkarBarisLama([]BarisLama{a, b})

	if !adaTemuan(lap, TemuanAtributBerbeda, "BUSINESSNAME") {
		t.Fatal("perbedaan atribut klaim tidak dilaporkan")
	}
	for _, x := range lap.Temuan {
		if x.Jenis == TemuanAtributBerbeda && x.Medan == "BUSINESSNAME" {
			if !strings.Contains(x.Nilai, "UJI BISNIS") || !strings.Contains(x.Nilai, "UJI BISNIS LAIN") {
				t.Errorf("kedua nilai tidak ikut dilaporkan: %q", x.Nilai)
			}
		}
	}
}

// Nilai yang berasal dari hardcode di sumbernya dilaporkan apa adanya.
func TestNilaiHardcodeDilaporkan(t *testing.T) {
	b := contohBaris("R1", "UJI-CASE-1", "006", "1")
	b.STS_REJECT = "0"
	_, lap := BongkarBarisLama([]BarisLama{b})

	if lap.BarisHardcode != 1 {
		t.Errorf("BarisHardcode = %d, mau 1", lap.BarisHardcode)
	}
	if !adaTemuan(lap, TemuanNilaiHardcode, "ACCEPTATION_DATE, STS_REJECT") {
		t.Error("nilai hardcode tidak dilaporkan")
	}
}

// Baris tanpa CASEID tidak dapat digantung ke klaim mana pun: dilaporkan, dan
// tidak dipindahkan diam-diam.
func TestBarisTanpaCaseIDDilaporkan(t *testing.T) {
	a := contohBaris("R1", "", "006", "1")
	b := contohBaris("R2", "UJI-CASE-1", "006", "1")
	pohon, lap := BongkarBarisLama([]BarisLama{a, b})

	if len(pohon) != 1 {
		t.Fatalf("klaim terbentuk %d, mau 1", len(pohon))
	}
	if lap.AdjustmentTerbentuk != 1 {
		t.Errorf("adjustment terbentuk %d, mau 1", lap.AdjustmentTerbentuk)
	}
	if lap.BarisMasuk != 2 {
		t.Errorf("BarisMasuk = %d, mau 2", lap.BarisMasuk)
	}
	if !adaTemuan(lap, TemuanIdentitasKosong, "CASEID") {
		t.Error("baris tanpa CASEID tidak dilaporkan")
	}
}

// Kode status dibawa apa adanya sebagai teks, tidak diterjemahkan saat migrasi.
func TestKodeStatusDibawaSebagaiTeks(t *testing.T) {
	b := contohBaris("R1", "UJI-CASE-1", "006", "1")
	b.STS_REJECT = "006"
	pohon, _ := BongkarBarisLama([]BarisLama{b})
	if got := pohon[0].Klaim.Peserta[0].Baris[0].KodeStatus; got != "006" {
		t.Errorf("kode status = %q, mau %q", got, "006")
	}
}

func adaTemuan(l LaporanRekonsiliasi, jenis, medan string) bool {
	for _, t := range l.Temuan {
		if t.Jenis == jenis && t.Medan == medan {
			return true
		}
	}
	return false
}

// Penulisan ganda: pohon -> baris datar -> pohon tidak kehilangan baris dan
// tidak mengubah satu digit uang pun.
func TestPulangPergiPohonDanBarisDatar(t *testing.T) {
	masuk := []BarisLama{
		contohBaris("R1", "UJI-CASE-1", "006", "1234567890.12345678"),
		contohBaris("R2", "UJI-CASE-1", "006", "0.00000001"),
		contohBaris("R3", "UJI-CASE-1", "010", "250000"),
	}
	pohon, _ := BongkarBarisLama(masuk)
	if len(pohon) != 1 {
		t.Fatalf("klaim %d, mau 1", len(pohon))
	}

	datar := BarisLamaDari(pohon[0])
	if len(datar) != len(masuk) {
		t.Fatalf("baris datar %d, mau %d", len(datar), len(masuk))
	}

	balik, lap := BongkarBarisLama(datar)
	if len(balik) != 1 {
		t.Fatalf("klaim balik %d, mau 1", len(balik))
	}
	if lap.AdjustmentTerbentuk != len(masuk) {
		t.Errorf("adjustment balik %d, mau %d", lap.AdjustmentTerbentuk, len(masuk))
	}

	for i, ps := range balik[0].Klaim.Peserta {
		asal := pohon[0].Klaim.Peserta[i]
		if ps.NomorSertifikat != asal.NomorSertifikat {
			t.Errorf("peserta %d: sertifikat %q, mau %q", i, ps.NomorSertifikat, asal.NomorSertifikat)
		}
		if len(ps.Baris) != len(asal.Baris) {
			t.Fatalf("peserta %d: %d baris, mau %d", i, len(ps.Baris), len(asal.Baris))
		}
		for j, adj := range ps.Baris {
			mau := utils.FormatDecimal(asal.Baris[j].JumlahKlaim.Amount)
			got := utils.FormatDecimal(adj.JumlahKlaim.Amount)
			if got != mau {
				t.Errorf("peserta %d baris %d: jumlah %q, mau %q", i, j, got, mau)
			}
			if adj.KodeStatus != asal.Baris[j].KodeStatus {
				t.Errorf("peserta %d baris %d: kode status berubah", i, j)
			}
		}
	}
}

// AC 11: ketiga kolom bank ikut pindah, dengan nama barunya.
//
// Nama warisannya berbeda dari nama baru: NAME_OF_BANK tetap, IDBANK menjadi
// ID_BANK, ACCOUNTNO menjadi ACCOUNT_NO. Ronde 1 membuat kolomnya di tabel
// tetapi pembongkarnya tidak pernah mengisi - test ini yang menguncinya.
func TestKolomBankIkutPindah(t *testing.T) {
	b := contohBaris("R1", "UJI-CASE-1", "006", "100.00")
	b.NAME_OF_BANK = "UJI BANK NUSANTARA"
	b.IDBANK = "UJI-BANK-014"
	// ⛔ Berawalan nol dengan sengaja: nomor rekening adalah TEKS. Kalau ia
	// pernah menjadi bilangan, nol di depannya hilang dan test ini gagal.
	b.ACCOUNTNO = "0012345678"

	pohon, lap := BongkarBarisLama([]BarisLama{b})
	if len(pohon) != 1 || len(pohon[0].Klaim.Peserta) != 1 ||
		len(pohon[0].Klaim.Peserta[0].Baris) != 1 {
		t.Fatalf("bentuk pohon tidak seperti yang diharapkan: %+v", pohon)
	}
	adj := pohon[0].Klaim.Peserta[0].Baris[0]

	if adj.NamaBank != "UJI BANK NUSANTARA" {
		t.Errorf("NamaBank = %q, mau %q", adj.NamaBank, "UJI BANK NUSANTARA")
	}
	if adj.IDBank != "UJI-BANK-014" {
		t.Errorf("IDBank = %q, mau %q", adj.IDBank, "UJI-BANK-014")
	}
	if adj.NomorRekening != "0012345678" {
		t.Errorf("NomorRekening = %q, mau %q - nol di depan hilang?",
			adj.NomorRekening, "0012345678")
	}
	for _, tm := range lap.Temuan {
		if strings.Contains(tm.Medan, "BANK") || strings.Contains(tm.Medan, "ACCOUNT") {
			t.Errorf("kolom bank menghasilkan temuan yang tidak diharapkan: %+v", tm)
		}
	}
}

// Ketiga kolom bank kembali ke nama warisannya saat baris datar disusun ulang.
func TestKolomBankPulangPergi(t *testing.T) {
	b := contohBaris("R1", "UJI-CASE-1", "006", "100.00")
	b.NAME_OF_BANK = "UJI BANK NUSANTARA"
	b.IDBANK = "UJI-BANK-014"
	b.ACCOUNTNO = "0012345678"

	pohon, _ := BongkarBarisLama([]BarisLama{b})
	if len(pohon) != 1 {
		t.Fatalf("klaim terbentuk %d, mau 1", len(pohon))
	}
	kembali := BarisLamaDari(pohon[0])
	if len(kembali) != 1 {
		t.Fatalf("dapat %d baris, mau 1", len(kembali))
	}
	if kembali[0].NAME_OF_BANK != b.NAME_OF_BANK {
		t.Errorf("NAME_OF_BANK = %q, mau %q", kembali[0].NAME_OF_BANK, b.NAME_OF_BANK)
	}
	if kembali[0].IDBANK != b.IDBANK {
		t.Errorf("IDBANK = %q, mau %q", kembali[0].IDBANK, b.IDBANK)
	}
	if kembali[0].ACCOUNTNO != b.ACCOUNTNO {
		t.Errorf("ACCOUNTNO = %q, mau %q", kembali[0].ACCOUNTNO, b.ACCOUNTNO)
	}
}
