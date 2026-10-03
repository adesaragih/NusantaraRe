# Register OQ — Treaty Contract Out

Satu tempat untuk seluruh pertanyaan terbuka modul ini. Rincian dan buktinya tetap di tiket terkait dan di bab
`LAPORAN-GILIRAN.md` tempat OQ itu dibuka; di sini status dan keputusannya.

| OQ | Pokok | Dibuka | Status |
| --- | --- | --- | --- |
| OQ-TCO-01 | bentuk teks tanggal warisan (tujuh bentuk dikenal pengurai) **+ bentuk TULIS** `STARTDATE`/`ENDDATE` tahun & reinsurer (tco4: stempel Pega 00:00 WIB, `[dugaan kuat]`) | tiket 01, lanjutan 3 | **ditutup** (lanjutan 4, dari data DEV) — `TREATYYEAR.STARTDATE/ENDDATE` `YYYYMMDD` (182/182), pembaca menolak bentuk lain; `TREATYREINSURER.STARTDATE/ENDDATE` tidak ditulis (430/430 kosong) — `708a351` |
| OQ-TCO-02 | arti/bentuk `IUDATE` | tiket 01 | terbuka |
| OQ-TCO-03 | isi hidup `PROPORTIONALLIST`/`OBJECT` | tiket 01 | terbuka |
| OQ-TCO-04 | arti `PROPORTION` | tiket 01 | **ditutup 30-09-2026** (work owner) — dua nilai `Proportional` / `NonProportional` (label "Proportional" / "Non Proportional"), bukan master jenis reasuransi; daftar `associated` properti `.Proportion` tidak diekspor, selaras param `SetTreatyArrangementDesc_Act` dan `ViewDetailDescription.xml` b3614/b6719 (tiket 03) |
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
| OQ-TCO-18 | `KURS` diisi, skala 8, dua kurs = master rusak | tiket 11 | **ditutup** — *"setuju"*; **dipersempit 29-09-2026** — baris KEMBAR (TOIDR sama) = satu kurs, TOIDR berbeda tetap 503 (tiket 11) |
| OQ-TCO-19 | tombol simpan tunggal | tiket 09 | **ditutup** — *"tidak perlu"*; rute simpan utuh dibuang (kelompok 3) |
| OQ-TCO-20 | "klausul milik kontrak ini" di popup hapus | tiket 10 | **ditutup** — *"dari induknya"* (kelompok 2); angkanya tidak lagi tampil di popup (01-10-2026) |
| OQ-TCO-21 | hapus kontrak yang kombinasinya dipakai bersama | tinjauan lanjutan 1 | **ditutup** — *"hapus saja, samain dengan pega"* (kelompok 2) |
| OQ-TCO-22 | `Folder` / `Durasi` / `Namafile` unggahan penyimpanan nyata | tiket 12 (lanjutan 2, kelompok 5) | terbuka — untuk work owner; nilai `TreatyContractOut/` / `60` / `IMAGEID` dipertahankan berlabel `[terbuka — OQ-TCO-22]` |
| OQ-TCO-23 | bentuk TULIS desimal teks `PROPORTIONALARRG.RP/USD/PCT/PCTME/KURS`, `MTREATYSECURITY.PCT_SHARE` (tco4: titik, tanpa ribuan, `[dugaan kuat]` hasil `@toDecimal`) | lanjutan 3 | **ditutup** (lanjutan 4, dari data DEV) — titik: `RP` 930 bertitik / 0 berkoma, `PCT` 520 / 0; kode sudah menulis titik |
| OQ-TCO-24 | badan `PEGA_M_ATTACHMENT` (penulis lampiran Treaty Contract Out, `InsertAtatchment_Sql` b60), tipe kolom dan PK/indeks `M_ATTACHMENTTREATY_2` / `T_STORAGE_IMAGE` (daftar lampiran membaca `T_STORAGE_IMAGE` per `IMAGEID`) | lanjutan 3 | **ditutup** (lanjutan 5, dari `ALL_SOURCE`) — prosedur meng-upsert simpanan JSON halaman `M_ATTACHMENTTREATY.DATA_JSON` per `ID`; `id_count` variabel lokal, bukan tabel; `_2` tabel relasional yang dibaca/dihapus seluruh RDB modul. Keputusan `[asisten dari bukti; veto work owner]`: tetap menulis dan membaca `_2`, simpanan JSON **tidak** diteruskan (spec §15). Temuan warisan: penulis hidup Pega tidak menulis `_2` — `dba-procedures.md` |
| OQ-TCO-25 | `USERID`/`TGLUPDATE` reinsurer dan business diisi layanan walau Pega mengosongkannya (jejak modul dibuang tco4) | lanjutan 3 | **ditutup** (lanjutan 4) — kosongkan seperti Pega (data DEV 0/430, 2/4.621); pelaku di log aplikasi — `8426880` |
| OQ-TCO-26 | `Update_T_Storage_SQL` (menyegarkan `URLPUBLIC`/`EXPDATE`/`TANGGAL_UPLOAD` sesudah `geturl`, `GetUrlGoogleStorage_Act`) belum ditiru — unduhan kita memakai URL bertanda tangan segar tiap kali; pembaca lain tabel bersama melihat URL saat unggah | /code-review lanjutan 3 | **ditutup** (lanjutan 4) — ditiru: sesudah tiap geturl yang berhasil, transaksi pendek sendiri; `exp` diubah seperti Pega, juga di jalur unggah — `8426880` |
| OQ-TCO-27 | kolom uang teks warisan (`PROPORTIONALARRG.RP`/`USD`, dll.) yang memuat ribuan bertitik SATU (`500.000`) — desimal atau ribuan? Sejak 29-09-2026 bertitik ≥ 2 (`1.000.000`) dibaca ribuan (data DEV Limit MB); bertitik satu tetap desimal dan akan terbaca 1000× lebih kecil TANPA galat bila sebenarnya ribuan. Sensus baca-saja yang menjawabnya: `SELECT COUNT(*) FROM POOLDATA.PROPORTIONALARRG WHERE REGEXP_LIKE(RP,'^[0-9]{1,3}\.[0-9]{3}$') OR REGEXP_LIKE(USD,'^[0-9]{1,3}\.[0-9]{3}$')` | tiket 08 (penyisiran layar 29-09-2026) | **terbuka** — untuk work owner / DBA |

## Keputusan work owner 29-09-2026

Jawaban atas OQ-TCO-08 … 21 (brief `PROMPT-LANJUTAN-TREATY-CONTRACT-OUT-2.md` §1). Label kode yang semula
`[keputusan kami]` / `[dugaan kuat]` untuk OQ yang ditutup kini `[keputusan work owner 29-09-2026]`.

OQ spec yang tetap terbuka (di luar register ini): aturan Portfolio (AC 36; LimitMB ditutup 02-10-2026 — ikuti XML), arti `QUARTER = '0'`, arti bisnis
istilah klausul.
