# PROMPT — EKSEKUSI SELURUH MODUL DARI HULU KE HILIR *(peta ketergantungan, gelombang, dan brief per modul)*

> Dokumen ini menggantikan "satu modul demi satu modul". Ia menetapkan **rumpun modul yang saling berhubungan**, **urutan hulu ke
> hilir** di dalam tiap rumpun, **gelombang** yang boleh berjalan paralel, dan **satu templat brief** yang dipakai tiap sesi
> executor. Aturan yang berlaku untuk seluruh modul: `PROMPT-INDUK-TIGA-MODUL.md` §0.3–§0.5 dan §2, lanjutan 8 §1, larangan
> keamanan, bab TELEMETRI EKSEKUSI. Disusun 28-09-2026 dari `spec.md` tiap modul di `.scratch/` *(baris "Hilir", "Keluar",
> "Konteks hilir")* dan cacah tiket di disk.

---

## 1. PETA KETERGANTUNGAN — empat rumpun

Bukti arah diambil dari spec, bukan tebakan: `master-product-name-life/spec.md` b74 *"dikonsumsi Claim Life, PremiumList Life,
Endorsement Life"*; `master-contract-retro-life/spec.md` b80/b223 *"dikonsumsi Claim Life, Komite Claim Life, Master Product
Name Life"*; `endorsement-life/spec.md` b9/b89 *"menulis versi baru di tujuh tabel PremiumList Life"*; `edm-treaty-in/spec.md`
b93/b112 *(polis induk dari NB Treaty In)*; `claim-prop/spec.md` b191 *(nomor polis dari master treaty)*; `claim-facin/spec.md`
b359/b406 *(penawaran fakultatif diisi NB FacIn, RNW Fac In, Endorsment Fac In)*.

```
RUMPUN LIFE        Master Contract Retro Life ─► Master Product Name Life ─► PremiumList Life ─► Endorsement Life
                                                                                   │
                                                                                   ▼
                                                                              Claim Life ─► Komite Claim Life

RUMPUN TREATY IN   NB Treaty In ─► EDM Treaty In ─► Claim Prop ─► Komite Claim Prop
                                                  └► Claim Non Prop ─► Komite Claim Non Prop        (belum ada bahan)

RUMPUN FAC IN      NB FacIn ─► RNW Fac In ─► Endorsment Fac In ─► Claim Fac In ─► Komite Claim FacIn
                   (tiga modul hulu belum ada bahan)

HULU TREATY/FAC   Treaty Contract Out ─► dibaca Claim Prop, Komite Claim Prop, Claim Fac In   (ralat 28-09: folder korpus ADA;
                                          ia penulis tunggal enam tabel master arrangement)
```

⚠️ **Ketergantungan adalah kontrak DATA, bukan kode.** Oracle tetap dipakai dan tabel warisan `POOLDATA` sudah berisi data
master. Maka modul hilir boleh dibangun **sebelum** hulunya dimigrasi selama ia **membaca tabel warisan** yang sudah ada
*(contoh nyata: Claim Life membaca view `PRODUCTINWARD_LIFE` sebelum Master Product Name Life dimigrasi, keputusan bh)*. Yang
**wajib** berurutan hanya bila hilir membaca **tabel baru** milik hulu *(contoh: Claim Life membaca `T_PREMIUM_LIST`, keputusan av)*.

## 2. KEADAAN BAHAN PER MODUL *(disk, 28-09-2026)*

| Rumpun | Modul | spec | grilling | tiket | Keadaan kode |
| --- | --- | :---: | ---: | ---: | --- |
| Life | Master Contract Retro Life | ✅ | 3 | 13 | belum |
| Life | Master Product Name Life | ✅ | 2 | 9 | belum |
| Life | PremiumList Life | ✅ | 1 | 10 | ±58% *(05a bagian 2 → 09)* |
| Life | Endorsement Life | ✅ | 2 | 13 | belum |
| Life | Claim Life | ✅ | 4 | 15 | ±92% *(§3.1 sensus)* |
| Life | Komite Claim Life | ✅ | 2 | 11 | ±10% *(01 → 09)* |
| Treaty In | NB Treaty In | ✅ | 4 | 31 | belum |
| Treaty In | EDM Treaty In | ✅ | 1 | 13 | belum |
| Treaty In | Claim Prop | ✅ | 6 | 16 | belum |
| Treaty In | Komite Claim Prop | ✅ | 8 | 15 | belum |
| Treaty In | Claim Non Prop · Komite Claim Non Prop | ⛔ | 0 | 0 | **perlu grilling** |
| Fac In | NB FacIn · RNW Fac In · Endorsment Fac In | ⛔ | 0 | 0 | **perlu grilling** |
| Fac In | Claim Fac In | ✅ | 6 | 16 | belum |
| Fac In | Komite Claim FacIn | ✅ | 4 | 14 | belum |
| Terpisah | Treaty Contract Out | ✅ | 2 | 12 | belum |

Tiket siap dieksekusi di luar tiga modul yang berjalan: **152**. Lima modul tanpa bahan tidak dieksekusi sampai grilling selesai.

## 3. GELOMBANG — yang boleh berjalan bersamaan

| Gelombang | Isi *(urutan di dalam baris = hulu ke hilir)* | Paralel dengan | Syarat mulai |
| ---: | --- | --- | --- |
| **G0** | PremiumList Life *(sisa)* → Komite Claim Life → Claim Life §3.1 | G1 | berjalan sekarang *(brief GILIRAN-10)* |
| **G1** | Master Contract Retro Life → Master Product Name Life | G0 | tabel warisan `*_LIFE` terbaca; nol tabrakan nama |
| **G2** | Endorsement Life | G3 | PremiumList Life tiket 09 menyatu *(ia menulis versi di tujuh tabel PL)* |
| **G3** | NB Treaty In → EDM Treaty In → Claim Prop → Komite Claim Prop | G2 | G0 menyatu *(pola klaim + komite dari rumpun Life dipakai ulang)* |
| **G4** | Claim Fac In → Komite Claim FacIn | Treaty Contract Out | Claim Prop menyatu *(spec Claim Fac In merujuk Claim Prop 38 kali: pola yang sama)* |
| **G5** | Treaty Contract Out | G4 | keanggotaan korpus dipastikan *(§5)*; Master Contract Retro Life menyatu |
| **Grilling** | Claim Non Prop, Komite Claim Non Prop, NB FacIn, RNW Fac In, Endorsment Fac In | kapan saja | **work owner** hadir *(grilling adalah wawancara, bukan eksekusi)* |

Paralel maksimum yang aman: **tiga jalur** sekaligus *(konflik kode bersama naik tajam di atas itu — pelajaran penjaga kaskade)*.

## 4. TATA LETAK PER MODUL — cabang, worktree, port, rentang migrasi

| Modul | Cabang / worktree `.worktrees\…` | Backend | Vite | Migrasi |
| --- | --- | ---: | ---: | --- |
| Claim Life | `main` | 8080 | 5173 | 001–029 |
| Komite Claim Life | `modul/komite-claim-life` | 8083 | 5175 | 030–049 |
| PremiumList Life | `modul/premiumlist-life` | 8082 | 5174 | 050–079 |
| Endorsement Life | `modul/endorsement-life` | 8084 | 5176 | 080–099 |
| Master Contract Retro Life | `modul/master-contract-retro-life` | 8085 | 5177 | 100–119 |
| Master Product Name Life | `modul/master-product-name-life` | 8086 | 5178 | 120–139 |
| NB Treaty In | `modul/nb-treaty-in` | 8087 | 5179 | 140–179 |
| EDM Treaty In | `modul/edm-treaty-in` | 8088 | 5180 | 180–199 |
| Claim Prop | `modul/claim-prop` | 8089 | 5181 | 200–229 |
| Komite Claim Prop | `modul/komite-claim-prop` | 8090 | 5182 | 230–249 |
| Claim Fac In | `modul/claim-facin` | 8091 | 5183 | 250–279 |
| Komite Claim FacIn | `modul/komite-claim-facin` | 8092 | 5184 | 280–299 |
| Treaty Contract Out | `modul/treaty-contract-out` | 8093 | 5185 | 300–319 |

Pembuatan worktree *(asisten atau work owner, sekali per modul)*:

```
Set-Location D:\XML\RNM_BRD\OUTPUT_HASIL_RNM
git worktree add -b modul/<modul> .worktrees\<modul> main
Copy-Item APP_RNM\.env          .worktrees\<modul>\APP_RNM\.env            # lalu HTTP_ADDR=:<backend>
Copy-Item APP_RNM\frontend\.env .worktrees\<modul>\APP_RNM\frontend\.env   # lalu DEV_PROXY_TARGET=http://localhost:<backend>
```

Pembagian berkas *(induk §2, diperluas)*: `internal/*/<prefiks>_*.go`, `handlers/rute_<modul>.go`, `frontend/src/pages/<modul>/`,
`components/<modul>/`, `assets/labels.<modul>.ts`. Shell dan `labels.ts` hanya **ditambah** butir menu kelompoknya. Kode
bersama hanya aditif dan dilaporkan. Penjaga uji bersama dipersempit **per rentang migrasi modul**, tidak dilonggarkan.

## 5. TEMPLAT BRIEF — satu sesi executor BARU per modul *(isi `<…>` dari §2–§4)*

```
# PROMPT — MODUL <Nama Modul> (sesi baru, folder .worktrees\<modul>, cabang modul/<modul>)

Baca: PROMPT-EKSEKUSI-HULU-HILIR.md (§1 rumpun, §4 tata letak), PROMPT-INDUK-TIGA-MODUL.md (§0.3–§0.5, §2),
.scratch/<modul>/spec.md, seluruh .scratch/<modul>/grilling-*.md, seluruh .scratch/<modul>/issues/*.md.
Korpus: D:\XML\RNM_BRD\<Folder Korpus>\ (READ-ONLY). REFERENSI_UI hanya disalin, tidak disunting.

0. LANGKAH 0: git status bersih; git log -1 = SHA main terbaru; uji hijau; .env port <backend>/<vite>.
1. PETA XML dulu, ditulis ke tiket 00 atau 01 sebelum kode:
   - flow utama: seluruh konektor (pyTo mendahului pyFrom di DOM), assignment, status kerja VERBATIM;
   - flow action -> section -> tombol -> activity (baris b<n> dari sed -e 's/></>\n</g');
   - harness portal bila ada -> menu. Menu HANYA yang berbukti; kelompok = nama folder korpus.
2. KONTRAK HULU (§1): baca data hulu dari tabel warisan POOLDATA yang ada, atau dari tabel baru hulu bila sudah
   menyatu. Tulis kontraknya di tiket; kunci dengan uji dua sisi.
3. TIKET BERURUTAN 00 -> terakhir, satu commit per tiket "<modul>: tiket NN — <judul>":
   pembacaan ulang XML (activity sebagai pohon; WHEN dibaca dari urutan aksi; penulis medan dicari, bukan hanya
   pembacanya); ralat bertanggal bila XML membantah tiket; uji murni + handler + JS; PARITAS milik modul;
   migrasi hanya <rentang> dan hanya dari keputusan tercatat; penjaga kata cadangan Oracle hijau.
4. LARANGAN: prosedur tidak dipanggil (logika ditiru); nol COMMIT di teks SQL; endpoint luar tidak dipanggil di DEV
   (outbox T_LOG_SERVICE_RNM + stub); -migrate tidak dijalankan; nol nama orang, nomor polis, kredensial, alamat
   layanan, nama operator Pega di artefak; fixture UJI-*.
5. BERHENTI: hanya bila tiket terakhir selesai, atau "konteks menipis: kira-kira N% tersisa" pada batas tiket.
6. LAPORAN (satu pesan): tabel tiket -> commit -> tombol/rule XML -> rute/komponen; ralat tiket; OQ dibuka/ditutup;
   kontrak hulu yang dipakai; kode bersama aditif; angka uji tiap commit; bab TELEMETRI EKSEKUSI.
7. TIDAK merge ke main. Asisten memverifikasi log lalu menyatukan, hulu lebih dulu.
```

Templat yang sama dipakai untuk **workflow ultracode**: satu agen per tiket di dalam satu jalur modul, agen verifikasi
terpisah per tiket, paling banyak tiga jalur modul bersamaan *(§3)*.

## 6. PENYATUAN — hulu lebih dulu

Asisten menyatukan per gelombang dengan urutan §1: **Master Contract Retro → Master Product Name → PremiumList → Endorsement →
Claim Life → Komite Claim Life**; **NB Treaty In → EDM → Claim Prop → Komite Claim Prop**; **Claim Fac In → Komite Claim FacIn**;
**Treaty Contract Out**. Tiap penyatuan: uji pohon terpadu di worktree bersih, lalu work owner menjalankan `-migrate` sekali
dari `main`. Konflik di kode bersama diselesaikan dengan mengambil kedua sisi dan menyempitkan penjaga per rentang migrasi.

## 7. UNTUK WORK OWNER — keputusan yang dibutuhkan sebelum gelombang berikut

| Kapan | Keputusan |
| --- | --- |
| Sebelum G1 | setuju Master Contract Retro Life dan Master Product Name Life dikerjakan **paralel** dengan G0 |
| Sebelum G3 | urutan rumpun Treaty In atau Fac In lebih dulu *(rekomendasi: Treaty In, bahan paling lengkap: 75 tiket)* |
| Sebelum G5 | Treaty Contract Out termasuk lingkup migrasi atau tidak *(tidak ada folder korpusnya di antara 17 modul)* |
| Kapan saja | jadwal **grilling** lima modul tanpa bahan |
| Tiap penyatuan | `-migrate` ke DEV; skema uji DBA *(G1)* agar 38 uji `db` berhenti SKIP |

---

*Disusun 28 September 2026 dari spec tiga belas modul di `.scratch/` (baris hilir/keluar/konsumen), hitungan rujukan silang antar
spec, cacah grilling dan tiket di disk, dan tata letak worktree yang sudah dipakai tiga modul pertama.*

## 8. RALAT 28-09-2026 — tiga folder korpus terlewat

Peta di atas disusun dari daftar folder yang terpotong. Korpus memuat **20** folder modul, bukan 17: **Treaty Contract Out**,
**Treaty In**, dan **Treaty In Adjustment** terlewat. Akibatnya:

- **Treaty Contract Out adalah HULU**, bukan terpisah: `spec.md` b106 menyatakannya satu-satunya penulis enam tabel master
  arrangement yang dibaca Claim Prop, Komite Claim Prop, dan Claim Fac In. Ia boleh dieksekusi **sekarang**, paralel dengan G0
  *(brief `PROMPT-IMPLEMENTASI-MODUL-TREATY-CONTRACT-OUT.md`, worktree `modul/treaty-contract-out`, port 8093, migrasi 300–319)*.
- **Treaty In** *(51 tiket)* dan **Treaty In Adjustment** *(13 tiket, addendum atas Treaty In)* serta **Claim Non Prop** *(37)* dan
  **Komite Claim Non Prop** *(12)* punya bahan lengkap di `dastin/_migration-docs/` *(hanya-baca)*, disusun dari korpus
  `D:\XML_NURE\` dengan konvensi nama tabel berbeda. Syarat sebelum dieksekusi: **paket rekonsiliasi** — cocokkan dengan korpus
  `D:\XML\RNM_BRD\`, tetapkan hubungan Treaty In dengan NB Treaty In dan EDM, selaraskan konvensi, salin ke `.scratch/`.
- Modul yang benar-benar tanpa bahan tinggal **tiga**: NB FacIn, RNW Fac In, Endorsment Fac In.
