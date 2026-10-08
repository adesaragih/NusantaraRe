package services

// Untuk apa berkas ini: PENANGAN AKSI YANG MENULIS LEBIH DARI HALAMAN - pilihan pop-up yang dibaca ulang di server,
// Save to issue RNM (SaveOutstanding_Act), Save (SetOutstanding_Act), Send to Acceptation (CheckNopolicy_Act +
// UpdateTableOS), Submit (finishAssignment OutstandingClaim + ProteksiData_act), penyerahan komite
// (AttachmentProtect_ACT + ProteksiInitialandDate_Act + AddKomiteTreatyChild_ACT), Acceptation (HitServiceToKasir_Act),
// PLA (TryMakePLA_Act), DLA (PrintDLATreatyIn), Close Claim (CloseClaimProp).

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/repository"
)

// pesanKosong - padanan validasi klien Pega (`pyClientValidation=true`, medan wajib kosong). `[dugaan]` teks bawaan
// Pega "Value cannot be blank"; rule pesannya tidak ada di ekspor.
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

// aksiPilihMaster - tombol Choose grid "Data Master TreatyIn" (`SetValueToClaim_Act`). Param = TREATYID|TREATYGROUPID|
// CLASSOFBUSINESSID; baris dibaca ULANG dari view (layar tidak menulis medan master).
func aksiPilihMaster(j *jalanAksi) error {
	p := strings.SplitN(j.r.Param, "|", 3)
	for len(p) < 3 {
		p = append(p, "")
	}
	b, ada, err := j.l.a.BarisMasterDari(j.ctx, p[0], p[1], p[2])
	if err != nil {
		return err
	}
	if !ada {
		return fmt.Errorf("%w: master treaty %q tidak ada", ErrPermintaanTidakSah, p[0])
	}
	return models.SetValueToClaim(j.k, j.h, models.ParamPilihMaster{COB: b.ClassOfBusiness, COBID: b.ClassOfBusinessID,
		TreatyContractName: b.TreatyContractName, ProportionType: b.ProportionType, IDMaster: b.TreatyID,
		TreatyGroupID: b.TreatyGroupID, TreatyGroupName: b.TreatyGroup})
}

// aksiPilihPolis - double-click grid "Data Polis" (`CheckNoPolicy`). Param = nomor polis; harus ada di daftar
// SetPolicyTreatyProp master ini.
func aksiPilihPolis(j *jalanAksi) error {
	rows, err := j.l.a.DaftarPolis(j.ctx, models.AwalanMaster(j.h.Ambil(models.CD+"IDMaster")),
		j.h.Ambil(models.CD+"TreatyGroupName"))
	if err != nil {
		return err
	}
	for _, r := range rows {
		if r.PolicyNo == j.r.Param {
			if err := models.CheckNoPolicy(j.k, j.h, r.PolicyNo, r.Quarter, r.TreatyYear, r.Prodke); err != nil {
				return err
			}
			// spreading atas dari polis, bawah dari SpreadingList master (keputusan work owner 08-10-2026)
			m, err := master(j)
			if err != nil {
				return err
			}
			return models.IsiSpreadingPolis(j.k, j.h, r.PolicyNo, m)
		}
	}
	return fmt.Errorf("%w: polis %q bukan polis master ini", ErrPermintaanTidakSah, j.r.Param)
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

// aksiSimpanKatastrofe - tombol Save form katastrofe baru (`SaveCatasrtope_Act` lalu `InputCatastrope(Cancel)`). Param =
// `ParamCat.NOTE`.
func aksiSimpanKatastrofe(j *jalanAksi) error {
	if strings.TrimSpace(j.r.Param) == "" {
		return &GalatValidasi{Pesan: []string{"Note: " + pesanKosong}}
	}
	k := models.SusunKatastrofe(j.k, j.h, j.r.Param)
	if err := j.l.g.SisipKatastrofe(j.ctx, j.tx, k, j.k.Sekarang); err != nil {
		return err
	}
	models.InputCatastrope(j.h, "Cancel")
	return nil
}

// aksiShareRNM - change dropdown RNM Share (`DisableEditRNMShare`): nilai harus salah satu pilihan GetRNMShareTreaty.
func aksiShareRNM(j *jalanAksi) error {
	pilih := j.r.Param
	if pilih == "" {
		pilih = j.h.Ambil("TreatyShareTemp.CARI1")
	}
	m, err := master(j)
	if err != nil {
		return err
	}
	for _, s := range models.PilihanShareRNM(j.h, m) {
		if s == pilih {
			models.DisableEditRNMShare(j.h, pilih)
			return nil
		}
	}
	return fmt.Errorf("%w: RNM Share %q bukan pilihan master", ErrPermintaanTidakSah, pilih)
}

// ---------------------------------------------------------------- Outstanding Claim

// aksiSetOutstanding - tombol "Save" (`SetOutstanding_Act` lalu save). Langkah 4 (SaveOSKlaimTreaty_SQL ->
// PEGA_JSON_OS_AKSEP_KLAIMTRT) = OQ-CP-08 (procedure tidak ada di DEV); langkah 5 InsertJsonClaimTreaty_act.
func aksiSetOutstanding(j *jalanAksi) error {
	if !models.KlaimSementaraBerpolis(j.h) {
		return nil
	}
	return j.l.g.SalinJSONKlaim(j.ctx, j.tx, models.KunciInstans(j.kasus.ID), j.h.Ambil(models.CD+"NoClaim"),
		j.h.Ambil(models.CD+"PolicyData.PolicyNo"), j.k.Sekarang)
}

// aksiSaveOutstanding - tombol "Save to issue RNM" (`SaveOutstanding_Act`, Section OutstandingClaim dan
// InputAcceptation_Est).
func aksiSaveOutstanding(j *jalanAksi) error {
	h, k := j.h, j.k
	models.ProteksiData(h) // 1
	// ⚠️ PENYIMPANGAN SADAR AC 30 (keputusan work owner): total share spreading setiap mata uang tepat 100 % - Pega
	// hanya memberi tanda sisi > 100 (CountSpreading_act 6.2.4).
	tepat, err := models.PeriksaShareTepat100(h)
	if err != nil {
		return err
	}
	for _, p := range tepat {
		h.TambahPesan(models.DaftarSpreading, p)
	}
	if err := models.PeriksaDOLSerupa(k, h); err != nil { // AC 102: duplikat DOL memblokir penyimpanan
		return err
	}
	oldID, err := j.l.a.KodeLamaBisnis(j.ctx, h.Ambil(models.OQ+"BusinessCode")) // 2
	if err != nil {
		return err
	}
	models.PeriksaOldID(h, oldID)       // 3-5
	if err := validasi(h); err != nil { // 6
		return err
	}
	rencana := models.RencanakanNomor(h)
	if rencana.Sementara { // 7-10
		tahun := strconv.Itoa(k.Sekarang.In(models.Jakarta).Year())
		no, err := j.l.g.NomorSementara(j.ctx, j.tx, tahun)
		if err != nil {
			return err
		}
		models.TerapkanNomorSementara(h, no)
	}
	if rencana.Klaim { // 15-20
		b, err := j.l.g.UrutNomor(j.ctx, j.tx, models.JenisNomorKlaim, k.Sekarang)
		if err != nil {
			return err
		}
		models.TerapkanNomorKlaim(h, models.RakitNomorKlaim(b.Jenis, h.Ambil(models.OQ+"BusinessOldId"), b.MMYYYY, b.Urut))
	}
	kunci := models.KunciInstans(j.kasus.ID)
	for _, os := range models.KirimEstimasi(k, h, kunci) { // 23
		if err := j.l.g.SisipOS(j.ctx, j.tx, os, k.Sekarang); err != nil {
			return err
		}
	}
	models.BekukanDataLama(h)                           // 24-28
	models.SelesaiOutstanding(k, h)                     // 29-30
	if err := models.IsiRetroDanPLA(k, h); err != nil { // 31-36
		return err
	}
	if err := j.l.g.SalinJSONKlaim(j.ctx, j.tx, kunci, h.Ambil(models.CD+"NoClaim"), // 37
		h.Ambil(models.CD+"PolicyData.PolicyNo"), k.Sekarang); err != nil {
		return err
	}
	if err := j.antreKonversi(models.StsOSEstimasi); err != nil { // 39 (prakondisi nonaktif - selalu)
		return err
	}
	if err := j.l.g.CatatLogLayanan(j.ctx, j.tx, repository.LogLayanan{IDPega: kunci, // 40-41
		Parameter: kunci + " / " + h.Ambil(models.CD+"PolicyData.PolicyNo"), JenisService: JenisServiceEstimasi},
		k.Sekarang); err != nil {
		return err
	}
	if h.Ambil(models.CD+"IsPLA") == "1" { // local action PrintFile [hanya jika IsPLA = 1]
		j.info = models.PesanPrintPla
	}
	return nil
}

// JenisServiceEstimasi - `Param.JenisService = "ESTIMASI"` (SaveOutstanding_Act langkah 40).
const JenisServiceEstimasi = "ESTIMASI"

// aksiPLA - local action GeneratePLA, post-activity `TryMakePLA_Act`: nomor PLA (kode NONLIFE + "H") bila kosong, nomor
// PLA ke estimasi terkirim. Berkas PDF = OQ-CP-05.
func aksiPLA(j *jalanAksi) error {
	h := j.h
	if h.Ambil(models.CD+"NoPla") == "" {
		oldID, err := j.l.a.KodeLamaBisnis(j.ctx, h.Ambil(models.OQ+"BusinessCode"))
		if err != nil {
			return err
		}
		b, err := j.l.g.UrutNomor(j.ctx, j.tx, models.JenisNomorPLA, j.k.Sekarang)
		if err != nil {
			return err
		}
		models.TerapkanNomorPLA(h, models.RakitNomorPLA(b.Jenis, oldID, b.MMYYYY, b.Urut))
	} else {
		models.TerapkanNomorPLA(h, "")
	}
	j.info = models.OQDokumenPDF
	return nil
}

// aksiKirimAkseptasi - tombol "Send to Acceptation" (`CheckNopolicy_Act` -> `UpdateTableOS`).
func aksiKirimAkseptasi(j *jalanAksi) error {
	h := j.h
	if !models.PeriksaPolisAkseptasi(h) { // 2-3
		return validasi(h)
	}
	if h.Ambil(models.CD+"NoClaim") == "" { // UpdateTableOS 1-5 (KODE_BIS = BRANCH_CODE - tidak pernah terisi)
		no, err := j.l.g.NomorCFS(j.ctx, j.tx, "", j.k.Sekarang)
		if err != nil {
			return err
		}
		models.TerapkanNomorKlaim(h, no)
	}
	if models.KlaimSementaraBerpolis(h) { // 7-8 (KLAIMTRT = OQ-CP-08)
		if err := j.l.g.SalinJSONKlaim(j.ctx, j.tx, models.KunciInstans(j.kasus.ID), h.Ambil(models.CD+"NoClaim"),
			h.Ambil(models.CD+"PolicyData.PolicyNo"), j.k.Sekarang); err != nil {
			return err
		}
	}
	models.SelesaiKirimAkseptasi(j.k, h) // UpdateTableOS 9 + CheckNopolicy_Act 5
	return nil
}

// aksiKirimKeTeknik - tombol "Send to Acceptation" (perintah work owner 08-10-2026 "send to acceptation nya langsung
// kirim ke teknik, ga usah klik submit lagi"): CheckNopolicy_Act lalu, bila IsAcceptation tersetel, Submit
// (finishAssignment -> Input Acceptation, workbasket ReasKlaimTeknik) dalam SATU transaksi aksi - validasi Submit yang
// gagal membatalkan seluruhnya (IsAcceptation, salinan JSON_KLAIM). Nomor klaim biasanya sudah diambil Save to issue
// RNM (tombol ini aktif hanya bila IsOutstanding = 1). Tombol Submit tetap untuk berkas yang sudah ber-IsAcceptation
// tetapi belum pindah tahap (dikirim sebelum perubahan ini / data lama).
func aksiKirimKeTeknik(j *jalanAksi) error {
	if err := aksiKirimAkseptasi(j); err != nil {
		return err
	}
	if j.h.Ambil("IsAcceptation") != "1" {
		return nil // CheckNopolicy_Act berhenti (polis belum ada) - tidak ada yang dikirim
	}
	return aksiSubmit(j)
}

// aksiSubmit - tombol "Submit" (finishAssignment FlowAction OutstandingClaim): [hanya jika Policy No dan NoClaim
// terisi], validasi klien (wajib), post-activity `ProteksiData_act`, lalu Transition4 -> Assignment1 Input Acceptation
// (workbasket ReasKlaimTeknik).
func aksiSubmit(j *jalanAksi) error {
	if !models.BolehSubmitOutstanding(j.h) {
		return fmt.Errorf("%w: Submit hanya bila Policy No dan Claim No terisi", ErrAksiTertutup)
	}
	ts := models.Evaluasi(j.h, models.LayarOutstanding(), false)
	if err := wajibTerisi(j, ts); err != nil {
		return err
	}
	models.ProteksiData(j.h)
	if err := validasi(j.h); err != nil {
		return err
	}
	return j.l.g.PindahTahap(j.ctx, j.tx, j.kasus.ID, models.TahapOutstanding, models.TahapAcceptation,
		models.WorkbasketAcceptation, j.k.Sekarang)
}

// ---------------------------------------------------------------- akseptasi

// aksiMataUangAdjustment - change Currency AdjustmentDetail (`SetNameCurrency_Act(CurrID=.CurrencyID)`).
func aksiMataUangAdjustment(j *jalanAksi) error {
	cur := models.AmbilJalur(j.h, models.JalurAdj(j.r.Indeks, "CurrencyID"))
	return models.SetNameCurrency(j.k, j.h, j.r.Indeks, cur)
}

// aksiPilihRekening - pilihan autocomplete "Name of Bank". Param = ACCOUNTNO|NAMEOFBANK; rekening dibaca ulang dari
// daftar pilihan baris itu.
func aksiPilihRekening(j *jalanAksi) error {
	p := strings.SplitN(j.r.Param, "|", 2)
	for len(p) < 2 {
		p = append(p, "")
	}
	daftar, err := models.PilihanRekening(j.k, j.h, j.r.Indeks)
	if err != nil {
		return err
	}
	for _, r := range daftar {
		if r.AccountNo == p[0] && r.NameOfBank == p[1] {
			return models.PilihRekening(j.h, j.r.Indeks, r)
		}
	}
	return fmt.Errorf("%w: rekening bukan pilihan baris ini", ErrPermintaanTidakSah)
}

// LampiranKlaim - cacah berkas per kategori (`AttachCategory.pxResults.CountAttach`). ⚠️ Sumber halaman itu tidak ada
// di ekspor dan unggahan berkas Claim Prop menunggu persetujuan storage (OQ-CP-12): selalu kosong - penyerahan komite
// tertahan pesan "Please Upload Attachment ...", persis perilaku XML tanpa lampiran.
func LampiranKlaim() map[string]int { return map[string]int{} }

// j0BukaKomite - `AttachmentProtect_ACT` + `CekPremiLunas_Act` (penanda Protect.CARI1 / CARI2) baris n.
func j0BukaKomite(l *Layanan, ctx context.Context, k *models.Konteks, h *models.Halaman, n int) error {
	lim, _, err := l.a.LimitDirekturUtama(ctx)
	if err != nil {
		return err
	}
	lolos, err := models.AttachmentProtect(h, n, LampiranKlaim(), lim)
	if err != nil {
		return err
	}
	premi, err := models.PremiLunas(k, h, n)
	if err != nil {
		return err
	}
	h.Setel("Protect.CARI1", nolSatu(lolos))
	h.Setel("Protect.CARI2", nolSatu(premi))
	return nil
}

func nolSatu(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// aksiBukaKomite - tombol "Send to Committe": AttachmentProtect_ACT (refresh), lalu - hanya bila Protect.CARI1 = 1 dan
// CARI2 = 1 - harness CommitteeTreaty dengan activity `ProteksiInitialandDate_Act`.
func aksiBukaKomite(j *jalanAksi) error {
	if err := j0BukaKomite(j.l, j.ctx, j.k, j.h, j.r.Indeks); err != nil {
		return err
	}
	if j.h.Ambil("Protect.CARI1") != "1" || j.h.Ambil("Protect.CARI2") != "1" {
		return nil
	}
	j.h.Setel("TempCommiteClaim.DateOfComitee", j.k.Hari())
	j.h.Setel("TreatyExchangeYearly.UserName", j.k.Pelaku)
	return models.ProteksiInitialandDate(j.k, j.h, j.r.Indeks)
}

// aksiKomite - tombol "Send Claim to Committee" (`AddKomiteTreatyChild_ACT`). Gerbang DITEGAKKAN di layanan (AC 57,
// 58): Protect.CARI1/CARI2 dihitung ulang, tombol harus tampil dan aktif pada isian sesudah digabung (Payable,
// Remarks, Circumstanses / Salvage / Adjuster Fee), baris belum diserahkan.
func aksiKomite(j *jalanAksi) error {
	h, n := j.h, j.r.Indeks
	if h.Ambil("Protect.CARI1") != "1" || h.Ambil("Protect.CARI2") != "1" {
		if err := validasi(h); err != nil {
			return err
		}
		return fmt.Errorf("%w: penyerahan komite ditolak (lampiran / premi)", ErrAksiTertutup)
	}
	ts := models.Evaluasi(h, models.LayarKomite(n, true), false)
	if !models.AksiTerbuka(ts, "AddKomiteTreatyChild", n) {
		if err := wajibTerisi(j, ts); err != nil {
			return err
		}
		return fmt.Errorf("%w: tombol Send Claim to Committee tidak tampil untuk isian ini", ErrAksiTertutup)
	}
	b, err := adjBaris(h, n)
	if err != nil {
		return err
	}
	if b["IsKomite"] == "1" || b[models.PropKomiteID] != "" {
		return fmt.Errorf("%w: baris adjustment sudah diserahkan ke komite", ErrAksiTertutup)
	}
	if err := models.TandaiKirimKomite(j.k, h, n); err != nil { // 1-15, 31
		return err
	}
	roster, err := j.l.a.RosterKomite(j.ctx, models.NilaiRosterKomite(b), models.STSKlaimProp) // 18-22
	if err != nil {
		return err
	}
	if len(roster) == 0 {
		return &GalatValidasi{Pesan: []string{"Roster komite (EMAILKOMITE STS_KLAIM PROP) kosong untuk nilai adjustment ini"}}
	}
	if err := j.l.g.SimpanHalaman(j.ctx, j.tx, j.kasus.ID, h); err != nil { // ID baris adjustment stabil
		return err
	}
	var anggota []repository.AnggotaTangga
	for i, r := range roster {
		anggota = append(anggota, repository.AnggotaTangga{Urut: i + 1, OperatorID: r.OperatorID, Jabatan: r.Jabatan, Email: r.Email})
	}
	nama, err := j.l.a.NamaPelaku(j.ctx, j.k.Pelaku)
	if err != nil {
		return err
	}
	komite, err := j.l.g.BuatKasusKomite(j.ctx, j.tx, j.kasus.ID, b[models.PropID], j.k.Pelaku, nama, anggota, j.k.Sekarang) // 25
	if err != nil {
		return err
	}
	if err := j.l.g.SetelKomiteAdjustment(j.ctx, j.tx, b[models.PropID], komite); err != nil {
		return err
	}
	b[models.PropKomiteID] = komite
	return j.antre(JenisEfekEmailKomite, komite, map[string]string{"klaim": j.kasus.ID, "komite": komite}) // 34 SendEmailKlaim
}

func adjBaris(h *models.Halaman, n int) (models.Baris, error) {
	d := h.AmbilDaftar(models.DaftarAdjustment)
	if n < 1 || n > len(d) {
		return nil, models.ErrBarisTidakAda
	}
	return d[n-1], nil
}

// aksiKasir - tombol "Acceptation" (`HitServiceToKasir_Act`; SaveAcceptationTreaty_Act di tombol yang sama
// ber-[hanya jika: When NEVER] - tidak pernah jalan). Langkah 3: hanya bila status konversi "1" (getStatusKonversi -
// dibaca HANYA di produksi); di luar produksi aktivitas keluar tanpa efek, persis XML.
func aksiKasir(j *jalanAksi) error {
	b, err := adjBaris(j.h, j.r.Indeks)
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
	if b["IDOfBank"] == "" { // 11-12
		id, err := j.l.a.IDBankRekening(j.ctx, b["NameOfBank"], b["BranchOfBank"], b["NoAccount"])
		if err != nil {
			return err
		}
		b["IDOfBank"] = id
	}
	if !models.PanjangNoAksepCLMP(b["AcceptedNo"]) { // 13.1 F=6
		return nil
	}
	email, err := j.l.a.EmailCeding(j.ctx, j.h.Ambil(models.TM+"CedingID")) // 13.1.2
	if err != nil {
		return err
	}
	m := models.SusunMuatanKasir(j.k, j.h, b, email, false)             // 13.1.3-13.3 (IsPEGASyariah = OQ-CP-14)
	if err := j.antre(JenisEfekKasir, b["AcceptedNo"], m); err != nil { // 13.4 (hanya produksi)
		return err
	}
	if b["IsFacRetro"] == "1" && b["DLA_No"] == "" { // local action PrintFileDLA
		j.info = models.PesanCetakDLA
	}
	return nil
}

// aksiDLA - local action GenerateDLATreaty, post-activity `PrintDLATreatyIn`: riwayat "Print DLA", nomor DLA (kode
// NONLIFE + "P") bila kosong. Berkas PDF = OQ-CP-05.
func aksiDLA(j *jalanAksi) error {
	n := j.r.Indeks
	if err := wajibTerisi(j, models.Evaluasi(j.h, models.LayarDLA(n), false)); err != nil {
		return err
	}
	b, err := adjBaris(j.h, n)
	if err != nil {
		return err
	}
	j.k.Riwayat(j.h, models.TeksCetakDLA) // 2-3
	if b["DLA_No"] == "" {                // 9-13
		oldID, err := j.l.a.KodeLamaBisnis(j.ctx, j.h.Ambil(models.OQ+"BusinessCode"))
		if err != nil {
			return err
		}
		bn, err := j.l.g.UrutNomor(j.ctx, j.tx, models.JenisNomorDLA, j.k.Sekarang)
		if err != nil {
			return err
		}
		b["DLA_No"] = models.RakitNomorPLA(bn.Jenis, oldID, bn.MMYYYY, bn.Urut)
	}
	j.info = models.OQDokumenPDF
	return nil
}

// aksiTutupKlaim - tombol "Yes" PreventRejectClaimProp (tanpa centang Close Without Payment) -> `CloseClaimProp`.
func aksiTutupKlaim(j *jalanAksi) error {
	h := j.h
	if err := wajibTerisi(j, models.Evaluasi(h, models.LayarTutupKlaim(), false)); err != nil {
		return err
	}
	if h.Ambil("TempCommiteClaim.AllocationShareSalvage") == "true" {
		return fmt.Errorf("%w: %s", ErrTertunda, models.OQTutupTanpaBayar)
	}
	if !models.PeriksaTutupKlaim(h) { // 1-4
		return validasi(h)
	}
	k := j.k
	models.SelesaiTutupKlaim(k, h, h.Ambil("TempCommiteClaim.Remarks")) // 5, 8
	kunci := models.KunciInstans(j.kasus.ID)
	if err := j.l.g.SisipOS(j.ctx, j.tx, models.BarisOSTutup(k, h, kunci), k.Sekarang); err != nil { // 6
		return err
	}
	if err := j.antreKonversi(models.StsOSTutupBerkas); err != nil { // 7
		return err
	}
	if err := j.l.g.SalinJSONKlaim(j.ctx, j.tx, kunci, h.Ambil(models.CD+"NoClaim"), // 9
		h.Ambil(models.CD+"PolicyData.PolicyNo"), k.Sekarang); err != nil {
		return err
	}
	if err := models.HitungTurunan(h); err != nil {
		return err
	}
	if err := j.l.g.SimpanHalaman(j.ctx, j.tx, j.kasus.ID, h); err != nil {
		return err
	}
	j.selesai = true
	return j.l.g.TutupKasus(j.ctx, j.tx, j.kasus.ID, j.kasus.Tahap, k.Sekarang) // 10 ASMForceCaseClose
}

// ---------------------------------------------------------------- efek keluar

// Jenis efek outbox Claim Prop.
const (
	JenisEfekKonversi    = "konversi-klaim" // KonversiKlaim_Act -> Connect-REST KonversiKlaimNonLife (Klaim / insertClaimAccept)
	JenisEfekKasir       = "kasir"          // HitServiceToKasir_Act -> SendAcceptationToKasir (Kasir / insertAllPaymentKasir)
	JenisEfekEmailKomite = "email-komite"   // AddKomiteTreatyChild_ACT 34 -> SendEmailKlaim
)

// MuatanOutbox - isi baris T_LOG_SERVICE_RNM.MUATAN (teks JSON milik outbox inti, bukan data klaim).
type MuatanOutbox struct {
	Kategori1 string `json:"kategori1,omitempty"`
	Kategori2 string `json:"kategori2,omitempty"`
	KlaimID   string `json:"klaimId"`
	Isi       any    `json:"isi,omitempty"`
}

// kunciEfek - kunci M_LINK_SERVICE per jenis (Kategori_1/Kategori_2 VERBATIM dari activity).
var kunciEfek = map[string][2]string{
	JenisEfekKonversi: {"Klaim", "insertClaimAccept"},
	JenisEfekKasir:    {"Kasir", "insertAllPaymentKasir"},
}

// antre - efek keluar hanya di produksi (`When/IsPEGAPROD` di setiap Connect-REST / email).
func (j *jalanAksi) antre(jenis, rujukan string, isi any) error {
	if !j.l.produksi {
		return nil
	}
	k := kunciEfek[jenis]
	teks, err := json.Marshal(MuatanOutbox{Kategori1: k[0], Kategori2: k[1], KlaimID: j.kasus.ID, Isi: isi})
	if err != nil {
		return err
	}
	_, err = j.l.g.AntreEfek(j.ctx, j.tx, jenis, rujukan, string(teks), j.k.Sekarang)
	return err
}

// antreKonversi = KonversiKlaim_Act (INSKEY, NOPOLIS, STS_REJECT).
func (j *jalanAksi) antreKonversi(sts string) error {
	return j.antre(JenisEfekKonversi, j.kasus.ID, map[string]string{"CASEID": models.KunciInstans(j.kasus.ID),
		"NOPOLIS": j.h.Ambil(models.CD + "PolicyData.PolicyNo"), "STS_REJECT": sts})
}
