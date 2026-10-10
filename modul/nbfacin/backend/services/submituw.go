package services

// Submit tahap UNDERWRITING (tiket 61) - Email Type UW 1 Accept / 2 Reject / 3 Ask / 7 Decline (`DDL\EmailTypeUW.xml`) +
// Comment + Validity Of Acceptance. Keputusan work owner 08-10-2026 (diteruskan sesi 9d, D1-D7; rincian di docs/issues/61):
// cakupan "jenjang UW dulu" - Reject / Ask / Decline penuh, Accept hanya naik jenjang UW lewat limit; Accept ke jalur akhir
// (Decision24: Kadiv / Direktur / FacRetro / Binding / SaveToProduction) DITAHAN. Decline juga di JUW A (penyimpangan sadar:
// flow tanpa konektor Decline dari Assignment21).
//
// Bukti `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\` (semua langkah tanpa label `//` kecuali disebut):
//   - FlowAction\InwardFacultative_IsUW.xml: pre InputOfferFacInEngineerUW_preACT, post SetValidateDateUW_PostAct, section
//     EmailSection_IsUW (EmailTypeUW wajib bila StatusBusiness != 3 && FlagOnGoingPolicy != 1; Comment wajib; DaysValidity
//     -> GetDateValidity_ACT edit; DateValidity / WaitingBindDate baca-saja).
//   - Activity\GetDateValidity_ACT.xml langkah 2 (NB bawaan) / 4 (edit) -> validitasBawaan / validitasEdit.
//   - Activity\SetValidateDateUW_PostAct.xml langkah 5.3 (ViewSuggest), 8 (InsertHistoryAkseptasiPega), 9
//     (InsertFacInPerformance_Act: Accept ACCEPT + langkah 13 OUTSTANDING tahap berikutnya, Decline DECLINE), NB status Decline
//     L3474 (dua spasi). Langkah 3, 4, 10-14 di luar cakupan (D6).
//   - Activity\GetLimitAkseptasi_JUW_UW.xml + RDBList\GetLimitAkseptasi*JUWA_SQL.xml + Flow\InputInwardFacultativeOffer.xml
//     Decision42 -> acceptUW.
//   - Activity\CheckSpreadingProtect_ACT.xml (radio EmailTypeUW) -> proteksiSpreadingUW.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services/acceptance"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

const (
	// pesanFacRetro - Accept ke jalur yang belum dibangun (tiket 62, K62-1: FacRetro = tiket 65): teks SISTEM BARU (bukan
	// korpus). pesanTanggaTakPasti - tangga limit tidak pasti (galatTanggaElse): sistem baru menolak menebak (tiket 63;
	// sebelumnya ditahan dengan pesan "jalur akhir belum tersedia" tiket 62) `[keputusan agent]`.
	pesanFacRetro       = "Accept meneruskan case ke jalur FacRetro yang belum tersedia."
	pesanTanggaTakPasti = "Jenjang approval berikutnya tidak dapat ditentukan dari tabel limit; hubungi administrator."
	posisiSeniorUW      = "ReasFacInSeniorUnderwriting"
	posisiDepHeadUW     = "ReasFacInDepHeadUnderwriting"
	posisiJuniorUWA     = "ReasFacInJuniorUnderwritingA"
	posisiUWUtama       = "ReasFacInUnderwriting"
	posisiUWFinancial   = "ReasFacInUnderwritingFinancial"
	// Tiket 62 - jenjang lanjutan (Flow\InputInwardFacultativeOffer.xml Decision23: A1 / A8 / A11 / A6 / A4 / A5).
	posisiKadivTeknik   = "ReasFacInGroupLeader"
	posisiKadivFac      = "ReasFacInFacultativeDivHead"
	posisiKadivFin      = "ReasFacInFinDivHead"
	posisiDirMarketing  = "ReasFacInMarketingDirector"
	posisiDirTeknik     = "ReasFacInTechnicalDirector"
	posisiManagerTeknik = "ReasFacInManagerTeknik"
	// batasHariBawaan - GetDateValidity_ACT 2.5 (hanya mode bawaan NB). maksHariValiditas - kontrak api.ts (0-9999).
	batasHariBawaan   = 30
	maksHariValiditas = 9999
	// lebarCommentUW - T_GENERAL_POLIS.COMMENT_TEXT VARCHAR2(500) (rancangan, migrasi 212).
	lebarCommentUW = 500
	// bentukTanggalRiwayat - TO_CHAR(TGL_INP, 'DD-MM-YYYY HH24:MI') (sqlRiwayatAkseptasi), jam dinding WIB.
	bentukTanggalRiwayat = "02-01-2006 15:04"
)

// statusInboxPosisi - NB status (sesudah awalan "NB") tiap posisi tujuan Accept, SATU tempat. SET korpus untuk Senior / DepHead
// dan jenjang lanjutan memuat NAMA ORANG (tidak disalin): Senior / DepHead = keputusan work owner 08-10-2026 D7 (tiket 61),
// jenjang lanjutan = K62-2 (nama jabatan). JUW_A / UNDERWRITER verbatim (spasi sebelum 'S apa adanya).
var statusInboxPosisi = map[string]string{
	posisiSeniorUW:      " IS IN SENIOR UNDERWRITER'S INBOX",
	posisiDepHeadUW:     " IS IN DEPHEAD UNDERWRITER'S INBOX",
	posisiJuniorUWA:     " IS IN JUNIOR UNDERWRITER (A)'S INBOX",
	posisiUWUtama:       " IS IN UNDERWRITER 'S INBOX",
	posisiKadivTeknik:   " IS IN KADIV TEKNIK'S INBOX",
	posisiKadivFac:      " IS IN KADIV FACULTATIVE'S INBOX",
	posisiKadivFin:      " IS IN KADIV FINANCIAL'S INBOX",
	posisiDirMarketing:  " IS IN DIREKTUR MARKETING'S INBOX",
	posisiDirTeknik:     " IS IN DIREKTUR TEKNIK'S INBOX",
	posisiManagerTeknik: " IS IN MANAGER TEKNIK'S INBOX",
}

// tujuanAcceptUW - Flow Decision42 (LetterNo GetLimitAkseptasi_JUW_UW -> posisi). ToDepHeadUW mensyaratkan IsB2B != "ASM" -
// IsB2B tidak ada di sistem baru: non-ASM. Else -> Decision24.
var tujuanAcceptUW = map[string]string{"SENIORUW": posisiSeniorUW, "DEPHEADUNDERWRITER": posisiDepHeadUW, "JUW_A": posisiJuniorUWA,
	"UNDERWRITER": posisiUWUtama}

// jabatanPosisiUW - GetLimitAkseptasi_JUW_UW langkah 4-7: PositionNote -> CARI5 (JUW dibaca B, ADR-0032).
var jabatanPosisiUW = map[string]string{posisiJuniorUWB: "JUW_B", posisiJuniorUWA: "JUW_A", posisiUWUtama: "UNDERWRITER"}

// posisiDecision24 - Accept Senior (A7) / UW Financial (A10) / DepHead (A17) -> Decision24 (cek FacRetro) -> Utility2.
var posisiDecision24 = map[string]bool{posisiSeniorUW: true, posisiUWFinancial: true, posisiDepHeadUW: true}

// posisiJenjangLanjut - Assignment1 / 8 / 11 / 6 / 4 / 5: Accept LANGSUNG ke Utility2 GetLimitAkseptasi_ActFlow (tanpa
// Decision24) `[terverifikasi]`; Submit tahap uw meliputi posisi ini juga (tiket 62), Spreading / Akumulasi tidak.
var posisiJenjangLanjut = map[string]bool{posisiKadivTeknik: true, posisiKadivFac: true, posisiKadivFin: true, posisiDirMarketing: true,
	posisiDirTeknik: true, posisiManagerTeknik: true}

// jabatanPengguna - KEPUTUSAN-30-09-2026 butir 33 (limit pengguna = MASUKAN dari model peran, BUKAN `LOGIN =
// {OperatorID.pyUserIdentifier}`; dikonfirmasi tiket 62): jabatan pengguna = peran POSITION case saat ini, kebalikan konektor
// Decision42 / Decision23 `[terverifikasi]`. UW Financial tidak punya jabatan di tabel limit -> Else (batas pengetahuan).
var jabatanPengguna = map[string]acceptance.Jabatan{posisiJuniorUWB: "JUW_B", posisiJuniorUWA: "JUW_A", posisiUWUtama: "UNDERWRITER",
	posisiSeniorUW: "SENIORUW", posisiDepHeadUW: "DEPHEADUNDERWRITER", posisiKadivTeknik: "KADIVTEKNIK", posisiKadivFac: "KADIVFACULTATIVE",
	posisiKadivFin: "KADIVFINANCIAL", posisiDirMarketing: "DIREKTURMARKETING", posisiDirTeknik: "DIREKTURTEKNIK", posisiManagerTeknik: "MANAGERTEKNIK"}

// tabelPreferred - langkah 11-15 cabang IsFire: `.OfferFacIn.IsPreferredRisk` -> tabel limit (RDB *JUWA_SQL); nilai lain ->
// tanpa kueri -> LetterNo "" -> Else.
var tabelPreferred = map[string]string{"Preferred Risk": repository.TabelLimitProperty,
	"Non-Preferred Risk": "M_LIMIT_PROPERTY_NON_PREFERREDD", "Preferred Risk Commercial": "M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL"}

// ---- Validity Of Acceptance (GetDateValidity_ACT) ----------------------------------------------------------------
//
// `[dugaan]` diuji empiris (D2; perintah audit `docs/issues/61-audit/cekvaliditas.py` dan `cekhari.py`, dijalankan dengan
// `PYTHONIOENCODING=utf-8 py <skrip>` atas `D:\migrasi\RNM\DDL\CONTOH`, hanya mencetak hitungan): @DateTime.addCalendar(tanggal, tahun, bulan, ?, hari, jam, …) - arg1 tahun
// (SystemSetOneYear_DT), arg2 bulan (tiket 51 "iya ikut pega"), arg4 hari (DateValidity = tanggal DateSuggest + DaysValidity
// pada 69/69 ekspor DDL\CONTOH), arg5 jam (Set12Hour_DT). WaitingBind rumus L701 cocok 40/56 (ekspor mayoritas EDM) - batas
// pengetahuan. Selisih hari dipotong (2/2 contoh yang membedakan). Tanggal = 00:00 WIB, disimpan teks Pega DateTime UTC
// ("yyyyMMddT170000.000 GMT", seperti ekspor).

// tanggalWIB - tanggal kalender WIB `t` (00:00 WIB) = @FormatDateTime(t, "yyyyMMdd", "Asia/Jakarta").
func tanggalWIB(t time.Time) time.Time {
	y, m, d := t.In(WIB).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, WIB)
}

// teksTanggalPega - tanggal 00:00 WIB -> teks Pega DateTime UTC.
func teksTanggalPega(d time.Time) string { return d.UTC().Format(BentukWaktuPega) }

// hariBawaanNB - langkah 2.2-2.6: Begindate < sekarang (atau kosong) -> 0; selain itu selisih hari dipotong, +1 bila > 0,
// paling banyak 30.
func hariBawaanNB(sekarang time.Time, mulai string) int {
	b, err := time.Parse(BentukWaktuPega, mulai)
	if err != nil || b.Before(sekarang) {
		return 0
	}
	hari := int(b.Sub(sekarang) / (24 * time.Hour))
	if hari > 0 {
		hari++
	}
	if hari > batasHariBawaan {
		hari = batasHariBawaan
	}
	return hari
}

// isianValiditas - masukan rumus WaitingBindDate L701.
type isianValiditas struct {
	sekarang         time.Time
	mulai, keUW      string    // StartDateTime / DateToUW teks Pega DateTime
	aksepTerakhir    time.Time // TGL_INP riwayat terakhir SEBELUM Submit ini (jam dinding WIB)
	adaAksepTerakhir bool
	tba              bool // QuotationData.IsTBA - tidak ada di sistem baru: selalu false (D2)
}

// waitingBind - L701: @if(DateToUW <= Begindate (teks), @if(IsTBA, LastAksep + 3 bulan, Begindate + 1 bulan), LastAksep + 1
// bulan); "" bila tanggal dasarnya kosong / tak terurai (FormatDateTime("") = "").
func waitingBind(v isianValiditas) string {
	dasarAksep := func(n int) string {
		if !v.adaAksepTerakhir {
			return ""
		}
		return teksTanggalPega(tambahBulan(tanggalWIB(v.aksepTerakhir), n))
	}
	if v.keUW <= v.mulai {
		if v.tba {
			return dasarAksep(3)
		}
		b, err := time.Parse(BentukWaktuPega, v.mulai)
		if err != nil {
			return ""
		}
		return teksTanggalPega(tambahBulan(tanggalWIB(b), 1))
	}
	return dasarAksep(1)
}

// hasilValiditas - DaysValidity (teks bulat) dan dua tanggal teks Pega DateTime.
type hasilValiditas struct{ hari, tanggal, waitingBind string }

// hitungValiditas - mode bawaan (hariIsian nil, langkah 2) atau edit (langkah 4): DateValidity = hari ini (WIB) + hari.
func hitungValiditas(v isianValiditas, hariIsian *int) hasilValiditas {
	hari := 0
	if hariIsian != nil {
		hari = *hariIsian
	} else {
		hari = hariBawaanNB(v.sekarang, v.mulai)
	}
	return hasilValiditas{hari: strconv.Itoa(hari), tanggal: teksTanggalPega(tanggalWIB(v.sekarang).AddDate(0, 0, hari)),
		waitingBind: waitingBind(v)}
}

// uraiHariValiditas - "" -> nil (mode bawaan); bilangan bulat 0-9999 tanpa tanda; selain itu 400.
func uraiHariValiditas(teks string) (*int, error) {
	teks = strings.TrimSpace(teks)
	if teks == "" {
		return nil, nil
	}
	n, err := strconv.Atoi(teks)
	if err != nil || n < 0 || n > maksHariValiditas || strings.ContainsAny(teks, "+-") {
		return nil, fmt.Errorf("%w: daysValidity wajib bilangan bulat 0-%d", ErrMasukanSubmit, maksHariValiditas)
	}
	return &n, nil
}

// isianValiditasKasus - LastAksep = TGL_INP riwayat HISTORYAKSEPTASIPRODUCTION terbaru (sebelum Submit ini).
func (s *Service) isianValiditasKasus(ctx context.Context, id string, k models.KonteksSubmit) (isianValiditas, error) {
	v := isianValiditas{sekarang: s.sekarang(), mulai: k.StartDateTime, keUW: k.DateToUW}
	riwayat, err := s.submit.RiwayatAkseptasi(ctx, idPega(id))
	if err != nil {
		return v, err
	}
	if len(riwayat) > 0 {
		if t, err := time.ParseInLocation(bentukTanggalRiwayat, riwayat[0].Tanggal, WIB); err == nil {
			v.aksepTerakhir, v.adaAksepTerakhir = t, true
		}
	}
	return v, nil
}

// ValiditasSubmit - jawaban POST …/submit/validitas (hari teks bulat; tanggal kabel DD-MM-YYYY).
type ValiditasSubmit struct{ DaysValidity, DateValidity, WaitingBindDate string }

// HitungValiditasSubmit - POST …/submit/validitas, tidak menyimpan; penjaga = Submit tahap uw (401 / 503 / 404 / 403
// workbasket / 409 bukan tahap uw / 400 isian). daysValidity "" = mode bawaan (GetDateValidity_ACT langkah 2; Pega: klik radio
// `.EmailTypeUW` memanggilnya dengan edit="" - perubahan kontrak 08-10-2026, permintaan work owner lewat sesi 9d); angka =
// mode edit (langkah 4).
func (s *Service) HitungValiditasSubmit(ctx context.Context, pelaku inti.Pelaku, id, hariTeks string) (ValiditasSubmit, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return ValiditasSubmit{}, err
	}
	if err := s.siapSubmit(id); err != nil {
		return ValiditasSubmit{}, err
	}
	k, err := s.konteksSubmit(ctx, id)
	if err != nil {
		return ValiditasSubmit{}, err
	}
	if _, err := alasanTidakBoleh(k, pelaku); err != nil {
		return ValiditasSubmit{}, err
	}
	if !tahapUWForm(tahapKonteks(k)) { // tiket 69a: juga UW tahap Policy (EmailSection_IsUW tanpa syarat flag)
		return ValiditasSubmit{}, ErrKasusBacaSaja
	}
	hari, err := uraiHariValiditas(hariTeks)
	if err != nil {
		return ValiditasSubmit{}, err
	}
	v, err := s.isianValiditasKasus(ctx, id, k)
	if err != nil {
		return ValiditasSubmit{}, err
	}
	h := hitungValiditas(v, hari)
	return ValiditasSubmit{DaysValidity: h.hari, DateValidity: tanggalKabelValiditas(h.tanggal),
		WaitingBindDate: tanggalKabelValiditas(h.waitingBind)}, nil
}

// tanggalKabelValiditas - teks Pega DateTime -> DD-MM-YYYY WIB; "" tetap "".
func tanggalKabelValiditas(s string) string {
	if s == "" {
		return ""
	}
	return waktuKeKabel(s)
}

// validitasTampilan - GET …/submit tahap uw: isian TERSIMPAN saja ("" bila belum) - Pega mengisi field terikat properti itu
// baru saat radio Email Type diklik (mode bawaan, POST …/submit/validitas daysValidity ""); perubahan kontrak 08-10-2026.
func validitasTampilan(k models.KonteksSubmit) (hari, tanggal, wb string) {
	return k.DaysValidity, tanggalKabelValiditas(k.DateValidity), tanggalKabelValiditas(k.WaitingBindDate)
}

// ---- Accept: jenjang UW (GetLimitAkseptasi_JUW_UW + Decision42) ---------------------------------------------------

// tujuanUW - hasil Accept tahap uw: posisi + NB status, atau pesan (ditahan).
// tujuanUW - posisi + NB status tujuan Accept, atau pesan bila ditahan; binding = jalur akhir (tiket 63: JSON_OFFER +
// FACINOFFER + FLAG 2).
type tujuanUW struct {
	posisi, status, pesan string
	binding               bool
	// selesai - Accept berujung selesai Policy (tiket 69a-2, Utility1 SaveJsonPolicyFacIn_Act).
	selesai bool
}

// jabatanBerikutUW - GetLimitAkseptasiJUWA_SQL atas baris tabel limit (tanpa SQL baru, D3): JABATAN pertama (urut
// LIMIT_BOTTOM naik; seri -> urutan tabel) dengan team_group = TeamGroup, LIMIT_BOTTOM > bottom jabatan sendiri dan <
// CARID2. "" bila jabatan sendiri tidak ada / ganda dengan bottom berbeda (Pega ORA-01427) / tidak ada jenjang.
func jabatanBerikutUW(baris []models.BarisLimitA, tabel, teamGroup, jabatan string, carid2 *apd.Decimal) (string, string) {
	tg := strings.TrimSpace(teamGroup)
	if tg == "" { // `team_group = {TeamGroup}` dengan TeamGroup kosong tidak cocok baris mana pun (review tiket 61)
		return "", "TeamGroup kosong"
	}
	var sendiri *apd.Decimal
	type calon struct {
		jabatan string
		bottom  *apd.Decimal
	}
	var semua []calon
	for _, b := range baris {
		if b.Tabel != tabel || strings.TrimSpace(b.TeamGroup) != tg || b.LimitBottom == nil {
			continue
		}
		if strings.TrimSpace(b.Jabatan) == jabatan {
			if sendiri != nil && sendiri.Cmp(b.LimitBottom) != 0 {
				return "", "jabatan " + jabatan + " ganda dengan LIMIT_BOTTOM berbeda"
			}
			sendiri = b.LimitBottom
		}
		semua = append(semua, calon{strings.TrimSpace(b.Jabatan), b.LimitBottom})
	}
	if sendiri == nil {
		return "", "jabatan " + jabatan + " tidak ada di " + tabel + " team " + tg
	}
	sort.SliceStable(semua, func(i, j int) bool { return semua[i].bottom.Cmp(semua[j].bottom) < 0 })
	for _, c := range semua {
		if c.bottom.Cmp(sendiri) > 0 && c.bottom.Cmp(nolBila(carid2)) < 0 {
			return c.jabatan, ""
		}
	}
	return "", "tidak ada jenjang di antara jabatan sendiri dan TSI"
}

// preferredKasus - `.OfferFacIn.IsPreferredRisk` (SetFlagOccupation_ACT; derivasi yang sama dengan proteksiLokasi 9.2.27 /
// 9.2.28 - langkahPreferred).
func preferredKasus(objek []models.ObjekFire, kategori map[string]string) string {
	topRiskKasus := false
	for _, o := range objek {
		topRiskKasus = topRiskKasus || o.IsTopRisk
	}
	preferred, topRiskLokal := "", false
	for _, o := range objek {
		topRiskLokal = topRiskLokal || o.IsTopRisk
		for _, ok := range o.Occupations {
			preferred = langkahPreferred(preferred, kategori[ok.OccupationID], topRiskKasus == topRiskLokal)
		}
	}
	return preferred
}

// dasarJenjang - masukan jenjang Accept tahap uw: IsPreferredRisk, spreading tersimpan (top risk, FacRetro), dua total TSI,
// dan baris tabel limit (TeamGroup / Jabatan di-TRIM).
type dasarJenjang struct {
	preferred         string
	sp                models.KasusSpreading
	adaTop            bool
	totalNR, totalTop *apd.Decimal // totalNR nil bila ada top risk (tidak dipakai) atau kurs tidak ada
	pesanKurs         string
	baris             []models.BarisLimitA
	oldID             string // BUSINESS.OLDID Class of Business case; "" = tidak diketahui (oldIDBisnisKasus)
	// Tiket 73 P4: barisB = M_LIMIT_FINANCIALINS (K73-34 NextFinancial); aneka = case keluarga Aneka (nil = FIRE, perilaku tiket
	// 61 / 62 tidak berubah selain K73-33).
	barisB []models.BarisLimitB
	aneka  *jenjangAneka
}

// jenjangAneka - case keluarga Aneka: BUSINESS case (kasusTangga), lini, dan pohon objek Aneka (TSI / FacRetro).
type jenjangAneka struct {
	kb     kasusBisnis
	lini   liniSpreadingAneka
	lokasi []models.LokasiAneka
}

// barisTSIJenjang - satu baris TempCurrencyValue CountTotalTSIPremiNusaRe_Act: TSI (CARI22) dan TSI top risk (CARID1) dalam mata
// uang asli; idr = kurs 1 tanpa pembacaan kurs.
type barisTSIJenjang struct {
	mataUang  string
	idr       bool
	tsi, top  *apd.Decimal
	lokasiTop bool
}

// totalTSIUW - CountTotalTSIPremiNusaRe_Act (dipanggil langkah 1 GetLimitAkseptasi_JUW_UW / _ActFlow) langkah 13-14:
// TotalTSINusaRe = PS × Σ CARI22 × kurs / 100 (L8826 / L9206; TabaruFund 100 - syariah di luar cakupan) dan TotalTSITopRisk =
// PS × Σ kurs × CARID1 / 100 (L9060 / L9254 - DENGAN kurs; K73-33 keputusan work owner 09-10-2026, FIRE juga: semula tiket 61
// tanpa kurs mengikuti CountPremiAndTSIRNMFireMBU_ACT 6.4-6.5). Kurs TERBARU (GetCurrencyToIDR_SQL), IDR = 1. Baris: FIRE
// langkah 5 (CoverageList(1).TSILiability item; top risk = sama, lokasi top risk) atau Aneka / Golf langkah 6 / 7
// (barisTSIAneka). Pesan = kurs tidak ada (ditahan, teks agent).
func (s *Service) totalTSIUW(ctx context.Context, d *dasarJenjang) error {
	var baris []barisTSIJenjang
	if d.aneka != nil {
		baris = barisTSIAneka(d.aneka)
	} else {
		for _, l := range d.sp.Lokasi {
			for _, it := range l.Items {
				if len(it.Coverages) == 0 {
					continue
				}
				tsi := it.Coverages[0].TSILiability
				baris = append(baris, barisTSIJenjang{mataUang: it.Currency, idr: it.Currency == namaIDR, tsi: tsi, top: tsi, lokasiTop: l.IsTopRisk})
			}
		}
	}
	ps := nolBila(d.sp.PercentShare)
	cache, total, top := map[string]*apd.Decimal{}, nol(), nol()
	for _, b := range baris {
		kurs := satu
		if !b.idr {
			k, ada, err := kursKeIDR(ctx, s.spreading, cache, b.mataUang)
			if err != nil {
				return err
			}
			if !ada || k == nil {
				d.pesanKurs = "Kurs " + b.mataUang + " ke IDR tidak ada di TREATYEXCHANGEYEARLY - limit akseptasi tidak dapat dihitung"
				return nil
			}
			kurs = k
		}
		total = tambah(total, kali(nolBila(b.tsi), kurs))
		if b.lokasiTop {
			d.adaTop = true // IsAdaTopRisk "true"
			top = tambah(top, kali(nolBila(b.top), kurs))
		}
	}
	d.totalTop = persenDari(ps, top)
	d.totalNR = persenDari(ps, total)
	return nil
}

// barisTSIAneka - CountTotalTSIPremiNusaRe_Act langkah 6 Aneka (IsAneka || IsBondingAndCustomBonds, L3695 / L3718; objek okupasi:
// CARI22 = CoverageList(1).TSI L2968, kurs mata uang objek L2394, IDR objek = 1 L3257; top risk lokasi = CoverageList(1).
// TSILiability L3422) dan langkah 7 Golf (isGolfInsurance L4673; objek Golf, SETIAP coverage: CARI22 = .TSI L4325, kurs mata
// uang objek L3913, kurs 1 bila Currency coverage "IDR" L4538 - objek IDR juga 1 `[keputusan agent]`; tanpa top risk).
// FlagDelete tidak ada di sistem baru.
func barisTSIAneka(a *jenjangAneka) []barisTSIJenjang {
	var baris []barisTSIJenjang
	for _, l := range a.lokasi {
		if a.lini.aneka {
			for _, ok := range l.Okupasi {
				for _, o := range ok.Objek {
					if len(o.Coverages) == 0 {
						continue
					}
					c := o.Coverages[0]
					baris = append(baris, barisTSIJenjang{mataUang: o.Currency, idr: o.Currency == namaIDR, tsi: c.TSI, top: c.TSILiability,
						lokasiTop: l.IsTopRisk})
				}
			}
		}
		if a.lini.golf {
			for _, o := range l.ObjekGolf {
				for _, c := range o.Coverages {
					baris = append(baris, barisTSIJenjang{mataUang: o.Currency, idr: o.Currency == namaIDR || strings.TrimSpace(c.Currency) == namaIDR,
						tsi: c.TSI})
				}
			}
		}
	}
	return baris
}

// carid2 - DataSearch.CARID2: @Math.divide(TotalTSITopRisk bila ada top risk, selain itu TotalTSINusaRe, 1, 0).
func (d dasarJenjang) carid2() *apd.Decimal {
	if d.adaTop {
		return bagiSkala(d.totalTop, satu, 0)
	}
	return bagiSkala(nolBila(d.totalNR), satu, 0)
}

// carid2UW - CARID2 kasus (dipakai uji); pesan = kurs tidak ada.
func (s *Service) carid2UW(ctx context.Context, id string) (*apd.Decimal, string, error) {
	d := dasarJenjang{}
	var err error
	if d.sp, err = s.bacaSpreading(ctx, nil, id); err != nil {
		return nil, "", err
	}
	if err := s.totalTSIUW(ctx, &d); err != nil || d.pesanKurs != "" {
		return nil, d.pesanKurs, err
	}
	return d.carid2(), "", nil
}

// adaFacRetro - `.OfferFacIn.IsFacRetro = 1` `[terverifikasi]` Activity\SetDataFacOutFire_Act.xml langkah 5 (aktif): ada baris
// spreading ber-TreatyType "10015" (FACOUT, langkah 2.3.4 / 2.3.5); langkah 6 mengosongkan IsInputFacRetro (0) - input
// FacRetro belum ada di sistem baru, jadi Decision22 selalu ke jalur FacRetro.
func (d dasarJenjang) adaFacRetro() bool {
	if d.aneka != nil { // SetDataFacOutAnekaGolf_Act (L1032 / L1909); K73-36: Golf SEMUA lokasi (Pega: lokasi terakhir, L1655)
		ada := false
		setiapCoverageAneka(d.aneka.lokasi, d.aneka.lini.aneka, d.aneka.lini.golf, func(_ int, _ *models.ObjekAneka, _ int, c *models.CoverageAneka) {
			for _, b := range c.Spreading {
				ada = ada || strings.TrimSpace(b.TreatyType) == treatyFACOUT
			}
		})
		return ada
	}
	for _, l := range d.sp.Lokasi {
		for _, it := range l.Items {
			for _, c := range it.Coverages {
				for _, b := range c.Spreading {
					if strings.TrimSpace(b.TreatyType) == treatyFACOUT {
						return true
					}
				}
			}
		}
	}
	return false
}

// kasusTangga - rules.Kasus untuk acceptance.Next: properti clipboard yang dibaca GetLimitAkseptasi_ActFlow dan When-nya.
// BusinessType "Fire" = cakupan modul (NB FIRE); StatusBusiness "1" = NB; IsBanding / IsFlagReject / IsB2B tidak ada di sistem
// baru (kosong: langkah 18 selalu membuang KADIVFACULTATIVE, ToDepHeadUW non-ASM).
type kasusTangga map[string]string

func (k kasusTangga) Nilai(jalur string) (string, bool) {
	v, ada := k[jalur]
	return v, ada
}

func (d dasarJenjang) kasusTangga(k models.KonteksSubmit) kasusTangga {
	t := d.kasusTanggaFire(k)
	if d.aneka != nil { // tiket 73 P4: BusinessType / BusinessCode / BusinessName case (kasusBisnis) untuk When COB
		for _, jalur := range []string{"pyWorkPage.Quotation.BusinessType", "pyWorkPage.OfferFacIn.QuotationData.BusinessType",
			"pyWorkPage.Quotation.BusinessCode", "pyWorkPage.OfferFacIn.QuotationData.BusinessCode", "pyWorkPage.Quotation.BusinessName",
			"pyWorkPage.OfferFacIn.QuotationData.BusinessName"} {
			t[jalur], _ = d.aneka.kb.Nilai(jalur)
		}
	}
	return t
}

func (d dasarJenjang) kasusTanggaFire(k models.KonteksSubmit) kasusTangga {
	return kasusTangga{
		"pyWorkPage.OfferFacIn.TotalTSINusaRe":               utils.FormatDecimal(nolBila(d.totalNR)),
		"pyWorkPage.OfferFacIn.TotalTSITopRisk":              utils.FormatDecimal(nolBila(d.totalTop)),
		"pyWorkPage.IsAdaTopRisk":                            teksBoolPega(d.adaTop),
		"pyWorkPage.OfferFacIn.QuotationData.StatusBusiness": "1",
		"pyWorkPage.OfferFacIn.QuotationData.BusinessType":   "Fire",
		"pyWorkPage.OfferFacIn.IsPreferredRisk":              d.preferred,
		"pyWorkPage.OfferFacIn.QuotationData.TeamGroup":      strings.TrimSpace(k.TeamGroup),
		"pyWorkPage.PositionNote":                            k.Position,
		"pyWorkPage.OfferFacIn.QuotationData.BusinessOldId":  d.oldID,
	}
}

// oldIDBisnisKasus - BusinessOldId case (repository.sqlOldIDBisnis). Class of Business / Group Business kosong, 0 atau > 1
// baris BUSINESS, OLDID kosong -> "" + log: predikat bentuk B (ISLIMIT*, BusinessOldId lawan angka) lalu bergalat
// rules.ErrTafsirBerbeda -> Else (actFlow), bukan menebak "bukan bond".
func (s *Service) oldIDBisnisKasus(ctx context.Context, id string) (string, error) {
	if s.kasus == nil || s.tableOfLimit == nil {
		return "", nil
	}
	k, err := s.kasus.BacaKasus(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return "", ErrKasusTidakAda
	}
	if err != nil {
		return "", err
	}
	nama, grup := k.Opportunity.ClassOfBusiness, k.Opportunity.GroupBusinessID
	if strings.TrimSpace(nama) == "" || strings.TrimSpace(grup) == "" {
		log.Printf("nbfacin: kasus %s: Class of Business / Group Business kosong - BusinessOldId tidak diketahui", id)
		return "", nil
	}
	old, err := s.tableOfLimit.OldIDBisnis(ctx, nama, grup)
	if err != nil {
		return "", err
	}
	if len(old) != 1 || strings.TrimSpace(old[0]) == "" {
		log.Printf("nbfacin: kasus %s: BUSINESS.OLDID %q group %q = %d baris / kosong - BusinessOldId tidak diketahui", id, nama, grup, len(old))
		return "", nil
	}
	return strings.TrimSpace(old[0]), nil
}

// bacaDasarJenjang - lihat dasarJenjang.
func (s *Service) bacaDasarJenjang(ctx context.Context, id string) (dasarJenjang, error) {
	var d dasarJenjang
	if s.limit == nil {
		return d, ErrTanpaDatabase
	}
	objek, err := s.BacaObjek(ctx, id)
	if err != nil {
		return d, err
	}
	var kode []string
	for _, o := range objek {
		for _, ok := range o.Occupations {
			kode = append(kode, ok.OccupationID)
		}
	}
	kategori, err := s.submit.KategoriOkupasi(ctx, unik(kode))
	if err != nil {
		return d, err
	}
	d.preferred = preferredKasus(objek, kategori)
	if d.sp, err = s.bacaSpreading(ctx, nil, id); err != nil {
		return d, err
	}
	if d.aneka, err = s.bacaJenjangAneka(ctx, id); err != nil {
		return d, err
	}
	if err := s.totalTSIUW(ctx, &d); err != nil {
		return d, err
	}
	if d.oldID, err = s.oldIDBisnisKasus(ctx, id); err != nil {
		return d, err
	}
	baris, barisB, err := s.limit.MuatLimit(ctx)
	if err != nil {
		return d, err
	}
	d.barisB = barisB
	for _, b := range baris {
		b.TeamGroup, b.Jabatan = strings.TrimSpace(b.TeamGroup), strings.TrimSpace(b.Jabatan)
		d.baris = append(d.baris, b)
	}
	return d, nil
}

// bacaJenjangAneka - nil bila case bukan keluarga Aneka (lini tak terbaca / penyimpan objek Aneka tidak terpasang).
func (s *Service) bacaJenjangAneka(ctx context.Context, id string) (*jenjangAneka, error) {
	if s.objekAneka == nil {
		return nil, nil
	}
	kb, ada := s.kasusBisnisDari(ctx, id)
	if !ada {
		return nil, nil
	}
	l := liniSpreadingDari(id, kb)
	if !l.aneka && !l.golf {
		return nil, nil
	}
	lokasi, err := s.objekAneka.BacaObjekAneka(ctx, id)
	if err != nil {
		return nil, galatKasus(err)
	}
	return &jenjangAneka{kb: kb, lini: l, lokasi: lokasi}, nil
}

// tujuanPosisi - posisi tujuan + NB status (statusInboxPosisi).
func tujuanPosisi(posisi string) tujuanUW {
	return tujuanUW{posisi: posisi, status: awalanNBStatus() + statusInboxPosisi[posisi]}
}

// acceptUW - posisi tujuan Accept tahap uw, atau pesan bila ditahan. Flow\InputInwardFacultativeOffer.xml:
//   - JUW B / JUW A / UW: GetLimitAkseptasi_JUW_UW + Decision42; Else -> Decision24 (tiket 62; tiket 61 menahannya);
//   - Senior / UW Financial / DepHead: Decision24;
//   - Kadiv Teknik / Kadiv Facultative / Kadiv Financial / Direktur Marketing / Direktur Teknik / Manager Teknik: langsung
//     GetLimitAkseptasi_ActFlow (LOOP: Accept di tiap level menjalankan ulang tangga).
func (s *Service) acceptUW(ctx context.Context, id string, k models.KonteksSubmit, pelaku inti.Pelaku) (tujuanUW, error) {
	d, err := s.bacaDasarJenjang(ctx, id)
	if err != nil {
		return tujuanUW{}, err
	}
	switch {
	case jabatanPosisiUW[k.Position] != "":
		if d.pesanKurs != "" {
			return tujuanUW{pesan: d.pesanKurs}, nil
		}
		letter, alasan := s.letterJUWUW(k, pelaku, d, jabatanPosisiUW[k.Position])
		if posisi, ok := tujuanAcceptUW[letter]; ok {
			return tujuanPosisi(posisi), nil
		}
		log.Printf("nbfacin: Accept UW kasus %s: Decision42 LetterNo %q (%s) - Else -> Decision24", id, letter, alasan)
		return s.decision24(id, k, d, pelaku)
	case posisiDecision24[k.Position]:
		return s.decision24(id, k, d, pelaku)
	case posisiJenjangLanjut[k.Position]:
		return s.actFlow(id, k, d, pelaku)
	}
	return tujuanUW{pesan: pesanTanggaTakPasti}, nil // tidak terjadi: tahapSubmitUW hanya meloloskan posisi di atas
}

// letterJUWUW - LetterNo GetLimitAkseptasi_JUW_UW (Utility8): langkah 8 When IsGroup -> exit, LetterNo "" (Else; tiket 69a,
// diselaraskan ke tahap Offer - keputusan work owner 08-10-2026 "Selaraskan, ikut Pega"); selain itu GetLimitAkseptasiJUWA_SQL
// atas tabel IsPreferredRisk. alasan hanya untuk log.
func (s *Service) letterJUWUW(k models.KonteksSubmit, pelaku inti.Pelaku, d dasarJenjang, jabatan string) (string, string) {
	if s.grup.kasusGroup(k, pelaku) {
		return "", "IsGroup (GetLimitAkseptasi_JUW_UW langkah 8)"
	}
	tabel := tabelPreferred[d.preferred]
	if d.aneka != nil { // tiket 73 P4: langkah 14 IsEngineering (L2652) / 15 IsNonPropertyandNonEngineering (L2836), tanpa IsPreferredRisk
		tabel = string(acceptance.TabelNonPropEng)
		if d.aneka.lini.engineering {
			tabel = string(acceptance.TabelEngineering)
		}
	}
	if tabel == "" {
		return "", "IsPreferredRisk tidak dikenal"
	}
	return jabatanBerikutUW(d.baris, tabel, k.TeamGroup, jabatan, d.carid2())
}

// decision24 - `IsFacRetro` (Decision22 `IsInputFacRetro` selalu tidak di sistem baru) -> ditahan pesanFacRetro (tiket 65);
// selain itu Utility2 GetLimitAkseptasi_ActFlow.
func (s *Service) decision24(id string, k models.KonteksSubmit, d dasarJenjang, pelaku inti.Pelaku) (tujuanUW, error) {
	if d.adaFacRetro() {
		log.Printf("nbfacin: Accept UW kasus %s: IsFacRetro (spreading FACOUT) - jalur FacRetro ditahan", id)
		return tujuanUW{pesan: pesanFacRetro}, nil
	}
	return s.actFlow(id, k, d, pelaku)
}

// galatTanggaElse - galat acceptance.Next yang berarti "tidak ada tujuan pasti" -> Else (ditahan + log), bukan 500.
// rules.ErrTafsirBerbeda: BusinessOldId tidak diketahui (oldIDBisnisKasus) dibandingkan dengan literal angka predikat bentuk B.
var galatTanggaElse = []error{acceptance.ErrDaftarGanda, acceptance.ErrTabelTakTerpilih, acceptance.ErrUrutanTakPasti,
	acceptance.ErrLimitKosong, acceptance.ErrBarisPenggunaGanda, acceptance.ErrBentukB, rules.ErrTafsirBerbeda}

// actFlow - GetLimitAkseptasi_ActFlow + Decision23 lewat acceptance.Next (tiket 11): jabatan pengguna = jabatanPengguna
// (butir 33); langkah 3 "exit kalau group" = When IsGroup (grup.go, tiket 63 D5). Decision23 Else (Selesai) -> jalur akhir
// tujuanBinding (tiket 63); tangga tak pasti -> ditahan pesanTanggaTakPasti. Konektor ke Kadiv Facultative (A8) tidak
// men-SET PositionNote di Pega (cacat korpus) - di sini POSITION tetap ditulis (penyimpangan sadar, tiket 62).
func (s *Service) actFlow(id string, k models.KonteksSubmit, d dasarJenjang, pelaku inti.Pelaku) (tujuanUW, error) {
	jabatan, ok := jabatanPengguna[k.Position]
	if !ok {
		// Pega: login tanpa baris tabel limit -> SQL tanpa hasil -> LetterNo "" -> Decision23 Else (keputusan NB-11 "pengguna
		// tak ada di tabel = selesai seperti sistem lama").
		log.Printf("nbfacin: Accept UW kasus %s: posisi %s tanpa jabatan tabel limit (butir 33) - Decision23 Else", id, k.Position)
		return s.tujuanBinding(id, k, pelaku), nil
	}
	if d.pesanKurs != "" {
		return tujuanUW{pesan: d.pesanKurs}, nil
	}
	tangga, err := susunTangga(d.baris, d.barisB)
	if err != nil {
		return tujuanUW{}, err
	}
	pengguna := acceptance.Pengguna{Jabatan: jabatan, AnggotaGrup: s.grup.kasusGroup(k, pelaku)}
	t, err := acceptance.Next(d.kasusTangga(k), tangga.BentukA, pengguna)
	if errors.Is(err, acceptance.ErrBentukB) { // K73-34: Bond / Kredit Senior ke atas -> M_LIMIT_FINANCIALINS (kontrakfacin.Tangga.Langkah)
		t, err = acceptance.NextFinancial(d.kasusTangga(k), tangga.BentukB, pengguna)
	}
	for _, g := range galatTanggaElse {
		if errors.Is(err, g) {
			log.Printf("nbfacin: Accept UW kasus %s: tangga %v - ditahan", id, err)
			return tujuanUW{pesan: pesanTanggaTakPasti}, nil
		}
	}
	if err != nil {
		return tujuanUW{}, err
	}
	if t.Selesai {
		log.Printf("nbfacin: Accept UW kasus %s: Decision23 Else (LetterNo %q) - SAVE JSON_OFFER", id, t.JabatanTujuan)
		return s.tujuanBinding(id, k, pelaku), nil
	}
	return tujuanPosisi(string(t.Antrean)), nil
}

// ---- Proteksi Accept tahap uw (CheckSpreadingProtect_ACT, D4) ----------------------------------------------------

// proteksiSpreadingUW - pesan verbatim, urutan Pega: CheckSpreadingProtectFire_ACT 2.2.3 per coverage (indeks lokasi /
// item / coverage = .pxListSubscript, mulai 1): M1 baris share / treaty kosong, M2 Σ share ≠ 100, M3 Σ PremiumSpreaded ≠
// PremiNusantaraRe (keduanya @divide(x,1,4)); lalu cekSpreadingFactIn CopyAll: M4 "Total TSI Spreaded for …" (pesan ringkasan
// spreading yang ada) dan M5 "Invalid Treaty Type Spreading!" (korpus: teks kosong di jalur ini - teks K59-8, D4).
// CoverageBasis 5 (layering) di luar cakupan (K48). M6 ORS / M7 special acceptance (data tidak ada) dan M8 (EDM) dilewati.
// Tiket 73 P4: langkah 6 CheckSpreadingProtectAnekaGolf_ACT (pesanSpreadingAneka) sesudah cabang FIRE, dan M4 / M5 juga dari
// ringkasan spreading Aneka - untuk case keluarga Aneka saja (lini tak terbaca -> perilaku FIRE tidak berubah). Dipakai tahap uw
// / UW Policy dan tahap Binding (K73-29).
func (s *Service) proteksiSpreadingUW(ctx context.Context, id string) ([]string, error) {
	sp, err := s.bacaSpreading(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	pesan := pesanSpreadingFire(sp)
	var aneka *dataSpreadingAneka
	if s.objekAneka != nil {
		if kb, ada := s.kasusBisnisDari(ctx, id); ada {
			if l := liniSpreadingDari(id, kb); l.aneka || l.golf {
				d, err := s.bacaDataSpreadingAneka(ctx, id)
				if err != nil {
					return nil, err
				}
				pesan = append(pesan, pesanSpreadingAneka(d.lokasi, d.lini)...)
				aneka = &d
			}
		}
	}
	t, err := s.tampilan(ctx, id, sp, nil, true)
	if err != nil {
		return nil, err
	}
	ringkas := t.Pesan
	if aneka != nil {
		ta, err := s.tampilanSpreadingAneka(ctx, id, *aneka, nil, true)
		if err != nil {
			return nil, err
		}
		ringkas = append(append([]string{}, ringkas...), ta.Pesan...)
	}
	for _, p := range ringkas {
		if strings.HasPrefix(p, "Total TSI Spreaded for ") || p == pesanJenisTreaty {
			pesan = append(pesan, p)
		}
	}
	return pesan, nil
}

// pesanSpreadingFire - CheckSpreadingProtectFire_ACT 2.2.3 (M1-M3) atas pohon FIRE.
func pesanSpreadingFire(sp models.KasusSpreading) []string {
	pesan := []string{}
	for i, l := range sp.Lokasi {
		for j, it := range l.Items {
			for m, c := range it.Coverages {
				if strings.TrimSpace(c.CoverageBasis) == "5" {
					continue
				}
				letak := fmt.Sprintf("location %d object item %d coverage %d", i+1, j+1, m+1)
				total, premi := nol(), nol()
				for _, b := range c.Spreading {
					total, premi = tambah(total, b.SharePercentage), tambah(premi, b.PremiumSpreaded)
					if b.SharePercentage == nil || strings.TrimSpace(b.TreatyType) == "" {
						pesan = append(pesan, "% spreading on "+letak+" can't be empty!")
					}
				}
				if total.Cmp(seratus) != 0 {
					pesan = append(pesan, "% Total spreading on "+letak+" must be 100%!")
				}
				if bagiSkala(premi, satu, 4).Cmp(bagiSkala(nolBila(c.PremiNusantaraRe), satu, 4)) != 0 {
					pesan = append(pesan, "Total premium spreaded on "+letak+" is not equals premium RNM!")
				}
			}
		}
	}
	return pesan
}

// periksaAcceptUW - proteksi spreading lalu jenjang; pesan gabungan (kosong = lolos) dan tujuan.
func (s *Service) periksaAcceptUW(ctx context.Context, id string, k models.KonteksSubmit, pelaku inti.Pelaku) ([]string, tujuanUW, error) {
	pesan, err := s.proteksiSpreadingUW(ctx, id)
	if err != nil {
		return nil, tujuanUW{}, err
	}
	t, err := s.acceptUW(ctx, id, k, pelaku)
	if err != nil {
		return nil, tujuanUW{}, err
	}
	if t.pesan != "" {
		pesan = append(pesan, t.pesan)
	}
	return pesan, t, nil
}

// ---- tulisan Submit tahap uw ---------------------------------------------------------------------------------------

// lengkapiTulisanUW - Accept: posisi / NB status dari jenjang (Decision42); riwayat Pega Workbasket = PositionNote
// (InsertHistoryAkseptasiPega, SetValidateDateUW_PostAct 8); Validity Of Acceptance (mode bawaan bila daysValidity "",
// selain itu edit) + Comment (migrasi 212); FACINPERFORMANCE (langkah 9, D5).
func (s *Service) lengkapiTulisanUW(ctx context.Context, w *models.TulisanSubmit, id string, pelaku inti.Pelaku, isi models.IsiSubmit,
	k models.KonteksSubmit, tujuan tujuanUW) error {
	if isi.EmailType == models.EmailTypeAccept {
		w.PosisiBaru, w.NBStatus = tujuan.posisi, tujuan.status
		if tujuan.binding {
			b, err := s.tulisanBinding(ctx, id, k)
			if err != nil {
				return err
			}
			w.Binding = b
		}
	}
	w.Pega.Workbasket = k.Position
	hari, err := uraiHariValiditas(isi.DaysValidity)
	if err != nil {
		return err
	}
	v, err := s.isianValiditasKasus(ctx, id, k)
	if err != nil {
		return err
	}
	h := hitungValiditas(v, hari)
	w.Validitas = &models.ValiditasUW{Days: h.hari, Date: h.tanggal, WaitingBind: h.waitingBind, Comment: isi.Comment}
	w.Performance, w.PerformanceLanjut = performanceUW(isi.EmailType, pelaku, k, id)
	return nil
}

// performanceUW - InsertFacInPerformance_Act tahap uw (FlagOnGoingPolicy 0 -> POSITION_DOC "OFFER"): Accept -> ACCEPT (langkah
// 9) + langkah 13 OUTSTANDING tahap berikutnya "BINDING" ber-USERNAME pembuat case (pxCreateOperator; penggantian nama
// ter-hardcode langkah 13.4 TIDAK diport - penyimpangan, D5); Decline -> DECLINE (langkah 10); Reject / Ask -> tidak ada
// (CARI3 kosong). Baris dengan kolom wajib kosong dilewati (langkah 12 / 13.6).
func performanceUW(emailType string, pelaku inti.Pelaku, k models.KonteksSubmit, id string) (utama, lanjut *models.PerformanceFacIn) {
	dasar := models.PerformanceFacIn{PositionDoc: "OFFER", Username: pelaku.AkunID, InsuredName: k.InsuredName, SOB: k.SobName,
		MOLeader: k.MOLeader, COB: k.ClassOfBusiness}
	lengkap := func(f models.PerformanceFacIn) *models.PerformanceFacIn {
		if f.Username == "" || f.InsuredName == "" || f.SOB == "" || f.MOLeader == "" || f.COB == "" {
			log.Printf("nbfacin: Submit UW kasus %s: FACINPERFORMANCE %s/%s dilewati (kolom wajib kosong)", id, f.PositionDoc, f.StatusDoc)
			return nil
		}
		return &f
	}
	switch emailType {
	case models.EmailTypeAccept:
		a, b := dasar, dasar
		a.StatusDoc = "ACCEPT"
		b.PositionDoc, b.StatusDoc, b.Username = "BINDING", "OUTSTANDING", k.CreateOp
		return lengkap(a), lengkap(b)
	case models.EmailTypeDecline:
		d := dasar
		d.StatusDoc = "DECLINE"
		return lengkap(d), nil
	}
	return nil, nil
}
