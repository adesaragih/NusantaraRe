package services_test

// Uji layanan lampiran tahun treaty - TANPA Oracle (tiket 12).
//
// Gudang, antrean outbox, dan penyimpanan dipalsukan. Penyimpanan adalah
// satu-satunya yang dipalsukan di uji HTTP; di sini ketiganya, supaya yang
// diuji adalah EFEKnya: berapa berkas ada di penyimpanan, status apa yang
// tampil, jejak apa yang tertulis, dan berkas antrean apa yang tertinggal.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/repository"
	"nusantarare/modul/treatycontractout/services"
)

// --- pemalsuan -------------------------------------------------------------

type efekLampiranUji struct {
	id, jenis, rujukan, muatan, status, galat string
	percobaan                                 int
	jadwal                                    time.Time
}

type antreanLampiranUji struct {
	mu   sync.Mutex
	efek []*efekLampiranUji
}

func (a *antreanLampiranUji) Antre(_ context.Context, _ *db.Tx, jenis, rujukan, muatan string, saat time.Time) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.efek = append(a.efek, &efekLampiranUji{id: fmt.Sprintf("E%03d", len(a.efek)+1), jenis: jenis,
		rujukan: rujukan, muatan: muatan, status: outbox.StatusEfekAntre, jadwal: saat})
	return nil
}

func (a *antreanLampiranUji) Pungut(_ context.Context, _ *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var pilih *efekLampiranUji
	for _, e := range a.efek {
		if e.status == outbox.StatusEfekAntre && !e.jadwal.After(saat) &&
			(pilih == nil || e.jadwal.Before(pilih.jadwal)) {
			pilih = e
		}
	}
	if pilih == nil {
		return outbox.BarisEfekKeluar{}, outbox.ErrEfekTidakAda
	}
	// Meniru `PungutEfek`: PERCOBAAN yang dikembalikan adalah nilai SEBELUM
	// dinaikkan.
	sebelum := pilih.percobaan
	pilih.percobaan++
	pilih.status = outbox.StatusEfekJalan
	return outbox.BarisEfekKeluar{ID: pilih.id, Jenis: pilih.jenis, Rujukan: pilih.rujukan,
		Muatan: pilih.muatan, Percobaan: sebelum}, nil
}

func (a *antreanLampiranUji) PungutRujukan(_ context.Context, _ *db.Tx, rujukan string, saat time.Time) (outbox.BarisEfekKeluar, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var pilih *efekLampiranUji
	for _, e := range a.efek {
		if e.rujukan == rujukan && e.status == outbox.StatusEfekAntre && !e.jadwal.After(saat) &&
			(pilih == nil || e.jadwal.Before(pilih.jadwal)) {
			pilih = e
		}
	}
	if pilih == nil {
		return outbox.BarisEfekKeluar{}, outbox.ErrEfekTidakAda
	}
	sebelum := pilih.percobaan
	pilih.percobaan++
	pilih.status = outbox.StatusEfekJalan
	return outbox.BarisEfekKeluar{ID: pilih.id, Jenis: pilih.jenis, Rujukan: pilih.rujukan,
		Muatan: pilih.muatan, Percobaan: sebelum}, nil
}
func (a *antreanLampiranUji) Tuntaskan(_ context.Context, _ *db.Tx, id, status string, jadwal time.Time,
	galat string, _ time.Time) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, e := range a.efek {
		if e.id == id {
			e.status, e.jadwal, e.galat = status, jadwal, galat
			return nil
		}
	}
	return errors.New("efek tidak ada")
}

// terakhir mengembalikan efek unggah terakhir sebuah lampiran.
func (a *antreanLampiranUji) terakhir(rujukan string) *efekLampiranUji {
	a.mu.Lock()
	defer a.mu.Unlock()
	var hasil *efekLampiranUji
	for _, e := range a.efek {
		if e.rujukan == rujukan && e.jenis == unggah.JenisEfekStorageUnggah {
			hasil = e
		}
	}
	return hasil
}

func (a *antreanLampiranUji) cacahStatus(jenis, status string) int {
	a.mu.Lock()
	defer a.mu.Unlock()
	n := 0
	for _, e := range a.efek {
		if e.jenis == jenis && e.status == status {
			n++
		}
	}
	return n
}

type gudangLampiranUji struct {
	// objek / objekDihapus - catatan T_STORAGE_IMAGE (tco4).
	objek        []models.ObjekPenyimpananTCO
	objekDihapus []string
	mu           sync.Mutex
	baris        map[string]models.LampiranTCO
	urut         int
	gagalSisip   error
	antrean      *antreanLampiranUji
}

func (g *gudangLampiranUji) barisDengan(l models.LampiranTCO) repository.BarisLampiranTCO {
	b := repository.BarisLampiranTCO{LampiranTCO: l}
	if e := g.antrean.terakhir(l.ID); e != nil {
		b.StatusEfek, b.GalatEfek, b.PercobaanEfek = e.status, e.galat, e.percobaan
	}
	return b
}

func (g *gudangLampiranUji) Daftar(_ context.Context, tahunID string) ([]repository.BarisLampiranTCO, error) {
	g.mu.Lock()
	var semua []models.LampiranTCO
	for _, l := range g.baris {
		if l.IDTreatyYear == tahunID {
			semua = append(semua, l)
		}
	}
	g.mu.Unlock()
	sort.Slice(semua, func(i, j int) bool { return semua[i].ID > semua[j].ID })
	var hasil []repository.BarisLampiranTCO
	for _, l := range semua {
		hasil = append(hasil, g.barisDengan(l))
	}
	return hasil, nil
}

func (g *gudangLampiranUji) Ambil(_ context.Context, tahunID, id string) (repository.BarisLampiranTCO, error) {
	g.mu.Lock()
	l, ada := g.baris[id]
	g.mu.Unlock()
	if !ada || l.IDTreatyYear != tahunID {
		return repository.BarisLampiranTCO{}, repository.ErrLampiranTidakAda
	}
	return g.barisDengan(l), nil
}

func (g *gudangLampiranUji) AmbilUntukKirim(_ context.Context, _ *db.Tx, id string) (models.LampiranTCO, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	l, ada := g.baris[id]
	if !ada {
		return models.LampiranTCO{}, repository.ErrLampiranTidakAda
	}
	return l, nil
}

func (g *gudangLampiranUji) Sisip(_ context.Context, _ *db.Tx, l models.LampiranTCO) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.gagalSisip != nil {
		return "", g.gagalSisip
	}
	g.urut++
	l.ID = fmt.Sprintf("1%09d", g.urut)
	g.baris[l.ID] = l
	return l.ID, nil
}

func (g *gudangLampiranUji) Hapus(_ context.Context, _ *db.Tx, tahunID, id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	l, ada := g.baris[id]
	if !ada || l.IDTreatyYear != tahunID {
		return repository.ErrLampiranTidakAda
	}
	delete(g.baris, id)
	return nil
}

// SimpanObjek meniru `Insert_T_Storage_SQL`: objek tercatat = lampiran
// ber-T_STORAGE_ID itu terkirim.
func (g *gudangLampiranUji) SimpanObjek(_ context.Context, _ *db.Tx, o models.ObjekPenyimpananTCO) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.objek = append(g.objek, o)
	for id, l := range g.baris {
		if l.ImageID == o.ImageID {
			l.TStorageID = o.ImageID
			g.baris[id] = l
		}
	}
	return nil
}

// HapusObjek meniru `DeleteStorage_SQL`.
func (g *gudangLampiranUji) HapusObjek(_ context.Context, _ *db.Tx, imageID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.objekDihapus = append(g.objekDihapus, imageID)
	for id, l := range g.baris {
		if l.ImageID == imageID {
			l.TStorageID = ""
			g.baris[id] = l
		}
	}
	return nil
}

// penyimpananLampiranUji meniru penyimpanan berkas: peta kunci -> isi.
type penyimpananLampiranUji struct {
	mu     sync.Mutex
	isi    map[string][]byte
	simpan int
	// gagal - galat Simpan; tulisDulu - tulis berkasnya LALU gagal (jawaban hilang).
	gagal     error
	tulisDulu bool
	// ekstensi - argumen `ekstensi` tiap Simpan (OQ-TCO-08: dari nama berkas asli).
	ekstensi []string
}

func (p *penyimpananLampiranUji) Simpan(_ context.Context, kunci string, isi io.Reader, _, ekstensi string) (models.ObjekPenyimpananTCO, error) {
	data, err := io.ReadAll(isi)
	if err != nil {
		return models.ObjekPenyimpananTCO{}, err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.simpan++
	p.ekstensi = append(p.ekstensi, ekstensi)
	if p.gagal != nil && !p.tulisDulu {
		return models.ObjekPenyimpananTCO{}, p.gagal
	}
	p.isi[kunci] = data
	return models.ObjekPenyimpananTCO{ImageID: kunci, Namafile: kunci, URLPublic: "UJI-URL-" + kunci}, p.gagal
}

func (p *penyimpananLampiranUji) Buka(_ context.Context, kunci string) (io.ReadCloser, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	d, ada := p.isi[kunci]
	if !ada {
		return nil, services.ErrBerkasTidakAdaDiPenyimpanan
	}
	return io.NopCloser(bytes.NewReader(d)), nil
}

func (p *penyimpananLampiranUji) Hapus(_ context.Context, kunci string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ada := p.isi[kunci]; !ada {
		return services.ErrBerkasTidakAdaDiPenyimpanan
	}
	delete(p.isi, kunci)
	return nil
}

func (p *penyimpananLampiranUji) Ada(_ context.Context, kunci string) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ada := p.isi[kunci]
	return ada, nil
}

func (p *penyimpananLampiranUji) cacah() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.isi)
}

func (p *penyimpananLampiranUji) setelGagal(err error, tulisDulu bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.gagal, p.tulisDulu = err, tulisDulu
}

type kategoriLampiranUji []string

func (k kategoriLampiranUji) Daftar(context.Context) ([]string, error) { return k, nil }

type tahunLampiranUji struct{}

func (tahunLampiranUji) Ambil(_ context.Context, id string) (models.TahunTreaty, error) {
	if id != "1000001" && id != "1000002" {
		return models.TahunTreaty{}, repository.ErrTahunTreatyTidakAda
	}
	return models.TahunTreaty{ID: id}, nil
}

// --- rakitan ---------------------------------------------------------------

type rakitLampiran struct {
	l       *services.LampiranTahunTCO
	gudang  *gudangLampiranUji
	antrean *antreanLampiranUji
	simpan  *penyimpananLampiranUji
	folder  string
	jam     *time.Time
}

func rakitanLampiran(t *testing.T) *rakitLampiran {
	t.Helper()
	a := &antreanLampiranUji{}
	g := &gudangLampiranUji{baris: map[string]models.LampiranTCO{}, antrean: a}
	p := &penyimpananLampiranUji{isi: map[string][]byte{}}
	saat := time.Date(2026, 9, 29, 9, 0, 0, 0, time.UTC)
	r := &rakitLampiran{gudang: g, antrean: a, simpan: p, folder: t.TempDir(), jam: &saat}
	r.l = services.New(nil).LampiranTahunTCO().
		DenganGudang(g).DenganAntrean(a).DenganPenyimpanan(p).
		DenganKategori(kategoriLampiranUji{"CLAUSES", "R/I SLIP", "OTHERS"}).
		DenganTahun(tahunLampiranUji{}).DenganFolder(r.folder).
		DenganTransaksi(transaksiUji).DenganJam(func() time.Time { return *r.jam })
	return r
}

func (r *rakitLampiran) maju(d time.Duration) { *r.jam = r.jam.Add(d) }

func (r *rakitLampiran) unggah(t *testing.T, tahun, nama, isi string) services.HasilLampiranTCO {
	t.Helper()
	h, err := r.l.Unggah(context.Background(), pelakuUjiTCO, tahun, unggah.BerkasMasuk{
		NamaFile: nama, Kategori: "r/i slip", Isi: strings.NewReader(isi)})
	if err != nil {
		t.Fatalf("unggah %s: %v", nama, err)
	}
	return h
}

// berkasAntre mencacah berkas di folder antrean modul ini.
func (r *rakitLampiran) berkasAntre(t *testing.T) []string {
	t.Helper()
	var nama []string
	_ = filepath.Walk(r.folder, func(jalur string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			nama = append(nama, jalur)
		}
		return nil
	})
	return nama
}

// --- uji -------------------------------------------------------------------

func TestLampiranTanpaIdentitasDitolak(t *testing.T) {
	r := rakitanLampiran(t)
	ctx := context.Background()
	if _, err := r.l.Daftar(ctx, inti.Pelaku{}, "1000001"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("daftar: %v", err)
	}
	if _, err := r.l.Unggah(ctx, inti.Pelaku{}, "1000001", unggah.BerkasMasuk{NamaFile: "a.pdf",
		Kategori: "CLAUSES", Isi: strings.NewReader("x")}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("unggah: %v", err)
	}
	if _, err := r.l.Hapus(ctx, inti.Pelaku{}, "1000001", "1000000001"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("hapus: %v", err)
	}
	if _, _, err := r.l.Unduh(ctx, inti.Pelaku{}, "1000001", "1000000001"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("unduh: %v", err)
	}
	if len(r.berkasAntre(t)) != 0 {
		t.Error("berkas mendarat tanpa identitas")
	}
}

func TestLampiranBawaanGagalTerang(t *testing.T) {
	_, err := services.New(nil).LampiranTahunTCO().DenganTransaksi(transaksiUji).
		Daftar(context.Background(), pelakuUjiTCO, "1000001")
	if !errors.Is(err, services.ErrGudangLampiranBelumDisuntik) {
		t.Errorf("bawaan harus gagal terang: %v", err)
	}
}

// AC 54, 56: berkas dilampirkan pada tahun treaty, kategori dari master,
// dan efeknya berjalan sampai terkirim.
func TestLampiranUnggahSampaiTerkirim(t *testing.T) {
	r := rakitanLampiran(t)
	h := r.unggah(t, "1000001", "kontrak.pdf", "ISI-UJI")
	if h.Peringatan != "" {
		t.Errorf("peringatan tak terduga: %s", h.Peringatan)
	}
	l := h.Lampiran
	if l.Status != models.StatusLampiranTerkirim || l.Category != "R/I SLIP" || l.UserID != "UJI-ADMIN" ||
		l.IDTreatyYear != "1000001" || l.FileMimeType == "" {
		t.Errorf("hasil: %+v", l)
	}
	if r.simpan.cacah() != 1 {
		t.Fatalf("penyimpanan berisi %d berkas, mau 1", r.simpan.cacah())
	}
	for _, d := range r.simpan.isi {
		if string(d) != "ISI-UJI" {
			t.Errorf("isi berkas di penyimpanan %q", d)
		}
	}
	if sisa := r.berkasAntre(t); len(sisa) != 0 {
		t.Errorf("berkas antrean tertinggal sesudah terkirim: %v", sisa)
	}
	if r.antrean.cacahStatus(unggah.JenisEfekStorageUnggah, outbox.StatusEfekSelesai) != 1 {
		t.Error("efek unggah tidak selesai")
	}
}

// AC 55, 58: kegagalan unggah TIDAK membatalkan apa pun - rekamnya ada,
// statusnya tertunda, galatnya terlihat, berkas antreannya tetap untuk
// dicoba lagi.
func TestLampiranUnggahGagalTidakMembatalkan(t *testing.T) {
	r := rakitanLampiran(t)
	r.simpan.setelGagal(errors.New("jaringan putus"), false)
	h := r.unggah(t, "1000001", "slip.xlsx", "ISI")
	if h.Lampiran.Status != models.StatusLampiranTertunda {
		t.Errorf("status %q, mau tertunda", h.Lampiran.Status)
	}
	if !strings.Contains(h.Lampiran.Galat, "jaringan putus") || h.Lampiran.Percobaan != 1 {
		t.Errorf("kegagalan tidak terlihat: %+v", h.Lampiran)
	}
	if len(r.gudang.baris) != 1 {
		t.Error("rekam lampiran hilang karena unggahan gagal")
	}
	if len(r.berkasAntre(t)) != 1 {
		t.Error("berkas antrean dibuang padahal belum terkirim")
	}
}

// AC 58: pengulangan TIDAK menggandakan berkas. Percobaan pertama menulis
// berkasnya lalu jawabannya hilang; percobaan kedua menulis ke KUNCI yang sama.
func TestLampiranPengulanganTidakMenggandakan(t *testing.T) {
	r := rakitanLampiran(t)
	r.simpan.setelGagal(errors.New("jawaban hilang"), true)
	h := r.unggah(t, "1000001", "kontrak.pdf", "ISI")
	if h.Lampiran.Status != models.StatusLampiranTertunda || r.simpan.cacah() != 1 {
		t.Fatalf("keadaan awal: status %q, berkas %d", h.Lampiran.Status, r.simpan.cacah())
	}
	r.simpan.setelGagal(nil, false)
	r.maju(24 * time.Hour)
	if _, err := r.l.JalankanAntrean(context.Background(), "UJI-PEKERJA", 5); err != nil {
		t.Fatal(err)
	}
	b, err := r.l.Daftar(context.Background(), pelakuUjiTCO, "1000001")
	if err != nil || len(b) != 1 || b[0].Status != models.StatusLampiranTerkirim {
		t.Fatalf("sesudah ulang: %+v %v", b, err)
	}
	if r.simpan.cacah() != 1 || r.simpan.simpan != 2 {
		t.Errorf("berkas di penyimpanan %d (mau 1), panggilan simpan %d (mau 2)", r.simpan.cacah(), r.simpan.simpan)
	}
	// Mengulang yang sudah terkirim ditolak, dan tetap satu berkas.
	if _, err := r.l.Ulangi(context.Background(), pelakuUjiTCO, "1000001", b[0].ID); !errors.Is(err, services.ErrLampiranSudahTerkirim) {
		t.Errorf("ulangi terkirim: %v", err)
	}
	if r.simpan.cacah() != 1 {
		t.Error("ulangi menggandakan berkas")
	}
}

// ADR-0015: galat permanen menyerah seketika dan masuk jejak; galat sementara
// menyerah sesudah jatahnya habis - tidak berputar selamanya.
func TestLampiranMenyerahTerlihat(t *testing.T) {
	r := rakitanLampiran(t)
	r.simpan.setelGagal(outbox.ErrPenyimpananBelumDisetujui, false)
	h := r.unggah(t, "1000001", "a.pdf", "ISI")
	// tco4: nol jejak modul - "menyerah" terlihat dari status outbox.
	if h.Lampiran.Status != models.StatusLampiranGagal || r.antrean.cacahStatus(unggah.JenisEfekStorageUnggah, outbox.StatusEfekGagalPermanen) != 1 {
		t.Errorf("permanen: status %q", h.Lampiran.Status)
	}

	r2 := rakitanLampiran(t)
	r2.simpan.setelGagal(errors.New("sementara"), false)
	h2 := r2.unggah(t, "1000001", "b.pdf", "ISI")
	for i := 0; i < 30 && r2.antrean.terakhir(h2.Lampiran.ID).status == outbox.StatusEfekAntre; i++ {
		r2.maju(48 * time.Hour)
		if _, err := r2.l.JalankanAntrean(context.Background(), "UJI-PEKERJA", 1); err != nil {
			t.Fatal(err)
		}
	}
	e := r2.antrean.terakhir(h2.Lampiran.ID)
	if e.status != outbox.StatusEfekGagalPermanen || e.percobaan > 20 {
		t.Errorf("sementara: status %q sesudah %d percobaan", e.status, e.percobaan)
	}
	if r2.antrean.cacahStatus(unggah.JenisEfekStorageUnggah, outbox.StatusEfekGagalPermanen) != 1 {
		t.Error("outbox tidak berstatus gagal permanen")
	}
}

// AC 58: yang gagal dapat diulang oleh manusia, dan hasilnya satu berkas.
func TestLampiranUlangiSesudahGagal(t *testing.T) {
	r := rakitanLampiran(t)
	r.simpan.setelGagal(outbox.ErrPenyimpananBelumDisetujui, false)
	h := r.unggah(t, "1000001", "a.pdf", "ISI")
	r.simpan.setelGagal(nil, false)
	u, err := r.l.Ulangi(context.Background(), pelakuUjiTCO, "1000001", h.Lampiran.ID)
	if err != nil || u.Lampiran.Status != models.StatusLampiranTerkirim {
		t.Fatalf("ulangi: %+v %v", u, err)
	}
	if r.simpan.cacah() != 1 {
		t.Errorf("berkas %d", r.simpan.cacah())
	}
}

// AC 56: kategori wajib dari master; master kosong adalah keadaan server.
func TestLampiranKategoriDariMaster(t *testing.T) {
	r := rakitanLampiran(t)
	_, err := r.l.Unggah(context.Background(), pelakuUjiTCO, "1000001", unggah.BerkasMasuk{
		NamaFile: "a.pdf", Kategori: "KARANGAN", Isi: strings.NewReader("x")})
	if !errors.Is(err, services.ErrKategoriLampiranTidakDikenal) {
		t.Errorf("kategori asing: %v", err)
	}
	if len(r.gudang.baris) != 0 || len(r.berkasAntre(t)) != 0 {
		t.Error("kategori asing tetap menulis")
	}
	kosong := r.l.DenganKategori(kategoriLampiranUji{})
	if _, err := kosong.Kategori(context.Background(), pelakuUjiTCO); !errors.Is(err, services.ErrKategoriLampiranKosong) {
		t.Errorf("master kosong: %v", err)
	}
	daftar, err := r.l.Kategori(context.Background(), pelakuUjiTCO)
	if err != nil || len(daftar) != 3 {
		t.Errorf("kategori: %v %v", daftar, err)
	}
}

// b376 "Tidak ada file yg diattach"; batas ukuran; folder belum disetel.
func TestLampiranGerbangBerkas(t *testing.T) {
	r := rakitanLampiran(t)
	ctx := context.Background()
	_, err := r.l.Unggah(ctx, pelakuUjiTCO, "1000001", unggah.BerkasMasuk{NamaFile: "a.pdf",
		Kategori: "CLAUSES", Isi: strings.NewReader("")})
	if !errors.Is(err, unggah.ErrBerkasKosong) || !strings.Contains(err.Error(), "No file attached") {
		t.Errorf("berkas kosong: %v", err)
	}
	_, err = r.l.DenganBatas(4).Unggah(ctx, pelakuUjiTCO, "1000001", unggah.BerkasMasuk{NamaFile: "a.pdf",
		Kategori: "CLAUSES", Isi: strings.NewReader("12345")})
	if !errors.Is(err, unggah.ErrBerkasTerlaluBesar) {
		t.Errorf("batas: %v", err)
	}
	_, err = r.l.DenganFolder("").Unggah(ctx, pelakuUjiTCO, "1000001", unggah.BerkasMasuk{NamaFile: "a.pdf",
		Kategori: "CLAUSES", Isi: strings.NewReader("1")})
	if !errors.Is(err, unggah.ErrUnggahanDirBelumDisetel) {
		t.Errorf("folder: %v", err)
	}
	_, err = r.l.Unggah(ctx, pelakuUjiTCO, "9999999", unggah.BerkasMasuk{NamaFile: "a.pdf",
		Kategori: "CLAUSES", Isi: strings.NewReader("1")})
	if !errors.Is(err, services.ErrTahunTreatyTidakAda) {
		t.Errorf("tahun tidak ada: %v", err)
	}
	if len(r.gudang.baris) != 0 || len(r.berkasAntre(t)) != 0 {
		t.Errorf("gerbang yang menolak tetap menulis: %v", r.berkasAntre(t))
	}
}

// Transaksi yang gagal membuang berkas antrean: tidak ada berkas tanpa rekam.
func TestLampiranTransaksiGagalMembuangBerkasAntre(t *testing.T) {
	r := rakitanLampiran(t)
	r.gudang.gagalSisip = errors.New("sisip gagal")
	_, err := r.l.Unggah(context.Background(), pelakuUjiTCO, "1000001", unggah.BerkasMasuk{
		NamaFile: "a.pdf", Kategori: "CLAUSES", Isi: strings.NewReader("ISI")})
	if err == nil {
		t.Fatal("galat sisip ditelan")
	}
	if sisa := r.berkasAntre(t); len(sisa) != 0 {
		t.Errorf("berkas yatim: %v", sisa)
	}
}

// Nama unggahan tidak menentukan jalur di disk.
func TestLampiranBerkasAntreTetapDiDalamFolder(t *testing.T) {
	r := rakitanLampiran(t)
	r.simpan.setelGagal(errors.New("tahan"), false)
	r.unggah(t, "1000001", `..\..\keluar.pdf`, "ISI")
	sisa := r.berkasAntre(t)
	if len(sisa) != 1 {
		t.Fatalf("berkas antrean: %v", sisa)
	}
	rel, err := filepath.Rel(r.folder, sisa[0])
	if err != nil || strings.HasPrefix(rel, "..") || !strings.HasSuffix(rel, ".pdf") {
		t.Errorf("berkas antrean di %q", rel)
	}
}

// Hapus berkas yang SUDAH tidak ada di penyimpanan tidak menggagalkan
// penghapusan rekamnya; jejak tertulis.
func TestLampiranHapusBerkasSudahTidakAda(t *testing.T) {
	r := rakitanLampiran(t)
	h := r.unggah(t, "1000001", "a.pdf", "ISI")
	r.simpan.isi = map[string][]byte{}
	if _, err := r.l.Hapus(context.Background(), pelakuUjiTCO, "1000001", h.Lampiran.ID); err != nil {
		t.Fatal(err)
	}
	if len(r.gudang.baris) != 0 {
		t.Errorf("rekam %d", len(r.gudang.baris))
	}
	if r.antrean.cacahStatus(unggah.JenisEfekStorageHapus, outbox.StatusEfekSelesai) != 1 {
		t.Error("efek hapus tidak selesai - berkas yang sudah tidak ada dianggap kegagalan")
	}
}

// Hapus sebelum terkirim: berkas antrean ikut dibuang, efek unggah yang masih
// menunggu tuntas tanpa kerja.
func TestLampiranHapusSebelumTerkirim(t *testing.T) {
	r := rakitanLampiran(t)
	r.simpan.setelGagal(errors.New("tahan"), false)
	h := r.unggah(t, "1000001", "a.pdf", "ISI")
	if _, err := r.l.Hapus(context.Background(), pelakuUjiTCO, "1000001", h.Lampiran.ID); err != nil {
		t.Fatal(err)
	}
	r.simpan.setelGagal(nil, false)
	r.maju(48 * time.Hour)
	if _, err := r.l.JalankanAntrean(context.Background(), "UJI-PEKERJA", 5); err != nil {
		t.Fatal(err)
	}
	if sisa := r.berkasAntre(t); len(sisa) != 0 {
		t.Errorf("berkas antrean tertinggal: %v", sisa)
	}
	if r.simpan.cacah() != 0 {
		t.Error("lampiran yang sudah dihapus tetap diunggah")
	}
	if r.antrean.cacahStatus(unggah.JenisEfekStorageUnggah, outbox.StatusEfekSelesai) != 1 {
		t.Error("efek unggah lampiran terhapus tidak tuntas")
	}
}

// AC 61: rekam tanpa berkas terdeteksi dan dapat diperbaiki.
func TestLampiranSelarasMendeteksiRekamTanpaBerkas(t *testing.T) {
	r := rakitanLampiran(t)
	ctx := context.Background()
	h := r.unggah(t, "1000001", "a.pdf", "ISI")
	r.simpan.isi = map[string][]byte{}
	temuan, err := r.l.PeriksaSelaras(ctx, pelakuUjiTCO, "1000001")
	if err != nil || len(temuan) != 1 || temuan[0].LampiranID != h.Lampiran.ID {
		t.Fatalf("temuan: %+v %v", temuan, err)
	}
	if _, _, err := r.l.Unduh(ctx, pelakuUjiTCO, "1000001", h.Lampiran.ID); !errors.Is(err, services.ErrLampiranTanpaBerkas) {
		t.Errorf("unduh rekam tanpa berkas: %v", err)
	}
	if _, err := r.l.Ulangi(ctx, pelakuUjiTCO, "1000001", h.Lampiran.ID); !errors.Is(err, services.ErrBerkasSumberLampiranHilang) {
		t.Errorf("ulangi tanpa sumber: %v", err)
	}
	if _, err := r.l.Hapus(ctx, pelakuUjiTCO, "1000001", h.Lampiran.ID); err != nil {
		t.Fatal(err)
	}
	if temuan, _ := r.l.PeriksaSelaras(ctx, pelakuUjiTCO, "1000001"); len(temuan) != 0 {
		t.Errorf("sesudah diperbaiki: %+v", temuan)
	}
}

// AC 61 sisi lain: rekam terkirim yang berkasnya hilang, tetapi berkas
// antreannya masih ada, diunggah ulang oleh "ulangi".
func TestLampiranUlangiMemulihkanDariAntrean(t *testing.T) {
	r := rakitanLampiran(t)
	r.simpan.setelGagal(errors.New("jawaban hilang"), true)
	h := r.unggah(t, "1000001", "a.pdf", "ISI")
	// Berkas sudah di penyimpanan tetapi rekam belum tahu; lalu penyimpanan
	// kehilangannya. Berkas antrean masih ada.
	r.simpan.isi = map[string][]byte{}
	r.simpan.setelGagal(nil, false)
	u, err := r.l.Ulangi(context.Background(), pelakuUjiTCO, "1000001", h.Lampiran.ID)
	if err != nil || u.Lampiran.Status != models.StatusLampiranTerkirim || r.simpan.cacah() != 1 {
		t.Errorf("ulangi: %+v %v, berkas %d", u, err, r.simpan.cacah())
	}
}

// AC 57: unduh satu dan unduh semua; batas tahun treaty ditegakkan.
func TestLampiranUnduh(t *testing.T) {
	r := rakitanLampiran(t)
	ctx := context.Background()
	a := r.unggah(t, "1000001", "kontrak.pdf", "ISI-A")
	b := r.unggah(t, "1000001", `C:\x\kontrak.pdf`, "ISI-B")
	meta, rc, err := r.l.Unduh(ctx, pelakuUjiTCO, "1000001", a.Lampiran.ID)
	if err != nil {
		t.Fatal(err)
	}
	isi, _ := io.ReadAll(rc)
	_ = rc.Close()
	if string(isi) != "ISI-A" || meta.FileName != "kontrak.pdf" {
		t.Errorf("unduh: %q %+v", isi, meta)
	}
	if _, _, err := r.l.Unduh(ctx, pelakuUjiTCO, "1000002", a.Lampiran.ID); !errors.Is(err, services.ErrLampiranTidakAda) {
		t.Errorf("lampiran tahun lain terunduh: %v", err)
	}
	entri := map[string]string{}
	err = r.l.UnduhSemua(ctx, pelakuUjiTCO, "1000001", func(nama string, isi io.Reader) error {
		d, _ := io.ReadAll(isi)
		entri[nama] = string(d)
		return nil
	})
	if err != nil || len(entri) != 2 || entri[a.Lampiran.ID+"_kontrak.pdf"] != "ISI-A" ||
		entri[b.Lampiran.ID+"_kontrak.pdf"] != "ISI-B" {
		t.Errorf("unduh semua: %v %v", entri, err)
	}
	r.simpan.setelGagal(errors.New("tahan"), false)
	c := r.unggah(t, "1000001", "c.pdf", "ISI-C")
	if _, _, err := r.l.Unduh(ctx, pelakuUjiTCO, "1000001", c.Lampiran.ID); !errors.Is(err, services.ErrLampiranBelumTerkirim) {
		t.Errorf("unduh tertunda: %v", err)
	}
}

// Temuan /code-review: aksi seorang pemakai hanya menjalankan efek MILIK
// lampirannya - efek lampiran lain (pemakai lain) tetap di antrean, dan tidak
// ada jejak yang tercatat atas nama pemakai yang salah.
func TestLampiranAksiHanyaMenjalankanEfekMiliknya(t *testing.T) {
	r := rakitanLampiran(t)
	if err := r.antrean.Antre(context.Background(), nil, unggah.JenisEfekStorageUnggah, "UJI-LAIN", "{}", *r.jam); err != nil {
		t.Fatal(err)
	}
	r.unggah(t, "1000001", "kontrak.pdf", "ISI-UJI")
	for _, e := range r.antrean.efek {
		if e.rujukan == "UJI-LAIN" && (e.status != outbox.StatusEfekAntre || e.percobaan != 0) {
			t.Errorf("efek lampiran lain ikut dijalankan: %+v", e)
		}
	}
	if r.antrean.cacahStatus(unggah.JenisEfekStorageUnggah, outbox.StatusEfekSelesai) != 1 {
		t.Error("efek lampiran sendiri tidak selesai")
	}
}

// tco4: objek terkirim DICATAT seperti Insert_T_Storage_SQL (URL dari jawaban
// penyimpanan), dan dibuang seperti DeleteStorage_SQL sesudah hapus berhasil.
func TestLampiranObjekPenyimpananDicatatDanDibuang(t *testing.T) {
	r := rakitanLampiran(t)
	h := r.unggah(t, "1000001", "a.pdf", "ISI")
	r.gudang.mu.Lock()
	objek := append([]models.ObjekPenyimpananTCO(nil), r.gudang.objek...)
	r.gudang.mu.Unlock()
	if len(objek) != 1 || objek[0].ImageID == "" || objek[0].URLPublic != "UJI-URL-"+objek[0].ImageID ||
		objek[0].Namafile != objek[0].ImageID {
		t.Fatalf("objek tercatat: %+v", objek)
	}
	if _, err := r.l.Hapus(context.Background(), pelakuUjiTCO, "1000001", h.Lampiran.ID); err != nil {
		t.Fatal(err)
	}
	r.gudang.mu.Lock()
	defer r.gudang.mu.Unlock()
	if len(r.gudang.objekDihapus) != 1 || r.gudang.objekDihapus[0] != objek[0].ImageID {
		t.Errorf("catatan objek dibuang: %v", r.gudang.objekDihapus)
	}
}
