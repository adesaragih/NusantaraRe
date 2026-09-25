> Modul  : Treaty In Adjustment · Ronde A · 2026-09-24
> Peran  : penilai
> Masukan: `00-LINGKUP.md` · `01-TEMUAN.md` · `02-SIDANG.md` · `03-LUBANG.md` · `04-PARITAS.md` · `05-GERBANG.md` — seluruhnya ronde A · bukti mentah pada 708 XML
> Status : TERBUKA
> Sifat  : TAMBAH-SAJA

# 06 · PUTUSAN

Ditulis dengan peran berbeda dari penemunya. Keputusan di bawah **milik pemilik proses**; peran
penilai di sini mencatatnya lengkap dengan dasar, alternatif yang ditolak, dan arah dampaknya
(`METODE` §6.7). Tiga koreksi metodologis dijatuhkan atas cara interogasi ronde ini sendiri.

---

## 1. Koreksi metodologis

### MA-01 · Precondition dibaca dengan perkakas yang belum membaca sarangnya

`pyStepsPreCondParamsWhen` tersarang di bawah `pyStepsPreCondParams/rowdata`, bukan anak langsung
baris langkah. Perkakas awal memakai `findtext` pada baris langkah, sehingga **setiap** langkah
tampak tanpa syarat. Akibatnya dua klaim dilaporkan lalu dicabut: "mesin selisih tanpa percabangan"
dan "penyalinan `ActualValue` tanpa syarat". Pelanggaran `METODE` §2.1.

Ditambah satu hal yang tidak ada di §2.0: **`pyStepsPreCondition = false` mematikan SYARATNYA
sementara langkahnya tetap berjalan** — bentuk mati kedua yang berbeda dari blok `//`. Teks
syaratnya tetap tersimpan dan terbaca meyakinkan.

### MA-02 · Pemanggil dihitung sebagai rujukan, bukan sebagai langkah hidup

TDA-07 dilaporkan tanpa memeriksa **parameter** yang dikirim pemanggilnya. Pelanggaran `METODE`
§2.0a, yang menyebut kesalahan ini sudah terjadi tiga kali di sesi Treaty In. Ia terjadi lagi di
sini, dan baru tertangkap ketika pemilik proses meminta daftar pemanggil hidupnya.

### MA-03 · Sebab disimpulkan dari dua bukti yang berdampingan, bukan dari yang menentukan

NA-02 semula menyebut "dua pengurai yang berbeda" sebagai sebab patahnya penomoran. Keduanya memang
berbeda, tetapi yang satu **nilainya dibuang** — jadi ia bukan sebab. Bentuk kesalahannya sama
dengan §3.5: dua hal yang sama-sama benar, hanya satu yang memisahkan.

### MA-04 · Salinan terbungkus dihitung sebagai isi aturan

Dua hitungan ronde ini ikut menghitung kemunculan di dalam `pyIncludedRuleXML` dan
`pyRuleVersionsList` (NA-09). Bentuk kesalahannya sama dengan `METODE` §2.0a — salinan bukan
rujukan, apalagi panggilan — dan ia lolos karena hitungannya dilakukan atas teks berkas, bukan atas
pohon aturan.

**Aturan bukti yang ditambahkan ke daftar periksa sesi ini** (`00-LINGKUP` §5a): kemunculan di
dalam kedua bungkus itu **tidak dihitung dan tidak dipakai sebagai bukti**; setiap hitungan
menyebut badan dan salinan terpisah.

---

### MA-05 · Sasaran tulis dicari dengan nama tag yang salah, dan nol dibaca sebagai jawaban

Menjawab A7 saya menyapu penulis `TreatyIn.Position` memakai tag `pyTarget`, `pyName`,
`pyPropertyTarget`. Sapuan itu mengembalikan **nol**, dan nol itu saya laporkan sebagai *"tidak ada
aturan yang menulis nilai itu"*. Sasaran tulis sebuah **Data Transform** bernama
`pyPropertiesName` / `pyPropertiesValue`; tidak satu pun tag yang saya pakai dapat menemukannya.
Nol itu bukan jawaban — ia **perkakas yang tidak menengok ke tempat yang benar**, bentuk lain dari
*"diam bukan bukti"* (`METODE` §3.1).

Kesimpulannya bertahan setelah disapu ulang dengan tag yang benar, tetapi bertahan karena
keberuntungan, bukan karena caranya.

Dua kekeliruan turunannya, keduanya dikoreksi di tempat:

* Saya menyebut `Section/TreatyInActionButtons` sebagai *"tombol yang memasukkan keadaan itu ke
  daftar peran"*. Ia **pembaca**: kondisi tampil sebuah tombol kolom, yang membaca
  `TreatyIn.Position`. Penggrill yang menangkapnya benar; tabel saya keliru.
* Saya menyebut cabang 1.4 `Akseptasi_DT` mati **"blok `//`"**. Penandanya untuk Data Transform
  adalah **`pyDisabled = true`** — bentuk "mati" yang **ketiga**, berbeda dari dua yang sudah
  dicatat `TA-01`. Kesimpulan matinya benar; nama penandanya salah, dan nama yang salah akan
  menyesatkan orang yang memverifikasi ulang.

**Yang ditambahkan ke daftar periksa:** sebelum melaporkan sebuah sapuan bernilai nol, buktikan
dulu perkakasnya dapat menemukan **kasus positif yang diketahui**. Sapuan penulis `Position` yang
benar menemukan 24 penulisan; sapuan yang salah menemukan nol. Selisih itu seharusnya terlihat
sebelum laporannya ditulis, bukan sesudah dibantah.
---

### MA-06 · Temuan disebut baru tanpa membaca ADR yang masih di daftar rekonsiliasi sendiri

Menjawab A7 saya menyusun daftar penulis `TreatyIn.Position` dan menuliskannya sebagai temuan.
**ADR-0055 §4 butir 1 sudah menetapkan hal yang sama pada 23 September 2026**, dengan sapuan yang
sama dan kesimpulan yang sama — dan ADR-0055 ada di daftar rekonsiliasi cabang A yang saya susun
sendiri, sebagai `QA-8`. Hal yang sama berlaku untuk ketiga tombol "(dev)": ADR-0055 §4 sudah
mendaftarnya, termasuk catatan bahwa `TreatyInForceEdit` terlihat oleh setiap pengguna.

Ini bentuk kesalahan yang sama dengan yang sudah saya lakukan sekali ronde ini — berargumen tentang
ADR-0052 sebelum membacanya. Kali itu tertangkap oleh penggrill; kali ini tertangkap sendiri, tetapi
**sesudah tertulis ke berkas**.

**Yang ditambahkan ke daftar periksa:** sebelum menuliskan sesuatu sebagai temuan, **baca lebih
dulu setiap ADR yang masih terbuka di daftar rekonsiliasi cabang yang sedang dikerjakan**. Daftar
itu bukan antrean pekerjaan saja — ia juga daftar tempat jawabannya mungkin sudah ada.

Akibatnya pada isi: NA-10 dan NA-11 ditandai ulang sebagai **penguatan mandiri**, bukan temuan
baru, dan bagian yang benar-benar baru dipisahkan — jejak nama `"BERNARD"`, keberadaan pintu samping
di layar addendum, `RevisionState` yang ditinggalkan `Force Resolve Complete(dev)`, dan matinya
`TreatyInSetToDirector`.
---

### MA-07 · Syarat tampil dibaca dari sel, tanpa menelusuri wadah di atasnya

Saya melaporkan `Force Edit (dev)` *"terlihat oleh semua orang"* atas dasar `pyVisible = ALWAYS`
pada **selnya**, tanpa menelusuri rantai wadah sampai akar seksi. Di atasnya ada layout `S3` dengan
`pyContainerVisibleWhen = OperatorID.pyOrgDivision = 'IT'`.

Ini persis kesalahan yang sudah saya hindari dengan benar pada NA-01 - di sana rantai tampil
diperiksa sampai akar seksi dan dilaporkan "nol syarat". Yang membedakan hanya kedisiplinannya,
bukan pengetahuannya: **prosedur yang benar sudah ada dan tidak saya jalankan.**

> **KOREKSI 25 September 2026 - MA-07 sendiri setengah keliru.** Sesudah menemukan gerbang divisi,
> saya menyimpulkan ketiga tombol "terlihat operator divisi IT", seolah gerbang wadah
> **menggantikan** syarat sel. Keduanya berlaku **bersamaan**. Arti `pyVisible` diperiksa dari
> sebarannya (553 dari 553 `OTHER` menyimpan kondisi; hanya 37 dari 6663 `ALWAYS` menyimpannya),
> dan hasilnya per tombol: `Force Resolve Complete(dev)` dan `ReturnToInputor(dev)` terlihat oleh
> **operator divisi IT yang bernama ALDO SAPUTRA atau Daniel Suhana**; hanya `Force Edit (dev)`
> yang terlihat oleh **setiap** operator divisi IT, sebab `pyVisible = ALWAYS` membuang kondisi
> namanya.
>
> Bentuk kesalahannya sama dua kali berturut-turut: **satu tingkat dibaca sebagai seluruh rantai** -
> mula-mula sel tanpa wadah, lalu wadah tanpa sel. `TA-06` karena itu diperluas: klaim keterlihatan
> menyebut **hasil perkalian** seluruh tingkat, dan menyebut nilai `pyVisible` tiap tingkat supaya
> pembacanya tahu mana yang dievaluasi.

Yang membuatnya lolos: `pyVisible = ALWAYS` terasa seperti kesimpulan akhir, padahal ia hanya
menyatakan sel itu tidak menambah syarat - bukan bahwa tidak ada syarat.

**Yang ditambahkan ke daftar periksa:** setiap klaim keterlihatan menyebut **rantainya**, bukan
satu tag - sel, baris, tabel, dan setiap layout di atasnya, dengan `pyContainerVisibleWhen` dan
`pyAssociatedPrivileges` ikut dicetak. Klaim yang hanya menyebut satu tingkat ditolak sendiri
sebelum dilaporkan.
---

### MA-08 · Pembalikan klaim yang hanya muncul di dalam sebuah diff

Dua putaran lalu saya melaporkan kontrol revisi di `Section/InputTreatyInOffer` sebagai **tidak
berpembatas** — "nol `pyVisible`, `pyCondition`, `pyVisibleWhen`, `pyDisabledWhen`, `pyPrivilege`",
dan "label tombolnya tidak terbaca" — lalu menuliskannya ke daftar eskalasi induk sebagai
*"keduanya tidak dibatasi peran dan selalu terlihat oleh siapa pun"*. Dua hari kemudian saya
menemukan yang sebaliknya, dan **menaruhnya hanya di dalam sebuah diff** — tanpa menyatakan bahwa
sebuah klaim dicabut (`METODE` §5.4).

**1. Apa yang keliru pada pemeriksaan sebelumnya.** Pemeriksaan itu menyapu **subpohon sel** untuk
nama-nama tag pembatas dan mendapat nol. Yang saya sapu adalah sel **tombolnya**; syarat perannya
tidak tersimpan sebagai `pyDisabledWhen` atau `pyPrivilege` melainkan sebagai
`pyUserData/pyCondition` berpasangan dengan `pyVisible = OTHER` — satu tingkat di dalam sel yang
sama, pada cabang `pyUserData` yang tidak masuk daftar tag yang saya cari. Bentuknya sama persis
dengan `MA-05`: **daftar tag yang tidak lengkap mengembalikan nol, dan nol dibaca sebagai
jawaban.** Dan "label tidak terbaca" berasal dari kesalahan yang sama — `pyLabel` ada di dalam sel,
tidak saya cari.

**2. Di tingkat mana pembatasnya berada.**

| Tingkat | Isi | Baris `Section/InputTreatyInOffer.xml` |
|---|---|---|
| **sel** | `pyUserData/pyVisible = OTHER` + `pyCondition` menuntut `OperatorID.pyWorkBasketList(2).pyWorkBasketName = 'ReasTreatyInAdmin'` | 268820, 269268, 269652, 270086 |
| **wadah** | layout `S168`, `pyIsVisibilityOption = CONDITION`, `pyContainerVisibleWhen = OutputParam.DATASHOW != 1` | — |
| **privilege** | `pyAssociatedPrivileges = false` di **setiap** tingkat | — |

**3. Dari mana labelnya terbaca.** `pyLabel` di dalam sel: `Edit` (268415), `View` (268955),
`Copy` (269396), `Revision` (269782).

**4. Dan pembalikan itu sendiri memuat kesalahan baru, yang saya temukan saat menelusuri
akibatnya.** Saya menulis `Edit` mengirim `revisionstate=1`, atas dasar **keberadaan kata**
`revisionstate` di dalam subpohon selnya. Dibaca sebagai pasangan nama-nilai, `Edit` mengirim
`viewstate=` dan `revisionstate=` **kosong**. Precondition `param.revisionstate==1` karena itu
tidak terpenuhi, dan `Edit` **tidak** memendekkan rantai persetujuan kontrak yang belum disetujui.
Dugaan bahwa kontrak baru dapat disetujui tanpa kepala departemen dan direktur **tidak berdasar**.

**Keluarganya sama dengan MA-05 s.d. MA-07, dan ini yang keempat berturut-turut:** sebuah
pemeriksaan yang **lebih sempit daripada klaimnya** — daftar tag yang kurang, satu tingkat rantai,
keberadaan dibaca sebagai nilai. `TA-04` sampai `TA-06` masing-masing menyerang satu sisi; yang
belum tertutup adalah sisi keempat, dan `TA-07` menutupnya:

> **`TA-07` (usulan).** Sebuah parameter dianggap terkirim hanya bila **nilainya** dibaca, bukan
> namanya. Dan setiap pencarian berbasis teks di dalam subpohon dinyatakan sebagai pencarian teks,
> bukan sebagai pembacaan struktur.

**Yang juga ditambahkan, dan ini soal cara melapor, bukan cara membaca:** sebuah pembalikan klaim
**tidak boleh** hanya muncul sebagai perubahan di dalam artefak. Ia dinyatakan terang di badan
laporan, dengan klaim lamanya dikutip. Diusulkan sebagai penajaman `METODE` §5.4.
---

## 2. Putusan

### GRL-01 — Satu model, bukan satu berkas: addendum adalah `VERSI_KONTRAK` milik spesifikasi induk

* **Label** — **PERUBAHAN**, berubah dari *"addendum berdampingan dengan kontrak, dan kontrak tidak
  pernah digantikan"*. Di sistem lama `SaveTreatyIn_EDM_Act` langkah 12-13 mati (§6.3).
* **Yang tersirat, kini dinyatakan** — kata "versi" berarti **versi yang disetujui menggantikan
  pendahulunya**, dan hilir membaca versi yang berlaku. Itu memang maksudnya; sejalan ADR-0040.
* **Keputusan**
  1. Addendum adalah `VERSI_KONTRAK`; entitas `PENYESUAIAN` dibuang.
  2. Satu ERD, satu daftar invarian, satu daftar keadaan, satu rantai persetujuan — milik
     spesifikasi induk.
  3. Susunan dokumennya adalah **bahan to-spec**, bukan keputusan grilling.
  4. **RALAT 24 Sep 2026, dari penggrill:** `NILAI_SELISIH` **tidak sementara** — ADR-0048 butir 3
     mengikat, dan bentuk sempitnya sudah diputuskan sesi ERD induk (`METODE` §8.3, §8.5). Yang
     terbuka hanya apa yang beku, bentuk, dan cakupannya (cabang E). `BESARAN_DAPAT_DISESUAIKAN`
     tetap bergantung C1.
  5. Temuan yang mengubah entitas induk dicatat sebagai **usulan untuk langkah 8-10 induk**;
     `SPEC-MODEL-DATA.md` tidak ditulis dari sesi ini.
  6. Embargo Adjustment berakhir 23 September 2026.
* **Dasar** — begitu addendum adalah versi, tidak ada entitas Adjustment yang dapat dispesifikasikan
  sendiri. Spesifikasi terpisah berarti dua daftar keadaan yang harus dijaga sama selamanya. Dan
  karena langkah 8-10 induk belum dijalankan, temuan sesi ini masih dapat masuk ke sana alih-alih
  menjadi tambalan.
* **Ditolak** — (b) spesifikasi terpisah lewat seam; (c) menunggu induk selesai.
* **Arah dampak** — memilih (a) keliru: spesifikasi induk membengkak, ongkosnya jadwal, dapat
  dibatalkan dengan memecah berkas. Memilih (b)/(c) keliru: invarian versi ditulis dua kali dengan
  bunyi berbeda, dan yang ketahuan belakangan bukan perbedaannya melainkan akibatnya di data.
* **Status** — **terkunci**; `PP-1` menunggu, hanya memengaruhi arah dampak.
* **Titipan** — cabang D (rencana penyalaan kemampuan baru), cabang I (versi berlaku untuk kontrak
  warisan).

### GRL-02 — ADR-0036 berlaku dengan penyesuaian

* **Label — dibagi dua, karena satu label menyembunyikan cara mengujinya** (`METODE` §6.5).
  Rumusan lama *"angka selisih dapat bergeser"* adalah sisa G2 yang sudah dibantah §2.1, dan
  **dicabut**: di sistem lama `OLDDATA` dan `ValueDifference` justru **beku** di dalam `JSONDATA`
  baris addendum itu sendiri.

  | Bagian | Label | Berubah dari | Sumber |
  |---|---|---|---|
  | niat "angka dasar persetujuan beku" | **PELESTARIAN** | — sistem lama sudah membekukannya | §2.1 |
  | baris addendum yang sudah disetujui **tidak boleh dapat ditimpa** | **PERUBAHAN** | tabrakan pengenal masuk cabang `UPDATE` dan menimpa baris lain sambil melapor berhasil | **TDA-01** (§4.4, §6.2) |
  | **setiap** angka ringkasan tersimpan | **PERUBAHAN** | sebagian ditulis ke `ActualValue` lalu ditimpa pada penyimpanan berikutnya | **TDA-06** (§5.5, §2.2) |
  | baris addendum yang sudah disetujui **tidak dapat disunting di tempat** | **PERUBAHAN** *(ditambahkan 25 Sep)* | `Force Edit (dev)` membuka kuncinya dan `Save EDM(dev)` menyimpannya **tanpa syarat status apa pun**; keduanya terlihat operator divisi `IT` | **NA-19** |

  Konsekuensi pengujiannya: bagian pertama **dapat** diuji terhadap data lama; tiga bagian
  berikutnya tidak. Bagian keempat memperjelas batas bagian pertama: pembekuan di sistem lama
  berlaku untuk **mesinnya** — `OLDDATA` dan `ValueDifference` memang beku di dalam `JSONDATA` —
  tetapi **tidak terlindung dari penyuntingan manual** (NA-19).
* **Keputusan**
  1. Klasifikasi: **berlaku dengan penyesuaian**.
  2. **Prinsip terkunci:** angka yang menjadi dasar persetujuan sebuah addendum — **termasuk selisih
     yang ditampilkan kepada penyetuju** — tidak berubah sesudah persetujuan; berlaku juga bila
     rumus selisih diperbaiki kemudian.
  3. **RALAT 24 Sep 2026, dari penggrill:** pertanyaan "disimpan atau dihitung saat dibaca" tidak
     pernah terbuka — ADR-0048 butir 3 sudah mengikat. Dan sesudah GRL-03, mekanismenya tinggal
     **dua**: bekukan hasil selisih beserta pemadanannya, atau versikan rumusnya.
  4. Alasan ADR-0036 ditulis ulang sebagai **usulan revisi** `REV-1`, bukan suntingan.
* **Dasar** — prinsip ini bukan tambahan atas ADR-0036 melainkan **paruh pertama kalimat
  keputusannya sendiri**: *"Angka yang menjadi dasar suatu persetujuan dibekukan pada saat
  persetujuan itu, **bersama** seluruh masukan yang membentuknya."* Usulan penggrill semula hanya
  membekukan paruh keduanya.
* **Ditolak** — (a) berlaku apa adanya; (b) bekukan dasar selisih — **salah sasaran**, karena
  memperbaiki TDA-04 mengubah **pemadanan baris**, bukan angka dasarnya.
* **Arah dampak** — tidak setangkup: salah ke arah menyimpan berarti sejumlah angka mubazir, dapat
  dibatalkan; salah ke arah tidak menyimpan berarti angka dasar persetujuan hilang, tidak dapat
  dipulihkan.
* **Status** — prinsip **terkunci**; mekanisme menunggu cabang E; premisnya menunggu bisnis lewat
  `DB-1`.

### GRL-03 — ADR-0048 berlaku dengan penyesuaian: bingkainya mati, ketiga kewajibannya hidup

| # | Butir | Label | Berubah dari |
|---|---|---|---|
| 1 | Lapis *Batas* ADR-0048 dinyatakan **mati** | **bukan perilaku — tidak diuji** | pernyataan cakupan dokumen, bukan perilaku sistem; tidak ada yang dapat diuji terhadap data lama |
| 2 | Butir 3 **mengikat**: selisih disimpan di tabel tersendiri | **PERUBAHAN** | selisih disimpan sebagai sub-pohon `ValueDifference` di dalam `JSONDATA`, dan sebagiannya tidak pernah tersimpan (TDA-06) |
| 3 | Butir 2 (SELECT, tidak menyalin) **bertahan** | **PERUBAHAN** | sistem lama menyalin potret penuh ke `TreatyIn.OLDDATA` (§2.1) |
| 4 | Butir 1 dibaca sebagai **ID versi** | **PELESTARIAN, sementara** | `OLDID` lama pun menunjuk **baris yang disesuaikan**, yang bisa berupa revisi (§9.3). Yang akan menjadi PERUBAHAN adalah bunyi harfiah aturan pemilik proses ("ID Treaty In" = kontrak), **atau** rantai linear bila cabang B memilihnya — karena pilihan "baris mana pun" hilang. **Label ditinjau ulang sesudah cabang B.** |
| 5 | Cabang E hanya memutuskan apa yang beku, bentuk, dan cakupan | — | konsekuensi butir 2 dan 3 |

* **Dasar** — yang mati adalah bingkainya, bukan ketiga hal yang diwajibkan disediakan Treaty In,
  dan ketiganya masih dipakai GRL-01 dan GRL-02. Butir 1 dan 2 bukan milik grilling untuk ditulis
  ulang: keduanya aturan bisnis pemilik proses (`METODE` §8.5).
* **Ditolak** — **(b) buka kembali butir 3 di cabang E**, ditolak karena alasannya (TDA-06) adalah
  **fakta sistem lama**, dan menurut `METODE` §3.7 fakta tidak dapat membatalkan keputusan
  rancangan. Alasan ongkos yang semula dipakai penggrill **dicabut**. **(c) pensiunkan ADR-0048**,
  ditolak karena hanya bingkainya yang mati.
* **Arah dampak** — bila butir 4 ternyata dimaksudkan ID kontrak, rantai addendum kehilangan
  bentuknya dan seluruh cabang B berubah; karena itu ia tidak dianggap terverifikasi dan masuk
  `DB-2`.
* **Status** — **terkunci**; butir 4 menunggu bisnis lewat `DB-2`, dan labelnya menunggu cabang B.

### GRL-04 — Penyesuaian bukan entitas; `PENYESUAIAN` tidak kembali

* **Label** — **PELESTARIAN** terhadap putusan `METODE` §8.3, yang ronde ini **tidak berhasil
  dibantah**.
* **Keputusan** — (a) bertahan. `PENYESUAIAN` tidak kembali. Jenis, pemohon, dan alasan menjadi
  atribut versi.
* **Dasar, dan ia penting karena bukan yang semula diajukan penggrill** — putusan ini bertahan
  **sebagai gerbang** (`METODE` §8.3: 19 properti ditelusuri satu per satu, nol fakta tersisa),
  **bukan atas dasar uji ekspor ronde ini**. Ketiga uji kardinalitas yang diajukan penggrill
  **tidak memisahkan** (`METODE` §3.4, §3.5): sistem lama tidak punya tempat untuk penyesuaian di
  luar baris versinya, jadi ia pasti satu-ke-satu — baik penyesuaian itu peristiwa bisnis tersendiri
  maupun bukan. Bentuk simpan lama tidak dapat menjawab pertanyaan tentang jenis. Uji pertama juga
  salah sasaran: `UPDATE` atas baris yang sama adalah penyimpanan draf berulang, bukan bukti
  kardinalitas.
* **Uji keempat ditolak dengan alasan tersendiri** — `EDMEffective` yang terkunci ke `Commencement`
  dipakai penggrill untuk **membuang**, dan itu arah yang berbahaya (`METODE` §3.2). Faktanya pun
  menunjuk ke arah lain: `EDMEffective` dibaca `TreatyCalculateProratePct`, sementara lima langkah
  pro rata **mati** — sehingga tanggal yang terkunci itu bisa kehendak bisnis, bisa pula akibat
  kemampuan pro rata yang tidak pernah selesai. Ekspor tidak dapat memisahkannya. Dan tanggal
  berlaku pun tidak memisahkan: ia bisa menjadi atribut versi tanpa entitas baru.
* **Sumber jawaban yang sah** — orang (`METODE` §4.6). §8.6 Q1 adalah satu-satunya pembatalnya.
* **Pembatal, dirumuskan terang** — satu dokumen addendum mengubah **lebih dari satu** kontrak atau
  versi, **atau** satu revisi menggabungkan **lebih dari satu** dokumen addendum.
* **Status** — **terkunci dari sisi ekspor**, menunggu bisnis sebagai satu-satunya pembatal.
  **Rancangan boleh berjalan di atas keputusan ini.**
* **Catatan yang menyeberang** — bila `DB-5` (tanggal berlaku) dibantah, GRL-01 **tidak** batal,
  tetapi cabang C dan E berubah: "berlaku sejak" menjadi atribut versi, dan pro rata menjadi
  kemampuan yang harus diputuskan — dibangun sebagai kemampuan baru (`METODE` §3.8) atau dibuang.
  Dititipkan ke E bersama temuan `EDMEffective` -> `TreatyCalculateProratePct`.

### GRL-05 — ADR-0040 berlaku dengan penyesuaian; butir 5 dilepas ke C1

* **Label per butir**

  | Butir ADR-0040 | Label | Berubah dari |
  |---|---|---|
  | 1 — `OLDID` pecah jadi dua hubungan; "addendum atas" menjadi hubungan versi | **PELESTARIAN** | — sejalan GRL-01 dan GRL-03 butir 4 |
  | 3 — lapisan beku lima field | **PERUBAHAN** | sistem lama **membuka tiga di antaranya di layar addendum**, dan menulis kelimanya per baris addendum |
  | 5 — materialitas diturunkan, tidak diketik | **tidak dilabeli di sini** | dilepas ke **C1** |

* **Butir 1** — berlaku apa adanya.

* **Butir 5** — **tidak diputuskan di cabang A.** Ia satu dari tiga definisi materialitas yang
  hidup berdampingan (diketik pengguna di sistem lama, diturunkan dari jenis di ADR-0049,
  diturunkan dari hasil selisih di ADR-0040), dan C1 gerbangnya karena ia juga masukan cabang E.

  **Pencabutan:** penggrill meminta bukti atas frasa *"butir 5 bertentangan dengan butir 3-nya
  sendiri"*. Frasa itu **dicabut** (`METODE` §5.1). Diperiksa ulang terhadap kutipannya: butir 3
  menetapkan **field mana yang boleh berubah**, butir 5 menetapkan **bagaimana sebuah perubahan
  digolongkan**. Keduanya berbicara tentang hal yang berbeda dan tidak saling meniadakan. Yang
  benar-benar bertentangan dengan butir 5 adalah **NA-01** (materialitas diketik) dan **ADR-0049**
  (diturunkan dari jenis) — bukan butir 3.

* **Butir 3 — berlaku dengan penyesuaian, berlabel PERUBAHAN**, dengan empat syarat yang dipenuhi
  sebagai berikut:

  **(i) Bukti sebagai daftar, diperiksa sampai akar** (`METODE` §3.9). Kontrol yang dapat disunting
  pada pohon layar addendum:

  | Field | Berkas | Bentuk | Syarat di sel, baris, tabel, `SectionBody`, `Section` |
  |---|---|---|---|
  | `Ceding` | `Section/InputTreatyInAdjustment.xml` | (dua kontrol) | **tanpa syarat apa pun** |
  | `Commencement` | `Section/InputTreatyInAdjustment.xml` | satu `pxDateTime` + satu | **tanpa syarat apa pun** |
  | `Termination` | `Section/InputTreatyInAdjustment.xml` | satu `pxDateTime` + satu | **tanpa syarat apa pun** |
  | `Ceding` | `Section/TreatyInNONProportional.xml` | `pxAutoComplete` | **tanpa syarat apa pun** |

  Karena **tidak ada syarat sama sekali**, pemeriksaan untuk kedua nilai `EDMMaterialType` terjawab
  sekaligus: materialitas **tidak** menggerbangi ketiga field ini.

  **`ProportionType`: dinyatakan TIDAK dapat dicapai dari layar addendum.** Kata "praktis" dicabut.
  Satu-satunya kontrol yang dapat disuntingnya (`pxTextInput`, `pyReadOnly = false`) ada di
  `Section/ShowSummary.xml`, dan `ShowSummary` dirujuk **hanya** oleh `InputTreatyInOffer` beserta
  harness dan flow action-nya — **bukan** oleh `InputTreatyInAdjustment`. Pada jalur penawaran
  kontrak ia terbuka; pada jalur addendum tidak.

  **`CedingID`:** tidak ada kontrol yang dapat disuntingnya — **dan itu justru masalahnya.** Dua
  dari tiga kontrol `Ceding` tidak menulis `CedingID` (NA-08, §4.6), sehingga mengubah nama cedant
  lewat layar addendum meninggalkan pengenalnya pada nilai lama. Ketiadaan kontrol **bukan** bukti
  ketiadaan penyimpangan (`METODE` §2.0). **`UA-10` karena itu mengukur nama dan ID terpisah**, dan
  mendeteksi baris yang namanya berubah tanpa ID-nya berubah.

  **Source of business:** tidak ada kontrol yang dapat disunting pada pohon layar addendum; tetap
  diukur `UA-10` (`METODE` §3.1).

  **Tambahan 25 September 2026, dari audit `TA-04`:**

  * **Hitungannya diralat `2 + 2 + 2` -> `1 + 1 + 1`.** Satu kontrol per field, bukan dua.
    Kesimpulannya tidak berubah — ketiganya tetap dapat disunting tanpa syarat apa pun.
  * **Dasarnya diperkuat (NA-17b).** Ketiga field itu juga tidak bersyarat di
    `InputTreatyInOffer`, sehingga sesudah `Revision` atau `Force Edit (dev)` membuka kembali
    kontrak yang **sudah disetujui**, ketiganya dapat disunting. Lapisan beku ADR-0040 berlaku
    antar versi, dan justru di situlah ia dilanggar.
  * **Satu kekhawatiran dibantah (NA-20).** Penggeser `Termination` otomatis
    (`TreatyInSetTreatyYear`) **ada** di seksi yang dipakai layar addendum, tetapi sel pemicunya
    `pyReadOnly = true` dan sel tanggal yang dapat disunting tidak punya peristiwa. Menyunting
    `Commencement` pada addendum **tidak** menggeser `Termination`.

  **(ii) `UA-10` mengukur kelima field**, bukan tiga. Tidak ditemukan terbuka **bukan** bukti
  datanya tidak pernah berubah (`METODE` §3.1) — dan prosedur `PEGA_M_TREATY_IN_EDM` memang menulis
  kelimanya per baris addendum.

  **(iii) Pertanyaan pokok `METODE` §1.2 diajukan**, karena membuka periode di addendum bisa
  **kebutuhan bisnis**, bukan sekadar akibat cara sistem dibangun. Masuk daftar untuk dibantah
  sebagai `DB-6` dan `DB-7`. Bila dibantah, dicatat sebagai **usulan revisi ADR-0040**, bukan
  penyimpangan diam-diam.

  **(iv) `METODE` §3.8 berlaku** — menegakkan lapisan beku menolak tindakan yang hari ini bisa
  dilakukan. Dititipkan ke cabang D (tahap peringatan dan angka pemicu blokir) dan cabang I (nasib
  addendum warisan yang melanggar, besarnya `UA-10`).

* **Ditolak** — **(b)** memutuskan butir 5 sekarang: mengunci gerbang C1 dari samping, tanpa daftar
  jenis yang masih menunggu `EXP-1`. **(c)** melepas lapisan beku karena sistem lama membukanya:
  ditolak dengan `METODE` §3.7 — itu **fakta**, dan fakta tidak membatalkan keputusan rancangan.

* **Arah dampak** — melepas lapisan beku tidak dapat dibatalkan setelah data tertulis: kontrak yang
  cedant-nya berubah di tengah jalan menjadi sah, dan tidak ada lagi dasar untuk menyatakan dua
  kontrak itu berbeda. Menahannya padahal bisnis memang memperpanjang periode lewat addendum
  berarti menolak pekerjaan yang sah — dapat dibatalkan lewat revisi ADR, dan `DB-6`/`DB-7` ada
  untuk menemukannya sebelum cut-over.

* **Status** — **terkunci** untuk butir 1 dan 3; butir 3 menunggu bisnis lewat `DB-6` dan `DB-7`.
  Butir 5 **dipindahkan** ke C1.

### GRL-06 — Klasifikasi materialitas warisan tidak dihitung ulang

* **Label** — **PELESTARIAN**.
* **Cabang** — A (rekonsiliasi ADR-0049), diangkat dari klausa "Data warisan".
* **Keputusan**
  1. **A6 ditutup dengan rujukan.** ADR-0049 hampir seluruhnya tentang materialitas — empat dari
     lima bagiannya — dan pertentangan tiga definisinya sudah tercatat di **GRL-05** dan dilepas ke
     **C1**. Tidak ada yang digali ulang di cabang A.
  2. **Satu klausanya bukan materialitas, dan dikunci di sini:** nilai materialitas warisan
     **dibawa apa adanya, beserta asal-usulnya** (ADR-0042). Ia tidak dipakai menurunkan apa pun
     dan tidak menjadi kolom hidup.
  3. **Batas larangannya, dinyatakan terang.** Yang dilarang adalah **menimpa atau mengganti**
     klasifikasi warisan. Yang **tidak** dilarang adalah **menghitung pembanding** untuk keperluan
     deteksi — dan `UA-3` memang membutuhkannya: addendum non-material yang nilai uangnya berubah.
     Hasil hitungan semacam itu **disimpan terpisah** dan tidak pernah menggantikan nilai warisan.
  4. **Klausa ini TIDAK diperluas** ke seluruh data warisan (`METODE` §3.6). Cabang I **harus**
     menjalankan aturan baru atas baris warisan setidaknya sekali — untuk menentukan **versi
     berlaku** kontrak warisan, karena sistem lama tidak pernah mengganti kontrak (GRL-01). GRL-06
     berlaku **hanya** untuk materialitas.
* **Dasar** — menjalankan aturan penurunan yang baru atas baris lama akan **mengklasifikasi ulang
  sejarah** menurut aturan yang tidak berlaku saat baris itu dibuat. Hasilnya tetap terlihat
  konsisten, sehingga kesalahannya tidak terdeteksi, dan nilai aslinya sudah tertimpa — tidak dapat
  dipulihkan.
* **Ditolak** — **(b)** menggali ADR-0049 penuh sekarang: melitigasi ulang perkara yang sudah
  diparkir, dan `EXP-1` belum datang. **(c)** menutup A6 seluruhnya dengan rujukan ke C1: C1
  memutuskan **definisi**, bukan **migrasi**, sehingga klausa ini dapat lolos tanpa pemilik.
* **Arah dampak** — mengunci di sini padahal C1 pemiliknya: satu keputusan tercatat dua kali, yang
  kedua tinggal merujuk — murah. Membiarkannya lolos: sejarah terklasifikasi ulang diam-diam,
  tidak terdeteksi, tidak dapat dipulihkan.
* **Status** — **terkunci**.
* **Titipan ke cabang I** — pertanyaan umum yang lebih besar dari materialitas: **aturan baru mana
  yang boleh dijalankan atas baris warisan dan mana yang tidak**, dan hasilnya dicatat dengan
  asal-usul yang bagaimana.
* **Titipan ke cabang C — bingkai pemasangan jenis, diperjelas.** Yang dipasangkan **bukan tiga
  dengan tiga**. Sistem lama punya **dua sumbu** — `EDMState` × `EDMMaterialType`, hingga **enam**
  kombinasi — sedangkan ADR-0049 punya **satu sumbu** berisi tiga jenis. Dua sumber memberi tahu
  kombinasi mana yang nyata: **`EXP-1`** menunjukkan yang **ditawarkan** layar, **`UA-2`**
  menunjukkan yang **pernah terbentuk**. Pemasangannya diputuskan di C, tidak sekarang.

### GRL-07 — ADR-0052 berlaku dengan penyesuaian: satu rantai, dan penyimpangan yang berjalan hanya SATU

* **Cabang** — A (rekonsiliasi ADR-0052), `QA-7`.
* **Label per keputusan** — dipisah, karena satu label menyembunyikan dua risiko yang berbeda:

  | # | Keputusan ADR-0052 | Label | Berubah dari |
  |---|---|---|---|
  | 1 | satu rantai persetujuan bawaan | **PELESTARIAN** | — sistem lama memang punya satu rantai empat tingkat yang hidup |
  | 2 | tabel perutean kosong dari pengecualian | **PELESTARIAN** | — tidak ada pengecualian yang berjalan untuk dibawa |
  | 3 | jalur penerima tugas kelompok **tidak dibawa** | **PELESTARIAN** — *diubah dari PERUBAHAN* | keadaannya **tidak dapat dimasuki**: tidak ada aturan yang menulis `Position = "ReasTreatyInGroupLeader"` — ditetapkan ADR-0055 §4.1, dikuatkan mandiri oleh NA-10 |
  | 4 | jalur revisi **tidak dipendekkan** | **PERUBAHAN** | `RevisionState = 1` memotong rantai menjadi dua tingkat, dan pada addendum nilainya **diwarisi**, bukan dipilih (NA-12, TDA-08) |

* **Keputusan**
  1. **Keempat keputusan ADR-0052 berlaku.** Tidak satu pun dibatalkan.
  2. **Konteks ADR-0052 keliru secara faktual dan dikoreksi lewat `REV-2`** (usulan revisi ADR,
     bukan suntingan diam-diam): ADR itu menyatakan ada **dua** penyimpangan yang berjalan.
     Penyimpangan yang berjalan **satu** — jalur revisi. Jalur kelompok adalah **kode untuk keadaan
     yang tidak dapat dicapai**.

     **`REV-2` lebih ringan daripada yang tampak**, dan itu harus disebut: **ADR-0055 §4 butir 1**
     sudah membuktikan hal itu pada 23 September 2026 dan sudah menyatakan bahwa dasar ADR-0052
     kini *"lebih kuat: bukan tidak ditemukan pemakaiannya, melainkan tidak ada jalan masuknya sama
     sekali"*. Yang belum dikerjakan hanyalah **ADR-0052 sendiri masih berbunyi "dua
     penyimpangan"**. `REV-2` karena itu berbentuk **rujukan silang**, bukan penyelidikan baru.
  3. **Keputusan 4 dibelah menurut sisinya**, karena risikonya tidak sama:

     | Sisi | Asal rantai pendek | Menghapusnya berarti |
     |---|---|---|
     | **addendum** | **diwarisi** dari baris yang dimuat; tidak dipilih siapa pun, dan membawa akibat kedua — addendum lahir terkunci (§7.4) | membuang kebetulan — **PERUBAHAN berisiko rendah**, dan memang harus dibuang |
     | **revisi-di-tempat kontrak** | **dipilih** lewat kontrol `Revision` di `InputTreatyInOffer`, berpasangan dengan cabang 2 `Akseptasi_DT` yang dirancang khusus untuknya — **dan disertai penguncian 96 sel uang** (NA-21, §7.7) | mungkin membuang kehendak bisnis — pertanyaan §1.2 belum terjawab |

  4. **Sisi kedua tidak ditutup di sini.** Tidak ada bukti bahwa butir 4 ADR-0052 pernah diputuskan
     atas jawaban bisnis; ADR-nya tidak menyebut sumber semacam itu. Maka ia masuk **daftar untuk
     dibantah** sebagai `DB-10` (`METODE` §4.7), dan bila dibantah dicatat sebagai **usulan revisi
     ADR-0052**, bukan sebagai penyimpangan diam-diam.

     > **Penajaman 26 September 2026, lalu DIRALAT pada giliran yang sama.** `Revision` bukan hanya
     > memendekkan rantai; ia juga mengunci sebagian kontrol lewat `ViewState = 1` (NA-21). Saya
     > sempat membaca itu sebagai kebijakan — "perubahan tanpa uang ditinjau lebih pendek". **Bacaan
     > itu gugur** sesudah kuncinya dicocokkan ke daftar 33 akar §5.6 (NA-22): **tidak satu pun dari
     > dua belas akar uang terkunci seluruhnya**, dan `LayersEDM` tidak menjaga apa yang dijaga
     > kembarannya `Layers`. Sesudah `Revision`, angka uang **tetap** dapat diubah.
     >
     > Maka sisi revisi-di-tempat **tidak** memperoleh pembenaran bisnis dari bentuk layarnya.
     > `DB-10` ditulis ulang untuk kedua kalinya tanpa mengandaikan kebijakan apa pun, dan
     > **`DB-10b` dicabut** — ia lahir dari pemisahan "niat versus kebocoran" yang ternyata tidak
     > ada niatnya.
  5. **Akibatnya ke orang ditarik ke hari pertama** (`METODE` §4.5, §3.8) dan dititipkan ke
     **cabang D**: revisi yang hari ini berhenti di kepala seksi akan naik ke kepala departemen dan
     direktur sejak cut-over. Siapa yang harus diberi tahu, dan berapa tambahan beban penyetuju,
     diputuskan di D. Ukurannya **`UA-13`**.
* **Dasar** — `METODE` §3.8 memperingatkan kemampuan mati yang **dinyalakan** terhitung sebagai
  kemampuan baru. Kebalikannya juga berlaku: kemampuan mati yang **dibuang** akan terhitung sebagai
  perubahan padahal tidak ada yang berubah. Bila Konteks dibiarkan, pembaca berikutnya akan mengira
  ada jalur kelompok yang pernah dipakai, mencari datanya, lalu merancang pemetaan untuk sesuatu
  yang tidak pernah berjalan.
* **Ditolak** — **(b)** berlaku apa adanya dengan Konteks dibiarkan: menyimpan pernyataan fakta yang
  salah di dalam ADR, dan `METODE` §3.7 melarang ADR memutuskan fakta tentang sistem lama.
  **(c)** bertentangan: tidak ada satu pun keputusannya yang gugur.
* **Arah dampak** — melabeli keputusan 3 PELESTARIAN padahal jalur kelompok sebenarnya hidup:
  sebuah jalur persetujuan yang dipakai orang hilang tanpa peringatan, dan ketahuannya baru saat
  pekerjaan seseorang tertahan — mahal, dan di hari yang paling rapuh. Karena itu label ini
  **tidak** bersandar pada diamnya kode saja: `UA-12` yang diperluas mengukur sejarahnya, bukan
  hanya nilai sekarang. Melabelinya PERUBAHAN padahal mati: pekerjaan migrasi dibesarkan untuk
  sesuatu yang nol barisnya — murah, dan dapat dibatalkan.
* **Syarat pembalikan** — bila `UA-12` butir (a) tidak nol, keadaan itu ada pada data meski tak
  dapat dimasuki lewat aturan, dan penanganannya milik **cabang I** (ADR-0054, keadaan warisan tak
  terpetakan); label PELESTARIAN tetap berdiri, sebab yang diukur adalah baris lama, bukan jalan
  masuk baru. Bila butir (b) atau (c) tidak nol, ada **sejarah persetujuan dengan tingkat yang
  tidak ada di model baru**, dan itu mengubah apa yang harus disimpan di riwayat — juga cabang I.
* **Status** — **terkunci**, dengan `DB-10` terbuka pada sisi revisi-di-tempat.

### GRL-08 — ADR-0055 berlaku dengan penyesuaian: daftar keadaannya dipakai apa adanya untuk versi addendum

* **Cabang** — A (rekonsiliasi ADR-0055), `QA-8`. Pertanyaan terakhir cabang A.
* **Keputusan**
  1. **Daftar tujuh keadaan sah + `WARISAN_TAK_TERPETAKAN` dan tiga belas perpindahan dipakai
     untuk versi addendum tanpa keadaan tambahan.** Pemetaan `Position` + `StatusAkseptasi`
     berlaku apa adanya: addendum yang baru lahir membawa `Position = "ReasTreatyInAdmin"`,
     `StatusAkseptasi = ""` (`TreatyInSetEdit` langkah 3 dan 6) -> `DRAFT`.
  2. **Pilihan "tambah keadaan terkunci" ditolak**, dan alasannya konsistensi: itu akan membawa
     `ViewState` masuk lewat pintu belakang dengan nama lain. `ViewState` **turunan**, bukan keadaan
     - dan ADR-0055 §4.2 memakai persis penalaran itu untuk membuang `DIKEMBALIKAN`. Menerapkannya
     tidak konsisten akan meruntuhkan alasan keduanya.
  3. **`REV-3` diajukan** sebagai usulan revisi ADR-0055, memuat empat butir — tiga dari analisis
     dan satu dari bagian 4 pemeriksaan tombol:

     | Butir | Isi |
     |---|---|
     | i | §4 mendaftar `TreatyInSetToDirector` sebagai pintu samping yang berjalan. **Seluruh langkah tingkat atasnya ber-blok `//`** — tidak ada satu langkah hidup pun, dan langkah 4.5 menuntut keadaan yang tidak dapat dimasuki. Pintunya ada, temboknya utuh |
     | ii | §4 menyebut "layar penawaran". Seksi `TreatyInActionButtons` juga disertakan `Section/InputTreatyInAdjustment` — **pintu samping ada di layar addendum juga** |
     | iii | §4 mencatat `TreatyInForceEdit` "terlihat oleh **setiap pengguna** layar penawaran". Yang benar: **setiap operator divisi `IT`** — layout `S3` ber-`pyContainerVisibleWhen = OperatorID.pyOrgDivision = 'IT'` (NA-16) |
     | iv | Di layar addendum, `Force Resolve Complete(dev)` memanggil prosedur **kontrak**; `UPDATE` tidak mengenai baris dan **tidak ada yang tersimpan**, sementara layar melapor berhasil (NA-15). Baris kontrak asalnya **tidak** tersentuh |

* **Label per perpindahan, bukan satu label untuk seluruh daftar** (`METODE` §6.5). Untuk jalur
  addendum:

  | Perpindahan / sifat | Label | Berubah dari |
  |---|---|---|
  | `LAHIR` -> `DRAFT`, `AJUKAN`, `SETUJUI` ×3, `KEMBALIKAN` ×3 | **PELESTARIAN** | — bentuknya sama dengan `Akseptasi_DT` cabang 1 |
  | `TOLAK` -> `DITOLAK` | **PERUBAHAN** | **TDA-02**: penolakan addendum **menghapus barisnya** dari `M_TREATY_IN_EDM` dan `TREATY_IN_EDM` (`TreatyInDeclineConfirmation_postactEDM` langkah 6-7, keduanya hidup). ADR-0055 syarat (a) melarangnya |
  | `DITOLAK` terminal | **PERUBAHAN** | di sistem lama `Decline` mengosongkan `Position`, sehingga pengajuan berikutnya masuk lagi dari Sec Head tanpa jejak |
  | rantai empat tingkat **tanpa pengecualian** | **PERUBAHAN** | **GRL-07 butir 4** / TDA-08: `RevisionState = 1` memotongnya menjadi dua, dan pada addendum nilainya **diwarisi** |
  | `DRAFT` -> `DIBATALKAN` | **BARU** | tidak ada padanannya; sistem lama memakai jalur penolakan sebagai tempat sampah |
  | `WARISAN_TAK_TERPETAKAN` -> `PERBAIKAN_WARISAN` | **BARU** | tidak ada padanannya |

  **Padanan "SecHead Reject yang menulis ulang `RevisionState = "1"`" di model baru: tidak ada, dan
  memang tidak diperlukan.** Di sistem lama nilai itu dipertahankan supaya perkara yang dikembalikan
  tetap menempuh rantai pendek pada pengajuan berikutnya. Di model baru rantainya selalu empat
  tingkat (GRL-07 butir 4), sehingga tidak ada yang perlu diingat: `KEMBALIKAN` mengembalikan versi
  itu ke `DRAFT`, dan `AJUKAN` berikutnya menempuh rantai yang sama seperti semula. **Penanda yang
  hilang tanpa pengganti** — itu akibat langsung GRL-07, bukan keputusan baru.

* **Butir (iii) analisis diterima sebagai PERUBAHAN.** Membuang `ViewState` berarti addendum yang
  **terkunci warisan** — mewarisi `RevisionState = 1`, lalu dikunci `TreatyInSetEdit` langkah 2.1
  — menjadi `DRAFT` biasa yang dapat disunting pembuatnya. Kuncinya lahir dari **penularan, bukan
  kehendak siapa pun**, sehingga membukanya dapat dipertanggungjawabkan. Dua akibat dititipkan:
  **cabang I** (nasib addendum yang masih `DRAFT` saat peralihan) dan **cabang D** (pemberitahuan
  kepada pembuatnya, `METODE` §4.5). Besarannya **`UA-9(a)`**.

* **Preseden GRL-06 diperiksa dan TIDAK berlaku di sini** (`METODE` §3.6). GRL-06 melarang
  **menghitung ulang** sebuah nilai warisan yang merupakan **keputusan manusia** — klasifikasi
  materialitas dipilih orang, dan menimpanya berarti mengarang ulang sejarah. Kelima penanda lama
  bukan itu: `ViewState`, `IsEditData`, `RevisionState`, `Position`, dan `StatusAkseptasi` adalah
  **mekanisme**, bukan keputusan — tidak seorang pun pernah "memutuskan `ViewState = 1"`, ia akibat
  sebuah langkah. Membuangnya tidak menghapus keputusan siapa pun; keputusannya tersimpan di
  `CATATAN_PERSETUJUAN` dan di keadaan versinya.

  **Tetapi dua di antaranya tetap dibutuhkan selama peralihan, dan itu harus tertulis:** `Position`
  dan `StatusAkseptasi` adalah **satu-satunya** sumber untuk memetakan baris yang sedang berada di
  tengah rantai persetujuan saat cut-over ke keadaan yang benar (tabel padanan ADR-0055). Keduanya
  **dibaca saat migrasi** dan **tidak disimpan sesudahnya** — itu bukan pengecualian terhadap
  klausa penutup ADR-0055, melainkan penjelasan bagaimana klausa itu dijalankan.

* **Ditolak** — **(b)** berlaku apa adanya: menyimpan tiga pernyataan fakta yang keliru di dalam
  ADR, dan `METODE` §3.7 melarang ADR memutuskan fakta tentang sistem lama. **(c)** daftar keadaan
  perlu tambahan: alasan konsistensi di butir 2.
* **Arah dampak** — bila daftar keadaan ternyata **tidak** cukup untuk addendum, ketahuannya di
  to-spec saat menyusun kelengkapan per-perpindahan, dan ongkosnya **menambah baris pada sebuah
  daftar** — persis seperti `DIBATALKAN` yang ditambahkan 24 September. Bila sebaliknya addendum
  diberi daur hidup tersendiri sekarang, mesin persetujuannya tergandakan, GRL-01 dan GRL-04
  keduanya terlanggar, dan pembatalannya menuntut pembongkaran model. Ongkosnya tidak setangkup;
  sisi yang murah diambil tanpa menunggu bukti.
* **Status** — **terkunci**. `REV-3` menunggu keputusan atas usulan revisi ADR.

---

## 3. Temuan turunan

### TA-01 · Dua bentuk "mati" yang berbeda, dan hanya satu ada di `METODE` §2.0

`pyStepsBlockName = "//"` mematikan **langkahnya**; `pyStepsPreCondition = false` mematikan
**syaratnya** sementara langkahnya berjalan. Yang kedua belum tercatat di tabel §2.0, dan ia sudah
menyesatkan dua kali dalam satu ronde. **Usulan: tambahkan barisnya ke `METODE` §2.0.**

### TA-02 · Dua ekspor dari satu ruleset dapat memuat versi aturan yang berbeda

NA-03. Akibatnya untuk metode: kalimat "kedua ekspor memuat aturan yang sama" harus selalu
dinyatakan **per berkas**, tidak pernah sebagai selimut. **Usulan: tajamkan `METODE` §8.1.**

### TA-04 · Bentuk "mati" yang ketiga, dan aturan sapuan-nol

Dua usulan untuk `METODE`, keduanya lahir dari `MA-05`:

1. **`TA-01` diperluas.** Penanda "mati" ada **tiga**, bukan dua: `pyStepsBlockName = "//"` pada
   langkah aktivitas, `pyStepsPreCondition = "false"` pada kondisinya, dan **`pyDisabled = true`
   pada baris Data Transform**. Ketiganya berbeda maknanya dan hanya yang pertama ada di
   `METODE` §2.0.
2. **Aturan baru, usulan untuk §3.1.** Sebuah sapuan yang mengembalikan **nol** tidak boleh
   dilaporkan sebelum perkakasnya dibuktikan dapat menemukan **kasus positif yang diketahui** pada
   korpus yang sama. Nol dari perkakas yang belum diuji bukan bukti ketiadaan — ia diam yang
   berasal dari perkakasnya, bukan dari korpusnya.
