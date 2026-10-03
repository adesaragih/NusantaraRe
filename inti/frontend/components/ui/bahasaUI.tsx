/**
 * Bahasa teks BAWAAN komponen bersama — lihat `lib/teksUI.ts`.
 *
 * Bawaan `id`: layar yang tidak dibungkus Provider tidak berubah sama sekali.
 */
import { createContext, useContext } from 'react'

import { TEKS_UI, type Bahasa, type TeksUI } from '../../lib/teksUI'

export const BahasaUI = createContext<Bahasa>('id')

/** Bahasa layar ini. */
export function useBahasaUI(): Bahasa {
  return useContext(BahasaUI)
}

/** Teks bawaan komponen bersama dalam bahasa layar ini. */
export function useTeksUI(): TeksUI {
  return TEKS_UI[useContext(BahasaUI)]
}
