# Pertanyaan untuk **DBA**
## Migrasi Treaty Inward — Realisasi & Endorsement · modul *NB Treaty In*

Anda pemegang basis data Oracle. Pertanyaan di bawah hanya dapat dijawab dari sana — tidak satu pun dapat kami turunkan dari berkas sistem lama.

> ⭐ **8 pertanyaan untuk Anda** *(⭐ 2 baru — P61 dan P62, lahir dari naskah procedure yang Anda kirim)*. ⛔ Sisanya bukan urusan Anda dan tidak disertakan.
>
> ⭐⭐ **P1 dan P29 SUDAH TERJAWAB, dan keduanya tidak lagi menahan apa pun.** `[terverifikasi]` 2026-09-22 — naskah **empat** stored procedure dan **dua contoh** `DATA_JSON` diterima. **Penahan proyek kini nol.**
>
> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat ini semula berbunyi:
> > *"**P1** dan **P29** menahan seluruh penulisan spesifikasi penyimpanan data. Tiga sisanya penting tetapi tidak menahan."*
>
> ⭐ **Diurutkan dari yang paling menahan pekerjaan**, bukan menurut nomor.

## Cara menjawab

- ⭐ **Jawab langsung di bawah tiap pertanyaan.** Tidak perlu format khusus.
- ⭐ **Baris `rujukan:` paling bawah tiap pertanyaan boleh Anda abaikan** — itu catatan teknis
  untuk tim migrasi, bukan bagian dari pertanyaan.
- ⚠️ **Tidak tahu adalah jawaban yang sah**, dan lebih berguna daripada perkiraan. Tulis
  *"tidak tahu"* dan sebutkan siapa yang mungkin tahu.
- ⚠️ **Jangan menebak.** Sebuah tebakan yang masuk ke spesifikasi akan dibangun sebagai fakta.
- ⛔ Nomor pertanyaan **jangan diubah** — ia dipakai untuk melacak jawaban ke temuan aslinya.

---

## 1. P1 — Kami butuh isi tiga program penyimpan data di basis data

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

**Jawaban:**

> **Naskah keempat procedure diterima dari work owner.** `[data DBA]` `[terverifikasi]` 2026-09-22
>
> Butir ini **tertutup**. Penahan terakhir proyek NB Treaty In lepas.
>
> ---
>
> ## 1 · `POOLDATA.PEGA_TREATY_IN` — menulis DUA tabel, bukan satu
>
> | Tabel | Kolom |
> | --- | --- |
> | `M_TREATY_IN` | `ID`, `JSONDATA` |
> | `POOLDATA.TREATY_IN` | `ID` + **20 kolom datar** |
>
> Kedua puluh kolom itu: `PROPORTIONTYPE` `TREATYCONTRACTNAME` `TERITORIALSCOPE` `COMMENCEMENT`
> `TERMINATION` `CLASSOFBUSINESS` `LEADINGREINSSOURCE` `LEADINGREINSSOURCEID` `CEDING` `CEDINGID`
> `LEADINGREINSID` `NUSARESHAREPCT` `BROKERAGEPCT` `INFORMATION` `POSITIONUSERNAME` `POSITION`
> `STATUSAKSEPTASI` `CHOOSESTATUSAKSEPTASI` `TREATYYEAR`.
>
> ⭐ **Data kontrak treaty SUDAH relasional.** Yang belum hanyalah data polis. Ini memperkecil
> pekerjaan perancangan tabel secara berarti — lihat dampaknya pada **P29**.
>
> **Pembentukan ID:** `M_SITE_DATABASE.ID` (baris ber-`CURRENT_SITE='1'`) disambung
> `M_TREATY_IN_SEQ.nextval` berbantalan nol enam digit.
>
> **Sisip atau perbarui** ditentukan nilai penanda: `IDPega = 'UnknownId'` berarti sisip, selain itu
> perbarui. Pada penyisipan, teks `'UnknownId'` di dalam JSON **diganti** dengan ID yang baru
> dibentuk.
>
> `COMMIT` **di dalam**, pada kedua cabang. Galat menghasilkan `ROLLBACK` dan `StsSave = 0`.
>
> ---
>
> ## 2 · `POOLDATA.PEGA_JSON_POLIS_TREATYIN` — satu tabel, delapan kolom
>
> `POOLDATA.json_polis` — `IDPEGA`, `DATA_JSON` (CLOB), `TGL_INPUT`, `NOPOLIS`, `NOENDORS`,
> `PRODKE`, `TGL_PROD`, `USERNAME`.
>
> Hanya **menyisip**, tidak pernah memperbarui. `TGL_INPUT` diisi `SYSDATE`; `TGL_PROD` diurai dari
> teks berformat `DD/MM/YYYY HH24:MI:SS`.
>
> `COMMIT` **di dalam**.
>
> ⭐ **Inilah tabel yang menyimpan data polis sebagai dokumen** — sasaran penggantian oleh tabel
> datar. Tujuh kolom lainnya sudah datar dan dapat dipertahankan apa adanya.
>
> ---
>
> ## 3 · `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`
>
> Membaca `POOLDATA.TANGGAL_CLOSING` dengan `SELECT TO_NUMBER(tanggal) … WHERE ROWNUM = 1`.
>
> ⭐ **Satu baris berlaku global** — menjawab pertanyaan bentuk tabel pada **P54**. Bukan per periode,
> bukan per lini bisnis.
>
> Aturan geser bulan: bila tanggal hari ini **lebih besar** dari tanggal closing, periode digeser
> satu bulan ke depan (`ADD_MONTHS`). Sesuai **P54**.
>
> Deret nomor disimpan di `POOLDATA.GENERATE_SEQUENCE_NUMBER` — `class`, `jenis`, `tahun`, `no_seq`,
> `tanggal`, `mm_yyyy`. Deretnya **per (class, jenis, tahun)**, dikunci `FOR UPDATE`, dan bila
> barisnya belum ada disisipkan dengan `no_seq = 1`. Keluarannya berbantalan nol lima digit.
>
> **Tidak ada `COMMIT` di dalam** — sesuai **P2**. Pemanggilnya yang meng-commit.
>
> `[terbuka]` ⛔ **Ada tanggal yang ditulis mati di dalam procedure:**
> `IF TRUNC(v_now) <= TO_DATE('02/01/2026','DD/MM/YYYY') THEN v_mm_yyyy := '12.2025'`.
> Seluruh penomoran sebelum 2 Januari 2026 dipaksa ke periode Desember 2025. Tampak sebagai
> penanganan perpindahan tahun yang dibuat sekali jalan. **Perlu dipastikan apakah masih
> dikehendaki**, dan apakah pola yang sama akan diperlukan pada perpindahan tahun berikutnya.
>
> ---
>
> ## 4 · `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` — penghapusannya BERSYARAT
>
> ⭐ **Tidak menghapus apa pun bila konversi sudah berhasil.** Penjaganya:
> `STS_KONVERSI` harus `NULL` atau bukan `'1'`. Ini memperbaiki pemahaman pada **P51** — penghapusan
> tidak dapat menyentuh data yang sudah sah.
>
> Lingkupnya ditentukan parameter `Bisnis`, dan **selalu berdasarkan `IDPEGA`** — satu kasus, bukan
> lebih luas:
>
> | `Bisnis` | Tabel yang dihapus |
> | --- | --- |
> | `'T'` — treaty | `JSON_POLIS` · `TREATYINPRODUCTION` · `TREATYINPRODUCTION_BACKUP` |
> | `'F'` — fakultatif | `JSON_POLIS` · `FACINPRODUCTION` · `FACINPRODUCTION_BACKUP` · `FACOUTPRODUCTION` |
>
> `COMMIT` **di dalam**, sesudah seluruh penghapusan berhasil. Galat pada tabel mana pun
> menghasilkan `ROLLBACK` dan penghentian — jadi penghapusannya **utuh atau tidak sama sekali**.
>
> ⚠️ **Tabel `_BACKUP` ikut dihapus.** Namanya menyiratkan jaring pengaman, tetapi ia dibersihkan
> bersama data utamanya. **Tidak ada pemulihan** sesudah penghapusan ini berjalan. Perlu diketahui
> sebelum jalur ini dibangun ulang.
>
> ---
>
> ## Yang ditegaskan atau diperbaiki oleh naskah ini
>
> | Butir | Akibat |
> | --- | --- |
> | **P2** | ditegaskan tepat: dua meng-commit sendiri, penomoran tidak, penghapus meng-commit |
> | **P51** | diperbaiki: penghapusan **bersyarat**, tidak menyentuh konversi yang sudah berhasil |
> | **P54** | ditegaskan: `TANGGAL_CLOSING` satu baris global, aturan geser bulan sesuai |
> | **P29** | lingkup mengecil: data kontrak sudah relasional; hanya `DATA_JSON` yang perlu dipecah |
>
> ---
>
> ## TAMBAHAN 2026-09-22 — dua procedure lagi diterima
>
> Ronde keadaan EDM menemukan dua procedure yang **tidak pernah diminta** pada P1. Naskah keduanya
> diberikan work owner pada hari yang sama. **Jumlah procedure yang dikenal kini enam, bukan empat.**
>
> ### 5 · `POOLDATA.INSERTJSONPOLISMONITORING`
>
> Satu tabel: `JSON_POLIS_MONITORING` — `IDPEGA`, `DATA_JSON`, `TGL_INPUT`, `NOPOLIS`, `PRODKE`,
> `OLDNOPOLIS`, `TGL_PROD`.
>
> ⭐ **Idempoten, tetapi bukan "insert-or-update".** Bila sudah ada baris ber-`IDPEGA` yang sama, ia
> **berhenti tanpa berbuat apa-apa** — cabang bertanda *"Update"* isinya hanya `RETURN`. Namanya
> menjanjikan pembaruan; perilakunya tidak.
>
> `PRODKE` **selalu diisi `0`**, bukan dari parameter.
>
> `COMMIT` **di dalam**.
>
> `[terbuka]` ⛔ **Pesan galatnya memuat markah HTML** — `<span style="color:red">…</p>`. Lapisan
> basis data menghasilkan tampilan, dan tag-nya pun tidak berpasangan. Di sistem baru pesan galat
> tidak boleh membawa markah; penyajian milik lapisan tampilan.
>
> `[terbuka]` ⛔ **Parameter `p_ProductionLevel` mengubah bunyi pesan** — bernilai `'5'`
> menyembunyikan `SQLERRM` dan menggantinya dengan *"Harap Hubungin IT Terkait"*. Penanda lingkungan
> merembes ke dalam procedure. Menyentuh **ADR-0005**; perlu ditinjau ketika jalur ini ditulis ulang.
>
> `[terverifikasi]` Kolom `STS_KONVERSI` **tidak diisi** procedure ini. Ia ditetapkan terpisah oleh
> langkah 16 `serviceInsertArasapas_act` lewat `UPDATE` — sejalan temuan **P51** dan **P52**.
>
> Blok besar yang dikomentari — `c_counter_prodke`, `PRODKE_SEQ`, `new_uuid` — **tidak dimigrasi**.
>
> ### 6 · `POOLDATA.InsertUpdateAchievment`
>
> Satu tabel: `POOLDATA.ACHIEVEMENT`, 18 kolom — `IDPEGA` `NOPOLIS` `NOOFFER` `SOBNAME`
> `TREATYGROUPNAME` `TREATYTYPE` `QUARTER` `QUARTERYEAR` `IDCURRENCY` `CURRENCY` `PREMIUM` `RICOMM`
> `BROKERAGE` `NETPREMIUM` `PAIDCLAIM` `OUTSTANDINGCLAIM` `PROPORTIONALTYPE` `TGL_PROD`.
>
> ⭐ **Namanya menyebut Update, tetapi ia hanya menyisip.** Tidak ada jalur pembaruan sama sekali.
>
> **Tidak ada `COMMIT` di dalam** — pemanggilnya yang meng-commit. Sejalan **P2**.
>
> ⭐⭐ **Dan inilah yang paling penting: parameter uangnya bertipe `NUMBER`, bukan `VARCHAR2`.**
> `P_PREMIUM` `P_RICOMM` `P_BROKERAGE` `P_NETPREMIUM` `P_PAIDCLAIM` `P_OUTSTANDINGCLAIM` seluruhnya
> `NUMBER`. Ini satu-satunya procedure yang diketahui menerima uang sebagai bilangan; lima lainnya
> menerimanya sebagai teks atau di dalam dokumen.
>
> ⇒ Di sinilah presisi ditetapkan saat melintasi batas. Bila Pega mengirim bilangan mengambang,
> galatnya dibekukan di kolom `NUMBER` pada titik ini.
>
> `[terbuka]` ⛔⛔ **Penjaga anti-duplikat dipatahkan oleh galat uang.** Pemeriksaannya
> membandingkan **seluruh 17 kolom**, termasuk keenam kolom uang:
>
> ```
> WHERE … AND PREMIUM = P_PREMIUM AND RICOMM = P_RICOMM AND BROKERAGE = P_BROKERAGE
>         AND NETPREMIUM = P_NETPREMIUM AND PAIDCLAIM = P_PAIDCLAIM
>         AND OUTSTANDINGCLAIM = P_OUTSTANDINGCLAIM
> ```
>
> Selisih sekecil `2,76 x 10^-7` — galat yang **terbukti ada di data produksi**, lihat butir terbuka
> galat uang pada **P29** — membuat baris yang seharusnya sama dinilai **berbeda**, sehingga baris
> kedua **tersisip**. Penjaganya gagal justru pada kasus yang paling perlu dijaga.
>
> ⇒ Sistem baru **tidak boleh** meniru pemeriksaan duplikat berbasis kesamaan seluruh kolom uang.
> Kuncinya harus medan pengenal — bukan nilai uang. **Perlu keputusan work owner**, dan berkaitan
> langsung dengan pilihan (a)/(b)/(c) pada butir galat uang.
>
> ### Ringkas enam procedure
>
> | # | Procedure | Tabel yang disentuh | `COMMIT` di dalam |
> | ---: | --- | --- | :---: |
> | 1 | `PEGA_TREATY_IN` | `M_TREATY_IN` · `TREATY_IN` | ya |
> | 2 | `PEGA_JSON_POLIS_TREATYIN` | `json_polis` | ya |
> | 3 | `PROC_GENERATE_SEQUENCE_NUMBER` | `TANGGAL_CLOSING` *(baca)* · `GENERATE_SEQUENCE_NUMBER` | **tidak** |
> | 4 | `PEGA_DELETE_ERROR_KONVERSI` | `JSON_POLIS` · `TREATYINPRODUCTION` · `…_BACKUP` · *(sisi Fac)* | ya |
> | 5 | `INSERTJSONPOLISMONITORING` | `JSON_POLIS_MONITORING` | ya |
> | 6 | `InsertUpdateAchievment` | `ACHIEVEMENT` | **tidak** |
>
> `[terbuka]` Belum dipastikan apakah **enam sudah seluruhnya**. Korpus mencatat lebih dari sepuluh
> nama `POOLDATA.*` yang dipanggil; sisanya milik modul lain dan belum ditelusuri.

---

## 2. P29 — Sebuah kolom yang namanya berarti "jenis usaha" ternyata menyimpan seluruh data kontrak. Benarkah?

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

**Jawaban:**

> **Terjawab — dan sumber datanya dipindahkan.** `[keputusan work owner]` `[data DBA]` 2026-09-22
>
> **JSON tidak dipakai di sistem baru.** Pembacaan `JSONDATA` beserta seluruh mekanisme
> `adoptJSONObject` + `getString("CLASSOFBUSINESS")` **tidak dimigrasi**. Sebagai gantinya, data
> ditarik dari **view relasional** yang sudah ada di Oracle:
>
> ```
> POOLDATA.TREATYINDETAILJOINEDM
>    = pooldata.treatyindetail  UNION ALL  pooldata.treatyindetailedm
> ```
>
> Naskah view diberikan work owner pada 2026-09-22. **39 kolom**, di antaranya:
>
> | Kelompok | Kolom |
> | --- | --- |
> | pengenal | `ID` `TREATYID` `TREATYYEAR` `SOBID` `SOB` `CEDINGID` `CEDING` |
> | kontrak | `TREATYCONTRACTNAME` `PROPORTIONTYPE` `TREATYTYPE` `TREATYGROUPID` `TREATYGROUP` |
> | penggolongan | `CLASSOFBUSINESSID` `CLASSOFBUSINESS` |
> | uang | `LIMITCURRENCY`/`LIMITVALUE` · `RETENTIONCURRENCY`/`RETENTIONVALUE` · `EPICURRENCY`/`EPIVALUE` · `NETPREMICURRENCY`/`NETPREMIVALUE` · `MDPCURRENCY`/`MDPVALUE` · `SHARECURRENCY`/`SHAREVALUE` |
> | potongan | `DEDUCTION1` `DEDUCTION2` `BROKERAGE` |
> | lapisan | `LAYERTYPE` `LAYER` `LAYERPARTTYPE` `LAYERPART` `INSTALLMENTNO` |
> | lain | `RIOGR` `RIONR` `RNM_SHARE` `COMMENCEMENT` `TERMINATION` |
>
> **Uji kecukupan — semua terpenuhi.** `[terverifikasi]`
> `ReportDefinition\BrowseTreatyInDetail.xml` merujuk **33 medan**; ke-**33**-nya ada di view.
> **Nol medan yang dipakai tetapi tidak tersedia.** Dari 31 berkas layar, **29 kolom view** dipakai
> langsung.
>
> Dengan demikian pertanyaan asal tertutup: nama `CLASSOFBUSINESS` dipakai bermacam-macam di sistem
> lama — sebagai alias untuk `JSONDATA`, untuk `BIZCODE`, untuk `BIZNAME`, sebagai medan di dalam
> dokumen JSON, dan sebagai kolom asli. Di sistem baru ia **satu hal saja**: kolom
> `CLASSOFBUSINESS` pada view di atas.
>
> **Enam aktivitas yang membongkar JSON tidak dimigrasi:** `FetchMasterTreatyIn` ·
> `SetTreatyIn_Act` · `InputPolicyTreatyInDetail_preACT` · `InputPolicyTreatyInDetail_NonProp` ·
> `InputPolicyTreatyOutDetail_preACT` · `InputPolicyTreatyOutDetail_NonProp`. Keduanya yang pertama
> memetik `CLASSOFBUSINESS` sepuluh kali; sisanya sekali.
>
> `[terbuka]` **`RNM_SHARE` tidak dipakai sama sekali** — nol kemunculan di seluruh modul
> NB Treaty In. Perlu dipastikan apakah kolom itu dipakai modul lain sebelum diabaikan.
>
> **Bentuk data di sistem baru — ditetapkan 2026-09-22.** `[keputusan work owner]`
>
> | Hal | Ketetapan |
> | --- | --- |
> | presisi kolom uang | **tidak dibatasi.** Nilai disimpan dengan presisi penuh sebagaimana Oracle mengembalikannya; pembulatan hanya boleh di titik penyajian, tidak pernah di repository |
> | `COMMENCEMENT` · `TERMINATION` | **dibaca apa adanya.** Bila kolomnya `DATE`, dibaca sebagai tanggal; bila teks, teksnya tidak diformat ulang. Berlaku hanya untuk pembacaan data lama — keputusan **P32** tetap berdiri untuk penulisan tanggal baru |
> | `DEDUCTION1` · `DEDUCTION2` · `BROKERAGE` · `RNM_SHARE` | **persentase, bukan nilai uang.** Disimpan apa adanya: nilai `12.5` berarti 12,5 persen, bukan 0,125. Tipe terpisah dari uang, tanpa mata uang — view memang tidak punya kolom mata uang untuk keempatnya |
>
> Konsekuensi yang mengikat implementasi: nilai uang **tidak boleh** dibaca sebagai `float`
> (CLAUDE.md §7), dan pengali persentase **tidak boleh** dipakai langsung terhadap nilai uang —
> perkalian tanpa pembagian seratus menghasilkan angka seratus kali lipat tanpa menimbulkan galat.
>
> `[terbuka]` **`RNM_SHARE` persentase dari apa?** Ketiga persentase lain jelas dihitung terhadap
> premi, tetapi `RNM_SHARE` nol dipakai di seluruh modul NB Treaty In sehingga tidak ada satu pun
> rumus yang menunjukkan penyebutnya. **Tidak memblokir** — kolomnya dibaca dan disimpan apa adanya.
> Wajib dijelaskan sebelum kolom itu dipakai dalam perhitungan mana pun.
>
> ⛔ **RALAT 2026-09-22, sore.** `[penyimpangan sadar]` Jawaban ini semula ditutup dengan:
>
> > *"Permintaan contoh isi `JSONDATA` dibatalkan untuk keperluan bentuk data, sebab bentuknya kini
> > diambil dari view. Contoh masih berguna hanya bila kelak diperlukan pemindahan data lama."*
>
> **Keliru, dan pembatalannya dicabut.** Pernyataan itu benar untuk **sisi baca** — view
> `TREATYINDETAILJOINEDM` memang mencukupi di sana. Tetapi **sisi tulis adalah hal yang berbeda**:
> `pooldata.treatyindetail` dan `treatyindetailedm` bukan tabel penyimpanan modul ini, melainkan
> data kontrak yang dibaca. Yang selama ini **menyimpan** data NB Treaty In adalah JSON polis,
> lewat `PEGA_TREATY_IN` dan `PEGA_JSON_POLIS_TREATYIN`.
>
> Karena JSON tidak dipakai lagi, **tabel datar baru harus dirancang** untuk menampung apa yang
> selama ini masuk ke JSON itu.
>
> `[terbuka]` **Contoh isi `JSONDATA` adalah bahan utama perancangan tabel itu, bukan pelengkap.**
> Sebabnya sudah tercatat di **P15**: `GetPageJSONString()` menyalin **seluruh halaman kerja apa
> adanya**, tanpa memilih medan. Tidak ada satu pun aturan di korpus yang menyebutkan medan apa saja
> ada di dalam dokumen itu — hanya datanya yang dapat menunjukkannya.
>
> Preseden proyek: `premiumlist-life` merancang tujuh tabel barunya dengan keterangan
> *"terbukti dari `DATA_JSON` nyata"*, dan `claim-life` menempuh jalan yang sama. Keduanya
> mendahulukan rancangan penyimpanan **sebelum** spec utama, bukan sesudahnya.
>
> **Yang diminta, dinaikkan bobotnya:** 3-5 isi kolom JSON polis apa adanya, beserta **ukuran
> terbesarnya**. Tanpa itu, tabel datar penyimpanan NB Treaty In tidak dapat dirancang — hanya
> ditebak.
>
> ---
>
> **DITERIMA 2026-09-22 sore — satu contoh `DATA_JSON` sungguhan dari work owner.**
> `[data DBA]` `[terverifikasi]`
>
> ⛔ **Contohnya TIDAK disimpan di dalam berkas proyek mana pun.** Ia memuat nama orang pada medan
> pemasar, operator, dan tertanggung. Yang dicatat di bawah hanya **bentuk** dan **nama medan**.
>
> **Bersarang lima tingkat:**
>
> ```
> PolicyTreatyIn
>   +- LocationList[]
>   |    +- OccupationList[]
>   |         +- AnekaList[]
>   |              +- CoverageList[]
>   |                   +- ClauseList[]
>   |                   +- DeductibleList[]
>   +- SuggestList[]
>   +- QuotationData{}        1:1, bukan daftar
> ```
>
> ⇒ perkiraan **delapan tabel**, bukan dua atau tiga. Pohon objek pertanggungan sendiri lima tingkat.
>
> **Medan skalar tingkat atas** (38): `BalanceDueTo` `IDNewBisnis` `BizCode` `BizName` `CedingCo`
> `CedingCoName` `Claim` `Currency` `Deduction1` `Deduction2` `DueTo` `EndDate` `ExcessLoss`
> `IDCurrency` `IsApproved` `InsuredID` `InsuredName` `MarketingOfficer` `NetPremium` `NoOffer`
> `OveriddingCommOgp` `OveriddingCommOnp` `PremiOgp` `PremiOnp` `ResultOgp1` `ResultOgp2`
> `ResultOnp1` `ResultOnp2` `RiCommOgp` `RiCommOnp` `Show` `SOB` `SOBName` `StartDate`
> `StatementDate` `TreatyGroupName` `TreatyYear` — ditambah `pxObjClass`.
>
> **Lima sifat yang mengikat perancangan:**
>
> 1. `[terverifikasi]` **Seluruh nilai bertipe teks**, termasuk uang (`"27545312.5"`), tanggal, dan
>    penanda. Menguatkan ketetapan P29 nomor 3: penguraian di Go **eksplisit**, tipe dari sumber
>    tidak dipercaya.
> 2. `[terverifikasi]` **Dua format tanggal dalam satu dokumen** — `"20171130"` dan
>    `"20170930T170000.000 GMT"`. Persis persoalan **P32**, kini terbukti dari data.
> 3. `[terverifikasi]` **`IsApproved` bernilai `"1"` sebagai teks** — menguatkan **P24**.
> 4. `[terverifikasi]` **`pxObjClass` ada di setiap simpul** — internal Pega, **tidak dimigrasi**.
>    Menegaskan **P15**: dokumen ini potret halaman kerja, bukan skema yang dirancang.
> 5. `[terverifikasi]` **`CedingCo` muncul dua kali dengan nilai berbeda**, dan yang di dalam
>    `QuotationData` **berakhiran `"; "`** — daftar yang digabung dengan pemisah titik-koma.
>    Ceding dapat lebih dari satu, disimpan sebagai teks gabungan, bukan baris terpisah.
>
> **`SuggestList` membawa `IsApproved` per baris**, terpisah dari `IsApproved` tingkat atas —
> menegaskan dua properti senama yang berbeda, sebagaimana tercatat di **P24** dan **P42**.
>
> ⭐ **Uji silang yang paling menguatkan.** Contoh ini ber-`GroupPanel` = `006` dan
> `QuotationData.BusinessOldId` = `01`. Disimulasikan pada `DecisionTable\BusinessType_DeT.xml`
> yang direkonstruksi di **P19**: baris 29 tidak cocok, baris 30 tidak cocok, **baris 31 cocok**
> lewat daftar OR-nya, menghasilkan `"FireStyle2"` — dan dokumennya memang berisi `"FireStyle2"`.
>
> ⇒ Rekonstruksi 36 baris itu **terbukti benar dari data produksi**, bukan hanya dari pembacaan XML.
> Ini juga menutup keraguan ronde 4 yang menyebut tabel itu tinggal 9 baris.
>
> ---
>
> ## ⭐⭐ TIGA CONTOH DITERIMA 2026-09-22 — dan ketiganya mengubah rancangannya
>
> `[data DBA]` `[terverifikasi]`
> ⛔ **Ketiganya TIDAK disimpan di berkas mana pun.** Semuanya memuat nama orang pada medan
> pemasar, operator, dan tertanggung. Yang dicatat hanya **bentuk**, **nama medan**, dan **jumlah**.
>
> ### ⛔⛔ Dokumen JSON ini TIDAK berbentuk tetap — dan bentuknya mengikuti JENIS TREATY
>
> | | Contoh 1 | Contoh 2 | Contoh 3 |
> | --- | ---: | ---: | ---: |
> | jenis | proporsional | proporsional *(SOA)* | ⛔ **non-proporsional XOL** |
> | medan skalar tingkat atas | **37** | **64** | **47** |
> | ⭐ ada di **ketiganya** | | **23** | |
> | ⭐ gabungan ketiganya | | **74** | |
>
> ⚠️ *(Daftar 38 medan yang tercatat di atas untuk contoh 1 termasuk `pxObjClass`; tanpa properti
> internal Pega itu **37**, sesuai tabel ini.)*
>
> ⛔ **Hanya 23 dari 74 medan muncul di ketiga contoh.** Lebih dari dua pertiga gabungan itu
> **tidak ada di sebagian kasus**.
>
> **Daftar bersarangnya berbeda-beda:**
>
> ```
> 1 : LocationList > OccupationList > AnekaList > CoverageList > (ClauseList, DeductibleList)
> 2 : ListInstallment (datar) . SpreadingRiskList
> 3 : ListInstallment > InstallmentList . TreatyXOLList > ValueList . OldData > TreatyXOLList
>     BreakDownSpreadList
> ```
>
> ⛔ **Contoh 2 tidak punya `LocationList` sama sekali** — padahal seluruh pohon objek
> pertanggungan lima tingkat yang dicatat di atas bergantung padanya.
>
> ### ⛔⛔ Bahaya terbesarnya: satu nama daftar, DUA kedalaman
>
> **`ListInstallment` bersarang di contoh 3 tetapi datar di contoh 2.** Nama sama, kedalaman
> berbeda. ⛔ **Pengurai yang menganggapnya seragam akan patah.**
>
> ⭐⭐ **Dan ini dapat diperiksa dari korpus, bukan hanya dari contoh.** `[terverifikasi]`
> Disapu atas NB + EDM, blok penyunting dibuang:
>
> | Bentuk yang dirujuk aturan | Rujukan | Berkas |
> | --- | ---: | ---: |
> | ⛔ `ListInstallment(n).InstallmentList` — **bersarang** | **58** | 6 |
> | ⛔ `ListInstallment(n).<medan skalar>` — **datar** | **101** | 11 |
> | `TreatyXOLList(n).ValueList` | 207 | 8 |
> | `OldData.TreatyXOLList` | 130 | 9 |
> | `BreakDownSpreadList` | 11 | 5 |
>
> ⇒ **Kedua bentuk `ListInstallment` sama-sama nyata di dalam aturan sistem lama.** Contoh 2 dan
> contoh 3 bukan kebetulan — keduanya sah.
>
> ### ⭐ Penentu bentuknya terlihat
>
> `QuotationData.ProportionalType` bernilai `"Proportional"` lawan `"NonProportional"`, dan
> `IsNewPolicyNonProp` bernilai `"0"` lawan `"1"`. **Bentuk dokumen mengikuti jenis treaty.**
> `[terverifikasi]` Korpus menyebut `ProportionalType` **170 kali** dan `IsNewPolicyNonProp`
> **87 kali**.
>
> ⚠️ **Akibatnya bagi rancangan:** tabel penyimpanan tidak dapat dirancang sebagai satu bentuk
> tunggal tanpa lebih dulu memutuskan bagaimana kedua jenis treaty itu ditampung — satu tabel
> dengan kolom opsional, atau pemisahan menurut jenis. ⛔ **Itu keputusan rancangan, dan ronde ini
> tidak merancang apa pun.**
>
> ### ⭐ `OldData` muncul di dalam berkas POLIS BARU
>
> Contoh 3 memuat **`OldData`** — halaman data lama yang selama ini dikenal hanya dari EDM —
> **di dalam dokumen polis baru**, berisi `TreatyXOLList` kosong. `[terverifikasi]` Korpus
> mendukungnya: `OldData.TreatyXOLList` dirujuk **130 kali di 9 berkas**.
>
> ⇒ Halaman pembanding endorsemen **sudah ikut tersimpan sejak polis dibuat**, walau kosong.
>
> ### ⭐ Dan korpus mengenal lebih banyak lagi daripada kedua contoh
>
> `[terverifikasi]` Diukur **dua cara** atas `NB Treaty In` + `EDM Treaty In`, blok penyunting
> dibuang lebih dulu:
>
> | Cara | Hasil |
> | --- | ---: |
> | **A** — sapuan rujukan teks `PolicyTreatyIn.<nama>` di seluruh berkas | **94** nama unik, 1.869 rujukan |
> | **B** — hanya sisi-kiri penetapan nilai *(`PropertiesName` + `pyPropertiesName`)* | **86** nama unik, 867 penulisan |
> | irisan · hanya di A · hanya di B | 86 · 8 · **0** |
>
> ⇒ **94 properti**, di antaranya ⭐ **9 simpul bersarang**:
>
> | Simpul bersarang | Rujukan | Muncul di kedua contoh? |
> | --- | ---: | --- |
> | `OldData` | 210 | ⛔ tidak |
> | `TreatyXOLList` | 169 | ⛔ tidak |
> | `ListInstallment` | 101 | hanya contoh 2 |
> | `TreatyXOLDifferenceList` | 92 | ⛔ tidak |
> | `TreatyDifference` | 89 | ⛔ tidak |
> | `SpreadingRiskList` | 65 | hanya contoh 2 |
> | `QuotationData` | 37 | hanya contoh 2 |
> | `PolicyTreatyInDetail` | 2 | ⛔ tidak |
> | `PolicyNo` | 3 | ⚠️ **diduga artefak pola**, bukan simpul sungguhan |
>
> ⚠️ **Ini BUKAN sensus.** Ia hasil satu bentuk penelusuran, dijalankan dua cara — **batas bawah**,
> bukan angka akhir. Properti yang hanya muncul di dalam dokumen JSON dan **tidak pernah dirujuk**
> satu aturan pun **tidak akan tertangkap cara mana pun di atas**.
>
> ### ⛔⛔ `[terbuka]` Akibatnya bagi perancangan tabel
>
> **Rancangan tabel TIDAK dapat dibangun dari contoh.** Dua contoh saja sudah berselisih **30 medan**
> dan **tidak berbagi satu pun daftar bersarang**; korpus menunjukkan sedikitnya **94 properti** dan
> **9 simpul bersarang**, empat di antaranya — `OldData`, `TreatyXOLList`,
> `TreatyXOLDifferenceList`, `TreatyDifference` — **tidak muncul di kedua contoh** padahal dirujuk
> ratusan kali.
>
> ⭐ **Urutan yang benar:** rancangan dibangun dari **sensus properti korpus**, dengan contoh
> `DATA_JSON` sebagai **penguji** — bukan sebaliknya.
>
> ⛔ **Sensus itu ronde tersendiri, dan belum dikerjakan.** Ia yang kini ditunggu tiket
> `00-skema-penyimpanan-dan-migrasi`, menggantikan P1 dan P29.
>
> ---
>
> ## ⛔⛔ `[terbuka]` GALAT ANGKA UANG YANG SUDAH TERSIMPAN PERMANEN DI PRODUKSI
>
> `[terverifikasi]` `[data DBA]` 2026-09-22 — dari contoh 3.
> ⛔⛔ **Butir ini menyentuh ketetapan P29 nomor 1** *(presisi penuh, pembulatan hanya di titik
> penyajian)*, dan **tidak boleh diputuskan tim migrasi.**
>
> ```
> premium angsuran   148157378.220000069     x 4  =  592629512.880000276
> NetPremium         592629512.880000276          <- sama persis
> nilai bulat 2 desimal                           =  592629512.88
> selisih                                         =  2,76 x 10^-7
> ```
>
> ⭐ **Hitungannya saya ulang sendiri dengan `Decimal` presisi 40**, bukan disalin: perkaliannya
> **cocok persis**, dan selisihnya **tepat 2,76 x 10^-7**.
>
> Ekor `...0000276` **bukan presisi** — ia galat yang **berlipat empat** karena premi terpecah ke
> empat angsuran, tiap angsuran membawa galatnya, lalu dijumlahkan kembali.
>
> ### ⚠️ Satu koreksi atas dugaan sebabnya — dan ia memperburuk, bukan meringankan
>
> Galat ini **bukan** sekadar galat bilangan mengambang ganda *(IEEE-754 `float64`)*. Diperiksa:
>
> | | Nilai |
> | --- | --- |
> | `float64` terdekat dari `592629512.88` | `592629512.87999999523...` — meleset **ke bawah**, sekitar **4,8 x 10^-9** |
> | galat yang benar-benar tersimpan | meleset **ke ATAS**, **2,76 x 10^-7** |
> | ⭐ nisbahnya | sekitar **57 kali lebih besar**, dan **arahnya berlawanan** |
> | `592629512.88 / 4` dalam desimal | `148157378.22` — **tepat, tanpa sisa** |
>
> ⇒ ⛔ **Pembagian empat yang bersih TIDAK akan menghasilkan ekor ini sama sekali.** Galatnya
> lahir di **rantai perhitungan**, bukan di penyimpanan — kemungkinan dari pembagian persentase
> berpresisi terbatas *(`@divide(...,n)` yang membulatkan di `n` desimal)*, lalu dikalikan kembali.
> ⛔ **Sebab pastinya tidak dapat dinyatakan dari satu kasus, dan tidak dikarang di sini.**
>
> ⭐ **Kenapa ini penting bagi keputusan di bawah:** bila galatnya ada di rantai hitung, maka
> sistem baru yang memakai aritmetika desimal yang benar **akan menghasilkan angka yang BERBEDA
> dari sistem lama** — bukan hanya pada data lama, tetapi pada perhitungan baru. Itu menyentuh
> **uji paritas** di tiket **15**.
>
> ### ⛔ Tiga pilihan, dan akibat masing-masing — **tidak dipilih di sini**
>
> | | Pilihan | Akibat |
> | --- | --- | --- |
> | **a** | **ikuti apa adanya** | paritas sempurna dengan sistem lama; ⛔ cacatnya **diwariskan**, dan tiap penjumlahan menambahnya |
> | **b** | **bulatkan saat migrasi** | bersih ke depan; ⛔ **angka historis berubah**, laporan lama tidak lagi cocok |
> | **c** | **simpan apa adanya, bulatkan saat dihitung** | sejarah utuh, galat **berhenti berlipat**; ⚠️ angka hasil hitung berbeda dari sistem lama |
>
> ⛔⛔ **Tim migrasi TIDAK memilih.** Pemutusnya `[work owner]`, dan **presisi uang yang wajar per
> mata uang** adalah pertanyaan `[Finance]`.
>
> ⭐ **ADR-0003 tetap ditegakkan apa pun pilihannya:** sistem baru **tidak memakai `float`**,
> sehingga ia **tidak menambah galat baru** — apa pun yang diputuskan tentang galat yang sudah ada.
> ⛔ Tidak ada ADR baru yang diperlukan untuk butir ini.
>
> ⚠️ **Selama butir ini terbuka, ketetapan P29 nomor 1 tetap berlaku apa adanya** — nilai dibaca
> dan disimpan dengan presisi penuh. Butir ini **tidak menahan** penulisan lapisan penyimpanan; ia
> menahan **keputusan migrasi data lama** dan **ambang uji paritas**.

---

## 3. P2 — Apakah ketiga program itu menyelesaikan penyimpanannya sendiri, atau menunggu aplikasi?

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

**Dampak bila salah:** ⛔ separuh data menjadi permanen sementara separuh lainnya batal — dan
selisihnya hanya ketahuan berbulan-bulan kemudian.

rujukan: `COMMIT` di dalam blok PL/SQL Pega vs di dalam procedure — OQ-013, butir baru #2

**Jawaban:**

> **Terjawab per procedure.** `[keputusan work owner]` `[data DBA]` 2026-09-22
>
> | Procedure | `COMMIT` di dalam procedure |
> | --- | --- |
> | `POOLDATA.PEGA_TREATY_IN` | **ada** |
> | `POOLDATA.PEGA_JSON_POLIS_TREATYIN` | **ada** |
> | `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | **tidak ada** |
>
> **Akibatnya pada sistem lama.** `[terverifikasi]` Blok PL/SQL pemanggilnya juga memuat `COMMIT` —
> lihat `RDBList\SaveTreatyIn.xml`, `RDBList\SavePolisTreatyIn_SQL.xml`,
> `RDBList\GetSequenceNumber_SQL.xml`. Jadi untuk dua yang pertama terjadi **commit ganda**:
> procedure meng-commit, lalu blok meng-commit lagi. Yang kedua tidak berpengaruh.
>
> Konsekuensi nyatanya: untuk `PEGA_TREATY_IN` dan `PEGA_JSON_POLIS_TREATYIN`, **aplikasi tidak
> dapat membatalkan penyimpanan** setelah procedure dipanggil. Kegagalan pada langkah berikutnya
> meninggalkan data yang sudah permanen — sistem lama **tidak punya jaminan keutuhan** untuk urutan
> yang memanggil keduanya.
>
> Untuk `PROC_GENERATE_SEQUENCE_NUMBER`, `COMMIT` pada blok pemanggilnya yang berlaku.
>
> **Untuk sistem baru** `[penyimpangan sadar]`: karena ketiganya ditulis ulang di Go tanpa stored
> procedure, keputusan commit menjadi milik aplikasi sepenuhnya. Seluruh urutan penyimpanan
> dibungkus **satu transaksi**, dan kegagalan di langkah mana pun membatalkan seluruhnya. Ini
> **lebih ketat** daripada sistem lama, dan disengaja.
>
> `[terbuka]` Perlu dipastikan apakah ada proses lain — laporan, integrasi, atau pekerjaan
> terjadwal — yang selama ini **mengandalkan** data menjadi permanen lebih awal. Bila ada, perubahan
> ini mengubah waktu munculnya data bagi proses tersebut.

---

## 4. P3 — Empat nama tabel muncul dalam dua ejaan. Satu tabel atau dua?

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

**Jawaban:**

> **Satu tabel. Seluruhnya milik skema `POOLDATA`.** `[keputusan work owner]` `[data DBA]` 2026-09-22
>
> Keempat nama menunjuk tabel yang sama, baik ditulis dengan awalan maupun tanpa. Penulisan tanpa
> awalan berarti skema bawaan pengguna yang menyambung, dan pengguna itu adalah `POOLDATA`.
>
> | Nama | Pemilik | `POOLDATA.` | tanpa awalan |
> | --- | --- | ---: | ---: |
> | `CURRENCY` | `POOLDATA` | 4 | 24 |
> | `REINSURANCETYPE` | `POOLDATA` | 18 | 10 |
> | `PROPORTIONALARRG` | `POOLDATA` | 2 | 22 |
> | `HISTORYAKSEPTASIPEGA` | `POOLDATA` | 0 | 15 |
>
> Angka dihitung dari 1.151 naskah SQL di 21 modul, setelah 54 rujukan properti yang menyerupai
> nama tabel dibuang. `[terverifikasi]`
>
> **Mengikat implementasi:** setiap query di sistem baru menulis skema **secara eksplisit** —
> `POOLDATA.CURRENCY`, bukan `CURRENCY`. Mengandalkan skema bawaan pengguna membuat perilaku
> berbeda antara lingkungan uji dan produksi, dan perbedaan itu tidak menimbulkan pesan galat.

---

## 5. P4 — Satu kolom pada tabel riwayat tidak pernah diisi dari modul ini. Siapa mengisinya?

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

**Jawaban:**

> **Kolom dipertahankan; diisi dari identitas akses login.** `[keputusan work owner]` 2026-09-22
>
> `POOLDATA.HISTORYAKSEPTASIPEGA.OPERATORID` **tetap dibawa** ke sistem baru, dan diisi dari
> **identitas akses login** pengguna yang melakukan tindakan — bukan dari nama tampilan.
>
> Sejalan dengan pola yang sudah ada di tabel riwayat satunya: `RDBList\InsertViewSuggest_SQL.xml`
> mengisi kolom `AKSES_LOGIN` dari `OperatorID.pyUserIdentifier`, terpisah dari kolom `PIC` yang
> memuat nama tampilan. Lihat **P33**, tempat pemisahan itu ditetapkan.
>
> **Perubahan perilaku yang disadari.** `[penyimpangan sadar]` `[terverifikasi]` Sistem lama
> **tidak pernah** menulis ke kolom ini dari aplikasi — nol dari 1.151 naskah SQL di 21 modul
> menyebutnya sebagai kolom. Ke-67 kemunculan `OPERATORID` di korpus seluruhnya adalah halaman
> Pega `OperatorID.pyXxx`, bukan nama kolom. Jadi sistem baru **mulai mengisi kolom yang selama ini
> kosong** — itu perbaikan jejak audit, bukan peniruan.
>
> `[terbuka]` **Belum dicacah** berapa baris yang `OPERATORID`-nya sudah terisi sekarang.
> **Tidak memblokir.** Bila ternyata tidak nol, artinya ada penulis di luar aplikasi — pemicu,
> penjadwal, atau sistem lain — yang belum masuk peta migrasi, dan itu wajib ditelusuri sebelum
> go-live. Satu perintah cukup:
> `SELECT COUNT(*) total, COUNT(operatorid) terisi FROM POOLDATA.HISTORYAKSEPTASIPEGA;`

---

## 6. P49 — Satu kolom riwayat tidak pernah diisi dari sistem lama. Apakah ada yang lain yang mengisinya?  ⭐ *(BARU — ronde 4)*


Kami sudah menanyakan kolom penampung identitas operator pada tabel riwayat akseptasi **(P4)**.
⭐ **Sekarang kami dapat menjawab separuhnya sendiri, dan jawabannya tegas:**

⛔ **Dari seluruh 1.151 pernyataan SQL di seluruh sistem lama — bukan hanya modul ini — tidak satu
pun menulis ke kolom itu.** Tidak satu pun bahkan menyebutnya sebagai kolom. Ke-67 kemunculan
namanya ternyata **nama halaman pengguna**, bukan nama kolom.

**Konteks:** yang kami periksa hanyalah pernyataan SQL yang ditulis di dalam sistem lama.
⛔ Kami **tidak** dapat memeriksa program yang tersimpan di dalam basis data, pemicu basis data,
pekerjaan terjadwal, atau sistem lain di luar Pega.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** kolom itu memang tidak pernah terisi,
boleh ditinggalkan · **(b)** diisi dari luar Pega, sebutkan oleh apa · **(c)** dulu terisi, sekarang
tidak.
⭐ Yang paling menentukan, dan cepat: ⭐ **berapa baris yang kolom itu tidak kosong** di basis data
sekarang. Nol berarti (a); lebih dari nol berarti (b) atau (c).

**Dampak bila salah:** ⚠️ sebuah keterangan jejak audit yang selama ini ada **hilang** dari sistem
baru, atau sebaliknya kami membawa kolom mati selamanya.

rujukan: sisiran 1.151 `pyBrowseSQL` di 21 modul; 67 kemunculan teks `OPERATORID`, **seluruhnya**
berbentuk `OperatorID.pyXxx`; ⛔ 0 penulisan, 0 rujukan kolom. ⚠️ Lingkup **tidak** mencakup stored
procedure, pemicu, penjadwal, atau sistem luar — ronde 4 Bab D

**Jawaban:**

> **DITUTUP SEBAGAI RANGKAP — lihat P4.** `[keputusan work owner]` 2026-09-22
>
> Butir ini menanyakan hal yang sama dengan **P4**, dan ditulis pada ronde 4 tanpa mengetahui P4
> sudah ada. Jawaban yang berlaku ada di P4: kolom `OPERATORID` **dipertahankan** dan diisi dari
> **identitas akses login**.
>
> Temuan ronde 4 yang menyertainya tetap sah dan sudah pindah ke P4: dari 1.151 pernyataan SQL di
> 21 modul, **nol** menulis ke kolom itu, dan ke-67 kemunculan namanya seluruhnya berbentuk
> `OperatorID.pyXxx` — nama halaman Pega, bukan nama kolom. `[terverifikasi]`
>
> `[terbuka]` **Satu hitungan masih layak diminta, tetapi tidak memblokir**, dan ditumpangkan ke
> surat **P1** sebagai satu baris tambahan, bukan sebagai pertanyaan tersendiri:
>
> ```sql
> SELECT COUNT(*) total, COUNT(operatorid) terisi FROM POOLDATA.HISTORYAKSEPTASIPEGA;
> ```
>
> Nol berarti tidak ada tindak lanjut. Bukan nol berarti ada penulis **di luar Pega** — pemicu basis
> data, penjadwal, atau sistem lain — yang belum masuk peta migrasi dan wajib ditelusuri sebelum
> go-live.

---

## 7. P61 — Tabel bernama `_BACKUP` ikut terhapus bersama data utamanya. Sengaja?  ⭐ *(BARU — dari naskah procedure)*

Naskah `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` yang Anda kirim menghapus, untuk satu kasus yang
gagal dikonversi, **tabel utama beserta tabel bernama `_BACKUP`-nya sekaligus**:

| `Bisnis` | Tabel yang dihapus |
| --- | --- |
| `'T'` — treaty | `JSON_POLIS` · `TREATYINPRODUCTION` · **`TREATYINPRODUCTION_BACKUP`** |
| `'F'` — fakultatif | `JSON_POLIS` · `FACINPRODUCTION` · **`FACINPRODUCTION_BACKUP`** · `FACOUTPRODUCTION` |

**Konteks:** namanya menyiratkan **jaring pengaman** — salinan yang disimpan supaya data dapat
dipulihkan. Tetapi ia dibersihkan bersama data utamanya, sehingga **sesudah penghapusan berjalan
tidak ada apa pun yang tersisa untuk dipulihkan**. `[terverifikasi]`

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** sengaja, `_BACKUP` bukan
cadangan pemulihan melainkan tabel kerja, sebutkan fungsinya · **(b)** sengaja, dan memang tidak
ada kebutuhan pemulihan · **(c)** tidak sengaja, seharusnya `_BACKUP` dipertahankan.
⭐ Dan: **apakah ada cadangan lain di luar tabel ini** — ekspor berkala, jurnal basis data, atau
salinan di sistem lain.

**Dampak bila salah:** ⛔ sistem baru menyalin perilaku yang **menghapus satu-satunya salinan
cadangan** setiap kali konversi gagal; kesalahan pada jalur konversi menjadi **tidak dapat
dipulihkan**, dan itu baru diketahui saat pemulihan benar-benar dibutuhkan.

rujukan: naskah `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` §jawaban **P1** nomor 4; jalurnya
`Activity\serviceInsertArasapas_act` langkah 12 → `RDBList\DeleteDataProduction` — lihat **P51**

⚠️ Pemutusnya `[work owner]`; dicatat di lembar ini karena faktanya berasal dari basis data.

**Jawaban:**

> *(tulis di sini)*

---

## 8. P62 — Sebuah tanggal ditulis mati di dalam program penomoran. Masih dikehendaki?  ⭐ *(BARU — dari naskah procedure)*

Naskah `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` memuat satu baris yang memaksa periode penomoran:

```
IF TRUNC(v_now) <= TO_DATE('02/01/2026','DD/MM/YYYY') THEN v_mm_yyyy := '12.2025'
```

⇒ **seluruh penomoran sebelum 2 Januari 2026 dipaksa ke periode Desember 2025**, mengabaikan
tanggal closing maupun aturan geser bulan yang berlaku di sisa procedure. `[terverifikasi]`

**Konteks:** bentuknya seperti **penanganan perpindahan tahun yang dibuat sekali jalan** — ditambahkan
untuk satu pergantian tahun tertentu, lalu tertinggal di dalam naskah.

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** aturan tetap, periode
Desember memang berlaku sampai 2 Januari setiap tahun · **(b)** penanganan sekali jalan untuk
pergantian 2025–2026 saja, boleh **tidak** dibawa ke sistem baru · **(c)** lainnya.
⭐ Dan: **apakah pola yang sama akan diperlukan pada pergantian ke 2027**, atau aturannya harus
dibuat umum.

**Dampak bila salah:** ⛔ bila ditiru apa adanya, sistem baru memuat tanggal mati yang **tidak
pernah berlaku lagi**, dan tidak seorang pun akan mengingat asalnya. ⛔ Bila dibuang padahal
aturannya tetap, **nomor polis pada awal Januari jatuh ke periode yang salah** — dan nomor polis
tidak dapat ditarik kembali sesudah terbit.

rujukan: naskah `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` §jawaban **P1** nomor 3; bersinggungan
dengan **P54** *(`TANGGAL_CLOSING` satu baris global)*

⚠️ Pemutusnya `[work owner]` bersama `[Finance]`; dicatat di lembar ini karena faktanya berasal
dari basis data.

**Jawaban:**

> *(tulis di sini)*

---

## Sesudah Anda menjawab

⭐ Kembalikan lembar ini apa adanya — **tidak perlu dirapikan**. Jawaban setengah lengkap tetap berguna; yang tidak berguna adalah lembar yang ditahan sampai lengkap.

⭐⭐ **Tidak ada lagi yang tertahan.** `[terverifikasi]` 2026-09-22 — P1 dan P29 terjawab; **penahan proyek nol**. Dua pertanyaan baru di bawah *(P61, P62)* **tidak menahan** — keduanya perlu dipastikan sebelum jalur terkait dibangun ulang, bukan sebelum pekerjaan dimulai.

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat ini semula berbunyi:
> > *"⚠️ **Yang tertahan selama pertanyaan ini belum terjawab:** **P1** dan **P29** menahan seluruh penulisan spesifikasi penyimpanan data. Tiga sisanya penting tetapi tidak menahan."*

---

*Disusun dari pembacaan berkas ekspor sistem lama, tiga putaran. Tidak ada pertanyaan di lembar ini yang dijawab sendiri oleh tim migrasi.*
