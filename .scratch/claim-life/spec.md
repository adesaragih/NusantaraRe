# Spec — Claim — Life (migrasi Pega → Go + React + Oracle)

Status: ready-for-agent
Konteks: `claim-life` (Claim — Life)
Tanggal: 2026-09-14 · penyimpanan ditetapkan 2026-09-16
Sumber: grilling Ronde 1–4 (`grilling-ronde-1.md` … `-4.md`),
`revisi-penyimpanan-json-dibuang.md` (bentuk penyimpanan), `CONTEXT.md`,
`docs/adr/ADR-0001`–`ADR-0015`, `discovery/flows/Claim Life.md`,
`discovery/modules/Claim Life.md`, `discovery/flows/_SUMMARY-claim.md`,
`discovery/open-questions.md`, `kesiapan-to-spec.md`
Skill: `/mattpocock-skills:to-spec`

> **Peta cepat.** Perilaku klaim ada di **§3–§11** dan **§13–§16**; **bentuk penyimpanannya** ada di
> **§2b**, dan **migrasinya** di **§14**. Seluruh acceptance criteria bernomor **1–57**; tiket
> merujuk nomor itu apa adanya, jadi penomoran **tidak pernah digeser**.

> **Konvensi penandaan.** Setiap pernyataan perilaku membawa salah satu dari:
> `[terverifikasi]` — terbaca di korpus, disertai path + rule;
> `[keputusan work owner]` — dinyatakan work owner, **tidak** terbukti di korpus;
> `[terbuka]` — belum diketahui, nomor OQ dari `discovery/open-questions.md`.
> Tidak ada pernyataan tanpa salah satu penanda. ADR dirujuk, **tidak** diulang isinya.

> **Sumber tunggal.** Spec ini bersandar **hanya** pada korpus Pega `D:\XML\RNM_BRD\` (READ-ONLY)
> dan artefak di `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\`, yang sekaligus **repo target tunggal** untuk
> dokumen maupun kode. Sumber ADR tunggal: `docs/adr/ADR-0001`…`ADR-0012`.
> ⛔ `D:\XML\nusantara-re\` **di-blacklist** (keputusan work owner 2026-09-14) — tidak dibaca,
> tidak dijadikan pembanding, tidak dijadikan target. Tidak ada satu pun klaim di spec ini yang
> bersandar padanya.

---

## Problem Statement

Penanganan klaim jiwa (*life*) Nusantara Re hari ini berjalan di Pega. Tiga peran — admin klaim,
penasihat medis, dan supervisor — mendaftarkan klaim, menelaah sisi medis, lalu memutuskan
akseptasi, dengan Komite sebagai tangga persetujuan di luar sistem ini.

Masalah yang dihadapi pengguna bila status quo dipertahankan:

1. **Aturannya tidak terbaca.** Logika penting tersebar di 51 Activity, 29 RDBList, 20 Section, dan
   16 FlowAction; sebagian berada di stored procedure Oracle yang badannya tidak ada di mana pun
   yang dapat dibaca tim. Tidak ada satu tempat yang menyatakan "kapan klaim selesai".
2. **Jejak siapa-melakukan-apa tidak lengkap.** Pengembalian kasus ke admin atau ke medis tidak
   merekam pelakunya sama sekali — hanya penanda bernilai `1`.
3. **Nama menyesatkan.** Kolom status klaim bernama `STS_REJECT`, padahal nilai `1` berarti
   **diaksep**. Siapa pun yang membaca sistem berisiko salah paham.
4. **Wewenang bergantung pada atribut data, bukan peran.** Untuk dua tipe polis tertentu, siapa pun
   dapat mengirim kasus ke Komite.

## Solution

Membangun konteks **Claim — Life** sebagai layanan Go + antarmuka React yang menulis ke Oracle
`POOLDATA` yang sama, dengan:

- **Mesin status yang dinyatakan eksplisit** dan berunit **baris `AdjustmentList`** (**ADR-0011**).
- **Jejak audit siapa + kapan untuk setiap transisi dan setiap jalur balik** (**ADR-0007**) —
  perbaikan sadar atas sistem lama.
- **Peran ditegakkan di lapisan layanan**, bukan sekadar di visibilitas layar (**ADR-0002**).
- **Uang tanpa `float`** (**ADR-0003**).
- **Integrasi luar bersifat asinkron, tidak memblokir, dengan antre-ulang** (**ADR-0008**).
- **Komite Life tetap sistem luar**, dihubungkan tiga kontrak eksplisit (**ADR-0001**).

Perilaku bisnis dibawa apa adanya (paritas), kecuali empat penyimpangan sadar yang sudah
diputuskan: jejak audit (**ADR-0007**), antre-ulang (**ADR-0008**), muatan kontrak Komite yang
diperluas (**ADR-0001**), dan representasi uang (**ADR-0003**).

---

## User Stories

### Pendaftaran dan Outstanding — `ReasLifeAdmin`

1. Sebagai **ReasLifeAdmin**, saya ingin mendaftarkan klaim life baru, supaya klaim masuk ke siklus
   penanganan. `[terverifikasi]` `Claim Life/Flow/Register_Flow.xml`
   (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW`), shape `Assignment2`
   "Input Register".
2. Sebagai **ReasLifeAdmin**, saya ingin memperoleh nomor klaim otomatis, supaya penomoran konsisten
   dan tidak bentrok. `[terverifikasi]` `GetSequenceNumber_SQL` → `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`
   (**ADR-0006**).
3. Sebagai **ReasLifeAdmin**, saya ingin menginput baris **AdjustmentList** pertama dan menyimpannya
   ke Outstanding, supaya klaim tercatat sebagai sedang berjalan. `[keputusan work owner]` — langkah 1
   mesin status (**ADR-0011**); `[terverifikasi]` penulis `STS_REJECT = 0` adalah
   `Claim Life/Activity/SaveOutStandingLife_Act.xml`.
4. Sebagai **ReasLifeAdmin**, saya ingin sistem menolak *Date of Loss* yang berada di luar jendela
   valuasi polis, supaya klaim yang tidak tertanggung tidak masuk ke siklus. `[terverifikasi]`
   `Claim Life/Activity/ValidasiDOL_Act.xml`
   (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `VALIDASIDOL_ACT` / `RULE-OBJ-ACTIVITY`).
5. Sebagai **ReasLifeAdmin**, saya ingin jenis klaim terisi otomatis dari kode produk, supaya saya
   tidak salah memilih. `[terverifikasi]` `ContentNote` diturunkan dari `BusinessCode`
   (`Claim Life/Activity/SaveOutStandingLife_Act.xml`); tabel `L1`–`L21` di `CONTEXT.md`.
6. Sebagai **ReasLifeAdmin**, saya ingin menolak baris adjustment yang saya input sendiri tanpa
   melalui Komite, supaya kesalahan input dapat dibatalkan cepat. `[keputusan work owner]` —
   langkah 1b; `[terverifikasi]` tombol "Reject Outstanding" di
   `Claim Life/Section/AdjustmentDetail_Section.xml` bergerbang
   `pyPosition =='ReasLifeAdmin' && …CLAIM_NO !='' && .STS_REJECT == 0`.
7. Sebagai **ReasLifeAdmin**, saya ingin penolakan saya hanya membatalkan **baris itu**, supaya
   klaimnya tetap hidup dan saya dapat menginput baris pengganti. `[keputusan work owner]`
   (**ADR-0011**).
8. Sebagai **ReasLifeAdmin**, saya ingin melanjutkan klaim ke tahap Medical Check, supaya sisi medis
   ditelaah oleh yang berwenang. `[terverifikasi]` shape `Assignment3` "Medical Check" pada
   `Register_Flow`.
9. Sebagai **ReasLifeAdmin**, saya ingin menerima kembali kasus yang dikembalikan dari medis atau
   supervisor, supaya kekurangan data dapat saya lengkapi. `[terverifikasi]`
   `Claim Life/When/IsSendtoAdmin.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` /
   `RULE-OBJ-WHEN`); menggerbangi pengembalian dari **tiga titik** di `Register_Flow`.
10. Sebagai **ReasLifeAdmin**, saya ingin melihat siapa yang mengembalikan kasus dan kapan, supaya
    saya tahu apa yang diminta. `[keputusan work owner]` — penyimpangan sadar (**ADR-0007**); hari
    ini tidak terekam sama sekali.

### Telaah medis — `ReasLifeMedicalAdvisor`

11. Sebagai **ReasLifeMedicalAdvisor**, saya ingin menerima klaim yang menunggu telaah medis, supaya
    antrean kerja saya jelas. `[terverifikasi]` FlowAction `MEDICALCHECK`; shape `Assignment3`.
12. Sebagai **ReasLifeMedicalAdvisor**, saya ingin mencatat hasil telaah medis, supaya keputusan
    akseptasi punya dasar. `[terverifikasi]` `Claim Life/Section/MedicalCheckClaimLife.xml`
    (memuat gerbang `pyPosition`).
13. Sebagai **ReasLifeMedicalAdvisor**, saya ingin mengembalikan kasus ke admin bila data kurang,
    supaya saya tidak memutuskan di atas data tidak lengkap. `[terverifikasi]` `IsSendtoAdmin`.
14. Sebagai **ReasLifeMedicalAdvisor**, saya ingin meneruskan kasus ke supervisor, supaya keputusan
    akseptasi diambil oleh yang berwenang. `[keputusan work owner]` — langkah 3 mesin status.
15. Sebagai **ReasLifeMedicalAdvisor**, saya **tidak** ingin dapat mengubah status akseptasi, supaya
    batas wewenang jelas. `[keputusan work owner]` — tahap saya tidak mengubah `STS_REJECT`.

### Keputusan akseptasi — `ReasLifeSPV`

16. Sebagai **ReasLifeSPV**, saya ingin menerima klaim yang siap dianalisis, supaya saya dapat
    memutuskan. `[terverifikasi]` FlowAction `AKSEPTASICLAIMLIFE`; shape `Assignment4` "Claim Analis".
17. Sebagai **ReasLifeSPV**, saya ingin mengirim adjustment ke Komite, supaya keputusan akseptasi
    diambil sesuai tangga persetujuan. `[terverifikasi]` `Claim Life/Activity/CreateKMTLife_Act.xml`
    (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`) — `Call
    pxAddChildWork` berkelas `ASM-FW-GCNMFW-Work-KomiteLife`.
18. Sebagai **ReasLifeSPV**, saya ingin pengiriman ke Komite berlaku untuk **semua** adjustment,
    bukan hanya di atas ambang tertentu, supaya tidak ada jalur keputusan yang melewati Komite.
    `[keputusan work owner]` — langkah 4.
19. Sebagai **ReasLifeSPV**, saya ingin Komite menerima nilai klaim dan mata uangnya, supaya mereka
    tidak perlu membaca balik ke sistem ini. `[keputusan work owner]` — muatan diperluas
    (**ADR-0001**).
20. Sebagai **ReasLifeSPV**, saya ingin hasil keputusan Komite terpantul ke baris adjustment yang
    bersangkutan, supaya status di sistem ini selalu mencerminkan keputusan terakhir.
    `[terverifikasi]` `Komite Claim Life/Activity/KomitePostAdjustment.xml`
    (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`).
21. Sebagai **ReasLifeSPV**, saya ingin menambah baris adjustment baru setelah Komite menolak,
    supaya klaim dapat diajukan ulang dengan angka yang diperbaiki. `[keputusan work owner]` —
    langkah 6.
22. Sebagai **ReasLifeSPV**, saya ingin baris baru mewarisi angka dari baris pertama, supaya saya
    tidak mengetik ulang data yang sama. `[terverifikasi]`
    `Claim Life/Activity/SetIndexAdjustmentList.xml` menyalin 8 kolom dari `AdjustmentList(1)` ke
    `AdjustmentList(<LAST>)` — **tanpa** `STS_REJECT`.
23. Sebagai **ReasLifeSPV**, saya ingin mengembalikan kasus ke medis bila telaahnya perlu diulang,
    supaya keputusan tidak diambil di atas telaah yang meragukan. `[terverifikasi]`
    `Claim Life/When/IsSendtoMedical.xml`.
24. Sebagai **ReasLifeSPV**, saya ingin mengembalikan kasus ke admin, supaya kekurangan administratif
    dapat diperbaiki. `[terverifikasi]` `IsSendtoAdmin`.

### Status dan kejelasan — semua peran

25. Sebagai **pengguna mana pun**, saya ingin melihat status tiap baris adjustment secara terpisah,
    supaya saya tahu persis apa yang sudah diputuskan dan apa yang belum. `[keputusan work owner]`
    (**ADR-0011**).
26. Sebagai **pengguna mana pun**, saya ingin status ditampilkan dengan kata yang benar
    (Outstanding / Aksep / Ditolak), **bukan** nama field `STS_REJECT`, supaya saya tidak salah
    paham. `[terverifikasi]` nilai `1` = diaksep meski namanya "reject" — `CONTEXT.md`.
27. Sebagai **pengguna mana pun**, saya ingin melihat riwayat lengkap putaran Komite pada satu
    klaim, supaya saya paham mengapa klaim berada di keadaannya sekarang. `[keputusan work owner]` —
    jumlah baris = jumlah putaran.
28. Sebagai **pengguna mana pun**, saya ingin klaim yang seluruh barisnya sudah diputus tampak
    "selesai", supaya antrean kerja saya bersih. `[keputusan work owner]` — "selesai" adalah
    **keadaan turunan**, bukan status tersimpan; aturannya ditetapkan di §Implementation Decisions.

### Audit dan kepatuhan

29. Sebagai **auditor**, saya ingin mengetahui siapa dan kapan untuk setiap transisi status, supaya
    keputusan dapat dipertanggungjawabkan. `[keputusan work owner]` (**ADR-0007**).
30. Sebagai **auditor**, saya ingin setiap jalur balik (`SendtoAdmin`, `SendtoMedical`) terekam
    pelakunya, supaya pengembalian kasus dapat ditelusuri. `[keputusan work owner]` (**ADR-0007**);
    `[terverifikasi]` hari ini kedua penanda hanya menyimpan nilai `1`.
31. Sebagai **auditor**, saya ingin kegagalan integrasi luar terlihat di jalur audit, bukan hanya di
    log layanan, supaya kegagalan diam tidak terjadi. `[keputusan work owner]` (**ADR-0008**).

### Integrasi luar

32. Sebagai **ReasLifeAdmin**, saya ingin mengunggah dokumen pendukung klaim, supaya berkasnya
    lengkap. `[terverifikasi]` `Claim Life/Activity/InsertGoogleStorage_Act.xml`
    (`ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `INSERTGOOGLESTORAGE_ACT`) (**ADR-0010**).
33. Sebagai **pengguna mana pun**, saya ingin mengunduh dokumen yang sudah diunggah, supaya saya
    dapat memeriksanya. `[terverifikasi]` `Claim Life/Activity/GetUrlGoogleStorage_Act.xml`,
    `DownloadDocumentClaim.xml`.
34. Sebagai **ReasLifeSPV**, saya ingin anggota Komite menerima notifikasi saat kasus diserahkan,
    supaya keputusan tidak tertunda. `[terverifikasi]` `Claim Life/Activity/SendEmailKlaimLF.xml`.
35. Sebagai **pengguna mana pun**, saya ingin alur klaim **tidak tertahan** saat layanan luar gagal,
    supaya pekerjaan saya tidak terhenti karena persoalan teknis. `[keputusan work owner]`
    (**ADR-0008**); `[terverifikasi]` paritas — `Claim Life` hanya punya pencatatan
    (`RDBList/InsertLogServiceClaim.xml` → `pooldata.monitoring_klaim_log`), bukan gerbang
    keberhasilan seperti `IsSuccessHitService` di konteks facultative.
36. Sebagai **operator sistem**, saya ingin efek keluar yang gagal dapat diantre ulang, supaya tidak
    hilang diam-diam. `[keputusan work owner]` (**ADR-0008**) — kemampuan **baru**.

### Uang

37. Sebagai **Finance**, saya ingin nilai uang tidak kehilangan presisi, supaya angka reasuransi
    dapat dipertanggungjawabkan. `[keputusan work owner]` (**ADR-0003**) — 8 kolom uang.
38. Sebagai **Finance**, saya ingin `EM_PERCENT` diperlakukan sebagai persen, bukan uang, supaya
    maknanya tidak tercampur. `[keputusan work owner]` (**ADR-0003**).
39. Sebagai **Finance**, saya ingin mata uang menyertai nilai uang ke Komite, supaya pita nilai
    dibandingkan atas dasar yang benar. `[keputusan work owner]` (**ADR-0001**).

### Migrasi

40. Sebagai **work owner**, saya ingin seluruh data klaim life dipindahkan, supaya tidak ada
    pekerjaan tertinggal di Pega. `[keputusan work owner]` (**ADR-0009**).
41. Sebagai **work owner**, saya ingin klaim yang sedang berjalan terbawa beserta seluruh barisnya,
    supaya riwayat putaran Komite tidak hilang. `[keputusan work owner]` (**ADR-0011**, **ADR-0009**).
42. Sebagai **work owner**, saya ingin penomoran klaim tidak melompat atau mengulang setelah
    migrasi, supaya nomor tetap unik. `[terbuka]` **OQ-002** — kontrak procedure belum diketahui.

### Penyimpanan relasional — §2b

43. Sebagai **organisasi**, saya ingin setiap atribut klaim menjadi **kolom bernama** dengan tipe
    yang benar, supaya bentuk klaim dapat diperiksa, dicari, dan divalidasi — bukan tersembunyi di
    dalam satu dokumen. `[keputusan work owner]`
44. Sebagai **`ReasLifeAdmin`**, saya ingin kesalahan isian ketahuan **saat menyimpan**, bukan
    berbulan-bulan kemudian ketika laporan tidak cocok.
45. Sebagai **organisasi**, saya ingin klaim punya **tabelnya sendiri**, supaya ia berdiri sebagai
    entitas, bukan tempelan pada rekam akseptasi. `[keputusan work owner]`
46. Sebagai **organisasi**, saya ingin satu klaim tersimpan **utuh atau tidak sama sekali**, supaya
    tidak pernah ada klaim yang tersimpan separuh.
47. Sebagai **`ReasLifeSPV`**, saya ingin **setiap peserta punya daftar putaran keputusannya
    sendiri**, supaya putaran satu peserta tidak tercampur dengan peserta lain. `[terverifikasi]`
48. Sebagai **`ReasLifeAdmin`**, saya ingin **memilih peserta mana** dari premium list yang akan
    diklaim, dan pilihan itu **terekam**, supaya dapat dipertanggungjawabkan.
    `[keputusan work owner]`
49. Sebagai **`ReasLifeSPV`**, saya ingin melihat tanggal diterima, tanggal konfirmasi, dan tanggal
    penyelesaian **per peserta**, supaya peserta yang tertahan terlihat — bukan tertutup satu
    tanggal untuk seluruh klaim. `[keputusan work owner]`
50. Sebagai **organisasi**, saya ingin data polis yang melekat pada klaim adalah **snapshot saat
    klaim dibuat**, supaya perubahan polis kemudian tidak mengubah klaim yang sedang berjalan.
    `[keputusan work owner]`
    ⚠️ **RALAT 2026-09-18 — US 50 DIBALIK.** Data polis **bukan lagi snapshot**. Perubahan
    polis sesudah klaim dibuat **akan** mengubah tampilan klaim lama, dan untuk polis
    ber-endorsement itu **pasti terjadi**. ⚠️ **penyimpangan sadar**, diterima work owner.
    Lihat §2b RALAT A.
51. Sebagai **organisasi**, saya ingin atribut polis tersimpan **sekali per klaim**, supaya
    mengubah satu atribut tidak menuntut menyentuh setiap baris peserta.
    ⚠️ **RALAT 2026-09-18 — US 51 KOSONG ARTINYA.** Atribut polis **tidak lagi disimpan di
    klaim sama sekali**, jadi tidak ada yang "tersimpan sekali per klaim". Lihat §2b RALAT A.
52. Sebagai **`ReasLifeAdmin`**, saya ingin mengunggah dokumen **per peserta**, dan saya ingin
    **dicegah menyimpan ke Outstanding** bila masih ada peserta yang berkasnya belum lengkap —
    dengan pesan yang menyebut **peserta mana**. `[terverifikasi]`
53. Sebagai **organisasi**, saya ingin **status akseptasi dan tanggal akseptasi berisi nilai
    sebenarnya** menurut aksi yang terjadi, bukan nilai yang distempel otomatis pada setiap insert.
    `[keputusan work owner]`
54. Sebagai **`ReasLifeAdmin`**, saya ingin menghapus klaim **beserta seluruh isinya** setelah
    **diberi peringatan berisi jumlah baris yang akan ikut terhapus**, supaya saya dapat
    membatalkan. `[keputusan work owner]`
55. Sebagai **sistem produksi/Arasapas**, saya ingin **membaca langsung dari tabel klaim**, supaya
    saya tidak perlu mengurai dokumen JSON. `[keputusan work owner]`
56. Sebagai **tim integrasi**, saya ingin bentuk yang saya baca **stabil dan bertipe**, supaya
    perubahan tampilan tidak memecahkan integrasi.
57. Sebagai **tim migrasi**, saya ingin nilai uang dan tanggal yang hari ini berupa teks menjadi
    **tipe yang semestinya**, dan yang **tidak dapat diurai dilaporkan** — bukan didiamkan atau
    diam-diam menjadi nol.
58. Sebagai **tim migrasi**, saya ingin atribut polis yang hari ini **berulang di setiap baris
    peserta** menjadi **satu baris** per klaim, dengan **perbedaan antar baris dilaporkan** bila
    ada.
59. Sebagai **Finance**, saya ingin **rekening pembayaran** (nama bank, id bank, nomor rekening)
    melekat pada baris adjustment, supaya pembayaran klaim punya tujuan yang jelas.
    `[terverifikasi]`
60. Sebagai **organisasi**, saya ingin sebuah baris **tidak dapat diserahkan ke Komite** selama
    rekening pembayarannya belum lengkap, supaya keputusan Komite tidak jatuh pada baris yang belum
    dapat dibayar. `[terverifikasi]` — paritas perilaku existing.

---

## Implementation Decisions

### 1. Batas konteks

Claim — Life dispesifikasikan sendiri; **Komite Life adalah sistem luar** (**ADR-0001**). Batasnya
tiga kontrak: penyerahan *child work*, jalur balik hasil keputusan, dan tabel akseptasi bersama.
Spesifikasi ini **berhenti** di ketiga batas itu — isi dan tangga internal Komite tidak dibentuk di
sini.

### 2. Penempatan modul

`[terverifikasi]` Struktur mengikat ada di **`OUTPUT_HASIL_RNM/CLAUDE.md` §5**: kode Go dan React
di-scaffold **di dalam `OUTPUT_HASIL_RNM\`** — `cmd/`, `internal/{config,handlers,models,repository,services}`,
`pkg/utils/`, `frontend/`, `go.mod`, `Makefile` — dengan arah dependency
**`handlers` → `services` → `repository`**, tidak boleh terbalik dan tidak boleh memotong lapisan.
Belum satu pun di antaranya ada hari ini.

Modul yang dibangun untuk konteks ini:

| Lapisan | Tanggung jawab |
| --- | --- |
| `handlers` | Endpoint REST Claim — Life; penegakan peran (**ADR-0002**) sebelum memanggil service |
| `services` | Mesin status per baris adjustment; aturan turunan "klaim selesai"; orkestrasi efek keluar; penulisan jejak audit |
| `repository` | Akses `POOLDATA`; penulisan **delapan tabel klaim relasional** (§2b) dan `T_WORK_CLAIM`; pembacaan sumber snapshot |
| `models` | Agregat klaim Life: header → peserta → baris adjustment & dokumen (§2b) |
| `frontend` | Layar Register / Outstanding / Medical Check / Claim Analis; kontrol bergerbang peran |

`[pertanyaan terbuka]` **Pembagian paket domain di dalam `internal/` belum ditetapkan.** `CLAUDE.md`
§5 menetapkan lapisan, bukan bagaimana keempat lini Claim dibagi antar paket. Ini penting karena
lini Life **menyimpang secara struktural** dari tiga lini lain — dua buktinya ada di §"Bukti
struktural". Pemilik: **Lead Engineer**; diambil saat scaffolding, sebelum tiket pertama yang
menyentuh model data.

### 2b. Bentuk penyimpanan — **delapan tabel relasional + `T_WORK_CLAIM`**

> ⚠️ **RALAT 2026-09-18 — DUA TABEL DIHAPUS, JADI ENAM. Judul §2b di atas ("delapan tabel")
> dan seluruh pohon di bawahnya SUDAH TIDAK BERLAKU.** Teks lama dibiarkan utuh sebagai jejak;
> yang mengikat adalah ralat ini. Berlaku juga untuk **US 50 · US 51** dan **AC 32 · 35 · 37 ·
> 38 · 48 · 51** — masing-masing ditandai di tempatnya.

#### RALAT 2026-09-18 · A — `T_CLAIM_POLICY` dan `T_CLAIM_MARKETING` DIHAPUS

`[keputusan work owner]` **`T_CLAIM_POLICY` dan `T_CLAIM_MARKETING` DIHAPUS.** Bukan diganti nama —
tabelnya **memang tidak ada**. Skema klaim Life turun dari **8 tabel menjadi 6**: lima tabel klaim
(`T_GENERAL_CLAIM`, `T_CLAIMLF_PREMIUMLIST_DETAIL`, `T_CLAIMLF_ADJUSTMENT`,
`T_CLAIMLF_ADJUSTMENT_SPREADING`, `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`) ditambah
**`DOCUMENT_CLAIM`**. Data polis dan marketing **dibaca dari tabel polis** (`T_PREMIUM_LIST` dkk,
modul **PremiumList Life**) lewat **penunjuk di `T_GENERAL_CLAIM`**.

⛔ Jangan membuat `T_CLAIMLF_POLICY` maupun `T_CLAIMLF_MARKETING`. Penghapusan ini **bukan**
penggantian nama.

**Mengapa boleh dihapus** `[terverifikasi]`:

| Tabel | Bukti |
| --- | --- |
| `T_CLAIM_POLICY` | **27 dari 32 kolomnya tersedia di `T_PREMIUM_LIST`.** (32 = 23 kolom dasar + 9 kolom tambahan hasil audit; `ID`/`CLAIM_ID` tidak dihitung. Lima sisanya kehilangan rumah — tabel C di bawah.) |
| `T_CLAIM_MARKETING` | **Terbukti turunan, nol nilai asli milik klaim.** `Claim Life/Activity/SetMOClaim_Act.xml` mengisi `MarketingData` dari `PolicyDataLife` (`MOID` · `MarketingCode` · `MarketingName` · `TeamGroup` · `BranchCode` · `BranchName`) **atau** dari master marketing officer (`ASM-FW-GISFW-Int-marketingofficer`, lewat `Claim Life/ReportDefinition/BrowseMarketingOfficer_RD.xml`). Tidak satu pun nilainya lahir di klaim. |

⚠️ **Penyimpangan sadar — potret berubah menjadi baca hidup.** `[keputusan work owner]` Keputusan
lama menyebut `PolicyDataLife` sebagai **snapshot** — *"disalin, bukan baca live"* — dan itu **kini
dibalik**. Akibatnya nyata dan diterima: **polis yang berubah sesudah klaim dibuat akan mengubah
tampilan klaim lama.** Untuk polis ber-endorsement hal itu **pasti terjadi**, karena endorsement
membuat versi baru. Ini **penyimpangan sadar**, bukan paritas.

#### RALAT 2026-09-18 · B — kolom yang bertambah dan yang pindah

| Tabel | Aksi | Kolom |
| --- | --- | --- |
| `T_GENERAL_CLAIM` | **TAMBAH** | `CASEID_POLICY` · `POLICY_NO` · `ENDORSMENT_NO` — ketiganya ⚠️ `[terbuka]`, lihat C |
| `T_GENERAL_CLAIM` | **BUANG** | `PL_NUMBER` → diganti nama menjadi **`POLICY_NO`** |
| `T_GENERAL_CLAIM` | **BUANG** | `CREATE_OP` · `CREATE_OP_NAME` · `TGL_UPDATE` → **pindah** ke `T_WORK_CLAIM` |
| `T_GENERAL_CLAIM` | **BUANG** | `CASEID` → **pindah** ke `T_WORK_CLAIM` *(DIPUTUSKAN 2026-09-18, lihat C2)* |
| `T_WORK_CLAIM` | **TAMBAH** | `CREATE_OP` · `CREATE_OP_NAME` · `TGL_UPDATE` (pindahan) |
| `T_WORK_CLAIM` | **TAMBAH** | `CASEID` (pindahan) *(DIPUTUSKAN 2026-09-18)* |
| `T_WORK_CLAIM` | **TAMBAH** | `LINI` — kolom **ADA**; nilai Life = konstanta lini Life. `[terbuka — Non-Life]` hanya **daftar enum lintas-lini** |

#### RALAT 2026-09-18 · C — yang MASIH TERBUKA. Jangan dijawab sendiri.

**C1 · Ketiga penunjuk polis belum menunjuk apa pun yang ada.** Ketiganya ditulis apa adanya di
atas, **bukan** karena sudah beres:

| Penunjuk | Mengapa belum menunjuk |
| --- | --- |
| `CASEID_POLICY` | Menunjuk **kolom apa?** `T_PREMIUM_LIST` **sengaja membuang** `pyID`/`px*`/`py*` — catatannya berbunyi *"ID work Pega diganti ID sequence baru"*. **Case id Pega milik polis tidak disimpan di sisi sana.** |
| `POLICY_NO` | **Tidak ada di header polis.** Ia ada di `T_PREMIUM_LIST_DETAIL` — **per peserta**. Header polis punya `PL_NUMBER`. Join header-klaim → header-polis lewat `POLICY_NO` **tidak nyambung.** |
| `ENDORSMENT_NO` | **Belum menunjuk satu versi.** Di sisi polis namanya **`NOENDORS`** `[terverifikasi]` (`Endorsement Life/RDBList/Generate_NoEndorsmentLife.xml`), dan **versi berjalan ditentukan `PRODKE` terbesar** — bukan oleh nomor endorsement saja. |

**C2 · Nasib `CASEID` — ✅ DIPUTUSKAN 2026-09-18, `[terbuka]` DITUTUP.**
`[keputusan work owner]` `CASEID` **PINDAH** dari `T_GENERAL_CLAIM` ke `T_WORK_CLAIM`. Alasannya
**sama** dengan `CREATE_OP`/`CREATE_OP_NAME`/`TGL_UPDATE`: ia identitas **work object**, bukan
atribut klaim. `T_GENERAL_CLAIM` **tidak lagi memuat** `CASEID`.

> Teks lama: *"`CASEID` masih tertinggal di `T_GENERAL_CLAIM` … Ikut pindah atau tetap? Belum
> diputuskan."* — **dicabut**, sudah dijawab.

**C3 · Lima kolom kehilangan rumah — JANGAN DITEBAK.** Inilah 5 dari 32 kolom `T_CLAIM_POLICY` yang
**tidak** tersedia di `T_PREMIUM_LIST`:

| Kolom | Keadaan |
| --- | --- |
| `TEAM_GROUP` | Tidak ada di `T_PREMIUM_LIST`; **ada di master marketing officer** — diambil lewat `MO_ID` |
| `BUSINESS_ID` | Tidak ada di `T_PREMIUM_LIST`; dipakai sebagai **parameter** `Claim Life/RDBList/Generate_NoAccept_Life.xml`. **Sumbernya wajib ditetapkan** |
| `TANGGAL_RESPON` · `TANGGAL_REALISASI` · `TANGGAL_KONFIRMASI_BALIK` | **Tanggal proses klaim, bukan atribut polis.** Tidak ada di `T_PREMIUM_LIST`, dan **tidak lagi ada di klaim**. Rumahnya `T_GENERAL_CLAIM` atau `T_WORK_CLAIM` — **belum diputuskan** |

**C4 · ✅ PEMBLOKIR TIKET 14 DICABUT 2026-09-18 — relasi nomor 2 dipecahkan lewat SHARED PK.**
`[keputusan work owner]` `T_GENERAL_CLAIM` dan baris klaim `T_WORK_CLAIM` memakai **`ID` yang sama
persis** — **shared primary key**, 1:1. **Tidak ada kolom `WORK_CLAIM_ID`**; hubungannya dijamin
oleh **ID yang identik**, bukan oleh kolom penyambung. Begitu pula `T_GENERAL_KOMITE.ID` =
`T_WORK_CLAIM.ID` baris komite. Rinciannya di **RALAT F** di bawah.

> Teks lama: *"⚠️ MEMBLOKIR TIKET 14 … tidak ada `WORK_CLAIM_ID` di `T_GENERAL_CLAIM`, tidak ada
> `CLAIM_ID` di `T_WORK_CLAIM` … Tiket 14 tidak dapat dinyatakan selesai sebelum ini dijawab."*
> — **dicabut**. Relasi nomor 2 **bukan lagi** alasan tiket 14 belum selesai.

#### RALAT 2026-09-18 · D — pohon yang berlaku: enam tingkat, berakar di `T_WORK_CLAIM`

`T_WORK_CLAIM` **naik menjadi akar**. Ia bukan lagi tabel di samping pohon klaim; klaim dan kasus
komite **sama-sama** baris di dalamnya, dibedakan oleh `COVER_KEY`.

```
TINGKAT 1   T_WORK_CLAIM ─────────────────── akar · LINTAS-LINI · satu baris per work object
            PK  ID
            FK  COVER_KEY → T_WORK_CLAIM.ID   (menunjuk dirinya sendiri; NULL bila tak punya induk)
            LINI · PY_POSITION · ACCEPT_STATUS · SENDTO_ADMIN · SENDTO_MEDICAL · TYPE
            CASEID · CREATE_OP · CREATE_OP_NAME · TGL_UPDATE  <- PINDAHAN dari header klaim
            |
   +--------+-------------------------------------------+
   | baris KLAIM  COVER_KEY = NULL                      | baris KOMITE  COVER_KEY = ID baris klaim
   |                                                    |
TINGKAT 2                                           TINGKAT 2
   +--1:1-- T_GENERAL_CLAIM                             +--1:1-- T_GENERAL_KOMITE
           PK ID = T_WORK_CLAIM.ID   <- SHARED PK               PK ID = T_WORK_CLAIM.ID baris komite
             (tidak ada kolom FK terpisah)                        <- SHARED PK, tidak ada kolom FK
           penunjuk polis (ke LUAR, bukan anak):                FK ADJUSTMENT_ID -> T_CLAIMLF_ADJUSTMENT.ID
             CASEID_POLICY   [terbuka]                            (tutup lingkar)
             POLICY_NO       [terbuka]                          KOMITE_LOOP · KOMITE_COUNT · ACCEPT_STATUS
             ENDORSMENT_NO   [terbuka]                          |
           |                                         TINGKAT 3  |
TINGKAT 3  |                                                    +--1:N-- T_KOMITE_KOMITELIST
           +--1:N-- T_CLAIMLF_PREMIUMLIST_DETAIL                        PK ID
           |        PK ID                                              FK DATA_KOMITE_ID
           |        FK CLAIM_ID -> T_GENERAL_CLAIM.ID  CASCADE             -> T_GENERAL_KOMITE.ID  CASCADE
           |        |                                                   KOMITE_URUT · KOMITE_ID · ID_KOMITE
TINGKAT 4  |        +--1:N-- T_CLAIMLF_ADJUSTMENT                       KOMITE_EMAIL · KOMITE_APROVAL
           |        |        PK ID                                      KOMITE_COMMENT · DATE_APPROVE
           |        |        FK PREMIUM_LIST_DETAIL_ID
           |        |           -> T_CLAIMLF_PREMIUMLIST_DETAIL.ID  CASCADE
           |        |        FK KOMITE_ID -> T_WORK_CLAIM.ID  <- penunjuk balik KE ATAS, nullable
           |        |        |
TINGKAT 5  |        |        +--1:N-- T_CLAIMLF_ADJUSTMENT_SPREADING
           |        |                 FK ADJUSTMENT_ID -> T_CLAIMLF_ADJUSTMENT.ID  CASCADE
           |        |                 |
TINGKAT 6  |        |                 +--1:N-- T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO
           |        |                          FK SPREADING_ID
           |        |                             -> T_CLAIMLF_ADJUSTMENT_SPREADING.ID  CASCADE
TINGKAT 4  |        +--1:N-- DOCUMENT_CLAIM                      <- LINTAS-LINI
           |                 FK PREMIUM_LIST_DETAIL_ID -> ...DETAIL.ID   (Life saja)
           |                 ON DELETE di Go -- induk beda tabel per lini
           |
           +-- DIBACA dari luar, BUKAN anak, TIDAK ikut cascade:
                 T_PREMIUM_LIST dkk (polis) · master marketing officer · EMAILKOMITE ·
                 master retro · security reinsurer · currency
```

✅ **Ejaan penaut induk-anak — DIPUTUSKAN 2026-09-18, `[terbuka]` DITUTUP.**
`[keputusan work owner]` **Ejaan final `COVER_KEY` (snake_case);** `CoverKey` warisan Pega
**tidak dipakai** sebagai nama kolom.
Alasan: seluruh kolom SQL baru proyek ini snake_case (`DATA_KOMITE_ID`, `ADJUSTMENT_ID`,
`PREMIUM_LIST_DETAIL_ID`, dst), dan Oracle melipat identifier tanpa kutip menjadi huruf besar —
sehingga
`CoverKey` akan menjadi identifier **tanpa garis bawah**, berbeda dari kolom-kolom sekitarnya.
Konsistensi SQL menang. Seluruh artefak kedua modul sudah diseragamkan.

#### RALAT 2026-09-18 · E — sebelas relasi

| # | Induk | Anak | Kunci tamu | Kard | Hapus |
| --- | --- | --- | --- | --- | --- |
| 1 | `T_WORK_CLAIM` | `T_WORK_CLAIM` | `COVER_KEY` | 1:N | di Go |
| 2 | `T_WORK_CLAIM` | `T_GENERAL_CLAIM` | **tidak ada kolom terpisah** — `T_GENERAL_CLAIM.ID` = `T_WORK_CLAIM.ID` (**shared PK**) | 1:1 | — |
| 3 | `T_GENERAL_CLAIM` | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `CLAIM_ID` | 1:N | CASCADE |
| 4 | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `T_CLAIMLF_ADJUSTMENT` | `PREMIUM_LIST_DETAIL_ID` | 1:N | CASCADE |
| 5 | `T_CLAIMLF_ADJUSTMENT` | `T_CLAIMLF_ADJUSTMENT_SPREADING` | `ADJUSTMENT_ID` | 1:N | CASCADE |
| 6 | `T_CLAIMLF_ADJUSTMENT_SPREADING` | `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` | `SPREADING_ID` | 1:N | CASCADE |
| 7 | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `DOCUMENT_CLAIM` | `PREMIUM_LIST_DETAIL_ID` | 1:N | **di Go** |
| 8 | `T_WORK_CLAIM` | `T_GENERAL_KOMITE` | **tidak ada kolom terpisah** — `T_GENERAL_KOMITE.ID` = `T_WORK_CLAIM.ID` baris komite (**shared PK**) | 1:1 | di Go |
| 9 | `T_GENERAL_KOMITE` | `T_KOMITE_KOMITELIST` | `DATA_KOMITE_ID` | 1:N | CASCADE |
| 10 | `T_CLAIMLF_ADJUSTMENT` | `T_GENERAL_KOMITE` | `ADJUSTMENT_ID` | 1:1 | di Go |
| 11 | `T_CLAIMLF_ADJUSTMENT` | `T_WORK_CLAIM` | `KOMITE_ID` | N:1 | penunjuk |

**Seluruh kunci tamu ber-index.** Relasi **2** dan **8** tidak punya kunci tamu untuk di-index —
keduanya **shared primary key**, dan PK sudah ber-index dengan sendirinya.

✅ **`DOCUMENT_CLAIM` — `ON DELETE` DIPUTUSKAN 2026-09-18, `[terbuka]` DITUTUP.**
`[keputusan work owner]` Penghapusannya **ditangani di Go**, **bukan** cascade basis data. Alasan:
ia **LINTAS-LINI** dan induknya **tabel yang berbeda per lini**, sehingga satu `ON DELETE CASCADE`
tidak dapat seragam. Untuk Life induknya tetap `T_CLAIMLF_PREMIUMLIST_DETAIL`.

#### RALAT 2026-09-18 · F — identitas: SHARED PK dan nomor bisnis berformat

**F1 · Shared primary key** `[keputusan work owner]`. Hubungan akar↔anak tingkat 2 **tidak memakai
kolom penyambung**. `T_GENERAL_CLAIM` dan baris klaim `T_WORK_CLAIM` memakai **`ID` yang sama
persis**; begitu pula `T_GENERAL_KOMITE` dan baris komite `T_WORK_CLAIM`. **Kolom `WORK_CLAIM_ID`
dibuang** dari seluruh rancangan — bukan diganti nama, **tidak ada**.

```
T_WORK_CLAIM  (akar, satu baris per work object)
  ID="CLM-123456"  COVER_KEY=NULL            <- baris KLAIM
  ID="KMT-000789"  COVER_KEY="CLM-123456"    <- baris KOMITE, menunjuk klaim induknya

T_GENERAL_CLAIM   ID="CLM-123456"    <- SHARED PK, = T_WORK_CLAIM baris klaim   (1:1)
T_GENERAL_KOMITE  ID="KMT-000789"    <- SHARED PK, = T_WORK_CLAIM baris komite  (1:1)
                  ADJUSTMENT_ID -> T_CLAIMLF_ADJUSTMENT.ID   (tutup lingkar, TETAP ada)
```

`ADJUSTMENT_ID` pada `T_GENERAL_KOMITE` **tetap ada** — ia penutup lingkar ke baris adjustment, dan
bukan bagian dari shared PK.

**F2 · Tipe identitas: TEKS BERFORMAT, bukan sequence** `[keputusan work owner]`.
`T_WORK_CLAIM.ID` adalah **teks berformat**: baris klaim `CLM-xxxxxx` (contoh `CLM-123456`), baris
komite `KMT-xxxxxx` (contoh `KMT-000789`). Kolom yang **mengikuti tipe teks ini**:

| Kolom | Isi |
| --- | --- |
| `T_WORK_CLAIM.ID` | `CLM-xxxxxx` (klaim) atau `KMT-xxxxxx` (komite) |
| `T_WORK_CLAIM.COVER_KEY` | `T_WORK_CLAIM.ID` baris induk; `NULL` bila tak punya induk |
| `T_GENERAL_CLAIM.ID` | **shared PK** — sama dengan baris klaim |
| `T_GENERAL_KOMITE.ID` | **shared PK** — sama dengan baris komite |
| `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` | `T_WORK_CLAIM.ID` baris komite |

⚠️ **Penyimpangan sadar dari ADR-0006 — dicatat, bukan dilanggar diam-diam.**
`[keputusan work owner]` **ADR-0006** menetapkan identitas dari **sequence**. Untuk `T_WORK_CLAIM`
dan kedua tabel ber-shared-PK, identitasnya adalah **nomor bisnis berformat** (`CLM-`/`KMT-`),
**bukan** sequence murni. Ini **penyimpangan sadar khusus tabel-tabel itu**; ADR-0006 tetap berlaku
untuk identitas tabel klaim lainnya (`T_CLAIMLF_*`, `DOCUMENT_CLAIM`).

⚠️ `[terbuka]` **Dua hal TETAP terbuka — tipe sudah diputuskan, dua ini belum. Jangan tebak:**

| `[terbuka]` | Pemilik |
| --- | --- |
| **Generator** nomor `CLM-`/`KMT-`: siapa yang membuatnya, apakah ada sequence di belakang prefiks, apakah di-reset per tahun | **DBA / work owner** |
| Apakah `T_WORK_CLAIM.COVER_KEY` dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` dipasangi `REFERENCES T_WORK_CLAIM(ID)`, atau dibiarkan tanpa constraint | **DBA / work owner** |

**F3 · Yang ditutup dan yang tidak, per 2026-09-18.**

| Butir | Status |
| --- | --- |
| Relasi 2 tanpa kolom (pemblokir tiket 14) | ✅ **DITUTUP** — shared PK (F1) |
| Relasi 8 memakai `WORK_CLAIM_ID` | ✅ **DITUTUP** — shared PK, kolomnya dibuang (F1) |
| Nasib `CASEID` | ✅ **DITUTUP** — pindah ke `T_WORK_CLAIM` (C2) |
| `ON DELETE` `DOCUMENT_CLAIM` | ✅ **DITUTUP** — di Go (RALAT E) |
| Tipe `T_WORK_CLAIM.ID` | ✅ **DITUTUP** — teks berformat (F2) |
| Isi `LINI` | ✅ kolom **ADA**, nilai Life = konstanta lini Life; `[terbuka — Non-Life]` hanya daftar enum |
| **Generator** `CLM-`/`KMT-` | ⚠️ `[terbuka]` |
| `REFERENCES` untuk `COVER_KEY`/`KOMITE_ID` | ⚠️ `[terbuka]` |
| **Ketiga penunjuk polis** (C1) | ⚠️ `[terbuka]` — **tidak disentuh** |
| **Lima kolom yatim** (C3) | ⚠️ `[terbuka]` — **tidak disentuh** |
| **`PREMIUM_SPREADED_NET` dua rumus** | ⚠️ `[terbuka]` — **tidak disentuh**, milik Product + UW |


⚠️ **Penyimpangan sadar 1 — seluruh JSON dibuang.** `[keputusan work owner]` Tidak ada blob JSON di
sistem baru. `JSON_KLAIM` dan serialisasi `ClaimData` dibuang; isi klaim disimpan **relasional**.
Alasannya bukan selera: bentuk JSON tidak dapat divalidasi, tidak dapat di-index, dan memaksa setiap
pembaca hilir mengurai dokumen.

`[keputusan work owner]` Klaim punya **tabelnya sendiri**, terpisah dari tabel existing.

⚠️ **Yang dibuang HANYA JSON.** `[keputusan work owner]` Sistem baru menulis **dua tempat**:
tabel relasional di bawah, **dan** `INSERT` flat ke **`OS_AKSEPTASI_KLAIM_LIFE`** — karena hilir
(Arasapas / produksi) masih membaca dari sana. `M_LIFE_PREMIUM_DETAIL` tetap **dibaca saja**,
sebagai sumber snapshot peserta.

```
T_GENERAL_CLAIM (PK ID)                          header — PremiumListSummary (1:1) dilipat ke sini
  ├─ T_CLAIM_POLICY               1:1   FK CLAIM_ID → T_GENERAL_CLAIM.ID              ON DELETE CASCADE
  ├─ T_CLAIM_MARKETING            1:1   FK CLAIM_ID → T_GENERAL_CLAIM.ID              ON DELETE CASCADE
  └─ T_CLAIMLF_PREMIUMLIST_DETAIL  1:N   FK CLAIM_ID → T_GENERAL_CLAIM.ID              ON DELETE CASCADE
        ├─ T_CLAIMLF_ADJUSTMENT     1:N   FK PREMIUM_LIST_DETAIL_ID → …DETAIL.ID    ON DELETE CASCADE
        │     └─ T_CLAIMLF_ADJUSTMENT_SPREADING        1:N  FK ADJUSTMENT_ID → …ADJUSTMENT.ID   ON DELETE CASCADE
        │            └─ T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO  1:N  FK SPREADING_ID → …SPREADING.ID  ON DELETE CASCADE
        └─ DOCUMENT_CLAIM         1:N   FK PREMIUM_LIST_DETAIL_ID → …DETAIL.ID    ON DELETE CASCADE

T_WORK_CLAIM (PK ID)   ⬅ tabel work mandiri, LINTAS-LINI (Life + Non-Life) — bukan anak T_GENERAL_CLAIM
```

⚠️ **Penyimpangan sadar 2 — induk `T_CLAIMLF_ADJUSTMENT` adalah PESERTA, bukan klaim.**
`[terverifikasi]` `Claim Life/Activity/SavePesertaClaim.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEPESERTACLAIM` / `RULE-OBJ-ACTIVITY`) menulis delapan field ke
`…PremiumListSummary.PremiumListDetail(<LAST>).AdjustmentList(<LAST>).<field>` — **setiap peserta
punya daftar adjustment sendiri**. Skema kandidat semula menggantungkannya ke header; itu **menghapus
informasi peserta pemilik** dan mematahkan mesin status (**ADR-0011**) yang menjadikan baris
adjustment sebagai unit. Ini **perbaikan relasi**, bukan peniruan.

#### Isi tiap tabel

**`T_GENERAL_CLAIM`** — skalar `ClaimData` + `PremiumListSummary` yang dilipat:
`CLAIM_NO`, `PY_ID`, `STS_KATASTROFE`, `KATASTROFE_NOTE`, `IS_KPR`, `STNC_CLAIM`, `ACCEPTED_NO`,
`STS_REJECT`, `RISLIPRNM`, `PL_NUMBER`, `BUSINESS_NAME`, `CASEID`, `CLAIM_RETRO`, `CREATE_OP`,
`CREATE_OP_NAME`, `TGL_UPDATE`.

`[terverifikasi]` `PremiumListSummary` **tidak pernah ber-subscript** di seluruh modul → **1:1**
dengan klaim, karena itu dilipat. `PL_NUMBER` ikut ke header — sumbernya
`PolicyDataLife.PremiumListSummary.PL_NUMBER`, jadi ia milik **summary**, bukan milik polis.

⚠️ Aturan pemilihan kolom **diperlebar**: `CASEID` dan ketiga kolom audit **bukan** property
`ClaimData` — asalnya work object dan tabel warisan. Aturan sebenarnya: **kolom = nilai yang dipakai
dan perlu bertahan**, dari mana pun asalnya.

> ⚠️ **RALAT 2026-09-18 — `T_CLAIM_POLICY` DIHAPUS**, dan kata **snapshot** di bawah
> **DICABUT**: potret berubah menjadi **baca hidup** ⚠️ **penyimpangan sadar**. **27 dari 32**
> kolomnya tersedia di `T_PREMIUM_LIST`; **lima** kehilangan rumah — termasuk ketiga
> `TANGGAL_*` dan `TEAM_GROUP` di daftar ini ⚠️ `[terbuka]`. Lihat §2b RALAT A dan C3.

**`T_CLAIM_POLICY`** (1:1) — snapshot `PolicyDataLife`, **+8 kolom** hasil audit:
`TANGGAL_RESPON`, `TANGGAL_REALISASI`, `TANGGAL_KONFIRMASI_BALIK`, `PRODUCT_NAME_ID`,
`PRODUCT_NAME`, `TEAM_GROUP`, `DATE_RECEIVED`, `BRANCH_NAME`, `BRANCH_CODE`; **−`PL_NUMBER`**.
`[terverifikasi]` Ketiga `TANGGAL_*` tampil di **empat section aktif** — `InputRegisterClaimLife`,
`InputOSClaimLife`, `InputAkseptasiClaimLife`, `MedicalCheckClaimLife` — **tanpa gerbang
visibilitas**.

`[keputusan work owner]` Isinya **snapshot** saat klaim dibuat — disalin, bukan dibaca live.
Perubahan polis kemudian **tidak** mengubah klaim berjalan.

> ⚠️ **RALAT 2026-09-18 — `T_CLAIM_MARKETING` DIHAPUS.** Seluruh kolom di bawah **terbukti
> turunan**: `Claim Life/Activity/SetMOClaim_Act.xml` mengisinya dari `PolicyDataLife` atau
> dari master marketing officer — **nol nilai asli milik klaim**. `TEAM_GROUP` adalah satu
> dari lima kolom yang ⚠️ `[terbuka]` kehilangan rumah. Lihat §2b RALAT A dan C3.

**`T_CLAIM_MARKETING`** (1:1) — `MARKETING_ID`, `CLIENT_ID`, `CLIENT_NAME`, `BRANCH_DETAIL_ID`,
`BRANCH_DETAIL_NAME`, `TEAM_GROUP`.

**`T_CLAIMLF_PREMIUMLIST_DETAIL`** (1:N) — satu baris per peserta yang diklaim, **+9 kolom** hasil
audit: **`IS_CHECK`**, `STATUS`, `RECOMMENDATION`, `STS_REJECT`, `SOURCE_ID`, `CONFIRMATION_DATE`,
`COMPLETE_DATE`, `CLAIM_RECEIVED_DATE`, `CEDING_RETENTION`.

⚠️ **Penyimpangan sadar 3 — `IS_CHECK` eksplisit.** `[keputusan work owner]` Aturan
*"`PremiumListSummary` = hanya peserta yang diklaim"* tidak punya penyimpan tanpa kolom ini. Karena
klaim menyimpan **snapshot**, `IS_CHECK`-lah yang merekam siapa yang dipilih pada saat pembacaan.

⚠️ `[keputusan work owner]` Ketiga tanggal berada **di peserta, bukan di header** — di situlah
korpus menulisnya.

**`T_CLAIMLF_ADJUSTMENT`** (1:N di bawah peserta) — **+4 kolom** hasil audit: `CLAIM_AMOUNT`,
`STS_REJECT`, `ACCEPTEDNO`, `ACCEPTATION_DATE`; **+3 kolom bank**: `NAME_OF_BANK`, `ID_BANK`,
`ACCOUNT_NO`; **+1 kolom rujukan**: **`KOMITE_ID`**.

⚠️ **Penyimpangan sadar 9 — roster & keputusan komite TIDAK disimpan di Claim Life.**
`[keputusan work owner]` `KOMITE_ID` adalah **rujukan** (nullable — `NULL` bila baris belum pernah
dikirim ke Komite). Roster dan keputusan per anggota milik konteks **Komite Claim Life**; Claim Life
**melihat**nya lewat join. **Tidak ada `T_CLAIMLF_ADJUSTMENT_KOMITE`.**

`[terverifikasi]` Di Pega justru sebaliknya — `KomiteList` adalah **anak `AdjustmentList`**:
`Komite Claim Life/Activity/KomitePostAdjustment.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` /
`KOMITEPOSTADJUSTMENT`) menulis `…AdjustmentList(idx).KomiteList(k).KomiteAproval` /
`.KomiteComment` / `.DateApprove`; `Claim Life/Activity/CreateKMTLife_Act.xml` mengisi
`childPageKomite.KomiteList`, `IndexPremiumList`, dan `IndexAdjustment = .pxListSubscript`.

⚠️ **Penyimpangan sadar 10 — rujukan memakai ID stabil, bukan indeks posisi.**
`[keputusan work owner]` Pega merujuk baris lewat **subscript posisi**, yang rusak begitu urutan
bergeser. Sistem baru memakai **`ADJUSTMENT_ID`** dan **`KOMITE_ID`**.

**`T_CLAIMLF_ADJUSTMENT_SPREADING`** (1:N di bawah **baris adjustment**) — hasil spreading per
treaty-year: `ADJUSTMENT_ID` (FK), `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`, `TREATY_YEAR_LIFE`,
`RETROCADED_SHARE`, `RATE`, `IDR`, `USD`, `CURRENCY`.

**`T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`** (1:N di bawah spreading) — per reinsurer: `SPREADING_ID`
(FK), `REINSURER_NAME`, `PERCENT_SHARE`, `AMOUNT`, `RATE`, `PREMIUM_SPREADED_GROSS`,
`PREMIUM_SPREADED_NET`, `COMMISION` (sic), `OVR_COMM`, `TREATY_TYPE_ID`, `TREATY_TYPE_NAME`.
Seluruh nilai uang dan persen **desimal** (**ADR-0003**).

⚠️ **Penyimpangan sadar 11 — spreading adjustment DIBEKUKAN dan DISIMPAN.**
`[keputusan work owner]` `[terverifikasi]` `AdjustmentList` (class
`ASM-FW-GISFW-Data-AdjustmentLife`) punya anak `SpreadingList` → `RetroLifeList`; keduanya diisi
`Claim Life/Activity/SpreadingClaimLife_Act.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
`SPREADINGCLAIMLIFE_ACT`), dengan kelas baris `ASM-FW-GISFW-Int-TREATYYEAR_LIFE` dan
`ASM-FW-GISFW-Int-RETROCESSIONLIFE`. Nilainya disimpan apa adanya; perubahan master treaty
sesudahnya **tidak** mengubahnya.

⚠️ `[terbuka]` **`PREMIUM_SPREADED_NET` punya dua rumus** di rule yang sama — `GROSS − Comm` dan
`GROSS − Discount − Comm`. Mana yang berlaku **tidak terbaca dari korpus**; rumusnya ditetapkan
Product + UW. **Jangan tebak.** Rincian dan AC-nya di tiket **14** dan **03**.

`[terverifikasi]` Ketiga field bank milik class **`ASM-FW-GISFW-Data-AdjustmentLife`** —
`.NameOfBank`, `.IDOfBank`, `.NoAccount` — dan tampil di
`Claim Life/Section/AdjustmentDetail_Section.xml`. Kolom fisiknya sudah ada di tabel warisan:
`NAME_OF_BANK`, `IDBANK`, `ACCOUNTNO` di
`Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml`.

⚠️ **Ketiganya adalah GERBANG penyerahan ke Komite — paritas, bukan penyimpangan.**
`[terverifikasi]` `Claim Life/Activity/GetListKomiteLife.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`) memuat prasyarat
`.NameOfBank=="" || .NoAccount=="" || .IDOfBank==""` dengan pesan
`"Name of bank cannot be empty"`. Perilaku ini **ditiru apa adanya**: baris adjustment yang salah
satu field banknya kosong **tidak dapat diserahkan ke Komite**.

⚠️ **Penyimpangan sadar 4 — `STS_REJECT` dan `ACCEPTATION_DATE` diisi NILAI SEBENARNYA.**
`[keputusan work owner]` `[terverifikasi]` INSERT warisan
(`Claim Life/RDBList/InsertJsonKlaimLife_sql.xml`, `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` /
`RNM!INSERTJSONKLAIMLIFE_SQL`) **meng-hardcode** `ACCEPTATION_DATE = SYSDATE` dan
`STS_REJECT = '0'` — keduanya tidak di-bind. Sistem baru **tidak meniru**:

| Aksi | `STS_REJECT` |
| --- | --- |
| `ReasLifeAdmin` insert ke Outstanding | `0` |
| `ReasLifeAdmin` reject langsung | `2` |
| `ReasLifeSPV` tambah baris Outstanding | `0` |

`ACCEPTATION_DATE` diisi **tanggal akseptasi sebenarnya** saat baris benar-benar diaksep — bukan
distempel `SYSDATE` pada setiap insert.

**`DOCUMENT_CLAIM`** (1:N di bawah peserta) — **tabel baru**.

⚠️ **Penyimpangan sadar 5 — dokumen jadi tabel sendiri, bukan lampiran bawaan.**
`[keputusan work owner]` `[terverifikasi]` `Claim Life/Activity/SaveOutStandingLife_Act.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`) beriterasi atas
`.DocumentList` **per peserta** dan **menolak simpan** dengan
`"The document hasn't been uploaded person number "+<nomor>` serta
`"Documents are incomplete, please complete the documents"`. Dokumen adalah **gerbang simpan**,
bukan pelengkap — karena itu ia tabel, dan pembacaannya `SELECT` biasa.

**`T_WORK_CLAIM`** — tabel work mandiri.

⚠️ **Penyimpangan sadar 6 — keadaan tangga TIDAK disimpan di header klaim.**
`[keputusan work owner]` `[terverifikasi]` `pyPosition` (42 rujukan), `AcceptStatus`, `SendtoAdmin`
(9), `SendtoMedical` (6), dan `Type` (25) hidup di **work object**, bukan di `ClaimData`. Keduanya
dipakai bersama **Life dan Non-Life**, sehingga ditempatkan di tabel **lintas-lini** tersendiri —
bukan dilipat ke `T_GENERAL_CLAIM`.

`[terbuka — Non-Life]` Relasi dan cascade formal `T_WORK_CLAIM` **ditetapkan saat konteks Non-Life
digarap**. Untuk Claim Life: menghapus klaim life **menghapus juga** baris work-nya. **Tidak
memblokir** konteks ini.

#### Tipe, kunci, dan kaskade

⚠️ **Penyimpangan sadar 7.** `[keputusan work owner]` Uang/share/premi → **desimal presisi
arbitrer** (**ADR-0003**, tidak pernah `float`); tanggal → **`DATE`**; **seluruh kolom nullable**
(wajib-isi ditegakkan di Go); identitas via **sequence** (**ADR-0006**).

⚠️ **Penyimpangan sadar 8 — kaskade + popup konfirmasi.** `[keputusan work owner]` Seluruh FK antar
kedelapan tabel klaim **`ON DELETE CASCADE`**. Menghapus klaim menghapus polis, marketing, peserta,
seluruh adjustment, **seluruh spreading dan spreading retro**, dan seluruh dokumen — **setelah popup
konfirmasi Ya/Batal** yang menyebut jumlah baris tiap jenis (pola Master Contract Retro Life &
Treaty Contract Out). ⚠️ Kaskade menyentuh **lima tingkat**; uji sampai **cicit**
(`T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`).

#### Yang dibaca, bukan dimiliki

- Data polis life / PremiumList NB/EDM → sumber snapshot `T_CLAIM_POLICY` dan peserta.
  ⚠️ **RALAT 2026-09-18** — baris ini DICABUT. `T_CLAIM_POLICY` dihapus; tidak ada lagi
  snapshot. Data polis **dibaca hidup** dari `T_PREMIUM_LIST` dkk lewat penunjuk di
  `T_GENERAL_CLAIM` ⚠️ **penyimpangan sadar**. Lihat §2b RALAT A.
- ⚠️ **`product_life` relasional** (hasil migrasi Master Product Name Life) → `PRODUCT_NAME` /
  `PRODUCT_NAME_ID`. `[keputusan work owner]` **Bukan** `m_product_life.JSONDATA`.
- Master retro, security reinsurer, currency → dibaca saat pengisian.
- `[terverifikasi]` **Master diagnosa** (class `ASM-FW-GISFW-DATA-DIAGNOSELIFE`, 132 rujukan) adalah
  **master pencarian**, bukan sub-object klaim: hasilnya mendarat sebagai **nilai tunggal**
  `DISEASE` / `ICD_CODE` pada peserta. **Tidak** menjadi tabel.

### 3. Mesin status — unit adalah **baris `AdjustmentList`**

Ditetapkan **ADR-0011**. Ringkasnya, untuk dipakai sebagai kontrak:

| # | Langkah | Pelaku | Akibat pada baris |
| --- | --- | --- | --- |
| 1 | Input baris pertama → Save ke Outstanding | `ReasLifeAdmin` | `0` |
| 1b | Reject Outstanding atas baris yang ia input — tanpa Komite; **membatalkan baris itu saja** | `ReasLifeAdmin` | `2` |
| 2 | Submit ke Medical Check | → `ReasLifeMedicalAdvisor` | tetap `0` |
| 3 | Submit ke SPV | → `ReasLifeSPV` | tetap `0` |
| 4 | Send ke Komite (semua adjustment) | `ReasLifeSPV`, atau **siapa pun bila `Type` `TP`/`TR`** | tetap `0` |
| 5 | Komite memutus | Komite Life (luar) | aksep → `1`; tolak → `2` |
| 6 | Setelah tolak: tambah baris **baru** → Save ke Outstanding → send Komite. Berulang. | `ReasLifeSPV` (atau `ReasLifeAdmin` setelah reject sendiri) | baris baru `0` |

**Aturan yang mengikat:**

- **`STS_REJECT` adalah status baris**, bukan status klaim. Nilai: `0` Outstanding, `1` Aksep,
  `2` Ditolak. ⚠️ Namanya menyesatkan — `1` berarti **diaksep**. Antarmuka **tidak boleh**
  menampilkan nama field ini.
- **Baris bersifat terminal.** Sekali `1` atau `2`, tidak berubah. `[terverifikasi]` tidak ada rule
  di korpus yang menulis `0` setelah `1`/`2`.
- **Klaim tidak terminal.** Revisi dilakukan dengan **menambah baris baru**, bukan mengubah baris
  lama.
- **Nilai `2` selalu berarti "baris ini ditolak", tidak pernah "klaim selesai"** — berlaku sama
  untuk penolakan Admin maupun Komite. `[keputusan work owner]`
- **`PremiumListDetail.STS_REJECT` adalah cerminan** baris adjustment terakhir, bukan unit keputusan
  tersendiri. Header klaim mengikuti adjustment terakhir. Di sistem baru ia **turunan**, bukan
  kolom yang ditulis mandiri. `[keputusan work owner]`; `[terverifikasi]` di Pega kedua tingkat
  ditulis berbarengan dengan nilai sama oleh `RejectOSClaimLife_Act` dan `KomitePostAdjustment`.

**Aturan turunan "klaim selesai" — ditetapkan di sini, bukan diwarisi.** Karena tidak ada nilai
status yang berarti "klaim selesai", keadaan itu **dihitung**:

> Sebuah klaim **selesai** bila tidak ada lagi baris `AdjustmentList` bernilai `0`
> **dan** terdapat sekurang-kurangnya satu baris bernilai `1`.
> Bila tidak ada baris bernilai `0` dan tidak ada pula yang bernilai `1`, klaim berada dalam keadaan
> **ditolak seluruhnya** — dan tetap dapat dilanjutkan dengan baris baru.

`[keputusan work owner]` untuk premisnya (nilai `2` tidak pernah berarti klaim selesai);
**aturan turunan di atas adalah keputusan spesifikasi** — konsekuensi logis dari premis itu, bukan
temuan korpus. Ia **wajib ditinjau work owner** sebelum diimplementasikan.

### 4. Aktor dan penegakan peran

Tiga peran (**ADR-0002**): `ReasLifeAdmin`, `ReasLifeMedicalAdvisor`, `ReasLifeSPV`. Rangkap peran
tidak diperbolehkan kecuali ditambahkan eksplisit pada role akun.

**Wewenang kirim ke Komite bergantung `Type`** (**ADR-0012**) — dibawa sebagai **paritas**:

| `Type` | Siapa yang boleh mengirim |
| --- | --- |
| `TP` (Payable) / `TR` (Receivable) | **siapa pun**, `ReasLifeAdmin` termasuk |
| `QP` / `QR` | **hanya `ReasLifeSPV`** |

⚠️ **Risiko RBAC yang diterima**, tercatat di **ADR-0012**: wewenang bergantung pada atribut data,
bukan peran. Ditinjau ulang saat konteks Komite Life / IAM digarap — **bukan** diperketat sekarang.

**Penegakan ada di lapisan layanan, bukan di visibilitas layar.** `[terverifikasi]` Di Pega
penegakan tidak seragam: `pyPosition` tidak ada sama sekali di `Section/InputOSClaimLife.xml`,
`Section/RejectOSClaimLife_Sec.xml`, `Section/ClaimComite.xml`, dan `Harness/Committe_Life.xml`;
yang menegakkan sesungguhnya adalah penugasan tahap di `Claim Life/Flow/Register_Flow.xml`.
Ketidakseragaman itu **tidak ditiru**.

### 5. `Type` menjadi satu field

`[terverifikasi]` Di Pega nilai yang sama hidup dalam **dua salinan**:

| Yang digerbangi | Properti yang dibaca |
| --- | --- |
| Wewenang send-Komite (`AdjustmentDetail_Section.xml`) | `pyWorkPage.Type` |
| Jendela validasi DOL (`ValidasiDOL_Act.xml`) | `pyWorkPage.PolicyDataLife.Type` |

Penyalinannya di `Claim Life/Activity/LoadDataPeserta_Act.xml` —
`pyWorkPage.Type = pyWorkPage.PolicyDataLife.Type`. **Sumber otoritatif = `PolicyDataLife.Type`**
(tipe polis).

**Keputusan: di sistem baru keduanya menjadi satu field**, bersumber dari data polis. Keduanya —
wewenang dan validasi — membaca field yang sama, sehingga tidak dapat berbeda seperti di Pega.

`Type` karena itu **menyentuh keamanan**, bukan sekadar data: perubahan nilainya wajib masuk jejak
audit (**ADR-0007**).

### 6. Validasi Date of Loss

`[terverifikasi]` `Claim Life/Activity/ValidasiDOL_Act.xml`:

| Cabang | Jendela | Argumen `@addCalendar` |
| --- | --- | --- |
| `Type=="QR" \|\| Type=="QP"` | `GROSS_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `(.DATE_OF_LOSS,0,0,0,0,0,0,0)` |
| `Type=="TR" \|\| Type=="TP"` | `RETROCESSION_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `(.DATE_OF_LOSS,0,0,0,1,0,0,0)` |

Gagal → pesan `"Invalid DOL"`, dengan `local.Begin==false || local.Expired==true`.

`[dugaan]` Argumen keempat bernilai `1` pada cabang `TP`/`TR` (nol pada `QP`/`QR`). Menurut tanda
tangan `@addCalendar` Pega yang lazim itu berarti **+1 hari** — tetapi definisi fungsinya **tidak
ada di korpus**, jadi satuannya belum terbukti. **Pergeseran ini dibawa apa adanya sebagai paritas**;
bila kelak dinyatakan keliru, itu keputusan terpisah. Satuannya perlu dikonfirmasi sebelum
implementasi.

### 7. Jenis klaim dari kode produk

`[terverifikasi]` `ContentNote` (jenis klaim) diturunkan dari `BusinessCode`;
`Claim Life/Activity/SaveOutStandingLife_Act.xml` menguji `L1`…`L11` dalam satu precondition.
`[keputusan work owner]` Daftar lengkap `L1`–`L21` → `DEATH` / `HEALTH` / `CI` / `TPD` / `TI` ada di
`CONTEXT.md`.

**Keputusan: pemetaan ini menjadi data acuan, bukan `if` berderet.** `[terverifikasi]` Modul Pega
hanya menguji `L1`–`L11`; `L12`–`L21` ada di master produk, bukan di modul klaim — sehingga
menyalin precondition apa adanya akan membawa daftar yang sudah tidak lengkap.

### 8. Kontrak dengan Komite Life

Ditetapkan **ADR-0001**, tiga kontrak. Yang mengikat implementasi:

- **Penyerahan** membuat *child work* berkelas `ASM-FW-GCNMFW-Work-KomiteLife`, membawa penunjuk
  baris (`IndexAdjustment`, `IndexPremiumList`), `KomiteLoop`, dan daftar anggota komite.
  `[keputusan work owner]` **ditambah**: nilai klaim, `CURRENCY`, dan `STS_REJECT` saat penyerahan.
- **Jalur balik** mengubah `STS_REJECT` baris. `[terverifikasi]` perubahan hanya terjadi pada
  **rung terakhir** — seluruh penulisan di `KomitePostAdjustment` digerbangi
  `pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop`.
- **Tabel akseptasi bersama** `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`, ditulis **satu rule yang sama**
  dari kedua sisi (`…!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL`, hash ternormalisasi `c50bfd9a12`
  identik).
- `[keputusan work owner]` `KomiteLoop` **ditentukan Claim — Life**; apa yang menentukan nilainya
  `[terbuka]` **OQ-032**.

### 9. Uang

**ADR-0003**: delapan kolom uang, `EM_PERCENT` persen, `CURRENCY` kode mata uang. Tidak ada `float`
di lapisan mana pun maupun di kontrak API. `[terverifikasi]` nama kolom tidak dapat dipakai menebak
sifatnya — empat kolom bernama `SHARE_*` / `*_SHARE` ternyata **uang**, bukan rasio.

`[terverifikasi]` `CURRENCY` dan `CURRENCYID` ada di **tingkat baris** dan disalin dari baris 1 oleh
`SetIndexAdjustmentList` — jadi dalam praktiknya seragam per klaim, meski strukturnya membolehkan
campur. Bentuk tipe akhirnya bergantung **OQ-060**.

Bentuk tipe konkretnya (tipe desimal berskala tetap, atau bilangan bulat satuan minor) ditetapkan
saat scaffolding; **ADR-0003** mengikat hanya syaratnya: presisi desimal eksak, bukan biner
floating-point, di seluruh lapisan dan di kontrak API.

### 10. Penomoran

**ADR-0006**: panggil `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`; **jangan replikasi logikanya** —
badan procedure tidak ada di korpus, mereplikasinya berarti menebak. Dua rule penomoran lama
(`Generate_NoKlaim_Life`, `Generate_NoKlaim_LifeRetro`) **tidak dimigrasikan**; `[terverifikasi]`
keduanya tidak terindeks sebagai rujukan aktif di `SaveOutStandingLife_Act`.

Konsekuensi: sistem baru **bergantung pada Oracle** untuk penomoran. Kontraknya `[terbuka]`
**OQ-002**.

### 11. Environment

**ADR-0005**: `IsPEGAPROD` menjadi flag lingkungan.

⚠️ **Konsekuensi yang wajib disadari:** `[terverifikasi]` `IsPEGAPROD` **juga menggerbangi simpan
utama** — bukan hanya efek samping. Di `Komite Claim Life/Activity/KomitePostAdjustment.xml`
precondition `IsPEGAPROD` muncul 3×, berdampingan dengan langkah penyimpanan.

**Keputusan: flag lingkungan hanya menggerbangi EFEK KELUAR, tidak pernah menggerbangi
penyimpanan.** Di lingkungan non-production, klaim tetap tersimpan; yang tidak berjalan adalah
unggah berkas, email, Arasapas, dan konversi. Ini **penyimpangan sadar** dari perilaku Pega —
tanpa itu, lingkungan non-production tidak dapat dipakai menguji apa pun.

### 12. Integrasi dan efek keluar

**Tiga** efek keluar (**ADR-0008**) — konversi produksi lewat payload JSON dibuang
(lihat di bawah). Seluruhnya **asinkron, tidak memblokir, dengan antre-ulang**:

| Efek | Rule sumber |
| --- | --- |
| Unggah/unduh berkas ke Google Storage (**ADR-0010**) | `Claim Life/Activity/InsertGoogleStorage_Act.xml`, `GetUrlGoogleStorage_Act.xml`, `DeleteGoogleStorage_Act.xml` |
| Email | `Claim Life/Activity/SendEmailKlaimLF.xml` |
| Arasapas | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` — satu-satunya salinan di korpus (**OQ-035**) |

Alamat layanan **di-lookup runtime** dari `M_LINK_SERVICE` lewat kunci `(KATEGORI_1, KATEGORI_2)`
(**ADR-0013**; menggantikan ADR-0004). **Konsekuensi yang diterima:** sebuah klaim
mencapai keadaan akhir sementara efek keluarnya masih tertunda.

⚠️ `[keputusan work owner]` **Efek keluar tinggal tiga.** Konversi ke
produksi **tidak lagi mengirim payload JSON**: Arasapas/produksi **`SELECT` langsung dari tabel
klaim** (§2b). `convertJsonNusareToProductionClaimLife` dan seluruh serialisasi JSON keluar **tidak
dimigrasikan**. Konsekuensinya bentuk yang dibaca hilir menjadi **tabel bertipe**, bukan dokumen —
dan tiga efek keluar yang tersisa (berkas, email, Arasapas) **tidak berubah aturannya**.

**Keputusan: seluruh integrasi berada di balik interface Go sejak awal**, bukan panggilan langsung
dari service. Alasannya bukan kerapian melainkan kemungkinan pengujian — kontrak dua di antaranya
belum diketahui (**OQ-002** `GET_TOKEN_STORAGE`, **OQ-047** daftar endpoint), sehingga satu-satunya
cara menguji perilaku "tidak memblokir + antre-ulang" adalah menggantinya di batas itu.

### 13. Jejak audit

**ADR-0007**: siapa + kapan untuk **setiap** transisi status **dan setiap** jalur balik. Karena unit
status adalah baris, **jejak direkam per baris**, bukan per klaim.

`[terverifikasi]` Yang ada sekarang hanya `CREATEOPNAME` + empat kolom tanggal pada rekam akseptasi;
`SendtoAdmin` dan `SendtoMedical` tidak merekam pelaku sama sekali. `[keputusan work owner]` Jejak
lama **tidak dapat direkonstruksi ke belakang** — riwayat lengkap hanya ada untuk kejadian setelah
cutover.

### 14. Migrasi

**ADR-0009**: seluruh data dipindahkan; koeksistensi ditolak. **ADR-0011**: seluruh **baris** ikut
pindah, bukan hanya keadaan terakhir — kalau tidak, riwayat putaran Komite hilang.

⚠️ **Bentuknya kini berubah total** `[keputusan work owner]`: dari dokumen JSON + satu tabel flat
warisan → **delapan tabel relasional** (§2b). Karena itu **irisan skema + migrasi adalah PREFACTOR,
tiket pertama** — tidak ada irisan lain yang berdiri sebelum bentuk barunya ada.

**Sumber → tujuan.** `[terverifikasi]` `OS_AKSEPTASI_KLAIM_LIFE` menyimpan **satu baris per
peserta** dan **mengulang 14 atribut polis** (ceding, source of business, retro, security
reinsurer, produk, type) pada **setiap** baris peserta — bukti: daftar 51 kolom di
`Claim Life/RDBList/InsertJsonKlaimLife_sql.xml`. Migrasi karena itu wajib **menormalkan**: atribut
polis yang berulang menjadi **satu baris `T_CLAIM_POLICY`** per klaim.

> ⚠️ **RALAT 2026-09-18 — normalisasi ini tidak lagi berlaku.** Tidak ada `T_CLAIM_POLICY`
> untuk ditulisi, dan atribut polis yang berulang di `OS_AKSEPTASI_KLAIM_LIFE` **tidak
> dipindahkan ke mana pun** — ia dibaca dari tabel polis. Yang **tetap berlaku**: perbedaan
> nilai antar baris peserta dalam satu klaim **dilaporkan**. Lihat §2b RALAT A dan AC 51.

⚠️ Bila atribut polis **berbeda antar baris peserta** dalam satu klaim, migrasi **melaporkannya** —
tidak diam-diam memilih salah satu.

⚠️ `[terverifikasi]` Dua nilai di-hardcode di INSERT warisan (`ACCEPTATION_DATE = SYSDATE`,
`STS_REJECT = '0'`). Data lama karena itu **tidak memuat tanggal akseptasi sebenarnya**. Migrasi
**melaporkan** hal ini; ia **tidak menafsirkan** dan tidak mengarang tanggal.

⚠️ **Nama sumber migrasi menyesatkan (OQ-066).** `[terverifikasi]`
`Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` /
`RULE-CONNECT-SQL`) bernama "Update" tetapi isinya **`INSERT INTO POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`
dengan 55 kolom** — **bukan** `UPDATE`. Dicatat di sini supaya tim migrasi **tidak mencari jalur
`UPDATE` yang memang tidak ada**. Ketiga kolom bank (`NAME_OF_BANK`, `IDBANK`, `ACCOUNTNO`) berasal
dari rule ini.

### 15. Rule yang tidak dimigrasikan

| Rule | Alasan |
| --- | --- |
| `Claim Life/Activity/SaveAdjustment_Act.xml` | `[keputusan work owner]` **dead rule** — sudah tidak dipakai |
| `Generate_NoKlaim_Life`, `Generate_NoKlaim_LifeRetro` | `[terverifikasi]` tidak terindeks sebagai rujukan aktif |
| ⚠️ `Claim Life/ConnectREST/convertJsonNusareToProductionClaimLife.xml` | `[keputusan work owner]` **payload JSON dibuang** — hilir `SELECT` dari tabel (§2b, §12) |
| ⚠️ `Claim Life/RDBList/InsertJsonClaimLifeGCNM.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONCLAIMLIFEGCNM`) → `POOLDATA.PEGA_JSON_KLAIM_PNC` | `[keputusan work owner]` **tidak dipakai** — sistem baru tidak memanggil procedure JSON ini. Nama/kode cabang diambil **relasional** dan disimpan di `T_CLAIM_POLICY` |
| ⚠️ `Claim Life/RDBList/GetJsonProductLife.xml` (`ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `ASM!GETJSONPRODUCTLIFE`) **sebagai baca-JSON** | `[keputusan work owner]` ia `json_table` atas `m_product_life.JSONDATA`; digantikan pembacaan **`product_life` relasional** (hasil migrasi Master Product Name Life) |

> ⚠️ **RALAT 2026-09-18** — pada baris `InsertJsonClaimLifeGCNM` di tabel ini, frasa
> **"disimpan di `T_CLAIM_POLICY`"** DICABUT: tabelnya dihapus. Nama/kode cabang **dibaca
> dari tabel polis**, tidak disimpan di klaim. Keputusan **tidak memakai** procedure JSON itu
> **tidak berubah**. Lihat §2b RALAT A.

⚠️ **Titik tinjau.** Status "dead" `SaveAdjustment_Act` **tidak dapat diverifikasi dari korpus** —
yang terbaca justru sebaliknya: `[terverifikasi]` rule itu **terpasang di UI**, dirujuk 2× sebagai
`<pyActivity>` dari `Claim Life/Section/ClaimLifeDetailGCNM.xml` (824.562 byte), dan menulis
`.AdjustmentList(<LAST>).STS_REJECT = 1` tanpa precondition Komite. Bila kelak ditemukan jalur itu
masih dipakai di produksi, keputusan ini perlu ditinjau ulang — dan bukti di atas titik mulanya
(**ADR-0011** §`SaveAdjustment_Act`).

---

### 16. Membaca peserta klaim — **hanya peserta hidup**

⚠️ `[keputusan work owner]` **Penajaman kontrak hilir (2026-09-15, verdict V14 grilling Endorsement
Life).** Konteks Endorsement Life menulis **baris bernilai negatif** (jurnal balik) ke tabel yang
**sama** dengan yang dibaca klaim. Aturannya:

> **Peserta yang sudah dibatalkan (EDM Batal) atau di-soft-delete TIDAK BOLEH MUNCUL di Claim Life.**

**Satu tabel, dua sudut pandang — dan itu disengaja:**

| Pembaca | Melihat |
| --- | --- |
| **Akuntansi / ringkasan premium** | **seluruh** baris — positif **dan** negatif; nettonya dari penjumlahan |
| **Klaim (konteks ini)** | **hanya peserta hidup** — yang belum dibatalkan dan belum dihapus |

#### Celah yang harus ditutup `[terverifikasi]`

`Claim Life/RDBList/GetPesertaClaim_sql1.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` /
`RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL`) hari ini berbunyi:

```sql
SELECT * FROM POOLDATA.M_LIFE_PREMIUM_DETAIL
WHERE PL_NUMBER = {…PremiumListSummary.PL_NUMBER}
  AND CERTIFICATE_NO LIKE '%'||{…CARI2}||'%'
  AND UPPER(NAME_OF_INSURED) LIKE '%'||{…CARI3}||'%'
```

**Tidak ada penyaring status apa pun** — ia mengambil seluruh baris, termasuk baris negatif.
`[terverifikasi]` Sensus korpus: **nol** kemunculan `EDMSTATUS` di seluruh modul `Claim Life/`.

⚠️ Apakah Pega menyaringnya di lapisan lain **tidak terbukti dari korpus**. **Jangan menebak
mekanisme Pega.** Perlakukan ini sebagai **kontrak yang wajib ditegakkan sistem baru**, apa pun cara
Pega dahulu.

#### Kolom penandanya **terbaca dari korpus** `[terverifikasi]`

`M_LIFE_PREMIUM_DETAIL` memiliki **tiga** kolom berakhiran status. Bukti: daftar kolom `INSERT` di
`SaveMasterLPDet` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` /
`RULE-CONNECT-SQL`), beserta sumber nilainya di klausa `VALUES`:

| Kolom | Diisi dari | Nilai terbaca | Penanda hidup/mati? |
| --- | --- | --- | --- |
| **`EDMSTATUS`** | `TempValue.EDMStatus` | `Old` / `New` / `Delete` / `Batal` | ✅ **ya — inilah penandanya** |
| `STATUS` | `CARI48` | `0` untuk `QR`/`QP`, `1` untuk `TP`/`TR` | ❌ tidak — ia penanda **jenis transaksi** |
| `STATUSOLD` | `CARI47` | `1` bila "Statusold = Old", `0` bila bukan | ❌ tidak |

`[terverifikasi]` Nilai `STATUS` pada jalur new business: `@if(Type=="QR","0",@if(Type=="QP","0",
"1"))` di `InsertLifePremiumDetail_act` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT` /
`RULE-OBJ-ACTIVITY`) — sejalan dengan sub-langkah "Status QR/QP" (`0`) dan "Status TR/TP" (`1`) di
`Endorsement Life/Activity/InsertJsonPolisLife_Act.xml`.

⚠️ **Jebakan yang harus ditangani: baris new business tidak mengisi `EDMSTATUS` sama sekali.**
`[terverifikasi]` Sensus `InsertLifePremiumDetail_act`: **nol** kemunculan `EDMStatus`. Jadi baris
NB masuk dengan `EDMSTATUS` **kosong/NULL**. Penyaring yang ditulis naif sebagai
`EDMSTATUS NOT IN ('Delete','Batal')` akan **membuang seluruh peserta new business**, karena di
Oracle perbandingan dengan `NULL` tidak pernah bernilai benar.

**Aturan yang berlaku:** peserta hidup = `EDMSTATUS` **kosong/NULL**, `Old`, atau `New`. Peserta
mati = `Delete` atau `Batal`. Penyaring wajib memperlakukan **kosong/NULL sebagai hidup**.

#### Yang masih perlu DBA

`[terbuka]` **Tipe dan nullability kolom `EDMSTATUS`** belum terbaca — apakah benar `NULL` atau
justru string kosong `''` untuk baris new business. Keduanya menuntut penyaring yang berbeda.
Masuk **OQ-001 (sisa)**, pemilik **DBA**. **Tidak memblokir aturan bisnisnya** — hanya bentuk
akhir penyaringnya.

---

## Testing Decisions

### Apa yang membuat test baik di sini

Test menguji **perilaku yang dapat diamati dari luar modul**, bukan detail implementasi. Untuk
konteks ini artinya: test menyatakan **apa yang terjadi pada baris adjustment dan pada klaim**, dan
**siapa yang boleh melakukannya** — bukan bahwa suatu fungsi dipanggil atau suatu struct punya
bentuk tertentu.

Test yang hanya membenarkan implementasi — misalnya memeriksa bahwa repository dipanggil dengan
argumen tertentu — **tidak diterima**: ia mengunci bentuk kode, bukan perilaku, sehingga refaktor
yang benar pun membuatnya merah sementara bug perilaku lolos.

### Seam — satu, di titik tertinggi, dibuat dari nol

`[terverifikasi]` **Tidak ada seam lama yang dapat dipakai ulang: belum ada kode sama sekali.**
`OUTPUT_HASIL_RNM\` hari ini hanya berisi dokumen — `cmd/`, `internal/`, `pkg/`, `frontend/`,
`go.mod`, dan `Makefile` **belum satu pun ada**. Seluruh kode di-scaffold **baru di dalam
`OUTPUT_HASIL_RNM\`** (`CLAUDE.md` §5), sehingga seam pun dirancang dari nol — dan sengaja
**hanya satu**:

> **Seam utama: API HTTP Claim — Life.** Test menggerakkan sistem lewat endpoint REST dan memeriksa
> hasilnya lewat endpoint REST, dengan rangkaian `handlers → services → repository` terpasang
> sungguhan.

Alasan memilih titik ini:

- Ia **titik tertinggi** yang masih menghasilkan perilaku teramati ujung-ke-ujung — satu slice
  vertikal menghasilkan satu perilaku yang dapat diuji utuh.
- Mesin status per baris, wewenang per `Type`, dan validasi DOL **semuanya** teramati di sini.
- Ia **tidak mengunci pembagian paket di dalam `internal/`** — pembagian itu memang belum
  diputuskan (§Bukti struktural), sehingga menguji di bawahnya berarti mengunci sesuatu yang masih
  terbuka.

**Batas proses difake, bukan dijadikan seam kedua:**

| Batas | Perlakuan dalam test |
| --- | --- |
| Google Storage, email, Arasapas, konversi | *fake* di balik interface Go — keempatnya memang efek keluar asinkron (**ADR-0008**) dengan alamat di-lookup runtime dari `M_LINK_SERVICE` (**ADR-0013**); test memeriksa **efeknya** (tercatat, tidak memblokir, dapat diantre ulang), bukan panggilannya |
| Oracle | skema uji nyata — bukan mock. `[terverifikasi]` seluruh penulisan Claim — Life bermuara ke tabel `POOLDATA` lewat `RULE-CONNECT-SQL`, dan penomoran bergantung stored procedure (**ADR-0006**); memocknya berarti tidak menguji apa pun yang penting |
| Jam | dapat dikendalikan, karena validasi DOL dan jejak audit bergantung waktu |

**Satu seam kedua yang mungkin tidak terhindarkan** — perlu konfirmasi:
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` dan `GET_TOKEN_STORAGE` badannya tidak ada di korpus
(**OQ-002**). Selama kontraknya tidak diketahui, keduanya **tidak dapat difake dengan jujur**.
Pilihannya: (a) test memakai procedure sungguhan di skema uji, atau (b) seam kedua di batas
repository khusus untuk kedua procedure itu. **Usulan saya (a)** — satu seam tetap satu, dan
perilaku procedure ikut teruji. Perlu persetujuan work owner + DBA.

### Modul yang diuji

| Yang diuji | Lewat seam | Catatan |
| --- | --- | --- |
| Mesin status per baris (6 langkah + 2 sumber reject) | API HTTP | inti spesifikasi |
| Aturan turunan "klaim selesai" | API HTTP | **menunggu tinjauan work owner** |
| Wewenang per peran dan per `Type` | API HTTP | termasuk kasus `TP`/`TR` bebas peran |
| Validasi DOL per `Type` | API HTTP | dua cabang jendela |
| `ContentNote` dari `BusinessCode` | API HTTP | `L1`–`L21` |
| Jejak audit tiap transisi & jalur balik | API HTTP | perilaku baru (**ADR-0007**) |
| Efek keluar tidak memblokir + antre-ulang | API HTTP + fake | perilaku baru (**ADR-0008**) |
| Representasi uang tanpa `float` | API HTTP | termasuk bentuk kontrak API |
| ⚠️ **Induk adjustment = peserta** | API HTTP + Oracle nyata | dua peserta × dua putaran → empat baris tertelusur benar (AC 33, 34) |
| ⚠️ **Kaskade hapus tiga tingkat** | API HTTP + Oracle nyata | adjustment & dokumen **cucu** ikut hilang; popup menyebut jumlah (AC 48) |
| ⚠️ **Gerbang dokumen per peserta** | API HTTP + Oracle nyata | simpan Outstanding tertahan, pesan menyebut peserta mana (AC 45) |
| ⚠️ **Atomisitas lintas delapan tabel** | API HTTP + Oracle nyata | kegagalan pada anak mana pun membatalkan seluruhnya (AC 49) |
| ⚠️ **Normalisasi atribut polis saat migrasi** | skrip migrasi + Oracle nyata | perbedaan antar baris peserta **dilaporkan**, bukan diam-diam dipilih (AC 51) |
| ⚠️ **Keadaan tangga di `T_WORK_CLAIM`** | API HTTP + Oracle nyata | bukan kolom header klaim; hapus klaim life → baris work ikut (AC 46, 47) |

⚠️ **Oracle tidak pernah dipalsukan untuk ketujuh baris bertanda di atas.** Kaskade tiga tingkat,
atomisitas lintas delapan tabel, presisi desimal, dan konversi tanggal **hanya berperilaku benar pada
basis data sungguhan**; memalsukannya berarti tidak menguji apa pun yang penting. Yang di-*fake*
tetap hanya **klien efek keluar** (berkas, email, Arasapas).

### Prior art

`[terverifikasi]` **Tidak ada — nol test, nol kode** di `OUTPUT_HASIL_RNM\`. Slice Claim — Life ini
akan **menjadi** prior art bagi konteks berikutnya, sehingga bentuk testnya perlu dipilih dengan
sadar sejak awal: apa pun yang ditulis di sini akan disalin oleh slice sesudahnya.

**Perintah verifikasi wajib ditulis eksplisit di tiap tiket** selama `Makefile` belum ada. Setelah
toolchain di-scaffold, target akhirnya `go test ./internal/...` untuk backend dan
`cd frontend && npm test` untuk frontend, sesuai struktur di `CLAUDE.md` §5.

---

## Acceptance Criteria

Dapat diverifikasi tanpa pengetahuan pribadi.

**Mesin status**

1. Baris adjustment yang baru disimpan ke Outstanding berstatus `0`.
2. Baris yang sudah `1` atau `2` tidak dapat berubah lagi melalui jalur mana pun.
3. Penolakan oleh `ReasLifeAdmin` membuat **hanya baris itu** bernilai `2`; baris lain dan klaimnya
   tidak berubah.
4. Penolakan oleh Komite membuat baris yang diserahkan bernilai `2`, dan klaim tetap dapat menerima
   baris baru.
5. Baris baru yang ditambahkan setelah penolakan berstatus `0`, dan mewarisi delapan kolom dari
   baris pertama **tanpa** mewarisi status.
6. Perubahan status oleh Komite hanya terjadi ketika putaran mencapai rung terakhir.
7. `PremiumListDetail` dan header klaim selalu mencerminkan baris adjustment terakhir.
8. Klaim dilaporkan "selesai" **hanya** bila tidak ada baris bernilai `0` dan ada sekurangnya satu
   bernilai `1`. *(menunggu tinjauan work owner — §Implementation Decisions 3)*

**Wewenang**

9. `ReasLifeMedicalAdvisor` tidak dapat mengubah status akseptasi baris mana pun.
10. Untuk klaim ber-`Type` `QP` atau `QR`, hanya `ReasLifeSPV` yang dapat mengirim ke Komite;
    upaya oleh peran lain ditolak di lapisan layanan.
11. Untuk klaim ber-`Type` `TP` atau `TR`, `ReasLifeAdmin` **dapat** mengirim ke Komite.
12. Penolakan wewenang terjadi di lapisan layanan, dan tetap terjadi meskipun kontrol UI-nya
    ditampilkan.

**Validasi**

13. Klaim ber-`Type` `QP`/`QR` dengan *Date of Loss* di luar `GROSS_VALUATION_BEGIN_DATE` …
    `_EXPIRED_DATE` ditolak dengan pesan yang setara `"Invalid DOL"`.
14. Klaim ber-`Type` `TP`/`TR` diuji terhadap `RETROCESSION_VALUATION_*`, dengan pergeseran tanggal
    yang sama seperti Pega. *(satuan pergeseran menunggu konfirmasi — §Implementation Decisions 6)*
15. `ContentNote` terisi sesuai tabel `BusinessCode` `L1`–`L21` di `CONTEXT.md`.

**Audit**

16. Setiap transisi status baris menghasilkan catatan berisi pelaku dan waktu.
17. Setiap pengembalian (`SendtoAdmin`, `SendtoMedical`) menghasilkan catatan berisi pelaku dan
    waktu.
18. Perubahan nilai `Type` menghasilkan catatan berisi pelaku dan waktu.

**Efek keluar dan lingkungan**

19. Kegagalan efek keluar mana pun tidak menahan transisi status klaim.
20. Kegagalan efek keluar tercatat di jalur audit, bukan hanya di log layanan, dan dapat diantre
    ulang.
21. Di lingkungan non-production, klaim **tetap tersimpan**; **ketiga** efek keluar tidak berjalan.


**Uang**

22. Tidak ada nilai uang yang direpresentasikan sebagai *binary floating point* di lapisan mana pun
    maupun di kontrak API.
23. `EM_PERCENT` tidak diperlakukan sebagai uang.
24. Nilai klaim yang menyeberang ke Komite membawa mata uangnya.

**Membaca peserta klaim — hanya peserta hidup** `[keputusan work owner]`

25. Pencarian peserta untuk klaim **tidak menampilkan** peserta yang sudah dibatalkan
    (`EDMSTATUS = 'Batal'`) maupun yang di-soft-delete (`EDMSTATUS = 'Delete'`). *(§16; verdict V14
    grilling Endorsement Life)*
26. Peserta **new business** — yang `EDMSTATUS`-nya **kosong/NULL** karena jalur NB tidak mengisinya
    — **tetap muncul**. Test wajib memuat kasus ini: penyaring naif `NOT IN ('Delete','Batal')`
    membuang seluruh peserta NB di Oracle dan **harus gagal**.
27. Peserta ber-`EDMSTATUS` **`Old`** dan **`New`** **tetap muncul** — keduanya peserta hidup.
28. Baris **bernilai negatif** hasil jurnal balik endorsement **tidak pernah** sampai ke layar
    klaim maupun ke perhitungan klaim.
29. Penyaringan terjadi di **satu tempat** — lapisan repository pembaca peserta — bukan tersebar di
    services atau handlers. Test yang menemukan jalur baca peserta tanpa penyaring **gagal**.
30. `STATUS` dan `STATUSOLD` **tidak** dipakai sebagai penanda hidup/mati; `STATUS` adalah penanda
    jenis transaksi (`0` untuk `QR`/`QP`, `1` untuk `TP`/`TR`). Test yang menyaring dengan kolom itu
    **gagal**.

**Penyimpanan relasional** — §2b, §14

31. ⚠️ **Tidak ada blob JSON** sebagai penyimpan isi klaim di lapisan mana pun. Test yang menemukan
    kolom JSON menyimpan atribut klaim **gagal**. *(§2b; penyimpangan sadar 1)*
32. ⚠️ Menyimpan klaim menulis **dua tempat**: **kedelapan tabel klaim** *dan* `INSERT` flat ke
    **`OS_AKSEPTASI_KLAIM_LIFE`** — hilir masih membaca dari sana. Test yang **menolak** penulisan
    `OS_AKSEPTASI_KLAIM_LIFE` justru **gagal**. `M_LIFE_PREMIUM_DETAIL` tetap **dibaca saja**.
    Yang **dibuang hanya JSON**. *(§2b; `[keputusan work owner]` — koreksi 2026-09-16)*
    ⚠️ **RALAT 2026-09-18** — **"kedelapan tabel klaim"** menjadi **ENAM**. Sisa AC 32 —
    `OS_AKSEPTASI_KLAIM_LIFE` tetap ditulis, yang dibuang hanya JSON — **tidak berubah**.
    Lihat §2b RALAT A.
33. ⚠️ Setiap baris `T_CLAIMLF_ADJUSTMENT` menunjuk **satu peserta** lewat `PREMIUM_LIST_DETAIL_ID`.
    Test yang menemukan FK adjustment menunjuk **header klaim** **gagal**. *(§2b; penyimpangan
    sadar 2)*
34. ⚠️ Dua peserta dengan masing-masing dua putaran adjustment menghasilkan **empat baris yang
    seluruhnya dapat ditelusuri ke peserta yang benar**. *(§2b; penyimpangan sadar 2)*
35. Atribut polis tersimpan **sekali per klaim**. Test yang menemukan atribut polis berulang per
    peserta **gagal**. *(§2b)*
    ⚠️ **RALAT 2026-09-18 — AC 35 KOSONG ARTINYA.** Atribut polis **tidak lagi disimpan di
    klaim sama sekali**, jadi "sekali per klaim" tidak lagi punya yang dihitung. Yang
    menggantikan: atribut polis **dibaca hidup** dari tabel polis ⚠️ **penyimpangan sadar**.
    AC ini **tidak dihapus** dan **tidak diganti** — ia menunggu keputusan work owner.
    Lihat §2b RALAT A.
36. Keempat field `PremiumListSummary` — `CLAIM_NO`, `PL_NUMBER`, `RISLIPRNM`, `BUSINESS_NAME` —
    berada di **header**, tidak terpecah ke tabel lain. *(§2b)*
37. `T_CLAIM_POLICY` menyimpan **tanggal respon, tanggal realisasi, tanggal konfirmasi balik**,
    **produk** (id + nama), **team group**, **tanggal diterima**, dan **cabang** (kode + nama).
    *(§2b)*
    ⚠️ **RALAT 2026-09-18 — AC 37 KOSONG ARTINYA.** `T_CLAIM_POLICY` **dihapus**; tabel yang
    disebut AC ini tidak ada. Lebih dari itu, **empat dari tujuh** hal yang didaftarnya —
    `TANGGAL_RESPON`, `TANGGAL_REALISASI`, `TANGGAL_KONFIRMASI_BALIK`, dan **team group** —
    ⚠️ `[terbuka]` **kehilangan rumah**: bukan pindah, tetapi belum punya tempat. AC ini
    **tidak dihapus**. Lihat §2b RALAT A dan C3.
38. ⚠️ `PRODUCT_NAME` / `PRODUCT_NAME_ID` diisi dari **`product_life` relasional**. Test yang
    menemukan pembacaan `m_product_life.JSONDATA` **gagal**. *(§2b, §15; penyimpangan sadar 1)*
    ⚠️ **RALAT 2026-09-18 — AC 38 berpindah pemilik, tidak batal.** `PRODUCT_NAME` /
    `PRODUCT_NAME_ID` **bukan lagi kolom klaim** — keduanya termasuk 27 kolom yang tersedia di
    `T_PREMIUM_LIST`, sehingga larangan membaca `m_product_life.JSONDATA` kini mengikat
    **modul PremiumList Life**, bukan Claim Life. Lihat §2b RALAT A.
39. ⚠️ Peserta menyimpan **penanda dipilih-untuk-diklaim** (`IS_CHECK`). Fitur "hanya peserta yang
    diklaim" **tidak dapat dinyatakan selesai** tanpa kolom ini. *(§2b; penyimpangan sadar 3)*
40. ⚠️ Tanggal **diterima**, **konfirmasi**, dan **penyelesaian** tersimpan **per peserta**, bukan
    di header. *(§2b)*
41. Peserta menyimpan `STATUS`, `RECOMMENDATION`, `SOURCE_ID`, dan `CEDING_RETENTION`. *(§2b)*
42. ⚠️ `STS_REJECT` diisi **nilai sebenarnya menurut aksi**: Admin insert Outstanding → `0`; Admin
    reject langsung → `2`; SPV tambah Outstanding → `0`. Test yang menemukan nilai di-hardcode
    **gagal**. *(§2b; penyimpangan sadar 4)*
43. ⚠️ `ACCEPTATION_DATE` diisi **tanggal akseptasi sebenarnya** saat baris diaksep — **bukan**
    `SYSDATE` pada setiap insert. *(§2b; penyimpangan sadar 4)*
44. ⚠️ Dokumen tersimpan **per peserta** di `DOCUMENT_CLAIM`, dan dapat dibaca dengan `SELECT`
    biasa. *(§2b; penyimpangan sadar 5)*
45. ⚠️ Menyimpan ke Outstanding **ditolak** bila ada peserta yang dokumennya belum lengkap, dengan
    pesan yang **menyebut peserta mana**. *(§2b; penyimpangan sadar 5)*
46. ⚠️ Keadaan tangga klaim — posisi, status akseptasi, pengembalian ke Admin/Medical, dan `Type` —
    tersimpan di **`T_WORK_CLAIM`**, **bukan** sebagai kolom `T_GENERAL_CLAIM`. *(§2b; penyimpangan
    sadar 6)*
47. Menghapus klaim life **menghapus juga** baris `T_WORK_CLAIM`-nya. *(§2b)*
48. ⚠️ Menghapus klaim **mengkaskade** ke polis, marketing, peserta, seluruh adjustment, dan seluruh
    dokumen — didahului **popup konfirmasi Ya/Batal** yang menyebut **jumlah baris tiap jenis**;
    **Batal** tidak mengubah apa pun. Seluruhnya dalam **satu transaksi**. *(§2b; penyimpangan
    sadar 8)*
    ⚠️ **RALAT 2026-09-18** — kata **"polis, marketing"** DICABUT dari daftar kaskade: kedua
    tabelnya dihapus. Sisa AC 48 **tetap berlaku**, dan kaskadenya kini menyentuh peserta,
    seluruh adjustment, **spreading**, **spreading retro**, dan dokumen.
    Lihat §2b RALAT A dan E.
49. Satu klaim beserta seluruh anaknya ditulis dalam **satu transaksi**; kegagalan di mana pun
    **membatalkan seluruhnya**. *(§2b)*
50. ⚠️ Seluruh nilai uang dan share bertipe **desimal presisi arbitrer**; seluruh tanggal bertipe
    `DATE`; seluruh kolom **nullable** dengan wajib-isi ditegakkan di Go; identitas dari
    **sequence**. *(§2b; **ADR-0003**, **ADR-0006**; penyimpangan sadar 7)*
51. ⚠️ Migrasi **menormalkan** atribut polis yang hari ini berulang di setiap baris peserta menjadi
    **satu baris** per klaim; **perbedaan nilai antar baris dalam satu klaim dilaporkan**, tidak
    diam-diam diambil salah satu. *(§14)*
    ⚠️ **RALAT 2026-09-18 — AC 51 KOSONG ARTINYA.** `T_CLAIM_POLICY` dihapus, jadi tidak ada
    lagi "satu baris per klaim" untuk dinormalkan ke mana pun; atribut polis **tidak
    dipindahkan ke mana pun** — ia dibaca dari tabel polis. Yang **tetap berlaku**: bila
    atribut polis berbeda antar baris peserta dalam satu klaim, migrasi **melaporkannya**.
    AC ini **tidak dihapus**. Lihat §2b RALAT A.
52. ⚠️ Migrasi **melaporkan** bahwa data lama ber-`ACCEPTATION_DATE` = waktu insert dan
    `STS_REJECT` = `0` karena **di-hardcode**, bukan karena nilainya sebenarnya. *(§14)*
53. Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara tepat**,
    bukan dengan toleransi. *(§14; **ADR-0003**)*
54. Tanggal yang berupa teks menjadi `DATE` **tanpa pergeseran zona waktu**; yang **tidak dapat
    diurai dilaporkan**, bukan didiamkan. *(§14)*
55. ⚠️ Sistem hilir membaca **langsung dari tabel klaim**; **tidak ada** payload JSON yang dikirim.
    *(§12; penyimpangan sadar 1)*
56. Baris adjustment menyimpan **nama bank**, **id bank**, dan **nomor rekening**. *(§2b;
    `[terverifikasi]` class `ASM-FW-GISFW-Data-AdjustmentLife`)*
57. **Penyerahan ke Komite ditolak** bila salah satu dari ketiga field bank pada baris yang
    diserahkan **kosong**, dengan pesan yang setara `"Name of bank cannot be empty"`. Ini
    **paritas perilaku existing**, bukan penyimpangan. *(§2b; `[terverifikasi]`
    `GetListKomiteLife`, prasyarat `.NameOfBank=="" || .NoAccount=="" || .IDOfBank==""`)*
58. ⚠️ **`T_CLAIMLF_ADJUSTMENT_SPREADING` dan `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` ada**, masing-masing
    menggantung pada **baris adjustment** dan pada **baris spreading** — bukan pada peserta maupun
    header klaim. Test yang menemukan induk yang salah **gagal**. *(§2b; penyimpangan sadar 11;
    rincian di tiket 14)*
59. ⚠️ Menyimpan baris adjustment **menghitung dan menyimpan** spreading-nya, dan nilainya
    **dibekukan**: perubahan master treaty sesudahnya tidak mengubah angka yang sudah tersimpan.
    *(§2b; penyimpangan sadar 11; tiket 03)*
60. ⚠️ Kaskade hapus menyentuh **lima tingkat** — klaim → peserta → adjustment → spreading →
    spreading retro. Test wajib memeriksa **cicit** ikut hilang. *(§2b; penyimpangan sadar 8)*
61. ⚠️ `T_CLAIMLF_ADJUSTMENT` menyimpan **`KOMITE_ID`** (nullable) sebagai **rujukan**; roster dan
    keputusan komite **tidak** disimpan di konteks ini, dan **tidak ada**
    `T_CLAIMLF_ADJUSTMENT_KOMITE`. *(§2b; penyimpangan sadar 9; tiket 10)*
62. ⚠️ Rujukan antar baris memakai **ID stabil**, bukan indeks posisi. Test yang menemukan padanan
    `IndexPremiumList` / `IndexAdjustment` sebagai kunci rujukan **gagal**. *(§2b; penyimpangan
    sadar 10)*

---

## Pertanyaan terbuka di dalam spec

Sepuluh OQ **non-pemblokir**. Spec tetap ditulis; bagian yang bergantung padanya ditandai di atas.
**Jangan ditebak.** Nomor diambil dari `discovery/open-questions.md`.

| OQ | Pertanyaan | Pemilik | Bagian spec yang menunggu |
| --- | --- | --- | --- |
| **OQ-020** (sisa) | Arti `QP` dan `QR`. Perilakunya sudah terbaca; hanya namanya belum. `[dugaan]` huruf pertama memisahkan gross (`Q*`) dari retrosesi (`T*`) | Product+UW | §4 tabel wewenang, §6 |
| **OQ-032** | Apa yang menentukan `KomiteLoop` — berapa putaran sebelum keputusan final | Product+UW | §8 |
| **OQ-037** | Ambang `Param.LIMIT_BOTTOM` yang menentukan komposisi roster Komite | Product+UW | §8 |
| **OQ-060** | Apakah keseragaman mata uang per klaim aturan atau kebetulan implementasi | Product+UW | §9, AC 24 |
| **OQ-035** | Kepemilikan `serviceInsertArasapasClaimLife_act` — satu salinan, dua konteks | Product+UW | §12 |
| ~~**OQ-001**~~ | ~~Tidak ada DDL — tipe, presisi, skala kolom; bentuk penyimpanan baris adjustment~~ | — | ✅ **DITUTUP 2026-09-16** — tabel klaim **baru dirancang sendiri** (§2b); tipe, kunci, dan kaskade ditetapkan di sana. Sisa DDL tabel **existing** yang dibaca sebagai sumber snapshot tetap DBA, **tidak memblokir** |
| ~~**OQ-002**~~ (bagian penyimpanan) | ~~Kontrak procedure penulis~~ | — | ✅ **DITUTUP 2026-09-16** — jalur tulis klaim **tidak lagi lewat procedure**; Go menulis tabelnya sendiri. Kontrak `PROC_GENERATE_SEQUENCE_NUMBER` (§10) dan `GET_TOKEN_STORAGE` (§12) **tetap terbuka**, tidak memblokir |
| **OQ-013** | Batas transaksi — `COMMIT` berada di dalam blok PL/SQL | DBA | §13 atomisitas |
| **OQ-047** | Isi `M_LINK_SERVICE` — daftar endpoint sebenarnya | DBA / IT-infra | §12 |
| **OQ-018** (sisa) | Lingkungan `pega-nusre`; bucket production vs dev | IT-infra | §11, §14 |

Tiket yang bergantung pada OQ di atas **bukan** `ready-for-agent` sampai OQ-nya terjawab — aturan
`docs/agents/triage-labels.md`.

**Satu OQ baru dari revisi penyimpanan — tidak memblokir konteks ini:**

| OQ | Pertanyaan | Pemilik | Bagian yang menunggu |
| --- | --- | --- | --- |
| `[terbuka — Non-Life]` | **Relasi dan cascade formal `T_WORK_CLAIM`.** Ia tabel work **lintas-lini** (Life + Non-Life), jadi bentuk relasinya ditetapkan saat konteks Non-Life digarap. Untuk Claim Life yang mengikat hanya: klaim life dihapus → baris work-nya ikut (AC 47). | work owner, saat konteks Non-Life | §2b |

⚠️ **Lima OQ audit penyimpanan ditutup 2026-09-16** `[keputusan work owner]`: tanggal per peserta;
dokumen jadi tabel sendiri; keadaan tangga di `T_WORK_CLAIM`; `STS_REJECT`/`ACCEPTATION_DATE` nilai
sebenarnya; `PEGA_JSON_KLAIM_PNC` tidak dipakai dan produk dibaca relasional.

---

## Bukti struktural: lini Life menyimpang dari tiga lini Claim lain

Dua fakta korpus yang membentuk batas modul dan menguatkan **ADR-0011**. Keduanya **terverifikasi
langsung dari korpus** — tidak bersandar pada dokumen mana pun di luar `D:\XML\RNM_BRD\`.

### 1. Class `Adjustment` Life berada di **framework berbeda** `[terverifikasi]`

| Modul | Class adjustment yang dipakai |
| --- | --- |
| `Claim Life` | `ASM-FW-**GISFW**-DATA-ADJUSTMENTLIFE` — 765 kemunculan |
| `Claim Prop`, `Claim Non Prop`, `Claim Fac In` | `ASM-FW-**GCNMFW**-DATA-ADJUSTMENT` — 6.343 kemunculan |

**Tidak ada tumpang tindih sama sekali**: Claim Life memakai *hanya* yang pertama, ketiga lini lain
memakai *hanya* yang kedua. Bedanya bukan nama, melainkan **framework** (`GISFW` vs `GCNMFW`).

Perintah audit:
```
for m in "Claim Life" "Claim Prop" "Claim Non Prop" "Claim Fac In"; do
  echo -n "$m : "
  grep -rhoE "ASM-FW-G[A-Z]+-DATA-ADJUSTMENT[A-Z]*" "$m" --include="*.xml" | sort -u | tr '\n' ' '
  echo
done
```

### 2. Akseptasi Life memposting ke **tabel berbeda** `[terverifikasi]`

| Tabel | Kemunculan | Modul penulis |
| --- | ---: | --- |
| `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` | 7 | `Claim Life`, `Komite Claim Life` |
| `POOLDATA.OS_AKSEPTASI_KLAIM` | 191 | `Claim Prop`, `Claim Non Prop`, `Claim Fac In`, dan tiga modul Komite-nya |

Penulisnya pun berbeda kelas: sisi Life lewat
`ASM-FW-GISFW-INT-**LIFE_PREMIUM_DETAIL**!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` —
yaitu **kelas integrasi premium detail Life**, bukan kelas akseptasi klaim umum.

Perintah audit:
```
grep -rhoE "OS_AKSEPTASI_KLAIM[A-Z_]*" . --include="*.xml" | sort | uniq -c | sort -rn
grep -rl "OS_AKSEPTASI_KLAIM[^_]" . --include="*.xml" | cut -d/ -f2 | sort -u
```

### Akibat untuk desain

Lini Life **bukan sekadar nilai diskriminator** pada model klaim bersama — ia menyimpang secara
struktural: struct adjustment sendiri, dan sasaran penyimpanan akseptasi sendiri. Model data yang
memaksakan satu bentuk untuk keempat lini akan bertabrakan dengan kedua fakta di atas.

`[pertanyaan terbuka]` **Bagaimana penyimpangan itu diwujudkan di kode** — satu paket dengan strategi
per lini, paket terpisah per lini, atau bentuk lain — **belum diputuskan**. `CLAUDE.md` §5 hanya
menetapkan lapisan (`handlers` → `services` → `repository`), bukan pembagian paket domain di
dalamnya. Keputusan ini **milik Lead Engineer**, dan sebaiknya diambil saat scaffolding, sebelum
tiket pertama yang menyentuh model data.

## Out of Scope

- **Isi dan tangga internal Komite Life.** Konteks luar (**ADR-0001**). Spec ini berhenti di tiga
  kontrak batas.
- **Tiga modul Claim lain** (`Claim Prop`, `Claim Non Prop`, `Claim Fac In`). Pola `STS_REJECT` +
  `AdjustmentList` ada di sana juga (`[terverifikasi]` 15/126, 17/66, 41/12) tetapi **belum**
  dinyatakan work owner — OQ-039, OQ-061, OQ-062 tetap terbuka untuk ketiganya.
- **Merapikan alur.** Non-goal yang sudah disepakati (Ronde 1 Q2): urutan tahap, jalur balik, dan
  percabangan dibawa apa adanya. Empat penyimpangan sadar sudah didaftar di §Solution; selain itu,
  paritas.
- **Memperketat kelonggaran wewenang `TP`/`TR`.** Dibawa apa adanya; peninjauan saat Komite/IAM
  digarap (**ADR-0012**).
- **Mengganti penyimpanan berkas.** Google Storage dipertahankan (**ADR-0010**).
- **Memindahkan penomoran ke aplikasi.** Tetap di Oracle (**ADR-0006**).
- **Identity & Access sebagai konteks.** `[terverifikasi]` `discovery/context-map.md` menandainya
  **ABSENT** dari korpus — nol rule otorisasi. Ia dibangun dari nol, di luar spec ini.
- **Scaffolding kode.** `cmd/`, `internal/`, `pkg/`, `frontend/`, `go.mod`, `Makefile` **belum ada
  satu pun** di `OUTPUT_HASIL_RNM`; membuatnya adalah pekerjaan terpisah yang mendahului tiket mana
  pun dari spec ini (`CLAUDE.md` §5).

---

## Further Notes

**Ukuran pekerjaan.** `[terverifikasi]` Modul `Claim Life` memuat 51 Activity, 29 RDBList,
20 Section, 16 FlowAction, 8 ReportDefinition, 3 Harness, 3 When, 2 ConnectREST, 1 Flow,
1 DataTransform, 1 DecisionTable, 1 SystemSettings. Satu Flow tunggal (`Register_Flow`) memegang
empat tahap.

**Urutan yang saya sarankan untuk `/to-tickets`** — vertical slice, masing-masing menghasilkan
perilaku teramati:

1. Scaffolding + seam API pertama (satu endpoint, satu test ujung-ke-ujung) — prasyarat semuanya.
2. Register + Outstanding + penomoran (menyentuh OQ-002 → kemungkinan `needs-info`).
3. Mesin status per baris + aturan turunan "klaim selesai" — inti spesifikasi.
4. Wewenang per peran dan per `Type`.
5. Validasi DOL per `Type` + `ContentNote` dari `BusinessCode`.
6. Kontrak Komite (penyerahan + jalur balik).
7. Jejak audit.
8. Efek keluar + antre-ulang.
9. Antarmuka React per tahap.

**Yang perlu keputusan sebelum tiket dibuat:**

1. **Aturan turunan "klaim selesai"** (§Implementation Decisions 3) — keputusan spesifikasi saya,
   bukan temuan korpus. Perlu tinjauan work owner.
2. **Flag lingkungan tidak boleh menggerbangi penyimpanan** (§Implementation Decisions 11) —
   keputusan spesifikasi saya, bukan temuan korpus dan bukan keputusan work owner. Di Pega
   `IsPEGAPROD` **juga** menggerbangi simpan utama; membawanya apa adanya membuat lingkungan
   non-production tidak dapat dipakai menguji apa pun. Perlu tinjauan work owner — bila ditolak,
   konsekuensinya harus dinyatakan.
3. **Seam kedua untuk stored procedure** (§Testing Decisions) — usulan saya (a) memakai procedure
   sungguhan di skema uji. Perlu work owner + DBA.
4. **Satuan pergeseran tanggal pada validasi DOL** (§Implementation Decisions 6) — `[dugaan]` +1 hari.
5. **Pembagian paket domain di dalam `internal/`** (§Bukti struktural) — lini Life menyimpang secara
   struktural dari tiga lini Claim lain; bentuk perwujudannya di kode belum diputuskan. Pemilik:
   Lead Engineer. Diambil saat scaffolding.

**Catatan sumber.** Spec ini bersandar **hanya** pada korpus Pega `D:\XML\RNM_BRD\` dan artefak di
`OUTPUT_HASIL_RNM\`. Tidak ada klaim yang bersandar pada repo lain. Sumber ADR tunggal:
`OUTPUT_HASIL_RNM\docs\adr\ADR-0001`…`ADR-0012`.
