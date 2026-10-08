package services

// Untuk apa berkas ini: AKSI LAYAR. Setiap aksi (refresh ber-activity sebuah sel, tombol, tombol Choose pop-up) berjalan
// dalam SATU transaksi:
//
//	kunci baris kasus (tahap belum berubah, belum tertutup) -> pelaku pemegang assignment
//	-> halaman TERSIMPAN dimuat + pra-proses + turunan -> tata dievaluasi
//	-> aksi harus tampil dan aktif di tata itu (AksiTerbuka) -> kiriman layar digabung HANYA untuk jalur terbuka
//	-> port activity dijalankan -> turunan dihitung ulang -> halaman disimpan
//
// ⚠️ PENYIMPANGAN SADAR (PARITAS): Pega menyimpan clipboard hanya pada Save / Obj-Save; di sini setiap aksi menyimpan
// halaman (write-through) - baris yang ditambah tombol Add tidak hilang di antara permintaan. Aksi yang keluar karena
// pesan validasi (`GalatValidasi`) dibatalkan seluruhnya; layar tetap menampilkan isian dan pesannya.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimprop/backend/models"
)

// PermintaanAksi - satu aksi layar.
type PermintaanAksi struct {
	// Aksi - nama aksi di tata (nama activity tanpa akhiran, atau ID tombol).
	Aksi string `json:"aksi"`
	// Indeks - Param.Index (baris grid / baris adjustment), berbasis 1.
	Indeks int `json:"indeks"`
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
	lama  *models.Halaman // halaman tersimpan sebelum kiriman digabung (TWorkPage CopyOldataCurr_act)
	h     *models.Halaman
	r     PermintaanAksi
	// info - pemberitahuan local action (PrintFile / PrintFileDLA).
	info string
	// tutup - kasus sudah ditutup / dipindah oleh aksi ini (halaman disimpan sebelum pemindahan).
	selesai bool
}

type penanganAksi func(j *jalanAksi) error

// aksiBarisAdj - aksi yang bekerja atas satu baris AdjustmentList (`Indeks`).
var aksiBarisAdj = map[string]bool{"SetPayableTreaty": true, "SetNameCurrency": true, "SetDLACedingSOB": true,
	"PilihAllocation": true, "CountGrossAdjTreaty": true, "PilihRekening": true, "SetKomiteNo": true, "BukaKomite": true,
	"AddKomiteTreatyChild": true, "HitServiceToKasir": true, "PrintDLATreatyIn": true}

// aksiInduk - aksi pop-up (tombol Choose / Save di harness) yang terbuka bila aksi PEMBUKA pop-up-nya terbuka di tata.
var aksiInduk = map[string]string{
	"SetValueToClaim":    "PilihMaster",
	"CheckNoPolicy":      "PilihPolis",
	"GetNameCauseofLoss": "PilihSebab",
	"SetCatastrope":      "BukaKatastrofe",
	"SaveCatasrtope":     "BukaKatastrofe",
	"BukaKatastrofe":     "BukaKatastrofe",
	"TryMakePLA":         "BukaPLA",
	"PrintDLATreatyIn":   "BukaDLA",
	"CloseClaimProp":     "BukaTutupKlaim",
	// tombol modal CommitteeTreaty tampil menurut isian modal - diperiksa ulang SESUDAH kiriman digabung (aksiKomite)
	"AddKomiteTreatyChild": "BukaKomite",
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
	if param != "" && r.Param == "" {
		r.Param = param
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
		if err := l.siapkan(kt, kasus, h); err != nil {
			return err
		}
		if err := l.turunkan(ctx, kt, h); err != nil {
			return err
		}
		// Pesan pra-proses tidak dibawa ke hasil aksi (pesan aksi yang tampil).
		h.BersihkanPesan()
		if nama == "AddKomiteTreatyChild" { // gerbang Protect dihitung ulang di server, bukan dari kiriman (AC 58)
			if err := j0BukaKomite(l, ctx, kt, h, r.Indeks); err != nil {
				return err
			}
		}
		utama, adj, modal := tataKasus(kasus, h, false)
		semua := semuaTata(utama, adj, modal)
		if aksiBarisAdj[nama] { // aksi baris adjustment: hanya tata baris itu (dan modalnya)
			semua = append(append(append([]models.Tata{}, adj[r.Indeks]...), modal[fmt.Sprintf("dla:%d", r.Indeks)]...),
				modal[fmt.Sprintf("komite:%d", r.Indeks)]...)
		}
		cek := r.Aksi
		if induk, ada := aksiInduk[nama]; ada {
			cek = induk
		}
		if !models.AksiTerbuka(semua, cek, r.Indeks) && !models.AksiTerbuka(semua, nama, r.Indeks) {
			return fmt.Errorf("%w: %q", ErrAksiTertutup, r.Aksi)
		}
		lama := h.Salin()
		models.GabungMasukan(h, r.Masukan, models.MedanTerbuka(semua))
		j := &jalanAksi{l: l, ctx: ctx, tx: tx, k: kt, kasus: kasus, p: p, lama: lama, h: h, r: r}
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
			if err := l.siapkan(kt, k2, h2); err != nil {
				return err
			}
		}
		if err := l.turunkan(ctx, kt, h2); err != nil {
			return err
		}
		h2.Pesan = pesan
		for _, x := range append([]string{"Protect.CARI1", "Protect.CARI2", "ParamInput.CARI1", "IsError"}, models.ModeLayar...) {
			if v := h.Ambil(x); v != "" {
				h2.Setel(x, v)
			}
		}
		hasil = l.layar(k2, h2, Pemegang(p, k2))
		hasil.Info = j.info
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

// validasi mengubah pesan halaman menjadi GalatValidasi (aksi dibatalkan).
func validasi(h *models.Halaman) error {
	if !h.AdaPesan() {
		return nil
	}
	return &GalatValidasi{Pesan: h.SemuaPesan()}
}

// ---------------------------------------------------------------- penangan sederhana

func halamanSaja(f func(j *jalanAksi)) penanganAksi {
	return func(j *jalanAksi) error { f(j); return nil }
}

func dariModel(f func(k *models.Konteks, h *models.Halaman) error) penanganAksi {
	return func(j *jalanAksi) error { return f(j.k, j.h) }
}

func indeks(f func(k *models.Konteks, h *models.Halaman, n int) error) penanganAksi {
	return func(j *jalanAksi) error { return f(j.k, j.h, j.r.Indeks) }
}

// mataUangBaris - CurrID baris setelah kiriman digabung (`ActivityParams(CurrID=.CurrencyID)`).
func mataUangBaris(j *jalanAksi, daftar string) string {
	d := j.h.AmbilDaftar(daftar)
	if j.r.Indeks < 1 || j.r.Indeks > len(d) {
		return ""
	}
	return d[j.r.Indeks-1]["CurrencyID"]
}

func master(j *jalanAksi) (models.MasterTreaty, error) {
	m, _, err := j.l.a.MasterTreaty(j.ctx, j.h.Ambil(models.CD+"IDMaster"))
	return m, err
}

// penangan - aksi -> port activity. Nama = nama activity tanpa akhiran `_Act`/`_act`.
var penangan map[string]penanganAksi

func init() {
	penangan = map[string]penanganAksi{
		// registrasi (tiket 01-04)
		"SetValueToClaim":     aksiPilihMaster,
		"CheckNoPolicy":       aksiPilihPolis,
		"GetNameCauseofLoss":  aksiPilihSebab,
		"SetEndDate":          halamanSaja(func(j *jalanAksi) { models.SetEndDate(j.h) }),
		"CheckPeriodPolicy":   halamanSaja(func(j *jalanAksi) { models.CheckPeriodPolicy(j.h) }),
		"CheckDateDOL":        dariModel(models.CheckDateDOL),
		"CheckReportDate":     halamanSaja(func(j *jalanAksi) { models.CheckReportDate(j.k, j.h) }),
		"CheckDateReceived":   halamanSaja(func(j *jalanAksi) { models.CheckDateReceived(j.k, j.h) }),
		"GetReportStatus":     dariModel(models.GetReportStatus),
		"SetDefNonCatastrope": halamanSaja(func(j *jalanAksi) { models.SetDefNonCatastrope(j.h) }),
		"SetEditCatastrope":   halamanSaja(func(j *jalanAksi) { models.SetEditCatastrope(j.h, j.r.Param) }),
		"BukaKatastrofe":      halamanSaja(func(j *jalanAksi) {}),
		"SetCatastrope":       aksiPilihKatastrofe,
		"SaveCatasrtope":      aksiSimpanKatastrofe,
		"SetConsultant":       dariModel(models.SetConsultant),
		"SetAdjsuter":         dariModel(models.SetAdjuster),
		"MakeLowercase":       halamanSaja(func(j *jalanAksi) { models.MakeLowercase(j.h) }),
		"GetAdders":           dariModel(models.GetAdders),
		"GetRNMShareTreaty": func(j *jalanAksi) error {
			m, err := master(j)
			models.PilihanShareRNM(j.h, m)
			return err
		},
		"DisableEditRNMShare": aksiShareRNM,
		// interest, deductible, claim amount, loss allocation, estimasi, spreading (tiket 05-07, 09)
		"AddInterest": halamanSaja(func(j *jalanAksi) { models.AddInterest(j.k, j.h) }),
		"SetCurencyInterest": func(j *jalanAksi) error {
			return models.SetCurencyInterest(j.k, j.lama, j.h, j.r.Indeks, mataUangBaris(j, models.DaftarInterest))
		},
		"CountTotalInsterest": func(j *jalanAksi) error { return models.CountTotalInsterest(j.lama, j.h, j.r.Indeks) },
		"DeleteInterest":      indeks(models.DeleteInterest),
		"SetFormat":           halamanSaja(func(j *jalanAksi) { models.SetFormat(j.k, j.h) }),
		"CountDeductible":     func(j *jalanAksi) error { return models.HitungTurunan(j.h) },
		"AddListClaimAmount":  dariModel(models.AddListClaimAmount),
		"SetCurencyList": func(j *jalanAksi) error {
			return models.SetCurencyList(j.k, j.h, j.r.Indeks, mataUangBaris(j, models.DaftarClaimAmount))
		},
		"CountListClaimAmountIDR": dariModel(models.CountListClaimAmountIDR),
		"DeleteListClaim":         indeks(models.DeleteListClaim),
		"AddLossAllocation":       dariModel(models.AddLossAllocation),
		"SetCurrency": func(j *jalanAksi) error {
			return models.SetCurrency(j.k, j.h, j.r.Indeks, mataUangBaris(j, models.DaftarLossAlloc))
		},
		"SetNameTreaty":       func(j *jalanAksi) error { return models.SetNameTreaty(j.k, j.h, j.r.Indeks) },
		"CountPersen":         dariModel(models.CountPersen),
		"RemoveLossAlloction": indeks(models.RemoveLossAlloction),
		"AddEstimation":       dariModel(models.AddEstimation),
		"CheckEstimateDate":   halamanSaja(func(j *jalanAksi) { models.CheckEstimateDate(j.k, j.h) }),
		"CurencyEstimation": func(j *jalanAksi) error {
			return models.CurencyEstimation(j.k, j.h, j.r.Indeks, mataUangBaris(j, models.DaftarEstimasi))
		},
		"CountEstimation": dariModel(models.CountEstimation),
		"DeleteEstimation": func(j *jalanAksi) error {
			m, err := master(j)
			if err != nil {
				return err
			}
			return models.DeleteEstimation(j.k, j.h, j.r.Indeks, m)
		},
		// tabel bawah = SpreadingList master (SetTreatyNameSpreading_Act langkah 11); Add / Delete aktif (keputusan work
		// owner 08-10-2026; AddSpreading_Act / DeleteSpreading_Act tidak diekspor)
		"SetTreatyNameSpreading": func(j *jalanAksi) error {
			m, err := master(j)
			if err != nil {
				return err
			}
			return models.SetTreatyNameSpreading(j.k, j.h, j.r.Indeks, m)
		},
		"AddSpreading": func(j *jalanAksi) error {
			m, err := master(j)
			if err != nil {
				return err
			}
			return models.AddSpreading(j.k, j.h, m)
		},
		"DeleteSpreading": func(j *jalanAksi) error {
			m, err := master(j)
			if err != nil {
				return err
			}
			return models.DeleteSpreading(j.k, j.h, j.r.Indeks, m)
		},
		"CountSpreading": dariModel(models.CountSpreading),
		// tombol Outstanding Claim (tiket 01, 07, 12)
		"SetOutstanding":    aksiSetOutstanding,
		"SaveOutstanding":   aksiSaveOutstanding,
		"TryMakePLA":        aksiPLA,
		"CheckNopolicy":     aksiKirimAkseptasi,
		"SubmitOutstanding": aksiSubmit,
		"Simpan":            halamanSaja(func(j *jalanAksi) {}),
		// akseptasi (tiket 08, 10, 11, 12, 13)
		"AddAdjustment":        halamanSaja(func(j *jalanAksi) { models.AddAdjustment(j.k, j.h) }),
		"DeleteAjsutment":      func(j *jalanAksi) error { return models.DeleteAjsutment(j.h, j.r.Indeks) },
		"SetPayableTreaty":     indeks(models.SetPayableTreaty),
		"SetNameCurrency":      aksiMataUangAdjustment,
		"SetDLACedingSOB":      func(j *jalanAksi) error { return models.SetDLACedingSOB(j.h, j.r.Indeks) },
		"PilihAllocation":      func(j *jalanAksi) error { return models.PilihAllocation(j.h, j.r.Indeks, j.r.Param) },
		"CountGrossAdjTreaty":  indeks(models.CountGrossAdjTreaty),
		"PilihRekening":        aksiPilihRekening,
		"SetKomiteNo":          halamanSaja(func(j *jalanAksi) {}),
		"BukaKomite":           aksiBukaKomite,
		"AddKomiteTreatyChild": aksiKomite,
		"HitServiceToKasir":    aksiKasir,
		"PrintDLATreatyIn":     aksiDLA,
		"CloseClaimProp":       aksiTutupKlaim,
	}
}
