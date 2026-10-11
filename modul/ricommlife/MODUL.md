# Modul `ricommlife` — R/I Comm Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner 06-10-2026:
*"membuat modul baru di Master Treaty dengan nama R/I Comm Life, panduannya baca dari xml di D:\NUSARE DEV\Menu RI Comm,
konsepnya hampir sama dengan menu R/I Rate, membuat CRUD, dan detail bisa di save dan edit"*. Panduan: section Pega
`InboxSummaryRIComm` (kelas `ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY`, judul "R/I COMM SUMMARY" b382) dan `InboxRIComm`
(kelas `ASM-FW-GISFW-Int-RI_COMM_LIFE`, judul "R/I COMM DETAIL" b349). Templat: modul saudara `riratelife`.

⛔ **Tanpa migrasi sendiri** (`—` di bawah, `tandaTanpaMigrasi` di `inti/backend/penjaga`): tabel flat `RICOMM_LIFE` =
migrasi inti `924` (riwayat - dibuang 934), baris menunya = migrasi inti `925`, satu tabel per jenis data = migrasi inti
`931`-`934` (RALAT R1).

| Kunci | Nilai |
| --- | --- |
| Nama modul | `ricommlife` |
| Folder korpus | `R/I Comm Life` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-RICOMMLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `—` |
| Slot menu | `—` |
| Prefix rute API | `/api/ri-comm-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Keputusan work owner 08-10-2026 — SATU tabel per jenis data (RALAT R1)

*"R/I COMM LIFE DIBUAT SATU TABEL PER JENIS DATA, sama seperti R/I Rate Life (ringkasan 927/928, detail 929/930)"*:
`M_RICOMM_LIFE_SUMMARY` = tabel flat `ID, USEDBY, MODIFIEDDATE, OPERATORID` (kolom view `RICOMM_LIFE_SUMMARY`; JSONDATA
dan view dibuang); `M_RICOMM_LIFE` = tabel flat rincian berkolom PERSIS tabel `RICOMM_LIFE` 924 (`IDUSEDBY
VARCHAR2(10)`, `USEDBY VARCHAR2(200)`, `CONTRACT NUMBER(5)`, `YEAR NUMBER(5)`, `COMM NUMBER(38,8)`; JSONDATA dibuang);
tabel `RICOMM_LIFE` (924) DIHAPUS. Nol tabel baru. Sequence dan rumus ID TETAP (`site || LPAD(seq, 6, '0')`).

Fakta POOLDATA (WO 08-10-2026, baca-saja): `M_RICOMM_LIFE_SUMMARY` = ID PK + JSONDATA (`ENSURE_M_RICOMM_LIFE_SUMMARY_JSON`,
33 byte), 3 baris; `M_RICOMM_LIFE` = ID PK + JSONDATA (`ENSURE_M_RICOMM_LIFE_JSON`), 0 baris; `RICOMM_LIFE` 2 baris
(`PK_RICOMM_LIFE`, `IX_RICOMM_LIFE_IDUSEDBY`, satu CHECK sistem `SYS_C0015437`); panjang maksimum ID 7, IDUSEDBY 7,
USEDBY 18; prosedur `PEGA_M_RICOMM_LIFE` dan `PEGA_M_RICOMM_LIFE_SUMMARY` akan INVALID (diterima WO).

Penerapan (pola dan pelajaran riratelife 927-930 + tinjauan independennya):

- Migrasi inti `931_m_ricomm_life_summary_kolom.sql` (ADD tiga kolom, berdiri sendiri), `932_m_ricomm_life_summary_satu_tabel.sql`
  (isi dari JSONDATA + buang JSONDATA lewat blok berpelindung katalog, indeks `IX_M_RICOMM_LIFE_SUMMARY_NAMA`, DROP VIEW
  terakhir), `933_m_ricomm_life_kolom.sql` (ADD lima kolom), `934_m_ricomm_life_satu_tabel.sql` (isi dari JSON bila ada,
  buang JSONDATA, timpa + sisip dari `RICOMM_LIFE` SESUDAH JSONDATA dibuang, indeks `IX_M_RICOMM_LIFE_IDUSEDBY`, DROP
  TABLE `RICOMM_LIFE` terakhir). USEDBY ringkasan 200 byte (= salinan di rincian dan `BatasNama`). CHECK `SYS_C0015437`
  tidak disalin (= `ID` NOT NULL 924; ID sudah PK). Jalur mundur aman diulang; `931_down`/`933_down` mulai dengan
  pelindung gagal-keras `UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0` (ORA-00904 bila JSONDATA sudah dibuang langkah
  maju yang tidak tercatat). JSONDATA dikembalikan NULLABLE (bukti NOT NULL tidak ada di repo; LANGKAH-WO (a) mencatat).
- Kode: ringkasan dan rincian dibaca/ditulis lewat kolom (`backend/repository/ricl.go`, `ricl_tabel.go`); `pxObjClass`
  tidak lagi disimpan (bukan kolom view); ID terpakai satu tabel per jenis; transaksi ringkasan + rincian tetap satu
  transaksi Go.
- **Dihapus** (alasan: objek yang mereka layani tidak ada lagi sesudah 932/934): `backend/alat/pindahflat/` dan
  `repository/pindah.go` (+ uji) - memindah JSON `M_RICOMM_LIFE` ke `RICOMM_LIFE`, kini dikerjakan 934 sekali jalan;
  `repository/ricl_json.go` (+ uji) - baca-ubah-tulis JSONDATA ringkasan; `docs/DBA-LEPAS-VIEW-RICOMM_LIFE.sql` - langkah
  DBA prasyarat 924 yang SUDAH dijalankan di DEV (komentar 924, yang tidak diubah, masih merujuknya; isinya di riwayat
  git). `docs/LANGKAH-WO-RICOMMLIFE.md` diarsipkan (spanduk ARSIP; 924 merujuknya).
- `inti/backend/db/koneksi.go` (`db.Koneksi`) kini TIDAK dipakai kode mana pun (satu-satunya pemakai = alat pindahflat).
  TIDAK dihapus - kandidat bersih-bersih tim inti (`docs/PR-RICOMMLIFE.md`).
- Langkah WO: `docs/LANGKAH-WO-RICOMMLIFE-SATU-TABEL.md` + berkas SQL\*Plus `docs/sql/satu_tabel_*.sql`.

| # | Bunyi lama | Bunyi baru | Bukti |
| --- | --- | --- | --- |
| R1 (08-10-2026) | butir 2 (DROP VIEW lewat berkas DBA), butir 3 (*"M_RICOMM_LIFE tidak disentuh; alat pindahflat"*), butir 6 (*"ringkasan tetap JSON"*), bab "Tabel warisan" (`M_RICOMM_LIFE_SUMMARY`, `M_RICOMM_LIFE`), A8 (alat pindah) | **Satu tabel per jenis data**: ringkasan = kolom `M_RICOMM_LIFE_SUMMARY`, rincian = kolom `M_RICOMM_LIFE`; JSONDATA, view `RICOMM_LIFE_SUMMARY`, dan tabel `RICOMM_LIFE` dibuang (931-934); alat pindahflat, `ricl_json.go`, berkas DBA dihapus; A8 dicabut | keputusan WO 08-10-2026; uji `TestMigrasi931KolomRingkasan`, `TestMigrasi932SatuTabel`, `TestMigrasi933KolomKomisi`, `TestMigrasi934SatuTabel`, `TestSqlRingkasan`, `TestSqlKomisi`, `TestNolJSONDanObjekLamaDiRepository`, `-tags=db` `TestDBMigrasiSatuTabelDanMundur`, `TestDBKembarKepemilikanDanDeleteBerantai` |

## Keputusan work owner 06-10-2026 (dikutip, diikuti persis)

> Butir 2, 3, 6 dan A8 digantikan RALAT R1 (08-10-2026) di atas; tabel ini riwayat.

Fakta DEV (dicek WO 06-10-2026):

- View `RICOMM_LIFE` = `SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.CONTRACT, a.JSONDATA.YEAR,
  a.JSONDATA.COMM FROM M_RICOMM_LIFE a`. `M_RICOMM_LIFE` 0 baris. Nol kode repo membaca RICOMM_LIFE; satu-satunya
  dependensi view = M_RICOMM_LIFE; `PEGA_M_RICOMM_LIFE` menulis M_RICOMM_LIFE; nol sinonim/grant.
- `PEGA_M_RICOMM_LIFE`: `SELECT ID INTO id_site FROM M_SITE_DATABASE WHERE CURRENT_SITE='1'` (b13);
  `id := id_site || lpad(to_Char(M_RICOMM_LIFE_SEQ.nextval),6,'0')` (b21); `INSERT INTO M_RICOMM_LIFE (ID, JSONDATA)`
  (b29); COMMIT di dalam prosedur. `M_SITE_DATABASE`: ID NUMBER; baris (1,'1') dan (2,'0') -> ID 7 karakter.
- View `RICOMM_LIFE_SUMMARY` = `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID FROM
  POOLDATA.M_RICOMM_LIFE_SUMMARY a`. `PEGA_M_RICOMM_LIFE_SUMMARY` = rumus sama dengan `M_RICOMM_LIFE_SUMMARY_SEQ`. Isi
  1 baris: ID 1000003, JSON {MODIFIEDDATE, OPERATORID, pxObjClass `ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY`, USEDBY "RI COMM
  RETRO"}.
- Sequence `M_RICOMM_LIFE_SEQ` last 43; `M_RICOMM_LIFE_SUMMARY_SEQ` last 4. Menu MASTER TREATY: … reinsurancetype 8,
  riratelife 9 (10 kosong).

Arahan → penerapan:

| # | Arahan (ringkas) | Penerapan |
| ---: | --- | --- |
| 1 | Tabel flat `RICOMM_LIFE` menggantikan VIEW, migrasi inti, tipe ikut XML, K6 bila XML diam, semua NULLABLE kecuali PK | `inti/backend/migrations/924_ricomm_life.sql`; bukti per kolom + **[penyimpangan sadar]**: `docs/STRUKTUR-TABEL-RICOMMLIFE.md` |
| 2 | DROP VIEW TIDAK di run -migrate yang sama; berkas DBA terpisah; `_down` = DROP TABLE + CREATE VIEW asli | `docs/DBA-LEPAS-VIEW-RICOMM_LIFE.sql`; `924_ricomm_life_down.sql`; uji `TestMigrasi924TeruraiDanBerpasangan`, `TestBerkasDBALepasView`. Penjaga inti menerima `CREATE VIEW` di `_down` tanpa entri baru |
| 3 | `M_RICOMM_LIFE` tidak disentuh; alat pindahflat uji-kering bawaan, tolak `IS_PEGA_PROD=true` | `backend/alat/pindahflat`, `repository/pindah.go`; SELURUH pemindahan di SATU koneksi terkunci (`db.Koneksi` = `*sql.Conn`, `SesiPindah`: ALTER SESSION lalu transaksi di koneksi yang sama) dan angka diurai Go tanpa NLS (`AngkaOracle`); uji `TestPeriksaMode`, `TestRencanaPindah`, `TestTolakProduksiSebelumKoneksi`, `TestAngkaOracleTanpaNLS`, `-tags=db` `TestDBSesiPindahSatuKoneksi` |
| 4 | ID baru = site saat berjalan `||` LPAD(seq warisan, 6); nol sequence baru; > VARCHAR2(10) = galat | `models.BentukID`, `services.pemberiID`; uji `TestBentukIDRumusProsedur` (site 1, seq 44 → `1000044`) |
| 5 | Transaksi milik Go, satu per simpan, PEGA_* tidak dipanggil | setiap tulis lewat `Layanan.tx`; Save detail = rincian + ringkasan satu transaksi |
| 6 | Ringkasan tetap JSON; sisip `JSON_OBJECT` + `pxObjClass`; ubah baca-ubah-tulis; MODIFIEDDATE format Pega | `repository/ricl_json.go`; uji `TestTerapkanKunci`, `TestNolJSONMergepatchDiSQL` |
| 7 | Detail flat per ringkasan, tambah + edit, kembar (CONTRACT, YEAR) ditolak, Delete beserta detail, Edit nama ikut ganti; Upload gaya riratelife | `services/detail.go`, `ringkasan.go`, `unggah.go` |
| 8 | Menu MASTER TREATY URUTAN 10 bila ada bukti; tanpa bukti tetap dibuat atas perintah WO | `925_m_nav_menu_ricommlife.sql`, MASTER TREATY URUTAN 10 - **[keputusan work owner 06-10-2026]** sebagai dasarnya. Korpus XML menu (portal/navigasi) tidak ada di mesin ini; kedua XML yang tersedia hanya section, tanpa struktur menu |
| 9 | Prefix `/api/ri-comm-life`, nama `ricommlife`, hak Full / View only | `handlers/rute.go`, `HakLihat` di `backend/modul.go` |

## Urutan langkah work owner (DEV)

➡️ **RALAT R1 (931-934): [`docs/LANGKAH-WO-RICOMMLIFE-SATU-TABEL.md`](docs/LANGKAH-WO-RICOMMLIFE-SATU-TABEL.md)** -
prasyarat ⛔ (Pega R/I Comm dan backend dihentikan, nol kunci DML), angka acuan dua kali, cadangan CSV, `-migrate`,
verifikasi, restart; tabel pemulihan K0-K5; jalur mundur.

RIWAYAT (924, sudah jalan di DEV): [`docs/LANGKAH-WO-RICOMMLIFE.md`](docs/LANGKAH-WO-RICOMMLIFE.md) (ARSIP) - (a) berkas
DBA pelepas view, (b) `-migrate` dari cabang ini (muat-env), (c) kueri verifikasi (`ALL_OBJECTS`, `T_MIGRASI`,
`M_NAV_MENU`), (d) INSERT hak menu superadmin (pemegang `kelolauser`), (e) restart backend, dan pemulihan bila (b)
terjalankan sebelum (a). Alat `pindahflat` (M_RICOMM_LIFE 0 baris DEV) opsional sesudah (c).

⛔ Bila (a) dilewati: pra-terbang menemukan view bernama sama dengan kolom yang sama, CREATE TABLE dijawab ORA-00955 dan
dilewati, lalu `CREATE INDEX` atas view gagal (ORA-01702) dan `-migrate` berhenti; pemulihan di berkas langkah itu.

**Jalur mundur** `924_ricomm_life_down.sql` (DROP TABLE lalu CREATE VIEW asli atas `M_RICOMM_LIFE`) hanya aman
SEBELUM aplikasi menulis rincian: baris yang ditulis aplikasi tidak ada di `M_RICOMM_LIFE` dan hilang. `cmd/api
-migrate-down` dipagari skema uji - bukan jalur mundur DEV.

## Urutan merge

**`modul/riratelife/implementasi` di-merge LEBIH DULU, lalu cabang ini (`modul/ricommlife/implementasi`)** - keputusan
work owner 06-10-2026. Kedua cabang membawa `inti/backend/db/koneksi.go` identik (`349be34d` di sini = `91221b2f` di
riratelife), migrasi inti 924/925 di sini dan 926-928 di riratelife (nomor tidak bertabrakan). PR: `docs/PR-RICOMMLIFE.md`.

**Uji coba merge (07-10-2026, diulang sesudah riratelife RALAT R6 - migrasi 927/928)** - worktree sementara dari ujung riratelife, `git merge --no-commit
modul/ricommlife/implementasi`: `inti/backend/db/koneksi.go` bersih (identik). SATU konflik, di
`inti/backend/penjaga/rentang_test.go` (`TestSlotMenuBerjalanSesudah900`, daftar urutan pelari - kedua cabang menambah
berkasnya sendiri). Penyelesaian: satu daftar berurutan `… "923_seq_rate_life.sql", "924_ricomm_life.sql",
"925_m_nav_menu_ricommlife.sql", "926_rate_life_summary_flat.sql", "927_m_rate_life_summary_kolom.sql",
"928_m_rate_life_summary_satu_tabel.sql", "952_menu_tiruan.sql"`. Sesudahnya `go vet ./...`
bersih dan `go test ./...` = baseline (3 paket claimlife), ditambah `TestNolAlamatLayananDiKode` yang gagal HANYA karena
worktree sementara tidak memuat berkas `.env` (tidak masuk git) - bukan akibat merge.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-RICOMMLIFE.md` (bukti tipe kolom, indeks), `LANGKAH-WO-RICOMMLIFE-SATU-TABEL.md` + `sql/satu_tabel_*.sql` (urutan WO RALAT R1), `LANGKAH-WO-RICOMMLIFE.md` (ARSIP 924), `PR-RICOMMLIFE.md` |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `modul.go` (tanpa `migrations/`) |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `ricommlife.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Rule XML → kode

| XML | Perilaku | Kode |
| --- | --- | --- |
| `InboxSummaryRIComm`: `InputParam.ERRMSG` b915 | pesan galat di atas form | `Gagal` di `RICommLife.tsx` |
| `.ID` b1091 disabled b1109-b1110 | ID tampil, tidak dapat diisi; baru = site + LPAD(seq, 6) | `services.pemberiID` |
| `.USEDBY` b1274 wajib b1289 (label "USEDBY" b1243, grid "R/I COMM NAME" b7220) | wajib, dipangkas, tidak kembar tanpa beda huruf, ≤ 200 byte | `services.Simpan` |
| Save b1799 → `AddToListSummary_Act` b1823; Cancel b2092 → `NewDataSummary_DT` b2120 (`DATASHOW = 'IsEdit'` b2261) | Add / Edit; Cancel hanya saat Edit | `POST` / `PUT /api/ri-comm-life` |
| Grid `BrowseRICommSummary`: ID · R/I COMM NAME · MODIFY OPERATOR b7377 (`.OPERATORID` b8348) · MODIFY DATE b7521 (`.MODIFIEDDATE` b8526, `Date-Short-Custom-YYYY` b8573); sort ID ASC b9857; 50 baris b10081 | filter ID / nama, urut per kolom, bawaan ID menaik, tanggal `DD-MM-YYYY` WIB | `GET /api/ri-comm-life` |
| Edit b8754 → `EditListSummary_DT` b8781 (ID, UsedBy) | isi form dari baris; ganti nama = ganti USEDBY rinciannya (satu transaksi) | `RICommLife.tsx`, `services.Simpan` |
| Detail b9042 → `setIDUsedBy_Act` b9059 + harness `InboxRIComm` b9096 (popup, `pyWindowName` "Ri Comm") | popup R/I COMM DETAIL | `KomisiDetail.tsx` |
| Delete b9680 → `DeleteSummaryDetail` b9704, DeleteID=.ID b9719 | ringkasan BESERTA rinciannya, satu transaksi, konfirmasi menyebut jumlah | `DELETE /{id}`, `services.Hapus` |
| Upload CSV b3189 → `UploadCSV_RICOMM` b3258; View Upload b3706 → `ViewCSVResult_RIComm` b3725; Simpan Upload b4771 → `SubmitRIComm_Act` b4795; label b5569 | pilih berkas; pratinjau tanpa menulis; simpan semua-atau-tidak satu transaksi | `UnggahCSV.tsx`, `POST /unggah/pratinjau`, `POST /unggah` |
| `InboxRIComm`: ID b1325; IDUSEDBY `TempIDUsedBy.ID` b1509 (`1=2` b1633); USEDBY `TempIDUsedBy.USEDBY` b1718 disabled wajib | IDUSEDBY dan USEDBY dari ringkasan yang dilihat, bukan isian (badan berisi `usedby` = 400) | `services.SimpanKomisi` |
| CONTRACT b1904 pxNumber wajib, placeholder 0 b1927 (`ChangeDotToPoint_DT` b1953) · YEAR b2185 wajib placeholder 0 b2205 pyMax 4 b2206 · COMM b2392 wajib | CONTRACT bulat 0-99999; YEAR tepat 4 angka; COMM desimal tidak negatif NUMBER(38,8), koma atau titik (A2) | `models.PeriksaIsianKomisi` |
| Save b2931 → `AddToList_Act` b2955; Cancel b3230 → `NewData_DT` b3258 (`IsEdit` b3396); Clear Field b1100 → `clearInputFieldRIComm_act` b1123 | tambah / ubah satu baris; Cancel hanya saat Edit; Clear Field mengosongkan form | `POST /{id}/detail`, `PUT /{id}/detail/{detailId}` |
| Upload container `InboxRIComm` `1=2` b4041 | tidak ditampilkan di popup (ada di halaman ringkasan) | — |
| Grid `BrowseRICommLife_RD` (idusedby b7396): ID · R/I COMM NAME · CONTRACT · YEAR · COMM (pxNumber b9101) · Edit b9295 → `EditList_DT` b9323; sort ID ASC b9515; 50 baris b9716 | ID menaik; paging 50 | `GET /{id}/detail`, `SqlDaftarKomisi` |

## Asumsi A1-A9 — bukti XML atau bawaan work owner

Arahan work owner 06-10-2026: setiap butir dicari buktinya di XML yang tersedia (`D:\NUSARE DEV\Menu RI Comm\
InboxSummaryRIComm.xml` = **S**, `…\InboxRIComm.xml` = **D**: validasi, When, kontrol layar, Activity simpan/hapus).
Hasil pencarian: kedua XML hanya rule **section** - Activity `AddToList_Act`, `AddToListSummary_Act`, `SubmitRIComm_Act`,
`DeleteSummaryDetail`, `setIDUsedBy_Act` hanya **dirujuk** (nama + parameter), isinya tidak ada; nol rule validasi,
nol When selain `1=1`, `1=2`, `InputParam.DATASHOW` (S b1966/b2261/b3582/b4647, D b1633/b3101/b3396/b4718). Ada bukti
→ **[terverifikasi]**; tidak ada → bawaan WO, **[penyimpangan sadar - menunggu WO]**.

| # | Aturan yang berlaku | Status dan bukti |
| --- | --- | --- |
| A1 | Kombinasi (CONTRACT, YEAR) kembar dalam satu ringkasan DITOLAK (Detail Save dan Upload; nol depan setara) | **[penyimpangan sadar - menunggu WO]** - isi `AddToList_Act` (D b2955) dan `SubmitRIComm_Act` (S b4795) tidak ada |
| A2 | CONTRACT bulat 0-99999; YEAR TEPAT 4 angka (1000-9999); COMM desimal TIDAK NEGATIF, tanpa batas 100, ≤ 30 + 8 angka | **[penyimpangan sadar - menunggu WO]** - CONTRACT pxNumber (D b1907) dengan `Precision` dan `NegativeNumberFormat` kosong (D b2089, b2095); YEAR hanya `pyMax` 4 (D b2206) = batas atas, bukan "tepat"; COMM pxTextInput (D b2395) tanpa pyMax/rentang |
| A3 | Save detail dan Simpan Upload memperbarui MODIFIEDDATE dan OPERATORID ringkasan (transaksi yang sama) | **[penyimpangan sadar - menunggu WO]** - isi Activity tidak ada |
| A4 | Edit nama ringkasan ikut mengganti USEDBY rinciannya (satu transaksi) | **[keputusan work owner 06-10-2026 butir 7]**; `EditListSummary_DT` (S b8781) hanya memuat ID/UsedBy ke form |
| A5 | CSV: kepala USEDBY, CONTRACT, YEAR, COMM (urutan bebas); pemisah `;`/`,`; 4 MB, 10.000 baris, 100 nama; nama dicocokkan ke ringkasan atau dibuat baru | kolom **[terverifikasi]** S b5569 (`Format excel : USEDBY, CONTRACT, YEAR, COMM`); sisanya **[penyimpangan sadar - menunggu WO]** (butir 7 "gaya riratelife") |
| A6 | R/I COMM DETAIL tanpa Delete per baris; View only tanpa Add/Edit/Delete/Upload/View Upload | tanpa Delete **[terverifikasi]**: grid D hanya tombol Edit b9295 (`EditList_DT` b9323), nol `Delete` di D; hak View only = **[keputusan work owner 06-10-2026 butir 9]** |
| A7 | ID yang ternyata terpakai dilewati (≤ 100 nomor); nomor sequence > 6 angka = galat | **[penyimpangan sadar - menunggu WO]** - rumus `site || LPAD(seq, 6)` = keputusan WO butir 4 (prosedur `PEGA_M_RICOMM_LIFE` b13/b21) |
| A8 | ~~Alat pindah: kosong = NULL; normalisasi teks angka menahan `-jalankan`~~ | **DICABUT RALAT R1** - alat pindahflat dihapus; pemindahan = migrasi inti 934 |
| A9 | Delete ringkasan menghapus rinciannya dalam SATU transaksi, TANPA cek rujukan produk | berantai **[terverifikasi]**: S b9704 `DeleteSummaryDetail` (nama = summary + detail), param `DeleteID=.ID` b9719; tanpa cek rujukan **[penyimpangan sadar - menunggu WO]**. ⚠️ Risiko: produk life Pega yang memegang `RICOMMID` ringkasan itu tertinggal menunjuk ID yang sudah tidak ada (modul lain tidak disentuh) |
| — | Grid urut ID menaik (ringkasan dan detail) | **[terverifikasi]** sort kolom pertama ASC, `pySortOrder` 1: S b9857/b9863, D b9515/b9520 |

## Migrasi

Nol migrasi modul. Migrasi inti `924_ricomm_life.sql` (tabel flat - riwayat, sudah jalan di DEV, dibuang 934),
`925_m_nav_menu_ricommlife.sql` (baris menu MASTER TREATY URUTAN 10, langsung menyala), dan RALAT R1:
`931_m_ricomm_life_summary_kolom.sql`, `932_m_ricomm_life_summary_satu_tabel.sql`, `933_m_ricomm_life_kolom.sql`,
`934_m_ricomm_life_satu_tabel.sql`. ⚠️ Karena 931-934 mengubah bentuk `M_RICOMM_LIFE_SUMMARY` dan `M_RICOMM_LIFE`,
keduanya TIDAK lagi dinyatakan "Tabel warisan" di bawah (preseden `adjusterconsultant` 870, riratelife 927/929); kolom
yang dibuat migrasi tercatat di `docs/STRUKTUR-TABEL-RICOMMLIFE.md`.

## Menjalankan uji modul ini saja

Dari folder akar repo:

```powershell
go test ./modul/ricommlife/...
go vet -tags db ./modul/ricommlife/...   # uji seam Oracle (skema uji) - butuh ORACLE_DSN skema uji
npx vitest run modul/ricommlife
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `RICOMM_LIFE_SUMMARY` | view warisan Pega atas `M_RICOMM_LIFE_SUMMARY.JSONDATA`; tidak pernah dibuat migrasi maju - DIBUANG `932_m_ricomm_life_summary_satu_tabel.sql` (RALAT R1), dibangun ulang hanya oleh jalur mundur 932; nol pembaca sesudah 932 |
