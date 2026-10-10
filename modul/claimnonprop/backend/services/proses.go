package services

// Untuk apa berkas ini: PENANGAN AKSI YANG MENULIS LEBIH DARI HALAMAN - pilihan pop-up yang dibaca ulang di server, Save
// (`SaveDataToJClaim_Act`), Save to issue RNM (`SaveDataToOSAksep_Act`), Print CFS (`GenerateCFS_act`), Print PLA
// (`GeneratePlaCNP_Act`), Submit (finishAssignment OutstandingClaim), Save To OS (`SaveToOS`), akseptasi baru
// (`AddAkseptasiCNP_Act`), rekening, penyerahan komite (`ProteksiSendKomiteCNP_Act` + `CreateChildKomiteCNP_Act`),
// Acceptation (`HitServiceToKasir_Act`), Save Previously Paid (`SaveCNPLayerList_Act`), Close Claim
// (`CloseClaimTNonProp`), Close Without Payment (`CloseClaimNP_preAct`).

import (
	"errors"
	"fmt"
	"strings"

	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/repository"
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

// ---------------------------------------------------------------- pilihan pop-up

// aksiPilihMaster - tombol Choose popup ChooseMasterTNonProp (`SetValueClaimTNP_Act`). Param =
// TREATYID|CLASSOFBUSINESSID|TREATYGROUP|sumber (IN / INEDM); baris dibaca ULANG dari master (layar tidak menulis medan
// master). Sumber INEDM hanya bila tombol Choose Master (INEDM) yang terbuka.
func aksiPilihMaster(j *jalanAksi) error {
	p := strings.SplitN(j.r.Param, "|", 4)
	for len(p) < 4 {
		p = append(p, "")
	}
	sumber := models.SumberMasterIN
	if p[3] == models.SumberMasterINEDM {
		sumber = models.SumberMasterINEDM
	}
	b, ada, err := j.l.a.BarisMasterDari(j.ctx, sumber, p[0], p[1], p[2])
	if err != nil {
		return err
	}
	if !ada {
		return fmt.Errorf("%w: master treaty %q tidak ada", ErrPermintaanTidakSah, p[0])
	}
	return models.PilihMaster(j.k, j.h, models.ParamPilihMaster{ID: b.TreatyID, BusinessName: b.ClassOfBusiness,
		BusinessCode: b.ClassOfBusinessID, TreatyGroup: b.TreatyGroup, TreatyName: b.TreatyContractName})
}

// aksiPilihPolis - autocomplete Policy No / tautan popup ViewListPolicyCNP (`CheckNoPolicy`, NOPOLIS + COB = CARI12).
// Nomor polis harus ada di daftar GetDataPolisNonProp_SQL master ini; COB dibaca dari barisnya.
func aksiPilihPolis(j *jalanAksi) error {
	nopol := strings.TrimSpace(j.r.Param)
	if nopol == "" {
		nopol = j.h.Ambil(models.CD + "PolicyData.PolicyNo")
	}
	rows, err := j.l.a.DaftarPolis(j.ctx, j.h.Ambil(models.CD+"IDMaster"))
	if err != nil {
		return err
	}
	for _, r := range rows {
		if r.PolicyNo == nopol {
			models.CheckNoPolicy(j.h, r.PolicyNo, r.TreatyGroup)
			return models.SetMOClaimTreaty(j.k, j.h)
		}
	}
	return fmt.Errorf("%w: polis %q bukan polis master ini", ErrPermintaanTidakSah, nopol)
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

// ---------------------------------------------------------------- Outstanding Claim

// salinJSONKlaim = InsertJsonClaimTreaty_act -> PEGA_JSON_KLAIM_PNC. MNK_NO_KLAIM dan NOPOLIS JSON_KLAIM NOT NULL: tanpa
// nomor klaim / polis procedure gagal dan galatnya ditelan (ErrMsg) - di sini tidak ditulis.
func (j *jalanAksi) salinJSONKlaim() error {
	no, pol := j.h.Ambil(models.CD+"NoClaim"), j.h.Ambil(models.CD+"PolicyData.PolicyNo")
	if no == "" || pol == "" {
		return nil
	}
	return j.l.g.SalinJSONKlaim(j.ctx, j.tx, models.KunciInstans(j.kasus.ID), no, pol, j.k.Sekarang)
}

// aksiSimpan - tombol "Save" (`SaveDataToJClaim_Act` = InsertJsonClaimTreaty_act, lalu save).
func aksiSimpan(j *jalanAksi) error { return j.salinJSONKlaim() }

// aksiSimpanOS - tombol "Save to issue RNM" (`SaveDataToOSAksep_Act`).
func aksiSimpanOS(j *jalanAksi) error {
	h, k := j.h, j.k
	h.Setel("FlagPrintPla.CARI28", "0") // 2
	oldID, err := j.l.a.KodeLamaBisnis(j.ctx, h.Ambil(models.OQ+"BusinessName"))
	if err != nil {
		return err
	}
	noOS, err := j.l.g.NomorKlaimOS(j.ctx, j.tx, j.kasus.ID) // 10
	if err != nil {
		return err
	}
	if !models.PeriksaSimpanOS(h, oldID, noOS) { // 3-13
		return validasi(h)
	}
	if h.Ambil(models.CD+"PolicyData.PolicyNo") == "" { // 15.4 -> END (26)
		h.TambahPesan(models.CD+"PolicyData.PolicyNo", models.PesanNoPolisKosong)
		return validasi(h)
	}
	if h.Ambil(models.CD+"NoClaim") == "" { // 15.6
		b, err := j.l.g.UrutNomor(j.ctx, j.tx, models.JenisNomorKlaim, k.Sekarang)
		if err != nil {
			return err
		}
		h.Setel(models.CD+"NoClaim", models.RakitNomorKlaim(b.Jenis, h.Ambil(models.CD+"QuotationData.BusinessOldId"),
			b.MMYYYY, b.Urut))
	}
	for _, l := range models.KelompokLayerOS(h) { // 15.7-15.8
		p := models.ParamOSLayer(h, l)
		lama, err := j.l.g.JumlahOS(j.ctx, j.tx, j.kasus.ID, p["TypeLoss"], p["Currency"])
		if err != nil {
			return err
		}
		if err := models.KurangiOS(h, p, lama); err != nil {
			return err
		}
		if models.NilaiNol(p["Value"]) { // 15.8.3 trans: Value 0 -> label B (baris tidak ditulis)
			continue
		}
		if err := j.l.g.SisipOS(j.ctx, j.tx, models.SusunBarisOS(h, j.kasus.ID, models.StsOSOutstanding, p), k.Sekarang); err != nil {
			return err
		}
		if err := j.salinJSONKlaim(); err != nil { // 15.8.7
			return err
		}
		if err := j.antre(JenisEfekOutstanding, j.kasus.ID, map[string]string{"layer": p["TypeLoss"],
			"currency": p["Currency"]}); err != nil { // 15.8.8 InsertClaimOutstanding_NP (IsPEGAPROD)
			return err
		}
		h.Setel("FlagPrintPla.CARI28", "1") // 15.8.9
	}
	models.TandaiOutstanding(h) // 16, 19-22
	return j.salinJSONKlaim()   // 25
}

// aksiCFS - tombol "Print CFS" (`GenerateCFS_act`). Berkas PDF = OQ-CNP-22.
func aksiCFS(j *jalanAksi) error {
	if _, ok, err := models.CetakCFS(j.h); err != nil || !ok {
		if err != nil {
			return err
		}
		return validasi(j.h)
	}
	j.info = models.PesanBerkasOQ
	return nil
}

// aksiPLA - local action GeneratePLACNP (`GeneratePlaCNP_Act`): nomor PLA (revisi bila nilai klaim berubah sejak nomor
// terbit), kronologi "PRINT PLA", IsPLA = 1, JSON_KLAIM. Berkas PDF per reasuradur = OQ-CNP-22. ⚠️ Nomor memakai
// `OfferFacIn.QuotationData.BusinessOldId` (langkah 6) - properti tanpa penulis di korpus (PARITAS, OQ); PLATNP_SEQ
// tidak ada di DEV (OQ DBA) - aksi berhenti terang.
func aksiPLA(j *jalanAksi) error {
	h := j.h
	switch no := h.Ambil(models.CD + "NoPla"); {
	case no != "" && h.Ambil("FlagPrintPla.CARI28") == "1": // 7
		h.Setel(models.CD+"NoPla", models.RevisiNomorPLA(no))
	case no == "": // 8-9
		nomor, err := j.l.g.NomorPLA(j.ctx, j.tx, h.Ambil(models.OQ+"BusinessOldId"), j.k.Sekarang)
		if err != nil {
			if strings.Contains(err.Error(), "ORA-02289") {
				return fmt.Errorf("%w: %s", ErrTertunda, OQSequencePLA)
			}
			return err
		}
		h.Setel(models.CD+"NoPla", nomor)
	}
	j.k.Riwayat(h, TeksPrintPLA)    // 19-20
	h.Setel(models.CD+"IsPLA", "1") // 19
	j.info = models.PesanBerkasOQ
	return j.salinJSONKlaim() // 21
}

// TeksPrintPLA - GeneratePlaCNP_Act langkah 19 (VERBATIM).
const TeksPrintPLA = "PRINT PLA"

// OQSequencePLA - sequence GenerateNoPLATNP tidak ada di DEV.
const OQSequencePLA = "OQ-CNP-39: sequence PLATNP_SEQ (GenerateNoPLATNP) tidak ada di DEV - nomor PLA belum dapat diterbitkan"

// aksiSubmit - tombol "Submit" (finishAssignment FlowAction OutstandingClaim): validasi klien (wajib), post-activity
// `InputOutStandingCTNP_PostAct`, lalu Assignment1 Input Acceptation (workbasket ReasKlaimTeknik, OQ-CNP-08).
func aksiSubmit(j *jalanAksi) error {
	if err := wajibTerisi(j, models.Evaluasi(j.h, models.LayarOutstanding(), false)); err != nil {
		return err
	}
	models.PascaOutstanding(j.h)
	if err := validasi(j.h); err != nil {
		return err
	}
	return j.l.g.PindahTahap(j.ctx, j.tx, j.kasus.ID, models.TahapOutstanding, models.TahapAcceptation,
		models.WorkbasketAcceptation, j.k.Sekarang)
}

// ---------------------------------------------------------------- Input Acceptation

// aksiSaveToOS - tombol "Save To OS" (`SaveToOS`).
func aksiSaveToOS(j *jalanAksi) error {
	h, k := j.h, j.k
	h.BersihkanPesan()                                  // 3
	h.Setel("FlagPrintPla.CARI28", "0")                 // 4
	if h.Ambil(models.CD+"PolicyData.PolicyNo") == "" { // 6.1 -> END (11)
		h.TambahPesan(models.CD+"PolicyData.PolicyNo", models.PesanNoPolisKosong)
		return validasi(h)
	}
	trtSebelum := ""
	for _, l := range models.KelompokLayerOS(h) { // 6.2-6.3
		lama, err := j.l.g.JumlahOS(j.ctx, j.tx, j.kasus.ID, l["TreatyName"], l["Currency"])
		if err != nil {
			return err
		}
		p, tulis, err := models.ParamOSSaveToOS(h, l, trtSebelum, lama)
		if err != nil {
			return err
		}
		trtSebelum = l["TreatyName"] // 6.3.6
		if !tulis {
			continue
		}
		if err := j.l.g.SisipOS(j.ctx, j.tx, models.SusunBarisOS(h, j.kasus.ID, models.StsOSOutstanding, p), k.Sekarang); err != nil {
			return err
		}
		if err := j.antreKonversi(models.StsOSOutstanding); err != nil { // 6.3.10 (hardcode URL / akun dibuang, OQ-CNP-03)
			return err
		}
		if !models.NilaiNol(p["Value"]) { // 6.3.11
			h.Setel("FlagPrintPla.CARI28", "1")
		}
	}
	if _, _, err := models.CetakCFS(h); err != nil { // 7 GenerateCFS_act
		return err
	}
	models.SelesaiSaveToOS(h) // 8-9
	return nil
}

// aksiTambahAkseptasi - tombol Add grid Acceptation List (`AddAkseptasiCNP_Act`).
func aksiTambahAkseptasi(j *jalanAksi) error {
	ok, err := models.AddAkseptasi(j.k, j.h, j.m)
	if err != nil {
		return err
	}
	if !ok {
		return validasi(j.h)
	}
	return nil
}

// rekeningAkseptasi - calon rekening "Name of Bank" akseptasi `n` (halaman `Result` SetPayableTreatyNP_Act: Payable 1 =
// rekening ceding, 2 = SOB, 3 = menurut nama penerima; satu mata uang = CURRENCYID Claim Acceptation pertama, multi
// mata uang = semua rekening klien).
func (j *jalanAksi) rekeningAkseptasi(n int) ([]models.RekeningBank, error) {
	h := j.h
	multi := models.MultiMataUang(h)
	cur := ""
	if acc := h.AmbilDaftar(models.JalurAdj(n, models.AnakClaimAccept)); len(acc) > 0 && !multi {
		cur = acc[0]["CurrencyID"]
	}
	var rek []models.RekeningBank
	var err error
	switch h.Ambil(models.CD + "Payable") {
	case "3":
		rek, err = j.l.a.RekeningBankNama(j.ctx, h.Ambil(models.CD+"PayableTo"), cur)
	case "1", "2":
		klien := h.Ambil(models.TM + "LeadingReinsSourceID")
		if h.Ambil(models.CD+"Payable") == "1" {
			klien = h.Ambil(models.TM + "CedingID")
		}
		if multi {
			rek, err = j.l.a.RekeningBankKlien(j.ctx, klien)
		} else {
			rek, err = j.l.a.RekeningBank(j.ctx, klien, cur)
		}
	}
	return models.RekeningSah(rek), err
}

// aksiPilihRekening - pilihan autocomplete "Name of Bank" (bank 1 / bank 2) lalu `SetAccoutNo_Act`. Param =
// ACCOUNTNO|NAMEOFBANK; rekening dibaca ulang dari calon akseptasi itu.
func aksiPilihRekening(j *jalanAksi, sfx string) error {
	n := j.r.Indeks
	b, err := adjBaris(j.h, n)
	if err != nil {
		return err
	}
	p := strings.SplitN(j.r.Param, "|", 2)
	for len(p) < 2 {
		p = append(p, "")
	}
	rek, err := j.rekeningAkseptasi(n)
	if err != nil {
		return err
	}
	for _, r := range rek {
		if r.AccountNo == p[0] && r.NameOfBank == p[1] {
			models.SalinRekening(j.h, b, r, sfx)
			return models.SetAccountNo(j.h, n)
		}
	}
	return fmt.Errorf("%w: rekening bukan pilihan akseptasi ini", ErrPermintaanTidakSah)
}

// aksiBukaKomite - tombol "Send to Committe": SetAccoutNo_Act, ProteksiSendKomiteCNP_Act, ProteksiNilaiClaim, lalu
// harness KomiteCNP (tanggal, inisial; gerbang tombol Send = `CekError.Lolos`).
func aksiBukaKomite(j *jalanAksi) error {
	n := j.r.Indeks
	if err := models.SetAccountNo(j.h, n); err != nil {
		return err
	}
	email, err := j.l.a.EmailPelaku(j.ctx, j.k.Pelaku)
	if err != nil {
		return err
	}
	c, err := models.ProteksiKirimKomite(j.k, j.h, n, email)
	if err != nil {
		return err
	}
	nama, err := j.l.a.NamaPelaku(j.ctx, j.k.Pelaku)
	if err != nil {
		return err
	}
	j.h.Setel(models.JalurLolosKomite, nolSatu(c.Lolos(j.h)))
	j.h.Setel(models.JalurTanggalKomite, j.k.Hari())
	j.h.Setel(models.JalurPICKomite, nama)
	j.bukaModal = fmt.Sprintf("komite:%d", n)
	return nil
}

func nolSatu(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// aksiKomite - tombol "Send Claim to Committee" (`CreateChildKomiteCNP_Act`). Gerbang DITEGAKKAN di layanan: proteksi
// dihitung ulang pada salinan halaman (tanpa efek kronologi), validasi Gross Value vs Spreading In SEBELUM kasus komite
// dibuat (perbaikan OQ-CNP-05 butir 4), tangga menurut OQ-CNP-01.
func aksiKomite(j *jalanAksi) error {
	h, n := j.h, j.r.Indeks
	b, err := adjBaris(h, n)
	if err != nil {
		return err
	}
	if !models.BolehKirimKomite(b) {
		return fmt.Errorf("%w: akseptasi sudah diserahkan ke komite", ErrAksiTertutup)
	}
	email, err := j.l.a.EmailPelaku(j.ctx, j.k.Pelaku)
	if err != nil {
		return err
	}
	salinan := h.Salin()
	kSalinan := *j.k
	c, err := models.ProteksiKirimKomite(&kSalinan, salinan, n, email)
	if err != nil {
		return err
	}
	if !c.Lolos(salinan) {
		if err := validasi(salinan); err != nil {
			return err
		}
		return fmt.Errorf("%w: penyerahan komite ditolak proteksi", ErrAksiTertutup)
	}
	if err := wajibTerisi(j, models.Evaluasi(h, models.LayarKomite(n), false)); err != nil {
		return err
	}
	ok, err := models.ValidasiGrossSpreadingIn(h, n) // 21-22 (sebelum 29)
	if err != nil {
		return err
	}
	if !ok {
		return &GalatValidasi{Pesan: []string{models.PesanGrossSpreadingIn}}
	}
	if len(h.AmbilDaftar(models.JalurAdj(n, models.AnakXOL))) == 0 { // 28 -> END (36): tanpa XOL tidak ada kasus komite
		return nil
	}
	roster, err := j.l.a.RosterKomite(j.ctx, models.HanyaTingkat1(h, n)) // 26
	if err != nil {
		return err
	}
	if len(roster) == 0 {
		return &GalatValidasi{Pesan: []string{"Roster komite (EMAILKOMITE STS_KLAIM NONPROP) kosong"}}
	}
	if err := j.l.g.SimpanHalaman(j.ctx, j.tx, j.kasus.ID, h); err != nil { // ID baris akseptasi stabil
		return err
	}
	var anggota []repository.AnggotaTangga
	for i, r := range roster { // roster terpilih (RosterKomite, urut DEGREE) = tangga kasus komite
		anggota = append(anggota, repository.AnggotaTangga{Urut: i + 1, OperatorID: r.OperatorID, Jabatan: r.Jabatan, Email: r.Email})
	}
	nama, err := j.l.a.NamaPelaku(j.ctx, j.k.Pelaku)
	if err != nil {
		return err
	}
	komiteLama := b[models.PropKomiteID]
	komite, err := j.l.g.BuatKasusKomite(j.ctx, j.tx, j.kasus.ID, b[models.PropID], j.k.Pelaku, nama, anggota, j.k.Sekarang) // 29
	if err != nil {
		return err
	}
	if err := j.l.g.SetelKomiteAdjustment(j.ctx, j.tx, b[models.PropID], komite, komiteLama); err != nil {
		return err
	}
	if err := models.TandaiKirimKomite(h, n, komite); err != nil { // 8, 20, 31-34
		return err
	}
	return j.antre(JenisEfekEmailKomite, komite, map[string]string{"klaim": j.kasus.ID, "komite": komite}) // 35 SendEmailKlaim
}

// aksiKasir - tombol "Acceptation" (`HitServiceToKasir_Act`). Langkah 3: hanya bila status konversi "1"
// (getStatusKonversi - dibaca HANYA di produksi); di luar produksi aktivitas keluar tanpa efek, persis XML.
func aksiKasir(j *jalanAksi) error {
	n := j.r.Indeks
	b, err := adjBaris(j.h, n)
	if err != nil {
		return err
	}
	if !(b["DirectToKasir"] == "true" && b["StatusKasir"] == "") { // 2
		return nil
	}
	sts, err := j.l.a.StatusKonversi(j.ctx, strings.ReplaceAll(b["AcceptedNo"], ".", "")) // 3
	if err != nil {
		return err
	}
	if sts != "1" {
		return nil
	}
	if !models.PanjangNoAksepCNP(b["AcceptedNo"]) { // 9 F=6
		return nil
	}
	email, err := j.l.a.EmailCeding(j.ctx, j.h.Ambil(models.TM+"CedingID")) // 9.2
	if err != nil {
		return err
	}
	// IsPEGASyariah (9.4) = node server Pega syariah; aplikasi ini satu instans konvensional (OQ-CNP-40).
	m, err := models.SusunMuatanKasir(j.h, n, email, j.k.Pelaku, j.l.kasir, false) // 9.3-9.6
	if err != nil {
		return err
	}
	if err := j.antre(JenisEfekKasir, b["AcceptedNo"], m); err != nil { // 9.7 (hanya produksi)
		return err
	}
	if b["IDOfBank"] == "" { // 11-12
		id, err := j.l.a.IDBankRekening(j.ctx, b["NameOfBank"], b["BranchOfBank"], b["NoAccount"])
		if err != nil {
			return err
		}
		b["IDOfBank"] = id
	}
	return nil
}

// aksiSimpanDibayar - tombol "Save Previously Paid" (`SaveCNPLayerList_Act`): satu baris OS STS 5 (CNPLayerList
// akseptasi dikurangi Previously Calculated).
func aksiSimpanDibayar(j *jalanAksi) error {
	b, err := models.BarisOSDibayar(j.h, j.kasus.ID, j.r.Indeks)
	if err != nil {
		return err
	}
	return j.l.g.SisipOS(j.ctx, j.tx, b, j.k.Sekarang)
}

// aksiTutupKlaim - tombol "Yes" local action CloseClaimMD (`CloseClaimTNonProp`). Remarks (`.Message`) wajib; disimpan
// ke REMARK_CLOSE (`ClaimData.Remark_Close`) supaya alasan penutupan tidak hilang (PARITAS).
func aksiTutupKlaim(j *jalanAksi) error {
	h := j.h
	if err := wajibTerisi(j, models.Evaluasi(h, models.LayarTutupKlaim(), false)); err != nil {
		return err
	}
	if !models.TutupKlaim(j.k, h) { // 1-4, 6
		return validasi(h)
	}
	h.Setel(models.CD+"Remark_Close", h.Ambil("Message"))
	k := j.k
	if err := j.l.g.SisipOS(j.ctx, j.tx, models.SusunBarisOS(h, j.kasus.ID, models.StsOSFinal, models.ParamOSTutup(h)),
		k.Sekarang); err != nil { // 5
		return err
	}
	if err := j.salinJSONKlaim(); err != nil { // 7
		return err
	}
	if err := j.antre(JenisEfekTutup, j.kasus.ID, map[string]string{"STS_REJECT": models.StsOSFinal}); err != nil { // 8
		return err
	}
	if err := models.HitungTurunan(h); err != nil {
		return err
	}
	if err := j.l.g.SimpanHalaman(j.ctx, j.tx, j.kasus.ID, h); err != nil {
		return err
	}
	j.selesai = true
	return j.l.g.TutupKasus(j.ctx, j.tx, j.kasus.ID, j.kasus.Tahap, k.Sekarang) // 9 ASMForceCaseClose
}

// aksiBukaCWP - tombol "Close Without Payment (CWP)": pra-proses `CloseClaimNP_preAct`, lalu local action CloseClaimNP.
// Tombol Yes-nya nonaktif (OQ-CNP-36).
func aksiBukaCWP(j *jalanAksi) error {
	nama, err := j.l.a.NamaPelaku(j.ctx, j.k.Pelaku)
	if err != nil {
		return err
	}
	models.PraCWP(j.k, j.h, nama)
	j.bukaModal = "cwp"
	return nil
}

// ---------------------------------------------------------------- efek keluar

// Jenis efek outbox Claim Non Prop (Connect-REST korpus; hanya produksi).
const (
	JenisEfekOutstanding = "outstanding-np" // SaveDataToOSAksep_Act 15.8.8 -> InsertClaimOutstanding_NP
	JenisEfekKonversi    = "konversi-klaim" // KonversiKlaim_Act -> KonversiKlaimNonLife (Klaim / insertClaimAccept)
	JenisEfekTutup       = "tutup-klaim-np" // CloseClaimTNonProp 8 -> insertClaimFinalOrClosed_NP
	JenisEfekKasir       = "kasir"          // HitServiceToKasir_Act 9.7 -> SendAcceptationToKasir (Kasir / insertAllPaymentKasir)
	JenisEfekEmailKomite = "email-komite"   // CreateChildKomiteCNP_Act 35 -> SendEmailKlaim
)

// MuatanOutbox - isi baris T_LOG_SERVICE_RNM.MUATAN (teks JSON milik outbox inti).
type MuatanOutbox struct {
	Kategori1 string `json:"kategori1,omitempty"`
	Kategori2 string `json:"kategori2,omitempty"`
	KlaimID   string `json:"klaimId"`
	Isi       any    `json:"isi,omitempty"`
}

// kunciEfek - kunci M_LINK_SERVICE per jenis (Kategori_1/Kategori_2 VERBATIM dari activity). Connector
// InsertClaimOutstanding_NP / insertClaimFinalOrClosed_NP ber-URL tertulis mati tanpa kategori M_LINK_SERVICE
// (OQ-CNP-31) - kuncinya kosong, pelaksana berhenti terang.
var kunciEfek = map[string][2]string{
	JenisEfekKonversi: {"Klaim", "insertClaimAccept"},
	JenisEfekKasir:    {"Kasir", "insertAllPaymentKasir"},
}

// ErrEfekTakDikenal - jenis efek tanpa penanganan.
var ErrEfekTakDikenal = errors.New("services: jenis efek outbox tidak dikenal")
