# 03: Sisi inward (kolom pada satu tabel) + simpan atomik produk + anak

**Status:** selesai (01-10-2026) — paket 4 `b179aa2`, layar paket 10 (`1ada8d7`); uji `db` ditulis dan MELEWATI di sesi implementasi (tanpa `ORACLE_DSN`; POOLDATA/DEV bukan sasaran)

**Blocked by:** 02 (sisi umum harus ada — keduanya berbagi identitas yang sama)

## Hasil & nilai pengguna

Sebagai **underwriter**, saya menetapkan **syarat penerimaan** produk — rentang usia, batas uang
pertanggungan, batas retensi ceding, share Nusantara Re, brokerage, masa berlaku; dan sebagai
**organisasi**, saya ingin kedua sisi produk tersimpan **bersama atau tidak sama sekali**, sehingga
tidak pernah ada produk yang tersimpan separuh. *(User story 5–9 di spec)*

⚠️ Inilah tiket dengan penyimpangan paling dalam di konteks ini.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas produk sisi inward — 40 field |
| `internal/repository` | **Satu transaksi** menulis baris `product_life` beserta tabel anak |
| `internal/services` | Orkestrasi simpan atomik; rollback menyeluruh |
| `internal/handlers` | Endpoint simpan produk utuh |
| `frontend/` | Layar produk — bagian syarat underwriting |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SetProductNameInward` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SETPRODUCTNAMEINWARD` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SetProductNameInward.xml` (125.938 byte) | **40 field sisi inward** |
| `SaveInwardProductName_Act` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEINWARDPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SaveInwardProductName_Act.xml` (64.796 byte, 6 langkah) | orkestrator simpan inward |
| `SaveProductNameLIfe` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!SAVEPRODUCTNAMELIFE` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/SaveProductNameLIfe.xml` | ⚠️ **rujukan — tidak dimigrasikan** |
| `SaveProductNameInwardLIfe` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!SAVEPRODUCTNAMEINWARDLIFE` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/SaveProductNameInwardLIfe.xml` | ⚠️ **rujukan — tidak dimigrasikan** |

`[fakta bisnis — work owner]` **SATU produk.** Existing: dua tabel (1:1 by `ID`). Skema baru:
**digabung jadi satu tabel `product_life`** (`[keputusan work owner]`).

⚠️ **Penyimpangan sadar 1 — satu transaksi atomik, procedure lama dibuang.**
`[keputusan work owner]`

`[data DBA]` Kedua procedure penulis (`POOLDATA.PEGA_M_PRODUCT_LIFE`,
`POOLDATA.PEGA_M_PRODUCT_INWARD_LIFE`) adalah **upsert JSON** yang **`COMMIT` sendiri**. Selama
keduanya dipakai, "tersimpan bersama atau tidak sama sekali" **mustahil**.

⚠️ **Penggabungan tabel induk `[keputusan work owner]`.** Sisi umum + sisi inward (dulu 2 tabel
existing, 1:1 by `ID`) **digabung jadi SATU tabel `product_life`** (tiket 01). Maka "sisi inward" =
sekumpulan kolom pada baris produk yang sama, bukan tabel terpisah.

**Keputusan:** sistem baru **tidak memanggil dan tidak memigrasikan** kedua procedure. Go menulis
**satu baris `product_life`** (memuat sisi umum + inward) **beserta seluruh baris tabel anak** dalam
**satu transaksi**; kegagalan di mana pun **membatalkan seluruhnya**.

⚠️ Ini **berbeda** dari pola "panggil procedure apa adanya" yang berlaku di lima konteks Life
sebelumnya (**ADR-0006**). Alasannya tunggal: **procedure itu menghalangi atomik yang diminta.**
Efek samping baik penggabungan: atomik induk jadi **trivial** (satu baris), bukan lintas dua tabel.

## ADR terkait

**ADR-0003** (uang non-float — batas pertanggungan, retensi, limit), **ADR-0006** (identitas dari
sequence; ⚠️ pengecualian sadar pada pemanggilan procedure), **ADR-0007** (jejak audit),
**ADR-0015** (kegagalan ditangani eksplisit, tidak ditelan).

## Acceptance criteria

- [ ] Sisi inward dapat diisi dengan **keempat puluh field**, termasuk rentang usia, batas uang
      pertanggungan, batas retensi ceding, share Nusantara Re, brokerage, dan masa berlaku.
      *(User story 5–6; AC 41 spec)*
- [ ] ⚠️ Baris `product_life` (sisi umum + inward) **beserta seluruh baris tabel anak** ditulis
      dalam **satu transaksi**: kegagalan pada tabel anak mana pun **membatalkan seluruhnya**.
      Dibuktikan dengan menyuntikkan kegagalan pada satu tabel anak lalu memastikan **tidak ada baris**
      `product_life`. *(AC 6, 7 spec; `[keputusan work owner]` — penyimpangan sadar 1)*
- [ ] ⚠️ Sisi umum & sisi inward = **satu baris** `product_life` (dulu dua tabel existing digabung);
      field yang muncul di kedua sisi jadi **satu kolom**. *(AC 8 spec; `[keputusan work owner]`)*
- [ ] ⚠️ **Procedure JSON lama tidak dipanggil.** Test yang menemukan pemanggilan
      `PEGA_M_PRODUCT_LIFE` atau `PEGA_M_PRODUCT_INWARD_LIFE` **gagal**. *(AC 9 spec;
      `[keputusan work owner]`)*
- [ ] Nilai uang (`*LIMIT*`, `*SUMINSURED*`, `*SUMREASURED*`) diperlakukan sebagai **desimal presisi
      arbitrer**; **tidak** melewati `float`. *(AC 28 spec; **ADR-0003**)*
- [ ] Persentase (share, brokerage, retensi) juga desimal — **tidak** dibulatkan ke bilangan bulat.
      *(AC 30 spec)*
- [ ] Tanggal (`BEGIN`, `MATURE`, `STNC`, `BIRTHDAY`) tersimpan sebagai **tanggal**, bukan teks.
      *(AC 31 spec)*
- [ ] Batas bawah yang **lebih besar** dari batas atas — usia maupun uang pertanggungan — **ditolak**,
      dengan pesan yang menyebut field-nya. *(AC 32 spec)*
- [ ] Menyimpan produk yang sudah ada memperbarui **kedua** sisi, bukan menambah baris baru di salah
      satunya. *(AC 4 spec)*
- [ ] ⚠️ Kolom pemegang polis bernama **`POLICYHOLDER`**, bukan `POLICYHODER`. *(AC 51 spec)*

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Mengapa konteks ini menyimpang dari lima konteks Life sebelumnya.** Di Claim Life, Komite,
PremiumList Life, Endorsement Life, dan Master Contract Retro Life keputusannya selalu **"panggil
procedure apa adanya"**. Di sini tidak — dan D1 (atomik) menarik D2 (relasional) sebagai akibat:
bila Go menulis sendiri, menulis **JSON** tidak lagi masuk akal. **Satu keputusan, dua akibat.**

⚠️ **Rujukan yang tidak dipakai** `[data DBA]`: `StsSave` 100/99, `ErrMsg` (terisi juga saat
sukses), `IDPegaOut`. Ketiganya milik procedure yang dibuang; dicatat sebagai jejak, bukan perilaku.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — atomik lintas satu tabel induk
`product_life` dan lima tabel anak **hanya berperilaku benar pada basis data sungguhan**;
memalsukannya berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Ralat bertanggal 01-10-2026 — sesi implementasi (paket 4)

> Sumber: `../RALAT-DEV-30-09-2026.md` (P1, P4, P6, R11, R14, R17) dan `../PARITAS-LAYAR-DAN-AKSI.md` §3.2, §7. Kalimat di atas **tidak dihapus**.

| Kalimat lama | Ralat |
| --- | --- |
| *"**digabung jadi satu tabel `product_life`**"*; *"Go menulis **satu baris `product_life`**"* | **P1/P4**: dua tabel lama tetap — `M_PRODUCT_LIFE` dan `M_PRODUCTINWARD_LIFE` — ditulis dalam **satu transaksi**; gagal sisi inward membatalkan sisi umum (uji `TestSimpanAtomikDuaTabel`, tiruan dan `db`) |
| Judul *"Sisi inward (kolom pada satu tabel)"*; rule sumber `SaveInwardProductName_Act` | **R11**: harness/section `InwardProductName` tak terjangkau (`Inward` b75322 `1=2`); medan inward diisi di `InboxProductName` dan disimpan `SaveProductName_Act` 13–16 — satu jalur simpan |
| *"Sisi inward dapat diisi dengan **keempat puluh field**"* | 40 kunci view `PRODUCTINWARD_LIFE` ditulis (yang tanpa medan form dari JSON lama, `""` untuk baris baru), ditambah medan form di luar view (`EXPIRYAGE`, `PREMIUMFACTOR`, `AnnuityInterest`, `PremiumRefundFactor`, `CURRENCYID`) — Pega menulisnya lewat `@GetPageJSONString()` tetapi tidak memuatnya kembali; di sini dimuat dari `JSONDATA` sendiri (P2) |
| *"Tanggal (`BEGIN`, `MATURE`, `STNC`) tersimpan sebagai **tanggal**, bukan teks"* | **P1**: di `JSONDATA` teks `dd/MM/yyyy` seperti Pega (`BrowseReinstypeOR_SQL` b84); kolom datar `BEGIN_DATE` bertipe `DATE` (OQ-MPNL-08). Masukan diperiksa sebagai tanggal |
| *"Menyimpan produk yang sudah ada memperbarui **kedua** sisi"* | tetap: baris inward lama diperbarui menurut **ID-nya sendiri** (data lama ber-ID sequence inward, bertaut `PRODUCTID`); produk tanpa baris inward mendapat baris ber-ID produk. ID inward produk baru = ID produk (**R14**, OQ-MPNL-02) |
| *"Kolom pemegang polis bernama **`POLICYHOLDER`**"* | kunci JSON tetap `POLICYHODER` (dibaca view); nilainya dari pemilih `Policy Holder` (master `CLIENT`) dan disalin ke JSON umum |
| *"Batas bawah yang **lebih besar** dari batas atas … **ditolak**"* | ditegakkan (**R17**): `Minimum Age (Years)` ≤ `Maximum Age (Years)`, `Min Sum Insured` ≤ `Max Sum Insured` |
