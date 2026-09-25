# Nusantara Re — Migrasi Pega → Go + React + Oracle

> **Jalankan sesi agent dari folder ini (`OUTPUT_HASIL_RNM`), bukan dari root korpus
> `D:\XML\RNM_BRD\`.** Berkas ini hanya terbaca otomatis bila sesi dimulai dari sini. Menaruhnya di
> root akan melanggar batas READ-ONLY yang dijaga sejak awal proyek.

## 1. Tujuan proyek

**Migrasi**, bukan greenfield. Aplikasi reasuransi Nusantara Re yang berjalan di **Pega** dipindahkan
ke **Go (backend) + React/Vite (frontend) + Oracle (database existing)**.

Konsekuensinya mengikat: perilaku yang ada **direkam dari korpus**, bukan direka ulang dari asumsi
bagaimana reasuransi "seharusnya" bekerja. Bila korpus tidak menjelaskan sesuatu, itu **pertanyaan
terbuka**, bukan ruang untuk berimprovisasi.

## 2. Sumber kebenaran — READ-ONLY

| Lokasi | Status | Isi |
| --- | --- | --- |
| `D:\XML\RNM_BRD\` | **READ-ONLY** | Korpus ekspor rule Pega — **20 modul, 9.369 berkas `.xml`, 17 tipe rule** |
| `D:\XML\nusantara-re\` | ⛔ **DI-BLACKLIST** | **Jangan dibaca, jangan dijadikan pembanding, jangan dijadikan repo target.** Perlakukan seolah tidak ada |
| `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\` | **repo target tunggal** | Seluruh keluaran: dokumen (discovery/spec/ADR/tiket) **dan** kode Go/React |

**Tulis HANYA ke `OUTPUT_HASIL_RNM\`.** Korpus `D:\XML\RNM_BRD\` READ-ONLY — jangan membuat,
mengubah, atau menghapus apa pun di sana, termasuk berkas baru di root korpus.

⛔ **`D:\XML\nusantara-re\` di-blacklist total (keputusan work owner, 2026-09-14 — Opsi A).**
Sebelumnya ia diizinkan sebagai referensi silang READ-ONLY. Izin itu **dicabut**: jangan membukanya,
jangan mengutipnya, jangan memakainya sebagai pembanding temuan, dan jangan menjadikannya target
implementasi. Bila sebuah klaim hanya bersandar pada `nusantara-re`, klaim itu **tidak berlaku** —
turunkan menjadi pertanyaan terbuka berpemilik, atau verifikasi ulang ke korpus Pega.

## 3. Hasil FASE A Discovery — sudah selesai, pakai sebagai bukti

FASE A (D0–D4) **ditutup 2026-09-14**. Seluruh keluarannya di `discovery\`:

| Butuh | Baca |
| --- | --- |
| Bagaimana sistem bekerja (narasi per konteks & antar-konteks) | `discovery/understanding-report.md` |
| Pembagian bounded context + cakupan bukti | `discovery/context-map.md` |
| Glossary domain (169 entri, berkolom Bukti) | `discovery/glossary.md` |
| Pertanyaan terbuka | `discovery/open-questions.md` |
| Perilaku satu konteks (alur, guard, objek Oracle, batas pengetahuan) | `discovery/flows/<konteks>.md` (20) |
| Peran & ketergantungan satu modul | `discovery/modules/<modul>.md` (20) |
| Isi satu modul per tipe rule | `discovery/inventory/<modul>.md` (20) |
| Identitas rule yang isinya berbeda antar modul | `discovery/inventory/_oq011-konflik-isi.md` (533 entri) |
| Penutup tiap step | `D1-CLOSING-REPORT.md`, `D2-CLOSING-REPORT.md`, `D3-D4-CLOSING-REPORT.md` |
| Konvensi telusur yang mengikat | `discovery/flows/_METHOD.md`, `_METHOD-noflow.md` |

Angka yang sudah terverifikasi: **20 konteks ditelusur, 428 rule**; **9 bounded context** (5 `full`,
4 `partial`) **+ Identity & Access = `absent`**.

## 4. Aturan anti-halusinasi — berlaku penuh di FASE B

Aturan ini dibawa dari discovery dan **tidak berubah**.

1. **Setiap klaim perilaku wajib bukti `path + rule`.** Sebut berkas yang dibaca, bukan hanya nama
   rule (korpus dirakit dari >1 server Pega dan memuat hostname DEV).
2. **Identitas rule = `class / nama` dari `<pxInsName>`.** Tipe rule dibaca dari field pertama
   `<pzOriginalInstanceKey>` (atau `<pxObjClass>` sebagai fallback).
3. **`<pzOriginalInstanceKey>` adalah asal salinan, BUKAN identitas.** Field 2–3 menyimpan
   provenance Save-As. Terukur di D1: namanya berbeda pada 410 dari 1.149 berkas (35,7 %),
   class pada 62. Memakainya sebagai identitas akan salah.
4. **Grep dulu, baca rentang seperlunya.** Ada Harness 8,2 MB dan activity 996 KB. Jangan membaca
   berkas besar utuh. Catat bila sebuah berkas **belum habis dibaca**.
5. **Jangan menebak** skema/tipe kolom, arti kode/status/enumerasi, kepanjangan singkatan, rumus,
   isi stored procedure, atau struktur JSON. Yang tidak dijelaskan korpus ditulis
   **`belum terverifikasi`**.
6. **Jangan menyatukan identitas yang berkonflik.** Cek `_oq011-konflik-isi.md` lebih dulu; telusur
   **per-varian**; jangan meminjam varian dari modul lain.
7. **Perbandingan isi wajib memakai hash ternormalisasi**, bukan `md5sum` mentah — ekspor Pega
   mengacak urutan elemen dan menyisipkan timestamp. Daftar tag volatil ada di `_METHOD.md`.
8. **Nomor OQ diambil dari register `discovery/open-questions.md`**, bukan dari ingatan.
9. **Label wajib** di setiap pernyataan: `[terverifikasi]` / `[dugaan]` / `[pertanyaan terbuka]`.
   **Setiap angka disertai perintah audit** yang menghasilkannya.
10. **Nilai nama orang tidak disalin** ke artefak mana pun. Catat rule, tag, dan jumlahnya saja.

## 4a. Aturan sensus

**Alasan bagian ini ada:** jendela sensus yang dikarang ulang per pertanyaan akan berbeda-beda, dan
perbedaannya **baru ketahuan setelah ada yang menghitung ulang**.

> ⛔ **KEWAJIBAN MEMAKAI `sensus.py` DICABUT** — `[keputusan work owner]` **2026-09-19**.
>
> Judul lama bagian ini dikutip utuh, tidak dihapus: *"## 4a. **Sensus wajib lewat `sensus.py`**"*,
> beserta kalimat kewajibannya: *"Setiap sensus atas `spec.md`, `issues/`, atau korpus XML
> dijalankan dengan **`OUTPUT_HASIL_RNM\sensus.py`**, **bukan** dengan `grep` atau regex yang
> ditulis ulang tiap kali."*
>
> **Sebabnya — dan ia jujur:** `sensus.py` **pernah ada**. Arsip `PROMPT-AI1-AI4-PENUTUP-CLAIM-PROP`
> mencatatnya *"sudah ada di OUTPUT_HASIL_RNM"*, dan `PROMPT-AJ1-AJ3-RALAT-EJAAN-CLAIM-PROP`
> mencatat ia sempat **diperbarui sampai versi ketiga**. `[terverifikasi]` **Kini berkasnya hilang
> tanpa satu pun salinan** — dicari di seluruh pohon dan nihil — dan proyek ini **bukan repositori
> git**, sehingga **tidak dapat dipulihkan**. ⛔ Aturan yang menunjuk berkas yang tak dapat
> dipulihkan **tidak dapat ditegakkan siapa pun**.
>
> ⚠️ **Dan sebut harganya, supaya pencabutan ini tidak dibaca sebagai "tidak penting":** dalam dua
> ronde beruntun sebuah kesimpulan runtuh karena sensus dihitung dengan jendela yang salah —
> **ronde 4 Claim Fac In** *(bentuk alur: empat jenis entri dijumlahkan sebagai satu)* dan
> **ronde 5** *(tipe parameter: nama dipasangkan ke tipe menyilang tujuh wadah)*. ⭐ Yang kedua
> sempat **menggeser dasar vonis ADR-0003** sebelum diralat.

### Aturan yang menggantikannya

⭐ **Tiap sensus WAJIB dihitung DUA CARA yang berbeda, dan jendelanya disebut.** Bila kedua cara
berselisih, tulis **"belum punya data"** dan **cantumkan keduanya** — ⛔ **jangan memilih salah
satu**. ⚠️ Bila sebuah angka sensus lama diralat, **sebut cara mana yang keliru dan kenapa** —
bukan hanya angka barunya.

⚠️ **"Dua cara" berarti dua jalan yang benar-benar berbeda**, bukan satu skrip dijalankan dua kali.
Yang sudah terbukti berguna: mengurai vs. mencacah baris · sensus penuh vs. **rekonsiliasi delta
terhadap angka lama yang tertulis** · menghitung maju vs. menguji kelengkapan dua arah *(A\B dan
B\A)*.

### Empat jebakan sensus yang sudah terbukti menggigit

1. unit AC hanya baris pertama — frasa di baris lanjutan terlewat;
2. blok kutipan `>` RALAT ikut terbaca — teks lama yang sudah dibatalkan terhitung sebagai aktif;
3. frasa terlipat di pergantian baris — `penyimpangan` di satu baris, `sadar` di baris berikutnya;
4. nama rule/tabel dicocokkan **peka** huruf besar-kecil — satu objek bisa dieja beberapa cara.

⭐ **Ditambah dua yang lahir sesudahnya:**

5. **entri dijumlahkan lintas wadah** — sebuah tag yang sama tinggal di beberapa wadah berbeda, dan
   menjumlahkan semuanya memberi angka yang **bukan jawaban pertanyaannya**. Sebut wadahnya, bukan
   hanya tagnya;
6. **dua daftar dipasangkan menurut urutan** — memasangkan nama ke-*n* dengan tipe ke-*n* hanya sah
   bila keduanya datang dari wadah yang sama. Pasangkan lewat **kunci**, bukan lewat urutan.

**Satu lagi yang diperingatkan tetapi tidak bisa dipaksakan:** properti korpus wajib dicari dengan
**awalan halaman** — `InputData.CARI16`, **bukan** `CARI16` telanjang. Yang telanjang memungut
halaman lain dan memberi angka palsu.

⚠️ **Dua hal operasional di mesin ini** — tetap berlaku untuk skrip apa pun yang ditulis:

- Pemanggilnya **`py`**, bukan `python3` atau `python` — keduanya hanya shim Microsoft Store yang
  gagal. `py --version` → Python 3.14.7.
- Awali dengan **`PYTHONIOENCODING=utf-8`**. Konsol Windows default cp1252 dan skrip akan melempar
  `UnicodeEncodeError` saat mencetak ⚠️.

⚠️ **Kalimat yang tetap berlaku walau alatnya gugur:** *"Verifikator yang instrumennya sendiri
belum diverifikasi memproduksi tuduhan palsu."* ⛔ Karena itu **uji instrumennya dulu** — jalankan
atas beberapa butir yang jawabannya **sudah diketahui**, dan laporkan hasilnya bersama sensusnya.

## 5. Struktur folder aplikasi target

**Lokasinya: `OUTPUT_HASIL_RNM\APP_RNM\`** *(sejak 25 September 2026 sore; sebelumnya langsung di
`OUTPUT_HASIL_RNM\`)*. Seluruh kode Go dan React ada di dalam `APP_RNM\` — `cmd/`, `internal/`,
`pkg/`, `frontend/`, `go.mod`, `Makefile` — sedangkan `discovery/`, `docs/`, `CONTEXT.md`,
`.scratch/`, `dastin/`, `jefri/` tetap di `OUTPUT_HASIL_RNM\`. **Satu repo git, berakar di
`OUTPUT_HASIL_RNM\`.** Frontend memakai **TypeScript**: komponen `.tsx`, modul lain `.ts`, nol
`.jsx`. Kode ditulis untuk pembaca yang **baru mengenal Go dan React** — mulai dari
`APP_RNM\README-BACA-DULU.md`.

```
APP_RNM/
├── README-BACA-DULU.md
├── cmd/api/main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   └── services/
├── pkg/utils/
├── frontend/               ← React via Vite, TypeScript (.tsx)
│   ├── src/assets/
│   ├── src/components/
│   ├── src/hooks/
│   ├── src/pages/
│   ├── src/services/
│   ├── src/store/
│   ├── src/App.tsx · src/main.tsx · src/vite-env.d.ts
│   └── package.json · tsconfig.json · vite.config.ts
├── go.mod · go.sum
└── Makefile
```

**Arah dependency: `handlers` → `services` → `repository`.** Tidak boleh terbalik, tidak boleh
memotong lapisan.

## 6. Pertanyaan terbuka — 61 terbuka, 38 memblokir FASE B

`discovery/open-questions.md`: **61 terbuka / 2 terjawab** (OQ-001…OQ-063, per 2026-09-14).
**38 di antaranya memblokir FASE B** (daftar per pemilik peran di `D3-D4-CLOSING-REPORT.md` §3,
dengan penanda ⛔) — angka itu dihitung saat D3/D4 ditutup, sebelum OQ-060…OQ-063 dibuat.

Empat OQ terbaru berasal dari grilling Claim — Life: **OQ-060** (cakupan `CURRENCY`),
**OQ-061** (tiga tingkat `STS_REJECT`, terjawab untuk Claim — Life), **OQ-062** (kefinalan,
terjawab untuk Claim — Life), **OQ-063** (gerbang `TP`/`TR`, terjawab).

Pemilik utama: **Product+Underwriting 28, DBA 11, IAM 10, Finance 4, pemilik export Pega 3,
Arsitektur Pega 1**.

**FASE B hanya boleh menutup OQ yang dijawab oleh pemilik pekerjaan (work owner).** Sisanya
**tetap terbuka** — jawaban harus datang dari **DBA / Product+Underwriting / Actuarial / Finance /
IAM**, bukan dari agent dan bukan dari pembacaan ulang korpus.

Empat OQ dengan dampak terluas: **OQ-020** (arti kode/enumerasi) dan **OQ-002** (body 67 stored
procedure) menyentuh **9 dari 9 konteks**; **OQ-021** dkk (RBAC) 8 konteks; **OQ-011** (versi rule
production) 7 konteks.

Batas pengetahuan yang tidak bisa ditembus dengan membaca lebih banyak: 67 SP tanpa body,
8 objek db-link, 10 objek JSON, **49 `DecisionTable` tanpa baris keputusan**, 5 keluarga `When`
tak terbaca, 533 identitas berkonflik, nol DDL, nol rule otorisasi.

## 7. Uang material — jangan `float`

Nilai uang **tidak boleh** direpresentasikan sebagai `float`. Representasi yang dipakai
(mis. `decimal`, integer minor unit, atau `NUMERIC` Oracle) adalah **keputusan ADR**, bukan pilihan
ad-hoc per handler.

`[terverifikasi]` Korpus memuat nilai uang ter-hardcode **tanpa penyebutan mata uang**: pangsa
`0.45` / `0.05`, ambang `"3000000000"` (**dibandingkan sebagai string**), limit `30000000.00` /
`50000000.00` / `57750000.00`. Mata uangnya adalah **OQ terbuka** (OQ-037, OQ-040, OQ-046) dan
harus dijawab sebelum angka mana pun dipindahkan.

## 8. Skill Matt Pocock — hanya dipicu manusia

`[terverifikasi]` Plugin `mattpocock-skills@claude-plugins-official` v1.2.3 terpasang (scope user).
Skill berikut ber-frontmatter **`disable-model-invocation: true`** sehingga **agent tidak dapat
memanggilnya**: `wayfinder`, `grill-with-docs`, `to-spec`, `to-tickets`,
`setup-matt-pocock-skills`, `implement`, `triage`, `handoff` — **20 dari 35** skill, terverifikasi
dari frontmatter 25 September 2026.

⚠️ **Ralat 25 September 2026** *(temuan sesi Fase 0, `LAPORAN-FASE-0.md` §4)*: `tdd` dan
`code-review` **tidak** ber-flag itu — keduanya boleh dipanggil agent, dan memang dipanggil oleh
`implement`. Kalimat "berhenti" di bawah berlaku untuk yang ber-flag saja.

**Bila diminta menjalankan skill Matt Pocock dan tidak bisa: BERHENTI.** Jangan diam-diam
mereplikasi alur internal skill. Laporkan:

1. bahwa skill tidak dapat dipanggil agent (sertakan pesan penolakannya sebagai bukti);
2. command persis yang harus diketik manusia sendiri, mis. `/mattpocock-skills:grill-with-docs`;
3. tunggu keputusan — jangan berasumsi.

## 9. Alur FASE B

```
/mattpocock-skills:grill-with-docs   →   :to-spec   →   :to-tickets
```

| Keluaran | Lokasi |
| --- | --- |
| Spec | `.scratch/<konteks-slug>/spec.md` |
| Tiket | `.scratch/<konteks-slug>/issues/NN-<slug>.md` |
| ADR | `docs/adr/` |
| `CONTEXT.md` | dibuat **lazy** oleh `/domain-modeling`; seed = `discovery/glossary.md` |

**Urutan konteks yang diusulkan** (berdasarkan jumlah OQ pemblokir + cakupan bukti,
`D3-D4-CLOSING-REPORT.md` §4):

1. **Claim — Life** — 5 pemblokir, cakupan `full`, 136 berkas, paling terpisah (overlap 12–14 %)
2. **Life — Penawaran & Premium List** — 9 pemblokir, `full`, berbagi mesin dengan #1
3. **Treaty Arrangement** — 5 pemblokir, batas konteks sudah jelas (OQ-022 terjawab)

**Menunggu jawaban bisnis:** Facultative Inward (20 pemblokir, 6.071 berkas) dan Treaty Inward —
Master & Akseptasi (18). **Komite** cakupannya `full` tetapi hampir seluruh pemblokirnya RBAC —
baru produktif setelah model peran diputuskan.

## 10. Persetujuan manusia diperlukan sebelum

Operasi destruktif; migration DB; perubahan authn/authz; akses credential atau regulated data;
perubahan sistem eksternal.

**Jangan** memasukkan secret, token, data pelanggan, atau dump production ke prompt, tiket, spec,
kode, maupun test. Endpoint & host internal = **konfigurasi/env var**, bukan literal — daftar
endpoint sesungguhnya ada di tabel Oracle `M_LINK_SERVICE` yang **isinya tidak ada di korpus**
(OQ-047).

---

## Agent skills

### Issue tracker

Local markdown di `OUTPUT_HASIL_RNM/.scratch/<konteks-slug>/` — bukan repo git, tanpa tracker
eksternal. Sembilan slug konteks sejajar dengan `discovery/context-map.md` §1.
See `docs/agents/issue-tracker.md`.

### Triage labels

Default: `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`; dicatat
sebagai baris `Status:` di tiap berkas issue (tracker ini tanpa label API).
See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: `OUTPUT_HASIL_RNM/CONTEXT.md` (belum ada; seed = `discovery/glossary.md`)
+ `OUTPUT_HASIL_RNM/docs/adr/`. See `docs/agents/domain.md`.
