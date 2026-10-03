# Tipe kolom `OS_AKSEPTASI_KLAIM_LIFE` — katalog instance pengembangan

`[data DBA — dibaca sendiri 26 September 2026, belum dikonfirmasi DBA]`

Sumber: `ALL_TAB_COLUMNS` pada instance **pengembangan** (Oracle 12.2.0.1), skema `POOLDATA`,
dibaca oleh sesi verifikasi lewat kredensial yang diberikan work owner. Hanya **katalog** yang
dibaca — nol baris data. Tabel itu berisi 13.694 baris pada tanggal pembacaan.

Dokumen ini menutup prasyarat brief ronde 3 §1-2 dan temuan ronde 3 §3-7 (uji tabel warisan yang
mengonfirmasi dirinya sendiri): kini penggolongan kolom di `barislamakolom.go` dapat dibandingkan
dengan bentuk sungguhan, bukan dengan tebakan atas nama.

⚠️ Produksi **belum** dibaca. Bentuk di produksi diasumsikan sama; DBA yang dapat memastikannya.

## Kolom — 62, bukan 55

Rule `UPDATEOSAKSEPTASICLAIMLIFE_SQL` menulis **55** kolom. Tabelnya memuat **62**. Tujuh kolom
yang tidak pernah ditulis rule itu diberi tanda ⚠️ di kolom "Catatan".

| # | Kolom | Tipe | Panjang | Presisi | Skala | Null | Catatan |
| ---: | --- | --- | ---: | ---: | ---: | :-: | --- |
| 1 | `CASEID` | VARCHAR2 | 100 | | | Y | |
| 2 | `NO_CLAIM` | VARCHAR2 | 100 | | | Y | |
| 3 | `NO_ACCEPTATION` | VARCHAR2 | 100 | | | Y | |
| 4 | `STS_REJECT` | **NUMBER** | 22 | 38 | 0 | Y | ⛔ ditebak teks di `barislamakolom.go` |
| 5 | `ACCEPTATION_DATE` | DATE | 7 | | | Y | |
| 6 | `POLICY_NO` | VARCHAR2 | 100 | | | Y | |
| 7 | `POLICY_HOLDER` | VARCHAR2 | 1000 | | | Y | |
| 8 | `CERTIFICATE_NO` | VARCHAR2 | 100 | | | Y | |
| 9 | `NAME_OF_INSURED` | VARCHAR2 | 1000 | | | Y | |
| 10 | `SEX` | VARCHAR2 | 50 | | | Y | |
| 11 | `DOB` | DATE | 7 | | | Y | |
| 12 | `AGE` | NUMBER | 22 | 38 | 0 | Y | tiruan memakai `NUMBER(5)` |
| 13 | `PLAN` | VARCHAR2 | 255 | | | Y | |
| 14 | `BEGIN_DATE` | DATE | 7 | | | Y | |
| 15 | `LAPSE_DATE` | DATE | 7 | | | Y | |
| 16 | `EXPIRED_DATE` | DATE | 7 | | | Y | |
| 17 | `STATUS` | VARCHAR2 | 100 | | | Y | |
| 18 | `CURRENCY` | VARCHAR2 | 100 | | | Y | |
| 19 | `STS_KONVERSI` | CHAR | 1 | | | Y | ⚠️ bukan bagian 55 |
| 20 | `TGL_KONVERSI` | DATE | 7 | | | Y | ⚠️ bukan bagian 55 |
| 21 | `WPC` | **DATE** | 7 | | | Y | ⛔ ditebak teks; DDL baru `003` menulis `VARCHAR2(32)` |
| 22 | `PL_NUMBER` | VARCHAR2 | 100 | | | Y | |
| 23 | `DISEASE` | VARCHAR2 | 1000 | | | Y | |
| 24 | `ICD_CODE` | VARCHAR2 | 10 | | | Y | DDL baru `003`: `VARCHAR2(32)` |
| 25 | `NOTES` | VARCHAR2 | 1000 | | | Y | |
| 26 | `CEDINGCO` | VARCHAR2 | 100 | | | Y | |
| 27 | `CEDINGCONAME` | VARCHAR2 | 1000 | | | Y | |
| 28 | `SOB` | VARCHAR2 | 100 | | | Y | |
| 29 | `SOBNAME` | VARCHAR2 | 1000 | | | Y | |
| 30 | `BUSINESSID` | VARCHAR2 | 100 | | | Y | |
| 31 | `BUSINESSNAME` | VARCHAR2 | 100 | | | Y | |
| 32 | `KETERANGAN` | VARCHAR2 | 1000 | | | Y | |
| 33 | `EM_PERCENT` | NUMBER | 22 | | | Y | presisi arbitrer |
| 34 | `SUM_INSURED` | NUMBER | 22 | | | Y | presisi arbitrer |
| 35 | `CEDING_RETENTION` | NUMBER | 22 | | | Y | presisi arbitrer |
| 36 | `SUM_REASURED` | NUMBER | 22 | | | Y | presisi arbitrer |
| 37 | `SHARE_NUSANTARA_RE` | NUMBER | 22 | | | Y | presisi arbitrer |
| 38 | `CLAIM_AMOUNT` | NUMBER | 22 | | | Y | presisi arbitrer |
| 39 | `SHARE_RETRO` | NUMBER | 22 | | | Y | presisi arbitrer |
| 40 | `CLAIM_RETRO` | **NUMBER** | 22 | | | Y | ⛔ ditebak teks; `Simpan` menulis teks ke sini |
| 41 | `RETROID` | VARCHAR2 | 100 | | | Y | |
| 42 | `RETRONAME` | VARCHAR2 | 100 | | | Y | |
| 43 | `SECURITYREINSURERID` | VARCHAR2 | 100 | | | Y | |
| 44 | `SECURITYREINSURER` | VARCHAR2 | 100 | | | Y | |
| 45 | `TYPECEDING` | VARCHAR2 | 50 | | | Y | |
| 46 | `TYPE` | VARCHAR2 | 10 | | | Y | |
| 47 | `CONFIRMATION_DATE` | DATE | 7 | | | Y | |
| 48 | `CLAIM_RECEIVED_DATE` | DATE | 7 | | | Y | |
| 49 | `COMPLETE_DATE` | DATE | 7 | | | Y | |
| 50 | `ID` | VARCHAR2 | 100 | | | Y | nullable, tanpa PK di katalog kolom |
| 51 | `NAME_OF_BANK` | VARCHAR2 | 100 | | | Y | |
| 52 | `IDBANK` | VARCHAR2 | 100 | | | Y | |
| 53 | `ACCOUNTNO` | VARCHAR2 | 100 | | | Y | |
| 54 | `PRODUCTNAMEID` | VARCHAR2 | 100 | | | Y | |
| 55 | `PRODUCTNAME` | VARCHAR2 | 1000 | | | Y | |
| 56 | `CREATEOPNAME` | VARCHAR2 | 100 | | | Y | |
| 57 | `RETROCEDED_SHARE` | NUMBER | 22 | | | Y | presisi arbitrer |
| 58 | `LAYER_1` | VARCHAR2 | 10 | | | Y | ⚠️ bukan bagian 55 |
| 59 | `LAYER_2` | VARCHAR2 | 10 | | | Y | ⚠️ bukan bagian 55 |
| 60 | `LAYER_3` | VARCHAR2 | 10 | | | Y | ⚠️ bukan bagian 55 |
| 61 | `LAYER_4` | VARCHAR2 | 10 | | | Y | ⚠️ bukan bagian 55 |
| 62 | `INDEXLIST` | VARCHAR2 | 10 | | | Y | ⚠️ bukan bagian 55 |

## Selisih dengan tebakan `barislamakolom.go` — di dalam 55 kolom

| Kolom | Ditebak | Sungguhan | Akibat hari ini |
| --- | --- | --- | --- |
| `STS_REJECT` | teks | `NUMBER(38,0)` | dibaca tanpa `TO_CHAR`; driver mengubah angka ke teks — masih benar untuk `0`/`1`, tetapi tiruan bertipe `VARCHAR2(255)` tidak menguji jalur ini |
| `CLAIM_RETRO` | teks | `NUMBER` | ⛔ `Simpan` menulis teks apa adanya ke kolom angka: `ORA-01722` bila isinya bukan angka; dibaca tanpa `TO_CHAR` sehingga pemisah desimal bergantung NLS sesi |
| `WPC` | teks | `DATE` | dibaca tanpa `TO_CHAR` sehingga bentuknya bergantung `NLS_DATE_FORMAT`; ⛔ DDL baru `003` membuat `WPC VARCHAR2(32)` — STRUKTUR menulis "teks", dan itu keliru terhadap tabel warisan |

Sebelas kolom angka sungguhan = sembilan yang ditebak + `STS_REJECT` + `CLAIM_RETRO`.
Sembilan kolom tanggal di dalam 55 = delapan yang ditebak + `WPC`.

## Yang berubah karena dokumen ini

1. `kolomAngkaLama` bertambah `STS_REJECT`, `CLAIM_RETRO`; `kolomTanggalLama` bertambah `WPC`;
   `TipeKolomBarisLama` mengikuti tabel di atas, termasuk panjang `VARCHAR2`.
2. Tabel tiruan skema uji memakai **62** kolom dengan tipe di atas, supaya `INSERT` 55 kolom dan
   `SELECT` 55 kolom diuji terhadap bentuk yang sama dengan produksi.
3. `WPC` di `003_t_claimlf_premiumlist_detail.sql` dan di STRUKTUR: **keputusan work owner** —
   `DATE` mengikuti tabel warisan, atau tetap teks dengan alasan yang ditulis.
4. `CLAIM_RETRO` di `T_GENERAL_CLAIM` (`002`): periksa tipenya terhadap `NUMBER` warisan.
5. Test murni yang membandingkan ketiga peta di `barislamakolom.go` dengan tabel dokumen ini,
   kolom demi kolom, supaya selisih berikutnya terdengar.
