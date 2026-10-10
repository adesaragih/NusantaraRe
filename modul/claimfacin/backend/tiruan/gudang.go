// Package tiruan adalah gudang dan acuan Claim Fac In DALAM MEMORI untuk uji seam HTTP (handlers -> services). Ia meniru
// apa yang Oracle simpan: halaman disimpan lewat proyeksi katalog (hanya medan berkolom), angka dan tanggal dinormalkan
// seperti hasil TO_CHAR repository, kronologi hanya bertambah, ID objek / item / adjustment stabil. Fixture berawalan
// `UJI-`.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// Gudang - penyimpanan tiruan.
type Gudang struct {
	mu sync.Mutex
	// GagalHapusDok - galat HapusDokumenKlaim (baris terkunci / ORA-) untuk uji urutan hapus lampiran.
	GagalHapusDok error
	seq           map[string]int
	Kasus         map[string]models.Kasus
	halaman       map[string]*models.Halaman
	// OS - baris OS_AKSEPTASI_KLAIM tertulis.
	OS []models.BarisOS
	// JSONKlaim - IDPEGA -> [MNK_NO_KLAIM, NOPOLIS].
	JSONKlaim  map[string][2]string
	Katastrofe []models.BarisKatastrofe
	Log        []repository.LogLayanan
	Progres    []models.Progres
	SubProgres []models.SubProgres
	Efek       []string
	// NomorTerbit - jenis -> urut terakhir (penghitung bersama).
	NomorTerbit map[string]int
	// DLAOS - "AcceptedNo=NoDLA" setiap InsertDLA_OS_SQL.
	DLAOS []string
	// KategoriDok - master T_KATEGORI_DOC_KLAIM FAC (ID -> LABEL); Dokumen - baris DOCUMENT_CLAIM.
	KategoriDok map[string]string
	Dokumen     []models.BarisDokumenKlaim
	// Storage - catatan T_STORAGE_IMAGE (`Berkas.Catat`).
	Storage []penyimpanan.Objek
	// Komite - kasus komite lahir: ID -> (adjustment, tangga).
	Komite map[string]KasusKomite
}

// KasusKomite - satu kasus komite tiruan.
type KasusKomite struct {
	// Transfer - TRANSFER_TYPE (2 TT2, 3 TT3, 4 TT4).
	Transfer                      string
	KlaimID, AdjustmentID, Posisi string
	Anggota                       []repository.AnggotaTangga
}

// Baru membuat gudang kosong.
func Baru() *Gudang {
	return &Gudang{seq: map[string]int{}, Kasus: map[string]models.Kasus{}, halaman: map[string]*models.Halaman{},
		JSONKlaim: map[string][2]string{}, NomorTerbit: map[string]int{}, KategoriDok: map[string]string{},
		Komite: map[string]KasusKomite{}}
}

type cadangan struct {
	seq         map[string]int
	kasus       map[string]models.Kasus
	halaman     map[string]*models.Halaman
	os          []models.BarisOS
	json        map[string][2]string
	kat         []models.BarisKatastrofe
	log         []repository.LogLayanan
	progres     []models.Progres
	sub         []models.SubProgres
	efek        []string
	nomorTerbit map[string]int
	dlaOS       []string
	dokumen     []models.BarisDokumenKlaim
	storage     []penyimpanan.Objek
}

// Transaksi - tiruan: salinan keadaan dipulihkan bila fn gagal (rollback).
func (g *Gudang) Transaksi(_ context.Context, fn func(tx *db.Tx) error) error {
	g.mu.Lock()
	c := cadangan{seq: map[string]int{}, kasus: map[string]models.Kasus{}, halaman: map[string]*models.Halaman{},
		os: append([]models.BarisOS{}, g.OS...), json: map[string][2]string{},
		kat: append([]models.BarisKatastrofe{}, g.Katastrofe...), log: append([]repository.LogLayanan{}, g.Log...),
		progres: append([]models.Progres{}, g.Progres...), sub: append([]models.SubProgres{}, g.SubProgres...),
		efek: append([]string{}, g.Efek...), nomorTerbit: map[string]int{}, dlaOS: append([]string{}, g.DLAOS...),
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
	for k, v := range g.NomorTerbit {
		c.nomorTerbit[k] = v
	}
	g.mu.Unlock()
	if err := fn(nil); err != nil {
		g.mu.Lock()
		g.seq, g.Kasus, g.halaman, g.OS, g.JSONKlaim = c.seq, c.kasus, c.halaman, c.os, c.json
		g.Katastrofe, g.Log, g.Progres, g.SubProgres, g.Efek, g.NomorTerbit = c.kat, c.log, c.progres, c.sub, c.efek,
			c.nomorTerbit
		g.DLAOS, g.Dokumen, g.Storage = c.dlaOS, c.dokumen, c.storage
		g.mu.Unlock()
		return err
	}
	return nil
}

func (g *Gudang) nomor(seq string) int {
	g.seq[seq]++
	return g.seq[seq]
}

// IDKasusBerikut - "CLM-000001".
func (g *Gudang) IDKasusBerikut(_ context.Context, _ *db.Tx, awalan string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return fmt.Sprintf("%s%06d", awalan, g.nomor("SEQ_WORK_CLAIM")), nil
}

// SisipKasus melahirkan kasus di tahap Input Register.
func (g *Gudang) SisipKasus(_ context.Context, _ *db.Tx, id, pembuat, nama string, saat time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Kasus[id] = models.Kasus{ID: id, Tahap: models.TahapRegister, PembuatID: pembuat, PembuatNama: nama,
		TglCreate: saat, TglUpdate: saat, Sumber: models.SumberGo}
	g.halaman[id] = models.HalamanBaru()
	return nil
}

// Keadaan membaca kasus.
func (g *Gudang) Keadaan(_ context.Context, _ *db.Tx, id string) (models.Kasus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k, ada := g.Kasus[id]
	if !ada || k.Tahap == models.TahapKomite {
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
func (g *Gudang) TutupKasus(ctx context.Context, tx *db.Tx, id, tahap string, saat time.Time) error {
	return g.TutupKasusStatus(ctx, tx, id, tahap, models.StatusSelesai, saat)
}

// TutupKasusStatus - penutupan kasus berstatus `status` (kontrak Komite Claim Fac In).
func (g *Gudang) TutupKasusStatus(_ context.Context, _ *db.Tx, id, tahap, status string, saat time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := g.Kasus[id]
	if k.Tahap != tahap || k.Tertutup() {
		return repository.ErrTahapBerubah
	}
	k.StatusWork, k.TglUpdate = status, saat
	g.Kasus[id] = k
	return nil
}

// TutupKomiteAnak - lihat `repository.Gudang.TutupKomiteAnak`.
func (g *Gudang) TutupKomiteAnak(_ context.Context, _ *db.Tx, klaimID, kecuali, status string, saat time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	for id, km := range g.Komite {
		k := g.Kasus[id]
		if km.KlaimID != klaimID || id == kecuali || k.Tahap != models.TahapKomite || k.Tertutup() {
			continue
		}
		k.StatusWork, k.TglUpdate, k.Posisi = status, saat, ""
		g.Kasus[id] = k
	}
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

// DaftarKasus - daftar kerja (lini disaring lewat tahap, bukan awalan).
func (g *Gudang) DaftarKasus(_ context.Context, s repository.SaringanKasus) ([]repository.RingkasanKasus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := []repository.RingkasanKasus{}
	for _, k := range g.Kasus {
		if k.Tahap == models.TahapKomite || s.Selesai != k.Tertutup() {
			continue
		}
		if !s.Selesai {
			if s.Tahap != "" && k.Tahap != s.Tahap {
				continue
			}
			if len(s.TahapIn) > 0 && !berisi(s.TahapIn, k.Tahap) {
				continue
			}
		}
		if s.Pembuat != "" && !strings.EqualFold(k.PembuatID, s.Pembuat) {
			continue
		}
		h := g.halaman[k.ID]
		out = append(out, repository.RingkasanKasus{ID: k.ID, Tahap: k.Tahap, Label: models.LabelTahap[k.Tahap],
			StatusWork: k.StatusWork, PembuatNama: k.PembuatNama, NoClaim: h.Ambil(models.CD + "NoClaim"),
			PolicyNo: h.Ambil(models.JalurNoPolis)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	return out, nil
}

func berisi(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// SimpanHalaman - proyeksi katalog + normalisasi; ID stabil baru; kronologi bertambah.
func (g *Gudang) SimpanHalaman(_ context.Context, _ *db.Tx, id string, h *models.Halaman) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	lama := g.halaman[id]
	if err := jagaKomite(lama, h); err != nil {
		return err
	}
	for o, ob := range h.AmbilDaftar(models.DaftarObjek) {
		if ob[models.PropID] == "" {
			ob[models.PropID] = fmt.Sprintf("%d", g.nomor("SEQ_T_CLAIM"))
		}
		for i, it := range h.AmbilDaftar(models.DaftarItem(o + 1)) {
			if it[models.PropID] == "" {
				it[models.PropID] = fmt.Sprintf("%d", g.nomor("SEQ_T_CLAIM"))
			}
			for _, a := range h.AmbilDaftar(models.DaftarAdj(o+1, i+1)) {
				if a[models.PropID] == "" {
					a[models.PropID] = fmt.Sprintf("%d", g.nomor("SEQ_T_CLAIM"))
				}
			}
		}
	}
	s := models.ProyeksiKatalog(h)
	Normalkan(s)
	pertahankanKolomKomite(lama, s)
	kr := models.SalinDaftar(lama.AmbilDaftar(models.DaftarKronologi))
	for _, b := range h.AmbilDaftar(models.DaftarKronologi) {
		if b[models.PropRiwayatBaru] == "1" {
			nb := b.Salin()
			delete(nb, models.PropRiwayatBaru)
			nb["No"] = fmt.Sprintf("%d", len(kr)+1)
			kr = append(kr, nb)
			delete(b, models.PropRiwayatBaru)
		}
	}
	s.SetelDaftar(models.DaftarKronologi, kr)
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

// SetelKasus menimpa kasus (fixture uji).
func (g *Gudang) SetelKasus(k models.Kasus) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Kasus[k.ID] = k
	if g.halaman[k.ID] == nil {
		g.halaman[k.ID] = models.HalamanBaru()
	}
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

// BacaOS - baris OS nomor klaim kasus (status 0 / 1).
func (g *Gudang) BacaOS(_ context.Context, kasusID, noKlaim string) ([]models.BarisOSTersimpan, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.BarisOSTersimpan
	for _, b := range g.OS {
		if b.NoClaim == noKlaim && b.CaseID == models.KunciInstans(kasusID) && (b.StsReject == "0" || b.StsReject == "1") {
			out = append(out, models.BarisOSTersimpan{NoClaim: b.NoClaim, NoPolis: b.NoPolis, StsReject: b.StsReject,
				DataJSON: b.DataJSON})
		}
	}
	return out, nil
}

// SalinJSONKlaim - baris IDPEGA yang ada dibiarkan.
func (g *Gudang) SalinJSONKlaim(_ context.Context, _ *db.Tx, idPega, noKlaim, noPolis string, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.JSONKlaim[idPega]; !ada {
		g.JSONKlaim[idPega] = [2]string{noKlaim, noPolis}
	}
	return nil
}

// SisipKatastrofe menyimpan katastrofe baru.
func (g *Gudang) SisipKatastrofe(_ context.Context, _ *db.Tx, k models.BarisKatastrofe, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Katastrofe = append(g.Katastrofe, k)
	return nil
}

// CatatLogLayanan menyimpan log layanan.
func (g *Gudang) CatatLogLayanan(_ context.Context, _ *db.Tx, l repository.LogLayanan, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Log = append(g.Log, l)
	return nil
}

// TulisProgres - PROGRESSCLAIM (sisip bila posisi belum ada) + SUBPROGRESSCLAIM (sisip bila belum ada Auto Create).
func (g *Gudang) TulisProgres(_ context.Context, _ *db.Tx, p models.Progres, s models.SubProgres) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	ada := false
	for _, x := range g.Progres {
		if x.IDPega == p.IDPega && x.Posisi == p.Posisi {
			ada = true
		}
	}
	if !ada {
		for i, x := range g.Progres {
			if x.IDPega == p.IDPega {
				g.Progres[i].Status = "Done"
			}
		}
		g.Progres = append(g.Progres, p)
	}
	for _, x := range g.SubProgres {
		if x.IDPega == s.IDPega && x.IDProgres == s.IDProgres && x.Jenis == s.Jenis {
			return nil
		}
	}
	g.SubProgres = append(g.SubProgres, s)
	return nil
}

// AntreEfek mencatat jenis efek.
func (g *Gudang) AntreEfek(_ context.Context, _ *db.Tx, jenis, rujukan, _ string, _ time.Time) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Efek = append(g.Efek, jenis+":"+rujukan)
	return fmt.Sprintf("EFK-%d", len(g.Efek)), nil
}

// TandaiDLAOS - InsertDLA_OS_SQL tiruan.
func (g *Gudang) TandaiDLAOS(_ context.Context, _ *db.Tx, noAkseptasi, noDLA string, _ time.Time) error {
	if noAkseptasi == "" {
		return nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.DLAOS = append(g.DLAOS, noAkseptasi+"="+noDLA)
	return nil
}

// Normalkan meniru TO_CHAR repository atas kolom katalog: angka direduksi, tanggal "2006-01-02", waktu
// "2006-01-02 15:04:05".
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
	var jalan func(t *models.Tabel, daftar string)
	jalan = func(t *models.Tabel, daftar string) {
		for n, b := range h.AmbilDaftar(daftar) {
			for _, k := range t.Kolom {
				if v := b[k.Properti]; v != "" {
					b[k.Properti] = nilai(k, v)
				}
			}
			for _, a := range t.Anak {
				jalan(a, models.JalurAnak(daftar, n+1, a.Daftar))
			}
		}
	}
	jalan(&models.TabelObjek, models.DaftarObjek)
	jalan(&models.TabelRetroTreaty, models.DaftarFacRetroTreaty)
}

// BuatKasusKomite - kelahiran kasus komite tiruan (T_WORK_CLAIM KMT- + header + tangga).
func (g *Gudang) BuatKasusKomite(ctx context.Context, tx *db.Tx, klaimID, adjID, transfer, pembuat, namaPembuat string,
	anggota []repository.AnggotaTangga, saat time.Time) (string, error) {
	if len(anggota) == 0 {
		return "", fmt.Errorf("tiruan: tangga komite kosong")
	}
	id, err := g.IDKasusBerikut(ctx, tx, models.AwalanKomite)
	if err != nil {
		return "", err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Kasus[id] = models.Kasus{ID: id, Tahap: models.TahapKomite, Posisi: anggota[0].OperatorID, PembuatID: pembuat,
		PembuatNama: namaPembuat}
	g.Komite[id] = KasusKomite{KlaimID: klaimID, AdjustmentID: adjID, Transfer: transfer, Posisi: anggota[0].OperatorID,
		Anggota: anggota}
	return id, nil
}

// AdaKomiteTutupTerbuka - kasus komite TT3 / TT4 klaim yang belum selesai.
func (g *Gudang) AdaKomiteTutupTerbuka(_ context.Context, _ *db.Tx, klaimID string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for id, k := range g.Komite {
		if k.KlaimID == klaimID && (k.Transfer == models.TransferTolak || k.Transfer == models.TransferTutup) &&
			!g.Kasus[id].Tertutup() {
			return true, nil
		}
	}
	return false, nil
}

// SetelKomiteAdjustment - KOMITE_ID baris adjustment tersimpan (UNIQUE ditiru: satu KMT satu adjustment).
func (g *Gudang) SetelKomiteAdjustment(_ context.Context, _ *db.Tx, adjID, komiteID, _ string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, h := range g.halaman {
		for jalur, rows := range h.Daftar {
			if !strings.HasSuffix(jalur, "."+models.AnakAdj) {
				continue
			}
			for _, b := range rows {
				if b[models.PropKomiteID] == komiteID && b[models.PropID] != adjID {
					return fmt.Errorf("tiruan: KOMITE_ID %s sudah tertaut (UQ_CLAIM_ADJUSTMENT_KOMITE)", komiteID)
				}
				if b[models.PropID] == adjID {
					b[models.PropKomiteID] = komiteID
					return nil
				}
			}
		}
	}
	return fmt.Errorf("tiruan: adjustment %s tidak ada", adjID)
}

// TanggaKomite - anggota tangga kasus komite tiruan (approval menunggu).
func (g *Gudang) TanggaKomite(komiteID string) []models.AnggotaKomite {
	g.mu.Lock()
	defer g.mu.Unlock()
	k, ada := g.Komite[komiteID]
	if !ada {
		return nil
	}
	var out []models.AnggotaKomite
	for _, a := range k.Anggota {
		out = append(out, models.AnggotaKomite{OperatorID: a.OperatorID, Jabatan: a.Jabatan, Email: a.Email})
	}
	return out
}

// jagaKomite meniru repository.jagaKomite: adjustment tersimpan ber-KOMITE_ID yang hilang dari halaman (dihapus langsung
// atau lewat item / objek) -> repository.ErrAdjustmentBerkomite.
func jagaKomite(lama, baru *models.Halaman) error {
	if lama == nil {
		return nil
	}
	semua := func(h *models.Halaman, f func(models.Baris)) {
		for o := range h.AmbilDaftar(models.DaftarObjek) {
			for i := range h.AmbilDaftar(models.DaftarItem(o + 1)) {
				for _, a := range h.AmbilDaftar(models.DaftarAdj(o+1, i+1)) {
					f(a)
				}
			}
		}
	}
	ada := map[string]bool{}
	semua(baru, func(a models.Baris) { ada[a[models.PropID]] = true })
	var hilang error
	semua(lama, func(a models.Baris) {
		if a[models.PropKomiteID] != "" && a[models.PropID] != "" && !ada[a[models.PropID]] {
			hilang = repository.ErrAdjustmentBerkomite
		}
	})
	return hilang
}

// pertahankanKolomKomite meniru `sqlUbahSimpul`: baris adjustment yang sudah tersimpan (UPDATE) tidak menulis kolom milik
// komite (`Kolom.MilikKomite`) - nilainya tetap nilai tersimpan; baris baru (INSERT) menulisnya.
func pertahankanKolomKomite(lama, baru *models.Halaman) {
	if lama == nil {
		return
	}
	tersimpan := map[string]models.Baris{}
	for o := range lama.AmbilDaftar(models.DaftarObjek) {
		for i := range lama.AmbilDaftar(models.DaftarItem(o + 1)) {
			for _, a := range lama.AmbilDaftar(models.DaftarAdj(o+1, i+1)) {
				tersimpan[a[models.PropID]] = a
			}
		}
	}
	for o := range baru.AmbilDaftar(models.DaftarObjek) {
		for i := range baru.AmbilDaftar(models.DaftarItem(o + 1)) {
			for _, a := range baru.AmbilDaftar(models.DaftarAdj(o+1, i+1)) {
				l, ada := tersimpan[a[models.PropID]]
				if !ada {
					continue
				}
				for _, k := range models.TabelAdjustment.Kolom {
					if !k.MilikKomite {
						continue
					}
					if v := l[k.Properti]; v != "" {
						a[k.Properti] = v
					} else {
						delete(a, k.Properti)
					}
				}
			}
		}
	}
}

// UbahAdjustmentKomite - tiruan `repository.Gudang.UbahAdjustmentKomite`: kolom milik komite baris adjustment `adjID`.
func (g *Gudang) UbahAdjustmentKomite(_ context.Context, _ *db.Tx, klaimID, adjID string, nilai map[string]string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	milik := map[string]bool{}
	for _, k := range models.TabelAdjustment.Kolom {
		if k.MilikKomite {
			milik[k.Properti] = true
		}
	}
	for p := range nilai {
		if !milik[p] {
			return fmt.Errorf("tiruan: %q bukan kolom milik komite baris adjustment", p)
		}
	}
	h := g.halaman[klaimID]
	if h == nil {
		return repository.ErrKasusTidakAda
	}
	for o := range h.AmbilDaftar(models.DaftarObjek) {
		for i := range h.AmbilDaftar(models.DaftarItem(o + 1)) {
			for _, a := range h.AmbilDaftar(models.DaftarAdj(o+1, i+1)) {
				if a[models.PropID] != adjID {
					continue
				}
				for p, v := range nilai {
					if v == "" {
						delete(a, p)
					} else {
						a[p] = v
					}
				}
				Normalkan(h)
				return nil
			}
		}
	}
	return fmt.Errorf("tiruan: baris adjustment %q klaim %q tidak ada", adjID, klaimID)
}
