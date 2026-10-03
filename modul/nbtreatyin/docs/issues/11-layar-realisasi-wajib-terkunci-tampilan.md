# 11: Layar realisasi — medan wajib, medan terkunci, dan bagian yang tidak dibangun

**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
**Blocked by:** 02
**Menutup:** AC 45 · 46 · 49 · 50 · 51 · 53 · 54 · 55 · 56 · 77 *(10 AC)* — US 27 · 28 · 29 · 31 · 32

## Hasil & nilai pengguna

Hari ini layar realisasi menuntut sejumlah medan diisi, mengunci sebagian lainnya, dan
menyembunyikan bagian yang tidak berlaku. ⚠️ `[terverifikasi]` **Delapan puluh** bagian layar
diberi syarat tampil yang **tidak mungkin pernah benar** — dan **tidak satu pun** dari yang 80 itu
menyembunyikan sebuah medan; seluruhnya label dan hiasan.

Sesudah tiket ini, pengguna melihat layar yang **menandai medan wajibnya dengan jelas**, mengunci
yang tidak boleh ia ubah, dan ⭐ **tidak memuat 80 bagian mati** yang selama ini ada tanpa pernah
tampil.

## Area codebase

- Lapisan handler: validasi medan wajib
- Antarmuka: penandaan medan wajib dan terkunci
- Antarmuka: daftar pilihan mata uang

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Medan wajib | **27** medan berbeda di **6** layar; sebarannya di `spec.md` §5.11 |
| Medan terkunci | **38**, seluruhnya lewat syarat hanya-baca yang **selalu benar**; setiap kunci mengenai **tepat satu medan** |
| Bagian mati | **91**; ⭐ **80** menempel pada sel tunggal, **11** pada wadah yang memuat medan |
| Pengecualian mata uang | `ReportDefinition\BrowseCurrency_RD.xml` · `BrowseCurrencyTreatyIn_RD.xml` |
| Pembersihan pesan galat | **12** tempat yang menghapus **sebelum** pesan dipasang |

## ADR terkait

- **ADR-0002** — RBAC memakai peran yang sudah ada

## Acceptance criteria

- [ ] 🟡 **AC 45** — layar menuntut medan wajibnya; **27** medan berbeda di **6** layar
- [x] **AC 46** — layar jenjang ketiga **tidak** mewajibkan enam medan yang wajib di layar admin
- [x] **AC 49** — **38** medan terkunci permanen
- [x] **AC 50** — **36** dari 38 di layar jenjang ketiga; **2** di layar biasa
- [x] **AC 51** — setiap kunci mengenai **tepat satu medan**, ⛔ bukan bagian atau tab
- [x] **AC 53** — **80** bagian mati **tidak dibangun**
- [x] **AC 54** — daftar pilihan mata uang **tidak memuat** kode yang dikecualikan
- [x] **AC 55** — kode jenis kontrak ditampilkan **apa adanya** bila keterangannya belum tersedia
- [x] **AC 56** — nilai kode di luar daftar yang dikenal **tetap diterima dan disimpan**
- [x] **AC 77** — pembersihan pesan galat di awal diterima apa adanya — **12** tempat

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **17** | **empat wadah** berisi **104** medan — masih dipakai atau ditinggalkan | ⚠️ menahan **keempatnya saja**, bukan tiketnya |
| **14** | apakah ada mata uang mati lain yang tidak disembunyikan | tidak menahan |

## Perintah verifikasi

1. Kosongkan satu medan wajib, simpan — ⭐ **ditolak**, dengan pesan yang menyebut medannya.
2. Coba sunting medan terkunci — ⭐ **tidak bisa**, dan tampil sebagai terkunci.
3. Cari kode mata uang yang dikecualikan di daftar pilihan — ⭐ **tidak ada**.
4. Kirim kode jenis kontrak yang tidak dikenal — ⭐ **tetap tersimpan**, ditampilkan apa adanya.

## Catatan

⭐ **Delapan puluh bagian mati dibuang tanpa ditanyakan kepada siapa pun** `[keputusan work owner]`
P44 — `[terverifikasi]` sebabnya struktural: membangunnya berbiaya nol, membuangnya juga.
⛔ **Empat wadah berisi 104 medan** menunggu Product & Underwriting dan **tidak dibangun** sebelum
dijawab.

## Hasil implementasi 2026-10-03

- Medan wajib per layar dari `pyRequired` (`models.DaftarMedanWajib`): 25 medan di dua layar
  realisasi + `ProductionDate` (tempat berperan, tiket 05) + `DateofSurvey` (layar survei historis,
  tidak dibangun) = 27 (AC 45). Ditegakkan pada Submit DAN Save (AC 48).
- Pilihan radio/dropdown bersumber "associated values" rule Property — **tidak ada di korpus**
  (TypeTax, IsSurveyReport, StatementType, ClaimType, ClaimPaymentType, DueTo) ⇒ isian teks apa adanya
  (AC 56), tidak dikarang. Approval memakai nilai `1`/`0` dengan teks dari `SaveViewSuggest`
  (Accept/Reject).
- Tidak dibangun: tombol/layar Survey Report (penyimpanan survei tidak dirancang, `[terbuka]`),
  pemilih SOB (hanya untuk XOL Retro), Choose Business R (treaty keluar, JSON).
