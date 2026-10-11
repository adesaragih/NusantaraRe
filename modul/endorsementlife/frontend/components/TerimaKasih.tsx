// Popup hasil Confirm (Accept) - "Endorsement confirmed" (permintaan work owner 09-10-2026, pola "PL Number created"
// PremiumList; menggantikan `Section/ConfirmSubmitEDM.xml` `Thank you for Submit !` b526 / `No. Endorsement` b654 /
// `Close` b1366). Endorsement No, EDM Type, Sub-type (bila terisi), Case ID; OK / X / Esc → kembali ke kotak masuk
// HANYA sesudah popup ditutup (`onTutup`). Nomor kosong (tak terduga) tampil "—", bukan gagal diam.

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import type { HasilPutusanEDM } from '../api'
import { HASIL_EDM, UMUM_EDM } from '../labels'
import { labelEdmType, sel, subJenisTerisi } from '../tampilan'

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

/** Baris label/nilai popup hasil (Sub-type hanya bila terisi). */
export function barisHasil(
  hasil: HasilPutusanEDM,
  kasus: { id: string; edmType: string; edmTypePerubahanData?: string; edmTypeBatal?: string },
): ReadonlyArray<readonly [string, string]> {
  const sub = subJenisTerisi(kasus)
  return [
    [HASIL_EDM.noEndorsement, sel(hasil.noEndorsement)],
    [HASIL_EDM.edmType, labelEdmType(kasus.edmType)],
    ...(sub !== '' ? [[HASIL_EDM.subType, sub] as const] : []),
    [HASIL_EDM.caseId, sel(kasus.id)],
  ]
}

export default function TerimaKasih({
  hasil,
  kasus,
  onTutup,
}: {
  hasil: HasilPutusanEDM
  kasus: { id: string; edmType: string; edmTypePerubahanData?: string; edmTypeBatal?: string }
  onTutup: () => void
}) {
  return (
    <Modal judul={HASIL_EDM.judul} onTutup={onTutup} labelBatal={HASIL_EDM.ok}>
      <dl className="edm-hasil" role="status">
        {barisHasil(hasil, kasus).map(([label, nilai]) => (
          <div key={label} className="edm-hasil__baris">
            <dt>{label}</dt>
            <dd>{nilai}</dd>
          </div>
        ))}
      </dl>
      {hasil.efekKeluar !== undefined && <AlarmEfek efek={hasil.efekKeluar} />}
    </Modal>
  )
}
