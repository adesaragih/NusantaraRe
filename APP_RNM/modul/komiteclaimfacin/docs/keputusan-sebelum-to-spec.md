# Komite Claim Fac In — Keputusan sebelum `to-spec`
## Empat keputusan work owner, satu sensus yang saya perbesar sendiri, dan satu tabrakan dengan ADR-0014

⛔ **Berkas ini MENCATAT.** Ia tidak menulis spec, tidak membuat tiket, tidak membuat DDL.
⛔ **Satu butir `[terbuka]` ditutup di sini — butir 24 — dan hanya karena keputusan work owner
memang menutupnya.**

> **SENSUS BERKAS INI** — *(dihitung ulang sesudah RALAT §2 + K10–K13, 2026-09-20)*
>
> ⛔ **Jendelanya disebut lebih dulu: seluruh berkas ini DIKURANGI blok sensus ini sendiri**
> — **964 baris** menurut cacah baris-baru, **965** menurut pemisahan teks. ⚠️ **Selisih satu,
> sebabnya diketahui:** berkas berakhir dengan baris baru, sehingga pemisahan teks menghasilkan satu
> baris kosong penutup yang bukan baris sungguhan. ⭐ **Angka yang dipakai: 964.**
> ⭐ Blok sensus dikecualikan dengan sengaja — bila ia menghitung dirinya sendiri, angkanya
> berkejaran tanpa henti. ⭐ **Cara yang sama dipakai ronde 3 dan 4.**
>
> `[terverifikasi]` **27** · `[keputusan work owner]` **8** · `[terbuka]` **8** ·
> `[penyimpangan sadar]` **3** · `[data DBA]` **1** · `[data work owner]` **1**
> `[work owner]` **8** · `[pengembang Pega lama]` **1** · `[tim operasi]` **1**
> ⚠️ **60** · ⭐ **174** · ⛔ **170** · ✅ **38** · ❓ **5** · RALAT **14** · tabel **26**
>
> ⭐ **Dihitung DUA CARA:** *(a)* cacah penanda atas badan berkas sebagai satu teks utuh;
> *(b)* cacah ulang **baris demi baris**. ✅ **Keduanya sepakat pada keenam belas angka.**
>
> ⚠️ **Tiga catatan jujur.** *(1)* Penanda `[terbuka]` terhitung **8**, sedangkan register §4
> mencatat **14 butir terbuka** — ⭐ jendelanya berbeda: yang satu **penanda tertulis**, yang lain
> **butir dalam register**. *(2)* ⛔ **Angka BERGESER dua kali hari ini.** Sebelum RALAT pertama:
> `[terverifikasi]` 16 · ⚠️ 45 · ⭐ 77 · ⛔ 99 · RALAT 1 · badan **648**. Sesudah RALAT pertama:
> `[terverifikasi]` 19 · ⚠️ 49 · ⭐ 85 · ⛔ 110 · RALAT 3 · badan **687**. ⭐ Keduanya dikutip,
> tidak dihapus. *(3)* ⭐ **RALAT melonjak 3 → 14** — itu ukuran seberapa besar §2 harus dibetulkan.

---

## LANGKAH 0 — gerbang masuk

⭐ **Seluruh pencarian dijalankan atas teks yang dinormalkan** — penanda `* _ \` # > [ ]` dibuang,
spasi dan pergantian baris diciutkan jadi satu spasi.

| # | Yang diperiksa | Harap | Dapat | |
| --- | --- | ---: | ---: | --- |
| 0a | berkas `.md` di folder | 4 | **4** | ✅ |
| 0b | baris `grilling-ronde-4.md` | 810 | **810** | ✅ |
| 0c | cacah `## ` | 9 | **9** | ✅ |
| 0d | `14 butir terbuka` *(dinormalkan)* | 2 | **2** | ✅ |
| 0e | `Yang MEMBLOKIR: 1 · 6 · 26` *(dinormalkan)* | 1 | **1** | ✅ |
| 0f | `Monika` | 0 | **0** | ✅ |
| 0g | `ADR-0014` | 1 | **1** | ✅ |
| 0h | `spec.md` ada? | TIDAK ADA | **TIDAK ADA** | ✅ |

⭐ **Delapan gerbang lolos.** ⚠️ Perhatikan 0f: `grilling-ronde-4.md` memang **nol** menyebut nama
orang — ⭐ itu disengaja sejak ronde 3, dan berkas **ini** adalah tempat pertama nama-nama itu
ditulis utuh, atas perintah eksplisit work owner *(§2a)*.

---

## §1 — Empat keputusan work owner

> ⛔ **RALAT terhadap prompt yang memerintahkan berkas ini.** Prompt berbunyi *"T1 — Catat **lima**
> keputusan work owner"*, ⛔ **tetapi hanya EMPAT yang diberikan** — Q1, Q2, Q4, Q5 — dan prompt
> yang sama **menahan Q3 secara eksplisit** lewat kalimatnya sendiri: *"Q3 ⛔ BELUM DIPUTUSKAN …
> ⛔ JANGAN memilih sendiri."*
>
> ⚠️ Judul berkas dan judul §1 sempat berbunyi *"**Lima** keputusan work owner"* — ⭐ **yang benar
> EMPAT, ditambah satu yang sengaja ditahan.** ⭐ Kalimat lamanya dikutip di sini, tidak dihapus.
>
> ⭐ **Cacah pembuktinya:** bab keputusan di bawah ini berisi **empat** — K5 · K6 · K9 · K8 —
> ditambah bab Q3 yang **tidak memuat keputusan**. ⛔ **Isinya sudah benar sejak awal; yang salah
> hanya judulnya.**
>
> ⚠️ **Jangan memakai cacah penanda di kepala sensus untuk memeriksa ini.** Kepala sensus
> menghitung **penanda yang tertulis**, sedangkan yang diralat di sini adalah **cacah BAB keputusan**.
> ⛔ Kedua jendela itu **tidak wajib sama**, dan menyamakannya pernah membuat saya menulis kalimat
> yang membatalkan dirinya sendiri — ⭐ kalimat itu dibuang sebelum berkas ini selesai, dan
> kejadiannya dicatat di sini.

Semuanya `[keputusan work owner]` **2026-09-20**. ⛔ **Ditulis apa adanya, tidak ditafsirkan ulang,
tidak dihaluskan.**

### K5 → Q1 · penjaga ganda-bayar

> ⭐ **Penjaga ganda-bayar dibangun EKSPLISIT sekarang. Penanda tersimpan dalam transaksi yang sama
> dengan pengiriman pembayaran; panggilan kedua DITOLAK.**

✅ **Sejalan rekomendasi asisten.**

⭐ **Akibat pada register:** butir **26** — *tanda `=` tunggal pada penjaga ganda-bayar* —
⭐ **TIDAK LAGI MEMBLOKIR.** ⛔ Ia **tidak ditutup**: sifatnya **berubah**, dari **keputusan
rancangan** menjadi **pemeriksaan DATA LAMA**. Rancangan sistem baru tidak lagi bergantung pada
jawabannya; yang masih bergantung hanyalah pertanyaan *"apakah sistem lama pernah membayar dua
kali"*.

⚠️ **Pertanyaan ke pengembang Pega lama tetap dicari, dan DIPECAH DUA:**

| | Pertanyaan |
| --- | --- |
| **1a** | Apakah tanda `=` **tunggal** dibaca sebagai **PEMBANDINGAN** di dalam ekspresi `when` |
| ⭐ **1b** | Apa arti **kode 2, 3, dan 6** pada `pyStepsPreCondParamsWhenTrue` / `…WhenFalse` |

⚠️⚠️ **1b lebih menentukan daripada 1a**, dan ini wajib dibaca pelan: ⛔ **selama 1b belum
terjawab, SETIAP vonis "gerbang mati" di modul mana pun berstatus `[terbuka]`.** ⭐ Empat ronde
grilling modul ini bersandar pada pembacaan kode arah itu.

### K6 → Q2 · surel galat kasir

> ⭐ **Pemberitahuan galat kasir dikirim HANYA pada kegagalan. Perilaku lama TIDAK ditiru.**

✅ **Sejalan rekomendasi asisten.** ⚠️ **`[penyimpangan sadar]`** — `[terverifikasi]` ronde 4 §3c:
sistem lama mengirimnya **setiap kali**, sebab gerbang `.StatusServiceKasir.ReponseCode != "1"`
ber-bendera `false`.

⭐ **Cacah surel galat yang benar-benar sampai = `[data work owner]`**, dicari dari kotak masuk
dengan subjek **`"Warning!, error Direct to Kasir"`**.

### K9 → Q4 · larangan menyetujui klaim yang diajukan sendiri

> ⭐ **TIRU PEGA APA ADANYA.**
> **Fac In:** pemeriksaan **hanya pada anggota komite urutan PERTAMA.**
> **Prop:** ⛔ **NOL pemeriksaan — tidak di aktivitas, tidak di layar.**

⛔ **BERBEDA dari rekomendasi asisten**, yang mengusulkan memberlakukannya untuk semua anggota di
kedua modul. ⛔ **Disengaja.**

⚠️ **Akibatnya ditulis terang-terangan, sekali, tanpa dihaluskan:**

> ⛔ **Anggota komite di urutan KEDUA ke bawah tetap dapat menyetujui klaim yang ia ajukan sendiri
> di Fac In. Di Prop, di urutan MANA PUN. Ini lubang yang dibawa masuk DENGAN SADAR.**

⭐ **Butir 24 DITUTUP oleh keputusan ini.** Kalimat penutupannya:

> ⭐ **Butir 24 ditutup 2026-09-20 oleh `[keputusan work owner]` Q4: larangan menyetujui klaim yang
> diajukan sendiri ditiru apa adanya — anggota pertama saja di Fac In, nihil di Prop — dan lubang
> pada urutan kedua ke bawah diterima dengan sadar.**

### K8 → Q5 · data lama MBU dan Travel

> ⭐ **BIARKAN, tetapi TANDAI.** Ringkasan akseptasi historis MBU dan Travel **tidak dibangun
> ulang**.

✅ **Sejalan rekomendasi asisten.**

⭐⭐ **Tandanya KOLOM, bukan catatan prosa** — mengikuti keputusan penanda 2026-09-19:

| | |
| --- | --- |
| nama kolom | **`DIBANGUN_ATURAN_LAMA`** |
| jenis | `CHAR(1)` |
| nilai | `'1'` / `'0'` |

⭐ **Alasannya dinyatakan:** supaya **laporan dapat menyaringnya**, bukan hanya pembaca manusia yang
kebetulan membaca catatan kaki.

⭐ **Cacah baris terdampak = `[data DBA]`.**

### Q3 → ⛔ **BELUM DIPUTUSKAN**

⛔ **Jawaban work owner bermakna ganda, dan saya TIDAK memilih sendiri.**
⭐ Diajukan utuh sebagai **Pertanyaan A** di bab **PERTANYAAN BARU**.
⛔ **Tidak ada keputusan Q3 yang ditulis di berkas ini.**

---

## §2 — Nama orang ditanam keras di gerbang

> ⛔⛔⛔ **RALAT BESAR — 2026-09-20. TEMUAN BAB INI MENGECIL, BUKAN MEMBESAR.**
>
> ⛔ **Seluruh §2 ini terlalu besar pada TIGA hal. Ketiganya diralat di bawah; ⭐ kalimat lamanya
> TIDAK dihapus, hanya dikutip berdampingan dengan koreksinya.**
>
> ---
>
> ⭐⭐ **ATURAN BACA BARU — sebab ketiganya.**
>
> `[terverifikasi]` Medan `pyExpression` yang dicacah bab ini **duduk di dalam** elemen
> **`pyExpressionGadget`** yang ber-`pxObjClass` **`PegaGadget-ExpressionBuilder`**.
> ⛔ **Itu keadaan kotak dialog *expression builder* Pega — apa yang TERAKHIR DIKETIK pengembang
> ke dalam editor — BUKAN nilai yang dieksekusi.**
>
> ⭐ Nilai yang **benar-benar berjalan** ada di **`PropertiesValue`** pada `rowdata` di dalam
> `pyParamArray` sebuah langkah, atau di **`pyStepsPreCondParamsWhen`**.
>
> > ⭐⭐ **ATURAN.** Elemen di bawah `pyExpressionGadget` / `PegaGadget-ExpressionBuilder` adalah
> > **SISA EDITOR**, bukan rule yang berjalan. ⛔ **Wajib dikeluarkan dari setiap sensus gerbang.**
> > ⚠️ Sensus yang tidak memisahkannya **melebihkan cacah gerbang sampai lebih dari dua kali
> > lipat** — itu yang terjadi pada bab ini.
>
> ---
>
> **RALAT SATU — angkanya.** `[terverifikasi]` sensus ulang atas **1.005 berkas**, dipisah menurut
> aturan di atas:
>
> | | §2 *(salah)* | ⭐ Yang benar |
> | --- | ---: | ---: |
> | gerbang bernama-orang **hidup** | 21 | ⭐ **9** |
> | di antaranya `pyMemo` *(catatan pengembang, tidak berjalan)* | — | **1** |
> | ⭐ **benar-benar dieksekusi** | — | ⭐ **8** |
> | **sisa editor** *(dikeluarkan)* | 0 | ⭐ **24** |
> | berkas ber-gerbang hidup | 10 | ⭐ **5** |
> | modul | 4 | **4** *(tetap)* |
> | akun unik | 6 | **6** *(tetap)* |
>
> **Kelima berkas yang gerbangnya HIDUP:**
>
> | Modul | Berkas | Tag | Cacah |
> | --- | --- | --- | ---: |
> | Claim Fac In | `Activity\SendEmailKlaimRejectClose.xml` | `pyStepsPreCondParamsWhen` | 2 |
> | Claim Prop | `Activity\SendEmailKlaimRejectClose.xml` | `pyStepsPreCondParamsWhen` | 2 |
> | ⭐ Komite Claim FacIn | `Activity\SetProteksiSubmiteKomite.xml` | `PropertiesValue` **1** + `pyMemo` **1** | 2 |
> | ⛔ Komite Claim Prop | `Activity\KomitePost_Close.xml` | `PropertiesValue` | 2 |
> | ⛔ Komite Claim Prop | `Activity\KomitePost_Reject.xml` | `PropertiesValue` | 2 |
>
> ---
>
> **RALAT DUA — Komite Claim FacIn TIDAK punya tabel jabatan hidup.**
> ⛔ Bab ini menuduh `KomitePost_Adjustment`, `KomitePost_CloseClaim`, dan `KomitePost_Reject`
> milik Komite Claim FacIn membawa tabel jabatan ditanam keras. ⭐ **Itu SALAH — ketiganya sisa
> editor.**
>
> `[terverifikasi]` **Jalur hidupnya memakai ROSTER, bukan nama:**
>
> | Modul | Rule | Nilai yang berjalan |
> | --- | --- | --- |
> | Komite Claim FacIn | `KomitePost_Adjustment` · `_CloseClaim` · `_Reject` | `Data.CARI12 = "Accepted by " + pyWorkPage.KomiteList(pyWorkPage.KomiteCount).IDKomite` |
> | Komite Claim Prop | `KomitePostAdjustment` lgk **8 · 9** | `DataChronology.CARI1 = "Accepted by " + …KomiteList(KomiteCount).IDKomite` |
>
> ⭐ **Tabel jabatan yang HIDUP hanya ada di DUA rule, keduanya milik Komite Claim Prop:**
> `KomitePost_Close` lgk **2 · 3** dan `KomitePost_Reject` lgk **2 · 3**, bergerbang
> `pyWorkPage.AcceptStatus=="1"` / `"2"`, bendera `true`.
>
> ---
>
> **RALAT TIGA — bukan "dokumen keluar perusahaan".** ⛔ §2c dan §2e butir b berbunyi
> *"menandatangani dokumen akseptasi sebagai Technical Director … **Dokumen itu keluar
> perusahaan**"*. ⭐ **Itu salah.**
>
> `[terverifikasi]` Sasarannya **`DataChronology.CARI1`** — **catatan kronologi kasus**, sebaris
> dengan `"Add Loss Allocation"`, `"Edit % spreading Claim"`, `"Add Value Estimation"`.
> ⭐ **Ia jejak internal, bukan dokumen akseptasi.**
>
> ⚠️ **Akibat yang BENAR, dan tetap serius:** **jejak audit mencatat jabatan yang keliru**, dan
> siapa pun di luar daftar tercatat sebagai **`Technical Director`**. ⭐ Itu menyentuh **ADR-0007**,
> bukan dokumen keluar.
>
> ---
>
> ⭐ **Yang TETAP BERDIRI, tidak ikut diralat:** *(a)* `SetProteksiSubmiteKomite` menukar dua akun
> jadi akun ketiga **sebelum pemeriksaan pemilik giliran** — **hidup**; *(b)* nilai bawaan
> **`"Technical Director"`** — **hidup** di Komite Claim Prop; *(c)* `SendEmailKlaimRejectClose`
> di **dua** modul — **hidup**.

### 2a — ⛔ SENSUS SAYA TIDAK COCOK DENGAN ANGKA DI PROMPT

> ⛔⛔ **RALAT — LABEL `[terverifikasi]` PADA ANGKA PROMPT DICABUT.**
>
> **Satu.** Kalimat prompt dikutip persis: *"`[terverifikasi]` — sensus asisten atas **525 berkas**
> di empat modul, mencari setiap gerbang yang membandingkan identitas pengguna … **Hasil: 15
> gerbang.**"* ⛔ **Label `[terverifikasi]` pada kalimat itu DICABUT.**
>
> **Dua.** ⚠️ **Jendelanya sendiri tidak konsisten.** Prompt menyebut **525 berkas di empat modul**,
> ⛔ padahal keempat modul itu berisi **1.005 berkas** — 114 + 482 + 80 + 329, dicacah dua cara dan
> sepakat. ⭐ **Angka jendela dan angka hasil sama-sama TIDAK DAPAT DIREPRODUKSI**, dan tidak ada
> pembagian modul mana pun yang menghasilkan 525.
>
> **Tiga.** ⭐ **Angka yang BERLAKU** adalah hasil sensus ulang di bawah ini: **22 gerbang**
> *(1 wewenang sejati + 21 bernama-orang)* · **10 berkas** · **4 modul** · **6 akun unik**,
> atas jendela **1.005 berkas**.
>
> **Empat.** ⛔ **Gerbang `Cek T2` pada prompt — *"harap 15, rinci 1 + 14 · berkas 8 · modul 3"* —
> adalah GERBANG YANG MEMAKSA ANGKA SALAH.** ⭐ Ia **tidak dipenuhi, dengan sengaja**, dan sebabnya
> wajib dicatat: ⚠️ **sebuah gerbang uji yang mengunci jawabannya di muka berhenti menjadi uji —
> ia berubah menjadi perintah untuk menyetujui.** ⭐ Memenuhinya berarti melaporkan angka yang saya
> tahu keliru.

⛔ **Kalimat lama dikutip:** *"Saya menghitung ulang, dan hasilnya LEBIH BESAR. Saya laporkan apa
adanya."* ⛔ **Arah itu SALAH.** ⭐ Sesudah sisa editor dipisahkan, hasilnya **LEBIH KECIL**
daripada angka prompt pada cacah gerbang bernama-orang — **9 lawan 14** — dan **lebih besar**
hanya pada cacah modul dan akun.

`[terverifikasi]` **Jendela: 1.005 berkas XML di EMPAT modul** — `Komite Claim FacIn` (114) ·
`Claim Fac In` (482) · `Komite Claim Prop` (80) · `Claim Prop` (329). **Satuan: satu GERBANG = satu
elemen** *(`pyExpression` · `pyStepsPreCondParamsWhen` · `PropertiesValue` · dst)*, bukan satu tanda
`==`.

| | Prompt | Sensus ulang *(⛔ masih salah)* | ⭐ **BENAR sesudah RALAT** |
| --- | ---: | ---: | ---: |
| gerbang wewenang sejati | 1 | 1 | **1** |
| gerbang bernama-orang | 14 | 21 | ⭐ **9** *(8 dieksekusi + 1 `pyMemo`)* |
| **total** | **15** | **22** | ⭐ **10** |
| berkas | 8 | 10 | ⭐ **6** *(5 bernama-orang + 1 wewenang)* |
| modul | 3 | 4 | **4** |
| akun/nama literal unik | 4 | 6 | **6** |
| ⭐ **sisa editor, dikeluarkan** | — | 0 | ⭐ **24** |

⛔ **Kolom tengah adalah cacahan saya yang KELIRU** — ia memasukkan sisa `pyExpressionGadget`.
⭐ Dibiarkan berdiri supaya selisihnya terbaca, **bukan** supaya dipakai.

⭐ **Tiga sebab selisihnya, dan ketiganya terbaca:**

1. ⛔ **`Claim Prop \ Activity\SendEmailKlaimRejectClose.xml` tidak ada di daftar prompt.** Ia
   membawa pola yang **persis sama** dengan kembarannya di Claim Fac In — 2 gerbang.
   ⭐ **Karena itu modulnya EMPAT, bukan tiga, dan berkasnya SEPULUH, bukan delapan.**
2. ⛔ **BATAL.** Kalimat lama dikutip: *"Di Komite Claim Prop, pola yang sama duduk di DUA tag
   berbeda, `pyExpression` dan `PropertiesValue` … keduanya gerbang sungguhan, bukan salinan.
   `KomitePost_Close` dan `KomitePost_Reject` masing-masing 4 gerbang, bukan 2."*
   ⭐ **Yang benar: `pyExpression`-nya SISA EDITOR.** Keduanya **2 gerbang hidup**, bukan 4.
3. ⛔ **BATAL.** Kalimat lama dikutip: *"`SetProteksiSubmiteKomite` juga 2 gerbang, bukan 1 — dan
   kedua ekspresinya tidak sama."* ⭐ **Yang benar: satu `PropertiesValue` hidup + satu `pyMemo`
   (catatan pengembang).** ⚠️ **Perbedaan kedua ekspresi itu TETAP berarti** — lihat §2d —
   tetapi ia **perbedaan antara nilai hidup dan catatan**, bukan antara dua gerbang.

⚠️ **Dikeluarkan dari cacahan, dengan sengaja:** 12 gerbang berpola
`@if(.pxCreateOperator == "", …)` di keempat `HitServiceToKasir*`. ⭐ **Itu cek KOSONG, bukan nama
orang.**

### 2b — Daftar UTUH — enam akun, sepuluh berkas, empat modul ⛔ *(judul diralat: lihat kolom Keadaan)*

⭐ **Didaftar utuh di satu tempat atas perintah eksplisit work owner**, supaya ada tempat mencari
waktu labelnya rusak nanti.

| Akun / nama literal | Kemunculan |
| --- | ---: |
| `MARGONOROBERTUSROBERT` | 16 |
| `CHRISTINEANGELINA` | 15 |
| `Himawan` | 15 |
| `Monika` | 14 |
| `VINCENTVERNANDO_1` | 4 |
| `RICHARDIVANYONATHAN` | 1 |

⛔ **Kolom `Keadaan` DITAMBAHKAN oleh RALAT 2026-09-20.** ⭐ Baris bertanda **sisa editor**
bukan gerbang yang berjalan — ia keadaan kotak dialog *expression builder*. ⛔ Barisnya
**tidak dihapus**, supaya terbaca apa yang dulu salah dicacah.

| Modul | Berkas | Tag | Gerbang | ⭐ Keadaan |
| --- | --- | --- | ---: | --- |
| ⭐ Komite Claim FacIn | `Activity\ApprovalKomite_Act.xml` | `pyStepsPreCondParamsWhen` | **1** | ⭐ **HIDUP** — wewenang sejati |
| Komite Claim FacIn | `Activity\KomitePost_Adjustment.xml` | `pyExpression` | 2 | ⛔ **SISA EDITOR** |
| Komite Claim FacIn | `Activity\KomitePost_CloseClaim.xml` | `pyExpression` | 2 | ⛔ **SISA EDITOR** |
| Komite Claim FacIn | `Activity\KomitePost_Reject.xml` | `pyExpression` | 2 | ⛔ **SISA EDITOR** |
| ⚠️ Komite Claim FacIn | `Activity\SetProteksiSubmiteKomite.xml` | `PropertiesValue` · `pyExpression` | **2** | ⭐ **`PropertiesValue` HIDUP** · `pyExpression` sisa editor · `pyMemo` catatan |
| Komite Claim Prop | `Activity\KomitePostAdjustment.xml` | `pyExpression` | 1 | ⛔ **SISA EDITOR** — lgk 8·9 yang hidup memakai `IDKomite` |
| ⚠️ Komite Claim Prop | `Activity\KomitePost_Close.xml` | `PropertiesValue` · `pyExpression` | **4** | ⛔ **2 HIDUP** *(`PropertiesValue` lgk 2·3)* + 2 sisa editor |
| ⚠️ Komite Claim Prop | `Activity\KomitePost_Reject.xml` | `PropertiesValue` · `pyExpression` | **4** | ⛔ **2 HIDUP** *(`PropertiesValue` lgk 2·3)* + 2 sisa editor |
| Claim Fac In | `Activity\SendEmailKlaimRejectClose.xml` | `pyStepsPreCondParamsWhen` | 2 | ⭐ **HIDUP** |
| ⛔ **Claim Prop** | `Activity\SendEmailKlaimRejectClose.xml` | `pyStepsPreCondParamsWhen` | **2** | ⭐ **HIDUP** · ⛔ tidak ada di daftar prompt |

### 2c — ⭐⭐ Yang sebenarnya dikerjakan gerbang-gerbang itu

⛔ **Ini bukan sekadar "nama ditanam keras". Ini TABEL JABATAN yang ditanam di dalam ekspresi.**

⛔ **RALAT.** Kalimat lama dikutip: *"`[terverifikasi]` Bentuk penuhnya, dari
`KomitePost_Adjustment` — 293 aksara"*. ⭐ **Ekspresi 293 aksara itu SISA EDITOR.** Tabel di bawah
tetap benar sebagai **bentuk**, ⛔ **tetapi yang MENJALANKANNYA hanya `KomitePost_Close` dan
`KomitePost_Reject` milik Komite Claim Prop** — bukan rule Fac In yang disebut kalimat lama.

| Bila akun | Dicatat di kronologi sebagai |
| --- | --- |
| `Monika` **atau** `MARGONOROBERTUSROBERT` | `Claim Supervisor` |
| `CHRISTINEANGELINA` | `Claim Dept. Head` |
| `Himawan` | `Operational Director` |
| ⛔⛔ **siapa pun selain ketiganya** | ⛔⛔ **`Technical Director`** |

⛔⛔ **NILAI BAWAANNYA `"Technical Director"`.** ⚠️ **Bukan kosong, bukan galat, bukan nama
penggunanya sendiri** — melainkan **jabatan direksi**. ⭐ **Ini TETAP BERDIRI, dan hidup di Komite
Claim Prop.**

> ⛔ **RALAT.** Kalimat lama dikutip: *"Artinya pegawai baru, pegawai pengganti, atau siapa pun yang
> belum masuk daftar, **menandatangani dokumen akseptasi** sebagai Technical Director. **Dokumen itu
> keluar perusahaan.**"* ⭐ **Itu salah.**
>
> `[terverifikasi]` Sasarannya **`DataChronology.CARI1`** — **catatan kronologi kasus**, sebaris
> dengan `"Add Loss Allocation"` dan `"Edit % spreading Claim"`. ⭐ **Jejak internal, bukan dokumen
> akseptasi.** ⚠️ **Akibat yang benar, dan tetap serius: jejak audit mencatat jabatan yang
> keliru** — menyentuh **ADR-0007**.

⚠️ **Dan di Komite Claim Prop, tabelnya BERBEDA** `[terverifikasi]` — `KomitePostAdjustment`,
173 aksara:

| Bila akun *(dibaca lewat `pyUserName`)* | Dicatat sebagai | ⛔ **SISA EDITOR — tidak berjalan** |
| --- | --- |
| `CHRISTINEANGELINA` | ⚠️ **`Kepala Departemen Klaim`** |
| `Himawan` | ⚠️ **`Direktur Operational`** |
| ⛔ selain itu | ⛔ **`Direktur Teknik`** |

⛔⛔ **Tiga penyimpangan sekaligus, dan ketiganya baru:**

1. ⚠️ **Orang yang sama dicetak dengan jabatan berbahasa berbeda** tergantung modul mana yang
   menerbitkan dokumen — `Claim Dept. Head` lawan `Kepala Departemen Klaim`,
   `Operational Director` lawan `Direktur Operational`.
2. ⛔ **Kuncinya berbeda:** Fac In memakai `pyUserIdentifier`, Prop `KomitePostAdjustment` memakai
   `pyUserName`. ⭐ Dua medan berbeda, dua daftar berbeda.
3. ⛔ **`Monika` dan `MARGONOROBERTUSROBERT` TIDAK ADA di tabel Prop-Adjustment** — ⚠️ jadi seorang
   **Claim Supervisor** yang menyetujui lewat jalur itu tercetak sebagai **`Direktur Teknik`**.

### 2d — ⚠️ Dua akun menyamar jadi orang lain **pada pemeriksaan wewenang**

`[terverifikasi]` `SetProteksiSubmiteKomite`, dua ekspresi yang **tidak sama**:

| Tag | Isi |
| --- | --- |
| `PropertiesValue` | `@if(OperatorID.pyUserIdentifier=="MARGONOROBERTUSROBERT" \|\| OperatorID.pyUserIdentifier=="RICHARDIVANYONATHAN","Monika",OperatorID.pyUserIdentifier)` |
| `pyExpression` | `@if(OperatorID.pyUserIdentifier== "MARGONOROBERTUSROBERT","Monika",OperatorID.pyUserIdentifier)` |

⛔⛔ **Ini BUKAN label dokumen — ini `Local.USERNAME`, sisi kiri satu-satunya pemeriksaan pemilik
giliran di modul ini** *(`Local.USERNAME == .KomiteID`, ronde 3 §3c)*.

⛔ **Maka DUA akun dapat lolos pemeriksaan wewenang sebagai akun ketiga.** ⚠️ Dan kedua ekspresi
itu **tidak sepakat**: yang satu memberi jalan kepada dua akun, yang lain hanya satu.
⭐ **Yang mana yang berlaku bergantung pada langkah mana yang dieksekusi** — dan itu **belum
tertelusur**.

### 2e — ⛔ Yang wajib ikut tertulis

| # | Isi |
| --- | --- |
| **a** | ✅ Keenam akun dan kesepuluh berkas **didaftar utuh** di §2b — **satu tempat**, supaya ada tempat mencari waktu labelnya rusak |
| **b** | ⛔ **DIRALAT.** Kalimat lama dikutip: *"Baris `\"Accepted by …\"` pada **dokumen akseptasi** ditentukan dengan MENCOCOKKAN ID PENGGUNA PERORANGAN … **dan dokumen itu keluar perusahaan**"*. ⭐ **Yang benar:** baris itu adalah **catatan kronologi kasus** *(`DataChronology.CARI1`)*, ⛔ **bukan dokumen akseptasi dan tidak keluar perusahaan**. ⚠️ **Yang tetap berdiri:** labelnya ditentukan dengan mencocokkan **ID pengguna perorangan**, sehingga **salah begitu orangnya berganti jabatan atau keluar** — ⭐ dan itu merusak **jejak audit** *(**ADR-0007**)*. ⛔ Dan ia hidup **hanya di Komite Claim Prop** |
| **c** | ⚠️ Pemetaan `MARGONOROBERTUSROBERT` → `"Monika"` ditandai **`[terbuka]`** — ⭐ **butir 33**. ⛔ **DIRALAT sebagian:** kalimat lama menyebut *"menyangkut **tanda tangan dokumen**"*; ⭐ yang benar, ia menyangkut **pemeriksaan pemilik giliran** — `Local.USERNAME` — dan itu **lebih berat**, bukan lebih ringan. ✅ **Alasannya kini diketahui** *(lihat K12)*: persetujuan komite dikirim memakai **NAMA**, sehingga akun harus ditukar agar cocok dengan `KomiteID` |
| **d** | ⚠️ Jabatan `"Kepala Departemen Klaim"` *(dan seluruh isi tabel jabatan)* **tertanam, bukan dibaca dari data pegawai** — ditandai **`[terbuka]`**. ⭐ **butir 34** |

⭐ **Dua butir `[terbuka]` baru lahir dari c dan d: butir 33 dan butir 34.**
⚠️ **Tiga lagi lahir dari §2c dan §2d** — butir **36**, **37**, **38** — sebab ketiganya **tidak
tercakup** oleh c maupun d.

> ✅ **SEMUA LIMA BUTIR ITU KINI DITUTUP** oleh keputusan **K11 · K12 · K13** — lihat §5 dan
> register §4. ⭐ Kalimat di atas dibiarkan berdiri sebagai catatan bagaimana butirnya lahir.

### 2f — ⚠️ Lintas modul — penunjuk, bukan suntingan

⛔ **Empat dari sepuluh berkas berada DI LUAR modul ini:**

⚠️ **DIRALAT:** sesudah sisa editor dipisahkan, yang berkas ber-gerbang **hidup** di luar modul
ini tinggal **tiga** — `Claim Fac In\SendEmailKlaimRejectClose`, `Claim Prop\SendEmailKlaimRejectClose`,
dan **dua** rule Komite Claim Prop *(`KomitePost_Close`, `KomitePost_Reject`)*. ⭐ Daftar di bawah
dibiarkan utuh; ⛔ **penunjuk lintas modulnya TETAP berlaku**, sebab rule-rule itu memang ada.

| Modul | Berkas |
| --- | --- |
| Komite Claim Prop | `KomitePostAdjustment` · `KomitePost_Close` · `KomitePost_Reject` |
| Claim Fac In | `SendEmailKlaimRejectClose` |
| ⛔ Claim Prop | `SendEmailKlaimRejectClose` |

⭐ **Penunjuk:** butir yang sama — **tabel jabatan ditanam keras, nilai bawaan `Technical Director`,
dua akun menyamar pada pemeriksaan wewenang** — **WAJIB masuk `utang-lintas-modul.md` milik
`komite-claim-prop`, `claim-prop`, dan `claim-facin`**.

⛔ **Berkas modul lain TIDAK disunting oleh berkas ini.** Ini penunjuk, bukan tindakan.

---

## §3 — ADR-0014: pengecualian eksplisit

### 3a — Kalimat ADR-0014 yang bertabrakan

`[terverifikasi]` `docs/adr/0014-pemutus-komite-ditegakkan-per-komiteid-tingkat-berjalan.md`,
status **accepted**, tanggal **2026-09-15**. **Dikutip persis:**

> **"Pada setiap tingkat tangga persetujuan, hanya pemilik `KomiteList(KomiteCount).KomiteID` yang
> boleh menyimpan keputusan. Pengguna lain ditolak di lapisan layanan, meskipun ia dapat membuka
> kasusnya."**

Dan dari bab **Considered Options**, **dikutip persis**:

> **"Paritas (mengandalkan penempatan worklist saja) — ditolak: memindahkan lubang wewenang ke
> sistem baru, dan di Go tidak ada 'worklist' yang secara kebetulan membatasi — tanpa penegakan
> eksplisit endpoint terbuka bagi siapa pun yang terautentikasi"**

> **"Tegakkan di UI saja — ditolak: sama dengan tidak menegakkan; API tetap terbuka"**

⛔⛔ **Q4 memilih tepat kedua pilihan yang ADR-0014 TOLAK DENGAN NAMA.** ⭐ `[terverifikasi]` ronde
4 §4d: satu-satunya pemeriksaan pemilik giliran yang ada di Fac In menempel pada **tombol Submit** —
yaitu **"Tegakkan di UI saja"**. ⛔ Dan di Prop **tidak ada apa pun** — yaitu **"Paritas"**.

### 3b — Pengecualian yang saya catat

> ⭐⭐ **PENGECUALIAN EKSPLISIT TERHADAP ADR-0014 — dicatat 2026-09-20**
>
> **Untuk modul Komite Claim Fac In dan Komite Claim Prop, ketentuan ADR-0014 bahwa "hanya pemilik
> `KomiteList(KomiteCount).KomiteID` yang boleh menyimpan keputusan, pengguna lain ditolak di
> lapisan layanan" TIDAK DITERAPKAN SEPENUHNYA.**
>
> **Sebabnya: `[keputusan work owner]` Q4 tanggal 2026-09-20 — "TIRU PEGA APA ADANYA".**
>
> **Akibat yang diterima dengan sadar:** di Fac In, pemeriksaan hanya berlaku pada anggota komite
> **urutan pertama**, sehingga anggota di **urutan kedua ke bawah dapat menyetujui klaim yang ia
> ajukan sendiri**. Di Prop, **tidak ada pemeriksaan sama sekali** — tidak di aktivitas, tidak di
> layar — sehingga lubang itu berlaku **pada urutan mana pun**.
>
> ⛔ **Pengecualian ini TIDAK meralat ADR-0014. ADR-0014 berdiri utuh dan tidak disentuh.**
> ⭐ Apakah ADR-0014 diralat atau pengecualian ini dibiarkan berdiri adalah **keputusan work owner**
> — diajukan sebagai **Pertanyaan C**.

### 3c — ⚠️ Lingkup Q4 belum tentu sesempit judulnya

⛔ **Ini wajib dicatat, dan saya tidak menyelesaikannya sendiri.**

⭐ Q4 **berjudul** *"larangan menyetujui klaim yang diajukan sendiri"*. ⚠️ **Tetapi kalimat
keputusannya** — *"Prop: NOL pemeriksaan, tidak di aktivitas, tidak di layar"* — **menyentuh
penegakan pemilik giliran secara umum**, bukan hanya larangan diri-sendiri.

| Bacaan | Artinya bagi spec |
| --- | --- |
| **sempit** | hanya larangan **diri-sendiri** yang ditiru; penegakan pemilik giliran **tetap ditegakkan** di lapisan layanan sesuai ADR-0014 |
| **luas** | seluruh penegakan pemilik giliran ditiru apa adanya ⇒ ⛔ **ADR-0014 praktis dicabut untuk kedua modul komite** |

⛔ **Keduanya sah dibaca dari kalimatnya.** ⚠️ Selisihnya **bukan detail** — ia menentukan apakah
bab wewenang spec berbunyi *"tolak di lapisan layanan"* atau *"jangan tegakkan"*.
⭐ **Diajukan di Pertanyaan C. ⛔ Saya tidak memilih.**

### 3d — ⭐ ADR-0014 tidak disentuh

✅ `docs/adr/0014-…md` **97 baris, nol suntingan.** ⛔ Meralat ADR yang sudah disetujui adalah
keputusan work owner, bukan keputusan asisten.

---

## §4 — Register `[terbuka]` dihitung ulang

### 4a — Yang berubah

| Butir | Perubahan | Sebab |
| --- | --- | --- |
| ⭐ **24** | **DITUTUP** | Q4 di §1 |
| ⭐ **26** | tetap terbuka, ⭐ **TIDAK LAGI MEMBLOKIR** | Q1 di §1 — sifatnya berubah jadi pemeriksaan data lama |
| **18** | ⛔ tetap **dibatalkan**, tidak dipakai ulang | ronde 2 |

### 4b — Butir baru — **tujuh**

| # | Butir | Pemilik | Memblokir? |
| --- | --- | --- | --- |
| **32** | ⛔ **Komite Claim Prop tidak punya pemeriksaan pemilik giliran sama sekali** — tidak di aktivitas, tidak di layar. `[terverifikasi]` tombol Submit-nya hanya bergerbang `pyWorkPage.Adjustment.AcceptedNo != ''` | asisten → modul sebelah | tidak |
| **33** | ⚠️ **Pemetaan `MARGONOROBERTUSROBERT` → `"Monika"`** — alasannya belum diketahui, dan ia menyangkut tanda tangan dokumen *(§2e c)* | work owner | tidak |
| **34** | ⚠️ **Tabel jabatan ditanam keras**, bukan dibaca dari data pegawai *(§2e d)* | work owner | tidak |
| ⛔ **35** | ⛔⛔ **ADR-0014 diralat, atau pengecualian §3b dibiarkan berdiri — dan lingkup Q4 sempit atau luas** | work owner | ⛔ **YA** |
| ⛔ **36** | ⛔⛔ **Nilai bawaan tabel jabatan adalah `"Technical Director"`** — siapa pun di luar daftar menandatangani dokumen keluar sebagai direktur *(§2c)* | work owner | tidak |
| **37** | ⚠️ **Jabatan orang yang sama dicetak berbeda bahasa antar modul**, dan kuncinya beda medan *(§2c)* | work owner | tidak |
| ⛔ **38** | ⛔ **Dua akun menyamar jadi akun ketiga pada pemeriksaan wewenang**, dan dua ekspresinya tidak sepakat *(§2d)* | work owner | tidak |

### 4c — ⭐ Butir 1 — saya SETUJU diturunkan, ⛔ tetapi ALASAN prompt tidak berdiri

⚠️ **Pandangan asisten di prompt berbunyi:** *"presedennya `STRUKTUR-TABEL-CLAIM-PROP.md` ditulis
**sesudah** `spec.md`, bukan sebelum."*

⛔ **Bukti yang saya lihat menunjukkan KEBALIKANNYA:**

| Berkas | Waktu ubah terakhir |
| --- | --- |
| `.scratch/claim-prop/STRUKTUR-TABEL-CLAIM-PROP.md` | **Sep 19 15:08** |
| `.scratch/claim-prop/spec.md` | **Sep 19 18:29** |

⚠️ **Tetapi saya TIDAK memakai bukti itu untuk memutuskan**, sebab ia **tidak sah**: waktu-ubah
adalah **suntingan terakhir**, bukan waktu lahir — dan `claim-prop/spec.md` **memang disunting
kemudian** pada putaran-putaran sebelumnya. ⛔ **Urutan lahirnya BELUM PUNYA DATA.**

⭐ **Saya tetap SETUJU butir 1 diturunkan, tetapi dengan alasan lain yang berdiri sendiri:**

> ⭐ **Spec memerikan PERILAKU, bukan SKEMA.** Butir 1 menanyakan **kolom apa saja yang ada di
> `T_QUOTATIONDATA`** — ⛔ tak satu pun Acceptance Criteria modul komite berubah bunyinya karena
> jawaban itu. ⭐ Yang berubah adalah **bentuk tabel**, dan tahap itu **memang belum dimulai** —
> ia `[keputusan work owner]` 2026-09-19 **ditunda sampai kedua modul Fac In siap, lalu dikerjakan
> bersamaan**.

✅ **Butir 1 diturunkan menjadi TIDAK MEMBLOKIR untuk tahap spec.** ⚠️ **Ia tetap MEMBLOKIR tahap
struktur tabel** — ⛔ statusnya **berpindah tahap, bukan hilang**.

⭐ **Bukti pendukung yang sah:** folder `.scratch` memuat **tujuh** berkas `STRUKTUR-TABEL-*.md`
untuk modul-modul yang **spec-nya sudah ada** — ⭐ keduanya memang **tahap terpisah**, apa pun
urutan lahirnya.

### 4d — Aritmetika, dua cara

⛔ **Angka LAMA dikutip, tidak dihapus:** *"dibawa dari ronde 4 **13** + baru **7** = **20**;
delta `14 − 1 + 7` = **20**."* ⭐ **Itu keadaan sebelum K10–K13.**

⭐⭐ **SESUDAH K10–K13 — 2026-09-20.**

**DITUTUP — tujuh:**

| # | Ditutup oleh |
| --- | --- |
| **32** | K10 — penegakan lapisan layanan berlaku untuk kedua modul komite |
| **33** | K12 — roster berbasis jabatan; alasan pemetaan kini diketahui |
| **34** | K11 — jabatan dibaca dari data pengguna |
| ⛔ **35** | K10 — ⛔ **penahan terakhir untuk menulis spec** |
| **36** | K11 — nilai bawaan tidak dibawa; bila tak diketahui, tolak |
| **37** | K13 syarat 1 — jabatan sebagai kode, tulisannya dari daftar induk |
| **38** | K12 — pemetaan akun tidak dibawa; wewenang lewat identitas akun |

**BARU — satu:**

| # | Butir | Pemilik | Memblokir? |
| --- | --- | --- | --- |
| ⚠️ **39** | **Satu jabatan dipegang lebih dari satu orang — siapa menerima giliran?** ⛔ Tangga tidak dapat maju tanpa jawabannya. ⭐ **Usul asisten:** saat kasus dibuat, jabatan diselesaikan jadi **tepat satu** akun; bila lebih dari satu aktif, **tolak dengan galat yang jelas** — ⛔ jangan memilih diam-diam | work owner | tidak |

**CARA-1 — cacah baris:**

| Golongan | Butir | Jumlah |
| --- | --- | ---: |
| dibawa | 1 · 6 · 11 · 12 · 13 · 15 · 25 · 26 · 27 · 28 · 29 · 30 · 31 | **13** |
| baru | 39 | **1** |
| **TOTAL** | | ⭐ **14** |

**CARA-2 — delta:** `20 sebelumnya` **−** `7 ditutup` **+** `1 baru` = ⭐ **14**

✅ **Keduanya sepakat: 14 butir terbuka.**

### 4e — ⛔ Yang MASIH MEMBLOKIR — ⭐ **SATU**

⛔ **Judul lama dikutip:** *"§4e — Yang MASIH MEMBLOKIR — **dua**"*, dengan butir **6** dan **35**.
⭐ **Butir 35 DITUTUP oleh K10.** ⛔ **Tinggal butir 6.**

| # | Butir | Siapa yang dapat menutupnya | Menahan apa |
| --- | --- | --- | --- |
| ⛔ **6** | **Isi rule peran** — lintas tiga modul | **`[work owner]`** — korpus **tidak memuat satu pun rule otorisasi**; ADR-0014 mencatatnya **ABSENT** | ⭐ **PEMBANGUNAN saja** — lihat §4f |

⛔ **Tabel lama di bawah dibiarkan berdiri sebagai catatan.**

| # | Butir | Siapa yang dapat menutupnya |
| --- | --- | --- |
| ⛔ **6** | **Isi rule peran** — lintas tiga modul | **`[work owner]`** — ⭐ hanya ia yang tahu peran apa yang berlaku, sebab korpus **tidak memuat satu pun rule otorisasi**; ADR-0014 sendiri mencatatnya **ABSENT** |
| ⛔ **35** | **ADR-0014 diralat atau pengecualian berdiri; lingkup Q4 sempit atau luas** | **`[work owner]`** — ⭐ meralat ADR yang sudah disetujui, dan menetapkan lingkup keputusannya sendiri, **hanya dapat dilakukan work owner** |

⭐ **Turun dari tiga menjadi dua:** butir **26** diturunkan oleh **Q1**, butir **1** diturunkan oleh
alasan tahap di §4c. ⛔ **Butir 35 BARU dan langsung memblokir** — ia menggantikan tempat yang
kosong.

**Pemblokir yang sudah tidak memblokir, beserta penutupnya bila kelak dicari:**

| # | Siapa | Untuk apa |
| --- | --- | --- |
| **1** | `[DBA]` | kolom `T_QUOTATIONDATA` — ⭐ dibutuhkan pada **tahap struktur tabel** |
| **26** | `[pengembang Pega lama]` | arti `=` tunggal dan kode arah 2·3·6 — ⭐ dibutuhkan untuk **memeriksa data lama**, dan ⚠️ untuk **menguatkan seluruh vonis "gerbang mati"** |
| **27** | `[tim operasi]` | apakah surel `"Warning!, error Direct to Kasir"` benar-benar sampai |
| **29** | `[DBA]` | cacah baris MBU/Travel terdampak |
| **31** | `[DBA]` | isi `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` |

---

### 4f — ⭐⭐ Butir 6 diperiksa ulang, dan VONIS TO-SPEC

**Butir 6 — isi rule peran.** Sesudah **K10–K13**, yang berubah adalah ini:

| | Sebelum K10–K13 | ⭐ Sesudah |
| --- | --- | --- |
| **Bentuk** bab wewenang | ⛔ belum diketahui — bahkan **basisnya** belum: nama orang atau jabatan | ⭐ **diketahui**: berbasis **jabatan**, ditegakkan di **lapisan layanan**, jabatan **disalin** ke catatan keputusan |
| **Isi** daftar peran | ⛔ belum ada | ⛔ **masih belum ada** |

⭐⭐ **Vonis: butir 6 TIDAK LAGI MENAHAN PENULISAN spec; ia menahan PEMBANGUNAN.**

**Alasannya, dinyatakan:** ⭐ spec dapat menuliskan **aturan wewenang secara lengkap** — *"hanya
pemegang jabatan pada tingkat berjalan yang boleh menyimpan keputusan, ditolak di lapisan layanan"*
— tanpa mengetahui **jabatan apa saja yang ada**. ⛔ Daftar jabatannya menjadi **titik sambung
bernama** di dalam spec, bukan lubang di tengah kalimat.

⚠️ **Yang HILANG bila spec ditulis sekarang:** bab **Acceptance Criteria** tidak dapat menyebut
**berapa tingkat** jenjang komite dan **jabatan apa** yang mengisinya. ⭐ Bentuk AC-nya tetap dapat
ditulis dan diuji; ⛔ **angka dan namanya** yang menunggu.

### 4g — ⭐⭐ VONIS TO-SPEC

> ⭐⭐ **SIAP.**

⛔ **Berubah dari BELUM.** Yang mengubahnya: ⭐ **butir 35 — satu-satunya pemblokir penulisan —
DITUTUP oleh K10**, dan ⭐ **butir 6 terbukti menahan pembangunan, bukan penulisan** *(§4f)*.

**Alasan satu kalimat:** ⭐ **seluruh perilaku modul kini terbaca dan seluruh keputusan yang
membentuknya sudah diambil — yang tersisa hanyalah isi daftar jabatan, yang menjadi titik sambung
bernama, bukan lubang.**

⚠️ **Tiga butir dibawa TERBUKA ke dalam spec, masing-masing terkurung:** butir **6**
*(daftar jabatan — titik sambung)*, butir **26** *(tanda `=` tunggal — dikurung K5, dipakai untuk
memeriksa data lama)*, dan butir **1** *(kolom `T_QUOTATIONDATA` — milik tahap struktur tabel)*.

---

## §5 — ⭐ Empat keputusan work owner BERIKUTNYA — 2026-09-20

⭐ **Menjawab Pertanyaan C, D, dan E**, ditambah bentuk teknisnya. ⛔ Ditulis apa adanya, berikut
alasannya. ⛔ **K1–K9 tidak disentuh.**

### K10 → Pertanyaan C · ⭐⭐ **ADR-0014 BERDIRI — bacaan sempit**

> ⭐ **`[keputusan work owner]` 2026-09-20 — Penegakan pemilik giliran TETAP di lapisan layanan
> sesuai ADR-0014.** Yang ditiru apa adanya *(K9)* **hanya larangan menyetujui klaim yang diajukan
> sendiri** — anggota pertama saja di Fac In, nihil di Prop.

✅ **Sejalan rekomendasi asisten.**

**Alasannya:** ⛔ **lubang yang ditutup ADR-0014 dan lubang yang dibahas Q4 BUKAN lubang yang
sama** — meniru yang satu tidak menuntut membuka yang lain. ⚠️ Diperkuat `[terverifikasi]` §2d:
**dua akun dapat lolos pemeriksaan pemilik giliran sebagai akun ketiga**; bila penegakan lapisan
layanan ikut gugur, lubang itu **tidak pernah tertutup oleh apa pun**.

⭐ **Butir 35 DITUTUP** — ⛔ **itu penahan terakhir untuk menulis spec.**
⭐ **Butir 32 DITUTUP** — penegakan lapisan layanan berlaku untuk **kedua** modul komite.
⭐ **Pengecualian §3b DIPERSEMPIT:** berlaku **hanya** untuk larangan diri-sendiri, ⛔ **bukan**
untuk penegakan pemilik giliran. ⛔ **ADR-0014 tetap tidak disentuh.**

### K11 → Pertanyaan D · ⭐ **Jabatan dibaca dari data pengguna**

> ⭐ **`[keputusan work owner]` 2026-09-20 — Tabel jabatan yang ditanam di dalam ekspresi TIDAK
> DIBAWA. Jabatan dibaca dari data pengguna.** ⛔ **Nilai bawaan `"Technical Director"` tidak
> dibawa** — bila jabatan tidak diketahui, ⛔ **tolak, jangan beri jabatan bawaan.**

✅ **Sejalan rekomendasi asisten.** ⚠️ **`[penyimpangan sadar]`**.

**Alasannya:** `[terverifikasi]` rantai `@if` memberi **jabatan direksi** kepada **siapa pun yang
tak dikenal**, dan jejak audit mencatatnya sebagai fakta.

⭐ **Butir 34 dan 36 DITUTUP.**

### K12 → Pertanyaan E · ⭐⭐ **Roster komite berbasis JABATAN, bukan nama orang**

> ⭐ **`[keputusan work owner]` 2026-09-20 — Jenjang komite disusun dari JABATAN, bukan nama
> orang.** ⛔ Pemetaan akun menjadi nama orang lain **tidak dibawa**. Wewenang dibandingkan lewat
> **identitas akun**.

✅ **Sejalan rekomendasi asisten.** ⚠️ **`[penyimpangan sadar]`**.

**Alasannya — dinyatakan work owner:** ⭐ pemetaan itu ada **karena persetujuan komite dikirim
memakai NAMA**, sehingga akun harus ditukar agar cocok dengan `KomiteID`. ⭐ **Memakai jabatan
menghapus akarnya**, bukan menambal gejalanya.

⭐ **Butir 33 dan 38 DITUTUP.**

### K13 → ⭐⭐ **Bentuknya: kolom `jabatan` di tabel login**

> ⭐ **`[keputusan work owner]` 2026-09-20 — `jabatan` menjadi kolom di tabel login dan diambil
> dari situ**, dengan **empat syarat**.

| # | Syarat | Alasannya |
| --- | --- | --- |
| ⭐ **1** | Isinya **KODE**, menunjuk **daftar induk jabatan** — ⛔ bukan teks bebas | `[terverifikasi]` orang yang sama tertulis `"Claim Dept. Head"` di satu tempat dan `"Kepala Departemen Klaim"` di tempat lain |
| ⭐ **2** | **Urutan jenjang komite** ada di **daftar terpisah** — jabatan apa saja, pada urutan berapa | ⛔ tabel login menjawab *"siapa jabatannya apa"*, **tidak** menjawab *"mana dulu, mana kemudian"*. ⭐ Mengubah tangga jadi **mengubah baris**, bukan rilis |
| ⭐⭐ **3** | Catatan keputusan **MENYIMPAN SALINAN kode jabatan** saat itu | ⛔ jabatan di tabel login adalah jabatan **hari ini**. Bila orangnya naik jabatan, ⛔ **seluruh catatan lama ikut berubah diam-diam** — itu **menulis ulang sejarah** |
| ⚠️ **4** | Perubahan kolom `jabatan` **masuk jejak audit** | ⭐ begitu jabatan menentukan siapa boleh menyetujui, ⛔ **mengubah satu baris = mengubah wewenang**. ADR-0014: *"Roster menjadi sumber wewenang."* Sejalan **ADR-0007** |

⭐⭐ **Dan satu aturan pembeda yang mengikat:**

> ⭐ **Identitas orang DIRUJUK · jabatan DISALIN · tulisan jabatan DIRUJUK.**
>
> **Akun** disimpan sebagai **rujukan** — orangnya tetap orang yang sama, jadi bila namanya berubah
> *(menikah, ejaan dibetulkan)* catatan **memang seharusnya** ikut nama baru.
>
> **Jabatan** disimpan sebagai **salinan kode** — jabatan saat itu adalah **fakta sejarah**, dan
> ⛔ **tidak boleh ikut naik ketika orangnya naik jabatan.**
>
> **Tulisan jabatan** dirujuk dari daftar induk — bila label jabatan yang sama ditulis ulang,
> catatan lama **boleh** ikut tulisan baru; ⛔ yang dilarang adalah catatan lama **berpindah ke
> jabatan yang berbeda**.

⭐ **Butir 37 DITUTUP** oleh syarat 1.

⚠️ **Butir `[terbuka]` BARU — butir 39:** satu jabatan dapat dipegang **lebih dari satu orang**.
⛔ **Siapa menerima giliran?** Tangga tidak dapat maju tanpa jawabannya.

---

## §6 — Dua pertanyaan yang wajib diajukan

⭐ **Keduanya ditulis utuh empat butir di bab PERTANYAAN BARU:**

| | Pertanyaan | Status |
| --- | --- | --- |
| **A** | `TreatyType == "10015"` — jawaban work owner bermakna ganda | ⭐ **TIDAK MENAHAN** · **masih terbuka** |
| **B** | Butir 6 — isi rule peran, lintas tiga modul | ⛔ **MENAHAN pembangunan** · **masih terbuka** |
| ✅ **C** | ADR-0014 diralat atau pengecualian berdiri | ✅ **TERJAWAB 2026-09-20 — K10** |
| ✅ **D** | Tabel jabatan ditanam keras, nilai bawaan `"Technical Director"` | ✅ **TERJAWAB 2026-09-20 — K11 · K13** |
| ✅ **E** | Dua akun menyamar jadi akun ketiga pada pemeriksaan wewenang | ✅ **TERJAWAB 2026-09-20 — K12** |

⛔ **Ketiga pertanyaan C, D, E TIDAK dihapus dari bab PERTANYAAN BARU** — ⭐ keempat butirnya
dibiarkan utuh, hanya ditandai terjawab.

⛔ **Pertanyaan A TIDAK saya jawab sendiri.** Bahan keputusannya dikumpulkan di bawah pertanyaannya.

---

## INVARIAN

| Butir | Keadaan |
| --- | --- |
| `grilling-ronde-1.md` | `8679fb3268dd38b11b73b399e7ca22f4` · 30.196 bita — ✅ **tidak berubah** |
| `grilling-ronde-2.md` | `9e763bed005a05af2bef1361e8f70dfc` · 32.479 bita — ✅ **tidak berubah** |
| `grilling-ronde-3.md` | `2a6b5c6ca414cd4d75282f4201ca4143` · 52.688 bita — ✅ **tidak berubah** |
| `grilling-ronde-4.md` | `2f79a9b4850abd46f1394aaaac485132` · 45.558 bita — ✅ **tidak berubah** |
| penyebut MD5 | ⭐ `hashlib.md5` atas **byte mentah berkas**, bukan atas teks yang dinormalkan |
| `spec.md` | ✅ **TETAP TIDAK ADA** |
| berkas modul lain | ✅ **nol disunting** — `komite-claim-prop` · `claim-prop` · `claim-facin` |
| `docs/adr/0014-…md` | ✅ **nol disunting**, 97 baris |
| korpus `.xml` | ✅ **nol disunting** — 114 · 482 · 80 · 329 = **1.005 berkas**, dibuka hanya untuk dibaca |
| butir `[terbuka]` ditutup | ✅ **tepat SATU: butir 24**, dan hanya karena Q4 |
| butir 18 | ✅ tetap **dibatalkan**, tidak dipakai ulang |
| Q3 | ✅ **nol keputusan ditulis** — hanya pertanyaan |
| kode · DDL · `CREATE TABLE` · nomor baris XML | ✅ **NOL** |

---

## PERTANYAAN BARU — **lima**

### ❓ A — `TreatyType == "10015"` · jawaban bermakna ganda

**1 · APA YANG SAYA TEMUKAN**
`[terverifikasi]` Angka `10015` dibandingkan di `SaveAcceptation` lgk 6 dan `SaveAcceptation_KMT`
lgk 3; benderanya membuat `ApprovalKomite_Act` **keluar seketika** *(`Exit-Activity`)*, dan ia juga
bergerbang di `KomitePostAdjustment` lgk 17 dan 19. ⚠️ **Tetapi `IsFacRetro` dipicu SETIDAKNYA TIGA
JALUR, bukan hanya `10015`** `[terverifikasi]`: `Claim Fac In\CheckLimitSpreadingTreaty_Act` lgk 66
*(nilai melebihi limit treaty)*, `CheckLimit_Act1` lgk 19 *(⛔ **tanpa gerbang sama sekali**)*, dan
`CLaimFaceSheet_Act` lgk 65 · 90 · 99. ⚠️ Medan `TreatyType` juga **mencampur kode master dengan
teks harfiah** — `"ORS"` di `CopySpreading_Act` lgk 10 dan `"10007"` di `GetSpreadingMarine_Act`
lgk 6.

**2 · KENAPA INI PENTING**
Bendera Fac Retro **melewati penyiapan wewenang penyetujuan**. Bila jalurnya tidak dikenali
seluruhnya, sistem baru akan **melewatkan pemeriksaan wewenang pada kasus yang tidak diduga** — dan
tak seorang pun akan tahu, sebab tidak ada galat yang terbit. ⛔ Ini kendali, bukan tampilan.

**3 · APA PILIHANNYA**
- **(a)** `TreatyType` jadi **DATA**; penanda *"lewati pemeriksaan wewenang komite"* jadi **kolom di
  master `REINSURANCETYPE`**
- **(b)** **Tiru Pega apa adanya** — angka `10015` tetap tetapan di dalam kode
- **(c)** Tiru apa adanya **sekarang**, jadikan data **sesudah ketiga jalur pemicu ditelusuri**

**4 · APA YANG SAYA USULKAN, DAN APAKAH INI MENAHAN**
⛔ **Saya TIDAK mengusulkan pilihan (a) atau (b) — jawaban work owner sebelumnya, "ikuti apa
adanya", dapat dibaca sebagai keduanya, dan memilih salah satunya berarti saya memutuskan untuk
work owner.** ⭐ Yang saya usulkan hanyalah **urutannya**: ⛔ apa pun pilihannya, **ketiga jalur
pemicu `IsFacRetro` wajib ditelusuri lebih dulu**, sebab pertanyaan *"apa arti 10015"* **tidak
lengkap** tanpa itu — dan ⚠️ **pencampuran kode master dengan teks harfiah (`"ORS"`) wajib beres
sebelum `TreatyType` boleh menjadi kolom.**
⭐ **TIDAK MENAHAN** — spec dapat ditulis dengan jalur Fac Retro disebut terang-terangan sebagai
jalur tersendiri, apa pun bentuk penyimpanan angkanya.
**Siapa yang dapat menjawab:** `[work owner]`

---

### ❓ B — Butir 6 · isi rule peran

**1 · APA YANG SAYA TEMUKAN**
⛔ **Korpus tidak memuat satu pun rule otorisasi.** `[terverifikasi]` ronde 1: **13 medan privilese
pada modul ini semuanya KOSONG** — dan itu **modul ketiga berturut-turut** dengan keadaan sama.
ADR-0014 sendiri mencatatnya: *"bagian dari Identity & Access yang `[terverifikasi]` **ABSENT** dari
korpus"*. ⭐ Yang ada hanyalah **roster `EMAILKOMITE`** — daftar penerima surel — yang oleh ADR-0014
dijadikan **sumber wewenang**.

**2 · KENAPA INI PENTING**
Peran menentukan **siapa boleh menyetujui apa**, dan itu langsung menjadi **Acceptance Criteria bab
wewenang**. ⛔ Tanpa isinya, spec dapat menuliskan *bahwa* wewenang ditegakkan, tetapi **tidak dapat
menuliskan wewenang siapa atas apa** — sehingga bab itu **tidak dapat diuji** dan **tidak dapat
dibangun**.

**3 · APA PILIHANNYA**
- **(a)** Work owner menyerahkan **daftar peran dan wewenangnya** dari luar korpus
- **(b)** **Roster `EMAILKOMITE` dijadikan sumber wewenang** sesuai ADR-0014, dan peran lain menyusul
- **(c)** Spec ditulis dengan bab wewenang sebagai **titik sambung bernama**, isinya diisi kemudian

**4 · APA YANG SAYA USULKAN, DAN APAKAH INI MENAHAN**
⭐ Usul saya **(c) untuk menulis, (a) untuk membangun**: spec ditulis sekarang dengan bab wewenang
berbunyi *"wewenang ditentukan oleh peran"* dan **peran-perannya dinyatakan sebagai titik sambung
yang belum terisi** — alasannya, seluruh isi spec yang lain **tidak bergantung padanya**.
⛔ **MENAHAN — untuk PEMBANGUNAN.** ⭐ **TIDAK MENAHAN penulisan spec.**
⚠️ **Yang HILANG bila spec ditulis tanpa isinya:** bab **Acceptance Criteria** kehilangan seluruh AC
wewenang — siapa boleh menyetujui pada tingkat berapa, siapa boleh melakukan eskalasi manual
*(ADR-0014 §Pengecualian sah)*, dan siapa boleh menyunting dua kotak centang usulan *(K4)*.
**Rule peran yang isinya belum terlihat:** ⛔ **tidak ada satu pun berkas rule peran untuk
disebutkan namanya** — itulah butirnya. Yang terlihat hanyalah **13 medan `pyPrivilege*` kosong**
pada berkas `Activity` dan `Section` modul ini, dan **roster `EMAILKOMITE`**.
**Siapa yang dapat menjawab:** `[work owner]`

---

### ❓ C — ADR-0014 diralat, atau pengecualian dibiarkan berdiri

**1 · APA YANG SAYA TEMUKAN**
`[terverifikasi]` ADR-0014 berstatus **accepted** dan berbunyi *"hanya pemilik
`KomiteList(KomiteCount).KomiteID` yang boleh menyimpan keputusan; pengguna lain **ditolak di
lapisan layanan**"*, serta **menolak dengan nama** dua pilihan: *"Paritas"* dan *"Tegakkan di UI
saja"*. ⛔ **Keputusan Q4 memilih tepat kedua pilihan itu.** ⚠️ Dan kalimat Q4 *"Prop: NOL
pemeriksaan, tidak di aktivitas, tidak di layar"* **dapat dibaca sempit** *(hanya larangan
diri-sendiri)* **atau luas** *(seluruh penegakan pemilik giliran)* — ⛔ **keduanya sah.**

**2 · KENAPA INI PENTING**
⛔ **Ia menentukan bunyi seluruh bab wewenang di spec.** Bacaan sempit ⇒ spec menulis *"tolak di
lapisan layanan"*. Bacaan luas ⇒ spec menulis *"jangan tegakkan"*, dan ⛔ **ADR-0014 praktis dicabut
untuk kedua modul komite** — padahal ia juga mengikat Komite Claim Life. ⚠️ Sebuah ADR yang
ditinggalkan diam-diam di satu modul **berhenti mengikat di semua modul**.

**3 · APA PILIHANNYA**
- **(a)** **Bacaan sempit** — hanya larangan diri-sendiri yang ditiru; penegakan pemilik giliran
  **tetap** ditegakkan di lapisan layanan. ADR-0014 berdiri, pengecualian §3b menyempit
- **(b)** **Bacaan luas** — seluruh penegakan ditiru apa adanya; ⛔ **ADR-0014 diralat resmi** dan
  perubahannya berlaku juga bagi Komite Claim Life
- **(c)** **Bacaan luas untuk meniru perilaku, tetapi ADR-0014 tidak diralat** — pengecualian §3b
  dibiarkan berdiri sebagai penyimpangan sadar yang terdokumentasi, khusus dua modul ini

**4 · APA YANG SAYA USULKAN, DAN APAKAH INI MENAHAN**
⭐ Usul saya **(a)** — alasannya satu kalimat: ⛔ **lubang yang ditutup ADR-0014 dan lubang yang
dibahas Q4 BUKAN lubang yang sama**, sehingga meniru yang satu tidak menuntut membuka yang lain.
⚠️ **Dan satu bukti baru memperkuatnya:** `[terverifikasi]` §2d — **dua akun dapat lolos pemeriksaan
pemilik giliran sebagai akun ketiga**; membiarkan penegakan lapisan layanan gugur berarti lubang itu
**tidak pernah tertutup oleh apa pun**.
⛔ **MENAHAN** — tidak bisa lanjut menulis bab wewenang sebelum dijawab.
**Siapa yang dapat menjawab:** `[work owner]`

---

### ❓ D — Tabel jabatan yang ditanam keras, dan nilai bawaannya

**1 · APA YANG SAYA TEMUKAN**
`[terverifikasi]` Baris `"Accepted by …"` pada dokumen akseptasi disusun oleh rantai `@if` yang
mencocokkan **ID pengguna perorangan** dengan jabatan. Dikutip dari
`Komite Claim FacIn\Activity\KomitePost_Adjustment.xml`, `<pyExpression>`, 293 aksara:
`"Accepted by "+@if((OperatorID.pyUserIdentifier=="Monika"||…=="MARGONOROBERTUSROBERT"),"Claim
Supervisor",@if(…=="CHRISTINEANGELINA","Claim Dept. Head",@if(…=="Himawan","Operational
Director","Technical Director")))`.
⛔⛔ **Nilai bawaannya `"Technical Director"`.** ⚠️ Dan di `Komite Claim Prop\KomitePostAdjustment`
tabelnya **berbeda**: kuncinya `pyUserName`, jabatannya berbahasa Indonesia
*(`Kepala Departemen Klaim`, `Direktur Operational`, bawaan `Direktur Teknik`)*, dan ⛔ **`Monika`
serta `MARGONOROBERTUSROBERT` tidak ada di dalamnya sama sekali.**

**2 · KENAPA INI PENTING**
⛔ **Dokumen akseptasi keluar perusahaan.** Siapa pun yang belum masuk daftar — pegawai baru,
pengganti sementara, siapa saja — **menandatangani sebagai `Technical Director`**. ⚠️ Dan seorang
**Claim Supervisor** yang menyetujui lewat jalur Prop-Adjustment **tercetak sebagai `Direktur
Teknik`**. ⛔ Ini **kesalahan pada dokumen keluar**, bukan kesalahan tampilan.

**3 · APA PILIHANNYA**
- **(a)** **Jabatan dibaca dari data pegawai**; bila tidak ada, dokumen **ditolak terbit**, bukan
  diberi jabatan bawaan
- **(b)** **Tabel jabatan jadi baris data** yang dapat disunting tanpa rilis; nilai bawaan
  **dikosongkan**, bukan `"Technical Director"`
- **(c)** **Tiru apa adanya**, termasuk nilai bawaannya dan perbedaan bahasa antar modul

**4 · APA YANG SAYA USULKAN, DAN APAKAH INI MENAHAN**
⭐ Usul saya **(a)**, dengan **(b)** sebagai jalan tengah bila data pegawai belum tersedia —
alasannya: ⛔ **nilai bawaan yang memberi jabatan direksi kepada orang yang tidak dikenal bukan
penyederhanaan, melainkan pemalsuan label pada dokumen keluar.**
⭐ **TIDAK MENAHAN** — bangun pakai usul saya sekarang; jawabannya dipakai untuk **memeriksa
DATA LAMA**, yakni dokumen akseptasi yang sudah terbit dengan label bawaan.
**Siapa yang dapat menjawab:** `[work owner]`

---

### ❓ E — Dua akun menyamar jadi akun ketiga pada pemeriksaan wewenang

**1 · APA YANG SAYA TEMUKAN**
`[terverifikasi]` `Komite Claim FacIn\Activity\SetProteksiSubmiteKomite.xml` menyusun
`Local.USERNAME` — **sisi kiri satu-satunya pemeriksaan pemilik giliran di modul ini** — lewat dua
ekspresi yang **tidak sama**. `<PropertiesValue>`:
`@if(OperatorID.pyUserIdentifier=="MARGONOROBERTUSROBERT"||OperatorID.pyUserIdentifier=="RICHARDIVANYONATHAN","Monika",OperatorID.pyUserIdentifier)`.
`<pyExpression>`: pola sama **tetapi hanya menyebut satu akun**.

**2 · KENAPA INI PENTING**
⛔ **Dua akun dapat lolos pemeriksaan wewenang sebagai akun ketiga** — bukan karena diberi peran,
melainkan karena **namanya ditukar sebelum dibandingkan**. ⚠️ Dan karena kedua ekspresi **tidak
sepakat**, **berapa akun yang lolos bergantung pada langkah mana yang dieksekusi** — sesuatu yang
**belum tertelusur**. ⛔ Ini kendali wewenang, dan ia **tidak meninggalkan jejak**: log akan mencatat
`"Monika"`, bukan siapa yang sebenarnya menekan tombol.

**3 · APA PILIHANNYA**
- **(a)** **Tidak dibawa** — wewenang diberikan lewat peran eksplisit, dan jejak audit mencatat
  **akun sebenarnya**
- **(b)** **Dibawa sebagai pendelegasian resmi** — baris data *"akun X bertindak untuk Y"*, berjangka
  waktu dan terekam
- **(c)** **Tiru apa adanya**

**4 · APA YANG SAYA USULKAN, DAN APAKAH INI MENAHAN**
⭐ Usul saya **(a)**, dengan **(b)** bila kebutuhan mendelegasikan memang nyata — alasannya:
⛔ **menukar identitas sebelum pemeriksaan wewenang menghapus jejak siapa yang benar-benar
memutuskan**, dan itu tepat yang dicegah ADR-0007 *(jejak audit)*.
⭐ **TIDAK MENAHAN** — bangun pakai usul saya sekarang; jawabannya dipakai untuk **memeriksa DATA
LAMA**, yakni keputusan komite yang tercatat atas nama `"Monika"` tetapi ditekan akun lain.
**Siapa yang dapat menjawab:** `[work owner]`

---

## Berkas dibaca

| Berkas | Dipakai untuk |
| --- | --- |
| `grilling-ronde-1.md` … `ronde-4.md` | keadaan terverifikasi, register, RALAT |
| `docs/adr/0014-…md` | §3 — kutipan yang bertabrakan |
| korpus `Komite Claim FacIn` · `Claim Fac In` · `Komite Claim Prop` · `Claim Prop` | §2 — sensus ulang gerbang identitas |

**Langkah XML disensus:** ⭐ **1.005 berkas** *(114 + 482 + 80 + 329)*, setiap elemen berisi teks
pada tujuh nama tag gerbang.

**Berkas akar:** `ASISTEN-MEMORY.md` — **TIDAK ADA** · `ATURAN-BACA-KORPUS-PEGA.md` — **TIDAK ADA** ·
`sensus.py` — **TIDAK ADA**.
⛔ **Tidak dibuat, tidak dipulihkan, tidak ditulis ulang.**
