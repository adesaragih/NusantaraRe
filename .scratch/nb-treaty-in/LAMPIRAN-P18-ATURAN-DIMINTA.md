# ⛔ DIBATALKAN — lampiran ini tidak jadi dikirim

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22
>
> **Lampiran ini dibatalkan bersama suratnya. Tidak ada satu pun aturan di bawah yang perlu
> diekspor ulang.**
>
> Daftar ini disusun atas dasar pernyataan tim migrasi bahwa isi langkah penetapan nilai tidak
> ikut terekspor. **Pernyataan itu keliru, dan kekeliruannya milik tim migrasi sepenuhnya.**
> Isi langkah **ada di dalam ekspor yang sudah kami terima**, di tag `PropertiesName` dan
> `PropertiesValue` — **tanpa awalan `py`**. Tim migrasi memeriksa `pyPropRef`,
> `pyPropertiesName`, dan `pyPropertiesValue` — **dengan** awalan. ⭐ **Sebabnya dua huruf.**
>
> Terukur pada ke-51 aturan di bawah: **269 langkah `Property-Set`**, **707 pasangan nama=nilai
> terisi**, **nol** tanda nilai terpotong. Rinciannya di `VERIFIKASI-P18.md`.
>
> Bunyi permintaan yang ditarik, dikutip utuh:
> > *"Lampiran P18 — daftar aturan yang diminta diekspor ulang. **51 aturan, 268 langkah
> > penetapan nilai.** Seluruhnya di `NB Treaty In\Activity\`. Diurutkan dari yang paling berat.
> > Bila tidak seluruhnya dapat dikerjakan, kerjakan dari atas — sepuluh teratas sudah mencakup
> > 56 % dari seluruh langkah, dua puluh teratas 76 %."*
>
> ⛔ **Tidak ada yang perlu dikirim. Nol aturan.**

---

## ⚠️ Tabel di bawah disimpan sebagai arsip, bukan sebagai permintaan

Daftarnya **tidak dihapus** — ia tetap berguna sebagai peta 51 aturan terjangkau yang memuat
rantai perhitungan uang, dan urutan beratnya tetap sahih. Yang gugur adalah **permintaannya**,
bukan daftarnya.

## Migrasi Treaty Inward — modul *NB Treaty In*

**51 aturan, 269 langkah penetapan nilai terukur** *(daftar ini semula menulis 268; selisih 1 pada
`SetValidateInstallment_Act` — lampiran 2, terukur 3)*. Seluruhnya di `NB Treaty In\Activity\`.

| # | Aturan | Langkah | Kumulatif |
| ---: | --- | ---: | ---: |
| 1 | `InsertToTreatyOutXOLList` | 30 | 11 % |
| 2 | `InsertToTreatyXOLListRetroShare` | 23 | 19 % |
| 3 | `GeneratePolicyNoTreaty_Act` | 17 | 26 % |
| 4 | `InsertToTreatyXOLList` | 16 | 32 % |
| 5 | `FillPaymentInstallment` | 15 | 37 % |
| 6 | `InputPolicyTreatyInPre_Act` | 12 | 42 % |
| 7 | `InputPolicyTreatyEDMDetail_NP` | 12 | 46 % |
| 8 | `TreatyNonPropSetSpreading` | 10 | 50 % |
| 9 | `Protection_Act` | 9 | 53 % |
| 10 | `InsertHistoryAkseptasiPega` | 8 | 56 % |
| 11 | `TreatyNonPropOutSetSpreading` | 7 | 59 % |
| 12 | `TreatyInNonSetTotal` | 7 | 61 % |
| 13 | `CountSpreading_Act` | 6 | 64 % |
| 14 | `SaveViewSuggest` | 5 | 66 % |
| 15 | `SaveJsonPolisTreatyIn_Act` | 5 | 67 % |
| 16 | `CountNetPremi_act` | 5 | 69 % |
| 17 | `CheckDuplicateOffer` | 5 | 71 % |
| 18 | `CekLimitTreatyAcc_Act` | 5 | 73 % |
| 19 | `SetReinstatementPct` | 4 | 75 % |
| 20 | `CountResult1_Act` | 4 | 76 % |
| 21 | `TreatyRealizationCheckDuplicate` | 3 | 77 % |
| 22 | `SetPPNPPH` | 3 | 78 % |
| 23 | `SearchHierarkiSourceBizAgentTreatyIn_Act` | 3 | 79 % |
| 24 | `FetchTreatyGroupOldID` | 3 | 80 % |
| 25 | `CountOGPONP_Act` | 3 | 82 % |
| 26 | `BreakDownSpreading_Act` | 3 | 83 % |
| 27 | `AgentSourceBizTreatyIn_Act` | 3 | 84 % |
| 28 | `TreatyRealizationCheckXOLList` | 2 | 85 % |
| 29 | `TreatyInInputVis` | 2 | 85 % |
| 30 | `SetValidateInstallment_Act` | 2 | 86 % |
| 31 | `SetTreatyCurrencyID` | 2 | 87 % |
| 32 | `SetDueTo_act` | 2 | 88 % |
| 33 | `SetCurrency_act` | 2 | 88 % |
| 34 | `SetCategoryAttach` | 2 | 89 % |
| 35 | `FetchTreatyGroupOJK` | 2 | 90 % |
| 36 | `CountRiCommOnp_act` | 2 | 91 % |
| 37 | `CountRiCommOgp_act` | 2 | 91 % |
| 38 | `CountResult2Onp_act` | 2 | 92 % |
| 39 | `CountResult2Ogp_act` | 2 | 93 % |
| 40 | `CountResult1Onp_Act` | 2 | 94 % |
| 41 | `CountOverridingCommOnp_Act` | 2 | 94 % |
| 42 | `CountOverridingCommOgp_Act` | 2 | 95 % |
| 43 | `ConcatSlipOfferNo_Act` | 2 | 96 % |
| 44 | `CheckDataMkt` | 2 | 97 % |
| 45 | `CalculatePremi_Act` | 2 | 97 % |
| 46 | `TreatyInputPctCommSpreading` | 1 | 98 % |
| 47 | `SetSurveyReport_Act` | 1 | 98 % |
| 48 | `ProtectDate` | 1 | 98 % |
| 49 | `InputParamUploadReas_act` | 1 | 99 % |
| 50 | `CountPctInstallment_Act` | 1 | 99 % |
| 51 | `ConvertHistoryDate` | 1 | 100 % |

---

## Cara daftar ini disusun

Folder `NB Treaty In` adalah **dependency closure** — ia memuat setiap aturan yang disentuh
modul ini, termasuk milik modul lain. Daftar ini sudah disaring dua kali:

| Tahap | Aturan | Langkah |
| --- | ---: | ---: |
| seluruh isi folder | 92 | 938 |
| dibuang: tidak terjangkau dari titik masuk mana pun | −28 | −561 |
| dibuang: tidak dimigrasi karena penyimpanan JSON ditinggalkan | −5 | −109 |
| **diminta** | **51** | **268** |

Penjangkauan ditelusuri dari titik masuk nyata — `Flow`, `Section`, `Harness`, `FlowAction`,
dan `DataTransform` — lalu diikuti sampai tidak ada lagi aturan baru yang tercapai.

Yang dibuang pada tahap kedua adalah aturan yang hanya ada untuk membongkar dan menyusun
dokumen JSON. Penyimpanan pindah ke tabel relasional, sehingga kelimanya tidak dibangun ulang.

*Disusun 2026-09-22. Angka dapat diuji ulang dari korpus.*
