// Pilihan bersegmen (radio) modul Treaty Description - dipakai saringan daftar dan isian form.

export interface OpsiSegmen<T extends string> {
  nilai: T
  label: string
}

export default function Segmen<T extends string>({
  label,
  opsi,
  nilai,
  onPilih,
}: {
  label: string
  opsi: readonly OpsiSegmen<T>[]
  nilai: T
  onPilih: (v: T) => void
}) {
  return (
    <div className="treatydescription__segmen" role="radiogroup" aria-label={label}>
      {opsi.map((o) => (
        <button
          key={o.nilai === '' ? 'semua' : o.nilai}
          type="button"
          role="radio"
          aria-checked={nilai === o.nilai}
          className={
            nilai === o.nilai
              ? 'treatydescription__segmen-butir treatydescription__segmen-butir--aktif'
              : 'treatydescription__segmen-butir'
          }
          onClick={() => onPilih(o.nilai)}
        >
          {o.label}
        </button>
      ))}
    </div>
  )
}
