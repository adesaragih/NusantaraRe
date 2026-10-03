# Daftar issue terbuka — keadaan **25 September 2026**

> **Apa ini.** Satu tempat untuk melihat **semua** yang masih menggantung, dikumpulkan dari tujuh
> dokumen yang selama ini menyimpannya terpisah-pisah. Disusun menurut **siapa yang bisa
> membukanya**, bukan menurut urutan kemunculannya.
>
> ⛔ **Ini BUKAN sumber kebenaran baru.** Tiap butir menunjuk dokumen aslinya. Bila isinya berbeda
> dengan dokumen asal, **yang asal yang benar** dan berkas ini yang salah.
>
> ⛔ **Angka di sini DIKUTIP ULANG, tidak diukur ulang hari ini.** Tiap angka menyebut dokumen yang
> mengukurnya, dan di sanalah perintah auditnya tersimpan (`CLAUDE.md` §3 butir 6).
>
> ✅ **Yang sudah TIDAK terbuka lagi:** seluruh sembilan butir rancangan tabel flat + lima
> turunannya, ditutup **K-063 … K-069** dan **ADR-0007**. Lembar `Butir Terbuka` di
> `08-flat\Tabel-Flat-Lintas-Siklus.xlsx` **sudah sesuai** — baris pertamanya berbunyi
> *"TIDAK ADA BUTIR TERBUKA."* Jangan buka lagi tanpa bukti baru.

---

## 0. Satu layar — semuanya sekaligus

| Golongan | Jumlah | Siapa yang membuka | Berhenti kalau tidak dibuka |
| --- | :-: | --- | --- |
| **A** Keputusan work owner | **6** *(A-1 ✅ tertutup K-070)* | **Anda, hari ini** | perbandingan rumus EDM · fixture · lingkup Life |
| **J** turunan K-070 | **2** *(J-1…J-15 ✅ tertutup K-071 · K-072 · K-073)* | **Anda** | ⛔ **tidak ada lagi** yang menahan spec pemuatan |
| **B** Harus dilihat di layar Pega | **7** | Anda / admin Pega | arti kode transisi · vonis rule mati · premis K-050 |
| **C** Jawaban DBA | **5** *(+T-10)* | DBA | F13 · jalur banding · klep gerbang EDM · aturan "versi terakhir" |
| **D** Ekspor ulang IT | **5** | IT / admin Pega | F10 · F11 · rekonsiliasi · jalur uang EDM |
| **E** Arti kode Product | **1** | Product | klasifikasi produk · layar input |
| **F** Pekerjaan kita yang belum ada spesifikasinya | **3** | kita sendiri | pemuatan data lama |
| **G** Utang kerapian dokumen | **6** | kita sendiri | dokumen saling bertentangan |
| **H** Yang akan `panic` bila tercapai | **7** | — (konsekuensi, bukan tugas) | aplikasi berhenti, bukan salah diam-diam |

✅ **Tidak ada lagi keputusan yang menahan spec program pemuatan.** Bahannya sudah diperbarui
(`08-flat\BAHAN-SPEC-PEMUATAN.md`, revisi 25 September). ⛔ `/to-spec` ber-`disable-model-invocation`
— bahan disiapkan, lalu berhenti.

---

# A. Keputusan work owner — bisa diputuskan hari ini, tanpa menunggu siapa pun

## A-1 ✅ **DITUTUP 25 September 2026 — K-070.** Migrasi **dua fase**

**Jawaban work owner:** seluruh polis **tetap** sasaran akhir, dikerjakan bertahap —
**fase 1 versi terakhir saja** (sekarang), **fase 2 generasi sebelumnya** (setelah sistem jadi).

⛔ **Bukan "versi terakhir saja" sebagai jawaban akhir.** Konsekuensi *"nilai sebelum endorsement
hilang permanen"* yang tertulis di tabel asli **tidak berlaku** — yang berubah hanya urutan waktu,
bukan cakupan.

⚠️ **Harga yang dibayar:** fase 2 harus **menjodohkan ke belakang** terhadap `ROW_UID` yang sudah
terlanjur dibangkitkan di fase 1. Itu **lebih sulit** daripada memuat seluruh versi sekaligus.

📌 **Sepuluh pertanyaan turunan lahir dari keputusan ini → golongan J**, plus satu permintaan DBA
baru (**T-10**).

---

## A-2 Metode deteksi *drift* — nomor versi rule atau hash 23 tag?

**Bahasa sederhana:** untuk tahu apakah dua salinan sebuah rule benar-benar sama, selama ini kita
membandingkan **nomor versinya**. Ternyata cara itu bocor besar.

`[terverifikasi]` (diukur dan dicatat di **K-058**): cara lama **melewatkan 267 dari 354** perbedaan
NB↔EDM — **75,4 %** luput. Dan angka lama **"95 rule berbeda versi"** yang tertulis di
`CLAUDE.md` §4.5 serta `steering\PANDUAN-KERJA.md` §5 **tidak tereproduksi** — hasil ukur ulang
**87 / 85 / 0**.

**Dua akibat kalau dibiarkan:** (a) kita akan terus memakai alat ukur yang meleset tiga dari empat
kali; (b) dua berkas aturan proyek memuat angka yang sudah terbukti keliru.

**Sumber:** **K-058** *(`[pertanyaan terbuka]` — TIDAK diputuskan di sana)*.

---

## A-3 `premium.Calculate` — satu fungsi bercabang, atau empat implementasi?

**Bahasa sederhana:** mesin premi ternyata punya **4 bentuk perhitungan** (rumus dasar · dua-bagian
Start+End · tabel tarif Travel · dekomposisi delta coverage). Semuanya lewat **satu pintu** (Seam 3,
tetap 5 seam — **tidak ada seam baru**). Yang belum diputus: di balik pintu itu satu fungsi
bercabang, atau empat implementasi terpisah?

Ini murni **keputusan desain**, tidak menunggu data apa pun.

**Sumber:** **K-057** area Seam 3 · `04-spec\03-spec-modul-terverifikasi.md`.

---

## A-4 Kapan rumus `Calculate*` EDM dibandingkan dengan NB — sebelum spec, atau saat implementasi?

`[dugaan]` `CalculatePremiPA_FacIn` NB↔EDM diklaim identik secara fungsional, **belum diverifikasi
dengan kontrak 23 tag**. Juga: 5 dari 6 salinan EDM mesin premi berbeda isi — rumus dasarnya sama,
bedanya di tempat lain.

**Sumber:** `04-spec\09-spec-edm-only.md` §7.2 · **K-057** bagian *Belum tuntas*.

---

## A-5 Batas himpunan **Life** untuk modul 6 — **52** atau **66** berkas?

Angkanya bergantung pada penyaring mana yang dipakai. Keduanya sah; memilihnya keputusan lingkup.

**Sumber:** `04-spec\09-spec-edm-only.md` §6.

---

## A-6 Dua rujukan menggantung di `SetErrorMessage_Act`

**Bahasa sederhana:** dua activity (`CountASMGrossPremi_ACT`, `CountASMNetPremi_ACT`) dinyatakan
dikeluarkan. Tetapi `[terverifikasi]` keduanya **masih dipanggil** dari `SetErrorMessage_Act.xml`
baris 518 & 618. Jadi ada dua langkah yang menunjuk ke tempat yang tidak ada lagi.

| Kemungkinan | Akibat |
| --- | --- |
| Langkah pemanggilnya **ikut usang** — sisa kelewat | Diport apa adanya; jalurnya mati sendiri |
| `SetErrorMessage_Act` **perlu ditinjau** | Butuh keputusan tersendiri sebelum diport |

⛔ **Sampai diputus, dilarang:** menghapus langkah pemanggil diam-diam, atau menyimpulkan activity
itu kode mati. Pola sama dengan K-019 dan K-024.

**Sumber:** register baris 424–440 · `_CHECKLIST-KESIAPAN.md` B.1.

---

## A-7 Lima berkas EDM yang tidak dirujuk dari mana pun (**E-Q24**)

`[terverifikasi]` lima berkas EDM tidak disebut rule mana pun. Yang paling menonjol
`SFAPortalEndorsement` (Harness) — `[dugaan]` mungkin dipanggil dari **konfigurasi portal yang tidak
ikut diekspor**, bukan mati.

**Sumber:** `04-spec\09-spec-edm-only.md` §5.1.

---

# B. Hanya bisa dijawab dengan **membuka layar Pega** — selagi sistem lama masih hidup

⚠️ **Golongan ini paling mendesak secara waktu.** Begitu Pega dimatikan, jawabannya hilang selamanya.

| # | Pertanyaan | Yang bergantung padanya |
| :-: | --- | --- |
| **U-1** | Arti **kode transisi 5** — `[terverifikasi]` **908 kemunculan**, terbanyak di antara yang belum tertambat | **F10** · seluruh pemetaan alur activity |
| **U-2** | Arti **kode transisi 1 dan 4** — `[terverifikasi]` **347** dan **120** baris bergerbang | `10-audit\02` §2.3 · setiap activity bercabang |
| **U-3** | Ada berapa rule bernama `IsFire`, `IsPA`, `IsUW`, di kelas apa saja | ⛔ premis **K-050** *"196 predikat identik"* · **F01** |
| **U-4** | Apakah `CountRateRetroCov` kelas `Data-Cargo` pernah ada | **K-060** butir 1 · **F06** *(kembaran tertulis: **T-7**)* |
| **U-5** | Apakah urutan `REPEATINGINDEX` = urutan eksekusi, **sebagai aturan umum** | apakah dua rumus prorata benar-benar mati |
| **U-6** | ✅ *Arti `//` terjawab 01-10-2026: di-remark, tidak dipakai (`KEPUTUSAN-30-09-2026.md` butir 43); sisa: membaca label tiap pemanggil.* Status remark `//` pada **setiap** pemanggil `GetLimitAkseptasi_Act` / `_Act2` | vonis tangga akseptasi EDM · cabang `IsBonding` |
| **U-7** | Resolusi `IsFacRetro` / `OfferFacRetro` — indeks flow menyebut app & workType yang **berbeda** dari kelas di `Embed-Reference-Rule` | ketertelusuran predikat Fac Out |

> ✅ **U-1 terjawab 1 Oktober 2026** (`KEPUTUSAN-30-09-2026.md` butir 39, lewat UI Pega): **`5` = Skip
> Whens**; sekaligus **`2` = Continue Whens**, **`3` = Skip Step** naik ke `[terverifikasi]`. U-2 (kode `1`,
> `4`) dan kode `6` masih terbuka — langkah yang disarankan untuk dilihat: `PERTANYAAN-AKSEPTASI.md` bab E.

> ✅ **U-2 terjawab 1 Oktober 2026** (`KEPUTUSAN-30-09-2026.md` butir 40): **`1` = Jump To Later Step**,
> **`4` = Exit Iteration**, **`6` = Exit Activity**. Peta kode `1`–`6` seluruhnya `[terverifikasi]`. Sisa:
> arti isian kosong (5.953 kemunculan), `PERTANYAAN-AKSEPTASI.md` bab F.
> ✅ **Terjawab juga** (butir 41): **kosong = Continue Whens**. Peta kode transisi selesai.

📌 **U-1 dan U-2 saling menopang.** Kode `2`, `3`, `6` sudah tertambat pada tingkat `[dugaan kuat]`.
Bila U-1 dan U-2 terjawab, **seluruh peta kode transisi naik ke `[terverifikasi]` sekaligus**.

⚠️ **Jangan tawarkan lagi hipotesis "kode 5 = kendali perulangan"** — sudah diuji dan gagal: langkah
pembawa kode 5 beriterasi hanya **9,8 %**, lebih rendah daripada kode 2 (**30,8 %**).

**Sumber:** `10-audit\11-daftar-periksa-ui-pega.md` (U-1…U-4) · **K-060** butir 2 (U-5) ·
**K-045** / P-12 W-3 (U-6) · register baris 2953 (U-7).

---

# C. Menunggu **DBA**

| # | Yang diminta | Memblokir |
| :-: | --- | --- |
| **T-1** | Isi `M_LINK_SERVICE` — ⛔ **hanya** `URL`, `KATEGORI_1`, `KATEGORI_2` | **F13** (cetak RI Slip + email) · **setiap** panggilan endpoint |
| **T-2** | Dua pertanyaan `HISTORYAKSEPTASIPEGA`: (a) apakah dipangkas berkala · (b) proses apa yang mengisi `ID_KOMITE` | jalur banding · flag reject |
| **T-3** | Konfirmasi status **8 nama** yang tiba sebagai `CREATE TABLE`, bukan kode prosedur | hanya ketepatan daftar — ⛔ **bukan** blocker jalur produksi |
| **T-4** | Isi baris `OPENPROTEKSI_EDM` untuk `Param.Type` 1–4 | klep 4 gerbang pembuatan EDM (**K-049**) |
| **T-10** | ⚠️ **Bobotnya TURUN sesudah K-071.** Berapa nilai **`PRODKE` tertinggi** di `JSON_POLIS` dan berapa baris yang dua digit? | ⛔ **tidak lagi memblokir** — J-2 sudah menetapkan urutan memakai `TGL_INPUT`. Angka ini kini **pemeriksa silang**: ia menunjukkan seberapa sering jebakan urut-teks berpeluang terwujud |

**Query T-10** — tanpa data pelanggan, hanya hitungan:

```sql
SELECT MAX(TO_NUMBER(PRODKE)) AS prodke_tertinggi,
       COUNT(CASE WHEN LENGTH(TRIM(PRODKE)) >= 2 THEN 1 END) AS baris_dua_digit,
       COUNT(*) AS total_baris
FROM   POOLDATA.JSON_POLIS
WHERE  PRODKE IS NOT NULL;
```

⛔ **`USERNAME`, `PASSWORD`, `NAMA`, `LOGIN`, `OPERATORID` tidak diminta dan tidak akan disimpan.**
Nama kolomnya boleh disebut; isinya tidak.

**Sumber:** `_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md` — sudah **siap kirim**, tinggal dikirim.

---

# D. Menunggu **IT / administrator Pega**

| # | Yang diminta | Memblokir |
| :-: | --- | --- |
| **T-5** | **10** rule `DecisionTable` diekspor ulang beserta seluruh barisnya, **berikut `pyClassName`** | arah keputusan UW · pemetaan coverage & spreading (F04, F05, F08) |
| **T-6** | **10** rule golongan **G2** diekspor ulang | **F10** · **F11** · rekonsiliasi — ⛔ `CountEdmAdjTSI_Act` **menyentuh jalur uang** |
| **T-7** | Apakah `CountRateRetroCov` kelas `Data-Cargo` pernah ada | **K-060** butir 1 · **F06** |
| **T-8** | Satu ekspor **rule** produksi pada satu titik waktu | ketidaksepadanan versi antar folder — ⚠️ **prioritas turun**, bukan gugur |
| **E-Q1** | Apakah resolusi rule Pega di instalasi ini **peka huruf**? | **12 berkas** EDM (`IsCar` ↔ `IsCAR`) · menentukan angka resmi **1.707 / 354** atau **1.719 / 342** |

⛔ **Sertakan `pyClassName` pada tiap ekspor.** Dua berkas bernama sama di kelas berbeda adalah **dua
rule**, bukan satu (Koreksi R1).

⛔ **Golongan G3 TIDAK diminta** — nama sama, kelas berbeda, keduanya sah, **keduanya diport**.
Mengekspor ulang tidak menyelesaikan apa pun di sana.

⚠️ **E-Q1 belum masuk paket T-.** Ia tercatat di `_CHECKLIST-KESIAPAN.md` H.3 dan register baris
2228, tetapi tidak punya nomor T. **Kandidat T-10.**

---

# E. Menunggu **Product**

| # | Yang diminta | Memblokir |
| :-: | --- | --- |
| **T-9** | Arti `BusinessCode` (**98** kode) + `BusinessOldId` (**87** kode) — tabel `kode → arti` | klasifikasi produk · layar input |

⛔ **Jangan ditebak.** Sampai terjawab, keduanya **`panic` bila tercapai** (`CLAUDE.md` §3 butir 4).

---

# F. Pekerjaan kita sendiri yang **belum punya spesifikasi**

## F-1 ✅ **SELESAI 25 September 2026** — spec pemuatan ditulis

`04-spec\11-spec-pemuatan-data-lama.md`, lewat `/to-spec` yang **diinvokasi work owner**. Seam
keempat `loader.Flatten` disetujui → **K-074**.

⛔ Terbit di `04-spec\`, bukan `.scratch\` — larangan menulis ke `.scratch\` mengalahkan konvensi
issue tracker; penyimpangan itu tertulis di kepala spec.

*Rumusan asal, dipertahankan sebagai jejak:*

## ~~F-1~~ Program pemuatan data lama — `ROW_UID`

Bagian **tersulit** dari seluruh sisa pekerjaan, dan satu-satunya yang belum berspesifikasi sama
sekali. Bahannya sudah siap di `08-flat\BAHAN-SPEC-PEMUATAN.md` (9 bagian), tetapi **spec-nya belum
ditulis** — ✅ **tidak lagi terhalang keputusan apa pun** sejak K-073.

📌 **Bobot butir ini TURUN, bukan naik.** `ROW_UID` fase 1 bersifat **sementara**: barisnya tetap
butuh identitas, tetapi tidak perlu kekal, karena fase 2 mengosongkan dan memuat ulang seluruh
riwayat dalam satu jalan. Rekonstruksi UID yang V-42 sebut bagian tersulit **tidak lagi diperlukan** —
tinggal pembangkitan biasa.

⛔ **Syaratnya mengikat:** fase 2 wajib selesai **sebelum** sistem dipakai sungguhan. Bila terlewat,
kita turun ke Jalan A dan bobotnya naik lagi. Enam syarat lengkap: **K-072 J-12**.

`[terverifikasi]` Di ekspor Pega **tidak ada satu pun identitas baris yang stabil** —
`pxListSubscript` hanya posisi, `OBJECT_NO` sekadar 1..N. Untuk baris **baru** mudah: bangkitkan UID
saat baris lahir. Untuk baris **lama** itulah masalahnya.

Dua pertanyaan turunan yang harus dijawab sebelum loader ditulis:
1. Seluruh versi atau versi terakhir saja? → **A-1**
2. Baris yang penjodohannya **meragukan** diapakan — diberi UID baru, atau ditandai?

## F-2 `08-flat\02-<slug>.md` — dokumen rancangan versi 115 contoh

⛔ Kuncinya **sudah dibuka**. Ini pekerjaan yang paling wajar dikerjakan berikutnya, dan sekalian
tempat memperbaiki bunyi **V-38** dan **V-40** (lihat **G-1**).

## F-3 Bahan tiket jalur produksi

Empat tiket yang dulu terhalang struktur tabel flat kini **tidak terhalang lagi**: jalur simpan
produksi · `repository` tulis · `services/production` · `models` tulis.

⚠️ `/to-tickets` dan `/to-spec` ber-`disable-model-invocation` — **bahannya disiapkan, lalu
berhenti**.

---

# J. ✅ **SELURUHNYA DIJAWAB 25 September 2026 — K-071**

> **J-1 … J-10 tertutup.** Badan pertanyaannya **tidak dihapus** — ia merekam mengapa pertanyaan itu
> pernah ada dan atas bukti apa (`PANDUAN-KERJA` §7).

| # | Jawaban ringkas |
| :-: | --- |
| **J-1** | `JSON_POLIS` tidak dipakai lagi; yang dipakai **tabel flat baru** |
| **J-2** | **`PRODKE` dan `TGL_INPUT` sama-sama tertinggi** → loader urut `TGL_INPUT`, `PRODKE` jadi pemeriksa silang, **`panic` bila tidak sepakat** |
| **J-3** | `PROD_KE` di tabel flat bertipe **`NUMBER`** → loader mengubah tipe, **`panic`** untuk nilai bukan angka |
| **J-4** | Seluruh data ikut dimigrasikan begitu sistem Go selesai |
| **J-5** | ⛔ `OLD_POLIS_ID` **hanya untuk RENEWAL** — NB dan EDM kosong → **amandemen sebagian K-068** |
| **J-6** | Rekomendasi diminta → ⭐ **Jalan B: muat ulang dari nol**, jangan menjodohkan ke belakang |
| **J-7** | Data lama dimasukkan **sambil menguji** sistem baru |
| **J-8** | Kunci riwayat sama (`ID_PEGA` = `IDPEGA`) → ⛔ **pertanyaan jalur bandingnya BUBAR** |
| **J-9** | Diperiksa → ⛔ `_MENJADI`/`_SELISIH` **tidak dapat** menggantikan nilai sebelum endorsement |
| **J-10** | Yang dimuat **yang sudah jadi** |

## ✅ J-11 … J-14 dijawab 25 September — **K-072**

| # | Jawaban |
| :-: | --- |
| **J-11** | Tautannya **`PRODKE`** → aturan: kelompokkan dengan **`POLICY_NO`** (terisi 115/115), urutkan dengan **`PRODKE`**. ⛔ `POLICY_MASTER_NUMBER` **bukan** kuncinya — terisi **1/115** dan berbeda dari `PolicyNo` |
| **J-12** | Syaratnya **enam butir**, seluruhnya tertulis di K-072 — ⬜ **persetujuannya sendiri masih ditunggu** |
| **J-13** | `STS_KONVERSI` = status **konversi ke tim lain** (hilir), ⛔ **bukan** penanda "sudah jadi" — mengoreksi rencana saya di K-071 |
| **J-14** | `_SELISIH` tanpa `_MENJADI` menandai **penambahan lokasi** → diport apa adanya sebagai kolom tunggal; ⛔ jangan dibuatkan `_MENJADI` baru |

## ✅ J-15 dijawab 25 September — **K-073: Jalan B DISETUJUI**

Fase 2 **mengosongkan tabel** dan memuat seluruh generasi dalam satu jalan. ⛔ Konsekuensi yang
mengikat: **`ROW_UID` fase 1 tidak boleh dijadikan sandaran apa pun di luar tabel flat** — ia akan
berganti.

📌 **V-42 turun bobotnya drastis:** tidak ada lagi UID yang perlu **direkonstruksi**, hanya
**dibangkitkan**.

## ⬜ Dua sisa

| # | Sisa | Sifatnya |
| :-: | --- | --- |
| **J-16** | Penyaring **"endorsement sudah jadi"** untuk fase 1 — `STS_KONVERSI` terbukti bukan jawabannya | ⚠️ tidak menentukan hasil akhir, **tetapi** menyesatkan pengujian |
| **J-17** | Arti `POLICY_MASTER_NUMBER` / `POLICY_MASTER_ID_PEGA` — terisi **1/115**, berbeda dari `PolicyNo` | `belum terverifikasi`, tidak mendesak |

⛔ **Aturan pengelompokan J-11 TIDAK dapat diuji dari korpus.** `[terverifikasi]` 115 berkas memuat
**113 `PolicyNo` unik**; 2 nilai berulang, masing-masing tepat 2× — dan itu **persis** dua pasang
duplikat byte-identik (V-45). **Nol pasangan generasi sejati** ada di contoh, jadi aturannya
bersandar pada jawaban work owner + bentuk skema, bukan pengukuran.

---

## Badan pertanyaan asal *(arsip — jangan dibaca sebagai keadaan sekarang)*

## J-1 ⛔ Apa yang menyatakan "ini polis yang **sama**"?

**Bahasa sederhana:** untuk memilih "yang terakhir", kita harus tahu lebih dulu **mana saja yang satu
rombongan**. Baris mana yang dianggap generasi dari polis yang sama?

`[terverifikasi]` `DDL\JSON_POLIS.txt`:

```
CONSTRAINT "PK_JSON_POLIS"  PRIMARY KEY ("IDPEGA")
CONSTRAINT "JSON_POLIS_U01" UNIQUE ("NOPOLIS", "IDPEGA")
"NOPOLIS"    VARCHAR2(100) NOT NULL
"NOENDORS"   VARCHAR2(100)
"OLDNOPOLIS" VARCHAR2(100)
```

Karena kunci utamanya **`IDPEGA` tunggal** dan yang unik adalah **pasangan** `(NOPOLIS, IDPEGA)`,
satu `NOPOLIS` **boleh** punya banyak baris — tiap generasi ber-`IDPEGA` sendiri. Jadi `[dugaan]`
pengelompokannya lewat **`NOPOLIS`**.

⚠️ **Tapi ada `OLDNOPOLIS`** — kolom itu ada justru karena nomor polis **bisa berganti**
(`[dugaan]` saat renewal). Kalau begitu, dua baris dengan `NOPOLIS` berbeda bisa jadi **polis yang
sama juga**.

**Pertanyaannya:** rantai generasi dijalin lewat **`NOPOLIS` saja**, atau **`NOPOLIS` +
`OLDNOPOLIS`** — dan apakah pergantian nomor polis (renewal) dihitung sebagai **generasi baru dari
polis lama**, atau **polis baru**?

## J-2 ⛔ "Versi terakhir" ditentukan oleh apa — `PRODKE` tertinggi atau `TGL_INPUT` terbaru?

⛔ **Ini bukan pertanyaan sepele.** `[terverifikasi]` `DDL\JSON_POLIS.txt` baris 7:

```
"PRODKE" VARCHAR2(5)
```

**`PRODKE` bertipe TEKS, bukan angka.** Mengurutkan teks membuat **`'9'` lebih besar daripada
`'10'`** — persis jebakan yang sama dengan `FACINOFFER.RATE` (`CLAUDE.md` §4.1). Kalau loader
memilih "PRODKE tertinggi" dengan urutan teks, polis yang sudah lewat endorsement ke-10 akan
**memuat generasi yang salah** — dan salahnya **diam-diam**.

📌 **Sistem lama sendiri tidak mengurutkan dengan `PRODKE`.** `[terverifikasi]` (dicatat pada
jawaban P-10) rule `GetProdKeOldData_SQL` memakai **`ORDER BY TGL_INPUT DESC`**.

⚠️ `[terverifikasi]` `PRODKE` **nol kemunculan** pada 115 contoh `DDL\CONTOH\` — ia hanya hidup
sebagai **kolom Oracle**, tidak ikut di dalam berkas kasus. Artinya **korpus tidak dapat
membuktikan** apakah ada polis yang pernah mencapai dua digit → **T-10**.

```powershell
# menghasilkan: berkas dipindai = 115, nol tag PRODKE, nol nilai dua digit
$n=0; $h=0
foreach ($p in [IO.Directory]::EnumerateFiles('D:\migrasi\RNM\DDL\CONTOH','*.*')) {
  $n++; $t=[IO.File]::ReadAllText($p)
  $h += ([regex]::Matches($t,'(?i)<(ProdKe|PRODKE)>[^<]*</\1>')).Count }
"berkas=$n tagProdKe=$h"
```

## J-3 Nilai `PRODKE` disimpan **apa adanya** atau **dinomori ulang**?

**Bahasa sederhana:** sebuah polis sudah sampai endorsement ke-7. Di fase 1 hanya generasi ke-7 yang
dimuat. Kolom `PROD_KE` di tabel flat diisi **7**, atau **1**?

| Pilihan | Akibatnya |
| --- | --- |
| **Apa adanya (7)** | Jujur, tetapi tabel punya **lubang** — ada generasi 7 tanpa 1…6, sampai fase 2 mengisinya |
| **Dinomori ulang (1)** | Tabel rapi, tetapi ⛔ **menulis ulang sejarah** — dan fase 2 nanti harus membongkarnya kembali |

`[dugaan]` **apa adanya** lebih sejalan dengan `CLAUDE.md` §1, tetapi ini keputusan Anda.

## J-4 Endorsement **baru** di sistem baru mulai dari nomor berapa?

Kalau yang dimuat generasi ke-7, endorsement berikutnya harus jadi **8** — bukan **1**. Berarti
penghasil nomornya **wajib membaca nilai yang dimuat**, bukan mulai dari nol.

⚠️ `[terverifikasi]` ada tabel `C_COUNTER_PRODKE` di `DDL\`. **`belum terverifikasi`** apakah ia
penghasil nomor itu dan bagaimana ia diselaraskan sesudah pemuatan fase 1.

---

## J-5 `OLD_POLIS_ID` **kosong seluruhnya** di fase 1 — dikonfirmasi?

Tidak ada generasi sebelumnya untuk ditunjuk, jadi kolom K-068 itu `NULL` di **100 %** baris fase 1.
`[terverifikasi]` di Oracle `UNIQUE` mengizinkan **banyak** `NULL`, jadi kendala uniknya **tidak
dilanggar**. Kolomnya baru bermakna setelah fase 2.

## J-6 ⛔ Fase 2 harus **menjodohkan ke belakang** — diterima?

**Bahasa sederhana:** di fase 1, tiap baris dapat `ROW_UID` baru yang kita bangkitkan sendiri. Nanti
di fase 2, saat generasi lama dimuat, baris lama **tidak boleh** dapat UID sembarangan — ia harus
dicocokkan ke UID baris keturunannya yang **sudah terlanjur ada**.

⛔ Ini **lebih sulit** daripada memuat semuanya sekaligus, karena arah penjodohan terbalik dan
sasarannya sudah beku.

**Pertanyaannya:** apakah fase 1 perlu **menyimpan sesuatu sekarang** untuk memudahkan fase 2 —
misalnya menyimpan jejak posisi asli tiap baris — atau cukup mengandalkan `SEQ_NO`, `OBJECT_NO` dan
`SRC_PATH` yang sudah ada di rancangan?

## J-7 Rekonsiliasi paralel run hanya mencakup **generasi terakhir** — dikonfirmasi?

ADR-0001 menuntut **nol selisih sampai digit terakhir**. Di fase 1 hanya generasi terakhir yang ada,
jadi yang dapat dibandingkan juga hanya itu. ⛔ **ADR-0001 tidak dilonggarkan** — hanya populasinya
yang lebih sempit. Perlu ditegaskan agar tidak terbaca sebagai pelanggaran.

---

## J-8 ⚠️ Riwayat akseptasi **utuh**, tetapi polisnya **sebagian** — apakah aman?

**Bahasa sederhana:** tabel `HISTORYAKSEPTASIPEGA` dibaca untuk **mengambil keputusan** — jalur
banding dan flag reject (`CLAUDE.md` §4.3), bukan sekadar catatan. `[terverifikasi]` **K-028**: tabel
itu **tidak pernah dipangkas**, jadi riwayatnya lengkap untuk **semua** generasi.

⛔ Di fase 1 muncul **ketimpangan**: riwayatnya lengkap, polisnya tidak. Baris riwayat milik generasi
lama akan menunjuk `ID_PEGA` yang **tidak ada** di tabel polis.

`[terverifikasi]` kolom tabel itu: `ID_PEGA`, `TGL_TRANSFER`, `STATUS`, `USERNAME`, `WORKBASKET`,
`ID_KOMITE`, `OPERATORID`.

**Pertanyaannya:** apakah hitungan riwayat berstatus reject — yang menentukan tersedianya jalur
banding — dihitung atas **seluruh generasi** (termasuk yang belum dimuat), atau **hanya generasi yang
ada**? Dua jawaban itu memberi jalur banding yang berbeda untuk kasus yang sama.

## J-9 Nilai **sebelum endorsement** memang tidak ada di fase 1 — dikonfirmasi?

Cabang `OldData` **tetap dibuang** (V-19), dan ⛔ ia **bukan** pengganti: `[terverifikasi]` V-43a
menyatakan `OldData` adalah potret salinan kerja **sesudah** baris baru dibentuk, bukan versi
sebelumnya. Jadi di fase 1 nilai sebelum endorsement memang **tidak tersedia sama sekali** sampai
fase 2 berjalan.

**Pertanyaan turunan:** apakah pasangan kolom **`_MENJADI` / `_SELISIH`** di tabel produksi
(`CLAUDE.md` §4.3, dipertahankan berpasangan) dapat memikul sebagian kebutuhan itu selama fase 1?
`[pertanyaan terbuka]` — belum diperiksa.

## J-10 Endorsement yang **belum selesai** — ikut dimuat sebagai "terakhir"?

Kalau generasi paling akhir sebuah polis ternyata endorsement yang **masih dalam proses**, apakah
itu yang dimuat, atau versi terakhir yang **sudah jadi**?

`[terverifikasi]` `JSON_POLIS` punya kolom `STS_KONVERSI NUMBER` dan `TGL_KONVERSI`.
**`belum terverifikasi`** arti nilai `STS_KONVERSI` — jadi **jangan ditebak** sebagai penyaring
sebelum artinya dipastikan.

---

# G. Utang kerapian dokumen — bukan blocker, tetapi menyesatkan bila dibiarkan

## G-1 ⚠️ Bunyi **V-38** dan **V-40** masih menyebut pencarian `PROD_KE` tertinggi

Keduanya digantikan **`OLD_POLIS_ID`** (**K-068**). Dua keputusan yang sama-sama sah akan terus
**terbaca bertentangan** sampai bunyinya diperbaiki. Pola sama dengan V-40/V-40b setelah butir 3.

## G-2 `08-flat\BAHAN-SPEC-PEMUATAN.md` sudah basi di tiga tempat

| Bagian | Yang salah sekarang |
| --- | --- |
| §3 tabel kolom sistem | Masih mencantumkan baris `LINI` — **K-069** menyatakan penanda lini **tidak ada di proyek ini** dan kolomnya sudah dihapus |
| §6 *Yang masih terbuka* | Kelima barisnya (**8 · 5b · 7b · 1c · 10**) **semuanya sudah tertutup** K-067/K-068/K-069 |
| §9 butir 4 & 5 | Butir 4 sudah gugur lewat K-065; butir 5 kini = **A-1** |

## G-3 `_LANJUTKAN-DARI-SINI.md` bertentangan dengan dirinya sendiri

§4 masih menandai butir **5b · 7b · 8** sebagai ⬜ **terbuka**, sementara §8 di berkas yang sama
menyatakan **seluruhnya tertutup**. §9 juga masih menulis register sebagai **K-001…K-062**, padahal
sudah **K-001…K-069**.

## G-4 Angka **"95 rule berbeda versi"** tidak tereproduksi

Tertulis di `CLAUDE.md` §4.5 dan `steering\PANDUAN-KERJA.md` §5; ukur ulang menghasilkan
**87 / 85 / 0**. ⛔ `CLAUDE.md` berada **di luar** `OUTPUT\` — menyuntingnya bukan wewenang saya.
Diangkat di sini supaya keputusannya milik Anda. Menyatu dengan **A-2**.

## G-5 Selisih **1.202** simpul `rowdata` belum terekonsiliasi

Dua cara menghitung memberi **22.139** dan **20.937**, padahal hanya ada **satu** ejaan `rowdata`.
Angka yang dipakai (**20.937 + 1.685 = 22.622**) adalah yang **konsisten** secara internal.

⛔ **Memadai untuk menentukan ukuran pekerjaan; JANGAN dipakai sebagai angka rekonsiliasi.**

## G-6 Dua sisa dari `_LANJUTKAN` §6

| | Butir | Keadaan sekarang |
| :-: | --- | --- |
| a | `Diagram-Skema-Tabel-NusantaraRe.xlsx` tidak ada lagi di bawah `RNM\` | **K-065**: diagram tidak dipakai, jadi ini **bukan blocker**. ⚠️ Yang tersisa: butir **1, 5b, 7** sudah tertutup di atas sumber yang **tidak dapat diaudit ulang** — risiko yang diterima sadar |
| b | Minta ekspor ulang `NB-176005` dan `NB-184183` | Keduanya **byte-identik** dengan `NB-172576` dan `NB-181622`. Bila ekspor ulang berbeda → cacat **pengumpulan contoh**, dan tafsir **V-45** ikut gugur |

---

# H. Yang akan `panic` bila tercapai — konsekuensi, bukan tugas

Ini **bukan** daftar yang harus dijawab sekarang. Ia dicatat supaya tidak ada yang terkejut saat
aplikasi berhenti — dan berhenti memang yang **diinginkan**: aplikasi yang jalan tapi salah
diam-diam jauh lebih berbahaya (`CLAUDE.md` §4.5).

| # | Yang `panic` | Kapan dicabut |
| :-: | --- | --- |
| 1 | `IsPKSASM` | ⛔ daftar nilai `IsB2B` **tidak** mencabutnya. Sebabnya tiga hal lain: `<pyConditionString>` masih placeholder · label **belum ter-resolve** · `<pyTempText>` = `true`, di **ketiga** folder. Hanya dicabut oleh ekspor produksi yang labelnya sudah ter-resolve |
| 2 | `BusinessCode` · `BusinessOldId` | **T-9** |
| 3 | Arti enumerasi tipe email `1` / `2` / `7` pada jalur retro | **K-061** — tidak dijelaskan korpus |
| 4 | `TreatyName = "SPL"` | belum ada pemetaan angka→label; `M_TREATY_IN` menyimpan labelnya **di dalam CLOB JSON** |
| 5 | Bentuk **tabel tarif Travel** | `ViewPremiTravel_SQL` **NIHIL di korpus** — masukan dari repository |
| 6 | Rule yang dirujuk tetapi berkasnya **tidak ada di folder mana pun** | ekspor produksi tunggal (**T-8**) |
| 7 | Status jalur produksi **Bonding** (`SaveFacinProdEDMBonding_Act`) | ⛔ cabang `IsBonding` belum dikonfirmasi aktif. **Tidak dihapus** (K-006) |

⚠️ **Dua celah lama yang tetap berdiri, dicatat agar tidak dianggap terlewat:**

- **Celah K-004** — rule asli kelas `Work` (`serviceInsertArasapas_act` tanpa sufiks) **tetap tidak
  ada** di korpus mana pun. K-035 memutuskan **arah rancangan**, bukan menemukan rule yang hilang.
- **`SetOldData`** — `[pertanyaan terbuka]` apakah ia proses salin yang dimaksud **K-039**, atau
  mekanisme terpisah. `GetPolicyData_ACT` yang memanggilnya **belum ditelusuri**.
- **Guard berbasis identitas** — apakah ada rule yang **membaca** kolom `LOGIN` untuk mengambil
  keputusan: **belum diperiksa**. Ketiga query tabel limit memilih `NAMA AS CARI5`; `[terverifikasi]`
  `CARI5` adalah slot serbaguna yang disebut **63 Activity** di NB, jadi kehadirannya **bukan bukti**
  ia menggerakkan alur. ⛔ Apa pun jawabannya, **K-025 tetap penuh**: nilai `NAMA` tidak pernah
  disalin ke artefak mana pun — hanya **jumlah dan mekanismenya**.

---

## Urutan yang saya sarankan

1. ~~**A-1**~~ ✅ **terjawab 25 September — K-070.** Penggantinya: **J-1 … J-4**, empat pertanyaan
   yang kini memegang kunci program pemuatan.
2. **Kirim golongan C, D, E** — `_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md` sudah siap kirim; makin cepat
   dikirim makin cepat kembali, dan pengirimannya tidak menghalangi pekerjaan lain. ⭐ **T-10 baru**
   menyertainya, dan ia menjawab separuh **J-2**.
3. **Golongan B** — paling terikat waktu. Begitu Pega mati, jawabannya hilang selamanya.
4. **F-2** sambil menunggu — menulis `08-flat\02-` tidak bergantung pada satu pun jawaban di atas,
   dan sekalian membereskan **G-1**.
5. **G-2, G-3** — suntingan kecil, tetapi tiap hari dibiarkan adalah hari dokumen saling
   bertentangan.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
