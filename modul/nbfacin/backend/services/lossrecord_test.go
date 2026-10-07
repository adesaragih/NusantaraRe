package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

type lossRecordTiruan struct {
	ada   models.LossRecord
	tulis []models.LossRecord
}

func (l *lossRecordTiruan) BacaLossRecord(context.Context, string) (models.LossRecord, error) {
	return l.ada, nil
}

func (l *lossRecordTiruan) TulisLossRecord(_ context.Context, _ *db.Tx, _ string, x models.LossRecord) error {
	l.tulis = append(l.tulis, x)
	return nil
}

type jabatanTiruan map[string]string

func (j jabatanTiruan) JabatanAkun(_ context.Context, akun string) (string, error) {
	if akun == "UJI-GALAT" {
		return "", errors.New("UJI galat baca")
	}
	return j[akun], nil
}

var (
	akunIT    = inti.Pelaku{AkunID: "UJI-IT"}
	akunBukan = inti.Pelaku{AkunID: "UJI-UW"}
)

// layananLossRecord - Begin 01-01-2025 WIB s.d. 01-01-2026 WIB, RI Commission 25 %, premi spread IDR 1.000 + USD 0,5
// (kurs 16.000) -> Premi PAR bawaan 9.000.
func layananLossRecord(lt *lossRecordTiruan) *Service {
	kasus := kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": {CaseID: "UJI-NB-1", General: models.General{
		StartDateTime: "20241231T170000.000 GMT", EndDateTime: "20251231T170000.000 GMT"}}}}
	cov := func(premi string) []models.CoverageSpreading {
		return []models.CoverageSpreading{{Spreading: []models.BarisSpreading{{TreatyType: "UJI-QS", PremiumSpreaded: ds(premi)}}}}
	}
	sp := &spreadingTiruan{baca: []models.KasusSpreading{{StartDateTime: "20241231T170000.000 GMT", Lokasi: []models.LokasiSpreading{
		{ObjectNo: "1", Items: []models.ItemSpreading{{Currency: "IDR", Coverages: cov("1000")}, {Currency: "USD", Coverages: cov("0.5")}}}}}},
		idUang: map[string]string{"IDR": idMataUangIDR, "USD": "UJI-USD", "EUR": "UJI-EUR"}, kurs: map[string]*apd.Decimal{"UJI-USD": ds("16000")}}
	return Baru(nil).DenganLossRecord(lt, jabatanTiruan{"UJI-IT": " it ", "UJI-UW": "Underwriter"}).DenganKasus(kasus).
		DenganSpreading(sp).DenganDeduction(&deductionTiruan{ada: models.KasusDeduction{RICommision: ds("25")}}).DenganTransaksi(tanpaTx).
		DenganKelasBisnis(kelasBisnisTiruan{})
}

// TestJendelaUwYear - Get5YearAgo: ADD_MONTHS −12L (hari terakhir bulan -> hari terakhir bulan), Uw Year = tahun awal.
func TestJendelaUwYear(t *testing.T) {
	tg := func(s string) time.Time { v, _ := time.Parse(BentukTanggalKabel, s); return v }
	j := jendelaUw(tg("01-01-2025"), tg("01-01-2026"))
	if len(j) != 5 || j[0].uwYear != "2024" || !j[0].mulai.Equal(tg("01-01-2024")) || !j[0].akhir.Equal(tg("01-01-2025")) ||
		j[4].uwYear != "2020" || j[4].ke != 5 {
		t.Errorf("%+v", j)
	}
	if got := tambahBulanOracle(tg("28-02-2025"), -12); !got.Equal(tg("29-02-2024")) {
		t.Errorf("hari terakhir Februari: %s", got.Format(BentukTanggalKabel))
	}
	if got := tambahBulanOracle(tg("15-03-2025"), -12); !got.Equal(tg("15-03-2024")) {
		t.Errorf("hari biasa: %s", got.Format(BentukTanggalKabel))
	}
}

// TestHitungLossRecord - Uw Year / TahunKe (jendela terakhir menang), DOL bawaan, huruf besar, kurs (USD ada, EUR tidak),
// total per COB + currency, Premi PAR bawaan, rasio 1 / 3 / 5 tahun dibetulkan, urutan baris dipertahankan.
func TestHitungLossRecord(t *testing.T) {
	lt := &lossRecordTiruan{}
	svc := layananLossRecord(lt)
	m := models.LossRecord{PremiMb: ds("1000"), Baris: []models.BarisLossRecord{
		{ClassOfBusiness: "UJI KEBAKARAN", DateOfLoss: "15-06-2024", CauseOfLoss: "uji api", Currency: "idr", LossAmount: ds("100")},
		{ClassOfBusiness: "UJI KEBAKARAN", DateOfLoss: "15-06-2022", CauseOfLoss: "uji", Currency: "USD", LossAmount: ds("1")},
		{ClassOfBusiness: "UJI LAIN", DateOfLoss: "", CauseOfLoss: "uji", Currency: " eur ", LossAmount: ds("5")},
		{ClassOfBusiness: "UJI KEBAKARAN", DateOfLoss: "01-01-2015", CauseOfLoss: "uji", Currency: "IDR", LossAmount: ds("7")},
	}}
	tp, err := svc.HitungLossRecord(context.Background(), akunIT, "UJI-NB-1", m)
	if err != nil {
		t.Fatal(err)
	}
	b := tp.LossRecord.Baris
	type ringkas struct{ uw, ke, dol, sebab, mu, idr string }
	got := []ringkas{}
	for _, x := range b {
		got = append(got, ringkas{x.UwYear, string(rune('0' + x.TahunKe)), x.DateOfLoss, x.CauseOfLoss, x.Currency, teksUji(x.LossAmountIDR)})
	}
	mau := []ringkas{{"2024", "1", "15-06-2024", "UJI API", "IDR", "100"}, {"2022", "3", "15-06-2022", "UJI", "USD", "16000"},
		{"2023", "2", "01-01-2024", "UJI", "EUR", ""}, {"", "0", "01-01-2015", "UJI", "IDR", "7"}}
	if !reflect.DeepEqual(got, mau) {
		t.Errorf("baris:\n%v\nmau\n%v", got, mau)
	}
	if teksUji(tp.LossRecord.PremiPar) != "9000" || !reflect.DeepEqual(tp.Pesan,
		[]string{"Kurs EUR pada Begin date tidak ada di TREATYEXCHANGEYEARLY - Loss Amount (IDR) tidak diisi"}) {
		t.Errorf("premi PAR %s / pesan %v", teksUji(tp.LossRecord.PremiPar), tp.Pesan)
	}
	r := tp.LossRecord.Rasio // Gross 10.000; Y1 100; Y3 = Y5 = 16.100; net faktor 0,75
	if teksUji(r.Gross1) != "1" || teksUji(r.Gross3) != "53.66666667" || teksUji(r.Gross5) != "32.2" ||
		teksUji(r.Net1) != "1.33333333" || teksUji(r.Net3) != "71.55555556" || teksUji(r.Net5) != "42.93333333" {
		t.Errorf("rasio %+v", []string{teksUji(r.Gross1), teksUji(r.Gross3), teksUji(r.Gross5), teksUji(r.Net1), teksUji(r.Net3), teksUji(r.Net5)})
	}
	tot := []string{}
	for _, c := range tp.TotalCob {
		tot = append(tot, c.ClassOfBusiness+"|"+c.Currency+"|"+teksUji(c.LossAmount)+"|"+teksUji(c.LossAmountIDR))
	}
	if !reflect.DeepEqual(tot, []string{"UJI KEBAKARAN|IDR|107|107", "UJI KEBAKARAN|USD|1|16000", "UJI LAIN|EUR|5|"}) {
		t.Errorf("total %v", tot)
	}
	if len(lt.tulis) != 0 {
		t.Error("hitung tidak boleh menyimpan")
	}
	// Premi PAR terisi dipakai apa adanya; Gross 0 -> rasio kosong + pesan.
	tp, _ = svc.HitungLossRecord(context.Background(), akunIT, "UJI-NB-1", models.LossRecord{PremiPar: ds("5")})
	if teksUji(tp.LossRecord.PremiPar) != "5" || teksUji(tp.LossRecord.Rasio.Gross1) != "0" {
		t.Errorf("premi terisi: %+v", tp.LossRecord)
	}
	svc.DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{{}}})
	tp, _ = svc.HitungLossRecord(context.Background(), akunIT, "UJI-NB-1", models.LossRecord{})
	if tp.LossRecord.Rasio.Gross1 != nil || !reflect.DeepEqual(tp.Pesan, []string{"Gross premi (PAR + MB + PL + Other) 0 - loss ratio tidak dihitung"}) {
		t.Errorf("gross 0: %+v %v", tp.LossRecord.Rasio, tp.Pesan)
	}
}

func teksUji(d *apd.Decimal) string {
	if d == nil {
		return ""
	}
	return d.Text('f')
}

// TestSimpanLossRecord - PUT: urut DOL terbaru dulu (tanpa DOL di akhir), DOL tersimpan teks Pega; GET mengembalikan DOL
// DD-MM-YYYY dan menghitung ulang.
func TestSimpanLossRecord(t *testing.T) {
	lt := &lossRecordTiruan{}
	svc := layananLossRecord(lt)
	m := models.LossRecord{PremiPar: ds("100"), Baris: []models.BarisLossRecord{
		{ClassOfBusiness: "UJI KEBAKARAN", DateOfLoss: "01-03-2024", LossAmount: ds("1")},
		{ClassOfBusiness: "UJI LAIN", DateOfLoss: "01-05-2024", LossAmount: ds("2")},
		{ClassOfBusiness: "UJI KEBAKARAN", DateOfLoss: "01-04-2024", LossAmount: ds("3")}}}
	tp, err := svc.SimpanLossRecord(context.Background(), akunIT, "UJI-NB-1", m)
	if err != nil || len(lt.tulis) != 1 {
		t.Fatalf("%v %d", err, len(lt.tulis))
	}
	urut := []string{}
	for i, b := range tp.LossRecord.Baris {
		urut = append(urut, b.DateOfLoss+"/"+lt.tulis[0].Baris[i].DateOfLoss)
	}
	if !reflect.DeepEqual(urut, []string{"01-05-2024/20240501T050000.000 GMT", "01-04-2024/20240401T050000.000 GMT",
		"01-03-2024/20240301T050000.000 GMT"}) {
		t.Errorf("urut %v", urut)
	}
	lt.ada = lt.tulis[0]
	g, err := svc.TampilanLossRecordKasus(context.Background(), akunIT, "UJI-NB-1")
	if err != nil || g.LossRecord.Baris[0].DateOfLoss != "01-05-2024" || g.LossRecord.Baris[0].UwYear != "2024" {
		t.Errorf("GET: %v %+v", err, g.LossRecord.Baris)
	}
}

// TestAksesLossRecord - JOB_POSITION "IT" (dipangkas, tanpa beda huruf); jalur DITOLAK: tanpa identitas 401, jabatan lain
// 403, tanpa pembaca 403; AksesTabIT lunak (galat -> false).
func TestAksesLossRecord(t *testing.T) {
	ctx := context.Background()
	svc := layananLossRecord(&lossRecordTiruan{})
	if !svc.AksesTabIT(ctx, akunIT) || svc.AksesTabIT(ctx, akunBukan) || svc.AksesTabIT(ctx, inti.Pelaku{}) ||
		svc.AksesTabIT(ctx, inti.Pelaku{AkunID: "UJI-GALAT"}) || Baru(nil).AksesTabIT(ctx, akunIT) {
		t.Error("AksesTabIT")
	}
	if _, err := svc.TampilanLossRecordKasus(ctx, inti.Pelaku{}, "UJI-NB-1"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	for _, p := range []inti.Pelaku{akunBukan, {AkunID: "UJI-TAK-ADA"}} {
		if _, err := svc.HitungLossRecord(ctx, p, "UJI-NB-1", models.LossRecord{}); !errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("403 %s: %v", p.AkunID, err)
		}
		if _, err := svc.SimpanLossRecord(ctx, p, "UJI-NB-1", models.LossRecord{}); !errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("403 simpan %s: %v", p.AkunID, err)
		}
	}
	if _, err := svc.HitungLossRecord(ctx, inti.Pelaku{AkunID: "UJI-GALAT"}, "UJI-NB-1", models.LossRecord{}); err == nil ||
		errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("galat baca jabatan harus diteruskan: %v", err)
	}
	if _, err := svc.HitungLossRecord(ctx, akunIT, "UJI-NB-1", models.LossRecord{Baris: []models.BarisLossRecord{{DateOfLoss: "2024-01-01"}}}); !errors.Is(err, ErrMasukanLossRecord) ||
		!strings.Contains(err.Error(), "baris[0].dateOfLoss") {
		t.Errorf("400: %v", err)
	}
}

// TestSimpanLossRecordKelasBisnis - Save: COB wajib dari daftar BUSINESS.NOTE (dipangkas, tanpa beda huruf), disimpan
// berejaan tabel; kosong / tidak ada -> 400 pesan per baris (nomor mulai 1), tanpa simpan. Hitung TIDAK memeriksanya.
func TestSimpanLossRecordKelasBisnis(t *testing.T) {
	ctx := context.Background()
	lt := &lossRecordTiruan{}
	svc := layananLossRecord(lt)
	baris := func(cob ...string) models.LossRecord {
		m := models.LossRecord{PremiPar: ds("100")}
		for _, c := range cob {
			m.Baris = append(m.Baris, models.BarisLossRecord{ClassOfBusiness: c, DateOfLoss: "01-03-2024", LossAmount: ds("1")})
		}
		return m
	}
	tp, err := svc.SimpanLossRecord(ctx, akunIT, "UJI-NB-1", baris("  uji kebakaran ", "UJI LAIN"))
	if err != nil || len(lt.tulis) != 1 || lt.tulis[0].Baris[0].ClassOfBusiness != "Uji Kebakaran" ||
		lt.tulis[0].Baris[1].ClassOfBusiness != "Uji Lain" || tp.TotalCob[0].ClassOfBusiness != "Uji Kebakaran" {
		t.Fatalf("kanon: %v %+v", err, lt.tulis)
	}
	_, err = svc.SimpanLossRecord(ctx, akunIT, "UJI-NB-1", baris("UJI LAIN", "", "UJI TAK ADA"))
	if !errors.Is(err, ErrMasukanLossRecord) || len(lt.tulis) != 1 ||
		!strings.Contains(err.Error(), "baris 2: Class Of Business wajib diisi; baris 3: Class Of Business \"UJI TAK ADA\" tidak ada di daftar BUSINESS") {
		t.Errorf("400: %v", err)
	}
	if tp, err := svc.HitungLossRecord(ctx, akunIT, "UJI-NB-1", baris("", "UJI TAK ADA")); err != nil || tp.LossRecord.Baris[1].ClassOfBusiness != "UJI TAK ADA" {
		t.Errorf("hitung tidak memeriksa COB: %v", err)
	}
	if _, err := svc.DenganKelasBisnis(nil).SimpanLossRecord(ctx, akunIT, "UJI-NB-1", baris("UJI LAIN")); !errors.Is(err, ErrKelasBisnisTanpaDatabase) {
		t.Errorf("503: %v", err)
	}
}
