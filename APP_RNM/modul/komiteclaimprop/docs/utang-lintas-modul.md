# Utang lintas modul — **Komite Claim Prop**

> **Ronde TEMUAN, bukan ronde keputusan.** ⛔ Nol berkas lama disunting. ⛔ Nol butir `[terbuka]`
> ditutup. ⛔ Nol ADR direvisi. Semua usulan di bawah menunggu keputusan work owner.

**Sensus blok ini:** vonis ADR **15** · RALAT dibutuhkan **3** · butir `[terbuka]` baru **2** ·
usulan teks **1** · premis brief yang terbukti **keliru 3** · ⭐ **lubang wewenang 1**

**Jendela sensus:** seluruh **80** berkas `.xml` `Komite Claim Prop`, ditambah **329** berkas
`Claim Prop` dan **482** berkas `Claim Fac In` sebagai pembanding. Dokumen: **26** berkas
`.scratch\komite-claim-prop` dan **15** ADR.

---

## §1 — Sensus alur: **tiga premis brief runtuh, satu demi satu**

### 1a — ⛔ Premis pertama **KELIRU**: baris 307 spec ternyata **BENAR**

Brief ini lahir dari tuduhan saya sendiri bahwa `spec.md` baris 307 — *"Bentuknya **empat**,
penghubungnya **empat**"* — **salah**, karena sensus `pxObjClass` memberi **12**.

⛔ **Tuduhan itu salah. Baris 307 benar, kata demi kata.**

`[terverifikasi]` Angka **12** yang saya pakai untuk menuduh ternyata **empat jenis entri berbeda
yang dijumlahkan**: bentuk **4** + penghubung **4** + pengubah **1** + tiket **3**. Bentuk alur
tinggal di wadah **`pyShapes`**; `pyConnectors`, `pyModifiers`, dan `pyTicketShapes` adalah wadah
**lain**, dan `pyRouterProp`/`pyNotifyProp` adalah **milik** sebuah bentuk, bukan bentuk.

| | cara 1 `pyShapeType` | mentah `pxObjClass` | ⭐ **bentuk** | penghubung | pengubah | tiket |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `KomiteTreaty_Flow` | 6 | 17 | ⭐ **4** | ⭐ **4** | 1 | 3 |
| `Flow_TreatyIn` *(Claim Prop, pembanding)* | 8 | 22 | 5 | 5 | 1 | 4 |
| `Register_Flow` *(Claim Fac In, pembanding)* | 16 | 38 | 8 | 10 | 2 | 7 |

**Uji tutup:** `bentuk 17 + penghubung 19 + pengubah 4 + tiket 14 + perute 17 + pemberitahu 6 = 77`,
dan sapuan mentah `pxObjClass` juga **77**. Tidak ada entri yang hilang atau terhitung dua kali.

⭐ **Keempat bentuknya:** `Start1` *(mulai)* · `ASSIGNMENT63` *(**`KomiteRouter`**, `WorkList`,
`pyRouteTo = Custom`)* · `Decision1` *(**`KomiteLoop`**, `DataXOR`)* · `END52` *(selesai)*.

`[terverifikasi]` **`END52` ber-`pyWorkStatus` `Resolved-Completed`** — sejalan dengan **AC 5**
*(kasus selesai berstatus akhir selesai-tuntas)*.

### 1b — ⛔ Premis kedua **KELIRU**: keempat perkecualian itu **bukan empat**

`[terverifikasi]` Dibaca satu per satu:

| Menggantung di | Nama | Isi |
| --- | --- | --- |
| `END52` | *(tanpa nama)* | ⛔ **hanya `pxObjClass`** — cangkang kosong |
| `ASSIGNMENT63` | *(tanpa nama)* | ⛔ cangkang kosong |
| ⭐ `Decision1` | **`komiteAccept_ticket`** | tiket bernama |
| `pyModifiers` | **`komiteAccept_ticket`** | **definisi** tiket yang sama |

⭐ **Jadi tiketnya SATU, bukan empat** — `komiteAccept_ticket`, menggantung pada gerbang
**`KomiteLoop`**. Dua sisanya cangkang tanpa nama dan tanpa pengenal.

### 1c — ⛔⛔ Premis ketiga **KELIRU, dan ini yang paling telak**: tiket itu **sudah tercatat, sudah ditanyakan, dan sudah diputuskan**

Saya menuduh ada *"empat `Data-MO-Event-Exception` yang tidak pernah tercatat di mana pun"*.
⛔ **Salah.** Yang satu-satunya bermakna **sudah digarap tuntas tiga ronde lalu**:

| Kapan | Di mana | Apa |
| --- | --- | --- |
| ronde 6 | `grilling-ronde-6.md` | ⭐ **ditemukan** — *"ada satu ticket bernama `komiteAccept_ticket`"* |
| ronde 6 | pertanyaan **14** | ditanyakan ke work owner — *"Masih dipakai?"* |
| ronde 6 | butir **8** register | *"siapa pemicunya"* — `[terbuka]` |
| ⭐ ronde 7 §A1 | `[keputusan work owner]` **2026-09-18** | ✅ **dibuang — jalur mati** |
| sekarang | `spec.md` *Titik lompat darurat* | tertulis sebagai keputusan |
| sekarang | ⭐ **AC 7** | *"Titik lompat darurat `komiteAccept_ticket` **tidak dibuat**"* |

⛔ **Tidak ada lubang. Modul ini sudah menutupnya sendiri.**

### 1d — ⭐⭐ Tetapi blok ini **menambahkan bukti yang dulu tidak ada**

Ronde 6 menutup temuannya dengan kalimat jujur: *"⛔ **Siapa yang memicunya** [belum diketahui]"*,
dan butir **8** register berbunyi *"siapa pemicunya"*. Keputusan ronde 7 *"jalur mati"* diambil
**tanpa bukti korpus** — ia keputusan, bukan pembacaan.

⭐ **Sekarang buktinya ada.**

Tiket Pega adalah **sasaran lompatan**: siapa pun yang memanggil metode `Obj-Set-Tickets` dapat
melempar alur ke titik itu, **dari mana saja**.

`[terverifikasi]` **Jendela: seluruh 409 berkas `.xml` dua modul, SELURUH tag — bukan daftar medan
pilihan.**

| Yang dicari | Komite Claim Prop | Claim Prop |
| --- | --- | --- |
| metode langkah **`Obj-Set-Tickets`** | ⭐ **NOL** | ⭐ **NOL** |
| activity **`SetTicket`** *(pembangkit baku Pega)* | ⭐ **NOL** | ⭐ **NOL** |
| pemanggil `SetTicket` | ⭐ **NOL** | ⭐ **NOL** |
| `komiteAccept_ticket` di luar berkas alur | ⭐ **NOL** | — |

`[terverifikasi]` `komiteAccept_ticket` muncul di **tepat tiga berkas** pada **seluruh 20 korpus** —
`Komite Claim Prop`, `Komite Claim Non Prop`, `Komite Claim FacIn` — dan **ketiganya berkas alur itu
sendiri**.

⭐⭐ **Jadi jawaban atas pertanyaan ronde 6 adalah: TIDAK ADA YANG MEMICUNYA.** Keputusan work owner
*"jalur mati"* kini **bukan lagi asumsi, melainkan terbaca dari korpus**.

⭐ **Akibat halus tetapi berguna untuk AC 7:** karena Pega sendiri **tidak pernah melemparnya**,
tidak membangunnya adalah **PARITAS**, bukan penyimpangan. AC 7 bertanda `[keputusan work owner]`;
ia kini **berhak juga atas `[terverifikasi]`**. Lihat usulan §1f.

⚠️ **Batas kejujuran vonis ini, dan ia nyata:** mesin pembangkitnya **ada di aplikasi yang sama**.
`Activity\SetTicket.xml` hidup di **tiga modul lain** — `NB FacIn`, `RNW Fac In`,
`Endorsment Fac In` — dan ia rule **`@baseclass`** bawaan Pega. ⛔ **Tidak ada buktinya menyentuh
modul ini**, dan saya tidak menyimpulkan ada. `[terbuka]`

### 1e — ⭐ Vonis **AC 4 · 5 · 6 · 7**: **AC MASIH BERDIRI** — keempatnya

| AC | Isi | Vonis | Dasar |
| --- | --- | --- | --- |
| **4** | kasus baru **langsung** masuk kotak kerja penyetuju | ✅ **MASIH BERDIRI** | `Start1` → `ASSIGNMENT63` tanpa bentuk antara; **nol** bentuk lain di kanvas |
| **5** | kasus tanpa penyetuju berikutnya **selesai**, status akhir selesai-tuntas | ✅ **MASIH BERDIRI** | `END52` ber-`pyWorkStatus` **`Resolved-Completed`**; **satu** bentuk selesai |
| **6** | ⭐ **satu-satunya jalan kembali = tahap penyetuju yang sama** | ✅ **MASIH BERDIRI** | ⭐ **inilah yang brief ini curigai.** Keempat bentuk terbaca; **nol** jalur perkecualian hidup; tiketnya **label mati** *(§1d)* |
| **7** | `komiteAccept_ticket` **tidak dibuat** | ✅ **MASIH BERDIRI — dan menguat** | §1c dan §1d |

⭐⭐ **Kecemasan yang melahirkan brief ini — bahwa jalur perkecualian mungkin jalan kembali yang lain
— terbukti TIDAK BERDASAR.**

### 1f — **USULAN** untuk baris 307 dan AC 7

> ⛔ **Usulan. Belum dipasang. Tidak ada kalimat lama yang perlu diganti** — keduanya sudah benar.
> Yang diusulkan hanyalah **penambahan**.

**Baris 307 — tambahan sesudah kalimat yang sudah ada** *(kalimatnya sendiri tidak diubah)*:

> ⭐ `[terverifikasi]` Angka **empat/empat** itu dihitung dari wadah **`pyShapes`** dan
> **`pyConnectors`**. Berkas alurnya juga memuat **satu pengubah** dan **tiga tiket**, yang
> **bukan** bentuk: dua di antaranya cangkang tanpa nama, dan yang bernama adalah
> **`komiteAccept_ticket`** — sudah diputuskan **tidak dibuat** *(AC 7)*. ⚠️ Sensus yang menghitung
> seluruh entri `Data-MO-*` tanpa memisahkan wadahnya akan memberi **12**, dan angka itu
> **bukan jumlah bentuk**.

**AC 7 — usulan penguatan tanda golongan:**

> **7.** `[keputusan work owner]` **+ `[terverifikasi]`** Titik lompat darurat `komiteAccept_ticket`
> **tidak dibuat**. ⭐ **Dan tidak membangunnya adalah paritas, bukan penyimpangan**: sapuan seluruh
> 409 berkas dua modul menemukan **nol** pemanggil `Obj-Set-Tickets` dan **nol** pemakaian
> `SetTicket`, sehingga di Pega pun tiket itu **tidak pernah dilempar**. Test yang menemukan jalan
> masuk ke alur komite selain penyerahan dari kasus klaim **gagal**. *(Bab 2)*

---

## §2 — Lima belas ADR diadu dengan Komite Claim Prop

### 2a — Di sinilah utangnya nyata

`[terverifikasi]` Sementara Claim Prop sudah merujuk **9 dari 15** ADR, Komite Claim Prop merujuk
⭐ **SATU** — ADR-0003, **sekali**.

⭐⭐ **Empat belas ADR belum pernah disebut sama sekali di modul ini** — termasuk **ADR-0014 dan
ADR-0015**, yang justru **lahir dari modul komite saudaranya** *(Komite Claim Life)*. Untuk modul
ini, premis brief **benar**.

### 2b — Lima belas vonis

| ADR | Pokok | Vonis Komite Claim Prop | Bukti |
| --- | --- | --- | --- |
| **0001** | batas konteks lewat *child work* | **TIDAK MENYENTUH** | `pxAddChildWork`: **NOL** — modul ini **anak**, bukan induk. Batasnya ditegakkan dari sisi Claim Prop |
| **0002** | RBAC tiga peran | **TIDAK MENYENTUH** ⚠️ | ⭐ lihat 2c — **jebakan nama** |
| **0003** | uang bukan `float` | ✅ **MENDUKUNG** ⚠️ | `Double`/`Float`: ⭐ **NOL** di 80 berkas. ⚠️ **tetapi pintu masuknya tidak bersih** — lihat 2d |
| **0004** | alamat jadi env var | *(tidak dinilai — berstatus `superseded`, digantikan 0013)* | — |
| **0005** | `IsPEGAPROD` jadi flag lingkungan | ✅ **MENDUKUNG** | `IsPEGAPROD` di **10 berkas**, a.l. `HitServiceToKasirKMT_Act`, `KomitePostAdjustment`, `KomitePost_Close`, `KomitePost_Reject` |
| **0006** | penomoran lewat *stored procedure* | ✅ **MENDUKUNG** | `GetSequenceNumber_SQL` di **2 berkas**, dipanggil `KomitePostAdjustment` |
| **0007** | jejak audit tiap transisi + jalur balik | ✅ **MENDUKUNG** ⚠️ | jejak operator **97 kemunculan / 11 berkas** — tetapi ⚠️ butir `[terbuka]` **1** modul ini sudah mencatat *"kolom akun operator di riwayat akseptasi **tidak pernah diisi** jalur komite"*. ⭐ **ADR-0007 memperkuat butir itu** |
| **0008** | efek keluar asinkron, **tidak memblokir** | ⛔ **MENENTANG — sengaja** | `IsSuccessHitService` / `StsKonversi`: ⭐ **NOL**, sama seperti Life. Tetapi ADR-0015 **sengaja menyimpang** darinya untuk konteks komite, dan modul ini mengikuti 0015 |
| **0009** | migrasi penuh, tanpa koeksistensi | ⚠️ **belum punya data** | modul ini **banyak** memutuskan tentang **data lama** *(keputusan 75, butir 14, tiket 13)*, tetapi ⛔ **tidak pernah memutuskan penuh vs koeksistensi** |
| **0010** | berkas tetap di Google Storage | ✅ **MENDUKUNG** ⚠️ | **396 kemunculan / 12 berkas** — **sama persis** dengan Claim Prop. ⚠️ lihat §3 |
| **0011** | unit status = baris `AdjustmentList` | ✅ **MENDUKUNG** | `AdjustmentList` **114 / 8 berkas**; sejalan AC 1 *(satu kasus komite ↔ satu baris penyesuaian)* |
| **0012** | wewenang kirim komite bergantung `Type` | **TIDAK MENYENTUH** | ⭐ `.Type` = `TP`/`TR`/`QP`/`QR`: **NOL** — sama seperti Claim Prop dan Claim Fac In |
| **0013** | alamat di-*lookup* dari `M_LINK_SERVICE` | ✅ **MENDUKUNG** | `GetLinkService` + kunci kategori, **10 berkas** |
| **0014** | ⭐⭐ keputusan komite ditegakkan per `KomiteID` tingkat berjalan | ⛔⛔ **MENENTANG** | ⭐ lihat 2e — **temuan terberat blok ini** |
| **0015** | efek keluar komite **wajib berhasil** | ✅ **MENDUKUNG** | `KomitePostAdjustment` + `KomitePost` ada; delapan efek keluar sudah dipetakan modul ini |

**Rekap:** ✅ MENDUKUNG **8** · ⛔ MENENTANG **2** · TIDAK MENYENTUH **3** · belum punya data **1** ·
tidak dinilai **1** *(superseded)*.

### 2c — ⭐ ADR-0002: **jebakan nama `pyPosition`**

`pyPosition` muncul **33 kali / 7 berkas** — terkesan ada pemeriksaan peran. ⛔ **Keliru.**

`[terverifikasi]` Menurut jenis rule dan nilai pembandingnya:

| Nilai | Berapa | Di mana | Artinya |
| --- | ---: | --- | --- |
| `Top` | **8** | ReportDefinition · Section | ⛔ **tata letak layar** |
| `AFTER` | **7** | ReportDefinition · Section | ⛔ **tata letak layar** |
| ⭐ `IT Developer` | **2** | `DataTransform\InsertChronology_DT` | ⭐ **wewenang sejati** |

⭐ **Nama peran gaya Life — `ReasLifeAdmin` dan kerabatnya — NOL.**

⚠️⚠️ Satu-satunya pemeriksaan wewenang sejati berbunyi **`OperatorID.pyPosition != "IT Developer"`**
— **pintu belakang pengembang** yang digantungkan pada **teks jabatan**, mengatur penghapusan
catatan kronologi. ⭐ Rule yang **sama persis** ada di Claim Prop. ⛔ Belum pernah tercatat di
kedua modul. `[terbuka]`

### 2d — ADR-0003: bersih di dalam, **tidak bersih di pintu masuk**

`[terverifikasi]` **Nol** parameter bertipe `Double` atau `Float` di seluruh **80** berkas modul
ini. ✅ Sejalan ADR-0003, dan sejalan keputusan modul bahwa uang bertipe desimal.

⚠️⚠️ **Tetapi uang MASUK ke modul ini lewat gerbang `Double`.** `Claim Prop\Activity\
AddKomiteTreatyChild_ACT` — **rule yang melahirkan kasus komite** — mendeklarasikan **`TSISpread`**
dan **`ClaimSpread`** bertipe **`Double`**, dan ⭐ **kedua sumber tipe sepakat** *(tanda tangan dan
medan parameter)*. `Claim Prop\Activity\SetKomiteTreaty_ACT` melakukan hal sama untuk
**`TotalAdjustment`**.

⭐ **Jadi kebersihan modul ini bergantung pada modul di seberangnya.** ⛔ Perlakuannya belum
diputuskan. `[terbuka]` — dan rincian penuhnya ada di berkas kembar `.scratch\claim-prop\
utang-lintas-modul.md` §2d.

### 2e — ⭐⭐⭐ ADR-0014: **MENENTANG** — lubang wewenang di modul yang sudah dinyatakan tuntas

**Yang ADR-0014 wajibkan:** pada tiap tingkat tangga, **hanya pemilik
`KomiteList(KomiteCount).KomiteID`** yang boleh **menyimpan keputusan**; pengguna lain **ditolak di
lapisan layanan**, meskipun dapat membuka kasusnya. ADR itu menyebutnya **penyimpangan sadar** —
karena *"Pega hanya **menempatkan** tugas, ia tidak **menegakkan** siapa yang memutuskan."*

**Keadaan Pega di modul ini — sama persis dengan yang ADR-0014 gambarkan:**

`[terverifikasi]` `KomiteID` muncul di **satu berkas saja**, `Activity\KomiteRouter` — dan di sana
ia dipakai untuk **menetapkan tujuan rute**, bergerbang *"keputusan penyetuju masih kosong"*. Itu
**penempatan**. Alurnya memakai `WorkList` dengan `pyRouteTo = Custom` — juga penempatan.
`AcceptStatus` ditulis oleh **tepat satu** penugasan properti di seluruh 80 berkas.

⭐ **Nol pemeriksaan bahwa pengguna yang menyimpan keputusan adalah pemilik `KomiteID` tingkat
berjalan.** Lubang yang sama, persis.

**Keadaan spec modul ini — dan inilah pertentangannya:**

`[terverifikasi]` Sapuan seluruh bab AC dan register: ⛔ **tidak ada satu pun AC** yang mewajibkan
penegakan itu.

| AC | Apa yang ia jamin | Apakah ia menutup lubangnya? |
| --- | --- | --- |
| **10** | kotak kerja **hanya** memuat kasus yang sedang gilirannya | ⛔ **tidak** — itu **penempatan** |
| **11** | tujuan rute diambil dari **akun operator** penyetuju giliran | ⛔ **tidak** — itu **penempatan** |
| **25** | tanggal dan identitas penyetuju **terisi otomatis** saat ia memutuskan | ⛔ **tidak** — itu **pengisian**, bukan pemeriksaan |
| ⭐ **67** | *"Wewenang di sistem baru ditentukan **aturan peran sistem baru**; penegakan **giliran** bukan penegakan **wewenang**"* | ⛔⛔ **justru sebaliknya** — ia **menyerahkan** wewenang ke aturan yang **belum ditulis di mana pun** |

⭐⭐ **AC 67 membaca gejalanya dengan benar** — ia sadar betul bahwa giliran bukan wewenang. ⛔
**Tetapi ia berhenti di situ**, dan menyerahkan penegakannya kepada *"aturan peran sistem baru"*
yang **tidak ada AC-nya, tidak ada tiketnya, dan tidak menyebut `KomiteID`**.

⚠️⚠️ **Akibatnya tegas:** modul ini **dapat dibangun lengkap, lulus seluruh 80 AC-nya, dan tetap
membiarkan siapa pun yang bisa membuka layar menyimpan keputusan untuk tingkat mana pun.** Itu
persis keadaan yang ADR-0014 dibuat untuk mencegah.

⛔ **Yang TIDAK saya lakukan:** menulis AC baru, mengubah AC 67, atau menyatakan ini penyimpangan
sadar. ⭐ **Ini butir `[terbuka]` untuk work owner**, dan menurut saya **butir pemblokir** — karena
ia menyentuh siapa yang boleh memutuskan uang. Penggolongannya tetap keputusanmu.

---

## §3 — Empat rule beda versi: modul ini memakai **versi yang sama tuanya dengan Claim Prop**

`[terverifikasi]` Untuk keempat rule yang berbeda versi antara Claim Prop dan Claim Fac In,
**Komite Claim Prop selalu memakai salinan yang sama dengan Claim Prop** — cap waktu, versi
ruleset, dan ukuran berkasnya **sama persis**.

⚠️ **Cara membacanya.** `md5` mentah **selalu berbeda** — bahkan antara Claim Prop dan Komite yang
ukuran dan cap waktunya identik. Sebabnya metadata ekspor. ⛔ **`md5` berkas BUKAN alat yang sah**
untuk pertanyaan ini. Yang dipakai: **perbandingan medan demi medan dengan medan-waktu dan
medan-ekspor dibuang.**

### 3a — ⭐⭐ Kurs standar: **isinya IDENTIK — kecemasannya GUGUR**

`[terverifikasi]` Claim Prop lawan Komite Claim Prop: ⭐ **hanya 2 medan berbeda dari 74**, dan
**keduanya cap waktu** *(`pyRuleFormStatusTime`, `pyShowJavaWindowName`)*. Claim Prop lawan Claim
Fac In: **12 medan berbeda**, **nol** di antaranya logika.

⭐ **Perintah SQL-nya sama kata demi kata di ketiga modul** — memanggil fungsi basis data
`getcurrencystandard` dengan mata uang dan `sysdate`, menghasilkan satu kolom `nilaiKurs`.

⭐⭐ **VONIS: kesimpulan kurs standar TIDAK BERUBAH — untuk kedua modul.** Selisih 4 tahun 3 bulan
itu **selisih label versi ruleset**, bukan selisih perilaku.

### 3b — ⭐⭐ Pembuat pengenal berkas: **cacat NYATA, dan modul ini ikut memikulnya**

`[terverifikasi]` `…!RNM!GENERATEIMAGEID_SQL` — **di sini isinya benar-benar berbeda**:

| | Bahan yang di-*hash* |
| --- | --- |
| **Komite Claim Prop** *(dan Claim Prop — sama persis)* | awalan tetap **+ cap waktu berketelitian MILIDETIK** |
| ⭐ **Claim Fac In** *(lebih baru)* | awalan tetap **+ cap waktu berketelitian NANODETIK + nilai unik sejagat** |

⭐⭐ **Bentuknya persis sebuah TAMBALAN TABRAKAN.** Versi yang dipakai modul ini **dapat
menghasilkan pengenal berkas kembar** bila dua unggahan jatuh dalam milidetik yang sama.

⚠️ **Modul ini memakai rule itu lewat `InsertGoogleStorage_Act`** — dan `RELASI-TABEL` modul ini
sudah mendaftarkannya sebagai salah satu **12 rule berpasangan** dengan Claim Prop. ⛔ **Cacat
tabrakannya belum pernah tercatat di dokumen modul ini** *(ia tercatat di `claim-prop\
grilling-ronde-2`)*. `[terbuka]`

### 3c — Jenis berkas dan daftar mata uang

`[terverifikasi]` `…!GETMIMETYPE` — versi Fac In memuat **48** baris tabel keputusan lawan **42**;
catatan pengembangnya berbunyi **`"add avi"`**. ⭐ **Menerima lebih banyak jenis berkas**, bukan
mengubah perilaku yang ada.

`[terverifikasi]` `…!BROWSECURRENCY_RD` — **jumlah baris penyaring sama: 18 lawan 18**. Selisihnya
hampir seluruhnya **peringatan mutu bawaan Pega**. ⛔ **Bukan selisih perilaku.**

---

## §4 — Kalimat *"52 dari 61 rule `When`"*

`[terverifikasi]` **Nol kemunculan** di seluruh **26** berkas `.scratch\komite-claim-prop`.
Kalimat itu **milik Claim Prop**, dan usulannya ada di berkas kembar §4c.

⚠️ **Satu hal yang menyentuh modul ini:** usulan itu bertumpu pada keputusan A5-6 Claim Fac In
*"rule dipisahkan dari kehidupannya"*. Bila diterima, prinsip yang sama berlaku bagi **setiap rule
bersama** yang dipakai modul ini — dan modul ini **berbagi 12 rule penyimpanan berkas** dengan
Claim Prop. ⛔ Belum ada yang perlu diubah di sini sekarang.

---

## §5 — Penutup

### 5a — RALAT yang DIBUTUHKAN — **tiga**

| # | Berkas | Kalimat/angka LAMA, dikutip utuh | Usulan |
| --- | --- | --- | --- |
| **1** | *(lisan, ronde sebelum blok ini)* | *"`komite-claim-prop\spec.md` baris 307 … **salah** … true census is **12 Data-MO entries** including **4 `Data-MO-Event-Exception` never mentioned anywhere**"* | ⛔⛔ **Ralat penuh: baris 307 BENAR.** Bentuk **4**, penghubung **4**. Angka 12 adalah **empat jenis entri dijumlahkan**. Dan perkecualian yang bermakna **sudah tercatat, ditanyakan, dan diputuskan** — ronde 6 → ronde 7 §A1 → **AC 7** |
| **2** | `spec.md` **AC 7** | *"`[keputusan work owner]` Titik lompat darurat `komiteAccept_ticket` **tidak dibuat**"* | ⭐ **tetap benar**; usulan §1f **menambah `[terverifikasi]`** — ia **paritas**, bukan penyimpangan, karena Pega pun tak pernah melemparnya |
| **3** | daftar utang `claim-facin\grilling-ronde-3/4` | *"**Baca ulang kesimpulan kurs standar Claim Prop** — rule yang dipakai di sana 4 tahun 3 bulan lebih tua"* | ⭐ **sudah dibaca: SQL-nya identik.** Kecemasannya **gugur**, untuk kedua modul |

### 5b — Butir `[terbuka]` **BARU** — dua

| # | Butir | Pemilik | Menahan? |
| --- | --- | --- | --- |
| ⭐⭐ **1** | **ADR-0014 tidak ditegakkan** — tidak ada AC yang mewajibkan hanya pemilik `KomiteID` tingkat berjalan boleh menyimpan keputusan; **AC 67 menyerahkannya ke aturan peran yang belum ditulis** | work owner + IAM | ⚠️ **menurut saya YA** — ia menyentuh siapa boleh memutuskan uang |
| **2** | ⚠️ **`OperatorID.pyPosition != "IT Developer"`** — pintu belakang pengembang bergantung **teks jabatan**; ada di **kedua** modul | work owner + IAM | belum dinilai |

⛔ **Nol butir `[terbuka]` lama ditutup.** ⚠️ Butir **8** ronde 6 *(siapa pemicu tiket)* sudah
ditutup ronde 7; blok ini **hanya menambahkan buktinya**, tidak menutup ulang.

### 5c — Yang blok ini **TIDAK** kerjakan

| # | Butir |
| --- | --- |
| 1 | **Menulis AC penegakan wewenang** untuk lubang §2e — itu keputusan work owner |
| 2 | **Menyisir 2 `FlowAction`** lebih dalam dari sensus bentuk — ⭐ **nol bentuk**, isinya tidak dibaca |
| 3 | **Membaca kedelapan efek keluar** terhadap tuntutan ADR-0015 satu per satu — hanya keberadaannya yang diuji |
| 4 | **Memeriksa kecocokan `END52`** antara modul ini dan Claim Fac In — keduanya memakai pengenal yang sama, dan itu **belum ditelusuri** |
| 5 | **Menguji 158 rule "identik"** — diterima apa adanya dari ronde 1 Claim Fac In |

### 5d — ⚠️ Kesimpulan blok ini yang **PALING RAWAN SALAH**

⚠️⚠️ **§2e — vonis ADR-0014 MENENTANG.**

**Kenapa rawan:** ia berdiri di atas **ketiadaan sebuah AC**, dan saya menyimpulkannya dari sapuan
kata kunci atas bab AC dan register — bukan dari membaca ke-80 AC satu per satu. ⛔ **Satu AC yang
menyatakan hal itu dengan kata lain sudah cukup membatalkannya.** Kata yang saya sapu: *menyimpan
keputusan · keputusan penyetuju · boleh memutus · hanya pemilik · KomiteID · identitas penyetuju ·
akun operator · lapisan layanan · wewenang · giliran*.

⚠️ **Rawan kedua, dan searah:** saya menyebutnya **pemblokir**. Itu **penilaian saya**, bukan
pembacaan — dan penggolongan pemblokir selama ini selalu keputusan work owner.

**Yang kokoh dan tidak rawan:** ⭐ **sensus bentuk** — ia **menutup secara aritmetika** *(77 = 77)*,
dan tiap entri terbaca **wadah demi wadah**. ⭐ Dan **§1c** — bahwa tiket itu sudah digarap tuntas —
terbaca langsung dari empat berkas dokumen, bukan disimpulkan.

---

## Lampiran — bukti berkas lain tidak disentuh

```
korpus Komite Claim Prop      80 berkas .xml   — nol dibuka untuk ditulis
korpus Claim Prop            329 berkas .xml   — nol dibuka untuk ditulis
korpus Claim Fac In          482 berkas .xml   — nol dibuka untuk ditulis
docs\adr\                     15 berkas        — nol disunting
.scratch\komite-claim-prop\   26 berkas lama   — nol disunting
.scratch\claim-prop\          26 berkas lama   — nol disunting
```

⛔ **Nol kode Go/React · nol `CREATE TABLE` · nol DDL · nol nomor baris XML dikutip · nol nilai
rahasia · nol butir `[terbuka]` ditutup · nol ADR direvisi.**
