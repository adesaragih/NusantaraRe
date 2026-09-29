# Struktur Tabel — Treaty Contract Out (master arrangement kontrak treaty non-life)

Acuan bentuk tabel untuk aplikasi Go. Dibuat 2026-09-28 di tiket 01 (PREFACTOR).
Presisi fisik adalah keputusan modul ini (NUMBER(38,8), pola PremiumList Life) dan menunggu pencocokan DBA.
Berkas ini menggambarkan BENTUK, bukan alasan — alasannya ada di `spec.md` dan `issues/01`.

⚠️ **Berkas ini adalah acuan TUNGGAL nama kolom.** Penjaga `TestKolomDDLCocokDenganStruktur` membandingkan
setiap tabel di sini dengan DDL migrasi 300–306 dua arah.

**Keputusan tco1 `[DIPUTUSKAN; veto work owner]`** — nama tabel berawalan **`T_`**. Keenam nama warisan
(`TREATYYEAR`, `TREATYCONTRACT`, `TREATYREINSURER`, `MTREATYSECURITY`, `TREATYBUSINESS`, `PROPORTIONALARRG`)
**sudah dipakai** tabel hidup di skema `POOLDATA` yang sama. Tabel warisan **tidak disentuh**: tidak ditulis,
tidak di-`ALTER`, hanya dibaca oleh skrip migrasi data (tiket 01) sebagai sumber.

**Nama kolom VERBATIM dari warisan** — termasuk yang dibaca hilir (`Claim Prop`, `Komite Claim Prop`,
`Claim Fac In`; spec b224–b232, kontrak `tco_kontrak_hilir.go`). Yang berubah hanya TIPE (uang/persen → angka
desimal, tanggal → DATE) dan satu PK surrogate baru di `T_MTREATYSECURITY`.

**Tipe ditulis sebagai kategori logis:** teks · angka desimal · DATE · timestamp.
Pemetaan fisik: teks → `VARCHAR2(255)` (pengenal `*ID` → `VARCHAR2(32)`; teks bebas panjang → `VARCHAR2(1000)`);
angka desimal → `NUMBER(38,8)`; DATE → `DATE`; timestamp → `TIMESTAMP`.

⚠️ Seluruh kolom **nullable** kecuali PK; wajib-isi ditegakkan di Go (ADR-U-0027). NOT NULL akan menolak baris
warisan yang memang kosong saat migrasi data.

⛔ Nol `COMMIT` di teks SQL (ADR-U-0029). Nol kolom dokumen: setiap atribut arrangement adalah kolom bernama
(penyimpangan sadar 1 — tabel `M_*` warisan MATI, tidak dibaca migrasi).

Rule korpus yang dirujuk (folder `D:\XML\RNM_BRD\Treaty Contract Out\`):

| Singkatan | Rule | Peran |
| --- | --- | --- |
| **SaveYear** | `RDBList/SaveMasterTreatyYear_SQL.xml` → `POOLDATA.PEGA_TREATYYEAR` (10 param + 2 out) | penulis `TREATYYEAR` |
| **SaveContract** | `RDBList/SaveMasterTreatyContract_SQL.xml` → `POOLDATA.PEGA_TREATYCONTRACT` (8 + 2 out) | penulis `TREATYCONTRACT` |
| **SaveReins** | `RDBList/SaveMasterTreatyReinsurer_SQL.xml` → `POOLDATA.PEGA_TREATYREINSURER` (19 + 2 out) | penulis `TREATYREINSURER` |
| **SaveBiz** | `RDBList/SaveMasterTreatyBusiness_SQL.xml` → `POOLDATA.PEGA_TREATYBUSINESS` (12 + 2 out) | penulis `TREATYBUSINESS` |
| **SaveArrg** | `RDBList/SaveMasterProportionalArrg.xml` → `POOLDATA.PEGA_PROPORTIONALARRG` (35 + 2 out) | penulis `PROPORTIONALARRG` induk |
| **SaveArrgChild** | `RDBList/SaveMasterProportionalArrgChild.xml` → `POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD` (26 + 2 out) | penulis `PROPORTIONALARRG` anak — tabel yang SAMA `[data DBA]` |
| **InsSec** | `RDBList/InsertToMTreatySecurity.xml` — `insert into mtreatysecurity values (…)` posisional 7 nilai | penulis `MTREATYSECURITY` |
| **DDL-DBA** | `dba-procedures.md` bab "DDL + 4 procedure master lain" `[data DBA]` | tipe warisan |

---

## T_TREATYYEAR

Tahun treaty — wadah seluruh kontrak satu tahun per grup. Sumber warisan `POOLDATA.TREATYYEAR`
(`[data DBA]` seluruh kolom `VARCHAR2`).

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | SaveYear param 1; warisan `'1' \|\| lpad(TreatyYear_seq.nextval, 6, '0')` |
| `TREATYYEAR` | teks | ya | | SaveYear param 2 — kode tahun, tetap teks (ADR-U-0022) |
| `UNDERWRITINGYEAR` | teks | ya | | SaveYear param 3 |
| `TREATYGROUPID` | teks | ya | | SaveYear param 4 — kunci gabungan anak |
| `TREATYGROUPNAME` | teks | ya | | SaveYear param 5 |
| `USERID` | teks | ya | | SaveYear param 6 (`OperatorID.pyUserName`) |
| `TGLUPDATE` | DATE | ya | | SaveYear param 7; warisan `VARCHAR2` → DATE (AC 53) |
| `PROPORTION` | teks | ya | | SaveYear param 8 — `[terbuka]` arti; layar `InputDtlTreatyContact.xml` b6829 mengisinya dari pilihan "Reinsurance Type" `.ID`; dibawa apa adanya sebagai teks |
| `STARTDATE` | DATE | ya | | SaveYear param 9; warisan `VARCHAR2` → DATE (AC 53) |
| `ENDDATE` | DATE | ya | | SaveYear param 10; warisan `VARCHAR2` → DATE (AC 53) |

**Sequence:** `SEQ_T_TREATYYEAR` — identitas `'1' + lpad(6)`; nilai berjalan diselaraskan migrasi data.

**Index:** tidak ada di luar PK. Anti-dobel `(STARTDATE, ENDDATE, TREATYGROUPID)` (AC 73) ditegakkan di Go,
**bukan** unique index: data warisan boleh sudah memuat duplikat.

---

## T_TREATYCONTRACT

Kontrak treaty (jenis reasuransi) di dalam satu tahun treaty. Sumber warisan `POOLDATA.TREATYCONTRACT`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | SaveContract param 1; warisan `'1' \|\| lpad(treatycontract_seq.nextval, 6, '0')` |
| `IDTREATYYEAR` | teks | ya | FK → `T_TREATYYEAR.ID` | SaveContract param 2 |
| `REINSTYPEID` | teks | ya | | SaveContract param 3 — dari master `REINSURANCETYPE` (tiket 02) |
| `REINSTYPENAME` | teks | ya | | SaveContract param 4 |
| `TREATYSTARTDATE` | DATE | ya | | SaveContract param 5; warisan sudah `DATE` |
| `TREATYENDDATE` | DATE | ya | | SaveContract param 6; warisan sudah `DATE` |
| `USERID` | teks | ya | | SaveContract param 7 (`OperatorID.pyUserName`) |
| `TGLUPDATE` | DATE | ya | | SaveContract param 8; warisan `VARCHAR2(1000)` → DATE (AC 53) |

**Sequence:** `SEQ_T_TREATYCONTRACT`.

**Index:** `IDX_T_TREATYCONTRACT_TAHUN (IDTREATYYEAR)`.

**Relasi:** FK ke tahun **tanpa** `ON DELETE CASCADE` — menghapus tahun yang masih berkontrak ditolak Oracle
(ORA-02292), gagal terang. Tidak ada jalur hapus tahun di tiket mana pun.

---

## T_TREATYREINSURER

Reinsurer pada kombinasi **(TREATYYEAR, TREATYGROUPID, REINSTYPEID)** — BUKAN FK ke kontrak
(`[fakta bisnis — work owner]`, spec §2). Sumber warisan `POOLDATA.TREATYREINSURER`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | SaveReins param 1; warisan `'1' \|\| lpad(M_TREATYREINSURER_SEQ.nextval, 6, '0')` |
| `TREATYYEAR` | teks | ya | | SaveReins param 2 — kunci gabungan; dibaca hilir `GetListRetro_Sql` |
| `TREATYGROUPID` | teks | ya | | SaveReins param 3 — kunci gabungan; dibaca hilir |
| `TREATYGROUPNAME` | teks | ya | | SaveReins param 4 |
| `REINSTYPEID` | teks | ya | | SaveReins param 5 — kunci gabungan; dibaca hilir |
| `REINSTYPENAME` | teks | ya | | SaveReins param 6 |
| `REINSURERID` | teks | ya | | SaveReins param 7; dibaca hilir |
| `CLIENTID` | teks | ya | | SaveReins param 8; dibaca hilir |
| `NAME` | teks | ya | | SaveReins param 9 — nama perusahaan reinsurer; dibaca hilir |
| `RICOMM` | angka desimal | ya | | SaveReins param 10; warisan sudah `NUMBER`; dibaca hilir |
| `PCTSHARE` | angka desimal | ya | | SaveReins param 11; warisan sudah `NUMBER`; dibaca hilir |
| `IUDATE` | teks | ya | | SaveReins param 12 — `[terbuka]` arti dan bentuk; tidak termasuk daftar tanggal AC 53, dibawa apa adanya |
| `USERID` | teks | ya | | SaveReins param 13 |
| `STARTDATE` | DATE | ya | | SaveReins param 14; warisan `VARCHAR2` → DATE (AC 53) |
| `ENDDATE` | DATE | ya | | SaveReins param 15; warisan `VARCHAR2` → DATE (AC 53) |
| `STATUSON` | teks | ya | | SaveReins param 16 |
| `STDRATING` | teks | ya | | SaveReins param 17 — field dipakai-ulang, isi belum terverifikasi; dibawa apa adanya |
| `OPERATORNAME` | teks | ya | | SaveReins param 18 (`OperatorID.pyUserName`) |
| `TGLUPDATE` | DATE | ya | | SaveReins param 19; warisan `VARCHAR2` → DATE (AC 53) |

**Sequence:** `SEQ_T_TREATYREINSURER`.

**Index:** `IDX_T_TREATYREINSURER_KOMB (TREATYYEAR, TREATYGROUPID, REINSTYPEID)` — kunci gabungan yang dipakai
`GetMasterReinsurerList`, kaskade hapus, dan hilir.

---

## T_MTREATYSECURITY

Security di bawah seorang reinsurer. Sumber warisan `POOLDATA.MTREATYSECURITY` (`[data DBA]` **tanpa PK**,
INSERT posisional). Penyimpangan sadar 5: PK surrogate, kolom bernama, `PCT_SHARE` desimal, `REAS_SECURITY`
atribut biasa.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | **baru** — surrogate dari `SEQ_T_MTREATYSECURITY`, diberikan saat migrasi (AC 68) |
| `THN_TREATY` | teks | ya | | InsSec posisi 1 (`InputData.CARI8` = `.TreatyYear` baris reinsurer, `ViewDetailTreatyReinsurerGrid1.xml` b5301) |
| `TOP_ID` | teks | ya | | InsSec posisi 2 — dikosongkan `''`; `[terbuka]` arti, dibawa sebagai kolom bernama |
| `TP_TREATY` | teks | ya | | InsSec posisi 3 — dikosongkan `''`; `[terbuka]` |
| `REAS_ID` | teks | ya | FK → `T_TREATYREINSURER.ID` `ON DELETE CASCADE` | InsSec posisi 4 (`InputData.CARI9` = `.ID` baris reinsurer, b5307) |
| `PCT_SHARE` | angka desimal | ya | | InsSec posisi 5; warisan `VARCHAR2(99)` → desimal (AC 51) |
| `USER_ID` | teks | ya | | InsSec posisi 6 — dikosongkan `''`; `[terbuka]` |
| `REAS_SECURITY` | teks | ya | | InsSec posisi 7 — nama security, **bukan** bagian kunci (AC 18) |

**Sequence:** `SEQ_T_MTREATYSECURITY`.

**Index:** `IDX_T_MTREATYSECURITY_REAS (REAS_ID)`.

**Relasi:** satu-satunya FK berkaskade di modul ini — security memang milik reinsurernya
(`DeleteFromTreatyReinsurer_Act.xml` menghapus security lalu reinsurer, tiket 06).

---

## T_TREATYBUSINESS

Jenis bisnis yang ditanggung pada kombinasi (TREATYYEAR, TREATYGROUPID, REINSTYPEID). Sumber warisan
`POOLDATA.TREATYBUSINESS` (`[data DBA]` seluruh kolom `VARCHAR2`).

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | SaveBiz param 1; warisan `'1' \|\| lpad(TREATY_BUSINESS_SEQ.nextval, 6, '0')`; dibaca hilir `GetTreatyGroup_Sql` |
| `ISACTIVE` | teks | ya | | SaveBiz param 2 — hilir menyaring `isactive='1'` (`Claim Prop/RDBList/GetTreatyGroupID.xml`) |
| `TREATYYEAR` | teks | ya | | SaveBiz param 3; dibaca hilir |
| `TREATYYEARID` | teks | ya | | SaveBiz param 4 — **boleh NULL di data lama** (`DeleteFromTREATYCONTRACT_SQL.xml`: `OR TREATYYEARID IS NULL`); tanpa FK |
| `TREATYGROUPID` | teks | ya | | SaveBiz param 5; dibaca hilir |
| `TREATYGROUPNAME` | teks | ya | | SaveBiz param 6 |
| `REINSTYPEID` | teks | ya | | SaveBiz param 7; dibaca hilir |
| `REINSTYPENAME` | teks | ya | | SaveBiz param 8 |
| `BIZCODE` | teks | ya | | SaveBiz param 9; dibaca hilir |
| `BIZNAME` | teks | ya | | SaveBiz param 10; dibaca hilir |
| `USERID` | teks | ya | | SaveBiz param 11 |
| `TGLUPDATE` | DATE | ya | | SaveBiz param 12; warisan `VARCHAR2` → DATE (AC 53) |

**Sequence:** `SEQ_T_TREATYBUSINESS`.

**Index:** `IDX_T_TREATYBUSINESS_KOMB (TREATYYEAR, TREATYGROUPID, REINSTYPEID)`.

---

## T_PROPORTIONALARRG

SATU tabel untuk 25 jenis klausul (penyimpangan sadar 2), dibedakan `TREATYDESCID`. Sumber warisan
`POOLDATA.PROPORTIONALARRG`. Baris "anak" (7 jenis, SaveArrgChild) mengisi NULL pada sembilan kolom khusus induk:
`ID_OCCUPATION`, `OCCUPATION`, `ID_CLAUSE`, `CLAUSE`, `TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`, `MORERP`, `MOREUSD`.

Kolom `PROPORTIONALLIST` dan `OBJECT` warisan **tidak dibawa** (AC 70): tidak di-set procedure mana pun;
migrasi data melaporkan bila ada isi hidup di sana.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | SaveArrg param 1; warisan `'1' \|\| lpad(PROPORTIONALARRG_SEQ.nextval, 7, '0')` — **7 digit** |
| `TREATYYEAR` | teks | ya | | SaveArrg param 2 — kunci gabungan; dibaca hilir |
| `TREATYYEARID` | teks | ya | | SaveArrg param 3; dibaca hilir `GetDataTreatyLimit_Sql` |
| `TREATYGROUPID` | teks | ya | | SaveArrg param 4 — kunci gabungan; dibaca hilir |
| `TREATYGROUPNAME` | teks | ya | | SaveArrg param 5 |
| `TREATYDESCID` | teks | ya | | SaveArrg param 6 — jenis klausul dari master `TREATYDESC`; dibaca hilir |
| `TREATYDESCNAME` | teks | ya | | SaveArrg param 7 |
| `REINSTYPEID` | teks | ya | | SaveArrg param 8 — kunci gabungan; dibaca hilir |
| `REINSTYPENAME` | teks | ya | | SaveArrg param 9; dibaca hilir |
| `LAYER` | teks | ya | | SaveArrg param 10 — istilah apa adanya (spec §4) |
| `LAYERPART` | teks | ya | | SaveArrg param 11 |
| `LAYERPARTTYPE` | teks | ya | | SaveArrg param 12 |
| `LAYERTYPE` | teks | ya | | SaveArrg param 13 |
| `KURS` | angka desimal | ya | | SaveArrg param 14 — nilai kurs, desimal (tiket 11 AC); diisi kurs yang dipakai menghitung `Usd` tujuh induk berkurs (warisan praktis selalu kosong, OQ-TCO-18) |
| `TGLUPDATE` | DATE | ya | | SaveArrg param 15; warisan sudah `DATE` |
| `USERID` | teks | ya | | SaveArrg param 16 |
| `LINE` | teks | ya | | SaveArrg param 17 |
| `PCT` | angka desimal | ya | | SaveArrg param 18; warisan `VARCHAR2(1000)` → desimal (AC 52); dibaca hilir |
| `PCTME` | angka desimal | ya | | SaveArrg param 19; warisan `VARCHAR2(1000)` → desimal (AC 52) |
| `YDCF` | teks | ya | | SaveArrg param 20 |
| `METHOD` | teks | ya | | SaveArrg param 21 |
| `TERRITORIALLIMIT` | teks | ya | | SaveArrg param 22 |
| `PARENTREINSTYPEID` | teks | ya | | SaveArrg param 23 — sentinel `"00"` = tanpa induk (`Activity/GetPeriode.xml` b889); dibaca hilir `GetQuotaShare` |
| `SPREADINGORDER` | teks | ya | | SaveArrg param 24 |
| `RP` | angka desimal | ya | | SaveArrg param 25; warisan `VARCHAR2(1000)` → desimal (AC 51); dibaca hilir |
| `USD` | angka desimal | ya | | SaveArrg param 26; warisan `VARCHAR2(1000)` → desimal (AC 51); dibaca hilir |
| `ID_OCCUPATION` | teks | ya | | SaveArrg param 27 — khusus induk |
| `OCCUPATION` | teks | ya | | SaveArrg param 28 — khusus induk |
| `ID_CLAUSE` | teks | ya | | SaveArrg param 29 — khusus induk |
| `CLAUSE` | teks | ya | | SaveArrg param 30 — khusus induk |
| `TREATYLIMIT` | angka desimal | ya | | SaveArrg param 31 — khusus induk; warisan sudah `NUMBER` |
| `COINS_MIN` | angka desimal | ya | | SaveArrg param 32 — khusus induk; warisan sudah `NUMBER` |
| `COINS_MAX` | angka desimal | ya | | SaveArrg param 33 — khusus induk; warisan sudah `NUMBER` |
| `MORERP` | angka desimal | ya | | SaveArrg param 34 — khusus induk; warisan sudah `NUMBER` |
| `MOREUSD` | angka desimal | ya | | SaveArrg param 35 — khusus induk; warisan sudah `NUMBER` |

**Sequence:** `SEQ_T_PROPORTIONALARRG` — 7 digit.

**Index:** `IDX_T_PROPARRG_KOMB (TREATYYEAR, TREATYGROUPID, REINSTYPEID)` dan
`IDX_T_PROPARRG_DESC (TREATYYEARID, TREATYGROUPID, TREATYDESCID, PARENTREINSTYPEID)` — bentuk `WHERE`
`GetMasterDescriptionLimitParentList` dan `TreatyTestChildTotal_Act`.

**Relasi:** TIDAK ada FK ke kontrak — klausul milik level tahun/grup/jenis, dan kaskade hapus kontrak
(tiket 10) **tidak** menyentuhnya.

---

## T_TREATYCO_JEJAK

Jejak audit modul ini (ADR-0007; AC 41 spec; tiket 10 "siapa, kapan, berapa baris tiap jenis"). Tabel BARU
tanpa padanan warisan: `USERID`/`TGLUPDATE` pada tiap baris hanya menyimpan penulis terakhir, dan baris yang
dihapus tidak dapat menyimpan siapa yang menghapusnya.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | `SEQ_T_TREATYCO_JEJAK` |
| `WAKTU` | timestamp | ya | | kapan (ADR-0007) — TIMESTAMP, dua tindakan dalam detik yang sama tetap terbedakan |
| `AKUN_ID` | teks | ya | | siapa — pengenal akun pelaku, nol nama orang |
| `TABEL` | teks | ya | | tabel yang disentuh (`T_TREATYYEAR`, … ) |
| `BARIS_ID` | teks | ya | | `ID` baris yang disentuh |
| `AKSI` | teks | ya | | `simpan` / `hapus` |
| `KETERANGAN` | teks | ya | | rincian, mis. cacah baris tiap jenis pada kaskade hapus (tiket 10) |

**Sequence:** `SEQ_T_TREATYCO_JEJAK`.

**Index:** `IDX_T_TREATYCO_JEJAK_BARIS (TABEL, BARIS_ID)`.


## T_TREATYYEAR_LAMPIRAN

Lampiran berkas pada tahun treaty — **fitur baru** tiket 12 (penyimpangan sadar 9), ditambahkan 29-09-2026. Tanpa
padanan warisan: `M_ATTACHMENTTREATY_2` berkunci ID treaty inward (`GetAllAttachment2_Sql.xml`) dan tetap milik konteks
itu. Nama kolom mengikuti tabel warisan di mana maknanya sama.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | `SEQ_T_TREATYYEAR_LAMPIRAN` (lebar 9) |
| `IDTREATYYEAR` | teks | tidak | FK → `T_TREATYYEAR.ID` tanpa kaskade | tahun treaty induk (bukan `TREATYID` inward) |
| `FILENAME` | teks | ya | | nama berkas unggahan (`GetAllAttachment2_Sql` `filename`) |
| `FILEMIMETYPE` | teks | ya | | jenis berkas (`FILEMIMETYPE`) |
| `CATEGORY` | teks | ya | | teks `NOTE` master `CATEGORY_ATTACH_REAS` (`SetCategoryAttachTreatyin.xml` b500) |
| `IMAGEID` | teks | tidak | unik | kunci berkas di penyimpanan; acak, lahir bersama baris, tidak pernah berubah |
| `T_STORAGE_ID` | teks | ya | | terisi sesudah penyimpanan memastikan berkasnya ada; kosong = tertunda |
| `UKURAN` | angka bulat | ya | | byte berkas |
| `USERID` | teks | ya | | pengenal akun pengunggah — nol nama orang |
| `TGLUPLOAD` | DATE | ya | | waktu unggah |

**Sequence:** `SEQ_T_TREATYYEAR_LAMPIRAN`.

**Index:** `IDX_T_TYLAMPIRAN_TAHUN (IDTREATYYEAR)`; unik `UQ_T_TYLAMPIRAN_IMAGEID (IMAGEID)`.
---

### Catatan tiket 01 — tabel yang TIDAK dibuat modul ini

| Objek warisan | Perlakuan |
| --- | --- |
| `REINSURANCETYPE`, `TREATYDESC`, `TREATYGROUP`, `TREATYEXCHANGEYEARLY`, `CATEGORY_ATTACH_REAS`, master mata uang | master **dibaca saja** dari `POOLDATA` (spec b107); tidak dibuat, tidak ditulis |
| `M_PROPORTIONALARRG`, `M_TREATYCONTRACT`, `M_TREATYYEAR`, `M_TREATYBUSINESS` | tabel dokumen MATI — tidak dibaca migrasi, tidak dibuat (AC 63, 64) |
| `TREATYYEAR`, `TREATYCONTRACT`, `TREATYREINSURER`, `MTREATYSECURITY`, `TREATYBUSINESS`, `PROPORTIONALARRG` | sumber migrasi data — dibaca sekali oleh `tco_migrasidata.go`, tidak pernah ditulis |
