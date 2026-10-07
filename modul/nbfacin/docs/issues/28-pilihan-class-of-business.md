# 28: Pilihan Class Of Business — `GET /api/nbfacin/class-of-business` (baca saja)

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah: permintaan sesi `nusantarare-0f`
> 02-10-2026 (form Opportunity, isian Class Of Business; frontend di sesi 0f, backend `nbfacin` di sesi ini), berdasar
> brief sesi itu dan DDL yang ditambahkan work owner.

**What to build:** pembaca baca-saja `POOLDATA.BUSINESS` dan satu endpoint yang mengembalikan SEMUA pilihan Class Of
Business milik satu group business, tanpa paging.

**Asal:**
- DDL READ-ONLY `D:\migrasi\RNM\DDL\BUSINESS.txt` `[terverifikasi]`: `POOLDATA.BUSINESS`, 18 kolom; yang dibaca hanya
  `ID`, `NOTE`, `BUSINESSGROUPID` — ketiganya `VARCHAR2(4000 BYTE)`, tanpa NOT NULL; tanpa PK; satu indeks `BUSINESS_INDEX1` atas `(ID)` — tidak ada indeks atas `BUSINESSGROUPID`
  (dampak kinerja saringan `belum terverifikasi`).
- RD Pega `D:\migrasi\RNM\NB FacIn\ReportDefinition\BrowseBusiness_RD.xml` (kelas `ASM-FW-GISFW-Int-BUSINESS`)
  `[terverifikasi]`: filter A `.ID = Param.ID`; filter B `.Note Contains Param.Note` dengan `pyCaseInsensitive` true;
  filter C `.BusinessGroupID = Param.Group`; logika `(A OR B) AND C`; urut `.ID` DESC; `pyPagingEnabled` false (cocok
  dengan "tanpa paging"); **tanpa** saringan NULL atas `.Note`.
  ⚠️ **Ralat atas brief:** filter B adalah **`Contains` tidak peka huruf**, bukan `=`.
- Pemakai RD `[terverifikasi]` (grep `BrowseBusiness_RD` di `D:\migrasi\RNM\`): di folder `NB FacIn` hanya
  `Section\InputLossRecord_Sec.xml` (autocomplete `.ClassOfBusiness`, `pyDisplayProperty .Note` — properti tampil ini
  milik section, bukan RD); ketiga parameternya (`ID`, `Note`, `Group`) di section itu **kosong** (`<pyValue/>`).
  Folder `RNW Fac In` memuat salinan RD dan section bernama sama — **tidak dibandingkan** (di luar lingkup NB).
  Section form Opportunity yang memakai RD ini **tidak ditemukan** di folder `NB FacIn`.
- Work owner (dikutip sesi 0f, tiket 27): Class Of Business "ditarik dari tabel business, berdasarkan treaty group yang
  dipilih". `[dugaan]` `groupBusinessId` = `GROUPBUSINESSID` akun yang dipilih di popup ChooseAccount (tiket 27) — kaitan
  ini **belum terverifikasi** dari korpus.
- Tangkapan layar Pega `docs/02-layar/tangkapan/class-of-business-pega-02-10-2026.png`: daftar tampak urut alfabetis NOTE.
  Isi daftarnya **tidak** disalin ke kode, uji, atau dokumen — uji memakai data sintetis `UJI-`.

**Blocked by:** —

**Status:** ready-for-human — dibangun 02-10-2026; uji tanpa Oracle hijau; uji terhadap Oracle sungguhan **tidak** dijalankan

## Kontrak

`GET /api/nbfacin/class-of-business?groupBusinessId=<id>` — `groupBusinessId` diteruskan apa adanya (tidak dipangkas,
peka huruf).

| Kode | Kapan | Badan |
| --- | --- | --- |
| 200 | berhasil | `{"baris":[{"id":"…","note":"…"}]}` — semua baris, tanpa paging; teks apa adanya; `baris` selalu larik (bisa kosong) |
| 400 | `groupBusinessId` tidak ada, kosong, hanya spasi, atau > 4000 byte | `{"galat": "..."}` (bentuk `galat.Tulis`; 400/503/500 sama) |
| 503 | aplikasi tanpa basis data | alasan: tabel bisnis `BUSINESS` tidak terbaca |
| 500 | galat Oracle / program | pesan tetap "galat server"; rinciannya hanya di log, tanpa nilai baris |

Nama parameter peka huruf: `groupbusinessid` tidak dikenal (= kosong → 400).

## Yang dibangun

- [x] `repository/bisnis.go` — `SELECT ID, NOTE FROM {skema}.BUSINESS WHERE BUSINESSGROUPID = :1 AND NOTE IS NOT NULL
      ORDER BY NOTE, ID`; tabel lewat `db.Qualify`; parameter terikat; tidak ada `SELECT *`
- [x] `services` — `KelasBisnis` (validasi sebelum basis data, 503 tanpa DB)
- [x] `handlers` — rute, bentuk jawaban di atas
- [x] `MODUL.md` *Tabel warisan* + `docs/STRUKTUR-TABEL-NB-FACIN.md` — `BUSINESS` dibaca, tidak dibuat
- [ ] Uji terhadap Oracle sungguhan — ⛔ tidak dijalankan (aturan kerja: tidak menjalankan apa pun ke Oracle)

## Keputusan agent — menunggu konfirmasi

| # | Keputusan | Dasar |
| --- | --- | --- |
| A74 | Urut **`NOTE`**, bukan `.ID` DESC seperti RD | brief sesi 0f; tangkapan layar Pega tampak alfabetis NOTE — selisih dengan RD dicatat, belum diputus work owner |
| A75 | `ID` sebagai pemutus seri sesudah `NOTE` | urutan deterministik bila dua baris ber-NOTE sama (pola A72) |
| A76 | `groupBusinessId` kosong/spasi saja **atau** > 4000 byte → 400; nilai tidak dipangkas | kosong dari brief; 4000 byte = lebar `BUSINESSGROUPID` (pola A73) |
| A77 | Filter A (`ID`) dan B (`Note Contains`) RD **tidak** dibangun — hanya filter C | brief meminta semua baris satu group; di satu-satunya section pemakai di folder `NB FacIn`, parameter `ID`/`Note` kosong. Penyaringan teks ketik di autocomplete (bila ada) = urusan frontend, `belum terverifikasi` |
| A78 | `NOTE IS NOT NULL` — baris tanpa NOTE dibuang | brief sesi 0f; **tidak ada di RD** (selisih dengan Pega, dicatat). Dugaan alasan: pilihan tanpa teks tak dapat ditampilkan — `belum terverifikasi` |

⚠️ **Batas pengujian:** SQL hanya diuji lewat **teks** (`TestSQLKelasBisnis`), bukan eksekusi Oracle. `[dugaan]` urutan
`ORDER BY NOTE` mengikuti `NLS_SORT` instance — bila `BINARY` (bawaan), huruf besar mendahului huruf kecil; bisa berbeda
dari urutan tangkapan layar Pega. `belum terverifikasi` (tanya DBA).

## Comments

*Tanpa nama orang, tanpa data pelanggan — uji memakai data sintetis berawalan `UJI-`.*
