> Modul  : Treaty In Adjustment · ronde TDA
> Dibuat : 2026-09-24
> Sifat  : sidang atas tujuh TDA tanpa nasib. **Setiap penutupan membawa kutipan beserta barisnya.**
>          Penutupan tanpa kutipan tidak sah (`00-LINGKUP.md` §3).

# Sidang ronde TDA — nasib tujuh TDA

| # | TDA | Keluaran | Oleh |
|---|---|---|---|
| **ST-1** | `TDA-01` | **diperbaiki** | `B-1` + `GRL-10` |
| **ST-2** | `TDA-12` | **diperbaiki** | `GRL-09` |
| **ST-3** | `TDA-15` | **diperbaiki** | ADR-0034 + `GRL-03` + `GRL-14` |
| **ST-4** | `TDA-09` | **KONFIRMASI** | ADR-0044 |
| **ST-5** | `TDA-07` | **diperbaiki** | `GRL-08` + INV-24 |
| **ST-6** | `TDA-11` | **diperbaiki** untuk model · sisa layar REKOMENDASI | `GRL-10` + `B-4` + INV-25 |
| **ST-7** | `TDA-13` | **satu residu — PERTANYAAN** | sebagian tertutup `GRL-12` + `GRL-13` |

---

## ST-1 — `TDA-01`: penjaga duplikat mati, tabrakan menimpa dan melapor berhasil → **diperbaiki**

**Kutipan, `GRILL-B/06-PUTUSAN.md` baris 48 — bahan to-spec `B-1`, yang menyebut TDA-01 di dalam
kolom "dari apa":**

> *"Menyimpan versi dengan **nomor urut** yang sudah dipakai di dalam kontrak yang sama **ditolak**,
> dan pesannya **menyebut nomor yang bentrok**"* — dari: *"tabrakan masuk cabang `UPDATE`, menimpa
> baris yang ada, lalu melapor berhasil (TDA-01)"*.

Dan sebab hulunya ikut hilang: `GRL-10` menetapkan **dasar = versi berlaku terakhir**, bukan baris
yang dipilih di picker, sehingga nomor tidak lagi dihitung dari baris pilihan.

**Sisa warisannya bukan nol, dan itu dinyatakan:** baris yang **sudah terlanjur tertimpa** di
produksi tidak dipulihkan oleh keputusan mana pun — melestarikan nomor tidak memulihkan isi yang
tertimpa (`GRL-09`). Diukur **`UA-14`** dan **`UA-21`**. Nasib sisa warisan: **ditunda**, pemilik
DBA.

---

## ST-2 — `TDA-12`: offset pengurai meleset satu → **diperbaiki**

**Kutipan, indeks, baris `GRL-09`:** label putusannya sendiri berbunyi

> *"pengenal warisan → PELESTARIAN; **bentuk turunan tanpa lebar tetap → PERUBAHAN (TDA-12)**"*.

Cacatnya sudah disebut **di dalam label putusan yang menutupnya**, dan tetap tercatat "terbuka" di
Lacak TDA selama dua ronde. Ini instans paling murni dari bentuk yang melahirkan ronde ini: **satu
keadaan diperbarui di satu tabel, tidak di tabel sebelahnya.**

---

## ST-3 — `TDA-15`: penggandaan pohon di dalam satu `JSONDATA` → **diperbaiki**

**Kutipan, ADR-0034:**

> *"Arsip JSON sistem lama **tidak punya jalur baca apa pun** di aplikasi hasil migrasi. Bukan
> 'dihindari', bukan 'hanya untuk kasus khusus' — **tidak ada kode yang membacanya**."*

Dan sekarang, untuk pertama kalinya, **keempat pohon punya tujuan bernama** — yang terakhir baru
didapat ronde C:

| Pohon lama | Ke mana | Oleh |
|---|---|---|
| nilai baru | `VERSI_KONTRAK` | `GRL-01` |
| `OLDDATA` | **tidak disimpan** — sisi lama dibaca lewat `ID_VERSI_KONTRAK_DASAR` | `GRL-03` butir 2, `GRL-10` |
| `ActualValue` | **tidak ada** — potret dibuang, masukan menjadi nilai versi | **`GRL-14`** |
| `ValueDifference` | `NILAI_SELISIH`, tabel tersendiri | `GRL-03` butir 3 |

> Penggandaan pohon bukan diperbaiki dengan aturan yang melarangnya. Ia hilang karena **tidak ada
> lagi dokumen yang dapat memuat empat pohon.**

---

## ST-4 — `TDA-09`: peran dari `pyWorkBasketList(2)` / `pyTelephone` / nama tersemat → **KONFIRMASI**

**Kutipan, ADR-0044 §Konsekuensi:**

> *"**Entitas peran dan penugasan bertanggal masuk gelombang 1.** Penugasan menyimpan sejak kapan
> sebuah peran melekat pada seseorang, bukan hanya keadaan hari ini… **Nama orang tidak pernah
> muncul di dalam aturan.**"*

Dan ADR-0044 §Konteks membaca sistem lama dengan tepat: *"Kita melihat sistem yang tidak pernah
punya cara menyatakan 'siapa', sehingga pertanyaan itu tidak pernah bisa diajukan."*

**Baris yang sedang berjalan saat peralihan** juga sudah punya jawaban — `GRILL-A/06-PUTUSAN.md`
baris 570: `Position` dan `StatusAkseptasi` adalah **satu-satunya** sumber untuk memetakannya.

**Sisa yang bukan pekerjaan teknis:** nama orang yang tertanam di aturan **yang hidup sejak 2019**
sudah menjadi butir daftar eskalasi induk. Tidak ditulis ulang di sini (`METODE` §6.6).

---

## ST-5 — `TDA-07`: dua kontrol layar mengosongkan status akseptasi kontrak → **diperbaiki**

**Kutipan, `GRILL-A/06-PUTUSAN.md` baris 563–565:**

> *"Kelima penanda lama bukan itu: `ViewState`, `IsEditData`, `RevisionState`, `Position`, dan
> `StatusAkseptasi` adalah **mekanisme**, bukan keputusan — tidak seorang pun pernah 'memutuskan
> `ViewState = 1`', ia akibat."*

Maka tidak ada `StatusAkseptasi` di sistem baru yang dapat dikosongkan sebuah kontrol layar.
Penggantinya berperilaku berlawanan: merevisi kontrak yang sudah disetujui **membuat versi baru
dalam keadaan `DRAFT`**, sementara versi yang disetujui tetap `DISETUJUI` dan **beku** — INV-24,
*"terminal berarti beku, bukan hanya berhenti berpindah"*, yang menjangkau entitas anaknya dengan
trigger `BEFORE INSERT OR UPDATE OR DELETE`.

> Cacat lama: satu kontrol **mengubah keadaan kontrak yang sudah disetujui**.
> Bentuk baru: keadaan itu **tidak dapat diubah oleh apa pun**, dan yang berubah adalah versi lain.

**Sisa:** wewenang membuka kunci — sudah ada di daftar eskalasi induk butir 1, dan `ESK-1` sudah
dinyatakan **tidak jadi butir baru**.

---

## ST-6 — `TDA-11`: picker menyatukan kontrak dan addendum → **diperbaiki untuk model**

**Kutipan, `GRILL-B/06-PUTUSAN.md` baris 123 dan 161** — keduanya menyebut TDA-11 di kolom "dari
apa":

> `GRL-10`: *"versi dasar = **versi berlaku terakhir pada saat versi baru dibuat**"* — **PERUBAHAN**,
> dari *"dasar = **baris mana pun yang dipilih di picker**, termasuk versi lama dan termasuk baris
> kontrak (§4.4 butir 2, TDA-11)"*.
>
> `B-4`: *"Versi yang ditunjuk `ID_VERSI_KONTRAK_DASAR` harus **versi berlaku terakhir**… dan **tidak
> boleh** berkeadaan `DITOLAK` atau `DIBATALKAN`"* — dari: *"dasar = baris mana pun yang dipilih di
> picker (TDA-11)"*.

**Kerusakan pokoknya hilang bukan dengan menyaring pickernya, melainkan dengan mencabut perannya:**
di sistem baru picker tidak lagi menentukan apa pun — dasar sebuah versi **diturunkan**. Dan
ketiadaan saringan keadaan ditutup dari arah lain oleh **INV-25**: paling banyak satu versi
tak-terminal per kontrak, sehingga kontrak yang sedang punya draf tidak dapat menerima draf kedua.

**Sisa yang tetap terbuka dan bukan pertanyaan grilling:** apa yang layar tampilkan dan bagaimana ia
menamai barisnya. Itu **REKOMENDASI TO-SPEC**, dan ia bertetangga dengan `R-H1`. Tidak ada
spesifikasi layar di mana pun (L-4 induk), jadi tidak ada yang dapat diputuskan di sini tanpa
mengarang bentuk antarmuka.

---

## ST-7 — `TDA-13`: jenis dan materialitas dipilih bebas di radio → **satu residu**

### Yang SUDAH tertutup

**Sisi materialitas: tertutup penuh.** `GRILL-B/06-PUTUSAN.md`, `GRL-12` butir 1 dan 4:

> *"**Untuk versi baru: material = versi itu memiliki sedikitnya satu baris `NILAI_SELISIH`.** Tidak
> ada atribut materialitas pada versi baru."*
>
> *"`TDA-10` dinyatakan diperbaiki. Materialitas turunan **tidak dapat berbohong**: ia bukan
> pernyataan yang perlu ditegakkan, melainkan pembacaan atas apa yang tersimpan."*

Radio materialitas **hilang bersama atributnya**. Tidak ada lagi yang dapat dipilih bebas.

**Sisi jenis: sebagian.** `GRL-13` menetapkan `JENIS_ADDENDUM` bernilai **dua** untuk versi baru, dan
`GRL-12` bahan to-spec `C-2` menegaskan kedua sumbu saling bebas:

> *"`JENIS_ADDENDUM` tidak menentukan materialitas, dan materialitas tidak menentukan jenis."*

### Yang TIDAK tertutup, dan itu residunya

Tidak satu pun putusan menyatakan **siapa yang menetapkan `JENIS_ADDENDUM`, dan apakah ia dapat
berubah sesudah versinya lahir.**

Di sistem lama jawabannya terbaca, dan ia bentuk yang tidak dapat dibawa:

| | Sistem lama |
|---|---|
| **disemai** | tombol — Revisi → `EDMState = 1`; Penyesuaian → `EDMState = 3` |
| **lalu** | radio picker menimpanya **tanpa syarat apa pun** (`NA-01`), antara 1 dan 2 |
| **kapan** | selama layarnya dapat disunting — termasuk sesudah addendum dibuat |
| **diperiksa di sisi simpan** | **nol** |

Maka jenis sebuah addendum di sistem lama adalah **nilai yang dapat berubah kapan saja tanpa jejak**,
dan `GRL-13` memberinya dua nilai tanpa menyatakan apakah sifat itu ikut.

**Ini bukan pertanyaan layar.** Ia menentukan apakah `JENIS_ADDENDUM` masuk lapisan beku versi,
apakah ia butuh baris di `JEJAK_PERUBAHAN` ketika berubah, dan apakah ia dapat berbeda antara saat
pengajuan dan saat persetujuan — tiga hal yang seluruhnya milik model.

**Dibawa ke frontier sebagai pertanyaan tunggal ronde TDA.**
