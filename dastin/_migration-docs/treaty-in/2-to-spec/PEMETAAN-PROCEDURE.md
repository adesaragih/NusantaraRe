# Pemetaan lima stored procedure — ke mana tiap potong logikanya pergi

**Tanggal:** 24 September 2026 · **Langkah 2 sesi to-spec** · **Pasangannya:** `docs/adr/0056-tanpa-stored-procedure.md`
**Sumber:** `PROCEDURE/*.txt` — lima berkas, 674 baris, dibaca seluruhnya

> **Tiga golongan, tidak ada yang keempat.**
> **PINDAH KE GOLANG** — aturan bisnis · **TURUN JADI CONSTRAINT** — penjagaan integritas yang
> lebih tepat di Oracle · **TIDAK DIBAWA** — cacat atau perancah, dan cacatnya disebut.

---

## 1. Rekap

| Procedure | Baris | Pindah | Constraint | Tidak dibawa |
|---|---:|---:|---:|---:|
| `PEGA_TREATY_IN` | 154 | 3 | 2 | **4** |
| `PEGA_M_TREATY_IN_EDM` | 150 | 2 | 2 | **4** |
| `PEGA_M_TREATY_IN_DETAIL` | 161 | 2 | 1 | **4** |
| `PEGA_M_TREATY_IN_DETAIL_EDM` | 152 | 2 | 1 | **4** |
| `GET_TOKEN_STORAGE` | 57 | 1 | 1 | **3** |

**Empat dari lima cacatnya sebangun**, dan itu bukan kebetulan: keempatnya disalin dari satu
cetakan yang sama. Cacatnya ikut tersalin.

---

## 2. `PEGA_TREATY_IN` — penyimpan kontrak

| # | Potong logika | Golongan | Keterangan |
|---|---|---|---|
| 1 | `IF IDPega = 'UnknownId'` memilih antara terbit-baru dan perbarui | **PINDAH KE GOLANG** | ini keputusan "kontrak ini sudah ada atau belum" — aturan bisnis, bukan penyimpanan. Di sistem baru ia keputusan lapisan aplikasi, dan tidak lagi ditandai lewat **nilai ajaib berupa teks** |
| 2 | `SELECT ID FROM M_SITE_DATABASE WHERE CURRENT_SITE='1'` | **PINDAH KE GOLANG** | penentuan situs penerbit; ia bagian aturan penomoran, bukan integritas |
| 3 | `id_site \|\| lpad(M_TREATY_IN_SEQ.nextval,6,'0')` | **PINDAH KE GOLANG** | bentuk pengenal. Catatan: `lpad(...,6,'0')` inilah yang membuat skema pengenal pecah saat urutannya melewati enam digit |
| 4 | `INSERT INTO M_TREATY_IN` dan `INSERT INTO TREATY_IN` dalam satu transaksi | **TURUN JADI CONSTRAINT** | yang dijaga: baris datar tidak boleh ada tanpa dokumennya. Di skema baru tidak ada dua bentuk, jadi penjagaannya menjadi **kunci asing biasa** |
| 5 | `COMMIT` di dalam procedure | **TURUN JADI CONSTRAINT** — sebagai **larangan** | batas transaksi dimiliki pemanggil, bukan penyimpan. Ini yang membuat cacat #8 tidak dapat dibatalkan |
| 6 | **`replace(DataPega,'UnknownId',id_TreatyIn_ins)`** | **TIDAK DIBAWA — cacat** | penyulihan **teks buta atas seluruh dokumen JSON**. Ia mengganti **setiap** kemunculan kata `UnknownId` di mana pun di dalam dokumen, termasuk di dalam nilai teks yang kebetulan memuatnya. Tidak ada yang membatasinya pada ruas pengenal |
| 7 | **`ErrMsg := 'Insert Rate Error : '`** pada penyimpan **treaty**, bukan *rate* | **TIDAK DIBAWA — cacat** | pesan galat menyebut modul yang salah di **keempat** procedure. Orang yang mencarinya di log akan mencari di modul yang keliru |
| 8 | **Dipanggil dari layar addendum, menghasilkan `UPDATE … WHERE ID = 'XXXXXXX/Rnn'`** | **TIDAK DIBAWA — cacat** | pengenal versi tidak pernah ada di `TREATY_IN`. `UPDATE` mengenai **nol baris**, Oracle **tidak** menimbulkan galat, `COMMIT` berhasil, `StsSave := 1`, dan pesannya berbunyi *"Data Sudah Disimpan Dengan ID : …"*. **Kegagalan senyap yang dilaporkan sebagai keberhasilan** |
| 9 | `EXCEPTION WHEN OTHERS THEN … RETURN` di tiga tingkat | **TIDAK DIBAWA — cacat** | `WHEN OTHERS` tanpa `RAISE` menelan **setiap** galat, termasuk yang tidak diantisipasi. Cacat #8 tetap senyap walaupun bukan karena penangkap ini — tetapi penangkap ini memastikan tidak ada cacat lain yang pernah terdengar |

---

## 3. `PEGA_M_TREATY_IN_EDM` — penyimpan addendum

Bentuknya kembar dengan §2. Yang **berbeda** dan karena itu dicatat:

| # | Potong logika | Golongan | Keterangan |
|---|---|---|---|
| 1 | **`IF STSINPUT = '0'`** memilih sisip atau perbarui | **PINDAH KE GOLANG** | penanda pemilihnya **berbeda dari saudaranya**: kontrak memakai `IDPega = 'UnknownId'`, addendum memakai parameter tersendiri. Dua cara menyatakan hal yang sama, di dua procedure yang disalin dari satu cetakan |
| 2 | `EDMDATE = SYSDATE` **pada sisip DAN pada perbarui** | **TIDAK DIBAWA — cacat** | kolom yang namanya berbunyi seperti *tanggal addendum* sebenarnya **tanggal sentuh terakhir**. Setiap penyimpanan ulang menggesernya, sehingga **tanggal pembuatan addendum tidak dapat dipulihkan dari kolom ini**. Di model baru, dibuat dan diubah adalah dua fakta (ADR-0006) |
| 3 | `OLDID` diperbarui **setiap kali** | **TURUN JADI CONSTRAINT** | penunjuk versi pendahulu seharusnya **beku** sesudah versi lahir. Di skema baru ia `ID_VERSI_KONTRAK_DASAR`, dan keberubahannya dilarang |
| 4 | tidak ada `M_SITE_DATABASE`, tidak ada urutan | **PINDAH KE GOLANG** | pengenal addendum **dirakit di Pega**, bukan di basis data — itulah sebabnya bentuk `/Rnn` sampai ke sini sebagai teks |
| 5–8 | `replace` buta · pesan *"Rate"* · `COMMIT` di dalam · `WHEN OTHERS` menelan | **TIDAK DIBAWA — cacat** | identik §2 |

---

## 4. `PEGA_M_TREATY_IN_DETAIL` dan `…_DETAIL_EDM` — proyeksi datar

### 4.1 Cacat yang sama beratnya dengan §2 butir 8, dan belum pernah tercatat

> **Kedua procedure ini TIDAK PUNYA CABANG `ELSE`.**

```sql
IF IDPega = 'UnknownId' THEN
   … sisip …
END IF;

StsSave := 1;
IDPegaOut := id_TreatyIn_out;
```

Bila `IDPega` **bukan** `'UnknownId'` — yaitu setiap kali baris rincinya **sudah ada dan hendak
diperbarui** — maka:

* **tidak satu pun pernyataan dijalankan**;
* `id_TreatyIn_out` **tidak pernah diberi nilai**, sehingga `IDPegaOut` keluar **NULL**;
* dan `StsSave := 1` **tetap dijalankan** — **berhasil**.

**Golongan: TIDAK DIBAWA — cacat.** Ia sekeluarga dengan §2 butir 8 dan dengan langkah Pega
berjudul *"if empty exit activity"* yang tidak memuat kode keluar: **pekerjaan yang tidak dilakukan,
dilaporkan sebagai selesai.**

Akibat yang dapat diukur: **`TREATYINDETAIL` tidak pernah diperbarui, hanya ditambah.** Setiap
penyimpanan ulang kontrak menyisipkan baris rinci baru atau tidak melakukan apa pun — dan mana di
antara keduanya bergantung pada apa yang Pega kirim sebagai `IDPega`. **Uji datanya**: apakah ada
`TREATYID` dengan lebih dari satu baris rinci untuk kombinasi layer/kelompok yang sama —
diusulkan sebagai **Uji AN**.

### 4.2 Dan proyeksi datar addendum TERTINGGAL TUJUH KOLOM

Kolom `INSERT` kedua procedure dibandingkan satu per satu:

| | Jumlah kolom |
|---|---:|
| `TREATYINDETAIL` | **62** |
| `TREATYINDETAILEDM` | **55** |

`TREATYINDETAILEDM` adalah **himpunan bagian murni** — tidak ada satu pun kolom yang hanya ada
padanya. Yang **hilang pada sisi addendum**:

```
DEDUCTIBLE   DEDUCTIBLE2   ADJ_RATE   PREMIUM_EARNED   MDP_PCT   ROL_PCT   MDP
```

**Ketujuhnya besaran bisnis non-proporsional**, dan ketujuhnya justru yang paling mungkin berubah
lewat addendum. Siapa pun yang membaca riwayat kontrak dari kedua tabel datar ini akan melihat
deductible dan MDP **hanya pada kontrak asalnya**, tidak pernah pada perubahannya.

**Golongan: TIDAK DIBAWA — cacat.** Di skema baru tidak ada dua proyeksi; `VERSI_KONTRAK` dan
anak-anaknya membawa besaran yang sama untuk versi mana pun, sehingga ketimpangan ini **lenyap
karena bentuknya, bukan karena ditambal**.

> **Ia juga memperkuat `ADR-0051`.** Kalau ada konsumen hilir yang membaca `TREATYINDETAILEDM`,
> konsumen itu selama ini **menerima data yang kurang tujuh kolom** — dan tidak ada yang tahu,
> karena tabelnya terlihat lengkap.

### 4.3 Sisanya

| # | Potong logika | Golongan |
|---|---|---|
| 1 | `SELECT COMMENCEMENT, TERMINATION FROM TREATY_IN WHERE ID = P_TREATYID` lalu `to_date(...,'YYYYMMDD')` | **PINDAH KE GOLANG** — pengambilan periode induk; dan di skema baru ia tidak perlu diambil, karena rincinya **berinduk** pada versinya |
| 2 | `WHEN NO_DATA_FOUND THEN v_COMMENCEMENT := NULL` | **TIDAK DIBAWA — cacat** | induk yang tidak ditemukan menjadi **tanggal kosong**, bukan penolakan. Baris rinci yatim tetap tersimpan. Di skema baru: **kunci asing**, dan ia menolak |
| 3 | `INSERT` dua tabel dalam satu transaksi | **TURUN JADI CONSTRAINT** |
| 4 | **`ErrMsg` dideklarasikan `OUT` dan tidak pernah diberi nilai sama sekali** | **TIDAK DIBAWA — cacat** | pemanggil menerima pesan kosong pada kegagalan maupun keberhasilan. Berbeda dari §2 butir 7, yang setidaknya **salah**; yang ini **tidak ada** |

---

## 5. `GET_TOKEN_STORAGE` — bukan Treaty In

| # | Potong logika | Golongan | Keterangan |
|---|---|---|---|
| 1 | `IF VAPPNAME IS NULL` menolak | **TURUN JADI CONSTRAINT** | `NOT NULL` biasa |
| 2 | ambil token belum kedaluwarsa, kalau tidak ada terbitkan baru | **PINDAH KE GOLANG** | kebijakan siklus hidup token |
| 3 | **`INPUTDATE > SYSDATE` sebagai uji "belum kedaluwarsa"** | **TIDAK DIBAWA — cacat penamaan** | kolom bernama *tanggal masukan* menyimpan **tanggal kedaluwarsa** (`v_expdate := SYSDATE + 1 menit`). Nama yang berarti lain daripada isinya — persis larangan `CONTEXT.md` §2.5 |
| 4 | **`STANDARD_HASH('ASMAPP' \|\| timestamp, 'MD5')` sebagai token akses** | **TIDAK DIBAWA — cacat keamanan** | token **dapat ditebak sepenuhnya**: awalannya tetap, sisanya waktu. Siapa pun yang tahu bentuknya dapat membangkitkan token yang sah untuk detik mana pun. **Tidak ada keacakan.** Ini **bukan** temuan Treaty In dan **tidak** ditambal di sini — ia dilaporkan, dan pemiliknya di §7 |
| 5 | masa berlaku **satu menit** | **TIDAK DIBAWA** | bukan keputusan yang dinyatakan di mana pun; ia tetapan di kode |

> **`GET_TOKEN_STORAGE` di luar lingkup Treaty In.** Ia ikut dibedah karena ada di folder yang
> sama, dan **butir 4 tidak boleh hanya duduk di berkas ini.**

---

## 6. Kolom yang dipanen — dan apa yang panenan itu TIDAK dapat katakan

`alat/panen-kolom-dari-penulis.py` dijalankan. **Batas kemampuannya dibawa apa adanya: ia dapat
menyatakan sebuah kolom ADA — karena ada yang menulisinya — dan TIDAK PERNAH bahwa sebuah kolom
TIDAK ADA.**

| Tabel | Kolom terpanen | Dari |
|---|---:|---|
| **`TREATY_IN`** | **20** | `PEGA_TREATY_IN` INSERT + UPDATE |
| **`TREATY_IN_EDM`** | **24** | `PEGA_M_TREATY_IN_EDM` INSERT + UPDATE |
| `TREATYINDETAIL` | 62 | INSERT |
| `TREATYINDETAILEDM` | 55 | INSERT |
| `M_TREATY_IN` · `M_TREATY_IN_DETAIL` · `M_TREATY_IN_DETAIL_EDM` | 2 masing-masing | INSERT |
| `M_TREATY_IN_EDM` | 4 | INSERT + UPDATE |
| `GCP_IMAGE` | 4 | INSERT |

**`TREATY_IN` — 20 kolom:**
`ID`, `PROPORTIONTYPE`, `TREATYCONTRACTNAME`, `TERITORIALSCOPE`, `COMMENCEMENT`, `TERMINATION`,
`CLASSOFBUSINESS`, `LEADINGREINSSOURCE`, `LEADINGREINSSOURCEID`, `CEDING`, `CEDINGID`,
`LEADINGREINSID`, `NUSARESHAREPCT`, `BROKERAGEPCT`, `INFORMATION`, `POSITIONUSERNAME`, `POSITION`,
`STATUSAKSEPTASI`, `CHOOSESTATUSAKSEPTASI`, `TREATYYEAR`

**`TREATY_IN_EDM` — 24 kolom:** kesembilan belas kolom bisnis yang sama, ditambah
`ID`, `OLDID`, `EDMDATE`, `EDMSTATE`, `EDMMATERIALTYPE`.

Pemeriksaan `PERMINTAAN-DBA-1-UJI-A-SAMPAI-H.sql` terhadap panenan ini: **42 rujukan kolom
terverifikasi tanpa DBA, 0 patut dicurigai.** Tiga tabel **tidak punya sumber apa pun** dan seluruh
rujukannya tetap **BELUM TERPERIKSA** — `M_TREATY_YEAR` (3), `TREATYEXCHANGEYEARLY` (1),
`TREATYINOFFER` (1). Itu **bukan** pernyataan bahwa namanya salah; itu pernyataan bahwa **belum ada
yang dapat memeriksanya**.

> **`TERITORIALSCOPE` dieja tanpa `R` kedua** — `TERITORIAL`, bukan `TERRITORIAL`. Ia ada di
> **kedua** tabel dan di kedua procedure, jadi ia **ejaan sebenarnya**, bukan salah ketik di satu
> tempat. Nama barunya di skema baru dieja benar; peta telusur yang mencatat pasangannya.

---

## 7. Lubang yang dilaporkan, tidak ditambal

| Lubang | Siapa menutup | Yang menagih |
|---|---|---|
| **`TREATYINDETAIL` tidak pernah diperbarui** (§4.1) — belum terukur di data | pemilik data, lewat **Uji AN** | `ADR-0051`: bila ada konsumen hilir, ia selama ini membaca baris ganda atau baris basi |
| **`TREATYINDETAILEDM` kurang tujuh kolom** (§4.2) | pemilik data — daftar konsumen hilir | setiap laporan yang menghitung deductible atau MDP dari tabel datar addendum |
| **Token dapat ditebak** (§5 butir 4) | **pemilik `GCP_IMAGE`, di luar Treaty In** | naik ke `DAFTAR-ESKALASI-MANAJEMEN.md` sebagai butir tersendiri — ia bukan utang migrasi, ia terbuka **sekarang**, di sistem yang berjalan |
| tiga tabel tanpa sumber pemeriksa | **DBA**, lewat Uji A-0 | berkas permintaan DBA tetap menandainya **BELUM TERPERIKSA**, bukan benar |
