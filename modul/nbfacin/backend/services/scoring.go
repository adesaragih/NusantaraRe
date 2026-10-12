package services

// Tab Scoring Risk kasus FIRE (tiket 54). Bukti `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\`:
//   - `Section\CallDualScoringRisk.xml` - "Write" bila PositionNote = 'ReasFacInMarketing' && FlagOnGoingPolicy = 0, selain
//     itu "Read";
//   - `Activity\ValueScoringRisk_act.xml` langkah 3-26 - poin tiap centang (poinScoring);
//   - `Activity\SetDataScoringRisk_act.xml` - isi otomatis kolom dari lokasinya: Status, okupasi (KDRISKEXPOSURE 01-05),
//     TSI IDR ≤ / > 30.000.000.000, TotalRNMShare (Σ PremiumSpreaded lokasi dalam IDR), zona gempa "Daerah Zona 1..5";
//   - `Activity\SetLossRatio_Act.xml` - band loss ratio ChechBox4..9;
//   - `Activity\SetScore_Act.xml` - ScorePerFactor = poin centang, Score = Σ, FinalScore;
//   - `Activity\ScoringResult.xml` - kelas, "Due to", "Approval From", rekomendasi;
//   - `Activity\SetDataOccupationScoringRisk_Act.xml` - Dominant Risk: "Can't Choose The Same Occupation".
// Poin, Score, dan persetujuan diuji ulang atas 5 kasus `DDL\CONTOH` ber-FinalScore (119/120 ScorePerFactor cocok; satu
// selisih = skor basi Pega yang dibetulkan di bawah).
//
// Keputusan work owner 06-10-2026 (docs/issues/54): hanya marketing (bolehUbah); cacat Pega DIBETULKAN - Earthquake kolom 2
// memakai kolom 1 / Tsunami menjumlah halaman induk (di sini setiap faktor dari kolomnya sendiri), skor tidak dihitung ulang
// saat muat (di sini GET dan setiap hitung SELALU menghitung ulang), "fasle", perbandingan TEKS `> "30000000000"` (di sini
// numerik); kurs TSI = TREATYEXCHANGEYEARLY per Begin date seperti Cedant (K49-11), bukan getcurrencystandard.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrMasukanScoring - isian Scoring Risk tidak sah. 400.
	ErrMasukanScoring = errors.New("services: isian Scoring Risk tidak sah")
	// ErrScoringBacaSaja - kasus bukan tahap marketing / FlagOnGoingPolicy bukan 0 (CallDualScoringRisk "Read"). 409.
	ErrScoringBacaSaja = errors.New("services: Scoring Risk hanya dapat diisi di tahap marketing dengan FlagOnGoingPolicy 0")
	// ErrScoringTanpaDatabase - pembaca Scoring / objek / spreading / okupasi tidak terpasang. 503.
	ErrScoringTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, Scoring Risk tidak terbaca")
)

// PesanOkupasiSama - `SetDataOccupationScoringRisk_Act` langkah 1 (local.Err), verbatim.
const PesanOkupasiSama = "Can't Choose The Same Occupation"

const (
	// flagScoringUbah - FlagOnGoingPolicy kasus yang boleh mengisi (CallDualScoringRisk "Write").
	flagScoringUbah = "0"
	// lebarIndexLocation - T_DATASCORINGRISKLIST.INDEX_LOCATION = T_PROPERTY.OBJECT_NO VARCHAR2(50).
	lebarIndexLocation = 50
)

// posisiMarketing - padanan `PositionNote = 'ReasFacInMarketing'`: posisi awal case NB "Offer" (butir 77, tahap marketing
// sistem baru) dan nilai Pega apa adanya bila data lama memuatnya `[dugaan]`.
var posisiMarketing = map[string]bool{repository.PosisiAwalCaseNB: true, "ReasFacInMarketing": true}

// poinScoring - TempScoringRisk.<faktor>.ScoreN (`ValueScoringRisk_act` langkah 3-26), urut ChechBox1..n.
var poinScoring = func() map[string][]*apd.Decimal {
	mentah := map[string][]string{
		"occupation":               {"15", "10", "5", "-5", "-15"},
		"riot":                     {"1", "1", "0", "-1"},
		"flood":                    {"5", "5", "0", "-1"},
		"earthquake":               {"2.5", "2.5", "2.25", "2", "1.5", "1", "-1"},
		"tsunami":                  {"1.5", "1.5", "0", "-1"},
		"classConstruction":        {"2.5", "2", "0", "0", "0", "-1"},
		"building":                 {"2.5", "0", "-1", "-1", "-1"},
		"machinery":                {"2.5", "0", "-1", "-1", "-1"},
		"stock":                    {"2.5", "2.5", "0", "-1"},
		"typeOfStock":              {"2.5", "2.5", "0", "-1"},
		"portableFireExtinguisher": {"2", "3", "2", "0"},
		"hydrant":                  {"2", "1.5", "0"},
		"sprinkler":                {"1", "1", "0"},
		"fireAlarm":                {"1.5", "0"},
		"securityGuard":            {"1.5", "1", "0"},
		"fireBrigade":              {"1.5", "1", "0.75", "0.5", "0"},
		"lossRatio":                {"22.5", "17.5", "15", "12.5", "10", "5", "-5", "-10", "-20", "-2.5"},
		"riskImprovement":          {"10", "10", "0"},
		"surveyReport":             {"10", "7.5", "5", "2.5", "1", "0"},
		"rate":                     {"2", "0", "-1"},
		"deductible":               {"2.5", "0", "-1"},
		"totalSumInsured":          {"-1", "1"},
		"otherInformation":         {"2.5", "0", "-1"},
		"clauses":                  {"0", "2", "-1"},
	}
	hasil := map[string][]*apd.Decimal{}
	for _, f := range models.DaftarFaktorScoring {
		p := mentah[f.Kunci]
		if len(p) != f.JumlahOpsi {
			panic("services: poin scoring " + f.Kunci + " tidak sejumlah opsinya")
		}
		for _, s := range p {
			d, _, err := apd.NewFromString(s)
			if err != nil {
				panic(err)
			}
			hasil[f.Kunci] = append(hasil[f.Kunci], d)
		}
	}
	return hasil
}()

var (
	// ambangScoring - 30.000.000.000 IDR (TSI `SetDataScoringRisk_act` 4.1.7 / 4.1.8; TotalRNMShare `ScoringResult` 2.2).
	ambangScoring = apd.New(30000000000, 0)
	// risikoOkupasi - Occupation.ChechBox1..5 = KDRiskExposure "01".."05" (SetDataScoringRisk_act 4.1.11).
	risikoOkupasi = map[string]int{"01": 1, "02": 2, "03": 3, "04": 4, "05": 5}
	// zonaGempa - Earthquake.ChechBox2..6 = Zone "Daerah Zona 1".."5" (SetDataScoringRisk_act langkah 5; ChechBox1 / 7 manual).
	zonaGempa = map[string]int{"Daerah Zona 1": 2, "Daerah Zona 2": 3, "Daerah Zona 3": 4, "Daerah Zona 4": 5, "Daerah Zona 5": 6}
)

// lokasiScoring - satu lokasi sebagaimana dibaca isi otomatis kolom. TSIIDR / ShareIDR nil = kurs mata uang KursTidakAda
// tidak ada.
type lokasiScoring struct {
	ObjectNo, OccupationID, OccupationName, Alamat, Zona, KDRisk string
	TSIIDR, ShareIDR                                             *apd.Decimal
	KursTidakAda                                                 []string
}

// LokasiDominan - satu pilihan Dominant Risk (autocomplete atas LocationList: ObjectNo -> IndexLocation, label
// OccupationList(1).OccupationName).
type LokasiDominan struct{ IndexLocation, Label string }

// TampilanScoring - jawaban seluruh endpoint tab Scoring Risk.
type TampilanScoring struct {
	BolehUbah                           bool
	InsuredName, Status, Mulai, Selesai string
	Scoring                             models.Scoring
	LokasiDominan                       []LokasiDominan
	Pesan                               []string
}

// DenganScoring memasang penyimpan Scoring Risk (tiket 54).
func (s *Service) DenganScoring(p repository.PenyimpanScoring) *Service { s.scoring = p; return s }

// ---- hitung (murni) ----------------------------------------------------------------------------------------------

// pilihanLossRatio - SetLossRatio_Act: persen > 0 -> band ChechBox4..9 (menggantikan pilihan lain - satu centang per
// faktor); persen kosong / ≤ 0 -> ChechBox4..9 dikosongkan, pilihan manual 1-3 / 10 tetap.
func pilihanLossRatio(persen *apd.Decimal, lama int) int {
	if persen == nil || persen.Sign() <= 0 {
		if lama >= 4 && lama <= 9 {
			return 0
		}
		return lama
	}
	for i, batas := range []int64{10, 25, 50, 75, 100} {
		if persen.Cmp(apd.New(batas, 0)) <= 0 {
			return 4 + i
		}
	}
	return 9
}

// hitungKolom - isi otomatis dari lokasi (okupasi dan TSI dikunci; zona gempa hanya bila belum dipilih), band loss ratio,
// ScorePerFactor dan Score. l nil = kolom tanpa lokasi (okupasi / TSI tidak dicentang).
func hitungKolom(m models.KolomScoring, l *lokasiScoring, pesan *[]string, namaKolom string) models.KolomScoring {
	k := models.KolomScoring{Ada: true, IndexLocation: m.IndexLocation, PersenLossRatio: m.PersenLossRatio,
		Pilihan: map[string]int{}, Skor: map[string]*apd.Decimal{}}
	for kunci, n := range m.Pilihan {
		k.Pilihan[kunci] = n
	}
	k.Pilihan[models.KunciOccupation], k.Pilihan[models.KunciTotalSumInsured] = 0, 0
	if l != nil {
		k.OccupationNote, k.OccupationCode, k.TopRiskLocation = l.OccupationName, l.OccupationID, l.Alamat
		k.Pilihan[models.KunciOccupation] = risikoOkupasi[l.KDRisk]
		switch {
		case l.TSIIDR == nil:
			mu := strings.Join(l.KursTidakAda, ", ")
			if mu == "" {
				mu = "mata uang lokasi"
			}
			*pesan = append(*pesan, "Kurs "+mu+" pada Begin date tidak ada di TREATYEXCHANGEYEARLY - Total Sum Insured "+namaKolom+
				" tidak diisi")
		case l.TSIIDR.Cmp(ambangScoring) <= 0:
			k.Pilihan[models.KunciTotalSumInsured] = 1
		default:
			k.Pilihan[models.KunciTotalSumInsured] = 2
		}
		if k.Pilihan[models.KunciEarthquake] == 0 {
			k.Pilihan[models.KunciEarthquake] = zonaGempa[l.Zona]
		}
	}
	k.Pilihan[models.KunciLossRatio] = pilihanLossRatio(m.PersenLossRatio, k.Pilihan[models.KunciLossRatio])
	k.Total = nol()
	for _, f := range models.DaftarFaktorScoring {
		skor := nol()
		if n := k.Pilihan[f.Kunci]; n > 0 {
			skor = poinScoring[f.Kunci][n-1]
		}
		k.Skor[f.Kunci] = skor
		k.Total = tambah(k.Total, skor)
	}
	return k
}

// hitungScoring - dua kolom (kolom 2 hanya bila lokasi dominan ber-okupasi, SetScore_Act langkah 3), FinalScore, hasil.
func hitungScoring(status string, masukan [2]models.KolomScoring, lok [2]*lokasiScoring) (models.Scoring, []string) {
	pesan := []string{}
	sc := models.Scoring{Status: status}
	sc.Kolom[0] = hitungKolom(masukan[0], lok[0], &pesan, "Top Risk")
	if lok[1] != nil && lok[1].OccupationName != "" {
		sc.Kolom[1] = hitungKolom(masukan[1], lok[1], &pesan, "Dominant Risk")
		if sc.Kolom[0].OccupationNote == sc.Kolom[1].OccupationNote {
			pesan = append(pesan, PesanOkupasiSama)
		}
	}
	sc.FinalScore = sc.Kolom[0].Total
	if sc.Kolom[1].Ada {
		// @divide(local.TotalScore, 2, 2) - pembulatan setengah ke atas `[dugaan]`.
		sc.FinalScore = bulat2(bagiSkala(tambah(sc.Kolom[0].Total, sc.Kolom[1].Total), apd.New(2, 0), 20))
	}
	sc.Hasil, sc.NoteFinalScore = hasilScoring(sc.FinalScore, sc.Kolom, lok)
	return sc, pesan
}

// bulat2 - 2 desimal setengah ke atas, nol di belakang dibuang.
func bulat2(d *apd.Decimal) *apd.Decimal { return bulatSetengahAtas(d, 2) }

// bulatSetengahAtas - `desimal` angka di belakang koma setengah ke atas (Pega @divide(x, y, n) `[dugaan]`), nol di belakang
// dibuang; nil tetap nil.
func bulatSetengahAtas(d *apd.Decimal, desimal int32) *apd.Decimal {
	if d == nil {
		return nil
	}
	h := new(apd.Decimal)
	ctx := apd.BaseContext.WithPrecision(60)
	ctx.Rounding = apd.RoundHalfUp
	if _, err := ctx.Quantize(h, d, -desimal); err != nil {
		return d // di luar presisi: nilai asli (pola simpan8), ditolak periksa muat bila melebihi kolom
	}
	h.Reduce(h)
	return h
}

// antara - a ≤ x ≤ b (atasTermasuk) / a ≤ x < b.
func antara(x *apd.Decimal, a, b int64, atasTermasuk bool) bool {
	if x.Cmp(apd.New(a, 0)) < 0 {
		return false
	}
	c := x.Cmp(apd.New(b, 0))
	return c < 0 || atasTermasuk && c == 0
}

// hasilScoring - ScoringResult langkah 1-3 VERBATIM, termasuk urutan evaluasi antar kolom (Primary.ScoringResult dibaca
// dan ditulis bergiliran kolom 1 lalu 2) dan TechnicalDirector yang dijaga OperationalDirector (salin-tempel Pega; hasilnya
// selalu sama dengan OperationalDirector). TotalRNMShare per kolom = premi spread lokasi kolom itu dalam IDR (Pega: satu
// halaman global terakhir diisi `[dugaan]` maksudnya per lokasi), dibandingkan NUMERIK.
func hasilScoring(fs *apd.Decimal, kolom [2]models.KolomScoring, lok [2]*lokasiScoring) (models.HasilScoring, string) {
	var h models.HasilScoring
	h.Good = antara(fs, 86, 100, true)
	h.AboveAverage = antara(fs, 76, 86, false)
	h.Average = antara(fs, 60, 76, false)
	h.BelowAverage = antara(fs, 30, 60, false)
	h.Poor = antara(fs, 0, 30, false)
	catatan := "Poor"
	switch {
	case h.Good:
		catatan = "Good"
	case h.AboveAverage:
		catatan = "Above Average"
	case h.Average:
		catatan = "Average"
	case h.BelowAverage:
		catatan = "Below Average"
	}
	for i, k := range kolom {
		if !k.Ada {
			continue
		}
		p := func(f string, n ...int) bool {
			for _, x := range n {
				if k.Pilihan[f] == x {
					return true
				}
			}
			return false
		}
		shareBesar := lok[i] != nil && lok[i].ShareIDR != nil && lok[i].ShareIDR.Cmp(ambangScoring) > 0
		// Due to
		h.Occupation = h.Occupation || p(models.KunciOccupation, 5, 4)
		h.ObjectConditions = h.ObjectConditions || p("classConstruction", 3, 4) || p("building", 4) || p("machinery", 4)
		h.LossRatio = h.LossRatio || p(models.KunciOccupation, 1, 2) && p(models.KunciLossRatio, 7, 8, 9) || p(models.KunciOccupation, 3, 4) && p(models.KunciLossRatio, 6, 7, 8, 9)
		h.Others = h.Others || p(models.KunciTotalSumInsured, 1)
		// Approval From
		h.OperationalDivHead = h.OperationalDivHead || p(models.KunciOccupation, 4) || h.ObjectConditions || h.LossRatio || h.Others ||
			h.OperationalDirector || h.TechnicalDirector || h.PresidentDirector
		h.OperationalDirector = h.OperationalDirector || p(models.KunciOccupation, 4) && shareBesar || h.ObjectConditions || h.LossRatio ||
			h.Others || h.TechnicalDirector || h.PresidentDirector
		if h.OperationalDirector {
			h.TechnicalDirector = h.OperationalDirector
		} else {
			h.TechnicalDirector = p(models.KunciOccupation, 4) && shareBesar || h.ObjectConditions || h.LossRatio || h.Others || h.PresidentDirector
		}
		h.PresidentDirector = h.PresidentDirector || p(models.KunciOccupation, 5)
	}
	setuju := h.OperationalDirector || h.TechnicalDirector || h.PresidentDirector || h.OperationalDivHead
	h.AcceptableNotApproval = antara(fs, 30, 100, true) && !setuju
	h.AcceptableWithApproval = antara(fs, 30, 100, true) && setuju
	h.NotAcceptable = antara(fs, 0, 30, false)
	return h, catatan
}

// jumlahOpsiScoring - kunci faktor -> banyak ChechBox.
var jumlahOpsiScoring = func() map[string]int {
	m := map[string]int{}
	for _, f := range models.DaftarFaktorScoring {
		m[f.Kunci] = f.JumlahOpsi
	}
	return m
}()

// periksaMasukanScoring - kunci faktor dikenal, nomor 0..JumlahOpsi, IndexLocation ≤ 50 bita (kolom 2).
func periksaMasukanScoring(masukan [2]models.KolomScoring) error {
	var masalah []string
	for i, k := range masukan {
		for kunci, n := range k.Pilihan {
			maks, ada := jumlahOpsiScoring[kunci]
			switch {
			case !ada:
				masalah = append(masalah, fmt.Sprintf("kolom[%d].pilihan.%s tidak dikenal", i, kunci))
			case n < 0 || n > maks:
				masalah = append(masalah, fmt.Sprintf("kolom[%d].pilihan.%s harus 0..%d", i, kunci, maks))
			}
		}
		if len(k.IndexLocation) > lebarIndexLocation {
			masalah = append(masalah, fmt.Sprintf("kolom[%d].indexLocation paling banyak %d bita", i, lebarIndexLocation))
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanScoring, strings.Join(masalah, "; "))
	}
	return nil
}

// ---- alur --------------------------------------------------------------------------------------------------------

// konteksScoring - semua yang dibaca untuk satu hitung.
type konteksScoring struct {
	kasus     models.Kasus
	tersimpan models.KasusScoring
	objek     []models.ObjekFire
	spreading models.KasusSpreading
	risiko    map[string]string
	// mulai / waktu - kinerja baca (scoring, kasus, objek, spreading, okupasi), untuk logScoringLambat.
	mulai time.Time
	waktu [5]time.Duration
	// kurs - sumber kurs satu permintaan (kursMemo: keIDRScoring dipanggil empat kali atas mata uang yang sama); nil = s.spreading.
	kurs sumberKurs
}

func (s *Service) bacaKonteksScoring(ctx context.Context, id string) (konteksScoring, error) {
	var k konteksScoring
	if s.scoring == nil || s.kasus == nil || s.objek == nil || s.spreading == nil || s.occupation == nil {
		return k, ErrScoringTanpaDatabase
	}
	if !idKasusSah(id) {
		return k, ErrKasusTidakAda
	}
	// Laporan work owner 12-10-2026 (tab Scoring "Memuat…" beberapa detik): empat bacaan independen dijalankan PARALEL (masing-masing
	// lewat basis data, tanpa tx bersama); hasil identik dengan urutan lama. Galat = yang PERTAMA terjadi (pembungkus sama dengan
	// alur berurutan lama), bacaan lain dibatalkan lewat ctx turunan.
	k.mulai = time.Now()
	ctxBaca, batal := context.WithCancel(ctx)
	defer batal()
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		pertama error
	)
	jalan := func(i int, f func() error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t0 := time.Now()
			err := f()
			k.waktu[i] = time.Since(t0)
			if err != nil {
				mu.Lock()
				if pertama == nil {
					pertama = err
					batal()
				}
				mu.Unlock()
			}
		}()
	}
	jalan(0, func() (err error) {
		if k.tersimpan, err = s.scoring.BacaScoring(ctxBaca, id); err != nil {
			return galatKasus(err)
		}
		return nil
	})
	jalan(1, func() (err error) {
		if k.kasus, err = s.kasus.BacaKasus(ctxBaca, id); err != nil {
			return galatKasus(err)
		}
		return nil
	})
	jalan(2, func() (err error) {
		k.objek, err = s.BacaObjek(ctxBaca, id)
		return err
	})
	jalan(3, func() (err error) {
		if k.spreading, err = s.spreading.BacaSpreading(ctxBaca, nil, id); err != nil {
			return galatKasus(err)
		}
		return nil
	})
	wg.Wait()
	if pertama != nil {
		return k, pertama
	}
	totalLokasi(&k.spreading)
	var kode []string
	lihat := map[string]bool{}
	for _, o := range k.objek {
		if len(o.Occupations) > 0 && o.Occupations[0].OccupationID != "" && !lihat[o.Occupations[0].OccupationID] {
			lihat[o.Occupations[0].OccupationID] = true
			kode = append(kode, o.Occupations[0].OccupationID)
		}
	}
	t0 := time.Now()
	var err error
	k.risiko, err = s.occupation.RisikoOkupasiFire(ctx, kode)
	k.waktu[4] = time.Since(t0)
	return k, err
}

// logScoringLambat - satu baris log bila konteks + hitung Scoring melewati ambangLogScoring: durasi tiap langkah (kurs = hitung
// kolom, termasuk bacaan kurs keIDRScoring). Tanpa data kasus selain ID.
func logScoringLambat(id string, k konteksScoring, kurs time.Duration) {
	if total := time.Since(k.mulai); total > ambangLogScoring {
		log.Printf("nbfacin: Scoring kasus %s lambat %v: scoring %v, kasus %v, objek %v, spreading %v, okupasi %v, kurs %v",
			id, total, k.waktu[0], k.waktu[1], k.waktu[2], k.waktu[3], k.waktu[4], kurs)
	}
}

// ambangLogScoring - konteks Scoring lebih lama dari ini dicatat per langkah (laporan work owner 12-10-2026).
const ambangLogScoring = 2 * time.Second

// lokasiKolom - lokasi kolom: Top Risk = satu-satunya lokasi, atau lokasi IsTopRisk pertama (SetDataScoringRisk_act
// langkah 4); Dominant = ObjectNo pertama = indexLocation (langkah 3). nil = tidak ada. TSI / premi spread lokasi dalam
// IDR per mata uang × kurs Begin date (totalKeIDR, K49-11); lokasi tanpa TSI dan premi = 0.
func (s *Service) lokasiKolom(ctx context.Context, k konteksScoring, top bool, indexLocation string) (*lokasiScoring, error) {
	var o *models.ObjekFire
	for i := range k.objek {
		if top && (len(k.objek) == 1 || k.objek[i].IsTopRisk) || !top && indexLocation != "" && k.objek[i].ObjectNo == indexLocation {
			o = &k.objek[i]
			break
		}
	}
	if o == nil {
		return nil, nil
	}
	l := &lokasiScoring{ObjectNo: o.ObjectNo, Alamat: o.RiskLocation}
	if len(o.Occupations) > 0 {
		l.OccupationID, l.OccupationName = o.Occupations[0].OccupationID, o.Occupations[0].OccupationName
		l.KDRisk = k.risiko[l.OccupationID]
	}
	for _, it := range o.Items {
		for _, c := range it.Coverages {
			if c.Zone != "" {
				l.Zona = c.Zone // coverage TERAKHIR berzona (Pega menimpa berurutan)
			}
		}
	}
	var jumlah []jumlahMataUang
	indeks := map[string]int{}
	tambahJumlah := func(mu string, tsi, premi *apd.Decimal) {
		if mu == "" {
			mu = mataUangKosong
		}
		n, ada := indeks[mu]
		if !ada {
			n, indeks[mu] = len(jumlah), len(jumlah)
			jumlah = append(jumlah, jumlahMataUang{Currency: mu, TSIRNM: nol(), Premi: nol()})
		}
		jumlah[n].TSIRNM, jumlah[n].Premi = tambah(jumlah[n].TSIRNM, tsi), tambah(jumlah[n].Premi, premi)
	}
	for _, t := range o.TotalPerCurrency {
		tambahJumlah(t.Currency, t.TSI, nil)
	}
	for _, sl := range k.spreading.Lokasi {
		if sl.ObjectNo == o.ObjectNo {
			for _, t := range sl.Total {
				tambahJumlah(t.Currency, nil, t.PremiumSpreaded)
			}
			break
		}
	}
	// TSI dan premi spread dikonversi TERPISAH: kurs mata uang yang hanya ada di premi tidak mengosongkan TSI (dan sebaliknya).
	var err error
	src := k.kurs
	if src == nil {
		src = s.spreading
	}
	if l.TSIIDR, l.KursTidakAda, err = keIDRScoring(ctx, src, k.spreading.StartDateTime, jumlah, true); err != nil {
		return nil, err
	}
	l.ShareIDR, _, err = keIDRScoring(ctx, src, k.spreading.StartDateTime, jumlah, false)
	return l, err
}

// bolehUbahScoring - CallDualScoringRisk "Write": tahap marketing dan FlagOnGoingPolicy 0.
func bolehUbahScoring(k models.KasusScoring) bool {
	return tahapMarketing(models.TahapKasus{Position: k.Position, FlagOnGoingPolicy: k.FlagOnGoingPolicy})
}

// keIDRScoring - Σ TSI (tsi true) atau premi spread mata uang yang bernilai bukan nol × kurs (totalKeIDR); tanpa mata uang
// bernilai -> 0; satu kurs tidak ada -> nil + mata uangnya.
func keIDRScoring(ctx context.Context, src sumberKurs, mulai string, jumlah []jumlahMataUang, tsi bool) (*apd.Decimal, []string, error) {
	var pilih []jumlahMataUang
	for _, j := range jumlah {
		v := j.Premi
		if tsi {
			v = j.TSIRNM
		}
		if nolBila(v).Sign() != 0 {
			pilih = append(pilih, jumlahMataUang{Currency: j.Currency, TSIRNM: v})
		}
	}
	if len(pilih) == 0 {
		return nol(), nil, nil
	}
	hasil, _, tidakAda, err := totalKeIDR(ctx, src, mulai, pilih)
	return hasil, tidakAda, err
}

// hitungScoringKasus - konteks + masukan -> tampilan (tanpa menyimpan).
func (s *Service) hitungScoringKasus(ctx context.Context, id string, k konteksScoring, masukan [2]models.KolomScoring) (TampilanScoring, error) {
	if k.kurs == nil {
		k.kurs = kursMemoBaru(s.spreading) // satu permintaan
	}
	var lok [2]*lokasiScoring
	var err error
	if lok[0], err = s.lokasiKolom(ctx, k, true, ""); err != nil {
		return TampilanScoring{}, err
	}
	if lok[1], err = s.lokasiKolom(ctx, k, false, masukan[1].IndexLocation); err != nil {
		return TampilanScoring{}, err
	}
	if masukan[1].IndexLocation != "" && lok[1] == nil {
		return TampilanScoring{}, fmt.Errorf("%w: kolom[1].indexLocation %q bukan ObjectNo lokasi kasus", ErrMasukanScoring,
			masukan[1].IndexLocation)
	}
	status := "Renewal"
	if strings.Contains(id, "NB") {
		status = "New Business"
	}
	sc, pesan := hitungScoring(status, masukan, lok)
	if lok[1] != nil && lok[1].OccupationName == "" {
		// Pega membuang kolom Dominant ber-OccupationNote kosong tanpa pesan; layar diberi tahu.
		pesan = append(pesan, "Dominant Risk lokasi "+lok[1].ObjectNo+" belum punya Occupation - kolom Dominant Risk tidak dihitung")
	}
	t := TampilanScoring{BolehUbah: bolehUbahScoring(k.tersimpan),
		InsuredName: k.kasus.InsuredName, Status: status, Mulai: waktuKeKabel(k.kasus.General.StartDateTime),
		Selesai: waktuKeKabel(k.kasus.General.EndDateTime), Scoring: sc, LokasiDominan: []LokasiDominan{}, Pesan: pesan}
	for _, o := range k.objek {
		d := LokasiDominan{IndexLocation: o.ObjectNo}
		if len(o.Occupations) > 0 {
			d.Label = o.Occupations[0].OccupationName
		}
		t.LokasiDominan = append(t.LokasiDominan, d)
	}
	return t, nil
}

// TampilanScoringKasus - GET: masukan tersimpan, dihitung ulang (keputusan work owner - skor Pega bisa basi).
func (s *Service) TampilanScoringKasus(ctx context.Context, id string) (TampilanScoring, error) {
	k, err := s.bacaKonteksScoring(ctx, id)
	if err != nil {
		return TampilanScoring{}, err
	}
	t0 := time.Now()
	hasil, err := s.hitungScoringKasus(ctx, id, k, k.tersimpan.Kolom)
	logScoringLambat(id, k, time.Since(t0))
	return hasil, err
}

// HitungScoring - POST …/scoring/hitung: tanpa menyimpan; kasus baca-saja -> 409.
func (s *Service) HitungScoring(ctx context.Context, id string, masukan [2]models.KolomScoring) (TampilanScoring, error) {
	if err := periksaMasukanScoring(masukan); err != nil {
		return TampilanScoring{}, err
	}
	k, err := s.bacaKonteksScoring(ctx, id)
	if err != nil {
		return TampilanScoring{}, err
	}
	if !bolehUbahScoring(k.tersimpan) {
		return TampilanScoring{}, ErrScoringBacaSaja
	}
	return s.hitungScoringKasus(ctx, id, k, masukan)
}

// SimpanScoring - PUT …/scoring: identitas (401) -> isian (400) -> baca-saja (409) -> okupasi sama (400, pesan Pega
// menghalangi Save) -> hitung dan simpan utuh.
func (s *Service) SimpanScoring(ctx context.Context, pelaku inti.Pelaku, id string, masukan [2]models.KolomScoring) (TampilanScoring, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanScoring{}, err
	}
	if err := s.wajibBolehUbah(ctx, pelaku, id); err != nil {
		return TampilanScoring{}, err
	}
	t, err := s.HitungScoring(ctx, id, masukan)
	if err != nil {
		return TampilanScoring{}, err
	}
	for _, p := range t.Pesan {
		if p == PesanOkupasiSama {
			return TampilanScoring{}, fmt.Errorf("%w: %s", ErrMasukanScoring, PesanOkupasiSama)
		}
	}
	if s.transaksi == nil {
		return TampilanScoring{}, ErrScoringTanpaDatabase
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error { return s.scoring.TulisScoring(ctx, tx, id, t.Scoring) })
	if err != nil {
		return TampilanScoring{}, galatKasus(err)
	}
	return t, nil
}
