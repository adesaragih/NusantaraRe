# NB Treaty In — Grilling Ronde 3
## Lima persen rantai nilai dapat dipulihkan, dan sebuah kolom bernama "kelas bisnis" ternyata memuat JSON utuh

Konteks: **Treaty Inward — Realisasi & Endorsement**. Modul: **hanya `NB Treaty In`**.
Korpus READ-ONLY `D:\XML\RNM_BRD\NB Treaty In\`. Aturan: `CLAUDE.md` §4 dan §4a.
Dasar: `grilling-ronde-1.md` dan `grilling-ronde-2.md`.

⛔ **Ronde ini MENEMUKAN dan BERTANYA.** Nol keputusan rancangan · nol kode · nol DDL · nol daftar
kolom usulan · nol nomor baris XML · ⛔ **nol butir ditutup sendiri** · ⛔ **nol nilai nama orang
disalin** · ⛔ **nol nomor OQ dikarang.**

⛔ Berkas ronde 1 dan ronde 2 **TIDAK disunting.** Ralat atas keduanya ditulis **di sini**.

> **SENSUS BERKAS INI**
>
> ⭐ **Jendelanya: berkas ini MINUS blok sensus ini sendiri.** ⛔ Disebut terang-terangan supaya
> blok ini tidak mengubah angka yang ia klaim.
>
> **1.079 baris di jendela** · **13 bab** `## ` · **34 sub-bab** `### ` · **12 sub-sub** `#### ` ·
> **30 tabel** · **16 pernyataan berpenanda terverifikasi** · **1 berpenanda dugaan** ·
> **11 berpenanda terbuka** · ⚠️ **0 berpenanda data DBA** — ⭐ ronde ini membaca korpus, bukan DBA.
>
> ⭐ **Register Bab C — DUA CARA.** *(a)* baris tabel ⇒ **C.1 = 10** butir lama berubah status,
> **C.2 = 16** butir baru; *(b)* nomor di kolom pertama C.2 ⇒ **1–16 tanpa nomor hilang**.
> ✅ **Sepakat.** ⛔ **Nol butir ditutup.**
>
> ⭐ **Pemilik C.2 — DUA CARA** *(baris yang memuat pemilik, dan kemunculan literalnya)*:
> `[pengembang Pega lama]` **8** · `[Product+Underwriting]` **6** · `[DBA]` **1** ·
> `[pemilik export Pega]` **1** · `[IAM]` **1** · `[Finance]` **0** · `[work owner]` **0**.
> ✅ **Sepakat.** ⚠️ Sebutan **17** lawan butir **16** — ⭐ satu butir berpemilik ganda.
>
> ⭐ **Bab D — DUA CARA:** *(a)* pola `**Pn ` ⇒ **12**; *(b)* judul `#### ` ⇒ **12**, bernomor
> **P29–P40 tanpa nomor hilang**. ✅ **Sepakat.** ⭐ Kumulatif tiga ronde: **P1–P40 = 40.**
>
> ⭐⭐ **Ejaan pemilik diperiksa SEBELUM merekap, sesuai peringatan brief:** `[Product+UW]`
> ⇒ **0 kemunculan**, `[Product+Underwriting]` ⇒ **12**. ✅ **Seragam — galat yang terulang di
> ronde 1 dan ronde 2 TIDAK terulang di sini.**
>
> ⭐ **Dan ini ronde pertama yang sensusnya menemukan NOL galat di rekap saya sendiri.**
> ⚠️ Bukan karena saya lebih teliti — ⭐ karena ejaannya diseragamkan **sejak awal**, bukan
> diperiksa sesudah ditulis.
>
> ⛔ **Nol keputusan rancangan · nol kode · nol DDL · nol daftar kolom usulan · nol nomor baris
> XML · nol butir ditutup sendiri · nol nilai nama orang · nol nomor OQ dikarang.**

---

## §0 — Pencocokan dan ujian instrumen

### 0.1 ⭐ Sembilan belas angka dicocokkan — SEMBILAN BELAS COCOK

| Yang dicocokkan | Harap | Terukur | |
| --- | --- | --- | :---: |
| md5 korpus *(isi disambung, urut path)* | `62a3735e…` | **`62a3735ebb2eabb788c8cd8abdda78ca`** | ✅ |
| `DataTransform` · `FlowAction` · `ReportDefinition` · `Flow` | 12 · 9 · 15 · 1 | **12 · 9 · 15 · 1** | ✅ |
| `pyStepsJavaSource` berisi | 8 | **8**, di **7 berkas** *(satu berkas punya 2)* | ✅ |
| `pyStepsObjectName` berisi di `Activity` | 457 | **457** | ✅ |
| `Property-Set` | 942 | **942** | ✅ |
| `Page-Set-Messages` | 105 | **105** | ✅ |
| `RDB-List` | 66 | **66** | ✅ |
| `Page-New` | 33 | **33** | ✅ |
| `Property-Remove` · `Page-Remove` | 25 · 25 | **25 · 25** | ✅ |
| `Property-Set-Messages` · `Page-Clear-Messages` | 17 · 17 | **17 · 17** | ✅ |
| `Java` | 8 | **8** | ✅ |
| `Page-Copy` · `Obj-Save` · `Obj-Browse` | 6 · 6 · 6 | **6 · 6 · 6** | ✅ |
| `pyPropertiesName` di 12 `DataTransform` | 80 | **80** | ✅ |

⭐ **Tambahan yang terukur dan belum pernah disebut:** **99 langkah `Call` bernama** di 92 `Activity`,
dan **92 jenis langkah berbeda** seluruhnya.

⚠️ **Satu keterangan atas angka brief, dan ia soal TAG bukan soal jumlah.** Brief menulis
*"80 nama properti, 62 nilai, 12 aksi"*. ⭐ **80 benar** — `pyPropertiesName`. ⛔ Tetapi **62** dan
**12** berasal dari **tag yang berbeda**: `pyValue` **62** *(⭐ seluruhnya bernilai `3` — kode arah
gerbang, bukan nilai bisnis)* dan `pyAction` **12** *(⭐ seluruhnya `load` — satu per berkas)*.
⭐ **Tag penetapan yang sebenarnya:** `pyPropertiesValue` **53** dan `pyActionName` **80**.
⛔ **Bukan alasan berhenti** — ⭐ dan perbedaan ini justru yang membuat Bab F dapat dihitung.

### 0.2 ⛔⛔ UJIAN INSTRUMEN — dua kegagalan lagi, dan yang KEDUA adalah jebakan yang DIPERINGATKAN

⭐ Ronde 1 gagal 2 kali, ronde 2 gagal 4 kali. ⛔ **Ronde 3 gagal 2 kali.** Total **delapan**,
dan **kedelapannya tertangkap ujian sebelum publikasi.**

| # | Yang alat saya katakan | Yang benar | Sebab galatnya |
| --- | --- | --- | --- |
| ⛔ **1** | *"Nilai `DataTransform` SELAMAT"* — diuji lewat `pyValue`, 62 dengan dan tanpa sisa editor | ⚠️ **Kesimpulannya benar, ujiannya salah.** `pyValue` **bukan** tag penetapan — isinya `3`, kode gerbang. Tag yang benar `pyPropertiesValue` **53** | ⛔ saya menguji **tag yang salah** dan menarik kesimpulan yang kebetulan benar — ⭐ **kebetulan bukan verifikasi** |
| ⛔⛔ **2** | 80 nama dipasangkan ke 53 nilai **menurut urutan**, menghasilkan `.PolicyTreatyIn.EndDate = "1"` dan `.EndDate <- @DateTime…` yang ngawur | ⛔ **SELURUHNYA palsu** | ⛔⛔ **Ini jebakan sensus nomor 6 `CLAUDE.md` §4a, yang berbunyi harfiah:** *"dua daftar dipasangkan menurut urutan — memasangkan nama ke-n dengan tipe ke-n hanya sah bila keduanya datang dari wadah yang sama. Pasangkan lewat KUNCI, bukan lewat urutan."* |

⭐ **Dua percobaan perbaikan, dan hanya yang kedua berhasil:**

| Percobaan | Hasil | Kenapa |
| --- | ---: | --- |
| ⚠️ regex per-`rowdata`, dipasangkan lewat `pyPropertyStepId` | **56 dari 80** baris | ⛔ `rowdata` **bersarang**, dan regex tidak-rakus memungut nilai **anak** ke baris **induk** |
| ⭐ **pengurai XML sungguhan, rekursif, menghormati sarang** | ⭐ **80 dari 80**, nilai **53 dari 53** | ✅ **kedua angka cocok tepat** |

⭐⭐ **Dan inilah bukti bahwa pengurai yang benar memang benar:** sesudah sarang dihormati,
⭐ **`SET` 53 baris dan SELURUH 53 bernilai**, sementara ⭐ `WHEN` 20, `APPLY_MODEL` 3,
`OTHERWISE_WHEN` 2, `APPEND_AND_MAP_TO` 1, `FOR_EACH_PAGE_IN` 1 — ⭐ **nol dari keenam jenis itu
bernilai, dan itu MEMANG SEHARUSNYA**, sebab keenamnya **struktur kendali**, bukan penetapan.
⛔ **Jadi "27 tanpa nilai" bukan data hilang** — ia **wadah**.

⭐ **Aturan baca yang lahir, dan ia mengikat:** ⛔ **Di ekspor Pega, dua daftar sejajar TIDAK
BOLEH dipasangkan menurut urutan, dan pengurai datar TIDAK CUKUP untuk aturan bersarang.**
⭐ Kedalaman sarang terukur di modul ini: **2**.

### 0.3 ⭐ Dua ujian yang LULUS, dilaporkan bersama yang gagal

| Butir uji | Jawaban yang sudah diketahui | Hasil | |
| --- | --- | --- | :---: |
| `pyPropertiesName`/`Value`/`ActionName` di **92 `Activity`** | ronde 2: nilai Activity tidak terekspor | ⭐ **0 · 0 · 0** — ⛔ tag penetapan `DataTransform` **tidak ada sama sekali** di `Activity` | ✅ |
| Jumlah baris terurai lawan cacahan tag | `pyPropertiesName` = 80 | ⭐ pengurai rekursif = **80** | ✅ |

---

## §1 — Sasaran 1: 12 `DataTransform`, satu-satunya penetapan nilai yang terbaca

`[terverifikasi]` **80 baris, 53 di antaranya `SET` bernilai.** Perintah audit: pengurai XML
rekursif atas `rowdata` ber-`pyPropertyStepId`, sesudah `<pyExpressionGadget>` dibuang.

**Sebaran aksi:** `SET` **53** · `WHEN` **20** · `APPLY_MODEL` **3** · `OTHERWISE_WHEN` **2** ·
`APPEND_AND_MAP_TO` **1** · `FOR_EACH_PAGE_IN` **1**.

### 1.1 ⭐⭐ `DeptHeadTreatyInUW_preDT` — dan ia MENJAWAB pertanyaan ronde 2

`[terverifikasi]` `DataTransform\DeptHeadTreatyInUW_preDT.xml`,
`ASM-FW-GISFW-WORK!DEPTHEADTREATYINUW_PREDT`, **7 langkah, seluruhnya `SET`, tanpa sarang:**

| Lgk | Menetapkan | Menjadi |
| ---: | --- | --- |
| 1 | `.PolicyTreatyIn.MasterID` | `""` |
| 2 | `.PolicyTreatyIn.OperatorName` | `OperatorID.pyUserIdentifier` |
| 3 | `.PolicyTreatyIn.SuggestDate` | `@DateTime.CurrentDate("MM/dd/yyyy hh:mm a","")` |
| ⭐ 4 | `.PolicyTreatyIn.IsApproved` | ⭐ **`""`** |
| ⭐ 5 | `.PolicyTreatyIn.isApprovedtoDeptHead` | ⭐ **`""`** |
| 6 | `.PolicyTreatyIn.Suggest` | `""` |
| ⭐⭐ 7 | ⭐ `pyWorkPage.PolicyTreatyIn.QuotationData` | ⭐⭐ **`pyWorkPage.Quotation`** |

⭐⭐ **Langkah 7 menjawab butir terbuka ronde 2 §2.7.** Ronde 2 menemukan `ProportionalType`
dibawa **dua halaman** — `Quotation` dan `QuotationData` — dan menyatakan *"apakah keduanya satu
nilai yang disalin, atau dua nilai yang dapat berbeda, tidak terbaca."*
⭐ **Kini terbaca:** `QuotationData` adalah **SALINAN** `Quotation`, dibuat **pada langkah pra
tahap Kepala Departemen**.
⛔ **Akibatnya, dan ia bukan hal kecil:** ⭐ **keduanya BERBEDA bila `Quotation` berubah sesudah
salinan dibuat** — jadi pembacaan `QuotationData` adalah pembacaan **keadaan pada saat masuk
tahap Kepala Departemen**, bukan keadaan sekarang.
⚠️ **Kapan tepatnya salinan itu dibuat ulang belum terbaca** — `[terbuka]`.

⭐ **Dan langkah 4–5 menemukan medan bendera KEDUA:** `isApprovedtoDeptHead`, ⛔ **berbeda dari
`IsApproved`** dan ⛔ **tidak pernah disebut ronde 1 maupun ronde 2**. Keduanya dikosongkan
bersamaan. ⚠️ Perbedaan keduanya **belum terbaca** — `[terbuka]`.

### 1.2 ⭐⭐ `InputPolicyTreatyIn_preDT` — nama antrean DITETAPKAN di sini

`[terverifikasi]` **20 baris, kedalaman 1.** ⭐ Rule ber-versi **01-01-84**, diperbarui
**23 Juni 2025** — salah satu yang terbaru.

**Lima gerbang pengisi-bila-kosong:**

| Lgk | Gerbang | Menetapkan menjadi |
| ---: | --- | --- |
| 1 | `.PolicyTreatyIn.StartDate==""` | `@DateTime.CurrentDate("dd/MM/yyyy","")` |
| 2 | `.PolicyTreatyIn.EndDate==""` | ⚠️ **`@DateTime.CurrentDate("dd/MM/yyyy","")`** — ⛔ **sama dengan tanggal mulai** |
| 3 | `.PolicyTreatyIn.StatementDate==""` | ⚠️ `@DateTime.CurrentDate("MM/dd/yyyy hh:mm a","")` — ⭐ **format bulan-hari, berbeda dari dua di atas** |
| ⭐ 4 | `.Quotation.ProportionalType=="NonProportional"` | `.PolicyTreatyIn.IsNewPolicyNonProp = "1"` |
| ⭐ 5 | `.PolicyTreatyIn.IsNewPolicyNonProp != "1"` | `= "0"` |

⚠️⚠️ **Langkah 2 layak disebut sendiri:** ⛔ bila tanggal akhir kosong, ia diisi **tanggal hari
ini** — ⭐ **sama dengan tanggal mulai**, bukan satu tahun sesudahnya. ⚠️ Sementara
`SystemSetOneYear_DT` *(§1.4)* memang menghitung satu tahun. ⛔ **Mana yang berlaku kapan belum
terbaca** — `[terbuka]`, dan ⚠️ ini menyentuh **masa berlaku kontrak**.

⚠️ **Dan tiga format tanggal bercampur** — `dd/MM/yyyy` dua kali, `MM/dd/yyyy hh:mm a` sekali di
sini, dan `MM/dd/yyyy hh:mm a` lagi di §1.1. ⛔ **Hari-bulan lawan bulan-hari di satu modul.**

**Sembilan `SET` lurus:**

| Lgk | Menetapkan | Menjadi |
| ---: | --- | --- |
| 6 | `.PolicyTreatyIn.MasterID` | `""` |
| 7 | `.PolicyTreatyIn.MarketingOfficer` | `pyWorkPage.Quotation.MarketingName` |
| 8 · 9 | `.Suggest` · `.SuggestDate` | `""` · tanggal kini |
| 10 | `.PolicyTreatyIn.OperatorName` | ⚠️ **`OperatorID.pyUserName`** |
| 11 | `pyWorkPage.FlagOnGoingPolicy` | `"1"` |
| ⭐⭐ 12 | ⭐ **`TempEmail.CARI28`** | ⛔⛔ **`OperatorID.pyWorkBasketList(2).pyWorkBasketName`** |
| ⭐⭐ 13 | gerbang `pyWorkPage.PositionNote==""` ⇒ **`pyWorkPage.PositionNote`** | ⛔⛔ **`OperatorID.pyWorkBasketList(2).pyWorkBasketName`** |
| 14 | `pyWorkPage.PolicyTreatyIn.QuotationData` | `pyWorkPage.Quotation` |

⭐⭐⭐ **Langkah 13 adalah temuan terpenting §1, dan ia menyentuh OQ-024 secara langsung.**
⭐ OQ-024 berbunyi *"pemetaan Assignment → workbasket tidak terbaca; routing `Custom`"*.
⭐ **Kini penetapannya terbaca:** nama antrean **ditulis ke `PositionNote`**, dan nilainya diambil
dari **entri KEDUA daftar antrean pengguna**. ⛔ **Menurut POSISI, bukan menurut nama.**

⭐⭐ **Dan langkah 12 menjelaskan gerbang yang ronde 2 tidak dapat terangkan.** Ronde 2 §2.5
menemukan gerbang layar `TempEmail.CARI28 != pyWorkPage.PositionNote` **4×** dan hanya dapat
menyebutnya *"slot generik menggerbangi layar"*. ⭐ **Kini terbaca penuh:** keduanya diisi dari
**sumber yang sama** pada langkah 12 dan 13, jadi gerbang itu membandingkan **antrean-ke-2
pengguna SEKARANG** dengan **yang tersimpan saat kasus masuk**. ⭐ Artinya: ⛔ **layar berubah
ketika keanggotaan antrean pengguna berubah.**

⚠️ **Ditambah satu selisih halus:** §1.1 lgk 2 menulis `OperatorName` dari
**`pyUserIdentifier`**, sedangkan §1.2 lgk 10 menulis medan **yang sama** dari
**`pyUserName`**. ⛔ **Dua sumber berbeda untuk satu medan** — `[terbuka]`.

### 1.3 ⭐ `DeptHeadTreatyIn_UW_postDT` — pesan status dipetakan per antrean

`[terverifikasi]` **13 baris, kedalaman 1.** Struktur: **1 `APPLY_MODEL`** memanggil
`AddToListCommentsPolicyTreatyIn_DT`, lalu **6 `WHEN`** yang masing-masing ber-satu `SET`
ke `.NBStatus`.

⭐ **Keenam gerbangnya bertumpu pada `pyWorkPage.PositionNote`** dan menyebut nama antrean sebagai
teks harfiah: `ReasTreatyInAdmin` · `ReasTreatyInDirector` *(2×)* · `ReasTreatyInGroupLeader`
*(2×)* · `ReasTreatyInDeptHead`.

⭐ **Bentuk pesannya terbaca**, dan ada **tiga pola**:

| Pola pesan | Kali |
| --- | ---: |
| `"NB WAS DECLINED BY " + @toUpperCase(OperatorID.pyUserName)` | 1 |
| `"NB IS IN " + @toUpperCase(pyWorkPage.pxCreateOpName) + "'S INBOX"` | 3 |
| ⛔ `"NB IS IN …'S INBOX"` dengan ⛔ **NAMA ORANG DITANAM LANGSUNG** | ⛔ **1** |
| `"NB POLICY NO : " + pyWorkPage.PolicyTreatyIn.PolicyNo` | 1 |

⛔⛔ **Satu pesan memuat nama orang yang ditanam keras di dalam teks pemberitahuan** —
⛔ **nilainya TIDAK saya salin** *(aturan 10)*. ⭐ Yang saya catat: **1 kemunculan**, di
`DataTransform\DeptHeadTreatyIn_UW_postDT.xml`, pada cabang `ReasTreatyInGroupLeader`, di dalam
**teks yang dibaca pengguna**. ⚠️ Ini **pola guard KELIMA** — ⛔ berbeda dari empat yang sudah
tercatat, sebab di sini nama orang bukan penjaga wewenang melainkan **isi pesan**, sehingga
⛔ **mengganti orangnya membuat pemberitahuan berbohong**, bukan membuat akses gagal.

### 1.4 Lima `DataTransform` kecil yang isinya lengkap terbaca

| Rule | Isinya |
| --- | --- |
| ⭐ `SystemSetOneYear_DT` | **1 `SET`**: `.EndDate = @DateTime.addCalendar(.StartDate,1,0,0,0,0,0,0)` — ⭐ **masa berlaku SATU TAHUN, rumusnya terbaca penuh** |
| ⭐ `TestTreatyToFacStatus` | `FOR_EACH_PAGE_IN .PolicyTreatyIn.SpreadingRiskList`, lalu **2 `WHEN`**: `.TreatyType=="10015"` dan `.TreatyType=="10218"` ⇒ `Primary.PolicyTreatyIn.HasFacOut = "1"` |
| ⭐ `TreatyEnableDisableInput` | `SET .QuotationData.ProportionalType = "NonProportional"`, lalu **3 `WHEN`** atas `IsNewPolicyNonProp`: `""` ⇒ `"1"` · `"0"` ⇒ `"1"` · `"1"` ⇒ `"0"` |
| `btnSOB_DT` · `btnCedingCO_DT` | masing-masing **2 `SET`**: `Quotation.btnQuotation` = `SOB` / `CedingCo`, dan `SearchSOB.CARI1 = ""` |
| ⚠️ `setCategoryAttachment_DT` | **1 `SET`**: `Attachment.pxResults(1).NOTE = "LIFE"` — ⚠️ berclass **`ASM-FW-GISFW-Work-LIFE`** |

⚠️ **`TestTreatyToFacStatus`: dua kode angka jenis treaty — `10015` dan `10218` — yang artinya
belum diketahui**, dan keduanya menandai treaty ini **punya fakultatif keluar**. ⛔ Tidak saya tebak.

⭐ **`TreatyEnableDisableInput` adalah pengalih (toggle), dan `""` diperlakukan SAMA dengan `"0"`**
— keduanya menjadi `"1"`. ⚠️ Jadi *"belum diisi"* dan *"tidak"* tidak dapat dibedakan di sini.

⚠️ **`setCategoryAttachment_DT` berclass lini LIFE, di dalam modul treaty inward**, dan menulis
teks harfiah `"LIFE"` ke catatan lampiran **baris pertama**. ⛔ Kenapa ia ada di sini
**belum terbaca** — `[terbuka]`.

### 1.5 `AddToListCommentsPolicyTreatyIn_DT` — mekanisme kronologi

`[terverifikasi]` **6 baris, kedalaman 2** — struktur paling dalam di antara 12:

```
1    APPEND_AND_MAP_TO  .PolicyTreatyIn.SuggestList
1.1     SET   .Suggest      = Param.Comment
1.2     WHEN  Param.Approved != ""
1.2.1      SET  .IsApproved = Param.Approved
1.3     SET   .Date         = Param.Date
1.4     SET   .OperatorName = Param.Operator
```

⭐ **Bacaannya:** setiap catatan dilampirkan ke **`SuggestList`** dengan **empat medan** —
komentar, penanda persetujuan *(hanya bila parameternya tidak kosong)*, tanggal, dan operator.
⭐ Dipanggil `APPLY_MODEL` oleh **`DeptHeadTreatyIn_UW_postDT`** dan **`InboxPolicyTreatyIn_postDT`**.

⚠️ **Tanggal dan operator datang dari PARAMETER**, bukan dari sistem — ⛔ jadi **pemanggilnya yang
menentukan** kapan dan oleh siapa, ⚠️ dan pemanggil dapat mengirim apa saja. `[terbuka]`.

### 1.6 Halaman sasaran penetapan

`[terverifikasi]` `.PolicyTreatyIn` **23** · `pyWorkPage.Quotation` **11** ·
`pyWorkPage.PolicyTreatyIn` **8** · `pyWorkPage` **3** · `pyWorkPage.OfferFacIn` **2** ·
`Primary.PolicyTreatyIn` **2** · `SearchSOB` **2** · `Param` **1** · `TempEmail` **1** ·
`Attachment` **1**.

⚠️ **`.PolicyTreatyIn` (relatif) 23 lawan `pyWorkPage.PolicyTreatyIn` (mutlak) 8** — ⛔ **halaman
yang sama ditulis dua cara**. ⭐ Bentuk relatif berarti "halaman langkah", dan halaman langkah
ditentukan konteks pemanggil. `[terbuka]`.

---

## §2 — Sasaran 2: delapan langkah Java, dan sisi BACA dari OQ-012

`[terverifikasi]` **8 langkah Java di 7 berkas** — `GetTreatyName` punya **2**, enam berkas lain
punya **1**. Perintah audit: `pyStepsJavaSource` berisi, sesudah sisa editor dibuang.

### 2.1 ⛔⛔ Enam dari delapan adalah KODE YANG SAMA — dan ia membaca JSON dari kolom bernama "kelas bisnis"

`[terverifikasi]` Enam langkah di `FetchMasterTreatyIn` · `InputPolicyTreatyInDetail_NonProp` ·
`InputPolicyTreatyInDetail_preACT` · `InputPolicyTreatyOutDetail_NonProp` ·
`InputPolicyTreatyOutDetail_preACT` · `SetTreatyIn_Act` **identik isinya**, berbeda hanya pada
halaman sasaran. Dikutip **dua baris intinya saja**, seperlunya untuk membuktikan klaim:

```
String IsiDataJson = DataJSONPage.getString("CLASSOFBUSINESS");
tempPage2.adoptJSONObject(IsiDataJson);
```

⭐⭐ **Inilah sisi BACA dari OQ-012.** Ronde 1 menemukan `@ASM.GetPageJSONString()` **menulis**
JSON ke basis data, dan source-nya tidak ada. ⭐ **Ronde 3 menemukan yang MEMBACANYA** — dan
hasilnya lebih keras dari dugaan:

| Yang terbaca | Artinya |
| --- | --- |
| ⛔⛔ Kolom bernama **`CLASSOFBUSINESS`** dibaca sebagai **teks**, lalu **diadopsi sebagai objek JSON** | ⛔ **sebuah kolom yang namanya berarti "kelas bisnis" memuat DOKUMEN JSON UTUH** |
| ⭐ Sasarannya halaman **`TreatyIn`** *(3 berkas)* atau **`pyWorkPage.TreatyIn`** *(3 berkas)* | ⭐ dua ejaan halaman yang sama |
| ⭐ Sumbernya daftar **`TempResult.pxResults`**, dilingkari satu per satu | ⭐ **berasal dari hasil query**, jadi kolomnya memang kolom basis data |
| ⛔⛔ `adoptJSONObject` | ⛔ **struktur halaman ditentukan oleh isi JSON pada saat jalan** — ⭐ **itulah sebabnya skema `TreatyIn` tidak dapat dinyatakan** |

⭐ **Dan `CLASSOFBUSINESS` terbukti memang sebuah kolom**, bukan nama karangan: ⭐ ia muncul di
antara **66 kolom** `BrowseTreatyInDetail` *(§3.2)*. ⛔ Artinya **kolom yang sama dipakai dua cara**
— sebagai kolom biasa di laporan, dan sebagai wadah JSON di Java. `[terbuka]`.

⛔⛔ **Penanganan kegagalannya menelan galat:**

```
} catch(Exception e){
  oLog.error("ReloadSection:Expection : "+e.getMessage());
}
```

⛔ **Galat hanya dicatat ke log, lalu aktivitas LANJUT.** ⭐ Jadi JSON yang rusak menghasilkan
halaman **kosong atau separuh terisi**, ⛔ **tanpa pesan kepada pengguna dan tanpa menghentikan
langkah berikutnya.** ⚠️ Dan pesan log-nya berbunyi `"ReloadSection:"` — ⭐ **nama konteks yang
berbeda**, tanda kode ini **disalin dari tempat lain**, sehingga log **menunjuk sumber yang salah**.

### 2.2 ⭐ Dua langkah Java sisanya: penghapus duplikat

`[terverifikasi]` Kedua langkah di `Activity\GetTreatyName.xml` *(268.638 B)* **membuang baris
kembar** dari sebuah daftar, memakai `HashSet` dan menelusuri **dari belakang ke depan**.

| Java | Daftar yang disaring | Kunci kembarnya |
| --- | --- | --- |
| #1 | `GetTreatyName.pxResults` | ⚠️ **`.CARI2`** — ⭐ slot generik lagi |
| #2 | `SpreadingList1.pxResults` | `.TreatyName` |

⚠️ **Komentarnya berbohong tentang kodenya:** komentar berbunyi *"read all 10 properties"* dan
*"concatinate all 10 properties… separated by #"*, ⛔ **tetapi kodenya membaca SATU properti saja**
dan tidak menyambung apa pun. ⭐ Sisa salin-tempel. ⚠️ ⛔ Ini **jenis kesalahan yang sama** seperti
temuan ronde 2 §4.2 *(label manusia berselisih dengan aturan yang dieksekusi)* — ⭐ **kali ini di
dalam komentar kode.**

⭐ **Nol dari 8 langkah Java menghitung uang**, dan ⭐ **nol memanggil apa pun di luar Pega** —
keduanya hanya memakai API klipboard Pega dan `HashSet` Java. ⛔ **Jadi Java BUKAN tempat rumus
uang bersembunyi.**

---

## §3 — Sasaran 3: 37 berkas yang tadinya nol dibuka

### 3.1 ⛔⛔ Sembilan `FlowAction` — dan SEMBILANNYA KOSONG

`[terverifikasi]` Kesembilan `FlowAction` disisir. ⛔ **Tidak satu pun memuat:**

| Yang dicari | Ditemukan |
| --- | ---: |
| ⛔ aturan validasi *(`pyValidateName`)* | ⛔ **0** |
| ⛔ hak-akses *(`pyPrivilegeName`)* | ⛔ **0** |
| ⛔ gerbang penampil *(`pyWhenName`)* | ⛔ **0** |
| ⛔ transformasi data pra / pasca *(`pyPreDataTransform` / `pyPostDataTransform`)* | ⛔ **0 berisi** — ⭐ tag-nya ada di kesembilan, tetapi isinya hanya `<pxObjClass>Embed-Invoke-DataTransform</pxObjClass>` dengan `<pyName/>` **kosong** |
| ⛔ aktivitas pra / pasca | ⛔ **0** |

⭐⭐ **Ini penegasan OQ-007 dari sudut KETIGA.** Ronde 2 membuktikan **295 `pyAssociatedPrivileges`
nol berisi** di layar. ⭐ Ronde 3 membuktikan **kesembilan aksi pindah-tahap juga tanpa hak-akses
dan tanpa validasi.** ⛔ **Jadi tidak ada satu pun titik di modul ini yang menegakkan wewenang
lewat mekanisme Pega** — ⭐ seluruhnya bertumpu pada gerbang dan nama orang yang sudah tercatat.

⚠️⚠️ **Dan ini melahirkan pertanyaan yang tidak dapat saya jawab:** ronde 1 mencatat bahwa alur
memakai `InboxPolicyTreatyIn_postDT` dan `DeptHeadTreatyIn_UW_postDT` sebagai transformasi pasca.
⛔ **Keduanya TIDAK terpasang di `FlowAction`-nya.** ⭐ Jadi ia terpasang **di aturan `Flow`**, bukan
di aksinya. ⚠️ **Konsekuensinya: aksi yang sama dipakai di beberapa tahap dengan transformasi
pasca yang berbeda**, dan mana-berpasangan-dengan-mana **hanya terbaca dari `Flow`**. `[terbuka]`.

⭐ **Label yang terbaca** — ⭐ **dua di antaranya berbeda dari nama rule-nya**, dan itulah yang
dilihat pengguna:

| Rule | Label di layar |
| --- | --- |
| ⭐ `PolicyTreatyInDeclineConfirm` | ⭐ **"Confirm Decline NB"** |
| ⭐ `ShowPolicyNoTreaty` | ⭐ **"Show PolicyNo"** |
| tujuh lainnya | sama dengan nama rule-nya |

⚠️ **Kesembilan `FlowAction` tinggal di TUJUH class berbeda** — `…WORK` *(2)*,
`…DATA-QUOTATION` *(2)*, `…DATA-POLICYTREATYIN` *(2)*, `…DATA-AGENT`, `…DATA-INSTALLMENT`,
`…DATA-TREATYININSTALLMENT`. ⭐ Hanya **dua** berclass kerja utama.

### 3.2 Lima belas `ReportDefinition`

`[terverifikasi]` Kelima belas disisir. ⭐ **Tiga di antaranya berkolom 66** dan itulah yang
memberi gambaran terluas tentang bentuk data treaty:

| Laporan | Kelas | Kolom |
| --- | --- | ---: |
| ⭐ `BrowseTreatyInDetail` | `…Int-TREATYINDETAIL` | ⭐ **66** |
| ⭐ `BrowseTreatyJoinEDM` | `…Int-TREATYINDETAILJOINEDM` | ⭐ **66** |
| ⚠️ `BrowseTreatyOutDetail` | ⚠️ **`…Int-TREATYOUTDETAIL`** | ⭐ **66** |
| ⭐ `BrowseTREATY_IN` | `…Int-TREATY_IN` | **32** |
| `BrowseAgentNusaRe_RD` | `…Int-AGENT` | 34 |
| `BrowseMarketingOfficer_RD` | `…Int-marketingofficer` | 24 |
| `BrowseCurrency_RD` · `BrowseCurrencyTreatyIn_RD` | `…Int-CURRENCY` | 16 · 16 |
| `BrowseClientName_RD` · `BrowseReinsuranceType_RD` | `…Int-AGENT` · `…Int-REINSURANCETYPE` | 14 · 14 |
| `BrowseTreatyGroup_RD` · `BrowseAgentHierarkiList_RD` | `…Int-TREATYGROUP` · `…Int-AGENT` | 12 · 10 |
| `BrowseCedingCo_RD` | `…Int-AGENT` | 8 |
| ⚠️ `GetListOpportunity` · `crmOpportunitiesList` | ⚠️ **`ASM-FW-SFAGISFW-Work-Opportunity`** | 51 · 43 |

⭐⭐ **`BrowseTREATY_IN` memperlihatkan bahwa medan tangga akseptasi adalah KOLOM**, bukan hanya
properti di memori: `.Position` · `.PositionUsername` · `.StatusAkseptasi` ·
`.ChooseStatusAkseptasi` · `.ProportionType` · `.Commencement` · `.Termination` ·
`.ClassofBusiness` · `.Ceding` · `.CedingID` · `.TreatyContractName` · `.Information`.

⚠️ **Tiga laporan berkolom 66 nyaris identik daftar kolomnya** — `.TREATYID`, `.PROPORTIONTYPE`,
`.TREATYTYPE`, `.TREATYGROUP`, ⭐ **`.CLASSOFBUSINESS`**, `.LIMITCURRENCY`, `.LIMITVALUE`,
`.RETENTIONCURRENCY`, `.RETENTIONVALUE`, `.EPICURRENCY` — dan **penyaringnya juga identik**
*(`Param.SOB`, `Param.TREATYYEAR`, `Param.CEDING`, `Param.TREATYTYPE`, …)*.
⭐ Jadi **treaty masuk, treaty masuk-gabung-EDM, dan treaty KELUAR dibaca dengan bentuk yang sama.**
⚠️ Ini memperkuat **OQ-022** *(modul ini penyentuh treaty outward terbanyak)* dengan bukti baru:
⭐ **bentuk datanya pun dibagi.**

⭐ **Pasangan mata uang yang bersanding:** `.LIMITCURRENCY`+`.LIMITVALUE`,
`.RETENTIONCURRENCY`+`.RETENTIONVALUE`, `.EPICURRENCY` — ⭐ **nilai uang di lapisan ini MEMANG
bermata uang**, berpasangan kolom.

⚠️ **Penyaring ter-hardcode yang terbaca:**

| Nilai | Di mana | Catatan |
| --- | --- | --- |
| ⚠️ **`"ITL"`** | `BrowseCurrency_RD` · `BrowseCurrencyTreatyIn_RD` | ⭐ sebuah kode mata uang **dikecualikan secara khusus** |
| ⚠️ `"LIFE INSURANCE"` | `BrowseAgentNusaRe_RD` | teks harfiah |
| `"NB-"` · `"Resolved-Completed"` · `"Resolved-Rejected"` | `GetListOpportunity` | ⭐ awalan nomor kasus dan **dua status penyelesaian** |
| `1` | `BrowseAgentNusaRe_RD` · `BrowseMarketingOfficer_RD` | ⭐ **angka telanjang**, artinya belum diketahui |

⛔ **Nol nomor polis ter-hardcode di kelima belas laporan** — ⭐ pola ronde 2 §4.4 **tidak berulang
di sini**. Perintah audit: pola `"RNM-[A-Z0-9.\-]{8,}"` atas `ReportDefinition\` ⇒ **0**.

⛔ **Nol nama orang di kelima belas laporan.** ⭐ Yang ada hanya `OperatorID.pyUserName` dan
`Param.UserIdentifier` — **rujukan, bukan nilai**.

### 3.3 Satu `Flow` — konfirmasi, bukan telusur ulang

`[terverifikasi]` `Flow\InputRealizationTreatyIn.xml`, `ASM-FW-GISFW-WORK!INPUTREALIZATIONTREATYIN`.
⭐ **Tidak ditelusur ulang** sesuai perintah brief. ⭐ Yang dikonfirmasi ronde 3: ⛔ **kesembilan
`FlowAction` tidak memuat transformasi pra/pasca**, sehingga ⭐ **penyambungan tahap ke transformasi
memang hanya hidup di berkas `Flow` ini** — sesuai telusur D2, dan **bukan** temuan baru.

---

## §4 — Sasaran 4: bentuk 92 `Activity`, dan dua aturan yang HILANG

### 4.1 Rantai panggil antar-`Activity`

`[terverifikasi]` **99 sisi panggil**; **41** `Activity` memanggil, **70** nama dipanggil.
Perintah audit: langkah ber-`pyStepsActivityName` berpola `call [<Class>.]<Nama>`, sesudah sisa
editor dibuang.

| | Jumlah |
| --- | ---: |
| nama yang dipanggil | **70** |
| ⭐ **ADA** di modul ini | ⭐ **66** |
| ⛔ **TIDAK ADA** di modul ini | ⛔ **4** |

⭐ **Pusat rantainya terbaca:** **`CountNetPremi_act` dipanggil 10 `Activity` berbeda** — ⭐ **titik
paling ramai di seluruh modul**. Disusul `CountSpreading_Act` **4**, lalu `CountResult1_Act` ·
`FetchTreatyGroupOldID` · `FetchTreatyGroupOJK` · `InsertToTreatyXOLList` masing-masing **3**.

⭐⭐ **Dan inilah kerangka yang brief harapkan:** ⭐ **rantai panggil terbaca utuh walau nilai
langkahnya tidak.** ⛔ Jadi bila ekspor ulang datang, **isinya tinggal ditempelkan ke kerangka
99 sisi ini** — ⭐ kerangkanya **tidak perlu dicari lagi.**

### 4.2 ⛔⛔ Empat target yang TIDAK ADA — dua di antaranya aturan bisnis

| Target | Dipanggil oleh | Sifatnya |
| --- | --- | --- |
| ⛔⛔ **`CopyToPolicy`** | ⭐ **`SaveJsonPolisTreatyIn_Act`** | ⛔⛔ **aturan bisnis, dan ia di JALUR SIMPAN INTI** |
| ⛔⛔ **`SumTSIPremiSpreadRNMMultiCob_Act`** | ⭐ **`SumTSIPremiSpreadedRNM_Act`** | ⛔⛔ **aturan bisnis, dan ia di RANTAI UANG** |
| `pxRetrieveReportData` | `FetchTreatyGroupOJK` · `GetCurrencyMaster` · `InputPolicyTreatyInDetail_preACT` | ⭐ bawaan Pega — **wajar tidak ada** |
| `pxShowReport` | `AgentSourceBizTreatyIn_Act` · `CheckDuplicateOffer` | ⭐ bawaan Pega — **wajar tidak ada** |

⛔⛔ **`CopyToPolicy` layak disebut sendiri.** ⭐ Telusur D2 sudah mendaftarkannya sebagai
*"dirujuk tapi belum ditelusur"*. ⭐ **Ronde 3 menaikkan statusnya:** ia bukan *belum ditelusur*,
melainkan ⛔ **tidak ada di modul ini sama sekali** — dan pemanggilnya adalah aktivitas yang
menyimpan JSON polis, ⭐ **langkah inti jalur realisasi**.

⛔⛔ **`SumTSIPremiSpreadRNMMultiCob_Act` sama beratnya.** ⭐ Pemanggilnya
`SumTSIPremiSpreadedRNM_Act` — ⭐ **aturan yang sama yang memuat 24 nomor polis ter-hardcode**
*(ronde 2 §4.4)*. ⛔ Jadi rantai uang **memanggil keluar ke aturan yang tidak ada**, dan
⭐ singkatan `MultiCob` `[dugaan]` menyarankan *"banyak class of business"* — ⛔ **tidak saya tebak.**

⭐ **Bersama OQ-025** *(`SERVICEINSERTARASAPAS_ACT`)*, ⛔ **kini ADA TIGA aturan bisnis yang
dipanggil modul ini dan tidak ada di dalamnya.**

### 4.3 Bentuk yang terbaca dari 92 `Activity`

`[terverifikasi]` **92 jenis langkah berbeda**; dua belas terbanyak **cocok seluruhnya** dengan
angka brief *(§0.1)*. ⭐ Yang dapat dibaca tentang bentuknya:

| Yang terbaca | Jumlah | Nilainya |
| --- | ---: | --- |
| `Property-Set` | **942** | ⛔ **tahu ADA, tidak tahu APA** |
| ⭐ `pyStepsObjectName` | ⭐ **457** | ⭐ **halaman tempat langkah bekerja — TERBACA** |
| `Page-Set-Messages` + `Property-Set-Messages` | **122** | ⭐ **122 tempat memasang pesan galat** |
| `RDB-List` | **66** | jembatan ke 41 rule SQL |
| `Page-New` + `Page-Remove` + `Page-Copy` | 33 + 25 + 6 | daur hidup halaman |
| `Property-Remove` | 25 | penghapusan properti |
| `Page-Clear-Messages` | 17 | ⚠️ **17 tempat MENGHAPUS pesan galat** |
| `Obj-Save` + `Obj-Browse` | 6 + 6 | ⭐ **hanya 6 titik simpan** di 92 aktivitas |
| `Java` | 8 | §2 |
| langkah `Call` bernama | **99** | §4.1 |

⚠️⚠️ **`Page-Clear-Messages` 17 lawan pemasang pesan 122** layak ditanyakan: ⛔ **tujuh belas
tempat menghapus pesan galat**, dan ⭐ **menghapus pesan galat berarti melanjutkan sesuatu yang
tadi ditandai salah.** ⛔ Apakah itu pembersihan yang sah atau penelanan galat **tidak terbaca** —
`[terbuka]`.

⛔ **Yang TIDAK saya coba, dan sengaja:** ⛔ **nol rumus aritmetika saya simpulkan dari 92
`Activity`**, sebab ronde 2 membuktikan nilainya tidak terekspor. ⭐ **Ronde 3 menghasilkan
STRUKTUR KEPUTUSAN, bukan rumus** — dan itu dinyatakan di sini supaya tidak ada yang menyangka
rantai uang sudah terbaca.

---

## Bab A — Tambahan stored procedure dan SQL mentah

⭐ **NIHIL tambahan, dibuktikan.**

`[terverifikasi]` Jendela: **12 `DataTransform` + 9 `FlowAction` + 15 `ReportDefinition` + 1 `Flow`
+ 8 langkah Java** = **37 berkas + 8 blok Java**, sisa editor dibuang.
Pola: `(POOLDATA|DATAPEGA|ARASAPAS|ASM)\.[A-Z0-9_]+\s*\(` dan `BEGIN`/`COMMIT`/`INSERT INTO`.

| Hasil | Jumlah |
| --- | ---: |
| ⭐ procedure Oracle **baru** | ⭐ **0** |
| ⭐ blok PL/SQL baru | ⭐ **0** |
| ⭐ panggilan luar-Pega dari Java | ⭐ **0** — ⭐ kedelapan langkah hanya memakai API klipboard Pega dan `HashSet` |

⭐ **Ketiga procedure ronde 1 tetap ketiganya, dan nol badannya masih ada.**

## Bab B — Tambahan objek Oracle

⭐ **NIHIL objek baru, tetapi ⭐ SATU KOLOM dikenali perannya.**

`[terverifikasi]` **0 nama tabel baru.** ⭐ **Tetapi `ReportDefinition` memperlihatkan bentuk
kolom** dari objek yang ronde 1 sudah daftarkan — dan satu di antaranya menentukan:

| Kolom | Yang baru diketahui |
| --- | --- |
| ⛔⛔ **`CLASSOFBUSINESS`** | ⭐ ia **salah satu dari 66 kolom** `BrowseTreatyInDetail`, ⛔ **dan Java membacanya sebagai DOKUMEN JSON** — §2.1 |
| ⭐ `.LIMITCURRENCY`+`.LIMITVALUE` · `.RETENTIONCURRENCY`+`.RETENTIONVALUE` · `.EPICURRENCY` | ⭐ **nilai uang berpasangan dengan kolom mata uangnya** di lapisan ini |
| ⭐ `.Position` · `.PositionUsername` · `.StatusAkseptasi` · `.ChooseStatusAkseptasi` | ⭐ **medan tangga akseptasi adalah KOLOM**, bukan hanya properti di memori |

⛔ **Nol daftar kolom usulan saya buat** — ⭐ yang di atas adalah **kolom yang terbaca**, bukan
usulan skema.

---

## Bab C — Register `[terbuka]` berpemilik

⭐ Nomor OQ dari `discovery/open-questions.md`. ⭐ **72 OQ terdaftar; nomor bebas berikutnya
OQ-073.** Butir ronde 1 dan 2 dirujuk **"ronde 1 baru #n"** dan **"ronde 2 baru #n"**.

⛔⛔ **NOL butir saya tutup.**

⚠️ **Ejaan pemilik diseragamkan SEJAK AWAL** menjadi `[Product+Underwriting]` — ⭐ sesuai
peringatan brief, sebab galat ini **sudah terulang dua ronde**. Dicek dengan cacahan sebelum
merekap, bukan sesudah.

### C.1 Butir lama yang BERUBAH STATUSNYA

⛔ **Berubah = diperdalam atau dipersempit. ⛔ BUKAN ditutup.**

| Butir | Perubahannya | Pemilik |
| --- | --- | --- |
| ⭐⭐ **OQ-012** | ⭐⭐ **Sisi BACA ketemu.** Kolom **`CLASSOFBUSINESS`** dibaca sebagai teks lalu **diadopsi sebagai JSON** ke halaman `TreatyIn`, di **6 tempat**. ⭐ Menjelaskan **kenapa skema tidak dapat dinyatakan**: struktur halaman ditentukan isi JSON saat jalan. ⛔ Struktur JSON-nya sendiri **tetap tidak diketahui** | `[DBA]` · `[pengembang Pega lama]` |
| ⭐⭐ **OQ-024** | ⭐⭐ **Dipersempit tajam.** Penetapan nama antrean **TERBACA**: `PositionNote` diisi dari `OperatorID.pyWorkBasketList(2).pyWorkBasketName` bila masih kosong. ⛔ **Menurut POSISI ke-2, bukan menurut nama** | `[IAM]` · `[Product+Underwriting]` |
| ⭐ **OQ-007** | ⭐ **Diperkuat dari sudut KETIGA.** Kesembilan `FlowAction` **nol hak-akses, nol validasi, nol gerbang** | `[IAM]` |
| ⭐ **OQ-022** | ⭐ **Diperkuat dengan bukti baru.** `BrowseTreatyOutDetail` berkolom **66 dan penyaring yang IDENTIK** dengan versi treaty masuk — ⭐ **bentuk datanya pun dibagi**, bukan hanya tabelnya | `[Product+Underwriting]` |
| ⭐ **OQ-025** | ⭐ **Konteksnya berubah, blocker tidak.** ⛔ Kini ADA **TIGA** aturan bisnis yang dipanggil tetapi tidak ada di modul — ⛔ `SERVICEINSERTARASAPAS_ACT`, **`CopyToPolicy`**, **`SumTSIPremiSpreadRNMMultiCob_Act`** | `[pemilik export Pega]` · `[Product+Underwriting]` |
| ⭐ **ronde 2 baru #1** *(942 langkah)* | ⭐ **Dikuantifikasi.** Lihat **Bab F** — ⭐ **5,3 % rantai nilai dapat dipulihkan**, ⛔ **94,7 % tetap hilang** | `[pemilik export Pega]` |
| ⭐ **ronde 2 baru #7** *(antrean menurut posisi)* | ⭐ **Diperkuat dari BACA menjadi TULIS.** Ronde 2 menemukan posisi-1 dan posisi-2 **dibaca**; ⭐ ronde 3 menemukan posisi-2 **DITULIS** ke `PositionNote` dan ke `TempEmail.CARI28` | `[IAM]` |
| ⭐ **ronde 1 baru #1** *(slot generik)* | ⭐ **Dua slot lagi terjelaskan:** `TempEmail.CARI28` = nama antrean ke-2 · `GetTreatyName…CARI2` = kunci pembuang-duplikat. ⛔ Sebagian saja | `[pengembang Pega lama]` |
| ⭐ **ronde 2 §2.7** *(dua halaman `ProportionalType`)* | ⭐⭐ **TERJAWAB dari korpus:** `QuotationData` adalah **SALINAN** `Quotation`, dibuat di langkah pra Kepala Departemen. ⛔ Kapan disalin ulang **tetap terbuka** | `[pengembang Pega lama]` |
| ⛔ **OQ-002 · OQ-013** | ⭐ tidak berubah. ⛔ **Bab F membuatnya lebih menekan** | `[DBA]` |

⚠️ **Butir lama yang TIDAK berubah di ronde ini:** OQ-026, OQ-059, OQ-023, OQ-028, OQ-020,
OQ-021, OQ-011, OQ-018, OQ-009, OQ-054, dan ronde 2 baru #2 · #3 · #4 · #5 · #6 · #8 · #9 · #10 ·
#11 · #12 · #13, serta ronde 1 baru #2 · #3 · #4 · #5 · #6 · #7. ⛔ Tidak diulang di sini.

### C.2 ⭐ Butir yang LAHIR di ronde 3

| # | Butir | Pemilik | Kenapa memblokir | Akibat bila dijawab salah |
| --- | --- | --- | --- | --- |
| ⛔⛔ **1** | **Kolom bernama `CLASSOFBUSINESS` memuat dokumen JSON utuh**, diadopsi ke halaman treaty di 6 tempat; kolom yang sama juga tampil sebagai kolom biasa di laporan | `[DBA]` | ⛔⛔ **struktur data treaty ditentukan saat jalan** — skema tidak dapat dinyatakan | ⛔ data lama tidak terbaca; atau kolom disalin sebagai teks biasa dan isinya hilang |
| ⛔⛔ **2** | **Dua aturan bisnis dipanggil tetapi TIDAK ADA di modul** — satu di jalur simpan inti, satu di rantai uang | `[pemilik export Pega]` | ⛔⛔ jalur simpan dan rantai uang **terputus di dua titik** | ⛔ langkah yang hilang tidak dibangun, dan tidak ada yang tahu ia pernah ada |
| ⛔⛔ **3** | **Galat JSON hanya dicatat ke log lalu aktivitas LANJUT** — `catch` tanpa penghentian, dan pesan log-nya menyebut konteks yang salah | `[pengembang Pega lama]` | ⛔ halaman separuh terisi **diperlakukan sebagai sah** | ⛔⛔ perhitungan jalan atas data tidak lengkap **tanpa satu pun pesan** |
| ⛔ **4** | **Tanggal akhir diisi TANGGAL HARI INI bila kosong**, sementara rule lain menghitung **satu tahun** | ⭐ `[Product+Underwriting]` | ⛔ **masa berlaku kontrak** | ⛔ kontrak tersimpan berakhir di hari ia dibuat |
| ⛔ **5** | **Tiga format tanggal bercampur** — `dd/MM/yyyy` dan `MM/dd/yyyy hh:mm a` di satu modul | `[pengembang Pega lama]` | ⛔ hari-bulan lawan bulan-hari | ⛔⛔ tanggal 3 Februari terbaca 2 Maret, **diam-diam** |
| ⛔ **6** | **Medan bendera KEDUA ditemukan** — `isApprovedtoDeptHead`, berbeda dari `IsApproved`, dikosongkan bersamaan | `[Product+Underwriting]` | ⛔ dua bendera persetujuan, arti bedanya tak terbaca | ⛔ persetujuan tercatat di bendera yang salah |
| ⛔ **7** | ⛔ **Nama orang ditanam di dalam TEKS PEMBERITAHUAN** — pola guard **kelima**, 1 kemunculan | ⭐ `[IAM]` · `[Product+Underwriting]` | ⚠️ berbeda dari empat pola sebelumnya: ini **isi pesan**, bukan penjaga akses | ⛔ pemberitahuan **berbohong** kepada pengguna setelah orangnya berganti |
| ⚠️ **8** | **Satu medan diisi dari dua sumber berbeda** — `OperatorName` dari `pyUserIdentifier` di satu rule, `pyUserName` di rule lain | `[pengembang Pega lama]` | ⚠️ satu kolom, dua arti | ⚠️ jejak audit mencatat dua jenis pengenal di kolom yang sama |
| ⚠️ **9** | **Tanggal dan operator kronologi datang dari PARAMETER**, bukan dari sistem | `[pengembang Pega lama]` | ⚠️ pemanggil menentukan siapa dan kapan | ⚠️ kronologi dapat mencatat waktu atau orang yang tidak benar |
| ⚠️ **10** | **Dua kode angka jenis treaty menandai "punya fakultatif keluar"** — `10015` dan `10218` | `[Product+Underwriting]` | ⚠️ arti kode belum diketahui | ⚠️ treaty salah ditandai, jalur fakultatif keluar salah jalan |
| ⚠️ **11** | **Kesembilan `FlowAction` tanpa transformasi pra/pasca**, padahal alur memakainya — penyambungannya hanya di berkas `Flow` | `[pengembang Pega lama]` | ⚠️ aksi yang sama dipakai beberapa tahap dengan pasangan berbeda | ⚠️ tahap dipasangkan ke transformasi yang salah |
| ⚠️ **12** | **17 langkah MENGHAPUS pesan galat**, lawan 122 yang memasangnya | `[pengembang Pega lama]` | ⚠️ menghapus pesan galat = melanjutkan yang tadi ditandai salah | ⚠️ validasi dilewati tanpa jejak |
| ⚠️ **13** | **Rule berclass lini LIFE di dalam modul treaty inward**, menulis teks harfiah ke catatan lampiran baris pertama | `[Product+Underwriting]` | ⚠️ batas konteks | ⚠️ lampiran lini lain tergolong salah |
| ⚠️ **14** | **Komentar kode berselisih dengan kodenya** — komentar menyebut 10 properti disambung, kode membaca 1 | `[pengembang Pega lama]` | ⚠️ jenis kekeliruan yang sama seperti ronde 2 §4.2, kali ini di komentar | ⚠️ pembaca berikutnya membangun dari komentar, bukan dari kode |
| ⚠️ **15** | **Satu kode mata uang dikecualikan khusus** di dua laporan | `[Product+Underwriting]` | ⚠️ pengecualian ter-hardcode | ⚠️ mata uang yang seharusnya tersedia hilang, atau sebaliknya |
| ⚠️ **16** | **Halaman yang sama ditulis dua cara** — relatif `.PolicyTreatyIn` **23×** lawan mutlak `pyWorkPage.PolicyTreatyIn` **8×** | `[pengembang Pega lama]` | ⚠️ bentuk relatif bergantung konteks pemanggil | ⚠️ penetapan mendarat di halaman yang salah |

### C.3 Rekap per pemilik — dihitung DUA CARA

⭐ **16 butir lahir di ronde 3.** ⛔ **0 ditutup.**

| Pemilik | Butir ronde 3 |
| --- | ---: |
| ⛔ `[pengembang Pega lama]` | ⭐ **8** |
| ⛔ `[Product+Underwriting]` | ⭐ **6** |
| ⛔ `[DBA]` | **1** |
| ⛔ `[pemilik export Pega]` | **1** |
| ⚠️ `[IAM]` | **1** |
| ⭐ `[Finance]` | ⭐ **0** — lihat C.4 |
| `[work owner]` | **0** |

### C.4 ⭐ `[Finance]` = 0 di ronde 3, dan sebabnya WAJIB disebut

⛔ **Brief menuntut ini eksplisit.**

⭐ **Ronde 1: 0** *(lubang)* → ⭐ **Ronde 2: 3** *(lubang terisi, isinya buruk)* → ⭐ **Ronde 3: 0.**

⛔ **Dan nol kali ini BERBEDA ARTINYA dari nol ronde 1.** Sebabnya:
⭐ **Ronde 3 memang tidak menyentuh uang, dan itu disengaja** — Sasaran 1 sampai 4 seluruhnya
tentang **bentuk**, bukan aritmetika. ⛔ Rumus uangnya **tidak dapat disentuh** sebab
ronde 2 membuktikan nilainya tidak terekspor.

⚠️ **Tetapi dua butir ronde 3 menyentuh uang secara TIDAK LANGSUNG**, dan keduanya berpemilik lain:
1. ⛔ **Butir #2** — rantai uang **memanggil aturan yang tidak ada** *(`SumTSIPremiSpreadRNMMultiCob_Act`)*
2. ⛔ **Butir #3** — data separuh terisi **diteruskan tanpa pesan**, dan perhitungan jalan di atasnya

⭐ **Jadi `[Finance]` tetap 3 butir dari ronde 2**, dan ⛔ **belum bertambah karena bahannya belum
ada** — bukan karena bersih.

---

## Bab D — PERTANYAAN SIAP KIRIM

⛔ **Hanya yang BARU.** Lanjut dari **P29**. Penjawabnya **tidak membaca XML**.
⛔ **Nol saya jawab sendiri.** Urut dari yang paling memblokir.

---

### Untuk DBA

#### **P29 — Sebuah kolom yang namanya berarti "jenis usaha" ternyata menyimpan seluruh data kontrak. Benarkah?**

Ada sebuah kolom yang **namanya menyebut jenis usaha**, tetapi sistem lama **tidak
memperlakukannya sebagai jenis usaha**. Ia membaca isi kolom itu sebagai **satu dokumen teks
berisi seluruh data kontrak**, lalu membongkarnya menjadi puluhan medan. Ini dilakukan di
**enam tempat berbeda**. Dan kolom yang sama **juga muncul sebagai kolom biasa** di daftar laporan.

**Konteks:** karena isi kolom itulah yang menentukan medan apa saja yang ada, kami **tidak dapat
menyatakan bentuk data kontrak** tanpa melihat isinya.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh contoh** — **tiga sampai lima isi kolom itu apa
adanya** dari basis data. ⭐ Dan satu keterangan: apakah kolom itu **memang dipakai dua cara**
*(kadang jenis usaha, kadang dokumen lengkap)*, atau kami salah membaca.

**Dampak bila salah:** ⛔⛔ seluruh data kontrak lama **tidak terbaca** oleh sistem baru, atau
kolom itu disalin sebagai teks pendek dan isinya **terpotong tanpa peringatan**.

rujukan: `getString("CLASSOFBUSINESS")` + `adoptJSONObject(...)` di 6 langkah Java
*(`FetchMasterTreatyIn`, `SetTreatyIn_Act`, `InputPolicyTreatyInDetail_*`,
`InputPolicyTreatyOutDetail_*`)*; `.CLASSOFBUSINESS` juga salah satu dari 66 kolom
`BrowseTreatyInDetail` — OQ-012, ronde 3 baru #1

---

### Untuk pemilik export Pega

#### **P30 — Dua langkah lagi yang dipanggil tetapi tidak ikut terkirim, dan keduanya di jalur penting**

Selain langkah yang sudah kami tanyakan, ada **dua langkah lagi** yang dipanggil sistem lama
tetapi **tidak ada di dalam berkas yang kami terima**. Satu dipanggil oleh langkah yang
**menyimpan data kontrak**; satu lagi dipanggil oleh langkah yang **menjumlahkan nilai
pertanggungan dan premi**.

**Konteks:** keduanya di jalur yang paling penting — menyimpan, dan menghitung uang.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — kedua langkah itu. Bila ekspornya memang
tidak menyertakan langkah dari bagian lain sistem, sebutkan **bagian mana** keduanya tinggal,
supaya kami dapat memintanya terpisah.

**Dampak bila salah:** ⛔ dua langkah yang selama ini berjalan **tidak dibangun**, dan tidak ada
yang tahu ia pernah ada — sebab yang hilang **tidak meninggalkan jejak**.

rujukan: `CopyToPolicy` dipanggil `SaveJsonPolisTreatyIn_Act`;
`SumTSIPremiSpreadRNMMultiCob_Act` dipanggil `SumTSIPremiSpreadedRNM_Act`; keduanya di luar
66 target yang ada — ronde 3 baru #2

---

### Untuk pengembang Pega lama

#### **P31 — Ketika data kontrak gagal dibaca, sistem hanya mencatat di log lalu LANJUT. Disengaja?**

Di enam tempat, sistem lama membongkar data kontrak dari satu kolom. Bila pembongkaran itu
**gagal** — misalnya isinya rusak — ⛔ **sistem hanya menulis catatan di log, lalu melanjutkan
langkah berikutnya.** Tidak ada pesan kepada pengguna, dan tidak ada penghentian.

**Konteks:** artinya perhitungan berikutnya dapat berjalan di atas data yang **kosong atau
separuh terisi**, dan hasilnya tampak normal.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** disengaja, kegagalan memang boleh
diabaikan · **(b)** seharusnya berhenti dan memberi pesan · **(c)** tidak tahu.
⭐ Dan: **seberapa sering catatan galat itu muncul** di log sistem lama, kalau bisa dilihat.

**Dampak bila salah:** ⛔⛔ angka yang salah **terlihat benar**, dan tidak ada satu pun tanda bahwa
datanya tidak lengkap.

rujukan: `catch(Exception e){ oLog.error("ReloadSection:Expection : "+...) }` — tanpa
penghentian; ⚠️ pesan log menyebut konteks lain *(section reload)* — ronde 3 baru #3

---

#### **P32 — Dua cara penulisan tanggal dipakai bersamaan. Mana yang benar di mana?**

Sistem lama menulis tanggal dengan **dua susunan berbeda** dalam modul yang sama: satu
**hari-bulan-tahun**, satu lagi **bulan-hari-tahun berikut jam**. Keduanya dipakai untuk medan
yang berbeda, dalam aturan yang berdekatan.

**Konteks:** tanggal 3 Februari dan 2 Maret **tidak dapat dibedakan** bila susunannya tertukar,
dan kekeliruan seperti ini **tidak menimbulkan pesan galat**.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh daftar** — medan mana memakai susunan mana.
⭐ Atau, bila lebih mudah: **contoh nilai nyata** dari tiga medan tanggal di basis data, supaya
kami dapat menurunkan susunannya sendiri.

**Dampak bila salah:** ⛔⛔ tanggal mulai dan berakhirnya kontrak **bergeser sampai sebelas bulan**,
dan tidak ada yang melihatnya sampai ada yang menagih.

rujukan: `@DateTime.CurrentDate("dd/MM/yyyy","")` lawan
`@DateTime.CurrentDate("MM/dd/yyyy hh:mm a","")` di `InputPolicyTreatyIn_preDT` dan
`DeptHeadTreatyInUW_preDT` — ronde 3 baru #5

---

#### **P33 — Satu medan "nama operator" diisi dari dua sumber berbeda. Mana yang dimaksud?**

Medan yang mencatat **nama operator** diisi dari **dua tempat berbeda** di dua aturan yang
berdekatan: satu memakai **pengenal akun**, satu memakai **nama tampilan orangnya**.

**Konteks:** keduanya terlihat serupa di layar tetapi berbeda isinya, dan medan ini masuk ke
jejak audit.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** seharusnya pengenal akun ·
**(b)** seharusnya nama tampilan · **(c)** memang berbeda per tahap, jelaskan.

**Dampak bila salah:** ⚠️ jejak audit memuat **dua jenis pengenal di satu kolom**, dan mencari
riwayat seseorang menjadi tidak dapat diandalkan.

rujukan: `DeptHeadTreatyInUW_preDT` lgk 2 `= OperatorID.pyUserIdentifier` lawan
`InputPolicyTreatyIn_preDT` lgk 10 `= OperatorID.pyUserName` — ronde 3 baru #8

---

#### **P34 — Tujuh belas langkah MENGHAPUS pesan kesalahan. Kenapa?**

Sistem lama memasang pesan kesalahan di **122 tempat**, dan ⛔ **menghapus pesan kesalahan di
17 tempat.**

**Konteks:** menghapus pesan kesalahan berarti melanjutkan sesuatu yang sesaat sebelumnya
ditandai salah. Itu bisa sah — misalnya membersihkan pesan lama sebelum memeriksa ulang — atau
bisa berarti validasi dilewati.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** pembersihan sebelum pemeriksaan ulang,
wajar · **(b)** memang untuk melewati validasi tertentu, sebutkan mana · **(c)** campuran, perlu
ditinjau.

**Dampak bila salah:** ⚠️ sistem baru **menolak** data yang selama ini diterima, atau **menerima**
data yang seharusnya ditolak.

rujukan: `Page-Set-Messages` 105 + `Property-Set-Messages` 17 lawan `Page-Clear-Messages` 17
di 92 `Activity` — ronde 3 baru #12

---

### Untuk Product + Underwriting

#### **P35 — Bila tanggal akhir kontrak dibiarkan kosong, sistem mengisinya dengan HARI INI. Benarkah begitu?**

Ketika sebuah kontrak dibuka dan **tanggal akhirnya belum diisi**, sistem lama mengisinya dengan
**tanggal hari itu** — ⭐ **sama dengan tanggal mulainya**. Sementara di bagian lain sistem ada
perhitungan yang menetapkan masa berlaku **satu tahun**.

**Konteks:** kalau yang berlaku adalah yang pertama, maka kontrak yang tanggal akhirnya terlupakan
**tersimpan sebagai kontrak yang berakhir di hari ia dibuat.**

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang harus satu tahun, pengisian
hari-ini adalah kekeliruan · **(b)** hari-ini benar, dan pengguna wajib mengubahnya ·
**(c)** tergantung jenis kontrak, jelaskan mana.

**Dampak bila salah:** ⛔⛔ **masa berlaku kontrak salah**, dan itu menentukan premi, klaim yang
ditanggung, dan pelaporan.

rujukan: `InputPolicyTreatyIn_preDT` lgk 2 — gerbang `.EndDate==""` ⇒
`@DateTime.CurrentDate("dd/MM/yyyy","")`; bandingkan `SystemSetOneYear_DT` ⇒
`@DateTime.addCalendar(.StartDate,1,0,0,0,0,0,0)` — ronde 3 baru #4

---

#### **P36 — Ada DUA penanda persetujuan, bukan satu. Apa bedanya?**

Kami menemukan **penanda persetujuan kedua** yang belum pernah kami laporkan — namanya menyebut
**persetujuan kepada Kepala Departemen**, terpisah dari penanda persetujuan umum. Keduanya
**dikosongkan bersamaan** ketika kasus masuk tahap Kepala Departemen.

**Konteks:** dua penanda berarti dua keadaan yang dapat berbeda, dan kami tidak dapat menurunkan
bedanya dari aturan.

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** dua tahap persetujuan yang
berbeda, jelaskan mana milik siapa · **(b)** salah satunya sudah tidak dipakai ·
**(c)** satu turunan dari yang lain.

**Dampak bila salah:** ⛔ persetujuan tercatat di penanda yang salah, dan kasus **maju atau tertahan
keliru**.

rujukan: `DeptHeadTreatyInUW_preDT` lgk 4–5 — `.PolicyTreatyIn.IsApproved = ""` dan
`.PolicyTreatyIn.isApprovedtoDeptHead = ""` — ronde 3 baru #6

---

#### **P37 — Dua kode angka menandai kontrak "punya penempatan keluar". Apa artinya?**

Sebuah aturan memeriksa **dua kode angka jenis kontrak** dan, bila cocok, menandai kontrak itu
**punya penempatan keluar**. Kedua kode itu **hanya angka**, tanpa keterangan.

**Konteks:** penandaan ini menentukan apakah sebuah kontrak masuk jalur penempatan keluar,
dan itu jalur yang berbeda sepenuhnya.

**Bentuk jawaban yang diharapkan:** dua baris — **apa arti masing-masing kode**, dan
⭐ **apakah masih ada kode lain** yang seharusnya ikut ditandai tetapi belum.

**Dampak bila salah:** ⛔ kontrak yang seharusnya masuk jalur penempatan keluar **tidak masuk**,
atau sebaliknya.

rujukan: `TestTreatyToFacStatus` — `.TreatyType=="10015"` dan `.TreatyType=="10218"` ⇒
`Primary.PolicyTreatyIn.HasFacOut = "1"` — ronde 3 baru #10

---

#### **P38 — Sebuah aturan lini jiwa hidup di dalam modul kontrak reasuransi masuk. Disengaja?**

Ada satu aturan yang **tergolong lini asuransi jiwa** tetapi tinggal di dalam modul kontrak
reasuransi masuk, dan ia menulis penanda **"jiwa"** ke catatan lampiran **baris pertama**.

**Konteks:** kalau modul ini tidak menangani lini jiwa, aturan itu tidak seharusnya ada; kalau
menangani, batas antar-modul perlu kami perbaiki.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** modul ini memang juga menangani lini
jiwa, jelaskan kapan · **(b)** peninggalan, tidak dipakai · **(c)** tidak tahu.

**Dampak bila salah:** ⚠️ lampiran tergolong ke lini yang salah, dan penggolongan itu ikut
terbawa ke laporan.

rujukan: `setCategoryAttachment_DT`, class `ASM-FW-GISFW-Work-LIFE`, ⇒
`Attachment.pxResults(1).NOTE = "LIFE"` — ronde 3 baru #13

---

#### **P39 — Satu mata uang dikecualikan secara khusus dari dua daftar pilihan. Masih perlu?**

Dua daftar pilihan mata uang **mengecualikan satu kode mata uang tertentu** secara khusus, ditulis
langsung di dalam aturan.

**Konteks:** pengecualian yang ditulis langsung tidak dapat diubah tanpa mengubah aturan, dan kami
perlu tahu apakah ia masih dimaksudkan.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** masih perlu dikecualikan, jelaskan
alasannya · **(b)** sudah tidak perlu · **(c)** seharusnya ada mata uang lain yang juga
dikecualikan, sebutkan.

**Dampak bila salah:** ⚠️ mata uang yang masih dipakai **hilang dari pilihan**, atau yang sudah
mati **muncul kembali**.

rujukan: penyaring bernilai `"ITL"` di `BrowseCurrency_RD` dan `BrowseCurrencyTreatyIn_RD` —
ronde 3 baru #15

---

### Untuk IAM

#### **P40 — Nama seseorang tertulis di dalam teks pemberitahuan yang dibaca pengguna**

Selain nama orang yang dipakai untuk menentukan wewenang, kami menemukan **satu nama orang
tertulis di dalam isi pesan pemberitahuan** — pesan yang memberi tahu pengguna di kotak masuk
siapa sebuah pengajuan sedang berada. ⛔ **Nilainya tidak kami salin ke berkas mana pun.**
Pesan-pesan lain di tempat yang sama **mengambil namanya dari data**, bukan menuliskannya.

**Konteks:** ini berbeda dari temuan sebelumnya — di sini nama itu **bukan penjaga akses**,
melainkan **isi pesan**. Jadi bila orangnya berganti, akses tetap jalan, ⛔ **tetapi pesannya
menjadi salah.**

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** seharusnya diambil dari data seperti
pesan lainnya · **(b)** memang harus nama itu, jelaskan alasannya · **(c)** pesan ini sudah tidak
dipakai.

**Dampak bila salah:** ⛔ pengguna diberi tahu bahwa pengajuannya ada di kotak masuk **orang yang
salah**, dan mereka menghubungi orang yang salah.

rujukan: `DeptHeadTreatyIn_UW_postDT`, cabang `PositionNote=="ReasTreatyInGroupLeader"`, `SET`
ke `.NBStatus` — 1 kemunculan; bandingkan 3 cabang lain yang memakai
`@toUpperCase(pyWorkPage.pxCreateOpName)` — ronde 3 baru #7

---

## Bab E — Lampiran bukti isolasi

| Yang dibuktikan | Keadaan | Perintah audit |
| --- | --- | --- |
| korpus `NB Treaty In` tidak berubah | ✅ **`62a3735ebb2eabb788c8cd8abdda78ca`** — cocok ronde 1 dan 2 | `find . -name '*.xml' \| sort \| xargs cat \| md5sum` |
| ⛔ **`EDM Treaty In` NOL dibuka** | ✅ **nol** — ⭐ OQ-025 tetap pertanyaan, varian **tidak dipinjam** | folder itu tidak pernah menjadi argumen; seluruh skrip ber-`os.chdir(r"D:\XML\RNM_BRD\NB Treaty In")` |
| ⛔ folder korpus modul lain | ✅ **NOL dibuka** | sama — satu direktori kerja di setiap skrip |
| ⛔ `D:\XML\nusantara-re\` | ✅ **NOL disentuh** | tidak pernah muncul sebagai argumen |
| ⛔ ronde 1 dan ronde 2 **tidak disunting** | ✅ **nol suntingan** — ⭐ ralat ditulis **di berkas ini** | `ls -la .scratch\nb-treaty-in\` — mtime ronde 1 **09:25**, ronde 2 **09:51** |
| berkas dibuat | ✅ **tepat SATU** — `grilling-ronde-3.md` | `ls .scratch\nb-treaty-in\` |
| ⛔ spec / tiket | ✅ **NOL** | tidak ada berkas ber-nama `spec`/`issue`/`tiket` |
| kode · DDL · `CREATE TABLE` · daftar kolom usulan | ✅ **NOL** | ⭐ kolom yang disebut Bab B adalah **kolom terbaca**, bukan usulan |
| ⛔ nilai nama orang | ✅ **NOL disalin** | §1.3 dan P40 menyebut **rule, cabang, medan, dan jumlah**; ⛔ **nol nama** |
| butir terbuka yang saya tutup | ✅ ⛔ **NOL** | Bab C — setiap butir tetap berpemilik |
| nomor OQ dikarang | ✅ **NOL** | dirujuk **"ronde 3 baru #n"**; OQ-073 disebut, **tidak dipakai** |
| ⚠️ ejaan pemilik seragam | ✅ **`[Product+Underwriting]` saja** | diseragamkan **sejak awal**, dicek dengan cacahan sebelum merekap |
| ⛔ kode Java **tidak ditulis ulang sebagai Go** | ✅ **nol** | ⭐ dikutip **dua baris** untuk membuktikan klaim, sesuai izin brief |

### ⚠️ Baris yang TIDAK bersih — dilaporkan sesuai perintah brief

⚠️ Sesi ini **panjang dan pernah disambung**. Pada penyambungan, harness **menyuntikkan ulang dua
bacaan lama** dari `.scratch` konteks lain. ⛔ **Bukan saya yang memanggilnya**, dan ⛔ **nol saya
panggil di ronde 3.**

⭐ **Apa yang tidak merembes:** ⛔ nol nama tabel, nol nama kolom, nol pola rancangan dari konteks
lain. ⭐ Seluruh angka ronde ini dari **sisiran korpus `NB Treaty In` sendiri** dan dari berkas
yang brief perintahkan dibaca.

---

## Bab F — ⭐⭐ PEMULIHAN RANTAI NILAI

⭐ **Satu pertanyaan, dijawab dengan angka:** dari seluruh penetapan nilai di modul ini, **berapa
persen dapat dipulihkan**, dan **berapa persen tetap hilang**?

### F.1 Hitungan — DUA CARA

⭐ **Cara A — menurut JUMLAH PENETAPAN.** Jendela: seluruh 278 berkas.

| Sumber penetapan | Jumlah | Terbaca? |
| --- | ---: | :---: |
| `Property-Set` di 92 `Activity` | ⛔ **942** | ⛔ **TIDAK** — nilainya nol terekspor |
| ⭐ `SET` di 12 `DataTransform` | ⭐ **53** | ✅ **YA, seluruhnya** |
| ⭐ langkah `Java` yang menetapkan | ⭐ **8** | ✅ **YA** — 6 adopsi JSON + 2 pembuang duplikat |
| ⭐ **TOTAL** | **1.003** | — |

> ⭐ **Dapat dipulihkan: 61 dari 1.003 = 6,1 %**
> ⛔ **Tetap hilang: 942 dari 1.003 = 93,9 %**

⭐ **Cara B — menurut JUMLAH ATURAN yang menetapkan nilai.** Jendela: berkas yang memuat langkah
penetapan.

| Golongan | Berkas | Terbaca? |
| --- | ---: | :---: |
| `Activity` ber-`Property-Set` | ⛔ **92** | ⛔ **TIDAK** |
| ⭐ `DataTransform` | ⭐ **12** | ✅ **YA** |
| ⭐ `Activity` ber-langkah `Java` | ⭐ **7** | ⚠️ **SEBAGIAN** — hanya langkah Java-nya |
| **TOTAL berkas berbeda** | **99** *(7 tumpang-tindih dengan 92)* | — |

> ⭐ **Dapat dipulihkan: 12 dari 99 aturan penetapan = 12,1 %**
> ⛔ **Tidak dapat: 87 dari 99 = 87,9 %**

⚠️ **Kedua cara BERSELISIH — 6,1 % lawan 12,1 % — dan selisihnya WAJAR, bukan cacat.**
⭐ Sebabnya: `DataTransform` adalah **berkas kecil ber-penetapan sedikit** *(rata-rata 4,4
`SET` per berkas)*, sedangkan `Activity` adalah **berkas besar ber-penetapan banyak**
*(rata-rata 10,2 `Property-Set` per berkas)*. ⭐ Jadi menurut **berkas** pemulihannya terlihat
dua kali lebih baik daripada menurut **penetapan**.

⭐ **Angka yang harus dipakai untuk keputusan adalah CARA A**, sebab yang dibangun ulang adalah
**penetapan**, bukan berkas.

> ⛔⛔ **VONIS ANGKA: sekitar SATU DARI ENAM BELAS penetapan nilai dapat dipulihkan.
> Lima belas dari enam belas tidak.**

### F.2 ⭐ Tetapi yang 6,1 % itu TIDAK merata — dan justru bagian terpenting

⭐ **Angka mentahnya menyesatkan bila berhenti di situ.** ⭐ Yang dapat dipulihkan **bukan potongan
acak** — ia **jenis tertentu**, dan jenisnya penting:

| Yang DAPAT dipulihkan | Isinya |
| --- | --- |
| ⭐⭐ **penyiapan tahap** | pengosongan bendera, penetapan tanggal awal, penetapan **nama antrean**, salinan `Quotation` → `QuotationData` |
| ⭐ **pengalihan cabang** | `IsNewPolicyNonProp`, `HasFacOut`, `btnQuotation` |
| ⭐ **kronologi** | struktur `SuggestList` berikut empat medannya |
| ⭐ **pesan status** | keenam pemetaan antrean → teks pemberitahuan |
| ⭐ **masa berlaku** | rumus satu tahun, **utuh** |
| ⭐ **pembacaan data kontrak** | mekanisme adopsi JSON |

| Yang TETAP hilang | Isinya |
| --- | --- |
| ⛔⛔ **seluruh aritmetika uang** | premi, bagian, pengurangan, pajak, komisi, limit, pemulihan |
| ⛔ **seluruh penetapan di dalam 92 `Activity`** | termasuk 30 aturan hitung uang ronde 2 |

⭐⭐ **Jadi jawabannya untuk keputusan `to-spec`:**

> ⭐ **Spec ALUR, TAHAP, WEWENANG, dan KRONOLOGI dapat ditulis SEKARANG** — bahannya ada, dan
> ronde 3 memulihkan bagian yang menentukan.
> ⛔⛔ **Spec PERHITUNGAN UANG tidak dapat ditulis sama sekali** — bukan sebagian, **nol**.
> ⭐ 6,1 % yang dipulihkan **tidak memuat satu pun rumus uang.**

⚠️ **Dan satu peringatan atas angka 93,9 %:** ⛔ ia **bukan** ukuran seberapa banyak yang belum
dibaca. ⭐ Ia ukuran seberapa banyak yang **tidak ada di ekspor**. ⛔ **Ronde keempat, kelima, atau
kedua puluh tidak akan menggerakkannya satu persen pun.**

---

## Penutup

### Cakupan sesudah ronde 3

| Golongan | Berkas | Keadaan |
| --- | ---: | --- |
| `RDBList` | 41 | ✅ **selesai** *(R1)* |
| blok PL/SQL | 5 | ✅ **dibaca utuh** *(R1)* |
| `Section` + `Harness` | 31 | ⚠️ **disisir** *(R2)*, ⛔ **nol dibaca utuh** |
| `When` | 75 | ✅ **syaratnya dibaca** *(R2)* |
| `DecisionTable` | 2 | ✅ **dibaca utuh** *(R2)* — ⛔ isinya tidak terekspor |
| ⭐ `DataTransform` | **12** | ⭐ **DIBACA UTUH, 80 dari 80 baris** *(R3)* |
| ⭐ langkah `Java` | **8** | ⭐ **DIBACA UTUH** *(R3)* |
| ⭐ `FlowAction` | **9** | ⭐ **disisir** *(R3)* — ⛔ kesembilan kosong |
| ⭐ `ReportDefinition` | **15** | ⭐ **disisir** *(R3)* |
| `Flow` | 1 | ✅ **dikonfirmasi**, tidak ditelusur ulang |
| ⚠️ `Activity` | 92 | ⚠️ **bentuk & rantai panggil terbaca** *(R3)*; ⛔ **nilai penetapan tidak ada** |
| `.xlsx` | 2 | ⛔ **tidak dibuka, sengaja** |

⭐ **Dibaca UTUH sepanjang tiga ronde: 28 berkas** — 5 blok PL/SQL + 2 `DecisionTable` +
1 `When` + 12 `DataTransform` + 8 langkah Java *(di 7 berkas)*.
⭐ **MASIH nol dibuka: 2 berkas** — ⭐ keduanya `.xlsx`, **artefak turunan, sengaja.**

> ⭐⭐ **Seluruh 278 berkas `.xml` kini sudah tersentuh sedikitnya satu kali.**

⛔ **Yang tetap tidak dibaca utuh: 123 berkas** — 31 layar, 92 `Activity` — ⭐ dan untuk 92
`Activity` itu **membaca utuh tidak akan menambah apa pun**, sebab nilainya tidak ada di sana.

### Rekap butir

| | Jumlah |
| --- | ---: |
| butir **lahir** ronde 3 | ⭐ **16** |
| butir lama **berubah status** | **10** |
| ⛔ butir **ditutup** | ⛔ **0** |
| butir kumulatif tiga ronde | **17 terdaftar + 7 (R1) + 13 (R2) + 16 (R3)** |
| pertanyaan siap kirim kumulatif | ⭐ **40** — P1–P40 |

### Apakah frontier masih terbuka?

⭐⭐ **Frontier BACAAN sudah HABIS. Frontier PENGETAHUAN tidak.**

⭐ **Itu dua hal berbeda, dan ronde 3 memisahkannya dengan bersih:**

| | Keadaan |
| --- | --- |
| ⭐ **Yang dapat dibaca dari korpus** | ✅ **HABIS.** 278 dari 278 berkas tersentuh; 28 dibaca utuh; nol golongan tak tersentuh |
| ⛔ **Yang tidak ada di korpus** | ⛔ **TETAP TIDAK ADA.** 942 nilai penetapan · baris 2 tabel keputusan · badan 3 procedure · 3 aturan bisnis yang dipanggil tetapi hilang · struktur JSON |

### ⭐⭐ VONIS — apakah masih perlu ronde 4?

> ⛔⛔ **TIDAK. Ronde 4 atas korpus ini tidak akan menghasilkan apa pun yang berarti.**

⭐ **Alasannya berupa angka, bukan perasaan:**

1. ⭐ **278 dari 278 berkas sudah tersentuh.** ⛔ Tidak ada golongan yang belum dibuka.
2. ⛔ **93,9 % rantai nilai tidak ada di ekspor** — ⭐ dan **nol persen** dari itu dapat digerakkan
   dengan membaca lebih lama.
3. ⛔ **Tiga aturan bisnis yang dipanggil tidak ada di modul** — ⭐ membacanya **mustahil**, bukan sulit.
4. ⛔ **Baris kedua tabel keputusan tidak terekspor** — sama.

⭐ **Yang masih dapat menambah pengetahuan, dan hanya ini:**

| Pekerjaan | Sifat | Menunggu |
| --- | --- | --- |
| ⭐ **Jawaban atas P1–P40** | ⭐ **satu-satunya jalur yang menggerakkan angka** | ⛔ **manusia** |
| ⭐ **Ekspor ulang dengan isi langkah disertakan** *(P18)* | ⭐ menggerakkan 93,9 % itu | ⛔ **pemilik export** |
| ⭐ **Spec ALUR + TAHAP + KRONOLOGI** | ⭐ **dapat dimulai SEKARANG** — Bab F.2 | ✅ **tidak menunggu** |
| ⛔ Spec **perhitungan uang** | ⛔ **nol persen mungkin** | ⛔ P18 |
| ⚠️ Membuka **`EDM Treaty In`** | ⚠️ ⛔ **keputusan work owner** | ⛔ izin |

### ⭐ Kesiapan `to-spec` — tiga syarat, dan satu jalan yang terbuka

⭐ Ronde 2 menyatakan **tiga syarat**. ⭐ **Ronde 3 tidak menambah syarat**, ⛔ tetapi
**memperluas syarat ketiga** dan ⭐ **membuka satu jalan yang sebelumnya tertutup:**

| Syarat | Keadaan sesudah ronde 3 |
| --- | --- |
| ⛔ **OQ-025** — efek keluar terakhir | ⛔ **diperluas menjadi TIGA aturan hilang**, bukan satu |
| ⛔ **Naskah tiga stored procedure** | ⛔ tidak berubah |
| ⛔ **Isi 942 langkah penetapan** *(P18)* | ⛔ **dikuantifikasi: 93,9 %** — dan ⛔ **nol dapat dipulihkan dengan membaca** |

⭐⭐ **Dan jalan yang terbuka, dinyatakan terang-terangan:**

> ⭐ **`to-spec` untuk ALUR, TAHAP, WEWENANG, KRONOLOGI, dan MASA BERLAKU dapat dimulai sekarang**
> — ⭐ bahannya lengkap, dan ronde 3-lah yang melengkapinya.
> ⛔⛔ **`to-spec` untuk PERHITUNGAN UANG tidak dapat dimulai sama sekali**, dan menunda seluruh
> spec sampai P18 terjawab akan **menahan pekerjaan yang sudah siap.**

⛔ **Pemisahan itu keputusan work owner, bukan saya.** ⭐ Saya hanya menyatakan **bahwa
pemisahannya mungkin**, dan **di mana garisnya**.
