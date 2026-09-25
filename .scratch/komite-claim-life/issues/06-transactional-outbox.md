# 06: Transactional outbox — keputusan + daftar efek dalam satu transaksi

**Status:** ready-for-agent

**Blocked by:** 04b (rekam akseptasi — satu jalur simpan)

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin keputusan komite dan **antrean efek keluarnya** tersimpan dalam
**satu transaksi** — sehingga tidak ada efek yang hilang bila proses mati tepat setelah keputusan
disimpan. *(User story 24, 28 di spec)*

Ini fondasi jaminan "wajib berhasil"; pengirimannya sendiri ada di tiket 07.

## Area codebase

`internal/models` (entri outbox: efek apa, muatan apa, keadaan apa), `internal/repository`
(penulisan outbox dalam transaksi yang sama dengan keputusan), `internal/services` (penyusunan
daftar empat efek saat keputusan final).

Belum menyentuh `frontend/` — permukaan manusia ada di tiket 08.

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` step **6** | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` `Obj-Save` — **persist**; keempat efek berjalan **sesudahnya** |

`[terverifikasi]` Urutan efek yang harus masuk antrean, dibaca dari `<pyStepPageReference>` versi
2026-09-15:

| Urutan | Step | Efek |
| ---: | ---: | --- |
| 1 | **8** | `Call InsertJsonClaimLife_Act` |
| 2 | **10** | `Call serviceInsertArasapasClaimLife_act` |
| 3 | **11** | `Call SendEmailKlaimLife` |
| 4 | **12** | `Call HitServiceToKasirKMTLife_Act` — Kasir/pembayaran |

`[terverifikasi]` Step **10** dan **12** digerbangi `AcceptStatus = 1 && KomiteCount == KomiteLoop`
(baris 8648, 8887) — hanya pada keputusan Setuju di tingkat final.

## ADR terkait

**ADR-0015** (efek keluar wajib berhasil — transactional outbox, at-least-once; **menyimpang dari
ADR-0008**), **ADR-0013** (alamat endpoint di-lookup runtime, bukan disimpan di outbox sebagai URL),
**OQ-013** (batas transaksi dipegang Go).

## Acceptance criteria

- [ ] Keputusan dan **daftar keempat efeknya** tersimpan dalam **satu transaksi database**.
      *(AC 19 spec)*
- [ ] Bila transaksi gagal, **tidak ada** keputusan tersimpan **dan tidak ada** entri outbox — tidak
      ada keadaan separuh.
- [ ] Bila proses mati tepat setelah commit, antrean efek **tetap ada** dan dapat diambil worker.
- [ ] Tiap entri outbox membawa **ID idempoten unik** sejak dibuat. *(AC 21 spec)*
- [ ] Entri outbox menyimpan **kunci kategori** endpoint, **bukan URL** — alamat di-resolve saat
      kirim (**ADR-0013**).
- [ ] Efek hanya diantrekan pada keputusan yang benar-benar final; keputusan di tingkat bukan-
      terakhir tidak menghasilkan entri outbox.
- [ ] Keputusan dilaporkan **tersimpan**, belum **tuntas** — kedua keadaan itu dibedakan.
      *(**ADR-0015**)*

## Catatan — mengapa Komite berbeda dari Claim Life

⚠️ **ADR-0015 menyimpang dari ADR-0008.** Di Claim — Life, efek keluar boleh gagal tanpa memblokir.
Di sini tidak — pemicunya **Kasir**, integrasi **pembayaran** yang `[terverifikasi]` tidak ada di
Claim Life.

Kegagalan diam di Kasir berarti **klaim disetujui tetapi tidak pernah sampai ke pembayaran**, dan
tidak ada yang tahu.

`[keputusan work owner]` Step **14** `SetInformationData` **bukan** efek keluar — hanya temporary
penampung hasil submit. **Tidak** masuk outbox.

## Perintah verifikasi

```
go test ./internal/...
make check
```
