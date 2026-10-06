// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { FORM_KONTRAK } from '../labels'

export default function MedanTakAda({ label }: { label: string }) {
  return (
    <div className="field">
      <label className="field__label">{label}</label>
      <input className="field__input" type="text" value="" readOnly disabled />
      <span className="trin__redup">{FORM_KONTRAK.takAdaDiWarisan}</span>
    </div>
  )
}
