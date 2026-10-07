package services

// Tab Payment kasus FIRE (tiket 51), per `.OfferFacIn.CurrencyList(n)` (nilai Payment mata uang dari tab Premium Deduction,
// tiket 50). Port `[terverifikasi]` (dibaca ulang sesi c3 06-10-2026):
//   - aksi "installment" = `Endorsment Fac In\Activity\FillPaymentInstallment.xml` (ASM-FW-GISFW-WORK!FILLPAYMENTINSTALLMENT;
//     NB FacIn tidak punya kelas ini - dipinjam dengan persetujuan work owner 06-10-2026, K51-2): tiap mata uang dibangun
//     ulang n cicilan, % = @Math.divide(100, n, 20), cicilan terakhir = 100/n + (100 − jumlah berjalan); DueDate =
//     @DateTime.addCalendar(@FormatDateTime(StartDateTime, "yyyyMMdd", "Asia/Jakarta"), 0, 3, ...) - SAMA untuk semua
//     cicilan, "" bila StartDateTime kosong (keputusan work owner "iya ikut pega").
//   - aksi "persen" = `Endorsment Fac In\Activity\SetValidateInstallment_Act.xml` (…-CURRENCY!SETVALIDATEINSTALLMENT_ACT):
//     Σ % > 100 -> "Total % Installment Can't Be More Than 100%", < 100 -> "Total % Installment Must Be 100%" (pesan, tidak
//     menghalangi).
//   - PUT = NB `SaveFillPaymentInstallment` + `CountPaymentInstallment_Act` (IsThereAnyObjectLocation_Act 12-13): % cicilan
//     dipertahankan bila Σ = 100, selain itu diratakan ulang 100/n; per cicilan Premium / RICommision / BrokerageFee / PPh /
//     PPN / PaymentTotal = nilai mata uang × % / 100 (skala 20, CountPaymentInstallment menimpa RICommision / PaymentTotal
//     skala 4 SaveFill), Deduction2 / Discount skala 4; lalu `ProtectPremiPolicy_Act` (pesan verbatim, ditolak 400 -
//     keputusan work owner, K51-5).
//   - NetPremium = SumTotalPayment = Premium − Commision − Deduction2 − BrokerageFee + PPh + PPN (SaveFillPaymentInstallment
//     3.1.2 "tambah deduction 2"); MENYIMPANG sadar dari CountPaymentInstallment_Act (yang membuang Deduction2) - keputusan
//     work owner 06-10-2026 "ikut dikurangkan" (K51-1).

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
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrMasukanPayment - isian / aksi / proteksi Save tab Payment. 400.
	ErrMasukanPayment = errors.New("services: isian tab Payment tidak sah")
	// ErrPaymentBerubah - mata uang badan tidak sama dengan mata uang objek tersimpan. 409.
	ErrPaymentBerubah = errors.New("Mata uang objek berubah sejak tab Payment dimuat - muat ulang tab")
	// ErrPaymentTanpaDatabase - tabel tab Payment tidak terbaca. 503.
	ErrPaymentTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tab Payment tidak terbaca")
)

const (
	// Pesan verbatim Pega.
	pesanPersenLebih100  = "Total % Installment Can't Be More Than 100%"                     // SetValidateInstallment_Act 3
	pesanPersenKurang100 = "Total % Installment Must Be 100%"                                // SetValidateInstallment_Act 4
	pesanPaymentKosong   = "List Payment can't be null!"                                     // ProtectPremiPolicy_Act 2 / 4
	pesanCicilanKosong   = "List Installment can't be null!"                                 // ProtectPremiPolicy_Act 5.7
	pesanPremiNol        = "Premi on paymnet can't be 0"                                     // ProtectPremiPolicy_Act 11
	pesanTotalTakCocok   = "Total premium RNM on installment didn't match with premium RNM!" // ProtectPremiPolicy_Act 5.5 / 10
	// AksiInstallment / AksiPersenCicilan - pemicu hitung (kontrak `AksiPayment`).
	AksiInstallment   = "installment"
	AksiPersenCicilan = "persen"
	// banyakCicilanMaks - layar `pyMax=2` (dua digit).
	banyakCicilanMaks = 99
	// bulanJatuhTempo - FillPaymentInstallment @DateTime.addCalendar(…, 0, 3, …).
	bulanJatuhTempo = 3
	// lebarKodeMataUangCicilan - T_LISTINSTALLMENT.CURRENCY_CODE VARCHAR2(10) (rancangan).
	lebarKodeMataUangCicilan = 10
)

// DenganPayment memasang penyimpan tab Payment (tiket 51).
func (s *Service) DenganPayment(p repository.PenyimpanPayment) *Service {
	s.payment = p
	return s
}

// TampilanPayment - jawaban GET / POST hitung / PUT tab Payment.
type TampilanPayment struct {
	Installment        string
	TarifPPH, TarifPPN *apd.Decimal
	MataUang           []models.MataUangPayment
	Pesan              []string
}

// IsianCicilanMataUang - cicilan badan satu mata uang (No, DueDate `YYYYMMDD`, persen).
type IsianCicilanMataUang struct {
	Currency string
	Cicilan  []models.CicilanPayment
}

// IsianPayment - badan POST hitung / PUT.
type IsianPayment struct {
	Installment string
	MataUang    []IsianCicilanMataUang
}

// ---- rumus ------------------------------------------------------------------------------------------------------

var satu = apd.New(1, 0)

// netPremium - SaveFillPaymentInstallment 3.1.2: Premium − Commision − Deduction2 − BrokerageFee + PPh + PPN, lalu
// @Math.divide(x, 1, 20).
func netPremium(b models.BarisDeduction) *apd.Decimal {
	n := kurang(kurang(kurang(nolBila(b.Premium), b.Commision), b.Deduction2), b.BrokerageFee)
	return bagiSkala(tambah(tambah(n, b.PPh), b.PPN), satu, 20)
}

// isiCicilan - nilai uang tiap cicilan dari persen (SaveFillPaymentInstallment 3.1.5.5 + CountPaymentInstallment_Act 3.3.1)
// dan SumTotalPayment / NetPremium mata uang.
func isiCicilan(m *models.MataUangPayment) {
	b := m.Baris
	m.NetPremium = netPremium(b)
	m.SumTotalPayment = m.NetPremium
	for i := range m.Cicilan {
		c := &m.Cicilan[i]
		p := nolBila(c.Persen)
		c.Premium = persenDari(b.Premium, p)
		c.RICommision = persenDari(b.Commision, p)
		c.Deduction2 = bagiSkala(kali(nolBila(b.Deduction2), p), seratus, 4)
		c.BrokerageFee = persenDari(b.BrokerageFee, p)
		c.PPh = persenDari(b.PPh, p)
		c.PPN = persenDari(b.PPN, p)
		c.Discount = nol() // .Policy.Payment.Diskon tidak pernah diisi di NB FacIn -> 0
		c.PaymentTotal = persenDari(m.NetPremium, p)
	}
}

// persenRata - FillPaymentInstallment 5.1.5.2-5.1.5.5: n cicilan @Math.divide(100, n, 20), terakhir menyerap sisa.
func persenRata(n int) []*apd.Decimal {
	bagian := bagiSkala(seratus, apd.New(int64(n), 0), 20)
	jalan := nol()
	out := make([]*apd.Decimal, n)
	for i := 1; i <= n; i++ {
		jalan = tambah(jalan, bagian)
		out[i-1] = bagian
		if i == n {
			out[i-1] = tambah(bagian, kurang(seratus, jalan))
		}
	}
	return out
}

// tambahBulan - Calendar.add(MONTH, n) Java: hari dipotong ke akhir bulan (31 Jan + 1 bulan = 28/29 Feb), bukan dilimpahkan
// seperti time.AddDate.
func tambahBulan(t time.Time, n int) time.Time {
	awal := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).AddDate(0, n, 0)
	akhir := awal.AddDate(0, 1, -1).Day()
	hari := t.Day()
	if hari > akhir {
		hari = akhir
	}
	return time.Date(awal.Year(), awal.Month(), hari, 0, 0, 0, 0, t.Location())
}

// jatuhTempoBawaan - FillPaymentInstallment @2693: Begin date (WIB) + 3 bulan, `YYYYMMDD`; "" bila Begin date kosong.
func jatuhTempoBawaan(mulai string) string {
	t, err := time.Parse(BentukWaktuPega, mulai)
	if err != nil {
		return ""
	}
	return tambahBulan(t.In(WIB), bulanJatuhTempo).Format("20060102")
}

// jumlahPersen - Σ persen cicilan.
func jumlahPersen(c []models.CicilanPayment) *apd.Decimal {
	j := nol()
	for _, x := range c {
		j = tambah(j, x.Persen)
	}
	return j
}

// pesanPersen - SetValidateInstallment_Act 3-5.
func pesanPersen(c []models.CicilanPayment) []string {
	switch j := jumlahPersen(c); {
	case lebihBesar(j, seratus):
		return []string{pesanPersenLebih100}
	case j.Cmp(seratus) < 0:
		return []string{pesanPersenKurang100}
	}
	return nil
}

func bulat4(x *apd.Decimal) *apd.Decimal { return bagiSkala(x, satu, 4) }

// selisihKecil - |a − b| < 0.001.
func selisihKecil(a, b *apd.Decimal) bool {
	d := kurang(a, b)
	if d.Sign() < 0 {
		d.Neg(d)
	}
	return d.Cmp(apd.New(1, -3)) < 0
}

// proteksiPayment - ProtectPremiPolicy_Act (NB, bukan Life / EDM / MOP): pesan verbatim urut Page-Set-Messages Pega
// (4 List Payment, per mata uang 5.7 List Installment, 10 Total premium, 11 Premi 0).
func proteksiPayment(mu []models.MataUangPayment) []string {
	if len(mu) == 0 {
		return []string{pesanPaymentKosong}
	}
	var m []string
	premiNol, totalPremiAll, tidakCocok, cicilanKosong := false, nol(), false, false
	for _, x := range mu {
		nilaiPremi := bulat4(x.Baris.Premium)
		if nilaiPremi.IsZero() {
			premiNol = true
		}
		totalPremiIns, commIns, totalIns := nol(), nol(), nol()
		for _, c := range x.Cicilan {
			if bulat4(c.Premium).IsZero() {
				premiNol = true
			}
			totalPremiIns, commIns, totalIns = tambah(totalPremiIns, c.Premium), tambah(commIns, c.RICommision), tambah(totalIns, c.PaymentTotal)
		}
		totalPremiIns, commIns, totalIns = bulat4(totalPremiIns), bulat4(commIns), bulat4(totalIns)
		totalPremiAll = tambah(totalPremiAll, nilaiPremi)
		totalComm, totalPayment := bulat4(x.Baris.Commision), bulat4(x.NetPremium)
		// 5.5: When 1 "ada yang beda" lalu When 2 berantai OR "ada satu saja yang selisihnya < 0.001" -> dilewati (verbatim:
		// pesan hanya bila KETIGA selisih >= 0.001, K51-6).
		beda := totalPremiIns.Cmp(nilaiPremi) != 0 || commIns.Cmp(totalComm) != 0 || totalIns.Cmp(totalPayment) != 0
		kecil := selisihKecil(totalPremiIns, nilaiPremi) || selisihKecil(commIns, totalComm) || selisihKecil(totalIns, totalPayment)
		if beda && !kecil {
			tidakCocok = true
		}
		if len(x.Cicilan) == 0 {
			cicilanKosong = true
		}
	}
	if cicilanKosong {
		m = append(m, pesanCicilanKosong)
	}
	if tidakCocok {
		m = append(m, pesanTotalTakCocok)
	}
	if premiNol && totalPremiAll.Sign() <= 0 {
		m = append(m, pesanPremiNol)
	}
	return m
}

func bulatkanPayment(mu []models.MataUangPayment) {
	for i := range mu {
		x := &mu[i]
		x.SumTotalPayment, x.NetPremium = simpan8(x.SumTotalPayment), simpan8(x.NetPremium)
		for j := range x.Cicilan {
			c := &x.Cicilan[j]
			c.Premium, c.RICommision, c.Deduction2, c.BrokerageFee = simpan8(c.Premium), simpan8(c.RICommision), simpan8(c.Deduction2), simpan8(c.BrokerageFee)
			c.PPh, c.PPN, c.Discount, c.PaymentTotal = simpan8(c.PPh), simpan8(c.PPN), simpan8(c.Discount), simpan8(c.PaymentTotal)
		}
	}
}

// ---- alur -------------------------------------------------------------------------------------------------------

// dasarPayment - nilai Payment per mata uang (tampilan tab Premium Deduction), data Payment tersimpan, Begin date.
func (s *Service) dasarPayment(ctx context.Context, id string) (TampilanDeduction, models.KasusPayment, string, error) {
	if s.payment == nil || s.deduction == nil || s.spreading == nil {
		return TampilanDeduction{}, models.KasusPayment{}, "", ErrPaymentTanpaDatabase
	}
	if !idKasusSah(id) {
		return TampilanDeduction{}, models.KasusPayment{}, "", ErrKasusTidakAda
	}
	k, err := s.payment.BacaPayment(ctx, id)
	if err != nil {
		return TampilanDeduction{}, k, "", galatKasus(err)
	}
	d, err := s.TampilanDeduction(ctx, id)
	if err != nil {
		return d, k, "", err
	}
	sp, err := s.spreading.BacaSpreading(ctx, nil, id)
	if err != nil {
		return d, k, "", galatKasus(err)
	}
	return d, k, sp.StartDateTime, nil
}

func tampilanPayment(d TampilanDeduction, installment string, mu []models.MataUangPayment, pesan []string) TampilanPayment {
	bulatkanPayment(mu)
	return TampilanPayment{Installment: installment, TarifPPH: d.TarifPPH, TarifPPN: d.TarifPPN, MataUang: mu,
		Pesan: append(append([]string{}, d.Pesan...), pesan...)}
}

// TampilanPayment - GET /api/nbfacin/kasus/{caseId}/payment: cicilan tersimpan (No, DueDate, persen) dengan nilai uang
// dihitung ulang dari nilai Payment mata uang terkini (Pega menghitung ulang cicilan di setiap Save, K51-4).
func (s *Service) TampilanPayment(ctx context.Context, id string) (TampilanPayment, error) {
	d, k, _, err := s.dasarPayment(ctx, id)
	if err != nil {
		return TampilanPayment{}, err
	}
	mu := make([]models.MataUangPayment, 0, len(d.MataUang))
	for _, b := range d.MataUang {
		m := models.MataUangPayment{Baris: b, Cicilan: append([]models.CicilanPayment{}, k.Cicilan[b.Currency]...)}
		isiCicilan(&m)
		mu = append(mu, m)
	}
	return tampilanPayment(d, k.Installment, mu, nil), nil
}

// jumlahCicilan - Installment wajib 1..99.
func jumlahCicilan(teks string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(teks))
	if err != nil || n < 1 || n > banyakCicilanMaks {
		return 0, fmt.Errorf("%w: installment wajib angka bulat 1-%d", ErrMasukanPayment, banyakCicilanMaks)
	}
	return n, nil
}

// cocokkanCicilan - cicilan badan dipasang ke mata uang server menurut KODE (himpunan wajib sama, 409).
func cocokkanCicilan(d TampilanDeduction, isian []IsianCicilanMataUang) ([]models.MataUangPayment, error) {
	if len(isian) != len(d.MataUang) {
		return nil, ErrPaymentBerubah
	}
	badan := map[string][]models.CicilanPayment{}
	for _, x := range isian {
		if _, ganda := badan[x.Currency]; ganda {
			return nil, fmt.Errorf("%w: mata uang %s ganda", ErrMasukanPayment, x.Currency)
		}
		badan[x.Currency] = x.Cicilan
	}
	mu := make([]models.MataUangPayment, 0, len(d.MataUang))
	for _, b := range d.MataUang {
		c, ada := badan[b.Currency]
		if !ada {
			return nil, ErrPaymentBerubah
		}
		mu = append(mu, models.MataUangPayment{Baris: b, Cicilan: append([]models.CicilanPayment{}, c...)})
	}
	return mu, nil
}

// HitungPayment - POST …/payment/hitung: FillPaymentInstallment (aksi installment, semua mata uang) atau
// SetValidateInstallment_Act (aksi persen, mata uang `currency` baris `idx`), TANPA simpan.
func (s *Service) HitungPayment(ctx context.Context, id, aksi, currency string, idx int, isian IsianPayment) (TampilanPayment, error) {
	if aksi != AksiInstallment && aksi != AksiPersenCicilan {
		return TampilanPayment{}, fmt.Errorf("%w: aksi harus installment atau persen", ErrMasukanPayment)
	}
	d, _, mulai, err := s.dasarPayment(ctx, id)
	if err != nil {
		return TampilanPayment{}, err
	}
	var pesan []string
	var mu []models.MataUangPayment
	if aksi == AksiInstallment {
		n, err := jumlahCicilan(isian.Installment)
		if err != nil {
			return TampilanPayment{}, err
		}
		rata, jatuh := persenRata(n), jatuhTempoBawaan(mulai)
		for _, b := range d.MataUang {
			m := models.MataUangPayment{Baris: b}
			for i := 0; i < n; i++ {
				m.Cicilan = append(m.Cicilan, models.CicilanPayment{InstallmentNo: i + 1, DueDate: jatuh, Persen: rata[i]})
			}
			mu = append(mu, m)
		}
	} else {
		if mu, err = cocokkanCicilan(d, isian.MataUang); err != nil {
			return TampilanPayment{}, err
		}
		ketemu := false
		for _, m := range mu {
			if m.Baris.Currency != currency {
				continue
			}
			ketemu = true
			if idx < 0 || idx >= len(m.Cicilan) {
				return TampilanPayment{}, fmt.Errorf("%w: idx %d di luar cicilan %s", ErrMasukanPayment, idx, currency)
			}
			pesan = pesanPersen(m.Cicilan)
		}
		if !ketemu {
			return TampilanPayment{}, ErrPaymentBerubah
		}
	}
	for i := range mu {
		isiCicilan(&mu[i])
	}
	return tampilanPayment(d, strings.TrimSpace(isian.Installment), mu, pesan), nil
}

// SimpanPayment - PUT …/payment: SaveFillPaymentInstallment + CountPaymentInstallment_Act, ProtectPremiPolicy_Act (400),
// simpan, jawab baca ulang.
func (s *Service) SimpanPayment(ctx context.Context, pelaku inti.Pelaku, id string, isian IsianPayment) (TampilanPayment, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanPayment{}, err
	}
	if err := s.wajibBolehUbah(ctx, pelaku, id); err != nil {
		return TampilanPayment{}, err
	}
	n, err := jumlahCicilan(isian.Installment)
	if err != nil {
		return TampilanPayment{}, err
	}
	if s.transaksi == nil {
		return TampilanPayment{}, ErrPaymentTanpaDatabase
	}
	d, _, _, err := s.dasarPayment(ctx, id)
	if err != nil {
		return TampilanPayment{}, err
	}
	mu, err := cocokkanCicilan(d, isian.MataUang)
	if err != nil {
		return TampilanPayment{}, err
	}
	var masalah []string
	for i := range mu {
		m := &mu[i]
		if len(m.Baris.Currency) > lebarKodeMataUangCicilan {
			masalah = append(masalah, fmt.Sprintf("kode mata uang %s melebihi %d bita (T_LISTINSTALLMENT.CURRENCY_CODE)", m.Baris.Currency,
				lebarKodeMataUangCicilan))
			continue
		}
		if len(m.Cicilan) > 0 && len(m.Cicilan) != n {
			masalah = append(masalah, fmt.Sprintf("jumlah cicilan %s (%d) tidak sama dengan Installment (%d)", m.Baris.Currency, len(m.Cicilan), n))
			continue
		}
		for j := range m.Cicilan {
			if m.Cicilan[j].DueDate == "" {
				masalah = append(masalah, fmt.Sprintf("cicilan %s ke-%d: Due Date wajib diisi", m.Baris.Currency, j+1))
			}
		}
		// SaveFillPaymentInstallment 3.1.5.3 / 3.1.5.4: persen dipertahankan bila Σ = 100, selain itu diratakan ulang.
		if jumlahPersen(m.Cicilan).Cmp(seratus) != 0 {
			rata := persenRata(n)
			for j := range m.Cicilan {
				m.Cicilan[j].Persen = rata[j]
			}
		}
		for j := range m.Cicilan {
			m.Cicilan[j].InstallmentNo = j + 1
		}
		isiCicilan(m)
		if m.NetPremium.Sign() < 0 {
			masalah = append(masalah, fmt.Sprintf("Total Payment %s negatif - potongan melebihi premi", m.Baris.Currency))
		}
	}
	masalah = append(proteksiPayment(mu), masalah...)
	if len(masalah) > 0 {
		return TampilanPayment{}, fmt.Errorf("%w: %s", ErrMasukanPayment, strings.Join(masalah, "; "))
	}
	bulatkanPayment(mu)
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		return s.payment.TulisPayment(ctx, tx, id, models.SimpanPayment{Installment: n, MataUang: mu})
	})
	if err != nil {
		return TampilanPayment{}, galatKasus(err)
	}
	return s.TampilanPayment(ctx, id)
}
