# Peta nama tabel datar — Treaty In

**Tanggal:** 23 September 2026 · **Skema:** `TREATY_MASUK`

> **TURUNAN.** Dibangkitkan [`alat/buat-peta-nama-tabel-treatyin.py`](../alat/buat-peta-nama-tabel-treatyin.py) dan **tidak pernah disunting tangan**. Bila berbeda dari definisinya, **alatnya yang salah**.
>
> Isinya diturunkan dari [`struktur-treatyin-lama.md`](struktur-treatyin-lama.md) — pohon `TreatyIn` hasil sapuan **329 berkas XML**.

Pendamping mesin-baca: [`peta-nama-tabel-treatyin.tsv`](peta-nama-tabel-treatyin.tsv).

---

## 0. Dua nama, dua lapisan

| Kolom | Lapisan | Tunduk pada |
|---|---|---|
| `NAMA_T` | **tabel datar** — bentuk pipih yang memetakan pohon clipboard Pega satu lawan satu | konvensi `T_` modul claim-non-prop |
| `PADANAN_DDL` | **objek Oracle** di skema `TREATY_MASUK` | §16 — kata utuh bahasa Indonesia, ≤30 bita |

Keduanya berlaku bersamaan dan **tidak saling menggantikan**. Awalan `T_` justru menandai bahwa barisnya adalah tabel datar, bukan entitas rancangan.

### Satu pengecualian: kepalanya memakai nama yang sudah ada

Kepala struktur ini **`TREATY_IN`**, tanpa awalan `T_`, karena ia **tabel Oracle yang sudah hidup sekarang** — bukan tabel datar yang kita turunkan. Awalan `T_` menandai sesuatu yang kita buat; menempelkannya pada tabel yang sudah ada akan mengaburkan justru hal yang paling perlu terlihat, yaitu bahwa tabel itu **bukan** milik kita.

Seluruh anaknya mengikuti struktur pohon dan memakai `T_…`, dengan kunci tamu `TREATY_IN_ID`.

> **20 lawan 99.** `TREATY_IN` yang sudah ada memuat **20 kolom bisnis** — terbaca dari argumen `POOLDATA.PEGA_TREATY_IN`, bukan ditebak. Halaman `TreatyIn` punya **99 skalar akar**. Selisih **79** itu tidak punya kolom di mana pun; ia hanya ada di dalam `M_TREATY_IN.JSONDATA`. Memakai nama tabel yang sudah ada **tidak** berarti memakai bentuk kolomnya yang sekarang.

---

## 1. Ringkasan

| | |
|---|---:|
| Tabel datar dibuat | **44** |
| Halaman yang **sengaja tidak** menjadi tabel | **2** |
| Kedalaman maksimum induk | **4** |

| Kelompok | Tabel |
|---|---:|
| Akar — tabel TREATY_IN yang sudah ada | 1 |
| Dipakai bersama modul lain | 2 |
| Layer dan rincian proporsional | 9 |
| Bagian dan penyebaran | 5 |
| Fakultatif | 6 |
| Daftar tingkat kontrak | 12 |
| Addendum dan selisih | 3 |
| Jembatan migrasi | 4 |
| Belum ada isinya | 2 |

---

## 2. Daftar tabel

### Akar — tabel TREATY_IN yang sudah ada

| `NAMA_T` | Padanan DDL | Induk | Kunci tamu | Asal di pohon Pega |
|---|---|---|---|---|
| `TREATY_IN` | KONTRAK + VERSI_KONTRAK | — | — | TreatyIn{} · 99 skalar akar · kelas ASM-FW-GISFW-Int-TREATY_IN · TABEL SUDAH ADA — 20 kolom bisnis nyata, 79 sisanya hanya di M_TREATY_IN.JSONDATA |

### Dipakai bersama modul lain

| `NAMA_T` | Padanan DDL | Induk | Kunci tamu | Asal di pohon Pega |
|---|---|---|---|---|
| `T_VIEW_COMMENT` | CATATAN_PERSETUJUAN | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.CommentList[] · OperatorName · Date · Suggest · IsApproved · kelas pinjaman SuggestList |
| `DOCUMENT_TREATY_IN` | DOKUMEN_KONTRAK | `TREATY_IN` | `TREATY_IN_ID` | M_ATTACHMENTTREATY_2 + T_STORAGE_IMAGE + CATEGORY_ATTACH_REAS — lewat RDB List, bukan lewat pohon clipboard |

### Layer dan rincian proporsional

| `NAMA_T` | Padanan DDL | Induk | Kunci tamu | Asal di pohon Pega |
|---|---|---|---|---|
| `T_TREATY_LIMITS` | LAYER | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.Limits[] · 18 field · ★ inti layer, n=307 |
| `T_TREATY_LIMIT_DETAIL` | DETAIL_PROPORSIONAL | `T_TREATY_LIMITS` | `LIMIT_ID` | TreatyIn.Limits[].Detail[] · 14 field · ★ inti proporsional, n=104 |
| `T_TREATY_LIMIT_COB` | KELAS_BISNIS_LAYER | `T_TREATY_LIMIT_DETAIL` | `LIMIT_DETAIL_ID` | Limits[].Detail[].COBList[] · ClassOfBusiness · ClassOfBusinessID |
| `T_TREATY_LIMIT_AMOUNT` | NILAI_LAYER | `T_TREATY_LIMIT_DETAIL` | `LIMIT_DETAIL_ID` | Detail[] · EPIList[] RetentionList[] CashLossList[] ClaimCoopList[] PLAList[] IOOLimitList[] — SATU tabel berkolom PERAN, bukan enam tabel |
| `T_TREATY_LIMIT_ACHIEVEMENT` | PENCAPAIAN | `T_TREATY_LIMIT_DETAIL` | `LIMIT_DETAIL_ID` | Limits[].Detail[].AchievementLists[] · n=1 · ⚠ kelas TIDAK dideklarasikan di mana pun |
| `T_TREATY_LIMIT_MEASURE` | BESARAN_LAYER | `T_TREATY_LIMITS` | `LIMIT_ID` | Limits[] · MDPList[] EgnpiTotalList[] PremiumEarnedList[] — SATU tabel berkolom PERAN |
| `T_TREATY_REINSTATEMENT` | PEMULIHAN_LIMIT | `T_TREATY_LIMITS` | `LIMIT_ID` | Limits[].Reinstatement_List[] · ReinstatementPct · AdditionalAmount1/2 · nama ber-garis-bawah di sumbernya |
| `T_TREATY_LIMIT_GROUP` | KELOMPOK_LAYER | `T_TREATY_LIMITS` | `LIMIT_ID` | Limits[].TreatyGroupList[] · TreatyGroup · TreatyGroupID |
| `T_TREATY_LIMIT_GROUP_COB` | KELAS_BISNIS_KELOMPOK | `T_TREATY_LIMIT_GROUP` | `LIMIT_GROUP_ID` | Limits[].TreatyGroupList[].ClassOfBusinessList[] ◄ KEDALAMAN 4 |

### Bagian dan penyebaran

| `NAMA_T` | Padanan DDL | Induk | Kunci tamu | Asal di pohon Pega |
|---|---|---|---|---|
| `T_TREATY_SHARE` | BAGIAN | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.Share[] · 10 field · ★ inti penyebaran, n=237 |
| `T_TREATY_SHARE_SPREADING` | PENYEBARAN_XOL | `T_TREATY_SHARE` | `SHARE_ID` | Share[].SpreadingListXOL[] · Pct · ReinsTypeName · ReinsTypeID, n=45 |
| `T_TREATY_SHARE_SPREAD_AMOUNT` | NILAI_TERSEBAR | `T_TREATY_SHARE_SPREADING` | `SHARE_SPREADING_ID` | 11 daftar RNMSpreadedList…[] di bawah SpreadingListXOL ◄ KEDALAMAN 4 — SATU tabel berkolom PERAN |
| `T_TREATY_SHARE_AMOUNT` | NILAI_BAGIAN | `T_TREATY_SHARE` | `SHARE_ID` | 17 daftar uang langsung di bawah Share[] — GrossPremiumList NetPremiumList RnmLimitList BrokerageList RNMSpreadedList… — SATU tabel berkolom PERAN |
| `T_TREATY_SHARE_DEDUCTION` | POTONGAN | `T_TREATY_SHARE` | `SHARE_ID` | Share[].DeductionList[] · Deduction · DeductionPct · Comment |

### Fakultatif

| `NAMA_T` | Padanan DDL | Induk | Kunci tamu | Asal di pohon Pega |
|---|---|---|---|---|
| `T_TREATY_FAC_SHARE` | BAGIAN_FAKULTATIF | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.FacultativeShareList[] · 9 field · kelas SAMA dengan Share[], n=58 |
| `T_TREATY_FAC_SHARE_AMOUNT` | NILAI_BAGIAN_FAKULTATIF | `T_TREATY_FAC_SHARE` | `FAC_SHARE_ID` | FacultativeShareList[] · GrossPremiumList[] NetPremiumList[] RnmLimitList[] RnmGrossPremiDisplay[] RnmLimitListDisplay[] — SATU tabel berkolom PERAN |
| `T_TREATY_FAC_SHARE_DEDUCTION` | POTONGAN_FAKULTATIF | `T_TREATY_FAC_SHARE` | `FAC_SHARE_ID` | FacultativeShareList[].DeductionList[] |
| `T_TREATY_FAC_REINSURER` | REASURADUR_FAKULTATIF | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.ShareFacultativeReinsurers[] (n=47) DAN FacultativeShareList[].ShareFacultativeReinsurers[] (n=6) — kelas sama, dua induk |
| `T_TREATY_FAC_LIMITS` | LAYER_FAKULTATIF | `T_TREATY_FAC_REINSURER` | `FAC_REINSURER_ID` | ShareFacultativeReinsurers[].FacultativeLimits[] · TreatyType · TreatyTypeID |
| `T_TREATY_FAC_LIMIT_DETAIL` | DETAIL_PROPORSIONAL_FAKULTATIF | `T_TREATY_FAC_LIMITS` | `FAC_LIMIT_ID` | FacultativeLimits[].Detail[] · 11 field ◄ KEDALAMAN 4 · kelas SAMA dengan Limits[].Detail[] |

### Daftar tingkat kontrak

| `NAMA_T` | Padanan DDL | Induk | Kunci tamu | Asal di pohon Pega |
|---|---|---|---|---|
| `T_TREATY_CURRENCY` | MATA_UANG_KONTRAK | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.CurrencyList[] · Currency · CurrencyID · Conversion · PeriodStart · PeriodEnd |
| `T_TREATY_EGNPI` | EGNPI | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.EGNPI[] · 9 field · Amount · AmountIDR · Proportion |
| `T_TREATY_RETENTION` | RETENSI | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.Retention[] · 6 field · ⚠ memakai kelas TreatyInEGNPI, bukan kelas sendiri |
| `T_TREATY_LIMIT_SUMMARY` | RINGKASAN_LIMIT | `TREATY_IN` | `TREATY_IN_ID` | LimitSummaryList[] LimitShareSummaryList[] LimitFacShareSummaryList[] MDPSummaryList[] — SATU tabel berkolom LINGKUP |
| `T_TREATY_INSTALLMENT` | ANGSURAN | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.Installment[] · Currency · ID |
| `T_TREATY_INSTALLMENT_ITEM` | RINCIAN_ANGSURAN | `T_TREATY_INSTALLMENT` | `INSTALLMENT_ID` | Installment[].InstallmentList[] · DueDate · InstallmentPct · PaymentDate · WPC   [bukti: identitas-kelas — bentuknya hanya terbaca lewat salinan OLDDATA/ValueDifference] |
| `T_TREATY_REPORTING_PERIOD` | PERIODE_PELAPORAN | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.ReportingPeriodList[] · Period · InitialDate · SubmissionDue · ConfirmationDue · SettlementDue |
| `T_TREATY_ACCUMULATION` | AKUMULASI | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.AccumulationList[] · Period · ReportDate |
| `T_TREATY_PORTFOLIO` | PORTOFOLIO | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.Portfolio[] · Type · TypePortfolio · Description |
| `T_TREATY_HAZARD_LIMIT` | BATAS_BAHAYA | `TREATY_IN` | `TREATY_IN_ID` | Earthquake FloodJab FloodNation RSMDLimit MaxCoGroup MaxCoNonGroup + Currency… — 10 skalar berpasangan DIPUTAR jadi daftar |
| `T_TREATY_TOTAL` | REKAP_KONTRAK | `TREATY_IN` | `TREATY_IN_ID` | 20 daftar Total…[] tingkat akar + TotalEgnpiAmountNP + TotalRetentionAmountNP — SATU tabel berkolom PERAN |
| `T_TREATY_RETRO_SHARE` | BAGIAN_RETRO | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.ShareReins[] + RetroList · ⚠ kelas TIDAK dideklarasikan di mana pun |

### Addendum dan selisih

| `NAMA_T` | Padanan DDL | Induk | Kunci tamu | Asal di pohon Pega |
|---|---|---|---|---|
| `T_TREATY_REVISION` | VERSI_KONTRAK · bagian addendum | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.OLDID · EDMState · EDMEffective · EDMMaterialType · RevisionState · RevisionDate · AddendumPremi · padanan lamanya TREATY_IN_EDM / M_TREATY_IN_EDM — TIDAK dipakai sebagai nama karena tabel itu bentuk SALINAN, bukan bagian addendum |
| `T_TREATY_VALUE_DIFFERENCE` | NILAI_SELISIH | `T_TREATY_REVISION` | `REVISION_ID` | TreatyIn.ValueDifference{} · 217 simpul · nilai sekarang − nilai lama · ditulis TreatyEDMCalculateDifference (11 langkah, SELURUHNYA hidup) |
| `T_TREATY_VALUE_BEFORE_PRORATE` | NILAI_SEBELUM_PRO_RATE | `TREATY_IN` | `TREATY_IN_ID` | TreatyIn.ValueBeforeProrate{} · 15 simpul · kelas sama dengan akar |
| — TIDAK DIBUAT — | OLDDATA tidak menjadi tabel | — | — | TreatyIn.OLDDATA{} · 142 simpul · ATURAN BISNIS: data lama TIDAK disimpan ulang — ia DISELECT |
| — TIDAK DIBUAT — | ActualValue tidak menjadi tabel | — | — | TreatyIn.ActualValue{} · 152 simpul · ⚠ namanya "nilai aktual" tetapi 8 isian menerima SELISIH — lihat struktur-treatyin-lama.md §5.3 |

### Jembatan migrasi

| `NAMA_T` | Padanan DDL | Induk | Kunci tamu | Asal di pohon Pega |
|---|---|---|---|---|
| `T_TREATY_MIG_CORRELATION` | MIGRASI_KORELASI | — | — | jembatan ke sistem lama — pzInsKey · TreatyIn.ID · TreatyIn.OLDID |
| `T_TREATY_MIG_LANDING` | MIGRASI_PENDARATAN | — | — | tempat bentuk lama mendarat utuh — M_TREATY_IN.JSONDATA apa adanya |
| `T_TREATY_MIG_REJECTED` | MIGRASI_NILAI_DITOLAK | `T_TREATY_MIG_LANDING` | `LANDING_ID` | nilai yang tidak dapat diurai — tercatat, tidak dibulatkan, tidak dibuang |
| `T_TREATY_OUTBOUND_ARCHIVE` | ARSIP_MUATAN_KELUAR | `TREATY_IN` | `TREATY_IN_ID` | 175 argumen ke PEGA_TREATY_IN / PEGA_M_TREATY_IN_EDM / PEGA_M_TREATY_IN_DETAIL / …_DETAIL_EDM |

### Belum ada isinya

| `NAMA_T` | Padanan DDL | Induk | Kunci tamu | Asal di pohon Pega |
|---|---|---|---|---|
| `T_TREATY_OFFER` | ⚠ TIDAK ADA PENULIS HIDUP | — | — | TREATYINOFFER — nol penulis terjangkau (sapuan dua tingkat). BUKAN sumber migrasi, BUKAN dasar rekonsiliasi |
| `T_TREATY_AUTHORITY_LIMIT` | BATAS_WEWENANG | — | — | TIDAK DITEMUKAN di pohon — batas wewenang persetujuan, menunggu dokumen |

---

## 3. Empat keputusan bentuk yang perlu dibaca, bukan disimpulkan

**(a) Daftar uang berbentuk sama digabung, dengan kolom PERAN.** Kelas `ASM-FW-GISFW-Data-TreatyInLimitsSpreading` muncul **120 kali** dan hampir seluruhnya berisi pasangan `Currency` + `Value`. Ia bukan entitas; ia bentuk pembawa satu nilai uang bermata uang. Membuat 120 tabel akan menyalin satu bentuk seratus dua puluh kali; membuat satu tabel tanpa kolom peran akan menumpuk besaran yang berbeda. Karena itu: satu tabel per **induk**, dengan kolom `PERAN` yang memegang nama daftar aslinya.

**(b) `OLDDATA` tidak menjadi tabel.** Aturan bisnisnya sudah ditetapkan: **data lama tidak disimpan ulang — ia DISELECT** dari versi yang disesuaikan. 142 simpul di bawah `TreatyIn.OLDDATA` karena itu tidak melahirkan satu kolom pun.

**(c) `ActualValue` juga tidak menjadi tabel — dan sebabnya berbeda.** Bukan karena aturan bisnis, melainkan karena **namanya tidak cocok dengan isinya**: delapan isian di `TreatyEDMDifferenceShare` menulis hasil pengurangan ke dalamnya. Membuat tabel bernama "nilai aktual" yang berisi selisih akan mewariskan kekeliruan itu. Adjudikasinya ditangguhkan ke [`TEMUAN-ADJUSTMENT-DITUNDA.md`](TEMUAN-ADJUSTMENT-DITUNDA.md).

**(d) Sepuluh skalar batas bahaya diputar menjadi daftar.** `Earthquake`, `FloodJab`, `FloodNation`, `RSMDLimit`, `MaxCoGroup`, `MaxCoNonGroup` masing-masing berpasangan dengan kolom mata uangnya sendiri. Bentuk lamanya memaksa satu kolom baru setiap kali ada jenis bahaya baru; `T_TREATY_HAZARD_LIMIT` membuatnya satu baris baru.

---

## 4. Yang TIDAK menjadi tabel, dan sebabnya

| Halaman / nama | Simpul | Sebab |
|---|---:|---|
| `TreatyIn.OLDDATA{}` | 142 | aturan bisnis — data lama **diselect**, tidak disalin |
| `TreatyIn.ActualValue{}` | 152 | nama tidak cocok dengan isinya (§3c) |
| `TREATYINOFFER` | — | **nol penulis terjangkau** di sistem berjalan |

Dua yang pertama **nyata dan hidup** di sistem lama; yang tidak dibawa adalah bentuk penyimpanannya, bukan datanya.
