package services

// Tab utama Loss Record kasus FIRE (tiket 55). Bukti `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\`:
//   - `Activity\SaveLossRecord.xml` - DOL kosong = Begin date − 1 tahun (7.3 / 7.6), CauseOfLoss huruf besar (7.4), Loss
//     Amount IDR = Loss Amount (IDR / mata uang kosong) atau × kurs (7.5), Uw Year / TahunKe dari jendela (7.7), Total Loss
//     per Class Of Business + Currency dan loss ratio (8), urut DOL terbaru (9);
//   - `RDBList\Get5YearAgo.xml` - jendela L = 1..5: [ADD_MONTHS(Begin, −12L), ADD_MONTHS(End, −12L)], Uw Year = tahun
//     ADD_MONTHS(Begin, −12L); kecocokan TERAKHIR (L terbesar) yang berlaku (loop Pega menimpa);
//   - `Activity\GetPremiPAR_ACT.xml` - Premi PAR kosong / 0 -> Σ premi × kurs, @divide(x, 1, 0).
// Keputusan work owner 06-10-2026 (docs/issues/55): tab untuk JOB_POSITION "IT" saja; cacat Pega DIBETULKAN - ember 3 / 5
// tahun (Pega `||` memasukkan semua), dua tombol simpan (di sini satu Save yang selalu menghitung), kurs per Begin date lewat
// pembaca Cedant (bukan BETWEEN atas VARCHAR), Premi PAR bawaan dari premi spread × kurs (Pega `.KursIDR` tak pernah diisi).

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrMasukanLossRecord - isian Loss Record tidak sah. 400.
	ErrMasukanLossRecord = errors.New("services: isian Loss Record tidak sah")
	// ErrLossRecordTanpaDatabase - pembaca Loss Record / kasus / spreading / deduction tidak terpasang. 503.
	ErrLossRecordTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, Loss Record tidak terbaca")
)

const (
	// jabatanTabIT - keputusan work owner 06-10-2026: Loss Record "untuk sementara untuk POSITION=IT saja" (tiket 55) dan
	// Correspondence "ini visibility menu nya juga cuma it ya" (tiket 56) - M_LOGIN_GO.JOB_POSITION, dipangkas, tanpa beda
	// huruf besar / kecil. Pega `When\IsIT.xml`: pyPosition "IT Developer" ATAU pyOrgDivision "IT".
	jabatanTabIT = "IT"
	// batasBarisLossRecord, lebar* - batas isian (CLASS_OF_BUSINESS = BUSINESS.NOTE 4000, CAUSE_OF_LOSS pola
	// T_LISTCAUSEOFLOSS 500, CURRENCY = kode mata uang 10).
	batasBarisLossRecord = 500
	lebarCOBLossRecord   = 4000
	lebarSebabLossRecord = 500
	lebarMataUangLossRec = 10
	// jendelaUwYear - Get5YearAgo `CONNECT BY LEVEL <= 5`.
	jendelaUwYear = 5
)

// TampilanLossRecord - jawaban seluruh endpoint tab Loss Record.
type TampilanLossRecord struct {
	LossRecord models.LossRecord
	TotalCob   []models.TotalLossRecord
	Pesan      []string
}

// DenganLossRecord memasang penyimpan Loss Record dan pembaca jabatan akun (tiket 55).
func (s *Service) DenganLossRecord(p repository.PenyimpanLossRecord, j repository.PembacaJabatan) *Service {
	s.lossRecord, s.jabatan = p, j
	return s
}

// ---- akses --------------------------------------------------------------------------------------------------------

// jabatanIT - akun pelaku ber-JOB_POSITION "IT"; tanpa identitas / tanpa pembaca -> false tanpa galat.
func (s *Service) jabatanIT(ctx context.Context, pelaku inti.Pelaku) (bool, error) {
	if strings.TrimSpace(pelaku.AkunID) == "" || s.jabatan == nil {
		return false, nil
	}
	j, err := s.jabatan.JabatanAkun(ctx, pelaku.AkunID)
	if err != nil {
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(j), jabatanTabIT), nil
}

// AksesTabIT - `aksesLossRecord` dan `aksesKorespondensi` GET kasus (satu pembacaan JOB_POSITION untuk keduanya - aturannya
// sama). LUNAK: galat baca jabatan -> false (dicatat), supaya layar kasus tidak pernah gagal karena tab ini.
func (s *Service) AksesTabIT(ctx context.Context, pelaku inti.Pelaku) bool {
	ok, err := s.jabatanIT(ctx, pelaku)
	if err != nil {
		log.Printf("nbfacin: akses tab IT (Loss Record / Correspondence) dianggap tidak ada: %v", err)
	}
	return ok
}

// wajibJabatanIT - penegakan di layanan (ADR-0030): tanpa identitas 401, bukan "IT" 403 (galat baca jabatan diteruskan).
func (s *Service) wajibJabatanIT(ctx context.Context, pelaku inti.Pelaku, tab string) error {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	ok, err := s.jabatanIT(ctx, pelaku)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%w: tab %s hanya untuk JOB_POSITION %s", inti.ErrTanpaWewenang, tab, jabatanTabIT)
	}
	return nil
}

// WajibAksesLossRecord - lihat wajibJabatanIT. Dipanggil handler SEBELUM mengurai badan (403 mendahului 400) dan sekali
// lagi oleh setiap metode tab.
func (s *Service) WajibAksesLossRecord(ctx context.Context, pelaku inti.Pelaku) error {
	return s.wajibJabatanIT(ctx, pelaku, "Loss Record")
}

// ---- jendela Uw Year ----------------------------------------------------------------------------------------------

// tambahBulanOracle - ADD_MONTHS: hari terakhir bulan -> hari terakhir bulan tujuan; selain itu hari dijepit.
func tambahBulanOracle(t time.Time, n int) time.Time {
	h := tambahBulan(t, n)
	if t.AddDate(0, 0, 1).Month() != t.Month() { // t hari terakhir bulannya
		return time.Date(h.Year(), h.Month()+1, 0, 0, 0, 0, 0, h.Location())
	}
	return h
}

type jendelaLossRecord struct {
	mulai, akhir time.Time
	uwYear       string
	ke           int
}

func jendelaUw(mulai, akhir time.Time) []jendelaLossRecord {
	var j []jendelaLossRecord
	for l := 1; l <= jendelaUwYear; l++ {
		m := tambahBulanOracle(mulai, -12*l)
		j = append(j, jendelaLossRecord{mulai: m, akhir: tambahBulanOracle(akhir, -12*l), uwYear: fmt.Sprintf("%04d", m.Year()), ke: l})
	}
	return j
}

// ---- hitung -------------------------------------------------------------------------------------------------------

// konteksLossRecord - yang dibaca untuk satu hitung.
type konteksLossRecord struct {
	mulaiPega          string    // Begin date teks Pega (kurs)
	mulai, akhir       time.Time // tanggal WIB
	adaMulai, adaAkhir bool
	riComm             *apd.Decimal // % RI Commission (Premium Deduction)
	premiSpread        []jumlahMataUang
}

// kursLossRecord - kurs mata uang ke IDR pada Begin date (pembaca Cedant: IDR 10026 -> 1); nil = tidak ada.
func (s *Service) kursLossRecord(ctx context.Context, cache map[string]*apd.Decimal, mataUang, mulai string) (*apd.Decimal, error) {
	if k, ada := cache[mataUang]; ada {
		return k, nil
	}
	id, err := s.spreading.IDMataUang(ctx, mataUang)
	if err != nil {
		return nil, err
	}
	var k *apd.Decimal
	switch {
	case id == idMataUangIDR:
		k = apd.New(1, 0)
	case id != "":
		if k, err = s.spreading.KursPada(ctx, id, mulai); err != nil {
			return nil, err
		}
	}
	cache[mataUang] = k
	return k, nil
}

// hitungLossRecord - masukan (DateOfLoss DD-MM-YYYY / "") -> baris terhitung (DateOfLoss tetap DD-MM-YYYY), Total per COB +
// Currency, rasio, pesan. Urutan baris DIPERTAHANKAN.
func (s *Service) hitungLossRecord(ctx context.Context, k konteksLossRecord, m models.LossRecord) (models.LossRecord, []models.TotalLossRecord, []string, error) {
	pesan := []string{}
	l := models.LossRecord{PremiPar: m.PremiPar, PremiMb: m.PremiMb, PremiPl: m.PremiPl, PremiOther: m.PremiOther,
		Baris: []models.BarisLossRecord{}}
	// Premi PAR bawaan (GetPremiPAR_ACT, dibetulkan): kosong / 0 -> Σ premi spread × kurs, 0 desimal.
	if nolBila(l.PremiPar).Sign() == 0 && len(k.premiSpread) > 0 {
		_, premi, tidakAda, err := totalKeIDR(ctx, s.spreading, k.mulaiPega, k.premiSpread)
		if err != nil {
			return l, nil, nil, err
		}
		if premi == nil {
			pesan = append(pesan, "Kurs "+strings.Join(tidakAda, ", ")+" pada Begin date tidak ada di TREATYEXCHANGEYEARLY - Premi PAR bawaan tidak diisi")
		} else {
			l.PremiPar = bulatSetengahAtas(premi, 0)
		}
	}
	var jendela []jendelaLossRecord
	if k.adaMulai && k.adaAkhir {
		jendela = jendelaUw(k.mulai, k.akhir)
	} else if len(m.Baris) > 0 {
		pesan = append(pesan, "Begin / End date kasus belum diisi - Uw Year tidak dihitung")
	}
	cache := map[string]*apd.Decimal{}
	tanpaKurs := map[string]bool{}
	tahun := map[int]*apd.Decimal{}
	for _, b := range m.Baris {
		b.CauseOfLoss = strings.ToUpper(b.CauseOfLoss)
		b.Currency = strings.ToUpper(strings.TrimSpace(b.Currency))
		b.UwYear, b.TahunKe, b.LossAmountIDR, b.KursIDR = "", 0, nil, nil
		if b.DateOfLoss == "" && k.adaMulai {
			b.DateOfLoss = tambahBulan(k.mulai, -12).Format(BentukTanggalKabel) // 7.6 Begin − 1 tahun (Calendar.add)
		}
		if tgl, err := time.Parse(BentukTanggalKabel, b.DateOfLoss); err == nil {
			for _, j := range jendela {
				if !tgl.Before(j.mulai) && !tgl.After(j.akhir) {
					b.UwYear, b.TahunKe = j.uwYear, j.ke // kecocokan terakhir menang
				}
			}
		}
		if b.Currency == "" || b.Currency == "IDR" {
			b.LossAmountIDR = b.LossAmount
		} else {
			kurs, err := s.kursLossRecord(ctx, cache, b.Currency, k.mulaiPega)
			if err != nil {
				return l, nil, nil, err
			}
			if kurs == nil {
				if !tanpaKurs[b.Currency] {
					tanpaKurs[b.Currency] = true
					pesan = append(pesan, "Kurs "+b.Currency+" pada Begin date tidak ada di TREATYEXCHANGEYEARLY - Loss Amount (IDR) tidak diisi")
				}
			} else {
				b.KursIDR = kurs
				if b.LossAmount != nil {
					b.LossAmountIDR = simpan8(kali(b.LossAmount, kurs))
				}
			}
		}
		if b.TahunKe > 0 && b.LossAmountIDR != nil {
			tahun[b.TahunKe] = tambah(tahun[b.TahunKe], b.LossAmountIDR)
		}
		l.Baris = append(l.Baris, b)
	}
	l.Rasio, pesan = rasioLossRecord(l, tahun, k.riComm, pesan)
	if err := periksaMuatLossRecord(l); err != nil {
		return l, nil, nil, err
	}
	return l, totalCobLossRecord(l.Baris), pesan, nil
}

// periksaMuatLossRecord - hasil yang disimpan wajib muat NUMBER(38,8) (30 digit bulat): 400, bukan ORA-01438 (pola
// coverage). Mis. Gross premi sangat kecil -> rasio sangat besar.
func periksaMuatLossRecord(l models.LossRecord) error {
	r := l.Rasio
	nilai := map[string]*apd.Decimal{"premiPar": l.PremiPar, "premiMb": l.PremiMb, "premiPl": l.PremiPl, "premiOther": l.PremiOther,
		"rasio.gross1": r.Gross1, "rasio.gross3": r.Gross3, "rasio.gross5": r.Gross5, "rasio.net1": r.Net1, "rasio.net3": r.Net3,
		"rasio.net5": r.Net5}
	for i, b := range l.Baris {
		nilai[fmt.Sprintf("baris[%d].lossAmountIdr", i)] = b.LossAmountIDR
	}
	var masalah []string
	for medan, d := range nilai {
		if !desimalSah(d) {
			masalah = append(masalah, "hasil "+medan+" melebihi 30 digit bulat")
		}
	}
	if len(masalah) > 0 {
		sort.Strings(masalah)
		return fmt.Errorf("%w: %s", ErrMasukanLossRecord, strings.Join(masalah, "; "))
	}
	return nil
}

// totalCobLossRecord - TotalLossCOB: per (Class Of Business, Currency) urut kemunculan; IDR kosong dilewati (nil bila
// semuanya kosong).
func totalCobLossRecord(baris []models.BarisLossRecord) []models.TotalLossRecord {
	t := []models.TotalLossRecord{}
	for _, b := range baris {
		n := -1
		for i := range t {
			if t[i].ClassOfBusiness == b.ClassOfBusiness && t[i].Currency == b.Currency {
				n = i
			}
		}
		if n < 0 {
			t = append(t, models.TotalLossRecord{ClassOfBusiness: b.ClassOfBusiness, Currency: b.Currency, LossAmount: nol()})
			n = len(t) - 1
		}
		t[n].LossAmount = tambah(t[n].LossAmount, b.LossAmount)
		if b.LossAmountIDR != nil {
			t[n].LossAmountIDR = tambah(t[n].LossAmountIDR, b.LossAmountIDR)
		}
	}
	return t
}

// rasioLossRecord - SaveLossRecord 8.6 DIBETULKAN: Y1 = TahunKe 1, Y3 = 1..3, Y5 = 1..5; Gross = PAR + MB + PL + Other;
// Deduction = RI Commission / 100 eksak (Pega @divide(…,100,2) - tanpa pembulatan, keputusan work owner); grossN =
// YN / (Gross × N) × 100, netN = YN / (Gross × (1 − Deduction) × N) × 100; skala 20 lalu 8 desimal.
func rasioLossRecord(l models.LossRecord, tahun map[int]*apd.Decimal, riComm *apd.Decimal, pesan []string) (models.RasioLossRecord, []string) {
	var r models.RasioLossRecord
	gross := nol()
	for _, p := range []*apd.Decimal{l.PremiPar, l.PremiMb, l.PremiPl, l.PremiOther} {
		gross = tambah(gross, p)
	}
	if gross.Sign() == 0 {
		return r, append(pesan, "Gross premi (PAR + MB + PL + Other) 0 - loss ratio tidak dihitung")
	}
	y := func(n int) *apd.Decimal {
		s := nol()
		for i := 1; i <= n; i++ {
			s = tambah(s, tahun[i])
		}
		return s
	}
	faktorNet := kurang(apd.New(1, 0), bagiSkala(nolBila(riComm), seratus, 20))
	rasio := func(n int, faktor *apd.Decimal) *apd.Decimal {
		return simpan8(kali(bagiSkala(y(n), kali(gross, faktor, apd.New(int64(n), 0)), 20), seratus))
	}
	r.Gross1, r.Gross3, r.Gross5 = rasio(1, apd.New(1, 0)), rasio(3, apd.New(1, 0)), rasio(5, apd.New(1, 0))
	if faktorNet.Sign() <= 0 {
		return r, append(pesan, "RI Commission 100% atau lebih - net loss ratio tidak dihitung")
	}
	r.Net1, r.Net3, r.Net5 = rasio(1, faktorNet), rasio(3, faktorNet), rasio(5, faktorNet)
	return r, pesan
}

// ---- alur ---------------------------------------------------------------------------------------------------------

func (s *Service) bacaKonteksLossRecord(ctx context.Context, id string) (konteksLossRecord, error) {
	var k konteksLossRecord
	kasus, err := s.kasus.BacaKasus(ctx, id)
	if err != nil {
		return k, galatKasus(err)
	}
	k.mulaiPega = kasus.General.StartDateTime
	k.mulai, k.adaMulai = tanggalJakarta(kasus.General.StartDateTime)
	k.akhir, k.adaAkhir = tanggalJakarta(kasus.General.EndDateTime)
	d, err := s.deduction.BacaDeduction(ctx, id)
	if err != nil {
		return k, galatKasus(err)
	}
	k.riComm = d.RICommision
	sp, err := s.spreading.BacaSpreading(ctx, nil, id)
	if err != nil {
		return k, galatKasus(err)
	}
	adaTop := totalLokasi(&sp)
	for _, t := range totalMataUang(sp, adaTop) {
		mu := t.Currency
		if mu == "" {
			mu = mataUangKosong
		}
		n := -1
		for i := range k.premiSpread {
			if k.premiSpread[i].Currency == mu {
				n = i
			}
		}
		if n < 0 {
			k.premiSpread = append(k.premiSpread, jumlahMataUang{Currency: mu, Premi: nol()})
			n = len(k.premiSpread) - 1
		}
		k.premiSpread[n].Premi = tambah(k.premiSpread[n].Premi, t.PremiumSpreaded)
	}
	return k, nil
}

func (s *Service) siapLossRecord(id string) error {
	if s.lossRecord == nil || s.kasus == nil || s.spreading == nil || s.deduction == nil {
		return ErrLossRecordTanpaDatabase
	}
	if !idKasusSah(id) {
		return ErrKasusTidakAda
	}
	return nil
}

// TampilanLossRecordKasus - GET: isi tersimpan, dihitung ulang (DOL tersimpan teks Pega -> DD-MM-YYYY).
func (s *Service) TampilanLossRecordKasus(ctx context.Context, pelaku inti.Pelaku, id string) (TampilanLossRecord, error) {
	if err := s.WajibAksesLossRecord(ctx, pelaku); err != nil {
		return TampilanLossRecord{}, err
	}
	if err := s.siapLossRecord(id); err != nil {
		return TampilanLossRecord{}, err
	}
	tersimpan, err := s.lossRecord.BacaLossRecord(ctx, id)
	if err != nil {
		return TampilanLossRecord{}, galatKasus(err)
	}
	for i := range tersimpan.Baris {
		tersimpan.Baris[i].DateOfLoss = tanggalKerugianKeKabel(tersimpan.Baris[i].DateOfLoss)
	}
	return s.hitungLossRecordKasus(ctx, id, tersimpan)
}

func (s *Service) hitungLossRecordKasus(ctx context.Context, id string, m models.LossRecord) (TampilanLossRecord, error) {
	k, err := s.bacaKonteksLossRecord(ctx, id)
	if err != nil {
		return TampilanLossRecord{}, err
	}
	l, total, pesan, err := s.hitungLossRecord(ctx, k, m)
	return TampilanLossRecord{LossRecord: l, TotalCob: total, Pesan: pesan}, err
}

// periksaLossRecord - batas baris / lebar teks, DateOfLoss DD-MM-YYYY atau kosong.
func periksaLossRecord(m models.LossRecord) error {
	var masalah []string
	if len(m.Baris) > batasBarisLossRecord {
		masalah = append(masalah, fmt.Sprintf("baris paling banyak %d", batasBarisLossRecord))
	}
	for i, b := range m.Baris {
		for _, c := range []struct {
			nama, nilai string
			lebar       int
		}{{"classOfBusiness", b.ClassOfBusiness, lebarCOBLossRecord}, {"causeOfLoss", b.CauseOfLoss, lebarSebabLossRecord},
			{"currency", strings.TrimSpace(b.Currency), lebarMataUangLossRec}} {
			if len(c.nilai) > c.lebar {
				masalah = append(masalah, fmt.Sprintf("baris[%d].%s paling banyak %d bita", i, c.nama, c.lebar))
			}
		}
		if b.DateOfLoss != "" {
			if t, err := time.Parse(BentukTanggalKabel, b.DateOfLoss); err != nil || t.Format(BentukTanggalKabel) != b.DateOfLoss {
				masalah = append(masalah, fmt.Sprintf("baris[%d].dateOfLoss harus DD-MM-YYYY", i))
			}
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanLossRecord, strings.Join(masalah, "; "))
	}
	return nil
}

// HitungLossRecord - POST …/loss-record/hitung: tanpa menyimpan, urutan baris dipertahankan.
func (s *Service) HitungLossRecord(ctx context.Context, pelaku inti.Pelaku, id string, m models.LossRecord) (TampilanLossRecord, error) {
	if err := s.WajibAksesLossRecord(ctx, pelaku); err != nil {
		return TampilanLossRecord{}, err
	}
	if err := periksaLossRecord(m); err != nil {
		return TampilanLossRecord{}, err
	}
	if err := s.siapLossRecord(id); err != nil {
		return TampilanLossRecord{}, err
	}
	return s.hitungLossRecordKasus(ctx, id, m)
}

// kanonKelasBisnis - Save saja (keputusan work owner 06-10-2026: COB WAJIB dari daftar BUSINESS.NOTE): tiap baris wajib
// terisi dan sama persis (dipangkas, tanpa beda huruf) dengan satu BUSINESS.NOTE; diganti ejaan NOTE tabel. Pesan per
// baris (nomor mulai 1). Hitung TIDAK memeriksanya (layar menghitung sambil mengetik).
func (s *Service) kanonKelasBisnis(ctx context.Context, m *models.LossRecord) error {
	if len(m.Baris) == 0 {
		return nil
	}
	if s.kelasBisnis == nil {
		return ErrKelasBisnisTanpaDatabase
	}
	var kunci []string
	lihat := map[string]bool{}
	for _, b := range m.Baris {
		if k := strings.ToUpper(strings.TrimSpace(b.ClassOfBusiness)); k != "" && !lihat[k] {
			lihat[k] = true
			kunci = append(kunci, k)
		}
	}
	note, err := s.kelasBisnis.NoteKelasBisnis(ctx, kunci)
	if err != nil {
		return err
	}
	var masalah []string
	baris := make([]models.BarisLossRecord, len(m.Baris))
	for i, b := range m.Baris {
		k := strings.ToUpper(strings.TrimSpace(b.ClassOfBusiness))
		n, ada := note[k]
		switch {
		case k == "":
			masalah = append(masalah, fmt.Sprintf("baris %d: Class Of Business wajib diisi", i+1))
		case !ada:
			masalah = append(masalah, fmt.Sprintf("baris %d: Class Of Business %q tidak ada di daftar BUSINESS", i+1, b.ClassOfBusiness))
		default:
			b.ClassOfBusiness = n
		}
		baris[i] = b
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanLossRecord, strings.Join(masalah, "; "))
	}
	m.Baris = baris
	return nil
}

// SimpanLossRecord - PUT …/loss-record: akses, isian, COB dari daftar BUSINESS (ejaan tabel), hitung, urut DOL terbaru dulu
// (langkah 9, stabil; tanpa DOL di akhir), ganti utuh.
func (s *Service) SimpanLossRecord(ctx context.Context, pelaku inti.Pelaku, id string, m models.LossRecord) (TampilanLossRecord, error) {
	if err := s.WajibAksesLossRecord(ctx, pelaku); err != nil {
		return TampilanLossRecord{}, err
	}
	if err := s.wajibTahapUbah(ctx, id); err != nil {
		return TampilanLossRecord{}, err
	}
	if err := periksaLossRecord(m); err != nil {
		return TampilanLossRecord{}, err
	}
	if err := s.kanonKelasBisnis(ctx, &m); err != nil {
		return TampilanLossRecord{}, err
	}
	t, err := s.HitungLossRecord(ctx, pelaku, id, m)
	if err != nil {
		return TampilanLossRecord{}, err
	}
	if s.transaksi == nil {
		return TampilanLossRecord{}, ErrLossRecordTanpaDatabase
	}
	tgl := func(b models.BarisLossRecord) time.Time {
		v, _ := time.Parse(BentukTanggalKabel, b.DateOfLoss)
		return v
	}
	sort.SliceStable(t.LossRecord.Baris, func(i, j int) bool { return tgl(t.LossRecord.Baris[i]).After(tgl(t.LossRecord.Baris[j])) })
	t.TotalCob = totalCobLossRecord(t.LossRecord.Baris) // urut kemunculan atas baris TERURUT - sama dengan GET berikutnya
	simpan := t.LossRecord
	simpan.Baris = make([]models.BarisLossRecord, len(t.LossRecord.Baris))
	for i, b := range t.LossRecord.Baris {
		b.DateOfLoss = tanggalKerugianKePega(b.DateOfLoss)
		simpan.Baris[i] = b
	}
	if err := s.transaksi(ctx, func(tx *db.Tx) error { return s.lossRecord.TulisLossRecord(ctx, tx, id, simpan) }); err != nil {
		return TampilanLossRecord{}, galatKasus(err)
	}
	return t, nil
}

// SaranKelasBisnis - GET /api/nbfacin/class-of-business/saran?q= (tiket 55).
func (s *Service) SaranKelasBisnis(ctx context.Context, kata string) ([]models.SaranAkumulasi, error) {
	if err := periksaParamKueri(ErrMasukanLossRecord, paramAkumulasi{"q", &kata}); err != nil {
		return nil, err
	}
	if s.kelasBisnis == nil {
		return nil, ErrKelasBisnisTanpaDatabase
	}
	return s.kelasBisnis.SaranKelasBisnis(ctx, kata)
}
