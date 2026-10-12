# 76: T_GENERAL_POLIS lengkap — kolom rancangan yang belum ada di DB (bertahap, mulai identitas polis)

Status: tahap 1 + keputusan ronde 2 (total kasus, DATE_OFFER, GENERATE_DATE_POLICY_NO) selesai di backend (belum di-commit; migrasi 252 ditulis, belum dijalankan). Kolom isian layar ditunda; tahap 3 lain menunggu work owner. Lanjutan QuotationData: issues/76b.

## Asal

Work owner 12-10-2026 (lewat sesi 9d): "ini udh jadi polis masa data yg disimpan tidak lengkap" — contoh NB-151029, data
PolicyData tidak ada di T_GENERAL_POLIS. Keputusan: **BERTAHAP, mulai identitas polis**.

## Sensus `[terverifikasi]` — dua cara, cocok

| | Cara 1: fungsi loader `tipeAkhirMigrasi` + `skemaTabel` (uji sementara) | Cara 2: pengurai python berkas migrasi + blok `skema_gen.go` |
| --- | --- | --- |
| T_GENERAL_POLIS fisik | 25 | 25 |
| Rancangan | 78 (75 + amandemen NB_STATUS / DATE_TO_UW / SUBMIT_TO_UW) | 75 mentah + 3 amandemen |
| Fisik bukan rancangan | 0 | 0 |
| Rancangan belum ada di DB | 53 | 53 |
| Tanpa 5 kolom Endorsement | 48 | 48 |

⚠️ Ralat instrumen cara 2: putaran pertama memberi fisik 24 / belum 54 karena memproses semua ADD lalu semua DROP per berkas.
Migrasi 202 justru membuang PPN_CHECK lalu MENAMBAHNYA ULANG, jadi urutannya menentukan. Sesudah diurutkan menurut posisi
pernyataan, hasilnya cocok dengan cara 1. Kolom Endorsement yang tidak ditambah: ENDORSEMENT_NO, ENDORSMENT_REASON,
PRORATE_START_EDM, PRORATE_EDM_END, ERRMSG.

## Migrasi 252 (`backend/migrations/252_general_polis_lengkap.sql` + `_down`)

- Satu `ALTER TABLE {skema}.T_GENERAL_POLIS ADD (…)` 48 kolom, satu blok `/`; `_down` = `DROP (…)` 48 kolom. Semua NULL-able.
- Rentang ketiga nbfacin 250-254; T_GENERAL_POLIS lahir di 182, jadi urutan nama aman. Ditulis, **tidak dijalankan** (DBA).
- Nama / tipe / lebar = rancangan. NUMBER polos rancangan ditolak penjaga (`TestNolNumberTanpaPresisi`), jadi presisinya
  `[keputusan agent]` mengikuti konvensi modul: uang / persen / share NUMBER(38,8) (ADR-0016) untuk CEDING_RETENTION,
  CEDING_RETENTION_NOMINAL, SHARE_RNML, SUM_TOTAL_TSI, TOTAL_PREMI_NUSA_RE, TOTAL_TSI_NUSA_RE, TOTAL_TSI_NUSA_RE_SPREADING,
  TOTAL_TSI_TOP_RISK, TREATY_CAPACITY; identitas OLD_POLIS_ID NUMBER(19) (A92); kode PRO_RATE_TYPE NUMBER(5) (A170).
- **Tidak wajib sebelum biner**: penulisan tahap 1 menoleransi ORA-00904 (log, lanjut — pola `inti/backend/login` galatHak),
  sehingga Selesai Policy yang sudah berhasil tidak bergantung pada kolom baru. Tanpa 252 identitas polis tidak tersimpan.
- Diperbarui: `docs/STRUKTUR-TABEL-NB-FACIN.md` (+48 baris T_GENERAL_POLIS; penjaga `TestKolomDDLCocokDenganStruktur` sempat
  menolak sebelum baris ditambah), `loader/migrasi_sebagian_test.go` (berkas 252, selisih presisi 11 kolom, hitungan 458 -> 506,
  `TestMigrasi252GeneralPolisLengkap` dua arah).

## Tahap 1 — identitas polis, ditulis Selesai Policy (`repository/selesaipolicy.go` `tulisPolicy`)

`UPDATE T_GENERAL_POLIS SET POLICY_NO, PROD_DATE_TIME, PROD_KE WHERE ID = case`, dalam transaksi Submit yang sama, SESUDAH
JSON_POLIS ditulis / dibaca dan SEBELUM CEDING_FACINPRODUCTION.

| Kolom | Properti Pega | Sumber `[terverifikasi]` | Diisi? |
| --- | --- | --- | --- |
| POLICY_NO | `pyWorkPage.OfferFacIn.PolicyData.PolicyNo` | `NB FacIn\Activity\SaveJsonPolicyFacIn_Act.xml` langkah 12 (rakit `ParamSeq.CARI2 + BusinessOldId + "." + HASIL1 + "." + HASIL2`, PRE PolicyNo kosong); = JSON_POLIS.NOPOLIS (langkah 15 INSERTJSON_JSONPOLIS_FACIN) | **ya** — nomor yang dipakai / dibuat `nomorPolis` |
| PROD_DATE_TIME | `pyWorkPage.OfferFacIn.PolicyData.ProdDateTime` | SaveJsonPolicyFacIn_Act langkah 5.1-5.3 (CurrentDateTime / StartDateTime / awal bulan berikut 05:00 GMT); disimpan JSON_POLIS.TGL_PROD lewat CARI21 langkah 13 (`dd/MM/yyyy HH:mm:ss` Asia/Jakarta) | **ya** — teks Pega UTC `YYYYMMDDTHHMMSS.000 GMT` (pola START_DATE_TIME) dari TGL_PROD; milidetik `.000` `[keputusan agent]` (TGL_PROD DATE = sumber tersimpan, sama dengan backfill) |
| PROD_KE | (tidak ada medan `ProdKe` di korpus NB, 0 hit) | `RDBList\INSERTJSON_JSONPOLIS_FACIN.xml` argumen PRODKE `'0'`; aturan loader `T_GENERAL_POLIS.PROD_KE` = JSON_POLIS.PRODKE (K-071 J-3) | **ya** — `MAX(PRODKE)` JSON_POLIS case (dua bentuk IDPEGA); bukan angka -> NULL (bukan ORA-01722) |
| GENERATE_DATE_POLICY_NO | `.PolicyData.GenerateDatePolicyNo` | hanya `Activity\GenerateNopolis_Act.xml` langkah 10 `@CurrentDateTime()` tanpa PRE — pre-processing `FlowAction\ShowPolis.xml`, local action modal dari tombol di `Section\EmailSection.xml`; BUKAN jalur Selesai Policy (SaveJsonPolicyFacIn_Act tidak men-SET-nya) | **tidak** — NULL |
| POLICY_MASTER_NUMBER | `pyWorkPage.OfferFacIn.PolicyMasterNumber` (akar, bukan PolicyData) | nol Property-Set / DataTransform di korpus NB FacIn; hanya dibaca (Section Periode / Periode_IsUW / Summary*, InsertUploadMarineCargo_act 3.12.1, InsertTreatyProduction_Sql) | **tidak** — NULL (port FACINPRODUCTION juga NULL) |
| POLICY_MASTER_ID_PEGA | `pyWorkPage.OfferFacIn.PolicyMasterIDPega` (akar) | nol setter; hanya `Obj-Open-By-Handle` di CekLimitSpreading_Act langkah 8 / SetIndexMarineCargo_Act langkah 2 | **tidak** — NULL |
| BINDER_RNM | `pyWorkPage.OfferFacIn.BinderRNM` (akar) | nol setter; Section Periode_IsUW + InsertTreatyProduction_Sql / InsertFacinProductionBackup_Sql | **tidak** — NULL (port FACINPRODUCTION juga NULL) |
| DATE_OFFER | `.OfferFacIn.Date` (akar, medan rancangan `Date`) | `DataTransform\ViewOfferFacInUW_PreDT.xml` baris 2 `@DateTime.CurrentDate("MM/dd/yyyy hh:mm a","GMT")`, dipanggil `InwardFacultative_PreDT` baris 1 = SETIAP form InwardFacultative(_IsUW) dibuka; dibaca AddToListSuggestOfferFacIn_DT (ViewSuggest) | **tidak** — NULL; bukan saat Selesai Policy. Bila work owner mau "waktu form terakhir dibuka", itu perilaku layar (tahap 3) |

Uji: `TestIdentitasPolisTahap1` (teks SQL persis, 5 bind, konversi WIB -> UTC termasuk lintas tanggal, ORA-00904 dikenali,
urutan / argumen di `tulisPolicy`), `TestBackfillIdentitasPolis252`.

## Backfill polis yang sudah terbit — `docs/alat/isi_identitas_polis_252.sql`

Sesudah 252. Hitung sebelum + pratinjau (termasuk NB-151029) -> satu UPDATE -> hitung sesudah; TANPA COMMIT. Cakupan:
T_WORK_POLIS LINI 'FAC' ID 'NB-%', POLICY_NO masih NULL (tidak menimpa hasil biner baru), punya baris JSON_POLIS (IDPEGA dua
bentuk, TGL_INPUT terbaru). Rumus = penulisan Selesai Policy. FACINPRODUCTION tidak dipakai (JSON_POLIS selalu ditulis lebih dulu
dan memuat ketiga nilai). Kolom tahap 1 tanpa sumber terverifikasi tidak disentuh.

## Peta tahap 2 / 3 — 40 kolom, BELUM diimplementasi (untuk keputusan work owner)

Riset dua agen baca-saja; tiap agen menghitung setter dua cara (pengurai langkah vs grep `PropertiesName`, cocok); sampel 6
klaim dicek ulang sesi c3 (SetFlagSaveFO_act, CheckPeriodLife_ACT, SetValidateDate_Act CurrentYear / ProRateType,
GetRISlipDataFromDB_Act, SearchTreatyGroup_Act, IsB2B nol setter) — cocok. **Temuan umum `[terverifikasi]`: ke-40 medan
berada di AKAR `pyWorkPage.OfferFacIn.<Medan>`, bukan `PolicyData`** (`PolicyData\.(<medan>)` = 0 hit di NB FacIn dan DDL).
Tak satu pun dibaca RDBList INSERT kecuali IsB2B (InsertTreatyProduction_Sql / InsertViewSuggest_SQL); Pega menyimpannya di
DATA_JSON (`@ASM.GetPageJSONString()`, SaveJsonPolicyFacIn_Act langkah 14 / SaveJsonOfferFacIn_Act langkah 15).

Path relatif `D:\migrasi\RNM\NB FacIn\`. L = langkah. MATI = blok `//`.

| Kolom | Di-SET oleh | Kapan | Format | Sudah diport? |
| --- | --- | --- | --- | --- |
| ADDITIONAL_CAPITAL | input "%Additional Capital" Section InwardFacIn / OfferFacIn_NusaRe (change -> CountFormulaRNM_Act); nol Property-Set | layar detail Nusantara Re | angka (persen) | belum |
| CEDING_RETENTION | changepercentrnm(Life)_act L2 `CedingRetentionNominal / PersonList(1).CoverageList(1).TSI * 100`; input "% Ceding retention" (InputInwardFacultativeDtl / InputEndorsementDtl) | ubah Ceding retention | angka (persen) | belum |
| CEDING_RETENTION_NOMINAL | input "Ceding retention" InputInwardFacultativeDtl / InputEndorsementDtl (change -> changepercentrnmLife_act) | layar detail | angka | belum |
| COMMENT_CEDING | InputOfferFacInEngineer_preACT L4 `""` (PRE ReasonOfDecline kosong); SetCauseOfDecline_ACT L9 `ReasonOfDecline + "\n" + ReasonOfDecline2` (PRE IsCedingConfirm Binding) | buka form; Cause of Decline | teks | belum |
| CONDITION | SearchTreatyGroup_Act L3.1 `.OUTPARAM1` = SearchTreatyGroupName_SQL (TREATYGROUPNAME treatybusiness tahun berjalan + bizcode; baris terakhir menang `[dugaan]`) | buka form (preACT L11 / UW L13) | teks | belum |
| CURRENT_YEAR | SetValidateDate_Act L2 `FormatDateTime(StartDateTime,"yyyy")` | ubah tanggal Periode | "yyyy" | belum |
| FLAG_SAVE_FO | SetFlagSaveFO_act L1 `1` | tombol Save Fac Out (InputFacOffer) | 1 | belum |
| FLAG_SPECIAL_ACCEPTANCE | Protection_Act L6 `"false"` (PRE IsSpecialAcceptance != "true"); satu-satunya setter | proteksi (EmailSection*) | "false" | belum (proteksi.go tanpa medan ini) |
| FLAG_SPREADING | CopyToAllLocSpreading_ACT L4.1 / CopyToAllSpreading_ACT L5.1 `.FlagSpreading` baris 10015; dihapus ResetFlagSpreading_ACT | tombol salin spreading | "SF-MR 2023" / "Open Fac Market" (DDL\FlagSpreading) | sebagian: salinlokasi.go memvalidasi per baris, nilai akar tidak disimpan |
| FOLLOWING_NB | dropdown "Following NB" Periode; CopyLocationFollowing_Act L9 `.CaseID`; CopyNopolisFollowing_Act L3 `Param.NB` | tab Periode | nomor NB | belum (following.go: "tanpa kolom") |
| ID_FOLLOWING_NB | rule sama: `"ASM-FW-GISFW-WORK " + Param.NB` | tab Periode | IDPEGA berawalan | belum |
| INPUT_COMMENT_DEDUCTIBLE | text area Section InputCommentDeductible (local action dari Periode); nol Property-Set | local action | teks | belum |
| IS_B2_B | **nol setter** di korpus NB FacIn dan DDL; dibaca When IsPKSASM / IsT1T3 / IsT2T4 / ToDepHeadUW, InsertTreatyProduction_Sql | — | "ASM" / "SRB" / "" (pembanding) | Go menganggap non-ASM (submituw.go) |
| IS_BANDING | SetBandingProposal_DT baris 2 `"true"` (+ ProposalAcceptStatus 4); checkbox EmailSection_IsUW | Banding (EmailSection*) | "true" | belum disimpan (status 4 diport acceptance/keputusan.go) |
| IS_EDIT_REFF_NUMBER | checkbox PeriodeEndorsement saja | layar Endorsement | "true" | belum (khusus EDM) |
| IS_FAC_RETRO | SetDataFacOut_Act L3 `0`, SetDataFacOut{Fire L5, AnekaGolf L6, CargoMBU L4, PATravel L5} `1` (PRE isFacOut); CheckSpreadingProtect_ACT L21.1 `1` (TreatyType 10015) | post-activity InwardFacultative(_IsUW) | 0 / 1 | logika setara (submituw.go adaFacRetro), tidak disimpan |
| IS_GROSS_NET_SHOW | checkbox FormulaTreatyCapacityDesc(_IsUW) saja | layar formula | 'True' | belum |
| IS_INPUT_FAC_RETRO | SetFlagSaveFO_act L1 `1`; OfferFacOut_PostAct L2 `@if(.EmailTypeRetro=1,1,0)`; SetDataFacOut* `0` (PRE UW & =1); connector Flow InputInwardFacultativeOffer / RISlip `1` | Fac Out / alur | 0 / 1 | belum (Go: selalu tidak) |
| IS_LIFE_PERIOD | CheckPeriodLife_ACT L1 `"0"`, L4 `"1"` (PRE selisih tahun = Period) | ubah Period / tanggal | "0" / "1" | belum |
| IS_OCCUP_EXCEPTION | hanya `0` di setter hidup (preACT L5 / L38, UW preACT L1 / L25, Protection_Act L4); setter `1` hanya di blok MATI | buka form; proteksi | 0 | belum |
| IS_PREFERRED_RISK | SetFlagOccupation_ACT L1 `""`, L2.2.3.4 / L2.2.4.4 `Occupation.pxResults(1).CARI1`; ProtectFIREMBUPA_Act L6 / L9.2.27.7 / L9.2.28.7 | proteksi | "Preferred Risk" / "Non-Preferred Risk" / "" | logika ya (proteksi.go, submituw.go), tidak disimpan |
| IS_PRO_RATE | dropdown "Prorate / Short period" Periode; preACT L37 `"Prorate"` (NB); CountPremi_ACT L9 `"ShortPeriod"` (EDM); InsertUploadFire_act L5.3 | buka form / Periode | "Prorate" / "ShortPeriod" | belum |
| IS_RI_SLIP | GetRISlipDataFromDB_Act L4 `1` | cetak RI Slip | "1" | belum (When membaca saja) |
| IS_SPECIAL_ACCEPTANCE | checkbox "Policy Type" Periode / Periode_IsUW; satu-satunya Property-Set MATI | tab Periode | "true" / "false" | belum |
| IS_TOP_RISK | SetFlagOccupation_ACT L1 `""`, L2.1.1 `"true"`; ProtectFIREMBUPA_Act L6 / L9.1.1 | proteksi | "true" | belum (top risk per lokasi diport, medan akar tidak) |
| OFFICER_FAC_OUT | OfferFacOut_PreAct L22 (pemilik antrean FacOut / operator); PrintRISlipPre_act L3 | Fac Out / cetak RI Slip | teks operator | belum |
| OLD_POLIS_ID | tanpa medan; renewal saja (aturan loader K-071 J-5) | — | — | NULL di NB |
| PARAM_CONDITION | SearchTreatyGroup_Act L4 `"FIRE"` (IsFire), L5 `"ANEKA"` (IsAneka) | buka form | "FIRE" / "ANEKA" | belum |
| PRODUCT_BRIGUNA | dropdown "Product" Periode; nol Property-Set | tab Periode | "Karya,Umum,Talangan" | belum |
| PRO_RATE_TYPE | SetValidateDate_Act L2 `1`; dropdown "Coverage Rate Type" Periode / PeriodeEndorsement | ubah tanggal / Periode | 1 / 4 | belum |
| SAVE_FAC_OUT | tidak ditemukan di korpus NB FacIn (hanya 1 sampel CONTOH EDM) | — | "0" | belum |
| SHARE_RNML | nol setter; Number read-only "Share RNM" InputInwardFacultativeDtl / InputEndorsementDtl | — | "0" | belum; sumber `[belum terverifikasi]` |
| SOURCE_ID | medan Pega `ID` diganti nama loader (aturan.go medanID): SaveJsonOfferFacIn_Act L19 `OfferFacIn.ID := IDPEGAOUT` | simpan offer | `ASM-FW-GISFW-WORK NB-…` | loader saja |
| SUM_TOTAL_TSI | SumTSIObjItem_act L3 (Σ TSIObjectItem); SetTotalFormula_ACT L4 (TSI top risk / KURS) | Section PropertyItemList / formula | angka | belum |
| TOTAL_PREMI_NUSA_RE | 22 setter di CountPremiAndTSIRNM{FireMBU, FireAnekaGolf, PALifeCargo, Cargo}_ACT + …EDM_ACT (dispatcher CountPremiAndTSINusantaraRe_ACT) | hitung premi / buka form | angka | sebagian, menyimpang (cedant.go per mata uang -> IDR, K49-11) |
| TOTAL_RI_COM_NUSA_RE | satu-satunya setter SetRIComm_Act L9 — MATI | — | angka | belum (deduction.go port L5 / L8) |
| TOTAL_TSI_NUSA_RE | CountTotalTSIPremiNusaRe_Act L14; CountPremiAndTSIRNMPALifeCargo_ACT L5.3.3 (Life) | GetLimitAkseptasi* / buka form UW | angka | ya (submituw.go totalTSIUW), tidak disimpan |
| TOTAL_TSI_NUSA_RE_SPREADING | 21 setter, rule sama dengan TOTAL_PREMI_NUSA_RE | hitung premi | angka | sebagian, menyimpang (cedant.go) |
| TOTAL_TSI_TOP_RISK | CountTotalTSIPremiNusaRe_Act L14; CountPremiAndTSIRNMFireMBU_ACT L6.5 / L6.6; …EDM_ACT L4.5 | hitung premi / limit | angka | ya (submituw.go), tidak disimpan |
| TREATY_CAPACITY | SumTreatyCapacity_Act L9 `lRp / KURS` (PRE IsPEGAPROD); L10.4 MATI | formula / tanggal | angka | belum (spreading.go KAPASITAS_TREATY = rule lain) |

## Pertanyaan untuk work owner (dari peta)

1. Kolom bernilai hitungan (TOTAL_*, SUM_TOTAL_TSI, TREATY_CAPACITY): Pega menyimpan nilai terakhir di halaman; Go menghitung
   ulang saat dipakai (sebagian menyimpang per mata uang, K49-11). Disimpan dari hitungan Go, atau dibiarkan NULL?
2. Kolom input layar yang belum ada di layar Go (ADDITIONAL_CAPITAL, CEDING_RETENTION*, PRODUCT_BRIGUNA, IS_PRO_RATE,
   IS_SPECIAL_ACCEPTANCE, FOLLOWING_NB, INPUT_COMMENT_DEDUCTIBLE, …): perlu diport sebagai isian layar dulu.
3. Kolom tanpa setter hidup (IS_B2_B, SHARE_RNML, SAVE_FAC_OUT, TOTAL_RI_COM_NUSA_RE, IS_OCCUP_EXCEPTION selalu 0): biarkan NULL?
4. DATE_OFFER / GENERATE_DATE_POLICY_NO: perilaku layar (buka form / ShowPolis), bukan Selesai Policy — diisi saat Submit sebagai
   pendekatan, atau dibiarkan NULL?

## Keputusan work owner 12-10-2026 (ronde 2) — diimplementasi sesi c3

1. **Kolom hitungan → SIMPAN dari hitungan Go yang sudah ada** (`services/ringkasanpolis.go`, `repository/submit.go`
   `sqlRingkasanTotalPolis`). Ditulis **tiap Submit** (semua tahap) `[keputusan agent]`: setter Pega
   (CountPremiAndTSINusantaraRe_ACT / CountTotalTSIPremiNusaRe_Act) berjalan tiap form detail dibuka / diubah dan di
   GetLimitAkseptasi*, jadi nilai halaman saat Submit = nilai terakhir; Go menghitung saat dipakai, bukan per ubahan isian.
   - TOTAL_TSI_NUSA_RE / TOTAL_TSI_TOP_RISK = `totalTSIUW` (CountTotalTSIPremiNusaRe_Act L14, kurs terbaru; FIRE / Aneka / Golf /
     Cargo);
   - TOTAL_TSI_NUSA_RE_SPREADING / TOTAL_PREMI_NUSA_RE = Total RNM IDR tab Cedant (`ringkasanRNM` + `totalKeIDR`) —
     **MENYIMPANG K49-11**: dijumlah per mata uang lalu ke IDR (Pega menjumlah angka lintas mata uang tanpa kurs);
   - dibulatkan 8 desimal (mask bind `PolaAngkaMasuk`); kurs tidak ada / galat baca → total itu NULL + log, **Submit tidak gagal**;
     ORA-00904 (252 belum jalan) → log + lanjut (`jalankanKolom252`);
   - SUM_TOTAL_TSI / TREATY_CAPACITY **tidak ditulis**: Go belum punya hitungannya (SumTSIObjItem_act / SetTotalFormula_ACT /
     SumTreatyCapacity_Act belum diport) — NULL.
2. Kolom input layar yang belum ada di layar Go → **ditunda** (NULL).
3. Kolom tanpa setter hidup (IS_B2_B, SHARE_RNML, SAVE_FAC_OUT, TOTAL_RI_COM_NUSA_RE, IS_OCCUP_EXCEPTION) → **NULL**.
4. Pendekatan `[keputusan work owner]`:
   - DATE_OFFER = Submit tahap **marketing**, hanya bila masih NULL (`NVL(DATE_OFFER, :1)` = Submit marketing pertama), format
     ViewOfferFacInUW_PreDT `"MM/dd/yyyy hh:mm a"` GMT (Go `01/02/2006 03:04 PM` UTC);
   - GENERATE_DATE_POLICY_NO = Selesai Policy saat nomor polis **terbit** (nomor baru; JSON_POLIS lama dipakai ulang → kolom
     tetap, `NVL(:5, GENERATE_DATE_POLICY_NO)`), format GenerateNopolis_Act L10 `@CurrentDateTime()` = teks Pega UTC.

Uji: `TestRingkasanTotalPolis` (IDR 500 / 0 / 500 / premi = ringkasanRNM; top risk 8 desimal; USD tanpa kurs → empat NULL; tanpa
pembaca → nil), `TestDateOfferSubmit` (lintas tengah hari / tengah malam UTC; tahap lain ""), `TestAjukanSubmit` (DATE_OFFER +
Ringkasan terpasang; fixture tanpa kurs → NULL, Submit tetap sukses), `TestTotalDanDateOffer252`, `TestIdentitasPolisTahap1`
(6 bind, `terbit` hanya nomor baru). Mutasi 10 / 12 mati; 2 hidup = **mutan ekuivalen** (`|| len(tidakAda) > 0` — totalKeIDR
sudah mengembalikan nil bila ada kurs hilang; `case d.pesanKurs` — totalTSIUW keluar dini dengan total nil, simpan8(nil) = nil).

**Tahap 3 yang murah (BELUM diimplementasi, menunggu keputusan)** — setter-nya sudah diport dan nilainya tersedia di
`dasarJenjang` (dibaca Submit tahap UW / Policy):
- IS_PREFERRED_RISK = `preferredKasus` (SetFlagOccupation_ACT; nilai "Preferred Risk" / "Non-Preferred Risk" / "");
- IS_TOP_RISK = `d.adaTop` → "true" / "" (SetFlagOccupation_ACT L2.1.1);
- IS_FAC_RETRO = `d.adaFacRetro()` → 1 / 0 (SetDataFacOut*_Act).
Ketiganya hanya tersedia saat Submit UW / Policy; tahap marketing tidak menghitungnya.
