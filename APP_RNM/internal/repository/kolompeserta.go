package repository

// Kolom peserta klaim - satu daftar untuk tulis dan baca.
//
// Pemilik: tiket 02 (pendaftaran menyalin peserta), dipakai tiket 06.
//
// Untuk apa berkas ini: T_CLAIMLF_PREMIUMLIST_DETAIL punya dua puluh tujuh
// kolom yang diisi saat pendaftaran. Menulis daftarnya dua kali - sekali di
// INSERT dan sekali di SELECT - adalah cara paling mudah membuat keduanya
// berselisih diam-diam, dan selisih urutan bind tidak terlihat sampai ada
// nilai yang mendarat di kolom yang salah.
//
// ⭐ Kenapa keempat tanggal valuasi ikut disalin: tiket 06 memvalidasi DOL
// terhadap jendela itu. Menyalinnya saat pendaftaran berarti validasi tidak
// perlu bertanya ulang ke M_LIFE_PREMIUM_DETAIL, yang berisi 66,8 juta baris.
//
// Dibaca sesudah: pohonklaim.go.

import (
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/internal/models"
	"nusantarare/pkg/utils"
)

// golonganKolom menyebut cara satu kolom peserta ditulis dan dibaca.
type golonganKolom int

const (
	kolomTeks    golonganKolom = iota // apa adanya
	kolomTanggal                      // TO_DATE saat tulis, TO_CHAR saat baca
	kolomAngka                        // apa adanya saat tulis, TO_CHAR saat baca
)

// kolomPeserta adalah SATU daftar kolom yang diisi pendaftaran, berurut.
//
// Urutannya mengikat: INSERT dan SELECT keduanya dibangun dari daftar ini.
var kolomPeserta = []struct {
	Nama     string
	Golongan golonganKolom
	// Ambil mengeluarkan nilainya sebagai teks; kosong berarti NULL.
	Ambil func(models.Peserta) string
}{
	{"PL_NUMBER", kolomTeks, func(p models.Peserta) string { return p.NomorPremiList }},
	{"POLICY_NO", kolomTeks, func(p models.Peserta) string { return p.NomorPolis }},
	{"CERTIFICATE_NO", kolomTeks, func(p models.Peserta) string { return p.NomorSertifikat }},
	{"CURRENCY", kolomTeks, func(p models.Peserta) string { return p.MataUang }},
	{"SOURCE_ID", kolomTeks, func(p models.Peserta) string { return p.SumberID }},
	{"IS_CHECK", kolomTeks, func(p models.Peserta) string { return p.IsCheck }},
	{"STS_REJECT", kolomTeks, func(p models.Peserta) string { return p.KodeStatus }},
	{"STNC_TREATY", kolomTeks, func(p models.Peserta) string { return p.STNC }},

	{"GROSS_VALUATION_BEGIN_DATE", kolomTanggal, func(p models.Peserta) string { return p.ValuasiGrossMulai }},
	{"GROSS_VALUATION_EXPIRED_DATE", kolomTanggal, func(p models.Peserta) string { return p.ValuasiGrossSelesai }},
	{"RETRO_VALUATION_BEGIN_DATE", kolomTanggal, func(p models.Peserta) string { return p.ValuasiRetroMulai }},
	{"RETRO_VALUATION_EXPIRED_DATE", kolomTanggal, func(p models.Peserta) string { return p.ValuasiRetroSelesai }},
	{"WPC", kolomTanggal, func(p models.Peserta) string { return p.WPC }},
	{"DATE_OF_LOSS", kolomTanggal, func(p models.Peserta) string { return p.TanggalKejadian }},
	// Tiga tanggal klaim dialog Edit Date (28-09-2026). Pendaftaran
	// menulisnya NULL - tanggal itu lahir sesudah klaim ada; yang
	// mengisinya `PerbaruiTanggalKlaim`.
	{"CLAIM_RECEIVED_DATE", kolomTanggal, func(p models.Peserta) string { return p.TanggalTerimaKlaim }},
	{"COMPLETE_DATE", kolomTanggal, func(p models.Peserta) string { return p.TanggalDokumenLengkap }},
	{"CONFIRMATION_DATE", kolomTanggal, func(p models.Peserta) string { return p.TanggalKonfirmasi }},
	{"BEGIN_DATE", kolomTanggal, func(p models.Peserta) string { return p.TanggalMulai }},
	{"EFFECTIVE_DATE", kolomTanggal, func(p models.Peserta) string { return p.TanggalEfektif }},
	{"LAPSE_DATE", kolomTanggal, func(p models.Peserta) string { return p.TanggalLapse }},
	{"EXPIRED_DATE", kolomTanggal, func(p models.Peserta) string { return p.TanggalExpired }},

	{"SUM_INSURED", kolomAngka, func(p models.Peserta) string { return utils.FormatDecimal(p.SumInsured.Amount) }},
	{"SUM_REASURED", kolomAngka, func(p models.Peserta) string { return utils.FormatDecimal(p.SumReasured.Amount) }},
	{"GROSS_PREMIUM", kolomAngka, func(p models.Peserta) string { return utils.FormatDecimal(p.GrossPremium.Amount) }},
	{"NET_PREMIUM", kolomAngka, func(p models.Peserta) string { return utils.FormatDecimal(p.NetPremium.Amount) }},
	{"CEDING_RETENTION", kolomAngka, func(p models.Peserta) string { return utils.FormatDecimal(p.CedingRetention.Amount) }},
	{"SHARE_NUSANTARA_RE", kolomAngka, func(p models.Peserta) string { return utils.FormatDecimal(p.ShareNusantaraRe.Amount) }},
	{"SHARE_RETRO", kolomAngka, func(p models.Peserta) string { return utils.FormatDecimal(p.ShareRetro.Amount) }},
	{"RETROCEDED_SHARE", kolomAngka, func(p models.Peserta) string { return utils.FormatDecimal(p.RetrocededShare.Amount) }},
	{"EM_PERCENT", kolomAngka, func(p models.Peserta) string { return utils.FormatDecimal(p.EMPercent.Value) }},

	// ⭐ AGE kolomnya sudah ada di migrasi 003 sejak awal, tetapi tidak
	// pernah terisi karena kolomSalin tidak membacanya. Diisi sejak
	// lanjutan 10 §1: `SaveInsuredClaim_Act` memilihnya, jadi ia memang
	// bagian pendaftaran.
	{"AGE", kolomAngka, func(p models.Peserta) string { return p.Umur }},
}

// insertPeserta menyusun pernyataan INSERT beserta nilainya.
//
// ID dan CLAIM_ID ditulis lebih dulu; sisanya datang dari kolomPeserta.
func insertPeserta(tabel, id, klaimID string, p models.Peserta) (string, []any) {
	nama := []string{"ID", "CLAIM_ID"}
	penampung := []string{":1", ":2"}
	nilai := []any{id, klaimID}

	for _, k := range kolomPeserta {
		n := len(nilai) + 1
		nama = append(nama, k.Nama)
		if k.Golongan == kolomTanggal {
			// Bentuk tanggalnya sama persis dengan yang dipakai saat membaca,
			// sehingga tulis dan baca bertemu di satu bentuk dan tidak satu
			// pun bergantung pada NLS_DATE_FORMAT sesi.
			penampung = append(penampung, fmt.Sprintf(`TO_DATE(:%d, 'YYYY-MM-DD HH24:MI:SS')`, n))
		} else {
			penampung = append(penampung, fmt.Sprintf(":%d", n))
		}
		nilai = append(nilai, kosongJadiNil(normalTanggal(k)(p)))
	}
	return fmt.Sprintf("INSERT INTO %s\n\t\t\t(%s)\n\t\t\tVALUES (%s)",
		tabel, strings.Join(nama, ", "), strings.Join(penampung, ",")), nilai
}

// normalTanggal menormalkan nilai tanggal ke satu bentuk sebelum TO_DATE.
//
// Tanpa ini, TO_DATE menemui bentuk yang tidak dikenalnya dan menjawab galat
// yang menuding format, bukan kolomnya.
func normalTanggal(k struct {
	Nama     string
	Golongan golonganKolom
	Ambil    func(models.Peserta) string
}) func(models.Peserta) string {
	if k.Golongan != kolomTanggal {
		return k.Ambil
	}
	return func(p models.Peserta) string {
		teks := strings.TrimSpace(k.Ambil(p))
		if teks == "" {
			return ""
		}
		t, err := utils.ParseTanggal(teks)
		if err != nil {
			// Dibiarkan apa adanya: TO_DATE yang akan menolaknya, dan
			// galatnya menyebut nilainya. Menelannya di sini menghilangkan
			// satu-satunya petunjuk.
			return teks
		}
		return utils.FormatTanggal(t)
	}
}

// selectPeserta menyusun daftar ekspresi SELECT untuk kolomPeserta.
//
// Kolom angka dan tanggal dibungkus TO_CHAR supaya tiba sebagai TEKS dengan
// bentuk yang tidak bergantung setelan sesi (ADR-U-0003, ADR-U-0016).
func selectPeserta() string {
	ekspresi := make([]string, 0, len(kolomPeserta))
	for _, k := range kolomPeserta {
		switch k.Golongan {
		case kolomTanggal:
			ekspresi = append(ekspresi, fmt.Sprintf(fmtTanggalOracle, k.Nama))
		case kolomAngka:
			ekspresi = append(ekspresi, fmt.Sprintf(fmtDesimal, k.Nama))
		default:
			ekspresi = append(ekspresi, k.Nama)
		}
	}
	return strings.Join(ekspresi, ", ")
}

// rakitPeserta menyusun satu peserta dari hasil SELECT.
//
// Urutannya mengikuti kolomPeserta, sama dengan yang dipakai insertPeserta -
// itulah sebabnya keduanya membaca daftar yang sama.
//
// ⚠️ Uang yang gagal diurai DILAPORKAN, tidak ditebak dan tidak didiamkan:
// angka yang salah jauh lebih berbahaya daripada pembacaan yang gagal
// (ADR-U-0003).
func rakitPeserta(id string, sel []sql.NullString) (models.Peserta, error) {
	if len(sel) != len(kolomPeserta) {
		return models.Peserta{}, fmt.Errorf(
			"repository: peserta %s terbaca %d kolom, daftar menyebut %d",
			id, len(sel), len(kolomPeserta))
	}
	p := models.Peserta{ID: id, Baris: []models.BarisAdjustment{}}
	nilai := map[string]sql.NullString{}
	for i, k := range kolomPeserta {
		nilai[k.Nama] = sel[i]
	}
	teks := func(n string) string { return nilai[n].String }

	p.NomorPremiList = teks("PL_NUMBER")
	p.NomorPolis = teks("POLICY_NO")
	p.NomorSertifikat = teks("CERTIFICATE_NO")
	p.MataUang = teks("CURRENCY")
	p.SumberID = teks("SOURCE_ID")
	p.IsCheck = teks("IS_CHECK")
	p.KodeStatus = teks("STS_REJECT")
	p.STNC = teks("STNC_TREATY")
	p.Umur = teks("AGE")

	p.ValuasiGrossMulai = teks("GROSS_VALUATION_BEGIN_DATE")
	p.ValuasiGrossSelesai = teks("GROSS_VALUATION_EXPIRED_DATE")
	p.ValuasiRetroMulai = teks("RETRO_VALUATION_BEGIN_DATE")
	p.ValuasiRetroSelesai = teks("RETRO_VALUATION_EXPIRED_DATE")
	p.WPC = teks("WPC")
	p.TanggalKejadian = teks("DATE_OF_LOSS")
	p.TanggalTerimaKlaim = teks("CLAIM_RECEIVED_DATE")
	p.TanggalDokumenLengkap = teks("COMPLETE_DATE")
	p.TanggalKonfirmasi = teks("CONFIRMATION_DATE")
	p.TanggalMulai = teks("BEGIN_DATE")
	p.TanggalEfektif = teks("EFFECTIVE_DATE")
	p.TanggalLapse = teks("LAPSE_DATE")
	p.TanggalExpired = teks("EXPIRED_DATE")

	uang := []struct {
		kolom string
		ke    *models.Money
	}{
		{"SUM_INSURED", &p.SumInsured}, {"SUM_REASURED", &p.SumReasured},
		{"GROSS_PREMIUM", &p.GrossPremium}, {"NET_PREMIUM", &p.NetPremium},
		{"CEDING_RETENTION", &p.CedingRetention},
		{"SHARE_NUSANTARA_RE", &p.ShareNusantaraRe},
		{"SHARE_RETRO", &p.ShareRetro}, {"RETROCEDED_SHARE", &p.RetrocededShare},
	}
	for _, u := range uang {
		m, err := uraiUang(id, u.kolom, nilai[u.kolom], p.MataUang)
		if err != nil {
			return models.Peserta{}, err
		}
		*u.ke = m
	}
	rasio, err := uraiRasio(id, "EM_PERCENT", nilai["EM_PERCENT"])
	if err != nil {
		return models.Peserta{}, err
	}
	p.EMPercent = rasio
	return p, nil
}
