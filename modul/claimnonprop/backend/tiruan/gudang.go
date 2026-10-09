// Package tiruan adalah gudang dan acuan Claim Non Prop DALAM MEMORI untuk uji seam HTTP (handlers -> services). Ia meniru
// apa yang Oracle simpan: halaman disimpan lewat proyeksi katalog (hanya medan berkolom), angka dan tanggal dinormalkan
// seperti hasil TO_CHAR repository, riwayat hanya bertambah, ID adjustment stabil. Fixture berawalan `UJI-`.
package tiruan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/repository"
)

// Gudang - penyimpanan tiruan.
type Gudang struct {
	mu      sync.Mutex
	urut    int
	seq     map[string]int
	Kasus   map[string]models.Kasus
	halaman map[string]*models.Halaman
	// OS - baris OS_AKSEPTASI_KLAIM tertulis.
	OS []models.BarisOS
	// JSONKlaim - IDPEGA -> [MNK_NO_KLAIM, NOPOLIS].
	JSONKlaim  map[string][2]string
	Katastrofe []models.KatastrofeBaru
	// Komite - ID kasus komite -> tangga; KomiteAdj - ID kasus komite -> ID baris adjustment.
	Komite    map[string][]repository.AnggotaTangga
	KomiteAdj map[string]string
	Efek      []string
	// NomorTerbit - jenis -> urut terakhir (penghitung bersama).
	NomorTerbit map[string]int
	// TanpaSeqPLA - meniru PLATNP_SEQ yang tidak ada (ORA-02289, OQ-CNP-39).
	TanpaSeqPLA bool
	// KategoriDok - master T_KATEGORI_DOC_KLAIM TYPE_KLAIM NONPROP (ID -> LABEL); Dokumen - dokumen klaim tertulis;
	// Storage - catatan T_STORAGE_IMAGE (`Berkas.Catat`).
	KategoriDok map[string]string
	Dokumen     []models.BarisDokumenKlaim
	Storage     []penyimpanan.Objek
}

// Baru membuat gudang kosong.
func Baru() *Gudang {
	return &Gudang{seq: map[string]int{}, Kasus: map[string]models.Kasus{}, halaman: map[string]*models.Halaman{},
		JSONKlaim: map[string][2]string{}, Komite: map[string][]repository.AnggotaTangga{}, KomiteAdj: map[string]string{},
		NomorTerbit: map[string]int{}, KategoriDok: map[string]string{}}
}

// Transaksi - tiruan: salinan keadaan dipulihkan bila fn gagal (rollback).
func (g *Gudang) Transaksi(_ context.Context, fn func(tx *db.Tx) error) error {
	g.mu.Lock()
	cad := g.salin()
	g.mu.Unlock()
	if err := fn(nil); err != nil {
		g.mu.Lock()
		g.pulihkan(cad)
		g.mu.Unlock()
		return err
	}
	return nil
}

type cadangan struct {
	urut        int
	seq         map[string]int
	kasus       map[string]models.Kasus
	halaman     map[string]*models.Halaman
	os          []models.BarisOS
	json        map[string][2]string
	kat         []models.KatastrofeBaru
	komite      map[string][]repository.AnggotaTangga
	komiteAdj   map[string]string
	efek        []string
	nomorTerbit map[string]int
	dokumen     []models.BarisDokumenKlaim
	storage     []penyimpanan.Objek
}

func (g *Gudang) salin() cadangan {
	c := cadangan{urut: g.urut, seq: map[string]int{}, kasus: map[string]models.Kasus{}, halaman: map[string]*models.Halaman{},
		os: append([]models.BarisOS{}, g.OS...), json: map[string][2]string{},
		kat: append([]models.KatastrofeBaru{}, g.Katastrofe...), komite: map[string][]repository.AnggotaTangga{},
		komiteAdj: map[string]string{}, efek: append([]string{}, g.Efek...), nomorTerbit: map[string]int{},
		dokumen: append([]models.BarisDokumenKlaim{}, g.Dokumen...), storage: append([]penyimpanan.Objek{}, g.Storage...)}
	for k, v := range g.seq {
		c.seq[k] = v
	}
	for k, v := range g.Kasus {
		c.kasus[k] = v
	}
	for k, v := range g.halaman {
		c.halaman[k] = v.Salin()
	}
	for k, v := range g.JSONKlaim {
		c.json[k] = v
	}
	for k, v := range g.Komite {
		c.komite[k] = v
	}
	for k, v := range g.KomiteAdj {
		c.komiteAdj[k] = v
	}
	for k, v := range g.NomorTerbit {
		c.nomorTerbit[k] = v
	}
	return c
}

func (g *Gudang) pulihkan(c cadangan) {
	g.urut, g.seq, g.Kasus, g.halaman, g.OS, g.JSONKlaim = c.urut, c.seq, c.kasus, c.halaman, c.os, c.json
	g.Katastrofe, g.Komite, g.KomiteAdj, g.Efek, g.NomorTerbit = c.kat, c.komite, c.komiteAdj, c.efek, c.nomorTerbit
	g.Dokumen, g.Storage = c.dokumen, c.storage
}

func (g *Gudang) nomor(seq string) int {
	g.seq[seq]++
	return g.seq[seq]
}

// IDKasusBerikut - "CLMNP-000001".
func (g *Gudang) IDKasusBerikut(_ context.Context, _ *db.Tx, awalan string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return fmt.Sprintf("%s%06d", awalan, g.nomor("SEQ_WORK_CLAIM")), nil
}

// SisipKasus melahirkan kasus.
func (g *Gudang) SisipKasus(_ context.Context, _ *db.Tx, id, pembuat, nama string, saat time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Kasus[id] = models.Kasus{ID: id, Tahap: models.TahapOutstanding, PembuatID: pembuat, PembuatNama: nama,
		TglCreate: saat, TglUpdate: saat, Sumber: models.SumberGo}
	g.halaman[id] = models.HalamanBaru()
	return nil
}

// Keadaan membaca kasus.
func (g *Gudang) Keadaan(_ context.Context, _ *db.Tx, id string) (models.Kasus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k, ada := g.Kasus[id]
	if !ada {
		return models.Kasus{}, repository.ErrKasusTidakAda
	}
	return k, nil
}

// KunciKasus - tahap belum berubah.
func (g *Gudang) KunciKasus(ctx context.Context, tx *db.Tx, id, tahap string) (models.Kasus, error) {
	k, err := g.Keadaan(ctx, tx, id)
	if err != nil {
		return k, err
	}
	if k.Tahap != tahap || k.Tertutup() {
		return models.Kasus{}, repository.ErrTahapBerubah
	}
	return k, nil
}

// PindahTahap memindah tahap.
func (g *Gudang) PindahTahap(_ context.Context, _ *db.Tx, id, lama, baru, posisi string, saat time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := g.Kasus[id]
	if k.Tahap != lama || k.Tertutup() {
		return repository.ErrTahapBerubah
	}
	k.Tahap, k.Posisi, k.TglUpdate = baru, posisi, saat
	g.Kasus[id] = k
	return nil
}

// TutupKasus menutup kasus.
func (g *Gudang) TutupKasus(_ context.Context, _ *db.Tx, id, tahap string, saat time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := g.Kasus[id]
	if k.Tahap != tahap || k.Tertutup() {
		return repository.ErrTahapBerubah
	}
	k.StatusWork, k.TglUpdate = models.StatusSelesai, saat
	g.Kasus[id] = k
	return nil
}

// SentuhKasus memperbarui waktu.
func (g *Gudang) SentuhKasus(_ context.Context, _ *db.Tx, id string, saat time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := g.Kasus[id]
	k.TglUpdate = saat
	g.Kasus[id] = k
	return nil
}

// DaftarKasus - daftar kerja.
func (g *Gudang) DaftarKasus(_ context.Context, s repository.SaringanKasus) ([]repository.RingkasanKasus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := []repository.RingkasanKasus{}
	for _, k := range g.Kasus {
		if !strings.HasPrefix(k.ID, models.AwalanKlaim) {
			continue
		}
		if s.Selesai != k.Tertutup() || (!s.Selesai && s.Tahap != "" && k.Tahap != s.Tahap) ||
			(s.Pembuat != "" && k.PembuatID != s.Pembuat) {
			continue
		}
		h := g.halaman[k.ID]
		out = append(out, repository.RingkasanKasus{ID: k.ID, Tahap: k.Tahap, Label: models.LabelTahap[k.Tahap],
			StatusWork: k.StatusWork, PembuatNama: k.PembuatNama, NoClaim: h.Ambil(models.CD + "NoClaim"),
			PolicyNo: h.Ambil(models.CD + "PolicyData.PolicyNo")})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	return out, nil
}

// SimpanHalaman - proyeksi katalog + normalisasi; ID adjustment baru; riwayat bertambah.
func (g *Gudang) SimpanHalaman(_ context.Context, _ *db.Tx, id string, h *models.Halaman) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	lama := g.halaman[id]
	for _, b := range h.AmbilDaftar(models.DaftarAdjustment) {
		if b[models.PropID] == "" {
			b[models.PropID] = fmt.Sprintf("%d", g.nomor("SEQ_T_CLAIM"))
		}
	}
	if lama != nil { // baris berkomite tidak boleh hilang
		ada := map[string]bool{}
		for _, b := range h.AmbilDaftar(models.DaftarAdjustment) {
			ada[b[models.PropID]] = true
		}
		for _, b := range lama.AmbilDaftar(models.DaftarAdjustment) {
			if b[models.PropKomiteID] != "" && !ada[b[models.PropID]] {
				return repository.ErrAdjustmentBerkomite
			}
		}
	}
	s := models.ProyeksiKatalog(h)
	Normalkan(s)
	// SuggestList hanya bertambah (T_VIEW_SUGGEST): baris ber-penanda Baru diberi NO berikutnya.
	riw := models.SalinDaftar(lama.AmbilDaftar(models.DaftarRiwayat))
	for _, b := range h.AmbilDaftar(models.DaftarRiwayat) {
		if b[models.PropRiwayatBaru] == "1" {
			nb := b.Salin()
			delete(nb, models.PropRiwayatBaru)
			nb["No"] = fmt.Sprintf("%d", len(riw)+1)
			riw = append(riw, nb)
			delete(b, models.PropRiwayatBaru)
		}
	}
	s.SetelDaftar(models.DaftarRiwayat, riw)
	// KOMITE_ID ditulis SetelKomiteAdjustment, bukan dari halaman.
	komite := map[string]string{}
	for _, b := range lama.AmbilDaftar(models.DaftarAdjustment) {
		if b[models.PropKomiteID] != "" {
			komite[b[models.PropID]] = b[models.PropKomiteID]
		}
	}
	for _, b := range s.AmbilDaftar(models.DaftarAdjustment) {
		delete(b, models.PropKomiteID)
		if k := komite[b[models.PropID]]; k != "" {
			b[models.PropKomiteID] = k
		}
	}
	g.halaman[id] = s
	return nil
}

// BacaHalaman membaca salinan halaman tersimpan.
func (g *Gudang) BacaHalaman(_ context.Context, _ *db.Tx, id string) (*models.Halaman, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	h, ada := g.halaman[id]
	if !ada {
		return nil, repository.ErrKasusTidakAda
	}
	return h.Salin(), nil
}

// Halaman - halaman tersimpan (untuk asersi uji).
func (g *Gudang) Halaman(id string) *models.Halaman {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.halaman[id].Salin()
}

// SetelHalaman menimpa halaman tersimpan (fixture uji).
func (g *Gudang) SetelHalaman(id string, h *models.Halaman) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.halaman[id] = h.Salin()
}

// KodeProduksiUji - kode produksi NONLIFE tiruan.
const KodeProduksiUji = "UJI-"

// UrutNomor - penghitung (CLASS, JENIS, TAHUN) tiruan: periode = bulan saat.
func (g *Gudang) UrutNomor(_ context.Context, _ *db.Tx, huruf string, saat time.Time) (repository.BahanNomor, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	jenis := KodeProduksiUji + huruf
	g.NomorTerbit[jenis]++
	l := saat.In(models.Jakarta)
	return repository.BahanNomor{Jenis: jenis, MMYYYY: fmt.Sprintf("%02d.%04d", int(l.Month()), l.Year()),
		Urut: g.NomorTerbit[jenis]}, nil
}

// SisipOS menyimpan baris OS.
func (g *Gudang) SisipOS(_ context.Context, _ *db.Tx, b models.BarisOS, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.OS = append(g.OS, b)
	return nil
}

// NomorPLA - GenerateNoPLATNP tiruan.
func (g *Gudang) NomorPLA(_ context.Context, _ *db.Tx, oldID string, saat time.Time) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.TanpaSeqPLA {
		return "", errors.New("ORA-02289: sequence does not exist")
	}
	return models.RakitNomorPLA(oldID, saat, g.nomor("PLATNP_SEQ")), nil
}

// JumlahOS - Σ DATA_JSON baris STS 0 kasus (kedua bentuk kunci) per layer x mata uang.
func (g *Gudang) JumlahOS(_ context.Context, _ *db.Tx, kasusID, layer, mataUang string) (models.NilaiOS, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var kal models.Kalkulator
	jml := map[string]*apd.Decimal{}
	for _, k := range []string{"Adjusterfee", "GrossValue", "CNPOthersFee", "Salvage", "Value"} {
		jml[k] = apd.New(0, 0)
	}
	for _, b := range g.OS {
		if (b.CaseID != models.KunciInstans(kasusID) && b.CaseID != models.KunciPegaLama(kasusID)) ||
			b.StsReject != models.StsOSOutstanding || b.TypeLoss != layer || b.Currency != mataUang {
			continue
		}
		var isi map[string]any
		if err := json.Unmarshal([]byte(b.DataJSON), &isi); err != nil {
			return models.NilaiOS{}, err
		}
		for k := range jml {
			if v, ok := isi[k].(string); ok && v != "" {
				jml[k] = kal.Tambah(jml[k], kal.Teks(k, v))
			}
		}
	}
	return models.NilaiOS{Adjusterfee: models.Teks(jml["Adjusterfee"]), GrossValue: models.Teks(jml["GrossValue"]),
		CNPOthersFee: models.Teks(jml["CNPOthersFee"]), Salvage: models.Teks(jml["Salvage"]),
		Value: models.Teks(jml["Value"])}, kal.Galat()
}

// NomorKlaimOS - nomor klaim baris OS terakhir kasus.
func (g *Gudang) NomorKlaimOS(_ context.Context, _ *db.Tx, kasusID string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i := len(g.OS) - 1; i >= 0; i-- {
		b := g.OS[i]
		if (b.CaseID == models.KunciInstans(kasusID) || b.CaseID == models.KunciPegaLama(kasusID)) && b.NoClaim != "" {
			return b.NoClaim, nil
		}
	}
	return "", nil
}

// SalinJSONKlaim - INSERT bila IDPEGA belum ada.
func (g *Gudang) SalinJSONKlaim(_ context.Context, _ *db.Tx, idPega, no, pol string, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.JSONKlaim[idPega]; !ada {
		g.JSONKlaim[idPega] = [2]string{no, pol}
	}
	return nil
}

// SisipKatastrofe menyimpan katastrofe.
func (g *Gudang) SisipKatastrofe(_ context.Context, _ *db.Tx, k models.KatastrofeBaru, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Katastrofe = append(g.Katastrofe, k)
	return nil
}

// SetelKomiteAdjustment menautkan baris adjustment ke komite.
func (g *Gudang) SetelKomiteAdjustment(_ context.Context, _ *db.Tx, adjID, komiteID, komiteLama string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, h := range g.halaman {
		for _, b := range h.AmbilDaftar(models.DaftarAdjustment) {
			if b[models.PropID] == adjID {
				if b[models.PropKomiteID] != "" && b[models.PropKomiteID] != komiteLama {
					return errors.New("tiruan: adjustment sudah berkomite")
				}
				b[models.PropKomiteID] = komiteID
				return nil
			}
		}
	}
	return errors.New("tiruan: baris adjustment tidak ada")
}

// BuatKasusKomite melahirkan kasus komite KMTNP-.
func (g *Gudang) BuatKasusKomite(_ context.Context, _ *db.Tx, klaimID, adjID, pembuat, nama string,
	anggota []repository.AnggotaTangga, saat time.Time) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(anggota) == 0 {
		return "", errors.New("tiruan: tangga komite kosong")
	}
	id := fmt.Sprintf("%s%06d", models.AwalanKomite, g.nomor("SEQ_WORK_CLAIM"))
	g.Kasus[id] = models.Kasus{ID: id, Tahap: models.TahapKomite, PembuatID: pembuat, PembuatNama: nama,
		TglCreate: saat, Sumber: models.SumberGo}
	g.Komite[id] = anggota
	g.KomiteAdj[id] = adjID
	return id, nil
}

// AntreEfek mencatat efek.
func (g *Gudang) AntreEfek(_ context.Context, _ *db.Tx, jenis, rujukan, muatan string, _ time.Time) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Efek = append(g.Efek, jenis+"|"+rujukan+"|"+muatan)
	return fmt.Sprintf("%d", len(g.Efek)), nil
}

// Normalkan meniru bentuk nilai yang keluar dari Oracle lewat katalog: desimal ringkas ("1.50" -> "1.5"), tanggal
// "2006-01-02", tanggal-waktu "2006-01-02 15:04:05".
func Normalkan(h *models.Halaman) {
	nilai := func(k models.Kolom, v string) string {
		if v == "" {
			return v
		}
		switch {
		case k.Golongan.Desimal():
			if d, err := models.AngkaTeks(k.Properti, v); err == nil {
				return models.Teks(d)
			}
		case k.Golongan == models.GolTanggal:
			if t, ok := models.UraiTanggal(v); ok {
				return t.Format("2006-01-02")
			}
		case k.Golongan == models.GolTanggalWaktu:
			if t, ok := models.UraiTanggal(v); ok {
				return t.Format("2006-01-02 15:04:05")
			}
		}
		return v
	}
	for _, k := range models.TabelHeaderKlaim.Kolom {
		if v := h.Ambil(k.Properti); v != "" {
			h.Setel(k.Properti, nilai(k, v))
		}
	}
	baris := func(t models.Tabel, rows []models.Baris) {
		for _, b := range rows {
			for _, k := range t.Kolom {
				if v := b[k.Properti]; v != "" {
					b[k.Properti] = nilai(k, v)
				}
			}
		}
	}
	for _, t := range models.TabelAnakKlaim {
		baris(t, h.AmbilDaftar(t.Daftar))
	}
	adj := h.AmbilDaftar(models.DaftarAdjustment)
	baris(models.TabelAdjustment, adj)
	for i := range adj {
		for _, c := range models.TabelCucuAdjustment {
			baris(c, h.AmbilDaftar(models.JalurAdj(i+1, c.Daftar)))
		}
	}
}
