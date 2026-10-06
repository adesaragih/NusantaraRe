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
  membaca grid.
- **K2 ID baru:** dari **sequence Oracle** (`SEQ_M_RATE_LIFE_SUMMARY`, `SEQ_M_RATE_LIFE`, migrasi inti 923).
- **K3 lingkup:** termasuk **Detail** (rincian rate) dan **Upload CSV / View Upload / Simpan Upload**. XML rule-rule
  itu tidak tersedia → dirancang dari petunjuk format di XML (`Format excel : USEDBY, CONTRACT, GENDER, AGE, RATE`
  b5305); perilaku di luar itu = ASUMSI (tabel di bawah).

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-RIRATELIFE.md` — peta tabel/view warisan, kunci JSON, asumsi |
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
| Edit b9773 → `EditList_DT` b9853 (ID, UsedBy, Gender, Contract, Age, Rate) | isi form dari baris; simpan = `JSON_MERGEPATCH` GENDER/CONTRACT/AGE/RATE | `SqlUbahRate` |
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
| A2 | Update = `JSON_MERGEPATCH(... RETURNING CLOB)`, sisip = `JSON_OBJECT(... ABSENT ON NULL)` | butuh Oracle 18c+; kunci JSON lain milik Pega tetap |
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
| R2 (06-10-2026) | Rule XML → kode, Detail: *"popup rincian rate (USEDBY, CONTRACT, GENDER, AGE, RATE), baca saja, 50 per halaman"* | Rate Detail bisa **tambah dan ubah** baris (form `InboxRIRate`), grid 20 per halaman, urut ID menurun, kolom ID · R/I RATE NAME · GENDER · CONTRACT · AGE · RATE · Edit | View Detail.xml (section `InboxRIRate`, work owner 06-10-2026): b3196, b9853, b10206, b7668; uji `TestRuteRateDetail`, `TestPeriksaIsianRate` |
| R1 (06-10-2026) | A1: *"`M_RATE_LIFE_SUMMARY (ID VARCHAR2(10), JSONDATA CLOB)`, kunci JSON `USEDBY`, `OPERATORID`, `MODIFIEDDATE` … definisi view ringkasan BELUM dibaca"* | View `RATE_LIFE_SUMMARY` = `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG FROM M_RATE_LIFE_SUMMARY a`; kolom `ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG` (semua VARCHAR2). Ketiga kunci yang ditulis **sudah benar** — nol perubahan nama di kode. `TYPE` dan `FLAG` ada di view tetapi **tidak dirujuk** XML `InboxSummaryRIRate` (form `.ID`/`.USEDBY`, grid `.ID`/`.USEDBY`/`.OPERATORID`/`.MODIFIEDDATE`) → tidak dibaca, tidak ditulis, tidak dikarang; ringkasan baru tanpa kunci `TYPE`/`FLAG` (pembaca `mastercontractretrolife` dan `masterproductnamelife` hanya membaca `ID, USEDBY`, tanpa saringan `FLAG`). Isi `TYPE`/`FLAG` pada data lama belum dilihat — bila ternyata dipakai menyaring, perlu keputusan WO | `ALL_VIEWS`, owner POOLDATA, dibaca WO 06-10-2026; konstanta `KolomViewRingkasan` (`backend/repository/rirl_tabel.go`); uji `TestKolomViewRingkasanEnamKolom`, `TestKunciDanKolomCocokDenganView` |

## Migrasi

Nol migrasi modul. Migrasi inti `922_m_nav_menu_riratelife.sql` (baris menu, langsung menyala) dan
`923_seq_rate_life.sql` (dua sequence, nilai awal = ID angka tertinggi + 1 dihitung di basis data tujuan).

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
| `M_RATE_LIFE_SUMMARY` | tabel fisik warisan Pega (ringkasan R/I rate life, JSON); modul ini menambah, mengubah, dan menghapus barisnya, tidak pernah membuat atau mengubah strukturnya (K1 keputusan work owner 05-10-2026: nol DDL) |
| `M_RATE_LIFE` | tabel fisik warisan Pega (baris rate, JSON); disisipkan Simpan Upload dan Rate Detail Save, diubah Rate Detail Edit, dihapus Delete ringkasan; nol DDL |
| `RATE_LIFE_SUMMARY` | view warisan atas `M_RATE_LIFE_SUMMARY`; dibaca grid; TIDAK di-DROP (K1) |
| `RATE_LIFE` | view warisan atas `M_RATE_LIFE`; dibaca Rate Detail, Delete, dan pemeriksa kembar upload |
