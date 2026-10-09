# Relasi Tabel — Claim Prop

> ⭐ `[keputusan work owner]` **2026-09-20.** **Awalan tabel klaim non-life yang lama — berhuruf `P` sesudah `CLAIM` — diganti menjadi `T_CLAIM_`**, sebab tabelnya kini **dipakai bersama lini FAC dan PROP** sehingga huruf **P** pada awalan menyesatkan. ⭐ **NONPROP nanti ikut tabel yang sama.**
>
> ⛔ **Nama awalan lamanya sengaja TIDAK dikutip harfiah** di berkas mana pun di luar lini Life — ⭐ supaya pencarian atas awalan lama itu **hanya** menemukan berkas yang memang belum diselaraskan, bukan kalimat yang menerangkan penggantiannya. ⭐ **Tujuh tabel dipakai bersama; lima di antaranya berkunci asing GANDA.** ⛔ **Tabel lini Life tidak disentuh** — awalannya berbeda dan tidak ikut terganti.
>
> ⚠️ Akibatnya **dua berkas lini Life masih menyebut nama tabel penyesuaian yang lama**. ⭐ Itu **dicatat sebagai `[terbuka]`**, ⛔ **bukan diperbaiki**.

> ⚠️ **RALAT 07-10-2026** — kalimat lama: *"Tujuh tabel dipakai bersama; lima di antaranya berkunci asing GANDA."* Tabel bersama berjumlah **lima** (`T_WORK_CLAIM`, `T_GENERAL_CLAIM`, `T_GENERAL_KOMITE`,
> `T_KOMITE_KOMITELIST`, `T_VIEW_SUGGEST`); "tujuh" basi.


**Tanggal:** 2026-09-19 · **Modul:** `Claim Prop` · Ujung komite dilacak di `Komite Claim Prop`
**Acuan nama kolom:** `STRUKTUR-TABEL-CLAIM-PROP.md` — berkas itu **tidak disunting** di sini.

> ⛔ **Berkas lama NOL disunting, tanpa kecuali** — termasuk `spec.md`, `issues/*`, berkas
> struktur, dan berkas relasi milik modul Komite. ⛔ Modul lain NOL.
> ⛔ Kode NOL · `CREATE TABLE` NOL · DDL NOL · nilai rahasia NOL.
> ⛔ **Tidak satu pun butir `[terbuka]` dinyatakan tertutup.**
> ⛔ **Nol nomor baris XML dikutip.** ⛔ **Tabel jejak audit tidak dinamai sendiri.**
> ⛔ **Kelima pertanyaan yang sedang menggantung tidak dijawab dan tidak diulang.**

---

## §A — Acuan yang mengikat

| # | Acuan | Dipakai di |
| --- | --- | --- |
| **A1** | Nama kolom dari `STRUKTUR-TABEL-CLAIM-PROP.md`; ejaan korpus = **bukti asal**, bukan nama (**AC 106**) | seluruh berkas |
| **A2** | Sepuluh nama `T_CLAIM_*` dipakai **apa adanya** | §B |
| **A3** | Lima tabel lintas-lini **tidak didefinisikan ulang**; sisi komite **sudah terkunci** dan punya berkas relasinya sendiri | §B relasi 14–17 |
| **A4** | Lini PROP: baris klaim **`CLMP-`** · baris komite **`TKMT-`** | §B relasi 1 · 2 |
| **A5** | Uang/persen/kurs: **angka desimal, 20 digit, 8 di belakang koma**, dihitung penuh tanpa pembulatan di tengah, tampil 4 desimal; **pencacah bilangan bulat biasa** | seluruh tabel bermuatan uang |

> ⚠️ **RALAT 07-10-2026** — kalimat lama: *"angka desimal, 20 digit, 8 di belakang koma"* Kolom uang, persen, share, dan kurs = **`NUMBER(38,10)`** (keputusan work owner 07-10-2026).
> Hitungan di Go penuh (`apd.Decimal`, nol float); di batas simpan Oracle membulatkan ke 10 angka di belakang
> koma. Bunyi "disimpan penuh tanpa pembulatan" berlaku untuk hitungan, bukan untuk kolom.

✅ **`STRUKTUR-TABEL-CLAIM-PROP.md` ADA** (662 baris) — dikonfirmasi di LANGKAH 0. Nol nama kolom
dikarang.

---

## §B6 — ⛔ JEBAKAN PENCARIAN, dan kesalahan saya sendiri yang terbukti

> ## ⚠️ RALAT — saya melakukan persis kesalahan yang diperingatkan
>
> Nama tabel `T_*` adalah **nama rancangan untuk Go**. Ia **tidak ada di korpus Pega** — dicari
> pun hasilnya **nol untuk kesepuluh tabel**, dan itu **bukan bukti apa pun**.
>
> ⛔ **Saya sudah melakukannya dua kali:**
>
> | Di mana | Yang saya tulis | Yang benar |
> | --- | --- | --- |
> | `RELASI-TABEL-KOMITE-CLAIM-PROP.md` §C5 | *"`T_VIEW_SUGGEST` — **TIDAK** menyentuh sisi komite … nol kemunculan di kedua modul"* | ⛔ **SALAH** |
> | `STRUKTUR-TABEL-CLAIM-PROP.md`, tabel lintas-lini | *"`T_VIEW_SUGGEST` — **tidak ketemu** di korpus … `[terbuka]` apa perannya"* | ⛔ **SALAH** |
>
> ⭐ **Halaman Pega-nya bernama `SuggestList`, dan ia ADA: 78 kemunculan di 4 berkas.**
>
> ⛔ Kedua berkas itu **tidak saya sunting** — ralatnya dicatat di sini saja.

### Peta nama tabel → nama halaman Pega, dipakai SEBELUM mencari

| Tabel rancangan | Halaman Pega | Kemunculan di korpus Claim Prop |
| --- | --- | --- |
| `T_WORK_CLAIM` | `pyWorkPage` *(objek kerja)* | 2 347 di 173 berkas |
| `T_GENERAL_CLAIM` | `pyWorkPage.ClaimData` | 2 338 di 114 berkas |
| **`T_VIEW_SUGGEST`** | ⭐ **`ClaimData.SuggestList`** | **78 di 4 berkas** |
| `T_CLAIM_ESTIMATION` | `ClaimData.EstimationList` | 105 di 21 berkas |
| `T_CLAIM_INTEREST` | `ClaimData.InterestList` | 101 di 17 berkas |
| `T_CLAIM_CLAIM_AMOUNT` | `ClaimData.ListClaimAmount` | 239 di 14 berkas |
| `T_CLAIM_SPREADING` | `ClaimData.SpreadingClaim` | 70 di 14 berkas |
| `T_CLAIM_BREAK_QS` | `ClaimData.SpreadingBreakQS` | 44 di 12 berkas |
| `T_CLAIM_FAC_RETRO` | `ClaimData.FacRetroList` | 56 di 6 berkas |
| `T_CLAIM_ADJUSTMENT` | `ClaimData.AdjustmentList` | 164 di 19 berkas |
| `T_CLAIM_ADJ_SPREADING` | `AdjustmentList(n).SpreadingAdjustment` | 65 di 11 berkas |
| `T_CLAIM_ADJ_QUOTA_SHARE` | `AdjustmentList(n).SpreadingQuotaShare` | 38 di 8 berkas |
| `T_CLAIM_ADJ_LOSS_ALLOCATION` | `AdjustmentList(n).LossAllocation` | 124 di 10 berkas |
| `T_KOMITE_KOMITELIST` | `KomiteList` | 16 di 3 berkas |

⛔ **Nama tabel `T_*` di korpus: NOL untuk keempat belasnya.** Itulah jebakannya.

### ⭐ `T_VIEW_SUGGEST` = daftar kronologi klaim — dan itu mengubah gambar

`[terverifikasi]` `ClaimData.SuggestList` diulang dan ditulisi oleh **dua penulis kronologi yang
sudah dikenal tiket 14**:

| Rule | Langkah | Yang dikerjakan |
| --- | --- | --- |
| `ASM-FW-GCNMFW-WORK!SETHISTORYKLAIMTREATY` | **1** | mengulang `pyWorkPage.ClaimData.SuggestList` dan menulis entri |
| `…WORK-CLAIMTREATY!INSERTCHRONOLOGY_DT` | — | penulis kedua, lewat data transform |
| `Section/InputAcceptation.xml` · `Section/OutstandingClaim.xml` | — | menampilkannya sebagai grid riwayat |

**Properti yang ditulis:** `.IsCedingConfirm` dan `.CommentSuggest`; `.PICSuggest` dibaca.

⭐ **`.IsCedingConfirm` menyimpan TINGKAT WEWENANG, bukan boolean** — nilainya
`"Admin Claim …"` · `"Claim Dept. Head"` · `"Operational Director"` · `"Technical Director"`.
Ini **persis** properti yang **AC 78** sebut *"properti berawalan `Is…` yang bukan boolean —
namanya berbohong dan isinya tidak dapat disaring"*.

⚠️ **Akibat yang wajib dicatat, tanpa menutup apa pun:** `T_VIEW_SUGGEST` adalah tabel
**lintas-lini**, dan **AC 4** menuntut jejak audit menjadi **tabel milik Claim Prop sendiri, bukan
menumpang struktur modul lain**. ⛔ Jadi keduanya hidup berdampingan dalam rancangan: `T_VIEW_SUGGEST`
**tempat sekarang**, dan tabel jejak audit **yang belum punya nama** tempat sesudah migrasi.
⛔ **Nama tabel itu tidak saya tetapkan** — lihat §C7.

---

## §B — Daftar relasi

### Ringkasan — **18 relasi**

| # | Induk | Anak | Kunci tamu | Kard. | ON DELETE | Index | Status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| **1** | `T_WORK_CLAIM` *(klaim)* | `T_WORK_CLAIM` *(komite)* | `COVER_KEY` | 1:N | **di Go** | biasa | ⛔ **tidak terbukti korpus** |
| **2** | `T_WORK_CLAIM` *(klaim)* | `T_GENERAL_CLAIM` | **tanpa kolom** — shared PK | 1:1 | **di Go** | PK saja | ⛔ **tidak terbukti korpus** |
| **3** | `T_GENERAL_CLAIM` | `T_CLAIM_ESTIMATION` | `CLAIM_ID` | 1:N | **CASCADE** | biasa | ✅ terbukti |
| **4** | `T_GENERAL_CLAIM` | `T_CLAIM_INTEREST` | `CLAIM_ID` | 1:N | **CASCADE** | biasa | ✅ terbukti |
| **5** | `T_GENERAL_CLAIM` | `T_CLAIM_CLAIM_AMOUNT` | `CLAIM_ID` | 1:N | **CASCADE** | biasa | ✅ terbukti |
| **6** | `T_GENERAL_CLAIM` | `T_CLAIM_SPREADING` | `CLAIM_ID` | 1:N | **CASCADE** | biasa | ✅ terbukti |
| **7** | `T_GENERAL_CLAIM` | `T_CLAIM_BREAK_QS` | `CLAIM_ID` | 1:N | **CASCADE** | biasa | ✅ terbukti |
| **8** | `T_GENERAL_CLAIM` | `T_CLAIM_FAC_RETRO` | `CLAIM_ID` | 1:N | **CASCADE** | biasa | ✅ terbukti |
| **9** | `T_GENERAL_CLAIM` | `T_VIEW_SUGGEST` | `CLAIM_ID` | 1:N | ⚠️ **JANGAN cascade** | biasa | ✅ terbukti — lewat `SuggestList` |
| **10** | `T_GENERAL_CLAIM` | `T_CLAIM_ADJUSTMENT` | `CLAIM_ID` | 1:N | **CASCADE** | biasa | ✅ terbukti |
| **11** | `T_CLAIM_ADJUSTMENT` | `T_CLAIM_ADJ_SPREADING` | `ADJUSTMENT_ID` | 1:N | **CASCADE** | biasa | ✅ terbukti |
| **12** | `T_CLAIM_ADJUSTMENT` | `T_CLAIM_ADJ_QUOTA_SHARE` | `ADJUSTMENT_ID` | 1:N | **CASCADE** | biasa | ✅ terbukti |
| **13** | `T_CLAIM_ADJUSTMENT` | `T_CLAIM_ADJ_LOSS_ALLOCATION` | `ADJUSTMENT_ID` | 1:N | **CASCADE** | biasa | ✅ terbukti |
| **14** | `T_CLAIM_ADJUSTMENT` | `T_GENERAL_KOMITE` | `ADJUSTMENT_ID` | 1:1 | di Go | ⭐ **UNIK** | ↗ **berkas Komite relasi 4** |
| **15** | `T_CLAIM_ADJUSTMENT` | `T_WORK_CLAIM` *(komite)* | `KOMITE_ID` | 1:1 | di Go | ⭐ **UNIK** | ↗ **berkas Komite relasi 5** — ⛔ tidak terbukti korpus |
| **16** | `T_WORK_CLAIM` *(komite)* | `T_GENERAL_KOMITE` | shared PK | 1:1 | di Go | PK saja | ↗ **berkas Komite relasi 2** |
| **17** | `T_GENERAL_KOMITE` | `T_KOMITE_KOMITELIST` | `DATA_KOMITE_ID` | 1:N | CASCADE | biasa | ↗ **berkas Komite relasi 3** |
| **18** ⭐ | `T_WORK_CLAIM` *(klaim)* | `DOCUMENT_CLAIM` | `IDPEGA` | 1:N | ⚠️ **JANGAN cascade** | biasa | ⭐ **RELASI BARU** — terbukti korpus |

⚠️ **Arah relasi 15 saya tulis `T_CLAIM_ADJUSTMENT → T_WORK_CLAIM`**, bukan seperti tertulis di
brief (`T_WORK_CLAIM → T_CLAIM_ADJUSTMENT`). Alasannya sama dengan yang sudah dilaporkan di berkas
relasi Komite: **kolom `KOMITE_ID` duduk di tabel penyesuaian**. ⛔ Dilaporkan, tidak dibetulkan
diam-diam.

### Rekap §B

| | |
| --- | --- |
| Jumlah relasi | **18** |
| Relasi **kedelapan belas** | ⭐ **ADA** — relasi **18**, `DOCUMENT_CLAIM.IDPEGA` |
| Relasi **TIDAK terbukti dari korpus** | **3** — relasi **1** (`COVER_KEY`) · **2** (shared PK `T_GENERAL_CLAIM`) · **15** (`KOMITE_ID`) |
| Bukti **tiga bagian lengkap** | **12 dari 18** — relasi **3–13** *(sebelas)* dan **18** |
| Dirujuk ke berkas Komite, tidak ditulis ulang | **4** — relasi **14 · 15 · 16 · 17** |

⭐ **Sama persis dengan sisi komite: 3 dari 6 tidak terbukti di sana, 3 dari 18 di sini — dan
ketiganya bertipe sama**, yaitu relasi yang lahir dari keputusan rancangan karena Pega tidak
memakai kolom untuk hubungan itu.

### B2 · B3 · B4 · B5 — per relasi

#### Relasi 1 — `COVER_KEY`

| | |
| --- | --- |
| **Bukti** | `[keputusan work owner 2026-09-18]`. ⛔ **Tidak ada kolom `COVER_KEY` di korpus** — Pega memakai mekanisme *cover/covered* bawaan. Jejak: `Activity/AddKomiteTreatyChild_ACT.xml` (`ASM-FW-GCNMFW-DATA-ADJUSTMENT!ADDKOMITETREATYCHILD_ACT`) langkah **12** menyalin halaman kasus anak |
| **DIISI** | **AddKomiteChild langkah 12**, saat kasus komite lahir |
| **DIBACA** | `Komite Claim Prop/Activity/KomitePostAdjustment.xml` langkah **4** |
| **Bila kosong** | ⛔ **tidak terbaca dari korpus** — langkah 4 tidak punya jalur kegagalan. `[terbuka]` |
| **ON DELETE · Index** | **di Go** — self-referencing, `CASCADE` berisiko · index **biasa** |

#### Relasi 2 — shared PK `T_WORK_CLAIM` ↔ `T_GENERAL_CLAIM`

| | |
| --- | --- |
| **Bukti** | `[keputusan work owner]`. ⛔ **Tidak ada bukti korpus** — di Pega `ClaimData` adalah **halaman tertanam** di dalam objek kerja, bukan tabel terpisah. Pemisahannya keputusan rancangan |
| **DIISI · DIBACA** | tidak berlaku — tidak ada kolom penyambung |
| **Bila kosong** | tidak mungkin — `ID`-nya PK kedua tabel |
| **ON DELETE · Index** | **di Go** — pasangan shared-PK lahir dan mati bersama · PK saja |

#### Relasi 3–8 · 10 — anak langsung kasus klaim

Ketujuhnya berbentuk sama: di Pega **daftar bersarang** di dalam `ClaimData`, di Go **tabel anak**
dengan `CLAIM_ID`.

| # | Tabel | DIISI *(rule · langkah)* | DIBACA | Bila kosong |
| --- | --- | --- | --- | --- |
| **3** | `T_CLAIM_ESTIMATION` | `Activity/AddEstimation_Act.xml` langkah **5.3** | `Activity/CountEstimation_Act.xml` langkah **3.1** · `Activity/CopyOldataCurr_act.xml` langkah **8** | ⛔ tidak terbaca — `CountEstimation_Act` mengulang daftar kosong tanpa pesan. `[terbuka]` |
| **4** | `T_CLAIM_INTEREST` | `Activity/AddInterest_act.xml` langkah **1** · `Activity/SetCurencyInterest_act.xml` langkah **1 · 3.1 · 5** | `Activity/CountTotalInsterest_Act.xml` langkah **5.2 · 7.1 · 11.1** | ⛔ tidak terbaca. `[terbuka]` |
| **5** | `T_CLAIM_CLAIM_AMOUNT` | `Activity/AddListClaimAmount.xml` langkah **6.1 · 7.1** | `Activity/CountPersen_act.xml` langkah **4.3.1 · 6.2.1** | ⚠️ `CountPersen_act` langkah **4.4** memberi total **nol** bila daftarnya kosong — **tanpa pesan**. `[terbuka]` |
| **6** | `T_CLAIM_SPREADING` | `Activity/CopyOldataCurr_act.xml` langkah **9.1** · `Activity/SetTreatyNameSpreading_Act.xml` langkah **6.1** | `Activity/CountSpreading_Act.xml` · `Activity/SaveOutstanding_Act.xml` langkah **26.1** | ⛔ tidak terbaca. `[terbuka]` |
| **7** | `T_CLAIM_BREAK_QS` | `Activity/CopyOldataCurr_act.xml` langkah **10.1** · `Activity/AddKomiteTreatyChild_ACT.xml` langkah **9.2** | `Activity/SaveOutstanding_Act.xml` langkah **27.1** | ⛔ tidak terbaca. `[terbuka]` |
| **8** | `T_CLAIM_FAC_RETRO` | `Activity/SaveAcceptationTreaty_Act.xml` langkah **9.3** | `Activity/AddKomiteTreatyChild_ACT.xml` langkah **8.1 · 10.1** · `Activity/PrintDLATreatyIn.xml` langkah **14.1 · 15.15** | ⛔ tidak terbaca. `[terbuka]` |
| **10** | `T_CLAIM_ADJUSTMENT` | `Activity/AddAdjustment_Act.xml` langkah **4 · 5 · 6** | banyak — a.l. `Activity/PrintFileAcceptance.xml` langkah **10**, `KomitePostAdjustment` langkah **3** | ⛔ tidak terbaca. `[terbuka]` |

**ON DELETE ketujuhnya: `CASCADE`** — ketujuh daftar itu **bagian dari kasus klaim** di Pega
*(halaman tertanam)*, jadi tidak punya arti tanpa induknya. **Index: biasa** pada `CLAIM_ID`.

#### Relasi 9 — `T_VIEW_SUGGEST` ⭐ terbukti lewat `SuggestList`

> ⚠️ **RALAT 07-10-2026** — relasi 9 **dibuat** (migrasi `532`, keputusan work owner "1 tabel aja gabung life dan non life"), dengan
> ON DELETE **CASCADE** seperti diagram sheet Claim Prop F38-F40 — bukan "JANGAN cascade" di tabel atas: riwayat
> hidup bersama klaimnya. CHECK `CK_VS_SATU_INDUK`: tepat satu dari `PREMIUM_LIST_ID` / `CLAIM_ID` terisi.

| | |
| --- | --- |
| **Bukti** | ✅ **terbukti korpus.** `Activity/SethistoryKlaimTreaty.xml` (`ASM-FW-GCNMFW-WORK!SETHISTORYKLAIMTREATY`) langkah **1** mengulang `pyWorkPage.ClaimData.SuggestList` dan menulis entri; `DataTransform/InsertChronology_DT.xml` penulis kedua |
| **DIISI** | **SethistoryKlaimTreaty langkah 1** — `.IsCedingConfirm` *(tingkat wewenang)* dan `.CommentSuggest` |
| **DIBACA** | `Section/InputAcceptation.xml` dan `Section/OutstandingClaim.xml` — grid riwayat. ⚠️ **AC 82**: grid itu **tidak punya konfigurasi pengurutan**; `SethistoryKlaimTreaty` langkah **1.11** menjalankan `Obj-Sort` atas `pyWorkPage.ClaimData` |
| **Bila kosong** | ⛔ tidak terbaca — grid tampil kosong tanpa pesan. `[terbuka]` |
| **ON DELETE** | ⚠️ **JANGAN cascade.** Ia **berwatak JEJAK** — menghapus klaim tidak boleh menghapus riwayatnya. ⚠️ Berbeda dari relasi 3–8·10 yang memang cascade |
| **Index** | **biasa** pada `CLAIM_ID`; ⚠️ **wajib ada index pada kolom tanggal** karena AC 82 menuntut urut tanggal |

#### Relasi 11 · 12 · 13 — anak baris penyesuaian

| # | Tabel | DIISI | DIBACA | Bila kosong |
| --- | --- | --- | --- | --- |
| **11** | `T_CLAIM_ADJ_SPREADING` | `Activity/CountSpreadingADJ_Act.xml` langkah **6.2 · 6.3** | `Activity/HitServiceToKasir_Act.xml` langkah **9.3** · `Activity/SendEmailKlaim.xml` langkah **8.2.1** | ⚠️ ⛔ tidak terbaca — **dan ia memberi makan jalur Kasir**. `[terbuka]` |
| **12** | `T_CLAIM_ADJ_QUOTA_SHARE` | `Activity/CountSpreadingADJ_Act.xml` langkah **7.2 · 7.3** | `Activity/PrintDLATreatyIn.xml` langkah **15.11.1** | ⛔ tidak terbaca. `[terbuka]` |
| **13** | `T_CLAIM_ADJ_LOSS_ALLOCATION` | `Activity/CountGrossAdjTreaty_Act.xml` langkah **5.2** | **CountGrossAdjTreaty_Act** langkah **2.1** | ⛔ tidak terbaca. `[terbuka]` |

**ON DELETE ketiganya: `CASCADE`** · **Index: biasa** pada `ADJUSTMENT_ID`.

⚠️ Relasi **11** dan **12** adalah **SEJAJAR**, bukan bersarang — `CountSpreadingADJ_Act` langkah
**6** dan **7** mengulang **dua daftar berbeda** pada tingkat yang sama.

#### Relasi 14 · 15 · 16 · 17 — sisi komite, **cukup dirujuk**

⛔ **Buktinya tidak ditulis ulang di sini.** Seluruhnya ada di
`.scratch/komite-claim-prop/RELASI-TABEL-KOMITE-CLAIM-PROP.md`:

| # di sini | Di berkas Komite | Ringkas |
| --- | --- | --- |
| **14** | relasi **4** | `ADJUSTMENT_ID` NOT NULL, index **UNIK**, **tanpa `REFERENCES`** — dua tabel tujuan menurut lini |
| **15** | relasi **5** | `KOMITE_ID` — ⛔ **tidak terbukti korpus**; yang ada hanya penanda `.IsKomite` |
| **16** | relasi **2** | shared PK baris komite ↔ `T_GENERAL_KOMITE` |
| **17** | relasi **3** | `DATA_KOMITE_ID`, 1:N, **CASCADE** |

#### Relasi 18 ⭐ — `DOCUMENT_CLAIM.IDPEGA` — **RELASI BARU**

| | |
| --- | --- |
| **Bukti** | ✅ **terbukti korpus, tiga bagian lengkap.** `Activity/InsertDocument_Act.xml` (`ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM!INSERTDOCUMENT_ACT`) langkah **3** menyusun halaman dokumen, langkah **5** menyimpannya (`Obj-Save`) |
| **Kolom yang terbaca** | `ID` · `TANGGAL` · **`IDPEGA`** · `NAMAFILE` · `MIME` · `KATEGORI_1` · `KATEGORI_2` · `NOAKSEP` · `NOPREKAS` · `PAYMENTDATE` · operator pembuat |
| **Kunci ke kasus klaim** | ⭐ **`IDPEGA`**, diisi dari `Param.IDPEGA` |
| **DIISI** | **InsertDocument_Act langkah 3**, disimpan langkah **5** |
| **DIBACA** | `Activity/GetBase64Attachment.xml` · `Activity/GetInvoiceAttachments.xml` · `Activity/GetPayAttachmentAdj_Act.xml` |
| **Bila kosong** | ⛔ tidak terbaca dari korpus — langkah 3 mengisi `IDPEGA` dari parameter **tanpa memeriksa**. `[terbuka]` |
| **ON DELETE** | ⚠️ **JANGAN cascade** — dokumen adalah **bukti**, dan sebagiannya sudah terbit keluar (DLA, acceptance note) |
| **Index** | **biasa** pada `IDPEGA` |

⚠️ **`IDPEGA` berisi kunci internal Pega** — sama persis dengan masalah `ID_PEGA` pada
`HISTORYAKSEPTASIPEGA` di sisi komite. Sesudah migrasi pengenal itu **tidak ada lagi**.
`[terbuka]` apa penggantinya.

⛔ **Kolom `DOCUMENT_CLAIM` tidak ditebak dan tidak ditulis DDL-nya** — kesebelas nama di atas
terbaca langsung dari penugasan langkah 3. Daftar kolom penuhnya tetap `[data DBA]`.

---

## §C — Relasi yang bukan kunci tamu

### C1 — ⭐ Empat properti kasus klaim yang ditulis modul Komite

`[terverifikasi]` Berkas relasi Komite §C3 membuktikan `KomitePostAdjustment` menulis ke kasus
klaim induk di **delapan titik**. Dua sudah punya rumah di sisi komite (`CLOSE_FILE`,
`RESERVED_CLAIM`). Enam sisanya menulis ke **empat properti**:

| Properti | Sudah jadi kolom di `STRUKTUR-TABEL-CLAIM-PROP.md`? | Penulis **di dalam Claim Prop** | Pembaca |
| --- | --- | --- | --- |
| `.IsAnyAcceptation` | ✅ **YA** — `IS_ANY_ACCEPTATION` di `T_GENERAL_CLAIM` | **1** — `Activity/CheckAnyAcceptationProp.xml` langkah **3** | `Section/InputAcceptation.xml` |
| `.ClaimData.IsSubjectivity` | ✅ **YA** — `IS_SUBJECTIVITY` di `T_GENERAL_CLAIM` | ⛔ **NOL** di Activity mana pun — **hanya modul Komite yang menulisnya** | layar akseptasi |
| `.ClaimData.FlagOnGoingCommitte` | ⛔ **DIBUANG** — ~~✅ **YA** — `FLAG_ON_GOING_COMMITTEE` di `T_GENERAL_CLAIM`~~; **tidak dibawa ke sistem baru** karena nilainya **dapat diturunkan** dari kasus komite yang berjalan. `[keputusan work owner]` **2026-09-19** · `spec.md` **14c**, **AC 129** | **1** — `Activity/AddKomiteTreatyChild_ACT.xml` langkah **3** (`"Send Commite"`) — **jejaknya dibiarkan terbaca** | layar akseptasi — **kini menurunkannya sendiri** |
| **`.AktifButton`** | ⛔ **BELUM** — **tidak ada kolomnya** | **5** — `AddAdjustment_Act` langkah **5** dan **12** · `AddKomiteTreatyChild_ACT` langkah **30** · `DeleteAjsutment_Act` langkah **3** · `SaveAcceptationTreaty_Act` langkah **15** | layar |

⭐ **Dua dari empat punya rumah · satu DIBUANG · satu belum.** Yang punya rumah:
`.IsAnyAcceptation` dan `.ClaimData.IsSubjectivity`. Yang **dibuang**:
`.ClaimData.FlagOnGoingCommitte` `[keputusan work owner]` 2026-09-19. Yang **belum**:
`.AktifButton`.

> ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip, tidak dihapus**: *"⭐ **Tiga dari empat sudah
> punya rumah. Satu belum: `.AktifButton`.**"* Ia benar sampai kolom `FLAG_ON_GOING_COMMITTEE`
> dibuang. ⛔ **Nol baris dihapus** — barisnya ditandai DIBUANG, bukan dicoret dari tabel.

⚠️ **`.ClaimData.IsSubjectivity` perlu perhatian tersendiri:** kolomnya ada, tetapi **nol penulis
di modul pemiliknya**. Satu-satunya penulis adalah modul Komite. Itu berarti kolom milik Claim Prop
yang **hanya diisi dari luar**.

#### Pilihan untuk `.AktifButton` — ⛔ **tidak saya tetapkan**

| Pilihan | Akibatnya |
| --- | --- |
| **A — tidak dijadikan kolom** | Konsisten dengan **AC 11**, yang membuang `pyExpanded` karena *"kebocoran lapisan UI ke penyimpanan tidak layak ditiru"*. `.AktifButton` bernama *button* dan bernilai `0`/`1` — bentuknya **persis** status tampilan. ⚠️ Tetapi ia ditulis **5 rule** termasuk sesudah penyimpanan dan sesudah penyerahan komite, jadi mungkin ia **penanda keadaan kasus**, bukan tampilan |
| **B — dijadikan kolom `T_GENERAL_CLAIM`** | Perilaku tertiru penuh; harganya satu kolom yang **namanya menyesatkan** dan harus dibetulkan menurut **AC 106** |
| **C — dijadikan kolom dengan nama sesuai isinya** | Sama dengan B, tetapi namanya menyebut keadaan *(misal: apakah tombol aksi masih boleh ditekan)*. ⚠️ Menuntut keputusan tentang **apa artinya**, yang belum terbaca dari korpus |

`[terbuka]` — **menunggu work owner.**

### C2 — Penunjuk posisional kedua: `CoverageSubscript`

`[terverifikasi]` Dari sisi Claim Prop, **satu titik saja** yang menyentuhnya:

```
AddKomiteTreatyChild_ACT langkah 13
    childPageKomite.Adjustment.CoverageSubscript  :=  .IndexCoverage
```

⚠️ **`.IndexCoverage` TIDAK PERNAH DIISI di modul Claim Prop.** Disisir `pyParamArray`,
`pyStepsCallParams`, `pyStepsJavaSource`, dan sel layar (`pyValue`) seluruh **329 berkas** —
**tidak ketemu di medan-medan itu** satu pun penulisnya. Sesuai aturan **I2**, kalimatnya dibatasi
pada medan yang disisir.

> ## ⛔ Daftar yang ditunjuknya: **TIDAK TERBACA DARI KORPUS**
>
> Karena pengisinya tidak ketemu, **daftar mana yang diindeks `CoverageSubscript` tidak dapat
> dipastikan**. ⛔ Tidak ditebak.

#### Calon kolom pengganti — **tiga pilihan**, ⛔ tidak ditetapkan

| Pilihan | Kalau `CoverageSubscript` menunjuk… | Kolom penggantinya |
| --- | --- | --- |
| **A** | **objek pertanggungan** — arti harfiah *coverage* | `T_CLAIM_INTEREST.ID`, disimpan sebagai kolom baru pada `T_GENERAL_KOMITE` |
| **B** | **baris estimasi** — `EstimationList` juga per objek dan per mata uang | `T_CLAIM_ESTIMATION.ID` |
| **C** | **tidak menunjuk apa pun yang masih hidup** — sisa dari bentuk data lama | **nol kolom**; penunjuknya dibuang saat migrasi |

⚠️ **Pilihan C tidak boleh diabaikan**: properti yang **tidak pernah diisi** di modul pemiliknya
adalah pola yang sudah dua kali terbukti sebagai sisa mati di proyek ini.

`[terbuka]` — **menunggu work owner.** ⛔ Ini **bukan** pertanyaan nomor 5 yang sedang menggantung;
yang menggantung adalah *penggantinya apa*, sedangkan yang saya laporkan di sini adalah fakta baru
bahwa **pengisinya pun tidak ketemu**.

### C3 — `DOCUMENT_CLAIM` — ⭐ kuncinya **TERBACA**

| | |
| --- | --- |
| **Rule penulis** | `Activity/InsertDocument_Act.xml` — `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM!INSERTDOCUMENT_ACT`, langkah **3** menyusun, langkah **5** `Obj-Save` |
| **Kuncinya ke kasus klaim** | ⭐ **`IDPEGA`** — diisi dari `Param.IDPEGA` |
| **Rule pembaca** | `GetBase64Attachment` · `GetInvoiceAttachments` · `GetPayAttachmentAdj_Act` · dipakai `PrintDLATreatyIn`, `PrintFileAcceptance`, `TryMakePLA_Act`, `CreatClaimAnalysis_Act` |
| **Induk untuk lini PROP** | ⛔ **masih `[terbuka]`** — **tidak saya tutup.** Yang kini terbaca hanyalah **kuncinya**; apakah induknya baris klaim atau baris penyesuaian bergantung isi `IDPEGA`, dan isi itu **kunci internal Pega** yang hilang sesudah migrasi |
| **Kolomnya** | ⛔ **`[data DBA]`** — sebelas nama di atas terbaca dari penugasan, tetapi daftar penuhnya tetap milik DBA. **Nol DDL ditulis** |

⚠️ **Jembatan ke `T_STORAGE_IMAGE` tidak terbaca.** `InsertDocument_Act` langkah **4** memanggil
rule penyimpanan berkas, tetapi **tidak ada kolom pada halaman dokumen yang memuat pengenal
berkas** di medan yang disisir. `[terbuka]`.

### C4 — `HISTORYAKSEPTASIPEGA` — ⛔ **Claim Prop TIDAK ikut menulis**

`[terverifikasi]` Disisir seluruh **329 berkas** `Claim Prop`, penyaringan **tidak peka huruf
besar-kecil** — **nol kemunculan**.

⭐ **Artinya relasi `HISTORYAKSEPTASIPEGA.ID_KOMITE` (berkas Komite relasi 6) adalah relasi
milik sisi komite saja, bukan lintas-lini sejauh kedua modul ini terbaca.**

⚠️ **`ID_PEGA` penggantinya:** ⛔ **tidak terbaca dari korpus.** Ia diisi dari `pzInsKey`, pengenal
internal Pega yang **hilang sesudah migrasi**. Calon penggantinya secara bentuk adalah
`T_WORK_CLAIM.ID` baris klaim (`CLMP-…`), tetapi **tidak ada apa pun di korpus yang menyatakan
pemetaan itu** — jadi ia tetap `[terbuka]`, **tidak saya tetapkan**.

### C5 — `T_STORAGE_IMAGE` — ⭐ dicoba dari sisi Claim Prop, dan **kuncinya terbaca sebagian**

`[terverifikasi]` Pernyataan SQL keempat rule-nya terbaca penuh:

| Rule | Yang terbaca |
| --- | --- |
| `…!RNM!INSERT_T_STORAGE_SQL` | kolom **`IMAGEID` · `URLPUBLIC` · `APPFOLDER` · `EXPDATE` · `FILENAME` · `APPNAME` · `STORAGE`** |
| `…!RNM!UPDATE_T_STORAGE_SQL` | menambah **`TANGGAL_UPLOAD`**; disaring `where imageid = …` |
| `…!RNM!GETLINKSTORAGE_SQL` | dibaca `where imageid = …` |
| `…!RNM!GENERATEIMAGEID_SQL` | `IMAGEID` dihasilkan dari **hash MD5 atas teks tetap + cap waktu** |

> ⭐ **Kunci tabelnya `IMAGEID`. Kunci ke KASUS KLAIM: ⛔ TIDAK ADA.**
>
> `[terverifikasi]` **Tidak satu pun kolom `t_storage_image` memuat pengenal klaim** — kedelapan
> kolom di atas seluruhnya tentang berkas, bukan tentang klaim. Ini **bukan** kesimpulan dari
> ketiadaan: keempat pernyataan SQL-nya terbaca **utuh**, jadi daftar kolomnya lengkap.
>
> ⚠️ Hubungannya ke klaim karena itu **harus lewat `DOCUMENT_CLAIM`** — tetapi kolom penyambungnya
> **tidak terbaca** *(C3)*. **BERHENTI di sini.** `[terbuka]`.

### C6 — Sepuluh halaman yang bukan tabel: penggantinya di Go

⛔ **Daftarnya tidak diubah.**

| Halaman Pega | Penggantinya di Go |
| --- | --- |
| `ClaimData.SpreadingRisk` | **query** ke master treaty saat dibutuhkan — **nol relasi** ke tabel klaim |
| `PaymentData` + `AcceptationList` + `CurrencyList` + `DetailPayment` | **query** ke layanan Kasir — **nol relasi**, data milik sistem lain |
| `ClaimData.SpreadingAdjustment` | **query agregat** atas `T_CLAIM_ADJ_SPREADING` — **nol relasi baru** |
| `ClaimData.SpreadingAdjustmentQS` | **query agregat** atas `T_CLAIM_ADJ_QUOTA_SHARE` — **nol relasi baru** |
| `AdjustmentList(n).FacRetroList` | **query** atas `T_CLAIM_FAC_RETRO` disaring baris penyesuaian — ⚠️ relasinya **sudah ada** *(relasi 8)*, tidak ditambah |
| `InterestListDtl` | **query** atas `T_CLAIM_INTEREST` — **nol relasi**, ia salinan |
| `ListTotalEstimation` + `TotalInterestInsured` | **variabel hasil hitung**, bukan relasi — dihitung saat dibaca |
| `ObjectList` + `ObjectItemList` | **variabel** penyusun tampilan dari `T_CLAIM_INTEREST` — **nol relasi** |

⭐ **Nol relasi tambahan** dari kesepuluh halaman itu. Empat menjadi **query**, dua menjadi
**query agregat**, dua menjadi **variabel**, dan dua *(SpreadingRisk, PaymentData)* menunjuk
**sistem lain**.

### C7 — Tabel jejak audit (AC 4) — relasinya `[terbuka]`

⛔ **Berkas struktur masih menandainya `[terbuka]` tanpa nama**, jadi relasinya ikut `[terbuka]`.
⛔ **Nama tidak saya tetapkan.**

**Bentuk relasinya yang sudah dapat dinyatakan, tanpa nama tabel:**

```
T_WORK_CLAIM (baris klaim)  ──1:N──►  «tabel jejak audit Claim Prop»
        kunci tamu: penunjuk ke baris klaim
        ON DELETE: ⚠️ JANGAN cascade — berwatak JEJAK
        index: biasa pada kunci tamu, DAN pada kolom tanggal (AC 82)
```

⚠️ **Hubungannya dengan relasi 9 wajib diputuskan bersama:** `T_VIEW_SUGGEST` adalah **tempat
jejak audit sekarang**, dan tabel ini adalah **tempatnya sesudah migrasi** menurut **AC 4**.
⛔ Apakah keduanya hidup berdampingan, atau yang satu menggantikan yang lain, **belum diputuskan**
— `[terbuka]`.

---

## §D — Periksa silang dengan sisi komite

### D1 — Relasi 1 · 4 · 5 di berkas Komite versus di sini

| Relasi | Berkas Komite | Berkas ini | Sama? |
| --- | --- | --- | --- |
| **`COVER_KEY`** | relasi **1** — 1:N, di Go, biasa, ⛔ tidak terbukti korpus | relasi **1** — **identik** | ✅ **SAMA** |
| **`ADJUSTMENT_ID`** | relasi **4** — 1:1, NOT NULL, **UNIK**, tanpa `REFERENCES`, di Go | relasi **14** — **identik**, dirujuk tidak ditulis ulang | ✅ **SAMA** |
| **`KOMITE_ID`** | relasi **5** — 1:1, **UNIK**, nullable, di Go, ⛔ tidak terbukti korpus | relasi **15** — **identik** | ✅ **SAMA** |

⭐ **Ketiganya sama, termasuk vonis "tidak terbukti korpus" pada dua di antaranya.** Nol
pertentangan. ⛔ Berkas Komite **tidak disunting**.

### D2 — `.IsKomite`: seluruh penulis dan pembaca — ⛔ **tidak saya simpulkan**

> ⛔ **Ini bahan keputusan work owner, bukan jawaban saya.** Pertanyaan *"apakah satu baris
> penyesuaian boleh masuk komite lebih dari sekali"* sedang menggantung dan **tidak saya jawab**.

`[terverifikasi]` Sensus **dua modul**, disisir `pyParamArray`, `pyStepsCallParams`,
`pyStepsJavaSource`, dan sel layar:

| Peran | Di mana | Nilai |
| --- | --- | --- |
| **PENULIS 1** | `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml` langkah **1** | `:= 1` — ⚠️ **tanpa memeriksa nilai sebelumnya** |
| **PENULIS 2** | `Komite Claim Prop/Activity/KomitePostAdjustment.xml` langkah **23** | `:= 0` — bergerbang penyetuju terakhir **dan** bukan persetujuan bersyarat |
| **PEMBACA 1** | `Claim Prop/Section/AdjustmentDetail.xml` | gerbang tampil/sunting |
| **PEMBACA 2** | `Claim Prop/Section/AdjustmentDetail_Section.xml` | gerbang tampil/sunting |
| **PEMBACA 3** | `Claim Prop/Section/InputAcceptation_Adjs.xml` | gerbang tampil/sunting |

⭐ **Penulis 2 · pembaca 3.** ⛔ **Nol rule yang memeriksa `.IsKomite` sebelum menulis `1`**, dan
**nol pembaca di luar layar** — ketiga pembacanya adalah gerbang tampilan, bukan penjaga aturan.

### D3 — Selisih daftar 17 properti: **keduanya benar, pertanyaannya berbeda**

`[terverifikasi]` Diperiksa dari sisi Claim Prop:

| Properti | Disalin ke kasus komite? | Muncul sebagai sel di `ShowTransfer`? |
| --- | --- | --- |
| `.Adjustment.CurrencyID` | ✅ **YA** — `AddKomiteTreatyChild_ACT` langkah **14** | ⛔ **TIDAK** |
| `.Adjustment.IndividualRiskRNM` | ✅ **YA** — **AddKomiteTreatyChild_ACT** langkah **14** | ✅ **YA** |

> ⭐ **Tidak ada yang salah — dua daftar menjawab dua pertanyaan berbeda.**
>
> Daftar di berkas struktur Komite menjawab *"apa yang **disalin** ke kasus komite"* — dan di situ
> **`CurrencyID` memang termasuk**. Hitungan saya di berkas relasi Komite menjawab *"apa yang
> **ditampilkan** di layar komite"* — dan di situ **`IndividualRiskRNM` yang termasuk**.
>
> ⚠️ **Yang keliru adalah label saya**, yang menyebutnya *"satu tukar"* seolah salah satunya salah.
> ⛔ Berkas struktur mana pun **tidak disunting**; ralat label ini dicatat di sini saja.

---

## §E — Pohon relasi dan penutup

### E1 — Pohon relasi — **lima tingkat**

```
T_WORK_CLAIM  (lintas-lini, satu baris per work object)
│   baris klaim   ID = CLMP-xxxxxx   COVER_KEY kosong      LINI = PROP
│   baris komite  ID = TKMT-xxxxxx   COVER_KEY = CLMP-…    LINI = PROP     [relasi 1]
│
├─ T_GENERAL_CLAIM                    shared PK, tanpa kolom penyambung     [relasi 2]
│   │
│   ├─ T_CLAIM_ESTIMATION            CLAIM_ID   1:N   CASCADE              [relasi 3]
│   ├─ T_CLAIM_INTEREST              CLAIM_ID   1:N   CASCADE              [relasi 4]
│   ├─ T_CLAIM_CLAIM_AMOUNT          CLAIM_ID   1:N   CASCADE              [relasi 5]
│   ├─ T_CLAIM_SPREADING             CLAIM_ID   1:N   CASCADE   +TREATY_ID [relasi 6]
│   ├─ T_CLAIM_BREAK_QS              CLAIM_ID   1:N   CASCADE   +TREATY_ID [relasi 7]
│   │     ⚠️ SEJAJAR dengan _SPREADING, bukan anaknya
│   ├─ T_CLAIM_FAC_RETRO             CLAIM_ID   1:N   CASCADE              [relasi 8]
│   ├─ T_VIEW_SUGGEST                 CLAIM_ID   1:N   ⚠️ JANGAN cascade    [relasi 9]
│   │     halaman Pega: ClaimData.SuggestList — daftar kronologi
│   │
│   └─ T_CLAIM_ADJUSTMENT            CLAIM_ID   1:N   CASCADE              [relasi 10]
│       │   + 7 kolom KOMITE_*   ·   KOMITE_ID → baris komite  1:1 UNIK     [relasi 15]
│       │
│       ├─ T_CLAIM_ADJ_SPREADING       ADJUSTMENT_ID 1:N CASCADE +TREATY_ID [relasi 11]
│       ├─ T_CLAIM_ADJ_QUOTA_SHARE     ADJUSTMENT_ID 1:N CASCADE +TREATY_ID [relasi 12]
│       │     ⚠️ SEJAJAR dengan _ADJ_SPREADING
│       └─ T_CLAIM_ADJ_LOSS_ALLOCATION ADJUSTMENT_ID 1:N CASCADE            [relasi 13]
│
├─ DOCUMENT_CLAIM                     IDPEGA     1:N   ⚠️ JANGAN cascade    [relasi 18] ⭐
│     kolom [data DBA] · jembatan ke T_STORAGE_IMAGE tidak terbaca
│
└─ «tabel jejak audit Claim Prop»     ⛔ BELUM PUNYA NAMA (AC 4)             [C7, terbuka]

   ┌────────── SISI KOMITE — ⛔ TIDAK didefinisikan ulang, lihat berkas Komite ──────────┐
   │  T_CLAIM_ADJUSTMENT.ID  ◄──1:1 UNIK──  T_GENERAL_KOMITE.ADJUSTMENT_ID  [relasi 14] │
   │  T_WORK_CLAIM (komite)   ──shared PK──► T_GENERAL_KOMITE                [relasi 16] │
   │  T_GENERAL_KOMITE        ──1:N────────► T_KOMITE_KOMITELIST             [relasi 17] │
   │  T_WORK_CLAIM (komite)   ──1:N────────► HISTORYAKSEPTASIPEGA   (komite saja)        │
   └───────────────────────────────────────────────────────────────────────────────────┘

   DIBACA dari luar, BUKAN anak, TIDAK ikut cascade:
      master treaty (SpreadingRisk) · sistem kasir (PaymentData) · T_STORAGE_IMAGE
```

**Tingkatnya lima:** `T_WORK_CLAIM` → `T_GENERAL_CLAIM` → `T_CLAIM_ADJUSTMENT` →
`T_CLAIM_ADJ_SPREADING` → baris di dalamnya.

⛔ **Tabel komite TIDAK didefinisikan ulang** — hanya sambungannya digambar.

### E2 — Butir `[terbuka]` berkas ini — **SEBELAS**

⛔ **Tidak satu pun dinyatakan tertutup.**

| # | Butir | Menunggu |
| --- | --- | --- |
| **1** | ⭐ **`.AktifButton` belum punya rumah kolom** — tiga pilihan disajikan di §C1 | **work owner** |
| **2** | ⭐ **`.IndexCoverage` tidak pernah diisi** di modul pemiliknya, jadi daftar yang ditunjuk `CoverageSubscript` **tidak terbaca**; tiga calon pengganti disajikan | **work owner** |
| **3** | **Induk `DOCUMENT_CLAIM` untuk lini PROP** — ⛔ butir lama, **tidak ditutup** | **work owner** |
| **4** | **Kolom penuh `DOCUMENT_CLAIM`** | **DBA** |
| **5** | **Jembatan `DOCUMENT_CLAIM` → `T_STORAGE_IMAGE`** tidak terbaca | korpus / DBA |
| **6** | **Pengganti `IDPEGA`** pada `DOCUMENT_CLAIM` sesudah `pzInsKey` hilang | **work owner** |
| **7** | **Pengganti `ID_PEGA`** pada `HISTORYAKSEPTASIPEGA` — bentuknya sama, pemiliknya sisi komite | **work owner** |
| **8** | **Hubungan `T_VIEW_SUGGEST` dengan tabel jejak audit AC 4** — berdampingan atau menggantikan? | **work owner** |
| **9** | **Nama tabel jejak audit** — ⛔ butir lama, **tidak dinamai sendiri** | **work owner** |
| **10** | **Apa yang terjadi bila daftar anak kosong** — sebelas relasi, **nol** yang punya penanganan di korpus | **work owner** |
| **11** | **`REFERENCES` pada `KOMITE_ID` dan `COVER_KEY`** — ⛔ butir lama | **DBA** |

### E3 — Relasi yang PALING RAWAN salah

⛔ **Ditunjuk, tidak diperbaiki.**

> ⚠️ **Paling rawan: relasi 9 — `T_VIEW_SUGGEST` sebagai anak `T_GENERAL_CLAIM` dengan
> `CLAIM_ID`.**

Ia baru saja lahir dari **koreksi atas kesalahan saya sendiri**, dan saya menetapkan bentuknya —
1:N lewat `CLAIM_ID`, jangan cascade — dari **satu** pembacaan: `SethistoryKlaimTreaty` langkah 1
mengulang `pyWorkPage.ClaimData.SuggestList`, yaitu daftar bersarang.

**Yang belum saya buktikan:** bahwa `T_VIEW_SUGGEST` memang **lintas-lini**. Ia masuk daftar lima
tabel bersama, jadi kemungkinan besar dipakai lini lain — dan kalau ya, `CLAIM_ID` saja **tidak
cukup**: ia perlu penanda lini, persis seperti `T_WORK_CLAIM`. ⛔ Saya **tidak memeriksanya**,
karena modul lain di luar lingkup.

⚠️ **Bentuknya sama dengan kesalahan yang baru saya ralat** — menyimpulkan dari satu sisi tanpa
menyisir sisi lain.

**Paling rawan kedua:** relasi **18**, `DOCUMENT_CLAIM.IDPEGA`. Kuncinya terbaca, tetapi **apa yang
dimuat `Param.IDPEGA` tidak saya lacak ke pemanggilnya** — kalau isinya kunci baris penyesuaian dan
bukan kunci klaim, induk relasi ini **salah**.

### E4 — Sebelum relasi tabel Claim Prop lengkap

| # | Yang kurang | Kenapa menahan |
| --- | --- | --- |
| **1** | **Vonis `T_VIEW_SUGGEST` lintas-lini atau tidak** *(E3)* | menentukan apakah relasi 9 perlu penanda lini |
| **2** | **Lacak `Param.IDPEGA` ke pemanggil `InsertDocument_Act`** *(E3)* | menentukan induk relasi 18 |
| **3** | **Nama tabel jejak audit** *(butir 9)* | relasinya tidak bisa ditulis tanpa nama |
| **4** | **Vonis `.CashLossList`** — butir terbuka berkas struktur | kalau ia tabel, pohon bertambah satu cabang |
| **5** | **Jawaban `.AktifButton`** *(butir 1)* | bukan relasi, tetapi mengubah isi `T_GENERAL_CLAIM` |
| **6** | **Keputusan `REFERENCES`** *(butir 11)* | menentukan apakah `ON DELETE` ditegakkan basis data atau Go |

### E5 — Pertanyaan BARU untuk work owner — **2**

> ⛔ **Kelima pertanyaan yang sedang menggantung TIDAK diulang** — `CountPersen_act` blok 5/5.1 ·
> proteksi periode polis · batas klaim TBA · baris penyesuaian masuk komite berulang ·
> pengganti penunjuk posisional kedua.

#### Pertanyaan 1 — `T_VIEW_SUGGEST` dan tabel jejak audit AC 4: berdampingan atau menggantikan?

**Apa yang ditanyakan.** `T_VIEW_SUGGEST` — halaman Pega `ClaimData.SuggestList` — adalah **tempat
jejak audit klaim berjalan sekarang**, dan ia **tabel lintas-lini**. **AC 4** menuntut jejak audit
menjadi **tabel milik Claim Prop sendiri, bukan menumpang struktur modul lain**. Jadi rancangan
memuat **dua** tempat untuk hal yang sama.

**Kenapa muncul.** Baru kelihatan hari ini, setelah kesalahan pencarian saya diralat. Sebelum ini
`T_VIEW_SUGGEST` dikira tidak dipakai sama sekali.

**Bedanya kalau A atau B.** **A — menggantikan:** seluruh entri kronologi pindah ke tabel baru milik
Claim Prop, `T_VIEW_SUGGEST` **tidak dipakai lini PROP**, dan migrasi memindahkan entri lama.
AC 4 terpenuhi bersih, tetapi lini lain mungkin masih memakai tabel bersama itu — perlu dipastikan.
**B — berdampingan:** `T_VIEW_SUGGEST` tetap untuk yang lintas-lini, tabel baru untuk yang khas
Claim Prop. ⚠️ Satu jenis peristiwa bisa tercatat **dua kali**, dan auditor melihat dua sumber.

**Apa yang tertahan.** Relasi **9**, relasi tabel jejak audit *(C7)*, dan aturan migrasi entri lama.

#### Pertanyaan 2 — Sebelas relasi anak, nol penanganan "daftar kosong": ditiru atau dijaga?

**Apa yang ditanyakan.** Untuk **kesebelas** relasi 1:N yang terbukti dari korpus, pertanyaan
*"apa yang terjadi bila kosong"* dijawab sama: ⛔ **tidak ada penanganannya di Pega**. Rule-rule
pembacanya mengulang daftar kosong dan menghasilkan nol, **tanpa pesan**. Contoh yang paling
menggigit: `CountPersen_act` langkah 4.4 memberi total estimasi **nol** bila daftar nilai klaim
kosong, dan angka nol itu **lolos ke hilir**.

**Kenapa muncul.** Ia muncul sebelas kali berturut-turut saat mengisi kolom *"bila kosong"* — pola
yang terlalu seragam untuk kebetulan.

**Bedanya kalau A atau B.** **A — ditiru:** Go juga menghasilkan nol diam-diam. Paritas penuh,
tetapi **nol adalah angka uang yang sah**, sehingga kesalahan data tidak akan pernah terlihat.
**B — dijaga:** daftar kosong pada titik yang menghasilkan uang menjadi **kegagalan yang terlihat**,
sejalan dengan **ADR-0015**. ⚠️ Ini **penyimpangan sadar** dan menambah jumlahnya.

**Apa yang tertahan.** Kolom *"bila kosong"* pada kesebelas relasi, dan apakah perlu AC baru.

### E6 — ⭐ Yang seharusnya dikerjakan tetapi TIDAK diperintahkan blok ini

⛔ **Disebutkan, tidak dikerjakan.**

| # | Butir |
| --- | --- |
| **1** | **Meralat dua berkas yang memuat kesalahan `T_VIEW_SUGGEST`** — `RELASI-TABEL-KOMITE-CLAIM-PROP.md` §C5 dan `STRUKTUR-TABEL-CLAIM-PROP.md`; keduanya menyatakan tabel itu tidak dipakai |
| **2** | **Menambahkan `T_VIEW_SUGGEST` ke `STRUKTUR-TABEL-CLAIM-PROP.md`** berikut ketiga properti yang terbaca (`IsCedingConfirm`, `CommentSuggest`, `PICSuggest`) |
| **3** | **Melacak `Param.IDPEGA`** ke seluruh pemanggil `InsertDocument_Act` — menentukan induk relasi 18 |
| **4** | **Memeriksa `T_VIEW_SUGGEST` di modul lain** — menentukan apakah ia perlu penanda lini |
| **5** | **Menyensus kolom `.CashLossList`** sebelum vonisnya diputuskan |
| **6** | **Memeriksa apakah tiket 14 perlu ditambal** — ia menyebut jejak audit tanpa menyebut `SuggestList` maupun `T_VIEW_SUGGEST` |
| **7** | **Menuliskan bentuk `T_GENERAL_CLAIM` secara penuh** — berkas ini hanya menyentuh kolom khas PROP |

---

## Lampiran — dua ralat atas pekerjaan saya sendiri

⛔ Berkas yang memuatnya **tidak disunting**. Dicatat di sini saja.

| # | Di mana | Yang tertulis | Yang benar |
| --- | --- | --- | --- |
| **1** | `RELASI-TABEL-KOMITE-CLAIM-PROP.md` §C5 · `STRUKTUR-TABEL-CLAIM-PROP.md` | *"`T_VIEW_SUGGEST` tidak ketemu / tidak menyentuh sisi komite"* | ⛔ **SALAH** — dicari dengan **nama tabel rancangan**, bukan nama halaman Pega. Halamannya `ClaimData.SuggestList`, **78 kemunculan di 4 berkas**, dan ia **daftar kronologi klaim** |
| **2** | `RELASI-TABEL-KOMITE-CLAIM-PROP.md` §C1 | daftar 17 properti disebut punya *"satu tukar"* terhadap berkas struktur | ⚠️ **Label keliru.** Keduanya benar — berkas struktur menjawab *"apa yang **disalin**"*, hitungan saya menjawab *"apa yang **ditampilkan**"*. Tidak ada yang salah, hanya dua pertanyaan berbeda |
