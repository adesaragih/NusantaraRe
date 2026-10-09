// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { Kosong, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { KOLOM_HISTORY, LAMPIRAN } from '../labels'

/**
 * Panel History — Date · PIC · Approval · Comment.
 *
 * ⛔ Sumbernya `T_VIEW_COMMENT` yang SUDAH terisi 11.365 baris, dibaca
 * `repository.BacaCatatan` yang sudah ada. Nol pembacaan baru, nol tabel
 * baru — keempat kolomnya berpadanan satu-satu.
 *
 * ⚠️ Teks kosongnya `No items`, disalin dari layar lama apa adanya.
 */
export default function PanelHistory({ baris }: { baris: readonly { tanggal: string; operator: string; disetujui: string; catatan: string }[] }) {
  return (
    <Panel judul={LAMPIRAN.judulHistory}>
      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              {KOLOM_HISTORY.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={KOLOM_HISTORY.length}>
                  <Kosong pesan={LAMPIRAN.tanpaRiwayat} petunjuk={LAMPIRAN.petunjukHistory} />
                </td>
              </tr>
            )}
            {baris.map((b, i) => (
              <tr key={i}>
                <td>{b.tanggal}</td>
                <td>{b.operator}</td>
                <td>{b.disetujui}</td>
                <td>{b.catatan}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Panel>
  )
}
