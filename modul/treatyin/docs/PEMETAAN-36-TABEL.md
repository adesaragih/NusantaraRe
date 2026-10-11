# Pemetaan 36 tabel `T_TREATY_*` — sumber, tabel model, dan keadaan hari ini

Dibangkitkan dari `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`, 6 Oktober 2026.
Letaknya: `_migration-docs/treaty-in/4-erd-dan-tabel-datar/` di dalam `D:\XML_NURE`.
Tiap baris menyebut **lembar dan selnya** supaya dapat diperiksa ulang tanpa mempercayai berkas ini.

⚠ Baris pertama xlsx itu berbunyi *"POTRET SISTEM LAMA, 24 September 2026. BUKAN RANCANGAN"*.
Pemilik proses diberi tahu dua kali dan dua kali tetap memilihnya. Dicatat sekali supaya
pilihannya terbaca SADAR, bukan terlewat.

## 1 · Ke-36 tabel

| `T_TREATY_*` | lembar · sel | sumber dokumen | tabel model `[…]` | DDL usulan | ada pada kita |
| --- | --- | --- | --- | :-: | :-: |
| `T_TREATY_ACCUMULATION` | Treaty In Prop · G126 | `TreatyIn.AccumulationList` | `AKUMULASI` | — | — |
| `T_TREATY_CURRENCY` | Treaty In Prop · O95 | `TreatyIn.CurrencyList` | `MATA_UANG_KONTRAK` | ✔ | — |
| `T_TREATY_EGNPI` | Treaty In Prop · G101 | `TreatyIn.EGNPI` | `EGNPI` | ✔ | ✔ |
| `T_TREATY_FAC_LIMITS` | Treaty In Prop · G87 | `TreatyIn.ShareFacultativeReinsurers.FacultativeLimits` | `LAYER_FAKULTATIF` | — | — |
| `T_TREATY_FAC_LIMIT_DETAIL` | Treaty In Prop · K91 | `TreatyIn.ShareFacultativeReinsurers.FacultativeLimits.Detail` | `DETAIL_PROPORSIONAL_FAKULTATIF` | — | — |
| `T_TREATY_FAC_REINSURER` | Treaty In Prop · K82 | `TreatyIn.ShareFacultativeReinsurers` | `REASURADUR_FAKULTATIF` | — | — |
| `T_TREATY_FAC_SHARE` | Treaty In Prop · K69 | `TreatyIn.FacultativeShareList` | `BAGIAN_FAKULTATIF` | — | — |
| `T_TREATY_FAC_SHARE_AMOUNT` | Treaty In Prop · G74 | `TreatyIn.FacultativeShareList.GrossPremiumList` | `NILAI_BAGIAN_FAKULTATIF` | — | — |
| `T_TREATY_FAC_SHARE_DEDUCTION` | Treaty In Prop · K78 | `TreatyIn.FacultativeShareList.DeductionList` | `POTONGAN_FAKULTATIF` | — | — |
| `T_TREATY_HAZARD_LIMIT` | Treaty In Prop · G134 | `TreatyIn` | `BATAS_BAHAYA` | — | — |
| `T_TREATY_INSTALLMENT` | Treaty In Prop · G113 | `TreatyIn.Installment` | `ANGSURAN` | — | — |
| `T_TREATY_INSTALLMENT_ITEM` | Treaty In Prop · G117 | `TreatyIn.ValueDifference.Installment.InstallmentList` | `RINCIAN_ANGSURAN` | ✔ | ✔ |
| `T_TREATY_LIMITS` | Treaty In Prop · K16 | `TreatyIn.Limits` | `LAYER` | ✔ | ✔ |
| `T_TREATY_LIMIT_ACHIEVEMENT` | Treaty In Prop · O33 | `TreatyIn.Limits.Detail.AchievementLists` | `PENCAPAIAN` | ✔ | ✔ |
| `T_TREATY_LIMIT_AMOUNT` | Treaty In Prop · O29 | `TreatyIn.Limits.Detail.EPIList` | `NILAI_LAYER` | — | — |
| `T_TREATY_LIMIT_COB` | Treaty In Prop · K25 | `TreatyIn.Limits.Detail.COBList` | `KELAS_BISNIS_LAYER` | ✔ | ✔ |
| `T_TREATY_LIMIT_DETAIL` | Treaty In Prop · G21 | `TreatyIn.Limits.Detail` | `DETAIL_PROPORSIONAL` | ✔ | ✔ |
| `T_TREATY_LIMIT_GROUP` | Treaty In Prop · O37 | `TreatyIn.Limits.TreatyGroupList` | `KELOMPOK_LAYER` | ✔ | ✔ |
| `T_TREATY_LIMIT_GROUP_COB` | Treaty In Prop · K42 | `TreatyIn.Limits.TreatyGroupList.ClassOfBusinessList` | `KELAS_BISNIS_KELOMPOK` | ✔ | ✔ |
| `T_TREATY_LIMIT_MEASURE` | Treaty In Non Prop · O37 | `TreatyIn.Limits.MDPList` | `BESARAN_LAYER` | — | — |
| `T_TREATY_LIMIT_SUMMARY` | Treaty In Prop · G109 | `TreatyIn.LimitSummaryList` | `RINGKASAN_LIMIT` | — | — |
| `T_TREATY_PORTFOLIO` | Treaty In Prop · G130 | `TreatyIn.Portfolio` | `PORTOFOLIO` | ✔ | ✔ |
| `T_TREATY_REINSTATEMENT` | Treaty In Non Prop · K42 | `TreatyIn.Limits.Reinstatement_List` | `PEMULIHAN_LIMIT` | ✔ | ✔ |
| `T_TREATY_REPORTING_PERIOD` | Treaty In Prop · K121 | `TreatyIn.ReportingPeriodList` | `PERIODE_PELAPORAN` | ✔ | ✔ |
| `T_TREATY_RETENTION` | Treaty In Prop · G105 | `TreatyIn.Retention` | `RETENSI` | — | — |
| `T_TREATY_RETRO_SHARE` | Treaty In Prop · G142 | `TreatyIn.ShareReins` | `BAGIAN_RETRO` | — | — |
| `T_TREATY_REVISION` | Treaty In Prop · C8 | `TreatyIn` | `VERSI_KONTRAK · bagian addendum` | ✔ | ✔ |
| `T_TREATY_SHARE` | Treaty In Prop · O46 | `TreatyIn.Share` | `BAGIAN` | ✔ | ✔ |
| `T_TREATY_SHARE_AMOUNT` | Treaty In Prop · O60 | `TreatyIn.Share.GrossPremiumList` | `NILAI_BAGIAN` | — | — |
| `T_TREATY_SHARE_DEDUCTION` | Treaty In Prop · K65 | `TreatyIn.Share.DeductionList` | `POTONGAN` | ✔ | ✔ |
| `T_TREATY_SHARE_SPREADING` | Treaty In Prop · G52 | `TreatyIn.Share.SpreadingListXOL` | `PENYEBARAN_XOL` | — | — |
| `T_TREATY_SHARE_SPREAD_AMOUNT` | Treaty In Prop · K56 | `TreatyIn.Share.SpreadingListXOL.RNMSpreadedListXOL` | `NILAI_TERSEBAR` | — | — |
| `T_TREATY_TOTAL` | Treaty In Prop · G138 | `TreatyIn.TotalShareNetNP` | `REKAP_KONTRAK` | — | — |
| `T_TREATY_VALUE_BEFORE_PRORATE` | Treaty In EDM Non Prop · G143 | `TreatyIn.ValueBeforeProrate` | `NILAI_SEBELUM_PRO_RATE` | — | ✔ |
| `T_TREATY_VALUE_DIFFERENCE` | Treaty In Prop · G12 | `TreatyIn.ValueDifference` | `NILAI_SELISIH` | — | ✔ |
| `T_VIEW_COMMENT` | Treaty In Prop · G146 | `TreatyIn.CommentList` | `CATATAN_PERSETUJUAN` | ✔ | ✔ |

**16 dari 36** punya DDL usulannya di `2-to-spec/ddl-usulan/` · **17 dari 36** sudah
berdiri pada kita.

---

## 2 · ⛔ RALAT atas bacaan pertama berkas ini

Bacaan pertamaku keliru, dan koreksinya datang dari pemilik proses: *"paling puncak
tetap `TREATY_IN` dan anak-anaknya mengikuti yang di xlsx."*

Aku semula membaca nama dalam kurung siku (`[AKUMULASI]`, `[LAYER]`) sebagai **tabel
model** kita, lalu menyimpulkan migrasi `436` menaruh nama `T_TREATY_*` pada tabel yang
salah. **Itu salah.** Kurung siku hanyalah penunjuk asal-usul; yang menentukan struktur
adalah **garis penggantungnya**, dan garis itu jelas:

```
TREATY_IN                                PK ID   ← AKAR, tabel warisan
├─ T_TREATY_REVISION            1:1   FK TREATY_IN_ID → TREATY_IN.ID
│   └─ T_TREATY_VALUE_DIFFERENCE 1:1  FK REVISION_ID → T_TREATY_REVISION.ID
├─ T_TREATY_LIMITS              1:N   FK TREATY_IN_ID
│   ├─ T_TREATY_LIMIT_DETAIL    1:N   FK LIMIT_ID
│   │   ├─ T_TREATY_LIMIT_COB          FK LIMIT_DETAIL_ID
│   │   ├─ T_TREATY_LIMIT_AMOUNT       FK LIMIT_DETAIL_ID
│   │   └─ T_TREATY_LIMIT_ACHIEVEMENT  FK LIMIT_DETAIL_ID
│   └─ T_TREATY_LIMIT_GROUP     1:N   FK LIMIT_ID
│       └─ T_TREATY_LIMIT_GROUP_COB    FK LIMIT_GROUP_ID
├─ T_TREATY_SHARE               1:N   FK TREATY_IN_ID
│   ├─ T_TREATY_SHARE_SPREADING       FK SHARE_ID
│   │   └─ T_TREATY_SHARE_SPREAD_AMOUNT  FK SHARE_SPREADING_ID
│   ├─ T_TREATY_SHARE_AMOUNT          FK SHARE_ID
│   └─ T_TREATY_SHARE_DEDUCTION       FK SHARE_ID
├─ T_TREATY_FAC_SHARE           1:N   FK TREATY_IN_ID
│   ├─ T_TREATY_FAC_SHARE_AMOUNT      FK FAC_SHARE_ID
│   └─ T_TREATY_FAC_SHARE_DEDUCTION   FK FAC_SHARE_ID
├─ T_TREATY_FAC_REINSURER       1:N   FK TREATY_IN_ID
│   └─ T_TREATY_FAC_LIMITS            FK FAC_REINSURER_ID
│       └─ T_TREATY_FAC_LIMIT_DETAIL     FK FAC_LIMIT_ID
├─ T_TREATY_INSTALLMENT         1:N   FK TREATY_IN_ID
│   └─ T_TREATY_INSTALLMENT_ITEM      FK INSTALLMENT_ID
└─ langsung di bawah TREATY_IN, masing-masing FK TREATY_IN_ID:
    T_TREATY_EGNPI · T_TREATY_RETENTION · T_TREATY_LIMIT_SUMMARY ·
    T_TREATY_REPORTING_PERIOD · T_TREATY_ACCUMULATION · T_TREATY_PORTFOLIO ·
    T_TREATY_HAZARD_LIMIT · T_TREATY_TOTAL · T_TREATY_RETRO_SHARE ·
    T_VIEW_COMMENT · T_TREATY_CURRENCY
```

⭐ Jadi migrasi `436` **BENAR** menaruh nama itu pada tabel yang menggantung pada
`TREATY_IN`. Tabel model kita (`PERIODE_AKUMULASI`, `BAGIAN`, `LAYER`, …) menggantung
pada `KONTRAK`/`VERSI_KONTRAK` — pohon yang berbeda, dan bukan yang xlsx ini gambar.

### ⛔ Satu penghalang yang mengenai SELURUH pohon

Setiap anak tingkat pertama membawa **`FK TREATY_IN_ID → TREATY_IN.ID CASCADE`**, dan
xlsx menulis `TREATY_IN   PK ID` di sel `C6`. Terukur di POOLDATA 6 Oktober 2026:

| | |
|---|---|
| kunci utama / unik pada `TREATY_IN` | **0** |
| index `UNIQUE` pada `TREATY_IN` | **0** |
| baris · `ID` kembar | **1.854** · **0** |

⭐ Datanya **sebenarnya unik** — kuncinya DAPAT dipasang. Yang menghalangi bukan data,
melainkan aturan: `TREATY_IN` tabel **warisan**, dan nol DDL boleh dikenakan padanya.

Dua jalan, dan keduanya milik pemilik proses:

1. **Pasang `UNIQUE (ID)` pada `TREATY_IN`** — satu-satunya cara pohon ini berdiri
   sebagaimana digambar. ⚠ Ia menyentuh tabel warisan, dan larangan itu berlaku sejak
   modul ini lahir.
2. **Anak menaut dengan `MASTERID` teks, tanpa kunci asing** — persis yang kesembilan
   tabel pendaratan lakukan hari ini. Pohon tetap terbentuk, penjaganya hilang: baris
   yatim tidak ditolak basis data, dan `CASCADE` harus ditiru aplikasi.

## 3 · Tiga bentrok dengan keputusan yang sudah berdiri

**3.1 · `T_TREATY_CURRENCY ← TreatyIn.CurrencyList [MATA_UANG_KONTRAK]`**
⛔ `MATA_UANG_KONTRAK` **DICABUT** migrasi `434`, keputusan §16 — kurs diambil
dari `TREATYEXCHANGEYEARLY`. DDL usulannya masih ada (`18_MATA_UANG_KONTRAK.sql`),
tetapi tabelnya sengaja dihapus. **Jangan dibangun sebelum §16 ditinjau.**

**3.2 · `T_TREATY_VALUE_DIFFERENCE` dan `T_TREATY_VALUE_BEFORE_PRORATE`**
Keduanya ada di lembar Treaty In, tetapi `NILAI_SELISIH` dan
`NILAI_SEBELUM_PRO_RATE` sudah **dipindahkan ke modul Adjustment** oleh §19
(migrasi `442`). **Jangan dibangun ulang di sini.**

**3.3 · Lima belas tabel yang belum berdiri — dan sebabnya BUKAN jalur terpotong**

⛔ **RALAT, 6 Oktober 2026.** Catatan ronde sebelumnya menyatakan tujuh tabel
tertahan karena *"jalur sumbernya terpotong di dalam sel xlsx"*. Itu **separuh
benar dan menyesatkan**: sel `←` memang memotong daftar jalur alternatifnya pada
sekitar 75 aksara (`… | TreatyIn.TotalEgnp`), tetapi **jalur PERTAMA setiap sel
utuh**, dan yang sesungguhnya menahan bukan panjang sel melainkan sebuah putusan
yang sudah tertulis di tempat lain.

Tempat itu `2-to-spec/PENELUSURAN-JSON-KE-KOLOM.md` — 377 jalur JSON, masing-masing
dengan putusan bertanda `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)`. Diadu dengan
ke-15 tabel yang belum berdiri, 6 Oktober 2026:

| Tabel | Putusan jalur sumbernya | Keadaan |
|---|---|---|
| `T_TREATY_LIMIT_AMOUNT` | 8 `DIPETAKAN`, 0 `DITURUNKAN` | **dapat dibangun** |
| `T_TREATY_FAC_SHARE_AMOUNT` | 2 `DIPETAKAN` | **dapat dibangun** |
| `T_TREATY_FAC_LIMITS` + `_DETAIL` | 5 `DIPETAKAN`, 3 `DITURUNKAN` | ⛔ **NOL ELEMEN** di seluruh 1.854 dokumen — jangan dibangun |
| `T_TREATY_SHARE_AMOUNT` | 2 `DIPETAKAN`, 1 `DITURUNKAN` | **dapat dibangun** tanpa kolom turunan |
| `T_TREATY_LIMIT_MEASURE` | **7 dari 7 `DITURUNKAN`** | menunggu putusan di bawah |
| `T_TREATY_LIMIT_SUMMARY` | **4 dari 4 `DITURUNKAN`** | menunggu putusan di bawah |
| `T_TREATY_REINSTATEMENT` | **2 dari 2 `DITURUNKAN`** | menunggu putusan di bawah |
| `T_TREATY_SHARE_SPREAD_AMOUNT` | **20 dari 20 `DITURUNKAN`** | menunggu putusan di bawah |
| `T_TREATY_TOTAL` | **16 dari 16 `DITURUNKAN`** | menunggu putusan di bawah |
| `T_TREATY_HAZARD_LIMIT` | sumbernya AKAR dokumen (`← TreatyIn`) | menunggu daftar medannya |
| `T_TREATY_REVISION` | sumbernya AKAR dokumen (`← TreatyIn`) | menunggu daftar medannya |
| `T_TREATY_CURRENCY` | `DIPETAKAN·TABEL` | tertahan §16 — lihat 3.1 |
| `T_TREATY_VALUE_DIFFERENCE` / `_BEFORE_PRORATE` | — | milik Adjustment, §19 — lihat 3.2 |

⚠️ **`DIPETAKAN` bukan izin membangun.** `T_TREATY_FAC_LIMITS` jalurnya utuh dan
lima kolomnya `DIPETAKAN`, tetapi sapuan 1.854 dokumen menemukan **nol elemen** — dan
tabel tanpa satu pun baris muatan bukan tabel. Cacah elemen setiap calon tabel karena
itu **diukur lebih dulu**, atas seluruh korpus, bukan atas `ROWNUM` sebagian.

**⛔ SATU PERTANYAAN memutuskan lima tabel sekaligus, dan ia belum pernah diajukan.**

`ADR-0037` berbunyi *"turunan — tidak disimpan"*, dan ia berbicara tentang **model
relasional baru** (`KONTRAK`, `VERSI_KONTRAK`, `LAYER`, …): nilai yang dapat dihitung
ulang tidak diberi kolom di sana. Tabel `T_TREATY_*` adalah hal yang **berbeda** — ia
**pendaratan**, yang mendaratkan dokumen Pega apa adanya supaya layar dapat
menampilkan apa yang sistem lama tampilkan.

⚠️ Dan kedua hal itu sudah terbukti berselisih di layar yang sedang berjalan:
`Limits[].MDPList` dan `Limits[].PremiumEarnedList` ditandai `DITURUNKAN` di dokumen
itu, **tetapi tab Limits menampilkan keduanya hari ini**. Jadi `DITURUNKAN` tidak dapat
dibaca sebagai "jangan didaratkan" tanpa seseorang menyatakannya.

> **Pertanyaannya:** apakah tabel pendaratan `T_TREATY_*` ikut membawa nilai
> `DITURUNKAN`, atau hanya yang `DIPETAKAN`?
>
> **Jawab "ikut"** → kelima tabel di atas dibangun, dan kolomnya diturunkan dari
> dokumen seperti ketiga belas tabel migrasi `437`.
> **Jawab "tidak"** → kelimanya **tidak pernah dibangun**, dan baris itu dicoret dari
> ke-36, bukan dibiarkan menggantung sebagai pekerjaan yang belum selesai.

⚠️ Selama pertanyaan itu belum dijawab, **menebak ke arah mana pun berbiaya**:
membangun lima tabel yang tidak boleh ada sama mahalnya dengan tidak membangun lima
yang diperlukan.

**3.4 · Kolom memang tidak ada di xlsx, dan itu disengaja**
xlsx menyatakannya sendiri: *"kotak memuat PK/FK dan sumber Pega saja, KOLOM SENGAJA
TIDAK DIMUAT"*. Karena itu kolom setiap tabel pendaratan **diturunkan dari dokumen**
dan diadu ke data — cara yang sama yang menghasilkan
[`TURUNAN-KOLOM-28-ANAK.md`](TURUNAN-KOLOM-28-ANAK.md) bagi ketiga belas tabel `437`.

⚠️ **Dan xlsx itu sendiri menyatakan dirinya BUKAN rancangan.** Lembar
"Catatan & Batas" baris pertama: *"POTRET SISTEM LAMA, 24 September 2026. BUKAN
RANCANGAN"*, lalu menunjuk `2-to-spec/KAMUS-KOLOM.md` dan `2-to-spec/ddl-usulan/`
sebagai tempat skema yang dibangun. Yang diambil darinya karena itu **nama tabel dan
bentuk pohonnya** — sesuai keputusan pemilik proses 5 Oktober 2026 — bukan bentuk
kolomnya.

**3.6 · ⭐ Cacah elemen ketiga tabel — DIUKUR 6 Oktober 2026, lalu dibangun**

Sapuan atas **SELURUH 1.855 baris** `POOLDATA.M_TREATY_IN` (nol `ROWNUM`, nol sampel; 1.855
dokumen terurai, nol gagal urai):

| Tabel | Jalur sumber | Elemen | Dokumen |
| --- | --- | ---: | ---: |
| `T_TREATY_LIMIT_AMOUNT` | `Limits[].Detail[].IOOLimitList` | 2.870 | 1.077 |
| | `Limits[].Detail[].RetentionList` | 2.870 | 1.078 |
| | `Limits[].Detail[].CessionList` | 2.868 | 1.077 |
| | `Limits[].Detail[].EPIList` | 2.867 | 1.039 |
| | **jumlah** | **11.475** | |
| `T_TREATY_SHARE_AMOUNT` | `Share[].GrossPremiumList` | 3.049 | 845 |
| | `Share[].NetPremiumList` | 3.048 | 845 |
| | **jumlah** | **6.097** | |
| `T_TREATY_FAC_SHARE_AMOUNT` | `FacultativeShareList[].GrossPremiumList` | 12 | 5 |
| | `FacultativeShareList[].NetPremiumList` | 12 | 5 |
| | **jumlah** | **24** | |

⭐ Ketiganya **lebih besar dari nol**, jadi ketiganya dibangun — berbeda dari `T_TREATY_FAC_LIMITS`
yang jalurnya `DIPETAKAN` tetapi cacahnya **nol**.

⚠️ `T_TREATY_FAC_SHARE_AMOUNT` **paling tipis** — 24 elemen di 5 dari 1.855 dokumen. Ia dibangun
sebab bukan nol, dan ketipisannya dicatat supaya tidak perlu diukur ulang.

⚠️ **Satu penambahan di luar jalur yang disebut tugas:** `FacultativeShareList[].NetPremiumList`
(12 elemen) ikut ditampung, sebab bentuknya sama persis dengan `GrossPremiumList` dan tabelnya
sudah berkolom `JENIS`. Tanpa itu, 12 elemen terukur tidak punya rumah. **Dapat dibalikkan tanpa
biaya** — tabelnya nol baris.

⚠️ **Satu kolom yang tidak ada di dokumen:** `JENIS`. Ia DIPAKSA oleh bentuk yang xlsx pilih —
empat larik sumber ke dalam satu tabel — bukan dipilih di sini. Tanpa pembeda, baris retensi dan
baris cession menjadi tidak dapat dibedakan.

**3.5 · Nama kolom kunci asing BERBEDA dari xlsx, dan itu disengaja**
Lembar "Daftar Relasi" menamai setiap kunci tamu menurut induknya: `LIMIT_ID`,
`LIMIT_DETAIL_ID`, `SHARE_ID`, `FAC_SHARE_ID`, `TREATY_IN_ID`. Migrasi `437` memakai
`IDINDUK` dan `MASTERID`, mengikuti pola kesembilan tabel pendaratan yang sudah berdiri
dan sudah terisi sejak migrasi `430`.

⚠️ Ini **selisih yang diketahui, bukan kelalaian**. Menyeragamkan ke bentuk xlsx
berarti menyunting `430`–`433` yang sudah dijalankan, pemuatnya, dan pembacanya
sekaligus. Bila pemilik proses menghendakinya, ia migrasi tersendiri — bukan sisipan.

---

## 4 · Yang dapat dikerjakan tanpa keputusan baru

| | |
|---|---|
| **24 tabel sudah berdiri** | migrasi `430`–`433`, `436`, `437`, **`438`** — seluruhnya nol baris |
| ⭐ **3 tabel DIBANGUN 6 Oktober 2026** | `LIMIT_AMOUNT` · `SHARE_AMOUNT` · `FAC_SHARE_AMOUNT` — migrasi **`438`**, sesudah cacah elemennya diukur atas SELURUH 1.855 dokumen (§3.6) |
| **2 tabel nol elemen** | `FAC_LIMITS`, `FAC_LIMIT_DETAIL` — jalurnya utuh, isinya tidak pernah ada |
| **5 tabel menunggu SATU pertanyaan** | `LIMIT_MEASURE`, `LIMIT_SUMMARY`, `REINSTATEMENT`, `SHARE_SPREAD_AMOUNT`, `TOTAL` — §3.3 |
| **2 tabel menunggu daftar medan akar** | `HAZARD_LIMIT`, `REVISION` |
| **`T_TREATY_CURRENCY`** | menunggu peninjauan §16 |
| **2 tabel nilai selisih** | milik modul Adjustment, §19 |
