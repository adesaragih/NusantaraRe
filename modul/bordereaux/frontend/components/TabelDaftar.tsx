// Grid daftar Bordereaux - padanan grid `PortalBordereaux_Sec` (S17, `InboxBordereaux_RD`): 11 kolom Pega dan kolom
// tombol Edit / View / Delete menurut hak. Kolom tombol MENEMPEL di kanan saat tabel digulir (pola Treaty Contract
// Out, permintaan work owner 04-10-2026).

import type { BarisDaftar } from '../api'
import { BDX, teksStatus } from '../labels'

export default function TabelDaftar({
  daftar,
  onBuka,
  onHapus,
}: {
  daftar: BarisDaftar[]
  onBuka: (id: string, lihat: boolean) => void
  onHapus: (b: BarisDaftar) => void
}) {
  return (
    <div className="bordereaux__gulir">
      <table className="inbox__tabel bordereaux__tabel">
        <thead>
          <tr>
            <th>{BDX.id}</th>
            <th>{BDX.type}</th>
            <th>{BDX.business}</th>
            <th>{BDX.reffNoSoa}</th>
            <th>{BDX.reffNoBdx}</th>
            <th>{BDX.cedingCo}</th>
            <th>{BDX.treatyName}</th>
            <th>{BDX.reportStart}</th>
            <th>{BDX.reportEnd}</th>
            <th>{BDX.position}</th>
            <th>{BDX.status}</th>
            <th className="table__actions">{BDX.aksi}</th>
          </tr>
        </thead>
        <tbody>
          {daftar.map((b) => (
            <tr key={b.bdxId} className="inbox__baris">
              <td>{b.bdxId}</td>
              <td>{b.type}</td>
              <td>{b.typeBusiness}</td>
              <td>{b.reffNoSoa}</td>
              <td>{b.reffNoBdx}</td>
              <td>{b.cedingName}</td>
              <td>{b.treatyName}</td>
              <td>{b.reportStart}</td>
              <td>{b.reportEnd}</td>
              <td>{b.position}</td>
              <td>{teksStatus(b.statusAksep)}</td>
              <td className="table__actions">
                <span className="bordereaux__aksi">
                  {b.hak.ubah && (
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => onBuka(b.bdxId, false)}>
                      {BDX.edit}
                    </button>
                  )}
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => onBuka(b.bdxId, true)}>
                    {BDX.view}
                  </button>
                  {b.hak.hapus && (
                    <button type="button" className="btn btn--ghost btn--sm bordereaux__hapus" onClick={() => onHapus(b)}>
                      {BDX.hapus}
                    </button>
                  )}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
