# 26: Form "Opportunity" (dibuka dari **Create opportunity** di portal)

> Disusun agent (sesi `nusantarare-0f`) dari **tangkapan layar Pega** kiriman work owner 02-10-2026 —
> **bukan hasil `/to-tickets`**. Bukti: `docs/02-layar/tangkapan/opportunity-pega-02-10-2026.png`
> (md5 `4233929bf124ea9fd111d30148b52325`).

**What to build:** form "Opportunity" persis seperti tangkapan layar: label, urutan, kelompok, tanda wajib,
dan nilai awal yang terlihat.

## ⛔ Batas bukti

`[terverifikasi]` Form ini **tidak ada di korpus**: label-labelnya (`Business Prospect Name`,
`Estimated Closing Date`, `Type Of Facultative`, `New Group Business`, `Opportunity Source`, …) dicari di
`D:\migrasi\RNM\` dan `D:\XML\RNM\NB Facin\` → **0 berkas**; properti padanannya juga tidak ditemukan di
section mana pun yang terkait portal. Tombol `Create opportunity` membuat kasus kelas
`D_crmAppExtPage.WorkClass_Opportunity` lewat `D_crmAppExtPage.CreateWork_Flow`
(`SFAPortalOpportunitiesHeader.xml` L2486–L2504), dan `D_crmAppExtPage` tidak diekspor.
`[dugaan]` layar Opportunity bawaan Pega CRM. Satu-satunya sumber = tangkapan layar.

## Isi layar (VERBATIM dari tangkapan layar)

| Kelompok | Label | Kontrol di gambar | Wajib | Nilai awal terlihat |
| --- | --- | --- | :-: | --- |
| judul | `Opportunity` | teks | | |
| baris atas | `Estimated Closing Date` | tanggal (ikon kalender) | ✱ | kosong |
| kiri | `Business Prospect Name` | kotak teks lebar | ✱ | kosong |
| kiri | `Group Business` | tiga tombol: `Search Group Business` · `New Company Detail` · `New Group Business` | | |
| kiri | `Class Of Business` | kotak isian (bersegitiga biru di pojok) | ✱ | kosong |
| kiri | `Type Of Inward` · `Type Of Facultative` | dua dropdown berdampingan | ✱ · ✱ | `Facultative` · `Facultative In` |
| kanan | `Phase` | dropdown | ✱ | `Proposal` |
| kanan | `Stage` | kotak berisi | | `Opportunity` |
| kanan | `Opportunity Source` | dropdown | | `Select...` — **13 pilihan** dari tangkapan layar dropdown terbuka (lihat C-1) |
| kanan | `Business Status` | dropdown | ✱ | `New Business` |
| bawah | `Description` | kotak teks panjang | | kosong |

## Keputusan agent (menunggu konfirmasi work owner)

- **C-1 — dropdown hanya berisi nilai yang TERLIHAT di gambar** (satu nilai, sebagai nilai awal).
  **Diperbarui 02-10-2026:** `Opportunity Source` kini berisi 13 pilihan VERBATIM dan berurutan dari tangkapan
  layar dropdown terbuka (`docs/02-layar/tangkapan/opportunity-source-pega-02-10-2026.png`, md5
  `d0df211cf73b4f227dcb4804c3b11c0f`): Iklan · Analisa Referral · Rujukan Pelanggan · Direct Mail · Email ·
  Rujukan Karyawan · Telepon Masuk · Partner · Seminar · Sosial Media · Pameran · Web · Whatsapp. Nilai awal
  tetap `Select...` (sorotan "Rujukan Karyawan" di gambar = posisi kursor, bukan nilai terpilih). ⚠️ `belum
  terverifikasi`: nilai yang DISIMPAN Pega per pilihan (teks atau kode) — di sini nilai = teks tampil. Dropdown
  lain masih menunggu tangkapan layar terbuka.
  Daftar lengkapnya `belum terverifikasi`; menulis daftar karangan di React dilarang pola `inti`
  (`Pilih`: "daftar yang ditulis di sini pasti karangan").
- **C-7 — keadaan awal dari `D:\migrasi\RNM\DDL\HALAMAN DEPAN NB.JPG`** (md5 `9f4813a2c9c56ba12306232436453ca7`;
  nama orang di gambar tidak disalin): ada medan **`Owner`** di kanan Estimated Closing Date berisi nama pengguna
  yang login (ditampilkan sebagai teks, dari sesi login); **`Type Of Inward` awalnya kosong** berteks `Choose Type
  Of Inward`, dan **`Type Of Facultative` belum tampil**. `[dugaan]` dari dua gambar: Type Of Facultative baru
  tampil sesudah Type Of Inward = `Facultative` (tangkapan layar pertama).
- **C-8 — `Search Group Business` membuka popup `ChooseAccount`** — VERBATIM dari tangkapan layar 02-10-2026
  (md5 `e4944778178e583eeac1ce1b7d87b0f3`; gambar TIDAK disalin ke repo karena memuat nama pelanggan): kotak
  `Search`, tombol `Search`, grid bernomor `Insured ID` · `Insured Name` · `Group Business` dengan tombol
  `Choose` per baris, paging angka + `Next`. Rule `ChooseAccount` tidak ada di korpus. Sumber data menurut work
  owner: **`T_M_ACCOUNT`** — DDL-nya belum ada di `D:\migrasi\RNM\DDL\`, jadi grid tampil `BelumTersedia` dan
  Search nonaktif sampai DDL + endpoint ada (backend nbfacin, sesi c3).
- **C-8 diperbarui 02-10-2026 (jawaban work owner):** sumber `POOLDATA.T_M_ACCOUNT` (DDL
  `D:\migrasi\RNM\DDL\T_M_ACCOUNT.txt`: ID, GROUPBUSINESSID, GROUPBUSINESS, INSUREDID, INSUREDNAME); Insured ID =
  `INSUREDID`, Insured Name = `INSUREDNAME`, Group Business = `GROUPBUSINESS`; pencarian "mengandung". Frontend
  memanggil `GET /api/nbfacin/account?cari=&halaman=` (kontrak di `api.ts`); endpoint dibangun sesi c3 (tiket 27).
  Sampai endpoint ada, popup menampilkan galat backend apa adanya.
- **C-9 — keputusan agent:** daftar dimuat saat popup dibuka dengan kotak kosong (= semua baris); Search memuat ulang
  dari halaman 1; paging memakai kata cari yang terakhir dikirim. Kolom yang dicari (tiga kolom tampil), ukuran halaman
  (20), dan urutan — keputusan agent di tiket 27, menunggu konfirmasi.
- **C-10 — sesudah Choose** (tangkapan layar work owner 02-10-2026): ketiga tombol Group Business diganti teks
  `GROUPBUSINESS` terpilih + ikon roda gigi; `[dugaan]` roda gigi membuka popup lagi untuk mengganti pilihan.
  ⛔ **Class Of Business belum diisi:** menurut work owner ia "ditarik dari tabel business, berdasarkan treaty group
  yang dipilih" — tabel dan DDL-nya belum ada di `D:\migrasi\RNM\DDL\`, dan tidak ditemukan di korpus (yang ada
  hanya `m_businessfield`, RDB `ConvertBusinessField`, untuk BusinessFieldNote). Menunggu work owner.
- **C-11 — Class Of Business** (jawaban work owner 02-10-2026: "ditarik dari tabel business, berdasarkan treaty group
  yang dipilih"; DDL `D:\migrasi\RNM\DDL\BUSINESS.txt`; tangkapan layar daftar terbuka
  `docs/02-layar/tangkapan/class-of-business-pega-02-10-2026.png`): kotak isian dengan daftar saran. `[terverifikasi]`
  korpus `NB FacIn\Section\InputLossRecord_Sec.xml` (~L2669) memberi `.ClassOfBusiness` pyUIElement `autocomplete` dari
  `BrowseBusiness_RD` (pyValue `.Note`; filter `(ID = Param.ID OR Note = Param.Note) AND BusinessGroupID = Param.Group`).
  `[dugaan]` form Opportunity memakai RD yang sama dan Param.Group = `GROUPBUSINESSID` akun terpilih. Frontend memanggil
  `GET /api/nbfacin/class-of-business?groupBusinessId=` (tiket 28, sesi c3); saran dimuat ulang tiap group berganti dan
  isian lama dikosongkan. Gambar memperlihatkan urutan abjad, sedangkan RD sort ID DESC — urutan dari backend.
- **Work owner 02-10-2026 atas tiket 27:** pencarian akun PEKA huruf besar-kecil, 15 baris per halaman.
- **C-2 — dua tombol Group Business lain nonaktif (`New Company Detail`, `New Group Business`)** (aksinya tidak diketahui).
- **C-3 — `Stage` ditampilkan read-only berisi `Opportunity`** `[dugaan]`: di gambar ia kotak berisi
  tanpa panah dropdown.
- **C-4 — `Class Of Business` kotak teks biasa**; arti segitiga biru (`[dugaan]` autocomplete / smart
  prompt) dan sumber pilihannya belum terverifikasi.
- **C-5 — tanpa tombol simpan.** Gambar tidak menampilkan tombol simpan/submit (`[dugaan]` terpotong atau
  milik harness). Tidak ada endpoint; isian tidak dikirim ke mana pun.
- **C-6 — nol catatan pengembang di layar** (sama dengan tiket 25 B-4).

## Untuk memverifikasi (minta ke work owner / pemilik export Pega)

Live UI atas form ini di Pega DEV → nama Section + Class, lalu Save As XML: section form, data page
`D_crmAppExtPage`, flow/flow action Create Opportunity, dan sumber dropdown (Field Value / Data Page /
Report Definition).

**Status:** ready-for-human — dibangun 02-10-2026 dari tangkapan layar, menunggu tinjauan work owner atas C-1…C-6

- [x] Label, urutan, kelompok, tanda wajib = tangkapan layar
- [x] Nilai awal terlihat terpasang; tanpa daftar dropdown karangan
- [x] Dibuka dari tombol `Create opportunity` di portal (tiket 25)
- [ ] Daftar dropdown, aksi tombol Group Business, validasi, dan penyimpanan — menunggu XML form

## Comments
