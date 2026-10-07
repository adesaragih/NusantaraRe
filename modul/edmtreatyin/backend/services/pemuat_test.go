package services_test

// Uji pemuat dokumen lama endorsemen (tiket EDM 10, 09) di atas gudang tiruan modul (`backend/tiruan`, yang
// menegakkan UNIQUE OLD_POLIS_ID, UNIQUE (NOPOLIS, PRODKE), generasi tertutup, transaksi batal, dan proyeksi 'PEGA'
// beku) ditambah pembaca JSON_POLIS dan penulis tambahan pemuat. ⭐ Uji utama tiket 10: rantai TIGA generasi dimuat
// dari dokumen lama, penjaga percabangan dan keutuhan benar-benar berjalan atasnya - bukan dilewati karena ini jalur
// migrasi. ⛔ Fixture fiktif UJI- (nomor polis, nama); nomor kasus EDMT-99xxxx fiktif.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/repository"
	"nusantarare/modul/edmtreatyin/backend/services"
	"nusantarare/modul/edmtreatyin/backend/tiruan"
)

// gudangUji - tiruan modul + JSON_POLIS + penulis tambahan pemuat (`repository/lama_edm.go`).
type gudangUji struct {
	*tiruan.Gudang
	dok     map[string]models.BarisJSONPolis   // ROWID -> baris JSON_POLIS endorsemen
	nb      map[string][]models.BarisJSONPolis // NOPOLIS -> baris JSON_POLIS generasi NB
	penanda map[string]models.PenandaMigrasi
	datar   map[string]models.KolomDatarLama
}

func gudangBaru() *gudangUji {
	return &gudangUji{Gudang: tiruan.Baru(), dok: map[string]models.BarisJSONPolis{}, nb: map[string][]models.BarisJSONPolis{},
		penanda: map[string]models.PenandaMigrasi{}, datar: map[string]models.KolomDatarLama{}}
}

func salinPeta[K comparable, V any](m map[K]V) map[K]V {
	out := make(map[K]V, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func (g *gudangUji) Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error {
	penanda, datar := salinPeta(g.penanda), salinPeta(g.datar)
	err := g.Gudang.Transaksi(ctx, fn)
	if err != nil {
		g.penanda, g.datar = penanda, datar
	}
	return err
}

// KunciJSONPolisEDM - urutan SQL (NOPOLIS, ROWID), BUKAN urutan generasi: pemuat wajib mengurutkannya sendiri.
func (g *gudangUji) KunciJSONPolisEDM(ctx context.Context) ([]models.KunciJSONPolis, error) {
	var out []models.KunciJSONPolis
	for k, b := range g.dok {
		out = append(out, models.KunciJSONPolis{Kunci: k, NoPolis: b.NoPolis, ProdKe: b.ProdKe})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].NoPolis != out[j].NoPolis {
			return out[i].NoPolis < out[j].NoPolis
		}
		return out[i].Kunci < out[j].Kunci
	})
	return out, nil
}

func (g *gudangUji) BacaJSONPolisEDM(ctx context.Context, kunci string) (models.BarisJSONPolis, error) {
	b, ada := g.dok[kunci]
	if !ada {
		return b, fmt.Errorf("uji: ROWID %s tidak ada", kunci)
	}
	return b, nil
}

func (g *gudangUji) BacaJSONPolisNB(ctx context.Context, nopolis string) ([]models.BarisJSONPolis, error) {
	return g.nb[nopolis], nil
}

func (g *gudangUji) KunciGenerasiLama(ctx context.Context, tx *db.Tx, id string) (models.KunciGenerasi, bool, error) {
	x := g.Generasi[id]
	if x == nil {
		return models.KunciGenerasi{}, false, nil
	}
	return models.KunciGenerasi{NoPolis: x.NoPolis, ProdKe: x.ProdKe, NoEndors: x.EDMNo}, true, nil
}

// GenerasiSebelumnya = repository.GenerasiSebelumnya.
func (g *gudangUji) GenerasiSebelumnya(ctx context.Context, tx *db.Tx, k models.Kasus, nopolis string) (string, error) {
	if k.OldPolisID != "" {
		return k.OldPolisID, nil
	}
	for id, x := range g.Generasi {
		if nopolis != "" && x.NoPolis == nopolis && x.ProdKe == k.ProdKe-1 {
			return id, nil
		}
	}
	return "", nil
}

func (g *gudangUji) SetelKolomDatarLamaEDM(ctx context.Context, tx *db.Tx, id string, k models.KolomDatarLama) error {
	x := g.Generasi[id]
	if x == nil || x.ProdKe < 1 {
		return errors.New("uji: kolom datar menyentuh 0 baris")
	}
	for _, y := range g.Generasi {
		if y.OldPolisID == id {
			return repository.ErrGenerasiTertutup
		}
	}
	g.datar[id] = k
	return nil
}

// SetelPenandaMigrasi = repository.SetelPenandaMigrasi: hanya induk SUMBER 'PEGA', setiap baris tepat satu kali.
func (g *gudangUji) SetelPenandaMigrasi(ctx context.Context, tx *db.Tx, polisID string, p models.PenandaMigrasi) error {
	s := g.Selisih[polisID]
	if s == nil || s.Kunci.Sumber != models.SumberPega {
		return errors.New("uji: penanda menyentuh 0 baris (bukan SUMBER PEGA)")
	}
	if len(p.Spreading) != len(s.Halaman.AmbilDaftar(models.DaftarSelisihSpreading)) ||
		len(p.Angsuran) != len(s.Halaman.AmbilDaftar(models.DaftarSelisihAngsuran)) ||
		len(p.Lapisan) != len(models.DatarSelisihLapisan(s.Halaman)) {
		return errors.New("uji: cacah penanda berbeda dari baris selisih")
	}
	g.penanda[polisID] = p
	return nil
}

func (g *gudangUji) SalinUsulanLamaEDM(ctx context.Context, tx *db.Tx, idPega string, baris []models.UsulanProduksi) (models.NasibUsulan, error) {
	if len(baris) == 0 {
		return models.UsulanTanpaBaris, nil
	}
	for _, u := range g.Usulan {
		if u.IDPega == idPega {
			return models.UsulanDilewati, nil
		}
	}
	return models.UsulanDisalin, g.CatatUsulan(ctx, tx, idPega, baris)
}

var _ services.GudangPemuat = (*gudangUji)(nil)

// ------------------------------------------------------------------ fixture

const (
	polisA = "UJI-QP.T1.10.2017.00001" // rantai tiga generasi + satu percabangan
	polisB = "UJI-QP.T1.10.2017.00002" // generasi pertama kehilangan baris spreading
	polisC = "UJI-QP.T1.10.2017.00003" // generasi NB belum dimuat
)

// dokUji - satu dokumen json_polis (struktur data guide; isi fiktif).
type dokUji struct {
	id, nopolis      string
	prodke           int
	oldEDMNo         string
	netPremium       string
	spreading        []string // TreatyType urut baris
	angsuran         []string // InstallmentNo urut baris
	selisihNet       string   // TreatyDifference.NetPremium (APA ADANYA)
	selisihSpreading []string // TreatyDifference.SpreadingRiskList().PremiumSpreaded
}

func (d dokUji) edmNo() string { return models.NomorEDM(d.nopolis, d.prodke) }

func (d dokUji) isi() map[string]any {
	var sp, sel, an, selAn []any
	for i, tt := range d.spreading {
		sp = append(sp, map[string]any{"TreatyType": tt, "SharePercentage": "50", "PremiumSpreaded": strconv.Itoa(10 * (i + 1))})
		v := "0"
		if i < len(d.selisihSpreading) {
			v = d.selisihSpreading[i]
		}
		sel = append(sel, map[string]any{"TreatyType": tt, "SharePercentage": "50", "PremiumSpreaded": v})
	}
	for _, n := range d.angsuran {
		an = append(an, map[string]any{"InstallmentNo": n, "DueDate": "20171101", "Premium": "5"})
		selAn = append(selAn, map[string]any{"InstallmentNo": n, "DueDate": "20171101", "Premium": "-0.000000276"})
	}
	m := map[string]any{
		"pxObjClass": models.KelasDokumenTreatyIn, "PolicyNo": d.nopolis, "NetPremium": d.netPremium,
		"QuotationData":     map[string]any{"ProportionalType": models.JenisProporsional, "OldPolicyNo": d.nopolis, "BusinessFac": "T"},
		"SpreadingRiskList": sp, "ListInstallment": an,
	}
	if d.prodke > 0 {
		m["EDMNo"], m["ProdKe"], m["EDMType"] = d.edmNo(), strconv.Itoa(d.prodke), "3"
		m["OldData"] = map[string]any{"EDMNo": d.oldEDMNo, "PolicyNo": d.nopolis}
		m["TreatyDifference"] = map[string]any{"NetPremium": d.selisihNet, "SpreadingRiskList": sel, "ListInstallment": selAn}
		m["SuggestList"] = []any{map[string]any{"Date": "20171002T010000.000 GMT", "Suggest": "UJI catatan " + d.id,
			"IsApproved": "1", "OperatorName": "UJI-OPERATOR"}}
	}
	return m
}

func (d dokUji) baris(t *testing.T) models.BarisJSONPolis {
	t.Helper()
	data, err := json.Marshal(d.isi())
	if err != nil {
		t.Fatal(err)
	}
	b := models.BarisJSONPolis{NoPolis: d.nopolis, ProdKe: strconv.Itoa(d.prodke), TglInput: "2017-10-02 08:00:00",
		Username: "UJI-AKUN", DataJSON: data}
	if d.prodke > 0 {
		b.IDPega, b.NoEndors = models.KelasKerjaEDM+" "+d.id, d.edmNo()
	} else {
		b.IDPega = "ASM-FW-GISFW-WORK-NB " + d.id
	}
	return b
}

// Generasi NB (PRODKE 0) - dimuat pemuat NB; di sini ditanam ke tiruan dan ke JSON_POLIS (pembanding uji-kering).
var nbUji = []dokUji{
	{id: "NB-990000", nopolis: polisA, netPremium: "150", spreading: []string{"UJI-QS", "UJI-SP"}, angsuran: []string{"1"}},
	{id: "NB-990010", nopolis: polisB, netPremium: "1", spreading: []string{"UJI-QS", "UJI-SP"}, angsuran: []string{"1"}},
}

// rowid -> dokumen endorsemen. ROWID sengaja tidak urut generasi.
var edmUji = map[string]dokUji{
	"A3": {id: "EDMT-990001", nopolis: polisA, prodke: 1, netPremium: "150", spreading: []string{"UJI-QS", "UJI-SP"},
		angsuran: []string{"1"}, selisihNet: "-130463146.760000276", selisihSpreading: []string{"-0.000000276", "12.5"}},
	// generasi 2: baris spreading bertukar tempat (pasangan bergeser); varian rumus berlapis Pega - 180 - 50 = 130.
	"A1": {id: "EDMT-990002", nopolis: polisA, prodke: 2, oldEDMNo: models.NomorEDM(polisA, 1), netPremium: "180",
		spreading: []string{"UJI-SP", "UJI-QS"}, angsuran: []string{"1", "2"}, selisihNet: "130",
		selisihSpreading: []string{"7.000000001", "-7"}},
	// percabangan: generasi 2 kedua atas generasi 1 yang sama.
	"A9": {id: "EDMT-990009", nopolis: polisA, prodke: 2, oldEDMNo: models.NomorEDM(polisA, 1), netPremium: "170",
		spreading: []string{"UJI-QS", "UJI-SP"}, angsuran: []string{"1"}, selisihNet: "20"},
	"A2": {id: "EDMT-990003", nopolis: polisA, prodke: 3, oldEDMNo: models.NomorEDM(polisA, 2), netPremium: "0",
		spreading: []string{"UJI-SP", "UJI-QS", "UJI-XL"}, angsuran: []string{"1", "2"}, selisihNet: "-180"},
	// keutuhan: generasi 1 polis B kehilangan baris spreading ke-2 milik generasi NB.
	"B1": {id: "EDMT-990011", nopolis: polisB, prodke: 1, netPremium: "1", spreading: []string{"UJI-QS"}, angsuran: []string{"1"},
		selisihNet: "0"},
	"B2": {id: "EDMT-990012", nopolis: polisB, prodke: 2, oldEDMNo: models.NomorEDM(polisB, 1), netPremium: "1",
		spreading: []string{"UJI-QS", "UJI-SP"}, angsuran: []string{"1"}, selisihNet: "0"},
	"C1": {id: "EDMT-990021", nopolis: polisC, prodke: 1, netPremium: "1", spreading: []string{"UJI-QS"}, angsuran: []string{"1"},
		selisihNet: "0"},
}

func siapkan(t *testing.T) *gudangUji {
	t.Helper()
	g := gudangBaru()
	for _, d := range nbUji {
		b := d.baris(t)
		g.nb[d.nopolis] = append(g.nb[d.nopolis], b)
		h, err := models.HalamanPembanding(b)
		if err != nil {
			t.Fatal(err)
		}
		g.TanamPolis(d.id, d.nopolis, 0, "", h, "")
	}
	for rowid, d := range edmUji {
		g.dok[rowid] = d.baris(t)
	}
	return g
}

func jalankanPemuat(t *testing.T, g *gudangUji, tulis bool) models.RingkasanPemuat {
	t.Helper()
	var arsip, gal bytes.Buffer
	lap, err := models.LaporanPemuatBaru(&arsip, &gal, tulis)
	if err != nil {
		t.Fatal(err)
	}
	if err := services.PemuatBaru(g).Jalankan(context.Background(), tulis, lap); err != nil {
		t.Fatal(err)
	}
	if err := lap.Tutup(); err != nil {
		t.Fatal(err)
	}
	return lap.Ringkasan()
}

// periksaGalatRantai - jenis galat ringkasan yang diharapkan rantai uji: percabangan (A9), keutuhan (B1), generasi
// sebelumnya belum dimuat (B2 - karena B1 gagal; C1 - generasi NB tidak ada).
func periksaGalatRantai(t *testing.T, r models.RingkasanPemuat) {
	t.Helper()
	harap := map[string]int{
		models.JenisGalat(models.ErrPercabangan):                1,
		models.JenisGalat(models.ErrKeutuhan):                   1,
		models.JenisGalat(models.ErrGenerasiSebelumnyaTidakAda): 2,
	}
	if !reflect.DeepEqual(r.GalatPerJenis, harap) || r.DokumenGagal != 4 {
		t.Errorf("galat %v (dokumen gagal %d), harap %v", r.GalatPerJenis, r.DokumenGagal, harap)
	}
}

// ------------------------------------------------------------------ uji

// TestPemuatRantaiTigaGenerasi - tiket 10 uji utama + AC 39-44.
func TestPemuatRantaiTigaGenerasi(t *testing.T) {
	g := siapkan(t)
	r := jalankanPemuat(t, g, true)
	if r.Dimuat != 3 || r.SudahDimuat != 0 || r.MedanBelumDiputuskan != 0 {
		t.Fatalf("ringkasan %+v", r)
	}
	periksaGalatRantai(t, r)

	// Rantai generasi: OLD_POLIS_ID = generasi TEPAT sebelumnya (ID-9, ID-14), PRODKE dan NOENDORS dari dokumen.
	sebelumnya := map[string]string{"EDMT-990001": "NB-990000", "EDMT-990002": "EDMT-990001", "EDMT-990003": "EDMT-990002"}
	for id, lama := range sebelumnya {
		x := g.Generasi[id]
		if x == nil {
			t.Fatalf("%s tidak dimuat", id)
		}
		n, _ := strconv.Atoi(id[len(id)-1:])
		if x.OldPolisID != lama || x.NoPolis != polisA || x.ProdKe != n || x.EDMNo != models.NomorEDM(polisA, n) || x.EDMType != "3" {
			t.Errorf("%s: generasi %+v", id, *x)
		}
		if k := g.Kasus[id]; k.StatusWork != models.StatusSelesai {
			t.Errorf("%s: kasus %+v - json_polis endorsemen hanya lahir sesudah disetujui", id, k)
		}
		if d := g.datar[id]; d != (models.KolomDatarLama{TglInput: "2017-10-02 08:00:00", Username: "UJI-AKUN"}) {
			t.Errorf("%s: kolom datar %+v", id, d)
		}
		s := g.Selisih[id]
		if s == nil || s.Kunci != (repository.KunciSelisih{NoPolis: polisA, ProdKe: n, EDMNo: models.NomorEDM(polisA, n),
			IDPega: id, Sumber: models.SumberPega}) {
			t.Errorf("%s: kunci selisih %+v", id, s)
		}
	}
	// Penjaga BENAR-BENAR berjalan: percabangan ditolak SisipKasus (UNIQUE OLD_POLIS_ID), keutuhan ditolak sebelum
	// menulis; dokumen yang ditolak tidak meninggalkan satu baris pun (satu transaksi).
	sisip := 0
	for _, p := range g.Panggil {
		if p == "SisipKasus" {
			sisip++
		}
	}
	if sisip != 4 {
		t.Errorf("SisipKasus dipanggil %d kali, harap 4 (tiga generasi + percabangan yang ditolak basis data)", sisip)
	}
	for _, id := range []string{"EDMT-990009", "EDMT-990011", "EDMT-990012", "EDMT-990021"} {
		if g.Generasi[id] != nil || g.Kasus[id].ID != "" || g.Selisih[id] != nil {
			t.Errorf("%s ditolak tetapi meninggalkan baris", id)
		}
		if _, ada := g.penanda[id]; ada {
			t.Errorf("%s ditolak tetapi penandanya tertulis", id)
		}
	}

	// ⛔ AC 39: selisih tersimpan PERSIS seperti dokumen - digit galat utuh, varian rumus berlapis TIDAK dihitung ulang.
	for id, harap := range map[string]string{"EDMT-990001": "-130463146.760000276", "EDMT-990002": "130", "EDMT-990003": "-180"} {
		if v := g.Selisih[id].Halaman.Ambil(models.HalamanPolis + ".TreatyDifference.NetPremium"); v != harap {
			t.Errorf("%s: NET_PREMIUM selisih %q, harap %q (dokumen)", id, v, harap)
		}
	}
	sp := g.Selisih["EDMT-990001"].Halaman.AmbilDaftar(models.DaftarSelisihSpreading)
	if len(sp) != 2 || sp[0]["PremiumSpreaded"] != "-0.000000276" || sp[1]["PremiumSpreaded"] != "12.5" {
		t.Errorf("selisih spreading generasi 1 %+v", sp)
	}
	if an := g.Selisih["EDMT-990002"].Halaman.AmbilDaftar(models.DaftarSelisihAngsuran); len(an) != 2 ||
		an[1]["Premium"] != "-0.000000276" || an[1]["DueDate"] != "2017-11-01" {
		t.Errorf("selisih angsuran generasi 2 %+v", an)
	}
	h2, err := models.PecahDokumenEDM(g.dok["A1"])
	if err != nil {
		t.Fatal(err)
	}
	g1, err := g.BacaGenerasi(context.Background(), nil, "EDMT-990001")
	if err != nil {
		t.Fatal(err)
	}
	hitung := h2.Halaman.Salin()
	models.PasangOldData(hitung, g1)
	if err := models.EDMTCalculateTreatyDifference(hitung); err != nil {
		t.Fatal(err)
	}
	if v := hitung.Ambil(models.HalamanPolis + ".TreatyDifference.NetPremium"); v == "130" || v == "" {
		t.Errorf("fixture tidak membedakan: hitung ulang %q", v)
	}

	// Penanda (AC 40-43): generasi 1 lawan NB - berpasangan, bukan berlapis; generasi 2 - spreading bertukar
	// (bergeser), sebelumnya ber-EDMNo (berlapis); generasi 3 - baris tambahan tidak bergeser.
	f, t1 := models.PenandaBaris{}, models.PenandaBaris{RumusBerlapis: true}
	gb := models.PenandaBaris{PasanganBergeser: true, RumusBerlapis: true}
	harapPenanda := map[string]models.PenandaMigrasi{
		"EDMT-990001": {Spreading: []models.PenandaBaris{f, f}, Angsuran: []models.PenandaBaris{f}},
		"EDMT-990002": {Spreading: []models.PenandaBaris{gb, gb}, Angsuran: []models.PenandaBaris{t1, t1}},
		"EDMT-990003": {Spreading: []models.PenandaBaris{t1, t1, t1}, Angsuran: []models.PenandaBaris{t1, t1}},
	}
	if !reflect.DeepEqual(g.penanda, harapPenanda) {
		t.Errorf("penanda %+v\nharap   %+v", g.penanda, harapPenanda)
	}
	for id := range g.penanda {
		if g.Selisih[id].Kunci.Sumber != models.SumberPega {
			t.Errorf("%s: penanda pada baris bukan PEGA (AC 43)", id)
		}
	}
	if r.Penanda != (models.RingkasanPenanda{BarisSpreading: 7, BarisAngsuran: 5, BergeserSpreading: 2, Berlapis: 9}) {
		t.Errorf("ringkasan penanda %+v", r.Penanda)
	}
	// ID-38: proyeksi PEGA beku - jalur biasa (SUMBER 'GO') tidak dapat membangunnya ulang.
	err = g.Transaksi(context.Background(), func(tx *db.Tx) error {
		return g.SimpanSelisih(context.Background(), tx, "EDMT-990002", h2.Halaman, repository.KunciSelisih{Sumber: models.SumberGo})
	})
	if !errors.Is(err, repository.ErrSelisihBeku) {
		t.Errorf("bangun ulang GO atas baris PEGA: %v", err)
	}
	// ... dan PEGA atas PEGA pun ditolak (bukan hapus anak lalu bentrok UNIQUE POLIS_ID).
	err = g.Transaksi(context.Background(), func(tx *db.Tx) error {
		return g.SimpanSelisih(context.Background(), tx, "EDMT-990002", h2.Halaman, repository.KunciSelisih{Sumber: models.SumberPega})
	})
	if !errors.Is(err, repository.ErrSelisihBeku) {
		t.Errorf("tulis ulang PEGA atas baris PEGA: %v", err)
	}

	// SuggestList -> riwayat produksi lewat CatatUsulan, kunci = ID kasus.
	if len(g.Usulan) != 3 || g.Usulan[0].TypePolis != "EDMT" || g.Usulan[0].AksesLogin != "" {
		t.Errorf("usulan %+v", g.Usulan)
	}
	if r.AngkaKasusTerbesar != 990003 || r.GenerasiTerbesar != 3 {
		t.Errorf("nomor kasus terbesar %d / generasi %d", r.AngkaKasusTerbesar, r.GenerasiTerbesar)
	}

	// Jalankan ulang aman: generasi yang sudah dimuat dilewati, tidak ada yang ganda.
	jmlUsulan, jmlGen := len(g.Usulan), len(g.Generasi)
	r2 := jalankanPemuat(t, g, true)
	if r2.Dimuat != 0 || r2.SudahDimuat != 3 || len(g.Usulan) != jmlUsulan || len(g.Generasi) != jmlGen {
		t.Errorf("jalankan ulang %+v (usulan %d, generasi %d)", r2, len(g.Usulan), len(g.Generasi))
	}
	periksaGalatRantai(t, r2)
}

// TestPemuatUjiKeringNolTulisan - uji-kering: hanya JSON_POLIS yang dibaca; galat dan penanda yang sama dilaporkan.
func TestPemuatUjiKeringNolTulisan(t *testing.T) {
	g := siapkan(t)
	jmlGen := len(g.Generasi)
	r := jalankanPemuat(t, g, false)
	if r.Dimuat != 3 || r.Tulis {
		t.Fatalf("ringkasan %+v", r)
	}
	periksaGalatRantai(t, r)
	if len(g.Panggil) != 0 || len(g.Generasi) != jmlGen || len(g.Selisih) != 0 || len(g.Usulan) != 0 ||
		len(g.penanda) != 0 || len(g.datar) != 0 {
		t.Errorf("uji-kering menulis: panggil %v, generasi %d, selisih %d", g.Panggil, len(g.Generasi), len(g.Selisih))
	}
	if r.Penanda != (models.RingkasanPenanda{BarisSpreading: 7, BarisAngsuran: 5, BergeserSpreading: 2, Berlapis: 9}) {
		t.Errorf("penanda uji-kering %+v - harus sama dengan -jalankan", r.Penanda)
	}
}

// TestPemuatIDKasusBentrokTidakMenimpa - ID kasus yang sudah dipakai generasi lain bukan "sudah dimuat".
func TestPemuatIDKasusBentrokTidakMenimpa(t *testing.T) {
	g := siapkan(t)
	g.TanamPolis("EDMT-990001", "UJI-QP.LAIN", 1, "", models.HalamanBaru(), "")
	r := jalankanPemuat(t, g, true)
	if r.GalatPerJenis[models.JenisGalat(models.ErrIDKasusBentrok)] != 1 || r.SudahDimuat != 0 {
		t.Errorf("ringkasan %+v", r)
	}
	if x := g.Generasi["EDMT-990001"]; x.NoPolis != "UJI-QP.LAIN" {
		t.Errorf("generasi lain tertimpa: %+v", *x)
	}
}
