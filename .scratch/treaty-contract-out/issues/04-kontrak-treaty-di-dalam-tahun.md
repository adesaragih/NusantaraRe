# 04: Kontrak treaty di dalam tahun treaty

**Status:** ready-for-agent

**Blocked by:** 02 (pemilih jenis reasuransi), 03 (tahun treaty sebagai induk)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin membuat **kontrak treaty** di dalam sebuah tahun
treaty — dengan **jenis reasuransinya** dan **masa berlakunya sendiri** — dan mengubahnya kemudian
tanpa membuat duplikat. *(User story 4–5, 7 di spec)*

Kontrak inilah yang membuka kombinasi **(tahun, grup, jenis reasuransi)** yang dipakai reinsurer,
business, dan seluruh klausul di tiket-tiket berikutnya.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas kontrak treaty |
| `internal/repository` | Upsert dikunci `ID`; rujukan ke tahun treaty |
| `internal/services` | Gerbang masa berlaku |
| `internal/handlers` | Endpoint kontrak |
| `frontend/` | Grid kontrak di dalam layar tahun treaty |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyContract_Act` | `@BASECLASS` / `SAVETREATYCONTRACT_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyContract_Act.xml` | orkestrator simpan |
| `NewInputTreatyContract_Act` | `@BASECLASS` / `NEWINPUTTREATYCONTRACT_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewInputTreatyContract_Act.xml` | baris baru |
| `SetUbahTreatyContract` | `@BASECLASS` / `SETUBAHTREATYCONTRACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetUbahTreatyContract.xml` | muat untuk diubah |
| `SetTanggalTreatyContract` | `@BASECLASS` / `SETTANGGALTREATYCONTRACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetTanggalTreatyContract.xml` | isi tanggal |
| `SaveMasterTreatyContract_SQL` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyContract_SQL.xml` | → `POOLDATA.PEGA_TREATYCONTRACT` |
| `InboxTreatyContract` | `DATA-PORTAL` / `INBOXTREATYCONTRACT` / `RULE-HTML-HARNESS` | `Treaty Contract Out/Harness/InboxTreatyContract.xml` | titik masuk modul |

`[terverifikasi]` Parameter `PEGA_TREATYCONTRACT`: `ID`, `IDTreatyYear`, `ReinsTypeID`,
`ReinsTypeName`, `TreatyStartDate`, `TreatyEndDate`, pengguna, `TglUpdate`, + 2 keluaran.

`[data DBA]` Upsert dikunci `ID`; identitas `'1' || lpad(treatycontract_seq.nextval, 6, '0')`;
`TREATYSTARTDATE`/`TREATYENDDATE` di-`to_date(…,'DD/MM/YYYY')`; **tidak `COMMIT` sendiri**;
`StsSimpan` **1 = sukses / 0 = gagal**.

## ADR terkait

**ADR-0006** (identitas lewat sequence), **ADR-0007** (jejak audit), **ADR-0015** (kegagalan
eksplisit).

## Acceptance criteria

- [ ] Kontrak treaty dibuat **di dalam** sebuah tahun treaty dan menyimpan **rujukan ke tahun itu**.
      *(AC 7 spec; User story 4)*
- [ ] Menyimpan kontrak yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec — sisi kontrak)*
- [ ] Masa berlaku kontrak yang **berakhir sebelum dimulai** **ditolak**, dengan pesan yang
      menyebut field-nya. *(AC 9 spec — sisi kontrak)*
- [ ] Jenis reasuransi kontrak dipilih dari daftar tersaring tiket **02** — tidak diketik bebas.
      *(AC 10, 11 spec)*
- [ ] Identitas kontrak **tidak pernah diketik pengguna**; berbentuk `'1'` + 6 digit.
      *(AC 5, 6 spec; **ADR-0006**)*
- [ ] Tanggal mulai dan akhir kontrak tersimpan sebagai **tanggal**, bukan teks. *(AC 53 spec)*
- [ ] Kontrak **tidak dapat dibuat** tanpa tahun treaty induk yang ada.
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-0007**)*

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Kombinasi, bukan foreign key.** `[fakta bisnis — work owner]` Reinsurer, business, dan klausul
**tidak** menyimpan `ID` kontrak. Mereka menggantung pada **(TreatyYear, TreatyGroupID,
ReinsTypeID)**. Kontrak "membuka" kombinasi itu, tetapi tidak memilikinya. Ini **berbeda** dari
Master Contract Retro Life, di mana anak-anak menggantung pada `ID` kontrak — jangan menyalin pola
dari sana.

⚠️ **Editor master, bukan proses berjenjang.** `[terverifikasi]` Nol rule `Flow`, nol rule `When`,
nol `StatusAkseptasi` di seluruh 303 berkas modul. Tidak ada Submit/Decline, tidak ada assignment,
tidak ada SLA.

`[terverifikasi]` **RBAC tidak dapat direkonstruksi** dari korpus — modul ini tidak memuat ekspresi
visibilitas ber-workbasket (`discovery/flows/_METHOD-noflow.md` §3.5). Siapa yang boleh mengedit
kontrak adalah keputusan terpisah, di luar tiket ini.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```
