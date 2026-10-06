package models

// Validasi unggahan CSV premium list - tiket 04 PremiumList Life.
//
// Untuk apa berkas ini: memeriksa satu berkas CSV baris demi baris dan
// mengumpulkan SELURUH penolakannya. Murni - nol Oracle, nol berkas, nol jam.
//
// Sumbernya `Activity/ValidasiUploadPL_act.xml`, dibaca sebagai pohon
// 28-09-2026.
//
// `[terverifikasi]` 43 langkah tingkat pertama. Jendela: elemen
// `<pyStepPageReference>` berbentuk `RH_1.pySteps(N)`, tingkat SATU saja.
// Dua cara, keduanya menjawab 43:
//
//	grep -c "<pyStepPageReference>RH_1\.pySteps([0-9]*)</pyStepPageReference>" \
//	  Activity/ValidasiUploadPL_act.xml
//	grep -o "...\(sama\)..." | grep -o "[0-9]*" | sort -n | tail -1
//
// Peta langkahnya: 1 menyiapkan pesan · 2 MENGISI 0 bagi 32 kolom uang yang
// kosong (tiap sub-langkah `.X = 0` bila `!@PropertyHasValue(.X)`, mis. 2.1
// b2146-2223) · 3-8 memeriksa rujukan master · 9 memeriksa SETIAP baris ·
// 14-43 menempelkan pesannya.
//
// ⛔ RALAT 28-09-2026 (sensus remark GILIRAN-12): langkah 9.2-9.5, 9.11, 10,
// 11, 12, 13, 19, dan 33 ter-remark (`//` b8025, b8198, b8368, b8517, b9667,
// b13018, b13171, b13322, b13475, b14393, b16535). Ronde pertama menegakkan
// empat penolakan dari langkah itu - sertifikat sudah dipakai, sertifikat
// ganda dalam berkas, `MEDICAL_STATUS` harus FCL/M/NM, dan lampiran wajib -
// ditambah duplikat peserta nama+DOB yang meminjam kalimat `local.err1`.
// Kelimanya DIBUANG: sistem lama tidak pernah menolak karenanya.
// (`CekDoubleInsured` memang hidup, tetapi di `Calculate1_Act` 7.5 untuk
// akumulasi retensi ceding - bukan penolakan unggah.)
//
// ⛔ OQ-PL-12 DITUTUP 29-09-2026 (GILIRAN-17) `[keputusan work owner]`: seperti
// langkah 2, kolom uang yang KOSONG menjadi 0 (`IsiNolUangKosong`) SEBELUM
// validasi, sehingga "HARUS ADA" langkah 9.12-9.17 tidak berbunyi atas sel
// kosong - persis sistem lama. Nilai 0 itu ikut tersimpan, seperti Pega
// (`SavePremiumList_Act` langkah 8 mengulang daftar yang sama).
//
// ⛔ PESAN DISALIN VERBATIM, termasuk yang menjanjikan hal yang tidak
// diperiksa. Lihat `PesanNetPremium` dan OQ-069.
//
// ⛔ ⚠️ NORMALISASI KOMA PEGA TIDAK DITIRU. Langkah "Rubah decimal dari koma
// jadi titik" memakai `@replaceAll(.KOLOM, ",", ".")` atas tiap kolom uang -
// MENGGANTI SETIAP KOMA, tanpa membedakan pemisah desimal dari pemisah
// ribuan. `1,234,567.89` menjadi `1.234.567.89`, yang bukan angka sama
// sekali, dan hasilnya tersimpan sebagai uang. Kami MENOLAK komanya alih-alih
// menggantinya - lihat `UangCSV`. ⚠️ Pengecualian sempit sejak 02-10-2026:
// berkas berpemisah `;` (Excel lokal Indonesia) memakai koma DESIMAL, dan
// SATU koma itu diubah menjadi titik - lihat `DesimalKomaKeTitik`.
//
// Dibaca sesudah: polis_nomor.go, polis_detail.go.

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
)

// Pesan penolakan - VERBATIM `local.err1`..`local.err32`.
//
// ⛔ DISALIN APA ADANYA, termasuk ejaan dan huruf besarnya. Orang yang sudah
// bertahun-tahun memakai layar lama mengenali kalimat ini; menerjemahkannya
// membuat mereka mengira ada galat jenis baru.
const (
	PesanNamaTertanggung  = "NAME OF INSURED HARUS ADA"
	PesanDOB              = "DOB HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanPlan             = "PLAN HARUS ADA"
	PesanBeginDate        = "BEGIN DATE HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanExpiredDate      = "EXPIRED DATE HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanSumInsured       = "SUM INSURED HARUS ADA"
	PesanCedingRetention  = "CEDING RETENTION HARUS ADA"
	PesanSumReasured      = "SUM REASURED HARUS ADA"
	PesanShareNusantaraRe = "SHARE NUSANTARA RE HARUS ADA"
	PesanGrossPremium     = "GROSS PREMIUM HARUS ADA"
	PesanCurrency         = "CURRENCY HARUS ADA"
	PesanStartDate        = "START_DATE HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanEffectiveDate    = "EFFECTIVE_DATE HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanSTNC             = "STNC HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanWPC              = "WPC HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanGrossValMulai    = "GROSS_VALUATION_BEGIN_DATE HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanGrossValSelesai  = "GROSS_VALUATION_EXPIRED_DATE HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanRetroValMulai    = "RETROCESSION_VALUATION_BEGIN_DATE HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanRetroValSelesai  = "RETROCESSION_VALUATION_EXPIRED_DATE HARUS ADA. FORMAT:dd/mm/yyyy"
	PesanPolicyNo         = "POLICY_NO HARUS ADA"
	PesanCertificateNo    = "CERTIFICATE_NO HARUS ADA"
)

// PesanNetPremium disalin VERBATIM, dan ia MENJANJIKAN ATURAN YANG TIDAK ADA.
//
// ⛔ OQ-069 TERBUKA. Kalimatnya berkata net premium harus "LEBIH BESAR DARI
// GROSS PREMIUM", tetapi satu-satunya precondition atas `NET_PREMIUM` di
// SELURUH `ValidasiUploadPL_act.xml` adalah `@PropertyHasValue(.NET_PREMIUM)`
// - nol perbandingan terhadap `GROSS_PREMIUM`. Arahnya pun janggal: lazimnya
// net lebih KECIL dari gross.
//
// ⛔ SELAMA OQ-069 TERBUKA, YANG DIPERIKSA HANYA KEBERADAANNYA - persis
// korpus. Mengarang perbandingannya akan menolak berkas yang di sistem lama
// diterima, dan arah yang salah akan menolak SETIAP berkas yang benar.
const PesanNetPremium = "NET PREMIUM HARUS ADA DAN LEBIH BESAR DARI GROSS PREMIUM"

// Pesan rujukan master - langkah 3-8.
const (
	PesanCedingCo     = "CEDING CO TIDAK TERDAFTAR"
	PesanPolicyHolder = "POLICY HOLDER TIDAK TERDAFTAR"
	PesanMarketing    = "MARKETING OFFICER TIDAK TERDAFTAR"
	PesanKelasBisnis  = "CLASS OF BUSINESS BELUM DIISI"
	PesanRetro        = "RETRO HARUS DIISI"
	PesanSOB          = "SOURCE OF BUSINESS TIDAK TERDAFTAR"
)

// Galat penguraian nilai CSV.
var (
	// ErrUangBerkoma - nilai uang memuat koma.
	//
	// ⛔ DITOLAK, BUKAN DIGANTI. `ValidasiUploadPL_act` menyatakan aturannya
	// enam kali di nama langkahnya ("SEPARATOR MENGGUNAKAN TITIK") dan
	// MENEGAKKANNYA hanya sekali, atas `GROSS_PREMIUM`
	// (`@contains(.GROSS_PREMIUM,",")`). Kami menegakkannya untuk SELURUH
	// kolom uang: aturan yang dinyatakan enam kali lalu dijaga sekali adalah
	// aturan yang dilanggar lima kali tanpa ada yang tahu.
	ErrUangBerkoma = errors.New(
		"models: nilai uang memuat koma; pemisah desimal HARUS titik dan " +
			"pemisah ribuan tidak diterima")
	// ⚠️ `ErrUangKosong` TIDAK dideklarasikan lagi di sini: `money.go` sudah
	// punya, dan artinya sama persis. Dua galat bernama sama untuk satu
	// keadaan membuat `errors.Is` pemanggil menjawab benar hanya untuk salah
	// satunya - dan pemanggil tidak punya cara menebak yang mana.
	// ErrUangTakTerurai - bukan bilangan desimal.
	ErrUangTakTerurai = errors.New("models: nilai uang tidak dapat diurai")
	// ErrTanggalTakBerbentuk - bukan `dd/mm/yyyy`.
	ErrTanggalTakBerbentuk = errors.New("models: tanggal bukan dd/mm/yyyy")
)

// BentukTanggalCSV adalah bentuk tanggal CSV - `dd/mm/yyyy`.
//
// `[terverifikasi]` Keduapuluh pesan tanggal menyebutnya harfiah
// (`FORMAT:dd/mm/yyyy`), dan pemeriksanya `@toDate(.X)!=0`.
const BentukTanggalCSV = "02/01/2006"

// TanggalEpochCSV adalah tanggal yang Pega KECUALIKAN secara khusus.
//
// ⛔ `@if(.DOB=="02/01/1970",true,@toDate(.DOB)!=0)` - langkah 9.
//
// Sebabnya: `@toDate` Pega mengembalikan NOL untuk tanggal yang setara epoch,
// dan pemeriksanya membandingkan dengan nol. Tanggal lahir 2 Januari 1970
// karena itu akan tertolak sebagai "tidak berbentuk" padahal ia sah - dan
// seseorang sudah pernah menambalnya satu per satu.
//
// ⚠️ Kami TIDAK punya cacat itu: pengurai Go tidak memakai nol sebagai
// penanda gagal. Pengecualiannya tetap dicatat di sini supaya pembaca
// berikutnya tahu kenapa korpus memuat baris yang tampak aneh - dan supaya
// tidak ada yang "merapikannya" menjadi aturan yang menolak tanggal itu.
const TanggalEpochCSV = "02/01/1970"

// UangCSV mengurai satu nilai uang dari CSV.
//
// ⛔ PEMISAHNYA DINYATAKAN, TIDAK DITEBAK (AC tiket 04): desimal memakai
// TITIK, dan pemisah ribuan TIDAK diterima sama sekali. `1,234,567.89` dan
// `1.234.567,89` keduanya DITOLAK - bukan diterima diam-diam sebagai angka
// yang mungkin sama, mungkin tidak.
//
// ⛔ NOL FLOAT (ADR-U-0003, ADR-U-0016). Teksnya diurai langsung menjadi
// desimal presisi arbitrer.
//
// ⚠️ Kosong DIBEDAKAN dari nol: pemanggil yang mewajibkan kolomnya memeriksa
// `ErrUangKosong` sendiri.
func UangCSV(teks string) (*apd.Decimal, error) {
	s := strings.TrimSpace(teks)
	if s == "" {
		return nil, uang.ErrUangKosong
	}
	// ⛔ Koma diperiksa LEBIH DAHULU, sebelum penguraian. Pengurai desimal
	// akan menolak `1,5` dengan pesannya sendiri - pesan yang tidak memberi
	// tahu orang bahwa masalahnya pemisah, bukan angkanya.
	if strings.ContainsRune(s, ',') {
		return nil, fmt.Errorf("%w: %q", ErrUangBerkoma, teks)
	}
	d, err := utils.ParseDecimal(s)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrUangTakTerurai, teks)
	}
	return d, nil
}

// TanggalCSV mengurai satu tanggal `dd/mm/yyyy`.
//
// ⛔ KETAT. `time.Parse` dengan bentuk ini menolak `2026-01-02` maupun
// `1/2/2026`, dan itu yang diinginkan: tanggal yang terbaca dua cara adalah
// tanggal yang suatu hari terbaca cara yang salah.
func TanggalCSV(teks string) (time.Time, error) {
	s := strings.TrimSpace(teks)
	if s == "" {
		return time.Time{}, ErrTanggalTakBerbentuk
	}
	t, err := time.Parse(BentukTanggalCSV, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %q", ErrTanggalTakBerbentuk, teks)
	}
	return t, nil
}

// KolomUangUnggah adalah kolom uang yang divalidasi, urut dokumen.
//
// `[terverifikasi]` 32 kolom. Jendela: precondition `@PropertyHasValue(.X)`
// milik langkah 2 di `Activity/ValidasiUploadPL_act.xml`, yaitu rentang
// antara `RH_1.pySteps(2)` dan `RH_1.pySteps(3)`.
//
// ⚠️ DUA CARA, DAN KEDUANYA BERBEDA - selisihnya disebut, tidak didiamkan:
//
//	cara A - cacah precondition `@PropertyHasValue(.X)` di langkah 2  -> 32
//	cara B - cacah `<pyStepsDescription>` sub-langkahnya yang berisi
//	         nama kolom huruf besar                                  -> 31
//
// Selisihnya SATU dan bernama: `FLEET_DISCOUNT` punya precondition tetapi
// deskripsi sub-langkahnya KOSONG. Yang benar 32 - precondition-lah yang
// BERJALAN; deskripsi hanya nama yang dibaca manusia. Sensus yang memakai
// cara B akan menghilangkan satu kolom uang dari validasi, dan kolom uang
// yang tidak divalidasi adalah kolom uang yang menerima apa saja.
//
// ⚠️ RALAT TIKET: tiket 04 menulis "33 kolom uang". Yang benar 32.
var KolomUangUnggah = []string{
	"NET_PREMIUM", "GROSS_PREMIUM", "SHARE_NUSANTARA_RE", "SUM_INSURED",
	"CEDING_RETENTION", "SUM_REASURED", "CLAIM", "TAX", "BROKERAGE_FEE",
	"OVR_COMM", "PROF_COMM", "EM_PERCENT", "COMM", "FLEET_DISCOUNT",
	"GROSS_PREMIUM_REFUND", "NET_PREMIUM_REFUND", "COMM_REFUND",
	"BROKERAGE_FEE_REFUND", "OVR_COMM_REFUND", "TAX_REFUND", "SHARE_RETRO",
	"GROSS_PREMIUM_RETRO", "DISCOUNT_PREMIUM_RETRO", "OVR_COMM_RETRO",
	"BROKERAGE_FEE_RETRO", "NET_PREMIUM_RETRO", "GROSS_PREMIUM_REFUND_RETRO",
	"DISCOUNT_PREMIUM_REFUND_RETRO", "OVR_COMM_REFUND_RETRO",
	"BROKERAGE_FEE_REFUND_RETRO", "NET_PREMIUM_REFUND_RETRO", "CLAIM_AMOUNT",
}

// IsiNolUangKosong - `ValidasiUploadPL_act` langkah 2 (2.1 b2121 ... 2.32
// b6691, `pyStepsBlockName` kosong): tiap kolom `KolomUangUnggah` yang
// `!@PropertyHasValue` disetel `0`. Kolom yang tidak ada di berkas pun kosong
// menurut Pega. OQ-PL-12 (GILIRAN-17).
//
// ⛔ KECUALI kolom uang yang WAJIB (`kolomUangWajib`, keputusan work owner
// 03-10-2026): sel kosongnya harus DITOLAK "HARUS ADA", dan mengisinya 0
// lebih dahulu akan membuat penolakan itu tidak pernah berbunyi.
//
// ⚠️ Ini PENAFSIRAN MASUKAN (sel kosong berarti nol), bukan pengisian tabel:
// ADR-U-0027 tetap berlaku untuk kolom yang memang tidak bernilai.
func IsiNolUangKosong(baris []BarisUnggah) {
	for i := range baris {
		if baris[i].Nilai == nil {
			baris[i].Nilai = map[string]string{}
		}
		for _, k := range KolomUangUnggah {
			if kolomUangWajib[k] {
				continue
			}
			if strings.TrimSpace(baris[i].Nilai[k]) == "" {
				baris[i].Nilai[k] = "0"
			}
		}
	}
}

// ——— Kolom wajib per Type ———
//
// ⛔ [keputusan work owner 03-10-2026] Daftar di bawah MENGGANTI kolom wajib
// `ValidasiUploadPL_act` langkah 9. Yang wajib kini bergantung pada `Type`
// polis: QR/QP memakai Gross Valuation, TR/TP memakai Retrocession
// Valuation. Kolom lain BOLEH kosong, tetapi bila terisi bentuknya tetap
// diperiksa (tanggal dd/mm/yyyy, angka ber-titik, bilangan bulat).
//
// ⚠️ PLAN, WPC, CEDING_RETENTION, SUM_REASURED, SHARE_NUSANTARA_RE,
// GROSS_PREMIUM, dan NET_PREMIUM - dulu wajib - TIDAK lagi wajib. Kolom uang
// yang kosong tetap diisi 0 (`IsiNolUangKosong`).

// jenisKolomUnggah menentukan pemeriksaan BENTUK satu kolom.
type jenisKolomUnggah int

const (
	jenisTeks jenisKolomUnggah = iota
	jenisTanggal
	jenisBulat
	jenisUang
)

// KolomWajibUnggahan adalah satu kolom wajib beserta pesan dan jenisnya.
type KolomWajibUnggahan struct {
	Kolom string
	Pesan string
	jenis jenisKolomUnggah
}

// PesanEntryAge - karangan layar baru berpola `<KOLOM> HARUS ADA` milik
// korpus; ENTRY_AGE tidak pernah diwajibkan `ValidasiUploadPL_act`.
const PesanEntryAge = "ENTRY_AGE HARUS ADA"

// ErrTipeUnggahTakDikenal - Type polis kosong atau di luar QR/QP/TP/TR, jadi
// kolom wajibnya tidak dapat ditentukan (keputusan work owner 03-10-2026:
// unggahan DITOLAK, bukan divalidasi dengan aturan tebakan).
var ErrTipeUnggahTakDikenal = errors.New(
	"choose the Type (QR, QP, TP or TR) and press Save Data before uploading the CSV")

// kolomUangWajib - kolom uang yang wajib di SEMUA Type; dikecualikan dari
// `IsiNolUangKosong`.
var kolomUangWajib = map[string]bool{"SUM_INSURED": true}

// KolomWajibPerTipe mengembalikan kolom wajib unggahan untuk satu Type, urut
// seperti daftar work owner.
func KolomWajibPerTipe(tipe string) ([]KolomWajibUnggahan, error) {
	var valuasi []KolomWajibUnggahan
	switch strings.TrimSpace(tipe) {
	case "QR", "QP":
		valuasi = []KolomWajibUnggahan{
			{"GROSS_VALUATION_BEGIN_DATE", PesanGrossValMulai, jenisTanggal},
			{"GROSS_VALUATION_EXPIRED_DATE", PesanGrossValSelesai, jenisTanggal},
		}
	case "TR", "TP":
		valuasi = []KolomWajibUnggahan{
			{"RETROCESSION_VALUATION_BEGIN_DATE", PesanRetroValMulai, jenisTanggal},
			{"RETROCESSION_VALUATION_EXPIRED_DATE", PesanRetroValSelesai, jenisTanggal},
		}
	default:
		if strings.TrimSpace(tipe) == "" {
			return nil, ErrTipeUnggahTakDikenal
		}
		return nil, fmt.Errorf("%w (Type %q)", ErrTipeUnggahTakDikenal, tipe)
	}
	k := []KolomWajibUnggahan{
		{"POLICY_NO", PesanPolicyNo, jenisTeks},
		{"CERTIFICATE_NO", PesanCertificateNo, jenisTeks},
		{"NAME_OF_INSURED", PesanNamaTertanggung, jenisTeks},
		{"DOB", PesanDOB, jenisTanggal},
		{"ENTRY_AGE", PesanEntryAge, jenisBulat},
		{"BEGIN_DATE", PesanBeginDate, jenisTanggal},
		{"EXPIRED_DATE", PesanExpiredDate, jenisTanggal},
	}
	k = append(k, valuasi...)
	return append(k, []KolomWajibUnggahan{
		{"CURRENCY", PesanCurrency, jenisTeks},
		{"SUM_INSURED", PesanSumInsured, jenisUang},
	}...), nil
}

// kolomTanggalUnggah - SELURUH kolom tanggal unggahan. Yang tidak wajib bagi
// Type-nya tetap diperiksa BENTUKNYA bila terisi.
//
// `[terverifikasi]` langkah 9: tiap kolom punya sepasang precondition -
// `@PropertyHasValue(.X)` lalu `@toDate(.X)!=0`.
var kolomTanggalUnggah = []struct{ Kolom, Pesan string }{
	{"DOB", PesanDOB},
	{"BEGIN_DATE", PesanBeginDate},
	{"EXPIRED_DATE", PesanExpiredDate},
	{"STNC", PesanSTNC},
	{"WPC", PesanWPC},
	{"START_DATE", PesanStartDate},
	{"EFFECTIVE_DATE", PesanEffectiveDate},
	{"GROSS_VALUATION_BEGIN_DATE", PesanGrossValMulai},
	{"GROSS_VALUATION_EXPIRED_DATE", PesanGrossValSelesai},
	{"RETROCESSION_VALUATION_BEGIN_DATE", PesanRetroValMulai},
	{"RETROCESSION_VALUATION_EXPIRED_DATE", PesanRetroValSelesai},
}

// KolomShareNusantaraRe dan KolomShareNusantaraReGross - dua nama kolom CSV
// untuk Share Nusantara Re.
//
// ⛔ [keputusan work owner 02-10-2026] `SHARE_NUSANTARA_RE_GROSS` DIDAHULUKAN;
// bila kolom itu tidak ada atau kosong, `SHARE_NUSANTARA_RE` yang dipakai.
// Keduanya berakhir di `SHARE_NUSANTARA_RE` - lihat `TerapkanAliasShare`.
const (
	KolomShareNusantaraRe      = "SHARE_NUSANTARA_RE"
	KolomShareNusantaraReGross = "SHARE_NUSANTARA_RE_GROSS"
)

// TerapkanAliasShare mengisi `SHARE_NUSANTARA_RE` dari
// `SHARE_NUSANTARA_RE_GROSS` bila yang terakhir terisi.
func TerapkanAliasShare(nilai map[string]string) {
	if v := strings.TrimSpace(nilai[KolomShareNusantaraReGross]); v != "" {
		nilai[KolomShareNusantaraRe] = v
	}
}

// DesimalKomaKeTitik mengubah SATU koma desimal menjadi titik: `10,5` →
// `10.5`.
//
// ⛔ HANYA untuk berkas berpemisah titik koma (`;`) - bentuk simpan Excel
// berlokal Indonesia, yang memakai koma sebagai pemisah DESIMAL
// (keputusan work owner 02-10-2026). Di berkas berpemisah koma, koma di
// dalam angka tetap DITOLAK `UangCSV`: di sana `1,234` bisa berarti seribu.
//
// ⛔ BUKAN `@replaceAll` Pega (lihat kepala berkas): nilai yang memuat titik
// DAN koma (`1.234,56`), atau lebih dari satu koma, DIBIARKAN apa adanya,
// sehingga `UangCSV` menolaknya - pemisah ribuan tetap tidak diterima.
func DesimalKomaKeTitik(teks string) string {
	s := strings.TrimSpace(teks)
	if strings.Count(s, ",") != 1 || strings.ContainsRune(s, '.') {
		return teks
	}
	return strings.Replace(s, ",", ".", 1)
}

// NormalisasiDesimalKoma menerapkan `DesimalKomaKeTitik` pada setiap kolom
// uang satu baris. Hanya untuk berkas berpemisah `;`.
func NormalisasiDesimalKoma(nilai map[string]string) {
	for _, k := range KolomUangTersimpan() {
		if v, ada := nilai[k]; ada {
			nilai[k] = DesimalKomaKeTitik(v)
		}
	}
}

// BarisUnggah adalah satu baris CSV, apa adanya.
type BarisUnggah struct {
	// Nomor adalah nomor baris CSV yang DILIHAT ORANG - 1 untuk baris data
	// pertama, dengan baris judul TIDAK dihitung.
	//
	// ⛔ Nomor baris ikut ke setiap penolakan (AC tiket 04). Penolakan tanpa
	// nomor baris memaksa orang mencocokkan pesan dengan berkas ratusan baris
	// dengan mata.
	Nomor int
	// Nilai berkunci NAMA KOLOM CSV, huruf besar.
	Nilai map[string]string
}

// Penolakan adalah satu alasan sebuah baris ditolak.
type Penolakan struct {
	Baris int    `json:"baris"`
	Kolom string `json:"kolom"`
	// Pesan VERBATIM dari korpus - yang dikenali pemakai lama.
	Pesan string `json:"pesan"`
	// Sebab adalah keterangan KAMI, tepat pada apa yang salah.
	//
	// ⛔ ADA DUA, dan keduanya perlu. Pesan verbatim menjaga pengenalan;
	// sebab menjaga kebenaran. `SUM INSURED HARUS ADA` yang muncul untuk
	// nilai `1,000` berbohong tentang penyebabnya - dan pesan yang berbohong
	// mengajari orang mengabaikan pesan.
	Sebab string `json:"sebab"`
}

// HasilUnggah adalah hasil pemeriksaan satu berkas.
type HasilUnggah struct {
	// CacahBaris adalah cacah SELURUH baris data yang terbaca.
	CacahBaris int `json:"cacahBaris"`
	// Ditolak memuat SELURUH penolakan, bukan yang pertama saja.
	//
	// ⛔ Berhenti di penolakan pertama memaksa orang mengunggah ulang
	// sebanyak jumlah kesalahannya. Berkas dengan dua puluh baris rusak
	// berarti dua puluh putaran.
	Ditolak []Penolakan `json:"ditolak"`
}

// Lolos menjawab apakah berkasnya boleh disimpan permanen.
func (h HasilUnggah) Lolos() bool { return len(h.Ditolak) == 0 }

// BarisDitolak mengembalikan nomor baris yang punya penolakan, unik dan urut.
func (h HasilUnggah) BarisDitolak() []int {
	lihat := map[int]bool{}
	keluar := []int{}
	for _, p := range h.Ditolak {
		if !lihat[p.Baris] {
			lihat[p.Baris] = true
			keluar = append(keluar, p.Baris)
		}
	}
	return keluar
}

// ValidasiUnggah memeriksa seluruh baris dan mengumpulkan setiap penolakan,
// dengan kolom wajib milik `tipe` (`KolomWajibPerTipe`).
//
// ⛔ SELURUH BARIS DIPERIKSA, bahkan sesudah satu ditolak - lihat
// `HasilUnggah.Ditolak`.
func ValidasiUnggah(tipe string, baris []BarisUnggah) (HasilUnggah, error) {
	wajib, err := KolomWajibPerTipe(tipe)
	if err != nil {
		return HasilUnggah{}, err
	}
	hasil := HasilUnggah{CacahBaris: len(baris), Ditolak: []Penolakan{}}
	hitungQR := HitungPesertaPerTipe(tipe)
	tolak := func(b int, kolom, pesan, sebab string) {
		hasil.Ditolak = append(hasil.Ditolak, Penolakan{
			Baris: b, Kolom: kolom, Pesan: pesan, Sebab: sebab})
	}
	// periksaBentuk menolak nilai TERISI yang bentuknya salah.
	periksaBentuk := func(b int, kolom, pesan string, jenis jenisKolomUnggah, v string) {
		switch jenis {
		case jenisTanggal:
			if _, err := TanggalCSV(v); err != nil {
				tolak(b, kolom, pesan, fmt.Sprintf("%q bukan dd/mm/yyyy", v))
			}
		case jenisBulat:
			if !polaBulat.MatchString(v) {
				tolak(b, kolom, pesan, fmt.Sprintf("%q bukan bilangan bulat", v))
			}
		case jenisUang:
			if _, err := UangCSV(v); err != nil {
				sebab := fmt.Sprintf("%q bukan desimal ber-pemisah titik", v)
				if errors.Is(err, ErrUangBerkoma) {
					sebab = fmt.Sprintf("%q memuat koma; pemisah desimal harus "+
						"titik dan pemisah ribuan tidak diterima", v)
				}
				tolak(b, kolom, pesan, sebab)
			}
		}
	}

	for _, b := range baris {
		ambil := func(k string) string { return strings.TrimSpace(b.Nilai[k]) }
		sudah := map[string]bool{}

		// 1. Kolom WAJIB Type ini - kosong ditolak, terisi diperiksa bentuknya.
		for _, w := range wajib {
			sudah[w.Kolom] = true
			v := ambil(w.Kolom)
			if v == "" {
				tolak(b.Nomor, w.Kolom, w.Pesan, "kolom kosong")
				continue
			}
			periksaBentuk(b.Nomor, w.Kolom, w.Pesan, w.jenis, v)
		}
		// 2. Tanggal lain - hanya bentuknya, bila terisi.
		for _, t := range kolomTanggalUnggah {
			if v := ambil(t.Kolom); !sudah[t.Kolom] && v != "" {
				periksaBentuk(b.Nomor, t.Kolom, t.Pesan, jenisTanggal, v)
			}
		}
		// 3. Uang lain - hanya bentuknya, bila terisi. ⛔ Korpus menyatakan
		// aturan titik enam kali dan menegakkannya sekali - lihat
		// `ErrUangBerkoma`.
		// Type QR: kolom yang DIHITUNG Calculate CSV tidak diperiksa - nilainya
		// di CSV diabaikan dan ditimpa (keputusan work owner 05-10-2026).
		for _, k := range KolomUangTersimpan() {
			if hitungQR && kolomHasilQR[k] {
				continue
			}
			if v := ambil(k); !sudah[k] && v != "" {
				periksaBentuk(b.Nomor, k, k+" bukan angka yang sah", jenisUang, v)
			}
		}
	}
	return hasil, nil
}

// polaBulat - bilangan bulat tak bertanda (umur).
var polaBulat = regexp.MustCompile(`^\d+$`)

// KolomWajibUnggah mengembalikan nama kolom yang HARUS ada di baris judul
// untuk satu Type.
//
// ⛔ Dipakai memeriksa JUDUL, bukan isi. Berkas yang kehilangan satu kolom
// sama sekali akan menghasilkan satu penolakan "kolom kosong" untuk SETIAP
// barisnya - seribu kalimat yang mengatakan satu hal. Memeriksanya sekali di
// judul memberi orang satu kalimat yang benar: kolom ini tidak ada di
// berkasmu.
func KolomWajibUnggah(tipe string) ([]string, error) {
	wajib, err := KolomWajibPerTipe(tipe)
	if err != nil {
		return nil, err
	}
	nama := make([]string, 0, len(wajib)+1)
	for _, w := range wajib {
		nama = append(nama, w.Kolom)
	}
	// Type QR: PERIOD_MM wajib ADA DI JUDUL (penentu CONTRACT rate/risk);
	// isinya diperiksa di tahap hitung, bukan sebagai "HARUS ADA" (keputusan
	// work owner 05-10-2026).
	if HitungPesertaPerTipe(tipe) {
		nama = append(nama, KolomPeriodeQR)
	}
	sort.Strings(nama)
	return nama, nil
}

// KolomPeriodeQR - kolom judul tambahan yang wajib untuk Type QR.
const KolomPeriodeQR = "PERIOD_MM"

// KolomUangUnggahTambahan - kolom uang CSV yang IKUT DISIMPAN di luar ke-32
// kolom `ValidasiUploadPL_act` (keputusan work owner 03-10-2026: summary
// harus memuat RI Admin Fee dan Deduction dari berkas). Kolomnya sudah ada di
// migrasi 052 - nol migrasi.
//
// ⚠️ TERPISAH dari `KolomUangUnggah`, dengan sengaja: daftar itu tiruan
// langkah 2 Pega (diisi 0 bila kosong, dijaga uji paritas). Kolom di sini
// TIDAK diisi 0 - kosong disimpan NULL, dan rekap membacanya sebagai nol.
//
// RATE ikut sejak 05-10-2026: Calculate CSV Type QR menghitungnya dan
// `T_PREMIUM_LIST_DETAIL.RATE` (migrasi 052) harus benar-benar tertulis.
var KolomUangUnggahTambahan = []string{"DEDUCTION", "RI_ADMIN_FEE", "SUM_AT_RISK_GROSS", "FACTOR", "RATE"}

// KolomUangTersimpan - seluruh kolom uang unggahan yang divalidasi bentuknya,
// dinormalkan desimal komanya, dan disimpan: `KolomUangUnggah` lalu
// `KolomUangUnggahTambahan`.
func KolomUangTersimpan() []string {
	k := make([]string, 0, len(KolomUangUnggah)+len(KolomUangUnggahTambahan))
	k = append(k, KolomUangUnggah...)
	return append(k, KolomUangUnggahTambahan...)
}
