package services

// Untuk apa berkas ini: LAYANAN dan bantuan bersamanya - galat yang dipetakan
// handlers, keanggotaan antrean, pra-proses flow action, dan pemuatan kasus.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

// Galat yang dipetakan handlers ke kode HTTP. Galat repository diteruskan
// lewat nama di sini supaya handlers tidak mengimpor repository.
var (
	ErrTanpaOracle         = db.ErrTanpaOracle
	ErrKasusTidakAda       = repository.ErrKasusTidakAda
	ErrTahapBerubah        = repository.ErrTahapBerubah
	ErrGenerasiTertutup    = repository.ErrGenerasiTertutup
	ErrNomorPolisSudahAda  = repository.ErrNomorPolisSudahAda
	ErrDataKontrakTidakAda = repository.ErrDataKontrakTidakAda
	ErrTipeNomorKosong     = repository.ErrTipeNomorKosong
	ErrOJKKosong           = repository.ErrOJKKosong
	ErrPermintaanTidakSah  = galat.ErrPermintaanTidakSah

	// ErrKasusTertutup - kasus sudah diselesaikan; tidak ada tindakan lagi.
	ErrKasusTertutup = errors.New("services: kasus sudah diselesaikan")
	// ErrBukanAnggotaAntrean - pelaku bukan anggota antrean tempat kasus
	// menunggu (AC 14: menurut NAMA antrean, bukan nomor urut).
	ErrBukanAnggotaAntrean = fmt.Errorf("%w: bukan anggota antrean tempat berkas menunggu", inti.ErrTanpaWewenang)
	// ErrTindakanTakAdaDiPosisi - tindakan itu tidak punya tombol di layar
	// posisi kasus (mis. Save di layar atasan, nomor polis di layar admin).
	ErrTindakanTakAdaDiPosisi = errors.New("services: tindakan ini tidak tersedia di posisi berkas")
)

// GalatValidasi - pesan validasi layar (medan wajib, pesan Property-Set-
// Messages). Handlers menjawabnya 422 beserta daftar pesannya.
type GalatValidasi struct{ Pesan []string }

func (e *GalatValidasi) Error() string {
	return "services: validasi layar gagal: " + strings.Join(e.Pesan, "; ")
}

// Layanan adalah pintu aturan dagang NB Treaty In.
type Layanan struct {
	g   Gudang
	jam func() time.Time
	// konversi dan produksi - efek keluar sesudah selesai (konversi.go).
	konversi PengirimKonversi
	produksi bool
}

// Baru menyusun layanan; `g` nil = tanpa Oracle (setiap tindakan 503).
func Baru(g Gudang, jam func() time.Time) *Layanan {
	if jam == nil {
		jam = time.Now
	}
	return &Layanan{g: g, jam: jam}
}

// AdaGudang - layanan tersambung ke penyimpanan.
func (l *Layanan) AdaGudang() bool { return l != nil && l.g != nil }

// periksaPelaku - layanan tersambung dan permintaan membawa identitas.
func (l *Layanan) periksaPelaku(p inti.Pelaku) error {
	if !l.AdaGudang() {
		return ErrTanpaOracle
	}
	return inti.WajibIdentitas(p)
}

// anggota menjawab: pelaku anggota antrean (workbasket) `posisi`.
func anggota(p inti.Pelaku, posisi string) bool { return posisi != "" && p.PunyaPeran(posisi) }

// ------------------------------------------------------------------ daftar dan buat

// DaftarKasus - daftar portal (`Section/SFAPortal_OpportunitiesList`).
func (l *Layanan) DaftarKasus(ctx context.Context, p inti.Pelaku, s models.SaringanKasus) ([]models.RingkasanKasus, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	if s.Posisi != "" && !models.AdalahPosisiTangga(s.Posisi) {
		return nil, fmt.Errorf("%w: posisi %q", ErrPermintaanTidakSah, s.Posisi)
	}
	return l.g.DaftarKasus(ctx, s)
}

// BuatKasus = tombol "Create opportunity" portal (`createWork` ->
// `Flow/InputRealizationTreatyIn`, connector Start1 -> Assignment2).
//
// ⛔ Tanpa gerbang peran: syarat tombolnya (`crmCreateOpportunity`,
// `isSellingMode*`, `IsNotAdmin`) milik CRM dan tidak memuat antrean NB
// Treaty In - mengarang gerbang berarti memutuskan siapa boleh membuat
// realisasi (pola premiumlistlife).
func (l *Layanan) BuatKasus(ctx context.Context, p inti.Pelaku) (models.Kasus, error) {
	if err := l.periksaPelaku(p); err != nil {
		return models.Kasus{}, err
	}
	nama, err := l.g.NamaTampilan(ctx, p.AkunID)
	if err != nil {
		return models.Kasus{}, err
	}
	var id string
	err = l.g.Transaksi(ctx, func(tx *db.Tx) error {
		var err error
		if id, err = l.g.IDKasusBerikut(ctx, tx); err != nil {
			return err
		}
		if err := l.g.SisipKasus(ctx, tx, id, p.AkunID, nama); err != nil {
			return err
		}
		h := models.HalamanBaru()
		h.Setel("PositionNote", models.PosisiAdmin)
		h.Setel(models.HalamanQuotation+".BusinessFac", models.BisnisTreaty)
		return l.g.SimpanHalaman(ctx, tx, id, h)
	})
	if err != nil {
		return models.Kasus{}, err
	}
	return l.g.Keadaan(ctx, nil, id)
}

// ------------------------------------------------------------------ buka

// Layar adalah satu kasus siap ditampilkan.
type Layar struct {
	Kasus   models.Kasus    `json:"kasus"`
	Halaman *models.Halaman `json:"halaman"`
	// BolehKerja - pelaku anggota antrean posisi kasus dan kasus terbuka;
	// selain itu layar hanya-baca.
	BolehKerja bool               `json:"bolehKerja"`
	Tombol     models.TombolKirim `json:"tombol"`
	// MedanWajib - jalur medan wajib layar posisi ini (berlaku saat ini).
	MedanWajib []string `json:"medanWajib"`
	// Tempat - tempat berperan (tiket 05) -> tampil atau tidak.
	Tempat map[string]bool `json:"tempat"`
	Pesan  []string        `json:"pesan,omitempty"`
}

// BukaKasus membuka satu kasus. Bila pelaku boleh bekerja, pra-proses flow
// action posisinya dijalankan atas halaman (tidak disimpan - di Pega
// pra-proses hanya mengubah clipboard).
func (l *Layanan) BukaKasus(ctx context.Context, p inti.Pelaku, id string) (Layar, error) {
	if err := l.periksaPelaku(p); err != nil {
		return Layar{}, err
	}
	k, h, err := l.muat(ctx, id)
	if err != nil {
		return Layar{}, err
	}
	boleh := !k.Tertutup() && anggota(p, k.PositionNote)
	if boleh {
		if err := l.siapkan(ctx, p, k, h); err != nil {
			return Layar{}, err
		}
	}
	return l.layar(ctx, p, k, h, boleh)
}

func (l *Layanan) layar(ctx context.Context, p inti.Pelaku, k models.Kasus, h *models.Halaman, boleh bool) (Layar, error) {
	tempat := l.tempat(p)
	var wajib []string
	for _, m := range models.DaftarMedanWajib(k.PositionNote) {
		if m.Syarat == nil || m.Syarat(h) {
			wajib = append(wajib, m.Jalur)
		}
	}
	if tempat[TempatTanggalProduksi] && h.Ambil(models.HalamanPolis+".IsApproved") == "1" {
		wajib = append(wajib, models.HalamanPolis+".ProductionDate")
	}
	ly := Layar{Kasus: k, Halaman: h, BolehKerja: boleh, MedanWajib: wajib, Tempat: tempat, Pesan: h.SemuaPesan()}
	if boleh {
		ly.Tombol = models.TombolUntuk(h, k.PositionNote)
	}
	return ly, nil
}

// muat membaca keadaan dan halaman tersimpan satu kasus.
func (l *Layanan) muat(ctx context.Context, id string) (models.Kasus, *models.Halaman, error) {
	k, err := l.g.Keadaan(ctx, nil, id)
	if err != nil {
		return models.Kasus{}, nil, err
	}
	h, err := l.g.BacaHalaman(ctx, nil, id)
	if err != nil {
		return models.Kasus{}, nil, err
	}
	h.Setel("Position", k.Position)
	if k.PositionNote != "" {
		h.Setel("PositionNote", k.PositionNote)
	}
	// Total spreading TURUNAN baris, tidak disimpan (models/katalog.go).
	if err := models.HitungTotalSpreading(h); err != nil {
		return models.Kasus{}, nil, err
	}
	return k, h, nil
}

// siapkan = pra-proses flow action posisi kasus:
//
//	Admin   `InputPolicyTreatyIn_preDT`, lalu `InputPolicyTreatyInPre_Act`
//	Atasan  `DeptHeadTreatyInUW_preDT`,  lalu `InputPolicyTreatyInPre_Act`
//
// Urutan Pega: data transform pra-proses lebih dulu, aktivitas pra-proses
// sesudahnya. Bagian `InputPolicyTreatyInPre_Act` yang dibangun: langkah 2
// (bisnis bila BizCode kosong), 3-4 dan 9 (tanggal; hari tutup buku dari
// TANGGAL_CLOSING - lihat `models.GeserTanggalProduksi`), dan 10
// (`TreatyRealizationCheckXOLList` - RALAT K8, `siapkanNonProp`). Langkah 1
// (`SetCategoryAttach` - lampiran, tidak ada di layar realisasi) tidak
// dibangun; langkah 5-8 menyiapkan daftar pilihan spreading (`Acuan`).
//
// Halaman master `TreatyIn` dimuat ulang dari view lewat `TreatyIn.ID` (P29).
func (l *Layanan) siapkan(ctx context.Context, p inti.Pelaku, k models.Kasus, h *models.Halaman) error {
	nama, err := l.g.NamaTampilan(ctx, p.AkunID)
	if err != nil {
		return err
	}
	sekarang := l.jam()
	if k.PositionNote == models.PosisiAdmin {
		models.PraprosesAdmin(h, sekarang, nama)
	} else {
		models.PraprosesAtasan(h, sekarang, nama)
	}
	if h.Ambil(models.HalamanPolis+".BizCode") == "" && h.Ambil(models.HalamanQuotation+".BusinessName") != "" {
		b, err := l.g.BisnisDariKunci(ctx, models.KunciCariBisnis(h.Ambil(models.HalamanQuotation+".BusinessName"), true))
		if err != nil {
			return err
		}
		models.TerapkanBisnisPra(h, b)
	}
	var hari int
	if err := l.g.Transaksi(ctx, func(tx *db.Tx) error {
		var err error
		hari, err = l.g.HariClosing(ctx, tx)
		return err
	}); err != nil {
		return err
	}
	models.PraprosesTanggal(h, sekarang, hari)
	if err := l.muatMaster(ctx, h); err != nil {
		return err
	}
	// langkah 10 dan master XOL (K8) - nonprop.go
	return l.siapkanNonProp(ctx, h)
}

// muatMaster mengisi halaman TreatyIn dari baris view kontrak terpilih.
func (l *Layanan) muatMaster(ctx context.Context, h *models.Halaman) error {
	id := h.Ambil(models.HalamanMaster + ".ID")
	if id == "" {
		return nil
	}
	b, err := l.g.DetailKontrak(ctx, id)
	if err != nil {
		return err
	}
	models.TerapkanMasterKontrak(h, b)
	return nil
}

// RiwayatKasus - riwayat akseptasi satu kasus, berurut waktu (AC 72).
func (l *Layanan) RiwayatKasus(ctx context.Context, p inti.Pelaku, id string) ([]models.Riwayat, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	if _, err := l.g.Keadaan(ctx, nil, id); err != nil {
		return nil, err
	}
	return l.g.DaftarRiwayat(ctx, models.KunciInstans(id))
}

// Acuan adalah daftar pilihan layar.
type Acuan struct {
	MataUang  []models.Pilihan `json:"mataUang"`
	MO        []models.Pilihan `json:"mo"`
	Spreading []models.Pilihan `json:"spreading"`
	JenisReas []models.Pilihan `json:"jenisReas"`
}

// DaftarAcuan membaca keempat daftar pilihan layar.
func (l *Layanan) DaftarAcuan(ctx context.Context, p inti.Pelaku) (Acuan, error) {
	if err := l.periksaPelaku(p); err != nil {
		return Acuan{}, err
	}
	var a Acuan
	var err error
	if a.MataUang, err = l.g.DaftarMataUang(ctx); err != nil {
		return Acuan{}, err
	}
	if a.MO, err = l.g.DaftarMO(ctx); err != nil {
		return Acuan{}, err
	}
	if a.Spreading, err = l.g.DaftarJenisSpreading(ctx); err != nil {
		return Acuan{}, err
	}
	if a.JenisReas, err = l.g.DaftarJenisReas(ctx); err != nil {
		return Acuan{}, err
	}
	return a, nil
}

// DaftarBisnis - grid popup `BusinessAndSOBList` (RD `BrowseTreatyInDetail`).
func (l *Layanan) DaftarBisnis(ctx context.Context, p inti.Pelaku, cari string) ([]models.BarisKontrak, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	return l.g.DaftarDetailKontrak(ctx, repository.SaringanDetail{TreatyID: strings.TrimSpace(cari)})
}
