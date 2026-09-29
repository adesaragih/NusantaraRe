package services

// Perhitungan spreading satu baris adjustment - tiket 03.
//
// Untuk apa berkas ini: meniru `Claim Life/Activity/SpreadingClaimLife_Act.xml`
// langkah demi langkah, dengan desimal eksak dan nol float. Tiap fungsi
// menyebut nomor langkah XML yang ditirunya di komentar kepalanya.
//
// Dibaca sesudah: adjustment.go.
//
// ⛔ Seluruh isi berkas ini MURNI: tidak menyentuh basis data, tidak membaca
// jam, dan hasilnya hanya bergantung masukannya. Pembacaan treaty dan rate
// milik repository; yang masuk ke sini sudah berupa nilai.
//
// Istilah:
//   - treaty-year : satu tahun perjanjian retrosesi, dengan kapasitasnya.
//   - kaskade     : sisa klaim yang mengalir ke treaty-year berikutnya.
//   - per mil     : per seribu. Rate di RATE_LIFE dinyatakan per mil.

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/uang"
	"nusantarare/inti/utils"
	"nusantarare/modul/claimlife/models"
)

// Pembulatan yang XML tuliskan, masing-masing dengan langkah asalnya.
// Angkanya tidak pernah ditulis telanjang di dalam rumus.
const (
	// .RATE = @divide(local.Rate,1,4) - langkah 8.2.1.9.2 baris 4396.
	desimalRateTersimpan = 4
	// local.Rate = @divide(@toDecimal(local.Rate),1000,10) - baris 4448.
	desimalRatePakai = 10
	// @divide(@toDecimal(.OVR_COMM),100,5) dan .COMMISION - baris 4468, 4796.
	desimalPersenKomisi = 5
	// @divide(.PERCENTSHARE,100,4) - baris 4488.
	desimalPersenShare = 4
	// ⛔ desimalUang TIDAK berasal dari XML - ia ADR-U-0016 Akibat 3: lapisan
	// services membulatkan pada desimal kedelapan SAAT MEMUAT. Tanpa itu,
	// Oracle yang membulatkan diam-diam saat menulis ke NUMBER(38,8), dan
	// angka yang dihitung berbeda dari angka yang tersimpan tanpa satu pun
	// galat. Kolom persen dan rate TIDAK ikut (Akibat 2) - keduanya mengikuti
	// pembulatan XML di atas.
	desimalUang = 8
)

// pembagiPerMil: RATE_LIFE menyatakan rate per seribu, dan XML membaginya
// seribu SESUDAH menyimpan nilai mentahnya.
const pembagiPerMil = 1000

// Galat yang dapat dikenali pemanggil dengan errors.Is.
var (
	// ErrRateTidakDitemukan: nol baris RATE_LIFE yang cocok. Pega memakai
	// rate baris sebelumnya di keadaan ini - lihat komentar PilihRate.
	ErrRateTidakDitemukan = errors.New("services: tidak ada baris rate yang cocok")
	// ErrRateBerganda: beberapa baris cocok dengan nilai BERBEDA.
	ErrRateBerganda = errors.New("services: baris rate yang cocok lebih dari satu dan berbeda nilai")
	// ErrKapasitasBelumDiketahui: kolom IDR/USD treaty-year belum ada isinya.
	ErrKapasitasBelumDiketahui = errors.New("services: kapasitas treaty-year belum diketahui")
	// ErrMataUangTanpaCabang: XML hanya punya cabang IDR dan USD.
	ErrMataUangTanpaCabang = errors.New("services: mata uang tanpa cabang kaskade di rule sumber")
	// ErrTreatyTidakTerurut: daftar treaty-year tiba tidak menaik kapasitas.
	ErrTreatyTidakTerurut = errors.New("services: treaty-year tidak terurut menaik menurut kapasitas")
)

// TahunTreaty adalah satu baris `treatyyear_life` beserta retrosesinya.
//
// ⚠️ Urutannya BERMAKNA: `GetJsonProductLife` menutup SQL-nya dengan
// `order by TO_NUMBER(IDR) asc`, jadi kapasitas terkecil dipakai lebih dulu.
// Pemanggil menyerahkan daftar yang sudah terurut, dan HitungSpreading
// MEMERIKSANYA - ia tidak mengurutkan ulang. Mengurutkan diam-diam akan
// menyembunyikan pembaca yang lupa `ORDER BY`; memeriksa membuatnya berbunyi.
type TahunTreaty struct {
	ID             string
	TreatyYearLife string
	// IDR dan USD adalah KAPASITAS treaty-year ini, dan boleh nil.
	//
	// ⛔ nil BUKAN nol. POOLDATA.TREATYYEAR_LIFE tidak punya kedua kolom ini
	// `[data DBA]`; memperlakukan yang kosong sebagai nol akan membuat seluruh
	// share nol tanpa seorang pun tahu sebabnya.
	IDR *uang.Money
	USD *uang.Money
	// Retro adalah baris RETROCESSIONLIFE milik treaty-year ini, urut ID -
	// `GetRetroLife_SQL`: `order by id asc`.
	Retro []BarisRetro
}

// BarisRetro adalah satu reinsurer pada satu treaty-year.
//
// ⚠️ PercentShare, Commision, dan OvrComm adalah PERSEN (0-100), bukan pecahan
// dan bukan uang. XML membagi ketiganya seratus sebelum memakainya.
type BarisRetro struct {
	ID             string
	ReinsurerName  string
	TreatyTypeID   string
	TreatyTypeName string
	PercentShare   uang.Ratio
	Commision      uang.Ratio
	OvrComm        uang.Ratio
}

// BarisRate adalah satu baris hasil `GetRateRetro` atas POOLDATA.RATE_LIFE.
//
// Seluruh medannya tiba sebagai TEKS dari view ber-CLOB, dan Rate sudah
// melewati penggantian koma-ke-titik di repository - meniru langkah 8.2.1.9.1
// `@toDecimal(@replaceAll(.CARI4,",","."))`.
type BarisRate struct {
	ID       string
	Umur     string // CARI1 = AGE
	Kontrak  string // CARI2 = CONTRACT
	JenisKel string // CARI3 = GENDER
	Rate     uang.Ratio
}

// MasukanSpreading adalah seluruh nilai yang perhitungan ini perlukan.
//
// ⚠️ Asalnya dua lapisan yang BERBEDA, dan XML-lah yang menentukan mana dari
// mana (`pyStepsObjectName` tiap langkah): EmPercent, Umur, dan TahunPolis dari
// PESERTA (langkah 8.1); Currency dan ClaimGross dari BARIS ADJUSTMENT
// (langkah 8.2). Menggabungkannya di satu struct membuat pemanggil menyatakan
// keduanya dengan sadar.
type MasukanSpreading struct {
	// Currency dan ClaimGross - baris adjustment. Kolom kita untuk CLAIM_GROSS
	// bernama CLAIM_AMOUNT: satu nilai, dua nama.
	Currency   string
	ClaimGross uang.Money
	// EmPercent dipakai sebagai PECAHAN LANGSUNG. Langkah 8.1 menyalin
	// .EM_PERCENT tanpa membaginya seratus - berbeda dari OVR_COMM dan
	// COMMISION, yang dibagi. Nama "percent" di korpus ini tidak dapat dipakai
	// menebak sifatnya.
	EmPercent uang.Ratio
	Umur      string
	JenisKel  string
	// TahunPolis adalah local.Year, dan ia yang MEMILIH cabang rumus NET.
	TahunPolis int
}

// TahunPolis menghitung local.Year - SpreadingClaimLife_Act langkah 8.1.
//
// XML-nya `@toDecimal(@substring(.GROSS_VALUATION_BEGIN_DATE,6,10)) -
// @toDecimal(@substring(.BEGIN_DATE,6,10)) + 1`: ia memotong BAGIAN TAHUN dari
// teks tanggal, lalu mengurangkan. Di sini tanggalnya sudah bertipe waktu,
// sehingga pemotongan teks - beserta ketergantungannya pada bentuk tanggal -
// hilang. Hasilnya sama; jebakannya tidak ikut.
func TahunPolis(mulai, valuasiKotor time.Time) int {
	return valuasiKotor.Year() - mulai.Year() + 1
}

// PilihRate memilih satu rate dari hasil GetRateRetro.
//
// SpreadingClaimLife_Act langkah 8.2.1.9.1 dan 8.2.1.9.2, dibaca beserta AKSI
// precondition-nya (2 = lanjut, 3 = lewati langkah):
//
//	GENDER memuat "U"  -> cocokkan UMUR saja
//	selain itu         -> cocokkan UMUR dan JENIS KELAMIN
//
// CONTRACT sengaja TIDAK ikut: ia hanya muncul di `pyExpression` yang tidak
// aktif (varian lama `PERIOD_YY==.CARI2`), bukan di precondition yang berjalan.
//
// ⛔ Dua penyimpangan sadar dari rule sumber, keduanya dilaporkan ke work owner:
//
//	(1) Pega tidak pernah mengosongkan local.Rate. Bila nol baris cocok, ia
//	    memakai rate baris retro sebelumnya - yang SUDAH dibagi seribu - lalu
//	    membaginya seribu sekali lagi. Di sini keadaan itu menjadi GALAT.
//	(2) Putaran Pega tidak berhenti pada yang pertama cocok, jadi yang menang
//	    adalah yang TERAKHIR, sedangkan GetRateRetro tidak memakai ORDER BY.
//	    Dengan 2.693 kombinasi ganda di RATE_LIFE `[data DBA]`, hasilnya tidak
//	    tertentu. Di sini ambiguitas DILAPORKAN beserta ID barisnya.
func PilihRate(baris []BarisRate, umur, jenisKel string) (uang.Ratio, error) {
	var cocok []BarisRate
	for _, b := range baris {
		if cocokRate(b, umur, jenisKel) {
			cocok = append(cocok, b)
		}
	}
	if len(cocok) == 0 {
		return uang.Ratio{}, fmt.Errorf("%w: umur %q jenis kelamin %q dari %d baris rate",
			ErrRateTidakDitemukan, umur, jenisKel, len(baris))
	}
	pertama := cocok[0]
	for _, b := range cocok[1:] {
		if samaNilai(pertama.Rate, b.Rate) {
			continue
		}
		var id []string
		for _, c := range cocok {
			id = append(id, c.ID)
		}
		return uang.Ratio{}, fmt.Errorf(
			"%w: umur %q jenis kelamin %q cocok pada baris %s dengan nilai berbeda",
			ErrRateBerganda, umur, jenisKel, strings.Join(id, ", "))
	}
	return pertama.Rate, nil
}

// cocokRate menerapkan kedua precondition langkah 8.2.1.9.1 dan 8.2.1.9.2.
func cocokRate(b BarisRate, umur, jenisKel string) bool {
	if !samaNilaiTeks(b.Umur, umur) {
		return false
	}
	if strings.Contains(strings.ToUpper(strings.TrimSpace(b.JenisKel)), "U") {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(b.JenisKel), strings.TrimSpace(jenisKel))
}

// samaNilaiTeks membandingkan dua teks angka tanpa terpeleset pada nol di depan
// maupun spasi tepi. "30", " 30 ", dan "030" adalah umur yang sama; yang bukan
// angka dibandingkan apa adanya.
func samaNilaiTeks(a, b string) bool {
	a, b = strings.TrimSpace(a), strings.TrimSpace(b)
	if a == b {
		return true
	}
	da, errA := utils.ParseDecimal(a)
	db, errB := utils.ParseDecimal(b)
	if errA != nil || errB != nil {
		return false
	}
	return da.Cmp(db) == 0
}

// samaNilai membandingkan dua rasio menurut NILAINYA, sehingga 2.5 dan 2.50
// tidak dianggap dua rate yang berbeda.
func samaNilai(a, b uang.Ratio) bool {
	if a.Value == nil || b.Value == nil {
		return a.Value == nil && b.Value == nil
	}
	return a.Value.Cmp(b.Value) == 0
}

// HitungSpreading membentuk seluruh baris spreading satu baris adjustment.
//
// Meniru SpreadingClaimLife_Act langkah 8.2.1 sampai 8.2.1.9.3. Urutannya:
// kaskade kapasitas per treaty-year, lalu satu baris retro per reinsurer.
//
// ⭐ Tanpa treaty-year hasilnya NOL BARIS - bukan galat. Treaty-year tanpa
// reinsurer tetap menghasilkan barisnya sendiri dengan nol anak: baris
// spreading lahir dari daftar treaty-year (`.SpreadingList = OutwardList
// .pxResults`, langkah 7.1), bukan dari daftar reinsurer.
func HitungSpreading(m MasukanSpreading, tahun []TahunTreaty,
	rate []BarisRate) ([]models.Spreading, error) {

	if len(tahun) == 0 {
		return nil, nil
	}
	// local.ClaimNet = .CLAIM_GROSS - langkah 8.2 baris 2465.
	sisa := nilaiAtauNol(m.ClaimGross.Amount)

	// Rate dicari SEKALI, dan malas. Hasilnya hanya bergantung pada masukan,
	// bukan pada baris retro, sehingga mencarinya per baris berarti pekerjaan
	// dan pesan galat yang sama berulang sebanyak jumlah reinsurer. Malas dan
	// bukan di muka: di XML pencarian rate ada DI DALAM putaran baris retro,
	// jadi treaty-year tanpa retrosesi memang tidak pernah memerlukannya.
	var rateTerpilih uang.Ratio
	var rateSudah bool
	// Kapasitas treaty-year sebelumnya, untuk memeriksa urutan menaik.
	var sebelumnya *apd.Decimal

	out := make([]models.Spreading, 0, len(tahun))
	for _, t := range tahun {
		kapasitas, err := kapasitasUntuk(t, m.Currency)
		if err != nil {
			return nil, err
		}
		// ⛔ Urutan menaik diperiksa, bukan diandaikan. Daftar yang tiba tanpa
		// ORDER BY membagi klaim ke treaty-year yang salah, dan angkanya tetap
		// tampak masuk akal - cacat yang tidak akan pernah terlihat sendiri.
		if sebelumnya != nil && kapasitas.Cmp(sebelumnya) < 0 {
			return nil, fmt.Errorf("%w: treaty-year %s berkapasitas lebih kecil "+
				"daripada pendahulunya; sumbernya wajib `order by TO_NUMBER(IDR) asc`",
				ErrTreatyTidakTerurut, t.ID)
		}
		sebelumnya = kapasitas

		// Langkah 8.2.1.4-7: bila sisa <= kapasitas, seluruh sisa diambil dan
		// sisanya menjadi nol; bila lebih besar, kapasitasnya yang diambil dan
		// selisihnya mengalir ke treaty-year berikutnya.
		share := new(apd.Decimal)
		if sisa.Cmp(kapasitas) <= 0 {
			share.Set(sisa)
			sisa = apd.New(0, 0)
		} else {
			share.Set(kapasitas)
			baru := new(apd.Decimal)
			if _, err := utils.DecimalContext().Sub(baru, sisa, kapasitas); err != nil {
				return nil, fmt.Errorf("services: mengurangi kapasitas treaty-year %s: %w",
					t.ID, err)
			}
			sisa = baru
		}
		if share, err = bulat(utils.DecimalContext(), share, desimalUang); err != nil {
			return nil, fmt.Errorf("services: membulatkan share treaty-year %s: %w", t.ID, err)
		}

		spr := models.Spreading{
			TreatyYearLife:  t.TreatyYearLife,
			RetrocadedShare: uang.Money{Amount: share, Currency: m.Currency},
			Currency:        m.Currency,
		}
		if t.IDR != nil {
			spr.IDR = *t.IDR
		}
		if t.USD != nil {
			spr.USD = *t.USD
		}
		// Jenis treaty hidup di baris retrosesi, dan seluruh baris satu
		// treaty-year membawa jenis yang sama di DEV. Yang pertama dipakai
		// sebagai jenis barisnya; bila nol retro, ia tetap kosong.
		if len(t.Retro) > 0 {
			spr.TreatyTypeID = t.Retro[0].TreatyTypeID
			spr.TreatyTypeName = t.Retro[0].TreatyTypeName
		}

		for _, r := range t.Retro {
			if !rateSudah {
				if rateTerpilih, err = PilihRate(rate, m.Umur, m.JenisKel); err != nil {
					return nil, fmt.Errorf("treaty-year %s: %w", t.ID, err)
				}
				rateSudah = true
			}
			anak, err := hitungRetro(m, share, r, rateTerpilih)
			if err != nil {
				return nil, fmt.Errorf("treaty-year %s reinsurer %s: %w",
					t.ID, r.ReinsurerName, err)
			}
			spr.Retro = append(spr.Retro, anak)
		}
		out = append(out, spr)
	}
	return out, nil
}

// kapasitasUntuk memilih kolom kapasitas menurut mata uang - langkah 8.2.1.4-7.
func kapasitasUntuk(t TahunTreaty, mataUang string) (*apd.Decimal, error) {
	var kapasitas *uang.Money
	switch strings.ToUpper(strings.TrimSpace(mataUang)) {
	case "IDR":
		kapasitas = t.IDR
	case "USD":
		kapasitas = t.USD
	default:
		// ⛔ Tidak ditebak. Rule sumber hanya punya dua cabang, dan mata uang
		// ketiga berarti pertanyaan yang belum dijawab siapa pun.
		return nil, fmt.Errorf("%w: %q pada treaty-year %s",
			ErrMataUangTanpaCabang, mataUang, t.ID)
	}
	if kapasitas == nil || kapasitas.Kosong() {
		return nil, fmt.Errorf("%w: treaty-year %s mata uang %s",
			ErrKapasitasBelumDiketahui, t.ID, mataUang)
	}
	return kapasitas.Amount, nil
}

// hitungRetro menghitung satu baris retro - langkah 8.2.1.9.2 dan 8.2.1.9.3.
//
// Kedua cabang berbagi rumus sampai GROSS; yang membedakan hanya NET, dan yang
// memilihnya adalah TAHUN POLIS:
//
//	tahun ke-1 : NET = GROSS - Discount - Comm,  Comm dihitung atas sisa
//	             sesudah Discount
//	selain itu : NET = GROSS - Comm
func hitungRetro(m MasukanSpreading, share *apd.Decimal, r BarisRetro,
	rateMentah uang.Ratio) (models.SpreadingRetro, error) {

	ctx := utils.DecimalContext()
	nol := models.SpreadingRetro{}
	var err error

	// .RATE = @divide(local.Rate,1,4) - yang TERSIMPAN adalah rate mentah.
	rateTersimpan, err := bulat(ctx, nilaiAtauNol(rateMentah.Value), desimalRateTersimpan)
	if err != nil {
		return nol, fmt.Errorf("membulatkan rate tersimpan: %w", err)
	}
	// local.Rate = @divide(local.Rate,1000,10) - yang DIPAKAI sudah dibagi.
	ratePakai := new(apd.Decimal)
	if _, err := ctx.Quo(ratePakai, nilaiAtauNol(rateMentah.Value),
		apd.New(pembagiPerMil, 0)); err != nil {
		return nol, fmt.Errorf("membagi rate per mil: %w", err)
	}
	if ratePakai, err = bulat(ctx, ratePakai, desimalRatePakai); err != nil {
		return nol, fmt.Errorf("membulatkan rate pakai: %w", err)
	}

	// .Amount = local.AmountRetroShare * @divide(.PERCENTSHARE,100,4)
	porsi, err := persenJadiPecahan(ctx, r.PercentShare, desimalPersenShare)
	if err != nil {
		return nol, fmt.Errorf("PERCENT_SHARE: %w", err)
	}
	jumlah := new(apd.Decimal)
	if _, err := ctx.Mul(jumlah, share, porsi); err != nil {
		return nol, fmt.Errorf("mengalikan share dengan porsi: %w", err)
	}
	if jumlah, err = bulat(ctx, jumlah, desimalUang); err != nil {
		return nol, fmt.Errorf("membulatkan AMOUNT: %w", err)
	}

	// .PREMIUM_SPREADED_GROSS = local.Rate * (1+local.EMPercent) * .Amount
	faktorEM := new(apd.Decimal)
	if _, err := ctx.Add(faktorEM, apd.New(1, 0), nilaiAtauNol(m.EmPercent.Value)); err != nil {
		return nol, fmt.Errorf("menjumlah EM_PERCENT: %w", err)
	}
	bruto := new(apd.Decimal)
	if _, err := ctx.Mul(bruto, ratePakai, faktorEM); err != nil {
		return nol, fmt.Errorf("mengalikan rate dengan EM: %w", err)
	}
	if _, err := ctx.Mul(bruto, bruto, jumlah); err != nil {
		return nol, fmt.Errorf("mengalikan dengan jumlah: %w", err)
	}
	// ⭐ GROSS dibulatkan SEBELUM NET diturunkan darinya. Alasannya bukan
	// ketelitian melainkan konsistensi: NET = GROSS - Comm harus berlaku pada
	// angka yang benar-benar masuk basis data, bukan hanya pada nilai antara
	// yang tidak pernah tersimpan.
	if bruto, err = bulat(ctx, bruto, desimalUang); err != nil {
		return nol, fmt.Errorf("membulatkan PREMIUM_SPREADED_GROSS: %w", err)
	}

	komisi, err := persenJadiPecahan(ctx, r.OvrComm, desimalPersenKomisi)
	if err != nil {
		return nol, fmt.Errorf("OVR_COMM: %w", err)
	}

	neto := new(apd.Decimal)
	if m.TahunPolis == 1 {
		// Langkah 8.2.1.9.3 "Tahun pertama".
		diskon, err := persenJadiPecahan(ctx, r.Commision, desimalPersenKomisi)
		if err != nil {
			return nol, fmt.Errorf("COMMISION: %w", err)
		}
		nilaiDiskon := new(apd.Decimal)
		if _, err := ctx.Mul(nilaiDiskon, bruto, diskon); err != nil {
			return nol, fmt.Errorf("mengalikan diskon: %w", err)
		}
		sesudahDiskon := new(apd.Decimal)
		if _, err := ctx.Sub(sesudahDiskon, bruto, nilaiDiskon); err != nil {
			return nol, fmt.Errorf("mengurangi diskon: %w", err)
		}
		nilaiKomisi := new(apd.Decimal)
		if _, err := ctx.Mul(nilaiKomisi, sesudahDiskon, komisi); err != nil {
			return nol, fmt.Errorf("mengalikan komisi: %w", err)
		}
		if _, err := ctx.Sub(neto, sesudahDiskon, nilaiKomisi); err != nil {
			return nol, fmt.Errorf("mengurangi komisi: %w", err)
		}
	} else {
		// Langkah 8.2.1.9.2 "Bukan tahun pertama".
		nilaiKomisi := new(apd.Decimal)
		if _, err := ctx.Mul(nilaiKomisi, bruto, komisi); err != nil {
			return nol, fmt.Errorf("mengalikan komisi: %w", err)
		}
		if _, err := ctx.Sub(neto, bruto, nilaiKomisi); err != nil {
			return nol, fmt.Errorf("mengurangi komisi: %w", err)
		}
	}
	if neto, err = bulat(ctx, neto, desimalUang); err != nil {
		return nol, fmt.Errorf("membulatkan PREMIUM_SPREADED_NET: %w", err)
	}

	return models.SpreadingRetro{
		ReinsurerName:        r.ReinsurerName,
		PercentShare:         r.PercentShare,
		Amount:               uang.Money{Amount: jumlah, Currency: m.Currency},
		Rate:                 uang.Ratio{Value: rateTersimpan},
		PremiumSpreadedGross: uang.Money{Amount: bruto, Currency: m.Currency},
		PremiumSpreadedNet:   uang.Money{Amount: neto, Currency: m.Currency},
		Commision:            r.Commision,
		OvrComm:              r.OvrComm,
		TreatyTypeID:         r.TreatyTypeID,
		TreatyTypeName:       r.TreatyTypeName,
	}, nil
}

// persenJadiPecahan membagi seratus lalu membulatkan - @divide(x,100,n).
func persenJadiPecahan(ctx *apd.Context, persen uang.Ratio, desimal int32) (*apd.Decimal, error) {
	hasil := new(apd.Decimal)
	if _, err := ctx.Quo(hasil, nilaiAtauNol(persen.Value), apd.New(100, 0)); err != nil {
		return nil, err
	}
	return bulat(ctx, hasil, desimal)
}

// bulat membulatkan ke sekian angka di belakang koma - padanan argumen ketiga
// @divide. Pembulatannya mengikuti konteks desimal aplikasi, satu tempat untuk
// seluruh modul.
func bulat(ctx *apd.Context, d *apd.Decimal, desimal int32) (*apd.Decimal, error) {
	hasil := new(apd.Decimal)
	if _, err := ctx.Quantize(hasil, d, -desimal); err != nil {
		return nil, err
	}
	hasil.Reduce(hasil)
	return hasil, nil
}
