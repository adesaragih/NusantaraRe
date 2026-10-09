# Struktur Tabel — Komite Claim Non Prop

Acuan bentuk tabel modul `komiteclaimnonprop` (dibangun 09-10-2026 atas perintah work owner, pola Komite Claim Prop).

⛔ **Modul ini tidak membuat satu tabel pun dan tidak menambah kolom.** Migrasinya hanya slot menu 988
(`UPDATE M_NAV_MENU SET DIMIGRASI = '1'`). Semua tabel di bawah milik modul lain atau tabel warisan; bentuk kolomnya
tidak didefinisikan ulang di sini — judul bab sengaja bukan nama tabel telanjang supaya penjaga STRUKTUR tidak membaca
berkas ini sebagai pemilik kolom.

Rule korpus yang dirujuk:

| Singkatan | Rule |
| --- | --- |
| **CreateChild** | `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml` |
| **KomitePost** | `Komite Claim Non Prop/Activity/KomitePostAdjustment.xml` |
| **KomiteRouter** | `Komite Claim Non Prop/Activity/KomiteRouter.xml` |
| **ShowTransfer** | `Komite Claim Non Prop/Section/ShowTransfer.xml` |

---

## Kasus komite — T_WORK_CLAIM, T_GENERAL_KOMITE, T_KOMITE_KOMITELIST (dipakai bersama)

Definisi kolom: dokumen STRUKTUR Claim Life, Komite Claim Life, dan Komite Claim Prop. Yang mengikat modul ini hanya
aturan baris NONPROP-nya:

| Aturan | Isi | Sumber |
| --- | --- | --- |
| `T_WORK_CLAIM.ID` | berawalan `KMTNP-` | CreateChild (Claim Non Prop tahap 1) |
| `T_WORK_CLAIM.COVER_KEY` | ID klaim induk berawalan `CLMNP-` | CreateChild |
| `T_WORK_CLAIM.LINI` | `NONPROP` — setiap kueri menyaring KETAT | pola Komite Claim Prop |
| `T_WORK_CLAIM.POSITION` | KomiteID tingkat berjalan (workbasket, migrasi claimnonprop 611); kosong sesudah selesai | perintah work owner 09-10-2026 |
| `T_GENERAL_KOMITE.KOMITE_USUL_TUTUP` / `_CADANG` / `KOMITE_SUBJECTIVITY` / `_NOTE` | isian tingkat 1 antar tingkat | kolom migrasi komiteclaimprop 680 / 682 |
| `T_KOMITE_KOMITELIST` | ditulis HANYA `backend/repository/tangga.go` (KomitePost S6 / S7 / S11.3) | penjaga Claim Life `komite_statik_test.go` |

## Klaim induk — hanya lewat kontrak

Tabel klaim Claim Non Prop (`T_GENERAL_CLAIM`, `T_CLAIM_ADJUSTMENT`, anak `T_CLAIM_NP_*`) dibaca dan ditulis HANYA lewat
`kontrak.KlaimTreatyNonPropKomite`; daftar kerja membaca `CLAIM_NO` dan `ACCEPTANCE_STATUS` (pola Komite Claim Prop).

## Tabel warisan — ditulis, tidak dibuat

Tanpa procedure dan tanpa COMMIT (isi procedure dibaca dari ALL_SOURCE DEV 09-10-2026):

| Tabel | Kolom yang ditulis | Asal |
| --- | --- | --- |
| `OS_AKSEPTASI_KLAIM` | CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, MASTERID | KomitePost S14.18 → InsertOSKlaimCNP S6 (`PEGA_JSON_OS_AKSEP_KLAIM`) |
| `CLAIMXOL2` | CASEID, "GrossAdjustment", "CNPReinstatement", "Currency", "KursIDR", XOL, TANGGAL | InsertOSKlaimCNP S7 → InsertXOLKlaimCNP S5 (`XOL2_AKSEP_KLAIM`) |
| `OS_AKSEPTASI_SUBJECTIVITY` | CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, STS_SUBJECTIVITY (UPDATE bila CASEID ada) | KomitePost S17 → InsertOSSubjectivityCNP S4 |
| `JSON_KLAIM` | MNK_NO_KLAIM, IDPEGA, TGL_INPUT, NOPOLIS, IDPROD (INSERT bila IDPEGA belum ada; tanpa DATA_JSON) | KomitePost S24 (`PEGA_JSON_KLAIM_PNC`) |
| `HISTORYAKSEPTASIPEGA` | ID_PEGA, TGL_TRANSFER, STATUS, USERNAME, WORKBASKET "KLAIM", ID_KOMITE | KomitePost S22-S23 |
| `T_LOG_SERVICE_RNM` | outbox `inti/backend/outbox` (konversi, Kasir, email — hanya produksi) | KomitePost S14.22 / S19 |

STS_KONVERSI / STS_DLA / CLAIMOLD tidak diisi InsertOSKlaimCNP (NULL); CLAIMOLD hanya untuk satu akun yang tertulis mati
di XML (dibuang).
