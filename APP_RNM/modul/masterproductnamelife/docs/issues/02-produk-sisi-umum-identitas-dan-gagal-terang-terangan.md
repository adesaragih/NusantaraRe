# 02: Produk sisi umum — CRUD, identitas dari sequence, gagal terang-terangan

**Status:** selesai (01-10-2026) — paket 1 `4ce0771` (baca) + paket 3 `2f39341` (tulis), layar paket 10 (`1ada8d7`); uji `db` ditulis dan MELEWATI di sesi implementasi (tanpa `ORACLE_DSN`; POOLDATA/DEV bukan sasaran)

**Blocked by:** 01 (skema relasional — bentuk barunya harus ada lebih dulu)

## Hasil & nilai pengguna

Sebagai **admin master produk life**, saya dapat membuat, mengubah, dan membaca **definisi produk**
tanpa menentukan nomor identitasnya sendiri — dan bila penyimpanan gagal, saya **diberi tahu**,
bukan dibiarkan mengira produk sudah tersimpan. *(User story 1–4, 10 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas produk sisi umum |
| `internal/repository` | Tulis/baca `product_life`; identitas dari sequence |
| `internal/services` | Orkestrasi buat/ubah; pemetaan kegagalan menjadi galat domain |
| `internal/handlers` | Endpoint CRUD produk |
| `frontend/` | Layar produk — bagian umum |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InwardProductName` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `INWARDPRODUCTNAME` / `RULE-HTML-HARNESS` | `Master Product Name Life/Harness/InwardProductName.xml` (540.636 byte) | **titik masuk, satu-satunya Harness** |
| `SetProductName` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SETPRODUCTNAME` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SetProductName.xml` | **20 field sisi umum** |
| `SaveProductName_Act` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SaveProductName_Act.xml` (165.682 byte, 16 langkah) | orkestrator simpan; set `CREATEOP`/`UPDATEOP` |
| `NewProductLife` | `ASM-FW-GISFW-…` / `NEWPRODUCTLIFE` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/NewProductLife.xml` | produk baru |

`[terverifikasi]` **Bukan proses berjenjang** — nol rujukan `StatusAkseptasi`, tanpa
`Akseptasi_DT`, tanpa tombol Submit/Decline, **nol rule `Flow`**. Ia editor master murni.

`[data DBA]` **Format identitas ditiru**: **`'1' + lpad(sequence, 5, '0')`** — lima digit, misalnya
`100001`. Sumbernya `M_PRODUCT_LIFE_SEQ`.

`[data DBA]` Kontrak keluaran procedure lama — `StsSave` **100 = sukses / 99 = gagal**, `ErrMsg`
teks yang **terisi juga saat sukses**, `IDPegaOut` identitas final — **tidak dipakai**, karena
procedure-nya dibuang (tiket 03). Dicatat sebagai **rujukan**, bukan perilaku.

## ADR terkait

**ADR-0006** (identitas lewat sequence basis data — aplikasi tidak menyusunnya), **ADR-0007**
(jejak audit), **ADR-0003** (uang non-float pada atribut produk).

## Acceptance criteria

- [ ] Produk dapat dibuat dengan nama, tipe, grup, dan atribut sisi umum lainnya. *(AC 1 spec)*
- [ ] **Aplikasi tidak menetapkan identitas produk** — identitas dibuat basis data lewat sequence.
      Test yang menemukan pembentukan identitas di sisi aplikasi **gagal**. *(AC 2 spec; **ADR-0006**)*
- [ ] ⚠️ Identitas berbentuk **`'1'` + lima digit** (misalnya `100001`) — format lama **ditiru**.
      *(AC 3 spec; `[keputusan work owner]` — penyimpangan sadar 3)*
- [ ] Menyimpan produk yang **sudah ada** memperbarui baris itu; **tidak** membuat duplikat.
      *(AC 4 spec)*
- [ ] Produk dapat dibaca kembali utuh lewat API. *(AC 5 spec)*
- [ ] ⚠️ **Gagal terang-terangan**: penyimpanan yang gagal **ditampilkan** dan **tidak pernah** tampak
      berhasil. *(AC 10 spec; `[keputusan work owner]` — penyimpangan sadar 3)*
- [ ] Penyimpanan mencatat **siapa** pembuat dan **siapa** pengubah terakhir. *(**ADR-0007**)*
- [ ] Nilai uang pada atribut produk diperlakukan sebagai **desimal presisi arbitrer**; **tidak**
      melewati `float` di lapisan mana pun maupun di JSON API. *(AC 28 spec; **ADR-0003**)*
- [ ] Nilai uang yang ditulis dan dibaca kembali **identik** — tidak ada pembulatan diam.
      *(AC 29 spec)*
- [ ] ⚠️ Tidak ada **`PoductName`** di kode baru — hanya `ProductName`. *(AC 50 spec;
      `[keputusan work owner]` — penyimpangan sadar 5)*

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Salah ketik yang mengenai langkah simpan.** `[terverifikasi]` Ejaan **`PoductName`** (tanpa `r`)
muncul **20 kali di dua berkas saja** — keduanya activity simpan — dan **tidak pernah dideklarasikan
sebagai halaman**, sementara `ProductName` muncul **714 kali**. Yang bersalah ketik justru menjaga
**langkah yang memanggil penyimpanan**.

⚠️ Membuangnya **mengubah perilaku**: guard yang selama ini mungkin tak pernah menyala akan **mulai
menolak** penyimpanan yang dahulu lolos. Itu disadari, bukan diselundupkan. Validasi lengkapnya ada
di tiket **05**.

⚠️ `ErrMsg` procedure lama **terisi juga saat sukses** `[data DBA]` — jangan jadikan keberadaan teks
sebagai penanda kegagalan. Di sistem baru, kegagalan ditentukan hasil transaksi Go sendiri.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — pembuatan identitas lewat sequence tidak
dapat difake dengan jujur.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Ralat bertanggal 01-10-2026 — sesi implementasi (paket 3)

> Sumber: `../RALAT-DEV-30-09-2026.md` (P1, P6, R7, R8, R14) dan `../PARITAS-LAYAR-DAN-AKSI.md` §3, §7. Kalimat di atas **tidak dihapus**.

| Kalimat lama | Ralat |
| --- | --- |
| *"Area codebase … Tulis/baca `product_life`"* | **P1**: tulis/baca `M_PRODUCT_LIFE.JSONDATA` berkunci Pega + kolom datar `RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE` (pernyataan yang sama, satu transaksi). Nol tabel baru |
| *"⚠️ Membuangnya **mengubah perilaku**: guard yang selama ini mungkin tak pernah menyala akan **mulai menolak**"* | **R7**: prakondisi ber-`PoductName` di `SaveProductName_Act` 8–11, 13–16 ber-**PRE=false** — tidak pernah dievaluasi; tidak ada perilaku yang berubah. Kode baru tetap tanpa `PoductName` |
| Identitas *"'1' + lima digit"* | tetap (**P6**): `'1' ‖ LPAD(M_PRODUCT_LIFE_SEQ.NEXTVAL, 5, '0')`; nomor > 5 digit dan ID yang sudah dipakai salah satu tabel **ditolak terang** (500 berkalimat), tidak dipotong/digandakan |
| (tidak disebut) medan `Product Code` b3894 dapat disunting dan terikat `ProductName.ID` | **R14**: baca-saja; `POST` yang membawa `id` ditolak 400 (ADR-0006) |
| *"Penyimpanan mencatat **siapa** pembuat dan **siapa** pengubah terakhir"* | `CREATEOP` = pelaku untuk produk baru — termasuk salinan `Copy`, yang di Pega mewarisi pembuat asal (OQ-MPNL-13); tetap untuk ubah. `UPDATEOP` = pelaku setiap simpan (`SaveProductName_Act` 1 b361, 6 b1370) |
| *"Endpoint CRUD produk"* | `POST /api/master-product-name-life/produk` (baru), `PUT /produk/{id}` (ubah). **Nol** rute hapus: korpus tidak punya tombol/aktivitas hapus produk |

---

## Ralat bertanggal 01-10-2026 — lanjutan 1 (katalog DEV)

| Kalimat lama | Ralat |
| --- | --- |
| ralat paket 3 *"+ kolom datar `RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE`"* | kolom datar hanya `RIRISKID`, `RIRISK` (DEV `ALL_TAB_COLUMNS`); `PRODUCTNAME` dan `BEGIN_DATE` tidak ada dan tidak ditulis — `6fd539c`. Tanggal tetap teks `dd/MM/yyyy` di `JSONDATA` |
