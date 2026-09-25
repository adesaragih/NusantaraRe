> Modul  : Komite Claim Non Prop · Ronde 05 · 2026-09-20
> Peran  : interogator
> Masukan: `PENGETAHUAN.md` (2026-09-19) · `PUTUSAN-01.md` (2026-09-20) · `INVENTARIS-BUKTI.md` (2026-09-20) · `FAKTA-A1B.md` (2026-09-20) · transkrip sesi `5696a074-…jsonl` (2026-09-20) · 59 XML `Komite Claim Non Prop/` + 279 XML `Claim Non Prop/` (ekspor 2026-09-08/09)
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

# 00 · LINGKUP RONDE 5

## 1. Apa yang digrilling

Ronde ini menggrilling **jalur yang menghubungkan klaim dengan sirkulasi komite**, dari
tiga arah yang sampai ronde 4 hanya disebut namanya:

1. **Pembentukan sirkulasi** (aliran A-1b) — kedua activity pembuat anak, dibaca langkah
   demi langkah, bukan dari ringkasan.
2. **Perutean jenjang** (aliran A-1a) — `KomiteRouter`, satu-satunya rule yang memilih
   siapa menerima giliran, yang sampai ronde 4 belum pernah dibuka.
3. **Tiga pertanyaan lama yang jawabannya ternyata sudah dipegang** — `Q-6`, `Q-7`, `Q-3`
   pada `PENGETAHUAN.md` §14, ditambah `Q-2`, yang menurut aturan biaya `03-LUBANG`
   langkah 1 wajib **dibaca dan dilaporkan sebagai fakta, bukan ditanyakan**.

Sebab ronde ini ada: `07-AUDIT` ronde sebelumnya menutup tahap 1 dengan status
`SELESAI BERSYARAT`, dan sisanya memuat empat perkara yang pemegangnya adalah **kamu
sendiri** — bukan DBA, bukan admin Pega. Perkara yang dapat dijawab dengan membaca tidak
boleh menginap satu ronde lagi.

## 2. Bukti yang dipakai

| Bukti | Asal | Tanggal | Dipakai untuk |
|---|---|---|---|
| `Komite Claim Non Prop/Activity/KomiteRouter.xml` | ekspor Pega | rule diubah 2020-06-15 | N-01, N-02, N-03, N-12 |
| `Komite Claim Non Prop/Activity/KomitePostAdjustment.xml` | ekspor Pega | rule diubah 2025-09-08 | N-02, gerbang `IsKomiteLoop` |
| `Komite Claim Non Prop/When/IsKomiteLoop.xml` | ekspor Pega | ver 01-01-07 | gerbang lingkar jenjang |
| `Komite Claim Non Prop/Flow/KomiteTreaty_Flow.xml` | ekspor Pega | 2020-02-27 | N-14 |
| `Komite Claim Non Prop/Section/ReinstatementPremiumDetails.xml` | ekspor Pega | — | N-15 |
| `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml` | ekspor Pega | rule diubah 2026-08-04 | N-04, N-06, N-07, N-08, N-09, N-10, N-11 |
| `Claim Non Prop/Activity/CreateChildKomiteCloseNP_Act.xml` | ekspor Pega | rule diubah 2025-10-01 | N-04, N-05, N-07 |
| Transkrip sesi `5696a074-4b6f-422e-a2aa-86e9a9c4b1ff.jsonl` | berkas di disk | 2026-09-20 | memulihkan bunyi ketetapan D/E/G/H/C ke `KETETAPAN.md` |

Pembaca langkah yang dipakai: `scratchpad/alat/dump_act.py`. Kode transisi dibaca mentah
dari XML bila penentu putusan — lihat N-01, yang berdiri di atas `pyStepsTransParamsWhenTrue`
dan `…WhenFalse` yang keduanya bernilai `2`.

## 3. Bukti yang sengaja TIDAK dipakai

| Tidak dipakai | Alasan |
|---|---|
| `HitServiceToKasirKMT_Act.xml`, `SaveOSClaim_SQL`, `SaveXOLClaim_SQL` | di balik pagar A-4 dan A-5 — lihat `REGISTER-PAGAR.md` |
| Isi surat dan dokumen (`SendEmailKlaimRejectClose_KMT.xml` selain fakta penulisnya) | di balik pagar A-6; satu akibat yang menyentuhnya dicatat sebagai lubang bertahan, bukan sebagai temuan berperlakuan |
| `MEMORI_PEMAHAMAN.MD` | ditulis sebelum folder Komite dibedah berkas-per-berkas; bukan bukti, melainkan ringkasan |
| Ingatan atas ronde 1–4 | bunyi ketetapan dipulihkan dari transkrip di disk, bukan ditulis ulang dari ingatan |

## 4. Aliran: yang terbuka dan yang beku

| Aliran | Isi | Keadaan ronde ini | Pagar |
|---|---|---|---|
| A-1a | mekanika keputusan & perutean jenjang | **TERBUKA** | — |
| A-1b | pembentukan sirkulasi | **TERBUKA** (sejak ronde 4) | — |
| A-2 | layar & aturan medan | TERBUKA, tidak digrilling ronde ini | — |
| A-3 | penomoran akseptasi | TERBUKA TERBATAS | PG-03 |
| A-4 | kiriman ke Kasir | **BEKU** | PG-04 |
| A-5 | akseptasi OS & pembalikan | **BEKU** | PG-05 |
| A-6 | isi dokumen & surat | **BEKU RINGAN** | PG-06 |

Nomor pagar merujuk `REGISTER-PAGAR.md`. Tidak ada pagar baru yang dipasang ronde ini.

## 5. Deret nomor ronde ini

| Deret | Arti | Catatan tabrakan |
|---|---|---|
| `N-01 … N-15` | temuan ronde 5 | deret baru; **tidak** memakai `G` yang sudah dipakai dua arti (temuan `G-01…G-20` dan ketetapan ronde 3 yang kini `J-1…J-4`) |
| `T5-01 …` | temuan turunan penilai, `06-PUTUSAN.md` | — |
| `Q5-1 …` | pertanyaan ronde 5, `03-LUBANG.md` | `Q-x` lama tetap milik `PENGETAHUAN.md` §14 dan `PUTUSAN-01.md` |
| `P5-01 …` | skenario paritas, `04-PARITAS.md` | — |
| `K5-1 …` | ketetapan ronde 5, masuk `KETETAPAN.md` | deret `D`/`E`/`J`/`H`/`C` tidak dipakai ulang |

## 6. Aturan berhenti

Penggalian berhenti dan **batas dinyatakan**, bukan ditafsir, pada empat keadaan:

1. **Arti nilai enum tidak ada di ekspor.** `Rule-Obj-FieldValue` tidak ikut diekspor
   (`INVENTARIS-BUKTI.md` §1 kolom "TIDAK diliput"). Maka domain `.Type` dan arti empat
   nilai `FlagProrate` dicatat sebagai **rentang yang terpakai**, bukan sebagai arti.
2. **Pemicu di luar 338 berkas.** Bila sebuah nama tidak muncul sebagai penulis di mana
   pun dalam dua folder, yang dicatat adalah "nol penulis **di dalam ekspor**", bukan
   "tidak ada penulis".
3. **Semantik mesin Pega.** Bila putusan bergantung pada cara Pega membandingkan `""`
   dengan `0`, temuan diturunkan ke `RAGU` — berkas tidak dapat memutuskannya.
4. **Pagar.** Begitu sebuah akibat mendarat di A-4, A-5, atau isi A-6, penggalian
   berhenti di situ dan akibatnya dicatat sebagai lubang bertahan. **Alasan pun tunduk
   pada pagar, bukan hanya kesimpulan.**
