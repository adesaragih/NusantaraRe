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
| `POSITION` | teks | ya | | keputusan tiket 14 — dari work object Pega `pyPosition` (NAMA PERAN); bernama `PY_POSITION` sampai migrasi `023` (**keputusan work owner 01-10-2026** — nama sama dengan `T_WORK_POLIS.POSITION`, artinya berbeda: polis menyimpan `Position` layar) |
| `SENDTO_ADMIN` | teks | ya | | keputusan tiket 14 — dari work object Pega `SendtoAdmin` |
| `SENDTO_MEDICAL` | teks | ya | | keputusan tiket 14 — dari work object Pega `SendtoMedical` |
| `CREATE_OP` | teks | ya | | keputusan tiket 14 |
| `CREATE_OP_NAME` | teks | ya | | korpus `CREATEOPNAME` — UpdOS, InsOS |
| `TGL_UPDATE` | DATE | ya | | keputusan tiket 14 — waktu UBAH, ditimpa tiap perpindahan |
| `TAHAP` | teks | ya | | **butir at** `[DIPUTUSKAN 27-09-2026]` — nama assignment VERBATIM `pyTaskName` (`Register_Flow.xml` 358 · 343 · 268 · 313). Ada karena `POSITION` tidak dapat membedakan **Input Register** dari **Outstanding Claim** (keduanya `ReasLifeAdmin`), sedangkan `Send Back to Register` (`InputOSClaimLife.xml:21404`) membuktikan keadaan itu dapat dituju kembali |
| `TGL_CREATE` | DATE | ya | | **butir au** — padanan `pxCreateDateTime`; kotak masuk diurutkan dengannya (`InboxPremiumList.xml:736`). `TGL_UPDATE` tidak dapat dipakai: ia ditimpa tiap perpindahan |
| `STATUS_WORK` | teks | ya | | **butir bb** `[DIPUTUSKAN 27-09-2026]` — status kerja kasus. **Satu-satunya nilai yang ditulis: `Resolved-Completed`**, VERBATIM `Register_Flow.xml` baris 899 *(shape `End1`, `rowdata REPEATINGINDEX="End1"` baris 883, `Data-MO-Event-End` baris 901; sembilan shape lain ber-`pyWorkStatus` kosong)*. **NULL = kasus belum ditutup** — status bawaan Pega untuk kasus berjalan tidak ada di ekspor dan tidak dikarang. Saat tutup, `TAHAP` **dikosongkan** *(`FinishAssignment` tanpa parameter, `ProtectCloseClaim_act` baris 838)* |

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

⛔ **`LINI` bukan `TYPE`.** Header klaim `T_GENERAL_CLAIM` punya kolom **`TYPE`** yang bersumber
korpus (UpdOS, InsOS) dan **artinya lain**. Kedua kolom itu **tidak boleh tertukar**.

### `TYPE` dan `CASE_ID` — KELUAR dari tabel ini (migrasi `023`)

⚠️ `[keputusan work owner]` 2026-10-01 — **`TYPE` pindah ke `T_GENERAL_CLAIM`** (dibaca
`TypeKlaim`, kontrak `inti/backend/kontrak/klaim.go`), dan **`CASE_ID` dibuang** — `ID` dipakai
sebagai gantinya: kasus baru memang `CASE_ID = ID`, dan migrasi klaim lama memakai `CASEID`
warisan sebagai `ID`. Migrasi `023` menolak berjalan (ORA-02293) bila satu baris saja punya
`CASE_ID` yang berbeda dari `ID`. Baris Komite (`KMTLF-`) tidak lagi menyimpan `TYPE`:
salinannya dulu tidak pernah dibaca — `TypeKlaim` membaca klaim induknya.

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
| `CLAIM_RETRO` | desimal | ya | | korpus `CLAIM_RETRO` — UpdOS, InsOS; ✅ **ralat 26-09-2026**, lihat bawah |
| `CURRENCY` | teks | ya | | ⭐ **BARU 26-09-2026** — `[keputusan work owner butir z1]`, lihat bawah |
| `CASEID_POLICY` | teks | ya | | keputusan tiket 14 — penunjuk polis |
| `POLICY_NO` | teks | ya | | keputusan tiket 14 (ganti nama dari `PL_NUMBER`); korpus `POLICY_NO` — UpdOS, InsOS |
| `ENDORSMENT_NO` | teks | ya | | keputusan tiket 14 — penunjuk polis |
| `BUSINESS_CODE` | teks | ya | | **temuan audit A0**, A1 — nomor akseptasi memuatnya (`Generate_NoAccept_Life` 85), model relasional tidak menyimpannya |
| `TYPE` | teks | ya | | korpus `TYPE` — UpdOS, InsOS; pindah dari `T_WORK_CLAIM` (migrasi `023`, **keputusan work owner 01-10-2026**). **Bukan** `LINI` |
| `SUMBER` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `MASTER_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `TREATY_NAME` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `TREATY_GROUP_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `TREATY_GROUP_NAME` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `PROPORTION_TYPE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `START_DATE_TREATY` | DATE | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `END_DATE_TREATY` | DATE | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `YEAR_OF_ACCOUNT` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `RNM_SHARE_PCT` | angka desimal | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `CLAIM_NO_TEMP` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `POLICY_START_DATE` | DATE | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `POLICY_END_DATE` | DATE | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `QUARTER` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `YEAR_OF_QUARTAL` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `TREATY_YEAR` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `POLICY_NO_CEDING` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `INSURED_NAME` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `PLA_NO_CEDING` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `PLA_NO_SOB` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `PERIOD_POLICY_TBA` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DATE_OF_LOSS` | DATE | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `REPORT_DATE` | DATE | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DATE_RECEIVED` | DATE | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `REPORTER_NAME` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `REPORTER_PHONE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `REPORT_TYPE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `REPORTER_STATUS` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `INSURED_RELATIONSHIP_OTHERS` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `REPORT_ADDRESS` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `CATASTROPHE_STATUS` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `NON_CATASTROPHE_TYPE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `CATASTROPHE_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `CAUSE_OF_LOSS` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `CAUSE_OF_LOSS_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `CONSULTANT_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `CONSULTANT_NAME` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `ADJUSTER_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `ADJUSTER_NAME` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `REPORT_DESCRIPTION` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `LOCATION_OF_LOSS` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `OCCUPATION` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `PROVINCE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `PROVINCE_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `POSTAL_CODE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `RW` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `RW_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DISTRICT` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DISTRICT_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `CITY` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `CITY_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `INSURED_INTEREST` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `SHARE_CEDING` | angka desimal | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DEDUCTIBLE_TYPE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DEDUCTIBLE_FORM_TYPE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DEDUCTIBLE_CURRENCY_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DEDUCTIBLE_PCT` | angka desimal | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DEDUCTIBLE_BASIS` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `TSI_DEDUCTIBLE` | angka desimal | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DEDUCTIBLE_VALUE` | angka desimal | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `NET_DEDUCTIBLE_VALUE` | angka desimal | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `IS_OUTSTANDING` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `IS_ACCEPTATION` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `IS_CFS` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `RE_CFS` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `IS_REALISATION` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `AKTIF_BUTTON` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `IS_PLA` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `NO_PLA` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `REMARK` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `REMARK_CLOSE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `PAYABLE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `PAYABLE_TO` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DLA_NO_CEDING` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `DLA_NO_SOB` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `MARKETING_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `MARKETING_CLIENT_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `MARKETING_CLIENT_NAME` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `MARKETING_TEAM_GROUP` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `MARKETING_BRANCH_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `MARKETING_BRANCH_NAME` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `RECEIVER_NAME` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `RECEIVER_BANK_NAME` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `RECEIVER_BANK_BRANCH` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `RECEIVER_ACCOUNT_NO` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `RECEIVER_SWIFT_CODE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `RECEIVER_BANK_ID` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `START_DATE_ESTIMATION` | DATE | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `END_DATE_ESTIMATION` | DATE | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `START_DATE_ADJUSTMENT` | DATE | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `IS_CLOSE_FILE` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `IS_RESERVED_CLAIM` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `IS_ANY_ACCEPTATION` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |
| `IS_SUBJECTIVITY` | teks | ya | | Claim Prop — migrasi `520` (lini `PROP`, keputusan work owner 07-10-2026 "Tabel bersama"); nullable, wajib-isi di Go |

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
| `WPC` | DATE | ya | | korpus `WPC` — SaveLPD, UpdOS; ⚠️ **ralat 26-09-2026**, lihat bawah |
| `STNC_TREATY` | teks | ya | | keputusan tiket 14; korpus `STNC` — SaveLPD |
| `IS_CHECK` | teks | ya | | keputusan tiket 14 — penanda peserta dipilih untuk diklaim |
| `STATUS` | teks | ya | | korpus `STATUS` — SaveLPD, UpdOS |
| `RECOMMENDATION` | teks | ya | | keputusan tiket 14 |
| `STS_REJECT` | teks | ya | | korpus `STS_REJECT` — UpdOS, InsOS |
| `SOURCE_ID` | teks | ya | | keputusan tiket 14 |
| `CONFIRMATION_DATE` | DATE | ya | | korpus `CONFIRMATION_DATE` — UpdOS, InsOS |
| `COMPLETE_DATE` | DATE | ya | | korpus `COMPLETE_DATE` — UpdOS, InsOS |
| `CLAIM_RECEIVED_DATE` | DATE | ya | | korpus `CLAIM_RECEIVED_DATE` — UpdOS, InsOS |
| `STS_HAPUS` | teks | ya | | migrasi 022 (OQ-M6, GILIRAN-17): penanda cabut peserta — NULL aktif, `'1'` dicabut; tombol `DELETE` `InputOSClaimLife` b17865; setiap pembaca menyaringnya (ADR-U-0031) |

**Index:** `CLAIM_ID`.

**Relasi:**

- induknya `T_GENERAL_CLAIM` lewat `CLAIM_ID` · 1:N · ON DELETE **CASCADE**
- anaknya `T_CLAIMLF_ADJUSTMENT` lewat `PREMIUM_LIST_DETAIL_ID` · 1:N · ON DELETE **CASCADE**
- anaknya `T_CLAIMLF_DOCUMENT` lewat `PREMIUM_LIST_DETAIL_ID` · 1:N · ON DELETE **di Go**

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
| `BRANCH_OF_BANK` | teks | ya | | **butir aj**, A1 — `AdjustmentDetail_Section` 6870 `BranchOfBank` |
| `SWIFT_CODE` | teks | ya | | **butir aj**, A1 — `AdjustmentDetail_Section` 6665 `SwiftCode` |
| `PAYABLE_TO` | teks | ya | | **butir aj**, A1 — `AdjustmentDetail_Section` 5929 `PayableTo` |

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

> **⚠️ RALAT 26 September 2026 — nama fisiknya `T_CLAIMLF_ADJ_SPREADING_RETRO`.**
>
> Nama di judul bab ini **36 byte**. `[keputusan work owner butir j, 26-09-2026]`: dipendekkan
> menjadi **`T_CLAIMLF_ADJ_SPREADING_RETRO` (29 byte)**, dan itulah nama yang dipakai DDL, kode, dan
> test sejak commit `8f5453b`.
>
> ⚠️ **Alasan aslinya — batas 30 byte — kini `[dugaan]`, bukan `[terverifikasi]`.** Oracle
> menolak pengenal di atas 30 byte dengan `ORA-00972` hanya bila `COMPATIBLE < 12.2`. Katalog
> instance **pengembangan** dibaca 26-09-2026 siang: versinya **12.2.0.1**, dan `COMPATIBLE`
> **tidak terbaca** dengan hak yang ada. Versi **produksi** belum diketahui sama sekali. Jadi belum
> dapat dipastikan apakah nama 36 byte itu sungguh ditolak. Keputusan j **tidak dibatalkan sendiri
> oleh executor**; bila DBA memastikan `COMPATIBLE >= 12.2` di kedua lingkungan, j dapat ditinjau
> ulang oleh work owner. Berkas migrasinya ikut diganti nama
> menjadi `006_t_claimlf_adj_spreading_retro.sql` pada 26-09-2026 `[butir n]`, aman karena migrasi
> belum pernah dijalankan di Oracle mana pun.
>
> `[terverifikasi]` Pemendekan ini sejalan dengan sistem berjalan, bukan menyimpang darinya:
> seluruh nama tabel fisik yang sungguh dipakai Pega di modul Claim Life berukuran **≤ 23 byte**
> (terpanjang `OS_AKSEPTASI_KLAIM_LIFE`). Nama 36 byte tidak pernah ada di korpus — ia lahir di
> dokumen ini, bukan di Pega.
>
> Judul bab, tabel relasi, dan diagram di bawah **sengaja tidak ditulis ulang**: yang berubah hanya
> nama fisiknya, dan menulis ulang dokumen akan menghapus jejak kenapa ia berubah. Di mana pun
> dokumen ini menulis `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`, yang ada di Oracle adalah
> `T_CLAIMLF_ADJ_SPREADING_RETRO`.

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

## T_CLAIMLF_DIAGNOSE

Diagnosa peserta. Satu baris mewakili **satu diagnosa pada satu peserta** — daftarnya banyak,
dan itu **rancangan**: grid `.DiagnoseList` punya tombol `Add` (b4690 → `addRow`) dan `Delete`
(b6160 → `deleteRow`).

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bilangan bulat | tidak | PK | keputusan butir bd — `SEQ_CLAIMLF_DIAGNOSE` (**ADR-0006**) |
| `PREMIUM_LIST_DETAIL_ID` | teks | tidak | FK | keputusan butir bd — → `T_CLAIMLF_PREMIUMLIST_DETAIL.ID`, `ON DELETE CASCADE` |
| `URUTAN` | bilangan bulat | tidak | | korpus `.pxListSubscript` — urutan grid; `SetSTS_Reject.xml` b241 memutarnya |
| `ICD_CODE` | teks | ya | | korpus `.ICDCODE` b5616 (read-only) — diisi `SetDisease.xml` b307 |
| `DISEASE` | teks | ya | | korpus `.DISEASE` b5422 (read-only) — diisi `SetDisease.xml` b260 |
| `GROUP_DIAGNOSE` | teks | ya | | korpus `.GROUPDIAGNOSE` b5860 — `pxDropdown` b5863, `pyListSource associated`; **daftar pilihannya `[tidak ada di korpus]`** |
| `STS_REJECT` | teks | ya | | korpus — `SetSTS_Reject.xml` b257 `.STS_REJECT = Primary.STS_REJECT` |

**Index:** `PREMIUM_LIST_DETAIL_ID`.

**Relasi:**

- induknya `T_CLAIMLF_PREMIUMLIST_DETAIL` lewat `PREMIUM_LIST_DETAIL_ID` · 1:N · ON DELETE **CASCADE**

⚠️ **Gerbang sunting**, dibaca dari section dan dicatat: `pyDisabledWhen`
`.STS_REJECT=='1' || .STS_REJECT=='2'` muncul **empat kali** pada kendali di dalam grid
(b4600–b6200). Baris yang sudah diaksep atau ditolak **tidak dapat disunting lagi**.

⛔ `DISEASE` berlebar **1000**, bukan 255. Nilainya datang dari `POOLDATA.DISEASE_LIFE.DISEASE`
`VARCHAR2(1000)` yang isi terpanjangnya **290** `[data DBA — katalog DEV]`; 255 **terbukti kurang**.
Migrasi `018` melebarkan `T_CLAIMLF_PREMIUMLIST_DETAIL.DISEASE` pula.

⚠️ **RALAT 27-09-2026 atas kalimat di atas:** kalimat gerbang itu muncul **tujuh** kali di
`ClaimLifeDetailGCNM.xml`, bukan empat. Empat di dalam grid *(b4682, b5059, b5870, b6152)*; tiga
lagi menjaga medan catatan `.ADMIN_NOTES` b2628, `.RECOMMENDATION` b7335, dan `.NOTES` b15234.
Artinya peserta yang sudah diputus membekukan **seluruh** isian layar Detail.

## T_CLAIMLF_STORAGE

Kartu berkas unggahan. Satu baris mewakili **satu berkas di penyimpanan**, dan
`T_CLAIMLF_DOCUMENT.T_STORAGE_ID` menunjuk `IMAGEID`-nya.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `IMAGEID` | teks | tidak | PK | korpus `Insert_T_Storage_SQL.xml` b86 — nilainya `models.IDDokumenBaru` |
| `URLPUBLIC` | teks | ya | | korpus b87 — dibaca `GetLinkStorage_SQL.xml` b85 sebagai `UploadDoc.URLImage` |
| `APPFOLDER` | teks | ya | | korpus b88 |
| `EXPDATE` | DATE | ya | | korpus b89 — `To_date(..., 'DD/MM/YYYY HH24:MI:SS')` b97 |
| `FILENAME` | teks | ya | | korpus b90 |
| `APPNAME` | teks | ya | | korpus b91 — padanan `GCP_IMAGE.APPNAME VARCHAR2(20)` `[data DBA]` |
| `STORAGE` | teks | ya | | korpus b92 — `Insert_T_Storage_SQL` b100 menulis literal `'standard'` |
| `TANGGAL_UPLOAD` | DATE | ya | | korpus `Update_T_Storage_SQL.xml` b89 — `To_date({UpdateDoc.DateTime}, 'MM/DD/YYYY HH24:MI:SS')`; sumbernya `UploadDoc.Response.DateTime` *(`GetUrlGoogleStorage_Act.xml` b2274)* |

**Index:** hanya PK.

**Relasi:** nol kunci tamu — lihat di bawah.

⛔ **TABEL KAMI SENDIRI, bukan `T_STORAGE_IMAGE`.** Kolomnya ditiru nama demi nama supaya
migrasi data kelak mekanis; yang tidak ditiru adalah **tempatnya**. `T_STORAGE_IMAGE` tabel
**bersama** lintas modul dan lintas aplikasi — kolom `APPNAME` ada justru karena itu — dan brief
menuntut persetujuan manusia sebelum penyambungan penyimpanan nyata. Sebab kedua: rule
penulisnya, `Insert_T_Storage_SQL.xml` b102, ber-`commit;`, yang **ADR-U-0029** larang.

⚠️ **Nol FK ke `T_CLAIMLF_DOCUMENT`**, dan itu meniru aslinya: `T_STORAGE_ID` menunjuk
`IMAGEID` tanpa constraint, sebab di sistem lama barisnya lahir di layanan luar dan boleh
mendahului maupun menyusul baris dokumennya. FK di sini akan menolak urutan yang sah.

⛔ **RALAT 28-09-2026 — `TANGGAL_UPLOAD` terlewat di migrasi 019.** Bab ini semula disusun dari
`Insert_T_Storage_SQL.xml` saja, yang memang tidak menyebut kolom itu: ia ditulis oleh **UPDATE**,
bukan oleh INSERT *(`Update_T_Storage_SQL.xml` b89)*. Membaca satu dari dua rule penulis lalu
menyimpulkan tentang tabelnya — bentuk kekeliruan yang sama untuk kelima kalinya. Migrasi **020**
menambahkannya.

⚠️ **Dua bentuk tanggal berbeda di satu pernyataan, dan itu ada di rule aslinya:**
`EXPDATE` memakai `'DD/MM/YYYY HH24:MI:SS'` sedangkan `TANGGAL_UPLOAD` memakai
`'MM/DD/YYYY HH24:MI:SS'`. Nilai warisan pada kedua kolom itu karena itu **tidak dapat dibedakan**
untuk tanggal 1–12 tiap bulan. Kami tidak mewarisi masalahnya — Go mengikat `time.Time`, bukan
teks — tetapi migrasi data **A4** harus tahu.

⛔ **`IMAGEID` punya rule pembangkitnya sendiri**, dan ia BUKAN pengenal dokumen:

```
RDBList/GenerateImageID_SQL.xml  b85-b88
  SELECT STANDARD_HASH(
           'ASMPP' || TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF9') || SYS_GUID(),
           'MD5') AS "InsertDoc.ImageID"
  FROM DUAL
```

dipanggil `InsertGoogleStorage_Act.xml` b2226, mengalir ke `Param.ImageID` b2784–2785, lalu ke
`Insert_T_Storage_SQL.xml` b94. Ditiru `models.ImageIDBaru` — MD5, **heksa huruf besar**, 32
karakter. Ronde pertama memakai pengenal dokumen *(cap waktu)*; kunci penyimpanan yang dapat
ditebak dari waktu unggah bukan kunci.

## T_CLAIMLF_DOCUMENT

> ⭐ **RALAT 26 September 2026 — bab ini dulu bernama `DOCUMENT_CLAIM`.**
>
> `[keputusan work owner butir v1]`. Katalog instance pengembangan dibaca 26-09-2026 dan
> **`POOLDATA.DOCUMENT_CLAIM` ternyata SUDAH ADA**: 14 kolom milik kelas Pega
> `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM` (`IDPEGA`, `NAMAFILE`, `MIME`, `KATEGORI_1/2`, `NOAKSEP`,
> `NOPREKAS`, `PAYMENTDATE`, `INSKEY_LINK`, `INSKEY_DATA`, `T_STORAGE_ID`, `PXCREATEOPERATOR`, …),
> berisi **295 baris**. Kolom itu bukan kolom yang didaftar bab ini.
>
> Membuat tabel bernama sama berarti salah satu dari dua hal, dan keduanya buruk: migrasi gagal
> `ORA-00955`, atau — lebih buruk — migrasi **melewatinya sebagai "sudah ada"** dan aplikasi
> berjalan di atas tabel yang kolomnya bukan miliknya. Karena itu tabel baru memakai nama
> sendiri, **`T_CLAIMLF_DOCUMENT`**, dan `POOLDATA.DOCUMENT_CLAIM` warisan **tidak disentuh sama
> sekali**.
>
> Berkas migrasinya ikut: `007_t_claimlf_document.sql` (+`_down`), beserta
> `SEQ_T_CLAIMLF_DOCUMENT`. ⚠️ Memakai tabel warisan apa adanya dengan pemetaan kolom
> **(v2, ADR-U-0042)** tetap terbuka untuk ditinjau saat tiket dokumen dikerjakan — bentuknya kini
> dapat direkayasa balik, dan itu keputusan tersendiri.

Dokumen pendukung klaim, **per peserta**. Tabel **LINTAS-LINI**.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | bulat | tidak | PK | keputusan tiket 14 — dari sequence |
| `PREMIUM_LIST_DETAIL_ID` | teks | ya | FK | keputusan tiket 14 — → `T_CLAIMLF_PREMIUMLIST_DETAIL.ID` |
| `NAMA_FILE` | teks | ya | | katalog `DOCUMENT_CLAIM.NAMAFILE`; ditulis `InsertDocument_Act` |
| `MIME` | teks | ya | | katalog `DOCUMENT_CLAIM.MIME`; `InsertDocument_Act` menulisnya HURUF KECIL (`@toLowerCase`) |
| `KATEGORI_1` | teks | ya | | katalog; `Param.KATEGORI_1` di `InsertDocument_Act` |
| `KATEGORI_2` | teks | ya | | katalog; satu-satunya kategori yang diisi manusia di `Section/DocumentLife.xml` |
| `TANGGAL` | DATE | ya | | katalog `DOCUMENT_CLAIM.TANGGAL` |
| `T_STORAGE_ID` | teks | ya | | penunjuk berkas di Google Storage (ADR-U-0010) |
| `PAYMENT_DATE` | DATE | ya | | katalog `DOCUMENT_CLAIM.PAYMENTDATE` |

⭐ **Ketujuh kolom isi ditambahkan 26-09-2026** lewat langkah migrasi `010`, `[keputusan work owner
butir ad]`. Sumbernya dan tujuh kolom warisan yang sengaja TIDAK ikut dijelaskan di blok akhir
dokumen ini, yang **dihitung dua kali dari sumber berbeda**.

~~⛔ **Kolomnya tidak ditulis di sini.**~~ `[data DBA]` Kelasnya `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`;
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
            LINI (FAC/PROP/NONPROP/LIFE) · POSITION (pyPosition) · SENDTO_ADMIN
            SENDTO_MEDICAL · CREATE_OP · CREATE_OP_NAME · TGL_UPDATE
            TAHAP · TGL_CREATE · STATUS_WORK
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
TINGKAT 4  |        +--1:N-- T_CLAIMLF_DOCUMENT                  <- LINTAS-LINI
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
| 7 | `T_CLAIMLF_PREMIUMLIST_DETAIL` | `T_CLAIMLF_DOCUMENT` | `PREMIUM_LIST_DETAIL_ID` | 1:N | **di Go** |
| 8 | `T_WORK_CLAIM` | `T_GENERAL_KOMITE` | **tidak ada kolom terpisah** — `T_GENERAL_KOMITE.ID` = `T_WORK_CLAIM.ID` baris komite (**shared PK**) | 1:1 | di Go |
| 9 | `T_GENERAL_KOMITE` | `T_KOMITE_KOMITELIST` | `DATA_KOMITE_ID` | 1:N | CASCADE |
| 10 | **tabel adjustment menurut `LINI`** — LIFE `T_CLAIMLF_ADJUSTMENT` · PROP `T_CLAIMP_ADJUSTMENT` | `T_GENERAL_KOMITE` | `ADJUSTMENT_ID` **NOT NULL, index UNIK, tanpa `REFERENCES`** | **1:1** | di Go |
| 11 | `T_CLAIMLF_ADJUSTMENT` | `T_WORK_CLAIM` | `KOMITE_ID` **nullable, index UNIK** | **1:1** | penunjuk |

Relasi **2** dan **8** tidak punya kunci tamu untuk di-index — keduanya **shared primary key**, dan
PK sudah ber-index dengan sendirinya. Seluruh kunci tamu lain **ber-index**.

---

## Catatan — belum ditetapkan, TIDAK menghambat berkas ini

- `[terbuka]` `COMMISION` kurang satu huruf S, tetapi `[terverifikasi]` **itu ejaan korpus** — 369 berkas XML, 2.691 kemunculan; `COMMISSION` (dua S) **juga** ada di korpus, 83 berkas. Korpus memakai kedua ejaan. Perbaikan belum diputuskan
- `[data DBA]` daftar kolom `T_CLAIMLF_DOCUMENT`
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


---

## ⚠️ RALAT 26 September 2026 — tiga tipe kolom yang ditebak dari namanya

Sampai hari ini tipe kolom `OS_AKSEPTASI_KLAIM_LIFE` tidak diketahui: korpus Pega tidak memuat DDL,
dan dokumen ini menurunkan tipenya dari **nama** kolom. Katalog `ALL_TAB_COLUMNS` instance
pengembangan dibaca 26-09-2026 dan disimpan sebagai `[data DBA]` di
`.scratch\claim-life\TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md`. Tiga tebakan meleset:

| Kolom | Ditulis dokumen ini | Katalog | Keadaan |
| --- | --- | --- | --- |
| `WPC` | teks | **`DATE`** | ✅ `[keputusan work owner butir w]` — baris di atas diralat menjadi **tanggal**, dan `003_t_claimlf_premiumlist_detail.sql` menjadi `WPC DATE`. Nol kode Go menyentuh kolom itu di tabel baru |
| `STS_REJECT` | teks | **`NUMBER(38,0)`** | ✅ peta tipe dan tabel tiruan skema uji mengikuti katalog. Baris dokumen ini tidak diubah: `STS_REJECT` di tabel **baru** memang teks atas keputusan terpisah *(kode status, ADR-U-0022)*, dan yang diralat hanyalah pembacaan tabel **warisan** |
| `CLAIM_RETRO` | teks | **`NUMBER`** | ✅ **`[terbuka]` DICABUT 26-09-2026 malam** — `[keputusan work owner butir w2]`: ia **UANG**, dan `002_t_general_claim.sql` menjadi `NUMBER(38,8)`. Lihat blok di bawah |

⛔ **Yang ikut tersingkap:** fixture test mengisi `CLAIM_RETRO` dengan teks `"UJI-RETRO"`. Selama
tabel tiruan bertipe `VARCHAR2` semuanya, itu lolos; terhadap tabel yang berbentuk sama dengan
warisan, Oracle menjawab `ORA-01722`. Fixture sudah diperbaiki, dan satu test murni kini mengunci
bahwa setiap nilai yang menuju kolom angka berupa angka atau kosong.

⚠️ Produksi **belum** dibaca. Bentuk di produksi diasumsikan sama dengan pengembangan; DBA yang
dapat memastikannya.


---

## ✅ RALAT 26 September 2026 malam — `CLAIM_RETRO` adalah UANG

`[keputusan work owner butir w2]`. Ronde 5 menahannya: tipe `NUMBER` sudah diketahui dari katalog,
**artinya belum** — uang atau perbandingan. Menebaknya akan mengulang persis kesalahan yang baru
saja dibongkar, jadi ia ditahan satu ronde.

Buktinya kini ada, dan menunjuk satu arah:

| Bukti | Angka |
| --- | --- |
| Nilai di rentang 0–1 *(ciri perbandingan pecahan)* | **0** |
| Nilai di rentang 1–100 *(ciri persen)* | **0** |
| Nilai di atas 100 | **4.778** |
| Nilai berdesimal | **475** |
| Baris dengan `CLAIM_RETRO` < `CLAIM_AMOUNT` | **8.755** *(782 sama; 220 lebih besar — anomali data, dicatat)* |
| Label Pega di `AdjustmentDetail_Section` | *"Claim Retro"*, di samping label uang lain |

⚠️ Seluruhnya **agregat** instance pengembangan: cacah dan rentang saja, **nol baris data dibaca**.

**Akibatnya:** `002_t_general_claim.sql` `CLAIM_RETRO` menjadi **`NUMBER(38,8)`**; di Go ia
`models.Money`, bukan `Ratio` — keduanya sengaja bertipe berbeda supaya tidak pernah terjumlahkan
(ADR-F-0004). `BarisLamaDari` kini menuliskannya ke tiap baris datar, dan nilai itu dibaca kembali
saat klaim dibongkar dari tabel warisan.

⚠️ **`[terbuka]` baru yang lahir dari keputusan ini:** `T_GENERAL_CLAIM` **tidak punya kolom mata
uang**, sedangkan uang tanpa mata uang tidak bermakna. Saat klaim dibongkar dari tabel warisan,
mata uangnya diambil dari baris adjustment; saat header dibaca sendirian lewat `AmbilHeader`, ia
**kosong**. Menambahkan kolom mata uang ke header adalah keputusan tersendiri, dan executor tidak
mengarangnya. Pemilik: work owner.


---

## ⭐ KOLOM BARU 26 September 2026 malam — `T_GENERAL_CLAIM.CURRENCY`

`[keputusan work owner butir z1]`. Kolom ini **tidak** ada di korpus Pega maupun di tabel warisan;
ia lahir dari `[terbuka]` yang dicatat ronde 6.

**Sebabnya:** keputusan **w2** menjadikan `CLAIM_RETRO` **uang**, dan uang tanpa mata uang tidak
bermakna — `Money` yang `Currency`-nya kosong dapat dijumlahkan dengan uang mata uang lain tanpa
ada yang mencegah. Sampai kolom ini ada, `AmbilHeader` mengembalikan `ClaimRetro` bermata uang
kosong: benar secara mekanis, tidak berguna bagi pembacanya.

**Kenapa di header, bukan per baris:** agregat instance pengembangan menunjukkan **nol** klaim yang
baris-barisnya bermata uang campur. Satu kolom di tingkat header karena itu cukup. ⚠️ Dan bila
kelak ada yang campur, `BongkarBarisLama` sudah melaporkannya sebagai `Temuan` sejak ronde 6 —
tidak dipilih diam-diam.

Langkah migrasinya `002` yang disunting langsung, sah selama `T_MIGRASI` belum pernah ada di
instance mana pun (brief §2 l) — dan pada 26-09-2026 itu masih benar.


---

## ⭐ KOLOM ISI `T_CLAIMLF_DOCUMENT` — 26 September 2026 malam

`[keputusan work owner butir ad]`, tiket 03. Bab `T_CLAIMLF_DOCUMENT` di atas menulis *"kolomnya
tidak ditulis di sini"* dan menyebutnya `[data DBA]`; sejak sensus ini kolomnya **dapat**
diturunkan, dan langkah migrasi `010` menambahkannya.

### Dari mana kolomnya datang — dan kenapa BUKAN dari tempat yang diminta

Brief meminta sensus `.DocumentList` di `SaveOutStandingLife_Act`. Di sana hanya ada **satu**
rujukan, dan kelasnya **`Link-Attachment`** — lampiran bawaan Pega, bukan tabel karangan. Sensus
itu karena itu tidak menghasilkan kolom apa pun, dan melaporkannya sebagai "nol kolom" akan
menyesatkan.

Kolomnya datang dari `InsertDocument_Act`, yang milik kelas `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM` dan
memakai **`Obj-Save`** — Pega menulis SELURUH properti kelas itu, sehingga tidak ada daftar kolom
eksplisit sama sekali. Kolomnya **adalah** keempat belas kolom tabel warisan `DOCUMENT_CLAIM` yang
sudah dibaca dari katalog instance pengembangan (brief ronde 4 §9).

Tabelnya ditulis di bab `T_CLAIMLF_DOCUMENT` di atas, bukan di sini: bab berjudul
bukan-nama-tabel tidak dibaca pembanding kolom, dan tabel yang tidak terbaca pembanding adalah
tabel yang boleh salah tanpa ada yang tahu.

### ⚠️ Tujuh kolom warisan yang TIDAK ikut, dan sebabnya

`IDPEGA`, `INSKEY_LINK`, `INSKEY_DATA`, `PXCREATEOPERATOR` — kunci internal mesin lama; sistem baru
bukan Pega, dan kolom yang isinya pengenal Pega tidak berarti apa-apa di sini. `NOAKSEP` dan
`NOPREKAS` — nomor akseptasi dan prekas hidup di **baris adjustment**; menyalinnya ke dokumen
membuat satu nilai punya dua rumah.

### ⭐ Cacah kedua, dari sumber yang berbeda — 26 September 2026 malam

Aturan sensus CLAUDE.md §4a menuntut tiap angka dihitung **dua cara yang sungguh berbeda**. Cacah
pertama di atas berasal dari **katalog** `POOLDATA.DOCUMENT_CLAIM` (14 kolom). Cacah kedua berasal
dari **rule yang menulisnya**: seluruh `Property-Set` di `InsertDocument_Act` sebelum `Obj-Save`
`[terverifikasi]` baris 647-940 berkas pecahan — `.ID`, `.TANGGAL`, `.IDPEGA`, `.NAMAFILE`,
`.MIME`, `.KATEGORI_1`, `.KATEGORI_2`, `.NOAKSEP`, `.NOPREKAS`, `.PAYMENTDATE`,
`.pxCreateOperator`, ditambah `.T_STORAGE_ID` yang diisi activity storage dan menjadi precondition
`Obj-Save`-nya. Dua belas properti, seluruhnya termuat di keempat belas kolom katalog; dua kolom
katalog yang tersisa (`INSKEY_LINK`, `INSKEY_DATA`) memang tidak pernah disentuh rule ini. Kedua
cacah **cocok**, dan itulah yang membuat daftar tujuh kolom di atas bukan tebakan.

`Section/DocumentLife.xml` menguatkannya dari sisi layar: daftar dokumen beriterasi
`.DocumentClaimList` dan menampilkan `.NAMAFILE` dan `.KATEGORI_2`, dengan pratinjau bergerbang
`.MIME` dan `.T_STORAGE_ID`.

### ⭐ Aturan "dokumen lengkap" — mekanismenya TERBACA, daftarnya tidak

Diralat 26-09-2026 malam. Blok ini sebelumnya menyatakan aturannya tidak dapat diturunkan sama
sekali. Pembacaan ulang `SaveOutStandingLife_Act` langkah 12 menunjukkan **mekanismenya**:
`GetCategoryLife_SQL` menghasilkan daftar kategori **wajib**; tiap `.pyCategory` pada dokumen
peserta dikumpulkan ke `Category2`, **di-dedup** oleh satu langkah Java atas `.CARI1`; pesan
*"Documents are incomplete, please complete the documents"* muncul bila **cacah kategori berbeda
yang terunggah tidak sama dengan cacah kategori wajib**.

⛔ Yang tetap `[terbuka — DBA/work owner]`: **isi** daftar wajibnya. `GetCategoryLife_SQL` tidak ada
di korpus — seluruh 29 berkas `Claim Life/RDBList/` sudah dicacah. Juga terbuka: kolom mana pada
tabel relasional yang memegang kategori pembanding, sebab gerbang Pega membacanya dari daftar
**lampiran** (`.pyCategory`), bukan dari `KATEGORI_1`/`KATEGORI_2`.

### ⚠️ Gerbang dokumen yang KEDUA, yang blok ini dulu lewatkan

Langkah 3 `SaveOutStandingLife_Act` memakai daftar yang **berbeda** (`.DocumentClaimList`) dan
syarat yang berbeda: ia hanya berjalan bila `Type` **bukan** `TP`/`TR`, peserta `.IsAccept=="true"`,
dan peserta itu **nol dokumen**; pesannya menyebut **nomor urut** peserta.

---

## T_CLAIMLF_JEJAK

**Butir am**, A1 — jejak audit setiap transisi dan setiap jalur balik (ADR-U-0007).

⚠️ `ADJUSTMENT_ID` **tidak** ber-FK: jejak harus selamat dari penghapusan apa pun, dan jejak yang
ikut terhapus bersama yang dijejakinya bukan jejak.

⚠️ `DARI` dan `KE` memuat **dua kosakata** dalam satu kolom: kode status (`"0"`, `"1"`, `"2"`) pada
transisi baris, dan nama peran (`"ReasLifeAdmin"`) pada jalur balik tahap. Disengaja — keduanya
adalah "keadaan sebelum" dan "keadaan sesudah".

| Kolom | Tipe | Boleh kosong | Kunci | Catatan |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | dari `SEQ_CLAIMLF_JEJAK` (ADR-U-0006) |
| `ADJUSTMENT_ID` | teks | ya | index | baris yang dijejaki; unit keputusan adalah baris (ADR-U-0011) |
| `KLAIM_ID` | teks | ya | index | dipakai jalur balik tahap, yang tidak menunjuk baris |
| `DARI` | teks | ya | | keadaan sebelum |
| `KE` | teks | ya | | keadaan sesudah |
| `AKUN_ID` | teks | tidak | | identitas AKUN, bukan nama orang (ADR-U-0002) |
| `WAKTU` | TIMESTAMP | tidak | | KAPAN-nya |
| `KOMENTAR` | teks | ya | | migrasi 021 (OQ-M5, GILIRAN-17): alasan penolakan Admin — `Remarks` `RejectOSClaimLife_Sec` b1687; lebar 4000 dari preseden `KOMITE_COMMENT` |

---

## T_GENERAL_KOMITE

**Butir af**, A1/A2 — bentuknya dari `.scratch/komite-claim-life/STRUKTUR-TABEL-KOMITE-CLAIM-LIFE.md`,
yang §1 sebut sebagai sumbernya. Tiket 00 Komite kelak **memverifikasi**, bukan membuat ulang.

⛔ **Shared primary key:** `ID` = `T_WORK_CLAIM.ID` baris komite, teks **`KMTLF-xxxxxx`**.

⛔ **`ADJUSTMENT_ID` TANPA `REFERENCES`** — tabel ini **lintas-lini**, jadi satu kolom punya **dua**
tabel tujuan menurut `T_WORK_CLAIM.LINI`. Oracle hanya dapat menunjuk satu; keutuhannya **dijaga
kode Go**. `NOT NULL` dan ber-index **UNIK**: satu baris adjustment = tepat satu kasus komite.

| Kolom | Tipe | Boleh kosong | Kunci | Catatan |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK, FK | → `T_WORK_CLAIM.ID` *(shared PK)*, awalan `KMTLF-` |
| `ADJUSTMENT_ID` | teks | tidak | index unik | LIFE → `T_CLAIMLF_ADJUSTMENT.ID`; **tanpa** `REFERENCES` |
| `KOMITE_LOOP` | angka bulat | ya | | tinggi tangga — `KomitePostAdjustment` 1322 |
| `KOMITE_COUNT` | angka bulat | ya | | tingkat sekarang — 1398, 6172, 9028 |
| `ACCEPT_STATUS` | teks | ya | | `"1"` aksep / `"2"` tolak — gerbang 5695, 8119, 8648, 8887 |

---

## T_KOMITE_KOMITELIST

**Butir af**, A1/A2 — satu baris = satu anggota pada satu jenjang tangga.

⛔ **Tiga nama kolom dibetulkan** `[keputusan work owner 2026-09-18 sore]`, sebab nama korpusnya
**menyesatkan ke dua arah sekaligus**. `[terverifikasi]` sensus 555 berkas: `.KomiteID :=
.OPERATOR_ID` *(8 penulisan, 5 berkas)* dan `.IDKomite := .JABATAN` *(12 penulisan, 7 berkas)*.

| Nama korpus | Kolom | Isinya |
| --- | --- | --- |
| `KomiteID` | **`KOMITE_OPERATORID`** | akun operator |
| `IDKomite` | **`KOMITE_JABATAN`** | jabatan |
| `KomiteAproval` *(satu P)* | **`KOMITE_APPROVAL`** | ejaan dibetulkan — penyimpangan sadar |

⚠️ `KOMITE_OPERATORID` dan `KOMITE_EMAIL` memuat **data orang**. Dibaca saat jalan dari
`EMAILKOMITE`; nol baris disalin ke fixture, tiket, maupun log.

| Kolom | Tipe | Boleh kosong | Kunci | Catatan |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | dari `SEQ_KOMITE_KOMITELIST` (ADR-U-0006) |
| `DATA_KOMITE_ID` | teks | ya | FK, index | → `T_GENERAL_KOMITE.ID` |
| `KOMITE_URUT` | angka bulat | ya | | jenjang tangga — `DEGREE` roster, menaik |
| `KOMITE_OPERATORID` | teks | ya | | **data orang** — `CreateKMTLife_Act` 866, 972 |
| `KOMITE_JABATAN` | teks | ya | | `CreateKMTLife_Act` 952, 1041 |
| `KOMITE_EMAIL` | teks | ya | | **data orang** — `CreateKMTLife_Act` 932, 1014 |
| `KOMITE_APPROVAL` | teks | ya | | `"0"` saat roster dibentuk (912, 993); `"1"`/`"2"` saat diputus |
| `KOMITE_COMMENT` | teks | ya | | `KomitePostAdjustment` 893, 5945 |
| `DATE_APPROVE` | DATE | ya | | 935, 5965 |
