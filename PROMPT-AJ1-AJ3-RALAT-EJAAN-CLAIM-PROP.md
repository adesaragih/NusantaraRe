# PERINTAH WAJIB — AJ1 s/d AJ3 · Ralat bukti AC 8 · Claim Prop

Tempel seluruh blok di bawah ke Claude eksekutor.

---

## GANTI KONTEKS

Modul: **Claim Prop**. Berkas yang diedit: **hanya**
`.scratch/claim-prop/issues/00-prefactor-skema-relasional-transaksi-uang.md` (406 baris).

Korpus `D:\XML\RNM_BRD\` READ-ONLY, **329 berkas**. Tulis HANYA ke
`D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\`. `D:\XML\nusantara-re\` DI-BLACKLIST. Jangan tulis kode
Go/React. Jangan sentuh korpus. **Jangan sentuh `spec.md` dan 15 tiket lain.**

## AI1–AI4 LULUS — saya jalankan `sensus.py` sendiri atas berkasmu

1638 · 124 AC · 41 tag = 41 frasa · 20 jebakan · ⚠️ 61 = 61 · selisih dua arah NIHIL · 16 tiket ·
124 rujukan unik · `Menutup` 33 AC · tiket 00 406 baris · `CLAUDE.md` 247 baris. Identik sampai
baris terakhir — dua mesin, dua sesi, satu alat.

## KENAPA ADA RALAT INI — kesalahan saya, bukan kamu

**Brief AI2 saya salah menimbang.** Work owner menunjukkannya: **Oracle melipat identifier tanpa
kutip menjadi huruf besar.** `treatyinproduction` dan `TREATYINPRODUCTION` menunjuk objek yang sama.
Beda huruf besar-kecil **bukan cacat dan tidak berakibat apa pun** pada sistem — di Go tinggal
ditulis satu cara.

Saya menggelembungkan **dua** objek menjadi **lima** dengan ikut menghitung beda huruf, lalu
menyuruhmu menanamnya sebagai bukti AC 8. Kamu mengerjakannya persis seperti yang saya minta.
Yang salah briefnya.

Hasil `sensus.py` sesudah kedua hal itu dipisah:

```
A. DITULIS DUA BENTUK — telanjang di satu rule, ber-prefiks di rule lain     2
   TREATYBUSINESS · TREATYINPRODUCTION
   ⚠️ INI yang punya akibat: nama telanjang bergantung pada schema bawaan koneksi,
      jadi ia bisa menunjuk objek lain bila Go menyambung sebagai pengguna berbeda.

B. HANYA BEDA HURUF BESAR-KECIL                                              3
   DIRECTTOKASIR_LOG · JSON_POLIS · T_STORAGE_IMAGE
   nol akibat. Yang terganggu hanya sensus yang peka huruf besar-kecil.
```

Bukti AC 8 yang sah adalah **A saja — dua objek** — dan itu sudah lengkap tertulis sejak AF1.

`sensus.py` di `OUTPUT_HASIL_RNM\` **sudah saya perbarui** (versi ketiga): kategori B kini dicetak
terpisah dengan keterangan nol akibat, dan objek dikunci pada **nama dasar** (prefiks dibuang).

## LANGKAH 0 — WAJIB

```
grep -c "Total LIMA objek" issues/00-*.md          → 1
grep -c "jsoN_polis" issues/00-*.md                → 2
grep -c "36 objek Oracle" issues/00-*.md           → 1
```

## ATURAN PELAPORAN — WAJIB

**Tepat 3 baris**: `AJ1 SELESAI <bukti yang bisa di-grep>` atau `AJ1 TIDAK <alasan>`.

---

# AJ1 — Cabut klaim "lima objek", ganti dengan pemisahan yang benar

**Ganti seluruh blok ini** — dari kalimat `**Dan bukan hanya dua objek itu.**` sampai dengan alinea
yang berakhir `…hanya catatan ejaannya.` (tabel lima baris dan kedua alinea ⚠️ sesudahnya ikut
terganti):

**menjadi:**

```
⚠️ **Yang berakibat hanya DUA objek di atas.** `[terverifikasi]` Pemeriksaan ulang `sensus.py --tabel`
atas teks SQL di **58 dari 58** berkas `RDBList/` memisahkan dua hal yang sempat tercampur:

| | Objek | Akibat |
| --- | --- | --- |
| **A. Ditulis dua bentuk** — telanjang di satu rule, ber-prefiks di rule lain | `TREATYINPRODUCTION` · `TREATYBUSINESS` | ⚠️ **NYATA.** Nama telanjang bergantung pada **schema bawaan koneksi**; ia dapat menunjuk objek lain bila Go menyambung sebagai pengguna yang berbeda. **Inilah isi AC 8.** |
| **B. Hanya beda huruf besar-kecil** | `JSON_POLIS` · `DIRECTTOKASIR_LOG` · `T_STORAGE_IMAGE` | **NOL akibat.** Oracle **melipat** identifier tanpa kutip menjadi huruf besar, jadi `treatyinproduction` dan `TREATYINPRODUCTION` adalah objek yang sama. Bukan cacat, bukan temuan. |

`[keputusan work owner]` **2026-09-18 — beda huruf besar-kecil bukan masalah.** Di Go seluruh nama
ditulis satu cara, konsisten. Yang wajib diperbaiki hanya **prefiks schema**, sesuai AC 8.

⚠️ **Satu akibat yang tetap berlaku, tetapi hanya bagi pembaca korpus:** sensus nama tabel **wajib
tidak peka huruf besar-kecil**. Yang peka akan melaporkan objek yang sama sebagai beberapa tabel
berbeda. Pakai `sensus.py --tabel`, yang mengunci pada **nama dasar** dan memisahkan A dari B.
```

⛔ **Jangan menyentuh tabel dua baris hasil AF1** (`treaty-in production` dan `treaty business`
dengan kolom *Telanjang di* / *Berprefiks di*) yang ada **di atas** blok ini. Itu bukti AC 8 yang
sah dan tetap berdiri apa adanya.

# AJ2 — 36 baris ejaan, 34 objek berbeda

Angka **36** menghitung `TREATYINPRODUCTION` dan `TREATYBUSINESS` **masing-masing dua kali**, karena
satu baris telanjang dan satu baris ber-prefiks. Objek Oracle yang benar-benar berbeda ada **34**.

Ini bukan soal huruf besar-kecil — ini akibat langsung dari temuan A.

Ganti kalimat pengantar tabel:

```
sekarang : **36 objek Oracle**, dengan aksi dan rule yang menyentuhnya:
jadi     : **36 baris ejaan — 34 objek Oracle yang berbeda**, dengan aksi dan rule yang
           menyentuhnya. Selisih duanya: `TREATYINPRODUCTION` dan `TREATYBUSINESS` muncul dua
           kali, sekali telanjang dan sekali ber-prefiks. Yang menyiapkan pemetaan skema
           menyiapkan **34**, bukan 36.
```

**Daftar tabelnya sendiri tidak diubah** — tetap 36 baris, supaya tiap ejaan tetap terlihat beserta
rule-nya.

# AJ3 — Jalankan ulang dan laporkan

```
py sensus.py .scratch\claim-prop --korpus "D:\XML\RNM_BRD\Claim Prop" --tabel --uji 124:PW
```

(dengan `PYTHONIOENCODING=utf-8` seperti yang sudah kamu catat di `CLAUDE.md` §4a)

Tempel keluaran bagian inventarisnya apa adanya — ia sekarang mencetak A dan B terpisah.

Invarian yang **tidak boleh berubah** — prompt ini tidak menyentuh AC mana pun:

```
AC                                    124
tag [penyimpangan sadar]               41
⚠️ spec = ⚠️ tiket                      61 = 61, selisih dua arah NIHIL
rujukan di 16 tiket                   124 · unik 124 · nol duplikat
Menutup: tiket 00                      33 AC
butir [terbuka] ringan                  5
AC ber-[terbuka]                        3  (104 · 115 · 120)
graf Blocked by                         tidak berubah
```

---

## YANG TIDAK BOLEH DILAKUKAN

- Jangan menyentuh `spec.md`. Ralat ini hanya di tiket 00.
- Jangan menyentuh tabel dua baris hasil AF1 — itu bukti AC 8 yang sah.
- Jangan mengubah isi daftar 36 baris, hanya kalimat pengantarnya.
- Jangan menambah, menghapus, atau memindahkan AC.
- Jangan menulis `CREATE TABLE`. Rancangan skema masih menunggu work owner.

## SETELAH AJ1–AJ3

Laporkan 3 baris sesuai ATURAN PELAPORAN, ditambah:

1. Keluaran inventaris `sensus.py` apa adanya — bagian A dan B.
2. Daftar invarian di AJ3, masing-masing lulus atau tidak.
3. Satu kalimat: apakah tiket 00 kini bebas dari klaim yang berlebihan.

Sesudah ini **Claim Prop tertutup**, tinggal menunggu rancangan skema relasional dari work owner.
Modul berikutnya **Claim Non Prop** — jangan dimulai sebelum briefnya ada.
