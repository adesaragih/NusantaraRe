# Struktur tabel — Endorsment Fac In (`endorsmentfacin`)

⛔ **Dibaca penjaga** `inti/backend/penjaga` (`TestKolomDDLCocokDenganStruktur`,
`TestGolonganTipeDDLCocokDenganStruktur`): setiap judul `## NAMA_TABEL` di bawah wajib sama persis dengan
kolom `CREATE TABLE` migrasi modul ini. Tabel yang hanya DIBACA ditulis dengan judul berskema
(`## POOLDATA.…`), sehingga tidak dibaca sebagai tabel milik modul ini.

Rentang migrasi `260-299`; slot menu `966` (UPDATE `DIMIGRASI`, bukan tabel). Nomor `260` dicadangkan
untuk `T_WORK_POLIS`/sequence kasus `EDM-` (irisan I3, keputusan work owner 11-10-2026).

## M_EDM_PENGECUALIAN_POLIS

Irisan I2 — **dibuat** modul ini, migrasi `261_m_edm_pengecualian_polis.sql`. Tabel konfigurasi
(ADR-0035, keputusan work owner 11-10-2026), **diisi DBA**; kosong = tanpa bypass. Pengganti literal
nomor polis yang tertanam di prakondisi Pega `[terverifikasi]` (nilai tidak disalin):

- `JENIS = 'TANGGAL'` — `Endorsment Fac In\Activity\CheckEDMPolisDate.xml` langkah 2 (1 literal):
  validasi tanggal endorsement dilewati.
- `JENIS = 'BATAL'` — `Endorsment Fac In\Activity\SetValueToEDMWork.xml` langkah 5 (2 literal): kasus
  tetap dibuat walau EdmType polis lama = 1. Belum dibaca kode (irisan I3).

Cacah dua cara atas `pyStepsPreCondParamsWhen` (urai XML / cacah baris): 1 dan 2, irisan 0. PK
`(NOPOLIS, JENIS)` = keputusan agent (satu polis boleh di kedua daftar). Bila tabel belum dibuat
(ORA-00942), pemeriksaan tanggal menganggap daftar kosong dan mengirim `catatan`.

| Kolom | Tipe | Tipe DDL (migrasi 261) | Isi |
| --- | --- | --- | --- |
| `NOPOLIS` | teks | VARCHAR2(100) NOT NULL | nomor polis persis seperti di `JSON_POLIS.NOPOLIS` (lebar sama) |
| `JENIS` | teks | VARCHAR2(10) NOT NULL | `TANGGAL` atau `BATAL` (CHECK) |
| `KETERANGAN` | teks | VARCHAR2(1000) | alasan pengecualian, diisi DBA |
| `TGL_BUAT` | DATE | DATE DEFAULT SYSDATE | saat baris dibuat |

## T_EDM_FACIN_KASUS

Irisan I3 — **dibuat** modul ini, migrasi `262_t_edm_facin_kasus.sql` (+ indeks `IX_T_EDM_FACIN_KASUS_POLIS` atas
`OLD_POLICY_NO` untuk gerbang 2). Kepala kasus `EDM-<5 digit>` (sequence `SEQ_WORK_POLIS_EDM`, migrasi 260);
**berbagi PK** dengan `T_WORK_POLIS.ID`, tanpa FK (pola `T_NB_OPPORTUNITY` nbfacin). Isi = medan yang
`SetValueToEDMWork` / `DataToEDM` salin dari kasus portal (`services/bukakasus.go`) dan konektor `Start2` flow
`InputAddendumFacultativeIn`. Lebar = keputusan agent; panjang diperiksa di Go dalam BYTE sebelum menulis.

| Kolom | Tipe | Tipe DDL (migrasi 262) | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(32) NOT NULL | PK = `T_WORK_POLIS.ID` (`EDM-00001`…) |
| `OLD_POLICY_NO` | teks | VARCHAR2(100) | `OfferFacIn.QuotationData.OldPolicyNo` = `.PolicyNo` portal (langkah 12) |
| `EDM_TYPE` | teks | VARCHAR2(10) | `.Quotation.EdmType` (K-029: 1 / 2 / 4) |
| `EDM_TYPE_NEW` | teks | VARCHAR2(10) | `.Quotation.EdmTypeNew` |
| `DETAIL_TYPE` | teks | VARCHAR2(10) | `.Quotation.Type` — Detail EDM Type, daftar `SetEdmType` |
| `ENDORSEMENT_DATE` | DATE | DATE | tanggal endorsement, tanggal kalender WIB |
| `ENDORSEMENT_REASON` | teks | VARCHAR2(4000) | `.Note` → `OfferFacIn.EndorsmentReason` (langkah 13) |
| `ENDORSEMENT_SOURCE` | teks | VARCHAR2(255) | `.EndorsementSource` |
| `ENDORSEMENT_SOURCE_NOTE` | teks | VARCHAR2(1000) | `.EndorsementSourceNote` (wajib di layar) |
| `ENDORSEMENT_STATUS` | teks | VARCHAR2(255) | `.EndorsementStatus` |
| `ADMIN_CHARGE` | teks | VARCHAR2(5) | `.AdminCharge`, checkbox `true` / `false` |
| `SURVEY` | teks | VARCHAR2(5) | `.Survey`, checkbox `true` / `false` |
| `INTERNAL_RETRO` | teks | VARCHAR2(10) | `.EndorsementInternalRetro` |
| `NB_STATUS` | teks | VARCHAR2(255) | `NEW EDM` (konektor `Start2`) |
| `IS_CEDING_CONFIRM` | teks | VARCHAR2(20) | `Policy` (konektor `Start2`) |

### T_WORK_POLIS (milik premiumlistlife, ditulis bersama)

Satu baris per kasus EDM, dalam transaksi yang sama dengan `T_EDM_FACIN_KASUS` (persetujuan work owner
11-10-2026): `LINI = 'FACEDM'` (BUKAN `FAC` — NB FacIn menyaring lininya hanya lewat kolom ini),
`FLAG_ONGOING_POLICY = '1'`, `STATUS_WORK = 'Pending-Policy'`, `CREATE_OP` = `CREATE_OP_NAME` = pengenal akun
pembuat, `TGL_CREATE`/`TGL_UPDATE` = SYSDATE. Tidak dibuat modul ini.

Posisi (keputusan work owner W-EDM-18, mengganti W-EDM-11; pola NB tiket 77):

- `POSITION = 'Policy'` (tahap) dan `POSITION_NOTE = 'ReasFacInMarketing'` (workbasket). Kolom `POSITION_NOTE`
  VARCHAR2(255) dibuat migrasi **nbfacin 254**.
- Selama 254 belum dijalankan, biner menulis bentuk lama (`POSITION = 'ReasFacInMarketing'`, tanpa `POSITION_NOTE`) —
  deteksi kolom toleran, `repository/posisinote.go`. Migrasi `263_posisi_kasus_edm.sql` mengonversi baris bentuk lama
  sesudah 254 (urutan nama berkas lintas modul).

## Tabel warisan yang dibaca, tidak dibuat

Judul di bawah sengaja berskema. Nol DDL; lebar dari `D:\migrasi\RNM\DDL\` bila ada.

### POOLDATA.JSON_POLIS

`GetEDMStatus_SQL` (gerbang 1, `DATA_JSON.QuotationData.EdmType`), `GetBusinessType_Sql` (langkah 7,
`DATA_JSON.QuotationData.BusinessType`) — baris `NOPOLIS` dengan `PRODKE = COUNT(NOPOLIS) - 1`;
`GetListRNWbyNopolis_SQL` (gerbang 5, `OLDNOPOLIS`). Sudah dinyatakan tabel warisan oleh `nbfacin`.

### POOLDATA.OS_AKSEPTASI_KLAIM

`GetDataClaim_SQL` (gerbang 3): `STS_REJECT` baris `NOPOLIS` terakhir menurut `TANGGAL`.

### POOLDATA.FACOUTPRODUCTION

`GetFacoutList_SQL` (gerbang 6, hanya EDM Internal Retro): ada baris `POLICYNO`.

### POOLDATA.FACINPRODUCTION

`GetStartDate` (`CheckEDMPolisDate` langkah 3): `BEGINDATE`/`ENDDATE` (DATE) baris `SPREAD_DATE` terakhir.

### POOLDATA.BUSINESS

Jenis bisnis polis NB sistem baru (DATA_JSON minimal): `FACINPRODUCTION.BUSINESSCODE` → `BUSINESS.ID` →
`GROUPPANEL`, `OLDID` (saringan `GROUPPANEL IS NOT NULL`, pola nbfacin), lalu kontrak `JenisBisnisFacIn`
(BusinessType_DeT milik nbfacin; keputusan work owner 12-10-2026 opsi 1a). Semua kolom VARCHAR2
(`D:\migrasi\RNM\DDL\BUSINESS.txt`).

### POOLDATA.OPENPROTEKSI_EDM

`BrowseOpenProteksiEdm_RD` (klep gerbang 1/3/4/5, `Param.Type` 1 Batal / 2 Klaim / 3 Pembayaran / 4 Renewal):
ada baris `POLICY_NO` + `TYPE` dengan `STS_AKSEP` bukan `'2'` atau NULL. DDL dan XML report kiriman work owner
11-10-2026 (`D:\migrasi\RNM\DDL\OPENPROTEKSI_EDM.txt`, `BrowseOpenProteksiEdm_RD.xml`); saringan F1 `IDPEGA` tidak
diport (parameternya tidak pernah diisi pemanggil). Tabel dipakai bersama Claim Prop / Claim Fac In (`TYPE '5'`).

### ARASAPAS.DETAIL_INVOICE

`SearcStatusBayarArasaps_SQL` (gerbang 4, hanya EdmType 1/2): ada baris `INV_INV_NO` = nomor polis tanpa
titik dan `IVD_JR_ID = '5'`. Skema lain — tidak terbaca (ORA-00942/01031) = gerbang `belum-diperiksa`.
