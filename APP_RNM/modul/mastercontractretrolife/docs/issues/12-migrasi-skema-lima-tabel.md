# 12: Migrasi skema — DDL lima tabel, PK, FK, dan sequence

**Status:** ready-for-agent — ✅ **OQ-001 ditutup** `[data DBA]`

**Blocked by:** None (can start immediately)

## Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin skema Oracle sistem baru menyimpan setiap nilai retro life
**tanpa berubah satu digit pun**, dengan **kunci primer dan kunci asing** yang menjamin keunikan dan
keterhubungan — sehingga rekonsiliasi tidak menemukan selisih dan tidak ada baris yatim yang mungkin.
*(User story 36–38 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL kelima tabel + PK + FK + sequence |
| `internal/repository` | Pemetaan tipe kolom ↔ desimal presisi arbitrer; `RIRATE` sebagai teks |
| — | Skrip rekonsiliasi |

## Tabel dan tipe `[data DBA]`

Sumber: `.scratch/master-contract-retro-life/ddl-tables-from-dba.md`. Tablespace `TBS_POOLDATA`.

| Tabel | Kolom kunci |
| --- | --- |
| `POOLDATA.TREATYYEAR_LIFE` | `ID`, `TREATYYEAR`, `UNDERWRITINGYEAR`, `USERID`, `TGLUPDATE`, `STARTDATE`, `ENDDATE` |
| `POOLDATA.TREATYCONTRACT_LIFE` | `ID`, `IDTREATYYEAR`, `REINSTYPEID`, `REINSTYPENAME`, `USERID`, `TGLUPDATE`, `IDR`, `USD`, `B_IDR`, `B_USD`, `IDR_SELISIH`, `USD_SELISIH`, `TREATYSTARTDATE`, `TREATYENDDATE` |
| `POOLDATA.TREATYREINSURER_LIFE` | `ID`, `TREATYYEARID`, `TREATYCONTRACTID`, `REINSTYPEID`, `REINSTYPENAME`, `REINSURERID`, `REINSURERNAME`, `PCTSHARE`, `COMMISION`, `OVR_COMM`, `USERID`, `TGLUPDATE` |
| `POOLDATA.TREATYSECURITYREINSURER_LIFE` | `ID`, `TREATYYEARID`, `TREATYCONTRACTID`, `TREATYREINSURERID`, `REINSURERID`, `REINSURERNAME`, `PCTSHARE`, `USERID`, `TGLUPDATE` |
| `POOLDATA.TREATYBUSINESS_LIFE` | `ID`, `TREATYYEARID`, `TREATYYEAR`, `TREATYCONTRACTID`, `REINSTYPEID`, `REINSTYPENAME`, `BIZCODE`, `BIZNAME`, `RIRATEID`, `RIRATE`, `USERID`, `TGLUPDATE` |

| Kelompok | Tipe `[data DBA]` | Catatan |
| --- | --- | --- |
| Uang: `IDR`, `USD`, `B_IDR`, `B_USD`, `IDR_SELISIH`, `USD_SELISIH` | **`NUMBER`** (tanpa presisi) | ✅ **ADR-0003 aman** — angka presisi penuh Oracle, bukan teks |
| Share/komisi: `PCTSHARE`, `COMMISION`, `OVR_COMM` | **`NUMBER`** | desimal |
| ⚠️ `RIRATE` | **`VARCHAR2(1000)`** | **teks** — bawa apa adanya |
| Tanggal: `TGLUPDATE`, `STARTDATE`, `ENDDATE`, `TREATYSTARTDATE`, `TREATYENDDATE` | **`DATE`** | menguatkan gerbang tahun (tiket 03) |
| Identitas: `ID`, `*ID` | `VARCHAR2(100)` | |
| Nama: `REINSTYPENAME`, `REINSURERNAME`, `BIZNAME` | `VARCHAR2(1000)` | lebar untuk cache nama |
| `USERID` | `VARCHAR2(100)` di year/contract/business; `VARCHAR2(1000)` di reinsurer/security | **tidak seragam** |
| Nullability | **semua kolom nullable, nol `NOT NULL`** | wajib-isi ditegakkan **di Go** |

## Kunci dan sequence `[data DBA]`

⚠️ **Penyimpangan sadar 7 — `ID` menjadi PRIMARY KEY di kelima tabel.** Semula hanya
`TREATYBUSINESS_LIFE` yang punya (`TREATYBUSINESS_LIFE_PK`, unique index, `ENABLE VALIDATE`); empat
lainnya tidak — padahal upsert berkunci `ID`.

⚠️ **Penyimpangan sadar 8 — FK antar tabel, mode `ON DELETE CASCADE`** (selaras kaskade tiket 09):

| FK | Kolom anak | → induk |
| --- | --- | --- |
| 1 | `TREATYCONTRACT_LIFE.IDTREATYYEAR` | `TREATYYEAR_LIFE.ID` |
| 2 | `TREATYREINSURER_LIFE.TREATYCONTRACTID` | `TREATYCONTRACT_LIFE.ID` |
| 3 | `TREATYSECURITYREINSURER_LIFE.TREATYREINSURERID` | `TREATYREINSURER_LIFE.ID` |
| 4 | `TREATYBUSINESS_LIFE.TREATYCONTRACTID` | `TREATYCONTRACT_LIFE.ID` |

`[data DBA]` **Sequence**: `TREATYYEAR_LIFE_SEQ` `START WITH 44`, `MINVALUE 1`, `NOCACHE`,
`NOCYCLE`, `NOORDER` → identitas tahun berikutnya `1000044`. Empat sequence lain dirujuk procedure:
`TREATYCONTRACT_LIFE_seq`, `TREATYREINSURER_LIFE_SEQ`, `TREATYSECURITYREINSURER_LIFE_SEQ`,
`TREATYBUSINESS_LIFE_SEQ`.

## ADR terkait

**ADR-0003** (uang non-float — `NUMBER` tanpa presisi → desimal presisi arbitrer di aplikasi),
**ADR-0006** (identitas dari sequence basis data), **ADR-0009** (migrasi penuh; koeksistensi
ditolak).

## Acceptance criteria

- [ ] DDL kelima tabel target menyalin **tipe dan panjang kolom** dari sumber apa adanya; **tidak
      ada** kolom uang yang menjadi `FLOAT`/`BINARY_DOUBLE`. *(**ADR-0003**)*
- [ ] ⚠️ **`ID` adalah PRIMARY KEY di kelima tabel.** Skrip migrasi **memverifikasi keberadaannya**
      dan **menambahkannya bila belum ada** — tidak berasumsi. *(AC 51 spec;
      `[data DBA]` + `[keputusan work owner]`)*
- [ ] ⚠️ **Keempat FK terpasang** dengan mode **`ON DELETE CASCADE`**. Skrip **memverifikasi mode
      `ON DELETE`-nya**, bukan hanya keberadaan FK-nya — mode yang keliru (`RESTRICT`/`NO ACTION`)
      **bertentangan** dengan kaskade tiket 09. *(AC 52 spec)*
- [ ] ⚠️ `RIRATE` **tetap `VARCHAR2`** — **tidak** diubah menjadi numerik. *(AC 26 spec;
      `[data DBA]`)*
- [ ] Kelima **sequence** pindah dengan **nilai berjalan yang benar**, sehingga identitas baru
      **tidak pernah bertabrakan** dengan yang lama. Diuji dengan membuat satu baris di tiap tabel
      sesudah migrasi.
- [ ] Format identitas `'1' || lpad(seq, 6, '0')` **tetap berlaku** sesudah migrasi.
- [ ] Nilai uang dan share pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan nilai
      lama dan baru **secara tepat**, bukan dengan toleransi.
- [ ] Nilai `DATE` pindah **tanpa pergeseran zona waktu**.
- [ ] Wajib-isi **tidak** ditambahkan sebagai `NOT NULL` tanpa keputusan terpisah — `[data DBA]`
      basis data lama **nol `NOT NULL`**, dan data lama mungkin memuat kolom kosong yang akan
      menolak migrasi. *(AC 53 spec — penegakan tetap di Go)*
- [ ] Index pendukung ada untuk kolom yang dipakai kueri hilir — khususnya kolom FK, dan
      `TREATYCONTRACTID` yang dipakai penghitungan total share.
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.
- [ ] Skema uji yang dipakai seluruh tiket lain dibangun **dari DDL yang sama** dengan produksi —
      bukan dari tiruan yang ditulis terpisah. *(AC 54 spec)*

## Blocker

**Tidak ada.** ✅ **OQ-001 ditutup** `[data DBA]` 2026-09-15 — tipe, nullability, PK, FK, dan
sequence seluruhnya diketahui.

⚠️ **Verifikasi keadaan nyata, jangan berasumsi** `[keputusan work owner]`. Keadaan yang dinyatakan:
**kelima tabel sudah punya PK `ID`** (business dari DBA; empat lainnya ditambahkan work owner) dan
**keempat FK sudah terpasang dengan mode `ON DELETE CASCADE`**. Meski demikian, skrip migrasi
**memeriksa keadaan sebenarnya lebih dulu** lalu menambahkan yang belum ada — termasuk
**memverifikasi mode `ON DELETE`**, karena mode `RESTRICT`/`NO ACTION` akan **bertentangan langsung**
dengan kaskade tiket 09. Ini sudah tercermin di AC di atas.

⚠️ **Catatan sumber yang sudah usang.** `.scratch/master-contract-retro-life/ddl-tables-from-dba.md`
pada bagian akhirnya masih memuat dua pernyataan yang **tidak lagi berlaku**: (1) FK "masih nol", dan
(2) Q5 sebagai "tolak hapus induk berpunya-anak" dengan `ON DELETE CASCADE` disebut *bertentangan*.
Keduanya **sudah direvisi** — lihat `grilling-ronde-1-jawaban.md` §Q5 (kaskade + popup konfirmasi)
dan §Status OQ. Bila membaca berkas DDL itu, **pakai keadaan final di sini**.

## Catatan

⚠️ **`USERID` tidak seragam** `[data DBA]`: `VARCHAR2(100)` di tiga tabel, `VARCHAR2(1000)` di dua
lainnya. Diseragamkan atau dibawa apa adanya — **putuskan di tiket ini**, jangan diam-diam
mengubahnya.

⚠️ **Anti-dobel logis belum diputuskan.** `[data DBA]` Basis data **nol `UNIQUE`** selain PK `ID`.
Bila anti-dobel logis dikehendaki (mis. satu reinsurer hanya sekali per kontrak), constraint-nya
ditambahkan **di sini** — tetapi keputusannya milik tiket 05/06.

## Seam & perintah verifikasi

Migrasi diuji terhadap **skema uji Oracle nyata** — bukan mock.

```
go test ./internal/...
```

---

## Ralat bertanggal 30-09-2026 — sesi implementasi (paket 0)

> Sumber: `RALAT-DEV-30-09-2026.md` (K1–K8 katalog DEV, R1–R12 pembacaan ulang XML) dan `PARITAS-LAYAR-DAN-AKSI.md`. Kalimat di atas **tidak dihapus**; yang berlaku adalah ralat ini.

| Kalimat lama | Ralat |
| --- | --- |
| *"Status: ready-for-agent — ✅ OQ-001 ditutup"*; *"DDL kelima tabel + PK + FK + sequence"*; *"Keadaan yang dinyatakan: kelima tabel sudah punya PK `ID` … keempat FK sudah terpasang dengan mode `ON DELETE CASCADE`"* | ⛔ **tiket DICABUT** (K1): kelima tabel di DEV nol constraint P/R/U. **Nol migrasi, nol DDL, nol tabel baru** — preseden tco4 Treaty Contract Out `[keputusan asisten dari preseden tco4; veto work owner]`. Rentang 100–139 tetap kosong; penjaga modul menolak berkas migrasi di rentang itu. Keunikan `ID` dan kaskade di Go (K2, K3); pertanyaan asal DDL → OQ-MCRL-02 |

**Status:** `wontfix` — dicabut 30-09-2026.
