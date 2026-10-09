# Struktur Tabel — Claim Fac In
## Sembilan tabel lini FAC, lima berkunci ganda, dan sisi komite tanpa tabel baru

> ⭐ `[keputusan work owner]` **2026-09-20.** **Awalan tabel klaim non-life yang lama —
> berhuruf `P` sesudah `CLAIM` — diganti menjadi `T_CLAIM_`**, sebab tabelnya kini **dipakai
> bersama lini FAC dan PROP** sehingga huruf **P** pada awalan menyesatkan.
> ⭐ **NONPROP nanti ikut tabel yang sama.**
>
> ⛔ **Nama awalan lamanya sengaja TIDAK dikutip harfiah** di berkas mana pun di luar lini
> Life — ⭐ supaya pencarian atas awalan lama itu **hanya** menemukan berkas yang memang belum
> diselaraskan, bukan kalimat yang menerangkan penggantiannya.
>
> ⭐ **Tujuh tabel dipakai bersama; lima di antaranya berkunci asing GANDA.**
> ⛔ **Tabel lini Life tidak disentuh** — awalannya berbeda dan tidak ikut terganti.
>
> ⚠️ Akibatnya **dua berkas lini Life masih menyebut nama tabel penyesuaian yang lama**.
> ⭐ Itu **dicatat sebagai `[terbuka]`**, ⛔ **bukan diperbaiki**.
>
> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Tujuh tabel dipakai bersama; lima di
> antaranya berkunci asing GANDA.**"* → **kunci asing ganda DITOLAK** work owner (OQ-CFI-01, 09-10-2026): dua tabel
> **baru** (`T_CLAIM_OBJECT` 561, `T_CLAIM_OBJECT_ITEM` 562), tabel `T_CLAIM_*` Claim Prop dipakai ulang dengan
> `ALTER … ADD`, baris FAC mengisi `CLAIM_ID` **dan** induk tingkatnya; nol `MODIFY` / `DROP`. Rincian di RALAT §2.

> **SENSUS BERKAS INI** — *(dihitung ulang sesudah RALAT + putaran kolom, 2026-09-20)*
>
> ⛔ **Jendelanya disebut: seluruh berkas DIKURANGI blok sensus ini sendiri** — **695 baris**
> menurut cacah baris-baru, **696** menurut pemisahan teks. ⚠️ **Selisih satu, sebabnya diketahui:**
> berkas berakhir dengan baris baru. ⭐ **Angka yang dipakai: 695.**
>
> bab `## ` **11** · tabel **32** · ⭐ **kolom 138** · ⚠️ **tanpa penulis 35** ·
> `[terverifikasi]` **6** · `[keputusan work owner]` **9** · `[terbuka]` **12** ·
> `[work owner]` **6** · ⚠️ **63** · ⭐ **183** · ⛔ **143** · ✅ **19** · ★ **4**
>
> ⭐ **Dihitung DUA CARA** — *(a)* teks utuh, *(b)* baris demi baris. ✅ **Sepakat pada kesebelas
> angka penanda.**
>
> ⛔ ⭐ **Awalan tabel yang lama dan awalan lini Life: nol kemunculan** di berkas ini.
>
> ⚠️ **Bergeser dari keadaan sebelumnya** — badan **308 → 695**, sebab bab kolom ditambahkan.
>
> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Angka yang dipakai: 695.**"* → sensus ini
> **tidak dihitung ulang**: ia mengukur badan berkas 20-09-2026. Blok RALAT 10-10-2026 dan lampiran pengikat
> 10-10-2026 di akhir berkas **berada di luar jendelanya**.

---

## Cara membaca berkas ini

Berkas ini memerikan **bentuk tabel dan asalnya**. ⛔ **Bukan DDL, bukan tipe data, bukan daftar
kolom** — daftar kolom adalah putaran berikutnya.

**Penanda:**

| Penanda | Artinya |
| --- | --- |
| `[terverifikasi]` | Terbaca langsung dari korpus ekspor |
| `[keputusan work owner]` | Diputuskan work owner, bertanggal |
| `[terbuka]` | Belum diketahui — ⛔ **tidak ditebak** |
| ★ | Khas lini FAC — tidak ada padanannya di lini PROP |
| ⭐ ⚠️ ⛔ | menentukan · mudah salah · larangan |

### ⭐⭐ Aturan baca gerbang Pega — **wajib, dan baru pertama kali tertulis**

⛔ **Aturan ini sudah dipakai di seluruh ronde grilling, tetapi belum pernah ditulis. Ia ditulis di
sini supaya dapat diperiksa orang lain.**

| Medan | Artinya |
| --- | --- |
| ⭐ `pyStepsPreCondParamsWhenFalse = 3` | ⭐ **Langkah jalan bila gerbang BENAR.** Nilai `3` berarti *lewati langkah*, dan ia dipasang pada sisi **salah** — jadi sisi **benar** yang meneruskan |
| ⭐ `pyStepsPreCondParamsWhenTrue = 3` | ⭐ **Langkah jalan bila gerbang SALAH.** ⚠️ **Kebalikannya** — dan inilah bentuk yang paling sering salah dibaca |

⚠️ **Dua aturan turunan yang menyertainya**, dipakai di seluruh berkas ini:

| # | Aturan |
| --- | --- |
| ⭐ **1** | Bendera **gerbang mati** mematikan **GERBANG**, bukan **LANGKAH** — ⚠️ langkahnya justru berjalan **lebih sering** |
| ⭐ **2** | Kode arah **5** mengubah rantai **DAN** menjadi **ATAU** |

⛔ **Rujukan korpus memakai nama berkas + nomor langkah.** ⛔ **Nol nomor baris XML** — ia berubah
setiap ekspor ulang.

### Letak enam butir wajib

| Butir | Bab |
| --- | --- |
| **a** — sembilan tabel FAC + sumber Pega | **§1** |
| **b** — lima tabel berkunci ganda + `CHECK` + preseden | **§2** |
| **c** — sebelas halaman yang tidak dibawa | **§3** |
| **d** — aturan keseragaman berbagi per treaty | **§4** |
| **e** — aturan baca gerbang | **bab ini**, di atas |
| **f** — register `[terbuka]` | **§6** |
| ⭐ **kolom tiap tabel** | ⭐ **§2b** |

---

## §0 — Keputusan pokok

`[keputusan work owner]` **2026-09-20** — empat keputusan yang mengikat seluruh berkas ini.

| # | Keputusan | Sebabnya |
| --- | --- | --- |
| ⭐ **K1** | **Awalan tabel klaim non-life diganti menjadi `T_CLAIM_`** | ⭐ tabelnya kini dipakai bersama lini **FAC** dan **PROP**, jadi awalan berhuruf **P** menyesatkan. ⭐ **NONPROP nanti ikut tabel yang sama** |
| ⭐ **K2** | **Tujuh tabel dipakai bersama FAC dan PROP.** ⭐ **Lima** di antaranya **induknya berbeda antar lini**, jadi **berkunci asing GANDA** | bentuk data kedua lini berbeda tingkat |
| ⛔ **K3** | **Tabel dan berkas lini Life TIDAK disentuh sama sekali** | ⚠️ Akibatnya **dua berkas Life masih menyebut nama tabel penyesuaian yang lama**. ⭐ **Dicatat sebagai `[terbuka]`**, ⛔ **bukan diperbaiki** |
| ⭐ **K4** | **Modul Komite tidak punya tabel baru** | `T_GENERAL_KOMITE` dan `T_KOMITE_KOMITELIST` **sudah terkunci sebagai tabel lintas-lini**; ⭐ lini FAC memakainya **apa adanya** |

> ⛔ **RALAT 10-10-2026 atas K2.** Kalimat lamanya dikutip utuh, tidak dihapus: *"**Tujuh tabel dipakai bersama FAC dan
> PROP.** ⭐ **Lima** di antaranya **induknya berbeda antar lini**, jadi **berkunci asing GANDA**"* → **ditolak** work
> owner, OQ-CFI-01 (09-10-2026): tabel `T_CLAIM_*` Claim Prop dipakai ulang lewat `ALTER … ADD` (`563`–`567`), baris
> FAC mengisi `CLAIM_ID` **dan** `OBJECT_ITEM_ID` / `ADJUSTMENT_ID`; nol kolom induk nullable, nol `CHECK`, nol
> `MODIFY`. Rincian di RALAT §2. K1, K3, K4 tidak berubah.

### Tabel lintas-lini — ⛔ sudah terkunci, tidak didefinisikan ulang di sini

`T_WORK_CLAIM` · `T_GENERAL_CLAIM` · `T_GENERAL_KOMITE` · `T_KOMITE_KOMITELIST` ·
`T_VIEW_SUGGEST`

### Awalan pengenal

| Lini | Baris klaim | Baris komite | Penanda lini |
| --- | --- | --- | --- |
| ⭐ **FAC** | `CLM-` | `KMT-` | `LINI = FAC` |
| **PROP** | `CLMP-` | `TKMT-` | `LINI = PROP` |

> ⛔ **RALAT 10-10-2026.** Sel lamanya dikutip utuh, tidak dihapus: *"`LINI = FAC`"* → nilai yang dibangun
> **`T_WORK_CLAIM.LINI = 'FACIN'`** — sama dengan `STS_KLAIM` `EMAILKOMITE` FACIN. Awalan `CLM-` / `KMT-` sama dengan
> Pega, nomor dari `SEQ_WORK_CLAIM`; lini **disaring lewat `LINI`, tidak pernah lewat awalan**; `TAHAP` = nama FlowAction
> (`InputRegister` / `InputEstimasi` / `InputSurveyor`) — OQ-CFI-02, `PARITAS.md` §10.

---

## §1 — Sembilan tabel lini FAC, dan asalnya di Pega

⭐ **Sembilan tabel dipakai lini FAC.** ⛔ Tiga tabel lain pada keluarga `T_CLAIM_` **hanya dipakai
lini PROP** dan **tidak digambar sebagai milik FAC** — lihat catatan di bawah tabel.

| # | Tabel | Sumber Pega — lini **FAC** | Sumber Pega — lini PROP | |
| --- | --- | --- | --- | --- |
| **1** | `T_CLAIM_OBJECT` | `ClaimData.ObjectList` | ⛔ tidak ada | ★ **khas FAC** |
| **2** | `T_CLAIM_OBJECT_ITEM` | `ObjectList[].ObjectItemList` | ⛔ tidak ada | ★ **khas FAC** |
| **3** | `T_CLAIM_ESTIMATION` | `ObjectItemList[].EstimationList` | `ClaimData.EstimationList` | ⭐ kunci ganda |
| **4** | `T_CLAIM_SPREADING` | `ObjectItemList[].SpreadingClaim` | `ClaimData.SpreadingClaim` | ⭐ kunci ganda |
| **5** | `T_CLAIM_BREAK_QS` | `ObjectItemList[].SpreadingAdjustment` | `ClaimData.SpreadingBreakQS` | ⭐ kunci ganda |
| **6** | `T_CLAIM_ADJUSTMENT` | `ObjectItemList[].Adjustment` **+** `DataCommitteFacin` | `ClaimData.AdjustmentList` **+** `DataCommitteeTreaty` | ⭐ kunci ganda |
| **7** | `T_CLAIM_ADJ_SPREADING` | `Adjustment[].SpreadingAdjustment` | `AdjustmentList[].SpreadingAdjustment` | kunci **sama** dua lini |
| **8** | `T_CLAIM_ADJ_QUOTA_SHARE` | `Adjustment[].SpreadingQuotaShare` | `AdjustmentList[].SpreadingQuotaShare` | kunci **sama** dua lini |
| **9** | `T_CLAIM_FAC_RETRO` | `Adjustment[].FacRetroList` | `ClaimData.FacRetroList` | ⚠️ kunci ganda, **beda TINGKAT** |

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Sembilan tabel dipakai lini FAC.**"* →
> sembilan tabel itu benar, tetapi **bentuk "kunci ganda" pada kolom terakhir tidak dibangun** (OQ-CFI-01). Yang
> dibangun 10-10-2026 (`MODUL.md` bab Migrasi, `backend/models/katalog_tabel.go`):
>
> | Migrasi | Tabel | Perubahan |
> | --- | --- | --- |
> | `560` | `T_GENERAL_CLAIM` | **19 kolom kepala FAC** `ALTER … ADD` nullable |
> | `561` | `T_CLAIM_OBJECT` | **tabel baru**; `CLAIM_ID` → `T_GENERAL_CLAIM` ON DELETE CASCADE |
> | `562` | `T_CLAIM_OBJECT_ITEM` | **tabel baru**; `OBJECT_ID` → `T_CLAIM_OBJECT` CASCADE, `CLAIM_ID` ikut diisi (CASCADE) |
> | `563`–`566` | `T_CLAIM_ESTIMATION` · `T_CLAIM_SPREADING` (+ `JENIS`) · `T_CLAIM_BREAK_QS` · `T_CLAIM_ADJUSTMENT` | `ADD OBJECT_ITEM_ID` nullable, FK CASCADE ke item, + kolom khas FAC nullable |
> | `567` | `T_CLAIM_FAC_RETRO` | `ADD ADJUSTMENT_ID` nullable, FK CASCADE ke adjustment |
>
> Baris FAC mengisi **`CLAIM_ID` dan `OBJECT_ITEM_ID`** sekaligus (`CLAIM_ID` tetap NOT NULL seperti DDL Claim Prop);
> `T_CLAIM_FAC_RETRO` FAC mengisi `CLAIM_ID` + `ADJUSTMENT_ID`, baris retro tingkat klaim (`ClaimData.FacRetroTreaty`)
> ber-`ADJUSTMENT_ID` kosong. `T_CLAIM_ADJ_SPREADING` / `T_CLAIM_ADJ_QUOTA_SHARE` dipakai tanpa perubahan. **Nol
> `MODIFY` / `DROP`.**

### ⛔ Tiga tabel yang HANYA dipakai lini PROP — catatan, bukan milik FAC

| Tabel | Sumber PROP | Lini FAC |
| --- | --- | --- |
| `T_CLAIM_INTEREST` | `ClaimData.InterestList` | ⛔ **tidak ada** |
| `T_CLAIM_CLAIM_AMOUNT` | `ClaimData.ListClaimAmount` | ⛔ **tidak ada** |
| `T_CLAIM_ADJ_LOSS_ALLOCATION` | `AdjustmentList[].LossAllocation` | ⛔ **tidak ada** |

⭐ **Ketiganya disebut di sini hanya supaya pembaca tahu keluarga `T_CLAIM_` lebih besar daripada
sembilan** — ⛔ **tidak digambar pada pohon lini FAC.**

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Ketiganya disebut di sini hanya supaya
> pembaca tahu keluarga `T_CLAIM_` lebih besar daripada sembilan**"* → keluarga itu **lebih besar lagi**:
> `T_CLAIM_LOSS_ALLOCATION` (Claim Prop `524`) tidak disebut, dan Claim Non Prop menambah `T_CLAIM_NP_LOSS_ALLOC` ·
> `T_CLAIM_NP_XOL_ALLOC` · `T_CLAIM_NP_CLAIM_ACCEPT` (`608`–`610`) serta kolom NONPROP di tabel bersama (`600`–`607`).
> ⛔ **Tak satu pun dipakai lini FAC** — katalog `backend/models/katalog_tabel.go` tidak menyebutnya.

### ⭐ Kronologi klaim — `T_VIEW_SUGGEST`

> ⛔⛔ **RALAT — 2026-09-20.** Bab ini sebelumnya berjudul *"⚠️ Tabel ketiga belas — jejak
> audit"* dan berisi tabel bernama **`T_CLAIM_AUDIT_TRAIL`** dengan keterangan
> *"⚠️ belum ditelusuri"* dan *"⚠️ namanya belum ditetapkan"*.
>
> ⛔ **Tabel itu TIDAK ADA.** `[keputusan work owner]` **2026-09-20**.
> ⭐ **Yang ada adalah `T_VIEW_SUGGEST`, dan isinya diambil dari KRONOLOGI.**
>
> ⭐ Kalimat lama dikutip di sini, ⛔ tidak dihapus.

`[terverifikasi]` `T_VIEW_SUGGEST` adalah **daftar kronologi klaim** — tabel **lintas-lini** yang
sudah terkunci.

| | |
| --- | --- |
| Sumber Pega | ⭐ **`ClaimData.SuggestList`** |
| Penulisnya | `Activity\SethistoryKlaimTreaty` langkah **1** · `DataTransform\InsertChronology_DT` |
| Induknya | `T_GENERAL_CLAIM` lewat `CLAIM_ID` |
| ⚠️ Perilaku hapus | ⛔⛔ **JANGAN cascade** |

⚠️⚠️ **Sebab "jangan cascade" wajib dibaca:** ⭐ kronologi adalah **jejak**, dan jejak yang ikut
terhapus bersama induknya **berhenti menjadi jejak**. ⛔ Menghapus klaim **tidak boleh** menghapus
riwayat siapa mengerjakan apa.

⭐ **Lini FAC memakainya** — ia lintas-lini, sama seperti `T_WORK_CLAIM` dan `T_GENERAL_CLAIM`.
⛔ **Bukan tabel baru**, jadi **tidak didefinisikan ulang di sini**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⚠️⚠️ **Sebab "jangan cascade" wajib dibaca:**
> ⭐ kronologi adalah **jejak**, dan jejak yang ikut terhapus bersama induknya **berhenti menjadi jejak**."* → perilaku
> hapusnya **CASCADE**: migrasi Claim Prop `532` membuat `T_VIEW_SUGGEST.CLAIM_ID` → `T_GENERAL_CLAIM` **ON DELETE
> CASCADE** (`FK_VS_CLAIM`, dengan `CHECK` satu induk `CK_VS_SATU_INDUK`), dan Claim Fac In memakainya apa adanya.
> Modul ini tidak menghapus baris `T_GENERAL_CLAIM` (nol `DELETE` atasnya di `backend/repository/`; penutupan hanya
> mengisi `STATUS_WORK`). Pemetaan medan kronologi → kolom bersama: OQ-CFI-19.
>
> ⛔ **RALAT 10-10-2026 — rencana tulisan ke tabel warisan yang tidak pernah tertulis di berkas ini.** Bab §1 hanya
> mendaftar tabel `T_CLAIM_*`; yang juga **ditulis** Claim Fac In (dibangun 10-10-2026, `repository/tulisan.go`,
> `repository/lampiran.go`, `PARITAS.md`), semuanya **di dalam satu transaksi per aksi**, nol procedure, nol `COMMIT`:
>
> - `OS_AKSEPTASI_KLAIM` — CFS: baris **STS 0** per estimasi; Close Claim: **STS 4**; DLA: `UPDATE DLA_NO, DLA_DATE`
>   menurut `AcceptedNo` (`InsertDLA_OS_SQL`); `DATA_JSON` diisi format Pega `GetPageJSONString`. **STS 1** = akseptasi
>   komite = tahap 2 (`komiteclaimfacin`).
> - `JSON_KLAIM` — baris `IDPEGA` yang ada dibiarkan, selainnya `INSERT`, **tanpa** `DATA_JSON`.
> - `PROGRESSCLAIM` / `SUBPROGRESSCLAIM` — logika procedure `PEGA_PROGRESSCLAIM` / `PEGA_SUBPROGRESSCLAIM` ditulis
>   ulang, procedure tidak dipanggil.
> - `CATASTROPHE` (`SaveCatasrtope_Act`); `DOCUMENT_CLAIM` + `T_STORAGE_IMAGE` (lampiran, `inti/backend/penyimpanan`);
>   `MONITORING_KLAIM_LOG`; kronologi di `T_VIEW_SUGGEST`.
> - Efek keluar hanya lewat outbox `T_LOG_SERVICE_RNM`, dan hanya di produksi (`IS_PEGA_PROD`).

---

## §2 — Kunci asing: lima tabel berkunci GANDA

⭐ **Sebabnya satu kalimat:** pada lini **FAC**, data menggantung di **item objek**; pada lini
**PROP**, data yang sama menggantung langsung di **klaim**. ⭐ Tabelnya satu; induknya dua.

| Tabel | Induk **FAC** | Induk **PROP** |
| --- | --- | --- |
| `T_CLAIM_ESTIMATION` | `OBJECT_ITEM_ID` | `CLAIM_ID` |
| `T_CLAIM_SPREADING` | `OBJECT_ITEM_ID` | `CLAIM_ID` |
| `T_CLAIM_BREAK_QS` | `OBJECT_ITEM_ID` | `CLAIM_ID` |
| `T_CLAIM_ADJUSTMENT` | `OBJECT_ITEM_ID` | `CLAIM_ID` |
| ⚠️ `T_CLAIM_FAC_RETRO` | **`ADJUSTMENT_ID`** | `CLAIM_ID` |

⚠️⚠️ **`T_CLAIM_FAC_RETRO` berbeda dari empat lainnya:** induknya bukan hanya **tabel** yang
berbeda, melainkan **TINGKAT** yang berbeda — di FAC ia menggantung pada **penyesuaian**, di PROP
pada **klaim**.

### Aturan penjaga — tepat satu induk terisi

Untuk keempat tabel pertama:

> `CHECK ( (OBJECT_ITEM_ID IS NULL) <> (CLAIM_ID IS NULL) )`

Untuk `T_CLAIM_FAC_RETRO`:

> `CHECK ( (ADJUSTMENT_ID IS NULL) <> (CLAIM_ID IS NULL) )`

⭐ **Bacaannya:** tepat **satu** dari dua kolom induk terisi — ⛔ tidak boleh keduanya, tidak boleh
tak satu pun.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Bacaannya:** tepat **satu** dari dua kolom
> induk terisi — ⛔ tidak boleh keduanya, tidak boleh tak satu pun."* → **bentuk ini DITOLAK** keputusan work owner
> **OQ-CFI-01 (09-10-2026)**: ⛔ nol `CLAIM_ID` nullable, nol `CHECK` "tepat satu induk", nol `MODIFY`. Yang dibangun:
>
> - **Tabel baru** `T_CLAIM_OBJECT` (migrasi `561`, `CLAIM_ID` → `T_GENERAL_CLAIM` CASCADE) dan `T_CLAIM_OBJECT_ITEM`
>   (`562`, `OBJECT_ID` → `T_CLAIM_OBJECT` CASCADE, `CLAIM_ID` ikut diisi).
> - Tabel Claim Prop `T_CLAIM_ESTIMATION` · `T_CLAIM_SPREADING` (+ kolom `JENIS`) · `T_CLAIM_BREAK_QS` ·
>   `T_CLAIM_ADJUSTMENT` dipakai ulang dengan `ALTER … ADD OBJECT_ITEM_ID` (nullable, FK CASCADE ke item; `563`–`566`).
> - Baris FAC mengisi **KEDUANYA** — `CLAIM_ID` (tetap NOT NULL) **dan** `OBJECT_ITEM_ID`; baris Prop / Non Prop
>   ber-`OBJECT_ITEM_ID` kosong. Karena `UNIQUE (CLAIM_ID, NOURUT)` tabel Claim Prop, `NOURUT` baris FAC berurut
>   **se-klaim** (`backend/models/katalog.go`).
> - `T_CLAIM_FAC_RETRO` + `ADJUSTMENT_ID` (`567`, FK CASCADE ke adjustment); `CLAIM_ID` tetap diisi.
> - 19 kolom kepala FAC `ADD` ke `T_GENERAL_CLAIM` (`560`). **Nol `MODIFY` / `DROP`.**
>
> Kolom yang berlaku: **lampiran pengikat** di akhir berkas ini (tabel baru) dan lampiran pengikat
> `modul/claimprop/docs/STRUKTUR-TABEL-CLAIM-PROP.md` (tabel bersama).

### ⭐ Polanya BUKAN hal baru

`[keputusan work owner]` **2026-09-18** — `T_VIEW_SUGGEST` **sudah memakai pola yang sama**:
dua induk **nullable**, `PREMIUM_LIST_ID` dan `CLAIM_ID`, dengan `CHECK` yang menuntut **tepat satu
terisi**.

⭐ **Jadi ini bukan bentuk baru yang diperkenalkan lini FAC** — ia **pola yang sudah berlaku** di
basis data ini, dipakai ulang.

### Dua tabel yang kuncinya SAMA di kedua lini

| Tabel | Kunci | Penjaga |
| --- | --- | --- |
| `T_CLAIM_ADJ_SPREADING` | `ADJUSTMENT_ID` | ⛔ **satu kolom saja, tanpa `CHECK`** |
| `T_CLAIM_ADJ_QUOTA_SHARE` | `ADJUSTMENT_ID` | ⛔ **satu kolom saja, tanpa `CHECK`** |

⭐ **Sebabnya:** keduanya menggantung pada **penyesuaian** di **kedua** lini — jadi tidak ada
percabangan induk.

---

## §2b — ⭐⭐ Kolom tiap tabel

⛔ **Nol `CREATE TABLE`, nol DDL.** ⭐ Jenis ditulis dalam **kata**: *teks · bilangan · DATE ·
CHAR(1)*. ⛔ Panjang dan presisi **belum ditetapkan** — itu putaran DDL.

⭐ **Tiap kolom menyebut sumbernya**: medan Pega **dan** berkas + nomor langkah tempat ia **ditulis**.
⚠️ Kolom yang **tidak ditemukan penulisnya** ditandai **⛔ penulis tak ketemu** — ⭐ itu **bukan**
alasan membuangnya, sebab medannya **terbaca dipakai**; ia **butir `[terbuka]`**, sejalan butir 1.

⛔ **Kolom teknis yang TIDAK dibawa**, berlaku untuk semua tabel: `pxListSubscript` ·
`pxCreateOperator` · `pxCreateOpName` — ⭐ ketiganya **bawaan Pega**, digantikan kolom jejak sistem
baru.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⛔ Panjang dan presisi **belum ditetapkan** —
> itu putaran DDL."* → DDL sudah ditulis (`560`–`567`). **Uang, persen, share, dan kurs `NUMBER(38,10)`**; penanda =
> `VARCHAR2(16)` teks (bukan `CHAR(1)`: `""` berbeda dari `"0"`); `ID` / kolom induk = `VARCHAR2(32)` teks, ID stabil dari
> `SEQ_T_CLAIM`, + `NOURUT`. **Nama kolom mengikuti DDL Claim Prop `521`–`530`**, bukan nama di bab ini; satu-satunya
> pemetaan properti → kolom = `backend/models/katalog_tabel.go`, kolom yang berlaku = lampiran pengikat. Nama yang paling
> sering berbeda:
>
> | Nama di bab ini | Nama DDL yang dibangun |
> | --- | --- |
> | `CURRENCY` (item, estimasi, spreading, Break QS, adjustment) | `CURRENCY_NAME` (+ `CURRENCY_ID`) |
> | `TREATY_TYPE` (spreading, Break QS, spreading / QS adjustment) | `TREATY_ID` — isinya `TreatyType` Pega |
> | `GROSS_ESTIMATION_PCT` · `CONVERT_VALUE` · `ESTIMASI_MORE_THAN_TSI` · `GROSS_ESTIMATION_PCT_MBU` | `GROSS_ESTIMATION_VALUE` · `ESTIMATION_VALUE_IDR` · `IS_MORE_THAN_TSI` · `GROSS_ESTIMATION_MBU` |
> | `VALUE_ADJUSTMENT` · `DIRECT_TO_KASIR` | `VALUE_ADJUSTMENT_IDR` · `IS_DIRECT_TO_KASIR` |
> | `PCT_SHARE_ALL_OBJ` · `RI_COMM_ALL_OBJ` · `TOTAL_ESTIMASI_REAS` (retro) | `SHARE_PCT` · `RI_COMMISSION_PCT` · `TOTAL_ESTIMATION_REINS` |
> | `OBJECT_ID` (teks, `ObjectID` polis) di `T_CLAIM_OBJECT` | `OBJECT_REF`; letak objek di polis = `KUNCI_POLIS` |
>
> ⚠️ Dan kalimat di atas soal `pxCreateOperator` / `pxCreateOpName` **tidak berlaku untuk dua tabel**: adjustment
> menyimpannya sebagai `CREATED_BY` / `CREATED_BY_NAME` / `CREATED_AT` (kolom Claim Prop `528`), estimasi sebagai
> `CREATED_BY` (`563`). Kolom bab ini yang tidak ada di lampiran pengikat — atau, untuk tabel bersama, di lampiran
> pengikat `modul/claimprop/docs/STRUKTUR-TABEL-CLAIM-PROP.md` — **tidak dibuat**.

### T_CLAIM_OBJECT — objek yang diklaim

⭐ **Induk:** `T_GENERAL_CLAIM` lewat `CLAIM_ID`.

| Kolom | Jenis | Wajib | Arti | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bilangan | ya | kunci utama | ⭐ baru — Pega memakai indeks larik |
| `CLAIM_ID` | bilangan | ya | induk klaim | ⭐ baru — Pega memakai sarang |
| `OBJECT_ID` | teks | ya | pengenal objek dari polis | `ObjectID` — ⛔ penulis tak ketemu |
| `OBJECT_NAME` | teks | ya | nama objek | `ObjectName` — `PrintPDFAccep_MultiAksep` lgk **13** |
| `TYPE_NAME` | teks | tidak | jenis objek | `TypeName` — `GetProgresClaim_ACT` lgk **3.2.1** |
| `OBJECT_STATUS` | teks | tidak | status objek | `ObjectStatus` — ⛔ penulis tak ketemu |
| `OBJECT_LOCATION` | teks | tidak | lokasi objek | `ObjectLocation` — ⛔ penulis tak ketemu |
| `SHARE_RETRO` | bilangan | tidak | bagian retro atas objek | `ShareRetro` — `DraftGenerateDLAFacin_Act` lgk **3** |
| `NO_DLA` | teks | tidak | nomor DLA | `NoDLA` — `DLAFacintoTreaty_Act` lgk **9** |
| `DLA_STATUS` | teks | tidak | status DLA | `DLAStatus` — `SaveAcceptation` lgk **13** |
| `PLA_STATUS` | teks | tidak | status PLA | `PlaStatus` — ⛔ penulis tak ketemu |
| `REMARKS_DLA` | teks | tidak | catatan DLA | `RemarksDLA` — `GetProgresClaim_ACT` lgk **3.2.1** |
| `REMARKS_PLA` | teks | tidak | catatan PLA | `RemarksPLA` — `GetProgresClaim_ACT` lgk **3.2.1** |
| `IS_FAC_RETRO` | CHAR(1) | tidak | objek bersifat retro fakultatif | `IsFacretro` — `CheckLimit_Act1` lgk **9** |
| `IS_MORE_THAN_TREATY_LIMIT` | CHAR(1) | tidak | nilai melampaui limit treaty | `IsMoreThanTreatyLimit` — `CheckLimit_Act1` lgk **2** |
| `IS_KOMITE` | CHAR(1) | tidak | objek masuk jalur komite | `IsKomite` — `SetGrossAdjustment_act` lgk **15** |
| `IS_PRINT_ACCEPT` | CHAR(1) | tidak | dokumen akseptasi sudah dicetak | `IsPrintAccept` — `SaveAcceptation` lgk **13** |
| `PRINT_FACE_CLAIM` | CHAR(1) | tidak | lembar muka klaim sudah dicetak | `PrintFaceClaim` — `CheckEstimateValue` lgk **24** |
| `CFS` | teks | tidak | penanda CFS | `CFS` — `CheckEstimateValue` lgk **24** |
| `SAVE_SPREADING` | CHAR(1) | tidak | penanda simpan pembagian | `SaveSpreading` — `SetValueSaveSpreading` lgk **1** |
| `OBJECT_DATE` | DATE | tidak | tanggal objek | `ObjekTanggal` — `GetProgresClaim_ACT` lgk **3.2.1** |
| `FOLLOWUP_DATE` | DATE | tidak | tanggal tindak lanjut | `TanggalFolloup` — `GetProgresClaim_ACT` lgk **3.2.1** |
| `NOTES` | teks | tidak | catatan | `Notes` — `GetProgresClaim_ACT` lgk **3.2.1** |
| `COMMENT` | teks | tidak | komentar | `Comment` — `GetProgresClaim_ACT` lgk **3.2.1** |
| `SURVEY_ID` | teks | tidak | pengenal survei | `ObjectSurveyID` — ⛔ penulis tak ketemu |
| `SURVEY_LOCATION` | teks | tidak | lokasi survei | `ObjectSurveyLocation` — ⛔ penulis tak ketemu |
| `SURVEYOR` | teks | tidak | nama surveyor | `ObjectSurveyor` — ⛔ penulis tak ketemu |

⚠️⚠️ **Dua puluh sembilan medan bergantung JENIS OBJEK, dan semuanya butir `[terbuka]` tersendiri:**

| Kelompok | Medan Pega |
| --- | --- |
| **Orang** | `ObjectIDCard` · `ObjectDateOfBirth` · `ObjectJob` · `Gender` · `ObjectParticipantStatus` |
| **Kendaraan** | `Brand` · `BrandName` · `Model` · `ModelName` · `Type` · `LicensePlate` · `EngineNumber` · `ChassisNumber` · `ObjectVehicleChasis` · `ObjectVehicleType` |
| **Bangunan** | `BuildingCategory` · `BuildingStorey` · `ConstructionClass` |
| **Kapal** | `VesselName` · `VesselType` · `VesselUsage` · `VesselClassification` · `VesselConstruction` · `YearofBuild` · `GrossTonnage` |
| **Ukuran umum** | `ObjectWeight` · `ObjectHeight` · `ObjectLeftHanded` |
| **Lain** | `ASMClassID` |

⛔ **Semuanya bertanda penulis tak ketemu** — ⭐ artinya **dibaca dari modul polis**, bukan ditulis
klaim. ⚠️ **Apakah kesembilan-belas-dan-sepuluh medan itu menjadi kolom, atau dibaca lewat rujukan ke
polis, BELUM DIPUTUSKAN.** ⭐ **Butir `[terbuka]` BARU — butir 11.**

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⚠️ **Apakah kesembilan-belas-dan-sepuluh medan
> itu menjadi kolom, atau dibaca lewat rujukan ke polis, BELUM DIPUTUSKAN.**"* → yang dibangun: **dibaca lewat rujukan
> ke polis**. Medan salinan polis **tidak disimpan**; `T_CLAIM_OBJECT.KUNCI_POLIS` menyimpan letak objek di halaman
> polis (`LocationList(n)` / `VehicleList(n)` / `PersonList(n)`), dan halaman polis dibaca ulang dari `JSON_POLIS` setiap
> kasus dimuat (migrasi `561`, `LengkapiObjek`). Kolom `T_CLAIM_OBJECT` yang berlaku = lampiran pengikat (19 kolom);
> kolom tabel di atas yang tidak ada di sana (mis. `TYPE_NAME`, `OBJECT_STATUS`, `SURVEY_ID`) **tidak dibuat**.

### T_CLAIM_OBJECT_ITEM — item di dalam objek

⭐ **Induk:** `T_CLAIM_OBJECT` lewat `OBJECT_ID`.

| Kolom | Jenis | Wajib | Arti | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bilangan | ya | kunci utama | ⭐ baru |
| `OBJECT_ID` | bilangan | ya | induk objek | ⭐ baru |
| `OBJECT_ITEM_NAME` | teks | ya | nama item | `ObjectItemName` — ⛔ penulis tak ketemu |
| `COVERAGE_ID` | teks | tidak | pengenal jaminan | `CoverageID` — ⛔ penulis tak ketemu |
| `COVERAGE_NOTE` | teks | tidak | catatan jaminan | `CoverageNote` — `CekCoverageNote_Act` lgk **1.2.1** |
| `CURRENCY` | teks | ya | mata uang item | `Currency` — `GetAllData_Act` lgk **19** |
| ⭐ `KURS_OBJECT_ITEM` | bilangan | ya | ⭐ **kurs konversi** — pengganti salinan jaminan polis | `KursObjectItem` — ⛔ penulis tak ketemu |
| ⭐ `VALUE_TSI_NUSARE_IDR` | bilangan | ya | ⭐ **nilai TSI bagian Nusantara Re dalam rupiah** — pengganti salinan jaminan polis | `ValueTSINusareIDR` — ⛔ penulis tak ketemu |
| `TSI_NUSARE` | bilangan | tidak | nilai TSI bagian Nusantara Re | `TSINusare` — ⛔ penulis tak ketemu |
| `TSI_PER_OBJECT` | bilangan | tidak | nilai TSI per objek | `TSIPerObject` — ⛔ penulis tak ketemu |
| `TOTAL_ESTIMASI` | bilangan | tidak | total estimasi item | `TotalEstimasi` — `DLAFacintoTreaty_Act` lgk **13.3.12.2.2** |
| `TOTAL_ESTIMASI_REAS` | bilangan | tidak | total estimasi bagian reasuransi | `TotalEstimasiReas` — `DLAFacintoTreaty_Act` lgk **13.3.13** |
| `TOTAL_GROSS_ESTIMASI` | bilangan | tidak | total estimasi bruto | `TotalGrossEstimasi` — `GeneratePLATreaty_Act` lgk **15** |
| `TOTAL_ESTIMATION_VALUE_IDR` | bilangan | tidak | total estimasi dalam rupiah | `TotalEstimationValueinIDR` — ⛔ penulis tak ketemu |
| `TOTAL_HASIL_CLAIM` | bilangan | tidak | total hasil klaim | `TotalHasilClaim` — `SetSpreadingAjsutement_Act` lgk **13** |
| `ADJUSTMENT_VAL` | bilangan | tidak | nilai penyesuaian item | `AdjustmentVal` — `SetNilaiResikoSendiri` lgk **24** |
| `IS_ADJ_VAL` | CHAR(1) | tidak | penanda nilai penyesuaian terisi | `IsAdjVal` — `SetInitial_ACT` lgk **6** |
| `PAYMENT_TYPE` | bilangan | tidak | jenis pembayaran | `PaymentType` — `SetAdjTypePayment_act` lgk **13** |
| `PERCENT_REAS` | bilangan | tidak | persentase reasuransi | `PercentReas` — `DLAFacintoTreaty_Act` lgk **13.3.13** |
| `PERCENT_HANDLING_FEE` | bilangan | tidak | persentase biaya penanganan | `PercentHandlingFee` — `DLAFacintoTreaty_Act` lgk **13.3.13** |
| `PERCENT_FACOUT_CLAIM` | bilangan | tidak | persentase fac out | `PercentFacoutClaim` — `DLAFacintoTreaty_Act` lgk **13.3.13** |
| `IS_FAC_RETRO` | CHAR(1) | tidak | item bersifat retro fakultatif | `IsFacretro` — `CheckLimitSpreadingTreaty_Act` lgk **17.4** |
| `IS_KOMITE` | CHAR(1) | tidak | item masuk jalur komite | `IsKomite` — `CreateKMTNo_Act` lgk **12** |

⚠️ `IndexObject` · `IndexAjustment` · `IndexPropertyItem` **tidak menjadi kolom** — ⭐ ketiganya
**indeks larik Pega**, digantikan kunci utama.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⚠️ `IndexObject` · `IndexAjustment` ·
> `IndexPropertyItem` **tidak menjadi kolom**"* → dua di antaranya **menjadi kolom** di migrasi `562`:
> `INDEX_PROPERTY_ITEM` (`IndexPropertyItem`) dan `INDEX_ADJUSTMENT` (`IndexAjustment`), bersama `INDEX_COVERAGE` /
> `INDEX_ANEKA` / `INDEX_OCCUPATION` (penunjuk ke baris polis yang dibaca ulang). Item juga
> menyimpan `CLAIM_ID` (NOT NULL, CASCADE) di samping `OBJECT_ID`. Kolom yang berlaku = lampiran pengikat.

### T_CLAIM_ESTIMATION — estimasi nilai kerugian

⭐ **Induk lini FAC:** `T_CLAIM_OBJECT_ITEM` lewat `OBJECT_ITEM_ID`.
⭐ **Induk lini PROP:** `T_GENERAL_CLAIM` lewat `CLAIM_ID`. ⛔ **Kedua kolom nullable.**

| Kolom | Jenis | Wajib | Arti | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bilangan | ya | kunci utama | ⭐ baru |
| ⭐ `OBJECT_ITEM_ID` | bilangan | ⛔ **nullable** | induk lini **FAC** | ⭐ baru |
| ⭐ `CLAIM_ID` | bilangan | ⛔ **nullable** | induk lini **PROP** | ⭐ baru |
| `ESTIMATION_VALUE` | bilangan | ya | nilai estimasi | `EstimationValue` — ⛔ penulis tak ketemu |
| `CURRENCY` | teks | ya | mata uang | `Currency` — ⛔ penulis tak ketemu |
| `CURRENCY_ID` | teks | ya | pengenal mata uang | `CurrencyID` — ⛔ penulis tak ketemu |
| `CONVERT_VALUE` | bilangan | tidak | nilai hasil konversi | `ConvertValue` — ⛔ penulis tak ketemu |
| `DEDUCTIBLE` | bilangan | tidak | risiko sendiri | `Deductible` — ⛔ penulis tak ketemu |
| `GROSS_ESTIMATION_PCT` | bilangan | tidak | persentase estimasi bruto | `GrossEstimationPct` — ⛔ penulis tak ketemu |
| `GROSS_ESTIMATION_PCT_MBU` | bilangan | tidak | persentase estimasi bruto lini MBU | `GrossEstimationPctMBU` — ⛔ penulis tak ketemu |
| `PERSEN_RNM` | bilangan | tidak | bagian Nusantara Re | `PersenRNM` — ⛔ penulis tak ketemu |
| `ESTIMASI_MORE_THAN_TSI` | CHAR(1) | tidak | estimasi melampaui TSI | `EstimasiMoreThanTSI` — ⛔ penulis tak ketemu |
| `ESTIMATION_DATE` | DATE | ya | tanggal estimasi | `EstimationDate` — `ValidateInputEstimate_act` lgk **1** |

⛔ **`CHECK ( (OBJECT_ITEM_ID IS NULL) <> (CLAIM_ID IS NULL) )`**

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"`CHECK ( (OBJECT_ITEM_ID IS NULL) <>
> (CLAIM_ID IS NULL) )`"* → **tidak dibangun** (OQ-CFI-01): `CLAIM_ID` tetap NOT NULL, baris FAC mengisi `CLAIM_ID`
> **dan** `OBJECT_ITEM_ID` (`ADD` nullable, FK CASCADE ke item), nol `CHECK` — lihat RALAT 10-10-2026 di §2.

⚠️⚠️ **Sepuluh dari tiga belas kolom bertanda penulis tak ketemu** — ⭐ angka tertinggi di seluruh
berkas ini. ⛔ **Itu mungkin bukan kebetulan:** estimasi diduga diisi lewat **layar**, bukan lewat
aktivitas. ⛔ **Belum terbukti.** ⭐ **Butir `[terbuka]` BARU — butir 12.**

### T_CLAIM_SPREADING — pembagian klaim per treaty

⭐ **Induk FAC:** `OBJECT_ITEM_ID` · **Induk PROP:** `CLAIM_ID`. ⛔ Keduanya nullable.

| Kolom | Jenis | Wajib | Arti | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bilangan | ya | kunci utama | ⭐ baru |
| ⭐ `OBJECT_ITEM_ID` | bilangan | ⛔ nullable | induk lini FAC | ⭐ baru |
| ⭐ `CLAIM_ID` | bilangan | ⛔ nullable | induk lini PROP | ⭐ baru |
| `TREATY_NAME` | teks | ya | nama treaty | `TreatyName` — `DeleteValueEstimation` lgk **9.1.1** |
| ⚠️ `TREATY_TYPE` | teks | ya | jenis treaty | `TreatyType` — `CopySpreading_Act` lgk **7** |
| ⭐ `SHARE_PERCENTAGE` | bilangan | ya | ⭐ persentase bagian — **lihat §4** | `SharePercentage` — `CopySpreading_Act` lgk **7** |
| `TSI_SPREADED` | bilangan | tidak | TSI yang dibagi | `TSISpreaded` — `GetCoverageAneka_Act` lgk **7** |
| `CLAIM_SPREADED` | bilangan | tidak | klaim yang dibagi | `ClaimSpreaded` — `DeleteValueEstimation` lgk **9.1.1** |
| `CURRENCY` | teks | ya | mata uang | `Currency` — `DeleteValueEstimation` lgk **9.1.1** |
| `CURRENCY_ID` | teks | ya | pengenal mata uang | `CurrencyID` — `DeleteValueEstimation` lgk **9.1.1** |

⛔ **`CHECK ( (OBJECT_ITEM_ID IS NULL) <> (CLAIM_ID IS NULL) )`**

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"`CHECK ( (OBJECT_ITEM_ID IS NULL) <>
> (CLAIM_ID IS NULL) )`"* → **tidak dibangun** (OQ-CFI-01): `CLAIM_ID` tetap NOT NULL, baris FAC mengisi `CLAIM_ID`
> **dan** `OBJECT_ITEM_ID` (`ADD` nullable, FK CASCADE ke item), nol `CHECK` — lihat RALAT 10-10-2026 di §2.

⭐ **Tujuh medan Pega, tujuh penulis ketemu** — ⭐ salah satu dari empat tabel dengan **nol** kolom
tanpa penulis. ⚠️ **`TREATY_TYPE` mencampur kode master dan teks harfiah** — sudah tercatat
`[terbuka]` butir 6.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⚠️ **`TREATY_TYPE` mencampur kode master dan
> teks harfiah** — sudah tercatat `[terbuka]` butir 6."* → yang dibangun **tidak punya kolom `TREATY_TYPE`**: properti
> Pega `TreatyType` disimpan di **`TREATY_ID`** (kolom Claim Prop `525`), apa adanya — kode master maupun teks. Kolom
> **`TSI_SPREADED`** (dan `PREMIUM_SPREADED`) **ditambahkan** migrasi `564`, bersama **`JENIS`**: `'POLIS'` =
> `SpreadingList` polis item, `'KLAIM'` = `SpreadingClaim` item, kosong = baris Prop / Non Prop. Kedua halaman Pega
> **disimpan** sebagai baris tabel ini (RALAT §3 butir 7). Butir 6 register §6 terjawab oleh bentuk itu.

### T_CLAIM_BREAK_QS — pecahan quota share

⭐ **Induk FAC:** `OBJECT_ITEM_ID` · **Induk PROP:** `CLAIM_ID`. ⛔ Keduanya nullable.

| Kolom | Jenis | Wajib | Arti | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bilangan | ya | kunci utama | ⭐ baru |
| ⭐ `OBJECT_ITEM_ID` | bilangan | ⛔ nullable | induk lini FAC | ⭐ baru |
| ⭐ `CLAIM_ID` | bilangan | ⛔ nullable | induk lini PROP | ⭐ baru |
| `TREATY_NAME` | teks | ya | nama treaty | `TreatyName` — `CekExGratia` lgk **6.1** |
| `TREATY_TYPE` | teks | ya | jenis treaty | `TreatyType` — `CekExGratia` lgk **6.1** |
| `SHARE_PERCENTAGE` | bilangan | ya | persentase bagian | `SharePercentage` — `CekExGratia` lgk **6.1** |
| `CURRENCY` | teks | ya | mata uang | `Currency` — `CekExGratia` lgk **6.1** |
| `CURRENCY_ID` | teks | ya | pengenal mata uang | `CurrencyID` — `CekExGratia` lgk **6.1** |

⛔ **`CHECK ( (OBJECT_ITEM_ID IS NULL) <> (CLAIM_ID IS NULL) )`**

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"`CHECK ( (OBJECT_ITEM_ID IS NULL) <>
> (CLAIM_ID IS NULL) )`"* → **tidak dibangun** (OQ-CFI-01): `CLAIM_ID` tetap NOT NULL, baris FAC mengisi `CLAIM_ID`
> **dan** `OBJECT_ITEM_ID` (`ADD` nullable, FK CASCADE ke item), nol `CHECK` — lihat RALAT 10-10-2026 di §2.

⚠️⚠️ **Kelima medannya IDENTIK nama dengan `T_CLAIM_ADJ_SPREADING`**, dan keduanya bersumber halaman
Pega **bernama sama**, hanya **beda tingkat sarang**. ⭐ Itulah sebabnya keduanya **tabel terpisah**.
⭐ **Bukti pendukung:** keduanya ditulis **berkas yang sama** pada langkah **bertetangga** —
`CekExGratia` lgk **6.1** lawan lgk **8.1** — ⭐ tanda kuat keduanya memang **data berbeda**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Bukti pendukung:** keduanya ditulis
> **berkas yang sama** pada langkah **bertetangga** — `CekExGratia` lgk **6.1** lawan lgk **8.1**"* → berkas ini
> **bertentangan dengan dirinya sendiri**: `CekExGratia` 6.1 juga dicatat sebagai sumber `T_CLAIM_BREAK_QS` *dan*
> `T_CLAIM_ADJ_SPREADING`, sementara register §6 butir 1 menyatakan penulis `ObjectItemList[].SpreadingAdjustment`
> **belum ketemu**. Yang dipakai kode yang dibangun (`backend/models/katalog_tabel.go`, `cfs.go`, `adjustment.go`):
>
> | Halaman Pega | Tabel | Penulis yang ditiru |
> | --- | --- | --- |
> | `ObjectItemList(i).SpreadingAdjustment` (Break QS **item**) | `T_CLAIM_BREAK_QS` — `CLAIM_ID` + `OBJECT_ITEM_ID`, tanpa `ADJUSTMENT_ID` | `CheckLimit_Act1` lgk **8** (`GetQuotaShare`; `SusunCedingRetro`) |
> | `Adjustment(a).SpreadingAdjustment` (spreading **adjustment**) | `T_CLAIM_ADJ_SPREADING` — `ADJUSTMENT_ID` | `CountTotalEstimasi_Act` 17.2, `CheckCurrency_ACT` 4, `CekExGratia` 6, `SetSpreadingAjsutement_Act` |
> | `Adjustment(a).SpreadingQuotaShare` | `T_CLAIM_ADJ_QUOTA_SHARE` — `ADJUSTMENT_ID` | Break QS item bermata uang sama: `CountTotalEstimasi_Act` 11 / 17.2, `CheckCurrency_ACT` 7, `CekExGratia` 8; `ExGratia = 1` mengosongkannya |
>
> Jadi `CekExGratia` 6.1 / 8.1 menulis **tingkat adjustment**, bukan `T_CLAIM_BREAK_QS`.

### T_CLAIM_ADJUSTMENT — penyesuaian nilai klaim

⭐ **Induk FAC:** `OBJECT_ITEM_ID` · **Induk PROP:** `CLAIM_ID`. ⛔ Keduanya nullable.
⭐ **Dua sumber Pega:** halaman penyesuaian **dan** halaman bahan pertimbangan komite.

| Kolom | Jenis | Wajib | Arti | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bilangan | ya | kunci utama | ⭐ baru |
| ⭐ `OBJECT_ITEM_ID` | bilangan | ⛔ nullable | induk lini FAC | ⭐ baru |
| ⭐ `CLAIM_ID` | bilangan | ⛔ nullable | induk lini PROP | ⭐ baru |
| `ADJUSTMENT_VALUE` | bilangan | ya | nilai penyesuaian | `AdjustmentValue` — `DLAFacintoTreaty_Act` lgk **13.3.13** |
| `ADJUSTMENT_GROSS` | bilangan | tidak | penyesuaian bruto | `AdjustmentGross` — `DLAFacintoTreaty_Act` lgk **13.3.13** |
| `GROSS_ADJUSTMENT` | bilangan | tidak | bruto penyesuaian | `GrossAdjustment` — `DLAFacintoTreaty_Act` lgk **13.3.13** |
| `VALUE_ADJUSTMENT` | bilangan | tidak | nilai penyesuaian alternatif | `ValueAdjustment` — ⛔ penulis tak ketemu |
| `CURRENCY_ID` | teks | ya | pengenal mata uang | `CurrencyID` — ⛔ penulis tak ketemu |
| `PAYMENT_TYPE` | bilangan | ya | jenis pembayaran | `PaymentType` — ⛔ penulis tak ketemu |
| `ADJUSTER_FEE_VALUE` | bilangan | tidak | biaya adjuster | `AdjusterFeeValue` — ⛔ penulis tak ketemu |
| `SALVAGE_VALUE` | bilangan | tidak | nilai sisa | `SalvageValue` — ⛔ penulis tak ketemu |
| `SHARE_PERSEN_NUSARE` | bilangan | tidak | bagian Nusantara Re | `SharePersenNusare` — `SetTreatyname` lgk **6.5** |
| `TOTAL_SHARE_NUSARE` | bilangan | tidak | total bagian Nusantara Re | `TotalShareNusare` — `SetTreatyname` lgk **6.7** |
| `TOTAL_SPREADING` | bilangan | tidak | total pembagian | `TotalSpreading` — `SetTreatyname` lgk **6.2** |
| `TOTAL_SPREADING_CLAIM` | bilangan | tidak | total pembagian klaim | `TotalSpreadingClaim` — `SetTreatyname` lgk **7.2** |
| `ACCEPTANCE_STATUS` | teks | tidak | status akseptasi | `AcceptanceStatus` — ⛔ penulis tak ketemu |
| `ACCEPTED_NO` | teks | tidak | nomor akseptasi | `AcceptedNo` — ⛔ penulis tak ketemu |
| `IS_PRINT_ACCEPT` | CHAR(1) | tidak | dokumen akseptasi dicetak | `IsPrintAccept` — ⛔ penulis tak ketemu |
| `IS_KOMITE` | CHAR(1) | tidak | penyesuaian masuk jalur komite | `IsKomite` — ⛔ penulis tak ketemu |
| ⭐ `KOMITE_ID` | teks | ⛔ nullable | ⭐ **pintasan tampilan** ke kasus komite | `KomiteNo` — `ViewKomite_act` lgk **3** |
| ⚠️ `IS_FAC_RETRO` | CHAR(1) | tidak | ⚠️ **menandai jalur yang melompati wewenang** | `IsFacRetro` — `CountTotalEstimasi_Act` lgk **19** |
| `DIRECT_TO_KASIR` | CHAR(1) | tidak | dibayar langsung lewat kasir | `DirectToKasir` — `CountTotalEstimasi_Act` lgk **19** |

**Enam kolom dari halaman bahan pertimbangan komite — ditulis SATU langkah yang sama:**

| Kolom | Jenis | Wajib | Arti | Sumber |
| --- | --- | --- | --- | --- |
| `KOMITE_CIRCUM_CAUSE_OF_LOSS` | teks | tidak | duduk perkara penyebab kerugian | `CircumCauseOfLoss` — `CreateKMTNo_Act` lgk **6.8.1.2** |
| `KOMITE_EXTENT_OF_LOSS` | teks | tidak | luas kerugian | `ExtentOfLoss` — idem |
| `KOMITE_LEGAL_LIABILITY` | teks | tidak | tanggung jawab hukum | `LegalLiability` — idem |
| `KOMITE_SALVAGE` | teks | tidak | keterangan nilai sisa | `Salvage` — idem |
| `KOMITE_ADJUSTER_FEE` | teks | tidak | keterangan biaya adjuster | `AdjusterFee` — idem |
| `KOMITE_REMARKS` | teks | tidak | catatan komite | `Remarks` — idem |

⛔ **`CHECK ( (OBJECT_ITEM_ID IS NULL) <> (CLAIM_ID IS NULL) )`**

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"`CHECK ( (OBJECT_ITEM_ID IS NULL) <>
> (CLAIM_ID IS NULL) )`"* → **tidak dibangun** (OQ-CFI-01): `CLAIM_ID` tetap NOT NULL, baris FAC mengisi `CLAIM_ID`
> **dan** `OBJECT_ITEM_ID` (`ADD` nullable, FK CASCADE ke item), nol `CHECK` — lihat RALAT 10-10-2026 di §2.

⭐ **Keenam medan bahan pertimbangan ditulis SATU langkah yang sama** — ⭐ bukti kuat mereka **satu
kelompok**, dan **tidak perlu tabel sendiri**.

⚠️ **Jejak komite tingkat penyesuaian TIDAK menjadi kolom maupun tabel** — ⭐ `[keputusan work
owner]` **K7**; lihat §5c.

### T_CLAIM_ADJ_SPREADING — pembagian atas penyesuaian

⭐ **Induk:** `T_CLAIM_ADJUSTMENT` lewat `ADJUSTMENT_ID` — ⭐ **kunci SAMA di kedua lini**, ⛔ tanpa
`CHECK`.

| Kolom | Jenis | Wajib | Arti | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bilangan | ya | kunci utama | ⭐ baru |
| `ADJUSTMENT_ID` | bilangan | ya | induk penyesuaian | ⭐ baru |
| `TREATY_NAME` | teks | ya | nama treaty | `TreatyName` — `CekExGratia` lgk **6.1** |
| `TREATY_TYPE` | teks | ya | jenis treaty | `TreatyType` — `CekExGratia` lgk **6.1** |
| `SHARE_PERCENTAGE` | bilangan | ya | persentase bagian | `SharePercentage` — `CekExGratia` lgk **6.1** |
| `CURRENCY` | teks | ya | mata uang | `Currency` — `CekExGratia` lgk **6.1** |
| `CURRENCY_ID` | teks | ya | pengenal mata uang | `CurrencyID` — `CekExGratia` lgk **6.1** |

### T_CLAIM_ADJ_QUOTA_SHARE — quota share atas penyesuaian

⭐ **Induk:** `T_CLAIM_ADJUSTMENT` lewat `ADJUSTMENT_ID` — ⭐ **kunci SAMA di kedua lini**.

| Kolom | Jenis | Wajib | Arti | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bilangan | ya | kunci utama | ⭐ baru |
| `ADJUSTMENT_ID` | bilangan | ya | induk penyesuaian | ⭐ baru |
| `TREATY_NAME` | teks | ya | nama treaty | `TreatyName` — `CekExGratia` lgk **8.1** |
| `TREATY_TYPE` | teks | ya | jenis treaty | `TreatyType` — `CekExGratia` lgk **8.1** |
| `SHARE_PERCENTAGE` | bilangan | ya | persentase bagian | `SharePercentage` — `CekExGratia` lgk **8.1** |
| `CURRENCY` | teks | ya | mata uang | `Currency` — `CekExGratia` lgk **8.1** |
| `CURRENCY_ID` | teks | ya | pengenal mata uang | `CurrencyID` — `CekExGratia` lgk **8.1** |

⚠️ **Kelima medannya identik nama dengan `T_CLAIM_ADJ_SPREADING`**, ⭐ **tetapi ditulis langkah
bertetangga di berkas yang sama** — lgk **6.1** lawan lgk **8.1**. ⛔ **Itu bukti kuat keduanya
memang data yang berbeda**, bukan salinan.

### T_CLAIM_FAC_RETRO — retro fakultatif

⭐ **Induk lini FAC:** `T_CLAIM_ADJUSTMENT` lewat `ADJUSTMENT_ID`.
⭐ **Induk lini PROP:** `T_GENERAL_CLAIM` lewat `CLAIM_ID`. ⚠️ **Beda TINGKAT** — `[keputusan work
owner]` **K5**.

| Kolom | Jenis | Wajib | Arti | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bilangan | ya | kunci utama | ⭐ baru |
| ⭐ `ADJUSTMENT_ID` | bilangan | ⛔ nullable | induk lini **FAC** | ⭐ baru |
| ⭐ `CLAIM_ID` | bilangan | ⛔ nullable | induk lini **PROP** | ⭐ baru |
| `REINSURER_ID` | teks | ya | pengenal reasuradur | `ReinsurerID` — `CheckLimitSpreadingTreaty_Act` lgk **21.2** |
| `REINSURER_NAME` | teks | ya | nama reasuradur | `ReinsurerName` — idem |
| `PCT_SHARE_ALL_OBJ` | bilangan | ya | persentase bagian seluruh objek | `PctShareAllObj` — idem |
| `RI_COMM_ALL_OBJ` | bilangan | tidak | komisi reasuransi seluruh objek | `RiCommAllObj` — idem |
| `OFFICER` | teks | tidak | petugas | `Officer` — idem |
| `ATTENTION` | teks | tidak | ditujukan kepada | `Attention` — ⛔ penulis tak ketemu |
| `OUR_REF` | teks | tidak | nomor rujukan kami | `OurRef` — `DraftGenerateDLAFacin_Act` lgk **12.2.1.2.1** |
| `START_PERIOD` | DATE | tidak | awal periode | `StartPeriod` — idem |
| `END_PERIOD` | DATE | tidak | akhir periode | `EndPeriod` — idem |
| `ADDITIONAL_INFO` | teks | tidak | keterangan tambahan | `AdditionalInfo` — idem |
| `TOTAL_ESTIMASI_REAS` | bilangan | tidak | total estimasi bagian reasuransi | `TotalEstimasiReas` — `DLAFacintoTreaty_Act` lgk **13.3.13** |

⛔ **`CHECK ( (ADJUSTMENT_ID IS NULL) <> (CLAIM_ID IS NULL) )`**

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"`CHECK ( (ADJUSTMENT_ID IS NULL) <>
> (CLAIM_ID IS NULL) )`"* → **tidak dibangun** (OQ-CFI-01): `CLAIM_ID` tetap NOT NULL dan selalu diisi;
> `ADJUSTMENT_ID` (`567`, `ADD` nullable, FK CASCADE ke adjustment) terisi pada baris retro FAC tingkat adjustment dan
> kosong pada baris retro tingkat klaim (`ClaimData.FacRetroTreaty`) — lihat RALAT 10-10-2026 di §2.

⚠️ Daftar lokasi **tidak menjadi kolom** — ⭐ ia **larik bersarang**, dan ⛔ `[terbuka]` apakah ia
menjadi tabel tersendiri. ⭐ **Butir `[terbuka]` BARU — butir 13.**

### ⭐ Rekapitulasi kolom

| Tabel | Kolom | Tanpa penulis |
| --- | ---: | ---: |
| `T_CLAIM_OBJECT` | 27 | 8 |
| `T_CLAIM_OBJECT_ITEM` | 23 | 8 |
| `T_CLAIM_ESTIMATION` | 13 | 10 |
| `T_CLAIM_SPREADING` | 10 | 0 |
| `T_CLAIM_BREAK_QS` | 8 | 0 |
| `T_CLAIM_ADJUSTMENT` | 28 | 8 |
| `T_CLAIM_ADJ_SPREADING` | 7 | 0 |
| `T_CLAIM_ADJ_QUOTA_SHARE` | 7 | 0 |
| `T_CLAIM_FAC_RETRO` | 15 | 1 |
| ⭐ **TOTAL** | ⭐ **138** | ⚠️ **35** |

⚠️⚠️ **Tiga puluh lima kolom belum diketahui penulisnya.** ⛔ **Itu tidak menahan migrasi BACA** —
medannya terbaca dipakai. ⭐ Yang tertahan adalah **jalur TULIS di Go**, sejalan butir 1.

⚠️ **Dan dua puluh sembilan medan bergantung jenis objek belum dihitung di sini sama sekali** —
⭐ butir 11.
---

## §3 — Sebelas halaman Pega Fac In yang TIDAK dibawa

| # | Halaman Pega | Sebab tidak dibawa | Penggantinya di Go |
| --- | --- | --- | --- |
| **1** | `ObjectList[].CoverageList` | ⭐ salinan polis; klaim **hanya menulis hasil konversi kurs** | ⭐ kolom `KURS_OBJECT_ITEM` dan `VALUE_TSI_NUSARE_IDR` di `T_CLAIM_OBJECT_ITEM` |
| **2** | `ObjectList[].PropertyItemList` | salinan polis | dibaca dari **modul polis** |
| **3** | `ObjectList[].ObjectCoverageList` | salinan polis | dibaca dari **modul polis** |
| **4** | `ObjectList[].AnekaList` | salinan polis | dibaca dari **modul polis** |
| **5** | `ObjectList[].OccupationList` | salinan polis | dibaca dari **modul polis** |
| **6** | `ObjectList[].LocationList` | salinan polis | dibaca dari **modul polis** |
| ⭐ **7** | `ObjectItemList[].SpreadingList` | ⭐ **himpunan bagian** dari `SpreadingClaim` — **nol medan hilang** | ⭐ `SELECT DISTINCT` dari `T_CLAIM_SPREADING` |
| **8** | `ObjectList[].LimitBusinessList` | ⛔ **nol medan dipakai** di seluruh korpus | ⛔ tidak ada |
| **9** | `ObjectList[].TreatyLimits` | ⛔ **nol medan dipakai** di seluruh korpus | ⛔ tidak ada |
| **10** | medan `CARI1`–`CARI50` · `Test` · `result` | ⭐ **sampah layar** | ⛔ tidak ada |
| ⚠️ **11** | `PrintFaceClaiml` | ⚠️ **salah ketik korpus** — di Pega, nama yang salah ketik **membuat halaman baru diam-diam** | ⛔ tidak ada |

### ⭐ Bukti butir 7 — selisih NIHIL

| Halaman | Medannya |
| --- | --- |
| `SpreadingList` | `SharePercentage` · `TSISpreaded` · `TreatyName` · `TreatyType` |
| `SpreadingClaim` | **keempatnya** **+** `ClaimSpreaded` · `Currency` · `CurrencyID` |

⭐ **`SpreadingList` ⊆ `SpreadingClaim`, dan selisihnya NIHIL.** ⛔ Karena itu ia **tidak menjadi
tabel** — ia **pandangan** atas tabel yang sudah ada.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⛔ Karena itu ia **tidak menjadi tabel** — ia
> **pandangan** atas tabel yang sudah ada."* → tidak menjadi tabel **sendiri** tetap benar, tetapi ia **bukan
> pandangan `SELECT DISTINCT`**: `SpreadingList` polis item disimpan sebagai baris **`T_CLAIM_SPREADING` ber-`JENIS =
> 'POLIS'`**, `SpreadingClaim` sebagai baris ber-`JENIS = 'KLAIM'` (migrasi `564`; katalog `TabelSpreadPolis` /
> `TabelSpreadKlaim`). `SpreadingClaim` mula-mula disalin dari `SpreadingList`, lalu ditulis ulang per mata uang saat
> estimasi berubah (`backend/models/estimasi.go`, `estimasi_baris.go`) — keduanya ditiru sebagai baris tersimpan,
> bukan diturunkan dengan kueri.

---

## §4 — ⚠️ Aturan keseragaman berbagi per treaty

⚠️ **Aturan ini lahir dari §3 butir 7, dan ia mengikat lapisan layanan — bukan basis data.**

⭐ **Sebabnya:** `SpreadingClaim` berbaris **treaty × mata uang**, sedangkan `SpreadingList` berbaris
**treaty saja**. ⛔ Agar keduanya tetap sepadan sesudah `SpreadingList` dibuang, satu hal wajib
dijamin:

> ⭐⭐ **Di dalam satu item objek, semua baris ber-`TREATY_ID` sama WAJIB punya `SHARE_PERCENTAGE`
> yang sama.**

⭐ **Akibatnya pada perilaku:** ⛔ **menyunting persentase berbagi mengenai SELURUH baris treaty itu
sekaligus** — bukan satu baris. ⚠️ Antarmuka yang membiarkan pengguna menyunting satu baris saja
**akan melanggar aturan ini tanpa galat**.

⛔ **Penegakannya di lapisan layanan**, sebab basis data tidak dapat menyatakannya dengan satu
`CHECK` baris tunggal.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Akibatnya pada perilaku:** ⛔ **menyunting
> persentase berbagi mengenai SELURUH baris treaty itu sekaligus** — bukan satu baris."* → **`[penyimpangan sadar]`
> DIBUANG — tidak dibangun.** Aturan dan akibatnya adalah **rancangan asisten**, bukan perilaku XML: di XML medan
> **Share % spreading estimasi hanya-baca (RO=ALWAYS)** di `Section/Estimasi`, `EstimasiPA`, `EstimasiMarine`, sehingga
> `CheckTotalSpreadingPct_Act` tak terjangkau; grid spreading yang dapat disunting dan tombol **Save Spreading**
> bergerbang `ClaimData.ExGratia = 1`, yang **tidak pernah benar** — penulis satu-satunya `InsertObjects_dt` langkah 11
> menulis `0` (`PARITAS.md` §3 dan §7 butir 10, OQ-CFI-17). Tidak ada jalur sunting persentase, jadi tidak ada penegak;
> register §6 butir 9 gugur. Spreading **adjustment** Ex Gratia (`Adjustment.ExGratia = 1`) memang dapat disunting per
> baris di XML dan dibangun begitu (`CountSpreadingAdjustment`).

---

## §5 — Sisi Komite — ⭐ nol tabel baru

`[keputusan work owner]` **K4** — ⭐ **Lini FAC memakai `T_GENERAL_KOMITE` dan
`T_KOMITE_KOMITELIST` apa adanya.** ⛔ **Tidak ada tabel baru untuk modul komite.**

### 5a — Enam medan Pega cocok satu per satu

`[terverifikasi]` Medan `KomiteList` di lini FAC **cocok satu per satu** dengan kolom yang sudah
terkunci:

| Medan Pega — FAC | Kolom yang sudah terkunci |
| --- | --- |
| `KomiteID` | `KOMITE_OPERATORID` |
| `IDKomite` | `KOMITE_JABATAN` |
| `KomiteEmail` | `KOMITE_EMAIL` |
| `KomiteAproval` | `KOMITE_APPROVAL` |
| `KomiteComment` | `KOMITE_COMMENT` |
| `DateApprove` | `DATE_APPROVE` |

⭐ **Tautan induknya:** `T_GENERAL_KOMITE.ADJUSTMENT_ID` untuk lini FAC menunjuk
`T_CLAIM_ADJUSTMENT.ID`.

### 5b — ⭐ Yang berbeda dari Komite Claim Prop: **perilakunya, bukan tabelnya**

`[terverifikasi]` — enam perbedaan, semuanya **perilaku**:

| # | Hal | Bukti |
| --- | --- | --- |
| **1** | `ApprovalKomite_Act` **ada di FAC, tidak ada di PROP** | padanan PROP `SetKomiteList_Act` **12 langkah**, isinya **hanya menjumlah total penyesuaian** — ⛔ **nol `Obj-Browse` anggota**, ⛔ **nol pita wewenang** |
| ⚠️ **2** | **Larangan menyetujui klaim sendiri** | `ApprovalKomite_Act` lgk **12** · `KomiteList(1).KomiteID == .OPERATOR_ID` — ⛔ **hanya anggota PERTAMA**; ⛔ **di PROP nol pemeriksaan** |
| **3** | **Gerbang jabatan** | `ApprovalKomite_Act` lgk **9** · `OperatorID.pyPosition == "SPV B"` — **hidup di FAC**, ⛔ **nihil di PROP** |
| ⚠️ **4** | **Fac Retro melompati wewenang** | `ApprovalKomite_Act` lgk **4** · `Exit-Activity` when `.IsFacRetro==1` |
| ⭐ **5** | `KomiteRouter` | ⭐ **identik langkah demi langkah** dengan milik PROP — **8 langkah, gerbang sama** |
| **6** | `ProteksiKomite.CARI1` | ada di `ShowTransfer` **FAC** · ⛔ **nol** di `ShowTransfer` **PROP** |

⭐⭐ **Kesimpulan bentuk:** ⛔ keenam perbedaan itu **tidak menuntut satu kolom pun berbeda**. ⭐ Ia
**perbedaan aturan**, dan tempatnya di **spec modul komite**, bukan di struktur tabel.

⛔ **Berkas struktur terpisah untuk modul komite TIDAK dibuat di sini** — ia ditulis **sesudah**
spec-nya, mengikuti urutan modul lain.

---

### 5c — ⭐⭐ Jejak komite: SATU tabel — dan pembuktiannya

`[keputusan work owner]` **K7 · 2026-09-20** — ⭐ **Hanya `T_KOMITE_KOMITELIST`.** Jejak komite di
tingkat **objek** dan di tingkat **penyesuaian** menjadi **kueri**, ⛔ bukan tabel.

⚠️ **Rekomendasi asisten berpagar:** *"perlu dibuktikan dulu — kalau ternyata isinya berbeda, kita
kehilangan data."* ⭐ **Pagar itu DILEPAS di sini, dan hasilnya LULUS.**

`[terverifikasi]` **Jendela: 596 berkas** — `Claim Fac In` (482) dan `Komite Claim FacIn` (114),
mencari setiap rujukan bermedan pada kedua halaman jejak komite.

| | Medan unik | Berkas |
| --- | ---: | ---: |
| jejak tingkat **penyesuaian** | **6** | 4 |
| jejak tingkat **objek** | **6** | 6 |
| ⭐ **irisan** | ⭐ **6** | |
| ⛔ hanya di tingkat penyesuaian | ⛔ **0** | |
| ⛔ hanya di tingkat objek | ⛔ **0** | |

⭐ **Keenam medannya sama persis:** pengenal anggota · jabatan · surel · persetujuan · komentar ·
tanggal persetujuan.

⛔⛔ **Dan buktinya telak:** `[terverifikasi]` `KomitePost_Adjustment` menulis jejak tingkat
penyesuaian sebagai **SALINAN SATU HALAMAN UTUH** dari jejak tingkat objek — ⭐ satu penugasan,
satu langkah, tanpa memilih medan.

✅ **K7 BERDIRI.** ⭐ Butir 3 dan butir 4 **ditutup dengan bukti**, ⛔ bukan hanya dengan keputusan.

⚠️ **Satu nuansa yang wajib ikut tertulis:** salinan itu **kemudian diisi nilai keputusan** —
persetujuan, komentar, dan tanggal ditulis ke dalamnya saat komite memutuskan. ⭐ Artinya tabel
tunggal itu **wajib memuat ketiga medan keputusan**, dan ia **sudah memuatnya**.

---

### 5d — Empat keputusan work owner yang mengikat berkas ini

`[keputusan work owner]` **2026-09-20** — ⭐ lanjutan **K1–K4** di §0.

| # | Keputusan | Catatan wajib |
| --- | --- | --- |
| ⚠️ **K5** | **`T_CLAIM_FAC_RETRO` tetap berinduk BEDA TINGKAT** — lini **PROP** bersumber daftar retro tingkat **klaim**; lini **FAC** bersumber daftar retro tingkat **penyesuaian** | ⛔ **BERBEDA dari rekomendasi asisten**, yang mengusulkan induk penyesuaian di kedua lini. ⛔ **Disengaja** — ⭐ ia **meniru Pega apa adanya**. ⚠️ **Akibat:** di lini PROP, satu klaim dengan **dua penyesuaian** membuat baris retro **tidak dapat dibedakan milik penyesuaian mana**. ⭐ **Keterbatasan itu SUDAH ADA di sistem lama** — ⛔ bukan yang alih ini tambahkan |
| ⭐ **K6** | **Komite lini PROP dan lini FAC disimpan di TABEL YANG SAMA** | ✅ Menguatkan **K4**. ⭐ **Akibat pada berkas:** ⛔ berkas struktur dan relasi terpisah untuk modul komite **TIDAK dibuat** — tabelnya sama, jadi tempat mencatatnya satu. ⚠️ Bila work owner kelak menghendaki simetri berkas, itu **keputusan tersendiri** |
| ⭐ **K7** | **Jejak komite: SATU tabel saja** | ✅ Sejalan rekomendasi asisten, ⭐ **dan pagarnya sudah dilepas** — lihat §5c |
| ⛔ **K8** | **Tidak ada tabel jejak audit bernama `T_CLAIM_AUDIT_TRAIL`** — ⭐ yang ada tabel kronologi lintas-lini, isinya diambil dari kronologi klaim | ⭐ Butir 5 dan butir 8 **ditutup**; relasi **ke-16** ditambahkan di berkas relasi, ⚠️⚠️ bertanda **JANGAN cascade** |

> ⛔ **RALAT 10-10-2026 atas K8.** Sel lamanya dikutip utuh, tidak dihapus: *"relasi **ke-16** ditambahkan di berkas
> relasi, ⚠️⚠️ bertanda **JANGAN cascade**"* → relasi ke-16 **CASCADE** (Claim Prop `532`) — lihat RALAT 10-10-2026 di
> §1 bab Kronologi. K5 tetap: `T_CLAIM_FAC_RETRO` FAC berinduk adjustment (`567`), dengan `CLAIM_ID` tetap diisi.

---

## §6 — Register `[terbuka]`

⭐ **Sepuluh butir.** ⛔ **Nol ditutup oleh berkas ini.**

| # | Butir | Pemilik |
| --- | --- | --- |
| ⚠️ **1** | ⛔ **Penulis baris `ObjectItemList[].SpreadingAdjustment` belum ketemu.** Yang terbaca hanya **pembacanya** — `CekExGratia` · `CheckCurrency_ACT` · `CheckLimit_Act1` — dan pencetak PDF akseptasi `PrintPDFAccep_MultiAksep`. ⭐ **Migrasi tidak terganggu**; yang tertahan **jalur TULIS di Go**. ⚠️ **Pengurai datar tidak bisa menjawabnya** — `CountTotalEstimasi_Act` **bersarang**. Butuh **pengurai yang menghormati sarang**, atau **tangkapan layar Pega Designer** | asisten |
| ✅ **2** | ✅ **DITUTUP 2026-09-20 oleh K5** — induk **tetap beda tingkat**: FAC dari penyesuaian, PROP dari klaim. ⛔ **Berbeda dari rekomendasi asisten**, disengaja. — *(kalimat lama:)* ⛔ **`T_CLAIM_FAC_RETRO` — induknya beda TINGKAT antar lini.** Di FAC induknya **penyesuaian**, di PROP induknya **klaim**. ⚠️ Dan di `STRUKTUR-TABEL-CLAIM-PROP.md`, `AdjustmentList(n).FacRetroList` justru masuk daftar **"bukan tabel"**. ⛔ **Perlu satu keputusan work owner** | `[work owner]` |
| ✅ **3** | ✅ **DITUTUP 2026-09-20 oleh K7** — **satu tabel saja**, `T_KOMITE_KOMITELIST`; jejak tingkat objek menjadi **kueri**. ⚠️ **Pembuktiannya dikerjakan dan LULUS** — lihat §5c. — *(kalimat lama:)* **Jejak komite `ObjectList[].KomiteList[]`** — ★ khas FAC, **dua medan** (`DateApprove`, `KomiteAproval`), **plus** `IsKomiteApprove` dan `TotalKomite` di tingkat objek. ⭐ Muat di `T_KOMITE_KOMITELIST` yang sudah terkunci, **atau tabel sendiri?** | `[work owner]` |
| ✅ **4** | ✅ **DITUTUP 2026-09-20 oleh K7, dan PEMBUKTIANNYA LULUS** — lihat §5c: irisan **6 dari 6**, selisih **NOL**, dan `KomitePost_Adjustment` menulis **salinan satu halaman utuh**. ⭐ **Duplikasi Pega, terbukti** | — |
| ✅ **5** | ~~Apakah lini FAC memakai tabel kronologi lintas-lini — belum ditelusuri~~ — ✅ **DITUTUP 2026-09-20 oleh K8:** lini FAC **memakai** `T_VIEW_SUGGEST`; ia tabel lintas-lini | — |
| **6** | **`TREATY_TYPE` ada di FAC, tidak ada di PROP** · **`TSI_SPREADED` ada di PROP, belum terbukti di FAC** | asisten |
| ⚠️ **7** | ⚠️ **Dua berkas lini Life masih menyebut nama tabel penyesuaian yang lama** di tabel pemetaan **LINI → tabel**. ⛔ Berkas Life **tidak boleh disentuh** `[keputusan work owner]`, jadi ketidakcocokan nama ini **DICATAT, bukan diperbaiki** | `[work owner]` |
| ✅ **8** | ~~Tabel jejak audit belum ditelusuri di lini FAC, namanya belum ditetapkan di lini PROP~~ — ✅ **DITUTUP 2026-09-20 oleh K8:** ⛔ **tabel itu tidak ada**; yang ada `T_VIEW_SUGGEST` dari kronologi. ⭐ Pertanyaannya **gugur** | — |
| ⭐ **9** | ⭐ **BARU — aturan keseragaman berbagi per treaty *(§4)* belum punya penegak yang ditunjuk.** ⛔ Basis data tidak dapat menyatakannya dengan satu `CHECK` baris tunggal, jadi ia jatuh ke lapisan layanan — ⚠️ **tetapi belum diputuskan apakah ia dijaga saat tulis, atau diperiksa berkala** | `[work owner]` |
| ⭐ **10** | ⭐ **BARU — `DataCommitteFacin` ikut menjadi sumber `T_CLAIM_ADJUSTMENT` bersama `Adjustment`.** ⚠️ Dua halaman Pega yang berbeda mengisi **satu** tabel; ⛔ **belum diperiksa apakah medannya bertabrakan** atau saling melengkapi | asisten |

⭐ **Butir 8 · 9 · 10 ditambahkan sendiri oleh berkas ini.** ⛔ Tujuh butir pertama berasal dari
keputusan yang memerintahkannya.

> ⛔ **RALAT 10-10-2026 atas butir 1, 6, dan 9.** Baris tabel tidak dapat memuat blok kutipan, jadi ralatnya di sini;
> kalimat lamanya dikutip utuh, tidak dihapus:
>
> - **1** — *"⛔ **Penulis baris `ObjectItemList[].SpreadingAdjustment` belum ketemu.**"* → yang ditiru kode:
>   `CheckLimit_Act1` langkah 8 (`GetQuotaShare`; `backend/models/cfs.go` `SusunCedingRetro`) → `T_CLAIM_BREAK_QS`.
>   Lihat RALAT di `T_CLAIM_BREAK_QS` §2b.
> - **6** — *"**`TREATY_TYPE` ada di FAC, tidak ada di PROP** · **`TSI_SPREADED` ada di PROP, belum terbukti di FAC**"* →
>   kolomnya `TREATY_ID` (isi `TreatyType`); `TSI_SPREADED` ditambahkan `564` dan ditulis FAC. Lihat RALAT di
>   `T_CLAIM_SPREADING` §2b.
> - **9** — *"⭐ **BARU — aturan keseragaman berbagi per treaty *(§4)* belum punya penegak yang ditunjuk.**"* → **gugur**:
>   aturan §4 dibuang sebagai rancangan di luar XML (RALAT §4).
>
> Butir 11 (29 medan jenis objek) terjawab oleh bentuk yang dibangun: dibaca dari polis lewat `KUNCI_POLIS` (RALAT
> `T_CLAIM_OBJECT` §2b). Aritmetika register di bawah adalah keadaan 20-09-2026.

### ⭐ Butir BARU dari putaran kolom — 2026-09-20

| # | Butir | Pemilik | Memblokir? |
| --- | --- | --- | --- |
| ⚠️ **11** | ⚠️ **Dua puluh sembilan medan `T_CLAIM_OBJECT` bergantung JENIS OBJEK** — orang, kendaraan, bangunan, kapal. ⛔ **Semuanya tanpa penulis di korpus klaim**, artinya **dibaca dari modul polis**. ⭐ **Menjadi kolom, atau dibaca lewat rujukan ke polis?** | `[work owner]` | tidak |
| ⚠️ **12** | ⚠️ **Sepuluh dari tiga belas kolom `T_CLAIM_ESTIMATION` tanpa penulis** — ⭐ angka tertinggi di berkas ini. ⛔ **Diduga** diisi lewat layar, bukan aktivitas; ⛔ **belum terbukti** | asisten | tidak |
| **13** | **Daftar lokasi di dalam retro fakultatif** — larik bersarang; ⭐ menjadi tabel tersendiri atau tidak | `[work owner]` | tidak |

### ⭐ Aritmetika register, dua cara

**CARA-1 — cacah baris:** dibawa `1 · 6 · 7 · 9 · 10` = **5** · baru `11 · 12 · 13` = **3**
⇒ ⭐ **8**

**CARA-2 — delta:** `10 sebelumnya` − `5 ditutup (2·3·4·5·8)` + `3 baru` = ⭐ **8**

✅ **Keduanya sepakat: 8 butir terbuka.**

⛔ **Tak satu pun memblokir** putaran DDL berikutnya — ⭐ kesemuanya menambah kolom atau tabel,
tidak mengubah yang sudah ditulis.

---

## Lampiran — bukti berkas lain tidak disentuh

| Berkas / folder | Keadaan |
| --- | --- |
| ⛔ `claim-life\` · `komite-claim-life\` · `premiumlist-life\` · `endorsement-life\` | ✅ **NOL disentuh, NOL dibuka untuk disunting** — dibuktikan dengan md5 agregat |
| `claim-facin\spec.md` | ✅ utuh — hanya dibaca |
| `komite-claim-facin\` seluruhnya | ✅ **NOL berkas ditambah, NOL disunting** |
| `komite-claim-prop\spec.md` · `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md` | ✅ utuh — hanya dibaca |
| korpus `Claim Fac In` · `Komite Claim FacIn` · `Claim Prop` · `Komite Claim Prop` | ✅ md5 tidak berubah |

⛔ Kode **NOL** · DDL **NOL** · `CREATE TABLE` **NOL** · **daftar kolom NOL** · nomor baris XML
**NOL**.

---

## Lampiran pengikat penjaga — kolom DDL tabel milik Claim Fac In (10-10-2026)

Bab di bawah **mengikat**: `inti/backend/penjaga` (`TestKolomDDLCocokDenganStruktur`,
`TestGolonganTipeDDLCocokDenganStruktur`) membandingkan setiap judul `## T_…` dengan DDL migrasi `561`–`562`.
Bab §1–§6 di atas tetap catatan rancangan 20-09-2026 beserta RALAT bertanggalnya; bila keduanya berbeda, bab
di bawah inilah yang berlaku. Kolom FAC di tabel bersama dicatat di dokumen pemiliknya, satu tabel satu
dokumen: `T_GENERAL_CLAIM` (migrasi `560`) di `modul/claimlife/docs/STRUKTUR-TABEL-CLAIM-LIFE.md`;
`T_CLAIM_ESTIMATION` / `T_CLAIM_SPREADING` / `T_CLAIM_BREAK_QS` / `T_CLAIM_ADJUSTMENT` / `T_CLAIM_FAC_RETRO`
(migrasi `563`–`567`) di lampiran pengikat `modul/claimprop/docs/STRUKTUR-TABEL-CLAIM-PROP.md`.

Uang, persen, share, dan kurs `NUMBER(38,10)`. Nama kolom mengikuti penamaan Claim Prop (prompt §6 butir 2);
katalog properti -> kolom = `backend/models/katalog_tabel.go`.

## T_CLAIM_OBJECT

Migrasi `561`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `561` — `SEQ_T_CLAIM`, stabil |
| `CLAIM_ID` | teks | tidak | FK | migrasi `561` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `561` |
| `KUNCI_POLIS` | teks | ya |  | migrasi `561` |
| `OBJECT_REF` | teks | ya |  | migrasi `561` |
| `OBJECT_NAME` | teks | ya |  | migrasi `561` |
| `CFS` | teks | ya |  | migrasi `561` |
| `PRINT_FACE_CLAIM` | teks | ya |  | migrasi `561` |
| `PLA_STATUS` | teks | ya |  | migrasi `561` |
| `IS_FAC_RETRO` | teks | ya |  | migrasi `561` |
| `IS_MORE_THAN_TREATY_LIMIT` | teks | ya |  | migrasi `561` |
| `IS_KOMITE` | teks | ya |  | migrasi `561` |
| `IS_PRINT_ACCEPT` | teks | ya |  | migrasi `561` |
| `REMARKS_PLA` | teks | ya |  | migrasi `561` |
| `REMARKS_DLA` | teks | ya |  | migrasi `561` |
| `SHARE_RETRO` | teks | ya |  | migrasi `561` |
| `NO_DLA` | teks | ya |  | migrasi `561` |
| `DLA_STATUS` | teks | ya |  | migrasi `561` |
| `SAVE_SPREADING` | teks | ya |  | migrasi `561` |

**Unik:** `(CLAIM_ID, NOURUT)`.

## T_CLAIM_OBJECT_ITEM

Migrasi `562`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `562` — `SEQ_T_CLAIM`, stabil |
| `OBJECT_ID` | teks | tidak | FK | migrasi `562` — → `T_CLAIM_OBJECT.ID`, ON DELETE CASCADE |
| `CLAIM_ID` | teks | tidak | FK | migrasi `562` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `562` |
| `OBJECT_ITEM_NAME` | teks | ya |  | migrasi `562` |
| `OBJECT_ITEM_REF` | teks | ya |  | migrasi `562` |
| `INDEX_PROPERTY_ITEM` | teks | ya |  | migrasi `562` |
| `INDEX_COVERAGE` | teks | ya |  | migrasi `562` |
| `INDEX_ANEKA` | teks | ya |  | migrasi `562` |
| `INDEX_OCCUPATION` | teks | ya |  | migrasi `562` |
| `COVERAGE_ID` | teks | ya |  | migrasi `562` |
| `COVERAGE_OLD_ID` | teks | ya |  | migrasi `562` |
| `COVERAGE_NOTE` | teks | ya |  | migrasi `562` |
| `OCCUPATION_ID` | teks | ya |  | migrasi `562` |
| `OCCUPATION_NAME` | teks | ya |  | migrasi `562` |
| `CURRENCY_ID` | teks | ya |  | migrasi `562` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `562` |
| `KURS_OBJECT_ITEM` | angka desimal | ya |  | migrasi `562` |
| `TSI_PER_OBJECT` | angka desimal | ya |  | migrasi `562` |
| `TSI_NUSARE` | angka desimal | ya |  | migrasi `562` |
| `VALUE_TSI_NUSARE_IDR` | angka desimal | ya |  | migrasi `562` |
| `LIMIT_OF_LIABILITY` | angka desimal | ya |  | migrasi `562` |
| `PREMI_NUSARE` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_GROSS_PREMI` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_PREMIUM_NUSARE` | angka desimal | ya |  | migrasi `562` |
| `DEDUCTIBLE_TYPE` | teks | ya |  | migrasi `562` |
| `DEDUCTIBLE_FORM_TYPE` | teks | ya |  | migrasi `562` |
| `DEDUCTIBLE_CURRENCY_ID` | teks | ya |  | migrasi `562` |
| `DEDUCTIBLE_PCT` | angka desimal | ya |  | migrasi `562` |
| `DEDUCTIBLE_BASIS` | teks | ya |  | migrasi `562` |
| `TSI_DEDUCTIBLE` | angka desimal | ya |  | migrasi `562` |
| `CLAIM_DEDUCTIBLE` | angka desimal | ya |  | migrasi `562` |
| `DEDUCTIBLE_VALUE` | angka desimal | ya |  | migrasi `562` |
| `NET_DEDUCTIBLE_VALUE` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_ESTIMASI` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_GROSS_ESTIMASI` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_GROSS_ESTIMASI_IDR` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_CLAIM_SPREADED` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_ESTIMATION_VALUE_IDR` | angka desimal | ya |  | migrasi `562` |
| `ESTIMATION_RESERVE` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_SHARE_PCT` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_TSI_SPREADED` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_PREMIUM_SPREADED` | angka desimal | ya |  | migrasi `562` |
| `PRINT_FACE_CLAIM` | teks | ya |  | migrasi `562` |
| `IS_FAC_RETRO` | teks | ya |  | migrasi `562` |
| `IS_MORE_THAN_TREATY_LIMIT` | teks | ya |  | migrasi `562` |
| `IS_KOMITE` | teks | ya |  | migrasi `562` |
| `PLA_REV_TREATY` | teks | ya |  | migrasi `562` |
| `ADJUSTMENT_VAL` | teks | ya |  | migrasi `562` |
| `IS_ADJ_VAL` | teks | ya |  | migrasi `562` |
| `PAYMENT_TYPE` | teks | ya |  | migrasi `562` |
| `INDEX_ADJUSTMENT` | teks | ya |  | migrasi `562` |
| `TOTAL_HASIL_CLAIM` | angka desimal | ya |  | migrasi `562` |
| `TOTAL_ESTIMASI_REAS` | angka desimal | ya |  | migrasi `562` |
| `PERCENT_REAS` | angka desimal | ya |  | migrasi `562` |
| `PERCENT_HANDLING_FEE` | angka desimal | ya |  | migrasi `562` |
| `PERCENT_FACOUT_CLAIM` | angka desimal | ya |  | migrasi `562` |

**Unik:** `(CLAIM_ID, NOURUT)`. **Index:** `OBJECT_ID`.
