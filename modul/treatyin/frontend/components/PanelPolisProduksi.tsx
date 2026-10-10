// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { KOLOM_POLIS_PRODUKSI, POLIS_PRODUKSI } from '../labels'
import type { BarisPolisProduksi } from '../api'

export default function PanelPolisProduksi({ baris }: { baris: readonly BarisPolisProduksi[] }) {
  return (
    <section className="trin__polis" aria-label={POLIS_PRODUKSI.judul}>
      <h4 className="trin__polis-judul">{POLIS_PRODUKSI.judul}</h4>
      <div className="table-wrap">
        <table className="trin__tabel trin__tabel--pega">
          <thead>
            <tr>
              {KOLOM_POLIS_PRODUKSI.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={KOLOM_POLIS_PRODUKSI.length}>{POLIS_PRODUKSI.tanpaBaris}</td>
              </tr>
            )}
            {/* ⛔ Keempatnya TEKS — `Policy No` dan `Pega ID` pengenal,
                `Quarter` dan `Quarter Year` `VARCHAR2`. Pengenal tidak
                pernah diformat. */}
            {baris.map((b, i) => (
              <tr key={b.nomorPolis + '|' + b.pegaID + '|' + String(i)}>
                <td>{b.nomorPolis}</td>
                <td>{b.pegaID}</td>
                <td>{b.kuartal}</td>
                <td>{b.tahunKuartal}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  )
}
