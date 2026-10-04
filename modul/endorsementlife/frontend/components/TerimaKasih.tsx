// `Section/ConfirmSubmitEDM.xml` (FlowAction `SetJsonPolisEDMLife_Confirm` b91): `Thank you for Submit !`
// b526, `No. Endorsement` b654 `.PremiumListSummary.PL_NUMBER_EDM`, `Close` b1366 → `finishAssignment`.

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import type { HasilPutusanEDM } from '../api'
import { TERIMA_EDM, UMUM_EDM } from '../labels'
import { sel } from '../tampilan'

/** Alarm efek keluar di layar (E5): kalimat tetap dari server, tanpa alamat. */
function AlarmEfek({ efek }: { efek: NonNullable<HasilPutusanEDM['efekKeluar']> }) {
  const gagal = [...efek.gagal, ...efek.tidakDiantre]
  if (efek.dilewati) return <p className="edm-catatan">{UMUM_EDM.efekDilewati}</p>
  if (gagal.length === 0) return null
  return (
    <div role="alert">
      <p className="edm-catatan">{UMUM_EDM.efekGagal}</p>
      <ul className="edm-pesan">
        {gagal.map((g) => (
          <li key={g}>{g}</li>
        ))}
      </ul>
    </div>
  )
}

export default function TerimaKasih({ hasil, onTutup }: { hasil: HasilPutusanEDM; onTutup: () => void }) {
  const { noEndorsement } = hasil
  return (
    <Modal judul={TERIMA_EDM.terimaKasih} onTutup={onTutup} labelBatal={TERIMA_EDM.close}>
      <dl className="edm-kepala">
        <div className="edm-kepala__medan">
          <dt>{TERIMA_EDM.noEndorsement}</dt>
          <dd>{sel(noEndorsement)}</dd>
        </div>
      </dl>
      {hasil.efekKeluar !== undefined && <AlarmEfek efek={hasil.efekKeluar} />}
    </Modal>
  )
}
