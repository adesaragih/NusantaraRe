# 15: Hapus klaim — popup konfirmasi, kaskade tiga tingkat, dan baris work

**Status:** ready-for-agent

**Blocked by:** 14 (skema relasional klaim — PREFACTOR), 03 (peserta, adjustment, dan dokumen harus
ada agar dapat dihitung dan dihapus)

> Tiket ini menutup **AC 47** — yang sebelumnya tidak dirujuk tiket mana pun — dan **sisi perilaku**
> AC 48. Tiket 14 hanya membuat *constraint* `ON DELETE CASCADE`-nya; yang dilihat pengguna —
> peringatan berisi jumlah, dan pembatalan yang benar-benar tidak mengubah apa pun — belum ada
> pemiliknya.

## Hasil & nilai pengguna

Sebagai **`ReasLifeAdmin`**, saya dapat menghapus sebuah klaim **beserta seluruh isinya** —
peserta, seluruh baris adjustment, seluruh spreading dan spreading retro, dan seluruh dokumen —
tetapi **tidak sebelum diberi peringatan berisi jumlah baris yang akan ikut terhapus**, supaya saya
dapat membatalkan ketika angkanya tidak sesuai dugaan. *(User story 54 di spec)*

⚠️ **RALAT 2026-09-18** — kata **"polis, marketing"** **DICABUT** dari daftar di atas: kedua
tabelnya **dihapus**. Data polis dan marketing **dibaca hidup** dari tabel polis dan **bukan milik
klaim**, jadi menghapus klaim **tidak menyentuhnya**.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Hapus berkaskade tiga tingkat + baris `T_WORK_CLAIM`, **dalam satu transaksi** |
| `internal/services` | Hitung jumlah baris terdampak sebelum menghapus; orkestrasi |
| `internal/handlers` | Endpoint pratinjau dampak + endpoint hapus |
| `frontend/` | Popup konfirmasi Ya/Batal dengan rincian jumlah per jenis |

## Bentuk yang dihapus — spec §2b

⚠️ **RALAT 2026-09-18** `[keputusan work owner]` — **`T_CLAIM_POLICY` dan `T_CLAIM_MARKETING`
DIHAPUS**, bukan diganti nama, jadi keduanya **tidak ada lagi untuk ikut terhapus**. ⛔ Jangan
membuat `T_CLAIMLF_POLICY`/`T_CLAIMLF_MARKETING`. Data polis dan marketing **dibaca hidup** dari
tabel polis — menghapus klaim **tidak boleh menyentuhnya sama sekali**.

```
T_GENERAL_CLAIM                       ← yang dihapus pengguna
  └─ T_CLAIMLF_PREMIUMLIST_DETAIL  1:N   ikut
        ├─ T_CLAIMLF_ADJUSTMENT     1:N   ikut  ⬅ CUCU
        │     └─ T_CLAIMLF_ADJUSTMENT_SPREADING        1:N  ikut  ⬅ CICIT
        │            └─ T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO  1:N  ikut  ⬅ CICIT-CUCU
        └─ DOCUMENT_CLAIM         1:N   ikut  ⬅ CUCU

T_WORK_CLAIM                      ← baris work klaim life itu ikut terhapus
```

⚠️ **RALAT 2026-09-18 — klaim "tiga tingkat" di tiket ini SUDAH SALAH SEJAK 2026-09-16, bukan
karena ralat hari ini.** Dua tabel spreading ditemukan pada verifikasi 2026-09-16 dan **tidak pernah
masuk** ke pohon di tiket ini. Kaskade sebenarnya menyentuh **lima tingkat** — klaim → peserta →
adjustment → spreading → spreading retro — persis seperti **AC 60 spec** dan tiket **14**. Uji wajib
memeriksa **sampai cicit-cucu** (`T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`), bukan hanya dua tingkat
cucu. Pohon yang berlaku: `spec.md` §2b RALAT D.

## ADR terkait

**ADR-0007** (jejak audit), **ADR-0015** (kegagalan ditangani eksplisit, tidak ditelan).

## Acceptance criteria

- [ ] ⚠️ Menghapus klaim **mengkaskade** ke peserta, **seluruh baris adjustment**, **seluruh
      spreading**, **seluruh spreading retro**, dan **seluruh dokumen**. **REVISI 2026-09-18:** kata
      **"polis, marketing"** **DICABUT** — kedua tabelnya dihapus. Test wajib memeriksa **setiap
      tingkat** secara terpisah, **sampai cicit-cucu**. *(AC 48 spec; penyimpangan sadar 8)*
- [ ] ⚠️ Menghapus klaim **tidak menyentuh tabel polis maupun marketing** — data itu **dibaca
      hidup** dari `T_PREMIUM_LIST` dkk dan **bukan milik klaim**. Test yang menemukan penghapusan
      menyentuh tabel polis **gagal**. *(REVISI 2026-09-18; `[keputusan work owner]`)*
- [ ] ⚠️ Penghapusan **didahului popup konfirmasi Ya/Batal** yang menyebut **jumlah baris tiap
      jenis** yang akan ikut terhapus. *(AC 48 spec)*
- [ ] ⚠️ Memilih **Batal** **tidak mengubah apa pun** — tidak ada baris terhapus, tidak ada status
      berubah, tidak ada jejak audit penghapusan. *(AC 48 spec)*
- [ ] Jumlah yang ditampilkan popup **sama persis** dengan jumlah yang benar-benar terhapus —
      dihitung dari data, bukan dari perkiraan.
- [ ] ⚠️ Menghapus klaim life **menghapus juga baris `T_WORK_CLAIM`**-nya; tidak ada keadaan tangga
      yang tertinggal tanpa klaim. *(AC 47 spec; penyimpangan sadar 6)*
- [ ] Seluruh penghapusan berjalan dalam **satu transaksi**: kegagalan di tingkat mana pun
      **membatalkan seluruhnya**, dan klaim tetap utuh. *(AC 49 spec)*
- [ ] Penghapusan mencatat **jejak audit** — siapa, kapan, dan berapa baris tiap jenis.
      *(**ADR-0007**)*
- [ ] Penghapusan yang gagal menghasilkan kegagalan **terang-terangan**, bukan sebagian terhapus
      diam-diam. *(**ADR-0015**)*
- [ ] Menghapus klaim **tidak menyentuh** `M_LIFE_PREMIUM_DETAIL` — ia hanya dibaca sebagai sumber
      snapshot peserta.
- [ ] ⚠️ **Perlakuan `OS_AKSEPTASI_KLAIM_LIFE` saat klaim dihapus** ditetapkan eksplisit. Ia **tetap
      ditulis** saat klaim disimpan (koreksi 2026-09-16, AC 32 spec), sehingga menghapus klaim
      menimbulkan pertanyaan: barisnya ikut dihapus, atau ditinggal karena hilir sudah membacanya?
      `[terbuka]` — **keputusan work owner**, **jangan tebak**. Tiket ini **tidak dinyatakan selesai**
      sebelum jawabannya ada.
- [ ] Menghapus klaim **tidak menghapus** peserta di premium list sumbernya — yang terhapus hanya
      **snapshot** milik klaim itu.

## Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka — Non-Life]` **tidak memblokir:** relasi & cascade formal `T_WORK_CLAIM` ditetapkan saat
konteks Non-Life digarap. Untuk Claim Life yang mengikat hanya AC di atas: klaim life dihapus →
baris work-nya ikut.

## Catatan

⚠️ **Mengapa kaskade tidak cukup diuji lewat constraint basis data.** Tiket 14 memasang
`ON DELETE CASCADE`; itu menjamin *baris* hilang, bukan bahwa **angka di popup benar** dan bahwa
**Batal benar-benar membatalkan**. Keduanya perilaku layanan dan layar — dan keduanya tempat bug
biasanya muncul.

⚠️ **Pola yang diikuti.** Kaskade + popup konfirmasi sudah ditetapkan di Master Contract Retro Life
(tiket 09) dan Treaty Contract Out (tiket 10). Bedanya di sini: **tiga tingkat**, dan ada **satu
tabel di luar pohon** (`T_WORK_CLAIM`) yang ikut.

## Seam & perintah verifikasi

**Seam: API HTTP Claim — Life** terhadap **skema uji Oracle nyata** — kaskade tiga tingkat dan
keutuhan setelah Batal **hanya berperilaku benar pada basis data sungguhan**; memalsukannya berarti
tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```
