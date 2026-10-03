# 02: Jenis reasuransi — master dibaca + saringan non-life

**Status:** ready-for-agent

**Blocked by:** 01 (skema harus ada)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin memilih **jenis reasuransi** dari daftar master — dan
daftar itu hanya memuat jenis yang **berlaku untuk non-life** — sehingga saya tidak pernah mengetik
bebas dan tidak pernah salah memilih jenis milik lini lain. *(User story 8–10 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Baca master jenis reasuransi (read-only) |
| `internal/services` | Saringan non-life |
| `internal/handlers` | Endpoint daftar jenis reasuransi |
| `frontend/` | Pemilih jenis reasuransi yang dipakai layar kontrak dan seluruh grid klausul |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| **RD dominan** | `ASM-FW-GISFW-INT-REINSURANCETYPE` / `BROWSEREINSURANCETYPE_RD_OLD_LJT_ID_ISNOTNULL` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` | ⚠️ dipakai **11 grid** |
| RD kedua | `ASM-FW-GISFW-INT-REINSURANCETYPE` / `BROWSEREINSURANCETYPE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseReinsuranceType_RD.xml` | dipakai 1 layar — **layar master**, di luar konteks ini |
| `GetMasterReinsTypeContract` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterReinsTypeContract.xml` | ⚠️ membaca tabel **JSON** — mati |

⚠️ **Nama berbohong — dan yang "Old" justru yang benar** `[terverifikasi]`:

| Rule | RuleSet | Commit | Dipakai |
| --- | --- | --- | ---: |
| `…!BROWSEREINSURANCETYPE_RD_OLD_LJT_ID_ISNOTNULL` | **01-01-91** | **2026-02-05** | **11** grid |
| `…!BROWSEREINSURANCETYPE_RD` | 01-01-83 | 2025-06-02 | 1 layar |

`[keputusan work owner]` Yang **dimigrasikan** adalah perilaku yang dipakai 11 grid:

```
.ID NOT IN ("10004","10011","10012","10021","10022","10025","10026","10028",
            "10248","10249","10018","10217")
AND .Flag = "active"
AND .Type IN ("1","2","3")
```

⚠️ **Penyimpangan sadar 8 — nama jujur.** `[keputusan work owner]` Komponen di sistem baru
**tidak** mengandung kata `"Old"`.

`[terverifikasi]` Class `ASM-FW-GISFW-INT-REINSURANCETYPE` adalah **class yang sama** dengan master
jenis reasuransi Life (`Master Contract Retro Life/ReportDefinition/BrowseReinsuranceTypeLimit_RD.xml`).
Jadi master ini **melayani life dan non-life sekaligus**; saringanlah yang memisahkan.

## ADR terkait

**ADR-0015** (kegagalan ditangani eksplisit, tidak ditelan).

## Acceptance criteria

- [ ] Jenis reasuransi dipilih dari **master**, tidak pernah diketik bebas. *(AC 10 spec)*
- [ ] ⚠️ Daftar disaring **persis** seperti existing: **dua belas ID di-blacklist** (`10004`,
      `10011`, `10012`, `10021`, `10022`, `10025`, `10026`, `10028`, `10248`, `10249`, `10018`,
      `10217`) **dan** `Flag = "active"` **dan** `Type` termasuk `1`, `2`, `3`.
      *(AC 11 spec; `[keputusan work owner]` — ditiru apa adanya, **jangan digeneralkan**)*
- [ ] Master jenis reasuransi **tidak ditulis** oleh konteks ini. Test yang menemukan tulisan ke
      master itu **gagal**. *(AC 12 spec)*
- [ ] ⚠️ Nama komponen, endpoint, dan fungsi di sistem baru **tidak mengandung kata "Old"**
      meskipun rule sumbernya bernama demikian. *(AC 13 spec; penyimpangan sadar 8)*
- [ ] Daftar dibaca dari **satu tempat**, dipakai layar kontrak maupun seluruh grid klausul —
      bukan disalin per layar.
- [ ] Master yang **kosong atau tidak terbaca** menghasilkan kegagalan yang **terlihat**, bukan
      daftar kosong yang diam. *(**ADR-0015**)*

## Blocker

**Tidak ada.** **OQ-020 ditutup** — arti `ReinsTypeID` sudah jelas: master bersama life & non-life,
dipisahkan oleh saringan di atas.

## Catatan

⚠️ **Beda ejaan antar konteks.** `.Flag` bernilai **`1`** di Master Contract Retro Life, tetapi
**`"active"`** di sini. Jangan menyalin nilai dari konteks Life.

⚠️ **Jangan ikut memigrasikan `GetMasterReinsTypeContract`** — kueri itu membaca `M_TREATYCONTRACT`
dan `M_TREATYYEAR` lewat `a.JSONDATA.…`, dua tabel JSON yang sudah **mati** (tiket 01).

⚠️ **Layar master jenis reasuransi ada di luar konteks ini.** `[keputusan work owner]`
`Harness/InboxTreatyContractReinsType.xml` (`DATA-PORTAL!INBOXTREATYCONTRACTREINSTYPE`) dan RD kedua
adalah perilaku **layar master**, konteks/menu tersendiri. Di sini: **baca saja**.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — saringan blacklist hanya terbukti benar
bila diuji terhadap master yang benar-benar memuat ID yang dikecualikan.

```
go test ./internal/...
cd frontend && npm test
make check
```
