# 27: Lookup akun untuk popup ChooseAccount — `GET /api/nbfacin/account` (baca saja)

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah: permintaan sesi `nusantarare-0f`
> 02-10-2026 untuk tiket 26 C-8 (form Opportunity, popup ChooseAccount; frontend di sesi 0f, backend `nbfacin` di sesi
> ini), berdasar jawaban work owner yang dikutip sesi itu.

**What to build:** pembaca baca-saja `POOLDATA.T_M_ACCOUNT` dan satu endpoint pencarian berhalaman yang dipakai popup
ChooseAccount.

**Asal:**
- DDL READ-ONLY `D:\migrasi\RNM\DDL\T_M_ACCOUNT.txt` (ditambahkan work owner 02-10-2026) `[terverifikasi]`: `ID
  VARCHAR2(255 CHAR) NOT NULL`, `GROUPBUSINESSID VARCHAR2(32 CHAR)`, `GROUPBUSINESS VARCHAR2(64 CHAR)`, `INSUREDID
  VARCHAR2(255 CHAR) NOT NULL`, `INSUREDNAME VARCHAR2(64 CHAR)`; tanpa PK/indeks di berkas itu.
- Jawaban work owner 02-10-2026 (dikutip sesi 0f): kolom layar Insured ID = `INSUREDID`, Insured Name = `INSUREDNAME`,
  Group Business = `GROUPBUSINESS`; pencarian "mengandung".
- `[terverifikasi]` sesi 0f: rule Pega `ChooseAccount` **tidak ada** di korpus (dua folder) — perilaku hanya dari tangkapan
  layar + jawaban work owner. Bukan port rule; tidak ada rumus atau perilaku Pega yang ditiru selain itu.

**Blocked by:** —

**Status:** ready-for-human — dibangun 02-10-2026; diubah 02-10-2026 sore sesuai keputusan work owner atas A70/A71
(butir 73); uji tanpa Oracle hijau; uji terhadap Oracle sungguhan **tidak** dijalankan

## Kontrak

`GET /api/nbfacin/account?cari=<teks>&halaman=<n>` — `halaman` mulai 1 (kosong = 1). `cari` "mengandung", **peka huruf
besar-kecil** (`uji` ≠ `UJI`); 15 baris per halaman (`"ukuran":15`).

| Kode | Kapan | Badan |
| --- | --- | --- |
| 200 | berhasil | `{"baris":[{"id","insuredId","insuredName","groupBusinessId","groupBusiness"}],"total":<int>,"halaman":<int>,"ukuran":<int>}` — teks apa adanya, kolom NULL → `""`; `baris` selalu larik (bisa kosong) |
| 400 | `halaman` bukan bilangan bulat ≥ 1, atau `cari` > 255 karakter | `{"galat": "..."}` (bentuk `galat.Tulis`; 400/503/500 sama) |
| 503 | aplikasi tanpa basis data | alasan: tabel akun tidak terbaca |
| 500 | galat Oracle / program | pesan tetap "galat server"; rinciannya hanya di log, tanpa nilai baris |

## Yang dibangun

- [x] `repository/akun.go` — `SELECT ID, GROUPBUSINESSID, GROUPBUSINESS, INSUREDID, INSUREDNAME` saja, tabel lewat
      `db.Qualify`; parameter terikat; pencarian `kolom LIKE :n ESCAPE '\'` **tanpa** `UPPER` (peka huruf), pola huruf apa adanya
      dengan `\` `%` `_` diloloskan (escape sepola `PolaCari` modul master lain); cacah total dibaca terpisah (`COUNT(*)`, pola inbox claimlife)
- [x] `services` — `CariAkun` (validasi halaman/panjang cari, offset, 503 tanpa DB)
- [x] `handlers` — rute, bentuk jawaban di atas
- [x] `MODUL.md` *Tabel warisan* + `docs/STRUKTUR-TABEL-NB-FACIN.md` — `T_M_ACCOUNT` dibaca, tidak dibuat
- [ ] Uji terhadap Oracle sungguhan — ⛔ tidak dijalankan (aturan kerja: tidak menjalankan apa pun ke Oracle)

## Keputusan agent — A69, A72, A73 menunggu konfirmasi; A70, A71 sudah diputus work owner

| # | Keputusan | Dasar |
| --- | --- | --- |
| A69 | Kolom yang dicari: **tiga kolom yang tampil** — `INSUREDID`, `INSUREDNAME`, `GROUPBUSINESS` (OR) | Work owner hanya menjawab "mengandung", tidak menyebut kolomnya |
| ~~A70~~ | ~~Tidak peka huruf besar-kecil (`UPPER` di kedua sisi)~~ → **DIUBAH work owner (butir 73.1): PEKA huruf** — `LIKE` tanpa `UPPER`, kolom dan pola apa adanya | kutipan work owner (diteruskan sesi 0f): **"harus peka besar kecil dong, 15 baris per halaman"** |
| ~~A71~~ | ~~Ukuran halaman 20~~ → **DIUBAH work owner (butir 73.2): 15** | kutipan yang sama |
| A72 | Urutan **`INSUREDID`, lalu `ID`** (deterministik; `INSUREDID` tidak dijamin unik) | urutan di gambar tidak terlihat jelas |
| A73 | `cari` > 255 karakter → 400 | kolom terpanjang 255 karakter; masukan lebih panjang tidak mungkin cocok |

⚠️ **Batas pengujian:** sifat "peka huruf" dan escape di sisi Oracle hanya diuji lewat **teks SQL** (tanpa `UPPER`/`LOWER`)
dan pola bind (`TestSQLAkun`, `TestPolaCari`: `uji` dan `UJI` berpola berbeda), bukan dengan eksekusi Oracle. `[dugaan]`
`LIKE` Oracle peka huruf bila `NLS_COMP` = `BINARY` (bawaan); bila instance memakai `NLS_COMP=LINGUISTIC` dengan `NLS_SORT`
berakhiran `_CI`, hasilnya bisa tidak peka huruf — `belum terverifikasi` (tanya DBA).

## Di luar tiket ini (menunggu work owner)

**Class Of Business** — ~~DDL belum ada~~ → dibangun di **tiket 28** (`GET /api/nbfacin/class-of-business`) sesudah DDL
`BUSINESS.txt` ditambahkan work owner 02-10-2026.

## Comments

*Tanpa nama orang, tanpa data pelanggan — uji memakai data sintetis berawalan `UJI-`.*
