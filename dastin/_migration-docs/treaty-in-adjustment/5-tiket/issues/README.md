# Papan tiket — Treaty In Adjustment

**Tiket: aktif 10 · tertahan 3 · selesai 0 · mati 0** — **13 tiket**

**Yang sebenarnya dapat dimulai: 0.** Sebabnya di §3, dan ia **seluruhnya ada di papan sebelah**.

> ## PAPAN INI DIPISAHKAN DARI PAPAN INDUK — 25 September 2026
>
> Tiket `01`…`13` sebelumnya berdiri di `treaty-in/5-tiket/issues/` bersama tiket `14`…`64`.
> **Dipisahkan atas perintah pemilik proses.**
>
> ### Ini membalik `GRL-01`, dan itu dicatat, bukan disamarkan
>
> `GRL-01` berbunyi: ***"satu model, satu spesifikasi, satu papan."*** Perintah ronde-ronde
> sebelumnya menegaskannya: *"Jangan membuat papan kedua."*
>
> **Yang dibalik hanya PAPANNYA.** Yang **tidak** berubah:
>
> | Tetap | Di mana |
> |---|---|
> | satu model data | `treaty-in/SPEC-MODEL-DATA.md` — modul ini **tidak punya §10 sendiri** |
> | satu daftar invarian | `treaty-in/SPEC-INVARIAN.md` |
> | satu daftar kemampuan | `treaty-in/5-tiket/DAFTAR-PEKERJAAN.md` — `P-60`…`P-66` **tetap di sana**, tidak ikut pindah |
> | satu daftar asumsi | `treaty-in/5-tiket/ASUMSI-CLEAR.md` |
> | penomoran tiket | `01`…`13` di sini, `14`…`64` di sana — **tidak ada tabrakan nomor**, dan tidak satu pun disusun ulang |
>
> **Ongkos pemisahan ini nyata dan tertulis di §3:** papan ini **tidak dapat menghitung sendiri**
> tiket mana yang dapat dimulai, sebab lima dari tiga belas tiketnya diblokir tiket yang **tidak ada
> di sini**.

---

## 1. Ukuran selesai di papan ini

Sama dengan papan induk, dan **bukan** ukuran batch lapisan aplikasi.

`L-4` menyatakan *"tidak ada spesifikasi layar di mana pun"*, dan lubang itu milik **batch lapisan
aplikasi**. Irisan di sini tegak melalui **lapisan yang ada**:

```
skema  →  constraint dan invarian  →  jalur migrasi  →  uji negatif DAN positif
```

| | |
|---|---|
| **Sebuah tiket selesai bila** | perilakunya **dapat gagal** dan **dapat diperiksa** |
| **Bukan bila** | ia dapat diperagakan di layar |

**Jangan mengarang layar untuk memenuhi bentuk irisan tegak.** Layar yang dikarang akan dipakai
sebagai kriteria selesai oleh orang yang **tidak tahu ia karangan**.

Dan karena "dapat diperagakan" tidak tersedia sebagai bukti, **uji menjadi satu-satunya bukti**:
tiap tiket membawa **uji negatif dan uji positif**. Ketiga kekeliruan lingkup di modul ini **lolos
uji negatif**; yang menangkapnya hanya uji positif — uji yang menuntut **data sah diterima**.

---

## 2. Papan

| # | Tiket | Ditahan oleh |
|---|---|---|
| `01` | Versi baru lahir dari versi berlaku terakhir, rujukan dasarnya disimpan eksplisit | **`14`** ⟦luar⟧ |
| `02` | Materialitas mengunci ruas yang boleh disunting, ditegakkan di sisi simpan | **`14`** ⟦luar⟧ |
| `03` | Tanggal berlaku pada versi ditolak bila di luar periode kontrak | **`14`** ⟦luar⟧ |
| `04` | Dokumen addendum bernomor sendiri, menyentuh beberapa kontrak | **`14`** ⟦luar⟧ · **`DB-16a`** |
| `05` | *Perluas* — kolom nomor urut versi berdiri berdampingan | **`14`** ⟦luar⟧ |
| `06` | Baris selisih berkunci bisnis, mata uang di dalam kuncinya | `01` |
| `07` | Jenis dan materialitas beku sejak diajukan, dan berjejak | `02` · **`DB-20`** |
| `08` | Dokumen addendum membawa tanggal berlakunya sendiri | `04` · **`DB-16b`** |
| `09` | Nomor dokumen baris warisan diisi dari arsip kertas | `04` |
| `10` | *Pindahkan* — nomor urut warisan diisi menurut kronologi | `05` |
| `11` | `INV-69` dan `INV-70` ditegakkan, beserta pemantau kebasiannya | `02` · `06` |
| `12` | *Kerutkan* — urutan tidak lagi diturunkan dari pengenal | `10` |
| `13` | Perubahan reinstatement menghasilkan baris selisih | `11` |

**⟦luar⟧** berarti penahannya **tidak ada di papan ini**. Ia ditulis dengan **nomor aslinya**, tidak
diterjemahkan — yang menahan dari luar harus terlihat berasal dari luar.

---

## 3. Yang sebenarnya dapat dimulai: NOL — dan sebabnya seluruhnya di papan sebelah

Pemeriksaan `K-2` dijalankan: *untuk tiap tiket yang mengaku dapat dimulai, sebutkan entitas yang
disentuhnya dan tunjukkan tiket mana yang membuatnya.*

| Tiket | Entitas yang disentuhnya | Tiket yang membuatnya |
|---|---|---|
| `01`, `02`, `03`, `05` | `VERSI_KONTRAK` | **`14`** — di papan `treaty-in` |
| `04` | `DOKUMEN_ADDENDUM`, merujuk `VERSI_KONTRAK` | **`14`** — idem |

> **Papan ini tidak dapat menghitung angkanya sendiri.** Lima dari tiga belas tiketnya diblokir satu
> tiket yang tidak ada di sini, dan **delapan sisanya menggantung pada kelima itu**.
>
> Itu ongkos pemisahan papan, dan ia dicatat di sini supaya siapa pun yang membaca *"aktif 10"*
> tidak menyimpulkan sepuluh tiket dapat dikerjakan hari ini. **Nol dapat.**

**Begitu `14` mendarat di papan induk**, yang dapat dimulai di sini melompat ke **lima** — `01`,
`02`, `03`, `05`, dan `04` bila `DB-16a` sudah dijawab.

---

## 4. Tepi yang melintasi kedua papan

Dipetakan sebelum pemisahan, bukan sesudah. **Enam tepi, dua arah:**

| Arah | Tepi |
|---|---|
| papan ini **diblokir** papan induk | `01` ← `14` · `02` ← `14` · `03` ← `14` · `04` ← `14` · `05` ← `14` |
| papan induk **diblokir** papan ini | **`40` ← `01`** — *nilai versi sebelumnya berdampingan, di-SELECT bukan disalin* |

**Tepi terakhir yang paling mudah hilang:** ia satu-satunya arah sebaliknya, dan ia ada di papan
yang **tidak memuat tiket `01`**. Ia tertulis di `treaty-in/5-tiket/issues/README.md` §penahan luar.

### Rujukan lain antar-papan, di luar `Blocked by`

Tujuh, dan seluruhnya tetap sah — ia rujukan, bukan penghalang:

```
14 → 04     32 → 13     38 → 11     40 → 01, 06     44 → 08     46 → 07     53 → 07
```

---

## 5. Aturan papan

Sama dengan papan induk, dan keduanya **tidak boleh berbeda**:

- **Status adalah medan di dalam berkas tiket**, bukan lokasi foldernya. Papan **dibangkitkan dari
  berkas** dan **tidak pernah menghapus tiket**.
- **Empat golongan, tidak ada yang kelima.**
- **Penahan dari luar papan ditulis dengan nama aslinya** — `DB-16a`, `DB-20`, dan sekarang juga
  **nomor tiket papan induk**.
- **Penahan yang selesai dicoret, tidak dihapus.**
- **Nomor yang lompat bukan kekeliruan**, dan **penomoran tidak disusun ulang** — `01`…`13` tetap
  `01`…`13` meski papannya pindah.
- **Pencacah diberi label bendanya.** Papan ini mencacah **tiket Adjustment**; jangan dijumlahkan
  dengan pencacah papan induk tanpa menyebut keduanya.

## 6. Menunggui, tidak menahan

| # | Menunggu | Yang berubah bila jawabannya lain |
|---|---|---|
| `01` | `REV-3` | **keadaan mana yang dihitung "berlaku"** — bukan bahwa ada penunjuk dasar |
| `02` | `DB-20` | **titik beku** materialitas, dan titik beku milik `07` |
| `09` | `DB-16a` | **bentuk kolom** dokumennya, yang ikut `04` |

## 7. Menahan sungguhan — tiga, dan ketiga pertanyaannya BELUM DIKIRIM

| # | Penahan | Kenapa ia menahan |
|---|---|---|
| `04` | `DB-16a` | bila dokumen ternyata dikirim ke luar, **jenis bernilai tiga** dan penanda kembali. Itu kolom, bukan penyetelan |
| `07` | `DB-20` | titik beku adalah **pokok** tiketnya |
| `08` | `DB-16b` | tanggalnya **adalah** pokok tiketnya, dan `KTV-2` menuntut kolomnya **dicabut sebelum data masuk** bila dibantah |

## 8. Asumsi yang dianggap clear

Daftarnya **tetap satu**, di `treaty-in/5-tiket/ASUMSI-CLEAR.md` — ia tidak ikut dipisah, sebab
sebagian kodenya dipakai kedua papan (`KTV-A` dipakai 19 tiket papan induk, `DB-20` dipakai empat
tiket di sini).

Sapuan yang mengeluarkan daftarnya kini menuntut **dua folder**:

```
grep -r "DIASUMSIKAN-CLEAR(" treaty-in/5-tiket/issues/ treaty-in-adjustment/5-tiket/issues/
```

> **Itu ongkos kedua dari pemisahan papan**, dan ia lebih halus daripada yang pertama: sapuan yang
> hanya menyisir satu folder akan mengembalikan daftar yang **terlihat lengkap**.
