# Struktur Tabel — Komite Claim Fac In

Acuan bentuk tabel modul `komiteclaimfacin` (dibangun 10-10-2026, prompt work owner tahap 2, pola Komite Claim Prop / Non
Prop).

⛔ **Modul ini tidak membuat tabel.** Migrasinya (rentang 640–679) hanya mengubah baris roster warisan dan tabel bersama
`T_GENERAL_KOMITE`. Definisi kolom tabel bersama tinggal di dokumen pemiliknya (STRUKTUR Claim Life, Komite Claim Life,
Komite Claim Prop — bab `T_GENERAL_KOMITE` sudah memuat kolom yang ditambahkan modul ini, lampiran 10-10-2026). Judul bab
di berkas ini sengaja bukan nama tabel telanjang supaya penjaga STRUKTUR tidak membacanya sebagai pemilik kolom kedua.

Rule korpus yang dirujuk:

| Singkatan | Rule |
| --- | --- |
| **ApprovalKomite** | `Komite Claim FacIn/Activity/ApprovalKomite_Act.xml` |
| **SetValueKomite** | `Komite Claim FacIn/Activity/SetValueKomite.xml` |
| **KomitePost** | `Komite Claim FacIn/Activity/KomitePost_Adjustment.xml`, `KomitePost_Reject.xml`, `KomitePost_CloseClaim.xml` |
| **KomiteRouter** | `Komite Claim FacIn/Activity/KomiteRouter.xml` |
| **ShowTransfer** | `Komite Claim FacIn/Section/ShowTransfer.xml` |
| **CreateKMTNo** | `Claim Fac In/Activity/CreateKMTNo_Act.xml` (TT2) |
| **SendReject / SendClose** | `Claim Fac In/Activity/SendRejectClaimToKomite2.xml` (TT3), `SendCloseClaimToKomite.xml` (TT4) |

---

## Migrasi modul ini

| No | Isi | Mundur |
| --- | --- | --- |
| 640 | `EMAILKOMITE` STS_KLAIM FACIN: `OPERATOR_ID` / `NAME` = workbasket per DEGREE + JABATAN (1 `ReasClaimSPVA`, 2 `ReasClaimDeptHead`, 3 `ReasClaimTechDivHead`, 4 `ReasClaimOpsDir`, 5 `ReasClaimTechDir`), `EMAIL` kosong; `M_WORKBASKET` `ReasClaimSPVA` / `ReasClaimSPVB` disisipkan bila belum ada (keputusan work owner KCF-01) | baris FACIN ber-workbasket dikosongkan; workbasket tidak dibuang (sudah ada di DEV sebelum 640) |
| 641 | `T_GENERAL_KOMITE.ADJUSTMENT_ID` boleh kosong (`MODIFY ... NULL`, KCF-03) | `NOT NULL NOVALIDATE` — baris kosong yang ada dibiarkan |
| 642 | `T_GENERAL_KOMITE.TRANSFER_TYPE CHAR(1) DEFAULT '2' NOT NULL` + `CK_GENERAL_KOMITE_TRANSFER` (`'2'`/`'3'`/`'4'`, KCF-03) | kolom + CHECK dibuang |
| 643 | `T_GENERAL_KOMITE.KOMITE_CIRCUM_CAUSE_OF_LOSS` / `KOMITE_EXTENT_OF_LOSS` / `KOMITE_LEGAL_LIABILITY` `VARCHAR2(4000)` nullable — teks pop-up TT3 / TT4 (keputusan work owner OQ-KCFI-03) | tiga kolom dibuang |

Urutan: 640+ berjalan sesudah claimprop 537 (empat workbasket jenjang 2–5) dan claimfacin 560–567 (kolom FAC tabel klaim).

## Kasus komite — T_WORK_CLAIM, T_GENERAL_KOMITE, T_KOMITE_KOMITELIST (dipakai bersama)

Yang mengikat baris FACIN:

| Aturan | Isi | Sumber |
| --- | --- | --- |
| `T_WORK_CLAIM.ID` | berawalan `KMT-` (awalan Pega, OQ-CFI-02) | CreateKMTNo / SendReject / SendClose (Claim Fac In) |
| `T_WORK_CLAIM.COVER_KEY` | ID klaim induk berawalan `CLM-` | idem |
| `T_WORK_CLAIM.LINI` | `FACIN` — setiap kueri menyaring KETAT (fixture Claim Life memakai `KMT-` ber-LINI LIFE) | pola Komite Claim Prop |
| `T_WORK_CLAIM.TAHAP` | `Komite_Flow` | Claim Fac In `models.TahapKomite` |
| `T_WORK_CLAIM.POSITION` | KomiteID tingkat berjalan (workbasket, migrasi 640); kosong sesudah selesai | KomiteRouter (perbaikan §5 butir 1) |
| `T_GENERAL_KOMITE.ADJUSTMENT_ID` | TT2: baris `T_CLAIM_ADJUSTMENT` yang diputus (`KOMITE_ID` UNIK di tabel adjustment, OQ-CFI-28); TT3 / TT4: kosong | KCF-03 |
| `T_GENERAL_KOMITE.TRANSFER_TYPE` | `'2'` adjustment, `'3'` Reject Claim, `'4'` Close Without Payment | KomitePostAct S2 / S4 / S5 |
| `T_GENERAL_KOMITE.KOMITE_LOOP` | cacah tangga; tingkat 1 Submit menulis ulang sesudah perluasan | ApprovalKomite S8 (KCF-02) |
| `T_GENERAL_KOMITE.KOMITE_COUNT` | tingkat berjalan; tolak = `KOMITE_LOOP`, lalu +1 setiap Submit (akhir = loop + 1) | KomitePost_Adjustment S14 / S24, KomitePost_Reject S16 |
| `T_GENERAL_KOMITE.KOMITE_USUL_TUTUP` / `_CADANG` | isian "Propose To Close Case" / "Propose To Reserved" tingkat 1 TT2 | ShowTransfer LS45 (kolom komiteclaimprop 680) |
| `T_GENERAL_KOMITE.KOMITE_SUBJECTIVITY(_NOTE)` | tidak dipakai Fac In (bawaan `'0'`) | — |
| `T_GENERAL_KOMITE.KOMITE_CIRCUM_CAUSE_OF_LOSS` / `_EXTENT_OF_LOSS` / `_LEGAL_LIABILITY` | teks pop-up Chronology / Extent Of Loss / Policy Liability TT3 / TT4 (ditulis Claim Fac In saat KMT lahir; TT2 kosong) | SendRejectClaimToKomite2 / SendCloseClaimToKomite 7.2 (migrasi 643) |
| `T_KOMITE_KOMITELIST` | ditulis HANYA `backend/repository/tangga.go` (perluasan ApprovalKomite S7.1, keputusan KomitePost); `KOMITE_OPERATORID` = akun pemutus sesudah diputus | penjaga Claim Life `komite_statik_test.go` |

## Klaim induk — hanya lewat kontrak

Tabel klaim Fac In (`T_GENERAL_CLAIM`, `T_CLAIM_OBJECT`, `T_CLAIM_OBJECT_ITEM`, `T_CLAIM_ADJUSTMENT`, anak-anaknya,
`T_VIEW_SUGGEST`) dibaca dan ditulis HANYA lewat `kontrak.KlaimFacInKomite` (`inti/backend/kontrak/klaimfacin.go`,
disediakan claimfacin); daftar kerja membaca `T_GENERAL_CLAIM.CLAIM_NO` dan `T_CLAIM_ADJUSTMENT.ACCEPTANCE_STATUS` (pola
Komite Claim Prop). Daftar putih jalur tulis balik: lihat kontrak dan `docs/PARITAS.md`.

## Tabel warisan — ditulis, tidak dibuat

Tanpa procedure dan tanpa COMMIT (isi procedure dibaca dari ALL_SOURCE DEV 10-10-2026):

| Tabel | Kolom yang ditulis | Asal |
| --- | --- | --- |
| `OS_AKSEPTASI_KLAIM` | CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, STS_KONVERSI (hanya bila terisi, OQ-KCFI-08), STS_DLA (TGL_PROD = trigger) | KomitePost_Adjustment S8 (STS 1 / 4), KomitePost_Reject S12 → SaveReject_ACT_KMT (STS 2 per estimasi), KomitePost_CloseClaim S12 (STS 4) — `PEGA_JSON_OS_AKSEP_KLAIM` |
| `JSON_KLAIM` | MNK_NO_KLAIM, IDPEGA, TGL_INPUT, NOPOLIS, IDPROD (INSERT bila IDPEGA belum ada; tanpa DATA_JSON) | KomitePost_Adjustment S16, KomitePost_Reject S11, KomitePost_CloseClaim S10 (`PEGA_JSON_KLAIM_PNC`) |
| `MONITORING_KLAIM_LOG` | TGL_INPUT, IDPEGA, PARAMETER, JN_SERVICE "AKSEPATSI", NO_AKSEPTASI | KomitePost_Adjustment S18-S19 (bukan Fac Retro) |
| `HISTORYAKSEPTASIPEGA` | ID_PEGA, TGL_TRANSFER, STATUS, USERNAME, WORKBASKET "KLAIM", ID_KOMITE | KomitePost_Adjustment S20-S21 |
| `SUBPROGRESSCLAIM` | POSITION2 "Accepted" / "Rejected" (UPDATE `IDPEGA` = KMT) | KomitePost_Adjustment S22-S23 |
| `CLAIMREJECTED` | INSKEY, ID, INSNAME, LABEL, STATUSWORK, CREATEOPNAME, CREATEOPERATOR, OBJCLASS, UPDATEDATETIME, UPDATEOPNAME, UPDATEOPERATOR, REMARK | KomitePost_Reject S14 (TT3 disetujui) |
| `T_LOG_SERVICE_RNM` | outbox `inti/backend/outbox` (konversi, Kasir, email — hanya produksi) | KomitePost_Adjustment S17 / S25 / 7.2.1.15, KomitePost_Reject S13.10 / S15, KomitePost_CloseClaim S11.10 / S12.4 |

Dibaca: `EMAILKOMITE` (roster FACIN aktif), `REINSURANCETYPE`, `BANKACCOUNT`, `M_LOGIN_GO`, `M_LOGIN_GO_WORKBASKET`,
`M_WORKBASKET`, `POOLDATA.GETCURRENCYSTANDARD`, `KODE_PRODUKSI` + `TANGGAL_CLOSING` (lewat `inti/backend/penomor`),
`GL.F_GET_EMAIL` (hanya produksi).
