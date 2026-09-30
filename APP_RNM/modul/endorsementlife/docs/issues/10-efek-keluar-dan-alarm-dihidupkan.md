# 10: Efek keluar Arasapas + alarm yang dihidupkan

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 09 (efek keluar berjalan setelah penyimpanan selesai dan keadaannya diketahui)

## Hasil & nilai pengguna

Sebagai **tim operasi**, saya ingin endorsement yang tersimpan diteruskan ke Arasapas dengan alamat
yang **selalu dibaca runtime**, dan saya ingin **diberi tahu bila penyimpanan gagal** — karena hari
ini jalur endorsement berjalan **tanpa alarm sama sekali**. *(User story 37–40 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/services` | Orkestrasi efek keluar; gerbang lingkungan; deteksi keadaan separuh |
| `internal/repository` | Resolusi alamat dari `M_LINK_SERVICE`; pembacaan balik untuk deteksi |
| `internal/clients` | Klien Arasapas + pengirim email **di balik interface** |
| `internal/handlers` | Status efek keluar terbaca API |
| `frontend/` | Penanda "terkirim" / "tertunda" pada endorsement |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `serviceInsertArasapasLife_act` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SERVICEINSERTARASAPASLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/serviceInsertArasapasLife_act.xml` | efek bisnis — ⚠️ **rule berbeda** dari yang ber-class `ASM-FW-GISFW-WORK-LIFE` |
| `GetLinkService` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/GetLinkService.xml` | **resolusi alamat runtime** |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` | pemicu — step 13, 15, 16 |
| `IsPEGAPROD` | `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | `Endorsement Life/When/IsPEGAPROD.xml` | **flag lingkungan** |

`[terverifikasi]` Rantai di `InsertJsonPolisLife_Act`:

| Step | Isi | Penanda korpus | Nasib di sistem baru |
| ---: | --- | --- | --- |
| ~~13~~ | `RDB-List` "Cek sudah masuk atau blm datanya" | **REMARK** (6399) | ⚠️ **DIHIDUPKAN** |
| ~~15~~ | `call @baseclass.SendEmailNotification` | **REMARK** (6752) | ⚠️ **DIHIDUPKAN** |
| **16** | `call serviceInsertArasapasLife_act` | aktif, dijaga `IsPEGAPROD` (7230) | dipertahankan |

⚠️ `[terverifikasi]` **Endorsement hari ini berjalan tanpa alarm** — deteksi keadaan separuh dan
email kegagalan **keduanya mati**, sementara jalur new business memilikinya (**PL-06**).

## ADR terkait

**ADR-0013** (alamat di-resolve **runtime** dari `M_LINK_SERVICE`; **dilarang** sebagai literal,
konstanta, **maupun env var** — menggantikan **ADR-0004**), **ADR-0005** (`IsPEGAPROD` → flag
lingkungan), **ADR-0008** (efek keluar asinkron dengan antre ulang), **ADR-0007** (jejak audit).

## Acceptance criteria

- [ ] Di lingkungan **non-production**, efek keluar **tidak berjalan**; penyimpanan endorsement
      **tetap** berjalan penuh. *(AC 43 spec; **ADR-0005**)*
- [ ] Gerbang lingkungan adalah **satu** flag yang dibaca di **satu tempat**, bukan pemeriksaan
      tersebar.
- [ ] Alamat Arasapas di-resolve **runtime** dari `M_LINK_SERVICE` pada setiap panggilan. Test yang
      memindai kode untuk URL sebagai **literal, konstanta, atau pembacaan env var** **gagal** bila
      menemukannya. *(AC 44 spec; **ADR-0013**)*
- [ ] ⚠️ **Penyimpangan sadar — alarm dihidupkan.** Deteksi keadaan separuh **dan** email kegagalan
      **berjalan** di jalur endorsement, dengan perilaku **sama** seperti jalur new business.
      *(AC 45 spec; `[keputusan work owner]`)*
- [ ] Email terkirim **hanya** bila pembacaan balik menunjukkan penyimpanan **gagal** — ia **alarm**,
      bukan notifikasi bisnis. Penyimpanan yang berhasil **tidak** mengirim email.
- [ ] Kegagalan efek keluar **tidak** membatalkan endorsement yang sudah tersimpan; ia **tercatat,
      terlihat, dan dapat diulang**. *(AC 46 spec; **ADR-0008**)*
- [ ] Arasapas **kiriman keluar** dan pengirim email berada **di balik interface** dan **di-fake** di
      test; yang diperiksa adalah **efeknya** (panggilan terjadi/tidak, email terpicu/tidak).
- [ ] ⚠️ Pembacaan **pembayaran** Arasapas (tiket 01) **tidak** ikut di-fake — ia gerbang bisnis,
      diuji terhadap skema uji nyata. Keduanya **tidak** boleh tercampur di satu klien.

## Blocker

**Tidak ada.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
