# Spec — Endorsement Life (migrasi Pega → Go + React + Oracle)

> **Ralat 01-10-2026** (gelombang 2 brief, `RALAT-DEV-01-10-2026.md` — ralat mengalahkan isi di bawah). Teks lama yang tidak berlaku:
> - **E2** — §16 / Kontrak batas: *"⚠️ `POOLDATA.M_LIFE_PREMIUM_DETAIL` / `_SUMMARY` adalah bentuk **sistem lama** — bukan target tulis sistem baru"* → Endorsement **menulis** keduanya persis `SaveMasterLPDet`/`InsertPLSummary`, di samping versi baru tabel aplikasi.
> - **R01** — §7 *"`<nomor polis>/<PRODKE + 1>`"* → nomor polis = `PL_NUMBER` (`InsertJsonPolisEDM` b102).
> - **R04** — AC 31 dianggap sejalan dengan Pega → Pega menahan kasus `Decline` (`FilterProteksiEDMLife` b526); keputusan work owner diikuti, OQ-EDM-002.
> - **R05** — §1 *"di jalur endorsement tidak dijalankan"* → korpus memanggilnya di `SaveCSVEDMLife` langkah 3 b629 (`·`); keputusan work owner diikuti, OQ-EDM-003.
> - **R20** — §16 *"Endorsement adalah **baris versi baru** di `T_PREMIUM_LIST`"* beserta `T_WORK_POLIS` sejajar (STRUKTUR) → baris versi = kasus itu sendiri; **tanpa** baris `T_WORK_POLIS`.

Status: ready-for-agent
Konteks: `endorsement-life` — "Life — Endorsement" (konteks/menu **terpisah** dari PremiumList Life)
Modul: **Endorsement Life** (75 berkas), class work `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE`
Tanggal: 2026-09-15 · **direvisi 2026-09-16 — bentuk penyimpanan**

> ⚠️ **Revisi penyimpanan 2026-09-16** `[keputusan work owner]`. **Seluruh JSON dibuang**, dan
> endorsement **berbagi tujuh tabel PremiumList Life yang sama** — ia **versi baru**, bukan entitas
> tersendiri. Pencocokan peserta lama↔baru lewat **`PARENT_ID`** (pointer eksplisit), menggantikan
> loop indeks Pega.
> Yang berubah: **§4** (memuat & menyalin), **§11** (batas transaksi — titik potong lenyap),
> **§16** (baru — bentuk penyimpanan & versi), **AC 54–72** (baru), §Pertanyaan terbuka, §Out of
> Scope.
> **Yang sengaja TIDAK berubah**: gerbang kelayakan (§3), `EdmType`/`EDMStatus` (§5), jurnal balik
> (§6), penomoran (§7), mesin alur (§8), kunci field (§9), CSV (§10), efek keluar (§12), Arasapas
> (§13), penyaringan peserta klaim (§14), kode mati (§15), dan **AC 1–53**.
> Sumber: `keputusan-struktur-edm.md` + `.scratch/premiumlist-life/revisi-penyimpanan-premiumlist.md`.

Sumber: `.scratch/endorsement-life/keputusan-struktur-edm.md` (struktur & versi),
`.scratch/endorsement-life/grilling-ronde-1.md` (12 verdict),
`.scratch/endorsement-life/grilling-ronde-2.md` (6 verdict), `CONTEXT.md`,
`docs/adr/ADR-0001`–`ADR-0015`, `.scratch/premiumlist-life/spec.md` (mesin bersama),
`discovery/open-questions.md`
Skill: `/mattpocock-skills:to-spec`

> **Konvensi penandaan.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[keputusan desain]`; `[data DBA]`; `[terbuka]` = OQ.
> **Identitas rule wajib menyertakan class** — nama sama di class berbeda = rule berbeda.

> **Sumber tunggal.** Korpus Pega `D:\XML\RNM_BRD\` (READ-ONLY) dan artefak di
> `OUTPUT_HASIL_RNM\`, sekaligus repo target tunggal. ⛔ `D:\XML\nusantara-re\` di-blacklist.

---

## Problem Statement

Polis Life yang sudah berjalan **berubah**. Peserta bertambah atau keluar, datanya salah dan perlu
dikoreksi, atau seluruh polis dibatalkan. Hari ini perubahan itu dikerjakan di Pega lewat menu
Endorsement Life yang terpisah dari input polis baru — dan pengetahuan tentang cara kerjanya
tersebar di 75 berkas rule yang sebagian sudah mati, sebagian menyimpan aturan penting di dalam
deskripsi langkah yang **bertentangan dengan kodenya**.

Empat hal membuat konteks ini berisiko dipindahkan secara salah:

1. **Pembatalan bukan penghapusan.** Endorsement batal **tidak** menghapus baris dan **tidak**
   menulis nol — ia menulis **baris bernilai negatif** ke tabel yang sama. Deskripsi langkah Pega
   berbunyi "Set 0 jika EDM Batal", tetapi kodenya mengalikan **32 kolom uang dengan `-1`**. Siapa
   pun yang memigrasikan dari deskripsi akan menghasilkan saldo yang salah.
2. **Akuntansi dan klaim melihat data yang berbeda.** Baris negatif harus terlihat akuntansi (agar
   net benar), tetapi peserta yang batal atau dihapus **tidak boleh muncul** di Claim Life. Satu
   tabel, dua sudut pandang.
3. **Nomor endorsement dirakit dari data, bukan dari sequence.** `PL_NUMBER_EDM` berbeda mekanisme
   dari `PL_NUMBER`, dan di Pega **dua rule membacanya dengan urutan berbeda** — celah yang dapat
   menghasilkan nomor salah.
4. **Batas transaksi terpotong-potong.** Empat rule SQL commit sendiri, ditambah dua `Commit`
   eksplisit yang berjalan **sebelum** pengguna menekan simpan.

## Solution

Membangun ulang Endorsement Life sebagai **menu tersendiri** di atas Go + React + Oracle, dengan
perilaku dibawa apa adanya (paritas) kecuali **tujuh penyimpangan sadar** yang sudah diputus.

Alur inti yang **tidak ada** di konteks PremiumList Life:

**pilih polis new business lewat nomor polis → lima gerbang kelayakan → buat case + muat & salin
data polis lama → sunting (Perubahan Data) atau batalkan → `Confirm` / `Decline` → simpan.**

`[keputusan work owner]` **Endorsement adalah entri baru yang menunjuk polis NB** lewat nomor polis
+ `PRODKE` — **bukan mutasi in-place**. Baris polis lama tidak pernah dihapus fisik.

### Tujuh penyimpangan sadar

| # | Penyimpangan | Alasan | Sumber |
| --- | --- | --- | --- |
| 1 | Pembacaan `PRODKE` disatukan ke **`ORDER BY PRODKE DESC`** | Pega punya dua urutan berbeda → nomor EDM bisa salah | `[keputusan desain]` |
| 2 | **Penjaga anti-dobel** `(NOPOLIS, PRODKE)` sebelum menulis **versi baru** | asalnya: penulis JSON lama `INSERT` polos + commit sendiri. Sekarang alasannya berbeda — penyalinan dua kali menghasilkan **dua versi** (§11, §16) | `[keputusan desain]` |
| 3 | **Alarm dihidupkan** (deteksi keadaan separuh + email kegagalan) | Pega menjalankan endorsement **tanpa alarm**; jalur new business memilikinya | `[keputusan work owner]` |
| 4 | Pembacaan Arasapas **dikurung di satu repository** bertanda batas lintas sistem | mencegah kueri lintas skema tersebar | `[keputusan desain]` |
| 5 | Rule gerbang **diganti nama** — ia gerbang kelayakan, bukan pembatalan | nama Pega menyesatkan | `[keputusan work owner]` |
| 6 | **Tanpa batas jumlah baris CSV** (50.000 dibuang) | batas teknis, bukan aturan bisnis; unggahan besar diproses bertahap | `[keputusan work owner]` |
| 7 | **Peserta batal/delete disaring keluar** dari jalur baca klaim | akuntansi melihat seluruh baris; klaim hanya peserta hidup | `[keputusan work owner]` |

### Kontrak batas

| Arah | Isi |
| --- | --- |
| **Hulu** | Endorsement **membaca** polis new business dari `POOLDATA.JSON_POLIS` (berkunci `NOPOLIS`, termasuk `PRODKE` dan atribut `DATA_JSON.EdmType` di dalam CLOB) dan dari `POOLDATA.M_LIFE_PREMIUM_DETAIL` |
| **Hilir** | Endorsement **menulis** baris — **termasuk baris bernilai negatif** — sebagai **versi baru di tujuh tabel PremiumList Life** (`T_PREMIUM_LIST` + `_DETAIL` + `_SPREADING` + `_SPREADING_RETRO` + `_SUMMARY`), dengan **`PARENT_ID`** pada `_DETAIL` (**§16**), berkunci `PL_NUMBER_EDM` di samping `PL_NUMBER`. ⚠️ `POOLDATA.M_LIFE_PREMIUM_DETAIL` / `_SUMMARY` adalah bentuk **sistem lama** — bukan target tulis sistem baru |
| **Hilir — penyaringan** | ⚠️ Peserta ber-`EDMStatus` **`"Batal"`** atau **`"Delete"`** **tidak boleh muncul** di Claim Life. Jalur baca klaim **wajib menyaringnya keluar** |

`[terverifikasi]` Class integrasi `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` adalah tulang punggung
domain Life: Claim Life 46 rule, **Endorsement Life 25**, PremiumList Life 19, Komite Claim Life 11.

---

## User Stories

### Memilih polis dan kelayakan

1. Sebagai **inputor Life**, saya ingin memilih polis yang akan di-endorse dengan **nomor polis**,
   supaya saya bekerja dengan identitas yang sama seperti yang dipakai seluruh perusahaan.
   `[keputusan work owner]` — kunci di seluruh rule endorsement adalah `NOPOLIS`.
2. Sebagai **inputor**, saya ingin dapat **menelusuri daftar premium list** untuk menemukan nomor
   polis yang saya cari, supaya saya tidak perlu menghafalnya. `[terverifikasi]`
   `BrowsePremiumList_RD` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `BROWSEPREMIUMLIST_RD`) —
   `[keputusan work owner]` ia **alat bantu pencarian, bukan sumber kunci**.
3. Sebagai **organisasi**, saya ingin endorsement **ditolak bila nomor polis kosong**, supaya tidak
   ada case tanpa sasaran.
4. Sebagai **organisasi**, saya ingin endorsement **ditolak bila polisnya tidak ada** di rekam
   polis, supaya endorsement tidak menggantung tanpa induk.
5. Sebagai **organisasi**, saya ingin **hanya ada satu endorsement terbuka per polis**, supaya dua
   perubahan tidak saling menimpa. Ini **aturan integritas inti** konteks ini.
6. Sebagai **organisasi**, saya ingin polis yang **sudah pernah dibatalkan** tidak dapat di-endorse
   lagi, supaya polis mati tidak dihidupkan lewat pintu belakang.
7. Sebagai **Finance**, saya ingin polis yang **sudah dibayar** tidak dapat dibatalkan, supaya
   pembatalan tidak menciptakan selisih terhadap kas yang sudah masuk.
8. Sebagai **inputor**, saya ingin **tahu alasan penolakan** saat endorsement tidak dapat dibuat,
   supaya saya dapat memperbaikinya sendiri.

### Memuat data polis lama

9. Sebagai **inputor**, saya ingin data polis lama **otomatis termuat** ke endorsement yang baru
   dibuat, supaya saya tidak mengetik ulang ratusan peserta.
10. Sebagai **inputor**, saya ingin dapat **melihat polis lama** dalam tampilan tersendiri sebelum
    mengubah apa pun, supaya saya yakin memilih polis yang benar.
11. Sebagai **inputor**, saya ingin tampilan polis lama **menyesuaikan jenis transaksi** (`QP`,
    `QR`, `TP`, `TR`), supaya kolom yang saya lihat relevan. `[terverifikasi]` `QR` memakai layar
    induk; `QP`/`TP`/`TR` memakai varian.
12. Sebagai **organisasi**, saya ingin baris warisan polis lama **ditandai sebagai `Old`**, supaya
    dapat dibedakan dari baris yang ditambahkan endorsement.

### Dua maksud endorsement

13. Sebagai **inputor**, saya ingin dapat memilih maksud endorsement: **Perubahan Data** atau
    **Batal**, supaya sistem tahu perlakuan apa yang berlaku. `[keputusan work owner]` `EdmType`
    `1` = Perubahan Data, `3` = Batal; `2` dan `4` tidak dipakai.
14. Sebagai **inputor**, pada **Perubahan Data** saya ingin dapat **menambah peserta**, supaya
    peserta baru masuk pertanggungan. Baris bertanda `New`.
15. Sebagai **inputor**, pada **Perubahan Data** saya ingin dapat **menghapus peserta tertentu**,
    supaya peserta yang keluar tidak lagi ditanggung. Baris bertanda `Delete`.
16. Sebagai **inputor**, saya ingin **Batal** membatalkan **seluruh peserta sekaligus**, supaya saya
    tidak perlu menandai satu per satu. `[keputusan work owner]`
17. Sebagai **organisasi**, saya ingin maksud endorsement **terkunci** setelah case dibuat, supaya
    endorsement perubahan data tidak diam-diam berubah menjadi pembatalan.

### Jurnal balik

18. Sebagai **Finance**, saya ingin pembatalan menghasilkan **baris bernilai negatif**, bukan baris
    nol dan bukan penghapusan, supaya jejak transaksi asli tetap ada dan nettonya benar.
19. Sebagai **Finance**, saya ingin **penghapusan peserta** juga menghasilkan nilai negatif — hanya
    **untuk peserta itu**, supaya koreksi bersifat selektif. `[keputusan work owner]`
20. Sebagai **Finance**, saya ingin baris positif asli dan baris negatif **hidup berdampingan** di
    tabel yang sama, supaya laporan cukup menjumlahkan.
21. Sebagai **Finance**, saya ingin **seluruh komponen uang** ikut dibalik — premi, komisi,
    brokerage, tax, admin fee, klaim, dan seluruh kelompok refund dan retro — supaya tidak ada
    komponen yang tertinggal positif.
22. Sebagai **Finance**, saya ingin **tidak ada nilai uang yang berubah** saat menyeberang batas
    penyimpanan, termasuk nilai negatif.

### Nomor endorsement

23. Sebagai **organisasi**, saya ingin endorsement memperoleh **nomornya sendiri**
    (`PL_NUMBER_EDM`), supaya penomoran new business tidak tertimpa.
24. Sebagai **organisasi**, saya ingin nomor endorsement **lahir sekali**, supaya satu endorsement
    tidak pernah punya dua nomor.
25. Sebagai **organisasi**, saya ingin nomor produksi (`PRODKE`) **bertambah satu** setiap
    endorsement atas polis yang sama, supaya urutannya terbaca.
26. Sebagai **tim migrasi**, saya ingin `PRODKE` dibaca dengan **satu urutan yang sama** di seluruh
    sistem, supaya nomor endorsement tidak bergantung pada rule mana yang kebetulan dipanggil.

### Alur keputusan

27. Sebagai **atasan**, saya ingin **menyetujui** endorsement sebelum tersimpan permanen, supaya
    perubahan polis tidak terjadi tanpa persetujuan.
28. Sebagai **atasan**, saya ingin dapat **menolak** endorsement, dan penolakan itu **menutup**
    case-nya. `[terverifikasi]` Endorsement **tidak punya** keluaran `Reject` — hanya `Confirm` dan
    `Decline`, berbeda dari PremiumList Life.
29. Sebagai **inputor**, saya ingin **membatalkan endorsement yang salah** dengan menolaknya, supaya
    polisnya bebas di-endorse ulang. `[keputusan work owner]` — tidak ada tombol batal terpisah.
30. Sebagai **organisasi**, saya ingin setiap transisi tahap **meninggalkan jejak audit**, supaya
    dapat ditelusuri siapa mengubah apa dan kapan.

### Unggah CSV

31. Sebagai **inputor**, saya ingin **mengunggah CSV** untuk menambah peserta secara massal, supaya
    tidak mengetik satu per satu.
32. Sebagai **inputor**, saya ingin **mengunggah ulang** berkas yang diperbaiki dan hasilnya
    **mengganti**, bukan menumpuk. `[terverifikasi]` — baris `New` sebelumnya dibuang lebih dulu.
33. Sebagai **inputor**, saya ingin unggahan **tidak dibatasi jumlah barisnya**, supaya polis besar
    tidak tertolak karena alasan teknis. `[keputusan work owner]`
34. Sebagai **organisasi**, saya ingin seluruh peserta dalam satu endorsement **satu plan dan satu
    pemegang polis**, supaya tidak tercampur polis lain. `[terverifikasi]` — tiap baris diuji
    terhadap baris pertama.
35. Sebagai **inputor**, saya ingin **pesan kesalahan menyebut baris dan kolom** yang salah, supaya
    saya dapat memperbaikinya.
36. Sebagai **organisasi**, saya ingin CSV **tidak menambah peserta pada endorsement Batal**, supaya
    pembatalan tidak disusupi peserta baru. `[terverifikasi]` — penandaan `New` hanya berlaku pada
    `EdmType=1`.

### Efek keluar dan lingkungan

37. Sebagai **tim operasi**, saya ingin endorsement yang tersimpan **diteruskan ke Arasapas**,
    supaya sistem hilir menerimanya.
38. Sebagai **tim operasi**, saya ingin alamat layanan keluar **dibaca runtime dari tabel**, supaya
    perpindahan lingkungan tidak menuntut penempelan ulang.
39. Sebagai **tim operasi**, saya ingin **diberi tahu bila penyimpanan gagal**, supaya kegagalan
    tidak baru ketahuan saat tutup buku. `[keputusan work owner]` — alarm **dihidupkan**, seragam
    dengan new business.
40. Sebagai **tim operasi**, saya ingin efek keluar **tidak berjalan di lingkungan non-production**,
    supaya uji coba tidak mengotori sistem hilir.

### Kontrak hilir

41. Sebagai **admin klaim Life**, saya ingin peserta yang **sudah dibatalkan atau dihapus** **tidak
    muncul** saat saya mencari peserta untuk klaim, supaya saya tidak memproses klaim atas peserta
    yang tidak lagi ditanggung. `[keputusan work owner]`
42. Sebagai **Finance**, saya ingin peserta yang sama **tetap terlihat** dalam laporan akuntansi
    beserta baris negatifnya, supaya nettonya dapat dihitung.
43. Sebagai **admin klaim Life**, saya ingin rekam premium hasil endorsement **dapat ditemukan**
    lewat `PL_NUMBER_EDM` maupun `PL_NUMBER`, supaya klaim atas polis yang di-endorse tetap
    terhubung.

### Yang sengaja tidak dibawa

44. Sebagai **tim migrasi**, saya ingin **batas 50.000 baris** tidak ikut pindah.
45. Sebagai **tim migrasi**, saya ingin rule pemuat data lama yang sudah mati **tidak ditiru apa
    adanya** — hanya perilakunya yang dipertahankan.
46. Sebagai **tim migrasi**, saya ingin **mesin hitung milik new business tidak dipanggil** di jalur
    endorsement, karena di Pega pun ia tidak dijalankan di sini.
47. Sebagai **tim migrasi**, saya ingin **commit implisit Pega tidak ditiru** — Go memegang batas
    transaksi.

---

## Implementation Decisions

### 1. Batas konteks

`[keputusan work owner]` **Endorsement Life adalah konteks/menu terpisah** dari PremiumList Life.
Alasannya `[terverifikasi]`:

| Pemisah | PremiumList Life | Endorsement Life |
| --- | --- | --- |
| Class work | `ASM-FW-GISFW-WORK-LIFE` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` |
| Flow masuk | `PremiumList Life/InputPolicyHolder.xml` (`INPUTPOLICYHOLDER` / `RULE-OBJ-FLOW`) | `Endorsement Life/Flow/InputEDMLife.xml` (`INPUTEDMLIFE` / `RULE-OBJ-FLOW`) |
| Alur inti | penawaran baru → premium list → simpan polis | **pilih polis NB → muat data lama → endorse** |

**Mesin yang dipakai bersama** — komponen, bukan alasan menyatukan menu:

| Mesin bersama | Identitas | Sifat |
| --- | --- | --- |
| Penulis detail | `SaveMasterLPDet` — `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | **satu rule identik** di kedua modul |
| Penulis summary | `InsertPLSummary` — `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | **satu rule identik** |
| Tabel | `POOLDATA.M_LIFE_PREMIUM_DETAIL`, `M_LIFE_PREMIUM_SUMMARY`, `JSON_POLIS` | dipakai kedua jalur **di sistem lama**. ⚠️ Di sistem baru **`JSON_POLIS` hanya DIBACA** — polis new business lama, sampai migrasi selesai. Penyimpanan endorsement = **tujuh tabel PremiumList Life** (**§16**) |
| Orkestrator simpan | `InsertJsonPolisLife_Act` | ⚠️ **nama sama, class berbeda → rule BERBEDA** |

⚠️ `[keputusan work owner]` **`Calculate1_Act` BUKAN mesin bersama.** Berkasnya identik di kedua
folder (`ASM-FW-GISFW-WORK-LIFE` / `CALCULATE1_ACT` / `RULE-OBJ-ACTIVITY`, diff ternormalisasi nol
baris) **tetapi di jalur endorsement tidak dijalankan**. Perhitungan endorsement dikerjakan
`SetPremi_EDM` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM`).

Claim Life dan Komite Claim Life adalah **konteks hilir**.

### 2. Penempatan modul

`[terverifikasi]` Struktur mengikat `CLAUDE.md` §5: kode di-scaffold **di dalam
`OUTPUT_HASIL_RNM\`**, arah dependency `handlers` → `services` → `repository`.

| Lapisan | Tanggung jawab |
| --- | --- |
| `handlers` | Endpoint REST: cari polis, cek kelayakan, buat endorsement, sunting detail, unggah CSV, putuskan, simpan |
| `services` | Lima gerbang kelayakan; pemuatan & penyalinan data polis lama; mesin `EDMStatus`; jurnal balik; penomoran EDM; orkestrasi efek keluar |
| `repository` | Baca `JSON_POLIS` & detail lama (polis NB, untuk disalin); **tulis versi baru ke tujuh tabel PremiumList Life** (header + detail + spreading + retro + summary); **satu repository terpisah** untuk pembacaan Arasapas |
| `models` | Endorsement, baris detail beserta `EDMStatus`, ringkasan premium |
| `frontend` | Layar Input EDM, popup lihat polis lama (4 varian), grid detail, unggah + tinjau CSV |

### 3. Gerbang kelayakan — **lima pemeriksaan, sebelum case dibuat**

`[terverifikasi]` Sumbernya `SetErrorBatalEndorsement_Act` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` /
`SETERRORBATALENDORSEMENT_ACT` / `RULE-OBJ-ACTIVITY`).

⚠️ `[keputusan work owner]` **Namanya diganti** di sistem baru — ia **gerbang kelayakan
endorsement**, bukan pembatalan.

| # | Penolakan | Cara memeriksa |
| --- | --- | --- |
| 1 | nomor polis kosong | validasi masukan |
| 2 | polis tidak ada di rekam polis | lookup `JSON_POLIS` berkunci `NOPOLIS` |
| 3 | **sudah ada endorsement yang belum selesai** | pencarian endorsement berjalan atas polis itu |
| 4 | polis **sudah pernah dibatalkan** | `EdmType` polis terakhir bernilai `3` |
| 5 | polis **sudah dibayar** **dan** maksudnya Batal | lookup pembayaran Arasapas |

`[keputusan work owner]` Case endorsement dibuat **hanya setelah kelima gerbang lolos**.

### 4. Memuat & menyalin data polis lama

`[terverifikasi]` Mesinnya `MappingEDMLife` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `MAPPINGEDMLIFE` /
`RULE-OBJ-ACTIVITY`, **14 langkah, nol langkah ter-remark**): buat case → ambil `IDPEGA` polis
terakhir dari `JSON_POLIS` → simpan **ID Pega lama + nomor polis lama** → buka polis NB → salin
detail dan halamannya → tandai baris warisan **`Old`**.

`[keputusan work owner]` **Popup "lihat polis lama" tetap ada** dan dipicu per `.Type`:

| `.Type` | Layar |
| --- | --- |
| `QR` | section **induk** |
| `QP` / `TP` / `TR` | varian khusus |

⚠️ Rule Pega `GetOldDetail_EDM` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GETOLDDETAIL_EDM`,
**4 dari 5 langkah mati**) **tidak ditiru apa adanya**, tetapi **perilakunya wajib ada** — popup
tetap harus terisi data polis lama.

### 5. `EdmType` dan `EDMStatus` — dua maksud, empat status

`[keputusan work owner]` **`EdmType`**: `1` = **Perubahan Data**, `3` = **Batal**. Nilai `2` dan `4`
tidak dipakai — enum efektif `{1, 3}`. `[terverifikasi]` Nilai `3` terbukti dari kode: precondition
`pyWorkPage.EdmType==3` di `SetPremi_EDM`. `[terverifikasi]` `EdmType` tersimpan **di dalam CLOB**
`DATA_JSON`, **bukan** sebagai kolom.

`[keputusan work owner]` **`EDMStatus`** — status per baris detail. Bedanya adalah **cakupan minus**,
bukan jenis:

| Nilai | Kapan | Akibat pada nilai uang |
| --- | --- | --- |
| `Old` | warisan polis new business | tidak diubah |
| `New` | peserta ditambah — **hanya pada `EdmType=1`** | positif, baris baru |
| `Delete` | peserta dihapus dalam Perubahan Data | **diminuskan — selektif per peserta** |
| `Batal` | lewat `EdmType=3` | **seluruh peserta otomatis batal — diminuskan menyeluruh** |

⚠️ **`Delete` bukan sekadar penanda** — ia **juga** menghasilkan nilai negatif.

### 6. Jurnal balik — **`× -1`, bukan nol**

`[terverifikasi]` `SetPremi_EDM` (**9 langkah, nol langkah ter-remark**) mengalikan **32 kolom uang
dengan `-1`** saat maksudnya batal atau peserta dihapus: `SUM_INSURED`, `CEDING_RETENTION`,
`SUM_REASURED`, `SHARE_NUSANTARA_RE`, `SHARE_NUSANTARA_RE_GROSS`, `SUM_AT_RISK_GROSS`,
`SUM_AT_RISK_RETRO`, `RETROCEDED_SHARE`, `SHARE_RETRO`, `RATE`, `FACTOR`, `GROSS_PREMIUM`,
`NET_PREMIUM`, `DEDUCTION`, `CLAIM_AMOUNT`, `RI_ADMIN_FEE`, `BROKERAGE_FEE`, beserta seluruh
kelompok `*_REFUND`, `*_RETRO`, dan `*_REFUND_RETRO`.

⚠️ **Deskripsi langkah Pega berbunyi "Set 0 jika EDM Batal" — itu menyesatkan.** Kodenya yang
berlaku. `[keputusan work owner]` pembalikan tanda memang perilaku akuntansi yang dikehendaki.

`[keputusan work owner]` Baris negatif ditulis ke **tabel yang sama**; baris positif asli dan baris
negatif **hidup berdampingan** → net akunting dari penjumlahan, bukan dari penghapusan.

### 7. Penomoran `PL_NUMBER_EDM`

`[terverifikasi]` **Bukan** dari sequence terpusat. Ia dirakit **`<nomor polis>/<PRODKE + 1>`** oleh
`Generate_NoEndorsmentLife` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!GENERATE_NOENDORSMENTLIFE`
/ `RULE-CONNECT-SQL`). **ADR-0006 tidak berlaku pada nomor EDM** — jangan memaksanya lewat
`PROC_GENERATE_SEQUENCE_NUMBER`.

`[terverifikasi]` **Idempoten di korpus**: dua langkah berturut-turut di `GenerateNoEDM_Life`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GENERATENOEDM_LIFE`) dijaga precondition
`PL_NUMBER_EDM==""`.

⚠️ `[keputusan desain]` **Penyimpangan sadar.** Pega punya dua pembaca `PRODKE` dengan urutan
berbeda:

| Rule | Identitas | Urutan |
| --- | --- | --- |
| `GetProdkeNopolis` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETPRODKENOPOLIS` / `RULE-CONNECT-SQL` | **`PRODKE DESC`** |
| `GetProdKeOldData_SQL` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!GETPRODKEOLDDATA_SQL` / `RULE-CONNECT-SQL` | `TGL_INPUT desc` |

Sistem baru memakai **`ORDER BY PRODKE DESC` di kedua tempat**. `PRODKE` adalah urutan produksi yang
menjadi dasar nomor; `TGL_INPUT` dapat menyimpang karena entri susulan atau koreksi.

### 8. Mesin alur

`[terverifikasi]` Flow `InputEDMLife` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` /
`RULE-OBJ-FLOW`): assignment **Input EDM Detail** dan **Input EDM Summary**, keputusan
`IsLifeAccepted`, utility simpan, routing **`WorkList`**, status akhir **`Resolved-Completed`** dan
**`Resolved-Rejected`**.

⚠️ `[terverifikasi]` **`IsLifeAccepted` versi Endorsement** (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` /
`ISLIFEACCEPTED` / `RULE-OBJ-DECISIONTABLE`) hanya mengeluarkan **`Confirm`** dan **`Decline`** —
**tanpa `Reject`**, berbeda dari versi PremiumList Life. **Endorsement tidak punya jalur "kembali ke
input"**.

`[keputusan work owner]` **Membatalkan endorsement = `Decline`** → case tertutup
(`Resolved-Rejected`), polis bebas di-endorse ulang. Tidak ada tombol batal terpisah.

### 9. Kunci field permanen

`[terverifikasi]` Editabilitas per form **bukan** rule `When` bernama — modul ini hanya punya dua
rule `When` (`@BASECLASS` / `ISPEGAPROD`, `@BASECLASS` / `RECORDEVENT`), keduanya tak terkait.
Mekanismenya **kondisi sebaris `<pyDisabledWhen>`** atas properti `.EditInput` / `.EditInput1`.

`[terverifikasi]` `.EditInput` adalah **kunci satu arah** — tiga penulis, semuanya ke `1`, tidak ada
yang mengembalikannya ke `0`. Field terkunci: **nomor polis**, **jenis endorsement/batal**.

`[keputusan work owner + desain]` Sistem baru **mempertahankan kunci permanen** dan **menambahkan
pesan penjelas** — di Pega field hanya mati tanpa alasan. Jalan keluar bila salah pilih polis tetap
`Decline` lalu buat baru.

### 10. Unggah CSV

`[terverifikasi]` Impor dikerjakan `UploadCSVEDMLifePremium_Act` (`@BASECLASS` /
`UPLOADCSVEDMLIFEPREMIUM_ACT` — **rule base-class**, hanya tiga langkah pertama yang hidup);
penyimpanan oleh `SaveCSVEDMLife` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE`,
**5 langkah, nol ter-remark**).

Tiga sifat yang dibawa `[terverifikasi]`:

1. **Unggah ulang mengganti, bukan menumpuk** — baris `New` sebelumnya dibuang lebih dulu.
2. **Baris baru hanya pada `EdmType=1`** — pada endorsement Batal, CSV tidak menambah peserta.
3. **Validasi konsistensi** — `PLAN` dan `POLICY_HOLDER` tiap baris diuji terhadap **baris pertama**.

⚠️ `[keputusan work owner]` **Batas 50.000 baris dibuang.** Tidak ada batasan jumlah baris. Unggahan
sangat besar diproses **bertahap** (streaming/batch internal) — **tanpa menolak karena jumlah
baris**.

Nilai uang di CSV tunduk **ADR-0003**: parse langsung ke desimal presisi arbitrer, **tidak** lewat
`float`, dengan format pemisah yang dinyatakan eksplisit.

### 11. Batas transaksi — ⚠️ **DIREVISI 2026-09-16: titik potong lenyap**

⚠️ **Tiga dari empat titik potong adalah penulis JSON atau penulis tabel warisan — seluruhnya
dibuang** (§16). `InsertJsonPolisEDM` tidak dipanggil; `SaveLifeinProduction_SQL`,
`SaveMasterLPDet`, dan `InsertPLSummary` menulis tabel existing yang di sistem baru **hanya
dibaca saat migrasi**.

**Aturan mengikat yang berlaku sekarang** `[keputusan work owner]`:

> **Satu endorsement — penyalinan versi, seluruh peserta hasil salin beserta `PARENT_ID`-nya,
> seluruh spreading & spreading retro, kolom EDM pada header, dan rekap mata uang yang dihitung
> ulang — ditulis dalam SATU transaksi, lalu commit sekali.**

⚠️ **Penjaga anti-dobel tetap berlaku, dengan alasan baru.** Dulu ia diperlukan karena
`InsertJsonPolisEDM` adalah `INSERT` polos yang commit sendiri. Sekarang ia diperlukan karena
**endorsement membuat baris versi baru**: menjalankan penyalinan dua kali menghasilkan **dua versi**
untuk satu maksud endorsement. Pemeriksaan **`(NOPOLIS, PRODKE)` sebelum menulis** karena itu
**dipertahankan** — bukan warisan, melainkan kebutuhan versi.

⚠️ **Dua `Commit` dini dipertahankan** `[keputusan work owner]` — di `CreateCaseEMDL`
(`DATA-PORTAL` / `CREATECASEEMDL`) dan `MappingEDMLife`. Gerbang "satu endorsement terbuka per
polis" **membutuhkan** case terlihat pengguna lain begitu dibuat. Ia **di luar** transaksi simpan.

<details>
<summary>Keadaan lama (sudah tidak berlaku) — titik potong `[terverifikasi]`</summary>

`[terverifikasi]` Titik potong di jalur endorsement:

| Yang commit sendiri | Identitas |
| --- | --- |
| `InsertJsonPolisEDM` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLISEDM` / `RULE-CONNECT-SQL` |
| `SaveLifeinProduction_SQL` | `ASM-FW-GISFW-WORK-LIFE` / `ASM!SAVELIFEINPRODUCTION_SQL` / `RULE-CONNECT-SQL` |
| `SaveMasterLPDet` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` |
| `InsertPLSummary` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` |

Ditambah **dua `Commit` aktif** yang berjalan **sebelum** pengguna menekan simpan: di
`CreateCaseEMDL` (`DATA-PORTAL` / `CREATECASEEMDL`) dan `MappingEDMLife`.

`[keputusan work owner]` Commit dini itu **wajar dan dipertahankan** — gerbang "satu endorsement
terbuka per polis" justru **membutuhkan** case terlihat oleh pengguna lain begitu dibuat.

**Aturan mengikat:** Go memegang batas transaksi eksplisit dan memperlakukan tiap pihak yang commit
sendiri sebagai **titik potong**. Data yang harus atomik tidak boleh dipisahkan olehnya.

⚠️ `[keputusan desain]` **Penyimpangan sadar — penjaga anti-dobel.** `InsertJsonPolisEDM` adalah
**`INSERT` polos**, bukan upsert, dan commit sendiri. Pengulangan menghasilkan rekam polis kedua.
Jalur new business tidak punya masalah ini karena memakai procedure upsert berkunci `IDPEGA`.
Sistem baru **memeriksa `(NOPOLIS, PRODKE)` sebelum menulis** dan tidak menulis ulang bila sudah ada.

</details>

### 12. Efek keluar

`[terverifikasi]` `serviceInsertArasapasLife_act` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` /
`SERVICEINSERTARASAPASLIFE_ACT` / `RULE-OBJ-ACTIVITY` — ⚠️ **rule berbeda** dari yang ber-class
`ASM-FW-GISFW-WORK-LIFE`), dijaga flag lingkungan.

⚠️ `[keputusan work owner]` **Alarm dihidupkan.** Di Pega, deteksi keadaan separuh dan pengiriman
email **keduanya ter-remark** di jalur endorsement — **endorsement hari ini berjalan tanpa alarm**,
sementara new business memilikinya. Sistem baru **menyeragamkan**.

Alamat Arasapas di-resolve **runtime** dari `M_LINK_SERVICE` (**ADR-0013**) — dilarang sebagai
literal, konstanta, maupun env var.

### 13. Pembacaan lintas skema Arasapas

`[terverifikasi]` Gerbang kelayakan ke-5 membaca **langsung** `ARASAPAS.DETAIL_INVOICE` dengan
penyaring `IVD_JR_ID = '5'`. `[keputusan work owner]` `IVD_JR_ID = '5'` berarti
**pembayaran/pelunasan**.

`[keputusan desain]` **Penyimpangan sadar.** Sistem baru **tetap membaca langsung** — mengubahnya
menjadi pemanggilan layanan adalah perubahan kontrak dengan pihak lain, di luar cakupan migrasi.
Syaratnya: pembacaan itu **dikurung dalam satu repository yang ditandai batas lintas sistem**, tidak
tersebar.

### 14. Penyaringan peserta untuk klaim

⚠️ `[keputusan work owner]` **Penajaman kontrak lintas konteks.** Peserta ber-`EDMStatus` `Batal`
atau `Delete` **tidak boleh muncul** di Claim Life. Kueri klaim
(`Claim Life/RDBList/GetPesertaClaim_sql1.xml`, `ASM-FW-GCNMFW-WORK-CLAIMLIFE` /
`RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL`, berkunci `PL_NUMBER`) **wajib menyaringnya keluar**.

**Akuntansi melihat seluruh baris; klaim hanya peserta yang masih hidup.** Satu tabel, dua sudut
pandang — dan itu **sengaja**.

#### Kolom penandanya `[terverifikasi]`

`M_LIFE_PREMIUM_DETAIL` punya **tiga** kolom berakhiran status; hanya **satu** menandai hidup/mati:

| Kolom | Diisi dari | Nilai | Penanda hidup/mati? |
| --- | --- | --- | --- |
| **`EDMSTATUS`** | `TempValue.EDMStatus` | `Old` / `New` / `Delete` / `Batal` | ✅ **ya** |
| `STATUS` | `CARI48` | `0` untuk `QR`/`QP`, `1` untuk `TP`/`TR` | ❌ penanda **jenis transaksi** |
| `STATUSOLD` | `CARI47` | `1`/`0` | ❌ |

Bukti: daftar kolom `INSERT` dan klausa `VALUES` di kedua ekspor `SaveMasterLPDet`.

⚠️ `[terverifikasi]` **Jalur new business tidak mengisi `EDMSTATUS`** — sensus
`InsertLifePremiumDetail_act` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT`): **nol**
kemunculan. Baris NB masuk dengan `EDMSTATUS` kosong/NULL dan **tetap peserta hidup**. Penyaring
naif `EDMSTATUS NOT IN ('Delete','Batal')` **membuang seluruh peserta NB** di Oracle.

**Aturan yang berlaku:** hidup = kosong/NULL, `Old`, atau `New`. Mati = `Delete` atau `Batal`.

`[terbuka]` Tipe dan nullability `EDMSTATUS` belum terbaca — **OQ-001 (sisa)**, pemilik **DBA**.
Tidak memblokir aturan bisnisnya, hanya bentuk akhir penyaringnya.

⚠️ Ini menyentuh **tiket 08 PremiumList Life** dan **spec Claim Life** — keduanya **sudah
diperbarui** 2026-09-15 (spec Claim Life §16 + AC 25–30; tiket 02 Claim Life; tiket 08 PremiumList
Life).

### 15. Kode mati yang tidak dimigrasikan

| Bagian | Identitas | Catatan |
| --- | --- | --- |
| `GetOldDetail_EDM` **sebagai rule** | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GETOLDDETAIL_EDM` | ⚠️ **perilaku popup tetap wajib ada** |
| `CreateCaseEMDL` step **7–11** | `DATA-PORTAL` / `CREATECASEEMDL` | penyalinan data lama yang sudah mati |
| `UploadCSVEDMLifePremium_Act` step **4–6** | `@BASECLASS` / `UPLOADCSVEDMLIFEPREMIUM_ACT` | hanya tiga langkah pertama hidup |
| `Commit` eksplisit & `Connect-REST` yang ter-remark | — | Go memegang transaksi |
| **`Calculate1_Act` di jalur endorsement** | `ASM-FW-GISFW-WORK-LIFE` / `CALCULATE1_ACT` | berkas identik, **tetapi tidak dijalankan** di sini |
| Batas **50.000 baris** | `SetPremi_EDM` | dibuang |

⚠️ **Pelajaran metodologi (OQ-066), berlaku lintas modul.** Dua penanda **tidak dapat dipercaya
sendirian**: (1) `<pyStepsBlockName>` — ada langkah aktif yang sebenarnya mati, dan langkah mati
yang justru dihidupkan kembali; (2) **kesamaan berkas antar folder** — `Calculate1_Act` identik di
kedua modul tetapi tidak dipakai di endorsement. Rujukan `Call` / `<RequestType>` pun hanya
menunjukkan **kemungkinan** pemanggilan. **Status hidup/mati menuntut konfirmasi work owner.**

---

### 16. Bentuk penyimpanan — **berbagi tabel PremiumList Life; endorsement adalah VERSI** ⚠️ BARU 2026-09-16

⚠️ **Penyimpangan sadar 8 — seluruh JSON dibuang.** `[keputusan work owner]` `InsertJsonPolisEDM`
dan `convertJsonNusareToProduction` **tidak dimigrasikan**; hilir membaca **langsung dari tabel**.

⚠️ **Penyimpangan sadar 9 — endorsement TIDAK punya tabel sendiri.** `[keputusan work owner]`
Ia memakai **tujuh tabel PremiumList Life yang sama** (spec PremiumList §12). Endorsement adalah
**baris versi baru** di `T_PREMIUM_LIST`; seluruh versi hidup berdampingan. Class work-nya memang
berbeda (`ASM-FW-GISFW-Work-EndorsementLife` vs `Work-LIFE`), tetapi **keduanya dipetakan ke tabel
yang sama**.

#### Pembeda versi — kolom pada `T_PREMIUM_LIST`

`[terverifikasi dari pyFields `DATA_JSON` EDM nyata]` Kolom EDM, seluruhnya **nullable** — kosong
pada baris new business, terisi saat endorsement:

| Kolom | Arti |
| --- | --- |
| `EDM_TYPE` | **`1`=Perubahan Data, `3`=Batal** — dropdown `EDMTYPE` di korpus. Nilai `2` memang tidak ada |
| `OLD_POLICY_NO` | nomor polis yang di-endorse |
| `EDM_DATE`, `EDM_NOTE` | tanggal & catatan endorsement |
| `TYPE_CEDING` | `1`=QS, `2`=SURPLUS, `3`=QS+SURPLUS, `4`=XOL — dropdown korpus |
| `PREMI_PROPOSED`, `UANG_PERTANGGUNGAN`, `SUM_INSURED` | **NUMBER** |
| `JENIS_PRODUK`, `SISTEM_REASURANSI` | teks |
| `STATUSS`, `STATUS_UPDATE`, `STATUS_SERVICE` | status proses endorsement |
| `START_DATE`, `END_DATE` | **DATE** |

Ditambah pembeda versi umum: `NOENDORS` `[terverifikasi]` (`Generate_NoEndorsmentLife`),
`PL_NUMBER_EDM`, **`PRODKE`** (⚠️ **versi berjalan = `PRODKE` terbesar**, `[terverifikasi]`),
`EDMSTATUS`, `STATUSOLD`, dan identitas kerja **`EDMLF-<n>`** `[fakta bisnis — work owner]`
(paritas `NBLF-`; **nol kecocokan di korpus** — diterima sebagai fakta bisnis).

#### ⚠️ Penyimpangan sadar 10 — pencocokan peserta lewat `PARENT_ID`, bukan indeks

`[terverifikasi]` Pega mencocokkan peserta lama↔baru dengan **loop indeks** di `MappingEDMLife` —
`PremiumListDetail(Local.idxLocationEDM)` dibandingkan `(Param.EDMList = .pxListSubscript)`. Cara itu
**rawan salah pasang** begitu urutan baris bergeser (tambah/hapus peserta), dan **kunci
pencocokannya tidak terbaca dari korpus**.

**Keputusan** `[keputusan work owner]`: kolom **`PARENT_ID`** — FK self-reference ke `ID` peserta
versi sebelumnya, **hanya di `T_PREMIUM_LIST_DETAIL`**. `NULL` untuk new business dan untuk peserta
yang baru ditambahkan. Pointer eksplisit **tidak pernah salah pasang**, tidak bergantung
`CERTIFICATE_NO` (bisa berubah) maupun urutan baris.

Hanya tingkat peserta yang di-selisih, jadi hanya ia yang butuh pointer:

| Tabel | Disalin saat endorsement? | `PARENT_ID`? |
| --- | --- | --- |
| `T_PREMIUM_LIST` (header) | ✅ versi baru + kolom EDM | opsional, untuk jejak versi |
| `T_PREMIUM_LIST_DETAIL` (peserta) | ✅ | ✅ **wajib** |
| `T_PREMIUM_LIST_SPREADING` / `_SPREADING_RETRO` | ✅ ikut peserta | ❌ cukup FK ke induk versi baru |
| `T_PREMIUM_LIST_SUMMARY` (rekap mata uang) | ❌ **tidak** — **dihitung ulang** dari peserta versi baru | ❌ |
| `T_VIEW_SUGGEST` | ❌ **tidak** `[keputusan work owner]` — riwayat penawaran milik proses awal | ❌ |

**Selisih new − old** dibaca dengan menyandingkan baris pada `ID` dan `PARENT_ID`-nya; peserta baru
(`PARENT_ID` `NULL`) tidak punya pengurang.

#### ⚠️ Penyimpangan sadar 11 — hapus = flag + nilai minus, **tidak pernah** hapus fisik

`[keputusan work owner]` Status peserta adalah **turunan** dari `PARENT_ID` + aksi:

| Kondisi | `EDMSTATUS` | Pengurang |
| --- | --- | --- |
| hasil salin, nilai diubah | `Change` | selisih `new − old` |
| peserta baru setelah salin | `New` | tidak ada (`PARENT_ID` `NULL`) |
| baris salin ditandai keluar | `Delete` | **pengurang penuh** (`−old`) |
| pembatalan seluruh polis | `Batal` | **pengurang penuh** untuk **setiap** peserta |

Ini **memperkuat** §6 (jurnal balik `× −1`): baris tetap tersimpan dan **auditable**; yang berubah
hanya tandanya. Semangat sama dengan Claim Life.

**`EdmType 1` vs `EdmType 3` memakai jalur simpan yang sama** — salin + `PARENT_ID`. Bedanya hanya
nilai dan flag: `1` mengubah sebagian peserta; `3` **me-minus-kan seluruhnya**.

#### `OldData` tidak disimpan

`[keputusan work owner]` `Addendum.OldData` di Pega adalah **halaman kerja sementara** untuk
menampilkan before/after — ia **tidak di-persist**. Data lama dibaca dari **versi sebelumnya**
(`PRODKE` lebih kecil, lewat `PARENT_ID`).

---

## Testing Decisions

### Apa yang membuat test baik di sini

Test memeriksa **perilaku yang terlihat dari luar**, bukan susunan internal. Ia menggerakkan alur
lewat endpoint dan memeriksa hasilnya lewat endpoint — bukan memanggil fungsi dalam.

Tiga hal di konteks ini **wajib** diuji lewat data nyata, bukan lewat mock: **tanda nilai uang**
(positif versus negatif), **isi `EDMStatus`**, dan **apa yang bertahan saat gagal di tengah**.
Ketiganya adalah inti kebenaran endorsement dan tidak dapat difake dengan jujur.

### Seam — **memakai ulang seam yang sudah ada**

`[terverifikasi]` Repo target belum di-scaffold. Seam yang ditetapkan spec Claim — Life adalah
**API HTTP**, dan konteks ini memakainya kembali — **tidak menambah seam**:

> **Seam utama: API HTTP Endorsement Life.** Test menggerakkan alur lewat endpoint REST dan
> memeriksa hasilnya lewat endpoint REST, dengan `handlers → services → repository` terpasang
> sungguhan, terhadap skema uji Oracle.

**Batas proses difake:**

| Batas | Perlakuan |
| --- | --- |
| Arasapas — **kirim** (efek keluar) | *fake* di balik interface; test memeriksa **efeknya** |
| Arasapas — **baca pembayaran** (gerbang ke-5) | **skema uji nyata** — ia gerbang bisnis, bukan sekadar integrasi |
| Email alarm | *fake*; test memeriksa terpicu atau tidak |
| Oracle | **skema uji nyata, bukan mock** — tanda uang dan `EDMStatus` adalah inti perilaku, dan sifat commit tiap penulis tidak dapat difake dengan jujur |
| Jam | dapat dikendalikan |

**Tidak ada seam kedua.** Konteks ini tidak punya worker asinkron.

### Modul yang diuji

| Yang diuji | Lewat seam |
| --- | --- |
| Lima gerbang kelayakan, satu per satu dan gabungan | API HTTP + skema uji |
| Pemuatan & penyalinan data polis lama; penandaan `Old` | API HTTP + skema uji |
| Popup polis lama per `Type` (`QR` induk; `QP`/`TP`/`TR` varian) | API HTTP |
| Dua maksud: Perubahan Data versus Batal | API HTTP |
| **Jurnal balik** — 32 kolom uang menjadi negatif, cakupan selektif versus menyeluruh | API HTTP + skema uji |
| Penomoran `PL_NUMBER_EDM` dan kenaikan `PRODKE` | API HTTP + skema uji |
| `Confirm` / `Decline`, dan **ketiadaan** `Reject` | API HTTP |
| Unggah CSV: ganti bukan tumpuk, tanpa batas baris, validasi konsistensi | API HTTP + skema uji |
| Kunci field permanen dan pesannya | API HTTP |
| Batas transaksi campuran — apa yang bertahan saat gagal di tengah | API HTTP + skema uji |
| Anti-dobel versi — hitung baris `T_PREMIUM_LIST` untuk `(NOPOLIS, PRODKE)` | API HTTP + skema uji |
| Efek keluar di belakang flag lingkungan; alarm menyala saat gagal | API HTTP + fake |
| **Penyaringan peserta batal/delete** dari jalur baca klaim | API HTTP + skema uji |

### Prior art

`[terverifikasi]` **Tidak ada** — nol kode, nol test di `OUTPUT_HASIL_RNM`. Spec Claim — Life,
Komite Claim Life, dan PremiumList Life menetapkan bentuknya; konteks ini mengikuti bentuk yang sama.

Perintah verifikasi wajib ditulis eksplisit di tiap tiket selama `Makefile` belum ada. Target:
`go test ./internal/...` dan `cd frontend && npm test`.

---

## Acceptance Criteria

**Gerbang kelayakan**

1. Endorsement ditolak bila **nomor polis kosong**, dengan pesan yang menyebutnya.
2. Endorsement ditolak bila **polis tidak ditemukan** di rekam polis.
3. Endorsement ditolak bila polis itu **masih punya endorsement yang belum selesai** — **satu
   endorsement terbuka per polis**, aturan integritas inti.
4. Endorsement ditolak bila polis **sudah pernah dibatalkan** (`EdmType` terakhir = `3`).
5. Endorsement **Batal** ditolak bila polis **sudah dibayar**. Endorsement **Perubahan Data** atas
   polis yang sudah dibayar **tetap boleh** — gerbang ini hanya berlaku bila `EdmType=3`.
6. Case endorsement dibuat **hanya setelah kelima gerbang lolos**; gerbang yang gagal **tidak**
   meninggalkan case separuh.

**Memuat polis lama**

7. Data polis new business termuat **otomatis** ke endorsement baru; baris warisannya bertanda
   **`Old`**.
8. ID Pega lama dan nomor polis lama **tersimpan** pada case endorsement.
9. Popup "lihat polis lama" **ada dan terisi** — `QR` memakai layar induk; `QP`, `TP`, `TR` memakai
   varian masing-masing. Keempat jenis didukung.
10. Baris polis new business **tidak pernah dihapus fisik** oleh proses apa pun di konteks ini.

**Maksud dan status**

11. `EdmType` hanya menerima **`1`** (Perubahan Data) dan **`3`** (Batal); nilai lain ditolak.
12. `EDMStatus` hanya menerima **`Old`**, **`New`**, **`Delete`**, **`Batal`**.
13. Baris **`New`** hanya dapat lahir pada `EdmType=1`; pada `EdmType=3` penambahan peserta
    **ditolak**.
14. `EdmType=3` **otomatis** menandai **seluruh** peserta polis sebagai `Batal` — tanpa penandaan
    satu per satu.
15. `Delete` menandai **hanya peserta yang dipilih**.

**Jurnal balik**

16. Baris `Delete` dan `Batal` menghasilkan nilai **negatif**, **bukan nol** dan **bukan
    penghapusan**.
17. **Seluruh 32 kolom uang** ikut dibalik tandanya; tidak ada komponen yang tertinggal positif.
    Test memeriksa kolom demi kolom.
18. Baris positif asli **tetap ada** berdampingan dengan baris negatif di tabel yang sama.
19. **Jumlah** baris positif dan negatif untuk peserta yang dibatalkan = **nol**.
20. Cakupan minus benar: `Delete` menegatifkan **hanya peserta terpilih**; `Batal` menegatifkan
    **seluruh peserta**.
21. Nilai uang **tidak** melewati `float` di lapisan mana pun maupun di JSON API — termasuk nilai
    negatif. (**ADR-0003**)
22. Nilai yang ditulis dan dibaca kembali **identik**; tidak ada pembulatan diam pada nilai negatif.

**Penomoran**

23. Endorsement memperoleh `PL_NUMBER_EDM` dirakit **`<nomor polis>/<PRODKE+1>`** — **bukan** dari
    `PROC_GENERATE_SEQUENCE_NUMBER`.
24. Nomor endorsement **lahir sekali**; pemanggilan ulang tidak mengubahnya dan tidak menaikkan
    `PRODKE`.
25. Dua endorsement berurutan atas polis yang sama memperoleh `PRODKE` **berurutan**.
26. `PL_NUMBER` new business **tidak berubah** oleh endorsement apa pun.
27. `PRODKE` dibaca dengan **`ORDER BY PRODKE DESC`** di **setiap** tempat; ada test yang gagal bila
    ada jalur yang memakai urutan lain.

**Alur keputusan**

28. `Confirm` melanjutkan endorsement ke penyimpanan.
29. `Decline` **menutup** case; case tertutup tidak dapat dilanjutkan maupun diputuskan ulang.
30. **Tidak ada** keluaran `Reject` di konteks ini — test yang menemukan jalur "kembali ke input"
    **gagal**.
31. Setelah `Decline`, polis yang sama **dapat di-endorse ulang** — gerbang ke-3 tidak lagi menolak.
32. Setiap transisi tahap menulis **jejak audit**: siapa, kapan, dari tahap apa ke tahap apa.
    (**ADR-0007**)

**Kunci field**

33. Nomor polis dan jenis endorsement/batal **terkunci permanen** setelah case dibuat; percobaan
    mengubahnya ditolak di sisi server, bukan hanya di layar.
34. Layar **menampilkan pesan penjelas** mengapa field terkunci — bukan sekadar mematikannya.

**Unggah CSV**

35. Unggah ulang **mengganti** baris `New` sebelumnya, tidak menumpuk.
36. **Tidak ada batas jumlah baris**; unggahan besar diproses dan **tidak** ditolak karena
    jumlahnya. Test memuat berkas di atas 50.000 baris.
37. Baris dengan `PLAN` atau `POLICY_HOLDER` berbeda dari **baris pertama** ditolak, dengan pesan
    yang menyebut nomor baris dan kolomnya.
38. Nilai uang di CSV di-parse dengan **format pemisah yang dinyatakan eksplisit**, langsung ke
    desimal presisi arbitrer. Test memuat kasus pemisah ribuan.

**Transaksi**

39. Tiap penulis yang commit sendiri diperlakukan sebagai **titik potong**; data yang harus atomik
    tidak dipisahkan olehnya.
40. Kegagalan di tengah meninggalkan keadaan yang **terdeteksi dan dapat dipulihkan** — bukan senyap.
41. Menulis rekam polis dua kali untuk `(NOPOLIS, PRODKE)` yang sama **tidak** menggandakan baris —
    dibuktikan dengan memanggilnya dua kali lalu menghitung baris.
42. Urutan pemanggilan terdokumentasi di kode sebagai bagian kebenaran; ada test yang gagal bila
    urutannya diubah.

**Efek keluar**

43. Di lingkungan non-production, efek keluar **tidak berjalan**; penyimpanan **tetap** berjalan.
44. Alamat Arasapas di-resolve **runtime** dari `M_LINK_SERVICE`; test yang menemukan URL sebagai
    literal, konstanta, atau env var **gagal**. (**ADR-0013**)
45. **Alarm menyala** bila pembacaan balik menunjukkan penyimpanan gagal — perilaku **sama** dengan
    jalur new business.
46. Kegagalan efek keluar **tidak** membatalkan endorsement yang sudah tersimpan; ia tercatat,
    terlihat, dan dapat diulang.
47. Pembacaan pembayaran Arasapas terkurung di **satu** repository bertanda batas lintas sistem;
    test yang menemukan kueri skema Arasapas di luar itu **gagal**.

**Kontrak hilir**

48. Peserta ber-kolom **`EDMSTATUS`** bernilai **`Batal`** atau **`Delete`** **tidak muncul** pada
    jalur baca klaim.
48a. Peserta **new business** — yang `EDMSTATUS`-nya **kosong/NULL** karena jalur NB tidak mengisi
    kolom itu — **tetap muncul**. ⚠️ Penyaring naif `EDMSTATUS NOT IN ('Delete','Batal')` membuang
    seluruh peserta NB di Oracle; test wajib memuat kasus ini dan **harus gagal** bila penyaringnya
    naif.
48b. `STATUS` dan `STATUSOLD` **tidak** dipakai sebagai penanda hidup/mati — `STATUS` adalah penanda
    jenis transaksi (`0` untuk `QR`/`QP`, `1` untuk `TP`/`TR`).
49. Peserta yang sama **tetap terlihat** pada jalur baca akuntansi, lengkap dengan baris negatifnya.
50. Rekam premium hasil endorsement dapat ditemukan lewat **`PL_NUMBER_EDM`** maupun **`PL_NUMBER`**.
51. Perubahan bentuk rekam detail/summary diperlakukan sebagai **perubahan kontrak lintas konteks**
    dan ditandai demikian di kode. (**ADR-0001**)

**Kode mati**

52. Tidak ada padanan: rule pemuat data lama yang mati, langkah penyalinan mati di pembuat case,
    langkah unggah CSV yang mati, `Commit` eksplisit, maupun `Connect-REST` yang ter-remark.
53. **Mesin hitung milik new business tidak dipanggil** di jalur endorsement.

**Bentuk penyimpanan dan versi** ⚠️ BARU 2026-09-16 — §16

54. ⚠️ **Tidak ada blob JSON** di jalur endorsement; `InsertJsonPolisEDM` **tidak dipanggil** dan
    tidak ada padanan perakit payload. *(§16; penyimpangan sadar 8)*
55. ⚠️ Endorsement **tidak membuat tabel sendiri** — ia baris **versi baru** di tabel polis yang
    sama. Test yang menemukan tabel header khusus endorsement **gagal**. *(§16; penyimpangan
    sadar 9)*
56. Seluruh versi sebuah polis **hidup berdampingan**; versi berjalan adalah yang ber-`PRODKE`
    **terbesar**. *(§16)*
57. Kolom EDM **kosong** pada baris new business dan **terisi** pada baris endorsement — keduanya
    dalam tabel yang sama. *(§16)*
58. `EDM_TYPE` hanya menerima **`1`** (Perubahan Data) dan **`3`** (Batal); nilai lain **ditolak**.
    *(§16; `[terverifikasi]` dropdown `EDMTYPE`)*
59. ⚠️ Menekan Submit menyalin **header, seluruh peserta, seluruh spreading, dan seluruh spreading
    retro** ke versi baru — masing-masing dengan **identitas baru**. *(§16)*
60. ⚠️ Setiap peserta hasil salin membawa **`PARENT_ID`** yang menunjuk peserta versi sebelumnya.
    *(§16; penyimpangan sadar 10)*
61. ⚠️ Peserta yang **baru ditambahkan** setelah penyalinan ber-`PARENT_ID` **`NULL`**, dan **tidak
    punya pengurang**. *(§16)*
62. ⚠️ Pencocokan peserta **tidak** bergantung pada urutan baris maupun `CERTIFICATE_NO`. Test yang
    menyisipkan dan menghapus peserta sehingga urutan bergeser **tetap** menghasilkan pasangan yang
    benar. *(§16; penyimpangan sadar 10 — inilah kelas bug yang diperbaiki)*
63. ⚠️ **Rekap mata uang TIDAK disalin** — ia **dihitung ulang** dari peserta versi baru. Test yang
    menemukan baris rekap tersalin **gagal**. *(§16)*
64. ⚠️ **Riwayat penawaran tidak disalin** ke versi endorsement. *(§16)*
65. ⚠️ Selisih `new − old` per peserta dihitung dengan menyandingkan baris pada `PARENT_ID`-nya,
    bukan pada indeks. *(§16)*
66. ⚠️ Peserta yang dikeluarkan ditandai **`Delete`** dan nilainya menjadi **pengurang penuh** —
    barisnya **tetap tersimpan**. Test yang menemukan penghapusan fisik **gagal**. *(§16;
    penyimpangan sadar 11)*
67. ⚠️ `EdmType` **`3`** (Batal) menjadikan **setiap** peserta pengurang penuh, menandainya
    **`Batal`**, tanpa menghapus satu baris pun. *(§16; penyimpangan sadar 11)*
68. `EdmType` `1` dan `3` memakai **jalur simpan yang sama**; yang berbeda hanya nilai dan flag.
    *(§16)*
69. ⚠️ `OldData` **tidak disimpan**; data lama dibaca dari **versi sebelumnya**. Test yang menemukan
    tabel/kolom penyimpan `OldData` **gagal**. *(§16)*
70. ⚠️ Satu endorsement ditulis dalam **satu transaksi** — penyalinan, kolom EDM, dan rekap yang
    dihitung ulang; kegagalan di mana pun **membatalkan seluruhnya**. *(§11)*
71. ⚠️ Menjalankan penyalinan **dua kali** untuk maksud endorsement yang sama **tidak** menghasilkan
    dua versi — pemeriksaan `(NOPOLIS, PRODKE)` menahannya. *(§11)*
72. Identitas kerja endorsement berbentuk **`EDMLF-<n>`**, sejajar `NBLF-<n>` pada new business.
    *(§16; `[fakta bisnis — work owner]`)*

---

## Pertanyaan terbuka di dalam spec

**Nol OQ pemblokir.** `[terverifikasi]` OQ-070, OQ-071, dan OQ-072 ditutup 2026-09-15.

| OQ | Pertanyaan | Pemilik | Bagian yang menunggu |
| --- | --- | --- | --- |
| ~~**OQ-001**~~ (sisa) | ~~DDL fisik `M_LIFE_PREMIUM_SUMMARY`, `M_LIFE_PREMIUM_DETAIL`, `JSON_POLIS`~~ | — | ✅ **DITUTUP 2026-09-16** — tabel target **dirancang sendiri** bersama PremiumList Life (§16). ⚠️ Pertanyaan **nilai uang negatif** terjawab oleh rancangan: kolom uang **desimal presisi arbitrer bertanda**, dan jurnal balik memang menulis nilai negatif (§6, §16) |
| **OQ-020** (sisa) | Arti `QP` dan `QR` | Product+UW | **TETAP TERBUKA.** ⚠️ **KOREKSI 2026-09-16:** sempat ditandai ditutup "dari dropdown `TYPE` korpus EDM Life (QR=Receivable/QP=Payable)" — itu **KELIRU**. Sensus: `Receivable`/`Payable` **NOL berkas** di `Endorsement Life/`; label hanya di `pyLocalizedValue` payload `DATA_JSON` runtime (instance), **bukan rule korpus**. Yang terbukti cuma perilaku (`Q*`→gross, `T*`→retrosesi). Arti QP/QR belum terverifikasi |
| `[terbuka]` | **`NOENDORS`** — format & kontrak `Generate_NoEndorsmentLife`; **`EDMLF-`** — format & generator | DBA / work owner | penomoran; **tidak memblokir** |
| **OQ-066** | Penanda korpus tidak dapat dipercaya sendirian — pelajaran metodologi | Arsitektur Pega + Product+UW | §15; sudah diterapkan sebagai kehati-hatian |
| **OQ-028** | Seluruh Assignment memakai `WorkList` (bukan `WorkBasket`) — implikasi routing/RBAC | IAM | §2; belum menyentuh AC |

---

## Out of Scope

- **PremiumList Life (new business) — seluruhnya.** Konteks/menu terpisah; spec dan tiketnya di
  `.scratch/premiumlist-life/`. Mesin bersama (`SaveMasterLPDet`, `InsertPLSummary`, tabel premium)
  dipakai **sebagai komponen** dari sini, tetapi alur penawaran dan premium list new business bukan
  cakupan spec ini.
- **Claim Life dan Komite Claim Life.** Konteks hilir. ⚠️ **Kecuali** satu hal: syarat penyaringan
  peserta batal/delete (§14) adalah **kontrak** yang wajib dipenuhi sisi klaim — dan karena itu
  **tiket 08 PremiumList Life** dan **spec Claim Life** perlu diperbarui di luar spec ini.
- **Lini non-Life.** Treaty, facultative, dan claim non-life punya konteks sendiri.
- **Mengubah pembacaan Arasapas menjadi layanan.** Perubahan kontrak dengan pihak lain; di luar
  cakupan migrasi.
- **Merapikan alur.** Urutan tahap dan bentuk keputusan dibawa apa adanya; tujuh penyimpangan sadar
  sudah didaftar di §Solution.
- **Identity & Access.** `[terverifikasi]` ABSENT dari korpus; dibangun dari nol.
- **Scaffolding kode.** Pekerjaan terpisah yang mendahului tiket mana pun.
- **Migrasi data endorsement lama.** Menunggu DDL fisik (OQ-001 sisa).

---

## Further Notes

**Ukuran pekerjaan.** `[terverifikasi]` Endorsement Life: **75 berkas** — 16 Activity, 17 RDBList,
14 Section, 8 Harness, 6 FlowAction, 6 ReportDefinition, 2 DataTransform, 2 When, 1 Flow,
1 DecisionTable, 1 ConnectREST, 1 SystemSettings.

**Dua deskripsi Pega yang berbohong.** Konteks ini memuat dua nama/deskripsi yang bertentangan
dengan kodenya, dan keduanya berbahaya bila dipercaya:

| Yang tertulis | Yang sebenarnya |
| --- | --- |
| Rule bernama "SetErrorBatalEndorsement" | **gerbang kelayakan pembuatan** endorsement, bukan pembatalan |
| Langkah berdeskripsi "Set 0 jika EDM Batal" | mengalikan **32 kolom uang dengan `-1`** |

Keduanya sudah dikoreksi di spec ini. **Pelajaran yang berlaku lintas modul: baca kodenya, jangan
namanya.**

**Urutan yang saya sarankan untuk `/to-tickets`** — vertical slice:

1. Cari polis + **lima gerbang kelayakan** (termasuk pembacaan Arasapas terkurung).
2. Buat case + muat & salin data polis lama + penandaan `Old`.
3. Popup lihat polis lama per `Type`.
4. Penomoran `PL_NUMBER_EDM` + `PRODKE` satu urutan.
5. Perubahan Data: tambah/hapus peserta + `EDMStatus`.
6. **Jurnal balik**: `Delete` selektif dan `Batal` menyeluruh.
7. Unggah CSV tanpa batas baris + validasi konsistensi.
8. Alur `Confirm`/`Decline` + kunci field permanen + jejak audit.
9. Simpan: batas transaksi campuran + anti-dobel `(NOPOLIS, PRODKE)`.
10. Efek keluar + alarm dihidupkan.
11. **Kontrak hilir**: penyaringan peserta batal/delete dari jalur baca klaim.
12. Migrasi skema — **`needs-info`** menunggu OQ-001, pola sama dengan tiket 13 Claim Life dan
    tiket 09 PremiumList Life.

**Keputusan yang masih terbuka, tidak memblokir:** pembagian paket domain di dalam `internal/` —
sama seperti tiga spec sebelumnya.

**Catatan sumber.** Spec ini bersandar **hanya** pada korpus Pega `D:\XML\RNM_BRD\` dan artefak di
`OUTPUT_HASIL_RNM\`. Sumber ADR tunggal: `docs/adr/ADR-0001`…`ADR-0015`.
