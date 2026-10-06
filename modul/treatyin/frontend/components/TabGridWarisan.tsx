// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.
//
// ⭐ 6 Oktober 2026 — MODE UBAH. Permintaan pemakai: "apabila di klik edit
// maka tiap function itu bisa digunakan semuanya seperti add pada tiap
// table, kalau view baru tidak bisa". Di mode `ubah` sel menjadi kotak
// isian dan `Add`/`Delete` bekerja atas salinan baris di layar ini; di mode
// `lihat` grid baca-saja dan tombolnya tidak dirender — ekspor menjaga
// tombol-tombol itu dengan `TreatyIn.ViewState !='1'`.

import { useState } from 'react'

import { Kosong, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { FORM_KONTRAK, GRID_TAMBAH } from '../labels'
import type { ModeForm } from '../mode'

/**
 * Grid satu tab yang isinya dibaca dari dokumen warisan.
 *
 * ⛔ Kolomnya DISEBUT pemanggilnya, tidak disimpulkan dari barisnya: baris
 * boleh nol, dan kepala kolom harus tetap terlihat.
 *
 * ⚠️ Baris disalin ke keadaan SEKALI saat grid lahir; pemanggil memberi
 * `key` per kontrak supaya kontrak lain melahirkan grid baru.
 */
export default function TabGridWarisan({
  judul,
  kolom,
  baris,
  petunjukKosong,
  mode = 'lihat',
  bisaTambah = true,
}: {
  judul: string
  kolom: readonly string[]
  baris: readonly (readonly string[])[]
  petunjukKosong: string
  mode?: ModeForm
  /**
   * `false` hanya untuk grid yang barisnya LAHIR dari aksi lain, bukan
   * diketik — riwayat Information & Submit (baris baru lahir dari Submit).
   */
  bisaTambah?: boolean
}) {
  const [isi, setIsi] = useState<string[][]>(() => baris.map((b) => [...b]))
  const bisaUbah = mode === 'ubah'
  const tombol = bisaUbah && bisaTambah
  return (
    <Panel judul={judul}>
      {tombol && (
        <div className="trin__panel-kepala">
          <span className="trin__redup">{GRID_TAMBAH.petunjuk}</span>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              setIsi([...isi, kolom.map(() => '')])
            }}
          >
            {GRID_TAMBAH.tambah}
          </button>
        </div>
      )}
      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              {kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              {tombol && <th scope="col" aria-label={GRID_TAMBAH.hapus} />}
            </tr>
          </thead>
          <tbody>
            {isi.length === 0 && (
              <tr>
                <td colSpan={kolom.length + (tombol ? 1 : 0)}>
                  <Kosong pesan={FORM_KONTRAK.tanpaBaris} petunjuk={petunjukKosong} />
                </td>
              </tr>
            )}
            {isi.map((b, i) => (
              <tr key={i}>
                {b.map((v, j) => (
                  <td key={j}>
                    {bisaUbah ? (
                      <input
                        className="field__input"
                        type="text"
                        value={v}
                        aria-label={kolom[j] ?? ''}
                        onChange={(e) => {
                          setIsi(isi.map((r, x) => (x === i ? r.map((c, y) => (y === j ? e.target.value : c)) : r)))
                        }}
                      />
                    ) : (
                      v
                    )}
                  </td>
                ))}
                {tombol && (
                  <td>
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setIsi(isi.filter((_, x) => x !== i))
                      }}
                    >
                      {GRID_TAMBAH.hapus}
                    </button>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Panel>
  )
}
