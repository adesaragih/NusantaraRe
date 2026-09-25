# Keputusan tanpa verifikasi — Treaty In Adjustment

**Tanggal:** 24 September 2026 · **Memutuskan:** pemilik proses migrasi NuRe
**Butir:** **4** · **Terverifikasi: NOL.**

> ## Namanya menyatakan sifatnya, dan itu disengaja
>
> Berkas ini **tidak** diberi nama yang menyamarkan bahwa isinya belum diuji. Setiap butir di
> dalamnya diputuskan **tanpa** bukti yang memutuskannya — yang ada hanya **dasar** dan **syarat
> pembalikan**.
>
> **Tidak satu pun butir di berkas ini boleh dikutip sebagai fakta tentang sistem lama.**

## Kenapa keempatnya diputuskan sekarang

Ronde D menyebut empat butir *"menahan model"*. Keempatnya **diuji ulang terhadap satu pertanyaan —
*apakah ia menentukan letak kolom?*** — dan **tidak satu pun lolos**.

| Butir | Yang sebenarnya diblokirnya |
|---|---|
| `DB-20` | kapan sebuah kolom **membeku** — bukan di mana ia berada |
| `DB-16a` | apakah sebuah **nilai enum** perlu ditambah — bukan tabel mana yang ada |
| `DB-16b` | apakah sebuah **kolom nullable** terpakai — bukan di tabel mana ia duduk |
| paket `REV` | **alasan dan daftar keadaan** di ADR — bukan bentuk data |

> **Keempatnya memblokir KEPASTIAN, bukan MODEL.** Mencampur keduanya membuat proyek tampak jauh
> lebih terhambat daripada keadaannya — kekeliruan yang sudah tercatat di modul induk.

**Ketiga `DB` tetap dikirim.** Yang berubah hanya kedudukannya: jawabannya **menyempitkan**, tidak
lagi memblokir.

---

## Dua koreksi terhadap rumusan yang diterima — dinyatakan, bukan diselaraskan diam-diam

### Koreksi 1 — nama entitasnya `DOKUMEN_ADDENDUM`, bukan `DOKUMEN_PENYESUAIAN`

Perintah yang diterima menyebut **`DOKUMEN_PENYESUAIAN`** pada `KTV-2`. Nama yang **terkunci**
adalah **`DOKUMEN_ADDENDUM`** — `GRL-19` butir 4, dengan alasannya tertulis: ia dibedakan dari
`DOKUMEN_KONTRAK` milik modul induk, yang dipakai untuk lampiran dan berstatus *"dirujuk, tidak
dimiliki"*.

**Nama yang dipakai di berkas ini: `DOKUMEN_ADDENDUM`.** Alasan tambahan, dan ia bukan kerapian:
`GRL-04` yang baru dibatalkan menyatakan **`PENYESUAIAN` tidak kembali**, dan penelusuran 19
properti yang menolaknya **tetap berdiri**. Memakai kata `PENYESUAIAN` pada entitas yang baru lahir
akan membuat pembaca berikutnya mengira putusan itu dicabut — padahal yang lahir benda **di atas**
versi, bukan benda **di bawah**nya.

> **Bila pemilik proses memang menghendaki `DOKUMEN_PENYESUAIAN`**, itu perubahan pada `GRL-19`
> butir 4 dan ditulis sebagai koreksi putusan, bukan sebagai pemilihan nama di berkas ini.

### Koreksi 2 — pasangan `KTV-2`/`KTV-3` dengan `DB-16a`/`DB-16b` tertukar

Perintah yang diterima memasangkan `KTV-2` (tanggal berlaku) dengan `DB-16a`, dan `KTV-3` (penanda
dikirim ke luar) dengan `DB-16b`. Bunyi keduanya, dikutip dari `GRILL-D/05-GERBANG.md` §5:

> **`DB-16a`.** *"Revisi internal tidak pernah disertai dokumen yang dikirim ke cedant."*
> **`DB-16b`.** *"Dokumen addendum selalu punya tanggal berlaku sendiri, yang dapat berbeda dari
> tanggal berlaku versi yang dipayunginya."*

**Pasangan yang benar: `KTV-2` ↔ `DB-16b` (tanggal berlaku), `KTV-3` ↔ `DB-16a` (dikirim ke luar).**
Dipakai apa adanya di bawah.

---

## KTV-1 — Materialitas memakai titik beku dan jejak yang sama dengan jenis

**Label: PERUBAHAN.** Berubah dari: *"materialitas hanya ditegakkan di layar, dan tidak ada jejak
perubahannya"* (`TDA-10`).

### Keputusan

`SIFAT_MATERIAL_ADDENDUM` **beku sejak `AJUKAN`**, dan setiap perubahannya selama `DRAFT`
meninggalkan **baris jejak** — persis aturan yang `GRL-18` bagian 2 dan 3 tetapkan untuk
`JENIS_ADDENDUM`.

### Dasar

1. **Keseragaman dengan `GRL-18`.** Dua sumbu yang dipilih di layar yang sama, sebelum perubahan
   dikerjakan, dengan titik beku yang berbeda akan menuntut dua aturan dan dua penjelasan — dan
   pembacanya tidak punya cara menebak mana yang berlaku untuk mana.
2. **Ongkos salahnya tidak setangkup.** Membawa jejak yang ternyata tidak diperlukan berarti sebuah
   tabel jejak dengan lebih sedikit baris — **murah**. Memasangnya belakangan berarti **jejak untuk
   masa sebelum pemasangan tidak ada**, dan tidak dapat diadakan kemudian.

### Syarat pembalikan

Bila bisnis menyatakan pilihannya **beku sejak lahir** — tidak dapat diubah bahkan selama `DRAFT` —
maka jejaknya **dipangkas**. Itu **satu perubahan aturan**, bukan perubahan skema: barisnya tetap
ada, hanya tidak pernah lebih dari satu.

**Yang menyempitkannya:** `DB-20`.

### Arah dampak bila salah

Salah ke arah "terlalu longgar": ada jejak yang tidak dibaca siapa pun. Salah ke arah sebaliknya —
tidak dipasang, lalu dibutuhkan: **tidak ada cara memulihkan jejak yang tidak pernah ditulis.**

---

## KTV-2 — `DOKUMEN_ADDENDUM` membawa tanggal berlaku sendiri, boleh kosong

**Label: BARU.** Entitasnya sendiri baru (`GRL-19`); kolom ini tidak punya padanan di sistem lama.

### Keputusan

`DOKUMEN_ADDENDUM` punya kolom **tanggal berlaku**, **boleh kosong**. Bila kosong, tanggal berlaku
yang dipakai adalah tanggal berlaku **versi** yang dipayunginya.

### Dasar

Kolom *nullable* yang tidak terpakai **murah** — ia satu kolom kosong. Menambahkannya **sesudah
migrasi berjalan** mahal: setiap baris warisan harus disentuh ulang, dan tidak ada sumber untuk
mengisinya (`TD-01` — sistem lama tidak merekam dokumen sama sekali).

### Syarat pembalikan

Bila **`DB-16b` dibantah** — dokumen **tidak** punya tanggal berlaku sendiri, ia selalu mengikuti
versinya — kolomnya **dicabut sebelum data masuk**. Mencabut kolom kosong tidak berongkos.

### Arah dampak bila salah

Salah "terlalu longgar": satu kolom kosong selamanya. Salah sebaliknya: dokumen yang berlaku sejak
tanggal berbeda dari versinya **tidak dapat dinyatakan**, dan yang terjadi adalah orang memaksakan
tanggalnya ke versi — merusak tanggal versi untuk menyelamatkan tanggal dokumen.

---

## KTV-3 — Tidak ada penanda "dikirim ke luar"; jenis tetap bernilai dua

**Label: PELESTARIAN.** Berubah dari: tidak ada — sistem lama pun tidak punya penanda itu.

### Keputusan

**Tidak ada kolom yang menyatakan sebuah versi atau dokumen "dikirim ke luar".** `JENIS_ADDENDUM`
tetap bernilai **dua** (`GRL-13`, `GRL-18`).

### Dasar

Memilih **tempat** sebuah penanda sebelum tahu **benda apa** yang ditaruh lebih mahal daripada
menambahkan satu nilai enum kelak. Dan sesudah `GRL-19`, pertanyaannya berubah bentuk: bila
"dikirim ke luar" ternyata bermakna, tempatnya **mungkin bukan enum jenis sama sekali** — mungkin ia
sifat **dokumen**, yang kini punya entitasnya sendiri.

> **Menunggu justru memperbaiki keputusannya**, bukan sekadar menundanya.

### Syarat pembalikan

Bila **`DB-16a` dibantah** — revisi internal **memang** kadang disertai dokumen yang dikirim ke
cedant — maka penandanya ditambahkan, dan **tempatnya diputuskan saat itu**: satu nilai enum pada
jenis, atau satu kolom pada `DOKUMEN_ADDENDUM`.

### Arah dampak bila salah

Salah: versi baru yang lahir sementara itu **kehilangan klasifikasinya**, sebanyak jarak sampai
`DB-16a` terjawab. Data lama tidak terpengaruh — `EDMState` warisan dibawa apa adanya (`GRL-13`).

---

## KTV-4 — Paket `REV-1` … `REV-6` diserahkan, dan tidak menahan to-spec

**Label: PERUBAHAN.** Berubah dari: prasyarat cabang K butir 2 berbunyi *"paket `REV` diserahkan
**dan ditanggapi** pemilik ADR"*.

### Keputusan

Prasyarat cabang K butir 2 berbunyi **"diserahkan"**. **Tanggapannya ditagih sebelum to-ticket**,
bukan sebelum to-spec.

### Dasar

Keenam REV diperiksa terhadap pertanyaan *"apakah ia menentukan letak kolom?"*:

| REV | ADR | Yang direvisinya | Menentukan letak kolom? |
|---|---|---|---|
| `REV-1` | ADR-0048 | koreksi **pemerian**, bukan penalaran | tidak |
| `REV-2` | ADR-0052 | penyimpangan yang berjalan **satu**, bukan dua | tidak |
| `REV-3` | ADR-0055 | empat koreksi pada **daftar keadaan** | tidak — keadaan sudah menjadi kolom lewat `GRL-08` |
| `REV-4` | ADR-0049 | sumber turunan materialitas berubah | **tidak lagi** — `GRL-20` sudah memutuskan materialitas menjadi **masukan** |
| `REV-5` | ADR-0049 | himpunan jenis menyusut menjadi dua | tidak — **sudah didahului `GRL-18`**, yang mengunci jenis sebagai masukan bernilai dua |
| `REV-6` | ADR-0037 | label PERUBAHAN + pro rata lama tidak pernah dapat selain 100% | tidak — `GRL-15` sudah memutuskan atributnya dibawa, mesinnya tidak dibangun |

**Nol dari enam** menentukan letak kolom. Seluruhnya merevisi **alasan** dan **daftar keadaan**.

### Syarat pembalikan

Bila pemilik ADR **menolak** salah satu REV **dengan akibat pada bentuk data**, bagian to-spec yang
bersandar padanya **dibuka kembali**. Yang paling mungkin: penolakan `REV-3` yang mengubah daftar
keadaan, sebab keadaan **adalah** kolom.

### Arah dampak bila salah

Salah: satu bagian to-spec ditulis ulang. Salah ke arah sebaliknya — menunggu tanggapan yang tidak
kunjung datang — **seluruh to-spec tertahan oleh revisi alasan**, dan itu ongkos yang jauh lebih
besar.

---

## Ringkasan kedudukan

| Butir | Kedudukan | Yang menyempitkannya | Kapan ditagih |
|---|---|---|---|
| `KTV-1` | diputuskan | `DB-20` | sebelum to-ticket |
| `KTV-2` | diputuskan | `DB-16b` | **sebelum data masuk** — kolomnya masih dapat dicabut |
| `KTV-3` | diputuskan | `DB-16a` | sebelum to-ticket |
| `KTV-4` | diputuskan | tanggapan pemilik ADR | **sebelum to-ticket** |

**Nol di antaranya terverifikasi.**
