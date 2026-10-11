# Tab Share (non-prop) — aplikasi lawan ekspor Pega

Dibaca 7 Oktober 2026 dari `D:\XML_NURE\Treaty In`:

```
Section/TreatyInTabsNonProportional.xml   wadah TABBED ke-5 `Share`
Activity/TreatyInNonAddItem.xml           Type = share · sharereins · sharefacname
Activity/TreatyInXOLAddSpreading.xml      237 KB — membaca RD PROPORTIONALARRG
Activity/TreatyInSetBrokerage.xml         315 KB
Activity/TreatyInNPSetTotal.xml           param.type == "share"
Activity/TreatyInSummaryLimitShare.xml    175 KB
Activity/TreatyInShareListValue.xml
```

> ⚠️ **Berkas ini CATATAN PEMBANDING, bukan rencana kerja.** Pembangunan tab
> Share Non-Prop dipegang sesi lain (`nusantara-re-app-16`) atas permintaan
> pemakai. Lihat §6.

---

## 1 · Hasil perbandingan

### 1.1 Panel `Share` (form atas)

| Medan / tombol | Properti | Pemicu | Di aplikasi |
|---|---|---|---|
| `% RNM Share` | `TreatyIn.RNMShare` | `TreatyInXOLAddSpreading` | ⛔ nol di layar |
| `% Brokerage` | `TreatyIn.BrokeragePercent` | `TreatyInSetBrokerage` + `XOLAddSpreading` | ⛔ nol di layar |
| centang | `TreatyIn.RNMShareAcrossTheBoard` | keduanya | ⛔ nol di layar |
| `Share to Other Retro` | `TreatyIn.FacultativeShare` | `XOLAddSpreading` | ⛔ nol di layar |
| `Brokerage From Other Retro` | `TreatyIn.FacultativeShareBrokerage` | `XOLAddSpreading` | ⛔ nol di layar |
| **`Update Summary`** | — | `NonAddItem(share)` → `XOLAddSpreading` → `SetBrokerage` → `TreatyEDMCalculateDifference` | ⛔ nol di layar |

⭐ Label `Share Non Pro Rate (actual)` tampil bila `TreatyIn.IsProRate = True`.
⭐ `Brokerage From Other Retro` tampil bila `TreatyIn.FacultativeShare > 0`.

### 1.2 Dua grid reasuradur — **nol dibangun**

| Grid | Larik | Add | Kolom |
|---|---|---|---|
| `Reinsurer Name` | `TreatyIn.ShareReins` | `NonAddItem(sharereins)` | Reinsurer Name · Layer · % Share |
| `Facultative Reinsurers` | `TreatyIn.ShareFacultativeReinsurers` | `NonAddItem(sharefacname)` | Facultative Reinsurers · Layer · % Share |

Add/Delete keduanya HIDUP (`TreatyIn.ViewState != '1'`).

~~⚠️ Ketiga kolomnya `pyReadOnly = true`, padahal barisnya dapat ditambah.~~

⛔ **RALAT 7 Oktober 2026 — bacaan di atas KELIRU, dan ditarik.**

Sel-sel itu memang memuat `pyReadOnly = true`, tetapi sel yang SAMA juga
memuat `pyReadOnlyCondition = TreatyIn.ViewState = 1` — dan syarat itulah
yang berlaku. Jadi selnya aktif di mode Edit, terkunci di mode lihat.

Tiga hal membenarkannya, dan ketiganya lebih kuat daripada pembacaan tag
telanjang:
- tangkapan layar Pega memperlihatkan sel itu sebagai isian aktif
  (placeholder `Text` / `0,00`, panah autocomplete);
- `T_TREATY_RETRO_SHARE` dan `T_TREATY_FAC_REINSURER` berisi 35 baris
  TERISI — seseorang pernah mengisinya;
- `TreatyInNonAddItem(sharereins/sharefacname)` memang hanya menambah baris
  ber-`ID` kosong, jadi Add/Delete tidak menuntut aksi services sama sekali.

⚠️ Kekeliruannya satu jenis dengan §3.1: membaca SATU tag dan berhenti,
padahal tag kedua di sel yang sama membatalkannya. `pyReadOnly` telanjang
bukan jawaban sampai `pyReadOnlyCondition` ikut dibaca.

### 1.3 Panel `RNM Share`

| Bagian | Properti | Di aplikasi |
|---|---|---|
| `Share to RNM :` + `%` | `TreatyIn.RnmShareDeducted` (baca saja) | ⛔ nol |
| grid `RNM Share` | `.LayerType` · `.Layer` · `.LayerPartType` · `.LayerPart` · `RnmLimitListDisplay(1..2)` · `RnmGrossPremiDisplay(1..2)` · `.RNMShare` | ⚠️ ada, TETAPI kolomnya lain |
| grid `Summarry of RNM Share` *(ejaan Pega)* | `LimitShareSummaryList` | ⛔ nol |

⛔ **Grid `RNM Share` aplikasi memakai kolom yang BERBEDA.** `KOLOM_RNM_SHARE`
hari ini berisi nama kunci dokumen warisan (`LAYER`, `RNMSHARE`,
`LIABILITY_RNM`, `MDP_RNM_100`, `EPIRNMQS100`, `RNM_RETAINED_PREMI`,
`RNM_QS_PREMI`) — itu kolom tabel pendaratan, bukan kolom layar Pega.
Add/Delete-nya `1=2` di ekspor, dan pencabutannya di aplikasi **benar**.

### 1.4 `Total All Layers RNM Share` — sembilan grid

Urutan layar, dan larik yang mengisinya (`pyPageListProperty`):

| # | Judul grid | Larik | Diisi oleh |
|---|---|---|---|
| 1 | `Total RNM Limit (RNM Share)` | `TotalShareRnmNP` | `NPSetTotal(share)` |
| 2 | `Total OR Limit` | `TotalSpreadedRnmProp` | `XOLAddSpreading` |
| 3 | `Total R/I Limit` | `TotalSpreadedRnmRIProp` | `XOLAddSpreading` |
| 4 | `Total Gross Min Premium` | `TotalShareGrossMinNP` | ⛔ **nol langkah** — §3.1 |
| 5 | `Total Gross Premium (MDP)` | `TotalShareGrossNP` | `NPSetTotal(share)` |
| 6 | `Total Deduction` | `TotalShareDeductionNP` | `NPSetTotal(share)` |
| 7 | `Total Net Premium` | `TotalShareNetNP` | `NPSetTotal(share)` |
| 8 | `Total OR Net Premium` | `TotalSpreadedNetPremi` | `XOLAddSpreading` / `SetBrokerage` |
| 9 | `Total R/I Net Premium` | `TotalSpreadedNetPremiRI` | `XOLAddSpreading` / `SetBrokerage` |

Tombolnya: **`Update Total`** (`NPSetTotal(share)` + `SummaryLimitShare` +
`SummaryLimitFacShare`) dan **`Update Value in Share`**
(`TreatyInShareListValue`, syarat `ViewState != '1' && TreatyMasterInEDM`).

⛔ Empat tombol lain di ekspor MATI dan tidak dibangun: dua
`TreatyInNonSetTotal` (`1=2`), `Show Facultative Share (unused)` (`NEVER`),
`Show Share From Other Retro` (`1=2`).

---

## 2 · Keadaan kode saat perbandingan dibuat

⭐ **Backend-nya SUDAH ADA dan jauh lebih lengkap daripada layarnya** —
`services/hitung_share_np.go` (838 baris, 13 aksi), `share_np_muat.go`,
`models/share_np.go`, `repository/pendaratan_share.go`. Jalur bacanya pun
tersambung: `warisan_kontrak.go` mengisi `k.ShareNP`.

⛔ Yang kosong **layarnya**: `TabShare.tsx` masih 126 baris grid hanya-baca
dan nol memakai `shareNP`.

---

## 3 · ⛔ Temuan

### 3.1 `pyStepsBlockName == "//"` — tanda langkah DIKOMENTARI

Di `TreatyInNPSetTotal.xml`, antara langkah yang mengosongkan total dan
langkah yang mengisinya:

```
[22] Property-Remove  blockName `Share`     praAktif true   -> jalan
[23] (blok)           blockName `//`        praAktif false  -> DIKOMENTARI
[24] (blok)           blockName ``          praAktif true   -> jalan
[25] Property-Remove  blockName `FShare`    praAktif true   -> jalan
[26] (blok)           blockName ``          praAktif false  -> JALAN
```

⛔ **`praAktif=false` BUKAN tanda mati.** [26] membuktikannya: ia
`praAktif=false` dan tetap berjalan — ia yang mengisi `TotalFacShare*`.
`praAktif=false` hanya berarti "tanpa pra-syarat". Yang mematikan langkah
adalah `pyStepsBlockName == "//"`.

⚠️ Ronde ini SEMPAT keliru di titik itu: [23] ditambahkan ke
`NPSetTotalShare` atas dasar pembacaan `praAktif=false`, lalu dikembalikan
sesudah tanda `//` ditunjukkan dan diukur ulang. Dicatat di sini dan di
kepala `NPSetTotalShare` supaya tidak terulang.

⭐ Akibatnya: `TotalShareGrossMinNP` **tidak pernah terisi**, sebab [23]
satu-satunya langkah yang menjumlah `GrossPremiumMinList`. Grid
`Total Gross Min Premium` memang selalu kosong — di Pega maupun di sini.

⭐ Ini juga melengkapi catatan skill `baca-xml`, yang menyatakan medan penanda
langkah ter-remark **belum terkonfirmasi**. Sekarang terkonfirmasi.

### 3.2 Dua tanda mirip, dua akibat berlawanan

| Tanda | Contoh | Akibat |
|---|---|---|
| `pyStepsActivityName` kosong | `TotalEgnpi` [1.4.2] | langkah nol metode → mati |
| `pyStepsBlockName == "//"` | `NPSetTotal` [23] | dikomentari → mati |
| `pyStepsPreCondition == false` | `NPSetTotal` [26] | **tetap jalan** |

### 3.3 `SummaryLimitShare` mengodekan mata uang secara KERAS

`@If(local.currency=="IDR", …)` dan `@If(local.currency=="USD", …)`. Mata
uang ketiga **hilang tanpa jejak** dari grid `Summarry of RNM Share`. Judul
kolomnya pun menyebutkannya: `100% Limit (IDR)`, `100% Limit (USD)`, dst.

⭐ Dan baris yang `RnmLimitList`-nya kosong **nol masuk ringkasan** — hanya
langkah 4.3.3 yang menambah baris, dan ia di dalam loop `RnmLimitList`.

### 3.4 `SetBrokerage` membuang SATU baris deduksi, bukan semua

Langkah 6 menyimpan indeks baris ber-`Deduction < 1` ke `local.subscriptdelete`
di dalam loop, lalu menghapus **sekali** sesudah loop. Jadi yang terbuang
hanya yang TERAKHIR ditemukan. Total dan Net tidak dihitung ulang sesudahnya.

### 3.5 `XOLAddSpreading` membaca basis data

Ia memanggil `BrowseTreatyArrangement_ParentReinsMasterTrt` atas
`ASM-FW-GISFW-Int-PROPORTIONALARRG` dan `BrowseTreatyArrangement_Limit_RD`.
Jadi tab ini **tidak dapat** seluruhnya menjadi rute `/hitung/` nol-baca —
ia membaca tabel master. Membaca, bukan menulis: Aturan B tetap aman.

### 3.6 `idKontrakKhusus = "1000951"` — kode keras di Activity

`FetchQSfromMasterXOL` dan `TreatyInSetBrokerage` memperlakukan satu kontrak
secara khusus (faktor 85 %). Sudah terbawa di `hitung_share_np.go`.

---

## 4 · Yang diperiksa dan TIDAK ditemukan

- `TotalFacShare*` ([25]/[26]) dihitung Pega tetapi **nol grid di tab ini
  menampilkannya** — kesembilan `pyPageListProperty` menunjuk `TotalShare*`
  dan `TotalSpreaded*` saja.

---

## 5 · Uji yang ditambahkan ronde ini

`services/hitung_share_np_test.go` — 15 uji; sebelumnya **nol uji** untuk 838
baris rumus. Yang dipaku: `NPSetTotalShare` hanya tingkat spreading,
`TotalShareGrossMinNP` selalu kosong, saringan `SpreadingTypeXOL`, mata uang
kosong ditolak, `SummaryLimitShare` (IDR/USD keras, baris tanpa limit nol
masuk, format `Note`), `NonAddItemShare` (RNMShare nol → pesan;
`FacultativeShare` → `RnmShareDeducted`), kesembilan kunci Total selalu ada,
nol larik nil.

---

## 6 · ⛔ Dua sesi mengerjakan tab yang sama

7 Oktober 2026 sesi `nusantara-re-app-16` menyatakan memegang tab Share
Non-Prop **seluruhnya, backend dan frontend**, atas permintaan pemakai di
sesinya. Sesi ini berhenti menyunting kode Share sesudah itu.

⚠️ Satu kerusakan sudah terjadi sebelum koordinasi: `hitung_share_np_test.go`
milik sesi itu **tertimpa** oleh berkas uji sesi ini. Repo bukan git, jadi
isinya tidak dapat dipulihkan.

⭐ Berkas perbandingan ini sengaja dipisahkan dari kode: ia hasil pembacaan
ekspor, dan berguna bagi siapa pun yang membangun tabnya.
