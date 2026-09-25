# Prompt Awal untuk Claude Code — Migrasi Pega → Go + React (Nusantara Re)

Dokumen ini berisi **prompt siap copy-paste ke Claude Code**, disusun mengikuti metodologi
Agentic Development / rangkaian skill Matt Pocock (runbook tim di
`../1. Agentic Development Methodology/`).

Alur dibagi **dua besar**, sesuai permintaan:

- **FASE A — DISCOVERY / PEMAHAMAN (STEP D0–D4):** Claude **mempelajari dulu** aplikasi & seluruh
  flow XML Pega secara lengkap dan menuliskan pemahamannya sebagai catatan discovery **berbukti**.
  **Belum ada** keputusan desain, ADR, spec, atau tiket di fase ini.
- **FASE B — PENULISAN (STEP W0–W5):** setelah pemahaman matang dan diverifikasi manusia, baru
  masuk grilling → CONTEXT/ADR → BRD → spec → ticketing → konsolidasi steering.

> **Aturan urutan:** JANGAN mulai FASE B sebelum FASE A selesai dan di-review manusia.
> Satu STEP idealnya satu sesi Claude Code. Jaga FASE A dalam rangkaian konteks yang sama;
> pindah sesi baru saat konteks membengkak (pakai `/handoff`).

---

## Fakta yang sudah diverifikasi (bukan tebakan)

- **Sumber XML Pega (READ-ONLY):** `D:\XML\RNM_BRD\` — **20 modul BRD**:
  `Claim Fac In`, `Claim Life`, `Claim Non Prop`, `Claim Prop`, `EDM Treaty In`,
  `Endorsement Life`, `Endorsment Fac In`, `Komite Claim FacIn`, `Komite Claim Life`,
  `Komite Claim Non Prop`, `Komite Claim Prop`, `Master Contract Retro Life`,
  `Master Product Name Life`, `NB FacIn`, `NB Treaty In`, `PremiumList Life`, `RNW Fac In`,
  `Treaty Contract Out`, `Treaty In`, `Treaty In Adjustment`.
- **Panduan metodologi:** `D:\XML\RNM_BRD\1. Agentic Development Methodology\` + `D:\XML\RNM_BRD\.runbook.txt`.
- **Folder output SEMUA artefak:** `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\`.
- **Tipe rule Pega dalam korpus: 17 folder tipe** (dikoreksi di STEP D0, lihat
  `discovery/README.md` §6.1) — `Activity`, `ConnectREST`, `DataPage`, `DataTransform`,
  `DecisionTable`, `DecisionTree`, `Flow`, `FlowAction`, `HTMLRule`, `Harness`, `Menu`, `RDBList`,
  `ReportDefinition`, `Section`, `SystemSettings`, `When`, `excludeXML`.
  **`ConnectREST` (51 file) adalah tipe integrasi eksternal** — jangan asumsikan integrasi hanya
  terlihat dari dalam Activity.
- **Tag XML kunci (terverifikasi dari file nyata + dikoreksi di STEP D0):**
  - **`<pxInsName>` → sumber otoritatif identitas rule (class + nama).** Format `CLASS!NAMA`;
    untuk RDBList `CLASS!PREFIX!NAMA` (prefix `ASM`/`RNM`, lihat OQ-008). Cocok 100% dengan nama
    file di batch 1. Ini yang dipakai sebagai identitas.
  - **`<pzOriginalInstanceKey>` BUKAN identitas** — ini kunci rule yang di-"Save As" (asal
    salinan). Field pertamanya (`RULE-OBJ-<TIPE>`) valid sebagai **tipe rule**, tapi nama/class-nya
    berbeda dari identitas pada 410/1.149 file batch 1 (35,7%). Rekam sebagai **"asal salinan"**
    untuk melacak silsilah klon, jangan sebagai identitas.
  - `<pyStartActivity>` (Flow: activity awal), `<pyActivityType>`, `<pyActivityName>`.
  - Step/logika Activity: `<pyStepsActivityName>`, `<pyStepsPreCondParamsWhen>`,
    `<pyStepsJavaSource>` (Java tertanam).
  - RDBList: `<pyBrowseSQL>` (SQL asli). Section: `<pyInclude>` (embed = dependency nyata),
    `<pyType>` = `FIELD`/`LAYOUT`/`SUB_SECTION`. `<pySection>` di aksi Refresh **bukan** dependency.

## Target arsitektur (mengikat — dipakai FASE B)

- **Backend:** Go — **Frontend:** React (Vite) — **Database:** Oracle existing (tetap dipakai).
- **Struktur folder web app (WAJIB persis):**

```
my-web-app/
├── cmd/
│   └── api/
│       └── main.go            # Entry point backend Golang
├── internal/                  # Kode privat backend
│   ├── config/                # Env vars + setup database
│   ├── handlers/              # HTTP controllers (REST API)
│   ├── models/                # Struct data & entitas domain
│   ├── repository/            # Query & interaksi database
│   └── services/              # Business logic (jembatan handler-repo)
├── pkg/
│   └── utils/                 # Helper publik (hashing, formatter, dll)
├── frontend/                  # Root React (via Vite)
│   ├── public/
│   ├── src/
│   │   ├── assets/            # Gambar, font, stylesheet global
│   │   ├── components/        # Reusable UI (Button, Modal, dll)
│   │   ├── hooks/             # Custom hooks (useAuth, useFetch, dll)
│   │   ├── pages/             # Komponen per route halaman
│   │   ├── services/          # API client (Axios/Fetch) ke backend Go
│   │   ├── store/             # State management (Redux/Zustand/Context)
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── vite.config.js
├── go.mod
├── go.sum
└── Makefile                   # Otomasi run/build kedua environment
```

- Arah dependency backend: `handlers → services → repository`.

## Aturan anti-halusinasi (berlaku di SEMUA step)

1. **Ini migrasi, bukan greenfield.** Kebenaran ada di Pega existing. Setiap klaim perilaku
   **wajib** sebut bukti `path file + nama rule` (contoh: `NB Treaty In/Flow/InputRealizationTreatyIn.xml`).
2. **File XML besar. Grep dulu, baca range belakangan.** Jangan baca utuh file besar.
3. **Jangan simpulkan perilaku dari nama rule saja** — nama Pega sering menipu.
4. **Identitas rule = `class / nama / tipe`** — `class` + `nama` dari **`<pxInsName>`**, `tipe`
   dari field pertama `<pzOriginalInstanceKey>`. **JANGAN** memakai `<pzOriginalInstanceKey>`
   sebagai identitas (itu asal salinan). Nama rule tidak unik lintas tipe/class.
5. **Jangan mengarang skema/tipe data** (tidak ada DDL/`Property/` di export). Yang belum terlihat
   = **pertanyaan terbuka berpemilik**, bukan tebakan.
6. **Jangan hardcode identitas orang/environment** → jadi RBAC/config di FASE B.
7. **Jangan mengarang kepanjangan singkatan** atau arti field dari caption sebelahnya.
8. **Korpus `D:\XML\RNM_BRD\` READ-ONLY.** Menulis HANYA ke `OUTPUT_HASIL_RNM\`.
9. **Uang material** — jangan `float`; keputusan lewat ADR (FASE B).
10. **Hasil agent belum tepercaya sebelum diverifikasi.** Tandai asumsi & scope creep.

---

# FASE A — DISCOVERY / PEMAHAMAN APLIKASI (belajar dulu, belum menulis keputusan)

Tujuan FASE A: menghasilkan **peta pemahaman lengkap** aplikasi dan flow XML-nya, cukup detail
sehingga engineer/agent lain bisa paham cara kerja sistem **tanpa** membuka Pega — semuanya
berbukti path+rule. Output FASE A ditulis ke `OUTPUT_HASIL_RNM\discovery\`.

## STEP D0 — Setup ruang kerja discovery

**Sesi baru. Jalankan dari `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\`.**

```
Konteks: proyek migrasi sistem reasuransi Nusantara Re dari Pega ke Go + React + Oracle.
FASE saat ini = DISCOVERY (belajar dulu). Belum membuat keputusan desain/ADR/spec/tiket.

Semua tulisan HANYA ke D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\. Sumber kebenaran = korpus XML Pega
READ-ONLY di D:\XML\RNM_BRD\ (20 modul). JANGAN menulis/mengubah/menghapus apa pun di luar
OUTPUT_HASIL_RNM\.

Buat kerangka discovery:
- discovery/README.md            -> tujuan fase, aturan bukti (path+rule), status per modul
- discovery/inventory/           -> inventaris rule per modul (per tipe)
- discovery/flows/               -> hasil telusur flow per konteks
- discovery/modules/             -> catatan pemahaman per modul
- discovery/glossary-draft.md    -> istilah domain yang DITEMUKAN (draft, belum final)
- discovery/open-questions.md    -> hal yang tidak bisa dipastikan dari korpus (berpemilik: peran)

Tulis di discovery/README.md aturan main FASE A:
- setiap klaim wajib bukti path file + nama rule;
- identitas rule ditulis class/nama/tipe dari <pzOriginalInstanceKey>;
- grep dulu, baca range belakangan; jangan baca utuh file besar;
- jangan menebak schema/singkatan/arti caption; yang tak pasti masuk open-questions.md.

Laporkan struktur yang dibuat.
```

## STEP D1 — Inventaris rule (peta "apa saja yang ada")

```
Tujuan: inventarisasi SELURUH rule di 20 modul D:\XML\RNM_BRD\ TANPA menyimpulkan logika bisnis.
Batas: hanya struktur + metadata. Grep terarah; jangan baca isi penuh tiap file.

Untuk tiap modul, hasilkan discovery/inventory/<modul>.md yang mendaftar tiap rule sebagai:
- nama file, tipe rule (folder), dan class/nama/tipe dari <pzOriginalInstanceKey>;
- ringkasan 1 baris peran rule (dari tipe + nama file + class), diberi label
  [terverifikasi] bila didukung tag, atau [dugaan] bila hanya dari nama (tandai jelas).

Lalu discovery/inventory/_summary.md:
- jumlah file per tipe rule per modul (angka terukur + cara menghitungnya);
- daftar nama file yang dipakai >1 tipe rule (tabrakan lintas tipe), dengan class-nya;
- daftar When rule bernama sama di >1 class Pega.

Aturan: angka wajib bisa diaudit ulang (sertakan perintah grep/pencarian). Yang tak terukur
tulis "belum terukur". JANGAN menyimpulkan perilaku di step ini.
```

## STEP D2 — Telusur flow end-to-end (peta "bagaimana aplikasi berjalan")

Ulangi **per konteks**. Untuk domain treaty inward, titik masuk Flow yang benar-benar ada adalah
`NB Treaty In/Flow/InputRealizationTreatyIn.xml` — **`Treaty In` sendiri TANPA Flow** (lihat
`discovery/README.md` §6.2; 6 modul tanpa Flow). Untuk modul tanpa Flow, tetapkan titik masuk
alternatif lebih dulu (kandidat: `Harness` lalu `FlowAction`).

```
Tujuan: memahami dan mendokumentasikan flow END-TO-END konteks <NAMA> dari file Flow di
D:\XML\RNM_BRD\<modul>\Flow\ (bila ada), ditelusuri sampai activity, when, section, dan SQL yang
dipanggilnya. Tulis ke discovery/flows/<konteks>.md.

Bila modul TIDAK punya Flow: mulai dari Harness (titik masuk layar) dan/atau FlowAction, catat
metode titik masuk yang dipakai di awal file. Jangan memaksakan metode <pyStartActivity>.

Cara telusur (berbukti, bukan tebakan):
1. Bila ada Flow: baca <pyStartActivity> dan shape/assignment/connector di dalamnya.
   Catat urutan langkah, percabangan, dan workbasket/assignment (routing) yang terlihat.
2. Untuk tiap Activity yang dipanggil: baca <pyStepsActivityName> (method tiap step),
   <pyStepsPreCondParamsWhen> (precondition), dan <pyStepsJavaSource> bila ada Java tertanam.
   Rangkum APA yang dilakukan step, bukan menebak mengapa.
3. Untuk tiap When yang jadi guard: catat kondisinya apa adanya (nilai literal yang diuji).
4. Untuk RDBList yang dipanggil: kutип SQL dari <pyBrowseSQL> (ringkas), catat tabel & kolom
   yang di-SELECT/INSERT/UPDATE. JANGAN menyimpulkan tipe kolom.
5. Untuk UI: dari Harness & Section, catat <pyInclude> (embed = dependency nyata) dan struktur
   <pyType> FIELD/LAYOUT/SUB_SECTION. Ingat <pySection> di aksi Refresh BUKAN dependency.

Isi discovery/flows/<konteks>.md dengan:
- Diagram alur tekstual (langkah -> langkah, dengan guard When);
- State/status yang berubah + nilai literalnya (apa adanya, tanpa mengarang artinya);
- Tabel/kolom Oracle yang disentuh + operasinya (read/write) beserta rule sumbernya;
- Integrasi eksternal yang terlihat (REST/service call) sebagai konfigurasi, bukan endpoint literal;
- Daftar rule terlibat sebagai class/nama/tipe + path;
- Hal yang tidak bisa dipastikan (arti kode status, formula, DDL) -> catat di open-questions.md.

Setiap poin wajib menyebut path file + nama rule. Bila arti sebuah nilai/kode tidak tertulis di
korpus, tulis "arti belum terverifikasi" — jangan ditebak.
```

## STEP D3 — Catatan pemahaman per modul + istilah domain (draft)

```
Tujuan: konsolidasi pemahaman per modul dan istilah domain yang DITEMUKAN selama D1-D2.

1. discovery/modules/<modul>.md: ringkasan peran modul dalam sistem, fitur/proses utamanya,
   entitas data yang disentuh, ketergantungan ke modul lain (dengan bukti), dan batasan yang
   terlihat. Tandai [terverifikasi]/[dugaan]/[pertanyaan terbuka].
2. discovery/glossary-draft.md: istilah domain (Indonesia & Inggris apa adanya dari korpus)
   dengan definisi ringkas + kolom Bukti (path+rule). Aturan bahasa: jangan mengarang
   terjemahan Inggris untuk istilah Indonesia; jangan mengarang kepanjangan singkatan; jangan
   ambil arti field dari caption sebelahnya. Singkatan yang tak dijabarkan korpus -> tulis
   "kepanjangan belum terverifikasi".
3. Perbarui discovery/open-questions.md: kelompokkan per pemilik peran (DBA / Product+Underwriting
   / Actuarial / Finance / IAM) dan sebut apa yang diblokir tiap pertanyaan.
```

## STEP D4 — Peta konteks & laporan pemahaman menyeluruh (checkpoint manusia)

```
Tujuan: sintesis SELURUH discovery menjadi pemahaman menyeluruh tingkat sistem.

Hasilkan:
- discovery/context-map.md: bounded context yang teridentifikasi dari 20 modul, modul mana milik
  konteks mana, kebocoran batas antar-konteks yang terlihat (dengan bukti), dan cakupan bukti
  tiap konteks: full/partial/absent. Konteks yang tak bisa diturunkan dari korpus (mis. Identity
  & Access bila memang tak ada rule-nya) -> tandai absent + jadikan pertanyaan terbuka.
- discovery/understanding-report.md: ringkasan naratif "bagaimana aplikasi ini bekerja" per
  konteks, alur utama antar-modul, dan daftar risiko migrasi yang sudah terlihat.

Di akhir, laporkan: apa yang sudah dipahami dengan yakin (berbukti), apa yang masih dugaan, dan
daftar pertanyaan terbuka berpemilik. JANGAN menutup gap dengan tebakan. Berhenti di sini dan
minta review manusia sebelum masuk FASE B (penulisan).
```

**GATE (wajib):** manusia me-review `understanding-report.md`, `context-map.md`, dan
`open-questions.md`. FASE B tidak dimulai sebelum pemahaman ini dianggap cukup.

---

# FASE B — PENULISAN (grilling → CONTEXT/ADR → BRD → spec → ticketing → steering)

Fase ini memakai hasil FASE A sebagai bukti. Jaga grilling→spec→tiket dalam rangkaian konteks
yang sama; pisahkan per konteks.

## STEP W0 — Setup steering repository (Tahap 0 penulisan)

```
/setup-matt-pocock-skills

Konteks: FASE penulisan migrasi Nusantara Re. Semua tulisan HANYA ke
D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\. Gunakan hasil FASE A di OUTPUT_HASIL_RNM\discovery\ sebagai
bukti. Korpus D:\XML\RNM_BRD\ tetap READ-ONLY.

Siapkan struktur agentic:
- CLAUDE.md          -> instruksi operasional repo (baca dulu). Isi: tujuan (migrasi, bukan
  greenfield), lokasi korpus READ-ONLY, ringkasan hasil discovery, aturan anti-halusinasi,
  struktur folder web app target (salin persis dari dokumen ini), arah dependency
  handlers->services->repository, alur skill.
- CONTEXT.md         -> glossary domain final; SEED dari discovery/glossary-draft.md, tiap
  istilah wajib kolom Bukti (path+rule).
- docs/adr/          -> ADR + ADR-TEMPLATE.md + README.md
- docs/agents/issue-tracker.md, triage-labels.md, domain.md

Laporkan struktur + status "repo siap workflow agentic".
```

## STEP W1 — Grilling per konteks

```
/grill-with-docs

Tujuan: mematangkan requirement migrasi konteks <NAMA>. Bukti: discovery/flows/<konteks>.md,
discovery/modules/<modul>.md, CONTEXT.md, discovery/context-map.md, korpus D:\XML\RNM_BRD\<modul>.

Uji ambiguitas & asumsi tersembunyi. Untuk tiap perilaku: siapa penggunanya (peran), perilaku
diharapkan & cara verifikasinya, non-goal, cara sistem gagal, state/transisi + guard, serta
security/authorization(SoD)/audit/data/migration/rollback/observability.

Wajib: klaim menyebut bukti path+rule; rule ditulis class/nama/tipe; istilah baru -> CONTEXT.md
(dengan bukti); keputusan sulit dibalik -> ADR (pakai ADR-TEMPLATE.md); hal yang butuh DDL/body
SP/bentuk JSON/arti kode/formula yang TIDAK ada di korpus -> pertanyaan terbuka berpemilik,
JANGAN ditebak.
```

## STEP W2 — BRD + Spec

```
/to-spec

Sintesis grilling <konteks> menjadi:
1. docs/brd/BRD-<konteks>.md: latar belakang (migrasi dari Pega), ruang lingkup & non-goal,
   aktor/peran, proses bisnis & state transition (bukti path+rule), aturan bisnis, kebutuhan
   data (rujuk pertanyaan terbuka bila DDL belum ada), integrasi eksternal, non-fungsional
   (security, audit, uang non-float), asumsi & pertanyaan terbuka berpemilik, dan traceability
   ke rule Pega sumber + ke catatan discovery.
2. docs/spec/spec-<konteks>.md: behaviour terverifikasi, state machine, authorization, audit,
   error behavior, rollout, non-goal.

Gunakan CONTEXT.md + ADR. Jangan menambah perilaku tak berbukti. Tandai bagian yang bergantung
pada pertanyaan terbuka.
```

## STEP W3 — Ticketing

```
/to-tickets

Pecah spec <konteks> jadi vertical slice demoable. Simpan tiket Markdown di
.scratch/<konteks>/issues/NN-<slug>.md.

Tiap tiket WAJIB: hasil & nilai pengguna; acceptance criteria terverifikasi tanpa pengetahuan
pribadi; ruang lingkup & non-goal; area codebase (paket Go / folder React sesuai struktur wajib);
rule Pega sumber (class/nama/tipe + path); dependency & blocker; constraint keamanan/data/operasional;
perintah verifikasi eksplisit; untuk perubahan data/production: migration/rollout/rollback.

Definition of Ready: tandai 'ready-for-agent' hanya bila lengkap dan tidak bergantung pada
pertanyaan terbuka; jika bergantung -> 'needs-info' + sebut pertanyaan + pemiliknya. Tiket
/to-tickets tidak perlu /triage.
```

## STEP W4 — Konsolidasi steering & paket BRD final

```
Konsolidasikan FASE A + W0-W3 jadi paket "BRD lengkap" yang koheren di OUTPUT_HASIL_RNM:
1. Pastikan CLAUDE.md, CONTEXT.md, discovery/context-map.md, docs/agents/*.md konsisten dengan ADR.
2. docs/README.md: index seluruh artefak (discovery/, ADR, BRD, spec, tiket, context-map,
   open-questions terkonsolidasi).
3. Verifikasi silang: tiap ADR dirujuk >=1 BRD/spec/tiket; tiap istilah CONTEXT.md yang dipakai
   BRD punya bukti; tak ada tiket 'ready-for-agent' yang bergantung pada pertanyaan terbuka.

Laporkan ringkas: yang lengkap, yang masih pertanyaan terbuka (+pemilik), dan yang tidak bisa
diselesaikan tanpa DDL Oracle / body stored procedure / data nyata. Jangan menutup gap dengan tebakan.
```

---

## Alur skill tambahan (sesuai kebutuhan)

| Kondisi | Skill |
| --- | --- |
| Perlu bukti runnable atas satu pertanyaan desain | `/handoff` → `/prototype` → `/handoff` |
| Tiket siap dikerjakan (fase implementasi, di luar cakupan prompt ini) | sesi baru → `/implement <path tiket>` |
| Bug sulit / regression | `/diagnosing-bugs` |
| Sebelum merge | `/code-review <base>` |
| Konteks membengkak / ganti fase | `/handoff` lalu sesi baru |

## Batas keamanan & approval manusia

Butuh approval manusia sebelum: operasi destruktif, migration DB, perubahan authn/authz, akses
credential/regulated data, perubahan sistem eksternal. Jangan masukkan secret/token/data
pelanggan/production dump ke prompt/kode/test. Endpoint & host internal = konfigurasi/env var,
bukan literal.

---

### Catatan

- **FASE A (discovery) wajib selesai & di-review manusia sebelum FASE B.** Ini yang memastikan
  Claude "mempelajari dulu aplikasi & flow XML secara lengkap" sebelum menulis BRD/ADR/spec/tiket.
- Prompt ini mencakup **sampai BRD lengkap** (steering + ADR + spec + ticketing). Implementasi
  kode (`/implement`) dimulai kemudian, satu tiket per sesi, dari konteks baru.
```