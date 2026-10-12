package services_test

// Layar Outstanding dirombak mengikuti Pega (10-10-2026): gerbang tombol bawah, isian tombol Save, dan gerbang dokumen
// Save to RNM yang membaca IsAccept. TANPA Oracle.

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"nusantarare/inti/backend/galat"
	intiuang "nusantarare/inti/backend/uang"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/services"
)

func TestPindahDariOutstandingMengikutiGerbangPega(t *testing.T) {
	o, m, r := models.TahapOutstanding, models.TahapMedicalCheck, models.TahapInputRegister
	for _, u := range []struct {
		nama        string
		ke          models.Tahap
		nomor, tipe string
		mau         error
	}{
		// b21832: Send to Medical Check nonaktif selama CLAIM_NO kosong -> server 409.
		{"medical tanpa CLAIM_NO", m, "", "QP", services.ErrPindahMenungguNomorKlaim},
		// b22123: hanya Type QP/QR.
		{"medical Type TP", m, "UJI-CLM-1", "TP", services.ErrPindahTypeTanpaMedis},
		{"medical Type TR", m, "UJI-CLM-1", "TR", services.ErrPindahTypeTanpaMedis},
		{"medical QP bernomor", m, "UJI-CLM-1", "QP", nil},
		{"medical QR bernomor", m, "UJI-CLM-1", "qr", nil},
		// b21699: Send Back to Register hanya sebelum CLAIM_NO.
		{"register bernomor", r, "UJI-CLM-1", "QP", services.ErrPindahSudahBernomor},
		{"register belum bernomor", r, "", "QP", nil},
	} {
		err := services.PeriksaPindahDariOutstanding(o, u.ke, u.nomor, u.tipe)
		if (u.mau == nil) != (err == nil) || (u.mau != nil && !errors.Is(err, u.mau)) {
			t.Errorf("%s: galat %v, mau %v", u.nama, err, u.mau)
		}
	}
	// Tahap asal lain tidak digerbangi di sini.
	if err := services.PeriksaPindahDariOutstanding(models.TahapMedicalCheck, models.TahapClaimAnalis, "", "TP"); err != nil {
		t.Errorf("tahap asal Medical ikut digerbangi: %v", err)
	}
}

func TestPindahMemanggilGerbangOutstanding(t *testing.T) {
	isi, err := os.ReadFile("tahap.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	sah, gerbang, tulis := strings.Index(s, "models.SerahTerimaSah(asal, ke)"),
		strings.Index(s, "PeriksaPindahDariOutstanding(asal, ke, hdr.NomorKlaim, hdr.Type)"),
		strings.Index(s, "baca.PerbaruiTahap(")
	if sah < 0 || gerbang < 0 || tulis < 0 || !(sah < gerbang && gerbang < tulis) {
		t.Errorf("gerbang Outstanding tidak berada di antara SerahTerimaSah dan PerbaruiTahap: %d %d %d", sah, gerbang, tulis)
	}
}

func TestIsianOutstanding(t *testing.T) {
	aktif := map[string]bool{"UJI-P1": true, "UJI-P2": true}
	lama := models.Klaim{ID: "UJI-K1", StsKatastrofe: services.NonKatastrofe}
	ok := services.IsianOutstanding{StsKatastrofe: services.Katastrofe, IsKPR: services.KPR, KatastrofeNote: "UJI catatan",
		Accept: map[string]string{"UJI-P1": "true", "UJI-P2": "false"}}
	if err := services.PeriksaIsianOutstanding(ok, lama, aktif); err != nil {
		t.Fatalf("isian sah ditolak: %v", err)
	}
	for _, u := range []struct {
		nama string
		ubah func(*services.IsianOutstanding)
	}{
		{"catastrophe asing", func(i *services.IsianOutstanding) { i.StsKatastrofe = "Banjir" }},
		{"accept bukan true/false", func(i *services.IsianOutstanding) { i.Accept = map[string]string{"UJI-P1": "1"} }},
		{"peserta bukan milik klaim", func(i *services.IsianOutstanding) { i.Accept = map[string]string{"UJI-X": "true"} }},
		{"is kpr terlalu panjang", func(i *services.IsianOutstanding) { i.IsKPR = "123456789" }},
		// Keputusan work owner 11-10-2026: hanya KPR / Non-KPR.
		{"is kpr Yes", func(i *services.IsianOutstanding) { i.IsKPR = "Yes" }},
		{"catatan terlalu panjang", func(i *services.IsianOutstanding) { i.KatastrofeNote = strings.Repeat("x", 1001) }},
	} {
		is := ok
		u.ubah(&is)
		if err := services.PeriksaIsianOutstanding(is, lama, aktif); !errors.Is(err, galat.ErrPermintaanTidakSah) {
			t.Errorf("%s: galat %v, mau ErrPermintaanTidakSah", u.nama, err)
		}
	}
	// b3515 / b5059: sesudah bernomor, Catastrophe dan Note terkunci; Is KPR dan Accept tetap boleh.
	bernomor := models.Klaim{ID: "UJI-K1", NomorKlaim: "UJI-CLM-1", StsKatastrofe: services.NonKatastrofe}
	if err := services.PeriksaIsianOutstanding(ok, bernomor, aktif); !errors.Is(err, services.ErrKatastrofeTerkunci) {
		t.Errorf("Catastrophe berubah sesudah bernomor: %v", err)
	}
	tetap := services.IsianOutstanding{StsKatastrofe: services.NonKatastrofe, IsKPR: services.NonKPR,
		Accept: map[string]string{"UJI-P1": "true"}}
	if err := services.PeriksaIsianOutstanding(tetap, bernomor, aktif); err != nil {
		t.Errorf("Is KPR / Accept sesudah bernomor ditolak: %v", err)
	}
}

// Gerbang dokumen Save to RNM langkah 3 (b1181) membaca IsAccept, bukan IsCheck.
func TestGerbangDokumenMembacaIsAccept(t *testing.T) {
	tanpaDok := []models.Peserta{{IsCheck: "true", IsAccept: "true"}}
	if err := services.PeriksaDokumenAda("QP", tanpaDok); err == nil {
		t.Error("peserta di-Accept tanpa dokumen lolos")
	}
	belumAccept := []models.Peserta{{IsCheck: "true", IsAccept: ""}, {IsCheck: "true", IsAccept: "false"}}
	if err := services.PeriksaDokumenAda("QP", belumAccept); err != nil {
		t.Errorf("peserta yang belum di-Accept ikut dituntut dokumen: %v", err)
	}
}

// Data warisan Is KPR (mis. "Y") tampil apa adanya: Save yang tidak mengubahnya tidak gagal; menggantinya ke nilai asing
// lain tetap ditolak.
func TestIsKPRWarisanTidakMerusakSave(t *testing.T) {
	lama := models.Klaim{ID: "UJI-K1", IsKPR: "Y", StsKatastrofe: services.NonKatastrofe}
	aktif := map[string]bool{}
	tetap := services.IsianOutstanding{StsKatastrofe: services.NonKatastrofe, IsKPR: "Y"}
	if err := services.PeriksaIsianOutstanding(tetap, lama, aktif); err != nil {
		t.Errorf("nilai warisan yang tidak diubah ditolak: %v", err)
	}
	asing := services.IsianOutstanding{StsKatastrofe: services.NonKatastrofe, IsKPR: "N"}
	if err := services.PeriksaIsianOutstanding(asing, lama, aktif); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("nilai asing baru lolos: %v", err)
	}
	for _, v := range []string{services.KPR, services.NonKPR, ""} {
		ganti := services.IsianOutstanding{StsKatastrofe: services.NonKatastrofe, IsKPR: v}
		if err := services.PeriksaIsianOutstanding(ganti, lama, aktif); err != nil {
			t.Errorf("Is KPR %q ditolak: %v", v, err)
		}
	}
}

// 11-10-2026: Claim RNM + bank lewat tombol Save grid - aturan 422 / 409 / 400 tetap.
func TestPeriksaIsianBaris(t *testing.T) {
	m := func(s string) intiuang.Money {
		v, err := intiuang.NewMoney(s, "IDR")
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	baris := models.BarisAdjustment{ID: "UJI-A1", JumlahKlaim: m("10"), ShareNusantaraRe: m("100")}
	ok, err := services.PeriksaIsianBaris(baris, services.IsianBarisAdjustment{Jumlah: "50", NamaBank: " UJI BANK "})
	if err != nil || ok.Jumlah != "50" || ok.NamaBank != "UJI BANK" {
		t.Fatalf("isian sah: %+v, %v", ok, err)
	}
	// 11-10-2026 "Save hanya save progress": Claim RNM > Share TETAP tersimpan; penegakannya di Save to RNM.
	if lebih, err := services.PeriksaIsianBaris(baris, services.IsianBarisAdjustment{Jumlah: "101"}); err != nil || lebih.Jumlah != "101" {
		t.Errorf("melebihi share harus tetap tersimpan: %+v, %v", lebih, err)
	}
	berstatus := baris
	berstatus.KodeStatus = "0"
	if _, err := services.PeriksaIsianBaris(berstatus, services.IsianBarisAdjustment{Jumlah: "50"}); !errors.Is(err, services.ErrAdjustmentTerkunci) {
		t.Errorf("baris berstatus: %v, mau ErrAdjustmentTerkunci (409)", err)
	}
	// 11-10-2026: Claim RNM KOSONG boleh pada Save grid (NULL); wajibnya di Save to RNM.
	if kosong, err := services.PeriksaIsianBaris(baris, services.IsianBarisAdjustment{Jumlah: " ", NamaBank: "UJI BANK"}); err != nil || kosong.Jumlah != "" || kosong.NamaBank != "UJI BANK" {
		t.Errorf("Claim RNM kosong: %+v, %v", kosong, err)
	}
	if _, err := services.PeriksaIsianBaris(berstatus, services.IsianBarisAdjustment{}); !errors.Is(err, services.ErrAdjustmentTerkunci) {
		t.Errorf("kosong pada baris berstatus: %v", err)
	}
	for _, is := range []services.IsianBarisAdjustment{{Jumlah: "-1"}, {Jumlah: "abc"}, {Jumlah: "1", NamaBank: strings.Repeat("x", 256)}} {
		if _, err := services.PeriksaIsianBaris(baris, is); !errors.Is(err, galat.ErrPermintaanTidakSah) {
			t.Errorf("%+v: %v, mau ErrPermintaanTidakSah (400)", is, err)
		}
	}
}

// 11-10-2026 (ORA-12899 di DEV): Catastrophe dan Is KPR dibatasi 20 byte (VARCHAR2(20), migrasi 029) SEBELUM bind; nilai
// sah terpanjang "Non-Catastrophe" (15) lolos.
func TestBatasLebarKatastrofeDanKPR(t *testing.T) {
	lama := models.Klaim{ID: "UJI-K1"}
	if err := services.PeriksaIsianOutstanding(services.IsianOutstanding{StsKatastrofe: services.NonKatastrofe, IsKPR: services.NonKPR}, lama, nil); err != nil {
		t.Errorf("nilai sah ditolak: %v", err)
	}
	for _, is := range []services.IsianOutstanding{
		{StsKatastrofe: strings.Repeat("x", 21)},
		{StsKatastrofe: services.NonKatastrofe, IsKPR: strings.Repeat("y", 21)},
	} {
		if err := services.PeriksaIsianOutstanding(is, lama, nil); !errors.Is(err, galat.ErrPermintaanTidakSah) {
			t.Errorf("%+v: %v, mau ErrPermintaanTidakSah (400)", is, err)
		}
	}
	// Warisan IS_KPR > 20 tidak mungkin tersimpan (kolom 20), jadi jalan pintas "sama dengan tersimpan" tetap aman.
	if len(services.NonKatastrofe) != 15 {
		t.Errorf("Non-Catastrophe = %d byte, mau 15 (b3728)", len(services.NonKatastrofe))
	}
}

// Keputusan work owner 11-10-2026: Admin Notes (dan Location of Death kelak) hanya dapat diubah Admin selama kasus di inbox
// Admin. Save grid menegakkannya SEBELUM menulis apa pun: tahap harus Outstanding (409 ErrIsianBukanOutstanding) dan pelaku
// pemegang tahap itu = ReasLifeAdmin (403).
func TestCatatanAdminHanyaAdminDiTahapOutstanding(t *testing.T) {
	if peran, ada := models.PeranPemegangTahap(models.TahapOutstanding); !ada || peran != "ReasLifeAdmin" {
		t.Fatalf("pemegang tahap Outstanding = %q, mau ReasLifeAdmin", peran)
	}
	isi, err := os.ReadFile("isianoutstanding.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	simpan := s[strings.Index(s, "func (x *SimpanOutstanding) Simpan("):]
	tahap := strings.Index(simpan, "if tahap != models.TahapOutstanding {")
	pemegang := strings.Index(simpan, "WajibPemegangTahap(pelaku, tahap)")
	tulis := strings.Index(simpan, "SetelCatatanAdmin(")
	if tahap < 0 || pemegang < 0 || tulis < 0 || tahap > tulis || pemegang > tulis {
		t.Errorf("gerbang tahap %d / pemegang %d harus mendahului SetelCatatanAdmin %d", tahap, pemegang, tulis)
	}
}

// 11-10-2026: kunci cabangBank / nomorRekening dari klien DIABAIKAN (tanpa medan; Branch of Bank / Account No baca saja).
func TestBankHanyaNameOfBankDariIsian(t *testing.T) {
	var is services.IsianOutstanding
	badan := `{"adjustment":{"UJI-A1":{"jumlah":"5","namaBank":"UJI BANK","cabangBank":"UJI CABANG","nomorRekening":"999"}}}`
	if err := json.Unmarshal([]byte(badan), &is); err != nil {
		t.Fatal(err)
	}
	if got := is.Adjustment["UJI-A1"]; got != (services.IsianBarisAdjustment{Jumlah: "5", NamaBank: "UJI BANK"}) {
		t.Errorf("isian %+v", got)
	}
}

// Save grid = simpan progres: tidak menomori, tidak spreading, tidak outbox, tidak memvalidasi DOL / Share.
func TestSaveGridHanyaSimpanProgres(t *testing.T) {
	isi, err := os.ReadFile("isianoutstanding.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	simpan := s[strings.Index(s, "func (x *SimpanOutstanding) Simpan("):]
	simpan = simpan[:strings.Index(simpan, "\n}\n")]
	for _, dilarang := range []string{"NomorBerikut", "IsiNomorKlaimKosong", "Spreading", "outbox", "Arasapas", "ValidasiDOL(",
		"PeriksaClaimGrossTerhadapShare"} {
		if strings.Contains(simpan, dilarang) {
			t.Errorf("Save grid memanggil %q - harus hanya menyimpan progres", dilarang)
		}
	}
}

// 11-10-2026: proteksi Send to Medical Check - setiap peserta Accept punya >= 1 dokumen; SEMUA kekurangan dikumpulkan;
// peserta tanpa Accept tidak diperiksa; tanpa peserta Accept lolos.
func TestPeriksaDokumenAcceptMedis(t *testing.T) {
	peserta := []models.Peserta{
		{ID: "P1", IsAccept: "true", NamaTertanggung: "UJI-A", NomorSertifikat: "S1"},
		{ID: "P2", IsAccept: "true", NomorSertifikat: "S2"},
		{ID: "P3", IsAccept: "false", NamaTertanggung: "UJI-C"},
		{ID: "P4", IsAccept: "true", NamaTertanggung: "UJI-D"},
	}
	err := services.PeriksaDokumenAcceptMedis(peserta, map[string]int{"P1": 0, "P2": 0, "P4": 2})
	var p *services.PelanggaranPindah
	if !errors.As(err, &p) || len(p.Kekurangan) != 2 {
		t.Fatalf("galat %v", err)
	}
	if p.Kekurangan[0].NomorPeserta != 1 || p.Kekurangan[0].Pesan != "Insured number 1 (UJI-A): at least one claim document is required because the insured is accepted." {
		t.Errorf("#1 %+v", p.Kekurangan[0])
	}
	if p.Kekurangan[1].NomorPeserta != 2 || !strings.Contains(p.Kekurangan[1].Pesan, "Insured number 2 (S2)") {
		t.Errorf("#2 %+v", p.Kekurangan[1])
	}
	if err := services.PeriksaDokumenAcceptMedis(peserta, map[string]int{"P1": 1, "P2": 3, "P4": 1}); err != nil {
		t.Errorf("semua Accept berdokumen: %v", err)
	}
	if err := services.PeriksaDokumenAcceptMedis([]models.Peserta{{ID: "P3", IsAccept: "false"}}, nil); err != nil {
		t.Errorf("tanpa Accept: %v", err)
	}
	isi, err := os.ReadFile("tahap.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	a, b, c := strings.Index(s, "PeriksaPindahDariOutstanding(asal, ke"), strings.Index(s, "baca.SudahSaveRNM(ctx, klaimID)"), strings.Index(s, "PeriksaDokumenAcceptMedis(klaim.Peserta, cacah)")
	if a < 0 || b < a || c < b {
		t.Errorf("urutan gerbang Send to Medical: %d %d %d", a, b, c)
	}
}
