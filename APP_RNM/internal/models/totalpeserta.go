package models

// Enam total uang SATU PESERTA — layar Detail.
//
// Untuk apa berkas ini: `Section/ClaimLifeDetailGCNM.xml` menampilkan enam
// medan `.Total*` pada halaman PESERTA (kelas `Int-LIFE_PREMIUM_DETAIL` b84,
// dan setiap ikatannya berawalan TITIK). Berkas ini menghitungnya.
//
// Sumber, dibaca sebagai pohon langkah utuh:
//
//	`Activity/SavePesertaClaim.xml`
//	  b4002 langkah 8   ULANG(EMBEDDED b4843) atas
//	                    `pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail`
//	                    = PESERTA; tanpa prasyarat (`...PreCondParamsWhen/` b4835)
//	       b4027..b4159 reset ENAM `local.Total*` ke literal `0`
//	       b4179        `Local.IndexPremium = .pxListSubscript`
//	    b4221 langkah 8.1 ULANG(EMBEDDED b4570) atas `.AdjustmentList` b4226;
//	                    tanpa prasyarat (`...PreCondParamsWhen/` b4562)
//	       b4246 `local.TotalCedingRetention = .CEDING_RETENTION + local.TotalCedingRetention`
//	       b4292 `.SHARE_NUSANTARA_RE`  b4312 `.SUM_INSURED`
//	       b4332 `.SUM_REASURED`        b4352 `.SHARE_RETRO`
//	       b4372 `.CLAIM_AMOUNT`
//	    b4592 langkah 8.2 TANPA ulang (`...HasRepeat/` b4801), tanpa prasyarat:
//	       b4616..b4743 keenam `local.Total*` ditulis ke halaman PESERTA
//
// Rumus yang sama persis ada pada `Activity/SaveOutStandingLife_Act.xml`
// langkah 23 b10638 / 23.1 b10841 (EMBEDDED b11046) / 23.2 b11067.
//
// ⛔ RALAT 27-09-2026, DAN INI RALAT ATAS KESIMPULAN KAMI SENDIRI. Ronde
// sebelumnya menyatakan keenam total "tidak dapat ditiru" sebab
// `CheckTotalAdjustmentClaim` tidak ada di korpus, lalu melaporkannya OQ-H.
// Rujukan menggantung itu NYATA, tetapi kesimpulannya keliru: yang hilang
// hanya pemanggil REFRESH di layar. NILAInya dihitung oleh kedua activity di
// atas, dan keduanya ADA. Sebab salahnya: berhenti pada rule yang namanya
// tertulis di section, tanpa menanyakan siapa lagi yang menulis medan itu.
//
// Dibaca sesudah: klaimlife.go (BarisAdjustment), money.go.

import (
	"errors"
	"fmt"
)

// ErrTotalMataUangBeragam menandai peserta yang baris adjustment-nya
// bermata uang lebih dari satu.
//
// ⛔ Ini GALAT, bukan peringatan, dan bukan pula "ambil yang pertama". Enam
// medan di layar ini angka uang (ADR-U-0003); menjumlahkan dua mata uang
// menghasilkan angka yang tampak benar dan tidak berarti apa-apa. Lewat jalur
// normal ia tidak dapat terjadi - `CURRENCY_ID` termasuk kolom WARISAN yang
// `SetIndexAdjustmentList` langkah 3 salin dari baris pertama - sehingga
// munculnya berarti datanya memang rusak, dan itu harus berbunyi.
var ErrTotalMataUangBeragam = errors.New("models: baris adjustment satu peserta bermata uang beragam")

// TotalPeserta adalah keenam total uang satu peserta.
//
// ⛔ TIDAK DISIMPAN. Tidak ada kolom `TOTAL_*` di `T_CLAIMLF_ADJUSTMENT`
// maupun `T_CLAIMLF_PREMIUMLIST_DETAIL`, dan tidak ada yang ditambahkan:
// ia turunan penuh dari baris adjustment, dan turunan yang disimpan adalah
// turunan yang suatu hari berbeda dari sumbernya.
//
// ⚠️ PENYIMPANGAN SADAR yang DICATAT: Pega menghitungnya saat `Submit`
// (`SavePesertaClaim`) dan `Save to RNM` (`SaveOutStandingLife_Act`) lalu
// MENULIS hasilnya ke halaman peserta. Layar Pega karena itu dapat basi -
// sesudah putaran atau akseptasi, angkanya tetap angka lama sampai tombol
// simpan ditekan lagi. Di sini ia dihitung saat DIBACA, sehingga tidak pernah
// basi. Arah selisihnya disengaja: angka yang selalu mutakhir tidak pernah
// lebih berbahaya daripada angka yang diam-diam tertinggal.
type TotalPeserta struct {
	// Nama JSON DITULIS, tidak diserahkan pada nama medan Go. Tanpa tag,
	// Go mengirim `CedingRetention` berhuruf besar dan React membaca
	// `undefined` tanpa satu pun galat - medan uang yang diam-diam kosong.
	CedingRetention  Money `json:"cedingRetention"`
	ShareNusantaraRe Money `json:"shareNusantaraRe"`
	SumInsured       Money `json:"sumInsured"`
	SumReasured      Money `json:"sumReasured"`
	ShareRetro       Money `json:"shareRetro"`
	JumlahKlaim      Money `json:"jumlahKlaim"`
}

// kolomTotal memasangkan tiap total dengan pengambil nilainya dari satu baris.
//
// ⚠️ Satu daftar, bukan enam blok berurutan: enam blok yang mirip adalah enam
// kesempatan menyalin-tempel lalu lupa mengganti satu nama, dan nama yang
// tertukar di sini menukar dua ANGKA UANG tanpa satu pun tipe yang berubah.
var kolomTotal = []struct {
	// nama adalah nama kolom Pega, dipakai di pesan galat supaya orang dapat
	// mencari kolom yang sama di kedua sistem.
	nama   string
	ambil  func(BarisAdjustment) Money
	simpan func(*TotalPeserta, Money)
}{
	{"CEDING_RETENTION", func(b BarisAdjustment) Money { return b.CedingRetention },
		func(t *TotalPeserta, m Money) { t.CedingRetention = m }},
	{"SHARE_NUSANTARA_RE", func(b BarisAdjustment) Money { return b.ShareNusantaraRe },
		func(t *TotalPeserta, m Money) { t.ShareNusantaraRe = m }},
	{"SUM_INSURED", func(b BarisAdjustment) Money { return b.SumInsured },
		func(t *TotalPeserta, m Money) { t.SumInsured = m }},
	{"SUM_REASURED", func(b BarisAdjustment) Money { return b.SumReasured },
		func(t *TotalPeserta, m Money) { t.SumReasured = m }},
	{"SHARE_RETRO", func(b BarisAdjustment) Money { return b.ShareRetro },
		func(t *TotalPeserta, m Money) { t.ShareRetro = m }},
	{"CLAIM_AMOUNT", func(b BarisAdjustment) Money { return b.JumlahKlaim },
		func(t *TotalPeserta, m Money) { t.JumlahKlaim = m }},
}

// JumlahKolomTotal adalah cacah total yang layar harus tampilkan.
//
// ⚠️ Ia dikunci uji di kedua sisi. Total yang DIHILANGKAN dari layar sama
// merusaknya dengan total yang dikarang, dan yang pertama tidak berbunyi
// sendiri - itulah sebab `Total Ceding Retention` sempat hilang enam hari:
// ia satu-satunya yang tidak punya aksi refresh, sehingga tidak ikut terhitung
// ketika rujukan `CheckTotalAdjustmentClaim` yang dicacah.
const JumlahKolomTotal = 6

// HitungTotalPeserta menjumlahkan SELURUH baris adjustment satu peserta.
//
// ⛔ SELURUHNYA - tanpa memandang `STS_REJECT`, dan tanpa memandang
// `IsCheck`. Bukan kelonggaran: langkah 8.1 b4221 dan langkah 23.1 b10841
// sama-sama BERPRASYARAT KOSONG. Menyaring baris yang ditolak akan
// menghasilkan angka yang lebih "masuk akal" dan BUKAN angka sistem lama.
//
// ⛔ Baris berkolom KOSONG tidak menggagalkan penjumlahan dan tidak dianggap
// bernilai nol: ia tidak menyumbang. Penampungnya mulai dari literal `0`
// (b4027..b4159), sehingga peserta TANPA baris adjustment berjumlah `0` -
// bukan kosong. Itu yang Pega tulis, dan ADR-U-0027 tidak membantahnya:
// ADR itu mengatur KOLOM yang kosong, bukan unsur identitas penjumlahan.
//
// mataUang adalah mata uang peserta. Bila kosong, mata uang baris pertama
// yang mengisinya; bila baris pun tidak menyebut, hasilnya bermata uang
// kosong - dan itu jujur, sebab memang tidak ada yang menyebutnya.
func HitungTotalPeserta(baris []BarisAdjustment, mataUang string) (TotalPeserta, error) {
	kurs := mataUang
	for _, b := range baris {
		if b.CurrencyID == "" {
			continue
		}
		if kurs == "" {
			kurs = b.CurrencyID
			continue
		}
		if b.CurrencyID != kurs {
			return TotalPeserta{}, fmt.Errorf("%w: %q dan %q",
				ErrTotalMataUangBeragam, kurs, b.CurrencyID)
		}
	}

	var total TotalPeserta
	for _, kolom := range kolomTotal {
		jumlah, err := NewMoney("0", kurs)
		if err != nil {
			return TotalPeserta{}, fmt.Errorf("total %s: %w", kolom.nama, err)
		}
		for _, b := range baris {
			nilai := kolom.ambil(b)
			if nilai.Kosong() {
				continue
			}
			// ⛔ MATA UANG TIAP KOLOM DIPERIKSA SENDIRI, tidak dilabeli
			// ulang. Ronde pertama menulis `nilai.Currency = kurs` dengan
			// komentar "yang berbeda sudah ditolak di atas" - dan komentar
			// itu KELIRU: pemeriksaan di atas hanya membaca `CurrencyID`
			// BARIS, tidak pernah membaca Currency tiap nilai uang. Satu
			// kolom bermata uang lain karena itu dijumlahkan diam-diam di
			// bawah mata uang yang salah, dan `Money.Add` yang seharusnya
			// menolaknya justru dilucuti lebih dulu (ADR-U-0003).
			//
			// Yang KOSONG mengikut penampung: nilai uang tanpa label mata
			// uang bukan mata uang lain, ia hanya belum berlabel.
			if nilai.Currency == "" {
				nilai.Currency = kurs
			} else if nilai.Currency != kurs {
				return TotalPeserta{}, fmt.Errorf("%w: kolom %s bermata uang %q, peserta %q",
					ErrTotalMataUangBeragam, kolom.nama, nilai.Currency, kurs)
			}
			if jumlah, err = jumlah.Add(nilai); err != nil {
				return TotalPeserta{}, fmt.Errorf("total %s: %w", kolom.nama, err)
			}
		}
		kolom.simpan(&total, jumlah)
	}
	return total, nil
}
