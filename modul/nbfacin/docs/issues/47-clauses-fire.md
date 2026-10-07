# 47: Tab Clauses kasus FIRE

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 05-10-2026 "selanjutnya lanjut
> ke tab Clauses", disertai gambar layar Pega DEV (popup Choose Clause, dan daftar sesudah Submit). Frontend sesi 0f;
> backend sesi c3.

**What to build:** tab Clauses di Inward Facultative untuk kasus FIRE: daftar klausa kasus, popup Choose Clause
(multi-pilih), isi klausa + argumen, Save.

**Status:** frontend selesai 05-10-2026 (uji hijau). ⛔ Backend menunggu berkas dari work owner (lihat "Yang tidak ada di
korpus").

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\InputInwardFacultativeDtl.xml` tab "Clauses" (`pyContainerVisibleWhen = IsOilGas || IsFire`), sel 54 →
  `InputDtlClause_FacIn` (OTHER `IsFire`); Save sel 57 → `SaveFacIn_Act` (refresh `InputDtlObject_FacIn`).
- `Section\InputDtlClause_FacIn.xml`: S27 sel 209 "Choose Clause" (runActivity `SearchClauseFireSQL_PreAct`, harness
  popup `ChooseClauseFire`); S29 grid `.ClauseList` (tingkat kasus, `pyWorkPage.ClauseList`, kelas Data-Clause) → per baris
  `InputClauseFire_FacIn` + tombol hapus (deleteRow; Hapus / Batal ber-`SetFlagDelete_DT` hanya EDM).
- `Section\InputClauseFire_FacIn.xml`: 3 Clause Code (baca-saja) · 5 Total Argument (`.ArgumentCount > 0`) · 6 Title
  (baca-saja) · 8 Description (baca-saja) · 9 "See Clause and Argument" → local action `InputClauseFire_ViewDtl`.
- `FlowAction\InputClauseFire_ViewDtl` ("Isi Klasula"): pre `ViewClauseArgFireSQL` (argumen dari RDBList
  `SearchClauseArgFireSQL` atas M_ARGCLAUSEFIRE bila ArgumentList kosong), post `ReplaceClauseArgumentFireAct` (isi = isi
  asli, setiap `_&<nomor>` diganti Isi Argumen); section: `.ClauseContent` (rich text, baca-saja) + grid No. Argumen /
  Deskripsi Argumen / Isi Argumen (dapat diisi).
- `Section\ChooseClauseFire.xml`: 3 Language (`Clause.ClauseLanguageID`, dropdown) · 4 Key word · 7 Search · 8 Submit;
  grid `SearchClauseListOutput.pxResults` (kelas Int-CLAUSE): centang `.IsSelected`, Clause ID `.ID`, Title `.Title`,
  Description `.Info`.
- `Activity\SearchClauseFireSQL_PostAct.xml` (Data-Clause): per hasil tercentang, bila kodenya belum ada di ClauseList →
  `RetrieveArgumentNumberClauseSQL` (jumlah argumen) lalu tambah baris: ClauseCode=.ID, ClauseDescription=.Info,
  ClauseTitle=.Title, ClauseLanguage=Indonesia "0" / Inggris "1" / lain "2", ClauseLanguageID=.Language,
  ClauseContent=ClauseContentTemp=.Text, ArgumentCount; lalu Obj-Save.
- `SaveJsonOfferFacIn_Act`: `OfferFacIn.ClauseList = pyWorkPage.ClauseList`. Contoh kasus `DDL\P-5 *.txt`: ClauseList ada di
  akar, kosong di semua contoh.

## Yang tidak ada di korpus (diminta ke work owner 05-10-2026)

1. Activity `SearchClauseFireSQL_PreAct` kelas **Data-Clause** (atau RDB-List / RD yang mengisi `SearchClauseListOutput`) —
   SQL pencarian klausa. Di korpus hanya versi kelas Work yang memanggilnya.
2. DDL `M_CLAUSE`, `M_ARGCLAUSEFIRE` (+ view `V_CLAUSES` / `V_CLAUSES_ARG` bila ada). `DDL\` hanya memuat `M_KLAUSUL_PLAN`.
3. Aturan properti `ClauseLanguageID` (daftar Language).
4. (bila ada) `SearchClauseFireSQL_PostAct` kelas Work, RDB `ReplaceClauseFireSQL`.

## Frontend (sesi 0f)

- `components/TabClauses.tsx` (`tambahKlausa`, `gantiArgumen`, `kodeBahasa`, `adaArgumen`; popup Choose Clause, popup Isi
  Klasula) + uji; `pages/InwardFacultative.tsx` (tab Clauses kasus FIRE); `api.ts` (`cariKlausa`, `argumenKlausa`,
  `ambilKlausaKasus`, `simpanKlausaKasus`, tipe `KlausaKasus`, `HasilKlausa`, `ArgumenKlausa`); `labels.ts` (`KLAUSA`,
  `OPSI_BAHASA_KLAUSA`, `TEKS_KLAUSA`).

## Kontrak backend (untuk sesi c3; bagian sumber data menunggu berkas)

- `GET /api/nbfacin/klausa?bahasa=&q=&halaman=` → `{ baris: [{ id, title, info, text, language, argumentCount }], total,
  halaman, ukuran: 10 }` — mengikuti pencarian Pega (berkas 1).
- `GET /api/nbfacin/klausa/{id}/argumen` → `{ baris: [{ argumentNumber, argumentDescription, argumentValue }] }` —
  `SearchClauseArgFireSQL` (argumentValue = DefaultValue).
- `GET` / `PUT /api/nbfacin/kasus/{caseId}/klausa` → `{ baris: KlausaKasus[] }` (clauseCode, clauseTitle,
  clauseDescription, clauseLanguage, clauseLanguageId, clauseContent, clauseContentTemp, argumentCount, argumentList[]).
  Penyimpanan: tabel rancangan ClauseList (dicek c3 di skema loader).

## Backend (sesi c3, 05-10-2026)

Berkas 1–3 dikirim work owner 05-10-2026 ke `D:\migrasi\RNM\DDL\`: `SearchClauseFireSQL_PreAct.xml`, `RetrieveClauseSQL.xml`,
`M_CLAUSE.txt` (TABEL ID / OLDID VARCHAR2(7), JSONDATA), `CLAUSE.txt` (VIEW atas JSON M_CLAUSE; INFO FIRE = Description),
`ClauseLanguageID.xml` (PromptList "0" Indonesia, "1" Inggris, "2" Dual Bahasa) — `[terverifikasi]` dibaca c3.

- `GET /api/nbfacin/klausa?bahasa=&q=&halaman=` → `{ baris: [{ id, title, info, text, language, argumentCount }], total,
  halaman, ukuran: 10 }` — PreAct (`.Type = "FIRE"`, `.Language = ClauseLanguageID`, `.pyNote` = kata kunci) +
  RetrieveClauseSQL PERSIS: `title` / `text` = `decode(bahasa, '1', TITLEENG / TEXTENG, '2', TITLEDUAL / TEXTDUAL, '0',
  TITLEINA / TEXTINA)` atas view CLAUSE ber-`TYPE = 'FIRE'`, judul tidak kosong, `"Info" like '%' || q || '%'` (peka huruf,
  tanpa ESCAPE, seperti asal). `bahasa` kode bawaan "0"; selain 0 / 1 / 2 → 400. `language` = LABEL ClauseLanguageID kode
  itu ("Indonesia" / "Inggris" / "Dual Bahasa") — `.Language` yang PostAct petakan ke ClauseLanguage ("Indonesia" "0",
  "Inggris" "1", lain "2"). `argumentCount` = RetrieveArgumentNumberClauseSQL (yang PostAct pakai untuk ArgumentCount,
  `HASIL1`), BUKAN `SUMOFARGUMENT`.
Dibangun lebih dulu (tidak bergantung berkas 1–2):

- `GET /api/nbfacin/kasus/{caseId}/klausa` → `{ baris: KlausaKasus[] }` urut simpan; `argumentList` selalu larik. 404 case
  tidak ada / bukan LINI Fac In, 503 tanpa basis data.
- `PUT /api/nbfacin/kasus/{caseId}/klausa` badan `{ baris: KlausaKasus[] }` (kunci asing → 400) → ganti utuh dalam satu
  transaksi (sentuh case, pastikan General, hapus argumen lalu klausa, sisip ulang urut), jawab baca ulang. Identitas
  wajib (401).
- `GET /api/nbfacin/klausa/{id}/argumen` → `{ baris: [{ argumentNumber, argumentDescription, argumentValue }] }` —
  `RDBList\SearchClauseArgFireSQL.xml` PERSIS (`[terverifikasi]` 05-10-2026): JSON dot-notation `ArgumentNumber` /
  `Description` / `DefaultValue` dari `M_ARGCLAUSEFIRE` ber-`OLDID` = `OLDID` `M_CLAUSE` ber-ID itu, urut JSON
  `ArgumentNumber` apa adanya (teks: "10" sebelum "2", seperti SQL asal). Struktur M_ARGCLAUSEFIRE / M_CLAUSE yang
  dipakai (ID, OLDID, JSONDATA) cukup dari SQL itu; DDL tetap diminta.
- Migrasi `197_t_clauselist.sql` (+ `_down`): SATU tabel `T_CLAUSELIST` (induk `T_GENERAL_POLIS`), klausa beserta
  argumennya (K47-4) + sequence; `docs/STRUKTUR-TABEL-NB-FACIN.md`. Ditulis, belum dijalankan. Kontrak JSON GET / PUT
  tidak berubah (`argumentList` larik per klausa).
- Berkas: `models/klausa.go`, `repository/klausa.go` (+ uji SQL), `services/klausa.go`, `services/layanan.go`,
  `handlers/klausa.go` (+ uji), `handlers/handlers.go` (rute, peta galat 400 / 501 / 503).

## Keputusan agent (menunggu konfirmasi)

- **K47-1** (c3) ClauseList TIDAK ada di workbook rancangan (`loader/skema_gen.go`) dan kosong di seluruh 108 contoh
  `DDL\CONTOH` → tabel rancangan agent berpola tabel rancangan (ID dari sequence, IDPEGA, COB_GROUP, PARENT_ID,
  SEQ_NO, ROW_UID), induk `T_GENERAL_POLIS` seperti daftar akar lain (LocationList). Nomor 197 (196 dilewati).
- **K47-2** (c3) Lebar teks VARCHAR2(4000) — DDL M_CLAUSE / M_ARGCLAUSEFIRE belum ada, nilai dot-notation = VARCHAR2(4000).
  Isi klausa juga VARCHAR2(4000): kolom dokumen besar (CLOB) dilarang penjaga migrasi `TestKolomUangDesimalDanNolJSON`
  (keputusan tim inti), dan isi asal lewat dot-notation maksimal 4000 bita. Isi sesudah argumen diganti yang melebihi
  4000 bita **ditolak 400**, tidak dipotong. `ARGUMENT_COUNT` teks angka (frontend teks).
- **K47-4** (c3) **Satu tabel** (work owner 05-10-2026: "T_CLAUSELIST dan T_CLAUSEARGUMENTLIST bukannya bisa jadi 1 tabel
  aja ? T_CLAUSEARGUMENTLIST disimpan di T_CLAUSELIST juga"): **satu baris per argumen**, kolom klausa berulang,
  `ARGUMENT_SEQ_NO` 1..n; klausa tanpa argumen = satu baris ber-`ARGUMENT_SEQ_NO` kosong. Kelebihan: satu tabel, tanpa
  FK anak, simpan / baca satu pernyataan, lolos penjaga. Kekurangan: kolom klausa (isi sampai 4000 bita) tersalin per
  argumen; satu klausa = beberapa baris (identitas klausa = `PARENT_ID` + `SEQ_NO`, bukan `ID`); kolom klausa harus
  ditulis sama di setiap barisnya (aplikasi menulis ulang utuh tiap Save, jadi tidak menyimpang). DITOLAK: daftar
  argumen sebagai teks terstruktur di satu kolom — penjaga `TestKolomUangDesimalDanNolJSON` melarang atribut di dalam
  dokumen; pengecualiannya keputusan tim inti. Bila 197 **sudah** dijalankan dengan bentuk dua tabel: migrasi baru
  (198) — `ALTER TABLE T_CLAUSELIST ADD (ARGUMENT_SEQ_NO, ARGUMENT_NUMBER, ARGUMENT_DESCRIPTION, ARGUMENT_VALUE)`, salin
  argumen (baris klausa pertama mendapat argumen 1, argumen berikut disisip sebagai baris baru berkolom klausa sama),
  lalu `DROP TABLE T_CLAUSEARGUMENTLIST` + sequence-nya — dan 197 harus dikembalikan ke bentuk lama (nama tercatat di
  `T_MIGRASI`).
- **K47-5** (c3) Pencarian: urut ID (SQL asal tanpa urutan) untuk halaman 10 yang stabil; LIKE peka huruf dan tanpa
  ESCAPE dibiarkan seperti asal (`%` / `_` di kata kunci bertindak sebagai pola). ⚠️ `.Language` hasil PostAct
  `[dugaan]`: RetrieveClauseSQL tidak mengembalikan Language, sedangkan PostAct membandingkan `.Language` dengan
  "Indonesia" / "Inggris" — dibaca sebagai label bahasa yang dicari.
- **K47-7** (c3, 05-10-2026) **Permintaan pengguna — menyimpang dari RetrieveClauseSQL**: (1) "ini keyword buatin bisa
  huruf besar atau kecil dong" → kata kunci dicocokkan `UPPER("Info") LIKE '%' || UPPER(q) || '%'` (asal peka huruf;
  K47-5 diganti untuk bagian ini); (2) "kalau pilihnya bahasa inggris, ... ambil yg bahasa inggris juga dong, kn di
  tabel nya ada, kolom TEXTENG, kalau kolom itu kosong baru ambil yg bahasa indonesia ( kolom TEXTINA)" → isi bahasa "1"
  = `TEXTENG`, atau `TEXTINA` bila `TEXTENG` kosong / spasi saja. Judul dan bahasa "0" / "2" tetap seperti asal. Isi
  klausa yang SUDAH tersimpan di kasus tidak berubah (diambil saat Choose Clause).
- **K47-6** (c3, 05-10-2026) **M_ARGCLAUSEFIRE tidak ada di DEV** — bukti query work owner: `ALL_OBJECTS` `OWNER =
  'POOLDATA'` untuk CLAUSE / M_CLAUSE / M_ARGCLAUSEFIRE hanya mengembalikan CLAUSE (VIEW) dan M_CLAUSE (TABLE); M_CLAUSE
  tanpa ID ganda. Pencarian di DEV menjawab 500 — `[dugaan kuat]` ORA-00942 dari subkueri jumlah argumen (baris log belum
  dikirim). Kini keberadaan M_ARGCLAUSEFIRE ditanya ke katalog sekali per permintaan: tidak ada → pencarian tetap jalan
  dengan `argumentCount` KOSONG (bukan "0") dan satu peringatan log yang menyebut tabelnya; `GET …/klausa/{id}/argumen`
  → **200 `{"baris":[]}`** + peringatan log. Pengguna 05-10-2026 (gambar popup Isi Klasula, klausa `1000340`):
  "M_ARGCLAUSEFIRE memang tidak ada, saat di klik See Clause and Argument, yg muncul itu adalah TEXTINA atau TEXTENG" —
  popup hanya menampilkan isi klausa (Text menurut bahasa), tanpa grid argumen; sempat 503 sebelum keterangan ini. Galat
  lain tidak ditangkap (tetap 500 + log "nbfacin: galat server").
- **K47-3** (c3) Pemeriksaan Save: `clauseCode` wajib dan tidak ganda (PostAct hanya menambah kode yang belum ada);
  `clauseLanguage` kosong atau "0" / "1" / "2"; `argumentCount` kosong atau angka; paling banyak 500 klausa.

- **C-1** Isi klausa ditampilkan sebagai teks (bukan HTML rich text) — aman dari skrip.
- **C-2** Penggantian argumen mengikuti urutan Pega (nomor menaik, ganti-semua harfiah; `_&1` juga mengenai awal `_&10`).
- **C-3** Hasil Choose Clause berhalaman 10 baris.
- ~~**C-4** Language sementara hanya "Indonesia"~~ — aturan `ClauseLanguageID` dikirim 05-10-2026 (`DDL\ClauseLanguageID.xml`):
  kode 0 Indonesia / 1 Inggris / 2 Dual Bahasa; awal 0; yang dikirim ke pencarian = kode.
- **C-5** Hapus baris = langsung dari daftar (deleteRow); varian Hapus / Batal ber-FlagDelete khusus EDM tidak dibawa.

## Acceptance criteria

- [x] Tab Clauses tampil untuk kasus FIRE; daftar per baris sesuai `InputClauseFire_FacIn`.
- [x] Choose Clause multi-pilih, tanpa ganda; See Clause and Argument mengisi argumen dan menerapkannya ke isi klausa.
- [x] Backend: argumen klausa, baca / simpan ClauseList kasus (sesi c3, 05-10-2026; migrasi 197 belum dijalankan).
- [x] Backend: cari klausa (sesi c3, 05-10-2026; PreAct + RetrieveClauseSQL).

### Popup Isi Klasula tanpa argumen (sesi 0f, 05-10-2026)

Keterangan work owner (diteruskan sesi c3): "M_ARGCLAUSEFIRE memang tidak ada, saat di klik See Clause and Argument, yg
muncul itu adalah TEXTINA atau TEXTENG" — gambar popup "Isi Klasula" hanya teks klausa, tombol Cancel dan Submit.
Frontend: argumen kosong → popup hanya isi klausa (tanpa grid, tanpa catatan); tombol Cancel / Submit (label Pega); isi
tidak berubah saat Submit tanpa argumen.
