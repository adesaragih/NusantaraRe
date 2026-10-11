package services

// Untuk apa berkas ini: PENANGAN AKSI YANG MENULIS LEBIH DARI HALAMAN atau MEMBACA ULANG PILIHAN - pilihan pop-up /
// autocomplete yang dibaca ulang di server (polis, ceding, wilayah, cause of loss, katastrofe), Submit Input Register
// (finishAssignment), Download Claim Face Sheet (`CLaimFaceSheet_Act`), Print PLA (`SetIndexObj_Act` +
// `GeneratePLA`), Send to PIC Claim (`SetDisable_ACT` + finishAssignment), Back (`BackFromRegister` +
// `BackToRegister_act` + finishAssignment).

import (
	"fmt"
	"strings"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// pesanKosong - padanan validasi klien Pega (`pyClientValidation=true`, medan wajib kosong). `[dugaan]` teks bawaan
// Pega "Value cannot be blank"; rule pesannya tidak ada di ekspor (pola Claim Prop).
const pesanKosong = "Value cannot be blank"

func wajibTerisi(j *jalanAksi, ts []models.Tata) error {
	kosong := models.WajibKosong(j.h, ts, nil)
	if len(kosong) == 0 {
		return nil
	}
	var p []string
	for _, l := range kosong {
		p = append(p, l+": "+pesanKosong)
	}
	return &GalatValidasi{Pesan: p}
}

// pindah = finishAssignment: halaman disimpan, tahap dipindah, pra-proses tahap baru dijalankan dan disimpan,
// InsertProgressClaim tahap baru (pra-proses assignment berikutnya, Pega menjalankannya saat assignment dibuka).
func (j *jalanAksi) pindah(baru, posisi string) error {
	models.BuangTurunan(j.h)
	if err := j.l.g.SimpanHalaman(j.ctx, j.tx, j.kasus.ID, j.h); err != nil {
		return err
	}
	if err := j.l.g.PindahTahap(j.ctx, j.tx, j.kasus.ID, j.kasus.Tahap, baru, posisi, j.k.Sekarang); err != nil {
		return err
	}
	j.kasus.Tahap = baru
	j.k.Langkah = models.LabelTahap[baru]
	if err := j.l.siapkan(j.ctx, j.k, j.kasus, j.h); err != nil {
		return err
	}
	j.h.BersihkanPesan()
	if err := j.l.g.SimpanHalaman(j.ctx, j.tx, j.kasus.ID, j.h); err != nil {
		return err
	}
	j.selesai = true
	return j.l.tulisProgres(j.ctx, j.tx, j.k, j.kasus.ID, j.h, models.ParamProgres{})
}

// ---------------------------------------------------------------- pilihan pop-up

// aksiPilihPolis - tautan nomor polis grid ViewPolis (`CopyNB_Act`, Param PolicyNo + Prodke = "nopolis|prodke"). Baris
// dibaca ULANG: pasangan nomor + prodke harus ada di FACINPRODUCTION / JSON_POLIS.
func aksiPilihPolis(j *jalanAksi) error {
	nopolis, prodke, _ := strings.Cut(j.r.Param, "|")
	nopolis, prodke = strings.TrimSpace(nopolis), strings.TrimSpace(prodke)
	if nopolis == "" {
		return fmt.Errorf("%w: nomor polis kosong", ErrPermintaanTidakSah)
	}
	ada, err := j.l.a.AdaPolis(j.ctx, nopolis, prodke)
	if err != nil {
		return err
	}
	if !ada {
		return fmt.Errorf("%w: polis %q prodke %q tidak ada di FACINPRODUCTION", ErrPermintaanTidakSah, nopolis, prodke)
	}
	dok, ada, err := j.l.a.DokumenPolis(j.ctx, nopolis, prodke)
	if err != nil {
		return err
	}
	if !ada {
		return fmt.Errorf("%w: %s", ErrTertunda, OQPolisTanpaJSON)
	}
	j.bukaModal = ModalPilihPolis
	return models.PilihPolis(j.h, nopolis, prodke, dok)
}

// OQPolisTanpaJSON - polis FACINPRODUCTION tanpa dokumen JSON_POLIS (NB Fac In sistem baru tidak menulis JSON_POLIS).
const OQPolisTanpaJSON = "OQ-CFI-11: polis ini tidak punya dokumen JSON_POLIS - halaman polis (OfferFacIn) tidak dapat disalin"

// aksiPilihCeding - tombol Choose grid CedingCoList (`GetCeding_act` Param.Ceding = .CedingCoName baris).
func aksiPilihCeding(j *jalanAksi) error {
	d := j.h.AmbilDaftar(models.DaftarCedingCo)
	if j.r.Indeks < 1 || j.r.Indeks > len(d) {
		return fmt.Errorf("%w: baris ceding %d", ErrPermintaanTidakSah, j.r.Indeks)
	}
	return models.GetCeding(j.k, j.h, d[j.r.Indeks-1]["CedingCoName"])
}

// aksiPilihWilayah - autocomplete Country / Province / City / District / Region (`isi .ID -> ...ID`): pilihan dibaca
// ULANG dari master menurut induknya, nama + ID ditulis.
func aksiPilihWilayah(sumber string) penanganAksi {
	return func(j *jalanAksi) error {
		h, a, ctx, nilai := j.h, j.l.a, j.ctx, strings.TrimSpace(j.r.Param)
		var opsi []models.Pilihan
		var err error
		jalur := ""
		switch sumber {
		case models.SumberNegara:
			jalur = models.CD + "Country"
			opsi, err = a.DaftarNegara(ctx, nilai)
		case models.SumberProvinsi:
			jalur = models.CD + "Province"
			opsi, err = a.DaftarProvinsi(ctx, h.Ambil(models.CD+"Country"), nilai)
		case models.SumberKota:
			jalur = models.CD + "City"
			opsi, err = a.DaftarKota(ctx, h.Ambil(models.CD+"Province"), nilai)
		case models.SumberDistrik:
			jalur = models.CD + "District"
			opsi, err = a.DaftarDistrik(ctx, h.Ambil(models.CD+"CityID"), nilai)
		case models.SumberRW:
			jalur = models.CD + "RW"
			opsi, err = a.DaftarRW(ctx, h.Ambil(models.CD+"DistrictID"), nilai)
		}
		if err != nil {
			return err
		}
		for _, o := range opsi {
			if o.Nilai == nilai {
				h.Setel(jalur, o.Nilai)
				for p, v := range o.Tambahan {
					h.Setel(models.CD+p, v)
				}
				return nil
			}
		}
		return fmt.Errorf("%w: pilihan %s %q", ErrPermintaanTidakSah, sumber, nilai)
	}
}

// aksiCekKembar - change Claim Estimate (`CheckDoubleClaim_Act`): IsError 1 / 0.
func aksiCekKembar(j *jalanAksi) error {
	_, err := models.CheckDoubleClaim(j.ctx, j.l.a, j.h)
	return err
}

// aksiPilihSebab - tombol Choose grid "List Cause of Loss" (`GetNameCauseofLoss_Act`).
func aksiPilihSebab(j *jalanAksi) error {
	s, ada, err := j.l.a.BarisSebabID(j.ctx, j.r.Param)
	if err != nil {
		return err
	}
	if !ada {
		return fmt.Errorf("%w: cause of loss %q tidak ada", ErrPermintaanTidakSah, j.r.Param)
	}
	models.GetNameCauseofLoss(j.h, s.Description, s.ID)
	return nil
}

// aksiPilihKatastrofe - tombol Choose grid katastrofe (`SetCatastrope_act`).
func aksiPilihKatastrofe(j *jalanAksi) error {
	c, ada, err := j.l.a.BarisKatastrofeID(j.ctx, j.r.Param)
	if err != nil {
		return err
	}
	if !ada {
		return fmt.Errorf("%w: katastrofe %q tidak ada", ErrPermintaanTidakSah, j.r.Param)
	}
	models.SetCatastrope(j.h, c.ID, c.Note)
	return nil
}

// aksiSimpanKatastrofe - tombol Save form katastrofe baru (`SaveCatasrtope_Act` lalu `InputCatastrope(Cancel)`).
// Param = `ParamCat.NOTE`.
func aksiSimpanKatastrofe(j *jalanAksi) error {
	if strings.TrimSpace(j.r.Param) == "" {
		return &GalatValidasi{Pesan: []string{"Note: " + pesanKosong}}
	}
	if err := j.l.g.SisipKatastrofe(j.ctx, j.tx, models.SusunKatastrofe(j.k, j.h, j.r.Param), j.k.Sekarang); err != nil {
		return err
	}
	models.InputCatastrope(j.h, "Cancel")
	return nil
}

// ---------------------------------------------------------------- Input Register

// aksiSubmitRegister - tombol Submit Input Register (IsError 0) / Submit pop-up ProtectDOL: validasi klien (wajib),
// pre-DT `InsertObjectItemList_DT` + `CheckListEstimasi_Act` + `ProteksiDataRegister_Act` (`PascaRegister`), lalu
// finishAssignment bila `TempClaimData.ClaimData.Test != ""` -> Assignment7 Input Estimasi (worklist pembuat).
func aksiSubmitRegister(j *jalanAksi) error {
	if err := wajibTerisi(j, models.Evaluasi(j.h, models.LayarRegister(), false)); err != nil {
		return err
	}
	models.PascaRegister(j.k, j.h)
	if err := validasi(j.h); err != nil {
		return err
	}
	if j.h.Ambil(models.JalurUjiCalon) == "" {
		return nil
	}
	j.bukaModal = ""
	return j.pindah(models.TahapEstimasi, j.kasus.PembuatID)
}

// ---------------------------------------------------------------- Input Estimasi

// aksiCFS - tombol "Download Claim Face Sheet" baris objek (`CLaimFaceSheet_Act`, `models` berkas cfs.go).
func aksiCFS(j *jalanAksi) error {
	if err := j.butuh(true, false, false); err != nil {
		return err
	}
	h, k := j.h, j.k
	if err := models.ProtectDownloadFaceClaim(k, h, j.o, j.kasus.ID); err != nil { // 1
		return err
	}
	if err := validasi(h); err != nil { // ProtectDownloadFaceClaim 6: pesan menghentikan aksi
		return err
	}
	k.Kronologi(h, models.TeksCFS)          // 2-3
	if h.Ambil(models.CD+"NoClaim") == "" { // 12-17
		b, err := j.l.g.UrutNomor(j.ctx, j.tx, models.JenisNomorKlaim, k.Sekarang)
		if err != nil {
			return err
		}
		h.Setel(models.CD+"NoClaim", models.RakitNomorKlaim(b.Jenis, h.Ambil(models.OQ+"BusinessOldId"), b.MMYYYY, b.Urut))
	}
	if err := models.NamaTreatySpreading(k, h, j.o); err != nil { // 27
		return err
	}
	rows, err := models.SetelCFS(k, h, j.o, j.kasus.ID) // 43-46
	if err != nil {
		return err
	}
	for _, b := range rows { // 43.7.2 SaveCFS_ACT
		if err := j.l.g.SisipOS(j.ctx, j.tx, b, k.Sekarang); err != nil {
			return err
		}
	}
	if err := j.salinJSONKlaim(); err != nil { // 43.9
		return err
	}
	if err := j.antreKonversi(models.StsOSOutstanding); err != nil { // 47 (Connect-REST hanya IsPEGAPROD)
		return err
	}
	kunci := models.KunciInstans(j.kasus.ID)
	if err := j.l.g.CatatLogLayanan(j.ctx, j.tx, repository.LogLayanan{IDPega: kunci, // 48-49
		Parameter: kunci + " / " + h.Ambil(models.JalurNoPolis), JenisService: models.JenisServiceEstimasi},
		k.Sekarang); err != nil {
		return err
	}
	for i := range h.AmbilDaftar(models.DaftarItem(j.o)) { // 50
		if err := models.CheckLimit(k, h, j.o, i+1); err != nil {
			return err
		}
	}
	j.info = models.OQDokumenPDF
	return nil
}

// salinJSONKlaim = InsertJsonClaimNonMBU_act -> PEGA_JSON_KLAIM_PNC. MNK_NO_KLAIM dan NOPOLIS JSON_KLAIM NOT NULL: tanpa
// nomor klaim / polis procedure gagal dan galatnya ditelan (ErrMsg) - di sini tidak ditulis.
func (j *jalanAksi) salinJSONKlaim() error {
	no, pol := j.h.Ambil(models.CD+"NoClaim"), j.h.Ambil(models.JalurNoPolis)
	if no == "" || pol == "" {
		return nil
	}
	return j.l.g.SalinJSONKlaim(j.ctx, j.tx, models.KunciInstans(j.kasus.ID), no, pol, j.k.Sekarang)
}

// aksiBukaPLA - tombol "Print PLA" baris objek: `SetIndexObj_Act` + pra-proses harness `CreateRemarksNotePLA_Act`, lalu
// pop-up Pla_Dtl.
func aksiBukaPLA(j *jalanAksi) error {
	if err := j.butuh(true, false, false); err != nil {
		return err
	}
	if err := models.BukaPLA(j.h, j.o); err != nil {
		return err
	}
	j.bukaModal = fmt.Sprintf("%s:%d", AwalanModalPLA, j.o)
	return nil
}

// aksiPLA - tombol Submit pop-up Pla_Dtl (`GeneratePLA`, lalu `SendEmailDLA_ACT` yang keluar karena
// ProtectPrint.CARI50 false): nomor PLA treaty "G" per spreading treaty (GeneratePLATreaty_Act 5-13), PlaStatus 1.
func aksiPLA(j *jalanAksi) error {
	if err := j.butuh(true, false, false); err != nil {
		return err
	}
	if !pembukaTerbuka(j, "", "BukaPLA", j.o) { // modal pla:o lahir untuk setiap objek; pembukanya yang menyaring
		return fmt.Errorf("%w: Print PLA objek %d tertutup", ErrAksiTertutup, j.o)
	}
	if err := wajibTerisi(j, models.Evaluasi(j.h, models.LayarPLA(j.o), false)); err != nil {
		return err
	}
	h, k := j.h, j.k
	ada := models.NomorPLATreatyAda(h, j.o) // GeneratePLA 4.2 (sebelum loop)
	r, err := models.GeneratePLA(k, h, j.o)
	if err != nil {
		return err
	}
	for _, i := range r.Treaty {
		for _, s := range h.AmbilDaftar(models.DaftarDiItem(j.o, i, models.AnakSpreadPolis)) {
			if !models.SpreadTreatyPLA(s) {
				continue
			}
			nomor := ada
			if nomor == "" { // GeneratePLATreaty_Act 5-8
				b, err := j.l.g.UrutNomor(j.ctx, j.tx, models.JenisNomorPLATreaty, k.Sekarang)
				if err != nil {
					return err
				}
				nomor = models.RakitNomorKlaim(b.Jenis, h.Ambil(models.OQ+"BusinessOldId"), b.MMYYYY, b.Urut)
			}
			if err := models.TerapkanPLATreaty(k, h, j.o, i, nomor, s); err != nil {
				return err
			}
		}
	}
	if err := models.SelesaiPLA(h, j.o); err != nil { // 9
		return err
	}
	j.info = models.OQDokumenPDF
	if len(r.FacOQ) > 0 {
		// `[inferensi]` (PARITAS, OQ-CFI-21): ObjectItem.GeneratePLA (PLA fac retro) tidak diekspor; pasangannya
		// GeneratePLATreaty_Act langkah 19 menyetel IsPicTransfer 1 - tanpa itu tombol Send to PIC Claim kasus retro
		// fakultatif tidak pernah terbuka.
		h.Setel(models.JalurIsPicTransfer, "1")
		j.info = models.OQPLAFac
	}
	return nil
}

// aksiKirimPIC - tombol "Send to PIC Claim" (`SetDisable_ACT`), lalu finishAssignment bila objek pertama bukan retro
// atau PLA-nya sudah dicetak -> Decision5 (pyNote kosong) -> Assignment3 Choose Surveyor (workbasket ReasKlaimTeknik).
func aksiKirimPIC(j *jalanAksi) error {
	models.SetDisable(j.k, j.h)
	if err := validasi(j.h); err != nil {
		return err
	}
	if !models.BolehKirimPIC(j.h) {
		return nil
	}
	return j.pindah(models.TahapSurveyor, models.WorkbasketSurveyor)
}

// aksiKembaliRegister - tombol Back layar Input Estimasi (pre-DT `BackFromRegister`, `BackToRegister_act`,
// finishAssignment) -> Decision5 IsBackStage (pyNote "Back") -> Assignment1 Input Register.
func aksiKembaliRegister(j *jalanAksi) error {
	models.BackFromRegister(j.k, j.h)
	return j.pindah(models.TahapRegister, j.kasus.PembuatID)
}
