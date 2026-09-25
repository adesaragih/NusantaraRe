> Modul  : Treaty In Adjustment · ronde TDA
> Dibuat : 2026-09-24
> Sifat  : **nama properti sistem lama, APA ADANYA.** Ejaan, huruf besar-kecil, singkatan, dan
>          ketidakkonsistenannya **tidak diperbaiki, tidak diterjemahkan, tidak diganti.**
> Cara  : ditarik dari ekspor dengan pemindai teks atas berkas aturan yang disentuh tiap putusan.
>         Tidak satu nama pun diketik dari ingatan.

# Properti asli yang disentuh putusan ronde C dan ronde TDA

## 0. Kenapa berkas ini ada, dan apa yang ia BUKAN

Putusan `GRL-14` … `GRL-17` ditulis dengan nama sistem **baru** — `VERSI_KONTRAK`,
`NILAI_SELISIH`, `JENIS_ADDENDUM`. Nama-nama itu **usulan**, dan usulan yang menggantikan nama lama
di setiap kalimat akan membuat orang berikutnya tidak dapat lagi menemukan bendanya di ekspor.

Berkas ini jangkarnya: **nama lama, verbatim.**

> Ia **bukan** peta pemetaan lama→baru. Ia daftar nama asli. Pemetaannya ada di putusan
> masing-masing, dan di `PETA-TELUSUR-JSON.md` induk.

Dan ia **akar properti** — segmen pertama sesudah `TreatyIn.`. Jalur daun lengkap untuk
`ValueDifference` (33 akar, 165 jalur daun) sudah ada di `PENGETAHUAN.md` §5.6 dan tidak diulang.

Kolom sisi: **`ADJ`** muncul di ekspor Treaty In Adjustment, **`TI`** muncul di ekspor Treaty In.

---

## 1. `GRL-14` — `ActualValue`

Aturan yang dipindai: `SaveTreatyIn_EDM_Act`, `TreatyInSetEditPre`, `TreatyEDMDifferenceShare`,
`TreatyEDMDifferencePremium`. **16 akar.**

| Properti asli | Sisi |
|---|---|
| `TreatyIn.ActualValue` | ADJ + TI |
| `TreatyIn.AddendumPremi` | ADJ |
| `TreatyIn.Commencement` | ADJ |
| `TreatyIn.Comment` | ADJ + TI |
| `TreatyIn.CommentList` | ADJ |
| `TreatyIn.EDMEffective` | ADJ |
| `TreatyIn.EDMState` | ADJ + TI |
| `TreatyIn.EGNPI` | ADJ + TI |
| `TreatyIn.ID` | ADJ + TI |
| `TreatyIn.Information` | ADJ + TI |
| `TreatyIn.OLDDATA` | ADJ + TI |
| `TreatyIn.OLDID` | ADJ + TI |
| `TreatyIn.Share` | ADJ + TI |
| `TreatyIn.StatusAkseptasi` | ADJ + TI |
| `TreatyIn.TotalEgnpiAmount` | ADJ + TI |
| `TreatyIn.ValueDifference` | ADJ + TI |

---

## 2. `GRL-15` — pro rata

Aturan yang dipindai: `TreatyCalculateProratePct`, `TreatyEDMProRateCalculation`,
`TreatyEDMCalculateDifference`. **10 akar, seluruhnya ADJ + TI.**

| Properti asli | Sisi |
|---|---|
| `TreatyIn.Commencement` | ADJ + TI |
| `TreatyIn.EDMEffective` | ADJ + TI |
| `TreatyIn.InstallmentNo` | ADJ + TI |
| `TreatyIn.IsProRate` | ADJ + TI |
| `TreatyIn.ProRateDays` | ADJ + TI |
| `TreatyIn.ProRatePercent` | ADJ + TI |
| `TreatyIn.ProRateTotalDays` | ADJ + TI |
| `TreatyIn.Termination` | ADJ + TI |
| `TreatyIn.ValueBeforeProrate` | ADJ + TI |
| `TreatyIn.ValueDifference` | ADJ + TI |

---

## 3. `GRL-16` — share fakultatif

Aturan yang dipindai: `TreatyEDMDifferenceDeduction`, `TreatyInDifferenceFacShare`,
`TreatyInDifferenceFacShareTotal`. **7 akar, seluruhnya ADJ + TI.**

| Properti asli | Sisi |
|---|---|
| `TreatyIn.ActualValue` | ADJ + TI |
| `TreatyIn.FacultativeShare` | ADJ + TI |
| `TreatyIn.FacultativeShareBrokerage` | ADJ + TI |
| `TreatyIn.FacultativeShareList` | ADJ + TI |
| `TreatyIn.LimitShareSummaryList` | ADJ + TI |
| `TreatyIn.OLDDATA` | ADJ + TI |
| `TreatyIn.ValueDifference` | ADJ + TI |

---

## 4. `GRL-17` — penomoran dan identitas warisan

Aturan yang dipindai: `TreatyInRevisi_post`, `TreatyCreateEDM`, `TreatyInEDMSetValue`.
**6 akar, seluruhnya hanya ADJ** — ketiganya termasuk 56 aturan khas Adjustment.

| Properti asli | Sisi |
|---|---|
| `TreatyIn.CommentList` | ADJ |
| `TreatyIn.EDMMaterialType` | ADJ |
| `TreatyIn.EDMState` | ADJ |
| `TreatyIn.ID` | ADJ |
| `TreatyIn.OLDID` | ADJ |
| `TreatyIn.RevisionDate` | ADJ |

---

## 5. `ST-7` — residu jenis addendum

Aturan yang dipindai: `TreatyInSetEditPre`, `TreatyCreateEDM`, `PickerTreatyInMasterRevisi`.
**8 akar, seluruhnya hanya ADJ.**

| Properti asli | Sisi |
|---|---|
| `TreatyIn.ActualValue` | ADJ |
| `TreatyIn.AddendumPremi` | ADJ |
| `TreatyIn.Commencement` | ADJ |
| `TreatyIn.CommentList` | ADJ |
| `TreatyIn.EDMEffective` | ADJ |
| `TreatyIn.EDMMaterialType` | ADJ |
| `TreatyIn.EDMState` | ADJ |
| `TreatyIn.EGNPI` | ADJ |

---

## 6. Ejaan yang TIDAK diperbaiki, dan sebabnya disebut satu per satu

Empat pola penamaan hidup berdampingan di dalam satu kelas yang sama. Dicatat supaya tidak ada yang
"membetulkannya" kelak dan memutus rujukan:

| Pola | Contoh apa adanya |
|---|---|
| **HURUF BESAR SELURUHNYA** | `OLDDATA` · `OLDID` |
| **CamelCase** | `ActualValue` · `ValueDifference` · `FacultativeShareBrokerage` |
| **bahasa Indonesia di tengah korpus Inggris** | `StatusAkseptasi` · `AddendumPremi` |
| **singkatan tanpa kepanjangan di mana pun** | `EGNPI` · `MDP` · `RSMD` |

Dan tiga hal yang layak dibaca dua kali sebelum dipercaya:

1. **`AddendumPremi` terpotong.** Ia bendera yang disetel `TreatyInSetEditPre` untuk `EDMState = 3`.
   Namanya bukan `AddendumPremium`. **Jangan dilengkapi.**
2. **`ValueBeforeProrate` versus label langkahnya.** `TreatyEDMProRateCalculation` langkah 1
   berbunyi *"Copy value from ValueDifference to **ValueBeforeProrat**"* — labelnya terpotong satu
   huruf, propertinya tidak. Yang benar propertinya; labelnya sekadar memo penulis.
3. **`TreatyIn.OLDID` berarti TIGA hal**, dan pembedanya bukan ejaan melainkan penulisnya — versi
   pendahulu, kontrak asal salinan, dan nomor penawaran warisan. Aturan pembeda per barisnya ada di
   `../treaty-in/KEPUTUSAN-TANPA-VERIFIKASI.md` §6. **Satu nama, jangan disatukan menjadi satu
   kolom.**

---

## 7. Batas perkakas — dinyatakan di dalam berkas yang ia hasilkan

> Pemindai ini **tidak dapat menyatakan sebuah properti DIPAKAI.** Ia hanya menyatakan namanya
> **MUNCUL** di berkas itu.

- Pemindaiannya **berbasis teks, bukan nama tag.** Itu disengaja — ia kebal terhadap kelima bentuk
  penulis properti dan terhadap bentuk yang belum dikenal. Ongkosnya: kemunculan di dalam **komentar
  penulis aturan** dan di dalam **salinan terbungkus** (`pyIncludedRuleXML`, `pyRuleVersionsList`)
  **ikut terpungut**, dan tidak dipisahkan di sini.
- **Berkas aturan yang tidak ditemukan: 0.** Seluruh 13 nama aturan yang diminta ada di ekspor.
- Yang **tidak** termasuk: properti berawalan selain `TreatyIn.` — antara lain `Param.*`,
  `local.*`, `SaveData1.*`, dan `InputData.*`. Ketiadaannya di sini **bukan** pernyataan bahwa ia
  tidak ada.
