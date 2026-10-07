// Kotak teks panjang (text area) modul PremiumList Life — DAPAT diisi atau
// HANYA dibaca (permintaan work owner 03-10-2026: Underwriting Policy dan
// Marketing Note menjadi text area).
//
// ⛔ Kenapa bukan `Area` inti saja: `Area` tidak punya mode baca-saja, dan
// menggantinya dengan `Field` sebaris saat terkunci memotong teks panjang
// menjadi satu baris. Markupnya SAMA dengan `Area` (`field` > label >
// `textarea.field__input`), jadi gaya inti tetap berlaku.

export default function AreaTeks({
  label,
  value,
  onChange,
  readOnly,
  required,
  baris = 3,
}: {
  label: string
  value: string
  onChange?: (v: string) => void
  readOnly?: boolean
  required?: boolean
  baris?: number
}) {
  return (
    <div className="field">
      <label className="field__label">
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      <textarea
        className={'field__input' + (readOnly ? ' field__input--readonly' : '')}
        aria-label={label}
        rows={baris}
        value={value}
        readOnly={readOnly}
        onChange={(e) => onChange?.(e.target.value)}
      />
    </div>
  )
}
