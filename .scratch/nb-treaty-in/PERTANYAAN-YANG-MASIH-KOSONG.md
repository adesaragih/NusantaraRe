# Pertanyaan yang MASIH KOSONG
## Migrasi Treaty Inward — Realisasi & Endorsement · modul *NB Treaty In*

⛔⛔ **CUPLIKAN USANG — jangan dipakai sebagai neraca.** `[penyimpangan sadar]` 2026-09-22

Berkas ini adalah cuplikan **51 pertanyaan** *(28 terjawab, 21 kosong)*. Sejak itu proyek sudah
sampai **P60**, dan **P18 ditarik**. ⭐ **Neraca yang sahih ada di lembar per pemilik**, bukan di
sini. Cuplikan ini dipertahankan sebagai catatan, dan hanya butir P18 yang disamakan.

⭐ Bunyi semula: *"⭐ **Cuplikan keadaan 2026-09-22.** Dari **51 pertanyaan**, **28 sudah terjawab**
dan ⛔ **21 masih kosong**. Berkas ini memuat **hanya yang kosong**."*

> ⚠️ **Berkas ini SALINAN, bukan sumbernya.** Sumbernya tetap lembar per pemilik *(`PERTANYAAN-untuk-….md`)*, dan **28 jawaban yang sudah ada tinggal di sana.**
> ⭐ **Boleh dijawab di sini** — jawabannya akan digabungkan kembali. ⛔ Yang tidak boleh: menjawab pertanyaan yang sama di dua tempat dengan isi berbeda.

### ⭐⭐ Yang MENAHAN seluruh pekerjaan — **NIHIL**  — ⚠️ *semula tiga: P18 ditarik, lalu P1 dan P29 terjawab*

| # | Pemilik | Yang ditahan |
| --- | --- | --- |
| ✅ ~~**P1**~~ | ~~DBA~~ | ✅ **TERJAWAB 2026-09-22** — empat naskah diterima |
| ⭐ ~~**P18**~~ | ~~pemilik export Pega~~ | ⭐ **DITARIK 2026-09-22** — isinya ada di ekspor |
| ⛔ **P29** | **DBA** | kolom pemuat JSON |

⭐ **Selama keduanya kosong, spesifikasi penyimpanan tidak dapat ditulis.** ⚠️ Bagian lain pekerjaan tetap berjalan.

> ⚠️ **RALAT** 2026-09-22 — semula: *"⭐ **Selama ketiganya kosong, spesifikasi penyimpanan dan perhitungan uang tidak dapat ditulis sama sekali.**"* ⭐ **Perhitungan uang tidak lagi tertahan** — rumusnya terbaca di ekspor.

### Sebaran

| Pemilik | Kosong | Nomornya |
| --- | ---: | --- |
| ⛔ **DBA** | **4** | P2 · P3 · P4 · P49  — ⚠️ *semula 6; **P1** dan **P29** terjawab. ⭐ Dua butir **baru** P61 dan P62 ada di lembar DBA, tidak di cuplikan ini* |
| ⛔ **pemilik export Pega** | **2** | P30 · P42  — ⚠️ *semula 3; **P18 ditarik***  |
| **Finance** | **3** | P20 · P21 · P43 |
| **Product & Underwriting** | **6** | P20 · P21 · P22 · P44 · P45 · P46 |
| **pengembang Pega lama** | **3** | P15 · P17 · P47 |

**Cara menjawab:** tulis di bawah tiap pertanyaan. Baris `rujukan:` boleh diabaikan. ⚠️ *"Tidak tahu"* adalah jawaban yang sah. ⛔ Nomor jangan diubah.

---

# Untuk DBA  ·  5 kosong  — ⚠️ *semula 6; **P1 terjawab***

## P1 — ✅ TERJAWAB 2026-09-22  ⭐ **TIDAK MENAHAN**

> ✅ **Naskah empat stored procedure diterima work owner.** `[data DBA]` `[terverifikasi]`
> Judul butir ini semula berbunyi ⛔ *"P1 — Kami butuh isi tiga program penyimpan data di basis
> data ⛔⛔ **MENAHAN**"*. ⭐ **Penahan proyek kini nol.** Jawaban lengkapnya di
> `PERTANYAAN-untuk-DBA.md` §P1 — berkas ini hanya cuplikan.

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

## P29 — ✅ TERJAWAB 2026-09-22  ⭐ **TIDAK MENAHAN**

> ✅ **Terjawab, dan lingkupnya mengecil.** `[keputusan work owner]` `[data DBA]` Judul butir ini
> semula berbunyi ⛔ *"P29 — Sebuah kolom yang namanya berarti "jenis usaha" ternyata menyimpan
> seluruh data kontrak. Benarkah?  ⛔⛔ **MENAHAN**"*. Pembacaan pindah ke view relasional
> `POOLDATA.TREATYINDETAILJOINEDM` *(39 kolom)*; data kontrak **sudah relasional** di
> `POOLDATA.TREATY_IN`. Jawaban lengkapnya di `PERTANYAAN-untuk-DBA.md` §P29.

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

---

## P2 — Apakah ketiga program itu menyelesaikan penyimpanannya sendiri, atau menunggu aplikasi?

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

## P3 — Empat nama tabel muncul dalam dua ejaan. Satu tabel atau dua?

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

## P4 — Satu kolom pada tabel riwayat tidak pernah diisi dari modul ini. Siapa mengisinya?

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

## P49 — Satu kolom riwayat tidak pernah diisi dari sistem lama. Apakah ada yang lain yang mengisinya?

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

# Untuk pemilik export Pega  ·  2 kosong  — ⚠️ *semula 3; **P18 ditarik***

## P18 — ⛔ DITARIK OLEH TIM MIGRASI  ⭐ **TIDAK MENAHAN**

> ⛔⛔ **DITARIK.** `[penyimpangan sadar]` 2026-09-22 — **tidak perlu dijawab siapa pun.**
>
> Isi langkah penetapan nilai **ada di dalam ekspor**, di tag `PropertiesName`/`PropertiesValue`
> — **tanpa awalan `py`**. Tim migrasi memeriksa varian **dengan** awalan, menemukan nol, dan
> mempercayainya. ⭐ **Sebabnya dua huruf.** Lihat `VERIFIKASI-P18.md`.
>
> Terukur: **938 dari 939** langkah `Property-Set` membawa isinya sendiri; **2.481** pasangan
> nama=nilai terisi; **707** pada ke-51 aturan yang semula diminta; **nol** tanda terpotong.
>
> ⚠️ ⭐ **Sidik jari kekeliruannya ada di baris `rujukan:` pertanyaan ini sendiri** —
> *"`PropertiesValue` = 74 dengan salinan editor, **0** tanpa"*. Angka 74 itu dihitung dari atribut
> `REPEATINGINDEX="PropertiesValue"` di dalam blok penyunting, **bukan** dari tag `PropertiesValue`.
> Tag yang sebenarnya tidak pernah diperiksa.
>
> ---
>
> **Bunyi pertanyaan yang ditarik — dikutip utuh:**
>
> > ## P18 — Isi 268 langkah penghitungan tidak ada di dalam berkas yang kami terima. Bisa dikirim ulang?  ⛔⛔ **MENAHAN**
>
> > Sistem lama punya langkah-langkah yang **menetapkan nilai** — hasil perhitungan premi, bagian,
> > pengurangan, dan pajak. Kami dapat melihat **bahwa** langkah itu ada, dan **urutannya**,
> > ⛔ **tetapi tidak satu pun rumus atau nilai yang ditetapkannya ikut terkirim.** Yang ada hanya
> > salinan sementara yang tertinggal dari layar penyunting, dan salinan itu **terpotong di tengah**
> > sehingga tidak dapat dipakai.
>
> > **Konteks:** tanpa isi langkah-langkah ini, perhitungan uang sistem lama tidak dapat ditiru —
> > hanya ditebak.
>
> > ⭐ **Lingkupnya sudah kami persempit lebih dulu, supaya Anda tidak mengerjakan yang tidak kami
> > butuhkan.** Folder yang kami terima memuat 92 aturan dengan 938 langkah, tetapi sebagian besar
> > bukan milik modul ini:
>
> > | Tahap penyaringan | Aturan | Langkah |
> > | --- | ---: | ---: |
> > | seluruh isi folder | 92 | 938 |
> > | dibuang — tidak terjangkau dari titik masuk mana pun | −28 | −561 |
> > | dibuang — tidak dibangun ulang karena penyimpanan JSON ditinggalkan | −5 | −109 |
> > | ⭐ **yang kami minta** | ⭐ **51** | ⭐ **268** |
>
> > ⭐ **Turun 72 % dari angka yang mungkin Anda dengar sebelumnya.**
>
> > **Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — ekspor ulang **51 aturan** pada lampiran
> > terpisah, dengan **isi langkah disertakan**.
>
> > ⭐ Bila tidak seluruhnya dapat dikerjakan, kerjakan **dari urutan teratas lampiran**: sepuluh
> > teratas sudah mencakup **56 %** dari seluruh langkah, dua puluh teratas **76 %**.
>
> > ⚠️ Bila ekspor jenis ini memang **tidak pernah** menyertakan isi langkah, katakan begitu, dan
> > sebutkan **cara lain yang tersedia** *(cetakan layar aturan, dokumentasi, atau akses baca ke
> > sistem lama)*.
>
> > **Dampak bila salah:** ⛔⛔ seluruh perhitungan uang dibangun dari tebakan, dan **selisihnya baru
> > muncul di laporan keuangan berbulan-bulan kemudian**.
>
> > ⛔ **Satu akibat lagi yang baru diketahui:** sembilan medan pada layar Kepala Departemen
> > **wajib diisi tetapi terkunci** — pengguna tidak dapat mengetiknya. Yang mengisinya adalah
> > langkah-langkah ini. Selama isinya belum ada, ⛔ **layar Kepala Departemen tidak akan dapat
> > disimpan sama sekali**, sehingga butir ini menahan bukan hanya perhitungan uang tetapi juga
> > alur persetujuan.
>
> > rujukan: langkah `Property-Set` di 51 `Activity` terjangkau; `PropertiesValue` = 74 dengan salinan
> > editor, **0** tanpa; `pyStepsCallParams` hanya placeholder; penjangkauan ditelusuri dari `Flow`,
> > `Section`, `Harness`, `FlowAction`, `DataTransform` — ronde 2 baru #1, lingkup diperbarui 2026-09-22
>
> > ⭐ **Lampiran:** `LAMPIRAN-P18-ATURAN-DIMINTA.md` — 51 nama aturan, jumlah langkah masing-masing,
> > diurutkan dari yang terberat.
>
> > **Jawaban:**
>
> > > *(tulis di sini)*
>
> > ---
>
>
> ---
>
> ⛔ **Tidak ada yang perlu dikirim.**

---

## P30 — Dua langkah lagi yang dipanggil tetapi tidak ikut terkirim, dan keduanya di jalur penting

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

**Jawaban:**

> **Separuh tertutup: `CopyToPolicy` tidak dipakai.** `[keputusan work owner]` 2026-09-22
>
> Pemanggilan `CopyToPolicy` dari `Activity\SaveJsonPolisTreatyIn_Act` **tidak digunakan** dan
> **tidak dimigrasi**. Tidak perlu diminta kepada pemilik export.
>
> `[terverifikasi]` Berkasnya memang tidak ada di NB Treaty In; yang ada hanya di tiga modul
> Fac In — `Endorsment Fac In`, `NB FacIn`, `RNW Fac In`. Di NB Treaty In namanya disebut satu kali,
> yaitu pada pemanggilan yang kini dinyatakan tidak dipakai itu.
>
> ---
>
> `[terbuka]` **`SumTSIPremiSpreadRNMMultiCob_Act` masih dicari.**
>
> Dipanggil oleh `Activity\SumTSIPremiSpreadedRNM_Act` — penjumlah nilai pertanggungan dan premi.
> **Tidak ada di seluruh korpus 21 modul**, padahal dipanggil oleh empat modul: `NB Treaty In`,
> `NB FacIn`, `RNW Fac In`, `Endorsment Fac In`.
>
> Berbeda dari `CopyToPolicy`, tidak ada satu pun salinan di modul mana pun untuk dibandingkan,
> dan letaknya di rantai uang. **Permintaan ke pemilik export tetap berdiri untuk butir ini:**
> di bagian mana aturan itu tinggal, supaya dapat diminta terpisah.

---

## P42 — Sembilan belas medan tampil di layar, tetapi tidak satu pun aturan mengisinya. Dari mana datangnya?  ⛔ **penting**

Layar sistem lama menampilkan **sembilan belas keterangan** yang **tidak disentuh oleh satu pun
aturan** yang kami terima — tidak diisi, tidak dibaca, tidak dihitung. Delapan di antaranya adalah
**nilai uang beserta mata uangnya**: batas, retensi, premi bersih, dan pendapatan premi
diperkirakan.

Kami menduga keterangan itu **datang dari satu kolom basis data yang isinya berupa dokumen**, yang
dibongkar saat layar dibuka — bukan dari aturan. Kalau dugaan itu benar, sistem baru yang dibangun
hanya dari aturan **tidak akan pernah menampilkannya.**

**Konteks:** kami tidak dapat memastikan sendiri, sebab isi kolom dokumen itu tidak ikut terkirim
kepada kami.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** benar, keterangan itu datang dari
kolom dokumen · **(b)** datang dari bagian sistem lain yang belum kami terima, sebutkan bagian
mana · **(c)** medan itu memang kosong dan tidak pernah terisi.
⭐ Bila **(a)** atau **(b)**: ⭐ **butuh contoh** — tiga isi kolom itu apa adanya.

**Dampak bila salah:** ⛔ delapan medan uang **selalu kosong** di sistem baru, dan ⛔ **tidak ada
pesan galat** yang memberi tahu — hanya angka yang hilang dari layar.

rujukan: 19 properti layar tanpa kemunculan di 92 `Activity` + 12 `DataTransform` + 75 `When` +
41 `RDBList`; `LIMITVALUE`, `RETENTIONVALUE`, `EPIVALUE`, `NETPREMIVALUE` + empat mata uang
pasangannya; bandingkan kolom `BrowseTreatyInDetail` dan mekanisme `adoptJSONObject` ronde 3 §2.1

**Jawaban:**

> **Terjawab — bukan dari kolom dokumen, melainkan dari view relasional.** `[keputusan work owner]` 2026-09-22
>
> Kesembilan belas medan itu **kolom nyata** pada view `POOLDATA.TREATYINDETAILJOINEDM`, bukan
> hasil perhitungan. Itulah sebabnya tidak satu pun `Activity` atau `DataTransform` menghasilkannya.
>
> Delapan medan uang yang dipersoalkan seluruhnya ada sebagai kolom:
> `LIMITCURRENCY`/`LIMITVALUE` · `RETENTIONCURRENCY`/`RETENTIONVALUE` · `EPICURRENCY`/`EPIVALUE` ·
> `NETPREMICURRENCY`/`NETPREMIVALUE`.
>
> **Uji kecukupan:** `ReportDefinition\BrowseTreatyInDetail.xml` merujuk 33 medan; ke-33-nya ada di
> view; nol yang hilang. `[terverifikasi]`
>
> **Koreksi atas bunyi pertanyaan.** Kalimat "tidak disentuh oleh satu pun aturan" terlalu keras.
> Keempat medan uang dirujuk **33-37 kali** di `ReportDefinition`, sebagai `pyTargetProperty` pada
> kelas integrasi `ASM-FW-GISFW-Int-TREATYINDETAIL` dan kerabatnya. Yang benar: **tidak ada Activity
> atau DataTransform yang menghitungnya** — medan itu dibaca, bukan dihasilkan.
>
> `[terbuka]` Pada sisi Pega keempat medan uang ber-`pyStringType` = **`Text`** — uang disimpan
> sebagai teks, bukan bilangan. Tipe Oracle-nya perlu dipastikan sebelum penguraian di Go
> ditetapkan. Dicatat juga pada P29.


# Untuk Finance

---

# Untuk Finance  ·  3 kosong

## P20 — Dua puluh nomor polis tertulis langsung di dalam aturan penghitung uang. Apakah itu disengaja?  ⛔ **penting**

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

**Jawaban:**

> **Dipindahkan ke modul Fac In — rantai ini bukan milik Treaty.** `[keputusan work owner]` 2026-09-22
>
> Kedua puluh nomor polis itu berada di dalam rantai spreading/premi yang **berasal dari Fac In**.
> Rantai itu tampak di folder NB Treaty In karena layarnya **ditanam bersama** (embed), bukan karena
> logikanya milik Treaty. Butir ini **tidak dijawab di lingkup NB Treaty In**; dipindahkan ke daftar
> pertanyaan modul Fac In.
>
> **Bukti — berkasnya identik sampai byte terakhir.** `[terverifikasi]`
>
> | Aturan | NB Treaty In | NB FacIn | Kelas pemilik |
> | --- | ---: | ---: | --- |
> | `Section\ListSuggest` | 218.982 B | 218.982 B | `ASM-FW-GISFW-Data-PolicyTreatyIn` |
> | `Activity\Protection_Act` | 343.271 B | 343.271 B | `ASM-FW-GISFW-Work` |
> | `Activity\CheckSpreadingProtect_ACT` | 306.703 B | 306.703 B | `ASM-FW-GISFW-Work` |
> | `Activity\cekSpreadingFactIn` | — | — | `ASM-FW-GISFW-Data-Coverage` |
> | `Activity\SumTSIPremiSpreadedRNM_Act` | — | — | `ASM-FW-GISFW-Data-Coverage` |
>
> Penguat: **`EDM Treaty In` memiliki `Protection_Act` sendiri** — 157.304 B, kelas
> `ASM-FW-GISFW-Data-PolicyTreatyIn`, jauh lebih kecil. Versi khusus Treaty memang ada, dan
> NB Treaty In tidak memakainya.
>
> Penguat lain: nama pemanggilnya `cekSpreadingFactIn` — "cek Spreading **FacIn**"; halaman
> `OfferFacIn` muncul 1.635 kali di dalam folder NB Treaty In; dan `SumTSIPremiSpreadedRNM_Act`
> juga dipanggil `IsThereAnyObjectLocation_Act` yang **hanya ada di tiga modul Fac**.
>
> **Sebaran nomor per tahun**, untuk dibawa ke modul Fac: 2023 — 12 · 2024 — 35 · 2025 — 72.
> 119 kemunculan, 20 nomor unik, di `SumTSIPremiSpreadedRNM_Act` (96), `ProtectFIREMBUPA_Act` (12),
> `When\IsErrorSpreading.xml` (11). Pertumbuhannya menunjukkan **kebiasaan berulang**, bukan
> tambalan sekali jalan; di modul Fac nanti pertanyaannya adalah **ciri bersama apa** yang membuat
> kedua puluh polis itu diperlakukan berbeda.
>
> `[terbuka]` **Belum dipastikan dengan pengamatan.** Tidak ada satu pun syarat di sepanjang rantai
> pemanggilan — langkah 21, 17, dan 11 seluruhnya berjalan selalu — dan titik masuknya adalah aksi
> refresh pada kelas Treaty. Secara struktur rantai itu **terjangkau** dari berkas Treaty.
> Dugaan yang menjelaskan keduanya: berkas Treaty tidak memiliki halaman `Coverage` berisi, sehingga
> rantai berjalan di atas halaman kosong dan tidak mengerjakan apa pun. **Wajib dibuktikan** dengan
> satu pengamatan di sistem berjalan sebelum 290 langkah `Data-Coverage` dikeluarkan dari lingkup.

---

## P21 — Dua aturan uang hanya tertulis sebagai catatan. Apakah keduanya benar-benar berlaku?  ⛔ **penting**

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

**Jawaban:**

> **Dipindahkan ke modul Fac In — sama seperti P20.** `[keputusan work owner]` 2026-09-22
>
> Kedua aturan uang itu berada di rantai yang sama dengan P20, dan rantai itu **tidak terjangkau
> dari NB Treaty In**. Butir ini tidak dijawab di lingkup NB Treaty In.
>
> **Bukti jangkauan, sesudah `Protection_Act` versi benar dimasukkan 2026-09-22.** `[terverifikasi]`
>
> | Aturan | Pemanggil di NB Treaty In |
> | --- | --- |
> | `ProtectShareCedant_Act` | **tidak ada — yatim** |
> | `ProtectPremiPolicy_Act` | hanya `CheckSpreadingProtect_ACT`, yang sendirinya **yatim** |
>
> Versi `Protection_Act` yang semula ada di folder ini ternyata **berkas yang salah dimasukkan** —
> versi berkelas `ASM-FW-GISFW-Work` milik Fac, 37 langkah, 10 pemanggilan. Versi Treaty yang benar
> berkelas `ASM-FW-GISFW-Data-PolicyTreatyIn`, 16 langkah, satu pemanggilan
> (`ProtectionNonProp_Act`), dan **tidak memanggil satu pun** aturan pada rantai spreading.
>
> Penguat: kedua berkas **identik byte** dengan versi di NB FacIn dan RNW Fac In —
> `ProtectPremiPolicy_Act` 207.132 B, `ProtectShareCedant_Act` 130.178 B — dan keduanya berkelas
> `ASM-FW-GISFW-Work`, bukan kelas Treaty.
>
> **Yang ikut pindah ke modul Fac In**, dua aturan uang yang hanya tertulis sebagai catatan
> pengembang:
>
> | Catatan | Di dalam |
> | --- | --- |
> | `set protect kalau master polis premi harus 0` | `ProtectPremiPolicy_Act` |
> | `TOTAL PREMI SHARE CEDANT HARUS = PREMI RNM` | `ProtectShareCedant_Act` |
>
> Pertanyaan aslinya tetap berlaku di sana: apakah keduanya **kewajiban keras** yang harus ditolak
> bila dilanggar, **peringatan saja**, atau **bukan aturan** — dan bila keras, siapa yang boleh
> menyimpang.
>
> `[terbuka]` Satu butir ikut pindah: `ProtectShareCedant_Act` juga memuat catatan
> *"set precision 4 buat cek nilai total premi share cedant"*. Presisi 4 desimal itu perlu diperiksa
> bersama ketidakseragaman 4-lawan-8 desimal yang tercatat pada **P46**, ketika modul Fac digarap.

---

## P43 — Delapan angka uang tampil di layar tanpa ada perhitungan yang menghasilkannya. Apakah angkanya benar?  ⛔ **penting**

Di layar kontrak terdapat **delapan medan uang** — batas, retensi, premi bersih, dan pendapatan
premi diperkirakan, masing-masing berikut mata uangnya. ⛔ **Tidak satu pun aturan perhitungan
yang kami terima menghasilkan angka-angka itu.**

Artinya salah satu dari dua hal: angka itu **datang jadi dari basis data** tanpa dihitung ulang,
atau ia **dihitung di tempat yang belum kami lihat**.

**Konteks:** kami perlu tahu apakah angka yang dilihat pengguna itu **dihitung** atau **disalin**,
sebab keduanya berperilaku berbeda ketika data dasarnya berubah.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** angka itu disalin apa adanya dari data
kontrak, tidak pernah dihitung ulang · **(b)** dihitung, dan seharusnya berubah bila data dasarnya
berubah · **(c)** sebagian disalin sebagian dihitung, sebutkan mana.
⭐ Dan satu hal yang sangat membantu: **apakah pernah ditemukan angka di layar berbeda dari angka
di laporan?**

**Dampak bila salah:** ⛔ sistem baru **menghitung ulang** angka yang seharusnya tetap, atau
sebaliknya **membekukan** angka yang seharusnya mengikuti perubahan — dan selisihnya baru terlihat
saat penagihan.

rujukan: `pxCurrency` 113 medan di 31 berkas layar; delapan properti uang dalam daftar 19
properti hanya-layar — ronde 4 Bab A.6

**Jawaban:**

> **(a) Angka disalin apa adanya, tidak pernah dihitung ulang.** `[keputusan work owner]` 2026-09-22
>
> Sistem baru **menampilkan nilai yang tersimpan**, tidak menghitung ulang dari data dasar.
>
> **Pertanyaan pendamping — "pernahkah angka di layar berbeda dari angka di laporan?" — dijawab
> TIDAK.** Belum pernah ditemukan selisih. Artinya layar dan laporan selama ini membaca sumber yang
> sama, dan menyalin apa adanya tidak melestarikan penyimpangan apa pun.
>
> **Dasar dari korpus.** `[terverifikasi]` Kedelapan medan uang adalah **kolom tersimpan** pada view
> `POOLDATA.TREATYINDETAILJOINEDM` — `LIMITCURRENCY`/`LIMITVALUE`,
> `RETENTIONCURRENCY`/`RETENTIONVALUE`, `EPICURRENCY`/`EPIVALUE`,
> `NETPREMICURRENCY`/`NETPREMIVALUE`. Tidak satu pun `Activity` atau `DataTransform` di NB Treaty In
> menghasilkannya. Lihat **P29** dan **P42**.
>
> **Akibat yang menguntungkan:** butir ini **tidak lagi bergantung pada P18**.
> Dengan menyalin apa adanya, penyajian kedelapan medan dapat ditulis ke dalam spec sekarang juga.
>
> ⚠️ **RALAT** 2026-09-22 — satu kalimat di sini ditarik: ⛔ *"Bila angka harus dihitung ulang,
> rumusnya wajib diketahui lebih dulu — dan rumus itu justru yang belum terkirim."* ⭐ **Rumusnya
> terkirim**, dan terbaca penuh. **P18 ditarik.**
>
> `[terbuka]` Nilai disimpan dengan presisi penuh, tanpa pembulatan di lapisan repository
> (lihat **P29**). Format penyajian di layar — jumlah desimal dan pemisah ribuan — belum ditetapkan
> dan perlu dipastikan sebelum bab Acceptance Criteria ditulis. **Tidak memblokir.**

---

# Untuk Product & Underwriting  ·  6 kosong

## P20 — Dua puluh nomor polis tertulis langsung di dalam aturan penghitung uang. Apakah itu disengaja?  ⛔ **penting**

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

**Jawaban:**

> **Dipindahkan ke modul Fac In — rantai ini bukan milik Treaty.** `[keputusan work owner]` 2026-09-22
>
> Kedua puluh nomor polis itu berada di dalam rantai spreading/premi yang **berasal dari Fac In**.
> Rantai itu tampak di folder NB Treaty In karena layarnya **ditanam bersama** (embed), bukan karena
> logikanya milik Treaty. Butir ini **tidak dijawab di lingkup NB Treaty In**; dipindahkan ke daftar
> pertanyaan modul Fac In.
>
> **Bukti — berkasnya identik sampai byte terakhir.** `[terverifikasi]`
>
> | Aturan | NB Treaty In | NB FacIn | Kelas pemilik |
> | --- | ---: | ---: | --- |
> | `Section\ListSuggest` | 218.982 B | 218.982 B | `ASM-FW-GISFW-Data-PolicyTreatyIn` |
> | `Activity\Protection_Act` | 343.271 B | 343.271 B | `ASM-FW-GISFW-Work` |
> | `Activity\CheckSpreadingProtect_ACT` | 306.703 B | 306.703 B | `ASM-FW-GISFW-Work` |
> | `Activity\cekSpreadingFactIn` | — | — | `ASM-FW-GISFW-Data-Coverage` |
> | `Activity\SumTSIPremiSpreadedRNM_Act` | — | — | `ASM-FW-GISFW-Data-Coverage` |
>
> Penguat: **`EDM Treaty In` memiliki `Protection_Act` sendiri** — 157.304 B, kelas
> `ASM-FW-GISFW-Data-PolicyTreatyIn`, jauh lebih kecil. Versi khusus Treaty memang ada, dan
> NB Treaty In tidak memakainya.
>
> Penguat lain: nama pemanggilnya `cekSpreadingFactIn` — "cek Spreading **FacIn**"; halaman
> `OfferFacIn` muncul 1.635 kali di dalam folder NB Treaty In; dan `SumTSIPremiSpreadedRNM_Act`
> juga dipanggil `IsThereAnyObjectLocation_Act` yang **hanya ada di tiga modul Fac**.
>
> **Sebaran nomor per tahun**, untuk dibawa ke modul Fac: 2023 — 12 · 2024 — 35 · 2025 — 72.
> 119 kemunculan, 20 nomor unik, di `SumTSIPremiSpreadedRNM_Act` (96), `ProtectFIREMBUPA_Act` (12),
> `When\IsErrorSpreading.xml` (11). Pertumbuhannya menunjukkan **kebiasaan berulang**, bukan
> tambalan sekali jalan; di modul Fac nanti pertanyaannya adalah **ciri bersama apa** yang membuat
> kedua puluh polis itu diperlakukan berbeda.
>
> `[terbuka]` **Belum dipastikan dengan pengamatan.** Tidak ada satu pun syarat di sepanjang rantai
> pemanggilan — langkah 21, 17, dan 11 seluruhnya berjalan selalu — dan titik masuknya adalah aksi
> refresh pada kelas Treaty. Secara struktur rantai itu **terjangkau** dari berkas Treaty.
> Dugaan yang menjelaskan keduanya: berkas Treaty tidak memiliki halaman `Coverage` berisi, sehingga
> rantai berjalan di atas halaman kosong dan tidak mengerjakan apa pun. **Wajib dibuktikan** dengan
> satu pengamatan di sistem berjalan sebelum 290 langkah `Data-Coverage` dikeluarkan dari lingkup.

---

## P21 — Dua aturan uang hanya tertulis sebagai catatan. Apakah keduanya benar-benar berlaku?  ⛔ **penting**

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

**Jawaban:**

> **Dipindahkan ke modul Fac In — sama seperti P20.** `[keputusan work owner]` 2026-09-22
>
> Kedua aturan uang itu berada di rantai yang sama dengan P20, dan rantai itu **tidak terjangkau
> dari NB Treaty In**. Butir ini tidak dijawab di lingkup NB Treaty In.
>
> **Bukti jangkauan, sesudah `Protection_Act` versi benar dimasukkan 2026-09-22.** `[terverifikasi]`
>
> | Aturan | Pemanggil di NB Treaty In |
> | --- | --- |
> | `ProtectShareCedant_Act` | **tidak ada — yatim** |
> | `ProtectPremiPolicy_Act` | hanya `CheckSpreadingProtect_ACT`, yang sendirinya **yatim** |
>
> Versi `Protection_Act` yang semula ada di folder ini ternyata **berkas yang salah dimasukkan** —
> versi berkelas `ASM-FW-GISFW-Work` milik Fac, 37 langkah, 10 pemanggilan. Versi Treaty yang benar
> berkelas `ASM-FW-GISFW-Data-PolicyTreatyIn`, 16 langkah, satu pemanggilan
> (`ProtectionNonProp_Act`), dan **tidak memanggil satu pun** aturan pada rantai spreading.
>
> Penguat: kedua berkas **identik byte** dengan versi di NB FacIn dan RNW Fac In —
> `ProtectPremiPolicy_Act` 207.132 B, `ProtectShareCedant_Act` 130.178 B — dan keduanya berkelas
> `ASM-FW-GISFW-Work`, bukan kelas Treaty.
>
> **Yang ikut pindah ke modul Fac In**, dua aturan uang yang hanya tertulis sebagai catatan
> pengembang:
>
> | Catatan | Di dalam |
> | --- | --- |
> | `set protect kalau master polis premi harus 0` | `ProtectPremiPolicy_Act` |
> | `TOTAL PREMI SHARE CEDANT HARUS = PREMI RNM` | `ProtectShareCedant_Act` |
>
> Pertanyaan aslinya tetap berlaku di sana: apakah keduanya **kewajiban keras** yang harus ditolak
> bila dilanggar, **peringatan saja**, atau **bukan aturan** — dan bila keras, siapa yang boleh
> menyimpang.
>
> `[terbuka]` Satu butir ikut pindah: `ProtectShareCedant_Act` juga memuat catatan
> *"set precision 4 buat cek nilai total premi share cedant"*. Presisi 4 desimal itu perlu diperiksa
> bersama ketidakseragaman 4-lawan-8 desimal yang tercatat pada **P46**, ketika modul Fac digarap.

---

## P22 — Sembilan puluh satu bagian layar dibuat agar tidak pernah muncul, dan tiga puluh delapan dikunci permanen. Masih perlu?

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

**Jawaban:**

> **DITUTUP — terwakili sepenuhnya oleh P44 dan P45.** `[keputusan work owner]` 2026-09-22
>
> Butir ini adalah induk dari dua pertanyaan ronde 4, dan keduanya sudah diputuskan:
>
> | Bagian P22 | Diputuskan di |
> | --- | --- |
> | 91 bagian layar yang dimatikan | **P44** — 80 dibuang tanpa ditanyakan; 4 wadah berisi 104 medan dibawa ke Product & Underwriting |
> | 38 medan yang dikunci permanen | **P45** — ikuti apa adanya, seluruhnya tetap terkunci |
>
> ---
>
> **Koreksi atas bunyi pertanyaan.** `[terverifikasi]` P22 menyatakan *"seluruh tiga puluh delapan
> ada di layar Kepala Departemen"*. Yang benar: **36 di layar Kepala Departemen, 2 di layar biasa** —
> satu di `Section\DetailPolicyTreatyIn`, satu di `Section\GeneralPolicyTreatyIn`. Penguncian itu
> karenanya bukan semata-mata sifat jabatan.
>
> ---
>
> **Pertanyaan turunan — "apakah Kepala Departemen hanya boleh melihat, tidak mengubah?"**
>
> > **Tidak.** Kepala Departemen mengisi **tujuh medan** sendiri, lalu memutuskan.
>
> `[terverifikasi]` Disilangkan dari daftar medan wajib (**P47**) dan daftar medan terkunci (**P45**)
> pada layar yang sama:
>
> | | Jumlah |
> | --- | ---: |
> | medan wajib diisi | 16 |
> | di antaranya terkunci | 9 |
> | **wajib dan dapat diketik** | **7** |
>
> Ketujuhnya: `Claim` · `Deduction1` · `EndDate` · `PremiOgp` · `PremiOnp` · `StartDate` ·
> `StatementDate`.
>
> Sembilan medan yang terkunci terisi sendiri oleh perhitungan — dan rumusnya **terbaca di ekspor**.
>
> ⚠️ **RALAT** 2026-09-22 — kalimat berikut ditarik: ⛔ *"rumusnya berada di dalam langkah yang
> diminta pada **P18**. Selama itu belum ada, ketujuh medan ini dapat diisi tetapi layar **tetap
> tidak dapat disimpan**, karena sembilan medan wajib lainnya kosong."* ⭐ `[terverifikasi]` Layar
> **dapat disimpan**: enam dari sembilan diisi langkah yang terbaca, tiga sisanya diketik
> underwriter di layar sebelumnya. **P18 ditarik.**
>
> **Untuk perancangan peran di sistem baru:** Kepala Departemen adalah peran **penyunting sekaligus
> pemutus**, bukan pengamat. Tujuh medan itu wewenangnya, dan tidak boleh dijadikan hanya-baca.

---

## P44 — Sembilan puluh satu bagian layar dimatikan, tetapi delapan puluh di antaranya kosong. Yang mana yang masih perlu?

Kami sebelumnya melaporkan bahwa **91 bagian layar** diberi syarat tampil yang tidak mungkin
pernah benar — setara "tampilkan jika 1 sama dengan 2". ⭐ **Sekarang kami tahu isi masing-masing**,
dan gambarannya jauh lebih ringan daripada kesannya:

- ⭐ **80 dari 91 tidak menyembunyikan satu medan pun** — hanya wadah kosong atau label.
- **11 sisanya menyembunyikan 128 medan**, dan ⭐ **104 dari 128 itu berada di empat bagian saja**,
  di dua layar yang nyaris kembar — masing-masing menyembunyikan **26 medan**.

**Konteks:** pertanyaannya menjadi jauh lebih sempit. Kami tidak perlu keputusan tentang 91 hal —
⭐ **hanya tentang sebelas**, dan terutama tentang **empat bagian berisi 26 medan** itu.

**Bentuk jawaban yang diharapkan:** untuk **empat bagian besar** *(daftar jenis usaha dan sumber
bisnis, dan kembarannya untuk retro)* — pilihan ganda: **(a)** memang sudah tidak dipakai, boleh
hilang · **(b)** dimatikan sementara dan seharusnya kembali · **(c)** tidak tahu.
⭐ Untuk 80 yang kosong: **boleh kami abaikan seluruhnya?**

**Dampak bila salah:** ⛔ dua puluh enam medan yang masih dibutuhkan **hilang** dari sistem baru,
atau kami membangun dan merawat bagian yang **tidak pernah dipakai siapa pun**.

rujukan: 91 elemen *(`1=2` 78 · `1=2`/`NEVER` pada `pyContainerVisibleWhen` 11 · 2 gabungan
`&& NEVER`)*; sel di bawahnya `FIELD` 128 · `LABEL` 132 · `LAYOUT` 20 · `SUB_SECTION` 5;
empat elemen 26-medan di `BusinessAndSOBList` dan `BusinessAndSOBListRetro` — ronde 4 Bab B.2

**Jawaban:**

> **Dipersempit: 80 elemen dibuang tanpa ditanyakan; hanya 4 yang dibawa ke Product & Underwriting.**
> `[keputusan work owner]` 2026-09-22
>
> **Sebabnya terbukti struktural, bukan kebetulan.** `[terverifikasi]` Kedua tag pengendali
> kemunculan mengerjakan hal yang berbeda:
>
> | Tag | Mustahil | Menempel pada | Menyembunyikan medan? |
> | --- | ---: | --- | --- |
> | `pyCondition` | 80 | sel tunggal — umumnya label dan hiasan | **tidak satu pun** |
> | `pyContainerVisibleWhen` | 11 | wadah yang memuat medan | **128 medan** |
>
> Karena itu 80 elemen `pyCondition` **tidak perlu ditanyakan kepada siapa pun**: membangunnya
> berbiaya nol, membuangnya juga berbiaya nol. Dibuang dari daftar pekerjaan.
>
> `[terbuka]` **Yang dibawa ke Product & Underwriting hanya empat wadah, berisi 104 medan:**
>
> | Medan | Berkas |
> | ---: | --- |
> | 26 | `BusinessAndSOBList` — wadah pertama |
> | 26 | `BusinessAndSOBList` — wadah kedua |
> | 26 | `BusinessAndSOBListRetro` — wadah pertama |
> | 26 | `BusinessAndSOBListRetro` — wadah kedua |
>
> Pertanyaannya satu kalimat: **"Empat bagian layar ini berisi 104 medan dan sudah dimatikan.
> Masih dipakai, atau sudah ditinggalkan?"**
>
> `[terbuka]` Tujuh wadah sisanya menyembunyikan 24 medan — `Installments_ReadOnly` (6),
> `DetailPolicyTreatyIn` (4), `GeneralPolicyTreatyIn` (4), `SpreadingRiskList` (4),
> `DetailPolicyTreatyInNonProportionalEDM` (2), `SFAPortalOpportunities` dan
> `SFAPortalOpportunitiesHeader`. Dilampirkan sebagai daftar untuk ditandai ya/tidak, **tanpa
> pembahasan** — bobotnya terlalu kecil untuk memakan waktu pertemuan.
>
> **Bila jawabannya "sudah ditinggalkan": 128 medan tidak perlu dibangun.**

---

## P45 — Tiga puluh delapan medan dikunci permanen. Haruskah tetap terkunci?

**Tiga puluh delapan medan** diberi setelan hanya-baca yang **selalu berlaku** — tidak pernah dapat
diisi siapa pun, dalam keadaan apa pun. ⭐ **Tiga puluh enam di antaranya ada di dua layar Kepala
Departemen.**

Yang terbaca namanya, sebagai contoh: **"Statement Period"**, **"Survey Report"**, dan
**"Statement Date"**.

**Konteks:** sebelumnya kami menduga seluruh layar Kepala Departemen memang dibuat hanya-lihat.
⭐ **Ternyata tidak** — yang dikunci adalah **medan tertentu, satu per satu**, dan dua di antaranya
bahkan berada di layar biasa, bukan layar Kepala Departemen.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** benar, Kepala Departemen memang hanya
melihat, tidak mengubah · **(b)** seharusnya ada yang dapat diubah, sebutkan medan mana ·
**(c)** penguncian ini peninggalan, sudah tidak dimaksudkan.
⭐ Khususnya: **apakah "Statement Period" dan "Statement Date" memang tidak boleh diubah?**

**Dampak bila salah:** ⛔ pengguna **tidak dapat memperbaiki** keterangan yang keliru dan harus
menempuh jalan memutar, atau sebaliknya medan yang seharusnya terkunci **menjadi dapat diubah**.

rujukan: `pyReadOnlyCondition` selalu-benar **38**, pada elemen `pyUserData` yang induknya sel
`pyType=FIELD`; `DetailDeptHeadTreatyIn_UW` 18 · `GeneralDeptHeadTreatyIn_UW` 18 ·
`DetailPolicyTreatyIn` 1 · `GeneralPolicyTreatyIn` 1 — ronde 4 Bab B.3

**Jawaban:**

> **Ikuti apa adanya — 38 medan tetap terkunci permanen.** `[keputusan work owner]` 2026-09-22
>
> Perilaku sistem lama disalin tanpa perubahan. Tidak ada medan yang dibuka, tidak ada kewajiban
> yang dicabut.
>
> **Sensus dipastikan ulang — 38 tepat.** `[terverifikasi]`
>
> | Berkas | Terkunci |
> | --- | ---: |
> | `DetailDeptHeadTreatyIn_UW` | 18 |
> | `GeneralDeptHeadTreatyIn_UW` | 18 |
> | `DetailPolicyTreatyIn` | 1 |
> | `GeneralPolicyTreatyIn` | 1 |
>
> `pyReadOnlyCondition` yang selalu benar: `1==1` 32x, `1=1` 2x, `1 = 1` 2x, `ALWAYS` 2x.
> Setiap kunci mengenai **tepat satu medan**, bukan bagian atau tab.
>
> Medan yang terkunci di kedua layar Dept Head: `Deduction2` `ExcessLoss` `NetPremium`
> `OutstandingClaim` `OveriddingCommOgp` `OveriddingCommOnp` `PPHValue` `ResultOgp1` `ResultOgp2`
> `ResultOnp1` `ResultOnp2` `RiCommOgp` `RiCommOnp` `SalvageValue` `StatementType`.
>
> ---
>
> `[terbuka]` **Sembilan medan wajib diisi TETAPI terkunci.** Disilangkan dengan daftar medan wajib
> pada **P47**: `Deduction2` · `ExcessLoss` · `OutstandingClaim` · `OveriddingCommOgp` ·
> `OveriddingCommOnp` · `ResultOnp1` · `RiCommOgp` · `RiCommOnp` · `SalvageValue`.
>
> Kesembilannya wajib diisi tetapi tidak dapat diketik pengguna. Itu hanya dapat bekerja bila ada
> yang **mengisinya secara otomatis**. Rumus pengisinya berada di dalam **268 langkah yang belum
> terkirim (P18)**.
>
> **Akibat yang harus diketahui sebelum go-live:** selama P18 kosong, sistem baru akan menampilkan
> sembilan medan wajib yang selamanya kosong, dan layar Dept Head **tidak akan dapat disimpan**.
> Ini menjadikan P18 penahan bukan hanya bagi perhitungan uang, tetapi juga bagi **alur persetujuan
> tingkat Dept Head**.
>
> Enam medan terkunci lainnya — `NetPremium` `PPHValue` `ResultOgp1` `ResultOgp2` `ResultOnp2`
> `StatementType` — tidak wajib, sehingga tidak menahan penyimpanan.

---

## P46 — Sebuah medan terisi otomatis dengan kata "Inclusive". Apa artinya, dan kapan berlaku?

Layar sistem lama hampir **tidak pernah** mengisi medan dengan nilai awal — hanya **sebelas medan**
di seluruh modul yang punya nilai bawaan. Sembilan di antaranya diisi angka **nol**; dua sisanya
diisi kata **"Inclusive"**.

**Konteks:** kata itu satu-satunya nilai bawaan yang berupa **istilah bisnis**, bukan angka. Kami
tidak dapat menurunkan artinya dari berkas, dan nilai bawaan menentukan apa yang tersimpan ketika
pengguna tidak mengubah apa pun.

**Bentuk jawaban yang diharapkan:** satu paragraf — **apa arti "Inclusive"** pada medan itu,
**apa nilai lain yang mungkin**, dan **apakah "Inclusive" memang bawaan yang benar** atau
seharusnya kosong sampai pengguna memilih.

**Dampak bila salah:** ⚠️ kontrak yang pengguna tidak menyentuh medan itu **tersimpan dengan
pilihan yang tidak pernah ia buat**.

rujukan: `pyDefaultValue` berisi pada 11 sel `FIELD`; nilai `0` **9×**, `Inclusive` **2×**
*(`DetailPolicyTreatyIn` + 1 berkas)* — ronde 4 Bab A.5

**Jawaban:**

> **Ikuti apa adanya — perilaku sistem lama disalin utuh, termasuk ketidakseragamannya.**
> `[keputusan work owner]` 2026-09-22
>
> Medan `PolicyTreatyIn.TypeTax` bernilai awal `"Inclusive"` dan menggerakkan rumus potongan pajak
> pada brokerage:
>
> ```
> BrokerageFeeSebenarnya = @if(TypeTax == "Inclusive", Deduction / (102.2/100), Deduction)
> ```
>
> Bila `"Inclusive"`, potongan dibagi 1,022 — pajak 2,2 % dikeluarkan dari angkanya. Bila bukan,
> dipakai apa adanya. Rumus ini ada di **12 tempat pada 7 Activity**; angka `2.2` juga muncul di
> 4 berkas layar. `[terverifikasi]`
>
> **Tiga sifat yang ikut disalin, seluruhnya disadari:**
>
> `[terverifikasi]` **`Exclusive` tidak pernah tertulis di mana pun** — nol kemunculan di seluruh
> modul NB Treaty In. Cabang "bukan Inclusive" hanya tercapai bila medan berisi sesuatu yang lain,
> termasuk kosong.
>
> `[terverifikasi]` **Pembanding berupa teks persis.** `TypeTax == "Inclusive"`. Huruf kecil, spasi
> di belakang, atau ejaan lain jatuh ke cabang kedua dan brokerage **tidak** dibagi 1,022 —
> selisih 2,2 % pada setiap potongan, tanpa pesan galat.
>
> `[terbuka]` **Presisi tidak seragam untuk rumus yang sama.** Empat desimal di
> `InputPolicyTreatyInDetail_NonProp` dan `InputPolicyTreatyOutDetail_NonProp`; delapan desimal di
> `InputPolicyTreatyInDetail_preACT`, `InsertToTreatyXOLList`, `InsertToTreatyOutXOLList`, dan
> `SetPPNPPH`. Akibatnya **kontrak yang sama dapat menghasilkan brokerage yang sedikit berbeda,
> bergantung layar mana yang menyimpannya.**
>
> Keputusan "ikuti apa adanya" berarti ketidakseragaman ini **ditiru**, bukan diperbaiki. Perbedaan
> nilainya sudah ada di data sekarang, bukan risiko baru. Dicatat sebagai butir terbuka karena bila
> kelak ditemukan selisih brokerage yang dipersoalkan, **di sinilah asalnya**, dan keputusan ini
> wajib ditinjau ulang bersama Finance.
>
> Angka 2,2 % tetap ditulis di tempat pemakaiannya, tidak dijadikan setelan — sesuai keputusan yang
> sama. Bila tarif berubah, perubahannya menyentuh 12 tempat.

---

# Untuk pengembang Pega lama  ·  3 kosong

## P15 — Sebuah fungsi yang membentuk seluruh muatan data tidak ada naskahnya. Apa isinya?

Ada satu fungsi buatan sendiri yang **mengubah data satu kasus menjadi satu teks besar**, dan teks
itulah yang disimpan ke basis data sebagai muatan utama. ⛔ **Naskah fungsi itu tidak ada di dalam
ekspor mana pun.**

**Konteks:** data lama tersimpan dalam bentuk yang dihasilkan fungsi ini. Sistem baru harus dapat
**membacanya kembali**.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — naskah fungsi tersebut, **atau** beberapa
**contoh nyata** hasilnya dari basis data *(3–5 contoh sudah cukup untuk menurunkan bentuknya)*.

**Dampak bila salah:** ⛔ **seluruh data lama tidak terbaca** oleh sistem baru.

rujukan: `@ASM.GetPageJSONString()` di `Activity\SaveJsonPolisTreatyIn_Act.xml` langkah 6 — OQ-012

**Jawaban:**

> **Naskah diterima, dan fungsinya tidak dimigrasi.** `[keputusan work owner]` 2026-09-22
>
> Isi `@ASM.GetPageJSONString()`, diberikan work owner:
>
> ```java
> PublicAPI tools = (PegaAPI) ThreadContainer.get().getPublicAPI();
> ClipboardPage stepPage = tools.getStepPage();
> String retValue = stepPage.getJSON(false);
> return retValue;
> ```
>
> **Artinya: fungsi ini tidak memilih apa pun.** Ia mengubah **seluruh halaman langkah** menjadi
> JSON apa adanya. Tidak ada penyaringan medan, tidak ada pemetaan, tidak ada bentuk yang dirancang.
>
> Itu menjelaskan mengapa isi kolom `JSONDATA` tidak dapat diterka dari aturan mana pun: bentuknya
> bukan skema, melainkan **potret halaman kerja pada saat penyimpanan** — apa pun yang kebetulan ada
> di sana ikut masuk.
>
> **Keputusan: penyimpanan pindah ke tabel relasional. JSON tidak dipakai lagi, baik untuk membaca
> maupun menulis.** Melengkapi keputusan pada **P29** yang semula hanya menyangkut pembacaan.
> Fungsi ini tidak dimigrasi.
>
> ---
>
> **Dua akibat yang perlu dicatat.**
>
> ~~`[terbuka]`~~ ✅ **P1 TERJAWAB 2026-09-22.** Butir ini semula berbunyi:
> > *"`[terbuka]` **P1 menjadi lebih penting, bukan berkurang.** Karena JSON tidak ditulis lagi,
> > sistem baru harus menulis langsung ke tabel-tabel yang selama ini diisi
> > `POOLDATA.PEGA_JSON_POLIS_TREATYIN`. Untuk itu kita tetap perlu tahu **tabel dan kolom mana
> > saja yang disentuh procedure itu** — permintaan P1 berdiri utuh untuk ketiga procedure."*
>
> ⭐ **Naskahnya diterima — dan bukan tiga, melainkan EMPAT.** `PEGA_TREATY_IN` menulis **dua**
> tabel *(`M_TREATY_IN` dan `POOLDATA.TREATY_IN` berisi 20 kolom datar)*;
> `PEGA_JSON_POLIS_TREATYIN` menulis **satu** tabel delapan kolom; `PROC_GENERATE_SEQUENCE_NUMBER`
> tidak menulis data bisnis; `PEGA_DELETE_ERROR_KONVERSI` — yang keempat — menghapus.
> Rinciannya di `PERTANYAAN-untuk-DBA.md` §P1.
>
> `[terbuka]` **Isi `JSONDATA` lama kemungkinan memuat jauh lebih banyak daripada data bisnis.**
> Karena `getJSON(false)` menyalin seluruh halaman, dokumen lama dapat memuat properti internal
> Pega dan — perlu diperiksa — **nilai berupa nama orang**. Bila kelak ada pemindahan data dari
> kolom itu, isinya wajib disaring lebih dulu, bukan disalin utuh.

---

## P17 — Aturan yang hanya hidup di layar belum kami baca sama sekali. Apa yang ada di sana?

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

**Jawaban:**

> **TIDAK BERLAKU LAGI — sudah dikerjakan sendiri, tidak perlu ditanyakan.** `[keputusan work owner]` 2026-09-22
>
> Pertanyaan ini meminta **urutan prioritas** pembacaan lapisan layar. Grilling ronde 4 sudah
> mengurai **seluruhnya**: 25 `Section` + 6 `Harness`, **16.635.834 byte, 100 %, nol berkas tidak
> terurai.** `[terverifikasi]`
>
> Yang dihasilkan sudah menjadi pertanyaan tersendiri — **P44** sampai **P47** — jadi tidak ada lagi
> yang perlu diminta dari pengembang Pega lama untuk butir ini.
>
> Sisa pekerjaannya bukan pertanyaan kepada siapa pun, melainkan **penulisan spec**: syarat tampil,
> medan wajib, medan hanya-baca, validasi sisi layar, dan urutan pengisian dituangkan ke dalam
> bab User Stories dan Acceptance Criteria.
>
> Ditutup bukan karena dijawab, melainkan karena **pertanyaannya tidak berlaku lagi**.

---

## P47 — Seratus enam puluh setelan "wajib isi", tetapi tidak satu pun menempel pada medan. Di mana kewajibannya berlaku?

Di ke-31 layar terdapat **160 setelan "wajib isi" yang berisi**. ⛔ **Tetapi tidak satu pun dari
870 medan** di layar-layar itu membawanya. Artinya setelan-setelan itu menempel pada **sesuatu yang
lain** — dan kami tidak dapat memastikan pada apa.

**Konteks:** kami tidak ingin menyimpulkan bahwa modul ini **tidak punya medan wajib** hanya karena
kami mencari di tempat yang salah. Ini persis jenis kekeliruan yang sudah beberapa kali menjerat
kami.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** setelan itu memang menempel pada
wadah/baris, bukan pada medan, jelaskan bagaimana ia berlaku · **(b)** ia peninggalan yang tidak
lagi bekerja · **(c)** ada tempat lain yang seharusnya kami periksa, sebutkan.
⭐ Dan yang paling membantu: ⭐ **daftar medan yang benar-benar wajib diisi** menurut sistem lama —
walau hanya untuk satu layar utama.

**Dampak bila salah:** ⛔ sistem baru **menerima** kontrak yang selama ini ditolak karena ada medan
kosong, atau sebaliknya **menolak** yang selama ini diterima.

rujukan: `pyRequired` 160 berisi di 31 berkas layar *(ronde 2 §2.4)* lawan 0 pada `pyUserData` dari
870 sel `pyType=FIELD` *(ronde 4 Bab A.3)*

**Jawaban:**

> **(a) Setelan itu MEMANG menempel pada medan. Premis pertanyaan keliru.** `[terverifikasi]` 2026-09-22
> Terjawab dari korpus — tidak perlu ditanyakan kepada siapa pun.
>
> **RALAT.** Pertanyaan ini menyatakan *"tidak satu pun dari 870 medan membawanya"*. Penelusuran
> yang menghasilkan kesimpulan itu mencari di `pyUserData`. Setelan `pyRequired` tidak tinggal di
> sana — ia tinggal di dalam blok sel layar, berdampingan dengan `pyValue` yang menyebut nama
> medannya.
>
> Diperiksa ulang: dari 160 `pyRequired=true`, **seluruhnya dapat dipasangkan dengan nama medan**.
> Pemasangan dilakukan dalam jendela blok yang sama, bukan menurut urutan.
>
> **27 medan wajib berbeda, tersebar di 6 layar:**
>
> | Layar | Jumlah | Medan |
> | --- | ---: | --- |
> | `GeneralPolicyTreatyIn` | 22 | `Claim` `ClaimPaymentType` `ClaimType` `Deduction1` `Deduction2` `EndDate` `ID` `IDCurrency` `OutstandingClaim` `OveriddingCommOgp` `OveriddingCommOnp` `PremiOgp` `PremiOnp` `Quartal` `RiCommOgp` `RiCommOnp` `SalvageValue` `StartDate` `StatementDate` `TypeTax` `YearOfQuartal` |
> | `DetailPolicyTreatyIn` | 21 | sama, tanpa `ID` |
> | `DetailDeptHeadTreatyIn_UW` | 16 | `Claim` `Deduction1` `Deduction2` `EndDate` `ExcessLoss` `OutstandingClaim` `OveriddingCommOgp` `OveriddingCommOnp` `PremiOgp` `PremiOnp` `ResultOnp1` `RiCommOgp` `RiCommOnp` `SalvageValue` `StartDate` `StatementDate` |
> | `GeneralDeptHeadTreatyIn_UW` | 16 | sama dengan di atas |
> | `ListSuggest` | 3 | `IsApproved` `ProductionDate` `Suggest` |
> | `InputHistoricalSurveyReportDtl` | 1 | `DateofSurvey` |
>
> **Yang perlu diperhatikan saat menulis spec:**
>
> Layar Dept Head **tidak** mewajibkan `ClaimPaymentType`, `ClaimType`, `IDCurrency`, `Quartal`,
> `TypeTax`, `YearOfQuartal` — enam medan yang wajib di layar admin. Sebaliknya layar Dept Head
> mewajibkan `ResultOnp1`, yang tidak wajib di layar admin.
>
> Jadi kewajiban **berbeda menurut tingkat**, bukan seragam. Itu harus dipertahankan di sistem baru,
> dan tidak boleh disederhanakan menjadi satu daftar tunggal.
>
> `[terbuka]` Beberapa medan wajib itu juga punya syarat `pyRequiredWhen` — kewajiban bersyarat,
> bukan mutlak. Contoh yang sudah diketahui: `ListSuggest.IsApproved` wajib hanya bila dua pengguna
> tertentu yang membuka (lihat **P25**, tempat ID operator ditulis mati diganti peran). Syarat
> bersyarat itu perlu disalin satu per satu ketika bab Acceptance Criteria ditulis.

---

## Sesudah Anda menjawab

⭐ Kembalikan apa adanya — **tidak perlu lengkap**. Jawaban sebagian tetap berguna.

*Cuplikan 2026-09-22, disusun dari lembar per pemilik. Tidak ada pertanyaan yang dijawab sendiri oleh tim migrasi.*
