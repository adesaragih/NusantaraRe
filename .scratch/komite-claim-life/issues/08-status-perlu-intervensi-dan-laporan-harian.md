# 08: Status "perlu intervensi" di UI Komite + laporan harian

**Status:** ready-for-agent

**Blocked by:** 07 (worker pengirim — retry + anti-dobel)

## Hasil & nilai pengguna

Sebagai **anggota komite**, saya ingin melihat kasus yang efek keluarnya **gagal dan perlu
intervensi** — dan sebagai **operator**, saya ingin daftar itu sampai ke saya setiap hari tanpa
harus mencarinya. *(User story 26 di spec)*

Tanpa tiket ini, "wajib berhasil" (**ADR-0015**) hanya **memindahkan kegagalan diam ke antrean** —
tidak ada yang melihatnya. Itulah sebabnya ia tiket terpisah dari 07: **ini permukaan untuk manusia,
bukan mesin.**

## Area codebase

`internal/models` (keadaan "perlu intervensi" pada entri outbox), `internal/services` (kriteria
kapan sebuah entri masuk keadaan itu; penyusunan laporan harian), `internal/handlers` (endpoint
daftar kasus bermasalah), `frontend/` (penanda status di inbox dan layar kasus).

## Rule Pega sumber

**Tidak ada padanan di korpus.** `[terverifikasi]` Sistem lama **tidak memuat** permukaan pemantauan
apa pun untuk kegagalan efek keluar. Yang ada hanya pencatatan di sisi database:

| Objek | Peran |
| --- | --- |
| `POOLDATA.DIRECTTOKASIR_LOG` | `[terverifikasi]` jejak panggilan Kasir |
| `POOLDATA.MONITORING_KLAIM_LOG` | `[terverifikasi]` log layanan — dipakai jalur Claim Life |

Keduanya **log**, bukan antrean kerja: tidak ada yang menampilkannya kepada manusia, dan tidak ada
keadaan "belum selesai" yang dapat ditindaklanjuti.

**Tiket ini karena itu perilaku baru sepenuhnya**, bukan paritas — ia konsekuensi langsung
**ADR-0015**.

## ADR terkait

**ADR-0015** (status "perlu intervensi" di UI Komite + laporan harian adalah **syarat pengaman
wajib**, bukan tambahan), **ADR-0007** (kegagalan masuk jalur audit, bukan hanya log layanan),
**ADR-0014** (siapa yang melihat daftar itu mengikuti wewenang).

## Acceptance criteria

- [ ] Entri outbox yang **tidak pulih** setelah retry masuk keadaan **"perlu intervensi"** yang
      eksplisit — bukan diam di antrean. *(AC 23 spec)*
- [ ] Kasus ber-status itu **terlihat di UI Komite**, menempel pada kasusnya, bukan di halaman
      terpisah yang harus dicari.
- [ ] Tampilan menyebut **efek mana** yang gagal (InsertJson / Arasapas / Email / Kasir) dan
      **sejak kapan**.
- [ ] Ada **laporan harian** berisi seluruh kasus ber-status itu.
- [ ] Laporan harian tetap terkirim **meskipun kosong** — ketiadaan laporan tidak boleh ambigu
      dengan ketiadaan masalah.
- [ ] Keputusan yang efeknya belum tuntas dilaporkan **tersimpan**, bukan **tuntas**.
      *(AC 25 spec; **ADR-0015**)*
- [ ] Kegagalan juga tercatat di **jalur audit**, bukan hanya di log layanan. *(**ADR-0007**)*

## Catatan

⚠️ **Kasus yang paling penting terlihat adalah Kasir.** `[terverifikasi]` `HitServiceToKasirKMTLife_Act`
tidak ada di Claim Life — ia khas Komite, dan ia memicu jalur **pembayaran**. Klaim yang disetujui
tetapi gagal sampai ke Kasir adalah kegagalan yang paling mahal bila tidak terlihat.

`[keputusan work owner]` Siapa yang menindaklanjuti keadaan "perlu intervensi" — anggota komite,
admin, atau operator — **belum ditetapkan**. Tiket ini menampilkannya; **alur penanganannya**
keputusan terpisah. Tidak memblokir: yang penting kegagalan berhenti tersembunyi.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
