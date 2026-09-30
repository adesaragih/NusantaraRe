# Issue tracker: Local Markdown

Issues dan spec untuk proyek ini hidup sebagai berkas markdown di **folder modulnya**:
**`OUTPUT_HASIL_RNM/APP_RNM/modul/<nama>/docs/`**.

> ⚠️ **Struktur tim satu folder per modul (30-09-2026).** Dulu `OUTPUT_HASIL_RNM/.scratch/<nama-panjang>/`;
> kini setiap dokumen modul — spec, tiket, grilling — tinggal bersama kode modulnya, dan dua puluh folder
> modul (satu per folder korpus) sudah berdiri. `<nama>` = nama folder `.scratch` lama tanpa tanda hubung
> (`claim-life` → `claimlife`, `nb-treaty-in` → `nbtreatyin`). Dokumen lintas modul (ADR, `CONTEXT.md`)
> di `OUTPUT_HASIL_RNM/docs/bersama/`.

> **Catatan lokasi (penting).** Template bawaan skill memakai `.scratch/` di **root repo**.
> Di proyek ini root (`D:\XML\RNM_BRD\`) adalah **korpus ekspor Pega yang READ-ONLY** — 20 folder
> modul, 9.369 berkas `.xml`, tidak boleh disentuh. Seluruh keluaran agent karena itu berada di
> bawah `OUTPUT_HASIL_RNM\`. Setiap path di berkas ini sudah disesuaikan.
> `D:\XML\RNM_BRD\` bukan repository git dan tidak punya issue tracker eksternal.

## Conventions

- Satu modul per direktori: **`OUTPUT_HASIL_RNM/APP_RNM/modul/<nama>/docs/`**
- Spec berada di `OUTPUT_HASIL_RNM/APP_RNM/modul/<nama>/docs/spec.md`
- Issue implementasi satu berkas per tiket di
  `OUTPUT_HASIL_RNM/APP_RNM/modul/<nama>/docs/issues/<NN>-<slug>.md`, bernomor dari `01`,
  **tidak pernah** satu berkas gabungan
- Status triage dicatat sebagai baris `Status:` di dekat bagian atas tiap berkas issue
  (lihat `triage-labels.md` untuk string perannya)
- Komentar dan riwayat percakapan di-*append* di bagian bawah berkas di bawah heading `## Comments`

### Slug konteks yang dipakai proyek ini

`[terverifikasi]` Sembilan bounded context sudah ditetapkan di
`discovery/context-map.md` §1. Slug yang dianjurkan, agar tracker sejajar dengan peta konteks:

| Slug | Bounded context | Modul |
| --- | --- | --- |
| `claim-life` | Claim — Life | Claim Life |
| `life-offer` | Life — Penawaran & Premium List | PremiumList Life, Endorsement Life |
| `treaty-arrangement` | Treaty Arrangement | Treaty Contract Out |
| `claim-nonlife` | Claim — Non-Life | Claim Fac In, Claim Prop, Claim Non Prop |
| `life-master` | Life — Master | Master Product Name Life, Master Contract Retro Life |
| `komite` | Komite | 4 modul Komite |
| `treaty-realisasi` | Treaty Inward — Realisasi & Endorsement | NB Treaty In, EDM Treaty In |
| `treaty-master` | Treaty Inward — Master & Akseptasi | Treaty In, Treaty In Adjustment |
| `facultative-inward` | Facultative Inward | NB FacIn, RNW Fac In, Endorsment Fac In |

## When a skill says "publish to the issue tracker"

Buat berkas baru di bawah `OUTPUT_HASIL_RNM/APP_RNM/modul/<nama>/docs/` (buat direktori `docs/`
bila belum ada; folder modulnya sudah berdiri untuk kedua puluh modul korpus).

## When a skill says "fetch the relevant ticket"

Baca berkas pada path yang dirujuk. Pengguna biasanya memberikan path atau nomor issue langsung.

## Wayfinding operations

Dipakai `/wayfinder`. **Map** adalah satu berkas dengan satu berkas **child** per tiket.

- **Map**: `OUTPUT_HASIL_RNM/.scratch/<effort>/map.md` (bagian Notes / Decisions-so-far / Fog).
- **Child ticket**: `OUTPUT_HASIL_RNM/.scratch/<effort>/issues/NN-<slug>.md`, bernomor dari `01`,
  dengan pertanyaannya di badan berkas. Baris `Type:` mencatat jenis tiket
  (`research`/`prototype`/`grilling`/`task`); baris `Status:` mencatat `claimed`/`resolved`.
- **Blocking**: baris `Blocked by: NN, NN` di dekat bagian atas. Tiket terbuka blokirnya bila setiap
  berkas yang disebut sudah `resolved`.
- **Frontier**: pindai `OUTPUT_HASIL_RNM/.scratch/<effort>/issues/` untuk berkas yang terbuka,
  tidak terblokir, dan belum diklaim; nomor terkecil menang.
- **Claim**: set `Status: claimed` dan simpan sebelum mulai bekerja.
- **Resolve**: *append* jawabannya di bawah heading `## Answer`, set `Status: resolved`, lalu
  *append* penunjuk konteks (ringkas + tautan) ke Decisions-so-far di `map.md`.

## Batas yang tetap berlaku di proyek ini

- Korpus `D:\XML\RNM_BRD\` **READ-ONLY**; `D:\XML\nusantara-re\` ⛔ **DI-BLACKLIST** — jangan
  dibaca atau dijadikan pembanding (keputusan work owner 2026-09-14). Tulis **hanya** ke
  `OUTPUT_HASIL_RNM\`, yang sekaligus **repo target tunggal** untuk dokumen dan kode.
- Setiap klaim perilaku wajib bukti **`path + rule`** (rule ditulis `class / nama / tipe` dari
  `<pxInsName>`) — lihat `discovery/flows/_METHOD.md`.
- Nomor OQ diambil dari register `discovery/open-questions.md`, **bukan dari ingatan**.
- Jangan memasukkan secret/token/data pelanggan/dump production ke tiket, spec, kode, atau test.
- Endpoint & host internal = konfigurasi/env var, **bukan** literal.
- Uang material **jangan** `float`.
