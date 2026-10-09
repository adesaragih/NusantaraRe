# Relasi Tabel — Claim Fac In
## Enam belas relasi: sebelas sisi klaim, empat sisi komite, satu kronologi

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
> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Tujuh tabel dipakai bersama; lima di
> antaranya berkunci asing GANDA.**"* → kunci asing ganda **DITOLAK** work owner (OQ-CFI-01, 09-10-2026); baris FAC
> mengisi `CLAIM_ID` **dan** induk tingkatnya, nol `CHECK`, nol `MODIFY`. Rincian di RALAT §B dan §D.

> **SENSUS BERKAS INI** — *(dihitung ulang sesudah relasi ke-16, 2026-09-20)*
>
> ⛔ **Jendelanya disebut: seluruh berkas DIKURANGI blok sensus ini sendiri** — **227 baris**
> menurut cacah baris-baru, **228** menurut pemisahan teks. ⚠️ **Selisih satu, sebabnya diketahui:**
> berkas berakhir dengan baris baru. ⭐ **Angka yang dipakai: 227.**
>
> bab `## ` **10** · tabel **8** · ⭐ **relasi 16** · `[terverifikasi]` **2** ·
> `[keputusan work owner]` **3** · `[terbuka]` **1** · ⚠️ **13** · ⭐ **56** · ⛔ **27** ·
> ✅ **5** · ★ **4**
>
> ⭐ **Dihitung DUA CARA** — *(a)* teks utuh, *(b)* baris demi baris. ✅ **Sepakat pada kesepuluh
> angka penanda.**
>
> ⛔ ⭐ **Awalan tabel yang lama dan awalan lini Life: nol kemunculan** di berkas ini.
>
> ⚠️ **Bergeser dari keadaan sebelumnya** — badan **201 → 227**, relasi **15 → 16**.
>
> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Angka yang dipakai: 227.**"* → sensus ini
> **tidak dihitung ulang**; ia mengukur badan berkas 20-09-2026, dan blok RALAT 10-10-2026 berada di luar jendelanya.

---

## §A — Acuan yang mengikat

⭐ Berkas ini memerikan **relasi** — arah, kardinalitas, dan perilaku hapus. ⛔ **Bukan DDL, bukan
tipe data, bukan daftar kolom.**

| Acuan | Isinya |
| --- | --- |
| `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` | ⭐ sembilan tabel lini FAC, asalnya di Pega, dan lima kunci ganda |
| `claim-prop\RELASI-TABEL-CLAIM-PROP.md` | ⭐ bentuk berkas ini ditiru dari sana |
| `komite-claim-prop\STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md` | ⭐ `T_GENERAL_KOMITE` dan `T_KOMITE_KOMITELIST` **sudah terkunci** di sana |

> ⛔ **RALAT 10-10-2026.** Sel lamanya dikutip utuh, tidak dihapus: *"`claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md`"* ·
> *"`claim-prop\RELASI-TABEL-CLAIM-PROP.md`"* · *"`komite-claim-prop\STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md`"* → letak
> folder `.scratch` lama; kini `APP_RNM/modul/claimfacin/docs/`, `APP_RNM/modul/claimprop/docs/`, dan
> `APP_RNM/modul/komiteclaimprop/docs/`. Kolom yang **mengikat** ada di lampiran pengikat STRUKTUR masing-masing.

**Cara membaca kolomnya:**

| Kolom | Artinya |
| --- | --- |
| **Kunci** | nama kolom kunci tamu pada tabel anak |
| **Kardinalitas** | `1:N` atau `1:1` |
| **Ditegakkan** | `CASCADE` = kunci tamu sungguhan dengan hapus berantai · **`di Go`** = ⛔ **tanpa `REFERENCES`**, dijaga lapisan layanan · **`—`** = tanpa kolom kunci tamu sama sekali |
| **Catatan** | ⭐ penanda kunci ganda dan keunikan |

⭐ **Aturan indeks yang berlaku untuk seluruh enam belas relasi:**

> ⭐ **Seluruh kunci tamu WAJIB ber-index.** ⭐ Yang bertanda **UNIK** wajib **index UNIK**.

---

## §B — Sebelas relasi sisi klaim

| # | Induk | Anak | Kunci | Kardinalitas | Ditegakkan | Catatan |
| ---: | --- | --- | --- | --- | --- | --- |
| **1** | `T_WORK_CLAIM` | `T_WORK_CLAIM` | `COVER_KEY` | `1:N` | ⭐ **di Go** | ⭐ **self** · **nullable** |
| **2** | `T_WORK_CLAIM` | `T_GENERAL_CLAIM` | **SHARED PK** | `1:1` | **—** | ⛔ **tanpa kolom kunci tamu** |
| **3** | `T_GENERAL_CLAIM` | `T_CLAIM_OBJECT` | `CLAIM_ID` | `1:N` | `CASCADE` | ★ khas FAC |
| **4** | `T_CLAIM_OBJECT` | `T_CLAIM_OBJECT_ITEM` | `OBJECT_ID` | `1:N` | `CASCADE` | ★ khas FAC |
| **5** | `T_CLAIM_OBJECT_ITEM` | `T_CLAIM_ESTIMATION` | `OBJECT_ITEM_ID` | `1:N` | `CASCADE` | ⭐ **kunci ganda** |
| **6** | `T_CLAIM_OBJECT_ITEM` | `T_CLAIM_SPREADING` | `OBJECT_ITEM_ID` | `1:N` | `CASCADE` | ⭐ **kunci ganda** |
| **7** | `T_CLAIM_OBJECT_ITEM` | `T_CLAIM_BREAK_QS` | `OBJECT_ITEM_ID` | `1:N` | `CASCADE` | ⭐ **kunci ganda** |
| **8** | `T_CLAIM_OBJECT_ITEM` | `T_CLAIM_ADJUSTMENT` | `OBJECT_ITEM_ID` | `1:N` | `CASCADE` | ⭐ **kunci ganda** |
| **9** | `T_CLAIM_ADJUSTMENT` | `T_CLAIM_ADJ_SPREADING` | `ADJUSTMENT_ID` | `1:N` | `CASCADE` | kunci **SAMA** dua lini |
| **10** | `T_CLAIM_ADJUSTMENT` | `T_CLAIM_ADJ_QUOTA_SHARE` | `ADJUSTMENT_ID` | `1:N` | `CASCADE` | kunci **SAMA** dua lini |
| **11** | `T_CLAIM_ADJUSTMENT` | `T_CLAIM_FAC_RETRO` | `ADJUSTMENT_ID` | `1:N` | `CASCADE` | ⭐ **kunci ganda** |

> ⛔ **RALAT 10-10-2026.** Sel lamanya dikutip utuh, tidak dihapus: *"⭐ **kunci ganda**"* (relasi 5–8 dan 11) →
> relasinya sendiri **dibangun** (FK `ON DELETE CASCADE`, ber-index, migrasi `561`–`567`), tetapi **bukan kunci ganda**
> (OQ-CFI-01): `CLAIM_ID` di kelima tabel itu tetap NOT NULL dan **juga diisi baris FAC** — jadi setiap baris FAC
> menggantung pada **dua** relasi sekaligus: ke tingkatnya (`OBJECT_ITEM_ID` / `ADJUSTMENT_ID`) **dan** ke
> `T_GENERAL_CLAIM` lewat `CLAIM_ID` (CASCADE, DDL Claim Prop `521`–`528`); pengecualian: baris retro tingkat klaim
> (`ClaimData.FacRetroTreaty`) ber-`ADJUSTMENT_ID` kosong. Relasi 4 pun begitu: `T_CLAIM_OBJECT_ITEM`
> membawa `CLAIM_ID` (NOT NULL, CASCADE ke `T_GENERAL_CLAIM`) di samping `OBJECT_ID`. Nol `CHECK` "tepat satu induk".

### ⚠️ Catatan atas relasi 1 dan 2

**Relasi 1** — ⭐ `T_WORK_CLAIM` menunjuk **dirinya sendiri** lewat `COVER_KEY`, dan kolomnya
**nullable**: baris **induk** tidak punya penutup. ⛔ Ditegakkan **di Go**, bukan dengan
`REFERENCES` — ⭐ sebab kunci tamu yang menunjuk tabel sendiri **mempersulit penghapusan berantai**.

**Relasi 2** — ⭐ `T_GENERAL_CLAIM` memakai **kunci utama yang sama** dengan `T_WORK_CLAIM`. ⛔ **Tak
ada kolom kunci tamu** — keduanya satu baris logis yang dipecah dua tabel.

---

## §C — Empat relasi sisi Komite Fac In

| # | Induk | Anak | Kunci | Kardinalitas | Ditegakkan | Catatan |
| ---: | --- | --- | --- | --- | --- | --- |
| **12** | `T_CLAIM_ADJUSTMENT` | `T_GENERAL_KOMITE` | `ADJUSTMENT_ID` | `1:1` | ⭐ **di Go** | ⭐ **UNIK** · **NOT NULL** · ⛔ **tanpa `REFERENCES`** |
| **13** | `T_WORK_CLAIM` | `T_CLAIM_ADJUSTMENT` | `KOMITE_ID` | `1:1` | ⭐ **di Go** | ⭐ **UNIK** · **nullable** · ⭐ **pintasan tampilan** |
| **14** | `T_WORK_CLAIM` | `T_GENERAL_KOMITE` | **SHARED PK** | `1:1` | **—** | ⭐ baris **`KMT-`** |
| **15** | `T_GENERAL_KOMITE` | `T_KOMITE_KOMITELIST` | `DATA_KOMITE_ID` | `1:N` | `CASCADE` | ⭐ **bersama lini PROP** |

---

## §C2 — ⭐ Relasi ke-16: kronologi klaim

`[keputusan work owner]` **K8 · 2026-09-20**

| # | Induk | Anak | Kunci | Kardinalitas | Ditegakkan | Catatan |
| ---: | --- | --- | --- | --- | --- | --- |
| ⭐ **16** | `T_GENERAL_CLAIM` | `T_VIEW_SUGGEST` | `CLAIM_ID` | `1:N` | ⚠️⚠️ **JANGAN cascade** | ⭐ **lintas-lini** · sumber `ClaimData.SuggestList` |

⚠️⚠️ **"Jangan cascade" adalah satu-satunya perilaku hapus yang berbeda di seluruh berkas ini**,
dan sebabnya wajib dibaca: ⭐ **kronologi adalah JEJAK.** ⛔ Jejak yang ikut terhapus bersama
induknya **berhenti menjadi jejak** — menghapus klaim **tidak boleh** menghapus riwayat siapa
mengerjakan apa.

`[terverifikasi]` Bentuk ini **sudah berlaku** di lini PROP — relasi **9** di
`claim-prop\RELASI-TABEL-CLAIM-PROP.md`, bertanda **JANGAN cascade** dengan alasan yang sama.

⭐ **Penulisnya dua:** `Activity\SethistoryKlaimTreaty` langkah **1** yang mengulang
`ClaimData.SuggestList`, dan `DataTransform\InsertChronology_DT` sebagai penulis kedua.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⚠️⚠️ **"Jangan cascade" adalah satu-satunya
> perilaku hapus yang berbeda di seluruh berkas ini**"* dan *"`[terverifikasi]` Bentuk ini **sudah berlaku** di lini
> PROP — relasi **9** di `claim-prop\RELASI-TABEL-CLAIM-PROP.md`, bertanda **JANGAN cascade** dengan alasan yang sama."*
> → relasi 16 **CASCADE**: migrasi Claim Prop `532` membuat `T_VIEW_SUGGEST.CLAIM_ID` → `T_GENERAL_CLAIM` **ON DELETE
> CASCADE** (`FK_VS_CLAIM`, ber-index `IDX_VS_CLAIM`, `CHECK` satu induk `CK_VS_SATU_INDUK` bersama `PREMIUM_LIST_ID`);
> Claim Fac In memakai tabel itu apa adanya; relasi 16 kini berperilaku sama dengan relasi `CASCADE` lain. Modul ini
> tidak menghapus baris `T_GENERAL_CLAIM`; penutupan kasus hanya mengisi `STATUS_WORK`.

### ⚠️ Catatan atas relasi 12 dan 13 — keduanya saling menunjuk

⭐ **Relasi 12** menunjuk **dari komite ke penyesuaian**: tiap kasus komite lahir dari **tepat satu**
penyesuaian, karena itu **UNIK** dan **NOT NULL**.

⭐ **Relasi 13** menunjuk **kembali**, dari penyesuaian ke kasus komite — ⭐ **pintasan tampilan**,
supaya layar klaim dapat menampilkan status komite **tanpa satu putaran kueri tambahan**.

⛔⛔ **Keduanya ditegakkan di Go, bukan dengan `REFERENCES`, dan sebabnya wajib dibaca:** ⚠️ dua
kunci tamu yang **saling menunjuk** membuat urutan penyisipan mustahil — ⛔ baris mana pun yang
ditulis lebih dulu akan melanggar kunci tamu yang satunya. ⭐ Karena itu keduanya dijaga **lapisan
layanan**, dalam **satu transaksi**.

⚠️ **Akibat yang wajib diketahui pembangun:** ⛔ **basis data TIDAK akan menolak** kasus komite yang
menggantung tanpa penyesuaian, maupun penyesuaian yang menunjuk kasus komite yang sudah terhapus.
⭐ **Penjagaannya ada di lapisan layanan, dan wajib punya uji.**

### ⭐ Relasi 14 — baris `KMT-`

⭐ Kasus komite **bukan tabel tersendiri**: ia **baris di `T_WORK_CLAIM`** berawalan **`KMT-`**,
dengan `T_GENERAL_KOMITE` memakai **kunci utama yang sama**. ⭐ Bentuknya **sama persis** dengan
relasi 2 pada sisi klaim.

### ⭐ Relasi 15 — dipakai bersama lini PROP

⛔ `T_KOMITE_KOMITELIST` **sudah terkunci** sebagai tabel lintas-lini di
`komite-claim-prop\STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md`. ⭐ Lini FAC memakainya **apa adanya** —
⛔ **nol kolom baru, nol tabel baru.**

---

## §D — Lima kunci ganda, dikumpulkan

⭐ **Lima dari enam belas relasi menuju tabel berkunci ganda.** ⛔ Pada lini **FAC** hanya **satu**
dari dua kolom induk yang terisi; kolom yang satunya **selalu kosong**.

| # relasi | Tabel anak | Kolom terisi di **FAC** | Kolom kosong di FAC *(dipakai PROP)* |
| ---: | --- | --- | --- |
| **5** | `T_CLAIM_ESTIMATION` | `OBJECT_ITEM_ID` | `CLAIM_ID` |
| **6** | `T_CLAIM_SPREADING` | `OBJECT_ITEM_ID` | `CLAIM_ID` |
| **7** | `T_CLAIM_BREAK_QS` | `OBJECT_ITEM_ID` | `CLAIM_ID` |
| **8** | `T_CLAIM_ADJUSTMENT` | `OBJECT_ITEM_ID` | `CLAIM_ID` |
| ⚠️ **11** | `T_CLAIM_FAC_RETRO` | **`ADJUSTMENT_ID`** | `CLAIM_ID` |

⭐ **Penjaganya adalah `CHECK` "tepat satu terisi"**, diuraikan di
`STRUKTUR-TABEL-CLAIM-FACIN.md` §2, ⭐ dan polanya **sudah berlaku** sejak
`T_VIEW_SUGGEST` — `[keputusan work owner]` **2026-09-18**.

> ⛔ **RALAT — 2026-09-20.** Berkas ini sempat menulis nama tabel itu sebagai
> **`T_CLAIM_VIEW_SUGGEST`**. ⭐ **Yang benar `T_VIEW_SUGGEST`**, tanpa sisipan `CLAIM_`.
> ⚠️ **Sebabnya perlu dicatat:** awalan itu **ditambahkan karena pencocokan pola**, bukan karena
> dibaca dari sumbernya — ⛔ tepat jenis kekeliruan yang aturan sensus ada untuk mencegah.

⚠️⚠️ **Relasi 11 adalah satu-satunya yang induknya beda TINGKAT antar lini** — di FAC menggantung
pada **penyesuaian**, di PROP pada **klaim**. ⛔ `[terbuka]` — lihat register §6
`STRUKTUR-TABEL-CLAIM-FACIN.md`.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Penjaganya adalah `CHECK` "tepat satu
> terisi"**, diuraikan di `STRUKTUR-TABEL-CLAIM-FACIN.md` §2"* dan sel tabel *"Kolom kosong di FAC *(dipakai PROP)*"* →
> **tidak dibangun** (OQ-CFI-01, 09-10-2026): nol `CHECK`, `CLAIM_ID` **tidak kosong di FAC** — baris FAC mengisi
> `CLAIM_ID` **dan** `OBJECT_ITEM_ID` (`563`–`566`); `T_CLAIM_FAC_RETRO` FAC mengisi `CLAIM_ID` + `ADJUSTMENT_ID`
> (`567`). Relasi 11 bukan lagi `[terbuka]`: K5 menutupnya (FAC berinduk adjustment, `CLAIM_ID` tetap diisi). Lihat
> RALAT 10-10-2026 di `STRUKTUR-TABEL-CLAIM-FACIN.md` §2.

---

## §E — Pohon relasi lini FAC

```
T_WORK_CLAIM                                    LINI = FAC · CLM- / KMT-
  └ T_GENERAL_CLAIM              SHARED PK
      └ T_CLAIM_OBJECT           CLAIM_ID              ★ khas FAC
          └ T_CLAIM_OBJECT_ITEM  OBJECT_ID             ★ khas FAC
              ├ T_CLAIM_ESTIMATION      OBJECT_ITEM_ID
              ├ T_CLAIM_SPREADING       OBJECT_ITEM_ID
              ├ T_CLAIM_BREAK_QS        OBJECT_ITEM_ID
              └ T_CLAIM_ADJUSTMENT      OBJECT_ITEM_ID
                    ├ T_CLAIM_ADJ_SPREADING     ADJUSTMENT_ID
                    ├ T_CLAIM_ADJ_QUOTA_SHARE   ADJUSTMENT_ID
                    └ T_CLAIM_FAC_RETRO         ADJUSTMENT_ID
```

> ⛔ **RALAT 10-10-2026.** Baris pohon lamanya dikutip utuh, tidak dihapus:
> `T_WORK_CLAIM                                    LINI = FAC · CLM- / KMT-` →
> nilainya **`T_WORK_CLAIM.LINI = 'FACIN'`** (sama dengan `STS_KLAIM` `EMAILKOMITE` FACIN); lini disaring lewat `LINI`,
> tidak pernah lewat awalan (OQ-CFI-02). Pohon ini juga tidak menggambar bahwa `T_CLAIM_OBJECT_ITEM`,
> `T_CLAIM_ESTIMATION`, `T_CLAIM_SPREADING` (baris `JENIS` `POLIS` / `KLAIM`), `T_CLAIM_BREAK_QS`, `T_CLAIM_ADJUSTMENT`,
> dan `T_CLAIM_FAC_RETRO` **juga** membawa `CLAIM_ID` (CASCADE ke `T_GENERAL_CLAIM`), serta retro tingkat klaim
> (`ADJUSTMENT_ID` kosong) langsung di bawah `T_GENERAL_CLAIM`.

⭐ **Sisi komite menggantung di dua titik:**

```
T_CLAIM_ADJUSTMENT  ──(12, UNIK NOT NULL)──▶  T_GENERAL_KOMITE
T_CLAIM_ADJUSTMENT  ◀──(13, UNIK nullable)──  pintasan KOMITE_ID
T_WORK_CLAIM        ──(14, SHARED PK)──────▶  T_GENERAL_KOMITE      baris KMT-
T_GENERAL_KOMITE    ──(15, DATA_KOMITE_ID)─▶  T_KOMITE_KOMITELIST
```

⭐ **Kedalamannya lima tingkat** pada sisi klaim — ⚠️ **satu tingkat lebih dalam daripada lini
PROP**, sebab lini FAC menyisipkan **objek** dan **item objek** yang tidak ada di PROP.

---

## §F — Yang BUKAN relasi kunci tamu

⛔ **Tiga hal yang mudah dikira relasi, tetapi bukan:**

| Yang dikira relasi | Kenyataannya |
| --- | --- |
| `ObjectItemList[].SpreadingList` | ⛔ **bukan tabel** — ⭐ himpunan bagian dari `T_CLAIM_SPREADING`, diambil dengan `SELECT DISTINCT`. Lihat `STRUKTUR-TABEL-CLAIM-FACIN.md` §3 butir 7 |
| `ObjectList[].CoverageList` dan lima saudaranya | ⛔ **salinan polis** — ⭐ dibaca dari **modul polis**, bukan dari tabel klaim |
| `T_CLAIM_INTEREST` · `T_CLAIM_CLAIM_AMOUNT` · `T_CLAIM_ADJ_LOSS_ALLOCATION` | ⛔ **tidak dipakai lini FAC** — ⭐ ketiganya milik lini PROP, disebut hanya sebagai catatan |

> ⛔ **RALAT 10-10-2026.** Sel lamanya dikutip utuh, tidak dihapus: *"⛔ **bukan tabel** — ⭐ himpunan bagian dari
> `T_CLAIM_SPREADING`, diambil dengan `SELECT DISTINCT`."* → `SpreadingList` memang bukan tabel **sendiri**, tetapi
> **disimpan** sebagai baris `T_CLAIM_SPREADING` ber-`JENIS = 'POLIS'` (migrasi `564`), bukan diturunkan dengan `SELECT
> DISTINCT` — lihat RALAT 10-10-2026 di `STRUKTUR-TABEL-CLAIM-FACIN.md` §3. Baris ketiga juga kurang: selain ketiga
> tabel itu, `T_CLAIM_LOSS_ALLOCATION` (Claim Prop `524`) dan tabel Claim Non Prop `T_CLAIM_NP_LOSS_ALLOC` ·
> `T_CLAIM_NP_XOL_ALLOC` · `T_CLAIM_NP_CLAIM_ACCEPT` (`608`–`610`) ada, dan **tak satu pun dipakai lini FAC**.

---

## §G — Periksa silang dengan sisi komite

`[terverifikasi]` ⭐ **Enam medan `KomiteList` lini FAC cocok satu per satu** dengan kolom
`T_KOMITE_KOMITELIST` yang sudah terkunci — ⛔ **nol medan tersisa, nol kolom baru.**

⭐ **Dan enam perbedaan antara komite lini FAC dan lini PROP seluruhnya PERILAKU** — diuraikan di
`STRUKTUR-TABEL-CLAIM-FACIN.md` §5b. ⛔ **Tak satu pun menuntut kolom atau relasi yang berbeda.**

⭐ **Karena itu §C berisi empat relasi yang sama bentuknya untuk kedua lini**, dan yang berbeda
hanyalah **ke tabel penyesuaian mana** `ADJUSTMENT_ID` menunjuk.

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
