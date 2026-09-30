# NB Treaty In — Grilling Ronde 4
## Sembilan puluh satu elemen mati, dan delapan puluh di antaranya tidak menyembunyikan apa pun

Modul: **hanya `NB Treaty In`**. Korpus `D:\XML\RNM_BRD\` **READ-ONLY**; menulis hanya ke
`OUTPUT_HASIL_RNM\`. ⛔ `D:\XML\nusantara-re\` **tidak disentuh**.
⭐ Modul lain dibaca **hanya** untuk membuktikan ada/tidaknya sesuatu di luar NB Treaty In —
⛔ temuannya **tidak masuk sensus modul ini**, dan selalu disebut jendelanya.

⛔ **Nol kode · nol DDL · nol `CREATE TABLE` · nol usulan daftar kolom · nol butir ditutup sendiri ·
nol nilai nama orang disalin.**

⭐ **Pertanyaan baru TIDAK ada di berkas ini** — ia di `PERTANYAAN-RONDE-4.md`, bernomor mulai
**P42**.

---

## §0 — Pemeriksaan atas brief, dan dua ralat

### 0.1 ⭐ Tiga angka brief diperiksa — ketiganya BENAR, dan saya sempat salah menyangkanya keliru

| Klaim brief | Terukur | |
| --- | --- | :---: |
| `Property-Set` **25.986** | ⭐ **25.986** — jendela: **seluruh korpus, 21 folder** | ✅ |
| `pyBrowseSQL` **41** di NB · **1.151** di korpus | ⭐ **41** · **1.151** | ✅ |
| `Section`+`Harness` **16.635.834 B** dari **35.678.284 B** | ⭐ keduanya **tepat** | ✅ |

⚠️ **Saya sempat menyangka angka 25.986 bertabrakan dengan 942 milik ronde 3.** ⛔ **Tidak.**
⭐ **942 adalah NB Treaty In; 25.986 adalah seluruh korpus.** ⭐ Keduanya benar — **jendelanya
berbeda**, dan itu persis jebakan sensus `CLAUDE.md` §4a **#5** *(entri dijumlahkan lintas wadah)*.
⭐ Perintah audit: cacah `<pyStepsActivityName>Property-Set</...>` per folder, lalu dijumlahkan.

### 0.2 ⛔ RALAT ATAS BRIEF — dua pernyataannya tentang ronde 1–3 tidak tepat

> ⛔ **Brief menulis:** *"Sasaran 1 — LAPISAN LAYAR, BELUM DIBUKA SAMA SEKALI … **Nol berkas
> dibuka** di ronde 1–3."*
>
> ⚠️ **Tidak tepat.** ⭐ **Ronde 2 §2 menyisir ke-31 berkas layar** dan menghasilkan
> **733 syarat berisi, 229 bermakna, 11 berkas murni tata letak**, berikut vonis per-berkas.
> ⛔ **Yang BENAR:** ronde 2 **tidak membaca satu pun secara utuh** — ia hanya **mengekstrak nilai
> tag**, ⛔ tidak pernah menelusuri strukturnya. ⭐ **Sasaran 1 tetap sah dan tetap baru**, sebab
> yang diminta ronde 4 — *elemen ini menyembunyikan apa* — **memang belum pernah dijawab.**

> ⛔ **Brief menulis:** *"Ronde sebelumnya pernah menyimpulkan naskah SQL tidak ikut terekspor.
> **Itu keliru.**"*
>
> ⛔ **Ronde 1–3 tidak pernah menyimpulkan itu.** ⭐ Ronde 1 **mengurai ke-41 `pyBrowseSQL`** dan
> **mengutip utuh kelima blok PL/SQL**, termasuk ketiga pemanggilan procedure.
> ⭐ **Yang ronde 2 simpulkan adalah hal yang BERBEDA:** ⛔ **nilai langkah `Property-Set` tidak
> terekspor** — bukan naskah SQL. ⭐ Keduanya tidak boleh dicampur.

### 0.3 ⭐⭐ RALAT ATAS RONDE 2 — dan brief-lah yang memancingnya, walau tag yang ia sebut juga meleset

> ⛔⛔ **RALAT.** Kalimat ronde 2 §3.1 berikut **terlalu luas**, dan dikutip utuh:
>
> *"⛔ **Nol `pyCriteriaValue`, nol `pyResult`, nol `pyReturnValue`, nol `pyPropertyName`, nol
> `pyOtherwiseResult`.** ⭐ Yang bertahan hanya `pyLabel` berisi nama aturan dan daftar versi."*
>
> ⭐ **Bagian yang benar:** kelima tag itu memang **nol**, dan **nilai serta hasil baris keputusan
> memang tidak terekspor**.
> ⛔ **Bagian yang keliru:** ronde 2 menyimpulkan *"yang bertahan hanya `pyLabel`"*. ⚠️ **Itu
> salah** — ⭐ **tiga hal lain bertahan**, dan ketiganya berarti.

⚠️ **Dan tag yang brief sebut juga bukan tempatnya.** `[terverifikasi]` Diperiksa langsung:

| Tag | `isApproved` | `BusinessType_DeT` |
| --- | ---: | ---: |
| `pyCondition` · `pyOrConditions` · `pyResults` *(disebut brief)* | ⛔ **0 berisi** | ⛔ **0 berisi** |
| ⭐ **`pyProperty`** | ⭐ **berisi** | ⭐ **berisi** |
| ⭐ **`pyColumnDataType`** · **`pyDefaultOperator`** | ⭐ berisi | ⭐ berisi |
| ⭐ **`pyRowNum`** | **1 baris** | ⭐ **9 baris** |

⭐⭐ **Yang kini terbaca, dan ronde 2 melewatkannya:**

| Tabel keputusan | Properti yang diuji | Tipe kolom | Operator | Baris |
| --- | --- | --- | --- | ---: |
| ⭐ `isApproved` | ⭐ **`pyWorkPage.PolicyTreatyIn.IsApproved`** | ⭐⭐ **`text`** | `=` | **1** |
| ⭐ `BusinessType_DeT` | ⭐ **`.Quotation.GroupPanel`** dan **`.Quotation.BusinessOldId`** | `text` | `=` | ⭐ **9** |

⭐⭐ **Dan `pyColumnDataType = text` untuk `IsApproved` menjawab butir ronde 2 baru #5 dari korpus.**
⭐ Ronde 2 menemukan bendera itu dibandingkan **sebagai teks berkutip** di satu tempat dan
**sebagai angka telanjang** di tempat lain, dan menyebutnya `[terbuka]`.
⭐ **Kini terbaca: tabel keputusan memperlakukannya sebagai TEKS.** ⇒ ⭐ ejaan berkutip `'0'` yang
sejalan dengan aturan hidup; ⛔ pembandingan angka telanjang-lah yang menyimpang.
⭐ `[keputusan work owner]` Brief menyatakan aturan hidup ada di `DecisionTable`, dan
**`0` = ditolak, selain itu = disetujui** — ⭐ **sejalan dengan tipe `text`.**

> ⛔⛔ **RALAT — dan yang menangkapnya BUKAN saya.** Kalimat berikut dikutip utuh, lalu ditarik
> sebagian:
>
> *"⚠️⚠️ Satu temuan baru yang tidak diminta siapa pun: nomor baris `BusinessType_DeT` adalah
> −1, 2, 3, 5, 7, 22, 29, 30, 34 — ⛔ tidak berurutan, dan tertinggi 34 padahal hanya 9 baris
> tersisa. ⭐ `[dugaan]` tabel ini **pernah punya sedikitnya 34 baris** dan sebagian besar
> **dihapus**, sementara penomorannya menyimpan lubangnya."*
>
> ⭐ **Pertanyaan P48 yang lahir dari kalimat ini SUDAH DITARIK** oleh pemeriksa, dengan alasan
> `[terverifikasi]`: *"Tidak ada baris yang hilang. Tabelnya utuh, 36 baris, seluruh nilainya
> terbaca."*
>
> ⛔ **Yang saya TARIK: dugaannya.** ⭐ Kesimpulan *"sebagian besar baris dihapus"* **tidak
> berdasar** — ia dugaan yang saya bangun dari lubang penomoran, dan lubang penomoran **bukan
> bukti penghapusan**.
>
> ⚠️ **Yang saya TIDAK tarik, sebab diukur ulang dan hasilnya tetap:** di **berkas ekspor ini**,
> `pyRowNum` muncul **9 kali** dengan nilai −1, 2, 3, 5, 7, 22, 29, 30, 34 — dicacah **dua cara**
> *(tag pembuka; tag berpasangan)*, **keduanya 9**. ⛔ Dan **nol tag mana pun** di berkas itu memuat
> nilai kriteria atau hasil: `pyValue` seluruhnya `3`, dan tidak ada teks bisnis di tag mana pun.
>
> ⚠️⚠️ **Jadi dua pernyataan ini berselisih, dan saya cantumkan keduanya alih-alih memilih:**
>
> | | Baris | Nilai kriteria |
> | --- | ---: | --- |
> | ⭐ catatan pemeriksa `[terverifikasi]` | **36** | **terbaca seluruhnya** |
> | ⭐ ukuran saya atas **berkas ekspor** | **9** penanda | ⛔ **nol** |
>
> ⭐ **Penjelasan yang paling mungkin, dan ia `[terbuka]`:** kedua angka benar untuk **jendela yang
> berbeda** — **36 baris hidup di sistem Pega**, sementara **ekspor ini hanya membawa 9 penanda
> baris dan nol nilainya**. ⛔ **Belum dipastikan**, dan pemastiannya bukan milik saya.
>
> ⭐⭐ **Yang wajib dicatat tentang cara kerja, bukan tentang tabelnya:** ⛔ **ini kegagalan
> instrumen kesembilan sepanjang empat ronde, dan satu-satunya yang TIDAK saya tangkap sendiri.**
> ⭐ Delapan sebelumnya tertangkap oleh ujian instrumen; ⛔ **yang ini lolos sampai ke berkas
> pertanyaan**, dan ditangkap pemeriksa manusia. ⚠️ Sebabnya: saya menarik **kesimpulan sejarah**
> *("dihapus")* dari **bukti struktur** *(lubang penomoran)*, dan keduanya tidak sejenis.
>
> ⚠️ **Jebakan #2 juga terlewat di berkas ini:** `<rowdata />` yang menutup sendiri berjumlah
> **52** di `BusinessType_DeT` — ⛔ saya memeriksanya untuk berkas layar, **tidak** untuk tabel
> keputusan.

---

## Bab A — Sasaran 1: lapisan layar, kini DIURAI

### A.1 ⭐ Berapa byte yang benar-benar dibaca

⭐ **Brief menuntut angka ini, jadi diberikan tanpa dibulatkan dan tanpa dikaburkan:**

| | Nilai |
| --- | ---: |
| berkas layar | **31** *(25 `Section` + 6 `Harness`)* |
| ⭐ **byte DIURAI oleh pengurai XML** | ⭐ **16.635.834 — 100 %** |
| ⭐ persen modul | ⭐ **46,6 %** dari 35.678.284 B |
| berkas yang **tidak** terurai | ⛔ **0** |

⚠️ **Dan bedakan dua hal, sebab brief memang menuntut kejujuran di sini:**
⭐ **Seluruh 16.635.834 byte MASUK ke pengurai** — tiap berkas dibuka penuh, di-`ElementTree`,
dan ditelusuri seluruh simpulnya. ⛔ **Yang tidak terjadi** adalah **saya membaca ke-16,6 MB itu
dengan mata**; yang masuk ke laporan ini adalah **ekstraknya**. ⭐ Jadi klaim *"dibuka"* di sini
berarti **terurai utuh oleh program**, ⛔ **bukan** *"terbaca utuh oleh manusia"*.
⚠️ Akibatnya: apa yang **tidak saya cari** tetap tidak ditemukan, walau berkasnya terurai penuh.

### A.2 ⛔ Sebelas dari dua belas nama tag yang saya duga TIDAK ADA

⭐ `CLAUDE.md` jebakan **#6** menuntut daftar tag diperiksa dulu. ⭐ **Diperiksa, dan hasilnya keras:**

| Tag yang mungkin ditebak | Ada? |
| --- | --- |
| ⛔ `pyFieldName` · `pyPropertyName` · `pyTabName` · `pyLayoutName` · `pyCellProperty` · `pyControlName` · `pyReference` · `pyDefaultValue`\* · `pyInitialValue` · `pyValidateName` · `pyMaxLength` | ⛔ **nol berisi di lapisan atas** |
| ⭐ `pyFormatType` | ⭐ **ada, 32 berisi** |

\* ⭐ `pyDefaultValue` **ada**, tetapi bukan di tempat yang saya duga — ia **anak `pyUserData`**,
bukan anak sel. Lihat A.5.

⭐ **Tag struktur yang BENAR-BENAR membawa arti**, ditemukan dengan menyisir daftar tag lebih dulu:

| Tag | Berisi | Isinya |
| --- | ---: | --- |
| ⭐ **`pyType`** | **2.554** | ⭐ jenis sel: `LABEL` **939** · **`FIELD` 870** · `LAYOUT` **238** · `SUB_SECTION` **86** · `HIERARCHY` **2** · *(kosong 419)* |
| ⭐ **`pyFormat`** | **820** | bentuk kendali — lihat A.4 |
| ⭐ **`pyValue`** | **1.769** | rujukan properti dan teks harfiah |
| ⭐ **`pyLabelFieldValue`** | **329** | ⭐ **label yang dibaca pengguna** |
| `pyFormatType` | 112 | `number` 80 · `date` 24 · `text` 8 |

⚠️ **Jebakan #2 diperiksa:** `<rowdata …/>` yang **menutup sendiri** = **12**, lawan **24.490**
berpasangan. ⭐ Jadi pola `<rowdata…>(.*?)</rowdata>` akan **melewatkan 12 sel kosong**.
⛔ **Saya tidak memakai pola itu** — dipakai `ElementTree`, yang menangkap keduanya.

### A.3 Delapan ratus tujuh puluh medan

`[terverifikasi]` **870 sel `FIELD`**, **754** membawa rujukan properti, ⭐ **140 properti berbeda**.

| Yang dicari brief | Terukur |
| --- | ---: |
| medan wajib isi pada sel `FIELD` | ⛔ **0** |
| medan hanya-baca pada sel `FIELD` | ⛔ **0** |
| ⭐ medan **terkunci permanen** | ⭐ **36** |
| ⭐ medan ber-**nilai awal** | ⭐ **11** |

⚠️⚠️ **Dan angka pertama itu bertabrakan dengan ronde 2, jadi saya nyatakan keduanya.**
⭐ Ronde 2 mengukur **`pyRequired` 160 berisi** di ke-31 berkas layar. ⭐ Ronde 4 mengukur
**0 di antaranya menempel pada sel `FIELD`**. ⛔ **Keduanya dapat benar bersamaan** — artinya
ke-160 itu menempel pada elemen **selain sel medan**. ⛔ **Di mana tepatnya belum saya telusur**,
dan ⛔ **saya tidak menyimpulkan bahwa modul ini tidak punya medan wajib.** `[terbuka]` → **P45**.

### A.4 Bentuk kendali, dan medan uang di layar

`[terverifikasi]` `pyFormat` pada sel `FIELD`:

| Kendali | Jumlah | |
| --- | ---: | --- |
| `pxNumber` | **219** | angka |
| `pxDisplayText` | **154** | ⭐ **tampil saja, tidak dapat diisi** |
| `pxTextInput` | **114** | teks |
| ⭐ **`pxCurrency`** | ⭐ **113** | ⭐ **medan UANG** |
| `pxButton` | **68** | tombol |
| `pxDateTime` | 38 | tanggal |
| `pxDropdown` · `pxRadioButtons` · `pxTextArea` · `pxCheckbox` | 36 · 11 · 9 · 6 | pilihan dan teks panjang |
| `pxIcon*` *(History, Attachments, ExpandCollapse, Cancel)* | 6 masing-masing | ikon |

⭐ **113 medan uang di layar**, dan ⭐ **154 medan hanya tampil** — ⚠️ yang kedua berarti
**seperlima layar adalah keluaran perhitungan**, bukan masukan.

### A.5 ⭐ Nilai awal yang ditetapkan LAYAR — dan jumlahnya sangat sedikit

`[terverifikasi]` `pyDefaultValue` berisi pada **11 medan saja**, dengan **2 nilai berbeda**:

| Nilai awal | Medan | Berkas |
| --- | ---: | --- |
| `0` | **9** | `DetailPolicyTreatyInNonProportional` + 1 berkas lain |
| ⚠️ **`Inclusive`** | **2** | `DetailPolicyTreatyIn` + 1 berkas lain |

⭐ **Bacaannya:** ⛔ **layar hampir tidak menetapkan nilai awal apa pun** — ⭐ penyiapan nilai
hidup di `DataTransform` *(ronde 3 §1)*, bukan di layar.
⚠️ **`Inclusive` adalah satu-satunya nilai awal berupa istilah bisnis**, dan ⛔ artinya
**belum terverifikasi** → **P46**.

### A.6 ⭐⭐ Sembilan belas properti yang HANYA ada di layar

⭐⭐ **Temuan terpenting Bab A.**

`[terverifikasi]` Dari **137 properti layar** yang diperiksa, ⛔ **19 tidak muncul di satu pun**
`Activity` *(92)*, `DataTransform` *(12)*, `When` *(75)*, atau `RDBList` *(41)* — ⭐ **jendela:
220 berkas aturan NB Treaty In, teks ter-unescape dua kali, sisa editor dibuang.**

| Kelompok | Properti |
| --- | --- |
| ⛔⛔ **UANG — nilai berpasangan mata uang** | **`LIMITVALUE`** · **`RETENTIONVALUE`** · **`EPIVALUE`** · **`NETPREMIVALUE`** · `EPICURRENCY` · `MDPCURRENCY` · `NETPREMICURRENCY` · `RETENTIONCURRENCY` |
| ⚠️ **klaim & survei** | `OutstandingClaim` · `DateofSurvey` · `SurveyedBy` · `LossPrevention` |
| periode & laporan | `Quartal` · `YearOfQuartal` · `StatementType` |
| lain | `Remarks` · `CurrentMode` · `FilterTermForOpportunity` · `TextNoQuotation` |

⛔⛔ **Delapan di antaranya medan UANG**, dan ⛔ **tidak satu pun aturan di modul ini menulis
atau membacanya.**

⭐ **Dan penjelasan yang paling mungkin sudah ada di ronde 3** — ⭐ `[dugaan]`, **belum
terverifikasi**: ronde 3 §2.1 menemukan kolom **`CLASSOFBUSINESS` memuat dokumen JSON** yang
**diadopsi** ke halaman `TreatyIn`, sehingga **struktur halaman terbentuk saat jalan**.
⭐ Delapan nama uang di atas **muncul persis sebagai kolom** `BrowseTreatyInDetail` *(ronde 3 §3.2:
`.LIMITCURRENCY`, `.LIMITVALUE`, `.RETENTIONCURRENCY`, `.RETENTIONVALUE`, `.EPICURRENCY`)*.
⭐ ⇒ `[dugaan]` **layar menampilkan medan yang lahir dari JSON, bukan dari aturan Pega.**
⛔ **Tidak saya tegaskan** — → **P42**.

⚠️ **Kenapa ini penting, dan bukan sekadar rapi:** ⛔ bila sistem baru dibangun dari **aturan**
saja, **kedelapan medan uang ini tidak akan pernah terisi** — dan ⛔ **tidak ada pesan galat**
yang memberi tahu.

---

## Bab B — Sasaran 2: sensus ulang elemen mati dan terkunci

### B.1 ⭐ Ronde 2 DIKONFIRMASI — 91 dan 38, dan sekarang diketahui kenapa 89 juga muncul

⭐ **Dihitung dua cara yang benar-benar berbeda, dan keduanya dilaporkan.**

| Cara | MATI | TERKUNCI |
| --- | ---: | ---: |
| ⭐ **A — ungkapan ATOM saja** *(`1=2`, `NEVER`, `1==2`)* | **89** | — |
| ⭐ **B — atom DITAMBAH gabungan `&&` yang memuat bagian mustahil** | ⭐ **91** | — |
| `pyReadOnlyCondition` yang **selalu benar** | — | ⭐ **38** |
| **ronde 2 menulis** | **91** | **38** |

⭐⭐ **Yang dipercaya: CARA B — 91.** ⭐ Sebabnya: dua ungkapan berbentuk
**`<sesuatu> && NEVER`**, dan ⭐ **sebuah gabungan DAN yang memuat bagian mustahil juga tidak
pernah benar.** ⛔ Cara A membuangnya hanya karena bentuknya tidak persis sama.
✅ **Jadi angka ronde 2 BENAR, dan selisih 89 lawan 91 kini terjelaskan seluruhnya.**

**Sebarannya** `[terverifikasi]`: `pyCondition` **78** + `pyContainerVisibleWhen` **11** = 89 atom,
⭐ ditambah **2 gabungan** pada `pyCondition` ⇒ **91**.
`pyReadOnlyCondition` selalu-benar ⇒ **38**.

### B.2 ⭐⭐ APA yang disembunyikan tiap elemen mati — dan jawabannya mengubah keputusannya

⭐ **Inilah yang brief minta, dan yang belum pernah dijawab.**

`[terverifikasi]` Untuk tiap elemen mati, ditelusuri **seluruh sel `pyType` di bawahnya**:

| Yang disembunyikan | Jumlah sel |
| --- | ---: |
| ⭐ **`FIELD`** | ⭐ **128** |
| `LABEL` | 132 |
| `LAYOUT` | 20 |
| `SUB_SECTION` | 5 |

⭐⭐ **Dan sebarannya adalah temuan sebenarnya:**

| Besarnya yang disembunyikan | Elemen |
| --- | ---: |
| ⛔⛔ **kosong — nol medan, atau hanya label** | ⛔ **80** |
| **2–5 medan** | 6 |
| **6–20 medan** | 1 |
| ⭐ **lebih dari 20 medan** | ⭐ **4** |

⭐⭐ **Bacaannya, dan ia menghemat pekerjaan besar:**
⛔ **80 dari 91 elemen mati tidak menyembunyikan satu medan pun.** ⭐ Membangunnya berbiaya
**nol**, dan membuangnya juga berbiaya nol.
⭐ **Seluruh persoalan sebenarnya ada pada 11 elemen yang menyembunyikan 128 medan** —
⛔ dan **104 dari 128 itu terkumpul di EMPAT elemen saja**, di dua berkas yang nyaris kembar:

| Medan | Berkas | Tag |
| ---: | --- | --- |
| **26** | `BusinessAndSOBList` | `pyContainerVisibleWhen` = `1=2` |
| **26** | `BusinessAndSOBListRetro` | `pyContainerVisibleWhen` = `1=2` |
| **26** | `BusinessAndSOBList` *(elemen kedua)* | `pyContainerVisibleWhen` = `1=2` |
| **26** | `BusinessAndSOBListRetro` *(elemen kedua)* | `pyContainerVisibleWhen` = `1=2` |
| 6 | `Installments_ReadOnly` | `1=2` |
| 4 · 4 · 4 | `DetailPolicyTreatyIn` · `GeneralPolicyTreatyIn` · `SpreadingRiskList` | `1=2` |
| 2 | `DetailPolicyTreatyInNonProportionalEDM` | ⭐ `NEVER` |

**Per berkas** `[terverifikasi]`: `DetailDeptHeadTreatyIn_UW` **18** · `DetailPolicyTreatyIn`
**18** · `GeneralDeptHeadTreatyIn_UW` **18** · `GeneralPolicyTreatyIn` **18** ·
`BusinessAndSOBList`/`Retro` **4** masing-masing · `SFAPortalOpportunities` 3 · sisanya 1–2.

⭐ ⇒ **72 dari 91 ada di empat berkas utama**, dan ⭐ **hampir seluruhnya termasuk 80 yang kosong.**

### B.3 ⭐ Elemen TERKUNCI mengunci SATU medan masing-masing — dan medannya bernama

⭐ Ronde 2 hanya dapat mengatakan *"38 elemen hanya-baca permanen"*. ⭐ **Ronde 4 tahu apa yang
dikunci.**

`[terverifikasi]` `pyReadOnlyCondition` tinggal di elemen **`<pyUserData>`**, yang **induknya
adalah sel ber-`pyType = FIELD`**. ⭐ ⇒ ⭐ **setiap kunci mengunci tepat SATU medan**, ⛔ bukan
bagian atau tab.

| Berkas | Terkunci |
| --- | ---: |
| ⭐ `DetailDeptHeadTreatyIn_UW` | **18** |
| ⭐ `GeneralDeptHeadTreatyIn_UW` | **18** |
| `DetailPolicyTreatyIn` | 1 |
| `GeneralPolicyTreatyIn` | 1 |

⭐ **Tiga medan yang terbaca namanya**, sebagai contoh — label apa adanya seperti dilihat pengguna:

| Properti | Label di layar |
| --- | --- |
| `.StartDate` | ⭐ **"Statement Period"** |
| `.QuotationData.IsSurveyReport` | ⭐ **"Survey Report"** |
| `.StatementDate` | ⭐ **"Statement Date"** |

⭐⭐ **Dan ini menajamkan ronde 2 §2.2.** Ronde 2 menyimpulkan `[dugaan]` *"layar Kepala
Departemen dibuat hanya-baca dengan cara dipaku"*. ⭐ **Kini lebih tepat:** ⛔ bukan **layarnya**
yang dipaku — ⭐ **36 medan tertentu yang dipaku, satu per satu**, dan **dua di antaranya bahkan
di layar biasa, bukan layar Kepala Departemen.**

⚠️ **Selisih 38 lawan 36 dinyatakan apa adanya:** **38** bila dihitung dari elemen pembawa syarat;
**36** bila dihitung dari sel `FIELD` yang tersentuh. ⭐ Selisih **2** berarti dua kunci menempel
pada elemen yang **bukan** sel medan. ⛔ **Belum ditelusur** — `[terbuka]`.

---

## Bab C — Sasaran 3: ejaan nama tabel, rujukan TABEL dipisah dari rujukan PROPERTI

⭐ **Jendela: 1.151 naskah `pyBrowseSQL` di 21 modul korpus.** ⛔ `nusantara-re` tidak disentuh.
⭐ Rujukan **tabel** diambil hanya yang muncul sesudah **`FROM` · `JOIN` · `INTO` · `UPDATE`**;
⭐ rujukan **properti** *(berbentuk `<sesuatu>.NAMA`)* **dibuang**, dan jumlahnya dilaporkan.

| Nama | ber-`POOLDATA.` | tanpa awalan | ⚠️ rujukan PROPERTI yang dibuang |
| --- | ---: | ---: | ---: |
| `REINSURANCETYPE` | **18** | **10** | 0 |
| `PROPORTIONALARRG` | **2** | **22** | 0 |
| `HISTORYAKSEPTASIPEGA` | ⭐ **0** | **15** | 0 |
| ⭐ `CURRENCY` | **4** | ⭐ **24** | ⛔⛔ **54** |

⭐⭐ **Peringatan brief terbukti benar, dan besarannya besar:** ⛔ **54 kemunculan `…​.CURRENCY`
adalah rujukan PROPERTI, bukan nama tabel** — ⭐ **lebih banyak daripada ke-28 rujukan tabelnya
digabung.** ⛔ Penelusuran naif akan **melipatgandakan** angkanya.

⚠️ **Dua selisih kecil terhadap angka brief, dilaporkan apa adanya:**

| Nama | Brief | Ronde 4 | Sebab |
| --- | ---: | ---: | --- |
| `REINSURANCETYPE` tanpa awalan | 11 | **10** | ⭐ ronde 4 hanya menghitung yang berada di **posisi tabel** |
| `PROPORTIONALARRG` tanpa awalan | 25 | **22** | sama |

⭐ **Yang dipercaya: angka ronde 4**, sebab jendelanya **dinyatakan dan lebih sempit** — hanya
posisi tabel. ⛔ **Tetapi angka brief tidak saya sebut salah** — ⚠️ selisihnya **belum saya
telusur satu per satu**, dan bisa jadi brief menangkap bentuk rujukan yang pola saya lewatkan.
⭐ **Keduanya dicantumkan**, sesuai aturan sensus.

⭐ **Berkas mana memakai ejaan mana** — yang menyentuh **NB Treaty In**:

| Nama | Di NB Treaty In |
| --- | --- |
| `REINSURANCETYPE` | ⚠️ **tanpa awalan** — `RDBList\GetTreatyName.xml` |
| `PROPORTIONALARRG` | ⭐ **ber-`POOLDATA.`** — `RDBList\GetBreakDownSpread_SQL.xml` |
| `CURRENCY` | ⭐ **ber-`POOLDATA.`** — `RDBList\GetCurrencyIDByName.xml` |
| `HISTORYAKSEPTASIPEGA` | ⚠️ **tanpa awalan** *(ronde 1 §A.2)* |

⚠️⚠️ **`REINSURANCETYPE` dan `PROPORTIONALARRG` dieja BERBEDA arah di modul yang sama** — satu
tanpa awalan, satu dengan. ⛔ **Apakah keduanya tabel yang sama TIDAK saya simpulkan** — ⭐ itu
keputusan DBA, dan sudah menjadi **P3**.

---

## Bab D — Sasaran 4: kolom `OPERATORID`

⭐ **Jawaban tegas, dari korpus, dengan lingkup disebut.**

⭐ **Jendela: 1.151 naskah `pyBrowseSQL` di 21 modul.**

| Yang dicacah | Jumlah |
| --- | ---: |
| naskah SQL yang menyebut teks `OPERATORID` **dalam bentuk apa pun** | **67** |
| ⚠️ di antaranya rujukan **halaman Pega** `OperatorID.pyXxx` | ⛔ **67 kemunculan — seluruhnya** |
| ⭐ **naskah SQL yang MENULIS ke kolom `OPERATORID`** | ⭐⭐ **0** |
| naskah SQL yang menyebutnya **sebagai kolom** tanpa menulis | ⭐ **0** |

⭐⭐ **Jadi P4 terjawab dari korpus:** ⛔ **tidak satu pun dari 1.151 naskah SQL di seluruh korpus
menulis ke kolom `OPERATORID`, dan tidak satu pun bahkan menyebutnya sebagai kolom.**
⭐ Ke-67 kemunculannya **seluruhnya** adalah **parameter dari halaman pengguna Pega** — persis
jebakan yang brief peringatkan.

⚠️ **Lingkup penelusuran dinyatakan, sesuai jebakan #7:** ⭐ yang disisir adalah **isi
`<pyBrowseSQL>`** saja. ⛔ **Tidak** disisir: stored procedure *(badannya tidak ada — OQ-002)*,
pemicu basis data, pekerjaan terjadwal, dan sistem lain di luar Pega.
⭐ ⇒ ⛔ **Pernyataan yang sah hanyalah: kolom itu tidak ditulis dari Pega.** ⛔ **Bukan** bahwa ia
tidak pernah terisi. → **P47**.

---

## Bab E — Bukti isolasi

| Yang dibuktikan | Keadaan |
| --- | --- |
| ⛔ `D:\XML\nusantara-re\` | ✅ **NOL disentuh** — tidak pernah menjadi argumen; daftar modul menyaringnya keluar |
| korpus `D:\XML\RNM_BRD\` | ✅ **READ-ONLY** — hanya dibaca; penulisan hanya ke `OUTPUT_HASIL_RNM\` |
| ⭐ modul lain dibaca | ⚠️ **ya, 20 modul** — ⭐ **hanya** untuk Bab C dan D, ⭐ **dan brief mengizinkannya** untuk membuktikan ada/tidaknya sesuatu |
| ⛔ temuan modul lain masuk sensus NB Treaty In | ✅ **NOL** — Bab A dan B **murni 31 berkas NB Treaty In**; Bab C dan D menyebut jendelanya **1.151 / 21 modul** setiap kali |
| ronde 1 · 2 · 3 disunting | ✅ **NOL** — ralat ditulis **di berkas ini** |
| berkas dibuat | ✅ **tepat DUA**, sesuai perintah brief |
| kode · DDL · `CREATE TABLE` · usulan kolom | ✅ **NOL** |
| ⛔ nilai nama orang | ✅ **NOL disalin** |
| butir terbuka ditutup sendiri | ✅ **NOL** |
| rahasia · token · data nasabah | ✅ **NOL** |

⚠️ **Baris tidak bersih:** sesi ini panjang dan pernah disambung; harness pernah menyuntik ulang
dua bacaan `.scratch` konteks lain. ⛔ **Nol nama tabel, kolom, atau pola rancangannya merembes** —
seluruh angka ronde 4 dari sisiran korpus sendiri.

---

## TELEMETRI EKSEKUSI

⛔⛔ **Pengukuran dari luar TIDAK dilakukan.** ⭐ Brief meminta `claude --print --output-format json`
dijalankan **dari luar sesi**; ⛔ ronde ini dijalankan **di dalam sesi interaktif**, sehingga
`total_cost_usd` dan `duration_ms` **tidak tersedia**. ⭐ **Dinyatakan apa adanya, tidak ditaksir.**

⭐ **Yang DAPAT diukur dari dalam** — selisih terhadap baseline yang diambil sebelum ronde ini,
dibaca dari transkrip sesi:

| Yang dicatat | Nilai | Sumber |
| --- | ---: | --- |
| ⭐ **token keluaran** | ⭐ **154.518** | transkrip `.jsonl`, SELISIH terhadap baseline |
| ⭐ **token cache-read** | ⭐ **26.189.892** | sama |
| token cache-write | **218.953** | sama |
| token masuk | **94** | sama |
| ⭐ **jumlah panggilan alat** | ⭐ **47** | sama |
| ⛔ durasi | ⛔ **TIDAK DAPAT DIUKUR** | butuh `--print` dari luar sesi |
| ⛔ biaya | ⛔ **TIDAK DAPAT DIUKUR** | butuh `--print` dari luar sesi |
| ⭐ **berkas dibuka / byte diurai** | ⭐ **31 berkas layar · 16.635.834 B (100 %)** · ditambah 2 `DecisionTable`, 220 berkas aturan NB, dan **1.151 naskah SQL di 21 modul** | perintah audit di tiap bab |

⚠️ **Tiga keterbatasan pengukuran dari dalam, disebut apa adanya:**
*(a)* pesan terakhir **belum tertulis** ke transkrip saat pengukuran kedua — angkanya **kurang satu
pesan**; *(b)* perintahnya mengambil transkrip **paling baru diubah** — jalurnya **diperiksa** dan
cocok dengan sesi ini; *(c)* yang dilaporkan adalah **SELISIH**, bukan total sesi.

⭐ **Bandingkan dengan tiga ronde sebelumnya** — seluruhnya SELISIH, alat yang sama:

| | R1 | R2 | R3 | ⭐ **R4** |
| --- | ---: | ---: | ---: | ---: |
| keluaran | 192.021 | 225.898 | 178.857 | ⭐ **154.518** |
| cache-read | 13.532.360 | 22.688.258 | 16.725.320 | ⚠️ **26.189.892** |
| panggilan | 69 | 72 | 39 | **47** |

⭐ **Keluaran terendah dari empat ronde**, ⚠️ **cache-read tertinggi** — ⭐ sebab ronde ini
membaca **tiga berkas grilling sebelumnya** di tiap giliran, dan menyisir **1.151 naskah SQL di
21 modul** untuk Bab C dan D.

