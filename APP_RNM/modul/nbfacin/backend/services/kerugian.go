package services

// Sub-tab Loss Record - tiket 42: pemeriksaan catatan kerugian, Insured Name, tanggal, dan hitung ulang Loss Ratio.
//
// Loss Ratio `[terverifikasi]` `D:\migrasi\RNM\DDL\SetLossRatio_Act.xml` (kelas OfferFacIn-LocationReinsurance; versi
// NB FacIn\Activity hanya kelas ScoringRisk), diteruskan sesi 0f + keputusan work owner W-4 (03-10-2026):
//   - langkah 2: keempat LR = 0;
//   - langkah 3 (loop .Property.ListCauseOfLoss; prasyarat CARI8 MATI): .CoinsData.CoinsName = QuotationData.InsuredName;
//     bila Date - DateOfLoss <= 365: ΣAmount1 += Amount, ΣClaim1 += Claim; bila <= 1825: ΣAmount2, ΣClaim2 (kumulatif);
//   - langkah 4/5, hanya bila ΣClaim != 0: LR Amount = ΣClaim/ΣAmount, LR Percent = (ΣClaim/ΣAmount) * 100;
//   - W-4: ΣAmount = 0 -> LR tetap 0 (Pega tidak menjaga pembagi nol).
// Date = @DateTime.CurrentDate (ViewOfferFacInUW_PreDT) -> tanggal hari ini Asia/Jakarta (WIB).
// Tidak diport: langkah 1 (IsB2B "ASM" -> Detail catatan terakhir "No Info").

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// skalaLossRatio - A148: hasil bagi LR dibulatkan setengah-ke-atas pada desimal ke-8 (ADR-0016: delapan desimal).
const skalaLossRatio = 8

// hariLR1, hariLR35 - ambang hari langkah 3 (365 = 1 tahun, 1825 = 5 tahun).
const hariLR1, hariLR35 time.Duration = 365, 1825

// batasKerugian - T_LISTCAUSEOFLOSS.SEQ_NO NUMBER(5).
const batasKerugian = 99999

// lebarKerugian - lebar kolom (BYTE) migrasi 191. coinsName tidak diperiksa: selalu ditimpa nama tertanggung case.
var lebarKerugian = []struct {
	nama  string
	nilai func(models.CatatanKerugian) string
	n     int
}{
	{"lossObject", func(c models.CatatanKerugian) string { return c.LossObject }, 500},
	{"currency", func(c models.CatatanKerugian) string { return c.Currency }, 50},
	{"causeOfLoss", func(c models.CatatanKerugian) string { return c.CauseOfLoss }, 500},
	{"remarks", func(c models.CatatanKerugian) string { return c.Remarks }, 500},
	{"detail", func(c models.CatatanKerugian) string { return c.Detail }, 500},
}

// uangKerugian - tiga medan uang catatan kerugian (NUMBER(38,8)).
var uangKerugian = []struct {
	nama  string
	nilai func(models.CatatanKerugian) string
}{
	{"amount", func(c models.CatatanKerugian) string { return c.Amount }},
	{"claim", func(c models.CatatanKerugian) string { return c.Claim }},
	{"preventionOfLoss", func(c models.CatatanKerugian) string { return c.PreventionOfLoss }},
}

// periksaKerugian - catatan kerugian baris objek ke-`n` (tanpa basis data). Currency wajib (K-069/K-012, pola A133);
// Loss Detail tidak dipaksa wajib (N-4); Remarks tanpa enumerasi.
func periksaKerugian(n int, rugi []models.CatatanKerugian) []string {
	if len(rugi) > batasKerugian {
		return []string{fmt.Sprintf("baris[%d].lossRecords paling banyak %d", n, batasKerugian)}
	}
	var masalah []string
	for m, c := range rugi {
		awal := fmt.Sprintf("baris[%d].lossRecords[%d].", n, m)
		if strings.TrimSpace(c.Currency) == "" {
			masalah = append(masalah, awal+"currency wajib diisi")
		}
		for _, u := range uangKerugian {
			if v := u.nilai(c); v != "" && !polaDesimal.MatchString(v) {
				masalah = append(masalah, awal+u.nama+" harus angka >= 0 dengan paling banyak 8 desimal")
			}
		}
		var tgl []string
		uraiKabel(awal+"dateOfLoss", c.DateOfLoss, &tgl)
		masalah = append(masalah, tgl...)
		for _, l := range lebarKerugian {
			if len(l.nilai(c)) > l.n {
				masalah = append(masalah, fmt.Sprintf("%s%s paling banyak %d byte", awal, l.nama, l.n))
			}
		}
	}
	return masalah
}

// desimalAtauNol - teks desimal (sudah lolos periksa) -> apd; kosong = 0.
func desimalAtauNol(s string) *apd.Decimal {
	if s == "" {
		return apd.New(0, 0)
	}
	d, err := utils.ParseDecimal(s)
	if err != nil {
		return apd.New(0, 0)
	}
	return d
}

// konteksLR - pembagian berskala tetap, setengah-ke-atas (A148). Tidak diubah sesudah dibuat (aman dipakai bersama).
var konteksLR = func() *apd.Context {
	c := apd.BaseContext.WithPrecision(60)
	c.Rounding = apd.RoundHalfUp
	return c
}()

// errLRTerlalu - hasil Loss Ratio tidak muat di NUMBER(38,8) (ΣClaim jauh melebihi ΣAmount); dipetakan 400.
var errLRTerlalu = errors.New("loss ratio melebihi NUMBER(38,8)")

// bagiLR - pembilang * kali / penyebut, dibulatkan ke skalaLossRatio desimal. Penyebut 0 -> "0" (W-4). Hasil wajib
// muat di NUMBER(38,8) (polaDesimal: <= 30 digit bulat) - selain itu errLRTerlalu, bukan galat Oracle saat simpan.
func bagiLR(pembilang, penyebut *apd.Decimal, kali int64) (string, error) {
	if penyebut.IsZero() {
		return "0", nil
	}
	var h apd.Decimal
	if _, err := konteksLR.Mul(&h, pembilang, apd.New(kali, 0)); err != nil {
		return "", fmt.Errorf("services: loss ratio: %w", err)
	}
	if _, err := konteksLR.Quo(&h, &h, penyebut); err != nil {
		return "", fmt.Errorf("services: loss ratio: %w", err)
	}
	if _, err := konteksLR.Quantize(&h, &h, -skalaLossRatio); err != nil {
		return "", fmt.Errorf("services: loss ratio: %w", err)
	}
	h.Reduce(&h)
	s := utils.FormatDecimal(&h)
	if !polaDesimal.MatchString(s) {
		return "", errLRTerlalu
	}
	return s, nil
}

// hitungLossRatio - langkah 2-5 SetLossRatio_Act atas catatan sebuah objek pada tanggal `hari` (WIB). Catatan tanpa
// dateOfLoss TIDAK dijumlah (A147: selisih tanggal tidak terdefinisi). Tanggal kerugian sesudah `hari` (selisih
// negatif) ikut, seperti perbandingan `<=` Pega.
func hitungLossRatio(rugi []models.CatatanKerugian, hari time.Time) (models.LossRatio, error) {
	nol := models.LossRatio{OneYearAmount: "0", OneYearPercent: "0", ThreeFiveYearAmount: "0", ThreeFiveYearPercent: "0"}
	acuan := time.Date(hari.Year(), hari.Month(), hari.Day(), 0, 0, 0, 0, WIB)
	amount1, claim1, amount2, claim2 := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for _, c := range rugi {
		t, err := time.ParseInLocation(BentukTanggalKabel, c.DateOfLoss, WIB)
		if c.DateOfLoss == "" || err != nil {
			continue
		}
		// Dua tengah malam WIB (tanpa musim panas): selisih tepat kelipatan 24 jam - pembagian bulat, tanpa float.
		selisih := acuan.Sub(t) / (24 * time.Hour)
		a, k := desimalAtauNol(c.Amount), desimalAtauNol(c.Claim)
		for _, j := range []struct {
			batas         time.Duration
			amount, claim *apd.Decimal
		}{{hariLR1, amount1, claim1}, {hariLR35, amount2, claim2}} {
			if selisih > j.batas {
				continue
			}
			if _, err := konteksLR.Add(j.amount, j.amount, a); err != nil {
				return nol, fmt.Errorf("services: loss ratio: %w", err)
			}
			if _, err := konteksLR.Add(j.claim, j.claim, k); err != nil {
				return nol, fmt.Errorf("services: loss ratio: %w", err)
			}
		}
	}
	for _, h := range []struct {
		amount, claim *apd.Decimal
		nilai, persen *string
	}{{amount1, claim1, &nol.OneYearAmount, &nol.OneYearPercent}, {amount2, claim2, &nol.ThreeFiveYearAmount, &nol.ThreeFiveYearPercent}} {
		if h.claim.IsZero() {
			continue
		}
		var err error
		if *h.nilai, err = bagiLR(h.claim, h.amount, 1); err != nil {
			return nol, err
		}
		if *h.persen, err = bagiLR(h.claim, h.amount, 100); err != nil {
			return nol, err
		}
	}
	return nol, nil
}

// tanggalKerugianKePega - DD-MM-YYYY (sudah lolos periksa) -> teks Pega DateTime pukul 12:00 WIB (pola Begin date,
// butir 78.1); kosong = "".
func tanggalKerugianKePega(s string) string {
	t, err := time.ParseInLocation(BentukTanggalKabel, s, WIB)
	if s == "" || err != nil {
		return ""
	}
	return t.Add(jamTulisWaktu).UTC().Format(BentukWaktuPega)
}

// tanggalKerugianKeKabel - teks Pega DateTime ATAU Date -> DD-MM-YYYY (WIB); teks lain apa adanya (A89).
func tanggalKerugianKeKabel(s string) string {
	if _, err := time.Parse(BentukWaktuPega, s); err == nil {
		return waktuKeKabel(s)
	}
	return tanggalKeKabel(s)
}

// siapkanKerugian - salinan daftar objek untuk disimpan (tiket 42): per objek Loss Ratio dihitung ulang (lossRatio
// badan PUT diabaikan), CoinsName = nama tertanggung case (langkah 3.1), dateOfLoss -> teks Pega. Nama tertanggung
// dibaca hanya bila ada catatan kerugian.
func (s *Service) siapkanKerugian(ctx context.Context, id string, baris []models.ObjekFire) ([]models.ObjekFire, error) {
	hari := s.sekarang().In(WIB)
	var tertanggung *string
	hasil := make([]models.ObjekFire, len(baris))
	for i, o := range baris {
		lr, err := hitungLossRatio(o.LossRecords, hari)
		if errors.Is(err, errLRTerlalu) {
			return nil, fmt.Errorf("%w: baris[%d].lossRatio: %v (ΣClaim / ΣAmount)", ErrMasukanObjek, i, err)
		}
		if err != nil {
			return nil, err
		}
		o.LossRatio = lr
		rugi := make([]models.CatatanKerugian, len(o.LossRecords))
		for j, c := range o.LossRecords {
			if tertanggung == nil {
				if s.kasus == nil {
					return nil, ErrObjekTanpaDatabase
				}
				k, err := s.kasus.BacaKasus(ctx, id)
				if errors.Is(err, repository.ErrKasusTidakAda) {
					return nil, ErrKasusTidakAda
				}
				if err != nil {
					return nil, err
				}
				tertanggung = &k.InsuredName
			}
			c.CoinsName = *tertanggung
			c.DateOfLoss = tanggalKerugianKePega(c.DateOfLoss)
			rugi[j] = c
		}
		o.LossRecords = rugi
		hasil[i] = o
	}
	return hasil, nil
}
