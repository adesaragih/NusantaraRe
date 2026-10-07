// Beranda — layar awal sesudah identitas, butir **bg**.
//
// ⚠️ `[kerangka aplikasi, bukan menu Pega]`. Isinya keputusan work owner, bukan tiruan sebuah rule: salam dan panel
// kotak masuk akun per workbasket (06-10-2026). Kartu tahap Claim Life, tabel Modul dan kartu Peran Anda DIBUANG
// (perintah work owner 06-10-2026: "buang aja, ga perlu") - modul dibuka dari sidebar.
//
// ⛔ Beranda tidak mengenal isi modul mana pun: angka dan daftar berkas datang dari kontrak `MenuModul.antreanBeranda`
// / `daftarBeranda` yang didaftarkan modulnya.

import { useEffect, useState } from 'react'

import { BERANDA } from '../inti/frontend/labels'
import { IkonCari } from '../inti/frontend/components/ui/dasar'
import { pesanGalat } from '../inti/frontend/klien'
import { type Sesi } from '../inti/frontend/store/sesi'
import { modulDipasang } from '../inti/frontend/lib/daftarMenu'
import type { AntreanBeranda, DaftarBeranda, ModulFrontend } from '../inti/frontend/modul'
import { MODUL_FRONTEND, type Halaman } from './daftar'

// Kotak masuk akun per workbasket (keputusan work owner 06-10-2026: "beranda menunjukkan berapa banyak case yang
// masuk di akun dia ... mengikuti workbasket"). Modul ikut dengan mendaftarkan `antreanBeranda` di PENDAFTARAN_MENU
// - Beranda tidak mengenal isi modul mana pun; modul tanpa penghitung tidak tampil di sini.

/** Modul AKTIF yang mendaftarkan penghitung kotak masuk. */
export function modulBerkotakMasuk(
  modul: readonly ModulFrontend<Halaman>[],
  aktif: readonly string[] | null = null,
): ModulFrontend<Halaman>[] {
  return modul.filter((m) => m.antreanBeranda !== undefined && modulDipasang(m.nama, aktif))
}

/** Kotak masuk satu modul yang sudah dimuat. */
export interface KotakModul {
  modul: ModulFrontend<Halaman>
  baris: AntreanBeranda[]
}

/** Banyak warna palet grafik kotak masuk (`.beranda__warna--0..5` di styles.css, dari token tema yang ada). */
export const WARNA_KOTAK_MASUK = 6

/** Satu batang panel kotak masuk: satu workbasket, dijumlah dari semua modul (jenis) yang memegangnya. */
export interface BatangKotakMasuk {
  /** Kode workbasket - kunci batang. */
  workbasket: string
  nama: string
  jumlah: number
  /** Lebar batang, persen terhadap workbasket terbanyak. */
  persen: number
  /** Porsi dari total berkas menunggu, persen bulat. */
  porsi: number
  /** Indeks warna palet, berputar menurut peringkat. */
  warna: number
}

/** Batang per workbasket, urut jumlah terbanyak (sama banyak: menurut nama). Jumlah NOL tidak tampil (perintah work
 *  owner 06-10-2026: "kalo jumlahnya 0 ga usah muncul"). Workbasket yang sama di beberapa modul = SATU batang; jenisnya
 *  dirinci `jenisKotakMasuk` sesudah batang diklik. */
export function batangKotakMasuk(kotak: readonly KotakModul[]): BatangKotakMasuk[] {
  const per = new Map<string, { workbasket: string; nama: string; jumlah: number }>()
  for (const { baris } of kotak) {
    for (const a of baris) {
      if (a.jumlah <= 0) continue
      const ada = per.get(a.workbasket)
      if (ada === undefined) per.set(a.workbasket, { workbasket: a.workbasket, nama: a.nama, jumlah: a.jumlah })
      else ada.jumlah += a.jumlah
    }
  }
  const semua = [...per.values()]
  const puncak = Math.max(0, ...semua.map((b) => b.jumlah))
  const total = semua.reduce((n, b) => n + b.jumlah, 0)
  return semua
    .sort((a, b) => b.jumlah - a.jumlah || (a.nama < b.nama ? -1 : a.nama > b.nama ? 1 : 0))
    .map((b, i) => ({
      ...b,
      persen: Math.round((b.jumlah / puncak) * 100),
      porsi: Math.round((b.jumlah / total) * 100),
      warna: i % WARNA_KOTAK_MASUK,
    }))
}

/** Cacah per jenis (modul) di SATU workbasket - kolom yang terbuka sesudah batangnya diklik (permintaan work owner
 *  06-10-2026: klik workbasket -> per jenis -> daftar berkas). Urut terbanyak; modul bernilai nol tidak tampil. */
export function jenisKotakMasuk(
  kotak: readonly KotakModul[],
  workbasket: string,
): { modul: ModulFrontend<Halaman>; jumlah: number }[] {
  return kotak
    .map(({ modul, baris }) => ({
      modul,
      jumlah: baris.filter((a) => a.workbasket === workbasket).reduce((n, a) => n + a.jumlah, 0),
    }))
    .filter((j) => j.jumlah > 0)
    .sort((a, b) => b.jumlah - a.jumlah)
}

/** Pilihan di panel kotak masuk: satu workbasket (atau semua - `workbasket` null) satu modul. */
export interface PilihanKotak {
  modul: ModulFrontend<Halaman>
  workbasket: string | null
  judul: string
}

/** Kunci pilihan - baris panel yang sedang terpilih. */
export function kunciPilihan(p: PilihanKotak | null): string {
  return p === null ? '' : `${p.modul.nama}|${p.workbasket ?? '*'}`
}

/** Baris daftar berkas yang cocok dengan filter (perintah work owner 06-10-2026: "tambahkan filter disini"): kueri
 *  dicari di SEMUA kolom yang tampil, huruf besar-kecil diabaikan; kueri kosong = semua baris. */
export function saringDaftar(daftar: DaftarBeranda, kueri: string): DaftarBeranda['baris'] {
  const q = kueri.trim().toLowerCase()
  if (q === '') return daftar.baris
  return daftar.baris.filter((r) => daftar.kolom.some((k) => (r.sel[k.kunci] ?? '').toLowerCase().includes(q)))
}

/** Pilihan panel yang diingat App selama sesi login: workbasket yang dibuka, jenis (daftar berkas), dan filter -
 *  tampil lagi sesudah Back dari layar kasus (perintah work owner 06-10-2026: "saat di back, ini jangan ilang"). */
export interface IngatanBeranda {
  bukaWb: string | null
  pilihan: PilihanKotak | null
  kueri: string
}

export default function Beranda({
  masuk,
  onBuka,
  onBukaKasus,
  modulAktif = null,
  ingatan,
  onIngat,
}: {
  masuk: Sesi
  onBuka: (modul: Halaman) => void
  /** Membuka SATU berkas di modulnya (daftar kotak masuk, `PropsRute.bukaKasus`). */
  onBukaKasus?: (modul: ModulFrontend<Halaman>, id: string) => void
  /** Modul aktif dari backend; `null` = semua (lihat `halamanAktif`). */
  modulAktif?: readonly string[] | null
  /** Pilihan panel terakhir - Beranda dibongkar selama layar kasus tampil, jadi App yang menyimpannya. */
  ingatan?: IngatanBeranda
  onIngat?: (i: IngatanBeranda) => void
}) {
  const [kotak, setKotak] = useState<KotakModul[]>([])
  // kotak masuk sudah dimuat: sebelum itu panel tidak tampil, sesudahnya nol berkas = panel berkata begitu
  const [kotakSiap, setKotakSiap] = useState(false)
  const [galatKotak, setGalatKotak] = useState<string | null>(null)
  // kunci isi daftar modul aktif - array baru tiap render tidak boleh memicu ulang permintaan
  const kunciAktif = modulAktif === null ? '*' : modulAktif.join(',')
  // daftar berkas pilihan panel kotak masuk (keputusan work owner 06-10-2026: tanpa masuk menu modul)
  const [pilihan, setPilihan] = useState<PilihanKotak | null>(ingatan?.pilihan ?? null)
  const [daftar, setDaftar] = useState<DaftarBeranda | null>(null)
  const [galatDaftar, setGalatDaftar] = useState<string | null>(null)
  // workbasket yang dibuka: kolom per jenis-nya tampil (permintaan work owner 06-10-2026)
  const [bukaWb, setBukaWb] = useState<string | null>(ingatan?.bukaWb ?? null)
  // filter daftar berkas - disaring di layar, daftarnya sudah utuh di tangan; dikosongkan saat daftar lain dibuka
  const [kueri, setKueri] = useState(ingatan?.kueri ?? '')

  // Laporkan pilihan ke App supaya tampil lagi sesudah Back dari layar kasus.
  useEffect(() => {
    onIngat?.({ bukaWb, pilihan, kueri })
  }, [bukaWb, pilihan, kueri, onIngat])

  // Daftar berkas dimuat ulang setiap kali tampil (juga sesudah Back): isinya bisa sudah berubah.
  useEffect(() => {
    setDaftar(null)
    setGalatDaftar(null)
    const muat = pilihan?.modul.daftarBeranda
    if (pilihan === null || muat === undefined) return
    let hidup = true
    void (async () => {
      try {
        const d = await muat(pilihan.workbasket)
        if (hidup) setDaftar(d)
      } catch (e) {
        if (hidup) setGalatDaftar(pesanGalat(e) ?? BERANDA.galatDaftar)
      }
    })()
    return () => {
      hidup = false
    }
  }, [pilihan])

  /** Klik baris panel: tampilkan daftar berkasnya di sini; klik lagi = tutup; modul tanpa daftar = buka modulnya. */
  const pilih = (p: PilihanKotak) => {
    if (p.modul.daftarBeranda === undefined) {
      onBuka(p.modul.halamanAwal)
      return
    }
    setPilihan((lama) => (kunciPilihan(lama) === kunciPilihan(p) ? null : p))
    setKueri('')
  }

  /** Klik batang workbasket: buka kolom per jenis-nya; klik lagi = tutup. Daftar berkas baru tampil sesudah jenis diklik. */
  const bukaWorkbasket = (workbasket: string) => {
    setBukaWb((lama) => (lama === workbasket ? null : workbasket))
    setPilihan(null)
    setKueri('')
  }

  useEffect(() => {
    let hidup = true
    const modul = modulBerkotakMasuk(MODUL_FRONTEND, kunciAktif === '*' ? null : kunciAktif.split(','))
    void (async () => {
      const hasil: KotakModul[] = []
      const gagal: string[] = []
      // berbarengan; satu modul yang gagal tidak menghapus kotak masuk modul lain
      await Promise.all(
        modul.map(async (m) => {
          try {
            const baris = (await m.antreanBeranda?.()) ?? []
            hasil.push({ modul: m, baris })
          } catch (e) {
            gagal.push(`${m.kelompok}: ${pesanGalat(e) ?? BERANDA.galatKotak}`)
          }
        }),
      )
      if (!hidup) return
      hasil.sort((a, b) => (a.modul.nama < b.modul.nama ? -1 : a.modul.nama > b.modul.nama ? 1 : 0))
      setKotak(hasil)
      setGalatKotak(gagal.length > 0 ? gagal.join(' · ') : null)
      setKotakSiap(true)
    })()
    return () => {
      hidup = false
    }
  }, [kunciAktif])

  const batang = batangKotakMasuk(kotak)
  const totalKotak = batang.reduce((n, b) => n + b.jumlah, 0)
  // workbasket yang sudah tak berisi sesudah kotak masuk dimuat ulang = kolom per jenis tertutup sendiri
  const wbBuka = batang.find((b) => b.workbasket === bukaWb)
  const tampil = daftar === null ? [] : saringDaftar(daftar, kueri)

  return (
    <div className="beranda">
      <section className="kepala-halaman">
        <div>
          <h2 className="beranda__judul">
            {BERANDA.salam}, {masuk.akunID}
          </h2>
        </div>
      </section>

      {galatKotak !== null && (
        <p className="alert alert--error" role="alert">
          {galatKotak}
        </p>
      )}

      {/* Panel kotak masuk akun (keputusan work owner 06-10-2026): angka total, batang proporsi bersegmen (pengganti donut
          dasbor lama), daftar peringkat per workbasket, cacah per jenis; jumlah nol tidak tampil; klik membuka modul. */}
      {kotakSiap && (
        <section className="kartu beranda__kotak" aria-labelledby="beranda-kotak">
          <div className="beranda__kotak-kepala">
            <div>
              <h3 id="beranda-kotak" className="beranda__kotak-judul">
                {BERANDA.kotakMasuk}
              </h3>
              <p className="beranda__kotak-total">
                <strong>{totalKotak}</strong>
                <span>{BERANDA.menunggu}</span>
              </p>
            </div>
          </div>
          {batang.length === 0 && <p className="beranda__daftar-catatan">{BERANDA.daftarKosong}</p>}
          {batang.length > 0 && (
            <div
              className="beranda__proporsi"
              role="img"
              aria-label={batang.map((b) => `${b.nama} ${b.jumlah}`).join(', ')}
            >
              {batang.map((b) => (
                <span
                  key={b.workbasket}
                  className={`beranda__proporsi-ruas beranda__warna--${b.warna}`}
                  style={{ flexGrow: b.jumlah }}
                  title={`${b.nama}: ${b.jumlah} (${b.porsi}%)`}
                />
              ))}
            </div>
          )}
          {batang.length > 0 && (
            <div className="beranda__kotak-isi">
              <div>
                <h4 className="beranda__kotak-subjudul">{BERANDA.perWorkbasket}</h4>
                <ul className="beranda__batang-daftar">
                  {batang.map((b) => (
                    <li key={b.workbasket}>
                      <button
                        type="button"
                        className={`beranda__batang beranda__warna--${b.warna}${
                          wbBuka?.workbasket === b.workbasket ? ' beranda__batang--aktif' : ''
                        }`}
                        aria-label={`${b.nama}: ${b.jumlah} ${BERANDA.antrean}, ${b.porsi}%`}
                        aria-expanded={wbBuka?.workbasket === b.workbasket}
                        onClick={() => {
                          bukaWorkbasket(b.workbasket)
                        }}
                      >
                        <span className="beranda__batang-kepala">
                          <span className="beranda__batang-titik" aria-hidden="true" />
                          <span className="beranda__batang-nama">{b.nama}</span>
                          <span className="beranda__batang-angka">{b.jumlah}</span>
                          <span className="beranda__batang-porsi">{b.porsi}%</span>
                        </span>
                        <span className="beranda__batang-jalur" aria-hidden="true">
                          <span className="beranda__batang-isi" style={{ width: `${b.persen}%` }} />
                        </span>
                      </button>
                    </li>
                  ))}
                </ul>
              </div>
              {/* Per jenis BARU tampil sesudah batang workbasket diklik, dan hanya jenis di workbasket itu; klik jenis =
                daftar berkasnya di bawah panel (permintaan work owner 06-10-2026). */}
              {wbBuka !== undefined && (
                <div className={`beranda__jenis-kolom beranda__warna--${wbBuka.warna}`}>
                  <h4 className="beranda__kotak-subjudul">{BERANDA.perJenis}</h4>
                  <ul className="beranda__jenis" aria-label={`${BERANDA.perJenis}: ${wbBuka.nama}`}>
                    {jenisKotakMasuk(kotak, wbBuka.workbasket).map((j) => {
                      const p: PilihanKotak = { modul: j.modul, workbasket: wbBuka.workbasket, judul: wbBuka.nama }
                      const aktif = kunciPilihan(pilihan) === kunciPilihan(p)
                      return (
                        <li key={j.modul.nama}>
                          <button
                            type="button"
                            className={`beranda__jenis-butir${aktif ? ' beranda__jenis-butir--aktif' : ''}`}
                            aria-pressed={aktif}
                            onClick={() => {
                              pilih(p)
                            }}
                          >
                            <span className="beranda__jenis-nama">{j.modul.kelompok}</span>
                            <span className="beranda__jenis-angka">{j.jumlah}</span>
                          </button>
                        </li>
                      )
                    })}
                  </ul>
                </div>
              )}
            </div>
          )}
        </section>
      )}

      {/* Daftar berkas pilihan panel kotak masuk: kolom dari modulnya; ID membuka layar kasus langsung. */}
      {pilihan !== null && (
        <section
          className="kartu beranda__daftar"
          aria-labelledby="beranda-daftar"
          aria-busy={daftar === null && galatDaftar === null}
        >
          <div className="kartu__kepala beranda__daftar-kepala">
            <h3 id="beranda-daftar">
              {pilihan.judul}
              {pilihan.workbasket !== null && (
                <span className="beranda__daftar-modul"> · {pilihan.modul.kelompok}</span>
              )}
            </h3>
            <div className="beranda__daftar-alat">
              {/* kapsul filter, saring saat mengetik (perintah work owner 06-10-2026: "tambahkan filter disini") */}
              <form className="beranda__saring" role="search" onSubmit={(e) => e.preventDefault()}>
                <IkonCari />
                <input
                  className="beranda__saring-isian"
                  aria-label={BERANDA.saring}
                  placeholder={BERANDA.saringPetunjuk}
                  value={kueri}
                  disabled={daftar === null || daftar.baris.length === 0}
                  onChange={(e) => setKueri(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Escape') setKueri('')
                  }}
                />
                <button
                  type="button"
                  className="beranda__saring-hapus"
                  aria-label={BERANDA.bersihkanSaring}
                  disabled={kueri === ''}
                  onClick={() => setKueri('')}
                >
                  ×
                </button>
              </form>
              <button
                type="button"
                className="beranda__daftar-tutup"
                aria-label={BERANDA.tutupDaftar}
                onClick={() => {
                  setPilihan(null)
                  setKueri('')
                }}
              >
                ×
              </button>
            </div>
          </div>
          {galatDaftar !== null && (
            <p className="alert alert--error" role="alert">
              {galatDaftar}
            </p>
          )}
          {daftar === null && galatDaftar === null && <p className="beranda__daftar-catatan">{BERANDA.memuatDaftar}</p>}
          {daftar !== null && daftar.baris.length === 0 && (
            <p className="beranda__daftar-catatan">{BERANDA.daftarKosong}</p>
          )}
          {daftar !== null && daftar.baris.length > 0 && tampil.length === 0 && (
            <p className="beranda__daftar-catatan">{BERANDA.tidakCocok}</p>
          )}
          {daftar !== null && tampil.length > 0 && (
            <div className="table-wrap beranda__daftar-tabel">
              <table>
                <thead>
                  <tr>
                    {daftar.kolom.map((k) => (
                      <th key={k.kunci} scope="col">
                        {k.label}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {tampil.map((r) => (
                    <tr key={r.id}>
                      {daftar.kolom.map((k, i) => (
                        <td key={k.kunci}>
                          {i === 0 && onBukaKasus !== undefined ? (
                            <button
                              type="button"
                              className="beranda__kartu-tautan"
                              onClick={() => onBukaKasus(pilihan.modul, r.id)}
                            >
                              {r.sel[k.kunci] ?? r.id}
                            </button>
                          ) : (
                            (r.sel[k.kunci] ?? '')
                          )}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      )}
    </div>
  )
}
