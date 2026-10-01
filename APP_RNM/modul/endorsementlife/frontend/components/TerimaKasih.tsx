// `Section/ConfirmSubmitEDM.xml` (FlowAction `SetJsonPolisEDMLife_Confirm` b91): `Thank you for Submit !`
// b526, `No. Endorsement` b654 `.PremiumListSummary.PL_NUMBER_EDM`, `Close` b1366 → `finishAssignment`.

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import { TERIMA_EDM } from '../labels'
import { sel } from '../tampilan'

export default function TerimaKasih({ noEndorsement, onTutup }: { noEndorsement: string; onTutup: () => void }) {
  return (
    <Modal judul={TERIMA_EDM.terimaKasih} onTutup={onTutup} labelBatal={TERIMA_EDM.close}>
      <dl className="edm-kepala">
        <div className="edm-kepala__medan">
          <dt>{TERIMA_EDM.noEndorsement}</dt>
          <dd>{sel(noEndorsement)}</dd>
        </div>
      </dl>
    </Modal>
  )
}
