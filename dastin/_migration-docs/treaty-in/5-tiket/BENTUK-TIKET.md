# Bentuk tiket — perbandingan dengan folder contoh

**Tanggal:** 24 September 2026
**Contoh:** `_migration-docs/komite-claim-non-prop/3-to-tickets/issues` — **12 tiket + README papan**

> **Contoh bentuk, BUKAN sumber isi.** Tidak ada satu pun klaim domain dari folder itu yang berlaku
> di Treaty In, persis seperti perlakuan atas `4-erd-dan-tabel-datar` milik modul yang sama.
>
> **Dan bila bentuk contoh tidak memuat medan yang prompt kita wajibkan, YANG MENANG PROMPT KITA.**
> Contohnya lebih tua daripada aturan yang lahir di sesi ini; ketiadaan sebuah medan di sana bukan
> pernyataan bahwa medan itu tidak perlu.

**Perbedaannya dilaporkan ke dua arah, tidak dipilih diam-diam.**

---

## 1. Bentuk contoh, apa adanya

Sembilan medan, **dipakai di kedua belas tiket tanpa kecuali** — keseragaman itu sendiri patut
ditiru:

```
---
status: aktif | tertahan | selesai | selesai-sebagian
---

# NN: Judul kalimat penuh

*Asal: ‹berkas spesifikasi› bagian ‹bagian›, ‹cakupan›.*

**What to build:**   kemampuannya, lalu artefak yang dibangun
**Persyaratan:**     daftar `S-xxx`
**Tidak termasuk:**  pengecualian, masing-masing dengan alasannya
**Jalur gagal:**     masukan salah → akibat yang teramati, dipisah titik-tengah
**Uji:**             uji yang sudah ada, dan yang **BARU**
**Menggantikan:**    aturan sistem lama yang digantikannya, dengan kode cacatnya
**Blocked by:**      daftar, atau "None (can start immediately)"
**Dasar:**           DECIDED(…) · EVIDENCED(…)

- [ ] kriteria selesai, satu per baris

**Ketidakpastian:** apa yang belum pasti, dan apa yang berubah bila jawabannya lain
```

Nama berkas: `NN-judul-dipotong-sekitar-60-karakter.md`.

---

## 2. Medan yang prompt kita WAJIBKAN dan contoh TIDAK PUNYA

**Lima. Seluruhnya lahir dari sesi ini, dan seluruhnya tetap dipakai.**

| Medan | Kenapa ia tidak ada di contoh | Kenapa ia tetap wajib di sini |
|---|---|---|
| **GOLONGAN** — di frontmatter; "dari apa" dipegang `MENGGANTIKAN:` (§6.1) | modul contoh membangun sesuatu yang **belum pernah ada**, jadi seluruh tiketnya otomatis baru | Treaty In **memindahkan sistem yang berjalan**. Tanpa golongan, pembaca tidak dapat tahu apakah sebuah tiket memindahkan perilaku atau menciptakannya — dan itu menentukan cara mengujinya |
| **KENAPA BEGINI** | sebagian isinya ada di `Dasar:` sebagai kode keputusan | kode keputusan **menunjuk**, ia tidak **menjelaskan**. Pembaca yang tidak ikut sesi mana pun tidak akan membuka tujuh ADR untuk satu tiket. Ini medan terpenting di seluruh berkas |
| **CARA MENYALAKANNYA** — empat butir | tidak ada golongan BARU di sana, jadi tidak ada yang perlu dinyalakan bertahap | tujuh tiket Treaty In bergolongan BARU, dan **mode peringatan tanpa pembaca dan tanpa ambang berangka berlangsung selamanya** |
| **PEMBUAT PERTAMA** | tiket `01` di sana **adalah** pembuat pertama, dan tidak ada yang membandingkan taksirannya | irisan tegak membuat tiket pertama sebuah entitas 3–5× lebih besar. Tanpa penandanya, taksirannya akan dibandingkan dengan tiket biasa dan seluruh perkiraan melenceng |
| **PENGHALANG** — apa yang ditunggu, **siapa dapat menjawab**, apa yang berubah | `Ketidakpastian:` memuat dua dari tiga | **"siapa dapat menjawab" hilang**, dan itu yang mengubah catatan menjadi tagihan. Sejalan dengan `CONTEXT.md` §2.9b: lubang membawa pemilik dan saat penagihan |

**`INVARIAN YANG HARUS DIPENUHI`** ada padanannya tetapi tidak setara: `Dasar:` menyebut ADR dan
kode keputusan, bukan invarian bernomor. Treaty In punya **63 invarian bernomor** dan sesi DDL
menulis constraint dari kalimatnya — jadi tiketnya menyebut `INV-xx`, bukan hanya ADR.

---

## 3. Medan BAGUS yang contoh punya dan prompt kita TIDAK

**Lima, dan saya usulkan keempatnya diambil.** Pemilik proses menyatakan mungkin mengambilnya.

| Medan contoh | Kenapa ia layak diambil |
|---|---|
| **`status:` di frontmatter** | *"Status adalah field di dalam berkas, bukan lokasi foldernya"* — papan dibangkitkan dari berkas dan **tidak pernah menghapus tiket**. Ia membuat papan dapat dihitung ulang kapan saja, dan itu persis disiplin "hitungan yang diverifikasi ulang, bukan diingat" |
| **`Jalur gagal:`** | masukan salah → akibat yang **teramati**. Prompt kita menuntut kriteria selesai *"yang dapat gagal"*, dan medan ini adalah **tempat kegagalannya ditulis sebagai kalimat**, bukan tersembunyi di dalam daftar centang. Untuk Treaty In ia juga tempat alami menuliskan **invarian mana yang menolak** |
| **`Menggantikan:`** | menyebut aturan sistem lama yang digantikan, **dengan kode cacatnya**. Untuk modul migrasi ini nilainya lebih besar daripada di modul contoh: ia jawaban langsung atas *"apakah ini PELESTARIAN"*, dan ia mencegah dua tiket diam-diam menggantikan aturan yang sama |
| **`Dasar: DECIDED(…) · EVIDENCED(…)`** | **memisahkan yang diputuskan dari yang dibuktikan**, di dalam satu baris. Itu persis pembedaan yang `KEPUTUSAN-TANPA-VERIFIKASI.md` tegakkan di tingkat berkas, dan membawanya ke tingkat tiket murah sekali |
| **`Uji:`** | uji yang sudah ada dan yang **BARU**, hidup **di dalam** tiketnya. Sejalan dengan aturan kita bahwa uji data **bukan** tiket tersendiri |

---

## 4. Aturan papan yang diambil apa adanya

Kelimanya ada di README contoh dan tidak satu pun bertabrakan dengan prompt kita:

| Aturan | Kenapa |
|---|---|
| status adalah field, papan dibangkitkan, **tiket tidak pernah dihapus** | papan dapat dihitung ulang; tidak ada keadaan yang hanya ada di susunan folder |
| **penahan dari luar papan memakai nama aslinya** — `L-3`, `CO-3` — tidak diterjemahkan jadi nomor tiket | yang menahan dari luar **harus terlihat berasal dari luar**, dan itu langsung berguna: empat kemampuan Treaty In tertahan hal di luar papan |
| penahan yang selesai **dicoret, tidak dihapus** | rantainya tetap terbaca |
| **nomor yang lompat bukan kekeliruan** | menyusun ulang penomoran memutus setiap rujukan yang sudah ada |
| **"Menunggui, tidak menahan"** dipisahkan dari penahan sungguhan | *"supaya papan tidak berteriak serigala"*. Ini persis pembedaan yang §1.1 `DAFTAR-PEKERJAAN.md` buat antara **tidak terhalang luar** dan **boleh dikerjakan** |

---

## 5. Satu hal yang TIDAK ditiru

Contoh memakai kode HTTP (`422`, `403`) di `Jalur gagal:`. Untuk Treaty In itu **terlalu dini**:
bentuk antarmuka belum diputuskan di mana pun, dan **L-4 menyatakan tidak ada spesifikasi layar**.
Jalur gagalnya ditulis sebagai **akibat yang teramati dan invarian yang menolaknya** —
*"ditolak INV-50, dan galatnya menyebut induk penyebaran mana yang tidak berjumlah seratus"* —
bukan sebagai kode status.

Menyalin kode HTTP dari modul lain berarti mengarang bentuk antarmuka Treaty In tanpa menyadarinya,
dan itu **klaim domain**, bukan bentuk.

---

## 6. Bentuk yang berlaku untuk Treaty In — **TIGA BELAS medan**

Gabungan keduanya, dengan prompt kita menang di setiap tabrakan.

```
---
status:   aktif | tertahan | selesai | selesai-sebagian
golongan: pelestarian | perubahan | baru
---

# NN: Judul

**SATU KALIMAT:**   kemampuannya, satu kalimat, pelakunya disebut
**KENAPA BEGINI:**  satu sampai tiga kalimat — medan terpenting, §7 menjaganya
**ASALNYA DARI MANA:** berkas dan bagian, kemampuan P-xx
**INVARIAN:**       `INV-xx` yang harus dipenuhi
**MENGGANTIKAN:**   aturan sistem lama + kode cacatnya — atau **"tidak ada"**
**JALUR GAGAL:**    masukan salah → akibat teramati + invarian yang menolak
**UJI:**            yang sudah ada, dan yang BARU
**BERGANTUNG PADA:** tiket lain · penahan luar dengan nama aslinya
**YANG TEGAS BUKAN BAGIAN TIKET INI:** pengecualian + alasannya
**PENGHALANG:**     apa yang ditunggu · SIAPA dapat menjawab · apa yang berubah
**PEMBUAT PERTAMA:** ya/tidak — taksirannya tidak dibandingkan dengan tiket biasa
**CARA MENYALAKANNYA:** empat butir — wajib bila `golongan: baru`
**DASAR:**          DECIDED(…) · EVIDENCED(…)

- [ ] kriteria selesai — pemeriksaan yang DAPAT GAGAL
```

### 6.1 `GOLONGAN` hanya di frontmatter — satu fakta, satu tempat

Bentuk yang saya usulkan semula memuatnya **dua kali**: `golongan:` di frontmatter **dan** medan
`**GOLONGAN:**` di badan. **Dua tempat untuk satu fakta, dan salah satunya akan basi** — kegagalan
yang sama yang sudah dua kali diperbaiki minggu ini di berkas spesifikasi, hanya kali ini **di dalam
satu berkas**.

- `golongan:` **tinggal di frontmatter saja** — ia yang dibaca papan, dan papan harus dapat dihitung
  ulang dari berkas.
- Medan `**GOLONGAN:**` di badan **dihapus**.

**Dan bagian "PERUBAHAN — dari apa" sudah dipegang medan lain.** `MENGGANTIKAN:` menyebut aturan
sistem lama yang digantikan **beserta kode cacatnya** — itu jawaban yang sama atas *"berubah dari
apa"*, hanya lebih lengkap. Maka yang bertahan yang lebih lengkap, dan gantinya ia diberi
**kewajiban keterisian**:

| `golongan:` | `MENGGANTIKAN:` |
|---|---|
| `perubahan` | **wajib terisi** — aturan lama yang digantikan, dengan kode cacatnya |
| `pelestarian` | **wajib terisi** — aturan lama yang dipindahkan perilakunya |
| `baru` | **wajib berbunyi "tidak ada"** — kekosongan **yang dinyatakan**, bukan kekosongan |

### 6.2 Dua pasang medan yang bertetangga dekat, dan bedanya

Ditulis supaya yang mengisi tidak menaruh hal yang sama di dua tempat.

| Pasangan | Bedanya |
|---|---|
| `ASALNYA DARI MANA` vs `DASAR:` | **`ASALNYA` menyebut DI MANA TERTULIS** — berkas, bagian, nomor kemampuan. Ia penunjuk tempat. **`DASAR:` menyebut KENAPA DIPERCAYA** — keputusan mana yang mengikat, bukti mana yang membacanya. Satu tiket dapat berasal dari satu bagian tetapi berdasar lima ADR, dan sebaliknya |
| `BERGANTUNG PADA` vs `PENGHALANG` | **`BERGANTUNG PADA` adalah PEKERJAAN yang belum selesai** — ada yang dapat menyelesaikannya dengan bekerja. **`PENGHALANG` adalah PERTANYAAN yang belum dijawab** — tidak ada jumlah pekerjaan yang menutupnya; seseorang harus menjawab. Itu sebabnya `PENGHALANG` wajib menyebut **siapa dapat menjawab**, dan `BERGANTUNG PADA` tidak |

---

## 7. Penjaga supaya tiga belas medan tidak menjadi tiga belas formalitas

Tiga belas medan wajib itu banyak, dan yang selalu terjadi pada borang panjang adalah **medan
terpenting diisi dengan mengulang medan lain** — justru karena ia yang paling penting.

### 7.1 Uji `KENAPA BEGINI`

> **`KENAPA BEGINI` tidak boleh dapat ditulis oleh orang yang hanya membaca medan lain di tiket
> itu.** Bila isinya dapat disusun dari `SATU KALIMAT` + `INVARIAN` + `ASALNYA DARI MANA`, ia
> **kosong**.

Yang **sah**: menyebut **apa yang terjadi di sistem lama**, atau **akibat yang akan mengejutkan**
pembaca yang tidak ikut sesi mana pun.

| | |
|---|---|
| **lolos** | *"Kurs dibekukan pada versi, tidak dibaca saat tampil. Sistem lama membaca kurs tahun berjalan, sehingga angka yang sudah disetujui berubah sendiri setiap ganti tahun."* |
| **gagal** | *"Supaya memenuhi INV-43."* |

Ujinya dapat dijalankan peninjau mana pun, dan itu yang membuatnya berguna: ia tidak menuntut tahu
domainnya, hanya menuntut membaca tiketnya.

### 7.2 Aturan `DASAR:` — masa lalu tidak diputuskan, ia dibuktikan

Pembedaan DECIDED versus EVIDENCED hanya berguna bila ada yang tidak boleh masuk salah satunya:

> **Setiap klaim tentang apa yang SISTEM LAMA LAKUKAN menuntut `EVIDENCED`, bukan `DECIDED`.**
> `DECIDED` sah untuk putusan bisnis dan untuk rancangan kita. Ia **tidak pernah** sah sebagai dasar
> sebuah pernyataan tentang masa lalu — **masa lalu tidak diputuskan, ia dibuktikan.**

**Dan `EVIDENCED(...)` menyebut EKSPOR MANA, bukan hanya berkasnya:**

```
DASAR: DECIDED(ADR-0036, ADR-0049) · EVIDENCED(SetSpreadName@ekspor-2026-09,
       FetchQSfromMasterXOL@ekspor-2026-09)
```

Sebabnya **L-9**: asal lingkungan ekspor belum diketahui, dan bila jawabannya kelak "QA", seseorang
harus tahu **tiket mana** yang buktinya berasal dari sana. Tanpa penanda ekspor, jawabannya menuntut
membaca ulang setiap tiket; dengan penanda itu, **satu sapuan atas nama ekspor** mengeluarkan
daftarnya.

Dan bila ekspor kedua diperoleh (`LUBANG-SPESIFIKASI.md` §8), penanda itu langsung berguna dengan
cara lain: tiket yang buktinya berasal dari salah satu **lima belas berkas yang sudah dibandingkan**
dapat ditandai **terverifikasi**, dan sisanya tidak.

Akibat langsungnya, dan ia uji yang tajam:

> **Tiket bergolongan `pelestarian` yang seluruh `DASAR:`-nya DECIDED adalah tiket yang mengaku
> melestarikan sesuatu tanpa satu pun bukti bahwa sesuatu itu ada.**

---

## 8. Satu hal yang TIDAK ditiru, dan aturan umum di baliknya

Contoh memakai kode HTTP (`422`, `403`) di `Jalur gagal:`. Untuk Treaty In itu **terlalu dini**:
bentuk antarmuka belum diputuskan di mana pun, dan **L-4 menyatakan tidak ada spesifikasi layar**.
Jalur gagalnya ditulis sebagai **akibat yang teramati dan invarian yang menolaknya** —
*"ditolak INV-50, dan galatnya menyebut induk penyebaran mana yang tidak berjumlah seratus"*.

Aturan umumnya, dan ia berlaku setiap kali folder modul lain ditiru:

> **Yang ditiru BENTUKNYA; yang menyebut DUNIA NYATA adalah ISI, dan isi tidak ikut.**

Ujinya mudah saat ragu: **apakah ia menyebut sesuatu di luar berkas itu?** Nama medan, susunan
bagian, cara menulis ketergantungan, konvensi penamaan berkas — semuanya **bentuk**. Kode status,
nama tabel, nama peran, angka ambang, nama sistem hilir — semuanya **menyebut dunia nyata**, dan
dunia nyata modul contoh bukan dunia nyata modul ini.
