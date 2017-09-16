package crc32combine

/*
zlib CRC32 combine (https://github.com/madler/zlib)

Copyright (C) 1995-2006, 2010, 2011, 2012 Mark Adler <madler@alumni.caltech.edu>

This software is provided 'as-is', without any express or implied
warranty.  In no event will the authors be held liable for any damages
arising from the use of this software.
Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

1. The origin of this software must not be misrepresented; you must not
   claim that you wrote the original software. If you use this software
   in a product, an acknowledgment in the product documentation would be
   appreciated but is not required.

2. Altered source versions must be plainly marked as such, and must not be
   misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.

Ported from C to Go in 2016 by Justin Ruggles, with minimal alteration.
Used uint for unsigned long. Used uint32 for input arguments in order to match
the Go hash/crc32 package.
*/

func gf2MatrixTimes(mat []uint, vec uint) uint {
    var sum uint

    for vec != 0 {
        if vec&1 != 0 {
            sum ^= mat[0]
        }
        vec >>= 1
        mat = mat[1:]
    }
    return sum
}

func gf2MatrixSquare(square, mat []uint) {
    for n := 0; n < 32; n++ {
        square[n] = gf2MatrixTimes(mat, mat[n])
    }
}

func CRC32Combine(poly uint32, crc1, crc2 uint32, len2 int64) uint32 {
    /* degenerate case (also disallow negative lengths) */
    if len2 <= 0 {
        return crc1
    }

    even := make([]uint, 32) /* even-power-of-two zeros operator */
    odd := make([]uint, 32)  /* odd-power-of-two zeros operator */

    /* put operator for one zero bit in odd */
    odd[0] = uint(poly) /* CRC-32 polynomial */
    row := uint(1)
    for n := 1; n < 32; n++ {
        odd[n] = row
        row <<= 1
    }

    /* put operator for two zero bits in even */
    gf2MatrixSquare(even, odd)

    /* put operator for four zero bits in odd */
    gf2MatrixSquare(odd, even)

    /* apply len2 zeros to crc1 (first square will put the operator for one
       zero byte, eight zero bits, in even) */
    crc1n := uint(crc1)
    for {
        /* apply zeros operator for this bit of len2 */
        gf2MatrixSquare(even, odd)
        if len2&1 != 0 {
            crc1n = gf2MatrixTimes(even, crc1n)
        }
        len2 >>= 1

        /* if no more bits set, then done */
        if len2 == 0 {
            break
        }

        /* another iteration of the loop with odd and even swapped */
        gf2MatrixSquare(odd, even)
        if len2&1 != 0 {
            crc1n = gf2MatrixTimes(odd, crc1n)
        }
        len2 >>= 1

        /* if no more bits set, then done */
        if len2 == 0 {
            break
        }
    }

    /* return combined crc */
    crc1n ^= uint(crc2)
    return uint32(crc1n)
}
