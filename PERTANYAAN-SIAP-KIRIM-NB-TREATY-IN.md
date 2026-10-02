# Pertanyaan siap kirim — NB Treaty In

> ⛔⛔ **RALAT MENYELURUH ATAS BERKAS INI — P18 DITARIK.** `[penyimpangan sadar]` 2026-09-22
>
> Berkas ini memuat **1 pernyataan** yang bersandar pada premis ⛔ *"isi langkah penetapan nilai
> tidak ikut terekspor"* — di baris **30 — ⚠️ cuplikan lama, masih menulis *942 langkah***. **Premis itu keliru.**
>
> Isi langkah **ada di dalam ekspor**, di tag `PropertiesName`/`PropertiesValue` — **tanpa awalan
> `py`**. Tim migrasi memeriksa varian **dengan** awalan, menemukan nol, dan mempercayainya.
> ⭐ **Sebabnya dua huruf.** **938 dari 939** langkah `Property-Set` membawa isinya sendiri;
> **2.481** pasangan nama=nilai terisi; **nol** tanda terpotong. Lihat
> `VERIFIKASI-P18.md` dan `KOREKSI-P18-DIJALANKAN.md`.
>
> ⚠️ **Kalimat-kalimat itu sengaja TIDAK disunting satu per satu**, sebab sebagian berada di dalam
> **jawaban yang sudah diberikan pemiliknya** — mengubah jawaban orang lain bukan wewenang tim
> migrasi. Ralat ini berlaku atas seluruhnya.

---


**Tanggal:** 2026-09-22 · **Sumber:** grilling ronde 1, 2, dan 3

Berkas ini adalah gabungan Bab D dari ketiga ronde grilling, disusun ulang **per pemilik**
agar tiap orang menerima satu daftar, bukan tiga.

Tiap pertanyaan memuat: pertanyaannya dalam bahasa bisnis · konteks singkat · bentuk jawaban
yang diharapkan · dampak bila dijawab salah · dan satu baris `rujukan:` berisi bukti teknis
**yang boleh Anda abaikan** bila tidak diperlukan.

⛔ Tidak satu pun pertanyaan di sini dijawab sendiri oleh tim migrasi. Semuanya menunggu Anda.

## Ringkas

| Pemilik | Pertanyaan |
| --- | ---: |
| pemilik export Pega | 3 |
| DBA | 5 |
| Product + Underwriting | 11 |
| Finance dan Product + Underwriting | 3 |
| pengembang Pega lama | 13 |
| IAM | 5 |
| **Total** | **40** |

---

# Untuk pemilik export Pega — 3 pertanyaan

### ⛔ ~~**P18 — Isi 942 langkah penghitungan tidak ada di dalam berkas yang kami terima. Bisa dikirim ulang?**~~  — **DITARIK 2026-09-22**

> ⛔ **Pertanyaan ini ditarik.** Isinya ada di ekspor; angka 942 juga sudah usang *(939 seluruhnya,
> 938 berisi)*. Tidak perlu dikirim.

*(dari grilling ronde 2)*

Sistem lama punya **942 langkah** yang menetapkan nilai — hasil perhitungan premi, bagian,
pengurangan, dan pajak. Kami dapat melihat **bahwa** langkah itu ada, dan **urutannya**,
⛔ **tetapi tidak satu pun rumus atau nilai yang ditetapkannya ikut terkirim.** Yang ada hanya
salinan sementara yang tertinggal dari layar penyunting, dan salinan itu **terpotong di tengah**
sehingga tidak dapat dipakai.

**Konteks:** tanpa isi langkah-langkah ini, perhitungan uang sistem lama tidak dapat ditiru —
hanya ditebak.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — ekspor ulang dengan **isi langkah
disertakan**. Bila ekspor jenis ini memang tidak pernah menyertakannya, katakan begitu, dan
sebutkan **cara lain yang tersedia** *(cetakan layar aturan, dokumentasi, atau akses baca ke
sistem lama)*.

**Dampak bila salah:** ⛔⛔ seluruh perhitungan uang dibangun dari tebakan, dan **selisihnya baru
muncul di laporan keuangan berbulan-bulan kemudian**.

rujukan: 942 langkah `Property-Set` di 92 `Activity`; `PropertiesValue` = 74 dengan salinan editor,
**0** tanpa; `pyStepsCallParams` hanya placeholder — ronde 2 baru #1

---

---

### **P19 — Isi dua tabel keputusan juga tidak terkirim, dan satu di antaranya masih sering diubah**

*(dari grilling ronde 2)*

Ada dua **tabel keputusan** — daftar "kalau begini maka begitu" — dan ⛔ **baris-barisnya tidak
ikut terkirim.** Yang tersisa hanya namanya. Salah satunya **terakhir diubah 22 September 2025**
dan punya **delapan versi sebelumnya**, jadi ia jelas masih dipakai dan masih disesuaikan.

**Konteks:** salah satu dari keduanya adalah **pemeriksaan pertama** yang dilewati setiap
pengajuan.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — kedua tabel keputusan itu beserta
barisnya. Bila tidak bisa, **cetakan layarnya** sudah cukup.

**Dampak bila salah:** ⛔ penggolongan bisnis menjadi berbeda, dan itu menentukan **treaty mana
yang dipakai** untuk sebuah penutupan.

rujukan: `DecisionTable\isApproved.xml` · `DecisionTable\BusinessType_DeT.xml`, nol
`pyCriteriaValue`/`pyResult` — OQ-026, ronde 2 baru #3

---

---

### **P30 — Dua langkah lagi yang dipanggil tetapi tidak ikut terkirim, dan keduanya di jalur penting**

*(dari grilling ronde 3)*

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

---

# Untuk DBA — 5 pertanyaan

### ✅ ~~**P1 — Kami butuh isi tiga program penyimpan data di basis data**~~  — **TERJAWAB 2026-09-22**

> ✅ **Naskah EMPAT stored procedure diterima work owner**, ditambah dua contoh `DATA_JSON`.
> Tidak perlu dikirim. Jawabannya di `.scratch
b-treaty-in\PERTANYAAN-untuk-DBA.md` §P1.

*(dari grilling ronde 1)*

Ada tiga program yang tersimpan **di dalam basis data** — bukan di aplikasi — dan seluruh
penyimpanan data treaty inward melewatinya. Salah satunya menerima **24 keterangan sekaligus**.
Kami dapat melihat keterangan apa yang **dikirim masuk**, ⛔ **tetapi tidak dapat melihat apa yang
dilakukan program itu** — tabel mana yang diisi, kolom mana, pemeriksaan apa yang dijalankan.

**Konteks:** sistem baru harus menyimpan data yang sama, dan tanpa isi program ini kami akan
**menebak** ke mana data pergi.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — naskah ketiga program berikut namanya:
`PEGA_TREATY_IN`, `PEGA_JSON_POLIS_TREATYIN`, `PROC_GENERATE_SEQUENCE_NUMBER`.
⚠️ Untuk yang ketiga, sebagian keterangannya sudah kami terima pada **15 September 2026**; yang
belum adalah naskah lengkapnya.

**Dampak bila salah:** data treaty tersimpan di tempat atau bentuk yang berbeda dari sekarang,
dan **tidak seorang pun akan tahu sampai laporan tidak cocok**.

rujukan: `RDBList\SaveTreatyIn.xml` · `RDBList\SavePolisTreatyIn_SQL.xml` ·
`RDBList\GetSequenceNumber_SQL.xml` — OQ-002

---

---

### **P2 — Apakah ketiga program itu menyelesaikan penyimpanannya sendiri, atau menunggu aplikasi?**

*(dari grilling ronde 1)*

Ketika data disimpan, ada satu titik di mana penyimpanan menjadi **permanen dan tidak bisa
dibatalkan**. Kami menemukan bahwa perintah "jadikan permanen" itu ditulis **di sisi aplikasi**,
⛔ tetapi untuk salah satu dari tiga program itu keterangan Anda sebelumnya menyebut bahwa
**program itu sendiri tidak menjadikannya permanen**. ⚠️ Untuk dua program lainnya kami **belum
tahu**.

**Konteks:** bila sebuah program menjadikan datanya permanen sendiri **di tengah** pekerjaan,
maka kegagalan pada langkah berikutnya **tidak dapat membatalkan** yang sudah tersimpan.

**Bentuk jawaban yang diharapkan:** pilihan ganda, **untuk masing-masing dari tiga program**:
**(a)** program menjadikan permanen sendiri · **(b)** program menyerahkannya ke aplikasi ·
**(c)** tergantung jalur, jelaskan.

**Dampac bila salah:** ⛔ separuh data menjadi permanen sementara separuh lainnya batal — dan
selisihnya hanya ketahuan berbulan-bulan kemudian.

rujukan: `COMMIT` di dalam blok PL/SQL Pega vs di dalam procedure — OQ-013, butir baru #2

---

---

### **P3 — Empat nama tabel muncul dalam dua ejaan. Satu tabel atau dua?**

*(dari grilling ronde 1)*

Empat nama tabel — mata uang, jenis reasuransi, pengaturan proporsional, dan tabel riwayat —
muncul **dua kali** di dalam aturan: sekali **dengan nama gudang data di depannya**, sekali
**tanpa**. Kami tidak dapat memastikan dari aturan itu apakah keduanya menunjuk **tabel yang sama**.

**Konteks:** bila tanpa awalan berarti "gudang milik pengguna yang sedang menyambung", maka
jawabannya bisa berbeda antara lingkungan uji dan lingkungan sebenarnya.

**Bentuk jawaban yang diharapkan:** untuk tiap nama, satu baris — **nama gudang yang benar**, dan
apakah ada tabel lain bernama sama di gudang berbeda. Contoh nilai: `POOLDATA.CURRENCY` dan
`CURRENCY` ⇒ *"sama"* atau *"berbeda, yang kedua ada di gudang X"*.

**Dampak bila salah:** sistem baru **membaca tabel yang salah** atau **menulis ke gudang yang
salah**, dan keduanya tidak menimbulkan pesan galat.

rujukan: `CURRENCY` · `REINSURANCETYPE` · `PROPORTIONALARRG` · `HISTORYAKSEPTASIPEGA` /
`HISTORYAKSEPTASIPRODUCTION` — butir baru #3

---

---

### **P4 — Satu kolom pada tabel riwayat tidak pernah diisi dari modul ini. Siapa mengisinya?**

*(dari grilling ronde 1)*

Tabel riwayat akseptasi punya **tujuh kolom**. Bagian sistem yang kami periksa mengisi **enam**.
Kolom ketujuh — penampung identitas operator — ⛔ **tidak diisi sama sekali dari sini**.

**Konteks:** kami perlu tahu apakah kolom itu diisi bagian lain sistem, atau memang tidak pernah
dipakai — sebab yang tidak pernah dipakai tidak perlu dibawa ke sistem baru.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** diisi bagian lain, sebutkan mana ·
**(b)** tidak pernah dipakai, boleh ditinggalkan · **(c)** dipakai untuk data lama saja.
⭐ Bila memungkinkan: **berapa baris yang kolom itu tidak kosong** di basis data sekarang.

**Dampak bila salah:** jejak audit sistem baru **kehilangan satu keterangan** yang selama ini ada,
atau sebaliknya kami membawa kolom mati selamanya.

rujukan: `POOLDATA.HISTORYAKSEPTASIPEGA` kolom `OPERATORID` ·
`RDBList\InsertHistoryAkseptasiPega_Sql.xml` — butir baru #4

---

---

### **P29 — Sebuah kolom yang namanya berarti "jenis usaha" ternyata menyimpan seluruh data kontrak. Benarkah?**

*(dari grilling ronde 3)*

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

---

# Untuk Product + Underwriting — 11 pertanyaan

### **P5 — Apakah persetujuan Direktur dan jalur klaim benar-benar berjalan?**

*(dari grilling ronde 1)*

Di dalam gambar alur kerja ada **lima kotak** yang menggambarkan sebuah jalur: "Klaim" →
"Manajer Klaim" → "Persetujuan Direktur". ⛔ **Tidak ada satu pun garis yang masuk ke jalur itu.**
⚠️ Jadi entah jalur itu **sudah lama tidak dipakai**, atau **garisnya hilang** ketika sistem lama
diekspor.

**Konteks:** kalau jalur itu hidup, sistem baru harus punya jenjang persetujuan Direktur; kalau
mati, membangunnya berarti menambah jenjang yang tidak pernah ada.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** jalur ini masih dipakai, jelaskan
kapan · **(b)** sudah lama tidak dipakai · **(c)** tidak tahu, perlu diperiksa di sistem berjalan.
⭐ Bila (a): **berapa kali dipakai dalam setahun terakhir.**

**Dampak bila salah:** ⛔ satu jenjang persetujuan **hilang** dari sistem baru, atau sebuah jenjang
**dibangun dan diwajibkan** padahal tidak pernah ada.

rujukan: `Flow\InputRealizationTreatyIn.xml` — `Decision5`/`Assignment5`/`Assignment1`, sub-graf
tanpa connector masuk — OQ-023

---

---

### **P6 — Ada dua aturan berbeda dengan nama sama untuk memutuskan "sudah disetujui atau belum". Mana yang berlaku?**

*(dari grilling ronde 1)*

Gerbang **pertama** seluruh alur adalah pemeriksaan "apakah ini sudah benar/disetujui". Di sistem
lama ada **dua aturan terpisah dengan nama yang persis sama** untuk itu, dibuat **berselisih
sekitar 26 menit** pada hari yang sama, dan ⛔ **kami tidak dapat menentukan mana yang sebenarnya
dipakai.**

**Konteks:** enam titik percabangan memakai pemeriksaan ini; bila kami memilih yang salah, arah
seluruh alur bisa berubah.

**Bentuk jawaban yang diharapkan:** **butuh berkas atau demonstrasi** — satu contoh kasus nyata
beserta keterangan apakah ia lolos gerbang pertama, sehingga kami dapat mencocokkan. Atau, bila
diketahui: **aturan mana yang berlaku**, berbentuk daftar syarat.

**Dampak bila salah:** ⛔ pengajuan yang seharusnya **berhenti** malah **diteruskan**, atau
sebaliknya pengajuan sah **tertahan**.

rujukan: `When\isApproved.xml` dan `DecisionTable\isApproved.xml`, keduanya
`ASM-FW-GISFW-WORK!ISAPPROVED` — OQ-026, butir baru #7

---

---

### **P7 — Apa arti kode-kode yang menentukan percabangan?**

*(dari grilling ronde 1)*

Beberapa percabangan alur diputuskan oleh **kode tanpa keterangan**. Yang kami temukan:
angka **1** berarti "lolos"; teks **`TREATYINDEPTHEAD`** disimpan **di kolom nomor surat** dan
dipakai sebagai penanda arah; ⛔ kode bisnis **`40`** menghentikan satu langkah penyimpanan; dan
sebuah penanda bernilai **`EDM`** atau **`POLICY`** membedakan dua cara penyimpanan.

**Konteks:** kode-kode ini menggerakkan alur, jadi artinya menentukan perilaku, bukan tampilan.

**Bentuk jawaban yang diharapkan:** satu baris per kode — **apa artinya**, dan **nilai lain apa
yang mungkin muncul**. Contoh: `40` ⇒ *"lini usaha …, dikecualikan karena …"*.

**Dampak bila salah:** ⛔ pengajuan diarahkan ke jenjang yang salah, atau **dilewatkan dari
penyimpanan** tanpa ada yang menyadarinya.

rujukan: `.IsApproved`, `LetterNo`, `.Quotation.BusinessCode != "40"`, `param.isFOR` — OQ-020

---

---

### **P8 — Apa yang sebenarnya dikirim ke layanan "ARASAPAS" di akhir proses?**

*(dari grilling ronde 1)*

Langkah **terakhir** jalur realisasi memanggil sesuatu bernama "HIT SERVICE ARASAPAS".
⛔ **Isi langkah itu tidak ada di dalam ekspor yang kami terima untuk modul ini** — hanya
pembungkusnya. Salinan isinya ada di **dua modul lain**, dan ⛔ **kami sengaja tidak meminjamnya**,
sebab keduanya terdaftar berisi berbeda.

**Konteks:** ini efek keluar terakhir — bila ia mengirim data ke sistem keuangan, kegagalannya
berarti transaksi tidak tercatat di sana.

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** kirim ke sistem
faktur/keuangan, sebutkan sistemnya · **(b)** kirim ke sistem lain · **(c)** sudah tidak dipakai.
⭐ Dan: **apa yang terjadi sekarang bila pengiriman itu gagal** — dicoba ulang, atau diabaikan?

**Dampak bila salah:** ⛔ transaksi **tidak sampai** ke sistem tujuan, atau ⛔ **terkirim dua kali**.

rujukan: `Flow\InputRealizationTreatyIn.xml` `Utility2` →
`Activity\serviceInsertArasapas_act.xml` — OQ-025 *(BLOCKER)*

---

---

### **P9 — Mengapa modul realisasi treaty masuk menyentuh data treaty KELUAR?**

*(dari grilling ronde 1)*

Modul yang menangani **treaty masuk** ternyata menyentuh tabel **treaty keluar** di **delapan
berkas** — sementara modul yang **namanya** treaty keluar tidak menyentuhnya sama sekali.

**Konteks:** kami perlu tahu apakah ini memang alur bisnis *(misalnya bagian treaty masuk
diteruskan keluar)*, atau peninggalan sejarah.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang satu proses, jelaskan ·
**(b)** peninggalan, tidak dipakai lagi · **(c)** dipakai untuk laporan saja.

**Dampak bila salah:** tanggung jawab treaty keluar **dibangun di tempat yang salah**, dan
perubahan di satu sisi merusak sisi lain.

rujukan: `POOLDATA.M_TREATY_OUT`, class `…INT-TREATYOUTDETAIL` — OQ-022

---

---

### **P10 — Apakah semua pekerjaan masuk ke antrean BERSAMA, tidak pernah ke antrean pribadi?**

*(dari grilling ronde 1)*

Keenam titik penugasan di alur ini memakai **antrean bersama** — ⛔ **tidak satu pun** menugaskan
ke **kotak masuk pribadi** seseorang.

**Konteks:** ini menentukan bentuk kotak masuk di sistem baru: daftar bersama yang siapa pun dalam
peran itu boleh ambil, atau tugas milik satu orang.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** benar, semua antrean bersama ·
**(b)** seharusnya ada yang pribadi, sebutkan tahap mana · **(c)** dulu pribadi, diubah.

**Dampak bila salah:** pekerjaan **menumpuk tanpa pemilik**, atau sebaliknya **terkunci pada satu
orang** yang sedang tidak ada.

rujukan: keenam Assignment `<pyImplementation>WorkBasket` / `<pyRouteTo>Custom` — OQ-028

---

---

### **P35 — Bila tanggal akhir kontrak dibiarkan kosong, sistem mengisinya dengan HARI INI. Benarkah begitu?**

*(dari grilling ronde 3)*

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

---

### **P36 — Ada DUA penanda persetujuan, bukan satu. Apa bedanya?**

*(dari grilling ronde 3)*

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

---

### **P37 — Dua kode angka menandai kontrak "punya penempatan keluar". Apa artinya?**

*(dari grilling ronde 3)*

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

---

### **P38 — Sebuah aturan lini jiwa hidup di dalam modul kontrak reasuransi masuk. Disengaja?**

*(dari grilling ronde 3)*

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

---

### **P39 — Satu mata uang dikecualikan secara khusus dari dua daftar pilihan. Masih perlu?**

*(dari grilling ronde 3)*

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

---

# Untuk Finance dan Product + Underwriting — 3 pertanyaan

### **P20 — Dua puluh nomor polis tertulis langsung di dalam aturan penghitung uang. Apakah itu disengaja?**

*(dari grilling ronde 2)*

⛔ Di dalam tiga aturan sistem lama terdapat **dua puluh nomor polis tertentu yang ditulis
langsung**, bukan dibaca dari data. ⚠️ **Dua dari tiga aturan itu menghitung nilai pertanggungan
dan premi.** Nomor-nomornya bertanggal **2023, 2024, dan 2025** — jadi ditambahkan bertahap
selama tiga tahun.

**Konteks:** artinya dua puluh penutupan itu **dihitung dengan cara berbeda** dari semua penutupan
lain, dan perbedaannya tidak tercatat di mana pun selain di dalam aturan itu.

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** memang perlakuan khusus
yang disengaja, jelaskan alasannya · **(b)** tambalan sementara untuk memperbaiki data yang salah,
sudah tidak perlu · **(c)** tidak tahu, perlu diperiksa satu per satu.
⭐ Bila (a) atau (c): ⭐ **kami perlu daftar itu dibaca bersama**, sebab sistem baru harus tahu
apakah dua puluh polis ini tetap istimewa.

**Dampak bila salah:** ⛔⛔ dua puluh penutupan berpindah ke perhitungan yang berbeda dari
sekarang, **tanpa satu pun pesan galat** yang memberi tahu.

rujukan: `Activity\SumTSIPremiSpreadedRNM_Act.xml` 24× · `Activity\ProtectFIREMBUPA_Act.xml` 12× ·
`When\IsErrorSpreading.xml` 11×; 47 kemunculan, 20 nomor unik — ronde 2 baru #2

---

---

### **P21 — Dua aturan uang hanya tertulis sebagai catatan. Apakah keduanya benar-benar berlaku?**

*(dari grilling ronde 2)*

Di dalam catatan pengembang pada langkah perhitungan, ada **dua aturan uang** yang dinyatakan
sebagai kewajiban: bahwa **premi polis induk harus nol** dalam keadaan tertentu, dan bahwa
**total bagian premi yang diberikan kepada pemberi bisnis harus sama dengan premi perusahaan**.
⛔ Keduanya **hanya ada sebagai catatan** — kami tidak dapat memastikan apakah sistem benar-benar
menolak ketika keduanya dilanggar.

**Konteks:** kalau ini benar-benar kewajiban, sistem baru harus menegakkannya; kalau bukan,
membangunnya akan menolak data yang sah.

**Bentuk jawaban yang diharapkan:** untuk masing-masing, pilihan ganda — **(a)** kewajiban keras,
harus ditolak bila dilanggar · **(b)** peringatan saja, boleh dilanjutkan · **(c)** bukan aturan,
hanya catatan lama. ⭐ Bila (a): **siapa yang boleh menyimpang**, kalau ada.

**Dampak bila salah:** ⛔ ketidakseimbangan premi **lolos tanpa terdeteksi**, atau sebaliknya
penutupan yang sah **ditolak**.

rujukan: `pyStepsDescription` — `set protect kalau master polis premi harus 0` ·
`TOTAL PREMI SHARE CEDANT HARUS = PREMI RNM` — ronde 2 baru #10

---

---

### **P22 — Sembilan puluh satu bagian layar dibuat agar tidak pernah muncul, dan tiga puluh delapan dikunci permanen. Masih perlu?**

*(dari grilling ronde 2)*

⛔ **Sembilan puluh satu** bagian layar diberi syarat tampil yang **tidak mungkin pernah benar** —
setara dengan menuliskan "tampilkan jika 1 sama dengan 2". Dan ⛔ **tiga puluh delapan** medan
diberi syarat **hanya-baca yang selalu benar**, sehingga tidak pernah dapat disunting.
⭐ **Seluruh tiga puluh delapan itu ada di layar Kepala Departemen.**

**Konteks:** kami perlu tahu mana yang memang tidak boleh dipakai lagi, dan mana yang dimatikan
sementara lalu terlupakan — sebab yang pertama tidak perlu dibangun, yang kedua perlu.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang sudah tidak dipakai, boleh
hilang · **(b)** dimatikan sementara, seharusnya kembali · **(c)** campuran, perlu ditinjau
bersama. ⭐ Untuk layar Kepala Departemen: **apakah jabatan itu memang hanya boleh melihat, tidak
mengubah?**

**Dampak bila salah:** ⛔ medan mati ikut dibangun dan dirawat selamanya, atau medan yang masih
dibutuhkan **hilang** dari sistem baru.

rujukan: `1=2` 72× · `NEVER` 15× · `1==2` 2×; `pyReadOnlyCondition` `1==1` 32× + varian 6×,
seluruhnya di `DetailDeptHeadTreatyIn_UW` dan `GeneralDeptHeadTreatyIn_UW` — ronde 2 baru #6

---

---

# Untuk pengembang Pega lama — 13 pertanyaan

### **P14 — Satu wadah keterangan dipakai bersama oleh 18 aturan, dan isinya berbeda-beda. Bagaimana cara membacanya?**

*(dari grilling ronde 1)*

Ada sebuah **wadah keterangan bernomor** — slot 1, slot 2, slot 3, dan seterusnya — yang dipakai
**bersama oleh 18 aturan berbeda**. ⛔ **Slot yang sama membawa hal yang berbeda di aturan yang
berbeda**: slot ke-3 membawa **seluruh muatan data dalam bentuk teks** di satu aturan, dan membawa
**posisi/jabatan** di aturan lain.

**Konteks:** kami hendak memetakan setiap slot ke kolom yang tepat. Bila pemetaannya dibuat satu
kali untuk semua, ⛔ **ia akan salah di sebagian besar aturan**.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang per-aturan, tidak ada makna
tetap · **(b)** ada dokumen pemetaan, **butuh berkas** · **(c)** ada aturan tak tertulis,
jelaskan. ⭐ Bila (a): **aturan mana yang mengisi slot itu sebelum penyimpanan dijalankan.**

**Dampak bila salah:** ⛔⛔ satu pemetaan yang salah **merusak enam aturan lain** yang memakai slot
yang sama — dan kesalahannya **tidak menimbulkan pesan galat**, hanya data yang salah tempat.

rujukan: halaman `InputData` 15 slot / 18 rule; `InputData.CARI3` di 7 rule —
OQ-059, butir baru #1

---

---

### **P15 — Sebuah fungsi yang membentuk seluruh muatan data tidak ada naskahnya. Apa isinya?**

*(dari grilling ronde 1)*

Ada satu fungsi buatan sendiri yang **mengubah data satu kasus menjadi satu teks besar**, dan teks
itulah yang disimpan ke basis data sebagai muatan utama. ⛔ **Naskah fungsi itu tidak ada di dalam
ekspor mana pun.**

**Konteks:** data lama tersimpan dalam bentuk yang dihasilkan fungsi ini. Sistem baru harus dapat
**membacanya kembali**.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — naskah fungsi tersebut, **atau** beberapa
**contoh nyata** hasilnya dari basis data *(3–5 contoh sudah cukup untuk menurunkan bentuknya)*.

**Dampak bila salah:** ⛔ **seluruh data lama tidak terbaca** oleh sistem baru.

rujukan: `@ASM.GetPageJSONString()` di `Activity\SaveJsonPolisTreatyIn_Act.xml` langkah 6 — OQ-012

---

---

### **P16 — Satu perhitungan membaca tabel internal Pega secara langsung. Apa penggantinya?**

*(dari grilling ronde 1)*

Sebuah perhitungan membaca **tabel kerja internal sistem lama** secara langsung, seolah tabel
biasa. ⛔ **Tabel itu akan hilang** begitu sistem lama ditinggalkan.

**Konteks:** apa pun yang bergantung pada perhitungan itu akan **berhenti bekerja tanpa pesan
galat**.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** perhitungan ini masih dipakai,
jelaskan untuk apa · **(b)** sudah tidak dipakai · **(c)** tidak tahu. ⭐ Bila (a): **angka apa
yang sebenarnya dicari** dari situ.

**Dampak bila salah:** sebuah angka di layar menjadi **selalu nol** dan tidak ada yang menyadarinya.

rujukan: `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` di `RDBList\GetCountClaim.xml` — butir baru #5

---

---

### **P17 — Aturan yang hanya hidup di layar belum kami baca sama sekali. Apa yang ada di sana?**

*(dari grilling ronde 1)*

⛔ **Kami belum membuka satu pun dari 31 berkas layar** modul ini — dan **kelima berkas terbesar
seluruh modul ada di antaranya**. Artinya pemeriksaan yang dilakukan **di layar** — medan wajib,
tombol yang hanya muncul untuk jabatan tertentu, perhitungan yang berjalan saat mengetik — ⛔ **belum
terbaca.**

**Konteks:** kami menyatakan ini terang-terangan agar tidak ada yang menyangka ronde ini sudah
meliput modulnya. ⭐ Ronde berikutnya perlu arahan **layar mana yang paling penting**.

**Bentuk jawaban yang diharapkan:** ⭐ **urutan prioritas** — dari daftar layar utama
*(input polis treaty, rincian polis, rincian non-proporsional, layar Kepala Departemen)*, mana yang
**paling banyak memuat aturan bisnis**, bukan hanya tampilan.

**Dampak bila salah:** ronde berikutnya membaca **8 MB layar** dan menemukan sebagian besarnya
tata letak, sementara aturan yang penting tetap terlewat.

rujukan: 25 `Section` + 6 `Harness`, **nol dibuka**; lima terbesar semuanya `Section` — butir baru #6

---

---

### **P23 — Pada sepertiga penggolong, keterangan yang tertulis berbeda dari yang dikerjakan aturan. Mana yang benar?**

*(dari grilling ronde 2)*

Ada 75 aturan kecil yang menggolongkan jenis bisnis. Masing-masing punya **dua hal**: keterangan
yang ditulis manusia, dan syarat yang benar-benar dijalankan. ⛔ **Pada 18 dari 52 yang punya
keduanya — sepertiga — keduanya tidak berbicara tentang hal yang sama.** Satu contoh: keterangannya
berbunyi *"Kode Bisnis 02 dan 58"*, sedangkan yang dijalankan memeriksa **jenis** bisnis bernama
*"Bonding"*.

**Konteks:** kami menulis spesifikasi dari aturan ini. Bila kami membaca keterangannya, kami salah
pada sepertiga kasus; bila kami membaca yang dijalankan, mungkin keterangannya yang mencerminkan
maksud sebenarnya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** yang dijalankan selalu benar,
keterangannya basi · **(b)** keterangannya menyatakan maksud, yang dijalankan mungkin bug ·
**(c)** kasus per kasus, ⭐ **perlu ditinjau bersama 18 aturan itu**.

**Dampak bila salah:** ⛔ sistem baru menggolongkan bisnis **berbeda dari sistem lama**, dan
penggolongan itu menentukan treaty, premi, dan laporan.

rujukan: 52 `When` punya `pyConditionString` + bentuk bertanda kurung; **18 tidak cocok (35 %)`;
contoh `When\IsBondingAndCustomBonds.xml` — ronde 2 baru #4

---

---

### **P24 — Bendera "sudah disetujui" dibandingkan sebagai teks di satu tempat dan sebagai angka di tempat lain. Mana yang dimaksud?**

*(dari grilling ronde 2)*

Pemeriksaan **pertama** yang dilewati setiap pengajuan membaca sebuah penanda "sudah disetujui".
⛔ Di satu tempat penanda itu dibandingkan **sebagai teks** — nol di dalam tanda kutip — dan di
tempat lain **sebagai angka** — nol tanpa kutip.

**Konteks:** bila penanda itu belum pernah diisi, "kosong" dan "nol" adalah dua hal berbeda dalam
satu perbandingan dan satu hal yang sama dalam perbandingan lainnya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** penanda ini angka, pembandingan teks
adalah kekeliruan · **(b)** penanda ini teks · **(c)** boleh keduanya, tidak pernah jadi masalah.
⭐ Dan: **apa arti penanda ini saat belum pernah diisi** — belum diperiksa, atau ditolak?

**Dampak bila salah:** ⛔ pengajuan yang belum diperiksa **dianggap ditolak**, atau sebaliknya
**diteruskan** seolah sudah disetujui.

rujukan: `.IsApproved = '0'` berkutip 2× vs `.IsApproved==0` telanjang 2× di layar;
`When\isApproved.xml` menguji `= 1` — ronde 2 baru #5

---

---

### **P25 — Wewenang di beberapa tempat ditentukan oleh POSISI seseorang dalam sebuah daftar. Apakah urutan itu dijamin?**

*(dari grilling ronde 2)*

Beberapa aturan memutuskan apa yang boleh dilihat seseorang dengan melihat **antrean kerja
nomor satu** atau **nomor dua** dalam daftar antrean pengguna itu — ⛔ **bukan dengan menyebut
nama antreannya.**

**Konteks:** kalau seseorang ditambahkan ke antrean baru, atau daftarnya diurutkan ulang,
⛔ **perilaku aturan berubah tanpa ada yang menyentuh aturan itu.**

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** urutannya memang dijamin, jelaskan
bagaimana · **(b)** tidak dijamin, ini memang rapuh · **(c)** tidak tahu. ⭐ Bila (a) atau (c):
⭐ **antrean mana yang seharusnya dimaksud** pada posisi satu dan posisi dua.

**Dampak bila salah:** ⛔ wewenang seseorang **berpindah diam-diam** ketika keanggotaan antreannya
berubah karena alasan yang sama sekali lain.

rujukan: `pxRequestor.OperatorID.pyWorkBasketList(1)` 3× di `When`;
`OperatorID.pyWorkBasketList(2).pyWorkBasketName` 1× di layar — ronde 2 baru #7

---

---

### **P26 — Satu nilai pengaturan tampak salah ketik. Apakah ia bekerja?**

*(dari grilling ronde 2)*

Di empat tempat, sebuah pengaturan yang mengunci medan diberi nilai yang **tampaknya salah
ketik** — satu huruf hilang dari kata yang seharusnya. Di tempat lain nilai yang sama ditulis
dengan benar.

**Konteks:** kalau sistem lama **mengabaikan** nilai yang salah ketik itu, maka keempat medan itu
sebenarnya **tidak pernah terkunci** — padahal seseorang berniat menguncinya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** salah ketik, penguncian tidak bekerja ·
**(b)** nilai itu sah, bekerja normal · **(c)** tidak tahu, perlu dicoba di sistem lama.

**Dampak bila salah:** ⚠️ medan yang seharusnya terkunci **dapat disunting**, atau sebaliknya kami
mengunci medan yang selama ini bebas.

rujukan: `pyDisabledNew` = `truewhn` 4×, di samping `true` 17× dan `always` 12× — ronde 2 baru #8

---

---

### **P27 — Wadah slot bernomor: sebagian artinya tertulis di catatan. Apakah catatan itu masih benar, dan di mana sisanya?**

*(dari grilling ronde 2)*

Melanjutkan pertanyaan sebelumnya tentang wadah keterangan bernomor: ⭐ **kami menemukan sebagian
artinya tertulis di dalam catatan pengembang** — slot 44 untuk pengurangan, 31 untuk mata uang,
47 untuk saldo terutang, 13 untuk premi neto, 32 untuk premi bruto.

**Konteks:** catatan bisa basi tanpa ada yang tahu, dan kami hanya menemukan arti **sebagian**
slot — bukan semuanya.

**Bentuk jawaban yang diharapkan:** ⭐ **konfirmasi + kelengkapan** — apakah kelima arti di atas
masih benar, dan **di mana arti slot lainnya dapat dibaca**. Bila tidak ada tempatnya,
katakan begitu.

**Dampak bila salah:** ⛔ nilai uang masuk ke kolom yang salah, dan **kesalahannya tidak
menimbulkan pesan galat** — hanya angka yang salah tempat.

rujukan: `pyStepsDescription` — `Deduction (CARI44[deduction])` ·
`Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])` ·
`Gross Premi (CARI32, CARI31[CURRENCY])` — OQ-059, ronde 1 baru #1, ronde 2 pendalaman

---

---

### **P31 — Ketika data kontrak gagal dibaca, sistem hanya mencatat di log lalu LANJUT. Disengaja?**

*(dari grilling ronde 3)*

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

---

### **P32 — Dua cara penulisan tanggal dipakai bersamaan. Mana yang benar di mana?**

*(dari grilling ronde 3)*

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

---

### **P33 — Satu medan "nama operator" diisi dari dua sumber berbeda. Mana yang dimaksud?**

*(dari grilling ronde 3)*

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

---

### **P34 — Tujuh belas langkah MENGHAPUS pesan kesalahan. Kenapa?**

*(dari grilling ronde 3)*

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

---

# Untuk IAM — 5 pertanyaan

### **P11 — Peran seseorang disimpan di kolom NOMOR TELEPON. Di mana peran yang sebenarnya?**

*(dari grilling ronde 1)*

⛔ Satu-satunya penentu peran yang dapat kami baca adalah **kolom nomor telepon** pada data
pengguna. Isinya bukan nomor telepon, melainkan **kode seperti `TREATY1` dan `SPVTREATY1`**.
Percabangan alur diputuskan dari kode itu.

**Konteks:** sistem baru butuh daftar peran yang sebenarnya. Kami tidak dapat menurunkannya dari
kolom yang dipakai untuk hal lain.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh daftar** — peran apa saja yang ada di proses treaty
inward, **siapa saja pemegangnya sekarang**, dan **apa arti tiap kode** yang muncul di kolom itu
*(`TREATY1`, `TREATY2`, `SPVTREATY1`, `SPVTREATY2`)*. ⚠️ Serta: **siapa yang berwenang mengubah
isi kolom itu hari ini.**

**Dampak bila salah:** ⛔ orang yang tidak berwenang **dapat menyetujui**, atau orang yang
berwenang **terkunci di luar**.

rujukan: `OperatorID.pyTelephone` — `When\IsTreaty1.xml`, `When\IsSPVTreaty1.xml` — OQ-027

---

---

### **P12 — Ada nama orang tertentu tertulis di dalam aturan. Apakah itu masih benar?**

*(dari grilling ronde 1)*

Di beberapa aturan, wewenang diputuskan dengan **mencocokkan nama orang tertentu** yang tertulis
langsung di dalam aturan itu. ⛔ **Nilainya tidak kami salin ke berkas mana pun.**

**Konteks:** aturan seperti ini berhenti bekerja saat orangnya pindah jabatan atau keluar, dan
tidak seorang pun mendapat pemberitahuan.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** ganti dengan peran, sebutkan perannya ·
**(b)** memang harus orang tertentu, jelaskan alasannya · **(c)** sudah tidak relevan, hapus.

**Dampak bila salah:** ⛔ wewenang tetap menempel pada orang yang sudah **tidak berhak**.

rujukan: `pyWorkPage.pxCreateOperator` di `When\IsSPVCreate.xml`;
`OperatorID.pyUserIdentifier` di 12 berkas — OQ-021

---

---

### **P13 — Siapa yang memegang pekerjaan di tiap tahap?**

*(dari grilling ronde 1)*

Alur menyebut **lima nama antrean** — Admin, Kepala Seksi, Ketua Kelompok, Kepala Departemen, dan
Direktur. ⛔ Tetapi **hubungan antara tahap dan antrean tidak tertulis** — penentuan antreannya
dilakukan secara khusus, dan kolom yang biasanya menyimpannya **kosong**.

**Konteks:** tanpa pemetaan ini, kami tidak dapat menyatakan siapa memegang giliran di tahap mana.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh tabel sederhana** — untuk setiap tahap
*("Input Realisasi", "Akseptasi Kepala Treaty", "Akseptasi Kepala Departemen")*, **antrean mana**
yang menerimanya.

**Dampak bila salah:** pekerjaan masuk ke **kotak yang salah** dan tidak ada yang merasa
bertanggung jawab.

rujukan: `<pyWorkBasket>` kosong, `<pyRouteTo>Custom`; nama antrean di `<pyRuleName>` — OQ-024

---

---

### **P28 — Nama orang juga dipakai di layar, lewat medan keempat yang belum pernah kami laporkan**

*(dari grilling ronde 2)*

Selain tempat-tempat yang sudah kami tanyakan, ⛔ **nama orang tertentu juga menentukan apa yang
muncul di layar** — di **empat layar**, **dua belas tempat**, memuat **empat nama berbeda**.
⛔ **Nilainya tidak kami salin ke berkas mana pun.** Dan salah satu layar memakai **medan identitas
yang berbeda lagi** dari yang pernah kami laporkan.

**Konteks:** pada dua layar, nama yang sama dipakai **dua arah** — ada bagian yang muncul hanya
untuk orang itu, dan bagian lain yang muncul untuk **semua kecuali** orang itu. ⭐ Mengeluarkan
orang itu mengubah **dua** perilaku sekaligus.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** ganti dengan peran, sebutkan perannya ·
**(b)** memang harus orang tertentu, jelaskan alasannya · **(c)** sudah tidak relevan.
⭐ Dan: **apa yang seharusnya terjadi** pada bagian layar itu untuk orang lain.

**Dampak bila salah:** ⛔ bagian layar **hilang** bagi orang yang berhak, atau **terlihat** oleh
yang tidak berhak.

rujukan: `pyUserIdentifier` di `DetailDeptHeadTreatyIn_UW` 3× · `GeneralDeptHeadTreatyIn_UW` 3× ·
`ListSuggest` 4×; ⚠️ `pxInsName` di `DetailPoliciesNonProportional` 2× — OQ-021 diperluas,
ronde 2 baru *(pola guard keempat)*

---

---

### **P40 — Nama seseorang tertulis di dalam teks pemberitahuan yang dibaca pengguna**

*(dari grilling ronde 3)*

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

---
