# Grilling Ronde 5 — Claim Prop

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Claim Prop\` (READ-ONLY)
**Sasaran:** empat keputusan work owner (§A), `Local.Value` yang menggantung (§B), ketelitian lima
titik kelas Adjustment dan hulu `.TotalClaim` (§C), tiga butir goyah sisa (§D).

> ⛔ **Nilai rahasia NOL.** Aturan §A1 ronde 4 tetap berlaku penuh.
> ⛔ **Tiket NOL disentuh.** Keempatnya sudah ditambal di ronde 4.
> ⛔ Yang disunting hanya **`spec.md` (§E)** dan **`grilling-ronde-2.md` (satu baris, §A4)**.

> ⚠️ **Tabrakan dilaporkan.** Ronde `/grilling` yang sedang saya susun — frontier-nya arti `//`,
> nasib nilai lama, dan cakupan A2 — **digantikan** blok GANTI KONTEKS ini, yang menjawab ketiganya
> sebagai keputusan. Blok ini menang; ronde grilling itu dibatalkan. Ini **kali kedua berturut-turut**
> ronde grilling tergantikan brief; dicatat sebagai pola, bukan keluhan.

---

## §A — Empat keputusan work owner, dicatat

### A1 — ⭐ Label `//` MEMANG mematikan langkah `[data work owner 2026-09-19]`

`[data work owner 2026-09-19]` Work owner membuka `KomitePostAdjustment` di Pega dan memeriksa
step **16** anak **1–4**: keempatnya **tampak dicoret di layar**.

**Akibatnya, berurutan:**

| Akibat | Status |
| --- | --- |
| Aturan baca ronde 1–3 (*`//` = remark*) | ✅ **BERDIRI** |
| 68 langkah ber-`//` di Claim Prop | ✅ memang **mati** |
| Tiga `Obj-Save` + satu `Commit` di antaranya | ⭐ **TEMUAN: jalur simpan mati** — `SETPAYABLETO_ACT` 4 · `SETPAYABLETREATY_ACT` 4 · `GETPAYATTACHMENTADJ_ACT` 3.3.4 · `Commit` 3.3.9. ⛔ Dicatat sebagai **temuan**, **bukan cacat** |
| Pertanyaan 2 ronde 4 — `CHECKNOPOLICY` menembak pencarian polis tiga kali | ⛔ **GUGUR** — langkah 3 dan 4 mati, pencarian **tidak rangkap** |
| Tiket **01** | ⛔ **TIDAK ditambal** — tidak ada yang berubah untuknya |
| Butir goyah *"aturan `//` tidak berbukti"* | ✅ **TUTUP** |

⚠️ **Catatan bentuk bukti.** Bukti ini **bukan dari korpus** — ekspor Pega tidak memuat medan apa
pun yang mematikan langkah. Ia dari **layar Pega**, dan karena itu ditandai `[data work owner]`,
bukan `[terverifikasi]`. Dicatat begitu di `spec.md` §Aturan baca.

### A2 — Perkalian ganda share ceding TIDAK PERNAH BERJALAN `[keputusan work owner 2026-09-19]`

`[keputusan work owner 2026-09-19]` Blok yang memuat perkalian kedua atas share ceding di
`Activity/CountPersen_act.xml` **ber-remark**, jadi ia **tidak pernah dieksekusi di produksi**.

| | |
| --- | --- |
| ⭐ **Data lama** | **AMAN** — nol nilai tersimpan yang perlu diperbaiki |
| Butir `[terbuka]` *"nasib nilai lama"* | ✅ **TUTUP**, dengan alasan itu |
| Langkah itu | ⛔ **TIDAK dimigrasikan** |

### ⛔ A2-SELISIH — dilaporkan, tidak dikejar supaya cocok

Catatan lisan work owner menyebut: *"yang di-remark adalah **5.1**, BUKAN step 5. Step 5 labelnya
kosong."* **Ekspor korpus menunjukkan kebalikannya.**

`[terverifikasi]` Dibaca ulang dari `Activity/CountPersen_act.xml`, medan `pyStepsBlockName`:

| Langkah | `pyStepsBlockName` |
| --- | --- |
| 4.2.4 | `//` |
| **5** *(induk blok)* | ⭐ **`//`** |
| **5.1** | ⛔ **kosong** |
| 5.2 · 5.2.1 · 6 · 6.1 · 6.2 · 6.2.1 | kosong |

⛔ **Saya tidak memilih satu di antaranya.** Keduanya dikerjakan berdampingan di §B sebagai
**bacaan A** dan **bacaan B**. ⭐ **Kabar baiknya: kesimpulan hilirnya sama dalam kedua bacaan** —
lihat §B4. Yang berbeda hanya apakah ada cacat tambahan di dalam blok 5, dan itu digantung
`[terbuka]` di §B5.

### A3 — Lima titik perkalian kelas Adjustment: perkaliannya TETAP `[keputusan work owner 2026-09-19]`

`[keputusan work owner 2026-09-19]` `COUNTGROSSADJTREATY_ACT` langkah **3 · 4 · 5.2** dan
`COUNTVALUEADJTREATY_ACT` langkah **12** (dua titik).

- ⛔ **Perkaliannya JANGAN dibuang.** ⛔ **Rumusnya JANGAN diubah.**
- ✅ Yang diseragamkan hanya **ketelitiannya**: hitung dengan **ketelitian penuh**, **nol pemotongan
  di tengah jalan**.
- ⭐ Akibatnya angka berubah **hanya di belakang koma**, bukan sebesar satu faktor.

### A4 — Nama profil autentikasi dibersihkan dari ronde 2 ✅

✅ Dilakukan. **Satu baris**, tidak lebih. Buktinya di §E-lampiran.

---

## §B — `Local.Value` yang menggantung di `CountPersen_act`

### B1 — Setiap pengisi `Local.Value`: **dua**, dan keduanya di dalam blok

`[terverifikasi]` Sensus 100% atas seluruh langkah `Activity/CountPersen_act.xml` — `pyParamArray`,
`pyStepsCallParams`, dan `pyStepsJavaSource`:

| Peran | Langkah | Label | Rumus | Hidup? |
| --- | --- | --- | --- | --- |
| **isi** | **5.1** | kosong | `@toDecimal(.ClaimEstimation) * @divide(ShareCeding,100,4)` | ⚠️ bergantung bacaan |
| pakai | 5.2.1 | kosong | `.ClaimSpreaded := .SharePercentage * Local.Value / 100` | ⚠️ bergantung bacaan |
| **isi** | **6.1** | kosong | `@toDecimal(.Value)` | ✅ **HIDUP** *(gerbang `.Note!="Yes"`)* |
| pakai | 6.2.1 | kosong | `.ClaimSpreaded := @divide(.SharePercentage,100,10) * Local.Value` | ✅ **HIDUP** |

**Pengisi: 2. Yang pasti hidup: 1 (langkah 6.1).** ⛔ Nol pengisi di langkah 1–4, nol di langkah
`Java` (langkah 3), nol lewat parameter panggilan.

### B2 — Isi `Local.Value` saat blok 5 berjalan

⚠️ Pertanyaan ini **hanya ada dalam bacaan B**. Dijawab untuk keduanya:

#### Bacaan A — langkah **5 (induk)** ber-remark *(sesuai ekspor)*

Blok 5 **tidak berjalan sama sekali**: 5.1, 5.2, dan 5.2.1 ikut mati bersama induknya.
⛔ **Pertanyaan B2 gugur** — tidak ada saat di mana "blok 5 berjalan".

#### Bacaan B — hanya **5.1** yang ber-remark *(sesuai catatan lisan)*

`[terverifikasi]` Urutannya dibuktikan dari nomor langkah dan gerbangnya:

```
langkah 1..3   -> Local.Value TIDAK pernah disentuh
langkah 4.1    -> Local.Currency := .CurrencyID      (Local.Value tetap tidak disentuh)
langkah 4.4    -> .ClaimEstimation := Local.TotalList
langkah 5.1    -> MATI  (bacaan B)   <-- satu-satunya pengisi Local.Value di blok 5
langkah 5.2.1  -> .ClaimSpreaded := .SharePercentage * Local.Value / 100
```

⭐ **`Local.Value` belum pernah diisi apa pun ketika 5.2.1 membacanya.** `Local.*` adalah halaman
lokal activity, dibuat baru setiap pemanggilan — jadi ia **tidak** membawa sisa dari pemanggilan
sebelumnya, dan **tidak** ada langkah lain yang mengisinya lebih dulu.

**Akibatnya dalam bacaan B:** langkah 5.2.1 menulis `.ClaimSpreaded` dari nilai **kosong** —
hasilnya nol, atau gagal konversi. ⚠️ Lihat §B5.

⚠️ **Ikutan yang sama:** langkah 5.1 juga satu-satunya yang menyegarkan `Local.Currency` di blok 5.
Dalam bacaan B, gerbang 5.2.1 (`.CurrencyID==Local.Currency && .IsOldData!="Yes"`) dibandingkan
dengan mata uang **sisa dari blok 4**, yaitu baris terakhir `ListUang.pxResults` — bukan baris yang
sedang diulang.

### B3 — Blok 5 versus blok 6: gerbang, tabrakan, dan siapa menimpa siapa

`[terverifikasi]` Keduanya menulis **properti yang sama** (`.ClaimSpreaded`, `.ClaimEstimation`)
pada **halaman yang sama** (`pyWorkPage.ClaimData.SpreadingRisk`):

| | Blok 5 | Blok 6 |
| --- | --- | --- |
| Gerbang induk | `2/2` — **tanpa syarat** | `2/2` — **tanpa syarat** |
| Daftar yang diulang (luar) | `ListUang.pxResults` | `pyWorkPage.ClaimData.ListClaimAmount` |
| Langkah pengisi `Local.Value` | 5.1 — **tanpa syarat** | 6.1 — **`.Note!="Yes"`** |
| Daftar yang diulang (dalam) | `pyWorkPage.ClaimData.SpreadingRisk` | `pyWorkPage.ClaimData.SpreadingRisk` |
| Gerbang penulis | 5.2.1 — `.CurrencyID==Local.Currency && .IsOldData!="Yes"` | 6.2.1 — **syarat yang sama persis** |
| Urutan | lebih dulu | **belakangan** |

**Bisakah keduanya berjalan pada pemanggilan yang sama?**

- **Bacaan A:** ⛔ **tidak** — blok 5 mati. Hanya blok 6 yang berjalan.
- **Bacaan B:** ✅ **ya** — keduanya berjalan, **tanpa apa pun yang mencegahnya**. Gerbang induk
  keduanya `2/2`.

**Siapa menimpa siapa?** ⭐ **Blok 6 menimpa blok 5** — ia berjalan belakangan, menulis properti yang
sama, pada halaman yang sama, dengan **gerbang penulis yang identik**. Untuk setiap baris
`SpreadingRisk` yang lolos gerbang itu, nilai akhirnya **selalu milik blok 6**.

⚠️ `[terbuka]` **Satu celah tersisa dalam bacaan B.** Bila sebuah baris `ListClaimAmount` punya
`.Note=="Yes"`, langkah **6.1 dilewati** sehingga `Local.Value` **tidak disegarkan**, sedangkan
6.2.1 tetap berjalan dengan gerbangnya sendiri. Baris `SpreadingRisk` yang dihitung pada putaran itu
memakai `Local.Value` **basi**. ⛔ Ini berlaku di **kedua bacaan** — ia cacat blok 6, bukan blok 5.

### B4 — Nilai `.ClaimSpreaded` yang benar-benar dipakai hilir

`[terverifikasi]` Sesudah blok 6, langkah **7–11** tidak menyentuh `.ClaimSpreaded` maupun
`.ClaimEstimation` sama sekali:

| Langkah | Metode | Menyentuh nilai uang? |
| --- | --- | --- |
| 7 | `Page-Clear-Messages` | ⛔ tidak |
| 8 · 9 | `Property-Set` → `DataChronology.CARI1` | ⛔ tidak — teks kronologi |
| 10 | `Apply-DataTransform` | ⛔ tidak menulis ketiganya |
| 11 | `Call AddEstimation_Act` | ⛔ tidak menulis ketiganya |

> ⭐ **Jawabannya: nilai hilir SELALU berasal dari BLOK 6 — dalam kedua bacaan.**
> Blok 6 adalah blok yang **tidak** mengalikan share ceding, yaitu rumus yang **benar**.

Angka itu kemudian terbawa ke hilir lewat `pyWorkPage.ClaimData.SpreadingRisk`, dan satu-satunya
efek keluar yang membacanya adalah `Activity/SendEmailKlaim.xml` — badan surat pemberitahuan.

### B5 — ⚠️ Temuan yang dinaikkan ke work owner

> ⛔ **Tidak saya perbaiki sendiri.** Ditulis, ditandai `[terbuka]`, dinaikkan.

**Apa.** Dalam **bacaan B** (kalau yang ber-remark benar-benar 5.1 dan bukan langkah 5), langkah
**5.2.1** `Activity/CountPersen_act.xml` menulis `.ClaimSpreaded` memakai `Local.Value` yang
**belum pernah diisi oleh langkah mana pun** — satu-satunya pengisinya di blok itu adalah langkah
yang mati.

**Kenapa.** Kode lama dimatikan **sebagian**: perkalian share ceding-nya dicabut, tetapi pemakainya
di 5.2.1 dibiarkan hidup. Dalam bacaan A hal ini tidak terjadi, karena induknya ikut mati.

**Beda A/B.**

| | Bacaan A — langkah 5 mati *(ekspor)* | Bacaan B — hanya 5.1 mati *(catatan lisan)* |
| --- | --- | --- |
| Blok 5 berjalan? | ⛔ tidak | ✅ ya |
| 5.2.1 menulis nilai kosong? | ⛔ tidak terjadi | ⚠️ **ya** |
| Nilai akhir hilir | blok 6 | blok 6 *(menimpa)* |
| Cacat tersisa | **nol** | menulis-lalu-ditimpa; terlihat hanya bila blok 6 tidak menjangkau baris itu |

**Apa yang tertahan.** ⛔ **Tidak ada pekerjaan yang tertahan.** Keputusan A2 (*blok itu tidak
dimigrasikan*) berlaku sama dalam kedua bacaan, dan nilai hilirnya sama. Yang tertahan hanya
**ketepatan catatan**: `spec.md` butir `[terbuka]` ke-6 dan §A2-SELISIH di atas.

---

## §C — Ketelitian lima titik kelas Adjustment

### C1 — Kelima rumus APA ADANYA, berikut ketelitian yang berlaku sekarang

`[terverifikasi]` Disalin utuh dari `pyParamArray` kedua rule. **Seluruh lima titik memakai
ketelitian `20` pada setiap `@divide`** — seragam, tidak ada yang menyimpang.

#### `ASM-FW-GCNMFW-DATA-ADJUSTMENT!COUNTGROSSADJTREATY_ACT`

| Langkah | Sasaran | Rumus | Ketelitian tiap `@divide` |
| --- | --- | --- | --- |
| **3** | `.GrossValue` | `.GrossAdjustment × (PersenRNM/100) × (PersenLossAllocation/100) × (ShareCeding/100)` | `20` · `20` · `20` |
| **4** | `.GrossValue` | **rumus identik dengan langkah 3** | `20` · `20` · `20` |
| **5.2** | `.ClaimSpreaded` | `Local.TreatyGross × (ShareCeding/100) × (SharePercentage/100)` | `20` · `20` |

⚠️ Langkah **3** dan **4** berisi rumus **sama persis**; yang membedakan hanya gerbangnya —
`.Type==1` pada langkah 3, `.Type!=1` pada langkah 4. Dicatat karena ia mudah disangka rangkap.

#### `ASM-FW-GCNMFW-DATA-ADJUSTMENT!COUNTVALUEADJTREATY_ACT`

| Langkah | Sasaran | Rumus | Ketelitian tiap `@divide` |
| --- | --- | --- | --- |
| **12** | `.AdjustmentValue` | `.ProposeAdjustmentValue × (PersenRNM/100) × (PersenLossAllocation/100) × (ShareCeding/100)` | `20` · `20` · `20` |
| **12** | `.IndividualRiskRNM` | `.IndividualRiskValue × (PersenRNM/100) × (PersenLossAllocation/100) × (ShareCeding/100)` | `20` · `20` · `20` |

Di kedua rule, `local.ShareCeding` disalin lebih dulu dari `pyWorkPage.ClaimData.ShareCeding` —
langkah **1** dan langkah **10** masing-masing.

### C2 — Bentuk yang berlaku sesudah A3

Kelima titik itu **tetap mengalikan share ceding, tepat sekali, seperti sekarang**. Yang berubah
hanya cara angkanya dihitung: **seluruh perkalian dan pembagian dikerjakan dengan ketelitian penuh
dari nilai asalnya, dan tidak ada hasil antara yang dipotong ke sejumlah angka di belakang koma
sebelum dipakai langkah berikutnya**. Pemotongan hanya boleh terjadi sekali, pada saat angka itu
ditampilkan atau dicetak. Bentuk ini mengikuti keputusan uang yang sudah berdiri di `spec.md`
bab 4 — kali dulu, bagi terakhir, dan rantai *"hasil dipakai sebagai masukan"* diputus.

### C3 — Angka hilir yang ikut berubah, dan besarnya

| Angka hilir | Berubah? | Besarnya |
| --- | --- | --- |
| `.GrossValue` | ✅ ya | **hanya di belakang koma** |
| `.ClaimSpreaded` *(jalur Adjustment)* | ✅ ya | **hanya di belakang koma** |
| `.AdjustmentValue` | ✅ ya | **hanya di belakang koma** |
| `.IndividualRiskRNM` | ✅ ya | **hanya di belakang koma** |

⭐ **Tidak satu pun berubah sebesar satu faktor.** Ketelitian sekarang sudah `20` di seluruh lima
titik — yang paling dalam di modul ini — jadi selisihnya berada jauh di belakang koma dan hanya
muncul lewat pemotongan berantai, bukan lewat perubahan rumus.

**Silang dengan §B4 ronde 4 — jalur Kasir.** Ronde 4 menyatakan jalur Kasir *"tidak terkena"*, dan
saya sendiri menandainya **paling rawan** karena `.TotalClaim` belum dilacak ke hulu. §C4 di bawah
melacaknya, dan hasilnya **mengubah bentuk jawabannya**.

### C4 — ⭐ `.TotalClaim` dilacak ke hulu — hasilnya bukan yang saya duga

`[terverifikasi]` Sensus 100% atas **329 berkas**, seluruh wadah (`pyParamArray`,
`pyStepsCallParams`, `pyStepsJavaSource`, dan teks mentah untuk jenis rule non-Activity):

| | |
| --- | --- |
| Langkah yang **menulis** `.TotalClaim` | ⛔ **NOL** |
| Langkah yang **membaca** `.TotalClaim` | **2** |

Kedua pembacanya:

| Rule | Langkah | Pemakaian |
| --- | --- | --- |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT!HITSERVICETOKASIR_ACT` | **9.3** | `TempKasir.CARI9 := .TotalClaim − .PremiumSpreaded` |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT!SENDEMAILKLAIM` | **8.2.2** | dimasukkan ke muatan email |

> ⭐ **`.TotalClaim` diisi DI LUAR korpus Claim Prop.** Tidak ada satu pun dari 329 berkas yang
> menulisinya. Ia datang dari luar modul ini — basis data, modul lain, atau rule deklaratif yang
> tidak ikut diekspor.

**Apa artinya untuk kesimpulan ronde 4 yang saya tandai paling rawan:**

- ⛔ Kesimpulan itu **tidak runtuh**. Di dalam korpus, **nol jalur** menghubungkan kelima titik
  perkalian kelas Adjustment ke `.TotalClaim`. Kelimanya menulis `.GrossValue`, `.ClaimSpreaded`,
  `.AdjustmentValue`, `.IndividualRiskRNM` — **tidak satu pun menulis `.TotalClaim`**.
- ⚠️ Tetapi ia **juga tidak dapat ditutup**. Hulunya berada di luar jangkauan korpus, jadi
  *"tidak terkena"* hanya berlaku **sejauh modul ini terbaca**.
- ✅ **Bentuk jawaban yang berlaku:** `[terverifikasi]` di dalam Claim Prop, jalur Kasir **tidak
  terkena** perkalian share ceding; `[terbuka]` apa yang mengisi `.TotalClaim` — **di luar
  cakupan modul ini**, dan tidak ditebak.

⭐ **Ini menurunkan kerawanannya**, karena sebabnya berpindah: dari *"saya belum menyisir"* menjadi
*"tidak ada di korpus untuk disisir"*. Yang pertama adalah kelalaian; yang kedua adalah batas
cakupan yang jujur.

---

## §D — Tiga butir goyah sisa

### D1 — Ketiganya, dan apa persisnya yang kurang

> ⛔ **Tidak satu pun dinyatakan tertutup.** Penutupan adalah keputusan work owner.

#### Goyah 1 — *"FIX ERROR HANDLING"* tanpa jejak

| | |
| --- | --- |
| **Apa yang kurang** | Catatan pengembang pada `ASM-FW-GISFW-INT-T_STORAGE_IMAGE!INSERTGOOGLESTORAGE_ACT` dan `…!GETURLGOOGLESTORAGE_ACT` menyebut pekerjaan penanganan kegagalan yang **jejaknya tidak ditemukan** di keempat keluarga wadah |
| **Apa yang sudah cukup** | ✅ Ronde 4 membuktikan ia **satu rule bersama**, bukan dua temuan: `pxInsName` **dan** `pxUpdateDateTime` identik di **14 modul**. Jadi cakupannya pasti, dan membetulkannya sekali cukup untuk semua |
| **Apa yang masih perlu** | ⚠️ Apakah pekerjaan itu **pernah dikerjakan lalu dicabut**, atau **belum pernah dikerjakan sama sekali**. Korpus tidak dapat menjawabnya — ia hanya memuat keadaan terakhir, bukan riwayat. **Perlu work owner atau riwayat rule di Pega** |

#### Goyah 2 — *"89 memakai kode 5/6"* sebagai satu ember

| | |
| --- | --- |
| **Apa yang kurang** | Ronde 2 menggabung kode **5** (*Skip Whens*) dan kode **6** (*Exit Activity*) menjadi satu ember bernama *"belum diterangkan"*, padahal artinya berbeda jauh |
| **Apa yang sudah cukup** | ✅ Embernya sudah **dipecah dan dihitung ulang**: kode 5 = **7** baris, kode 6 = **16** baris, total **23** — bukan 89. Arti keenam kode arah juga sudah tercatat `[data work owner]` |
| **Apa yang masih perlu** | ⚠️ Hanya **penyelarasan angka**: ronde 2 masih memuat "89" dan "380 normal" di badan teksnya, sedangkan hitungan yang berlaku adalah **23** dan **363**. ⛔ Ronde 2 read-only, jadi selisih itu harus **diserap `spec.md`**, bukan ditambal di tempat asalnya. **Nol kesimpulan turunan** bersandar padanya |

#### Goyah 3 — bentuk rantai syarat tujuh langkah berkode 5

| | |
| --- | --- |
| **Apa yang kurang** | Gerbang **langkah induk** 7, 8, 9, 10 `CopyOldataCurr_act` **tidak pernah dituliskan**. Tanpa itu, rantai syarat anaknya terbaca sebagai AND dan menjadi mustahil — `CurrencyIDOld≠"" AND SpreadingRisk>0 AND (Interest OR ClaimAmount)` |
| **Apa yang sudah cukup** | ✅ Arti kode **5** sudah tegak: satu baris berkode 5 menjadikan rantainya **OR**, bukan AND. Jadi rantai itu **tidak mustahil** |
| **Apa yang masih perlu** | ⚠️ Gerbang keempat langkah induk itu **belum dibaca dan dituliskan**. Ia menyentuh **tiket 00 dan 06** serta **satu bab `spec.md`**, semuanya tentang **angka uang**. Pembaca yang menurunkan gerbangnya sendiri dengan aturan AND akan menyimpulkan langkahnya mati lalu **membuang rumusnya**. **Bisa dijawab dari korpus** — satu sisiran |

**Jumlah goyah sesudah ronde 5: 3.** *(Yang keempat — aturan `//` — ditutup §A1.)*

### D2 — Aturan parser baru, dan pola yang sudah menggigit **tiga kali**

> ### ⭐ ATURAN PARSER — *ketiadaan*
>
> **Ketiadaan tidak boleh disimpulkan sebelum SELURUH medan yang mungkin memuatnya disisir.**
>
> Sebelum menulis *"nol"*, *"tidak ada"*, atau *"tidak terbaca"*, wajib disebutkan **medan mana
> saja yang sudah disisir** dan **atas dasar apa daftar itu dianggap lengkap**. Kalau daftarnya
> tidak dapat dibuktikan lengkap, yang ditulis adalah **`[terbuka]`**, bukan *"nol"*.

**Ketiga kalinya, berurutan — polanya sama persis:**

| # | Kapan | Yang dinyatakan tidak ada | Medan yang belum disisir | Kenyataannya |
| --- | --- | --- | --- | --- |
| **1** | Ronde 2 — Komite Claim Prop | *"Sel tidak pernah bergerbang sendiri"* | `pyVisible` berada di dalam **sub-halaman** `pyUserData`, bukan anak langsung | Sel **memang** bergerbang sendiri — **108 gerbang hidup** |
| **2** | Ronde 5 — Komite Claim Prop | *"`M_LINK_SERVICE` membaca baris pertama **tanpa saringan**"* | saringannya ada di **`pyParamArray`**, wadah yang belum dibuka | Saringan **ada** — `.KATEGORI_1` + `.KATEGORI_2`. Pertanyaan saya sendiri ditarik |
| **3** | Ronde 3 → 4 — Claim Prop | *"Halaman yang diulang **tidak terbaca** dari struktur"* | `pyStepsObjectName` (Step Page) — medan **ketiga**, tidak ikut disisir | Terbaca untuk **194 dari 224** (`EMBEDDED` **192 dari 192**) |

**Bentuk yang berulang, dinyatakan sekali supaya kelihatan:** ketiganya menyisir **sebagian** wadah,
menemukan kosong, lalu menulis **kesimpulan positif tentang ketiadaan**. Yang salah bukan
sisirannya — yang salah adalah **melompat dari "saya tidak menemukan" ke "tidak ada"**.

⚠️ **Ada kemungkinan keempat yang berbentuk sama** — ronde 5 Komite, *"PDF nol penanganan gagal"*,
gugur setelah sumber Java dibaca. Tidak dimasukkan ke tabel karena brief meminta tiga; dicatat
supaya hitungannya jujur.

---

## §E — `spec.md` ditambal

`[terverifikasi]` **Lima sisipan, seluruhnya di bagian yang disebut brief.** Bagian lain utuh.

| Butir | Di mana | Isi |
| --- | --- | --- |
| **E1 · E3** | bab **4** *(Uang — satu jalur perhitungan)* | keempat keputusan §A, bertanggal, berikut A2-SELISIH sebagai `[terbuka]` |
| **E2** | *Further Notes → Aturan baca korpus yang berlaku* | `//` disahkan sebagai aturan baca, **berikut bentuk buktinya**: layar Pega, `[data work owner]`, **bukan** korpus |
| **E4** | bab **14b** *(baru, sebelum bab 14)* | penggolongan *"perubahan sadar ketiga"* **DICABUT** |
| **E5** | *Pertanyaan terbuka → Status sekarang* | daftar `[terbuka]` disusun ulang + butir ringan **ke-6** ditambahkan |

### E4 — Titik yang sengaja diubah: **tetap DUA**, bukan tiga

Ronde 4 mencatat *"rumus share ceding diperbaiki"* sebagai **perubahan sadar ketiga**. Ronde 5
**mencabutnya**, dengan alasan satu baris:

> Perkalian gandanya **tidak pernah berjalan** (blok ber-remark), jadi **tidak ada perilaku berjalan
> yang disimpangi**; dan yang benar-benar berubah menurut A3 hanyalah **ketelitian**, yang sudah
> tercakup keputusan uang di bab 4 butir 1–3.

⛔ **Jumlah perubahan sadar: tetap DUA.**

### E5 — Daftar `[terbuka]` sesudah disusun ulang

| | |
| --- | --- |
| **HILANG** | ⛔ **NOL** dari tabel butir ringan. Dua butir yang digantung ronde 4 ditutup §A dan **tidak pernah masuk tabel**: *nasib nilai lama* (A2) dan *arti `//`* (A1) |
| **BERTAMBAH** | **satu** — butir ke-6: label `//` di langkah **5** versus **5.1** `CountPersen_act` |
| **JUMLAH AKHIR** | ⭐ **ENAM** butir `[terbuka]` ringan · **`[terbuka]` yang memblokir: NIHIL** |

⚠️ Blok ringkasan lama di kepala `spec.md` masih menulis *"empat"*; RALAT 2026-09-19 yang saya
sisipkan menyatakan angka yang berlaku adalah **enam** dan menerangkan selisihnya. ⛔ Angka lama
**tidak dihapus** — ia milik RALAT 2026-09-18 (b) dan dibiarkan apa adanya.

---

## §F — Apa lagi

### F1 — Sesudah ronde 5, urut dari yang paling menahan

| # | Yang dikerjakan | Kenapa | Besar | Menunggu |
| --- | --- | --- | --- | --- |
| **1** | **Gerbang induk `CopyOldataCurr_act` 7 · 8 · 9 · 10** | ⭐ satu-satunya goyah sisa yang **bisa dijawab dari korpus**; menyentuh tiket **00** dan **06** serta satu bab spec, semuanya angka uang | satu sisiran | korpus |
| **2** | **87 catatan `pyMemo` yang belum diuji** | dari 4 larangan langsung yang diuji, **2 masih dikerjakan** — separuh. Sisanya belum dinilai sama sekali | sedang | korpus |
| **3** | **Selaraskan angka ronde 2 ke `spec.md`** | ronde 2 masih memuat "89" dan "380"; yang berlaku **23** dan **363**. Ronde 2 read-only, jadi penyelarasannya harus di spec | sekali jalan | — |
| **4** | **11 `RDBList` + 11 `Section` + sisa 10 Activity** yang belum dinilai | penyaringan otomatis tidak dapat menilainya; perlu dilihat satu per satu | sedang | korpus |
| **5** | **Riwayat rule `INSERTGOOGLESTORAGE_ACT` di Pega** | satu-satunya cara menutup goyah 1 | kecil | **work owner** |

### F2 — Kesimpulan yang sekarang PALING RAWAN SALAH

⛔ **Ditunjuk, tidak diperbaiki.**

> ⚠️ **Yang paling rawan: §B3 — *"blok 6 selalu menimpa blok 5"*.**

Ia bersandar pada satu penalaran urutan: blok 6 bernomor lebih besar, menulis properti yang sama,
pada halaman yang sama, dengan gerbang penulis yang identik — **jadi** ia menimpa.

**Yang belum saya buktikan:** bahwa kedua blok benar-benar menyentuh **baris `SpreadingRisk` yang
sama**. Daftar luarnya **berbeda** — blok 5 mengulang `ListUang.pxResults`, blok 6 mengulang
`ListClaimAmount`. Kalau kedua daftar luar itu menghasilkan **himpunan mata uang yang berbeda**,
maka gerbang `.CurrencyID==Local.Currency` bisa meloloskan **baris yang berbeda** di tiap blok, dan
ada baris `SpreadingRisk` yang **hanya** disentuh blok 5 — yang berarti nilai blok 5 **bertahan**.

⚠️ **Ini bentuk yang sama dengan §D2 nomor 3**: saya menyimpulkan dari **satu** sumbu (nomor
langkah) tanpa menyisir sumbu kedua (himpunan baris yang benar-benar terjangkau). ⛔ Dalam
**bacaan A** hal ini tidak berakibat apa-apa; dalam **bacaan B** ia menentukan apakah nilai kosong
pernah bertahan sampai hilir.

Yang paling rawan kedua: **§C3 *"hanya di belakang koma"***. Ia benar bila ketelitian `20` memang
yang dipakai di seluruh rantai. Saya membuktikan `20` **di kelima titik itu**, tetapi **tidak** di
langkah-langkah yang memasok `.GrossAdjustment`, `.PersenRNM`, dan `.PersenLossAllocation`.

### F3 — Pertanyaan BARU untuk work owner — **1**

> ⛔ Tidak mengulang apa pun yang sudah dijawab §A.

#### Pertanyaan 1 — `CountPersen_act`: yang ber-remark langkah 5 atau 5.1?

**Apa yang ditanyakan.** Catatan lisan menyebut yang dicoret adalah **5.1**, dan bahwa langkah 5
berlabel kosong. Ekspor korpus menunjukkan kebalikannya: label `//` ada di **langkah 5**, dan 5.1
kosong. Mohon dilihat sekali lagi di layar Pega — apakah yang tercoret **satu baris** (5.1) atau
**seluruh blok** beserta anak-anaknya.

**Kenapa ini penting.** Kalau yang mati hanya 5.1, maka langkah **5.2.1** tetap berjalan dan menulis
nilai uang dari `Local.Value` yang **tidak pernah diisi apa pun** — nol atau gagal konversi. Kalau
induknya yang mati, hal itu tidak pernah terjadi dan blok 5 bersih sepenuhnya.

**Pilihan yang saya lihat.** (a) langkah 5 mati → blok 5 bersih, nol cacat tersisa; (b) hanya 5.1
mati → ada penulisan nilai kosong yang kemudian ditimpa blok 6, dan perlu dipastikan apakah ia
selalu tertimpa.

**Yang saya sarankan.** ⚠️ Saya **tidak menyarankan** — ini pertanyaan tentang apa yang terlihat di
layar, bukan tentang apa yang sebaiknya dilakukan. ⛔ Yang penting dicatat: **keputusan A2 tidak
berubah dalam kedua jawaban**, dan **angka hilirnya sama**. Pertanyaan ini menentukan **ketepatan
catatan**, bukan pekerjaan.

---

## §E-lampiran — bukti berkas yang tidak bergerak

Sidik jari MD5 diambil **sebelum** penyuntingan dan dibandingkan **sesudahnya**.

### ✅ BERUBAH — tepat **2** berkas, keduanya disebut brief

```
spec.md                 <-- §E saja (lima sisipan)
grilling-ronde-2.md     <-- §A4 saja (satu baris)
```

### ⛔ TIDAK BERGERAK — MD5 identik sebelum dan sesudah

```
grilling-ronde-1.md · grilling-ronde-3.md · grilling-ronde-4.md
periksa-ulang-aturan-baru.md
issues/00 · 01 · 02 · 03 · 04 · 05 · 06 · 07 · 08 · 09 · 10 · 11 · 12 · 13 · 14 · 15   (16 tiket)
```

**⭐ Tiket: NOL disentuh — keenam belasnya terbukti identik.**
**4 berkas ronde/periksa-ulang terbukti identik.** Total **20 berkas** tidak bergerak.

### Bukti §A4 — `grilling-ronde-2.md` hanya satu baris berubah

| | Sebelum | Sesudah |
| --- | --- | --- |
| Jumlah baris | 3 391 | 3 392 *(+1 — kalimat dibungkus dua baris agar muat lebar)* |
| Nama profil autentikasi | ada, 1 kemunculan | ⭐ **NOL** |
| `InsertLOGDirectKasir_SQL` *(nama **rule**, bukan profil)* | 1 | **1 — tidak disentuh** |

Teks penggantinya: *"profil autentikasi bernama — namanya sengaja tidak dicatat"*.
⛔ Tidak ada perubahan lain di berkas itu.

### Berkas BARU yang ditulis ronde ini — **satu**

```
grilling-ronde-5.md
```
