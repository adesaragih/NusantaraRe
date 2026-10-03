# Grilling Ronde 4 — Claim Prop

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Claim Prop\` (READ-ONLY)
**Sasaran:** dua keputusan work owner (§A), sensus 100% share ceding (§B), letak kredensial Kasir (§C),
gabungan kesimpulan goyah (§D).

> ⛔ **Nilai rahasia NOL.** Tidak satu pun nilai kredensial disalin ke berkas ini. Perintah yang
> dipakai untuk §C ditulis ulang supaya secara struktural **tidak bisa** mencetak isi medan —
> ia hanya melaporkan nama medan dan status *TERISI/kosong*.

> ⛔ Berkas ronde 1, ronde 2, `periksa-ulang-aturan-baru.md`, `grilling-ronde-3.md`, `spec.md`,
> dan 12 tiket lain **TIDAK disunting**. Yang disunting hanya tiket **00 · 04 · 11 · 13**,
> hanya pada bagian yang disebut — lihat §E.

> ⚠️ **Tabrakan dilaporkan.** Ronde `/grilling` yang sedang saya susun (dua pertanyaan: kredensial
> Kasir, dan share ceding) **digantikan** blok GANTI KONTEKS ini, yang menjawab keduanya sebagai
> keputusan. Blok ini menang; ronde grilling itu dibatalkan.

---

## §A — Dua keputusan work owner, dicatat

### A1 — Kredensial autentikasi Kasir ada di dalam rule `[keputusan work owner]`

`[keputusan work owner]` Kredensial autentikasi ke Kasir **tertulis di dalam rule**, sehingga ada
rahasia yang ikut beredar di dalam berkas ekspor korpus.

**Aturan yang berlaku sejak sekarang, dicatat supaya tidak terlanggar di ronde berikutnya:**

1. ⛔ **Nilainya tidak disalin ke berkas mana pun** — tidak ke ronde ini, tidak ke tiket, tidak ke
   spec, tidak ke laporan. Yang boleh disebut hanya: **rule apa**, **di bagian apa**, dan
   **bentuknya apa**.
2. ⛔ **Tidak ditampilkan di keluaran perintah.** Perintah yang akan mencetaknya diganti.

⚠️ **Peringatan satu baris — dicatat, tidak ditindaklanjuti sendiri:** karena kredensial itu berada
di dalam ekspor yang beredar, ia **sebaiknya diganti sesudah migrasi**. Itu urusan **tim Kasir**;
dicatat di sini semata supaya tidak terlewat.

### A2 — Rumus share ceding DIPERBAIKI, bukan disalin `[keputusan work owner]`

`[keputusan work owner]` Perkalian ganda share ceding **tidak ditiru** di aplikasi Go. Rumusnya
**diperbaiki**. Ini menjadi **perubahan sadar yang ketiga** di modul Claim Prop.

⚠️ `[terbuka]` **Nasib nilai lama belum diputuskan** — angka yang sudah terlanjur tersimpan dengan
perkalian ganda. Menunggu work owner; **tidak ditebak di sini**.

---

## §B — Share ceding: sensus 100%

### B1 — Setiap titik kali/bagi, tanpa kecuali

`[terverifikasi]` Penyaringan dilakukan **tidak peka huruf besar-kecil** atas seluruh **329 berkas**.
Properti yang dipakai: `.ClaimData.ShareCeding` dan salinan lokalnya `local.ShareCeding`.

**Seluruh korpus memuat share ceding di 7 berkas saja** — 5 Activity dan 2 Section.

| # | Rule (`pxInsName`) | Langkah | Perlakuan |
| --- | --- | --- | --- |
| 1 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY!ADDLISTCLAIMAMOUNT` | **1** | **mengisi**: dipaku `"100"` bila kosong/nol *(bergerbang)* |
| 2 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY!ADDLISTCLAIMAMOUNT` | **7.1** | ✖ kali |
| 3 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY!COUNTLISTCLAIMAMOUNTIDR` | **1.1** | ✖ kali |
| 4 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY!COUNTPERSEN_ACT` | **5.1** | ✖ kali ⭐ **yang dilarang catatan pengembang** |
| 5 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT!COUNTGROSSADJTREATY_ACT` | **1** | salin ke `local.ShareCeding` |
| 6 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT!COUNTGROSSADJTREATY_ACT` | **3** · **4** · **5.2** | ✖ kali (3×) |
| 7 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT!COUNTVALUEADJTREATY_ACT` | **10** | salin ke `local.ShareCeding` |
| 8 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT!COUNTVALUEADJTREATY_ACT` | **12** | ✖ kali (2×) |
| 9 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY!INPUTACCEPTATION_EST` | — | Section: medan masukan |
| 10 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY!OUTSTANDINGCLAIM_EST` | — | Section: medan masukan |

**Hitungan sendiri:** **8 titik perkalian**, **0 titik pembagian oleh share ceding**, 2 penyalinan
ke variabel lokal, 1 pengisian nilai, 2 medan masukan. Seluruh perkalian berbentuk
`@divide(ShareCeding,100,N) * …` — persen diubah ke pecahan lebih dulu, presisi `N` **tidak seragam**:
**10** (dua titik), **20** (lima titik), **4** (satu titik — justru titik 5.1).

### B2 — Di mana ia "sudah dikalikan" — hulunya, langkah demi langkah

Catatan pengembang pada `COUNTPERSEN_ACT` berbunyi *"perbaiki struktur hitungan dan **jgn kali share
ceding lagi karna udah dikalikan**"*. Hulunya terlacak penuh:

```
ADDLISTCLAIMAMOUNT langkah 7.1      (mengulang pyWorkPage.ClaimData.ListClaimAmount)
    .Value  :=  ShareCeding/100  x  .NetDeductibleValue      <-- SHARE CEDING MASUK DI SINI
            atau
COUNTLISTCLAIMAMOUNTIDR langkah 1.1
    .Value  := (ShareCeding/100  x  .ClaimAmount) - .NetDeductibleValue   <-- ATAU DI SINI
                     |
                     v
COUNTPERSEN_ACT langkah 4.3         (mengulang daftar yang SAMA: ListClaimAmount)
  langkah 4.3.1   Local.TotalList := Local.TotalList + .Value    [syarat .CurrencyID==Local.Currency]
  langkah 4.4     .ClaimEstimation := Local.TotalList            <-- sudah membawa share ceding
                     |
                     v
COUNTPERSEN_ACT langkah 5.1
    Local.Value := .ClaimEstimation  x  ShareCeding/100          <-- DIKALIKAN KEDUA KALINYA
```

`[terverifikasi]` Daftar yang diulang langkah **4.3** dan yang ditulisi langkah **7.1** adalah daftar
yang sama: `pyWorkPage.ClaimData.ListClaimAmount`. **Rantainya tertutup.**

### B3 — Rumus yang benar dan rumus yang berjalan, berdampingan

| | Rumus |
| --- | --- |
| **Yang berjalan** *(langkah 5.1 → 5.2.1)* | `ClaimSpreaded = SharePercentage × ( ClaimEstimation × ShareCeding/100 ) / 100` |
| **Yang benar** *(blok kembarnya, langkah 6.1 → 6.2.1)* | `ClaimSpreaded = (SharePercentage/100) × Value` — **tanpa** share ceding |

⭐ **Blok 6 adalah kembaran blok 5 yang sudah diperbaiki.** Keduanya berada di activity yang sama,
mengulang daftar yang sama pada langkah `5.2` dan `6.2` (`pyWorkPage.ClaimData.SpreadingRisk`), dan
keduanya menulis properti yang sama (`.ClaimSpreaded`, `.ClaimEstimation`). Blok 6 berjalan
**sesudah** blok 5.

**Faktor dan arah:**

- Faktor kesalahan = **ShareCeding / 100**.
- Arah = **terlalu kecil** (understated) selama ShareCeding < 100.
- ShareCeding = 100 → faktor **1** → **tidak terlihat sama sekali**.
- ShareCeding = 80 → hasil **80%** dari yang benar, yaitu **kurang 20%**.

⚠️ `[terverifikasi]` `ADDLISTCLAIMAMOUNT` langkah **1** memaku `ShareCeding := "100"` ketika nilainya
kosong atau nol. **Pada jalur nilai bawaan, cacat ini tidak menimbulkan selisih.** Ia hanya menggigit
ketika share ceding benar-benar diisi bukan 100.

### ⛔ B3-RALAT — *"langkah ber-remark"* tidak dapat dibuktikan dari ekspor

Ronde-ronde sebelumnya memakai aturan **`pyStepsBlockName` diawali `//` = langkah diremark**.
Ronde 4 mengujinya, dan **aturan itu tidak berdiri di atas bukti**:

| Yang diuji | Hasil |
| --- | --- |
| Tag khusus untuk mematikan langkah di seluruh korpus | ⛔ **tidak ada** — yang ada hanya `pyDisabled*` milik elemen UI, bukan langkah activity |
| Langkah berlabel diawali `//` | **68** dari **1366** |
| Di antaranya **tanpa syarat apa pun** pada keluarga gerbang pertama | **32** |
| Gerbang `COUNTPERSEN_ACT` langkah **5** | `2/2` — **tanpa syarat, jadi selalu jalan** |

`[terbuka]` **Apakah `//` benar-benar mematikan langkah di mesin Pega tidak dapat dijawab dari
korpus.** Yang ada di ekspor hanyalah **label**. Kalau ia bukan saklar, ke-32 langkah itu **berjalan**
— dan di antaranya ada **`Obj-Save`** (`SETPAYABLETO_ACT` 4, `SETPAYABLETREATY_ACT` 4,
`GETPAYATTACHMENTADJ_ACT` 3.3.4) dan **`Commit`** (`GETPAYATTACHMENTADJ_ACT` 3.3.9).
**Menunggu work owner.**

⚠️ Ini **tidak mengubah keputusan A2** — rumusnya diperbaiki dalam kedua bacaan. Yang bergantung
padanya adalah **nasib nilai lama** (`[terbuka]` di §A2): kalau `//` bukan saklar, perkalian ganda
benar-benar pernah berjalan pada data produksi.

### B4 — Angka hilir yang terpengaruh

| Hilir | Terpengaruh? | Bukti |
| --- | --- | --- |
| **Email komite** — `ASM-FW-GCNMFW-DATA-ADJUSTMENT!SENDEMAILKLAIM` | ⚠️ **YA** | memuat `.ClaimSpreaded` ke badan surat |
| **Kasir** — `ASM-FW-GCNMFW-DATA-ADJUSTMENT!HITSERVICETOKASIR_ACT` langkah 9.3 | ⛔ **TIDAK langsung** | medan uangnya `TempKasir.CARI9 := .TotalClaim − .PremiumSpreaded` — bukan `.ClaimSpreaded`/`.ClaimEstimation` |
| **Dokumen / cetakan** | ⛔ tidak ditemukan | ketiga properti tidak muncul di `ReportDefinition/` maupun `RDBList/` |
| **Sebaran keseluruhan** | 13 Activity + 5 Section | penyaringan 100% atas 329 berkas |

⚠️ `[terbuka]` **Jalur ke Kasir melewati keluarga rumus yang berbeda** — `.AdjustmentValue`,
`.GrossValue`, `.IndividualRiskRNM` pada `COUNTGROSSADJTREATY_ACT` dan `COUNTVALUEADJTREATY_ACT`.
Kelima titik perkalian di kedua rule itu **masing-masing hanya sekali mengalikan share ceding**, dan
**tidak ada catatan pengembang yang melarangnya**. Apakah kelimanya juga terkena pola yang sama
**belum dapat dibuktikan** — tidak ditebak di sini.

### B5 — Bab spec dan tiket yang terpengaruh

| Sasaran | Bagian | Sifatnya |
| --- | --- | --- |
| **`spec.md`** Implementation Decisions — rumus uang | perkalian ganda share ceding | ⭐ perubahan sadar **ketiga** — ⛔ **belum ditambal**, spec tidak disunting ronde ini |
| **tiket `00`** prefactor skema transaksi uang | AC 19–20, pembulatan dan satu fungsi spreading | ✅ **ditambal** §E |
| **tiket `13`** efek keluar | autentikasi Kasir | ✅ **ditambal** §E |
| **tiket `11`** penyerahan komite | pengiriman Kasir | ✅ **ditambal** §E |
| **tiket `04`** klasifikasi lini bisnis | empat rule `When` bercatatan | ✅ **ditambal** §E — hasilnya **negatif**, lihat §E |

### B6 — Catatan pengembang KEDUA yang melarang sesuatu yang masih dikerjakan

`[terverifikasi]` Sensus 100% `pyMemo`: **61 catatan** berbentuk larangan atau perintah perbaikan.
Yang berbentuk **larangan langsung** ada **empat**; keempatnya diuji:

| Rule | Catatan | Masih dikerjakan? |
| --- | --- | --- |
| `COUNTPERSEN_ACT` | *"jgn kali share ceding lagi karna udah dikalikan"* | ⚠️ **YA** — langkah 5.1, gerbang `2/2` |
| **`CHECKNOPOLICY`** | *"jangan pake yg di query"* | ⚠️ **YA** — langkah **3** (`Call Rule-Obj-Report-Definition`) dan **4** (`RDB-List`) berlabel `//`, **tanpa syarat apa pun**; langkah 7 menjalankan `RDB-List` pengganti ⭐ **TEMUAN BARU** |
| `INSERTJSONCLAIMTREATY_ACT` | *"hapus yg hapus attachment list"* | ⛔ **TIDAK** — nol langkah `//`, nol langkah penghapus lampiran |
| `SETCATASTROPE_ACT` · `SETDEFNONCATASTROPE_ACT` | *"hapus save"* | ⛔ **TIDAK** — sudah diuji ronde 3, nol langkah simpan |

⭐ **`CHECKNOPOLICY` adalah kembaran struktural `COUNTPERSEN_ACT`**: catatan melarang, kode lama
dibiarkan dengan label `//`, kode baru ditaruh di bawahnya. Ini **pola**, bukan kejadian tunggal —
dan ia membuat B3-RALAT di atas jauh lebih menentukan.

---

## §C — Kredensial Kasir: letak, bentuk, dan apa yang TIDAK ada di korpus

### C1 — Rule yang memakainya

`[terverifikasi]` **`ASM-FW-GCNMFW-DATA-ADJUSTMENT!SENDACCEPTATIONTOKASIR`**
(`Claim Prop\ConnectREST\SendAcceptationToKasir.xml`). Catatan pengembangnya: *"tambah auth"*.

### C2 — Bentuk autentikasinya, dan medan mana yang membawanya

⛔ **Nilai tidak dicetak, tidak disalin, dan tidak dibaca.**

| Bagian rule | Medan | Bentuk |
| --- | --- | --- |
| pengaturan Connect-REST | `pyUseAuthentication` | **`true`** — satu-satunya di seluruh 329 berkas Claim Prop |
| pengaturan Connect-REST | `pyAuthProfileSelectionType` | **profil autentikasi bernama** (bukan kredensial sebaris) |
| pengaturan Connect-REST | `pyAuthenticationProfile` | berisi **nama profil** — ⛔ namanya pun tidak ditulis di sini |
| pengaturan proksi | `pyProxyAuthTypeSelection` | `NO_AUTH` — proksi tidak berautentikasi |
| blok tertanam | `pxObjClass = Data-Admin-Security-AuthenticationProfile` | **rujukan** ke rule profil |

**Bentuknya: autentikasi lewat profil bernama, bukan header yang ditulis di badan rule.**

### C3 — Rule lain dengan pola yang sama

`[terverifikasi]` Sensus **seluruh 22 modul korpus** (bukan hanya Claim Prop):

| Pola | Jumlah | Nama rule |
| --- | --- | --- |
| `pyUseAuthentication = true` | **6 berkas**, **1 identitas rule** | `ASM-FW-GCNMFW-DATA-ADJUSTMENT!SENDACCEPTATIONTOKASIR` — salinannya di `Claim Prop`, `Claim Fac In`, `Claim Non Prop`, `Komite Claim FacIn`, `Komite Claim Non Prop`, `Komite Claim Prop` |
| Connect-REST lain di Claim Prop | 5 | `GETDTLPAYMENTCLAIM` · `KONVERSIKLAIMNONLIFE` · `SERVICEGOOGLE` · `GETPAYATTACHMENT` · `GETPREMIUMPAIDONTREATYIN` — semuanya `pyUseAuthentication = **false**` |
| Token diambil dari basis data (pola berbeda, **bukan** kredensial tertanam) | 1 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE!RNM!GETTOKENSTORAGE_SQL` |

### ⚠️ C-RALAT — dua angka saya sendiri, dihitung ulang dan **berbeda**

1. ⛔ **Ronde 3 §C2 menulis *"15 kemunculan auth dan satu `Basic`"*.** Dihitung ulang: kata `Basic`
   itu adalah nilai medan **`pyDisplayMode = basic`**, yaitu medan tampilan UI yang muncul di
   **puluhan berkas biasa** dan **tidak ada hubungannya dengan autentikasi**. **Bukan bukti
   autentikasi dasar.** Kalimat ronde 3 itu **ditarik**.
2. ⚠️ **Nilai kredensialnya sendiri tidak saya temukan di korpus.** Yang diperiksa, seluruh 22 modul:
   - medan `pyUserName` / `pyPassword` / `pyUserID` / `pyAuthType` / `pyClientSecret` terisi → **0 berkas**
   - literal `Basic <base64>` → **tidak ada** · `Bearer <token>` → **tidak ada**
   - pola `user:pass@` di dalam URL → **tidak ada**
   - rule profil autentikasi sebagai berkas tersendiri (`pxInsName` bertipe profil) → **tidak ada**

   `[terbuka]` **Premis §A1 dan isi korpus berbeda di titik ini.** Di ekspor yang saya baca, rule
   Kasir hanya **merujuk** profil autentikasi bernama; **rule profilnya sendiri tidak ikut
   diekspor**. Mungkin nilainya ada di ekspor lain yang tidak saya pegang, atau tersimpan di dalam
   aliran terkode `pyRuleData` yang **sengaja tidak saya buka**. Dilaporkan apa adanya, **tidak
   dikejar supaya cocok**. ⛔ Ini **tidak mengubah** aturan §A1 — aturan itu tetap saya patuhi penuh.

### C4 — Apa yang KURANG dari tiket 13

| Kurang | Akibatnya bila tidak ditambahkan |
| --- | --- |
| Tiket 13 **sama sekali tidak menyebut autentikasi** | orang yang membangun klien Kasir dari tiket ini akan membuat panggilan **tanpa autentikasi**, lalu ditolak Kasir |
| Tidak menyebut bahwa autentikasi berbentuk **profil bernama** | kredensialnya akan dipaku di kode, bukan di konfigurasi |
| Tidak menyebut bahwa **hanya panggilan Kasir** yang berautentikasi | kelima Connect-REST lain akan ikut diberi autentikasi yang tidak perlu |
| Tidak mencatat bahwa kredensial **sebaiknya diganti sesudah migrasi** | peringatan itu hilang bersama ronde ini |

✅ Keempatnya ditambahkan di §E.

---

## §D — Kesimpulan goyah, digabung dan diperiksa

### D1 — Tujuh butir mentah, digabung dan dibuang rangkapnya

| # | Kesimpulan goyah | Asal | Status sesudah ronde 4 |
| --- | --- | --- | --- |
| **1** | Bab efek keluar tidak menyebut autentikasi Kasir | ronde 3 G1-1 | ✅ **bukti lengkap** — §C, tiket 13 ditambal. ⛔ tetapi bukti pendukungnya (*"satu `Basic`"*) **ditarik**, lihat C-RALAT |
| **2** | `CountPersen` 5.1 mengalikan share ceding padahal dilarang | ronde 3 G1-2 | ✅ **bukti lengkap** — §B2/B3 melacak rantainya penuh; keputusan A2 menjawabnya |
| **3** | *"84 berkas belum dibuka"* → angkanya **80** | ronde 3 G1-3 | ✅ **bukti lengkap** — sudah diralat di ronde 3 |
| **4** | *"FIX ERROR HANDLING"* tanpa jejak | ronde 3 G1-4 | ⚠️ **MASIH GOYAH** — lihat D2 |
| **5** | *"89 memakai kode 5/6"* menggabung dua hal | periksa-ulang GOYAH 1 | ⚠️ **MASIH GOYAH** — angka sensus; nol kesimpulan turunan bersandar padanya |
| **6** | Bentuk rantai syarat tujuh langkah berkode 5 tak dinyatakan | periksa-ulang GOYAH 2 | ⚠️ **MASIH GOYAH** — menyentuh tiket 00 dan 06 |
| **7** | Jalur kegagalan `SaveAdjusterConsultant_Act` | periksa-ulang GOYAH 3 | ✅ **bukti lengkap** — tiket 08 ditambal 2026-09-19 |

**Rangkap yang dibuang:** butir 1 dan ronde-3 G1-1 adalah butir yang sama; butir 7 dan ronde-3 §A
adalah butir yang sama. **7 mentah → 7 unik → 3 masih goyah**, ditambah **satu goyah BARU**:

| **BARU** | Aturan *"`//` = langkah diremark"* dipakai sejak ronde 1 tanpa bukti | ronde 4 §B3-RALAT | ⛔ **GOYAH BERAT** — menyentuh **68 langkah**, termasuk `Obj-Save` dan `Commit` |

**Jadi: 4 goyah sesudah ronde 4.**

⛔ **Catatan bentuk:** *"bukti lengkap"* **bukan** *"butir ditutup"*. Menutup butir adalah keputusan
work owner.

### D2 — `INSERTGOOGLESTORAGE_ACT`: identitas rule dicocokkan

`[terverifikasi]` Dicocokkan dengan **dua** penanda sekaligus, bukan nama berkas:

| Penanda | Salinan Komite Claim Prop | Salinan Claim Prop | Cocok? |
| --- | --- | --- | --- |
| `pxInsName` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE!INSERTGOOGLESTORAGE_ACT` | **sama persis** | ✅ |
| `pxUpdateDateTime` | `2026-03-13 07:47:20.538 GMT` | **sama persis** | ✅ |
| `pxCreateDateTime` | — | `2026-03-13 07:46:56.266 GMT` | — |
| `pyMemo` | `FIX ERROR HANDLING` | **sama persis** | ✅ |

⭐ **Rule yang sama, satu identitas, tersebar di 14 modul** — `Claim Fac In` · `Claim Life` ·
`Claim Non Prop` · `Claim Prop` · `Endorsment Fac In` · `Komite Claim FacIn` · `Komite Claim Life` ·
`Komite Claim Non Prop` · `Komite Claim Prop` · `Master Product Name Life` · `NB FacIn` ·
`RNW Fac In` · `Treaty In` · `Treaty In Adjustment`. Ketiga penanda **identik di keempat belasnya**.

**Kesimpulannya: temuan *"FIX ERROR HANDLING tanpa jejak"* bukan dua temuan di dua modul — ia SATU
temuan pada SATU rule bersama.** Membetulkannya di satu tempat membetulkannya untuk semua.

⚠️ `[terbuka]` Ukuran berkasnya **tidak** identik: **159 896 byte** di delapan modul, **160 027 byte**
di enam modul — selisih **131 byte** dengan `pxUpdateDateTime` yang sama. Selisih itu ada pada
**pembungkus ekspor**, bukan pada rule-nya. Tidak dikejar lebih jauh.

### D3 — `SETNAMECURRENCY_ACT`: kapan kurs terkunci

`[terverifikasi]` Kelas `ASM-FW-GCNMFW-DATA-ADJUSTMENT` · 18 langkah. Catatan pengembangnya: *"ok"*.

**Kuncinya ada di langkah 6, dan bentuknya "isi hanya bila kosong":**

```
langkah 6:   .KursIDR  :=  @if( .KursIDR == "" , <kurs dari pencarian> , .KursIDR )
```

⭐ **Artinya kurs dikunci pada pengisian pertama.** Sekali `.KursIDR` terisi, pencarian kurs
**tidak lagi menimpanya** — nilai lama dipertahankan, berapa pun kurs hari ini.

Jalur yang dilewati sebelum dan sesudahnya:

| Langkah | Yang terjadi |
| --- | --- |
| **1** | menyiapkan pencarian kurs untuk mata uang dari parameter |
| **2** | ⚠️ berlabel `//`, **tanpa syarat** — mengisi `Primary.KursIDR := .KursValue` dari daftar penyesuaian bila mata uangnya cocok. Lihat B3-RALAT |
| **3 · 5** | mengambil kurs (laporan + `RDB-List`) |
| **6** | ⭐ **penguncian** |
| **7 · 8** | menjumlahkan estimasi bermata-uang sama → `.GrossAdjustment`, `.AdjustmentValue` |
| **12 · 13** | **membuang** baris spreading yang mata uangnya berbeda (`Property-Remove`) |
| **16 · 17** | memanggil `COUNTGROSSADJTREATY_ACT` lalu `SETPAYABLETREATY_ACT` |

⚠️ Modul Komite Claim Prop menutup butir kurs (K8) dengan bersandar pada rule ini. `[terverifikasi]`
Sandarannya **berdiri**: kunci "isi hanya bila kosong" memang ada. ⚠️ Tetapi langkah **2** berlabel
`//` juga menulis kurs — kalau `//` bukan saklar, ada **dua** penulis kurs, bukan satu.
`[terbuka]` — bergantung pada B3-RALAT.

### D4 — Halaman yang diulang: sumbu lain, dan **prediksi P3 TERBALIK**

Ronde 3 menyatakan halaman yang diulang **tidak terbaca**, karena `pyStepsRepeatType` kosong di
seluruh 1366 langkah dan `pyStepPage` kosong di seluruh 224. **Sumbu itu memang buntu.**

⭐ **Sumbu yang benar adalah `pyStepsObjectName` (Step Page), dan ia terisi.**

| | |
| --- | --- |
| Langkah berulang | **224** — `EMBEDDED` **192** · `REPEAT` **32**, di **63** Activity |
| `pyStepsObjectName` **TERISI** | **194** (87%) |
| `pyStepsObjectName` kosong | **30** — ⭐ **ke-30-nya bertipe `REPEAT`** |
| Halaman berbeda yang diulang | **58 bentuk** |

**Pembagiannya bersih dan masuk akal:**
- **`EMBEDDED` 192 dari 192 punya halaman** → perulangan atas daftar halaman; halamannya terbaca.
- **`REPEAT` 30 dari 32 tidak punya halaman** → perulangan atas **pencacah**, bukan daftar; ia memang
  memakai `pyStepsRepeatDefStart` / `…Limit` / `…Iteration` (terisi di **33** langkah).

Sepuluh halaman terbanyak: `ClaimData.EstimationList` **27** · `ClaimData.SpreadingRisk` **16** ·
`ClaimData.ListClaimAmount` **14** · `pyReportContentPage.pxResults` **12** · `ClaimData.InterestList` **10** ·
`ClaimData.SpreadingClaim` **10** · `ClaimData.FacRetroList` **7** · `ClaimData.AdjustmentList` **7** ·
`ClaimData.SpreadingBreakQS` **5** · `.SpreadingAdjustment` **5**.

> ⛔ **P3 MELESET, dan itu hasil.** Prediksi ronde 3 berbunyi *"halaman yang diulang TIDAK terbaca
> dari struktur, jadi tetap `[terbuka]`"*. **Salah.** Ia terbaca untuk **194 dari 224**. Sebabnya:
> saya menyimpulkan ketiadaan dari **dua** medan tanpa menyisir medan ketiga. Ini **pelajaran yang
> sama persis** dengan ronde 2 Komite (*"sel tidak pernah bergerbang sendiri"*) — **"tidak ketemu"
> hanya sah bila seluruh kedalaman sudah dicari.**

⚠️ **Sekalian: cara baca saya salah jalur.** Pembaca `pyStepsRepeatDef` yang saya tulis mencari
`rowdata` di dalamnya; elemen-elemennya ternyata **anak langsung**. Angka **224 / 192 / 32** tetap
benar karena ronde 3 menghitungnya lewat jalan lain, tetapi pembaca yang salah itu sempat
menghasilkan **0** di ronde ini sebelum ketahuan.

---

## §E — Tambalan: tiket 00 · 04 · 11 · 13

| Tiket | Bagian yang disentuh | Isi tambalan |
| --- | --- | --- |
| **00** | *Rule Pega sumber* + satu AC baru | perkalian ganda share ceding, rantai hulunya, faktor `ShareCeding/100`, dan `[terbuka]` nasib nilai lama |
| **04** | *Catatan `[terbuka]` ringan* | hasil **negatif**: keempat rule `When` bercatatan **tidak** mengubah penggolongan |
| **11** | *Rule Pega sumber* | pengiriman Kasir memakai autentikasi; menunjuk tiket 13 |
| **13** | *Rule Pega sumber* + dua AC baru | autentikasi profil bernama, hanya pada panggilan Kasir, dan peringatan penggantian kredensial |

### E-04 — hasilnya NEGATIF, dan itu dicatat apa adanya

Ronde 3 menduga *"tiga rule `When` yang catatannya menyebut perubahan penggolongan"* perlu diperiksa.
Sensus 100% menemukan **empat**, dan **tidak satu pun** mengubah pembacaan tiket 04:

| Rule | Catatan | Hasil pemeriksaan |
| --- | --- | --- |
| `ASM-FW-GISFW-DATA!ISANEKA` | *"ganti IsBonding jadi IsBondingAndCustomBonds"* | ✅ **sudah dikerjakan** — rujukan yang berlaku adalah `IsBondingAndCustomBonds`; kata `IsBonding` polos hanya tersisa **di dalam teks catatan itu sendiri** |
| `ASM-FW-GISFW-WORK!ISEDMADJSHARECEDANT` | *"save as ganti value"* | ⛔ bukan klasifikasi lini bisnis — menguji status/jenis penawaran |
| `@BASECLASS!ISPEGAPROD` | *"ganti"* | ⛔ bukan klasifikasi — penanda lingkungan |
| `ASM-FW-GCNMFW-WORK-PNC!ISPA_PNC` | *"tambah pyWorkPage dan hapus policy"* | ⛔ kelas `WORK-PNC`, bukan objek kerja Claim Prop — termasuk 52 rule sisa impor |

⛔ **Tidak ada AC tiket 04 yang berubah.** Yang ditambahkan hanya catatan bahwa jalur ini
**sudah diperiksa dan tertutup**, supaya tidak diperiksa ulang.

---

## §F — Sesudah ronde 4

### F1 — Yang dikerjakan berikutnya, urut dari yang paling menahan

| # | Yang dikerjakan | Kenapa | Besar | Menunggu |
| --- | --- | --- | --- | --- |
| **1** | ⛔ **Putuskan arti `//`** | **68 langkah**, termasuk `Obj-Save` dan `Commit`; ia menentukan apakah perkalian ganda pernah benar-benar berjalan, dan apakah `CHECKNOPOLICY` menembak query dua kali | satu keputusan | **work owner** |
| **2** | **Tambal `spec.md`** — perubahan sadar ketiga | spec belum disunting sejak 2026-09-18; sekarang ada tiga perubahan sadar, bukan dua | sekali jalan | — |
| **3** | **Lima titik share ceding di kelas Adjustment** | `COUNTGROSSADJTREATY_ACT` dan `COUNTVALUEADJTREATY_ACT` — jalur ke Kasir, dan **tidak ada catatan yang melarang** | sedang | korpus |
| **4** | **87 catatan `pyMemo` yang belum diuji** | dari 4 larangan langsung yang diuji, **2 masih dikerjakan** — separuh | sedang | korpus |
| **5** | **11 `RDBList` + 11 `Section` + sisa 10 Activity** yang belum dinilai | penyaringan otomatis tidak dapat menilainya | sedang | korpus |

### F2 — Kesimpulan ronde 4 yang PALING RAWAN SALAH

⛔ **Ditunjuk, tidak diperbaiki.**

> ⚠️ **Yang paling rawan: §B4 — *"jalur Kasir tidak terkena perkalian ganda"*.**

Ia bersandar pada **satu** pembacaan: medan uang `TempKasir.CARI9` berisi
`.TotalClaim − .PremiumSpreaded`, dan kedua properti itu bukan `.ClaimSpreaded`. Tetapi saya
**tidak melacak `.TotalClaim` sampai ke hulunya** — hanya memeriksa siapa yang menulisnya secara
langsung. Kalau `.TotalClaim` ternyata diturunkan dari keluarga `Adjustment` yang **juga** mengalikan
share ceding, kesimpulan ini **runtuh**, dan angka yang dikirim ke Kasir ikut terpengaruh.

**Ini pola yang sama yang sudah dua kali membakar proyek ini** (ronde 2 Komite; P3 di ronde ini):
menyatakan sesuatu **tidak** terpengaruh sesudah menyisir **satu** lapis saja.

Yang paling rawan kedua: **§D1 menandai 3 dari 7 butir "bukti lengkap"**. Kata itu berarti *bukti
sudah cukup untuk diputuskan*, **bukan** *butirnya ditutup* — penutupannya keputusan work owner.

### F3 — Pertanyaan BARU untuk work owner — **3**

> ⚠️ Pertanyaan tentang **nasib nilai lama share ceding** sengaja **tidak diulang** — ia sudah
> tergantung sebagai `[terbuka]` di §A2.

#### Pertanyaan 1 ⛔ — Apakah label `//` mematikan langkah di Pega?

**Apa yang ditanyakan.** Di seluruh ekspor **tidak ada satu pun medan yang mematikan sebuah langkah
activity**. Yang ada hanya label `//` pada **68** langkah, dan **32** di antaranya **tanpa syarat apa
pun** — jadi menurut isi ekspor, langkah-langkah itu **berjalan**.

**Kenapa ini menahan.** Di antara ke-32 itu ada **tiga `Obj-Save`** dan **satu `Commit`**. Kalau `//`
bukan saklar, ada penyimpanan ke basis data yang selama ini dianggap mati tetapi sebenarnya hidup.
Ia juga menentukan apakah perkalian ganda share ceding pernah berjalan di produksi — yaitu apakah
data lama perlu diperbaiki.

**Pilihan yang saya lihat.** (a) `//` memang mematikan langkah → 68 langkah itu mati, aturan
ronde 1–3 berdiri; (b) `//` hanya catatan mata pengembang → 32 langkah itu hidup, dan beberapa
kesimpulan lama perlu dibaca ulang.

**Yang saya sarankan.** ⚠️ Saya **tidak menyarankan** — saya tidak punya bukti ke arah mana pun, dan
menebak di sini merusak lebih banyak daripada menunggu.

#### Pertanyaan 2 — `CHECKNOPOLICY` menembak pencarian polis dua kali?

**Apa yang ditanyakan.** Catatan pengembangnya berbunyi *"jangan pake yg di query"*. Langkah **3**
(`Call Rule-Obj-Report-Definition`) dan **4** (`RDB-List`) berlabel `//` **tanpa syarat**; langkah
**7** menjalankan `RDB-List` pengganti. Kalau `//` bukan saklar, pencarian nomor polis berjalan
**tiga kali** per pemanggilan, dan yang dipakai adalah hasil **terakhir**.

**Kenapa ini penting.** Ia menyentuh **tiket 01** (registrasi klaim dan nomor polis), yang belum
pernah menyebut adanya pencarian rangkap. Dan langkah **9** activity itu adalah `Obj-Save`.

**Yang saya sarankan.** Jawab Pertanyaan 1 lebih dulu — pertanyaan ini **sepenuhnya** turunan darinya.

#### Pertanyaan 3 — Lima perkalian share ceding di kelas Adjustment: benar atau ikut cacat?

**Apa yang ditanyakan.** Di luar `COUNTPERSEN_ACT`, share ceding dikalikan **lima kali lagi** —
tiga di `COUNTGROSSADJTREATY_ACT` (langkah 3, 4, 5.2) dan dua di `COUNTVALUEADJTREATY_ACT`
(langkah 12). Kelimanya mengalikan **sekali saja**, dan **tidak ada catatan pengembang yang
melarangnya**.

**Kenapa ini penting.** Keputusan A2 berbunyi *"rumusnya diperbaiki"*. ⚠️ **Tidak jelas apakah A2
mencakup kelima titik ini atau hanya titik 5.1.** Kalau kelimanya ikut "diperbaiki" padahal memang
benar, angka yang dikirim ke **Kasir** akan berubah — dan itu efek keluar.

**Yang saya sarankan.** **A2 hanya mencakup `COUNTPERSEN_ACT` langkah 5.1.** Alasannya: hanya di
situ ada bukti perkalian ganda, dan hanya di situ ada catatan pengembang yang melarang. Kelima titik
lain **disalin apa adanya** sampai ada bukti sebaliknya. ⛔ Tetapi ini **keputusan work owner**.

---

## §E-lampiran — bukti berkas yang tidak bergerak

Sidik jari MD5 diambil **sebelum** penyuntingan dan dibandingkan **sesudahnya**.

### ✅ BERUBAH — tepat **4** berkas, semuanya disebut brief

```
issues/00-prefactor-skema-relasional-transaksi-uang.md
issues/04-klasifikasi-lini-bisnis.md
issues/11-penyerahan-komite-dan-penutupan-klaim.md
issues/13-efek-keluar-kasir-arasapas-konversi-email.md
```

### ⛔ TIDAK BERGERAK — sidik jari MD5 **identik** sebelum dan sesudah

```
spec.md                                       <-- tidak disunting ronde ini
grilling-ronde-1.md
grilling-ronde-2.md
grilling-ronde-3.md
periksa-ulang-aturan-baru.md
issues/01-registrasi-klaim-dan-nomor-polis.md
issues/02-kelunasan-premi-dan-pembebasan-proteksi.md
issues/03-cause-of-loss-dan-katastrofa.md
issues/05-insured-interest-tsi.md
issues/06-loss-allocation-dan-spreading.md
issues/07-estimasi.md
issues/08-baris-adjustment-dan-adjuster.md
issues/09-deductible.md
issues/10-wewenang.md
issues/12-dokumen-pla-dla-acceptance-note.md
issues/14-jejak-audit.md
issues/15-tiga-prefix-klaim-dan-varian-syariah.md
```

**12 tiket + `spec.md` + 4 berkas ronde = 17 berkas terbukti tidak bergerak.**

### Berkas BARU yang ditulis ronde ini — **satu**

```
grilling-ronde-4.md
```
