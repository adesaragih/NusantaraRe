# PERINTAH WAJIB — AI1 s/d AI4 · PENUTUP Claim Prop

Tempel seluruh blok di bawah ke Claude eksekutor.

---

## GANTI KONTEKS

Modul: **Claim Prop**. Berkas yang diedit: `.scratch/claim-prop/spec.md`,
`.scratch/claim-prop/issues/00-prefactor-skema-relasional-transaksi-uang.md`, dan `CLAUDE.md`.

Korpus `D:\XML\RNM_BRD\` READ-ONLY, **329 berkas** (angka dari `spec.md`). Tulis HANYA ke
`D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\`. `D:\XML\nusantara-re\` DI-BLACKLIST. Bukti wajib
`path berkas dari akar modul + nama rule + nomor step Pega`. **JANGAN kutip nomor baris XML.**
Jangan hapus teks lama di `spec.md` — ralat sebagai blok `⚠️ RALAT 2026-09-18`. Jangan tulis kode
Go/React. Jangan sentuh korpus. Jangan sentuh 15 tiket lain.

## AH1–AH5 LULUS — dihitung ulang, tiap angka cocok

123 AC · nol AC tanpa penanda · 40 tag = 40 frasa, himpunan identik · 20 jebakan · ⚠️ spec 60 =
⚠️ tiket 60 · selisih dua arah NIHIL · frasa terlipat NIHIL · AC ber-`[terbuka]` tepat 104 · 115 ·
120 · butir ringan 5 · 16 tiket · 123 rujukan unik · graf tidak berubah · nol `CREATE TABLE`.

**Catatan 328 versus 329: kamu benar, saya salah — dan saya menulis dua angka berbeda di dalam satu
prompt.** Sandbox saya memang 328; korpus 329. Yang tidak ada di salinan saya persis satu berkas,
**`When/IsCustomBonds.xml`**, ditambahkan 17 September sesudah saya menyalin. Ketiga hasil "nol
penulis" tidak berubah karena rule `When` tidak melakukan `Property-Set`. Catatanmu di `spec.md`
sudah benar, tidak perlu diapa-apakan.

## LANGKAH 0 — WAJIB

```
grep -c "^124\. " .scratch/claim-prop/spec.md    → 0
grep -c "^123\. " .scratch/claim-prop/spec.md    → 1
grep -c "jsoN_polis" .scratch/claim-prop/spec.md → 0
ls sensus.py                                      → ada di OUTPUT_HASIL_RNM
```

Bila AC 124 sudah ada, **berhenti** dan laporkan.

## ATURAN PELAPORAN — WAJIB

**Tepat 4 baris**: `AI1 SELESAI <bukti yang bisa di-grep>` atau `AI1 TIDAK <alasan>`.

---

# AI1 — AC 124 disetujui, tulis ke spec dan tiket 00

⚠️ `[keputusan work owner]` **2026-09-18 — usulanmu disetujui.** Tulis sebagai **AC 124** di akhir
bab `## Acceptance Criteria`, dengan tanda berkurung tepat sesudah nomor:

```
124. `[penyimpangan sadar]` ⚠️ `[data DBA]` **Panggilan stored procedure penulis adalah langkah
     TERAKHIR sebelum `COMMIT`** dalam satu transaksi, tanpa pekerjaan penting yang masih
     menggantung sebelumnya. Aplikasi memeriksa `StsSimpan` — **`1` berhasil, `0` gagal** — dan
     pada `0` melaporkan **gagal**. **Alasan menyimpang:** ketiga procedure memasang
     `EXCEPTION WHEN OTHERS THEN ROLLBACK` di dalam dirinya sendiri; `ROLLBACK` telanjang di Oracle
     membatalkan **seluruh transaksi**, termasuk pekerjaan aplikasi sebelumnya, dan
     `ROLLBACK TO SAVEPOINT` tidak menolong karena yang dipasang bukan itu. Pada `StsSimpan = 0`
     transaksinya sudah dibatalkan procedure, sehingga `COMMIT` sesudahnya tidak menyimpan apa pun.
     AC 6 menuntut **hasil**, bukan **urutan**; urutan itulah yang wajib di sini.
```

Aturannya:

- Frasa `penyimpangan sadar` **wajib utuh dalam satu baris**, tidak terlipat.
- Tautkan ke **AC 6** dan ke **Implementation Decision §17** yang sudah kamu tulis di AH3 —
  sebagai rujukan silang, jangan ubah bunyi keduanya.
- Perbarui ringkasan di kepala `spec.md`.

Lalu di **tiket 00**:

- Tambahkan satu baris AC di sub-bab `### Batas transaksi penomoran`, bertanda ⚠️, berlabel
  `**Alasan menyimpang:**`, ditutup rujukan `*(AC 124 spec)*`.
- `**Menutup:**` bertambah `· AC 124`, dan jumlahnya **naik satu**.

# AI2 — Dua bukti ejaan baru untuk AC 8

`[terverifikasi]` Temuan `sensus.py`, dari teks SQL di **58 dari 58** berkas `RDBList/`, dengan
pencocokan **tidak peka huruf besar-kecil**:

```
JSON_POLIS                   TIGA ejaan   JSON_POLIS · json_polis · jsoN_polis
POOLDATA.DIRECTTOKASIR_LOG   DUA ejaan    POOLDATA.DIRECTTOKASIR_LOG · pooldata.directtokasir_log
```

`jsoN_polis` — huruf `N` kapital di tengah — tidak akan pernah tertebak oleh siapa pun yang
mengetikkan namanya dari ingatan.

Tambahkan keduanya ke bukti AC 8 di **tiket 00**, di bawah tabel ejaan yang sudah ada dari AF1.
Bunyi yang ditambahkan:

> Dan bukan hanya dua objek itu. `JSON_POLIS` muncul dalam **tiga** ejaan
> (`JSON_POLIS` · `json_polis` · `jsoN_polis`) dan `POOLDATA.DIRECTTOKASIR_LOG` dalam **dua**.
> Total **empat objek** dieja lebih dari satu cara di dalam satu modul.
>
> ⚠️ Daftar 36 objek yang ditanam di AE2 **meng-`upper()` seluruh nama**, sehingga perbedaan ejaan
> ini tidak terlihat dari daftar itu. Siapa pun yang menyensus ulang wajib memakai `sensus.py`
> `--tabel`, yang menampilkan ejaan apa adanya.

**Jangan mengubah daftar 36 objek itu sendiri** — jumlahnya tetap 36, yang bertambah hanya
catatan ejaannya.

# AI3 — Pakai `sensus.py` mulai sekarang

`OUTPUT_HASIL_RNM\sensus.py` sudah ada. Jendelanya tertanam di dalam kode, dicetak di setiap hasil,
dan **ia menolak jalan tanpa uji instrumen**.

Tambahkan bagian ini ke `CLAUDE.md` (dibaca semua Claude di repo ini):

> ## Sensus wajib lewat `sensus.py`
>
> Setiap sensus atas `spec.md`, `issues/`, atau korpus XML dijalankan dengan
> `OUTPUT_HASIL_RNM\sensus.py`, bukan dengan grep atau regex yang ditulis ulang tiap kali.
> Alasannya: jendela sensus yang dikarang ulang per pertanyaan akan berbeda-beda, dan perbedaannya
> baru ketahuan setelah ada yang menghitung ulang.
>
> ```
> python3 sensus.py .scratch/<modul> --uji 70:PW 60:- 82:PW
> python3 sensus.py .scratch/<modul> --korpus "D:\XML\RNM_BRD\<Modul>" --tabel
> python3 sensus.py .scratch/<modul> --korpus "..." --cari Halaman.Properti
> ```
>
> `--uji` **wajib** — dua tiga AC yang jawabannya sudah diketahui. Tanpa itu skripnya tidak
> mencetak apa pun. Verifikator yang instrumennya sendiri belum diverifikasi memproduksi tuduhan
> palsu.
>
> Empat jebakan yang sudah ditutup skrip ini: unit AC hanya baris pertama · blok kutipan `>` RALAT
> ikut terbaca · frasa terlipat di pergantian baris · nama dicocokkan peka huruf besar-kecil.
> Satu lagi yang diperingatkan tetapi tidak bisa dipaksakan: properti korpus **wajib** dicari
> dengan awalan halaman (`InputData.CARI16`, bukan `CARI16`).

# AI4 — Sensus penutup, dijalankan dengan `sensus.py`

```
python3 sensus.py .scratch/claim-prop --uji 70:PW 60:- 82:PW 124:PW
python3 sensus.py .scratch/claim-prop --korpus "D:\XML\RNM_BRD\Claim Prop" --tabel --uji 124:PW
```

Tempel keluarannya apa adanya ke laporan — termasuk baris JENDELA dan baris jumlah berkas.

Angka yang **berubah karena AI1**, beserta sebabnya:

```
                              sebelum   sesudah   sebab
AC                              123       124     AI1 menambah satu
tag [penyimpangan sadar]         40        41     AC 124 bertanda
⚠️ spec = ⚠️ tiket                60        61     AC 124 ber-⚠️
rujukan di 16 tiket             123       124     AI1
Menutup: tiket 00                32        33     AI1
```

Angka yang **tidak boleh berubah**:

```
butir [terbuka] ringan            5
AC ber-[terbuka] aktif            3   (104 · 115 · 120)
jebakan                          20
jumlah tiket                     16
rujukan unik = rujukan total          → nol duplikat
cakupan AC 1–124                      → nol yang belum tertaut
graf Blocked by                       → tidak berubah, 00 satu-satunya titik mulai
frasa penyimpangan sadar terlipat     → NIHIL
CREATE TABLE di tiket mana pun        → 0
```

---

## YANG TIDAK BOLEH DILAKUKAN

- Jangan menambah AC selain AC 124.
- Jangan mengubah daftar 36 objek Oracle — yang bertambah hanya catatan ejaannya.
- Jangan mengubah bunyi AC 6 dan Implementation Decision §17; AC 124 hanya menautkan.
- Jangan mengarang nama tabel relasional dan jangan menulis `CREATE TABLE` — rancangan skema masih
  menunggu work owner.
- Jangan menutup kelima butir `[terbuka]` ringan. Semuanya tetap terbuka dan tetap tidak memblokir.

## SETELAH AI1–AI4

Laporkan 4 baris sesuai ATURAN PELAPORAN, ditambah:

1. Keluaran `sensus.py` apa adanya, kedua perintah.
2. Tabel angka berubah dan angka tetap, terisi hasil sebenarnya.
3. Satu kalimat: apakah Claim Prop tertutup dan hanya menunggu rancangan skema dari work owner.

Sesudah itu **modul berikutnya: Claim Non Prop.** Briefnya menyusul terpisah — jangan dimulai
sebelum briefnya ada.
