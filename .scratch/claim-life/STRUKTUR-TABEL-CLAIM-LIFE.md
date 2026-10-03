# Struktur Tabel — Claim Life

Acuan bentuk tabel untuk aplikasi Go. Dibuat 2026-09-18 atas perintah work owner.
Presisi fisik (panjang teks, presisi desimal) adalah `[data DBA]` dan TIDAK ditetapkan di sini.
Berkas ini menggambarkan BENTUK, bukan alasan — alasannya ada di `spec.md` dan `issues/`.

⚠️ **Berkas ini adalah acuan TUNGGAL nama kolom. Seluruh nama kolom snake_case**
`[keputusan work owner 2026-09-18]`. Ejaan yang muncul di `spec.md`, `issues/`, dan `revisi-*.md`
adalah **ejaan korpus** — itu **bukti asal kolom**, bukan nama kolom.

**Tipe ditulis sebagai kategori logis:** teks · angka desimal · bilangan bulat · DATE.
Seluruh uang dan share adalah **angka desimal**, tidak pernah float (**ADR-0003**).

**Sumber tiap kolom** adalah salah satu dari dua:

- **korpus** — kolom yang dipakai sistem berjalan, disertai nama rule-nya
- **keputusan** — kolom yang sudah ditetapkan di `spec.md` / `issues/`, disertai nomor tiketnya

Rule korpus yang dirujuk berulang di berkas ini:

| Singkatan | Rule | Tabel |
| --- | --- | --- |
| **UpdOS** | `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` (55 kolom) | `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` |
| **InsOS** | `Claim Life/RDBList/InsertJsonKlaimLife_sql.xml` (51 kolom) | `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` |
| **SaveLPD** | `PremiumList Life/RDBList/SaveMasterLPDet.xml` (80 kolom) | `POOLDATA.M_LIFE_PREMIUM_DETAIL` |

⚠️ Seluruh kolom **nullable** kecuali PK; wajib-isi ditegakkan di Go (`spec.md` §2b, tiket 14).

---

## T_WORK_CLAIM

Tabel work **lintas-lini** (Life + Non-Life). Satu baris mewakili **satu work object** — bisa baris
klaim, bisa baris kasus komite.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | keputusan tiket 14 — teks berformat, **delapan awalan** (lihat di bawah) |
| `COVER_KEY` | teks | ya | FK | keputusan tiket 14 + tiket 00 Komite — menunjuk `T_WORK_CLAIM.ID` induk |
| `LINI` | teks | ya | | keputusan tiket 14 + **keputusan work owner 2026-09-18** — tepat satu dari `FAC` · `PROP` · `NONPROP` · `LIFE` |
| `PY_POSITION` | teks | ya | | keputusan tiket 14 — dari work object Pega `pyPosition` |
| `SENDTO_ADMIN` | teks | ya | | keputusan tiket 14 — dari work object Pega `SendtoAdmin` |
| `SENDTO_MEDICAL` | teks | ya | | keputusan tiket 14 — dari work object Pega `SendtoMedical` |
| `TYPE` | teks | ya | | korpus `TYPE` — UpdOS, InsOS |
| `CASE_ID` | teks | ya | | korpus `CASEID` — UpdOS, InsOS |
| `CREATE_OP` | teks | ya | | keputusan tiket 14 |
| `CREATE_OP_NAME` | teks | ya | | korpus `CREATEOPNAME` — UpdOS, InsOS |
| `TGL_UPDATE` | DATE | ya | | keputusan tiket 14 |

### Delapan awalan `ID` — satu pasang per lini

`[keputusan work owner]` 2026-09-18. Setiap lini punya sepasang awalan sendiri:

| LINI | baris klaim | baris komite |
| --- | --- | --- |
| **FAC** | `CLM-` | `KMT-` |
| **PROP** | `CLMP-` | `TKMT-` |
| **NONPROP** | `CLMNP-` | `KMTNP-` |
| **LIFE** | `CLMLF-` | `KMTLF-` |

`[terverifikasi]` **`TKMT-` bukan salah ketik.** Kelas kerja Claim Prop di Pega adalah
`ASM-FW-GCNMFW-Work-ClaimTreaty` dan kelas komitenya `Work-KomiteTreaty` — huruf **T = Treaty**,
nama Pega untuk lini Prop.

`[keputusan work owner]` **Syariah di luar lingkup** — sistem dan basis datanya terpisah. Varian
`CLMS-` / `CLMPS-` / `CLMNPS-` **tidak** ikut dirancang di sini.

**Aturan baris** `[keputusan work owner]` — ditulis apa adanya:

```
CHECK (
     ( (ID LIKE 'CLM-%' OR ID LIKE 'CLMP-%' OR ID LIKE 'CLMNP-%' OR ID LIKE 'CLMLF-%')
       AND COVER_KEY IS NULL )
  OR ( (ID LIKE 'KMT-%' OR ID LIKE 'TKMT-%' OR ID LIKE 'KMTNP-%' OR ID LIKE 'KMTLF-%')
       AND COVER_KEY IS NOT NULL )
)
```

`[terverifikasi]` **Kedelapan pola tidak saling tumpang tindih** karena **tanda hubung ikut
dihitung**: `CLMP-` tidak cocok dengan `CLM-%` (huruf keempat `P`, bukan `-`), `KMTNP-` tidak
cocok dengan `KMT-%`, dan `TKMT-` tidak berawalan `KMT`. Aturannya **tidak ambigu** — jangan
"memperbaikinya".

⚠️ **Yang TIDAK bisa dijaga `CHECK`.** Aturan `CHECK` hanya melihat **satu baris**. Ia **tidak
sanggup** memastikan baris `TKMT-` menunjuk baris `CLMP-`, atau baris `KMTLF-` menunjuk
`CLMLF-`. **Kecocokan pasangan lini itu aturan Go, bukan aturan basis data.** Cara termurahnya:
**nilai `LINI` baris komite wajib sama dengan `LINI` baris yang ditunjuk `COVER_KEY`.**
⚠️ `[terbuka]` siapa yang memasang aturan itu di Go.

`COVER_KEY` **nullable**, diisi saat baris komite dibuat.

### `LINI` — isinya ditetapkan

`[keputusan work owner]` 2026-09-18. `LINI` berisi **tepat satu** dari **`FAC` · `PROP` ·
`NONPROP` · `LIFE`**, dan **diisi di kedua jenis baris** — baris klaim maupun baris komite.
Baris komite memakai **nilai yang sama** dengan baris klaim induknya, supaya Go cukup membaca
**satu baris** untuk tahu lininya.

⛔ **`LINI` bukan `TYPE`.** `T_WORK_CLAIM` sudah punya kolom **`TYPE`** yang bersumber korpus
(UpdOS, InsOS) dan **artinya lain**. Kedua kolom itu **tidak boleh tertukar**.

### `ACCEPT_STATUS` — DIBUANG dari tabel ini

⚠️ `[keputusan work owner]` 2026-09-18 — kolom **`ACCEPT_STATUS` dihapus dari `T_WORK_CLAIM`**.
Hasil akseptasi adalah **milik kasus komite**, dan **`T_GENERAL_KOMITE.ACCEPT_STATUS` sudah
menyimpannya** — yang itu **tidak disentuh**.

`[terverifikasi]` Korpus mendukung: `.AcceptStatus` **ditimpa tiap putaran** oleh keputusan
penyetuju yang sedang bertugas, dan saat ditolak `KomiteCount` langsung dipaksa sama dengan
`KomiteLoop` sehingga putaran berhenti — **nilai yang tersisa memang keputusan penentu**.

**Index:** `COVER_KEY`.

**Relasi:**

- induk dari dirinya sendiri lewat `COVER_KEY` → `T_WORK_CLAIM.ID` · 1:N · ON DELETE **di Go**
- 1:1 dengan `T_GENERAL_CLAIM` lewat **shared PK** — `ID` sama persis, **tanpa kolom penyambung**
- 1:1 dengan `T_GENERAL_KOMITE` lewat **shared PK** — `ID` sama persis, **tanpa kolom penyambung**
- ditunjuk oleh `T_CLAIMLF_ADJUSTMENT.KOMITE_ID`

---

## T_GENERAL_CLAIM

Header klaim Life. Satu baris mewakili **satu klaim**. `ID`-nya **sama persis** dengan baris klaim di
`T_WORK_CLAIM` (shared PK).

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | keputusan tiket 14 — shared PK = `T_WORK_CLAIM.ID` baris klaim |
| `CLAIM_NO` | teks | ya | | korpus `NO_CLAIM` — UpdOS, InsOS |
| `PY_ID` | teks | ya | | keputusan tiket 14 |
| `STS_KATASTROFE` | teks | ya | | keputusan tiket 14 |
| `KATASTROFE_NOTE` | teks | ya | | keputusan tiket 14 |
| `IS_KPR` | teks | ya | | keputusan tiket 14 |
| `STNC_CLAIM` | teks | ya | | keputusan tiket 14 |
| `ACCEPTED_NO` | teks | ya | | korpus `NO_ACCEPTATION` — UpdOS |
| `STS_REJECT` | teks | ya | | korpus `STS_REJECT` — UpdOS, InsOS |
| `RI_SLIP_RNM` | teks | ya | | keputusan tiket 14 |
| `BUSINESS_NAME` | teks | ya | | korpus `BUSINESSNAME` — UpdOS, InsOS |
| `CLAIM_RETRO` | teks | ya | | korpus `CLAIM_RETRO` — UpdOS, InsOS |
| `CASEID_POLICY` | teks | ya | | keputusan tiket 14 — penunjuk polis |
| `POLICY_NO` | teks | ya | | keputusan tiket 14 (ganti nama dari `PL_NUMBER`); korpus `POLICY_NO` — UpdOS, InsOS |
| `ENDORSMENT_NO` | teks | ya | | keputusan tiket 14 — penunjuk polis |

**Index:** tidak ada di luar PK.

**Relasi:**

- induknya `T_WORK_CLAIM` · **shared PK**, 1:1 · tidak ada kolom penyambung
- anaknya `T_CLAIMLF_PREMIUMLIST_DETAIL` lewat `CLAIM_ID` · 1:N · ON DELETE **CASCADE**
- `CASEID_POLICY` / `POLICY_NO` / `ENDORSMENT_NO` menunjuk **ke luar** (tabel polis, modul
  PremiumList Life) — **bukan** anak, **tidak** ikut cascade

---

## T_CLAIMLF_PREMIUMLIST_DETAIL

Peserta yang diklaim. Satu baris mewakili **satu peserta di dalam satu klaim**.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | keputusan tiket 14 |
| `CLAIM_ID` | teks | ya | FK | keputusan tiket 14 — → `T_GENERAL_CLAIM.ID` |
| `PL_NUMBER` | teks | ya | | korpus `PL_NUMBER` — SaveLPD, UpdOS |
| `POLICY_NO` | teks | ya | | korpus `POLICY_NO` — SaveLPD, UpdOS |
| `POLICY_HOLDER` | teks | ya | | korpus `POLICY_HOLDER` — SaveLPD, UpdOS |
| `CERTIFICATE_NO` | teks | ya | | korpus `CERTIFICATE_NO` — SaveLPD, UpdOS |
| `NAME_OF_INSURED` | teks | ya | | korpus `NAME_OF_INSURED` — SaveLPD, UpdOS |
| `DOB` | DATE | ya | | korpus `DOB` — SaveLPD, UpdOS |
| `AGE` | bilangan bulat | ya | | korpus `AGE` — SaveLPD, UpdOS |
| `SEX` | teks | ya | | korpus `SEX` — SaveLPD, UpdOS |
| `PLAN` | teks | ya | | korpus `PLAN` — SaveLPD, UpdOS |
| `DISEASE` | teks | ya | | korpus `DISEASE` — UpdOS, InsOS |
| `ICD_CODE` | teks | ya | | korpus `ICD_CODE` — UpdOS, InsOS |
| `DESCRIPTION` | teks | ya | | korpus `DESCRIPTION` — SaveLPD |
| `NOTES` | teks | ya | | korpus `NOTES` — UpdOS, InsOS |
| `KETERANGAN` | teks | ya | | korpus `KETERANGAN` — UpdOS, InsOS |
| `DATE_OF_LOSS` | DATE | ya | | keputusan tiket 14 |
| `RECEIVED_DATE` | DATE | ya | | keputusan tiket 14 |
| `BEGIN_DATE` | DATE | ya | | korpus `BEGIN_DATE` — SaveLPD, UpdOS |
| `EFFECTIVE_DATE` | DATE | ya | | korpus `EFFECTIVE_DATE` — SaveLPD |
| `EXPIRED_DATE` | DATE | ya | | korpus `EXPIRED_DATE` — SaveLPD, UpdOS |
| `LAPSE_DATE` | DATE | ya | | korpus `LAPSE_DATE` — SaveLPD, UpdOS |
| `GROSS_VALUATION_BEGIN_DATE` | DATE | ya | | korpus `GROSS_VALUATION_BEGIN_DATE` — SaveLPD |
| `GROSS_VALUATION_EXPIRED_DATE` | DATE | ya | | korpus `GROSS_VALUATION_EXPIRED_DATE` — SaveLPD |
| `RETROCESSION_VALUATION_BEGIN_DATE` | DATE | ya | | keputusan tiket 14; korpus `RETRO_VALUATION_BEGIN_DATE` — SaveLPD |
| `RETROCESSION_VALUATION_EXPIRED_DATE` | DATE | ya | | keputusan tiket 14; korpus `RETRO_VALUATION_EXPIRED_DATE` — SaveLPD |
| `CURRENCY` | teks | ya | | korpus `CURRENCY` — SaveLPD, UpdOS |
| `SUM_INSURED` | angka desimal | ya | | korpus `SUM_INSURED` — SaveLPD, UpdOS |
| `SUM_REASURED` | angka desimal | ya | | korpus `SUM_REASURED` — SaveLPD, UpdOS |
| `GROSS_PREMIUM` | angka desimal | ya | | korpus `GROSS_PREMIUM` — SaveLPD |
| `NET_PREMIUM` | angka desimal | ya | | korpus `NET_PREMIUM` — SaveLPD |
| `CLAIM_AMOUNT` | angka desimal | ya | | korpus `CLAIM_AMOUNT` — SaveLPD, UpdOS |
| `EM_PERCENT` | angka desimal | ya | | korpus `EM_PERCENT` — SaveLPD, UpdOS |
| `SHARE_NUSANTARA_RE` | angka desimal | ya | | korpus `SHARE_NUSANTARA_RE` — SaveLPD, UpdOS |
| `SHARE_RETRO` | angka desimal | ya | | korpus `SHARE_RETRO` — SaveLPD, UpdOS |
| `RETROCEDED_SHARE` | angka desimal | ya | | korpus `RETROCEDED_SHARE` — SaveLPD, UpdOS |
| `CEDING_RETENTION` | angka desimal | ya | | korpus `CEDING_RETENTION` — SaveLPD, UpdOS |
| `WPC` | teks | ya | | korpus `WPC` — SaveLPD, UpdOS |
| `STNC_TREATY` | teks | ya | | keputusan tiket 14; korpus `STNC` — SaveLPD |
| `IS_CHECK` | teks | ya | | keputusan tiket 14 — penanda peserta dipilih untuk diklaim |
| `STATUS` | teks | ya | | korpus `STATUS` — SaveLPD, UpdOS |
| `RECOMMENDATION` | teks | ya | | keputusan tiket 14 |
| `STS_REJECT` | teks | ya | | korpus `STS_REJECT` — UpdOS, InsOS |
| `SOURCE_ID` | teks | ya | | keputusan tiket 14 |
| `CONFIRMATION_DATE` | DATE | ya | | korpus `CONFIRMATION_DATE` — UpdOS, InsOS |
| `COMPLETE_DATE` | DATE | ya | | korpus `COMPLETE_DATE` — UpdOS, InsOS |
| `CLAIM_RECEIVED_DATE` | DATE | ya | | korpus `CLAIM_RECEIVED_DATE` — UpdOS, InsOS |

**Index:** `CLAIM_ID`.

**Relasi:**

- induknya `T_GENERAL_CLAIM` lewat `CLAIM_ID` · 1:N · ON DELETE **CASCADE**
- anaknya `T_CLAIMLF_ADJUSTMENT` lewat `PREMIUM_LIST_DETAIL_ID` · 1:N · ON DELETE **CASCADE**
- anaknya `DOCUMENT_CLAIM` lewat `PREMIUM_LIST_DETAIL_ID` · 1:N · ON DELETE **di Go**

---

## T_CLAIMLF_ADJUSTMENT

Putaran keputusan atas seorang peserta. Satu baris mewakili **satu baris `AdjustmentList`** — inilah
unit keputusan mesin status (**ADR-0011**).

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | keputusan tiket 14 |
| `PREMIUM_LIST_DETAIL_ID` | teks | ya | FK | keputusan tiket 14 — → `T_CLAIMLF_PREMIUMLIST_DETAIL.ID` |
| `SHARE_NUSANTARA_RE` | angka desimal | ya | | keputusan tiket 14 — halaman Pega `AdjustmentList` |
| `CEDING_RETENTION` | angka desimal | ya | | keputusan tiket 14 — halaman Pega `AdjustmentList` |
| `SUM_INSURED` | angka desimal | ya | | keputusan tiket 14 — halaman Pega `AdjustmentList` |
| `SUM_REASURED` | angka desimal | ya | | keputusan tiket 14 — halaman Pega `AdjustmentList` |
| `SHARE_RETRO` | angka desimal | ya | | keputusan tiket 14 — halaman Pega `AdjustmentList` |
| `RETROCEDED_SHARE` | angka desimal | ya | | keputusan tiket 14 — halaman Pega `AdjustmentList` |
| `CURRENCY_ID` | teks | ya | | keputusan tiket 14 — halaman Pega `AdjustmentList` |
| `CURRENCY` | teks | ya | | keputusan tiket 14 — halaman Pega `AdjustmentList` |
| `CLAIM_AMOUNT` | angka desimal | ya | | korpus `CLAIM_AMOUNT` — UpdOS, InsOS |
| `STS_REJECT` | teks | ya | | korpus `STS_REJECT` — UpdOS, InsOS |
| `ACCEPTED_NO` | teks | ya | | korpus `NO_ACCEPTATION` — UpdOS |
| `ACCEPTATION_DATE` | DATE | ya | | korpus `ACCEPTATION_DATE` — UpdOS, InsOS |
| `NAME_OF_BANK` | teks | ya | | korpus `NAME_OF_BANK` — UpdOS |
| `ID_BANK` | teks | ya | | korpus `IDBANK` — UpdOS |
| `ACCOUNT_NO` | teks | ya | | korpus `ACCOUNTNO` — UpdOS |
| `KOMITE_ID` | teks | ya | FK | keputusan tiket 14 — → `T_WORK_CLAIM.ID` baris komite |

**Index:** `PREMIUM_LIST_DETAIL_ID` · `KOMITE_ID` **(UNIK)**.

**Relasi:**

- induknya `T_CLAIMLF_PREMIUMLIST_DETAIL` lewat `PREMIUM_LIST_DETAIL_ID` · 1:N · ON DELETE **CASCADE**
- anaknya `T_CLAIMLF_ADJUSTMENT_SPREADING` lewat `ADJUSTMENT_ID` · 1:N · ON DELETE **CASCADE**
- menunjuk `T_WORK_CLAIM` lewat `KOMITE_ID` · **1:1** · nullable · penunjuk, bukan kaskade
- ditunjuk `T_GENERAL_KOMITE.ADJUSTMENT_ID` · **1:1** — ⚠️ kolom itu punya **dua tabel tujuan** menurut `T_WORK_CLAIM.LINI` (LIFE ke sini, PROP ke `T_CLAIMP_ADJUSTMENT`), sehingga **tidak dipasangi `REFERENCES`**; keutuhannya dijaga Go

⚠️ `[keputusan work owner]` **Satu baris `AdjustmentList` = TEPAT satu kasus komite, dan sebaliknya.**
Karena itu `KOMITE_ID` ber-index **UNIK** meski nullable: baris yang belum pernah dikirim ke Komite
bernilai `NULL`, dan dua baris adjustment **tidak boleh** menunjuk kasus komite yang sama.

---

## T_CLAIMLF_ADJUSTMENT_SPREADING

Hasil spreading sebuah baris adjustment, dipecah per treaty-year. Satu baris mewakili **satu
treaty-year** dari satu baris adjustment. Nilainya **dibekukan** saat adjustment disimpan.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | keputusan tiket 14 |
| `ADJUSTMENT_ID` | teks | ya | FK | keputusan tiket 14 — → `T_CLAIMLF_ADJUSTMENT.ID` |
| `TREATY_TYPE_ID` | teks | ya | | keputusan tiket 14 |
| `TREATY_TYPE_NAME` | teks | ya | | keputusan tiket 14 |
| `TREATY_YEAR_LIFE` | teks | ya | | keputusan tiket 14 |
| `RETROCADED_SHARE` | angka desimal | ya | | keputusan tiket 14 |
| `RATE` | angka desimal | ya | | keputusan tiket 14 |
| `IDR` | angka desimal | ya | | keputusan tiket 14 |
| `USD` | angka desimal | ya | | keputusan tiket 14 |
| `CURRENCY` | teks | ya | | keputusan tiket 14 |

**Index:** `ADJUSTMENT_ID`.

**Relasi:**

- induknya `T_CLAIMLF_ADJUSTMENT` lewat `ADJUSTMENT_ID` · 1:N · ON DELETE **CASCADE**
- anaknya `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` lewat `SPREADING_ID` · 1:N · ON DELETE **CASCADE**

---

## T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO

Pecahan spreading per reinsurer. Satu baris mewakili **satu reinsurer** pada satu baris spreading.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | keputusan tiket 14 |
| `SPREADING_ID` | teks | ya | FK | keputusan tiket 14 — → `T_CLAIMLF_ADJUSTMENT_SPREADING.ID` |
| `REINSURER_NAME` | teks | ya | | keputusan tiket 14 |
| `PERCENT_SHARE` | angka desimal | ya | | keputusan tiket 14 |
| `AMOUNT` | angka desimal | ya | | keputusan tiket 14 |
| `RATE` | angka desimal | ya | | keputusan tiket 14 |
| `PREMIUM_SPREADED_GROSS` | angka desimal | ya | | keputusan tiket 14 |
| `PREMIUM_SPREADED_NET` | angka desimal | ya | | keputusan tiket 14 |
| `COMMISION` | angka desimal | ya | | keputusan tiket 14 — ejaan `COMMISION` (sic), dipertahankan |
| `OVR_COMM` | angka desimal | ya | | keputusan tiket 14 |
| `TREATY_TYPE_ID` | teks | ya | | keputusan tiket 14 |
| `TREATY_TYPE_NAME` | teks | ya | | keputusan tiket 14 |

**Index:** `SPREADING_ID`.

**Relasi:**

- induknya `T_CLAIMLF_ADJUSTMENT_SPREADING` lewat `SPREADING_ID` · 1:N · ON DELETE **CASCADE**
- tidak punya anak

---

## DOCUMENT_CLAIM

Dokumen pendukung klaim, **per peserta**. Tabel **LINTAS-LINI**.

⛔ **Kolomnya tidak ditulis di sini.** `[data DBA]` Kelasnya `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`;
SQL-nya dibuat Pega sendiri dan **nol kemunculan** di rule SQL mana pun, sehingga daftar kolomnya
tidak dapat diturunkan dari korpus maupun dari keputusan yang sudah ada. Menuliskannya berarti
mengarang.

**Index:** `PREMIUM_LIST_DETAIL_ID`.

**Relasi:**

- induknya `T_CLAIMLF_PREMIUMLIST_DETAIL` lewat `PREMIUM_LIST_DETAIL_ID` · 1:N · ON DELETE **di Go**
  (lintas-lini: induknya tabel berbeda per lini, sehingga cascade basis data tidak dapat seragam)

---

## Pohon relasi — enam tingkat, berakar di `T_WORK_CLAIM`

```
TINGKAT 1   T_WORK_CLAIM ─────────────────── akar · LINTAS-LINI · satu baris per work object
            PK  ID   TEKS BERFORMAT, 8 awalan: klaim CLM- CLMP- CLMNP- CLMLF-
                                        komite KMT- TKMT- KMTNP- KMTLF-
            FK  COVER_KEY → T_WORK_CLAIM.ID   (menunjuk dirinya sendiri; NULL bila tak punya induk)
            LINI (FAC/PROP/NONPROP/LIFE) · PY_POSITION · SENDTO_ADMIN
            SENDTO_MEDICAL · TYPE   <- TYPE bukan LINI, artinya lain
            CASE_ID · CREATE_OP · CREATE_OP_NAME · TGL_UPDATE
            |
   +--------+-------------------------------------------+
   | baris KLAIM  COVER_KEY = NULL                      | baris KOMITE  COVER_KEY = ID baris klaim
   |                                                    |
TINGKAT 2                                           TINGKAT 2
   +--1:1-- T_GENERAL_CLAIM                             +--1:1-- T_GENERAL_KOMITE
           PK ID = T_WORK_CLAIM.ID   <- SHARED PK               PK ID = T_WORK_CLAIM.ID baris komite
             (tidak ada kolom FK terpisah)                        <- SHARED PK, tidak ada kolom FK
           penunjuk polis (ke LUAR, bukan anak):                FK ADJUSTMENT_ID -> tabel adj. menurut LINI
             CASEID_POLICY                                        (tutup lingkar, 1:1, NOT NULL)
             POLICY_NO                                          KOMITE_LOOP · KOMITE_COUNT · ACCEPT_STATUS
             ENDORSMENT_NO                                      |
           |                                         TINGKAT 3  |
TINGKAT 3  |                                                    +--1:N-- T_KOMITE_KOMITELIST
           +--1:N-- T_CLAIMLF_PREMIUMLIST_DETAIL                        PK ID
           |        PK ID                                              FK DATA_KOMITE_ID
           |        FK CLAIM_ID -> T_GENERAL_CLAIM.ID  CASCADE             -> T_GENERAL_KOMITE.ID  CASCADE
           |        |                                                   KOMITE_URUT · KOMITE_OPERATORID
TINGKAT 4  |        +--1:N-- T_CLAIMLF_ADJUSTMENT                       KOMITE_JABATAN
           |        |        PK ID                                      KOMITE_EMAIL · KOMITE_APPROVAL
           |        |        FK PREMIUM_LIST_DETAIL_ID                  KOMITE_COMMENT · DATE_APPROVE
           |        |           -> T_CLAIMLF_PREMIUMLIST_DETAIL.ID  CASCADE
           |        |        FK KOMITE_ID -> T_WORK_CLAIM.ID  <- penunjuk balik KE ATAS,
           |        |                                            nullable, 1:1, index UNIK
           |        |        |
TINGKAT 5  |        |        +--1:N-- T_CLAIMLF_ADJUSTMENT_SPREADING
           |        |                 FK ADJUSTMENT_ID -> T_CLAIMLF_ADJUSTMENT.ID  CASCADE
           |        |                 |
TINGKAT 6  |        |                 +--1:N-- T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO
           |        |                          FK SPREADING_ID
           |        |                             -> T_CLAIMLF_ADJUSTMENT_SPREADING.ID  CASCADE
TINGKAT 4  |        +--1:N-- DOCUMENT_CLAIM                      <- LINTAS-LINI
           |                 FK PREMIUM_LIST_DETAIL_ID -> ...DETAIL.ID   (Life saja)
           |                 ON DELETE di Go -- induk beda tabel per lini
           |
           +-- DIBACA dari luar, BUKAN anak, TIDAK ikut cascade:
                 T_PREMIUM_LIST dkk (polis) · master marketing officer · EMAILKOMITE ·
                 master retro · security reinsurer · currency
```

## Sebelas relasi

| # | Induk | Anak | Kunci tamu | Kard | Hapus |
| --- | --- | --- | --- | --- | --- |
| 1 | `T_WORK_CLAIM` | `T_WORK_CLAIM` | `COVER_KEY` | 1:N | di Go |
| 2 | `T_WORK_CLAIM` | `T_GENERAL_CLAIM` | **tidak ada kolom terpisah** — `T_GENERAL_CLAIM.ID` = `T_WORK_CLAIM.ID` (**shared PK**) | 1:1 | — |
| 3 | `T_GENERAL_CLAIM` | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `CLAIM_ID` | 1:N | CASCADE |
| 4 | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `T_CLAIMLF_ADJUSTMENT` | `PREMIUM_LIST_DETAIL_ID` | 1:N | CASCADE |
| 5 | `T_CLAIMLF_ADJUSTMENT` | `T_CLAIMLF_ADJUSTMENT_SPREADING` | `ADJUSTMENT_ID` | 1:N | CASCADE |
| 6 | `T_CLAIMLF_ADJUSTMENT_SPREADING` | `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` | `SPREADING_ID` | 1:N | CASCADE |
| 7 | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `DOCUMENT_CLAIM` | `PREMIUM_LIST_DETAIL_ID` | 1:N | **di Go** |
| 8 | `T_WORK_CLAIM` | `T_GENERAL_KOMITE` | **tidak ada kolom terpisah** — `T_GENERAL_KOMITE.ID` = `T_WORK_CLAIM.ID` baris komite (**shared PK**) | 1:1 | di Go |
| 9 | `T_GENERAL_KOMITE` | `T_KOMITE_KOMITELIST` | `DATA_KOMITE_ID` | 1:N | CASCADE |
| 10 | **tabel adjustment menurut `LINI`** — LIFE `T_CLAIMLF_ADJUSTMENT` · PROP `T_CLAIMP_ADJUSTMENT` | `T_GENERAL_KOMITE` | `ADJUSTMENT_ID` **NOT NULL, index UNIK, tanpa `REFERENCES`** | **1:1** | di Go |
| 11 | `T_CLAIMLF_ADJUSTMENT` | `T_WORK_CLAIM` | `KOMITE_ID` **nullable, index UNIK** | **1:1** | penunjuk |

Relasi **2** dan **8** tidak punya kunci tamu untuk di-index — keduanya **shared primary key**, dan
PK sudah ber-index dengan sendirinya. Seluruh kunci tamu lain **ber-index**.

---

## Catatan — belum ditetapkan, TIDAK menghambat berkas ini

- `[terbuka]` `COMMISION` kurang satu huruf S, tetapi `[terverifikasi]` **itu ejaan korpus** — 369 berkas XML, 2.691 kemunculan; `COMMISSION` (dua S) **juga** ada di korpus, 83 berkas. Korpus memakai kedua ejaan. Perbaikan belum diputuskan
- `[data DBA]` daftar kolom `DOCUMENT_CLAIM`
- `[data DBA]` presisi fisik seluruh kolom
- `[terbuka]` `REFERENCES T_WORK_CLAIM(ID)` pada `KOMITE_ID` dan `COVER_KEY`
- `[terbuka]` kolom nomor klaim (`NO_CLAIM`) dan nomor akseptasi (`NO_ACCEPTATION`) belum punya rumah
- `[terbuka]` isi `CASEID_POLICY`, `POLICY_NO`, `ENDORSMENT_NO`
- `[terbuka]` rumah lima kolom yatim bekas `T_CLAIM_POLICY`
- `[terbuka]` generator nomor untuk **delapan awalan** (`CLM-` `CLMP-` `CLMNP-` `CLMLF-` / `KMT-` `TKMT-` `KMTNP-` `KMTLF-`) pada `T_WORK_CLAIM.ID`
- `[terbuka]` `PREMIUM_SPREADED_NET` punya dua rumus di rule yang sama
- `[terbuka]` `ACCEPTED_NO` (`T_CLAIMLF_ADJUSTMENT`) dan `ACCEPTED_NO` (`T_GENERAL_CLAIM`) dieja berbeda untuk nilai yang bersumber sama
- `[terbuka]` `STS_REJECT` ada di tiga tabel sekaligus (header, peserta, adjustment)
- `[terbuka]` `CLAIM_AMOUNT`, `SHARE_RETRO`, `RETROCEDED_SHARE`, `SUM_INSURED`, `SUM_REASURED`, `SHARE_NUSANTARA_RE`, `CEDING_RETENTION`, `CURRENCY` ada di peserta dan di adjustment sekaligus
- `[terbuka]` `T_GENERAL_CLAIM.POLICY_NO` dan `T_CLAIMLF_PREMIUMLIST_DETAIL.POLICY_NO` bernama sama
- `[terbuka]` `RETROCADED_SHARE` (spreading) dan `RETROCEDED_SHARE` (peserta, adjustment) dieja berbeda
- `[terbuka]` siapa yang memasang aturan **`LINI` baris komite == `LINI` baris `COVER_KEY`** di Go — `CHECK` tidak sanggup menjaganya
- `[terbuka]` bila Claim Non Prop kelak punya tabel adjustment sendiri, tujuan `T_GENERAL_KOMITE.ADJUSTMENT_ID` menjadi **tiga**
- `[terbuka]` tipe `AGE` sebagai bilangan bulat belum dikonfirmasi DBA
