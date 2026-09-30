# 04: Tujuh pemilih master — dan RI Rate yang bukan RI Risk

**Status:** ready-for-agent

**Blocked by:** 02 (pemilih mengisi field pada produk)

## Hasil & nilai pengguna

Sebagai **admin master**, saya memilih ceding, pemegang polis, mata uang, sumber bisnis, dan penyebab
klaim dari **master masing-masing**; dan sebagai **underwriter** saya memilih **RI Risk** dan
**RI Rate** dari dua master yang **berbeda** — sehingga tidak ada nama yang diketik bebas dan tidak
ada dua master yang tertukar. *(User story 11–17 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Tujuh jenis nilai master beserta identitasnya |
| `internal/repository` | Pembacaan tiap master |
| `internal/services` | Penolakan nilai di luar master |
| `internal/handlers` | Tujuh endpoint daftar master |
| `frontend/` | Tujuh dialog pemilih |

## Rule Pega sumber

`[terverifikasi]` Sebelas FlowAction modul ini seluruhnya bercorak **pemilih dan konfirmasi**, bukan
langkah proses. Tujuh di antaranya pemilih master:

| Pemilih | Class / Nama / Tipe | Arti `[fakta bisnis — work owner]` |
| --- | --- | --- |
| `ChooseCeding` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECEDING` / `RULE-OBJ-FLOW-ACTION` | perusahaan ceding |
| `ChoosePolicyHolder` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSEPOLICYHOLDER` / `RULE-OBJ-FLOW-ACTION` | pemegang polis |
| `ChooseCurrency` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECURRENCY` / `RULE-OBJ-FLOW-ACTION` | mata uang |
| `ChooseSOB` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSESOB` / `RULE-OBJ-FLOW-ACTION` | **SOB = Source of Business** |
| `ChooseCauseOfLoss` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECAUSEOFLOSS` / `RULE-OBJ-FLOW-ACTION` | penyebab klaim/rugi |
| `ChooseRIRisk` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSERIRISK` / `RULE-OBJ-FLOW-ACTION` | master **RI Risk** — mengisi `RIRISKID`/`RIRISK` |
| **`ChooseRIRate`** | ⚠️ **`ASM-FW-GISFW-DATA-PLAN`** / `CHOOSERIRATE` / `RULE-OBJ-FLOW-ACTION` | master **RI Rate** — **class berbeda** |

⚠️ `[fakta bisnis — work owner]` **RI Rate ≠ RI Risk.** Dua master berbeda; **jangan digabung**.
**RI Risk** melekat pada **produk induk**; **RI Rate** melekat pada **baris plan** (tabel anak,
tiket 06).

`[terverifikasi]` Pendukung: `SearchPolicyHolder_act` (`ASM-FW-GISFW-…` / `SEARCHPOLICYHOLDER_ACT` /
`RULE-OBJ-ACTIVITY`), `SetRIRate` (`…` / `SETRIRATE`), `SetParamRate` (`…` / `SETPARAMRATE`).

## ADR terkait

**ADR-0001** (master pendukung dikelola konteks lain — modul ini **memilih dari** mereka),
**ADR-0003** (nilai bermata uang mengikuti mata uang terpilih).

## Acceptance criteria

- [ ] Ketujuh nilai dipilih dari **masternya masing-masing** — **tidak** diketik bebas.
      *(AC 18 spec)*
- [ ] ⚠️ **RI Rate dan RI Risk adalah dua master berbeda** dan **tidak pernah tertukar**. RI Risk
      mengisi produk induk; RI Rate mengisi **baris plan**. Test yang menemukan satu sumber untuk
      keduanya **gagal**. *(AC 19 spec; `[fakta bisnis — work owner]`)*
- [ ] Nilai yang dipilih tersimpan **beserta identitas masternya**, bukan hanya namanya —
      `CEDING`+`CEDINGID`, `SOBNAME`+`SOBID`, `CAUSE`+`CAUSEID`, `RIRISK`+`RIRISKID`. *(AC 20 spec)*
- [ ] Pilihan yang **tidak ada** di master **ditolak**. *(AC 21 spec)*
- [ ] Nama yang ditampilkan dapat **dibangun ulang** dari identitasnya — nama tersimpan adalah
      salinan tampilan, bukan sumber kebenaran.
- [ ] Pemilih pemegang polis dapat **dicari**, tidak hanya digulir — daftar master dapat besar.
- [ ] Mata uang terpilih berlaku pada **seluruh nilai uang** produk itu. *(**ADR-0003**)*
- [ ] Ketujuh daftar dibaca dari basis data, **tidak ditanam** sebagai konstanta di kode.

## Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **OQ kecil — tidak memblokir:** kepanjangan **"RI Rate"** dan **"RI Risk"** belum
diketahui. Pemilik: **Product+UW**. Tiket ini dapat selesai tanpa jawaban itu — yang mengikat adalah
**keduanya master terpisah**, bukan namanya.

## Catatan

⚠️ **Satu pemilih ber-class berbeda.** `[terverifikasi]` Enam pemilih ber-class
`ASM-FW-GISFW-INT-PRODUCT_LIFE`; **`ChooseRIRate` ber-class `ASM-FW-GISFW-DATA-PLAN`**. Ini sejalan
dengan perannya: RI Rate dipakai pada **plan**, bukan pada produk induk.

⚠️ **Master pendukung di luar cakupan.** Modul ini **memilih dari** ketujuh master; ia **tidak
mengelolanya**. Lihat §Out of Scope spec.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```
