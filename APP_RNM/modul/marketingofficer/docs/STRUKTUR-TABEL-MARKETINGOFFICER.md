# Struktur Tabel — Marketing Officer: PETA TABEL WARISAN yang dipakai

**Keputusan work owner 03-10-2026:** modul `marketingofficer` mengelola (tambah dan ubah) tabel warisan
`POOLDATA.MARKETINGOFFICER` — *"buatkan di modul baru dengan nama marketingofficer, gunanya untuk insert update
table marketingofficer"* — dan **tidak membuat tabel baru** (*"aku ga mau tambah tabel baru"*). Modul ini **tidak
membuat satu tabel, kolom, index, trigger, atau sequence pun**; satu-satunya migrasinya slot menu `990`.

Berkas ini **peta**, bukan DDL. Tipe dari katalog DEV (`ALL_TAB_COLUMNS`, agregat, 03-10-2026). Kedua tabel
terdaftar `Tabel warisan: dibaca, tidak dibuat` di `MODUL.md` — migrasi mana pun yang membuatnya gagal penjaga.

Sumber perilaku (korpus `D:\XML\RNM_BRD\NB FacIn\`, sama di `RNW Fac In` dan `Endorsment Fac In`):

| Rule | Isi yang ditiru |
| --- | --- |
| `Section/InputMarketingOfficer.xml` | caption form: Code, Name Marketing, Set as a leader, Branch, Sub Branch, Leader, Active |
| `Activity/SaveMarketingOfficer_Act.xml` | leader: `ClientID2 = "LEADER"`, `MOLeader = ClientName`; anggota: `MOLeader` = `ClientName` baris `ID = ClientID2`; `TeamGroup = KanwilGroup`, `BranchDetailName = Name` cabang; tidak menyimpan bila `MOLeader` atau `ClientName` kosong |
| `RDBList/UpdateMasterMarketingOfficer.xml` | memanggil `POOLDATA.PEGA_MARKETINGOFFICER` — **tidak dipanggil** modul ini; INSERT/UPDATE-nya ditiru (`repository/mo.go`) |
| `ReportDefinition/SelectLeader_RD.xml` | dropdown Leader: `ClientID2 = "LEADER"`, aktif |
| `ReportDefinition/BrowseBranchDetail_RD.xml` | Sub Branch: `.Status = 1`, disaring `ParentID` |
| `Activity/SendEmailPolicy.xml` + `RDBList/GetMKTandLeader_SQL.xml` | `AKSES_LOGIN` = Operator ID; email polis ke MO, CC leader |

## MARKETINGOFFICER

Tabel warisan Pega, **tanpa PK, unique, maupun FK** (DEV: 74 baris, 5 index non-unik). Trigger warisan
`TRG_MARKETINGOFFICER_LOG` mencatat baris LAMA ke `MARKETINGOFFICER_LOG` setiap UPDATE (tanpa `AKSES_LOGIN`).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(100) | ya | | kunci baris (tanpa PK) | `1` + 7 digit `CURRENCY_SEQ`, seperti `PEGA_MARKETINGOFFICER`; tidak berubah |
| `CLIENTID` | VARCHAR2(100) | ya | | Marketing Code — disalin ke quotation, `GENERALOFFER`, `JSON_FOLLOWING`, `LIFEINPRODUCTION`, `T_PREMIUM_LIST` | baris baru: `M_LOGIN_GO.CONTACT_ID`, atau `CLIENTID` lama bila akun itu sudah punya baris MO; **tidak pernah diubah**. Baris lama Pega: kode kontak SFAGIS |
| `CLIENTNAME` | VARCHAR2(100) | ya | | Name Marketing; view `INBOXPENDING` | `M_LOGIN_GO.NAME` saat baris dibuat; tidak diubah |
| `CLIENTID2` | VARCHAR2(100) | ya | | leader | `LEADER` bila leader, selain itu `ID` baris leader |
| `MOLEADER` | VARCHAR2(100) | ya | | nama leader | salinan `CLIENTNAME` leader (leader: nama sendiri) |
| `MOSTATUS` | VARCHAR2(100) | ya | | Active | `1` aktif, `2` nonaktif — satu baris AKTIF per `CLIENTID` dan per `AKSES_LOGIN` |
| `BRANCHPARENT` | VARCHAR2(10) | ya | | Branch | `BRANCH.BRANCHPARENTID` Sub Branch yang dipilih |
| `BRANCHDETAILID` | VARCHAR2(100) | ya | | Sub Branch | `BRANCH.ID` |
| `BRANCHDETAILNAME` | VARCHAR2(100) | ya | | nama Sub Branch | salinan `BRANCH.NAME` saat disimpan |
| `TEAMGROUP` | VARCHAR2(100) | ya | | team group | salinan `BRANCH.KANWILGROUP` |
| `BRANCHSTATUS` | VARCHAR2(100) | ya | | — | tidak ada di form Pega; dibaca, **tidak pernah ditulis** |
| `TANGGAL` | DATE | ya | | jejak | `SYSDATE` setiap simpan |
| `USERUPDATE` | VARCHAR2(100) | ya | | jejak | `LOGIN_ID` pelaku |
| `AKSES_LOGIN` | VARCHAR2(50) | ya | | email polis (Pega `SendEmailPolicy`) | `M_LOGIN_GO.LOGIN_ID`; akun > 50 karakter ditolak. Nilai lama Pega yang tidak ada di `M_LOGIN_GO` dibiarkan dan ditandai di layar |

## BRANCH

Dibaca saja, hanya kolom di bawah. Tabel fisik kelas Pega `ASM-FW-GISFW-Int-BRANCHDETAIL` **tidak tercatat di XML**;
`BRANCH` dipilih work owner 03-10-2026 (bukti DEV: `BRANCHDETAILID` MO 65/65 ada di `BRANCH.ID`, `BRANCHPARENT`
66/66 = `BRANCHPARENTID`, `TEAMGROUP` = `KANWILGROUP` 57/65). ⚠️ **Belum dikonfirmasi DBA** — kandidat lain `BRANCH2`
(758 baris); nama tabelnya satu konstanta `repository.TabelCabang`.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(4000) | ya | | Sub Branch | — |
| `NAME` | VARCHAR2(4000) | ya | | nama Sub Branch | — |
| `BRANCHPARENTID` | VARCHAR2(4000) | ya | | Branch | — |
| `KANWILGROUP` | VARCHAR2(4000) | ya | | team group | — |
| `STATUS` | VARCHAR2(4000) | ya | | hanya `1` yang ditawarkan | `BrowseBranchDetail_RD` |

## Dibaca dari inti

`M_LOGIN_GO` (milik inti, digambarkan `inti/docs/STRUKTUR-TABEL-INTI.md`): `LOGIN_ID`, `NAME`, `CONTACT_ID`,
`EMAIL`, `IS_ACTIVE` — pilihan Login Account hanya akun aktif. Butuh migrasi inti 905 (`CONTACT_ID`).
