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
	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/repository"
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
	// ErrAkseptasiDiKomite - tulisan ke akseptasi yang sedang di komite (prompt §6 butir 10).
	ErrAkseptasiDiKomite = errors.New("services: akseptasi ini sedang di komite")
	// ErrTertunda - aksi bergantung pada rule / objek yang belum tersedia (OQ).
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

// Layanan adalah pintu aturan dagang Claim Non Prop.
type Layanan struct {
	g        Gudang
	a        Acuan
	jam      func() time.Time
	produksi bool
	kasir    models.KonfigKasir
	berkas   PenyimpananBerkas
}

// Baru menyusun layanan; `g` nil = tanpa Oracle (setiap aksi 503).
func Baru(g Gudang, a Acuan, jam func() time.Time, produksi bool) *Layanan {
	if jam == nil {
		jam = time.Now
	}
	return &Layanan{g: g, a: a, jam: jam, produksi: produksi}
}

// DenganKasir memasang kode tetap muatan Kasir (`konfigurasi/kasir.json`).
func (l *Layanan) DenganKasir(k models.KonfigKasir) *Layanan { l.kasir = k; return l }

// AdaGudang - layanan tersambung ke penyimpanan.
func (l *Layanan) AdaGudang() bool { return l != nil && l.g != nil && l.a != nil }

func (l *Layanan) periksaPelaku(p inti.Pelaku) error {
	if !l.AdaGudang() {
		return ErrTanpaOracle
	}
	return inti.WajibIdentitas(p)
}

// Pemegang - pelaku memegang assignment kasus (ditegakkan di layanan, prompt §6 butir 10).
//
//	Assignment2 Outstanding Claim  `ToCurrentOperator` -> pembuat kasus (CREATE_OP)
//	Assignment1 Input Acceptation  `ToWorkbasket` ReasKlaimTeknik (OQ-CNP-08) -> anggota workbasket
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

// konteks menyusun konteks aksi pelaku (tingkat wewenang dari roster).
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
	// Adjustment - AdjustmentDetailNP_Section per baris AdjustmentList (indeks 1..n).
	Adjustment map[int][]models.Tata `json:"adjustment,omitempty"`
	// Modal - local action / harness: "tutupKlaim", "cwp", "pla", "komite:<n>", "interest:<i>",
	// "reinstatement:<n>:<i>".
	Modal map[string][]models.Tata `json:"modal,omitempty"`
	Pesan []string                 `json:"pesan,omitempty"`
	// PesanMedan - pesan per jalur (Property-Set-Messages).
	PesanMedan map[string][]string `json:"pesanMedan,omitempty"`
	// Info - pemberitahuan sesudah aksi (mis. berkas dokumen menunggu OQ-CNP-22).
	Info string `json:"info,omitempty"`
	// BukaModal - modal yang dibuka layar sesudah aksi (harness KomiteCNP, local action CloseClaimNP).
	BukaModal string `json:"bukaModal,omitempty"`
	// Mode - penanda mode layar (models.ModeLayar); layar mengembalikannya di setiap aksi.
	Mode map[string]string `json:"mode,omitempty"`
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
		for i := range h.AmbilDaftar(models.DaftarInterest) {
			modal[fmt.Sprintf("interest:%d", i+1)] = models.Evaluasi(h, models.LayarInterest(i+1), kunci)
		}
	default:
		utama = models.Evaluasi(h, models.LayarAkseptasi(), kunci)
		modal["tutupKlaim"] = models.Evaluasi(h, models.LayarTutupKlaim(), kunci)
		modal["cwp"] = models.Evaluasi(h, models.LayarCWP(), kunci)
		for i := range h.AmbilDaftar(models.DaftarAdjustment) {
			n := i + 1
			adj[n] = models.Evaluasi(h, models.LayarDetailAkseptasi(n), kunci)
			modal[fmt.Sprintf("komite:%d", n)] = models.Evaluasi(h, models.LayarKomite(n), kunci)
			for x := range h.AmbilDaftar(models.JalurAdj(n, models.AnakXOL)) {
				modal[fmt.Sprintf("reinstatement:%d:%d", n, x+1)] = models.Evaluasi(h, models.LayarReinstatement(n, x+1), kunci)
			}
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
	return &Layar{Kasus: k, Label: models.LabelTahap[k.Tahap], Halaman: h, BolehKerja: boleh, Tata: utama,
		Adjustment: adj, Modal: modal, Pesan: h.SemuaPesan(), PesanMedan: h.Pesan, Mode: models.AmbilMode(h)}
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

// master membaca master treaty kasus (kosong bila belum dipilih).
func (l *Layanan) master(ctx context.Context, h *models.Halaman) (models.MasterTreaty, error) {
	id := h.Ambil(models.CD + "IDMaster")
	if id == "" {
		return models.MasterTreaty{}, nil
	}
	m, _, err := l.a.MasterTreaty(ctx, id)
	return m, err
}

// turunkan menyusun medan TURUNAN halaman (tidak disimpan): halaman master `TreatyInMaster` dimuat ulang dari master
// (RNM Share dan Treaty Group tersimpan dipertahankan), `HitungTurunan`, grid komite akseptasi (tangga tersimpan atau
// calon roster menurut OQ-CNP-01), status kasir, grid Claim History.
func (l *Layanan) turunkan(ctx context.Context, h *models.Halaman) error {
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
	for i, b := range h.AmbilDaftar(models.DaftarAdjustment) {
		n := i + 1
		var tangga, calon []models.AnggotaKomite
		var err error
		if kid := b[models.PropKomiteID]; kid != "" {
			if tangga, err = l.a.TanggaKomite(ctx, kid); err != nil {
				return err
			}
			if tangga == nil {
				tangga = []models.AnggotaKomite{}
			}
			b["KomiteNo"] = kid // CreateChildKomiteCNP_Act 31: nomor kasus komite
		} else if calon, err = l.a.RosterKomite(ctx, models.HanyaTingkat1(h, n)); err != nil {
			return err
		}
		if err := models.SusunKomiteAkseptasi(h, n, tangga, calon); err != nil {
			return err
		}
		if b["AcceptedNo"] != "" {
			ket, ada, err := l.a.StatusKasir(ctx, b["AcceptedNo"])
			if err != nil {
				return err
			}
			if ada {
				b["StatusKasir"] = ket
			}
		}
	}
	h.SetelDaftar(models.DaftarRiwayatTampil, models.TampilRiwayat(h.AmbilDaftar(models.DaftarRiwayat)))
	return nil
}

// siapkan = pra-proses flow action tahap kasus bagi pemegangnya:
//
//	OutstandingClaim  InputOutStandingClmTNP_PreAct (daftar polis / treaty dibaca saat dropdown dibuka; MO)
//	InputAcceptation  InputAkseptasi_PreAct (CommentLOD, IsAcceptation, pra Outstanding, mesin XoL bila kosong)
func (l *Layanan) siapkan(ctx context.Context, k *models.Konteks, kasus models.Kasus, h *models.Halaman) error {
	switch kasus.Tahap {
	case models.TahapOutstanding:
		return models.PraOutstanding(k, h)
	case models.TahapAcceptation:
		m, err := l.master(ctx, h)
		if err != nil {
			return err
		}
		return models.PraAkseptasi(k, h, m)
	}
	return nil
}

// ---------------------------------------------------------------- buka, daftar, buat

// BukaKasus membuka satu kasus. Pemegang assignment: pra-proses dijalankan (tidak disimpan); selainnya hanya-baca.
func (l *Layanan) BukaKasus(ctx context.Context, p inti.Pelaku, id string, lihat bool) (*Layar, error) {
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
	boleh := Pemegang(p, k) && !lihat
	if boleh {
		if err := l.siapkan(ctx, kt, k, h); err != nil {
			return nil, err
		}
		h.BersihkanPesan()
	}
	if err := l.turunkan(ctx, h); err != nil {
		return nil, err
	}
	return l.layar(k, h, boleh), nil
}

// HakPelaku - hak halaman awal: switch Teknik aktif hanya bagi anggota workbasket Assignment1 (pola Claim Prop).
// Komite - tabel komite di bawah inbox (menu Komite Claim Non Prop disembunyikan, perintah work owner 09-10-2026, pola
// Claim Prop): akun memegang workbasket yang tercantum sebagai KomiteID roster EMAILKOMITE NONPROP aktif.
type HakPelaku struct {
	WorkbasketTeknik bool `json:"workbasketTeknik"`
	Komite           bool `json:"komite"`
}

// Hak membaca hak halaman awal pelaku (workbasket dari sesi; roster komite NONPROP dari basis data).
func (l *Layanan) Hak(ctx context.Context, p inti.Pelaku) (HakPelaku, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HakPelaku{}, err
	}
	h := HakPelaku{WorkbasketTeknik: p.PunyaPeran(models.WorkbasketAcceptation)}
	if len(p.Peran) == 0 || l.a == nil {
		return h, nil
	}
	roster, err := l.a.RosterKomite(ctx, false)
	if err != nil {
		return HakPelaku{}, err
	}
	for _, r := range roster {
		if r.OperatorID != "" && p.PunyaPeran(r.OperatorID) {
			h.Komite = true
			break
		}
	}
	return h, nil
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

// BuatKasus = Start -> Assignment2 (`Flow_TreatyIn`): kasus CLMNP- baru di worklist pembuatnya. Harness New tidak
// diekspor (OQ-CNP-23) - tombol Add Claim halaman awal, pola Claim Prop.
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
	return l.BukaKasus(ctx, p, id, false)
}

// AcuanStatis - daftar pilihan bersama layar.
type AcuanStatis struct {
	MataUang  []models.Pilihan             `json:"mataUang"`
	LossAlloc []models.Pilihan             `json:"lossAlloc"`
	Kode      map[string][]string          `json:"kode"`
	LabelKode map[string]map[string]string `json:"labelKode"`
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
	for _, v := range models.PilihanLossAllocation {
		out.LossAlloc = append(out.LossAlloc, models.Pilihan{Nilai: v, Label: v})
	}
	out.Kode = models.KodePilihan
	out.LabelKode = models.LabelKode
	return out, nil
}
