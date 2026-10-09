package services

// RUMUS TAB INSTALLMENT (cabang NON-PROPORSIONAL) - dari ekspor Pega.
//
// ---------------------------------------------------------------------
// ⭐ TIGA ACTIVITY, LIMA PEMICU
// ---------------------------------------------------------------------
//
//	Activity/TreatyInSetValueInstallment.xml
//	    isian `Installment` - perilaku `change`, TANPA `status`
//	    tombol `Update Value` - `status=update`
//	Activity/SetTotalInstallment.xml (rincian baris, Section/Installments.xml)
//	    sel `% Installment` - `status=editpercentage`: Amount dihitung ulang
//	    sel `Amount` - tanpa parameter: hanya total
//	Activity/TreatyInNPSetTotal.xml langkah 27-28
//	    tombol `Update Total` - `type=installment`
//
// Ketiganya IDENTIK di korpus Treaty In dan Treaty In Adjustment (diadu
// langkah demi langkah 7 Oktober 2026); rute ini melayani kedua layar.
//
// ---------------------------------------------------------------------
// TreatyInSetValueInstallment
// ---------------------------------------------------------------------
//
//	 1  status=update → TempDate.Installment = TreatyIn.Installment
//	 2  Property-Remove TreatyIn.Installment, TreatyIn.TotalInstallmentNP
//	 3  TreatyIn.InstallmentNo = Param.Installment
//	 4  Installment < 1 → pesan, KELUAR (grid sudah kosong: langkah 2)
//	 5  TotalShareNetNP kosong → pesan, KELUAR
//	 6  satu Installment{Currency} per baris TotalShareNetNP
//	 7  Percentage = @Math.divide(100, InstallmentNo, 2)
//	 8  per Installment: value = TotalShareNetNP.Value bermata uang sama
//	    (8.1.1 dan 8.1.2 berprasyarat FacultativeShare ==0 / !=0 dan
//	    menulis HAL YANG SAMA); REPEAT 1..InstallmentNo:
//	      n++, pct += Percentage, baris {Installment n, DueDate hari ini, Currency}
//	      n == N && pct != 100 → InstallmentPct = Percentage + (100 - pct)
//	      n <  N || pct == 100 → InstallmentPct = Percentage
//	      Amount = @divide(InstallmentPct, 100, 4) * value
//	 9  status=update → DueDate, WPC, PaymentDate, InstallmentPct DIPULIHKAN
//	    dari TempDate, indeks sama
//	10  EDMState bukan "" dan bukan "0" (Adjustment) → dipulihkan dari
//	    TreatyIn.OLDDATA.Installment, indeks sama
//	11  SetTotalInstallment TANPA parameter → lompat ke label `calc`
//	    (langkah 4): AmountTotal = Σ Amount, PctTotal = Σ InstallmentPct.
//	    ⛔ Amount TIDAK dihitung ulang: persen yang dipulihkan di langkah 9/10
//	    tidak mengubah Amount hasil langkah 8. Disalin apa adanya.
//	12  TotalInstallmentNP = {Currency, Value: AmountTotal} per Installment
//
// ---------------------------------------------------------------------
// ⚠️ YANG TIDAK BERBUKTI DI EKSPOR - DIPUTUSKAN DI SINI, DINYATAKAN
// ---------------------------------------------------------------------
//   - Langkah 9/10 membaca `(local.idx).InstallmentList(<CURRENT>)` yang bisa
//     TIDAK ADA (jumlah angsuran naik, atau sisi Old lebih pendek). Perilaku
//     Pega atas halaman yang tidak ada tidak terekspor; di sini nilai hasil
//     hitungan DIPERTAHANKAN.
//   - `param.Installment < 1` atas isian kosong / bukan angka: diperlakukan
//     sebagai < 1 (pesan yang sama).
//   - DueDate `@CurrentDate("dd/MM/yyyy","Asia/Jakarta")` ditulis dalam
//     bentuk TERSIMPAN properti Date (`YYYYMMDD`): panjang maksimum
//     `InstallmentList.DueDate` di data Pega = 8
//     (`docs/STRUKTUR-TABEL-TREATY-IN.md`, tabel ukuran).
//
// Pembulatan `@Math.divide`/`@divide` = HALF_UP pada skala yang disebut -
// dikonfirmasi modul endorsmentfacin (`pembagian.go`, A37).
//
// ⛔ Pagar: lebih dari 1000 angsuran DITOLAK (`ErrMasukanTidakSah`) - REPEAT
// Pega tidak berbatas, rute HTTP harus.

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/penomor"
)

// Aksi tab Installment.
const (
	// AksiAngsuranNilai - `TreatyInSetValueInstallment` (isian Installment / Update Value).
	AksiAngsuranNilai = "nilai"
	// AksiAngsuranTotalBaris - `SetTotalInstallment` atas SATU halaman Installment.
	AksiAngsuranTotalBaris = "total-baris"
	// AksiAngsuranTotal - tombol `Update Total` (`TreatyInNPSetTotal` installment).
	AksiAngsuranTotal = "total"
	// AksiAngsuranTanggalBayar - `TreatyInUpdatePaymentDate` (sel Due Date):
	// SATU halaman Installment (`Indeks`), SETIAP baris.
	AksiAngsuranTanggalBayar = "tanggal-bayar"
	// AksiAngsuranTanggalBayarSemua - `TreatyInUpdatePaymentDate_Act` (sel
	// WPC): SEMUA halaman `TreatyIn.Installment`, baris ber-Due Date.
	AksiAngsuranTanggalBayarSemua = "tanggal-bayar-semua"
)

// Status parameter Activity.
const (
	statusAngsuranUpdate = "update"
	statusEditPersen     = "editpercentage"
)

// Pesan `Property-Set-Messages` langkah 4 dan 5, apa adanya.
const (
	pesanAngsuranKurang = "Installment Value cannot be less than 1"
	pesanNetPremiKosong = "Please update Total Net Premium!"
)

// batasAngsuran - pagar jumlah putaran REPEAT.
const batasAngsuran = 1000

// BarisAngsuran - satu baris `InstallmentList` (Section/Installments.xml).
type BarisAngsuran struct {
	Installment    string `json:"Installment"`
	DueDate        string `json:"DueDate"`
	WPC            string `json:"WPC"`
	PaymentDate    string `json:"PaymentDate"`
	Currency       string `json:"Currency"`
	InstallmentPct string `json:"InstallmentPct"`
	Amount         string `json:"Amount"`
}

// Angsuran - satu halaman `TreatyIn.Installment` (satu mata uang).
type Angsuran struct {
	Currency        string          `json:"Currency"`
	AmountTotal     string          `json:"AmountTotal"`
	PctTotal        string          `json:"PctTotal"`
	InstallmentList []BarisAngsuran `json:"InstallmentList"`
}

// MasukanAngsuran - satu aksi tab Installment beserta isian layarnya.
type MasukanAngsuran struct {
	Aksi string `json:"aksi"`
	// Status - `update` (tombol Update Value), `editpercentage` (sel
	// % Installment), atau kosong.
	Status string `json:"status"`
	// InstallmentNo - `Param.Installment` = isian `TreatyIn.InstallmentNo`.
	InstallmentNo string `json:"installmentNo"`
	// EDMState - `TreatyIn.EDMState`; bukan ""/"0" = layar Adjustment.
	EDMState string `json:"edmState"`
	// Angsuran - `TreatyIn.Installment` saat ini.
	Angsuran []Angsuran `json:"angsuran"`
	// AngsuranLama - `TreatyIn.OLDDATA.Installment` (hanya Adjustment).
	AngsuranLama []Angsuran `json:"angsuranLama"`
	// NetPremium - `TreatyIn.TotalShareNetNP`.
	NetPremium []NilaiMataUang `json:"netPremium"`
	// Indeks - halaman Installment (mulai 0) untuk `total-baris`.
	Indeks int `json:"indeks"`
}

// HasilAngsuran - keadaan tab sesudah aksi.
type HasilAngsuran struct {
	Angsuran []Angsuran `json:"angsuran"`
	// TotalInstallmentNP - `nil` (JSON `null`) bila aksinya tidak
	// menyentuhnya (`total-baris`).
	TotalInstallmentNP []NilaiMataUang `json:"TotalInstallmentNP"`
	InstallmentNo      string          `json:"InstallmentNo"`
	Pesan              []string        `json:"pesan"`
}

// HitungAngsuran - bentuk ber-pelaku untuk handler. Tanggal hari ini dibaca
// di zona Jakarta - zona yang `@CurrentDate(..., "Asia/Jakarta")` sebut.
func (l *Layanan) HitungAngsuran(p inti.Pelaku, m MasukanAngsuran) (HasilAngsuran, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilAngsuran{}, err
	}
	return HitungAngsuran(m, penomor.DiJakarta(time.Now()))
}

// HitungAngsuran menjalankan satu aksi tab Installment. MURNI: nol baca, nol
// tulis basis data; `kini` = hari ini di Jakarta.
func HitungAngsuran(m MasukanAngsuran, kini time.Time) (HasilAngsuran, error) {
	switch m.Aksi {
	case AksiAngsuranNilai:
		return nilaiAngsuran(m, kini)
	case AksiAngsuranTotalBaris:
		h := HasilAngsuran{Angsuran: salinAngsuranSemua(m.Angsuran), InstallmentNo: m.InstallmentNo, Pesan: []string{}}
		if m.Indeks >= 0 && m.Indeks < len(h.Angsuran) {
			h.Angsuran[m.Indeks] = totalHalamanAngsuran(h.Angsuran[m.Indeks], m.Status, m.NetPremium)
		}
		return h, nil
	case AksiAngsuranTanggalBayar, AksiAngsuranTanggalBayarSemua:
		h := HasilAngsuran{Angsuran: salinAngsuranSemua(m.Angsuran), InstallmentNo: m.InstallmentNo, Pesan: []string{}}
		semua := m.Aksi == AksiAngsuranTanggalBayarSemua
		for i := range h.Angsuran {
			if !semua && i != m.Indeks {
				continue
			}
			for j := range h.Angsuran[i].InstallmentList {
				tanggalBayar(&h.Angsuran[i].InstallmentList[j])
			}
		}
		return h, nil
	case AksiAngsuranTotal:
		// TreatyInNPSetTotal langkah 27 (Property-Remove) + 28 (append).
		return HasilAngsuran{
			Angsuran:           salinAngsuranSemua(m.Angsuran),
			TotalInstallmentNP: totalAngsuranNP(m.Angsuran),
			InstallmentNo:      m.InstallmentNo,
			Pesan:              []string{},
		}, nil
	}
	return HasilAngsuran{}, fmt.Errorf("%w: aksi installment %q", ErrMasukanTidakSah, m.Aksi)
}

// tanggalBayar - `TreatyInUpdatePaymentDate` langkah 1.1 dan
// `TreatyInUpdatePaymentDate_Act` langkah 1.1.1 (diunggah pemakai 7 Oktober
// 2026; kedua Activity kelas `ASM-FW-GISFW-Data-TreatyInInstallment`):
//
//	local.adddate = .WPC
//	.PaymentDate  = @addCalendar(.DueDate, 0,0,0, local.adddate + 1, 0,0,0)
//
// Parameter keempat `@addCalendar` = HARI — urutan yang sama dengan
// `TreatyInSetReport` (`@addCalendar(TempDate,'0','0','0',hari,...)`).
// `DueDate`/`PaymentDate` tersimpan 8 aksara (`YYYYMMDD`, tipe tanggal, bukan
// stempel DateTime GMT), jadi tidak ada geseran zona waktu.
//
// ⚠️ Tidak berbukti, diputuskan:
//   - WPC kosong / bukan bilangan bulat = 0 (`local.adddate + 1` = 1 hari);
//   - Due Date kosong / tak terbaca: Payment Date DIBIARKAN. `_Act`
//     menjaganya dengan `.DueDate != ""`; versi sel tanpa penjaga dan hasil
//     `@addCalendar` atas tanggal kosong tidak terbaca dari ekspor.
//
// Belum diukur atas data: tabel pendaratan `T_TREATY_INSTALLMENT_ITEM` sedang
// dikosongkan.
func tanggalBayar(b *BarisAngsuran) {
	jatuh, ok := tanggalMasukan(b.DueDate)
	if !ok {
		return
	}
	hari := 0
	if w, err := strconv.Atoi(strings.TrimSpace(b.WPC)); err == nil {
		hari = w
	}
	b.PaymentDate = jatuh.AddDate(0, 0, hari+1).Format("20060102")
}

// nilaiAngsuran = `TreatyInSetValueInstallment` langkah 1-12.
func nilaiAngsuran(m MasukanAngsuran, kini time.Time) (HasilAngsuran, error) {
	// Langkah 2-3: tab dikosongkan SEBELUM pemeriksaan - pesan langkah 4/5
	// meninggalkan grid kosong, persis Pega.
	h := HasilAngsuran{
		Angsuran:           []Angsuran{},
		TotalInstallmentNP: []NilaiMataUang{},
		InstallmentNo:      m.InstallmentNo,
		Pesan:              []string{},
	}
	n, ok := angkaSah(m.InstallmentNo)
	if !ok || n.Cmp(apd.New(1, 0)) < 0 {
		h.Pesan = append(h.Pesan, pesanAngsuranKurang)
		return h, nil
	}
	if len(m.NetPremium) == 0 {
		h.Pesan = append(h.Pesan, pesanNetPremiKosong)
		return h, nil
	}
	if n.Cmp(apd.New(batasAngsuran, 0)) > 0 {
		return HasilAngsuran{}, fmt.Errorf("%w: lebih dari %d angsuran", ErrMasukanTidakSah, batasAngsuran)
	}

	// Langkah 6.
	for _, np := range m.NetPremium {
		h.Angsuran = append(h.Angsuran, Angsuran{Currency: np.Currency, InstallmentList: []BarisAngsuran{}})
	}
	// Langkah 7.
	persen := bagiBulatDes(seratus, n, 2)
	tanggal := kini.Format("20060102")
	// Langkah 8.
	for i := range h.Angsuran {
		a := &h.Angsuran[i]
		nilai := apd.New(0, 0)
		for _, np := range m.NetPremium {
			if np.Currency == a.Currency {
				nilai = angka(np.Value)
			}
		}
		ke, jumlah := apd.New(0, 0), apd.New(0, 0)
		// REPEAT Start 1, Iteration 1, Limit TreatyIn.InstallmentNo.
		for putaran := int64(1); apd.New(putaran, 0).Cmp(n) <= 0; putaran++ {
			ke = tambah(ke, apd.New(1, 0))
			jumlah = tambah(jumlah, persen)
			pctBaris := persen // 8.2.3: n < N || pct == 100
			if ke.Cmp(n) == 0 && jumlah.Cmp(seratus) != 0 {
				pctBaris = tambah(persen, kurang(seratus, jumlah)) // 8.2.2
			}
			a.InstallmentList = append(a.InstallmentList, BarisAngsuran{
				Installment:    teks(ke),
				DueDate:        tanggal,
				Currency:       a.Currency,
				InstallmentPct: teks(pctBaris),
				Amount:         teks(kali(bagiBulat(pctBaris, 100, 4), nilai)),
			})
		}
	}
	// Langkah 9 lalu 10 - urutannya penting: OLDDATA menimpa TempDate.
	if m.Status == statusAngsuranUpdate {
		pulihkanAngsuran(h.Angsuran, m.Angsuran)
	}
	if m.EDMState != "" && m.EDMState != "0" {
		pulihkanAngsuran(h.Angsuran, m.AngsuranLama)
	}
	// Langkah 11 (tanpa parameter = hanya total) dan 12.
	for i := range h.Angsuran {
		h.Angsuran[i] = totalHalamanAngsuran(h.Angsuran[i], "", m.NetPremium)
	}
	h.TotalInstallmentNP = totalAngsuranNP(h.Angsuran)
	return h, nil
}

// pulihkanAngsuran - langkah 9/10: empat medan dari halaman sumber berindeks
// sama. Halaman sumber yang tidak ada: nilai hitungan dipertahankan (lihat
// kepala berkas).
func pulihkanAngsuran(tuju, sumber []Angsuran) {
	for i := range tuju {
		if i >= len(sumber) {
			continue
		}
		for j := range tuju[i].InstallmentList {
			if j >= len(sumber[i].InstallmentList) {
				continue
			}
			s := sumber[i].InstallmentList[j]
			r := &tuju[i].InstallmentList[j]
			r.DueDate, r.WPC, r.PaymentDate, r.InstallmentPct = s.DueDate, s.WPC, s.PaymentDate, s.InstallmentPct
		}
	}
}

// totalHalamanAngsuran = `SetTotalInstallment` atas satu halaman Installment
// (`Param.treaty` bukan "out").
//
//	status=editpercentage: langkah 1 (tanpa metode) lalu 3 - per baris
//	  premivalue = TotalShareNetNP.Value bermata uang sama (BERTAHAN antarbaris
//	  bila tidak ada yang cocok), Amount = @divide(InstallmentPct,100,20) * premivalue
//	selain itu: lompat ke label `calc` (langkah 4)
//	langkah 4-5: AmountTotal = Σ Amount, PctTotal = Σ InstallmentPct
func totalHalamanAngsuran(a Angsuran, status string, net []NilaiMataUang) Angsuran {
	b := salinAngsuran(a)
	if status == statusEditPersen {
		premi := apd.New(0, 0)
		for i := range b.InstallmentList {
			r := &b.InstallmentList[i]
			for _, np := range net {
				if np.Currency == r.Currency {
					premi = angka(np.Value)
				}
			}
			r.Amount = teks(kali(bagiBulat(angka(r.InstallmentPct), 100, 20), premi))
		}
	}
	nilai, pct := apd.New(0, 0), apd.New(0, 0)
	for _, r := range b.InstallmentList {
		nilai = tambah(nilai, angka(r.Amount))
		pct = tambah(pct, angka(r.InstallmentPct))
	}
	b.AmountTotal, b.PctTotal = teks(nilai), teks(pct)
	return b
}

// totalAngsuranNP - `TotalInstallmentNP` = {Currency, Value: AmountTotal}
// per halaman Installment (SetValueInstallment 12 = NPSetTotal 28).
func totalAngsuranNP(xs []Angsuran) []NilaiMataUang {
	out := []NilaiMataUang{}
	for _, a := range xs {
		out = append(out, NilaiMataUang{Currency: a.Currency, Value: a.AmountTotal})
	}
	return out
}

// angkaSah membaca teks desimal; kosong / bukan angka → false.
func angkaSah(s string) (*apd.Decimal, bool) {
	d, _, err := apd.NewFromString(strings.TrimSpace(s))
	if err != nil || d.Form != apd.Finite {
		return nil, false
	}
	return d, true
}

func salinAngsuran(a Angsuran) Angsuran {
	b := a
	b.InstallmentList = append([]BarisAngsuran{}, a.InstallmentList...)
	return b
}

func salinAngsuranSemua(xs []Angsuran) []Angsuran {
	out := make([]Angsuran, 0, len(xs))
	for _, a := range xs {
		out = append(out, salinAngsuran(a))
	}
	return out
}
