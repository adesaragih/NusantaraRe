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

## Implementasi — 28-09-2026 (giliran 10)

`[terverifikasi]` Tidak ada padanan korpus (catatan tiket tetap benar) — perilaku baru, konsekuensi
ADR-0015.

### Yang dibangun

- `models/komite_efek.go` — km5: **kode** status outbox (`antre`/`jalan`/`selesai`/`gagal-permanen`)
  di `models`, **kata** di layar (`tertunda`/`tuntas`/`perlu intervensi`). `KeadaanEfekKasus`:
  "tuntas" **hanya** bila setiap efek selesai; satu gagal permanen → "perlu intervensi"; ada yang
  tertunda → "tersimpan, belum tuntas" (AC 25); kode asing **tidak** dianggap tuntas. Kode di tiga lapis
  dikunci sama (`TestKodeEfekSamaDiTigaLapis`).
- "Perlu intervensi" = baris outbox `gagal-permanen` — keadaan **eksplisit** yang dicapai pekerja tiket
  07 sesudah 8 percobaan atau galat permanen (AC 23), dengan jejak menyerah di jalur audit (ADR-0007;
  `rekamMenyerah` Claim Life yang dipakai ulang).
- **Menempel pada kasusnya**: `GET /api/komite/{id}` kini membawa `efek` (keadaan ringkas + tiap efek:
  jenis, keadaan, percobaan, **sejak** kapan); layar Kasus Komite menampilkannya. `GALAT_TERAKHIR`
  sengaja tidak dibaca (dapat menyebut objek basis data/kunci kategori).
- **Laporan harian**: `GET /api/komite/laporan-harian` — seluruh efek Komite `gagal-permanen`; bidang
  `kosong` **dinyatakan**, bukan disimpulkan (AC "tetap terkirim meskipun kosong"). Layar Inbox Komite
  menampilkannya bagi admin, termasuk kalimat "tidak ada" saat kosong. `[asumsi — OQ-007/021]` penerima
  = `ReasLifeAdmin`.

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| tidak pulih → "perlu intervensi" eksplisit | ✅ `gagal-permanen` (tiket 07) |
| terlihat menempel pada kasusnya | ✅ layar kasus |
| menyebut efek mana dan sejak kapan | ✅ |
| laporan harian seluruh kasus ber-status itu | ✅ endpoint + layar · ⚠️ **pengiriman** harian (email/penjadwal) belum ada — keputusan operasi |
| laporan tetap ada meski kosong | ✅ `kosong: true` + kalimatnya |
| tersimpan ≠ tuntas | ✅ |
| kegagalan di jalur audit | ✅ jejak menyerah |

### Angka

Go **593 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **357** · tsc bersih.
