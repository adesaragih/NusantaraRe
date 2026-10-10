package services

// Untuk apa berkas ini: AKSI LAYAR CHOOSE SURVEYOR / INPUT ADJUSTMENT - tambah / hapus adjustment, perhitungan panel
// InputAdjustment (Payment Type, gross, deductible, fee, salvage, mata uang, Ex Gratia, spreading), payable / rekening,
// Cedant Panel, Print DLA (ChooseDla_Act), Acceptation (SaveAcceptation + HitServiceToKasir_Act), Back.
//
// Konteks (`sasaran`):
//
//	"adjitem:..(o)..(i)"        Indeks = baris Adjustment (hapus) / 0 (tombol "+")
//	"adjdtl:..(o)..(i)..(a)"    panel adjustment a; Indeks = baris grid di dalamnya (spreading adjustment)
//	"cedant:..(o)..(i)..(a)"    pop-up ViewCedantPanel; Indeks = baris CedingCedantList polis
//	"dla:o"                     pop-up PrintDLA_dtl objek o

import (
	"fmt"
	"strconv"
	"strings"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// AwalanModalDLA - kunci modal pop-up Print DLA ("dla:o").
const AwalanModalDLA = "dla"

// hitungAdj - aksi hitung panel adjustment (o, i, e = adjustment): port lalu pesan validasi membatalkan aksi.
func hitungAdj(f func(k *models.Konteks, h *models.Halaman, o, i, a int) error) penanganAksi {
	return diBaris(func(j *jalanAksi) error {
		if err := f(j.k, j.h, j.o, j.i, j.e); err != nil {
			return err
		}
		return validasi(j.h)
	})
}

// murniAdj - aksi panel adjustment tanpa acuan.
func murniAdj(f func(h *models.Halaman, o, i, a int) error) penanganAksi {
	return hitungAdj(func(_ *models.Konteks, h *models.Halaman, o, i, a int) error { return f(h, o, i, a) })
}

// aksiTambahAdjustment - tombol "+" grid Adjustment (SetIndexObject_ACT + CountTotalEstimasi_Act).
func aksiTambahAdjustment(j *jalanAksi) error {
	if err := models.TambahAdjustment(j.k, j.h, j.o, j.i); err != nil {
		return err
	}
	return validasi(j.h)
}

// aksiSetPayable - change Payable To / Specify (SetPayable_Act): payable baris + klaim, rekening pertama klien payable
// (Payable 1 / 2) ke baris adjustment.
func aksiSetPayable(j *jalanAksi) error {
	r, err := models.RencanaPayable(j.h, j.o, j.i, j.e)
	if err != nil {
		return err
	}
	if r.SemuaKlien || r.Klien == "" { // 8.2 hanya daftar pilihan Name of Bank
		return nil
	}
	rek, err := j.l.a.RekeningBank(j.ctx, r.Klien, r.MataUang) // 6
	if err != nil {
		return err
	}
	return models.TerapkanRekening(j.h, j.o, j.i, j.e, rek) // 7
}

// aksiPilihRekening - autocomplete Name of Bank (Param = models.KunciRekening): baris dibaca ulang dari daftar yang
// sama dengan pilihan layar; SetPayableTo_act sesudahnya tanpa akibat tersimpan (1-4 ber-remark, 6 ReceiverClaim).
func aksiPilihRekening(j *jalanAksi) error {
	rek, err := j.l.rekeningAdj(j.ctx, j.h, j.o, j.i, j.e)
	if err != nil {
		return err
	}
	for _, r := range rek {
		if models.KunciRekening(r) == j.r.Param {
			return models.PilihRekening(j.h, j.o, j.i, j.e, r)
		}
	}
	return fmt.Errorf("%w: rekening %q tidak ada di pilihan", ErrPermintaanTidakSah, j.r.Param)
}

// aksiSetCedant - tombol Choose ("SetCedant:1") / Choose All ("SetCedant:2") pop-up ViewCedantPanel (SetCedant_act
// Param.ChooseType). Jenis dibaca dari nama aksi; Choose membaca CedingCo baris polis Indeks di server.
func aksiSetCedant(j *jalanAksi) error {
	if err := j.butuh(true, true, true); err != nil {
		return err
	}
	_, tipe, _ := strings.Cut(j.r.Aksi, ":")
	id := ""
	if tipe == "1" {
		rows := j.h.AmbilDaftar(models.AwalanPolis + ".CedingCedantList")
		if j.r.Indeks < 1 || j.r.Indeks > len(rows) {
			return fmt.Errorf("%w: baris cedant %d", ErrPermintaanTidakSah, j.r.Indeks)
		}
		id = rows[j.r.Indeks-1]["CedingCo"]
	}
	return models.SetCedant(j.h, j.o, j.i, j.e, tipe, id)
}

// aksiHapusSpreadAdj - tombol Delete baris spreading adjustment (Indeks = baris).
func aksiHapusSpreadAdj(j *jalanAksi) error {
	if err := models.HapusSpreadAdj(j.h, j.o, j.i, j.e, j.r.Indeks); err != nil {
		return err
	}
	return validasi(j.h)
}

// aksiBukaDLA - tombol Print DLA baris objek (CheckLimitSpreadingTreaty_Act + ProtectPrint + CreateRemarksNote_Act).
func aksiBukaDLA(j *jalanAksi) error {
	if err := j.butuh(true, false, false); err != nil {
		return err
	}
	if err := models.BukaDLA(j.h, j.o); err != nil {
		return err
	}
	j.bukaModal = fmt.Sprintf("%s:%d", AwalanModalDLA, j.o)
	return nil
}

// bolehDLA - tombol pembuka Print DLA terbuka (`VIS .IsFacretro = 1`, `NA .DLAStatus != 0`): Submit pop-up hanya sah
// selama pembukanya terbuka.
func bolehDLA(h *models.Halaman, o int) bool {
	ob, err := models.Objek(h, o)
	if err != nil {
		return false
	}
	return ob["IsFacretro"] == "1" && (ob["DLAStatus"] == "0" || ob["DLAStatus"] == "")
}

// aksiDLA = Submit PrintDLA_dtl (`ChooseDla_Act` Param.Index = o, lalu SendEmailDLA_ACT):
//
//	1-3    models.ChooseDLA (DLAStatus 1, IsTreatyOut, jenis nomor)
//	3.1.3  GenerateDLAFacin_Act - nomor "P", models.TerapkanDLAFac, InsertDLA_OS_SQL per AcceptedNo
//	3.1.4  DLAFacintoTreaty_Act - nomor "S" (gerbang CekLimit, OQ-CFI-24), models.TerapkanDLATreaty
//	6      KonversiKlaim_Act STS 1 (pre=false: selalu; efek hanya produksi)
//	7-8    log "AKSEPATSI"; 9-10 HitDLAClaimFacin (perbaikan prompt §5 butir 5: hanya produksi); 11-12 log "DLA"
//
// Perbaikan prompt §5 butir 4: transisi kosong langkah 9 yang melompat mundur ke JMP1 dibuang - 6-12 berjalan sekali;
// `DLAList` (RetroList.TotalClaim "123") tidak dibangun - nol pembaca di korpus (OQ-CFI-25). Surel (SendEmailDLA_ACT)
// = OQ-CFI-22, berkas PDF = OQ-CFI-20.
func aksiDLA(j *jalanAksi) error {
	if err := j.butuh(true, false, false); err != nil {
		return err
	}
	if !bolehDLA(j.h, j.o) {
		return fmt.Errorf("%w: Print DLA objek %d", ErrAksiTertutup, j.o)
	}
	if err := wajibTerisi(j, models.Evaluasi(j.h, models.LayarDLA(j.o), false)); err != nil {
		return err
	}
	h, k := j.h, j.k
	r, err := models.ChooseDLA(h, j.o)
	if err != nil {
		return err
	}
	if r.Keluar {
		return nil
	}
	nomor := ""
	for _, jenis := range r.Urutan {
		b, err := j.l.g.UrutNomor(j.ctx, j.tx, jenis, k.Sekarang) // GenerateDLAFacin_Act 7-10 / DLAFacintoTreaty_Act 5-8
		if err != nil {
			return err
		}
		nomor = models.RakitNomorKlaim(b.Jenis, h.Ambil(models.OQ+"BusinessOldId"), b.MMYYYY, b.Urut)
		if jenis == models.JenisNomorDLATreaty {
			if err := models.TerapkanDLATreaty(k, h, j.o, nomor); err != nil {
				return err
			}
			continue
		}
		hasil, err := models.TerapkanDLAFac(k, h, j.o, nomor)
		if err != nil {
			return err
		}
		for _, no := range hasil.Akseptasi { // 20.1.2 InsertDLA_OS_SQL
			if err := j.l.g.TandaiDLAOS(j.ctx, j.tx, no, nomor, k.Sekarang); err != nil {
				return err
			}
		}
	}
	if err := j.antreKonversi("1"); err != nil { // 6
		return err
	}
	kunci := models.KunciInstans(j.kasus.ID)
	if err := j.l.g.CatatLogLayanan(j.ctx, j.tx, repository.LogLayanan{IDPega: kunci, // 7-8
		Parameter: kunci + " / " + h.Ambil(models.JalurNoPolis) + " / 1", JenisService: models.JenisServiceAkseptasi,
		NoAkseptasi: nomor}, k.Sekarang); err != nil {
		return err
	}
	if err := j.antre(JenisEfekDLA, j.kasus.ID, map[string]string{"CASEID": kunci, "NO_DLA": nomor}); err != nil { // 9-10
		return err
	}
	if err := j.l.g.CatatLogLayanan(j.ctx, j.tx, repository.LogLayanan{IDPega: kunci, Parameter: kunci, // 11-12
		JenisService: models.JenisServiceDLA, NoDLA: nomor}, k.Sekarang); err != nil {
		return err
	}
	j.info = models.OQDokumenPDF
	return nil
}

// aksiAkseptasi = tombol "Acceptation" (`SaveAcceptation` lalu `HitServiceToKasir_Act`):
//
//	SaveAcceptation 1-7   models.SaveAcceptation (keluar bila sudah Print Acceptation)
//	8                     nama treaty spreading adjustment kosong (SetTreatNameAdjustment)
//	9.2                   kasir bila AcceptanceStatus 1, IsPrintAccept kosong, StatusKasir kosong
//	13                    models.SelesaiAkseptasi
//	15                    JSON_KLAIM (InsertJsonClaimNonMBU_act)
//	16                    KonversiKlaim_Act STS 1 (pre=false)
//
// Tombol kedua (HitServiceToKasir_Act) hanya berjalan ketika SaveAcceptation keluar di langkah 1 (klik ulang - jalur
// kirim ulang kasir). `[penyimpangan sadar]` (PARITAS): Pega menjalankan HitServiceToKasir_Act DUA kali pada klik
// pertama (9.2 lalu tombol kedua); di sini sekali per klik - efek kasir lewat outbox.
func aksiAkseptasi(j *jalanAksi) error {
	h, o, i, a := j.h, j.o, j.i, j.e
	ok, err := models.SaveAcceptation(j.k, h, o, i, a)
	if err != nil {
		return err
	}
	if !ok {
		return j.kasir(o, i, a)
	}
	for _, s := range h.AmbilDaftar(models.DaftarDiAdj(o, i, a, models.AnakAdjSpread)) { // 8
		if s["TreatyName"] != "" {
			continue
		}
		if s["TreatyName"], err = j.l.a.NamaJenisReas(j.ctx, s["TreatyType"]); err != nil {
			return err
		}
	}
	b, err := models.Adj(h, o, i, a)
	if err != nil {
		return err
	}
	if b["AcceptanceStatus"] == "1" && b["IsPrintAccept"] == "" && b["StatusKasir"] == "" { // 9.2
		if err := j.kasir(o, i, a); err != nil {
			return err
		}
	}
	if err := models.SelesaiAkseptasi(h, o, i, a); err != nil { // 13
		return err
	}
	if err := j.salinJSONKlaim(); err != nil { // 15
		return err
	}
	if err := j.antreKonversi("1"); err != nil { // 16
		return err
	}
	j.info = models.OQDokumenPDF
	return nil
}

// kasir = HitServiceToKasir_Act jalur `IsCLM` (13.2): muatan kasir ke outbox (hanya produksi). DIRECTTOKASIR_LOG ditulis
// bersama tanggapan kasir (13.5-13.6) - panggilan nyata menunggu persetujuan (OQ-CFI-26), pola Claim Prop.
func (j *jalanAksi) kasir(o, i, a int) error {
	b, err := models.Adj(j.h, o, i, a)
	if err != nil {
		return err
	}
	if !models.BolehKasir(b) { // 2
		return nil
	}
	sts, err := j.l.a.StatusKonversi(j.ctx, strings.ReplaceAll(b["AcceptedNo"], ".", "")) // 3 (hanya IsPEGAPROD)
	if err != nil {
		return err
	}
	if sts != "1" { // 3 transisi PASCA-langkah: status konversi 1 lanjut (T=2), selainnya keluar (F=6)
		return nil
	}
	if b["IDOfBank"] == "" { // 11-12
		id, err := j.l.a.IDBankRekening(j.ctx, b["NameOfBank"], b["BranchOfBank"], b["NoAccount"])
		if err != nil {
			return err
		}
		b["IDOfBank"] = id
	}
	if !models.PanjangNoAksepCLM(b["AcceptedNo"]) { // 13.2 F=6
		return nil
	}
	email, err := j.l.a.EmailCeding(j.ctx, models.KunciCedingKasir(j.h)) // 13.2.2.1.1.2
	if err != nil {
		return err
	}
	m, err := models.SusunMuatanKasir(j.k, j.h, o, b, email)
	if err != nil {
		return err
	}
	return j.antre(JenisEfekKasir, b["AcceptedNo"], m) // 13.4
}

// aksiKembaliEstimasi - tombol Back layar Input Adjustment (pre-DT `BackToEstimasi`, finishAssignment) -> Decision8
// IsBackStage -> Assignment7 Input Estimasi (worklist pembuat).
func aksiKembaliEstimasi(j *jalanAksi) error {
	models.BackToEstimasi(j.h)
	return j.pindah(models.TahapEstimasi, j.kasus.PembuatID)
}

// konteksAdj - pengurai konteks panel / modal adjustment ke (o, i, a).
func konteksAdj(k string) (o, i, a int, ok bool) {
	prefiks, idx, ok := models.UraiPanel(k)
	if !ok {
		return 0, 0, 0, false
	}
	switch prefiks {
	case models.PanelObjekAdj:
		return idx[0], 0, 0, true
	case models.PanelItemAdj:
		if len(idx) < 2 {
			return 0, 0, 0, false
		}
		return idx[0], idx[1], 0, true
	case models.PanelAdj, models.ModalCedant, models.ModalKomite:
		if len(idx) < 3 {
			return 0, 0, 0, false
		}
		return idx[0], idx[1], idx[2], true
	}
	return 0, 0, 0, false
}

// modalDLA - nomor objek kunci modal "dla:o".
func modalDLA(k string) (int, bool) {
	s, ok := strings.CutPrefix(k, AwalanModalDLA+":")
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	return n, err == nil && n > 0
}

// penanganAdjustment - aksi layar Input Adjustment (digabung ke `penangan` oleh init aksi.go).
func penanganAdjustment() map[string]penanganAksi {
	return map[string]penanganAksi{
		"CountTotalEstimasi":       diItem(aksiTambahAdjustment),
		"DisableSendComite":        hitungAdj(models.HapusAdjustment),
		"SetAdjTypePayment":        hitungAdj(models.SetAdjTypePayment),
		"SetGrossAdjustment":       hitungAdj(models.SetGrossAdjustment),
		"SetNilaiResikoSendiri":    hitungAdj(models.SetNilaiResikoSendiri),
		"SetValueAdjusterFee":      hitungAdj(models.SetValueAdjusterFee),
		"CheckCurrency":            hitungAdj(models.CheckCurrency),
		"CekExGratia":              hitungAdj(models.CekExGratia),
		"CountSpreadingAdjustment": murniAdj(models.CountSpreadingAdjustment),
		"TambahSpreadAdj":          murniAdj(models.TambahSpreadAdj),
		"HapusSpreadAdj":           diBaris(aksiHapusSpreadAdj),
		"SetDLACedingSOB":          murniAdj(models.SetDLACedingSOB),
		"SetPayable":               diBaris(aksiSetPayable),
		"SetPayableTo":             diBaris(func(*jalanAksi) error { return nil }), // 1-4 ber-remark; 6 ReceiverClaim
		"PilihRekening":            diBaris(aksiPilihRekening),
		"BukaCedant": diBaris(func(j *jalanAksi) error {
			j.bukaModal = models.KunciPanel(models.ModalCedant, models.DaftarAdj(j.o, j.i), j.e)
			return nil
		}),
		"SetCedant":        aksiSetCedant,
		"SimpanAdjustment": diBaris(func(*jalanAksi) error { return nil }), // click:save
		"Acceptation":      diBaris(aksiAkseptasi),
		"BukaDLA":          aksiBukaDLA,
		"ChooseDla":        aksiDLA,
		"BackToEstimasi":   aksiKembaliEstimasi,
		// komite TT2 (komite.go), penutupan / penolakan (tutup.go)
		"SendPICProtect":      diBaris(aksiBukaKomite),
		"KirimKomite":         diBaris(aksiKirimKomite),
		"PreventRejectClaim":  bukaModal(func(*jalanAksi) string { return ModalTutup }),
		"SetelAlokasiSalvage": halamanSaja(func(*jalanAksi) {}), // postValue + muat ulang pilihan konfirmasi
		"CloseClaim":          aksiTutupKlaim,
		"BukaRejectClaim":     aksiBukaTolak,
		// TT3 / TT4 (KCF-03)
		"SendRejectClaimToKomite2": aksiKirimTolakKomite,
		"SendCloseClaimToKomite":   aksiKirimTutupKomite,
	}
}
