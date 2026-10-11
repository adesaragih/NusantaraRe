package services

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// cargoCsvTiruan - PenulisCargoCsv tanpa basis data: ganti = cargo case dikosongkan dulu; gagalKe = panggilan ke-n gagal (sekali).
type cargoCsvTiruan struct {
	ada     map[string][]models.ObjekCargo
	ganti   []bool
	ukuran  []int
	gagalKe int
}

func (c *cargoCsvTiruan) TulisCargoCsv(_ context.Context, _ *db.Tx, id string, cargo []models.ObjekCargo, ganti bool) error {
	c.ganti, c.ukuran = append(c.ganti, ganti), append(c.ukuran, len(cargo))
	if c.gagalKe == len(c.ganti) {
		c.gagalKe = 0
		return errors.New("ORA-UJI koneksi putus")
	}
	if c.ada == nil {
		c.ada = map[string][]models.ObjekCargo{}
	}
	if ganti {
		c.ada[id] = nil
	}
	c.ada[id] = append(c.ada[id], cargo...)
	return nil
}

// kapalTiruan - PenyimpanKapalMaster tanpa basis data; daftar = nama kapal yang didaftarkan (urut).
type kapalTiruan struct {
	id     map[string]bool
	nama   map[string]models.KapalCargo
	daftar []string
}

func (k *kapalTiruan) AdaKapalID(_ context.Context, id string) (bool, error) { return k.id[id], nil }

func (k *kapalTiruan) KapalDariNama(_ context.Context, nama string) (models.KapalCargo, bool, error) {
	x, ok := k.nama[nama]
	return x, ok, nil
}

func (k *kapalTiruan) DaftarKapal(_ context.Context, _ *db.Tx, nama, operator string) (models.KapalCargo, error) {
	if k.nama == nil {
		k.nama = map[string]models.KapalCargo{}
	}
	k.daftar = append(k.daftar, nama+"/"+operator)
	x := models.KapalCargo{ID: fmt.Sprintf("10000000%02d", len(k.daftar)), NmShip: nama, Remark: "Incomplete"}
	k.nama[nama] = x
	return x, nil
}

// masterCargoCsvUji - masterUji + agen A1.
func masterCargoCsvUji() masterCargoTiruan {
	m := maps.Clone(masterUji())
	m["Agen/A1"] = "UJI AGEN"
	return m
}

// layananUnggahCargoUji - layananUnggahUji untuk case MARINE CARGO (GroupPanel 004): mata uang IDR / CNY, kapal "UJI LAMA" ada di
// master (ID S1), ID kapal "S9" ada.
func layananUnggahCargoUji() (*Service, *unggahTiruan, *cargoCsvTiruan, *kapalTiruan, *csvTiruan) {
	svc, u, _, _, c := layananUnggahUji()
	svc.tableOfLimit = tolTiruan{bisnis: satuBisnis("004", "24")}
	ct := &cargoCsvTiruan{}
	kt := &kapalTiruan{id: map[string]bool{"S9": true}, nama: map[string]models.KapalCargo{"UJI LAMA": {ID: "S1", NmShip: "UJI LAMA",
		GRT: "1.5", YMake1: "2001", FlagName: "ID"}}}
	svc = svc.DenganObjekCargo(cargoTiruan{ada: map[string][]models.ObjekCargo{"UJI-NB-1": nil}, mu: map[string]models.MataUangCsv{
		"IDR": {ID: "10026", OldID: "1", Currency: "IDR"}, "CNY": {ID: "10050", OldID: "9", Currency: "CNY"}}}).
		DenganMasterCargo(masterCargoCsvUji()).DenganUnggahCsvCargo(ct, kt)
	return svc, u, ct, kt, c
}

// barisCargoUji - satu baris cargo sah (per nama kolom).
func barisCargoUji() map[string]string {
	return map[string]string{"TRADING": "DOMESTIK", "DEPARTURE": "jakarta", "DESTINATION": "surabaya", "GOODS": "BARANG",
		"PACKING": "PETI", "CONVEYENCE": "KAPAL LAUT", "TSI": "1.000,5", "CURRENCY": "IDR", "AGENT": "A1", "DATASHIP": "uji lama",
		"INVOICE": "INV-1", "ETD": "05/03/2026", "BLNUMBER": "BL-1", "TOP_RISK": "TRUE"}
}

// judulCargoUji - judul cargo urutan TERBALIK + satu kolom tak dikenal (NO) di depan.
func judulCargoUji() []string {
	j := slices.Clone(models.KolomCsvCargo)
	slices.Reverse(j)
	return append([]string{"NO"}, j...)
}

// csvCargoUji - berkas CSV cargo n baris (judulCargoUji); ubah(i, baris) mengubah baris ke-i (1..n).
func csvCargoUji(n int, ubah func(int, map[string]string)) string {
	judul := judulCargoUji()
	var b strings.Builder
	b.WriteString("\uFEFF" + strings.Join(judul, ";") + "\r\n")
	for i := 1; i <= n; i++ {
		r := barisCargoUji()
		r["NO"] = fmt.Sprint(i)
		if ubah != nil {
			ubah(i, r)
		}
		sel := make([]string, len(judul))
		for k, j := range judul {
			sel[k] = r[j]
		}
		b.WriteString(strings.Join(sel, ";") + "\r\n")
	}
	return b.String()
}

// TestPetaJudulCargo - judul per NAMA: urutan bebas, tak peka huruf / spasi tepi, kolom tak dikenal diabaikan, ganda / kurang ->
// 400 berpesan; peta FIRE tetap posisi.
func TestPetaJudulCargo(t *testing.T) {
	p, err := petaJudulCargo(strings.Join(judulCargoUji(), ";"))
	if err != nil || p.n != 18 || p.idx[0] != -1 || p.kolom[0] != "TOP_RISK" || len(p.kolom) != 18 || p.idx[1] != 17 {
		t.Fatalf("%v %+v", err, p)
	}
	isi, ada := selDataPeta("9;TRUE;;;;;;;;;;;;;;;;;DOMESTIK", p)
	if !ada || isi[0] != "DOMESTIK" || isi[17] != "TRUE" {
		t.Errorf("sel: %v", isi)
	}
	if _, ada := selDataPeta("x;;", p); !ada {
		t.Error("baris hanya berisi kolom tak dikenal dianggap kosong")
	}
	if _, err := petaJudulCargo(" trading ;" + strings.Join(models.KolomCsvCargo, ";")); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), "TRADING dua kali") {
		t.Errorf("ganda: %v", err)
	}
	if _, err := petaJudulCargo("TRADING;GOODS"); !errors.Is(err, ErrMasukanObjek) || !strings.Contains(err.Error(), "DEPARTURE, DESTINATION") {
		t.Errorf("kurang: %v", err)
	}
	if _, err := petaJudulFire(strings.Join(models.KolomCsvCargo, ";")); !errors.Is(err, ErrMasukanObjek) {
		t.Error("judul cargo diterima sebagai FIRE")
	}
}

// isiCargoUji - baris per nama -> isi kanonik.
func isiCargoUji(r map[string]string) []string {
	isi := make([]string, len(models.KolomCsvCargo))
	for i, k := range models.KolomCsvCargo {
		isi[i] = r[k]
	}
	return isi
}

// TestUraiBarisCargoCsv - langkah 3.1: TSI berkoma (titik ribuan dibuang) / bertitik, RMB -> CNY, DEPARTURE / DESTINATION / DATASHIP
// huruf besar, TOP_RISK "TRUE" persis, ETD dd?MM?yyyy -> SailDate bentuk tanggal kita; angka / tanggal / lebar tidak sah = galat.
func TestUraiBarisCargoCsv(t *testing.T) {
	r := barisCargoUji()
	r["CURRENCY"] = "RMB"
	b, m := uraiBarisCargoCsv(4, isiCargoUji(r))
	if len(m) != 0 || utils.FormatDecimal(b.tsi) != "1000.5" || b.mataUang != "CNY" || b.departure != "JAKARTA" || b.kapal != "UJI LAMA" ||
		!b.topRisk || b.sailDate != tanggalKerugianKePega("05-03-2026") || b.kunci() != [4]string{"DOMESTIK", "BARANG", "KAPAL LAUT", "PETI"} {
		t.Fatalf("%v %+v", m, b)
	}
	r["TSI"], r["TOP_RISK"], r["ETD"] = "2500.75", "true", "05-03-2026"
	if b, m := uraiBarisCargoCsv(4, isiCargoUji(r)); len(m) != 0 || utils.FormatDecimal(b.tsi) != "2500.75" || b.topRisk ||
		b.sailDate == "" {
		t.Errorf("titik / true kecil: %v %+v", m, b)
	}
	r["TSI"], r["ETD"], r["INVOICE"] = "abc", "5/3/2026", strings.Repeat("X", 51)
	_, m = uraiBarisCargoCsv(4, isiCargoUji(r))
	if len(m) != 3 || !strings.Contains(strings.Join(m, "|"), "baris 4 kolom TSI: angka tidak sah") ||
		!strings.Contains(strings.Join(m, "|"), "kolom ETD: tanggal harus dd/MM/yyyy") || !strings.Contains(m[0], "INVOICE: paling banyak 50") {
		t.Errorf("galat: %q", m)
	}
	r = barisCargoUji()
	r["TSI"], r["ETD"] = "", ""
	if b, m := uraiBarisCargoCsv(1, isiCargoUji(r)); len(m) != 0 || b.tsi != nil || b.sailDate != "" {
		t.Errorf("kosong: %v %+v", m, b)
	}
}

// TestProtectBarisCargoCsv - ProtectUploadMarineCargo_act: Trading / Goods / Conveyence tidak ada = berhenti; PACKING = lanjut; ID
// kapal di DATASHIP (atau kosong) MENGGANTI pesan; agen tidak ada = kosong tanpa pesan; mata uang tidak ada = galat.
func TestProtectBarisCargoCsv(t *testing.T) {
	svc, _, _, _, _ := layananUnggahCargoUji()
	ctx := context.Background()
	cek := func(ubah func(map[string]string)) ([]string, barisCargoCsv) {
		r := barisCargoUji()
		ubah(r)
		b, m := uraiBarisCargoCsv(1, isiCargoUji(r))
		if len(m) > 0 {
			t.Fatal(m)
		}
		m, err := svc.periksaBarisCargoCsv(ctx, &b, cacheCargoCsvBaru())
		if err != nil {
			t.Fatal(err)
		}
		return m, b
	}
	if m, b := cek(func(map[string]string) {}); len(m) != 0 || b.tradingID != "T1" || b.packingID != "P1" ||
		b.agenSurvei.FullName != "UJI AGEN" || b.mu.ID != "10026" {
		t.Errorf("sah: %q %+v", m, b)
	}
	for nama, c := range map[string]struct {
		ubah func(map[string]string)
		mau  []string
	}{
		"trading": {func(r map[string]string) { r["TRADING"], r["GOODS"] = "X", "Y" }, []string{"* Trading X Not Available in Master"}},
		"goods":   {func(r map[string]string) { r["GOODS"], r["PACKING"] = "Y", "Z" }, []string{"* Goods Y Not Available in Master"}},
		"packing+conveyence": {func(r map[string]string) { r["PACKING"], r["CONVEYENCE"] = "Z", "W" },
			[]string{"* PACKING Z Not Available in Master", "* Conveyence W Not Available in Master"}},
		"packing+id kapal": {func(r map[string]string) { r["PACKING"], r["DATASHIP"] = "Z", "s9" }, []string{pesanKapalIDCsv}},
		"datship kosong":   {func(r map[string]string) { r["DATASHIP"] = "" }, []string{pesanKapalIDCsv}},
		"trading kosong":   {func(r map[string]string) { r["TRADING"] = "" }, []string{"* Trading  Not Available in Master"}},
		"mata uang":        {func(r map[string]string) { r["CURRENCY"] = "XYZ" }, []string{`baris 1 kolom CURRENCY: mata uang "XYZ" tidak ada di master CURRENCY`}},
	} {
		if m, _ := cek(c.ubah); !reflect.DeepEqual(m, c.mau) {
			t.Errorf("%s: %q", nama, m)
		}
	}
	if m, b := cek(func(r map[string]string) { r["AGENT"] = "A404" }); len(m) != 0 || b.agenSurvei != (models.AgenSurvei{}) {
		t.Errorf("agen tak ada: %q %+v", m, b.agenSurvei)
	}
}

// TestUnggahCsvCargoAlur - 2 500 baris, pasangan baris berurutan berkunci sama (cargo = baris ganjil + genap, 1 250 cargo):
// validasi -> simpan per kelompok >= 1000 baris yang TIDAK memotong cargo; cargo gabungan = kunci / invoice / ETD baris pertama,
// kapal / rute / agen / coverage baris terakhir; kapal baru didaftarkan sekali per nama; CurrencyList dihapus sekali; View Upload
// memakai judul urut berkas.
func TestUnggahCsvCargoAlur(t *testing.T) {
	ctx := context.Background()
	svc, _, ct, kt, c := layananUnggahCargoUji()
	// Pasangan (1,2), (3,4), … berkunci sama; pasangan berurutan dibedakan TRADING (DOMESTIK / EKSPOR bergantian). Baris genap
	// berganti kapal (baru, didaftarkan), DEPARTURE, dan mata uang (RMB -> CNY).
	svc.masterCargo.(masterCargoTiruan)["Trading/T2"] = "EKSPOR"
	berkas := csvCargoUji(2500, func(i int, r map[string]string) {
		r["INVOICE"], r["TSI"] = fmt.Sprint("INV-", i), fmt.Sprint(i)
		if (i+1)/2%2 == 0 {
			r["TRADING"] = "EKSPOR"
		}
		if i%2 == 0 {
			r["DATASHIP"], r["DEPARTURE"], r["CURRENCY"] = "kapal baru", "medan", "RMB"
		}
	})
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "cargo.csv", strings.NewReader(berkas))
	if err != nil {
		t.Fatal(err)
	}
	x := st.Unggah
	if x.Status != StatusUnggahValid || x.Jenis != models.JenisCsvCargo || x.TotalBaris != 2500 || x.JumlahGalat != 0 ||
		!slices.Equal(x.Kolom, judulCargoUji()[1:]) {
		t.Fatalf("validasi: %+v %v", x, st.Galat)
	}
	h, err := svc.BarisUnggahCsv(ctx, "UJI-NB-1", x.ID, 1, 2)
	if err != nil || h.Jenis != models.JenisCsvCargo || h.Kolom[0] != "TOP_RISK" || h.Baris[1].Isi[0] != "DOMESTIK" || h.Baris[1].Isi[1] != "medan" {
		t.Fatalf("baris: %v %+v", err, h)
	}
	st, err = svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", x.ID)
	if err != nil || st.Unggah.Status != StatusUnggahSelesai || st.Unggah.BarisTersimpan != 2500 {
		t.Fatalf("simpan: %v %+v", err, st.Unggah)
	}
	cargo := ct.ada["UJI-NB-1"]
	if len(cargo) != 1250 || !slices.Equal(ct.ganti, []bool{true, false, false}) || !slices.Equal(ct.ukuran, []int{500, 500, 250}) ||
		len(c.hapus) != 1 || !slices.Equal(kt.daftar, []string{"KAPAL BARU/" + pelakuCsv.AkunID}) {
		t.Fatalf("cargo %d ganti %v ukuran %v hapus %d daftar %v", len(cargo), ct.ganti, ct.ukuran, len(c.hapus), kt.daftar)
	}
	g := cargo[1] // baris 3 + 4
	if g.TradingNote != "EKSPOR" || g.TradingID != "T2" || g.InvoiceNumber != "INV-3" || g.Ship.NmShip != "KAPAL BARU" ||
		g.FromRute != "MEDAN" || g.ToRute != "SURABAYA" || len(g.Coverages) != 1 || utils.FormatDecimal(g.Coverages[0].TSI) != "4" ||
		g.Coverages[0].Currency != "CNY" || g.Coverages[0].CurrencyID != "10050" || g.Ship.AdditionalShip[0].NmShip != "KAPAL BARU" ||
		g.SurveyAgent.AgentID != "A1" || !g.IsTopRisk || g.Coverages[0].Coverage != "" {
		t.Errorf("cargo gabungan: %+v", g)
	}
	if cargo[0].Ship.ID != "1000000001" || cargo[0].InvoiceNumber != "INV-1" {
		t.Errorf("cargo 1: %+v", cargo[0])
	}
}

// TestUnggahCsvCargoGagalDanUlangi - galat Protect -> gagal_validasi (semua baris dilaporkan, simpan ditolak); putus di kelompok 2
// -> gagal + Ulangi melanjutkan tanpa cargo ganda; selama menyimpan, PUT objek-cargo 409.
func TestUnggahCsvCargoGagalDanUlangi(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, _ := layananUnggahCargoUji()
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "c.csv", strings.NewReader(csvCargoUji(5, func(i int, r map[string]string) {
		if i == 2 || i == 4 {
			r["GOODS"] = "TIDAK ADA"
		}
	})))
	if err != nil || st.Unggah.Status != StatusUnggahGagalValidasi || st.Unggah.JumlahGalat != 2 ||
		st.Galat[0] != (models.PesanUnggahCsv{NoBaris: 2, Teks: "* Goods TIDAK ADA Not Available in Master"}) {
		t.Fatalf("%v %+v %v", err, st.Unggah, st.Galat)
	}
	if _, err := svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", st.Unggah.ID); !errors.Is(err, ErrUnggahStatus) {
		t.Errorf("simpan gagal_validasi: %v", err)
	}

	svc, u, ct, _, _ := layananUnggahCargoUji()
	ct.gagalKe = 2
	// Setiap baris cargo sendiri (TRADING bergantian DOMESTIK / EKSPOR).
	svc.masterCargo.(masterCargoTiruan)["Trading/T2"] = "EKSPOR"
	st, err = svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "c.csv", strings.NewReader(csvCargoUji(2100, func(i int, r map[string]string) {
		if i%2 == 0 {
			r["TRADING"] = "EKSPOR"
		}
	})))
	if err != nil || st.Unggah.Status != StatusUnggahValid {
		t.Fatalf("validasi: %v %+v %v", err, st.Unggah, st.Galat)
	}
	id := st.Unggah.ID
	st, _ = svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", id)
	if st.Unggah.Status != StatusUnggahGagal || st.Unggah.BarisTersimpan != 1000 || len(ct.ada["UJI-NB-1"]) != 1000 ||
		!strings.Contains(st.Unggah.GalatProses, "gagal di baris 1001") {
		t.Fatalf("putus: %+v (%d cargo)", st.Unggah, len(ct.ada["UJI-NB-1"]))
	}
	st, err = svc.UlangiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", id)
	if err != nil || st.Unggah.Status != StatusUnggahSelesai || len(ct.ada["UJI-NB-1"]) != 2100 ||
		!slices.Equal(ct.ganti, []bool{true, false, false, false}) {
		t.Fatalf("ulangi: %v %+v (%d cargo, %v)", err, st.Unggah, len(ct.ada["UJI-NB-1"]), ct.ganti)
	}

	u.kerja[id].Status = StatusUnggahMenyimpan
	if _, err := svc.GantiObjekCargo(ctx, pelakuCsv, "UJI-NB-1", nil); !errors.Is(err, ErrUnggahBerjalan) {
		t.Errorf("kunci objek-cargo: %v", err)
	}
}

// TestUnggahCsvCargoTanpaMaster - master Marine belum terpasang (nama tabel DEV menunggu) -> 503 ErrMasterCargoBelumAda sebelum
// berkas dibaca; case FIRE tetap memakai importer FIRE.
func TestUnggahCsvCargoTanpaMaster(t *testing.T) {
	ctx := context.Background()
	svc, u, _, _, _ := layananUnggahCargoUji()
	svc.masterCargo = nil
	if _, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "c.csv", strings.NewReader(csvCargoUji(1, nil))); !errors.Is(err, ErrMasterCargoBelumAda) {
		t.Errorf("tanpa master: %v", err)
	}
	if len(u.kerja) != 0 {
		t.Error("pekerjaan dibuat")
	}
	fire, _, _, _, _ := layananUnggahUji()
	if _, err := fire.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "c.csv", strings.NewReader(csvCargoUji(1, nil))); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), "template objek") {
		t.Errorf("berkas cargo di case FIRE: %v", err)
	}
}

// TestKelompokCargoTidakMemotongCargo - cargo 3 baris (TRADING bergantian per 3 baris): baris 1000-1002 satu cargo melintasi batas
// 1000 baris -> kelompok 1 = baris 1..1002 (334 cargo), kelompok 2 = 1003..1200 (66 cargo); tidak ada cargo terbelah.
func TestKelompokCargoTidakMemotongCargo(t *testing.T) {
	ctx := context.Background()
	svc, _, ct, _, _ := layananUnggahCargoUji()
	svc.masterCargo.(masterCargoTiruan)["Trading/T2"] = "EKSPOR"
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "c.csv", strings.NewReader(csvCargoUji(1200, func(i int, r map[string]string) {
		if (i-1)/3%2 == 1 {
			r["TRADING"] = "EKSPOR"
		}
	})))
	if err != nil || st.Unggah.Status != StatusUnggahValid {
		t.Fatalf("%v %+v", err, st.Unggah)
	}
	st, err = svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", st.Unggah.ID)
	if err != nil || st.Unggah.Status != StatusUnggahSelesai || len(ct.ada["UJI-NB-1"]) != 400 || !slices.Equal(ct.ukuran, []int{334, 66}) {
		t.Errorf("%v %+v %d %v", err, st.Unggah, len(ct.ada["UJI-NB-1"]), ct.ukuran)
	}
}

// TestJudulTemplatePegaCargo - K74-18 (keputusan work owner 11-10-2026): judul persis template Pega (17 wajib + INVOICEDATE, tanpa
// COVERAGE_NAME) diterima; INVOICEDATE diabaikan (tidak di kolom); TSI "1.237.675,00" berkutip terurai; tanpa salah satu dari 17
// -> 400.
func TestJudulTemplatePegaCargo(t *testing.T) {
	judul := "TRADING;ETD;DEPARTURE;DESTINATION;INVOICE;INVOICEDATE;BLNUMBER;CONVEYENCE;CONVEYENCE_NAME;DATASHIP;GOODS;GOODS_NAME;" +
		"PACKING;AGENT;COVERAGE_IDX;CURRENCY;TSI;TOP_RISK"
	p, err := petaJudulCargo(judul)
	if err != nil || len(p.kolom) != 17 || slices.Contains(p.kolom, "INVOICEDATE") || p.kolom[1] != "ETD" {
		t.Fatalf("%v %+v", err, p.kolom)
	}
	isi, _ := selDataPeta(`DOMESTIK;21/09/2021;Rengat;Medan;INV;21/09/2021;BL;KAPAL LAUT;;UJI LAMA;BARANG;;PETI;;1;IDR;"1.237.675,00";`, p)
	b, m := uraiBarisCargoCsv(1, isi)
	if len(m) != 0 || utils.FormatDecimal(b.tsi) != "1237675" || b.departure != "RENGAT" || b.topRisk {
		t.Errorf("%v %+v", m, b)
	}
	if _, err := petaJudulCargo(strings.Replace(judul, ";TOP_RISK", "", 1)); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), "kolom TOP_RISK tidak ada") {
		t.Errorf("kurang: %v", err)
	}
}

// TestAgenBelumAdaCargo - AGENT terisi selagi master agen belum terpasang: validasi BERHENTI di baris itu (gagal_validasi, galat
// jelas), tidak ditebak kosong; AGENT kosong tetap sah (K74-20).
func TestAgenBelumAdaCargo(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, _ := layananUnggahCargoUji()
	svc.masterCargo.(masterCargoTiruan)["Agen/*"] = "BELUM"
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "c.csv", strings.NewReader(csvCargoUji(1500, func(i int, r map[string]string) {
		if i < 3 {
			r["AGENT"] = ""
		}
	})))
	if err != nil || st.Unggah.Status != StatusUnggahGagalValidasi || st.Unggah.JumlahGalat != 1 || st.Unggah.BarisDiproses != 3 ||
		st.Galat[0].NoBaris != 3 || !strings.Contains(st.Galat[0].Teks, "baris 3 kolom AGENT \"A1\": master agen (LLOYDAGENT) belum tersedia") {
		t.Fatalf("%v %+v %v", err, st.Unggah, st.Galat)
	}
	st, err = svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "c.csv", strings.NewReader(csvCargoUji(2, func(_ int, r map[string]string) {
		r["AGENT"] = ""
	})))
	if err != nil || st.Unggah.Status != StatusUnggahValid {
		t.Errorf("tanpa agen: %v %+v %v", err, st.Unggah, st.Galat)
	}
}

// TestDaftarMasterCargo - GET /marine/{jenis}: empat master terpasang; coverage / max-age / agent (DDL menunggu) dan master belum
// terpasang -> ErrMasterCargoBelumAda (503); catatan MaxAge dengan master MaxAge belum ada: ID sama -> catatan lama, ID baru -> 503.
func TestDaftarMasterCargo(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, _ := layananUnggahCargoUji()
	d, err := svc.DaftarMasterCargo(ctx, "trading")
	if err != nil || len(d) != 1 || d[0] != (models.BarisMasterCargo{ID: "T1", Description: "DOMESTIK"}) {
		t.Errorf("%v %+v", err, d)
	}
	for _, j := range []string{"coverage", "max-age", "agent", "x"} {
		if _, err := svc.DaftarMasterCargo(ctx, j); !errors.Is(err, ErrMasterCargoBelumAda) {
			t.Errorf("%s: %v", j, err)
		}
	}
	if _, err := Baru(nil).DaftarMasterCargo(ctx, "goods"); !errors.Is(err, ErrMasterCargoBelumAda) {
		t.Errorf("tanpa master: %v", err)
	}
	cari := func() (string, bool, error) { return "", false, repository.ErrMasterBelumAda }
	if isi, err := svc.catatanMaster("1", "M1", true, "W LAMA", cari, "MaxAge"); err != nil || isi != "W LAMA" {
		t.Errorf("sama: %v %q", err, isi)
	}
	if _, err := svc.catatanMaster("1", "M2", false, "", cari, "MaxAge"); !errors.Is(err, ErrMasterCargoBelumAda) {
		t.Errorf("baru: %v", err)
	}
}

// TestLebarMasterCargoCsv - audit lebar 11-10-2026: nilai MASTER melebihi kolom tujuan (BYTE = len UTF-8, bukan jumlah rune) ->
// galat jelas per baris; tepat di batas sah. ID 50 (T_CARGOLIST.*_ID), alamat / kota / negara agen 500 (T_SURVEYAGENT).
func TestLebarMasterCargoCsv(t *testing.T) {
	sah := barisCargoCsv{no: 7, tradingID: strings.Repeat("T", 50), goodID: strings.Repeat("é", 25), packingID: "P1",
		conveyanceID: "C1", agenSurvei: models.AgenSurvei{Address: strings.Repeat("A", 500), City: strings.Repeat("é", 250), Country: "ID"}}
	if m := lebarMasterCargoCsv(sah); len(m) != 0 {
		t.Fatalf("tepat di batas: %q", m)
	}
	for nama, c := range map[string]struct {
		ubah func(*barisCargoCsv)
		mau  string
	}{
		"trading 51":             {func(b *barisCargoCsv) { b.tradingID += "T" }, "baris 7: master ID Trading 51 byte melebihi kolom T_CARGOLIST.TRADING_ID (paling banyak 50 byte)"},
		"goods 26 rune 52 bita":  {func(b *barisCargoCsv) { b.goodID += "é" }, "master ID Goods 52 byte melebihi kolom T_CARGOLIST.GOOD_ID"},
		"packing 51":             {func(b *barisCargoCsv) { b.packingID = strings.Repeat("P", 51) }, "T_CARGOLIST.PACKING_ID"},
		"conveyance 51":          {func(b *barisCargoCsv) { b.conveyanceID = strings.Repeat("C", 51) }, "T_CARGOLIST.CONVEYANCE_ID"},
		"alamat 501":             {func(b *barisCargoCsv) { b.agenSurvei.Address += "A" }, "master Address agen 501 byte melebihi kolom T_SURVEYAGENT.ADDRESS (paling banyak 500 byte)"},
		"kota 251 rune 502 bita": {func(b *barisCargoCsv) { b.agenSurvei.City += "é" }, "master City agen 502 byte melebihi kolom T_SURVEYAGENT.CITY"},
		"negara 501":             {func(b *barisCargoCsv) { b.agenSurvei.Country = strings.Repeat("N", 501) }, "T_SURVEYAGENT.COUNTRY"},
	} {
		b := sah
		c.ubah(&b)
		if m := lebarMasterCargoCsv(b); len(m) != 1 || !strings.Contains(m[0], c.mau) {
			t.Errorf("%s: %q", nama, m)
		}
	}
	// lewat periksaBarisCargoCsv: ID master Trading 51 byte -> galat baris.
	svc, _, _, _, _ := layananUnggahCargoUji()
	svc.masterCargo.(masterCargoTiruan)["Trading/"+strings.Repeat("X", 51)] = "PANJANG"
	r := barisCargoUji()
	r["TRADING"] = "PANJANG"
	b, _ := uraiBarisCargoCsv(3, isiCargoUji(r))
	m, err := svc.periksaBarisCargoCsv(context.Background(), &b, cacheCargoCsvBaru())
	if err != nil || len(m) != 1 || !strings.Contains(m[0], "baris 3: master ID Trading 51 byte") {
		t.Errorf("periksa: %v %q", err, m)
	}
}
