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

## Keputusan work owner 06-10-2026 — ringkasan menjadi tabel flat (RALAT R4)

Perintah: *"tabel M_RATE_LIFE_SUMMARY buat jadi flat menampilkan data detail yang ada pada tabel view
RATE_LIFE_SUMMARY"*. Sesudah ditanya:

- **K-F1:** VIEW `RATE_LIFE_SUMMARY` **DIGANTI TABEL FLAT bernama sama** `RATE_LIFE_SUMMARY` (migrasi inti `926`, pola
  `RICOMM_LIFE`): berkas DBA pelepas view terpisah SEBELUM `-migrate`; `_down` = DROP TABLE lalu CREATE VIEW definisi
  asli; alat pindah uji-kering bawaan, `-jalankan` ditolak bila `IS_PEGA_PROD=true`. `M_RATE_LIFE_SUMMARY` (JSON, cacah
  masih berubah) TIDAK disentuh: cadangan + sumber alat pindah.
- **K-F2:** SEMUA 6 kolom view ikut, isi apa adanya (`ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG`, semua
  VARCHAR2). Modul tidak menulis `TYPE`/`FLAG` (tidak di XML): baris baru NULL, baris pindahan apa adanya.
  ⚠️ **FLAG dikeluarkan RALAT R5** (keputusan work owner 06-10-2026: FLAG tidak digunakan) - tabel flat 5 kolom.

Penerapan: `inti/backend/migrations/926_rate_life_summary_flat.sql` (+ `_down`), `docs/DBA-LEPAS-VIEW-RATE_LIFE_SUMMARY.sql`,
`backend/alat/pindahflat` + `repository/pindah.go`, repository ringkasan membaca DAN menulis tabel flat (kolom bernama;
jalur JSON ringkasan dibuang - `rirl_json.go` tetap untuk `M_RATE_LIFE`). Bukti lebar per kolom + **[penyimpangan
sadar]**: `docs/STRUKTUR-TABEL-RIRATELIFE.md`. **Urutan WO siap-tempel: [`docs/LANGKAH-WO-RIRATELIFE-FLAT.md`](docs/LANGKAH-WO-RIRATELIFE-FLAT.md).**

**Pembaca lain ikut membaca tabel flat:** `mastercontractretrolife` (autocomplete `R/I RATE`) dan
`masterproductnamelife` (`Choose R/I Rate`) menjalankan `SELECT ID, USEDBY FROM RATE_LIFE_SUMMARY` - kodenya TIDAK
diubah dan tetap berjalan atas tabel (nama dan kolom sama).

⚠️ **Risiko cutover:** ringkasan baru yang ditulis Pega (`PEGA_M_RATE_LIFE_SUMMARY` → `M_RATE_LIFE_SUMMARY`) SESUDAH
pemindahan tidak masuk tabel flat, sehingga tidak tampil di layar ini, di MCRL, maupun MPNL. Saran: bekukan penulisan
Pega ke R/I Rate (layar Pega `InboxSummaryRIRate` / upload Pega) sejak langkah DBA; bila belum dapat dibekukan,
jalankan alat pindah ULANG berkala (aman diulang: baris sama dilewati; ID sama berisi beda menghentikan putaran untuk
diperiksa). ID baru aplikasi tidak bertabrakan dengan ID Pega yang belum dipindah: `MaksID`/`AdaID` ringkasan memeriksa
tabel flat DAN `M_RATE_LIFE_SUMMARY`. Rincian rate yang ditulis Pega tetap terlihat (view `RATE_LIFE` tidak berubah),
tetapi ringkasannya baru tampil sesudah dipindah.

**Fakta DEV (dicek WO 06-10-2026, baca-saja):** `M_RATE_LIFE_SUMMARY` 347 baris dan MASIH BERUBAH (MODIFIEDDATE terbaru
`20261006T040628.169 GMT`) - dokumen dan alat tidak menanam angka baris; panjang maksimum isi ID 7, USEDBY 99,
MODIFIEDDATE 23, OPERATORID 17, TYPE kosong di semua baris (FLAG 2 - tidak dipakai, R5) (lebar 926 cukup - **[terverifikasi data DEV
06-10-2026]**); `UPPER(TRIM(USEDBY))` kembar = 0; dependensi / sinonim / grant view = nol.

### FLAG tidak digunakan — RALAT R5

Keputusan work owner 06-10-2026: **kolom FLAG tidak digunakan di R/I Rate Life.** Fakta (dicek WO, baca-saja): FLAG
diisi layar Pega saat Submit (nol trigger / prosedur); isi DEV `AP` 230, `PM` 106, `PY` 1, kosong 10; FLAG juga ada di
tiap baris detail `M_RATE_LIFE` (putaran ini tidak diubah); nol kode repo membacanya. Maka tabel flat 926 TIDAK memuat
FLAG, modul tidak membaca / menulisnya, dan alat pindah tidak menyalinnya. Nilai FLAG lama TIDAK hilang: tetap di
`M_RATE_LIFE_SUMMARY.JSONDATA` (cadangan), dan jalur mundur 926 memulihkan view lengkap dengan FLAG. TYPE tetap:
ringkasan baru NULL, Edit tidak menimpa, pindahan apa adanya.

### Cutover delta (alat pindah)

Cacah dibaca SAAT berjalan (sumber dan flat), tidak ada angka tetap. Putaran pertama tanpa `-sejak` (penuh); putaran
berikutnya `-sejak="<batas delta berikutnya>"` dari laporan putaran sebelumnya. Per ID:

| Keadaan | Tindakan |
| --- | --- |
| tidak ada di flat, MODIFIEDDATE sumber lebih baru dari `-sejak` (atau tanpa `-sejak`) | **baru** - disisip |
| tidak ada di flat, MODIFIEDDATE sumber tidak lebih baru dari `-sejak` / kosong | **dilewati** - dianggap dihapus aplikasi, tidak dihidupkan lagi, dilaporkan |
| kelima kolom flat sama | **sama** - dilewati |
| beda, MODIFIEDDATE sumber lebih baru dari flat | **berubah** - flat diperbarui kelima kolom (Pega lebih baru), ID dilaporkan |
| beda, MODIFIEDDATE flat sama / lebih baru / salah satu kosong atau tidak terbaca | **konflik** - TIDAK ditimpa, ID dilaporkan untuk diperiksa WO |
| hanya di flat | dibiarkan (ringkasan baru aplikasi) |

Aplikasi tetap hanya menulis tabel flat (MODIFIEDDATE = waktu simpan), sehingga tulisan aplikasi sesudah pemindahan
selalu lebih baru dari salinan Pega yang lama → konflik, bukan ditimpa. ⚠️ Batas aturan: bila Pega DAN aplikasi
sama-sama mengubah ringkasan yang sama, versi yang MODIFIEDDATE-nya lebih baru menang dan ID-nya tercantum di laporan
(`berubah` atau `konflik`). Alasan utama tetap: bekukan Pega.

### Urutan merge

**`modul/riratelife/implementasi` di-merge LEBIH DULU, lalu `modul/ricommlife/implementasi`.** Kedua cabang membawa
`inti/backend/db/koneksi.go` identik (commit `91221b2f` di sini = `349be34d` di ricommlife), migrasi inti 926 di sini
dan 924/925 di ricommlife. PR: `docs/PR-RIRATELIFE-FLAT.md`.

**Uji coba merge (06-10-2026)** - worktree sementara dari ujung riratelife, `git merge --no-commit
modul/ricommlife/implementasi`: `inti/backend/db/koneksi.go` bersih (identik). SATU konflik, di
`inti/backend/penjaga/rentang_test.go` (`TestSlotMenuBerjalanSesudah900`, daftar urutan pelari - kedua cabang menambah
berkasnya sendiri). Penyelesaian: satu daftar berurutan `… "923_seq_rate_life.sql", "924_ricomm_life.sql",
"925_m_nav_menu_ricommlife.sql", "926_rate_life_summary_flat.sql", "952_menu_tiruan.sql"`. Sesudahnya `go vet ./...`
bersih dan `go test ./...` = baseline (3 paket claimlife), ditambah `TestNolAlamatLayananDiKode` yang gagal HANYA karena
worktree sementara tidak memuat berkas `.env` (tidak masuk git) - bukan akibat merge.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-RIRATELIFE.md` — tabel flat (bukti lebar), tabel/view warisan, kunci JSON; `DBA-LEPAS-VIEW-RATE_LIFE_SUMMARY.sql`; `LANGKAH-WO-RIRATELIFE-FLAT.md` |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `alat/pindahflat/` `modul.go` (tanpa `migrations/`) |
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
| R5 (06-10-2026) | R4/K-F2: tabel flat memuat keenam kolom view termasuk `FLAG`; FLAG berstatus [penyimpangan sadar - menunggu WO] + pertanyaan arti `AP`/`PM`/`PY` | **FLAG tidak digunakan - keputusan WO 06-10-2026; nilai lama tetap di `M_RATE_LIFE_SUMMARY.JSONDATA`.** FLAG diisi layar Pega saat Submit (nol trigger/prosedur). Kolom FLAG dibuang dari CREATE TABLE 926 (belum dijalankan di DEV - berkas yang sama diubah, bukan migrasi baru); `_down` tetap memulihkan view asli lengkap. Modul dan alat pindah tidak membaca/menulis/menyalin FLAG; delta, COUNT, dan MINUS memakai 5 kolom. Pertanyaan FLAG untuk WO dicabut | keputusan WO 06-10-2026; `926_rate_life_summary_flat.sql`; uji `TestMigrasi926TeruraiDanBerpasangan` (nol FLAG di DDL), `TestSqlTulisRingkasan` (nol FLAG di SQL repository dan alat), `TestRencanaPindahPenuh` |
| R4 (06-10-2026) | K1: *"CRUD menulis ke `M_RATE_LIFE_SUMMARY`, view TIDAK di-DROP, nol DDL pada tabel/view warisan"* | **Keputusan work owner 06-10-2026 K-F1/K-F2**: VIEW `RATE_LIFE_SUMMARY` diganti TABEL FLAT bernama sama (migrasi inti 926, keenam kolom view, isi apa adanya); ringkasan dibaca DAN ditulis di tabel flat (kolom bernama, `TYPE`/`FLAG` tidak ditulis); `M_RATE_LIFE_SUMMARY` dibaca saja (cadangan, sumber alat pindah, ID terpakai); DROP VIEW = berkas DBA terpisah sebelum `-migrate` | perintah WO *"tabel M_RATE_LIFE_SUMMARY buat jadi flat …"*; `926_rate_life_summary_flat.sql`; uji `TestMigrasi926TeruraiDanBerpasangan`, `TestSqlTulisRingkasan`, `TestRencanaPindah`, `TestNilaiJSON`, `-tags=db` `rirl_db_test.go` |
| R3 (06-10-2026) | A2: *"Update = `JSON_MERGEPATCH(... RETURNING CLOB)` … butuh Oracle 18c+"* | Edit ringkasan, salinan nama rate, dan Edit / Save Rate Detail **tidak lagi memakai `JSON_MERGEPATCH`**: JSONDATA dibaca `SELECT … FOR UPDATE` di dalam transaksi, kunci yang berubah diganti di Go (`TerapkanKunci`: kunci Pega lain dan bentuk angka tetap, nilai kosong = kunci dibuang), lalu `UPDATE … SET JSONDATA = :1` - pola yang sudah berjalan di DEV (`masterproductnamelife`). Sisip tetap `JSON_OBJECT` (terbukti berjalan) | log server DEV 06-10-2026 10:55: `riratelife: menyimpan rate: repository: mengubah ringkasan: ORA-00907: missing right parenthesis` (Save Rate Detail; transaksi dibatalkan, nol baris setengah jadi); `backend/repository/rirl_json.go`; uji `TestTerapkanKunci`, `TestNolJSONMergepatchDiSQL` |
| R2 (06-10-2026) | Rule XML → kode, Detail: *"popup rincian rate (USEDBY, CONTRACT, GENDER, AGE, RATE), baca saja, 50 per halaman"* | Rate Detail bisa **tambah dan ubah** baris (form `InboxRIRate`), grid 20 per halaman, urut ID menurun, kolom ID · R/I RATE NAME · GENDER · CONTRACT · AGE · RATE · Edit | View Detail.xml (section `InboxRIRate`, work owner 06-10-2026): b3196, b9853, b10206, b7668; uji `TestRuteRateDetail`, `TestPeriksaIsianRate` |
| R1 (06-10-2026) | A1: *"`M_RATE_LIFE_SUMMARY (ID VARCHAR2(10), JSONDATA CLOB)`, kunci JSON `USEDBY`, `OPERATORID`, `MODIFIEDDATE` … definisi view ringkasan BELUM dibaca"* | View `RATE_LIFE_SUMMARY` = `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG FROM M_RATE_LIFE_SUMMARY a`; kolom `ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG` (semua VARCHAR2). Ketiga kunci yang ditulis **sudah benar** — nol perubahan nama di kode. `TYPE` dan `FLAG` ada di view tetapi **tidak dirujuk** XML `InboxSummaryRIRate` (form `.ID`/`.USEDBY`, grid `.ID`/`.USEDBY`/`.OPERATORID`/`.MODIFIEDDATE`) → tidak dibaca, tidak ditulis, tidak dikarang; ringkasan baru tanpa kunci `TYPE`/`FLAG` (pembaca `mastercontractretrolife` dan `masterproductnamelife` hanya membaca `ID, USEDBY`, tanpa saringan `FLAG`). Isi `TYPE`/`FLAG` pada data lama belum dilihat — bila ternyata dipakai menyaring, perlu keputusan WO | `ALL_VIEWS`, owner POOLDATA, dibaca WO 06-10-2026; konstanta `KolomViewRingkasan` (`backend/repository/rirl_tabel.go`); uji `TestKolomViewRingkasanEnamKolom`, `TestKunciDanKolomCocokDenganView` |

## Migrasi

Nol migrasi modul. Migrasi inti `922_m_nav_menu_riratelife.sql` (baris menu, langsung menyala),
`923_seq_rate_life.sql` (dua sequence, nilai awal = ID angka tertinggi + 1 dihitung di basis data tujuan), dan
`926_rate_life_summary_flat.sql` (tabel flat ringkasan, RALAT R4 - prasyarat: view dibuang DBA; 924/925 dipakai
`ricommlife` di cabangnya).

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
| `M_RATE_LIFE_SUMMARY` | tabel fisik warisan Pega (ringkasan R/I rate life, JSON, cacah masih berubah); sejak RALAT R4 (K-F1 06-10-2026) DIBACA saja - cadangan, sumber alat pindahflat, dan pemeriksa ID terpakai; nol DDL, nol tulis |
| `M_RATE_LIFE` | tabel fisik warisan Pega (baris rate, JSON); disisipkan Simpan Upload dan Rate Detail Save, diubah Rate Detail Edit, dihapus Delete ringkasan; nol DDL |
| `RATE_LIFE` | view warisan atas `M_RATE_LIFE`; dibaca Rate Detail, Delete, dan pemeriksa kembar upload |
