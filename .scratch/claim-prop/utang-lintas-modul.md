# Utang lintas modul — **Claim Prop**

> **Ronde TEMUAN, bukan ronde keputusan.** ⛔ Nol berkas lama disunting. ⛔ Nol butir `[terbuka]`
> ditutup. ⛔ Nol ADR direvisi. Semua usulan di bawah menunggu keputusan work owner.

**Sensus blok ini:** vonis ADR **15** · RALAT dibutuhkan **4** · butir `[terbuka]` baru **3** ·
usulan teks **1** · premis brief yang terbukti **keliru 2**

**Jendela sensus:** seluruh **329** berkas `.xml` `Claim Prop`, ditambah **80** berkas
`Komite Claim Prop` dan **482** berkas `Claim Fac In` sebagai pembanding. Dokumen: **26** berkas
`.scratch\claim-prop` dan **15** ADR.

---

## §1 — Sensus alur dua cara, dan **pencabutan kecemasan yang saya sendiri timbulkan**

### 1a — ⛔ Premis brief ini **KELIRU**, dan itu kekeliruan saya

Brief ini lahir dari peringatan yang saya sampaikan sendiri: bahwa sensus bentuk alur memakai
`pyShapeType` itu **buta**, dan `pxObjClass` memberi **15** untuk `Flow_TreatyIn`, bukan **8**.
⛔ **Kesimpulan itu salah.** Yang saya hitung sebagai "bentuk" ternyata **empat jenis entri berbeda
yang dijumlahkan menjadi satu angka**.

`[terverifikasi]` Bentuk alur tinggal di **`pyModelProcess`**, dan di sana ada **tiga wadah
sejajar** ditambah **dua milik**:

| Wadah | Isinya | Apakah ia "bentuk"? |
| --- | --- | --- |
| `pyShapes` | bentuk di kanvas | ⭐ **ya** |
| `pyConnectors` | penghubung antar bentuk | tidak — ia **garis** |
| `pyModifiers` | tiket tingkat alur | tidak — ia **label** |
| `pyTicketShapes` *(di dalam sebuah bentuk)* | tiket yang menggantung pada bentuk itu | tidak — ia **milik** bentuk |
| `pyRouterProp` · `pyNotifyProp` *(di dalam sebuah bentuk)* | perute dan pemberitahu | tidak — ia **milik** bentuk |

### 1b — Angka yang benar

`[terverifikasi]` Dihitung **dua cara**, dan keduanya menutup:

| | cara 1 `pyShapeType` | mentah `pxObjClass` | ⭐ **bentuk** | penghubung | pengubah | tiket |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `Flow_TreatyIn` **(Claim Prop)** | 8 | 22 | ⭐ **5** | **5** | 1 | 4 |
| `KomiteTreaty_Flow` *(pembanding)* | 6 | 17 | **4** | **4** | 1 | 3 |
| `Register_Flow` *(Claim Fac In, pembanding)* | 16 | 38 | **8** | **10** | 2 | 7 |

**Uji tutup:** `bentuk 17 + penghubung 19 + pengubah 4 + tiket 14 + perute 17 + pemberitahu 6 = 77`,
dan sapuan mentah `pxObjClass` juga **77**. Tidak ada entri yang hilang atau terhitung dua kali.

⭐ **Kelima bentuk `Flow_TreatyIn`:** `Start1` *(mulai)* · `Assignment1` *(**Input Acceptation**,
`WorkBasket`)* · `Assignment2` *(**Outstanding Claim**, `WorkList`, `pyRouteTo = Current operator`)*
· `Decision3` *(**`IsBack`**, `DataXOR`)* · `End1` *(selesai)*.

`[terverifikasi]` **Bentuk selesai ADA**, dan `pyWorkStatus`-nya **`Resolved-Completed`** — satu
bentuk selesai, satu status.

⚠️ `Assignment1` ber-`pyBaseClass` **`ASM-FW-GCNMFW-Work-PNC`** sedangkan alurnya berkelas
`…-Work-ClaimTreaty`. Itu **sisa impor model lain di tingkat ALUR** — sekerabat dengan temuan
"52 dari 61 rule `When`" di §4. ⛔ Belum pernah tercatat. `[terbuka]`

### 1c — ⭐ Kelima "perkecualian" itu: **satu tiket bernama, empat cangkang kosong**

`[terverifikasi]` Dibaca satu per satu:

| Menggantung di | Nama | Isi |
| --- | --- | --- |
| `End1` | *(tanpa nama)* | ⛔ **hanya `pxObjClass`** — cangkang kosong |
| `Decision3` | *(tanpa nama)* | ⛔ cangkang kosong |
| `Assignment2` | *(tanpa nama)* | ⛔ cangkang kosong |
| ⭐ `Assignment1` | **`AcceptanceClaim`** | tiket bernama |
| `pyModifiers` | **`AcceptanceClaim`** | **definisi** tiket yang sama |

⭐ **Jadi tiketnya SATU, bukan lima** — `AcceptanceClaim`, menggantung pada langkah **Input
Acceptation**. Empat sisanya cangkang tanpa nama dan tanpa pengenal.

⭐ **`AcceptanceClaim` adalah rule sungguhan**, bukan sekadar label gambar: alur merujuknya sebagai
`pxRuleObjClass = Rule-Obj-Ticket` di kelas **`ASM-FW-GCNMFW-Work`** — kelas **induk bersama**,
sehingga ia **terlihat lintas modul**. Catatan pengembang pada alurnya berbunyi
*"add tiket AcceptanceClaim"*.

### 1d — ⭐⭐ Vonis: **tiketnya LABEL MATI** — tak ada yang membangkitkannya

Tiket Pega adalah **sasaran lompatan**: siapa pun yang memanggil metode `Obj-Set-Tickets` dapat
melempar alur ke titik itu, **dari mana saja**. Kalau ada yang membangkitkannya, ia **jalan masuk
tersembunyi**.

`[terverifikasi]` **Jendela: seluruh 409 berkas `.xml` dua modul, SELURUH tag — bukan daftar medan
pilihan.**

| Yang dicari | Claim Prop | Komite Claim Prop |
| --- | --- | --- |
| metode langkah **`Obj-Set-Tickets`** | ⭐ **NOL** | ⭐ **NOL** |
| activity **`SetTicket`** *(pembangkit baku Pega)* | ⭐ **NOL** | ⭐ **NOL** |
| pemanggil `SetTicket` | ⭐ **NOL** | ⭐ **NOL** |
| `AcceptanceClaim` di luar berkas alurnya sendiri | ⭐ **NOL** | — |

`[terverifikasi]` `AcceptanceClaim` muncul di **tepat tiga berkas** pada **seluruh 20 korpus** —
`Claim Prop`, `Claim Non Prop`, `Claim Fac In` — dan **ketiganya berkas alur itu sendiri**.

⭐⭐ **Karena itu: tidak ada jalan kembali tersembunyi lewat tiket di Claim Prop.** Label itu
terpasang di kanvas dan **tidak pernah dilempar**.

⚠️ **Batas kejujuran vonis ini, dan ia nyata:** mesin pembangkitnya **ada di aplikasi yang sama**.
`Activity\SetTicket.xml` hidup di **tiga modul lain** — `NB FacIn`, `RNW Fac In`,
`Endorsment Fac In` — dan ia rule **`@baseclass`** bawaan Pega. Karena `AcceptanceClaim` terdaftar
di kelas induk bersama, sebuah rule **di luar ekspor ini** secara teknis **dapat** melemparnya.
⛔ **Tidak ada buktinya di korpus**, dan saya tidak menyimpulkan ada. `[terbuka]`

---

## §2 — Lima belas ADR diadu dengan Claim Prop

### 2a — ⛔ Premis kedua brief ini juga **KELIRU**

Brief menyatakan kedua modul *"tuntas tanpa pernah sekalipun diadu dengan ke-15 dokumen itu"*.
`[terverifikasi]` Untuk Claim Prop itu **tidak benar**: **9 dari 15 ADR sudah dirujuk** di
dokumennya — ADR-0003 *(13 sebutan)*, 0006 *(8)*, 0011 *(8)*, 0014 *(8)*, 0015 *(8)*, 0010 *(5)*,
0013 *(5)*, 0005 *(1)*, 0012 *(1)*.

⭐ **Yang benar-benar belum pernah disebut: ENAM** — **ADR-0001 · 0002 · 0004 · 0007 · 0008 ·
0009**. Di situlah nilai blok ini, bukan di seluruh lima belas.

### 2b — Lima belas vonis

| ADR | Pokok | Vonis Claim Prop | Bukti |
| --- | --- | --- | --- |
| **0001** | batas konteks lewat *child work* | ✅ **MENDUKUNG** | `pxAddChildWork` di **2 activity** — `AddKomiteTreatyChild_ACT`, `SendCloseClaimToKomite`. Bentuk kontraknya sama dengan Life. ⛔ **belum pernah disebut** |
| **0002** | RBAC tiga peran | **TIDAK MENYENTUH** ⚠️ | ⭐ lihat 2c — **jebakan nama** |
| **0003** | uang bukan `float` | ⛔⛔ **MENENTANG** | ⭐ lihat 2d — **12 parameter uang bertipe `Double`** |
| **0004** | alamat jadi env var | *(tidak dinilai — berstatus `superseded`, digantikan 0013)* | — |
| **0005** | `IsPEGAPROD` jadi flag lingkungan | ✅ **MENDUKUNG** | `IsPEGAPROD` di **7 berkas**, a.l. `HitServiceToKasir_Act`, `InsertGoogleStorage_Act`, `KonversiKlaim_Act`, `SendEmailKlaim` — keadaan yang persis dijelaskan ADR |
| **0006** | penomoran lewat *stored procedure* | ✅ **MENDUKUNG** | `GetSequenceNumber_SQL` di **4 berkas** — `SaveOutstanding_Act`, `PrintDLATreatyIn`, `TryMakePLA_Act` |
| **0007** | jejak audit tiap transisi + jalur balik | ✅ **MENDUKUNG** ⚠️ | kolom jejak **hanya nama operator**: `UPDATEOPNAME` 11 · `CREATEOPNAME` 3 · `CREATEDATETIME` 1. **Kekurangan yang sama dengan Life** ⇒ perbaikan ADR berlaku. ⛔ **belum pernah disebut** |
| **0008** | efek keluar asinkron, **tidak memblokir** | ⛔ **MENENTANG — sengaja** | `IsSuccessHitService` / `StsKonversi`: ⭐ **NOL** *(sama seperti Life)*, tetapi spec Claim Prop bab 12 **memilih outbox transaksional (0015)**. Pertentangannya **sadar dan sudah tertulis** |
| **0009** | migrasi penuh, tanpa koeksistensi | ⚠️ **belum punya data** | ⛔ spec tidak pernah memutuskan **penuh vs koeksistensi**. Yang ada hanya *"tidak dimigrasikan"* tentang **RULE**, bukan **DATA**. ⭐ **lubang sejati** — lihat §5 |
| **0010** | berkas tetap di Google Storage | ✅ **MENDUKUNG** ⚠️ | **396 kemunculan / 12 berkas**. ⚠️ tetapi lihat §3 — versi rule yang dipakai **lebih tua dan bisa bertabrakan** |
| **0011** | unit status = baris `AdjustmentList` | ✅ **MENDUKUNG** | `AdjustmentList` **145 / 19 berkas**; sudah dirujuk 8× |
| **0012** | wewenang kirim komite bergantung `Type` | **TIDAK MENYENTUH** | ⭐ `.Type` = `TP`/`TR`/`QP`/`QR`: **NOL** — **sama persis dengan Claim Fac In** |
| **0013** | alamat di-*lookup* dari `M_LINK_SERVICE` | ✅ **MENDUKUNG** | `GetLinkService` + kunci kategori, **14 berkas** |
| **0014** | keputusan komite ditegakkan per `KomiteID` | ✅ **MENDUKUNG** — ⛔ **DIRALAT 2026-09-19** | ⭐ **AC 57** modul ini berbunyi *"Seluruh gerbang wewenang **ditegakkan di lapisan layanan** (**ADR-0014**)"*, dan **AC 58** mengujinya lewat **endpoint tanpa UI**. ⛔ Vonis lama dikutip utuh, tidak dihapus: ~~**TIDAK MENYENTUH** — wilayahnya modul komite; Claim Prop hanya **memasok** `KomiteID` *(3 berkas)*. `AcceptStatus` ditulis **satu** kali~~. ⚠️ **Sebab kekeliruan saya:** saya menyapu **korpus Pega** dan menyimpulkan dari ketiadaan pemeriksaan di sana — padahal ADR menanyakan **apa yang dibangun di sistem baru**, dan **spec-nya sudah menjawab**. ⭐ **Temuan utama tidak melemah, justru menajam:** modul saudara **sudah menegakkan**, dan **modul komite yang tertinggal**. |
| **0015** | efek keluar komite **wajib berhasil** | ✅ **MENDUKUNG** | dipilih eksplisit di spec bab 12 sebagai penyimpangan sadar |

**Rekap:** ✅ MENDUKUNG **8** · ⛔ MENENTANG **2** · TIDAK MENYENTUH **2** · belum punya data **1** ·
tidak dinilai **1** *(superseded)*.

> ⛔ **RALAT 2026-09-19** — rekap lamanya **dikutip, tidak dihapus**: *"✅ MENDUKUNG **7** ·
> ⛔ MENENTANG **2** · TIDAK MENYENTUH **3** · belum punya data **1** · tidak dinilai **1**"*.
> Satu vonis bergeser: **ADR-0014 TIDAK MENYENTUH → MENDUKUNG**. ⛔ **Tidak ada vonis lain yang
> diubah**, dan berkas ini **tidak disunting untuk hal lain** — ia berkas temuan, bukan berkas hidup.

### 2c — ⭐ ADR-0002: **jebakan nama `pyPosition`**

Sapuan pertama memberi `pyPosition` **115 kemunculan / 28 berkas** — terkesan modul ini penuh
pemeriksaan peran. ⛔ **Itu keliru.**

`[terverifikasi]` Dibongkar menurut **jenis rule** dan **nilai pembandingnya**:

| Nilai | Berapa | Di mana | Artinya |
| --- | ---: | --- | --- |
| `Top` | **38** | ReportDefinition · Section | ⛔ **tata letak layar** |
| `AFTER` | **18** | ReportDefinition · Section | ⛔ **tata letak layar** |
| ⭐ `IT Developer` | **2** | `DataTransform\InsertChronology_DT` | ⭐ **wewenang sejati** |

⭐ **Nama peran gaya Life — `ReasLifeAdmin` dan kerabatnya — NOL.** Sapuan pola `Reas*` tidak
menemukan satu pun. Model tiga peran ADR-0002 **tidak ada** di sini.

⚠️⚠️ **Dan satu temuan yang belum pernah tercatat:** satu-satunya pemeriksaan wewenang sejati
berbunyi **`OperatorID.pyPosition != "IT Developer"`** — sebuah **pintu belakang pengembang** yang
digantungkan pada **teks jabatan**, bukan pada peran. Ia mengatur penghapusan catatan kronologi.
⭐ Rule yang sama **ada persis** di Komite Claim Prop. ⛔ Perlakuannya belum diputuskan. `[terbuka]`

### 2d — ⭐⭐ ADR-0003: **MENENTANG**, dan buktinya lebih kuat daripada di Claim Fac In

`[terverifikasi]` **12 parameter berbau uang bertipe `Double`** di **5 berkas**:

| Berkas | Parameter | Sumber tipe |
| --- | --- | --- |
| `Activity\AddKomiteTreatyChild_ACT` | **`TSISpread`** · **`ClaimSpread`** | ⭐ **tanda tangan DAN medan — keduanya** |
| `Activity\SetKomiteTreaty_ACT` | **`TotalAdjustment`** | ⭐ **tanda tangan DAN medan** |
| `Activity\CountValueADJTreaty_Act` | **`TotalAdjustment`** | medan |
| `Section\ComiteeClaimTreaty` | `TSISpread` · `ClaimSpread` *(2×)* | medan |
| `Harness\CommitteeTreaty` | `ClaimSpread` | medan |

⭐ **Ini bukti yang lebih kuat daripada Claim Fac In.** Di Fac In kedua sumber tipe **berselisih**,
sehingga vonisnya harus bersandar pada tambalan koma-ke-titik di SQL. **Di sini keduanya sepakat**
pada tiga parameter — `TSISpread`, `ClaimSpread`, `TotalAdjustment` — sehingga tidak ada celah
"sumber mana yang sah".

⚠️⚠️ **Yang membuatnya berat:** `AddKomiteTreatyChild_ACT` adalah **rule yang melahirkan kasus
komite**. Jadi **uang menyeberang batas modul lewat gerbang `Double`** — Claim Prop → Komite Claim
Prop. Komite sendiri **bersih** *(nol `Double`/`Float`)*, tetapi **pintu masuknya tidak**.

⛔ **Yang TIDAK saya simpulkan:** bahwa presisi benar-benar hilang. `Double` adalah tipe **parameter
Pega**, dan apakah nilainya sempat melewati aritmetika biner **belum dibaca**. Yang tegas: ADR-0003
melarang **representasi** uang sebagai bilangan pecahan biner, dan **representasi itu ada di sini**.
`[terbuka]`

---

## §3 — Empat rule beda versi: **satu kecemasan gugur, satu cacat nyata muncul**

Dasar utang ini: dari **162** rule beridentitas sama antara Claim Prop dan Claim Fac In, **158
identik** dan **4 beda versi** — keempatnya **lebih baru di Fac In**.

⚠️ **Cara membacanya.** `md5` mentah keempat berkas **selalu berbeda** — bahkan antara Claim Prop
dan Komite Claim Prop yang **ukuran dan cap waktunya sama persis**. Sebabnya metadata ekspor
*(`pxHostId`, `pxMoveImport*`, `pzChecksum`, stempel penguncian)*. ⛔ **`md5` berkas karena itu
BUKAN alat yang sah** untuk pertanyaan ini. Yang dipakai: **perbandingan medan demi medan dengan
medan-waktu dan medan-ekspor dibuang.**

### 3a — ⭐⭐ Kurs standar: **isinya IDENTIK. Kecemasan 4 tahun 3 bulan GUGUR.**

`[terverifikasi]` `ASM-FW-GISFW-INT-CURRENCYSTANDARD!ASM!CURRENCYSTANDARD` — **12 medan berbeda dari
74**, dan **nol** di antaranya logika:

> `pxHostId` · `pxMoveImportDateTime` · `pxMoveImportOperId` · `pxMoveImportOperName` ·
> `pxOriginalRuleSet` · `pyRuleFormStatusTime` · `pyLabel` · `pySPRuleSetName` ·
> `pyShowJavaWindowName` · `pzChecksum` · `pzOriginalInstanceKey` · `pzRuleSetVersionPatch`

⭐ **Perintah SQL-nya sama kata demi kata di ketiga modul** — memanggil fungsi basis data
`getcurrencystandard` dengan mata uang dan `sysdate`, menghasilkan satu kolom `nilaiKurs`.

⭐⭐ **VONIS: kesimpulan kurs standar Claim Prop TIDAK BERUBAH.** Selisih 4 tahun 3 bulan itu
**selisih label versi ruleset**, bukan selisih perilaku. Butir utang ini **boleh ditutup** — ⛔ dan
penutupannya tetap keputusan work owner, bukan keputusan saya.

⚠️ **Satu ganjilan yang saya catat tanpa menyimpulkan:** `pzOriginalInstanceKey` salinan Claim Prop
menunjuk sebuah rule **pemutakhir master mata uang**, sedangkan salinan Fac In menunjuk
**`CURRENCYSTANDARD`** itu sendiri — pertanda salinan Claim Prop lahir dari **"simpan-sebagai" rule
lain**. ⛔ Tidak mengubah SQL-nya. `[terbuka]`

### 3b — ⭐⭐ Pembuat pengenal berkas: **cacat NYATA, dan tambalannya sudah ada**

`[terverifikasi]` `…!RNM!GENERATEIMAGEID_SQL` — **di sini isinya benar-benar berbeda**:

| | Bahan yang di-*hash* |
| --- | --- |
| **Claim Prop** *(dan Komite — sama persis)* | awalan tetap **+ cap waktu berketelitian MILIDETIK** |
| ⭐ **Claim Fac In** *(lebih baru)* | awalan tetap **+ cap waktu berketelitian NANODETIK + `SYS_GUID()`** |

⭐ **Dua hal ditambahkan versi baru:** ketelitian waktu naik **seribu juta kali** *(milidetik →
nanodetik)*, dan ditambahkan **nilai unik sejagat** dari basis data.

⭐⭐ **Itu persis bentuk sebuah TAMBALAN TABRAKAN** — dan ia **membenarkan temuan Claim Prop
sendiri.** `grilling-ronde-2` sudah menuliskan: *"keunikan bergantung resolusi **milidetik**. **Dua
unggahan dalam milidetik sama menghasilkan ImageID identik**"*.

✅ **Jadi kesimpulan Claim Prop bukan gugur — ia DIKUATKAN, dan kini punya jalan keluar yang sudah
terbukti dipakai di modul lain.**

⚠️ **VONIS: kesimpulan penyimpanan berkas BERUBAH — bukan dibatalkan, melainkan bertambah.** Versi
yang dipakai Claim Prop dan Komite adalah versi **yang belum ditambal**.

### 3c — Jenis berkas: bertambah, bukan berubah

`[terverifikasi]` `…!GETMIMETYPE` — **18 medan berbeda dari 162**, dan yang bermakna hanya satu:
baris tabel keputusan **42 → 48**. Catatan pengembang versi baru berbunyi **`"add avi"`**.

⭐ **Versi Fac In menerima lebih banyak jenis berkas.** ⛔ Perilaku yang sudah ada **tidak diubah**.
Apakah Claim Prop perlu ikut menerima jenis tambahan itu: **keputusan work owner**. `[terbuka]`

### 3d — Daftar mata uang: nol selisih logika

`[terverifikasi]` `…!BROWSECURRENCY_RD` — **69 medan berbeda dari 363**, tetapi **jumlah baris
penyaringnya sama: 18 lawan 18**. Hampir seluruh selisih adalah **peringatan mutu bawaan Pega**
*(`pxWarning*`)*: versi baru memikul **2** peringatan kinerja, versi lama **1**. ⛔ **Bukan selisih
perilaku.**

---

## §4 — Kalimat *"52 dari 61 rule `When`"*

### 4a — Di mana ia hidup: **10 kemunculan di 5 berkas**, bukan empat tempat

`[terverifikasi]`

| Berkas | Kemunculan | Sifatnya |
| --- | ---: | --- |
| `spec.md` | **3** | ⭐ **hidup** — perlu ikut berubah |
| `issues\04-klasifikasi-lini-bisnis.md` | **4** | ⭐ **hidup** — perlu ikut berubah |
| `grilling-ronde-2.md` *(sebagai **51 dari 60**)* | 1 | jejak sejarah — **biarkan** |
| `grilling-ronde-4.md` | 1 | jejak sejarah — **biarkan** |
| `grilling-ronde-6.md` | 1 | jejak sejarah — **biarkan** |

⛔ Brief menyebut *"empat tempat"*. Yang benar **10 kemunculan**, dan yang **perlu disunting 7**, di
**2 berkas**.

### 4b — ⭐ Keputusan A5-6 ternyata **menjawab butir `[terbuka]` yang sudah menggantung**

AC 72 memikul RALAT bertanggal 2026-09-18 yang berbunyi: *"`[terbuka]` **penggolongannya belum
diputuskan.** … Apakah ini penerapan aturan berdiri atau penyimpangan tersendiri diserahkan ke
**work owner**"*.

⭐ **Keputusan A5-6 Claim Fac In menjawabnya langsung** — *"rule dipisahkan dari KEHIDUPANNYA"*.
Dengan itu pertanyaan **"dimigrasikan atau tidak"** menjadi **salah pertanyaan**: rule-nya
dipindahkan **sekali**; yang berbeda per modul adalah **apakah halaman yang diujinya ada**.

⛔ **Saya tidak menutup butir itu** — penutupannya keputusan work owner.

### 4c — **USULAN** teks pengganti

> ⛔ **Usulan. Belum dipasang. Kalimat lamanya dikutip utuh, tidak dihapus.**

**Kalimat lama, `spec.md` AC 72:**

> *"Rule klasifikasi warisan yang menguji model modul lain **tidak dimigrasikan** — 52 dari 61 rule
> `When` (dulu **51 dari 60**) menguji properti yang tidak ada pada objek kerja Claim Prop."*

**Usulan pengganti:**

> `[terverifikasi]` Ke-61 rule klasifikasi lini **dipindahkan sekali sebagai satu himpunan**, dan
> **hidup-matinya tidak ikut dipindahkan**. Sebuah rule bernilai benar hanya bila halaman yang
> diujinya ada pada objek kerja modul yang memanggilnya. **Pada Claim Prop, 52 dari 61 menguji
> `BusinessType` pada halaman `OfferFacIn`/`Quotation` yang tidak ada di sini — sehingga di modul
> ini ke-52-nya tidak pernah bernilai benar.** Pada Claim Fac In halaman itu **ada**, dan rule yang
> sama **hidup**. ⭐ Karena itu ia **bukan** rule yang "tidak dimigrasikan", melainkan rule yang
> **mati di satu modul dan hidup di modul lain** — satu salinan, dua nasib.
>
> ⚠️ Klasifikasi lini di sistem baru tetap memakai **satu kunci `TreatyGroupID`** *(AC 71)*; ke-61
> rule itu **tidak dipakai untuk klasifikasi**, dan keberadaannya semata agar satu salinan logika
> tetap ada bagi modul yang memerlukannya.

⭐ **Dengan begitu kedua keputusan berdiri bersama:** Claim Prop tetap benar, Claim Fac In tetap
benar, dan sistem baru tetap punya **satu salinan**.

### 4d — Tempat yang ikut berubah bila usulan diterima

| # | Berkas | Yang menyesuaikan |
| --- | --- | --- |
| 1 | `spec.md` AC 72 | teks utama + **RALAT 2026-09-18 dinyatakan terjawab** |
| 2 | `spec.md` bab klasifikasi | kalimat *"52 dari 61 … sisa impor model Fac In"* |
| 3 | `spec.md` penutup | *"Kelimanya termasuk 52 rule sisa impor"* |
| 4–7 | `issues\04-…` **4 kemunculan** | termasuk label **`Alasan menyimpang:`** yang kini punya jawaban |

⚠️ **Dan satu akibat yang harus disebut:** bila usulan diterima, **penggolongan AC 72 berubah** dari
`[terbuka]` menjadi **penerapan aturan berdiri** — sehingga **cacah penyimpangan sadar Claim Prop
mungkin ikut berubah**. ⛔ Belum saya hitung ulang; itu pekerjaan blok berikutnya.

---

## §5 — Penutup

### 5a — RALAT yang DIBUTUHKAN — **empat**

| # | Berkas | Kalimat/angka LAMA, dikutip utuh | Usulan |
| --- | --- | --- | --- |
| **1** | *(lisan, ronde sebelum blok ini)* | *"`Flow_TreatyIn.xml` … cara-2 `pxObjClass` = **15**"* dan *"**5 `Event-Exception`** yang tidak pernah tercatat"* | ⭐ **bentuknya 5, penghubungnya 5**; angka 15 adalah **empat jenis entri dijumlahkan**. Perkecualiannya **1 tiket bernama + 4 cangkang kosong**, dan ⭐ **label mati** |
| **2** | `spec.md` AC 72 + `issues\04-…` | *"52 dari 61 rule `When` … **tidak dimigrasikan**"* | §4c — *"satu salinan, dua nasib"* |
| **3** | `grilling-ronde-2.md` | *"keunikan bergantung resolusi **milidetik** … **Dua unggahan dalam milidetik sama menghasilkan ImageID identik**"* | ⭐ **tetap benar dan kini dikuatkan**; ditambah: versi lebih baru rule yang **sama** sudah menambalnya |
| **4** | daftar utang `claim-facin\grilling-ronde-3/4` | *"**Baca ulang kesimpulan kurs standar Claim Prop** — rule yang dipakai di sana 4 tahun 3 bulan lebih tua"* | ⭐ **sudah dibaca: SQL-nya identik.** Kecemasannya **gugur** |

### 5b — Butir `[terbuka]` **BARU** — tiga

| # | Butir | Pemilik |
| --- | --- | --- |
| **1** | ⚠️ **`OperatorID.pyPosition != "IT Developer"`** — pintu belakang pengembang bergantung **teks jabatan**; ada di **kedua** modul | work owner + IAM |
| **2** | ⚠️ **ADR-0009 tidak pernah dijawab** — Claim Prop tidak memutuskan **migrasi penuh vs koeksistensi** | work owner |
| **3** | ⚠️ `Assignment1` ber-`pyBaseClass` **`…-Work-PNC`** di alur berkelas `…-Work-ClaimTreaty` — sisa impor **di tingkat alur** | work owner |

⛔ **Nol butir `[terbuka]` lama ditutup.**

### 5c — Yang blok ini **TIDAK** kerjakan

| # | Butir |
| --- | --- |
| 1 | **Membaca isi `AcceptanceClaim` sebagai rule** — hanya rujukannya yang terbaca; wujud rule-nya **tidak ada di ekspor** |
| 2 | **Menelusuri apakah `Double` benar-benar melewati aritmetika biner** — yang dibaca **deklarasi tipe**, bukan jalannya nilai |
| 3 | **Menghitung ulang cacah penyimpangan sadar** sesudah usulan §4c |
| 4 | **Menyisir 12 `FlowAction` lebih dalam** dari sensus bentuk — ⭐ **nol bentuk**, dan isinya tidak dibaca |
| 5 | **Menguji 158 rule "identik"** — diterima apa adanya dari ronde 1 Claim Fac In, **tidak dihitung ulang di sini** |

### 5d — ⚠️ Kesimpulan blok ini yang **PALING RAWAN SALAH**

⚠️⚠️ **§1d — vonis "tiket itu label mati".**

**Kenapa rawan:** ia berdiri di atas sebuah **ketiadaan**, dan ketiadaan hanya sekuat jendela yang
menyapunya. Jendela saya **409 berkas dua modul**. `AcceptanceClaim` terdaftar di kelas **induk
bersama `ASM-FW-GCNMFW-Work`**, dan pembangkit bakunya **terbukti ada di tiga modul lain**. ⛔ Satu
rule di luar ekspor ini sudah cukup membatalkan vonisnya.

**Yang kokoh dan tidak rawan:** ⭐ **sensus bentuk** — ia **menutup secara aritmetika** *(77 = 77)*,
dan tiap entri terbaca **wadah demi wadah**, bukan disimpulkan dari pola.

**Kedua paling rawan:** §2b vonis **ADR-0007 MENDUKUNG** — ia bersandar pada **cacah nama kolom
jejak**, dan saya **tidak membaca** apakah tiap transisi status benar-benar tercatat.

---

## Lampiran — bukti berkas lain tidak disentuh

```
korpus Claim Prop            329 berkas .xml   — nol dibuka untuk ditulis
korpus Komite Claim Prop      80 berkas .xml   — nol dibuka untuk ditulis
korpus Claim Fac In          482 berkas .xml   — nol dibuka untuk ditulis
docs\adr\                     15 berkas        — nol disunting
.scratch\claim-prop\          26 berkas lama   — nol disunting
.scratch\komite-claim-prop\   26 berkas lama   — nol disunting
```

⛔ **Nol kode Go/React · nol `CREATE TABLE` · nol DDL · nol nomor baris XML dikutip · nol nilai
rahasia · nol butir `[terbuka]` ditutup · nol ADR direvisi.**
