package services

// Tab Premium Deduction / Brokerage Fee kasus FIRE (tiket 50). Port `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\Activity\`
// (dibaca 05-10-2026), per `.OfferFacIn.CurrencyList(n)` dengan Premium = `.Policy.Payment.Premium` = Σ PremiNusantaraRe
// item bermata uang itu (`CountPremiAndTSIRNMFireMBU_ACT` 6.2.5):
//   - SetRIComm_Act (aksi "persen", PUT): langkah 5 `PPnCheck = true` SELALU; 8.1 Commision = @Math.divide(RIComm ×
//     Premium, 100, 20); 8.3 BrokerageFee = @Math.divide(PctBrokerageFee × Premium, 100, 20); 8.4 / 8.5 dasar pajak =
//     BrokerageFee bila PPnCheck true, = @divide(BrokerageFee, @divide(PPN + 100, 100, 8), 8) bila false; 8.6 PPh / PPN =
//     @Math.divide(tarif × dasar, 100, 20). Langkah 8.2 / 9 berlabel `//` (mati); EDM (langkah 1, 6) di luar cakupan.
//   - SetTotalPremium_Act (aksi "commision", idx): RICommision = @Math.divide(Commision[idx], Premium[idx]) × 100, lalu
//     semua baris Commision = @Math.divide(RICommision × Premium, 100, 20).
//   - ProtekDeduction2_Act (aksi "deduction2", idx): Deduction2[idx] > Premium[idx] -> "Deduction 2 cannot be greater than
//     the Premium" (berhenti); selain itu PctDeduction2 = @Math.divide(Deduction2[idx], Premium[idx]) × 100, semua baris
//     Deduction2 = @Math.divide(PctDeduction2 × Premium, 100, 20).
//   - TotalBrokerageFee_Act (aksi "brokerage", idx): PctBrokerageFee = @Math.divide(BrokerageFee[idx], Premium[idx], 20) ×
//     100; semua baris BrokerageFee = @Math.divide(Pct × Premium, 100, 20), PPh / PPN dihitung ulang (8.3-8.5).
//   - Tarif PPH / PPN = POOLDATA.TAX baris TANGGAL terbaru (Obj-Browse Sort D, MaxRecords 1).
// `@Math.divide(a, b)` tanpa skala ditulis skala 20 seperti saudaranya (K50-4). SaveFillPaymentInstallment = tab Payment.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrMasukanDeduction - isian / aksi tab Premium Deduction tidak sah. 400.
	ErrMasukanDeduction = errors.New("services: isian tab Premium Deduction tidak sah")
	// ErrDeductionBerubah - mata uang badan tidak sama dengan mata uang objek tersimpan. 409.
	ErrDeductionBerubah = errors.New("Mata uang objek berubah sejak tab Premium Deduction dimuat - muat ulang tab")
	// ErrDeductionTanpaDatabase - tabel tab Premium Deduction tidak terbaca. 503.
	ErrDeductionTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tab Premium Deduction tidak terbaca")
)

const (
	pesanDeduction2Lebih = "Deduction 2 cannot be greater than the Premium" // ProtekDeduction2_Act langkah 2
	// AksiPersen / AksiCommision / AksiDeduction2 / AksiBrokerage - pemicu hitung (kontrak `AksiDeduction`).
	AksiPersen     = "persen"
	AksiCommision  = "commision"
	AksiDeduction2 = "deduction2"
	AksiBrokerage  = "brokerage"
	banyakMataUang = 50
)

// DenganDeduction memasang penyimpan tab Premium Deduction (tiket 50).
func (s *Service) DenganDeduction(p repository.PenyimpanDeduction) *Service {
	s.deduction = p
	return s
}

// TampilanDeduction - jawaban GET / POST hitung / PUT tab Premium Deduction.
type TampilanDeduction struct {
	PctPremiumDeduction, PctBrokerageFee *apd.Decimal
	PPnCheck                             bool
	TarifPPH, TarifPPN                   *apd.Decimal // nil bila TAX kosong
	MataUang                             []models.BarisDeduction
	Pesan                                []string
}

// IsianDeduction - badan POST hitung / PUT: persen atas dan nilai grid per mata uang (Premium, PPh, PPN diabaikan).
type IsianDeduction struct {
	PctPremiumDeduction, PctBrokerageFee *apd.Decimal
	MataUang                             []models.BarisDeduction
}

// ---- rumus ------------------------------------------------------------------------------------------------------

// premiPerMataUang - satu baris per mata uang objek (urut kemunculan pertama, ringkasanRNM), Premium = Σ premi NR; item
// tanpa mata uang tidak punya CurrencyList. Kosong bila % Share RNM kosong.
func premiPerMataUang(sp models.KasusSpreading) []models.BarisDeduction {
	_, jumlah, _ := ringkasanRNM(sp)
	baris := []models.BarisDeduction{}
	for _, j := range jumlah {
		if j.Currency != mataUangKosong {
			baris = append(baris, models.BarisDeduction{Currency: j.Currency, Premium: j.Premi})
		}
	}
	return baris
}

// hitungPajak - SetRIComm_Act 8.4-8.6 / TotalBrokerageFee_Act 8.3-8.5 atas satu baris.
func hitungPajak(b *models.BarisDeduction, tarif models.TarifPajak, ppnCheck bool) {
	dasar := nolBila(b.BrokerageFee)
	if !ppnCheck {
		inc := bagiSkala(tambah(tarif.PPN, seratus), seratus, 8)
		dasar = bagiSkala(dasar, inc, 8)
	}
	b.PPh, b.PPN = persenDari(tarif.PPH, dasar), persenDari(tarif.PPN, dasar)
}

// hitungRIComm - SetRIComm_Act langkah 8: Commision, BrokerageFee, PPh, PPN tiap baris dari persen; Deduction2 tidak
// disentuh.
func hitungRIComm(baris []models.BarisDeduction, ri, pctBrok *apd.Decimal, tarif models.TarifPajak, ppnCheck bool) {
	for i := range baris {
		b := &baris[i]
		b.Commision = persenDari(ri, b.Premium)
		b.BrokerageFee = persenDari(pctBrok, b.Premium)
		hitungPajak(b, tarif, ppnCheck)
	}
}

// persenBaris - @Math.divide(nilai, Premium, 20) × 100 (K50-4); Premium 0 -> galat (Pega tanpa penjaga, K50-5).
func persenBaris(nilai *apd.Decimal, b models.BarisDeduction, medan string) (*apd.Decimal, error) {
	if nolBila(b.Premium).IsZero() {
		return nil, fmt.Errorf("%w: Premium %s 0 - %s tidak dapat dihitung menjadi persen", ErrMasukanDeduction, b.Currency, medan)
	}
	return kali(bagiSkala(nilai, b.Premium, 20), seratus), nil
}

// segarkanCicilan - K50-12: Save tab Premium Deduction juga menghitung ulang cicilan tab Payment yang tersimpan (persen dan
// Due Date tetap; nilai uang, Total Payment, Net Premium dari nilai Payment baru) - Pega menjalankan
// SaveFillPaymentInstallment + CountPaymentInstallment_Act di setiap Save (IsThereAnyObjectLocation_Act 12-13), sehingga
// kolom tersimpan tidak basi bagi pembaca Oracle langsung. Tanpa Installment / cicilan, atau bila Net Premium negatif
// (K51-7), tidak ada yang ditulis; GET tetap menghitung ulang (K51-4). Proteksi Payment tidak dijalankan di sini.
func (s *Service) segarkanCicilan(ctx context.Context, tx *db.Tx, id string, bayar models.KasusPayment, baris []models.BarisDeduction) error {
	n, err := strconv.Atoi(strings.TrimSpace(bayar.Installment))
	if s.payment == nil || err != nil || len(bayar.Cicilan) == 0 {
		return nil
	}
	mu := make([]models.MataUangPayment, 0, len(baris))
	for _, b := range baris {
		m := models.MataUangPayment{Baris: b, Cicilan: append([]models.CicilanPayment{}, bayar.Cicilan[b.Currency]...)}
		isiCicilan(&m)
		if m.NetPremium.Sign() < 0 {
			return nil
		}
		mu = append(mu, m)
	}
	bulatkanPayment(mu)
	return s.payment.TulisPayment(ctx, tx, id, models.SimpanPayment{Installment: n, MataUang: mu})
}

// persenPenuh - persen berpresisi penuh (skala 20) yang menghasilkan nilai grid, bila persen 8 desimal yang dikirim layar
// adalah pembulatannya; selain itu persen kiriman apa adanya (K50-9). Pega menyimpan persen hasil aksi baris
// (SetTotalPremium_Act 1 / TotalBrokerageFee_Act 6) berpresisi penuh di clipboard dan Save menghitung ulang darinya, jadi
// nilai yang diketik tidak bergeser; di sini persen tampil / tersimpan 8 desimal (NUMBER(38,8)), dan presisi penuhnya
// dibangun ulang dari baris pertama berpremi bukan nol.
func persenPenuh(pct8 *apd.Decimal, baris []models.BarisDeduction, nilai func(models.BarisDeduction) *apd.Decimal) *apd.Decimal {
	if pct8 == nil {
		return nil
	}
	for _, b := range baris {
		if nolBila(b.Premium).IsZero() {
			continue
		}
		penuh := kali(bagiSkala(nilai(b), b.Premium, 20), seratus)
		if samaDesimal(simpan8(penuh), simpan8(pct8)) {
			return penuh
		}
		return pct8
	}
	return pct8
}

func nilaiCommision(b models.BarisDeduction) *apd.Decimal    { return b.Commision }
func nilaiBrokerageFee(b models.BarisDeduction) *apd.Decimal { return b.BrokerageFee }

// pesanDeduction2 - satu pesan verbatim bila ada baris Deduction2 > Premium.
func pesanDeduction2(baris []models.BarisDeduction) []string {
	for _, b := range baris {
		if lebihBesar(b.Deduction2, b.Premium) {
			return []string{pesanDeduction2Lebih}
		}
	}
	return nil
}

func bulatkanDeduction(baris []models.BarisDeduction) {
	for i := range baris {
		b := &baris[i]
		b.Premium, b.Commision, b.Deduction2 = simpan8(b.Premium), simpan8(b.Commision), simpan8(b.Deduction2)
		b.BrokerageFee, b.PPh, b.PPN = simpan8(b.BrokerageFee), simpan8(b.PPh), simpan8(b.PPN)
	}
}

// ---- alur -------------------------------------------------------------------------------------------------------

// dasarDeduction - premi per mata uang (objek tersimpan), data tersimpan, tarif pajak, pesan dasar.
func (s *Service) dasarDeduction(ctx context.Context, id string) ([]models.BarisDeduction, models.KasusDeduction, models.TarifPajak, bool, []string, error) {
	var k models.KasusDeduction
	var tarif models.TarifPajak
	if s.deduction == nil || s.spreading == nil {
		return nil, k, tarif, false, nil, ErrDeductionTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, k, tarif, false, nil, ErrKasusTidakAda
	}
	k, err := s.deduction.BacaDeduction(ctx, id)
	if err != nil {
		return nil, k, tarif, false, nil, galatKasus(err)
	}
	sp, err := s.bacaSpreadingRingkas(ctx, nil, id)
	if err != nil {
		return nil, k, tarif, false, nil, galatKasus(err)
	}
	tarif, ada, err := s.deduction.TarifPajak(ctx)
	if err != nil {
		return nil, k, tarif, false, nil, err
	}
	var pesan []string
	if sp.PercentShare == nil {
		pesan = append(pesan, "% Share RNM belum diisi (tab Spreading) - premi Nusantara Re per mata uang belum ada")
	}
	if !ada {
		pesan = append(pesan, "Tarif PPH / PPN (tabel TAX) tidak ada - PPh / PPN dihitung 0")
	}
	return premiPerMataUang(sp), k, tarif, ada, pesan, nil
}

func tampilanDeduction(ri, pct *apd.Decimal, tarif models.TarifPajak, adaTarif bool, baris []models.BarisDeduction, pesan []string) TampilanDeduction {
	bulatkanDeduction(baris)
	t := TampilanDeduction{PctPremiumDeduction: simpan8(ri), PctBrokerageFee: simpan8(pct), PPnCheck: true, MataUang: baris,
		Pesan: append(append([]string{}, pesan...), pesanDeduction2(baris)...)}
	if adaTarif {
		t.TarifPPH, t.TarifPPN = simpan8(tarif.PPH), simpan8(tarif.PPN)
	}
	return t
}

// TampilanDeduction - GET /api/nbfacin/kasus/{caseId}/deduction. Commision / BrokerageFee / PPh / PPN tersimpan dipakai
// bila premi mata uang itu sama dengan saat Save; selain itu dihitung ulang dari persen tersimpan (K50-3, Pega menghitung
// ulang di setiap Save tab lain lewat IsThereAnyObjectLocation_Act langkah 11). Deduction2 = tersimpan.
func (s *Service) TampilanDeduction(ctx context.Context, id string) (TampilanDeduction, error) {
	baris, k, tarif, ada, pesan, err := s.dasarDeduction(ctx, id)
	if err != nil {
		return TampilanDeduction{}, err
	}
	ppn := k.PPnCheck != "false"
	hitungRIComm(baris, k.RICommision, k.PctBrokerageFee, tarif, ppn)
	simpan := map[string]models.BarisDeduction{}
	for _, b := range k.MataUang {
		simpan[b.Currency] = b
	}
	for i := range baris {
		b := &baris[i]
		lama, adaLama := simpan[b.Currency]
		if !adaLama {
			continue
		}
		b.Deduction2 = lama.Deduction2
		if lama.Premium != nil && samaDesimal(lama.Premium, simpan8(b.Premium)) {
			b.Commision, b.BrokerageFee, b.PPh, b.PPN = lama.Commision, lama.BrokerageFee, lama.PPh, lama.PPN
		}
	}
	t := tampilanDeduction(k.RICommision, k.PctBrokerageFee, tarif, ada, baris, pesan)
	t.PPnCheck = k.PPnCheck == "true"
	return t, nil
}

// cocokkanIsian - nilai grid badan dipasang ke baris server menurut KODE mata uang (bukan urutan); himpunan mata uang
// wajib sama (409). Deduction2 / Commision / BrokerageFee kosong = 0.
func cocokkanIsian(baris []models.BarisDeduction, isian []models.BarisDeduction) error {
	if len(isian) != len(baris) {
		return ErrDeductionBerubah
	}
	badan := map[string]models.BarisDeduction{}
	for _, b := range isian {
		if _, ganda := badan[b.Currency]; ganda {
			return fmt.Errorf("%w: mata uang %s ganda", ErrMasukanDeduction, b.Currency)
		}
		badan[b.Currency] = b
	}
	for i := range baris {
		b, ada := badan[baris[i].Currency]
		if !ada {
			return ErrDeductionBerubah
		}
		baris[i].Commision, baris[i].Deduction2, baris[i].BrokerageFee = nolBila(b.Commision), nolBila(b.Deduction2), nolBila(b.BrokerageFee)
	}
	return nil
}

// indeksMataUang - posisi baris server bermata uang isian[idx].
func indeksMataUang(baris []models.BarisDeduction, isian []models.BarisDeduction, idx int) (int, error) {
	if idx < 0 || idx >= len(isian) {
		return 0, fmt.Errorf("%w: idx %d di luar 0..%d", ErrMasukanDeduction, idx, len(isian)-1)
	}
	for i := range baris {
		if baris[i].Currency == isian[idx].Currency {
			return i, nil
		}
	}
	return 0, ErrDeductionBerubah
}

func periksaIsianDeduction(isian IsianDeduction) error {
	if len(isian.MataUang) > banyakMataUang {
		return fmt.Errorf("%w: mataUang paling banyak %d baris", ErrMasukanDeduction, banyakMataUang)
	}
	return nil
}

// HitungDeduction - POST …/deduction/hitung: satu aksi Pega atas isian layar, TANPA simpan. PPnCheck = true (SetRIComm_Act
// langkah 5; TotalBrokerageFee membaca nilai yang sama).
func (s *Service) HitungDeduction(ctx context.Context, id, aksi string, idx int, isian IsianDeduction) (TampilanDeduction, error) {
	if err := periksaIsianDeduction(isian); err != nil {
		return TampilanDeduction{}, err
	}
	switch aksi {
	case AksiPersen, AksiCommision, AksiDeduction2, AksiBrokerage:
	default:
		return TampilanDeduction{}, fmt.Errorf("%w: aksi harus persen, commision, deduction2, atau brokerage", ErrMasukanDeduction)
	}
	baris, _, tarif, ada, pesan, err := s.dasarDeduction(ctx, id)
	if err != nil {
		return TampilanDeduction{}, err
	}
	if err := cocokkanIsian(baris, isian.MataUang); err != nil {
		return TampilanDeduction{}, err
	}
	ri := persenPenuh(isian.PctPremiumDeduction, baris, nilaiCommision)
	pct := persenPenuh(isian.PctBrokerageFee, baris, nilaiBrokerageFee)
	switch aksi {
	case AksiPersen:
		hitungRIComm(baris, ri, pct, tarif, true)
	case AksiCommision:
		i, err := indeksMataUang(baris, isian.MataUang, idx)
		if err != nil {
			return TampilanDeduction{}, err
		}
		if ri, err = persenBaris(baris[i].Commision, baris[i], "Total Premium Deduction"); err != nil {
			return TampilanDeduction{}, err
		}
		for j := range baris {
			baris[j].Commision = persenDari(ri, baris[j].Premium)
			hitungPajak(&baris[j], tarif, true)
		}
	case AksiDeduction2:
		i, err := indeksMataUang(baris, isian.MataUang, idx)
		if err != nil {
			return TampilanDeduction{}, err
		}
		if lebihBesar(baris[i].Deduction2, baris[i].Premium) {
			return TampilanDeduction{}, fmt.Errorf("%w: %s", ErrMasukanDeduction, pesanDeduction2Lebih)
		}
		pctD2, err := persenBaris(baris[i].Deduction2, baris[i], "Total Premium Deduction 2")
		if err != nil {
			return TampilanDeduction{}, err
		}
		for j := range baris {
			baris[j].Deduction2 = persenDari(pctD2, baris[j].Premium)
			hitungPajak(&baris[j], tarif, true)
		}
	case AksiBrokerage:
		i, err := indeksMataUang(baris, isian.MataUang, idx)
		if err != nil {
			return TampilanDeduction{}, err
		}
		if pct, err = persenBaris(baris[i].BrokerageFee, baris[i], "Total Brokerage Fee"); err != nil {
			return TampilanDeduction{}, err
		}
		for j := range baris {
			baris[j].BrokerageFee = persenDari(pct, baris[j].Premium)
			hitungPajak(&baris[j], tarif, true)
		}
	}
	return tampilanDeduction(ri, pct, tarif, ada, baris, pesan), nil
}

// SimpanDeduction - PUT …/deduction: SetRIComm_Act dari persen badan (PPnCheck = true), Deduction2 badan apa adanya
// (wajib <= Premium, pesan verbatim), simpan, jawab baca ulang.
func (s *Service) SimpanDeduction(ctx context.Context, pelaku inti.Pelaku, id string, isian IsianDeduction) (TampilanDeduction, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanDeduction{}, err
	}
	if err := s.wajibBolehUbah(ctx, pelaku, id); err != nil {
		return TampilanDeduction{}, err
	}
	if err := periksaIsianDeduction(isian); err != nil {
		return TampilanDeduction{}, err
	}
	if s.transaksi == nil {
		return TampilanDeduction{}, ErrDeductionTanpaDatabase
	}
	_, _, tarif, _, _, err := s.dasarDeduction(ctx, id)
	if err != nil {
		return TampilanDeduction{}, err
	}
	// Cicilan tab Payment yang tersimpan ikut dihitung ulang di transaksi yang sama (K50-12).
	var bayar models.KasusPayment
	if s.payment != nil {
		if bayar, err = s.payment.BacaPayment(ctx, id); err != nil {
			return TampilanDeduction{}, galatKasus(err)
		}
	}
	// Premi per mata uang dibaca di transaksi yang sama dengan penulisan (bukan dari baca sebelumnya), seperti tab Cedant.
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		sp, err := s.bacaSpreadingRingkas(ctx, tx, id)
		if err != nil {
			return err
		}
		baris := premiPerMataUang(sp)
		if err := cocokkanIsian(baris, isian.MataUang); err != nil {
			return err
		}
		var lebih []string
		for _, b := range baris {
			if lebihBesar(b.Deduction2, b.Premium) {
				lebih = append(lebih, b.Currency)
			}
		}
		if len(lebih) > 0 {
			return fmt.Errorf("%w: %s (%s)", ErrMasukanDeduction, pesanDeduction2Lebih, strings.Join(lebih, ", "))
		}
		ri := persenPenuh(isian.PctPremiumDeduction, baris, nilaiCommision)
		pct := persenPenuh(isian.PctBrokerageFee, baris, nilaiBrokerageFee)
		hitungRIComm(baris, ri, pct, tarif, true)
		bulatkanDeduction(baris)
		if err := s.deduction.TulisDeduction(ctx, tx, id, models.KasusDeduction{RICommision: simpan8(isian.PctPremiumDeduction),
			PctBrokerageFee: simpan8(isian.PctBrokerageFee), PPnCheck: "true", MataUang: baris}); err != nil {
			return err
		}
		return s.segarkanCicilan(ctx, tx, id, bayar, baris)
	})
	if err != nil {
		return TampilanDeduction{}, galatKasus(err)
	}
	return s.TampilanDeduction(ctx, id)
}
