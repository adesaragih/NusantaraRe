# Spec — Komite Claim Life (migrasi Pega → Go + React + Oracle)

Status: ready-for-agent
Konteks: `komite-claim-life` (Komite Claim Life)
Tanggal: 2026-09-15
Sumber: grilling Ronde 1–2 (`grilling-ronde-1.md`, `grilling-ronde-2.md`), `CONTEXT.md`,
`docs/adr/ADR-0001`–`ADR-0015`, `discovery/modules/Komite Claim Life.md`,
`discovery/flows/Komite Claim Life.md`, `discovery/inventory/Komite Claim Life.md`,
`discovery/context-map.md` §2.7, `discovery/open-questions.md`
Skill: `/mattpocock-skills:to-spec`

> **Konvensi penandaan.** Setiap pernyataan perilaku membawa salah satu dari:
> `[terverifikasi]` — terbaca di korpus, disertai path + rule;
> `[keputusan work owner]` — dinyatakan work owner, **tidak** terbukti korpus;
> `[data DBA]` — dari dump DBA; `[terbuka]` — OQ dari register.
> Tidak ada pernyataan tanpa penanda.

> **Sumber tunggal.** Korpus Pega `D:\XML\RNM_BRD\` (READ-ONLY) dan artefak di
> `OUTPUT_HASIL_RNM\`, yang sekaligus **repo target tunggal**. ⛔ `D:\XML\nusantara-re\`
> di-blacklist — tidak dibaca, tidak dirujuk.

> ⚠️ **Catatan audit.** `Komite Claim Life/Activity/KomitePostAdjustment.xml` **diperbarui work
> owner 2026-09-15** (515.675 byte, stempel simpan `20260915T023136`). Seluruh nomor langkah di spec
> ini dari versi itu, dibaca lewat `<pyStepPageReference>`.

---

## Problem Statement

Keputusan akseptasi klaim jiwa bernilai besar tidak diambil satu orang — ia menaiki **tangga
persetujuan komite**. Hari ini tangga itu berjalan di Pega, dan bentuknya menyulitkan siapa pun yang
ingin memahaminya:

1. **Tangganya tidak terlihat di proses.** `[terverifikasi]` Grafnya hanya **1 Assignment + 1
   Decision + 4 connector** — satu assignment yang di-loop. Berapa tingkat dan siapa penyetujunya
   ditentukan **data**, bukan struktur. Membaca flow tidak memberi tahu apa pun tentang tangganya.
2. **Wewenang tidak ditegakkan.** `[terverifikasi]` Pega hanya *menempatkan* kasus di antrean
   pemilik `KomiteID`. Tidak ada satu pun pemeriksaan bahwa yang menyimpan keputusan adalah orang
   itu. Keputusan masuk lewat dropdown wajib di layar, dan **tidak ada rule yang menulisnya**.
3. **Efek keluar bisa gagal diam-diam** — termasuk **Kasir**, jalur pembayaran. Klaim dapat disetujui
   tanpa pernah sampai ke pembayaran, dan tidak ada yang tahu.
4. **Kode mati bercampur kode hidup.** `[terverifikasi]` Lima langkah ter-remark, satu cabang routing
   yang tidak pernah menyala, satu gerbang EXIT dengan tiga identitas ter-hardcode, dan dua blok tulis
   kembar sepanjang ribuan baris. Membedakan yang berlaku dari yang tidak menuntut penelusuran
   langkah demi langkah.

## Solution

Membangun **Komite Claim Life** sebagai konteks tersendiri di layanan Go + antarmuka React, yang
menulis ke Oracle `POOLDATA` yang sama, dengan:

- **Tangga persetujuan yang dinyatakan eksplisit** — jumlah tingkat dihitung dari roster, bukan dari
  struktur proses.
- **Wewenang ditegakkan di lapisan layanan** (**ADR-0014**), dengan satu pengecualian sah: eskalasi
  manual naik satu tingkat.
- **Efek keluar wajib berhasil** lewat transactional outbox (**ADR-0015**) — menyimpang dari
  **ADR-0008** karena Kasir memindahkan uang.
- **Satu jalur simpan berparameter status**, menggantikan dua blok kembar.
- **Kode mati tidak ikut pindah** — `TransferType`, gerbang EXIT retro, dan lima langkah ter-remark
  dibuang.

Perilaku bisnis dibawa apa adanya (paritas), kecuali **tiga penyimpangan sadar** yang sudah
diputuskan: penegakan wewenang (**ADR-0014**), jaminan efek keluar (**ADR-0015**), dan penyatuan blok
tulis.

### Batas dengan Claim — Life

Konteks ini adalah **sisi lain** dari kontrak yang sudah dibangun di Claim — Life (**ADR-0001**):

| Arah | Kontrak | Tiket Claim Life |
| --- | --- | --- |
| **Masuk** | Claim Life menyerahkan baris `AdjustmentList` sebagai *child work* berkelas `ASM-FW-GCNMFW-Work-KomiteLife`, membawa `KomiteLoop`, `KomiteList`, penunjuk baris, nilai klaim, `CURRENCY`, dan `STS_REJECT` saat penyerahan | **tiket 10** |
| **Keluar** | Hasil keputusan memantul ke baris `AdjustmentList` sebagai `STS_REJECT` (`1` aksep / `2` tolak), hanya pada tingkat terakhir | **tiket 11** |
| **Bersama** | `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` — `[terverifikasi]` ditulis **satu rule yang sama** dari kedua sisi, hash ternormalisasi `c50bfd9a12` identik | — |

`[terverifikasi]` Komite **bukan** lapisan tipis di atas Claim: **59 dari 133 identitas** modul
Komite (**44,4 %**) tidak muncul di modul Claim mana pun (`discovery/inventory/_summary.md` §15.4).
Ia punya class kerja sendiri, model data sendiri (`KomiteList`), dan roster sendiri.

---

## User Stories

### Menerima dan merutekan kasus

1. Sebagai **anggota komite**, saya ingin kasus yang menunggu keputusan saya muncul di inbox saya,
   supaya saya tahu apa yang harus saya kerjakan. `[terverifikasi]`
   `Komite Claim Life/Flow/KomiteLife_Flow.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITELIFE_FLOW`
   / `RULE-OBJ-FLOW`) — Assignment "KomiteRouter", `pyImplementation = WorkList`,
   `pyRouteTo = Custom`.
2. Sebagai **anggota komite**, saya ingin inbox saya **hanya** berisi kasus sesuai posisi saya di
   roster, supaya saya tidak melihat pekerjaan orang lain. `[keputusan work owner]` (**ADR-0014**).
3. Sebagai **organisasi**, saya ingin kasus dirutekan ke **baris roster pertama yang belum
   menyetujui**, supaya tangga bergerak sendiri tanpa perlu satu shape per tingkat.
   `[terverifikasi]` `Komite Claim Life/Activity/KomiteRouter.xml`
   (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY`) —
   `param.AssignTo = .KomiteID` (baris ~294) bergerbang `.KomiteAproval == 0` (baris ~382).
4. Sebagai **organisasi**, saya ingin jumlah tingkat tangga **dihitung dari roster**, bukan dari
   konstanta, supaya perubahan kebijakan cukup dilakukan di data. `[terverifikasi]`
   `Claim Life/Activity/CreateKMTLife_Act.xml` —
   `childPageKomite.KomiteLoop = @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)`.
5. Sebagai **organisasi**, saya ingin roster dipilih berdasarkan **pita nilai klaim**, supaya klaim
   besar menaiki tangga yang lebih panjang. `[terverifikasi]`
   `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml`
   (`ASM-FW-GCNMFW-INT-EMAILKOMITE` / `FILTEREMAILKOMITEWITHLIMIT` / `RULE-OBJ-REPORT-DEFINITION`) —
   `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM AND .STS_KLAIM = Param.STS_KLAIM AND .STS_AKTIF = "1"`.
6. Sebagai **organisasi**, saya ingin klaim bernilai **negatif** tetap menemukan rosternya, supaya
   pembalikan nilai tidak melewati tangga persetujuan. `[terverifikasi]`
   `Claim Life/Activity/GetListKomiteLife.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
   `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`) —
   `Param.LIMIT_BOTTOM = @if(Local.IsADj<0, Local.IsADj * -1, Local.IsADj)`.

### Memutuskan

7. Sebagai **anggota komite**, saya ingin melihat rincian klaim beserta baris adjustment-nya sebelum
   memutuskan, supaya keputusan saya berdasar. `[terverifikasi]` FlowAction `ViewTransferDtl`
   (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `VIEWTRANSFERDTL` / `RULE-OBJ-FLOWACTION`).
8. Sebagai **anggota komite**, saya ingin memilih **Setuju** atau **Tolak**, supaya keputusan saya
   tercatat. `[keputusan work owner]` enum tertutup `{1 = Setuju, 2 = Tolak}`; `[terverifikasi]`
   kontrol dropdown wajib di `Komite Claim Life/Section/ShowTransfer.xml`
   (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION`) baris 32607 —
   `pyValue = .AcceptStatus`, `pyFormat = pxDropdown`, `pyRequired = true`.
9. Sebagai **anggota komite**, saya ingin menuliskan komentar atas keputusan saya, supaya alasannya
   terbaca tingkat berikutnya. `[terverifikasi]` `KomitePostAdjustment` menulis
   `KomiteList(KomiteCount).KomiteComment = pyWorkPage.Comment`.
10. Sebagai **anggota komite**, saya ingin keputusan saya tercatat beserta **waktunya**, supaya
    riwayat tangga dapat ditelusuri. `[terverifikasi]`
    `KomiteList(KomiteCount).DateApprove = @CurrentDateTime()`.
11. Sebagai **organisasi**, saya ingin nilai keputusan di luar `{1, 2}` **ditolak terang-terangan**,
    supaya tangga tidak berhenti diam-diam. `[keputusan work owner]`; `[terverifikasi]`
    `Komite Claim Life/When/IsKomiteLoop.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `ISKOMITELOOP` /
    `RULE-OBJ-WHEN`) hanya melanjutkan pada `"1"`.
12. Sebagai **organisasi**, saya ingin **hanya pemilik `KomiteID` tingkat berjalan** yang dapat
    menyimpan keputusan, supaya wewenang tidak dapat dilewati lewat API. `[keputusan work owner]`
    (**ADR-0014**) — penyimpangan sadar.

### Menaiki tangga

13. Sebagai **organisasi**, saya ingin tangga berlanjut selama keputusan **Setuju** dan tingkat belum
    habis, supaya persetujuan berjenjang berjalan. `[terverifikasi]` `IsKomiteLoop`:
    `.AcceptStatus = "1"` **DAN** `.KomiteCount <= .KomiteLoop`.
14. Sebagai **organisasi**, saya ingin tangga **berhenti** saat keputusan Tolak, supaya klaim yang
    ditolak tidak diteruskan ke atas. `[terverifikasi]` guard yang sama — `AcceptStatus = "1"` adalah
    syarat lanjut.
15. Sebagai **organisasi**, saya ingin naik tingkat **tanpa** membuat nomor akseptasi, supaya nomor
    hanya lahir sekali. `[terverifikasi]` di tingkat bukan-terakhir hanya
    `KomiteAproval = AcceptStatus` yang dicatat (baris ~792, ~908), lalu `KomiteCount + 1`
    (baris ~9020).
16. Sebagai **admin komite**, saya ingin memindahkan kasus **naik satu tingkat** bila anggota tingkat
    berjalan berhalangan, supaya kasus tidak macet. `[keputusan work owner]` (**ADR-0014**).

### Keputusan final

17. Sebagai **organisasi**, saya ingin rekam akseptasi dibuat **sekali**, di tingkat terakhir, supaya
    tidak ada nomor ganda. `[terverifikasi]` kedua blok tulis digerbangi `KomiteCount == KomiteLoop`
    (baris 5695 dan 8119).
18. Sebagai **ReasLifeSPV** di Claim — Life, saya ingin hasil keputusan komite memantul ke baris
    `AdjustmentList` saya, supaya status di sistem saya mencerminkan keputusan terakhir.
    `[terverifikasi]` `KomitePostAdjustment` menulis `STS_REJECT` pada **dua tingkat baris**
    (`PremiumListDetail` dan `AdjustmentList`) — kontrak **tiket 11** Claim Life.
19. Sebagai **organisasi**, saya ingin nomor akseptasi dibuat lewat jalur yang sama dengan Claim
    Life, supaya penomoran konsisten di kedua sisi. `[terverifikasi]` rantai aktif **4.7**
    `GetKodeProdLife_SQL` → **4.9** `GetSequenceNumber_SQL` → **4.11/4.12** cabang `QR,QP` / `TR,TP`
    → **4.15** `Insert ke OS` (**ADR-0006**).
20. Sebagai **organisasi**, saya ingin dokumen akseptasi tercetak saat disetujui, supaya ada bukti
    keputusan. `[terverifikasi]` step **4.17** `Call PrintAkseptasiPDF`, **4.18**
    `Call LoadDocumentLife_ACT`.
21. Sebagai **Finance**, saya ingin nilai uang pada rekam akseptasi tidak kehilangan presisi, supaya
    angkanya dapat dipertanggungjawabkan. `[data DBA]` kolom uang Oracle `NUMBER` tanpa presisi →
    desimal presisi arbitrer (**ADR-0003**).

### Efek keluar

22. Sebagai **organisasi**, saya ingin keputusan final terkirim ke sistem hilir, supaya klaim yang
    disetujui benar-benar diproses. `[terverifikasi]` empat efek setelah `Obj-Save` (step 6).
23. Sebagai **Finance**, saya ingin klaim yang disetujui **pasti** sampai ke Kasir, supaya pembayaran
    tidak hilang diam-diam. `[keputusan work owner]` (**ADR-0015**) — menyimpang dari **ADR-0008**.
24. Sebagai **operator sistem**, saya ingin efek yang gagal **diantre ulang sampai berhasil**, supaya
    tidak ada kiriman yang menguap. `[keputusan work owner]` — transactional outbox, at-least-once.
25. Sebagai **Finance**, saya **tidak** ingin pembayaran terkirim dua kali, supaya tidak ada dobel
    bayar. `[keputusan work owner]` — ID idempoten unik + **cek status sukses sebelum kirim ulang**
    khusus Email dan Kasir.
26. Sebagai **anggota komite**, saya ingin melihat kasus yang efek keluarnya **perlu intervensi**,
    supaya kegagalan tidak tersembunyi di antrean. `[keputusan work owner]` — status eksplisit di UI
    Komite + laporan harian.
27. Sebagai **operator sistem**, saya ingin alamat layanan keluar dibaca dari tabel saat dipanggil,
    supaya pemisahan dev–prod ditentukan isi database. `[terverifikasi]` kunci `Klaim` /
    `insertClaimLife` (**ADR-0013**).
28. Sebagai **organisasi**, saya ingin keputusan komite tersimpan **sebelum** efek keluar berjalan,
    supaya kegagalan jaringan tidak membatalkan keputusan. `[terverifikasi]` keempat efek berada
    setelah `Obj-Save` (step 6).

### Kejelasan dan jejak

29. Sebagai **pengguna mana pun**, saya ingin melihat riwayat lengkap tangga — siapa, kapan,
    keputusan apa, komentar apa — supaya saya paham mengapa klaim berada di keadaannya sekarang.
    `[terverifikasi]` entri `KomiteList(KomiteCount)` per tingkat.
30. Sebagai **auditor**, saya ingin perpindahan kasus naik tingkat lewat eskalasi tercatat pelakunya,
    supaya pengetatan wewenang tidak dilubangi diam-diam. `[keputusan work owner]` (**ADR-0014**,
    **ADR-0007**).
31. Sebagai **auditor**, saya ingin perubahan roster tercatat, supaya perubahan **siapa yang
    berwenang** dapat ditelusuri. `[keputusan work owner]` (**ADR-0014**).
32. Sebagai **pengguna mana pun**, saya ingin status ditampilkan sebagai kata, bukan nama field
    `STS_REJECT`, supaya saya tidak salah paham. `[terverifikasi]` nilai `1` berarti **diaksep**
    meski namanya "reject".

### Yang sengaja tidak dibawa

33. Sebagai **tim migrasi**, saya ingin cabang routing `TransferType` **tidak** ikut pindah, supaya
    kode mati tidak menular. `[terverifikasi]` `TransferType` **tidak pernah diisi** — nol
    `Property-Set` di Claim Life maupun Komite Claim Life; `TransferType == '2'` selalu FALSE.
34. Sebagai **tim migrasi**, saya ingin layar `ShowTransfer` yang bergantung `TransferType` ikut
    dibuang. `[keputusan work owner]`; `[terverifikasi]` visible-when `.TransferType==2`
    (baris ~29177).
35. Sebagai **tim migrasi**, saya ingin gerbang EXIT retro **tidak** direplikasi, supaya tidak ada
    klaim yang diam-diam melewati efek keluar. `[keputusan work owner]` (**OQ-064**).
36. Sebagai **tim migrasi**, saya ingin lima langkah ter-remark **tidak** ikut pindah.
    `[terverifikasi]` 4.4, 4.5, 4.6, 4.14, 5.5 ber-`<pyStepsBlockName>//`.
37. Sebagai **tim migrasi**, saya ingin nilai retro dipakai **apa adanya** dari data policy, tanpa
    logika penukaran. `[keputusan work owner]` (**OQ-065**).
38. Sebagai **tim migrasi**, saya ingin dua blok tulis kembar **disatukan**, supaya tidak ada dua
    tempat yang harus diubah bersamaan. `[terverifikasi]` diff struktural kedua blok: **nol beda**.
39. Sebagai **tim migrasi**, saya ingin `SetInformationData` **tidak** diperlakukan sebagai efek
    keluar. `[keputusan work owner]` — temporary penampung hasil submit.
40. Sebagai **tim migrasi**, saya ingin dua rule penomoran lama **tidak** dimigrasikan.
    `[terverifikasi]` `Generate_NoAccept_KMT_Life` / `_LifeRetro` — `RequestType` 1, terindeks
    `pyRuleName` **0**, dan ber-`blockname //`. Dua sinyal bebas menunjuk hal yang sama.

### Penyimpanan ⚠️ BARU 2026-09-16 — §9

41. Sebagai **organisasi**, saya ingin **roster dan keputusan tiap tingkat tersimpan di tabel
    relasional** yang dapat ditelusuri **per anggota**, supaya riwayat persetujuan dapat
    dipertanggungjawabkan dan tidak hilang bersama halaman kerja. `[keputusan work owner]`
42. Sebagai **pengguna**, dari sebuah baris adjustment saya ingin **langsung tahu kasus komite mana**
    yang memutuskannya — lewat `KOMITE_ID` — supaya saya tidak perlu menelusuri balik.
    `[keputusan work owner]`

---

## Implementation Decisions

### 1. Batas konteks

Komite Claim Life dispesifikasikan sebagai konteks tersendiri. **Claim — Life adalah konteks luar**,
dihubungkan tiga kontrak (**ADR-0001**) — lihat §"Batas dengan Claim — Life" di atas. Spesifikasi ini
berhenti di ketiga batas itu.

### 2. Penempatan modul

`[terverifikasi]` Struktur mengikat di `CLAUDE.md` §5: kode di-scaffold **di dalam
`OUTPUT_HASIL_RNM\`**, arah dependency **`handlers` → `services` → `repository`**.

| Lapisan | Tanggung jawab |
| --- | --- |
| `handlers` | Endpoint REST Komite; penegakan wewenang (**ADR-0014**) sebelum memanggil service |
| `services` | Mesin tangga; penyusunan roster; penulisan rekam akseptasi; penulisan outbox |
| `repository` | Akses `POOLDATA`; pemanggilan stored procedure; lookup `M_LINK_SERVICE` |
| `models` | Kasus komite + koleksi `KomiteList` per tingkat |
| worker | Pengirim outbox dengan retry (**ADR-0015**) — proses terpisah dari permintaan HTTP |
| `frontend` | Inbox komite per posisi; layar keputusan; status "perlu intervensi" |

`[pertanyaan terbuka]` Pembagian paket domain di dalam `internal/` belum ditetapkan — sama seperti
pada spec Claim — Life. Pemilik: **Lead Engineer**.

### 3. Mesin tangga persetujuan

**Unit keputusan = baris `AdjustmentList`** (**ADR-0011**), sama dengan Claim — Life.

| # | Langkah | Pelaku | Akibat |
| --- | --- | --- | --- |
| 1 | Kasus diterima dari Claim Life sebagai *child work* | — | `KomiteCount` mulai, `KomiteLoop` = jumlah roster |
| 2 | Kasus dirutekan ke baris roster pertama ber-`KomiteAproval == 0` | sistem | — |
| 3 | Anggota komite tingkat itu memutuskan **Setuju** / **Tolak** | pemilik `KomiteID` | `KomiteAproval`, `KomiteComment`, `DateApprove` terisi |
| 4a | **Setuju** dan `KomiteCount < KomiteLoop` | — | `KomiteCount + 1`, kembali ke langkah 2 |
| 4b | **Setuju** dan `KomiteCount == KomiteLoop` | — | rekam akseptasi + nomor dibuat; `STS_REJECT = 1`; efek keluar berjalan |
| 4c | **Tolak** | — | tangga berhenti; `STS_REJECT = 2` pada tingkat terakhir |

**Aturan mengikat:**

- `KomiteLoop` **dihitung**, bukan dikonfigurasi: COUNT baris roster `EMAILKOMITE` aktif
  ber-`LIMIT_BOTTOM <= |CLAIM_AMOUNT|`.
- Nilai keputusan adalah **enum tertutup** `{1, 2}`. Nilai lain ditolak terang-terangan.
- **Nomor dan rekam akseptasi lahir sekali**, di tingkat terakhir, pada keputusan Setuju.
- Tingkat bukan-terakhir **hanya** mencatat jejak persetujuan; tidak menyentuh tabel akseptasi.

### 4. Penegakan wewenang

**ADR-0014.** Hanya pemilik `KomiteList(KomiteCount).KomiteID` yang boleh menyimpan keputusan;
ditegakkan **di lapisan layanan**, bukan di visibilitas layar. Inbox hanya menampilkan kasus sesuai
posisi roster.

**Pengecualian sah:** eskalasi manual **naik satu tingkat** bila anggota berhalangan — aksi admin,
terekam. Pemutus di tingkat sama yang bukan pemilik `KomiteID`, dan eskalasi turun, tetap dilarang.

⚠️ **Penyimpangan sadar.** `[terverifikasi]` Pega tidak menegakkan apa pun — ia hanya menempatkan
kasus di worklist, dan `AcceptStatus` tidak ditulis rule mana pun (nol `<PropertiesName>` di 47
berkas modul).

### 5. Satu jalur simpan, bukan dua blok kembar

⚠️ **Penyimpangan sadar dari paritas struktural.** `[terverifikasi]` `UpdateOsAkseptasiClaimLife_sql`
dipanggil **dua kali** — step **4.15** (aksep, gerbang baris 5695) dan step **5.6** (tolak, gerbang
baris 8119) — dengan himpunan properti **identik** (diff struktural nol beda), beda hanya nilai
status.

**Sistem baru: satu jalur simpan berparameter status.** Dua blok kembar sepanjang ribuan baris adalah
tempat divergensi diam tumbuh.

### 6. Efek keluar — empat, wajib berhasil

**ADR-0015**, menyimpang dari **ADR-0008**.

| Urutan | Step | Efek |
| ---: | ---: | --- |
| 1 | **8** | `InsertJsonClaimLife_Act` |
| 2 | **10** | `serviceInsertArasapasClaimLife_act` — endpoint via `M_LINK_SERVICE` (**ADR-0013**) |
| 3 | **11** | `SendEmailKlaimLife` |
| 4 | **12** | `HitServiceToKasirKMTLife_Act` — **Kasir / pembayaran** |

Seluruhnya **setelah `Obj-Save`**. Pola: **transactional outbox** — keputusan + daftar efek dalam
satu transaksi; worker mengirim dengan retry sampai sukses.

**Syarat pengaman wajib:** ID idempoten unik per kiriman · **cek status sukses sebelum kirim ulang**
untuk Email dan Kasir · status **"perlu intervensi"** di UI Komite + laporan harian.

**Di luar lingkup:** step **14** `SetInformationData` — temporary, bukan efek keluar.

`[terverifikasi]` Step **10** dan **12** digerbangi `AcceptStatus = 1 && KomiteCount == KomiteLoop`
(baris 8648, 8887) — hanya berjalan pada keputusan Setuju di tingkat final.

### 7. Kode mati yang tidak dimigrasikan

| Yang dibuang | Bukti |
| --- | --- |
| `TransferType` + cabang routingnya + layar `ShowTransfer` | `[terverifikasi]` nol `Property-Set`; selalu FALSE |
| Gerbang EXIT step 9 + `1000013` / `L0000141` / `L0000134` | `[keputusan work owner]` (**OQ-064**) |
| Step 4.4, 4.5 `Generate No Akseptasi` | `[terverifikasi]` `blockname //` + tidak terindeks |
| Step 4.6 `Set Nilai Akseptasi` | `[terverifikasi]` `blockname //` |
| Step 4.14, 5.5 `Tukar SecurityReinsurer dengan RetroName` | `[terverifikasi]` `blockname //` |
| Logika penukaran retro | `[keputusan work owner]` (**OQ-065**) — nilai retro dipakai **apa adanya** dari data policy |

⚠️ **Konsekuensi membuang gerbang EXIT:** `[keputusan work owner]` klaim ber-retro tersebut yang
selama ini dikecualikan kini menjalankan **seluruh** efek keluar, **termasuk Kasir**. Ini
**perubahan perilaku yang menyentuh uang**, diterima sadar (**OQ-064**).

### 8. Data dan transaksi

| Objek | Peran | Catatan |
| --- | --- | --- |
| `POOLDATA.EMAILKOMITE` | roster + pita nilai | `[data DBA]` PK `ID`; `LIMIT_BOTTOM`/`LIMIT_TOP INTEGER`; `STS_AKTIF`; `STS_KLAIM`; **`STS_REJECT VARCHAR2(15)`** ⚠️ berbeda tipe dari tabel klaim |
| `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` | rekam akseptasi | `[data DBA]` `STS_REJECT NUMBER(38)`; uang `NUMBER` tanpa presisi → desimal presisi arbitrer (**ADR-0003**). `[terverifikasi]` 55 kolom terbaca langsung dari blok PL/SQL |
| `POOLDATA.GENERATE_SEQUENCE_NUMBER` | sequence | `[data DBA]` PK komposit `(CLASS, JENIS, TAHUN)` |
| `POOLDATA.M_LINK_SERVICE` | alamat endpoint | **ADR-0013** — lookup runtime, dilarang hardcode |
| `POOLDATA.GCP_IMAGE` | cache token | `[data DBA]` |
| `POOLDATA.DIRECTTOKASIR_LOG` | jejak Kasir di sisi database | `[terverifikasi]` |

**Batas transaksi dipegang Go** (OQ-013 tertutup untuk jalur Life). Outbox ditulis dalam transaksi
yang sama dengan keputusan.

⚠️ `[terverifikasi]` `UpdateOsAkseptasiClaimLife_sql` adalah **rule bersama** dengan Claim — Life
(hash ternormalisasi `c50bfd9a12` identik). Mengubah bentuknya adalah **perubahan kontrak lintas
konteks**.

---

### 9. Bentuk penyimpanan — **relasional** ⚠️ BARU 2026-09-16

⚠️ **Lubang yang ditutup di sini.** Sampai sebelum ini, spec memperlakukan `KomiteList` sebagai
**page runtime Pega** — riwayat tangga dianggap "sudah ada di data". Ia **tidak punya tabel** di
skema baru. Bagian ini memberinya satu.

⚠️ **Penyimpangan sadar 1 — roster & keputusan komite jadi tabel relasional.**
`[keputusan work owner]` Page/JSON Pega (`ClaimData…AdjustmentList.KomiteList`) **dibuang**.

#### Rantai penyimpanan

```
T_CLAIMLF_ADJUSTMENT (milik Claim Life) — KOMITE_ID = identitas kasus komite
     │                                   ⬅ penunjuk langsung: dari baris adjustment,
     │                                      pengguna tahu kasus komite mana yang memutuskannya
     │   kirim komite → LAHIR BARIS BARU di T_WORK_CLAIM
     ▼
T_WORK_CLAIM (lintas-lini, satu baris per work object)
     │   ID = identitas kasus komite ; COVER_KEY = ID baris klaim ; posisi/status tangga
     │   (TANPA KMT_NO · TANPA ADJUSTMENT_ID)
     ▼   penghubung: T_WORK_CLAIM.ID
T_GENERAL_KOMITE (header kasus komite) — WORK_CLAIM_ID → T_WORK_CLAIM.ID
     !!! RALAT 2026-09-18: kedua baris di atas DICABUT. Tidak ada kolom
     !!! WORK_CLAIM_ID. Penghubungnya ID YANG SAMA (shared primary key):
     !!! T_GENERAL_KOMITE.ID = T_WORK_CLAIM.ID baris komite.
     └─ T_KOMITE_KOMITELIST   1:N   FK DATA_KOMITE_ID → T_GENERAL_KOMITE.ID   ON DELETE CASCADE
```

⚠️ **Penunjuk DUA ARAH, dijaga konsisten dalam SATU transaksi saat kirim komite:**

| Arah | Kolom | Guna |
| --- | --- | --- |
| adjustment → komite | `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` = identitas kasus komite (= `T_WORK_CLAIM.ID`) | pintasan tampilan |
| komite → adjustment | `T_GENERAL_KOMITE.ADJUSTMENT_ID` → `T_CLAIMLF_ADJUSTMENT.ID` | menutup lingkar |

⚠️ **`ADJUSTMENT_ID` berada di `T_GENERAL_KOMITE`, BUKAN di `T_WORK_CLAIM`.** `T_WORK_CLAIM` adalah
tabel **lintas-lini**; ia memegang `ID`, `COVER_KEY`, beserta posisi/status tangga.

⚠️ **REVISI 2026-09-17** `[keputusan work owner]` — **`KMT_NO` DIBUANG.** `T_WORK_CLAIM` adalah tabel
**satu baris per work object**; kasus komite mendapat barisnya sendiri (`ID`), dan induknya ditunjuk
lewat `COVER_KEY`. Kolom nomor terpisah berarti menyimpan identitas yang sama dua kali. Test yang
menemukan kolom `KMT_NO` **gagal**.

#### Kolom

> ⚠️ **RALAT 2026-09-18 — SHARED PRIMARY KEY; kolom `WORK_CLAIM_ID` DIBUANG.**
> `[keputusan work owner]` Teks di bawah **tidak dihapus** sebagai jejak, tetapi yang **mengikat**
> adalah ralat ini.
>
> 1. ⛔ **Kolom `WORK_CLAIM_ID` tidak ada** — bukan diganti nama, **dibuang**. Hubungan
>    `T_WORK_CLAIM` ↔ `T_GENERAL_KOMITE` dijamin oleh **`ID` yang identik**: `T_GENERAL_KOMITE.ID`
>    **sama persis** dengan `T_WORK_CLAIM.ID` baris komite (**shared PK**, 1:1). Begitu pula di sisi
>    klaim: `T_GENERAL_CLAIM.ID` = `T_WORK_CLAIM.ID` baris klaim.
>    Setiap penyebutan *"`WORK_CLAIM_ID` penghubung"* atau *"nama usulan"* di bawah **dicabut**.
> 2. **`T_GENERAL_KOMITE.ADJUSTMENT_ID` TETAP ADA** — penutup lingkar ke baris adjustment. Ia
>    **bukan** bagian shared PK dan **tidak** ikut dibuang.
> 3. ✅ **Tipe `T_WORK_CLAIM.ID` DIPUTUSKAN** — **teks berformat**: baris klaim `CLM-xxxxxx`
>    (contoh `CLM-123456`), baris komite `KMT-xxxxxx` (contoh `KMT-000789`). **Bukan angka
>    sequence.** `COVER_KEY`, `T_GENERAL_CLAIM.ID`, `T_GENERAL_KOMITE.ID`, dan
>    `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` **mengikuti** tipe itu. Setiap `[terbuka]` *"teks atau angka
>    sequence — belum ditetapkan"* di bawah **dicabut**.
> 4. ⚠️ **Penyimpangan sadar dari ADR-0006**, dicatat bukan dilanggar diam-diam: identitas
>    `T_WORK_CLAIM` dan kedua tabel ber-shared-PK adalah **nomor bisnis berformat**, bukan sequence.
>    ADR-0006 **tetap berlaku** untuk `T_KOMITE_KOMITELIST.ID`.
> 5. ⚠️ `[terbuka]` **TETAP terbuka, jangan tebak:** (a) **generator** nomor `CLM-`/`KMT-` — siapa
>    yang membuatnya, sequence di belakang prefiks atau tidak, reset per tahun atau tidak;
>    (b) apakah `COVER_KEY` dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` dipasangi
>    `REFERENCES T_WORK_CLAIM(ID)`. Pemilik keduanya **DBA / work owner**.

**`T_GENERAL_KOMITE`** — header kasus komite, **lintas-lini**:
`ID` (PK, sequence), `WORK_CLAIM_ID` (penghubung ke `T_WORK_CLAIM.ID`), `ADJUSTMENT_ID`,
⚠️ *baris di atas **DICABUT** 2026-09-18 — `ID` **bukan** sequence melainkan **shared PK** =
`T_WORK_CLAIM.ID` baris komite (teks berformat `KMT-xxxxxx`); `WORK_CLAIM_ID` **tidak ada**.*
**`KOMITE_LOOP`** (jumlah tingkat = COUNT roster aktif), **`KOMITE_COUNT`** (tingkat berjalan),
**`ACCEPT_STATUS`** (hasil final: `1` aksep / `2` tolak), + audit (operator, tanggal).

**`T_KOMITE_KOMITELIST`** — satu baris **per anggota per jenjang**:
`ID` (PK), `DATA_KOMITE_ID` (FK), **`KOMITE_URUT`** (jenjang/urutan tangga), `KOMITE_ID` (anggota
pemutus), `ID_KOMITE` (jabatan), `KOMITE_EMAIL`, **`KOMITE_APROVAL`** (`0` belum / `1` setuju /
`2` tolak), `KOMITE_COMMENT`, **`DATE_APPROVE`** (**DATE**).

**`T_WORK_CLAIM`** — satu baris per work object: `ID` (identitas work object), **+`COVER_KEY`**
(`ID` baris induk; `NULL` bila tidak punya induk).
`[terbuka]` Tipe `T_WORK_CLAIM.ID` — teks atau angka sequence — **belum ditetapkan**; tipe `COVER_KEY`
dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` mengikutinya. **Jangan tebak.**

> ✅ **RALAT 2026-09-18 — baris `[terbuka]` di atas DICABUT, sudah dijawab.**
> `[keputusan work owner]` Tipe `T_WORK_CLAIM.ID` = **teks berformat**: `CLM-xxxxxx` (klaim) /
> `KMT-xxxxxx` (komite), **bukan** angka sequence. `COVER_KEY`, `T_GENERAL_CLAIM.ID`,
> `T_GENERAL_KOMITE.ID`, dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` mengikutinya.
> ⚠️ Penyimpangan sadar dari **ADR-0006** (identitas dari sequence) — dicatat, bukan dilanggar
> diam-diam; ADR-0006 tetap berlaku untuk `T_KOMITE_KOMITELIST.ID`.
> ⚠️ `[terbuka]` yang **tersisa**: **generator** nomor `CLM-`/`KMT-`, dan apakah `COVER_KEY` /
> `KOMITE_ID` dipasangi `REFERENCES T_WORK_CLAIM(ID)`. Pemilik **DBA / work owner**.
> **Jangan tebak.**
`[terbuka — Non-Life]` Kolom lintas-lini finalnya ditetapkan saat konteks Non-Life digarap.

#### Pemetaan kolom → field korpus `[terverifikasi]`

Field roster ber-class **`ASM-FW-GCNMFW-Data-Comitee`**.

| Kolom | Field korpus | Bukti |
| --- | --- | --- |
| `KOMITE_ID` | `KomiteID` | `Claim Life/Activity/CreateKMTLife_Act.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`): `KomiteList(<APPEND>).KomiteID = .OPERATOR_ID` |
| `KOMITE_EMAIL` | `KomiteEmail` | berkas sama: `.KomiteEmail = .EMAIL` |
| `KOMITE_APROVAL` | `KomiteAproval` | berkas sama: `.KomiteAproval = 0` saat roster dibentuk |
| `ID_KOMITE` | `IDKomite` | class `ASM-FW-GCNMFW-Data-Comitee` |
| `KOMITE_APROVAL` (saat putus) | `KomiteAproval` | `Komite Claim Life/Activity/KomitePostAdjustment.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`): `KomiteList(Local.Komite).KomiteAproval = pyWorkPage.AcceptStatus` |
| `KOMITE_COMMENT` | `KomiteComment` | berkas sama: `.KomiteComment = pyWorkPage.Comment` |
| `DATE_APPROVE` | `DateApprove` | berkas sama: `.DateApprove = @CurrentDateTime()` |

`[terverifikasi]` `AcceptStatus` ber-class `ASM-FW-GCNMFW-Work-KomiteLife`.

`[data DBA]` **`EMAILKOMITE` tetap tabel master yang DIBACA** — sumber `KomiteID`, email, jabatan,
dan pita nilai. Ia **bukan** tabel baru dan **tidak** diganti.

#### Penyimpangan sadar

⚠️ **2 — rujukan memakai `ID` / `COVER_KEY` yang stabil, bukan indeks posisi Pega.**
`[terverifikasi]` `CreateKMTLife_Act` memakai `IndexAdjustment` / `IndexPremiumList`
(`.pxListSubscript`), yang rusak begitu urutan baris bergeser.

⚠️ **3 — tidak ada hapus fisik.** Integritas dijaga PK/FK, dengan **cascade dari header ke anak**.

⚠️ **4 — `DATE_APPROVE` bertipe `DATE`**; kolom uang tetap **desimal presisi arbitrer**
(**ADR-0003**).

---

## Testing Decisions

### Apa yang membuat test baik di sini

Test menguji **perilaku yang teramati dari luar**: apa yang terjadi pada tangga, pada baris
`AdjustmentList`, dan pada antrean efek keluar — bukan bahwa suatu fungsi dipanggil. Test yang hanya
membenarkan implementasi tidak diterima.

### Seam — **memakai ulang seam Claim — Life**, tidak menambah

`[terverifikasi]` Repo target belum di-scaffold. Seam yang sudah ditetapkan spec Claim — Life adalah
**API HTTP**, dan konteks ini memakainya kembali — bukan seam baru:

> **Seam utama: API HTTP Komite Claim Life.** Test menggerakkan tangga lewat endpoint REST dan
> memeriksa hasilnya lewat endpoint REST, dengan `handlers → services → repository` terpasang
> sungguhan, terhadap skema uji Oracle.

Alasannya sama: titik tertinggi yang masih memberi perilaku ujung-ke-ujung, dan tidak mengunci
pembagian paket yang belum diputuskan.

**Satu seam kedua yang tidak terhindarkan: worker outbox.** `[keputusan work owner]` Efek keluar
dikirim oleh **proses terpisah**, bukan di dalam permintaan HTTP. Perilakunya — retry sampai sukses,
cek status sebelum kirim ulang, penandaan "perlu intervensi" — **tidak teramati** lewat API HTTP
saja. Usul: seam kedua di **batas pengirim outbox**, dengan efek luar difake.

Ini seam kedua yang **sungguh diperlukan**, bukan kenyamanan — dan perlu dikonfirmasi.

**Batas proses difake:**

| Batas | Perlakuan |
| --- | --- |
| Arasapas, Email, Kasir, konversi JSON | *fake* di balik interface; test memeriksa **efeknya**, termasuk perilaku retry dan anti-dobel |
| Oracle | skema uji nyata — penomoran dan roster adalah inti perilaku |
| Jam | dapat dikendalikan — `DateApprove` dan retry bergantung waktu |

### Modul yang diuji

| Yang diuji | Lewat seam |
| --- | --- |
| Tangga: routing ke baris `KomiteAproval == 0`, naik tingkat, berhenti saat Tolak | API HTTP |
| `KomiteLoop` = COUNT roster, termasuk nilai mutlak klaim negatif | API HTTP |
| Penegakan wewenang per `KomiteID` + eskalasi naik satu tingkat | API HTTP |
| Nomor & rekam akseptasi sekali di tingkat final | API HTTP |
| Jalur balik `STS_REJECT` ke dua tingkat baris (kontrak tiket 11) | API HTTP |
| Enum keputusan tertutup `{1,2}`; nilai lain ditolak | API HTTP |
| Outbox: retry, anti-dobel Email/Kasir, status "perlu intervensi" | seam worker + fake |
| Uang tanpa `float` | API HTTP |

### Prior art

`[terverifikasi]` **Tidak ada** — nol kode, nol test di `OUTPUT_HASIL_RNM`. Spec Claim — Life
menetapkan bentuknya; konteks ini **mengikuti bentuk yang sama**, dan slice pertama Claim Life akan
menjadi prior art nyata bila dikerjakan lebih dulu.

Perintah verifikasi wajib ditulis eksplisit di tiap tiket selama `Makefile` belum ada. Target:
`go test ./internal/...` dan `cd frontend && npm test`.

---

## Acceptance Criteria

Dapat diverifikasi tanpa pengetahuan pribadi.

**Tangga**

1. Kasus baru dirutekan ke baris roster **pertama** yang ber-`KomiteAproval == 0`.
2. `KomiteLoop` sama dengan **COUNT** baris roster `EMAILKOMITE` aktif (`STS_AKTIF = "1"`) ber-
   `LIMIT_BOTTOM <= CLAIM_AMOUNT`; tidak dibaca dari konstanta mana pun.
3. Klaim bernilai **negatif** dicari rosternya dengan **nilai mutlak**.
4. ⚠️ `[keputusan work owner]` Bila tidak ada baris roster yang cocok, penyerahan **gagal
   terang-terangan** — bukan tangga nol tingkat. **Ini penjaga defensif, bukan alur normal:**
   `[keputusan work owner]` roster dijamin **≥ 1** secara bisnis (limit berjenjang selalu menutup
   nilai klaim berapa pun). AC ini menjaga bila **data master limit rusak**. `[terverifikasi]` Pega
   **tidak** memuat gerbang ini — `Claim Life/Activity/CreateKMTLife_Act.xml`
   (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`) tetap men-set
   `childPageKomite.KomiteCount = 1` **tanpa syarat** meski
   `KomiteLoop = @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)` bernilai `0` — tangga
   nol tingkat, senyap. Nol `Page-Set-Messages` / `Exit-Activity` / `pxAddMessage` / `Obj-Validate`
   yang menjaganya. "Gagal terang-terangan" adalah **keputusan sistem baru**.
   **Dua sisi.** Sisi **kirim** (Claim Life, tiket 10) menggagalkan penyerahan. Sisi **terima**
   (Komite, tiket 01) **menolak** muatan ber-`KomiteLoop < 1` atau ber-`KomiteList` kosong —
   `T_GENERAL_KOMITE` dan `T_KOMITE_KOMITELIST` tidak dibuat. Penjaga **berlapis**: sisi terima tetap
   menolak walau sisi kirim bocor atau muatan masuk lewat API langsung. Test yang menemukan
   `T_GENERAL_KOMITE` lahir tanpa satu pun baris `T_KOMITE_KOMITELIST` **gagal**.
5. Keputusan **Setuju** pada tingkat bukan-terakhir menaikkan `KomiteCount` satu dan **tidak**
   menyentuh tabel akseptasi.
6. Keputusan **Tolak** menghentikan tangga pada tingkat mana pun ia terjadi.
7. Tiap tingkat menghasilkan satu entri berisi keputusan, komentar, dan waktu.

**Wewenang**

8. Pengguna yang **bukan** pemilik `KomiteID` tingkat berjalan **ditolak** saat menyimpan keputusan,
   meskipun ia dapat membuka kasusnya.
9. Penolakan terjadi **di lapisan layanan**, dan tetap terjadi meskipun kontrol UI ditampilkan.
10. Inbox seorang anggota komite hanya memuat kasus pada tingkat yang ia miliki.
11. Eskalasi memindahkan kasus **naik satu tingkat**; eskalasi turun ditolak.
12. Eskalasi tercatat: siapa memindahkan, kapan, dari tingkat mana ke tingkat mana.

**Keputusan final**

13. Rekam akseptasi dan nomornya dibuat **sekali**, pada keputusan Setuju di tingkat terakhir.
14. Nomor akseptasi diperoleh lewat rantai `GetKodeProdLife_SQL` → `GetSequenceNumber_SQL`; aplikasi
    **tidak** memuat logika pembentukan format nomor (**ADR-0006**).
15. `STS_REJECT` tertulis pada **dua tingkat baris** — `PremiumListDetail` dan `AdjustmentList` —
    dengan nilai yang sama (kontrak **tiket 11** Claim Life).
16. Penyimpanan rekam akseptasi memakai **satu jalur** berparameter status; tidak ada dua jalur
    kembar di kode.
17. Tidak ada nilai uang sebagai *binary floating point* di lapisan mana pun maupun di JSON.

**Efek keluar**

18. Keempat efek dikirim **setelah** keputusan tersimpan; kegagalan jaringan tidak membatalkan
    keputusan.
19. Keputusan dan daftar efeknya tersimpan dalam **satu transaksi**.
20. Efek yang gagal **diantre ulang** sampai berhasil atau ditandai **perlu intervensi**.
21. Tiap kiriman membawa **ID idempoten unik**.
22. Sebelum mengirim ulang **Email** atau **Kasir**, sistem memeriksa status "sudah terkirim sukses";
    bila sudah, tidak dikirim lagi.
23. Kasus ber-status "perlu intervensi" terlihat di UI Komite dan masuk laporan harian.
24. Alamat endpoint di-resolve **runtime** dari `M_LINK_SERVICE` lewat `(KATEGORI_1, KATEGORI_2)`;
    tidak ada URL sebagai literal, konstanta, maupun env var (**ADR-0013**).
25. Keputusan komite dilaporkan **tuntas** hanya setelah keempat efek berhasil.

**Kode mati**

26. Tidak ada padanan `TransferType`, gerbang EXIT retro, maupun ketiga identitas ter-hardcode di
    kode.
27. Nilai retro yang ditulis ke rekam akseptasi adalah nilai **apa adanya** dari data policy; tidak
    ada logika penukaran.
28. Tidak ada padanan `Generate_NoAccept_KMT_Life` / `_LifeRetro`.
29. Status ditampilkan sebagai kata, bukan nama field `STS_REJECT` dan bukan angka.

**Penyimpanan** ⚠️ BARU 2026-09-16 — §9

30. ⚠️ Kasus komite disimpan sebagai **`T_GENERAL_KOMITE`** (header) + **`T_KOMITE_KOMITELIST`** (satu baris
    **per anggota per jenjang**). **Tidak ada page maupun JSON** sebagai penyimpan roster atau
    keputusan. Test yang menemukannya **gagal**. *(§9; penyimpangan sadar 1)*
31. ⚠️ Setiap keputusan anggota menulis **`KOMITE_APROVAL`**, **`KOMITE_COMMENT`**, dan
    **`DATE_APPROVE`** ke baris `T_KOMITE_KOMITELIST` yang bersesuaian; `DATE_APPROVE` bertipe **`DATE`**.
    *(§9; penyimpangan sadar 4)*
32. ⚠️ **`T_GENERAL_KOMITE.ADJUSTMENT_ID` dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` diisi dalam transaksi yang
    sama** saat kirim komite — penunjuk dua arah **tidak pernah** setengah terisi. Test yang
    menemukan salah satunya kosong sementara yang lain terisi **gagal**. *(§9)*
33. ⚠️ Riwayat tangga dibaca dari **`T_KOMITE_KOMITELIST` diurut `KOMITE_URUT`** — **bukan** dari page
    runtime. *(§9; penyimpangan sadar 1)*
34. ⚠️ Rujukan antar tabel memakai **`ID` / `COVER_KEY`**, bukan **indeks posisi**. Test yang menemukan
    padanan `IndexAdjustment` / `IndexPremiumList` sebagai kunci rujukan **gagal**. *(§9;
    penyimpangan sadar 2)*

**Nilai keputusan**

35. Nilai keputusan di luar enum tertutup **`{1 = Setuju, 2 = Tolak}`** **ditolak terang-terangan**
    di lapisan layanan — tangga **tidak** berhenti diam-diam, dan keputusan tidak tersimpan.
    `[keputusan work owner]` untuk penolakannya; `[terverifikasi]` bahwa Pega **tidak** menolak —
    `Komite Claim Life/When/IsKomiteLoop.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `ISKOMITELOOP` /
    `RULE-OBJ-WHEN`) hanya melanjutkan pada `"1"`; nilai lain menghentikan tangga **tanpa pesan**.
    Menaikkan **User story 11** menjadi kriteria yang dapat di-test.

---

## Pertanyaan terbuka di dalam spec

Dua OQ **non-pemblokir**. Bagian yang bergantung padanya ditandai; **jangan ditebak**.

| OQ | Pertanyaan | Pemilik | Bagian yang menunggu |
| --- | --- | --- | --- |
| **OQ-035** | `UpdateWorkObject` dan `serviceInsertArasapasClaimLife_act` dipanggil dari Komite tetapi salinannya hanya ada di modul lain — resolusi lewat pewarisan class atau ekspor tidak lengkap? Perilakunya sudah terbaca; **kepemilikan rule-nya** belum | Arsitektur Pega | §6 efek keluar — **bukan pemblokir isi** |
| **OQ-007 / OQ-021** | Model RBAC lintas konteks belum ditetapkan; pemetaan **`KomiteID` → identitas akun** bergantung padanya | IAM | §4 penegakan wewenang (**ADR-0014**) |

**Asumsi eksplisit yang menyertainya:** spec ini mengasumsikan setiap `KomiteID` dapat dipetakan ke
satu identitas akun yang dapat diautentikasi. Bila pemetaan itu ternyata banyak-ke-banyak, aturan
penegakan di §4 perlu ditinjau ulang.

`[terverifikasi]` Identity & Access ditandai **ABSENT** dari korpus (`discovery/context-map.md`) —
nol rule otorisasi. Ia dibangun dari nol.

---

## Out of Scope

- **Isi dan alur internal Claim — Life.** Konteks luar (**ADR-0001**); spec ini berhenti di tiga
  kontrak batas.
- **Tiga modul Komite lain** (`Komite Claim FacIn`, `Komite Claim Prop`, `Komite Claim Non Prop`).
  `[terverifikasi]` Grafnya sama bentuk, tetapi **empat implementasi berbeda untuk satu konsep** —
  sasaran routing, guard identitas, dan rule penomoran semuanya berlainan. OQ-036 (sasaran
  ter-hardcode `komitepnc`…`komitepnc4`) tidak menyentuh Life.
- **Merapikan alur.** Urutan tangga dan bentuk keputusan dibawa apa adanya. Tiga penyimpangan sadar
  sudah didaftar di §Solution; selain itu, paritas.
- **Identity & Access sebagai konteks.** Dibangun dari nol, di luar spec ini.
- **Scaffolding kode.** `cmd/`, `internal/`, `frontend/`, `Makefile` belum ada — pekerjaan terpisah
  yang mendahului tiket mana pun.
- **Kontrak layanan Kasir.** `[terbuka]` OQ-002 — bentuk permintaan dan makna jawabannya tidak ada di
  korpus. Spec menetapkan **jaminan pengirimannya**, bukan bentuk pesannya.

---

## Further Notes

**Ukuran pekerjaan.** `[terverifikasi]` Modul `Komite Claim Life` memuat 47 berkas: 18 Activity,
14 RDBList, 3 ReportDefinition, 3 FlowAction, 3 Section, 2 When, 1 Flow, 1 DecisionTable,
1 ConnectREST, 1 SystemSettings. Berkas terbesar `ShowTransfer.xml` (1.125.234 byte) dan
`KomitePostAdjustment.xml` (515.675 byte).

**Urutan yang saya sarankan untuk `/to-tickets`** — vertical slice:

1. Penerimaan kasus dari Claim Life + inbox per posisi (kontrak masuk, tiket 10 Claim Life).
2. Mesin tangga: routing `KomiteAproval == 0`, naik tingkat, berhenti saat Tolak.
3. Penegakan wewenang per `KomiteID` + eskalasi naik satu tingkat.
4. Keputusan final: nomor + rekam akseptasi sekali di tingkat terakhir.
5. Jalur balik `STS_REJECT` ke Claim Life (kontrak keluar, tiket 11 Claim Life).
6. Transactional outbox + worker retry.
7. Anti-dobel Email & Kasir + status "perlu intervensi" + laporan harian.
8. Antarmuka React: inbox, layar keputusan, riwayat tangga.

**Yang perlu keputusan sebelum tiket dibuat:**

1. **Seam kedua untuk worker outbox** (§Testing Decisions) — perilaku retry tidak teramati lewat API
   HTTP saja. Perlu konfirmasi.
2. **Pembagian paket domain di dalam `internal/`** — sama seperti spec Claim — Life. Pemilik: Lead
   Engineer.
3. **Pemetaan `KomiteID` → identitas akun** (OQ-007/OQ-021) — asumsi satu-ke-satu dinyatakan di
   §Pertanyaan terbuka.

**Catatan sumber.** Spec ini bersandar **hanya** pada korpus Pega `D:\XML\RNM_BRD\` dan artefak di
`OUTPUT_HASIL_RNM\`. Sumber ADR tunggal: `docs/adr/ADR-0001`…`ADR-0015`.
