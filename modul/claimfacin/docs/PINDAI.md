# Pindai Claim Fac In (klaim fakultatif inward)

> Pindai baca-saja 09–10-10-2026 untuk tahap 1 (prompt "IMPLEMENTASI CLAIM FAC IN, TAHAP 1 DARI 2"). Sumber: korpus
> `D:\XML\RNM_BRD\Claim Fac In\` (483 berkas, baca-saja), `D:\XML\RNM_BRD\Komite Claim FacIn\` (115 berkas, hanya data
> yang diserahkan ke kasus komite), pembanding `modul/claimprop` dan `modul/claimnonprop`, DEV POOLDATA lewat SELECT
> agregat / katalog saja. Status tiap tombol: [`PARITAS.md`](PARITAS.md); pertanyaan: [`OQ.md`](OQ.md).
>
> Penanda: **[inferensi]** = disimpulkan, bukan tertulis di XML. **[DEV]** = fakta DEV 09-10-2026. **REMARK** = langkah
> activity berlabel `//`, tidak pernah jalan. Aturan baca layar R1–R5 sama dengan
> `modul/claimprop/docs/grilling-ronde-2.md` §13.
>
> Akun orang, alamat surel, kata sandi, alamat IP / URL, dan nomor polis asli yang tertulis mati di XML **sengaja tidak
> dicetak** di dokumen ini.

## 1. Sensus korpus

| Jenis rule | Berkas |
| --- | ---: |
| Activity | 179 |
| RDBList | 63 |
| When | 60 |
| Section | 57 |
| FlowAction | 33 |
| DataTransform | 28 |
| ReportDefinition | 27 |
| Harness | 16 |
| DataPage | 8 |
| ConnectREST | 8 |
| Flow / DecisionTable / SystemSettings / `Struktur_Register_Flow.xlsx` | 1 / 1 / 1 / 1 |
| **Jumlah** | **483** |

Lima rule diekspor dari **checkout privat** developer (2026-09-07): `Flow/Register_Flow`, `Section/ClaimSurvey`,
`Section/InputInwardFacultativeDtl`, `Section/ViewHistoryClaim`, `Section/ViewPolis`, ditambah `When/IsSPK`. Dibangun
menurut ekspor itu; versi yang berjalan di Pega hidup belum dipastikan.

`Activity/InsertLogServiceClaim` gagal diurai alat pindai sesi ini; log layanan (`MONITORING_KLAIM_LOG`) ditulis menurut
parameter pemanggilnya, pola Claim Prop.

## 2. Alur

`Register_Flow` (kelas `ASM-FW-GCNMFW-Work-PNC`, awalan `CLM-`):

1. Start1 → Decision4 **B2B** (`When/IsSPK`: `OfferFacIn.IsB2B = "SPK"`) → langsung Assignment3; selain itu
2. Assignment1 **"Input Register"** (operator saat ini, FlowAction `InputRegister`: section `InputRegister`, pra
   `CallActivityInputRegister` + DT `InsertObjects_dt`, pasca DT `InputRegisterPostDT`, local action
   `ProteksiDataRegister_Act`, validate `ValidateDate` tidak diekspor) →
3. Assignment7 **"Input Estimasi"** (WorkList Custom = pembuat, FlowAction `InputEstimasi`: section
   `InputEstimasiAdmin`, pra `InputEstimationPre` + DT `SetEstimation_DT`) → Decision5 `IsBackStage` ? Assignment1 :
4. Assignment3 **"Choose Surveyor"** (WorkBasket Custom, FlowAction `InputSurveyor` berlabel "Input Adjustment": section
   `ClaimSurvey`, pra `SetTypePDFAdjustment` + DT `SetStartDateAdjustment`) → Decision8 `IsBackStage` ? Assignment7 :
5. END52 Resolved-Completed (`CloseClaim` → `ASMForceCaseClose`).

Workbasket XML `ReasPNCTeknik` tidak ada di DEV **[DEV]**; diganti `ReasKlaimTeknik` (bawaan b, sama dengan Claim Prop).

## 3. Layar

| Assignment | Section | Isi |
| --- | --- | --- |
| Input Register | `InputRegister` | tab Register (`InputRegisterDetail`), Policy Detail & Claims History (`ViewHistoryClaim` + `InputInwardFacultativeDtl`), Progress Claim (`ProgresClaim_Sec` / `SubProgresClaim_Sec`); Save / Submit / Reject Claim / View Status Payment Premi; pop-up `ViewPolis` (Choose Polis), `ProtectDOL`, `ClaimComiteeReject`, `Catastrope_Sec`, Cause of Loss |
| Input Estimasi | `InputEstimasiAdmin` | No Claim, Send to PIC Claim, Back; tab Estimation (`InputEstimasiDetail`: grid objek per lini → panel item `PropertyItemListGridEstimation*` / `ShowItemPA` → panel estimasi `Estimasi` / `EstimasiPA` / `EstimasiMarine`), Policy Detail, Progress, View Registration; pop-up `Pla_Dtl`, Outstanding Summary |
| Choose Surveyor | `ClaimSurvey` | tab Claim Details (`InputEstimasi`: Registration / Estimation / Adjustment & Acceptation `ShowObjectAdj` → `ItemListEstimation*` → `Adjusment_SC` → `InputAdjustment`), Policy Details and Claim History, Progress Claim; Back, Close Claim; pop-up `PrintDLA_dtl`, `ViewCedantPanel`, `ClaimComite`, `PreventRejectClaim` |

Grid bertingkat **objek → item objek → estimasi / adjustment** dengan varian per lini (60 rule `When`, satu fungsi per
When di `backend/models/lini.go`): Fire / Aneka / Golf, Marine Cargo, MBU, PA, Travel.

## 4. Polis dan objek

- Choose Polis membaca `FACINPRODUCTION` (±1,88 juta baris **[DEV]**, berindeks `NOPOLIS` / `IDPEGA` / `BUSINESSCODE`)
  lewat `GetPolisForClaim_SQL`; XML menyisipkan teks cari langsung (`{ASIS:InputData.CARI1}`) — diperbaiki menjadi bind.
- `CopyNB_Act` menyalin `JSON_POLIS.DATA_JSON` (nopolis + prodke) ke `OfferFacIn` (Java). `JSON_POLIS` RNM-F: 85 baris /
  84 polis **[DEV]**. Polis tanpa JSON tidak dapat disalin (OQ-CFI-11); NB Fac In sistem baru tidak menulis keduanya
  (OQ-CFI-30).
- Objek calon: `InsertObjects_dt` dari `OfferFacIn.LocationList` / `VehicleList` / `PersonList` per lini. Halaman polis
  tidak disimpan di klaim — dibaca ulang dari `JSON_POLIS` setiap muat; objek menyimpan kunci letaknya di polis.

## 5. Hitungan dan penomoran

- Estimasi: kurs, TSI Nusare, deductible (`CountTSI_Act`), `CheckEstimateValue` (TSI / Limit of Liability / nol /
  negatif), spreading polis → spreading klaim per treaty, `CheckLimit_Act1` (limit treaty / fac retro).
- Adjustment: Payment Type 1–7, gross, deductible tipe 1/2/3, VAT, salvage, adjuster fee, Ex Gratia, spreading
  adjustment + Break QS, Cedant Panel, Payable / rekening, `CheckLimitSpreadingTreaty_Act`.
- Nomor `KODE_PRODUKSI NONLIFE + huruf + BusinessOldId + "." + MM.YYYY + "." + urut 5 digit`
  (`PROC_GENERATE_SEQUENCE_NUMBER` dibaca dari `ALL_SOURCE`, ditulis ulang lewat `inti/backend/penomor`): klaim **K**
  (CFS pertama), PLA treaty **G** (+ revisi `/n`), DLA fac **P**, DLA treaty **S**.

## 6. Penyimpanan Pega dan fakta DEV (09-10-2026)

| Objek | Keadaan |
| --- | --- |
| `T_WORK_CLAIM` | baris PROP / NONPROP / Life; kolom ID, COVER_KEY, LINI, SENDTO_ADMIN, SENDTO_MEDICAL, CREATE_OP, CREATE_OP_NAME, TGL_UPDATE, TAHAP, TGL_CREATE, STATUS_WORK, POSITION. ID kasus = awalan + LPAD(`SEQ_WORK_CLAIM`, 6) |
| `OS_AKSEPTASI_KLAIM` | 8.585 kasus / 28.908 baris `CLM-` (`CASEID` = `ASM-FW-GCNMFW-WORK CLM-n`), 2018–2026; STS 0 / 1 / 2 / 4 = 16.893 / 8.776 / 532 / 2.700 baris |
| `JSON_KLAIM` | 61 kasus `CLM-` |
| Pega DEV | `Work-PNC`: 144 kasus `CLM` (134 New); `Work-Komite`: 27 kasus `KMT` |
| `EMAILKOMITE` FACIN | 5 baris aktif DEGREE 1–5 (Claim Supervisor / Claim Dept. Head / Technic Div. Head / Operational Director / Technical Director); batas 57.750.000 / 189.750.000 / 495.000.000 / 660.000.000; roster REOPEN 3 baris |
| `T_KATEGORI_DOC_KLAIM` | FAC 18 baris |
| `PROGRESSCLAIM` / `SUBPROGRESSCLAIM` | 6.649 / 9.016 baris. `PEGA_PROGRESSCLAIM`: sisip bila (IDPEGA, POSITION) belum ada, lalu `STATUS = 'Done'` untuk posisi lain; `PEGA_SUBPROGRESSCLAIM`: sisip bila belum ada (IDPEGA, IDPROGRES, JENISPROGRES, STATUS `Auto Create%`) |
| `CLAIMREJECTED` / `HISTORYAKSEPTASIPEGA` | 508 / 7.832 baris `CLM-` |
| Prosedur | `PEGA_PROGRESSCLAIM`, `PEGA_SUBPROGRESSCLAIM`, `PEGA_JSON_OS_AKSEP_KLAIM`, `PEGA_JSON_KLAIM_PNC`, `PROC_GENERATE_SEQUENCE_NUMBER`, `PEGA_D_CAUSE_OF_LOSS` VALID — dibaca, tidak dipanggil |

## 7. Integrasi (ConnectREST)

`KonversiKlaimNonLife` (STS 0 CFS / 1 akseptasi-DLA / 4 close), `HitDLAClaimFacin`, `SendAcceptationToKasir`,
`ServiceGoogle` (lewat `inti/backend/penyimpanan`) — keluar lewat outbox, hanya produksi. `getPremiumPaidOn`,
`getPremiumPaidOnMarine`, `getPaymentClaim`, `getPayAttachment` — baca-luar belum disetujui (OQ-CFI-14).

## 8. Rule dirujuk tetapi tidak diekspor

`ValidateDate` (OQ-CFI-13), harness `New` / `NewSample` (OQ-CFI-29), `ASM-FW-GCNMFW-Data-ObjectItem.GeneratePLA`
(OQ-CFI-21), harness `ClaimCommittee` dan local action `ShowRetro` kelas SpreadingRisk (OQ-CFI-25), aliran HTML dokumen
(OQ-CFI-20), prompt values `associated` (OQ-CFI-18), langkah Java "hapus spreading / treaty / value yang sama"
(`[inferensi]`, PARITAS §8).

## 9. Aturan baca

- `pyStepsBlockName` dicetak setiap membaca activity; `//` = REMARK, tidak dibangun.
- `pyStepsPreCondition = false` = prakondisi tidak dievaluasi, langkah selalu jalan.
- Kode WHEN: 1 lompat, 2 lanjut, 3 lewati, 5 kerjakan tanpa when berikut, 6 keluar activity.
- Langkah loop bermetode "-" yang membawa SET tidak menjalankan SET-nya; repeat EMBEDDED = for-each.
- Container ALWAYS tetap tampil walau kondisinya NEVER (R1–R5).
- Penulis properti dicari, bukan hanya pembaca; bukti = path + rule + nomor langkah.
