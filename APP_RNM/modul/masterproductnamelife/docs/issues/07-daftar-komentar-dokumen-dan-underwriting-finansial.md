# 07: Komentar, dokumen klaim, dan underwriting finansial

**Status:** ready-for-agent

**Blocked by:** 03 (baris anak ikut dalam transaksi atomik produk)

## Hasil & nilai pengguna

Sebagai **admin master**, saya mencatat **dokumen klaim** yang disyaratkan produk dan meninggalkan
**komentar** beserta tanggal dan nama penulisnya; dan sebagai **underwriter** saya mencatat
**underwriting finansial** bila produk memerlukannya — sehingga riwayat pertimbangan dan syarat
pelengkap tidak hilang. *(User story 20–23 di spec)*

> ✅ **Blocker "tiga daftar belum berbentuk" DISELESAIKAN dari korpus** (bukan ditunda). Sensus
> Activity menunjukkan hanya **satu** yang benar-benar tabel anak nyata:
> - `[terverifikasi]` **`LIENCLAUSE` = field skalar**, bukan list — ia `ProductNameInward.LIENCLAUSE`
>   (satu dari 40 field inward, tiket 03). **Bukan tabel anak.** ❌ `product_life_lien` dicoret.
> - `[terverifikasi]` **`OutwardList` = dead code** — hanya diisi `GetReinsTypeOR_Life` (jalur OR,
>   mati oleh gerbang `1==2`). **Bukan tabel.** ❌ `product_life_outward` dicoret.
> - `[terverifikasi]` **`FinancialUnderwritingList` = tabel anak nyata**, kolom terbaca dari
>   `CopyFinancialWriting`: `MinInsured`, `MaxInsured`, `Employee`, `Non_Employee`.
>   ✅ `product_life_fin_uw` — kolom pasti, tidak ditebak.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Lima jenis baris anak |
| `internal/repository` | Tulis/baca kelima tabel anak **di dalam transaksi produk** |
| `internal/services` | Aturan per jenis baris |
| `internal/handlers` | Endpoint kelima daftar |
| `frontend/` | Lima grid di layar produk |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `AddCommentList_Act` | `ASM-FW-GISFW-…` / `ADDCOMMENTLIST_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/AddCommentList_Act.xml` | tambah komentar — dipanggil `SaveProductName_Act` step 7 |
| `ConvertHistoryDate` | `ASM-FW-GISFW-…` / `CONVERTHISTORYDATE` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/ConvertHistoryDate.xml` | format tanggal riwayat |
| `CopyFinancialWriting` | `ASM-FW-GISFW-…` / `COPYFINANCIALWRITING` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/CopyFinancialWriting.xml` | salin underwriting finansial |

`[terverifikasi]` Tabel anak nyata modul ini (kolom dari Activity, bukan tebakan):

| Tabel anak | Kolom baris | Sumber |
| --- | --- | --- |
| `product_life_comment` | `Date`, `OperatorName`, `Suggest` (isi saran) | `AddCommentList_Act` / contoh JSON |
| `product_life_document_claim` | `Document` (+ audit `pxCreate*`) | contoh JSON |
| `product_life_plan` | (di tiket 06) | — |
| `product_life_uw_limit` | (di tiket 06) | — |
| `product_life_fin_uw` | `MinInsured`, `MaxInsured`, `Employee`, `Non_Employee` | `[terverifikasi]` `CopyFinancialWriting` |

**Bukan tabel (dicoret):** `product_life_lien` (LIENCLAUSE = field skalar di `productinward_life`),
`product_life_outward` (OutwardList = dead code jalur OR).

## ADR terkait

**ADR-0007** (jejak audit — komentar adalah riwayat pertimbangan), **ADR-0003** (bila daftar yang
belum berbentuk ternyata memuat nilai uang).

## Acceptance criteria

- [ ] Produk dapat memuat **dokumen klaim** yang disyaratkan. *(AC 24 spec)*
- [ ] Produk dapat memuat **komentar**, tersimpan beserta **tanggal dan nama penulisnya**.
      *(AC 24–25 spec)*
- [ ] Produk dapat memuat **underwriting finansial** (`product_life_fin_uw`) dengan kolom
      `MinInsured`, `MaxInsured`, `Employee`, `Non_Employee`. *(AC 24 spec; `[terverifikasi]`)*
- [ ] `LIENCLAUSE` diperlakukan sebagai **field skalar** di sisi inward (tiket 03), **bukan** tabel.
- [ ] **Tidak ada** tabel `product_life_outward` — `OutwardList` milik jalur OR yang mati.
- [ ] Baris dapat **ditambah dan dihapus** tanpa menyentuh produk induk — kecuali ketika penyimpanan
      dilakukan sebagai satu kesatuan. *(AC 26 spec)*
- [ ] Daftar yang **kosong** tersimpan sebagai daftar kosong, **bukan** kegagalan. *(AC 27 spec)*
- [ ] ⚠️ Kegagalan pada baris mana pun **membatalkan seluruh simpan** (baris `product_life` +
      seluruh anak) — baris anak ikut dalam transaksi produk. *(AC 7 spec; `[keputusan work owner]`)*
- [ ] ⚠️ Komentar **tidak dapat diubah** setelah tersimpan — ia riwayat, bukan field.
      `[keputusan work owner]` (dikonfirmasi: immutable).
- [ ] Urutan baris dalam tiap daftar **dipertahankan** saat dibaca kembali.

## Blocker

**Tidak ada.** ✅ Bentuk seluruh tabel anak **terselesaikan dari korpus** — tidak ada yang ditunda
ke implementasi.

## Catatan

⚠️ **Aturan pemilihan kolom tetap berlaku** `[keputusan work owner]`: kolom = field yang **di-SET di
Activity**. Field yang hanya muncul di Section dengan **visibilitas mati** (`1=2` / `1==2` /
`never`) dan tidak di-set Activity **jangan diambil**. Modul ini `[terverifikasi]` memang memuat
banyak gerbang `1==2` di Harness-nya.

⚠️ **Komentar melekat pada jalur simpan.** `[terverifikasi]` `SaveProductName_Act` **step 7**
memanggil `AddCommentList_Act` — jadi komentar terbentuk saat menyimpan produk, bukan lewat jalur
terpisah. Pertahankan keterkaitan itu.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```
