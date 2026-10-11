package services

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ---- tiruan --------------------------------------------------------------------------------------------------------

// lampiranTiruan - DOCUMENT_POLIS + T_STORAGE_IMAGE dalam memori; status efek dibaca dari antreanTiruan.
type lampiranTiruan struct {
	kategori []string
	dok      map[string]models.DokumenPolis
	objek    map[string]models.ObjekStorage
	antrean  *antreanTiruan
	segar    []models.ObjekStorage
}

func lampiranUji(a *antreanTiruan) *lampiranTiruan {
	return &lampiranTiruan{kategori: []string{"EXTEND WPC", "OFFER", "SURVEYREPORT", "APPROVAL SPREADING"},
		dok: map[string]models.DokumenPolis{}, objek: map[string]models.ObjekStorage{}, antrean: a}
}

func (l *lampiranTiruan) KategoriLampiran(context.Context) ([]string, error) { return l.kategori, nil }
func (l *lampiranTiruan) JumlahLampiran(_ context.Context, idPega string) (map[string]int, error) {
	out := map[string]int{}
	for _, d := range l.dok {
		if repository.IDPegaProduksi(d.IDPega) == repository.IDPegaProduksi(idPega) { // dua bentuk (idpegaproduksi.go)
			out[d.Kategori2]++
		}
	}
	return out, nil
}
func (l *lampiranTiruan) lengkapi(d models.DokumenPolis) models.DokumenPolis {
	_, ada := l.objek[d.TStorageID]
	d.Tersimpan = ada && d.TStorageID != ""
	d.StatusEfek = l.antrean.statusTerakhir(d.ID, unggah.JenisEfekStorageUnggah)
	d.Tanggal = "08-10-2026 13:04"
	return d
}
func (l *lampiranTiruan) DaftarLampiran(_ context.Context, idPega, kategori, _, _ string) ([]models.DokumenPolis, error) {
	var out []models.DokumenPolis
	for _, d := range l.dok {
		if repository.IDPegaProduksi(d.IDPega) == repository.IDPegaProduksi(idPega) && d.Kategori2 == kategori {
			out = append(out, l.lengkapi(d))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func (l *lampiranTiruan) AmbilLampiran(_ context.Context, idPega, id, _, _ string) (models.DokumenPolis, error) {
	d, ada := l.dok[id]
	if !ada || repository.IDPegaProduksi(d.IDPega) != repository.IDPegaProduksi(idPega) {
		return models.DokumenPolis{}, repository.ErrLampiranTidakAda
	}
	return l.lengkapi(d), nil
}
func (l *lampiranTiruan) KunciLampiran(_ context.Context, _ *db.Tx, id string) (models.DokumenPolis, error) {
	d, ada := l.dok[id]
	if !ada {
		return models.DokumenPolis{}, repository.ErrLampiranTidakAda
	}
	return l.lengkapi(d), nil
}
func (l *lampiranTiruan) SisipLampiran(_ context.Context, _ *db.Tx, d models.DokumenPolis, saat time.Time) (string, error) {
	for {
		id := models.IDDokumenPolis(saat)
		if _, ada := l.dok[id]; !ada {
			d.ID = id
			l.dok[id] = d
			return id, nil
		}
		saat = saat.Add(time.Millisecond)
	}
}
func (l *lampiranTiruan) HapusLampiran(_ context.Context, _ *db.Tx, idPega, id string) error {
	if d, ada := l.dok[id]; !ada || d.IDPega != idPega {
		return repository.ErrLampiranTidakAda
	}
	delete(l.dok, id)
	return nil
}
func (l *lampiranTiruan) LinkStorage(_ context.Context, _ *db.Tx, imageID string) (models.LinkStorage, bool, error) {
	o, ada := l.objek[imageID]
	if !ada {
		return models.LinkStorage{}, false, nil
	}
	return models.LinkStorage{ImageID: imageID, URLPublic: o.URLPublic, AppFolder: o.AppFolder, App: o.App,
		Namafile: o.Namafile}, true, nil
}
func (l *lampiranTiruan) SimpanObjek(_ context.Context, _ *db.Tx, o models.ObjekStorage) error {
	if _, ada := l.objek[o.ImageID]; ada {
		return errors.New("ORA-00001 UJI")
	}
	l.objek[o.ImageID] = o
	return nil
}
func (l *lampiranTiruan) PerbaruiObjek(_ context.Context, _ *db.Tx, o models.ObjekStorage) error {
	l.segar = append(l.segar, o)
	return nil
}
func (l *lampiranTiruan) HapusObjek(_ context.Context, _ *db.Tx, imageID string) error {
	delete(l.objek, imageID)
	return nil
}
func (l *lampiranTiruan) AppStorage(context.Context) (string, error) { return "UJIAPP", nil }
func (l *lampiranTiruan) TokenStorageBerlaku(context.Context, *db.Tx, string, time.Time, time.Duration) (string, time.Duration, error) {
	return "", 0, nil
}
func (l *lampiranTiruan) PungutEfekRujukan(context.Context, *db.Tx, string, string, time.Time) (outbox.BarisEfekKeluar, error) {
	return outbox.BarisEfekKeluar{}, errors.New("tidak dipakai uji")
}

// antreanTiruan - outbox dalam memori (status / jadwal / percobaan per efek).
type efekTiruan struct {
	b      outbox.BarisEfekKeluar
	status string
	jadwal time.Time
	galat  string
}

type antreanTiruan struct {
	efek []*efekTiruan
	n    int
}

func (a *antreanTiruan) Antre(_ context.Context, _ *db.Tx, jenis, rujukan, muatan string, saat time.Time) error {
	a.n++
	a.efek = append(a.efek, &efekTiruan{b: outbox.BarisEfekKeluar{ID: string(rune('A' + a.n)), Modul: ModulOutboxNBFacIn,
		Jenis: jenis, Rujukan: rujukan, Muatan: muatan}, status: outbox.StatusEfekAntre, jadwal: saat})
	return nil
}
func (a *antreanTiruan) pungut(saat time.Time, rujukan string) (outbox.BarisEfekKeluar, error) {
	for _, e := range a.efek {
		if e.status == outbox.StatusEfekAntre && !e.jadwal.After(saat) && (rujukan == "" || e.b.Rujukan == rujukan) {
			e.status = outbox.StatusEfekJalan
			return e.b, nil
		}
	}
	return outbox.BarisEfekKeluar{}, outbox.ErrEfekTidakAda
}
func (a *antreanTiruan) Pungut(_ context.Context, _ *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error) {
	return a.pungut(saat, "")
}
func (a *antreanTiruan) PungutRujukan(_ context.Context, _ *db.Tx, rujukan string, saat time.Time) (outbox.BarisEfekKeluar, error) {
	return a.pungut(saat, rujukan)
}
func (a *antreanTiruan) Tuntaskan(_ context.Context, _ *db.Tx, id, status string, jadwal time.Time, galat string, _ time.Time) error {
	for _, e := range a.efek {
		if e.b.ID == id {
			e.status, e.jadwal, e.galat = status, jadwal, galat
			if status == outbox.StatusEfekAntre {
				e.b.Percobaan++
			}
			return nil
		}
	}
	return errors.New("efek tidak ada")
}
func (a *antreanTiruan) statusTerakhir(rujukan, jenis string) string {
	s := ""
	for _, e := range a.efek {
		if e.b.Rujukan == rujukan && e.b.Jenis == jenis {
			s = e.status
		}
	}
	return s
}

// storageGagal - membungkus stub lokal; galatSimpan / galatHapus menggantikan hasilnya bila terisi.
type storageGagal struct {
	KlienStorageLampiran
	galatSimpan error
	simpan      int
}

func (s *storageGagal) Simpan(ctx context.Context, lok models.LokasiObjek, isi io.Reader, ext, mime string) (models.ObjekStorage, error) {
	s.simpan++
	if s.galatSimpan != nil {
		return models.ObjekStorage{}, s.galatSimpan
	}
	return s.KlienStorageLampiran.Simpan(ctx, lok, isi, ext, mime)
}

const kasusLampiran = "UJI-NB-1"

var (
	pelakuPemegang = inti.Pelaku{AkunID: "UJI-MKT", Peran: []string{"ReasFacInMarketing"}}
	pelakuLain     = inti.Pelaku{AkunID: "UJI-UW", Peran: []string{"ReasFacInUnderwriting"}}
	pelakuIT       = inti.Pelaku{AkunID: "UJI-IT", Peran: []string{"ReasFacInUnderwriting"}}
	jamLampiran    = func() time.Time { return time.Date(2026, 10, 8, 6, 4, 5, 7e6, time.UTC) } // 13:04:05.007 WIB
)

type rangkaLampiran struct {
	svc     *Service
	repo    *lampiranTiruan
	antrean *antreanTiruan
	storage *storageGagal
	folder  string
}

func rangkaUji(t *testing.T, flag string) rangkaLampiran {
	t.Helper()
	folder := t.TempDir()
	a := &antreanTiruan{}
	repo := lampiranUji(a)
	st := &storageGagal{KlienStorageLampiran: StorageLampiranLokal(folder)}
	svc := Baru(nil).DenganTransaksi(tanpaTx).DenganJam(jamLampiran).
		DenganTahap(tahapTiruan{kasusLampiran: {Position: "ReasFacInMarketing", FlagOnGoingPolicy: flag}}).
		DenganLossRecord(nil, jabatanTiruan{"UJI-IT": "IT"}).
		DenganLampiranUji(repo, a, st, folder)
	return rangkaLampiran{svc: svc, repo: repo, antrean: a, storage: st, folder: folder}
}

func berkasUji(nama, kategori, isi string) unggah.BerkasMasuk {
	return unggah.BerkasMasuk{NamaFile: nama, Kategori: kategori, Isi: strings.NewReader(isi)}
}

func isiFolder(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			r, _ := filepath.Rel(dir, p)
			out = append(out, filepath.ToSlash(r))
		}
		return nil
	})
	sort.Strings(out)
	return out
}

// ---- uji -----------------------------------------------------------------------------------------------------------

// TestTampilanLampiran - kategori master apa adanya (urutan master), Count per KATEGORI_2 case (dokumen case lain tidak),
// bolehUnggah = siapa pun selama case belum jadi polis (K67-6), bolehHapus = flag "0" ATAU IT (siapa pun).
func TestTampilanLampiran(t *testing.T) {
	ctx := context.Background()
	r := rangkaUji(t, "0")
	r.repo.dok["X1"] = models.DokumenPolis{ID: "X1", IDPega: idPega(kasusLampiran), Kategori2: "OFFER"}
	r.repo.dok["X2"] = models.DokumenPolis{ID: "X2", IDPega: idPega(kasusLampiran), Kategori2: "OFFER", Kategori1: "reas"}
	r.repo.dok["X3"] = models.DokumenPolis{ID: "X3", IDPega: idPega("UJI-NB-2"), Kategori2: "OFFER"}
	tp, err := r.svc.TampilanLampiranKasus(ctx, pelakuPemegang, kasusLampiran)
	mau := []models.KategoriLampiran{{Kategori: "EXTEND WPC", Jumlah: 0}, {Kategori: "OFFER", Jumlah: 2}, {Kategori: "SURVEYREPORT", Jumlah: 0}, {Kategori: "APPROVAL SPREADING", Jumlah: 0}}
	if err != nil || !reflect.DeepEqual(tp.Kategori, mau) || !tp.BolehUnggah || !tp.BolehHapus {
		t.Fatalf("%v %+v", err, tp)
	}
	if tp, _ := r.svc.TampilanLampiranKasus(ctx, pelakuLain, kasusLampiran); !tp.BolehUnggah || !tp.BolehHapus {
		t.Errorf("bukan pemegang flag 0 (K67-6 boleh unggah): %+v", tp)
	}
	r = rangkaUji(t, "1")
	for _, c := range []struct {
		p     inti.Pelaku
		hapus bool
	}{{pelakuPemegang, false}, {pelakuLain, false}, {pelakuIT, true}} {
		if tp, err := r.svc.TampilanLampiranKasus(ctx, c.p, kasusLampiran); err != nil || tp.BolehHapus != c.hapus {
			t.Errorf("flag 1 %s: %v %+v", c.p.AkunID, err, tp)
		}
	}
	if _, err := r.svc.TampilanLampiranKasus(ctx, inti.Pelaku{}, kasusLampiran); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	if _, err := r.svc.TampilanLampiranKasus(ctx, pelakuPemegang, "UJI-NB-X"); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("404: %v", err)
	}
	if _, err := Baru(nil).TampilanLampiranKasus(ctx, pelakuPemegang, kasusLampiran); !errors.Is(err, ErrLampiranTanpaDatabase) {
		t.Errorf("503: %v", err)
	}
}

// TestUnggahLampiran - satu berkas: DOCUMENT_POLIS (ID cap waktu Pega jam 12-an, IDPEGA pzInsKey, MIME ekstensi, KATEGORI_1
// "Reas", KATEGORI_2 teks master, operator) + efek dijalankan seketika -> T_STORAGE_IMAGE (Namafile Pega), stub berisi,
// berkas antrean dibuang; daftar: status "", catatan = KATEGORI_2.
func TestUnggahLampiran(t *testing.T) {
	ctx := context.Background()
	r := rangkaUji(t, "1")
	id, err := r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji(`C:\fakepath\Slip Final.PDF`, "offer", "ISI-PDF-UJI"))
	if err != nil || id != "20261008010405007" {
		t.Fatalf("%v %q", err, id)
	}
	d := r.repo.dok[id]
	imageID := d.TStorageID
	if d.IDPega != "UJI-NB-1" || d.NamaFile != "Slip Final.PDF" || d.Mime != "pdf" || d.Kategori1 != "Reas" ||
		d.Kategori2 != "OFFER" || d.Operator != "UJI-MKT" || len(imageID) != unggah.PanjangImageID {
		t.Fatalf("dokumen %+v", d)
	}
	o, ada := r.repo.objek[imageID]
	if !ada || o.Namafile != "20261008-010405-7 - Slip Final.PDF" || o.ImageID != imageID {
		t.Fatalf("objek %+v", o)
	}
	if got := isiFolder(t, r.folder); !reflect.DeepEqual(got, []string{"nbfacin/penyimpanan/" + imageID}) {
		t.Errorf("folder %q", got)
	}
	if b, _ := os.ReadFile(filepath.Join(r.folder, "nbfacin", "penyimpanan", imageID)); string(b) != "ISI-PDF-UJI" {
		t.Errorf("isi stub %q", b)
	}
	if e := r.antrean.efek; len(e) != 1 || e[0].status != outbox.StatusEfekSelesai || e[0].b.Jenis != unggah.JenisEfekStorageUnggah ||
		e[0].b.Rujukan != id || strings.Contains(e[0].b.Muatan, "Slip") {
		t.Errorf("efek %+v", e[0])
	}
	baris, err := r.svc.DaftarBerkasLampiran(ctx, pelakuLain, kasusLampiran, "OFFER")
	if err != nil || len(baris) != 1 || baris[0] != (BarisLampiran{ID: id, NamaFile: "Slip Final.PDF", Catatan: "OFFER",
		Tanggal: "08-10-2026 13:04", Status: ""}) {
		t.Errorf("daftar %v %+v", err, baris)
	}
	// Unggahan kedua di milidetik yang sama: ID digeser satu milidetik.
	id2, err := r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("b.xlsx", "OFFER", "ISI-XLSX"))
	if err != nil || id2 != "20261008010405008" || r.repo.dok[id2].Mime != "xlsx" {
		t.Errorf("kedua %v %q", err, id2)
	}
	// Unduh: isi sama.
	meta, rc, err := r.svc.UnduhLampiran(ctx, pelakuLain, kasusLampiran, id)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(rc)
	_ = rc.Close()
	if meta.NamaFile != "Slip Final.PDF" || string(b) != "ISI-PDF-UJI" {
		t.Errorf("unduh %+v %q", meta, b)
	}
}

// TestUnggahLampiranDitolak - gerbang sebelum satu byte pun mendarat: bukan pemegang 403, jenis di luar GetMimeType 400
// berpesan, kategori kosong 400 / asing 422, berkas < 7 byte 400, tanpa UNGGAHAN_DIR 503; nol dokumen, nol efek, folder kosong.
func TestUnggahLampiranDitolak(t *testing.T) {
	ctx := context.Background()
	for nama, c := range map[string]struct {
		p      inti.Pelaku
		b      unggah.BerkasMasuk
		mau    error
		pesan  string
		folder bool
	}{
		"jenis asing":     {pelakuPemegang, berkasUji("a.exe", "OFFER", "ISI-PDF-UJI"), ErrMasukanLampiran, "Jenis berkas tidak didukung: .exe", true},
		"tanpa ekstensi":  {pelakuPemegang, berkasUji("README", "OFFER", "ISI-PDF-UJI"), ErrMasukanLampiran, "Jenis berkas tidak didukung: (tanpa ekstensi)", true},
		"kategori kosong": {pelakuPemegang, berkasUji("a.pdf", " ", "ISI-PDF-UJI"), ErrMasukanLampiran, "Kategori wajib dipilih.", true},
		"kategori asing":  {pelakuPemegang, berkasUji("a.pdf", "LAINNYA", "ISI-PDF-UJI"), ErrKategoriLampiranTidakDikenal, "", true},
		"kecil":           {pelakuPemegang, berkasUji("a.pdf", "OFFER", "123456"), ErrMasukanLampiran, "Berkas kosong atau terlalu kecil.", true},
		"kosong":          {pelakuPemegang, berkasUji("a.pdf", "OFFER", ""), ErrMasukanLampiran, "Berkas kosong atau terlalu kecil.", true},
		"tanpa nama":      {pelakuPemegang, berkasUji("", "OFFER", "ISI"), ErrMasukanLampiran, "Pilih berkas", true},
		"tanpa folder":    {pelakuPemegang, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"), unggah.ErrUnggahanDirBelumDisetel, "", false},
		"tanpa identitas": {inti.Pelaku{}, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"), inti.ErrTanpaIdentitas, "", true},
	} {
		r := rangkaUji(t, "0")
		if !c.folder {
			r.svc.unggahanDir = ""
		}
		_, err := r.svc.UnggahLampiran(ctx, c.p, kasusLampiran, c.b)
		if !errors.Is(err, c.mau) || !strings.Contains(errStr(err), c.pesan) || len(r.repo.dok) != 0 || len(r.antrean.efek) != 0 ||
			len(isiFolder(t, r.folder)) != 0 {
			t.Errorf("%s: %v (%d dok, %d efek, %q)", nama, err, len(r.repo.dok), len(r.antrean.efek), isiFolder(t, r.folder))
		}
	}
}

func errStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// TestUnggahLampiranStorageGagal - galat sementara: efek dijadwal ulang (status "menunggu", berkas antrean tetap); galat
// permanen: "gagal"; ulangi + penyimpanan pulih -> tersimpan, antrean dibuang. Unduh sebelum tersimpan 409.
func TestUnggahLampiranStorageGagal(t *testing.T) {
	ctx := context.Background()
	r := rangkaUji(t, "0")
	r.storage.galatSimpan = ErrStorageTakTerjangkau
	id, err := r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"))
	if err != nil {
		t.Fatal(err)
	}
	e := r.antrean.efek[0]
	if e.status != outbox.StatusEfekAntre || !e.jadwal.After(jamLampiran()) || e.b.Percobaan != 1 || !strings.Contains(e.galat, "tidak terjangkau") {
		t.Errorf("efek sementara %+v", e)
	}
	if b, _ := r.svc.DaftarBerkasLampiran(ctx, pelakuPemegang, kasusLampiran, "OFFER"); len(b) != 1 || b[0].Status != models.StatusLampiranMenunggu {
		t.Errorf("status menunggu %+v", b)
	}
	if _, _, err := r.svc.UnduhLampiran(ctx, pelakuPemegang, kasusLampiran, id); !errors.Is(err, ErrLampiranBelumTersimpan) {
		t.Errorf("unduh menunggu: %v", err)
	}
	// Permanen.
	r = rangkaUji(t, "0")
	r.storage.galatSimpan = ErrStorageMenolakPermintaan
	id, _ = r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"))
	if b, _ := r.svc.DaftarBerkasLampiran(ctx, pelakuPemegang, kasusLampiran, "OFFER"); len(b) != 1 || b[0].Status != models.StatusLampiranGagal {
		t.Errorf("status gagal %+v", b)
	}
	// Ulangi: case sudah jadi polis 403 (K67-6); belum polis -> efek baru, penyimpanan pulih -> tersimpan.
	tahapAsli := r.svc.tahap
	r.svc.DenganTahap(tahapTiruan{kasusLampiran: {Position: "ReasFacInMarketing", FlagOnGoingPolicy: "0", StatusWork: statusSelesaiPolicy}})
	if err := r.svc.UlangiLampiran(ctx, pelakuPemegang, kasusLampiran, id); !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("ulangi polis 403: %v", err)
	}
	r.svc.tahap = tahapAsli
	r.storage.galatSimpan = nil
	if err := r.svc.UlangiLampiran(ctx, pelakuPemegang, kasusLampiran, id); err != nil {
		t.Fatal(err)
	}
	if b, _ := r.svc.DaftarBerkasLampiran(ctx, pelakuPemegang, kasusLampiran, "OFFER"); len(b) != 1 || b[0].Status != "" ||
		len(r.antrean.efek) != 2 || len(isiFolder(t, r.folder)) != 1 {
		t.Errorf("sesudah ulangi %+v %d %q", b, len(r.antrean.efek), isiFolder(t, r.folder))
	}
	if err := r.svc.UlangiLampiran(ctx, pelakuPemegang, kasusLampiran, id); !errors.Is(err, ErrLampiranSudahTersimpan) {
		t.Errorf("ulangi tersimpan: %v", err)
	}
}

// TestUlangiSumberHilang - berkas antrean hilang dan penyimpanan tidak memilikinya: ulangi 409; pelaksana menyerah
// PERMANEN. Berkas antrean hilang tetapi objeknya ADA: pelaksana mencatat objek tanpa mengunggah ulang.
func TestUlangiSumberHilang(t *testing.T) {
	ctx := context.Background()
	r := rangkaUji(t, "0")
	r.storage.galatSimpan = ErrStorageTakTerjangkau
	id, _ := r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"))
	d := r.repo.dok[id]
	antre := filepath.Join(r.folder, "nbfacin", "antre", d.TStorageID+".pdf")
	if err := os.Remove(antre); err != nil {
		t.Fatal(err)
	}
	if err := r.svc.UlangiLampiran(ctx, pelakuPemegang, kasusLampiran, id); !errors.Is(err, ErrBerkasSumberLampiranHilang) {
		t.Errorf("ulangi: %v", err)
	}
	r.svc.jam = func() time.Time { return jamLampiran().Add(time.Hour) }
	if n, err := r.svc.JalankanAntreanLampiran(ctx, 5); err != nil || n != 1 || r.antrean.efek[0].status != outbox.StatusEfekGagalPermanen {
		t.Errorf("pelaksana: %v %d %+v", err, n, r.antrean.efek[0])
	}
	// Objek ada di penyimpanan (percobaan lama menulisnya lalu gagal dicatat).
	r = rangkaUji(t, "0")
	r.storage.galatSimpan = ErrStorageTakTerjangkau
	id, _ = r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"))
	d = r.repo.dok[id]
	lok, _ := models.LokasiUnggah(d.ID, d.TStorageID, d.NamaFile)
	if _, err := r.storage.KlienStorageLampiran.Simpan(ctx, lok, strings.NewReader("ISI-PDF-UJI"), "pdf", "application/pdf"); err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(filepath.Join(r.folder, "nbfacin", "antre", d.TStorageID+".pdf"))
	r.svc.jam = func() time.Time { return jamLampiran().Add(time.Hour) }
	if n, err := r.svc.JalankanAntreanLampiran(ctx, 5); err != nil || n != 1 || r.repo.objek[d.TStorageID].Namafile != lok.Namafile ||
		r.storage.simpan != 1 {
		t.Errorf("pulih: %v %d %+v %d", err, n, r.repo.objek, r.storage.simpan)
	}
}

// TestHapusLampiran - hapus FISIK (K-A): DOCUMENT_POLIS dibuang di transaksi, objek stub + T_STORAGE_IMAGE lewat efek
// storage-hapus; flag 1 bukan IT 403, IT boleh; dokumen Pega tanpa T_STORAGE_ID hanya barisnya; objek nyata (APPFOLDER)
// tidak disentuh stub (efek gagal permanen, catatan tetap); dokumen case lain 404.
func TestHapusLampiran(t *testing.T) {
	ctx := context.Background()
	r := rangkaUji(t, "0")
	id, _ := r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"))
	img := r.repo.dok[id].TStorageID
	if err := r.svc.HapusLampiran(ctx, pelakuLain, kasusLampiran, id); err != nil {
		t.Fatal(err)
	}
	if _, ada := r.repo.dok[id]; ada || len(r.repo.objek) != 0 || len(isiFolder(t, r.folder)) != 0 ||
		r.antrean.efek[1].b.Jenis != unggah.JenisEfekStorageHapus || r.antrean.efek[1].status != outbox.StatusEfekSelesai ||
		!strings.Contains(r.antrean.efek[1].b.Muatan, img) {
		t.Errorf("hapus: %+v %+v %q %+v", r.repo.dok, r.repo.objek, isiFolder(t, r.folder), r.antrean.efek[1])
	}
	if err := r.svc.HapusLampiran(ctx, pelakuLain, kasusLampiran, id); !errors.Is(err, ErrLampiranTidakAda) {
		t.Errorf("404: %v", err)
	}
	// Flag 1: bukan IT 403 (dokumen tetap), IT boleh.
	r = rangkaUji(t, "1")
	id, _ = r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"))
	if err := r.svc.HapusLampiran(ctx, pelakuPemegang, kasusLampiran, id); !errors.Is(err, inti.ErrTanpaWewenang) || len(r.repo.dok) != 1 {
		t.Errorf("flag 1 bukan IT: %v", err)
	}
	if err := r.svc.HapusLampiran(ctx, pelakuIT, kasusLampiran, id); err != nil || len(r.repo.dok) != 0 {
		t.Errorf("flag 1 IT: %v", err)
	}
	// Dokumen Pega: tanpa T_STORAGE_ID -> tanpa efek; dengan objek nyata -> stub menolak, catatan tetap.
	r = rangkaUji(t, "0")
	r.repo.dok["P1"] = models.DokumenPolis{ID: "P1", IDPega: idPega(kasusLampiran), Kategori2: "OFFER"}
	r.repo.dok["P2"] = models.DokumenPolis{ID: "P2", IDPega: idPega(kasusLampiran), Kategori2: "OFFER", TStorageID: "ABCDEF0123456789"}
	r.repo.objek["ABCDEF0123456789"] = models.ObjekStorage{ImageID: "ABCDEF0123456789", AppFolder: models.AppFolderUji("UJIAPP", "Policy/Doc/2025/01/x.pdf")}
	r.repo.dok["Q1"] = models.DokumenPolis{ID: "Q1", IDPega: idPega("UJI-NB-2"), Kategori2: "OFFER"}
	if err := r.svc.HapusLampiran(ctx, pelakuLain, kasusLampiran, "P1"); err != nil || len(r.antrean.efek) != 0 {
		t.Errorf("P1: %v %d", err, len(r.antrean.efek))
	}
	if err := r.svc.HapusLampiran(ctx, pelakuLain, kasusLampiran, "P2"); err != nil || len(r.antrean.efek) != 1 ||
		r.antrean.efek[0].status != outbox.StatusEfekGagalPermanen || len(r.repo.objek) != 1 {
		t.Errorf("P2: %v %+v %+v", err, r.antrean.efek, r.repo.objek)
	}
	if err := r.svc.HapusLampiran(ctx, pelakuLain, kasusLampiran, "Q1"); !errors.Is(err, ErrLampiranTidakAda) {
		t.Errorf("case lain: %v", err)
	}
}

// TestUnduhLampiranCaseLain - ID dokumen dari jalur URL dibatasi case-nya; objek tercatat tanpa berkas -> 409.
func TestUnduhLampiranCaseLain(t *testing.T) {
	ctx := context.Background()
	r := rangkaUji(t, "0")
	id, _ := r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"))
	r.svc.DenganTahap(tahapTiruan{kasusLampiran: {Position: "ReasFacInMarketing", FlagOnGoingPolicy: "0"}, "UJI-NB-2": {}})
	if _, _, err := r.svc.UnduhLampiran(ctx, pelakuPemegang, "UJI-NB-2", id); !errors.Is(err, ErrLampiranTidakAda) {
		t.Errorf("case lain: %v", err)
	}
	_ = os.Remove(filepath.Join(r.folder, "nbfacin", "penyimpanan", r.repo.dok[id].TStorageID))
	if _, _, err := r.svc.UnduhLampiran(ctx, pelakuPemegang, kasusLampiran, id); !errors.Is(err, ErrLampiranTanpaBerkas) {
		t.Errorf("tanpa berkas: %v", err)
	}
	if _, err := r.svc.DaftarBerkasLampiran(ctx, pelakuPemegang, kasusLampiran, ""); !errors.Is(err, ErrMasukanLampiran) {
		t.Errorf("kategori kosong: %v", err)
	}
}

// TestLayakUlangLampiran - keadaan permanen tidak diulang; jaringan diulang.
func TestLayakUlangLampiran(t *testing.T) {
	for _, e := range []error{ErrBerkasSumberLampiranHilang, ErrStorageMenolakPermintaan, ErrObjekNyataDiStub,
		repository.ErrAppStorageKosong, errMuatanLampiranRusak} {
		if layakUlangLampiran(e) {
			t.Errorf("%v diulang", e)
		}
	}
	if !layakUlangLampiran(ErrStorageTakTerjangkau) || !layakUlangLampiran(ErrStorageGagal) {
		t.Error("jaringan tidak diulang")
	}
	if g := ringkasGalatLampiran(errors.New(strings.Repeat("x", 1500))); len([]rune(g)) != 1000 {
		t.Errorf("ringkas %d", len(g))
	}
}

// TestUnggahLampiranK676 - K67-6 + ralat (keputusan work owner 09-10-2026): unggah oleh siapa pun yang membuka case (bukan
// pemegang workbasket pun boleh), kecuali STATUS_WORK "Resolved-Completed" / "Resolved-Decline" / "Resolved-Rejected" (ejaan
// Pega, ralat 10-10-2026) -> bolehUnggah false dan 403; status terbuka dan "Resolved-Reject" (bukan status Pega) tidak dikunci.
func TestUnggahLampiranK676(t *testing.T) {
	ctx := context.Background()
	r := rangkaUji(t, "1")
	if _, err := r.svc.UnggahLampiran(ctx, pelakuLain, kasusLampiran, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI")); err != nil {
		t.Fatalf("bukan pemegang: %v", err)
	}
	for status, boleh := range map[string]bool{"Resolved-Completed": false, "Resolved-Decline": false, "Resolved-Rejected": false,
		"": true, "Open-Offer": true, "Resolved-Reject": true} {
		r := rangkaUji(t, "1")
		r.svc.DenganTahap(tahapTiruan{kasusLampiran: {Position: "ReasFacInMarketing", FlagOnGoingPolicy: "1", StatusWork: status}})
		tp, err := r.svc.TampilanLampiranKasus(ctx, pelakuLain, kasusLampiran)
		if err != nil || tp.BolehUnggah != boleh {
			t.Errorf("%q tampilan: %v %+v", status, err, tp)
		}
		_, err = r.svc.UnggahLampiran(ctx, pelakuPemegang, kasusLampiran, berkasUji("a.pdf", "OFFER", "ISI-PDF-UJI"))
		if boleh != (err == nil) || (!boleh && (!errors.Is(err, inti.ErrTanpaWewenang) || !strings.Contains(err.Error(), pesanUnggahPolis))) {
			t.Errorf("%q unggah: %v", status, err)
		}
	}
}
