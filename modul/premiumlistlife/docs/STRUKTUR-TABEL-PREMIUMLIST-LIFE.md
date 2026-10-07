# Struktur Tabel — PremiumList Life (New Business + Endorsement)

Acuan bentuk tabel untuk aplikasi Go. Dibuat 2026-09-18 atas perintah work owner.
NB dan EDM memakai TABEL YANG SAMA — yang membedakan BARIS, bukan tabel.
Presisi fisik adalah `[data DBA]` dan TIDAK ditetapkan di sini.
Berkas ini menggambarkan BENTUK, bukan alasan — alasannya ada di `spec.md` dan `issues/`.

⚠️ **Berkas ini adalah acuan TUNGGAL nama kolom. Seluruh nama kolom snake_case**
`[keputusan work owner 2026-09-18]`. Ejaan yang muncul di `spec.md`, `issues/`, dan `revisi-*.md`
adalah **ejaan korpus** — itu **bukti asal kolom**, bukan nama kolom.

**Tipe ditulis sebagai kategori logis:** teks · angka desimal · bilangan bulat · DATE.
Seluruh uang dan share adalah **angka desimal**, tidak pernah float (**ADR-0003**).

**Kolom `Dipakai`** menyatakan baris versi mana yang mengisinya: **NB + EDM** · **NB saja** ·
**EDM saja**. Kolom ber-`EDM saja` **kosong** pada baris new business, dan sebaliknya.

**Sumber tiap kolom** adalah salah satu dari dua:

- **korpus** — kolom yang dipakai sistem berjalan, disertai nama rule-nya
- **keputusan** — kolom yang sudah ditetapkan di `spec.md` / `issues/`, disertai nomor tiketnya

Rule korpus yang dirujuk berulang di berkas ini:

| Singkatan | Rule | Tabel |
| --- | --- | --- |
| **SaveLIP-NB** | `PremiumList Life/RDBList/SaveLifeinProduction_SQL.xml` (36 kolom) | `POOLDATA.LIFEINPRODUCTION` |
| **SaveLIP-EDM** | `Endorsement Life/RDBList/SaveLifeinProduction_SQL.xml` (27 kolom) | `POOLDATA.LIFEINPRODUCTION` |
| **SaveLPD** | `RDBList/SaveMasterLPDet.xml` (80 kolom) — **identik byte demi byte** di kedua modul | `POOLDATA.M_LIFE_PREMIUM_DETAIL` |
| **InsUpload** | `PremiumList Life/RDBList/InsertDataUploadLife.xml` (7 kolom) | `POOLDATA.M_TEMPUPLOADLIFE` |

⚠️ Seluruh kolom **nullable** kecuali PK; wajib-isi ditegakkan di Go (`spec.md` §12, tiket 00).

---

## T_WORK_POLIS

Tabel work **lintas-lini** (semua polis, Life + Non-Life). Satu baris mewakili **satu work object
polis**. Bukan anak `T_PREMIUM_LIST`.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | NB + EDM | keputusan tiket 00 PremiumList |
| `LINI` | teks | ya | | NB + EDM | keputusan tiket 00 PremiumList — "identitas polis + lini" |
| `POSITION` | teks | ya | | NB + EDM | keputusan tiket 00 PremiumList, `spec.md` §12 — nilai connector `Confirm`/`Decline`/`Reject`/`Offer`/`Premium` |
| `STATUS_WORK` | teks | ya | | NB + EDM | korpus `pyWorkStatus` — keputusan tiket 00 PremiumList, `spec.md` §12; bernama `STATUS` sampai migrasi `059` (seragam dengan `T_WORK_CLAIM`, keputusan work owner 01-10-2026) |
| `FLAG_ONGOING_POLICY` | teks | ya | | NB | korpus `FlagOnGoingPolicy` — `CreateInputLife` b618, VERBATIM `"0"` (Input Offer) / `"1"` (Input Premium); migrasi `057`, butir **bn** (GILIRAN-13) |
| `COVER_KEY` | teks | ya | FK | NB + EDM | keputusan work owner 01-10-2026 — penunjuk kasus induk, sama dengan `T_WORK_CLAIM` butir d; **kosong** sampai ada modul yang terbukti mengisinya (nol `pxCoverInsKey` di XML PremiumList dan Endorsement); migrasi `059` |
| `CREATE_OP` | teks | ya | | NB + EDM | keputusan work owner 01-10-2026 — akun pembuat (`pxCreateOperator`), sama dengan `T_WORK_CLAIM`; kosong untuk baris yang lahir sebelum `059` |
| `CREATE_OP_NAME` | teks | ya | | NB + EDM | keputusan work owner 01-10-2026 — sama dengan `T_WORK_CLAIM`; diisi akun pembuat sampai login menyediakan nama tampilan; baris lama dari `T_PREMIUM_LIST.CREATE_OP_NAME` |
| `TGL_CREATE` | DATE | ya | | NB + EDM | keputusan work owner 01-10-2026 — waktu lahir kasus, sama dengan `T_WORK_CLAIM`; baris lama dari `T_PREMIUM_LIST.TGL_INPUT` |
| `TGL_UPDATE` | DATE | ya | | NB + EDM | keputusan work owner 01-10-2026 — waktu ubah terakhir baris kasus, sama dengan `T_WORK_CLAIM` |

⛔ Kolom **audit** yang tiket 00 sebut tanpa nama kini **dinamai** oleh keputusan work owner 01-10-2026:
`CREATE_OP`, `CREATE_OP_NAME`, `TGL_CREATE`, `TGL_UPDATE` — nama dan tipe sama dengan `T_WORK_CLAIM`.

**Index:** `IX_WORK_POLIS_COVER_KEY` (`COVER_KEY`).

**Relasi:**

- 1:1 dengan `T_PREMIUM_LIST` lewat **shared PK** — `ID` sama persis, **tanpa kolom penyambung**
- baris NB dan baris EDM **sejajar**, tidak saling menunjuk
- `COVER_KEY` → `T_WORK_POLIS.ID` (menunjuk dirinya sendiri), `FK_WORK_POLIS_COVER_KEY`, **tanpa `ON DELETE`** —
  menghapus induk yang masih ditunjuk ditolak Oracle, bukan ikut menghapus anak

---

## T_PREMIUM_LIST

Header polis. Satu baris mewakili **satu versi polis** — new business maupun endorsement. Seluruh
versi hidup berdampingan; **versi berjalan adalah baris ber-`PROD_KE` terbesar**.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | NB + EDM | keputusan tiket 00 PremiumList — shared PK = `T_WORK_POLIS.ID` |
| `ID_PEGA` | teks | ya | | NB + EDM | korpus `IDPEGA` — SaveLIP-NB, SaveLIP-EDM |
| `NO_POLIS` | teks | ya | | NB + EDM | korpus `NOPOLIS` — SaveLIP-NB, SaveLIP-EDM |
| `BUSINESS_CODE` | teks | ya | | NB + EDM | korpus `BUSINESSCODE` — SaveLIP-NB, SaveLIP-EDM |
| `BUSINESS_NAME` | teks | ya | | NB + EDM | korpus `BUSINESSNAME` — SaveLIP-NB, SaveLIP-EDM |
| `CEDING_CO` | teks | ya | | NB + EDM | korpus `CEDINGCO` — SaveLIP-NB, SaveLIP-EDM |
| `CEDING_CO_NAME` | teks | ya | | NB + EDM | korpus `CEDINGCONAME` — SaveLIP-NB, SaveLIP-EDM |
| `DATE_RECEIVED` | DATE | ya | | NB + EDM | korpus `DATERECEIVED` — SaveLIP-NB, SaveLIP-EDM |
| `MARKETING_CODE` | teks | ya | | NB + EDM | korpus `MARKETINGCODE` — SaveLIP-NB, SaveLIP-EDM |
| `MARKETING_NAME` | teks | ya | | NB + EDM | korpus `MARKETINGNAME` — SaveLIP-NB, SaveLIP-EDM |
| `POLICY_HOLDER` | teks | ya | | NB + EDM | korpus `POLICYHOLDER` — SaveLIP-NB, SaveLIP-EDM |
| `POLICY_HOLDER_NAME` | teks | ya | | NB + EDM | korpus `POLICYHOLDERNAME` — SaveLIP-NB, SaveLIP-EDM |
| `PRO_RATE_TYPE` | teks | ya | | NB + EDM | korpus `PRORATETYPE` — SaveLIP-NB, SaveLIP-EDM |
| `CREATE_OP_NAME` | teks | ya | | NB + EDM | korpus `CREATEOPNAME` — SaveLIP-NB, SaveLIP-EDM |
| `SOB` | teks | ya | | NB + EDM | korpus `SOB` — SaveLIP-NB, SaveLIP-EDM |
| `SOB_NAME` | teks | ya | | NB + EDM | korpus `SOBNAME` — SaveLIP-NB, SaveLIP-EDM |
| `TYPE` | teks | ya | | NB + EDM | korpus `TYPE` — SaveLIP-NB, SaveLIP-EDM |
| `TYPE_CEDING` | teks | ya | | NB + EDM | korpus `TYPECEDING` — SaveLIP-NB, SaveLIP-EDM · domain nilai `1`=QS `2`=SURPLUS `3`=QS+SURPLUS `4`=XOL — keputusan tiket 00 Endorsement |
| `TYPE_CEDING_NAME` | teks | ya | | NB + EDM | korpus `TYPECEDINGNAME` — SaveLIP-NB, SaveLIP-EDM |
| `MO_ID` | teks | ya | | NB + EDM | korpus `MOID` — SaveLIP-NB, SaveLIP-EDM |
| `NO_OFFER` | teks | ya | | NB + EDM | korpus `NOOFFER` — SaveLIP-NB, SaveLIP-EDM |
| `RI_SLIP_RNM` | teks | ya | | NB + EDM | korpus `RISLIPRNM` — SaveLIP-NB, SaveLIP-EDM |
| `RETRO_ID` | teks | ya | | NB + EDM | korpus `RETROID` — SaveLIP-NB, SaveLIP-EDM |
| `RETRO_NAME` | teks | ya | | NB + EDM | korpus `RETRONAME` — SaveLIP-NB, SaveLIP-EDM |
| `SECURITY_REINSURER_ID` | teks | ya | | NB + EDM | korpus `SECURITYREINSURERID` — SaveLIP-NB, SaveLIP-EDM |
| `SECURITY_REINSURER` | teks | ya | | NB + EDM | korpus `SECURITYREINSURER` — SaveLIP-NB, SaveLIP-EDM |
| `TGL_INPUT` | DATE | ya | | NB + EDM | korpus `TGL_INPUT` — SaveLIP-NB, SaveLIP-EDM |
| `DESCRIPTION` | teks | ya | | **NB saja** | korpus `DESCRIPTION` — SaveLIP-NB |
| `WPC` | DATE | ya | | **NB saja** | korpus `WPC` — SaveLIP-NB |
| `PRODUCT_NAME_ID` | teks | ya | | **NB saja** | korpus `PRODUCTNAMEID` — SaveLIP-NB |
| `PRODUCT_NAME` | teks | ya | | **NB saja** | korpus `PRODUCTNAME` — SaveLIP-NB |
| `LAYER_1` | teks | ya | | **NB saja** | korpus `LAYER_1` — SaveLIP-NB |
| `LAYER_2` | teks | ya | | **NB saja** | korpus `LAYER_2` — SaveLIP-NB |
| `LAYER_3` | teks | ya | | **NB saja** | korpus `LAYER_3` — SaveLIP-NB |
| `LAYER_4` | teks | ya | | **NB saja** | korpus `LAYER_4` — SaveLIP-NB |
| `ANNUITY_INTEREST` | angka desimal | ya | | **NB saja** | korpus `ANNUITYINTEREST` — SaveLIP-NB |
| `PREMIUM_REFUND_FACTOR` | angka desimal | ya | | **NB saja** | korpus `PREMIUMREFUNDFACTOR` — SaveLIP-NB |
| `NO_ENDORS` | teks | ya | | **EDM saja** | korpus `NOENDORS` — SaveLIP-EDM |
| `EDM_TYPE` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement — `1`=Perubahan Data, `3`=Batal |
| `OLD_POLICY_NO` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `EDM_DATE` | DATE | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `EDM_NOTE` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `PREMI_PROPOSED` | angka desimal | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `UANG_PERTANGGUNGAN` | angka desimal | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `SUM_INSURED` | angka desimal | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `JENIS_PRODUK` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `SISTEM_REASURANSI` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `STATUSS` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `STATUS_UPDATE` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `STATUS_SERVICE` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `START_DATE` | DATE | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `END_DATE` | DATE | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `PL_NUMBER_EDM` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `PROD_KE` | bilangan bulat | ya | | **EDM saja** | keputusan tiket 00 Endorsement — versi berjalan = `PRODKE` terbesar |
| `EDM_STATUS` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement — `Old`/`New`/`Delete`/`Batal` |
| `STATUS_OLD` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `BATAS_USIA_PESERTA` | bilangan bulat | ya | | **NB saja** | korpus `.BatasUsiaPeserta` "Age Limit" — `Section/InputOfferLife.xml`; migrasi 059 |
| `PERIODE_PERTANGGUNGAN` | teks | ya | | **NB saja** | korpus `.PeriodePertanggungan` "Coverage Period" — InputOfferLife.xml; migrasi 059 |
| `TANGGAL_PENAWARAN` | DATE | ya | | **NB saja** | korpus `.TanggalPenawaran` "Offering Date" — InputOfferLife.xml; migrasi 059 |
| `TANGGAL_RESPON` | DATE | ya | | **NB saja** | korpus `.TanggalRespon` "Response Date" — InputOfferLife.xml; migrasi 059 |
| `TANGGAL_KONFIRMASI` | DATE | ya | | **NB saja** | korpus `.TanggalKonfirmasi` "Confirmation Date" — InputOfferLife.xml; migrasi 059 |
| `TBC` | bilangan bulat | ya | | **NB saja** | korpus `.TBC` "Input TBC" — InputOfferLife.xml; migrasi 059 |
| `TANGGAL_TBC` | DATE | ya | | **NB saja** | korpus `.TanggalTBC` "Max TBC" = TanggalKonfirmasi + TBC hari — `SetMaxTBCLife_Act`; migrasi 059 |
| `KETERANGAN_MARKETING` | teks | ya | | **NB saja** | korpus `.KeteranganMarketing` "Marketing Note" — InputOfferLife.xml; migrasi 059 |
| `QQ_NAME` | teks | ya | | **NB saja** | korpus `.QQName` "Insured Name" — InputOfferLife.xml; migrasi 060 |
| `JENIS_USAHA` | teks | ya | | **NB saja** | korpus `.JenisUsaha` "Occupation" — InputOfferLife.xml; migrasi 060 |
| `KETENTUAN_UNDERWRITING` | teks | ya | | **NB saja** | korpus `.KetentuanUnderwriting` "Underwriting Policy" — InputOfferLife.xml; migrasi 060 |
| `TANGGAL_KONFIRMASI_BALIK` | DATE | ya | | **NB saja** | korpus `.TanggalKonfirmasiBalik` "Re-Confirmation Date" — InputOfferLife.xml; migrasi 060 |
| `TANGGAL_REALISASI` | DATE | ya | | **NB saja** | korpus `.TanggalRealisasi` "Realization Date" — InputOfferLife.xml; migrasi 060 |
| `TANGGAL_BIND` | DATE | ya | | **NB saja** | korpus `.TanggalBind` "Binding Date" — InputOfferLife.xml; migrasi 060 |
| `STATUS_FINAL` | teks | ya | | **NB saja** | korpus `.StatusFinal` "Final Status" — InputOfferLife.xml; migrasi 060 |
| `JENIS_ASURANSI` | teks | ya | | **NB saja** | korpus `.JenisAsuransi` "Reinsurance Type" — InputOfferLife.xml; turunan `SetReinsuranceType` (`TypeCeding "4"` → Non Proportional, selain itu Proportional), disimpan seperti `Obj-Save` InputOfferLife_ACT langkah 8; migrasi 061 |
| `STATUS_PENAWARAN` | teks | ya | | **NB saja** | radio "Status" layar Input Offer (`pyWorkPage.Status` / `EmailTypePL`, kode radio) — keputusan work owner 01-10-2026, tidak berkolom di Pega; SATU BARIS PER STATUS per kasus (baris utama `ID = ID_PEGA` = status terakhir); migrasi 062 |


`[keputusan work owner]` 2026-09-18 — **`TYPE_CEDING` adalah satu kolom.** Tiket 00 Endorsement
mendokumentasikan **domain nilainya** (`1`–`4`), bukan menambah kolom. `[terverifikasi]`
`PremiumList Life/Activity/InputOfferLife_ACT.xml`.

**Index:** `NO_POLIS` · `PROD_KE` (rantai versi dicari lewat pasangan keduanya).

**Relasi:**

- induknya `T_WORK_POLIS` · **shared PK**, 1:1 · tidak ada kolom penyambung
- anaknya `T_PREMIUM_LIST_DETAIL` lewat `PREMIUM_LIST_ID` · 1:N · ON DELETE **CASCADE**
- anaknya `T_PREMIUM_LIST_SUMMARY` lewat `PREMIUM_LIST_ID` · 1:N · ON DELETE **CASCADE**
- anaknya `T_VIEW_SUGGEST` lewat `PREMIUM_LIST_ID` · 1:N · ON DELETE **CASCADE**

⚠️ `[keputusan work owner]` **Kolom `WORK_POLIS_ID` DIBUANG** — bukan diganti nama, **tidak ada**.
Hubungan `T_WORK_POLIS` ↔ `T_PREMIUM_LIST` dijamin oleh **`ID` yang identik**.

⚠️ `[keputusan work owner 01-10-2026]` **Satu baris per status penawaran** (migrasi 062): satu kasus
dapat punya beberapa baris. **Baris utama** `ID = ID_PEGA` = nomor kasus memuat status TERAKHIR dan
tetap satu-satunya yang dibaca kotak masuk, Premium List Detail, penomoran, summary, Claim Life, dan
FK tabel anak. **Baris status** (`ID` = 32 heksa dari nomor kasus + status, `ID_PEGA` = nomor kasus)
hanya salinan isian penawaran — tanpa `NO_POLIS`/`TYPE`/kolom tahap lain. Lihat
`backend/repository/polis_barisstatus.go`.

⚠️ **Bentuk `NO_ENDORS`** `[terverifikasi]` — dicatat sebagai keterangan, **bukan** kolom baru:

```
NO_ENDORS = <NO_POLIS> / <PROD_KE dua digit>

rantainya:  GetProdKeOldData_SQL      -> Local.Prodke
            GenerateNoEDM_Life        -> CARI4 = Prodke + 1
                                         CARI14 = CARI4 dipad nol jadi 2 digit
            Generate_NoEndorsmentLife -> NO_POLIS || '/' || CARI14
```

---

## T_PREMIUM_LIST_DETAIL

Peserta polis. Satu baris mewakili **satu peserta pada satu versi polis**.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | NB + EDM | korpus `ID` — SaveLPD |
| `PREMIUM_LIST_ID` | teks | ya | FK | NB + EDM | keputusan tiket 00 PremiumList — → `T_PREMIUM_LIST.ID` |
| `PARENT_ID` | teks | ya | FK | **EDM saja** | keputusan tiket 00 Endorsement — self-reference → `T_PREMIUM_LIST_DETAIL.ID` |
| `ID_PEGA` | teks | ya | | NB + EDM | korpus `IDPEGA` — SaveLPD |
| `PL_NUMBER` | teks | ya | | NB + EDM | korpus `PL_NUMBER` — SaveLPD |
| `POLICY_NO` | teks | ya | | NB + EDM | korpus `POLICY_NO` — SaveLPD |
| `POLICY_HOLDER` | teks | ya | | NB + EDM | korpus `POLICY_HOLDER` — SaveLPD |
| `CERTIFICATE_NO` | teks | ya | | NB + EDM | korpus `CERTIFICATE_NO` — SaveLPD |
| `NAME_OF_INSURED` | teks | ya | | NB + EDM | korpus `NAME_OF_INSURED` — SaveLPD |
| `DESCRIPTION` | teks | ya | | NB + EDM | korpus `DESCRIPTION` — SaveLPD |
| `SEX` | teks | ya | | NB + EDM | korpus `SEX` — SaveLPD |
| `DOB` | DATE | ya | | NB + EDM | korpus `DOB` — SaveLPD |
| `AGE` | bilangan bulat | ya | | NB + EDM | korpus `AGE` — SaveLPD |
| `ENTRY_AGE` | bilangan bulat | ya | | NB + EDM | korpus `ENTRY_AGE` — SaveLPD |
| `CURRENT_AGE` | bilangan bulat | ya | | NB + EDM | korpus `CURRENT_AGE` — SaveLPD |
| `PLAN` | teks | ya | | NB + EDM | korpus `PLAN` — SaveLPD |
| `RISK` | teks | ya | | NB + EDM | korpus `RISK` — SaveLPD |
| `MEDICAL_STATUS` | teks | ya | | NB + EDM | korpus `MEDICAL_STATUS` — SaveLPD |
| `STNC` | teks | ya | | NB + EDM | korpus `STNC` — SaveLPD |
| `WPC` | teks | ya | | NB + EDM | korpus `WPC` — SaveLPD |
| `CURRENCY` | teks | ya | | NB + EDM | korpus `CURRENCY` — SaveLPD |
| `PERIOD_YY` | bilangan bulat | ya | | NB + EDM | korpus `PERIOD_YY` — SaveLPD |
| `PERIOD_MM` | bilangan bulat | ya | | NB + EDM | korpus `PERIOD_MM` — SaveLPD |
| `PASSED_PERIOD` | bilangan bulat | ya | | NB + EDM | korpus `PASSED_PERIOD` — SaveLPD |
| `BEGIN_DATE` | DATE | ya | | NB + EDM | korpus `BEGIN_DATE` — SaveLPD |
| `EFFECTIVE_DATE` | DATE | ya | | NB + EDM | korpus `EFFECTIVE_DATE` — SaveLPD |
| `EXPIRED_DATE` | DATE | ya | | NB + EDM | korpus `EXPIRED_DATE` — SaveLPD |
| `LAPSE_DATE` | DATE | ya | | NB + EDM | korpus `LAPSE_DATE` — SaveLPD |
| `GROSS_VALUATION_BEGIN_DATE` | DATE | ya | | NB + EDM | korpus `GROSS_VALUATION_BEGIN_DATE` — SaveLPD |
| `GROSS_VALUATION_EXPIRED_DATE` | DATE | ya | | NB + EDM | korpus `GROSS_VALUATION_EXPIRED_DATE` — SaveLPD |
| `RETRO_VALUATION_BEGIN_DATE` | DATE | ya | | NB + EDM | korpus `RETRO_VALUATION_BEGIN_DATE` — SaveLPD |
| `RETRO_VALUATION_EXPIRED_DATE` | DATE | ya | | NB + EDM | korpus `RETRO_VALUATION_EXPIRED_DATE` — SaveLPD |
| `SUM_INSURED` | angka desimal | ya | | NB + EDM | korpus `SUM_INSURED` — SaveLPD |
| `SUM_REASURED` | angka desimal | ya | | NB + EDM | korpus `SUM_REASURED` — SaveLPD |
| `SUM_AT_RISK_GROSS` | angka desimal | ya | | NB + EDM | korpus `SUM_AT_RISK_GROSS` — SaveLPD |
| `SUM_AT_RISK_RETRO` | angka desimal | ya | | NB + EDM | korpus `SUM_AT_RISK_RETRO` — SaveLPD |
| `CEDING_RETENTION` | angka desimal | ya | | NB + EDM | korpus `CEDING_RETENTION` — SaveLPD |
| `CEDING_CO` | teks | ya | | NB + EDM | korpus `CEDING_CO` — SaveLPD |
| `SHARE_NUSANTARA_RE` | angka desimal | ya | | NB + EDM | korpus `SHARE_NUSANTARA_RE` — SaveLPD |
| `SHARE_NUSANTARA_RE_GROSS` | angka desimal | ya | | NB + EDM | korpus `SHARE_NUSANTARA_RE_GROSS` — SaveLPD |
| `SHARE_RETRO` | angka desimal | ya | | NB + EDM | korpus `SHARE_RETRO` — SaveLPD |
| `RETROCEDED_SHARE` | angka desimal | ya | | NB + EDM | korpus `RETROCEDED_SHARE` — SaveLPD |
| `RATE` | angka desimal | ya | | NB + EDM | korpus `RATE` — SaveLPD |
| `FACTOR` | angka desimal | ya | | NB + EDM | korpus `FACTOR` — SaveLPD |
| `EM_PERCENT` | angka desimal | ya | | NB + EDM | korpus `EM_PERCENT` — SaveLPD |
| `PRO_RATE_TYPE` | teks | ya | | NB + EDM | korpus `PRORATETYPE` — SaveLPD |
| `GROSS_PREMIUM` | angka desimal | ya | | NB + EDM | korpus `GROSS_PREMIUM` — SaveLPD |
| `NET_PREMIUM` | angka desimal | ya | | NB + EDM | korpus `NET_PREMIUM` — SaveLPD |
| `COMM` | angka desimal | ya | | NB + EDM | korpus `COMM` — SaveLPD |
| `PROF_COMM` | angka desimal | ya | | NB + EDM | korpus `PROF_COMM` — SaveLPD |
| `OVR_COMM` | angka desimal | ya | | NB + EDM | korpus `OVR_COMM` — SaveLPD |
| `BROKERAGE_FEE` | angka desimal | ya | | NB + EDM | korpus `BROKERAGE_FEE` — SaveLPD |
| `TAX` | angka desimal | ya | | NB + EDM | korpus `TAX` — SaveLPD |
| `FLEET_DISCOUNT` | angka desimal | ya | | NB + EDM | korpus `FLEET_DISCOUNT` — SaveLPD |
| `DEDUCTION` | angka desimal | ya | | NB + EDM | korpus `DEDUCTION` — SaveLPD |
| `RI_ADMIN_FEE` | angka desimal | ya | | NB + EDM | korpus `RI_ADMIN_FEE` — SaveLPD |
| `CLAIM` | angka desimal | ya | | NB + EDM | korpus `CLAIM` — SaveLPD |
| `CLAIM_AMOUNT` | angka desimal | ya | | NB + EDM | korpus `CLAIM_AMOUNT` — SaveLPD |
| `GROSS_PREMIUM_REFUND` | angka desimal | ya | | NB + EDM | korpus `GROSS_PREMIUM_REFUND` — SaveLPD |
| `NET_PREMIUM_REFUND` | angka desimal | ya | | NB + EDM | korpus `NET_PREMIUM_REFUND` — SaveLPD |
| `COMM_REFUND` | angka desimal | ya | | NB + EDM | korpus `COMM_REFUND` — SaveLPD |
| `OVR_COMM_REFUND` | angka desimal | ya | | NB + EDM | korpus `OVR_COMM_REFUND` — SaveLPD |
| `BROKERAGE_FEE_REFUND` | angka desimal | ya | | NB + EDM | korpus `BROKERAGE_FEE_REFUND` — SaveLPD |
| `TAX_REFUND` | angka desimal | ya | | NB + EDM | korpus `TAX_REFUND` — SaveLPD |
| `DEDUCTION_REFUND` | angka desimal | ya | | NB + EDM | korpus `DEDUCTION_REFUND` — SaveLPD |
| `RI_ADMIN_FEE_REFUND` | angka desimal | ya | | NB + EDM | korpus `RI_ADMIN_FEE_REFUND` — SaveLPD |
| `GROSS_PREMIUM_RETRO` | angka desimal | ya | | NB + EDM | korpus `GROSS_PREMIUM_RETRO` — SaveLPD |
| `NET_PREMIUM_RETRO` | angka desimal | ya | | NB + EDM | korpus `NET_PREMIUM_RETRO` — SaveLPD |
| `DISCOUNT_PREMIUM_RETRO` | angka desimal | ya | | NB + EDM | korpus `DISCOUNT_PREMIUM_RETRO` — SaveLPD |
| `OVR_COMM_RETRO` | angka desimal | ya | | NB + EDM | korpus `OVR_COMM_RETRO` — SaveLPD |
| `BROKERAGE_FEE_RETRO` | angka desimal | ya | | NB + EDM | korpus `BROKERAGE_FEE_RETRO` — SaveLPD |
| `RI_ADMIN_FEE_RETRO` | angka desimal | ya | | NB + EDM | korpus `RI_ADMIN_FEE_RETRO` — SaveLPD |
| `GROSS_PREMIUM_REFUND_RETRO` | angka desimal | ya | | NB + EDM | korpus `GROSS_PREMIUM_REFUND_RETRO` — SaveLPD |
| `NET_PREMIUM_REFUND_RETRO` | angka desimal | ya | | NB + EDM | korpus `NET_PREMIUM_REFUND_RETRO` — SaveLPD |
| `DISCOUNT_PREMIUM_REFUND_RETRO` | angka desimal | ya | | NB + EDM | korpus `DISCOUNT_PREMIUM_REFUND_RETRO` — SaveLPD |
| `OVR_COMM_REFUND_RETRO` | angka desimal | ya | | NB + EDM | korpus `OVR_COMM_REFUND_RETRO` — SaveLPD |
| `BROKERAGE_FEE_REFUND_RETRO` | angka desimal | ya | | NB + EDM | korpus `BROKERAGE_FEE_REFUND_RETRO` — SaveLPD |
| `RI_ADMIN_FEE_REFUND_RETRO` | angka desimal | ya | | NB + EDM | korpus `RI_ADMIN_FEE_REFUND_RETRO` — SaveLPD |
| `PL_NUMBER_EDM` | teks | ya | | **EDM saja** | korpus `PL_NUMBER_EDM` — SaveLPD |
| `EDM_STATUS` | teks | ya | | **EDM saja** | korpus `EDMSTATUS` — SaveLPD |
| `STATUS_OLD` | teks | ya | | **EDM saja** | korpus `STATUSOLD` — SaveLPD |
| `STATUS` | teks | ya | | **EDM saja** | korpus `STATUS` — SaveLPD |

**Index:** `PREMIUM_LIST_ID` · `PARENT_ID`.

**Relasi:**

- induknya `T_PREMIUM_LIST` lewat `PREMIUM_LIST_ID` · 1:N · ON DELETE **CASCADE**
- induk dari dirinya sendiri lewat `PARENT_ID` → `T_PREMIUM_LIST_DETAIL.ID` · 1:N · nullable
- anaknya `T_PREMIUM_LIST_SPREADING` lewat `DETAIL_ID` · 1:N · ON DELETE **CASCADE**
- anaknya `DOCUMENT_POLIS` · 1:N · ON DELETE `[data DBA]`

---

## T_PREMIUM_LIST_SPREADING

Hasil spreading seorang peserta, dipecah per treaty-year/jenis. Satu baris mewakili **satu
treaty-year** dari satu peserta. Nilainya **dibekukan** saat polis disimpan.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | NB + EDM | keputusan tiket 00 PremiumList |
| `DETAIL_ID` | teks | ya | FK | NB + EDM | keputusan tiket 00 PremiumList — → `T_PREMIUM_LIST_DETAIL.ID` |
| `TREATY_TYPE_ID` | teks | ya | | NB + EDM | keputusan `spec.md` §12 |
| `TREATY_TYPE_NAME` | teks | ya | | NB + EDM | keputusan `spec.md` §12 — QS / 2ND QS / SURPLUS / 2ND SURPLUS / OR |
| `TREATY_YEAR_LIFE` | teks | ya | | NB + EDM | keputusan `spec.md` §12 |
| `RETROCADED_SHARE` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `IDR` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `USD` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `IDR_SELISIH` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `USD_SELISIH` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `B_IDR` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `B_USD` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `TGL_UPDATE` | DATE | ya | | NB + EDM | keputusan `spec.md` §12 |
| `USER_ID` | teks | ya | | NB + EDM | keputusan `spec.md` §12 |

**Index:** `DETAIL_ID`.

**Relasi:**

- induknya `T_PREMIUM_LIST_DETAIL` lewat `DETAIL_ID` · 1:N · ON DELETE **CASCADE**
- anaknya `T_PREMIUM_LIST_SPREADING_RETRO` lewat `SPREADING_ID` · 1:N · ON DELETE **CASCADE**
- ⛔ **tidak** punya `PARENT_ID` — ikut tersalin di bawah peserta versi baru

---

## T_PREMIUM_LIST_SPREADING_RETRO

Pecahan spreading per reinsurer. Satu baris mewakili **satu reinsurer** pada satu baris spreading.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | NB + EDM | keputusan tiket 00 PremiumList |
| `SPREADING_ID` | teks | ya | FK | NB + EDM | keputusan tiket 00 PremiumList — → `T_PREMIUM_LIST_SPREADING.ID` |
| `REINSURER_NAME` | teks | ya | | NB + EDM | keputusan `spec.md` §12 |
| `PERCENT_SHARE` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `AMOUNT` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `COMMISION` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 — ejaan `COMMISION` (sic), dipertahankan |
| `OVR_COMM` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `RATE` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `PREMIUM_SPREADED_GROSS` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `PREMIUM_SPREADED_NET` | angka desimal | ya | | NB + EDM | keputusan `spec.md` §12 |
| `TREATY_TYPE_ID` | teks | ya | | NB + EDM | keputusan `spec.md` §12 |
| `TREATY_TYPE_NAME` | teks | ya | | NB + EDM | keputusan `spec.md` §12 |
| `TREATY_START_DATE` | DATE | ya | | NB + EDM | keputusan `spec.md` §12 |
| `TREATY_END_DATE` | DATE | ya | | NB + EDM | keputusan `spec.md` §12 |
| `TGL_UPDATE` | DATE | ya | | NB + EDM | keputusan `spec.md` §12 |
| `USER_ID` | teks | ya | | NB + EDM | keputusan `spec.md` §12 |

**Index:** `SPREADING_ID`.

**Relasi:**

- induknya `T_PREMIUM_LIST_SPREADING` lewat `SPREADING_ID` · 1:N · ON DELETE **CASCADE**
- tidak punya anak · ⛔ **tidak** punya `PARENT_ID`

---

## T_PREMIUM_LIST_SUMMARY

Rekap uang polis **per mata uang**. Satu baris mewakili **satu mata uang pada satu versi polis**.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | NB + EDM | keputusan tiket 00 PremiumList |
| `PREMIUM_LIST_ID` | teks | ya | FK | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` — → `T_PREMIUM_LIST.ID` |
| `CURRENCY` | teks | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `PL_NUMBER` | teks | ya | | NB | nomor PL rekap, diisi saat Confirm — keputusan work owner 03-10-2026, menggantikan "tanpa PL_NUMBER" tiket 05a; migrasi 064 |
| `BALANCE` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `PREMIUM` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `COMMISSION` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `PROF_COMM` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `OVR_COMM` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `TAX` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `CLAIM` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `CLAIM_AMOUNT` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `DEDUCTION` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `BROKERAGE_FEE` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `RI_ADMIN_FEE` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `CEDING_RETENTION` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `SUM_REASURED` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `SUM_AT_RISK_GROSS` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `SHARE_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `SHARE_NUSANTARA_RE_GROSS` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `GROSS_PREMIUM_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `NET_PREMIUM_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `DISCOUNT_PREMIUM_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `OVR_COMM_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `BROKERAGE_FEE_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `RI_ADMIN_FEE_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `GROSS_PREMIUM_REFUND` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `NET_PREMIUM_REFUND` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `COMM_REFUND` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `OVR_COMM_REFUND` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `TAX_REFUND` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `BROKERAGE_FEE_REFUND` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `DEDUCTION_REFUND` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `RI_ADMIN_FEE_REFUND` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `GROSS_PREMIUM_REFUND_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `NET_PREMIUM_REFUND_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `DISCOUNT_PREMIUM_REFUND_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `OVR_COMM_REFUND_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `BROKERAGE_FEE_REFUND_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |
| `RI_ADMIN_FEE_REFUND_RETRO` | angka desimal | ya | | NB + EDM | keputusan `revisi-penyimpanan-premiumlist.md` §`T_PREMIUM_LIST_SUMMARY` |

⚠️ **Kolom di atas berasal dari sisi KEPUTUSAN, bukan dari korpus.**
`PremiumList Life/RDBList/InsertPLSummary.xml` memanggil `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY`
**secara posisional** — 37 argumen masuk + 2 keluar — sehingga **nama kolomnya ada di dalam
procedure** dan **tidak terbaca dari korpus**. Yang terbaca adalah daftar di
`revisi-penyimpanan-premiumlist.md`. Pemetaan argumen posisional → kolom tetap `[data DBA]`.

**Index:** `PREMIUM_LIST_ID`.

**Relasi:**

- induknya `T_PREMIUM_LIST` lewat `PREMIUM_LIST_ID` · 1:N · ON DELETE **CASCADE**
- ⛔ **tidak** punya `PARENT_ID` — **tidak disalin** antar versi

---

## T_VIEW_SUGGEST

Riwayat penawaran dan konfirmasi ceding. Satu baris mewakili **satu langkah penawaran** pada satu
versi polis.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | NB + EDM | keputusan tiket 00 PremiumList |
| `PREMIUM_LIST_ID` | teks | ya | FK | NB + EDM | keputusan tiket 00 PremiumList — → `T_PREMIUM_LIST.ID` |
| `NO` | bilangan bulat | ya | | NB + EDM | keputusan `spec.md` §12; `PremiumList Life/Activity/AddHistorySuggest.xml` |
| `DATE_SUGGEST` | DATE | ya | | NB + EDM | keputusan `spec.md` §12; `AddHistorySuggest.xml` |
| `PIC_SUGGEST` | teks | ya | | NB + EDM | keputusan `spec.md` §12; `AddHistorySuggest.xml` |
| `IS_CEDING_CONFIRM` | teks | ya | | NB + EDM | keputusan `spec.md` §12 — `Accept`/`Reject`/`Decline` |
| `COMMENT_SUGGEST` | teks | ya | | NB + EDM | keputusan `spec.md` §12; `AddHistorySuggest.xml` |
| `INITIAL_SUGGEST` | teks | ya | | NB + EDM | keputusan `spec.md` §12 — `Offer`/`Bind`. ⛔ **Diralat 28-09-2026 (pl6)**: namanya semula `INITIAL`, **kata cadangan Oracle** — migrasi `056` gagal di DEV. Nama baru mengikuti pola saudaranya `DATE_SUGGEST`/`PIC_SUGGEST`/`COMMENT_SUGGEST` |

**Index:** `PREMIUM_LIST_ID`.

**Relasi:**

- induknya `T_PREMIUM_LIST` lewat `PREMIUM_LIST_ID` · 1:N · ON DELETE **CASCADE**

---

# Di luar pohon

## DOCUMENT_POLIS

Dokumen pendukung polis. Tabel **lintas-lini**.

⛔ **Kolomnya tidak ditulis di sini.** `[data DBA]` Kelasnya `ASM-FW-GISFW-Int-DOCUMENT_POLIS`;
SQL-nya dibuat Pega sendiri, sehingga daftar kolomnya tidak dapat diturunkan dari korpus maupun dari
keputusan yang sudah ada.

**Relasi:** menggantung pada `T_PREMIUM_LIST_DETAIL` · 1:N · ON DELETE `[data DBA]`.

---

## M_TEMPUPLOADLIFE

Tabel **singgah** unggah CSV peserta. Satu baris mewakili **satu baris CSV yang belum diproses**.
Ia **di luar pohon polis** — tidak punya FK ke `T_PREMIUM_LIST` dan tidak ikut cascade.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `NO` | bilangan bulat | ya | | NB + EDM | korpus `NO` — InsUpload |
| `INSURED` | teks | ya | | NB + EDM | korpus `INSURED` — InsUpload |
| `DOB` | DATE | ya | | NB + EDM | korpus `DOB` — InsUpload |
| `BEGIN_DATE` | DATE | ya | | NB + EDM | korpus `BEGINDATE` — InsUpload |
| `END_DATE` | DATE | ya | | NB + EDM | korpus `ENDDATE` — InsUpload |
| `POLIVYHOLDER` | teks | ya | | NB + EDM | korpus `POLIVYHOLDER` — InsUpload (ejaan korpus, sic) |
| `CEDING_RETENTION` | angka desimal | ya | | NB + EDM | korpus `CEDING_RETENTION` — InsUpload |

**Index:** tidak ada.

**Relasi:** tidak punya induk maupun anak. Isinya dibersihkan oleh
`RDBList/DeleteTempUploadDataLife.xml` — **identik di kedua modul**.

---

## Pohon relasi — lima tingkat

```
TINGKAT 1   T_WORK_POLIS ─────────── mandiri · LINTAS-LINI · satu baris per work object
            PK  ID
            LINI · POSITION · STATUS_WORK · FLAG_ONGOING_POLICY
            CREATE_OP · CREATE_OP_NAME · TGL_CREATE · TGL_UPDATE   <- seragam T_WORK_CLAIM (059)
            FK COVER_KEY -> T_WORK_POLIS.ID   <- self, nullable, tanpa ON DELETE
            |
            |  shared PK: ID sama persis, TANPA kolom penyambung (WORK_POLIS_ID dibuang)
            v
TINGKAT 2   T_PREMIUM_LIST                         satu baris per VERSI polis
            PK ID = T_WORK_POLIS.ID  <- SHARED PK
            NO_POLIS · PROD_KE   <- versi berjalan = PROD_KE terbesar
            NO_ENDORS = NO_POLIS/<PROD_KE dua digit>   (baris EDM saja)
            |
TINGKAT 3   +--1:N-- T_PREMIUM_LIST_DETAIL                      CASCADE
            |        PK ID · FK PREMIUM_LIST_ID
            |        FK PARENT_ID -> T_PREMIUM_LIST_DETAIL.ID   <- self, nullable, EDM saja
            |        |
TINGKAT 4   |        +--1:N-- T_PREMIUM_LIST_SPREADING          CASCADE
            |        |        FK DETAIL_ID
            |        |        |
TINGKAT 5   |        |        +--1:N-- T_PREMIUM_LIST_SPREADING_RETRO   CASCADE
            |        |                 FK SPREADING_ID
            |        |
            |        +--1:N-- DOCUMENT_POLIS        <- kolom [data DBA]
            |
TINGKAT 3   +--1:N-- T_PREMIUM_LIST_SUMMARY                     CASCADE
            |        FK PREMIUM_LIST_ID · kolom [data DBA]
            |        (TIDAK disalin antar versi)
            |
TINGKAT 3   +--1:N-- T_VIEW_SUGGEST                             CASCADE
                     FK PREMIUM_LIST_ID

DI LUAR POHON:  M_TEMPUPLOADLIFE  (tabel singgah unggah CSV, tanpa FK)
```

## Tabel relasi

| # | Induk | Anak | Kunci tamu | Kard | Hapus |
| --- | --- | --- | --- | --- | --- |
| 1 | `T_WORK_POLIS` | `T_PREMIUM_LIST` | **tidak ada kolom terpisah** — `T_PREMIUM_LIST.ID` = `T_WORK_POLIS.ID` (**shared PK**) | 1:1 | — |
| 2 | `T_PREMIUM_LIST` | `T_PREMIUM_LIST_DETAIL` | `PREMIUM_LIST_ID` | 1:N | CASCADE |
| 3 | `T_PREMIUM_LIST` | `T_PREMIUM_LIST_SUMMARY` | `PREMIUM_LIST_ID` | 1:N | CASCADE |
| 4 | `T_PREMIUM_LIST` | `T_VIEW_SUGGEST` | `PREMIUM_LIST_ID` | 1:N | CASCADE |
| 5 | `T_PREMIUM_LIST_DETAIL` | `T_PREMIUM_LIST_SPREADING` | `DETAIL_ID` | 1:N | CASCADE |
| 6 | `T_PREMIUM_LIST_SPREADING` | `T_PREMIUM_LIST_SPREADING_RETRO` | `SPREADING_ID` | 1:N | CASCADE |
| 7 | `T_PREMIUM_LIST_DETAIL` | `T_PREMIUM_LIST_DETAIL` | `PARENT_ID` (self-reference, nullable) | 1:N | `[data DBA]` |
| 8 | `T_PREMIUM_LIST_DETAIL` | `DOCUMENT_POLIS` | `[data DBA]` | 1:N | `[data DBA]` |
| 9 | `T_WORK_POLIS` | `T_WORK_POLIS` | `COVER_KEY` (self-reference, nullable) | 1:N | ditolak (tanpa `ON DELETE`, sama dengan `T_WORK_CLAIM`) |

Relasi **1** tidak punya kunci tamu untuk di-index — **shared primary key**, dan PK sudah ber-index
dengan sendirinya. Seluruh kunci tamu lain **ber-index**.

---

## Catatan — belum ditetapkan, TIDAK menghambat berkas ini

- `[terverifikasi]` `COMMISION` (satu S) dan `COMMISSION` (dua S) **SAMA-SAMA ejaan korpus** — 369
  berkas lawan 83 berkas. Bukan salah ketik artefak. `T_PREMIUM_LIST_SPREADING_RETRO` memakai
  `COMMISION`, `T_PREMIUM_LIST_SUMMARY` memakai `COMMISSION`, keduanya mengikuti sumbernya
  masing-masing. **JANGAN diseragamkan tanpa keputusan work owner.**
- `[terbuka]` `POLIVYHOLDER` `[terverifikasi]` ejaan korpus — `InsertDataUploadLife.xml` di **kedua** modul. Ditiru atau diperbaiki belum diputuskan
- `[terbuka]` `STATUSS` dua huruf S di akhir, tetapi `[terverifikasi]` **itu ejaan korpus** — 73 berkas XML. Bukan salah ketik artefak. Perbaikan belum diputuskan
- `[data DBA]` daftar kolom `M_LIFE_PREMIUM_SUMMARY` + isi procedure `PEGA_M_LIFE_PREMIUM_SUMMARY`
- `[data DBA]` daftar kolom `DOCUMENT_POLIS`
- `[data DBA]` presisi fisik seluruh kolom
- `[terbuka]` delapan kolom produk/layer (`PRODUCT_NAME_ID`, `PRODUCT_NAME`, `LAYER_1..4`, `ANNUITY_INTEREST`, `PREMIUM_REFUND_FACTOR`) tidak ditulis jalur EDM — kosong di baris versi EDM
- `[terbuka]` rantai versi dicari lewat (`NO_POLIS`, `PROD_KE`), belum ada kolom penunjuk versi sebelumnya
- `[terbuka]` polis lama hasil migrasi belum punya baris `T_WORK_POLIS`, padahal PK-nya diambil dari sana
- `[terbuka]` format identitas kerja `EDMLF-<n>`
- `[terbuka]` kolom `T_PREMIUM_LIST_SPREADING`, `_SPREADING_RETRO`, dan `T_VIEW_SUGGEST` **sudah ditetapkan** di `spec.md` §12 dan `revisi-penyimpanan-premiumlist.md`; yang belum ditetapkan hanya presisi dan nullability per kolom
- ~~`[terbuka]` nama kolom audit pada `T_WORK_POLIS` belum ditetapkan; tiket 00 hanya menyebut "audit"~~ —
  **ditutup 01-10-2026** (keputusan work owner): `CREATE_OP`, `CREATE_OP_NAME`, `TGL_CREATE`, `TGL_UPDATE`, migrasi `059`
- `[terbuka]` `DESCRIPTION` dan `WPC` juga tidak ditulis jalur EDM, di luar delapan kolom produk/layer
- `[terbuka]` header memakai ejaan korpus (`BUSINESS_CODE`, `CEDING_CO`, `ANNUITY_INTEREST`) sedangkan `revisi-penyimpanan-premiumlist.md` memakai snake_case (`BUSINESS_CODE`, `CEDING_CO`, `ANNUITY_INTEREST`)
- `[terbuka]` `WPC` bertipe DATE di header polis tetapi teks di peserta
- `[terbuka]` `SUM_INSURED`, `PL_NUMBER_EDM`, `EDM_STATUS`, `STATUS_OLD` ada di header **dan** di peserta
- `[terbuka]` `ON DELETE` untuk `PARENT_ID` dan untuk `DOCUMENT_POLIS` belum ditetapkan
- `[terbuka]` `M_TEMPUPLOADLIFE` tidak punya penunjuk ke polis maupun ke pengunggah
