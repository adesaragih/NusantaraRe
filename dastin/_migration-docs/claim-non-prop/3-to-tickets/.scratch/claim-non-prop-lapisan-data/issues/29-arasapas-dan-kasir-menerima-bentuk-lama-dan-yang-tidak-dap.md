---
status: selesai
menunggu-luar: [REQ-033]
---

# 29: Arasapas dan kasir menerima bentuk lama, dan yang tidak dapat diringkas tidak hilang diam-diam


> **SELESAI 19 September 2026 — dan view-nya belum menghasilkan selisih sama sekali.**
>
> ### Yang hilang: seluruh pokok tiket ini
>
> Tiket ini berbunyi *"keduanya membaca tabel kanonik **dan arsip muatan keluar**, lalu menghasilkan **selisih**"*. `V_AKSEPTASI_KOMPATIBEL` **tidak menyentuh `ARSIP_MUATAN_KELUAR` sama sekali** — ia mengeluarkan **nilai mutlak**.
>
> Itu bukan kekurangan kecil. Bukti di tiket ini sendiri menunjukkan sistem hilir **tidak menerima keadaan, ia menerima tambahan**: `GetDataOS` menjumlahkan seluruh baris berkunci sama ber-`STS_REJECT = 0`, lalu `SaveDataToOSAksep_Act` **mengurangkan** hasilnya; prosedur di seberang **selalu `INSERT`**.
>
> > **Mengeluarkan nilai mutlak ke jalur yang menjumlahkan tambahan berarti melipatgandakan setiap nilai pada pengiriman kedua.**
>
> Diperbaiki: `V_AKSEPTASI_KOMPATIBEL` kini `nilai sekarang − SUM(arsip WHERE DITOLAK = 0)`, per kunci kasar.
>
> ### Cacat kedua: kedua view mengelompokkan dengan cara berbeda
>
> Keduanya satu janji, bukan dua: **"tidak ada yang hilang diam-diam"**. Janji itu hanya berlaku bila setiap kelompok yang jatuh dari `V01` muncul di `V02`. Keduanya **tidak sama**, di dua tempat:
>
> | | `V01` | `V02` (lama) |
> |---|---|---|
> | Penyambungan ke `ALOKASI_LAYER` | empat field layer + mata uang | **tanpa `LAYER_BAGIAN` dan `LAYER_BAGIAN_JENIS`** |
> | Penyaring | `KEADAAN_BARIS = 'LENGKAP'` | **tidak ada** |
>
> Akibatnya ada kelompok yang **jatuh dari keduanya** — keadaan terburuk: hilang, dan tidak terlihat hilang. Penyambungan, penyaring, dan pengelompokan kini **identik**; yang berbeda hanya `HAVING`-nya, dan keduanya saling melengkapi tepat.
>
> ### Satu keadaan yang tetap tidak tertangkap, dan dicatat terbuka
>
> Kelompok yang **seluruh** barisnya bukan `'LENGKAP'` tidak muncul di `V01` maupun `V02` — penyaringnya membuang barisnya lebih dulu. Itu bukan kelalaian: sebabnya **"belum lengkap"**, bukan **"tidak sepakat"** — sebab berbeda, yang tempatnya di `V_PARITAS_SHADOW` dan uji tiket `37`. Ditulis di kaki `V02` supaya yang membacanya tidak menyimpulkan sendiri bahwa cakupan kedua view itu lengkap.
>
> **REQ-033 tetap menunggui, tidak menahan**: ia mengubah **berapa banyak** kelompok yang ditolak, bukan apakah penolakannya ada.

*Asal: `T-16` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** View kompatibilitas menghasilkan bentuk yang dikenal hilir; kelompok yang besaran per-layernya tidak sepakat **ditolak** dan muncul di view penolakan.

`V_AKSEPTASI_KOMPATIBEL`, `V_AKSEPTASI_DITOLAK`. Keduanya membaca tabel kanonik **dan arsip muatan keluar**, lalu menghasilkan **selisih**.

**Tidak termasuk:** Menegakkan kunci kasar sebagai `UNIQUE` kedua — **ditolak dengan alasan**: itu menjadikan keterbatasan bentuk lama sebagai aturan bisnis sistem baru.

**Blocked by:**

- **REQ-033** — **menunggui, tidak menahan** (dari luar papan): ia menentukan apakah agregasi halus ke kasar pernah menemui baris yang tidak sepakat — itu mengubah **aturan penolakan di view**, bukan keberadaan view-nya
- ~~`18`~~ *(selesai)* — Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran
- ~~`15`~~ *(selesai)* — Jembatan ke sistem lama berdiri sebagai tabel terpisah
- ~~`25`~~ *(selesai)* — Arsip muatan keluar: apa yang benar-benar dikirim, tersimpan sebagai catatan


**Dasar:** DECIDED(ADR-0023, ADR-0021 pengecualian `CASEID`). EVIDENCED: sapuan S1 — 17 parameter `InputParamOs.*` di `Activity\SaveDataToOSAksep_Act.xml`; `TypeLoss = .TreatyName` adalah **satuan kasar**.

- [x] Kolom dan tipenya sama dengan yang diterima hilir hari ini, termasuk `CASEID`.
- [x] Besaran aditif **dijumlahkan**; Limit Layer, Premi Deposit, persen premi pemulihan, kurs, dan porsi **tidak** — kelimanya justru menjadi penjaga kesepakatan.
- [x] Bila besaran per-layer berbeda antar baris yang dilebur, view **tidak menghasilkan baris**, dan kelompoknya muncul di view penolakan beserta sebabnya — **kini benar-benar berpasangan**: pengelompokan dan penyaring kedua view identik.
- [x] `LayerPart` dan `LayerPartType` tidak muncul di keluaran.
- [x] Keluarannya **selisih terhadap arsip muatan keluar**, bukan nilai mutlak — dan hanya arsip ber-`DITOLAK = 0` yang terhitung. **Sebelumnya tidak ada sama sekali.**

**Ketidakpastian:** **REQ-033** menentukan apakah agregasi halus→kasar pernah menemui baris yang tidak sepakat.

**Temuan S1 sudah TERTUTUP, dan tidak lewat REQ-018.** Sumber `OutOSAcc` ditemukan: `RDBList\GetDataOS.xml` dan `RDBList\GetDataCNPOS.xml` — EVIDENCED:

```sql
select sum(nvl(a.data_json.Value,0)) as "Value", ...
  from OS_AKSEPTASI_KLAIM a
 WHERE CASEID = {pyWorkPage.pzInsKey}
   AND a.data_json.TypeLoss = {InputParamOs.TypeLoss}
   AND a.data_json.Currency = {InputParamOs.Currency}
   AND STS_REJECT = 0
```

Penyaringnya **persis kunci alami ADR-0024**, dan agregatnya `SUM` atas seluruh baris berkunci sama. Jadi baris memang **tambahan**, bukan keadaan — dan itu pasangan wajib dari prosedur yang **selalu `INSERT`**. Naik dari DERIVED ke **EVIDENCED**.

**Satu syarat yang tidak terduga dan ikut mengikat**: penyaringnya memuat **`STS_REJECT = 0`**. Muatan yang ditolak **tidak** ikut dijumlahkan. Karena itu `ARSIP_MUATAN_KELUAR` berkolom `DITOLAK`, dan view menjumlahkan hanya yang tidak ditolak.

Yang tersisa di **REQ-018** hanya rekonsiliasi: apakah jumlah seluruh baris lama benar-benar sama dengan nilai sekarang di produksi. Bila tidak, itu temuan tentang **data lama**, bukan alasan mengubah rancangan.
