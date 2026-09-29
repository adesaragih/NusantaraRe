# Register OQ — Treaty Contract Out

Satu tempat untuk seluruh pertanyaan terbuka modul ini. Rincian dan buktinya tetap di tiket terkait dan di bab
`LAPORAN-GILIRAN.md` tempat OQ itu dibuka; di sini status dan keputusannya.

| OQ | Pokok | Dibuka | Status |
| --- | --- | --- | --- |
| OQ-TCO-01 | bentuk teks tanggal warisan (tujuh bentuk dikenal pengurai) | tiket 01 | terbuka |
| OQ-TCO-02 | arti/bentuk `IUDATE` | tiket 01 | terbuka |
| OQ-TCO-03 | isi hidup `PROPORTIONALLIST`/`OBJECT` | tiket 01 | terbuka |
| OQ-TCO-04 | arti `PROPORTION` | tiket 01 | terbuka |
| OQ-TCO-05 | label "Underwriting Year"/"Transaction Year" bersilang | tiket 03 | terbuka |
| OQ-TCO-06 | pemilih kontrak memakai blacklist | tiket 02 | terbuka |
| OQ-TCO-07 | nama fisik `SOANote`/`Code` | tiket 02 | terbuka |
| OQ-TCO-08 | token penyimpanan (`GET_TOKEN_STORAGE` + garam) | tiket 12 | **ditutup** — *"sekarang"*; pelaksana nyata di balik `PELAKSANA_STORAGE=nyata`, bawaan `stub` (kelompok 5) |
| OQ-TCO-09 | pekerja latar `SatuPutaran` | tiket 12 | **ditutup** — *"perlu"*; `JalankanPekerja` dari `cmd/api`, interval env, bawaan mati (kelompok 5) |
| OQ-TCO-10 | anomali 366 hari tanggal akhir kontrak | tiket 04 | **ditutup** — *"mulai + 1 tahun kalender"*, penyimpangan sadar (kelompok 4) |
| OQ-TCO-11 | satu jenis satu kontrak per tahun | tiket 04 | **ditutup** — *"benar"* |
| OQ-TCO-12 | kolom `STATUSACTIVE` master `AGENT` | tiket 05 | **ditutup** — *"benar"* |
| OQ-TCO-13 | nilai nonaktif business | tiket 07 | **ditutup** — *"0 berarti nonaktif"* |
| OQ-TCO-14 | anti-dobel per jenis klausul, wajib Object/Periode | tiket 08 | **ditutup** — *"setuju"* |
| OQ-TCO-15 | daftar `ReinsTypeID` form klausul | tiket 08 | **ditutup** — *"benar"* |
| OQ-TCO-16 | kolom master `TREATYDESC`/`OCCUPATION`/`CLAUSE` | tiket 08 | **ditutup** — *"benar"* |
| OQ-TCO-17 | security dobel ditolak, `%Share` 0..100 | tiket 06 | **ditutup** — *"setuju"* (penyimpangan sadar dari Pega) |
| OQ-TCO-18 | `KURS` diisi, skala 8, dua kurs = master rusak | tiket 11 | **ditutup** — *"setuju"* |
| OQ-TCO-19 | tombol simpan tunggal | tiket 09 | **ditutup** — *"tidak perlu"*; rute simpan utuh dibuang (kelompok 3) |
| OQ-TCO-20 | "klausul milik kontrak ini" di popup hapus | tiket 10 | **ditutup** — *"dari induknya"* (kelompok 2) |
| OQ-TCO-21 | hapus kontrak yang kombinasinya dipakai bersama | tinjauan lanjutan 1 | **ditutup** — *"hapus saja, samain dengan pega"* (kelompok 2) |
| OQ-TCO-22 | `Folder` / `Durasi` / `Namafile` unggahan penyimpanan nyata | tiket 12 (lanjutan 2, kelompok 5) | terbuka — `[keputusan kami]` `TreatyContractOut` / `60` / `IMAGEID` |

## Keputusan work owner 29-09-2026

Jawaban atas OQ-TCO-08 … 21 (brief `PROMPT-LANJUTAN-TREATY-CONTRACT-OUT-2.md` §1). Label kode yang semula
`[keputusan kami]` / `[dugaan kuat]` untuk OQ yang ditutup kini `[keputusan work owner 29-09-2026]`.

OQ spec yang tetap terbuka (di luar register ini): aturan LimitMB & Portfolio (AC 36), arti `QUARTER = '0'`, arti bisnis
istilah klausul.
