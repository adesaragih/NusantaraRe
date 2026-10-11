package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// unggahTiruan - PenyimpanUnggahCsv dalam memori.
type unggahTiruan struct {
	n      int
	kerja  map[string]*models.UnggahCsv
	baris  map[string][]models.BarisUnggahCsv
	pesan  map[string]map[string][]models.PesanUnggahCsv
	adaObj func(string) bool
}

func barUnggahTiruan(adaObj func(string) bool) *unggahTiruan {
	return &unggahTiruan{kerja: map[string]*models.UnggahCsv{}, baris: map[string][]models.BarisUnggahCsv{},
		pesan: map[string]map[string][]models.PesanUnggahCsv{}, adaObj: adaObj}
}

func (u *unggahTiruan) BuatUnggah(_ context.Context, _ *db.Tx, x models.UnggahCsv) (string, error) {
	u.n++
	x.ID = fmt.Sprint(u.n)
	if x.Jenis == "" { // JENIS DEFAULT 'fire' (migrasi 258)
		x.Jenis = models.JenisCsvFire
	}
	u.kerja[x.ID] = &x
	return x.ID, nil
}

func (u *unggahTiruan) UbahUnggah(_ context.Context, _ *db.Tx, x models.UnggahCsv) error {
	if _, ok := u.kerja[x.ID]; !ok {
		return repository.ErrUnggahTidakAda
	}
	x.Basi = u.kerja[x.ID].Basi
	u.kerja[x.ID] = &x
	return nil
}

func (u *unggahTiruan) BacaUnggah(_ context.Context, id string) (models.UnggahCsv, error) {
	x, ok := u.kerja[id]
	if !ok {
		return models.UnggahCsv{}, repository.ErrUnggahTidakAda
	}
	return *x, nil
}

func (u *unggahTiruan) UnggahTerakhir(_ context.Context, caseID string) (models.UnggahCsv, bool, error) {
	for i := u.n; i >= 1; i-- {
		if x := u.kerja[fmt.Sprint(i)]; x != nil && x.CaseID == caseID {
			return *x, true, nil
		}
	}
	return models.UnggahCsv{}, false, nil
}

func (u *unggahTiruan) SisipBarisUnggah(_ context.Context, _ *db.Tx, id, _ string, b []models.BarisUnggahCsv) error {
	for _, x := range b {
		u.baris[id] = append(u.baris[id], models.BarisUnggahCsv{No: x.No, Isi: append([]string{}, x.Isi...)})
	}
	return nil
}

func (u *unggahTiruan) BacaBarisUnggah(_ context.Context, id, _ string, dari, sampai int) ([]models.BarisUnggahCsv, error) {
	var h []models.BarisUnggahCsv
	for _, x := range u.baris[id] {
		if x.No >= dari && x.No <= sampai {
			h = append(h, x)
		}
	}
	return h, nil
}

func (u *unggahTiruan) SisipPesanUnggah(_ context.Context, _ *db.Tx, id, jenis string, _ int, p []models.PesanUnggahCsv) error {
	if u.pesan[id] == nil {
		u.pesan[id] = map[string][]models.PesanUnggahCsv{}
	}
	u.pesan[id][jenis] = append(u.pesan[id][jenis], p...)
	return nil
}

func (u *unggahTiruan) BacaPesanUnggah(_ context.Context, id, jenis string, batas int) ([]models.PesanUnggahCsv, error) {
	p := u.pesan[id][jenis]
	if len(p) > batas {
		p = p[:batas]
	}
	return append([]models.PesanUnggahCsv{}, p...), nil
}

func (u *unggahTiruan) HapusPesanUnggah(_ context.Context, _ *db.Tx, id string) error {
	delete(u.pesan, id)
	return nil
}

func (u *unggahTiruan) AdaObjek(_ context.Context, caseID, _ string) (bool, error) {
	return u.adaObj(caseID), nil
}

// penambahObjekTiruan - TambahObjek menambah ke objekTiruan; gagalKe = panggilan ke-n gagal (sekali).
type penambahObjekTiruan struct {
	ob      objekTiruan
	awal    []int
	gagalKe int
	panggil int
}

func (p *penambahObjekTiruan) TambahObjek(_ context.Context, _ *db.Tx, id string, baris []models.ObjekFire, awal int) error {
	p.panggil++
	if p.panggil == p.gagalKe {
		return errors.New("UJI: koneksi putus")
	}
	p.awal = append(p.awal, awal)
	p.ob.ada[id] = append(p.ob.ada[id], baris...)
	return nil
}

// csvUji - teks berkas: judul 61 kolom (huruf kecil, BOM, pemisah ";") + n baris barisCsvUji (NO berurut); ubah(i, baris) boleh
// mengubah baris ke-i (1-based).
func csvUji(n int, ubah func(int, map[string]string)) string {
	var b strings.Builder
	b.WriteString("\uFEFF" + strings.ToLower(strings.Join(KolomCsvObjek, ";")) + ";\r\n")
	for i := 1; i <= n; i++ {
		r := barisCsvUji("NO", fmt.Sprint(i))
		if ubah != nil {
			ubah(i, r)
		}
		sel := make([]string, len(KolomCsvObjek))
		for j, k := range KolomCsvObjek {
			sel[j] = r[k]
		}
		b.WriteString(strings.Join(sel, ";") + "\r\n")
		if i == 1 {
			b.WriteString("\r\n;;;;\r\n") // baris kosong + baris sel kosong: dilewati, tidak dinomori
		}
	}
	return b.String()
}

func layananUnggahUji() (*Service, *unggahTiruan, *penambahObjekTiruan, objekTiruan, *csvTiruan) {
	c := csvMasterUji()
	svc, _, ob, _ := layananCsvUji(c, &penambahCsvTiruan{})
	u := barUnggahTiruan(func(id string) bool { return len(ob.ada[id]) > 0 })
	p := &penambahObjekTiruan{ob: ob}
	svc.DenganUnggahCsv(u, p)
	svc.latar = func(f func()) { f() }
	return svc, u, p, ob, c
}

// TestBacaCsvObjekMengalir - dialek layar tiket 66: BOM, judul tak peka huruf (61 / 55 kolom, sel kosong ujung), baris kosong /
// sel kosong semua dilewati, kutip ber-";", sel kurang = ""; judul salah / tidak ada -> ErrMasukanObjek; baris raksasa ditolak.
func TestBacaCsvObjekMengalir(t *testing.T) {
	var no []int
	var isi [][]string
	judul55 := strings.Join(KolomCsvObjek[:55], ";")
	n, err := BacaCsvObjekMengalir(strings.NewReader("\uFEFF"+judul55+"\n\n1;\"A;B\"\"C\";x\n;;\n  \n2\n"), nil,
		func(i int, s []string) error { no, isi = append(no, i), append(isi, s); return nil })
	if err != nil || n != 2 || len(no) != 2 || isi[0][1] != `A;B"C` || isi[0][2] != "x" || isi[1][0] != "2" || len(isi[1]) != 61 || isi[1][60] != "" {
		t.Fatalf("%d %v %q", n, err, isi)
	}
	for _, s := range []string{"", "\n\n", "SALAH;JUDUL\n1;2\n", strings.Join(KolomCsvObjek[:60], ";") + "\n"} {
		if _, err := BacaCsvObjekMengalir(strings.NewReader(s), nil, func(int, []string) error { return nil }); !errors.Is(err, ErrMasukanObjek) {
			t.Errorf("%q: %v", s[:min(len(s), 20)], err)
		}
	}
	raksasa := strings.Join(KolomCsvObjek, ";") + "\n" + strings.Repeat("x", batasPanjangBarisCsv+10) + "\n"
	if _, err := BacaCsvObjekMengalir(strings.NewReader(raksasa), nil, func(int, []string) error { return nil }); !errors.Is(err, ErrMasukanObjek) {
		t.Errorf("baris raksasa: %v", err)
	}
}

// pembangkitCsv - io.Reader berkas sintetis n baris tanpa menyimpan berkas utuh di memori.
type pembangkitCsv struct {
	n, i  int
	baris []byte
	sisa  []byte
}

func (g *pembangkitCsv) Read(p []byte) (int, error) {
	if len(g.sisa) == 0 {
		if g.i > g.n {
			return 0, io.EOF
		}
		if g.i == 0 {
			g.sisa = []byte(strings.Join(KolomCsvObjek, ";") + "\r\n")
		} else {
			g.sisa = g.baris
		}
		g.i++
	}
	n := copy(p, g.sisa)
	g.sisa = g.sisa[n:]
	return n, nil
}

// TestBacaCsvObjekSejutaBaris - tiket 75: 1 000 000 baris dibaca mengalir dan setiap baris divalidasi (uraiBarisCsv) dengan
// memori terbatas: puncak heap (dicuplik tiap 50 000 baris) di bawah 64 MiB, tidak tumbuh bersama jumlah baris.
func TestBacaCsvObjekSejutaBaris(t *testing.T) {
	if testing.Short() {
		t.Skip("uji 1 juta baris")
	}
	r := barisCsvUji()
	sel := make([]string, len(KolomCsvObjek))
	for j, k := range KolomCsvObjek {
		sel[j] = r[k]
	}
	g := &pembangkitCsv{n: 1_000_000, baris: []byte(strings.Join(sel, ";") + "\r\n")}
	runtime.GC()
	var puncak uint64
	var ms runtime.MemStats
	galat := 0
	n, err := BacaCsvObjekMengalir(g, nil, func(i int, isi []string) error {
		if _, m := uraiBarisCsv(i, petaBarisCsv(isi)); len(m) > 0 {
			galat++
		}
		if i%50_000 == 0 {
			runtime.ReadMemStats(&ms)
			puncak = max(puncak, ms.HeapAlloc)
		}
		return nil
	})
	if err != nil || n != 1_000_000 || galat != 0 {
		t.Fatalf("%d baris, %d galat, %v", n, galat, err)
	}
	t.Logf("1 000 000 baris; puncak heap %.1f MiB", float64(puncak)/(1<<20))
	if puncak > 64<<20 {
		t.Errorf("puncak heap %d byte > 64 MiB", puncak)
	}
}

// TestUnggahCsvAlur - F1 + F2: 2 500 baris (3 kelompok) -> valid (pesan / galat kosong, header tersimpan, baris kosong tidak
// dinomori) -> simpan: kelompok 1 MENGGANTI (General + CurrencyList), kelompok 2-3 MENAMBAH dengan SEQ mulai 1000 / 2000 ->
// selesai, "Success Insert Data", 2 500 lokasi; View Upload per halaman; adaObjek.
func TestUnggahCsvAlur(t *testing.T) {
	ctx := context.Background()
	svc, u, p, ob, c := layananUnggahUji()
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "uji.csv", strings.NewReader(csvUji(2500, nil)))
	if err != nil {
		t.Fatal(err)
	}
	x := st.Unggah
	if x.Status != StatusUnggahValid || x.TotalBaris != 2500 || x.BarisDiproses != 2500 || x.JumlahGalat != 0 || x.NamaBerkas != "uji.csv" ||
		x.NoOfferSlip != "UJI-SLIP-1" || x.BeginDate != "19/10/2026" || x.Ukuran == 0 || !st.AdaObjek || x.DibuatOleh != pelakuCsv.AkunID {
		t.Fatalf("validasi: %+v", x)
	}
	h, err := svc.BarisUnggahCsv(ctx, "UJI-NB-1", x.ID, 3, 1000)
	if err != nil || h.Total != 2500 || len(h.Baris) != 500 || h.Baris[0].No != 2001 || h.Baris[0].Isi[0] != "2001" ||
		h.Jenis != models.JenisCsvFire || len(h.Kolom) != len(KolomCsvObjek) {
		t.Fatalf("halaman 3: %d %d %v %+v", h.Total, len(h.Baris), err, h.Kolom)
	}
	st, err = svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", x.ID)
	if err != nil {
		t.Fatal(err)
	}
	x = st.Unggah
	if x.Status != StatusUnggahSelesai || x.KelompokTerakhir != 3 || x.BarisTersimpan != 2500 || len(st.Pesan) != 1 || st.Pesan[0].Teks != pesanSuksesCsv {
		t.Fatalf("simpan: %+v %v", x, st.Pesan)
	}
	if len(ob.ada["UJI-NB-1"]) != 2500 || ob.ada["UJI-NB-1"][0].ObjectNo != "1" || ob.ada["UJI-NB-1"][2499].ObjectNo != "2500" ||
		fmt.Sprint(p.awal) != "[1000 2000]" || len(c.hapus) != 1 {
		t.Errorf("objek %d, awal %v, hapus currency %d", len(ob.ada["UJI-NB-1"]), p.awal, len(c.hapus))
	}
	if _, err := svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", x.ID); !errors.Is(err, ErrUnggahStatus) {
		t.Errorf("simpan dua kali: %v", err)
	}
	if _, err := svc.StatusUnggahCsvKasus(ctx, "LAIN", x.ID); !errors.Is(err, ErrUnggahTidakAda) {
		t.Errorf("case lain: %v", err)
	}
	if akt, ada, err := svc.UnggahCsvAktif(ctx, "UJI-NB-1"); err != nil || !ada || akt.Unggah.ID != x.ID {
		t.Errorf("aktif: %v %v", ada, err)
	}
	_ = u
}

// TestUnggahCsvGagalValidasi - galat baris (angka, mata uang) ketat -> gagal_validasi, galat bernomor baris; simpan ditolak; judul
// salah -> 400 tanpa pekerjaan; berkas tanpa baris -> "Please Upload CSV".
func TestUnggahCsvGagalValidasi(t *testing.T) {
	ctx := context.Background()
	svc, u, _, _, _ := layananUnggahUji()
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "", strings.NewReader(csvUji(3, func(i int, r map[string]string) {
		if i == 2 {
			r["LOL"] = "X1"
		}
		if i == 3 {
			r["CURRENCY"] = "ZZZ"
		}
	})))
	if err != nil {
		t.Fatal(err)
	}
	if st.Unggah.Status != StatusUnggahGagalValidasi || st.Unggah.JumlahGalat != 2 || len(st.Galat) != 2 || st.Galat[0].NoBaris != 2 ||
		!strings.Contains(st.Galat[0].Teks, "kolom LOL") || !strings.Contains(st.Galat[1].Teks, "ZZZ") {
		t.Fatalf("%+v %v", st.Unggah, st.Galat)
	}
	if _, err := svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", st.Unggah.ID); !errors.Is(err, ErrUnggahStatus) {
		t.Errorf("simpan: %v", err)
	}
	n := u.n
	if _, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "", strings.NewReader("SALAH\n1\n")); !errors.Is(err, ErrMasukanObjek) || u.n != n {
		t.Errorf("judul salah: %v (%d pekerjaan baru)", err, u.n-n)
	}
	st, err = svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "", strings.NewReader(strings.Join(KolomCsvObjek, ";")+"\n"))
	if err != nil || st.Unggah.Status != StatusUnggahGagalValidasi || len(st.Galat) != 1 || st.Galat[0].Teks != pesanUnggahCsv {
		t.Errorf("tanpa baris: %v %+v", err, st)
	}
}

// TestUnggahCsvUlangi - simpan gagal di kelompok 3 (koneksi putus) -> gagal "gagal di baris 2001", kelompok 1-2 tetap; Ulangi
// melanjutkan dari kelompok 3 (tanpa baris ganda) -> selesai.
func TestUnggahCsvUlangi(t *testing.T) {
	ctx := context.Background()
	svc, _, p, ob, _ := layananUnggahUji()
	p.gagalKe = 2 // TambahObjek ke-2 = kelompok 3
	st, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "", strings.NewReader(csvUji(2500, nil)))
	if err != nil {
		t.Fatal(err)
	}
	st, err = svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", st.Unggah.ID)
	if err != nil {
		t.Fatal(err)
	}
	x := st.Unggah
	if x.Status != StatusUnggahGagal || x.Fase != FaseUnggahSimpan || x.KelompokTerakhir != 2 || x.BarisTersimpan != 2000 ||
		x.GalatDiBaris != 2001 || !strings.Contains(x.GalatProses, "gagal di baris 2001") || len(ob.ada["UJI-NB-1"]) != 2000 {
		t.Fatalf("gagal: %+v (%d objek)", x, len(ob.ada["UJI-NB-1"]))
	}
	st, err = svc.UlangiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", x.ID)
	if err != nil || st.Unggah.Status != StatusUnggahSelesai || len(ob.ada["UJI-NB-1"]) != 2500 || st.Unggah.BarisTersimpan != 2500 {
		t.Errorf("ulangi: %v %+v (%d objek)", err, st.Unggah, len(ob.ada["UJI-NB-1"]))
	}
	if _, err := svc.UlangiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", x.ID); !errors.Is(err, ErrUnggahStatus) {
		t.Errorf("ulangi selesai: %v", err)
	}
}

// TestKunciUnggahCsv - K75-2: pekerjaan berjalan (berdetak) -> PUT objek / Save Data lama / unggah baru 409; basi -> dianggap
// terputus: status menandainya gagal, kunci lepas, Ulangi boleh.
func TestKunciUnggahCsv(t *testing.T) {
	ctx := context.Background()
	svc, u, _, _, _ := layananUnggahUji()
	id, _ := u.BuatUnggah(ctx, nil, models.UnggahCsv{CaseID: "UJI-NB-1", Status: StatusUnggahMenyimpan, Fase: FaseUnggahSimpan, TotalBaris: 10})
	if _, err := svc.GantiObjek(ctx, pelakuCsv, "UJI-NB-1", []models.ObjekFire{{ObjectType: "Office"}}); !errors.Is(err, ErrUnggahBerjalan) {
		t.Errorf("PUT objek: %v", err)
	}
	if _, err := svc.SimpanCsvObjek(ctx, pelakuCsv, "UJI-NB-1", []map[string]string{barisCsvUji()}); !errors.Is(err, ErrUnggahBerjalan) {
		t.Errorf("Save Data lama: %v", err)
	}
	if _, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "", strings.NewReader(csvUji(1, nil))); !errors.Is(err, ErrUnggahBerjalan) {
		t.Errorf("unggah baru: %v", err)
	}
	u.kerja[id].Basi = true
	st, err := svc.StatusUnggahCsvKasus(ctx, "UJI-NB-1", id)
	if err != nil || st.Unggah.Status != StatusUnggahGagal || !strings.Contains(st.Unggah.GalatProses, "terputus") {
		t.Errorf("basi: %v %+v", err, st.Unggah)
	}
	if _, err := svc.GantiObjek(ctx, pelakuCsv, "UJI-NB-1", []models.ObjekFire{{ObjectType: "Office"}}); err != nil {
		t.Errorf("sesudah basi: %v", err)
	}
}
