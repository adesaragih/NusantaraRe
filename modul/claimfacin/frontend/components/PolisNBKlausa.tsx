// Isi tab polis NB Fac In hanya-baca Clauses (FIRE; `PolisNB.tsx`): kartu klausa `modul/nbfacin/frontend/components/
// TabClauses.tsx` (pola, bukan impor) tanpa tombol. Data `polisNBKlausa.ts`.

import { Kosong } from '../../../../inti/frontend/components/ui/dasar'
import type { Halaman } from '../api'
import { KLAUSA_NB as K, TEKS_KLAUSA_NB as T, adaArgumenNB, klausaNB } from '../polisNBKlausa'

export function TabClausesNB({ h }: { h: Halaman }) {
  const daftar = klausaNB(h)
  if (daftar.length === 0) return <Kosong pesan={T.kosong} />
  return (
    <div className="claimfacin__nb-klausa">
      {daftar.map((k, i) => (
        <article key={`${k.clauseCode}-${i}`} className="claimfacin__nb-klausa-kartu">
          <div className="claimfacin__nb-klausa-identitas">
            <span className="claimfacin__nb-klausa-label">{K.clauseCode}</span>
            <code className="claimfacin__nb-kode">{k.clauseCode}</code>
            {k.clauseLanguageId !== '' && <span className="claimfacin__nb-klausa-jumlah">{k.clauseLanguageId}</span>}
            {adaArgumenNB(k.argumentCount) && (
              <span className="claimfacin__nb-klausa-jumlah">
                {K.totalArgumen}: {k.argumentCount}
              </span>
            )}
          </div>
          <div>
            <span className="claimfacin__nb-klausa-label">{K.title}</span>
            <strong>{k.clauseTitle}</strong>
          </div>
          <div>
            <span className="claimfacin__nb-klausa-label">{K.description}</span>
            <p className="claimfacin__nb-klausa-deskripsi">{k.clauseDescription}</p>
          </div>
        </article>
      ))}
    </div>
  )
}
