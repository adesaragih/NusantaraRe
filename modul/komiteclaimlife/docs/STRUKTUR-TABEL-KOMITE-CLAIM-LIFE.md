# Struktur Tabel — Komite Claim Life

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

Rule korpus yang dirujuk di berkas ini:

| Singkatan | Rule | Kelas |
| --- | --- | --- |
| **CreateKMT** | `Claim Life/Activity/CreateKMTLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` |
| **KomitePost** | `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` |
| **Comitee** | kelas data roster `ASM-FW-GCNMFW-Data-Comitee` | — |

⚠️ Seluruh kolom **nullable** kecuali PK dan yang disebut NOT NULL; wajib-isi ditegakkan di Go.

---

## T_WORK_CLAIM — baris komite

⚠️ **Tabel ini BUKAN milik modul Komite dan tidak didefinisikan ulang di sini.** Definisi penuhnya —
seluruh kolom, tipe, index, dan aturan `CHECK` — ada di
**`.scratch/claim-life/STRUKTUR-TABEL-CLAIM-LIFE.md` §`T_WORK_CLAIM`**. Ia tabel **lintas-lini**,
satu baris per work object.

Yang mengikat modul ini hanya **aturan baris komitenya**:

| Aturan | Isi | Sumber |
| --- | --- | --- |
| `ID` | berawalan **`KMTLF-`** (teks berformat `KMTLF-xxxxxx`) | keputusan tiket 00 Komite + **work owner 2026-09-18** |
| `COVER_KEY` | **wajib terisi** — `ID` baris klaim induknya (berawalan **`CLMLF-`**) | keputusan tiket 00 Komite + **work owner 2026-09-18** |
| `LINI` | **`LIFE`** — sama dengan baris klaim induknya | **work owner 2026-09-18** |

Baris komite **lahir saat penyerahan ke Komite diterima**; sebelum itu ia tidak ada.

⚠️ **`T_WORK_CLAIM` lintas-lini**, jadi definisi kolom `ID`-nya menyebut **delapan awalan** —
sepasang per lini:

| LINI | baris klaim | baris komite |
| --- | --- | --- |
| **FAC** | `CLM-` | `KMT-` |
| **PROP** | `CLMP-` | `TKMT-` |
| **NONPROP** | `CLMNP-` | `KMTNP-` |
| **LIFE** | `CLMLF-` | `KMTLF-` |

Berkas ini modul **Life**, jadi contoh dan ilustrasinya memakai pasangan Life saja
(`CLMLF-` / `KMTLF-`). Definisi penuhnya di `STRUKTUR-TABEL-CLAIM-LIFE.md`.

`T_GENERAL_KOMITE` memakai **`ID` yang sama persis** dengan baris komite ini — **shared primary
key**, tanpa kolom penyambung.

---

## T_GENERAL_KOMITE

Header kasus komite. Satu baris mewakili **satu kasus komite**, yaitu satu penyerahan sebuah baris
adjustment ke Komite. `ID`-nya **sama persis** dengan baris komite di `T_WORK_CLAIM` (shared PK).

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | keputusan tiket 00 Komite — shared PK = `T_WORK_CLAIM.ID` baris komite, berawalan **`KMTLF-`** |
| `ADJUSTMENT_ID` | teks | **tidak** | FK | keputusan tiket 00 Komite + **work owner 2026-09-18** — penutup lingkar, **dua tabel tujuan menurut `LINI`**: LIFE → `T_CLAIMLF_ADJUSTMENT.ID`, PROP → `T_CLAIMP_ADJUSTMENT.ID`. **Tanpa `REFERENCES`** |
| `KOMITE_LOOP` | bilangan bulat | ya | | keputusan tiket 00 Komite — jumlah tingkat tangga |
| `KOMITE_COUNT` | bilangan bulat | ya | | keputusan tiket 00 Komite — tingkat yang sedang berjalan |
| `ACCEPT_STATUS` | teks | ya | | keputusan tiket 00 Komite — hasil final: `1` aksep / `2` tolak |
| `KOMITE_USUL_TUTUP` | teks | tidak | | migrasi `komiteclaimprop/680` (RALAT 08-10-2026) — usul tutup klaim, `'1'`/`'0'` bawaan `'0'`; kasus Life tidak menulisnya |
| `KOMITE_USUL_CADANG` | teks | tidak | | migrasi `komiteclaimprop/680` (RALAT 08-10-2026) — usul cadangkan klaim, `'1'`/`'0'` bawaan `'0'`; kasus Life tidak menulisnya |

**Index:** `ADJUSTMENT_ID` **(biasa — `IX_GENERAL_KOMITE_ADJ`, migrasi `komiteclaimprop/681`)**.

> **RALAT 08-10-2026** (izin work owner 08-10-2026, Komite Claim Prop). Kalimat lama: *"**Index:** `ADJUSTMENT_ID` **(UNIK)**."* — ID adjustment Prop (`SEQ_T_CLAIM`) dan Life (`SEQ_CLAIMLF_ADJ`) sama-sama angka polos, sehingga indeks unik lintas lini suatu saat menolak penyerahan sah (ORA-00001). Keunikan satu adjustment ↔ satu kasus komite kini dijaga `KOMITE_ID` UNIK di `T_CLAIMLF_ADJUSTMENT` dan `T_CLAIM_ADJUSTMENT`. Dua kolom usul di atas ditambahkan migrasi `komiteclaimprop/680`.

**Relasi:**

- induknya `T_WORK_CLAIM` · **shared PK**, 1:1 · tidak ada kolom penyambung · ON DELETE **di Go**
- anaknya `T_KOMITE_KOMITELIST` lewat `DATA_KOMITE_ID` · 1:N · ON DELETE **CASCADE**
- menunjuk **tabel adjustment menurut `LINI`** lewat `ADJUSTMENT_ID` · **1:1** · ON DELETE **di Go** — LIFE → `T_CLAIMLF_ADJUSTMENT`, PROP → `T_CLAIMP_ADJUSTMENT`

⚠️ `[keputusan work owner]` **Satu baris `AdjustmentList` = TEPAT satu kasus komite, dan sebaliknya.**
Karena itu `ADJUSTMENT_ID` **NOT NULL** dan ber-index **UNIK**: setiap kasus komite wajib menunjuk
satu baris adjustment, dan dua kasus komite **tidak boleh** menunjuk baris adjustment yang sama.
Pasangannya di sisi klaim adalah `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` — nullable, index **UNIK**.

⚠️ **`ADJUSTMENT_ID` menunjuk DUA tabel, bukan satu.** `[keputusan work owner]` 2026-09-18 —
`T_GENERAL_KOMITE` lintas-lini, jadi satu kolom ini punya **dua tabel tujuan**:

| Lini kasus komite | `ADJUSTMENT_ID` menunjuk |
| --- | --- |
| **LIFE** | `T_CLAIMLF_ADJUSTMENT.ID` |
| **PROP** | `T_CLAIMP_ADJUSTMENT.ID` |

Tabel tujuannya ditentukan oleh **`T_WORK_CLAIM.LINI` pada baris yang sama**.

⛔ **Konsekuensinya, ditulis apa adanya:** Oracle hanya bisa memasang `REFERENCES` ke **satu**
tabel, jadi **klausa itu tidak dipasang sama sekali**. **Keutuhan `ADJUSTMENT_ID` dijaga kode Go,
bukan basis data.** Baris komite yang menunjuk adjustment **terhapus tidak akan ditolak** basis
data.

⚠️ **`NOT NULL` dan index `UNIK` TETAP** — keduanya masih bisa ditegakkan basis data, dan
keduanyalah yang menjamin **1:1**.

⚠️ `[terbuka]` Bila kelak **Claim Non Prop** punya tabel adjustment sendiri, tujuannya menjadi
**tiga**. Dicatat, **tidak dirancang sekarang**.


---

## T_KOMITE_KOMITELIST

Roster dan keputusan komite. Satu baris mewakili **satu anggota komite pada satu jenjang tangga**
dari satu kasus komite.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | keputusan tiket 00 Komite |
| `DATA_KOMITE_ID` | teks | ya | FK | keputusan tiket 00 Komite — → `T_GENERAL_KOMITE.ID` |
| `KOMITE_URUT` | bilangan bulat | ya | | keputusan tiket 00 Komite — jenjang tangga |
| `KOMITE_OPERATORID` | teks | ya | | korpus `KomiteID` — CreateKMT (`KomiteList(<APPEND>).KomiteID = .OPERATOR_ID`); **isinya akun operator** |
| `KOMITE_JABATAN` | teks | ya | | korpus `IDKomite` — kelas Comitee (`IDKomite := .JABATAN`); **isinya jabatan** |
| `KOMITE_EMAIL` | teks | ya | | korpus `KomiteEmail` — CreateKMT (`.KomiteEmail = .EMAIL`) |
| `KOMITE_APPROVAL` | teks | ya | | ⚠️ korpus `KomiteAproval` (**satu P**) — CreateKMT (`= 0` saat roster dibentuk), KomitePost (`= pyWorkPage.AcceptStatus`); nilai `0`/`1`/`2` |
| `KOMITE_COMMENT` | teks | ya | | korpus `KomiteComment` — KomitePost (`= pyWorkPage.Comment`) |
| `DATE_APPROVE` | DATE | ya | | korpus `DateApprove` — KomitePost (`= @CurrentDateTime()`) |

### Tiga nama kolom diganti — `[keputusan work owner]` 2026-09-18 sore

| Lama | Baru | Isinya |
| --- | --- | --- |
| `KOMITE_ID` | **`KOMITE_OPERATORID`** | akun **operator** (Pega `KomiteID := .OPERATOR_ID`) |
| `OPERATORID_KOMITE` | **`KOMITE_JABATAN`** | **jabatan** (Pega `IDKomite := .JABATAN`) |
| `KOMITE_APROVAL` | **`KOMITE_APPROVAL`** | ⚠️ ejaan **dibetulkan** |

`[terverifikasi]` Sensus **555 berkas korpus**: `.KomiteID := .OPERATOR_ID` (**8 penulisan, 5
berkas**) dan `.IDKomite := .JABATAN` (**12 penulisan, 7 berkas**). Dua properti bernama nyaris
sama, **isinya berbeda**. Kolom `OPERATORID_KOMITE` bersumber dari `IDKomite`, jadi namanya
**menyesatkan** — isinya jabatan, bukan id operator.

⚠️ **Jejak nama lengkap, jangan diputus:**
`ID_KOMITE` → `OPERATORID_KOMITE` *(keputusan 18 Sep pagi)* → **`KOMITE_JABATAN`**
*(keputusan 18 Sep sore)*. **Keputusan pagi dibatalkan.**

⚠️ **`KOMITE_APPROVAL` adalah penyimpangan sadar.** Korpus mengeja **`KomiteAproval` dengan
satu huruf P**. Ejaannya **sengaja dibetulkan** — berbeda dari preseden `COMMISION` dan
`STATUSS` yang dipertahankan apa adanya. **Jangan "merapikan" balik.** Skrip migrasi memetakan
properti Pega **`KomiteAproval` → kolom `KOMITE_APPROVAL`**.

**Index:** `DATA_KOMITE_ID`.

**Relasi:**

- induknya `T_GENERAL_KOMITE` lewat `DATA_KOMITE_ID` · 1:N · ON DELETE **CASCADE**
- tidak punya anak

### Tiga keputusan yang sebelumnya belum tertulis di mana pun

**8a · Jalur baca selalu berangkat dari adjustment, bukan sebaliknya.**
`[keputusan work owner]`

```
T_CLAIMP_ADJUSTMENT.KOMITE_ID   (berisi "TKMT-...")
   -> SELECT T_GENERAL_KOMITE      WHERE ID = KOMITE_ID
   -> SELECT T_KOMITE_KOMITELIST   WHERE DATA_KOMITE_ID = ID
```

Akibatnya **tidak ada query yang berangkat dari `ADJUSTMENT_ID`**, sehingga masalah dua tabel
tujuan di atas **tidak pernah menyentuh jalur baca** — ia murni soal **keutuhan data**.

**8b · Satu kolom tanggal persetujuan, bukan dua.** `[keputusan work owner]`
Di Pega, **satu langkah yang sama** mengisi `DateApprove` **dan** `DateApproval` dengan
`@CurrentDateTime()` — satu peristiwa ditulis dua kali, nilainya **identik sampai detik**.
`[terverifikasi]` Sensus **555 berkas**: `DateApprove` dibaca **5 berkas tampilan**;
`DateApproval` dibaca **NOL berkas tampilan**.

⛔ **Keputusan: simpan satu, `DATE_APPROVE`.** `DateApproval` **sengaja TIDAK dijadikan kolom**.
Jangan menambahkannya kembali karena "ada di korpus".

**8c · `T_GENERAL_KOMITE` ditinjau dan diterima apa adanya.** `[keputusan work owner]`
Kelima kolomnya diperiksa dan **tidak ada yang diubah namanya**. `ACCEPT_STATUS` **sengaja
dipertahankan di sini** — alasannya: supaya **hasil akseptasi komite terakhir terbaca dari satu
tempat**.

---

### Daftar penutup — bentuk final kedua tabel

**`T_GENERAL_KOMITE` — 5 kolom**

| Kolom | Keterangan |
| --- | --- |
| `ID` | PK, **shared PK** dengan baris komite di `T_WORK_CLAIM` |
| `ADJUSTMENT_ID` | FK **NOT NULL**, index **UNIK**, **dua tabel tujuan**, dijaga **Go** |
| `KOMITE_LOOP` | jumlah penyetuju yang dibutuhkan |
| `KOMITE_COUNT` | penyetuju ke berapa yang sedang berjalan |
| `ACCEPT_STATUS` | `1` setuju · `2` tolak |

**`T_KOMITE_KOMITELIST` — 9 kolom**

| Kolom | Keterangan | Asal Pega |
| --- | --- | --- |
| `ID` | PK | |
| `DATA_KOMITE_ID` | FK → `T_GENERAL_KOMITE.ID` | |
| `KOMITE_URUT` | penyetuju ke berapa | |
| `KOMITE_OPERATORID` | akun operator | `KomiteID` |
| `KOMITE_JABATAN` | jabatan | `IDKomite` |
| `KOMITE_EMAIL` | email | `KomiteEmail` |
| `KOMITE_APPROVAL` | ⚠️ `0` belum · `1` setuju · `2` tolak | `KomiteAproval`, **ejaan dibetulkan** |
| `KOMITE_COMMENT` | catatan penyetuju | `KomiteComment` |
| `DATE_APPROVE` | tanggal diputuskan | `DateApprove` |

`[data DBA]` **`POOLDATA.EMAILKOMITE` tetap tabel master yang DIBACA** — sumber `KomiteID`, email,
jabatan, dan pita nilai. Ia **bukan** tabel baru dan **tidak** diganti.

---

## Pohon relasi — sisi komite

```
T_WORK_CLAIM  (akar, lintas-lini, satu baris per work object)
  ID="CLMLF-123456"  COVER_KEY=NULL              <- baris KLAIM  (LINI=LIFE)
  ID="KMTLF-000789"  COVER_KEY="CLMLF-123456"    <- baris KOMITE (LINI=LIFE)
     |
     |  shared PK: ID yang sama persis, tanpa kolom penyambung
     v
T_GENERAL_KOMITE   ID="KMTLF-000789"
     |             FK ADJUSTMENT_ID -> tabel adjustment menurut LINI   1:1, NOT NULL, UNIK
     |                                 LIFE: T_CLAIMLF_ADJUSTMENT.ID
     |                                 PROP: T_CLAIMP_ADJUSTMENT.ID   (tanpa REFERENCES)
     |             KOMITE_LOOP · KOMITE_COUNT · ACCEPT_STATUS
     |
     +--1:N-- T_KOMITE_KOMITELIST                        ON DELETE CASCADE
              PK ID
              FK DATA_KOMITE_ID -> T_GENERAL_KOMITE.ID
              KOMITE_URUT · KOMITE_OPERATORID · KOMITE_JABATAN
              KOMITE_EMAIL · KOMITE_APPROVAL · KOMITE_COMMENT · DATE_APPROVE

  DIBACA dari luar, BUKAN anak, TIDAK ikut cascade:
     POOLDATA.EMAILKOMITE (master roster)
```

## Relasi yang menyentuh modul ini

| # | Induk | Anak | Kunci tamu | Kard | Hapus |
| --- | --- | --- | --- | --- | --- |
| 8 | `T_WORK_CLAIM` | `T_GENERAL_KOMITE` | **tidak ada kolom terpisah** — `T_GENERAL_KOMITE.ID` = `T_WORK_CLAIM.ID` baris komite (**shared PK**) | 1:1 | di Go |
| 9 | `T_GENERAL_KOMITE` | `T_KOMITE_KOMITELIST` | `DATA_KOMITE_ID` | 1:N | CASCADE |
| 10 | **tabel adjustment menurut `LINI`** — LIFE `T_CLAIMLF_ADJUSTMENT` · PROP `T_CLAIMP_ADJUSTMENT` | `T_GENERAL_KOMITE` | `ADJUSTMENT_ID` **NOT NULL, index UNIK, tanpa `REFERENCES`** | **1:1** | di Go |
| 11 | `T_CLAIMLF_ADJUSTMENT` | `T_WORK_CLAIM` | `KOMITE_ID` **nullable, index UNIK** | **1:1** | penunjuk |

Penomoran relasi mengikuti daftar sebelas relasi di
`.scratch/claim-life/STRUKTUR-TABEL-CLAIM-LIFE.md`.

---

## Catatan — belum ditetapkan, TIDAK menghambat berkas ini

- `[data DBA]` daftar kolom `DOCUMENT_CLAIM`
- `[data DBA]` presisi fisik seluruh kolom
- `[terbuka]` `REFERENCES T_WORK_CLAIM(ID)` pada `KOMITE_ID` dan `COVER_KEY`
- `[terbuka]` kolom nomor klaim (`NO_CLAIM`) dan nomor akseptasi (`NO_ACCEPTATION`) belum punya rumah
- `[terbuka]` isi `CASEID_POLICY`, `POLICY_NO`, `ENDORSMENT_NO`
- `[terbuka]` rumah lima kolom yatim bekas `T_CLAIM_POLICY`
- `[terbuka]` nama kolom audit (operator, tanggal) pada `T_GENERAL_KOMITE` — keputusan tiket 00 menyebutnya "+ audit (operator, tanggal)" tanpa menamainya, sehingga tidak ditulis di tabel atas
- `[terbuka]` generator nomor untuk **delapan awalan** (`CLM-` `CLMP-` `CLMNP-` `CLMLF-` / `KMT-` `TKMT-` `KMTNP-` `KMTLF-`) pada `T_WORK_CLAIM.ID`
- `[terbuka]` siapa yang memasang aturan **`LINI` baris komite == `LINI` baris `COVER_KEY`** di Go — `CHECK` tidak sanggup menjaganya
- `[terbuka]` bila Claim Non Prop kelak punya tabel adjustment sendiri, tujuan `ADJUSTMENT_ID` menjadi **tiga**
- `[terbuka]` tipe `KOMITE_LOOP`, `KOMITE_COUNT`, `KOMITE_URUT` sebagai bilangan bulat belum dikonfirmasi DBA
