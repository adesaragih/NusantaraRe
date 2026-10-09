package services

// Untuk apa berkas ini: AKSI LAYAR. Setiap aksi (refresh ber-activity sebuah sel, tombol, tombol Choose pop-up) berjalan
// dalam SATU transaksi:
//
//	kunci baris kasus (tahap belum berubah, belum tertutup) -> pelaku pemegang assignment
//	-> halaman TERSIMPAN dimuat + pra-proses + turunan -> tata dievaluasi
//	-> aksi harus tampil dan aktif di tata itu (AksiTerbuka) -> kiriman layar digabung HANYA untuk jalur terbuka
//	-> port activity dijalankan -> turunan dihitung ulang -> halaman disimpan
//
// ⚠️ PENYIMPANGAN SADAR (PARITAS, pola Claim Prop): Pega menyimpan clipboard hanya pada Save / Obj-Save; di sini setiap
// aksi menyimpan halaman (write-through). Aksi yang keluar karena pesan validasi (`GalatValidasi`) dibatalkan
// seluruhnya; layar tetap menampilkan isian dan pesannya.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimnonprop/backend/models"
)

// PermintaanAksi - satu aksi layar.
type PermintaanAksi struct {
	// Aksi - nama aksi di tata (nama activity tanpa akhiran, atau ID tombol); ":param" = parameter activity.
	Aksi string `json:"aksi"`
	// Indeks - baris grid tingkat klaim (Param.Index / idx), atau nomor akseptasi untuk aksi panel detail akseptasi.
	Indeks int `json:"indeks"`
	// Baris - baris grid DI DALAM panel detail akseptasi (Param.idx), 0 bila bukan sel grid panel.
	Baris int `json:"baris"`
	// Param - parameter pilihan pop-up (ID baris terpilih) atau parameter activity.
	Param string `json:"param"`
	// Tahap - tahap kasus yang dilihat layar (penjaga balapan: tahap berubah = 409).
	Tahap string `json:"tahap"`
	// Masukan - nilai medan terbuka layar (jalur -> teks).
	Masukan map[string]string `json:"masukan"`
	// Mode - penanda mode layar terakhir (`Layar.Mode`, models.ModeLayar) yang dikembalikan layar.
	Mode map[string]string `json:"mode,omitempty"`
}

// jalanAksi - konteks satu aksi untuk penanganannya.
type jalanAksi struct {
	l     *Layanan
	ctx   context.Context
	tx    *db.Tx
	k     *models.Konteks
	kasus models.Kasus
	p     inti.Pelaku
	h     *models.Halaman
	r     PermintaanAksi
	m     models.MasterTreaty
	// info / bukaModal - pemberitahuan dan modal yang dibuka layar sesudah aksi.
	info, bukaModal string
	// selesai - kasus sudah ditutup / dipindah oleh aksi ini (halaman disimpan sebelum pemindahan).
	selesai bool
}

type penanganAksi func(j *jalanAksi) error

// aksiPanel - aksi panel detail akseptasi `Indeks` (AdjustmentDetailNP_Section dan modal komite / reinstatement).
var aksiPanel = map[string]bool{"SetInterimXOL": true, "HapusKlaimAkseptasi": true, "HapusLossAkseptasi": true,
	"AdjClaimCNP": true, "SetPayableTreatyNP": true, "PilihRekening": true, "PilihRekening2": true, "BukaKomite": true,
	"CreateChildKomiteCNP": true, "HitServiceToKasir": true, "GenerateCACNP": true, "SaveCNPLayerList": true}

// aksiPanelBaca - aksi panel yang tidak menulis akseptasi (boleh saat akseptasi di komite).
var aksiPanelBaca = map[string]bool{"GenerateCACNP": true}

// aksiInduk - aksi pop-up (tombol Choose / Save di harness) yang terbuka bila salah satu tombol PEMBUKA pop-up-nya
// terbuka di tata.
var aksiInduk = map[string][]string{
	"SetValueClaimTNP":   {"ChooseMasterIn", "ChooseMasterInEDM"},
	"CheckNoPolicy":      {"ViewListPolicy", "CheckNoPolicy"},
	"GetNameCauseofLoss": {"ChooseCauseOfLoss"},
	"SetCatastrope":      {"CatastrofeList"},
	"SaveCatasrtope":     {"CatastrofeList"},
	"InputCatastrope":    {"CatastrofeList"},
}

// Aksi menjalankan satu aksi dan mengembalikan layar sesudahnya.
func (l *Layanan) Aksi(ctx context.Context, p inti.Pelaku, id string, r PermintaanAksi) (*Layar, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	nama, param := r.Aksi, ""
	if i := strings.Index(nama, ":"); i >= 0 { // "SetEditCatastrope:Edit"
		nama, param = nama[:i], nama[i+1:]
	}
	f, ada := penangan[nama]
	if !ada {
		return nil, fmt.Errorf("%w: aksi %q", ErrPermintaanTidakSah, r.Aksi)
	}
	var hasil *Layar
	var gagal *GalatValidasi
	err := l.g.Transaksi(ctx, func(tx *db.Tx) error {
		kasus, err := l.g.Keadaan(ctx, tx, id)
		if err != nil {
			return err
		}
		if kasus.Tertutup() {
			return ErrKasusTertutup
		}
		if r.Tahap != "" && r.Tahap != kasus.Tahap {
			return ErrTahapBerubah
		}
		if !Pemegang(p, kasus) {
			return ErrBukanPemegang
		}
		if kasus, err = l.g.KunciKasus(ctx, tx, id, kasus.Tahap); err != nil {
			return err
		}
		h, err := l.g.BacaHalaman(ctx, tx, id)
		if err != nil {
			return err
		}
		h.Setel("pyID", kasus.ID)
		models.PasangMode(h, r.Mode)
		kt, err := l.konteks(ctx, p, l.jam())
		if err != nil {
			return err
		}
		if err := l.siapkan(ctx, kt, kasus, h); err != nil {
			return err
		}
		if err := l.turunkan(ctx, h); err != nil {
			return err
		}
		h.BersihkanPesan() // pesan pra-proses tidak dibawa ke hasil aksi
		utama, adj, modal := tataKasus(kasus, h, false)
		semua := semuaTata(utama, adj, modal)
		// tataMasuk - medan yang kiriman layarnya diterima: aksi panel menerima medan layar utama (postValue tanpa
		// activity, mis. Occupation) dan panel akseptasinya sendiri, bukan panel akseptasi lain.
		tataMasuk := semua
		idxCek := r.Indeks
		if aksiPanel[nama] {
			b, err := adjBaris(h, r.Indeks)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrPermintaanTidakSah, err)
			}
			if !aksiPanelBaca[nama] && b["IsKomite"] == "1" && b["AcceptanceStatus"] == "0" {
				return ErrAkseptasiDiKomite
			}
			semua = append(append([]models.Tata{}, adj[r.Indeks]...), modal[fmt.Sprintf("komite:%d", r.Indeks)]...)
			tataMasuk = append(append([]models.Tata{}, utama...), semua...)
		}
		// Sel grid di dalam panel akseptasi (Baris > 0): gerbang sel dibaca pada baris grid itu - juga untuk aksi tingkat
		// klaim dari panel (Rate of Exchange -> CountClaimTNP).
		if r.Baris > 0 {
			idxCek = r.Baris
		}
		if !aksiTerbuka(semua, nama, r.Aksi, idxCek) {
			return fmt.Errorf("%w: %q", ErrAksiTertutup, r.Aksi)
		}
		models.GabungMasukan(h, r.Masukan, models.MedanTerbuka(tataMasuk))
		if param != "" && r.Param == "" {
			r.Param = param
		}
		m, err := l.master(ctx, h)
		if err != nil {
			return err
		}
		j := &jalanAksi{l: l, ctx: ctx, tx: tx, k: kt, kasus: kasus, p: p, h: h, r: r, m: m}
		if err := f(j); err != nil {
			var gv *GalatValidasi
			if errors.As(err, &gv) {
				gagal = gv
				gv.Layar = l.layar(kasus, h, true)
			}
			return err
		}
		if err := models.HitungTurunan(h); err != nil {
			return err
		}
		if !j.selesai {
			if err := l.g.SimpanHalaman(ctx, tx, id, h); err != nil {
				return err
			}
			if err := l.g.SentuhKasus(ctx, tx, id, kt.Sekarang); err != nil {
				return err
			}
		}
		pesan := h.Pesan
		k2, h2, err := l.muat(ctx, tx, id)
		if err != nil {
			return err
		}
		if Pemegang(p, k2) {
			if err := l.siapkan(ctx, kt, k2, h2); err != nil {
				return err
			}
		}
		if err := l.turunkan(ctx, h2); err != nil {
			return err
		}
		h2.BersihkanPesan()
		h2.Pesan = pesan
		for _, x := range append([]string{"ParamInput.CARI1"}, models.ModeLayar...) {
			if v := h.Ambil(x); v != "" {
				h2.Setel(x, v)
			}
		}
		hasil = l.layar(k2, h2, Pemegang(p, k2))
		hasil.Info, hasil.BukaModal = j.info, j.bukaModal
		return nil
	})
	if gagal != nil {
		return nil, gagal
	}
	if err != nil {
		return nil, err
	}
	return hasil, nil
}

// aksiTerbuka - aksi (nama, atau nama berparameter apa adanya) atau salah satu pembuka pop-up-nya terbuka di tata.
func aksiTerbuka(ts []models.Tata, nama, utuh string, indeks int) bool {
	if models.AksiTerbuka(ts, utuh, indeks) || models.AksiTerbuka(ts, nama, indeks) {
		return true
	}
	for _, induk := range aksiInduk[nama] {
		if models.AksiTerbuka(ts, induk, indeks) {
			return true
		}
	}
	return false
}

// validasi mengubah pesan halaman menjadi GalatValidasi (aksi dibatalkan).
func validasi(h *models.Halaman) error {
	if !h.AdaPesan() {
		return nil
	}
	return &GalatValidasi{Pesan: h.SemuaPesan()}
}

func adjBaris(h *models.Halaman, n int) (models.Baris, error) {
	d := h.AmbilDaftar(models.DaftarAdjustment)
	if n < 1 || n > len(d) {
		return nil, models.ErrBarisTidakAda
	}
	return d[n-1], nil
}

// ---------------------------------------------------------------- penangan sederhana

func halamanSaja(f func(j *jalanAksi)) penanganAksi {
	return func(j *jalanAksi) error { f(j); return nil }
}

func dariModel(f func(k *models.Konteks, h *models.Halaman) error) penanganAksi {
	return func(j *jalanAksi) error { return f(j.k, j.h) }
}

// mataUangBaris - `.Currency` baris grid tingkat klaim sesudah kiriman digabung (`Param.Currency = .Currency`).
func mataUangBaris(j *jalanAksi, daftar string) string {
	d := j.h.AmbilDaftar(daftar)
	if j.r.Indeks < 1 || j.r.Indeks > len(d) {
		return ""
	}
	return d[j.r.Indeks-1]["Currency"]
}

// hitungKlaim = CountClaimTNP_Act (tingkat klaim, jalur absolut pyWorkPage).
func hitungKlaim(j *jalanAksi) error { return models.HitungKlaim(j.k, j.h, j.m) }

// penangan - aksi -> port activity. Nama = nama activity tanpa akhiran `_Act`/`_act`/`_ACT`.
var penangan map[string]penanganAksi

func init() {
	penangan = map[string]penanganAksi{
		// registrasi
		"SetValueClaimTNP":    aksiPilihMaster,
		"CheckNoPolicy":       aksiPilihPolis,
		"SetEndDate":          halamanSaja(func(j *jalanAksi) { models.SetEndDate(j.h) }),
		"CheckPeriodPolicy":   halamanSaja(func(j *jalanAksi) { models.CheckPeriodPolicy(j.h) }),
		"CheckDateDOL":        dariModel(models.CheckDateDOL),
		"CheckReportDate":     halamanSaja(func(j *jalanAksi) { models.CheckReportDate(j.k, j.h) }),
		"CheckDateReceived":   halamanSaja(func(j *jalanAksi) { models.CheckDateReceived(j.k, j.h) }),
		"GetReportStatus":     dariModel(models.GetReportStatus),
		"SetDefNonCatastrope": halamanSaja(func(j *jalanAksi) { models.SetDefNonCatastrope(j.h) }),
		"SetEditCatastrope":   halamanSaja(func(j *jalanAksi) { models.SetEditCatastrope(j.h, j.r.Param) }),
		"BukaKatastrofe":      halamanSaja(func(*jalanAksi) {}),
		"InputCatastrope":     halamanSaja(func(j *jalanAksi) { models.InputCatastrope(j.h, j.r.Param) }),
		"SetCatastrope":       aksiPilihKatastrofe,
		"SaveCatasrtope":      aksiSimpanKatastrofe,
		"GetNameCauseofLoss":  aksiPilihSebab,
		"SetConsultant":       dariModel(models.SetConsultant),
		"SetAdjsuter":         dariModel(models.SetAdjuster),
		"MakeLowercase":       halamanSaja(func(j *jalanAksi) { models.MakeLowercase(j.h) }),
		"GetAdders":           dariModel(models.GetAdders),
		// interest, deductible, klaim 100%, Loss Allocation, XoL
		"AddInterestListCNP": halamanSaja(func(j *jalanAksi) { models.AddInterest(j.h) }),
		"HapusInterest": func(j *jalanAksi) error {
			return models.HapusBarisDaftar(j.h, models.DaftarInterest, j.r.Indeks)
		},
		"SetCurrency": func(j *jalanAksi) error {
			return models.SetCurrency(j.k, j.h, j.m, j.r.Indeks, j.r.Param)
		},
		"CountTotalInterest": halamanSaja(func(*jalanAksi) {}), // turunan (models.HitungTurunan)
		"SetTPLNote":         func(j *jalanAksi) error { return models.SetTPLNote(j.k, j.h, j.m, j.r.Indeks) },
		"SetFormat":          halamanSaja(func(j *jalanAksi) { models.SetFormat(j.k, j.h) }),
		"SetDataDeductible": func(j *jalanAksi) error {
			return models.SetDataDeductible(j.h, j.r.Param, j.h.Ambil(models.CD+"TypeDeductible"))
		},
		"AddListClaimNP": halamanSaja(func(j *jalanAksi) { models.AddListClaim(j.k, j.h, j.r.Param) }),
		"CountClaimTNP":  hitungKlaim,
		"HapusKlaim": func(j *jalanAksi) error {
			if err := models.HapusBarisDaftar(j.h, models.DaftarClaimAmount, j.r.Indeks); err != nil {
				return err
			}
			return hitungKlaim(j)
		},
		"AddLossAlocation": halamanSaja(func(j *jalanAksi) { models.AddLossAllocation(j.k, j.h) }),
		"CountLossAllocation": func(j *jalanAksi) error {
			return models.HitungXOL(j.k, j.h, j.m, j.r.Param, j.r.Indeks, mataUangBaris(j, models.DaftarLossAlloc))
		},
		"HapusLossAlloc": func(j *jalanAksi) error {
			return models.HapusBarisDaftar(j.h, models.DaftarLossAlloc, j.r.Indeks)
		},
		"AdjClaimAmount":   func(j *jalanAksi) error { return models.AdjClaimAmount(j.h, j.r.Indeks) },
		"IntIsSaveToOs":    halamanSaja(func(j *jalanAksi) { j.h.Setel("IsSaveToOs", "0") }),
		"SetActualPremium": func(j *jalanAksi) error { return models.SetActualPremium(j.h) },
		// tombol Outstanding Claim
		"SaveDataToJClaim":  aksiSimpan,
		"SaveDataToOSAksep": aksiSimpanOS,
		"GenerateCFS":       aksiCFS,
		"GeneratePlaCNP":    aksiPLA,
		"Submit":            aksiSubmit,
		"Simpan":            halamanSaja(func(*jalanAksi) {}),
		// Input Acceptation
		"SaveToOS":           aksiSaveToOS,
		"AddAkseptasiCNP":    aksiTambahAkseptasi,
		"DeleteAkseptasi":    func(j *jalanAksi) error { return models.DeleteAkseptasi(j.h, j.r.Indeks) },
		"CloseClaimTNonProp": aksiTutupKlaim,
		"BukaCWP":            aksiBukaCWP,
		// detail akseptasi (Indeks = akseptasi, Baris = baris grid panel)
		"SetInterimXOL": func(j *jalanAksi) error { return models.SetInterimXOL(j.h, j.r.Indeks) },
		"HapusKlaimAkseptasi": func(j *jalanAksi) error {
			return models.HapusBarisDaftar(j.h, models.JalurAdj(j.r.Indeks, models.AnakClaimAccept), j.r.Baris)
		},
		"HapusLossAkseptasi": func(j *jalanAksi) error {
			return models.HapusBarisDaftar(j.h, models.JalurAdj(j.r.Indeks, models.AnakLossAlloc), j.r.Baris)
		},
		"AdjClaimCNP": func(j *jalanAksi) error { return models.AdjClaimCNP(j.h, j.r.Indeks, j.r.Baris) },
		"SetPayableTreatyNP": func(j *jalanAksi) error {
			if err := models.SetPayableNP(j.k, j.h, j.r.Indeks, j.r.Param); err != nil {
				return err
			}
			return models.SetAccountNo(j.h, j.r.Indeks)
		},
		"PilihRekening":        func(j *jalanAksi) error { return aksiPilihRekening(j, "") },
		"PilihRekening2":       func(j *jalanAksi) error { return aksiPilihRekening(j, "2") },
		"BukaKomite":           aksiBukaKomite,
		"CreateChildKomiteCNP": aksiKomite,
		"HitServiceToKasir":    aksiKasir,
		"GenerateCACNP": func(j *jalanAksi) error {
			models.CetakCA(j.h)
			j.info = models.PesanBerkasOQ
			return nil
		},
		"SaveCNPLayerList": aksiSimpanDibayar,
	}
}
