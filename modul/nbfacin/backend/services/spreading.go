package services

// Tab Spreading kasus FIRE (tiket 48). Port jalur NB FIRE rantai Pega `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\` (dibaca
// 05-10-2026; rincian per langkah di docs/issues/48):
//   - hitungNR            = CountPremiAndTSIRNMFireMBU_ACT langkah 6 (TSINusantaraRe, PremiNusantaraRe, TotalPremiumNusantaraRe);
//   - spreadingOtomatis   = GetKapasitasTreaty cabang IsFire (QS 100 %, atau QS + SPL dipotong KAPASITAS_TREATY);
//   - salinTemplate       = CopyToAllSpreadingFire_ACT (bukan EDM);
//   - terapkanPersen      = DT SetTSIPremiSpreaded_FacIn mode "percent";
//   - hitungTotal         = SumTSIPremiSpreadedRNM_FIRE_Act cabang 3.5 (tanpa top risk) / 3.6 (ada top risk) +
//                           SumTSIPremiSpreadedRNM_Act langkah 18-22 (kurs, TotalSpreadingCurrency, batas treaty, TotalSpreadAll);
//   - daftarTreaty        = GetTreatyName (SpreadingList1 = dropdown; GetTreatyName.pxResults = batas treaty);
//   - periksaJenisTreaty  = cekSpreadingFactIn langkah 8-9.
// Skala @Math.divide(x, y, 20) = bagi (setengah-ke-atas, skala 20); disimpan 8 desimal (simpan8).
// Di luar cakupan (K48-6): cabang EDM / Group, syariah (Tabarufund; di Pega pun mati), layering CoverageBasis 5, daftar ID
// case hardcode Pega, PositionNote (Marketing / Team Leader), FlagDelete.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrMasukanSpreading - isian spreading tidak sah. 400.
	ErrMasukanSpreading = errors.New("services: isian spreading tidak sah")
	// ErrSpreadingBerubah - bentuk lokasi / item / coverage badan PUT tidak sama dengan objek tersimpan. 409.
	ErrSpreadingBerubah = errors.New("Objek berubah sejak spreading dimuat - muat ulang tab Spreading")
	// ErrSpreadingTanpaDatabase - tabel spreading / master treaty tidak terbaca. 503.
	ErrSpreadingTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, spreading tidak terbaca")
)

const (
	// Pesan verbatim Pega.
	pesanShareLebih100   = "Share Percentage can't be more than 100"        // cekSpreadingFactIn(_Act)
	pesanKapasitasTop    = "Total TSI RNM Top Risk exceeds treaty capacity" // GetKapasitasTreaty
	pesanTanpaQS         = "Can not proceed spreading without QS"           // CalcultePersentageSpeading_Act 2.17
	pesanJenisKembar     = "Treaty Type can't be same"                      // CalcultePersentageSpeading_Act 2.14
	pesanJenisTreaty     = "Invalid Treaty Type Spreading!"                 // cekSpreadingFactIn langkah 3 / 9
	pesanObjekKosong     = "Object can't be empty!"                         // IsThereAnyObjectLocation_Act langkah 7
	catatanTerorisme     = "TERRORISM & SABOTAGE"                           // SumTSIPremiSpreadedRNM_FIRE_Act .4.4
	treatyORS, namaORS   = "10007", "ORS"                                   // GetTreatyName langkah 8
	limitORS             = "150000000000"                                   // GetTreatyName langkah 8 (CARI3)
	treatyFACOUT         = "10015"                                          // GetTreatyName langkah 19
	namaFACOUT           = "FACOUT"                                         //
	idMataUangIDR        = "10026"                                          // SumTSIPremiSpreadedRNM_Act 19.4 "SET KURS = 1, WHEN IDR"
	namaIDR              = "IDR"                                            // GetKapasitasTreaty @if(.Currency=="IDR")
	lebarTreatyType      = 50                                               // T_SPREADINGLIST.TREATY_TYPE
	lebarTreatyName      = 500                                              // T_SPREADINGLIST.TREATY_NAME
	banyakBarisSpreading = 50                                               // batas baris per coverage / template (penjaga badan)
)

var (
	// batasORS150 / batasORS1815 - SumTSIPremiSpreadedRNM_Act 20.2 / 20.3. Mata uang tidak tertulis di korpus - `[dugaan]` IDR
	// sesudah kurs, `[pertanyaan terbuka]` sejenis OQ-046 (nilai hardcode tanpa mata uang; belum ada OQ sendiri untuk batas
	// ORS). Arti kode 10007 / 10015 / 10026 hanya dari keterangan langkah Pega - sejenis OQ-020 (arti kode).
	batasORS150  = apd.New(150000000000, 0)
	batasORS1815 = apd.New(181500000000, 0)
	// mulaiORS1815 - "20260630T170000.000 GMT" (1 Juli 2026 00:00 WIB) - 20.3.
	mulaiORS1815 = time.Date(2026, 6, 30, 17, 0, 0, 0, time.UTC)
	// jendelaSFHRE - GetTreatyName 10.2 (SF-HRE dipakai FIRE hanya 1-Jun-2026 .. 31-Mei-2027).
	jendelaSFHREMulai = time.Date(2026, 6, 1, 5, 0, 0, 0, time.UTC)
	jendelaSFHREAkhir = time.Date(2027, 5, 31, 5, 0, 0, 0, time.UTC)
	sepuluhRibu       = apd.New(10000, 0)
)

// DenganSpreading memasang penyimpan spreading (tiket 48).
func (s *Service) DenganSpreading(p repository.PenyimpanSpreading) *Service {
	s.spreading = p
	return s
}

// TampilanSpreading - jawaban seluruh endpoint tab Spreading.
type TampilanSpreading struct {
	PercentShare      *apd.Decimal
	Treaty            []models.TreatySpreading // ID / Name (dropdown SpreadingList1)
	Template          []models.TemplateSpreading
	Lokasi            []models.LokasiSpreading
	RingkasanTreaty   []models.TotalSpreading
	RingkasanMataUang []models.RingkasanMataUangSpreading
	Pesan             []string
}

// ---- aritmetika -------------------------------------------------------------------------------------------------

func nol() *apd.Decimal { return apd.New(0, 0) }

func tambah(a, b *apd.Decimal) *apd.Decimal {
	h := new(apd.Decimal)
	konteksCoverage.Add(h, nolBila(a), nolBila(b))
	return h
}

func kurang(a, b *apd.Decimal) *apd.Decimal {
	h := new(apd.Decimal)
	konteksCoverage.Sub(h, nolBila(a), nolBila(b))
	return h
}

// bagiSkala - @Math.divide(a, b, 20); pembagi nol -> 0 (Pega: hasil kosong -> 0 di penjumlahan berikutnya `[dugaan]`).
func bagiSkala(a, b *apd.Decimal, skala int32) *apd.Decimal {
	h, err := bagi(nolBila(a), nolBila(b), skala)
	if err != nil {
		return nol()
	}
	return h
}

// persenDari - @Math.divide(x*y, 100, 20).
func persenDari(x, y *apd.Decimal) *apd.Decimal {
	return bagiSkala(kali(nolBila(x), nolBila(y)), seratus, 20)
}

// lebihBesar - a > b (nil = 0).
func lebihBesar(a, b *apd.Decimal) bool { return nolBila(a).Cmp(nolBila(b)) > 0 }

// ---- tanggal ----------------------------------------------------------------------------------------------------

// waktuMulai - Begin date teks Pega -> waktu; false bila kosong / tak terurai.
func waktuMulai(k models.KasusSpreading) (time.Time, bool) {
	t, err := time.Parse(BentukWaktuPega, k.StartDateTime)
	return t, err == nil
}

// tanggalMulai - @FormatDateTime(StartDateTime, "dd/MM/yyyy", "Asia/Jakarta") (GetKapasitasTreaty langkah 4,
// GetTreatyName langkah 2-3); "" bila Begin date kosong.
func tanggalMulai(k models.KasusSpreading) string {
	t, ok := waktuMulai(k)
	if !ok {
		return ""
	}
	return t.In(WIB).Format("02/01/2006")
}

// ---- CountPremiAndTSIRNMFireMBU_ACT langkah 6 ------------------------------------------------------------------

// hitungNR - TSINusantaraRe = PS × TSILiability / 100; PremiNusantaraRe = PS × Premium / 100 − PS × Discount / 100;
// TotalPremiumNusantaraRe item = Σ PremiNusantaraRe.
func hitungNR(k *models.KasusSpreading) {
	ps := nolBila(k.PercentShare)
	for i := range k.Lokasi {
		for j := range k.Lokasi[i].Items {
			it := &k.Lokasi[i].Items[j]
			total := nol()
			for m := range it.Coverages {
				c := &it.Coverages[m]
				c.TSINusantaraRe = persenDari(ps, c.TSILiability)
				c.PremiNusantaraRe = kurang(persenDari(ps, c.Premium), persenDari(ps, c.Discount))
				total = tambah(total, c.PremiNusantaraRe)
			}
			it.TotalPremiumNusantaraRe = total
		}
	}
}

// ---- GetKapasitasTreaty cabang IsFire ---------------------------------------------------------------------------

// sumberKapasitas - master yang dibaca GetKapasitasTreaty.
type sumberKapasitas interface {
	KapasitasTreaty(ctx context.Context, tsi *apd.Decimal, tanggal string) (models.KapasitasTreaty, bool, error)
	KursTerbaru(ctx context.Context, idMataUang string) (*apd.Decimal, error)
	IDMataUang(ctx context.Context, nama string) (string, error)
}

// kursKeIDR - GetCurrencyToIDR_SQL per mata uang item (cache); ok=false bila mata uang bukan IDR dan kursnya tidak ada.
func kursKeIDR(ctx context.Context, src sumberKapasitas, cache map[string]*apd.Decimal, mataUang string) (*apd.Decimal, bool, error) {
	if k, ada := cache[mataUang]; ada {
		return k, k != nil || mataUang == namaIDR, nil
	}
	id, err := src.IDMataUang(ctx, mataUang)
	if err != nil {
		return nil, false, err
	}
	var kurs *apd.Decimal
	if id != "" {
		if kurs, err = src.KursTerbaru(ctx, id); err != nil {
			return nil, false, err
		}
	}
	if kurs != nil && kurs.IsZero() {
		kurs = nil
	}
	cache[mataUang] = kurs
	return kurs, kurs != nil || mataUang == namaIDR, nil
}

// spreadingOtomatis - GetKapasitasTreaty cabang IsFire (lihat docs/issues/48 K48-3): SELURUH spreading coverage diganti.
// Pesan kapasitas verbatim; berhenti di langkah END seperti Pega. KAPASITAS_TREATY tanpa baris = Exit-Activity (tanpa
// pesan di Pega; di sini pesan agent K48-4).
func spreadingOtomatis(ctx context.Context, src sumberKapasitas, k *models.KasusSpreading) ([]string, error) {
	ps := nolBila(k.PercentShare)
	tanggal := tanggalMulai(*k)
	cache := map[string]*apd.Decimal{}
	// Kurs seluruh item diperiksa lebih dulu: tanpa kurs, spreading tersimpan dibiarkan utuh (K48-4).
	for _, l := range k.Lokasi {
		for _, it := range l.Items {
			_, ok, err := kursKeIDR(ctx, src, cache, it.Currency)
			if err != nil {
				return nil, err
			}
			if !ok {
				return []string{"Kurs " + it.Currency + " ke IDR tidak ada di TREATYEXCHANGEYEARLY - spreading otomatis tidak dibuat"}, nil
			}
		}
	}
	// Pra-pindai (DDL 6.1): TSI RNM 100 % dan top risk dalam IDR, coverage pertama tiap item; spreading lama dibuang.
	tsiRNM100, tsiTopRisk, tsiRNMTopRisk1 := nol(), nol(), nol()
	for i := range k.Lokasi {
		l := &k.Lokasi[i]
		for j := range l.Items {
			it := &l.Items[j]
			kurs, _, err := kursKeIDR(ctx, src, cache, it.Currency)
			if err != nil {
				return nil, err
			}
			for m := range it.Coverages {
				c := &it.Coverages[m]
				if m == 0 {
					dasar := nolBila(c.TSILiability)
					if it.Currency != namaIDR {
						dasar = kali(dasar, kurs)
					}
					tsiRNM100 = tambah(tsiRNM100, persenDari(dasar, ps))
					if l.IsTopRisk {
						tsiTopRisk = tambah(tsiTopRisk, dasar)
						tsiRNMTopRisk1 = tambah(tsiRNMTopRisk1, persenDari(dasar, ps))
					}
				}
				c.Spreading = []models.BarisSpreading{}
			}
		}
	}
	tsiRNMTopRisk, pctQS := nol(), nol()
	for i := range k.Lokasi {
		l := &k.Lokasi[i]
		var templat []models.BarisSpreading // SpreadingList dibuang tiap lokasi (DDL 6.2.1)
		for j := range l.Items {
			it := &l.Items[j]
			kurs, _, err := kursKeIDR(ctx, src, cache, it.Currency)
			if err != nil {
				return nil, err
			}
			for m := range it.Coverages {
				c := &it.Coverages[m]
				tsi := nolBila(c.TSILiability)
				if l.IsTopRisk {
					tsi = tsiTopRisk
				}
				tsiRNM := persenDari(c.TSILiability, ps)
				tsiGross := persenDari(c.TSI, ps)
				premi := nolBila(c.PremiNusantaraRe)
				kap, ada, err := src.KapasitasTreaty(ctx, simpan8(tsi), tanggal)
				if err != nil {
					return nil, err
				}
				if !ada {
					return []string{"Kapasitas treaty (KAPASITAS_TREATY) untuk TSI " + utils.FormatDecimal(simpan8(tsi)) + " pada " +
						tanggal + " tidak ada - spreading otomatis berhenti"}, nil
				}
				maxQS, maxSPL := nolBila(kap.MaxLimitQSIDR), nolBila(kap.MaxLimitSPLIDR)
				if l.IsTopRisk {
					tsiRNMTopRisk = tsiRNMTopRisk1
					if it.Currency != namaIDR {
						tsiRNMTopRisk = bagiSkala(tsiRNMTopRisk1, kurs, 20)
					}
				}
				if it.Currency != namaIDR {
					maxQS, maxSPL = bagiSkala(maxQS, kurs, 20), bagiSkala(maxSPL, kurs, 20)
					tsiRNM100 = bagiSkala(tsiRNM100, kurs, 20) // DDL 6.2.4.3.8 - dibagi ulang tiap coverage, seperti Pega
				}
				if len(templat) == 0 {
					templat = []models.BarisSpreading{{TreatyName: kap.TreatyNameQS, TreatyType: kap.IDTreatyQS, SharePercentage: apd.New(100, 0)}}
				}
				c.Spreading = salinBaris(templat)
				maxTSI := tambah(maxQS, maxSPL)
				if lebihBesar(tsiRNMTopRisk, maxTSI) {
					return []string{pesanKapasitasTop}, nil
				}
				pecah := func(dasar *apd.Decimal) {
					pctQS = kali(bagiSkala(maxQS, dasar, 20), seratus)
					pctSPL := bagiSkala(kurang(seratus, pctQS), apd.New(1, 0), 20)
					if lebihBesar(dasar, maxTSI) {
						pctSPL = kali(bagiSkala(maxSPL, dasar, 20), seratus)
					}
					templat = []models.BarisSpreading{
						{TreatyName: kap.TreatyNameQS, TreatyType: kap.IDTreatyQS, SharePercentage: pctQS, TSISpreaded: maxQS},
						{TreatyName: kap.TreatyNameSPL, TreatyType: kap.IDTreatySPL, SharePercentage: pctSPL},
					}
					c.Spreading = salinBaris(templat)
				}
				if l.IsTopRisk && lebihBesar(tsiRNMTopRisk, maxQS) {
					pecah(tsiRNMTopRisk)
				}
				if !l.IsTopRisk && lebihBesar(tsiRNM100, maxQS) {
					pecah(tsiRNM100)
				}
				// Tiap baris: TSI / premi / gross menurut share; QS dipotong maxQS, SPL ditambahkan (DDL 6.2.4.3.19 / .20).
				pembagiQS := tsiRNM
				if l.IsTopRisk {
					pembagiQS = tsiRNMTopRisk
				}
				dipotong := false
				n := len(c.Spreading)
				for s := 0; s < n; s++ {
					b := &c.Spreading[s]
					isiBaris(b, b.SharePercentage, tsiRNM, premi, tsiGross)
					if lebihBesar(b.TSISpreaded, maxQS) {
						dipotong = true
						pctQS = kali(bagiSkala(maxQS, pembagiQS, 20), seratus)
						isiBaris(b, pctQS, tsiRNM, premi, tsiGross)
						b.TSISpreaded = maxQS
					}
					tambahSPL := dipotong && lebihBesar(tsiRNM, maxTSI)
					if l.IsTopRisk {
						tambahSPL = dipotong && !lebihBesar(tsiRNMTopRisk, maxQS)
					}
					if tambahSPL {
						c.Spreading = append(c.Spreading, barisSPL(kap, pctQS, tsiRNM, pembagiQS, maxSPL, premi, tsiGross, l.IsTopRisk, tsiRNMTopRisk))
					}
					if !l.IsTopRisk && lebihBesar(tsiRNM, maxTSI) {
						return []string{"TSI in Location " + strconv.Itoa(i+1) + " exceeds treaty capacity"}, nil
					}
				}
			}
		}
	}
	return nil, nil
}

// isiBaris - TSISpreaded / PremiumSpreaded / TSIGrossSpreaded = share × dasar / 100 (skala 20).
func isiBaris(b *models.BarisSpreading, share, tsiRNM, premi, tsiGross *apd.Decimal) {
	b.SharePercentage = share
	b.TSISpreaded = persenDari(share, tsiRNM)
	b.PremiumSpreaded = persenDari(share, premi)
	b.TSIGrossSpreaded = persenDari(share, tsiGross)
}

// barisSPL - baris SPL tambahan GetKapasitasTreaty (DDL 6.2.4.3.19.4 tanpa top risk / .20.4 dengan top risk, verbatim:
// top risk memakai TSIRNMTopRisk bila TSI SPL tidak melebihi maxSPL). Baris yang ditambahkan tidak ikut diiterasi lagi
// (For Each Embedded Page; bila ikut, Pega berputar tanpa henti saat TSIRNM > maxtsi - `[dugaan]`, K48-3).
func barisSPL(kap models.KapasitasTreaty, pctQS, tsiRNM, pembagi, maxSPL, premi, tsiGross *apd.Decimal, topRisk bool,
	tsiRNMTopRisk *apd.Decimal) models.BarisSpreading {
	pctSPL := bagiSkala(kurang(seratus, pctQS), apd.New(1, 0), 20)
	b := models.BarisSpreading{TreatyName: kap.TreatyNameSPL, TreatyType: kap.IDTreatySPL, SharePercentage: pctSPL}
	tsi := persenDari(pctSPL, tsiRNM)
	switch {
	case lebihBesar(tsi, maxSPL):
		tsi = maxSPL
	case topRisk:
		tsi = tsiRNMTopRisk
	}
	if tsi.Cmp(nolBila(maxSPL)) == 0 {
		pctSPL = kali(bagiSkala(maxSPL, pembagi, 20), seratus)
	}
	b.TSISpreaded = tsi
	if !topRisk {
		b.SharePercentage = pctSPL // .20.4 menulis ulang share; .21.4 (top risk) tidak - share tetap 100 − pctQS (verbatim)
	}
	b.PremiumSpreaded = persenDari(pctSPL, premi)
	b.TSIGrossSpreaded = persenDari(pctSPL, tsiGross)
	return b
}

func salinBaris(b []models.BarisSpreading) []models.BarisSpreading {
	return append([]models.BarisSpreading{}, b...)
}

// ---- CopyToAllSpreadingFire_ACT / DT SetTSIPremiSpreaded_FacIn "percent" ---------------------------------------

// salinTemplate - setiap coverage mendapat template: TSISpreaded = share × (TSILiability × PS / 100) / 100, PremiumSpreaded
// = share × PremiNusantaraRe / 100, TSIGrossSpreaded = share × (TSI × PS / 100) / 100.
func salinTemplate(k *models.KasusSpreading, templat []models.TemplateSpreading) {
	ps := nolBila(k.PercentShare)
	setiapCoverage(k, func(_ *models.ItemSpreading, c *models.CoverageSpreading) {
		c.Spreading = make([]models.BarisSpreading, 0, len(templat))
		for _, t := range templat {
			b := models.BarisSpreading{TreatyType: t.TreatyType, TreatyName: t.TreatyName}
			isiBaris(&b, t.SharePercentage, persenDari(c.TSILiability, ps), c.PremiNusantaraRe, persenDari(c.TSI, ps))
			c.Spreading = append(c.Spreading, b)
		}
	})
}

// terapkanPersen - DT mode "percent": TSISpreaded = share × TSINusantaraRe / 100, PremiumSpreaded = share ×
// PremiNusantaraRe / 100; TSIGrossSpreaded dihitung ulang seperti cekSpreadingFactIn (share × TSI × PS / 10000, K48-5).
func terapkanPersen(k *models.KasusSpreading) {
	ps := nolBila(k.PercentShare)
	setiapCoverage(k, func(_ *models.ItemSpreading, c *models.CoverageSpreading) {
		for s := range c.Spreading {
			b := &c.Spreading[s]
			b.TSISpreaded = persenDari(b.SharePercentage, c.TSINusantaraRe)
			b.PremiumSpreaded = persenDari(b.SharePercentage, c.PremiNusantaraRe)
			b.TSIGrossSpreaded = bagiSkala(kali(nolBila(b.SharePercentage), nolBila(c.TSI), ps), sepuluhRibu, 20)
		}
	})
}

func setiapCoverage(k *models.KasusSpreading, fn func(*models.ItemSpreading, *models.CoverageSpreading)) {
	for i := range k.Lokasi {
		for j := range k.Lokasi[i].Items {
			it := &k.Lokasi[i].Items[j]
			for m := range it.Coverages {
				fn(it, &it.Coverages[m])
			}
		}
	}
}

// ---- SumTSIPremiSpreadedRNM_FIRE_Act ---------------------------------------------------------------------------

// cariTotal - baris ber-(treatyType, currency); -1 bila tidak ada.
func cariTotal(t []models.TotalSpreading, jenis, mataUang string) int {
	for i := range t {
		if t[i].TreatyType == jenis && t[i].Currency == mataUang {
			return i
		}
	}
	return -1
}

// totalLokasi - TotalTSIPremiSpreadRNM tiap lokasi (cabang 3.5 bila tidak satu pun lokasi top risk, 3.6 bila ada) dan
// ClaimEstimation baris spreading. Mengembalikan apakah ada top risk.
func totalLokasi(k *models.KasusSpreading) bool {
	ps := nolBila(k.PercentShare)
	adaTop := false
	for _, l := range k.Lokasi {
		adaTop = adaTop || l.IsTopRisk
	}
	for i := range k.Lokasi {
		l := &k.Lokasi[i]
		l.Total = []models.TotalSpreading{}
		for j := range l.Items {
			it := &l.Items[j]
			for m := range it.Coverages {
				c := &it.Coverages[m]
				pertama := m == 0
				pctFL := nol()
				if fl := nolBila(c.FirstLoss); fl.Sign() > 0 && fl.Cmp(seratus) < 0 {
					pctFL = fl
				}
				tsiLoL := persenDari(c.LimitOfLiability, ps)
				tsiGross := persenDari(c.TSI, ps)
				tsiFL := bagiSkala(kali(nolBila(c.TSI), pctFL, ps), sepuluhRibu, 20)
				teror := c.CoverageNote == catatanTerorisme
				for s := range c.Spreading {
					b := &c.Spreading[s]
					share := nolBila(b.SharePercentage)
					tsiObj, lolSpread, tsiTop := nol(), nol(), nol()
					flSpread := persenDari(tsiFL, share)
					if !adaTop {
						tsiObj = nolBila(b.TSISpreaded)
						lolSpread = persenDari(tsiLoL, share)
						b.ClaimEstimation = lolSpread
					} else {
						if pertama {
							tsiObj = nolBila(b.TSISpreaded)
						}
						lolSpread = persenDari(tsiLoL, share)
						if !teror {
							b.ClaimEstimation = lolSpread
						}
					}
					if pctFL.Sign() > 0 && pertama {
						tsiObj = persenDari(tsiGross, share)
					}
					if pertama && l.IsTopRisk {
						tsiTop = nolBila(b.TSISpreaded)
					}
					if n := cariTotal(l.Total, b.TreatyType, it.Currency); n >= 0 {
						t := &l.Total[n]
						t.PremiumSpreaded = tambah(t.PremiumSpreaded, b.PremiumSpreaded)
						// .4.8.2 (coverage pertama) dan .4.8.3 (terorisme) langkah terpisah: coverage pertama ber-catatan
						// terorisme dijumlah DUA kali (verbatim).
						for _, kena := range []bool{pertama, teror} {
							if kena {
								t.TSISpreaded = tambah(t.TSISpreaded, tsiObj)
								t.ClaimEstimation = tambah(t.ClaimEstimation, lolSpread)
								t.ClaimAmountIDR = tambah(t.ClaimAmountIDR, flSpread)
							}
						}
						if pertama && l.IsTopRisk {
							t.ClaimSpreaded = tambah(t.ClaimSpreaded, tsiTop)
						}
						if share.Cmp(nolBila(t.SharePercentage)) != 0 {
							t.SharePercentage = nol() // "Pct Share Beda"
						}
						continue
					}
					l.Total = append(l.Total, models.TotalSpreading{Currency: it.Currency, TreatyType: b.TreatyType,
						SharePercentage: share, ClaimSpreaded: tsiTop, TSISpreaded: tsiObj, ClaimEstimation: lolSpread,
						ClaimAmountIDR: flSpread, PremiumSpreaded: nolBila(b.PremiumSpreaded)})
				}
			}
		}
	}
	return adaTop
}

// totalMataUang - TotalSpreading per mata uang × treaty dari total lokasi (3.5.2.4 / 3.6.2.4): tanpa top risk seluruh
// medan dijumlah; dengan top risk ClaimEstimation hanya dari lokasi top risk. Urutan mata uang = kemunculan pertama di
// lokasi (Pega: urutan OfferFacIn.CurrencyList, yang tidak dibaca tab ini - `[dugaan]` sama).
func totalMataUang(k models.KasusSpreading, adaTop bool) []models.TotalSpreading {
	hasil := []models.TotalSpreading{}
	for _, l := range k.Lokasi {
		for _, t := range l.Total {
			n := cariTotal(hasil, t.TreatyType, t.Currency)
			if n < 0 {
				baru := t
				if adaTop && !l.IsTopRisk {
					baru.ClaimEstimation = nol()
				}
				hasil = append(hasil, baru)
				continue
			}
			h := &hasil[n]
			h.PremiumSpreaded = tambah(h.PremiumSpreaded, t.PremiumSpreaded)
			h.TSISpreaded = tambah(h.TSISpreaded, t.TSISpreaded)
			h.ClaimSpreaded = tambah(h.ClaimSpreaded, t.ClaimSpreaded)
			h.ClaimAmountIDR = tambah(h.ClaimAmountIDR, t.ClaimAmountIDR)
			if !adaTop || l.IsTopRisk {
				h.ClaimEstimation = tambah(h.ClaimEstimation, t.ClaimEstimation)
			}
			if nolBila(h.SharePercentage).Cmp(nolBila(t.SharePercentage)) != 0 {
				h.SharePercentage = nol() // 3.5.2.4.2.2 / 3.6 (pengecualian case NB-80462 tidak dibawa, K48-6)
			}
		}
	}
	return hasil
}

// ---- GetTreatyName ----------------------------------------------------------------------------------------------

// daftarTreaty - GetTreatyName: (dropdown SpreadingList1, batas GetTreatyName.pxResults). Baris master + ORS; "SF-HRE"
// dibuang di luar jendela 1-Jun-2026..31-Mei-2027; ganda nama -> kemunculan TERAKHIR dipertahankan (langkah Java 15/17);
// dropdown tanpa nama ber-"TRT", + FACOUT.
func daftarTreaty(master []models.TreatySpreading, mulai time.Time, adaMulai bool) (dropdown, batas []models.TreatySpreading) {
	semua := append(append([]models.TreatySpreading{}, master...), models.TreatySpreading{ID: treatyORS, Name: namaORS, Limit: limitORS})
	dalamJendela := adaMulai && !mulai.Before(jendelaSFHREMulai) && !mulai.After(jendelaSFHREAkhir)
	var saring []models.TreatySpreading
	for _, t := range semua {
		if strings.Contains(t.Name, "SF-HRE") && !dalamJendela {
			continue
		}
		saring = append(saring, t)
	}
	batas = unikTerakhir(saring)
	for _, t := range batas {
		if !strings.Contains(t.Name, "TRT") {
			dropdown = append(dropdown, models.TreatySpreading{ID: t.ID, Name: t.Name})
		}
	}
	dropdown = append(unikTerakhir(dropdown), models.TreatySpreading{ID: treatyFACOUT, Name: namaFACOUT})
	return dropdown, batas
}

// unikTerakhir - nama ganda: kemunculan terakhir yang tinggal, urutan asal dipertahankan.
func unikTerakhir(t []models.TreatySpreading) []models.TreatySpreading {
	sudah := map[string]bool{}
	var balik []models.TreatySpreading
	for i := len(t) - 1; i >= 0; i-- {
		if sudah[t[i].Name] {
			continue
		}
		sudah[t[i].Name] = true
		balik = append(balik, t[i])
	}
	out := make([]models.TreatySpreading, 0, len(balik))
	for i := len(balik) - 1; i >= 0; i-- {
		out = append(out, balik[i])
	}
	return out
}

// namaTreaty - nama dropdown ber-ID itu ("" bila tidak ada).
func namaTreaty(dropdown []models.TreatySpreading, id string) string {
	for _, t := range dropdown {
		if t.ID == id {
			return t.Name
		}
	}
	return ""
}

// ---- SumTSIPremiSpreadedRNM_Act langkah 18-22 -------------------------------------------------------------------

// sumberKurs - master kurs yang dibaca SumTSIPremiSpreadedRNM_Act.
type sumberKurs interface {
	KursPada(ctx context.Context, idMataUang, waktu string) (*apd.Decimal, error)
	IDMataUang(ctx context.Context, nama string) (string, error)
}

// ringkasan - TotalSpreadingCurrency (+ nama treaty), batas treaty (pesan verbatim), TotalSpreadAll.
func ringkasan(ctx context.Context, src sumberKurs, k models.KasusSpreading, perMataUang []models.TotalSpreading,
	dropdown, batas []models.TreatySpreading) ([]models.TotalSpreading, []models.RingkasanMataUangSpreading, []string, error) {
	cekTop := false // langkah 18.1.1 "TopRisk"
	for _, t := range perMataUang {
		cekTop = cekTop || lebihBesar(t.ClaimSpreaded, nol())
	}
	var pesan []string
	kurs := map[string]*apd.Decimal{}
	totalSpread := map[string]*apd.Decimal{} // treatyType -> Σ TotalSpread (IDR)
	treaty := make([]models.TotalSpreading, 0, len(perMataUang))
	namaSebelum := ""
	for _, t := range perMataUang {
		k1, ada := kurs[t.Currency]
		if !ada {
			id, err := src.IDMataUang(ctx, t.Currency)
			if err != nil {
				return nil, nil, nil, err
			}
			switch {
			case id == idMataUangIDR:
				k1 = apd.New(1, 0)
			case id != "":
				if k1, err = src.KursPada(ctx, id, k.StartDateTime); err != nil {
					return nil, nil, nil, err
				}
			}
			if k1 == nil {
				pesan = append(pesan, "Kurs "+t.Currency+" pada Begin date tidak ada di TREATYEXCHANGEYEARLY - batas treaty mata uang ini tidak diperiksa")
			}
			kurs[t.Currency] = k1
		}
		dasar := t.TSISpreaded
		if cekTop {
			dasar = t.ClaimSpreaded
		}
		if k1 != nil { // 20.5.1.1: TsiObj = @divide(TsiObj + TotalSpread, 1, 4) tiap baris
			totalSpread[t.TreatyType] = bagiSkala(tambah(totalSpread[t.TreatyType], kali(nolBila(dasar), k1)), apd.New(1, 0), 4)
		}
		// 19.5.5: Local.TreatyName hanya diganti bila jenis treaty ada di SpreadingList1 - bila tidak, nama baris sebelumnya
		// terbawa (verbatim).
		if n := namaTreaty(dropdown, t.TreatyType); n != "" {
			namaSebelum = n
		}
		baris := t
		baris.TreatyName = namaSebelum
		treaty = append(treaty, baris)
	}
	mulai, adaMulai := waktuMulai(k)
	for _, b := range batas {
		limit, err := utils.ParseDecimal(strings.TrimSpace(b.Limit))
		if err != nil {
			limit = nol()
		}
		if b.ID == treatyORS {
			limit = batasORS150
			if adaMulai && !mulai.Before(mulaiORS1815) {
				limit = batasORS1815
			}
		}
		if limit.Sign() > 0 && lebihBesar(totalSpread[b.ID], limit) {
			pesan = append(pesan, "Total TSI Spreaded for "+b.Name+" can't be more than "+b.Limit)
		}
	}
	semua := []models.RingkasanMataUangSpreading{}
	for _, t := range treaty {
		lol := t.ClaimEstimation
		if lebihBesar(t.ClaimAmountIDR, nol()) {
			lol = t.ClaimAmountIDR
		}
		premi := bagiSkala(t.PremiumSpreaded, apd.New(1, 0), 4) // 22.1 (dijumlah); baris baru 22.3 memakai nilai utuh
		n := -1
		for i := range semua {
			if semua[i].Currency == t.Currency {
				n = i
			}
		}
		if n < 0 {
			semua = append(semua, models.RingkasanMataUangSpreading{Currency: t.Currency, TSITopRisk: nolBila(t.ClaimSpreaded),
				TSI: nolBila(t.TSISpreaded), LoL: nolBila(lol), Premium: nolBila(t.PremiumSpreaded)})
			continue
		}
		s := &semua[n]
		s.TSI, s.Premium = tambah(s.TSI, t.TSISpreaded), tambah(s.Premium, premi)
		s.TSITopRisk, s.LoL = tambah(s.TSITopRisk, t.ClaimSpreaded), tambah(s.LoL, lol)
	}
	return treaty, semua, pesan, nil
}

// periksaJenisTreaty - cekSpreadingFactIn langkah 8-9: jenis treaty spreading yang bukan GetTreatyName dan bukan ORS /
// FACOUT -> "Invalid Treaty Type Spreading!" (sekali).
func periksaJenisTreaty(perMataUang []models.TotalSpreading, batas []models.TreatySpreading) []string {
	kenal := map[string]bool{treatyORS: true, treatyFACOUT: true}
	for _, b := range batas {
		kenal[b.ID] = true
	}
	for _, t := range perMataUang {
		if !kenal[t.TreatyType] {
			return []string{pesanJenisTreaty}
		}
	}
	return nil
}

// ---- alur -------------------------------------------------------------------------------------------------------

// bacaSpreading - BacaSpreading + pemetaan galat.
func (s *Service) bacaSpreading(ctx context.Context, tx *db.Tx, id string) (models.KasusSpreading, error) {
	if s.spreading == nil {
		return models.KasusSpreading{}, ErrSpreadingTanpaDatabase
	}
	if !idKasusSah(id) {
		return models.KasusSpreading{}, ErrKasusTidakAda
	}
	k, err := s.spreading.BacaSpreading(ctx, tx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return k, ErrKasusTidakAda
	}
	return k, err
}

// treatyKasus - dropdown dan batas treaty case (bizcode = BUSINESS.ID Class of Business case, cara Table of Limit);
// tanpa bizcode / Begin date -> hanya ORS + FACOUT (GetTreatyName langkah 7 dilewati).
func (s *Service) treatyKasus(ctx context.Context, id string, k models.KasusSpreading) ([]models.TreatySpreading, []models.TreatySpreading, error) {
	var master []models.TreatySpreading
	tanggal := tanggalMulai(k)
	bizcode, err := s.kodeBisnisKasus(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if bizcode != "" && tanggal != "" {
		if master, err = s.spreading.DaftarTreaty(ctx, bizcode, tanggal); err != nil {
			return nil, nil, err
		}
	}
	mulai, ada := waktuMulai(k)
	dropdown, batas := daftarTreaty(master, mulai, ada)
	return dropdown, batas, nil
}

// kodeBisnisKasus - BUSINESS.ID tunggal Class of Business + Group Business case; "" bila tidak dapat ditentukan (Class of
// Business / Group Business kosong, BUSINESS nol / ganda). Galat basis data diteruskan.
func (s *Service) kodeBisnisKasus(ctx context.Context, id string) (string, error) {
	if s.kasus == nil || s.tableOfLimit == nil {
		return "", nil
	}
	k, err := s.kasus.BacaKasus(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return "", ErrKasusTidakAda
	}
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(k.Opportunity.ClassOfBusiness) == "" || strings.TrimSpace(k.Opportunity.GroupBusinessID) == "" {
		return "", nil
	}
	kode, err := s.tableOfLimit.KodeBisnis(ctx, k.Opportunity.ClassOfBusiness, k.Opportunity.GroupBusinessID)
	if err != nil || len(kode) != 1 {
		return "", err
	}
	return kode[0], nil
}

// tampilan - total, ringkasan, pesan; periksaJenis = cekSpreadingFactIn langkah 8 (hanya Copy To All - Save Pega tidak
// memanggilnya). TotalPremiumNusantaraRe item = Σ PremiNusantaraRe coverage (FireMBU 6.2.2.1.4; tidak disimpan).
func (s *Service) tampilan(ctx context.Context, id string, k models.KasusSpreading, pesan []string, periksaJenis bool) (TampilanSpreading, error) {
	dropdown, batas, err := s.treatyKasus(ctx, id, k)
	if err != nil {
		return TampilanSpreading{}, err
	}
	setiapCoverage(&k, func(it *models.ItemSpreading, _ *models.CoverageSpreading) { it.TotalPremiumNusantaraRe = nil })
	setiapCoverage(&k, func(it *models.ItemSpreading, c *models.CoverageSpreading) {
		if c.PremiNusantaraRe != nil {
			it.TotalPremiumNusantaraRe = tambah(it.TotalPremiumNusantaraRe, c.PremiNusantaraRe)
		}
	})
	adaTop := totalLokasi(&k)
	perMataUang := totalMataUang(k, adaTop)
	for i := range k.Lokasi {
		for j := range k.Lokasi[i].Total {
			t := &k.Lokasi[i].Total[j]
			t.TreatyName = namaTreaty(dropdown, t.TreatyType)
		}
	}
	treaty, semua, pesanBatas, err := ringkasan(ctx, s.spreading, k, perMataUang, dropdown, batas)
	if err != nil {
		return TampilanSpreading{}, err
	}
	pesan = append(append([]string{}, pesan...), pesanBatas...)
	if periksaJenis {
		pesan = append(pesan, periksaJenisTreaty(perMataUang, batas)...)
	}
	if pesan == nil {
		pesan = []string{}
	}
	bulatkanSimpan(&k)
	for i := range k.Lokasi {
		k.Lokasi[i].Total = bulatkanTotal(k.Lokasi[i].Total)
		for j := range k.Lokasi[i].Items {
			it := &k.Lokasi[i].Items[j]
			it.TotalPremiumNusantaraRe = simpan8(it.TotalPremiumNusantaraRe)
		}
	}
	for i := range semua {
		r := &semua[i]
		r.TSITopRisk, r.TSI, r.LoL, r.Premium = simpan8(r.TSITopRisk), simpan8(r.TSI), simpan8(r.LoL), simpan8(r.Premium)
	}
	return TampilanSpreading{PercentShare: k.PercentShare, Treaty: dropdown, Template: []models.TemplateSpreading{}, Lokasi: k.Lokasi,
		RingkasanTreaty: bulatkanTotal(treaty), RingkasanMataUang: semua, Pesan: pesan}, nil
}

// bulatkanTotal - total tampil 8 desimal (nilai hitungan skala 20).
func bulatkanTotal(t []models.TotalSpreading) []models.TotalSpreading {
	for i := range t {
		x := &t[i]
		x.SharePercentage, x.ClaimSpreaded, x.TSISpreaded = simpan8(x.SharePercentage), simpan8(x.ClaimSpreaded), simpan8(x.TSISpreaded)
		x.ClaimEstimation, x.ClaimAmountIDR, x.PremiumSpreaded = simpan8(x.ClaimEstimation), simpan8(x.ClaimAmountIDR), simpan8(x.PremiumSpreaded)
	}
	return t
}

// periksaShareKasus - % Share RNM wajib, 0..100.
func periksaShareKasus(ps *apd.Decimal) error {
	if ps == nil {
		return fmt.Errorf("%w: percentShare wajib diisi", ErrMasukanSpreading)
	}
	if lebihBesar(ps, seratus) {
		return fmt.Errorf("%w: percentShare paling besar 100", ErrMasukanSpreading)
	}
	return nil
}

// periksaBarisShare - satu daftar baris (template / spreading coverage): share wajib, Σ <= 100, lebar teks.
func periksaBarisShare(jalur string, jenis, nama []string, share []*apd.Decimal) []string {
	var masalah []string
	if len(share) > banyakBarisSpreading {
		return []string{fmt.Sprintf("%s paling banyak %d baris", jalur, banyakBarisSpreading)}
	}
	total := nol()
	for i := range share {
		p := fmt.Sprintf("%s[%d].", jalur, i)
		switch {
		case share[i] == nil:
			masalah = append(masalah, p+"sharePercentage wajib diisi")
		case lebihBesar(share[i], seratus):
			masalah = append(masalah, p+"sharePercentage paling besar 100")
		}
		if strings.TrimSpace(jenis[i]) == "" {
			masalah = append(masalah, p+"treatyType wajib diisi")
		}
		if len(jenis[i]) > lebarTreatyType || len(nama[i]) > lebarTreatyName {
			masalah = append(masalah, fmt.Sprintf("%streatyType / treatyName paling banyak %d / %d bita", p, lebarTreatyType, lebarTreatyName))
		}
		total = tambah(total, share[i])
	}
	if lebihBesar(total, seratus) {
		masalah = append(masalah, jalur+": "+pesanShareLebih100)
	}
	return masalah
}

// TampilanSpreading - GET /api/nbfacin/kasus/{caseId}/spreading.
func (s *Service) TampilanSpreading(ctx context.Context, id string) (TampilanSpreading, error) {
	k, err := s.bacaSpreading(ctx, nil, id)
	if err != nil {
		return TampilanSpreading{}, err
	}
	return s.tampilan(ctx, id, k, nil, false)
}

// HitungShareSpreading - POST …/spreading/hitung-share: NR + spreading otomatis + total; TIDAK menyimpan. NR dihitung
// SEBELUM spreading otomatis (K48-3; Pega memakai PremiNusantaraRe lama di GetKapasitasTreaty).
func (s *Service) HitungShareSpreading(ctx context.Context, id string, ps *apd.Decimal) (TampilanSpreading, error) {
	if err := periksaShareKasus(ps); err != nil {
		return TampilanSpreading{}, err
	}
	k, err := s.bacaSpreading(ctx, nil, id)
	if err != nil {
		return TampilanSpreading{}, err
	}
	k.PercentShare = ps
	hitungNR(&k)
	pesan, err := spreadingOtomatis(ctx, s.spreading, &k)
	if err != nil {
		return TampilanSpreading{}, err
	}
	return s.tampilan(ctx, id, k, pesan, false)
}

// SalinSpreading - POST …/spreading/salin: Copy To All Spreading + total; TIDAK menyimpan.
func (s *Service) SalinSpreading(ctx context.Context, id string, ps *apd.Decimal, templat []models.TemplateSpreading) (TampilanSpreading, error) {
	if err := periksaShareKasus(ps); err != nil {
		return TampilanSpreading{}, err
	}
	if len(templat) == 0 {
		return TampilanSpreading{}, fmt.Errorf("%w: template wajib berisi paling sedikit satu baris", ErrMasukanSpreading)
	}
	jenis, nama, share := make([]string, len(templat)), make([]string, len(templat)), make([]*apd.Decimal, len(templat))
	for i, t := range templat {
		jenis[i], nama[i], share[i] = t.TreatyType, t.TreatyName, t.SharePercentage
	}
	if m := periksaBarisShare("template", jenis, nama, share); len(m) > 0 {
		return TampilanSpreading{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(m, "; "))
	}
	k, err := s.bacaSpreading(ctx, nil, id)
	if err != nil {
		return TampilanSpreading{}, err
	}
	m, err := s.proteksiTemplate(ctx, templat)
	if err != nil {
		return TampilanSpreading{}, err
	}
	if len(m) > 0 {
		return TampilanSpreading{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(m, "; "))
	}
	k.PercentShare = ps
	hitungNR(&k)
	salinTemplate(&k, templat)
	t, err := s.tampilan(ctx, id, k, nil, true)
	if err != nil {
		return t, err
	}
	t.Template = templat
	return t, nil
}

// proteksiTemplate - CalcultePersentageSpeading_Act (change Type Treaty grid Copy Spreading), langkah 2.13-2.17:
//   - jenis treaty sama di dua baris atau lebih -> "Treaty Type can't be same" (Local.Counter >= 2);
//   - baris pertama yang bukan ORS (10007) / FACOUT (10015) wajib ber-REINSURANCETYPE.NOTE memuat "QS" (@contains, peka
//     huruf) -> "Can not proceed spreading without QS".
//
// Langkah 2.3-2.9 "jenis treaty ini tidak terdapat untuk bisnis ini " (M_TREATYBUSINESS JSON + PROPORTIONALARRG, tahun
// berjalan) belum di-port (K48-9).
func (s *Service) proteksiTemplate(ctx context.Context, templat []models.TemplateSpreading) ([]string, error) {
	var m []string
	jumlah := map[string]int{}
	for _, t := range templat {
		jumlah[strings.TrimSpace(t.TreatyType)]++
	}
	for _, t := range templat {
		if jumlah[strings.TrimSpace(t.TreatyType)] >= 2 {
			m = append(m, pesanJenisKembar)
			break
		}
	}
	if pertama := strings.TrimSpace(templat[0].TreatyType); pertama != treatyORS && pertama != treatyFACOUT {
		note, err := s.spreading.CatatanJenisTreaty(ctx, pertama)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(note, "QS") {
			m = append(m, pesanTanpaQS)
		}
	}
	return m, nil
}

// SimpanSpreading - PUT …/spreading (Save tab, IsThereAnyObjectLocation_Act): bentuk lokasi wajib sama dengan objek
// tersimpan (409); hanya % Share RNM dan treatyType / treatyName / sharePercentage baris yang diambil; NR + DT percent +
// total dihitung ulang; simpan; jawab baca ulang.
func (s *Service) SimpanSpreading(ctx context.Context, pelaku inti.Pelaku, id string, ps *apd.Decimal, lokasi []models.LokasiSpreading) (TampilanSpreading, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanSpreading{}, err
	}
	if err := periksaShareKasus(ps); err != nil {
		return TampilanSpreading{}, err
	}
	k, err := s.bacaSpreading(ctx, nil, id)
	if err != nil {
		return TampilanSpreading{}, err
	}
	if len(k.Lokasi) == 0 { // IsThereAnyObjectLocation_Act 7 / 16: tanpa lokasi tidak disimpan
		return TampilanSpreading{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, pesanObjekKosong)
	}
	if !samaBentuk(k.Lokasi, lokasi) {
		return TampilanSpreading{}, ErrSpreadingBerubah
	}
	var masalah []string
	for i := range k.Lokasi {
		for j := range k.Lokasi[i].Items {
			for m := range k.Lokasi[i].Items[j].Coverages {
				baris := lokasi[i].Items[j].Coverages[m].Spreading
				jenis, nama, share := make([]string, len(baris)), make([]string, len(baris)), make([]*apd.Decimal, len(baris))
				salin := make([]models.BarisSpreading, len(baris))
				for n, b := range baris {
					jenis[n], nama[n], share[n] = b.TreatyType, b.TreatyName, b.SharePercentage
					salin[n] = models.BarisSpreading{TreatyType: strings.TrimSpace(b.TreatyType), TreatyName: strings.TrimSpace(b.TreatyName),
						SharePercentage: b.SharePercentage}
				}
				masalah = append(masalah, periksaBarisShare(fmt.Sprintf("lokasi[%d].items[%d].coverages[%d].spreading", i, j, m), jenis, nama, share)...)
				k.Lokasi[i].Items[j].Coverages[m].Spreading = salin
			}
		}
	}
	if len(masalah) > 0 {
		return TampilanSpreading{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(masalah, "; "))
	}
	k.PercentShare = ps
	hitungNR(&k)
	terapkanPersen(&k)
	totalLokasi(&k) // ClaimEstimation baris
	bulatkanSimpan(&k)
	if s.transaksi == nil {
		return TampilanSpreading{}, ErrSpreadingTanpaDatabase
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error { return s.spreading.TulisSpreading(ctx, tx, id, k, true) })
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return TampilanSpreading{}, ErrKasusTidakAda
	}
	if err != nil {
		return TampilanSpreading{}, err
	}
	return s.TampilanSpreading(ctx, id)
}

// samaBentuk - jumlah lokasi / item / coverage dan oldId coverage sama menurut posisi.
func samaBentuk(simpan, badan []models.LokasiSpreading) bool {
	if len(simpan) != len(badan) {
		return false
	}
	for i := range simpan {
		if len(simpan[i].Items) != len(badan[i].Items) {
			return false
		}
		for j := range simpan[i].Items {
			a, b := simpan[i].Items[j].Coverages, badan[i].Items[j].Coverages
			if len(a) != len(b) {
				return false
			}
			for m := range a {
				if a[m].OldID != b[m].OldID {
					return false
				}
			}
		}
	}
	return true
}

// bulatkanSimpan - nilai yang ditulis: 8 desimal (NUMBER(38,8)).
func bulatkanSimpan(k *models.KasusSpreading) {
	setiapCoverage(k, func(_ *models.ItemSpreading, c *models.CoverageSpreading) {
		c.TSINusantaraRe, c.PremiNusantaraRe = simpan8(c.TSINusantaraRe), simpan8(c.PremiNusantaraRe)
		for s := range c.Spreading {
			b := &c.Spreading[s]
			b.SharePercentage, b.TSISpreaded, b.TSIGrossSpreaded = simpan8(b.SharePercentage), simpan8(b.TSISpreaded), simpan8(b.TSIGrossSpreaded)
			b.PremiumSpreaded, b.ClaimEstimation = simpan8(b.PremiumSpreaded), simpan8(b.ClaimEstimation)
		}
	})
}

// adaIsiSpreading - % Share RNM terisi atau ada baris spreading.
func adaIsiSpreading(k models.KasusSpreading) bool {
	ada := k.PercentShare != nil
	setiapCoverage(&k, func(_ *models.ItemSpreading, c *models.CoverageSpreading) { ada = ada || len(c.Spreading) > 0 })
	return ada
}

// pertahankanSpreading - Save tab Object (K48-7): spreading coverage di posisi (lokasi, item, coverage) yang oldId-nya sama
// dipasang ke coverage baru, lalu NR + DT percent dihitung ulang dari % Share RNM tersimpan (Pega menghitung ulang di setiap
// Save, IsThereAnyObjectLocation_Act). Dipanggil di dalam transaksi Save Object, sesudah objek ditulis ulang.
func (s *Service) pertahankanSpreading(ctx context.Context, tx *db.Tx, id string, lama models.KasusSpreading) error {
	if s.spreading == nil || !adaIsiSpreading(lama) {
		return nil
	}
	baru, err := s.spreading.BacaSpreading(ctx, tx, id)
	if err != nil {
		return err
	}
	baru.PercentShare = lama.PercentShare
	for i := range baru.Lokasi {
		for j := range baru.Lokasi[i].Items {
			for m := range baru.Lokasi[i].Items[j].Coverages {
				c := &baru.Lokasi[i].Items[j].Coverages[m]
				c.Spreading = []models.BarisSpreading{}
				if i < len(lama.Lokasi) && j < len(lama.Lokasi[i].Items) && m < len(lama.Lokasi[i].Items[j].Coverages) &&
					lama.Lokasi[i].Items[j].Coverages[m].OldID == c.OldID {
					c.Spreading = salinBaris(lama.Lokasi[i].Items[j].Coverages[m].Spreading)
				}
			}
		}
	}
	if baru.PercentShare != nil {
		hitungNR(&baru)
		terapkanPersen(&baru)
		totalLokasi(&baru)
	}
	bulatkanSimpan(&baru)
	return s.spreading.TulisSpreading(ctx, tx, id, baru, false)
}
