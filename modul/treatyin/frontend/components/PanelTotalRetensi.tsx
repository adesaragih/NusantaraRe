// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { Panel } from '../../../../inti/frontend/components/ui/dasar'
import { TOTAL_RETENSI } from '../labels'
import type { BarisTotalRetensiWarisan } from '../api'
import { selAngka } from './angka'

export default function PanelTotalRetensi({ baris }: { baris: readonly BarisTotalRetensiWarisan[] }) {
  return (
    <Panel judul={TOTAL_RETENSI.judul}>
      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              {/* Judul kolom pertama ADALAH nama panelnya — @152591 — dan
                  isinya mata uang (`.Currency` @161105). */}
              <th scope="col">{TOTAL_RETENSI.judul}</th>
              <th scope="col">{TOTAL_RETENSI.kolomNilai}</th>
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={2}>{TOTAL_RETENSI.tanpaBaris}</td>
              </tr>
            )}
            {baris.map((b) => (
              <tr key={b.mataUang}>
                <td>{b.mataUang}</td>
                {/* UANG — empat desimal, pemisah ribuan titik. */}
                <td>{selAngka('uang', b.nilai)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="trin__aksi" role="group" aria-label={TOTAL_RETENSI.judul}>
        {/* ⛔ MATI, dan sebabnya tertulis di sebelahnya. */}
        <button type="button" className="btn" disabled>
          {TOTAL_RETENSI.perbarui}
        </button>
        <span className="trin__redup">{TOTAL_RETENSI.perbaruiPetunjuk}</span>
      </div>
    </Panel>
  )
}
