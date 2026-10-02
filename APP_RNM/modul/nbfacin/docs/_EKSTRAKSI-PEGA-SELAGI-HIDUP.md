# Permintaan Ekstraksi dari Pega — **selagi sistem masih hidup**

**Kepada:** Tim IT · Tim DBA  ·  **Tanggal:** 15 September 2026
**Konteks:** migrasi Facultative Inward dari Pega ke Go + React + Oracle.

> ## ⏳ Mengapa ini mendesak
>
> Enam butir di bawah **hanya dapat diambil selama Pega masih berjalan**. Setelah sistem dimatikan,
> tidak ada sumber penggantinya — bukan "sulit", melainkan **tidak mungkin**. Tidak satu pun dapat
> direkonstruksi dari ekspor rule yang sudah kami terima.
>
> Tanpa keenamnya, sistem baru **tidak dapat direkonsiliasi** terhadap sistem lama, dan sebagian
> logika harus ditebak — yang dilarang oleh aturan proyek.

Seluruh angka di dokumen ini dihitung langsung dari korpus ekspor yang ada di tangan kami; perintah
auditnya tersedia bila diperlukan. **Dokumen ini tidak memuat nama orang maupun data pelanggan.**

---

## 1 · Ekspor produksi **tunggal**, satu titik waktu, tiga siklus sekaligus

**Pemilik: IT (Pega)**

`[terverifikasi]` Ekspor yang kami terima **tidak sepadan antar folder**. Dari 1.719 rule bernama
sama yang hadir di NB dan Endorsement, **95 berbeda versi ruleset** — dan arahnya **tidak
konsisten**: 62 salinan Endorsement lebih tua, 33 justru lebih baru.

Contoh konkret — `Activity\SetToInbox_ACT.xml`:

| Folder | Versi ruleset | Commit |
| --- | --- | --- |
| `NB FacIn` | `01-01-95` | 2026-08-07 |
| `RNW Fac In` | `01-01-95` | 2026-08-07 |
| `Endorsment Fac In` | **`01-01-87`** | **2025-09-18** |

Selisih **11 bulan pada rule yang sama**, padahal ketiga folder diekspor dalam rentang 8 hari. Salinan
Endorsement bahkan kehilangan satu cabang routing yang ada di salinan NB.

**Konsekuensi bila tidak diperbaiki:** selisih versi **tidak dapat dipisahkan** dari perbedaan bisnis
yang sesungguhnya. Artinya rekonsiliasi paralel run **mustahil** — setiap selisih angka akan punya
dua kemungkinan penyebab dan kami tidak dapat menentukan mana.

**Yang kami minta:** satu ekspor dari **lingkungan produksi**, diambil pada **satu titik waktu**,
mencakup **ketiga siklus** sekaligus. Ini juga menyelesaikan butir 2.

⚠️ **Butir ini sekaligus menjadi alat verifikasi untuk sebelas activity** yang dirujuk tetapi tidak
ada di ekspor yang kami terima (daftar namanya di `00-KEPUTUSAN-WORK-OWNER.md` K-006). Kami tidak
meminta ekspor khususnya — tetapi karena ketiadaan sebuah rule di ekspor **bukan bukti** rule itu
tidak dipakai (ekspor kami berasal dari titik waktu berbeda antar folder), enam cabang alur yang
memanggilnya kami **tangguhkan**, bukan buang. Ekspor produksi tunggal akan menjawabnya: bila
activity-nya ada di sana, cabang-cabang itu wajib kami port. Di antaranya **jalur Special Acceptance**
dan **tangga akseptasi putaran kedua** — keduanya inti aplikasi.

---

## 2 · `DecisionTable/IsUWAccepted` — baris pemetaan yang tidak ikut terekspor

**Pemilik: IT (Pega)**

`[terverifikasi]` Rule ini mendeklarasikan **enam hasil** — `confirm` · `reject` · `ask` · `banding` ·
`revise` · `decline` — sementara Data Transform menulis nilai numerik `1`, `2`, `3`, `4`, `7`, `9`
ke properti `ProposalAcceptStatus`. **Baris yang memetakan angka ke hasil tidak ikut terekspor.**

Bahwa menebaknya berbahaya sudah terbukti: nilai `4` **ditulis** oleh rule bernama
`SetBandingProposal_DT` (banding) tetapi **diuji** oleh `When/IsFacout` (fac out) — dua nama pemakai
dengan tema berlawanan untuk satu nilai yang sama.

⚠️ Isi `IsUWAccepted` juga **berbeda** antara siklus endorsement dan NB/RNW, jadi ekspornya
dibutuhkan **untuk ketiga siklus**.

**Konsekuensi bila hilang:** arah **setiap** keputusan underwriting di seluruh node persetujuan tidak
dapat ditentukan. Ini butir yang memblokir paling luas di seluruh proyek.

**Yang kami minta:** ekspor `IsUWAccepted` **beserta seluruh barisnya** (tampilan tabel keputusan
lengkap, bukan hanya header), untuk ketiga siklus. Berlaku sama untuk **11 DecisionTable lain** yang
baris hasilnya juga tidak terekspor, termasuk `isApproved`.

---

## 3 · `ALL_SOURCE` untuk 35 prosedur/fungsi `POOLDATA`

**Pemilik: DBA**

`[terverifikasi]` **Seluruh jalur tulis ke produksi melewati prosedur `POOLDATA`**, dan tidak satu
pun badan prosedurnya ada di korpus. Kami hanya melihat nama dan parameter panggilannya.

35 objek yang terbukti dipanggil sebagai prosedur/fungsi (dideteksi dari pola `POOLDATA.NAMA(`):

```
ERRORFACINPROD           FACINFORBACKUP            FACINLIFE
FACINOFFER               FACINOFFERLIFE            FACINPERFORMANCE
FACINSPREADLIFE          GENERATE_FACRETRO_NO      GET_TOKEN_STORAGE
GETCURRENCYSTANDARD      HISTORYAKSEPTASIPRODUCTION INSERTJSONPOLIS
INSERTJSONPOLISMONITORING INSERTUPDATECEDINGPRODUCTION INSERTUPDATERISKADDRESS
MONITORING_PROD_LOG      PEGA_DELETE_ERROR_KONVERSI PEGA_JSON_POLIS_TREATYIN
PEGA_M_ACCUMULATION_LIFE PEGA_M_JSON_OFFER         PEGA_MARKETINGOFFICER
PEGA_TREATY_IN           PROC_GENERATE_SEQUENCE_NUMBER PROSESCOPY
RDBINSERTCLIENT          RDBMASTERACCUMULATEDTYPE  RDBMASTERACCUMULATION
RDBMASTERBRANCH          RDBMASTERCITY             RDBMASTERDISTRICT
RDBMASTERNATION          RDBMASTERPROVINCE         RDBMASTERRW
RDBMASTERSHIP            TREATYPRODUCTION_BACKUP
```

Yang paling kritis: `PEGA_M_JSON_OFFER` (satu-satunya jalur persistensi offer), `INSERTJSONPOLIS`,
`INSERTUPDATECEDINGPRODUCTION`.

**Konsekuensi bila hilang:** logika penulisan produksi **tidak dapat direplikasi**, dan angka
produksi sistem baru tidak akan cocok dengan sistem lama.

**Yang kami minta:** keluaran `ALL_SOURCE` (atau `DBMS_METADATA.GET_DDL`) untuk ketiga puluh lima
objek di atas. **Definisi kode saja — tanpa data.**

---

## 4 · Isi 7 tabel `M_LIMIT_*` + ejaan pasti kolom `JABATAN`

**Pemilik: DBA**

`[terverifikasi]` Mekanisme tangga akseptasi **sudah kami pahami penuh** dari rule — query, urutan,
cara berhenti. Yang tidak ada adalah **datanya**.

Tujuh tabel limit per lini bisnis:

```
M_LIMIT_PROPERTYY                        M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL
M_LIMIT_PROPERTY_NON_PREFERREDD          M_LIMIT_NONPROPANDENGG
M_LIMIT_ENGINEERINGG                     M_LIMIT_FINANCIALINS
M_LIMIT_LIFE
```

Kolom yang perlu: `JABATAN`, `JABATAN_ATASAN`, `LIMIT_BOTTOM`, `LIMIT_BOTTOM2`, `MAX_LIMIT_IDR`,
`MAX_LIMIT_USD`, `BATAS_WAKTU`, `LOGIN` — **beserta DDL-nya** (tipe kolom), karena korpus memuat
**nol DDL**.

⚠️ **Yang paling menentukan: ejaan pasti nilai kolom `JABATAN`.** Token routing di dalam rule ditulis
**tanpa spasi** (`DIREKTURTEKNIK`), sementara ejaan di Oracle belum kami ketahui. Bila keduanya tidak
cocok, **tidak ada approver yang pernah ditemukan** dan seluruh tangga macet.

**Konsekuensi bila hilang:** tangga akseptasi **tidak dapat direkonsiliasi sama sekali**.

**Yang kami minta:** ekstraksi isi ketujuh tabel (data master wewenang, **bukan data pelanggan**) +
DDL-nya.

---

## 5 · Riwayat akseptasi — `HISTORYAKSEPTASIPEGA` dan `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK`

**Pemilik: DBA + IT**

`[terverifikasi]` Kedua tabel ini **dibaca untuk mengambil keputusan alur**, bukan sekadar audit:

- `GetAksepBanding_SQL` mengambil workbasket terakhir dari riwayat lalu men-join ke tabel limit untuk
  menentukan **jalur banding**.
- `GetFlagReject_SQL` menghitung baris berstatus reject untuk menentukan **flag penolakan**.
- Dua query membaca `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK` untuk **routing**.

⚠️ `DATAPEGA.PC_*` adalah **tabel internal Pega** — ia **hilang bersama Pega**. Ini satu-satunya butir
di dokumen ini yang datanya musnah, bukan sekadar sulit diambil.

**Konsekuensi bila hilang:** case yang sedang berjalan kehilangan jalur bandingnya, dan riwayat
keputusan untuk case historis tidak dapat direkonstruksi. **Tidak dapat dibatalkan.**

**Yang kami minta:**

1. **Ekstraksi penuh** `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK` + DDL-nya, sebelum dekomisioning.
2. DDL `HISTORYAKSEPTASIPEGA`, dan konfirmasi **apakah tabel ini pernah diarsip atau dipangkas**.
   Bila pernah — pemangkasan itu **sudah mengubah keputusan alur**, dan kami perlu tahu sejak kapan.
3. Penjelasan siapa yang mengisi kolom `ID_KOMITE`; tidak ada rule di korpus yang menulisinya.

---

## Ringkasan permintaan

| # | Butir | Pemilik | Sifat |
| ---: | --- | --- | --- |
| 1 | Ekspor produksi tunggal, 1 titik waktu, 3 siklus | IT | ekspor rule |
| 2 | `IsUWAccepted` + 11 DecisionTable lain, **beserta barisnya** | IT | ekspor rule |
| 3 | `ALL_SOURCE` 35 prosedur `POOLDATA` | DBA | definisi kode |
| 4 | Isi + DDL 7 tabel `M_LIMIT_*`, ejaan `JABATAN` | DBA | data master + DDL |
| 5 | Ekstraksi `DATAPEGA.PC_HISTORY_*` + DDL `HISTORYAKSEPTASIPEGA` | DBA + IT | **data, musnah bersama Pega** |

**Butir 1 menyelesaikan sebagian besar butir 2 sekaligus** — bila satu ekspor produksi penuh dapat
diambil, itu langkah paling efisien.

**Butir 5 adalah satu-satunya yang datanya benar-benar musnah.** Bila hanya satu yang dapat
dikerjakan hari ini, kerjakan butir 5.

---

### Sudah ditarik dari permintaan ini

**11 activity yang dirujuk tetapi berkasnya tidak ada** — semula butir 1. **Tidak diminta sebagai
permintaan ekspor tersendiri**: work owner menyatakan kesebelasnya sudah tidak digunakan (K-006,
15 September 2026). Daftar namanya tercatat di `00-KEPUTUSAN-WORK-OWNER.md` K-006.

⚠️ Penarikan ini **bukan** berarti alur yang memanggilnya dibuang. Enam cabang pemanggil —
termasuk **Special Acceptance** dan **tangga akseptasi putaran kedua** — berstatus **ditangguhkan**
dan akan diverifikasi terhadap ekspor produksi tunggal (butir 1). Lihat butir 1.

---

*Tidak ada nama orang, kredensial, endpoint, maupun data pelanggan di dokumen ini. Permintaan data
terbatas pada data master wewenang dan riwayat keputusan internal.*
