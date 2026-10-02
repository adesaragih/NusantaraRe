# inventory/ — STEP D1

Status: **SELESAI** — 20 dari 20 modul, 9.369 dari 9.369 file (100%).
Laporan penutup D1: `../D1-CLOSING-REPORT.md`.

## Isi folder

| Berkas | Isi |
| --- | --- |
| `<modul>.md` | inventaris satu modul, format **7 bagian** yang seragam |
| `_summary.md` | rekap lintas modul per batch + **rekap final korpus (§21–§22)**. **FINAL** |
| `_oq011-konflik-isi.md` | register identitas rule yang isinya berbeda antar modul (OQ-011) |

## Format 7 bagian tiap `<modul>.md`

1. Ringkasan terukur — jumlah file (dicocokkan dengan `../README.md` §5.1), identitas unik, class
   unik, jenis SQL RDBList, prefix kunci RDBList, tabel jumlah per tipe rule.
2. Sebaran class Pega — ditandai mana yang bukan class aplikasi (→ OQ-009).
3. Daftar rule per tipe — `File | Class | Nama rule | Fakta dari tag | Asal salinan`.
4. Tabrakan nama di dalam modul — nama file dipakai >1 tipe; `When` bernama sama di >1 class.
5. Keluarga klon — satu `<pzOriginalInstanceKey>` dipakai >1 rule.
6. Objek database — tabel/view dari `<pyBrowseSQL>` + stored procedure yang dipanggil.
7. Integrasi eksternal — `ConnectREST` sebagai konfigurasi; URL literal ditandai sebagai temuan.

## Aturan D1 (tidak berubah sejak batch 1)

- **Hanya struktur + metadata.** Tidak menyimpulkan logika bisnis — itu STEP D2.
- **Identitas rule** dari `<pxInsName>`, bukan `<pzOriginalInstanceKey>` (lihat `../README.md` §6.4).
  Tipe rule dari field pertama `<pzOriginalInstanceKey>`, atau `<pxObjClass>` bila tag itu tidak ada.
  **Bukan dari nama folder** — folder terbukti bisa berbeda dari tipe rule sebenarnya.
- **Jangan baca file utuh.** Korpus 1.319,3 MB, file terbesar 8,2 MB. Grep terarah, baca range.
- **Perbandingan isi rule wajib ternormalisasi** — buang tag volatil/provenance, urutkan baris,
  lalu hash. `md5sum` mentah **dilarang** menjadi dasar klaim "berbeda". Daftar tag lengkap ada di
  `_oq011-konflik-isi.md`.
- **Pertahankan `@`** saat mengekstrak nama objek SQL — menandai Oracle database link
  (`TABEL@DBLINK`).
- Setiap angka disertai perintah audit. Yang tidak terukur ditulis `belum terukur`.

## Cakupan per batch

| Batch | Domain | Modul | File | Status |
| ---: | --- | ---: | ---: | --- |
| 1 | Treaty inward | 4 | 1.149 | selesai |
| 2 | Facultative inward | 5 | 6.667 | selesai |
| 3 | Claim non-facultative | 6 | 871 | selesai |
| 4 | Life, master & outward | 5 | 682 | selesai |
