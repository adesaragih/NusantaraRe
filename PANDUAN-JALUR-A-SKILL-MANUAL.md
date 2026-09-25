# Panduan Jalur A — menjalankan skill Matt Pocock secara manual

> **INI DOKUMEN PANDUAN, BUKAN PROMPT.** Jangan menempel seluruh isi file ini ke Claude.
> Ikuti langkahnya; salin HANYA blok prompt yang ditandai pada langkah yang sedang Anda kerjakan.

> **STATUS TERKINI (diverifikasi di disk):**
> - LANGKAH 1 (D3/D4) — SUDAH SELESAI (modules/*.md, glossary.md, context-map.md,
>   understanding-report.md, D3-D4-CLOSING-REPORT.md sudah ada).
> - LANGKAH 2 (setup tracker) — SUDAH SELESAI (docs/agents/*.md, docs/adr/, .scratch/ sudah ada).
> - SISA SETUP — tinggal membuat CLAUDE.md (tempel "prompt keputusan CLAUDE.md" ke sesi setup
>   yang menunggu). Setelah itu setup + FASE A tuntas.
> - BERIKUTNYA — langsung ke **LANGKAH 3** (grilling FASE B). LANGKAH 1 & 2 di bawah TINGGAL ARSIP.

Anda memilih Jalur A: Anda yang mengetik command skill sendiri (skill ber-`disable-model-invocation:
true`, jadi hanya bisa dipicu manusia), lalu Claude melanjutkan di dalam kerangka skill itu.

**Penting (hasil membaca SKILL.md wayfinder langsung):**
- `wayfinder` adalah skill **perencanaan AWAL** — memetakan pekerjaan besar yang masih "berkabut"
  sebagai peta tiket keputusan di **issue tracker**. Discovery Anda **sudah jauh melewati** titik
  itu (D0-D2 tuntas, 57 OQ terpetakan, 20 konteks tertelusur). Jadi **jangan** menjalankan
  `/wayfinder` untuk D3/D4 — D3/D4 hanya menyintesis hasil, bukan memetakan kabut.
- Wayfinder & grill-with-docs **butuh issue tracker**. Folder ini bukan repo git dan belum di-setup.
  SKILL.md: bila tracker tidak tersedia, jalankan `/setup-matt-pocock-skills` dulu; default-nya
  "local-markdown tracker".

Urutan yang benar untuk Jalur A ada di bawah. Ikuti berurutan.

---

## LANGKAH 1 — Selesaikan D3/D4 dulu (TANPA skill)

D3/D4 adalah sintesis discovery; tidak butuh skill apa pun. Tempel isi blok kode dari file
`PROMPT-D3-D4-CLOSING.md` ke Claude Code. Setelah selesai, FASE A ditutup dan Anda punya:
`glossary.md`, `context-map.md`, `understanding-report.md`, `D3-D4-CLOSING-REPORT.md`, dan 20
`modules/*.md`.

Jangan lompat ke skill sebelum langkah ini selesai — FASE B butuh artefak ini sebagai bukti.

---

## LANGKAH 2 — Setup tracker sekali (Anda ketik command-nya)

Di Claude Code, ketik sendiri (autocomplete: ketik `/` lalu cari):

```
/mattpocock-skills:setup-matt-pocock-skills
```

Setelah command aktif, tempel prompt pengarah ini:

```
Konteks: proyek migrasi Nusantara Re (Pega -> Go + React + Oracle). Folder kerja
D:\XML\RNM_BRD\OUTPUT_HASIL_RNM bukan repo git dan tidak punya issue tracker eksternal.

Setup skill Matt Pocock untuk repo ini dengan pilihan yang PALING SEDIKIT menambah dependensi:
- Issue tracker: gunakan LOCAL-MARKDOWN tracker, disimpan di OUTPUT_HASIL_RNM\.scratch\<konteks>\
  (issues/ untuk tiket, spec.md untuk spesifikasi). Jangan memasang tracker eksternal.
- Domain docs: single-context. Kandidat seed CONTEXT.md sudah ada di discovery\glossary.md.
- ADR: docs\adr\ di dalam OUTPUT_HASIL_RNM.
- Triage labels: default (needs-triage, needs-info, ready-for-agent, ready-for-human, wontfix).

Batas: tulis HANYA ke OUTPUT_HASIL_RNM. Korpus D:\XML\RNM_BRD\ dan D:\XML\nusantara-re\ READ-ONLY.
Gunakan hasil FASE A di OUTPUT_HASIL_RNM\discovery\ sebagai sumber bukti (bukti path+rule).
Laporkan lokasi tracker/docs yang disepakati, lalu berhenti - jangan mulai grilling.
```

Ini membereskan syarat "tracker" yang diminta wayfinder/grill sekaligus, tanpa memasang git/issue
tracker eksternal.

---

## LANGKAH 3 — FASE B: grilling per konteks (Anda ketik command-nya) — LANGKAH AKTIF SEKARANG

FASE B mematangkan requirement satu konteks pada satu waktu. Urutan direkomendasikan di
`discovery\D3-D4-CLOSING-REPORT.md` §4: **konteks pertama = Claim — Life** (cakupan full, hanya 5
OQ pemblokir, modul terkecil & paling terpisah).

**Sebelum grilling:** siapkan jawaban 5 OQ pemblokir Claim — Life (OQ-002, OQ-018, OQ-020, OQ-021,
OQ-029) dari pihak terkait (DBA / IT-infra / Product+UW / IAM). Yang tak terjawab tetap OQ terbuka.

Di Claude Code, **dari dalam folder OUTPUT_HASIL_RNM\**, ketik sendiri:

```
/mattpocock-skills:grill-with-docs
```

Setelah command aktif, tempel PROMPT SIAP-PAKAI ini (khusus konteks pertama, sudah terisi — tidak
perlu diedit):

```
Tujuan: mematangkan requirement migrasi konteks Claim — Life (migrasi Nusantara Re
Pega -> Go + React + Oracle), sampai bisa diimplementasikan tanpa penjelasan lisan.

Bahan bukti (READ-ONLY, dari FASE A): discovery\modules\Claim Life.md,
discovery\flows\Claim Life.md, discovery\flows\_SUMMARY-claim.md,
discovery\context-map.md §2.6, discovery\understanding-report.md §2.4,
discovery\glossary.md, discovery\open-questions.md,
discovery\inventory\_oq011-konflik-isi.md. Korpus D:\XML\RNM_BRD\ dan D:\XML\nusantara-re\
READ-ONLY. Tulis HANYA ke OUTPUT_HASIL_RNM\.

Disiplin (WAJIB): klaim menyebut bukti path+rule; rule class/nama/tipe dari <pxInsName>;
nomor OQ dari register open-questions.md; jangan tebak schema/kode/rumus/isi SP/JSON;
keputusan sulit dibalik -> ADR di docs\adr\; istilah -> seed dari glossary.md.

5 OQ pemblokir Claim — Life yang perlu diselesaikan di sesi ini: OQ-002 (isi stored procedure
klaim Life), OQ-018 (korpus dari dev atau prod), OQ-020 (arti kode status/PaymentType),
OQ-021 (peran di balik identitas hardcode), OQ-029 (IsPEGAPROD). Grill saya satu per satu untuk
kelima OQ ini plus perilaku/aktor/non-goal/state&guard/error/audit/migration konteks Claim — Life.
Untuk tiap jawaban saya, catat sebagai keputusan (CONTEXT.md/ADR) dengan bukti; yang tidak saya
jawab tetap OQ terbuka (jangan ditebak). Jangan lanjut ke spec/tiket sampai grilling matang.
```

**Untuk konteks BERIKUTNYA** (setelah Claim — Life), pakai template generik ini — ganti `<KONTEKS>`
dan daftar bahan sesuai konteksnya (lihat urutan di D3-D4-CLOSING-REPORT.md §4):

```
Tujuan: mematangkan requirement migrasi konteks <KONTEKS> (Nusantara Re, Pega -> Go + React +
Oracle). Bahan bukti (READ-ONLY): discovery\modules\<modul>.md, discovery\flows\<konteks>.md,
discovery\context-map.md, discovery\understanding-report.md, discovery\glossary.md,
discovery\open-questions.md, discovery\inventory\_oq011-konflik-isi.md. Korpus & nusantara-re
READ-ONLY; tulis HANYA ke OUTPUT_HASIL_RNM\.
Disiplin: bukti path+rule; rule class/nama/tipe dari <pxInsName>; nomor OQ dari register; jangan
tebak schema/kode/rumus/SP/JSON; keputusan sulit dibalik -> ADR; istilah -> seed dari glossary.md.
Selesaikan OQ pemblokir konteks ini yang saya jawab; sisanya tetap terbuka. Grill saya soal
perilaku/aktor/non-goal/state&guard/error/audit/migration. Jangan lanjut spec/tiket sampai matang.
```

---

## LANGKAH 4 — FASE B lanjutan (Anda ketik command-nya)

Setelah grilling satu konteks matang, lanjutkan (Anda ketik sendiri, satu per satu):

```
/mattpocock-skills:to-spec
```
lalu arahkan: sintesis keputusan grilling <KONTEKS> jadi docs\brd\BRD-<konteks>.md +
docs\spec\spec-<konteks>.md (behaviour terverifikasi, state machine, authorization, audit, error,
non-goal), pakai CONTEXT.md + ADR, tandai bagian yang bergantung pada OQ yang belum terjawab.

```
/mattpocock-skills:to-tickets
```
lalu arahkan: pecah spec <KONTEKS> jadi vertical slice ke .scratch\<konteks>\issues\NN-<slug>.md;
tiap tiket punya acceptance criteria terverifikasi, area codebase (struktur folder Go/React target),
rule Pega sumber (class/nama/tipe), dependency, perintah verifikasi; tandai 'ready-for-agent' hanya
bila tidak bergantung OQ terbuka, selain itu 'needs-info'.

---

## Aturan tetap (berlaku semua langkah)

- Skill Matt Pocock HANYA Anda yang picu (ketik command). Claude tidak memanggilnya.
- Bila Claude melaporkan sebuah command skill tidak bisa dijalankan / minta prasyarat (mis.
  tracker), ia harus BERHENTI dan memberi tahu, bukan mereplikasi alur skill diam-diam.
- GATE manusia: FASE B hanya menutup OQ yang Anda jawab; OQ pemblokir yang butuh DBA/Product+UW/
  Actuarial/Finance/IAM tetap terbuka sampai dijawab pihak itu.
```
