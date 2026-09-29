# PROMPT — LANJUTAN 6 MODUL **Treaty Contract Out** *(folder `OUTPUT_HASIL_RNM`, cabang `main`)*: **perbaikan cacat "Backend tidak terhubung" di List Description → Show TreatyDesc**

> Laporan work owner 29-09-2026: membuka `List Description` → `Show TreatyDesc` menampilkan *"Backend tidak terhubung (127.0.0.1:8080)"*
> padahal backend sudah dijalankan. Commit dengan jalur eksplisit. Jangan menyentuh suntingan frontend sesi lain yang belum di-commit
> *(`index.html`, `App.tsx`, `PagarGalat.tsx`, `labels.ts`, `styles.css`, `Shell.tsx`, `dasar.tsx`, `hooks/useTema.ts`, `lib/tema*.ts`)*.

## 0. DIAGNOSIS — asisten, 29-09-2026 *(backend `:8080` hidup sejak 16.51, `healthz` sehat)*

| Permintaan layar | Jawaban |
| --- | --- |
| `GET /api/treaty-contract-out/jenis-klausul?isXol=0/1` | 200 |
| `GET …/tahun/1000682/klausul?descId=10001&induk=` | 200 |
| `GET …/tahun/1000682` | 200 |
| **`GET …/tahun/1000682/kurs`** dan **`…/kurs/konversi`** | **503** `{"galat":"services: master kurs atau mata uang tidak dapat dipakai: models: nilai master kurs tidak dapat diurai: STARTDATE \"20190801T00000.000 GMT\" bukan bentuk YYYYMMDD\"T\"HH24MISS.FF3 TZR"}` |

**Dua cacat, keduanya bentuk lintas-lapis yang sudah dikenal:**

1. **Pengurai tanggal kurs lebih ketat dari Oracle.** Data warisan *(agregat bentuk, DEV)*: `STARTDATE` **129** baris
   `99999999T999999.999 GMT` dan **11** baris `99999999T99999.999 GMT` *(jam hanya lima angka)*; `ENDDATE` 140/140 normal. Pega tidak pernah
   mengurainya di Java: `RDBList/GetMasterKursList.xml` menyerahkannya ke **Oracle** —
   `to_date({InputData.CARI1},'YYYYMMDD') BETWEEN trunc(TO_TIMESTAMP_TZ(STARTDATE,'YYYYMMDD"T"HH24MISS.FF3 TZR')) …` — dan Oracle
   menerima lima angka jam. `models/tco_kurs.go:104` menolaknya, sehingga **satu** baris cacat mematikan kurs **seluruh** tahun.
2. **Klien menyebut setiap 503 "backend tidak terhubung".** `frontend/src/lib/keadaanGalat.ts:98–104` memperlakukan 502/503/504 sebagai
   backend mati, **termasuk** 503 yang datang dari backend sendiri dengan badan `{"galat": …}`. Pesan backend yang tepat tertelan, dan
   pemakai disuruh menyalakan backend yang sudah menyala.

## 1. PERBAIKAN

| # | Isi | Commit |
| ---: | --- | --- |
| 1 | **Kurs seperti Pega**: pembacaan dan pembandingan tanggal master kurs dilakukan **di SQL** dengan ekspresi RDB yang sama *(`trunc(TO_TIMESTAMP_TZ(STARTDATE,'YYYYMMDD"T"HH24MISS.FF3 TZR'))`, dan padanannya untuk `ENDDATE`)*, sehingga Go menerima `DATE`, bukan teks; pengurai teks Go untuk kolom ini dibuang atau hanya dipakai bila SQL tidak dapat dipakai. **Periksa tabelnya**: RDB membaca `treatyexchange` dengan `Quarter='0'`, sedangkan agregat di atas dari `TREATYEXCHANGEYEARLY` — pakai tabel yang RDB hidup pakai, catat. Uji: contoh bentuk lima angka jam diterima; bentuk rusak lain *(huruf, panjang salah)* tetap ditolak dengan galat berkata-kata; satu baris cacat **tidak** mematikan baris lain bila Oracle sendiri menolaknya *(saring per baris, laporkan cacahnya)*. Ralat tiket 11 | `treaty-contract-out: kurs dibaca seperti GetMasterKursList — tanggal diurai Oracle` |
| 2 | **Klien membedakan 503 backend dari 503 proxy**: `backendMati()` hanya bila jawaban **tidak** membawa badan galat dari backend *(kode `BACKEND_TIDAK_TERJANGKAU`, atau 502/503/504 tanpa `galat`)*; 503 **dengan** `{"galat": …}` → `galat-api` dengan pesan backend. Uji kontrak dua sisi *(handler 503 menulis `galat`; klien menampilkan pesannya)*; uji yang ada untuk proxy mati tetap hijau. Berkas bersama — perubahan aditif, dilaporkan | `fix: 503 dari backend menampilkan pesannya, bukan "backend tidak terhubung"` |

## 2. LAPORAN

Satu pesan pendek: dua commit, bukti uji *(mutasi merah dulu)*, angka uji dengan dan tanpa tag `db`, dan langkah uji ulang di layar
*(`List Description` → `Show TreatyDesc` pada tahun `1000682`)*.

---

*Disusun 29 September 2026 dari panggilan langsung ke `:8080` dan lewat proxy Vite `:5173` (enam rute layar klausul), agregat bentuk tanggal
`TREATYEXCHANGEYEARLY`, `RDBList/GetMasterKursList.xml`, `models/tco_kurs.go:104`, dan `lib/keadaanGalat.ts:98–104`.*
