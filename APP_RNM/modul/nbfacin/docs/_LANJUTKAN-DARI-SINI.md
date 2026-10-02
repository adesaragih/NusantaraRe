# LANJUTKAN DARI SINI — status kerja migrasi Facultative Inward

> **Disimpan 2026-09-24 sore WIB.** Berkas ini **menggantikan** versi 19 September, yang berhenti di
> register K-043 dan Prompt 52 dan belum memuat Fac Out maupun pekerjaan tabel flat.
>
> Versi lama **tidak diarsipkan terpisah** — isinya yang masih berlaku (aturan kerja §2, kontrak
> metode §3, jebakan §10) dibawa utuh ke sini; yang basi diperbarui di tempatnya.
>
> Untuk **memulai sesi baru**, berkas inilah pintu masuknya. Baca §7 lebih dulu, lalu **ukur ulang
> sendiri** angka-angka di sana sebelum mempercayainya — dokumen bisa tertinggal, korpus tidak.

---

## 0. Satu layar — di mana kita

| | |
| --- | --- |
| **Pekerjaan berhenti di** | `OUTPUT\08-flat\` — rekonsiliasi rancangan tabel flat |
| **Menunggu** | keputusan work owner atas **sembilan titik bentrok** |
| **Boleh dikerjakan tanpa menunggu** | tidak ada di jalur flat; `08-flat\02-…` **terkunci** |
| **Butir terbuka** | ✅ **TIDAK ADA di jalur tabel flat** — seluruhnya ditutup 24-25 September. ⛔ **Issue di luar jalur flat ada dan masih banyak** → `_DAFTAR-ISSUE-TERBUKA.md` (A-1…A-7 · U-1…U-7 · T-1…T-9 · E-Q1) |
| **Butir 9 (mata uang)** | ✅ **DITUTUP 24-09-2026** — **K-063** + **ADR-0007** |
| **Butir 1b (tabel akar)** | ✅ **DITUTUP 24-09-2026** — **K-064**: tabel **SAMA** dengan Treaty In |
| **Butir 8 (tabel Total*)** | ✅ **DITUTUP 24-09-2026** — **K-067**: informasi saja, aman dihapus. ADR-0001 **tidak diubah** |
| **Butir 5b (tautan versi)** | ✅ **DITUTUP 24-09-2026** — **K-068**: `OLD_POLIS_ID`, UNIK, tanpa FK |
| **`K-` berikutnya** | **K-075** — K-070…K-073 migrasi dua fase · **K-074 seam ke-4 `loader.Flatten`** (25-09) |
| **ADR berikutnya** | **0008** |
| **Token log terakhir** | **Prompt 78**, kumulatif lihat `_TOKEN-LOG.md` |

---

## 1. Riwayat program — apa yang sudah terjadi

`[terverifikasi]` Diturunkan dari `_TOKEN-LOG.md` dan isi folder `OUTPUT\`.

| Tahap | Prompt | Hasil yang tersimpan |
| --- | --- | --- |
| Orientasi + discovery awal | 1–3 | Pass korpus NB; `_ARSIP-lintas-siklus\` (kemudian diarsipkan) |
| Keputusan lingkup + kuesioner | 4–15 | **K-001…K-021**; 4 kuesioner; ADR-0001…0006 |
| DDL masuk + enumerasi ditutup | 16–22 | **K-023…K-029**; enumerasi tertutup K-029 |
| **Tiket NB** | 23–27 | **16 tiket** `05-tickets\` + indeks |
| Pergeseran ke RNW + discovery | 28–33 | **K-030…K-033**; NB↔RNW nol beda |
| Spec + **tiket RNW** | 34–42 | **8 tiket** `05-tickets\rnw\` + indeks |
| **Discovery EDM** E-1…E-6 | 43–52 | 7 dokumen `07-edm\`; kontrak 23 tag (K-042/K-043) |
| Bahan spec EDM + **tiket EDM** | 53–64 | 10 berkas `04-spec\`; **22 tiket** `05-tickets\edm\` + indeks |
| **Fac Out** — discovery → tiket | 65–70 | **14 tiket** `05-tickets\facout\` + indeks; **K-057…K-062**; audit 01–10 |
| Ralat prorata + paket pihak luar | 71 | Catatan pelengkap K-060; `10-audit\11` daftar periksa UI Pega; paket **T-1…T-9** |
| **Rekonsiliasi tabel flat** (Tugas 1) | 72 | `08-flat\00-rekonsiliasi-rancangan-flat.md` |
| **Verifikasi independen** Tugas 1 | 73 | `08-flat\00a-verifikasi-independen-tugas-1.md` |

**Hitungan sekarang** `[terverifikasi]`: `OUTPUT\` **148** berkas `.md` + **1** `.xlsx` · `DDL\` **130** berkas +
`DDL\CONTOH\` **115** · `adr\` **6** · `04-spec\` **11** · `10-audit\` **11** ·
tiket **16 + 8 + 22 + 14** (masing-masing plus satu indeks).

---

## 2. Aturan kerja yang MENGIKAT

- Bahan kerja **HANYA** dari `D:\migrasi\RNM\`. Korpus Treaty `RNM_BRD\` **tidak dibaca** (K-005).
- `NB FacIn\`, `RNW Fac In\`, `Endorsment Fac In\` = **READ-ONLY mutlak**.
- **SELURUH output ke `D:\migrasi\RNM\OUTPUT\`**. ⛔ **JANGAN tulis ke `.scratch\`**.
- **DILARANG HALUSINASI**: tiap klaim sebut path + nama rule + label
  `[terverifikasi]` / `[dugaan]` / `[pertanyaan terbuka]`. Tiap angka sertakan perintah audit.
- **NAMA BUKAN BUKTI** — periksa `pxObjClass` / tag pembawa, bukan nama berkas atau folder.
- **Identitas rule = basis `pzInsKey` + `pyClassName`**, bukan nama berkas (R3 / Koreksi R1 di
  `_ARSIP-lintas-siklus\_BACA-INI.md`). Stempel waktu di ujung `pzInsKey` dibuang dulu.
- **Nama orang / email / nilai DOB / plat / alamat TIDAK PERNAH disalin** (K-025) — **nama kolom
  tetap ada**. ⛔ `DDL\CONTOH\` memuat **data pribadi nyata**; skrip yang menyapunya hanya boleh
  mengeluarkan hitungan dan nama tag.
- **Jangan menulis berkas `.go`** sampai diminta.
- **GUARD SKILL**: `/to-spec`, `/to-tickets`, `/to-questionnaire` ber-`disable-model-invocation`.
  Siapkan bahannya lalu **BERHENTI**. Jangan improvisasi penggantinya.
- **Koreksi hanya sah bila yang LAMA dibuktikan tidak didukung korpus** (R1). Menemukan alternatif
  yang juga ada **bukan** bukti yang lama salah.
- **Sensus label dijalankan PALING AKHIR**, setelah semua suntingan (R2).
- Keputusan yang dibalik **tidak dihapus** — ditandai dibatalkan dan menunjuk penggantinya.
- Uang tidak pernah `float`. Kode usang **diport apa adanya** (`CLAUDE.md` §1).
- Skrip temp di `OUTPUT\_tmp_*.ps1`, tulis hasil ke berkas lalu baca, **HAPUS setelah selesai**.
- Catat START/END + estimasi token di `OUTPUT\_TOKEN-LOG.md` tiap tahap.

---

## 3. Kontrak metode 23 tag — SALIN PERSIS

Dasar K-042 dan K-043. Rinci di `07-edm\02-e1-pola-perbedaan-when.md` §1.

**23 tag dibuang seluruh barisnya:**

```
pxCreateDateTime  pxUpdateDateTime  pxSaveDateTime  pxCommitDateTime  pxMoveImportDateTime
pxOriginalCreateDateTime  pyRuleFormStatusTime  pxWarningCreatedTime  pxCreateOperator
pxCreateOpName  pxCreateSystemID  pxUpdateOperator  pxUpdateOpName  pxUpdateSystemID
pxMoveImportOperId  pxMoveImportOperName  pxOriginalCreateOperator  pxOriginalCreateOpName
pxOriginalCreateSystemID  pxHostId  pzChecksum  pzIndexCount  pyShowJavaWindowName
```

**6 tag berkunci — tag DIPERTAHANKAN, stempel waktu di dalamnya diganti `<TS>`:**

```
pzInsKey  pzIndexOwnerKey  pzOriginalInstanceKey  pzDocumentKey  pyJavaClassName  pxInsName
```

**+ isi blok `pzIndexes` diabaikan** (amandemen K-043). ⚠️ Tag ini **selalu ber-atribut** — cari
dengan `<pzIndexes` **tanpa** `>`.

⛔ **JANGAN buang** tag pembawa kelas: `pxRuleClassName` · `pxRuleFamilyName` · `pxRuleObjClass` ·
`pxObjClass` · `pyDesignatedClass` · `pyClassName` · `pyRuleName` · `rowdata`.

**Wajib** urutkan `[StringComparer]::Ordinal` sebelum hash.

**Uji silang wajib:** `When\FlagOldData` harus **BEDA**; `When\IsAdmin` · `IsBonding` ·
`IsAddButton` · `DataPage\D_AnekaList` harus **identik**.

---

## 4. ⛔ PEKERJAAN AKTIF — rekonsiliasi tabel flat

Dua rancangan bersaing menggantikan CLOB `JSON_POLIS.DATA_JSON`:

| | Berkas | Dasar |
| --- | --- | --- |
| **lama** | `08-flat\01-rancangan-tabel-flat-dari-contoh-nb.md` | **5** contoh · nama `FLAT_*` |
| **baru** | `Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx` | **114** contoh · 75 tabel · 1.290 kolom · V-1…V-50 |
| **rekonsiliasi** | `08-flat\00-rekonsiliasi-rancangan-flat.md` | sembilan titik bentrok |
| **verifikasi** | `08-flat\00a-verifikasi-independen-tugas-1.md` | 25/27 angka tereproduksi; **dua butir berubah** |
| **lintas siklus** | `08-flat\Tabel-Flat-Lintas-Siklus.xlsx` | **11 lembar** · **78 tabel · 1.310 kolom** · diukur ulang dari **115** contoh (24-09) · lembar `Audit Mata Uang` · kolom **`NULL`** |

⛔ **`Tabel-Flat-Lintas-Siklus.xlsx` BUKAN pengganti rancangan dan BUKAN `02-`.** Ia tidak memutuskan
satu butir pun: isi rancangan baru dibawa apa adanya, yang ditambahkan hanya **dimensi siklus yang
diukur** plus lembar `Butir Terbuka` yang mencatat kelima butir sebagai **masih terbuka**. Kunci
`08-flat\02-…` tetap berlaku.

⛔ **Dokumen lama MASIH BERLAKU** sampai kesembilan butir diputuskan. `08-flat\02-…` **terkunci**.

### Keadaan kesembilan butir sesudah verifikasi

| # | Butir | Status |
| :-: | --- | --- |
| 1 | Nama `FLAT_*` vs `T_*` | ✅ **baru menang** — skema rumah 43 `T_*`, nol `FLAT_*`. Lanjutannya masih terbuka: tabel fisik bersama Treaty In atau terpisah |
| 2 | Kunci alami vs V-47 FK | ✅ baru menang, bersyarat — 10 kolom bukti bercampur tetap tidak disentuh |
| 3 | `pxListSubscript` vs `SEQ_NO` | ✅ **TIDAK lagi memblokir** — pertentangan V-41/V-40b **semu**; 18.140 baris, nol menyimpang. Adopsi `SEQ_NO`, **perbaiki bunyi V-40/V-40b** |
| 4 | `*Old` disimpan atau tidak | ✅ kaidah baru menang (9 tag, nol di NB) · ⚠️ risiko: nilai lama hilang bila versi sebelumnya tidak ikut dimigrasikan |
| 5 | Kunci versi | ✅ **baru menang** — `JSON_POLIS` ber-`PRIMARY KEY (IDPEGA)` tunggal |
| 5b | **`OLD_POLIS_ID`?** | ⬜ **terbuka** — Treaty In memakainya; preseden `JSON_POLIS.OLDNOPOLIS` sudah ada |
| 6 | `PRODKE` tipe | ✅ baru menang (`NUMBER(5)`) — sebagian gugur bila 5b diadopsi |
| 7 | ViewSuggest | ✅ **DITUTUP mengikuti V-31** — pembukaan ulang **salah baca**; `T_VIEW_SUGGEST` tabel **klaim**, dan Treaty In memakai `HISTORYAKSEPTASIPRODUCTION` (baris relasi 58) |
| 7b | **`IsCedingConfirm`** | ⬜ **terbuka, menggantikan butir 7** — jangan ditulis ke `POSISI`; itu bentrok dengan V-48 |
| 8 | Tabel `Total*` dibuang | ⬜ **terbuka** — V-28 minta toleransi, **ADR-0001 menolaknya**. Empat jalan A/B/C/**D** |
| 9 | `*_CCY` per kolom uang | ✅ **DITUTUP 24-09-2026 — K-063 + ADR-0007.** Empat kelompok: **(a) AMAN** 8 tabel/25 kolom · **(b) PULIHKAN** 10/61 · **(c) WARISI** 8/33 · **(d) TANPA mata uang** 5 total skalar akar. Sembilan kolom terbukti **bukan uang** dan keluar → **33 tabel / 147 kolom** |

### Angka kunci tabel flat `[terverifikasi]`

| Ukuran | Nilai |
| --- | ---: |
| Contoh di `DDL\CONTOH\` | **115** = **102 NB** + 13 EDM + **0 RNW** |
| — ⛔ **populasi efektif** (2 pasang duplikat byte-identik) | **113** |
| Tabel · kolom di rancangan baru | **78** · **1.329** |
| Tabel tanpa kolom mata uang | **65** |
| `<ID>` akar: angka murni / kelas+nomor / tidak ada | **38 / 50 / 26** |

⛔ **Setiap statistik berbentuk "N dari 115" menghitung dua kasus dua kali** — termasuk cakupan
kolom V-48 dan V-49a.

---

## 5. Menunggu pihak luar — `_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md`

| # | Butir | Pemilik |
| :-: | --- | --- |
| T-1 | Isi `M_LINK_SERVICE` — **hanya** `URL`, `KATEGORI_1`, `KATEGORI_2` | DBA |
| T-2 | Dua pertanyaan `HISTORYAKSEPTASIPEGA` (pemangkasan · pengisi `ID_KOMITE`) | DBA |
| T-3 | Status 8 nama yang tiba sebagai `CREATE TABLE`, bukan kode prosedur | DBA |
| T-4 | Isi baris `OPENPROTEKSI_EDM` | DBA |
| T-5 | **10** `DecisionTable` diekspor ulang | IT |
| T-6 | **10** rule golongan **G2** diekspor ulang | IT |
| T-7 | Apakah `CountRateRetroCov` kelas `Data-Cargo` pernah ada | IT / admin Pega |
| T-8 | Satu ekspor **rule** produksi pada satu titik waktu | IT |
| T-9 | Arti `BusinessCode` (98) + `BusinessOldId` (87) | Product |

⛔ **`ALL_SOURCE` BUKAN lagi butir terbuka** — 33 nama sudah masuk `DDL\`. Dokumen lain yang masih
menyebutnya sebagai penunggu jalur produksi **sudah basi**.

📌 Pertanyaan yang hanya bisa dijawab di layar Pega: `10-audit\11-daftar-periksa-ui-pega.md`
(U-1 kode transisi 5 · U-2 kode 1 dan 4 · U-3 rule kembar `IsFire`/`IsPA`/`IsUW` · U-4 `Data-Cargo`).

---

## 6. ⬜ Dua tindakan tertunda, di luar keputusan

**[a] Kembalikan `Diagram-Skema-Tabel-NusantaraRe.xlsx` ke `DDL\`.**
`[terverifikasi]` Berkas itu **tidak ada lagi** di bawah `D:\migrasi\RNM\` — dihapus ke Recycle Bin
24-09-2026 15.31. Isinya tidak hilang: salinan ber-SHA256 identik ada di Recycle Bin dan di
`C:\Users\Administrator\Downloads\RNM\DDL\`. ⛔ Selama ia tidak kembali, **butir 1, 5b dan 7
bersandar pada sumber yang tidak dapat diaudit ulang** (`CLAUDE.md` §3.1).

**[b] Minta ekspor ulang `NB-176005` dan `NB-184183`.**
`[terverifikasi]` Keduanya **byte-identik** dengan `NB-172576` dan `NB-181622`. Bila hasil ekspor
ulang berbeda, ini cacat pengumpulan contoh — **bukan** cacat data Pega, dan tafsir V-45 tentang
"penyalinan data" ikut gugur.

---

## 7. Gerbang penomoran + angka resmi

`[terverifikasi]` **Verifikasi sendiri, jangan percaya angka ini begitu saja.**

| | Nilai |
| --- | --- |
| `K-` berikutnya | **K-075** — **72** judul, rentang K-001…K-074 |
| K-017 · K-022 | **memang tanpa judul `##`** — K-017 kosong permanen, K-022 dicadangkan (register baris 1047–1050). **Bukan cacat** |
| ADR berikutnya | **0008** — ADR-0007 ditulis 24-09-2026 |
| Token log terakhir | **Prompt 78** |
| **Contoh `DDL\CONTOH\`** | **115** = **102 NB + 13 EDM + 0 RNW** · populasi efektif **113** (2 pasang byte-identik) |
| Korpus | NB **2.083** · RNW **1.927** · EDM **2.061** · total **6.071** |
| Rule `When` | **601**; **170** menyembunyikan kondisi; **0** tidak terbaca |
| Lingkup EDM (K-043 diamandemen) | **708** = 354 berbeda + 354 EDM-only |
| Rule kembar seluruh korpus | G1 **4.590** · G2 **483** · G3 **215** (109 nama unik) |

⚠️ **Dua angka 354 itu KEBETULAN dan TIDAK berhubungan.**

---

## 8. LANGKAH BERIKUTNYA

**Menunggu keputusan work owner** — terurut menurut dampak:

## ✅ SELURUH BUTIR RANCANGAN TABEL FLAT SUDAH DITUTUP

| Butir | Ditutup oleh | Hasil |
| --- | --- | --- |
| 1 · 2 · 3 · 4 · 5 · 6 · 7 | rekonsiliasi 24-09 | — |
| **9** mata uang | **K-063** + **ADR-0007** | empat kelompok (a)(b)(c)(d) |
| **1b** tabel akar | **K-064** | SAMA dengan Treaty In |
| **1d** rekonsiliasi kolom | **K-065** | gugur, diagram tidak dipakai |
| JSON=XML · RNW=NB | **K-066** | sepadan seluruh COB |
| **8** tabel `Total*` | **K-067** | informasi saja, ADR-0001 **tidak diubah** |
| **5b** tautan versi | **K-068** | `OLD_POLIS_ID`, UNIK, tanpa FK |
| **7b · 9b · 1c · usulan 10** | **K-069** | lihat di bawah |
| fan-out mata uang | diukur | 432 baris, seluruhnya satu kode |

### K-069 — empat yang terakhir

- **7b** `IsCedingConfirm` → **kolom sendiri**. Bukan `POSISI`, yang sudah dipakai V-48.
- **9b** empat kolom kembaran → **dikeluarkan**. `T_LOCATIONLIST` keluar dari (a);
  `T_PERSONLIST` keluar dari (b). **33/147 → 31/143**.
- **1c** kolom penanda lini → ⛔ **TIDAK ADA di proyek ini**. Kolom `LINI` yang sempat ditambahkan
  24-09 **dihapus kembali**. K-064 tetap berlaku; yang ditetapkan adalah **lingkup**.
- **usulan 10** → **disetujui**. **24 kolom** mata uang jadi `NOT NULL DEFAULT 'UNKNOWN'`.
  ⛔ Kelompok (c) menuntut **8 kolom BARU** — tanpa kolomnya, "bersuara" mustahil.

⚠️ **Bunyi V-38 dan V-40 masih perlu ditulis ulang** — keduanya menyebut pencarian `PROD_KE`
tertinggi, yang digantikan `OLD_POLIS_ID` (K-068). Dua keputusan sah yang akan terus terbaca
bertentangan sampai bunyinya diperbaiki.

### ✅ Butir 8 — cara tertutupnya penting diingat

**K-067:** penjumlahan di ketiga tabel `Total*` **informasi saja**, tidak dipakai untuk pembayaran.
**Aman dihapus.**

⛔ **Pertentangan dengan ADR-0001 BUBAR, dan tidak satu pun jalan A/B/C/D dipakai.** Keempatnya
berangkat dari premis bahwa nilainya akan **dihitung ulang lalu dibandingkan**. Premis itu gugur:
nilainya tidak disimpan dan tidak dihitung ulang, jadi **tidak ada yang dibandingkan**. Selisih
2,6 per sepuluh juta **tidak pernah terwujud**. **ADR-0001 tetap utuh, tanpa amandemen.**

⚠️ **Dua pagar yang tetap berlaku:**
- Kelima skalar akar di `T_GENERAL_POLIS` **tidak dihitung ulang** — nilai terekam (ADR-0007).
- Angka total hasil hitung ulang boleh untuk **tampilan**, **tidak boleh** dipakai dalam
  **rekonsiliasi**.

`[terverifikasi]` Rancangan sudah sesuai: ketiga tabel `Total*` **tidak ada** di 78 tabel;
`T_FR_TOTALTSIPREMIRETRO` tetap ada (dikecualikan V-28); kelima skalar akar tetap ada.
**Tidak ada yang perlu disunting.**

✅ **Butir 1b DITUTUP** 24-09-2026 — **K-064**: `T_WORK_POLIS` dan `T_GENERAL_POLIS` adalah **tabel
yang SAMA** dengan Treaty In dan PremiumList, dibedakan **kolom penanda lini**. Kolom khas satu lini
**wajib nullable**. Dua butir baru lahir darinya:

- **1c** ⬜ nama dan nilai sah kolom penanda lini. `[dugaan]` `LINI`, mengikuti preseden
  `T_WORK_CLAIM.LINI`. Kolomnya **sudah ditambahkan** ke `T_WORK_POLIS`, ditandai menunggu konfirmasi.
- **1d** ⬜ rekonsiliasi **18 kolom** `T_WORK_POLIS` dan **74 kolom** `T_GENERAL_POLIS` sisi Fac In
  terhadap sisi Treaty In. ⛔ **Terhalang berkas hilang.**

**Butir 9b** — ⬜ **terbuka**, lahir dari K-063: apakah koreksi "bukan uang" juga berlaku bagi
**kembaran di sisi Fac In**. `[terverifikasi]` diukur tanpa `OldData` atas 115 contoh:

| Kolom | Nilai | Isi |
| --- | ---: | --- |
| `T_LOCATIONLIST.LOSS_RATIO1_YEAR_AMOUNT` | **194** | semuanya 0 |
| `T_LOCATIONLIST.LOSS_RATIO35_YEAR_AMOUNT` | **194** | semuanya 0 |
| `T_LISTINSTALLMENT.STAMP` | **160** | semuanya 0 |
| `T_PERSONLIST.ASMCC_AMOUNT` | **5** | 4 bernilai `1` + 1 besar |

⛔ Bila ya, `T_LOCATIONLIST` **keluar dari kelompok (a)**.

**Butir 10** — ⬜ **USULAN menunggu persetujuan**: kolom mata uang kelompok (a)(b)(c) jadi
**`NOT NULL` + sentinel `'UNKNOWN'`** sesuai ADR-0006, supaya ketiadaan mata uang **bersuara**.
**Belum diterapkan** — lembar `Kolom` saat ini menandai semuanya `NULL`.

### ✅ Diagram skema rumah TIDAK dipakai — **K-065**

Sumber pembuatan tabel flat adalah **`D:\migrasi\RNM\Claude outputs\`**. Butir **1d** (rekonsiliasi
kolom terhadap Treaty In) **gugur**, dan penghalang yang dicatat K-064 **dicabut**.

⚠️ **Risiko yang diterima sadar:** K-064 tetap berlaku — `T_WORK_POLIS` dan `T_GENERAL_POLIS` adalah
tabel yang **sama** dengan Treaty In, tetapi kolomnya dibangun **hanya dari sisi Fac In**. Irisan
kolom, tabrakan ruang `IDPEGA`, dan nilai sah penanda lini **baru akan muncul saat integrasi**.

### Pekerjaan yang TIDAK menunggu keputusan siapa pun

1. ✅ **SELESAI 24-09** — kelompok (b) K-063 diterapkan: **10 kolom `CURRENCY_CODE VARCHAR2(10)`**
   pada 10 tabel, menaungi 61 kolom uang. Lembar `Kolom` **1.310 → 1.321** kolom.
2. ✅ **SELESAI 24-09** — **DDL draf seluruh 78 tabel**: `08-flat\DDL-tabel-flat-draf.sql`
   (88.028 B, 1.836 baris). Bentuk `T_WORK_POLIS` dan `T_GENERAL_POLIS` dikunci **sesuai lembar
   yang ada** (keputusan work owner 24-09).
3. ⬜ **Satu pertanyaan ke penyusun V-28** untuk butir 8. **Kandidat langkah berikutnya.**

### Isi `DDL-tabel-flat-draf.sql`

| | |
| --- | ---: |
| `CREATE TABLE` | **78** |
| Kolom | **1.329** |
| `PRIMARY KEY` | **78** (surrogate `ID`) |
| `FOREIGN KEY` | **64** |
| Anak berinduk **ganda** — ⛔ tanpa FK | **12** |
| `CREATE INDEX` | **77** (3 `UNIQUE`) |

`[terverifikasi]` Nol identifier berkutip · nol nama constraint/index melebihi 30 karakter · nol nama
ganda · nol `BINARY_DOUBLE` sebagai tipe kolom.

⛔ **Dua belas anak berinduk ganda tidak dapat diberi `FOREIGN KEY`** — satu kolom `PARENT_ID`
menunjuk sampai **5** tabel berbeda (`T_COVERAGELIST`). Keutuhannya ditegakkan di Go, dan
`[terverifikasi]` **12 dari 12** punya kolom `PARENT_TABLE` (V-23) untuk menyimpan tabel mana yang
dimaksud.

⚠️ **Draf. Belum pernah dijalankan terhadap Oracle mana pun** — yang dijamin hanya konsistensi
internalnya. ✅ **Tidak ada lagi butir terbuka yang dapat mengubahnya.**

**Boleh dikerjakan sekarang:** `08-flat\02-<slug>.md` (rancangan versi 115 contoh) → entri
`K-070…` → bahan tiket jalur produksi. ⛔ Kunci `02-` **sudah boleh dibuka**.

⚠️ **Yang belum punya spesifikasi sama sekali:** pemuatan data lama — `SEQ_NO` dibangkitkan saat muat
(V-41) dan **`ROW_UID` direkonstruksi** untuk seluruh riwayat polis (V-42, sendirinya masih terbuka).
Itu, bukan DDL-nya, yang akan memakan waktu.

⛔ **Jangan menambah seam baru.** **Empat** seam disetujui: `premium.Calculate`,
`acceptance.Next`, `rules.Eval`, dan **`loader.Flatten`** (K-074, 25-09) — ditambah seam `repository`
yang **dicadangkan sejak 16-09** dan kini aktif.

---

## 9. Peta berkas penting

| Berkas | Isi |
| --- | --- |
| `00-KEPUTUSAN-WORK-OWNER.md` | Register **K-001…K-062**. **Sumber otoritatif** — mengalahkan `[pertanyaan terbuka]` di dokumen mana pun |
| `steering\PANDUAN-KERJA.md` · `steering\GLOSARIUM.md` | Aturan operasional · kosakata |
| `adr\0001…0006` | Keputusan rancangan |
| `08-flat\00a` → `00` → `01` | **Pekerjaan aktif**, urutan baca begitu |
| `_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md` | T-1…T-9 siap kirim |
| `10-audit\01…11` | Audit Fac Out + daftar periksa UI Pega |
| `05-tickets\` · `\rnw` · `\edm` · `\facout` | 16 · 8 · 22 · 14 tiket |
| `07-edm\07-e6-penutup.md` | Ringkasan discovery endorsement |
| `_ARSIP-lintas-siklus\_BACA-INI.md` | Temuan lintas-siklus yang masih berlaku + **Koreksi R1** (rule kembar) |
| `_TOKEN-LOG.md` | Ledger token, terakhir Prompt 73 |

**Sudah basi — sejarah saja, jangan dipakai sebagai keadaan sekarang:**
`_HANDOFF-SESI.md` (16-09) · `_HANDOFF-EDM.md` (21-09) · `_ARSIP-lintas-siklus\` (baca `_BACA-INI.md`
dulu).

---

## 10. Jebakan yang sudah memakan biaya — jangan diulang

1. **Nama rule bukan bukti** · nama sama tipe beda · posisi kata · sufiks · **kapitalisasi** ·
   **kelas sumber data** · **nama folder bukan tipe rule** (`RDBList\` berisi `Rule-Connect-SQL`).
2. **Nama berkas sama ≠ rule sama.** Periksa basis `pzInsKey` + `pyClassName` lebih dulu.
3. Tag ber-atribut **tidak** tertangkap pencarian berpenutup: `<pzIndexes>` → 0, `<pzIndexes` → 82.
4. **Nomor baris tidak menentukan keanggotaan langkah** — telusuri pohon `pySteps/rowdata`, dan
   bedakan **anak LANGSUNG** dari keturunan.
5. **Ekspresi gerbang saja belum menentukan arah** — baca `pyStepsPreCondParamsWhenTrue/False`.
6. **PowerShell**: variabel **tidak peka huruf** (`$B` vs `$b` — sudah berkali-kali kena);
   `Get-Content -Raw` wajib `-Encoding utf8`; `\"` bukan escape (pakai backtick); `-f` di dalam
   pemanggilan metode memecah argumen; `@()` membungkus ulang List yang dikembalikan fungsi.
7. **Regex Oracle**: `CREATE OR REPLACE **EDITIONABLE** PROCEDURE` — kata `EDITIONABLE` menyisip di
   tengah; regex tanpa itu melaporkan 9, bukan 25.
8. Kolom `pyStepsRepeatDef` dan `pyStepsPreCondParams` **memuat nama orang + email** — jangan dicetak.
9. **Sensus label paling akhir** — sudah tiga kali meleset karena dijalankan sebelum suntingan selesai.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
