package services

// Untuk apa berkas ini: TINDAKAN LAYAR - refresh berhitung, Save, popup Choose Business, Submit (finishAssignment).
// Pola: salinan `modul/nbtreatyin/backend/services/tindakan.go` (06-10-2026), isinya menurut section / flow EDM.
//
// Setiap tindakan bekerja atas halaman TERSIMPAN yang dipra-proses ulang (`siapkan`), lalu nilai kiriman layar
// DIGABUNG hanya untuk medan yang boleh diubah di posisi itu (`models.GabungMasukanLayar`).

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/repository"
)

// kerjakan memuat kasus untuk ditindak: terbuka, pelaku anggota antreannya, halaman dipra-proses dan digabung
// dengan kiriman layar, medan turunan dihitung ulang, action set sel yang dipicu kiriman diputar ulang.
func (l *Layanan) kerjakan(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) (models.Kasus, *models.Halaman, error) {
	if err := l.periksaPelaku(p); err != nil {
		return models.Kasus{}, nil, err
	}
	k, h, err := l.muat(ctx, id)
	if err != nil {
		return models.Kasus{}, nil, err
	}
	if k.Tertutup() {
		return models.Kasus{}, nil, ErrKasusTertutup
	}
	if k.GenerasiTertutup {
		return models.Kasus{}, nil, ErrGenerasiTertutup
	}
	if !anggota(p, k.PositionNote) {
		return models.Kasus{}, nil, ErrBukanAnggotaAntrean
	}
	if err := l.siapkan(ctx, p, k, h); err != nil {
		return models.Kasus{}, nil, err
	}
	pemicu := models.GabungMasukanLayar(h, masuk, k.PositionNote)
	if err := l.turunkan(ctx, h); err != nil {
		return models.Kasus{}, nil, err
	}
	if pemicu.Pajak {
		// keputusan WO 07-10-2026: pajak grid XOL NonProp baru dihitung ulang (rantai Choose Business, master sama)
		if err := models.HitungUlangPajakNonPropEDM(h, l.g.MasterEDM(ctx), l.jam()); err != nil {
			return models.Kasus{}, nil, rusak(err)
		}
	}
	if err := models.TerapkanPemicu(h, pemicu, l.jam()); err != nil {
		return models.Kasus{}, nil, err
	}
	return k, h, nil
}

// turunkan menghitung ULANG di server medan turunan dari isian - layar tidak pernah menjadi sumber nilai turunan:
//
//	SetPPNPPH langkah 1-3   STS_PKP agen sumber bisnis
//	CountNetPremi_act       NetPremium, Balance*, PPN/PPH, BrokerageFee*, DueTo (tab New Data, bukan NonProp baru)
//	CountSpreading_Act      total spreading (jumlah baris)
func (l *Layanan) turunkan(ctx context.Context, h *models.Halaman) error {
	if err := l.pasangStsPKP(ctx, h); err != nil {
		return err
	}
	if !models.PolisNonPropBaru(h) {
		if err := models.CountNetPremi(h); err != nil {
			return err
		}
	}
	return models.HitungTotalSpreading(h)
}

// pasangStsPKP = `SetPPNPPH` langkah 1-3: RD `BrowseClientName_RD` atas `QuotationData.SourceOfBusiness`.
func (l *Layanan) pasangStsPKP(ctx context.Context, h *models.Halaman) error {
	sts, err := l.g.StsPKPAgen(ctx, h.Ambil(models.HalamanPolis+".QuotationData.SourceOfBusiness"))
	if err != nil {
		return err
	}
	h.Setel(models.JalurStsPKP, sts)
	return nil
}

// tulis menjalankan fn di satu transaksi SESUDAH mengunci baris kasus dan memastikan tahapnya masih tahap yang
// dibaca `kerjakan`.
func (l *Layanan) tulis(ctx context.Context, k models.Kasus, fn func(tx *db.Tx) error) error {
	return l.g.Transaksi(ctx, func(tx *db.Tx) error {
		if err := l.g.KunciKasus(ctx, tx, k.ID, k.StatusWork); err != nil {
			return err
		}
		return fn(tx)
	})
}

// simpanGenerasi = Obj-Save pyWorkPage: halaman generasi + proyeksi selisihnya di transaksi yang SAMA
// (spec-penyimpanan ID-26, AC 25). Selisih proporsional dihitung ulang dulu (`models.HitungSelisihGenerasi`).
// `selesai` - Utility1: kunci saring NOPOLIS proyeksi mengikuti generasinya, KOSONG selama berjalan (tinjauan kode
// 06-10-2026; NOPOLIS generasi baru diisi `SetelNomorPolisSelesai`).
func (l *Layanan) simpanGenerasi(ctx context.Context, tx *db.Tx, k models.Kasus, h *models.Halaman, selesai bool) error {
	if err := models.HitungSelisihGenerasi(h); err != nil {
		return err
	}
	models.RapikanBentukSimpan(h)
	if err := l.g.SimpanHalaman(ctx, tx, k.ID, h); err != nil {
		return err
	}
	nopolis := ""
	if selesai {
		nopolis = h.Ambil(models.HalamanPolis + ".PolicyNo")
	}
	return l.g.SimpanSelisih(ctx, tx, k.ID, h, repository.KunciSelisih{
		NoPolis: nopolis,
		ProdKe:  k.ProdKe,
		EDMNo:   h.Ambil(models.HalamanPolis + ".EDMNo"),
		IDPega:  models.KunciInstans(k.ID),
		Sumber:  models.SumberGo,
	})
}

// ------------------------------------------------------------------ hitung

// PermintaanHitung - action set satu sel / tombol layar: SATU bentuk, selalu `Urutan` (pola NB).
type PermintaanHitung struct {
	Urutan []LangkahHitung `json:"urutan"`
	// Indeks - Param.Index / Param.idx, berbasis 1.
	Indeks  int             `json:"indeks"`
	Halaman *models.Halaman `json:"halaman"`
}

// LangkahHitung - satu refresh berhitung.
type LangkahHitung struct {
	// Aksi - nama aktivitas Pega tanpa akhiran `_Act`.
	Aksi string `json:"aksi"`
	// Param - parameter aktivitas (Data "Pct"/"Amount", DiscountType, Result, Overidding).
	Param string `json:"param"`
}

type aksiFn = func(l *Layanan, ctx context.Context, h *models.Halaman, param string, indeks int) error

func aksiHalaman(f func(*models.Halaman) error) aksiFn {
	return func(_ *Layanan, _ context.Context, h *models.Halaman, _ string, _ int) error { return f(h) }
}

func aksiParam(f func(*models.Halaman, string) error) aksiFn {
	return func(_ *Layanan, _ context.Context, h *models.Halaman, param string, _ int) error { return f(h, param) }
}

func aksiIndeks(f func(*models.Halaman, int) error) aksiFn {
	return func(_ *Layanan, _ context.Context, h *models.Halaman, _ string, indeks int) error {
		return f(h, indeks)
	}
}

func tanpaGalat(f func(*models.Halaman)) func(*models.Halaman) error {
	return func(h *models.Halaman) error { f(h); return nil }
}

// aksiHitung - aksi refresh yang tertulis di sel / tombol layar EDM (INVENTARIS-XML.md bab 5).
var aksiHitung = map[string]aksiFn{
	"CountOGPONP":     aksiHalaman(models.CountOGPONP),
	"CountResult1":    aksiParam(models.CountResult1),
	"CountResult1Onp": aksiParam(models.CountResult1Onp),
	"CountResult2Ogp": aksiParam(models.CountResult2Ogp),
	"CountResult2Onp": aksiParam(models.CountResult2Onp),
	"CountSpreading":  aksiIndeks(models.CountSpreading),
	"SetDueTo":        aksiHalaman(models.SetDueTo),
	"RemoveTypeTax":   aksiHalaman(tanpaGalat(models.RemoveTypeTax)),
	// tombol "Calculate Value Difference" (PropNewData2 S24) dan refresh otherSection
	// `DetailPolicyTreatyInPropValueDifference` (defer-load) sel uang tab New Data
	"EDMTCalculateTreatyDifference": aksiHalaman(models.EDMTCalculateTreatyDifference),
	"FillPaymentInstallment": func(l *Layanan, _ context.Context, h *models.Halaman, _ string, _ int) error {
		return models.FillPaymentInstallment(h, l.jam())
	},
	// sel `.Installment` S12 DetailPolicyTreatyInAddendum (NonProp baru): Param.Installment = .Installment
	"FillPaymentInstallmentEDMT": func(_ *Layanan, _ context.Context, h *models.Halaman, _ string, _ int) error {
		return models.FillPaymentInstallmentEDMT(h, h.Ambil(models.HalamanPolis+".Installment"))
	},
	// radio Approval ListSuggestEDM: refresh + Protection_Act
	"Protection": func(l *Layanan, _ context.Context, h *models.Halaman, _ string, _ int) error {
		models.ProtectionAct(h, l.jam())
		return nil
	},
	// CheckDataMkt: Obj-Browse marketing officer (dilewati bila MOID kosong); XML menutupnya Obj-Save langkah 5
	"CheckDataMkt": func(l *Layanan, ctx context.Context, h *models.Halaman, _ string, _ int) error {
		var mo models.BarisMO
		if id := h.Ambil(models.HalamanPolis + ".QuotationData.MOID"); id != "" {
			var err error
			if mo, err = l.g.MO(ctx, id); err != nil {
				return err
			}
		}
		models.TerapkanMO(h, mo)
		return nil
	},
	// sel `.TypeTax`: XML hanya postValue; ketetapan NB pajak dihitung ulang (`kerjakan`). Aksi ini putaran server.
	"HitungPajak": aksiHalaman(func(*models.Halaman) error { return nil }),
}

// Hitung menjalankan action set sel (`Urutan`) dan mengembalikan layarnya. Hasil TIDAK disimpan - kecuali
// `CheckDataMkt` (Obj-Save langkah 5).
func (l *Layanan) Hitung(ctx context.Context, p inti.Pelaku, id string, r PermintaanHitung) (Layar, error) {
	if len(r.Urutan) == 0 {
		return Layar{}, fmt.Errorf("%w: urutan hitung kosong", ErrPermintaanTidakSah)
	}
	fs := make([]aksiFn, len(r.Urutan))
	simpan := false
	for i, s := range r.Urutan {
		f, ada := aksiHitung[s.Aksi]
		if !ada {
			return Layar{}, fmt.Errorf("%w: aksi hitung %q", ErrPermintaanTidakSah, s.Aksi)
		}
		fs[i] = f
		simpan = simpan || s.Aksi == "CheckDataMkt"
	}
	k, h, err := l.kerjakan(ctx, p, id, r.Halaman)
	if err != nil {
		return Layar{}, err
	}
	for _, s := range r.Urutan {
		if !aksiTerbuka(k.PositionNote, s.Aksi, h) {
			return Layar{}, fmt.Errorf("%w: aksi %q", ErrTindakanTakAdaDiPosisi, s.Aksi)
		}
	}
	for i, s := range r.Urutan {
		if err := fs[i](l, ctx, h, s.Param, r.Indeks); err != nil {
			return Layar{}, err
		}
	}
	if simpan {
		if err := l.tulis(ctx, k, func(tx *db.Tx) error { return l.simpanGenerasi(ctx, tx, k, h, false) }); err != nil {
			return Layar{}, err
		}
	}
	return l.layar(k, h, true), nil
}

// ------------------------------------------------------------------ simpan

// SimpanDraf = tombol "Save" (`PropNewData2` / S19 `DetailPolicyTreatyInAddendum`, `click -> save`) - layar admin.
func (l *Layanan) SimpanDraf(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) (Layar, error) {
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return Layar{}, err
	}
	if k.PositionNote != models.PosisiAdmin {
		return Layar{}, ErrTindakanTakAdaDiPosisi
	}
	models.IsiNolWajibUang(h, k.PositionNote)
	if kosong := models.MedanWajibKosongSimpan(h, k.PositionNote); len(kosong) > 0 {
		return Layar{}, &GalatValidasi{Pesan: kosong}
	}
	if err := l.tulis(ctx, k, func(tx *db.Tx) error { return l.simpanGenerasi(ctx, tx, k, h, false) }); err != nil {
		return Layar{}, err
	}
	return l.layar(k, h, true), nil
}

// ------------------------------------------------------------------ pilih bisnis

// bolehPilihBisnis - tombol `Choose Business` (`DetailPolicyTreatyInAddendum`) tampil bila PositionNote =
// ReasTreatyInAdmin.
func bolehPilihBisnis(k models.Kasus) error {
	if k.PositionNote != models.PosisiAdmin {
		return ErrTindakanTakAdaDiPosisi
	}
	return nil
}

// DaftarBisnis - isi grid popup `BusinessAndSOBListEDM` (repository.DaftarBisnisEDM). Tanpa simpan
// (showHarness `pySubmitData` - kiriman layar digabung seperti tindakan lain).
func (l *Layanan) DaftarBisnis(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) ([]models.Baris, error) {
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return nil, err
	}
	if err := bolehPilihBisnis(k); err != nil {
		return nil, err
	}
	return l.g.DaftarBisnisEDM(ctx, repository.SaringanPopupEDM{
		OldNoOffer: h.Ambil(models.HalamanPolis + ".OldData.NoOffer"),
		EDMType:    h.Ambil(models.HalamanPolis + ".EDMType"),
		Retro:      h.Ambil(models.HalamanPolis+".ClaimType") == models.KlaimXOLRetro,
	})
}

// PilihBisnis = tombol "Choose" grid popup -> `Activity/EDMChooseBusiness_Act` langkah 1-12
// (`models.PilihBisnisEDM`) lalu 13-14 Obj-Save + Commit.
//
// ⚠️ `[terverifikasi]` rantai itu TIDAK membaca `.ID` baris yang diklik: master ditentukan
// `PolicyTreatyIn.OldData.NoOffer` (SetValueEDM_Act 4-7). Baris apa pun menghasilkan muatan yang sama - ditiru
// apa adanya (pertanyaan untuk work owner).
func (l *Layanan) PilihBisnis(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) (Layar, error) {
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return Layar{}, err
	}
	if err := bolehPilihBisnis(k); err != nil {
		return Layar{}, err
	}
	if err := models.PilihBisnisEDM(h, l.g.MasterEDM(ctx), l.jam()); err != nil {
		return Layar{}, rusak(err)
	}
	if err := l.tulis(ctx, k, func(tx *db.Tx) error { return l.simpanGenerasi(ctx, tx, k, h, false) }); err != nil {
		return Layar{}, err
	}
	return l.layar(k, h, true), nil
}

// ------------------------------------------------------------------ kirim

// Kirim = `finishAssignment` flow action posisi kasus.
//
//	Admin  `InboxPolicyTreatyInAddendum`: Submit (IsApproved 1: SetDueTo_act lalu finish bila Protect.CARI1/2 = 0;
//	       0: `PolicyTreatyInDeclineConfirm`, local action SetDueTo_act, Yes = finish) -> pasca-DT
//	       `InboxPolicyTreatyInAddendum_postDT` -> pasca-activity `InsertHistoryAkseptasiPega`
//	Atasan `DeptHeadTreatyIn_UWAddendum`: Submit -> pasca-DT `InboxPolicyTreatyIn_UW_postDT` ->
//	       `InsertHistoryAkseptasiPega`; Dept Head IsApproved 1: popup `ShowPolicyNoTreaty` (OK = finish)
//
// lalu connector flow (`models.Langkah`). Dept Head menyetujui = Utility1 `SaveJsonPolisTreatyInEDM_Act` di
// transaksi yang sama (NOPOLIS generasi, json_polis tanpa DATA_JSON, ACHIEVEMENT, TREATYINPRODUCTION). Sesudah
// transaksi: Utility2 `serviceInsertArasapas_act` (`konversikan`) - gagalnya tidak membatalkan apa pun.
// ⚠️ SuggestList -> HISTORYAKSEPTASIPRODUCTION (ketetapan NB K4; XML EDM tidak menulisnya - models/usulan.go).
func (l *Layanan) Kirim(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) (HasilKirim, error) {
	k, err := l.kirim(ctx, p, id, masuk)
	if err != nil {
		return HasilKirim{}, err
	}
	hasil := HasilKirim{Kasus: k}
	if k.StatusWork == models.StatusSelesai {
		h, err := l.g.BacaHalaman(ctx, nil, id)
		if err != nil {
			log.Printf("edmtreatyin: konversi %s tidak disusun: %v", id, err)
			hasil.PesanKonversi = models.PesanGagalKonversi
			return hasil, nil
		}
		hasil.PesanKonversi = l.konversikan(ctx, id, h, k.TglCreate)
	}
	return hasil, nil
}

func (l *Layanan) kirim(ctx context.Context, p inti.Pelaku, id string, masuk *models.Halaman) (models.Kasus, error) {
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return models.Kasus{}, err
	}
	posisi := k.PositionNote
	if models.TombolUntuk(h, posisi) == models.TombolTidakAda {
		return models.Kasus{}, &GalatValidasi{Pesan: []string{"Approval"}}
	}
	if err := models.SetDueTo(h); err != nil {
		return models.Kasus{}, err
	}
	if err := l.pasangStsPKP(ctx, h); err != nil {
		return models.Kasus{}, err
	}
	if err := l.validasiKirim(h, posisi); err != nil {
		return models.Kasus{}, err
	}
	nama, err := l.g.NamaTampilan(ctx, p.AkunID)
	if err != nil {
		return models.Kasus{}, err
	}
	if posisi == models.PosisiAdmin {
		models.PascaAdmin(h, p.AkunID, nama)
	} else {
		models.PascaAtasan(h, p.AkunID)
	}
	tr, err := models.Langkah(posisi, h.Ambil(models.HalamanPolis+".IsApproved"))
	if err != nil {
		return models.Kasus{}, err
	}
	namaKotak := ""
	if tr.KembaliKePembuat {
		if namaKotak, err = l.g.NamaTampilan(ctx, k.CreateOp); err != nil {
			return models.Kasus{}, err
		}
	}
	if tr.NBStatusKePosisi || (tr.KembaliKePembuat && strings.TrimSpace(namaKotak) == "") {
		pk, err := l.g.PemegangKotakMasuk(ctx, tr.PosisiBaru)
		if err != nil {
			return models.Kasus{}, err
		}
		namaKotak = models.NamaKotakMasuk(tr.PosisiBaru, pk)
	}
	err = l.tulis(ctx, k, func(tx *db.Tx) error {
		// InsertHistoryAkseptasiPega (pasca-activity kedua flow action) - WORKBASKET = posisi putusan diambil.
		if err := l.g.CatatRiwayat(ctx, tx, models.Riwayat{
			IDPega:     models.KunciInstans(id),
			Status:     models.StatusRiwayat(h.Ambil(models.HalamanPolis + ".IsApproved")),
			Username:   nama,
			Workbasket: posisi,
			OperatorID: p.AkunID,
		}); err != nil {
			return err
		}
		if baru := models.UsulanBelumTersimpan(h); len(baru) > 0 {
			if err := l.g.CatatUsulan(ctx, tx, models.KunciInstans(id), baru); err != nil {
				return err
			}
		}
		if tr.Ditutup() {
			if !tr.Simpan {
				// admin menolak: Decision3 No -> End3, tanpa Utility1 - generasi dilepas dari rantai polis
				if err := l.simpanGenerasi(ctx, tx, k, h, false); err != nil {
					return err
				}
				if err := l.g.LepasGenerasi(ctx, tx, id); err != nil {
					return err
				}
				// generasi lepas tidak punya "generasi sebelumnya" lagi - proyeksi selisihnya dibuang
				if err := l.g.BuangSelisihGenerasi(ctx, tx, id); err != nil {
					return err
				}
				return l.g.TutupKasus(ctx, tx, id, k.StatusWork, tr.StatusTutup)
			}
			// Utility1 SaveJsonPolisTreatyInEDM_Act
			hari, err := l.g.HariClosing(ctx, tx)
			if err != nil {
				return err
			}
			models.PrasimpanPolis(h, id, l.jam(), hari)
			// produksi = selisih: dihitung ulang SEBELUM baris produksi disusun (tinjauan kode 06-10-2026)
			if err := models.HitungSelisihGenerasi(h); err != nil {
				return err
			}
			simpanan, err := models.SusunSimpananPolis(h, id, p.AkunID)
			if err != nil {
				return err
			}
			if err := l.simpanGenerasi(ctx, tx, k, h, true); err != nil {
				return err
			}
			if err := l.g.SetelNomorPolisSelesai(ctx, tx, id, h.Ambil(models.HalamanPolis+".PolicyNo")); err != nil {
				return err
			}
			if err := l.g.SimpanPolisProduksi(ctx, tx, simpanan); err != nil {
				return err
			}
			return l.g.TutupKasus(ctx, tx, id, k.StatusWork, tr.StatusTutup)
		}
		switch {
		case tr.NBStatusKePosisi:
			h.Setel("NBStatus", models.TeksNBStatusKotakMasuk(namaKotak))
		case tr.KembaliKePembuat:
			h.Setel("NBStatus", models.TeksNBStatusKembali(namaKotak))
		}
		h.Setel("PositionNote", tr.PosisiBaru)
		if err := l.simpanGenerasi(ctx, tx, k, h, false); err != nil {
			return err
		}
		return l.g.PindahPosisi(ctx, tx, id, k.StatusWork, tr.PosisiBaru, tr.PositionBaru)
	})
	if err != nil {
		return models.Kasus{}, err
	}
	return l.g.Keadaan(ctx, nil, id)
}

// validasiKirim - medan wajib layar dan pesan validasi yang dipasang rantai layar (Property-Set-Messages /
// Page-Set-Messages): halaman berpesan tidak dapat di-submit, sama dengan Pega.
//
//	semua posisi  `Protection_Act` (radio Approval ListSuggestEDM) atas SALINAN: pesannya saja
//	admin         `CountOGPONP_Act` atas SALINAN (tab New Data, bukan NonProp baru) - pesan validasi uang
//	admin setuju  ⛔ `[keputusan work owner]` spec-penyimpanan AC 8 / ID-15: generasi baru yang kehilangan baris
//	              spreading generasi sebelumnya ditolak (`models.BarisSpreadingHilang`)
func (l *Layanan) validasiKirim(h *models.Halaman, posisi string) error {
	models.IsiNolWajibUang(h, posisi)
	pesan := models.MedanWajibKosong(h, posisi)
	salin := h.Salin()
	salin.BersihkanPesan()
	if posisi == models.PosisiAdmin && !models.PolisNonPropBaru(h) {
		if err := models.CountOGPONP(salin); err != nil {
			return err
		}
	}
	models.ProtectionAct(salin, l.jam())
	pesan = append(pesan, salin.SemuaPesan()...)
	if posisi == models.PosisiAdmin && models.Disetujui(h.Ambil(models.HalamanPolis+".IsApproved")) {
		pesan = append(pesan, models.BarisSpreadingHilang(h)...)
	}
	// ProductionDate yang diisi Protection_Act ikut halaman yang disimpan
	if v := salin.Ambil(models.HalamanPolis + ".ProductionDate"); h.Ambil(models.HalamanPolis+".ProductionDate") == "" {
		h.Setel(models.HalamanPolis+".ProductionDate", v)
	}
	if len(pesan) > 0 {
		return &GalatValidasi{Pesan: pesan}
	}
	return nil
}

// rusak membungkus galat hitung atas nilai master (bukan angka, bagi nol) - sumbernya dokumen master.
func rusak(err error) error {
	if err == nil || errors.Is(err, ErrMasterXOLRusak) {
		return err
	}
	if errors.Is(err, models.ErrBukanAngka) || errors.Is(err, models.ErrBagiNol) {
		return fmt.Errorf("%w: %w", ErrMasterXOLRusak, err)
	}
	return err
}
