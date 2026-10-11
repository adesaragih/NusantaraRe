package services

// Untuk apa berkas ini: LAYANAN - galat yang dipetakan handlers, wewenang pemegang assignment, pemuatan kasus beserta
// halaman polis dan medan turunannya, pra-proses flow action, pembentukan layar (layar utama, panel masterDetail
// bersarang, modal), daftar kerja halaman awal, dan pembuatan kasus.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// Galat yang dipetakan handlers ke kode HTTP.
var (
	ErrTanpaOracle        = db.ErrTanpaOracle
	ErrKasusTidakAda      = repository.ErrKasusTidakAda
	ErrTahapBerubah       = repository.ErrTahapBerubah
	ErrPermintaanTidakSah = galat.ErrPermintaanTidakSah
	// ErrKasusTertutup - kasus Resolved-Completed: server menolak setiap tulisan.
	ErrKasusTertutup = errors.New("services: kasus sudah diselesaikan")
	// ErrBukanPemegang - pelaku bukan pemegang assignment kasus (Input Register / Input Estimasi = worklist pembuat,
	// Choose Surveyor = workbasket ReasKlaimTeknik).
	ErrBukanPemegang = fmt.Errorf("%w: bukan pemegang assignment kasus ini", inti.ErrTanpaWewenang)
	// ErrAksiTertutup - aksi tidak tampil atau nonaktif di layar kasus saat ini.
	ErrAksiTertutup = errors.New("services: aksi ini tidak tersedia di layar kasus saat ini")
	// ErrAdjustmentDiKomite - tulisan ke adjustment yang sedang di komite.
	ErrAdjustmentDiKomite = errors.New("services: adjustment ini sedang di komite")
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

// Layanan adalah pintu aturan dagang Claim Fac In.
type Layanan struct {
	g        Gudang
	a        Acuan
	jam      func() time.Time
	produksi bool
	// berkas - penyimpanan berkas lampiran (`DenganPenyimpanan`); nil = unggah / unduh 503.
	berkas PenyimpananBerkas
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

// Pemegang - pelaku memegang assignment kasus (ditegakkan di layanan, prompt §6 butir 10).
//
//	Assignment1 Input Register   `Current operator` -> pembuat kasus (CREATE_OP)
//	Assignment7 Input Estimasi   `ToWorklist` tanpa parameter operator -> operator yang menyerahkan = pembuat
//	Assignment3 Choose Surveyor  `ToWorkbasket` tanpa nama workbasket -> ReasKlaimTeknik (bawaan b prompt §3)
func Pemegang(p inti.Pelaku, k models.Kasus) bool {
	if k.Tertutup() {
		return false
	}
	switch k.Tahap {
	case models.TahapRegister, models.TahapEstimasi:
		return p.AkunID != "" && strings.EqualFold(p.AkunID, k.PembuatID)
	case models.TahapSurveyor:
		return p.PunyaPeran(models.WorkbasketSurveyor)
	}
	return false
}

// konteks menyusun konteks aksi pelaku (tingkat wewenang dari roster, label langkah assignment).
func (l *Layanan) konteks(ctx context.Context, p inti.Pelaku, kasus models.Kasus, saat time.Time) (*models.Konteks, error) {
	tingkat, err := l.a.TingkatPelaku(ctx, p.AkunID)
	if err != nil {
		return nil, err
	}
	return &models.Konteks{Ctx: ctx, Acuan: l.a, Pelaku: p.AkunID, Tingkat: tingkat, Langkah: models.LabelTahap[kasus.Tahap],
		Sekarang: saat, Produksi: l.produksi}, nil
}

// ---------------------------------------------------------------- layar

// Layar - satu kasus siap ditampilkan.
type Layar struct {
	Kasus      models.Kasus    `json:"kasus"`
	Label      string          `json:"label"`
	Halaman    *models.Halaman `json:"halaman"`
	BolehKerja bool            `json:"bolehKerja"`
	// Tata - section assignment kasus (InputRegister / InputEstimasiAdmin / ClaimSurvey).
	Tata []models.Tata `json:"tata"`
	// Panel - panel baris grid masterDetail, kunci `models.KunciPanel` ("est:ClaimData.ObjectList(1)").
	Panel map[string][]models.Tata `json:"panel,omitempty"`
	// Modal - local action / harness: "pilihPolis", "protectDOL", "pla:<o>".
	Modal map[string][]models.Tata `json:"modal,omitempty"`
	Pesan []string                 `json:"pesan,omitempty"`
	// PesanMedan - pesan per jalur (Property-Set-Messages).
	PesanMedan map[string][]string `json:"pesanMedan,omitempty"`
	// Info - pemberitahuan sesudah aksi (mis. berkas dokumen menunggu OQ).
	Info string `json:"info,omitempty"`
	// BukaModal - modal yang dibuka layar sesudah aksi.
	BukaModal string `json:"bukaModal,omitempty"`
	// Mode - penanda mode layar (models.ModeLayar); layar mengembalikannya di setiap aksi.
	Mode map[string]string `json:"mode,omitempty"`
}

// Kunci modal.
const (
	ModalPilihPolis = "pilihPolis"
	ModalProtectDOL = "protectDOL"
	AwalanModalPLA  = "pla"
)

// tataKasus mengevaluasi seluruh section yang terbuka bagi kasus ini.
func tataKasus(k models.Kasus, h *models.Halaman, kunci bool) ([]models.Tata, map[string][]models.Tata, map[string][]models.Tata) {
	panel := map[string][]models.Tata{}
	modal := map[string][]models.Tata{}
	var utama []models.Tata
	estimasi := func() {
		for o := range h.AmbilDaftar(models.DaftarObjek) {
			on := o + 1
			panel[models.KunciPanel(models.PanelObjekEst, models.DaftarObjek, on)] =
				models.Evaluasi(h, models.LayarObjekEstimasi(on), kunci)
			modal[fmt.Sprintf("%s:%d", AwalanModalPLA, on)] = models.Evaluasi(h, models.LayarPLA(on), kunci)
			for i := range h.AmbilDaftar(models.DaftarItem(on)) {
				panel[models.KunciPanel(models.PanelItemEst, models.DaftarItem(on), i+1)] =
					models.Evaluasi(h, models.LayarItemEstimasi(on, i+1), kunci)
			}
		}
	}
	lolos := func(h *models.Halaman) bool {
		return h.Ambil(models.JalurProtect1) == "1" && h.Ambil(models.JalurProtect2) == "1"
	}
	adjustment := func() {
		for o := range h.AmbilDaftar(models.DaftarObjek) {
			on := o + 1
			panel[models.KunciPanel(models.PanelObjekAdj, models.DaftarObjek, on)] =
				models.Evaluasi(h, models.LayarObjekAdj(on), kunci)
			modal[fmt.Sprintf("%s:%d", AwalanModalDLA, on)] = models.Evaluasi(h, models.LayarDLA(on), kunci)
			for i := range h.AmbilDaftar(models.DaftarItem(on)) {
				in := i + 1
				panel[models.KunciPanel(models.PanelItemAdj, models.DaftarItem(on), in)] =
					models.Evaluasi(h, models.LayarItemAdj(on, in), kunci)
				for a := range h.AmbilDaftar(models.DaftarAdj(on, in)) {
					panel[models.KunciPanel(models.PanelAdj, models.DaftarAdj(on, in), a+1)] =
						models.Evaluasi(h, models.LayarAdjustment(on, in, a+1), kunci)
					modal[models.KunciPanel(models.ModalCedant, models.DaftarAdj(on, in), a+1)] =
						models.Evaluasi(h, models.LayarCedant(on, in, a+1), kunci)
					modal[kunciModalKomite(on, in, a+1)] = models.Evaluasi(h, models.LayarKomite(on, in, a+1, lolos), kunci)
				}
			}
		}
	}
	switch k.Tahap {
	case models.TahapRegister:
		utama = models.Evaluasi(h, models.LayarRegister(), kunci)
		modal[ModalPilihPolis] = models.Evaluasi(h, models.LayarPilihPolis(), kunci)
		modal[ModalProtectDOL] = models.Evaluasi(h, models.LayarProtectDOL(), kunci)
		modal[ModalTolak] = models.Evaluasi(h, models.LayarTolak(), kunci)
	case models.TahapEstimasi:
		utama = models.Evaluasi(h, models.LayarEstimasi(), kunci)
		modal[ModalPilihPolis] = models.Evaluasi(h, models.LayarPilihPolis(), kunci)
		estimasi()
	case models.TahapSurveyor:
		utama = models.Evaluasi(h, models.LayarSurveyor(), kunci)
		modal[ModalTutup] = models.Evaluasi(h, models.LayarTutup(), kunci)
		estimasi()
		adjustment()
	default: // Resolved: layar terakhir, terkunci
		kunci = true
		utama = models.Evaluasi(h, models.LayarSurveyor(), true)
		estimasi()
		adjustment()
	}
	for n := range h.AmbilDaftar(models.DaftarProgres) {
		panel[models.KunciPanel(models.PanelProgres, models.DaftarProgres, n+1)] =
			models.Evaluasi(h, models.LayarSubProgres(n+1), kunci)
	}
	return utama, panel, modal
}

// semuaTata - gabungan seluruh tata (untuk medan terbuka).
func semuaTata(utama []models.Tata, panel, modal map[string][]models.Tata) []models.Tata {
	out := append([]models.Tata{}, utama...)
	for _, t := range panel {
		out = append(out, t...)
	}
	for _, t := range modal {
		out = append(out, t...)
	}
	return out
}

func (l *Layanan) layar(k models.Kasus, h *models.Halaman, boleh bool) *Layar {
	utama, panel, modal := tataKasus(k, h, !boleh)
	return &Layar{Kasus: k, Label: models.LabelTahap[k.Tahap], Halaman: h, BolehKerja: boleh, Tata: utama,
		Panel: panel, Modal: modal, Pesan: h.SemuaPesan(), PesanMedan: h.Pesan, Mode: models.AmbilMode(h)}
}

// ---------------------------------------------------------------- pemuatan

// muat membaca kasus, halamannya, dan halaman polis (tx boleh nil).
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
	if err := l.muatPolis(ctx, h); err != nil {
		return models.Kasus{}, nil, err
	}
	return k, h, nil
}

// muatPolis membaca ulang halaman polis `OfferFacIn` dari JSON_POLIS (POLICY_NO + PRODKE kasus), memasang pengganti
// FacRetroList yang tersimpan, IsTreatyIn, dan medan salinan polis setiap objek.
func (l *Layanan) muatPolis(ctx context.Context, h *models.Halaman) error {
	nopolis, prodke := h.Ambil(models.JalurNoPolis), h.Ambil(models.JalurProdke)
	if nopolis != "" {
		dok, ada, err := l.a.DokumenPolis(ctx, nopolis, prodke)
		if err != nil {
			return err
		}
		if ada {
			if err := models.TerapkanPolisTersimpan(h, dok); err != nil {
				return err
			}
			if h.Ambil(models.JalurNoPolis) == "" { // dokumen tanpa PolicyNo: nomor kasus dipertahankan
				h.Setel(models.JalurNoPolis, nopolis)
			}
		}
		h.Setel(models.JalurIsTreatyIn, models.IsTreatyInDari(nopolis))
	}
	models.TerapkanRetroTreaty(h)
	models.LengkapiObjek(h)
	models.LengkapiCedant(h)
	return nil
}

// turunkan menyusun medan TURUNAN halaman (tidak disimpan): indeks dan nomor baris, grid Claim Status (kronologi urut
// menurun), Progress Claim, kasus kembar ProtectDOL.
func (l *Layanan) turunkan(ctx context.Context, k models.Kasus, h *models.Halaman) error {
	models.HitungTurunan(h)
	h.SetelDaftar(models.DaftarKronologiTampil, models.TampilKronologi(h.AmbilDaftar(models.DaftarKronologi)))
	pr, sub, err := l.a.ProgresKlaim(ctx, k.ID)
	if err != nil {
		return err
	}
	h.SetelDaftar(models.DaftarProgres, pr)
	for n, rows := range sub {
		h.SetelDaftar(models.JalurAnak(models.DaftarProgres, n, models.AnakSubProg), rows)
	}
	if err := l.turunkanAdjustment(ctx, h); err != nil {
		return err
	}
	if k.Tahap == models.TahapRegister && h.Ambil(models.JalurIsError) == "1" {
		salin := h.Salin()
		kembar, err := models.CheckDoubleClaim(ctx, l.a, salin)
		if err != nil {
			return err
		}
		var rows []models.Baris
		for _, r := range kembar {
			rows = append(rows, models.Baris{"BRANCH_NAME": r.ClaimNo})
		}
		h.SetelDaftar(models.DaftarKembar, rows)
	}
	return nil
}

// turunkanAdjustment - medan turunan baris adjustment: daftar komite calon (`SetListKomite_act`, roster FACIN) bagi
// adjustment yang sudah dihitung dan belum berkomite; status kasir (`GetStatusKasir_Act`, defer load blok kasir) bagi
// adjustment ber-AcceptedNo.
func (l *Layanan) turunkanAdjustment(ctx context.Context, h *models.Halaman) error {
	var roster []models.AnggotaKomite
	dibaca := false
	for o := range h.AmbilDaftar(models.DaftarObjek) {
		for i := range h.AmbilDaftar(models.DaftarItem(o + 1)) {
			for a, b := range h.AmbilDaftar(models.DaftarAdj(o+1, i+1)) {
				if kmt := b[models.PropKomiteID]; kmt != "" { // tangga kasus komite tersimpan
					tangga, err := l.a.TanggaKomite(ctx, kmt)
					if err != nil {
						return err
					}
					if tangga == nil {
						tangga = []models.AnggotaKomite{}
					}
					if err := models.SusunKomiteAdjustment(h, o+1, i+1, a+1, nil, tangga); err != nil {
						return err
					}
				} else if b["PaymentType"] != "" {
					if !dibaca {
						var err error
						if roster, err = l.a.RosterKomite(ctx); err != nil {
							return err
						}
						dibaca = true
					}
					if err := models.SusunKomiteAdjustment(h, o+1, i+1, a+1, roster, nil); err != nil {
						return err
					}
				}
				if b["AcceptedNo"] != "" {
					ket, ada, err := l.a.StatusKasir(ctx, b["AcceptedNo"])
					if err != nil {
						return err
					}
					models.TerapkanStatusKasir(b, ket, ada)
				}
			}
		}
	}
	return nil
}

// siapkan = pra-proses flow action tahap kasus bagi pemegangnya:
//
//	InputRegister  CallActivityInputRegister + InsertObjects_dt (riwayat klaim polis, calon objek)
//	InputEstimasi  InputEstimationPre + SetEstimation_DT (SetMOClaim_Act)
//	InputSurveyor  SetTypePDFAdjustment + SetStartDateAdjustment (SetMOClaim_Act)
//
// InsertProgressClaim (bagian ketiga pra-proses) menulis PROGRESSCLAIM - dijalankan di transaksi aksi (`tulisProgres`).
func (l *Layanan) siapkan(ctx context.Context, k *models.Konteks, kasus models.Kasus, h *models.Halaman) error {
	switch kasus.Tahap {
	case models.TahapRegister:
		var riwayat []models.RiwayatKlaimPolis
		if nopol := h.Ambil(models.JalurNoPolis); nopol != "" {
			var err error
			if riwayat, err = l.a.RiwayatKlaimPolis(ctx, nopol); err != nil {
				return err
			}
		}
		models.PraRegister(k, h, riwayat)
	case models.TahapEstimasi:
		models.PraEstimasi(h)
		if err := models.LengkapiMarine(k, h); err != nil {
			return err
		}
		return models.SetMOClaim(k, h)
	case models.TahapSurveyor:
		models.PraAdjustment(k, h)
		return models.SetMOClaim(k, h)
	}
	return nil
}

// tulisProgres = InsertProgressClaim (pra-proses setiap assignment; idempoten per posisi).
func (l *Layanan) tulisProgres(ctx context.Context, tx *db.Tx, k *models.Konteks, id string, h *models.Halaman,
	p models.ParamProgres) error {
	pr, sub, ada := models.RencanaProgres(k, h, models.KunciInstans(id), p)
	if !ada {
		return nil
	}
	return l.g.TulisProgres(ctx, tx, pr, sub)
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
	kt, err := l.konteks(ctx, p, k, l.jam())
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
	if err := l.turunkan(ctx, k, h); err != nil {
		return nil, err
	}
	return l.layar(k, h, boleh), nil
}

// HakPelaku - hak halaman awal: tab workbasket aktif hanya bagi pemegang workbasket Choose Surveyor; tabel komite di
// bawah inbox (modul Komite Claim Fac In tanpa menu, prompt tahap 2 §2 butir 2) hanya bagi anggota roster komite FACIN.
type HakPelaku struct {
	WorkbasketSurveyor bool `json:"workbasketSurveyor"`
	Komite             bool `json:"komite"`
}

// WorkbasketSPVB - cadangan SPV A di tingkat 1 komite (keputusan work owner 10-10-2026 KCF-01; bukan baris roster).
const WorkbasketSPVB = "ReasClaimSPVB"

// Hak membaca hak halaman awal pelaku (workbasket dari sesi; roster komite FACIN dari basis data).
func (l *Layanan) Hak(ctx context.Context, p inti.Pelaku) (HakPelaku, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HakPelaku{}, err
	}
	h := HakPelaku{WorkbasketSurveyor: p.PunyaPeran(models.WorkbasketSurveyor), Komite: p.PunyaPeran(WorkbasketSPVB)}
	if h.Komite || l.a == nil {
		return h, nil
	}
	roster, err := l.a.RosterKomite(ctx)
	if err != nil {
		return HakPelaku{}, err
	}
	for _, r := range roster {
		if r.OperatorID != "" && (r.OperatorID == p.AkunID || p.PunyaPeran(r.OperatorID)) {
			h.Komite = true
			break
		}
	}
	return h, nil
}

// Jenis daftar halaman awal.
const (
	DaftarSaya       = "saya"       // Input Register + Input Estimasi - worklist pembuat
	DaftarWorkbasket = "workbasket" // Choose Surveyor - workbasket ReasKlaimTeknik
	DaftarSelesai    = "selesai"    // Resolved-Completed
)

// DaftarKasus - halaman awal: pintu masuk Pega yang berbukti (Register_Flow Assignment1 / Assignment7 ke worklist,
// Assignment3 ke workbasket).
func (l *Layanan) DaftarKasus(ctx context.Context, p inti.Pelaku, jenis, cari string) ([]repository.RingkasanKasus, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	s := repository.SaringanKasus{Cari: cari}
	switch jenis {
	case DaftarSaya, "":
		s.TahapIn, s.Pembuat = []string{models.TahapRegister, models.TahapEstimasi}, p.AkunID
	case DaftarWorkbasket:
		if !p.PunyaPeran(models.WorkbasketSurveyor) {
			return []repository.RingkasanKasus{}, nil
		}
		s.Tahap = models.TahapSurveyor
	case DaftarSelesai:
		s.Selesai = true
	default:
		return nil, fmt.Errorf("%w: daftar %q", ErrPermintaanTidakSah, jenis)
	}
	return l.g.DaftarKasus(ctx, s)
}

// BuatKasus = Start -> Decision4 (IsSPK salah: polis belum dipilih) -> Assignment1 Input Register di worklist
// pembuatnya. Harness New tidak diekspor - tombol Add Claim halaman awal, pola Claim Prop.
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
		if err := l.g.SisipKasus(ctx, tx, id, p.AkunID, nama, saat); err != nil {
			return err
		}
		k, h, err := l.muat(ctx, tx, id)
		if err != nil {
			return err
		}
		kt, err := l.konteks(ctx, p, k, saat)
		if err != nil {
			return err
		}
		if err := l.siapkan(ctx, kt, k, h); err != nil {
			return err
		}
		h.BersihkanPesan()
		if err := l.g.SimpanHalaman(ctx, tx, id, h); err != nil {
			return err
		}
		return l.tulisProgres(ctx, tx, kt, id, h, models.ParamProgres{})
	})
	if err != nil {
		return nil, err
	}
	return l.BukaKasus(ctx, p, id, false)
}

// AcuanStatis - daftar pilihan bersama layar.
type AcuanStatis struct {
	MataUang  []models.Pilihan             `json:"mataUang"`
	JenisReas []models.Pilihan             `json:"jenisReas"`
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
	if out.JenisReas, err = l.a.DaftarJenisReas(ctx); err != nil {
		return out, err
	}
	out.Kode = models.KodePilihan
	out.LabelKode = models.LabelKode
	return out, nil
}
