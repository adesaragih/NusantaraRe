# Modul `riratelife` — R/I Rate Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner 05-10-2026:
*"Buat Menu baru Namanya R/I Rate Life pada Master Treaty, menu ini bisa CRUD untuk simpan data ke tabel
RATE_LIFE_SUMMARY, panduannya xml yang saya berikan"*. Panduan: section Pega `InboxSummaryRIRate` (kelas
`ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY`, judul "R/I RATE SUMMARY" b382). Baris menunya dibuat migrasi inti `922`
(golongan MASTER TREATY, URUTAN 9), sequence ID-nya migrasi inti `923`.

⛔ **Tanpa migrasi sendiri** (`—` di bawah, `tandaTanpaMigrasi` di `inti/backend/penjaga`): seluruh nomor modul
001-899 dan slot menu 950-999 sudah terbagi, dan modul lain tidak boleh disentuh.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `riratelife` |
| Folder korpus | `R/I Rate Life` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-RIRATELIFE` |
| Status | dimigrasi |
| Rentang migrasi | `—` |
| Slot menu | `—` |
| Prefix rute API | `/api/ri-rate-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Keputusan work owner 05-10-2026 (sesudah ditanya)

- **K1 tabel tujuan:** `RATE_LIFE_SUMMARY` di DEV adalah VIEW (6 kolom, dibaca `mastercontractretrolife` dan
  `masterproductnamelife`) atas tabel fisik `M_RATE_LIFE_SUMMARY` (339 baris). *"CRUD menulis ke
  `M_RATE_LIFE_SUMMARY`, view TIDAK di-DROP, nol DDL pada tabel/view warisan."* View `RATE_LIFE_SUMMARY` dipakai
  membaca grid. ⚠️ **Digantikan RALAT R4 (06-10-2026)** - lihat bab di bawah.
- **K2 ID baru:** dari **sequence Oracle** (`SEQ_M_RATE_LIFE_SUMMARY`, `SEQ_M_RATE_LIFE`, migrasi inti 923).
- **K3 lingkup:** termasuk **Detail** (rincian rate) dan **Upload CSV / View Upload / Simpan Upload**. XML rule-rule
  itu tidak tersedia → dirancang dari petunjuk format di XML (`Format excel : USEDBY, CONTRACT, GENDER, AGE, RATE`
  b5305); perilaku di luar itu = ASUMSI (tabel di bawah).

## Keputusan work owner 07-10-2026 — ringkasan SATU tabel `M_RATE_LIFE_SUMMARY` (RALAT R6, menggantikan R4)

*"Ringkasan R/I Rate Life cukup SATU tabel - M_RATE_LIFE_SUMMARY."* Sasaran akhir POOLDATA: `M_RATE_LIFE_SUMMARY`
berkolom persis seperti tabel flat `RATE_LIFE_SUMMARY` (926): `ID VARCHAR2(10)` PK, `USEDBY VARCHAR2(500)`,
`TYPE VARCHAR2(100)`, `MODIFIEDDATE VARCHAR2(50)`, `OPERATORID VARCHAR2(200)`; kolom `JSONDATA` dan constraint
`ENSURE_M_RATE_LIFE_SUMMARY_JSON` dibuang; tabel flat `RATE_LIFE_SUMMARY` (926) DIHAPUS; nol tabel baru.

Keadaan DEV (WO 07-10-2026, baca-saja): 926 sudah jalan - `RATE_LIFE_SUMMARY` = TABLE, SUMBER KEBENARAN (aplikasi
menulis ke sana sejak 09:09), identik dengan JSON `M_RATE_LIFE_SUMMARY` (MINUS dua arah = 0); `M_RATE_LIFE_SUMMARY` = ID
PK + JSONDATA CLOB (IS JSON). Dependensi DB `M_RATE_LIFE_SUMMARY`: prosedur `PEGA_M_RATE_LIFE_SUMMARY` dan
`PEGA_M_PLAN_LIFE_SUMMARY` - akan INVALID. ⚠️ **Pemutusan Pega:** WO menerima R/I Rate TIDAK lagi disimpan lewat Pega;
satu-satunya penulis ringkasan = modul ini.

[keputusan work owner 07-10-2026] DELETE 928 langkah 2 (buang baris `M_RATE_LIFE_SUMMARY` yang tidak ada di tabel flat) DISETUJUI; tabel flat tetap sumber kebenaran. Bukti POOLDATA (baca-saja, WO 07-10-2026): JSON 347 baris, flat 347 baris; ID hanya di JSON = `1000469` ("TEST RATE LIFE", sudah dihapus pengguna lewat aplikasi); ID hanya di flat = `1000471` (baru dari aplikasi); isi beda pada ID yang sama = 0; tidak ada tulisan Pega sesudah cutover 09:09. Pengecualian penjaga `TestKolomUangDesimalDanNolJSON` untuk `DROP COLUMN JSONDATA` 928 DISETUJUI WO, tetap ditinjau tim inti (`docs/PR-RIRATELIFE-FLAT.md`).

Penerapan:

- Migrasi inti `927_m_rate_life_summary_kolom.sql` (ALTER … ADD keempat kolom, berdiri sendiri - ADD tidak aman
  diulang) dan `928_m_rate_life_summary_satu_tabel.sql` (isi dari flat, buang baris yang sudah dihapus aplikasi, buang
  JSONDATA + constraint IS JSON lewat blok berpelindung katalog, sisip baris flat-saja, indeks
  `IX_M_RATE_LIFE_SUMMARY_NAMA`, DROP TABLE flat TERAKHIR - setiap pernyataan aman diulang). Jalur mundur keduanya
  mengembalikan keadaan sesudah 926 (FLAG lama hanya dari cadangan CSV).
- Repository: ringkasan dibaca DAN ditulis di kolom `M_RATE_LIFE_SUMMARY` (satu tabel; ID terpakai diperiksa di sini
  saja). Alat `pindahflat` dan berkas `DBA-LEPAS-VIEW-RATE_LIFE_SUMMARY.sql` DIHAPUS (tidak diperlukan lagi; riwayatnya
  di git). `db.Koneksi` (inti) tetap - dipakai ricommlife.
- Pembaca lain: `masterproductnamelife` (`Choose R/I Rate`) dan `mastercontractretrolife` (autocomplete `R/I RATE`)
  kini membaca `SELECT ID, USEDBY FROM M_RATE_LIFE_SUMMARY` (diizinkan WO), tetap baca-saja.
- **Urutan WO siap-tempel + pemulihan + jalur mundur: [`docs/LANGKAH-WO-RIRATELIFE-SATU-TABEL.md`](docs/LANGKAH-WO-RIRATELIFE-SATU-TABEL.md).**

Riwayat: RALAT R4 (06-10-2026, view → tabel flat `RATE_LIFE_SUMMARY`, migrasi 926, alat pindah delta) dan RALAT R5
(FLAG tidak digunakan) - tabel RALAT di bawah. FLAG tetap tidak digunakan; nilai FLAG lama hanya ada di cadangan CSV
`ID + JSONDATA` yang dibuat sebelum 928 (LANGKAH-WO (a)).

### Urutan merge

**`modul/riratelife/implementasi` di-merge LEBIH DULU, lalu `modul/ricommlife/implementasi`.** Kedua cabang membawa
`inti/backend/db/koneksi.go` identik (commit `91221b2f` di sini = `349be34d` di ricommlife); migrasi inti 926-928 di
sini dan 924/925 di ricommlife. PR: `docs/PR-RIRATELIFE-FLAT.md`.

**Uji coba merge (07-10-2026)** - worktree sementara dari ujung riratelife, `git merge --no-commit
modul/ricommlife/implementasi`: `inti/backend/db/koneksi.go` bersih (identik). SATU konflik, di
`inti/backend/penjaga/rentang_test.go` (`TestSlotMenuBerjalanSesudah900`, daftar urutan pelari). Penyelesaian: satu
daftar berurutan `… "923_seq_rate_life.sql", "924_ricomm_life.sql", "925_m_nav_menu_ricommlife.sql",
"926_rate_life_summary_flat.sql", "927_m_rate_life_summary_kolom.sql", "928_m_rate_life_summary_satu_tabel.sql",
"952_menu_tiruan.sql"`. Sesudahnya `go vet ./...` bersih dan `go test ./...` = baseline (3 paket claimlife), ditambah
`TestNolAlamatLayananDiKode` yang gagal HANYA karena worktree sementara tidak memuat berkas `.env` (tidak masuk git).

## Keputusan work owner 07-10-2026 — rincian rate tabel flat `M_RATE_LIFE` (RALAT R7)

*"Detail R/I Rate Life (M_RATE_LIFE) menjadi tabel flat, satu tabel saja (pola sama dengan ringkasan, migrasi
927/928). M_RATE_LIFE TIDAK dihapus."* Sasaran akhir POOLDATA: `M_RATE_LIFE` berkolom PERSIS seperti view `RATE_LIFE`
lama - `ID` (PK warisan), `IDUSEDBY`, `USEDBY`, `TYPE`, `GENDER`, `CONTRACT`, `AGE`, `RATE`, semua teks (RATE tetap
teks: nilai Pega apa adanya, berkoma maupun bertitik desimal); `JSONDATA` + constraint `ENSURE_M_RATE_LIFE_JSON`
dibuang; view `RATE_LIFE` DIBUANG; indeks `IX_M_RATE_LIFE_IDUSEDBY`; nol tabel baru.

Keadaan DEV (WO 07-10-2026, baca-saja): `M_RATE_LIFE` = `ID VARCHAR2(10)` PK + `JSONDATA CLOB` (IS JSON), 98.306 baris;
view `RATE_LIFE` 8 kolom VARCHAR2; panjang maksimum isi ID 7, IDUSEDBY 7, USEDBY 95, TYPE 0 (kosong), GENDER 1,
CONTRACT 3, AGE 3, RATE 20 byte; RATE 90.437 baris berkoma desimal, 3.226 bertitik desimal; 539 baris yatim (IDUSEDBY
tanpa ringkasan); `SEQ_M_RATE_LIFE` terakhir 1117016. Kunci JSON `FLAG` TIDAK dipindah (RALAT R5) - hanya di cadangan CSV.
Prosedur `PEGA_M_RATE_LIFE` menjadi INVALID (diterima WO).

Penerapan:

- Migrasi inti `929_m_rate_life_kolom.sql` (ALTER … ADD ketujuh kolom, berdiri sendiri) dan
  `930_m_rate_life_satu_tabel.sql` (blok berpelindung katalog: isi kolom dari JSONDATA apa adanya, lalu
  `DROP COLUMN JSONDATA CASCADE CONSTRAINTS`; CREATE INDEX; DROP VIEW `RATE_LIFE` terakhir). Jalur mundur 930 + 929
  membangun ulang JSONDATA (`JSON_OBJECT … ABSENT ON NULL`), constraint, dan view (FLAG hanya dari cadangan CSV).
- Modul ini membaca dan menulis kolom bernama `M_RATE_LIFE` (`backend/repository/rirl.go`, `rirl_tabel.go`
  `KolomRate`); `rirl_json.go` (baca `FOR UPDATE` + ganti kunci di Go) dihapus. Perilaku tetap: kosong = NULL (dulu
  kunci dibuang), RATE disimpan berkoma desimal (A8), `TYPE` tidak ditulis.
- Pembaca lain ikut membaca `M_RATE_LIFE` berkolom sama (baca-saja, nama konstanta tetap): `claimlife`
  (`NamaViewRateLife`, spreading retro), `premiumlistlife` (`ViewRateProduk`, rate produk / Hitung QR),
  `masterproductnamelife` (`MasterRate`, View Rate), `mastercontractretrolife` (`MasterRate`, Rate List) - ⚠️ tinjauan
  tim inti (`docs/PR-RIRATELIFE-FLAT.md`).
- Langkah WO: `docs/LANGKAH-WO-RIRATELIFE-DETAIL-FLAT.md` + berkas SQL\*Plus `docs/sql/detail_flat_*.sql` (prasyarat:
  Pega R/I Rate dan backend dihentikan, nol kunci DML, angka acuan dibandingkan lagi tepat sebelum `-migrate`).
- Tinjauan independen 08-10-2026: `929_down` / `927_down` mulai dengan pelindung gagal-keras
  (`UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0` → ORA-00904 bila JSONDATA sudah dibuang); `930_down` aman diulang;
  JSONDATA dikembalikan NULLABLE (bukti NOT NULL tidak ada di repo).

### Item terbuka work owner — 539 baris yatim

539 baris `M_RATE_LIFE` (WO 07-10-2026) ber-IDUSEDBY yang tidak ada di `M_RATE_LIFE_SUMMARY`. Keputusan: TIDAK dihapus,
dipindah 930 apa adanya. Tidak tampil di layar mana pun modul ini (Rate Detail dibuka dari ringkasan), tetapi tetap
terbaca pembaca lain menurut IDUSEDBY. Contoh ID BELUM dilihat executor (tidak dikarang); kueri baca-saja untuk WO:

```sql
SELECT COUNT(*) FROM POOLDATA.M_RATE_LIFE r WHERE NOT EXISTS
  (SELECT 1 FROM POOLDATA.M_RATE_LIFE_SUMMARY s WHERE s.ID = r.IDUSEDBY);
SELECT r.ID, r.IDUSEDBY, r.USEDBY FROM POOLDATA.M_RATE_LIFE r WHERE NOT EXISTS (SELECT 1 FROM POOLDATA.M_RATE_LIFE_SUMMARY s WHERE s.ID = r.IDUSEDBY) ORDER BY r.IDUSEDBY, r.ID FETCH FIRST 20 ROWS ONLY;
```

Perlu keputusan WO: dibiarkan, dihapus, atau disambungkan ke ringkasan - di luar cakupan R7.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-RIRATELIFE.md` — kolom ringkasan, tabel/view warisan, kunci JSON; `LANGKAH-WO-RIRATELIFE-SATU-TABEL.md`; `LANGKAH-WO-RIRATELIFE-DETAIL-FLAT.md` (+ `sql/detail_flat_*.sql`); `PR-RIRATELIFE-FLAT.md` |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `modul.go` (tanpa `migrations/`) |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `riratelife.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Rule XML → kode

| XML | Perilaku | Kode |
| --- | --- | --- |
| `.ID` b1089, pxTextInput disabled selalu b1104-b1105 | ID tampil, tidak dapat diisi; baru = sequence | `FormRingkasan.tsx`, `services.pemberiID` |
| `.USEDBY` b1273, label "R/I RATE NAME" b1242, wajib b1285 | wajib, dipangkas, **tidak kembar** tanpa beda huruf | `services.Simpan` |
| `InputParam.ERRMSG` b896 | pesan galat di atas form | `Gagal` di `FormRingkasan.tsx` |
| Save b1809 → `AddToListSummary_Act` b1833; Cancel b2109 (visible `DATASHOW='IsEdit'` b2275) → `NewDataSummary_DT` b2137 | container `1=2` b1529 disembunyikan Pega; **ditampilkan** (perintah CRUD). Cancel hanya saat Edit = kembali ke form Add | `POST` / `PUT /api/ri-rate-life` |
| Grid `BrowseRateLifeSummary` b9568, param `id` / `idusedby` b9308-b9314, 50 baris b12645, read-only b12575 | filter ID / R/I RATE NAME, urut per kolom, paging 50 | `GET /api/ri-rate-life` |
| Kolom ID · R/I RATE NAME b9588 · MODIFY OPERATOR b9748 (`.OPERATORID` b10709) · MODIFY DATE b9894 (`.MODIFIEDDATE` b10879, `Date-Short-Custom-YYYY` b10929) | tanggal tampil `DD-MM-YYYY` WIB | `TampilTanggal` |
| Edit b11109 → `EditListSummary_DT` b11136 (ID, UsedBy) | isi form dari baris | `RIRateLife.tsx` |
| Detail b11398 → `setIDUsedBy_Act` b11415 + harness `InboxRIRate` b11444, "Rate Detail" b11446 | popup section `InboxRIRate` (baris di bawah, View Detail.xml) | `RateDetail.tsx` |
| `InboxRIRate` judul "R/I RATE DETAIL" b367; `.ID` b1342 disabled; R/I RATE NAME = `TempIDUsedBy.USEDBY` b1747 (disabled, wajib b1766); IDUSEDBY = `TempIDUsedBy.ID` b1526 (`1=2` b1657) | nama dan IDUSEDBY **dari ringkasan yang dilihat**, bukan isian; ID baru = sequence | `services.SimpanRate` |
| `.GENDER` b1938 radio (tidak wajib) · `.CONTRACT` b2122 wajib b2141, placeholder 0 · `.AGE` b2400 tidak wajib · `.RATE` b2565 wajib b2579, placeholder `0,0000` b2583; `ChangeDotToPoint_DT` b2171/b2609 | GENDER U/M/F atau kosong; CONTRACT dan AGE bulat 0-120; RATE disimpan berkoma desimal | `models.PeriksaIsianRate` |
| Save b3118 → `AddToList_Act` b3196 (visible ALWAYS b3336); Cancel b3418 → `NewData_DT` b3497 (visible `DATASHOW='IsEdit'` b3638); Clear Field b1064 → `clearInputFieldRIRate_act` b1142 | tambah / ubah satu baris; Cancel hanya saat Edit; Clear Field mengosongkan form | `POST /{id}/rate`, `PUT /{id}/rate/{rateId}` |
| Grid `BrowseRateLife_RD` b10200, idusedby b7568: ID · R/I RATE NAME · GENDER · CONTRACT · AGE · RATE; urut ID DESC b7668 lalu RATE ASC b8436; 20 per halaman b10206 | ID menurun (angka); paging 20 | `GET /{id}/rate`, `SqlDaftarRate` |
| Edit b9773 → `EditList_DT` b9853 (ID, UsedBy, Gender, Contract, Age, Rate) | isi form dari baris; simpan = ganti GENDER/CONTRACT/AGE/RATE di JSONDATA (RALAT R3) | `repository.UbahRate` |
| Upload CSV / View Upload / Simpan Upload di `InboxRIRate` (container `1=2` b4029) | tidak ditampilkan di popup (sudah ada di halaman ringkasan) | — |
| Delete b12238 → `DeleteSummaryDetail` b12262, DeleteID=.ID b12279 | hapus ringkasan **beserta** rate ber-IDUSEDBY sama, satu transaksi, konfirmasi menyebut jumlah rate | `DELETE /{id}`, `services.Hapus` |
| "Upload CSV" b2949 → local action `UploadCSV_RIRATE` b3019 | pilih berkas CSV (popup) | `UnggahCSV.tsx` |
| "View Upload" b3467 → popup `ViewCSVResult_RIRate` b3486 | pratinjau hasil urai + galat per baris, tanpa menulis | `POST /unggah/pratinjau` |
| "Simpan Upload" b4511 → `SubmitRIRate_Act` b4535 | urai ulang di server, tolak seluruhnya bila ada galat (422 per baris), satu transaksi | `POST /unggah`, `services.SimpanUnggah` |
| Label `Format excel : USEDBY, CONTRACT, GENDER, AGE, RATE` b5305 | kepala CSV wajib lima kolom itu | `models.UraiCSV` |

## Asumsi terbuka

| # | Asumsi | Alasan / risiko |
| --- | --- | --- |
| ~~A1~~ | ~~kunci JSON `M_RATE_LIFE_SUMMARY` dianggap `USEDBY`, `OPERATORID`, `MODIFIEDDATE`~~ | ✅ **terbukti — RALAT R1** di bawah |
| ~~A2~~ | ~~Update = `JSON_MERGEPATCH(... RETURNING CLOB)`~~ | ❌ **gugur — RALAT R3**: Oracle DEV menolaknya; ubah = baca-ubah-tulis JSONDATA di Go |
| A3 | MODIFIEDDATE ditulis format Pega `YYYYMMDDTHHMMSS.mmm GMT`; OPERATORID = akun login | seperti TGLUPDATE `reinsurancetype`; format data DEV belum dilihat |
| A4 | Urutan bawaan grid: ID angka menurun | urutan RD Pega tidak ada di XML |
| A5 | Delete menghapus rate ber-IDUSEDBY sama tanpa memeriksa `RIRATEID` di Contract Retro Life (`TREATYBUSINESS_LIFE`) dan Product Name Life (`M_PRODUCTNAME_LIFE_PLAN`) | modul lain tidak disentuh; rujukan yang tertinggal menunjuk ringkasan yang sudah tidak ada |
| A6 | Edit nama ikut mengganti salinan `USEDBY` di baris `M_RATE_LIFE` ringkasan itu | supaya Rate Detail konsisten; salinan nama `RIRATE` di MCRL/MPNL tidak diganti |
| A7 | CSV: kepala wajib (urutan bebas, tanpa beda huruf); pemisah `;` bila kepala memuatnya, selain itu `,` — dengan `,` RATE berkoma desimal wajib dikutip (`"0,5"`), tanpa kutip baris ditolak (tidak ditebak) | label b5305 tidak menyebut pemisah |
| A8 | GENDER U/M/F; AGE wajib, CONTRACT boleh kosong, keduanya bulat 0-120; RATE desimal tak bertanda, koma atau titik, disimpan berkoma desimal (`0.50` → `0,50`) | isi `RATE_LIFE` DEV (`claimlife/docs/KATALOG-TABEL-PESERTA-DAN-TREATY.md` b118-b140) |
| A9 | Kombinasi (USEDBY, GENDER, AGE, CONTRACT) tidak kembar di berkas maupun terhadap baris yang sudah ada | DEV memuat 2.693 kembar warisan yang membuat `GetRateRetro` ambigu |
| A10 | USEDBY dicocokkan ke ringkasan bernama sama (tanpa beda huruf) atau dibuat baru; ringkasan lama ikut diperbarui OPERATORID/MODIFIEDDATE; nama yang cocok ke >1 ringkasan warisan ditolak | rule `SubmitRIRate_Act` tidak tersedia |
| A11 | Batas: 4 MB, 10.000 baris data, 100 nama berbeda per unggah; nama ≤ 200 byte | `RATE_LIFE` DEV 98.305 baris / 348 ringkasan |
| A12 | View Upload pun tertutup bagi View only | brief: View only tanpa Upload |
| A13 | Rate Detail: kombinasi (GENDER, AGE, CONTRACT) tidak boleh kembar di ringkasan yang sama (tambah dan ubah) | sama dengan A9; rule `AddToList_Act` tidak tersedia |
| A14 | Rate Detail Save ikut memperbarui MODIFY OPERATOR / MODIFY DATE ringkasan | sama dengan A10 |
| A15 | Pilihan radio GENDER = U, M, F (boleh kosong) | prompt list properti `.GENDER` tidak ada di XML; nilai dari isi `RATE_LIFE` DEV |
| A16 | Rate Detail tanpa Delete per baris | XML `InboxRIRate` tidak punya tombol Delete; baris hilang hanya lewat Delete ringkasan |

## RALAT

| # | Bunyi lama | Bunyi baru | Bukti |
| --- | --- | --- | --- |
| R7 (07-10-2026) | Bab "Tabel warisan": *"`M_RATE_LIFE` tabel fisik warisan Pega (baris rate, JSON) … nol DDL"*, *"`RATE_LIFE` view warisan atas `M_RATE_LIFE`; dibaca Rate Detail, Delete, dan pemeriksa kembar upload"*; R3 (baca `FOR UPDATE` + ganti kunci JSON di Go) untuk rincian | **Keputusan work owner 07-10-2026: rincian rate SATU tabel flat `M_RATE_LIFE`** berkolom ID, IDUSEDBY, USEDBY, TYPE, GENDER, CONTRACT, AGE, RATE (teks, isi apa adanya; RATE tetap teks); JSONDATA + constraint IS JSON dibuang, view `RATE_LIFE` dibuang (930), indeks `IX_M_RATE_LIFE_IDUSEDBY`; `M_RATE_LIFE` bukan lagi tabel warisan; rincian dibaca/ditulis lewat kolom (`rirl_json.go` dihapus); FLAG tidak dipindah (R5); 539 baris yatim dipindah apa adanya (item terbuka WO); claimlife, premiumlistlife, MPNL, MCRL membaca `M_RATE_LIFE` berkolom sama; `PEGA_M_RATE_LIFE` INVALID (diterima WO) | migrasi inti `929_m_rate_life_kolom.sql`, `930_m_rate_life_satu_tabel.sql`; uji `TestMigrasi929KolomRate`, `TestMigrasi930SatuTabel`, `TestSqlRate`, `TestNolJSONDanViewDiRepository`, `-tags=db` `TestDBMigrasiRincianFlatDanMundur`; pengecualian penjaga `TestKolomUangDesimalDanNolJSON` untuk `DROP COLUMN JSONDATA` 930 (pola 928) - tinjauan tim inti (`docs/PR-RIRATELIFE-FLAT.md`); `docs/LANGKAH-WO-RIRATELIFE-DETAIL-FLAT.md` |
| R6 (07-10-2026) | R4: *"VIEW `RATE_LIFE_SUMMARY` diganti TABEL FLAT bernama sama … `M_RATE_LIFE_SUMMARY` (JSON) TIDAK disentuh: cadangan + sumber alat pindah"* | **Keputusan work owner 07-10-2026: ringkasan cukup SATU tabel `M_RATE_LIFE_SUMMARY`** berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID (lebar = 926); JSONDATA + constraint IS JSON dibuang; tabel flat 926 dihapus (928); Pega tidak lagi menyimpan R/I Rate (prosedur `PEGA_M_RATE_LIFE_SUMMARY`, `PEGA_M_PLAN_LIFE_SUMMARY` INVALID - diterima WO); alat pindahflat dan berkas DBA lepas view dihapus; MPNL dan MCRL membaca kolom `M_RATE_LIFE_SUMMARY` | migrasi inti `927_m_rate_life_summary_kolom.sql`, `928_m_rate_life_summary_satu_tabel.sql`; uji `TestMigrasi927KolomRingkasan`, `TestMigrasi928SatuTabel`, `TestSqlTulisRingkasan`, `-tags=db` `TestDBMigrasiSatuTabelDanMundur`; [keputusan work owner 07-10-2026] DELETE 928 langkah 2 (buang baris `M_RATE_LIFE_SUMMARY` yang tidak ada di tabel flat) DISETUJUI; tabel flat tetap sumber kebenaran. Bukti POOLDATA (baca-saja, WO 07-10-2026): JSON 347 baris, flat 347 baris; ID hanya di JSON = `1000469` ("TEST RATE LIFE", sudah dihapus pengguna lewat aplikasi); ID hanya di flat = `1000471` (baru dari aplikasi); isi beda pada ID yang sama = 0; tidak ada tulisan Pega sesudah cutover 09:09. Pengecualian penjaga `TestKolomUangDesimalDanNolJSON` untuk `DROP COLUMN JSONDATA` 928 DISETUJUI WO, tetap ditinjau tim inti (`docs/PR-RIRATELIFE-FLAT.md`). |
| R5 (06-10-2026) | R4/K-F2: tabel flat memuat keenam kolom view termasuk `FLAG`; FLAG berstatus [penyimpangan sadar - menunggu WO] + pertanyaan arti `AP`/`PM`/`PY` | **FLAG tidak digunakan - keputusan WO 06-10-2026; nilai lama tetap di `M_RATE_LIFE_SUMMARY.JSONDATA`.** FLAG diisi layar Pega saat Submit (nol trigger/prosedur). Kolom FLAG dibuang dari CREATE TABLE 926 (belum dijalankan di DEV - berkas yang sama diubah, bukan migrasi baru); `_down` tetap memulihkan view asli lengkap. Modul dan alat pindah tidak membaca/menulis/menyalin FLAG; delta, COUNT, dan MINUS memakai 5 kolom. Pertanyaan FLAG untuk WO dicabut | keputusan WO 06-10-2026; `926_rate_life_summary_flat.sql`; uji `TestMigrasi926TeruraiDanBerpasangan` (nol FLAG di DDL), `TestSqlTulisRingkasan` (nol FLAG di SQL repository dan alat), `TestRencanaPindahPenuh` |
| R4 (06-10-2026) | K1: *"CRUD menulis ke `M_RATE_LIFE_SUMMARY`, view TIDAK di-DROP, nol DDL pada tabel/view warisan"* | **Keputusan work owner 06-10-2026 K-F1/K-F2**: VIEW `RATE_LIFE_SUMMARY` diganti TABEL FLAT bernama sama (migrasi inti 926, keenam kolom view, isi apa adanya); ringkasan dibaca DAN ditulis di tabel flat (kolom bernama, `TYPE`/`FLAG` tidak ditulis); `M_RATE_LIFE_SUMMARY` dibaca saja (cadangan, sumber alat pindah, ID terpakai); DROP VIEW = berkas DBA terpisah sebelum `-migrate` | perintah WO *"tabel M_RATE_LIFE_SUMMARY buat jadi flat …"*; `926_rate_life_summary_flat.sql`; uji `TestMigrasi926TeruraiDanBerpasangan`, `TestSqlTulisRingkasan`, `TestRencanaPindah`, `TestNilaiJSON`, `-tags=db` `rirl_db_test.go` |
| R3 (06-10-2026) | A2: *"Update = `JSON_MERGEPATCH(... RETURNING CLOB)` … butuh Oracle 18c+"* | Edit ringkasan, salinan nama rate, dan Edit / Save Rate Detail **tidak lagi memakai `JSON_MERGEPATCH`**: JSONDATA dibaca `SELECT … FOR UPDATE` di dalam transaksi, kunci yang berubah diganti di Go (`TerapkanKunci`: kunci Pega lain dan bentuk angka tetap, nilai kosong = kunci dibuang), lalu `UPDATE … SET JSONDATA = :1` - pola yang sudah berjalan di DEV (`masterproductnamelife`). Sisip tetap `JSON_OBJECT` (terbukti berjalan) | log server DEV 06-10-2026 10:55: `riratelife: menyimpan rate: repository: mengubah ringkasan: ORA-00907: missing right parenthesis` (Save Rate Detail; transaksi dibatalkan, nol baris setengah jadi); `backend/repository/rirl_json.go`; uji `TestTerapkanKunci`, `TestNolJSONMergepatchDiSQL` |
| R2 (06-10-2026) | Rule XML → kode, Detail: *"popup rincian rate (USEDBY, CONTRACT, GENDER, AGE, RATE), baca saja, 50 per halaman"* | Rate Detail bisa **tambah dan ubah** baris (form `InboxRIRate`), grid 20 per halaman, urut ID menurun, kolom ID · R/I RATE NAME · GENDER · CONTRACT · AGE · RATE · Edit | View Detail.xml (section `InboxRIRate`, work owner 06-10-2026): b3196, b9853, b10206, b7668; uji `TestRuteRateDetail`, `TestPeriksaIsianRate` |
| R1 (06-10-2026) | A1: *"`M_RATE_LIFE_SUMMARY (ID VARCHAR2(10), JSONDATA CLOB)`, kunci JSON `USEDBY`, `OPERATORID`, `MODIFIEDDATE` … definisi view ringkasan BELUM dibaca"* | View `RATE_LIFE_SUMMARY` = `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG FROM M_RATE_LIFE_SUMMARY a`; kolom `ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG` (semua VARCHAR2). Ketiga kunci yang ditulis **sudah benar** — nol perubahan nama di kode. `TYPE` dan `FLAG` ada di view tetapi **tidak dirujuk** XML `InboxSummaryRIRate` (form `.ID`/`.USEDBY`, grid `.ID`/`.USEDBY`/`.OPERATORID`/`.MODIFIEDDATE`) → tidak dibaca, tidak ditulis, tidak dikarang; ringkasan baru tanpa kunci `TYPE`/`FLAG` (pembaca `mastercontractretrolife` dan `masterproductnamelife` hanya membaca `ID, USEDBY`, tanpa saringan `FLAG`). Isi `TYPE`/`FLAG` pada data lama belum dilihat — bila ternyata dipakai menyaring, perlu keputusan WO | `ALL_VIEWS`, owner POOLDATA, dibaca WO 06-10-2026; konstanta `KolomViewRingkasan` (`backend/repository/rirl_tabel.go`); uji `TestKolomViewRingkasanEnamKolom`, `TestKunciDanKolomCocokDenganView` |

## Migrasi

Nol migrasi modul. Migrasi inti `922_m_nav_menu_riratelife.sql` (baris menu, langsung menyala),
`923_seq_rate_life.sql` (dua sequence, nilai awal = ID angka tertinggi + 1 dihitung di basis data tujuan),
`926_rate_life_summary_flat.sql` (tabel flat ringkasan, RALAT R4 - sudah jalan di DEV), `927_m_rate_life_summary_kolom.sql`
dan `928_m_rate_life_summary_satu_tabel.sql` (satu tabel, RALAT R6), `929_m_rate_life_kolom.sql` dan
`930_m_rate_life_satu_tabel.sql` (rincian flat, RALAT R7). 924/925 dipakai `ricommlife` di cabangnya.
⚠️ Karena 927/928 mengubah bentuk `M_RATE_LIFE_SUMMARY` dan 929/930 mengubah bentuk `M_RATE_LIFE`, kedua tabel itu
TIDAK lagi dinyatakan "Tabel warisan" di bawah (preseden `adjusterconsultant` 870); kolom yang dibuat migrasi tercatat
di `docs/STRUKTUR-TABEL-RIRATELIFE.md`.

## Menjalankan uji modul ini saja

Dari folder akar repo:

```powershell
go test ./modul/riratelife/...
npx vitest run modul/riratelife
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `RATE_LIFE` | view warisan Pega atas `M_RATE_LIFE.JSONDATA`; tidak pernah dibuat migrasi maju - DIBUANG `930_m_rate_life_satu_tabel.sql` (RALAT R7), dibangun ulang hanya oleh jalur mundur 930; nol pembaca sesudah 930 |
