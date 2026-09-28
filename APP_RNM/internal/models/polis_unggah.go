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
// ⚠️ CELAH TERCATAT - OQ-PL-12: karena langkah 2, kolom uang wajib yang
// KOSONG di sistem lama menjadi 0 dan lolos "HARUS ADA" langkah 9.12-9.17;
// di sini ia ditolak (ADR-U-0027: kosong bukan nol). Keputusannya milik work
// owner.
//
// ⛔ PESAN DISALIN VERBATIM, termasuk yang menjanjikan hal yang tidak
// diperiksa. Lihat `PesanNetPremium` dan OQ-069.
//
// ⛔ ⚠️ NORMALISASI KOMA PEGA TIDAK DITIRU. Langkah "Rubah decimal dari koma
// jadi titik" memakai `@replaceAll(.KOLOM, ",", ".")` atas tiap kolom uang -
// MENGGANTI SETIAP KOMA, tanpa membedakan pemisah desimal dari pemisah
// ribuan. `1,234,567.89` menjadi `1.234.567.89`, yang bukan angka sama
// sekali, dan hasilnya tersimpan sebagai uang. Kami MENOLAK komanya alih-alih
// menggantinya - lihat `UangCSV`.
//
// Dibaca sesudah: polis_nomor.go, polis_detail.go.

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/pkg/utils"
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
		return nil, ErrUangKosong
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

// kolomUangWajib adalah enam kolom uang yang HARUS ada isinya.
//
// `[terverifikasi]` langkah 9 mewajibkan keenamnya per baris, dan masing-
// masing punya pesannya sendiri. Kolom uang lain divalidasi BENTUKNYA bila
// terisi, tetapi tidak diwajibkan - persis korpus.
var kolomUangWajib = map[string]string{
	"SUM_INSURED":        PesanSumInsured,
	"CEDING_RETENTION":   PesanCedingRetention,
	"SUM_REASURED":       PesanSumReasured,
	"SHARE_NUSANTARA_RE": PesanShareNusantaraRe,
	"GROSS_PREMIUM":      PesanGrossPremium,
	"NET_PREMIUM":        PesanNetPremium,
}

// kolomTanggalWajib adalah kolom tanggal yang harus ada dan berbentuk.
//
// `[terverifikasi]` langkah 9: tiap kolom punya sepasang precondition -
// `@PropertyHasValue(.X)` lalu `@toDate(.X)!=0`.
var kolomTanggalWajib = []struct{ Kolom, Pesan string }{
	{"DOB", PesanDOB},
	{"BEGIN_DATE", PesanBeginDate},
	{"EXPIRED_DATE", PesanExpiredDate},
	{"START_DATE", PesanStartDate},
	{"EFFECTIVE_DATE", PesanEffectiveDate},
	{"STNC", PesanSTNC},
	{"WPC", PesanWPC},
	{"GROSS_VALUATION_BEGIN_DATE", PesanGrossValMulai},
	{"GROSS_VALUATION_EXPIRED_DATE", PesanGrossValSelesai},
	{"RETROCESSION_VALUATION_BEGIN_DATE", PesanRetroValMulai},
	{"RETROCESSION_VALUATION_EXPIRED_DATE", PesanRetroValSelesai},
}

// kolomTeksWajib adalah kolom teks yang harus ada isinya.
var kolomTeksWajib = []struct{ Kolom, Pesan string }{
	{"CERTIFICATE_NO", PesanCertificateNo},
	{"NAME_OF_INSURED", PesanNamaTertanggung},
	{"PLAN", PesanPlan},
	{"CURRENCY", PesanCurrency},
	{"POLICY_NO", PesanPolicyNo},
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

// ValidasiUnggah memeriksa seluruh baris dan mengumpulkan setiap penolakan.
//
// ⛔ SELURUH BARIS DIPERIKSA, bahkan sesudah satu ditolak - lihat
// `HasilUnggah.Ditolak`.
func ValidasiUnggah(baris []BarisUnggah) HasilUnggah {
	hasil := HasilUnggah{CacahBaris: len(baris), Ditolak: []Penolakan{}}
	tolak := func(b int, kolom, pesan, sebab string) {
		hasil.Ditolak = append(hasil.Ditolak, Penolakan{
			Baris: b, Kolom: kolom, Pesan: pesan, Sebab: sebab})
	}

	for _, b := range baris {
		ambil := func(k string) string { return strings.TrimSpace(b.Nilai[k]) }

		for _, w := range kolomTeksWajib {
			if ambil(w.Kolom) == "" {
				tolak(b.Nomor, w.Kolom, w.Pesan, "kolom kosong")
			}
		}

		for _, w := range kolomTanggalWajib {
			v := ambil(w.Kolom)
			if v == "" {
				tolak(b.Nomor, w.Kolom, w.Pesan, "kolom kosong")
				continue
			}
			if _, err := TanggalCSV(v); err != nil {
				tolak(b.Nomor, w.Kolom, w.Pesan,
					fmt.Sprintf("%q bukan dd/mm/yyyy", v))
			}
		}

		// ⛔ SETIAP kolom uang diperiksa BENTUKNYA bila terisi; keenam yang
		// wajib juga diperiksa KEBERADAANNYA. Korpus menyatakan aturan titik
		// enam kali dan menegakkannya sekali - lihat `ErrUangBerkoma`.
		for _, k := range KolomUangUnggah {
			v := ambil(k)
			pesan, wajib := kolomUangWajib[k]
			if v == "" {
				if wajib {
					tolak(b.Nomor, k, pesan, "kolom kosong")
				}
				continue
			}
			if !wajib {
				pesan = k + " bukan angka yang sah"
			}
			if _, err := UangCSV(v); err != nil {
				sebab := fmt.Sprintf("%q bukan desimal ber-pemisah titik", v)
				if errors.Is(err, ErrUangBerkoma) {
					sebab = fmt.Sprintf("%q memuat koma; pemisah desimal harus "+
						"titik dan pemisah ribuan tidak diterima", v)
				}
				tolak(b.Nomor, k, pesan, sebab)
			}
		}
	}
	return hasil
}

// KolomWajibUnggah mengembalikan nama kolom yang HARUS ada di baris judul.
//
// ⛔ Dipakai memeriksa JUDUL, bukan isi. Berkas yang kehilangan satu kolom
// sama sekali akan menghasilkan satu penolakan "kolom kosong" untuk SETIAP
// barisnya - seribu kalimat yang mengatakan satu hal. Memeriksanya sekali di
// judul memberi orang satu kalimat yang benar: kolom ini tidak ada di
// berkasmu.
func KolomWajibUnggah() []string {
	nama := make([]string, 0,
		len(kolomTeksWajib)+len(kolomTanggalWajib)+len(kolomUangWajib))
	for _, w := range kolomTeksWajib {
		nama = append(nama, w.Kolom)
	}
	for _, w := range kolomTanggalWajib {
		nama = append(nama, w.Kolom)
	}
	for k := range kolomUangWajib {
		nama = append(nama, k)
	}
	sort.Strings(nama)
	return nama
}
