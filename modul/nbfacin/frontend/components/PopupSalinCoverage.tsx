// Popup "Copy Coverage From" (tiket 53) - `NB FacIn\Section\CopyCoverageFrom.xml` (flow action `CopyCoverageFrom`):
// Location Address (`.SelectedLocationAddress`, daftar `pyWorkPage.OfferFacIn.LocationList`, tampil
// `.Property.RiskLocation.ASMAddress`) -> Object Item Name (`.SelectedObjectItem`, item lokasi terpilih, tampil
// `.ItemType`; tampil bila lokasi terpilih) -> tombol "Copy to Another Object Item" (tampil bila keduanya terpilih;
// `CopyCoverage_Act`: coverage item terpilih menggantikan coverage item yang sedang dibuka, lalu popup ditutup).

import { useState } from 'react'

import { Modal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import type { ObjekFire } from '../api'
import { SALIN_COVERAGE as S } from '../labels'
import type { PosisiItem } from './salinCoverage'

export default function PopupSalinCoverage({
  objek,
  onTutup,
  onSalin,
}: {
  objek: ObjekFire[]
  onTutup: () => void
  onSalin: (dari: PosisiItem) => void
}) {
  const [lokasi, setLokasi] = useState('')
  const [item, setItem] = useState('')
  const o = lokasi === '' ? -1 : Number.parseInt(lokasi, 10)
  const items = o >= 0 ? (objek[o]?.items ?? []) : []
  return (
    <Modal
      judul={S.dari}
      onTutup={onTutup}
      aksi={
        lokasi !== '' && item !== '' ? (
          <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary" onClick={() => onSalin({ o, i: Number.parseInt(item, 10) })}>
            {S.keItemLain}
          </button>
        ) : undefined
      }
    >
      <div className="nbf-salin-cov">
        <Pilih
          label={S.lokasi}
          value={lokasi}
          onChange={(v) => {
            setLokasi(v)
            setItem('')
          }}
          opsi={objek.map((x, k) => ({ value: String(k), label: x.riskLocation || x.roadName || `${S.objek} ${x.objectNo}` }))}
        />
        {lokasi !== '' && (
          <Pilih label={S.item} value={item} onChange={setItem} opsi={items.map((it, k) => ({ value: String(k), label: it.itemType || String(k + 1) }))} />
        )}
      </div>
    </Modal>
  )
}
