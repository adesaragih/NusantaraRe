# PROMPT — lanjutan 7: **MENU DAN BERANDA HANYA DARI BUKTI XML**, lalu F0.2–F0.5, lalu A3 di dalam shell

> Dibaca sesudah `PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE.md` dan lanjutan 1–6. Semua aturan di sana
> tetap berlaku kecuali yang **diralat di bab 1**. Perintah work owner 27 September 2026 yang
> melahirkan brief ini: *"JANGAN HALU DAN JANGAN MEMBUAT MENU YANG TIDAK ADA DI SISTEM SEBELUMNYA
> (PEGA). XML ITU PATOKAN DASAR UNTUK MEMBUATNYA! … KALAU DI TIKET TIDAK BENAR ATAU TIDAK SESUAI MAKA
> PERBAIKI JUGA DI TIKET NYA."* Satu giliran = **seluruh F0.2–F0.5**, lalu A3 kelompok demi kelompok
> tanpa berhenti di antaranya *(mekanisme lanjutan 3: laporan review masuk berkas, satu laporan akhir)*.

---

## 0. KEADAAN SESUDAH `0ee1c33` — F0.1 DIVERIFIKASI ULANG

| Klaim laporan F0.1 | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 15 berkas baru, fondasi tampilan berdiri | `git show --stat`: 15 berkas, +6.258 | ✅ |
| `npm run build` hijau 89 modul, `tsc` bersih | dijalankan ulang: 89 modul, CSS 29,86 kB, `tsc` exit 0 | ✅ |
| 39 uji JS | vitest: 13 desimal + 10 label + 16 api = 39, nol gagal | ✅ |
| Go tetap hijau | `go test -tags=db ./...`: **261 PASS · 0 FAIL · 34 SKIP** + 41 subtest PASS | ✅ |
| 28 ekspor `dasar.tsx` | terhitung **26** *(25 nilai + `interface Opsi`)* — selisih hitung, bukan substansi | ⚠️ kecil |
| `labels.test.ts` membaca korpus | `readFileSync` saja, `describe.skipIf(!adaKorpus)` — READ-ONLY, tidak menggagalkan CI tanpa korpus | ✅ |
| backslash `muat-env.ps1` diperiksa pada teks terender | `keadaanGalat.ts:79` memakai `\\` | ✅ |
| nol kebocoran | grep alamat/sandi/service/e-mail di `frontend/src`: nol | ✅ |

Sisa yang sudah dinyatakan executor dan **tetap dikerjakan di brief ini**: `axios` masih di
`package.json` *(dilepas di F0.2)*; `styles.css` 3.956 baris utuh *(dipangkas sesudah layar ada)*.

---

## 1. RALAT BRIEF 6 §2 ATURAN 7 — MENU HANYA YANG BERBUKTI XML

### 1.1 Bukti yang dibaca ulang hari ini *(semua path relatif `D:\XML\RNM_BRD\`, READ-ONLY)*

| # | Bukti | Isi | Arti |
| ---: | --- | --- | --- |
| 1 | `Claim Life\Struktur_InboxClaimLife.xlsx` sheet `Section` baris 2 | `Flow · Register_Flow · (root / entry point)`; **tidak ada** baris bersumber MENU/portal | Pengekspor Pega sendiri menetapkan titik masuk modul = pembuatan kasus lewat `Register_Flow` |
| 2 | `NB Treaty In\Struktur_MenuNBTreatyIn.xlsx` baris 2–3; `Endorsment Fac In\BRD_Endorsement_MenuKeFlow.xlsx` baris 2 | `Harness · SFAPortalOpportunities · menu portal (entry point)` class `PegaCRM-Portal`; `Flow … createWork dari tombol Create di SFAPortalOpportunitiesHeader`; `Sumber MENU · Harness · SFAPortalEndorsement · (entry point / portal)` | Modul yang punya menu portal **diekspor bersama harness portalnya**. Claim Life tidak → portal Claim Life **tidak ada di ekspor** |
| 3 | `Claim Life\Flow\Register_Flow.xml` | `<pyLabel>Register</pyLabel>` **155**; `<pyWorkTypeName>ClaimLife</pyWorkTypeName>` **270**; `<pyTaskName>` Input Register **358**, Outstanding Claim **343**, Medical Check **268**, Claim Analis **313**; Input Register: `WorkList` **1511** + `Current operator` **1508**; Outstanding Claim: `WorkList` **1300** + `Custom` **1319** + `ToWorklist` **1351**; Medical Check: `WorkBasket` **993** + `ToWorkbasket` **1044** + `<Workbasket>ReasLifeMedicalAdvisor` **1072**; Claim Analis: `WorkBasket` **1119** + `ToWorkbasket` **1170** + `<Workbasket>ReasLifeSPV` **1198** | Empat antrian: dua worklist milik operator, dua workbasket per peran. ⛔ Baris **147, 866, 948, 1008, 1135, 1254** memuat nama operator — **jangan dikutip** ke mana pun |
| 4 | `Section\InputRegisterClaimLife.xml` **3845–3847**, **4006–4007**; `Section\ClaimLifeDetailGCNM.xml` **5080–5081**, **5224–5225**; `Section\AdjustmentDetail_Section.xml` **15557–15558**, **15865–15867** | `<pyTarget>popup` + `<pyHarnessName>SearchPolicy_Harness` / `Diagnose_Harness` / `Committe_Life` | Ketiga harness adalah **POPUP dari dalam layar**, bukan menu. Sheet struktur menaruhnya di bawah section pemanggil *(baris 20, 67, 128)* |
| 5 | `ReportDefinition\InboxPremiumList.xml` | class `Assign-Worklist` **39**; kolom `Case ID` **721**, `Create Date/Time` **736** *(urut DESC **733**)*, `Create Operator Name` **751**, `Work Status` **765**; `pyPageSize 50` **592**; `pyMaxRecords 500` **942**; dipanggil `SearchPolicy_Section` *(popup Register)* | Satu-satunya pola **inbox ber-worklist** yang diekspor bersama modul ini. Kosakata korpus untuk antrian: **Inbox** *(`InboxPremiumList`, `Struktur_InboxClaimLife`)*, **Worklist/Workbasket** *(`ToWorklist`, `ToWorkbasket`)* |
| 6 | `Komite Claim Life\Struktur_KomiteClaimLife.xlsx` baris 2; folder `Komite Claim Life` | `KomiteLife_Flow · (root / entry point)`; **tidak ada** folder `Harness\` | Komite pun tanpa portal; menunya kelak dari root flow itu |
| 7 | `OUTPUT_HASIL_RNM\discovery\open-questions.md` tabel titik masuk | modul **tanpa** flow punya harness portal `DATA-PORTAL!…`; modul **ber**-flow titik masuknya flow | Konsisten dengan 1–2 |

Cara membaca `.xlsx` tanpa Excel dan tanpa Python *(mesin ini tidak punya Python)*:
`& .\.scratch\alat\baca-xlsx.ps1 -Path "D:\XML\RNM_BRD\Claim Life\Struktur_InboxClaimLife.xlsx"`
*(skrip baru, dibaca dulu isinya; hanya `ZipFile.OpenRead`, korpus tetap READ-ONLY)*.

### 1.2 Keputusan menu `[DIPUTUSKAN — 27 September 2026, dari perintah work owner]`

Sidebar memuat **hanya** yang berbukti:

| Kelompok | Butir | Bukti | Membuka |
| --- | --- | --- | --- |
| **Claim Life** `[nama modul korpus; pyWorkTypeName ClaimLife, Register_Flow.xml:270]` | **Inbox Claim Life** | nama berkas struktur pengekspor + kosakata `InboxPremiumList`; **tidak ada label portal di korpus** → di `labels.ts` ditandai `[tidak ada di korpus]` | halaman beranda sesudah masuk *(bab 1.3)* |
| | **Register** | VERBATIM `Register_Flow.xml:155` | layar `LAYAR.register` = `Input Register Claim Life` *(halaman `RegisterKlaim`)* |
| **Komite Claim Life** | *(belum)* | lahir hanya saat Bagian B dimulai, dari root `KomiteLife_Flow` dan section-nya | — |

Yang **TIDAK** ada di menu, dan **aturan 7 brief 6 dicabut** karenanya: 11 modul lain, `BelumTersedia`
untuk mereka, "Beranda" terpisah, serta *Dokumen / Komite / Detail & Tutup / Cari Polis / Medical
Check / Claim Analis* sebagai menu. Semua itu dibuka **dari dalam kasus**: 16 flow action + 3 popup.
Aturan 6 brief 6 *(beranda = antrian kerja pemegang peran)* **tetap**, bentuknya di bab 1.3.

Uji yang wajib ada di F0.3: `labels.test.ts` membuktikan `MENU.register` ada di
`Register_Flow.xml` baris 155; satu uji menegaskan sidebar **tidak** memuat satu pun nama modul lain
*(`Treaty`, `Endorsement`, `Fac`, `Prop`, `Master`, `Premium`)* dan tidak memuat `BelumTersedia`.

**OQ baru untuk pemilik ekspor Pega** *(tambahkan ke `.scratch/claim-life/OQ-untuk-tim.md`)*:
*"Apakah ada harness portal/navigasi untuk Claim Life (padanan `SFAPortalOpportunities` di NB Treaty
In) beserta tombol Create-nya? Bila ada, mohon diekspor."* Sampai terjawab, menu di atas adalah
**minimum berbukti**; bila jawaban datang, menu **mengikuti XML-nya**, bukan sebaliknya.

Tiket: `grep -rEi 'sidebar|beranda|homepage|menu portal' .scratch/claim-life` → **nol**. Tidak ada
tiket yang menyebut menu, jadi tidak ada ralat tiket untuk bab ini; yang diralat adalah brief 6.
Bila saat A3 sebuah layar ternyata berbeda dari tiketnya → ralat **di tiketnya** *(brief modul §1.2)*.

### 1.3 Bentuk **Inbox Claim Life** *(beranda)*

| Bagian | Ketentuan | Bukti |
| --- | --- | --- |
| Tab | **4 tab = 4 assignment**, judul VERBATIM `TAHAP` *(sudah di `labels.ts`)*; lencana cacah per tab | `Register_Flow.xml` 358 · 343 · 268 · 313 |
| Tab per peran | `ReasLifeAdmin` → *Input Register* + *Outstanding Claim* *(worklist: kasus yang `T_WORK_CLAIM.CREATE_OP` = pelaku)*; `ReasLifeMedicalAdvisor` → *Medical Check*; `ReasLifeSPV` → *Claim Analis*; pelaku multi-peran melihat gabungan; tab tanpa peran **tidak** dirender | 1508 *(Current operator)*, 1351 *(ToWorklist ke pxCreateOperator — lanjutan 4 bab 7)*, 1072, 1198 |
| Isi tab *Input Register* | kasus yang **dikembalikan ke Admin dari Outstanding** lewat `IsSendtoAdmin` di `Decision3` *(sheet baris 172)*; kasus baru hasil `POST /api/klaim-life` langsung berada di *Outstanding Claim* *(Assignment2 selesai saat register disimpan)*. ⚠️ `models.TahapDariPeran` **tidak membedakan** kedua tahap Admin *(`models/tahap.go:98–105`)*: executor **memeriksa** bagaimana tiket 08 *(`services/tahap.go`)* menyimpan kembalinya kasus ke Admin *(`PY_POSITION`, `SENDTO_ADMIN`)* dan memakainya. Bila tidak dapat dibedakan: tab tetap ada, isinya `Kosong` dengan keterangan `[terbuka — work owner]` di layar, **bukan** kolom baru yang dikarang | `Register_Flow.xml` Decision3; `When/IsSendtoAdmin.xml` |
| Kolom Pega-standar | `Case ID`, `Create Date/Time`, `Create Operator Name`, `Work Status` — label VERBATIM, ditandai `[dipinjam dari InboxPremiumList — report PremiumList Life yang diekspor bersama Claim Life]`; peta kolom: `CASE_ID`, kolom waktu pembuatan yang **ada** di skema kita *(executor menetapkan dari migrasi 001/002/012 dan mencatatnya; bila hanya `TGL_UPDATE` yang ada, itu yang dipakai dan labelnya menyebutnya)*, `CREATE_OP_NAME`, status turunan klaim | `InboxPremiumList.xml` 721 · 736 · 751 · 765 |
| Kolom Claim Life sendiri | `CLAIM_NO`, `POLICY_NO`, `BUSINESS_NAME`, `CURRENCY` dari `T_GENERAL_CLAIM`; labelnya diambil VERBATIM dari `pyLabel`/`pyCaption` section `InputOSClaimLife` atau `InputRegisterClaimLife` *(executor mencari dan mencantumkan barisnya)*; yang tidak ada → `[tidak ada di korpus]` | section terkait |
| Halaman & urutan | 50 baris per halaman, maksimum 500, urut waktu **DESC** | `InboxPremiumList.xml` 592 · 942 · 733 |
| Klik baris | membuka **layar tahap** kasus itu sesuai flow action assignment-nya: *Input Register* → `InputRegisterClaimLife`; *Outstanding Claim* → `OSClaimLife`; *Medical Check* → `MedicalCheck`; *Claim Analis* → `AkseptasiClaimLife`. Sampai kelompok A3-nya selesai, baris membuka halaman Detail yang ada *(`KlaimLife.tsx`)* dengan `BelumTersedia` **di dalam halaman** yang menyebut nama flow action-nya — bukan di menu | sheet baris 3 dan 26; `Register_Flow` |
| Tombol kepala | **`Register`** → halaman Register | `Register_Flow.xml:155` |
| Popup | `SearchPolicy_Harness` *(Register)*, `Diagnose_Harness` *(Detail)*, `Committe_Life` *(Adjustment)* dibangun sebagai `Modal` referensi — karena XML berkata `popup` | bab 1.1 butir 4 |

**Backend baru** *(F0.4)*: `GET /api/klaim-life?tahap=<1..4>&halaman=&ukuran=` — `tahap` = `models.Tahap`;
sumber `T_WORK_CLAIM ⋈ T_GENERAL_CLAIM` *(`PY_POSITION`, `CREATE_OP`)*; berbatas dan berurut; jawaban
memuat `total` untuk lencana. Gerbang **di services**, bukan di handler: tahap worklist → `CREATE_OP =
pelaku.AkunID` **dan** `WajibPeran(models.PeranAdminLife)`; tahap workbasket → `WajibPemegangTahap`.
`tahap` di luar 1–4 → 400; pelaku kosong → 403 *(bukan daftar kosong)*. Uang tidak ikut di daftar
kecuali sebagai teks *(ADR-U-0003)*; kode status hanya lewat `models` *(literal tetap satu rumah)*.

---

## 2. PAKET F0.2–F0.5 *(dikoreksi; satu commit tiap paket, `frontend: kerangka — F0.x …`)*

Sebelum **setiap** paket: baca ulang XML yang disebut di kolom bukti bab 1, dan catat barisnya di
`labels.ts` + `labels.test.ts` *(pola F0.1 yang sudah benar)*.

| Paket | Isi | Bukti selesai |
| --- | --- | --- |
| **F0.2** | `store/sesi` *(pelaku + peran, pola referensi `sessionStorage`)*; `services/api.ts` **dimigrasi utuh** ke klien `fetch` referensi *(`ApiFailure` yang sudah ada, `bolehUlang`)*, `axios` **dilepas** dari `package.json` dan `api.test.ts` disesuaikan; **Login stub**: akun teks bebas `UJI-*`, peran **boleh lebih dari satu** *(`X-Peran` dipisah koma — `handlers/pelaku.go:46`)*, hanya saat `AUTH_STUB=true`, pita "mode stub", nol sandi | uji: masuk → `X-Pelaku`/`X-Peran` terkirim; keluar → hilang; galat jaringan → `BACKEND_TIDAK_TERJANGKAU` |
| **F0.3** | **Shell** referensi *(sidebar terlipat/laci, `KelompokMenu`, `PaletMenu` Ctrl+K, topbar, menu profil, `PagarGalat`)* dengan menu **bab 1.2 saja**; `labels.ts` `MENU` + bukti; judul aplikasi "Nusantara Re" `[tidak ada di korpus]` | uji: menu per peran; palet; lipat/laci; uji "tidak ada modul lain / tidak ada BelumTersedia di menu"; uji bukti `MENU.register` |
| **F0.4** | **Inbox Claim Life** bab 1.3 + rute list + gerbang services + `Halaman`/`useHalaman` | uji murni services *(worklist milik sendiri; workbasket per peran; 403; 400)*; uji handler; uji JS per peran *(tab yang tampil, lencana, klik baris → rute)* |
| **F0.5** | `RegisterKlaim` dan `KlaimLife` **di dalam shell** *(Register lewat menu dan tombol kepala; Detail lewat klik baris)*; `App.tsx` lama dibuang; `PANDUAN-MENJALANKAN.txt` bab 4 diperbarui *(login stub, peran, Inbox, cacah uji)* | `npm run dev`: masuk → Inbox → Register → Detail, nol galat konsol; `npm run build`; `tsc`; vitest; Go hijau |

Aturan yang tetap dari brief 6 §2: 1 *(referensi disalin lalu diadaptasi; tanpa aset e-Treaty)*,
2 *(pilihan `fetch` sudah dinyatakan F0.1 — jalankan)*, 3 *(desimal tanpa float — sudah dibuktikan)*,
4 *(setiap teks membawa bukti)*, 5 *(login stub)*, 6 *(beranda = antrian, bentuk bab 1.3)*, 8
*(perilaku shell persis referensi)*. Aturan 7 **dicabut**.

---

## 3. A3 DI DALAM SHELL — urutan tetap brief 6 §3

Register → Outstanding → Detail & Tutup → Dokumen → Medis → Akseptasi → Komite. Tiap kelompok
**menggantikan** halaman Detail sementara pada tab yang bersangkutan, mengambil himpunan field, gerbang
`<pyCondition>`, dan tombol **persis** dari section-nya *(dibaca sebagai pohon; lanjutan 4 §7)*, dan
memperbarui `PARITAS-LAYAR-DAN-AKSI.md` per commit. Ketiga popup dibangun di kelompoknya masing-masing.
Perbedaan terhadap tiket → ralat **di tiket** dengan blok bertanggal dan bukti path + baris.

## 4. A4 DAN BAGIAN B

Persis lanjutan 5 §4–§5. Kelompok menu **Komite Claim Life** lahir dari `Struktur_KomiteClaimLife.xlsx`
root `KomiteLife_Flow` dan flow action-nya — bukan dari daftar yang dikarang.

## 5. LANGKAH 0

`git add PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-7.md .scratch/alat/baca-xlsx.ps1` → commit
`docs: brief lanjutan 7 — menu dan beranda hanya dari bukti XML; alat baca xlsx` → `git status
--porcelain` kosong → uji hijau *(261 · 34 SKIP · 39 JS · 89 modul)* → SHA = titik tetap **F0.2**.

## 6. LAPORAN AKHIR GILIRAN

Persis brief 6 §6, ditambah tabel **menu → bukti XML → halaman → rute** dan tabel **tab Inbox → gerbang
→ flow action tujuan**. Satu laporan, di akhir, sesudah F0.5 *(atau sesudah kelompok A3 terakhir yang
sempat selesai)*.

## 7. TELEMETRI EKSEKUSI

Per paket/kelompok, satu baris: SHA · berkas diubah · uji *(Go PASS/SKIP, JS, modul build)* · berkas
XML yang dibaca ulang *(cacah + nama)* · taksiran token *(±)* · waktu. Tutup dengan total giliran.
Angka token sejati tidak terlihat dari dalam sesi — tulis **taksiran** dan sebut begitu.

---

*Disusun 27 September 2026 sesudah verifikasi `0ee1c33` (build, tsc, vitest, `go test -tags=db`
dijalankan ulang; 15 berkas dibaca), pembacaan ulang `Struktur_InboxClaimLife.xlsx`,
`Struktur_MenuNBTreatyIn.xlsx`, `BRD_Endorsement_MenuKeFlow.xlsx`, `Struktur_KomiteClaimLife.xlsx`,
`Register_Flow.xml` (routing empat assignment), tiga harness dan enam titik `popup` pemanggilnya,
`InboxPremiumList.xml`, serta `discovery/open-questions.md` tabel titik masuk.*
