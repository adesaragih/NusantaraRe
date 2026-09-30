# 03: Popup "lihat polis lama" — satu tampilan per jenis transaksi

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 02 (case endorsement harus sudah memegang rujukan polis lama)

## Hasil & nilai pengguna

Sebagai **inputor Life**, sebelum mengubah apa pun saya ingin **membuka polis lama dalam tampilan
tersendiri** dan melihat kolom yang relevan bagi jenis transaksinya, sehingga saya yakin sedang
meng-endorse polis yang benar. *(User story 10–11 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Baca polis NB beserta detailnya berkunci ID pega lama / nomor polis lama |
| `internal/services` | Pemilihan tampilan menurut `Type` |
| `internal/handlers` | Endpoint "lihat polis lama" |
| `frontend/` | Popup polis lama — tampilan induk + tiga varian |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `ShowLifePremiumSummary_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SHOWLIFEPREMIUMSUMMARY_EDM` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/ShowLifePremiumSummary_EDM.xml` (2.194.694 byte) | **pemicu, memilih varian per `Type`** |
| `ViewOldPolicy_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWOLDPOLICY_EDM` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/ViewOldPolicy_EDM.xml` | **tampilan induk — melayani `QR`** |
| `ViewOldPolicy_EDM_QP` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWOLDPOLICY_EDM_QP` | `Endorsement Life/Section/ViewOldPolicy_EDM_QP.xml` | varian `QP` |
| `ViewOldPolicy_EDM_TP` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWOLDPOLICY_EDM_TP` | `Endorsement Life/Section/ViewOldPolicy_EDM_TP.xml` | varian `TP` |
| `ViewOldPolicy_EDM_TR` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWOLDPOLICY_EDM_TR` | `Endorsement Life/Section/ViewOldPolicy_EDM_TR.xml` | varian `TR` |
| `GetOldDetail_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GETOLDDETAIL_EDM` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/GetOldDetail_EDM.xml` (53.517 byte) | ⚠️ pemuat data — **rule tidak dimigrasikan, perilakunya wajib ada** |

`[terverifikasi]` **Pemasangan varian per `Type`** — di `ShowLifePremiumSummary_EDM`, tiap pasang
tombol diikuti `<pyCondition>` atas `.Type`:

| Harness (baris) | `<pyCondition>` (baris) |
| --- | --- |
| `ViewOldPolicy_EDM` (65043, 65282) | **`.Type=='QR'`** (65394) |
| `ViewOldPolicy_EDM_QP` (65575, 65816) | `.Type=='QP'` (65952) |
| `ViewOldPolicy_EDM_TR` (66136, 66374) | `.Type=='TR'` (66510) |
| `ViewOldPolicy_EDM_TP` (66694, 66932) | `.Type=='TP'` (67068) |

`[terverifikasi]` Tiap tombol adalah **pasangan aksi**: `Run Activity` → `GetOldDetail_EDM`
(`pyActivityClass` = `ASM-FW-GISFW-Work-EndorsementLife`), lalu `showHarness` dengan
`pyTarget = popup`.

## ADR terkait

**ADR-0001** (batas konteks — membaca polis NB adalah pembacaan hulu), **ADR-0003** (uang non-float
pada nilai yang ditampilkan).

## Acceptance criteria

- [ ] Popup "lihat polis lama" **ada dan terisi** data polis new business. *(AC 9 spec)*
- [ ] Tampilan dipilih menurut `Type`: **`QR` memakai tampilan induk**; `QP`, `TP`, `TR` memakai
      variannya masing-masing. **Keempat jenis didukung.** *(AC 9 spec)*
- [ ] Membuka popup **tidak mengubah apa pun** pada endorsement maupun pada polis lama — ia murni
      baca.
- [ ] Nilai uang yang ditampilkan **persis** seperti tersimpan; tidak ada pembulatan tampilan yang
      menyesatkan. (**ADR-0003**)
- [ ] Popup dapat dibuka **berulang kali** tanpa efek samping.

## Catatan — rule mati, perilaku hidup

⚠️ `[keputusan work owner]` **Rule `GetOldDetail_EDM` tidak dimigrasikan** — 4 dari 5 langkahnya
`<pyStepsBlockName>//` (baris 244, 539, 841, 975); hanya step 2 `Obj-Open-By-Handle` "Get data Life
Old" (page `WorkLife`) yang hidup.

**Tetapi perilakunya WAJIB ADA.** `[terverifikasi]` Rule itu dijalankan oleh **keempat** tombol
sebelum popup tampil — membuangnya begitu saja **mengosongkan keempat popup**. Di sistem baru,
pemuatan polis lama memakai **jalur baca biasa** berkunci ID pega lama / nomor polis lama yang sudah
dipegang case sejak tiket 02 — tanpa activity khusus.

⚠️ **OQ-066 berlaku** — penetapan "rule mati" bersandar pada **keputusan work owner**; penanda
`<pyStepsBlockName>` hanya pendukung.

## Blocker

**Tidak ada.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
