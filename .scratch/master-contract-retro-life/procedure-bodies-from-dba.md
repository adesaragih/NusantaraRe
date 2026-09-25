# Body 5 Stored Procedure Penulis — dari DBA (menutup OQ-002)

Tanggal: 2026-09-15
Sumber: **DBA / work owner** (dikirim langsung; bukan dari korpus Pega).
Status: `[data DBA]` — menutup **OQ-002** untuk Master Contract Retro Life.

> Konteks: 5 procedure ini adalah **seluruh jalur tulis** modul Master Contract Retro Life.
> Grilling Bagian A membuktikan nol logika perhitungan di sisi Pega; body ini mengonfirmasi
> nol logika bisnis di sisi DB juga — hanya **upsert mentah + audit**.

---

## Temuan lintas kelima procedure (pola IDENTIK)

| Aspek | Perilaku terbukti | Implikasi untuk Go |
| --- | --- | --- |
| **Mode simpan** | **UPSERT dikunci `ID`**: `SELECT COUNT(1) … WHERE ID=p_ID` → ada = UPDATE, tidak ada = INSERT | Repository "save" = upsert; bukan insert murni |
| **Generasi ID** | Saat INSERT, ID dibuat DB: `'1' \|\| lpad(seq.nextval,6,'0')` (mis. `1000001`). `p_ID` **diabaikan** saat INSERT, dipakai hanya untuk mencari saat UPDATE | ID diciptakan DB; aplikasi kirim ID kosong untuk baris baru. Sequence per tabel. Konsisten ADR-0006 (penomoran lewat procedure/DB) |
| **`HASIL1` / `o_message`** | **Pesan galat teks**, bukan kode angka. **Kosong/NULL = sukses**; berisi teks = gagal. ⚠️ Berisi **HTML** (`<span style="color:red">…</span>`) + `SQLERRM` | Q11: periksa `o_message` — **non-empty = gagal**. ⚠️ **Strip/abaikan HTML** dari Pega; bangun pesan bersih di Go |
| **`TGLUPDATE`** | Selalu `SYSDATE` — parameter `p_TGLUPDATE` **diabaikan** | Jangan andalkan cap waktu aplikasi; DB yang menetapkan |
| **`USERID`** | Disimpan apa adanya dari `p_USERID` (audit) | Kirim identitas pengguna |
| **Transaksi** | **`COMMIT` di dalam tiap procedure**; `ROLLBACK` pada exception terluar | Tiap simpan = transaksi mandiri. **Tidak ada** transaksi lintas-baris di sisi DB |
| **Validasi bisnis** | **NOL** — tidak ada cek 100%, tidak ada cek anak/kaskade, tidak ada unique selain PK `ID` | Semua aturan (total share display, hapus berpenjaga, gerbang tahun) **wajib di lapisan Go** — DB tidak menegakkannya |
| **Tipe parameter uang** | ⚠️ Kolom uang & share diterima procedure sebagai **`VARCHAR2`** (`p_IDR`, `p_USD`, `p_B_IDR`, `p_PCTSHARE`, … semua VARCHAR2) | Perlu DDL untuk tahu tipe **kolom** sebenarnya (NUMBER vs VARCHAR2). Menentukan ADR-0003 |

⚠️ **Konsekuensi penting — `USD` boleh kosong:** `INSERTTREATYCONTRACT_LIFE` meng-INSERT `p_USD`
tanpa cek NULL/kosong. Jadi dari sisi procedure, `USD` **tidak wajib** — konsisten dengan temuan
Pega (`SaveTreatyLimit_Act`: `USD` tidak wajib). Tinggal konfirmasi nullability kolom via DDL.

⚠️ **`o_message` mengandung HTML** — ini "pesan galat bergaya UI Pega". Di Go, `o_message` non-empty
= gagal; teks HTML **jangan** diteruskan mentah ke API/UI baru — ekstrak maknanya, sajikan bersih.

⚠️ **Anti-dobel:** upsert dikunci `ID` **tidak** mencegah duplikat logis (mis. dua reinsurer sama di
satu kontrak) — hanya mencegah dua baris ber-`ID` sama. Bila anti-dobel logis diinginkan, itu aturan
Go + (idealnya) unique constraint DB. → catat sebagai keputusan tiket.

---

## Ringkasan per procedure

Kelima berpola sama; perbedaan hanya kolom.

### 1. `INSERTTREATYYEAR_LIFE`
Upsert `TREATYYEAR_LIFE` by `ID`. INSERT-ID = `'1'+lpad(TREATYYEAR_LIFE_seq,6)`.
Kolom: `ID, TREATYYEAR, UNDERWRITINGYEAR, USERID, TGLUPDATE(SYSDATE), STARTDATE, ENDDATE`.

### 2. `INSERTTREATYCONTRACT_LIFE`
Upsert `TREATYCONTRACT_LIFE` by `ID`. INSERT-ID = `'1'+lpad(TREATYCONTRACT_LIFE_seq,6)`.
Kolom: `ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME, TREATYSTARTDATE, TREATYENDDATE, USERID,
TGLUPDATE(SYSDATE), IDR, USD, B_IDR, B_USD, IDR_SELISIH, USD_SELISIH`.
⚠️ `IDR_SELISIH`/`USD_SELISIH` **ditulis apa adanya dari parameter** — DB **tidak** menghitungnya.
Penyimpangan sadar 2 (selisih dihitung di Go) tetap berlaku: Go yang menghitung sebelum mengirim,
**atau** Go tidak mengirim kolom itu. → keputusan tiket.

### 3. `INSERTREINSURER_LIFE`
Upsert `TREATYREINSURER_LIFE` by `ID`. INSERT-ID = `'1'+lpad(TREATYREINSURER_LIFE_SEQ,6)`.
Kolom: `ID, TREATYYEARID, TREATYCONTRACTID, REINSTYPEID, REINSTYPENAME, REINSURERID, REINSURERNAME,
PCTSHARE, COMMISION, OVR_COMM, USERID, TGLUPDATE(SYSDATE)`.

### 4. `INSERTSECURITYREINSURER_LIFE`
Upsert `TREATYSECURITYREINSURER_LIFE` by `ID`. INSERT-ID = `'1'+lpad(TREATYSECURITYREINSURER_LIFE_SEQ,6)`.
Kolom: `ID, TREATYYEARID, TREATYCONTRACTID, TREATYREINSURERID, REINSURERID, REINSURERNAME, PCTSHARE,
USERID, TGLUPDATE(SYSDATE)`.
⚠️ `TREATYREINSURERID` = FK logis ke reinsurer induk (dasar retrosesi-atas-retrosesi, Q3) —
disimpan mentah, **tidak** ditegakkan procedure.

### 5. `INSERTBUSINESS_LIFE`
Upsert `TREATYBUSINESS_LIFE` by `ID`. INSERT-ID = `'1'+lpad(TREATYBUSINESS_LIFE_SEQ,6)`.
Kolom: `ID, TREATYYEARID, TREATYYEAR, TREATYCONTRACTID, REINSTYPEID, REINSTYPENAME, BIZCODE, BIZNAME,
RIRATEID, RIRATE, USERID, TGLUPDATE(SYSDATE)`.

---

## Dampak ke keputusan spec

| Keputusan | Status setelah body |
| --- | --- |
| Q9 panggil apa adanya | ✅ dikuatkan — procedure aman dipanggil; upsert deterministik |
| Q11 `HASIL1` diperiksa | ✅ diperjelas — **non-empty `o_message` = gagal**; strip HTML |
| Q1 normalisasi | ⚠️ **tetap Go-side** — procedure TETAP menulis REINSTYPEID/TREATYYEAR ke tabel anak; normalisasi = keputusan skema baru kita, bukan tiru procedure |
| Q2 selisih dihitung | ⚠️ **tetap Go-side** — procedure menulis SELISIH dari parameter; Go yang menghitung/menghilangkan |
| Q4 total 100% | ✅ dikuatkan — DB tak menegakkan; display-only di Go benar |
| Q5 hapus berpenjaga | ⚠️ **wajib Go-side** — procedure/DB tak punya kaskade; penjaga di Go |
| Q6 tahun abadi | ✅ konsisten — tak ada procedure hapus tahun |
| ADR-0006 penomoran | ✅ dikuatkan — ID dibuat DB via sequence |

## OQ setelah ini

- **OQ-002** → ✅ **DITUTUP** (body diketahui; aturan simpan = upsert + audit, nol logika bisnis).
- **OQ-001** → **masih terbuka sebagian**. Yang TERJAWAB oleh body: format ID (`'1'+lpad(seq,6)`),
  `USD` boleh kosong (sisi procedure), uang lewat sebagai VARCHAR2 di parameter. Yang MASIH perlu
  **DDL 4 tabel**: tipe kolom sebenarnya (NUMBER(p,s) vs VARCHAR2 — presisi uang ADR-0003),
  NOT NULL/nullable tiap kolom, ada/tidak PK/FK/unique/index. Sequence terkonfirmasi ada
  (dirujuk procedure): `TREATYYEAR_LIFE_seq`, `TREATYCONTRACT_LIFE_seq`, `TREATYREINSURER_LIFE_SEQ`,
  `TREATYSECURITYREINSURER_LIFE_SEQ`, `TREATYBUSINESS_LIFE_SEQ`.
- **OQ kecil** (nama tabel fisik REINSURANCETYPE, isi `.Code`/`.Type`) → masih terbuka, tak memblokir.
