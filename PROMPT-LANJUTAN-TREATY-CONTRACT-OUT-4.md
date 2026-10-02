# PROMPT — LANJUTAN 4 MODUL **Treaty Contract Out** *(folder `OUTPUT_HASIL_RNM`, cabang `main` @ `eee4eaf` atau lebih baru)*: **bentuk teks tanggal mengikuti DATA warisan, jawaban OQ-TCO-01/23/24/25/26 dari katalog DEV**

> Hanya konteks Treaty Contract Out. Brief modul dan lanjutan 1–3 tetap berlaku kecuali yang dibantah di sini. Commit selalu dengan jalur
> eksplisit (`git commit --only <berkas>`). `App.tsx` berisi suntingan work owner — lihat §3.

## 0. VERIFIKASI LANJUTAN 3 *(asisten, 29-09-2026)*

| Klaim | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 6 commit kerja + 1 dokumen `e6905ed` → `eee4eaf` | ada | ✅ |
| nol migrasi Treaty; nol nama `T_TREATY…` di kode | 0 berkas `3xx_*`; 0 rujukan di luar uji | ✅ |
| menu satu butir `Treaty Contract Out` | `daftarMenu.ts:71` satu entri `tco-tahun` | ✅ |
| "commit sesi lain menyapu penghapusan" | **benar, dan itu kesalahan asisten**: commit brief `5415aab` dibuat dengan `git commit` tanpa jalur, sehingga 21 penghapusan yang sudah Anda *stage* ikut masuk | ✅ diakui; sejak `51faf89` asisten juga memakai jalur eksplisit |
| "`.scratch/cadangan/` hilang" | dihapus **asisten** *(salinan patch `App.tsx` ada di luar repo)* supaya pohon `main` bersih | ✅ bukan kehilangan |

## 1. TEMUAN DARI DATA DEV *(agregat bentuk teks — digit diganti `9`, nol nilai disalin)*

| Kolom warisan | Bentuk di DEV | Yang kode tulis sekarang | Akibat |
| --- | --- | --- | --- |
| `TREATYYEAR.STARTDATE`, `ENDDATE` | **182/182 baris `99999999`** *(`YYYYMMDD`)* | stempel Pega `YYYYMMDDTHHMMSS.mmm GMT` pukul 00:00 WIB *(`StempelTanggalJakartaTCO`, `tco_tahun.go:228/252`)* | ⛔ **salah bentuk**: baris baru berbeda dari 182 baris lama; hilir yang membaca `YYYYMMDD` gagal diam-diam |
| `TREATYYEAR.TGLUPDATE` | 182/182 `99999999T999999.999 GMT` | `StempelPegaTCO` | ✅ cocok |
| `TREATYREINSURER.STARTDATE`, `ENDDATE` | **430/430 kosong** | stempel 00:00 WIB *(`tco_reinsurer.go:328/353`)* | ⛔ Pega **tidak** menulisnya |
| `TREATYREINSURER.USERID`/`TGLUPDATE` | **0/430 terisi** | diisi layanan | ⛔ = OQ-TCO-25 |
| `TREATYBUSINESS.USERID`/`TGLUPDATE` | **2/4.621 terisi** | diisi layanan | ⛔ = OQ-TCO-25 |
| `PROPORTIONALARRG.RP`, `PCT` *(teks desimal)* | `RP`: 930 bertitik, **0 berkoma**, 2.678 terisi; `PCT`: 520 bertitik, **0 berkoma** | `TulisDesimalWarisanTCO` bertitik | ✅ **OQ-TCO-23 terjawab: titik** |
| prosedur `PEGA_M_ATTACHMENT` *(dipanggil `InsertAtatchment_Sql`)* | `ALL_SOURCE` **terbaca** dari akun `POOLDATA`: 35 baris; tabel yang disebut `M_ATTACHMENTTREATY` dan `ID_COUNT` | kode menulis `M_ATTACHMENTTREATY_2` *(dari `GetAllAttachment2_Sql`)* | ⚠️ **OQ-TCO-24**: rekonsiliasi — baca badan prosedur dari `ALL_SOURCE` *(`SELECT text FROM all_source WHERE owner='POOLDATA' AND name='PEGA_M_ATTACHMENT' ORDER BY line`)*, catat di `dba-procedures.md`; tentukan tabel yang benar-benar ditulis dan yang dibaca `GetAllAttachment2`; bila keduanya berbeda, itu temuan warisan yang dicatat, bukan disatukan diam-diam |

## 2. KEPUTUSAN `[asisten dari data; veto work owner]`

| Butir | Isi |
| --- | --- |
| **OQ-TCO-01** | **ditutup dari data**: `TREATYYEAR.STARTDATE/ENDDATE` ditulis **`YYYYMMDD`** *(delapan angka, tanggal kalender, tanpa jam/zona)*; pembaca menerima `YYYYMMDD` *(bentuk warisan)* dan menolak bentuk lain dengan galat berkata-kata; `TGLUPDATE` tetap stempel Pega; `TREATYREINSURER.STARTDATE/ENDDATE` **tidak ditulis** *(dibiarkan kosong seperti 430 baris lama)* kecuali RDB penulis yang hidup mengisinya — buktikan dari `SaveMasterTreatyReinsurer_SQL`/prosedurnya |
| **OQ-TCO-23** | **ditutup**: titik |
| **OQ-TCO-25** | **kosongkan seperti Pega** *(data: 0/430 dan 2/4.621 terisi)* — `USERID`/`TGLUPDATE` reinsurer dan business tidak diisi layanan; identitas pelaku tetap tercatat di log aplikasi, bukan di kolom warisan |
| **OQ-TCO-26** | **tiru** `Update_T_Storage_SQL` *(konsisten dengan Claim Life migrasi 020 yang mengisi `TANGGAL_UPLOAD`)* |
| **OQ-TCO-24** | §1 baris terakhir — dikerjakan executor dari `ALL_SOURCE`, bukan menunggu DBA |
| **OQ-TCO-22** | tetap untuk work owner |

## 3. HALAMAN YATIM DI `App.tsx`

Rute `tco-kontrak`/`tco-klausul` dan dua halaman pembungkusnya kini tanpa butir menu. Suntingan work owner di `App.tsx` tinggal **satu baris
kosong**. Bila work owner menulis *"buang suntingan App.tsx"*, asisten mengembalikannya ke HEAD dan executor membuang kedua rute serta
halamannya di paket 3; bila tidak, paket 3 dilewati dan dicatat.

## 4. URUTAN — satu commit per paket

| # | Paket | Commit |
| ---: | --- | --- |
| 1 | bentuk tanggal `YYYYMMDD` + reinsurer tanggal kosong + uji dua arah *(contoh bentuk dari data: 8 angka; stempel `…T…GMT` ditolak di kolom tanggal tahun)*; ralat tiket 03/05 | `treaty-contract-out: OQ-TCO-01 — tanggal tahun YYYYMMDD seperti data warisan` |
| 2 | OQ-25 kosongkan · OQ-26 tiru · OQ-24 rekonsiliasi lampiran | `treaty-contract-out: OQ-TCO-24/25/26 — lampiran menurut PEGA_M_ATTACHMENT, kolom pelaku kosong seperti Pega` |
| 3 | halaman yatim *(bila diizinkan §3)* | `treaty-contract-out: buang halaman yatim tco-kontrak/tco-klausul` |
| 4 | dokumen + register OQ + uji penuh | `docs: treaty-contract-out — lanjutan 4` |

## 5. LAPORAN

Satu pesan: tabel **paket → commit → RDB/katalog → kode**; OQ ditutup; angka uji dengan dan tanpa tag `db`; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 29 September 2026 sesudah verifikasi `e6905ed..eee4eaf` dan tiga kueri agregat ke DEV (bentuk teks tanggal `TREATYYEAR`/`TREATYREINSURER`,
cacah kolom pelaku, pemisah desimal `PROPORTIONALARRG`, keterbacaan `ALL_SOURCE` `PEGA_M_ATTACHMENT`).*
