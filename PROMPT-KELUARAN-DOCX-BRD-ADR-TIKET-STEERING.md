# PROMPT — keluarkan BRD, ADR, Tiket, dan Steering sebagai `.docx`, sekarang

> Tempel berkas ini sebagai brief. Muat skill **`docx`** sebelum menulis berkas Word.
> Bukan `/to-spec`, bukan `/to-tickets` — keduanya sudah selesai dan keluarannya markdown.
>
> Empat dokumen, seluruhnya **dikeluarkan sekarang** dari bahan yang ada. Tidak ada yang ditunda,
> tidak ada yang menunggu keputusan siapa pun. Boleh satu sesi, atau satu dokumen per sesi —
> tiap bab §2–§5 berdiri sendiri.
>
> Revisi 25 September 2026 siang: pohon **`jefri\`** ditemukan sesudah versi pertama ditulis.
> Ia menutup tiga modul Fakultatif. Seluruh angka di bawah sudah memuatnya.

---

## 0. BAHAN — hasil audit `OUTPUT_HASIL_RNM\` 25 September 2026

Bahan tersebar di **tiga pohon**; ketiganya dipakai. **20 dari 20 modul** punya spec dan tiket.

| Pohon | Modul | Spec | Tiket | ADR |
| --- | ---: | --- | ---: | ---: |
| `.scratch\<modul>\` | 13 | `spec*.md` | `issues\NN-*.md` — **188** | `docs\adr\` — **42** |
| `dastin\_migration-docs\<modul>\` | 4 — Claim Non Prop · Komite Claim Non Prop · Treaty In · Treaty In Adjustment | `2-to-spec\SPEC-*.md` | `3-to-tickets\` atau `5-tiket\issues\` — **117** | `docs\adr\` per modul — **56** |
| `jefri\OUTPUT FIX\` | 3 — NB FacIn · RNW Fac In · Endorsment Fac In, plus **Fac Out** yang dipisah jadi menu sendiri | `04-spec\` — 11 berkas | `05-tickets\` — **60** *(akar 16 · `rnw\` 8 · `edm\` 22 · `facout\` 14; berkas `00-INDEKS*` bukan tiket)* | `adr\` — **7** |
| **Jumlah** | **20** | | **365** | **105** |

Pendukung: `discovery\` *(80 berkas)* · `CLAUDE.md` · `CONTEXT.md` ·
`Diagram-Skema-Tabel-NusantaraRe.xlsx` · `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md` ·
`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md` · `jefri\OUTPUT FIX\03-keputusan\` · `04-kuesioner\` ·
`08-flat\` *(78 tabel, DDL, workbook)* · `10-audit\` · `steering\GLOSARIUM.md`.

**Keadaan tiket:** `.scratch` nol tertahan. `jefri` **24 siap, 36 tertahan** — seluruh EDM dan
Fac Out `blocked`, sebagian besar oleh tiket lain di rantainya. `dastin` dihitung executor.

**Keadaan bahan Fac In, dinyatakan spec-nya sendiri** *(`04-spec\03`, Out of Scope 13 butir)*:
isi **35 stored procedure** jalur tulis produksi tidak ada di ekspor · **nol DDL** · **604 dropdown**
tanpa daftar nilai · baris tabel keputusan underwriting tidak terekspor · **77 rule** perlu ekspor
ulang · **27 pertanyaan** kuesioner ke UW, Product, IT, DBA. Ini **fakta yang dicatat** di dokumen,
bukan alasan menunda dokumen.

⚠️ Angka di atas hasil audit; executor **mengukur ulang** dari berkas dan memakai angkanya
sendiri. Selisih dilaporkan, bukan disesuaikan.

---

## 1. KELUARAN

Folder: **`OUTPUT_HASIL_RNM\keluaran-docx\`**

| | Berkas | Pembaca |
| --- | --- | --- |
| A | `BRD-NusantaraRe.docx` | manajemen · Product+Underwriting · vendor |
| B | `ADR-NusantaraRe.docx` | arsitek · pengembang senior |
| C | `Tiket-<bounded-context>.docx` — **sembilan** berkas | pengembang · pemimpin tim |
| D | `Steering-NusantaraRe.docx` | steering committee |

Tiap `.docx` disusun dari **markdown gabungan** yang disimpan di samping-nya dengan nama sama.

---

## 2. A · BRD

| Bab | Isi | Sumber |
| ---: | --- | --- |
| 1 | Latar dan tujuan: Pega → Go + React, Oracle tetap, **nol pemanggilan procedure dari SQL** | `CLAUDE.md` · `CONTEXT.md` |
| 2 | Lingkup: 20 modul, 9 bounded context, tiga pohon bahan | `discovery\context-map.md` §1 · `docs\agents\issue-tracker.md` |
| 3 | Peta konteks dan hubungan antar konteks | `discovery\context-map.md` |
| 4–23 | **Satu bab per modul**: masalah · solusi · user story ringkas · keputusan dagang · penyimpangan sadar · butir terbuka | `spec*.md` tiap modul; untuk Fac In: `jefri\04-spec\01`, `03`, `04`, `05–10` |
| 24 | ⭐ **Kepastian bahan per modul** — apa yang terverifikasi dari korpus, apa yang menunggu pihak lain, dan pihak mana | bab *Out of Scope* tiap spec · `jefri\04-spec\03` · `jefri\04-kuesioner\` |
| 25 | Model data ringkas: tabel dan relasinya per konteks, **tanpa DDL** | `Diagram-Skema-Tabel-NusantaraRe.xlsx` · `STRUKTUR-TABEL-*.md` · `4-erd-dan-tabel-datar\` · `jefri\08-flat\Tabel-Flat-Lintas-Siklus.xlsx` |
| 26 | Keputusan arsitektur yang paling mengikat — ringkasan, merujuk dokumen B | `docs\adr\00-INDEKS.md` |
| 27 | Butir terbuka per pemilik peran: DBA · Product+Underwriting · Actuarial · Finance · IAM · IT | `discovery\open-questions.md` · seluruh `[terbuka]` · `jefri\04-kuesioner\` |
| 28 | Risiko yang diterima sadar | seluruh `[penyimpangan sadar]` · `jefri\...\K-065` |

Bab 24 adalah yang membuat BRD ini jujur tanpa perlu ditunda: pembaca melihat **per modul** mana
yang sudah pasti dan mana yang belum, dengan pemiliknya.

User story **diringkas per modul**; tiap bab menutup dengan rujukan ke spec-nya.

---

## 3. B · ADR

**105 ADR dimuat utuh** — ADR pendek, memangkasnya menghilangkan alasannya.

Tiga seri bernomor sama dengan isi berbeda. Cara menyajikannya **di dalam dokumen saja** — nol
berkas diganti nama:

| Bagian | Seri | Heading |
| --- | --- | --- |
| I | `docs\adr\` 0001–0042 | `ADR-U-nnnn — judul` |
| II | `dastin\...\docs\adr\` per modul | `ADR-D-<modul>-nnnn — judul` *(`CNP` · `KCNP` · `TI`)* |
| III | `jefri\OUTPUT FIX\adr\` 0001–0007 | `ADR-F-nnnn — judul` |

**Lampiran 1 — konkordansi**, 105 baris: nomor asli → awalan seri → judul → berkas asal.

**Lampiran 2 — pasangan sebidang**: ADR dari seri berbeda yang membahas hal yang sama, disajikan
berdampingan dengan satu kalimat persamaan dan satu kalimat perbedaan. Sekurangnya tiga kelompok
sudah diketahui: **uang** *(`U-0003` · `U-0016` · `D-CNP-0003` · `D-CNP-0007` · `F-0004` ·
`F-0005` · `F-0007`)*, **tabel akar bersama** *(`U-0025` · `U-0026` lawan keputusan K-064 · K-065 ·
K-069 di `jefri\04-spec\11`)*, **penamaan tabel**. **Tidak diputuskan mana yang berlaku.**

Format frontmatter berbeda antar seri **diseragamkan tampilannya** di Word, bukan di sumbernya.

---

## 4. C · Tiket — sembilan berkas, badan lengkap

| Berkas | Modul di dalamnya | Sumber |
| --- | --- | --- |
| `Tiket-claim-life.docx` | Claim Life | `.scratch` |
| `Tiket-life-offer.docx` | PremiumList Life · Endorsement Life | `.scratch` |
| `Tiket-life-master.docx` | Master Product Name Life · Master Contract Retro Life | `.scratch` |
| `Tiket-claim-nonlife.docx` | Claim Fac In · Claim Prop · Claim Non Prop | `.scratch` · `dastin` |
| `Tiket-komite.docx` | Komite Claim FacIn · Life · Prop · Non Prop | `.scratch` · `dastin` |
| `Tiket-treaty-realisasi.docx` | NB Treaty In · EDM Treaty In | `.scratch` |
| `Tiket-treaty-master.docx` | Treaty In · Treaty In Adjustment | `dastin` |
| `Tiket-treaty-arrangement.docx` | Treaty Contract Out | `.scratch` |
| `Tiket-facultative-inward.docx` | NB FacIn 16 · RNW 8 · EDM 22 · Fac Out 14 — **60 tiket** | `jefri\05-tickets\` |

Tiap berkas:

- **Matriks status** di depan: per modul — siap · tertahan · `wontfix` · digantikan. Dihitung dari
  baris `Status:` tiap tiket. Untuk Fac In, 36 tiket tertahan **dicantumkan beserta `Blocked by`-nya**
  supaya rantainya terbaca.
- **Badan tiket lengkap**, per modul lalu per nomor. Heading `<modul> · NN — judul`; nomor jefri
  memakai awalan `E`/`F`/`R` — dipertahankan apa adanya.
- Tiket **digantikan** *(NB Treaty In 24–28)* dan `wontfix` **ikut dimuat**, ditandai di heading.
- `claim-non-prop\3-to-tickets\TICKETS.md` dan `jefri\05-tickets\00-INDEKS*.md` adalah indeks,
  bukan tiket — dipakai untuk mencocokkan cacah, bukan disalin.

---

## 5. D · Steering

Pendek — **6–10 halaman**. Keputusan dan angka, bukan cerita.

| § | Isi |
| ---: | --- |
| 1 | Satu halaman keadaan: 20/20 modul · 365 tiket · 105 ADR di tiga seri · Fac In 24 siap dari 60 |
| 2 | Yang siap dibangun sekarang, per bounded context — dari matriks status dokumen C |
| 3 | Urutan pembangunan yang diusulkan, dan gerbang antar tiket yang menentukannya |
| 4 | Butir terbuka yang **memblokir pembangunan**, per pemilik peran, dengan cacah — termasuk Fac In: `ALL_SOURCE` 35 procedure dan DDL *(DBA)* · 604 dropdown *(Product)* · 77 rule ekspor ulang *(pemilik export Pega)* |
| 5 | **Keputusan yang diminta komite** — tiap butir dengan pilihan dan akibatnya, tanpa rekomendasi: |
| 5a | ⛔ **`T_GENERAL_POLIS` dan `T_WORK_POLIS`** — satu tabel fisik, **dua rancangan kolom**: Treaty In *(79 medan + `REMARK`, lembar `Diagram-Skema`)* lawan Fac In *(71 kolom, dari 113 contoh dokumen)*; sisi Fac In menyatakan rekonsiliasi **gugur** *(K-065, "diagram skema rumah tidak dipakai")*. Siapa pemilik rancangan tabel itu |
| 5b | ⛔ **Kolom pembeda lini** — K-069 menyatakan di luar proyek dan mencabutnya; `ADR-U-0026` mewajibkannya. Satu tabel akan memuat baris Fac In dan Treaty In tanpa pembeda |
| 5c | **Penamaan tabel untuk properti Pega yang sama** — `T_POLIS_CEDING` lawan `T_CEDINGCOLIST`, `T_POLIS_INSTALMENT` lawan `T_LISTINSTALLMENT`, `T_POLIS_SPREADING` lawan `T_SPREADINGLIST`, `T_POLIS_QUOTATION` lawan `T_QUOTATIONDATA` |
| 5d | **Tiga seri ADR** — dibiarkan dengan awalan, atau dilebur; siapa yang memutus pasangan sebidang |
| 5e | **Jadwal Fakultatif** — 36 dari 60 tiket tertahan; kapan blocker §4 ditutup, dan oleh siapa |
| 6 | Risiko teratas — maksimal tujuh, satu paragraf tiap risiko, dengan angka |
| 7 | Langkah berikutnya, 30 hari |

---

## 6. CARA MENGERJAKAN

1. Muat skill `docx`. Untuk tiap dokumen: susun markdown gabungan → jadikan `.docx`.
2. **Lihat hasilnya** sesuai cara skill `docx`: konversi ke PDF, render halaman, baca gambarnya.
   Periksa daftar isi, tabel tidak terpotong, heading berjenjang. Dokumen yang belum dilihat belum
   selesai.
3. Berkas tersegel **dibaca, tidak disunting**: `grilling-ronde-*.md` · `VERIFIKASI-*.md` ·
   `KOREKSI-*-DIJALANKAN.md` · `PROMPT-*.md` · seluruh isi `dastin\` · seluruh isi `jefri\`.
4. ⚠️ `jefri\OUTPUT FIX\steering\` berisi **glosarium dan panduan kerja** — bukan dokumen steering
   committee. Ia sumber kosakata untuk bab Fac In, bukan bahan dokumen D.
5. Nol fakta baru dari korpus. Bila sesuatu tidak ada di spec, tiket, ADR, discovery, atau catatan
   keputusan, ia tidak ada di dokumen.

---

## 7. DISIPLIN

Nol nama orang · nol nomor polis harfiah · nol cuplikan data produksi · nol secret, token, alamat
layanan — ditulis "konfigurasi". Nol DDL di BRD dan Steering. Butir `[terbuka]` **tidak ditutup**;
pasangan ADR sebidang dan lima butir §5 Steering **tidak dimenangkan**.

---

## 8. TANDA BERHASIL

- Dua belas `.docx` di `keluaran-docx\` — A, B, D, dan sembilan C — masing-masing dengan markdown
  gabungannya, masing-masing sudah dilihat renderannya.
- BRD memuat 20 bab modul dan bab **Kepastian bahan per modul**.
- Dokumen B: 105 ADR utuh, konkordansi 105 baris, nol rujukan ADR tanpa awalan seri.
- Dokumen C: cacah tiket per berkas cocok dengan cacah berkas sumbernya; `Tiket-facultative-inward`
  memuat 60 tiket dengan `Blocked by` yang terbaca.
- Steering §5 memuat lima keputusan yang diminta, tiap butir dengan pilihan dan akibat — nol
  rekomendasi.
- Laporan akhir: cacah terukur lawan cacah §0, dan bab telemetri yang memisahkan terukur dari
  taksiran.

---

*Disusun 25 September 2026. Direvisi siang hari yang sama sesudah pohon `jefri\` ditemukan.*
