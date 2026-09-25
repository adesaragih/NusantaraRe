# Isi folder `.scratch/claim-non-prop-lapisan-data`

**Tanggal**: 19 September 2026 · **40 berkas** · satu README yang dibangkitkan, 39 tiket.

Folder ini **papan tiket** untuk pekerjaan lapisan data modul Claim Non Prop. Ia menyimpan **pekerjaannya** — apa yang harus dibangun, apa yang menahannya, dan apa yang sudah selesai.

> **Yang dihasilkan pekerjaan itu TIDAK ada di sini.** DDL, ADR, spec, temuan, dan kamus kolom hidup di `_migration-docs/claim-non-prop/`. Folder ini menyimpan **daftar kerjanya**, bukan **hasil kerjanya**.
>
> Bila Anda mencari "bagaimana bentuk tabelnya" atau "kenapa diputuskan begitu", Anda salah folder — mulailah dari `_migration-docs/claim-non-prop/KEADAAN-AKHIR.md`.

---

## 1. Keadaan papan

**aktif 0 · menunggu instance 8 · tertahan 1 · selesai 29 · mati 1 = 39**

```
cd issues
grep -h '^status:' *.md _selesai/*.md _tertahan/*.md _mati/*.md | sort | uniq -c
ls *.md _selesai/*.md _tertahan/*.md _mati/*.md | grep -v README | wc -l
```

Kedua perintah harus sepakat di **39**.

**Tidak ada tiket aktif, dan itu bukan berarti pekerjaannya kurang.** Seluruh tiket lapisan data selesai 19 September 2026. Yang tersisa **membutuhkan sesuatu yang belum ada** — instance Oracle yang dapat dipasangi, batch lapisan aplikasi, dan jawaban REQ-012.

---

## 2. Susunan berkas

```
.scratch/claim-non-prop-lapisan-data/
├── ISI-FOLDER.md          ← berkas ini
└── issues/
    ├── README.md          papan terbit — DIBANGKITKAN, jangan disunting tangan
    ├── 12…39 .md          26 tiket berstatus aktif atau menunggu instance
    ├── _selesai/          29 tiket selesai
    ├── _tertahan/         1 tiket tertahan
    └── _mati/             1 tiket dibatalkan
```

### `issues/README.md` dibangkitkan, bukan ditulis

```
python ../../_migration-docs/claim-non-prop/alat/papan.py
```

Alat itu **membaca** tiap berkas tiket dan **menulis satu berkas saja**: `issues/README.md`. Berkas tiketnya sendiri tidak pernah disentuh. Menyunting `README.md` dengan tangan akan hilang pada pembangkitan berikutnya — sunting **berkas tiketnya**, lalu jalankan ulang.

---

## 3. Empat aturan papan, dan kenapa masing-masing ada

**1. Status adalah field di dalam berkas, bukan lokasi folder.**

```
status: aktif | menunggu-instance | tertahan | selesai | selesai-sebagian | mati
```

Folder `_selesai/`, `_tertahan/`, `_mati/` hanya kerapian. Pembangkit membaca field-nya, dan **tidak pernah menghapus berkas tiket**. Aturan ini ada karena versi sebelumnya memakai folder sebagai status, dan pembangkitnya bisa menghilangkan tiket.

**2. Penahan dari luar papan ditulis dengan nama aslinya** — `REQ-032`, `ADR-0028`, `lingkup-batch` — tidak diterjemahkan jadi nomor tiket. Yang menahan dari luar **harus terlihat berasal dari luar**.

**3. Penahan yang sudah selesai dicoret, tidak dihapus**: `~~12~~ *(selesai)*`. Rantainya tetap terbaca.

**4. Nomor yang lompat bukan kekeliruan.** Tidak ada `27` dan `28` di daftar aktif karena keduanya sudah selesai. Penomoran **tidak pernah disusun ulang** — menyusunnya ulang memutus setiap rujukan yang sudah ada.

### Dua penahan yang dibedakan, dan bedanya penting

| Field | Artinya |
|---|---|
| `ditahan-luar: [x]` | **tidak boleh jalan** sampai `x` kembali |
| `menunggu-luar: [x]` | **boleh jalan.** Jawaban `x` mengubah satu hal yang tiketnya sudah sebut — sebuah constraint, sebuah jalur migrasi — bukan rancangannya |

Dipisahkan supaya papan **tidak berteriak serigala**. Tiga tiket ber-`menunggu-luar` REQ-018 semuanya selesai tanpa menunggunya.

---

## 4. Kedua puluh sembilan tiket selesai

### 4.1 Enam sapuan sumber — membaca sistem lama (`05`–`10`)

Tiket ini tidak membangun apa pun; ia **membaca ekspor XML** dan menutup pertanyaan.

| # | Asal | Apa yang dibacanya | Hasil yang paling berakibat |
|---|---|---|---|
| `05` | T-32 | kolom yang sungguh diterima Arasapas dan kasir | 25 medan `CARI` sisi kasir, belum pernah diinventarisasi (**D36**, **D46**) |
| `06` | T-33 | kelengkapan enum `CNPStatusCase` | dua dari empat nilai yang didaftar memori **tidak ada di mana pun** (**D35**) |
| `07` | T-34 | perilaku per nilai `PaymentType` | nilai `1`–`7` terbaca sebagai perbandingan **tanpa kutip** — penolakan lama saya dicabut (**D9**) |
| `08` | T-35 | padanan IDR untuk biaya penilaian, salvage, biaya lain | **sembilan nama tanpa padanan rupiah sama sekali** → ADR-0029 |
| `09` | T-36 | apa persisnya yang ditulis `EditXOLAlokasi` | `IsEditClaim` ditulis, **dibaca nol kali** (**D10**, **D28**) |
| `10` | T-37 | sapuan S5–S8 dan dua sumber yang belum pernah dibaca | `AlokasiXOLPaid` **nol penulis** di empat lapisan (**D11**) |

### 4.2 Dua keputusan bentuk (`27`, `28`)

`27` domain kolom keadaan · `28` pasangan IDR–kurs. Keduanya menetapkan pola yang lalu dipakai berulang di 22 tabel.

### 4.3 Dua puluh satu tiket rancangan — menghasilkan DDL

| # | Asal | Menghasilkan |
|---|---|---|
| `01` | T-01 | `00_SKEMA_DAN_AKUN.sql` — 2 tablespace, 3 akun, 3 peran |
| `03` | T-14 | **aturan penamaan final** — 105 objek ditulis ulang, satu tabel diganti nama |
| `04` | T-15 | index penopang setiap foreign key |
| `12` | T-03 | `01_KLAIM.sql` — aggregate root; kolom `LINI_USAHA` dan `PENUTUPAN_LAMA` |
| `13` | T-21 | `21_MIGRASI_PENDARATAN.sql` — gerbang `IS JSON` **dicabut** |
| `14` | T-22 | `20_MIGRASI_NILAI_DITOLAK.sql` — seam `ID_PENDARATAN` |
| `15` | T-23 | `19_MIGRASI_KORELASI.sql` — lima pengenal lama, dua index arah lama→baru |
| `16` | T-04 | `02_NILAI_KLAIM_MATA_UANG.sql` — **dua cacat constraint ditemukan di sini** |
| `17` | T-05 | `03_ALOKASI_LAYER.sql`, `04_RETENSI_CEDANT.sql` — kolom `_HITUNG`/`_SUNTING`/`_DIPAKAI` |
| `18` | T-06 | `05_AKSEPTASI.sql`, `06_ADJUSTMENT.sql`, `08_REKENING_PENERIMA.sql` |
| `19` | T-08 | `10`, `11`, `12` — masukan mesin alokasi tersimpan |
| `20` | T-09 | `09`, `13`, `14` — objek, kronologi, dokumen |
| `21` | T-10 | `15`, `16`, dan **`Z00_ISIAN_AWAL.sql`** — satu-satunya berkas DML |
| `22` | T-11 | `17_KOREKSI_NILAI.sql`, `18_KLAIM_PENJAGA_TANGGAL.sql` |
| `23` | T-19 | `V04_V_PARITAS_SHADOW.sql` — ditulis ulang dari 1 tabel jadi 19 |
| `25` | T-39 | `22_ARSIP_MUATAN_KELUAR.sql` — tulis-sekali, ditegakkan hak akses |
| `26` | T-07 | `07_PREMI_PEMULIHAN.sql` — seluruh masukan rumus di satu baris |
| `29` | T-16 | `V01`, `V02` — view kompatibilitas, **menghasilkan selisih** |
| `30` | T-17 | `V03_V_REKAP_KLAIM_MATA_UANG.sql` |
| `31` | T-18 | `V05`–`V09` — lima PageList lama jadi view |
| `35` | T-20 | `V00_HAK_AKSES.sql` — 44 pencabutan hak objek |

---

## 5. Sepuluh tiket yang belum selesai, dan sebabnya berbeda-beda

### 5.1 Delapan uji — `menunggu instance`

`24` `32` `33` `34` `36` `37` `38` `39`

**Rancangannya lengkap; yang kurang mesinnya.** Uji ini baru dapat hijau sesudah DDL benar-benar dijalankan di instance Oracle, dan itu di luar batch lapisan data.

Ditandai `menunggu-instance` — bukan `aktif` — supaya papan **tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun**. Itu pembedaan yang sengaja dibuat, bukan cara memperkecil angka.

> **Jangan jalankan `Z02_KUNCI_PEMILIK.sql` sebelum kedelapannya hijau.** Sesudah akun pemilik terkunci, setiap perbaikan kecil menjadi permintaan ke DBA.

### 5.2 Satu tertahan — `11`

Tipe desimal di sisi Golang. **Tertahan lingkup, bukan pengetahuan**: presisinya sudah tetap (`NUMBER(38,20)`, ADR-0003). Prasyaratnya tinggal satu — batch lapisan aplikasi dimulai.

### 5.3 Satu mati — `02`

"Tabel singkatan tertutup ditulis dan dibekukan". **Dibatalkan, bukan ditunda.**

Premisnya gugur: nama pengenal disangka menabrak batas 30 byte sehingga singkatan diperlukan. Yang menabrak batas ternyata **satu kebiasaan** — mencantumkan daftar kolom di nama constraint. Begitu itu berhenti, tekanannya hilang.

**Berkasnya tidak dihapus.** Tiket yang pernah ada adalah bukti premisnya **diperiksa**, bukan dilewatkan. Nomor `02` tidak dipakai ulang.

---

## 6. Yang masih menahan pelaksanaan

Empat REQ, seluruhnya di `_migration-docs/claim-non-prop/ORACLE-REQUESTS.md`. **Tidak satu pun menahan keputusan rancangan** — keempatnya menahan **pelaksanaan**:

| REQ | Menahan | Yang berubah bila jawabannya lain |
|---|---|---|
| REQ-018 | penanganan migrasi di `12`, `15`, `18`, `29`, `34` | constraint, **bukan** kolom |
| REQ-033 | pemetaan view di `29`, `32` | berapa banyak kelompok ditolak, bukan apakah penolakannya ada |
| REQ-021 | `35`, `39` | perlu-tidaknya pokok tersendiri untuk manajemen |
| REQ-037 | klaim ADR-0017 di `35` | apakah "satu pintu tulis" berlaku sungguhan |

---

## 7. Membaca satu berkas tiket

Tiap tiket berbentuk sama:

```
---
status: selesai
menunggu-luar: [REQ-018]
---

# 12: Klaim tersimpan sebagai aggregate root…

> catatan penutupan — apa yang berubah, dan apa yang ditemukan saat mengerjakannya

*Asal: `T-03` di `_migration-docs/claim-non-prop/TICKETS.md`.*

**What to build:** …        **Tidak termasuk:** …       **Blocked by:** …
**Dasar:** DECIDED(ADR-xxxx) / EVIDENCED: …
- [x] kriteria penerimaan
**Ketidakpastian:** …
```

**Bagian yang paling banyak dilewati orang, dan paling berharga: `Tidak termasuk` dan `Ketidakpastian`.** Keduanya menyatakan apa yang **sengaja** tidak dikerjakan beserta sebabnya — tanpa itu, orang berikutnya akan mengira ia menemukan pekerjaan yang terlupa, lalu mengerjakannya dan melemahkan rancangan.

Contohnya ada di `19`: `ESTIMASI_AWAL` sengaja **tidak** diberi `NOMOR_URUT`, karena kunci alaminya lebih kuat. Kriterianya sudah diubah supaya tidak terbaca sebagai pekerjaan yang belum selesai.

---

## 8. Ke mana selanjutnya

| Pertanyaan Anda | Berkas |
|---|---|
| Apa keadaan seluruh pekerjaan ini? | `_migration-docs/claim-non-prop/KEADAAN-AKHIR.md` |
| Bentuk tabelnya seperti apa? | `_migration-docs/claim-non-prop/ddl-usulan/` — 36 berkas |
| Kenapa diputuskan begitu? | `_migration-docs/claim-non-prop/docs/adr/` — 29 ADR |
| Apa yang belum diketahui? | `_migration-docs/claim-non-prop/_selesai/OPEN-QUESTIONS.md` |
| Bagaimana sampai ke sini? | `_migration-docs/claim-non-prop/RIWAYAT-SESI-2026-09-19.md` |
| Bentuk datar seluruh skema? | `_migration-docs/claim-non-prop/keluaran/datar-*.tsv` beserta `RINGKASAN-TABEL-DATAR.md` — **TURUNAN** dari `ddl-usulan/`, dibangkitkan `alat/buat-tabel-datar.py`. Untuk dibuka di spreadsheet dan disaring |
| Tabelnya berhubungan bagaimana? | `_migration-docs/claim-non-prop/ERD-KLAIMNP.html` — ERD kotak-entitas skema baru, 22 kotak dan 17 foreign key. **TURUNAN** dari `ddl-usulan/` |
| Rancangan tabel `T_*` Claim Non Prop | `_migration-docs/claim-non-prop/Diagram-Skema-Tabel-ClaimNonProp.xlsx` beserta `ERD-CLAIM-NON-PROP.html` — 30 tabel, hanya PK dan FK |
| Papan hari ini | `issues/README.md` |
