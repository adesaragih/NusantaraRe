package services

// Untuk apa berkas ini: PEMUAT DATA LAMA - kasus Claim Prop warisan Pega (OS_AKSEPTASI_KLAIM + JSON_KLAIM) ke tabel
// Claim Prop lewat penyimpanan yang sama dengan aplikasi (prompt implementasi §6 butir 11; AC 9-12, 123, 132). Alat
// baris perintahnya `backend/alat/pemuatlama`.
//
// Bawaan = uji-kering: hanya MEMBACA, nol pernyataan tulis, laporan dicetak. `tulis` = satu transaksi per kasus; kasus
// yang ID-nya sudah ada dilewati (aman diulang). Kasus bergalat tidak dimuat dan tidak menghentikan kasus lain.

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/repository"
)

// GudangPemuat - bacaan sumber lama dan tulisan kasus yang dipakai pemuat.
type GudangPemuat interface {
	Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error
	BacaOSLama(ctx context.Context) ([]models.BarisOSLama, error)
	BacaJSONKlaimLama(ctx context.Context) (map[string]string, error)
	Keadaan(ctx context.Context, tx *db.Tx, id string) (models.Kasus, error)
	SisipKasusLama(ctx context.Context, tx *db.Tx, k models.Kasus) error
	SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error
}

// Pemuat menjalankan pemuatan kasus lama.
type Pemuat struct {
	g GudangPemuat
	// nama - nama operator (M_LOGIN_GO); nil = tanpa nama.
	nama func(ctx context.Context, akun string) (string, error)
}

// PemuatBaru menyusun pemuat. `nama` boleh nil.
func PemuatBaru(g GudangPemuat, nama func(ctx context.Context, akun string) (string, error)) *Pemuat {
	return &Pemuat{g: g, nama: nama}
}

// PemuatDariDasar - pemuat di atas basis data bersama (transaksi `inti.Dasar`, sama dengan layanan).
func PemuatDariDasar(d *inti.Dasar) (*Pemuat, error) {
	if d == nil || !d.PunyaDatabase() {
		return nil, ErrTanpaOracle
	}
	g := repository.Baru(d.DB())
	return PemuatBaru(penyimpanOracle{Gudang: g, dasar: d}, repository.AcuanDari(g, false).NamaPelaku), nil
}

// Jalankan memuat seluruh kasus lama. Galat per kasus dicatat di laporan; galat baca sumber menghentikan jalannya.
func (p *Pemuat) Jalankan(ctx context.Context, tulis bool, lap *LaporanPemuat) error {
	rows, err := p.g.BacaOSLama(ctx)
	if err != nil {
		return err
	}
	dok, err := p.g.BacaJSONKlaimLama(ctx)
	if err != nil {
		return err
	}
	kasus := map[string][]models.BarisOSLama{}
	var urut []string
	for _, b := range rows {
		if _, ada := kasus[b.CaseID]; !ada {
			urut = append(urut, b.CaseID)
		}
		kasus[b.CaseID] = append(kasus[b.CaseID], b)
	}
	sort.Strings(urut)
	lap.r.Baris = len(rows)
	for _, kunci := range urut {
		if err := ctx.Err(); err != nil {
			return err
		}
		lap.r.Kasus++
		if err := p.satu(ctx, tulis, kunci, kasus[kunci], dok, lap); err != nil {
			return err
		}
	}
	return nil
}

// satu memuat satu kasus. Hanya galat basis data selain milik kasus ini yang dikembalikan.
func (p *Pemuat) satu(ctx context.Context, tulis bool, kunci string, rows []models.BarisOSLama, dok map[string]string,
	lap *LaporanPemuat) error {
	id := models.IDDariKunciPega(kunci)
	berlaku, riwayat := models.BarisBerlaku(rows)
	lap.r.Riwayat += riwayat
	teks, berJSON := dok[kunci]
	if berJSON {
		lap.r.DariJSON++
	}
	tahap, tutup, ok := models.TahapLama(berlaku.StsReject)
	switch {
	case !ok:
		lap.galat(id, "STS_REJECT", berlaku.StsReject, models.SebabStsTakDikenal)
		return nil
	}
	var (
		h       *models.Halaman
		dibuang []models.MedanDibuang
		err     error
	)
	if berJSON {
		h, dibuang, err = models.UraiKlaimLama(teks)
	} else {
		h, dibuang, err = models.HalamanDariOSLama(berlaku)
	}
	if err != nil {
		lap.galat(id, "DATA_JSON", "", "galat: "+err.Error())
		return nil
	}
	for _, m := range dibuang {
		lap.buang(id, m)
	}
	if g := models.GalatNilaiKatalog(h); len(g) > 0 {
		for _, m := range g {
			lap.galat(id, m.Jalur, m.Nilai, m.Sebab)
		}
		return nil
	}
	urutan := append([]models.BarisOSLama(nil), rows...)
	models.UrutkanOSLama(urutan)
	k := models.Kasus{ID: id, Tahap: tahap, PembuatID: urutan[0].InsertOp, TglCreate: urutan[0].Tanggal,
		TglUpdate: berlaku.Tanggal, Sumber: models.SumberPega, Posisi: urutan[0].InsertOp}
	if tahap == models.TahapAcceptation {
		k.Posisi = models.WorkbasketAcceptation
	}
	kunciTahap := tahap
	if tutup {
		k.StatusWork = models.StatusSelesai
		kunciTahap = models.StatusSelesai
	}
	lap.r.PerTahap[kunciTahap]++
	lap.r.Siap++
	if !tulis {
		return nil
	}
	if p.nama != nil && k.PembuatID != "" {
		if n, err := p.nama(ctx, k.PembuatID); err == nil {
			k.PembuatNama = n
		}
	}
	dimuat := false
	err = p.g.Transaksi(ctx, func(tx *db.Tx) error {
		if _, err := p.g.Keadaan(ctx, tx, id); err == nil {
			return nil
		} else if !errors.Is(err, repository.ErrKasusTidakAda) {
			return err
		}
		if err := p.g.SisipKasusLama(ctx, tx, k); err != nil {
			return err
		}
		dimuat = true
		return p.g.SimpanHalaman(ctx, tx, id, h)
	})
	switch {
	case err != nil:
		lap.galat(id, "", "", "galat tulis: "+err.Error())
	case dimuat:
		lap.r.Dimuat++
	default:
		lap.r.Dilewati++
	}
	return nil
}

// RingkasanPemuat - cacah satu jalannya pemuat.
type RingkasanPemuat struct {
	Tulis bool
	// Kasus - kasus lama (CASEID berbeda); Baris - baris OS_AKSEPTASI_KLAIM; Riwayat - baris selain yang berlaku.
	Kasus, Baris, Riwayat int
	// DariJSON - kasus berhalaman JSON_KLAIM; Siap - kasus lolos periksa (uji-kering: yang akan dimuat).
	DariJSON, Siap          int
	Dimuat, Dilewati, Gagal int
	// PerTahap - kasus siap per tahap (Resolved-Completed untuk kasus tutup berkas).
	PerTahap map[string]int
	// PerSebab - medan dibuang per sebab.
	PerSebab map[string]int
}

// Selesai - nol kasus gagal (kasus ditunda tidak menahan).
func (r RingkasanPemuat) Selesai() bool { return r.Gagal == 0 }

// Teks - ringkasan untuk layar (cacah saja, nol data kasus).
func (r RingkasanPemuat) Teks() string {
	var b strings.Builder
	mode := "uji-kering (baca saja)"
	if r.Tulis {
		mode = "jalankan (tulis)"
	}
	fmt.Fprintf(&b, "Pemuat data lama Claim Prop - %s\n", mode)
	fmt.Fprintf(&b, "Kasus lama          : %d\nBaris OS            : %d\nBaris berlaku       : %d\nBaris riwayat       : %d\n",
		r.Kasus, r.Baris, r.Kasus, r.Riwayat)
	fmt.Fprintf(&b, "Berhalaman JSON     : %d\nSiap dimuat         : %d\nDimuat              : %d\nDilewati (sudah ada): %d\n",
		r.DariJSON, r.Siap, r.Dimuat, r.Dilewati)
	fmt.Fprintf(&b, "Kasus gagal         : %d\n", r.Gagal)
	for _, k := range kunciPeta(r.PerTahap) {
		fmt.Fprintf(&b, "  tahap %-20s: %d\n", k, r.PerTahap[k])
	}
	for _, k := range kunciPeta(r.PerSebab) {
		fmt.Fprintf(&b, "  medan %s: %d\n", k, r.PerSebab[k])
	}
	return b.String()
}

func kunciPeta(m map[string]int) []string {
	k := make([]string, 0, len(m))
	for x := range m {
		k = append(k, x)
	}
	sort.Strings(k)
	return k
}

// LaporanPemuat - dua berkas CSV (arsip medan dibuang; galat) dan ringkasan cacah.
type LaporanPemuat struct {
	r            RingkasanPemuat
	arsip, salah *csv.Writer
	gagal        map[string]bool
}

// LaporanPemuatBaru menulis kepala kedua berkas.
func LaporanPemuatBaru(arsip, galat io.Writer, tulis bool) *LaporanPemuat {
	l := &LaporanPemuat{r: RingkasanPemuat{Tulis: tulis, PerTahap: map[string]int{}, PerSebab: map[string]int{}},
		arsip: csv.NewWriter(arsip), salah: csv.NewWriter(galat), gagal: map[string]bool{}}
	_ = l.arsip.Write([]string{"KASUS", "JALUR", "NILAI", "SEBAB"})
	_ = l.salah.Write([]string{"KASUS", "JALUR", "NILAI", "SEBAB"})
	return l
}

func (l *LaporanPemuat) buang(id string, m models.MedanDibuang) {
	l.r.PerSebab[m.Sebab]++
	_ = l.arsip.Write([]string{id, m.Jalur, m.Nilai, m.Sebab})
}

func (l *LaporanPemuat) galat(id, jalur, nilai, sebab string) {
	if !l.gagal[id] {
		l.gagal[id] = true
		l.r.Gagal++
	}
	_ = l.salah.Write([]string{id, jalur, nilai, sebab})
}

// Tutup mengosongkan penyangga kedua berkas.
func (l *LaporanPemuat) Tutup() error {
	l.arsip.Flush()
	l.salah.Flush()
	return errors.Join(l.arsip.Error(), l.salah.Error())
}

// Ringkasan - cacah jalannya.
func (l *LaporanPemuat) Ringkasan() RingkasanPemuat { return l.r }
