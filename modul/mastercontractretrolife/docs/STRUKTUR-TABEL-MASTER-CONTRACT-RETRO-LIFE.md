# Struktur Tabel — Master Contract Retro Life: PETA TABEL WARISAN yang dipakai

> Disusun 30-09-2026 (paket 10, brief `PROMPT-IMPLEMENTASI-MODUL-MASTER-CONTRACT-RETRO-LIFE.md`).

**K1 (`RALAT-DEV-30-09-2026.md`)** — modul ini **tidak membuat satu tabel, sequence, maupun constraint pun**
*(preseden tco4 Treaty Contract Out; tiket 12 dicabut)*. Ia menulis dan membaca lima tabel yang **sudah ada** di
`POOLDATA`, dengan nama tabel dan kolom **VERBATIM**. Rentang migrasi `100`–`139` sengaja kosong
(`TestMCRLNolMigrasiDiRentang`); satu-satunya migrasi modul ini adalah slot menu `958` (`UPDATE DIMIGRASI`).

Berkas ini **peta**, bukan DDL: tabel → kolom VERBATIM → tipe katalog → penulis Pega → cara kolom itu ditulis
modul ini. Penjaga yang memakainya: `TestKolomDDLCocokDenganStruktur` / `TestTabelBukanMilikKitaTidakDibuat`
(`inti/backend/penjaga/strukturkolom_test.go`) — setiap tabel di sini dinyatakan "Tabel warisan" di `MODUL.md`, jadi
migrasi mana pun yang **membuatnya** merah.

Sumber tipe: `docs/ddl-tables-from-dba.md` `[data DBA]`, dicocokkan dengan katalog DEV 30-09-2026. Kolom di katalog:
**nol** constraint P/R/U di kelima tabel (K1–K3) — keunikan `ID` dan kaskade hapus dijaga di Go, dalam satu transaksi.

⛔ Procedure **tidak dipanggil** (keputusan **o**): isinya ditiru di Go — UPSERT dikunci `ID`, ID baru
`'1' || LPAD(<sequence warisan>.NEXTVAL, 6, '0')`, `TGLUPDATE = SYSDATE`, nol `COMMIT` di teks SQL.

## TREATYYEAR_LIFE

Tahun treaty — grid halaman awal (`BrowseTreatyYear_Life_RD`). Penulis Pega `SaveMasterTreatyYear_Life_SQL` →
`INSERTTREATYYEAR_LIFE`. **Nol penghapus** (tahun abadi, AC 2).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(100) | `'1' || LPAD(TREATYYEAR_LIFE_SEQ.NEXTVAL, 6, '0')` — dari server, tidak pernah dari klien |
| `TREATYYEAR` | VARCHAR2(100) | label layar **TRANSACTION YEAR** |
| `UNDERWRITINGYEAR` | VARCHAR2(100) | `UNDERWRITING YEAR` |
| `USERID` | VARCHAR2(100) | akun pelaku |
| `TGLUPDATE` | DATE | `SYSDATE` |
| `STARTDATE` | DATE | `START DATE`; disalin ke `TREATYCONTRACT_LIFE.TREATYSTARTDATE` anak (K4) |
| `ENDDATE` | DATE | `END DATE`; disalin ke `TREATYENDDATE` anak (K4) |

## TREATYCONTRACT_LIFE

Kontrak (jenis reasuransi) dan batas proteksi di satu tahun — panel `Reins Type`. Penulis Pega
`SaveMasterTreatyContract_Life_SQL` → `INSERTTREATYCONTRACT_LIFE`; penghapus `DeleteTreatyLimit_SQL` (datar — di
sini berjenjang, K2).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(100) | `'1' || LPAD(TREATYCONTRACT_LIFE_SEQ.NEXTVAL, 6, '0')` |
| `IDTREATYYEAR` | VARCHAR2(100) | → `TREATYYEAR_LIFE.ID` (tanpa FK) |
| `REINSTYPEID` | VARCHAR2(100) | master `REINSURANCETYPE.ID` (`.Flag = 1`) |
| `REINSTYPENAME` | VARCHAR2(1000) | `REINSURANCETYPE.Note` — dibaca server dari master, bukan dari klien |
| `USERID` | VARCHAR2(100) | akun pelaku |
| `TGLUPDATE` | DATE | `SYSDATE` |
| `IDR` | NUMBER | `MAXIMUM LIMIT (IDR)` — `TO_NUMBER(:koef) / POWER(10, :skala)` (tahan NLS) |
| `USD` | NUMBER | `MAXIMUM LIMIT (USD)` — idem |
| `B_IDR` | NUMBER | `MINIMUM LIMIT (IDR)` — idem |
| `B_USD` | NUMBER | `MINIMUM LIMIT (USD)` — idem |
| `IDR_SELISIH` | NUMBER | **dihitung Go** `IDR − B_IDR` (K5: nol penulis di korpus modul) |
| `USD_SELISIH` | NUMBER | **dihitung Go** `USD − B_USD` (K5) |
| `TREATYENDDATE` | DATE | salinan `TREATYYEAR_LIFE.ENDDATE` (K4; layar baca-saja) |
| `TREATYSTARTDATE` | DATE | salinan `TREATYYEAR_LIFE.STARTDATE` (K4; layar baca-saja) |

## TREATYREINSURER_LIFE

Reinsurer satu kontrak — panel `Reinsurer List`. Penulis Pega `SaveMasterTreatyReinsurer_Life_SQL` →
`INSERTREINSURER_LIFE`; penghapus `DeleteSecurityReinsurer_SQL` (datar — di sini berjenjang ke security).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(100) | `'1' || LPAD(TREATYREINSURER_LIFE_SEQ.NEXTVAL, 6, '0')` |
| `TREATYYEARID` | VARCHAR2(100) | salinan `TREATYCONTRACT_LIFE.IDTREATYYEAR` induk |
| `TREATYCONTRACTID` | VARCHAR2(100) | → `TREATYCONTRACT_LIFE.ID` (tanpa FK) |
| `REINSTYPEID` | VARCHAR2(100) | salinan kontrak induk |
| `REINSTYPENAME` | VARCHAR2(1000) | salinan kontrak induk |
| `REINSURERNAME` | VARCHAR2(1000) | `AGENT.ClientName` — dibaca server dari master |
| `PCTSHARE` | NUMBER | `(%) SHARE`, 0..100 |
| `COMMISION` | NUMBER | label layar **(%) DISCOUNT**, 0..100 |
| `OVR_COMM` | NUMBER | `(%) OVR COMM` — tanpa gerbang (Pega tidak memeriksanya, OQ-MCRL-03) |
| `USERID` | VARCHAR2(1000) | akun pelaku |
| `TGLUPDATE` | DATE | `SYSDATE` |
| `REINSURERID` | VARCHAR2(100) | master `AGENT.ID` (`ID LIKE '%L0%'`, `STATUSACTIVE = 1`) |

## TREATYSECURITYREINSURER_LIFE

Security satu reinsurer — panel `Security Reinsurer`. Penulis Pega `SaveMasterTreatySecurityReinsurer_Life_SQL` →
`INSERTSECURITYREINSURER_LIFE`; penghapus `DeleteSecurityReinsurerLife_SQL`. DEV: 0 baris (30-09-2026).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(100) | `'1' || LPAD(TREATYSECURITYREINSURER_LIFE_SEQ.NEXTVAL, 6, '0')` |
| `TREATYYEARID` | VARCHAR2(100) | salinan reinsurer induk |
| `TREATYCONTRACTID` | VARCHAR2(100) | salinan reinsurer induk |
| `TREATYREINSURERID` | VARCHAR2(100) | → `TREATYREINSURER_LIFE.ID` (tanpa FK) |
| `REINSURERID` | VARCHAR2(100) | master `AGENT.ID` |
| `REINSURERNAME` | VARCHAR2(1000) | `AGENT.ClientName` |
| `PCTSHARE` | NUMBER | persen **dari share induk** (tiket 06); eksposur = share × share induk / 100, dihitung, tidak disimpan |
| `USERID` | VARCHAR2(1000) | akun pelaku |
| `TGLUPDATE` | DATE | `SYSDATE` |

## TREATYBUSINESS_LIFE

Business satu kontrak — panel `Business List`. Penulis Pega `SaveMasterTreatyBusiness_Life_SQL` →
`INSERTBUSINESS_LIFE` dan `SaveTreatyBusinessAll_Life_SQL` (`Copy to all Reinstype`);
penghapus `DeleteRowBusinessList`. Satu-satunya tabel ber-PK di DDL DBA (`TREATYBUSINESS_LIFE_PK`); katalog DEV
30-09-2026 nol constraint (OQ-MCRL-02).

| Kolom | Tipe katalog | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(100) | `'1' || LPAD(TREATYBUSINESS_LIFE_SEQ.NEXTVAL, 6, '0')` |
| `TREATYYEARID` | VARCHAR2(100) | salinan kontrak induk |
| `TREATYYEAR` | VARCHAR2(100) | `TREATYYEAR_LIFE.TREATYYEAR` tahun induk — juga pada baris `Copy to all Reinstype` (SQL Pega tidak menulisnya; K4: semua salinan ditulis dari induk) |
| `REINSTYPEID` | VARCHAR2(100) | salinan kontrak induk |
| `REINSTYPENAME` | VARCHAR2(1000) | salinan kontrak induk |
| `BIZCODE` | VARCHAR2(100) | master `BUSINESS.ID` (`OLDID LIKE 'L%'`) |
| `BIZNAME` | VARCHAR2(1000) | `BUSINESS.Note` |
| `USERID` | VARCHAR2(100) | akun pelaku |
| `TGLUPDATE` | DATE | `SYSDATE` |
| `TREATYCONTRACTID` | VARCHAR2(100) | → `TREATYCONTRACT_LIFE.ID` (tanpa FK) |
| `RIRATEID` | VARCHAR2(100) | ID ringkasan tabel rate — autocomplete `R/I RATE`, view `RATE_LIFE_SUMMARY.ID` (K1 keputusan work owner 01-10-2026) *(RALAT 07-10-2026: ringkasan rate kini tabel `M_RATE_LIFE_SUMMARY` berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026, `modul/riratelife/MODUL.md` RALAT R6; modul ini membacanya `SELECT ID, USEDBY`, tetap baca-saja)* |
| `RIRATE` | VARCHAR2(1000) | **nama tabel rate** (teks, Pertanyaan A terjawab — RALAT R7) |

## Master yang dibaca saja

Tidak ditulis modul ini (`repository.DaftarMasterDibacaSaja`, `TestMCRLMasterDibacaSaja`): `REINSURANCETYPE`
(dropdown `REINS TYPE`), `AGENT` (autocomplete `REINSURER NAME` / `SECURITY REINSURER NAME`), `BUSINESS`
(autocomplete `BUSINESS NAME`). Kedua sumber tabel rate (autocomplete `R/I RATE`, section `Rate List`) **tidak dibaca**
sampai work owner menyetujui sumbernya (OQ-MCRL-13) — rutenya menjawab 503 berkalimat.

> **Ralat 01-10-2026 (K1 keputusan work owner 01-10-2026, OQ-MCRL-13 + OQ-MCRL-05):** kalimat di atas tidak berlaku lagi. Kedua view rate
> dibaca **saja** dan masuk `DaftarMasterDibacaSaja`: `RATE_LIFE_SUMMARY` (`ID`, `USEDBY`) untuk autocomplete `R/I RATE`, *(RALAT 07-10-2026: ringkasan rate kini tabel `M_RATE_LIFE_SUMMARY` berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026, `modul/riratelife/MODUL.md` RALAT R6; modul ini membacanya `SELECT ID, USEDBY`, tetap baca-saja)*
> `RATE_LIFE` (`ID`, `USEDBY`, `GENDER`, `CONTRACT`, `AGE`, `RATE`, berkunci `IDUSEDBY`) untuk `Rate List`. Penjaga
> `periksaBacaSaja` menolak SQL selain SELECT ke objek mana pun di daftar itu sebelum sampai ke Oracle.
