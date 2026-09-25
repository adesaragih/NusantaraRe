# DDL 5 Tabel `_LIFE` — dari DBA (menutup OQ-001)

Tanggal: 2026-09-15
Sumber: **DBA / work owner** (dikirim langsung; bukan dari korpus Pega).
Status: `[data DBA]` — menutup **OQ-001** untuk Master Contract Retro Life.

> ✅ **KEADAAN FINAL (baca ini dulu).** Setelah pembaruan work owner di DB:
> - **PK `ID` ada di KELIMA tabel.**
> - **4 FK antar tabel ada, mode `ON DELETE CASCADE`** (contract→year, reinsurer→contract,
>   security→reinsurer, business→contract).
> - **NOT NULL: tetap nol** → wajib-isi ditegakkan di Go.
> - **Q5 FINAL = kaskade + popup konfirmasi Ya/Batal** (bukan "tolak"). FK CASCADE selaras.
>
> Beberapa kalimat di bawah ditulis pada tahap awal (sebelum PK/FK dipasang) dan memakai kata seperti
> "nol constraint" / "PK belum terkonfirmasi" / "tolak hapus" — **abaikan; keadaan final di atas yang
> berlaku.** Migrasi (tiket 12) tetap **verifikasi keadaan nyata**, bukan berasumsi.

---

## Temuan kunci

| # | Temuan | Implikasi |
| --- | --- | --- |
| 1 | **Kolom uang = `NUMBER`** (tanpa presisi) di `TREATYCONTRACT_LIFE`: `IDR, USD, B_IDR, B_USD, IDR_SELISIH, USD_SELISIH` | ✅ **ADR-0003 aman** — uang disimpan sebagai angka, bukan teks. `NUMBER` tanpa (p,s) = presisi penuh Oracle. Go pakai decimal presisi arbitrer |
| 2 | **Share/komisi = `NUMBER`**: `PCTSHARE, COMMISION, OVR_COMM` | ✅ desimal, bukan teks |
| 3 | ⚠️ **`RIRATE` = `VARCHAR2(1000)`** di `TREATYBUSINESS_LIFE` — **teks**, bukan NUMBER | Rate reasuransi disimpan sebagai string. Bawa sebagai teks (mungkin format "0.5%" / rate bertingkat); **jangan** paksa jadi angka tanpa konfirmasi. Catat kejanggalan |
| 4 | ⚠️ **NOL constraint** — tidak ada PK, FK, NOT NULL, UNIQUE, INDEX di kelima tabel; **semua kolom nullable** | Membuktikan grilling: **kaskade (Q5), anti-dobel, wajib-isi — semua di Go**. DB tak menjaga apa pun. Skema baru **boleh menambah** PK/FK/NOT NULL sebagai perbaikan sadar |
| 5 | **`USD` nullable** | ✅ konfirmasi: `USD` boleh kosong = **aturan sah**, bukan kelalaian (DB mengizinkan). AC 5 benar |
| 6 | **`ID` = `VARCHAR2(100)`. PK bervariasi per tabel** — lihat tabel PK di bawah | `TREATYBUSINESS_LIFE` **sudah punya PK** (`TREATYBUSINESS_LIFE_PK`, unique index, ENABLE VALIDATE). 4 tabel lain: PK belum terkonfirmasi (DDL awal tak menyertakan ALTER ADD CONSTRAINT) → OQ kecil |
| 7 | **`REINSTYPENAME`/`REINSURERNAME`/`BIZNAME` = `VARCHAR2(1000)`** | Lebar untuk cache nama (Q1). `USERID` bervariasi: 100 di year/contract/business, 1000 di reinsurer/security |
| 8 | **`TGLUPDATE, STARTDATE, ENDDATE, TREATYSTARTDATE, TREATYENDDATE` = `DATE`** | ✅ tanggal bertipe DATE — **menguatkan Q8**: gerbang tahun bisa banding tahun dari nilai DATE, bukan substring teks. Kolom tanggal memang DATE, jadi substring Pega memang keliru |
| 9 | **Sequence `TREATYYEAR_LIFE_SEQ` START WITH 44** | ID year baru berikutnya = `'1'+lpad(44,6)` = **`1000044`**. Format ID terkonfirmasi |

---

## Kolom per tabel (dari DDL)

### `TREATYYEAR_LIFE`
`ID VARCHAR2(100)`, `TREATYYEAR VARCHAR2(100)`, `UNDERWRITINGYEAR VARCHAR2(100)`,
`USERID VARCHAR2(100)`, `TGLUPDATE DATE`, `STARTDATE DATE`, `ENDDATE DATE`.

### `TREATYCONTRACT_LIFE`
`ID VARCHAR2(100)`, `IDTREATYYEAR VARCHAR2(100)`, `REINSTYPEID VARCHAR2(100)`,
`REINSTYPENAME VARCHAR2(1000)`, `USERID VARCHAR2(100)`, `TGLUPDATE DATE`,
`IDR NUMBER`, `USD NUMBER`, `B_IDR NUMBER`, `B_USD NUMBER`, `IDR_SELISIH NUMBER`, `USD_SELISIH NUMBER`,
`TREATYENDDATE DATE`, `TREATYSTARTDATE DATE`.

### `TREATYREINSURER_LIFE`
`ID VARCHAR2(100)`, `TREATYYEARID VARCHAR2(100)`, `TREATYCONTRACTID VARCHAR2(100)`,
`REINSTYPEID VARCHAR2(100)`, `REINSTYPENAME VARCHAR2(1000)`, `REINSURERNAME VARCHAR2(1000)`,
`PCTSHARE NUMBER`, `COMMISION NUMBER`, `OVR_COMM NUMBER`, `USERID VARCHAR2(1000)`, `TGLUPDATE DATE`,
`REINSURERID VARCHAR2(100)`.

### `TREATYSECURITYREINSURER_LIFE`
`ID VARCHAR2(100)`, `TREATYYEARID VARCHAR2(100)`, `TREATYCONTRACTID VARCHAR2(100)`,
`TREATYREINSURERID VARCHAR2(100)`, `REINSURERID VARCHAR2(100)`, `REINSURERNAME VARCHAR2(1000)`,
`PCTSHARE NUMBER`, `USERID VARCHAR2(1000)`, `TGLUPDATE DATE`.

### `TREATYBUSINESS_LIFE`
`ID VARCHAR2(100)`, `TREATYYEARID VARCHAR2(100)`, `TREATYYEAR VARCHAR2(100)`,
`REINSTYPEID VARCHAR2(100)`, `REINSTYPENAME VARCHAR2(1000)`, `BIZCODE VARCHAR2(100)`,
`BIZNAME VARCHAR2(1000)`, `USERID VARCHAR2(100)`, `TGLUPDATE DATE`,
`TREATYCONTRACTID VARCHAR2(100)`, `RIRATEID VARCHAR2(100)`, `RIRATE VARCHAR2(1000)`.

Tablespace `TBS_POOLDATA`.

### Status PK per tabel (final — dikonfirmasi work owner)

| Tabel | PK pada `ID` | Bukti |
| --- | --- | --- |
| `TREATYBUSINESS_LIFE` | ✅ ADA | `[data DBA]` DDL (`TREATYBUSINESS_LIFE_PK`, unique index, ENABLE VALIDATE) |
| `TREATYYEAR_LIFE` | ✅ ADA | `[keputusan work owner]` — ditambahkan work owner |
| `TREATYCONTRACT_LIFE` | ✅ ADA | `[keputusan work owner]` |
| `TREATYREINSURER_LIFE` | ✅ ADA | `[keputusan work owner]` |
| `TREATYSECURITYREINSURER_LIFE` | ✅ ADA | `[keputusan work owner]` |

✅ **Kelima tabel kini punya PK pada `ID`** (perbaikan sadar 7 — selesai di DB oleh work owner).
✅ **FK antar tabel SUDAH ditambahkan** (perbaikan sadar 8 — selesai di DB oleh work owner):

| FK | Kolom anak | → | Induk (`.ID`) |
| --- | --- | --- | --- |
| 1 | `TREATYCONTRACT_LIFE.IDTREATYYEAR` | → | `TREATYYEAR_LIFE` |
| 2 | `TREATYREINSURER_LIFE.TREATYCONTRACTID` | → | `TREATYCONTRACT_LIFE` |
| 3 | `TREATYSECURITYREINSURER_LIFE.TREATYREINSURERID` | → | `TREATYREINSURER_LIFE` |
| 4 | `TREATYBUSINESS_LIFE.TREATYCONTRACTID` | → | `TREATYCONTRACT_LIFE` |

`NOT NULL`: tetap tak ada di DB → wajib-isi tetap ditegakkan di Go.

✅ **FINAL: FK dipasang `ON DELETE CASCADE`** oleh work owner — **selaras dengan Q5 kaskade**
(hapus induk → anak ikut terhapus, dengan popup konfirmasi Ya/Batal di Go). Tidak ada lagi
pertentangan: Q5 = kaskade (bukan tolak), FK CASCADE mendukungnya.

→ OQ kecil PK & FK: **DITUTUP** — PK kelima tabel ada, 4 FK ada, mode `ON DELETE CASCADE`.

### Sequence
`TREATYYEAR_LIFE_SEQ` START WITH 44, MINVALUE 1, NOCACHE, NOCYCLE, NOORDER.
(4 sequence lain dirujuk procedure: `TREATYCONTRACT_LIFE_seq`, `TREATYREINSURER_LIFE_SEQ`,
`TREATYSECURITYREINSURER_LIFE_SEQ`, `TREATYBUSINESS_LIFE_SEQ` — DDL tak dikirim, tapi keberadaannya
pasti karena procedure memakainya.)

---

## Perbaikan sadar BARU yang muncul dari DDL (untuk skema Oracle target)

| # | Perbaikan | Alasan |
| --- | --- | --- |
| 7 | **`ID` dijadikan PRIMARY KEY** di kelima tabel | DDL lama `ID` VARCHAR2(100) tanpa PK — keunikan tak dijamin; upsert `WHERE ID=` rawan |
| 8 | **FK berkaskade eksplisit** (contract→year, reinsurer→contract, security→reinsurer, business→contract) — atau minimal FK penjaga | DDL lama nol FK; ini yang menyebabkan risiko yatim (Q5). Perbaikan sadar mendukung penegakan Q5 di dua lapisan |
| 9 | Pertimbangkan **`RIRATE` tetap teks** kecuali work owner memastikan formatnya angka | DDL: `RIRATE VARCHAR2(1000)`; jangan asal ubah tipe |

> Catatan: Q5 FINAL = **kaskade + popup konfirmasi Ya/Batal** (bukan tolak). Popup + hitung anak
> **tetap wajib di Go** (spec); FK `ON DELETE CASCADE` = lapis kedua yang menjalankan penghapusan
> anak, bukan pengganti popup.

## OQ setelah ini

- **OQ-001** → ✅ **DITUTUP**. Tipe kolom (uang/share=NUMBER, rate=VARCHAR2), nullability (semua
  nullable, `USD` boleh kosong = sah), constraint (nol) semuanya diketahui. **Tiket migrasi (tiket 12)
  naik ke `ready`.**
- **OQ-002** → ✅ DITUTUP (sesi sebelumnya).
- **OQ kecil** (nama tabel fisik REINSURANCETYPE, isi `.Code`/`.Type`) → masih terbuka, tidak
  memblokir tiket mana pun.
