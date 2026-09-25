# Tabel datar — **dibangkitkan dari `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`**

> ### POTRET SISTEM LAMA, 24 September 2026. BUKAN RANCANGAN.
>
> Skema yang dibangun ada di 2-to-spec/KAMUS-KOLOM.md dan 2-to-spec/ddl-usulan/. Gambar skema barunya: 4-erd-dan-tabel-datar/ERD-SKEMA-BARU.xlsx

**Dibangkitkan:** [`alat/bangkitkan-dari-erd-v2.py`](../alat/bangkitkan-dari-erd-v2.py) · **Sumber tunggal:** `alat/erd-v2.json` — sama dengan [`STRUKTUR-DATA-ERD-V2.md`](STRUKTUR-DATA-ERD-V2.md) dan [`ERD-TREATY-IN-DAN-EDM.html`](ERD-TREATY-IN-DAN-EDM.html)

> **Ketiganya tidak dapat berbeda isinya**, sebab ketiganya dibangkitkan dari satu berkas yang sama. Bila berbeda, **alatnya yang salah**.

---

## 1. Satu baris per tabel datar

| # | Tabel datar | Satu baris mewakili | Kunci | Induk | Kunci tamu | ON DELETE |
|---:|---|---|---|---|---|---|
| 1 | `TREATY_IN` | `TreatyIn` | `ID` | **akar** | — | — |
| 2 | `T_TREATY_REVISION` | `TreatyIn` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 3 | `T_TREATY_VALUE_DIFFERENCE` | `TreatyIn.ValueDifference` | *tidak dinyatakan* | `T_TREATY_REVISION` | `REVISION_ID` | CASCADE |
| 4 | `T_TREATY_LIMITS` | `TreatyIn.Limits` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 5 | `T_TREATY_LIMIT_DETAIL` | `TreatyIn.Limits.Detail` | *tidak dinyatakan* | `T_TREATY_LIMITS` | `LIMIT_ID` | CASCADE |
| 6 | `T_TREATY_LIMIT_COB` | `TreatyIn.Limits.Detail.COBList` | *tidak dinyatakan* | `T_TREATY_LIMIT_DETAIL` | `LIMIT_DETAIL_ID` | CASCADE |
| 7 | `T_TREATY_LIMIT_AMOUNT` | `TreatyIn.Limits.Detail.EPIList | TreatyIn.Limits.Detail.RetentionList | Tr` | *tidak dinyatakan* | `T_TREATY_LIMIT_DETAIL` | `LIMIT_DETAIL_ID` | CASCADE |
| 8 | `T_TREATY_LIMIT_ACHIEVEMENT` | `TreatyIn.Limits.Detail.AchievementLists` | *tidak dinyatakan* | `T_TREATY_LIMIT_DETAIL` | `LIMIT_DETAIL_ID` | CASCADE |
| 9 | `T_TREATY_LIMIT_GROUP` | `TreatyIn.Limits.TreatyGroupList` | *tidak dinyatakan* | `T_TREATY_LIMITS` | `LIMIT_ID` | CASCADE |
| 10 | `T_TREATY_LIMIT_GROUP_COB` | `TreatyIn.Limits.TreatyGroupList.ClassOfBusinessList` | *tidak dinyatakan* | `T_TREATY_LIMIT_GROUP` | `LIMIT_GROUP_ID` | CASCADE |
| 11 | `T_TREATY_SHARE` | `TreatyIn.Share` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 12 | `T_TREATY_SHARE_SPREADING` | `TreatyIn.Share.SpreadingListXOL` | *tidak dinyatakan* | `T_TREATY_SHARE` | `SHARE_ID` | CASCADE |
| 13 | `T_TREATY_SHARE_SPREAD_AMOUNT` | `TreatyIn.Share.SpreadingListXOL.RNMSpreadedListXOL | TreatyIn.Share.Spread` | *tidak dinyatakan* | `T_TREATY_SHARE_SPREADING` | `SHARE_SPREADING_ID` | CASCADE |
| 14 | `T_TREATY_SHARE_AMOUNT` | `TreatyIn.Share.GrossPremiumList | TreatyIn.Share.NetPremiumList | TreatyIn` | *tidak dinyatakan* | `T_TREATY_SHARE` | `SHARE_ID` | CASCADE |
| 15 | `T_TREATY_SHARE_DEDUCTION` | `TreatyIn.Share.DeductionList` | *tidak dinyatakan* | `T_TREATY_SHARE` | `SHARE_ID` | CASCADE |
| 16 | `T_TREATY_FAC_SHARE` | `TreatyIn.FacultativeShareList` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 17 | `T_TREATY_FAC_SHARE_AMOUNT` | `TreatyIn.FacultativeShareList.GrossPremiumList | TreatyIn.FacultativeShare` | *tidak dinyatakan* | `T_TREATY_FAC_SHARE` | `FAC_SHARE_ID` | CASCADE |
| 18 | `T_TREATY_FAC_SHARE_DEDUCTION` | `TreatyIn.FacultativeShareList.DeductionList` | *tidak dinyatakan* | `T_TREATY_FAC_SHARE` | `FAC_SHARE_ID` | CASCADE |
| 19 | `T_TREATY_FAC_REINSURER` | `TreatyIn.ShareFacultativeReinsurers | TreatyIn.FacultativeShareList.ShareF` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 20 | `T_TREATY_FAC_LIMITS` | `TreatyIn.ShareFacultativeReinsurers.FacultativeLimits` | *tidak dinyatakan* | `T_TREATY_FAC_REINSURER` | `FAC_REINSURER_ID` | CASCADE |
| 21 | `T_TREATY_FAC_LIMIT_DETAIL` | `TreatyIn.ShareFacultativeReinsurers.FacultativeLimits.Detail` | *tidak dinyatakan* | `T_TREATY_FAC_LIMITS` | `FAC_LIMIT_ID` | CASCADE |
| 22 | `T_TREATY_CURRENCY` | `TreatyIn.CurrencyList` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 23 | `T_TREATY_EGNPI` | `TreatyIn.EGNPI` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 24 | `T_TREATY_RETENTION` | `TreatyIn.Retention` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 25 | `T_TREATY_LIMIT_SUMMARY` | `TreatyIn.LimitSummaryList | TreatyIn.LimitShareSummaryList | TreatyIn.Limi` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 26 | `T_TREATY_INSTALLMENT` | `TreatyIn.Installment` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 27 | `T_TREATY_INSTALLMENT_ITEM` | `TreatyIn.ValueDifference.Installment.InstallmentList` | *tidak dinyatakan* | `T_TREATY_INSTALLMENT` | `INSTALLMENT_ID` | CASCADE |
| 28 | `T_TREATY_REPORTING_PERIOD` | `TreatyIn.ReportingPeriodList` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 29 | `T_TREATY_ACCUMULATION` | `TreatyIn.AccumulationList` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 30 | `T_TREATY_PORTFOLIO` | `TreatyIn.Portfolio` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 31 | `T_TREATY_HAZARD_LIMIT` | `TreatyIn` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 32 | `T_TREATY_TOTAL` | `TreatyIn.TotalShareNetNP | TreatyIn.TotalShareGrossNP | TreatyIn.TotalEgnp` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 33 | `T_TREATY_RETRO_SHARE` | `TreatyIn.ShareReins` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 34 | `T_VIEW_COMMENT` | `TreatyIn.CommentList` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |
| 35 | `T_TREATY_LIMIT_MEASURE` | `TreatyIn.Limits.MDPList | TreatyIn.Limits.EgnpiTotalList | TreatyIn.Limits` | *tidak dinyatakan* | `T_TREATY_LIMITS` | `LIMIT_ID` | CASCADE |
| 36 | `T_TREATY_REINSTATEMENT` | `TreatyIn.Limits.Reinstatement_List` | *tidak dinyatakan* | `T_TREATY_LIMITS` | `LIMIT_ID` | CASCADE |
| 37 | `T_TREATY_VALUE_BEFORE_PRORATE` | `TreatyIn.ValueBeforeProrate` | *tidak dinyatakan* | `TREATY_IN` | `TREATY_IN_ID` | CASCADE |

## 2. Cara mengisinya — asal tiap baris

| Tabel datar | Diisi dari jalur Pega | Padanan model baru | Tingkat bukti |
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

## 3. Batas — dibaca sebelum memakai daftar ini

| | |
|---|---|
| Bukti **kuat** saja — `JALUR-PENUH` | **6** tabel |
| Bukti **lemah** saja — `DAUN-RELATIF` | **17** tabel |
| Membawa **kedua** tanda | **14** tabel |
| Tanpa tingkat bukti tertulis | **0** tabel |
| **Jumlah** | **37** |

> **Keempat golongan saling lepas, dan jumlahnya diperiksa sama dengan cacah tabel.** Sebuah tabel dapat membawa **kedua** tanda — baris buktinya berbunyi `DAUN-RELATIF/JALUR-PENUH` — sehingga menjumlahkan "kuat" dan "lemah" begitu saja menghasilkan angka yang melampaui cacah tabelnya.

> **`DAUN-RELATIF` lemah, dan sebabnya disebut di sumbernya:** hanya nama ruas terakhir yang ditemukan — misalnya `.COBList` — sehingga ia **mungkin memungut nama milik entitas lain**.

> **Nol rujukan berarti BELUM TERPERIKSA, bukan tidak dipakai.** Sapuan yang menghasilkan daftar ini hanya membaca **Section**. Activity, Data Transform, Report Definition, dan RDB List **tidak disapu**.

> **`ON DELETE CASCADE` adalah USULAN rancangan, bukan perilaku sistem lama.** Ia mengikuti konvensi `contooh.xlsx`.

