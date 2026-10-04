// Halaman Master Data - KERANGKA (sesi c3, 04-10-2026). Digantikan layar sesi 0f (daftar per tabel master, tambah /
// ubah, aktif / nonaktif) di atas `GET /api/masterdata` dan `/api/masterdata/{tabel}`.

import { KELOMPOK_MASTERDATA } from '../menu'

export default function HalamanMasterData() {
  return (
    <section>
      <h1>{KELOMPOK_MASTERDATA}</h1>
    </section>
  )
}
