// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { useState } from 'react'

import { Kosong, Modal, Panel } from '../../../../inti/frontend/components/ui/dasar'
import type { BarisKategoriLampiran, BarisLampiranWarisan } from '../api'
import {
  KOLOM_LAMPIRAN,
  KOLOM_LIHAT_BERKAS,
  LAMPIRAN,
} from '../labels'

/**
 * Panel Attachment — Category · Count · Upload file · View File.
 *
 * ⛔ Disalin dari `Section/WorkAttachments.xml`: keempat kolomnya, kedua
 * tombolnya, dan spanduk birunya. Datanya dari `M_ATTACHMENTTREATY_2`,
 * tabel WARISAN — nol tabel baru dibuat untuk panel ini.
 *
 * ⛔ SELURUH kategori tampil, termasuk yang nol berkas: kolom `Count` tidak
 * akan pernah berbunyi `0` kalau barisnya disembunyikan saat kosong.
 *
 * ⚠️ Kategori yang pasangan kode↔namanya BELUM dipastikan ditandai, dan
 * kodenya ditampilkan menggantikan namanya. Menebak pasangannya menaruh
 * berkas di kategori yang salah, dan itu baru ketahuan bertahun kemudian.
 */
export default function PanelLampiran({
  kategori,
  berkas,
}: {
  kategori: readonly BarisKategoriLampiran[]
  berkas: readonly BarisLampiranWarisan[]
}) {
  /**
   * ⭐ PERMINTAAN PERUBAHAN, satu-satunya tempat ronde ini menyimpang dari
   * Pega — dan pemilik proses yang memintanya, pada keterangan gambar `25`:
   *
   *   "dan view upload saran dibuatkan pop up dan ini diubah menjadi
   *    design nya bagus"
   *
   * ⛔ Yang gambar 25 perlihatkan adalah JENDELA CHROME TERPISAH
   * (`ShowAttachmentTreaty`), lengkap dengan bilah alamat `appdev...`. Di
   * sini ia menjadi modal di dalam halaman.
   *
   * ⚠️ `Modal` DIPANGGIL dari `inti`, tidak ditulis ulang — modal kedua di
   * repositori ini berarti dua perilaku tutup, dua perangkap fokus, dan dua
   * tempat untuk salah.
   */
  const [berkasDilihat, setBerkasDilihat] = useState<string | null>(null)
  const kategoriDilihat = kategori.find((k) => k.kode === berkasDilihat) ?? null
  return (
    <Panel judul={LAMPIRAN.judul}>
      {/* Spanduk biru — ATURAN nama berkas, bukan hiasan. Ditegakkan di
          services (`NamaBerkasAman`) dan dinyatakan di sini. */}
      <span className="trin__spanduk" role="note">
        {LAMPIRAN.spanduk}
      </span>

      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              {KOLOM_LAMPIRAN.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {kategori.length === 0 && (
              <tr>
                <td colSpan={KOLOM_LAMPIRAN.length}>
                  <Kosong pesan={LAMPIRAN.tanpaLampiran} petunjuk={LAMPIRAN.petunjukLampiran} />
                </td>
              </tr>
            )}
            {kategori.map((k) => (
              <tr key={k.kode}>
                {/* ⭐ NAMA apa adanya — kesebelasnya, persis tangkapan layar
                    pemilik proses 6 Oktober 2026. Sebelumnya empat baris
                    berbunyi KODE-nya (`00003 nama kategori belum
                    dipastikan`), sebab panel dirender dari katalog basis
                    data alih-alih dari daftar layar.

                    ⛔ Pertanyaan kode↔nama TIDAK tertutup oleh ini, dan
                    tidak boleh terlihat tertutup: `dipastikan` tetap
                    dibawa, dan jalur UNGGAH — yang memerlukan kodenya —
                    masih mati. Lihat
                    `docs/PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md`. */}
                <td>{k.nama}</td>
                {/* ⛔ Cacah TIDAK diformat — ia butir, bukan uang. */}
                <td>{String(k.cacah)}</td>
                <td>
                  {/* ⛔ `Upload file` MATI — jalur unggahnya belum ada, dan
                      tombol hidup yang tidak mengunggah apa pun berbohong.
                      Gambar 24 memperlihatkan modal "ASM Attach Content"
                      di baliknya. */}
                  <button type="button" className="btn btn--sm" disabled>
                    {LAMPIRAN.unggah}
                  </button>
                </td>
                <td>
                  {/* ⭐ `View File` HIDUP — ia hanya menampilkan apa yang
                      sudah dibaca, dan menampilkan itu yang panel baca-saja
                      memang boleh lakukan. */}
                  <button
                    type="button"
                    className="btn btn--sm"
                    onClick={() => {
                      setBerkasDilihat(k.kode)
                    }}
                  >
                    {LAMPIRAN.lihatBerkas}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {berkasDilihat !== null && (
        <Modal
          judul={LAMPIRAN.judulLihatBerkas}
          onTutup={() => {
            setBerkasDilihat(null)
          }}
          labelBatal={LAMPIRAN.tutup}
          lebar
        >
          <div className="table-wrap">
            <table className="trin__tabel">
              <thead>
                <tr>
                  {/* ⛔ DUA kolom, persis gambar 25 — `File Name` dan
                      `Type`. Tidak tiga, tidak empat. */}
                  {KOLOM_LIHAT_BERKAS.map((k) => (
                    <th key={k} scope="col">
                      {k}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {berkasKategori(berkas, kategoriDilihat).length === 0 && (
                  <tr>
                    <td colSpan={KOLOM_LIHAT_BERKAS.length}>{LAMPIRAN.tanpaBaris}</td>
                  </tr>
                )}
                {berkasKategori(berkas, kategoriDilihat).map((b) => (
                  <tr key={b.id}>
                    <td>{b.namaBerkas}</td>
                    <td>{b.jenisMime}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Modal>
      )}
    </Panel>
  )
}

/**
 * Berkas milik SATU kategori.
 *
 * ⚠️ Disaring menurut `kode`, bukan menurut nama: nama kategori yang
 * pasangannya belum dipastikan ditampilkan sebagai kodenya, dan menyaring
 * menurut yang tampil akan menyaring menurut dua hal yang berbeda.
 */
function berkasKategori(
  berkas: readonly BarisLampiranWarisan[],
  kategori: BarisKategoriLampiran | null,
): readonly BarisLampiranWarisan[] {
  if (kategori === null) return []
  // ⛔ Kode bila ia diketahui, NAMA bila tidak. Keempat kategori yang
  // kodenya belum dipastikan akan selalu berbunyi `No items` kalau
  // penyaringnya hanya kode — dan kosong yang salah terbaca persis seperti
  // kosong yang benar.
  if (kategori.kode !== '') return berkas.filter((b) => b.kodeKategori === kategori.kode)
  return berkas.filter((b) => b.namaKategori === kategori.nama)
}
