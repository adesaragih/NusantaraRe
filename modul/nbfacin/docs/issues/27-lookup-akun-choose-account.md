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

**Status:** ready-for-human — dibangun 02-10-2026; diubah 02-10-2026 sore (butir 73) dan 03-10-2026 (butir 75: kembali
tidak peka huruf); uji tanpa Oracle hijau; uji terhadap Oracle sungguhan **tidak** dijalankan

## Kontrak

`GET /api/nbfacin/account?cari=<teks>&halaman=<n>` — `halaman` mulai 1 (kosong = 1). `cari` "mengandung", **tidak peka
huruf besar-kecil** (`uji` = `Uji` = `UJI`, butir 75); 15 baris per halaman (`"ukuran":15`, butir 73.2 tetap).

| Kode | Kapan | Badan |
| --- | --- | --- |
| 200 | berhasil | `{"baris":[{"id","insuredId","insuredName","groupBusinessId","groupBusiness"}],"total":<int>,"halaman":<int>,"ukuran":<int>}` — teks apa adanya, kolom NULL → `""`; `baris` selalu larik (bisa kosong) |
| 400 | `halaman` bukan bilangan bulat ≥ 1, atau `cari` > 255 karakter | `{"galat": "..."}` (bentuk `galat.Tulis`; 400/503/500 sama) |
| 503 | aplikasi tanpa basis data | alasan: tabel akun tidak terbaca |
| 500 | galat Oracle / program | pesan tetap "galat server"; rinciannya hanya di log, tanpa nilai baris |

## Yang dibangun

- [x] `repository/akun.go` — `SELECT ID, GROUPBUSINESSID, GROUPBUSINESS, INSUREDID, INSUREDNAME` saja, tabel lewat
      `db.Qualify`; parameter terikat; pencarian `UPPER(kolom) LIKE :n ESCAPE '\'` dengan pola huruf besar (tidak peka huruf,
      butir 75) dan `\` `%` `_` diloloskan (escape sepola `PolaCari` modul master lain); cacah total dibaca terpisah (`COUNT(*)`, pola inbox claimlife)
- [x] `services` — `CariAkun` (validasi halaman/panjang cari, offset, 503 tanpa DB)
- [x] `handlers` — rute, bentuk jawaban di atas
- [x] `MODUL.md` *Tabel warisan* + `docs/STRUKTUR-TABEL-NB-FACIN.md` — `T_M_ACCOUNT` dibaca, tidak dibuat
- [ ] Uji terhadap Oracle sungguhan — ⛔ tidak dijalankan (aturan kerja: tidak menjalankan apa pun ke Oracle)

## Keputusan agent — A69, A72, A73 menunggu konfirmasi; A70, A71 sudah diputus work owner

| # | Keputusan | Dasar |
| --- | --- | --- |
| A69 | Kolom yang dicari: **tiga kolom yang tampil** — `INSUREDID`, `INSUREDNAME`, `GROUPBUSINESS` (OR) | Work owner hanya menjawab "mengandung", tidak menyebut kolomnya |
| ~~A70~~ | ~~Tidak peka huruf besar-kecil (`UPPER` di kedua sisi)~~ → ~~DIUBAH work owner (butir 73.1): PEKA huruf — `LIKE` tanpa `UPPER`~~ → **DIUBAH LAGI work owner (butir 75): TIDAK peka huruf** — `UPPER(kolom) LIKE :n`, pola dibesarkan di Go (`strings.ToUpper`), bukan `UPPER(:n)` di SQL; bentuk yang dibangun semula untuk A70 dan diverifikasi sesi 0f | butir 73.1, kutipan **"harus peka besar kecil dong, 15 baris per halaman"** — **dibatalkan** butir 75, kutipan **"pada saat search Group Business, itukan ada isian untuk search, itu buatin tanpa liat huruf besar atau kecil"** (keduanya diteruskan sesi 0f) |
| ~~A71~~ | ~~Ukuran halaman 20~~ → **DIUBAH work owner (butir 73.2): 15** — tetap berlaku sesudah butir 75 | kutipan butir 73 |
| A72 | Urutan **`INSUREDID`, lalu `ID`** (deterministik; `INSUREDID` tidak dijamin unik) | urutan di gambar tidak terlihat jelas |
| A73 | `cari` > 255 karakter → 400 | kolom terpanjang 255 karakter; masukan lebih panjang tidak mungkin cocok |

⚠️ **Batas pengujian:** sifat "tidak peka huruf" dan escape di sisi Oracle hanya diuji lewat **teks SQL** (tepat tiga
`UPPER(kolom)`) dan pola bind (`TestSQLAkun`, `TestPolaCari`: `uji`, `Uji`, `UJI` berpola sama), bukan dengan eksekusi
Oracle. `[dugaan]` huruf besar pola dibuat `strings.ToUpper` Go sedangkan kolom `UPPER` Oracle — untuk huruf non-ASCII
keduanya dapat berbeda; belum diuji ke Oracle. ~~`[dugaan]` `LIKE` peka huruf bergantung `NLS_COMP`~~ — catatan masa butir 73.1, tidak
relevan lagi sejak `UPPER` dipakai (teksnya utuh di register butir 73).

## Di luar tiket ini (menunggu work owner)

**Class Of Business** — ~~DDL belum ada~~ → dibangun di **tiket 28** (`GET /api/nbfacin/class-of-business`) sesudah DDL
`BUSINESS.txt` ditambahkan work owner 02-10-2026.

## Comments

*Tanpa nama orang, tanpa data pelanggan — uji memakai data sintetis berawalan `UJI-`.*
