package services

// Untuk apa berkas ini: LAYANAN - galat yang dipetakan handlers, wewenang pemegang assignment, pemuatan kasus beserta
// medan turunannya, pra-proses flow action, pembentukan layar, daftar kerja halaman awal, dan pembuatan kasus.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/repository"
)

// Galat yang dipetakan handlers ke kode HTTP.
var (
	ErrTanpaOracle        = db.ErrTanpaOracle
	ErrKasusTidakAda      = repository.ErrKasusTidakAda
	ErrTahapBerubah       = repository.ErrTahapBerubah
	ErrPermintaanTidakSah = galat.ErrPermintaanTidakSah
	// ErrKasusTertutup - kasus Resolved-Completed: server menolak setiap tulisan (prompt §6 butir 10).
	ErrKasusTertutup = errors.New("services: kasus sudah diselesaikan")
	// ErrBukanPemegang - pelaku bukan pemegang assignment kasus (Assignment2 = worklist pembuat, Assignment1 =
	// workbasket ReasKlaimTeknik).
	ErrBukanPemegang = fmt.Errorf("%w: bukan pemegang assignment kasus ini", inti.ErrTanpaWewenang)
	// ErrAksiTertutup - aksi tidak tampil atau nonaktif di layar kasus saat ini.
	ErrAksiTertutup = errors.New("services: aksi ini tidak tersedia di layar kasus saat ini")
	// ErrTertunda - aksi bergantung pada rule / layanan yang belum tersedia (OQ).
	ErrTertunda = errors.New("services: aksi ini tertunda menunggu keputusan (OQ)")
)

// GalatValidasi - pesan validasi layar (Property-Set-Messages / wajib). Handlers menjawabnya 422.
type GalatValidasi struct {
	Pesan []string
	// Layar - layar dengan isian dan pesannya (tidak disimpan) untuk ditampilkan.
	Layar *Layar
}

func (e *GalatValidasi) Error() string {
	return "services: validasi layar gagal: " + strings.Join(e.Pesan, "; ")
}

// Layanan adalah pintu aturan dagang Claim Prop.
type Layanan struct {
	g        Gudang
	a        Acuan
	jam      func() time.Time
	produksi bool
}

// Baru menyusun layanan; `g` nil = tanpa Oracle (setiap aksi 503).
func Baru(g Gudang, a Acuan, jam func() time.Time, produksi bool) *Layanan {
	if jam == nil {
		jam = time.Now
	}
	return &Layanan{g: g, a: a, jam: jam, produksi: produksi}
}

// AdaGudang - layanan tersambung ke penyimpanan.
func (l *Layanan) AdaGudang() bool { return l != nil && l.g != nil && l.a != nil }

func (l *Layanan) periksaPelaku(p inti.Pelaku) error {
	if !l.AdaGudang() {
		return ErrTanpaOracle
	}
	return inti.WajibIdentitas(p)
}

// Pemegang - pelaku memegang assignment kasus (tiket 10, AC 57: ditegakkan di layanan).
//
//	Assignment2 Outstanding Claim  `ToCurrentOperator` -> pembuat kasus (CREATE_OP)
//	Assignment1 Input Acceptation  `ToWorkbasket` ReasKlaimTeknik -> anggota workbasket
func Pemegang(p inti.Pelaku, k models.Kasus) bool {
	if k.Tertutup() {
		return false
	}
	switch k.Tahap {
	case models.TahapOutstanding:
		return p.AkunID != "" && strings.EqualFold(p.AkunID, k.PembuatID)
	case models.TahapAcceptation:
		return p.PunyaPeran(models.WorkbasketAcceptation)
	}
	return false
}

// konteks menyusun konteks aksi pelaku (tingkat wewenang dari roster, AC 54).
func (l *Layanan) konteks(ctx context.Context, p inti.Pelaku, saat time.Time) (*models.Konteks, error) {
	tingkat, err := l.a.TingkatPelaku(ctx, p.AkunID)
	if err != nil {
		return nil, err
	}
	return &models.Konteks{Ctx: ctx, Acuan: l.a, Pelaku: p.AkunID, Tingkat: tingkat, Sekarang: saat,
		Produksi: l.produksi}, nil
}

// ---------------------------------------------------------------- layar

// Layar - satu kasus siap ditampilkan.
type Layar struct {
	Kasus      models.Kasus    `json:"kasus"`
	Label      string          `json:"label"`
	Halaman    *models.Halaman `json:"halaman"`
	BolehKerja bool            `json:"bolehKerja"`
	// Tata - section assignment kasus (OutstandingClaim / InputAcceptation).
	Tata []models.Tata `json:"tata"`
	// Adjustment - AdjustmentDetail_Section per baris AdjustmentList (indeks 1..n).
	Adjustment map[int][]models.Tata `json:"adjustment,omitempty"`
	// Modal - local action / harness: "tutupKlaim", "pla", "dla:<n>", "komite:<n>".
	Modal map[string][]models.Tata `json:"modal,omitempty"`
	Pesan []string                 `json:"pesan,omitempty"`
	// PesanMedan - pesan per jalur (Property-Set-Messages).
	PesanMedan map[string][]string `json:"pesanMedan,omitempty"`
	// Info - pemberitahuan local action sesudah aksi (PrintFile "Please Print Pla" / PrintFileDLA).
	Info string `json:"info,omitempty"`
}

// tataKasus mengevaluasi seluruh section yang terbuka bagi kasus ini.
func tataKasus(k models.Kasus, h *models.Halaman, kunci bool) ([]models.Tata, map[int][]models.Tata, map[string][]models.Tata) {
	modal := map[string][]models.Tata{}
	adj := map[int][]models.Tata{}
	var utama []models.Tata
	switch k.Tahap {
	case models.TahapOutstanding:
		utama = models.Evaluasi(h, models.LayarOutstanding(), kunci)
		modal["pla"] = models.Evaluasi(h, models.LayarPLA(), kunci)
	default:
		utama = models.Evaluasi(h, models.LayarAkseptasi(), kunci)
		modal["pla"] = models.Evaluasi(h, models.LayarPLA(), kunci)
		modal["tutupKlaim"] = models.Evaluasi(h, models.LayarTutupKlaim(), kunci)
		lolos := h.Ambil("Protect.CARI1") == "1" && h.Ambil("Protect.CARI2") == "1"
		for i := range h.AmbilDaftar(models.DaftarAdjustment) {
			n := i + 1
			adj[n] = models.Evaluasi(h, models.LayarAdjustment(n), kunci)
			modal[fmt.Sprintf("dla:%d", n)] = models.Evaluasi(h, models.LayarDLA(n), kunci)
			modal[fmt.Sprintf("komite:%d", n)] = models.Evaluasi(h, models.LayarKomite(n, lolos), kunci)
		}
	}
	return utama, adj, modal
}

// semuaTata - gabungan seluruh tata (untuk medan terbuka dan aksi terbuka).
func semuaTata(utama []models.Tata, adj map[int][]models.Tata, modal map[string][]models.Tata) []models.Tata {
	out := append([]models.Tata{}, utama...)
	for _, t := range adj {
		out = append(out, t...)
	}
	for _, t := range modal {
		out = append(out, t...)
	}
	return out
}

func (l *Layanan) layar(k models.Kasus, h *models.Halaman, boleh bool) *Layar {
	utama, adj, modal := tataKasus(k, h, !boleh)
	ly := &Layar{Kasus: k, Label: models.LabelTahap[k.Tahap], Halaman: h, BolehKerja: boleh, Tata: utama,
		Adjustment: adj, Modal: modal, Pesan: h.SemuaPesan(), PesanMedan: h.Pesan}
	return ly
}

// ---------------------------------------------------------------- pemuatan

// muat membaca kasus dan halamannya (tx boleh nil).
func (l *Layanan) muat(ctx context.Context, tx *db.Tx, id string) (models.Kasus, *models.Halaman, error) {
	k, err := l.g.Keadaan(ctx, tx, id)
	if err != nil {
		return models.Kasus{}, nil, err
	}
	h, err := l.g.BacaHalaman(ctx, tx, id)
	if err != nil {
		return models.Kasus{}, nil, err
	}
	h.Setel("pyID", k.ID)
	return k, h, nil
}

// turunkan menyusun medan TURUNAN halaman (tidak disimpan): halaman master `TreatyInMaster` dimuat ulang dari master
// (RNMShareP tersimpan dipertahankan), totals (`HitungTurunan`), mata uang adjustment, spreading adjustment total,
// grid komite, status kasir, grid Claim History, pesan realisasi.
func (l *Layanan) turunkan(ctx context.Context, k *models.Konteks, h *models.Halaman) error {
	// `pyWorkPage.OfferFacIn.QuotationData` - salinan kerja; Class of Business disimpan sekali (BUSINESS_CODE).
	if h.Ambil(models.OQ+"BusinessCode") == "" {
		h.Setel(models.OQ+"BusinessCode", h.Ambil(models.CD+"QuotationData.BusinessCode"))
	}
	if id := h.Ambil(models.CD + "IDMaster"); id != "" {
		m, ada, err := l.a.MasterTreaty(ctx, id)
		if err != nil {
			return err
		}
		if ada {
			models.TerapkanMaster(h, m, false)
		}
	}
	if err := models.HitungTurunan(h); err != nil {
		return err
	}
	models.SusunMataUangAdjustment(h)
	if err := models.HitungSpreadingAdjustmentTotal(h); err != nil {
		return err
	}
	for i, b := range h.AmbilDaftar(models.DaftarAdjustment) {
		var tangga []models.AnggotaKomite
		if kid := b[models.PropKomiteID]; kid != "" {
			var err error
			if tangga, err = l.a.TanggaKomite(ctx, kid); err != nil {
				return err
			}
			b["KomiteNo"] = kid // SetKomiteNo_Act: nomor kasus komite (pxCoveredInsKeys)
			if tangga == nil {
				tangga = []models.AnggotaKomite{}
			}
		}
		if err := models.SusunKomiteAdjustment(k, h, i+1, tangga); err != nil {
			return err
		}
		if b["AcceptedNo"] != "" { // GetStatusKasir_Act (defer load blok kasir)
			ket, ada, err := l.a.StatusKasir(ctx, b["AcceptedNo"])
			if err != nil {
				return err
			}
			models.TerapkanStatusKasir(b, ket, ada)
		}
	}
	models.CheckAnyAcceptationProp(h)
	if h.Ambil("IsEditRNMShare") == "" {
		h.Setel("IsEditRNMShare", "false")
	}
	h.SetelDaftar(models.DaftarRiwayatTampil, models.TampilRiwayat(h.AmbilDaftar(models.DaftarRiwayat)))
	return nil
}

// siapkan = pra-proses flow action tahap kasus bagi pemegangnya:
//
//	OutstandingClaim  DT SetDateOutstanding (.StartDateEstimation), Activity CheeckNoRNM_Act
//	InputAcceptation  DT SetDateAcceptation (.EndDateEstimation, .StartDateAdjustment), Activity GetPICAdjutment_Act
//
// ⚠️ Pega menjalankan DT pra-proses setiap assignment dibuka; di sini tanggalnya hanya diisi bila masih kosong
// (penyimpangan sadar - jejak waktu mulai tahap tidak bergeser setiap kali berkas dibuka; PARITAS).
func (l *Layanan) siapkan(k *models.Konteks, kasus models.Kasus, h *models.Halaman) error {
	switch kasus.Tahap {
	case models.TahapOutstanding:
		if h.Ambil("StartDateEstimation") == "" {
			h.Setel("StartDateEstimation", k.Waktu())
		}
		return models.CheeckNoRNM(k, h)
	case models.TahapAcceptation:
		if h.Ambil("EndDateEstimation") == "" {
			h.Setel("EndDateEstimation", k.Waktu())
		}
		if h.Ambil("StartDateAdjustment") == "" {
			h.Setel("StartDateAdjustment", k.Waktu())
		}
		return models.GetPICAdjutment(k, h)
	}
	return nil
}

// ---------------------------------------------------------------- buka, daftar, buat

// BukaKasus membuka satu kasus. Pemegang assignment: pra-proses dijalankan (pesannya tampil, tidak disimpan); selainnya
// hanya-baca.
func (l *Layanan) BukaKasus(ctx context.Context, p inti.Pelaku, id string) (*Layar, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	k, h, err := l.muat(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	kt, err := l.konteks(ctx, p, l.jam())
	if err != nil {
		return nil, err
	}
	boleh := Pemegang(p, k)
	if boleh {
		if err := l.siapkan(kt, k, h); err != nil {
			return nil, err
		}
	}
	if err := l.turunkan(ctx, kt, h); err != nil {
		return nil, err
	}
	return l.layar(k, h, boleh), nil
}

// JenisDaftar - tab halaman awal.
const (
	DaftarSaya       = "saya"       // Assignment2 Outstanding Claim - worklist pembuat (ToCurrentOperator)
	DaftarWorkbasket = "workbasket" // Assignment1 Input Acceptation - workbasket ReasKlaimTeknik
	DaftarSelesai    = "selesai"    // Resolved-Completed
)

// DaftarKasus - halaman awal: pintu masuk Pega yang berbukti (`Flow_TreatyIn` Assignment2 ke worklist operator,
// Assignment1 ke workbasket).
func (l *Layanan) DaftarKasus(ctx context.Context, p inti.Pelaku, jenis, cari string) ([]repository.RingkasanKasus, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	s := repository.SaringanKasus{Cari: cari}
	switch jenis {
	case DaftarSaya, "":
		s.Tahap, s.Pembuat = models.TahapOutstanding, p.AkunID
	case DaftarWorkbasket:
		if !p.PunyaPeran(models.WorkbasketAcceptation) {
			return []repository.RingkasanKasus{}, nil
		}
		s.Tahap = models.TahapAcceptation
	case DaftarSelesai:
		s.Selesai = true
	default:
		return nil, fmt.Errorf("%w: daftar %q", ErrPermintaanTidakSah, jenis)
	}
	return l.g.DaftarKasus(ctx, s)
}

// BuatKasus = Start1 -> Assignment2 (`Flow_TreatyIn` Transition3): kasus CLMP- baru di worklist pembuatnya.
// ⚠️ Harness New / NewSample (pintu pembuatan) tidak diekspor (`pyCanCreateWorkObject=false`) - tombol "New" halaman
// awal = OQ-CP-13 (label dan letak pintu pembuatan Pega).
func (l *Layanan) BuatKasus(ctx context.Context, p inti.Pelaku) (*Layar, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	nama, err := l.a.NamaPelaku(ctx, p.AkunID)
	if err != nil {
		return nil, err
	}
	saat := l.jam()
	var id string
	err = l.g.Transaksi(ctx, func(tx *db.Tx) error {
		var err error
		if id, err = l.g.IDKasusBerikut(ctx, tx, models.AwalanKlaim); err != nil {
			return err
		}
		return l.g.SisipKasus(ctx, tx, id, p.AkunID, nama, saat)
	})
	if err != nil {
		return nil, err
	}
	return l.BukaKasus(ctx, p, id)
}

// AcuanStatis - daftar pilihan bersama layar.
type AcuanStatis struct {
	MataUang   []models.Pilihan    `json:"mataUang"`
	JenisReas  []models.Pilihan    `json:"jenisReas"`
	JenisReas4 []models.Pilihan    `json:"jenisReas4"`
	Kode       map[string][]string `json:"kode"`
}

// Acuan membaca daftar pilihan bersama.
func (l *Layanan) Acuan(ctx context.Context, p inti.Pelaku) (AcuanStatis, error) {
	if err := l.periksaPelaku(p); err != nil {
		return AcuanStatis{}, err
	}
	var out AcuanStatis
	var err error
	if out.MataUang, err = l.a.DaftarMataUang(ctx); err != nil {
		return out, err
	}
	if out.JenisReas, err = l.a.DaftarJenisReas(ctx, ""); err != nil {
		return out, err
	}
	if out.JenisReas4, err = l.a.DaftarJenisReas(ctx, models.TipeJenisLossAllocation); err != nil {
		return out, err
	}
	out.Kode = models.KodePilihan
	return out, nil
}
