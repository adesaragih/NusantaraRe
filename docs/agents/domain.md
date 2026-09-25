# Domain Docs

Cara skill engineering mengonsumsi dokumentasi domain proyek ini ketika menjelajah basis kode.

**Layout: single-context.** Satu `CONTEXT.md` + `docs/adr/`.

> **Catatan lokasi (penting).** Template bawaan skill memakai root repo. Di proyek ini root
> (`D:\XML\RNM_BRD\`) adalah **korpus ekspor Pega yang READ-ONLY**; seluruh keluaran agent berada di
> bawah `OUTPUT_HASIL_RNM\`. Setiap path di berkas ini sudah disesuaikan.

## Before exploring, read these

- **`OUTPUT_HASIL_RNM/CONTEXT.md`** — glossary domain proyek.
  **Belum dibuat.** Kandidat seed-nya sudah siap: **`OUTPUT_HASIL_RNM/discovery/glossary.md`**
  (169 entri, seluruhnya berkolom Bukti `path + rule`). Lihat §"Seed CONTEXT.md" di bawah.
- **`OUTPUT_HASIL_RNM/docs/adr/`** — ADR yang menyentuh area yang akan dikerjakan.
  **Masih kosong**; ADR pertama baru ditulis di FASE B.

Bila salah satu berkas ini belum ada, **lanjutkan diam-diam**. Jangan menandai ketiadaannya, jangan
menyarankan membuatnya di muka. Skill `/domain-modeling` (dicapai lewat `/grill-with-docs` dan
`/improve-codebase-architecture`) membuatnya secara *lazy* ketika istilah atau keputusan benar-benar
terselesaikan.

## File structure

Single-context (layout proyek ini):

```
OUTPUT_HASIL_RNM/
├── CONTEXT.md                  ← belum ada; seed = discovery/glossary.md
├── docs/
│   ├── adr/                    ← kosong; ADR ditulis di FASE B
│   └── agents/
│       ├── issue-tracker.md
│       ├── triage-labels.md
│       └── domain.md
├── .scratch/                   ← issue tracker local-markdown
│   └── <konteks-slug>/
│       ├── spec.md
│       └── issues/NN-<slug>.md
└── discovery/                  ← keluaran FASE A (sumber bukti)
    ├── README.md
    ├── D1-CLOSING-REPORT.md
    ├── D2-CLOSING-REPORT.md
    ├── D3-D4-CLOSING-REPORT.md
    ├── context-map.md
    ├── understanding-report.md
    ├── glossary.md
    ├── open-questions.md
    ├── inventory/   (20 modul + 3 meta)
    ├── flows/       (20 konteks + 4 sintesis + 2 metode)
    └── modules/     (20 catatan modul)
```

`[terverifikasi]` Layout **single-context** dipilih meski korpus memuat **9 bounded context**
(`discovery/context-map.md` §1). Alasannya: sinyal monorepo tidak ada — tidak ada
`pnpm-workspace.yaml`, `package.json`, `packages/*`, maupun `src/` — dan basis kode target
(Go + React) **belum ada**. Bila kelak repo target terbentuk sebagai multi-package, layout ini dapat
dipromosikan menjadi `CONTEXT-MAP.md` + `CONTEXT.md` per konteks tanpa kehilangan apa pun:
`discovery/context-map.md` sudah memuat pembagiannya.

## Seed CONTEXT.md

`[terverifikasi]` `discovery/glossary.md` ditandai sebagai **kandidat seed `CONTEXT.md`** dan berisi
169 entri dalam empat bagian: 52 istilah domain, 51 kode/enumerasi, 25 nama workbasket,
41 singkatan belum terjabarkan — setiap entri berkolom **Bukti** (`path + rule`).

**Peringatan yang mengikat sebelum menyalinnya menjadi `CONTEXT.md`:** mayoritas entri masih
`arti belum terverifikasi`. Enam OQ menentukan isinya — **OQ-020** (arti seluruh kode/enumerasi),
**OQ-008** (prefix `ASM`/`RNM`/`GCNM`), **OQ-014** (kepanjangan `EDM`), **OQ-057** (`REINSTYPEID`
dan kode `OR` tanpa satu pun nilai literal), **OQ-049** (`SFAGIS`), **OQ-043** (arti hasil
`confirm`/`reject`/`ask`/`banding`/`revise`/`decline`).

**Jangan mengisi kolom mana pun dengan tebakan.** Yang belum terverifikasi tetap ditandai demikian
sampai pemilik peran menjawab. Lihat `discovery/glossary.md` §5.

## Use the glossary's vocabulary

Ketika keluaran Anda menamai konsep domain (judul issue, usulan refactor, hipotesis, nama test),
pakai istilah sebagaimana didefinisikan di `CONTEXT.md` — dan sebelum `CONTEXT.md` ada, di
`discovery/glossary.md`.

Jangan bergeser ke sinonim yang glossary hindari. **Jangan mengarang terjemahan istilah Indonesia,
kepanjangan singkatan, atau arti field dari caption** — aturan ini dibawa dari FASE A dan tetap
berlaku.

Bila konsep yang Anda butuhkan belum ada di glossary, itu sinyal: entah Anda mengarang bahasa yang
tidak dipakai proyek (pikirkan ulang), atau memang ada celah nyata (catat untuk `/domain-modeling`).

## Flag ADR conflicts

Bila keluaran Anda bertentangan dengan ADR yang ada, sampaikan eksplisit alih-alih diam-diam
menimpanya:

> _Bertentangan dengan ADR-0007 (…), tetapi layak dibuka ulang karena…_

## Sumber bukti FASE A

`[terverifikasi]` Seluruh keputusan FASE B harus dapat ditelusuri ke bukti FASE A:

| Butuh | Baca |
| --- | --- |
| Bagaimana sistem bekerja | `discovery/understanding-report.md` |
| Pembagian konteks + cakupan bukti | `discovery/context-map.md` |
| Perilaku satu konteks (alur, guard, objek Oracle) | `discovery/flows/<konteks>.md` |
| Peran & ketergantungan satu modul | `discovery/modules/<modul>.md` |
| Isi satu modul (17 tipe rule) | `discovery/inventory/<modul>.md` |
| Pertanyaan yang belum terjawab | `discovery/open-questions.md` (57 terbuka, 38 memblokir FASE B) |
| Urutan FASE B yang diusulkan | `discovery/D3-D4-CLOSING-REPORT.md` §4 |
