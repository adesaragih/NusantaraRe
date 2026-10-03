# Revisi Penyimpanan PremiumList Life — JSON DIBUANG, polis → tabel baru relasional

Tanggal: 2026-09-16
Sumber: grilling struktur (Claude) + verifikasi korpus Kiro + **contoh `DATA_JSON` nyata dari work
owner** (satu polis `NBLF-30607`, Type TP, EDM/auto-retro) + keputusan work owner.
Status: `[keputusan work owner]` — konteks BARU (PremiumList Life), buang JSON → relasional.

> Pola sama Claim Life: buang JSON, simpan relasional di tabel BARU terpisah dari existing.
> Existing (`M_LIFE_PREMIUM_SUMMARY`, `M_LIFE_PREMIUM_DETAIL`, `LIFEINPRODUCTION`, `JSON_POLIS`,
> `JSON_OFFER_LIFE`) hanya dibaca saat migrasi.

---

## Keputusan inti `[keputusan work owner]`

1. **SEMUA JSON dibuang.** `JSON_POLIS.DATA_JSON` (CLOB) + `JSON_OFFER_LIFE` (CLOB) → dibuang. Data
   disimpan relasional/flat. Perakit `@ASM.GetPageJSONString()` tidak direplikasi.
2. **Seluruh polis → tabel BARU terpisah dari existing.** 7 tabel baru (di bawah).
3. **Kolom header = property yang di-SET aplikasi ke work object** (`RULE-OBJ-PROPERTY`), BUKAN dari
   isi CLOB mentah. Buang bawaan Pega (`px*`/`py*`/`pz*`) dan single-page kosong
   (`Policy`/`Quotation`/`TempError`). Aturan sama Master Product Name Life.
4. **FK antar tabel = `ON DELETE CASCADE`** + popup konfirmasi.
5. Uang→**NUMBER** (ADR-0003), tanggal→**DATE**, kolom **nullable**, PK **sequence** (ADR-0006).
6. **`T_WORK_POLIS`** = tabel terpisah, **LINTAS-LINI (semua polis, Life + Non-Life)** — status/posisi
   tangga proses. Pola sama `T_WORK_CLAIM`. Bukan anak `T_PREMIUM_LIST`.

⚠️ **Catatan skala (jutaan baris):** `T_PREMIUM_LIST_DETAIL`, `_SPREADING`, `_SPREADING_RETRO` bisa
sangat besar. Wajib **index pada FK**; pertimbangkan **partisi** (per periode/tahun). Hapus fisik
polis besar itu berat (cascade jutaan baris) → pertimbangkan arsip/soft-delete sebagai jalur normal.

---

## Struktur (7 tabel baru + T_WORK_POLIS) — terbukti dari `DATA_JSON` nyata

```
T_PREMIUM_LIST  (akar/header polis — PK ID)                     ← work object (field bisnis)
  ├─ T_PREMIUM_LIST_SUMMARY       1:N  FK PREMIUM_LIST_ID        ← CurrencyList (rekap per mata uang)
  ├─ T_PREMIUM_LIST_DETAIL        1:N  FK PREMIUM_LIST_ID        ← PremiumListDetail (per peserta)
  │     └─ T_PREMIUM_LIST_SPREADING       1:N  FK DETAIL_ID      ← SpreadingList (per treaty-year/jenis)
  │            └─ T_PREMIUM_LIST_SPREADING_RETRO  1:N  FK SPREADING_ID  ← RetroLifeList (per reinsurer)
  └─ T_VIEW_SUGGEST               1:N  FK PREMIUM_LIST_ID        ← OfferFacIn.ViewSuggest

T_WORK_POLIS  (PK ID)   ⬅ TABEL BARU, TERPISAH, LINTAS-LINI — status/posisi tangga proses
```

Semua FK `ON DELETE CASCADE` + popup. Pohon terbukti sama persis dengan `DATA_JSON` contoh:
`PremiumListSummary → {CurrencyList[], PremiumListDetail[] → SpreadingList[] → RetroLifeList[]}`,
`OfferFacIn → ViewSuggest[]`.

---

## Kolom per tabel `[terverifikasi dari DATA_JSON nyata + pyFields RULE-OBJ-PROPERTY]`

### T_PREMIUM_LIST (header) — field bisnis di-set aplikasi (buang px*/py*)
`ID` (PK), `PL_NUMBER` (dari PremiumListSummary), `RISLIPRNM` (dari PremiumListSummary — **opsional**;
ada di retro/EDM contoh TP, KOSONG di NB contoh QP),
`ANNUITY_INTEREST`, `BRANCH_CODE`, `BRANCH_NAME`, `BUSINESS_CODE`, `BUSINESS_NAME`, `CEDING_CO`,
`CEDING_CO_NAME`, `DATE_RECEIVED` DATE, `DESCRIPTION`, `EDIT_INPUT`, `EMAIL_TYPE_PL`,
`FLAG_ON_GOING_POLICY`, `IS_JSON_POLIS`, `JENIS_ASURANSI`, `MARKETING_CODE`, `MARKETING_NAME`,
`MO_ID`, `NO_OFFER`, `POLICY_HOLDER`, `POLICY_HOLDER_NAME`, `POSITION`, `PREMIUM_REFUND_FACTOR`,
`PRODUCT_NAME`, `PRODUCT_NAME_ID`, `PROPOSAL_ACCEPT_STATUS`, `PRO_RATE_TYPE`, `RETRO_ID`,
`RETRO_NAME`, `SECURITY_REINSURER`, `SECURITY_REINSURER_ID`, `SOB_NAME`, `SOURCE_OF_BUSINESS`,
`TYPE`, `TYPE_CEDING`, `WPC` DATE.
> ⚠️ BUANG: seluruh `px*`/`py*` (pxApplication, pyID=NBLF-30607, pyStatusWork, pyOrigDivision, dst =
> bawaan Pega/routing). ID work Pega diganti ID sequence baru.
> ⚠️ Nama menyesatkan: `IS_JSON_POLIS`/`EmailTypePL`/`EditInput` = flag proses; simpan apa adanya,
> arti dikonfirmasi work owner bila perlu.
> `CEDING_CO`(=`CedingCo`), `CEDING_CO_NAME`(=`CedingCoName`) DI SINI (header), bukan di anak.

### T_PREMIUM_LIST_SUMMARY (CurrencyList, 1:N per mata uang) — REKAP UANG DISIMPAN
`ID` (PK), `PREMIUM_LIST_ID` (FK), `CURRENCY`, `PREMIUM` NUMBER, `COMMISSION` NUMBER,
`BROKERAGE_FEE` NUMBER, `OVR_COMM` NUMBER, `TAX` NUMBER, `PROF_COMM` NUMBER, `CLAIM` NUMBER,
`CLAIM_AMOUNT` NUMBER, `BALANCE` NUMBER, `DEDUCTION` NUMBER, `RI_ADMIN_FEE` NUMBER,
`CEDING_RETENTION` NUMBER, `SUM_REASURED` NUMBER, `SUM_AT_RISK_GROSS` NUMBER, `SHARE_RETRO` NUMBER,
`SHARE_NUSANTARA_RE_GROSS` NUMBER,
+ kelompok `*_REFUND`: `GROSS_PREMIUM_REFUND`, `NET_PREMIUM_REFUND`, `COMM_REFUND`,
`BROKERAGE_FEE_REFUND`, `OVR_COMM_REFUND`, `TAX_REFUND`, `DEDUCTION_REFUND`, `RI_ADMIN_FEE_REFUND`,
+ kelompok `*_RETRO`: `GROSS_PREMIUM_RETRO`, `NET_PREMIUM_RETRO`, `DISCOUNT_PREMIUM_RETRO`,
`OVR_COMM_RETRO`, `BROKERAGE_FEE_RETRO`, `RI_ADMIN_FEE_RETRO`,
+ `*_REFUND_RETRO`: `GROSS_PREMIUM_REFUND_RETRO`, `NET_PREMIUM_REFUND_RETRO`,
`DISCOUNT_PREMIUM_REFUND_RETRO`, `OVR_COMM_REFUND_RETRO`, `BROKERAGE_FEE_REFUND_RETRO`,
`RI_ADMIN_FEE_REFUND_RETRO`. Semua uang NUMBER.
> ⚠️ Temuan: rekap uang per mata uang TERNYATA ADA di data (bukan cuma CURRENCY) → tabel penuh, bukan
> daftar. Mengoreksi dugaan awal "uang cuma ke Param".

### T_PREMIUM_LIST_DETAIL (PremiumListDetail, 1:N per peserta) — ~81 kolom
`ID` (PK), `PREMIUM_LIST_ID` (FK), `NAME_OF_INSURED`, `CERTIFICATE_NO`, `POLICY_NO`,
`POLICY_HOLDER`, `PLAN`, `SEX`, `DOB` DATE, `ENTRY_AGE`, `CURRENT_AGE`, `PERIOD_MM`, `CURRENCY`,
`BEGIN_DATE` DATE, `EFFECTIVE_DATE` DATE, `EXPIRED_DATE` DATE, `LAPSE_DATE` DATE,
`GROSS_VALUATION_BEGIN_DATE` DATE, `GROSS_VALUATION_EXPIRED_DATE` DATE,
`RETROCESSION_VALUATION_BEGIN_DATE` DATE, `RETROCESSION_VALUATION_EXPIRED_DATE` DATE, `WPC` DATE,
`FACTOR` NUMBER, `RATE` NUMBER, `EM_PERCENT` NUMBER,
`SUM_INSURED` NUMBER, `SUM_REASURED` NUMBER, `SUM_AT_RISK` NUMBER, `SUM_AT_RISK_GROSS` NUMBER,
`SUM_AT_RISK_RETRO` NUMBER, `CEDING_RETENTION` NUMBER, `SHARE_NUSANTARA_RE` NUMBER,
`SHARE_NUSANTARA_RE_GROSS` NUMBER, `SHARE_RETRO` NUMBER, `RETROCEDED_SHARE` NUMBER,
`CLAIM` NUMBER, `CLAIM_AMOUNT` NUMBER, `PROF_COMM` NUMBER,
+ kelompok uang gross/refund/retro (GROSS_PREMIUM*, NET_PREMIUM*, BROKERAGE_FEE*, OVR_COMM*, COMM*,
TAX*, DEDUCTION*, RI_ADMIN_FEE*, DISCOUNT_PREMIUM_*_RETRO) semua NUMBER.
> Acuan lengkap: 81 kolom `SaveMasterLPDet` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!ASM!SAVEMASTERLPDET`).
> `[terverifikasi 2 contoh DATA_JSON]` **NB (QP) vs EDM (TP):**
> - NB pakai jendela **GROSS** (`GROSS_VALUATION_BEGIN/EXPIRED_DATE`) + `LAPSE_DATE` + `PERIOD_MM`;
>   EDM/retro pakai jendela **RETROCESSION** (`RETROCESSION_VALUATION_*`). Tabel WAJIB punya **keduanya**.
> - `FACTOR` = NUMBER (contoh QP `1.0000000`, 7 desimal), bukan integer.
> Kolom EDM: `PL_NUMBER_EDM`, `EDMSTATUS`, `STATUSOLD`, `STATUS` (dari SaveMasterLPDet; tak muncul di
> dua contoh JSON karena keduanya jalur baru — muncul saat proses EDM/old data).

### T_PREMIUM_LIST_SPREADING (SpreadingList, 1:N per treaty-year/jenis)
`ID` (PK), `DETAIL_ID` (FK → T_PREMIUM_LIST_DETAIL), `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`
(QS/2ND QS/SURPLUS/2ND SURPLUS/OR), `TREATY_YEAR_LIFE`, `RETROCADED_SHARE` NUMBER, `IDR` NUMBER,
`USD` NUMBER, `IDR_SELISIH` NUMBER, `USD_SELISIH` NUMBER, `B_IDR` NUMBER, `B_USD` NUMBER,
`TGL_UPDATE` DATE, `USER_ID`.
> class asli `ASM-FW-GISFW-Int-TREATYYEAR_LIFE`. `ID` asli (1000032..) = kunci treaty year master;
> di tabel baru pakai PK sequence, `TREATY_TYPE_ID`/`TREATY_YEAR_LIFE` jadi kolom rujukan.
> ⚠️ Ini HASIL spreading yang DIBEKUKAN per peserta (keputusan work owner: simpan).

### T_PREMIUM_LIST_SPREADING_RETRO (RetroLifeList, 1:N per reinsurer)
`ID` (PK), `SPREADING_ID` (FK → T_PREMIUM_LIST_SPREADING), `REINSURER_NAME`, `PERCENT_SHARE` NUMBER,
`AMOUNT` NUMBER, `COMMISION` NUMBER, `OVR_COMM` NUMBER, `RATE` NUMBER,
`PREMIUM_SPREADED_GROSS` NUMBER, `PREMIUM_SPREADED_NET` NUMBER, `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`,
`TREATY_START_DATE` DATE, `TREATY_END_DATE` DATE, `TGL_UPDATE` DATE, `USER_ID`.
> class asli `ASM-FW-GISFW-Int-RETROCESSIONLIFE`. `IDTREATYYEAR_LIFE` di JSON = penghubung ke
> spreading induk → jadi FK `SPREADING_ID`.

### T_VIEW_SUGGEST (OfferFacIn.ViewSuggest, 1:N) — riwayat penawaran/konfirmasi ceding
`ID` (PK), `PREMIUM_LIST_ID` (FK), `NO`, `DATE_SUGGEST` DATE, `PIC_SUGGEST`, `IS_CEDING_CONFIRM`
(Accept/Reject/Decline), `COMMENT_SUGGEST`, `INITIAL` (Offer/Bind).
> Sumber: `AddHistorySuggest.xml` + `DATA_JSON` (OfferFacIn.ViewSuggest[]). OfferFacIn dibuang;
> ViewSuggest NAIK jadi anak header (keputusan work owner).

### T_WORK_POLIS (PK ID) — LINTAS-LINI, status/posisi tangga
Status connector (`Confirm`/`Decline`/`Reject`/`Offer`/`Premium`), `pyPosition`, dst. Kolom final
disamakan dengan `T_WORK_CLAIM` saat penyelarasan lintas-lini. Bukan anak `T_PREMIUM_LIST`.

---

## Penyimpangan sadar (⚠️)
1. Buang JSON (blob → relasional).
2. Rekap uang currency disimpan penuh (`T_PREMIUM_LIST_SUMMARY`).
3. Spreading + spreading_retro per peserta **DISIMPAN** (dibekukan; Pega lama tak menyimpan di
   PremiumList — dihitung di konteks Treaty/Master Retro). Menyimpang demi jejak historis.
4. `T_WORK_POLIS` lintas-lini untuk status proses.
5. Tipe uang→NUMBER, tanggal→DATE, nullable, PK sequence, FK cascade.
6. Catatan skala: index FK + partisi + arsip/soft-delete (bukan penyimpangan, kebutuhan non-fungsional).

## OQ tersisa
- **OQ-1 (kelengkapan kolom): DITUTUP `[keputusan work owner]` — dua contoh `DATA_JSON` diterima**
  (TP/EDM `NBLF-30607` + QP/NB `NBLF-33693` + QR/NB `NBLF-32940`). Temuan yang sudah dimasukkan:
  detail NB menambah `GROSS_VALUATION_BEGIN/EXPIRED_DATE`, `EFFECTIVE_DATE`, `LAPSE_DATE`, `PERIOD_MM`;
  `FACTOR`=NUMBER; `RISLIPRNM` + `RETRO_ID`/`RETRO_NAME`/`SECURITY_REINSURER*` opsional (kosong NB
  tanpa retro); `SpreadingList`/`RetroLifeList` bisa KOSONG (NB murni → tabel 1:N boleh nol);
  `PremiumListDetail` bisa banyak baris (QR: 5 peserta, 1 CurrencyList rekap).
  Pola `Q*`=jendela GROSS, `T*`=jendela RETROCESSION (konsisten dgn validasi DOL Claim Life).
  Sisa kecil non-pemblokir: kolom EDM `EDMSTATUS`/`STATUSOLD`/`STATUS` (dari `SaveMasterLPDet`) muncul
  saat proses old-data/EDM — dipastikan saat konteks Endorsement.
- **OQ-2 (tipe presisi): DITUTUP `[keputusan work owner]` — ikuti rekomendasi/standar proyek.** Uang =
  `NUMBER` presisi arbitrer (ADR-0003, `float` dilarang), tanggal = `DATE`, kolom nullable, PK
  sequence. Presisi fisik persis dicocokkan DBA saat tiket migrasi (bukan pemblokir).
- **OQ-3 (T_WORK_POLIS): DITUTUP untuk sekarang `[keputusan work owner]` — buat kolom PENTING dulu.**
  Minimum: `ID` (PK), identitas polis (kunci ke polis + lini), posisi/status tangga
  (`POSITION`/`STATUS` dari nilai connector `Confirm`/`Decline`/`Reject`/`Offer`/`Premium`),
  `pyPosition` bila perlu, audit (op/tanggal). Cascade: hapus polis → baris work ikut terhapus.
  Kolom lain **ditambah saat implementasi / konteks Non-Life** — tidak diblokir sekarang.

## Dampak ke artefak (langkah selanjutnya)
- Tulis **spec PremiumList Life** (baru) bagian penyimpanan + AC.
- **Slice skema/migrasi = PREFACTOR (tiket pertama)** karena bentuk berubah total.
- EDM (Endorsement Life) berbagi tabel yang sama (`PL_NUMBER_EDM`/`NOENDORS`/`PRODKE`/`EDMSTATUS`
  membedakan) — diselaraskan saat konteks Endorsement digarap.
