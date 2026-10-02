# Papan tiket modul `treatyin` — empat belas dari 51

Papan **lengkap** modul ini — 51 tiket, `14`…`64`, beserta rantai penahannya — ada di
`D:\XML_NURE\_migration-docs\treaty-in\5-tiket\issues\README.md`. Folder ini memuat **tiket yang sudah
dikerjakan saja**, disalin apa adanya supaya kode dan tiketnya ditinjau bersama.

| # | Tiket | Keadaan |
| --- | --- | --- |
| `14` | Kontrak dan versi pertamanya berdiri | skema selesai |
| `15` | Himpunan acuan bertambah tanpa mengubah arti | skema selesai |
| `20` | Beberapa mata uang berlaku, masing-masing berkurs dan berperiode | skema selesai |
| `22` | Retensi cedant per kelompok treaty per mata uang | skema selesai |
| `23` | EGNPI per kelompok treaty per mata uang | skema selesai |
| `24` | Portofolio masuk dan keluar yang menyertai kontrak | skema selesai |
| `25` | Periode pelaporan beserta jatuh temponya | skema selesai; **INV-55 belum ditegakkan** |
| `26` | Periode akumulasi, ditolak bila di luar periode kontrak | skema selesai; **INV-56 belum ditegakkan** |
| `27` | Termin pembayaran premi bernomor urut, per mata uang | skema selesai |
| `28` | Skala ko-asuransi sebagai beberapa baris | skema selesai |
| `29` | Batas tanggungan untuk bahaya apa pun di daftar | skema selesai |
| `30` | Dokumen dilampirkan, dirujuk bukan disalin | skema selesai |
| `31` | Layer beserta limit, deductible, MDP, dan pemulihan limitnya | skema selesai |
| `39` | Jejak perubahan: siapa mengubah fakta apa dan kapan | entitas selesai; **mekanisme pengisiannya belum** |

`14` dan `15` adalah **PEMBUAT PERTAMA**, dan keduanya satu-satunya yang dapat dimulai hari pertama:
`14` membuat `KONTRAK` dan `VERSI_KONTRAK` yang seluruh papan menggantung padanya, `15` membuat keenam
tabel acuan yang setiap kunci asing menunggunya. Keduanya mendarat, lalu dua belas tiket yang hanya
tertahan keduanya ikut dikerjakan — **23 tabel, 192 kolom** seluruhnya.

**Yang terbuka berikutnya:** `21` (dari `20`), `32` `33` `34` (dari `31`), dan `16` `17` `18` `19`
`41` `42` `43` yang hanya tertahan `14`. Tiket `38` tetap tertahan `Uji AD` dan `L-3`; `44` tertahan
`KTV-A` dan bertenggat — ia yang memuat data pertama, dan sesudahnya presisi tidak dapat dipersempit
lagi.

## Ukuran selesai di papan ini

Papan ini **tidak** memakai irisan tegak skema · API · UI · uji. `L-4` menyatakan *"tidak ada
spesifikasi layar di mana pun"*, dan papan melarang mengarang layar untuk memenuhi bentuk itu. Yang
berlaku:

```
skema  →  constraint dan invarian  →  jalur migrasi  →  uji negatif DAN positif
```

| | |
| --- | --- |
| **Sebuah tiket selesai bila** | perilakunya **dapat gagal** dan **dapat diperiksa** |
| **Bukan bila** | ia dapat diperagakan di layar |

## Yang BELUM selesai dari kedua tiket ini

Daftar periksa keduanya menuntut **uji negatif dan uji positif**, dan uji itu menuntut Oracle yang
dapat dijangkau. `L-3` menyatakan tidak ada instans Oracle yang terjangkau — **tidak satu baris pun
DDL di sini pernah dijalankan di mana pun**. Yang sudah ada: bentuk tabelnya, constraint-nya, jalur
mundurnya, dan kecocokannya dengan `KAMUS-KOLOM.md` (diperiksa `TestKolomDDLCocokDenganStruktur`).

| Daftar periksa yang menunggu Oracle | Tiket |
| --- | --- |
| simpan versi yatim → ditolak kunci asing | `14` |
| dua versi bernomor urut sama → ditolak `INV-04` | `14` |
| kontrak bertanggal terbalik → ditolak `INV-53` *(di services, lihat DDL `401`)* | `14` |
| satu kontrak dengan **tiga** versi diterima, lapisan bekunya tidak disalin | `14` |
| kode ganda ditolak di **keenam** tabel acuan → `INV-68` | `15` |
| bahaya **kesembilan** ditambahkan dan langsung dapat dipakai, tanpa perubahan skema | `15` |
