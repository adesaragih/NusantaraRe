# Struktur data sistem lama — **Treaty In Adjustment**

> ### POTRET SISTEM LAMA, 24 September 2026. BUKAN RANCANGAN.
>
> Skema yang dibangun ada di 2-to-spec/KAMUS-KOLOM.md dan 2-to-spec/ddl-usulan/. Gambar skema barunya: 4-erd-dan-tabel-datar/ERD-SKEMA-BARU.xlsx

**Dibangkitkan:** [`alat/pisah-struktur-erd-v2-per-modul.py`](../treaty-in/alat/pisah-struktur-erd-v2-per-modul.py) · **Sumber tunggal:** `alat/erd-v2.json` — **sumber yang sama** dengan berkas gabungannya, `treaty-in/4-erd-dan-tabel-datar/STRUKTUR-DATA-ERD-V2.md`.

> **Berkas ini tidak ditulis tangan, dan tidak boleh disunting.** Ia bagian **Treaty In Adjustment** dari potret `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`. Sunting sumbernya, lalu jalankan ulang alatnya.

---

## 0. Bagaimana pemisahannya dilakukan, dan apa yang TIDAK dipisah

**Sumbu pemisahnya satu, dan ia mekanis:** nama lembar yang memuat **`EDM`** masuk modul
**Treaty In Adjustment**; selain itu masuk modul **Treaty In**. Tidak ada aturan kedua.

| Modul | Lembar |
|---|---|
| Treaty In | *Treaty In Prop* · *Treaty In Non Prop* |
| **Treaty In Adjustment** | *Treaty In EDM Prop* · *Treaty In EDM Non Prop* |

> ### YANG DIPISAH ADALAH GAMBARNYA, BUKAN MODELNYA
>
> **34 dari 37 tabel dipakai KEDUA modul**, dan karena itu muncul di **kedua** berkas.
> Itu **bukan penggandaan** melainkan kenyataan sistem lama: satu tabel yang sama dipakai
> jalur kontrak **dan** jalur addendum. Kolom **Bersama modul lain** menandainya per baris.
>
> | | |
> |---|---:|
> | tabel unik di **Treaty In Adjustment** | **35** |
> | — di antaranya **juga** di Treaty In | **34** |
> | — **hanya** di Treaty In Adjustment | **1** — `T_TREATY_VALUE_BEFORE_PRORATE` |
> | tabel unik di kedua modul digabung | **37** |

---

## 1. Hitungan

| | Jumlah |
|---|---:|
| Tabel unik di modul ini | **35** |
| Kotak tergambar (satu tabel dapat muncul di beberapa lembar) | **60** |
| Relasi yang menyentuh lembar modul ini | **34** |
| Relasi **tidak terklasifikasi** — dibawa ke kedua berkas | **3** |
| Lembar *Treaty In EDM Prop* | 25 kotak |
| Lembar *Treaty In EDM Non Prop* | 35 kotak |

## 2. Entitas — satu baris per tabel unik

| Tabel | Kard. | Kunci utama | Kunci asing | ON DELETE | Induk | Muncul di lembar | Bersama Treaty In |
|---|---|---|---|---|---|---|---|
| `TREATY_IN` | — | `ID` | — | — | **akar** | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_REVISION` | 1:1 | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_LIMITS` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_LIMIT_DETAIL` | 1:N | *tidak dinyatakan* | `LIMIT_ID` | CASCADE | `T_TREATY_LIMITS` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_LIMIT_COB` | 1:N | *tidak dinyatakan* | `LIMIT_DETAIL_ID` | CASCADE | `T_TREATY_LIMIT_DETAIL` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_LIMIT_AMOUNT` | 1:N | *tidak dinyatakan* | `LIMIT_DETAIL_ID` | CASCADE | `T_TREATY_LIMIT_DETAIL` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_LIMIT_ACHIEVEMENT` | 1:N | *tidak dinyatakan* | `LIMIT_DETAIL_ID` | CASCADE | `T_TREATY_LIMIT_DETAIL` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_LIMIT_GROUP` | 1:N | *tidak dinyatakan* | `LIMIT_ID` | CASCADE | `T_TREATY_LIMITS` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_LIMIT_GROUP_COB` | 1:N | *tidak dinyatakan* | `LIMIT_GROUP_ID` | CASCADE | `T_TREATY_LIMIT_GROUP` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_SHARE` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_SHARE_SPREADING` | 1:N | *tidak dinyatakan* | `SHARE_ID` | CASCADE | `T_TREATY_SHARE` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_SHARE_SPREAD_AMOUNT` | 1:N | *tidak dinyatakan* | `SHARE_SPREADING_ID` | CASCADE | `T_TREATY_SHARE_SPREADING` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_SHARE_AMOUNT` | 1:N | *tidak dinyatakan* | `SHARE_ID` | CASCADE | `T_TREATY_SHARE` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_SHARE_DEDUCTION` | 1:N | *tidak dinyatakan* | `SHARE_ID` | CASCADE | `T_TREATY_SHARE` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_FAC_SHARE` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_FAC_SHARE_AMOUNT` | 1:N | *tidak dinyatakan* | `FAC_SHARE_ID` | CASCADE | `T_TREATY_FAC_SHARE` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_FAC_SHARE_DEDUCTION` | 1:N | *tidak dinyatakan* | `FAC_SHARE_ID` | CASCADE | `T_TREATY_FAC_SHARE` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_RETENTION` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_LIMIT_SUMMARY` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_REPORTING_PERIOD` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_ACCUMULATION` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_PORTFOLIO` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_HAZARD_LIMIT` | 1:1 | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_VIEW_COMMENT` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_FAC_LIMIT_DETAIL` | 1:N | *tidak dinyatakan* | `FAC_LIMIT_ID` | CASCADE | `T_TREATY_FAC_LIMITS` | EDM Prop, EDM Non Prop | **ya** |
| `T_TREATY_VALUE_DIFFERENCE` | 1:1 | *tidak dinyatakan* | `REVISION_ID` | CASCADE | `T_TREATY_REVISION` | EDM Non Prop | **ya** |
| `T_TREATY_LIMIT_MEASURE` | 1:N | *tidak dinyatakan* | `LIMIT_ID` | CASCADE | `T_TREATY_LIMITS` | EDM Non Prop | **ya** |
| `T_TREATY_REINSTATEMENT` | 1:N | *tidak dinyatakan* | `LIMIT_ID` | CASCADE | `T_TREATY_LIMITS` | EDM Non Prop | **ya** |
| `T_TREATY_FAC_REINSURER` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Non Prop | **ya** |
| `T_TREATY_CURRENCY` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Non Prop | **ya** |
| `T_TREATY_EGNPI` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Non Prop | **ya** |
| `T_TREATY_INSTALLMENT` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Non Prop | **ya** |
| `T_TREATY_TOTAL` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Non Prop | **ya** |
| `T_TREATY_RETRO_SHARE` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Non Prop | **ya** |
| `T_TREATY_VALUE_BEFORE_PRORATE` | 1:1 | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Non Prop | **TIDAK — khas modul ini** |

> **Kolom *Kunci utama* berbunyi *tidak dinyatakan* untuk 34 dari 35 tabel, dan itu salinan setia.** Berkas sumbernya mencetak **PK hanya pada kotak akar**; pada kotak anak baris keduanya dipakai untuk **FK**. Maka kekosongan ini berarti **sumbernya diam**, bukan **tabelnya tanpa kunci utama**.

## 3. Asal di sistem lama dan padanannya di model baru

| Tabel | Jalur Pega | Padanan di model baru | Tingkat bukti |
|---|---|---|---|
| `TREATY_IN` | `TreatyIn` | **KONTRAK + VERSI_KONTRAK** | JALUR-PENUH |
| `T_TREATY_REVISION` | `TreatyIn` | **VERSI_KONTRAK · bagian addendum** | JALUR-PENUH |
| `T_TREATY_LIMITS` | `TreatyIn.Limits` | **LAYER** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_LIMIT_DETAIL` | `TreatyIn.Limits.Detail` | **DETAIL_PROPORSIONAL** | DAUN-RELATIF |
| `T_TREATY_LIMIT_COB` | `TreatyIn.Limits.Detail.COBList` | **KELAS_BISNIS_LAYER** | DAUN-RELATIF |
| `T_TREATY_LIMIT_AMOUNT` | `TreatyIn.Limits.Detail.EPIList | TreatyIn.Limits.Detail.RetentionList | Tr` | **NILAI_LAYER** | DAUN-RELATIF |
| `T_TREATY_LIMIT_ACHIEVEMENT` | `TreatyIn.Limits.Detail.AchievementLists` | **PENCAPAIAN** | DAUN-RELATIF |
| `T_TREATY_LIMIT_GROUP` | `TreatyIn.Limits.TreatyGroupList` | **KELOMPOK_LAYER** | DAUN-RELATIF |
| `T_TREATY_LIMIT_GROUP_COB` | `TreatyIn.Limits.TreatyGroupList.ClassOfBusinessList` | **KELAS_BISNIS_KELOMPOK** | DAUN-RELATIF |
| `T_TREATY_SHARE` | `TreatyIn.Share` | **BAGIAN** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_SHARE_SPREADING` | `TreatyIn.Share.SpreadingListXOL` | **PENYEBARAN_XOL** | DAUN-RELATIF |
| `T_TREATY_SHARE_SPREAD_AMOUNT` | `TreatyIn.Share.SpreadingListXOL.RNMSpreadedListXOL | TreatyIn.Share.Spread` | **NILAI_TERSEBAR** | DAUN-RELATIF |
| `T_TREATY_SHARE_AMOUNT` | `TreatyIn.Share.GrossPremiumList | TreatyIn.Share.NetPremiumList | TreatyIn` | **NILAI_BAGIAN** | DAUN-RELATIF |
| `T_TREATY_SHARE_DEDUCTION` | `TreatyIn.Share.DeductionList` | **POTONGAN** | DAUN-RELATIF |
| `T_TREATY_FAC_SHARE` | `TreatyIn.FacultativeShareList` | **BAGIAN_FAKULTATIF** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_FAC_SHARE_AMOUNT` | `TreatyIn.FacultativeShareList.GrossPremiumList | TreatyIn.FacultativeShare` | **NILAI_BAGIAN_FAKULTATIF** | DAUN-RELATIF |
| `T_TREATY_FAC_SHARE_DEDUCTION` | `TreatyIn.FacultativeShareList.DeductionList` | **POTONGAN_FAKULTATIF** | DAUN-RELATIF |
| `T_TREATY_RETENTION` | `TreatyIn.Retention` | **RETENSI** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_LIMIT_SUMMARY` | `TreatyIn.LimitSummaryList | TreatyIn.LimitShareSummaryList | TreatyIn.Limi` | **RINGKASAN_LIMIT** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_REPORTING_PERIOD` | `TreatyIn.ReportingPeriodList` | **PERIODE_PELAPORAN** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_ACCUMULATION` | `TreatyIn.AccumulationList` | **AKUMULASI** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_PORTFOLIO` | `TreatyIn.Portfolio` | **PORTOFOLIO** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_HAZARD_LIMIT` | `TreatyIn` | **BATAS_BAHAYA** | JALUR-PENUH |
| `T_VIEW_COMMENT` | `TreatyIn.CommentList` | **CATATAN_PERSETUJUAN** | JALUR-PENUH |
| `T_TREATY_FAC_LIMIT_DETAIL` | `TreatyIn.ShareFacultativeReinsurers.FacultativeLimits.Detail` | **DETAIL_PROPORSIONAL_FAKULTATIF** | DAUN-RELATIF |
| `T_TREATY_VALUE_DIFFERENCE` | `TreatyIn.ValueDifference` | **NILAI_SELISIH** | JALUR-PENUH |
| `T_TREATY_LIMIT_MEASURE` | `TreatyIn.Limits.MDPList | TreatyIn.Limits.EgnpiTotalList | TreatyIn.Limits` | **BESARAN_LAYER** | DAUN-RELATIF |
| `T_TREATY_REINSTATEMENT` | `TreatyIn.Limits.Reinstatement_List` | **PEMULIHAN_LIMIT** | DAUN-RELATIF |
| `T_TREATY_FAC_REINSURER` | `TreatyIn.ShareFacultativeReinsurers | TreatyIn.FacultativeShareList.ShareF` | **REASURADUR_FAKULTATIF** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_CURRENCY` | `TreatyIn.CurrencyList` | **MATA_UANG_KONTRAK** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_EGNPI` | `TreatyIn.EGNPI` | **EGNPI** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_INSTALLMENT` | `TreatyIn.Installment` | **ANGSURAN** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_TOTAL` | `TreatyIn.TotalShareNetNP | TreatyIn.TotalShareGrossNP | TreatyIn.TotalEgnp` | **REKAP_KONTRAK** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_RETRO_SHARE` | `TreatyIn.ShareReins` | **BAGIAN_RETRO** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_VALUE_BEFORE_PRORATE` | `TreatyIn.ValueBeforeProrate` | **NILAI_SEBELUM_PRO_RATE** | JALUR-PENUH |

> **Kolom *Padanan di model baru* adalah BACAAN, bukan pemetaan yang mengikat.** Yang mengikat soal nama kolom dan tipe: `treaty-in/2-to-spec/KAMUS-KOLOM.md`; soal daftar entitas: `treaty-in/4-erd-dan-tabel-datar/STRUKTUR-DATA.md`.

## 4. Relasi yang menyentuh modul ini — 34

| # | Induk | Kard. | Anak | Kunci tamu | ON DELETE | Muncul di lembar |
|---|---|---|---|---|---|---|
| 1 | `TREATY_IN` | 1:1 | `T_TREATY_REVISION` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 2 | `TREATY_IN` | 1:N | `T_TREATY_LIMITS` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 3 | `T_TREATY_LIMITS` | 1:N | `T_TREATY_LIMIT_DETAIL` | `LIMIT_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 4 | `T_TREATY_LIMIT_DETAIL` | 1:N | `T_TREATY_LIMIT_COB` | `LIMIT_DETAIL_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 5 | `T_TREATY_LIMIT_DETAIL` | 1:N | `T_TREATY_LIMIT_AMOUNT` | `LIMIT_DETAIL_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 6 | `T_TREATY_LIMIT_DETAIL` | 1:N | `T_TREATY_LIMIT_ACHIEVEMENT` | `LIMIT_DETAIL_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 7 | `T_TREATY_LIMITS` | 1:N | `T_TREATY_LIMIT_MEASURE` | `LIMIT_ID` | CASCADE | Non Prop, EDM Non Prop |
| 8 | `T_TREATY_LIMITS` | 1:N | `T_TREATY_REINSTATEMENT` | `LIMIT_ID` | CASCADE | Non Prop, EDM Non Prop |
| 9 | `T_TREATY_LIMITS` | 1:N | `T_TREATY_LIMIT_GROUP` | `LIMIT_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 10 | `T_TREATY_LIMIT_GROUP` | 1:N | `T_TREATY_LIMIT_GROUP_COB` | `LIMIT_GROUP_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 11 | `TREATY_IN` | 1:N | `T_TREATY_SHARE` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 12 | `T_TREATY_SHARE` | 1:N | `T_TREATY_SHARE_SPREADING` | `SHARE_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 13 | `T_TREATY_SHARE_SPREADING` | 1:N | `T_TREATY_SHARE_SPREAD_AMOUNT` | `SHARE_SPREADING_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 14 | `T_TREATY_SHARE` | 1:N | `T_TREATY_SHARE_AMOUNT` | `SHARE_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 15 | `T_TREATY_SHARE` | 1:N | `T_TREATY_SHARE_DEDUCTION` | `SHARE_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 16 | `TREATY_IN` | 1:N | `T_TREATY_FAC_SHARE` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 17 | `T_TREATY_FAC_SHARE` | 1:N | `T_TREATY_FAC_SHARE_AMOUNT` | `FAC_SHARE_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 18 | `T_TREATY_FAC_SHARE` | 1:N | `T_TREATY_FAC_SHARE_DEDUCTION` | `FAC_SHARE_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 19 | `TREATY_IN` | 1:N | `T_TREATY_FAC_REINSURER` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Non Prop |
| 21 | `T_TREATY_FAC_LIMITS` | 1:N | `T_TREATY_FAC_LIMIT_DETAIL` | `FAC_LIMIT_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 22 | `TREATY_IN` | 1:N | `T_TREATY_CURRENCY` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Non Prop |
| 23 | `TREATY_IN` | 1:N | `T_TREATY_EGNPI` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Non Prop |
| 24 | `TREATY_IN` | 1:N | `T_TREATY_RETENTION` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 25 | `TREATY_IN` | 1:N | `T_TREATY_LIMIT_SUMMARY` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 26 | `TREATY_IN` | 1:N | `T_TREATY_INSTALLMENT` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Non Prop |
| 28 | `TREATY_IN` | 1:N | `T_TREATY_REPORTING_PERIOD` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 29 | `TREATY_IN` | 1:N | `T_TREATY_ACCUMULATION` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 30 | `TREATY_IN` | 1:N | `T_TREATY_PORTFOLIO` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 31 | `TREATY_IN` | 1:1 | `T_TREATY_HAZARD_LIMIT` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 32 | `TREATY_IN` | 1:N | `T_TREATY_TOTAL` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Non Prop |
| 33 | `TREATY_IN` | 1:N | `T_TREATY_RETRO_SHARE` | `TREATY_IN_ID` | CASCADE | Prop, Non Prop, EDM Non Prop |
| 34 | `TREATY_IN` | 1:N | `T_VIEW_COMMENT` | `TREATY_IN_ID` | di Go | Prop, Non Prop, EDM Prop, EDM Non Prop |
| 36 | `T_TREATY_REVISION` | 1:1 | `T_TREATY_VALUE_DIFFERENCE` | `REVISION_ID` | CASCADE | Prop, Non Prop, EDM Non Prop |
| 37 | `TREATY_IN` | 1:1 | `T_TREATY_VALUE_BEFORE_PRORATE` | `TREATY_IN_ID` | CASCADE | EDM Non Prop |

### 4a. Relasi yang TIDAK muncul di lembar mana pun — 3

Ketiganya ada di *Daftar Relasi* sumbernya, dan kolom **Muncul di lembar**-nya berbunyi
*(tidak terklasifikasi)*. **Dibawa ke KEDUA berkas**, sebab membuangnya dari keduanya
akan menghilangkannya sama sekali — dan memilih salah satu berarti memutuskan sesuatu
yang sumbernya tidak nyatakan.

| # | Induk | Kard. | Anak | Kunci tamu | ON DELETE |
|---|---|---|---|---|---|
| 35 | `TREATY_IN` | — | `DOCUMENT_TREATY_IN` | `TREATY_IN_ID` | di Go |
| 38 | `T_TREATY_MIG_LANDING` | — | `T_TREATY_MIG_REJECTED` | `LANDING_ID` | di Go |
| 39 | `TREATY_IN` | — | `T_TREATY_OUTBOUND_ARCHIVE` | `TREATY_IN_ID` | di Go |

> **4 tabel di antaranya tidak pernah digambar sebagai kotak** di lembar mana pun: `DOCUMENT_TREATY_IN`, `T_TREATY_MIG_LANDING`, `T_TREATY_MIG_REJECTED`, `T_TREATY_OUTBOUND_ARCHIVE`. Ia ada di daftar relasi dan tidak di gambar — dan itu keadaan sumbernya, bukan kekeliruan perkakas ini.

## 5. Pohon per lembar

### Treaty In EDM Prop — 25 kotak

- **`TREATY_IN`** — PK `ID`
  - **`T_TREATY_REVISION`** *1:1* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_LIMITS`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_LIMIT_DETAIL`** *1:N* — FK `LIMIT_ID` → `T_TREATY_LIMITS`.`ID` · CASCADE
      - **`T_TREATY_LIMIT_COB`** *1:N* — FK `LIMIT_DETAIL_ID` → `T_TREATY_LIMIT_DETAIL`.`ID` · CASCADE
      - **`T_TREATY_LIMIT_AMOUNT`** *1:N* — FK `LIMIT_DETAIL_ID` → `T_TREATY_LIMIT_DETAIL`.`ID` · CASCADE
      - **`T_TREATY_LIMIT_ACHIEVEMENT`** *1:N* — FK `LIMIT_DETAIL_ID` → `T_TREATY_LIMIT_DETAIL`.`ID` · CASCADE
    - **`T_TREATY_LIMIT_GROUP`** *1:N* — FK `LIMIT_ID` → `T_TREATY_LIMITS`.`ID` · CASCADE
      - **`T_TREATY_LIMIT_GROUP_COB`** *1:N* — FK `LIMIT_GROUP_ID` → `T_TREATY_LIMIT_GROUP`.`ID` · CASCADE
  - **`T_TREATY_SHARE`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_SHARE_SPREADING`** *1:N* — FK `SHARE_ID` → `T_TREATY_SHARE`.`ID` · CASCADE
      - **`T_TREATY_SHARE_SPREAD_AMOUNT`** *1:N* — FK `SHARE_SPREADING_ID` → `T_TREATY_SHARE_SPREADING`.`ID` · CASCADE
    - **`T_TREATY_SHARE_AMOUNT`** *1:N* — FK `SHARE_ID` → `T_TREATY_SHARE`.`ID` · CASCADE
    - **`T_TREATY_SHARE_DEDUCTION`** *1:N* — FK `SHARE_ID` → `T_TREATY_SHARE`.`ID` · CASCADE
  - **`T_TREATY_FAC_SHARE`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_FAC_SHARE_AMOUNT`** *1:N* — FK `FAC_SHARE_ID` → `T_TREATY_FAC_SHARE`.`ID` · CASCADE
    - **`T_TREATY_FAC_SHARE_DEDUCTION`** *1:N* — FK `FAC_SHARE_ID` → `T_TREATY_FAC_SHARE`.`ID` · CASCADE
  - **`T_TREATY_RETENTION`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_LIMIT_SUMMARY`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_REPORTING_PERIOD`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_ACCUMULATION`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_PORTFOLIO`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_HAZARD_LIMIT`** *1:1* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_VIEW_COMMENT`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
- **`T_TREATY_FAC_LIMIT_DETAIL`** *1:N* — FK `FAC_LIMIT_ID` → `T_TREATY_FAC_LIMITS`.`ID` · CASCADE

### Treaty In EDM Non Prop — 35 kotak

- **`TREATY_IN`** — PK `ID`
  - **`T_TREATY_REVISION`** *1:1* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_VALUE_DIFFERENCE`** *1:1* — FK `REVISION_ID` → `T_TREATY_REVISION`.`ID` · CASCADE
  - **`T_TREATY_LIMITS`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_LIMIT_DETAIL`** *1:N* — FK `LIMIT_ID` → `T_TREATY_LIMITS`.`ID` · CASCADE
      - **`T_TREATY_LIMIT_COB`** *1:N* — FK `LIMIT_DETAIL_ID` → `T_TREATY_LIMIT_DETAIL`.`ID` · CASCADE
      - **`T_TREATY_LIMIT_AMOUNT`** *1:N* — FK `LIMIT_DETAIL_ID` → `T_TREATY_LIMIT_DETAIL`.`ID` · CASCADE
      - **`T_TREATY_LIMIT_ACHIEVEMENT`** *1:N* — FK `LIMIT_DETAIL_ID` → `T_TREATY_LIMIT_DETAIL`.`ID` · CASCADE
    - **`T_TREATY_LIMIT_MEASURE`** *1:N* — FK `LIMIT_ID` → `T_TREATY_LIMITS`.`ID` · CASCADE
    - **`T_TREATY_REINSTATEMENT`** *1:N* — FK `LIMIT_ID` → `T_TREATY_LIMITS`.`ID` · CASCADE
    - **`T_TREATY_LIMIT_GROUP`** *1:N* — FK `LIMIT_ID` → `T_TREATY_LIMITS`.`ID` · CASCADE
      - **`T_TREATY_LIMIT_GROUP_COB`** *1:N* — FK `LIMIT_GROUP_ID` → `T_TREATY_LIMIT_GROUP`.`ID` · CASCADE
  - **`T_TREATY_SHARE`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_SHARE_SPREADING`** *1:N* — FK `SHARE_ID` → `T_TREATY_SHARE`.`ID` · CASCADE
      - **`T_TREATY_SHARE_SPREAD_AMOUNT`** *1:N* — FK `SHARE_SPREADING_ID` → `T_TREATY_SHARE_SPREADING`.`ID` · CASCADE
    - **`T_TREATY_SHARE_AMOUNT`** *1:N* — FK `SHARE_ID` → `T_TREATY_SHARE`.`ID` · CASCADE
    - **`T_TREATY_SHARE_DEDUCTION`** *1:N* — FK `SHARE_ID` → `T_TREATY_SHARE`.`ID` · CASCADE
  - **`T_TREATY_FAC_SHARE`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_FAC_SHARE_AMOUNT`** *1:N* — FK `FAC_SHARE_ID` → `T_TREATY_FAC_SHARE`.`ID` · CASCADE
    - **`T_TREATY_FAC_SHARE_DEDUCTION`** *1:N* — FK `FAC_SHARE_ID` → `T_TREATY_FAC_SHARE`.`ID` · CASCADE
  - **`T_TREATY_FAC_REINSURER`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_CURRENCY`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_EGNPI`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_RETENTION`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_LIMIT_SUMMARY`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_INSTALLMENT`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_REPORTING_PERIOD`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_ACCUMULATION`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_PORTFOLIO`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_HAZARD_LIMIT`** *1:1* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_TOTAL`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_RETRO_SHARE`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_VIEW_COMMENT`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_VALUE_BEFORE_PRORATE`** *1:1* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
- **`T_TREATY_FAC_LIMIT_DETAIL`** *1:N* — FK `FAC_LIMIT_ID` → `T_TREATY_FAC_LIMITS`.`ID` · CASCADE

## 6. Catatan & Batas — disalin apa adanya dari sumbernya

| | |
|---|---|
| **CATATAN DAN BATAS — dibaca sebelum memakai angka mana pun** |  |
| **Sumber daftar tabel** | peta-nama-tabel-treatyin.tsv — 46 baris, 2 bertanda TIDAK DIBUAT dikeluarkan, sisa 44. |
| **Sumber penggolongan** | Sapuan atas 118 berkas Section: "Treaty In" 50 berkas, "Treaty In Adjustment" 68 berkas. |
| **Salinan tidak dihitung** | 733 blok pyIncludedRuleXML dibuang sebelum pencocokan (MA-04: seksi tersertakan bukan bukti pemakaian). |
| **Cara menggolongkan** | Dari NAMA SEKSI yang merujuk tabel itu. Proportional -> PROP. NonProportional / Layers / XOL / NP -> NON-PROP. OldData / EDM / ValueDifference / Adjust / Revisi / Picker -> EDM. Selain itu UMUM. |
| **Arti UMUM** | Seksi yang melayani kedua jenis treaty. Tabel yang hanya muncul di seksi UMUM masuk ke lembar Prop dan Non Prop sekaligus, lalu ditandai BERSAMA. |
| **YANG DAPAT DIBUKTIKAN** | Bahwa sebuah tabel DIRUJUK oleh seksi tertentu. |
| **YANG TIDAK DAPAT DIBUKTIKAN** | Bahwa sebuah tabel TIDAK DIPAKAI. Activity, Data Transform, Report Definition, dan RDB List tidak disapu. Nol berarti BELUM TERPERIKSA. |
| **Dua tingkat bukti** | JALUR-PENUH 20 tabel (kuat). DAUN-RELATIF 17 tabel (lemah). |
| **Tidak terklasifikasi** | 7 tabel tanpa jalur Pega: DOCUMENT_TREATY_IN, T_TREATY_MIG_CORRELATION, T_TREATY_MIG_LANDING, T_TREATY_MIG_REJECTED, T_TREATY_OUTBOUND_ARCHIVE, T_TREATY_OFFER, T_TREATY_AUTHORITY_LIMIT. Tidak digambar — dilaporkan di sini. |
| **Temuan yang layak dibaca** | Di tingkat TABEL, prop dan non-prop berbagi hampir seluruh entitas. Pemisah sesungguhnya ada di KOLOM, dan kolom sengaja tidak dimuat sesuai permintaan. |
| **Sifat berkas ini** | Potret dari artefak to-spec yang masih berjalan. peta-nama-tabel-treatyin.tsv belum memuat hasil grilling Adjustment GRL-12 sampai GRL-18, dan daftar entitas induk masih bertengkar antara 27 dan 28. |
| **Dibuat** | Sapuan mekanis 24 September 2026, mengikuti format contooh.xlsx. |

## 7. Batas berkas ini

| Ia TIDAK dapat menyatakan… | Sebabnya |
|---|---|
| bahwa himpunan tabelnya **lengkap** | ia salinan sebuah **gambar**. Gambar itu disusun dari pohon clipboard Pega, dan pohon itu **kurang 340 properti titik buta** (`L-8`, ditagih `M-4`) |
| bahwa pembagian dua modul ini **benar secara bisnis** | sumbu pemisahnya **nama lembar**, dan nama lembar adalah keputusan orang yang menggambar — bukan fakta yang dibaca dari ekspor |
| bahwa `ON DELETE` di sini **perilaku sistem lama** | Catatan & Batas sumbernya sendiri menyatakannya **usulan rancangan**, bukan perilaku yang terbaca |

