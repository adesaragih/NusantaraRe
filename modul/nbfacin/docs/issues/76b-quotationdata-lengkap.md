# 76b: T_QUOTATIONDATA lengkap — salinan nama / kode / status QuotationData (bertahap)

Status: tahap 1 selesai di backend (belum di-commit; migrasi 253 ditulis, belum dijalankan). Sisanya peta, menunggu work owner.

## Asal

Work owner 12-10-2026 (lewat sesi 9d): "Ya, bertahap". Sistem baru menyimpan bisnis / tertanggung sebagai RUJUKAN di
T_NB_OPPORTUNITY dan membaca nama dari master (T_M_ACCOUNT lewat ACCOUNT_ID, MARKETINGOFFICER lewat MOID); Pega menyimpan
SALINAN di QuotationData.

## Sensus `[terverifikasi]` — dua cara, cocok (instrumen yang sama dengan tiket 76, sudah diralat urutan pernyataan)

| | Cara 1: loader `tipeAkhirMigrasi` + `skemaTabel` | Cara 2: pengurai python berkas migrasi + blok `skema_gen.go` |
| --- | --- | --- |
| T_QUOTATIONDATA fisik | 16 (183 / 184 / 185 / 199) | 16 |
| Rancangan | 67 (nol amandemen) | 67 |
| Fisik bukan rancangan | 0 | 0 |
| Belum ada di DB | 51 | 51 |
| Tanpa 11 Endorsement / Renewal | 40 | 40 |

Sebelas yang tidak ditambah: EDM_CHARGE_FEE, EDM_DATE, EDM_NOTE, EDM_SOURCE, EDM_SOURCE_NOTE, EDM_STATUS, EDM_SURVEY, EDM_TYPE,
EDM_TYPE_NEW, ENDORSEMENT_INTERNAL_RETRO, RNW_DATE.

## Migrasi 253 (`backend/migrations/253_quotationdata_lengkap.sql` + `_down`)

Satu `ALTER TABLE {skema}.T_QUOTATIONDATA ADD (…)` 40 kolom, seluruhnya VARCHAR2 rancangan (nol NUMBER); `_down` = `DROP (…)`.
Rentang ketiga 250-254, T_QUOTATIONDATA lahir di 183 (urutan aman). Ditulis, tidak dijalankan. **Tidak wajib sebelum biner**:
penulisan tahap 1 menoleransi ORA-00904 (`jalankanKolomMigrasi`, dipakai bersama 252). Diperbarui: STRUKTUR (+40 baris),
`migrasi_sebagian_test` (berkas 253, hitungan 506 -> 546, `TestMigrasi253QuotationDataLengkap` dua arah + `_down`), MODUL.md
(nomor berikut 254).

## Kunci penelusuran `[terverifikasi]`

`NB FacIn\Activity\InputOfferFacInEngineer_preACT.xml` langkah 14 `Page-Copy pyWorkPage.Quotation -> pyWorkPage.OfferFacIn.QuotationData`
(pre-processing FlowAction InwardFacultative; varian UW langkah 14 sama): QuotationData sebagian besar adalah **salinan halaman
Quotation** tiap form dibuka. Pengisian awal Quotation (New Opportunity / SFA) **tidak ada di korpus** — format CONTOH
(`ASM-SFAGIS-WORK-ORG ORG-<n>`, `ASM-SFAGIS-WORK-CONTACT CON-<n>`) menunjuk modul SFA.

## Tahap 1 — ditulis Save General (`repository/kasus.go` `sqlSalinanQuotation`, `services/lini.go` `bisnisQuotation`)

Satu UPDATE T_QUOTATIONDATA per PARENT_ID dalam transaksi Save General, SESUDAH upsert T_QUOTATIONDATA (baris ada sejak Save
General pertama) dan SEBELUM daftar ceding. Titik setara `[keputusan agent]`: Pega menyalin ulang dari master tiap form dibuka;
Save General = saat rujukan (MOID, Policy Type) disimpan sistem baru. Nilai master lebih lebar dari kolom rancangan -> NULL (bukan
ORA-12899 yang menggagalkan Save General).

| Kolom | Sumber Pega `[terverifikasi]` | Sumber Go |
| --- | --- | --- |
| MARKETING_CODE / MARKETING_NAME / TEAM_GROUP / BRANCH_CODE / BRANCH_NAME | `Activity\CheckDataQuotation.xml` langkah 6: MARKETINGOFFICER per `ID = MOID` → CLIENTID / CLIENTNAME / TEAMGROUP / BRANCHDETAILID / BRANCHDETAILNAME (DDL `MARKETINGOFFICER.txt`) | subkueri MARKETINGOFFICER per `q.MOID` (tabel yang sudah dibaca Go: sqlKonteksSubmit, kotak masuk). ⚠️ Penyimpangan: langkah 5.2 Pega menjelajah ulang per ClientID bila MOSTATUS 2 (baris pertama tak tentu) — di sini baris ber-ID saja |
| INSURED_ID | Pega: GetInsuredID (POOLDATA.CLIENT per nama) / Change Insured Name | T_NB_OPPORTUNITY.INSURED_ID (= T_M_ACCOUNT.INSUREDID popup akun) `[keputusan agent]` — kesetaraan format dengan CLIENT.ID `[belum terverifikasi]` |
| INSURED_NAME | `SetInsuredName_Act` langkah 1 `param.Name` (Change Insured Name) | T_M_ACCOUNT.INSUREDNAME per ACCOUNT_ID (yang ditampilkan layar, sqlBacaKasus) `[keputusan agent]` |
| IS_MOP | `Activity\SetMOP_Act.xml` L1 PolicyType "0" → "NonMOP", L2 "1" / "2" → "MOP" (dipanggil Section Periode) | CASE atas POLICY_TYPE tersimpan; nilai lain → kolom tetap |
| BUSINESS_FAC | `InputQuotation_PreDT`: "F" untuk workgroup ReasFacIn (CONTOH: selalu F) | literal 'F' (modul Fac In) |
| STATUS_BUSINESS | tanpa setter di korpus (Quotation SFA); CONTOH 1 / 2 / 3 | literal '1' = `statusBusinessNB` (sudah dipakai FACINOFFER / produksi) |
| BUSINESS_CODE / BUSINESS_OLD_ID / GROUP_PANEL | BusinessOldId: `GetOldIDBusiness_SQL` (SaveJsonOfferFacIn_Act 5, GenerateNopolis_Act 7); BusinessCode / GroupPanel: Quotation SFA | `kasusBisnisDari` (BUSINESS NOTE + BUSINESSGROUPID, ber-GROUPPANEL, tepat satu); tak unik → kolom tetap (NVL) |
| BUSINESS_NAME | Quotation SFA (= BUSINESS.NOTE) | Class of Business opportunity |
| BUSINESS_TYPE / BUSINESS_TYPE2 | `SetBusinessType_Act` 2 (BusinessType_DeT, bawaan "UNKNOWN"); 3 / `InwardFacultative_PreDT` `@if(BusinessType = "Life","LIFE","NONLIFE")` | `businessTypeDari` (port DT) |

**Tidak diisi**: QQID (nol referensi korpus), TYPE (jalur EDM saja), PERIOD / PERIOD_MM (isian layar Life — ditunda),
OFFERING_DATE / INPUT_QUO_DATE (tanpa setter; dari SFA), GROUP_NAME (Pega dari POOLDATA.CLIENT.GROUPNAME per InsuredID —
CLIENT tidak dibaca Go), COB_GROUP (kolom sistem loader `kelompokLini`, tanpa padanan medan Pega).

Uji: `TestSalinanQuotation253` (teks SQL, 7 bind, urutan di SimpanGeneral), `TestSalinanBisnisQuotation` (FIRE / Life / BUSINESS
tak unik / lebar / tanpa pembaca), `TestBackfillSalinanQuotation253`, `TestMigrasi253QuotationDataLengkap`.

## Backfill — `docs/alat/isi_salinan_quotation_253.sql`

Sesudah 253; satu UPDATE tanpa COMMIT; cakupan LINI FAC NB-% dengan STATUS_BUSINESS masih NULL (penanda "salinan belum pernah
ditulis" — biner baru selalu menulis '1', jadi hasilnya tidak ditimpa). ⚠️ Kepala skrip memuat peringatan: **nilai master SEKARANG,
bukan saat case dibuat**. BUSINESS_TYPE / TYPE2 tidak diisi SQL (DT diport di Go) — terisi Save General berikutnya.

## Peta kolom sisanya — BELUM diimplementasi (untuk keputusan work owner)

Dua agen baca-saja (setter dihitung dua cara, cocok); klaim kunci dicek ulang sesi c3 (Page-Copy langkah 14, CheckDataQuotation 6,
SetMOP_Act). Path relatif `D:\migrasi\RNM\NB FacIn\`.

| Kolom | Setter | Kapan | Format (CONTOH) | Catatan |
| --- | --- | --- | --- | --- |
| BTN_QUOTATION | `Quotation.btnQuotation`: btnSOB_DT "SOB", btnCedingCO_DT "CedingCo", btnCedingCedant_DT "CedingCedant"; masuk lewat Page-Copy | tombol SOB / Ceding | penanda tombol | status UI, bukan data bisnis |
| COMMENT_TEXT | tidak ditemukan | — | teks | — |
| COPY_FROM | tidak ditemukan | — | "None" | — |
| EMAIL | ConcatAllEmail_DT (Quotation.Email = gabungan email QuotationList + CC) | EmailSection | daftar alamat dipisah koma | korespondensi Go (tiket 56) punya penerima sendiri |
| FAC_OUT_STATUS | OfferFacOut_PreAct 8 "OfferFacOut"; CheckProtectFacout_Act 7 "MemoPlacing"; SetDataFacOut* "" | Fac Out | teks | Fac Out belum diport |
| FRONTING_STATUS | nol di korpus | — | "No" | — |
| INPUT_MONTH | isian Periode (aktif bila ProRateType 4) | layar | angka teks | isian layar — ditunda |
| IS_EDM_SHARE | CountPaymentEdm_Act 14 (0 / 1) | EDM | 0 / 1 | jalur EDM |
| IS_GROUP | preACT 23-25 / UW preACT 10-11 / GetTeamGroup_Act 6: "Group" bila pembuat case ada di daftar login literal, selain itu "NonGroup" | buka form / Submit | Group / NonGroup | logika diport (`services/grup.go`, dari CREATE_OP) — **murah** bila ingin disimpan |
| IS_PROPOSAL | tidak ada di QuotationData (hanya halaman Vehicle) | — | 0 / 1 | — |
| IS_SURVEY_REPORT | isian radio Periode | layar | Yes / No | isian layar — ditunda |
| IS_TBA | tanpa setter (dibaca GetDateValidity_ACT) | — | "true" / "false" | Go konstanta false |
| IS_TEMPLATE | dropdown read-only Periode, tanpa setter | — | "NB" | — |
| OLD_POLICY_NO | tanpa setter (read-only Periode_IsUW) | EDM / renewal | nomor polis | NB: NULL |
| OPERATOR_ID | tanpa setter (dibaca LogActivity_Act) | — | login | — |
| SOB_LEADER0 | SetDataSobCeding_Act 3 (`BrowseSearchSobCeding_RD` Leader0) | pilih SOB | kode agen | SOB Go (tiket 33) membaca AGENT, Leader0 tidak |
| SOB_LSG | InputQuotation_PreDT `"001"` | Flow InputQuotation | "001" | konstanta |
| STATUS_SYARIAH | tanpa setter (dibaca CountPremiAndTSI*) | — | 0 / 1 | syariah di luar cakupan |
| PERIOD / PERIOD_MM | isian Periode (Life) + CheckPeriodLife_ACT 10 | layar | bulat / bulan | ditunda (isian layar) |
| OFFERING_DATE / INPUT_QUO_DATE | tanpa setter (SFA) | — | teks Pega DateTime | padanan T_GENERAL_POLIS.OFFERING_DATE / T_WORK_POLIS.TGL_CREATE `[dugaan]` |
| QQID / TYPE / GROUP_NAME / COB_GROUP | lihat "Tidak diisi" di atas | | | |

⚠️ Catatan: prakondisi IsGroup memuat literal login orang (tidak disalin ke artefak ini).
