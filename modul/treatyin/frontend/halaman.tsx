// PENAMPUNG HALAMAN `TreatyIn` — padanan clipboard Pega untuk satu kontrak
// yang sedang dibuka.
//
// ---------------------------------------------------------------------
// ⭐ MENGAPA ADA (permintaan pemakai 7 Oktober 2026)
// ---------------------------------------------------------------------
// Di Pega seluruh tab satu kasus menulis ke SATU halaman `TreatyIn`
// (properti skalar + page list). Berpindah tab tidak membuang apa pun, dan
// tab yang satu membaca isian tab lain dari halaman yang sama — Share dan
// Achievement In IDR membaca `TreatyIn.Limits` yang tab Limits isi;
// Accumulation membaca `ReportingStart`/`ReportingEnd` tab Reporting
// Period.
//
// Layar ini dahulu menyimpan isian di `useState` TIAP TAB. Tab yang tidak
// tampil dilepas React, dan isiannya ikut hilang; tab lain tidak pernah
// melihatnya. Penampung ini hidup di form (`FormKontrakTreatyIn`), di atas
// strip tab, sehingga:
//
//   1. isian bertahan selama form terbuka, ke tab mana pun pemakai pindah;
//   2. properti yang sama dibaca dan ditulis SEMUA tab yang memakainya.
//
// ⛔ BUKAN tabel. Tidak ada yang menulis ke basis data dari sini; isinya
// masuk tabel saat tombol Save/Submit — pekerjaan berikutnya. Kuncinya
// ejaan properti Pega (`Portfolio`, `CoInScale[].CoInShare`, …) supaya
// pemetaan Save nanti satu lawan satu.
//
// ⛔ Tanpa penyedia (uji yang merender satu tab), `useProperti` jatuh ke
// `useState` biasa — perilaku tab sebelum penampung ada.

import { createContext, useContext, useEffect, useState, type Dispatch, type ReactNode, type SetStateAction } from 'react'

/** Isi halaman: properti Pega → nilai (teks, page list, atau baris). */
export type Halaman = Readonly<Record<string, unknown>>

/**
 * Properti halaman `TreatyIn` yang tab cabang PROPORSIONAL pegang — dari
 * `Section/TreatyInTabsProportional.xml` (+ `TreatyInShareProp`,
 * `TreatyInfoSubmit`), satu baris per tab. Properti yang muncul di DUA tab
 * adalah jalur ketergantungannya.
 */
export const PROPERTI_TAB_PROP: Readonly<Record<string, readonly string[]>> = {
  'Reporting Period': [
    'ReportingStart', 'ReportingEnd', 'ReportingPeriod', 'ReportingInterval',
    'ReportingSubmission', 'ReportingConfirmation', 'ReportingSettlement', 'ReportingPeriodList',
  ],
  Portfolio: ['Portfolio'],
  Limits: ['Limits'],
  Share: [
    'RNMShareP', 'BrokeragePercentP', 'OptionLimit',
    'Limits', 'TotalShareRnmProp', 'TotalSpreadedRnmProp', 'TotalSpreadedRnmRIProp',
  ],
  'Co-Ins Scale': ['CoInScale', 'MaxCoNonGroup', 'MaxCoGroup'],
  Accumulation: ['AccumulationPeriod', 'AccumulationList', 'ReportingStart', 'ReportingEnd'],
  Exclusions: ['ExclusionsP'],
  'Special Conditions': ['SpecialConditionsP'],
  'Information & Submit': ['Information', 'Comment'],
  'Achievement In IDR': ['Limits'],
}

/** Nilai satu properti — `undefined` bila belum pernah diisi. */
export function bacaProperti(h: Halaman, kunci: string): unknown {
  return Object.prototype.hasOwnProperty.call(h, kunci) ? h[kunci] : undefined
}

/** Halaman baru dengan satu properti diganti — fungsi murni. */
export function setelProperti(h: Halaman, kunci: string, nilai: unknown): Halaman {
  return { ...h, [kunci]: nilai }
}

/** Pengubah halaman yang form berikan ke tab. */
export interface PenampungHalaman {
  halaman: Halaman
  /** Ganti satu properti dari nilai TERKINI-nya (aman terhadap pembaruan beruntun). */
  ubah: (kunci: string, f: (lama: unknown) => unknown) => void
}

const KonteksHalaman = createContext<PenampungHalaman | null>(null)

/** Penyedia halaman — form memegang keadaannya. */
export function PenyediaHalaman({ penampung, children }: { penampung: PenampungHalaman; children: ReactNode }) {
  return <KonteksHalaman.Provider value={penampung}>{children}</KonteksHalaman.Provider>
}

/**
 * Penampung yang form pegang: `useState` atas `Halaman`. `kosongkan` dipakai
 * saat kontrak lain dimuat — halaman kasus baru, isian lama tidak terbawa.
 */
export function usePenampungHalaman(): PenampungHalaman & { kosongkan: () => void } {
  const [halaman, setHalaman] = useState<Halaman>({})
  return {
    halaman,
    ubah: (kunci, f) => {
      setHalaman((h) => setelProperti(h, kunci, f(bacaProperti(h, kunci))))
    },
    kosongkan: () => {
      setHalaman({})
    },
  }
}

/** Penampung yang sedang berlaku, atau `null` di luar form. */
export function usePenampung(): PenampungHalaman | null {
  return useContext(KonteksHalaman)
}

/**
 * Satu properti halaman, bentuknya seperti `useState`.
 *
 * `awal` = nilai saat properti itu BELUM ada di halaman (data kontrak yang
 * dimuat, atau kosong). Properti yang belum ada DISEMAI ke halaman sekali —
 * tab lain yang membacanya mendapat nilai yang sama, bukan turunan sendiri.
 */
export function useProperti<T>(kunci: string, awal: T | (() => T)): [T, Dispatch<SetStateAction<T>>] {
  const p = useContext(KonteksHalaman)
  const [lokal, setLokal] = useState<T>(awal)
  const ada = p !== null && Object.prototype.hasOwnProperty.call(p.halaman, kunci)
  // Penyemaian: hanya selama properti belum ada; sesudahnya efek ini diam.
  // `lokal` tetap nilai awal tab — `setLokal` tidak dipakai bila ada penyedia.
  useEffect(() => {
    if (p === null || ada) return
    p.ubah(kunci, (lama) => (lama === undefined ? lokal : lama))
  }, [p, ada, kunci, lokal])
  if (p === null) return [lokal, setLokal]
  const nilai = (ada ? p.halaman[kunci] : lokal) as T
  const setel: Dispatch<SetStateAction<T>> = (v) => {
    p.ubah(kunci, (lama) => {
      const kini = (lama === undefined ? lokal : lama) as T
      return typeof v === 'function' ? (v as (x: T) => T)(kini) : v
    })
  }
  return [nilai, setel]
}
