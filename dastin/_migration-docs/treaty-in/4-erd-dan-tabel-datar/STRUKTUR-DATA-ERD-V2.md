# Struktur data — **dibangkitkan dari `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`**

> ### POTRET SISTEM LAMA, 24 September 2026. BUKAN RANCANGAN.
>
> Skema yang dibangun ada di 2-to-spec/KAMUS-KOLOM.md dan 2-to-spec/ddl-usulan/. Gambar skema barunya: 4-erd-dan-tabel-datar/ERD-SKEMA-BARU.xlsx

**Dibangkitkan:** [`alat/bangkitkan-dari-erd-v2.py`](../alat/bangkitkan-dari-erd-v2.py) · **Sumber tunggal:** `alat/erd-v2.json`

> **Berkas ini tidak ditulis tangan, dan tidak boleh disunting.** Ia salinan struktur yang ada di berkas sumbernya — nama tabel, kardinalitas, kunci, jalur Pega, padanan model baru, dan tingkat bukti seluruhnya apa adanya. Sunting sumbernya, lalu jalankan ulang alatnya.

> ### BERKAS INI GABUNGAN DUA MODUL — pecahannya ada, dan dibangkitkan dari sumber yang SAMA
>
> | Modul | Berkas | Lembar |
> |---|---|---|
> | **Treaty In** | [`../STRUKTUR-DATA-TREATY-IN.md`](../STRUKTUR-DATA-TREATY-IN.md) | *Treaty In Prop* · *Treaty In Non Prop* |
> | **Treaty In Adjustment** | [`../../treaty-in-adjustment/STRUKTUR-DATA-TREATY-IN-ADJUSTMENT.md`](../../treaty-in-adjustment/STRUKTUR-DATA-TREATY-IN-ADJUSTMENT.md) | *Treaty In EDM Prop* · *Treaty In EDM Non Prop* |
>
> Keduanya dibangkitkan [`alat/pisah-struktur-erd-v2-per-modul.py`](../alat/pisah-struktur-erd-v2-per-modul.py) dari `alat/erd-v2.json` yang sama. **34 dari 37 tabel dipakai kedua modul** dan karena itu muncul di kedua pecahan — itu kenyataan sistem lama, bukan penggandaan.

---

## 1. Hitungan

| | Jumlah |
|---|---:|
| Tabel unik di keempat lembar | **37** |
| Kotak tergambar (satu tabel dapat muncul di beberapa lembar) | **130** |
| Relasi di lembar *Daftar Relasi* | **41** |
| Tabel yang dipakai **bersama** lebih dari satu lembar | **36** |
| Lembar *Treaty In Prop* | 34 kotak |
| Lembar *Treaty In Non Prop* | 36 kotak |
| Lembar *Treaty In EDM Prop* | 25 kotak |
| Lembar *Treaty In EDM Non Prop* | 35 kotak |

## 2. Entitas — satu baris per tabel unik

| Tabel | Kard. | Kunci utama | Kunci asing | ON DELETE | Induk | Muncul di lembar |
|---|---|---|---|---|---|---|
| `TREATY_IN` | — | `ID` | — | — | **akar** | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_REVISION` | 1:1 | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_VALUE_DIFFERENCE` | 1:1 | *tidak dinyatakan* | `REVISION_ID` | CASCADE | `T_TREATY_REVISION` | Prop, Non Prop, EDM Non Prop |
| `T_TREATY_LIMITS` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_LIMIT_DETAIL` | 1:N | *tidak dinyatakan* | `LIMIT_ID` | CASCADE | `T_TREATY_LIMITS` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_LIMIT_COB` | 1:N | *tidak dinyatakan* | `LIMIT_DETAIL_ID` | CASCADE | `T_TREATY_LIMIT_DETAIL` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_LIMIT_AMOUNT` | 1:N | *tidak dinyatakan* | `LIMIT_DETAIL_ID` | CASCADE | `T_TREATY_LIMIT_DETAIL` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_LIMIT_ACHIEVEMENT` | 1:N | *tidak dinyatakan* | `LIMIT_DETAIL_ID` | CASCADE | `T_TREATY_LIMIT_DETAIL` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_LIMIT_GROUP` | 1:N | *tidak dinyatakan* | `LIMIT_ID` | CASCADE | `T_TREATY_LIMITS` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_LIMIT_GROUP_COB` | 1:N | *tidak dinyatakan* | `LIMIT_GROUP_ID` | CASCADE | `T_TREATY_LIMIT_GROUP` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_SHARE` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_SHARE_SPREADING` | 1:N | *tidak dinyatakan* | `SHARE_ID` | CASCADE | `T_TREATY_SHARE` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_SHARE_SPREAD_AMOUNT` | 1:N | *tidak dinyatakan* | `SHARE_SPREADING_ID` | CASCADE | `T_TREATY_SHARE_SPREADING` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_SHARE_AMOUNT` | 1:N | *tidak dinyatakan* | `SHARE_ID` | CASCADE | `T_TREATY_SHARE` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_SHARE_DEDUCTION` | 1:N | *tidak dinyatakan* | `SHARE_ID` | CASCADE | `T_TREATY_SHARE` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_FAC_SHARE` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_FAC_SHARE_AMOUNT` | 1:N | *tidak dinyatakan* | `FAC_SHARE_ID` | CASCADE | `T_TREATY_FAC_SHARE` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_FAC_SHARE_DEDUCTION` | 1:N | *tidak dinyatakan* | `FAC_SHARE_ID` | CASCADE | `T_TREATY_FAC_SHARE` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_FAC_REINSURER` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Non Prop |
| `T_TREATY_FAC_LIMITS` | 1:N | *tidak dinyatakan* | `FAC_REINSURER_ID` | CASCADE | `T_TREATY_FAC_REINSURER` | Prop, Non Prop |
| `T_TREATY_FAC_LIMIT_DETAIL` | 1:N | *tidak dinyatakan* | `FAC_LIMIT_ID` | CASCADE | `T_TREATY_FAC_LIMITS` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_CURRENCY` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Non Prop |
| `T_TREATY_EGNPI` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Non Prop |
| `T_TREATY_RETENTION` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_LIMIT_SUMMARY` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_INSTALLMENT` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Non Prop |
| `T_TREATY_INSTALLMENT_ITEM` | 1:N | *tidak dinyatakan* | `INSTALLMENT_ID` | CASCADE | `T_TREATY_INSTALLMENT` | Prop, Non Prop |
| `T_TREATY_REPORTING_PERIOD` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_ACCUMULATION` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_PORTFOLIO` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_HAZARD_LIMIT` | 1:1 | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_TOTAL` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Non Prop |
| `T_TREATY_RETRO_SHARE` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Non Prop |
| `T_VIEW_COMMENT` | 1:N | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | Prop, Non Prop, EDM Prop, EDM Non Prop |
| `T_TREATY_LIMIT_MEASURE` | 1:N | *tidak dinyatakan* | `LIMIT_ID` | CASCADE | `T_TREATY_LIMITS` | Non Prop, EDM Non Prop |
| `T_TREATY_REINSTATEMENT` | 1:N | *tidak dinyatakan* | `LIMIT_ID` | CASCADE | `T_TREATY_LIMITS` | Non Prop, EDM Non Prop |
| `T_TREATY_VALUE_BEFORE_PRORATE` | 1:1 | *tidak dinyatakan* | `TREATY_IN_ID` | CASCADE | `TREATY_IN` | EDM Non Prop |

> **Kolom *Kunci utama* berbunyi *tidak dinyatakan* untuk 36 dari 37 tabel, dan itu salinan setia.** Berkas sumbernya mencetak **PK hanya pada kotak akar**; pada kotak anak baris keduanya dipakai untuk **FK**. Maka kekosongan ini berarti **sumbernya diam**, bukan **tabelnya tanpa kunci utama** — lembar *Daftar Relasi* sumbernya sendiri merujuk `‹TABEL›.ID` sebagai sasaran tiap kunci asing.

## 3. Asal di sistem lama dan padanannya di model baru

| Tabel | Jalur Pega | Padanan di model baru | Tingkat bukti |
|---|---|---|---|
| `TREATY_IN` | `TreatyIn` | **KONTRAK + VERSI_KONTRAK** | JALUR-PENUH |
| `T_TREATY_REVISION` | `TreatyIn` | **VERSI_KONTRAK · bagian addendum** | JALUR-PENUH |
| `T_TREATY_VALUE_DIFFERENCE` | `TreatyIn.ValueDifference` | **NILAI_SELISIH** | JALUR-PENUH |
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
| `T_TREATY_FAC_REINSURER` | `TreatyIn.ShareFacultativeReinsurers | TreatyIn.FacultativeShareList.ShareF` | **REASURADUR_FAKULTATIF** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_FAC_LIMITS` | `TreatyIn.ShareFacultativeReinsurers.FacultativeLimits` | **LAYER_FAKULTATIF** | DAUN-RELATIF |
| `T_TREATY_FAC_LIMIT_DETAIL` | `TreatyIn.ShareFacultativeReinsurers.FacultativeLimits.Detail` | **DETAIL_PROPORSIONAL_FAKULTATIF** | DAUN-RELATIF |
| `T_TREATY_CURRENCY` | `TreatyIn.CurrencyList` | **MATA_UANG_KONTRAK** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_EGNPI` | `TreatyIn.EGNPI` | **EGNPI** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_RETENTION` | `TreatyIn.Retention` | **RETENSI** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_LIMIT_SUMMARY` | `TreatyIn.LimitSummaryList | TreatyIn.LimitShareSummaryList | TreatyIn.Limi` | **RINGKASAN_LIMIT** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_INSTALLMENT` | `TreatyIn.Installment` | **ANGSURAN** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_INSTALLMENT_ITEM` | `TreatyIn.ValueDifference.Installment.InstallmentList` | **RINCIAN_ANGSURAN** | DAUN-RELATIF |
| `T_TREATY_REPORTING_PERIOD` | `TreatyIn.ReportingPeriodList` | **PERIODE_PELAPORAN** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_ACCUMULATION` | `TreatyIn.AccumulationList` | **AKUMULASI** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_PORTFOLIO` | `TreatyIn.Portfolio` | **PORTOFOLIO** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_HAZARD_LIMIT` | `TreatyIn` | **BATAS_BAHAYA** | JALUR-PENUH |
| `T_TREATY_TOTAL` | `TreatyIn.TotalShareNetNP | TreatyIn.TotalShareGrossNP | TreatyIn.TotalEgnp` | **REKAP_KONTRAK** | DAUN-RELATIF/JALUR-PENUH |
| `T_TREATY_RETRO_SHARE` | `TreatyIn.ShareReins` | **BAGIAN_RETRO** | DAUN-RELATIF/JALUR-PENUH |
| `T_VIEW_COMMENT` | `TreatyIn.CommentList` | **CATATAN_PERSETUJUAN** | JALUR-PENUH |
| `T_TREATY_LIMIT_MEASURE` | `TreatyIn.Limits.MDPList | TreatyIn.Limits.EgnpiTotalList | TreatyIn.Limits` | **BESARAN_LAYER** | DAUN-RELATIF |
| `T_TREATY_REINSTATEMENT` | `TreatyIn.Limits.Reinstatement_List` | **PEMULIHAN_LIMIT** | DAUN-RELATIF |
| `T_TREATY_VALUE_BEFORE_PRORATE` | `TreatyIn.ValueBeforeProrate` | **NILAI_SEBELUM_PRO_RATE** | JALUR-PENUH |

## 4. Pohon per lembar

### Treaty In Prop — 34 kotak

- **`TREATY_IN`** — PK `ID`
  - **`T_TREATY_REVISION`** *1:1* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_VALUE_DIFFERENCE`** *1:1* — FK `REVISION_ID` → `T_TREATY_REVISION`.`ID` · CASCADE
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
  - **`T_TREATY_FAC_REINSURER`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_FAC_LIMITS`** *1:N* — FK `FAC_REINSURER_ID` → `T_TREATY_FAC_REINSURER`.`ID` · CASCADE
      - **`T_TREATY_FAC_LIMIT_DETAIL`** *1:N* — FK `FAC_LIMIT_ID` → `T_TREATY_FAC_LIMITS`.`ID` · CASCADE
  - **`T_TREATY_CURRENCY`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_EGNPI`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_RETENTION`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_LIMIT_SUMMARY`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_INSTALLMENT`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_INSTALLMENT_ITEM`** *1:N* — FK `INSTALLMENT_ID` → `T_TREATY_INSTALLMENT`.`ID` · CASCADE
  - **`T_TREATY_REPORTING_PERIOD`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_ACCUMULATION`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_PORTFOLIO`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_HAZARD_LIMIT`** *1:1* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_TOTAL`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_RETRO_SHARE`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_VIEW_COMMENT`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE

### Treaty In Non Prop — 36 kotak

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
    - **`T_TREATY_FAC_LIMITS`** *1:N* — FK `FAC_REINSURER_ID` → `T_TREATY_FAC_REINSURER`.`ID` · CASCADE
      - **`T_TREATY_FAC_LIMIT_DETAIL`** *1:N* — FK `FAC_LIMIT_ID` → `T_TREATY_FAC_LIMITS`.`ID` · CASCADE
  - **`T_TREATY_CURRENCY`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_EGNPI`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_RETENTION`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_LIMIT_SUMMARY`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_INSTALLMENT`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
    - **`T_TREATY_INSTALLMENT_ITEM`** *1:N* — FK `INSTALLMENT_ID` → `T_TREATY_INSTALLMENT`.`ID` · CASCADE
  - **`T_TREATY_REPORTING_PERIOD`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_ACCUMULATION`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_PORTFOLIO`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_HAZARD_LIMIT`** *1:1* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_TOTAL`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_TREATY_RETRO_SHARE`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE
  - **`T_VIEW_COMMENT`** *1:N* — FK `TREATY_IN_ID` → `TREATY_IN`.`ID` · CASCADE

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

## 5. Catatan & Batas — disalin apa adanya dari sumbernya

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

